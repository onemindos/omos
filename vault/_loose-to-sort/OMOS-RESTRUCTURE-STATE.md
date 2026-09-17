# OMOS Restructure — Current State & Source of Truth

> **Purpose:** Living doc for the OMOS restructure. There was no doc for the in-flight
> restructure; this is it. Update it as things change.
> **Owner:** Zeus. **Maintained with:** Letta (Letta Code agent).
> **Created:** 2026-07-01. **Last verified live:** 2026-07-01 ~03:45 EDT on `omos-mbp`.
> **Pairs with canon:** `omos-core/docs/ARCHITECTURE.md` (2026-06-13), `NATS-DOCTRINE.md`, `MASTER-PLAN.md`.

---

## 0. The one-liner
A local-first OS where humans, AI agents, computers, and machines are all first-class
**actors** on a single **NATS** nervous system. Every event is recorded forever in
**ClickHouse** (the Time Machine). Work flows through a universal **Task Engine**. **One app
(`omos-atoc`) renders and controls all of it.** The bus owns truth; the app is a window.

---

## 1. The two halves (dependency is one-way: frontend → backend)
| | Repo | Is | Language |
|---|---|---|---|
| **FRONTEND** | `omos-atoc` | Flutter app "OneMind Command" (iOS·Android·Web·Desktop). Single pane of glass, UI-state only. | Dart/Flutter |
| **BACKEND** | `omos-core` | The OS: NATS config, `omd` daemon, `services/`, schemas, installer, canon docs. | Go/Python + config |

Goal of the restructure: **condense the many local dashboards into the ONE frontend UI**,
fed by live data, **reachable from mobile**. Chosen exposure: **Tailnet-only** (private to
Zeus's devices; ClickHouse/NATS never hit the public internet).

---

## 2. Doctrine (non-negotiable)
1. Bus-first, app-stateless — app renders bus state, owns no DB.
2. **Push, never poll** — no polling loops; NATS delivers in real time.
3. Location transparency — same subjects localhost or remote (Tailscale).
4. Polymorphic actors — human | agent | node | device | squad, one identity model.
5. Record everything → ClickHouse.
6. App connects to **NATS directly over WebSocket (:8080)**, subscribes to subjects, and
   issues commands by publishing.

### Five layers
- **L1 Substrate:** NATS — :4222 client, :8080 ws, :1883 mqtt, :7422 leaf, :8222 monitor.
- **L2 Memory:** JetStream (hot/durable + KV) + ClickHouse (cold/analytical — Time Machine).
- **L3 Compute:** `omd` daemon per node (fs/shell, runs agent harnesses) + agent layer.
- **L4 Task Engine:** Issues (durable) + Tasks (runs) + Actors + Triggers.
- **L5 Interface:** the Flutter app (`omos-atoc`).

---

## 3. THIS MACHINE = `omos-mbp` (the core host today)
- Tailscale name/host: `omos-mbp` → `omos-mbp.tail717ba0.ts.net` (100.74.182.117).
- **CRITICAL: `omos-core` is NOT a tailnet host — it does not resolve.** Any app config
  pointing at `omos-core.tail717ba0.ts.net` (gateway / NATS / ClickHouse) is DEAD.
  All those services actually run right here on `omos-mbp`.
- Tailnet peers: omos-cloud (linux, online), omos-ha (online), omos-iphone17 (online),
  omos-ipad (offline), omos-dell/pad/tab/tablet/tv (offline).

---

## 4. Verified port map on omos-mbp (2026-07-01)
| Service | Port | State | Notes |
|---|---|---|---|
| NATS client | 4222 | up | |
| NATS websocket | 8080 | up | browser/app clients connect here (doctrine path) |
| NATS monitor | 8222 | up | `/healthz` OK after .DS_Store fix (was 503) |
| ClickHouse HTTP | **8123** | up | v26.6.1. **App wrongly configured for 8124.** |
| ClickHouse native | 9001 | up | |
| Gateway (node) | 18789 | up | Tailscale serve :8443 (tailnet-only) → this |
| Hermes webui (python) | 8787 | up | tailnet-only via serve :9443 |
| Public funnel target | 8091 | **DEAD** | funnel 443 → 8091 points at nothing |
| omos-atoc web (Caddy) | 9100 | up | LaunchAgent `com.onemind.atoc`, root `omos-atoc/build/web` |
| agenttap HUD | 8088 | up | reference UI the theme was lifted from |

### Tailscale serve/funnel today
- `https://omos-mbp.tail717ba0.ts.net` (Funnel ON, public) → 127.0.0.1:8091 **(dead)**
- `:8443` (tailnet-only) → localhost:18789 (gateway)
- `:9443` (tailnet-only) → localhost:8787 (hermes webui)

---

## 5. Backend services — `omos-core/nats/services/`
NATS `$SRV` micro-services. Publish `evt.{domain}.>`, accept `cmd.{domain}.>`, state in JetStream KV.

| Service | Status | Notes |
|---|---|---|
| `task-engine/` | ✅ designed LIVE, **currently DOWN** | `om_task.py` (`om-task` v2.0.0). LaunchAgent `com.onemind.om-task`. |
| `connect/` | ✅ LIVE | `om_presence.py` voice bridge (running as `com.onemind.om-presence-v2`). |
| `control-plane/` | 🔜 stub | accounts/streams/KV/nodes (Synadia Control Plane equiv). |
| `insights/` | 🔜 stub | topology/history/checks (Synadia Insights equiv). |
| `agents/ api/ auth/ bridge/ data/ edge/ grid/ intel/ meta/` | 📄 design | |

### Task Engine subject contract (om-task v2.0.0)
- Request-reply: `task.create`, `task.update.*`, `task.approve.*`, `task.list`, `task.get.*`,
  `task.calendar`, `task.inbox`, `task.project.create`, `task.project.list`. Subscribes `agents.hb.>`.
- Events: `evt.task.created / updated / due.soon / {id}.overdue / {id}.stalled / {id}.complete / {id}.approved`.
- KV: `task-index`, `projects`, `task-notify`(48h), `agent-registry`.
- Push: ntfy.sh topic `onemind-zeus-64c8cb63c3ee91bc`; Telegram `@omos_hermes_bot` chat 577573558.

### REST bridges (translate REST↔NATS for the Flutter app)
- **task-bridge** `services/task-engine/task-bridge/bridge.js` — **:8401**. REST + SSE `/events`.
  Feeds `/tasks/inbox`, `/tasks/calendar`, `/tasks/board`, `/tasks/:id`, `/notifications`.
- **session-bridge** `services/session-bridge/bridge.js` — **:8402**. LaunchAgent `io.onemind.session-bridge`.
  Read-only; unifies `~/.hermes/state.db` (SQLite) + `~/.claude/projects/**/*.jsonl`.
  Feeds `/agents/sessions`, `/agents/sessions/:id`.
- ⚠️ `omos-core/nats/omos-bridge/BRIDGE-RUNNER.md` is an OLDER retired HA/Hermes bridge — NOT these.

### ⚠️ 5b. DEEPER RECON (2026-07-01, later) — restructure has moved past the 06-13 doc
The actual `services/` tree today does NOT match the ARCHITECTURE.md §4/§4b bridge paths:
- **No `services/session-bridge/` dir. No `services/task-engine/task-bridge/` subdir.** The REST
  bridges for :8401 / :8402 described in canon **do not exist in the current tree.** Either the
  restructure will feed Tasks/Sessions over NATS directly (doctrine: push, NATS-ws), or the
  bridges must be rebuilt. **OPEN DECISION.**
- Current `services/` dirs (dated Jun 21–30, more built-out than README's "design" label):
  `agents/` (LEGACY.md, README, empty `extracted/`), `api/`, `auth/`, `bridge/`, `connect/`,
  `data/`, `edge/`, `grid/`, `intel/`, `meta/`, `task-engine/`, plus empty `control-plane/`, `insights/`.
- Most non-live services contain an **empty `extracted/`** dir → scaffolding, not implemented.
- **Live-service pattern:** the service keeps its own LaunchAgent plist IN its dir. Example:
  `connect/com.onemind.om-presence-v2.plist` runs `connect/om_presence_v2.py`
  (env `NATS_URL=nats://127.0.0.1:4222`). This is the template for standing up task-engine.
- `task-engine/om_task.py` EXISTS (38KB, Jun 13) but has **no installed LaunchAgent** → down.
- `meta/om_meta.py` + `meta/healthcheck.sh` exist (a meta-agent service).
- **Authoritative NATS config = `/Users/zeus/nats-server.conf`** (LaunchAgent `com.nats.server`,
  `nats-server -c /Users/zeus/nats-server.conf`). No auth/TLS block. JetStream store `~/nats-data`.
- Bus is alive (monitor: 22 conns, 74M msgs). `nats` CLI (v0.4.0, ctx `omos`→nats://localhost:4222)
  currently errors `no servers available` even via 127.0.0.1 — local CLI/handshake quirk to debug;
  NOT a server outage (services connect fine). TODO: fix CLI context so we can `nats stream ls`.

---

## 6. Running LaunchAgents (omos/onemind) — 2026-07-01
**Running:** `omos.heynats`, `omos.nats-dashboard`, `omos.clickhouse.timemachine`,
`omos.agenttap` (:8088), `com.onemind.atoc` (:9100), `com.onemind.om-presence-v2`,
`com.onemind.nex-node`, `omos.agent-web-ui`, `com.onemind.natsdata-guard` (new — see §8).
**Down (gaps):** `com.onemind.om-task`, task-bridge (:8401), `io.onemind.session-bridge` (:8402).
**Disabled:** `omos.orp`.

---

## 7. THE DATA GAPS (why frontend pages are blank)
1. **Dead host:** config points at `omos-core.*` which doesn't resolve → gateway/NATS/ClickHouse unreachable.
2. **Wrong ClickHouse port:** app says 8124, real is 8123 → Insights / Time Machine blank.
3. **Bridges/engine down:** task-engine + :8401 + :8402 → Tasks + Sessions blank.
4. **Web /api gap:** `/api/events` (SSE) + `/api/all` not proxied → gateway sync dead on web.
5. NATS-backed pages (Center, Wire, Observe, Manage, Secure, Topology) work once pointed at omos-mbp.

### Frontend page → data source (from nav_model.dart `kNav`, flags live/ready)
- COMMAND: Center✓live, Wire✓live, Meta(stub)
- BUS: Observe✓, Manage✓, Secure✓, Topology✓ (NATS monitor); Connect(stub)
- AGENTS: Roster, Sessions (session-bridge), Memory, Skills, Presence✓live
- FLEET: Map, Robotics(stub), Drones(stub), Edge(stub)
- OPS: Services✓, Insights✓(ClickHouse), Alerts(stub), Task Inbox/Calendar/Board✓(task-bridge), Notifications✓
- SYSTEM: Settings, Page Map

---

## 8. Decisions & changes made this session (2026-07-01)
- **★ TRANSPORT = PATH A: NATS-NATIVE (decided by Zeus 2026-07-01).** The frontend connects
  DIRECTLY to NATS over WebSocket (:8080) — subscribes to subjects + does request-reply. **NO
  REST bridges, no gateway REST/SSE.** This is the doctrine path (push, never poll). The dead
  task-bridge(:8401)/session-bridge(:8402) are NOT rebuilt. Tasks use `task.*` request-reply +
  `evt.task.*` subs; Sessions/agents/insights all move to NATS subjects. Retire GatewaySyncNotifier
  (REST/SSE) and the gateway(:18789) dependency. App keeps `dart_nats` (already in pubspec).
- **Exposure = Tailnet-only** (not public funnel). ClickHouse/NATS stay private; phone/iPad reach via tailnet.
- **Fixed NATS false FAULT:** removed stray `.DS_Store` from `~/nats-data/jetstream/` that NATS
  read as a JetStream account (`/healthz` was 503). Now `{"status":"ok"}`.
- **Added guard:** LaunchAgent `com.onemind.natsdata-guard` sweeps `.DS_Store` from `~/nats-data` every 300s.
- **Confirmed:** minimal Flutter build renders white only in headless (no GPU/WebGL) browser; real
  browsers with WebGL render fine. Not an app-code bug.

---

## 9. Consolidation architecture (target)
One Caddy on omos-mbp path-routes everything; Tailscale serve (tailnet-only) fronts it:
```
https://omos-mbp.tail717ba0.ts.net/        → OneMind Flutter dashboard (build/web)
  /natsmon/*      → 127.0.0.1:8222   (NATS monitor)
  /clickhouse/*   → 127.0.0.1:8123   (ClickHouse HTTP — real port)
  /gateway  (ws)  → 127.0.0.1:18789  (agent gateway)  [or NATS-ws :8080 per doctrine]
  /task-api/*     → 127.0.0.1:8401   (task-bridge)
  /session-api/*  → 127.0.0.1:8402   (session-bridge)
  /api/events,/api/all → sync proxy
```
Then repoint `omos-atoc/lib/core/config.dart` off dead `omos-core` → this base.
**Doctrine tension to resolve:** app should read NATS-ws (:8080) directly (push), but current
build uses gateway (:18789) + REST bridges. Decide the canonical path during restructure.

---

## 10. PATH A migration plan (NATS-native frontend)

### Good news: the transport layer already exists
`omos-atoc/lib/core/nats_service.dart` — singleton `NatsService.instance` — already implements
`subscribe()`, `subscribeMultiple()`, `publish()/publishJson()`, and `request()` (request-reply)
over NATS WebSocket via `dart_nats`. Primary URI `ws://127.0.0.1:8080`, fallback
`ws://omos-core...:8080` (DEAD — fix to omos-mbp). So calling Task Engine subjects is a one-liner.

### REST → NATS subject mapping (task_notifier.dart, representative)
| Current Dio call (REST bridge) | Target NATS request-reply |
|---|---|
| GET `/tasks/inbox` | `request('task.inbox', {})` |
| GET `/tasks?status&domain&project_id` | `request('task.list', {filters})` |
| GET `/tasks/calendar?from&to` | `request('task.calendar', {from_date,to_date})` |
| GET `/tasks/:id` | `request('task.get.{id}', {})` |
| POST `/tasks` | `request('task.create', {body})` |
| PATCH `/tasks/:id` | `request('task.update.{id}', {fields})` |
| DELETE `/tasks/:id` | `request('task.update.{id}', {status:'archived'})` — no delete in contract |
Response bodies already match (bridge was a thin wrapper). Same pattern for session_* and agent notifiers.

### Files to rewire (retire Dio/gateway, use NatsService)
- `stores/task_notifier.dart` → `task.*` request-reply.
- `stores/session_list_notifier.dart`, `stores/session_notifier.dart` → session subjects (NEEDS backend svc, see below).
- `stores/agent_list_notifier.dart` → agents discovery (agent_roster_service.dart already uses dart_nats + `agents.hb.>`).
- RETIRE: `stores/gateway_sync_notifier.dart`, `core/gateway_client.dart`, gateway(:18789) dep.
- `core/config.dart`: drop taskBridgeUrl/sessionBridgeUrl/gateway; fix ClickHouse 8124→8123; NATS URIs → omos-mbp.

### Backend prerequisites (a page has no data until its NATS service runs)
- **Tasks:** `task-engine/om_task.py` EXISTS but is DOWN → install LaunchAgent (copy `connect/` plist pattern). Serves `task.*`.
- **Sessions:** NO NATS service exists (old REST session-bridge deleted) → **BUILD a session `$SRV`**
  exposing Hermes (`~/.hermes/state.db`) + Claude (`~/.claude/projects/**/*.jsonl`) over e.g. `session.list`/`session.get`.
- **Insights (ClickHouse):** SQL-over-HTTP, not NATS. Either query ClickHouse :8123 directly (native/tailnet)
  OR build an `insights` `$SRV` that wraps queries (`svc.insights.query`). OPEN — likely direct-HTTP short term.
- **NATS reachability from mobile:** expose ws :8080 via Tailscale serve (tailnet-only). Update NatsServiceConfig
  primary=`ws://omos-mbp.tail717ba0.ts.net:8080` (or MagicDNS `ws://omos-mbp:8080`), drop dead omos-core fallback.

### Ordered execution
1. [ ] **Phase 1** — Stand up `om-task` as NATS service (LaunchAgent); verify `task.inbox` answers over NATS.
2. [ ] **Phase 2** — Expose NATS ws :8080 via Tailscale serve (tailnet-only); fix NatsServiceConfig + config.dart URIs.
3. [ ] **Phase 3** — Rewire `task_notifier` → NATS request-reply; rebuild; verify Tasks pages live end-to-end.
4. [ ] **Phase 4** — Build session `$SRV`; rewire session notifiers; verify Sessions pages.
5. [ ] **Phase 5** — Roster via `agents.hb.>`; Insights via ClickHouse :8123 (or insights svc); retire GatewaySyncNotifier.
6. [ ] **Phase 6** — Full sweep: rebuild, screenshot-verify every page, update this doc's status table.

## 11b. ★ CRITICAL FINDING — TWO Task Engines collide on the bus (2026-07-01)
There are **two `om-task` implementations** both claiming `task.*` subjects on NATS:
- **Go `om-task` (RUNNING via NEX)** — `omos-core/nats/nex/om-task/main.go` (327 lines). Deployed
  via NEX (`Nexfile`). Implements ONLY: `task.create`, `task.list`, `task.update`, `task.syntax.validate`.
  **MISSING:** `task.inbox`, `task.calendar`, `task.get.*`, `task.approve.*`, `task.project.create/list`.
- **Python `om_task.py` (NOT installed)** — `omos-core/nats/services/task-engine/om_task.py` (915 lines).
  COMPLETE: all endpoints + heartbeat + ntfy/Telegram push. Needs `aiohttp` (installed today) + LaunchAgent.
- **Collision:** both share service name `om-task` + queue group → `task.list` load-balances between them
  nondeterministically; 2-token subjects time out. Verified via `/connz?subs=1` (cid 113 Go, cid 127 py).
- **Consequence:** Task Inbox + Calendar pages CANNOT work on the Go service (no `inbox`/`calendar` responder).
- **DECISION NEEDED (see below).** Do NOT run both simultaneously.

## 11. Done
- [x] Fix NATS `.DS_Store` health false-alarm + add `com.onemind.natsdata-guard`.
- [x] Decide transport = PATH A (NATS-native).
- [x] Create this doc on Desktop.
- [x] Task Engine deps: installed `aiohttp` into `/opt/homebrew/bin/python3`.
- [x] Verified Python `om_task.py` runs + registers + answers `task.*` request-reply (when alone on bus).
- [x] Diagnosed dual-task-engine collision (Go NEX vs Python). See §11b.
- [x] **DECISION (Zeus 2026-07-01): GO is canonical.** Finish the Go `om-task` NEX port; retire Python `om_task.py`.
  Rationale: single static binary for the fleet (Pi/drones/edge), NEX-native deploy, no Python dep hell, Go ecosystem.
  Frontend is unaffected either way (same NATS `task.*` contract).
- [x] Ported ALL missing Go endpoints (inbox, calendar, get, update, approve, delete, project.create/list)
  in `nex/om-task/` (split across main.go, run.go, task_write.go, task_read.go, project.go). Uses flexible
  `map[string]interface{}` schema matching frontend `TaskItem` exactly (due_date, is_all_day, labels, etc.).
  Builds clean. Verified live: create/get/update/approve/delete/project.*/list all ✓.
- [x] Fixed Nexfile URI (nats/nex → nex path move); bumped to v3.1.0-nex.
- [x] **DATA BLOAT PURGED (user approved).** Deleted JetStream streams `KV_task-index` (441,648 keys /
  342 MB) and `KV_tasks` (223,113 keys / 111 MB). Engine auto-recreated a clean `task-index` on deploy.
  Empty responses now return `[]` not `null` (Dart-safe) — patched list/inbox/calendar/project.list.
- [x] **Go om-task LIVE via NEX** (workload `zsUBkXt5IOGVxMi5nh5cJ8`, v3.1.0-nex). All 11 subjects verified:
  create/list/get/update/approve/delete/calendar/inbox/project.create/project.list/syntax.validate — all
  instant, contract-complete (due_date, is_all_day, labels, progress, project_id, recurrence, reminders...).
  Python `om_task.py` retired (not running). Dual-engine collision RESOLVED.
- [ ] WATCH: whatever created 220K tasks — currently no active spammer (bucket stayed at test count).
  If task-index regrows, hunt the producer (old engine bug or an ingestion loop).

### NEXT PHASE — Frontend NATS wiring (omos-atoc)
- [ ] Rewrite `task_notifier.dart` to call `NatsService.instance.request('task.inbox'|'task.list'|'task.calendar', {})`
  instead of Dio/HTTP bridges. Same for session/agent notifiers.
- [ ] Fix `core/config.dart`: ClickHouse 8124→8123, NATS URIs → omos-mbp, drop taskBridgeUrl/sessionBridgeUrl.
- [ ] Verify every page renders live NATS data over Tailscale (:8080 ws).
