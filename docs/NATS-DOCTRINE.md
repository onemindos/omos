# NATS-DOCTRINE.md — OneMind NATS Canon

> **Authority:** This document is the binding NATS doctrine for all OneMind OS work.
> Worker agents MUST follow these rules. Deviations require Legacy (architect) approval.
> Maintained by: Legacy (default profile). Last updated: 2026-06-10.

---

## 1. Why NATS Is the Spine

One substrate replaces six: pub/sub, request/reply, durable streams (JetStream), KV store,
object store, and service discovery (`$SRV` micro). No Kafka, no RabbitMQ, no Redis pub/sub,
no separate service mesh. NATS runs the same binary on the VPS, a Mac, and a Pi Zero —
which is exactly the OneMind fleet shape (House / Barn / Office nodes, target: hundreds).

**Decision (locked):** NATS is the nervous system. Every service, agent, app, and sensor
speaks through it. Anything that can't speak NATS gets a bridge (see §7).

---

## 2. Subject Grammar (LOCKED)

External standard — Synadia Agent Protocol v0.3 (verb-first):

```
agents.{verb}.{agent}.{owner}.{session}    verb ∈ prompt | hb | status | attachments
$SRV.PING.agents / $SRV.INFO.agents        discovery — callers MUST read endpoint
                                           subjects from here, never hardcode
```

OneMind internal grammar:

```
evt.{domain}.{entity}.{event}      events (facts, past tense)    evt.so.task.created
cmd.{domain}.{entity}.{action}     commands (imperatives)        cmd.le.light.on
tel.{node}.{class}.{metric}        telemetry                     tel.barn-pi-1.temp.c
svc.{service}.{endpoint}           internal service RPC          svc.meta.task.list
```

Rules:
- Domains map to the vault pillars: `so` `hp` `le` `ge` (+ `sys` for infrastructure).
- Design subjects so **wildcard filters do the work** (`evt.so.>`, `tel.barn-pi-1.>`).
  A consumer needing >300 disjoint filters means the hierarchy is wrong — redesign it.
- Events are immutable facts. Commands may be rejected. Never publish a command as an event.

---

## 3. JetStream Rules (from Synadia production guidance)

1. **`consumer create` is idempotent — use it directly.** Never call `consumer info` to
   check-then-create. Same name/config → success no-op; changed config → update; absent →
   create. `consumer info` hits the meta-leader and computes state — monitoring only.
2. **Read pending counts from message metadata** (`NumPending` on the last fetched message),
   not from `consumer info`.
3. **Consumer budget: stay far under 100K total consumers** (cluster instability risk).
   Homestead scale should never exceed hundreds. If consumer count grows with device count,
   the design is wrong — use republish or Direct Get instead.
4. **Republish over filtered consumers** for real-time fan-out without delivery guarantees:
   stream republishes `orders` → `filtered.orders.>`, clients use plain core subscriptions.
5. **Direct Get / Batch Get** for "latest value" reads (mobile/IoT pattern):
   `$JS.API.DIRECT.GET.<stream>.<last_by_subj>`. No consumer needed.
6. **≤300 disjoint subject filters per consumer.** Beyond that, split consumers or fix the
   subject hierarchy.
7. **Flow control is an architectural decision, not a tuning knob.** Decide per-stream what
   happens under backpressure (drop, buffer, summarize) at design time. No firehoses.
8. Useful primitives to reach for before building anything custom: per-message TTL
   (ephemeral state), delayed message scheduling, atomic batch publish (2.12+),
   distributed CRDT counters, expected-sequence headers (optimistic concurrency).

---

## 4. Edge ↔ Core Streaming Patterns

Three patterns, three different jobs. Mature systems use all three simultaneously.

| Pattern | What | When | Caveats |
|---------|------|------|---------|
| **Consume** | Core pulls selectively from edge streams | Core needs a *subset* (alerts, threshold breaches) | Couples to edge reachability; filter at ingest boundary, not downstream |
| **Mirror** | Exact 1:1 read-only replica, sequence numbers preserved, auto-catchup after disconnection | Compliance/audit trail, backup/DR, read offload | Full copy = storage cost. Don't default to mirror because it "feels safe" — that's the firehose anti-pattern |
| **Merge (source)** | Many edge streams → one core stream, subject transforms normalize taxonomy in topology | Cross-site analytics (cross-node correlation) | **No global ordering** — per-source order only. Origin sequence numbers not preserved |

OneMind application: each fabric node (house, barn, office) runs a leaf node with local
JetStream. Node-local telemetry stays local; core merges `tel.>` for fleet analytics;
critical event streams mirror to core for the permanent record; Legacy consumes selectively
for attention-worthy signals.

**Edge resilience:** store-and-forward at the leaf, not retry loops in clients.
"Just retry" fails edge systems — the leaf node buffers during disconnection and
JetStream catch-up handles reconciliation.

---

## 5. Security Model

- **Accounts are security domains** (DDD boundary, not an org-chart artifact). Cross-account
  sharing is explicit: **stream exports** (read-only data) and **service exports** (RPC).
- **Two management models:** decentralized JWT (operator → account → user; scales, key
  rotation, bulk invalidation via signing keys) vs config-based (simple, requires
  config sync + restart across cluster). Hybrid allowed: operator-generated nkeys in config.
- **Auth callout** bridges to external auth (OAuth/LDAP/custom service). Mixed mode is legal
  across accounts, but never mix callout + JWT **within one account**.
- **Full resolver** is built into every server — JWTs sync mesh-style, no external component.
- **Leaf nodes authenticate independently** — each node can run its own auth model; hub
  still authorizes cross-account imports.
- Credential rotation disconnects clients by design; reconnect triggers fresh auth.
- Publisher identity is not a core NATS concept — embed `source_id` in payloads where
  audit matters.

OneMind staging: start config-based (single operator household, low churn), with subjects
and accounts designed so the cutover to operator/JWT mode is a config change, not a
redesign. Cutover trigger: fleet > ~20 nodes or first non-family device class.

Planned accounts:
```
SYS         system account (monitoring — required for Insights-equiv)
CORE        agents + meta services (Legacy, workers, meta-layer)
FABRIC      telemetry + IoT (sensors, bridges, robots)
INTERFACE   apps (Flutter, React cockpit) — imports from CORE/FABRIC via exports
```

---

## 6. Client Policy (Orbit Philosophy)

Synadia's stance, adopted: **core clients stay minimal/stable; opinionated patterns live in
extension libraries (Orbit).** Consequences for OneMind:

- Python services: `nats-py` + `orbit.py` extensions where useful.
- React cockpit: `nats.ws` + `orbit.js`.
- Flutter: `dart_nats` (community — no official Dart client, no orbit.dart). Treat it as a
  thin core client; implement needed patterns (request-many, sentinel) app-side in one
  `nats_service.dart`, mirroring Orbit semantics so future swaps are mechanical.
- **Request-many patterns** (multiple responders / chunked responses): terminate by
  max-messages (known count), sentinel (empty message), or stall timeout (discovery —
  "all responders that exist"). Use stall-timeout for agent discovery sweeps.

---

## 7. Synadia Product Map — Reuse vs Replicate

Synadia sells polish on public APIs. Their open engines we use as-is; their closed
binaries we replicate against the same public interfaces.

| Synadia | Status | OneMind move |
|---------|--------|--------------|
| **NEX / Workloads** | ✅ OSS (Go) | **Use as-is.** Deploy services to any fleet node over the bus (`nex node up` per machine). Workloads get scoped NATS creds automatically |
| **Connect / Wombat runtime** | ✅ OSS | **Use as-is.** Bridge engine for HA, webhooks, external APIs ↔ bus |
| **Orbit libraries** (7 langs) | ✅ OSS | Use `orbit.py`, `orbit.js` directly |
| **callout.go** | ✅ OSS | Auth callout when we harden (Phase 7) |
| **Agent SDK** (`synadia-ai/synadia-agents`) | ✅ OSS | **The agent layer.** Caller + host SDKs (Py/TS), prebuilt harness plugins incl. Hermes (WIP upstream), Claude Code, OpenClaw |
| **Control Plane** | ❌ closed | Replicate: `control-plane-sdk-go` exposes their full OpenAPI spec = free blueprint. Ours = Python service on `$SYS` + app Command tab |
| **Insights** | ❌ closed | Replicate: `nats-surveyor` + `prometheus-nats-exporter` (both OSS, nats-io) → SQLite history ("time travel" = SQL) → app Analytics tab. Their ~100 auto-checks = our reimplementation checklist |
| **Protect** (NATS WAF) | ❌ closed | Defer. Single-operator homestead doesn't need a protocol WAF on day one. When needed: auth callout policy at connect-time first; Expr-style rule proxy later |
| **HTTP Gateway** | ❌ closed | Trivial shim (or Wombat handles it) |
| **Schema Registry** | ❌ closed | JSON Schema files in a NATS KV bucket. That is all it is |
| **Cloud / NGS** | SaaS | Not used — sovereignty requirement. Self-hosted only |

---

## 8. Agent Layer Doctrine (Zeus directive — LOCKED)

> **The harness owns its config. The app never does.**

- Each harness (Hermes, Claude Code, OpenClaw, PI, Codex, …) keeps its own config in its
  own native format/location.
- The OneMind app is a **review + edit surface only**: read config through the agent's own
  NATS endpoints, write changes back through the same endpoints, natively into the
  harness's settings. The Agent SDK plugin at each harness translates.
- **No shadow configs. No app-side agent database. No sync conflicts.** Discovery comes
  from `$SRV.PING.agents` sweeps (stall-timeout request-many), never a hardcoded roster.
- Meta-agent topology per Derek Collison's thesis: heterogeneous harnesses, one secure
  substrate; meta-agent (Legacy) coordinates; A2A and realtime agent events arrive in the
  SDK roadmap — we adopt them when released, we don't pre-build them.

---

## 9. Geo-Envelope Standard (TACMAP — LOCKED)

> **The geo-envelope is the contract that turns the NATS bus into a Common Operating Picture.**
> Any actor or event that carries a position becomes renderable on TACMAP automatically —
> no map code changes when a new asset class arrives. Every device, robot, aircraft, vessel,
> conflict event, or world-feed entry that publishes to the geo subject space participates.

### 9.1 Subject Taxonomy

```
# ── INNER RING — your fleet (participatory — you ARE on the map) ──────────────
tel.fleet.{class}.{id}.position       live position of a fleet asset
tel.fleet.{class}.{id}.status         battery, mode, link-health, sensor state
cmd.fleet.{class}.{id}.{verb}         dispatch to an asset  verb ∈ goto|estop|rtl|…
cmd.mission.{id}.{verb}               mission lifecycle     verb ∈ create|waypoint|abort|complete

# ── OUTER RING — world intel (observational — the world around you) ────────────
evt.world.air.{id}                    aircraft (ADS-B, commercial + military)
evt.world.air.mil.{id}                military aircraft (surge-detected)
evt.world.sea.{id}                    vessels (AIS)
evt.world.conflict.{id}               conflict events (ACLED)
evt.world.military.{id}               military base positions
evt.world.fire.{id}                   wildfire perimeters (NASA FIRMS)
evt.world.seismic.{id}                earthquake events (USGS)
evt.world.chokepoint.{id}             strategic waterway status (Suez, Malacca, Hormuz…)
evt.world.news.{country}.{id}         geo-tagged news/intelligence events

# ── TAK BRIDGE — field devices (ATAK/iTAK/WinTAK clients) ────────────────────
evt.tak.cot.{uid}                     Cursor-on-Target events from TAK clients
evt.tak.status.{uid}                  TAK client connection state

# ── MAP CONTROL — agent/app interaction ───────────────────────────────────────
cmd.map.filter.set                    filter active layers / geographic focus
cmd.map.highlight.{uid}               highlight a specific actor on the map
cmd.map.geofence.{id}.{verb}          draw / remove a geofence
cmd.map.viewport.set                  center / zoom the map view
cmd.map.annotation.{id}.{verb}        drop / update / remove a map pin or note
cmd.map.mission.draw                  render a mission route on the map
evt.map.geofence.{id}.{enter|exit}    Watch Job fires on spatial crossing
```

Asset class vocabulary for `tel.fleet.{class}`:
`drone | robot | vehicle | node | edge | agent | aircraft | vessel`

### 9.2 Message Shape (every geo subject carries this envelope)

```json
{
  "uid":     "icao-A1B2C3",           // stable, globally unique identity
  "type":    "air.commercial",        // dot-path type (see §9.3)
  "lat":     40.7128,                 // WGS-84 decimal degrees
  "lon":     -74.0060,
  "hae":     10972,                   // height above ellipsoid (metres); 0 for ground assets
  "course":  270,                     // degrees true (0–360); -1 if stationary/unknown
  "speed":   450,                     // m/s; -1 if unknown
  "callsign": "UA123",                // human-readable name / callsign
  "stale":   30,                      // seconds until this reading is considered stale
  "t":       1718150400000,           // Unix epoch milliseconds (source timestamp)
  "meta":    {}                       // domain-specific payload (see §9.4)
}
```

Rules:
- `uid` is the map identity — same uid updating position = moving dot, never a new dot.
- `hae` = 0 for ground assets (robots, vehicles, nodes). Use metres above ellipsoid for air/sea.
- `stale` drives the map fade-out: a dot whose last update is older than `stale` seconds goes dim.
- `meta` is open — each domain adds its own fields (operator, ICAO type, vessel flag, conflict type…).
- All timestamps are **epoch milliseconds UTC** — no ISO strings in hot-path geo messages.
- Publishers MUST include `lat`, `lon`, `uid`, `type`, and `t`. All other fields default to -1 / "".

### 9.3 Type Taxonomy

```
# Fleet (inner ring)
fleet.drone          fleet.robot          fleet.vehicle
fleet.node           fleet.edge           fleet.agent
fleet.aircraft       fleet.vessel

# World (outer ring)
air.commercial       air.military         air.unknown
sea.cargo            sea.tanker           sea.military       sea.unknown
ground.conflict      ground.military      ground.fire
ground.seismic       ground.chokepoint    ground.news

# TAK (bridge)
tak.atak             tak.itak             tak.wintak
```

### 9.4 Domain-Specific `meta` Shapes

```jsonc
// air.*
{ "icao": "A1B2C3", "operator": "usaf", "model": "B-52H",
  "origin_iata": "BOS", "dest_iata": "LAX", "squawk": "7700" }

// sea.*
{ "mmsi": "123456789", "flag": "US", "vessel_type": "cargo",
  "dest_port": "Shanghai", "imo": "1234567" }

// ground.conflict
{ "country": "UA", "event_type": "Explosion/Remote violence",
  "fatalities": 2, "source": "ACLED", "notes": "..." }

// ground.chokepoint
{ "chokepoint_id": "suez", "transit_status": "disrupted",
  "vessels_waiting": 47, "avg_wait_hours": 14.2 }

// fleet.drone / fleet.robot
{ "battery_pct": 82, "mode": "autonomous", "link_rssi": -67,
  "mission_id": "m-001", "node_id": "omos-barn" }
```

### 9.5 Adding a New Asset Class

> **Zero map-code changes.** Adding a new asset class to the TACMAP requires only:
> 1. Add a `type` entry in §9.3 vocabulary.
> 2. Publish `tel.fleet.{class}.{id}.position` with the standard envelope.
> 3. The map renders it. The agent sees it. ClickHouse records it.
>
> That is the doctrine. A Go2 robot dog, a greenhouse sensor, a manned plane —
> each appears on the map the moment it publishes its first position message.

### 9.6 Time Machine Integration

Every geo subject is recorded by the ClickHouse NATS-engine table (`onemind.events`).
This means:

- **Replay** — query `SELECT * FROM onemind.events WHERE subject LIKE 'tel.fleet.%'
  AND timestamp BETWEEN {t1} AND {t2}` → scrub the map backward in time.
- **Correlation** — `JOIN` inner-ring fleet positions with outer-ring world events by
  proximity and timestamp. "Was my drone near any military surge activity last Tuesday?"
- **Geofence history** — `evt.map.geofence.{id}.{enter|exit}` records every spatial crossing.

No additional logging infra needed. The Time Machine is already live.

### 9.7 Agent Interaction Contract

Legacy (and any agent with bus authority) interacts with the map via `cmd.map.>`:

```
cmd.map.filter.set       { layers: ["air.*","sea.*"], bbox: [lon0,lat0,lon1,lat1] }
cmd.map.highlight.{uid}  { color: "#FF4FD8", ring: true, label: "FOCUS" }
cmd.map.viewport.set     { lat: 26.5, lon: 56.5, zoom: 7 }  // center on Strait of Hormuz
cmd.map.geofence.{id}.create  { polygon: [[lon,lat],...], watch: true, alert: "evt.map.geofence" }
cmd.map.annotation.{id}.create  { lat, lon, label, color, icon, note }
```

The Flutter app subscribes to `cmd.map.>` and responds in real-time.
**The map is an input device, not just a display.**

### 9.8 Open Decision — TAK Bridge Subjects (live vs §9.1)

Recorded 2026-10-27. **No runtime subject was changed.** This section records how the live
system differs from §9.1 until the owner decides.

| | Live today (TAK 5.8, NatsBridgePlugin) | §9.1 doctrine |
|---|---|---|
| TAK → bus | `tak.cot.out`: one subject for all EUD CoT (`natsProducerSubject`) | `evt.tak.cot.{uid}`: per-UID geo-envelope |
| bus → TAK | `tak.cot.inject` (`natsConsumerSubject`) | not specified |
| Client state | none | `evt.tak.status.{uid}` |
| Stream | `COT` (`tak.cot.>`): receives the traffic; `TAKCoTFlowStopped` watches it | `GEO` (`evt.tak.cot.>`, `evt.tak.status.>`): defined but has 0 messages |

Evidence: `services/takserver/plugin/tak.server.plugins.NatsBridgePlugin.yaml`,
`platform/nats/streams/cot.yaml` and `geo.yaml`, and the live GEO stream message count.
`CLAUDE.md` already lists `tak.cot.inject` and `tak.cot.>` in the grammar block, so the two
documents disagree.

Options for the owner:

1. **Ratify `tak.cot.*` as the raw TAK wire.** Amend §9.1 so that `evt.tak.cot.{uid}` is a
   derived projection: a separate normalizer subscribes to `tak.cot.out` and publishes
   per-UID geo-envelopes into `GEO`. This needs no TAK change.
2. **Move the plugin to `evt.tak.cot.{uid}`.** This changes the plugin config and image,
   restarts `takserver-plugins` (TIER-1), and requires migrating `COT` consumers and the
   `TAKCoTFlowStopped` alert.

Until the owner decides, new consumers of TAK CoT read `tak.cot.out` (stream `COT`), and
injectors publish to `tak.cot.inject`. Nobody should publish `evt.tak.cot.*` by hand.

The edge TAK mesh gateway (`services/bridge/bridges/tak`, ADR 0003) uses the same raw-wire
convention for traffic that never touches TAK Server: CoT XML heard on the mesh goes to
`tak.cot.mesh.rx`, and CoT XML published to `tak.cot.mesh.tx` is broadcast to the mesh. Both
are in stream `COT`.

---

## 10. Sources

- synadia.com/blog/orbit-and-the-future-of-the-clients (Pietrek, Dec 2024)
- synadia.com/blog/synadia-agents-agentic-ai-era-will-be-heterogeneous (Collison, May 2026)
- synadia.com/blog/jetstream-design-patterns-for-scale (Connolly, Jan 2025)
- synadia.com/blog/decentralized-security-webinar (Humulock, Jul 2025)
- synadia.com/blog/nats-edge-event-architecture-8-mirror-merge-or-consume (Baloi, May 2026)
- github.com/synadia-ai/synadia-agent-sdk-docs (protocol v0.3 spec)
- Full blog index reviewed 2026-06-10; deep-dive backlog: Edge series 1–9, message tracing,
  per-message TTL, delayed scheduling, CRDT counters, MQTT-vs-NATS fleet comparison.
