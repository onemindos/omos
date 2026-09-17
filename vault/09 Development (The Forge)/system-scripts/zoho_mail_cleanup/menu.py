#!/usr/bin/env python3
"""
Zoho Mail Cleanup — Interactive Menu
=====================================
Central hub for all cleanup operations.

Usage: python3 menu.py
"""

import subprocess
import sys
from pathlib import Path

SCRIPT_DIR = Path(__file__).parent


def run_script(name: str):
    """Run a script in this directory."""
    script = SCRIPT_DIR / name
    subprocess.run([sys.executable, str(script)])


def main():
    while True:
        print()
        print("=" * 60)
        print("  ⚡ ZOHO MAIL CLEANUP TOOLKIT")
        print("  zeus@onemindcodex.com")
        print("=" * 60)
        print()
        print("  1. 📊 Full Mailbox Audit        (audit.py)")
        print("  2. 🗑  Cleanup by Sender          (cleanup_senders.py)")
        print("  3. 🗓  Delete Old Emails          (cleanup_old.py)")
        print("  4. 🧹 Empty Folders (Trash/Spam)  (cleanup_folders.py)")
        print("  5. 📭 Unread Cleanup              (cleanup_unread.py)")
        print()
        print("  0. Exit")
        print()

        choice = input("  Choose operation: ").strip()

        if choice == "1":
            run_script("audit.py")
        elif choice == "2":
            run_script("cleanup_senders.py")
        elif choice == "3":
            run_script("cleanup_old.py")
        elif choice == "4":
            run_script("cleanup_folders.py")
        elif choice == "5":
            run_script("cleanup_unread.py")
        elif choice == "0":
            print("\n  👋 Done. Your inbox thanks you.\n")
            break
        else:
            print("  Invalid choice.")


if __name__ == "__main__":
    main()
