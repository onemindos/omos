---
id: ui-00-20260111053925
title: "Routing Rules"
type: reference
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/00]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Routing Rules"
---

# Routing Rules

> **How to route requests to the correct domain and agent**

---

## The Core Question Filters

Every request gets filtered through these questions:

| Question | If YES → | Domain | Agent |
|----------|----------|--------|-------|
| "Does this improve ME?" | Personal optimization | HP (25-49) | HP Agent |
| "Does this protect US?" | Family/legacy | LE (50-74) | LE Agent |
| "Does this create VALUE?" | Business/revenue | GE (75-99) | GE Agent |
| "How do we coordinate?" | Tech/AI/infrastructure | UI (00-24) | UI Agent |

---

## Detailed Routing Logic

### HP (25-49): Holistic Performance

**Route HERE if request involves:**
- ✅ Health, fitness, nutrition, sleep
- ✅ Learning, education, courses, books
- ✅ Personal development, habits, routines
- ✅ Personal finance, budgeting, personal credit
- ✅ Skills development, certifications, training
- ✅ Mental health, therapy, mindfulness
- ✅ Personal goals, identity, values
- ✅ Recreation, hobbies, travel (solo)
- ✅ Social relationships, friendships
- ✅ Personal sovereignty, privacy

**Examples:**
- "Help me design a workout plan" → HP
- "What books should I read about X?" → HP
- "Track my sleep quality" → HP
- "Budget for personal expenses" → HP

---

### LE (50-74): Legacy Evolution

**Route HERE if request involves:**
- ✅ Family (spouse, children, parents)
- ✅ Home, property, homestead, sanctuary
- ✅ Marriage, partnership, communication
- ✅ Parenting, children's education
- ✅ Elder care, honoring parents
- ✅ Heritage, traditions, family history
- ✅ Family health, wellness, medical
- ✅ Estate planning, trusts, wills
- ✅ Generational wealth preservation
- ✅ Family office, inheritance
- ✅ Community, extended family, neighbors
- ✅ Family travel, experiences

**Examples:**
- "Plan family vacation" → LE
- "Create estate plan" → LE
- "Homestead project planning" → LE
- "Kids' education strategy" → LE

---

### GE (75-99): Generational Entrepreneurship

**Route HERE if request involves:**
- ✅ Business operations, ventures, companies
- ✅ Revenue generation, sales, monetization
- ✅ Products, services, offerings
- ✅ Marketing, branding, content creation
- ✅ Business finance, cash flow, runway
- ✅ Investments, portfolio, assets
- ✅ Business partnerships, clients, deals
- ✅ Business strategy, analytics, metrics
- ✅ Business legal (contracts, entities)
- ✅ Career (as entrepreneurship - YOU Inc)
- ✅ Personal brand (revenue-generating)

**Examples:**
- "Launch new product" → GE
- "Review business metrics" → GE
- "Create marketing campaign" → GE
- "Negotiate client contract" → GE

---

### UI (00-24): Unified Intelligence

**Route HERE if request involves:**
- ✅ Tech infrastructure, servers, networks
- ✅ AI agents, automation, workflows
- ✅ Development, code, repositories
- ✅ Security, backups, credentials
- ✅ Integrations, APIs, webhooks
- ✅ Monitoring, telemetry, logs
- ✅ Cross-domain coordination
- ✅ Framework updates, architecture
- ✅ Agent configuration, memory blocks
- ✅ System health, performance

**Examples:**
- "Deploy new N8N workflow" → UI
- "Configure agent memory" → UI
- "Set up API integration" → UI
- "Check system status" → UI

---

## Edge Cases & Cross-Domain Routing

### Money/Finance Content

**Decision tree:**
```
Is this about money?
├─ EARNING → GE (business revenue)
├─ MANAGING → HP (personal cash flow)
└─ PRESERVING → LE (wealth & legacy)
```

**Examples:**
- "Business tax strategy" → GE (primary)
- "Personal budget" → HP (primary)
- "Family trust" → LE (primary)
- "Investment portfolio" → LE (preservation) or GE (active trading)

---

### Education/Learning Content

**Decision tree:**
```
Who benefits?
├─ ME (personal growth) → HP
├─ BUSINESS (capability) → GE
├─ FAMILY (knowledge transfer) → LE
└─ TECH (system capability) → UI
```

**Examples:**
- "Learn Python" (personal) → HP
- "Learn sales" (for business) → GE
- "Teach kids coding" → LE
- "Learn Docker" (for infrastructure) → UI

---

### Real Estate Content

**Decision tree:**
```
What's the purpose?
├─ FAMILY HOME → LE
├─ INVESTMENT PROPERTY → LE (wealth building)
├─ RENTAL BUSINESS → GE (revenue generation)
└─ PROPERTY AUTOMATION → UI (technical system)
```

---

### Health Content

**Decision tree:**
```
Who's it for?
├─ ZEUS (personal) → HP
├─ FAMILY MEMBERS → LE
└─ BUSINESS (team health) → GE
```

---

### Technology Content

**Decision tree:**
```
What's the purpose?
├─ PERSONAL PRODUCTIVITY → HP
├─ FAMILY COORDINATION → LE
├─ BUSINESS OPERATIONS → GE
└─ INFRASTRUCTURE/AI → UI
```

---

## Multi-Domain Content

When content spans multiple domains:

### Step 1: Identify Primary Domain
**Where does the PROJECT live?**

Use these criteria:
1. **Primary beneficiary** (who benefits most?)
2. **Primary purpose** (what's the main goal?)
3. **Time frame** (short-term or generational?)

### Step 2: Tag Secondary Domains
Add tags: `#pillar/HP`, `#pillar/LE`, `#pillar/GE`, `#pillar/UI`

### Step 3: Link, Don't Duplicate
One source of truth, linked from other contexts

**Example:**
- "Personal brand content creation"
  - **Primary:** GE (generates revenue)
  - **Tag:** `#pillar/HP` (develops skills)
  - **Location:** `75-99 GE/75 Zeus Delacruz Enterprise/`
  - **Link from:** `25-49 HP/28 Mastery/`

---

## Routing by Time Horizon

| Timeframe | Likely Domain |
|-----------|---------------|
| Today/This Week | HP (daily ops) or GE (current business) |
| This Quarter | GE (business goals) or HP (quarterly habits) |
| This Year | Balance across HP/LE/GE |
| 5+ Years | LE (generational) or GE (business scale) |
| 50+ Years | LE (legacy preservation) |

---

## Agent Handoff Protocol

### When to Hand Off

**UI Agent (Mother AI):**
- Receives all requests first
- Routes based on these rules
- Hands off to domain agent
- Coordinates cross-domain

**Domain Agents (HP/LE/GE):**
- Handle domain-specific requests
- Escalate to UI if cross-domain
- Update UI on significant changes

### Handoff Message Template

```
@[TargetAgent], routing this to you:

Request: [original request]
Domain: [HP/LE/GE/UI]
Reason: [why this domain]
Context: [relevant background]
Priority: [P1/P2/P3/P4]

Let me know if you need anything from me.
```

---

## Quick Routing Checklist

**Before routing, ask:**
1. ✅ Is this about Zeus personally? → HP
2. ✅ Is this about family/home/legacy? → LE
3. ✅ Is this about business/revenue? → GE
4. ✅ Is this about tech/AI/infrastructure? → UI
5. ✅ Does it span multiple? → UI coordinates

**If still unsure:**
- Default to UI (coordinator)
- Ask Zeus for clarification
- Use HITL for important decisions

---

## Routing Examples (50 Common Scenarios)

| Request | Route To | Reason |
|---------|----------|--------|
| "Help me lose weight" | HP | Personal health |
| "Plan family vacation" | LE | Family activity |
| "Launch new product" | GE | Business revenue |
| "Set up N8N workflow" | UI | Automation/tech |
| "Daily workout routine" | HP | Personal fitness |
| "Kids' college fund" | LE | Family wealth |
| "Sales pipeline review" | GE | Business operations |
| "Configure agent memory" | UI | AI infrastructure |
| "Read these 5 books" | HP | Personal learning |
| "Homestead planning" | LE | Legacy property |
| "Marketing campaign" | GE | Business marketing |
| "API integration" | UI | Technical integration |
| "Sleep optimization" | HP | Personal health |
| "Marriage counseling" | LE | Family relationship |
| "Client proposal" | GE | Business development |
| "System monitoring" | UI | Infrastructure |
| "Meal prep routine" | HP | Personal nutrition |
| "Estate planning" | LE | Legacy preservation |
| "Revenue forecast" | GE | Business finance |
| "Deploy server" | UI | Infrastructure |
| "Meditation practice" | HP | Personal mindfulness |
| "Family traditions" | LE | Heritage preservation |
| "Product pricing" | GE | Business strategy |
| "Backup configuration" | UI | Data security |
| "Personal budget" | HP | Personal finance |
| "Home renovation" | LE | Family property |
| "Business tax strategy" | GE | Business legal/finance |
| "Agent health check" | UI | AI monitoring |
| "Skill certification" | HP | Personal development |
| "Parent elder care" | LE | Family support |
| "Partnership agreement" | GE | Business legal |
| "Webhook setup" | UI | Technical integration |
| "Habit tracking" | HP | Personal behavior |
| "Family heritage doc" | LE | Legacy documentation |
| "Investor pitch" | GE | Business funding |
| "Security audit" | UI | Infrastructure security |
| "Personal goals review" | HP | Self assessment |
| "Children education plan" | LE | Family planning |
| "Business metrics dashboard" | GE | Business analytics |
| "N8N troubleshooting" | UI | Technical support |
| "Travel planning (solo)" | HP | Personal recreation |
| "Travel planning (family)" | LE | Family activity |
| "Content creation (business)" | GE | Business marketing |
| "Git repository setup" | UI | Development infrastructure |
| "Personal credit repair" | HP | Personal finance |
| "Family wealth trust" | LE | Legacy preservation |
| "Business credit building" | GE | Business finance |
| "System documentation" | UI | Technical documentation |
| "Therapy/counseling" | HP | Personal mental health |
| "Family office setup" | LE | Generational wealth |

---

*See also: [Navigation-Protocol.md](Navigation-Protocol.md) for where to find content after routing*
