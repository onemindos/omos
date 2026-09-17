#!/usr/bin/env python3
"""
Email Bloat Audit — Daily volume + top sources across ALL folders.
Answers: how many emails/day, from where, and what to unsubscribe from.
"""
import imaplib
import email
import email.utils
import email.header
import re
import json
from datetime import datetime, timedelta
from collections import Counter, defaultdict
from config import imap_session


def decode_header_value(raw):
    if not raw:
        return ""
    try:
        decoded_parts = email.header.decode_header(raw)
        result = []
        for part, charset in decoded_parts:
            if isinstance(part, bytes):
                result.append(part.decode(charset or "utf-8", errors="replace"))
            else:
                result.append(str(part))
        return " ".join(result)
    except:
        return str(raw)


def parse_date(date_str):
    """Try to parse email date into a datetime."""
    if not date_str:
        return None
    try:
        parsed = email.utils.parsedate_to_datetime(date_str)
        return parsed
    except:
        return None


def extract_domain(addr):
    """Extract domain from email address."""
    if "@" in addr:
        return addr.split("@")[-1].lower()
    return addr.lower()


def main():
    print("=" * 70)
    print("  📊 EMAIL BLOAT AUDIT — Daily Volume & Top Sources")
    print("=" * 70)

    all_emails = []

    with imap_session() as mail:
        # List all folders
        st, folders_raw = mail.list()
        folder_names = []
        for f in folders_raw:
            fname = f.decode().split(' "/" ')[-1].strip('"')
            folder_names.append(fname)

        for folder in sorted(folder_names):
            st, d = mail.select(f'"{folder}"', readonly=True)
            if st != "OK":
                continue
            total = int(d[0])
            if total == 0:
                continue

            print(f"  📂 {folder}: {total} emails...", flush=True)

            st, msgs = mail.fetch(
                f"1:{total}",
                "(BODY.PEEK[HEADER.FIELDS (FROM DATE)])"
            )
            if st != "OK":
                continue

            for item in msgs:
                if not isinstance(item, tuple):
                    continue
                header_data = item[1].decode("utf-8", errors="replace") if isinstance(item[1], bytes) else item[1]
                msg = email.message_from_string(header_data)

                from_raw = msg.get("From", "")
                name, addr = email.utils.parseaddr(from_raw)
                name = decode_header_value(name) if name else ""
                addr = addr.lower()
                domain = extract_domain(addr)
                date_str = msg.get("Date", "")
                dt = parse_date(date_str)

                all_emails.append({
                    "addr": addr,
                    "name": name,
                    "domain": domain,
                    "date": dt,
                    "folder": folder,
                })

    print(f"\n  Total emails scanned: {len(all_emails)}")

    # ─── DATE RANGE ───
    dated = [e for e in all_emails if e["date"]]
    if not dated:
        print("  No parseable dates found!")
        return

    dates_only = [e["date"].date() for e in dated]
    min_date = min(dates_only)
    max_date = max(dates_only)
    span_days = (max_date - min_date).days + 1

    print(f"  Date range: {min_date} → {max_date} ({span_days} days)")
    print(f"  Average: {len(dated) / span_days:.1f} emails/day")

    # ─── EMAILS PER DAY (last 30 days) ───
    print("\n" + "=" * 70)
    print("  📅 EMAILS PER DAY (last 30 days)")
    print("=" * 70)
    today = max_date
    daily = Counter()
    for e in dated:
        d = e["date"].date()
        if (today - d).days <= 30:
            daily[d] += 1

    for d in sorted(daily.keys()):
        bar = "█" * min(daily[d], 80)
        weekday = d.strftime("%a")
        print(f"  {d} {weekday} {daily[d]:>4}  {bar}")

    if daily:
        vals = list(daily.values())
        avg_30 = sum(vals) / len(vals)
        print(f"\n  Last 30 days avg: {avg_30:.1f} emails/day")
        print(f"  Peak: {max(vals)} emails on {max(daily, key=daily.get)}")

    # ─── EMAILS PER WEEK (last 12 weeks) ───
    print("\n" + "=" * 70)
    print("  📆 EMAILS PER WEEK (last 12 weeks)")
    print("=" * 70)
    weekly = Counter()
    for e in dated:
        d = e["date"].date()
        if (today - d).days <= 84:
            week_start = d - timedelta(days=d.weekday())
            weekly[week_start] += 1
    for w in sorted(weekly.keys()):
        bar = "█" * (weekly[w] // 5)
        print(f"  Week of {w}: {weekly[w]:>4}  {bar}")

    # ─── TOP DOMAINS (biggest bloat sources) ───
    print("\n" + "=" * 70)
    print("  🌐 TOP DOMAINS BY VOLUME (all-time)")
    print("=" * 70)
    domain_counter = Counter(e["domain"] for e in all_emails)
    domain_daily = {}
    for dom, cnt in domain_counter.most_common(40):
        per_day = cnt / span_days
        domain_daily[dom] = per_day
        bar = "█" * min(cnt // 2, 50)
        print(f"  {cnt:>5} ({per_day:>4.1f}/day)  {dom:<45} {bar}")

    # ─── TOP SENDERS (individual addresses) ───
    print("\n" + "=" * 70)
    print("  👤 TOP SENDERS BY VOLUME (all-time)")
    print("=" * 70)
    sender_counter = Counter()
    sender_names = {}
    for e in all_emails:
        sender_counter[e["addr"]] += 1
        if e["name"]:
            sender_names[e["addr"]] = e["name"]

    for addr, cnt in sender_counter.most_common(50):
        per_day = cnt / span_days
        name = sender_names.get(addr, "")
        display = f"{name} <{addr}>" if name else addr
        bar = "█" * min(cnt // 2, 30)
        print(f"  {cnt:>5} ({per_day:>4.1f}/day)  {display[:65]:<65} {bar}")

    # ─── SUBSTACK BREAKDOWN ───
    substack = [e for e in all_emails if "substack" in e["domain"]]
    print("\n" + "=" * 70)
    print(f"  📰 SUBSTACK BREAKDOWN — {len(substack)} emails ({len(substack)/span_days:.1f}/day)")
    print("=" * 70)
    sub_senders = Counter()
    for e in substack:
        base = e["addr"].split("@")[0].split("+")[0]
        sub_senders[base] += 1
    for sender, cnt in sub_senders.most_common(40):
        per_day = cnt / span_days
        print(f"  {cnt:>5} ({per_day:>4.1f}/day)  {sender}")
    sub_total = sum(sub_senders.values())
    print(f"\n  Substack total: {sub_total} emails = {sub_total/len(all_emails)*100:.1f}% of all mail")

    # ─── ACTIONABLE: WHAT TO UNSUBSCRIBE ───
    print("\n" + "=" * 70)
    print("  🔧 ACTIONABLE: TOP BLOAT SOURCES TO UNSUBSCRIBE/DISABLE")
    print("=" * 70)

    # Group by source type
    categories = {
        "Substack Newsletters": [],
        "Marketing/Promo": [],
        "Service Notifications": [],
        "Financial Alerts": [],
        "Dev Tool Notifications": [],
    }

    marketing_domains = {"beehiiv.com", "workweek.com", "forwardfuture.ai",
        "infinitusmail.com", "leadremail.com", "joinsuperhuman.ai",
        "gethalfbaked.com", "simplilearnmailer.com", "tldrnewsletter.com",
        "founderos.com", "mikeschmitz.com", "imiloainstitute.com",
        "phoenixgrowthconnect.help", "startuptycoon.com", "joinsecret.com",
        "liamottley.com", "deeplearning.ai", "theaiautomators.com"}

    service_domains = {"zoho", "linear.app", "circle.so", "circle.co",
        "slack.com", "github.com", "notification"}

    financial_domains = {"creditsesame.com", "myfico.com", "monarch.com",
        "capitalone.com", "navyfederal.org", "kikoff.com", "servicecu.org",
        "bankofamerica.com", "verizonfinancialservices.com",
        "plaid.com", "mercury.com"}

    dev_domains = {"microsoft.com", "google.com", "railway.app",
        "ollama.com", "retellai.com", "llamaindex.ai", "relevanceai.com",
        "heygen.com", "deepgram.com", "supabase.com", "tavily.com",
        "datadoghq.com", "databox.com", "verkada.com", "mindsdb.com",
        "wordware.ai", "toolhouse.ai"}

    for addr, cnt in sender_counter.most_common(200):
        if cnt < 2:
            continue
        domain = extract_domain(addr)
        per_day = cnt / span_days

        if "substack" in domain:
            base = addr.split("@")[0].split("+")[0]
            categories["Substack Newsletters"].append((base, cnt, per_day))
        elif any(d in domain for d in marketing_domains):
            categories["Marketing/Promo"].append((addr, cnt, per_day))
        elif any(d in domain for d in service_domains):
            categories["Service Notifications"].append((addr, cnt, per_day))
        elif any(d in domain for d in financial_domains):
            categories["Financial Alerts"].append((addr, cnt, per_day))
        elif any(d in domain for d in dev_domains):
            categories["Dev Tool Notifications"].append((addr, cnt, per_day))

    for cat_name, items in categories.items():
        if not items:
            continue
        # Dedupe substack
        if cat_name == "Substack Newsletters":
            seen = {}
            for name, cnt, pd in items:
                if name not in seen:
                    seen[name] = (cnt, pd)
                else:
                    seen[name] = (seen[name][0] + cnt, seen[name][1] + pd)
            items = [(n, c, p) for n, (c, p) in seen.items()]
            items.sort(key=lambda x: -x[1])

        cat_total = sum(x[1] for x in items)
        print(f"\n  {cat_name} — {cat_total} emails total:")
        for name, cnt, pd in sorted(items, key=lambda x: -x[1])[:20]:
            action = "UNSUBSCRIBE" if pd > 0.3 else "consider"
            print(f"    {cnt:>4} ({pd:.1f}/day) {name:<50} → {action}")

    # ─── SUMMARY ───
    print("\n" + "=" * 70)
    print("  📋 EXECUTIVE SUMMARY")
    print("=" * 70)
    print(f"  Total emails:        {len(all_emails)}")
    print(f"  Date span:           {span_days} days")
    print(f"  Avg emails/day:      {len(dated)/span_days:.1f}")

    # Top 5 domains by daily volume
    top5 = sorted(domain_daily.items(), key=lambda x: -x[1])[:5]
    print(f"\n  Top 5 daily bloat sources:")
    for dom, pd in top5:
        total = domain_counter[dom]
        pct = total / len(all_emails) * 100
        print(f"    {pd:.1f}/day  {dom} ({total} total, {pct:.0f}%)")

    sub_pct = len(substack) / len(all_emails) * 100
    print(f"\n  Substack alone: {len(substack)} emails = {sub_pct:.0f}% of all mail ({len(substack)/span_days:.1f}/day)")

    # Save report
    report = {
        "total_emails": len(all_emails),
        "date_range": f"{min_date} to {max_date}",
        "span_days": span_days,
        "avg_per_day": round(len(dated) / span_days, 1),
        "top_domains": domain_counter.most_common(50),
        "top_senders": sender_counter.most_common(100),
        "substack_senders": sub_senders.most_common(50),
    }
    with open("bloat_audit.json", "w") as f:
        json.dump(report, f, indent=2, default=str)
    print(f"\n  Full data saved to bloat_audit.json")


if __name__ == "__main__":
    main()
