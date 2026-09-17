---
title: "VPS — onemind-main"
type: machine
status: active
created: 2026-04-09
updated: 2026-04-09
owner: zeus
access_level: L0
subtype: vps
tags: [type/machine, domain/SO, status/active]
ai_summary: "Primary VPS on DigitalOcean — runs OpenClaw agents, NATS, onemind-web, n8n, and PostgreSQL. SSH: root@104.248.4.212."
---

# VPS — onemind-main

## Identity

| Field | Value |
|-------|-------|
| Name / hostname | onemind-main |
| Type | vps |
| Make / Model | DigitalOcean Droplet |
| OS | Ubuntu 24.04 |
| CPU | fill in |
| RAM | fill in |
| Storage | fill in |
| Owned / Leased | DigitalOcean |

## Network

| Field | Value |
|-------|-------|
| Public IP | 104.248.4.212 |
| Tailscale IP | 100.77.52.93 |
| Tailscale hostname | onemind-main |
| Local IP | N/A (cloud) |
| SSH access | `ssh root@104.248.4.212` |
| SSH key | In 1Password: SO vault |

## Services Running

| Service | Port | Purpose |
|---------|------|---------|
| openclaw | 3000 | Agent gateway (9 agents) |
| nats | 4222 / 9222 | Messaging backbone |
| onemind-web | 3456 | OneMind UI |
| postgres | 5432 | Database + pgvector |
| n8n | 5678 | Workflow automation |
| project-ux | 18800 | Project UI service |

## Access

| Method | Details |
|--------|---------|
| SSH | `ssh root@104.248.4.212` |
| Tailscale | `ssh root@100.77.52.93` |
| OneMind Web | `http://100.77.52.93:3456` |
| n8n | `http://100.77.52.93:5678` |
| Credentials | In 1Password: SO vault |

## Cost

| Field | Value |
|-------|-------|
| Provider | DigitalOcean |
| Plan | fill in |
| Monthly cost | fill in |
| Billing account | In 1Password: SO vault |

## Notes

- Systemd services: `openclaw`, `openclaw-reload`, `n8n`
- `codex-watch.service` — disabled (detached HEAD bug, needs fix)
- `codex-pull.timer` — 1min git pull from GitHub
- VPS path for codex: `/home/onemind/onemind-os/codex/`
- OpenClaw workspaces: `/home/onemind/onemind-os/openclaw/workspaces/`


---

## Related

- [[00-24 SO (Sovereign Operations)/08 OM Assets (Entity Registry) (The Registry) (Hive Mind)/08 OM Assets|08 OM Assets]]
