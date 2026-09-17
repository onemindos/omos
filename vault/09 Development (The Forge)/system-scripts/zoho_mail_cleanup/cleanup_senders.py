#!/usr/bin/env python3
"""
Cleanup by Sender — Delete/move emails from specific senders
=============================================================
Can target individual senders or load bulk sender list from audit.

Usage:
  python3 cleanup_senders.py                         # Interactive mode
  python3 cleanup_senders.py --sender "news@x.com"   # Single sender
  python3 cleanup_senders.py --top 20                 # Top 20 from audit
  python3 cleanup_senders.py --file senders.txt       # From file (one per line)
"""

import imaplib
import argparse
import json
import sys
from pathlib import Path

from config import imap_session, fmt_count, confirm, ZOHO_EMAIL

AUDIT_FILE = Path(__file__).parent / "audit_report.json"
BATCH_SIZE = 100  # Process messages in batches


def find_emails_from_sender(mail: imaplib.IMAP4_SSL, folder: str, sender: str) -> list[bytes]:
    """Search folder for all emails from a specific sender."""
    st, _ = mail.select(f'"{folder}"')
    if st != "OK":
        return []

    st, data = mail.search(None, f'FROM "{sender}"')
    if st != "OK" or not data[0]:
        return []

    return data[0].split()


def delete_messages(mail: imaplib.IMAP4_SSL, folder: str, msg_ids: list[bytes],
                    action: str = "delete") -> int:
    """Delete or move messages by ID. Returns count of processed messages."""
    if not msg_ids:
        return 0

    st, _ = mail.select(f'"{folder}"')
    if st != "OK":
        return 0

    processed = 0
    # Process in batches
    for i in range(0, len(msg_ids), BATCH_SIZE):
        batch = msg_ids[i:i + BATCH_SIZE]
        id_str = b",".join(batch)

        if action == "delete":
            # Mark as deleted
            st, _ = mail.store(id_str.decode(), "+FLAGS", "\\Deleted")
            if st == "OK":
                processed += len(batch)
        elif action == "trash":
            # Move to Trash
            for mid in batch:
                st, _ = mail.copy(mid.decode(), "Trash")
                if st == "OK":
                    mail.store(mid.decode(), "+FLAGS", "\\Deleted")
                    processed += 1

    # Expunge to finalize deletions
    if processed > 0:
        mail.expunge()

    return processed


def cleanup_sender(mail: imaplib.IMAP4_SSL, sender: str, folders: list[str] = None,
                   dry_run: bool = True) -> int:
    """Clean up all emails from a single sender across folders."""
    if folders is None:
        folders = ["INBOX"]

    total = 0
    for folder in folders:
        msg_ids = find_emails_from_sender(mail, folder, sender)
        count = len(msg_ids)
        if count == 0:
            continue

        if dry_run:
            print(f"  [DRY RUN]  {folder}: {fmt_count(count)} emails from {sender}")
            total += count
        else:
            print(f"  Deleting {fmt_count(count)} emails from {sender} in {folder}...", end="", flush=True)
            deleted = delete_messages(mail, folder, msg_ids, action="delete")
            print(f" ✓ {fmt_count(deleted)} deleted")
            total += deleted

    return total


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


def interactive_mode():
    """Interactive sender cleanup with preview."""
    print("=" * 60)
    print(f"  🗑  SENDER CLEANUP — {ZOHO_EMAIL}")
    print("=" * 60)
    print()

    # Load audit data if available
    top_senders = {}
    if AUDIT_FILE.exists():
        audit = json.loads(AUDIT_FILE.read_text())
        top_senders = audit.get("top_senders", {})
        if top_senders:
            print("Top senders from last audit:")
            for i, (addr, count) in enumerate(sorted(top_senders.items(), key=lambda x: -x[1])[:30], 1):
                print(f"  {i:>3}. {addr:<45} ({count} emails)")
            print()

    with imap_session() as mail:
        all_folders = get_all_folders(mail)

        while True:
            print("\nOptions:")
            print("  1. Enter sender email to cleanup")
            print("  2. Cleanup top N senders from audit")
            print("  3. Quit")
            choice = input("\nChoice: ").strip()

            if choice == "1":
                sender = input("Sender email: ").strip().lower()
                if not sender:
                    continue

                # Preview across all folders
                print(f"\n🔍 Searching for emails from: {sender}")
                total = cleanup_sender(mail, sender, all_folders, dry_run=True)

                if total == 0:
                    print(f"  No emails found from {sender}")
                    continue

                print(f"\n  Total: {fmt_count(total)} emails from {sender}")
                if confirm(f"  Delete all {fmt_count(total)} emails from {sender}?"):
                    deleted = cleanup_sender(mail, sender, all_folders, dry_run=False)
                    print(f"  ✅ Deleted {fmt_count(deleted)} emails")

            elif choice == "2":
                if not top_senders:
                    print("  No audit data. Run audit.py first.")
                    continue

                n = input("How many top senders to clean? [10]: ").strip()
                n = int(n) if n else 10

                sorted_senders = sorted(top_senders.items(), key=lambda x: -x[1])[:n]
                print(f"\n🔍 Preview: cleaning top {n} senders...")

                to_delete = []
                for addr, audit_count in sorted_senders:
                    total = cleanup_sender(mail, addr, all_folders, dry_run=True)
                    if total > 0:
                        to_delete.append((addr, total))

                if not to_delete:
                    print("  Nothing to clean!")
                    continue

                grand_total = sum(c for _, c in to_delete)
                print(f"\n  Grand total: {fmt_count(grand_total)} emails from {len(to_delete)} senders")

                if confirm(f"  Delete all {fmt_count(grand_total)} emails?"):
                    for addr, _ in to_delete:
                        cleanup_sender(mail, addr, all_folders, dry_run=False)
                    print(f"  ✅ Bulk sender cleanup complete")

            elif choice == "3":
                break


def main():
    parser = argparse.ArgumentParser(description="Delete emails by sender")
    parser.add_argument("--sender", help="Single sender email to clean")
    parser.add_argument("--top", type=int, help="Clean top N senders from audit")
    parser.add_argument("--file", help="File with sender emails (one per line)")
    parser.add_argument("--folder", default="INBOX", help="Folder to clean (default: INBOX)")
    parser.add_argument("--all-folders", action="store_true", help="Clean across all folders")
    parser.add_argument("--dry-run", action="store_true", help="Preview only, don't delete")
    parser.add_argument("--force", action="store_true", help="Skip confirmation prompts")
    args = parser.parse_args()

    # No args = interactive mode
    if not args.sender and not args.top and not args.file:
        interactive_mode()
        return

    with imap_session() as mail:
        folders = get_all_folders(mail) if args.all_folders else [args.folder]

        senders = []
        if args.sender:
            senders = [args.sender]
        elif args.file:
            senders = [l.strip() for l in Path(args.file).read_text().splitlines() if l.strip() and not l.startswith("#")]
        elif args.top:
            if not AUDIT_FILE.exists():
                print("❌ No audit data. Run audit.py first.")
                sys.exit(1)
            audit = json.loads(AUDIT_FILE.read_text())
            sorted_s = sorted(audit.get("top_senders", {}).items(), key=lambda x: -x[1])
            senders = [addr for addr, _ in sorted_s[:args.top]]

        grand_total = 0
        for sender in senders:
            total = cleanup_sender(mail, sender, folders, dry_run=True)
            grand_total += total

        print(f"\n{'[DRY RUN] ' if args.dry_run else ''}Total: {fmt_count(grand_total)} emails to delete")

        if args.dry_run:
            return

        if not args.force and not confirm(f"Delete {fmt_count(grand_total)} emails?"):
            print("Cancelled.")
            return

        for sender in senders:
            cleanup_sender(mail, sender, folders, dry_run=False)
        print(f"\n✅ Done — {fmt_count(grand_total)} emails deleted")


if __name__ == "__main__":
    main()
