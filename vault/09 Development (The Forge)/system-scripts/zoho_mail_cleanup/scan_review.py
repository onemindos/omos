#!/usr/bin/env python3
"""
Deep Review Scanner — Analyze ALL emails in INBOX/Review
=========================================================
Fetches every email's sender, subject, date and categorizes into:
  - KEEP: Real humans, homesteading/evergreen, attachments, actionable
  - TRASH: Stale news, newsletters, automated, marketing
"""
import imaplib
import email
import email.utils
import email.header
import re
import json
import sys
from datetime import datetime
from collections import Counter, defaultdict
from config import imap_session

FOLDER = "INBOX/Review"

# ─── EVERGREEN TOPIC PATTERNS (worth keeping) ───
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

# Known human senders (not companies/newsletters)
HUMAN_SENDER_HINTS = [
    # No @substack, no @marketing, no noreply, etc.
    # We detect "human" by ABSENCE of automated patterns
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
    r"\b\d{1,2}/\d{1,2}\b",  # dates like 3/15
]


def decode_header_value(raw):
    """Decode an email header that might be encoded."""
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
    """Check if sender looks automated/company."""
    addr_lower = addr.lower()
    return any(re.search(p, addr_lower) for p in AUTOMATED_SENDER_PATTERNS)


def classify_email(sender_addr, sender_name, subject):
    """Classify a single email. Returns (category, tags)."""
    subj_lower = subject.lower()
    addr_lower = sender_addr.lower()
    tags = []

    # Check homesteading
    if any(re.search(p, subj_lower) for p in HOMESTEAD_PATTERNS):
        tags.append("homestead")

    # Check evergreen
    if any(re.search(p, subj_lower) for p in EVERGREEN_PATTERNS):
        tags.append("evergreen")

    # Check stale news
    is_stale = any(re.search(p, subj_lower) for p in STALE_NEWS_SUBJECTS)
    if is_stale:
        tags.append("stale-news")

    # Check if sender is automated
    is_auto = is_automated_sender(addr_lower)
    if not is_auto:
        tags.append("possibly-human")

    # Decision logic
    if "homestead" in tags:
        return "KEEP-HOMESTEAD", tags
    if "possibly-human" in tags and "stale-news" not in tags:
        return "KEEP-HUMAN", tags
    if "evergreen" in tags and "stale-news" not in tags:
        return "KEEP-EVERGREEN", tags
    if is_stale and is_auto:
        return "TRASH-STALE", tags
    if is_auto:
        return "TRASH-AUTO", tags

    return "REVIEW", tags


def main():
    print("=" * 70)
    print(f"  🔍 DEEP REVIEW SCAN — {FOLDER}")
    print("  Fetching ALL emails (sender + subject + date)...")
    print("=" * 70)

    with imap_session() as mail:
        st, d = mail.select(f'"{FOLDER}"', readonly=True)
        if st != "OK":
            print(f"Cannot select {FOLDER}")
            return
        total = int(d[0])
        print(f"\n  Total emails: {total}\n")

        # Fetch all headers in one go
        print("  Fetching headers... (this may take a minute)")
        st, msgs = mail.fetch(
            f"1:{total}",
            "(BODY.PEEK[HEADER.FIELDS (FROM SUBJECT DATE)])"
        )
        if st != "OK":
            print("  Failed to fetch headers")
            return

        # Parse all emails
        results = {
            "KEEP-HOMESTEAD": [],
            "KEEP-HUMAN": [],
            "KEEP-EVERGREEN": [],
            "TRASH-STALE": [],
            "TRASH-AUTO": [],
            "REVIEW": [],
        }
        all_senders = Counter()
        human_senders = Counter()
        homestead_subjects = []

        parsed = 0
        for item in msgs:
            if not isinstance(item, tuple):
                continue

            seq_match = re.match(rb"(\d+)", item[0])
            if not seq_match:
                continue
            seq_num = seq_match.group(1).decode()

            header_data = item[1].decode("utf-8", errors="replace") if isinstance(item[1], bytes) else item[1]
            msg = email.message_from_string(header_data)

            from_raw = msg.get("From", "")
            name, addr = email.utils.parseaddr(from_raw)
            name = decode_header_value(name) if name else ""
            addr = addr.lower()
            subject = decode_header_value(msg.get("Subject", "(no subject)"))
            date_str = msg.get("Date", "")

            category, tags = classify_email(addr, name, subject)

            entry = {
                "seq": seq_num,
                "from": addr,
                "name": name,
                "subject": subject[:120],
                "date": date_str[:30],
                "tags": tags,
            }
            results[category].append(entry)
            all_senders[addr] += 1

            if "possibly-human" in tags:
                display = f"{name} <{addr}>" if name else addr
                human_senders[display] += 1

            if "homestead" in tags:
                homestead_subjects.append(f"  {subject[:80]} — {addr}")

            parsed += 1
            if parsed % 500 == 0:
                print(f"  Parsed {parsed}/{total}...", flush=True)

        print(f"  Parsed {parsed}/{total} — done.\n")

        # ─── SUMMARY ───
        print("=" * 70)
        print("  📊 CLASSIFICATION RESULTS")
        print("=" * 70)
        for cat in ["KEEP-HOMESTEAD", "KEEP-HUMAN", "KEEP-EVERGREEN", "TRASH-STALE", "TRASH-AUTO", "REVIEW"]:
            print(f"  {cat:<20} {len(results[cat]):>6}")
        print(f"  {'TOTAL':<20} {parsed:>6}")

        # ─── HOMESTEAD ───
        print("\n" + "=" * 70)
        print("  🌱 HOMESTEAD / OFF-GRID CONTENT (Evergreen)")
        print("=" * 70)
        if homestead_subjects:
            for s in homestead_subjects[:50]:
                print(s)
            if len(homestead_subjects) > 50:
                print(f"  ... and {len(homestead_subjects) - 50} more")
        else:
            print("  None found with homestead keywords in subject line")

        # ─── HUMAN SENDERS ───
        print("\n" + "=" * 70)
        print("  👤 POSSIBLY HUMAN SENDERS (not automated)")
        print("=" * 70)
        for sender, count in human_senders.most_common(50):
            print(f"  {count:>4}  {sender}")
        if len(human_senders) > 50:
            print(f"  ... and {len(human_senders) - 50} more")

        # ─── EVERGREEN ───
        print("\n" + "=" * 70)
        print("  📚 EVERGREEN CONTENT (how-tos, guides, frameworks)")
        print("=" * 70)
        for e in results["KEEP-EVERGREEN"][:30]:
            print(f"  {e['subject'][:70]}")
            print(f"       from: {e['from']}")
        if len(results["KEEP-EVERGREEN"]) > 30:
            print(f"  ... and {len(results['KEEP-EVERGREEN']) - 30} more")

        # ─── TOP TRASH SENDERS ───
        print("\n" + "=" * 70)
        print("  🗑️  TOP TRASH SENDERS (stale + automated)")
        print("=" * 70)
        trash_senders = Counter()
        for cat in ["TRASH-STALE", "TRASH-AUTO"]:
            for e in results[cat]:
                trash_senders[e["from"]] += 1
        for sender, count in trash_senders.most_common(30):
            print(f"  {count:>4}  {sender}")

        # ─── UNCATEGORIZED SAMPLE ───
        print("\n" + "=" * 70)
        print("  ❓ UNCATEGORIZED SAMPLE (first 40)")
        print("=" * 70)
        for e in results["REVIEW"][:40]:
            print(f"  {e['subject'][:65]}")
            print(f"       from: {e['from']}  date: {e['date'][:16]}")

        # Save full results
        outfile = "review_deep_scan.json"
        with open(outfile, "w") as f:
            json.dump(results, f, indent=2, default=str)
        print(f"\n  Full results saved to {outfile}")


if __name__ == "__main__":
    main()
