# MASTER-PLAN.md — OneMind Command & OS Build Program

> **The single source of truth for what we are building, why, in what order, and who does what.**
> Architect/reviewer: Legacy (Hermes default profile). Workers: Sonnet 4.6 agents (Opus 4.8 for flagged-hard tasks).
> Worker agents: read this file + `NATS-DOCTRINE.md` + `AGENTS.md` before touching anything.
> Status: PLANNING LOCKED → EXECUTION PHASED. Last updated: 2026-06-10.

---

## 0. Vision (Why This Exists)

OneMind OS is a **sovereign, local-first operating system for a homestead-scale fleet**:
hundreds of nodes (servers, Macs, Pis, ESP32 sensors, robots, vehicles, drones) across
House / Barn / Office, coordinated by heterogeneous AI agents, owned and operated by one
person. No cloud dependency on the critical path. Everything speaks NATS.

**OneMind Command** (this repo) is the human interface to that OS — one Flutter codebase
shipping to iPhone, Android, web (Brave), and Linux desktop.

We are deliberately building our own equivalents of Synadia's commercial suite
(Platform/Control Plane, Insights, Protect) on top of the open primitives they publish
(NATS, NEX, Wombat, Orbit, Agent SDK). See `NATS-DOCTRINE.md` §7 for the reuse-vs-replicate map.

### Non-negotiable principles

1. **Bus-first.** New capability = NATS subjects first, UI second. The app renders the bus.
2. **Harness sovereignty.** Agents keep their own configs. App is review/edit only (Doctrine §8).
3. **Local-first.** Mesh-internal (Tailscale) by default. Funnel exposure is per-route opt-in.
4. **SQLite everywhere** server-side state is needed. No Postgres until a real multi-writer need exists.
5. **Boring foundations, ambitious surface.** Stability budget goes to the bus + protocol; experimentation budget goes to UI.
6. **One codebase per surface.** Flutter app (mobile+web+desktop). React reserved for a future big-screen cockpit ONLY if Flutter web proves insufficient — not before.

---

## 1. System Architecture (Target State)

```
┌────────────────────────── INTERFACE LAYER ──────────────────────────┐
│  OneMind Command (Flutter: iOS / Android / Web / Linux)             │
│  Tabs: Home · Chat · Agents · Codex · Command                      │
│  Connects: NATS WebSocket (primary) · OpenClaw WS (legacy fallback) │
└──────────────────────────────┬──────────────────────────────────────┘
                               │ INTERFACE account (exports-only view)
┌──────────────────────────────▼──────────────────────────────────────┐
│                        NATS SPINE (per NATS-DOCTRINE.md)            │
│  Core VPS (omos-core) + leaf nodes per fabric site (house/barn/...)│
│  JetStream: evt.> cmd.> tel.> svc.>   KV: config, schema registry  │
│  Accounts: SYS / CORE / FABRIC / INTERFACE                          │
└───┬──────────────┬───────────────┬───────────────┬─────────────────┘
    │ CORE         │ CORE          │ FABRIC        │ SYS
┌───▼─────────┐ ┌──▼───────────┐ ┌─▼────────────┐ ┌▼────────────────┐
│ AGENT LAYER │ │ META LAYER   │ │ FABRIC LAYER │ │ OBSERVABILITY   │
│ Synadia     │ │ om-meta (Py) │ │ Wombat       │ │ om-insights (Py)│
│ Agent SDK   │ │ tasks        │ │ bridges:     │ │ surveyor +      │
│ v0.3:       │ │ missions     │ │  HomeAssist  │ │ $SYS events →   │
│  Legacy     │ │ knowledge    │ │  ESP32/IoT   │ │ SQLite history  │
│  Claude Code│ │ scheduling   │ │  webhooks    │ │ + auto-checks   │
│  OpenClaw   │ │ SQLite+evt.> │ │ NEX deploys  │ │ ("Insights")    │
│  workers... │ │              │ │              │ │                 │
└─────────────┘ └──────────────┘ └──────────────┘ └─────────────────┘
```

### The OneMind product suite (Synadia-equivalent naming)

| OneMind name | Replicates | What it is | Where it lives |
|---|---|---|---|
| **OM Control** | Synadia Control Plane | Manage/secure/configure the bus: accounts, streams, KV, nodes | Python svc + app **Command** tab |
| **OM Insights** | Synadia Insights | Topology, history ("time travel" = SQLite+SQL), auto-checks, traffic drill-down | Python svc + app **Analytics** page |
| **OM Protect** | Synadia Protect | Connect-time policy via auth callout; Expr-style rules later | Deferred to Phase 7 |
| **OM Connect** | Synadia Connect | Wombat runtime configs bridging HA/IoT/webhooks ↔ bus | Fabric layer |
| **OM Workloads** | Synadia Workloads | NEX nodes on every fleet machine; deploy services over the bus | Fleet-wide |
| **OM Agents** | Synadia Agents | Agent SDK v0.3 — discovery, prompt routing, review/edit config. **The SDK protocol is a first-class citizen INSIDE the app**: the Agents tab speaks v0.3 natively (`$SRV` discovery, `agents.{verb}.*` subjects, hb liveness) — not through any intermediary gateway | Agent layer + app **Agents** tab |

---

## 2. Current State (Verified 2026-06-10)

### Working
- Flutter app ~75% feature-complete: 5-tab shell, streaming chat, dashboard, agent roster,
  Codex browser, settings. ~9,000 lines Dart, 39 files. iOS deploy works (release mode).
- Web build path works (was deployed to VPS via nginx + Tailscale Funnel `/app/`).
- VPS proxy (`/opt/openclaw-proxy/server.js`) bridges OpenClaw WS → REST/SSE for web. 9/12 RPC methods live.
- NATS server on core VPS; Tailscale mesh (`tail717ba0.ts.net`) across fleet.
- Flutter SDK 3.44.1 now installed on omos-mbp (was the #1 blocker — app was built on the old laptop).

### Fixed this session
- `lib/core/config.dart`: hostname `omos-core-1` → `omos-core` (was dialing a nonexistent host — primary suspect for "connection issues").
- `web/index.html` + `web/manifest.json`: proper OneMind branding, dark theme color.

### Known issues / debt
- Git remote points at `onemindos/onemind-app.git`; local repo renamed `onemind-command`. → Rename on GitHub (auto-redirects old URL).
- Voice input + TTS are stubs (iOS 26 crashes forced plugin removal — needs LiveKit approach, not device plugins).
- Push notifications stubbed.
- Gamification: model exists (314 lines), no UI.
- No NATS path in the app yet — everything rides the OpenClaw proxy.
- 25 spec docs in `docs/specs/` are aspirational, unprioritized — this plan supersedes their ordering.

---

## 3. Phased Roadmap

Phases ship independently. Each has a **definition of done** and **delegation packets** (§4).
Order is dependency-driven: bus before consumers of the bus.

### Phase 0 — Baseline Green (NOW)
**Goal:** app builds + runs locally on this Mac; repo synced; docs authoritative.
- 0.1 ✅ Flutter SDK installed (3.44.1 via Homebrew)
- 0.2 ✅ `flutter build web` passes (lucide_icons vendored + patched for Flutter 3.44)
- 0.3 ✅ Served locally at `127.0.0.1:8090/app/`, dashboard verified in browser
- 0.3b ✅ OpenClaw setup gate removed from boot flow (Zeus directive) — app boots straight to dashboard; transports configured in Settings
- 0.4 ✅ Commit fixes + docs; push; GitHub repo renamed → `onemindos-command`
- 0.5 README/AGENTS.md updated to point at this plan (this session)
- **DoD:** Zeus opens `localhost:8090` in Brave, sees live gateway data. `git status` clean, remote matches.

### Phase 1 — NATS Foundation (Week 1)
**Goal:** the bus is real, secured, and observable; nothing app-side yet.
- 1.1 NATS server config: enable WebSocket listener (:8080 wss via Tailscale), JetStream domains for core
- 1.2 Account layout: SYS / CORE / FABRIC / INTERFACE (config-based; operator-ready subjects per Doctrine §5)
- 1.3 Stream topology: `EVT` (evt.>), `TEL` (tel.>), `CMD` (cmd.>) with retention/republish policies per Doctrine §3
- 1.4 KV buckets: `om-config`, `om-schema` (the "schema registry" — JSON Schemas keyed by subject)
- 1.5 Leaf node config template for fabric sites (house/barn/office) — store-and-forward per Doctrine §4
- **DoD:** `nats pub evt.sys.test.fired '{}'` from any mesh machine lands in EVT stream; accounts enforce isolation; leaf template documented.

### Phase 2 — App ↔ NATS (Week 1–2)
**Goal:** Flutter app reads the bus directly. OpenClaw path demoted to fallback.
- 2.1 Add `dart_nats` + single `lib/core/nats_service.dart` (connection, auth, request-many w/ stall timeout per Doctrine §6)
- 2.2 `NatsSyncNotifier` alongside existing `GatewaySyncNotifier` — subscribes `evt.>`, populates same providers
- 2.3 Agent discovery: `$SRV.PING.agents` sweep → live roster (replaces hardcoded `Agent.knownAgents`)
- 2.4 Connection settings UI: NATS URL + creds (secure storage), per-transport status chips
- 2.5 Web build: verify `dart_nats` over WebSocket works in browser (it must — if not, escalate to Legacy)
- **DoD:** App home shows live bus data with OpenClaw proxy stopped. Roster reflects whatever's actually registered on the bus.

### Phase 3 — Agent Layer v0.3 (Week 2)
**Goal:** heterogeneous agents discoverable + promptable from the app; config review/edit per Doctrine §8.
- 3.1 Hermes adapter → Synadia Agent SDK v0.3 (verb-first subjects, `$SRV` registration, hb)
- 3.2 Claude Code plugin (`claude-channel-nats`) wired to the bus
- 3.3 App Agents tab speaks Agent SDK v0.3 natively (in-app protocol implementation in Dart): roster from `$SRV.PING.agents` discovery, status from hb, prompt → `agents.prompt.{agent}.{owner}.{session}` streaming
- 3.4 Config review/edit: read harness config via agent endpoints, edit, write back natively. NO app-side storage.
- **DoD:** From the app: see Legacy + Claude Code live, prompt either, view+edit a config value that round-trips to the harness's own file.

### Phase 4 — Meta Layer: om-meta (Week 3–4)
**Goal:** tasks/missions/knowledge/scheduling as a bus-native Python service. Replaces Todoist-for-OS-work + scattered cron state.
- 4.1 `om-meta` service skeleton: nats-py + SQLite (WAL), svc.meta.* request handlers, evt.so.* emissions
- 4.2 Task + mission schema (steal SwarmClaw's data model shapes — reviewed earlier, MIT)
- 4.3 Scheduling: JetStream delayed messages for time-based triggers
- 4.4 App Tasks page binds to svc.meta.* + live evt.so.task.> updates
- **DoD:** Create task in app → SQLite row + evt.so.task.created on bus → visible in second client in <1s.

### Phase 5 — OM Insights (Week 4–5)
**Goal:** Synadia-Insights-equivalent observability.
- 5.1 Deploy `nats-surveyor` + enable $SYS account event capture → `om-insights` Python collector → SQLite history
- 5.2 Auto-checks engine: start with 10 rules (slow consumers, stream quorum, offline replica, high RTT, disk %, leaf disconnect, consumer lag, auth failures, cert expiry, heartbeat miss) — emit `evt.sys.check.failed`
- 5.3 App Analytics page: topology view, time-range queries ("time travel"), check status board
- **DoD:** Kill a leaf node → check fires → event on bus → red card in app within 30s.

### Phase 6 — Fabric Bridges: OM Connect (Week 5–6)
**Goal:** physical world on the bus.
- 6.1 Wombat runtime on core: Home Assistant ↔ `evt.le.>` / `cmd.le.>` bridge
- 6.2 ESP32/sensor ingestion path → `tel.>` (WireClaw pattern when hardware arrives)
- 6.3 NEX nodes on omos-mbp + core VPS (`nex node up`) — deploy om-meta/om-insights as NEX workloads
- **DoD:** Flip a Hue light from app Command tab via cmd.le.light.*; sensor telemetry visible in Analytics.

### Phase 7 — Hardening: OM Protect (Later — explicit trigger, not a date)
**Trigger:** fleet >20 nodes OR first non-family device class OR first external exposure beyond Tailscale.
- 7.1 Operator/JWT cutover (subjects already designed for it)
- 7.2 Auth callout service (callout.go) — central policy at connect-time
- 7.3 Per-message policy proxy (Expr rules) only if a real threat model demands it
- Also here: voice (LiveKit), push notifications, gamification UI — feature backlog resumes after foundations.

---

## 4. Delegation Protocol (How Work Gets Assigned)

**Roles**
- **Legacy (this profile):** architecture, task specs, review, merge authority. Does NOT write feature code.
- **Worker agents (Sonnet 4.6):** execute packets below. One packet = one branch = one PR.
- **Opus 4.8:** packets flagged 🔴 HARD only.

**Rules for workers**
1. Read `MASTER-PLAN.md`, `NATS-DOCTRINE.md`, `AGENTS.md` first. Doctrine violations = automatic reject.
2. Branch naming: `phase{N}/{packet-id}-short-name`. Conventional commits.
3. Every packet PR includes: what changed, how verified (real command output, not claims), doctrine sections touched.
4. Never modify harness configs by hand — that's the app's job via the SDK (Doctrine §8).
5. Blocked? Stop and report. Do not improvise around a blocker.
6. SQLite always WAL mode. Subjects always from Doctrine §2 grammar. No new dependencies without packet approval.

**Packet index** (specs written per-phase by Legacy as phases open)

| Packet | Phase | Scope | Difficulty |
|--------|-------|-------|------------|
| P0-A | 0 | Local web serve + verification script | 🟢 |
| P1-A | 1 | NATS server config: WS listener, accounts, streams, KV | 🟡 |
| P1-B | 1 | Leaf node template + store-and-forward validation | 🟡 |
| P2-A | 2 | `nats_service.dart` + NatsSyncNotifier | 🟡 |
| P2-B | 2 | Discovery sweep + roster wiring | 🟡 |
| P2-C | 2 | Web-platform NATS WS verification | 🟢 |
| P3-A | 3 | Hermes v0.3 adapter | 🔴 HARD |
| P3-B | 3 | Claude Code plugin wiring | 🟡 |
| P3-C | 3 | Agents tab: prompt + config round-trip UI | 🟡 |
| P4-A | 4 | om-meta service core | 🟡 |
| P4-B | 4 | Task/mission schema + handlers | 🟡 |
| P4-C | 4 | App Tasks page binding | 🟢 |
| P5-A | 5 | Insights collector ($SYS + surveyor → SQLite) | 🟡 |
| P5-B | 5 | Checks engine (10 rules) | 🟡 |
| P5-C | 5 | Analytics page: topology + time travel | 🟡 |
| P6-A | 6 | Wombat HA bridge | 🟡 |
| P6-B | 6 | NEX rollout + workload deploys | 🔴 HARD |

---

## 5. Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-06-10 | Flutter stays; no rewrite | App is 75% done; issues were connection config, not framework. Cross-platform requirement met by one codebase |
| 2026-06-10 | App connects to NATS directly; OpenClaw demoted to fallback | Kills protocol-mismatch fragility; bus is the stable contract |
| 2026-06-10 | Harness owns agent config; app = review/edit only | Zeus directive. No shadow config, no sync wars (Doctrine §8) |
| 2026-06-10 | Replicate Synadia closed products on open primitives | Their moats are UIs over public APIs; sovereignty requires self-hosted |
| 2026-06-10 | InsForge/Convex rejected | Cloud-coupled, Postgres-centric; our primitives (NATS+SQLite) already cover the need offline |
| 2026-06-10 | AionUI rejected; SwarmClaw/Mission Control = idea quarries only | Desktop toy / no NATS respectively. Steal schemas + panel layouts (MIT), run neither |
| 2026-06-10 | React cockpit deferred until Flutter web proves insufficient | One codebase discipline; avoid premature second frontend |
| 2026-06-10 | Config-based NATS auth now, operator-ready design, JWT cutover at >20 nodes | Right-size security ceremony to current fleet (Doctrine §5) |
| 2026-06-10 | No setup/onboarding gate in app boot | OpenClaw-era relic. Bus-first app renders dashboard always; transports live in Settings |
| 2026-06-10 | GitHub repo = `onemindos/onemindos-command` | Org-prefixed naming per Zeus |
| 2026-06-10 | Synadia OSS locked in: Wombat (bridges), NEX (workloads), orbit.py/js (clients), Agent SDK (agent layer IN-APP), synack (pattern ref only) | Use open engines as-is; build only UIs + glue |

---

## 6. Worker Quick-Start

```bash
# Environment
cd ~/dev/onemind-command
export PATH="/opt/homebrew/bin:$PATH"   # flutter + dart live here
flutter --version                        # expect 3.44.x

# Build web
flutter pub get && flutter build web --release --base-href /app/

# Run local dev (Chrome/web-server)
flutter run -d web-server --web-port 8090

# Read before coding (in order)
# 1. MASTER-PLAN.md  (this file — your packet's phase section)
# 2. NATS-DOCTRINE.md (binding rules)
# 3. AGENTS.md        (repo conventions, layout, gateway details)
```

---

*Maintained by Legacy. Workers: if reality contradicts this document, STOP and report — do not silently adapt.*
