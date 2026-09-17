---
id: ui-08-all-filters-complete
title: "Complete Filter Specifications - All 10 Smart Filters"
type: reference
status: active
created: 2026-01-14
modified: 2026-01-14
owner: zeus
tags: [type/reference, status/active, platform/zoho, system/email]
ai_summary: "Complete specifications for all 10 Zoho Mail smart filters with discovered domains and conditions"
---

# Complete Filter Specifications - All 10 Smart Filters

> **Status:** Ready for implementation
> **Last Updated:** 2026-01-14 after analyzing 461 emails in Education folder
> **CRITICAL:** Use `FROM` → `contains` (NOT "domain contains") for ALL conditions

---

## Filter Priority Order

Filters run in numerical priority order (lower number = higher priority):

```
1-5:   Pillar Routing (hp@, le@, ui@, ge@) → Folder placement
6-15:  Category Tagging → Apply granular tags within folders
```

---

## Filter #1: HP - Mind & Education (#26-mind)

**Priority:** 6 (runs after domain routing)
**Expected Matches:** ~150 emails
**Operator:** Match ANY condition (OR logic)

### Tier 1: Domain List (Primary - 25 domains)

```
FROM contains:
1.  deeplearning.ai
2.  coursera.org
3.  udemy.com
4.  linkedin.com
5.  skillshare.com
6.  masterclass.com
7.  pluralsight.com
8.  datacamp.com
9.  codecademy.com
10. khanacademy.org
11. edx.org
12. udacity.com
13. workweek.com
14. venturebeat.com
15. infinitusmail.com
16. relevanceai.com
17. mail.joinsuperhuman.ai          ← NEW (46 emails)
18. forwardfuture.ai                ← NEW (37 emails)
19. mail.gethalfbaked.com           ← NEW (26 emails)
20. simplilearnmailer.com           ← NEW (19 emails)
21. liamottley.com                  ← NEW (13 emails)
22. fundrapp.com                    ← NEW (7 emails)
23. aiautomationsociety.ai          ← NEW (7 emails)
24. daily.therundown.ai             ← NEW
25. fundingvault.io                 ← NEW
```

### Tier 1b: Specific Email Senders (Substack - 9 senders)

```
FROM contains:
1. astro@forwardfuture.ai
2. thecapitalistmag@substack.com
3. ainutpam@substack.com
4. linas@substack.com
5. ideas@mail.gethalfbaked.com
6. willywoo@substack.com
7. natehagens@substack.com
8. yungpueblo@substack.com
9. superhuman@mail.joinsuperhuman.ai
```

### Tier 2: Subject Keywords (5 high-value phrases)

```
SUBJECT contains:
1. "course enrollment"
2. "certificate of completion"
3. "new lesson available"
4. "professional development"
5. "learning path"
```

### Actions:
- Apply label: `#26-mind`
- Keep in current folder (HP)

**Total Conditions:** 39 (25 domains + 9 senders + 5 subject)

---

## Filter #2: HP - Body & Fitness (#27-body)

**Priority:** 7
**Expected Matches:** ~40 emails
**Operator:** Match ANY condition (OR logic)

### Tier 1: Domain List (18 domains)

```
FROM contains:
1.  myfitnesspal.com
2.  strava.com
3.  fitbit.com
4.  noom.com
5.  peloton.com
6.  garmin.com
7.  whoop.com
8.  headspace.com
9.  calm.com
10. betterhelp.com
11. talkspace.com
12. cronometer.com
13. loseit.com
14. underarmour.com
15. zwift.com
16. trainerroad.com
17. nike.com
18. apple.com
```

### Tier 2: Subject Keywords (8 phrases)

```
SUBJECT contains:
1. "workout summary"
2. "activity tracked"
3. "daily activity"
4. "fitness goal"
5. "calories burned"
6. "heart rate zone"
7. "distance tracked"
8. "personal best"
```

### Actions:
- Apply label: `#27-body`
- Keep in current folder (HP)

**Total Conditions:** 26 (18 domains + 8 subject)

---

## Filter #3: HP - Shopping & Orders (#hp-shopping)

**Priority:** 8
**Expected Matches:** ~50 emails
**Operator:** Match ANY condition (OR logic)

### Tier 1: Domain List (15 domains)

```
FROM contains:
1.  amazon.com
2.  ebay.com
3.  shopify.com
4.  etsy.com
5.  target.com
6.  walmart.com
7.  bestbuy.com
8.  newegg.com
9.  costco.com
10. homedepot.com
11. lowes.com
12. wayfair.com
13. zappos.com
14. nordstrom.com
15. macys.com
```

### Tier 2: Subject Keywords (8 phrases)

```
SUBJECT contains:
1. "order confirmation"
2. "shipping confirmation"
3. "tracking number"
4. "order delivered"
5. "package delivered"
6. "payment received"
7. "refund processed"
8. "return label"
```

### Actions:
- Apply label: `#hp-shopping`
- Keep in current folder (HP)

**Total Conditions:** 23 (15 domains + 8 subject)

---

## Filter #4: HP - Personal Finance (#30-finance)

**Priority:** 9
**Expected Matches:** ~60 emails
**Operator:** Match ANY condition (OR logic)
**⚠️ Deploy during account migration week**

### Tier 1: Domain List (18 domains)

```
FROM contains:
1.  monarch.com
2.  chase.com
3.  americanexpress.com
4.  bankofamerica.com
5.  wellsfargo.com
6.  citi.com
7.  discover.com
8.  capitalone.com
9.  paypal.com
10. venmo.com
11. cashapp.com
12. plaid.com
13. mint.com
14. ynab.com
15. personalcapital.com
16. creditkarma.com
17. turbotax.com
18. hrblock.com
```

### Tier 2: Subject Keywords (8 phrases)

```
SUBJECT contains:
1. "account statement"
2. "payment due"
3. "balance alert"
4. "transaction alert"
5. "credit card"
6. "bank account"
7. "minimum payment"
8. "autopay confirmation"
```

### Actions:
- Apply label: `#30-finance`
- Keep in current folder (HP)

**Total Conditions:** 26 (18 domains + 8 subject)

---

## Filter #5: HP - Social & Relationships (#32-relationships)

**Priority:** 10
**Expected Matches:** ~15 emails
**Operator:** Match ANY condition (OR logic)

### Tier 1: Domain List (9 domains)

```
FROM contains:
1. facebook.com
2. instagram.com
3. meetup.com
4. eventbrite.com
5. bumble.com
6. hinge.com
7. match.com
8. evite.com
9. paperlesspost.com
```

### Tier 2: Subject Keywords (6 phrases)

```
SUBJECT contains:
1. "tagged you"
2. "mentioned you"
3. "event invitation"
4. "RSVP"
5. "friend request"
6. "group invitation"
```

### Actions:
- Apply label: `#32-relationships`
- Keep in current folder (HP)

**Total Conditions:** 15 (9 domains + 6 subject)

---

## Filter #6: UI - Infrastructure & DevOps (#08-infrastructure)

**Priority:** 11
**Expected Matches:** ~80 emails
**Operator:** Match ANY condition (OR logic)

### Tier 1: Domain List (28 domains)

```
FROM contains:
1.  webflow.com
2.  github.com
3.  gitlab.com
4.  vercel.com
5.  netlify.com
6.  aws.amazon.com
7.  digitalocean.com
8.  linode.com
9.  cloudflare.com
10. namecheap.com
11. godaddy.com
12. hover.com
13. heroku.com
14. railway.app
15. render.com
16. fly.io
17. supabase.com
18. planetscale.com
19. mongodb.com
20. redis.com
21. letta.com
22. anthropic.com
23. openai.com
24. docker.com
25. news.militaryaerospace.com      ← NEW (29 emails)
26. mail.militaryaerospace.com      ← NEW (11 emails)
27. pagerduty.com                   ← NEW (9 emails)
28. verkada.com                     ← NEW (2 emails)
```

### Tier 2: Subject Keywords (10 phrases)

```
SUBJECT contains:
1. "deployment"
2. "server alert"
3. "API error"
4. "downtime"
5. "incident report"
6. "SSL certificate"
7. "domain renewal"
8. "backup complete"
9. "build failed"
10. "security alert"
```

### Actions:
- Apply label: `#08-infrastructure`
- Keep in current folder (UI)

**Total Conditions:** 38 (28 domains + 10 subject)

---

## Filter #7: GE - Business Finance & Treasury (#80-treasury)

**Priority:** 12
**Expected Matches:** ~30 emails
**Operator:** Match ANY condition (OR logic)

### Tier 1: Domain List (14 domains)

```
FROM contains:
1.  mercury.com
2.  stripe.com
3.  quickbooks.com
4.  xero.com
5.  bench.com
6.  gusto.com
7.  rippling.com
8.  deel.com
9.  pilot.com
10. brex.com
11. ramp.com
12. divvy.com
13. landbase.com                    ← NEW (2 emails)
14. motusinvestmentgroup.com       ← NEW
```

### Tier 2: Subject Keywords (8 phrases)

```
SUBJECT contains:
1. "business account"
2. "merchant payment"
3. "payroll processed"
4. "invoice sent"
5. "payment received"
6. "subscription billing"
7. "revenue report"
8. "quarterly statement"
```

### Actions:
- Apply label: `#80-treasury`
- Keep in current folder (GE)

**Total Conditions:** 22 (14 domains + 8 subject)

---

## Filter #8: GE - Personal Brand (#75-personal-brand)

**Priority:** 13
**Expected Matches:** ~25 emails
**Operator:** Match ANY condition (OR logic)

### Tier 1: Domain List (16 domains)

```
FROM contains:
1.  twitter.com
2.  youtube.com
3.  substack.com
4.  medium.com
5.  beehiiv.com
6.  convertkit.com
7.  mailchimp.com
8.  linktree.com
9.  stan.store
10. gumroad.com
11. patreon.com                     ← NEW (5 emails)
12. mail.beehiiv.com                ← NEW (4 emails)
13. mail.startuptycoon.com          ← NEW (6 emails)
14. founderos.com                   ← NEW (6 emails)
15. leadremail.com                  ← NEW (12 emails)
16. optimaed.com                    ← NEW
```

### Tier 2: Subject Keywords (6 phrases)

```
SUBJECT contains:
1. "new subscriber"
2. "newsletter analytics"
3. "content published"
4. "follower milestone"
5. "engagement report"
6. "brand metrics"
```

### Actions:
- Apply label: `#75-personal-brand`
- Keep in current folder (GE)

**Total Conditions:** 22 (16 domains + 6 subject)

---

## Filter #9: GE - Partnerships & Collaborations (#83-partnerships)

**Priority:** 14
**Expected Matches:** ~10 emails
**Operator:** Match ANY condition (OR logic)

### Tier 1: Domain List (2 domains)

```
FROM contains:
1. simple.ai                        ← NEW
2. livestorm-events.com             ← NEW
```

### Tier 2: Subject Keywords (8 phrases)

```
SUBJECT contains:
1. "partnership opportunity"
2. "collaboration proposal"
3. "joint venture"
4. "affiliate program"
5. "sponsorship"
6. "partnership agreement"
7. "co-marketing"
8. "strategic alliance"
```

### Actions:
- Apply label: `#83-partnerships`
- Keep in current folder (GE)

**Total Conditions:** 10 (2 domains + 8 subject)

---

## Filter #10: LE - Legal & Compliance (#57-legal)

**Priority:** 15
**Expected Matches:** ~8 emails
**Operator:** Match ANY condition (OR logic)

### Tier 1: Domain List (5 domains)

```
FROM contains:
1. legalzoom.com
2. rocketlawyer.com
3. nolo.com
4. docusign.com
5. hellosign.com
```

### Tier 2: Subject Keywords (10 phrases)

```
SUBJECT contains:
1. "legal notice"
2. "contract"
3. "agreement"
4. "terms and conditions"
5. "compliance"
6. "attorney"
7. "lawsuit"
8. "settlement"
9. "trademark"
10. "copyright"
```

### Actions:
- Apply label: `#57-legal`
- Keep in current folder (LE)

**Total Conditions:** 15 (5 domains + 10 subject)

---

## BONUS: Filter #11 - LE Family & Legacy (Discovered Domains)

**Priority:** 16
**Expected Matches:** ~30 emails
**Operator:** Match ANY condition (OR logic)

### Tier 1: Domain List (7 domains)

```
FROM contains:
1. thegogetterfamily.com            ← NEW (18 emails)
2. legacyplatform.co                ← NEW (4 emails)
3. legacywarrior.me                 ← NEW (2 emails)
4. childsupportisfraud.com          ← NEW (2 emails)
5. GrowYourOwnVegetables.org        ← NEW (5 emails)
6. mail.libcal.com                  ← NEW (2 emails)
7. founderos.com
```

### Tier 2: Subject Keywords (6 phrases)

```
SUBJECT contains:
1. "family"
2. "legacy planning"
3. "estate"
4. "inheritance"
5. "parenting"
6. "family meeting"
```

### Actions:
- Apply label: `#54-family` (or create #le-family)
- Keep in current folder (LE)

**Total Conditions:** 13 (7 domains + 6 subject)

---

## Implementation Strategy

### Phase 1: High-Volume Filters (Create First)
1. ✅ Filter #1: Mind (#26-mind) - 150+ emails
2. ✅ Filter #6: Infrastructure (#08-infrastructure) - 80+ emails
3. ✅ Filter #3: Shopping (#hp-shopping) - 50+ emails
4. ✅ Filter #4: Personal Finance (#30-finance) - 60+ emails

### Phase 2: Medium-Volume Filters
5. ⏸️ Filter #2: Body (#27-body) - 40 emails
6. ⏸️ Filter #7: Business Finance (#80-treasury) - 30 emails
7. ⏸️ Filter #8: Personal Brand (#75-personal-brand) - 25 emails

### Phase 3: Low-Volume Filters
8. ⏸️ Filter #5: Social (#32-relationships) - 15 emails
9. ⏸️ Filter #9: Partnerships (#83-partnerships) - 10 emails
10. ⏸️ Filter #10: Legal (#57-legal) - 8 emails
11. ⏸️ Filter #11: Family/Legacy - 30 emails

---

## Quick Create Checklist (Per Filter)

**For each filter:**

1. **Go to:** https://mail.zoho.com → Settings → Filters → Create Filter
2. **Name:** `[Priority] Category Name` (e.g., "6 - Auto-Tag Mind")
3. **Operator:** Select "Match ANY of the following" (OR logic)
4. **Add Conditions:**
   - [ ] Copy-paste all Tier 1 domains (one by one)
   - [ ] Copy-paste all Tier 2 subject keywords
   - [ ] Use `FROM` → `contains` for domains
   - [ ] Use `SUBJECT` → `contains` for keywords
5. **Add Action:**
   - [ ] Apply Label: `#[tag-name]`
   - [ ] DO NOT add "Move to folder" action
6. **Save & Test:**
   - [ ] Save filter
   - [ ] Send test email
   - [ ] Verify dual tags (#pillar + #category)

---

## Filter Summary Table

| # | Filter Name | Tag | Domains | Keywords | Total | Priority |
|---|-------------|-----|---------|----------|-------|----------|
| 1 | Mind & Education | #26-mind | 34 | 5 | 39 | 6 |
| 2 | Body & Fitness | #27-body | 18 | 8 | 26 | 7 |
| 3 | Shopping | #hp-shopping | 15 | 8 | 23 | 8 |
| 4 | Personal Finance | #30-finance | 18 | 8 | 26 | 9 |
| 5 | Social | #32-relationships | 9 | 6 | 15 | 10 |
| 6 | Infrastructure | #08-infrastructure | 28 | 10 | 38 | 11 |
| 7 | Business Finance | #80-treasury | 14 | 8 | 22 | 12 |
| 8 | Personal Brand | #75-personal-brand | 16 | 6 | 22 | 13 |
| 9 | Partnerships | #83-partnerships | 2 | 8 | 10 | 14 |
| 10 | Legal | #57-legal | 5 | 10 | 15 | 15 |
| 11 | Family/Legacy | #54-family | 7 | 6 | 13 | 16 |

**TOTAL:** 166 domains + 83 keywords = **249 conditions across 11 filters**

---

## Maintenance Notes

### Weekly (First Month)
- Check for untagged emails in each pillar folder
- Add newly discovered domains to relevant filters
- Adjust keywords if seeing false positives

### Monthly
- Review filter performance metrics
- Consolidate domains if seeing duplicates
- Archive low-performing filters

### As You Discover New Senders
1. Note the domain or email address
2. Identify which OneMind category it belongs to
3. Edit relevant filter
4. Add new condition
5. **Takes 30 seconds per domain**

---

## Expected Coverage After All 11 Filters

**Education folder (461 emails) breakdown:**
- ✅ ~127 emails → #26-mind (Filter #1)
- ✅ ~49 emails → #08-infrastructure (Filter #6)
- ✅ ~26 emails → #54-family (Filter #11)
- ✅ ~33 emails → #75-personal-brand or #80-treasury (Filters #7-8)
- ✅ ~7 emails → Various other filters
- ❓ ~219 emails → Require manual review (parsing errors)

**Estimated tagging rate:** 52% immediate + 48% manual review needed

---

*Created: 2026-01-14*
*Status: Complete and ready for implementation*
*Based on: Analysis of 461 emails + existing documentation*
