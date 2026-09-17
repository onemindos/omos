---
id: ui-00-20260111053925
title: "Decision Trees"
type: reference
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/00]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Decision Trees"
---

# Decision Trees

> **Quick routing and decision-making flowcharts**

---

## 1. Content Routing Decision Tree

```
New content arrives...
        │
        ▼
Does it involve MONEY?
  ├─ YES → What stage?
  │   ├─ EARNING (business revenue) → GE (75-99)
  │   ├─ MANAGING (personal budget) → HP (25-49)
  │   └─ PRESERVING (estate/trusts) → LE (50-74)
  │
  └─ NO → Who benefits?
      ├─ ME (personal) → HP (25-49)
      ├─ FAMILY (legacy) → LE (50-74)
      ├─ BUSINESS (revenue) → GE (75-99)
      └─ SYSTEMS (tech/AI) → UI (00-24)
```

---

## 2. Agent Routing Decision Tree

```
Request received...
        │
        ▼
Is this about...
├─ Health, fitness, learning, personal growth?
│  └─ YES → HP Agent (25-49)
│
├─ Family, home, children, estate, legacy?
│  └─ YES → LE Agent (50-74)
│
├─ Business, revenue, products, marketing?
│  └─ YES → GE Agent (75-99)
│
├─ Tech, automation, infrastructure, AI?
│  └─ YES → UI Agent (00-24)
│
└─ Cross-domain or unclear?
   └─ UI Agent coordinates (asks clarifying questions)
```

---

## 3. Autonomy Level Decision Tree

```
Agent wants to take action...
        │
        ▼
Is this IRREVERSIBLE?
  ├─ YES (delete, legal, health) → L0 - HUMAN ONLY
  │
  └─ NO → How much impact?
      ├─ HIGH IMPACT ($100+, external comms, significant changes)
      │  └─ L1 - ASK FIRST (HITL)
      │
      ├─ MEDIUM IMPACT (tasks, metrics, logging)
      │  └─ L2 - ACT THEN NOTIFY
      │
      └─ LOW IMPACT (read, search, route)
         └─ L3 - ACT WITHOUT ASKING
```

---

## 4. File Location Decision Tree

```
Where should this file go?
        │
        ▼
What type of content?
├─ PHILOSOPHY/ARCHITECTURE
│  └─ Is it in ONEMIND-CODEX.md?
│      ├─ YES → Leave it there (single source of truth)
│      └─ NO → Link to CODEX, don't duplicate
│
├─ IMPLEMENTATION/TEMPLATES
│  └─ Which domain?
│      ├─ Brain/Tech → 00-24 SO/
│      ├─ Self → 25-49 HP/
│      ├─ Family → 50-74 LE/
│      └─ Business → 75-99 GE/
│
├─ AGENT CONFIGURATION
│  └─ 00-24 SO/02 Agents/
│
├─ PROTOCOLS/SOPS
│  └─ 00-24 SO/04 Protocols/
│
└─ ARCHIVE/HISTORICAL
   └─ XX Archive/YYYY/
```

---

## 5. Priority Decision Tree

```
Task/project prioritization...
        │
        ▼
Is this URGENT (deadline <48h)?
  ├─ YES → Is it IMPORTANT (aligns with goals)?
  │   ├─ YES → P1 (Do NOW)
  │   └─ NO → P2 (Schedule soon or delegate)
  │
  └─ NO → Is it IMPORTANT (aligns with goals)?
      ├─ YES → P2 (Schedule intentionally)
      └─ NO → P3 (Backlog) or DELETE
```

**Eisenhower Matrix:**
```
               URGENT     NOT URGENT
IMPORTANT    |   P1    |    P2     |
             | Do Now  | Schedule  |
-------------|---------|-----------|
NOT          |   P2    |    P3     |
IMPORTANT    |Delegate | Eliminate |
```

---

## 6. Education/Learning Decision Tree

```
New learning content...
        │
        ▼
Who benefits?
├─ ME (personal growth)
│  └─ HP (25-49) - Personal development
│
├─ BUSINESS (builds capability)
│  └─ GE (75-99) - Business training
│
├─ FAMILY (knowledge transfer)
│  └─ LE (50-74) - Family education
│
└─ TECH (system capability)
   └─ UI (00-24) - Technical training
```

---

## 7. Real Estate Decision Tree

```
Real estate content...
        │
        ▼
What's the purpose?
├─ FAMILY HOME (live in)
│  └─ LE (50-74) - Home/sanctuary
│
├─ HOMESTEAD/PROPERTY (legacy land)
│  └─ LE (50-74) - Legacy estate
│
├─ INVESTMENT PROPERTY (wealth building)
│  └─ LE (50-74) - Generational wealth
│
├─ RENTAL BUSINESS (active income)
│  └─ GE (75-99) - Revenue generation
│
└─ AUTOMATION (smart home tech)
   └─ UI (00-24) - Infrastructure
```

---

## 8. Cross-Pillar Tagging Decision Tree

```
Content spans multiple domains...
        │
        ▼
Identify PRIMARY domain
(Where does the PROJECT live?)
        │
        ├─ Use 3 criteria:
        │  1. Primary beneficiary?
        │  2. Primary purpose?
        │  3. Time horizon?
        │
        ▼
Place file in PRIMARY domain
        │
        ▼
Add tags for secondary domains
Examples: #pillar/HP, #pillar/LE, #pillar/GE
        │
        ▼
Link (don't duplicate) from other contexts
```

---

## 9. Memory Sync Decision Tree

```
Should agent run memory sync?
        │
        ▼
Has vault been updated significantly?
  ├─ YES → What changed?
  │   ├─ New projects added → Sync domain context
  │   ├─ Priorities shifted → Sync current context
  │   ├─ Architecture changed → Re-read CODEX
  │   └─ Agent config updated → Re-read agent folder
  │
  └─ NO → Is it Sunday (weekly sync)?
      ├─ YES → Run weekly sync protocol
      └─ NO → Continue with current memory
```

---

## 10. Unified Score Status Decision Tree

```
Check Unified Score (PI + LI + WI) / 3
        │
        ▼
What's the score?
├─ < 60 (RED)
│  └─ Action: SIMPLIFY
│     - Reduce commitments
│     - Remove low-value projects
│     - Focus on essentials
│
├─ 60-79 (AMBER)
│  └─ Action: STABILIZE
│     - Maintain current course
│     - Optimize existing systems
│     - Don't add new initiatives
│
└─ 80+ (GREEN)
   └─ Action: SCALE
      - Add new initiatives
      - Expand capacity
      - Invest in growth
```

---

## Quick Reference: The 4 Domain Questions

```
┌─────────────────────────────────────────────┐
│  Does this improve ME?                      │
│  └─ YES → HP (25-49) Personal Path          │
├─────────────────────────────────────────────┤
│  Does this protect US?                      │
│  └─ YES → LE (50-74) Legacy Path            │
├─────────────────────────────────────────────┤
│  Does this create VALUE?                    │
│  └─ YES → GE (75-99) Wealth Path            │
├─────────────────────────────────────────────┤
│  How do we coordinate?                      │
│  └─ YES → UI (00-24) The Brain              │
└─────────────────────────────────────────────┘
```

---

*Use these trees for fast decision-making. When in doubt, consult [ONEMIND-CODEX.md](/ONEMIND-CODEX.md)*
