---
title: "Caddy"
type: reference
status: active
tags: [platform, infrastructure, networking, reverse-proxy, http]
---
# Caddy

> Lightweight reverse proxy (15MB). Routes `*.omos.lan:8888` traffic to local services by hostname.

## Role in OneMind

Local HTTP reverse proxy — single entry point on `:8888` for all services. Works with Unbound DNS (resolves hostnames) and NetBird (distributes DNS to all mesh peers).

## Architecture

```
NetBird peers → DNS query → Unbound :5335 → *.omos.lan → 192.168.1.182
Browser → hermes.omos.lan:8888 → Caddy → localhost:8787
```

## Active Routes

| Domain | Backend | Service |
|--------|---------|---------|
| hermes.omos.lan | :8787 | Hermes WebUI |
| hermes-dash.omos.lan | :9119 | Hermes Dashboard |
| dashboard.omos.lan | :3010 | Homepage |
| nats.omos.lan | :8222 | NATS Monitoring |
| nats-dash.omos.lan | :4321 | NATS Dashboard |
| uptime.omos.lan | :3001 | Uptime Kuma |
| netbird.omos.lan | :8081 | NetBird Dashboard |
| agent.omos.lan | :3300 | OneMind Dashboard |
| fleet.omos.lan | :3002 | Fleet View |
| legacy.omos.lan | :3000 | Legacy UI |
| netdata.omos.lan | :19999 | Netdata |
| langfuse.omos.lan | :3004 | Langfuse |

## Config

- `/opt/homebrew/etc/caddy/Caddyfile`
- Managed via brew service: `brew services restart caddy`

## Adding a New Service

1. Add route to Caddyfile
2. Add DNS record to Unbound config
3. Restart both: `brew services restart caddy unbound`
