---
id: so-02-omos-titan
title: "omos-titan — Titan Lightsail Node"
type: reference
status: active
created: 2026-05-08
modified: 2026-05-08
owner: zeus
tags: [domain/so, infra/vps, infra/lightsail, infra/tailscale]
---

# omos-titan — Titan Lightsail Node

> AWS Lightsail · Ubuntu · Heavy compute node

## Identity

| Attribute | Value |
|-----------|-------|
| Hostname | `omos-titan` |
| Public IP | `44.214.141.64` |
| Tailscale IP | `100.103.197.23` |
| Platform | AWS Lightsail |
| OS | Ubuntu |
| Uptime | Long-running (7d+ as of 2026-05-08) |
| Disk | 11% used |
| RAM | 4236 / 7818 MB |

## Access

```bash
# Via Tailscale (preferred)
ssh ubuntu@100.103.197.23

# Via public IP
ssh -i ~/.ssh/lightsail_default.pem ubuntu@44.214.141.64
```

## Security Posture

- **Public ports**: None — confirmed by nmap scan 2026-05-08
- All services accessible via Tailscale only
- Tailscale connection: active (connected)

## Service Management

```bash
ssh ubuntu@100.103.197.23
systemctl status
docker ps
```

---

## Related

- [[TAILSCALE-NETWORK]] — Full tailnet map
- [[TAILSCALE-SSH]] — SSH runbook
- [[omos-core]] — Primary node
- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
