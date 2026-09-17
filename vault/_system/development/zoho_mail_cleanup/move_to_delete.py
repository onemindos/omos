#!/usr/bin/env python3
"""
Move to Delete — Safe bulk move to INBOX/delete folder
========================================================
Moves emails to the INBOX/delete staging folder instead of
permanently deleting them. Zeus does the final purge manually.

NEVER deletes directly — only MOVES to INBOX/delete.

Usage:
  python3 move_to_delete.py --no-reply                    # Move all no-reply/automated
  python3 move_to_delete.py --sender "news@example.com"   # Move from specific sender
  python3 move_to_delete.py --senders-file kill.txt        # Move from senders in file
  python3 move_to_delete.py --newsletters                  # Move newsletter senders
  python3 move_to_delete.py --pattern "unsubscribe"        # Move by subject pattern
  python3 move_to_delete.py --from-scan no-reply           # Use scan_report.json data
"""

import imaplib
import email
import email.utils
import email.header
import re
import argparse
import json
import sys
from pathlib import Path
from collections import Counter

from config import imap_session, fmt_count, confirm, ZOHO_EMAIL

DELETE_FOLDER = "INBOX/delete"
SCAN_REPORT = Path(__file__).parent / "scan_report.json"

# Same patterns from scan_deep.py
NO_REPLY_PATTERNS = [
    r"no[-_]?reply", r"noreply", r"do[-_]?not[-_]?reply", r"donotreply",
    r"mailer[-_]?daemon", r"postmaster", r"bounce",
    r"notification[s]?@", r"alert[s]?@", r"automated", r"auto[-_]?confirm",
]

NEWSLETTER_PATTERNS = [
    r"newsletter", r"digest", r"weekly[-_]?update", r"daily[-_]?brief",
    r"substack\.com", r"mailchimp", r"sendgrid", r"campaign",
    r"marketing", r"promo", r"offer[s]?@", r"deal[s]?@", r"sale[s]?@",
]


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


def search_by_sender(mail: imaplib.IMAP4_SSL, folder: str, sender: str) -> list[bytes]:
    """Find all emails from a specific sender in a folder."""
    st, _ = mail.select(f'"{folder}"', readonly=True)
    if st != "OK":
        return []
    st, data = mail.search(None, f'FROM "{sender}"')
    if st != "OK" or not data[0]:
        return []
    return data[0].split()


def search_by_subject_pattern(mail: imaplib.IMAP4_SSL, folder: str, pattern: str) -> list[bytes]:
    """Find emails matching a subject pattern (IMAP SUBJECT search)."""
    st, _ = mail.select(f'"{folder}"', readonly=True)
    if st != "OK":
        return []
    st, data = mail.search(None, f'SUBJECT "{pattern}"')
    if st != "OK" or not data[0]:
        return []
    return data[0].split()


def search_by_sender_pattern(mail: imaplib.IMAP4_SSL, folder: str, patterns: list[str]) -> list[bytes]:
    """Find emails from senders matching regex patterns.

    IMAP doesn't support regex in FROM search, so we fetch all FROM headers
    and filter locally. Slower but accurate.
    """
    st, data = mail.select(f'"{folder}"', readonly=True)
    if st != "OK":
        return []

    total = int(data[0])
    if total == 0:
        return []

    # Fetch all FROM headers
    st, msgs = mail.fetch(f"1:{total}", "(BODY.PEEK[HEADER.FIELDS (FROM)])")
    if st != "OK":
        return []

    matching_ids = []
    compiled = [re.compile(p, re.IGNORECASE) for p in patterns]

    for item in msgs:
        if not isinstance(item, tuple):
            continue

        # Extract message sequence number from response
        # Format: b'123 (BODY[HEADER.FIELDS (FROM)] {XX}'
        seq_match = re.match(rb"(\d+)", item[0])
        if not seq_match:
            continue
        seq_num = seq_match.group(1)

        header_data = item[1].decode("utf-8", errors="replace") if isinstance(item[1], bytes) else item[1]
        msg = email.message_from_string(header_data)
        from_raw = msg.get("From", "")
        _, addr = email.utils.parseaddr(from_raw)
        addr = addr.lower()

        if any(p.search(addr) for p in compiled):
            matching_ids.append(seq_num)

    return matching_ids


def move_to_delete(mail: imaplib.IMAP4_SSL, folder: str, msg_ids: list[bytes],
                   dry_run: bool = True) -> int:
    """Move messages from folder to INBOX/delete. Returns count moved."""
    if not msg_ids:
        return 0

    if folder == DELETE_FOLDER:
        return 0  # Don't move from delete to delete

    if dry_run:
        print(f"  [DRY RUN]  {folder}: {fmt_count(len(msg_ids))} → INBOX/delete")
        return len(msg_ids)

    # Select folder (writable)
    st, _ = mail.select(f'"{folder}"')
    if st != "OK":
        print(f"  ❌ Cannot select {folder}")
        return 0

    BATCH = 50
    moved = 0
    total = len(msg_ids)

    for i in range(0, total, BATCH):
        batch = msg_ids[i:i + BATCH]
        for mid in batch:
            mid_str = mid.decode() if isinstance(mid, bytes) else str(mid)
            try:
                st, _ = mail.copy(mid_str, f'"{DELETE_FOLDER}"')
                if st == "OK":
                    mail.store(mid_str, "+FLAGS", "\\Deleted")
                    moved += 1
            except Exception as e:
                print(f"\n  ⚠ Error moving msg {mid_str}: {e}")

        pct = min(100, int((i + BATCH) / total * 100))
        print(f"\r  Moving {folder} → INBOX/delete... {pct}% ({moved}/{total})", end="", flush=True)

    # Expunge the source folder (removes the copies we flagged as deleted)
    mail.expunge()
    print(f"\r  Moving {folder} → INBOX/delete... ✓ {fmt_count(moved)} moved        ")
    return moved


def get_senders_from_scan(category: str) -> list[str]:
    """Load sender list from scan_report.json by category."""
    if not SCAN_REPORT.exists():
        print("❌ No scan report. Run scan_deep.py first.")
        sys.exit(1)

    data = json.loads(SCAN_REPORT.read_text())
    senders = set()
    key_map = {
        "no-reply": "no_reply_senders",
        "newsletter": "newsletter_senders",
        "valuable": "valuable_senders",
        "other": "other_senders",
    }

    key = key_map.get(category)
    if not key:
        print(f"❌ Unknown category: {category}")
        print(f"   Valid: {', '.join(key_map.keys())}")
        sys.exit(1)

    for folder_scan in data:
        folder_senders = folder_scan.get(key, {})
        senders.update(folder_senders.keys())

    return sorted(senders)


def interactive_mode():
    """Interactive move-to-delete."""
    print("=" * 70)
    print(f"  📦 MOVE TO DELETE — {ZOHO_EMAIL}")
    print(f"     Target: {DELETE_FOLDER}")
    print("=" * 70)
    print()

    with imap_session() as mail:
        all_folders = [f for f in get_all_folders(mail) if f != DELETE_FOLDER]

        while True:
            print("\nActions:")
            print("  1. Move all no-reply/automated emails → delete")
            print("  2. Move all newsletters → delete")
            print("  3. Move specific sender → delete")
            print("  4. Move by subject keyword → delete")
            print("  5. Move from scan report category → delete")
            print("  6. Show current delete folder count")
            print("  0. Exit")

            choice = input("\nChoice: ").strip()

            if choice == "1":
                print(f"\n🔍 Scanning for no-reply/automated emails across all folders...")
                grand_total = 0
                folder_results = []

                for folder in all_folders:
                    ids = search_by_sender_pattern(mail, folder, NO_REPLY_PATTERNS)
                    if ids:
                        print(f"  [DRY RUN]  {folder}: {fmt_count(len(ids))} no-reply emails")
                        folder_results.append((folder, ids))
                        grand_total += len(ids)

                if grand_total == 0:
                    print("  No no-reply emails found.")
                    continue

                print(f"\n  Total: {fmt_count(grand_total)} no-reply emails to move")
                if confirm(f"  Move all {fmt_count(grand_total)} to {DELETE_FOLDER}?"):
                    for folder, ids in folder_results:
                        move_to_delete(mail, folder, ids, dry_run=False)
                    print(f"\n  ✅ Moved {fmt_count(grand_total)} no-reply emails to delete")

            elif choice == "2":
                print(f"\n🔍 Scanning for newsletter/marketing emails...")
                grand_total = 0
                folder_results = []

                for folder in all_folders:
                    ids = search_by_sender_pattern(mail, folder, NEWSLETTER_PATTERNS)
                    if ids:
                        print(f"  [DRY RUN]  {folder}: {fmt_count(len(ids))} newsletter emails")
                        folder_results.append((folder, ids))
                        grand_total += len(ids)

                if grand_total == 0:
                    print("  No newsletter emails found.")
                    continue

                print(f"\n  Total: {fmt_count(grand_total)} newsletter emails to move")
                if confirm(f"  Move all {fmt_count(grand_total)} to {DELETE_FOLDER}?"):
                    for folder, ids in folder_results:
                        move_to_delete(mail, folder, ids, dry_run=False)
                    print(f"\n  ✅ Moved {fmt_count(grand_total)} newsletter emails to delete")

            elif choice == "3":
                sender = input("Sender email address: ").strip()
                if not sender:
                    continue

                grand_total = 0
                folder_results = []
                for folder in all_folders:
                    ids = search_by_sender(mail, folder, sender)
                    if ids:
                        print(f"  {folder}: {fmt_count(len(ids))} from {sender}")
                        folder_results.append((folder, ids))
                        grand_total += len(ids)

                if grand_total == 0:
                    print(f"  No emails from {sender}")
                    continue

                print(f"\n  Total: {fmt_count(grand_total)} emails from {sender}")
                if confirm(f"  Move to {DELETE_FOLDER}?"):
                    for folder, ids in folder_results:
                        move_to_delete(mail, folder, ids, dry_run=False)
                    print(f"  ✅ Moved {fmt_count(grand_total)} emails from {sender}")

            elif choice == "4":
                keyword = input("Subject keyword: ").strip()
                if not keyword:
                    continue

                grand_total = 0
                folder_results = []
                for folder in all_folders:
                    ids = search_by_subject_pattern(mail, folder, keyword)
                    if ids:
                        print(f"  {folder}: {fmt_count(len(ids))} matching '{keyword}'")
                        folder_results.append((folder, ids))
                        grand_total += len(ids)

                if grand_total == 0:
                    print(f"  No emails matching '{keyword}'")
                    continue

                print(f"\n  Total: {fmt_count(grand_total)} emails matching '{keyword}'")
                if confirm(f"  Move to {DELETE_FOLDER}?"):
                    for folder, ids in folder_results:
                        move_to_delete(mail, folder, ids, dry_run=False)
                    print(f"  ✅ Moved {fmt_count(grand_total)}")

            elif choice == "5":
                category = input("Category (no-reply/newsletter/other): ").strip()
                senders = get_senders_from_scan(category)
                print(f"\n  Found {len(senders)} senders in '{category}' from scan:")
                for s in senders[:20]:
                    print(f"    → {s}")
                if len(senders) > 20:
                    print(f"    ... and {len(senders) - 20} more")

                if not confirm(f"\n  Search and move all emails from these {len(senders)} senders?"):
                    continue

                grand_total = 0
                for sender in senders:
                    for folder in all_folders:
                        ids = search_by_sender(mail, folder, sender)
                        if ids:
                            moved = move_to_delete(mail, folder, ids, dry_run=False)
                            grand_total += moved

                print(f"\n  ✅ Moved {fmt_count(grand_total)} emails from {len(senders)} senders")

            elif choice == "6":
                st, data = mail.select(f'"{DELETE_FOLDER}"', readonly=True)
                if st == "OK":
                    print(f"\n  📦 {DELETE_FOLDER}: {fmt_count(int(data[0]))} emails staged for deletion")
                else:
                    print(f"  ❌ Cannot read {DELETE_FOLDER}")

            elif choice == "0":
                break


def main():
    parser = argparse.ArgumentParser(description="Move emails to INBOX/delete (safe staging)")
    parser.add_argument("--no-reply", action="store_true", help="Move all no-reply/automated")
    parser.add_argument("--newsletters", action="store_true", help="Move all newsletters")
    parser.add_argument("--sender", help="Move from specific sender")
    parser.add_argument("--senders-file", help="File with sender emails (one per line)")
    parser.add_argument("--pattern", help="Move by subject keyword")
    parser.add_argument("--from-scan", help="Use scan report category (no-reply/newsletter/other)")
    parser.add_argument("--folder", default=None, help="Only search in this folder")
    parser.add_argument("--dry-run", action="store_true", help="Preview only")
    parser.add_argument("--force", action="store_true", help="Skip confirmation")
    args = parser.parse_args()

    # No args = interactive
    if not any([args.no_reply, args.newsletters, args.sender, args.senders_file,
                args.pattern, args.from_scan]):
        interactive_mode()
        return

    with imap_session() as mail:
        all_folders = [f for f in get_all_folders(mail) if f != DELETE_FOLDER]
        if args.folder:
            all_folders = [args.folder]

        grand_total = 0
        folder_results = []

        if args.no_reply:
            print("🔍 Finding no-reply/automated emails...")
            for folder in all_folders:
                ids = search_by_sender_pattern(mail, folder, NO_REPLY_PATTERNS)
                if ids:
                    folder_results.append((folder, ids))
                    grand_total += len(ids)

        elif args.newsletters:
            print("🔍 Finding newsletter/marketing emails...")
            for folder in all_folders:
                ids = search_by_sender_pattern(mail, folder, NEWSLETTER_PATTERNS)
                if ids:
                    folder_results.append((folder, ids))
                    grand_total += len(ids)

        elif args.sender:
            print(f"🔍 Finding emails from {args.sender}...")
            for folder in all_folders:
                ids = search_by_sender(mail, folder, args.sender)
                if ids:
                    folder_results.append((folder, ids))
                    grand_total += len(ids)

        elif args.senders_file:
            senders = [l.strip() for l in Path(args.senders_file).read_text().splitlines()
                       if l.strip() and not l.startswith("#")]
            print(f"🔍 Finding emails from {len(senders)} senders...")
            for sender in senders:
                for folder in all_folders:
                    ids = search_by_sender(mail, folder, sender)
                    if ids:
                        folder_results.append((folder, ids))
                        grand_total += len(ids)

        elif args.pattern:
            print(f'🔍 Finding emails with subject matching "{args.pattern}"...')
            for folder in all_folders:
                ids = search_by_subject_pattern(mail, folder, args.pattern)
                if ids:
                    folder_results.append((folder, ids))
                    grand_total += len(ids)

        elif args.from_scan:
            senders = get_senders_from_scan(args.from_scan)
            print(f"🔍 Finding emails from {len(senders)} {args.from_scan} senders...")
            for sender in senders:
                for folder in all_folders:
                    ids = search_by_sender(mail, folder, sender)
                    if ids:
                        folder_results.append((folder, ids))
                        grand_total += len(ids)

        # Preview
        for folder, ids in folder_results:
            print(f"  {folder}: {fmt_count(len(ids))}")
        print(f"\n  Total: {fmt_count(grand_total)} emails → {DELETE_FOLDER}")

        if args.dry_run:
            print("[DRY RUN] No changes made.")
            return

        if not args.force and not confirm(f"Move {fmt_count(grand_total)} emails to {DELETE_FOLDER}?"):
            print("Cancelled.")
            return

        moved_total = 0
        for folder, ids in folder_results:
            moved_total += move_to_delete(mail, folder, ids, dry_run=False)

        print(f"\n✅ Moved {fmt_count(moved_total)} emails to {DELETE_FOLDER}")


if __name__ == "__main__":
    main()
