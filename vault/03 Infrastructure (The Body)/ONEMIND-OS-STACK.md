---
title: "OneMind OS — Canonical Stack"
type: architecture-doctrine
status: active
created: 2026-06-11
updated: 2026-06-12
owner: zeus
category: infrastructure
# Also in repo: onemindos/onemindos-command → docs/ONEMIND-OS-STACK.md
# Vault copy = infrastructure layer reference (NATS, services, fleet, data stores)
# Repo copy = full INTERFACE section + OneMind Command architecture
tags: [nats, stack, architecture, infrastructure, clickhouse, langfuse, livekit, doctrine, traefik, risingwave, flink, netbird]
ai_summary: "The canonical OneMind OS technology stack — NATS-native, local-first. Defines every layer (substrate, telemetry DB, analytical DB, memory, RAG, tracing, voice, services, SaaS bridge, backups, infra health, vault sync, edge, interface), what is LIVE vs INSTALLED vs PENDING as of the post-reset baseline, and the doctrine decisions behind each choice."
---

# OneMind OS — Canonical Stack

> NATS is the nervous system. Every byte flows through the bus. Everything builds up from it. Local-first, sovereign, portable. This document is the single source of truth for what the stack IS and what state each piece is in.

**Baseline date:** 2026-06-11 (post Mac-reset audit)

---

## The Stack at a Glance

| Layer | Choice | State |
|---|---|---|
| **Substrate** | NATS (hub: omos-mbp standalone; edge: embedded in Go binary) | 🟢 LIVE |
| **Telemetry DB** | ~~QuestDB~~ → **removed** (no native NATS; ClickHouse covers both) | ❌ REMOVED |
| **Analytical DB** | ClickHouse — long-term events, NATS engine table, SQL analytics | 🟢 LIVE |
| **Streaming SQL** | RisingWave — incremental materialized views from NATS subjects | 🟢 LIVE |
| **Stream Processing** | Apache Flink — CEP, cross-stream joins, stateful enrichment | 🟢 INSTALLED |
| **Agent Memory** | agentmemory (`:3111`) — BM25 + Vector + Graph hybrid, MCP→Hermes | 🟢 LIVE |
| **Meta Memory** | JetStream KV (zeus-context, om-config, om-schema) | 🟢 LIVE |
| **RAG** | Tantivy BM25 (vault) + agentmemory (sessions) — no Chroma | 🟡 PARTIAL |
| **Tracing** | Langfuse (local `:3000`, OTel ingest) | 🟢 LIVE |
| **Voice** | LiveKit (self-hosted, WebRTC SFU) + Python dual-conn agent | 🔴 PENDING |
| **Services** | Python `nats-py` asyncio micro services (`$SRV`-registered) | 🟡 PARTIAL |
| **SaaS Bridge** | bento v1.18.1 (WarpStream fork) — YAML pipelines → NATS | 🟢 INSTALLED |
| **Backups** | Litestream v0.5.12 → NATS Object Store (S3 gateway) | 🟢 INSTALLED |
| **Infra Health** | nats-surveyor (:7777) → Prometheus → Grafana + Netdata | 🟢 LIVE |
| **Ingress** | Traefik v3.7.5 — reverse proxy, HTTP :80, HTTPS :4443, dashboard :7070 | 🟢 LIVE |
| **Networking** | Tailscale → **migrating to NetBird** (cloud → self-hosted on VPS) | 🔄 MIGRATING |
| **IaC** | Terraform + nats-io/jetstream provider — streams, KV, consumers | 🟢 LIVE |
| **Vault Sync** | om-vault-sync (Python watchdog) → Object Store `OBJ_vault` | 🔴 PENDING |
| **Edge Nodes** | Go binary: embedded NATS + leaf + local cmd handler | 🔴 PENDING |
| **Interface** | OneMind Command (Flutter, NATS-WS) + React Cockpit (future) | 🔴 PENDING |

Legend: 🟢 LIVE = running · 🟡 INSTALLED/PARTIAL = binary present or partially built · 🔴 PENDING = not yet built

---

## Layer Detail

### SUBSTRATE — NATS  🟢
- Binary `/opt/homebrew/bin/nats-server` v2.14.1, LaunchAgent `io.nats.server`
- Ports: 4222 (client), 8222 (HTTP monitor), 8080 (WebSocket), 7422 (leafnode), 1883 (MQTT native)
- JetStream enabled. Config `~/.nats/config/onemind.conf`
- **Hub stays standalone** (omos-mbp, clamshell 24/7). Edge nodes embed NATS in a Go binary (`omd --embed`) as leaf nodes.
- See [[NATS-INFRASTRUCTURE]] for full server config, streams, KV buckets.

### TELEMETRY DB — ~~QuestDB~~  ❌ REMOVED (2026-06-13)
- **Decision:** QuestDB removed. Has no native NATS integration — requires a custom bridge process. ClickHouse's built-in NATS table engine covers the same ingest workload with zero glue code.
- ClickHouse `NATS` engine table subscribes directly to `tel.>` via JetStream durable consumer → Materialized View → MergeTree. No separate telemetry DB needed.
- See `STACK.md` for ClickHouse NATS ingest pattern.

### ANALYTICAL DB — ClickHouse  🟢
- LaunchAgent `com.onemind.clickhouse` (PID live). Ports 8124 (HTTP), 9001 (TCP) — separate from Langfuse's Docker ClickHouse (8123/9000).
- Native **NATS engine table** subscribes to `events.>` → `onemind.events` (zero-glue ingestion). Live: ~31K events, ~34M metrics.
- Role: long-term (10yr) event store, SQL analytics, "Time Machine" cross-domain queries.
- See `nats-server` skill → `references/clickhouse-experience-store.md`.

### AGENT MEMORY — agentmemory  🟢
- LaunchAgent, `localhost:3111` (API), `:3113` (viewer). v0.11.2.
- BM25 + Vector + Graph hybrid retrieval (RRF fusion). SQLite + iii-engine, no external deps.
- MCP connected to Hermes (`memory.provider=agentmemory`, 53 MCP tools).
- State store `~/data/` — **DO NOT DELETE.**

### META MEMORY — JetStream KV  🟢
- Buckets live: `zeus-context` (location/focus/status), `agent-state`, `config`, `notifications`. Planned: `om-config`, `om-schema`, `agent-facts`.
- Watchable → services react instantly to changes. Revisions, TTL, OCC.

### RAG — Tantivy + agentmemory  🟡 PARTIAL
- **Tantivy 0.26.0 installed** in Hermes venv (`~/.hermes/hermes-agent/venv`). For BM25 search over the Obsidian vault / documents.
- agentmemory handles session/agent-history hybrid retrieval.
- **No Chroma, no external vector DB** — benchmarks showed dense-only Chroma (56.4%) < pure BM25 (65.4%) < hybrid BM25+vector (75.1%). agentmemory already does hybrid correctly.
- Pending: `om-rag` Python service exposing `svc.rag.query`.

### TRACING — Langfuse  🟢
- Local self-hosted `:3000` (Node listener confirmed). OTel ingest.
- Agent traces, LLM calls, tool use, costs. Has its own Docker ClickHouse (8123/9000).

### VOICE — LiveKit  🔴 PENDING
- Self-hosted WebRTC SFU (target: omos-cloud VPS). **LiveKit does NOT expose NATS subjects natively** — its internal psrpc uses NATS but it's not user-facing; Redis is its coordination layer.
- Integration: a single Python LiveKit agent holds **two** connections — the LiveKit room AND the NATS bus.
- Transcripts → `voice.{session}.transcript.final`; intent → `cmd.{actor}.task.execute`; every utterance → `evt.so.voice.utterance` → ClickHouse.
- Realtime models (GPT-4o Realtime / Gemini Live) wire directly into the room; NATS publish still fires.
- Docs already in vault: `services/livekit/docs/VOICE-PIPELINE.md`, `LIVEKIT-SELF-HOSTED-REALITY.md`.

**DEPLOYMENT DECISION (2026-06-11): launchd/systemd binary, NOT Docker** — same doctrine as nats-server.
- LiveKit server is a single static Go binary (`brew install livekit`, v1.13.1). Run it as a service (launchd on Mac, systemd on the VPS), not a Docker container — Docker daemon crash = voice down, and the binary has zero runtime deps.
- **No Redis for single-node.** Redis is only required for distributed/multi-node LiveKit (multiple SFU replicas coordinating room state). One node = standalone, no Redis. OneMind voice is single-node → no Redis.
- **Exception — Egress (room recording):** LiveKit's *egress* service is normally run as a Docker container because it needs headless Chrome for compositing. This is OPTIONAL and only needed if we want recorded voice sessions. Core voice (server + Python agent + STT/LLM/TTS) needs no Docker at all.
- **Three-layer split:** (1) `livekit-server` binary as a service; (2) Python agent (`nats-py` + LiveKit SDK, dual connection) as a service; (3) egress in Docker only if recording is wanted later.

### SERVICES LAYER — Python micro services  🟡 PARTIAL
- Pattern: `nats-py` asyncio, registered as NATS micro services (`$SRV.PING/INFO/STATS`). nats-py present in `/opt/homebrew/lib/python3.13|3.14/site-packages`.
- **Live now:** 5 subscribers (agents/comms/health/home/security), notify-router, saas-action-worker, ha-bridge, synadia-agent-ui.
- **Target roster:** `om-meta`, `om-insights`, `om-rag`, `om-saas`, `om-voice`, `om-vault-sync`.
- **No Go daemon on the hub** — `omd` as originally scoped is redundant; `$SRV` micro API already provides discovery/registration/stats. Go is reserved for embedded-NATS edge binaries only.

### SaaS BRIDGE — Wombat  🟢 INSTALLED (not yet wired)
- **Binary:** `~/bin/wombat` (v1.0.7 source, built from `./cmd/wombat`, 344 MB static Go binary, built 2026-06-12).
- **Identity:** the MIT **community fork** of RedPanda Benthos (`wombatwisdom/wombat`) — NOT Synadia Connect (whose standalone mode requires Docker → rejected, violates no-Docker-on-critical-path).
- **Verified components:** NATS in/out (core, `nats_jetstream`, `nats_kv`, `nats_stream`), `http_server` input for inbound webhooks. 68 inputs / 80 outputs / 88 processors.
- **ROLE (locked): inbound ingest pipe ONLY.** webhook (Skool/GHL) → Wombat (`http_server` + Bloblang transform) → `saas.<provider>.<event>`. Replaces the `rpconnect` binary wiped in the Mac reset.
- **NOT a replacement for:** (a) `saas-action-worker` v4.1.0 — owns *outbound* actions at ~154 ms native, faster than Wombat; (b) **n8n** (`:5678`) — owns stateful multi-step workflows (branching, retries, human-in-loop) which Wombat cannot do.
- **Three-layer split:** Wombat = ingest pipe · saas-action-worker = outbound actions · n8n = workflow logic.
- **Pending:** write `~/.hermes/bridge/*.yaml` connector config + launchd plist `com.onemind.wombat-ingest`.

### BACKUPS — Litestream  🟢 INSTALLED
- v0.5.12 installed via brew. Binary `/opt/homebrew/bin/litestream`.
- Config: `~/.config/litestream/litestream.yml`
- NATS Object Store path: use S3-compatible gateway (natss3proxy) or local path + `nats object put`.
- Wire per-service: add each SQLite DB path to litestream config, run `litestream replicate`.

### INFRA HEALTH — Surveyor + Prometheus + Grafana  🟢 LIVE
- **nats-surveyor** installed (`go install`), LaunchAgent `io.nats.surveyor` running — exposes NATS `$SYS` internals as Prometheus metrics on `:7777/metrics`.
- **Prometheus v3.12** running (brew service). Scrapes nats-surveyor at `:7777`.
- **Grafana v13** running (brew service). Sources: Prometheus + ClickHouse direct query.
- **Netdata** running — real-time 1s host metrics. Install agent on every edge device; streams to parent Mac.
- **Uptime Kuma** running — endpoint uptime checks.
- SYS account user configured in `~/.nats/config/onemind.conf` — monitoring API unblocked.

### VAULT SYNC — om-vault-sync  🔴 PENDING
- Python watchdog on `~/onemind/` → NATS Object Store bucket `OBJ_vault`. Any node pulls any vault file on demand.
- Fires `evt.so.vault.file.updated` on change. Object Store: chunked, versioned, watchable.

### EDGE NODES — Go embedded binary  🔴 PENDING
- One static binary (~3-5 MB, cross-compiles to ARM): embedded NATS server + leaf node → hub `:7422` + local `cmd.node.{id}.exec` handler.
- `omd --embed` mode. Hub never embeds. First real use: first Raspberry Pi / barn node.
- See `nats-server` skill → `references/omd-embedded-nats.md`.

### INTERFACE — OneMind Command  🟢 SPEC LOCKED → BUILDING
- **Flutter app** (native macOS/iOS + `flutter build web` — one codebase, two outputs). Connects NATS over WebSocket (`ws://omos-mbp:8080`) directly — the bus IS the API.
- **Full spec:** [[ONEMIND-COMMAND-APP-SPEC]] — **19 pages**, every page's front-end + back-end + interactive verbs, organized by the **bus-relationship model**. All 6 decisions LOCKED 2026-06-12.
- **Bus-relationship model (anti-oversaturation):** every add-on has ONE relationship to the bus — **FEED** (services/Wombat/edge/robots/voice) · **PERSIST** (JetStream/ClickHouse/QuestDB/agentmemory/Litestream/ObjStore) · **BRIDGE** (saas-worker/n8n/LiveKit/ha-bridge) · **OBSERVE** (monitor:8222/$SYS/Surveyor→Prometheus/Langfuse). Add-ons surface as **tabs on the relevant NATS page, never as new pages.** Backups = tab on Manage; Metrics = tab on Observe; Tracing = tab on Insights.
- **IA uses Synadia pillar names:** BUS group = **Observe · Manage · Secure · Connect** (+ Topology). Groups: COMMAND, BUS, AGENTS, FLEET, OPS, SYSTEM.
- **5 data planes:** P1 Bus-live (NATS WS :8080) · P2 Bus-monitor (HTTP :8222) · P3 Service-HTTP (REST) · **P4 Gateway-RPC (Hermes gateway, wss)** · P5 Voice (LiveKit + WebRTC).
- **SOVEREIGNTY (decisions #4/#6):** Mission Control + SwarmClaw were **idea references only, NOT dependencies.** Every feature is **NATS-native or Hermes-gateway** — no third-party control plane in the critical path. OneMind is an OS, not just an agent system.
- **Meta Agent:** first-class page — **Legacy via the Hermes gateway** (Hermes authoritative); full bus visibility + command authority; the platform's own consciousness.
- **Watch service:** NATS subject/KV/stream watch is first-class; any watch promotable to an Alert Rule.
- **UI libs (approved):** LiveKit Flutter SDK (`livekit_client` + `livekit_components`) for Voice; Riverpod, go_router, fl_chart, dio, nats Dart client; CustomPainter for Topology/fleet graphs.
- **CODE HOME → monorepo `onemindos`** (APPROVED): merge `onemindos-core` (→ `core/`) + `onemindos-command` (→ `command/`) via git subtree. Single `schemas/` (envelope + subject taxonomy) → codegen to Dart/Go/Python so the message contract can never drift. Mission Control NOT in monorepo.
- **Scope:** ship ALL groups; phase by data-readiness P0→P5 (no scope cut).
- **First build (P1, data verified live today):** Command Center, The Wire, Observe, Manage, Secure, Connect, Topology, Services.
- Visual seed: [[nats-fabric-vision]]. Theme: agenttap HUD (Chakra Petch + IBM Plex Mono, cyan/orange/magenta, scanlines, grid floor, corner-bracket panels).

---

## Install / Build Ledger (post-reset)

### ✅ Done (2026-06-11 → 2026-06-13)
- Prometheus v3.12.0 (brew, running)
- Grafana v13.0.2 (brew, running)
- Tantivy 0.26.0 (Hermes venv)
- Litestream v0.5.12 (brew)
- nats-surveyor (go install, LaunchAgent running on :7777)
- bento v1.18.1 (brew — WarpStream fork of Benthos, binary `bento`)
- Traefik v3.7.5 (brew, LaunchAgent running — :80/:4443/:7070)
- RisingWave v2.8.0 (brew, LaunchAgent running — single-node, port 4566)
- Apache Flink v2.2.1 (brew, LaunchAgent running — JobManager port 8082)
- Terraform v1.15.6 + nats-io/jetstream provider (JetStream IaC at `~/onemind/config/nats/terraform/`)
- DuckDB v1.5.3 installed then **removed** (no NATS extension exists)
- QuestDB **removed** (no native NATS integration)
- n8n **removed** (blank install, no community nodes, Synadia nodes not installed)
- SYS account user added to NATS config (monitoring API unblocked)

### 🔧 Pending config
- Add **SAAS** + **VOICE** accounts to NATS config
- Wire Litestream per-service (add each SQLite DB to litestream.yml)
- `terraform apply` JetStream streams + KV buckets
- Wire bento pipelines (`~/.config/wombat/*.yaml`) per SaaS integration
- Migrate networking: Tailscale → NetBird Cloud (Free tier, 100 machines)
- Move configs into `~/onemind/config/` and symlink

### 🩹 Needs attention
- `com.onemind.graphite-bridge` — failing. Logs: `subscribers/logs/graphite-bridge-stderr.log`
- `com.onemind.notify-router` — failing. Logs: `subscribers/logs/notify-stderr.log`
- Several Hermes gateways flapping — separate from core stack
- `io.onemind.opentakserver` + `opentakserver-ui` — both exit 1, investigate or remove

---

## Doctrine Decisions (locked)

1. **NATS is the only substrate.** Every node speaks NATS. The subject is the address; identity in the subject, properties in the registry.
2. **Python is the primary service runtime.** Go is reserved for embedded-NATS edge binaries only. TypeScript for cockpit. Dart for the Flutter app.
3. **omd is over-scoped** — drop the hub Go daemon; `$SRV` micro API replaces it. Keep Go only for `omd --embed` on edge.
4. **Temporal.io: hard NO** — JetStream workqueue retention replaces it.
5. **ClickHouse only for storage** — QuestDB removed. ClickHouse native NATS table engine covers telemetry ingest with zero bridge code. No separate telemetry DB.
6. **No Chroma** — hybrid BM25+vector (agentmemory) + Tantivy BM25 (vault) is sufficient and local.
7. **Prometheus/Grafana ARE in the stack** for bus-internal health via nats-surveyor. Netdata for real-time host metrics per device.
8. **The bus is the API** — UIs subscribe to subjects, no REST/GraphQL layer.
9. **NetBird over Tailscale** — NetBird cloud now, self-hosted on VPS later. Same client, zero re-architecture migration. Headscale rejected (explicitly not for high-churn fleets).
10. **bento over n8n** — n8n removed. bento handles all NATS-native SaaS/IoT bridging via YAML. No runtime dependency.
11. **RisingWave for streaming SQL, Flink for CEP** — RisingWave (Phase 5) handles rolling aggregations as incremental materialized views. Flink (Phase 6) reserved for cross-stream temporal joins and complex event processing.
12. **Per-device footprint: 65MB** — NATS leaf (15MB) + Netdata agent (30MB) + NetBird client (20MB). Wombat and Litestream live on the hub only.

---

## Related

- [[nats-fabric-vision]] — interactive visual of this stack
- [[NATS-INFRASTRUCTURE]] — live NATS server documentation
- `~/onemindos-core/docs/NATS-DOCTRINE.md` — binding subject law
- `~/onemindos-core/docs/MASTER-PLAN.md` — phase plan
- `nats-server` skill — ops knowledge + all reference docs
