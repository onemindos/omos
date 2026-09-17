#!/usr/bin/env python3
"""
Deep Folder Scanner — Analyze email patterns in each folder
=============================================================
Goes into every folder and categorizes emails by sender patterns,
subject patterns, and content type. Helps decide what's trash vs
what's valuable for the Codex/mission.

Usage: python3 scan_deep.py
       python3 scan_deep.py --folder "INBOX"
       python3 scan_deep.py --folder "INBOX/25-49 HP (Holistic Performance)"
"""

import imaplib
import email
import email.utils
import email.header
import re
import argparse
import json
from collections import Counter, defaultdict
from datetime import datetime, timezone
from pathlib import Path

from config import imap_session, fmt_count, ZOHO_EMAIL

SCAN_OUTPUT = Path(__file__).parent / "scan_report.json"


# ---------------------------------------------------------------------------
# Pattern classification
# ---------------------------------------------------------------------------

# Obviously automated / no-reply / trash patterns
NO_REPLY_PATTERNS = [
    r"no[-_]?reply",
    r"noreply",
    r"do[-_]?not[-_]?reply",
    r"donotreply",
    r"mailer[-_]?daemon",
    r"postmaster",
    r"bounce",
    r"notification[s]?@",
    r"alert[s]?@",
    r"automated",
    r"auto[-_]?confirm",
]

NEWSLETTER_PATTERNS = [
    r"newsletter",
    r"digest",
    r"weekly[-_]?update",
    r"daily[-_]?brief",
    r"substack\.com",
    r"mailchimp",
    r"sendgrid",
    r"campaign",
    r"marketing",
    r"promo",
    r"offer[s]?@",
    r"deal[s]?@",
    r"sale[s]?@",
    r"news@",
    r"info@",
    r"updates@",
    r"hello@.*\.com",  # Generic hello@ from services
]

TRANSACTIONAL_PATTERNS = [
    r"receipt",
    r"invoice",
    r"order[-_]?confirm",
    r"shipping",
    r"tracking",
    r"billing",
    r"payment",
    r"subscription",
    r"verify",
    r"confirm.*email",
    r"password[-_]?reset",
    r"security[-_]?alert",
    r"two[-_]?factor",
    r"2fa",
    r"otp",
    r"login[-_]?attempt",
]

# Potentially valuable for Codex/mission
VALUABLE_KEYWORDS_SUBJECT = [
    r"proposal", r"partnership", r"collaboration", r"contract",
    r"invoice", r"meeting", r"schedule", r"call",
    r"onemind", r"codex", r"legacy", r"zeus",
    r"investment", r"business", r"opportunity",
    r"application", r"interview", r"offer letter",
    r"legal", r"documentation", r"agreement",
    r"attached", r"document", r"review",
]


def decode_header_value(raw: str) -> str:
    """Decode MIME encoded header value."""
    if not raw:
        return ""
    try:
        parts = email.header.decode_header(raw)
        decoded = []
        for part, charset in parts:
            if isinstance(part, bytes):
                decoded.append(part.decode(charset or "utf-8", errors="replace"))
            else:
                decoded.append(str(part))
        return " ".join(decoded)
    except Exception:
        return str(raw)


def classify_sender(addr: str) -> str:
    """Classify a sender address into a category."""
    addr_lower = addr.lower()
    for pattern in NO_REPLY_PATTERNS:
        if re.search(pattern, addr_lower):
            return "no-reply/automated"
    for pattern in NEWSLETTER_PATTERNS:
        if re.search(pattern, addr_lower):
            return "newsletter/marketing"
    for pattern in TRANSACTIONAL_PATTERNS:
        if re.search(pattern, addr_lower):
            return "transactional"
    return "other"


def classify_subject(subject: str) -> str:
    """Classify based on subject content."""
    subject_lower = subject.lower()
    for pattern in TRANSACTIONAL_PATTERNS:
        if re.search(pattern, subject_lower):
            return "transactional"
    for pattern in VALUABLE_KEYWORDS_SUBJECT:
        if re.search(pattern, subject_lower):
            return "potentially-valuable"
    if re.search(r"unsubscribe|opt[- ]?out|manage preferences", subject_lower):
        return "newsletter/marketing"
    return "uncategorized"


def extract_domain(addr: str) -> str:
    """Get domain from email address."""
    if "@" in addr:
        return addr.split("@")[-1].lower()
    return addr.lower()


def scan_folder(mail: imaplib.IMAP4_SSL, folder: str, max_fetch: int = 0) -> dict:
    """Deep scan a folder and return categorized analysis."""
    result = {
        "folder": folder,
        "total": 0,
        "unread": 0,
        "scanned": 0,
        "senders": Counter(),
        "domains": Counter(),
        "sender_categories": Counter(),
        "subject_categories": Counter(),
        "sample_subjects_by_sender": defaultdict(list),
        "no_reply_senders": Counter(),
        "newsletter_senders": Counter(),
        "valuable_senders": Counter(),
        "other_senders": Counter(),
    }

    st, data = mail.select(f'"{folder}"', readonly=True)
    if st != "OK":
        print(f"  ❌ Cannot select: {folder}")
        return result

    total = int(data[0])
    result["total"] = total
    if total == 0:
        return result

    # Get unread count
    st2, unread_data = mail.search(None, "UNSEEN")
    result["unread"] = len(unread_data[0].split()) if st2 == "OK" and unread_data[0] else 0

    # Determine fetch range
    fetch_count = max_fetch if max_fetch > 0 else total
    fetch_count = min(fetch_count, total)
    start = max(1, total - fetch_count + 1)
    msg_range = f"{start}:{total}"

    print(f"  Scanning {fmt_count(fetch_count)} of {fmt_count(total)} emails in {folder}...")

    # Fetch headers
    st, msgs = mail.fetch(msg_range, "(BODY.PEEK[HEADER.FIELDS (FROM SUBJECT DATE)])")
    if st != "OK":
        print(f"  ❌ Fetch failed")
        return result

    for item in msgs:
        if not isinstance(item, tuple):
            continue

        header_data = item[1].decode("utf-8", errors="replace") if isinstance(item[1], bytes) else item[1]
        msg = email.message_from_string(header_data)

        # Sender
        from_raw = msg.get("From", "")
        from_decoded = decode_header_value(from_raw)
        _, addr = email.utils.parseaddr(from_decoded)
        addr = addr.lower() if addr else from_decoded.lower()

        # Subject
        subject = decode_header_value(msg.get("Subject", "(no subject)"))

        # Classify
        sender_cat = classify_sender(addr)
        subject_cat = classify_subject(subject)
        domain = extract_domain(addr)

        result["senders"][addr] += 1
        result["domains"][domain] += 1
        result["sender_categories"][sender_cat] += 1
        result["subject_categories"][subject_cat] += 1

        # Bucket senders by category
        if sender_cat == "no-reply/automated":
            result["no_reply_senders"][addr] += 1
        elif sender_cat == "newsletter/marketing":
            result["newsletter_senders"][addr] += 1
        elif subject_cat == "potentially-valuable":
            result["valuable_senders"][addr] += 1
        else:
            result["other_senders"][addr] += 1

        # Collect sample subjects (up to 3 per sender)
        if len(result["sample_subjects_by_sender"][addr]) < 3:
            result["sample_subjects_by_sender"][addr].append(subject[:100])

        result["scanned"] += 1

    return result


def print_folder_report(scan: dict):
    """Pretty-print scan results for one folder."""
    folder = scan["folder"]
    total = scan["total"]
    scanned = scan["scanned"]
    unread = scan["unread"]

    print()
    print("━" * 70)
    print(f"  📂 {folder}")
    print(f"     Total: {fmt_count(total)} | Unread: {fmt_count(unread)} | Scanned: {fmt_count(scanned)}")
    print("━" * 70)

    # Category breakdown
    print(f"\n  📊 SENDER CATEGORIES:")
    for cat, count in scan["sender_categories"].most_common():
        pct = count / scanned * 100 if scanned > 0 else 0
        bar = "█" * int(pct / 2)
        print(f"     {cat:<25} {count:>6} ({pct:>5.1f}%) {bar}")

    # No-reply senders (TRASH candidates)
    if scan["no_reply_senders"]:
        print(f"\n  🗑  NO-REPLY / AUTOMATED ({sum(scan['no_reply_senders'].values())} emails):")
        for addr, count in scan["no_reply_senders"].most_common(20):
            subjects = scan["sample_subjects_by_sender"].get(addr, [])
            subj_preview = subjects[0][:60] if subjects else ""
            print(f"     {addr:<45} {count:>4}  │ {subj_preview}")

    # Newsletter senders (likely trash)
    if scan["newsletter_senders"]:
        print(f"\n  📰 NEWSLETTERS / MARKETING ({sum(scan['newsletter_senders'].values())} emails):")
        for addr, count in scan["newsletter_senders"].most_common(20):
            subjects = scan["sample_subjects_by_sender"].get(addr, [])
            subj_preview = subjects[0][:60] if subjects else ""
            print(f"     {addr:<45} {count:>4}  │ {subj_preview}")

    # Potentially valuable
    if scan["valuable_senders"]:
        print(f"\n  💎 POTENTIALLY VALUABLE ({sum(scan['valuable_senders'].values())} emails):")
        for addr, count in scan["valuable_senders"].most_common(20):
            subjects = scan["sample_subjects_by_sender"].get(addr, [])
            subj_preview = subjects[0][:60] if subjects else ""
            print(f"     {addr:<45} {count:>4}  │ {subj_preview}")

    # Other/uncategorized (need manual review)
    if scan["other_senders"]:
        print(f"\n  ❓ UNCATEGORIZED ({sum(scan['other_senders'].values())} emails — needs review):")
        for addr, count in scan["other_senders"].most_common(25):
            subjects = scan["sample_subjects_by_sender"].get(addr, [])
            subj_preview = subjects[0][:60] if subjects else ""
            print(f"     {addr:<45} {count:>4}  │ {subj_preview}")

    # Top domains
    print(f"\n  🌐 TOP DOMAINS:")
    for domain, count in scan["domains"].most_common(15):
        print(f"     {domain:<40} {count:>6}")


def get_all_folders(mail: imaplib.IMAP4_SSL) -> list[str]:
    """Get all folder names."""
    folders = []
    st, folder_list = mail.list()
    if st != "OK":
        return folders
    for item in folder_list:
        decoded = item.decode("utf-8") if isinstance(item, bytes) else item
        parts = decoded.rsplit('" "', 1)
        if len(parts) == 2:
            folders.append(parts[1].rstrip('"'))
        else:
            parts = decoded.rsplit(" ", 1)
            folders.append(parts[-1].strip('"'))
    return folders


def main():
    parser = argparse.ArgumentParser(description="Deep scan mail folders for patterns")
    parser.add_argument("--folder", help="Specific folder to scan (default: all non-empty)")
    parser.add_argument("--max", type=int, default=0, help="Max emails to scan per folder (0=all)")
    args = parser.parse_args()

    print("=" * 70)
    print(f"  🔬 DEEP FOLDER SCAN — {ZOHO_EMAIL}")
    print("=" * 70)

    all_scans = []

    with imap_session() as mail:
        if args.folder:
            target_folders = [args.folder]
        else:
            all_folders = get_all_folders(mail)
            # Skip Sent and folders with 0 messages
            target_folders = []
            for fname in all_folders:
                if fname.lower() in ("sent", "drafts"):
                    continue
                st, data = mail.select(f'"{fname}"', readonly=True)
                if st == "OK" and int(data[0]) > 0:
                    target_folders.append(fname)

        for folder in target_folders:
            scan = scan_folder(mail, folder, max_fetch=args.max)
            print_folder_report(scan)
            all_scans.append(scan)

    # Save full report
    serializable = []
    for scan in all_scans:
        serializable.append({
            "folder": scan["folder"],
            "total": scan["total"],
            "unread": scan["unread"],
            "scanned": scan["scanned"],
            "sender_categories": dict(scan["sender_categories"]),
            "no_reply_senders": dict(scan["no_reply_senders"].most_common()),
            "newsletter_senders": dict(scan["newsletter_senders"].most_common()),
            "valuable_senders": dict(scan["valuable_senders"].most_common()),
            "other_senders": dict(scan["other_senders"].most_common()),
            "top_domains": dict(scan["domains"].most_common(30)),
            "sample_subjects": {k: v for k, v in scan["sample_subjects_by_sender"].items()},
        })

    SCAN_OUTPUT.write_text(json.dumps(serializable, indent=2, default=str))
    print(f"\n💾 Full scan saved to: {SCAN_OUTPUT.name}")

    # Summary
    total_no_reply = sum(sum(s["no_reply_senders"].values()) for s in all_scans)
    total_newsletter = sum(sum(s["newsletter_senders"].values()) for s in all_scans)
    total_valuable = sum(sum(s["valuable_senders"].values()) for s in all_scans)
    total_other = sum(sum(s["other_senders"].values()) for s in all_scans)

    print()
    print("=" * 70)
    print("  📋 OVERALL SUMMARY")
    print("=" * 70)
    print(f"  🗑  No-reply/Automated:     {fmt_count(total_no_reply):>8}  ← Safe to move to delete")
    print(f"  📰 Newsletters/Marketing:   {fmt_count(total_newsletter):>8}  ← Review, mostly trash")
    print(f"  💎 Potentially Valuable:     {fmt_count(total_valuable):>8}  ← Keep / review")
    print(f"  ❓ Uncategorized:            {fmt_count(total_other):>8}  ← Needs manual review")
    print()
    print("  Next: python3 move_to_delete.py --no-reply     (move all no-reply)")
    print("        python3 move_to_delete.py --newsletters   (move newsletters)")
    print("        python3 move_to_delete.py --sender X      (move specific sender)")
    print()


if __name__ == "__main__":
    main()
