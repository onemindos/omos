# OneMind — Master Inventory
> Generated: 2026-07-06 | Complete picture of every workload, reference repo, and next step.

---

## LIVE ON MAC RIGHT NOW

### NEX Workloads (running via launchd)

| Binary | Location | What it does | Status |
|---|---|---|---|
| `om-unity` | `nex/om-unity/target/release/om-unity` | ORP — 39 protocol adapters, 5,938 AIS ships, DuckDB+Kuzu, ORP-QL, built-in COP | ✅ LIVE port 9090 |
| `om-overwatch-api` | `nex/om-overwatch-api/` | REST API for overwatch state | ✅ LIVE |
| `om-overwatch-interceptor` | `nex/om-overwatch-interceptor/` | Air track detection | ✅ LIVE |
| `om-overwatch-relay` | `nex/om-overwatch-relay/` | CoT relay bridge | ✅ LIVE |
| `om-overwatch-state` | `nex/om-overwatch-state/` | State machine — air tracks only today | ✅ LIVE |
| `om-overwatch-video` | `nex/om-overwatch-video/` | MediaMTX RTSP/WebRTC bridge | ✅ LIVE |
| `om-bridge-homeassistant` | `nex/om-bridge-homeassistant/` | Home Assistant → NATS | ✅ LIVE |
| `om-bridge-constellation` | `nats/nex/om-bridge-constellation/` | CoT → Constellation Overwatch | ✅ LIVE |
| `om-bridge-mavlink` | `nats/nex/om-bridge-mavlink/` | MAVLink → NATS | ✅ LIVE |
| `om-callout` | `nats/nex/om-callout/` | Outbound notifications | ✅ LIVE |
| `om-http-gateway` | `nats/nex/om-http-gateway/` | HTTP → NATS gateway | ✅ LIVE |
| `om-nex-ui` | `nats/nex/om-nex-ui/` | NEX workload UI | ✅ LIVE |
| `om-streaming-bridge` | `nats/nex/om-streaming-bridge/` | Video/stream bridge | ✅ LIVE |
| `resgate` | `nats/nex/resgate/` | RES protocol gateway | ✅ LIVE |
| `om-task-nex` | `nex/om-task/` | Task management on NATS | ✅ LIVE |
| `om_presence_v2.py` | `nats/services/connect/` | Presence tracking on NATS | ✅ LIVE |
| `mediamtx` | homebrew | RTSP/WebRTC/HLS video server | ✅ LIVE |
| `orp` | `ai-agents/orp/target/release/orp` | ORP binary (separate build) | ⛔ DISABLED |
| `apex-sentinel` | `defense/apex-sentinel/` | Counter-UAS sentinel | ✅ LIVE (launchd) |
| `overwatch-vision` | `defense/om-overwatch-vision/` | Computer vision detection | ✅ LIVE (launchd) |

### Other Running Services

| Service | Port | What it does |
|---|---|---|
| `nats-server` | 4222 | Mac leaf node + full JetStream |
| `nex node zeus-mac` | — | NEX workload runner |
| `hermes gateway` | — | Legacy AI (AWS Bedrock) |
| `heynats` | 5000 | NATS web explorer |
| `nats-dashboard` | 5173 | NATS UI dashboard |
| `caddy` | — | atoc reverse proxy |
| `nats-surveyor` | 7777 | NATS metrics scraper |

---

## REFERENCE REPOS ON MAC — TO BE EXTRACTED & PORTED

### Root Level (`/Users/zeus/`)

| Folder | Size | What it has | What to extract | Then |
|---|---|---|---|---|
| `worldwideview/` | 56MB | CesiumJS 3D globe, plugin system, DataBus (typed event bus), Agent Bus (MCP→fly globe), Next.js 16 | Plugin architecture + DataBus pattern + Agent Bus | Deploy to DOKS as global map, delete local |
| `worldmonitor/` | 185MB | globe.gl + deck.gl, 56 map layer types, 500+ news feeds, CII country scoring, 6 site variants | Layer catalog (56 types) + CII scoring algorithm + feed list | Extract to omos-blueprints docs, delete |
| `osiris/` | ~50MB | MapLibre GL, 50+ API routes (flights, CCTV, seismic, fires, maritime, cyber, sanctions, OSINT toolkit, GDELT, satellites, space-weather) | Every API route = a NATS publisher microservice | Port one by one, delete when done |
| `osiris-audit/` | ~50MB | Exact duplicate of osiris | Nothing new | **Delete immediately** |
| `Crucix/` | ~10MB | 27 data source adapters + LLM briefing engine + alert system | Each adapter = a NATS microservice | Port one by one to omos-core/nats/services/ |
| `tools/constellation-overwatch/` | ~5MB | Go C4ISR binary, entity model, NATS JetStream | Already bridged — entity schema docs | Archive to omos-core/docs/, delete |
| `tools/cot2constellation/` | <1MB | Python bridge: TAK CoT → Constellation | Already running — move to omos-core | Move, delete from tools/ |
| `tools/pulsar/` | <1MB | Go MAVLink+RTSP+YOLO edge agent | Replaced by om-bridge-mavlink | Delete |
| `tools/dashboards/` | ~5MB | NATS dashboard UI | Already running via launchd | Delete |

### Crucix — 27 Data Source Adapters (each becomes a NATS microservice)

| Adapter | Domain | NATS Subject |
|---|---|---|
| `opensky` | Live flights ADS-B | `tel.adsb.*` |
| `adsb` | ADS-B direct feed | `tel.adsb.*` |
| `ships` | AIS maritime vessels | `tel.ais.*` |
| `acled` | Conflict/violence events | `det.conflict.*` |
| `firms` | NASA active fires | `det.fire.*` |
| `noaa` | Weather/severe events | `det.weather.*` |
| `space` | Satellite tracks, solar | `det.space.*` |
| `gdelt` | Global events database | `intel.gdelt.*` |
| `opensanctions` | OFAC SDN, sanctions | `intel.sanctions.*` |
| `ofac` | US Treasury sanctions | `intel.sanctions.*` |
| `cisa-kev` | Exploited vulnerabilities | `intel.cyber.*` |
| `cloudflare-radar` | Internet traffic intel | `intel.cyber.*` |
| `fred` | Federal Reserve economic data | `intel.finance.*` |
| `bls` | Bureau of Labor Statistics | `intel.finance.*` |
| `eia` | Energy information | `intel.finance.*` |
| `yfinance` | Stock markets | `intel.finance.*` |
| `comtrade` | UN trade data | `intel.finance.*` |
| `treasury` | US Treasury data | `intel.finance.*` |
| `usaspending` | US federal spending | `intel.finance.*` |
| `telegram` | Geoparsed Telegram posts | `intel.osint.*` |
| `bluesky` | BlueSky social feed | `intel.osint.*` |
| `reddit` | Reddit intel | `intel.osint.*` |
| `epa` | Air quality | `det.environment.*` |
| `safecast` | Radiation monitoring | `det.environment.*` |
| `kiwisdr` | Software defined radio | `det.sigint.*` |
| `reliefweb` | Humanitarian alerts | `intel.humanitarian.*` |
| `who` | WHO disease alerts | `intel.health.*` |
| `gscpi` | Supply chain pressure | `intel.finance.*` |
| `patents` | Patent filings | `intel.tech.*` |

### Osiris — API Routes (each becomes a NATS publisher)

```
flights          → tel.adsb.*
maritime         → tel.ais.*
earthquakes      → det.seismic.*
fires            → det.fire.*
weather          → det.weather.*
space-weather    → det.space.*
satellites       → det.space.*
conflicts        → det.conflict.*
frontlines       → det.conflict.*
cctv             → tel.cctv.*
cyber-threats    → intel.cyber.*
malware          → intel.cyber.*
gdelt            → intel.gdelt.*
live-news        → intel.news.*
markets          → intel.finance.*
country-risk     → intel.cii.*
region-dossier   → intel.region.*
osint/sanctions  → intel.sanctions.*
osint/cve        → intel.cyber.*
osint/dns        → intel.osint.*
osint/ip         → intel.osint.*
osint/whois      → intel.osint.*
osint/bgp        → intel.osint.*
osint/shodan     → intel.osint.*
osint/threats    → intel.osint.*
osint/leaks      → intel.osint.*
osint/crypto     → intel.sanctions.*
air-quality      → det.environment.*
```

### Constellation Overwatch Ecosystem (GitHub: Constellation-Overwatch org)

| Repo | What it is | Value |
|---|---|---|
| `constellation-overwatch` | Go C4ISR data plane, entity model, NATS JetStream, MapLibre dashboard | Entity schema = OneMind entity model |
| `mavlink2constellation` | MAVLink UDP → Constellation NATS | Pattern already in om-bridge-mavlink |
| `ffmpeg2constellation` | H.264/MPEG-TS video publisher | Pattern for video feeds |
| `gocv2constellation` | Go object detection → Constellation | Pattern for vision feeds |
| `vision2constellation` | Python YOLO/Moondream → Constellation | Pattern for edge AI |
| `aero-arc-relay2constellation` | MAVLink relay via NATS | Replaced by om-bridge-mavlink |
| `pulsar` | Edge sync agent (MAVLink+RTSP+YOLO) | Replaced by om-bridge-mavlink |
| `toolbelt` | IaC templates | Reference only |

### Special Projects

| Project | What it is | Status | Priority |
|---|---|---|---|
| **ORP / om-unity** | 45MB Rust binary. Palantir-grade fusion engine. 39 protocol adapters. DuckDB+Kuzu+RocksDB. ORP-QL query language. mTLS federation mesh. 5,938 live AIS ships. Built-in COP map. | Built, running port 9090. `omos.orp.plist` DISABLED. | 🔥 **#1** — enable + wire NATS → `tel.ais.*` |
| **AMOS** | Python/Flask multi-domain C2. 537 tests + CI. 3 demo scenarios (border patrol, swarm, SAR). Mission orchestration. 200+ API routes. TAK/MAVLink/ROS2/MQTT/DDS/Link-16. Kill web. SIGINT. Wargaming. | Not extracted. Not running. | 🔥 **#2** — extract to omos-legacy, wire to NATS |

---

## INSIDE omos-core/integrations/ — Already Extracted Reference Repos

### ai-agents/
- `orp` — ORP source (built at `ai-agents/orp/target/release/orp`)
- `NyxID` — identity/auth patterns
- `cordum` — infrastructure patterns
- `iron-concept` — concept patterns
- `synadia` — Synadia agent patterns (also in hive/)

### defense/
- `apex-sentinel` — counter-UAS sentinel (LIVE in `defense/apex-sentinel/`)
- `SentinelX` — extended sentinel patterns
- `constellation-overwatch` — C4ISR reference
- `cot2constellation` — TAK→Constellation bridge (LIVE in `tools/`)
- `bachmach-pqc-iot-sentinel` — PQC IoT security patterns

### iot-edge/
- `IOnode`, `WireClaw`, `eclss-digital-twin`, `peat`, `rb-sensor-simulator`, `otel-edge`, `qiki-dtmp`
- Patterns: IoT sensor → NATS, OTel edge, digital twin

### microservices/
- `eventhorizon` — Go event sourcing framework
- `faststream` — Python async NATS microservice framework
- `fathom` — analytics patterns
- `go-atlas`, `go-events`, `go-intelli` — Go service patterns
- `ddd-hexagonal-cqrs-es-eda` — DDD/CQRS/ES patterns

### nats/
- `autonats` — auto-generated NATS bindings
- `UNS-Design-Studio` — Unified Namespace design
- `Nats-blitz` — NATS performance testing
- `flink-connector-nats` — Flink→NATS connector
- `gn` — Go NATS utilities

### robotics/
- `dimos` — robot fleet management
- `nats-ros-connector` — ROS2↔NATS bridge
- `rhombus` — robotics patterns
- `LLM_ROS2_agent-1` — LLM+ROS2 agent

### hive/ (NATS ecosystem)
- `synadia-agents/` — official Synadia agent SDK
- `synadia-agents-labs/` — experimental agent patterns
- `orbit.go` — Go NATS utility library
- `jwt-auth-builder.go` — NATS JWT auth
- `callout.go` — notification library
- `hermes-nats-gateway` — Hermes NATS gateway

---

## TAK SERVER 5.7 — BUILT-IN (Don't rebuild these)

| Capability | TAK Built-in | Status |
|---|---|---|
| Video streaming | Built-in video server (RTSP/HLS) | ✅ Don't need separate media server for TAK video |
| Mission packages | Marti API — full REST | ✅ Use directly |
| Data sync | Mission packages, attachments | ✅ Use directly |
| KML/GeoJSON overlay | Push map layers directly | ✅ Use directly |
| Federation | Server-to-server CoT | ✅ Available |
| Chat | Built-in TAK chat | ✅ Don't need takchat plugin |
| CoT ingest | NatsBridgePlugin (YOUR CODE) | ✅ LIVE — bus→TAK |
| CoT tap | cot-tap (deployed on DOKS) | ✅ LIVE — TAK→bus |

---

## THE GAP — What's Not Wired Yet

1. **No global map** — WorldWideView not deployed. No 3D globe showing anything.
2. **om-unity not publishing to NATS** — 5,938 ships in ORP, none on `tel.ais.*`
3. **Crucix adapters not running** — 27 sources, zero NATS publishers
4. **Osiris routes not ported** — 50 API domains, zero NATS publishers
5. **AMOS not extracted** — mission orchestration layer missing
6. **om-overwatch air-only** — ships/ground/space/cyber invisible on ATAK
7. **ORP plist disabled** — separate ORP binary not running

---

## STEP-BY-STEP PLAN (in order, nothing skipped)

| Step | Task | What it unlocks |
|---|---|---|
| 1 | Enable `omos.orp.plist` + wire om-unity NATS adapter → `tel.ais.*` | 5,938 ships on the bus |
| 2 | Deploy WorldWideView to DOKS (`wwv.onemindos.dev`) + NATS→WebSocket→DataBus bridge | Ships on 3D globe |
| 3 | Wire om-unity → CoT → ATAK (via om-overwatch-tak) | Ships on ATAK map |
| 4 | Port Crucix `opensky` → NATS `tel.adsb.*` | Live flights on globe + ATAK |
| 5 | Port Crucix `firms` → NATS `det.fire.*` | Active fires on globe |
| 6 | Port Crucix `acled` → NATS `det.conflict.*` | Conflict zones on globe |
| 7 | Port Crucix `ships` → NATS (supplement om-unity) | More maritime sources |
| 8 | Port Crucix `noaa` + `space` → NATS | Weather + satellites |
| 9 | Port Osiris `earthquakes`, `satellites`, `conflicts` → NATS | More domains |
| 10 | Generalize om-overwatch (all entity types, not just air) | Everything on ATAK |
| 11 | Extract AMOS → omos-legacy, wire to NATS | Mission orchestration |
| 12 | Continue Crucix + Osiris adapters one by one | Full intelligence picture |
| 13 | Port worldmonitor CII algorithm → Legacy/Oracle | Country risk scoring |
| 14 | Wire Agent Bus (WorldWideView MCP) → Legacy | Legacy can fly the globe |

---

## LOOSE FILES TO CLEAN UP

| File | Action |
|---|---|
| `Backup-codes-onemind.zeus.txt` | Move to 1Password, delete file |
| `zeuss-macbook-pro.tail717ba0.ts.net.crt/.key` | Managed by Tailscale — delete |
| `MumbleAutomaticCertificateBackup.p12` | Move to 1Password |
| `DEFENSE-*.md` | Stale session docs — delete |
| `TAKEOUT_ANALYSIS.md`, `TCL_20_XE_Screen_Replacement.md` | Delete |
| `bookmarks_6_16_26.html` | Old browser export — delete |
| `nats-server.conf.bak.*` | Backed up to 1Password + omos-infra — delete |
| `legacy-*.sh` | Move to omos-legacy/ |
| `ai_integration_setup.md` | Move to omos-core/docs/ |
| `Nomad_Orchestration_Deep_Dive.md` | Move to omos-infra/docs/ |
| `IMG_5418.JPG`, `libduckdb.a` | Delete |

---

*Last updated: 2026-07-06 by Legacy*
