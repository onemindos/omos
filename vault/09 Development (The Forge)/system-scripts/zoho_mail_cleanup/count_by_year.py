#!/usr/bin/env python3
"""Count emails by year in delete and Review folders."""
from config import imap_session

with imap_session() as mail:
    for folder in ["INBOX/delete", "INBOX/Review"]:
        st, d = mail.select(f'"{folder}"', readonly=True)
        if st != "OK":
            print(f"{folder}: cannot select")
            continue
        total = int(d[0])

        # 2025
        st, data = mail.search(None, "SINCE 01-Jan-2025 BEFORE 01-Jan-2026")
        y2025 = len(data[0].split()) if st == "OK" and data[0] else 0

        # Before 2025
        st, data = mail.search(None, "BEFORE 01-Jan-2025")
        older = len(data[0].split()) if st == "OK" and data[0] else 0

        # 2026
        st, data = mail.search(None, "SINCE 01-Jan-2026")
        y2026 = len(data[0].split()) if st == "OK" and data[0] else 0

        print(f"{folder} ({total} total):")
        print(f"  Before 2025: {older}")
        print(f"  2025:        {y2025}")
        print(f"  2026:        {y2026}")
        print()
