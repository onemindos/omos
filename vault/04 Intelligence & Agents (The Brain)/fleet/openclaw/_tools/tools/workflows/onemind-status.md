---
name: onemind-status
description: "ONE MIND STATUS — universal system check delivered every morning at 7am ET via the legacy-onemind-status cron. Covers infra, sync, tasks, vault, cron health, and sentinel."
updated: 2026-05-11
---

# Workflow: ONE MIND STATUS

Universal One Mind status check. Deterministic Lobster pipeline that gathers live data from all systems, renders a text board and HTML dashboard, then pauses for approval before delivering to Discord.

**Delivered to:** `#briefings` (Discord `1503277459565445201`)
**Dashboard:** `http://100.123.126.34:8099/status.html` (auto-refreshes every 60s, clickable link posted to Discord)
**Schedule:** `0 7 * * *` America/New_York (daily at 7am ET)
**Agent:** Legacy (`main`)
**Files on VPS:**
- Pipeline: `~/.openclaw/workflows/onemind-status.lobster`
- Data collector: `~/.openclaw/workflows/status-gather.sh`
- Text renderer: `~/.openclaw/workflows/status-render.py`
- HTML renderer: `~/.openclaw/workflows/status-html.py`

---

## What It Checks

| Section | Data Source |
|---|---|
| Infrastructure | `df`, `free`, `uptime`, `ps` on omos-core |
| Sync | `git` status, `onemind-watch.service`, `onemind-pull.timer` |
| Tasks | Todoist API v1 — total, overdue, P1 count, breakdown by pillar |
| Vault | Inbox item count (`99 Inbox (Queue) (The Dock)`) |
| Cron health | `~/.openclaw/cron/jobs.json` + `jobs-state.json` — error count per job |
| Sentinel | SSH health check to `omos-sentinel` — gateway status + memory |

---

## Output Format

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚡ ONE MIND STATUS — Monday, May 11 2026 — 07:00 EDT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## 🖥 INFRASTRUCTURE — omos-core
Disk: 4%  Mem: 1897MB/31720MB (6%)  up 3 days, 9 hours
OpenClaw: 9 processes running
Sentinel (omos-sentinel): gateway ✅  Mem: 963/3834MB

## 🔄 SYNC
HEAD: 485a0ed2  Last: "Obsidian Git: ..." (2 minutes ago)
watch ✅  pull-timer ✅  git: ✅ in sync
Commits (24h): 5

## 🎯 TASKS — 50 active
SO: 6  HP: 1  LE: 1  GE: 0  Other: 42
⚠️ Stalled pillars (0 tasks): GE

## 🗂 VAULT
Inbox: ✅ clear

## ⏰ CRON JOBS — 17/17 healthy  ✅

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ All systems operational
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## Pipeline Steps

```yaml
steps:
  - id: gather    # status-gather.sh → JSON (infra, sync, tasks, vault, crons, sentinel)
  - id: render    # status-render.py → formatted text status board
  - id: html      # status-html.py → writes status.html to ~/.openclaw/canvas/
  - id: approve   # APPROVAL GATE — pauses before posting to Discord
  - id: deliver   # runs only if approved
```

Steps 1–3 are read-only data collection and rendering. The approval gate fires before anything is posted — Legacy shows you the board and asks for sign-off. Approve → delivered to Discord + canvas URL included. Deny → nothing sent.

---

## Lobster Pipeline File

Location: `/home/ubuntu/.openclaw/workflows/onemind-status.lobster`

Two-step pipeline:
1. `gather` — runs `status-gather.sh`, outputs raw JSON
2. `render` — `status-render.py` reads `$gather.stdout`, outputs the formatted status board

---

## Data Collector

Location: `/home/ubuntu/.openclaw/workflows/status-gather.sh`

Python script (executed via bash) that outputs a single JSON object:

```json
{
  "timestamp": "Monday, May 11 2026 — 16:35 EDT",
  "infra": { "node": "omos-core", "disk_used": "4%", "mem_used_mb": 1897, ... },
  "sync": { "head": "485a0ed2", "watch_service": "active", "pull_timer": "active", ... },
  "tasks": { "total": 50, "overdue": 0, "by_pillar": {"SO": 6, "HP": 1, ...} },
  "vault": { "inbox_items": 0 },
  "crons": [ { "name": "legacy-onemind-status", "status": "idle", ... } ],
  "cron_errors": [],
  "sentinel": { "gateway": "active", "mem": "963/3834" }
}
```

---

## Cron Job Prompt (Legacy)

```
Run the ONE MIND STATUS check using Lobster: call the lobster tool with
action=run and pipeline=/home/ubuntu/.openclaw/workflows/onemind-status.lobster
— this collects system status across infra, sync, tasks, vault, cron,
and sentinel. When the pipeline completes, post the rendered output to
#briefings exactly as-is. Reply STATUS_OK when delivered.
```

---

## Troubleshooting

| Issue | Fix |
|---|---|
| `status-gather.sh` fails | Run manually: `bash ~/.openclaw/workflows/status-gather.sh` |
| Todoist error | Token is at `mcp.servers.todoist.env.TODOIST_API_TOKEN` in `openclaw.json` |
| Sentinel timeout | `ssh -o ConnectTimeout=5 omos-sentinel` — check Tailscale |
| Lobster timeout | Increase `timeoutMs` in cron prompt or split gather/render |
| Cron still failing | Check `journalctl --user -u openclaw-gateway -f` for error details |

---

## Related

- [lobster.md](../lobster.md) — Lobster workflow shell reference
- [automation/README.md](../../automation/README.md) — All cron jobs
- [CHANNELS.md](../../channels/CHANNELS.md) — Discord channel IDs
