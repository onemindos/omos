---
id: ui-00-20260111053925
title: "HIVE Protocol"
type: procedure
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/00]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "HIVE Protocol"
---

# HIVE Protocol

> **The Capture-to-Action System for One Mind**

---

## What is HIVE?

HIVE is the universal capture and routing protocol for the One Mind system. Every piece of information that enters the system flows through HIVE before reaching its destination.

```
H - HARVEST    → Capture raw input
I - INTEGRATE  → Route to correct pillar, merge with existing knowledge
V - VALIDATE   → Confirm placement, add metadata
E - EXECUTE    → Take action
```

---

## The Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   HARVEST   │ →  │  INTEGRATE  │ →  │  VALIDATE   │ →  │   EXECUTE   │
│             │    │             │    │             │    │             │
│ Voice note  │    │ AI routes   │    │ Confirm     │    │ Create task │
│ Text input  │    │ to pillar   │    │ metadata    │    │ File note   │
│ Email fwd   │    │ Match tags  │    │ Check dupes │    │ Trigger n8n │
│ Quick idea  │    │ Find links  │    │ Add context │    │ Notify user │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

---

## HARVEST — Capture Raw Input

**Purpose:** Get information INTO the system as fast as possible with zero friction.

### Input Channels

| Channel | Tool | Landing Zone |
|---------|------|--------------|
| Voice | Whisper/Vapi | Zeus Inbox |
| Text | Telegram/Discord | Zeus Inbox |
| Email | Forward to system | Zeus Inbox |
| Manual | Obsidian quick capture | Zeus Inbox |
| Web | Browser extension | Zeus Inbox |

### Rules

- Capture EVERYTHING — filter later
- Speed over perfection
- One universal inbox (Zeus)
- Timestamp automatically
- Source tracking (where it came from)

---

## INTEGRATE — Route to Correct Pillar

**Purpose:** Merge with existing knowledge and route to the correct location.

### Routing Logic

| If input contains... | Route to |
|----------------------|----------|
| Health, fitness, habits, personal growth | HP (10-19) |
| Family, home, homestead, legacy, heritage | LE (20-29) |
| Business, revenue, products, marketing, clients | GE (30-39) |
| Tech, automation, infrastructure, AI, code | IT (40-49) |
| Cross-domain or system-level | Command (00-09) |

### Integration Actions

1. **Match existing notes** — Find related content via embeddings/search
2. **Suggest links** — Propose wiki-links to existing notes
3. **Detect projects** — If part of active project, link to project note
4. **Identify type** — Task? Reference? Idea? Meeting note?

---

## VALIDATE — Confirm Placement

**Purpose:** Ensure the item is correctly filed with proper metadata.

### Validation Checklist

- [ ] Correct pillar assigned?
- [ ] Proper folder location?
- [ ] Required tags added?
- [ ] Duplicate check passed?
- [ ] Context sufficient?

### Metadata Requirements

| Field | Required | Example |
|-------|----------|---------|
| Pillar | Yes | `HP`, `LE`, `GE`, `IT` |
| Type | Yes | `task`, `note`, `reference`, `idea` |
| Source | Yes | `voice`, `telegram`, `manual` |
| Date | Auto | `2025-12-26` |
| Status | If task | `pending`, `active`, `done` |

### Human Review Triggers

Some items require human confirmation:
- Financial decisions
- External communications
- Ambiguous routing
- High-stakes actions

---

## EXECUTE — Take Action

**Purpose:** Turn validated input into concrete action.

### Action Types

| Type | Action | Destination |
|------|--------|-------------|
| Task | Create in Todoist | Synced to pillar |
| Note | File in Obsidian | Correct X2 folder |
| Reference | Save to knowledge base | X3 folder |
| Idea | Add to incubation | Project inbox |
| Meeting | Create event | Calendar |
| Reminder | Schedule notification | Todoist/Calendar |

### Automation Triggers

Execute can trigger workflows:
- n8n automation chains
- Discord/Telegram notifications
- Email responses
- Calendar events
- Task dependencies

---

## HIVE in Practice

### Example 1: Voice Capture

```
INPUT: "Remind me to call the accountant about Q4 taxes"

HARVEST  → Voice captured via Vapi → Zeus Inbox
INTEGRATE → Keywords: "accountant", "taxes" → Route to GE (Business)
VALIDATE  → Type: task, Pillar: GE, Priority: medium
EXECUTE  → Create Todoist task in GE project, due: tomorrow
```

### Example 2: Email Forward

```
INPUT: Forward email about server maintenance

HARVEST  → Email received → Zeus Inbox
INTEGRATE → Keywords: "server", "maintenance" → Route to IT
VALIDATE  → Type: reference, Pillar: IT, Tags: #infrastructure
EXECUTE  → File in IT/44 Stack & Development/docs/
Evolve 
```

### Example 3: Quick Idea

```
INPUT: "What if we added a family recipe database to the homestead app?"

HARVEST  → Telegram message → Zeus Inbox
INTEGRATE → Keywords: "family", "recipe", "homestead" → Route to LE
           → Also relates to: IT (app development)
VALIDATE  → Type: idea, Pillar: LE (primary), Cross-link: IT
EXECUTE  → Create idea note in LE inbox, link to IT project
```

---

## Implementation

### Where HIVE Lives

| Component | Location |
|-----------|----------|
| Protocol (this doc) | 00 Framework (Unified Governance) |
| Workflow automation | 40-49 IT/42 Workflows/HIVE/ |
| AI processing | 40-49 IT/41 Legacy-AI |
| Landing zone | 01 Zeus/_Inbox.md |

### Tech Stack

- **Capture:** Vapi (voice), Telegram bot, Email forwarding
- **Processing:** Legacy AI (Mother AI routes)
- **Automation:** n8n workflows
- **Storage:** Obsidian vault
- **Tasks:** Todoist (synced)

---

## HIVE vs Legacy AI

| HIVE | Legacy AI |
|------|-----------|
| **Protocol** — How information flows | **System** — The AI agents |
| Capture → Route → Validate → Execute | 5-agent architecture |
| The PROCESS | The INTELLIGENCE |

HIVE is the highway. Legacy AI is the driver.

---

## Quick Reference

```
H - HARVEST    Capture everything, filter nothing
I - INTEGRATE  Route to pillar, find connections
V - VALIDATE   Confirm metadata, check duplicates
E - EXECUTE    Create task, file note, trigger action
```

---

*Part of the One Mind Operating System*
*Last Updated: 2025-12-26*
