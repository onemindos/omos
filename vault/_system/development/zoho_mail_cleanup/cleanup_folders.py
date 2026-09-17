#!/usr/bin/env python3
"""
Cleanup Folders — Empty Trash, Spam, Junk, or any folder
==========================================================
Purges all messages from specified folders.

Usage:
  python3 cleanup_folders.py                        # Interactive
  python3 cleanup_folders.py --folder Trash          # Empty Trash
  python3 cleanup_folders.py --folder "Spam,Junk"    # Empty multiple
  python3 cleanup_folders.py --empty-all-trash       # All trash-like folders
"""

import imaplib
import argparse

from config import imap_session, fmt_count, confirm, ZOHO_EMAIL

TRASH_KEYWORDS = ("trash", "bin", "junk", "spam", "deleted", "bulk")


def get_all_folders(mail: imaplib.IMAP4_SSL) -> list[dict]:
    """Get folders with message counts."""
    folders = []
    st, folder_list = mail.list()
    if st != "OK":
        return folders

    for item in folder_list:
        decoded = item.decode("utf-8") if isinstance(item, bytes) else item
        parts = decoded.rsplit('" "', 1)
        if len(parts) == 2:
            name = parts[1].rstrip('"')
        else:
            parts = decoded.rsplit(" ", 1)
            name = parts[-1].strip('"')

        try:
            st, data = mail.select(f'"{name}"', readonly=True)
            count = int(data[0]) if st == "OK" else 0
        except:
            count = 0

        folders.append({"name": name, "total": count})

    return sorted(folders, key=lambda f: f["total"], reverse=True)


def empty_folder(mail: imaplib.IMAP4_SSL, folder: str, dry_run: bool = True) -> int:
    """Delete ALL messages in a folder."""
    st, data = mail.select(f'"{folder}"')
    if st != "OK":
        print(f"  ❌ Cannot select folder: {folder}")
        return 0

    total = int(data[0])
    if total == 0:
        print(f"  {folder}: already empty")
        return 0

    if dry_run:
        print(f"  [DRY RUN]  {folder}: {fmt_count(total)} messages")
        return total

    print(f"  Emptying {folder} ({fmt_count(total)} messages)...", end="", flush=True)

    # Select all messages
    st, msg_data = mail.search(None, "ALL")
    if st != "OK" or not msg_data[0]:
        print(" ❌ search failed")
        return 0

    msg_ids = msg_data[0].split()

    # Mark all as deleted in batches
    BATCH = 200
    deleted = 0
    for i in range(0, len(msg_ids), BATCH):
        batch = msg_ids[i:i + BATCH]
        id_str = b",".join(batch).decode()
        st, _ = mail.store(id_str, "+FLAGS", "\\Deleted")
        if st == "OK":
            deleted += len(batch)
        pct = min(100, int((i + BATCH) / len(msg_ids) * 100))
        print(f"\r  Emptying {folder} ({fmt_count(total)} messages)... {pct}%", end="", flush=True)

    mail.expunge()
    print(f"\r  Emptying {folder} ({fmt_count(total)} messages)... ✓ {fmt_count(deleted)} deleted")
    return deleted


def interactive_mode():
    """Interactive folder cleanup."""
    print("=" * 60)
    print(f"  🧹 FOLDER CLEANUP — {ZOHO_EMAIL}")
    print("=" * 60)
    print()

    with imap_session() as mail:
        folders = get_all_folders(mail)

        print(f"{'#':>3}  {'Folder':<35} {'Messages':>10}")
        print("-" * 55)
        for i, f in enumerate(folders, 1):
            marker = " 🗑" if any(k in f["name"].lower() for k in TRASH_KEYWORDS) else ""
            print(f"{i:>3}  {f['name']:<35} {fmt_count(f['total']):>10}{marker}")
        print()

        # Quick option for trash-like folders
        trash_folders = [f for f in folders if any(k in f["name"].lower() for k in TRASH_KEYWORDS) and f["total"] > 0]
        if trash_folders:
            trash_total = sum(f["total"] for f in trash_folders)
            print(f"  Quick option: Empty all trash/spam/junk ({fmt_count(trash_total)} messages total)")
            if confirm("  Empty all trash-like folders?"):
                for f in trash_folders:
                    empty_folder(mail, f["name"], dry_run=False)
                print(f"\n  ✅ Emptied {len(trash_folders)} folders")
                return

        # Manual selection
        print("\nEnter folder numbers to empty (comma-separated), or 'q' to quit:")
        picks = input("> ").strip()
        if picks.lower() == "q":
            return

        selected = []
        for p in picks.split(","):
            p = p.strip()
            if p.isdigit() and 1 <= int(p) <= len(folders):
                selected.append(folders[int(p) - 1])

        if not selected:
            print("No valid folders selected.")
            return

        # Preview
        total = 0
        for f in selected:
            total += empty_folder(mail, f["name"], dry_run=True)

        if confirm(f"\n  Delete {fmt_count(total)} messages from {len(selected)} folders?"):
            for f in selected:
                empty_folder(mail, f["name"], dry_run=False)
            print(f"\n  ✅ Folder cleanup complete")


def main():
    parser = argparse.ArgumentParser(description="Empty specific mail folders")
    parser.add_argument("--folder", help="Folder name(s), comma-separated")
    parser.add_argument("--empty-all-trash", action="store_true", help="Empty all trash/spam/junk folders")
    parser.add_argument("--dry-run", action="store_true", help="Preview only")
    parser.add_argument("--force", action="store_true", help="Skip confirmation")
    args = parser.parse_args()

    if not args.folder and not args.empty_all_trash:
        interactive_mode()
        return

    with imap_session() as mail:
        if args.empty_all_trash:
            all_folders = get_all_folders(mail)
            targets = [f for f in all_folders if any(k in f["name"].lower() for k in TRASH_KEYWORDS) and f["total"] > 0]
        else:
            targets = [{"name": f.strip(), "total": -1} for f in args.folder.split(",")]

        total = 0
        for f in targets:
            count = empty_folder(mail, f["name"], dry_run=True)
            total += count

        if total == 0:
            print("Nothing to clean.")
            return

        if args.dry_run:
            print(f"\n[DRY RUN] Would delete {fmt_count(total)} messages")
            return

        if not args.force and not confirm(f"Empty {len(targets)} folders ({fmt_count(total)} messages)?"):
            print("Cancelled.")
            return

        for f in targets:
            empty_folder(mail, f["name"], dry_run=False)
        print(f"\n✅ Done")


if __name__ == "__main__":
    main()
