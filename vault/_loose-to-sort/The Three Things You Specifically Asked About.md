# The Three Things You Specifically Asked About

## 1. Local Computer Access (files, shell — "the common structure with Claude/Hermes")

The answer is the `omd` daemon (L3). It runs on each machine and, on boot, announces itself to the bus as a **node** actor with capabilities — `fs.read`, `fs.write`, `shell.exec`, `agent.run:hermes`, `gpu`, `camera`, etc.

When the app or an agent needs something local, it doesn't reach into the machine directly. It publishes a request to that node's subject (`cmd.node.{id}.task`), the local daemon does it with real privilege, and replies over the bus.

**Why this matters:** a file read on *this* Mac and a file read on the *barn Pi* 200m away use the identical path. Location-transparent, permission-gated, and every action is recorded. The web app (browser sandbox, no raw filesystem) still gets full local access this way — it just talks to the daemon on localhost. Other agents already on the machine get discovered automatically over the bus. That's your "common structure."

---

## 2. ClickHouse as the Time Machine

It's *already doing this* — verified live. There's a NATS-engine table subscribing to the bus, feeding a partitioned table with a 10-year retention. `onemind.events` has 31,265 rows; `onemind.metrics` has 34 million. Every semantic event lands there permanently.

**The split that prevents confusion:**
- **JetStream** = operational memory (replay the last 30 days, exactly-once delivery to consumers)
- **ClickHouse** = analytical history (SQL: "what happened in the barn on March 3rd at 4pm")

One feeds operations, the other feeds insight. The future **Time Machine** feature in the app is just a UI scrubbing over ClickHouse.

**One cleanup item (tracked as card P2-T):** it currently listens on the old subject `events.>`. Needs to rebind to the new grammar (`evt.>` `cmd.>` `tel.>` `agents.>`) without dropping existing rows.

---

## 3. The OneMind Task Engine

Studied Multica's model in depth (read 18 of their docs). Their core insight is right — adopting it and generalizing past where they stopped.

### The Issue/Task Split (the key data-model decision)

- **Issue** = the durable thing to do. "Re-roof the barn." Lives forever, `OM-123`, assigned to any actor, has status/priority/comments/children.
- **Task** = one *attempt* at it. Ephemeral, its own lifecycle (queued → running → done/failed), with timeouts, retries, session resume. One issue spawns many tasks over its life.

This cleanly separates "the work" (durable, human) from "an attempt" (ephemeral, machine) — which is why retries, multi-agent reruns, and audit all work without hacks.

### Where I Go Past Multica — The Polymorphic Actor

Multica made *agents* first-class members alongside humans. OneMind makes the **physical world** first-class too:

```
Actor = human | agent | node | device | squad
        each with capabilities + an executor binding
```

**Critical distinction kept from Multica: NOTIFY vs TRIGGER.**
- Humans get *notified* (an inbox, like Todoist) and choose to act.
- Machines get *triggered* and execute immediately.
- Same task board, same issue, different executor.
- Adding "drone" or "CNC" later is just a new actor type with new capabilities — **never a schema rewrite.**

### And the Part That Makes It an OS, Not a To-Do App

Telemetry drives the board. `tel.greenhouse.soil.moisture < 20` fires an event → an autopilot creates an issue "water greenhouse" → assigned to the irrigation device → it executes. The physical world creates and completes its own tasks. NATS Connect (Wombat) later lets a task fire a SaaS action or move a robot.

### Why Ours Beats Multica's Core

They poll every 3 seconds; we push. An actor's daemon holds a live subscription and gets work the instant it's assigned — lower latency, and it scales to hundreds of nodes on the same mechanism.

---

## The "One Download" Vision (encoded in the design)

One installer ships `nats-server` + `clickhouse` + `omd` + agent harnesses + the app, pre-wired and supervised (launchd/systemd/s6). First run: daemon registers the node, app connects to localhost, existing agents auto-discovered. Add a second machine: same bundle, leaf-nodes into core over Tailscale, its node + agents appear in the app automatically. **Locality, portability, sovereignty** — your three values, baked into the architecture.

---

## What We Build Today — The 4 Foundation Pieces

Get these right and we never hit the wall:

| Card | Task |
|------|------|
| **P1-A** | NATS accounts + streams + KV (worker running now, 15 min in, healthy) |
| **P2-S** | Restructure the app's `lib/` into the feature-first layout (so every future feature is an additive folder, never surgery) |
| **P0-E** | Define the versioned event envelope + write it to the schema registry |
| **P2-T** | Plan the ClickHouse subject rebind, zero data loss |

Everything else — Time Machine, Fleet map, Insights, Connect, Workloads, Autopilots, Voice — sits on these four as additive modules.

The full detail is in `ARCHITECTURE.md` (committed, on main). It has the layer diagram, the complete Actor schema, the Task Engine subject grammar, the module layout, and an Open Decisions section (§8) with four calls — most notably: `omd` daemon in Go or Python? Lean Go (single static binary, cross-compiles clean to Pi/ARM, matches the NATS/NEX ecosystem) with a Python lane for agent glue.

---

## Architecture Layers

```
L5  PRESENCE      GPT-realtime + wearables (glasses HUD, watch, projector)
L4  ORCHESTRATOR  Hermes/Legacy — memory, tools, agent routing, L0–L3  ← your "manager"
L3  SPECIALISTS   domain LLM agents (infra, security, comms, business)
─────────────────────────────────────────────────────────────────────
L2  WM CONTROL    per-robot world models, latent MPC  ← NEW (V-JEPA 2-AC, LeWM, Cosmos)
L1  ALWAYS-ON     local WM perception, 24/7, token-free, emits SURPRISE  ← NEW
L0  FABRIC        robots, drones, vehicles, printers, cameras — via NATS
```

The omnipresent "manager" you want = **L4 Hermes, unchanged.** World models slot in at **L1–L2**, giving you the two things LLMs genuinely can't: token-free 24/7 perception, and fast physical action planning for robots/drones/vehicles. All on your existing NATS fabric — WMs become just more nodes on the mesh.

**2027 physical AI embodiment is not a guess — it's a roadmap with hardware already shipping.** You are positioning at the exact intersection of NATS + physical AI + sovereign infrastructure before anyone else has put it together as a unified personal OS.

---

## OneMind Cockpit

```
ONEMIND COCKPIT
│
├── FABRIC VIEW
│   ├── Live agent roster (heartbeats, sessions active)
│   ├── Node map (house, barn, office, drones, robots — live status)
│   ├── Event feed (evt.> scrolling live)
│   └── Command dispatch (publish to bus, governance pre-check)
│
├── OPERATIONS
│   ├── Tasks / Missions (Kanban — 6 columns)
│   ├── Scheduling (natural language → cron)
│   ├── Inbox / Feed
│   └── Knowledge / Skills
│
├── INTELLIGENCE
│   ├── Agent sessions (chat with any harness)
│   ├── Memory graph
│   ├── Cost / token tracking
│   └── Session replay (JetStream)
│
├── PHYSICAL
│   ├── Drone / robot status
│   ├── Camera feeds
│   ├── Sensor telemetry (live charts)
│   └── Location map (GPS assets)
│
└── GOVERNANCE
    ├── Approval queue (cmd.> pending)
    ├── Audit log
    ├── Security events
    └── Access levels (L0–L4)
```

---

## Think of It Like Building a House

Your house has three separate things:

1. **The foundation + electrical wiring** — you never see it but everything runs through it. This is **NATS**. It's the nervous system.
2. **The rooms and appliances** — the stuff that actually does work. These are your **Python services** (Legacy, voice agent, watch service, sensors, bridges).
3. **The walls, windows, and front door** — what you see and touch. This is your **app** (Flutter mobile, React web dashboard).

The confusion with Flutter is almost certainly **not Flutter's fault** — it's that the "electrical wiring" wasn't defined clearly enough before you started building rooms. You were wiring outlets before the panel was installed.

---

## Who Does What in OneMind OS

```
┌─────────────────────────────────────────────────────────────┐
│            WHO DOES WHAT IN ONEMIND OS                       │
├──────────────┬──────────────────────────────────────────────┤
│  PYTHON      │  THE BRAIN & NERVOUS SYSTEM                   │
│              │  • Legacy (Hermes) — the main agent           │
│              │  • Voice agent (LiveKit) — ears & mouth       │
│              │  • Watch service — autonomous governor        │
│              │  • NATS subscribers — event reactors          │
│              │  • IoT bridges — sensor readers               │
│              │  • ML / world models — perception             │
│              │  Lives 24/7 on your server/Pi                 │
├──────────────┼──────────────────────────────────────────────┤
│  DART/       │  THE EYES & HANDS (what you see & touch)      │
│  FLUTTER     │  • Mobile app (iPhone, Android, tablet)       │
│              │  • Cockpit panels (home dashboard on iPad)    │
│              │  • Connects directly to NATS :8080 WebSocket  │
│              │  Runs on every device you carry               │
├──────────────┼──────────────────────────────────────────────┤
│  REACT /     │  THE COMMAND CENTER (big screen / browser)    │
│  TYPESCRIPT  │  • Web dashboard (agenttap-style live view)   │
│              │  • Fabric map, drone feeds, sensor charts     │
│              │  • Connects to NATS :8080 WebSocket           │
│              │  Runs in Brave on any machine                 │
├──────────────┼──────────────────────────────────────────────┤
│  GO          │  THE BRIDGES (tiny always-on daemons)         │
│              │  • NATS ↔ Home Assistant bridge               │
│              │  • NATS ↔ ESP32/WireClaw bridge               │
│              │  • NATS ↔ drone MAVLink bridge                │
│              │  Small binary, fits on Pi Zero, runs forever  │
├──────────────┼──────────────────────────────────────────────┤
│  SQLITE      │  THE MEMORY (everything stored locally)       │
│              │  • Tasks, missions, knowledge, skills         │
│              │  • Session history, events log                │
│              │  • Works offline, works on Pi, no server      │
└──────────────┴──────────────────────────────────────────────┘
```

---

## The Real-Time NATS Front End + The Menu — Two Zones, One Discovery-Driven Registry

The app is *both* a real-time NATS console (the Synadia-suite equivalents) *and* the product shell for everything you'll add. Navigation splits into two zones:

**Zone A — OS Console** (your NATS single pane of glass):
Fleet → Streams → Insights → Control Plane → Protect → Agents — each one reading live from the bus (`$SYS`, `$JS.API.*`, `$SRV`, surveyor + ClickHouse).

**Zone B — OneMind Apps** (features riding the OS):
Tasks → Time Machine → Chat → Codex → Command → Voice → Smart-home → + everything future.

**The critical part — the menu is a registry, not a hardcoded list.** Sections are either *static* (built-in) or *discovered* (subscribe to a discovery subject). Add a node → Fleet gains an entry. Add an agent → Agents gains a row. Add a Connect bridge → a Connections entry appears. Each section carries a **live badge** (# online nodes, # running tasks) from a bus subscription. The menu literally renders the state of the system.

---

## Contact

**Wally Quevedo** — wally@nats.io
