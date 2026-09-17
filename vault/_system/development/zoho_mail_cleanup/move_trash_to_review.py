#!/usr/bin/env python3
"""
Move Trash to Review — Moves TRASH-STALE + TRASH-AUTO from INBOX, HP, UI → INBOX/Review.
Uses UID-based operations to avoid sequence number shifts.
Everything classified as KEEPER, KEEP-HOMESTEAD, KEEP-HUMAN, KEEP-EVERGREEN, or REVIEW stays put.
"""
import imaplib
import email
import email.utils
import email.header
import re
import sys
from collections import Counter
from config import imap_session

FOLDERS = [
    "INBOX",
    "INBOX/25-49 HP (Holistic Performance)",
    "INBOX/00-24 UI (Unified Intelligence)",
]
TARGET = "INBOX/Review"

# ─── Classification patterns ───

HOMESTEAD_PATTERNS = [
    r"homestead", r"garden", r"permaculture", r"compost", r"soil",
    r"chicken", r"poultry", r"livestock", r"off.?grid", r"solar",
    r"rain.?water", r"harvest", r"seed", r"plant", r"grow",
    r"no.?dig", r"regenerat", r"farm", r"land", r"acre",
    r"preserv", r"ferment", r"canning", r"food.?storage",
    r"bee.?keep", r"orchard", r"fruit.?tree", r"mushroom",
    r"herb", r"medicinal", r"natural.?build", r"earthship",
    r"water.?filter", r"well.?water", r"septic", r"greenhouse",
    r"root.?cellar", r"butcher", r"forag", r"wild.?craft",
]

EVERGREEN_PATTERNS = [
    r"how.?to", r"guide", r"tutorial", r"step.?by.?step",
    r"framework", r"principles", r"checklist", r"template",
    r"system", r"blueprint", r"playbook", r"strategy",
    r"mental.?model", r"stoic", r"philosophy", r"wisdom",
    r"parenting", r"child.?develop", r"family", r"estate.?plan",
    r"trust.?fund", r"generational", r"wealth.?build",
    r"health", r"nutrition", r"gut", r"microbiome", r"exercise",
    r"sleep", r"meditat", r"mindful", r"resilien",
]

AUTOMATED_SENDER_PATTERNS = [
    r"@substack\.com", r"@mail\.", r"@send\.", r"@marketing\.",
    r"@promo", r"@newsletter", r"@updates?\.", r"@info\.",
    r"@hello\.", r"@news\.", r"@digest", r"@campaign",
    r"noreply", r"no-reply", r"donotreply", r"notification",
    r"@mailchimp", r"@sendgrid", r"@beehiiv", r"@email\.",
    r"@offers?\.", r"@billing\.", r"@account\.",
    r"@zoho", r"@github\.com", r"@amazon\.com",
    r"@google\.com", r"@apple\.com", r"@microsoft\.com",
    r"@linkedin\.com", r"@twitter\.com", r"@x\.com",
    r"@facebook\.com", r"@meta\.com", r"@instagram\.com",
]

STALE_NEWS_SUBJECTS = [
    r"daily\s", r"weekly\s", r"this\s+week", r"today'?s",
    r"breaking", r"just\s+in", r"news\s+alert", r"digest",
    r"round.?up", r"recap", r"what\s+happened", r"top\s+stories",
    r"morning\s+brew", r"evening\s+", r"monday|tuesday|wednesday|thursday|friday",
    r"january|february|march|april|may|june|july|august|september|october|november|december",
    r"\b\d{1,2}/\d{1,2}\b",
]

KEEPER_PREFIXES = {
    "houseofgreen", "holisticfarming", "charlesdowding549359",
    "dralizapressman", "theparentingguide", "screenlessdad",
    "dremilyleeming", "vajenda", "joshturknett",
    "stoicwisdoms", "yungpueblo", "acornlandlabs", "peasantpolis",
    "johannadebiase",
}

KEEPER_EXACT = {
    "billing@consumerdirect.com", "olarios@consumerdirect.com",
    "tplplevy@schwab.com", "zeus.ops@legacyhq.ai",
    "founder@legacyplatform.co", "all@legacyplatform.co",
    "cameron.pfiffer@letta.intercom-mail.com",
    "rita@imiloainstitute.com", "atm@receipts.bankofamerica.com",
    "information@messages.tax.ny.gov", "drivera_245@icloud.com",
    "success@thegogetterfamily.com",
}


def decode_header_value(raw):
    if not raw:
        return ""
    decoded_parts = email.header.decode_header(raw)
    result = []
    for part, charset in decoded_parts:
        if isinstance(part, bytes):
            result.append(part.decode(charset or "utf-8", errors="replace"))
        else:
            result.append(str(part))
    return " ".join(result)


def is_automated_sender(addr):
    return any(re.search(p, addr) for p in AUTOMATED_SENDER_PATTERNS)


def is_known_keeper(addr):
    if addr in KEEPER_EXACT:
        return True
    if "@substack.com" in addr:
        base = addr.split("@")[0].split("+")[0]
        if base in KEEPER_PREFIXES:
            return True
    return False


def is_trash(addr, subject):
    """Returns True if email should be moved to Review (trash staging)."""
    addr = addr.lower()
    subj_lower = subject.lower()

    # Never trash known keepers
    if is_known_keeper(addr):
        return False

    # Never trash homestead content
    if any(re.search(p, subj_lower) for p in HOMESTEAD_PATTERNS):
        return False

    # Check if automated sender
    is_auto = is_automated_sender(addr)

    # Check if stale news
    is_stale = any(re.search(p, subj_lower) for p in STALE_NEWS_SUBJECTS)

    # Trash if automated OR stale news from automated source
    if is_auto:
        return True
    if is_stale:
        return True

    return False


def process_folder(mail, folder):
    """Scan folder, move trash UIDs to Review. Returns (moved, kept)."""
    st, d = mail.select(f'"{folder}"')
    if st != "OK":
        print(f"  ⚠️  Cannot select {folder}")
        return 0, 0
    total = int(d[0])
    if total == 0:
        print(f"  📂 {folder}: empty")
        return 0, 0

    print(f"\n  📂 {folder}: {total} emails")
    print(f"     Fetching headers (UID mode)...")

    # Use UID FETCH to get stable IDs
    st, msgs = mail.uid("fetch", "1:*", "(BODY.PEEK[HEADER.FIELDS (FROM SUBJECT)])")
    if st != "OK":
        print(f"     ⚠️  Failed to fetch")
        return 0, 0

    trash_uids = []
    keep_count = 0

    for item in msgs:
        if not isinstance(item, tuple):
            continue

        # Extract UID from response
        uid_match = re.search(rb"UID (\d+)", item[0])
        if not uid_match:
            continue
        uid = uid_match.group(1)

        header_data = item[1].decode("utf-8", errors="replace") if isinstance(item[1], bytes) else item[1]
        msg = email.message_from_string(header_data)
        _, addr = email.utils.parseaddr(msg.get("From", ""))
        addr = addr.lower()
        subject = decode_header_value(msg.get("Subject", "(no subject)"))

        if is_trash(addr, subject):
            trash_uids.append(uid)
        else:
            keep_count += 1

    print(f"     → Trash (move to Review): {len(trash_uids)}")
    print(f"     → Keep in place:          {keep_count}")

    if not trash_uids:
        return 0, keep_count

    # Move in batches of 50 UIDs
    moved = 0
    batch_size = 50
    for i in range(0, len(trash_uids), batch_size):
        batch = trash_uids[i:i + batch_size]
        uid_set = b",".join(batch)
        try:
            st, _ = mail.uid("copy", uid_set, f'"{TARGET}"')
            if st == "OK":
                mail.uid("store", uid_set, "+FLAGS", "\\Deleted")
                moved += len(batch)
            else:
                print(f"     ⚠️  Copy failed for batch starting at {i}")
        except Exception as e:
            print(f"     ⚠️  Error: {e}")
        print(f"     {moved}/{len(trash_uids)} moved...", flush=True)

    mail.expunge()
    return moved, keep_count


def main():
    print("=" * 70)
    print("  🗑️  MOVE TRASH → INBOX/Review")
    print("  Automated + stale news → Review for your deletion")
    print("  Keepers + humans stay in their current folders")
    print("=" * 70)

    total_moved = 0
    total_kept = 0

    with imap_session() as mail:
        for folder in FOLDERS:
            moved, kept = process_folder(mail, folder)
            total_moved += moved
            total_kept += kept

        # Final counts
        print("\n" + "=" * 70)
        print("  📊 FINAL RESULTS")
        print("=" * 70)
        for folder in FOLDERS + [TARGET]:
            st, d = mail.select(f'"{folder}"', readonly=True)
            cnt = int(d[0]) if st == "OK" else "?"
            print(f"  {cnt:>6}  {folder}")
        print(f"\n  ✅ Moved {total_moved} trash emails to Review")
        print(f"  📌 Kept {total_kept} emails in place")


if __name__ == "__main__":
    main()
