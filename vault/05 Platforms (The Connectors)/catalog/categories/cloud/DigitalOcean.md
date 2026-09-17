# DigitalOcean

> **Primary cloud VPS provider hosting all One Mind services**

## Overview

DigitalOcean hosts `onemind-main` — the production VPS running the entire One Mind infrastructure. 4 vCPU, 16 GB RAM, 193 GB disk, Ubuntu 24.04. All OpenClaw agents, N8N, PostgreSQL/TimescaleDB/pgvector, and supporting services run here.

## Status

| Attribute   | Value                              |
| ----------- | ---------------------------------- |
| Status      | ✅ Active                           |
| Category    | Cloud                              |
| Type        | VPS / Cloud Platform               |
| Integration | SSH, doctl CLI, REST API           |

## Use Cases in OneMind

1. **VPS Hosting** — `onemind-main` at `104.248.4.212` (Tailscale: `100.77.52.93`)
2. **Agent Runtime** — OpenClaw + all 9 agents run here
3. **Database Host** — PostgreSQL + TimescaleDB + pgvector (OmniDB)
4. **N8N Server** — Automation workflows run on this VPS

## Key Features

- Droplets (VPS) from $4/mo
- Managed Kubernetes (DOKS)
- Managed Databases (Postgres, Redis, MySQL)
- Spaces (S3-compatible object storage)
- App Platform (PaaS)
- Simple, developer-friendly UI

## Server Details

| Attribute | Value            |
| --------- | ---------------- |
| IP        | 104.248.4.212    |
| Tailscale | 100.77.52.93     |
| OS        | Ubuntu 24.04 LTS |
| vCPU      | 4                |
| RAM       | 16 GB            |
| Disk      | 193 GB SSD       |

## Links

- Website: [digitalocean.com](https://digitalocean.com)
- Docs: [docs.digitalocean.com](https://docs.digitalocean.com)
- CLI: `doctl` — [github.com/digitalocean/doctl](https://github.com/digitalocean/doctl)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
