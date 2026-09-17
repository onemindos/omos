---
title: "omos-cloud"
type: machine
status: active
created: 2026-06-14
updated: 2026-06-14
owner: zeus
subtype: vps
access_level: L2
tags: [type/machine, domain/SO, status/active, vps, aws, lightsail]
ai_summary: "AWS Lightsail VPS (Ubuntu 24.04) — primary OneMind cloud node. Runs Langfuse, WorldMonitor, WorldWideView, and Osiris via Docker. SSH via Tailscale: ubuntu@100.111.189.116."
---

# omos-cloud

> Primary OneMind OS cloud node. AWS Lightsail — always-on, Docker host for OSINT dashboards and LLM observability stack.

## Identity

| Field | Value |
|-------|-------|
| Name / hostname | omos-cloud |
| Type | VPS |
| Provider | AWS Lightsail |
| OS | Ubuntu 24.04.4 LTS |
| Storage | 309 GB (23 GB used, 8%) |

## Network

| Field | Value |
|-------|-------|
| Public IP | 3.224.147.141 |
| Tailscale IP | 100.111.189.116 |
| Tailscale hostname | omos-cloud |
| SSH access | `ssh -i ~/.ssh/lightsail_default.pem ubuntu@100.111.189.116` |
| SSH key | `~/.ssh/lightsail_default.pem` (also in 1Password: SO vault) |

## Services Running

| Service | Port | Stack | Health |
|---------|------|-------|--------|
| langfuse | 3000 | Docker (6 containers) | ✅ healthy |
| worldwideview (WWV) | 3001 | Docker (3 containers) | ✅ healthy |
| worldmonitor | 3002 | Docker (4 containers) | ✅ healthy |
| osiris | 3003 (+8080, 4000) | Docker (3 containers) | ✅ healthy |

## Access

| Method | Details |
|--------|---------|
| SSH (Tailscale) | `ssh -i ~/.ssh/lightsail_default.pem ubuntu@100.111.189.116` |
| Docker stacks | `/opt/omos/` (wwv, osiris) · `/home/ubuntu/worldmonitor` |

## Cost

| Field | Value |
|-------|-------|
| Provider | AWS Lightsail |
| Billing account | In 1Password: SO vault |

## Notes
- `ubuntu` user is in the docker group — no sudo needed for Docker
- All stacks in `/opt/omos/` except worldmonitor which lives in `/home/ubuntu/worldmonitor`
- Tailscale mesh node — not exposed directly to internet for SSH

---

## Related
- [[00-24 SO (Sovereign Operations)/01 The Fabric (Node Registry)/01 The Fabric|The Fabric]]
- [[langfuse]] · [[worldmonitor]] · [[worldwideview]] · [[osiris]]
