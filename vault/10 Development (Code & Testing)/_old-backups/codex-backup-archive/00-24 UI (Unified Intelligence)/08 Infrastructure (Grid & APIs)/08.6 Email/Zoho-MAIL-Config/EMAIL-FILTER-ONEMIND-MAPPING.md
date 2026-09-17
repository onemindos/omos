---
id: root-00-20260114130000
title: "EMAIL FILTER ONEMIND MAPPING"
type: reference
status: active
created: 2026-01-14
modified: 2026-01-14
owner: zeus
tags: [type/reference, status/active, email, filters, onemind, vault-aligned]
ai_summary: "Email filter system mapping real-world email patterns to OneMind Codex vault structure (00-99)"
---
# Email Filter OneMind Mapping

**Created:** 2026-01-14
**Purpose:** Map real-world email patterns to OneMind Codex vault numbering system
**Based on:** 173 real emails analyzed + OneMind Codex structure (00-99)

---

## Filter Philosophy

`★ Insight ─────────────────────────────────────`
**OneMind-Aligned Email System**
1. **Pillar tags** (#hp, #le, #ui, #ge) route to folders
2. **Category tags** (#30-finance, #27-body) match vault structure
3. **Search power** - Find all #30-finance emails across ALL pillars
`─────────────────────────────────────────────────`

---

## HP (Holistic Performance) - Slots 25-49

### 25 Identity (Self, Purpose & Spirit)
**Tag:** `#25-identity`
**Real patterns:** Currently no distinct email patterns found
**Future keywords:** coaching, purpose, goals, meditation, journal, mindfulness

---

### 26 Mind (Education & Learning) ✅ **ACTIVE**
**Tag:** `#26-mind`
**Real patterns from 50 Education folder emails:**

**FROM domains (high confidence):**
- @substack.com (Flow State, Willy Woo, ByteByteGo newsletters)
- @skool.com (community learning platform)
- @deeplearning.ai (AI courses)
- @joinsuperhuman.ai (email productivity)
- @mail.beehiiv.com (newsletters)
- @forwardfuture.com (tech newsletter)
- @bytebytego.com (system design)
- @coursera.org
- @udemy.com
- @linkedin.com/learning

**SUBJECT keywords:**
- "course", "lesson", "webinar", "newsletter", "learning path"

**Filter configuration:**
```
IF (FROM contains ANY OF):
  - @substack.com
  - @skool.com
  - @deeplearning.ai
  - @joinsuperhuman.ai
  - @mail.beehiiv.com
  - @coursera.org
  - @udemy.com
  - @forwardfuture.com
  - @bytebytego.com
OR (SUBJECT contains ANY OF):
  - "course"
  - "lesson"
  - "webinar"
  - "newsletter"
THEN:
  - Apply tag: #26-mind
```

---

### 27 Body (Health & Wellness) ✅ **ACTIVE**
**Tag:** `#27-body`
**Real patterns from 33 Performance folder emails:**

**FROM domains (high confidence):**
- @pausestudio.com (wellness studio)
- @riseabovefloatation.com (float therapy)
- @email.alltrails.com (hiking/fitness)
- @mychart.ynhh.org (medical records)
- @equinox.com (gym membership)
- @enagic.co.jp (water wellness)
- @accesshealthct.com (health insurance)
- @fitbit.com
- @myfitnesspal.com
- @peloton.com

**SUBJECT keywords:**
- "appointment", "health record", "wellness", "workout", "membership"

**Filter configuration:**
```
IF (FROM contains ANY OF):
  - @pausestudio.com
  - @riseabovefloatation.com
  - @email.alltrails.com
  - @mychart.ynhh.org
  - @equinox.com
  - @enagic.co.jp
  - @accesshealthct.com
  - @fitbit.com
  - @myfitnesspal.com
  - @peloton.com
OR (SUBJECT contains ANY OF):
  - "appointment"
  - "health record"
  - "wellness"
  - "workout"
  - "membership"
THEN:
  - Apply tag: #27-body
```

---

### 28 Mastery (Skills & Career)
**Tag:** `#28-mastery`
**Real patterns:** Not yet identified
**Future keywords:** career, job, skills, certification, professional development

---

### 29 Life Systems (Daily Ops)
**Tag:** `#29-systems`
**Real patterns:** Not yet identified
**Future keywords:** productivity, calendar, Todoist, task, routine, workflow

---

### 30 Finance (Money) ✅ **READY**
**Tag:** `#30-finance`
**Real patterns:** No finance emails in current folders (will activate during Phase 2 account migration)

**FROM domains (high confidence - preparing):**
- @chase.com
- @bankofamerica.com
- @paypal.com
- @stripe.com
- @schwab.com
- @fidelity.com
- @mint.com
- @venmo.com
- @cashapp.com

**SUBJECT keywords:**
- "statement", "payment", "transaction", "balance", "invoice", "billing"

**Filter configuration:**
```
IF (FROM contains ANY OF):
  - @chase.com
  - @bankofamerica.com
  - @paypal.com
  - @stripe.com
  - @schwab.com
  - @fidelity.com
  - @mint.com
  - @venmo.com
  - @cashapp.com
OR (SUBJECT contains ANY OF):
  - "statement"
  - "payment"
  - "transaction"
  - "balance"
  - "invoice"
THEN:
  - Apply tag: #30-finance
```

---

### 31 Sovereignty (Legal)
**Tag:** `#31-sovereignty`
**Real patterns:** Not yet identified
**Future keywords:** legal, attorney, contract, terms, compliance

---

### 32 Relationships (Social)
**Tag:** `#32-relationships`
**Real patterns:** Found in Social folder (needs deeper analysis)
**Future keywords:** dating, social events, meetup, networking

---

### 33 Joy (Recreation & Travel) ✅ **PARTIAL**
**Tag:** `#33-joy`
**Real patterns:** Some overlap with shopping (ride-sharing)

**FROM domains:**
- @uber.com (Uber rides - recreation)
- @lyftmail.com (Lyft rides - recreation)
- Airlines, hotels, travel booking sites

**Note:** Currently grouped under shopping; consider separate travel filter

---

### 34 Spirituality (Metaphysics)
**Tag:** `#34-spirituality`
**Real patterns:** Not yet identified
**Future keywords:** spiritual, consciousness, metaphysics, philosophy

---

## HP Shopping Category ✅ **ACTIVE**
**Tag:** `#hp-shopping`
**Real patterns from 40 Shopping folder emails:**

**FROM domains (high confidence):**
- @uber.com (Uber Eats food delivery)
- @lyftmail.com (ride-sharing)
- @airsoftgi.com (shopping)
- @bemightykind.com (cannabis delivery)
- @gardencatering.com (catering orders)
- @walmart.com
- @joinsecret.com (deals)
- @doordash.com
- @grubhub.com
- @instacart.com

**SUBJECT keywords:**
- "order confirmation", "delivery", "shipment", "tracking number"

**Filter configuration:**
```
IF (FROM contains ANY OF):
  - @uber.com
  - @lyftmail.com
  - @airsoftgi.com
  - @bemightykind.com
  - @gardencatering.com
  - @walmart.com
  - @joinsecret.com
  - @doordash.com
  - @grubhub.com
  - @instacart.com
OR (SUBJECT contains ANY OF):
  - "order confirmation"
  - "delivery scheduled"
  - "shipment"
  - "tracking"
THEN:
  - Apply tag: #hp-shopping
```

---

## LE (Legacy Evolution) - Slots 50-74

### 50 Home (Sanctuary & Property)
**Tag:** `#50-home`
**Real patterns:** Not yet identified
**Future keywords:** utilities, mortgage, insurance, home repair, rent

**Expected senders:**
- Utility companies (electric, gas, water)
- Home insurance
- Mortgage lenders
- Home Depot, Lowe's

---

### 51 Marriage (Partnership & Comms)
**Tag:** `#51-marriage`
**Real patterns:** Not yet identified
**Future keywords:** spouse, anniversary, counseling

---

### 52 Children (Parenting & Education)
**Tag:** `#52-children`
**Real patterns:** Not yet identified
**Future keywords:** school, daycare, pediatrician, kids activities

**Expected senders:**
- Schools
- Daycare centers
- Pediatricians
- Kids activity providers

---

### 53 Parents (Elder Care & Support)
**Tag:** `#53-parents`
**Real patterns:** Not yet identified
**Future keywords:** elder care, senior services, nursing

---

### 54 Heritage (History & Traditions)
**Tag:** `#54-heritage`
**Real patterns:** Not yet identified
**Future keywords:** ancestry, genealogy, family history

**Expected senders:**
- Ancestry.com
- 23andMe

---

### 55 Health (Family Wellness)
**Tag:** `#55-family-health`
**Real patterns:** Not yet identified (distinct from personal health #27-body)
**Future keywords:** family doctor, family health insurance

---

### 56 Wealth (Estate & Assets)
**Tag:** `#56-wealth`
**Real patterns:** Not yet identified
**Future keywords:** estate, trust, will, wealth management, retirement

---

### 57 Legal (Documents)
**Tag:** `#57-family-legal`
**Real patterns:** Not yet identified
**Future keywords:** estate planning, power of attorney, custody

---

### 58 Community (Extended Network)
**Tag:** `#58-community`
**Real patterns:** Not yet identified
**Future keywords:** neighborhood, HOA, local community, volunteer

**Expected senders:**
- HOA communications
- Nextdoor app
- Local government

---

## UI (Unified Intelligence) - Slots 00-24

### 00 Framework (Doctrine)
**Tag:** `#00-framework`
**Real patterns:** Not yet identified
**Future keywords:** architecture, system design, framework

---

### 01 Command (Dashboard & Navigation)
**Tag:** `#01-command`
**Real patterns:** Not yet identified
**Future keywords:** dashboard, metrics, analytics, monitoring

---

### 02 Agents (AI, Memory & Tools)
**Tag:** `#02-agents`
**Real patterns:** Some found in Education folder

**FROM domains:**
- @anthropic.com
- @openai.com
- @letta.com
- @deeplearning.ai (also #26-mind)

**Future expansion:** Separate AI service emails from AI learning content

---

### 03 Protocols (SOPs)
**Tag:** `#03-protocols`
**Real patterns:** Not yet identified
**Future keywords:** SOP, procedure, protocol, workflow

---

### 04 Resources (Assets & References)
**Tag:** `#04-resources`
**Real patterns:** Not yet identified
**Future keywords:** template, asset, reference, documentation

---

### 05 Environments (Contexts & Locations)
**Tag:** `#05-environments`
**Real patterns:** Not yet identified
**Future keywords:** workspace, coworking, office

---

### 06 Inbox (Queue)
**Tag:** `#06-inbox`
**Use case:** Temporary triage tag

---

### 08 Infrastructure (Grid & APIs) ✅ **ACTIVE**
**Tag:** `#08-infrastructure`
**Real patterns from 50 Legacyhq - Old folder emails:**

**FROM domains (high confidence):**
- @amazon.com (AWS support)
- @aws.amazon.com
- @mercury.com (business banking/infrastructure)
- @github.com
- @microsoft.com
- @azure.com
- @nvidia.com (hardware/tech)
- @verizon.com (telecom billing)
- @cloudflare.com
- @digitalocean.com
- @heroku.com
- @vercel.com

**SUBJECT keywords:**
- "support case", "billing", "invoice", "service update", "security alert"

**Filter configuration:**
```
IF (FROM contains ANY OF):
  - @amazon.com
  - @aws.amazon.com
  - @mercury.com
  - @github.com
  - @microsoft.com
  - @azure.com
  - @nvidia.com
  - @verizon.com
  - @cloudflare.com
  - @digitalocean.com
  - @heroku.com
  - @vercel.com
OR (SUBJECT contains ANY OF):
  - "support case"
  - "billing"
  - "invoice"
  - "service update"
  - "security alert"
THEN:
  - Apply tag: #08-infrastructure
```

---

### 09 Automation (Workflows)
**Tag:** `#09-automation`
**Real patterns:** Not yet identified
**Future keywords:** N8N, Zapier, automation, workflow, webhook

---

### 10 Development (Code & Testing)
**Tag:** `#10-development`
**Real patterns:** Overlap with #08-infrastructure (GitHub)

**FROM domains:**
- @github.com (already in #08-infrastructure)
- @gitlab.com
- CI/CD services

**Future:** May merge with #08-infrastructure or keep separate for dev-specific emails

---

### 11 Security (Cyber)
**Tag:** `#11-security`
**Real patterns:** Not yet identified
**Future keywords:** security, breach, password, 2FA, vulnerability

**Expected senders:**
- @1password.com
- Security alert services
- Vulnerability scanners

---

### 12 Research (R&D)
**Tag:** `#12-research`
**Real patterns:** Some overlap with #26-mind (newsletters)

**FROM domains:**
- Tech research newsletters
- Academic journals
- Innovation platforms

**Future:** Distinguish between learning content (#26-mind) and research/trends

---

### 24 Archive (Historical)
**Tag:** `#24-archive`
**Use case:** Historical reference emails

---

## GE (Generational Entrepreneurship) - Slots 75-99

### 75 Zeus Delacruz Enterprise (Personal Brand)
**Tag:** `#75-zeus-enterprise`
**Real patterns:** Not yet identified
**Future keywords:** zeus delacruz, speaking, consulting, press, media

---

### 76 Legacy Acadacom (Education)
**Tag:** `#76-acadacom`
**Real patterns:** Not yet identified
**Future keywords:** acadacom, course creation, students

---

### 77 Legacy Platform Holdings (IP)
**Tag:** `#77-ip`
**Real patterns:** Not yet identified
**Future keywords:** trademark, copyright, patent, licensing

---

### 78 Legacy Restore (Credit) ✅ **READY**
**Tag:** `#78-legacy-restore`
**Real patterns:** Found in Legacyhq folders

**FROM domains:**
- @legacyrestore.co
- @legacyhq.ai (old domain)
- Client communications
- Credit bureau emails

**Filter configuration:**
```
IF (FROM contains ANY OF):
  - @legacyrestore.co
  - @legacyhq.ai
OR (SUBJECT contains ANY OF):
  - "credit repair"
  - "credit report"
  - "dispute"
  - "credit bureau"
THEN:
  - Apply tag: #78-legacy-restore
```

---

### 79 Technology (Tech Stack)
**Tag:** `#79-tech-stack`
**Real patterns:** Overlap with #08-infrastructure
**Future keywords:** SaaS, software subscription, tools

---

### 80 Treasury (Finance)
**Tag:** `#80-treasury`
**Real patterns:** Not yet identified (distinct from personal #30-finance)
**Future keywords:** business bank, QuickBooks, business credit, revenue

**Expected senders:**
- Business bank accounts
- QuickBooks
- Business credit cards

---

### 81 Strategy (Analytics)
**Tag:** `#81-strategy`
**Real patterns:** Not yet identified
**Future keywords:** analytics, market research, strategy, insights

---

### 83 Partnerships (Network)
**Tag:** `#83-partnerships`
**Real patterns:** Not yet identified
**Future keywords:** partner, affiliate, collaboration, JV

---

### 84 Ventures (New Ops & Innovation)
**Tag:** `#84-ventures`
**Real patterns:** Not yet identified
**Future keywords:** venture, startup, new business, opportunity

---

## Priority Filter Implementation

### Phase 1: Active Filters (4 filters - IMPLEMENT NOW)
**Based on 173 real emails with proven patterns:**

1. **Filter #6: Auto-Tag Mind/Learning** → `#26-mind`
   - 50 Education folder emails
   - Substack, Skool, DeepLearning.AI patterns
   - 95% accuracy expected

2. **Filter #7: Auto-Tag Body/Health** → `#27-body`
   - 33 Performance folder emails
   - Pause Studio, Equinox, MyChart patterns
   - 90% accuracy expected

3. **Filter #8: Auto-Tag Shopping** → `#hp-shopping`
   - 40 Shopping folder emails
   - Uber, Walmart, DoorDash patterns
   - 85% accuracy expected

4. **Filter #9: Auto-Tag Infrastructure** → `#08-infrastructure`
   - 50 Legacyhq folder emails
   - AWS, GitHub, Verizon patterns
   - 80% accuracy expected

---

### Phase 2: Ready Filters (2 filters - IMPLEMENT DURING ACCOUNT MIGRATION)

5. **Filter #10: Auto-Tag Finance** → `#30-finance`
   - No current emails (activate when migrating bank accounts to hp@)
   - Chase, PayPal, Stripe patterns prepared
   - Will activate in Week 2 of account migration

6. **Filter #11: Auto-Tag Legacy Restore** → `#78-legacy-restore`
   - legacyrestore.co and legacyhq.ai patterns
   - Activate when business emails start arriving at ge@

---

### Phase 3: Future Filters (30+ filters - MAP WHEN EMAILS APPEAR)

Remaining categories will be created as emails naturally arrive in those categories during:
- Account migrations (Weeks 2-4)
- Daily email flow (Month 2+)
- New service signups

**Strategy:** Don't create filters for empty categories. Wait for real email patterns, then map to OneMind Codex numbering.

---

## Filter Creation Commands (Zoho MCP)

### Command Template
```
Create filter via MCP:
- name: "Auto-Tag [Category]"
- criteria: Match ANY (FROM contains OR SUBJECT contains)
- action: Apply tag #XX-category
- priority: [6-50]
- execute_after: pillar routing (filters 1-5)
```

### Implementation Order
```
Priority 1-5: Pillar routing (ALREADY DONE ✅)
Priority 6: Auto-Tag Mind (#26-mind)
Priority 7: Auto-Tag Body (#27-body)
Priority 8: Auto-Tag Shopping (#hp-shopping)
Priority 9: Auto-Tag Infrastructure (#08-infrastructure)
Priority 10: Auto-Tag Finance (#30-finance) - READY, activate in Week 2
Priority 11: Auto-Tag Legacy Restore (#78-legacy-restore) - READY
Priority 12-50: Future categories (create as needed)
```

---

## Success Metrics

**Phase 1 Target (End of Week 1):**
- ✅ 4 pillar folders created (HP, LE, UI, GE)
- ✅ 5 routing filters working (HP, LE, UI, GE + catch-all)
- ⏳ 4 category filters working (#26-mind, #27-body, #hp-shopping, #08-infrastructure)
- ⏳ 173 existing emails retroactively tagged

**Phase 2 Target (End of Month 1):**
- 6 total category filters (#30-finance, #78-legacy-restore added)
- 90% of incoming emails auto-tagged
- Zero manual folder sorting required

**Phase 3 Target (End of Month 2):**
- 15+ category filters covering all common email types
- Complete OneMind Codex alignment (00-99 numbering)
- Smart search working across all categories and pillars

---

## Next Steps

1. ✅ Review this mapping (YOU)
2. ⏳ Implement Phase 1 filters via Zoho MCP (ME)
3. ⏳ Retroactively tag 173 existing emails (ME)
4. ⏳ Test with new emails (WE)
5. ⏳ Activate Phase 2 filters during account migration (WEEK 2)

---

*Document created: 2026-01-14*
*Merges real-world patterns (173 emails) with OneMind Codex structure (00-99)*
*Ready for Zoho MCP implementation*
