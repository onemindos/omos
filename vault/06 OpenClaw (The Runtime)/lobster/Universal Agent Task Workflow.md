---
type: system-documentation
domain: SO/02
tags: [openclaw, lobster, workflows, agents, todoist, automation]
created: 2026-05-12
updated: 2026-05-12
status: active
---

# Universal Agent Task Workflow

One Lobster workflow for all agents (Zenith, Unity, Eternal, Sovereignty) — unified task management with config-driven agent mapping.

## Architecture Decision

**✅ One workflow with args** (not separate workflows per agent)

**Why:** Agents differ only in which Todoist project to query. Logic is identical.

```
agent-tasks.lobster (universal)
  ├── task-fetcher.py    → Query Todoist by agent_id
  ├── task-analyzer.py   → Prioritize actions
  ├── task-reporter.py   → Format report
  └── agent-config.json  → Agent/project mapping
```

## One Workflow, One Skill, Each Agent Gets Their Own Tasks

**Key Design Principle:**

```
┌─────────────────────────────────────────────┐
│  agent-tasks (ONE skill, ONE workflow)      │
└─────────────────────────────────────────────┘
                    │
        ┌───────────┼───────────┬───────────┐
        ▼           ▼           ▼           ▼
    ┌────────┐  ┌────────┐  ┌────────┐  ┌────────┐
    │ Zenith │  │ Unity  │  │Eternal │  │Sovrgty │
    │  (HP)  │  │  (LE)  │  │  (GE)  │  │  (SO)  │
    └────────┘  └────────┘  └────────┘  └────────┘
        │           │           │           │
        ▼           ▼           ▼           ▼
    11 projects 9 projects 10 projects 8 projects
```

**How it works:**

1. **Same skill** (`agent-tasks`) loaded by all agents
2. **Same workflow** (`agent-tasks.lobster`) executed by all agents
3. **Different agent_id** passed as argument → queries different projects
4. **Config-driven filtering** → `agent-config.json` maps each agent to their pillar's projects

**Result:**
- ✅ **Zenith** sees only HP tasks (domains 25-49, 11 projects)
- ✅ **Unity** sees only LE tasks (domains 50-74, 9 projects)
- ✅ **Eternal** sees only GE tasks (domains 75-99, 10 projects)
- ✅ **Sovereignty** sees only SO tasks (domains 00-24, 8 projects)
- ✅ **Legacy** sees ALL tasks (domains 00-99, 38 projects)

**Execution pattern:**

```javascript
// Zenith executes:
agent_id = "zenith"  // Gets HP projects (11 total)

// Unity executes:
agent_id = "unity"   // Gets LE projects (9 total)

// Eternal executes:
agent_id = "eternal" // Gets GE projects (10 total)

// Sovereignty executes:
agent_id = "sovereignty" // Gets SO projects (8 total)

// Legacy executes:
agent_id = "legacy"  // Gets ALL projects (38 total)
```

**DRY principle achieved:** One codebase, zero duplication, config-driven agent scoping.

## Usage

### Run for specific agent

```bash
# From OpenClaw tool
lobster(
  action="run",
  pipeline="/home/ubuntu/.openclaw/workflows/agent-tasks.lobster",
  argsJson='{"agent_id":"zenith"}'
)

# From CLI
openclaw invoke --tool lobster --action run \
  --args-json '{
    "pipeline": "/home/ubuntu/.openclaw/workflows/agent-tasks.lobster",
    "argsJson": "{\"agent_id\":\"zenith\"}"
  }'
```

### Agents

**All agents use the same workflow with different `agent_id` parameter:**

| Agent | ID | Pillar | Parent Project ID | Child Projects | Total | Description |
|-------|-----|--------|------------------|----------------|-------|-------------|
| **Legacy** | `legacy` | ALL | 4 parents | 34 | 38 | Master orchestrator - sees all pillars |
| Zenith | `zenith` | HP | 6g6P9w4gg33h49FX | 10 | 11 | Body, mind, money, spirit |
| Unity | `unity` | LE | 6g6PGPpm3XgJrPGR | 8 | 9 | Family, home, estate |
| Eternal | `eternal` | GE | 6g6PGWvG4f9C2Xr5 | 9 | 10 | Empire, brand, ventures |
| Sovereignty | `sovereignty` | SO | 6g6P9mcc7JMjfPrF | 7 | 8 | Infra, security, automation |

**Each agent manages ALL projects in their pillar** (parent + all children).

**Verified with real Todoist project IDs** (updated 2026-05-12).

## Workflow Steps

1. **Fetch** — Query Todoist API for agent's project (reads `agent-config.json`)
2. **Analyze** — Categorize tasks by priority and due date
3. **Report** — Generate human-readable markdown report
4. **Approve** — Human approval gate before actions
5. **Deliver** — Execute approved actions

## Task Priority Logic

```python
# Critical: Overdue P1 tasks
# High: P1 tasks due today
# Medium: Other overdue tasks
# Low: Future tasks

Priority levels:
🔴 Critical — Overdue P1 (needs immediate attention)
🟠 High — P1 due today
🟡 Medium — Other overdue
🟢 Low — Future work
```

## Output Example

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 TASK REPORT — ZENITH (HP)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## Summary
Total: 4 active
🔴 Overdue: 2
⚡ P1: 3
📅 Due today: 2

**Status:** Zenith: 2 overdue tasks need attention ⚠️

## Recommended Actions (3)

1. 🔴 **Morning workout routine**
   Reason: Overdue P1 task
   Link: https://todoist.com/showTask?id=...

2. 🔴 **Review weekly goals**
   Reason: Overdue P1 task
   Link: https://todoist.com/showTask?id=...

3. 🟠 **Schedule health checkup**
   Reason: P1 task due today
   Link: https://todoist.com/showTask?id=...

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  Action required
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Configuration

### Agent Config (`agent-config.json`)

```json
{
  "zenith": {
    "agent_id": "zenith",
    "pillar": "HP",
    "pillar_name": "Holistic Performance",
    "todoist_project_id": "2342567890",
    "todoist_project_name": "HP - Holistic Performance",
    "description": "Body, mind, money, spirit, identity, intelligence",
    "special_checks": ["health", "fitness", "learning"],
    "emoji": "⚔️"
  }
}
```

**To add new agent:** Add entry to config — no code changes needed.

### Environment

Set `TODOIST_API_TOKEN` in `openclaw.json`:

```json
{
  "mcp": {
    "servers": {
      "todoist": {
        "env": {
          "TODOIST_API_TOKEN": "your-token-here"
        }
      }
    }
  }
}
```

## Scheduling

### Daily morning task check (7 AM ET)

```bash
openclaw cron add \
  --name "zenith-daily-tasks" \
  --cron "0 7 * * *" \
  --tz "America/New_York" \
  --session session:zenith \
  --message "Run agent-tasks.lobster for zenith" \
  --announce \
  --channel discord \
  --to "channel:zenith-tasks"
```

### Add to agent Standing Orders

```markdown
## Standing Order: Daily Task Review
**Authority:** Check my Todoist project every morning
**Trigger:** 7 AM ET (cron-enforced)
**Action:** Run agent-tasks.lobster with my agent_id
**Approval gate:** None for reports, required for task execution
**Escalation:** If 3+ overdue P1 tasks, alert Zeus immediately
```

## Testing

```bash
# Test with mock data (no Todoist token required)
bash /home/ubuntu/.openclaw/workflows/test-agent-workflow.sh zenith

# Test individual steps
python3 task-fetcher.py zenith | \
  python3 task-analyzer.py | \
  python3 task-reporter.py
```

## Extending

### Add agent-specific conditional logic

```yaml
# In agent-tasks.lobster
steps:
  - id: health_check
    command: python3 spartan-health-audit.py
    stdin: $analyze.stdout
    condition: $agent_id == "zenith"
  
  - id: revenue_check
    command: python3 empire-revenue-audit.py
    stdin: $analyze.stdout
    condition: $agent_id == "eternal"
```

### Customize analysis rules

Edit `task-analyzer.py`:

```python
# Example: Flag finance tasks for Eternal
if agent_id == "eternal":
    finance_tasks = [t for t in tasks if "finance" in t.get("labels", [])]
    for task in finance_tasks:
        actions.append({
            "priority": "high",
            "action": "review_finances",
            "task_id": task["id"],
            "task_name": task["content"]
        })
```

## Files

| File | Location | Purpose |
|------|----------|---------|
| `agent-tasks.lobster` | `~/.openclaw/workflows/` | Universal workflow |
| `agent-config.json` | `~/.openclaw/workflows/` | Agent/project mapping |
| `task-fetcher.py` | `~/.openclaw/workflows/` | Todoist API query |
| `task-analyzer.py` | `~/.openclaw/workflows/` | Priority logic |
| `task-reporter.py` | `~/.openclaw/workflows/` | Report formatter |
| `test-agent-workflow.sh` | `~/.openclaw/workflows/` | Mock data test |
| `README-agent-tasks.md` | `~/.openclaw/workflows/` | Full documentation |

## Maintenance

### Update Todoist project IDs

```bash
vim ~/.openclaw/workflows/agent-config.json
# No code changes needed — config-driven
```

### Add new agent

1. Add entry to `agent-config.json`
2. Get Todoist project ID from Todoist UI
3. Test: `bash test-agent-workflow.sh <new-agent-id>`
4. Schedule: `openclaw cron add ...`

## Benefits

✅ **One workflow, all agents** — DRY principle  
✅ **Config-driven** — add agents without code changes  
✅ **Resumable** — approval gates survive restarts  
✅ **Schedulable** — integrates with cron  
✅ **Testable** — mock data for development  
✅ **Extensible** — conditional logic for special cases  
✅ **Auditable** — declarative workflow file  

## When to Split

Only create separate workflows if agents need:
- **Fundamentally different step sequences**
- **Different approval gates**
- **Incompatible data formats**

For simple variations (priorities, filters, special checks), use **conditional steps** in the universal workflow.

## Rule of Thumb

**If 80%+ of logic is shared → one workflow with args**

## Related

- [Lobster vs Manual Orchestration](/home/ubuntu/onemind/00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)/lobster/Lobster vs Manual Orchestration.md)
- [Lobster docs](https://docs.openclaw.ai/tools/lobster)
- [Task Flow](https://docs.openclaw.ai/automation/taskflow)
- [Cron Jobs](https://docs.openclaw.ai/automation/cron-jobs)
- Main workflow: `~/.openclaw/workflows/agent-tasks.lobster`
