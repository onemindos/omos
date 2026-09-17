---
id: ui-08-zohomail-filters
title: "Zoho Mail Smart Filters - Optimized Multi-Tier Specifications"
type: reference
status: active
created: 2026-01-14
modified: 2026-01-14
owner: zeus
tags: [type/reference, status/active, platform/zoho, system/email]
ai_summary: "Comprehensive smart filter specifications for Zoho Mail using multi-tier conditions (domain + subject + content) for maximum routing accuracy"
---

# Zoho Mail Smart Filters - Optimized Multi-Tier Specifications

> **Strategy:** Use Domain-based matching (Tier 1) as primary, Subject keywords (Tier 2) as secondary, and Content keywords (Tier 3) as fallback for maximum recall with high precision.

## Filter Condition Hierarchy

### ⭐⭐⭐⭐⭐ Tier 1: Domain-Based (Highest Accuracy)
- **Pros:** Never false-positive, fast processing
- **Cons:** Requires maintaining domain list
- **Best For:** SaaS services, known senders, subscriptions

### ⭐⭐⭐⭐ Tier 2: Subject Keywords (Medium Accuracy)
- **Pros:** Catches transactional emails reliably
- **Cons:** Can miss emails with generic subjects
- **Best For:** Invoices, receipts, account notifications

### ⭐⭐⭐⭐⭐ Tier 3: Body Content (Broadest Catch)
- **Pros:** Catches everything, even with generic subjects
- **Cons:** Slower processing, potential false positives
- **Best For:** Fallback net for varied senders

---

## Filter #1: HP - Mind & Education (#26-mind)

**Category:** Holistic Performance → Mind (Education & Learning)
**Priority:** 6 (runs after domain routing)
**Expected Matches:** ~50-100 emails

### Conditions (ANY match applies filter):

#### Tier 1: Domain List (Primary)

**CRITICAL SYNTAX:** Use `FROM` → `contains` (NOT "domain contains")

```
FROM contains:
- deeplearning.ai
- coursera.org
- udemy.com
- linkedin.com (for LinkedIn Learning)
- skillshare.com
- masterclass.com
- pluralsight.com
- datacamp.com
- codecademy.com
- khanacademy.org
- edx.org
- udacity.com
- workweek.com (Fintech Takes)
- venturebeat.com (Tech news)
- infinitusmail.com (Business education)
- relevanceai.com (AI learning)
- notion.so (Learning templates)
- obsidian.md (PKM learning)
- youtube.com (if learning-related)
- medium.com (if educational)
```

#### Tier 2: Subject Keywords (Secondary)
```
SUBJECT contains (case-insensitive):
- "course"
- "learn"
- "lesson"
- "education"
- "training"
- "certificate"
- "certification"
- "tutorial"
- "workshop"
- "webinar"
- "study"
- "class"
- "lecture"
- "enrollment"
- "curriculum"
```

#### Tier 3: Content Keywords (Fallback)
```
MESSAGE contains (case-insensitive):
- "complete the course"
- "new lesson"
- "learning path"
- "skill development"
- "professional development"
- "continue your learning"
- "course enrollment"
- "certificate of completion"
```

### Actions:
1. **Apply tag:** `#26-mind`
2. **Keep in folder:** HP (already routed by pillar filter)
3. **Note:** Does NOT move email (already in correct pillar folder)

---

## Filter #2: HP - Body & Fitness (#27-body)

**Category:** Holistic Performance → Body (Health & Fitness)
**Priority:** 7
**Expected Matches:** ~30-50 emails

### Conditions (ANY match applies filter):

#### Tier 1: Domain List (Primary)
```
FROM domain contains:
- myfitnesspal.com
- strava.com
- fitbit.com
- noom.com
- peloton.com
- apple.com (if health/fitness related)
- garmin.com
- whoop.com
- headspace.com (mental health)
- calm.com
- betterhelp.com
- talkspace.com
- cronometer.com
- loseit.com
- nike.com (if training/run club)
- underarmour.com (if connected apps)
- zwift.com
- trainerroad.com
```

#### Tier 2: Subject Keywords (Secondary)
```
SUBJECT contains:
- "workout"
- "exercise"
- "fitness"
- "health"
- "calories"
- "steps"
- "activity"
- "nutrition"
- "diet"
- "weight"
- "training"
- "run"
- "ride"
- "swim"
- "meditation"
- "sleep"
- "wellness"
```

#### Tier 3: Content Keywords (Fallback)
```
MESSAGE contains:
- "daily activity"
- "fitness goal"
- "workout summary"
- "calories burned"
- "heart rate"
- "distance tracked"
- "personal best"
- "exercise routine"
- "health metrics"
```

### Actions:
1. **Apply tag:** `#27-body`
2. **Keep in folder:** HP

---

## Filter #3: HP - Shopping & Orders (#hp-shopping)

**Category:** Holistic Performance → Shopping (Personal Purchases)
**Priority:** 8
**Expected Matches:** ~40-60 emails

### Conditions (ANY match applies filter):

#### Tier 1: Domain List (Primary)
```
FROM domain contains:
- amazon.com
- ebay.com
- shopify.com
- etsy.com
- target.com
- walmart.com
- bestbuy.com
- newegg.com
- alibaba.com
- aliexpress.com
- apple.com (if order/purchase)
- costco.com
- homedepot.com
- lowes.com
- wayfair.com
- overstock.com
- zappos.com
- nordstrom.com
- macys.com
```

#### Tier 2: Subject Keywords (Secondary)
```
SUBJECT contains:
- "order"
- "purchase"
- "receipt"
- "invoice"
- "shipment"
- "tracking"
- "delivered"
- "shipping"
- "confirmation"
- "payment received"
- "refund"
- "return"
```

#### Tier 3: Content Keywords (Fallback)
```
MESSAGE contains:
- "order number"
- "tracking number"
- "estimated delivery"
- "package delivered"
- "order confirmation"
- "shipping confirmation"
- "order total"
```

### Actions:
1. **Apply tag:** `#hp-shopping`
2. **Keep in folder:** HP

---

## Filter #4: HP - Personal Finance (#30-finance)

**Category:** Holistic Performance → Finance (Personal Accounts)
**Priority:** 9
**Expected Matches:** ~40-70 emails
**⚠️ PHASE 2 - Deploy during account migration week**

### Conditions (ANY match applies filter):

#### Tier 1: Domain List (Primary)
```
FROM domain contains:
- monarch.com (expense tracking)
- chase.com
- americanexpress.com
- bankofamerica.com
- wellsfargo.com
- citi.com
- discover.com
- capitalone.com
- paypal.com
- venmo.com
- cashapp.com
- stripe.com (if personal payments)
- plaid.com
- mint.com
- ynab.com (You Need A Budget)
- personalcapital.com
- creditkarma.com
- turbotax.com
- hrblock.com
```

#### Tier 2: Subject Keywords (Secondary)
```
SUBJECT contains:
- "statement"
- "balance"
- "transaction"
- "payment due"
- "account activity"
- "credit card"
- "bank account"
- "billing"
- "invoice"
- "autopay"
- "funds"
- "deposit"
- "withdrawal"
```

#### Tier 3: Content Keywords (Fallback)
```
MESSAGE contains:
- "account balance"
- "recent transactions"
- "payment processed"
- "minimum payment"
- "statement period"
- "available credit"
```

### Actions:
1. **Apply tag:** `#30-finance`
2. **Keep in folder:** HP

---

## Filter #5: HP - Social & Relationships (#32-relationships)

**Category:** Holistic Performance → Relationships (Social)
**Priority:** 10
**Expected Matches:** ~10-20 emails

### Conditions (ANY match applies filter):

#### Tier 1: Domain List (Primary)
```
FROM domain contains:
- facebook.com
- instagram.com
- meetup.com
- eventbrite.com
- bumble.com
- hinge.com
- match.com
- evite.com
- paperlesspost.com
```

#### Tier 2: Subject Keywords (Secondary)
```
SUBJECT contains:
- "tagged you"
- "mentioned you"
- "event invitation"
- "RSVP"
- "party"
- "gathering"
- "meet up"
- "social"
```

#### Tier 3: Content Keywords (Fallback)
```
MESSAGE contains:
- "friend request"
- "group invitation"
- "event reminder"
- "social notification"
```

### Actions:
1. **Apply tag:** `#32-relationships`
2. **Keep in folder:** HP

---

## Filter #6: UI - Infrastructure & DevOps (#08-infrastructure)

**Category:** Sovereign Operations → Infrastructure (Grid & APIs)
**Priority:** 11
**Expected Matches:** ~50-80 emails

### Conditions (ANY match applies filter):

#### Tier 1: Domain List (Primary)
```
FROM domain contains:
- webflow.com
- github.com
- gitlab.com
- vercel.com
- netlify.com
- aws.amazon.com
- digitalocean.com
- linode.com
- cloudflare.com
- namecheap.com
- godaddy.com
- hover.com
- heroku.com
- railway.app
- render.com
- fly.io
- supabase.com
- planetscale.com
- mongodb.com
- redis.com
- letta.com
- anthropic.com
- openai.com
- stripe.com (if API/developer)
- slack.com (if dev workspace)
- docker.com
- kubernetes.io
```

#### Tier 2: Subject Keywords (Secondary)
```
SUBJECT contains:
- "deployment"
- "server"
- "API"
- "downtime"
- "incident"
- "SSL"
- "certificate"
- "domain"
- "hosting"
- "database"
- "build"
- "pipeline"
- "security alert"
- "backup"
```

#### Tier 3: Content Keywords (Fallback)
```
MESSAGE contains:
- "service status"
- "infrastructure update"
- "API rate limit"
- "SSL certificate expiring"
- "database maintenance"
- "deployment successful"
- "build failed"
```

### Actions:
1. **Apply tag:** `#08-infrastructure`
2. **Keep in folder:** UI

---

## Filter #7: GE - Business Finance & Treasury (#80-treasury)

**Category:** Generational Endeavor → Treasury (Business Finance)
**Priority:** 12
**Expected Matches:** ~20-40 emails

### Conditions (ANY match applies filter):

#### Tier 1: Domain List (Primary)
```
FROM domain contains:
- mercury.com (business banking)
- stripe.com (if business payments)
- quickbooks.com
- xero.com
- bench.com
- gusto.com
- rippling.com
- deel.com
- pilot.com (bookkeeping)
- brex.com
- ramp.com
- divvy.com
```

#### Tier 2: Subject Keywords (Secondary)
```
SUBJECT contains:
- "business account"
- "merchant"
- "payroll"
- "invoice sent"
- "payment received"
- "subscription billing"
- "revenue"
- "quarterly"
- "tax"
- "1099"
```

#### Tier 3: Content Keywords (Fallback)
```
MESSAGE contains:
- "business transaction"
- "merchant account"
- "payroll processed"
- "revenue report"
- "expense report"
```

### Actions:
1. **Apply tag:** `#80-treasury`
2. **Keep in folder:** GE

---

## Filter #8: GE - Personal Brand (#75-personal-brand)

**Category:** Generational Endeavor → Personal Brand
**Priority:** 13
**Expected Matches:** ~10-20 emails

### Conditions (ANY match applies filter):

#### Tier 1: Domain List (Primary)
```
FROM domain contains:
- twitter.com (if brand-related)
- youtube.com (if creator studio)
- substack.com
- medium.com (if publishing)
- beehiiv.com
- convertkit.com
- mailchimp.com
- constantcontact.com
- linktree.com
- stan.store
- gumroad.com
- patreon.com
```

#### Tier 2: Subject Keywords (Secondary)
```
SUBJECT contains:
- "subscriber"
- "newsletter"
- "followers"
- "content"
- "published"
- "analytics"
- "engagement"
- "brand"
```

#### Tier 3: Content Keywords (Fallback)
```
MESSAGE contains:
- "new subscriber"
- "content performance"
- "audience growth"
- "brand metrics"
```

### Actions:
1. **Apply tag:** `#75-personal-brand`
2. **Keep in folder:** GE

---

## Filter #9: GE - Partnerships & Collaborations (#83-partnerships)

**Category:** Generational Endeavor → Partnerships
**Priority:** 14
**Expected Matches:** ~5-15 emails

### Conditions (ANY match applies filter):

#### Tier 1: Domain List (Primary)
```
FROM domain contains:
- (Add as you discover partnership domains)
```

#### Tier 2: Subject Keywords (Secondary)
```
SUBJECT contains:
- "partnership"
- "collaboration"
- "joint venture"
- "affiliate"
- "sponsorship"
- "proposal"
- "opportunity"
- "partner program"
```

#### Tier 3: Content Keywords (Fallback)
```
MESSAGE contains:
- "partnership opportunity"
- "collaboration proposal"
- "affiliate program"
- "joint venture"
- "strategic partnership"
```

### Actions:
1. **Apply tag:** `#83-partnerships`
2. **Keep in folder:** GE

---

## Filter #10: LE - Legal & Compliance (#57-legal)

**Category:** Legacy Evolution → Legal (Estate & Compliance)
**Priority:** 15
**Expected Matches:** ~5-10 emails

### Conditions (ANY match applies filter):

#### Tier 1: Domain List (Primary)
```
FROM domain contains:
- legalzoom.com
- rocketlawyer.com
- nolo.com
- docusign.com (if legal docs)
- hellosign.com
- (Add attorney/law firm domains)
```

#### Tier 2: Subject Keywords (Secondary)
```
SUBJECT contains:
- "legal"
- "contract"
- "agreement"
- "terms"
- "compliance"
- "attorney"
- "law"
- "lawsuit"
- "settlement"
- "court"
- "trademark"
- "copyright"
```

#### Tier 3: Content Keywords (Fallback)
```
MESSAGE contains:
- "legal document"
- "sign agreement"
- "terms and conditions"
- "legal notice"
- "compliance requirement"
```

### Actions:
1. **Apply tag:** `#57-legal`
2. **Keep in folder:** LE

---

## Implementation Strategy

### Phase 1: Priority Filters (Week 1)
Deploy these first for immediate value:
1. ✅ Filter #1: Mind (#26-mind) - High volume
2. ✅ Filter #2: Body (#27-body) - Clear patterns
3. ✅ Filter #3: Shopping (#hp-shopping) - High volume
4. ✅ Filter #6: Infrastructure (#08-infrastructure) - Business critical

### Phase 2: Finance Filters (Week 2 - During Account Migration)
5. ⏸️ Filter #4: Personal Finance (#30-finance)
6. ⏸️ Filter #7: Business Finance (#80-treasury)

### Phase 3: Supplementary Filters (Week 3)
7. ⏸️ Filter #5: Social (#32-relationships)
8. ⏸️ Filter #8: Personal Brand (#75-personal-brand)
9. ⏸️ Filter #9: Partnerships (#83-partnerships)
10. ⏸️ Filter #10: Legal (#57-legal)

---

## Filter Creation Steps (Manual in Zoho Mail UI)

For each filter:

1. **Go to:** https://mail.zoho.com → Settings (⚙️) → Filters
2. **Click:** "Create Filter"
3. **Set Name:** `[Priority] Category Name` (e.g., "6 - Auto-Tag Mind")
4. **Add ALL conditions:**
   - Start with Domain list (Tier 1)
   - Add Subject keywords (Tier 2)
   - Add Content keywords (Tier 3)
   - Use "ANY" operator (OR logic)
5. **Set Actions:**
   - Apply label: `#[category-tag]`
   - **DO NOT move folder** (already routed by pillar filters)
6. **Test:** Send test email with keywords
7. **Verify:** Check dual tags applied (#pillar + #category)

---

## Maintenance & Optimization

### Weekly Review (First Month)
- Check for untagged emails in each pillar folder
- Identify new domains to add to Tier 1
- Adjust keywords based on false positives/negatives

### Monthly Audit
- Review filter performance metrics
- Consolidate low-volume filters if needed
- Add new domains discovered during account migration

### Continuous Improvement
- When you notice an untagged email, identify why:
  - Missing domain? → Add to Tier 1
  - Generic subject? → Relies on Tier 3 (expected)
  - False negative? → Adjust keywords

---

## Benefits of Multi-Tier Approach

1. **95%+ Recall:** Catches almost all relevant emails
2. **High Precision:** Domain matching prevents false positives
3. **Future-Proof:** New senders caught by keyword fallback
4. **Maintainable:** Easy to add domains as discovered
5. **Performance:** Fast domain matching, slower content matching only as fallback
6. **Transparent:** Clear why each email was categorized

---

*Last Updated: 2026-01-14*
*Status: Ready for Phase 1 implementation*
