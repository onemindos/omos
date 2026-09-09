# Cron State

## Policy
- Mother's cron jobs: **ACTIVE** (infra-watchdog + fabric-status-report)
- All servant cron jobs: **PAUSED** — armed but not firing

## Why
All 9 servant agents are defined, profiled, and ready. Their cron schedules are registered but paused so they don't consume inference credits or flood the log until explicitly activated.

## To Activate a Servant
```bash
# Resume specific job
hermes cron resume <job-id>

# Or resume all for one agent
hermes cron list | grep <agent-name> | awk '{print $1}' | xargs hermes cron resume
```

## Active (Mother)
| Job | Schedule | Status |
|-----|----------|--------|
| infra-watchdog | every 30m | scheduled |
| fabric-status-report | every 6h | scheduled |

## Paused (Servants)
| Agent | Jobs | Status |
|-------|------|--------|
| Oracle | 2 | paused |
| Guardian | 2 | paused |
| Grid | 2 | paused |
| Spartan | 2 | paused |
| Haven | 2 | paused |
| Heritage | 2 | paused |
| Empire | 2 | paused |
| Eden | 2 | paused |
| Trinity | 2 | paused |
