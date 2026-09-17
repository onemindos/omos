# OneMind — The Unified Data Flow Pattern
**Date:** 2026-07-06 | **Status:** CONFIRMED as the canonical architecture route

---

## The Law: Everything follows ONE path

Once ORP is on the NATS bus, adding ANY new device, source, sensor, human, or intel feed follows this exact pattern. No exceptions.

```
REGISTER → INGEST → FUSE → PROJECT → DISPLAY → CORRELATE → ACT
```

---

## The Seven Steps

### 1. REGISTER
Publish entity metadata to `KV_device-registry`.

```json
{
  "id": "drone.x47",
  "kind": "drone",
  "name": "X-47 Reaper",
  "status": "online",
  "source": "om-bridge-mavlink",
  "subjects": {
    "telemetry": "tel.drone.x47.>",
    "commands": "cmd.drone.x47.>",
    "state": "ent.drone.x47.state"
  },
  "visibility": { "tak": true, "atoc": true, "feed": true, "legacy_watch": true },
  "tak": { "uid": "OM.drone.x47", "cot_type": "a-f-A-M-F-Q", "layer": "assets", "label": "X-47" },
  "capabilities": ["fly", "camera", "intercept", "rtl"],
  "location": { "lat": 30.267, "lon": -97.743, "hae": 120 },
  "metadata": { "airframe": "multirotor", "max_speed_ms": 22 }
}
```

This is the entity's birth certificate. It declares what it is, where it talks, what it can do, and how it appears on every display.

### 2. INGEST
The device (or its bridge) publishes raw telemetry to its registered subject.

```
tel.drone.x47.position    → { lat, lon, alt, speed, heading }
tel.drone.x47.battery     → { voltage, current, remaining_pct }
tel.drone.x47.heartbeat   → { mode, armed, gps_fix }
tel.ship.mmsi123456.pos   → { lat, lon, sog, cog, nav_status }
det.fire.firms.event      → { lat, lon, confidence, frp, satellite }
intel.conflict.acled.event→ { lat, lon, event_type, fatalities, actors }
tel.iot.livingroom.temp   → { state: "72.1", unit: "F" }
```

Every source publishes in its natural shape. No transformation required at ingest — ORP handles normalization.

### 3. FUSE (ORP — the Ontology engine)
ORP subscribes to `tel.>`, `det.>`, `evt.world.>`, `intel.>` via its NATS adapter.

For every event it:
- Creates or updates the entity in its graph (DuckDB + in-memory)
- Resolves identity (is this the same ship we saw from AIS and from visual detection?)
- Computes relationships (drone near geofence? ship heading to port?)
- Runs anomaly detection (orp-ml isolation forests)
- Maintains confidence scores per entity
- Signs events into the audit chain (Ed25519)

ORP is the SINGLE POINT OF FUSION. Nothing else fuses. Everything else is a view.

### 4. PROJECT
ORP publishes fused entity state back to the bus:

```
ent.ship.mmsi123456.state   → KV_world-state + stream
ent.drone.x47.state         → KV_world-state + stream
ent.fire.firms_12345.state  → KV_world-state + stream
ent.human.zeus.state        → KV_world-state + stream
```

Every projector (TAK, atoc, WWV, Grafana) subscribes to `ent.>` or specific subtrees.

### 5. DISPLAY
Multiple projectors consume fused state and render for their medium:

| Projector | Subscribes | Outputs | Display |
|-----------|-----------|---------|---------|
| om-tak-projector | `ent.*.*.state` | CoT XML → `tak.cot.inject` | iTAK, ATAK, WebTAK, CloudTAK |
| om-control-plane | `ent.*.*.state` | `ui.*` projections | atoc (Flutter) |
| WWV bridge | `ent.*.*.state` | WebSocket → DataBus events | WorldWideView (3D globe) |
| ClickHouse sink | `ent.*.*.state` | INSERT INTO entity_history | Grafana dashboards, time-travel |

ONE entity, ONE fused state, MANY views. The COP is not a product — it's a projection.

### 6. CORRELATE
Legacy (AI Commander) queries ORP-QL via NATS request/reply (`svc.orp.query`):

```sql
-- What ships are near our geofence?
MATCH (s:Ship) WHERE NEAR(s, 30.267, -97.743, 50) RETURN s.name, s.speed, s.properties

-- What hostile tracks are heading toward the port?
MATCH (t:Aircraft)-[:HEADING_TO]->(p:Port) WHERE t.properties.affiliation = 'hostile' RETURN t, p

-- What's within 10km of our drone?
MATCH (e:Entity) WHERE NEAR(e, 30.5, -97.2, 10) RETURN e.entity_type, e.name, e.confidence

-- Correlate: is this AIS track the same as this radar detection?
MATCH (a:Ship {mmsi: "123456789"})-[:NEAR]->(d:Detection) RETURN a, d
```

Legacy doesn't need to hold state. It QUERIES the state. ORP is the memory. Legacy is the reasoning.

### 7. ACT
Legacy issues commands through the approval-gated command path:

```json
{
  "v": 1,
  "id": "uuid",
  "ts": "2026-07-06T15:00:00Z",
  "cmd": "cmd.drone.x47.intercept",
  "source": "agent.legacy",
  "data": {
    "target_entity_id": "ent.drone.unknown_042",
    "intercept_mode": "follow",
    "max_altitude_m": 200
  },
  "approval": {
    "required": true,
    "status": "pending",
    "reason": "Hostile drone detected near perimeter"
  }
}
```

Protected commands (`cmd.drone.launch`, `cmd.intercept.*`, `cmd.robot.move`) require human approval. Non-protected commands (`cmd.entity.*.tak.show`, `cmd.layer.*.show`) execute immediately.

---

## The Universal Entity Lifecycle

```
BORN (register) → SEEN (ingest) → KNOWN (fuse) → SHOWN (project/display) → UNDERSTOOD (correlate) → ACTED UPON (cmd)
     ↑                                                                                                      |
     └──────────────────────── feedback (evt.command.executed, new telemetry) ───────────────────────────────┘
```

Every entity — drone, ship, fire, person, sensor, news event, home device — lives this lifecycle. No special cases.

---

## Adding a New Source (Recipe)

To add ANY new data source to OneMind:

1. **Write a connector** (or use ORP's built-in adapter if the protocol is supported)
   - Connector subscribes/polls the external source
   - Publishes to the appropriate NATS subject (`tel.*`, `det.*`, `intel.*`, `evt.world.*`)
   - Uses the Event Envelope v1 schema

2. **Register entities** it produces in `KV_device-registry`
   - Define kind, subjects, capabilities, TAK rendering hints

3. **Done.** ORP fuses it. Projectors display it. Legacy can query and act on it.

No new infrastructure. No new map. No new display code. The system absorbs it.

---

## Adding a New Device/Asset (Recipe)

To add a new controllable asset (drone, robot, vehicle, human operator):

1. **Register** in `KV_device-registry` with `capabilities` listing what it can do
2. **Bridge** publishes its telemetry to `tel.{kind}.{id}.>`
3. **Bridge** subscribes to `cmd.{kind}.{id}.>` for commands
4. **ORP** fuses it into the entity graph automatically
5. **Legacy** can now task it via `cmd.{kind}.{id}.{action}`

Example — adding a new ground robot:
```json
{
  "id": "robot.spot-03",
  "kind": "robot",
  "subjects": {
    "telemetry": "tel.robot.spot-03.>",
    "commands": "cmd.robot.spot-03.>"
  },
  "capabilities": ["move", "camera", "patrol", "rtb"],
  "tak": { "cot_type": "a-f-G-U-C-R", "layer": "assets", "label": "SPOT-03" }
}
```

---

## Why This Works (the properties that make it Palantir-grade)

1. **One graph, many views** — TAK/CloudTAK/atoc/WWV/Grafana all show the same truth
2. **Source-agnostic** — ORP normalizes everything into Entity + Relationship
3. **Queryable** — ORP-QL gives Legacy (and humans) full graph + geospatial query
4. **Correlatable** — entity resolution + relationships + NEAR/WITHIN = cross-source fusion
5. **Actionable** — command envelope with approval gates = safe autonomous action
6. **Auditable** — Ed25519 signed event chain = tamper-proof operational history
7. **Offline-capable** — NATS leaf nodes + ORP local instance = full capability without cloud
8. **Scalable** — add sources by adding connectors, not by changing architecture

---

*This is the route. Everything else is implementation detail.*
