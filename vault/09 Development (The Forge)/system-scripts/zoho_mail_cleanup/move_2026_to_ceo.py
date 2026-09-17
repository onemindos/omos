#!/usr/bin/env python3
"""Move 2026 emails from INBOX/Review to INBOX/CEO review."""
import imaplib
from config import imap_session

SOURCE = "INBOX/Review"
TARGET = "INBOX/CEO review"

with imap_session() as mail:
    st, d = mail.select(f'"{SOURCE}"')
    if st != "OK":
        print(f"Cannot select {SOURCE}")
        exit(1)

    total = int(d[0])
    print(f"{SOURCE}: {total} emails")

    # Find 2026 emails
    st, data = mail.search(None, "SINCE 01-Jan-2026")
    if st != "OK" or not data[0]:
        print("No 2026 emails found")
        exit(0)

    ids = data[0].split()
    print(f"Found {len(ids)} emails from 2026")
    print(f"Moving to {TARGET}...")

    moved = 0
    for i, mid in enumerate(ids):
        mid_str = mid.decode() if isinstance(mid, bytes) else str(mid)
        try:
            st, _ = mail.copy(mid_str, f'"{TARGET}"')
            if st == "OK":
                mail.store(mid_str, "+FLAGS", "\\Deleted")
                moved += 1
        except Exception as e:
            print(f"  Error on msg {mid_str}: {e}")

        if (i + 1) % 25 == 0 or i + 1 == len(ids):
            print(f"  {moved}/{len(ids)} moved...", flush=True)

    mail.expunge()
    print(f"\nDone. Moved {moved} emails from {SOURCE} to {TARGET}")
