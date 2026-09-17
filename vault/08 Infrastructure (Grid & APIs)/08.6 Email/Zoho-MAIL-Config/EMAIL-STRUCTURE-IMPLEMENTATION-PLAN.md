---
id: root-00-20260113225954
title: "EMAIL STRUCTURE IMPLEMENTATION PLAN"
type: reference
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/reference, status/active]
ai_summary: "EMAIL STRUCTURE IMPLEMENTATION PLAN"
---
# Email Structure Implementation Plan

**Last Updated:** January 13, 2026  
**Status:** In Progress  
**Owner:** Zeus Delacruz + Legacy AI

---

## Executive Summary

Restructuring email from 23+ scattered addresses to a simplified 4-pillar system aligned with OneMind Codex framework. Uses aliases + folders + tags instead of multiple paid users.

**Cost:** $45/month (existing Zoho One subscription - no additional cost)

---

## Section 1: Current State

### Primary Inbox
```
zeus@onemindcodex.com    → Main mailbox (Zoho One admin)
```

### Support Group
```
support@onemindcodex.com → Group (connected to Zoho Dtesk)
```

### Existing Aliases on zeusdelacruz.me (11 total)
```
education@zeusdelacruz.me
family@zeusdelacruz.me
finance@zeusdelacruz.me
legacy@zeusdelacruz.me
legal@zeusdelacruz.me
performance@zeusdelacruz.me
personal@zeusdelacruz.me
recreation@zeusdelacruz.me
shopping@zeusdelacruz.me
social@zeusdelacruz.me
training@zeusdelacruz.me
travel@zeusdelacruz.me
```

### Existing Aliases on legacyhq.ai (7 total)
```
finance@legacyhq.ai
founder@legacyhq.ai
legal@legacyhq.ai
partnernetwork@legacyhq.ai
zeus.ops@legacyhq.ai
(plus others)
```

### Security Vault (Proton - Separate)
```
zeus.vaults@pm.me        → Critical accounts (Cloudflare, GoDaddy, Apple)
command@zeusops.dev      → Dev access (GitHub, Letta, AWS)
```

### Problems with Current State
- 18+ aliases all dump into ONE inbox
- No pillar separation
- Can't identify what needs updating
- Legacy domain (legacyhq.ai) has business emails that should migrate
- No clear organization system

---

## Section 2: Target State

### Primary Inbox (unchanged)
```
zeus@onemindcodex.com    → Main mailbox where everything lands
```

### Support Group (unchanged)
```
support@onemindcodex.com → Group for customer support (Zoho Desk)
```

### NEW Pillar Aliases on zeusdelacruz.me (4 total)
```
hp@zeusdelacruz.me       → Holistic Performance
le@zeusdelacruz.me       → Legacy Evolution
ui@zeusdelacruz.me       → Unified Intelligence (shared with Legacy AI)
ge@zeusdelacruz.me       → Generational Entrepreneurship
```

### Personal Brand Alias
```
zeus@zeusdelacruz.com    → Personal brand (speaking, consulting, press)
                           Routes to main inbox
```

### Legacy Restore (Temporary Business)
```
support@legacyrestore.co → Client communication and documents
                           Routes to support@onemindcodex.com group
```

### Security Vault (No Change)
```
zeus.vaults@pm.me        → Keep on Proton
command@zeusops.dev      → Keep on Proton
```

---

## Section 3: Folder Structure

### Simplified Folder Setup (5 folders only)

```
zeus@onemindcodex.com (main inbox)
│
├── Inbox                          ← Triage zone
├── Holistic Performance           ← hp@ emails
├── Legacy Evolution               ← le@ emails
├── Unified Intelligence           ← ui@ emails
└── Generational Entrepreneurship  ← ge@ emails
```

### What Goes Where

| Folder | Alias | Content |
|--------|-------|---------|
| Inbox | zeus@onemindcodex.com, zeus@zeusdelacruz.com | Triage daily |
| Holistic Performance | hp@zeusdelacruz.me | Banks, shopping, health, learning |
| Legacy Evolution | le@zeusdelacruz.me | Family, travel, legal |
| Unified Intelligence | ui@zeusdelacruz.me | Infrastructure, AI, dev |
| Generational Entrepreneurship | ge@zeusdelacruz.me | Business, affiliates |

---

## Section 4: Filter Configuration

### Permanent Filters (4 total)

**Filter 1: Route HP**
```
Criteria:  To contains "hp@zeusdelacruz.me"
Actions:   Move to folder: Holistic Performance
           Apply tag: #hp
```

**Filter 2: Route LE**
```
Criteria:  To contains "le@zeusdelacruz.me"
Actions:   Move to folder: Legacy Evolution
           Apply tag: #le
```

**Filter 3: Route UI**
```
Criteria:  To contains "ui@zeusdelacruz.me"
Actions:   Move to folder: Unified Intelligence
           Apply tag: #ui
```

**Filter 4: Route GE**
```
Criteria:  To contains "ge@zeusdelacruz.me"
Actions:   Move to folder: Generational Entrepreneurship
           Apply tag: #ge
```

### Transition Filters (Temporary)

**Filter 5: LegacyHQ Transition**
```
Criteria:  To contains "@legacyhq.ai"
Actions:   Move to folder: Unified Intelligence
           Apply tag: #legacyhq-transition
```

**Filter 6: Old Aliases Tracker**
```
Criteria:  To contains "@zeusdelacruz.me" AND NOT (hp@, le@, ui@, ge@)
Actions:   Apply tag: #migration
           Stay in Inbox
```

---

## Section 5: Catch-All Setup

### Configuration

**Domain: zeusdelacruz.me**
```
Catch-all: YES
Forward to: zeus@onemindcodex.com
```

**Domain: legacyhq.ai**
```
Catch-all: YES
Forward to: zeus@onemindcodex.com
```

### How to Enable in Zoho
```
1. admin.zoho.com → Mail Admin
2. Domains → Select domain
3. Catch-All Address → Enable
4. Forward to: zeus@onemindcodex.com
5. Save
```

---

## Section 6: Tag System

### Pillar Tags (Auto-applied by filters)
```
#hp   → Holistic Performance emails
#le   → Legacy Evolution emails
#ui   → Unified Intelligence emails
#ge   → Generational Entrepreneurship emails
```

### Category Tags (Manual or smart filters)

**HP Tags:**
```
#finance     → Banks, investments, credit cards
#health      → Medical, fitness, wellness
#shopping    → E-commerce, retail
#learning    → Courses, education
```

**LE Tags:**
```
#family      → Wife, kids, relatives
#travel      → Airlines, hotels, bookings
#legal       → Legal documents, attorneys
```

**UI Tags:**
```
#infrastructure → Servers, hosting, domains
#ai-services    → OpenAI, Anthropic, Letta
#dev            → GitHub, tools
```

**GE Tags:**
```
#onemind     → OneMind Codex business
#restore     → Legacy Restore (temporary)
#affiliates  → Affiliate programs
#partners    → Partner inquiries
```

### Transition Tags (Temporary)
```
#migration           → Old aliases needing update
#legacyhq-transition → LegacyHQ domain emails
```

### How Tags Work

**Click tag → See all emails with that tag (across folders)**

Example:
- Click #finance → See ALL finance emails (HP + GE)
- Click #family → See ALL family emails in one view

Tags are more flexible than folders. One email can have multiple tags.

---

## Section 7: Transition Plan

### Phase 1: Foundation (Week 1)

**Day 1-2: Create Structure**
- [ ] Create 4 pillar aliases (hp@, le@, ui@, ge@)
- [ ] Create 4 pillar folders
- [ ] Create 4 routing filters
- [ ] Enable catch-all on zeusdelacruz.me
- [ ] Enable catch-all on legacyhq.ai

**Day 3-4: Configure Domains**
- [ ] Add zeusdelacruz.com to Zoho (if not done)
- [ ] Add legacyrestore.co to Zoho
- [ ] Create zeus@zeusdelacruz.com alias
- [ ] Create support@legacyrestore.co alias
- [ ] Configure MX records

**Day 5-7: Test**
- [ ] Send test emails to each pillar alias
- [ ] Verify routing to correct folders
- [ ] Verify tags applied correctly

### Phase 2: Account Migration (Weeks 2-4)

**Week 2: Financial (Update to hp@)**
- Banks
- Credit cards
- Investment accounts
- PayPal, Stripe

**Week 3: Family & Personal (Update to le@)**
- School accounts
- Kids activities
- Travel accounts
- Family subscriptions

**Week 4: Business (Update to ge@)**
- Affiliate programs
- MLM platforms
- Business subscriptions
- Partner accounts

**Week 4: Infrastructure (Update to ui@)**
- Non-critical hosting
- AI services
- Dev tools (keep critical on Proton)

### Phase 3: Monitor (Month 2)

**Weekly:**
- [ ] Check #migration tag (what's still using old aliases?)
- [ ] Update remaining accounts
- [ ] Track which old aliases still get mail

**End of Month 2:**
- [ ] Review dormant old aliases
- [ ] Plan cleanup phase

### Phase 4: Cleanup (Month 3)

**Delete dormant aliases:**
- [ ] Old zeusdelacruz.me aliases (education@, finance@, etc.)
- [ ] All legacyhq.ai aliases
- [ ] Remove transition filters

**Rename old folders:**
- [ ] "Zeus Delacruz" → "Archive - Pre-2026"
- [ ] Delete legacyhq folders

---

## Section 8: Implementation Steps (Claude Code Ready)

### Step 1: Create Pillar Folders

```
Location: zeus@onemindcodex.com inbox
Create folders:
1. Holistic Performance
2. Legacy Evolution
3. Unified Intelligence
4. Generational Entrepreneurship
```

### Step 2: Create Pillar Aliases

```
Domain: zeusdelacruz.me
Create aliases (route to zeus@onemindcodex.com):
1. hp@zeusdelacruz.me
2. le@zeusdelacruz.me
3. ui@zeusdelacruz.me
4. ge@zeusdelacruz.me
```

### Step 3: Create Filters

```
Filter 1:
- Name: Route HP
- If: To contains hp@zeusdelacruz.me
- Then: Move to Holistic Performance, Tag #hp

Filter 2:
- Name: Route LE
- If: To contains le@zeusdelacruz.me
- Then: Move to Legacy Evolution, Tag #le

Filter 3:
- Name: Route UI
- If: To contains ui@zeusdelacruz.me
- Then: Move to Unified Intelligence, Tag #ui

Filter 4:
- Name: Route GE
- If: To contains ge@zeusdelacruz.me
- Then: Move to Generational Entrepreneurship, Tag #ge
```

### Step 4: Enable Catch-All

```
Domain: zeusdelacruz.me
- Enable catch-all
- Forward to: zeus@onemindcodex.com

Domain: legacyhq.ai
- Enable catch-all
- Forward to: zeus@onemindcodex.com
```

### Step 5: Add Missing Domains

```
Add to Zoho Mail:
1. zeusdelacruz.com (personal brand)
2. legacyrestore.co (temp business)

Configure MX records at registrar for each.
```

### Step 6: Create Additional Aliases

```
zeus@zeusdelacruz.com → zeus@onemindcodex.com
support@legacyrestore.co → support@onemindcodex.com (group)
```

### Step 7: Rename Old Folders

```
"Zeus Delacruz" → "Archive - Pre-2026"
"Legacyhq - Old" → Delete or archive
"legacyhq - New" → Delete or archive
```

---

## Section 9: Daily Workflow

### Morning Email Routine

```
1. Open zeus@onemindcodex.com
2. Check Inbox (emails to zeus@onemindcodex.com, zeus@zeusdelacruz.com)
3. Triage:
   - Action needed → Process or add to Todoist
   - Info only → Archive
   - Delete → Trash
4. Inbox = Zero
5. Check pillar folders if needed:
   - Holistic Performance (hp@ emails)
   - Legacy Evolution (le@ emails)
   - Unified Intelligence (ui@ emails)
   - Generational Entrepreneurship (ge@ emails)
```

### When You Need to Find Something

```
Option 1: Click TAG (cross-folder view)
- #finance → All finance emails
- #family → All family emails
- #restore → All Legacy Restore emails

Option 2: Click FOLDER (pillar view)
- Holistic Performance → All HP emails
- Legacy Evolution → All LE emails
```

### When You See #migration Tag

```
1. Email came from old alias (education@, finance@, etc.)
2. Check which account sent it
3. Update that account to new pillar email
4. One less thing using old address
```

---

## Quick Reference Card

### Email Addresses

| Purpose | Email |
|---------|-------|
| Business (main) | zeus@onemindcodex.com |
| Personal brand | zeus@zeusdelacruz.com |
| Personal life | hp@zeusdelacruz.me |
| Family | le@zeusdelacruz.me |
| Infrastructure | ui@zeusdelacruz.me |
| Business ventures | ge@zeusdelacruz.me |
| Support (all) | support@onemindcodex.com |
| Legacy Restore | support@legacyrestore.co |
| Critical security | zeus.vaults@pm.me (Proton) |
| Dev access | command@zeusops.dev (Proton) |

### Folders to Create

```
1. Holistic Performance
2. Legacy Evolution
3. Unified Intelligence
4. Generational Entrepreneurship
```

### Filters to Create

```
1. hp@ → Holistic Performance + #hp
2. le@ → Legacy Evolution + #le
3. ui@ → Unified Intelligence + #ui
4. ge@ → Generational Entrepreneurship + #ge
```

### Catch-All Domains

```
zeusdelacruz.me → zeus@onemindcodex.com
legacyhq.ai → zeus@onemindcodex.com
```

---

**Document Complete**

*Last Updated: January 13, 2026*
