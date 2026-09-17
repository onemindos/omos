# Healthchecks.io

> **Cron Job & Scheduled Task Monitor**

## Overview

Healthchecks.io is a monitoring service for cron jobs, scheduled automations, and background tasks. It works by expecting a "ping" from your script at regular intervals — if a ping is missed, it alerts you. Think of it as a dead man's switch for your automation stack.

## Status

- **Active**: ⏸ Pending setup
- **Plan**: Free (up to 20 checks) or Hosted
- **Last Updated**: 2026-05-01

## Use Cases

1. **VPS Cron Monitoring** — Verify onemind-pull, backup scripts, and n8n jobs are running
2. **Agent Health Checks** — Alert if OpenClaw agents stop pinging on schedule
3. **Todoist Sync Checks** — Confirm nightly sync scripts completed
4. **Backup Verification** — Confirm Dropbox/USB backup jobs ran

## Key Features

- Simple HTTP ping integration: `curl https://hc-ping.com/{uuid}`
- Email, Slack, Telegram, Discord alerts on missing pings
- Cron syntax scheduling support
- Self-hostable (open source)
- REST API for creating/managing checks

## Integration Notes

- Add ping at end of any cron script: `curl -fsS --retry 3 https://hc-ping.com/{uuid}`
- VPS crontab entries should all have corresponding checks
- Telegram alert → same bot used by OpenClaw for notifications

## Links

- App: https://healthchecks.io
- GitHub: https://github.com/healthchecks/healthchecks
- Docs: https://healthchecks.io/docs/

## Related

- `_system/infrastructure/` — VPS cron jobs to monitor
- [[N8N]] — scheduled workflows to monitor
- `09 Automation` — automation health tracking
