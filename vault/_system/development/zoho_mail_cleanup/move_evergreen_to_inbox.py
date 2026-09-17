#!/usr/bin/env python3
"""Move evergreen emails from INBOX/Review to INBOX."""
import imaplib
import email
import email.utils
import email.header
import re
from config import imap_session

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

HOMESTEAD_PATTERNS = [
    r"homestead", r"garden", r"permaculture", r"compost", r"soil",
    r"chicken", r"poultry", r"livestock", r"off.?grid", r"solar",
    r"rain.?water", r"harvest", r"seed", r"plant", r"grow",
    r"no.?dig", r"regenerat", r"farm", r"land", r"acre",
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
    r"morning\s+brew", r"evening\s+",
    r"monday|tuesday|wednesday|thursday|friday",
    r"january|february|march|april|may|june|july|august|september|october|november|december",
    r"\b\d{1,2}/\d{1,2}\b",
]


def decode_hdr(raw):
    if not raw:
        return ""
    parts = email.header.decode_header(raw)
    res = []
    for p, ch in parts:
        if isinstance(p, bytes):
            res.append(p.decode(ch or "utf-8", errors="replace"))
        else:
            res.append(str(p))
    return " ".join(res)


def main():
    with imap_session() as mail:
        st, d = mail.select('"INBOX/Review"')
        if st != "OK":
            print("Cannot select INBOX/Review")
            return
        total = int(d[0])
        print(f"Review folder: {total} emails")

        st, data = mail.uid(
            "FETCH", "1:*",
            "(UID BODY.PEEK[HEADER.FIELDS (FROM SUBJECT)])"
        )
        if st != "OK":
            print("Fetch failed")
            return

        evergreen_uids = []
        for item in data:
            if not isinstance(item, tuple):
                continue
            uid_match = re.search(rb"UID (\d+)", item[0])
            if not uid_match:
                continue
            uid = uid_match.group(1).decode()
            hdr = item[1].decode("utf-8", errors="replace") if isinstance(item[1], bytes) else item[1]
            msg = email.message_from_string(hdr)
            _, addr = email.utils.parseaddr(msg.get("From", ""))
            addr = addr.lower()
            subject = decode_hdr(msg.get("Subject", ""))
            subj_lower = subject.lower()

            is_evergreen = any(re.search(p, subj_lower) for p in EVERGREEN_PATTERNS)
            if not is_evergreen:
                continue

            is_homestead = any(re.search(p, subj_lower) for p in HOMESTEAD_PATTERNS)
            if is_homestead:
                continue

            is_stale = any(re.search(p, subj_lower) for p in STALE_NEWS_SUBJECTS)
            is_auto = any(re.search(p, addr) for p in AUTOMATED_SENDER_PATTERNS)
            if is_stale and is_auto:
                continue
            if not is_auto and not is_stale:
                continue

            evergreen_uids.append(uid)
            print(f"  -> {subject[:80]}")

        print(f"\nFound {len(evergreen_uids)} evergreen emails to move to INBOX")

        if not evergreen_uids:
            print("Nothing to move.")
            return

        # Copy in batches
        moved = 0
        batch_size = 50
        for i in range(0, len(evergreen_uids), batch_size):
            batch = evergreen_uids[i:i + batch_size]
            uid_set = ",".join(batch)
            st, _ = mail.uid("COPY", uid_set, "INBOX")
            if st == "OK":
                moved += len(batch)
                print(f"  Copied {moved}/{len(evergreen_uids)}...")
            else:
                print(f"  COPY failed for batch starting at {i}")

        print(f"\nDone — {moved} evergreen emails copied to INBOX")


if __name__ == "__main__":
    main()
