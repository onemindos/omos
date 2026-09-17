# OneMind OS — NATS Fabric Architecture v2.0
**Author:** Legacy (OneMind NATS Meta-Agent)  
**Date:** 2026-06-30  
**Status:** ✅ LIVE — All streams, KV buckets, and object stores provisioned on fabric  
**Provisioning script:** `~/omos-core/nats/subscribers/provision-streams-v2.sh`

---

## Executive Summary

OneMind OS runs a sovereign, local-first NATS fabric across the homestead — connecting drones, robots, edge sensors, TAK operators, AI agents, SaaS services, and the operator (Zeus) on a single nervous system. This document defines the **v2.0 stream topology**: 14 streams, 17 KV buckets, and 5 object stores designed to handle the full operational scope of the fabric — from ESP32 soil sensors to AIS vessel tracking to ORP world intelligence to the operator's personal sovereign journal.

---

## Five-Layer Context

```
L5  INTERFACE    Flutter app (onemind-command) — stateless renderer
L4  TASK ENGINE  Issues + Tasks + Actors. JetStream + KV.
L3  COMPUTE      omd daemon + Agent layer (Legacy, Claude Code, workers)
L2  MEMORY       JetStream (operational) + ClickHouse (Time Machine, 10yr)
L1  SUBSTRATE    NATS — single transport for everything
```

**Rule:** The bus owns truth. The app is a window. Every byte flows through L1.

---

## Why v2.0 — What Changed from v1

| Problem in v1 (7 streams) | Solution in v2 (14 streams) |
|---|---|
| `tel.>` caught drones + soil sensors in one bucket — AIS/ADS-B from ORP would flood robot telemetry | TEL narrowed to fleet assets; ENV split for nature/bio (different retention, different consumers) |
| `evt.>` would absorb ORP world intel — AIS/ADS-B = millions of msgs/day at 90d retention | WORLD stream: 24h window + 2GB hard cap for raw ORP feeds |
| World intelligence had no distillation layer | INTEL stream: om-meta watch job writes curated alerts here — what Legacy monitors |
| Mission lifecycle scattered across CMD + EVT | MISSION stream: unified tactical log for every operation |
| No stream for Zeus's operator activity | ZEUS stream: 365d sovereign personal Time Machine |
| OM Insights (Phase 5) had no home for health events | SYS stream: system health + auto-checks foundation |
| COT was TAK-only; fleet positions had no durable home | GEO stream: TAK events + fleet position subjects |
| Zero KV buckets — all state in JetStream | 17 KV buckets: live state, fleet, agents, Zeus context, system |

---

## Stream Topology — Full Specification

### TIER 1 — FLEET + INTELLIGENCE

#### TEL — Fleet Telemetry
```
Subjects:  tel.fleet.>  tel.drone.>  tel.robot.>
           tel.edge.>   tel.iot.>    tel.industrial.>  tel.sat.>
Retention: 7 days
Storage:   file / discard=old
Max bytes: unlimited
```
**Purpose:** Continuous sensor/state data from all physical fleet assets. Drones report position 10x/sec. Robots report lidar, battery, mode. Edge nodes (ESP32, Pi) report sensor readings. Industrial machines via ORP Modbus/OPC-UA.  
**Who writes:** MAVLink bridge (drones), nats-ros-connector (ROS2 robots), IOnode firmware (ESP32), ORP industrial connectors.  
**Who reads:** TACMAP position layer, ClickHouse Time Machine, health watchdogs, robot GCS UI.  
**Subject grammar:**
```
tel.fleet.drone.{id}.position      → geo envelope (lat/lon/hae/course/speed)
tel.fleet.drone.{id}.battery       → battery percent
tel.fleet.robot.{id}.position      → robot position
tel.fleet.robot.{id}.lidar         → lidar scan data
tel.fleet.robot.{id}.battery       → robot battery
tel.edge.{node}.{sensor}           → ESP32/Pi sensor reading
tel.iot.{node}.{sensor}            → IoT device reading
tel.industrial.{machine}.{metric}  → industrial machine metric
tel.sat.{name}.pass                → satellite pass data
```

---

#### ENV — Environmental + Biometric
```
Subjects:  tel.nature.>  tel.bio.>  tel.weather.>
           tel.soil.>    tel.water.>
Retention: 30 days
Storage:   file / discard=old
```
**Purpose:** Long-horizon environmental and health data. Soil sensors report every 15 min — different access pattern from drones at 10Hz. Biometrics are personal/private. 30d retention for trend analysis.  
**Why not in TEL:** Soil sensor replay would flood drone telemetry consumers. Bio data needs account-level privacy separation (future FABRIC account isolation).  
**Who writes:** omos-rads `land_agent.py`, WireClaw ESP32 soil sensors, HRV wearable bridge, greenhouse-pi leaf.  
**Subject grammar:**
```
tel.soil.{location}.moisture       → soil moisture %
tel.soil.{location}.temperature    → soil temp °C
tel.weather.{station}.wind         → wind speed/direction
tel.weather.{station}.rain         → rainfall mm
tel.water.{location}.level         → water level
tel.nature.{sensor}.{metric}       → generic nature sensor
tel.bio.zeus.hr                    → Zeus heart rate (bpm)
tel.bio.zeus.hrv                   → Zeus heart rate variability
tel.bio.zeus.sleep                 → sleep quality score
tel.bio.zeus.stress                → stress index
```

---

#### WORLD — ORP World Intelligence (raw)
```
Subjects:  evt.world.>
Retention: 24 hours
Storage:   file / discard=old
Max bytes: 2 GB (hard cap — AIS/ADS-B is a firehose)
```
**Purpose:** Raw real-time world intelligence from ORP's 40+ protocol connectors. This is the operational window — the last 24 hours of global situational awareness. ClickHouse Time Machine keeps the permanent record.  
**⚠️ CRITICAL:** Do NOT increase retention here. AIS alone can produce millions of messages/day. The 2GB cap is a safety valve. ClickHouse absorbs everything for long-term storage.  
**Who writes:** ORP connectors (AIS, ADS-B, ACLED, NASA FIRMS, USGS).  
**Who reads:** om-meta watch job (filters → INTEL), TACMAP world intel overlay.  
**Subject grammar:**
```
evt.world.air.{icao}               → ADS-B aircraft position (geo envelope)
evt.world.air.mil.{icao}           → Military aircraft
evt.world.sea.{mmsi}               → AIS vessel position (geo envelope)
evt.world.conflict.{id}            → ACLED conflict event
evt.world.fire.{id}                → NASA FIRMS wildfire perimeter
evt.world.seismic.{id}             → USGS earthquake event
evt.world.military.{id}            → Military base/position
evt.world.chokepoint.{id}          → Strategic waterway status (Suez, Malacca, Hormuz)
evt.world.news.{country}.{id}      → Geo-tagged intelligence event
```

---

#### INTEL — Curated World Alerts
```
Subjects:  intel.>
Retention: 90 days
Storage:   file / discard=old
```
**Purpose:** Distilled world intelligence — events significant enough to review. WORLD is every ship. INTEL is the wildfire 200km from your homestead, the military surge in a watched region, the geofence crossing. This is what Legacy monitors and surfaces to Zeus.  
**Who writes:** om-meta watch job (subscribes WORLD → applies rules → publishes `intel.alert.*`).  
**Who reads:** Legacy (meta-agent), ZEUS stream (alerts), Flutter app Intel page.  
**Subject grammar:**
```
intel.alert.wildfire.{id}          → Wildfire within configured radius
intel.alert.seismic.{id}           → Earthquake above threshold magnitude
intel.alert.conflict.{id}          → Conflict event in watched region
intel.alert.military.surge.{id}    → Military activity surge detected
intel.watch.{zone}.{event}         → Watch zone trigger (geofence crossing)
intel.summary.daily                → Daily world intel summary (Legacy generated)
```

---

### TIER 2 — CONTROL

#### CMD — Commands
```
Subjects:  cmd.>
Retention: 30 days
Storage:   file / discard=old
```
**Purpose:** Every instruction issued to every actor — the complete 30-day command audit trail. Drone waypoints, robot tasks, node actuator commands, agent task dispatches — all here.  
**Subject grammar (NATS-DOCTRINE §2 + §9):**
```
cmd.fleet.drone.{id}.goto          → fly to waypoint
cmd.fleet.drone.{id}.estop         → emergency stop
cmd.fleet.drone.{id}.rtl           → return to launch
cmd.fleet.robot.{id}.patrol        → start patrol route
cmd.fleet.robot.{id}.estop         → emergency stop
cmd.fleet.vehicle.{id}.{verb}      → vehicle command
cmd.mission.{id}.create            → create mission
cmd.mission.{id}.waypoint          → add waypoint to mission
cmd.mission.{id}.abort             → abort mission
cmd.mission.{id}.complete          → mark mission complete
cmd.node.{id}.exec                 → execute command on edge node
cmd.agent.{id}.task                → dispatch task to agent
cmd.map.filter.set                 → set TACMAP layer filter
cmd.map.highlight.{uid}            → highlight actor on map
cmd.map.viewport.set               → pan/zoom TACMAP
cmd.map.geofence.{id}.create       → draw geofence
cmd.map.annotation.{id}.create     → drop map annotation
```

---

#### MISSION — Mission Lifecycle + Map Events
```
Subjects:  evt.mission.>  evt.map.>
Retention: 90 days
Storage:   file / discard=old
```
**Purpose:** Unified tactical log for every operation. When you run a mission — deploy drone, send Go2 on patrol, coordinate TAK — the mission narrative is a coherent queryable unit: status changes, waypoints reached, aborts, geofence events, map interactions. Query "show me everything about mission M-042" against this single stream.  
**Note:** `cmd.mission.>` and `cmd.map.>` live in CMD (full audit). MISSION captures only the **outcome events** (`evt.mission.>`) and **map state events** (`evt.map.>`).  
**Subject grammar:**
```
evt.mission.{id}.created           → mission created
evt.mission.{id}.status            → status change (queued/active/paused/aborted/complete)
evt.mission.{id}.waypoint_reached  → drone/robot hit waypoint
evt.mission.{id}.abort             → mission aborted (with reason)
evt.mission.{id}.complete          → mission completed
evt.map.geofence.{id}.enter        → actor entered geofence
evt.map.geofence.{id}.exit         → actor exited geofence
evt.map.annotation.{id}.created    → annotation dropped on map
```

---

### TIER 3 — EVENTS

#### EVT — Homestead Semantic Events
```
Subjects:  evt.so.>  evt.hp.>  evt.le.>  evt.ge.>
           evt.task.>  evt.cloud.>  evt.bio.>
           evt.robot.>  evt.drone.>  evt.sat.>
Retention: 90 days
Storage:   file / discard=old
```
**Purpose:** All semantic events from your own homestead — the domain-structured event log. NATS-DOCTRINE domains: `so` (sovereignty ops), `hp` (homestead/physical), `le` (living environment), `ge` (general).  
**Note:** `evt.world.>` → WORLD stream. `evt.mission.>/evt.map.>` → MISSION. `evt.sys.>` → SYS. `evt.tak.*` → GEO. These are deliberately NOT in EVT — this is what makes EVT manageable.  
**Subject grammar:**
```
evt.so.task.created                → om_task: new task
evt.so.task.{id}.status            → task status change
evt.so.task.{id}.complete          → task done
evt.le.home.{entity}.{event}       → Home Assistant: light on, door opened
evt.le.light.{id}.{state}          → specific light state change
evt.hp.greenhouse.soil.{alert}     → greenhouse soil alert
evt.hp.barn.{entity}.{event}       → barn sensor event
evt.cloud.ingress                  → cloud node inbound connection
evt.bio.alert                      → biometric threshold crossed
evt.robot.{id}.estop               → robot emergency stop
evt.drone.{id}.geofence_breach     → drone geofence violation
evt.sat.pass.start                 → satellite pass starting
```

---

#### SAAS — SaaS I/O
```
Subjects:  saas.>
Retention: 90 days
Storage:   file / discard=old
```
**Purpose:** Every SaaS interaction — inbound webhooks (Todoist, GHL, Skool, GitHub, etc.) normalized by Bento/Wombat, plus outbound API request/response pairs. The complete SaaS audit trail.  
**Subject grammar:**
```
saas.{provider}.{resource}.{event}           → inbound webhook (normalized)
saas.{provider}.{resource}.{action}.request  → outbound API request
saas.{provider}.{resource}.{action}.response → outbound API response
```
**Live connectors:** universal-inbound.yaml (:8091), todoist-outbound.yaml. Tailscale Funnel → public URL.

---

### TIER 4 — TACTICAL

#### COT — TAK Cursor-on-Target (preserved)
```
Subjects:  tak.cot.>
Retention: 30 days
Storage:   file / discard=old
Status:    LIVE — 378+ msgs (ATAK clients, NatsBridgePlugin)
```
**Purpose:** Active TAK CoT position reports from ATAK clients (Zulu1, field operators). Preserved from v1 — has live data. The NatsBridgePlugin writes here from the TAK server container. The cot2constellation bridge reads here to fan out to Constellation Overwatch.  
**Migration path:** GEO will eventually absorb this when it grows into the full TACMAP layer. For now both coexist — COT drains at 30d naturally.

---

#### GEO — TACMAP Geo Layer
```
Subjects:  evt.tak.cot.>  evt.tak.status.>
Retention: 30 days
Storage:   file / discard=old
```
**Purpose:** The normalized TACMAP data layer — TAK events and fleet position data for the COP map. As fleet assets (drones, robots) publish using the NATS-DOCTRINE §9 geo envelope, their position data lives in TEL (`tel.fleet.{class}.{id}.position`) and gets read by the TACMAP via Direct Get on `fleet-state` KV.  
**Subject grammar:**
```
evt.tak.cot.{uid}                  → normalized TAK CoT event (geo envelope shape)
evt.tak.status.{uid}               → TAK client connection state
```
**Note on fleet positions:** `tel.fleet.*.*.position` subjects are captured by TEL stream. The TACMAP map layer reads latest position via `fleet-state` KV (Direct Get — no consumer needed). GEO captures the event-level TAK layer.

---

### TIER 5 — HUMAN

#### VOICE — Voice Sessions
```
Subjects:  voice.>
Retention: 90 days (extended from 30d)
Storage:   file / discard=old
```
**Purpose:** Complete voice interaction record — transcripts, TTS outputs, intent routing, audio references. Extended from 30d to 90d: voice transcripts are decision history ("what did I tell Legacy six weeks ago that started this mission").  
**Subject grammar:**
```
voice.{session}.transcript.partial → streaming partial transcript
voice.{session}.transcript.final   → final transcript (STT output)
voice.{session}.audio.out          → TTS response audio ref (Object Store key)
voice.{session}.intent             → parsed intent
evt.voice.utterance                → captured for ClickHouse (10yr)
```

---

#### ZEUS — Operator Activity Log ⭐
```
Subjects:  zeus.>  op.zeus.>
Retention: 365 days
Storage:   file / discard=old
```
**Purpose:** The operator's sovereign personal Time Machine. Every other stream records what the *fabric* did. ZEUS records what *Zeus* did. This is the stream that answers:
- *"What was I working on March 15th at 3pm?"*
- *"What was my HRV the day I made that drone deployment decision?"*
- *"What did I approve or reject and why?"*
- *"What was the world context when I launched that mission?"*

**365d retention** — a full year of operator context. Completely private. Local-first. Never in a cloud.  
**Subject grammar:**
```
zeus.presence.{state}              → online | offline | field | home | barn
zeus.location.{context}            → where Zeus is physically (home/barn/field)
zeus.bio.hr                        → heart rate
zeus.bio.hrv                       → heart rate variability
zeus.bio.sleep.quality             → sleep score
zeus.bio.stress                    → stress index
zeus.focus.current                 → active work context
zeus.decision.{task_id}.approved   → task/mission approved by Zeus
zeus.decision.{task_id}.rejected   → task/mission rejected
zeus.decision.{task_id}.delegated  → delegated to another actor
zeus.session.{id}.start            → work session started
zeus.session.{id}.end              → work session ended
zeus.voice.{session_id}            → voice interaction reference
zeus.activity.{type}               → general activity log
zeus.context.{topic}               → context Legacy loaded for Zeus
op.zeus.>                          → operator namespace (future multi-operator expansion)
```

**Power queries (ClickHouse + ZEUS stream):**
```sql
-- What was happening when I made a decision?
SELECT z.*, i.* FROM zeus z
JOIN intel i ON ABS(zeus.ts - intel.ts) < 3600000
WHERE z.subject LIKE 'zeus.decision.%'
AND z.ts BETWEEN {T-1hr} AND {T+10min};

-- My HRV vs decision quality over last 90 days
SELECT date, avg(hrv), count(decisions) FROM zeus
WHERE subject IN ('zeus.bio.hrv', 'zeus.decision.%.approved')
GROUP BY date ORDER BY date;
```

---

### TIER 6 — AGENTS + OBSERVABILITY

#### AGENTS — Agent Fabric
```
Subjects:  agents.>
Retention: 7 days
Storage:   file / discard=old
Status:    LIVE — 1531+ msgs (Legacy heartbeats, prompts active)
```
**Purpose:** The agent communication layer — every prompt, heartbeat, status update, and tool call across all agents on the fabric. Legacy (this meta-agent) lives here. 7d is the right window — recent agent activity matters; ClickHouse keeps the permanent history.  
**Subject grammar (Synadia Agent Protocol v0.3):**
```
agents.prompt.{framework}.{owner}.{session}   → prompt to agent
agents.hb.{framework}.{owner}.{session}       → agent heartbeat
agents.status.{framework}.{owner}.{session}   → agent status
agents.attachments.{session}                  → file/data attachments
$SRV.PING.agents                              → discovery sweep
$SRV.INFO.agents                              → agent capability info
```

---

#### SYS — System Health + Observability
```
Subjects:  evt.sys.>  sys.>  tel.host.>
Retention: 7 days
Storage:   file / discard=old
```
**Purpose:** The foundation for Phase 5 OM Insights. Every auto-check result, service health event, and host metric goes here. When the `om-insights` service fires its 10 auto-check rules (slow consumers, stream quorum, offline replica, high RTT, disk %, leaf disconnect, consumer lag, auth failures, cert expiry, heartbeat miss), they emit to `evt.sys.check.failed` → lands in SYS → Flutter app shows red card within 30s.  
**Subject grammar:**
```
evt.sys.check.failed               → auto-check rule failed
evt.sys.check.recovered            → auto-check rule cleared
evt.sys.service.{name}.down        → service went offline
evt.sys.service.{name}.up          → service came online
evt.sys.leaf.{node}.disconnected   → leaf node lost connection
evt.sys.stream.{name}.lag          → consumer lag alert
sys.metric.{service}.{metric}      → service-level metric
tel.host.{node}.cpu                → host CPU percent
tel.host.{node}.memory             → host memory percent
tel.host.{node}.disk               → host disk percent
tel.host.{node}.network            → host network stats
```

---

## Complete Stream Summary Table

| # | Stream | Subjects | TTL | Volume | Status |
|---|--------|----------|-----|--------|--------|
| 1 | **TEL** | `tel.fleet.>` `tel.drone.>` `tel.robot.>` `tel.edge.>` `tel.iot.>` `tel.industrial.>` `tel.sat.>` | 7d | HIGH | ✅ LIVE |
| 2 | **ENV** | `tel.nature.>` `tel.bio.>` `tel.weather.>` `tel.soil.>` `tel.water.>` | 30d | LOW-MED | ✅ LIVE |
| 3 | **WORLD** | `evt.world.>` | 24h | VERY HIGH | ✅ LIVE |
| 4 | **INTEL** | `intel.>` | 90d | LOW | ✅ LIVE |
| 5 | **CMD** | `cmd.>` | 30d | MED | ✅ LIVE |
| 6 | **MISSION** | `evt.mission.>` `evt.map.>` | 90d | LOW | ✅ LIVE |
| 7 | **EVT** | `evt.so.>` `evt.hp.>` `evt.le.>` `evt.ge.>` `evt.task.>` `evt.cloud.>` `evt.bio.>` `evt.robot.>` `evt.drone.>` `evt.sat.>` | 90d | MED | ✅ LIVE |
| 8 | **SAAS** | `saas.>` | 90d | MED | ✅ LIVE (12 msgs) |
| 9 | **COT** | `tak.cot.>` | 30d | HIGH | ✅ LIVE (378+ msgs) |
| 10 | **GEO** | `evt.tak.cot.>` `evt.tak.status.>` | 30d | HIGH | ✅ LIVE |
| 11 | **VOICE** | `voice.>` | 90d | LOW | ✅ LIVE |
| 12 | **ZEUS** | `zeus.>` `op.zeus.>` | 365d | LOW | ✅ LIVE |
| 13 | **AGENTS** | `agents.>` | 7d | MED | ✅ LIVE (1531+ msgs) |
| 14 | **SYS** | `evt.sys.>` `sys.>` `tel.host.>` | 7d | LOW | ✅ LIVE |

---

## KV Store Architecture — 17 Buckets

KV stores hold the **current truth**. Streams hold **history**. These are different jobs. Never substitute one for the other.

### Operational
| Bucket | TTL | History | Purpose |
|--------|-----|---------|---------|
| `om-config` | none | 10 | System-wide OneMind OS configuration |
| `om-schema` | none | 5 | JSON Schema registry — subject envelope contracts |
| `node-registry` | none | 5 | All registered nodes + actors (omd registrations) |
| `task-index` | none | 10 | om_task: all tasks (permanent) |
| `projects` | none | 5 | om_task: project definitions |
| `task-notify` | 48h | 1 | om_task: reminder dedup log |

### Fleet Live State
| Bucket | TTL | History | Purpose |
|--------|-----|---------|---------|
| `fleet-state` | 30s | 3 | Latest position+status per fleet asset — **TTL expiry = asset offline** |
| `mission-state` | none | 5 | Active mission states and last known status |
| `geofences` | none | 10 | Defined geofence polygons + watch zones |

### Agent Layer
| Bucket | TTL | History | Purpose |
|--------|-----|---------|---------|
| `agent-registry` | none | 5 | Agent heartbeat registry (om_task aggregator) |
| `agent-config` | none | 10 | Agent harness configs (read/edit via SDK — no shadow copies) |
| `agent-caps` | 60s | 1 | Live capability advertisements — **TTL expiry = agent offline** |

### Zeus Context
| Bucket | TTL | History | Purpose |
|--------|-----|---------|---------|
| `zeus-context` | 1h | 5 | Zeus's current session context (auto-refresh) |
| `zeus-presence` | 120s | 1 | Zeus online/offline/location — **TTL expiry = offline** |
| `zeus-focus` | 24h | 3 | Zeus current focus + priority (daily reset) |

### System
| Bucket | TTL | History | Purpose |
|--------|-----|---------|---------|
| `sys-health` | 60s | 1 | Per-service health state — **TTL expiry = service crashed** |
| `world-state` | 300s | 1 | Latest-per-entity world intel — Direct Get for current AIS/ADS-B |

### The TTL-as-Failure-Signal Pattern

The TTL on `fleet-state`, `zeus-presence`, `agent-caps`, and `sys-health` is not a cache expiry — it's a **liveness heartbeat**.

- Drone must republish `fleet-state.drone.alpha` every 10s. Entry expires in 30s. Expiration = drone offline.
- Zeus device republishes `zeus-presence` every 60s. Expires in 120s. Expiration = Zeus offline.
- Agents republish `agent-caps` every 30s. Expires in 60s. Expiration = agent crashed.
- Services write `sys-health.{name}` every 30s. Expires in 60s. Expiration = service crashed.

No heartbeat checker daemon needed. NATS KV TTL expiry IS the alert. Watch jobs subscribe to `$KV.{bucket}.>` and listen for `OP: DEL` (expiration events) to trigger downstream alerts.

---

## Object Store Architecture — 5 Buckets

For binary assets too large for message payloads.

| Bucket | Purpose |
|--------|---------|
| `firmware` | ESP32/Pi/robot firmware blobs for OTA updates |
| `mission-plans` | Mission waypoints, KML files, MAVLink plan files |
| `atak-packages` | ATAK Data Packages (.zip) for client distribution |
| `media` | Drone/robot photos + video clips (mission artifacts) |
| `voice-recordings` | Voice session audio (VOICE stream refs by key) |

---

## Data Flow Diagrams

### World Intelligence Pipeline (ORP → NATS → Legacy)
```
AIS Feed ──► ORP connector ──► evt.world.sea.{mmsi}      ──► WORLD (24h/2GB)
ADS-B ─────► ORP connector ──► evt.world.air.{icao}      ──► WORLD (24h/2GB)
ACLED ─────► ORP connector ──► evt.world.conflict.{id}   ──► WORLD (24h/2GB)
NASA FIRMS ► ORP connector ──► evt.world.fire.{id}        ──► WORLD (24h/2GB)
USGS ──────► ORP connector ──► evt.world.seismic.{id}     ──► WORLD (24h/2GB)
                                                                    │
                                                         om-meta watch job
                                                         (subscribe + filter)
                                                                    │
                                              intel.alert.{type}.{id} ──► INTEL (90d)
                                                                    │
                                                Legacy reads INTEL ─► zeus.activity.alert
                                                                    │
                                                       Zeus notified (push/voice)
```

### Fleet Asset Pipeline (Drone → NATS → TACMAP)
```
Drone (MAVLink) ──► omos-bridge ──► tel.fleet.drone.alpha.position ──► TEL (7d)
                                 └─► tel.fleet.drone.alpha.battery  ──► TEL (7d)
                                                                            │
                                                            ClickHouse NATS engine
                                                            (permanent record)
                                                                            │
                                                            fleet-state KV write
                                                            Key: drone.alpha
                                                            TTL: 30s (liveness)
                                                                            │
                                                     Flutter TACMAP Direct Get
                                                     $JS.API.DIRECT.GET.KV_fleet-state
                                                     → instant map update, no consumer
```

### SaaS Integration Pipeline (Bento/Wombat)
```
Webhook POST /webhooks/todoist ──► Bento :8091 ──► saas.todoist.task.created ──► SAAS (90d)
                                                                                       │
                                                                              om-meta watch job
                                                                              ─► evt.so.task.created
                                                                              ─► Legacy notified
                                                                              ─► om_task entry created
```

### Voice Pipeline
```
Microphone ──► om_presence ──► voice.{sess}.transcript.final ──► VOICE (90d)
                                                                         │
                                                              Legacy receives transcript
                                                              ─► processes intent
                                                              ─► publishes cmd.* (if action)
                                                              ─► voice.{sess}.audio.out (TTS)
                                                              ─► zeus.voice.{sess} (ZEUS record)
                                                              ─► evt.voice.utterance (ClickHouse)
```

---

## Account Topology (4 Accounts)

```
SYS      → $SYS.>                    (NATS server observability — monitoring only)
CORE     → agents.> cmd.> evt.> etc  (AI agents + meta services — Legacy lives here)
FABRIC   → tel.> tak.cot.> saas.>    (IoT, sensors, bridges, robots)
INTERFACE → read-only imports         (Flutter app, React cockpit)
```

**Cross-account data flows:**
- FABRIC exports `tel.>` → CORE imports (ORP reads fleet telemetry)
- CORE exports `evt.>` + `svc.>` → INTERFACE imports (app renders events)
- INTERFACE can publish `cmd.fleet.>` (approved subset only)
- ORP straddles FABRIC (reads `tel.>`) + CORE (writes `evt.world.>`)

---

## Subject Grammar Additions for NATS-DOCTRINE.md

These subjects are now CANONICAL and should be added to `docs/NATS-DOCTRINE.md §2`:

```
# ZEUS operator namespace
zeus.presence.{state}               online | offline | field | home | barn
zeus.location.{context}             current physical location
zeus.bio.{metric}                   hr | hrv | sleep | stress
zeus.decision.{task_id}.{outcome}   approved | rejected | delegated
zeus.session.{id}.{verb}            start | end
zeus.focus.{context}                current active focus
zeus.voice.{session_id}             voice interaction reference
zeus.activity.{type}                general activity log
zeus.context.{topic}                Legacy context load for Zeus
op.{operator}.>                     future multi-operator expansion

# System observability
sys.metric.{service}.{metric}       service-level metric
sys.check.{rule}.{outcome}          auto-check result (pass|fail|warn)
tel.host.{node}.{metric}            cpu | memory | disk | network

# World intelligence
evt.world.air.{icao}                ADS-B aircraft (geo envelope)
evt.world.sea.{mmsi}                AIS vessel (geo envelope)
evt.world.conflict.{id}             ACLED conflict event
evt.world.fire.{id}                 NASA FIRMS wildfire
evt.world.seismic.{id}              USGS earthquake
evt.world.military.{id}             military base/position
evt.world.chokepoint.{id}           strategic waterway status
evt.world.news.{country}.{id}       geo-tagged news/intel

# Intel (distilled)
intel.alert.{type}.{id}             curated world alert
intel.watch.{zone}.{event}          watch zone trigger
intel.summary.daily                 daily intel summary

# Mission lifecycle
cmd.mission.{id}.{verb}             create | waypoint | abort | complete | pause
evt.mission.{id}.{event}            created | status | waypoint_reached | abort | complete

# TAK normalized events (GEO stream)
evt.tak.cot.{uid}                   normalized CoT event
evt.tak.status.{uid}                TAK client state
```

---

## Implementation Roadmap

### Immediate (this session — ✅ DONE)
- [x] 14 streams provisioned on live fabric
- [x] 17 KV buckets provisioned
- [x] Provisioning script at `~/omos-core/nats/subscribers/provision-streams-v2.sh`
- [x] This documentation

### Next: Publisher Wiring
Each new stream needs at least one publisher to become active:

| Stream | Publisher to wire | When |
|--------|------------------|------|
| ZEUS | `zeus-heartbeat.py` — publishes `zeus.presence.online` + `zeus.bio.*` | Now |
| SYS | `om-insights` service Phase 5 | Phase 5 |
| WORLD | ORP NATS connector (`orp connectors add`) | When ORP world feeds go live |
| INTEL | `om-meta` watch job (subscribe WORLD → publish `intel.*`) | After WORLD is live |
| GEO | TAK bridge update to publish `evt.tak.cot.>` | When TACMAP v2 ships |
| ENV | omos-rads `land_agent.py` update | When soil sensors connect |
| MISSION | `om-meta` mission lifecycle wiring | Phase 4 |

### Near term
- [ ] Add ZEUS + subject grammar to `docs/NATS-DOCTRINE.md §2`
- [ ] Write `zeus-heartbeat.py` — simple presence publisher (cron or daemon)
- [ ] Wire ORP NATS connector to publish `evt.world.>` (fills WORLD stream)
- [ ] Write `om-meta` watch job skeleton for INTEL filtering
- [ ] Update `provision-streams.sh` (original v1) to reference this script
- [ ] ClickHouse: add NATS engine tables for new streams (WORLD, INTEL, ZEUS, SYS)

---

## Pitfalls + Lessons Learned (2026-06-30)

1. **NATS CLI `--server` flag overridden by `NATS_URL` env var** — use `--context omos` or bare `nats` with the default context configured. `nats --server nats://localhost:4222` fails when `NATS_URL` is set in the environment.

2. **Subject overlap = hard error** — NATS will reject stream creation if any subject overlaps with an existing stream. Order matters: narrow EVT first (delete+recreate), then create WORLD/MISSION/SYS. The script handles this by deleting empty streams in correct dependency order.

3. **`>` wildcard is end-of-subject only** — `tel.fleet.drone.>.position` is INVALID. `>` must be the terminal token. Use `tel.fleet.drone.*.position` (single-level wildcard `*` in the middle) — but this will overlap with `tel.fleet.>`. For per-class position subjects, either put them all in TEL (where they belong) and use KV for current-state reads via Direct Get.

4. **`nats server ping --count` flag doesn't exist in v0.4.0** — use `nats server ping 1` (positional arg for expected server count).

5. **`set -u` + bash arrays** — when using `set -euo pipefail`, bash array expansions like `"${extra_args[@]}"` fail as "unbound variable" when the array is empty. Use a string variable instead: `extra_args=""` / `$extra_args` (unquoted, intentional word-splitting for flags).

6. **COT has live data — never delete it** — 378+ msgs from active ATAK clients. COT and GEO coexist until COT drains at its 30d TTL. Do not attempt to recreate COT.

7. **KV TTL is bucket-level max-age (stream max-age)** — `nats kv add bucket --ttl 30s` sets the JetStream stream's `max-age` to 30s. Individual key TTL is a different thing (per-message TTL, different flag). What we configured is the entire bucket expires entries at 30s. This is correct for the fleet-state/presence/caps liveness pattern.

---

## Files

| File | Purpose |
|------|---------|
| `~/omos-core/nats/subscribers/provision-streams-v2.sh` | Idempotent provisioning script |
| `~/omos-core/docs/NATS-DOCTRINE.md` | Binding rules (needs ZEUS + new subjects added) |
| `~/omos-core/docs/ARCHITECTURE.md` | Five-layer design (update stream count to 14) |
| `~/Desktop/ONEMIND-NATS-ARCHITECTURE-v2.md` | **This document** |

---

*Designed and provisioned by Legacy — OneMind OS NATS Meta-Agent*  
*2026-06-30 · omos-mbp · NATS Server 2.14.2 · JetStream file-backed*
