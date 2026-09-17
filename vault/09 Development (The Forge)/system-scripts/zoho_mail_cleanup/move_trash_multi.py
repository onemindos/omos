#!/usr/bin/env python3
"""
Move trash from INBOX, HP, UI to INBOX/Review.
Handles both auto-detected trash AND hidden newsletter senders.
"""
import imaplib
import email
import email.utils
import email.header
import re
from collections import Counter
from config import imap_session

FOLDERS = [
    "INBOX",
    "INBOX/25-49 HP (Holistic Performance)",
    "INBOX/00-24 UI (Unified Intelligence)",
]
TARGET = "INBOX/Review"

# ── Hidden newsletter senders (escaped auto-detection via custom domains) ──
HIDDEN_NEWSLETTER_SENDERS = {
    "astro@forwardfuture.ai",           # Forward Future / Matthew Berman
    "mike@mikeschmitz.com",             # Mike Schmitz (PKM marketing)
    "fintechtakes@workweek.com",        # Fintech Takes
    "dan@tldrnewsletter.com",           # ALL TLDR variants
    "hello@offers1.creditsesame.com",   # Credit Sesame
    "updates@simplilearnmailer.com",    # Simplilearn
    "admin@liamottley.com",             # Liam Ottley
    "scott@infinitusmail.com",          # Scott Oldford
    "mary@joinsecret.com",             # JoinSecret
    "matt@founderos.com",              # Matt Gray / Founderos
    "no_reply@learn.heygen.com",       # HeyGen marketing
    "hello@setapp.com",                # Setapp
    "hello@pausestudio.com",           # Pause
    "jack-roberts@mail.beehiiv.com",   # Beehiiv marketing
    "hi@mail.startuptycoon.com",       # Startup Tycoon
    "hello@deeplearning.ai",           # DeepLearning.AI
    "thebatch@deeplearning.ai",        # The Batch
    "hello@linkingyourthinking.com",   # Nick Milo / Linking Your Thinking
    "hi@retellai.com",                 # Retell AI
    "jakki.jakaj@tavily.com",          # Tavily
    "hello@ollama.com",                # Ollama updates
    "orders@bemightykind.com",         # Mighty Kind
    "drink@mudwtr.com",                # MudWtr
    "team@raycast.com",                # Raycast
    "support@relevanceai.com",         # Relevance AI
    "ant.wilson@supabase.com",         # Supabase updates
    "grace@info.clickup.com",          # ClickUp
    "team@mail.clickup.com",           # ClickUp
}

# ── Automated sender patterns (same as scanner) ──
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

# ── Keeper senders (DO NOT trash) ──
KEEPER_SENDERS = {
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
    "onlinebanking@ealerts.bankofamerica.com",
    "billing@companysage.com",
    "custservice@myfico.com",
    "servicecu@em.servicecu.org",
    "todd@advancedtaxsolutions.com",
    "ebay@ebay.com",
    "hello@gardencatering.com",
    "hello@mercury.com",
    "impact@imiloainstitute.com",
    "hello@circle.so",
}


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


def is_keeper(addr):
    addr = addr.lower()
    if addr in KEEPER_EXACT:
        return True
    if "@substack.com" in addr:
        base = addr.split("@")[0].split("+")[0]
        if base in KEEPER_SENDERS:
            return True
    return False


def is_auto(addr):
    return any(re.search(p, addr.lower()) for p in AUTOMATED_SENDER_PATTERNS)


def is_hidden_newsletter(addr):
    return addr.lower() in HIDDEN_NEWSLETTER_SENDERS


def main():
    print("=" * 70)
    print("  🗑️  MOVE TRASH → INBOX/Review")
    print("  Sources: INBOX, HP, UI")
    print("=" * 70)

    with imap_session() as mail:
        grand_moved = 0

        for folder in FOLDERS:
            st, d = mail.select(f'"{folder}"')
            if st != "OK":
                print(f"  ⚠️  Cannot select {folder}")
                continue
            total = int(d[0])
            if total == 0:
                print(f"  {folder}: empty — skip")
                continue

            print(f"\n  📂 {folder}: {total} emails")

            st, data = mail.uid(
                "FETCH", "1:*",
                "(UID BODY.PEEK[HEADER.FIELDS (FROM SUBJECT)])"
            )
            if st != "OK":
                print(f"  ⚠️  Fetch failed for {folder}")
                continue

            trash_uids = []
            reasons = Counter()

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

                # Never trash keepers
                if is_keeper(addr):
                    continue

                reason = None
                if is_auto(addr):
                    reason = "auto-detected"
                elif is_hidden_newsletter(addr):
                    reason = "hidden-newsletter"

                if reason:
                    trash_uids.append(uid)
                    reasons[reason] += 1

            print(f"  Found {len(trash_uids)} trash ({reasons.get('auto-detected', 0)} auto, {reasons.get('hidden-newsletter', 0)} hidden newsletters)")

            if not trash_uids:
                continue

            # Copy to Review in batches
            moved = 0
            batch_size = 50
            for i in range(0, len(trash_uids), batch_size):
                batch = trash_uids[i:i + batch_size]
                uid_set = ",".join(batch)
                st, _ = mail.uid("COPY", uid_set, f'"{TARGET}"')
                if st == "OK":
                    moved += len(batch)
                    print(f"    Copied {moved}/{len(trash_uids)}...")
                else:
                    print(f"    ⚠️  COPY failed for batch at {i}")

            grand_moved += moved
            print(f"  ✅ {folder}: {moved} moved to Review")

        print(f"\n{'=' * 70}")
        print(f"  TOTAL: {grand_moved} trash emails moved to INBOX/Review")
        print(f"{'=' * 70}")


if __name__ == "__main__":
    main()
