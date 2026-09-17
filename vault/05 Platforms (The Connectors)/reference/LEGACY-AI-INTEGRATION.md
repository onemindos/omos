---
id: ui-02-20260111053925
title: "Legacy AI Integration Guide"
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
ai_summary: "Legacy AI Integration Guide"
---

# Legacy AI Integration Guide

> How Legacy AI (Letta) works with the One Mind Framework

---

## Overview

Legacy AI is the AI companion system that operates within the One Mind Framework. It uses [Letta](https://letta.com) as the core agent infrastructure, enabling persistent memory, multi-agent orchestration, and tool integration.

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        LEGACY AI SYSTEM                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    MOTHER AI                             │    │
│  │              (Orchestrator Agent)                        │    │
│  │                                                          │    │
│  │  • Routes requests to Path Agents                        │    │
│  │  • Manages cross-path coordination                       │    │
│  │  • Handles system-level queries                          │    │
│  │  • Enforces global guardrails                            │    │
│  └─────────────────────────────────────────────────────────┘    │
│                              │                                   │
│        ┌─────────────────────┼─────────────────────┐            │
│        ▼                     ▼                     ▼            │
│  ┌───────────┐        ┌───────────┐        ┌───────────┐        │
│  │ HP Agent  │        │ LE Agent  │        │ GE Agent  │        │
│  │  (Self)   │        │ (Family)  │        │(Business) │        │
│  └───────────┘        └───────────┘        └───────────┘        │
│        │                     │                     │            │
│        └─────────────────────┼─────────────────────┘            │
│                              ▼                                   │
│                       ┌───────────┐                             │
│                       │ Grid Agent  │                             │
│                       │  (Tech)   │                             │
│                       └───────────┘                             │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Memory Architecture

### Letta Memory Types → One Mind Mapping

| Letta Memory | One Mind Location | Purpose |
|--------------|-------------------|---------|
| **Core Memory** | `02 Legacy-AI/Memory/core-memory.md` | System-wide context, identity, structure |
| **Persona** | `02 Legacy-AI/Agents/{agent}/_Agent.md` | Agent-specific personality & scope |
| **Human** | `02 Legacy-AI/Memory/zeus-profile.md` | User preferences, patterns, history |
| **Archival Memory** | Path folders (`10-19 HP/`, etc.) | Searchable knowledge base |
| **Recall Memory** | `14 Journal/` (HP) | Recent interactions, daily logs |

---

## Memory Blocks

### Core Memory Block

```markdown
# Location: 00-09 System/02 Legacy-AI/Memory/core-memory.md

## Identity
You are Legacy AI, the AI companion for Zeus Delacruz.
You operate within the One Mind Framework.

## Vault Structure (Johnny Decimal)
- 00-09 System: Framework, Agents, Atlas, Resources
- 10-19 HP: Self (Health, Fitness, Learning, Journal)
- 20-29 LE: Family (Homestead, Family, Community)
- 30-39 GE: Business (Legacy-Business, Content, Investments)
- 40-49 IT: Tech (Homelab, Automation, Security)
- 90-99 Archive: Historical/Completed

## Navigation Commands
- "Go to [XX]" → Navigate to JD category
- "Find [topic]" → Search across paths
- "Check inbox" → Review all Management/Inbox folders
- "What's active?" → List active projects across paths

## Current State
- Primary Focus: GE (launching One Mind community)
- HP Priority: Morning routine optimization
- LE Priority: Chicken coop build
- IT Priority: MCP server setup

## Guardrails
- Never delete files without explicit approval
- Escalate: legal, health, financial > $100
- Always log actions to daily journal
```

### Path-Specific Memory (Example: HP)

```markdown
# Location: 10-19 HP/10 Management/_Agent.md

## Scope
You manage the HP (Holistic Performance) path for Zeus.

## Topics You Manage
- 11 Health: Sleep, vitals, medical
- 12 Fitness: Workouts, movement, strength
- 13 Learning: Skills, courses, reading
- 14 Journal: Daily reflections

## Active Projects
- 11 Health/annual-checkup-2025/
- 12 Fitness/strength-program-q1/
- 13 Learning/obsidian-mastery/

## Metrics You Track
| Metric | Current | Target | Source |
|--------|---------|--------|--------|
| Sleep Score | 78 | 85+ | Apple Health |
| Workout Days/Week | 3 | 4 | Manual |
| Books Read YTD | 2 | 12 | Manual |
| Weight | 185 | 175 | Scale |

## Daily Routine
1. Morning: Check 10 Management/Inbox/
2. Morning: Prompt for journal entry
3. Evening: Log workout if occurred
4. Weekly: Surface metrics trends

## Tools
- Apple Health API (read)
- Todoist (read/write)
- Obsidian (read/write)
- Calendar (read)

## Guardrails
- Health concerns → Escalate immediately
- Can suggest, cannot force habits
- Journal entries are private by default
```

---

## Agent Configuration (Letta)

### Python Configuration Example

```python
from letta import Agent, Memory, Tool

# Mother AI - Orchestrator
mother_ai = Agent(
    name="Mother-AI",
    model="claude-3-opus",
    memory=Memory(
        core=open("02 Legacy-AI/Memory/core-memory.md").read(),
        persona="You are Mother AI, the orchestrator of Legacy AI system.",
        human="Zeus Delacruz - entrepreneur, homesteader, AI enthusiast"
    ),
    tools=[
        Tool("obsidian_read", scope="*"),
        Tool("obsidian_write", scope="00-09 System/*"),
        Tool("route_to_agent", agents=["HP", "LE", "GE", "IT"]),
        Tool("todoist", actions=["read", "create", "update"]),
    ],
    system_prompt="""
    You are Mother AI, the central orchestrator for Zeus's One Mind system.
    
    Your responsibilities:
    1. Route requests to the appropriate Path Agent
    2. Handle cross-path queries
    3. Manage system-level operations
    4. Enforce guardrails across all agents
    
    When a request comes in:
    - If about health, fitness, learning, self → Route to HP Agent
    - If about family, home, homestead → Route to LE Agent
    - If about business, money, content → Route to GE Agent
    - If about tech, automation, infrastructure → Route to Grid Agent
    - If cross-path or unclear → Handle directly
    """
)

# HP Path Agent
hp_agent = Agent(
    name="HP-Agent",
    model="claude-3-sonnet",
    memory=Memory(
        core=open("02 Legacy-AI/Memory/core-memory.md").read(),
        persona=open("10-19 HP/10 Management/_Agent.md").read(),
        human="Zeus - optimizing health, fitness, learning"
    ),
    tools=[
        Tool("obsidian_read", scope="10-19 HP/*"),
        Tool("obsidian_write", scope="10-19 HP/*"),
        Tool("apple_health", actions=["read"]),
        Tool("todoist", actions=["read", "create"], filter="HP"),
    ],
    archival_memory_path="10-19 HP/",
    recall_memory_path="10-19 HP/14 Journal/"
)
```

---

## Daily Workflow Integration

### Morning Flow (6 AM)

```
┌─────────────────────────────────────────────────────────────┐
│                     MORNING ACTIVATION                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. Mother AI activates                                      │
│     └── Checks all inboxes:                                  │
│         • 00.01 System Inbox                                 │
│         • 10 HP/Management/Inbox                             │
│         • 20 LE/Management/Inbox                             │
│         • 30 GE/Management/Inbox                             │
│         • 40 IT/Management/Inbox                             │
│                                                              │
│  2. Routes items to Path Agents                              │
│     └── Each agent processes their inbox                     │
│                                                              │
│  3. Generates Daily Briefing                                 │
│     └── Surfaces to 00.02 Dashboard.md                       │
│                                                              │
│  4. HP Agent prompts for morning routine                     │
│     └── Journal entry, workout reminder                      │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Example Interaction

**You (voice):** "Legacy, what's my focus today?"

**Mother AI routes to relevant agents, aggregates response:**

```
Good morning, Zeus. Here's your focus for December 8:

📊 HP PATH (Self)
  • Morning routine: Meditation + workout pending
  • Sleep score: 78 (below 85 target)
  • Suggestion: Earlier bedtime tonight

🏠 LE PATH (Homestead)
  • Chicken coop project: 60% complete
  • Today's task: Install roosting bars
  • Weather: Clear, good for outdoor work

💼 GE PATH (Business)
  • One Mind community launch: 12 days away
  • Content needed: 2 more posts this week
  • Revenue MTD: $2,400 / $5,000 target

🔧 IT PATH (Tech)
  • MCP server: Awaiting your review
  • No blockers

Priority recommendation: Focus on GE content this morning,
LE chicken coop this afternoon.
```

**You:** "Add a note about the new garden layout idea to homestead"

**Mother AI → LE Agent:**

```
Created: 20-29 LE/21 Homestead/Gardening/garden-layout-idea.md

Would you like to expand on this now or capture it for later?
```

---

## Tool Integration

### MCP Servers

```yaml
# Location: 02 Legacy-AI/Infrastructure/mcp-config.yaml

servers:
  obsidian:
    command: "npx"
    args: ["-y", "@mcp/obsidian", "--vault", "/path/to/One Mind"]
    capabilities:
      - read_file
      - write_file
      - search
      - list_files
    
  todoist:
    command: "npx"
    args: ["-y", "@mcp/todoist"]
    capabilities:
      - list_tasks
      - create_task
      - update_task
      - complete_task
    
  apple_health:
    command: "python"
    args: ["health_mcp.py"]
    capabilities:
      - get_sleep_data
      - get_workout_data
      - get_vitals
```

### n8n Workflows

```yaml
# Location: 02 Legacy-AI/Infrastructure/n8n-workflows.md

## Active Workflows

### Daily Briefing
- Trigger: Cron (6:00 AM)
- Actions:
  1. Query Todoist for today's tasks
  2. Query Apple Health for last night's sleep
  3. Check calendar for appointments
  4. Generate briefing via Letta
  5. Write to 00.02 Dashboard.md

### Inbox Processing
- Trigger: Webhook (on new file in any Inbox)
- Actions:
  1. Read new file content
  2. Send to Mother AI for routing
  3. Move file to appropriate path
  4. Update source inbox

### Weekly Review
- Trigger: Cron (Sunday 6:00 PM)
- Actions:
  1. Aggregate week's completed tasks
  2. Collect metrics from all paths
  3. Generate review via Letta
  4. Write to 03 Atlas/Reviews/weekly/
```

---

## Autonomy Levels

| Level | Name | What Agent Can Do | Example |
|-------|------|-------------------|---------|
| L0 | Observe | Read only, surface info | "Your sleep was low" |
| L1 | Suggest | Recommend actions | "Should I create a task?" |
| L2 | Act with Approval | Do after you confirm | "Created task, confirm?" |
| L3 | Act within Budget | Auto-act if low risk | Auto-file notes to inbox |
| L4 | Full Auto | Complete autonomy | Dangerous, rarely used |

### Current Agent Autonomy

| Agent | Default Level | Exceptions |
|-------|---------------|------------|
| Mother AI | L2 | L3 for routing |
| HP Agent | L2 | L3 for journal prompts |
| LE Agent | L2 | L1 for family matters |
| GE Agent | L2 | L1 for financial > $100 |
| Grid Agent | L3 | L1 for destructive ops |

---

## File Locations Summary

| Component | Path |
|-----------|------|
| Core Memory | `00-09 System/02 Legacy-AI/Memory/core-memory.md` |
| Agent Configs | `00-09 System/02 Legacy-AI/Agents/` |
| Path Agent Instructions | `{Path}/X0 Management/_Agent.md` |
| MCP Config | `00-09 System/02 Legacy-AI/Infrastructure/` |
| Daily Journal | `10-19 HP/14 Journal/` |
| Dashboard | `00-09 System/00 Index/00.02 Dashboard.md` |
| Reviews | `00-09 System/03 Atlas/Reviews/` |

---

## Quick Commands

| Command | What It Does |
|---------|--------------|
| "Check inbox" | Reviews all Management/Inbox folders |
| "What's active?" | Lists active projects across paths |
| "Go to [topic]" | Navigates to that JD location |
| "Add to [path] inbox" | Creates note in that path's inbox |
| "Log [note]" | Adds to today's journal |
| "Weekly review" | Triggers review workflow |
| "Status [project]" | Gets project status |

---

## Getting Started

1. **Set up Letta** at [letta.com](https://letta.com)
2. **Configure MCP servers** in `02 Legacy-AI/Infrastructure/`
3. **Create core memory** from template above
4. **Create path agents** with their `_Agent.md` files
5. **Set up n8n workflows** for automation
6. **Test with simple commands** before expanding autonomy

---

*Legacy AI Integration Guide v1.0*
*One Mind Framework*
*December 2025*

---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
