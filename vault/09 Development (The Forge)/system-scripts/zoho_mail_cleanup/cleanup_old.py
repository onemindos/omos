#!/usr/bin/env python3
"""
Cleanup Old Emails — Delete emails older than a threshold
==========================================================
Removes emails from specified folders that are older than N days.

Usage:
  python3 cleanup_old.py                           # Interactive
  python3 cleanup_old.py --days 365                # Delete >1yr old in INBOX
  python3 cleanup_old.py --days 90 --all-folders   # Delete >90d across all
  python3 cleanup_old.py --days 180 --dry-run      # Preview only
"""

import imaplib
import argparse
from datetime import datetime, timedelta, timezone
import sys

from config import imap_session, fmt_count, confirm, ZOHO_EMAIL


def get_all_folders(mail: imaplib.IMAP4_SSL) -> list[str]:
    """Get list of all folder names."""
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


def find_old_emails(mail: imaplib.IMAP4_SSL, folder: str, days: int) -> list[bytes]:
    """Find emails older than N days in a folder."""
    st, _ = mail.select(f'"{folder}"')
    if st != "OK":
        return []

    cutoff = datetime.now(timezone.utc) - timedelta(days=days)
    date_str = cutoff.strftime("%d-%b-%Y")  # IMAP date format

    st, data = mail.search(None, f'BEFORE "{date_str}"')
    if st != "OK" or not data[0]:
        return []

    return data[0].split()


def delete_old_in_folder(mail: imaplib.IMAP4_SSL, folder: str, days: int,
                         dry_run: bool = True) -> int:
    """Delete old emails in a folder. Returns count."""
    msg_ids = find_old_emails(mail, folder, days)
    count = len(msg_ids)

    if count == 0:
        return 0

    if dry_run:
        print(f"  [DRY RUN]  {folder}: {fmt_count(count)} emails older than {days} days")
        return count

    print(f"  Deleting {fmt_count(count)} old emails in {folder}...", end="", flush=True)

    # Select folder (writable)
    mail.select(f'"{folder}"')

    # Process in batches of 100
    BATCH = 100
    deleted = 0
    for i in range(0, len(msg_ids), BATCH):
        batch = msg_ids[i:i + BATCH]
        id_str = b",".join(batch).decode()
        st, _ = mail.store(id_str, "+FLAGS", "\\Deleted")
        if st == "OK":
            deleted += len(batch)
        # Progress
        pct = min(100, int((i + BATCH) / len(msg_ids) * 100))
        print(f"\r  Deleting {fmt_count(count)} old emails in {folder}... {pct}%", end="", flush=True)

    mail.expunge()
    print(f"\r  Deleting {fmt_count(count)} old emails in {folder}... ✓ {fmt_count(deleted)} deleted")
    return deleted


def interactive_mode():
    """Interactive old email cleanup."""
    print("=" * 60)
    print(f"  🗓  OLD EMAIL CLEANUP — {ZOHO_EMAIL}")
    print("=" * 60)
    print()

    with imap_session() as mail:
        all_folders = get_all_folders(mail)

        print("Available thresholds:")
        print("  1. Older than 30 days")
        print("  2. Older than 90 days (3 months)")
        print("  3. Older than 180 days (6 months)")
        print("  4. Older than 365 days (1 year)")
        print("  5. Custom")

        choice = input("\nChoose threshold: ").strip()
        days_map = {"1": 30, "2": 90, "3": 180, "4": 365}
        if choice in days_map:
            days = days_map[choice]
        elif choice == "5":
            days = int(input("Days: ").strip())
        else:
            print("Invalid choice")
            return

        print(f"\nTarget folders:")
        print("  1. INBOX only")
        print("  2. All folders")
        print("  3. Pick specific folders")

        fchoice = input("\nChoice: ").strip()
        if fchoice == "1":
            folders = ["INBOX"]
        elif fchoice == "2":
            folders = all_folders
        elif fchoice == "3":
            print("\nFolders:")
            for i, f in enumerate(all_folders, 1):
                print(f"  {i}. {f}")
            picks = input("Enter numbers (comma-separated): ").strip().split(",")
            folders = [all_folders[int(p.strip()) - 1] for p in picks if p.strip().isdigit()]
        else:
            return

        # Preview
        print(f"\n🔍 Scanning for emails older than {days} days...")
        grand_total = 0
        for folder in folders:
            count = delete_old_in_folder(mail, folder, days, dry_run=True)
            grand_total += count

        if grand_total == 0:
            print("  No old emails found!")
            return

        print(f"\n  Total: {fmt_count(grand_total)} emails to delete")
        if confirm(f"  Delete all {fmt_count(grand_total)} old emails?"):
            deleted_total = 0
            for folder in folders:
                deleted_total += delete_old_in_folder(mail, folder, days, dry_run=False)
            print(f"\n  ✅ Deleted {fmt_count(deleted_total)} old emails")


def main():
    parser = argparse.ArgumentParser(description="Delete emails older than N days")
    parser.add_argument("--days", type=int, help="Delete emails older than N days")
    parser.add_argument("--folder", default="INBOX", help="Folder (default: INBOX)")
    parser.add_argument("--all-folders", action="store_true", help="All folders")
    parser.add_argument("--dry-run", action="store_true", help="Preview only")
    parser.add_argument("--force", action="store_true", help="Skip confirmation")
    args = parser.parse_args()

    if not args.days:
        interactive_mode()
        return

    with imap_session() as mail:
        folders = get_all_folders(mail) if args.all_folders else [args.folder]

        # Preview
        grand_total = 0
        for folder in folders:
            count = delete_old_in_folder(mail, folder, args.days, dry_run=True)
            grand_total += count

        if grand_total == 0:
            print("No old emails found.")
            return

        print(f"\n{'[DRY RUN] ' if args.dry_run else ''}Total: {fmt_count(grand_total)} emails older than {args.days} days")

        if args.dry_run:
            return

        if not args.force and not confirm(f"Delete {fmt_count(grand_total)} emails?"):
            print("Cancelled.")
            return

        for folder in folders:
            delete_old_in_folder(mail, folder, args.days, dry_run=False)
        print(f"\n✅ Done")


if __name__ == "__main__":
    main()
