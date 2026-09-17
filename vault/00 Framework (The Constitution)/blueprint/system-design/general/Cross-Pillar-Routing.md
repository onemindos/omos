---
id: ui-00-20260111053925
title: "Cross-Domain Linking Strategy"
type: concept
domain: SO
area: "00 Framework"
agent: LEGACY
status: active
created: 2026-01-11
modified: 2026-04-30
owner: zeus
tags: [domain/so, type/concept, system/codex, status/active]
related:
  - system-primitives.md
  - Vault-Structure.md
  - OPENCLAW-FRAMEWORK-STRUCTURES.md
todoist_project_id: "6g6P9mmQ3Pf36qjG"
ai_summary: "How to link content across SO/HP/LE/GE domains without duplication — syntax, patterns, MOI template, decision matrix, and maintenance governance."
---

# Cross-Domain Linking Strategy

**Version**: 2.0
**Last Updated**: 2026-04-30
**Purpose**: Define how to link content across the four OneMind domains without duplication

---

## Core Philosophy

**One Source of Truth + Strategic References**

- Every piece of content has ONE primary home
- Other domains LINK to that content, never duplicate it
- Links should explain WHY you're crossing domains

---

## Domain Map

```
SO  (00-24) — Sovereign Operations   — The Brain (systems, AI, infra, doctrine)
HP  (25-49) — Holistic Performance   — Self Path (health, mind, skills, finance)
LE  (50-74) — Legacy Evolution       — Legacy Path (family, home, heritage, estate)
GE  (75-99) — Generational Empire  — Wealth Path (business, revenue, ventures)
```

---

## Linking Syntax & Standards

### Obsidian Wiki-Link Format

```markdown
[[Path/To/File|Display Text]]
```

### Pattern 1: Cross-Domain Reference (with context)
```markdown
**Related**: See [[50-74 LE/Wealth/Investment-Strategy|Family Investment Strategy (LE)]] for long-term wealth building.
```

### Pattern 2: Technical Implementation Link
```markdown
**Technical Details**: [[00-24 SO/03 Infrastructure/|Infrastructure Docs (SO)]]
```

### Pattern 3: Business Application Link
```markdown
**Business Application**: [[75-99 GE/Legacy-AI/Products|Legacy AI as Product (GE)]]
```

### Pattern 4: Bidirectional Reference

**File A** (GE — Legacy AI Business):
```markdown
## Technical Implementation
See [[00-24 SO/02 Intelligence/|SO Intelligence Docs]] for system architecture, deployment, and agent docs.
```

**File B** (SO — Intelligence):
```markdown
## Business Context
See [[75-99 GE/Legacy-AI/|Business Strategy (GE)]] for product positioning, market strategy, and sales materials.
```

---

## Common Cross-Domain Scenarios

### Scenario 1: Business That Uses Technology

**Example**: Legacy AI (the product)

| Domain | What Lives Here |
|--------|----------------|
| **GE** (`75-99 GE/`) | Business strategy, marketing, sales, customer management |
| **SO** (`00-24 SO/02 Intelligence/`) | Agent architecture, code, infrastructure |

### Scenario 2: Personal Finance Spanning Domains

| Domain | What Lives Here |
|--------|----------------|
| **GE** (`75-99 GE/Finance/Business-Credit/`) | Business credit building |
| **HP** (`25-49 HP/Finance/Credit/`) | Personal credit repair and management |
| **LE** (`50-74 LE/Family-Finance/Credit-Education/`) | Teaching family about credit |

### Scenario 3: Real Estate

| Domain | What Lives Here |
|--------|----------------|
| **GE** (`75-99 GE/Real-Estate/`) | Flipping, brokerage (business operation) |
| **LE** (`50-74 LE/Wealth/Real-Estate/`) | Investment properties and family estate |

### Scenario 4: Learning & Development

| Domain | What Lives Here |
|--------|----------------|
| **GE** (`75-99 GE/Education/`) | Courses to sell, community content |
| **HP** (`25-49 HP/Intelligence/Learning/`) | Personal development, skill building |
| **LE** (`50-74 LE/Children-Education/`) | Kids' education, family learning |
| **SO** (`00-24 SO/05 Development/`) | Technical certifications, dev skills |

### Scenario 5: Legal Matters

| Domain | What Lives Here |
|--------|----------------|
| **GE** (`75-99 GE/Entities-Legal/`) | Business structures, contracts, compliance |
| **HP** (`25-49 HP/Personal-Legal/`) | Personal ID, licenses |
| **LE** (`50-74 LE/Estate-Legal/`) | Estate planning, trusts, family legal |

---

## Map of Initiatives (MOI) Pattern

For major initiatives spanning multiple domains, create a **Map of Initiatives** doc in the primary domain.

### MOI Template

```markdown
# [Initiative Name] — Map of Initiatives

**Primary Domain**: [SO / HP / LE / GE]
**Status**: [Active / Planning / Complete]
**Owner**: Zeus
**Todoist Project**: [ID]

## Overview
[Brief description of the initiative and its goals]

## Cross-Domain Breakdown

### SO — Sovereign Operations (00-24)
**What**: [Systems/AI/infra aspects]
**Location**: [[link]]

### HP — Holistic Performance (25-49)
**What**: [Personal aspects]
**Location**: [[link]]

### LE — Legacy Evolution (50-74)
**What**: [Family/legacy aspects]
**Location**: [[link]]

### GE — Generational Empire (75-99)
**What**: [Business/revenue aspects]
**Location**: [[link]]

## Integration Points
[How domains connect and depend on each other]

## Timeline & Milestones
[Key dates and checkpoints]
```

### Example: Fund Launch MOI

```markdown
# Private Equity Fund Launch — Map of Initiatives

**Primary Domain**: GE | **Status**: Active

### SO (00-24)
Investor portal tech, automation, deal flow integrations
→ [[00-24 SO/05 Development/Fund-Tech-Stack|Fund Tech Stack (SO)]]

### HP (25-49)
Personal credibility, licensing, certifications
→ [[25-49 HP/Finance/Investment-Management|Investment Skills (HP)]]

### LE (50-74)
Wealth strategy, family participation, estate alignment
→ [[50-74 LE/Wealth/Investment-Strategy|Family Investment Strategy (LE)]]

### GE (75-99)
Business entity, operations, fundraising, deal flow
→ [[75-99 GE/Fund-Launch/|Fund Launch Operations (GE)]]
```

---

## Tagging for Cross-Domain Discovery

```yaml
---
domain: GE
related:
  - "25-49 HP/Finance/..."
  - "00-24 SO/02 Intelligence/..."
tags: [domain/ge, domain/so, topic/legacy-ai, status/active]
---
```

### Standard Tags by Domain

| Domain | Tags |
|--------|------|
| **SO** (00-24) | `domain/so`, `system`, `infrastructure`, `ai`, `automation`, `doctrine` |
| **HP** (25-49) | `domain/hp`, `personal`, `health`, `finance`, `skills`, `performance` |
| **LE** (50-74) | `domain/le`, `family`, `wealth`, `estate`, `legacy`, `homestead` |
| **GE** (75-99) | `domain/ge`, `business`, `revenue`, `marketing`, `founder`, `entity` |
| **Cross-domain** | `cross-domain`, `initiative`, `integration` |

---

## Callout Conventions

```markdown
> [!info] Related Content in Other Domains
> - **SO**: [[00-24 SO/...|System Context (SO)]]
> - **GE**: [[75-99 GE/...|Business Context (GE)]]

> [!link] Cross-Domain Initiative
> This project spans multiple domains. See the [[Map of Initiatives]] for full overview.

> [!warning] Avoid Duplication
> Don't recreate this content here. The authoritative source is:
> [[Primary Location|Display Name]]
```

---

## Decision Matrix: Link or Duplicate?

| Scenario | Link | Duplicate | Reason |
|----------|------|-----------|--------|
| Reference another domain's strategy | ✅ | ❌ | Single source of truth |
| Quick inline info (1-2 sentences) | ❌ | ✅ | Convenience over reference |
| Detailed process documentation | ✅ | ❌ | Maintain in one place |
| Frequently changing information | ✅ | ❌ | Prevent sync issues |
| Static historical record | ❌ | ✅ | Archive safety |
| Technical specs needed in business doc | ✅ | ❌ | SO owns technical truth |
| Business context needed in technical doc | ✅ | ❌ | GE owns business truth |

---

## Best Practices

**DO:**
- ✅ Link to primary source, never duplicate
- ✅ Explain WHY you're crossing domains
- ✅ Use descriptive link text
- ✅ Create bidirectional links for major connections
- ✅ Use MOIs for complex multi-domain initiatives
- ✅ Update links when moving files

**DON'T:**
- ❌ Copy content across domains "for convenience"
- ❌ Create circular link chains
- ❌ Link without context
- ❌ Use relative links that break on file moves

---

## Maintenance

| Cadence | Action |
|---------|--------|
| **Weekly** | Check for broken cross-domain links, update recent references |
| **Monthly** | Review MOIs for accuracy, update domain hub READMEs |
| **Quarterly** | Audit linking patterns, consolidate duplicated content |

---

## See Also

- [[system-primitives.md]] — Canonical vocabulary for all 55 system concepts
- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/blueprint/system-design/Vault-Structure]] — Current domain/folder layout
- [[OPENCLAW-FRAMEWORK-STRUCTURES.md]] — Agent roster and domain assignments

---

*v2.0 — Updated 2026-04-30 | Domain model updated from v1 pillars to SO/HP/LE/GE (00-99)*
