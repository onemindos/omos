---
id: root-00-20260114234500
title: "EMAIL CATEGORY ORGANIZATION STRATEGY"
type: guide
status: active
created: 2026-01-14
modified: 2026-01-14
owner: zeus
tags: [type/guide, status/active, email, organization]
ai_summary: "Complete strategy for organizing emails within pillar folders using category tags and smart filters"
---
# Email Category Organization Strategy

**Created:** 2026-01-14 23:45
**Status:** Ready for Implementation
**Owner:** Zeus Delacruz + Claude Code

---

## The Problem We're Solving

**Current State:**
- 15 old subfolders under "Zeus Delacruz" with emails mixed chronologically
- 5 new pillar folders (HP, LE, UI, GE, Inbox) with no internal organization
- New emails arrive but have no automatic categorization within folders
- No way to find "all finance emails" across both HP (personal) and GE (business)

**Target State:**
- Old emails migrated from 15 subfolders → 5 pillar folders with category tags
- New emails automatically tagged by sender/keyword patterns
- Cross-folder search enabled (e.g., #finance shows personal + business finance)
- Zero manual tagging for 90% of emails

---

## Part 1: Category Tag System Design

### The Two-Layer Tag Architecture

**Layer 1: Pillar Tags (Already Working)**
```
#hp   → Which domain? (Holistic Performance)
#le   → Which domain? (Legacy Evolution)
#ui   → Which domain? (Sovereign Operations)
#ge   → Which domain? (Generational Empire)
```

**Layer 2: Category Tags (NEW)**
```
#finance     → What type? (Banks, payments, billing)
#shopping    → What type? (E-commerce, purchases)
#family      → What type? (School, kids, relatives)
#travel      → What type? (Airlines, hotels, bookings)
#legal       → What type? (Attorneys, contracts)
... and more
```

**How They Work Together:**
```
Email from Chase Bank to hp@zeusdelacruz.me:
- Filter #1 → HP Folder + #hp (which domain)
- Filter #6 → + #finance (what type)
- Result: HP Folder with tags: #hp + #finance

Email from Business credit card to ge@zeusdelacruz.me:
- Filter #4 → GE Folder + #ge (which domain)
- Filter #6 → + #finance (what type)
- Result: GE Folder with tags: #ge + #finance

Click #finance → See BOTH personal and business finance across folders!
```

`★ Insight ─────────────────────────────────────`
**Why two-layer tagging is powerful**: Folders organize by life domain (personal, family, infrastructure, business). Tags organize by activity type (finance, shopping, travel). This matrix lets you answer BOTH questions: "What's happening in my business?" (click GE folder) AND "Where is all my money going?" (click #finance tag across all folders). One organizing system serves two completely different needs.
`─────────────────────────────────────────────────`

### Complete Category Tag List

**HP (Holistic Performance) Categories:**
```
#finance      - Banks, credit cards, investments, PayPal, Stripe
#shopping     - Amazon, eBay, retail stores, online purchases
#health       - Gym, doctors, pharmacy, fitness apps, wellness
#learning     - Courses, books, educational subscriptions, Udemy
#productivity - Software, apps, tools, subscriptions (Notion, Todoist)
#recreation   - Entertainment, hobbies, games, streaming services
```

**LE (Legacy Evolution) Categories:**
```
#family   - Schools, kids activities, relatives, family events
#travel   - Airlines, hotels, car rentals, trip bookings, Airbnb
#legal    - Attorneys, contracts, estate planning, notary
#home     - Utilities, repairs, home services, insurance
```

**UI (Sovereign Operations) Categories:**
```
#infrastructure - Hosting, servers, domains (Cloudflare, AWS, Railway)
#ai-services    - OpenAI, Anthropic, Letta, AI platform subscriptions
#dev-tools      - GitHub, IDEs, developer tools, code services
#monitoring     - Uptime alerts, system notifications, error logs
#automation     - N8N, Zapier, workflow tools
```

**GE (Generational Empire) Categories:**
```
#onemind       - OneMind Codex business emails
#legacyhq      - Legacy HQ emails (being phased out)
#restore       - Legacy Restore client emails
#affiliates    - Affiliate program emails and earnings
#partners      - Partner inquiries and communications
#revenue       - Payment notifications, earnings reports, invoices
#marketing     - Campaign updates, analytics, advertising
```

**Universal Categories (Any Folder):**
```
#automated     - Automated system emails, no-reply senders
#receipts      - Transaction receipts, confirmations
#newsletters   - Marketing newsletters, updates
#alerts        - Important alerts, notifications requiring attention
```

---

## Part 2: Old Email Migration Strategy

### Migration Mapping: Old Folders → New Folders + Tags

**Migrate to HP (Holistic Performance):**
```
Old Folder              New Folder    Category Tag     Notes
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Finance              →  HP         +  #finance        Personal banking
Shopping             →  HP         +  #shopping       Personal purchases
Performance          →  HP         +  #productivity   Self-optimization
Education            →  HP         +  #learning       Personal education
Recreation           →  HP         +  #recreation     Entertainment, hobbies
Productivity         →  HP         +  #productivity   Tools, apps
```

**Migrate to LE (Legacy Evolution):**
```
Old Folder              New Folder    Category Tag     Notes
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Family               →  LE         +  #family         School, kids, relatives
Travel               →  LE         +  #travel         Trip bookings, airlines
Legal                →  LE         +  #legal          Attorneys, estate
Legacy               →  LE         +  #family         Family legacy items
```

**Migrate to UI (Sovereign Operations):**
```
Old Folder              New Folder    Category Tag     Notes
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Training             →  UI         +  #dev-tools      Tech training
Legacyhq - Old       →  UI         +  #legacyhq       Business tech archives
legacyhq - New       →  UI         +  #legacyhq       Business tech current
```

**Migrate to GE (Generational Empire):**
```
Old Folder              New Folder    Category Tag     Notes
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Wealth               →  GE         +  #revenue        Business wealth building
Social (if business) →  GE         +  #marketing      Business social media
```

**Review Individually:**
```
Old Folder              Action                           Notes
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Save for Later       →  Review each email individually   Context-dependent
Social               →  Determine if personal or business First
```

### Migration Process (Step-by-Step)

**Phase 1: Bulk Move (Zoho Mail UI)**

```
Week 1: HP Migration (6 folders)
1. Open "Zeus Delacruz/Finance" folder
2. Select All emails (Ctrl+A / Cmd+A)
3. Click "Move" → Select "25-49 HP (Holistic Performance)"
4. Repeat for: Shopping, Performance, Education, Recreation, Productivity

Result: ~300+ emails moved to HP folder (but no tags yet)
```

**Phase 2: Bulk Tagging (Zoho Mail UI)**

```
Week 1 Continued: Add Category Tags
1. Open HP folder
2. Search: "from:@chase.com OR from:@bankofamerica.com"
3. Select all results
4. Click "Tag" → Create tag "#finance" → Apply
5. Repeat for other patterns (see Smart Filter Patterns below)

Result: Emails now have both #hp + #finance tags
```

**Phase 3: Repeat for Other Domains**

```
Week 2: LE Migration
- Move Family, Travel, Legal, Legacy folders → LE
- Bulk tag: #family, #travel, #legal

Week 3: UI Migration
- Move Training, Legacyhq folders → UI
- Bulk tag: #dev-tools, #legacyhq

Week 4: GE Migration
- Move Wealth → GE
- Review Social folder individually
- Bulk tag: #revenue, #marketing
```

---

## Part 3: Smart Filter Patterns (The Key to Automation)

### Philosophy: Patterns, Not Individual Senders

**❌ DON'T DO THIS (one filter per sender):**
```
Filter #6: Chase Bank
  FROM contains "@chase.com"
  Action: Apply tag #finance

Filter #7: Bank of America
  FROM contains "@bankofamerica.com"
  Action: Apply tag #finance

Filter #8: Wells Fargo
  FROM contains "@wellsfargo.com"
  Action: Apply tag #finance

... 50 more filters for every bank ...
```

**✅ DO THIS (pattern matching):**
```
Filter #6: Finance Pattern
  FROM contains "@chase.com"
  OR FROM contains "@bankofamerica.com"
  OR FROM contains "@wellsfargo.com"
  OR FROM contains "@paypal.com"
  OR FROM contains "@stripe.com"
  OR SUBJECT contains "statement"
  OR SUBJECT contains "invoice"
  OR SUBJECT contains "payment received"
  Action: Apply tag #finance

Result: ONE filter catches all banks + payment processors + keyword matches
```

### Recommended Filter Set (Filters #6-20)

**Filter #6: Finance Auto-Tag**
```
Name: Auto-Tag Finance
Priority: 6
Criteria:
  FROM contains "@chase.com"
  OR FROM contains "@bankofamerica.com"
  OR FROM contains "@wellsfargo.com"
  OR FROM contains "@paypal.com"
  OR FROM contains "@stripe.com"
  OR FROM contains "@square.com"
  OR FROM contains "@venmo.com"
  OR FROM contains "billing@"
  OR FROM contains "payments@"
  OR SUBJECT contains "statement"
  OR SUBJECT contains "invoice"
  OR SUBJECT contains "payment received"
  OR SUBJECT contains "bill due"
Actions:
  Apply tag: #finance
  (Do NOT move - pillar filters already handled routing)
```

**Filter #7: Shopping Auto-Tag**
```
Name: Auto-Tag Shopping
Priority: 7
Criteria:
  FROM contains "@amazon.com"
  OR FROM contains "@ebay.com"
  OR FROM contains "@walmart.com"
  OR FROM contains "@target.com"
  OR FROM contains "order-update@"
  OR FROM contains "shipping@"
  OR SUBJECT contains "order confirmation"
  OR SUBJECT contains "shipped"
  OR SUBJECT contains "delivered"
  OR SUBJECT contains "receipt"
  OR SUBJECT contains "purchase"
Actions:
  Apply tag: #shopping
```

**Filter #8: Family Auto-Tag**
```
Name: Auto-Tag Family
Priority: 8
Criteria:
  FROM contains ".edu"
  OR FROM contains "school@"
  OR FROM contains "teacher@"
  OR FROM contains "principal@"
  OR SUBJECT contains "parent"
  OR SUBJECT contains "student"
  OR SUBJECT contains "grades"
  OR SUBJECT contains "school event"
Actions:
  Apply tag: #family
```

**Filter #9: Travel Auto-Tag**
```
Name: Auto-Tag Travel
Priority: 9
Criteria:
  FROM contains "@delta.com"
  OR FROM contains "@united.com"
  OR FROM contains "@southwest.com"
  OR FROM contains "@marriott.com"
  OR FROM contains "@hilton.com"
  OR FROM contains "@airbnb.com"
  OR FROM contains "@expedia.com"
  OR SUBJECT contains "flight"
  OR SUBJECT contains "booking confirmed"
  OR SUBJECT contains "reservation"
  OR SUBJECT contains "itinerary"
Actions:
  Apply tag: #travel
```

**Filter #10: AI Services Auto-Tag**
```
Name: Auto-Tag AI Services
Priority: 10
Criteria:
  FROM contains "@openai.com"
  OR FROM contains "@anthropic.com"
  OR FROM contains "@letta.ai"
  OR FROM contains "usage@"
  OR FROM contains "api@"
  OR SUBJECT contains "API usage"
  OR SUBJECT contains "credits"
  OR SUBJECT contains "model"
Actions:
  Apply tag: #ai-services
```

**Filter #11: Infrastructure Auto-Tag**
```
Name: Auto-Tag Infrastructure
Priority: 11
Criteria:
  FROM contains "@cloudflare.com"
  OR FROM contains "@aws.amazon.com"
  OR FROM contains "@digitalocean.com"
  OR FROM contains "@railway.app"
  OR FROM contains "@heroku.com"
  OR FROM contains "@vercel.com"
  OR FROM contains "@netlify.com"
  OR SUBJECT contains "server"
  OR SUBJECT contains "downtime"
  OR SUBJECT contains "SSL certificate"
  OR SUBJECT contains "domain"
Actions:
  Apply tag: #infrastructure
```

**Filter #12: Health Auto-Tag**
```
Name: Auto-Tag Health
Priority: 12
Criteria:
  FROM contains "health@"
  OR FROM contains "doctor@"
  OR FROM contains "pharmacy@"
  OR FROM contains "@gym"
  OR FROM contains "@fitbit.com"
  OR FROM contains "@myfitnesspal.com"
  OR SUBJECT contains "appointment"
  OR SUBJECT contains "prescription"
  OR SUBJECT contains "medical"
Actions:
  Apply tag: #health
```

**Filter #13: Learning Auto-Tag**
```
Name: Auto-Tag Learning
Priority: 13
Criteria:
  FROM contains "@udemy.com"
  OR FROM contains "@coursera.org"
  OR FROM contains "@edx.org"
  OR FROM contains "@skillshare.com"
  OR FROM contains "@linkedin.com/learning"
  OR SUBJECT contains "course"
  OR SUBJECT contains "lesson"
  OR SUBJECT contains "certificate"
Actions:
  Apply tag: #learning
```

**Filter #14: Legal Auto-Tag**
```
Name: Auto-Tag Legal
Priority: 14
Criteria:
  FROM contains "attorney@"
  OR FROM contains "lawyer@"
  OR FROM contains "legal@"
  OR FROM contains "@law"
  OR SUBJECT contains "contract"
  OR SUBJECT contains "agreement"
  OR SUBJECT contains "legal notice"
Actions:
  Apply tag: #legal
```

**Filter #15: Dev Tools Auto-Tag**
```
Name: Auto-Tag Dev Tools
Priority: 15
Criteria:
  FROM contains "@github.com"
  OR FROM contains "@gitlab.com"
  OR FROM contains "@bitbucket.org"
  OR FROM contains "code@"
  OR FROM contains "developers@"
  OR SUBJECT contains "pull request"
  OR SUBJECT contains "commit"
  OR SUBJECT contains "deployment"
  OR SUBJECT contains "build"
Actions:
  Apply tag: #dev-tools
```

### Progressive Enhancement Strategy

**Month 1: Core Filters (Start with 5)**
```
Create only the most important:
1. #finance (catches 80% of important emails)
2. #shopping (high volume, easy to identify)
3. #family (school emails critical)
4. #travel (bookings need quick access)
5. #infrastructure (system alerts important)
```

**Month 2: Add as You Go**
```
Notice AI service emails not tagged?
→ Edit Filter #10 (AI Services), add new pattern

Notice new bank sending emails?
→ Edit Filter #6 (Finance), add "@newbank.com"

Notice recurring pattern not caught?
→ Create new filter for that pattern
```

**Month 3: Refinement**
```
Check false positives:
- Emails tagged #finance that aren't finance?
  → Refine filter criteria

- Emails NOT tagged that should be?
  → Add missing patterns

- Some filters never used?
  → Delete or consolidate
```

---

## Part 4: Handling New Emails (Automation Workflow)

### The Question: "Do I need a filter for every new sender?"

**Answer: NO - Use pattern expansion**

### Scenario 1: New Bank Account

**You open Wells Fargo account, emails arrive untagged**

```
Option A: Edit Existing Filter (Recommended)
1. Go to Zoho Mail → Settings → Filters
2. Find Filter #6 (Finance Pattern)
3. Click Edit
4. Add line: OR FROM contains "@wellsfargo.com"
5. Save

Result: All future Wells Fargo emails auto-tagged #finance
```

**Option B: Create Exception Filter (If needed)**
```
Only if Wells Fargo needs DIFFERENT handling:

Filter #16: Wells Fargo Investment
  FROM contains "@wellsfargo.com" AND SUBJECT contains "investment"
  Actions: Apply tags #finance + #investments

Filter #6 will still apply #finance
Filter #16 will ADD #investments for specific emails
```

### Scenario 2: New E-commerce Store

**You start shopping at Etsy, receipts arrive untagged**

```
1. First email arrives → Goes to HP folder (via hp@zeusdelacruz.me)
2. No category tag applied (not in filters yet)
3. Edit Filter #7 (Shopping Pattern)
4. Add: OR FROM contains "@etsy.com"
5. Save

Result: Next Etsy email auto-tagged #shopping
```

### Scenario 3: Completely New Email Type

**You sign up for meal delivery service (HelloFresh)**

```
Question: Is this shopping or health?

Decision Tree:
- Is it a purchase receipt? → #shopping
- Is it nutrition/wellness related? → #health
- Is it recurring subscription? → #shopping + #subscriptions (create new tag)

Add to Filter #7 (Shopping):
OR FROM contains "@hellofresh.com"
```

### Keyword-Based Automation (No Sender Required)

**Filters can match SUBJECTS without knowing sender:**

```
Filter #17: Receipts Universal
  SUBJECT contains "receipt for"
  OR SUBJECT contains "thanks for your purchase"
  OR SUBJECT contains "order #"
  Actions: Apply tag #receipts

Result: ANY email with receipt keywords gets tagged
```

**Filter #18: Alerts Universal**
```
  SUBJECT contains "alert:"
  OR SUBJECT contains "[URGENT]"
  OR SUBJECT contains "action required"
  OR SUBJECT contains "expires"
  Actions: Apply tag #alerts

Result: Urgent emails across all senders get flagged
```

---

## Part 5: Implementation Checklist

### Week 1: Foundation Setup

**Day 1-2: Create Core Category Filters**
```
☐ Create Filter #6: Finance Pattern (15 patterns)
☐ Create Filter #7: Shopping Pattern (10 patterns)
☐ Create Filter #8: Family Pattern (5 patterns)
☐ Create Filter #9: Travel Pattern (10 patterns)
☐ Create Filter #10: AI Services Pattern (5 patterns)
☐ Test each filter with existing emails
```

**Day 3-4: Migrate Old Emails (HP Domain)**
```
☐ Move "Zeus Delacruz/Finance" → HP Folder
☐ Move "Zeus Delacruz/Shopping" → HP Folder
☐ Move "Zeus Delacruz/Performance" → HP Folder
☐ Move "Zeus Delacruz/Education" → HP Folder
☐ Move "Zeus Delacruz/Recreation" → HP Folder
☐ Move "Zeus Delacruz/Productivity" → HP Folder
☐ Bulk tag finance emails: Search + Select All + Apply #finance
☐ Bulk tag shopping emails: Search + Select All + Apply #shopping
☐ Bulk tag learning emails: Search + Select All + Apply #learning
```

**Day 5-7: Test & Refine**
```
☐ Send test emails to hp@zeusdelacruz.me from Chase, Amazon, etc.
☐ Verify emails land in HP Folder with BOTH #hp + #category tags
☐ Check false positives (emails tagged incorrectly)
☐ Refine filter criteria as needed
```

### Week 2: LE Domain Migration

```
☐ Move "Zeus Delacruz/Family" → LE Folder
☐ Move "Zeus Delacruz/Travel" → LE Folder
☐ Move "Zeus Delacruz/Legal" → LE Folder
☐ Move "Zeus Delacruz/Legacy" → LE Folder
☐ Bulk tag family emails
☐ Bulk tag travel emails
☐ Bulk tag legal emails
☐ Test routing with emails to le@zeusdelacruz.me
```

### Week 3: UI & GE Domain Migration

```
☐ Move "Zeus Delacruz/Training" → UI Folder
☐ Move "Legacyhq - Old" → UI Folder
☐ Move "legacyhq - New" → UI Folder
☐ Move "Zeus Delacruz/Wealth" → GE Folder
☐ Create Filter #11: Infrastructure Pattern
☐ Create Filter #15: Dev Tools Pattern
☐ Bulk tag infrastructure emails
☐ Bulk tag dev tool emails
☐ Test routing
```

### Week 4: Cleanup & Optimization

```
☐ Review "Zeus Delacruz/Save for Later" - file individually
☐ Review "Zeus Delacruz/Social" - split personal vs business
☐ Delete empty old folders
☐ Rename "Zeus Delacruz" → "Archive - Pre-2026 (Empty)"
☐ Document which patterns work best
☐ Create quick reference card
```

---

## Part 6: Daily Workflow with Categories

### Morning Email Routine (Enhanced)

```
1. Open zeus@zeusdelacruz.com
2. Process Inbox (business emails requiring attention)
3. Check #alerts tag (urgent items across all folders)
4. Achieve Inbox Zero
5. Check pillar folders ONLY if needed:
   - Click HP → View personal life (already tagged)
   - Click #finance → View ALL finance (personal + business)
   - Click #family → View ALL family items
```

### Finding Specific Emails

**Scenario: Need all finance emails from last quarter**
```
Option 1: Tag-Based (Cross-Folder)
Click #finance → Shows ALL finance (HP personal + GE business)
Filter by date → Q4 2025
Result: Every finance email across entire system

Option 2: Folder-Based (Domain-Specific)
Click HP Folder → Search "tag:#finance"
Result: Only PERSONAL finance emails
```

**Scenario: Planning family vacation**
```
Click #travel tag → See all:
- Past trip bookings (receipts, confirmations)
- Airline frequent flyer programs
- Hotel loyalty programs
- Upcoming reservations
All in one view regardless of which folder they landed in
```

**Scenario: Tax time - need all receipts**
```
Click #receipts tag → Every transaction receipt
Filter by date → 2025
Export to PDF for accountant
```

### When to Create New Tags

**Create new tag when:**
- You search for same type repeatedly
- Existing tags don't capture specific need
- Category has 10+ emails

**Examples:**
```
#subscriptions - Recurring monthly charges
#investments - Investment accounts separate from banking
#insurance - Insurance policies and claims
#warranties - Product warranties and support
#events - Tickets, registrations, event confirmations
```

---

## Part 7: Advanced Tips & Tricks

### Combine Tags for Power Searches

```
#hp AND #finance → Personal finance only
#ge AND #finance → Business finance only
#finance AND #receipts → Finance receipts only
#alerts AND #infrastructure → Urgent server issues
```

### Smart Filter Ordering (Critical)

```
ORDER MATTERS:
1-5: Domain routing filters (must execute FIRST to move to folders)
6-20: Category tagging filters (execute AFTER to add tags)

Why: Category filters can't find emails if they've been moved already
Solution: Pillar filters move, category filters tag (no moving)
```

### Exception Handling

```
Problem: Bank emails tagged #finance but investment emails need #investments too

Solution: Create filter AFTER #finance filter:
Filter #21: Investments Sub-Tag
  FROM contains "@chase.com" AND SUBJECT contains "investment"
  Actions: Apply tag #investments (in addition to #finance)

Result: Investment emails get BOTH tags: #finance + #investments
```

### Temporary Tags During Migration

```
#migration - Emails from old aliases, needs account update
#review - Emails needing manual review
#archive-2025 - Old emails from specific year
#cleanup - Items pending cleanup decision

Delete these tags after migration complete
```

---

## Part 8: Success Metrics

### How to Know It's Working

**Week 1:**
- ✅ 5 core filters created and active
- ✅ 300+ old emails migrated to HP folder
- ✅ New emails auto-tagged with 80% accuracy
- ✅ Can find finance emails in < 5 seconds

**Week 4:**
- ✅ All 15 old folders migrated
- ✅ 10-15 category filters operational
- ✅ 90% of emails auto-tagged correctly
- ✅ Morning inbox triage: 10 minutes → 3 minutes

**Month 3:**
- ✅ Zero manual tagging required
- ✅ Can find any email type in seconds
- ✅ Cross-folder searches reveal patterns
- ✅ Filter patterns optimized and stable

---

## Part 9: Common Questions

**Q: Do I need to manually tag old emails?**
A: Only once during migration. Use bulk select → bulk tag. Future emails auto-tag.

**Q: What if an email matches multiple categories?**
A: That's perfect! One email can have multiple tags (#finance + #receipts + #automated).

**Q: How many filters is too many?**
A: 15-20 filters is ideal. More than 30 gets hard to manage. Use pattern matching to consolidate.

**Q: Can I change tags later?**
A: Yes! Tags are just labels. Add/remove anytime without affecting emails.

**Q: What if a filter tags incorrectly?**
A: Edit the filter criteria. Previous emails keep wrong tags, but you can bulk-fix them.

**Q: Should I tag newsletters separately?**
A: Optional. Create #newsletters if you want to bulk delete/archive them easily.

---

## Summary: Your Action Plan

**This Week:**
1. Create 5 core category filters (#finance, #shopping, #family, #travel, #ai-services)
2. Migrate HP domain old folders (Finance, Shopping, Performance, Education, Recreation)
3. Bulk tag migrated emails using search patterns
4. Test with new emails coming in

**Next 3 Weeks:**
1. Migrate LE, UI, GE domains
2. Add more category filters as patterns emerge
3. Refine filter criteria based on false positives
4. Delete old empty folders

**Month 2+:**
1. Monitor auto-tagging accuracy
2. Create additional tags as needed
3. Optimize filter patterns
4. Document what works best

**Result:** Email system that automatically organizes 90% of emails with zero ongoing effort.

---

*Document created: 2026-01-14 23:45*
*Status: Ready for Week 1 implementation*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
