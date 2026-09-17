---
id: sub-agent-dispatch-20260430
title: "Sub-Agent Dispatch System"
type: protocol
status: active
created: 2026-04-30
modified: 2026-04-30
owner: zeus
tags: [type/protocol, domain/so, system/agents, system/skills]
ai_summary: "How main agents delegate specialized work to sub-agents with skills for system maintenance, complex tasks, and parallel execution."
---

# Sub-Agent Dispatch System

> Main agents are generals. Sub-agents are specialists. Generals decide WHAT. Specialists decide HOW.

---

## Overview

Not every task should be handled by the 5 main agents directly. Complex, specialized, or maintenance work gets **delegated to sub-agents** — temporary workers with specific skills that execute and report back.

```text
ZEUS (Commander)
│
├── LEGACY (Meta-Agent / General)
│   ├── sub: CodeReviewer (skill: code-review)
│   ├── sub: Deployer (skill: ci-cd)
│   ├── sub: SystemMonitor (skill: infrastructure)
│   └── sub: DataBackup (skill: backup-restore)
│
├── SPARTAN (Health Agent / General)
│   ├── sub: WorkoutPlanner (skill: exercise-science)
│   ├── sub: NutritionTracker (skill: meal-planning)
│   └── sub: SleepAnalyzer (skill: sleep-optimization)
│
├── ORACLE (Mind/Money Agent / General)
│   ├── sub: Researcher (skill: deep-research)
│   ├── sub: PortfolioWatcher (skill: investment-analysis)
│   └── sub: LearningPathBuilder (skill: curriculum-design)
│
├── EDEN (Family/Home Agent / General)
│   ├── sub: GroceryManager (skill: shopping-lists)
│   ├── sub: HomeAutomation (skill: smart-home)
│   └── sub: EventPlanner (skill: scheduling)
│
└── EMPIRE (Business Agent / General)
    ├── sub: InvoiceProcessor (skill: accounting)
    ├── sub: ContentWriter (skill: copywriting)
    ├── sub: ClientComms (skill: email-drafting)
    └── sub: MarketResearcher (skill: competitive-analysis)
```

---

## When to Delegate vs. Execute Directly

| Criteria | Direct Execution | Delegate to Sub-Agent |
|----------|-----------------|----------------------|
| Complexity | Simple (< 5 min) | Complex (multi-step) |
| Specialization | General knowledge | Domain expertise needed |
| Parallelism | Sequential is fine | Can run alongside other work |
| Frequency | One-off | Recurring / maintenance |
| Risk | Low impact | High impact (needs specialist) |
| Duration | Quick | Long-running (> 15 min) |

---

## Sub-Agent Lifecycle

```text
1. SPAWN
   Main agent identifies task needing specialist
   ↓
2. BRIEF
   Main agent creates dispatch payload:
   - Task description
   - Required skill(s)
   - Context from Codex
   - Success criteria
   - Time budget
   - Reporting requirements
   ↓
3. EXECUTE
   Sub-agent runs with its skill set
   - Loads skill instructions from _codex/skills/
   - Has scoped access (only its domain)
   - Cannot escalate directly to Zeus (must go through parent)
   ↓
4. REPORT
   Sub-agent returns:
   - Result (success/failure/partial)
   - Artifacts created
   - Time spent
   - Learnings for future runs
   ↓
5. TERMINATE
   Main agent reviews result
   - Marks Todoist task accordingly
   - Logs metrics
   - Sub-agent context is released
```

---

## Dispatch Payload Schema

```yaml
dispatch:
  id: "dispatch-20260430-legacy-003"
  parent_agent: "legacy"
  sub_agent_type: "Deployer"
  skill_required: "ci-cd"
  
  task:
    todoist_id: "todoist-8234567890"
    title: "Deploy staging environment"
    description: "Run full CI/CD pipeline for onemind-web staging"
    priority: "p2"
    due: "2026-04-30T16:00:00Z"
  
  context:
    codex_files:
      - "00-24 SO/05 Development/deploy-playbook.md"
      - "00-24 SO/02 Infrastructure/staging-config.yaml"
    environment: "staging"
    last_deploy: "2026-04-28T14:30:00Z"
  
  constraints:
    time_budget_minutes: 30
    can_modify_production: false
    requires_approval: false
    rollback_on_failure: true
  
  success_criteria:
    - "All tests pass"
    - "Staging URL returns 200"
    - "No error logs in first 5 minutes"
  
  reporting:
    on_complete: "update todoist task + log to metrics"
    on_failure: "escalate to parent agent (Legacy)"
    artifacts: "save deploy log to _system/deploys/"
```

---

## System Maintenance Sub-Agents

These run on **recurring schedules** (self-scheduled crons) to keep the system healthy:

| Sub-Agent | Skill | Schedule | Task |
|-----------|-------|----------|------|
| **HealthChecker** | system-health | Every heartbeat | Verify all services responding |
| **BackupRunner** | backup-restore | Daily 03:00 | Backup Codex + configs |
| **LogRotator** | log-management | Daily 04:00 | Archive old heartbeat/event logs |
| **MetricsAggregator** | data-analysis | Daily 22:00 | Generate daily pulse report |
| **CronCleaner** | cron-management | Hourly | Remove completed/stale crons |
| **SyncChecker** | git-sync | Every 30 min | Ensure Codex git is in sync |
| **DependencyWatcher** | dependency-check | Weekly Monday | Check for outdated packages |
| **QuotaMonitor** | api-management | Every heartbeat | Track Todoist API usage |

---

## Skills System Integration

Sub-agents load their capabilities from the skill system:

```text
Skill Lookup Flow:
1. Sub-agent needs skill "ci-cd"
2. Check _codex/skills/ci-cd.md (personal override)
3. If not found → check PAI/Packs/ (upstream community)
4. If not found → check .claude/skills/ or .agents/skills/ (runtime)
5. If still not found → report "skill not available" to parent
```

### Skill File Format

```yaml
# _codex/skills/ci-cd.md (frontmatter)
---
skill_id: ci-cd
name: "CI/CD Deployment"
version: 1.2
agent_types: [Deployer, SystemMonitor]
requires:
  - git access
  - ssh to VPS
  - GitHub Actions API
capabilities:
  - run_pipeline
  - check_status
  - rollback
  - notify_on_complete
constraints:
  - never deploy to production without approval
  - always run tests first
  - max 3 retries before escalating
---
```

---

## Parallel Execution Rules

Sub-agents can run in parallel, but with guardrails:

```yaml
parallel_rules:
  max_concurrent_sub_agents: 5          # Per main agent
  max_total_sub_agents: 15              # Across all agents
  resource_locks:
    - "todoist_api"    # Max 2 concurrent callers
    - "git_push"       # Max 1 (sequential only)
    - "vps_deploy"     # Max 1 (one deploy at a time)
  conflict_resolution: "first-claimed wins, others queue"
```

---

## Error Handling & Escalation Chain

```text
Sub-Agent Error
│
├─► Retry (up to max_retries from skill config)
│   └─► Still failing?
│
├─► Escalate to Parent Agent
│   └─► Parent decides:
│       ├─► Try different approach
│       ├─► Assign different sub-agent
│       └─► Escalate to Zeus
│
└─► Zeus Intervention
    └─► Zeus decides:
        ├─► Manual fix + resume
        ├─► Cancel task
        └─► Redesign approach (new skill needed)
```

---

## Metrics Per Sub-Agent

Each sub-agent execution logs:

```yaml
sub_agent_metrics:
  dispatch_id: "dispatch-20260430-legacy-003"
  sub_agent_type: "Deployer"
  skill_used: "ci-cd"
  started_at: "2026-04-30T15:30:00Z"
  completed_at: "2026-04-30T15:42:00Z"
  duration_minutes: 12
  result: "success"
  retries: 0
  artifacts_created:
    - "_system/deploys/2026-04-30-staging.log"
  cost_estimate: "$0.02"  # LLM tokens used
  effectiveness: 1.0      # 1.0 = perfect, 0 = total failure
```

These aggregate into the Metrics Engine for the **agent effectiveness score**.


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
