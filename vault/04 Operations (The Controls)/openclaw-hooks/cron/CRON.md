---
title: "OpenClaw — Cron Jobs & Heartbeats"
type: reference
status: active
created: 2026-05-10
modified: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO-02, openclaw/cron]
---

# OpenClaw — Scheduled Cron Jobs & Heartbeats

All jobs live in `~/.openclaw/cron/jobs.json` on omos-core.
Timezone: `America/New_York`. Total jobs: **17**.

---

## The Three Job Types

| Type | Purpose | Channel |
|---|---|---|
| **Heartbeat ping** | 2-line alive signal — fleet status at a glance | `#agent-heartbeat` |
| **Daily/weekly report** | Substantive agent session — reads vault, synthesizes, writes report | Agent's own channel |
| **Scheduled brief** | Morning priorities, weekly review, silent infra check | Dedicated channels |

---

## Heartbeat Pings → `#agent-heartbeat`

Lightweight alive signals. Each agent posts a 4-line status block. Not a report — just a pulse.

| Job | Agent | Schedule | Model |
|---|---|---|---|
| `sentinel-heartbeat` | Sentinel 👁️ | Every 6 hours (`0 */6 * * *`) | Haiku 4.5 |
| `legacy-heartbeat` | Legacy ⚡ | Daily 8:00 AM (`0 8 * * *`) | Sonnet 4.6 |
| `zenith-heartbeat` | Zenith 🏔 | Daily 8:01 AM (`1 8 * * *`) | Sonnet 4.6 |
| `unity-heartbeat` | Unity 🌿 | Daily 8:02 AM (`2 8 * * *`) | Sonnet 4.6 |
| `eternal-heartbeat` | Eternal 👑 | Daily 8:03 AM (`3 8 * * *`) | Sonnet 4.6 |
| `sovereignty-heartbeat` | Sovereignty 🌐 | Daily 8:04 AM (`4 8 * * *`) | Sonnet 4.6 |

**Sentinel fires every 6 hours** (Haiku — its entire job is uptime monitoring, cheap and fast).
**All others fire daily at 8 AM ET**, staggered 1 minute apart to avoid simultaneous cold-starts.

Expected output format:
```
⚡ Legacy — ONLINE
Node: omos-core | 2026-05-11 08:00 ET
Model: claude-sonnet-4.6
Status: GREEN
```

---

## Daily Reports → Agent Channels

Full agent sessions. Agent reads vault + session logs, synthesizes pillar state, writes markdown report to vault `_codex/memory/`, posts to its own channel.

| Job | Agent | Schedule | Delivers to |
|---|---|---|---|
| `sovereignty-infra-check` | Sovereignty 🌐 | Daily 3:00 AM (`0 3 * * *`) | **Silent** — only posts to `#system-status` if flag |
| `legacy-morning-brief` | Legacy ⚡ | Daily 7:00 AM (`0 7 * * *`) | `#morning-brief` |
| `sentinel-daily-report` | Sentinel 👁️ | Daily 10:55 PM (`55 22 * * *`) | `#sentinel-chat` |
| `legacy-daily-report` | Legacy ⚡ | Daily 11:00 PM (`0 23 * * *`) | `#legacy-chat` |
| `sovereignty-daily-report` | Sovereignty 🌐 | Daily 11:05 PM (`5 23 * * *`) | `#sovereignty-chat` |
| `zenith-daily-report` | Zenith 🏔 | Daily 11:15 PM (`15 23 * * *`) | `#zenith-chat` |
| `eternal-daily-report` | Eternal 👑 | Daily 11:20 PM (`20 23 * * *`) | `#eternal-chat` |

**3 AM infra check is silent** — Sovereignty runs disk/backup/service checks but posts nothing unless something is wrong. Zero noise when GREEN.

---

## Weekly Reports → `#weekly-reports` + Agent Channels

| Job | Agent | Schedule | Delivers to |
|---|---|---|---|
| `zenith-weekly-scorecard` | Zenith 🏔 | Monday 8:00 AM (`0 8 * * 1`) | `#weekly-reports` |
| `eternal-weekly-scorecard` | Eternal 👑 | Monday 9:00 AM (`0 9 * * 1`) | `#weekly-reports` |
| `unity-weekly-report` | Unity 🌿 | Sunday 10:30 AM (`30 10 * * 0`) | `#unity-chat` |
| `legacy-weekly-review` | Legacy ⚡ | Sunday 9:00 AM (`0 9 * * 0`) | `#weekly-reports` |

---

## Report Output Locations

```
Reports:      /home/ubuntu/onemind/_codex/memory/agents/{id}/reports/YYYY-MM-DD.md
Briefs:       /home/ubuntu/onemind/_codex/memory/agents/main/briefs/YYYY-MM-DD-morning.md
Health logs:  /home/ubuntu/onemind/_codex/memory/agents/sovereignty/health/YYYY-MM-DD-3am.md
```

---

## Agent Models

| Agent | Model | Fallback |
|---|---|---|
| Legacy (`main`) | `github-copilot/claude-sonnet-4.6` | `amazon-bedrock/us.anthropic.claude-sonnet-4-6` |
| Zenith | `github-copilot/claude-sonnet-4.6` | `amazon-bedrock/us.anthropic.claude-sonnet-4-6` |
| Unity | `github-copilot/claude-sonnet-4.6` | `amazon-bedrock/us.anthropic.claude-sonnet-4-6` |
| Eternal | `github-copilot/claude-sonnet-4.6` | `amazon-bedrock/us.anthropic.claude-sonnet-4-6` |
| Sovereignty | `github-copilot/claude-sonnet-4.6` | `amazon-bedrock/us.anthropic.claude-sonnet-4-6` |
| Sentinel | `github-copilot/claude-haiku-4.5` | `amazon-bedrock/us.anthropic.claude-haiku-4-5-20251001-v1:0` |

Sentinel is Haiku by design — watchdog role, fast and lightweight.

---

## Managing Jobs

```bash
# List all jobs
ssh onemindos-core 'openclaw cron list'

# View run history
ssh onemindos-core 'openclaw cron runs'

# Trigger a job manually (useful for testing)
ssh onemindos-core 'openclaw cron trigger <job-id>'

# Enable / disable
ssh onemindos-core 'openclaw cron enable <job-id>'
ssh onemindos-core 'openclaw cron disable <job-id>'

# Edit directly
ssh onemindos-core 'nano ~/.openclaw/cron/jobs.json'
# Then reload:
ssh onemindos-core 'kill -HUP $(pgrep -f openclaw-gateway)'
```

Config: `~/.openclaw/cron/jobs.json`
Run history: `~/.openclaw/cron/runs/`

---

## Related

- [[../channels/CHANNELS|CHANNELS]] — Full channel ID reference
- [[../workspaces (Agent Workspaces)/|Agent Workspaces]] — HEARTBEAT.md per agent
- [[../README|OpenClaw Platform README]] — Fleet overview
