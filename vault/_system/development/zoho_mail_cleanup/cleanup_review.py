#!/usr/bin/env python3
"""
Move recommended keepers from INBOX/Review → INBOX/CEO review.
Everything else stays in Review — untouched.
"""
import imaplib
import email
import email.utils
import email.header
import re
from config import imap_session

SOURCE = "INBOX/Review"
TARGET = "INBOX/CEO review"

# Senders worth keeping — move to CEO review
KEEP_SENDERS = [
    # Homesteading (evergreen, LE domain) — including variant + addresses
    "houseofgreen@substack.com",
    "holisticfarming@substack.com",
    "holisticfarming+plant-profiles@substack.com",
    "charlesdowding549359@substack.com",
    "team@growyourownvegetables.org",
    "acornlandlabs@substack.com",
    "peasantpolis@substack.com",
    # Parenting (evergreen, LE domain)
    "dralizapressman@substack.com",
    "theparentingguide@substack.com",
    "screenlessdad@substack.com",
    # Financial records / billing
    "billing@consumerdirect.com",
    "olarios@consumerdirect.com",
    "tplplevy@schwab.com",
    # Your own emails
    "zeus.ops@legacyhq.ai",
    "founder@legacyplatform.co",
    "all@legacyplatform.co",
    # Letta AI (your stack)
    "cameron.pfiffer@letta.intercom-mail.com",
    # Retreats / education (LE)
    "rita@imiloainstitute.com",
    # Bank receipts / tax
    "atm@receipts.bankofamerica.com",
    "information@messages.tax.ny.gov",
    # Health (evergreen, HP domain)
    "dremilyleeming@substack.com",
    "vajenda@substack.com",
    "vajenda+hormone-therapy-from-a-to-z@substack.com",
    "vajenda+vulva-and-vagina@substack.com",
    "vajenda+supplements-and-such@substack.com",
    "joshturknett@substack.com",
    # Stoic / philosophy (evergreen, HP domain)
    "stoicwisdoms@substack.com",
    "stoicwisdoms+premium@substack.com",
    "yungpueblo@substack.com",
    # Credit repair client
    "drivera_245@icloud.com",
    # GoGetter Family (events / training)
    "success@thegogetterfamily.com",
]

# Build KEEP_SET with prefix matching for substack + variant addresses
KEEP_SET = {s.lower() for s in KEEP_SENDERS}

# For substack senders, also match any +variant address
KEEP_PREFIXES = set()
for s in KEEP_SENDERS:
    s = s.lower()
    if "@substack.com" in s:
        base = s.split("@")[0].split("+")[0]  # e.g. "johannadebiase"
        KEEP_PREFIXES.add(base)


def is_keeper(addr):
    """Check if an address matches KEEP_SET or a known prefix+variant."""
    addr = addr.lower()
    if addr in KEEP_SET:
        return True
    if "@substack.com" in addr:
        base = addr.split("@")[0].split("+")[0]
        if base in KEEP_PREFIXES:
            return True
    return False


def main():
    print("=" * 70)
    print("  📬 MOVE KEEPERS: Review → CEO review")
    print(f"  Only moving {len(KEEP_SET)} recommended senders")
    print(f"  Everything else stays in {SOURCE}")
    print("=" * 70)

    with imap_session() as mail:
        st, d = mail.select(f'"{SOURCE}"')
        if st != "OK":
            print(f"Cannot select {SOURCE}")
            return
        total = int(d[0])
        print(f"\n  Total in Review: {total}")
        print("  Fetching all FROM headers...")

        st, msgs = mail.fetch(f"1:{total}", "(BODY.PEEK[HEADER.FIELDS (FROM)])")
        if st != "OK":
            print("  Failed to fetch")
            return

        keep_ids = []
        skip_count = 0

        for item in msgs:
            if not isinstance(item, tuple):
                continue
            seq_match = re.match(rb"(\d+)", item[0])
            if not seq_match:
                continue
            seq_num = seq_match.group(1).decode()

            header = item[1].decode("utf-8", errors="replace") if isinstance(item[1], bytes) else item[1]
            msg = email.message_from_string(header)
            _, addr = email.utils.parseaddr(msg.get("From", ""))
            addr = addr.lower()

            if is_keeper(addr):
                keep_ids.append(seq_num)
            else:
                skip_count += 1

        print(f"\n  Moving to CEO review: {len(keep_ids)}")
        print(f"  Leaving in Review:    {skip_count}")

        if not keep_ids:
            print("  Nothing to move.")
            return

        # Move keepers to CEO review
        moved = 0
        for i, mid in enumerate(keep_ids):
            try:
                st, _ = mail.copy(mid, f'"{TARGET}"')
                if st == "OK":
                    mail.store(mid, "+FLAGS", "\\Deleted")
                    moved += 1
            except Exception as e:
                print(f"  Error on {mid}: {e}")
            if (i + 1) % 50 == 0 or i + 1 == len(keep_ids):
                print(f"  {moved}/{len(keep_ids)} moved...", flush=True)

        mail.expunge()

        # Final counts
        st, d = mail.select(f'"{SOURCE}"', readonly=True)
        review_now = int(d[0]) if st == "OK" else "?"
        st, d = mail.select(f'"{TARGET}"', readonly=True)
        ceo_now = int(d[0]) if st == "OK" else "?"

        print(f"\n  ✅ Done. Moved {moved} emails to CEO review")
        print(f"  📊 Review now:     {review_now}")
        print(f"  📊 CEO review now: {ceo_now}")


if __name__ == "__main__":
    main()
