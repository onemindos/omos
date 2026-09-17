---
id: so-02-omos-core
title: "omos-core — Primary Lightsail Node"
type: reference
status: active
created: 2026-05-04
modified: 2026-05-08
owner: zeus
tags: [domain/so, infra/vps, infra/lightsail, infra/tailscale]
---

# omos-core — Primary Lightsail Node

> AWS Lightsail · Ubuntu · Primary production node · All 9 OpenClaw agents

## Identity

| Attribute | Value |
|-----------|-------|
| Hostname | `omos-core` |
| Public IP | `44.193.145.121` |
| Tailscale IP | `100.123.126.34` |
| Platform | AWS Lightsail |
| OS | Ubuntu 24.04 |
| Uptime | Long-running (8d+ as of 2026-05-08) |
| Disk | 4% used |
| RAM | 1730 / 31720 MB |

## Access

```bash
# Via Tailscale (preferred)
ssh ubuntu@100.123.126.34

# Via public IP
ssh -i ~/.ssh/lightsail_default.pem ubuntu@44.193.145.121
```

## Core Platform

| Tool | Purpose |
|------|---------|
| Ubuntu 24.04 | Base OS |
| Docker | Container runtime |
| Python daemon | Background process management |
| cloud-init | VPS bootstrap |

## AI & Agents

| Tool | Purpose |
|------|---------|
| OpenClaw | AI agent gateway — all 9 agents |
| Ollama | Offline model runner (qwen, gemma, etc.) |
| Whisper (local) | Local speech-to-text |
| ffmpeg | Media processor |

## Networking & Security

| Tool | Purpose |
|------|---------|
| Tailscale | VPN mesh — zero public ports |
| fail2ban | Intrusion prevention |
| NATS | Message broker — inter-service pub/sub |

## Data & Storage

| Tool | Purpose |
|------|---------|
| rclone + Dropbox | Cloud sync — VPS ↔ Dropbox |
| boto3 | AWS SDK (S3, EC2, etc.) |
| onemind-watch.service | inotifywait → git commit + push on any vault change |
| onemind-pull.timer | Fires onemind-pull.service every 10s |
| onemind-pull.service | `git pull --rebase` on `/home/ubuntu/onemind` |

## Automation & Messaging

| Tool | Purpose |
|------|---------|
| n8n | Visual workflow automation |
| Matrix | Self-hosted messaging |

## Observability

| Tool | Purpose |
|------|---------|
| Datadog | Monitoring, metrics, APM, logs |

## Security Posture

- **Public ports**: None — confirmed by nmap scan 2026-05-08
- All services accessible via Tailscale only (`100.x.x.x`)
- Tailscale connection: active, direct (not relayed)

## Sync Architecture

```text
Mac (Obsidian / Claude Code)
    ↓ git push
GitHub (OneMind-OS/onemind-codex) ← SOURCE OF TRUTH
    ↓ onemind-pull.timer (every 10s)
    ↑ onemind-watch.service (inotifywait → auto-commit + push)
/home/ubuntu/onemind/ (VPS vault root)
```

| File | Path on VPS |
| ---- | ----------- |
| Pull service | `/etc/systemd/system/onemind-pull.service` |
| Pull timer | `/etc/systemd/system/onemind-pull.timer` |
| Watch service | `/etc/systemd/system/onemind-watch.service` |
| Watch script | `/usr/local/bin/onemind-watch.sh` |
| Vault root | `/home/ubuntu/onemind/` |
| Deploy key | `~/.ssh/codex_deploy` (GitHub deploy key ID 148743097) |

Full runbook: [[SYNC-SERVICES]]

## Service Management

```bash
# Sync services
systemctl status onemind-pull.timer onemind-watch.service
journalctl -u onemind-watch.service -f
systemctl start onemind-pull.service   # force immediate pull

# OpenClaw
systemctl status openclaw
journalctl -u openclaw -f
```

---

## Related

- [[TAILSCALE-NETWORK]] — Full tailnet map
- [[TAILSCALE-SSH]] — SSH runbook
- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
