# OneMind OS — Platform Architecture

## The Vision
OneMind is a polymorphic C2/ISR operating system. NATS is the nervous system. TAK is the display standard. Every sensor, AI agent, and operator interface is a first-class citizen on the same bus.

This platform replicates and extends the capabilities of:
- **Defender AI** → `omos-vision` (computer vision detections)
- **Leica BLK** → `omos-survey` (LIDAR / point cloud processing)
- **Octave** → `omos-site` (fixed facility / IoT monitoring)
- **ARCortex / TAKVR** → `omos-immersive` (3D/VR tactical environment)

Built on open protocols. Sovereign. No vendor lock-in. Runs air-gapped.

---

## How It All Comes Together

```
┌─────────────────────────────────────────────────────────────────┐
│                        SENSORS & INPUTS                         │
├──────────────┬──────────────┬──────────────┬────────────────────┤
│  omos-vision │ omos-survey  │  omos-site   │   External feeds   │
│  (cameras,   │  (LIDAR,     │  (IoT,       │   (ADS-B, AIS,     │
│  AI detect)  │  point cloud)│  facility)   │    ACLED, OSINT)   │
└──────┬───────┴──────┬───────┴──────┬───────┴────────┬───────────┘
       │              │              │                │
       ▼              ▼              ▼                ▼
┌──────────────────────────────────────────────────────────────────┐
│                    NATS BUS  (K8s HA cluster)                    │
│                                                                  │
│  det.visual.*    evt.survey.*    evt.facility.*    intel.*       │
│  tel.drone.*     tel.ship.*      tel.aircraft.*    ent.*.state   │
│  cmd.*           svc.*           agents.*                        │
└──────────┬───────────────────────────────┬───────────────────────┘
           │                               │
           ▼                               ▼
┌──────────────────────┐      ┌────────────────────────────────────┐
│   ORP (fusion)       │      │   Legacy AI (commander)            │
│                      │      │                                    │
│  Fuses all sensor    │      │  Reasons on the fused picture      │
│  inputs into unified │      │  Answers operator questions        │
│  entity states       │      │  Issues commands via cmd.*         │
│  ent.*.state →       │      │  Monitors missions                 │
│  NATS → TAK Server   │      │  Generates alerts + briefs         │
└──────────┬───────────┘      └────────────────┬───────────────────┘
           │                                   │
           └──────────────┬────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                      TAK SERVER                                 │
│                  (source of truth for all TAK content)          │
│                                                                  │
│  Receives CoT from ORP + Legacy AI                               │
│  Distributes to ALL connected clients                            │
│  Manages: missions, data packages, basemaps, channels           │
└────────────┬────────────────────────────────────────────────────┘
             │
             ├────────────────────────────────────────────────────┐
             ▼                                                    ▼
┌────────────────────────────┐              ┌─────────────────────────────┐
│   ATAK / iTAK / WinTAK    │              │   CloudTAK (web ops center) │
│                            │              │                             │
│  Field devices             │              │  Operators at computers     │
│  Superior TAK features     │              │  Planning + monitoring      │
│  Disconnected capable      │              │  Daily operations           │
│  Cesium 3D Tiles natively  │              │  AI agent interface         │
└────────────────────────────┘              └──────────────┬──────────────┘
                                                           │
                                            ┌──────────────▼──────────────┐
                                            │     CloudTAK Plugins        │
                                            │                             │
                                            │  cloudtak-plugin-nats       │
                                            │  → NATS bus console         │
                                            │                             │
                                            │  cloudtak-plugin-legacy-ai  │
                                            │  → AI agent chat            │
                                            │                             │
                                            │  cloudtak-plugin-ops        │
                                            │  → deck.gl 3D Tiles         │
                                            │  → terra-draw missions      │
                                            │  → Valhalla routing         │
                                            │  → LIDAR point clouds       │
                                            └─────────────────────────────┘

                        ┌────────────────────────┐
                        │   omos-immersive        │
                        │   (Unity / VR / AR)     │
                        │                         │
                        │  Connected to both:     │
                        │  TAK Server (CoT)       │
                        │  NATS bus (intelligence) │
                        │  3D Tiles (geo namespace)│
                        └────────────────────────┘
```

---

## Already Live — What's Running Today

### NATS Bus (K8s, 3-node HA)
The nervous system is live. All services below publish to it.

### Mac Leaf Node (hardware edge, cockpit)
Connected to K8s NATS as a forwarding leaf. Bridges hardware-bound inputs:
- **MAVLink bridge** → `tel.drone.*.position` / `tel.drone.*.attitude` (drone telemetry)
- **Home Assistant bridge** → smart home sensors, environmental data
- **Vision/cameras** → local camera feeds (feeds omos-vision)

### Redpanda Connect Pipelines (live data feeds)
Running in `data` namespace, publishing to NATS continuously:
- `feed-mta` → NYC transit (GTFS realtime) → `tel.transit.*`
- `feed-aircraft` → ADS-B aircraft positions → `tel.aircraft.*`
- `feed-ships` → AIS vessel positions → `tel.ship.*`
- `feed-hazards` → hazard alerts → `intel.hazard.*`

These are the precursors to `omos-feed-world` and `omos-feed-intel`.

### Legacy AI (Hermes) — live
Watching the bus, answering queries, injecting CoT. Connected to TAK Server via NATS leaf sidecar.

### CloudTAK — live
Web ops center at onemindos.dev. Latest commit: `157a0e83e` (Cesium imperative init fix).

### TAK Server 5.8 — live
Namespace `takserver`. NatsBridgePlugin wired: `tak.cot.inject` in, `tak.cot.>` out.

### ClickHouse + Grafana — live
Time-series event storage and dashboards at grafana.onemindos.dev.

---

## The Geo Namespace — Shared Data Infrastructure

All geospatial data services live in the `geo` Kubernetes namespace. Every layer feeds from here.

| Service | Role |
|---|---|
| **Martin** | Rust tile server — PostGIS vector tiles + PMTiles basemap |
| **TileServer GL** | MBTiles raster server — serves planet OSM to ATAK as WMS |
| **Tile38** | Real-time geofence engine — NATS native output |
| **TimescaleDB + PostGIS** | Time-series entity tracks + spatial queries |
| **STAC API** | Drone imagery + survey scan catalog (search by area/date) |
| **WebODM** | Drone photogrammetry → Cesium 3D Tiles mesh output |
| **Valhalla** | Routing engine — CloudTAK terra-route integration |
| **MinIO** | Object storage — raw LAS/LAZ scans, MBTiles, packages |

---

## Data Flow — End to End Example

**Scenario: Drone detects unknown vehicle, operator tasks response unit**

```
1. Drone camera stream → omos-vision
2. YOLO inference: "vehicle detected, 87% confidence"
3. Geo-registration: lat/lon computed from drone altitude + gimbal angle
4. NATS: det.visual.vehicle.track-42 published
5. ORP: fuses with existing AIS/ADS-B picture → "no friendly match"
6. ORP: ent.vehicle.track-42.state published
7. Legacy AI: "Unknown vehicle at grid 4QFJ1234, no friendly match, recommend intercept"
8. CoT generated → TAK Server
9. ATAK device shows hostile vehicle icon on map
10. CloudTAK web operator sees same icon + AI brief in Legacy AI panel
11. Operator draws response route in CloudTAK (terra-draw)
12. Route published to NATS cmd.route.publish
13. Mission pushed to TAK Server
14. ATAK unit receives route on their device
15. omos-immersive (VR): commander sees 3D visualization of vehicle + route
```

Total latency target: **< 3 seconds** from detection to all-client display.

---

## NATS Subject Grammar (locked)

```
tel.{node}.{class}.{metric}          telemetry          tel.drone.x47.position
det.{domain}.{type}.{id}             detections         det.visual.vehicle.track42
evt.{domain}.{entity}.{event}        semantic events    evt.facility.site01.alarm
cmd.{domain}.{entity}.{action}       commands           cmd.drone.x47.intercept
intel.{domain}.{source}.{type}       intelligence       intel.conflict.acled.event
ent.{entity_type}.{entity_id}.state  fused entity state ent.vehicle.track42.state
svc.{service}.{endpoint}             NATS RPC           svc.hermes.query
agents.{verb}.{agent}.{owner}.*      agent fabric       agents.prompt.hermes.zeus
```

---

## Services — Build Roadmap

| Wave | Service | Replicates | Description |
|---|---|---|---|
| Now | cloudtak-plugin-ops | — | deck.gl 3D Tiles + terrain config in CloudTAK |
| Now | cloudtak-plugin-nats fix | — | Proper Plugin API usage |
| 2 | omos-mission-engine | — | NATS mission orchestrator (workflow graphs) |
| 2 | omos-vision | Defender AI | Camera AI detection → NATS det.visual.* |
| 3 | omos-survey | Leica BLK | LIDAR/point cloud → 3D Tiles pipeline |
| 3 | omos-site | Octave | IoT/facility sensors → NATS evt.facility.* |
| 3 | omos-feed-world | — | AIS, ADS-B, OSINT automated feeds |
| 4 | omos-immersive | ARCortex/TAKVR | Unity VR/AR + NATS intelligence |
| 4 | Domain bridges | — | ROS (robotics), wellness, scheduler |

---

## Principles

1. **NATS is the nervous system.** Every service publishes and subscribes. No direct service-to-service HTTP in the hot path.
2. **TAK Server is the distribution layer.** CoT goes in, all TAK clients get it. We don't build separate distribution mechanisms.
3. **CloudTAK is the web ops center.** Minimal fork, maximum plugins. Upstream merges stay clean.
4. **Sovereign by default.** Everything runs air-gapped. No external dependencies in the hot path.
5. **Open protocols.** CoT, NATS, 3D Tiles, MQTT, LAS/LAZ — all open standards. No proprietary wire formats.
6. **One picture.** ORP ensures every client — ATAK, CloudTAK, VR, AI — sees the same fused reality.
