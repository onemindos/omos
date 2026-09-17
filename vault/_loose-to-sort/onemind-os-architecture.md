---
id: so-02-20260527-onemind-os-architecture
title: "OneMind OS — Unified Architecture Decision Record"
type: doctrine
status: active
created: 2026-05-27
modified: 2026-05-27
owner: zeus
pillar: SO
domain: "02 Infrastructure (The Body)"
ai_summary: "Canonical architecture: Hermes (Legacy) as the sole brain, NATS as nervous system, NetBird as network mesh, Flutter as single UI. Eliminates SwarmClaw, redundant dashboards, and standalone workers in favor of Hermes plugins and native NATS tools."
tags: [type/doctrine, zone/so, domain/so, status/active, infra/nats, infra/hermes, framework/core, architecture/canonical]
depth: 0
---

# OneMind OS — Unified Architecture

> **Core Equation:** `HUMAN (Zeus) + AI (Legacy) = ONE MIND`
> **Architecture:** Hermes (brain) + NATS (nervous system) + NetBird (network) + Flutter (interface)
>
> By Zeus Delacruz | May 2026

---

## 1. The Decision

**Hermes IS the operating system.** Not a component in a larger stack — the stack itself. Everything else is either a NATS subject Hermes talks to, a plugin inside Hermes, or a UI that connects to Hermes.

No middleman backends. No redundant orchestrators. No separate worker processes where a plugin suffices.

---

## 2. The Stack (Three Layers + Network)

```
┌─────────────────────────────────────────────────────────────────────────┐
│  LAYER 7 — APPLICATION                                                   │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  NATS Messaging + JetStream + KV + Object Store + Services      │    │
│  │                                                                  │    │
│  │  • Pub/Sub (fire & forget, fan-out, 1-to-many)                  │    │
│  │  • Request/Reply (synchronous agent-to-agent)                   │    │
│  │  • Queue Groups (load balance across instances)                  │    │
│  │  • JetStream (persistence, replay, exactly-once, retention)     │    │
│  │  • KV Store (live state — replaces Redis/etcd/Consul)           │    │
│  │  • Object Store (firmware, configs — internal S3 replacement)   │    │
│  │  • Services ($SRV — discovery, health, stats)                   │    │
│  │  • Consumers (filtered delivery to specific agents/nodes)       │    │
│  │  • Sources & Mirrors (cross-cluster replication, store & forward)│    │
│  │                                                                  │    │
│  │  THIS IS THE ENTIRE MIDDLEWARE. No Redis. No Kafka. No etcd.    │    │
│  │  No Consul. No RabbitMQ. No SQS. One binary.                   │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  LAYER 6 — PRESENTATION / API                                            │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  Hermes REST API (:9119) — sessions, costs, health              │    │
│  │  Hermes WebUI (:8787) — chat interface + WebSocket streaming    │    │
│  │  Hermes Gateway — messaging platforms (Telegram, NATS, etc.)    │    │
│  │  Hermes NATS Platform — Synadia Agent SDK v0.3 (any NATS client)│    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  LAYER 5 — SESSION / STATE                                               │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  Hermes state.db (SQLite + FTS5) — session transcripts, search  │    │
│  │  NATS KV (live state: agents, tasks, config, presence)           │    │
│  │  NATS JetStream (event history, telemetry, audit trail)          │    │
│  │  NATS Object Store (firmware, configs, small models)             │    │
│  │  Langfuse (ClickHouse) — observability, traces, cost tracking   │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  LAYER 4 — TRANSPORT SECURITY                                            │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  NATS TLS (in-transit encryption)                                │    │
│  │  NKeys + Accounts (authentication + namespace isolation)         │    │
│  │  Synadia Security (future: signed events, audit, compliance)     │    │
│  │  NetBird WireGuard (encrypted tunnels for all WAN links)         │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  LAYER 3 — NETWORK                                                       │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  NetBird (WireGuard mesh — NAT traversal, peer routing)          │    │
│  │                                                                  │    │
│  │  Connects: MBP ↔ Cloud ↔ Barn ↔ Office ↔ Drones ↔ Robots      │    │
│  │  NATS needs a network. NetBird IS that network.                  │    │
│  │  They are different layers. Both required.                       │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  LAYER 2 — DATA LINK                                                     │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  WiFi (onemindos) • Ethernet • LTE/5G (drones, mobile robots)   │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
│  LAYER 1 — PHYSICAL                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  MBP M2 Max • Dell server • Lightsail VMs • Unitree G1/Go2     │    │
│  │  DIMOS • Drone RECON • Raspberry Pis • Smart home hardware      │    │
│  └─────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Hermes = The Entire Brain

Hermes is not one component among many. It is the **operating intelligence** that owns:

| Capability | Implementation | Status |
|-----------|---------------|--------|
| **Agent orchestration** | Legacy (primary) + profiles (other agents) | ✅ Built |
| **NATS native tools** | `nats_request`, `nats_publish`, `nats_kv_*`, `nats_stream_read` | ✅ Built |
| **NATS platform** | Synadia Agent SDK v0.3 gateway adapter | ✅ Built |
| **SaaS execution** | NATS worker (Composio + Todoist via `saas.*` subjects) | ✅ Built |
| **Home automation** | HA-NATS bridge (`cmd.ha.*`) | ✅ Built |
| **Task management** | Kanban plugin (board, dispatch, workers) | ✅ Built |
| **Observability** | Langfuse plugin (traces all sessions) | ✅ Built |
| **Memory** | Persistent cross-session memory | ✅ Built |
| **Skills** | Procedural knowledge base | ✅ Built |
| **Scheduling** | Cron system (recurring autonomous jobs) | ✅ Built |
| **Multi-agent** | Profiles (same runtime, different identities) | ✅ Built |
| **Voice** | LiveKit integration (STT/TTS) | ✅ Built |
| **Delegation** | Subagents for parallel work | ✅ Built |
| **ACP** | Agent Communication Protocol | ✅ Built |
| **Chat interfaces** | WebUI, TUI, CLI, API, messaging platforms | ✅ Built |
| **Watch services** | NATS subscribers → trigger agent sessions | 🔨 To build |
| **Infrastructure tools** | Docker/service status as callable tools | 🔨 To build |
| **Langfuse query** | Read traces/costs from Langfuse API | 🔨 To build |
| **Object Store tools** | `nats_obj_get`, `nats_obj_put` | 🔨 To build |
| **$SRV discovery tool** | List all registered NATS services | 🔨 To build |

---

## 4. What Dies

| Service | Reason | Timeline |
|---------|--------|----------|
| **SwarmClaw** | Hermes does everything it did. Flutter talks to Hermes directly. | After Flutter is wired |
| **onemind-dashboard** (React/Bun) | Flutter replaces it | After Flutter fleet page works |
| **Homepage** (:3010) | Flutter home page replaces it | After Flutter is live |
| **Hermes Dashboard UI** (:9119 web) | Flutter consumes the API, no standalone UI needed | After Flutter analytics page |
| **Uptime Kuma** (:3001) | NATS `$SRV.PING` + health-watch service in Hermes | After watch services built |
| **Synadia Agent UI** | Already dead, remove plist | Immediate |

**Rule:** Nothing dies until its replacement is live in Flutter.

---

## 5. What Stays

| Service | Port | Role | Why It Stays |
|---------|------|------|-------------|
| **NATS Server** | 4222, 7422, 8222 | Nervous system | The backbone. Non-negotiable. |
| **Hermes** (Legacy) | — | Brain, orchestrator, API backend | IS the OS |
| **Hermes WebUI** | 8787 | Chat interface (until Flutter chat replaces) | Primary interaction |
| **Hermes Gateway** | — | Messaging platforms | Telegram etc. |
| **HA-NATS Bridge** | — | Home automation translator | Stateless, separate concern |
| **SaaS Action Worker** | — | Composio executor (may fold into plugin) | Could become plugin |
| **Langfuse** | 3004 | Observability across all sessions | Deep traces, costs |
| **NetBird** | 33080 | WAN mesh for remote/edge nodes | Layer 3 — required |
| **Caddy** | 8888 | Reverse proxy | TLS termination, routing |
| **Netdata** | 19999 | System metrics (bridged to NATS by Hermes) | Mature, free, deep |
| **OCIS** | 9200 | User-facing file storage | Evaluate later |
| **OpenClaw** | 3000 | Fleet agent harness (other agents) | Registers as NATS services |

---

## 6. Data Architecture

### Hot State (realtime, watchable)

**NATS KV Buckets:**
- `agents` — live registry, heartbeats, capabilities
- `tasks` — active task board (in-progress, pending)
- `config` — shared fleet configuration
- `presence` — who's online right now

KV is NOT a database. It's a key-value store for live state. No queries, no joins, no search. But: instant reads, realtime watches (push on change), TTL support.

### Warm State (queryable, historical)

**Hermes state.db (SQLite + FTS5):**
- All session transcripts (every provider, every model)
- Full-text search across all conversations
- Costs, tokens, tool calls

**Hermes Kanban DB:**
- Task history (completed, archived)
- Mission logs
- Dispatch records

### Cold State (observability, audit)

**Langfuse (ClickHouse):**
- LLM traces with latency, tokens, costs
- Tool execution spans
- Quality scores

**NATS JetStream Streams:**
- `ACTIVITY` — all events from all agents (durable, replayable)
- `TELEMETRY` — infra metrics (7-day retention)
- `COMMANDS` — command audit trail

### Blob State (files, distribution)

**NATS Object Store:**
- Firmware updates to drones/robots
- Config packages to edge nodes
- Small ML models (<500MB)
- Limitations: no partial reads, no CDN URLs, no versioning, storage = disk

**OCIS / Filesystem:**
- User-facing files
- Obsidian vault (knowledge base)
- Large media

---

## 7. NATS Object Store — Limitations

| Limitation | Detail |
|-----------|--------|
| **Chunk size** | Bounded by NATS `max_payload` (default 1MB, max 64MB) |
| **No partial reads** | Whole object or nothing. No byte-range requests. |
| **No public URLs** | Internal only. Can't generate download links. |
| **No versioning** | Overwrite = gone. No rollback. |
| **Storage = JetStream disk** | Replicates to each node in the replica set. |
| **No per-object ACL** | Access governed per-bucket by NATS accounts. |

**Use for:** firmware, configs, scripts, small models.
**Don't use for:** video, large datasets, user-facing downloads.

---

## 8. Watch Services — Making Hermes Proactive

Currently Hermes is reactive (waits for input). Watch services make it autonomous:

```yaml
# Hermes plugin: nats-watchers
# Subscribe to NATS subjects → trigger Hermes sessions

watchers:
  - subject: "$SYS.ACCOUNT.*.DISCONNECT"
    action: "Service disconnected. Check what went down, attempt remediation."
    
  - subject: "infra.docker.container.die"
    action: "Container crashed. Attempt restart, notify Zeus if critical."
    
  - subject: "events.ha.motion_detected"
    filter: { entity_id: "binary_sensor.front_door" }
    hours: "23:00-06:00"  # only at night
    action: "Night-time motion at front door. Check camera, assess threat."
    
  - subject: "agents.hb.>"
    timeout_s: 90  # no heartbeat in 90s = alert
    action: "Agent went silent. Investigate."
```

Implementation: A Hermes plugin that maintains NATS subscriptions. When triggered, it injects a user message into a session (new or existing). The brain reacts with full tool access — it can check docker, restart services, send notifications, escalate.

This replaces Uptime Kuma, custom alert scripts, and polling daemons with one plugin.

---

## 9. Flutter App Connection

The Flutter app (onemind-app) connects directly to Hermes. Three connection options:

| Method | Use Case | Protocol |
|--------|----------|----------|
| **Hermes WebUI WebSocket** | Chat, streaming responses, tool events | WebSocket JSON |
| **Hermes REST API** (:9119) | Sessions list, costs, analytics | HTTP REST |
| **NATS direct** (Synadia SDK) | Realtime fleet status, KV watches | NATS protocol |

The Flutter app is the **single interface** — web, iOS, macOS. It displays:
- Chat (talk to Legacy)
- Fleet (agent discovery via NATS $SRV)
- Tasks (from Kanban / KV)
- Infrastructure (docker, services, health — via Hermes tools or NATS direct)
- Memory (Hermes memory API)
- Analytics (Langfuse via Hermes tool or direct API)
- Voice (LiveKit WebRTC)

---

## 10. Edge Nodes & Robotics

```
┌─────────────────────────────────────────────────────────────────────────┐
│  HUB: MBP M2 Max (Zeus's machine)                                        │
│  NATS Server (hub) + Hermes (brain) + All services                       │
└───────────────────────┬─────────────────────────────────────────────────┘
                        │ Leaf node connections (over NetBird WireGuard)
        ┌───────────────┼───────────────┬───────────────────┐
        │               │               │                   │
┌───────▼───────┐ ┌─────▼─────┐ ┌───────▼───────┐ ┌───────▼───────┐
│ omos-dell     │ │ omos-core │ │ Barn Node     │ │ Drone RECON   │
│ (local server)│ │ (Lightsail)│ │ (future)      │ │ (future)      │
│               │ │           │ │               │ │               │
│ NATS leaf     │ │ NATS leaf │ │ NATS leaf     │ │ NATS leaf     │
│ + local JS    │ │           │ │ + local JS    │ │ + local JS    │
│               │ │           │ │ (offline ops) │ │ (store telem) │
└───────────────┘ └───────────┘ └───────────────┘ └───────────────┘

Key: "local JS" = local JetStream for offline operation.
Leaf nodes with JetStream can store-and-forward when disconnected.
When reconnected, mirrored streams sync automatically to hub.
```

For drones/robots:
- **NetBird** handles network roaming (WiFi → LTE → back)
- **NATS leaf** reconnects automatically when tunnel re-establishes
- **Local JetStream** stores telemetry/commands during disconnection
- **Object Store** distributes firmware updates fleet-wide
- **$SRV** registration makes them discoverable the moment they connect

---

## 11. Security Model

| Layer | Mechanism | Status |
|-------|-----------|--------|
| **Network** | NetBird WireGuard (all WAN encrypted) | ✅ Active |
| **Authentication** | NATS NKeys (per-node identity) | ✅ Active |
| **Isolation** | NATS Accounts (namespace separation) | ✅ Active |
| **In-transit** | NATS TLS | ✅ Active |
| **Event signing** | Synadia Security (signed events, audit) | 🔮 Future |
| **Per-device tokens** | NATS platform adapter device pairing | ✅ Built |

---

## 12. Build Priority

| # | What | Type | Effort |
|---|------|------|--------|
| 1 | **Watch services plugin** | Hermes plugin | 2-3 days |
| 2 | **Infrastructure tools** (docker, services) | Hermes tool additions | 1 day |
| 3 | **Langfuse query tool** | Hermes tool addition | 0.5 day |
| 4 | **NATS Object Store tools** | Hermes tool additions | 0.5 day |
| 5 | **$SRV discovery tool** | Hermes tool addition | 0.5 day |
| 6 | **Flutter → Hermes wiring** | Flutter app rework | 3-5 days |
| 7 | **Kill redundant services** | Cleanup | Progressive |

---

## 13. What NATS Replaces

| Previously | Now | How |
|-----------|-----|-----|
| Redis (cache/state) | NATS KV | Native, no extra service |
| Kafka (event streaming) | NATS JetStream | Streams with replay, consumers, retention |
| Consul/etcd (discovery) | NATS $SRV | Built-in service discovery + health |
| S3 internal (file dist) | NATS Object Store | For firmware/configs to edge nodes |
| RabbitMQ (task queues) | NATS Queue Groups | Automatic load balancing |
| Custom health checks | NATS $SRV.PING + $SYS events | Zero-config service health |
| HTTP between services | NATS Request/Reply | Sub-ms, no HTTP overhead |
| Custom sync daemons | JetStream Sources & Mirrors | Automatic cross-cluster replication |

---

## 14. The Rule

> **Everything communicates through NATS. If a service can't speak NATS natively, Hermes bridges it via a plugin or tool. There are no standalone workers where a plugin suffices. The Flutter app is the single UI. Hermes is the single brain. NATS is the single bus. NetBird is the single network.**
