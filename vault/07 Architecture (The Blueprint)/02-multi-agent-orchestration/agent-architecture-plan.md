---
id: ui-02-20260111053925
title: "OneMind Agent Architecture Plan"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "OneMind Agent Architecture Plan"
---

# OneMind Agent Architecture Plan

## Overview

Legacy AI serves as the Mother Agent, routing tasks to specialized Pillar Conductors who handle domain-specific work.

```
                        ┌─────────────┐
                        │  Legacy AI  │
                        │   (Mother)  │
                        └──────┬──────┘
                               │ delegates
           ┌───────────┬───────┴───────┬───────────┐
           ▼           ▼               ▼           ▼
    ┌──────────┐ ┌──────────┐   ┌──────────┐ ┌──────────┐
    │    HP    │ │    LE    │   │    GE    │ │    IT    │
    │Conductor │ │Conductor │   │Conductor │ │Conductor │
    └──────────┘ └──────────┘   └──────────┘ └──────────┘
         │            │              │            │
    ┌────┴────┐  ┌────┴────┐   ┌────┴────┐  ┌────┴────┐
    │HP Skills│  │LE Skills│   │GE Skills│  │IT Skills│
    │HP Tools │  │LE Tools │   │GE Tools │  │IT Tools │
    └─────────┘  └─────────┘   └─────────┘  └─────────┘
```

---

## Agent IDs

| Agent | ID | Role |
|-------|-----|------|
| Legacy-AI | agent-da585f6f-189f-472a-974c-29b9578f9b55 | Mother/Router |
| holistic-performance-ai | agent-1e1155fa-de36-475f-90c7-59e1761cb063 | HP Conductor |
| legacy-evolution-ai | agent-3ef12371-a11f-48c9-93f3-109dcf87e224 | LE Conductor |
| generational-entrepreneurship-ai | agent-052cf675-bcb6-468e-bbec-79e5be465c7a | GE Conductor |
| innovative-technology-ai | agent-2655ec04-d529-478e-b2b2-56b63b6f5f37 | IT Conductor |

---

## Memory Block Sharing Strategy

| Block | Legacy | HP | LE | GE | IT |
|-------|--------|----|----|----|----|
| onemind_documentation_system | ✓ | ✓ | ✓ | ✓ | ✓ |
| zeus_context | ✓ | ✓ | ✓ | ✓ | ✓ |
| family_context | ✓ | ✓ | ✓ | ✓ | - |
| pillar_routing | ✓ | - | - | - | - |
| hp_skills | - | ✓ | - | - | - |
| le_skills | - | - | ✓ | - | - |
| ge_skills | - | - | - | ✓ | - |
| it_skills | - | - | - | - | ✓ |

---

## Tool Distribution

### Core Tools (All Conductors)
- `todoist_get_tasks` - Get tasks from Todoist
- `todoist_get_task` - Get specific task details
- `todoist_complete_task` - Mark task complete
- `github_read_file` - Read vault documentation
- `github_edit_file` - Edit and commit changes

### Communication Tools (All Agents)
- `send_discord_alert` - Alert to #alerts
- `send_discord_notification` - Notify specific channel
- `telegram_send_alert` - Telegram alerts

### Pillar-Specific Tools (Future)
- HP: Apple Health sync, workout logging
- LE: Journal integration, relationship tracking
- GE: Financial APIs, CRM integration
- IT: Server management, deployment tools

---

## Skill Memory Block Format

Each skill is a memory block with this structure:

```markdown
## Domain
[What this skill covers]

## Tools Available
- tool_name: What it does and when to use it

## Protocols
- Protocol name: Step-by-step process

## Documentation Paths
- Vault paths for related documentation

## Keywords
[Trigger words that indicate this skill should be used]
```

---

## Day-to-Day Flow

### Morning Briefing
```
User → Legacy AI: "What's on my plate today?"
Legacy AI → Queries all 4 Conductors in parallel
Each Conductor → Checks their Todoist project
Legacy AI → Compiles unified daily briefing
```

### Task Assignment
```
User → Legacy AI: "Schedule my workouts and review the budget"
Legacy AI → Routes "workouts" → HP-Conductor
         → Routes "budget" → GE-Conductor
Both execute in parallel
Legacy AI → Compiles response
```

### Task Completion
```
Conductor completes task:
1. Marks task complete in Todoist
2. Updates documentation in vault
3. Sends Discord notification
4. Reports back to Legacy AI
```

---

## Implementation Phases

### Phase 1: HP Pilot - COMPLETE
- [x] Extract expertise from HP sub-agents
- [x] Create hp_skills memory block
- [x] Add core tools to HP-Conductor
- [x] Create pillar_routing block for Legacy
- [x] Delete empty HP sub-agents (7 deleted)
- [x] Test HP workflow

### Phase 2: Remaining Pillars - COMPLETE
- [x] LE-Conductor setup (le_skills block + tools)
- [x] GE-Conductor setup (ge_skills block + tools)
- [x] IT-Conductor setup (it_skills block + tools)
- [x] Delete all empty sub-agents (33 total deleted)
- [x] Delete CMD agents (3 deleted)

### Phase 3: Inter-Agent Communication - COMPLETE
- [x] Create delegation tools for Legacy
- [x] Test delegation to all 4 conductors
- [x] Verify response parsing

---

## Delegation Tools (Legacy AI Only)

| Tool | Target Agent | Use Case |
|------|--------------|----------|
| `delegate_to_hp` | HP-Conductor | Health, fitness, nutrition, sleep, energy, habits |
| `delegate_to_le` | LE-Conductor | Growth, relationships, journaling, spirituality |
| `delegate_to_ge` | GE-Conductor | Business, wealth, investments, sales, marketing |
| `delegate_to_it` | IT-Conductor | Infrastructure, automation, AI, development |

---

## Final Agent Count

**5 Agents Total** (down from 41):
1. Legacy-AI (Mother/Router)
2. HP-Conductor
3. LE-Conductor
4. GE-Conductor
5. IT-Conductor

---

## Memory Block IDs

| Block | ID |
|-------|-----|
| onemind_documentation_system | block-4d5f7d53-5cd8-459e-9647-0107e37354fb |
| pillar_routing | block-89b6ffa5-4d9f-4cc1-aece-ca7b15d14042 |
| hp_skills | block-b7268a61-275b-4c91-ae7f-d0c3d7b185fa |
| le_skills | block-44a0889e-56b5-424c-8b7d-9d4c4c151171 |
| ge_skills | block-532da1bd-6295-47fc-a0cf-ff1074f59b35 |
| it_skills | block-7e429424-800e-4ba4-aa5d-10b452035016 |

---

## Created
2025-12-28

## Updated
2025-12-28 - All phases complete, multi-agent delegation operational


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
