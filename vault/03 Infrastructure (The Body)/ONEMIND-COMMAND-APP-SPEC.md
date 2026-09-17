---
title: "OneMind Command App Spec — moved to repo"
type: redirect
updated: 2026-06-12
---

# Moved to repo

This document now lives in the **`onemindos-command`** repository alongside the code it describes:

👉 `onemindos/onemindos-command` → [`docs/ONEMIND-COMMAND-APP-SPEC.md`](https://github.com/onemindos/onemindos-command/blob/main/docs/ONEMIND-COMMAND-APP-SPEC.md)

Keeping architecture docs co-located with the code they describe — easier to keep in sync, shows up in PRs alongside code changes.

type: spec
status: approved
created: 2026-06-12
updated: 2026-06-12
owner: zeus
category: infrastructure
domain: SO/03
tags: [onemind-command, flutter, nats, ui-spec, architecture, hermes-gateway, livekit, monorepo]
related:
  - "[[ONEMIND-OS-STACK]]"
  - "[[SERVICE-AUDIT-KILL-KEEP]]"
ai_summary: "APPROVED production spec for OneMind Command — the NATS-native Flutter operator surface for OneMind OS. Defines the 5 data planes, the 4 bus-relationships organizing all add-ons (feed/persist/bridge/observe), the consolidated ~19-page IA using Synadia's pillar names (Secure/Manage/Observe/Connect + Insights), per-page front-end + back-end + interactive verbs, the Meta Agent (Legacy via Hermes gateway), the watch service, LiveKit Flutter SDK, the APPROVED monorepo (onemindos via git subtree), and the phased build of ALL groups. Mission Control + SwarmClaw were idea-references only, NOT dependencies — every feature is NATS-native or Hermes-gateway."
---

# OneMind Command — Full Application Specification

> **What this is:** the APPROVED blueprint for *every* page, its **front end**, its **back end**, and what it **does** (interactive verbs, not view-only). Built before writing Flutter.
>
> **Prime doctrine:** **NATS is the standard.** The bus IS the API. Every add-on supercharges NATS via one of four relationships (feed / persist / bridge / observe) and surfaces as a *feature/tab* on the relevant NATS page — never as page-bloat. Room is left for everything we'll add later without limiting us now.
>
> **Sovereignty rule (DECISIONS #4 + #6):** OneMind is an **OS, not an agent system**. Mission Control and SwarmClaw were *idea references only* — NOT runtime dependencies. Every feature routes **NATS-native** or through the **Hermes gateway**. No third-party control plane in the critical path.
>
> **Hard rule:** No mock data. Every panel binds to a real source; offline sources render an honest "offline/empty" state.

---

## 0. The bus-relationship model (why we won't oversaturate)

Every component of OneMind OS has exactly **one primary relationship** to the NATS bus. This decides whether something is a page, a tab, or a card.

| Relationship | Meaning | Components | Where it lives in UI |
|---|---|---|---|
| **FEED** | produces messages onto the bus | Python services, Wombat ingest, edge nodes, robots, voice agent, sensors | **Connect** + **Fleet** |
| **PERSIST** | makes bus state durable | JetStream, ClickHouse, QuestDB, agentmemory, **Litestream**, Object Store | tabs on **Manage** + **Insights** |
| **BRIDGE** | translates bus ↔ external world | saas-action-worker, n8n, **LiveKit**, ha-bridge | **Connect** + **Voice** |
| **OBSERVE** | reads bus/system state | `:8222` monitor, **$SYS**, **nats-surveyor→Prometheus→Grafana**, **Langfuse** | tabs on **Observe** + **Insights** |

**Rule:** a new add-on does NOT get a new top-level page by default. It attaches as a tab/feature to the page that owns its relationship. New *pages* are reserved for genuinely distinct operator functions. This holds us at ~19 focused pages even as the stack grows.

---

## 1. Data planes — how FE ↔ BE works

Five planes. This is the whole transport contract.

| Plane | Transport | Use | Direction |
|---|---|---|---|
| **P1 · Bus-live** | NATS over **WebSocket :8080** | live subscriptions to `>` subjects; request-reply; publish | bidirectional |
| **P2 · Bus-monitor** | HTTP **:8222** (`varz/connz/jsz/routez/leafz/gatewayz/accountz/healthz`) | server introspection | read |
| **P3 · Service-HTTP** | REST per service | ClickHouse :8124, QuestDB :9000, agentmemory :3111, Langfuse :3000, n8n :5678, Prometheus :9090, LiveKit :7880, crucix :3117 | read + write |
| **P4 · Gateway-RPC** | **wss** to **Hermes gateway** | agent control, chat, sessions, skills, cron, PTY | request-reply |
| **P5 · Voice** | **wss** LiveKit :7880 + WebRTC | realtime audio rooms, tracks, transcripts | bidirectional |

> **Note:** Mission Control REST (`/api/*`) is **removed** from the contract. Everything it would have served (sessions, skills, cron, memory, alerts, audit, search) is now **P1 (NATS-native)** or **P4 (Hermes gateway)**.

### The interactive contract — how a page "does things"
Not view-only. Every action maps to:
- **NATS request-reply** — `nats.request("cmd.agent.legacy.wake", …)` → reply. *Primary verb.*
- **NATS publish** — fire onto `cmd.>` / `evt.>`.
- **NATS JetStream API** — `$JS.API.STREAM.CREATE/PURGE/DELETE`, KV/OBJ ops.
- **Hermes gateway RPC** — agent/session/skill/cron control over wss.
- **Service REST write** — run an n8n workflow, query ClickHouse.
- **WebRTC** — publish mic track (voice).

### Watch service (cross-cutting feature)
NATS-native **watch** is first-class, surfaced wherever useful:
- **Subject watch** — subscribe to any subject pattern, live-tail (The Wire).
- **KV watch** — `bucket.watch(key)` → UI updates the instant a value changes (zeus-context, agent-state).
- **Stream/consumer watch** — lag/ack alerts.
- Any watch is **promotable to an Alert Rule** (Ops → Alerts) in one click → routes to `evt.alert.>` + notify.

### Meta Agent (DECISION #3 — Hermes)
A first-class **Meta Agent** page: **Legacy**, fronted by the **Hermes gateway**, with **full bus visibility and command authority** — reads any subject, queries any store, dispatches to any node. The platform's own consciousness: ask "what's happening across the OS?" or "spin down the barn node" and it reasons + acts across all planes. Distinct from per-agent chat (scoped to one agent). Lives in **COMMAND**.

### UI libraries (approved)
- **LiveKit Flutter SDK** `livekit_client` + `livekit_components` (prebuilt room/participant/audio widgets) for Voice. Python side: `livekit-agents` for the dual-conn voice agent.
- **Riverpod** (state), **go_router** (deep nested routing), **fl_chart** (Insights), **dio** (P3 REST), **nats**/`web_socket_channel` Dart client (P1). Custom `CustomPainter` for Topology/fleet graphs + waveforms.

---

## 2. Code location — monorepo `onemindos` (DECISION #1 — APPROVED ✅)

**Merging** `onemindos-core` (back end) + `onemindos-command` (Flutter) into one monorepo via `git subtree` (both histories preserved):

```
onemindos/
├── schemas/    ← SINGLE SOURCE OF TRUTH: subject taxonomy + message envelope + account model
│                  codegen → Dart (command) · Go (omd/edge) · Python (services)
├── core/       ← back end (was onemindos-core): nats config, services, omd edge binary, control-plane
├── command/    ← Flutter app (was onemindos-command): UI — native (macOS/iOS) + web build
├── ops/        ← launchd plists, install scripts, deploy
└── docs/       ← architecture, doctrine, this spec
```

**Why it wins:** in a NATS-native system **the message envelope IS the API.** Define `cmd.robot.{id}.estop` once in `schemas/`, regenerate, and the robot handler (Go/Python) *and* the Flutter E-STOP button (Dart) both get the typed contract — **they can never drift.** Cross-repo, that drift is the #1 production-bug source. One operator, one portable unit (clone once → whole OS) → matches sovereignty/portability.

**Not in the monorepo:** Mission Control and SwarmClaw — idea references only (DECISIONS #4/#6), not dependencies. Reusable services like n8n, agentmemory, LiveKit run as their own processes and are reached over P3/P5.

---

## 3. Information Architecture — ~19 pages, Synadia-pillar names (DECISION #2)

Deep, grouped, collapsible rail. The BUS group uses **Synadia's exact pillar names: Secure · Manage · Observe · Connect.** Add-ons are **tabs** (‹angle brackets›), not pages.

```
⌂ COMMAND
   • Command Center        /                 ← merged core+dashboard: live pulse, fleet, agents, services, alerts
   • The Wire              /wire             ‹Stream | Watch | Publish | Inspect›  (agenttap absorbed)
   • Meta Agent            /meta             ← Legacy via Hermes gateway, full bus authority

⚙ BUS  (NATS core — Synadia pillars)
   • Observe               /bus/observe      ‹Servers | Connections | Metrics(Surveyor/Prometheus) | Health›
   • Manage                /bus/manage       ‹Streams | Consumers | KV | Object Store | Backups(Litestream)›
   • Secure                /bus/secure       ‹Accounts | Users | Permissions | Audit›
   • Connect               /bus/connect      ‹Connectors(Wombat) | Workflows(n8n) | SaaS Actions | Integrations›
   • Topology              /bus/topology     ← cloud→edge mesh graph, animated flow

🧠 AGENTS  (intelligence)
   • Roster                /agents           → /agents/:id ‹Chat | Soul | Memory | Diagnostics | Activity›
   • Sessions              /agents/sessions  → /:id ‹Transcript | Control›
   • Memory                /agents/memory    ‹Layers | Graph | RAG Search›
   • Skills & Schedules    /agents/skills    ‹Skills | Cron | Runs/Evals›
   • Voice                 /agents/voice     → /room/:name  (LiveKit Flutter SDK)

🛰 FLEET  (nodes + hardware)
   • Map                   /fleet            → /node/:id
   • Robotics              /fleet/robotics   → /:id ‹Telemetry | Control | Camera | E-STOP›
   • Drones                /fleet/drones     → /:id ‹Mission | GPS | Camera›
   • Edge Nodes            /fleet/edge       → /:id  (embedded-NATS leaf)

🛡 OPS  (operations + governance)
   • Services              /ops/services     ← live health + start/stop (replaces Homepage/Kuma)
   • Insights              /ops/insights     ‹Events | Telemetry | Cost | Tracing(Langfuse)›
   • Alerts & Audit        /ops/alerts       ‹Rules | Firing | Approvals | Audit Trail›

⚙ SYSTEM
   • Settings              /settings         ‹Connections | Appearance | Access/Auth›
   • Page Map              /map
```

**Count:** 19 pages. The old standalone BRIDGE group is folded into **Connect** (Synadia's connector pillar). Robotics/Drones/Edge/Integrations are first-class but render honest-empty until hardware/keys arrive — room left open by design.

---

## 4. Page-by-page — Does · Frontend · Backend(plane+source) · Live

### ⌂ COMMAND

**Command Center `/`** — *merged core + dashboard*
- **Does:** single pane of glass; click any tile → drill to its page; quick actions (wake agent, new task, talk to Meta Agent, voice); pin/reorder tiles (persisted to KV).
- **Frontend:** `StatTile` grid + mini `WireFeed` + `ServiceHealthStrip` + `FleetMiniMap`; drag-reorder.
- **Backend:** P2 `varz` (msgs/sec, conns) · P1 `>` (wire) · P3 health-poll · P1 `tel.node.>` · P3 ClickHouse (24h).
- **Live:** ✓ msgs/sec, JetStream totals, agents, fleet — real.

**The Wire `/wire`** — *agenttap absorbed*
- **Does:** watch every message live; **filter** by subject/account/agent; **pause/resume + time-scrub**; click → payload JSON tree; **publish** test msg; **Watch** tab pins a subject/KV key → live alert; copy subject.
- **Frontend:** virtualized list (thousands/sec), filter bar, payload drawer, pub composer, scrubber (agenttap-style).
- **Backend:** P1 subscribe `>` / filtered over WS :8080; P1 publish; KV watch for the Watch tab.

**Meta Agent `/meta`** — *Legacy via Hermes gateway*
- **Does:** converse with the platform's own agent; it reads any subject, queries any store, commands any node; ask "what's happening" / "do X across the OS"; shows its reasoning + the bus actions it takes live.
- **Frontend:** chat + live "action tape" (subjects published, stores queried) + grammar-aware command bar.
- **Backend:** P4 Hermes gateway RPC to Legacy + P1 `cmd.>`/reply + read across P2/P3.

### ⚙ BUS  (Synadia pillars)

**Observe `/bus/observe`** — ‹Servers | Connections | Metrics | Health›
- **Does:** monitor every server CPU/mem/RTT/conns/subs/slow-consumers; **disconnect** a client; sort/filter conns; Metrics tab embeds Surveyor→Prometheus; Health tab shows JetStream recovery (the `$G` fault would show here).
- **Backend:** P2 `varz`/`connz` + **$SYS.REQ.SERVER.PING** (SYS creds, server-side only) + P3 Prometheus `:9090/api/v1/query_range` + P2 `healthz`.
- **Live:** ✓ 17 conns, 237 subs, RTT real.

**Manage `/bus/manage`** — ‹Streams | Consumers | KV | Object Store | Backups›  *(Synadia "Manage")*
- **Does:** **create/edit/delete/purge** streams; inspect retention/replicas/subjects; **consumers** (lag/ack, reset/delete); **KV** bucket CRUD + key editor + **watch**; **Object Store** browse/upload/download; **Backups** tab = Litestream replication status + restore trigger.
- **Backend:** P2 `jsz?streams&consumers` (read) + **$JS.API.*** (writes) + KV/OBJ APIs + P3 Litestream metrics.
- **Live:** ✓ 34 streams, 72,567 msgs real.

**Secure `/bus/secure`** — ‹Accounts | Users | Permissions | Audit›  *(Synadia "Secure")*
- **Does:** 6-account topology; per-account users/imports/exports/permissions; **add user**, **rotate creds** (redacted), edit subject perms; auth audit log.
- **Backend:** P2 `accountz` + config introspection + NSC/config push + reload + P1 `evt.audit.auth.>`.

**Connect `/bus/connect`** — ‹Connectors(Wombat) | Workflows(n8n) | SaaS Actions | Integrations›  *(Synadia "Connect" — absorbs old BRIDGE)*
- **Does:** manage **Wombat** ingest pipes + saas-action-worker out (status/throughput/last-event, YAML editor, start/stop, test-fire); **n8n** workflows (list/activate/run/executions, deep-link to editor); **SaaS Actions** (browse 2000+, execute, history/latency ~154ms); **Integrations** = future native per-platform API buttons (Skool/GHL) — catalog now, actions later.
- **Backend:** P3 Wombat `:4195/stats` + P1 `saas.>` + P1 `saas.action.>` request-reply + P3 n8n `:5678/rest/*`.

**Topology `/bus/topology`**
- **Does:** live mesh graph hub→leaf→routes→gateways; animated message-flow particles; click node → its conns; filter by account; clusters/zoom for hundreds of nodes.
- **Backend:** P2 `routez`/`leafz`/`gatewayz` + $SYS server map + P1 `tel.node.>`. `CustomPainter` force-graph.

### 🧠 AGENTS  *(all NATS-native or Hermes gateway — no Mission Control)*

**Roster `/agents` → `/agents/:id`**
- **Does:** all minds (Legacy + 5 profiles + coding agents + robots); live status from heartbeats; per-agent **Chat | Soul | Memory | Diagnostics | Activity**; **wake**, assign task. Merges old agents/chat/chatroom/codex.
- **Backend:** P1 `agents.hb.>` (heartbeats) + P4 Hermes gateway (chat/wake/diagnostics) + P1 `cmd.agent.>`.

**Sessions `/agents/sessions` → `/:id`**
- **Does:** list agent sessions; **monitor/pause/terminate**; **continue** a claude-code/codex session with a prompt; full transcript.
- **Backend:** **P4 Hermes gateway** (Hermes owns sessions natively) + P1 `evt.session.>` for live updates.

**Memory `/agents/memory`** — ‹Layers | Graph | RAG Search›
- **Does:** browse 4 layers (agentmemory semantic / JS-KV meta / ClickHouse analytical / QuestDB telemetry); link **graph**; **RAG search** (Tantivy BM25); add/edit/delete entries.
- **Backend:** **P3 agentmemory `:3111`** (53 tools / REST) + P1 JetStream KV + P3 Tantivy (via om-rag `svc.rag.query`) + P3 ClickHouse/QuestDB.

**Skills & Schedules `/agents/skills`** — ‹Skills | Cron | Runs/Evals›
- **Does:** skills list/view/**upsert/delete/check**; cron **create/pause/resume/run-now** + human-schedule parse; runs stream + provenance + eval scores.
- **Backend:** **P4 Hermes gateway** (skills + cronjob are Hermes-native) + P1 `evt.run.>` (run events) → ClickHouse for history/evals.

**Voice `/agents/voice` → `/room/:name`** — *LiveKit Flutter SDK*
- **Does:** list rooms/participants/tracks; **join + push-to-talk** (publish mic); live **transcript**; mute/remove participants; waveform; watch the dual-conn agent bridge LiveKit↔NATS.
- **Frontend:** `livekit_client` + `livekit_components`, waveform painter, transcript stream.
- **Backend:** P5 LiveKit wss `:7880` + token mint (server SDK) + P1 `voice.{session}.transcript.final`.
- **Live:** ⚠ LiveKit not booted → honest "offline" until its phase starts `livekit-server`.

### 🛰 FLEET

**Map `/fleet` → `/node/:id`** — geo+logical all-nodes; live telemetry; **send command**. P1 `tel.node.>` + P2 `leafz` + P1 `cmd.node.>`.
**Robotics `/fleet/robotics` → `/:id`** — ‹Telemetry | Control | Camera | E-STOP›; joints/battery/mode; **move/dock/task**; **E-STOP** (`cmd.robot.{id}.estop`, confirm-guarded); Guardian state; camera. P1 `tel.robot.>` + `cmd.robot.>` + P5 camera. DIMOS/OM1.
**Drones `/fleet/drones` → `/:id`** — ‹Mission | GPS | Camera›; upload/start/abort path; RTH; map overlay. P1 `tel.drone.>` + `cmd.drone.>`. Honest-empty until first drone.
**Edge Nodes `/fleet/edge` → `/:id`** — embedded-NATS leaf nodes; **provision/monitor/push**; path to hundreds. P2 `leafz` + P1 `tel.node.>` + `cmd.node.{id}.exec`.

### 🛡 OPS  *(NATS-native — no Mission Control)*

**Services `/ops/services` → `/:name`** — live health (port/status/uptime); **start/stop/restart** (launchctl); logs. Replaces Homepage + Uptime-Kuma. P3 health-poll each endpoint + P4 launchctl RPC (or P1 `cmd.svc.>`). ✓ 15 services real.
**Insights `/ops/insights`** — ‹Events | Telemetry | Cost | Tracing›; time-series + funnel + LLM cost + Langfuse traces; range picker; CSV export. P3 ClickHouse `:8124` (166 tables ✓) + QuestDB `:9000/exec` + Langfuse `:3000`. `fl_chart`.
**Alerts & Audit `/ops/alerts`** — ‹Rules | Firing | Approvals | Audit›; define watch→alert rules (stored in **KV** `om-alerts`); **ack/silence**; **approve/deny** risky execs (L1/L0 gate via `cmd.approval.>`); immutable **Audit** = JetStream stream `AUDIT` (`evt.audit.>`). All P1 NATS-native + P1 `evt.alert.>`.

### ⚙ SYSTEM

**Settings `/settings`** — ‹Connections | Appearance | Access›; edit all endpoints; theme/accent/density/scanlines; profiles; OIDC/tailnet. Local config (persisted to KV `om-config`) + P4 gateway caps.
**Page Map `/map`** — full IA reference, annotated with live-data source per page.

---

## 5. Build phases (DECISION #5 — ship ALL groups)

All six groups are in v1 scope. Phased purely by **data-readiness**, not by cutting scope.

| Phase | Scope | Why this order |
|---|---|---|
| **P0** | Monorepo `onemindos` set up + HUD theme + nav shell + 5 data clients (P1–P5) + Riverpod + `schemas/` codegen | foundation + contract |
| **P1** | COMMAND (Command Center, The Wire) + BUS (Observe, Manage, Secure, Connect, Topology) + OPS (Services) | data verified live today |
| **P2** | AGENTS (Roster, Sessions, Memory, Skills, **Meta Agent**) | Hermes gateway wiring |
| **P3** | Boot `livekit-server` → Voice; finish Connect tabs (Wombat/n8n live) | brings new services online |
| **P4** | FLEET (Map, Robotics, Drones, Edge) | grows with hardware; honest-empty until then |
| **P5** | OPS hardening (Insights, Alerts & Audit) + SYSTEM (Settings, Page Map) | governance + polish |

Each phase ends: `flutter analyze` clean + `flutter build web` green + live click-through. **Docs updated every phase.**

---

## 6. Decisions — LOCKED (2026-06-12)

1. **Monorepo** ✅ APPROVED — merge into `onemindos` via git subtree; single `schemas/` codegen.
2. **Naming** ✅ — Synadia pillar names: **Secure · Manage · Observe · Connect** (+ Topology, Insights).
3. **Meta Agent** ✅ — **Legacy via the Hermes gateway** (Hermes is authoritative; P4 = Hermes).
4. **Mission Control** ✅ — **NOT a dependency.** Idea reference only. All features → NATS-native or Hermes gateway.
5. **Scope** ✅ — **ship ALL groups**; phase by data-readiness (P0→P5), no scope cut.
6. **SwarmClaw** ✅ — **NOT a dependency.** Idea reference only. OneMind is an OS, not just an agent system.

**Next:** P0 — set up the `onemindos` monorepo + theme + nav shell + 5 data clients, then P1 pages whose data is verified live today.
