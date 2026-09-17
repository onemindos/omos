# One Mind NATS OS — Architecture Plan (Working Doc)

**Date:** 2026-06-05
**Author:** Zeus Delacruz + Legacy
**Status:** Active planning — companion to `ONEMIND-NATS-OS.md`
**Purpose:** Reconcile the original vision doc against the as-built reality, capture the honest assessment, the full coverage map (A–G), the hub-redundancy decision, and the universal subject taxonomy (the spine).

---

## 0. Reality Check — You Are Not "Pre-Architecture." You Are ~80% Built.

The original doc reads greenfield. The live server (probed 2026-06-05) says otherwise:

- **JetStream LIVE** — domain=`onemind`, v2.14.1, 5d19h uptime, 12+ streams + KV buckets
- **Edge fabric already wired** — leaf node (7422), MQTT (1883), WebSocket (8080) all listening
- **Nodes already on the bus** — ha-bridge (~1ms), saas-action-worker, notify-router, ClickHouse experience store (`events.>`, 5yr), n8n webhook gateway, nats-fabric + zeus-context plugins

| Original Phase | Reality | Status |
|---|---|---|
| 1 — JetStream + schema | Live, domain=onemind, 12+ streams | ✅ DONE |
| 2 — Automation nodes | n8n live, ha-bridge live, Redpanda replaced | ✅ MOSTLY |
| 3 — Wire Hermes agents | Native NATS toolset + 2 plugins live; only OpenClaw auto-registers as `$SRV` | 🟡 PARTIAL |
| 4 — Meta-Agent (brain) | Nothing subscribes to `>` and governs; Overwatch cron paused | 🔴 NOT BUILT |
| 5 — Synadia NEX migration | NEX not installed | 🔴 NOT STARTED |
| 6 — Experience engine | ClickHouse store LIVE; no pattern detection / replay-onboarding | 🟡 STORE YES, INTELLIGENCE NO |

**This is a reconciliation + completion effort, not a build-from-zero.**

---

## 1. The Thesis (Validated) + Honest Assessment

**The thesis:** Make NATS the central nervous system of the entire digital life. Every service, agent, app, device, robot, drone, human, location = an addressable node on one fabric. Legacy is the meta-agent. AI moves to the edge. Fleets, not harnesses.

**Where Zeus is right — fully:**
- Timing is real. Synadia pivoted their *entire company* to "the nervous system for agentic AI." NVIDIA routes GPU inference on a NATS supercluster. This aligns with where serious infra is going.
- "Agent fleets > agent harnesses" is the correct call. Everyone polishes single harnesses; the fleet-on-a-bus-with-wrapped-services model is a level up.
- Edge AI via leaf nodes is elegant and correct — same binary everywhere, local buffering, telemetry to hub.

**The 05 ↔ 08 unification insight (the core realization):**
- **05 Stack** = digital nodes — 166 services/tools (what your world can *do* — capabilities). Its integration map already says `All platforms → NATS`.
- **08 Entity Registry** = world nodes — "every noun in your world," 5 categories, access L0–L4 (what your world *is* — entities).
- **These are ONE universal node registry viewed from two sides.** The bus is the membrane that makes the digital catalog and the physical catalog the same addressable fabric. That is One Mind, literally: one namespace for the entire existence.

**Honest concerns (the risks that matter):**
1. **The risk is sprawl without a spine, not technical feasibility.** Services are "up but disorganized." Adding nodes onto a foundation with no identity, no canonical schema, no live registry = more chaos faster. **Identity + subject schema + registry-as-KV must land before a single new node.**
2. **"Wrap everything" silently becomes "maintain everything."** Every bespoke bridge is code owned forever. The native discipline (one wrapper pattern, micro-services, import/export) is the only thing between you and a 200-daemon swamp.
3. **Not everything deserves to be a streaming node.** Filter: does it emit/react in real time, or do you call it occasionally? Bus-native for emitters/reactors; on-demand wrap for monthly-touch APIs.
4. **Robots/drones: bridge the seam, do NOT replace their language.** Inner control loops (PX4/ArduPilot/MAVLink 100–400Hz, ROS2/DDS) cannot tolerate general-bus latency. Onboard loop stays native; only telemetry + high-level commands cross to the hub.
5. **The MBP is a single point of failure** for the entire nervous system. See §3.

---

## 2. The Two-Brain Model (Corrects the Doc's "AI subscribes to `>`")

An LLM seeing every message is financially and architecturally impossible — the bus does thousands of msgs/sec; 99.9% is `switch turned on` and needs a *rule*, not a *thought*. The brain splits in two (this maps to what already runs):

| | **Autonomic layer** (spinal cord) | **Cortex** (Legacy) |
|---|---|---|
| What | Daemon subscribed to `>`, deterministic rules | LLM reasoning, invoked on escalation only |
| Speed | sub-ms, zero token cost | seconds, full token cost |
| Job | Route, enforce hard ACLs, dedup, aggregate state, decide what's worth waking the cortex for | Judgment, novel situations, standing orders, fleet orchestration |
| Today | `subscriber.py` — already live | Legacy (me) — session-bound, not yet always-on |
| Bridge | escalates via `gov.escalate.>` → wakes cortex | — |

Phase 4 isn't a new always-on LLM. It's **`subscriber.py` evolved into a real router/policy engine + a clean escalation channel to Legacy.**

---

## 3. Hub Redundancy — Killing the SPOF

**Decision: 3-node JetStream cluster. Never a 2-node standby (it's worse than one).**

JetStream replication = Raft consensus = needs a *majority* alive to accept writes:

| Cluster size | Majority | Survives | Verdict |
|---|---|---|---|
| 1 (today) | 1 | 0 | SPOF |
| 2 (the trap) | 2 | 0 | Lose either → read-only. Two things can now kill you |
| **3 (answer)** | **2** | **1 failure** | **Lose any one, bus keeps writing ✅** |
| 5 | 3 | 2 | Overkill until fleet is huge |

**How it works:**
- Each node runs `nats-server` with a `cluster {}` block routing to the other two (port 6222, Tailscale/LAN).
- Critical streams + KV set `replicas: 3` — three physical copies.
- Clients list all three URLs; NATS gossips topology → auto-discover + auto-reconnect to survivors. Failover is automatic and invisible (ms).
- `domain: onemind` stays. Mixed-OS fine.

**Topology — two workhorses + one tiebreaker (not three equal nodes):**
```
MBP (hub, R3) ──cluster routes── SPARE (2nd node, R3) ──── TIEBREAKER (vote-only)
                                                            HA box / Pi / $5 VPS
```
MBP + spare carry load. The 3rd member just needs to be reliably online to cast the deciding vote. If the HA box is Pi-class, use it as tiebreaker only (modest storage), not a full data node. **Verify HA box specs before electing it.**

**Cluster vs Leaf — complementary, not the same:**
| Layer | Mechanism | Job | Lives |
|---|---|---|---|
| Cluster (3 nodes) | Raft, R3 | Hub HA — no SPOF | Local LAN (Raft is chatty — never over WAN) |
| Leaf nodes | Leaf + local JetStream | Edge survival — buffer offline, sync later | Drones, robots, barn, remote |

**Sequencing:**
- **Now (insurance):** NATS on spare as 2nd server + **Mirror** critical streams. Data survives MBP loss, manual failover. Low-risk.
- **When load-bearing (locks/security/presence):** Promote to full 3-node R3 cluster. Automatic failover.

**Cost: zero new hardware, zero subscription.** The idle box finally earns its keep.

---

## 4. Coverage Map (A–G) — Everything That Must Be Covered

The 05+08 union, organized as the full scope of "design my digital life around the bus." 36 items, seven layers.

### A — FOUNDATION (the spine — must come first)
1. **Identity & Security** — per-node NKeys/JWT, accounts, `$SYS`, hand-rolled policy (Synadia-Protect-equivalent, self-hosted)
2. **Subject Schema / Taxonomy** — canonical naming for every node + every event *(the single most important artifact — see §5)*
3. **Node Registry as KV** — live state layer over the 08 markdown registry
4. **Event Taxonomy** — the canonical `evt.>` firehose, what mirrors in, retention tiers

### B — NODE CLASSES (everything that becomes a node)
5. **Agents** — Legacy (meta) + cc, pi, openclaw, codex, gemini (ACP-over-NATS)
6. **Machines** — MBP hub, HA box, spare, VPS → leaf/cluster nodes
7. **Services/SaaS** — all 166 from 05, wrapped as micro-services or on-demand
8. **IoT/Smart Home** — HA entities, SwitchBot, cameras, sensors (bridged, live)
9. **Devices** — iPhone, tablets, peripherals
10. **Robots** — DIMOS, OM1, Unitree, FarmBot (ROS2/MAVLink seam-bridge)
11. **Drones** — PX4/ArduPilot/MAVLink (telemetry + command bridge only)
12. **Humans** — Zeus, Cielo, kids (presence, identity, comms routing, access L0–L4)
13. **Locations** — house, barn, office (HA areas, presence)
14. **Vehicles** — cars, motorcycle, plane (GPS/telemetry)
15. **Pets / Ecological** — Cerberus (trackers), land/garden/greenhouse (FarmBot, sensors)

### C — INGESTION (record everything)
16. **Webhook gateway** (n8n, live) — all SaaS webhooks → `evt.service.>`
17. **Outbound actions** (saas-action-worker) — agents act on SaaS via `svc.>`
18. **Polling bridges** for no-webhook services
19. **Device/system telemetry** — HA, iPhone, GPS, Netdata → `tel.>`

### D — INTELLIGENCE (the two-brain)
20. **Autonomic governor** — `subscriber.py` evolved into fast router/policy engine
21. **Legacy as Meta-Agent** — escalation target, judgment, fleet orchestration
22. **Experience engine** — ClickHouse pattern detection + replay-onboarding (store is live)
23. **Governance** — hard ACLs (identity) + soft rules (meta-agent)

### E — INTERFACE
24. **Flutter app** — the one interface (per doctrine)
25. **Real-time bus dashboard** — live view + audit
26. **Voice / CLI / device nodes**
27. **Notification router** (live)

### F — OBSERVABILITY & OPS
28. **Experience store** (ClickHouse, live)
29. **Server monitoring** ($SYS, /varz)
30. **Langfuse** (agent traces)
31. **Backup / DR / hub redundancy** (the SPOF fix — §3)

### G — EDGE & PHYSICAL FUTURE
32. **Leaf node topology** (HA box → spare → drones/robots)
33. **Local JetStream offline buffering**
34. **Sources & Mirrors** store-and-forward
35. **NEX workload deployment** (deferred until fleet is real)
36. **AI-on-the-edge** (local models/agents on nodes)

---

## 5. The Universal Subject Taxonomy (THE SPINE)

> The genome. Every node, event, command, and access rule derives from this. The naming law that unifies 05 (capabilities) and 08 (entities) into ONE address space.

### 5.1 The 5 planes (top-level message kinds)

Lead with **plane** (not node) because planes have different retention/QoS/persistence and map cleanly to JetStream stream boundaries and governance rules. `svc` was considered and deliberately collapsed into `cmd`/`qry` — a SaaS service is just a node of class `service`. This keeps governance to one rule: every mutation in the entire digital life is `cmd.>`.

```
<plane>.<class>.<id>.<path>

evt.<class>.<id>.<path>    EVENTS    — something happened. Immutable, fan-out, recorded. → JetStream + ClickHouse
cmd.<class>.<id>.<action>  COMMANDS  — do something (mutates state). Request/reply. GOVERNED — hard ACLs.
qry.<class>.<id>.<q>       QUERIES   — read something. Request/reply. Read-only.
tel.<class>.<id>.<metric>  TELEMETRY — continuous high-freq measurement. Ephemeral, sampled into evt.
gov.<function>.<detail>    GOVERNANCE— meta-agent control plane: route/approve/escalate/audit/alert

STATE  → KV buckets ($KV.<bucket>.<key>) — current values, Watch-reactive
SYSTEM → RESERVED: $SYS $JS $SRV $KV $O — app code never touches these
```

**Why plane-first (not entity-first):**
- Retention classes map to top-level prefixes — `evt.>` is one immutable stream, `tel.>` ephemeral, they never tangle
- Entity-first scatters retention across the namespace
- You don't lose the entity view: `*.drone.skyhawk-1.>` = everything about that drone, every plane

---

### 5.2 The 15 node classes (FINAL — ratified 2026-06-05)

**The class filter:** A class is a thing that PERSISTS and has identity. A thing that HAPPENS is an event owned by a node — not its own class.

```
TIER 1 — BEINGS
  human       People — Zeus, Cielo, Liam, Noah, Soluna, students, contacts
  pet         Companion animals + livestock — Cerberus, future farm animals

TIER 2 — PLACES
  location    Structures & named places — house, barn, office, greenhouse, TOC, event sites
  eco         Land & living systems — property, gardens, beds, soil, crops

TIER 3 — MACHINES YOU CONTROL
  device      Controllable hardware endpoint — plugs, switches, locks, cameras, phones,
              radios, printers, Oura/Apple Watch (Rule 1), GPS receiver, Starlink dish
  sensor      Read-only emitter — temp, humidity, motion, presence, leak, soil, CT clamp,
              industrial pressure/flow/vibration. NEVER appears in cmd.>
  equipment   Physical systems with lifecycle — CNC, 3D printer, tractor, welder,
              irrigation system, workshop tools, greenhouse HVAC, medical gear.
              Distinct from device: has jobs, maintenance records, material consumption,
              operational hours, safety states.
  robot       Ground robots & manipulators — DIMOS, OM1, Unitree G1, Go2, FarmBot
  drone       Aerial vehicles — current UAVs + future fixed-wing autonomous
  vehicle     Transport — cars (2), motorcycle, plane (manned)

TIER 4 — COMPUTE
  agent       AI reasoning entities — Legacy (meta), cc, pi, openclaw, codex, gemini
  node        Compute hosts — MBP hub, VPS, HA box, spare, Pi, all leaf nodes
              (RENAMED from 'machine' — frees 'machine' for IIoT; 'node' = NATS node, accurate)
  service     SaaS / APIs / software capabilities — all 166 from 05 Stack

TIER 5 — CONTROL PLANE (pseudo-classes, not real entities)
  mode        Global state presets — incognito, high-alert, watch-me
  fleet       Broadcast groups — all-agents, all-drones, tactical-unit

HORIZON — reserved names, add when it hurts (not before)
  gateway     Protocol bridges — Zigbee coordinator, Z-Wave, MAVLink↔NATS, ROS2↔NATS, ONVIF
  spacecraft  Orbital/suborbital — cubesat, when real
  model       Edge AI models in lifecycle — when fleet model-swapping is real
```

**Class count: 15 active, 3 reserved on the horizon.**

---

### 5.3 The 3 gaps that were found + fixed (the audit that produced the final list)

**Gap 1 — CRITICAL: `machine` naming collision with IIoT**
In every industrial/barn context "machine" means physical machinery (CNC, tractor, lathe).
Having compute hosts named `machine` = total ambiguity the moment IIoT gets real.
**Fix: renamed `machine` → `node`.** Accurate (these ARE NATS nodes), frees "machine" from the namespace.

**Gap 2 — `equipment` is a real class (not just `device`)**
A `device` is a controllable endpoint (turn on, adjust, query). `equipment` is a system with a lifecycle:
jobs, maintenance records, material consumption, operational hours, safety states.
Without this class, IIoT and the physical estate (barn, workshop, greenhouse) lose all lifecycle semantics.
Examples: 3D printer, CNC, tractor, welder, irrigation system, greenhouse HVAC, medical gear.

**Gap 3 — `gateway` (named, deferred)**
Protocol bridge nodes (Zigbee coord, Z-Wave, MAVLink↔NATS, ROS2↔NATS) are NOT compute,
NOT devices, NOT sensors. They're infrastructure for other nodes to exist on the bus.
Named as a reserved horizon class. Built as a pattern on `node` until 3+ gateways exist.

**Satellites: NOT a class (and that's correct)**
GPS, Starlink, ADS-B = infrastructure you don't own. Your *receiver* is a `device`.
Its readings flow as `tel.vehicle.plane-1.adsb` or `tel.vehicle.car-1.location`.
The satellite is Starlink's node, not yours.

---

### 5.4 Three modeling rules (resolve 90% of edge cases)

**RULE 1 — Telemetry publishes under what it's ABOUT, not what measured it.**
Oura ring = `device`. Its reading = `tel.human.zeus.heartrate` (subject is you, ring is the source in payload).
Bouncie GPS = `device`. Its reading = `tel.vehicle.car-1.location`.
Resolves every wearable, tracker, and meter.

**RULE 2 — Sense vs. act is the security boundary.**
`sensor` = read-only, NEVER in `cmd.>`. `device` = takes commands.
When a thing does both (camera that streams AND pans), class as `device` (higher privilege).
Makes NKey ACLs trivial: a sensor's key = zero `cmd` permissions, full stop.

**RULE 3 — Identity in the subject, relationships in the registry.**
The barn's temp meter = `sensor.meter-7616`, NOT `sensor.barn.meter-7616`.
"It's in the barn / belongs to Zeus / is pillar HP" = registry facts, not subject tokens.
Move it to the office → subject never changes, just update its registry card.
**This is the single most important rule in the entire taxonomy.**

---

### 5.5 The universal node address + `id` naming law

```
<plane>.<class>.<id>.<path...>

id = lowercase-hyphen, NEVER contains dots (dots are NATS token separators)
     Singletons:  legacy, mbp-hub, zeus, coral, main-house
     Instances:   skyhawk-1, plug-a886, meter-7616, car-1
```

**The 05↔08 unification:** Every entry in 05 Stack (capability) and 08 Entity Registry (entity) shares this address. A drone is `drone.skyhawk-1` across ALL planes — `evt`, `cmd`, `tel`, `qry`, KV registry. One identity, everywhere.

**Class mapping back to vault structure:**
- 08 biological → `human`, `pet`
- 08 digital → `node`, `agent`, `service`
- 08 mechanical → `device`, `equipment`, `robot`, `vehicle`, `drone`
- 08 physical → `location`
- 08 ecological → `eco`
- 05 services → `service`

---

### 5.6 IoT / IIoT / AI coverage matrix

| Domain | Classes | Complete? |
|---|---|---|
| Consumer IoT | `device` + `sensor` | ✅ |
| Industrial IoT (IIoT) | `sensor` + `device` + `equipment` + `node` | ✅ — `equipment` closes the gap |
| Agricultural IoT | `sensor` + `equipment` + `eco` + `robot` | ✅ — FarmBot=`robot`, bed=`eco`, probe=`sensor`, irrigation=`equipment` |
| Surveillance / Security | `device` (cameras) + `sensor` (motion) + `location` | ✅ |
| Automotive / Vehicle IoT | `vehicle` + `device` (OBD/GPS) + `sensor` | ✅ |
| Aviation IoT | `drone` (UAV) + `vehicle` (manned) + `device` (transponder, avionics) | ✅ |
| Robotics | `robot` + `drone` + `gateway` (ROS2 bridge, horizon) | ✅ |
| Edge AI | `agent` (reasoning) + `service` (model endpoint) on a `node` | ✅ — `model` splits when needed |
| Satellite / GPS | NOT your node — receiver=`device`, data=`tel` | ✅ correct by design |
| Tactical / ATAK / TOC | `human` + `device` (Meshtastic) + `location` (TOC) + `fleet` + `mode` | ✅ |
| Smart Home | `device` + `sensor` + `location` + `service` (HA) | ✅ |
| Farm / Greenhouse | `eco` + `equipment` + `sensor` + `robot` | ✅ |

---

### 5.7 Worked examples

```bash
# ONE DRONE — addressed across all planes (one identity, many planes)
cmd.drone.skyhawk-1.land
tel.drone.skyhawk-1.gps
evt.drone.skyhawk-1.landed
qry.drone.skyhawk-1.battery
$KV.node-registry.drone.skyhawk-1          # registry card + live state

# THE BARN — everything about one location
evt.location.barn.motion_detected
cmd.device.barn-door-1.open
cmd.equipment.air-compressor-1.power_on
tel.location.barn.temperature
sensor.barn-temp-1 → publishes tel.sensor.barn-temp-1.reading (Rule 1: about the sensor itself)
gov.approve.cmd.device.barn-door-1.open    # meta-agent sign-off required at night

# ZEUS — the human node
evt.human.zeus.left_home
tel.human.zeus.heartrate                   # Oura ring (Rule 1 — about Zeus, not device)
tel.human.zeus.location                    # iPhone GPS (Rule 1)
evt.human.zeus.sleep.ended

# SAAS (service class — was svc.todoist.*)
cmd.service.todoist.create_task
qry.service.github.list_issues
evt.service.stripe.payment_succeeded
evt.service.fastmail.email_received

# EQUIPMENT (lifecycle — the new class)
evt.equipment.bambu-x1.print_completed     { job: "part-42", duration_min: 94, material_g: 127 }
tel.equipment.tractor-1.hours              { hours: 342.7 }
cmd.equipment.irrigation-a.start_zone      { zone: "beds-north", duration_min: 20 }
evt.equipment.cnc-mill-1.job_started

# GOVERNANCE plane
gov.route.incoming_task
gov.approve.cmd.device.barn-door-1.open
gov.escalate.anomaly.agent-cc
gov.audit.cmd.device.barn-door-1.open
gov.alert.security.motion_barn_3am
```

---

### 5.8 Power subscriptions (proof the design works)

```bash
evt.>                        # record EVERYTHING into experience store (ClickHouse)
gov.escalate.>               # wake the cortex (Legacy) — only what needs a mind
tel.drone.>                  # all drone telemetry (flight dashboard)
cmd.drone.>                  # ALL commands to ANY drone — one gov rule covers the fleet
evt.sensor.>                 # every sensor event across every class
*.location.barn.>            # EVERYTHING about the barn — any plane, that node
*.human.zeus.>               # everything about Zeus — events, telemetry, commands
cmd.>                        # EVERY mutation in the digital life — the governance firehose
evt.service.>.payment*       # all payment events across all services
```

`*.<class>.<id>.>` = "everything about a node, all planes." One wildcard, total visibility.

---

### 5.9 The taxonomy IS the security model

Per-node NKey least-privilege is now a one-liner. Example ACLs:

```yaml
# Drone node — can only publish its own telemetry/events, subscribe to its own commands
drone-skyhawk-1:
  publish:   ["tel.drone.skyhawk-1.>", "evt.drone.skyhawk-1.>"]
  subscribe: ["cmd.drone.skyhawk-1.>"]

# Sensor node — read-only, zero command publish
sensor-barn-temp-1:
  publish:   ["tel.sensor.barn-temp-1.>", "evt.sensor.barn-temp-1.>"]
  subscribe: []   # sensors receive nothing

# Legacy (meta-agent) — full bus access
agent-legacy:
  publish:   [">"]
  subscribe: [">"]

# Read-only observer (L1 guest)
observer:
  publish:   []
  subscribe: ["evt.>", "tel.>"]   # can watch, cannot command or govern
```

A clean taxonomy makes cryptographic least-privilege mechanical, not bespoke.

---

### 5.10 Migration map (current sprawl → target spine)

| Current subject | Target subject | Notes |
|---|---|---|
| `agent.>` | `evt.agent.>` | rename |
| `agents.prompt.>` | `cmd.agent.<id>.prompt` | ACP over NATS |
| `events.>` | `evt.>` | rename — ClickHouse engine re-pointed |
| `saas.>` / `webhooks.>` | `evt.service.<service>.<event>` | normalize |
| `system.>` | `tel.node.<id>.*` + `evt.node.<id>.*` | node replaces machine |
| `health.>` | `tel.human.<id>.*` / `evt.human.<id>.health*` | |
| `security.>` | `evt.location.<id>.*` / `evt.sensor.<id>.*` | |
| `cmd.ha.*` | `svc.ha.*` → `cmd.service.ha.*` + `cmd.device.<id>.*` | ha = service |
| `cmd.notify` | `cmd.service.notify.send` | |
| `cmd.coral.query` | `qry.service.coral.sql` | |
| `cmd.tts.speak` | `cmd.service.tts.speak` | |
| `cmd.mode.*` | `cmd.mode.<name>` | stays — pseudo-class |
| `cmd.fleet.*` | `cmd.fleet.<group>.*` | stays — pseudo-class |

**Migration mechanism: NATS server-side `mappings {}`** — aliases old→new during cutover.
Nothing breaks while services migrate one at a time. Zero-downtime rename.

---

## 6. Reconciled Build Sequence

The original 6 phases are wrong for where things actually are. The real order:

1. **FOUNDATION — Identity & Accounts.** Single `$G`/password → decentralized JWT + NKeys (`nsc`). Per-agent cryptographic identity. Configure `$SYS`. *(Synadia's "Agent Identity," done natively — no paid product.)* Without it, governance is theater.
2. **SPINE — Ratify the subject taxonomy (§5)** and apply via `mappings {}`.
3. **REGISTRY — Node registry as KV** over the 08 markdown source of truth.
4. **REGISTRATION — Every node a `$SRV` micro-service.** `nats micro ls` = the whole fleet, live. Today only OpenClaw does this.
5. **THE BRAIN — Autonomic governor** (`subscriber.py` → router/policy engine) **+ escalation to Legacy.**
6. **GOVERNANCE — hand-rolled subject ACLs** first; evaluate Synadia Protect only if scale ever justifies a paid gateway.
7. **EXPERIENCE INTELLIGENCE** — pattern detection + replay-onboarding on the ClickHouse store that already exists.
8. **REDUNDANCY — Mirror now → 3-node cluster** when load-bearing (§3).
9. **NEX — last, optional.** Only when multi-node scheduling pain is real.

---

## 7. Sovereignty Principle (set once, stop re-litigating)

Synadia pivoted to exactly this vision and now sells the pieces: Cloud (multi-tenant hosted — violates locality), Platform (managed — external control), Protect + Insights (paid gateways). **Adopt the open protocol (NKeys/JWT, `$SRV`, NEX OSS) — the same 20MB binary already running — and hand-roll the governance/observability they sell.** Doctrine: *"Nothing leaves unless you say so."*

---

## 8. Integration Architecture — How Everything Connects

> The core principle: nothing talks to anything directly. Everything translates TO the bus. Once on the bus, everything can hear everything — regardless of language or protocol.

### 8.1 The 4 patterns (every integration is one of these)

**Pattern 1 — Native (zero bridge code)**
Three protocols are built into the NATS binary already running:
- MQTT :1883 → ESP32s, Arduino, most IoT sensors. MQTT topic slashes auto-become NATS dots. Server-side `mappings {}` normalizes to taxonomy.
- WebSocket :8080 → browsers, Flutter web, JS apps, onemind-command
- Plain NATS :4222 → anything using an SDK

**Pattern 2 — Protocol Bridge (one thin process per protocol family)**
One bridge serves all devices of that type. Bridge = translate only, never decide.
Registers as `$SRV` gateway node. Built once, serves N devices.

| Protocol | What uses it | Bridge | Language |
|---|---|---|---|
| Zigbee / Z-Wave | Smart home | HA → ha-bridge (already done) | — |
| MAVLink | Drones, plane | pymavlink + nats-py (companion computer) | Python |
| ROS2 / DDS | DIMOS, OM1, G1, Go2 | ros2-nats (on robot compute) | Python/Go |
| Modbus | Industrial PLCs, barn equipment | Node-RED Modbus node → NATS | Node-RED |
| OPC-UA | SCADA, heavy industrial | Node-RED OPC-UA node → NATS | Node-RED |
| BLE | Tile trackers, health sensors | Pi/HA BLE scanner → NATS | Python |
| ONVIF | IP cameras | HA camera integration → ha-bridge | — |
| CAN bus | OBD, some robots | python-can + nats-py | Python |
| Serial / GPIO | RPi GPIO, Arduino | pyserial + nats-py | Python |
| ADS-B | Plane transponder | dump1090 → Python → NATS | Python |
| Meshtastic | Off-grid radio mesh | meshtastic-python + nats-py | Python |

**Pattern 3 — Language SDK (for custom services and agents)**
Every language has an official NATS client. Language is irrelevant once on the bus.

| Language | SDK | Use |
|---|---|---|
| Python | nats-py | Bridges, governor, ML, sensor readers |
| JavaScript/TypeScript | nats.js + nats.ws | Node services, n8n, webhook handlers |
| Dart | dart_nats | onemind-command Flutter app |
| Swift | swift-nats | iOS native apps |
| Go | nats.go | Performance-critical bridges, fleet daemons |
| Rust | async-nats | High-perf sensor gateways, drone compute |

**Pattern 4 — SaaS Wrapper (config-driven, one process, N config files)**
For all 166 services in 05 Stack. NOT 166 daemons — ONE generic host reading YAML configs.
Config defines: subjects, REST endpoints, auth reference (1Password), body mapping, webhooks.
Add a new service = add a new YAML file, zero new code.
This is the evolved `saas-action-worker` — same pattern, made config-driven and schema-consistent.

---

### 8.2 The standard message envelope (solves "different languages = different payload shapes")

Every node, every language, every protocol bridge ALWAYS wraps payload in:

```json
{
  "id":      "01HZ8K...",
  "ts":      "2026-06-05T20:00:00.000Z",
  "source":  "service.todoist",
  "v":       1,
  "payload": { ...the actual data... }
}
```

Four fields, always present, always JSON. Any language deserializes the same envelope before reading payload. Language-agnostic interop without a schema registry.

**Exception — high-frequency telemetry:** `tel.>` subjects may use MessagePack (binary, ~30% smaller, same structure). Declared via NATS header `Content-Type: application/msgpack`. Consumers check header, pick deserializer. Everything else stays JSON.

---

### 8.3 The full integration stack

```
TIER 0 — Built into NATS server (zero code)
  MQTT bridge      ← ESP32s, IoT sensors, anything cheap
  WebSocket        ← browsers, Flutter, JS
  Leaf nodes       ← remote machines extend the same bus

TIER 1 — Protocol bridges (one per family, thin, config-driven)
  ha-bridge        ← Zigbee, Z-Wave, ONVIF, BLE via HA (LIVE)
  mavlink-bridge   ← drones, plane avionics
  ros2-bridge      ← DIMOS, OM1, G1, Go2
  modbus-bridge    ← industrial equipment (Node-RED)
  serial-bridge    ← Arduino, GPIO, raw serial

TIER 2 — SaaS wrapper (one process, N config files)
  service-wrapper  ← all 166 from 05 Stack, config-driven YAML

TIER 3 — Native NATS (SDK, direct connection)
  Agents           ← Legacy, cc, pi, openclaw
  Autonomic layer  ← governor daemon
  Flutter app      ← onemind-command
  Custom scripts   ← any language

TIER 4 — Server-side (no process at all)
  mappings {}      ← rename/normalize subjects without a bridge
  $SYS             ← server metrics + account events, built-in
```

---

### 8.4 Decision tree — how to add any new thing

```
Does it speak MQTT?            → connect to :1883 directly. Done.
Has a REST API / webhooks?     → add YAML to service-wrapper. Done.
Already goes through HA?       → ha-bridge already handles it. Done.
Speaks MAVLink?                → mavlink-bridge handles it. Done.
Speaks ROS2/DDS?               → ros2-bridge handles it. Done.
New protocol family?           → write ONE bridge (~100-200 lines), register as gateway.$proto.$id
Writing from scratch?          → import NATS SDK for your language, publish/subscribe correct subjects.
```

New code only for new protocol families. Your world has ~8-10 total, most already covered.

---

## 9. MCP vs NATS SaaS Wrapper — Relationship + Decision

> Sharp question with a sharp answer: they are NOT competing. They are different layers solving different problems. The smartest move is to use MCP as the implementation backend FOR the NATS wrapper.

### 9.1 What each one actually does

| | MCP | NATS SaaS Wrapper |
|---|---|---|
| **Designed for** | AI agents calling tools | ANY node calling services |
| **Callers** | LLMs only | Agents, automations, scripts, devices, Flutter app, governor |
| **Protocol** | MCP (stdio / HTTP SSE) | NATS subjects |
| **Schema** | Rich tool schema (LLM understands natively) | YAML config + JSON envelope |
| **Today in stack** | LIVE — Todoist MCP, Fastmail MCP, agentmemory (53 tools) | Partial — saas-action-worker not yet config-driven |

MCP answers: **"How does an LLM invoke a capability?"**
NATS SaaS wrapper answers: **"How does ANY node on the bus invoke a capability?"**

The autonomic governor, Flutter app, drone landing sequence, automation script — none are LLMs, none can use MCP natively. They publish `cmd.service.todoist.create_task` and the wrapper handles it.

### 9.2 The smart move — MCP servers ARE the implementation backend

Don't reimplement the REST layer twice. The wrapper calls the MCP server (or Composio CLI) that already handles the API complexity:

```
NATS Bus
  cmd.service.todoist.create_task
          ↓
  SaaS Wrapper (thin NATS↔MCP bridge)
          ↓
  MCP Server (Todoist) ← already handles auth, schema, REST
          ↓
  Todoist API
```

Priority order for wrapper backend:
```
1. MCP server exists?       → call MCP
2. Composio action exists?  → composio execute
3. Service CLI exists?      → shell exec (gh, himalaya, etc.)
4. REST API only?           → direct HTTP (last resort)
```

This turns "build 166 REST integrations" into "bridge 3 existing systems onto the bus."

### 9.3 The three-layer stack

```
LAYER 1 — IMPLEMENTATION (actual API calls)
  MCP Servers     ← Todoist, Fastmail, agentmemory, any MCP-native tool (LIVE)
  Composio CLI    ← 1000+ app actions (LIVE)
  Direct REST     ← fallback only

LAYER 2 — BUS BRIDGE (the SaaS wrapper)
  One process, subscribes to cmd.service.> and qry.service.>
  Translates NATS envelope → MCP call or Composio call or REST
  Publishes results as evt.service.> responses
  Registers as $SRV micro-service per service name

LAYER 3 — CALLERS
  LLM agents (Legacy, cc, pi)  → MCP natively OR via NATS
  Automations / governor       → NATS only
  Flutter app                  → NATS only
  Any future node              → NATS only
```

### 9.4 The CLI question

- **Composio CLI** → backend for the wrapper (shell exec for 1000+ apps). Not the bus interface.
- **Service CLIs** (`gh`, `himalaya`, etc.) → shell exec fallback from wrapper. Same pattern.
- **`nats` CLI** → ops/debug tool for the bus itself. Not a service backend.

### 9.5 Dual-path decision (ratified)

**Legacy uses MCP natively for interactive tool calls. Everything non-LLM goes via bus.**

Option A (dual path — adopted): Legacy's interactive MCP calls are real-time and user-facing. Extra bus hop adds latency for zero benefit. Automations, governor, app → bus always.

Option B (bus-first — deferred): Legacy publishes `cmd.service.*` like every other node. Pure but adds latency to interactive calls. Revisit if unified audit trail becomes critical.

Audit trail still works either way: the wrapper emits `evt.service.todoist.task_created` regardless of which path triggered it.

---

## 10. Management, Visibility + Control Plane

> The question: how do we see everything and manage all of it? Answer: 4 layers, built bottom-up. But first — the visibility problem is actually a naming problem.

### 10.1 The live server reality (probed 2026-06-05)

```
Connections: 13
  ha-bridge         nats.js    ← named ✅
  testui            nats.js    ← named ✅
  (unnamed) ×6      python3    ← anonymous ❌
  (unnamed) ×2      C          ← anonymous ❌
  (unnamed) ×2      unknown    ← anonymous ❌

Streams: 11
  ZEUS_EVENTS       2 messages  ← only live stream
  SAAS_EVENTS       8 messages  ← last activity 12 days ago
  AGENT_EVENTS, COMMANDS, COMMS, HEALTH_EVENTS,
  SECURITY_EVENTS, SYSTEM, TASKS ← ALL EMPTY, never used

KV Buckets: 8
  rules             13 values  (most alive)
  zeus-context      5 values
  config/agent-state/notifications ← sparse
  connections/schedules/triggers   ← ALL EMPTY

$SRV micro-services: ZERO
Leaf nodes: ZERO
```

**A beautiful dashboard of this is useless.** Anonymous connections, empty streams, zero registered services.
The visibility problem is a naming problem. Fix naming first, then any interface becomes meaningful.

### 10.2 The 3 fixes that make visibility work (cheapest wins in the project)

**Fix 1 — Name every connection (one parameter, one line per service)**
```python
nc = await nats.connect(servers=["nats://localhost:4222"], name="saas-worker-1")
```
Cost: 30 minutes across all Python workers. Payoff: connz shows exactly what's connected.

**Fix 2 — Register every service as $SRV micro-service**
Every node announces itself. `nats micro ls` = the whole live fleet.
```python
# ~20 lines per service, done once
srv = await nc.services.add(name="saas-worker", version="1.0.0", description="SaaS action worker")
```

**Fix 3 — Configure $SYS account**
Unlocks server-level event streams:
```
$SYS.SERVER.ACCOUNT.*.CONNECT     ← every connection event with full client detail
$SYS.SERVER.ACCOUNT.*.DISCONNECT  ← every disconnect
$SYS.SERVER.>.STATSZ              ← per-account message rates
```
Cost: 1 config block in nats-server.conf. Payoff: automatic audit trail of every node joining/leaving.

Do these 3 and the entire fleet is visible before building a single UI.

### 10.3 The 4 visibility layers

```
LAYER 4 — APPLICATION (what is the fleet doing?)
  onemind-command Flutter panel ← the "everything" view, long-term

LAYER 3 — EVENTS (what happened?)
  ClickHouse + evt.> firehose ← historical + live, already live

LAYER 2 — INFRASTRUCTURE METRICS (is the bus healthy?)
  $SYS + Surveyor → Prometheus → Netdata

LAYER 1 — OPS / REAL-TIME (immediate, zero install)
  nats CLI + HTTP :8222 ← available today
```

### 10.4 Layer 1 — What you have today (zero install)

```bash
nats sub ">"                      # see EVERYTHING live on the bus
nats stream ls && nats stream report   # stream health
nats stream view SAAS_EVENTS      # browse actual messages
nats kv ls && nats kv get rules <key>  # KV explorer
nats server report connections    # all connections with detail
nats micro ls                     # fleet (useful after $SRV fix)
nats micro info <service>
nats sub "evt.>"                  # all events
nats sub "cmd.>"                  # all commands (governance audit)
```

HTTP monitoring at :8222: varz, connz, jsz, leafz — raw JSON, complete, always available.

### 10.5 Layer 2 — Infrastructure metrics (install now, ~1hr, fully sovereign)

**NUI — visual NATS desktop explorer (install today)**
```bash
# https://github.com/nats-nui/nui — open source, local, free
# Download from GitHub releases or brew if tap available
```
Gives: visual stream browser, KV explorer, connection list, live subject monitor, consumer management.
Runs 100% local, connects to :4222. Sovereign.

**nats-surveyor → Netdata (you already have Netdata at :19999)**
```bash
go install github.com/nats-io/nats-surveyor@latest
# Configure Netdata to scrape Prometheus endpoint
# Result: message rates, connection counts, stream sizes on existing Netdata dashboard
```

### 10.6 Layer 3 — Events (already live)

ClickHouse `events.>` store with 5yr retention. Query with SQL.
After taxonomy rename: `evt.>` flows here automatically from every node.
For live view: subscribe to `evt.>` in NUI or CLI. No new tooling.

### 10.7 Layer 4 — The onemind-command Fabric View (build)

Yes, build a NATS interface. It lives inside onemind-command as a Fabric View panel — not standalone.
Flutter connects via WebSocket to NATS :8080 (already live) using dart_nats. No new server.

```
┌─────────────────────────────────────────────────────┐
│  ONEMIND COMMAND — FABRIC VIEW                      │
├──────────────────┬──────────────────────────────────┤
│  FLEET ($SRV)    │  LIVE EVENT FEED (evt.>)         │
│  ○ legacy    🟢  │  evt.human.zeus.left_home        │
│  ○ openclaw  🟢  │  evt.service.github.pr.opened    │
│  ○ ha-bridge 🟢  │  cmd.device.barn-door-1.open     │
│  ○ cc        🔴  │  gov.escalate.anomaly.cc         │
├──────────────────┼──────────────────────────────────┤
│  NODE REGISTRY   │  STREAMS                         │
│  (KV Watch)      │  evt.>   LIVE                    │
│  drone.skyhawk 🟢│  cmd.>   LIVE                    │
│  robot.om1     🔴│  tel.>   LIVE                    │
├──────────────────┼──────────────────────────────────┤
│  COMMAND DISPATCH│  GOVERNANCE                      │
│  subject: ______ │  pending approvals: 2            │
│  payload: ______ │  escalations: 0                  │
│  [PUBLISH]       │  [REVIEW]                        │
└──────────────────┴──────────────────────────────────┘
```

**How each panel works:**
- Fleet panel → `$SRV.*` discovery + `$KV.node-registry.*` Watch
- Live event feed → WebSocket :8080, subscribe `evt.>`
- Node registry → KV Watch on node-registry bucket (live, reactive)
- Command dispatch → publish to bus with governance pre-check
- Governance panel → subscribe to `gov.approve.>` pending queue

### 10.8 MCP + Composio as service backends (not competing — stacking)

MCP answers: "how does an LLM call a service?"
NATS SaaS wrapper answers: "how does ANY node call a service?"

The wrapper uses MCP/Composio as its implementation backend — no double REST implementation:
```
1. MCP server exists?       → call MCP     (Todoist, Fastmail, agentmemory — LIVE)
2. Composio action exists?  → composio execute  (1000+ apps — LIVE)
3. Service CLI exists?      → shell exec   (gh, himalaya, etc.)
4. REST only?               → direct HTTP  (last resort)
```

Legacy (LLM) keeps its direct MCP path for interactive calls — the extra bus hop adds latency for zero benefit. Everything non-LLM goes via bus. Audit trail works either way: wrapper always emits `evt.service.<name>.*`.

### 10.9 Full management stack — what to use for what

| Need | Tool | When |
|---|---|---|
| See the bus live | `nats sub ">"` CLI | Today |
| Browse streams/KV visually | NUI desktop app | Install today |
| Server health metrics | $SYS + Surveyor + Netdata | This week |
| Full fleet status | `nats micro ls` | After $SRV registration |
| Historical events | ClickHouse SQL | Already live |
| Live subject explorer | NUI | Today |
| Command dispatch + governance | onemind-command panel | Build |
| "See everything" | onemind-command Fabric View | Build |

### 10.10 The sequencing that makes this work

```
WEEK 1: Name connections + configure $SYS  → visibility unlocked
WEEK 2: Register $SRV per service           → fleet panel populated
WEEK 3: NUI + Surveyor installed            → visual + metrics
WEEK 4: Subject taxonomy rename             → clean evt.>/cmd.>/tel.> flows
BUILD:  onemind-command Fabric View panel   → the real control plane
```

---

## 11. Open Decisions / Next Up

- [ ] Ratify the subject taxonomy (§5) — the keystone. Refine class list + kind list together.
- [ ] Confirm identity-first sequencing (rolling cutover plan for the live `$G` services).
- [ ] Verify HA box specs → decide its cluster role (data node vs tiebreaker).
- [ ] Decide canonical firehose rename `events.>` → `evt.>` timing.
- [ ] Pick the first node class to fully migrate to the new schema as the reference implementation.

---

*One Mind Framework — NATS OS Initiative. Zeus Delacruz + Legacy, 2026-06-05.*
*Companion to `ONEMIND-NATS-OS.md`. Working doc — evolves as decisions are ratified.*
*Cognitive layer stack (L0–L5) + world-model layers L1/L2 subject contract: `ONEMIND-LAYER-STACK.md` (2026-06-07).*
