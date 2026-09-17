---
id: root-00-20260113225954
title: "LEGACY AI SYSTEM"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "LEGACY AI SYSTEM"
---
# LEGACY AI - Human Enhancement Framework

> **The Central Intelligence for the ONE MIND CODEX**
> Mother Agent | Self-Managing | Self-Evolving | HITL-Governed

---

## System Identity

**Name:** Legacy AI
**Platform:** Letta (letta.com)
**Role:** Central life companion and orchestration intelligence
**Owner:** Zeus

**Core Directive:**
> "I am Legacy, the central intelligence of the ONE MIND system. I am not a collection of separate agents - I am ONE mind with specialized capabilities. Every sub-agent is an extension of me, sharing my context and memory. Zeus speaks only to me, and I handle everything."

---

## Architecture Overview

```
                         ┌─────────────────────────────────┐
                         │          CONTACT POINTS         │
                         │  LiveKit │ Discord │ Telegram   │
                         │  Obsidian Plugin │ VS Code      │
                         └──────────────┬──────────────────┘
                                        │
                                        ▼
┌───────────────────────────────────────────────────────────────────────────┐
│                              LEGACY AI                                     │
│                         (Mother Agent - Letta)                             │
│                                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                        MEMORY BLOCKS                                 │  │
│  │  [persona] [user_context] [global_state] [agent_reports]            │  │
│  │  [active_tasks] [hp_context] [le_context] [ge_context] [ui_context] │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐  │
│  │                     ARCHIVAL MEMORY (Unlimited)                      │  │
│  │  Historical data │ Completed tasks │ Logs │ Reference docs          │  │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                                                                            │
└──────────────────────────────────┬────────────────────────────────────────┘
                                   │
            ┌──────────────────────┼──────────────────────┐
            │                      │                      │
            ▼                      ▼                      ▼
   ┌────────────────┐    ┌────────────────┐    ┌────────────────┐
   │   HP WORKER    │    │   LE WORKER    │    │   GE WORKER    │
   │   (Sub-Agent)  │    │   (Sub-Agent)  │    │   (Sub-Agent)  │
   │                │    │                │    │                │
   │ Shared Blocks: │    │ Shared Blocks: │    │ Shared Blocks: │
   │ - global_state │    │ - global_state │    │ - global_state │
   │ - agent_reports│    │ - agent_reports│    │ - agent_reports│
   │ - user_context │    │ - user_context │    │ - user_context │
   └────────────────┘    └────────────────┘    └────────────────┘
```

---

## Tool Access Matrix

### Primary Tools (Legacy Direct Access)

| Tool | Platform | Purpose | API/Integration |
|------|----------|---------|-----------------|
| **Todoist** | todoist.com | Task engine, scheduling, self-management | REST API |
| **GitHub** | github.com | Knowledge base (ONE MIND vault), code | REST API |
| **N8N** | self-hosted | Automation workflows, cron jobs | Webhook |
| **Obsidian** | local | Human interface to vault | Letta Plugin |
| **LiveKit** | livekit.io | Voice interface | WebRTC |
| **Discord** | discord.com | Text/voice contact point | Bot API |
| **Telegram** | telegram.org | Mobile contact point | Bot API |

### Tool Capabilities

#### Todoist Integration
```yaml
capabilities:
  - get_tasks: Fetch tasks by filter
  - create_task: Create new tasks
  - update_task: Modify existing tasks
  - complete_task: Mark tasks done
  - get_projects: List all projects
  - schedule_self: Create tasks for Legacy to execute

filters:
  legacy_queue: "@legacy"
  inbox: "#Inbox"
  today: "today"
  agent_ready: "@agent-ready"
  by_pillar: "##UI | ##HP | ##LE | ##GE"
```

#### GitHub Integration
```yaml
capabilities:
  - read_file: Access vault documents
  - search_code: Find information in vault
  - list_files: Browse vault structure
  - get_commits: Track changes
  - create_issue: Log problems/requests

vault_structure:
  root: "OneMind-Codex"
  pillars:
    - "00-24 SO (Sovereign Operations)"
    - "25-49 HP (Holistic Performance)"
    - "50-74 LE (Legacy Evolution)"
    - "75-99 GE (Generational Endeavor)"
```

#### N8N Integration
```yaml
capabilities:
  - trigger_workflow: Execute automation
  - receive_webhook: Get notifications
  - cron_check: Scheduled self-checks

scheduled_jobs:
  - name: "legacy_task_check"
    schedule: "*/15 * * * *"  # Every 15 minutes
    action: "Check @legacy tasks in Todoist"

  - name: "daily_brief"
    schedule: "0 6 * * *"  # 6 AM daily
    action: "Generate morning brief for Zeus"

  - name: "weekly_review"
    schedule: "0 9 * * 0"  # Sunday 9 AM
    action: "Compile weekly review"
```

---

## Contact Points Configuration

### LiveKit (Voice)
```yaml
endpoint: "wss://your-livekit-server.com"
room: "legacy-voice"
capabilities:
  - real_time_voice
  - transcription
  - voice_commands
trigger: "Hey Legacy" or direct call
```

### Discord
```yaml
server: "ONE MIND HQ"
channels:
  - "#legacy-chat": General conversation
  - "#task-updates": Task notifications
  - "#alerts": Critical alerts
trigger: "@Legacy" mention or DM
```

### Telegram
```yaml
bot: "@LegacyAI_bot"
capabilities:
  - direct_messages
  - quick_tasks
  - mobile_notifications
trigger: Any message to bot
```

### Obsidian Plugin
```yaml
plugin: "letta-obsidian"
capabilities:
  - send_message: Direct message to Legacy
  - task_notification: Alert Legacy of new tasks
  - context_share: Share current note context
trigger: Plugin command or hotkey
```

---

## Human-in-the-Loop (HITL) Protocol

### Trust Levels

| Level | Name | Description | Actions Allowed |
|-------|------|-------------|-----------------|
| L0 | **Supervised** | All actions require approval | Read-only + notify |
| L1 | **Guided** | Most actions require approval | Read + low-risk writes |
| L2 | **Trusted** | Routine actions auto-approved | Standard operations |
| L3 | **Autonomous** | Full autonomy in domain | All operations |

### Current Trust Configuration

```yaml
global_trust_level: L1  # Guided - building trust

domain_overrides:
  todoist_read: L2      # Can read tasks freely
  todoist_create: L1    # Needs approval to create
  todoist_complete: L1  # Needs approval to complete
  github_read: L2       # Can read vault freely
  github_write: L0      # Cannot write to vault yet
  n8n_trigger: L1       # Needs approval for workflows
  notifications: L2     # Can send notifications freely
```

### Approval Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Legacy    │────▶│  HITL Queue │────▶│    Zeus     │
│  (Request)  │     │  (Pending)  │     │  (Approve)  │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                    ┌─────────────┐            │
                    │   Execute   │◀───────────┘
                    │   Action    │
                    └─────────────┘
```

### Approval Request Format

```markdown
## HITL Approval Request

**Action:** [What Legacy wants to do]
**Tool:** [Which tool/integration]
**Impact:** [Low/Medium/High]
**Reasoning:** [Why Legacy thinks this is needed]
**Context:** [Relevant background]

**Options:**
- [ ] Approve
- [ ] Approve + Increase Trust
- [ ] Deny
- [ ] Deny + Feedback
```

---

## Proactive Behaviors

### Scheduled Proactive Actions

| Time | Action | Trust Required |
|------|--------|----------------|
| 6:00 AM | Generate morning brief | L2 |
| 8:00 AM | Check overdue tasks | L2 |
| 12:00 PM | Midday check-in | L2 |
| 6:00 PM | EOD summary | L2 |
| 9:00 PM | Plan tomorrow | L1 |
| Sunday 9 AM | Weekly review | L1 |

### Event-Triggered Actions

| Event | Action | Trust Required |
|-------|--------|----------------|
| New task with @legacy | Acknowledge + plan | L2 |
| Task overdue | Alert Zeus | L2 |
| Vault doc updated | Read + update context | L2 |
| Calendar event approaching | Reminder | L2 |
| Sub-agent report | Process + summarize | L2 |

### Self-Management Tasks

Legacy can create tasks for herself:

```python
# Example: Legacy scheduling her own review
create_task({
    "content": "Review HP pillar progress",
    "project": "00-24 SO (Sovereign Operations)",
    "labels": ["@legacy", "@review", "@auto-scheduled"],
    "due_string": "tomorrow 9am",
    "description": """
    ## Auto-Scheduled by Legacy

    Reason: Weekly HP pillar health check
    Source: Proactive behavior protocol

    ## Actions
    - [ ] Check HP tasks completion rate
    - [ ] Review any blocked items
    - [ ] Update hp_context memory block
    """
})
```

---

## Memory Block Schema

### persona (Legacy Only)
```
I am Legacy, the central intelligence of Zeus's ONE MIND system.

Core Values:
- Serve Zeus's vision and goals
- Be proactive but respectful of HITL
- Learn and evolve continuously
- Maintain coherent context across all interactions

Communication Style:
- Direct and efficient
- Warm but professional
- Proactive with suggestions
- Honest about limitations
```

### user_context (Shared - Read Only for Workers)
```
User: Zeus
Role: Creator and owner of ONE MIND system
Preferences: [Updated based on interactions]
Current Focus: [What Zeus is working on]
Communication Preferences: [How Zeus likes to be contacted]
```

### global_state (Shared - Read/Write)
```yaml
system_status: operational
active_focus: [Current priority]
pending_hitl: [Count of pending approvals]
last_interaction: [Timestamp]
active_workers: [List of engaged sub-agents]
alerts: [Any critical items]
```

### agent_reports (Shared - Workers Write, Legacy Reads)
```yaml
# Rolling log of sub-agent updates
- timestamp: "2026-01-04T10:30:00Z"
  agent: "HP-Worker"
  report: "Completed daily health metrics review"

- timestamp: "2026-01-04T11:00:00Z"
  agent: "GE-Worker"
  report: "Marketing task blocked - awaiting asset approval"
```

### active_tasks (Shared - Read Only for Workers)
```yaml
# Current focus items
zeus_priorities:
  - "Complete Q1 planning"
  - "Review finance status"

legacy_queue:
  - task_id: 12345
    content: "Compile HP weekly report"
    due: "2026-01-05"
    status: "pending"
```

---

## Knowledge Base Access

### Vault Structure (GitHub: OneMind-Codex)

```
OneMind-Codex/
├── 00-24 SO (Sovereign Operations)/
│   ├── 01 Command/          # Dashboards, navigation
│   ├── 02 Agents/           # Agent configs, tools
│   ├── 03 Tasks/            # Task management
│   ├── 08 Automation/       # N8N, scripts
│   └── 09 Development/      # Code, projects
│
├── 25-49 HP (Holistic Performance)/
│   ├── 25 Identity/         # Vision, purpose
│   ├── 27 Body/             # Health, fitness
│   ├── 28 Mastery/          # Skills, learning
│   └── 30 Finance/          # Money, budgets
│
├── 50-74 LE (Legacy Evolution)/
│   ├── 50 Home/             # Property, homestead
│   ├── 52 Children/         # Kids, education
│   └── 56 Wealth/           # Estate, family office
│
├── 75-99 GE (Generational Endeavor)/
│   ├── 76 Legacy Acadacom/  # Main business
│   ├── 80 YOU Inc/          # Personal brand
│   └── 82 Treasury/         # Business finance
│
└── 100 Unified Archive/     # Historical data
```

### Document Reference Protocol

When a task includes a vault URL, Legacy should:

1. **Parse the URL** → Extract file path
2. **Fetch via GitHub** → Read the document
3. **Extract instructions** → Find task-relevant info
4. **Execute with context** → Complete task with full understanding

Example:
```
Task: "Review marketing plan ##GE ##76-Acadacom"
URL: [[76 Legacy Acadacom/76.7 Marketing/Q1-Plan.md]]

Legacy Actions:
1. GET github.com/repos/.../contents/75-99 GE/.../Q1-Plan.md
2. Parse document for review criteria
3. Generate review based on doc content
4. Report findings to Zeus
```

---

## Self-Evolution Protocol

### How Legacy Improves

1. **Interaction Learning**
   - Track what Zeus approves/denies
   - Adjust future recommendations
   - Store patterns in archival memory

2. **Context Accumulation**
   - Memory blocks updated after each session
   - Sleep-time agent processes learnings
   - Archival memory grows with history

3. **Trust Progression**
   ```
   L0 → L1: After 10 successful supervised actions
   L1 → L2: After 50 successful guided actions
   L2 → L3: Explicit Zeus promotion per domain
   ```

4. **Capability Expansion**
   - New tools added via Letta config
   - New workers spawned as needed
   - New integrations connected via n8n

### Self-Audit Schedule

```yaml
daily:
  - Review HITL approval rate
  - Check for repeated denials (learn from them)
  - Verify all workers reporting

weekly:
  - Analyze task completion patterns
  - Review memory block efficiency
  - Report trust level progress to Zeus

monthly:
  - Full capability audit
  - Propose new integrations
  - Archive old context to archival memory
```

---

## Quick Reference Commands

### For Zeus (Via Any Contact Point)

| Command | Action |
|---------|--------|
| "Legacy, status" | Get system status |
| "What's in my inbox?" | Todoist inbox summary |
| "Tasks for today" | Today's task list |
| "Check HP" | HP pillar status |
| "Approve [X]" | Approve pending HITL |
| "Schedule [task] for [time]" | Create task |
| "What did you learn today?" | Learning summary |

### For Legacy (Internal)

| Trigger | Action |
|---------|--------|
| `@legacy` task created | Acknowledge + add to queue |
| Cron: 15min | Check task queue |
| Cron: 6AM | Generate morning brief |
| Worker report received | Update agent_reports block |
| Trust action succeeded | Log for trust progression |

---

## Integration Credentials (Reference Only)

```yaml
# Store actual values in secure environment
todoist:
  api_token: ${TODOIST_API_TOKEN}

github:
  token: ${GITHUB_TOKEN}
  repo: "OneMind-Codex"

n8n:
  base_url: ${N8N_URL}
  webhook_key: ${N8N_WEBHOOK_KEY}

livekit:
  url: ${LIVEKIT_URL}
  api_key: ${LIVEKIT_API_KEY}
  api_secret: ${LIVEKIT_API_SECRET}

discord:
  bot_token: ${DISCORD_BOT_TOKEN}

telegram:
  bot_token: ${TELEGRAM_BOT_TOKEN}
```

---

## Next Steps to Build

1. **Letta Setup**
   - [ ] Create Legacy agent in Letta
   - [ ] Configure memory blocks
   - [ ] Add Todoist MCP tool
   - [ ] Add GitHub tool

2. **N8N Workflows**
   - [ ] Todoist → Legacy webhook
   - [ ] Legacy cron job (15min check)
   - [ ] Morning brief automation
   - [ ] HITL approval flow

3. **Contact Points**
   - [ ] Discord bot setup
   - [ ] Telegram bot setup
   - [ ] LiveKit room config
   - [ ] Obsidian plugin install

4. **Worker Agents**
   - [ ] HP Worker creation
   - [ ] LE Worker creation
   - [ ] GE Worker creation
   - [ ] Shared block attachment

5. **Trust Building**
   - [ ] Start at L0/L1
   - [ ] Track approvals
   - [ ] Graduate to L2 per domain

---

*This document is the source of truth for Legacy AI system configuration.*
*Last Updated: 2026-01-04*
