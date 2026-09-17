#!/usr/bin/env python3
"""
Cleanup Unread — Bulk mark-as-read or delete unread emails
============================================================
Handles the unread email pile.

Usage:
  python3 cleanup_unread.py                          # Interactive
  python3 cleanup_unread.py --mark-read              # Mark all unread as read
  python3 cleanup_unread.py --mark-read --folder X   # Specific folder
  python3 cleanup_unread.py --delete-unread           # Delete all unread
  python3 cleanup_unread.py --delete-unread --days 90 # Delete unread >90 days
"""

import imaplib
import argparse
from datetime import datetime, timedelta, timezone

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


def get_unread_count(mail: imaplib.IMAP4_SSL, folder: str) -> int:
    """Count unread messages in folder."""
    st, data = mail.select(f'"{folder}"', readonly=True)
    if st != "OK":
        return 0
    st, unread_data = mail.search(None, "UNSEEN")
    if st != "OK" or not unread_data[0]:
        return 0
    return len(unread_data[0].split())


def find_unread(mail: imaplib.IMAP4_SSL, folder: str, older_than_days: int = 0) -> list[bytes]:
    """Find unread message IDs, optionally filtered by age."""
    st, _ = mail.select(f'"{folder}"')
    if st != "OK":
        return []

    if older_than_days > 0:
        cutoff = datetime.now(timezone.utc) - timedelta(days=older_than_days)
        date_str = cutoff.strftime("%d-%b-%Y")
        st, data = mail.search(None, f'UNSEEN BEFORE "{date_str}"')
    else:
        st, data = mail.search(None, "UNSEEN")

    if st != "OK" or not data[0]:
        return []

    return data[0].split()


def mark_as_read(mail: imaplib.IMAP4_SSL, folder: str, msg_ids: list[bytes] = None) -> int:
    """Mark messages as read. If msg_ids is None, marks ALL unread."""
    st, _ = mail.select(f'"{folder}"')
    if st != "OK":
        return 0

    if msg_ids is None:
        st, data = mail.search(None, "UNSEEN")
        if st != "OK" or not data[0]:
            return 0
        msg_ids = data[0].split()

    if not msg_ids:
        return 0

    BATCH = 200
    marked = 0
    for i in range(0, len(msg_ids), BATCH):
        batch = msg_ids[i:i + BATCH]
        id_str = b",".join(batch).decode()
        st, _ = mail.store(id_str, "+FLAGS", "\\Seen")
        if st == "OK":
            marked += len(batch)
        pct = min(100, int((i + BATCH) / len(msg_ids) * 100))
        print(f"\r  Marking as read in {folder}... {pct}%", end="", flush=True)

    print(f"\r  Marking as read in {folder}... ✓ {fmt_count(marked)} marked")
    return marked


def delete_unread(mail: imaplib.IMAP4_SSL, folder: str, msg_ids: list[bytes]) -> int:
    """Delete unread messages."""
    if not msg_ids:
        return 0

    st, _ = mail.select(f'"{folder}"')
    if st != "OK":
        return 0

    BATCH = 200
    deleted = 0
    for i in range(0, len(msg_ids), BATCH):
        batch = msg_ids[i:i + BATCH]
        id_str = b",".join(batch).decode()
        st, _ = mail.store(id_str, "+FLAGS", "\\Deleted")
        if st == "OK":
            deleted += len(batch)
        pct = min(100, int((i + BATCH) / len(msg_ids) * 100))
        print(f"\r  Deleting unread in {folder}... {pct}%", end="", flush=True)

    mail.expunge()
    print(f"\r  Deleting unread in {folder}... ✓ {fmt_count(deleted)} deleted")
    return deleted


def interactive_mode():
    """Interactive unread cleanup."""
    print("=" * 60)
    print(f"  📭 UNREAD CLEANUP — {ZOHO_EMAIL}")
    print("=" * 60)
    print()

    with imap_session() as mail:
        folders = get_all_folders(mail)

        # Show unread per folder
        print(f"{'Folder':<35} {'Unread':>8}")
        print("-" * 45)
        total_unread = 0
        folder_unreads = []
        for fname in folders:
            count = get_unread_count(mail, fname)
            if count > 0:
                print(f"  {fname:<33} {fmt_count(count):>8}")
                folder_unreads.append((fname, count))
                total_unread += count
        print("-" * 45)
        print(f"  {'TOTAL':<33} {fmt_count(total_unread):>8}")
        print()

        if total_unread == 0:
            print("  ✅ No unread emails!")
            return

        print("Options:")
        print("  1. Mark ALL unread as read (across all folders)")
        print("  2. Mark unread as read in INBOX only")
        print("  3. Delete unread emails older than N days")
        print("  4. Delete ALL unread emails (nuclear option)")
        print("  5. Quit")

        choice = input("\nChoice: ").strip()

        if choice == "1":
            if confirm(f"Mark {fmt_count(total_unread)} unread emails as read across all folders?"):
                total_marked = 0
                for fname, _ in folder_unreads:
                    total_marked += mark_as_read(mail, fname)
                print(f"\n  ✅ Marked {fmt_count(total_marked)} as read")

        elif choice == "2":
            inbox_unread = next((c for f, c in folder_unreads if f == "INBOX"), 0)
            if inbox_unread == 0:
                print("  No unread in INBOX")
                return
            if confirm(f"Mark {fmt_count(inbox_unread)} emails in INBOX as read?"):
                marked = mark_as_read(mail, "INBOX")
                print(f"\n  ✅ Marked {fmt_count(marked)} as read")

        elif choice == "3":
            days = input("Delete unread older than how many days? [90]: ").strip()
            days = int(days) if days else 90

            total_found = 0
            results = []
            for fname, _ in folder_unreads:
                ids = find_unread(mail, fname, older_than_days=days)
                if ids:
                    print(f"  [DRY RUN]  {fname}: {fmt_count(len(ids))} old unread emails")
                    results.append((fname, ids))
                    total_found += len(ids)

            if total_found == 0:
                print(f"  No unread emails older than {days} days")
                return

            if confirm(f"\nDelete {fmt_count(total_found)} old unread emails?"):
                for fname, ids in results:
                    delete_unread(mail, fname, ids)
                print(f"\n  ✅ Deleted {fmt_count(total_found)} old unread emails")

        elif choice == "4":
            if confirm(f"⚠️  DELETE {fmt_count(total_unread)} unread emails across ALL folders? This is irreversible!"):
                if confirm("  Are you SURE? Type 'y' to confirm permanent deletion"):
                    for fname, _ in folder_unreads:
                        ids = find_unread(mail, fname)
                        delete_unread(mail, fname, ids)
                    print(f"\n  ✅ Deleted all unread emails")

        elif choice == "5":
            return


def main():
    parser = argparse.ArgumentParser(description="Handle unread email bulk")
    parser.add_argument("--mark-read", action="store_true", help="Mark all unread as read")
    parser.add_argument("--delete-unread", action="store_true", help="Delete unread emails")
    parser.add_argument("--days", type=int, default=0, help="Only target unread older than N days")
    parser.add_argument("--folder", default=None, help="Specific folder (default: all)")
    parser.add_argument("--force", action="store_true", help="Skip confirmation")
    args = parser.parse_args()

    if not args.mark_read and not args.delete_unread:
        interactive_mode()
        return

    with imap_session() as mail:
        folders = [args.folder] if args.folder else get_all_folders(mail)

        if args.mark_read:
            total = 0
            for fname in folders:
                ids = find_unread(mail, fname, older_than_days=args.days)
                if ids:
                    print(f"  {fname}: {fmt_count(len(ids))} unread")
                    total += len(ids)

            if total == 0:
                print("No unread emails matching criteria.")
                return

            if not args.force and not confirm(f"Mark {fmt_count(total)} as read?"):
                print("Cancelled.")
                return

            for fname in folders:
                ids = find_unread(mail, fname, older_than_days=args.days)
                if ids:
                    mark_as_read(mail, fname, ids)
            print(f"\n✅ Done")

        elif args.delete_unread:
            total = 0
            for fname in folders:
                ids = find_unread(mail, fname, older_than_days=args.days)
                if ids:
                    print(f"  {fname}: {fmt_count(len(ids))} unread to delete")
                    total += len(ids)

            if total == 0:
                print("No unread emails matching criteria.")
                return

            if not args.force and not confirm(f"DELETE {fmt_count(total)} unread emails?"):
                print("Cancelled.")
                return

            for fname in folders:
                ids = find_unread(mail, fname, older_than_days=args.days)
                if ids:
                    delete_unread(mail, fname, ids)
            print(f"\n✅ Done")


if __name__ == "__main__":
    main()
