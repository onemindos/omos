---
id: so-02-omos-sentinel
title: "omos-sentinel — Sentinel Lightsail Node"
type: reference
status: active
created: 2026-05-08
modified: 2026-05-08
owner: zeus
tags: [domain/so, infra/vps, infra/lightsail, infra/tailscale]
---

# omos-sentinel — Sentinel Lightsail Node

> AWS Lightsail · Ubuntu · Rescue/failover gateway · Stripped-down OpenClaw

## Identity

| Attribute | Value |
|-----------|-------|
| Hostname | `omos-sentinel` |
| Public IP | `3.214.27.142` |
| Tailscale IP | `100.66.235.38` |
| Platform | AWS Lightsail |
| OS | Ubuntu |
| Uptime | Long-running (7d+ as of 2026-05-08) |
| Disk | 22% used |
| RAM | 915 / 3834 MB |

## Purpose

Sentinel is a **lightweight rescue node** — runs a stripped-down OpenClaw gateway (`openclaw-sentinel`) as failover if omos-core goes down. Legacy can be reached via Sentinel when the primary node is unreachable.

## Access

```bash
# Via Tailscale (preferred)
ssh ubuntu@100.66.235.38

# Via public IP
ssh -i ~/.ssh/lightsail_default.pem ubuntu@3.214.27.142
```

## Running Services

| Service | Description |
|---------|-------------|
| `openclaw-sentinel.service` | Sentinel OpenClaw gateway (port 19789, Tailscale-only) |
| `tailscaled.service` | Tailscale VPN |

**No Docker** — intentionally lean. No database, no NATS, no dashboards.

## OpenClaw Config

| Setting | Value |
|---------|-------|
| Gateway port | `19789` (bound to Tailscale only) |
| Profile | `sentinel` |

## Security Posture

- **Public ports**: None — confirmed by nmap scan 2026-05-08
- All services accessible via Tailscale only
- Tailscale connection: active, direct (not relayed)

## Service Management

```bash
ssh ubuntu@100.66.235.38
systemctl status openclaw-sentinel
journalctl -u openclaw-sentinel -f
```

## Sync

`onemind-pull.timer` pulls GitHub every 10s. Sentinel is stateless by design — no Docker volumes to back up.

---

## Related

- [[TAILSCALE-NETWORK]] — Full tailnet map
- [[TAILSCALE-SSH]] — SSH runbook
- [[omos-core]] — Primary node
- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
