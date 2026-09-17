---
id: OSSIM-001
title: OneMind OSSIM — Simulation Project
status: draft
owner: zeus
created: 2026-06-14
tags: [simulation, nats, ossim, drone, robot, tak, clickhouse, no-gpu]
---

# OneMind OSSIM
## Open Sovereign Simulation & Intelligence Mesh

> Simulate the full OneMind physical world — drones, robots, sensors, buildings, global events —
> without real hardware, without a GPU, on existing compute.
> The bus is real. The data flows are real. Only the physical devices are synthetic.

---

## What This Is

OSSIM is a software-only simulation layer that runs on top of the real NATS bus.
Every bridge, every subscriber, every agent, every UI feature works identically.
The only difference: instead of a real drone publishing telemetry, a Python script does it.
Instead of a real sensor reading temperature, a VM publishes a random value on the correct subject.

**The bus does not know the difference. Neither does the UI.**

This means:
- You can demo the full system before owning a single drone or robot
- You can stress-test the bus with 100 simulated assets
- You can build and validate every feature of the platform on a laptop
- When real hardware arrives, you swap the simulator script for the real bridge — nothing else changes

---

## Compute Architecture (No GPU Required)

### Hub (omos — your MBP)
- NATS server running now (healthy, 2.14.1)
- ClickHouse receiving events
- WorldMonitor feeding real-world data
- Legacy AI orchestrating
- No GPU needed — all CPU workloads

### Virtual Nodes (VMs or Docker containers)
One VM or container per simulated building/zone:

| VM | Simulates | Subjects Published |
|----|-----------|-------------------|
| sim-barn | barn sensors, barn drone, barn robot | tel.drone.sim-barn-drone.* tel.robot.sim-barn-bot.* tel.sensor.barn.* |
| sim-office | office nodes, IP cameras, access control | tel.sensor.office.* evt.security.office.* |
| sim-greenhouse | temperature, humidity, pump, irrigation | tel.sensor.greenhouse.* cmd.device.greenhouse-pump.* |
| sim-house | home automation, presence, plugs | evt.home.* tel.device.* |
| sim-field | outdoor drone patrol area | tel.drone.skyhawk-1.* |

Each VM is lightweight — 512MB RAM, 1 CPU core. No GPU. No 3D rendering.
Run them as Docker containers on your Mac or on omos-cloud.

### Simulation Engine Options (CPU only)

**Option A — Pure Python simulators (simplest, start here)**
Each simulated asset is a Python script using nats-py.
Runs on any machine. Zero dependencies beyond Python.
```
sim_drone.py      — publishes GPS, battery, heartbeat on a timer
sim_robot.py      — publishes pose, battery, patrol events
sim_sensor.py     — publishes temperature, humidity, motion
sim_vehicle.py    — publishes GPS track along a route
sim_plane.py      — publishes ADS-B style position updates
```

**Option B — ArduPilot SITL (Software In The Loop)**
ArduPilot's SITL runs a full drone autopilot simulation entirely in software.
No hardware. No GPU. CPU only.
Connects to MAVLink bridge → bridge connects to NATS → drone appears on TAK map.
This is an actual simulated flight with real PX4/ArduPilot physics.
Install: `pip install MAVProxy && sim_vehicle.py -v ArduCopter`

**Option C — Gazebo (optional, GPU helps but not required)**
3D world simulation. Heavier but gives visual feedback.
Can run CPU-only in headless mode (no display).
Worth adding later for robot navigation testing.
Skip for initial OSSIM — not needed for the bus demonstration.

**Recommended start: Option A.** Get everything flowing through the bus first.
Add SITL for realistic drone flights as a second step.

---

## The Five Real-Time Data Layers (all CPU, all working now)

### Layer 1 — World Events (LIVE NOW via WorldMonitor)
Real data. No simulation needed. Already flowing.

```
evt.world.air.{icao}          ADS-B aircraft — real planes, real positions, live
evt.world.sea.{mmsi}          AIS vessels — real ships, real routes, live
evt.world.seismology.{id}     Earthquakes — real USGS feed
evt.world.fire.{id}           NASA FIRMS wildfires
evt.world.conflict.{id}       ACLED conflict events
```

These are streaming into ClickHouse right now.
A plane flying over Greenwich CT is already a dot on your future map.
No simulation needed for this layer — it is real.

### Layer 2 — Fleet Assets (simulated until hardware arrives)
Python scripts publishing on correct subjects:

```
tel.drone.skyhawk-1.gps       { lat, lon, alt_m, course, speed }
tel.drone.skyhawk-1.battery   { voltage_v, remaining_pct }
tel.robot.om1.pose            { x, y, heading }
tel.robot.om1.battery         { remaining_pct }
tel.vehicle.zeus-truck.gps    { lat, lon, speed, course }
```

Any subscriber, any UI, any TAK client sees these as real assets.

### Layer 3 — Sensor Mesh (simulated, trivial to make real)
```
tel.sensor.barn-temp-1.reading    { temp_f: 74.2, humidity: 58 }
tel.sensor.office-motion.reading  { motion: true, zone: "entry" }
tel.sensor.greenhouse.reading     { temp_f: 82.1, soil_moisture: 0.4 }
```

When real SwitchBot / ESP32 sensors arrive, their bridges publish identical subjects.
Simulator scripts are removed. UI never knows.

### Layer 4 — TAK / Field Devices (simulated CoT)
CoT bridge publishes simulated ATAK positions to TAK server.
Simulated field operators appear as moving dots on the ATAK map.
```
evt.tak.cot.sim-operator-1    { lat, lon, callsign: "ZEUS-1", type: "a-f-G" }
evt.tak.cot.sim-drone-cas-1   { lat, lon, callsign: "SKYHAWK-1", type: "a-f-A" }
```

### Layer 5 — AI Agents and Commands (real, running now)
Legacy is real. NATS is real. Commands flow through the real bus.
Simulation layer just means the assets responding to commands are software, not hardware.

---

## OSSIM Demo Flow — The One That Goes Viral

**Setup: laptop, no hardware, 5 Python scripts running**

1. Legacy receives voice command: "Launch skyhawk-1 to patrol the south field"

2. Legacy publishes to NATS:
   `cmd.drone.skyhawk-1.takeoff { altitude_m: 30 }`

3. sim_drone.py receives the command, starts publishing moving GPS coordinates:
   `tel.drone.skyhawk-1.gps { lat: 41.024, lon: -73.628, alt_m: 30, course: 180 }`

4. CoT bridge translates to TAK:
   Dot labeled SKYHAWK-1 appears on ATAK map, moving south

5. Subscriber rule fires on takeoff event:
   Todoist task created: "Skyhawk-1 patrol initiated 14:32"

6. ClickHouse records entire flight telemetry permanently

7. Flutter app /fleet/map shows the drone moving on MapLibre

8. ATAK on your phone shows the same dot

9. WorldMonitor ADS-B layer shows real aircraft nearby for context

10. Legacy narrates via voice: "Skyhawk-1 is airborne, altitude 30m, heading south"

**Everything real except the drone. Bus, map, AI, TAK, tasks, database — all live.**

---

## Watch Jobs — The Intelligence Layer

This is where OSSIM becomes the event intelligence engine.

Watch jobs are persistent NATS subscribers with an intelligence layer.
Each watch job monitors a subject pattern and acts when conditions are met.

### Example 1 — Earthquake near Zeus's location
```python
# Watch: subscribe to evt.world.seismology.*
# Condition: magnitude > 3.0 AND distance from lat:41.0 lon:-73.6 < 200km
# Action: publish cmd.notify { message: "Earthquake M4.2 detected 45km away" }
# Latency: seconds from USGS detection to notification
```

### Example 2 — Wife arriving home
```python
# Watch: subscribe to tel.device.cielo-iphone.gps
# Condition: distance from home < 0.5km AND was > 0.5km 60s ago
# Action: publish cmd.ha.lights.on, cmd.notify { message: "Cielo is almost home" }
```

### Example 3 — Drone low battery in field
```python
# Watch: subscribe to tel.drone.*.battery
# Condition: remaining_pct < 20
# Action: publish cmd.drone.{id}.return_home, notify Zeus
```

### Example 4 — Unauthorized motion in barn at night
```python
# Watch: subscribe to tel.sensor.barn-motion.reading
# Condition: motion == true AND time between 23:00-06:00
# Action: HIGH ALERT — cmd.drone.skyhawk-1.takeoff, publish security event, notify
```

### Example 5 — Military aircraft near property
```python
# Watch: subscribe to evt.world.air.mil.*
# Condition: distance from home < 50km
# Action: log to ClickHouse, notify if persistence > 10 min
```

**The watch system is om-insights (Python, already planned).
It subscribes to everything, evaluates rules, acts instantly.
No GPU. No AI inference required for the rules engine.
AI only intervenes when a rule escalates to "needs reasoning."**

---

## NATS as Universal Event Stream — What Can Actually Be Streamed

Yes. Anything that is digital or has a digital representation can flow through NATS.

### Physical World (real-time)
| Source | Data | NATS Subject | Latency |
|--------|------|--------------|---------|
| ADS-B receivers | Every aircraft position update | evt.world.air.{icao} | 1-5 seconds |
| AIS receivers | Every vessel position | evt.world.sea.{mmsi} | 30 seconds |
| USGS API | Earthquake within seconds of detection | evt.world.seismology.{id} | 30-120 seconds |
| NASA FIRMS | Wildfire detection | evt.world.fire.{id} | 15 minutes |
| Weather APIs | Storm, wind, precipitation | evt.world.weather.{id} | Minutes |
| GPS trackers | Any device with GPS | tel.device.{id}.gps | Real-time |
| Drone | Position, camera, sensors | tel.drone.{id}.* | 10Hz |
| Robot | Pose, camera, battery | tel.robot.{id}.* | 1Hz |
| Vehicle OBD | Speed, RPM, fuel | tel.vehicle.{id}.* | 1Hz |

### Digital World (event-driven)
| Source | Data | NATS Subject |
|--------|------|--------------|
| GitHub | Every push, PR, issue | evt.service.github.* |
| Stripe | Every payment, subscription | evt.service.stripe.* |
| Todoist | Every task create/complete | evt.service.todoist.* |
| Home Assistant | Every device state change | evt.home.* |
| Email | Every incoming email | evt.service.fastmail.received |
| Webhook (any) | Any platform that sends webhooks | evt.service.{platform}.* |
| Social feeds | Posts, mentions (via API polling) | evt.social.{platform}.* |
| Stock/crypto | Price ticks | evt.market.{symbol} |
| News feeds | RSS/API articles | evt.news.{source}.{category} |

### The key insight
Every one of these becomes a NATS subject.
Once it is on the bus, ClickHouse captures it permanently.
Once it is in ClickHouse, you can query it with SQL.
Once it is queryable, you can find patterns:
- "Every time there is a military aircraft surge within 100km and a seismic event within 48 hours"
- "Every time Zeus's wife leaves work between 5-6pm it takes 23 minutes to arrive home"
- "Barn temperature spikes every Tuesday at 2pm — correlates with power equipment use"

That is the correlation and pattern engine. No special ML needed for most patterns.
SQL + time-series queries on ClickHouse handle 90% of it.
LLM reasoning handles the rest.

---

## Storage — How Many Events, What Size

### NATS JetStream limits
No hard limit on number of messages.
You set retention policy per stream:
- By age: keep last 30 days
- By size: keep up to 50GB
- By count: keep last 1M messages
- Work queue: delete after consumption

### ClickHouse compression — real numbers
ClickHouse uses LZ4 compression by default, ZSTD optional.
Real-world compression ratios for event data:

| Data Type | Raw Size | Compressed (LZ4) | Ratio |
|-----------|----------|------------------|-------|
| GPS telemetry | 1KB/msg | ~80 bytes | 12:1 |
| Sensor readings | 500 bytes/msg | ~40 bytes | 12:1 |
| ADS-B aircraft | 800 bytes/msg | ~65 bytes | 12:1 |
| JSON events | 2KB/msg | ~200 bytes | 10:1 |

### What that means at scale

1 billion GPS events:
- Raw: ~1TB
- ClickHouse compressed: ~83GB
- On a $20/month VPS with 100GB SSD: fits

10 billion mixed events:
- Raw: ~10TB
- ClickHouse compressed: ~800GB
- Tiered storage (hot SSD + cold object storage): fully manageable

ClickHouse handles 10+ billion rows routinely.
Cloudflare uses it for 600 billion rows.
At OneMind scale you will never hit a wall.

Retention strategy:
- Hot tier (SSD): last 90 days, full resolution
- Warm tier (cheaper disk): 1 year, 1-minute aggregates
- Cold tier (S3/R2 object storage): 10 years, hourly aggregates
- Cost at 10B events/year: approximately $50-100/month

---

## Voice + Instant Steering (Combat-Speed Response)

### The Architecture
```
Zeus voice (mic)
    → LiveKit STT (Whisper/Deepgram, self-hosted)
    → voice.{session}.transcript.final  [NATS — sub-second]
    → Intent Agent (classify command in <500ms)
    → cmd.drone.skyhawk-1.follow_target [NATS — <1ms delivery]
    → MAVLink bridge (on drone companion computer)
    → Drone adjusts in real-time
Total latency target: under 1.5 seconds voice to drone action
```

### "Follow me" with facial recognition
```
1. Camera feed from drone → inference on companion computer (EdgeTPU / Jetson Nano)
2. Face detected → lock target → publish tel.drone.skyhawk-1.target_locked
3. Zeus says "follow me"
4. Intent agent classifies → cmd.drone.skyhawk-1.follow { target_id: zeus-face }
5. Drone adjusts heading to keep Zeus centered in frame
6. Zeus says "go left 10 meters"
7. Instant heading adjustment via cmd.drone.skyhawk-1.adjust { dx: -10 }
```

### Combat override pattern (instant redirect)
Watch job monitors all active missions.
If Zeus speaks a steering command mid-mission, intent agent publishes with priority flag.
Subscriber sees priority flag, interrupts current mission, executes steering command.
Drone does not wait for mission to complete.
This is the "steer if necessary" model — always listening, always interruptible.

### Walkie-talkie / PTT integration
NATS has a voice subject namespace already defined:
`voice.{session}.audio.in` — PCM/Opus from mic or radio
STT runs on received audio.
Same pipeline. Works over radio relay.

---

## VR/AR Overlay with NATS

### How it works
VR/AR headset connects to NATS over WebSocket (ws://omos-mbp:8080).
Subscribes to:
- `tel.drone.*.gps` — drone positions → render 3D drone overlay in space
- `tel.robot.*.pose` — robot positions → render robot in physical space
- `tel.sensor.*.reading` — ambient data → floating data panels
- `evt.world.air.*` — aircraft → show planes in sky at correct bearing
- `cmd.*` — show active commands as visual indicators

No GPU on the server needed. GPU is in the headset.
NATS just streams the data. The headset renders it.

### Platforms
Meta Quest — WebXR, connects via ws:// to NATS, renders with Three.js or Babylon.js
HoloLens 2 — Unity + NATS WebSocket client, full mixed reality
Apple Vision Pro — Swift + nats.swift, spatial computing overlay
Custom (for field use) — Android-based AR + nats-py or nats.java

### The demo that goes viral
NATS streams drone GPS in real-time.
You put on the Quest.
You look toward the field.
A digital twin of your drone floats in the sky at the correct bearing and altitude.
Robot patrol path is drawn in AR on the ground.
Your AI agents appear as floating nodes in your space showing what they are doing.
World events (that aircraft 3km away) appear as a dot in the correct direction.

All of this is just NATS WebSocket subscriptions feeding a Three.js WebXR scene.
No special server. No special hardware beyond the headset.

---

## Anduril EagleEye — Integration Model

Anduril EagleEye is an AR combat helmet:
- Ballistic protection + AR visor
- AI mission command overlay
- Real-time battlefield awareness
- Connects to command network

For OneMind, you do not build this helmet.
You build the software layer it (or a civilian equivalent) connects to.
The integration pattern:
```
EagleEye / AR headset
    → connects to NATS WebSocket (same as any other client)
    → subscribes to tel.fleet.*.position
    → subscribes to evt.world.*
    → subscribes to evt.tak.cot.*
    → renders overlays based on NATS events
    → publishes operator position and status back to bus
```

Civilian equivalent now: Meta Quest 3 with passthrough AR.
Same WebSocket connection. Same NATS subjects.
Build the software bridge to Quest first, EagleEye-style hardware is the later step.

---

## OSSIM Simulator Scripts — What to Build First

Priority order for simulation scripts:

1. sim_drone.py — GPS patrol route, battery drain, event publishing
2. sim_robot.py — room-to-room navigation path, pose updates
3. sim_sensor.py — temperature/humidity/motion, configurable noise
4. sim_vehicle.py — GPS route replay from recorded tracks
5. sim_tak_operator.py — moving CoT contact on a map
6. sim_world_events.py — inject synthetic seismic/weather/conflict events for testing watch jobs
7. sim_aircraft_sitl.py — ArduPilot SITL integration for realistic flight simulation

Each script: ~50-100 lines of Python with nats-py.
Total to build: one focused day.

---

## Getting World Event Sources Directly (no WorldMonitor dependency)

| Feed | Direct Source | Method | Cost |
|------|---------------|--------|------|
| ADS-B (aircraft) | ADS-B Exchange API, FlightAware, OpenSky Network | REST polling or Websocket | Free tier available |
| AIS (vessels) | AISHub, MarineTraffic, aisstream.io | WebSocket stream | Free with account |
| Seismology | USGS Earthquake API | REST poll every 60s | Free |
| Wildfires | NASA FIRMS API | REST poll every 15min | Free |
| Weather | OpenWeatherMap, Open-Meteo | REST poll | Free |
| Conflict | ACLED API | REST poll daily | Free for researchers |
| Military | ADS-B Exchange (unfiltered) | REST | Free |
| Space weather | NOAA SWPC | REST | Free |

For each: write a small Python poller → transform to geo-envelope → publish to NATS.
Same pattern as WorldMonitor bridge but direct. No dependency on that deployment.
~100 lines per feed.

---

## ATAK VR Status

CivTAK site confirms: ATAK can run on Oculus Quest with AR passthrough mode.
The ATAK client floats in front of you and you can see your real environment through it.
This is not full 3D spatial AR — it is the 2D ATAK map rendered in the headset passthrough view.
Full spatial TAK AR (map overlaid on physical world at correct geo-coordinates) requires
additional plugin development or a custom WebXR client subscribing to NATS directly.

For OneMind: build the WebXR client approach.
Better UX than the floating 2D map. Fully controlled by you. Native NATS connection.

---

## Platform UI Decision — Flutter vs React/TypeScript

Synadia's own stack:
- nex: Go
- orbit.js: TypeScript (utility toolkit)
- insights: AI agent skills (not a UI)
- Agent Web UI (the Synadia agent dashboard): TypeScript + Vue 3 + Vite (you have this running at :3300)
- Platform control plane UI: TypeScript + React (their commercial product)

The honest Flutter vs TypeScript/React assessment:

| Factor | Flutter | TypeScript + React |
|--------|---------|-------------------|
| Mobile (iOS/Android) | Best — single codebase | Requires React Native or separate app |
| Desktop (Mac/Win/Linux) | Good | Electron or Tauri |
| Web | Works but not native feel | Best — this is its home turf |
| Data-dense dashboards | Harder — custom widgets | Easy — recharts, D3, AG Grid all native |
| Map (MapLibre, Deck.gl) | maplibre_gl works but limited | Full ecosystem, Deck.gl for 3D |
| 3D / WebXR | Very limited | Three.js, Babylon.js, native WebXR |
| Real-time streams | Good | Excellent — native WebSocket, NATS.ws |
| Component ecosystem | Smaller | Massive — MUI, shadcn, Radix, etc. |
| Build speed with AI | Moderate | Fast |
| Synadia compatibility | Requires dart_nats | Native — same stack as agent web UI |

Recommendation:
Keep Flutter for the mobile companion app — phones and tablets in the field.
Build the main ops dashboard / control plane in TypeScript + React.
This mirrors exactly what Synadia does and what every serious ops UI does.
React + shadcn/ui + Recharts + MapLibre + Three.js + NATS.ws covers everything.
The 3D map, VR overlay, agent dashboard, stream viewer — all are native in this stack.
Flutter trying to do dense infrastructure UI is fighting against its strengths.

---

## Scraping Synadia Platform for Reference

On the free trial question:
Synadia's control plane UI is a web app. With a trial account you can log in and inspect it.
The HTML, CSS, and JavaScript are all visible in browser DevTools.
Network tab shows every API call they make and the exact request/response structure.
This tells you exactly what endpoints they query and what the data shapes are.
You do not need to copy their code. You need to understand their data model and UX patterns.
That is 100% legal — you are looking at a web app in a browser, same as any user would.

What to extract from the trial:
- Every page and its URL structure
- Every API call (Network tab) — what NATS endpoints they query
- The data shapes returned
- The UI patterns for streams, consumers, accounts, connections
- What their "fleet" view looks like
- What their "insights" queries look like

Then you build your own version with your own code.

---

## Valuation — If You Go All Out

Comparable companies and their valuations at equivalent stages:

| Company | What They Do | Valuation |
|---------|-------------|-----------|
| Palantir | Data operations platform for government + enterprise | $41B (2024 revenue $2.8B) |
| Anduril | Defense hardware + software (autonomous systems) | $30-60B (raised $1.5B Series F) |
| Synadia | NATS commercial layer only (no hardware, no UI for end users) | $25M Series B raised (early) |
| Home Assistant / Nabu Casa | Consumer smart home open source + cloud | ~$50-100M (bootstrapped, private) |
| Grafana Labs | Open source observability + cloud | $6B |
| HashiCorp | Open source infrastructure tools + enterprise | Acquired $6.4B |
| n8n | Open source workflow automation + cloud | $170M ARR trajectory |
| DJI | Drone hardware + software ecosystem | ~$15-20B (private) |
| Axon | Police hardware (Taser) + software (cloud) | $45B public company |

The OneMind product is closest to:
**Anduril (hardware + software, sovereign) meets Palantir (data ops) meets Home Assistant (open, personal) meets Grafana (open source + cloud model)**

That category does not exist as a single company today.

Conservative path (open source + cloud, no hardware):
Series A at $20-50M valuation with $2-5M ARR
Series B at $150-300M valuation with $15-30M ARR
Exit or IPO at $1-3B with $100M+ ARR

Aggressive path (hardware + software bundles, fleet operators, defense-adjacent):
Think Anduril trajectory.
Series A at $50-100M
Series B at $500M-1B
Growth stage at $3-10B+
The hardware creates a moat the pure-software companies cannot replicate.

The community + live build angle accelerates all of it.
Home Assistant reached 500K+ active users with zero traditional marketing.
If you show people a nervous system for the AI physical world, live, as you build it —
the community builds before the product is finished.
That community becomes your Series A story.

The drones and robots you build and sell are the Raspberry Pi of the OneMind ecosystem.
Pre-wired. Pre-configured. Plug in and your world comes alive.
That is a product people will pay $2,000-10,000 for.
Not because of the hardware cost — because of what it connects to.

---

## Next Steps for OSSIM

1. Build sim_drone.py and sim_robot.py — 1 day
2. Get TAK server running on omos-cloud with CoT bridge
3. Wire Flutter app /fleet/map to subscribe to tel.> subjects
4. Connect WorldMonitor real feeds as Layer 1
5. Run the full demo: voice command → Legacy → NATS → simulated drone moves → TAK map updates → ClickHouse records it
6. Record it. That is the first video.

The SIM is real. The bus is real. The intelligence is real.
Only the hardware is software.
And when the hardware arrives, you remove the sim scripts and nothing else changes.

---

*Document: ONEMIND-OSSIM.md — created 2026-06-14*
*Location: vault root — accessible to all agents*
