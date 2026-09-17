---
id: root-00-20260114160000
title: "EMAIL SMART FILTERS - DUAL TAG SYSTEM"
type: reference
status: active
created: 2026-01-14
modified: 2026-01-14
owner: zeus
tags: [type/reference, status/active, email, filters, onemind, dual-tag]
ai_summary: "Smart consolidated email filters with dual-tag system (path + category) aligned with OneMind Codex structure (00-99)"
---
# Email Smart Filters - Dual Tag System

**Created:** 2026-01-14
**Purpose:** Consolidated filters with dual-tag system (path + category) for OneMind cross-system alignment
**Strategy:** SMART filters that combine related categories, apply 2 tags + move to pillar folder

---

## System Architecture

`★ Insight ─────────────────────────────────────`
**Dual-Tag Email System:**
1. **Path tag** (#hp/#le/#ui/#ge) → Domain (25-49, 50-74, 00-24, 75-99)
2. **Category tag** (#26-mind, #30-finance, #08-infrastructure) → Slot (00-99)
3. **Folder** → Physical organization (4 pillar folders)
4. **Cross-system alignment** → Same tags in 1Password, Todoist, GitHub
`─────────────────────────────────────────────────`

### Benefits
- ✅ **Domain search:** "Show ALL hp emails" → `tag:#hp`
- ✅ **Category search:** "Show ALL finance" → `tag:#30-finance` OR `tag:#80-business-finance`
- ✅ **Combined search:** "Show HP finance only" → `tag:#hp tag:#30-finance`
- ✅ **Cross-platform:** Same numbering in email, 1Password, Todoist, GitHub

---

## Filter Priority Order

```
Filters 1-5: Pillar Routing (ALREADY EXIST ✅)
├─ Priority 1: hp@zeusdelacruz.me → HP folder + #hp
├─ Priority 2: le@zeusdelacruz.me → LE folder + #le
├─ Priority 3: ui@zeusdelacruz.me → UI folder + #ui
├─ Priority 4: ge@zeusdelacruz.me → GE folder + #ge
└─ Priority 5: *@zeusdelacruz.me → UI folder + #catchall

Filters 6-11: Smart Category Filters (NEW - DUAL TAG)
├─ Priority 6: HP - Mind/Learning → HP folder + #hp + #26-mind
├─ Priority 7: HP - Body/Health → HP folder + #hp + #27-body
├─ Priority 8: HP - Shopping → HP folder + #hp + #hp-shopping
├─ Priority 9: HP - Finance (Personal) → HP folder + #hp + #30-finance
├─ Priority 10: HP - Social → HP folder + #hp + #32-relationships
├─ Priority 11: UI - Infrastructure → UI folder + #ui + #08-infrastructure
├─ Priority 12: GE - Business Finance → GE folder + #ge + #80-business-finance
├─ Priority 13: GE - Personal Brand → GE folder + #ge + #75-personal-brand
├─ Priority 14: GE - Partnerships → GE folder + #ge + #83-partnerships
├─ Priority 15: LE - Legal → LE folder + #le + #57-legal
└─ Priority 16: Cleanup Old Aliases → Delete old labels, redirect emails
```

---

## SMART FILTER SPECIFICATIONS

### 🔵 HP DOMAIN FILTERS (Priorities 6-10)

---

## FILTER #6: HP - Mind/Learning

### Configuration
- **Filter Name:** `HP - Mind (26)`
- **Priority:** 6
- **Status:** Active

### Criteria (Match ANY)
**FROM contains:**
- `substack.com`
- `skool.com`
- `deeplearning.ai`
- `joinsuperhuman.ai`
- `mail.beehiiv.com`
- `coursera.org`
- `udemy.com`
- `forwardfuture.com`
- `bytebytego.com`
- `linkingyourthinking.com`

**OR SUBJECT contains:**
- `course`
- `lesson`
- `webinar`
- `newsletter`
- `learning`
- `workshop`
- `certification`

### Actions (3 actions)
1. **Apply tag:** `#hp` (path)
2. **Apply tag:** `#26-mind` (category)
3. **Move to folder:** `25-49 HP (Holistic Performance)`

### Expected Results
- **~50 emails** will match
- **Replaces old label:** 8698543000000328109 (education@)
- **Cross-system:** Matches `category:learning` in 1Password

---

## FILTER #7: HP - Body/Health

### Configuration
- **Filter Name:** `HP - Body (27)`
- **Priority:** 7
- **Status:** Active

### Criteria (Match ANY)
**FROM contains:**
- `pausestudio.com`
- `mychart.ynhh.org`
- `accesshealthct.com`
- `email.alltrails.com`
- `riseabovefloatation.com`
- `equinox.com`
- `enagic.co.jp`
- `fitbit.com`
- `myfitnesspal.com`
- `peloton.com`
- `huskyhealthct.org`

**OR SUBJECT contains:**
- `appointment`
- `health record`
- `wellness`
- `workout`
- `membership`
- `medical`
- `fitness`
- `doctor`

### Actions (3 actions)
1. **Apply tag:** `#hp` (path)
2. **Apply tag:** `#27-body` (category)
3. **Move to folder:** `25-49 HP (Holistic Performance)`

### Expected Results
- **~33 emails** will match
- **Replaces old label:** 8698543000000328119 (performance@)
- **Cross-system:** Matches `category:health` in 1Password

---

## FILTER #8: HP - Shopping

### Configuration
- **Filter Name:** `HP - Shopping`
- **Priority:** 8
- **Status:** Active

### Criteria (Match ANY)
**FROM contains:**
- `uber.com`
- `lyftmail.com`
- `doordash.com`
- `grubhub.com`
- `instacart.com`
- `walmart.com`
- `bemightykind.com`
- `gardencatering.com`
- `airsoftgi.com`
- `joinsecret.com`
- `amazon.com` (order emails only)

**OR SUBJECT contains:**
- `order confirmation`
- `delivery scheduled`
- `shipment`
- `tracking`
- `your order`
- `receipt`

### Actions (3 actions)
1. **Apply tag:** `#hp` (path)
2. **Apply tag:** `#hp-shopping` (category - HP-specific, no slot number)
3. **Move to folder:** `25-49 HP (Holistic Performance)`

### Expected Results
- **~40 emails** will match
- **Replaces old label:** 8698543000000328121 (shopping@)
- **Cross-system:** Matches `category:shopping` in 1Password

---

## FILTER #9: HP - Finance (Personal)

### Configuration
- **Filter Name:** `HP - Finance (30)`
- **Priority:** 9
- **Status:** Active

### Criteria (Match ANY)
**FROM contains:**
- `monarch.com`
- `creditkarma.com`
- `creditsesame.com`
- `myfico.com`
- `chase.com`
- `bankofamerica.com`
- `paypal.com`
- `stripe.com`
- `venmo.com`
- `cashapp.com`
- `experianidworks.com`
- `meetava.com`
- `servicecu.org`

**OR SUBJECT contains:**
- `statement`
- `payment`
- `balance`
- `transaction`
- `credit score`
- `account alert`
- `bill due`

**AND NOT FROM:** (Exclude business finance)
- `mercury.com`
- `legacyhq.ai`

### Actions (3 actions)
1. **Apply tag:** `#hp` (path)
2. **Apply tag:** `#30-finance` (category)
3. **Move to folder:** `25-49 HP (Holistic Performance)`

### Expected Results
- **Activates during account migration**
- **Replaces old label:** 8698543000000328111 (finance@zeusdelacruz.me)
- **Cross-system:** Matches `category:finance` in 1Password HP vault

---

## FILTER #10: HP - Social/Relationships

### Configuration
- **Filter Name:** `HP - Social (32)`
- **Priority:** 10
- **Status:** Active

### Criteria (Match ANY)
**FROM contains:**
- `linkedin.com`
- `discord.com`
- `facebook.com`
- `instagram.com`
- `twitter.com`
- `meetup.com`
- `eventbrite.com`

**OR SUBJECT contains:**
- `message`
- `connection request`
- `invitation`
- `event`
- `networking`

### Actions (3 actions)
1. **Apply tag:** `#hp` (path)
2. **Apply tag:** `#32-relationships` (category)
3. **Move to folder:** `25-49 HP (Holistic Performance)`

### Expected Results
- **~5-10 emails** will match
- **Replaces old label:** 8698543000000439001 (social@)
- **Cross-system:** Matches `category:social` in 1Password

---

### 🔵 UI DOMAIN FILTER (Priority 11)

---

## FILTER #11: UI - Infrastructure

### Configuration
- **Filter Name:** `UI - Infrastructure (08)`
- **Priority:** 11
- **Status:** Active

### Criteria (Match ANY)
**FROM contains:**
- `aws.amazon.com`
- `no-reply-aws@amazon.com`
- `github.com`
- `microsoft.com`
- `azure.com`
- `openai.com`
- `anthropic.com`
- `cloudflare.com`
- `digitalocean.com`
- `vercel.com`
- `heroku.com`
- `nvidia.com`
- `verizon.com`
- `google.com` (cloud/workspace only)
- `toolhouseai.com`
- `crewai.com`
- `heygen.com`
- `webflow.com`

**OR SUBJECT contains:**
- `support case`
- `service update`
- `security alert`
- `quota increase`
- `billing` (tech services)
- `API`
- `deployment`

### Actions (3 actions)
1. **Apply tag:** `#ui` (path)
2. **Apply tag:** `#08-infrastructure` (category)
3. **Move to folder:** `00-24 UI (Unified Intelligence)`

### Expected Results
- **~50+ emails** will match
- **Replaces old labels:** 8698543000000351303 (zeus.ops@) + 8698543000000278003 (techops@)
- **Cross-system:** Matches `category:infrastructure` in 1Password

---

### 🟢 GE DOMAIN FILTERS (Priorities 12-14)

---

## FILTER #12: GE - Business Finance

### Configuration
- **Filter Name:** `GE - Treasury (80)`
- **Priority:** 12
- **Status:** Active

### Criteria (Match ANY)
**FROM contains:**
- `mercury.com`
- `@legacyhq.ai` (if contains "finance")
- `companysage.com`
- `consumerdirect.com` (partner billing)

**OR SUBJECT contains:**
- `business account`
- `virtual debit card`
- `merchant services`
- `1099`
- `business transaction`

### Actions (3 actions)
1. **Apply tag:** `#ge` (path)
2. **Apply tag:** `#80-business-finance` (category)
3. **Move to folder:** `75-99 GE (Generational Entrepreneurship)`

### Expected Results
- **Business finance only**
- **Replaces old label:** 8698543000000189019 (finance@legacyhq.ai)
- **Cross-system:** Matches `category:finance` in 1Password GE vault

---

## FILTER #13: GE - Personal Brand

### Configuration
- **Filter Name:** `GE - Personal Brand (75)`
- **Priority:** 13
- **Status:** Active

### Criteria (Match ANY)
**FROM contains:**
- `smartpassiveincome.com`
- `andyaudate.com`
- `imiloainstitute.com`

**OR (TO contains):**
- `founder@` (old alias)

**OR SUBJECT contains:**
- `speaking`
- `consulting`
- `interview`
- `partnership opportunity`
- `collaboration`

### Actions (3 actions)
1. **Apply tag:** `#ge` (path)
2. **Apply tag:** `#75-personal-brand` (category)
3. **Move to folder:** `75-99 GE (Generational Entrepreneurship)`

### Expected Results
- **Personal brand emails**
- **Replaces old label:** 8698543000000351315 (founder@legacyhq.ai)
- **Cross-system:** Matches `category:marketing` or new `category:personal-brand` in 1Password

---

## FILTER #14: GE - Partnerships

### Configuration
- **Filter Name:** `GE - Partnerships (83)`
- **Priority:** 14
- **Status:** Active

### Criteria (Match ANY)
**FROM contains:**
- `consumerdirect.com`
- `@legacyhq.ai` (if contains "partner")

**OR (TO contains):**
- `partnernetwork@`

**OR SUBJECT contains:**
- `partner`
- `affiliate`
- `referral`
- `commission`
- `network`

### Actions (3 actions)
1. **Apply tag:** `#ge` (path)
2. **Apply tag:** `#83-partnerships` (category)
3. **Move to folder:** `75-99 GE (Generational Entrepreneurship)`

### Expected Results
- **Partnership/affiliate emails**
- **Replaces old label:** 8698543000000189049 (partnernetwork@legacyhq.ai)
- **Cross-system:** Matches `category:partnernetwork` in 1Password

---

### 🟡 LE DOMAIN FILTER (Priority 15)

---

## FILTER #15: LE - Legal

### Configuration
- **Filter Name:** `LE - Legal (57)`
- **Priority:** 15
- **Status:** Active

### Criteria (Match ANY)
**FROM contains:**
- `companysage.com`
- `@legacyhq.ai` (if contains "legal")

**OR (TO contains):**
- `legal@`

**OR SUBJECT contains:**
- `LLC`
- `dissolution`
- `compliance`
- `registered agent`
- `corporate`
- `legal notice`

### Actions (3 actions)
1. **Apply tag:** `#le` (path)
2. **Apply tag:** `#57-legal` (category)
3. **Move to folder:** `50-74 LE (Legacy Evolution)`

### Expected Results
- **Business legal/compliance emails**
- **Replaces old label:** 8698543000000189009 (legal@legacyhq.ai)
- **Cross-system:** Matches `category:legal` in 1Password

---

## Implementation Checklist

### Step 1: Create Filters (15-20 minutes)
```
☐ Filter #6: HP - Mind (26)
☐ Filter #7: HP - Body (27)
☐ Filter #8: HP - Shopping
☐ Filter #9: HP - Finance (30)
☐ Filter #10: HP - Social (32)
☐ Filter #11: UI - Infrastructure (08)
☐ Filter #12: GE - Business Finance (80)
☐ Filter #13: GE - Personal Brand (75)
☐ Filter #14: GE - Partnerships (83)
☐ Filter #15: LE - Legal (57)
```

### Step 2: Test Each Filter (10 minutes)
```
☐ Send test email to hp@ with "course" → Should get #hp + #26-mind + HP folder
☐ Send test email to hp@ with "appointment" → Should get #hp + #27-body + HP folder
☐ Send test email to ui@ from AWS → Should get #ui + #08-infrastructure + UI folder
☐ Send test email to ge@ from mercury.com → Should get #ge + #80-business-finance + GE folder
```

### Step 3: Verify Old Labels Replaced
```
☐ Search emails with old label 8698543000000328109 → Should now have #26-mind
☐ Search emails with old label 8698543000000328119 → Should now have #27-body
☐ Search emails with old label 8698543000000351303 → Should now have #08-infrastructure
☐ All 16 old labels remapped to OneMind tags
```

### Step 4: Bulk Migrate Old Emails
```
☐ Move all Education folder emails → HP (filters auto-tag)
☐ Move all Performance folder emails → HP (filters auto-tag)
☐ Move all Shopping folder emails → HP (filters auto-tag)
☐ Move all Finance folder emails → HP or GE (filters auto-tag)
☐ Move all Legacyhq folder emails → UI or GE (filters auto-tag)
```

### Step 5: Delete Old Folders (After verification)
```
☐ Delete "Zeus Delacruz" folder (13 subfolders)
☐ Delete "Legacyhq - Old" folder
☐ Delete "legacyhq - New" folder (in Trash)
☐ Keep only 4 pillar folders: HP, LE, UI, GE
```

---

## Search Examples

### Domain Searches
```
tag:#hp                  → ALL personal emails (learning, health, shopping, finance, social)
tag:#ui                  → ALL tech emails (infrastructure, AI, development)
tag:#ge                  → ALL business emails (finance, brand, partnerships)
tag:#le                  → ALL family/legal emails
```

### Category Searches (Cross-Domain)
```
tag:#26-mind             → ALL learning emails (wherever they are)
tag:#30-finance          → Personal finance only
tag:#80-business-finance → Business finance only
tag:#08-infrastructure   → ALL infrastructure emails
```

### Combined Searches
```
tag:#hp tag:#30-finance  → Personal finance in HP domain
tag:#ge tag:#80-business-finance → Business finance in GE domain
tag:#ui tag:#08-infrastructure → Tech infrastructure in UI domain
```

### Cross-System Alignment
```
EMAIL: tag:#30-finance
1PASSWORD: path:hp category:finance
TODOIST: #30-finance
GITHUB: label:30-finance
```

---

## Success Metrics

**When complete, you'll have:**
- ✅ **4 pillar folders** (HP, LE, UI, GE) - all emails organized
- ✅ **10 smart filters** - dual-tag every email (path + category)
- ✅ **16 old labels** → replaced with OneMind tags (00-99)
- ✅ **Cross-system search** - same tags in email, 1Password, Todoist, GitHub
- ✅ **Zero manual sorting** - all emails auto-filed and auto-tagged

---

*Document created: 2026-01-14*
*Ready for implementation in Zoho Mail*
*Estimated time: 30-40 minutes (create + test + migrate)*
