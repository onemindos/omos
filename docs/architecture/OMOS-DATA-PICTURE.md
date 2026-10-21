# OMOS DATA PICTURE — What feeds the map, and how it all comes together

> **Review status (2026-10-17):** active OneMind OS design document, moved from the repo root and due for review and update. Until it is updated, where it disagrees with live state, `docs/STACK.md` (generated from the cluster) and `CLAUDE.md` win.

**Date:** 2026-07-04 · For Zeus. Plain terms. One pipeline. Everything is omos.

---

## 1. THE ONE PIPELINE (this is the whole architecture in one line)

Every piece of data — a plane, a ship, a fire, a drone, your kid's location, a news
alert — follows the SAME path:

```
SOURCE  →  CONNECTOR  →  BUS (NATS stream)  →  PROJECTOR  →  SCREEN
(world)    (omos-feed-*)   (WORLD/INTEL/TEL...)   (omos-tak-projector,   (ATAK map,
                                                    omos-control-plane)    atoc, Grafana)
```

- **Source** = where data is born (an API, a sensor, a drone, a repo's feed list)
- **Connector** = small service that pulls/receives it and publishes to the bus (one per source family)
- **Bus** = NATS streams you already provisioned (they're the landing zones waiting for publishers — THIS is what they were waiting for)
- **Projector** = picks what's relevant and converts it (CoT for TAK, ui.* for atoc)
- **Screen** = ATAK/iTAK tactical map, atoc console, Grafana

**Key insight for your "massive tactical map":** you already own the last mile —
`takserver-nats-bridge` (subject `tak.cot.inject`) is LIVE and proven. Anything we
publish as CoT lands on your ATAK map. Enriching the map = adding connectors on
the left. No new map infrastructure needed. This is why the focus on data is right.

---

## 2. WHAT WE ALREADY HAVE (live on the bus today)

| Source | Connector (live) | Stream | On the map? |
|---|---|---|---|
| Drone detections (acoustic/RF) | omos overwatch sentinel | EVT/det.* | ✅ via overwatch→TAK |
| Visual detections (YOLO) | omos overwatch vision | det.air.visual | ✅ |
| Friendly UAV telemetry (MAVLink) | mavlink-nats-bridge | TEL | ✅ |
| TAK devices (iTAK/ATAK positions) | takserver-nats-bridge | COT | ✅ native |
| AIS ships (demo feed) | om-unity (ORP engine) | HTTP only today | ⚠️ NOT on bus yet — task open |
| Voice transcripts | omos-voice | VOICE | n/a |
| Agent prompts/events | Hermes/Legacy | AGENTS | n/a |
| Tasks/projects | om-task | task.*/KV | atoc |
| Cluster/system health | omos-control-plane MVP | ui.* | atoc Systems |

**That's ~9 source families. Your streams ENV, WORLD, INTEL, GEO, ZEUS, SAAS are
still empty — provisioned landing zones. Sections 3 fills them.**

---

## 3. WHAT WE EXTRACTED (the enrichment — this answers "extract the sources")

### From Crucix — 29 ready-made source adapters (`/Users/zeus/Crucix/apis/sources/`)
Working fetch code, zero cloud. Grouped by which stream they'll feed:

| Group → Stream | Sources (adapter files) |
|---|---|
| **Air picture → WORLD** | `adsb.mjs` (ADS-B Exchange), `opensky.mjs` (OpenSky) |
| **Sea picture → WORLD** | `ships.mjs` (AIS) |
| **Space picture → WORLD** | `space.mjs` (satellites/launches) |
| **Fires/hazards → WORLD** | `firms.mjs` (NASA FIRMS), `noaa.mjs` (weather), `safecast.mjs` (radiation), `epa.mjs`, `who.mjs` |
| **Conflict/unrest → INTEL** | `acled.mjs` (conflict events), `gdelt.mjs` (global events), `reliefweb.mjs` |
| **Cyber → INTEL** | `cisa-kev.mjs` (exploited vulns), `cloudflare-radar.mjs` |
| **Sanctions/finance → INTEL** | `ofac.mjs`, `opensanctions.mjs`, `treasury.mjs`, `fred.mjs`, `yfinance.mjs`, `bls.mjs`, `eia.mjs`, `comtrade.mjs`, `usaspending.mjs`, `gscpi.mjs`, `patents.mjs` |
| **Social/chatter → INTEL** | `telegram.mjs`, `reddit.mjs`, `bluesky.mjs` |
| **RF → WORLD** | `kiwisdr.mjs` (global SDR receivers) |

### From worldmonitor — the curation layer (`/Users/zeus/worldmonitor/src/config/`)
- **~568 curated feed URLs** (`feeds.ts`) — news/OSINT feeds with categories
- **Geo registries** (static gold): `bases-expanded.ts` (military bases), `chokepoint-registry.ts` (maritime chokepoints), `airports.ts`, `apt-groups.ts` (threat actors), `ai-datacenters.ts`, `commodity-geo.ts` — these become **map overlay layers**
- **`basemap-styles.ts` / `basemap.ts`** — MapLibre basemap configs (feeds your sovereign map later)
- Source-failure resilience patterns (tests show how they handle dead feeds)

### From osiris — briefing/analyst workflow (`intel/server.js` + prompts)
- The **AI briefing pattern**: sources → LLM brief → structured report. This is a
  Hermes JOB, not a new service: Hermes reads INTEL stream → publishes daily brief.
  Crucix has the same idea (`BRIEFING_PROMPT.md`) — extract the prompts, done.

### From the pasted research (adjudicated in ONEMIND-MASTER-PICTURE.md)
- **AMOS** → mission contract patterns for future omos-mission. Pattern only.
- **tak-ml** → model packaging pattern for edge ML later. NOT a core. GPL + Java.
- **tak-gpt** → idea: Hermes as chat contact in iTAK (future omos-tak-chat organ).
- **ORP** → already ours (om-unity). Needs its NATS adapter turned on.
- **constellation-overwatch** → entity ontology reference. License bars adoption.
- **NIZAM/Palisade/ATAKRR/xTAK/etc.** → skip; nothing real to take.
- **TAK_OS/TAKpilot** → product-shape references (already driving atoc design).

**After this extraction the repos have given us everything they're going to give:
source adapters, feed catalogs, geo layers, briefing prompts, product shapes.**

---

## 4. WHERE IT GOES — the connector plan (few services, not source-count)

We do NOT build 29 connectors. We build **THREE**, each loading many adapters:

| New service | Loads | Publishes to | Why one service |
|---|---|---|---|
| **omos-feed-world** | adsb, opensky, ships, space, firms, noaa, safecast, kiwisdr | `evt.world.*` → **WORLD** stream (capped, short-lived by design) | All "moving world picture" pollers, same shape: poll → normalize → publish |
| **omos-feed-intel** | acled, gdelt, cisa-kev, sanctions group, telegram/reddit/bluesky, + worldmonitor's 568 feeds | `intel.*` → **INTEL** stream (curated, longer retention) | All "alerts & analysis" sources |
| **omos-tak-projector** (exists as MVP) | subscribes WORLD/INTEL/ENV + filters | `tak.cot.inject` → **your ATAK map** | The gate deciding WHAT deserves to be on the tactical map (geofence, type, threat filters — so the map doesn't drown) |

Supporting piece already designed: **omos-control-plane** projects the same data to
atoc (`ui.*`) — so atoc and ATAK show the same world, each at its own zoom.

**The result you asked for:** planes, ships, fires, radiation, conflict events,
chokepoints, military bases — live on YOUR iTAK/ATAK map, filtered so it stays
tactical, powered by streams you already provisioned. That's the massive tactical
map, built on the structures already designed. No new map platform.

---

## 5. WHY WE DON'T NEED THE REST (the cuts, with reasons)

| Item | Why it's no longer needed |
|---|---|
| Running worldmonitor/Crucix/osiris as apps | We took their source lists, adapters, geo layers, prompts. Their UIs duplicate atoc. Keep ONE reference copy each, stop thinking of them as services. |
| `osiris-audit/` | Byte-duplicate of osiris. Zero unique value. |
| Separate map infrastructure now | takserver-nats-bridge + omos-tak-projector already deliver to the map. Sovereign tiles (PMTiles "Earth download") = FUTURE item, one small container when we're ready. 3D globe (worldwideview pattern) = future mission-viz, sits on the SAME streams — nothing to build now. |
| ML platforms now | Vision organ already runs YOLO. tak-ml pattern noted for edge ML later. ML consumes the same bus data — adding it later costs nothing today. |
| NIZAM/Palisade/ATAKRR/xTAK clones | Nothing real to extract that we don't already have live. |
| Mac Docker stacks (go-events, lets-bet, peering, logtrace) | Superseded by bus + ClickHouse + cloud observability + future omos-mission/media. Cut after replacements proven (your rule). |
| Dead bridges, gateway :18789, test murmurs, nex scratch | Retired by NATS-native architecture; already flagged. |

---

## 6. DEPLOYMENT — one step at a time (each proven before the next)

```
STEP 0  NATS leaf link: Mac bus ⇆ cloud hub bus            ← everything depends on this
STEP 1  omos-control-plane → cloud hub (first migration, proves the pattern)
STEP 2  omos-feed-world → cloud hub (first NEW connector; WORLD stream comes alive;
        planes+ships+fires on the bus)
STEP 3  omos-tak-projector v2 → edge (WORLD/INTEL → CoT filters; YOUR MAP LIGHTS UP) ★
STEP 4  omos-feed-intel → cloud hub (INTEL stream alive; alerts/briefing data)
STEP 5  Hermes brief job (osiris/Crucix prompt pattern; daily INTEL brief via agents)
STEP 6  om-unity NATS adapter on (ORP entities → bus → map)
STEP 7  omos-router / omos-mission (workflows: "if X near geofence → alert/mission")
STEP 8  cuts executed one-by-one as replacements are confirmed (your GO each time)
FUTURE  omos-maps (PMTiles Earth), 3D globe viz, edge ML, omos-tak-chat, TAKVR/AR
```

★ Step 3 is the payoff step — that's when your ATAK map stops being only your own
assets and becomes the world picture.

---

## 7. NAMING LAW (recorded)
Everything is **omos-***, not om-*. New services use omos- prefix
(omos-feed-world, omos-tak-projector, omos-mission...). Existing om-* binaries get
renamed as they're touched/migrated — no big-bang rename.
