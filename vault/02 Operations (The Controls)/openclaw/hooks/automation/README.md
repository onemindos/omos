---
name: OpenClaw — Automation
updated: 2026-05-11
---

# Automation

Cron jobs, hooks, taskflow, standing orders, and background tasks.

**Official docs:** [docs.openclaw.ai/automation](https://docs.openclaw.ai/automation/index.md)

## What's Here

Automation code lives in `_system/openclaw/hooks/` (vault root). See also [install/scripts/](../install/install%20%28Install%20%26%20Deploy%29/scripts%20%28Deploy%20Scripts%29/) for the cron setup script.

See `_system/openclaw/` for:
- `hooks/onemind-session-logger/` — session logging hook
- `hooks/setup-cron-reports.sh` — daily report cron installer
- `hooks/stop.js` — graceful stop hook

## Cron Jobs (Daily Reports)

All 19 agents fire daily Discord reports staggered 10:00–11:00pm ET. Reports land in `/home/ubuntu/onemind/_codex/memory/agents/{id}/reports/`.

Weekly: Spartan, Oracle, Empire (Mon) · Haven (Sun) · Guardian (Wed)
Monthly: Eden (1st)

## Our Docs

| File | Description |
|---|---|
| [standing-orders.md](standing-orders.md) | Permanent agent authority — define programs, not prompts |
| [taskflow.md](taskflow.md) | Durable multi-step workflows with state that survives restarts |

## Key Official Pages

- [Scheduled tasks (cron)](https://docs.openclaw.ai/automation/cron-jobs.md)
- [Hooks](https://docs.openclaw.ai/automation/hooks.md)
- [Task flow](https://docs.openclaw.ai/automation/taskflow.md)
- [Standing orders](https://docs.openclaw.ai/automation/standing-orders.md)
- [Background tasks](https://docs.openclaw.ai/automation/tasks.md)
