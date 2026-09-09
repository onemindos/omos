# ARCHITECTURE.md — OneMind OS & OneMind Command

> **Authority:** Binding architecture for the OneMind application and the OS beneath it.
> Pairs with `NATS-DOCTRINE.md` (bus rules) and `MASTER-PLAN.md` (build order).
> Maintained by: Legacy (architect). Last updated: 2026-06-13.

---

## 0. The One-Sentence Architecture

**A local-first operating system where humans, AI agents, computers, and machines are all
first-class *actors* on a single NATS nervous system; every event is recorded forever in
ClickHouse (the Time Machine); work flows between actors through a universal Task Engine;
and one lightweight app renders and controls all of it.**

Nothing in the app owns state. The bus owns truth. The app is a window.

---

## 1. Design Principles (the anti-bottleneck rules)

These exist so features bolt on for years without a rewrite. Every one prevents a specific
bottleneck we would otherwise hit.

| # | Principle | Bottleneck it prevents |
|---|-----------|------------------------|
| 1 | **Bus-first, app-stateless.** App renders bus state; it never owns a private database to migrate. | App-owns-state migration hell |
| 2 | **Push, never poll.** NATS delivers in real time. No 3-second polling loops anywhere. | Polling latency + load (Multica's compromise) |
| 3 | **Location transparency.** Same subjects whether a node is localhost or in the barn 200m away. | Local-vs-remote special-casing |
| 4 | **Polymorphic actors.** Human / agent / node / device / squad share one identity model + capability sets. | Hardcoding "human vs AI"; adding "drone" later = a rewrite |
| 5 | **Record everything from day one.** ClickHouse captures all semantic events permanently. | "We wish we'd logged that" |
| 6 | **Harness sovereignty.** App reviews/edits agent configs through each harness's own endpoints; never a shadow copy. | Config fights, sync conflicts |
| 7 | **Feature-first code layout.** Each capability is a self-contained module, not a layer cut across files. | Monolith surgery for every feature |
| 8 | **Versioned event envelope + schema registry.** Every message carries a typed, versioned envelope. | Schema churn breaking old consumers |
| 9 | **One download installs the OS.** Installer ships NATS + ClickHouse + daemon + app, pre-wired and supervised. | "Works on my machine" / setup friction |

---

## 2. The Five Layers

```
┌──────────────────────────────────────────────────────────────────────┐
│  L5  INTERFACE      OneMind Command (Flutter: iOS·Android·Web·Desktop) │
│                     + React cockpit (ops). Renders & controls. No state.│
├──────────────────────────────────────────────────────────────────────┤
│  L4  TASK ENGINE    Issues (durable) · Tasks (runs) · Actors · Triggers │
│                     OneMind's Todoist-for-everything. JetStream + KV.    │
├──────────────────────────────────────────────────────────────────────┤
│  L3  COMPUTE        omd — the OneMind daemon on every node.            │
│                     Local file access, runs agent CLIs, executes tasks. │
│                     + Agent layer (Synadia Agent SDK v0.3 harnesses).   │
├──────────────────────────────────────────────────────────────────────┤
│  L2  MEMORY         JetStream (hot/operational replay, KV, object)      │
│                     ClickHouse (cold/analytical — the TIME MACHINE)     │
├──────────────────────────────────────────────────────────────────────┤
│  L1  SUBSTRATE      NATS — the nervous system. Every byte flows here.   │
└──────────────────────────────────────────────────────────────────────┘
```

Each layer talks only to the bus (L1), never directly to another layer's internals. That
is what lets us replace or scale any single layer without touching the others.

### L1 — Substrate (NATS)
The single transport. Pub/sub, request/reply, streams, KV, object store, service discovery —
one binary from VPS to Pi. Subjects defined in `NATS-DOCTRINE.md §2`. Status: **running**
(onemind-mac, :4222 client, :8080 ws, :1883 mqtt, :7422 leaf). Account split (CORE/FABRIC/
INTERFACE) is card **P1-A**, in progress.

### L2 — Memory (two stores, two jobs)
**Do not confuse these. They both hold events; they serve different masters.**

| | JetStream | ClickHouse |
|---|-----------|-----------|
| Role | Operational durability | Analytical history (**Time Machine**) |
| Retention | Hours→30d (per stream) | Years (10y TTL today) |
| Access | Replay by consumers, exactly-once, KV reads | SQL: aggregate, filter, time-travel |
| Feeds | Producers publish to subjects | NATS-engine table subscribes to subjects |
| Status | Streams = card P1-A | **LIVE** — `onemind.events` 31K rows, `onemind.metrics` 34M rows |

ClickHouse is already wired: a `NATS` engine table subscribes to the bus, a materialized
view fans it into a partitioned `MergeTree`. Every semantic event (`evt.>`, `cmd.>`,
`tel.>`, `agents.>`) lands here forever. This is the substrate for the **Insights/Time
Machine** product (replay "what happened in the barn on March 3rd at 4pm").

> **Migration note (card P2-T):** the live NATS-engine table currently binds the legacy
> subject `events.>`. Rebind to the doctrine semantic set (`evt.>, cmd.>, tel.>, agents.>`)
> and backfill. Until then both subject spaces coexist; producers should dual-publish or we
> alias during cutover. Do **not** drop the existing 31K rows.

### L3 — Compute (the daemon `omd` + the agent layer)
**This is the local-access layer — the answer to "access files on the local computer."**

`omd` (OneMind daemon) runs on every node (Mac, barn Pi, office box, eventually robots).
On boot it:
1. Connects to NATS (local socket, or Tailscale to core).
2. Registers the node as an **actor** and advertises its **capabilities**
   (`fs.read`, `fs.write`, `shell.exec`, `agent.run:hermes`, `gpu`, `camera`, …).
3. Subscribes to its command subject `cmd.node.{node-id}.>`.
4. Discovers and supervises local agent harnesses via the **Synadia Agent SDK** plugin
   (Hermes/Legacy, Claude Code, Codex, …) — exactly the "common structure" where agents
   interact with the computer locally.

When the app or an agent needs something local — read a file, run a command, drive a CLI —
it does **not** reach into the machine directly. It publishes a request to that node's
command subject; the node's `omd` performs the action with real local privilege and replies
over the bus. **Location-transparent, capability-gated, fully recorded.**

Web builds run in a browser sandbox (no raw FS). That's fine: the app talks to the daemon on
`localhost` over NATS-ws. Desktop builds can bundle `omd`. Either way the access path is
identical — through the daemon, never around it.

### L4 — Task Engine (OneMind's universal collaboration primitive)
See §4. The thing that makes humans + AI + computers + drones collaborate. Built on
JetStream (durable issues/tasks) + KV (current state) + the event bus (triggers & updates).

### L5 — Interface (OneMind Command)
The Flutter app (and a future React cockpit for dense ops views). Connects to NATS directly
over WebSocket. Subscribes to the subjects it cares about and renders them; issues commands
by publishing. Holds only UI state. Five tabs today (Home, Chat, Agents, Codex, Command),
growing into Tasks, Time Machine, and Fleet (§7).

---

## 3. The Actor Model (Zeus's extension of Multica)

Multica proved "agents as first-class members." OneMind generalizes the member into a
**polymorphic Actor** so the *physical world* joins the same model.

```
Actor {
  id            string          # stable, bus identity
  type          enum            # human | agent | node | device | squad
  display       { name, icon, color }
  capabilities  [string]        # what it can do (fs.read, shell.exec, fly, weld, notify…)
  endpoint      subject         # where to reach it (cmd.{type}.{id}.>)
  executor      enum            # how work reaches it: NOTIFY (humans) | TRIGGER (machines)
  status        enum            # online | offline | busy | degraded
  meta          json            # type-specific (model, hardware, location…)
}
```

| Type | Examples | Executor | Reached via |
|------|----------|----------|-------------|
| `human` | Zeus, Cielo | **NOTIFY** (inbox/push) | reminder; responds when ready |
| `agent` | Legacy, Claude Code | **TRIGGER** | `cmd.agent.{id}.task` (Agent SDK) |
| `node` | omos-mbp, barn-pi | **TRIGGER** | `cmd.node.{id}.task` (omd) |
| `device` | Go2 dog, drone, CNC | **TRIGGER** | `cmd.device.{id}.task` (omd or bridge) |
| `squad` | "FieldOps team" | **TRIGGER leader** | leader triages, delegates to members |

The defining Multica insight we keep: **NOTIFY vs TRIGGER.** Humans are *notified* and choose
to act; machines are *triggered* and execute. Same task, same board, different executor
binding. Adding a new machine class later is a new `type` + capabilities — never a schema
rewrite. (Principle #4.)

---

## 4. The OneMind Task Engine

> Todoist is human task management. Multica is human+AI code task management. **OneMind Task
> Engine is human + AI + computer + machine task management for an entire homestead OS.**

### Two objects (the Multica split — we adopt it wholesale)

- **Issue** — the durable unit of work. "Re-roof the barn section," "ship app v2," "water
  the greenhouse." Survives forever, human-readable id (`OM-123`), assignee = any Actor,
  status, priority, parent/children, comments, attachments. Stored in JetStream + KV.
- **Task** — one *execution run* against an issue. Ephemeral, its own state machine
  (`queued→dispatched→running→completed|failed|cancelled`), timeouts, retries, session
  resumption. One issue → many tasks over its life. Routed through the `CMD` workqueue stream.

This split is the single most important data-model decision. It cleanly separates "the thing
to do" (durable, human) from "an attempt to do it" (ephemeral, machine), and it's why retries,
multi-agent reruns, and audit all work without contortions.

### Subjects (per doctrine grammar)

```
evt.task.issue.created          evt.task.{id}.assigned        # facts → recorded in ClickHouse
evt.task.{id}.status            evt.task.{id}.comment
cmd.{actor-type}.{actor-id}.task.execute                      # dispatch to an actor
svc.task.issue.list / .get / .create                          # app↔engine RPC
```

### Triggers (how work starts — from Multica, generalized)

1. **Assign** an issue to an actor → `cmd.{actor}.task.execute` (heavy: actor owns it).
2. **Mention** an actor in a comment → light trigger, no ownership change.
3. **Chat** → a standalone task (sandboxed conversation).
4. **Autopilot** → time-driven (cron) or event-driven (a NATS subject pattern fires it).
   Example: `tel.greenhouse.soil.moisture < 20` publishes an event → autopilot creates an
   issue "water greenhouse" assigned to the irrigation device. **This is where the physical
   world drives the task board** — telemetry → task, automatically.

### Collaboration primitives (kept from Multica)
Comments (threaded), reactions, mentions, the `in_review` state, squads-as-routing, and the
anti-loop guards (self-mention skip, one-active-task-per-actor-per-issue dedup). Reviews are
convention over comments + status — Legacy reviews worker output exactly this way today.

### Why NATS beats Multica's REST-poll core
Multica's daemon polls every 3s. Ours is push: an actor's daemon holds a JetStream consumer
on its command subject and gets work the instant it's assigned — lower latency, no wasted
cycles, and the same mechanism scales to hundreds of nodes. Multica also runs server-side
data with BYO local compute; we make the **whole stack** local-first and self-hostable.

### Relationship to existing tools
- **Hermes Kanban** (what's building this repo right now) = the *engineering* board for the
  build. The Task Engine is the *product* — the homestead-wide board. The Kanban board can
  later become one view onto the Task Engine, or feed it. Not urgent; noted to avoid
  confusion.
- **Todoist** (Zeus's current human tool) = import source + a `human` notification channel.
  We can mirror human-assigned issues to Todoist for familiar reminders during transition.

### Implementation — om_task.py v2.0.0 (LIVE)

`om_task.py` (`services/task-engine/`) implements the Task Engine as a NATS `$SRV`
micro service (`om-task` v2.0.0). Running as LaunchAgent `com.onemind.om-task`.

**Subjects implemented** (`$SRV` endpoints + plain subscribes):

| Subject | Type | Notes |
|---------|------|-------|
| `task.create` | request-reply | returns `{task_id, status}` |
| `task.update.*` | request-reply | `*` = task_id; merges allowed fields |
| `task.approve.*` | request-reply | `{approved: bool, approver?, reason?}` |
| `task.list` | request-reply | filters: status, agent_id, domain, archived, limit |
| `task.get.*` | request-reply | `*` = task_id |
| `task.calendar` | request-reply | `{from_date, to_date}` → tasks with due_date in range |
| `task.inbox` | request-reply | inbox/pending + overdue + due today |
| `task.project.create` | request-reply | `{name, color, icon}` |
| `task.project.list` | request-reply | non-archived projects |
| `agents.hb.>` | plain subscribe | heartbeat aggregator → `agent-registry` KV |

**KV buckets:**

| Bucket | TTL | History | Purpose |
|--------|-----|---------|---------|
| `task-index` | none (permanent) | 10 | All tasks |
| `projects` | none (permanent) | 5 | Project definitions |
| `task-notify` | 48h | 1 | Dedup log for sent reminders |
| `agent-registry` | none | 5 | Agent heartbeat registry |

**Task schema fields:** `task_id`, `title`, `description`, `agent_id`, `priority` (normal/high/urgent/low), `domain` (so/hp/le/ge), `tags`, `labels`, `due_ts`, `due_date`, `start_date`, `is_all_day`, `recurrence`, `reminder_mins`, `project_id`, `requires_approval`, `status` (inbox/pending/in_progress/waiting/done/completed/failed), `progress` (0–100), `completed_at`, `archived`, `created_at`, `updated_at`.

**Heartbeat scan (60s):** detects stalled tasks (pending + agent + >10min stale), overdue tasks (`due_ts` passed), and fires due-soon reminders from `reminder_mins`.

**Push notifications (background delivery):**
All due-soon reminders, approval requests, and overdue alerts push to:
- **ntfy.sh** topic `onemind-zeus-64c8cb63c3ee91bc` — OS-level push to iOS/Android/web via ntfy app
- **Telegram** bot `@omos_hermes_bot` → chat ID `577573558` (requires `/start` sent to bot once)

Config via env vars: `NTFY_TOPIC`, `NTFY_BASE`, `TELEGRAM_TOKEN`, `TELEGRAM_CHAT_ID`.

**Events published:**
- `evt.task.created` — on task create
- `evt.task.updated` — on task update
- `evt.task.due.soon` — N minutes before due_ts
- `evt.task.{id}.overdue` — when due_ts passes
- `evt.task.{id}.stalled` — when assigned task hasn't updated in >10min
- `evt.task.{id}.complete` — on done/completed status
- `evt.task.{id}.approved` / `.rejected` — on approval

**REST Bridge** (`task-bridge/bridge.js`, port 8401): Node.js Express bridge that translates REST HTTP ↔ NATS for the Flutter app. Reads KV directly via `$JS.API.STREAM.MSG.GET` (raw JetStream API — avoids broken `@nats-io/jetstream` v3 JS client). Also streams task events to the Flutter app via SSE at `GET /events`. See `docs/WEB-DEPLOYMENT.md` for endpoints.

**Subject grammar debt:** Design specifies `svc.task.*` prefix; implementation uses flat `task.*`. Rename to `svc.task.*` when subject grammar is enforced globally.

---

## 4b. Session Bridge (LIVE — 2026-06-13)

The **Session Bridge** (`services/session-bridge/bridge.js`, port 8402) is a read-only REST API that exposes unified session history from two local sources:

| Source | Format | Key detail |
|--------|--------|------------|
| `~/.hermes/state.db` | SQLite (`better-sqlite3`) | Timestamps stored as epoch seconds (float), not ISO — converted on read |
| `~/.claude/projects/**/*.jsonl` | Claude Code JSONL | Sessions deduplicated by `sessionId` (keep most recent file) |

Sessions are presented in a unified shape regardless of source. The Flutter app uses this for the **Sessions** page under the AGENTS nav group.

Running as LaunchAgent `io.onemind.session-bridge`.

---

## 5. Local-First & the "One Download" Install

The vision: **one person downloads OneMind OS, runs one installer, and the whole nervous
system stands up.** Concretely the installer provisions and supervises:

```
OneMind OS bundle
├── nats-server         (substrate)        ─┐
├── clickhouse          (time machine)      │ supervised by launchd (mac) /
├── omd  (daemon)       (local compute)     │ systemd (linux) / s6 (container)
├── agent harnesses     (Hermes, etc., discovered or installed)
└── OneMind Command     (the app — desktop bundles omd; mobile/web connect to it)
```

On first run `omd` registers the node, the app connects to `localhost:4222`, and if other
agents are already on the machine they're discovered over the bus. Add a second machine: it
installs the same bundle, leaf-nodes into core over Tailscale, and its node + agents appear
in the app automatically. **Locality, portability, sovereignty** — your three values, encoded.

---

## 6. Application Module Structure (Flutter — feature-first)

Organized so every future feature is an additive folder, never a cross-cutting refactor
(Principle #7). Target layout for the restructure (card **P2-S**):

```
lib/
  core/          transport-agnostic: config, auth, envelope, result types
  bus/           NATS client wrapper, subscriptions, request-reply, $SRV discovery,
                 sentinel/request-many (Orbit-style patterns app-side)
  models/        Actor · Issue · Task · Event · Node · Device · Envelope  (pure data)
  features/      one self-contained module per capability:
    dashboard/     home, fleet health, live event ticker
    chat/          streaming agent conversations
    agents/        agent layer — roster (discovery), review/edit harness config via SDK
    tasks/         OneMind Task Engine UI (issues, tasks, board, autopilots)
    timemachine/   ClickHouse-backed history, replay, "what happened when"
    fleet/         nodes + devices + capabilities map
    command/       system control surface
  shared/        widgets, theme (domain colors), formatters, errors
  app.dart       shell + GoRouter (StatefulShellRoute)
  main.dart      boot → splash → shell (no setup gate; transports in Settings)
```

Rule: a feature module may depend on `core/`, `bus/`, `models/`, `shared/` — **never on
another feature module.** Cross-feature needs go through the bus or a shared model. This
keeps features independently buildable, testable, and removable.

---

## 7. Feature Roadmap (where the bolt-ons go)

Because the structure is set now, each of these is "add a feature module + subjects," not a
re-architecture:

- **Time Machine** — ClickHouse query UI: scrub a timeline, replay event windows, diff state.
- **Fleet** — live node/device map, capabilities, health from `tel.>`.
- **Insights** (Synadia-equiv) — `$SYS` + nats-surveyor → dashboards.
- **Connect** (Synadia-equiv, Wombat) — SaaS/IoT bridges in & out of the bus.
- **Workloads** (Synadia-equiv, NEX) — deploy services to any node from the app.
- **Control Plane** (Synadia-equiv) — security/account management surface.
- **Autopilots** — visual rule builder (telemetry/cron/webhook → task).
- **Voice** — already scaffolded; becomes a trigger + notification channel.

---

## 7b. Live Pages Built (2026-06-13)

All pages are wired into the Flutter app's nav rail and routing. Data comes from the bridges above — no SQLite in the app.

### Task pages (OPS group in nav rail)

| Route | Page | Data source |
|-------|------|-------------|
| `/tasks/inbox` | **Task Inbox** — overdue / today / upcoming / no-date sections, tap to complete | `GET /tasks/inbox` |
| `/tasks/calendar` | **Calendar** — month view with domain-colored event dots, day list below | `GET /tasks/calendar` |
| `/tasks/board` | **Task Board** — Kanban: Inbox → In Progress → Waiting → Done, FAB to create | `GET /tasks` |
| `/tasks/:id` | **Task Detail** — full editor: status, priority, domain, due date, project, recurrence, reminders, labels; 1s debounced auto-save | `GET/PATCH /tasks/:id` |
| `/notifications` | **Notifications** — live SSE-fed list, mark-all-read | SSE `/events` |

### Session pages (AGENTS group)

| Route | Page | Data source |
|-------|------|-------------|
| `/agents/sessions` | **Sessions** — unified list of Hermes + Claude Code sessions, filter chips (All/Hermes/Claude/Active) | `GET /sessions` |
| `/agents/sessions/:id` | **Session Detail** — full transcript, role-colored bubbles, collapsible thinking blocks, tool call chips | `GET /sessions/:id/transcript` |

### Push notifications (background — all devices)

Task reminders, approvals, and overdue alerts reach all devices even when the app is **closed**:
- **ntfy.sh** — install ntfy app, subscribe to `onemind-zeus-64c8cb63c3ee91bc`. Works on iOS, Android, and web browser.
- **Telegram** — `@omos_hermes_bot`. Send `/start` to the bot once to activate.

---

## 8. Open Decisions (need Zeus or a spike)

1. **Desktop shell:** Flutter desktop (one codebase) vs Tauri/React for the daemon-bundling
   build. Leaning Flutter desktop to keep one UI codebase; revisit if FFI/daemon supervision
   gets awkward.
2. **`omd` language:** Go (matches NATS/NEX ecosystem, single static binary, easy cross-compile
   to Pi/ARM) vs Python (matches our agent tooling). Leaning **Go** for the daemon core with a
   Python extension lane for agent glue.
3. **Issue/Task storage:** pure JetStream+KV vs a small embedded store (SQLite) per node synced
   over the bus. Leaning **JetStream+KV** to stay bus-native (Principle #1).
4. **Subject cutover** `events.>` → `evt.>` (card P2-T) — timing and dual-publish window.

---

## 9. What We Build Today (the foundation)

1. **P1-A** (running) — NATS accounts + streams + KV. The substrate's spine.
2. **P2-S** — restructure `lib/` to the feature-first layout above (no behavior change).
3. **P0-E** — define the **event envelope** standard + write it to the `om-schema` KV bucket.
4. **P2-T** — plan the ClickHouse subject rebind (`events.>` → semantic set), no data loss.

Everything else in the roadmap sits on these four. Get them right and we don't hit the wall.

## 10. Repository & Distribution Topology

> **Decision (2026-06-10): Split core from UI. Do NOT fork NATS. Ship a mirrorable bundle.**

### 10.1 Two repos, one contract

| Repo | Is | Language | Visibility | Mirrorable |
|------|-----|----------|-----------|-----------|
| **onemindos-core** | The OS: daemon, services, NATS config, schemas, installer, canon docs | Go (daemon/services) + config | Public-capable (others mirror & run) | ✅ yes — this is "the core other people can mirror" |
| **onemindos-command** | The UI: the Flutter app (single pane of glass) | Dart/Flutter | Private (your product face) | replicates a machine, not mirrored as a template |

**Dependency direction is one-way:** `command` depends on `core`'s published **contract**
(subject grammar + event envelope + service API schemas). `core` never depends on `command`.
The UI is a replaceable client; the OS is the durable asset. Anyone can write their own UI
against the same contract — that's the point of the split.

```
onemindos-core/                 ← the OS (mirror this to run OneMind)
  nats/         account layout, JetStream stream defs, leafnode templates (stock NATS, configured)
  omd/          the OneMind daemon (Go) — local fs/shell access, runs agent harnesses
  services/     meta-services on the bus:
    task-engine/      issues + tasks + actors + triggers
    control-plane/    accounts/security API over $SYS  (Synadia Control Plane equiv)
    insights/         surveyor + exporter → ClickHouse  (Synadia Insights equiv)
    connect/          Wombat bridges                     (Synadia Connect equiv)
  schemas/      THE CONTRACT — event envelope, subject grammar, JSON Schemas, codegen
  install/      one-download installer + supervisor units (launchd/systemd/s6)
  docs/         CANON — ARCHITECTURE.md, NATS-DOCTRINE.md, MASTER-PLAN.md live HERE

onemindos-command/              ← the UI (this repo)
  lib/          Flutter app
  contract/     generated Dart models from core/schemas (pinned version)
  docs/         UI-specific only; links to core canon
```

### 10.2 Do NOT fork NATS — distribute, configure, extend

Forking NATS would be a strategic mistake:
- You'd inherit maintaining the fork forever — merging upstream security patches by hand.
- You'd **lose the ecosystem**: NEX, Wombat, orbit, surveyor all target upstream NATS. A fork
  diverges and they stop working.
- Synadia ships fast; a fork falls behind within weeks.
- Zero upside — everything we need is achievable **without touching server source.**

What "a core others can mirror" actually means (three different senses, none is a fork):

| You said | Real mechanism | Layer |
|----------|---------------|-------|
| "a core other people can mirror" | They **clone & run `onemindos-core`** — stock NATS + our config + our services. Distribution, not forking. | repo |
| "the UI replicates what's on the machine" | App is **bus-first/stateless** — subscribes to live NATS state and renders it. Already our Principle #1. | app |
| data "mirroring" between machines | **NATS stream mirrors + leaf nodes** — built-in. Barn node mirrors to core. | bus |

We **extend** NATS by adding **services on the bus** (`$SRV` micro-services), never by changing
the server. The server stays stock and upgradeable; our intelligence lives in `services/`.

### 10.3 Migration from today's single repo

`onemindos-command` currently holds the canon docs because we bootstrapped here. Plan:
1. Create `onemindos-core`, scaffold the tree above.
2. **Move** ARCHITECTURE.md, NATS-DOCTRINE.md, MASTER-PLAN.md → `core/docs/` (they describe
   the whole system). Leave forwarding stubs in `command` pointing to core.
3. `command/AGENTS.md` keeps the orchestration section; add "canon lives in onemindos-core."
4. The in-flight P1-A NATS config + `docs/infra/NATS-SETUP.md` belong in **core**, not command.
   (P1-A ran in command for bootstrap convenience; its output relocates to core/nats + core/docs.)
5. Stand up `core/schemas/` as the shared contract; `command/contract/` consumes it.

---

## 11. Navigation — Discovery-Driven Single Pane of Glass

> "We need a menu for everything we add, because the connections must reflect the information."

The menu is **not a hardcoded list**. It's a **registry** populated partly by built-in modules
and partly by **live discovery from the bus**. Add a service, node, agent, stream, or bridge —
it appears in the menu automatically, because navigation reflects what is actually connected.

### Two zones

**Zone A — OS Console (the NATS single pane of glass; Synadia-suite equivalents):**

| Surface | Mirrors Synadia | Data source on the bus |
|---------|----------------|------------------------|
| **Fleet** | (Platform/topology) | node registrations + `tel.{node}.>` + leafnode/$SYS topology |
| **Streams** | JetStream UI | `$JS.API.STREAM.*`, `$JS.API.CONSUMER.*` |
| **Insights** | Insights | surveyor/exporter metrics → ClickHouse (the Time Machine) |
| **Control Plane** | Control Plane | `$SYS` accounts/connections; control-plane service API |
| **Protect** | Protect | auth callout policy + connection allow/deny (Phase 7) |
| **Agents** | Synadia Agents | `$SRV.PING.agents` discovery + Agent SDK; review/edit harness config |

**Zone B — OneMind Apps (the product features that ride the OS):**

Tasks (Task Engine) · Time Machine · Chat · Codex · Command · Voice · Home/Smart-home ·
+ every future module.

### The registry

```
NavSection {
  id, label, icon, zone (os|app),
  source: static | discovered,
  visible_when: capability/permission predicate,
  badge: live count from a subscription (e.g. # online nodes, # running tasks),
  route
}
```

- **Static sections** (built-in surfaces) register at app start.
- **Discovered sections** subscribe to discovery subjects and appear/disappear with reality:
  a new node → Fleet gains an entry; a new agent → Agents gains a row; a new Connect bridge →
  a Connections entry; a new stream → Streams updates. **Connections reflect information.**
- Sections render **live badges** from bus subscriptions, so the menu itself is real-time.

This is what makes the app a true real-time NATS front end *and* an extensible product shell:
one navigation system, two zones, driven by what the bus actually knows.
