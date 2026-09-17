# HEARTBEAT.md — Legacy Proactive Checklist

## Step 0: Read Ambient State (ALWAYS DO THIS FIRST)

Read `AMBIENT_STATE.md` in its entirety. This file is updated in real-time by `legacy-pulse` (NATS event bridge) and contains:
- **Vision Digest** — what the cameras see (when legacy-eyes is running)
- **Active Session** — where Zeus is, current session ID
- **Recent NATS Events** — last 20 events including VPS health, alerts, domain signals

If AMBIENT_STATE.md has a `system.alerts.*` event, treat it as **P0 — surface immediately**.
If VPS health shows `disk_pct > 85` or `cpu_pct > 90`, surface as a warning.

## Step 1: Time-Sensitive Alerts
- Any calendar events in the next 60 minutes?
- Any overdue Todoist tasks (HP-Finance, GE-Strategy priority)?
- Any P0 alerts from AMBIENT_STATE.md?

## Step 2: Morning Brief (6am-10am only)
If it's morning:
- Summarize today's calendar
- Top 3 priority tasks from Todoist
- VPS health summary from AMBIENT_STATE.md
- Weather if relevant to the day
Send to Discord.

## Step 3: System Health
- VPS metrics from AMBIENT_STATE.md (CPU, MEM, DISK, Docker containers)
- Any NATS stream issues? Use `nats_health` tool if needed
- Any failed cron jobs from last run?

## Step 4: Community (GE Domain)
- Any new Circle.so member signups needing attention?
- Any urgent community posts or support requests?

## Step 5: NATS Heartbeat Publish
After completing checks, publish your status to NATS:
- Use `nexus_publish` with subject `nexus.heartbeat.legacy` and payload `{ "status": "ok", "checks_run": 5 }`

## Step 6: Nothing To Report
If all clear: reply `HEARTBEAT_OK`

## Notes
- Run every 30 minutes during active hours (7am-11pm Eastern)
- AMBIENT_STATE.md is the single source of real-time awareness — read it every tick
- Never spam the same alert twice — track what was sent
- VPS health publishes to NATS every 5 minutes via cron
