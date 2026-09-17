# OneMind Android APK Roadmap
**Owner:** Zeus Delacruz | **Fabric:** OneMind OS | **Device:** Samsung S24 Ultra (SM-S928U1)
**Date:** 2026-07-16 | **Status:** Planning

---

## Architecture Vision

```
S24 Ultra
├── ATAK (tactical layer)          ← map, SA, missions, overlays
├── OneMind Agent (background)     ← NATS client, sensor publisher  ← BUILD FIRST
├── Legacy TAC (UI)                ← PTT to Legacy, mission status
└── ATAK Plugin (bridge)           ← Legacy intelligence on the map

         ↓ NATS WebSocket (Tailscale / wss://[tailscale]:8080)

OneMind Fabric
├── Legacy (cloud sovereign — agents.prompt.hermes.zeus.legacy)
├── ClickHouse Time Machine (all telemetry, 10yr TTL)
├── TAK Server 5.7 (CoT/SA — 152.42.155.95)
└── HA bridge, IoT, Dell edge, MacBook edge
```

The phone becomes a **first-class fabric node** — not just a viewer, but a publisher, subscriber, and actor on the NATS bus.

---

## APK Builds — Priority Order

---

### 🔴 TIER 1 — Build First (Foundation)

---

#### APK 1: OneMind Agent (Background Service)
**What it is:** Persistent background NATS client. The nervous system of the S24 on the fabric.
**Why first:** Everything else depends on it. Once the phone is a fabric node, all other APKs are additive.

**Core capabilities:**
- Maintains live NATS WebSocket connection to `wss://[tailscale-ip]:8080`
- Survives screen-off, battery optimization, Samsung Bixby killing it — runs as a foreground service
- Publishes device telemetry continuously:
  - `fabric.telemetry.zeus-s24.gps` — lat/lon/alt/speed/bearing at 1Hz
  - `fabric.telemetry.zeus-s24.battery` — %, charging state
  - `fabric.telemetry.zeus-s24.network` — WiFi/LTE/5G state, SSID
  - `fabric.telemetry.zeus-s24.motion` — walking/driving/stationary (ActivityRecognition API)
  - `fabric.telemetry.zeus-s24.sensors` — barometer, accelerometer summary
- Subscribes to `agents.cmd.zeus-s24.>` — Legacy can send commands to the phone:
  - Push notifications
  - Trigger TTS alerts
  - Request a camera snapshot
  - Trigger screen wake
- Subscribes to `agents.prompt.hermes.zeus.legacy` response channel → displays as system notification
- Registers as a NATS micro service (`nats micro ls` shows `zeus-s24` on the bus)

**Stack:** Kotlin · `io.nats:jnats` (NATS Java client) · WorkManager · ForegroundService
**Size estimate:** ~500 lines Kotlin for core. Small APK.
**Sideload:** Yes — ADB `adb install -r onemind-agent.apk`

---

#### APK 2: Legacy TAC (Command Interface)
**What it is:** The dedicated Legacy UI on the phone. Not a chat app — a sovereign command surface.
**Tagline:** "Legacy in your pocket."

**Core capabilities:**
- **PTT (Push-to-Talk) button** → mic → on-device Whisper transcription → NATS publish to `agents.prompt.hermes.zeus.legacy` → response TTS spoken back via earpiece/speaker
- **Mission Status Panel** — live pull from ClickHouse, shows current active mission, phase, contacts
- **Fabric Map** — who's on the NATS bus right now (all `nats micro ls` services rendered as nodes)
- **SA Dashboard** — all TAK contacts visible (fed from NATS bridge), last seen, battery, location
- **Marker Drop** — drops a CoT event on TAK AND sends context prompt to Legacy simultaneously ("Legacy, I just pinned this location — what do you know about it?")
- **Alert Feed** — all Legacy alerts, mission events, fabric anomalies in a scrollable timeline
- **Quick Actions** — configurable buttons for common prompts ("SITREP", "Route to waypoint", "Who's online", "Mission status")
- **TAK Chat relay** — read/send TAK server chat messages from within Legacy TAC (no need to open ATAK)

**Stack:** Kotlin · Jetpack Compose · NATS WebSocket · Whisper.cpp (Android port, on-device) · Android TTS
**Design:** Dark tactical UI, large touch targets for field use with gloves, AMOLED black background

---

### 🟡 TIER 2 — After Foundation

---

#### APK 3: ATAK Plugin — Legacy Overlay
**What it is:** A native ATAK plugin (Java/Kotlin, ATAK Plugin SDK). Legacy intelligence rendered directly on the tactical map.
**Why powerful:** Zero new UI — lives inside ATAK. Operators stay in one app.

**Core capabilities:**
- **"Ask Legacy" button** in ATAK toolbar → sends current map bounds + visible contacts + self position to Legacy → response rendered as CoT overlay on the map
- **AI annotations on markers** — tap any contact marker → "Ask Legacy about this contact" → callout rendered on map
- **Threat zone overlay** — Legacy draws geofence polygons on the map based on mission context
- **Route recommendation** — Legacy analyzes terrain/contacts and draws a recommended route as a CoT line
- **TAK Chat → Legacy** — ATAK chat channel `Legacy` becomes a two-way terminal to the mother
- **Mission briefing push** — at mission start, Legacy pushes a data package (markers, routes, zones) directly into ATAK via MARTI
- **Real-time CoT injection** — Legacy can inject arbitrary CoT events (markers, shapes, lines) at any time from cloud

**Stack:** Java/Kotlin · ATAK Plugin SDK · NATS via OneMind Agent IPC (bound service)
**Note:** ATAK Plugin SDK requires signing with the ATAK developer cert

---

#### APK 4: Sensor Feed Publisher
**What it is:** Dedicated high-frequency sensor telemetry publisher.

**Publishes:**
- GPS track history → `fabric.telemetry.zeus-s24.gps` (all points → ClickHouse → 10yr time machine)
- Accelerometer (fall detection, motion state) → `fabric.telemetry.zeus-s24.accel`
- Barometric pressure (altitude, weather prediction) → `fabric.telemetry.zeus-s24.baro`
- Ambient light/proximity
- Camera-based activity (optional: object detection via TFLite on-device)

**Value:** ClickHouse ingests everything. In 6 months you have a complete physical history — where you were, how fast, what you were doing. Mission replay works on your own timeline, not just TAK tracks.

---

#### APK 5: OneMind Comms
**What it is:** Sovereign encrypted messaging between all OneMind fabric nodes. Replace Signal for internal comms.

**Core capabilities:**
- NATS JetStream-backed persistent message store (not cloud, YOUR server)
- End-to-end: messages published to `fabric.comms.zeus.>`, only subscribers with NATS creds can read
- Push notifications via OneMind Agent
- Voice messages → Whisper transcription stored in ClickHouse
- File/image sharing via NATS object store
- Group threads per mission
- No phone number required — fabric identity only
- Works on any fabric node: S24, Dell edge, MacBook, future devices

**Stack:** Kotlin · JetStream · NATS Object Store · Whisper (transcription)

---

#### APK 6: HomeBase
**What it is:** OneMind home automation sovereign interface. Replaces the HA app.

**Core capabilities:**
- Controls all Home Assistant entities via NATS HA bridge (not direct HA API — through the bus)
- Hue lights, switches, sensors, locks — all NATS-routed
- Voice control: "Legacy, turn off the kitchen lights" → NATS → HA bridge → Hue
- Automation triggers publishable from phone: "I'm leaving" → NATS → HA → away mode
- Security feed: camera snapshots on demand via NATS
- Fabric-aware: if you're on Tailscale (away), routes through cloud; if on LAN, routes direct

---

### 🟢 TIER 3 — Future / Exploratory

---

#### APK 7: OneMind Dash (Family COP)
Common Operating Picture for the family. Not tactical — operational life.
- All family members as nodes on the map (with consent/opt-in)
- School pickup ETA, "I'm home" broadcasts, geofence alerts ("Noah arrived at school")
- Shared mission board (family tasks/errands as NATS-tracked items)
- Emergency PTT for family

#### APK 8: Intelligence Collector
- Passive intelligence: scan visible WiFi networks, BT devices, cell towers → publish to fabric
- Photo metadata extractor: take a photo, Legacy analyzes it and extracts intel (location, objects, faces if authorized)
- Voice memo → auto-transcribed → searchable via Legacy → stored in ClickHouse

#### APK 9: OneMind TAK Companion (iOS-equivalent)
For when an iPhone needs to join the fabric (family members). Lightweight NATS client + SA viewer.

---

## ATAK — Unused Capabilities to Activate

| Feature | Status | Action |
|---------|--------|--------|
| GeoFencing | ❌ Not set up | Draw zones → trigger Legacy alerts on NATS when contact enters/exits |
| Video feeds (RTSP) | ❌ Not set up | S24 Ultra camera → RTSP → ATAK video overlay |
| TAK Missions | ⚠️ Created, not operationalized | Legacy pushes mission data packages via MARTI |
| CoT custom events | ❌ Not used | Marker drops → trigger Legacy actions |
| Persistent Tracks | ❌ Not set up | All GPS breadcrumbs → ClickHouse |
| TAK Chat | ⚠️ Available | Wire to Legacy — chat channel becomes AI terminal |
| Data Packages | ❌ Not used | Push files/overlays to all contacts via MARTI |
| ADS-B Plugin | ⚠️ Check if installed | Real-time flights on map (see section below) |
| 9-line forms | ❌ Not used | Medical/CAS reporting forms |

---

## NYC Real-World Data Layers for ATAK/TAC

### Transit
| Source | Data | How to add |
|--------|------|------------|
| MTA GTFS-RT | Live subway car positions (all lines) | Feed into NATS → CoT events on map |
| MTA Bus Time | Live bus positions | Same pipeline |
| MTA Service Alerts | Delays, outages | NATS alert stream → Legacy TAC feed |
| LIRR / Metro-North | Train positions | GTFS-RT feed |

### Surveillance / Cameras
| Source | Data | How to add |
|--------|------|------------|
| NYC DOT Traffic Cameras | ~900 public cameras, MJPEG streams | WMS tile layer in ATAK or video overlay |
| NYC Open Data camera index | Camera locations as markers | CoT point overlay |
| Flock Safety | LPR camera network | Not public API — requires LE access or partnership |
| NYC LinkNYC | Kiosk locations (public WiFi nodes) | Static overlay |

### Traffic & Infrastructure
| Source | Data | How to add |
|--------|------|------------|
| NYC 511 / HERE Traffic | Real-time traffic speed/incidents | WMS overlay tile layer |
| NYC DOT incident feed | Road closures, construction | NATS alert stream |
| Con Edison outage map | Power outages | Scrape → NATS → map overlay |
| NOAA NWS | Weather radar, alerts | Already available as WMS — add as map layer |

### Aviation (see section below)
| Source | Data | How to add |
|--------|------|------------|
| OpenSky Network | ADS-B flights, free API | ATAK plugin or NATS poller → CoT |
| FAA SWIM | Full FAA flight data (requires registration) | Same |
| ADS-B Exchange | Unfiltered ADS-B, public | Best source, free |

### Maritime
| Source | Data | How to add |
|--------|------|------------|
| AIS (MarineTraffic / VesselFinder) | Ship positions in NYC harbor | NATS poller → CoT vessel markers |
| USCG NAIS | Official AIS feed | Free with registration |

---

## Flights in ATAK — What You Have and How to Get Them

**Plugins to check (if installed):**
- `ATAK Aviation Plugin` — adds ADS-B layer natively
- `Cursor on Target Bridge` — can receive ADS-B CoT from external sources
- `LayerManager` — can add WMS/WMTS layers including flight tracking tiles

**Best path if plugin isn't working:**
1. **ADS-B Exchange API** → free, unfiltered, global, 1–5s refresh
2. Build a small NATS poller (Python, runs on Legacy pod) that:
   - Polls `api.adsb.exchange.com` every 5s for flights in NYC area bounding box
   - Converts each aircraft to a CoT `a-f-A` (air friendly) event with callsign, altitude, speed, heading
   - Publishes to TAK server via UDP CoT multicast or direct TCP stream
   - All aircraft appear on ATAK map as plane icons in real time
3. ~50 lines of Python. No plugin needed.

---

## Build Order Summary

```
Phase 1 (Now):
  [ ] OneMind Agent APK — NATS background client on S24

Phase 2 (After Agent):
  [ ] Legacy TAC APK — PTT + mission status + fabric map
  [ ] ADS-B → NATS → TAK bridge (Python, runs on Legacy pod)
  [ ] MTA GTFS-RT → NATS → TAK bridge (Python, runs on Legacy pod)
  [ ] NYC DOT camera layer (WMS tile source XML for ATAK)

Phase 3 (After TAC):
  [ ] ATAK Plugin — Legacy Overlay (requires ATAK SDK cert)
  [ ] Sensor Feed Publisher APK
  [ ] OneMind Comms APK

Phase 4 (Homestead expansion):
  [ ] HomeBase APK
  [ ] Family COP (OneMind Dash)
  [ ] Intelligence Collector
```

---

## Notes
- All APKs sideloaded via ADB — no Play Store, no visibility, sovereign
- NATS transport: `wss://100.74.182.117:8080` (Tailscale) as primary, LAN fallback
- All telemetry → ClickHouse `onemind` database, 10yr TTL
- Legacy (mother) is the orchestration layer — all APKs report to her
- Device identity on fabric: `zeus-s24` / UID `R5CX11QVGGX`
