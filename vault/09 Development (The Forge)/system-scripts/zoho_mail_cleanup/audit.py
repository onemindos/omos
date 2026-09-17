#!/usr/bin/env python3
"""
Zoho Mail Audit — Full Mailbox Reconnaissance
==============================================
Scans all folders, counts emails, identifies top senders,
age distribution, and unread status. Run this first to
understand what you're working with.

Usage: python3 audit.py
"""

import imaplib
import email
import email.utils
from collections import Counter, defaultdict
from datetime import datetime, timezone
import json
import sys
from pathlib import Path

from config import imap_session, fmt_count, ZOHO_EMAIL

AUDIT_OUTPUT = Path(__file__).parent / "audit_report.json"


def list_folders(mail: imaplib.IMAP4_SSL) -> list[dict]:
    """Get all IMAP folders with message counts."""
    folders = []
    status, folder_list = mail.list()
    if status != "OK":
        print("❌ Failed to list folders")
        return folders

    for item in folder_list:
        # Parse folder name from IMAP response like: (\\HasNoChildren) "/" "INBOX"
        decoded = item.decode("utf-8") if isinstance(item, bytes) else item
        # Extract the folder name (last quoted string or last word)
        parts = decoded.rsplit('" "', 1)
        if len(parts) == 2:
            folder_name = parts[1].rstrip('"')
        else:
            parts = decoded.rsplit(" ", 1)
            folder_name = parts[-1].strip('"')

        # Get message count
        try:
            st, data = mail.select(f'"{folder_name}"', readonly=True)
            if st == "OK":
                count = int(data[0])
                # Get unread count
                st2, unread_data = mail.search(None, "UNSEEN")
                unread = len(unread_data[0].split()) if st2 == "OK" and unread_data[0] else 0
                folders.append({
                    "name": folder_name,
                    "total": count,
                    "unread": unread,
                })
        except Exception as e:
            folders.append({"name": folder_name, "total": -1, "unread": 0, "error": str(e)})

    return sorted(folders, key=lambda f: f["total"], reverse=True)


def sample_senders(mail: imaplib.IMAP4_SSL, folder: str, max_fetch: int = 500) -> list[str]:
    """Fetch sender addresses from a folder (samples up to max_fetch most recent)."""
    senders = []
    try:
        st, data = mail.select(f'"{folder}"', readonly=True)
        if st != "OK":
            return senders

        total = int(data[0])
        if total == 0:
            return senders

        # Fetch the most recent N messages (headers only)
        start = max(1, total - max_fetch + 1)
        msg_range = f"{start}:{total}"

        st, msgs = mail.fetch(msg_range, "(BODY.PEEK[HEADER.FIELDS (FROM DATE)])")
        if st != "OK":
            return senders

        for item in msgs:
            if isinstance(item, tuple):
                header_data = item[1].decode("utf-8", errors="replace") if isinstance(item[1], bytes) else item[1]
                msg = email.message_from_string(header_data)
                from_addr = msg.get("From", "")
                if from_addr:
                    # Extract just the email address
                    _, addr = email.utils.parseaddr(from_addr)
                    if addr:
                        senders.append(addr.lower())
    except Exception as e:
        print(f"  ⚠ Error sampling {folder}: {e}")

    return senders


def analyze_ages(mail: imaplib.IMAP4_SSL, folder: str, max_fetch: int = 500) -> dict:
    """Analyze email age distribution in a folder."""
    age_buckets = {"<7d": 0, "7-30d": 0, "1-3mo": 0, "3-6mo": 0, "6-12mo": 0, ">1yr": 0}
    now = datetime.now(timezone.utc)

    try:
        st, data = mail.select(f'"{folder}"', readonly=True)
        if st != "OK":
            return age_buckets

        total = int(data[0])
        if total == 0:
            return age_buckets

        start = max(1, total - max_fetch + 1)
        msg_range = f"{start}:{total}"

        st, msgs = mail.fetch(msg_range, "(BODY.PEEK[HEADER.FIELDS (DATE)])")
        if st != "OK":
            return age_buckets

        for item in msgs:
            if isinstance(item, tuple):
                header_data = item[1].decode("utf-8", errors="replace") if isinstance(item[1], bytes) else item[1]
                msg = email.message_from_string(header_data)
                date_str = msg.get("Date", "")
                if date_str:
                    try:
                        dt = email.utils.parsedate_to_datetime(date_str)
                        if dt.tzinfo is None:
                            dt = dt.replace(tzinfo=timezone.utc)
                        age_days = (now - dt).days
                        if age_days < 7:
                            age_buckets["<7d"] += 1
                        elif age_days < 30:
                            age_buckets["7-30d"] += 1
                        elif age_days < 90:
                            age_buckets["1-3mo"] += 1
                        elif age_days < 180:
                            age_buckets["3-6mo"] += 1
                        elif age_days < 365:
                            age_buckets["6-12mo"] += 1
                        else:
                            age_buckets[">1yr"] += 1
                    except Exception:
                        pass
    except Exception as e:
        print(f"  ⚠ Error analyzing ages in {folder}: {e}")

    return age_buckets


def run_audit():
    """Full mailbox audit."""
    print("=" * 60)
    print(f"  📬 ZOHO MAIL AUDIT — {ZOHO_EMAIL}")
    print("=" * 60)
    print()

    with imap_session() as mail:
        # --- STEP 1: Folder overview ---
        print("📁 Scanning folders...")
        folders = list_folders(mail)

        total_msgs = sum(f["total"] for f in folders if f["total"] > 0)
        total_unread = sum(f["unread"] for f in folders)

        print(f"\n{'Folder':<35} {'Total':>8} {'Unread':>8}")
        print("-" * 55)
        for f in folders:
            if f["total"] != 0:  # Skip empty folders
                unread_str = fmt_count(f["unread"]) if f["unread"] > 0 else "—"
                total_str = fmt_count(f["total"]) if f["total"] >= 0 else "ERROR"
                print(f"  {f['name']:<33} {total_str:>8} {unread_str:>8}")
        print("-" * 55)
        print(f"  {'TOTAL':<33} {fmt_count(total_msgs):>8} {fmt_count(total_unread):>8}")
        print()

        # --- STEP 2: Top senders (from INBOX) ---
        inbox_folders = [f for f in folders if f["name"].upper() in ("INBOX",) and f["total"] > 0]
        all_senders = Counter()

        for folder_info in folders:
            if folder_info["total"] > 0:
                fname = folder_info["name"]
                # Sample the biggest folders
                sample_size = min(500, folder_info["total"])
                print(f"📊 Sampling senders from {fname} ({sample_size} msgs)...")
                senders = sample_senders(mail, fname, max_fetch=sample_size)
                all_senders.update(senders)

        print(f"\n{'Top 30 Senders (sampled)':<45} {'Count':>6}")
        print("-" * 55)
        for addr, count in all_senders.most_common(30):
            print(f"  {addr:<43} {count:>6}")
        print()

        # --- STEP 3: Age distribution (INBOX) ---
        print("📅 Analyzing email ages in INBOX...")
        age_dist = analyze_ages(mail, "INBOX", max_fetch=1000)
        print(f"\n{'Age Bracket':<20} {'Count':>8}")
        print("-" * 30)
        for bracket, count in age_dist.items():
            bar = "█" * min(40, count // 5) if count > 0 else ""
            print(f"  {bracket:<18} {count:>8}  {bar}")
        print()

        # --- Save report ---
        report = {
            "email": ZOHO_EMAIL,
            "audit_date": datetime.now().isoformat(),
            "total_messages": total_msgs,
            "total_unread": total_unread,
            "folders": folders,
            "top_senders": dict(all_senders.most_common(100)),
            "inbox_age_distribution": age_dist,
        }
        AUDIT_OUTPUT.write_text(json.dumps(report, indent=2))
        print(f"💾 Full report saved to: {AUDIT_OUTPUT.name}")

        # --- Recommendations ---
        print("\n" + "=" * 60)
        print("  🎯 CLEANUP RECOMMENDATIONS")
        print("=" * 60)

        # Find newsletter/bulk senders (>10 emails sampled)
        bulk_senders = [(addr, c) for addr, c in all_senders.most_common(100) if c >= 10]
        if bulk_senders:
            print(f"\n  🗑  {len(bulk_senders)} bulk senders (10+ emails each):")
            for addr, c in bulk_senders[:15]:
                print(f"     → {addr} ({c} emails sampled)")

        if total_unread > 50:
            print(f"\n  📭  {fmt_count(total_unread)} unread emails — consider mark-all-read")

        old_count = age_dist.get(">1yr", 0)
        if old_count > 0:
            print(f"\n  🗓  {old_count} emails older than 1 year (sampled from INBOX)")

        trash_folders = [f for f in folders if any(k in f["name"].lower() for k in ("trash", "bin", "junk", "spam")) and f["total"] > 0]
        if trash_folders:
            print(f"\n  🧹  Folders to empty:")
            for f in trash_folders:
                print(f"     → {f['name']} ({fmt_count(f['total'])} messages)")

        print()
        print("Next steps — run specific cleanup scripts:")
        print("  python3 cleanup_senders.py    — Delete by sender")
        print("  python3 cleanup_old.py        — Delete old emails")
        print("  python3 cleanup_folders.py    — Empty folders")
        print("  python3 cleanup_unread.py     — Handle unread bulk")
        print("  python3 menu.py               — Interactive menu")
        print()


if __name__ == "__main__":
    run_audit()
