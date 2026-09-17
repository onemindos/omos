#!/usr/bin/env python3
"""
Multi-Folder Scanner — Scan INBOX, HP, UI in one pass
======================================================
Fetches every email's sender, subject, date and categorizes.
Reports per-folder and combined summary.
"""
import imaplib
import email
import email.utils
import email.header
import re
import json
from collections import Counter
from config import imap_session

FOLDERS = [
    "INBOX",
    "INBOX/25-49 HP (Holistic Performance)",
    "INBOX/00-24 UI (Unified Intelligence)",
]

# ─── Classification patterns (reused from scan_review.py) ───

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

# Known keeper senders (from previous cleanup)
KEEPER_SENDERS = {
    "houseofgreen", "holisticfarming", "charlesdowding549359",
    "dralizapressman", "theparentingguide", "screenlessdad",
    "dremilyleeming", "vajenda", "joshturknett",
    "stoicwisdoms", "yungpueblo", "acornlandlabs", "peasantpolis",
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
    addr_lower = addr.lower()
    return any(re.search(p, addr_lower) for p in AUTOMATED_SENDER_PATTERNS)


def is_known_keeper(addr):
    addr = addr.lower()
    if addr in KEEPER_EXACT:
        return True
    if "@substack.com" in addr:
        base = addr.split("@")[0].split("+")[0]
        if base in KEEPER_SENDERS:
            return True
    return False


def classify_email(addr, name, subject):
    subj_lower = subject.lower()
    addr_lower = addr.lower()
    tags = []

    if any(re.search(p, subj_lower) for p in HOMESTEAD_PATTERNS):
        tags.append("homestead")
    if any(re.search(p, subj_lower) for p in EVERGREEN_PATTERNS):
        tags.append("evergreen")

    is_stale = any(re.search(p, subj_lower) for p in STALE_NEWS_SUBJECTS)
    if is_stale:
        tags.append("stale-news")

    is_auto = is_automated_sender(addr_lower)
    if not is_auto:
        tags.append("possibly-human")

    if is_known_keeper(addr_lower):
        tags.append("known-keeper")

    # Decision
    if "known-keeper" in tags:
        return "KEEPER", tags
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


def scan_folder(mail, folder_name):
    """Scan a single folder and return classified results."""
    st, d = mail.select(f'"{folder_name}"', readonly=True)
    if st != "OK":
        print(f"  ⚠️  Cannot select {folder_name}")
        return None
    total = int(d[0])
    if total == 0:
        print(f"  {folder_name}: empty")
        return None

    print(f"  📂 {folder_name}: {total} emails — fetching headers...")

    st, msgs = mail.fetch(
        f"1:{total}",
        "(BODY.PEEK[HEADER.FIELDS (FROM SUBJECT DATE)])"
    )
    if st != "OK":
        print(f"  ⚠️  Failed to fetch from {folder_name}")
        return None

    results = {
        "KEEPER": [], "KEEP-HOMESTEAD": [], "KEEP-HUMAN": [],
        "KEEP-EVERGREEN": [], "TRASH-STALE": [], "TRASH-AUTO": [], "REVIEW": [],
    }
    senders = Counter()

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
            "folder": folder_name,
        }
        results[category].append(entry)
        senders[addr] += 1
        parsed += 1

    print(f"  ✅ {folder_name}: {parsed} parsed")
    return results, senders, total


def main():
    print("=" * 70)
    print("  🔍 MULTI-FOLDER SCAN — INBOX + HP + UI")
    print("=" * 70)

    combined = {
        "KEEPER": [], "KEEP-HOMESTEAD": [], "KEEP-HUMAN": [],
        "KEEP-EVERGREEN": [], "TRASH-STALE": [], "TRASH-AUTO": [], "REVIEW": [],
    }
    all_senders = Counter()
    folder_stats = {}

    with imap_session() as mail:
        for folder in FOLDERS:
            result = scan_folder(mail, folder)
            if result is None:
                folder_stats[folder] = {"total": 0, "breakdown": {}}
                continue
            results, senders, total = result
            folder_stats[folder] = {
                "total": total,
                "breakdown": {cat: len(results[cat]) for cat in results},
            }
            for cat in results:
                combined[cat].extend(results[cat])
            all_senders.update(senders)

    grand_total = sum(fs["total"] for fs in folder_stats.values())

    # ─── PER-FOLDER BREAKDOWN ───
    print("\n" + "=" * 70)
    print("  📊 PER-FOLDER BREAKDOWN")
    print("=" * 70)
    cats = ["KEEPER", "KEEP-HOMESTEAD", "KEEP-HUMAN", "KEEP-EVERGREEN", "TRASH-STALE", "TRASH-AUTO", "REVIEW"]
    header = f"  {'Category':<20}" + "".join(f"{f:>10}" for f in FOLDERS) + f"{'TOTAL':>10}"
    print(header)
    print("  " + "-" * (20 + 10 * (len(FOLDERS) + 1)))
    for cat in cats:
        row = f"  {cat:<20}"
        cat_total = 0
        for f in FOLDERS:
            val = folder_stats.get(f, {}).get("breakdown", {}).get(cat, 0)
            row += f"{val:>10}"
            cat_total += val
        row += f"{cat_total:>10}"
        print(row)
    # totals row
    row = f"  {'TOTAL':<20}"
    for f in FOLDERS:
        row += f"{folder_stats.get(f, {}).get('total', 0):>10}"
    row += f"{grand_total:>10}"
    print(row)

    # ─── COMBINED KEEPER / HOMESTEAD ───
    keepers = combined["KEEPER"] + combined["KEEP-HOMESTEAD"]
    if keepers:
        print("\n" + "=" * 70)
        print(f"  🌟 KNOWN KEEPERS + HOMESTEAD ({len(keepers)} emails)")
        print("=" * 70)
        keeper_senders = Counter()
        for e in keepers:
            display = f"{e['name']} <{e['from']}>" if e['name'] else e['from']
            keeper_senders[display] += 1
        for sender, count in keeper_senders.most_common(30):
            print(f"  {count:>4}  {sender}")

    # ─── HUMAN SENDERS ───
    humans = combined["KEEP-HUMAN"]
    if humans:
        print("\n" + "=" * 70)
        print(f"  👤 HUMAN / NON-AUTOMATED SENDERS ({len(humans)} emails)")
        print("=" * 70)
        human_senders = Counter()
        for e in humans:
            display = f"{e['name']} <{e['from']}>" if e['name'] else e['from']
            human_senders[display] += 1
        for sender, count in human_senders.most_common(40):
            folder_info = ""
            # show which folder(s)
            folders_for = set(x["folder"] for x in humans if x["from"] == e["from"])
            print(f"  {count:>4}  {sender}")
        if len(human_senders) > 40:
            print(f"  ... and {len(human_senders) - 40} more")

    # ─── EVERGREEN ───
    evergreen = combined["KEEP-EVERGREEN"]
    if evergreen:
        print("\n" + "=" * 70)
        print(f"  📚 EVERGREEN CONTENT ({len(evergreen)} emails)")
        print("=" * 70)
        for e in evergreen[:25]:
            print(f"  [{e['folder']}] {e['subject'][:65]}")
            print(f"         from: {e['from']}")
        if len(evergreen) > 25:
            print(f"  ... and {len(evergreen) - 25} more")

    # ─── TOP TRASH SENDERS ───
    trash = combined["TRASH-STALE"] + combined["TRASH-AUTO"]
    if trash:
        print("\n" + "=" * 70)
        print(f"  🗑️  TOP TRASH (stale + automated) — {len(trash)} emails")
        print("=" * 70)
        trash_senders = Counter()
        for e in trash:
            trash_senders[e["from"]] += 1
        for sender, count in trash_senders.most_common(40):
            # which folders?
            folders_for = sorted(set(x["folder"] for x in trash if x["from"] == sender))
            folder_str = ",".join(folders_for)
            print(f"  {count:>4}  [{folder_str}] {sender}")
        if len(trash_senders) > 40:
            print(f"  ... and {len(trash_senders) - 40} more senders")

    # ─── UNCATEGORIZED ───
    uncategorized = combined["REVIEW"]
    if uncategorized:
        print("\n" + "=" * 70)
        print(f"  ❓ UNCATEGORIZED ({len(uncategorized)} emails)")
        print("=" * 70)
        for e in uncategorized[:30]:
            print(f"  [{e['folder']}] {e['subject'][:60]}")
            print(f"         from: {e['from']}  date: {e['date'][:16]}")
        if len(uncategorized) > 30:
            print(f"  ... and {len(uncategorized) - 30} more")

    # ─── SAVE ───
    outfile = "multi_scan_report.json"
    with open(outfile, "w") as f:
        json.dump({
            "folder_stats": folder_stats,
            "combined": {k: v for k, v in combined.items()},
        }, f, indent=2, default=str)
    print(f"\n  Full results saved to {outfile}")


if __name__ == "__main__":
    main()
