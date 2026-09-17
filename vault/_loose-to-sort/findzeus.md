---
title: "Find Zeus"
type: reference
domain: SO
status: active
created: 2026-06-07
updated: 2026-06-07
tags: [find-zeus, emergency, location, presence, safety, sovereignty, nats]
---

# Find Zeus

Real-time location, presence, and welfare protocol for Zeus Delacruz. NATS is the central nervous system — every sensor layer publishes here. Legacy aggregates, reasons, and escalates.

---

## NATS Architecture — Center of Everything

All detection layers publish to NATS. Legacy subscribes and maintains a live Zeus status model.

```
┌──────────────────────────────────────────────────────────────────┐
│                         NATS BUS                                 │
│                                                                  │
│  tel.person.zeus.location     ← last known location              │
│  tel.person.zeus.presence     ← zone-level presence state        │
│  tel.person.zeus.vitals       ← breathing rate, heart rate       │
│  tel.person.zeus.heartbeat    ← personal node keepalive          │
│  evt.person.zeus.detected     ← zone entered / exited            │
│  evt.person.zeus.alert        ← anomaly / distress signal        │
│  qry.person.zeus.location     ← request last known location      │
│  gov.escalate.findzeus        ← trigger full Find Zeus sweep     │
└──────────────────────────────────────────────────────────────────┘
         ▲           ▲           ▲           ▲           ▲
    Cameras       RuView       iPhone      HA Sensors  Zeus Node
    Frigate       WiFi Radar   OpenClaw    Presence    Personal
```

Subject spine: `<plane>.<class>.<id>.<path>` — fully compliant with OneMind NATS OS grammar.

---

## Detection Layers

### 1. Facial Recognition — Cameras (Frigate + CodeProject.AI)

- **Stack:** Frigate NVR + CodeProject.AI object/face recognition running on home server
- **Hardware:** Amcrest / Reolink cameras — House, Barn, Office entry points
- **Flow:** Camera → Frigate detects face → HA automation fires → NATS `evt.person.zeus.detected`
- **Resolution:** Zone-level (which camera = which room/area)
- **Integrations:** HA entity `binary_sensor.zeus_detected_<zone>`
- **Notes:** Works in the dark with IR cameras. No cloud — fully local inference.

### 2. Presence Sensors (HA Entities)

- **Stack:** Home Assistant binary sensors
- **Hardware:** mmWave sensors (e.g. LD2410), PIR motion sensors, door/window contacts per zone
- **Flow:** Sensor → HA state change → NATS bridge → `tel.person.zeus.presence`
- **Zones:** House (room-level), Barn, Office, Driveway
- **Notes:** Passive, no device required. Covers all nodes. First layer to fire.

### 3. WiFi Radar — RuView

- **Repo:** https://github.com/ruvnet/ruview
- **Stack:** RuView (WiFi DensePose) + ESP32-S3 CSI nodes + Cognitum Seed
- **Hardware:** ESP32-S3 boards ($9/node) — one per room for full coverage
- **Capabilities:**
  - Through-wall presence detection (up to 5m)
  - Breathing rate (6–30 BPM, real-time)
  - Heart rate (40–120 BPM, real-time)
  - Fall detection (<200ms)
  - 17-keypoint pose estimation
  - Room-level occupancy count
  - Anomaly detection: `possible-distress`, `fall-risk-elevated`, `inactivity-anomaly`
- **Integration:** HA-native via MQTT (`--mqtt` flag). 21 entities per node published automatically.
- **Flow:** ESP32 CSI → RuView edge model → HA MQTT → NATS bridge → `tel.person.zeus.vitals` + `tel.person.zeus.presence`
- **Key feature:** Detects vital signs through walls — confirms Zeus is alive even without visual confirmation.
- **HA Blueprints:** 3 starter blueprints ship with RuView (distress detection, inactivity alert, sleep monitoring)

### 4. Home Assistant + iPhone App

- **Stack:** Home Assistant + HA Companion App on Zeus's iPhone
- **iPhone App:** Location sharing via HA iOS integration — room-level precision indoors, GPS outdoors
- **Flow:** HA iOS → HA `device_tracker.zeus_iphone` → NATS bridge → `tel.person.zeus.location`
- **HA Entities:** `person.zeus`, `device_tracker.zeus_iphone`
- **Notes:** Most reliable outdoor layer. Works on cellular/WiFi. Automatic zone detection (Home, Barn, Office, Away).

### 5. iPhone Find My (macOS)

- **Stack:** FindMy.app on macOS + Hermes `findmy` skill
- **Flow:** Agent polls FindMy.app → extracts Zeus's iPhone location → publishes `tel.person.zeus.location`
- **Skill:** `apple/findmy`
- **Precision:** GPS-level when iPhone is on, last-known when offline
- **Notes:** Works even if Zeus's iPhone has no signal (uses U1/Bluetooth mesh from other Apple devices nearby)

### 6. iPhone OpenClaw Node

- **Stack:** OpenClaw agent running on Zeus's iPhone
- **Flow:** OpenClaw iPhone → heartbeat + location → NATS `tel.person.zeus.heartbeat` + `tel.person.zeus.location`
- **Notes:** Active node — Zeus's phone reports itself directly to the NATS bus. No passive inference needed.
- **Failover:** If heartbeat stops → trigger `evt.person.zeus.alert`

### 7. NATS Swift App (Custom)

- **Stack:** Native Swift iOS/macOS app subscribing directly to NATS
- **Subjects:** `tel.person.zeus.>` — real-time Zeus status dashboard on any Apple device
- **Features:**
  - Live location on map
  - Vital signs from RuView nodes
  - Heartbeat status
  - Zone presence feed
  - One-tap `gov.escalate.findzeus` trigger
- **Status:** Planned — build target when NATS bus is fully wired
- **Notes:** Swift NATS client libraries available (NIO-based). Lightweight, no middleware.

### 8. Zeus Personal Node

- **Stack:** Dedicated device Zeus carries — Pi Zero 2W or custom ESP32 + cellular
- **Subjects published:** `tel.person.zeus.heartbeat` (every 60s), `tel.person.zeus.location` (GPS)
- **Features:**
  - Publishes heartbeat — if silent >5 min → `evt.person.zeus.alert`
  - GPS telemetry when outdoors
  - BLE advertising for room-level triangulation indoors
  - Physical panic button → `gov.escalate.findzeus`
- **Hardware refs:** `SO/03 Infrastructure/hardware/pi-zero-2w.md`
- **Notes:** Last-resort layer. Works when phone is dead or left behind.

### 9. Zeus Glasses Interface

- **Stack:** AR glasses with NATS uplink (Meta Ray-Ban, custom ESP32 frame, or future OneMind HW)
- **Subjects published:** `tel.person.zeus.location`, `tel.person.zeus.heartbeat`
- **Features:**
  - Always-on presence signal — Zeus is wearing them = Zeus is at that location
  - Future: visual command interface to publish intents to NATS
  - Future: camera feed → Frigate facial recog loop for self-ID
- **Status:** Roadmap — wires into the same NATS subjects as all other layers

---

## Additional Detection Layers

### Vehicle Telemetry

- GPS trackers in cars, motorcycle, manned plane → `tel.vehicle.zeus-car1.location`
- If vehicle is at Barn → high probability Zeus is in Barn
- Cross-reference with Barn presence sensors to confirm

### BLE Beacon / AirTag

- AirTag in Zeus's wallet/bag → FindMy mesh triangulation
- BLE anchor nodes at House/Barn/Office → room-level BLE positioning
- Feeds into `tel.person.zeus.location` via agent poll

### Tailscale Node Presence

- If Zeus's MBP or iPhone is reachable on Tailscale mesh → Zeus has connectivity
- Agent checks `tailscale status` → if `omos-mbp` is online → Zeus is reachable
- Loss of Tailscale presence = possible offline/travel state

### Smart Power Inference

- Barn smart plugs active (CNC, welder, etc.) → Zeus is likely in Barn
- HA energy monitoring entities → cross-reference with presence sensors
- Publishes inferred zone to `tel.person.zeus.presence`

### Dead-Man Switch (Heartbeat Watchdog)

- Personal node + OpenClaw iPhone publish heartbeat every 60s
- Legacy monitors — if no heartbeat for configurable window (default: 30 min during active hours)
- Triggers: `evt.person.zeus.alert` → notify Cielo → initiate Find Zeus sweep
- Active hours: 6am – 11pm ET. Suspended 11pm – 6am.

### Cerberus Proxy

- Cerberus (dog) AirTag tracked via FindMy
- Zeus and Cerberus are frequently co-located → if Cerberus is at Barn, check for Zeus
- Soft signal only — not authoritative

### Calendar + Schedule Inference

- Legacy checks Todoist / calendar for Zeus's schedule
- If HA shows no home presence + calendar shows "flight" or "meeting" → infer location
- Reduces false-positive alerts for planned absences

---

## Zeus JetStream Watchdog

JetStream sits between the sensor layers and Legacy. It persists state, detects staleness, and fires alerts — no agent needs to be running for this to work.

### KV Store — `KV_ZEUS_STATUS`

One key per dimension. Every sensor layer writes here. Legacy reads a single source of truth.

```
zeus.location          ← last known location {lat, lon, zone, source, ts}
zeus.presence          ← last active zone + sensor that confirmed it
zeus.heartbeat         ← last heartbeat {source: iphone|node|glasses, ts}
zeus.vitals            ← last RuView reading {breathing_bpm, heart_bpm, ts}
zeus.vehicle           ← last vehicle GPS {id, lat, lon, ts}
zeus.status            ← aggregated state: home|barn|office|away|unknown|alert
zeus.active_hours      ← bool — is it within 6am–11pm ET window
```

TTLs are set per key — if a key expires without renewal, JetStream fires a `$KV.DEL` event → watchdog catches the deletion → triggers alert.

### Watch Service / Script

Runs as a LaunchAgent on `omos-mbp`. Subscribes to:

```
tel.person.zeus.>      ← all sensor telemetry → write to KV_ZEUS_STATUS
evt.person.zeus.>      ← anomaly events → evaluate + escalate
$KV.zeus_status.>      ← key expirations → heartbeat stale detection
```

**Logic:**

```
on tel.person.zeus.heartbeat  → update KV zeus.heartbeat, reset TTL (5 min active / suspended 11pm–6am)
on tel.person.zeus.location   → update KV zeus.location + infer zeus.status
on tel.person.zeus.presence   → update KV zeus.presence
on tel.person.zeus.vitals     → update KV zeus.vitals, check for distress signals
on $KV.DEL zeus.heartbeat     → heartbeat TTL expired → publish evt.person.zeus.alert
on evt.person.zeus.alert      → run Agent Sweep → notify Cielo
on gov.escalate.findzeus      → force-run full sweep immediately
```

**Distress detection from vitals:**
- Breathing rate drops to 0 or stops publishing → `possible-distress`
- Heart rate anomaly (RuView `fall-risk-elevated` or `possible-distress` entities) → `evt.person.zeus.alert`
- Inactivity anomaly with no heartbeat → escalate

### Deployment

```bash
# LaunchAgent: ~/Library/LaunchAgents/com.onemind.findzeus-watchdog.plist
# Script: ~/onemind/_system/findzeus-watchdog.py (or .sh)
# NATS connection: nats://localhost:4222 (local) or Tailscale IP (remote)
```

- Lightweight — event-driven, not polling
- Survives agent restarts — JetStream KV holds state
- Works when Legacy is offline — JetStream holds the last known state until Legacy reconnects
- Single source of truth: any client (Swift app, OpenClaw, Legacy) queries `KV_ZEUS_STATUS` for instant status

---

## Agent Sweep — Find Zeus Protocol

When `gov.escalate.findzeus` fires, Legacy runs this sweep in order:

```
1. Query tel.person.zeus.location        → last published location + timestamp
2. Query tel.person.zeus.heartbeat       → time since last heartbeat
3. Query tel.person.zeus.presence        → last active zone
4. Poll FindMy (findmy skill)            → iPhone GPS position
5. Check HA person.zeus                  → HA device tracker state
6. Check Tailscale (zeus devices online?)→ connectivity status
7. Check calendar                        → any scheduled absence?
8. Compile full status report
9. Notify Cielo with report              → Cielo decides next action
```

Escalation hierarchy:
1. Legacy flags + notifies Cielo
2. Cielo attempts direct contact (iMessage, call)
3. If unreachable: Liam / Noah / Soluna notified
4. Local authorities — **Cielo authorization required**

---

## Known Locations

| ID | Location | Address | Nodes Present |
|----|----------|---------|---------------|
| `home` | House | Greenwich, CT | Cameras, RuView, Presence sensors, HA |
| `barn` | Barn | Greenwich, CT | Cameras, Presence sensors, smart plugs |
| `office` | Office | TBD | TBD |
| `away` | Unknown / traveling | — | Vehicle GPS, iPhone, Tailscale |

---

## Offline / Travel State

Zeus is considered "intentionally offline" if:
- Calendar shows travel/flight
- HA zone = `not_home` for >1 hour
- Last heartbeat timestamp matches expected departure time

In this state: queue non-urgent items, suppress heartbeat alerts, maintain location tracking via iPhone/vehicle GPS.

---

## Notes

<!-- Specific device IDs, AirTag serials, vehicle tracker credentials → 1Password only -->
<!-- NATS subject contract refinements → SO/01 The Fabric/ONEMIND-NATS-OS-PLAN.md -->
<!-- RuView ESP32 node placement map → SO/03 Infrastructure/hardware/ -->
