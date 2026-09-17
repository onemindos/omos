---
title: "Netdata"
type: reference
status: active
tags: [platform, infrastructure, monitoring, real-time, self-hosted]
---
# Netdata

> Real-time system monitoring. Per-second metrics, zero config, beautiful dashboards. Running locally via brew.

## Local Access
- `http://localhost:19999`

## Key Features
| Feature | Notes |
|---------|-------|
| Real-time | Per-second granularity |
| Zero config | Auto-discovers everything |
| Metrics | CPU, RAM, disk, network, containers, apps |
| Alerts | Built-in health checks + notifications |
| Lightweight | ~1% CPU usage |
| Retention | Local storage (configurable) |
| Docker | Monitors all containers automatically |
| NATS | Plugin available for NATS monitoring |

## Monitored
- System: CPU, RAM, disk I/O, network
- Docker containers (Langfuse, NetBird, etc.)
- NATS server metrics
- Brew services (Caddy, Cloudflared)

## CLI
```bash
brew services start netdata
brew services stop netdata
# Config: /opt/homebrew/etc/netdata/netdata.conf
```
