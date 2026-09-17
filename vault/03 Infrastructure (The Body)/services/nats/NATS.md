---
title: "NATS — OneMind Event Bus"
type: reference
status: archived
created: 2026-04-15
owner: zeus
tags: [type/reference, domain/so, status/archived, so/01-intelligence]
ai_summary: "NATS JetStream infrastructure that ran on 104.248.4.212 (onemind-main). Archived April 2026 when that VPS was decommissioned. Not migrated to onemindos-cloud. Documented for future re-deployment."
---

# NATS — OneMind Event Bus

> **Status: ARCHIVED** — This infrastructure ran on `104.248.4.212` (onemind-main). Decommissioned April 2026. Not migrated. Documented here for future re-deployment when needed.
>
> **Decision**: NATS was NOT moved to `onemindos-cloud` (`64.227.8.23`). Docker was NOT added. This folder is a clean reference for future use.

---

## What NATS Was

NATS JetStream was the **central event bus** for the OneMind OS infrastructure. It acted as the nervous system — routing IoT signals, system events, agent notifications, and state changes across all services.

```
IoT Devices (MQTT) ──┐
Home Assistant ───────┤
Agent Events ─────────┼──→ NATS JetStream (4222/1883/9222/7422) ──→ Subscribers
System Alerts ────────┤        (onemind-main VPS)
Location Updates ─────┘
```

### Services That Connected to NATS

| Consumer | Purpose | Protocol |
|---------|---------|---------|
| `legacy-pulse.py` | NATS → AMBIENT_STATE.md (Legacy reads this on heartbeat) | NATS native |
| `bridge.py` | NATS → OpenClaw webhook forwarder | NATS native |
| `project-ux` dashboard | Live event display in browser | NATS WebSocket (9222) |
| `openclaw-nats-ingester` | NATS → SQLite events DB | NATS native |
| Home Assistant | IoT sensor data | MQTT (1883) |
| Edge leaf nodes | Home/office/farm/van local NATS | Leaf node (7422) |

---

## Ports

| Port | Protocol | Purpose |
|------|----------|---------|
| `4222` | NATS native | Client connections (services, agents) |
| `1883` | MQTT | IoT devices, Home Assistant, Victron, sensors |
| `7422` | NATS leaf | Edge location leaf nodes (home, farm, office, van) |
| `8222` | HTTP | Monitoring dashboard |
| `9222` | WebSocket | Browser clients, onemind-ui, project-ux dashboard |

---

## NATS Subjects Schema

### Domain Mapping (matches 3+1 model)
```
system.alerts.>      — P0 system alerts (always surface to Legacy)
legacy.session.>     — Session lifecycle (Zeus opens/closes command center)
omcs.hp.>            — HP (Holistic Performance) signals
omcs.ge.>            — GE (Generational Empire) signals
omcs.it.>            — SO (Infrastructure/Tech) signals
omcs.le.>            — LE (Legacy Evolution) signals
```

### Infrastructure
```
it.droplet.>         — DigitalOcean droplet health events
location.update      — Zeus location changes
nexus.heartbeat.>    — Nexus edge node heartbeats
system.health.>      — Service health changes
```

### Physical/IoT
```
physical.ha.>        — Home Assistant events
physical.robot.>     — Robot/device status
physical.heartbeat.> — Physical sensor heartbeats
```

### Agent & Event Pipeline
```
digital.agent.*.result.legacy   — Agent task results routed to Legacy
events.community.member.joined  — Circle.so webhook events
events.task.completed.>         — ClickUp task completions
```

---

## Python Scripts (in this folder)

| File | Purpose |
|------|---------|
| `legacy-pulse.py` | NATS subscriber → writes AMBIENT_STATE.md so Legacy wakes contextually |
| `legacy-session-start.py` | Fire a session start event when Zeus opens command center |
| `bridge.py` | NATS → OpenClaw webhook bridge (routes events to Legacy as hook messages) |

---

## Docker Compose (in this folder)

`docker-compose.yml` — the infra stack that ran NATS + OmniDB + NATS ingester + legacy-bridge.

**Profiles:**
- `default` — NATS + legacy-bridge (lightweight)
- `--profile data` — adds OmniDB + NATS ingester
- `--profile voice` — adds LiveKit + voice-agent
- `--profile full` — everything

**Start commands (reference only):**
```bash
docker-compose up -d                   # Core (NATS + bridge)
docker-compose --profile data up -d    # + OmniDB + ingester
docker-compose --profile full up -d    # Everything
```

---

## NATS Config (`hub.conf`)

See `hub.conf` in this folder. Key settings:
- `server_name: onemind-cloud-hub`
- JetStream: `max_file_store: 2GB`, `max_memory_store: 512MB`, store at `/data`
- Auth: anonymous `onemind` user + `mqtt_device` for IoT
- Leaf node auth: `edge-home`, `edge-office`, `edge-farm` users (passwords via env)

---

## The Full Architecture (What Was Planned)

The docker-compose reveals the **full intended OneMind infrastructure**:

```
┌─────────── onemind-main VPS ────────────────────────┐
│                                                      │
│  NATS JetStream ←── MQTT (IoT/HA) ←── leaf nodes   │
│       ↓                                              │
│  nats-ingester → OmniDB (PostgreSQL + pgvector       │
│       ↓          + TimescaleDB)                      │
│  legacy-bridge → OpenClaw webhooks → Legacy agent    │
│       ↓                                              │
│  n8n (workflows)                                     │
│  LiveKit (WebRTC voice)                              │
│  voice-agent (Deepgram STT + Azure TTS)              │
└──────────────────────────────────────────────────────┘
```

OmniDB was designed as the **unified memory layer**:
- TimescaleDB → time-series event storage, historical queries, pattern detection
- pgvector → agent memory embeddings, codex RAG, semantic search

PostgreSQL + TimescaleDB + pgvector = ONE database with three superpowers.

---

## What Was Actually Running on 104.248.4.212

Only the NATS container was healthy (running 3+ weeks). OmniDB, n8n, and the other services were not active. The three Python services (`legacy-pulse`, `openclaw-nats-ingester`, `project-ux`) ran as systemd services outside Docker.

```
✅ onemind-nats (Docker) — 3 weeks healthy
✅ nats-dashboard (Docker) — port 18300
✅ legacy-pulse.service — NATS → AMBIENT_STATE.md
✅ project-ux.service — ambient dashboard port 18800
✅ openclaw-nats-ingester.service — NATS → events.db
❌ n8n — dead since Apr 12 (unit missing)
❌ OmniDB — not deployed
❌ onemind-web — files at /opt but container not running
```

---

## Future Re-deployment Notes

When NATS is needed again:

1. **Decide on host** — onemindos-cloud works if Docker is added, or a new dedicated droplet
2. **NATS only needs**: `docker-compose up -d nats` + the `hub.conf`
3. **Set leaf passwords**: `EDGE_HOME_PASSWORD`, `EDGE_OFFICE_PASSWORD`, `EDGE_FARM_PASSWORD` in `.env`
4. **Python deps**: `pip install nats-py` for legacy-pulse + bridge; `pip install nats-py httpx` for bridge
5. **OpenClaw webhook token**: Set `OPENCLAW_HOOKS_TOKEN` env var for bridge.py
6. **OmniDB**: separate deployment, uses TimescaleDB image `timescale/timescaledb:latest-pg16`

---

*Archived: 2026-04-15 when 104.248.4.212 (onemind-main) was decommissioned.*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/services (Running Services)/nats (NATS Config)/nats|nats]]
