# OTS + ATOC — Operational Runbook
> **Owner:** Zeus Delacruz | **Agent:** Legacy | **Status:** Active  
> **Last updated:** 2026-06-16  
> **Linked doctrine:** `~/onemind/00-24 SO (Sovereign Operations)/03 Infrastructure (The Body)/TACMAP-COP-STANDARD.md`

---

## Prime Thesis

> **NATS is the nervous system. TAK is the map layer. Legacy is the brain.**  
> One Mind + ATAK = a live sovereign command center. Every agent, sensor, asset, and alert — visible as icons on a tactical map. Bidirectional. Voice-controlled. Fully offline-capable. Tailscale-bound.

This doc is the operational runbook. The architecture doctrine lives in `TACMAP-COP-STANDARD.md`. This is the *build and run* reference.

---

## 1. Current Stack Status

### Services (as of 2026-06-16)

| Service | Status | Port | Notes |
|---|---|---|---|
| OpenTAKServer (OTS backend) | ✅ Running | 8080 (HTTP), 8443 (HTTPS) | PID managed by launchd `io.onemind.opentakserver` |
| OpenTAKServer UI | ✅ Running | 5174 | Node/Vite, launchd `io.onemind.opentakserver-ui` |
| PostgreSQL 16 | ✅ Running | 5432 | Homebrew, DB: `ots`, user: `ots` |
| RabbitMQ | ✅ Running | 5672 | Homebrew, was down — fixed by launchctl bootout/bootstrap |
| MediaMTX (video relay) | ✅ Configured | 9997 (API) | RTSP/WebRTC relay for drone/cam video |
| NATS | ✅ Running | 4222 | agenttap pipeline — the event bus |
| TAK SSL streaming | ✅ Listening | 8089 | ATAK client enrollment port |
| TAK TCP streaming | ✅ Listening | 8088 | Unencrypted fallback |

### Config File
```
~/ots/config.yml
```

### Log Files
```
~/ots/logs/ots-stdout.log
~/ots/logs/ots-stderr.log
```

### Data Directory
```
~/ots/                    — root
~/ots/ca/                 — Certificate Authority (certs for client enrollment)
~/ots/uploads/            — data packages, MBTiles maps, KML/KMZ files
~/ots/mediamtx/           — video relay config
```

### Known Non-Issues
- `Mumble auth not supported on this platform` — macOS limitation, not a bug, safe to ignore
- RabbitMQ had a broken launchd plist — fixed manually with `bootout/bootstrap`. Will survive reboots now.

### OTS Version
```
opentakserver 1.7.12
pytak 7.0.2
Python 3.14.5
```

### Web UI Access
```
http://localhost:5174          — admin console + basic map
http://100.114.142.55:5174     — Tailscale access from any node
```

---

## 2. What You CANNOT See Without a TAK Client

The OTS web UI (`localhost:5174`) is the **admin console only**. It shows:
- Connected clients
- CoT events (raw)
- Basic map with asset positions
- Data package management
- User/cert management

**Everything below requires ATAK on Android (or iTAK on iOS):**

| Feature | Requires TAK Client |
|---|---|
| Drawing tools (lines, polygons, threat rings, routes) | ✅ Yes |
| Offline map rendering with MBTiles | ✅ Yes |
| Real-time blue force tracking (all users moving on map) | ✅ Yes |
| GeoChat (map-linked messaging) | ✅ Yes |
| 9-Line / CASEVAC forms | ✅ Yes |
| SALUTE reports | ✅ Yes |
| Video feeds in-map | ✅ Yes |
| Data package sync (maps, icons, overlays) | ✅ Yes |
| Push-to-talk via Mumble | ✅ Yes (Android only) |

**ATAK is just an app.** Install it. No root. No MDM. No device takeover. Lives alongside everything else on the tablet. Standard Android permissions (GPS, camera, mic). Uninstall like any app.

---

## 3. ATAK Client — Setup on Android Tablet

### Install
- **Google Play Store:** Search "ATAK-CIV" — free, ~90% of features, no .mil email required
- **Full version (TAK.gov):** Requires .mil/.gov email — skip for now, Play Store version is sufficient

### Connect to OTS

#### Option A — Certificate Enrollment (recommended, secure)

1. On Mac, go to `http://localhost:5174` → login → **Users** → create a user for the tablet
2. Go to **Certificates** → generate a client cert for that user → download the `.p12` file
3. Transfer the `.p12` to the tablet (AirDrop won't work — use Tailscale file transfer, email, or USB)
4. On ATAK: Settings → Network → TAK Servers → Add Server
   - Host: `100.114.142.55` (your Tailscale IP)
   - Port: `8089`
   - Protocol: `SSL`
   - Import the `.p12` cert when prompted
5. Connect — your tablet appears as a blue dot on the OTS map immediately

#### Option B — TCP (quick test, no certs)
- Port: `8088`
- Protocol: `TCP`
- No cert needed — use this to verify connectivity first

#### Your OTS CA Password
```
atakatak   (set in config.yml → OTS_CA_PASSWORD)
```
Change this. It's the default and it's public knowledge.

---

## 4. Offline Maps — MBTiles Setup

### What it is
MBTiles = a single `.sqlite` file containing pre-rendered map tiles for a geographic area. OTS serves these to ATAK clients on enrollment/sync. Once on the tablet, works with **zero internet**.

### How to Get MBTiles

**Option A — MapTiler (paid, best quality)**
- https://www.maptiler.com/data/
- Download by region, export as MBTiles

**Option B — MOBAC (free, any source)**
- Mobile Atlas Creator: https://mobac.sourceforge.io/
- Pick area → export as MBTiles
- Sources: OpenStreetMap, USGS topo, satellite imagery

**Option C — OpenMapTiles (free, self-hosted)**
- https://openmaptiles.org/
- Download pre-built regional extracts
- US states available as individual files

**Option D — Quick test tile (US, small)**
```bash
# Download a small test MBTiles (10MB, sample area)
curl -L "https://github.com/mapbox/mapbox-gl-js/raw/main/test/fixtures/tiles/raster.mbtiles" \
  -o ~/Downloads/test-map.mbtiles
```

### Upload to OTS
1. `http://localhost:5174` → **Data Packages** → Upload
2. Select the `.mbtiles` file
3. OTS adds it to the sync pool
4. Next time an ATAK client connects → map downloads automatically

### Config flag (already enabled)
```yaml
OTS_PROFILE_MAP_SOURCES: true   # in ~/ots/config.yml — already set ✅
```

---

## 5. ADS-B (Live Aircraft Tracking)

### What it does
Pulls live aircraft positions from an API → converts to CoT → appears as aircraft icons on ATAK map, updating every 60 seconds.

### Current Config (WRONG — needs your actual location)
```yaml
OTS_ADSB_LAT: 40.744213    # ← currently set to NYC (Times Square area)
OTS_ADSB_LON: -73.986939   # ← change this to your actual location
OTS_ADSB_RADIUS: 10        # miles radius — increase if you want more coverage
OTS_ADSB_API_URL: https://api.airplanes.live/v2/point/
OTS_ADSB_GROUP: ADS-B
```

### Fix It
```bash
# Edit ~/ots/config.yml
# Update OTS_ADSB_LAT and OTS_ADSB_LON to your real coords
# Then restart OTS:
launchctl kickstart -k gui/501/io.onemind.opentakserver
```

### API Options
- **airplanes.live** — free tier, already configured, no key needed
- **AirNav RadarBox / FlightAware** — paid, better coverage, set `OTS_ADSB_API_KEY`

---

## 6. AIS (Live Ship Tracking)

### What it does
Pulls vessel positions from AISHub → CoT → ship icons on ATAK map.

### Setup
1. Create free account at https://www.aishub.net/register
2. Note your username
3. Edit `~/ots/config.yml`:
```yaml
OTS_AISHUB_USERNAME: your_username_here
OTS_AISHUB_NORTH_LAT: 45.0    # bounding box for your area of interest
OTS_AISHUB_SOUTH_LAT: 25.0
OTS_AISHUB_EAST_LON: -65.0
OTS_AISHUB_WEST_LON: -90.0
```
4. Restart OTS

---

## 7. Video Streaming (MediaMTX)

### What it does
Any RTSP/RTMP source (drone, GoPro, IP camera, phone) → OTS relays → clickable video icon on ATAK map. Tap the icon → full video feed in ATAK.

### MediaMTX is already running
- API: `http://localhost:9997`
- Auth token: `lvJOtKRXNHidug9kp9gXfUEBT7-MFQ` (in config.yml)

### Connect a source
- Point any RTSP camera at: `rtsp://100.114.142.55:8554/your-stream-name`
- ATAK clients will see a video icon appear on the map at that stream's position

### DJI Drone
- Use the FreeTAKUAS plugin in ATAK or the DJI RTMP output
- Stream → MediaMTX → appears on map as UAV icon with live video

---

## 8. OTS Plugin Ecosystem

### Plugins Flag (already enabled)
```yaml
OTS_ENABLE_PLUGINS: true    # already set ✅
OTS_PLUGIN_REPO: https://repo.opentakserver.io/brian/prod/
```

### Install any plugin
```bash
pipx inject opentakserver <plugin-package-name>
# Then restart OTS
launchctl kickstart -k gui/501/io.onemind.opentakserver
```

### Available Plugins

| Plugin | Source | What it does |
|---|---|---|
| `ots-aisstream-plugin` | Official repo | Real-time AIS via AISStream.io websocket (better than AISHub) |
| `ots-skyfi-plugin` | Official repo | SkyFi satellite imagery integration |
| `ots-arcgis-feed` | GitHub: teamvelociraptor | Pull ArcGIS feature layers as live CoT |
| `OTS-Debugger-Plugin` | GitHub: MMolineux | CoT inspection/debugging in UI |
| `OTS-Plugin-Template` | GitHub: brian7704 | Scaffold for building custom plugins |

### Building a Custom Plugin
Plugins are Flask blueprints. Minimal structure:
```
ots-myplugin/
├── ots_myplugin/
│   ├── __init__.py    # register blueprint with OTS
│   ├── plugin.py      # routes + CoT logic
│   └── models.py      # DB models (optional)
├── setup.cfg
└── pyproject.toml
```

Core of a plugin:
```python
from flask import Blueprint
blueprint = Blueprint('myplugin', __name__)

@blueprint.route('/api/myplugin/data')
def get_data():
    return {"status": "ok"}
```

Complexity: ~100-200 lines for a full functional plugin. Low barrier.

---

## 9. NATS → OTS CoT Bridge (THE BIG ONE)

### What This Is
A bridge that subscribes to NATS subjects, converts messages to CoT XML, and injects them into OTS's CoT pipeline. Result: **anything on your NATS bus appears as a live icon on your ATAK map.**

### Why This Is Huge
- Every One Mind agent, sensor, asset, alert → visible on map in real-time
- Bidirectional: ATAK sends CoT back → OTS → NATS → Legacy receives and acts
- Your tablet becomes a live command interface for the entire One Mind system

### Architecture
```
One Mind NATS Bus (port 4222)
         │
         ▼
  omos-tak bridge (Python, ~150 lines)
  - subscribes: onemind.cot.> and tel.{class}.{id}.position
  - converts to CoT XML via pytak
  - injects into OTS CoT pipeline
         │
         ▼
  OpenTAKServer (port 8089)
         │
         ▼
  ATAK on tablet — live map
```

### NATS Subject Taxonomy (from TACMAP doctrine)

**Outbound (One Mind → TAK map)**
| Subject | TAK Icon | Use Case |
|---|---|---|
| `onemind.location.zeus` | Blue force — person | Zeus moving on map |
| `onemind.location.home` | Friendly — structure | Home base fixed marker |
| `onemind.location.vehicle` | Friendly — vehicle | Car/motorcycle |
| `onemind.agent.legacy` | Equipment/sensor | Legacy agent heartbeat |
| `onemind.agent.spartan` | Equipment/sensor | Spartan agent status |
| `onemind.agent.oracle` | Equipment/sensor | Oracle agent status |
| `onemind.alert.geofence` | Unknown/hostile | Geofence breach alert |
| `onemind.alert.>` | Hostile marker | Any triggered alert |
| `onemind.sensor.>` | Equipment | Any sensor feed |
| `tel.robot.{id}.position` | Ground vehicle | Robot position |
| `tel.drone.{id}.position` | Air platform | Drone position |
| `evt.world.air` | Aircraft | ADS-B feed |
| `evt.world.sea` | Vessel | AIS feed |

**Inbound (TAK → NATS bus)**
| TAK Action | NATS Subject | Legacy Receives |
|---|---|---|
| Drop waypoint marker | `tak.cot.field.{uid}` | `{type: waypoint, lat, lon}` |
| Draw threat ring | `tak.cot.field.{uid}` | `{type: threat_zone, polygon}` |
| Send GeoChat message | `tak.cot.field.{uid}` | `{type: chat, text, location}` |
| Mark point of interest | `tak.cot.field.{uid}` | `{type: POI, uid, coords}` |
| Trigger CASEVAC form | `tak.cot.field.{uid}` | `{type: emergency_form}` |
| Client goes stale/offline | `tak.cot.field.{uid}` | `{uid: zeus, stale: true}` |

### CoT Message Format (what pytak generates)
```xml
<event version="2.0"
       uid="onemind-zeus-001"
       type="a-f-G-U-C"
       time="2026-06-16T12:00:00Z"
       start="2026-06-16T12:00:00Z"
       stale="2026-06-16T12:05:00Z"
       how="m-g">
  <point lat="40.7442" lon="-73.9867" hae="0" ce="9999999" le="9999999"/>
  <detail>
    <contact callsign="ZEUS"/>
    <track course="270" speed="0"/>
  </detail>
</event>
```

### MIL-STD-2525 Type Codes (what appears as what icon)
| Code | Icon Type | Use For |
|---|---|---|
| `a-f-G-U-C` | Friendly ground unit | Zeus, agents, people |
| `a-f-G-E-V` | Friendly ground vehicle | Cars, robots |
| `a-f-A-M-F-Q` | Friendly drone/UAV | Drones |
| `a-f-G-I-U-T-C` | Friendly installation | Home, fixed nodes |
| `a-f-G-U-C-I` | Friendly C2/intel | AI agents |
| `a-h-G` | Hostile ground | Threat/alert markers |
| `a-u-G` | Unknown ground | Unidentified assets |

### Existing Repo
```
~/omos-tak/
```
Check this — it may already have bridge scaffold code.

### Build Plan (T3 execution)
1. Check `~/omos-tak/` for existing code
2. Scaffold pytak NATS subscriber
3. Map NATS message fields to CoT geo-envelope
4. Inject into OTS via TCP on port 8088 (unencrypted, local)
5. Test: publish a fake position to NATS → see icon on ATAK map
6. Wire up real subjects: `onemind.location.zeus`, Home Assistant entities, agent heartbeats
7. Add inbound handler: OTS CoT → NATS publish

---

## 10. Legacy as ATOC — Autonomous Tactical Operations Controller

### What ATOC Is
Legacy monitoring the full NATS event stream 24/7 and acting autonomously on tactical events. Not reactive — **predictive and command-capable.**

### Full Architecture
```
┌─────────────────────────────────────────────────────┐
│                  LEGACY (ATOC)                       │
│                                                      │
│  Sees:   Full NATS event stream (all subjects)       │
│  Hears:  Voice via Mumble PTT or Telegram voice note │
│  Speaks: TTS → your earpiece or tablet speaker       │
│  Acts:   HA automations, agent tasks, CoT markers,   │
│          geofence triggers, mission commands          │
└──────────────────────┬──────────────────────────────┘
                       │ NATS Bus (port 4222)
        ┌──────────────┼──────────────────┐
        ▼              ▼                  ▼
      ATAK           Home             One Mind
   (TAK map)       Assistant           Agents
      │            185+ entities    Spartan/Oracle/
      │                             Eden/Empire/etc
      ▼
  Tablet display
  Voice in your ear
  Command interface
```

### ATOC Event Loop — What Legacy Does Autonomously

| NATS Event | Legacy Action |
|---|---|
| Zeus location changes | Tracks movement, detects anomaly, voice check-in if off-route |
| Geofence breach (from HA) | Drops red threat marker on map, fires Guardian agent, voice alert |
| You drop "RECON" POI in ATAK | Receives CoT, fires Oracle agent, briefs you via voice |
| Agent heartbeat stops | Drops yellow marker at last known, alerts on TAK + voice |
| CASEVAC form triggered | Full emergency protocol fires, contacts, HA alert |
| Home sensor anomaly | Marker on home node, HA automation, voice brief |
| Vehicle leaves geofence | Updates vehicle marker, logs, optional alert |
| Unknown CoT appears on network | Flags as unknown, Guardian agent investigates |

### ATOC Autonomous Actions by Domain

**Situational Awareness**
- Maintain current picture of all assets on NATS bus
- Detect stale/missing nodes (heartbeat monitoring)
- Cross-reference world events (ADS-B, AIS) with inner ring positions

**Command & Control**
- Accept voice commands via PTT → execute as NATS commands
- Drop mission waypoints on map programmatically
- Fire agent tasks from map events
- E-STOP robots/drones via `cmd.{class}.{id}.estop`

**Home Integration (HA → TAK)**
- 185 HA entities as potential map markers
- Geofence zones drawn as TAK overlays
- Door/sensor events → CoT alerts on map
- Presence detection → Zeus location on map

**Communications**
- Voice out via TTS (already have this) to tablet speaker or earpiece
- GeoChat messages to ATAK group
- Telegram alerts for critical events

### Voice Layer — Two Options

**Option A: Mumble PTT (TAK-native, full military radio feel)**
- ATAK has built-in Mumble push-to-talk support
- Run Mumble server on Mac (via Homebrew)
- Legacy bridges: Mumble audio → Whisper STT → process → TTS → back to Mumble channel
- Key PTT on tablet → Legacy hears it → responds in your ear via earpiece
- Full military-style comms with AI on the other end
- Setup: `brew install mumble-server` + Whisper bridge (~100 lines Python)

**Option B: Telegram voice notes (simpler, already wired)**
- Send voice note to Legacy → Hermes Whisper STT → processes → TTS response back
- Already works today, zero new infrastructure
- Lower latency on setup, same cognitive loop

**Both can run simultaneously.** Mumble for field ops, Telegram for quick commands.

### Build Plan (T5 execution)
1. Deploy Mumble server: `brew install mumble-server`
2. Configure Mumble channel: `onemind-atoc`
3. Build Whisper bridge: `mumble audio → whisper STT → Legacy NATS → TTS → mumble`
4. Connect ATAK Mumble plugin to local server
5. Test: "Legacy, status report" → voice response in earpiece
6. Wire ATOC event loop: subscribe to all relevant NATS subjects
7. Add autonomous action handlers per the table above

---

## 11. Full Build Roadmap

### Phase T1 — Foundation (THIS WEEK)
- [ ] Fix ADS-B coords in `~/ots/config.yml`
- [ ] Change OTS CA password (default `atakatak` is public)
- [ ] Enroll Android tablet as ATAK client
- [ ] Upload first MBTiles (offline map) via OTS UI
- [ ] Verify aircraft appear on ATAK map (ADS-B)
- [ ] Create AISHub account → enable ship tracking

### Phase T2 — NATS Bridge (NEXT)
- [ ] Inspect `~/omos-tak/` for existing scaffold
- [ ] Build NATS→CoT bridge (~150 lines pytak)
- [ ] Test: fake Zeus position on NATS → icon on ATAK
- [ ] Wire `onemind.location.zeus` (phone GPS → NATS → TAK)
- [ ] Wire Home node as fixed map marker
- [ ] Wire inbound CoT → NATS (TAK actions → Legacy)

### Phase T3 — One Mind on Map
- [ ] Wire all Home Assistant entities with lat/lon to map
- [ ] Add agent heartbeat markers (Spartan, Oracle, etc.)
- [ ] Wire geofence zones: HA → CoT → TAK overlay
- [ ] ADS-B + AIS on same map canvas as inner ring

### Phase T4 — ATOC Event Loop
- [ ] Build Legacy ATOC subscriber (all NATS subjects)
- [ ] Autonomous event handlers (geofence, heartbeat, anomaly)
- [ ] Map command interface: TAK action → NATS command → agent task
- [ ] HA bidirectional: TAK geofence draw → HA automation

### Phase T5 — Voice Layer
- [ ] Deploy Mumble server
- [ ] Build Whisper STT bridge
- [ ] Wire TTS output to Mumble channel
- [ ] ATAK PTT → Legacy hears → responds in ear
- [ ] Voice-dispatched missions: "send drone to north field" → `cmd.drone.1.goto`

---

## 12. Key Commands Reference

### Restart OTS backend
```bash
launchctl kickstart -k gui/501/io.onemind.opentakserver
```

### Restart RabbitMQ (if it dies again)
```bash
launchctl bootout gui/501 ~/Library/LaunchAgents/homebrew.mxcl.rabbitmq.plist
launchctl bootstrap gui/501 ~/Library/LaunchAgents/homebrew.mxcl.rabbitmq.plist
```

### Check all OTS services
```bash
launchctl list | grep onemind
brew services list | grep -E "rabbit|postgres"
ps aux | grep opentakserver | grep -v grep
```

### Watch OTS logs live
```bash
tail -f ~/ots/logs/ots-stdout.log
tail -f ~/ots/logs/ots-stderr.log
```

### Install an OTS plugin
```bash
pipx inject opentakserver <plugin-name>
launchctl kickstart -k gui/501/io.onemind.opentakserver
```

### Edit OTS config
```bash
nano ~/ots/config.yml
# After changes:
launchctl kickstart -k gui/501/io.onemind.opentakserver
```

### Test OTS is responding
```bash
curl -s http://localhost:8080/api/health
# Returns "Bad Request" = alive (no /health endpoint, but server is up)
```

---

## 13. Security Hardening (Do This)

| Item | Current | Fix |
|---|---|---|
| CA Password | `atakatak` (public default) | Change in `config.yml` → `OTS_CA_PASSWORD` |
| TAK TCP port 8088 | Open (no auth) | Disable: `OTS_ENABLE_TCP_STREAMING_PORT: false` |
| OTS_HOSTNAME | `100.114.142.55` (Tailscale) ✅ | Good — not public IP |
| IP whitelist | `127.0.0.1`, `100.114.142.55`, `100.122.203.68` ✅ | Add any new Tailscale nodes |
| RabbitMQ | Default guest/guest creds | Change: `OTS_RABBITMQ_USERNAME/PASSWORD` |
| OTS DB | `ots:ots_password` | Change in `SQLALCHEMY_DATABASE_URI` |

---

## 14. Reference Links

| Resource | URL |
|---|---|
| OTS Docs | https://docs.opentakserver.io |
| OTS GitHub | https://github.com/brian7704/OpenTAKServer |
| ATAK-CIV (Play Store) | https://play.google.com/store/apps/details?id=com.atakmap.app.civ |
| pytak (CoT bridge lib) | https://github.com/snstac/pytak |
| CoT Developer Guide | https://nps.edu/documents/104517539/109705106/COT+Developer+Guide.pdf |
| CoT Type Codes (MIL-STD-2525) | https://freetakteam.github.io/FreeTAKServer-User-Docs/About/architecture/mil_std_2525 |
| MapTiler (MBTiles) | https://www.maptiler.com/data/ |
| MOBAC (free MBTiles generator) | https://mobac.sourceforge.io/ |
| AISHub (free ship tracking) | https://www.aishub.net/register |
| MediaMTX (video relay) | https://github.com/bluenviron/mediamtx |
| Meshtastic ATAK plugin | https://github.com/meshtastic/ATAK-Plugin |
| CloudTAK (web COP reference) | https://github.com/dfpc-coe/CloudTAK |
| TACMAP Doctrine | `~/onemind/00-24 SO (Sovereign Operations)/03 Infrastructure (The Body)/TACMAP-COP-STANDARD.md` |

---

## 15. The Vision — One Sentence

> You are in the field. Tablet on your wrist. Legacy in your ear. Every agent, sensor, asset, and alert in your life — visible as icons on a live tactical map. You tap a point, Legacy briefs it. You speak a command, Legacy executes it. Fully offline. Fully sovereign. That is One Mind with a map layer.

---

*Written by Legacy | One Mind ATOC System | 2026-06-16*
