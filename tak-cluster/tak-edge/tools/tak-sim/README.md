# TAK-SIM — Autonomous Operations Simulation Stack

Full simulation of AI-powered tactical operations without any physical hardware.
Two virtual drones fly real GPS routes, geofences are drawn live on ATAK,
and the AI agent autonomously responds to breach events.

## What this gives you

- **Live drone tracks on ATAK** — 2 virtual UAVs moving on your map in real time
- **Geofenced areas drawn on map** — geocircle shapes visible to all ATAK clients
- **Autonomous breach response** — drone enters zone → hostile marker placed, response drone re-tasked
- **NLP map commands** — publish to NATS, map updates instantly
- **Full NATS backbone** — every event is observable, replayable, extensible
- **Hermes integration** — optionally route decisions through your local AI gateway

## Quick Start

### 1. Start infrastructure

```bash
cd tak-sim
docker compose up -d nats freetakserver
```

Wait ~20 seconds for FreeTAKServer to initialize.

### 2. Connect ATAK

On your ATAK device/emulator:
- Settings → Network → Add Server
- **Host:** your Mac's IP (e.g. 192.168.1.x — run `ipconfig getifaddr en0`)
- **Port:** 8087
- **Protocol:** TCP Streaming

### 3. Install dependencies and run

```bash
npm install
npm run sim
```

You'll immediately see:
- Two drone icons (SENTINEL-1, REAPER-2) appear and start moving
- Two red geofence circles (ALPHA, BRAVO) drawn on the map
- Three waypoint markers (CP-ALPHA, CP-BRAVO, FOB-MAIN)

### 4. Trigger a breach

As the drones fly their routes, they'll naturally enter the geofence zones.
Or manually trigger one:

```bash
nats pub ops.breach.zone-alpha '{
  "zoneId":"zone-alpha",
  "zoneName":"ALPHA",
  "entityId":"TEST-INTRUDER",
  "breachType":"INSIDE",
  "distanceM":0,
  "entityLat":38.8995,
  "entityLon":-77.038,
  "action":"launch_drone",
  "actionTarget":"REAPER-2",
  "severity":"HIGH",
  "timestamp":"2026-06-22T00:00:00Z"
}'
```

Watch on ATAK: red hostile marker appears, REAPER-2 gets re-routed.

### 5. NLP map commands

```bash
nats pub ai.tak.nlp "add hostile marker at 38.899, -77.037"
nats pub ai.tak.nlp "draw 400m geofence at 38.903, -77.032"
nats pub ai.tak.nlp "add waypoint at 38.895, -77.044"
```

## With Hermes AI

```bash
USE_HERMES=true HERMES_URL=http://localhost:8642 HERMES_PROFILE=tactical npm run sim
```

Every breach decision is routed through Hermes before action is taken.

## Configuration

| Env var | Default | Description |
|---|---|---|
| `NATS_URL` | `nats://localhost:4222` | NATS server |
| `TAK_HOST` | `127.0.0.1` | FreeTAKServer host |
| `TAK_PORT` | `8087` | TCP port |
| `TAK_PROTO` | `tcp` | `tcp` or `udp` |
| `HERMES_URL` | `http://localhost:8642` | Hermes gateway |
| `HERMES_PROFILE` | `default` | Hermes profile name |
| `USE_HERMES` | `false` | Enable AI decisions via Hermes |

## Edit the AO (area of operations)

In `src/index.ts`:
- `GEOFENCES` — zone definitions, actions, sizes
- `WAYPOINTS` — named points on the map

In `src/drone-simulator.ts`:
- `DEFAULT_DRONES` — drone IDs, speeds, waypoint routes

Change `lat/lon` to your actual location.

## Architecture

```
NATS JetStream (backbone)
    ↑                 ↓
DroneSimulator   GeofenceEngine
(GPS routes)     (breach detect)
    ↓                 ↓
CotRelay ←── AITakAgent (ai.tak.nlp, ops.breach.>)
    ↓              ↑
FreeTAKServer   Hermes (optional)
    ↓
ATAK (your device)
```

## NATS subjects

| Subject | Publisher | Subscriber | Data |
|---|---|---|---|
| `constellation.telemetry.{id}` | DroneSimulator | GeofenceEngine, AITakAgent | TelemetryEnvelope |
| `ops.breach.{zone-id}` | GeofenceEngine | AITakAgent | BreachEvent |
| `ai.tak.nlp` | you | AITakAgent | plain text command |
| `ai.tak.decisions` | AITakAgent | logging/Cordum | decision record |
| `ai.tak.log` | AITakAgent | dashboards | log line |
| `commands.{drone-id}.guidance` | AITakAgent | SentinelX/nats-ros-connector | GuidanceCommand |
| `tak.cot.outbound` | CotRelay | any listener | raw CoT XML |
