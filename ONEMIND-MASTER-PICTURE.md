# OneMind Master Picture — Full System Synthesis & Cut-the-Fat Plan
**Date:** 2026-07-04 · **Status:** For Zeus review before migration resumes
**Sources fused:** 6 Desktop architecture docs (6/30→7/4), omos-core tree, live cluster state, external repo research (AMOS, tak-ml, tak-gpt, TAK_OS, constellation-overwatch, ORP), four-plane infra doctrine, Synadia agent SDK.

---

## 1. THE ONE SYSTEM (what OneMind actually is)

```
                       ┌─────────────────────────────────────────┐
                       │   MANAGEMENT: Portainer (one screen)     │
                       └─────────────────────────────────────────┘
  ┌──────────────────────────────┐        ┌───────────────────────────────┐
  │  CLOUD HUB (DOKS K8s) ✅LIVE │        │  EDGE / SOVEREIGN (NEX+local)  │
  │  always-on brain             │  NATS  │  the moat — works offline      │
  │  - NATS HA 3-node (38 strms) │◀─leaf─▶│  - Mac NATS (until Proxmox)    │
  │  - om-control-plane          │  link  │  - om-overwatch kill chain     │
  │  - om-router / om-mission    │        │  - om-unity (ORP graph)        │
  │  - om-command-gateway        │        │  - TAK Server + bridge         │
  │  - Hermes/Legacy gateway     │        │  - omos-voice                  │
  │  - om-feed-* (world intel)   │        │  - drone/robot bridges         │
  │  - observability (3 pillars) │        │  - Home Assistant              │
  │  - Cloudflare tunnel (0-in)  │        │                                │
  └──────────────────────────────┘        └───────────────────────────────┘
              ▲                                        ▲
              │            ONE CONTRACT                │
              │   Envelope v1 + subject taxonomy       │
              ▼                                        ▼
  ┌─────────────────────────────────────────────────────────────────────┐
  │  INTERFACES: atoc (Flutter command console) · TAK (iTAK/ATAK/WebTAK  │
  │  tactical map) · voice (PTT) · future: TAKVR / AR overlay (ARKit)    │
  └─────────────────────────────────────────────────────────────────────┘
```

**The law (unchanged, reaffirmed):** every organ speaks the bus. TAK is a field
interface, not source of truth. Cloud is disposable coordination; the core loop
(detect→fuse→C2→intercept) never depends on it. Mac is cockpit, not server.

---

## 2. THE PATTERN YOU NAMED — "NATS + TAK + world view maps + intel services"

Your intuition is correct and it collapses into **three planes of awareness**, all
already in your architecture docs as streams:

| Awareness plane | Streams | Feeds it | Renders it |
|---|---|---|---|
| **TACTICAL** (meters, now) | COT, GEO, TEL, CMD, MISSION | overwatch organs, drones, TAK devices | TAK (iTAK/ATAK/WebTAK), atoc Overwatch |
| **WORLD** (planet, minutes) | WORLD (capped firehose), INTEL (curated) | om-feed-ais/adsb/weather/aprs + intel connectors | atoc map (MapLibre), future 3D globe |
| **PERSONAL/HOME** (family, presence) | EVT, ENV, ZEUS, VOICE | Home Assistant, sensors, presence, biometrics | atoc Team Control / family views |

The intel repos you named slot in as **sources and pattern references, NOT as
services to run**:

| Repo | What it is | Verdict |
|---|---|---|
| **worldmonitor** (AGPL) | AI news/geopolitical dashboard | EXTRACT: source catalog + freshness model → feeds INTEL stream via om-feed-* connectors. Don't run the app. |
| **Crucix** | 27-source intel terminal, zero cloud | EXTRACT: source list + fetch patterns → om-feed connectors. Don't run. |
| **Osiris** (Next.js) | OSINT recon workspace | PRODUCT-SHAPE reference for atoc intel/analyst pages. Don't run. (osiris-audit = duplicate clone → CUT) |
| **worldwideview** | Plugin-driven CesiumJS 3D globe | KEEP as the 3D globe reference. If atoc ever needs a cinematic world view, its plugin/layer architecture is the model. MapLibre stays the default 2D map. |

**"I want my own maps" — the sovereign maps decision:**
You already run GeoServer (:8600) for WMS layers. The gap is **base map tiles**
(the actual street/terrain imagery), which today come from public tile servers =
not sovereign, dies offline. The fix is cheap and proven:
- **Adopt Protomaps/PMTiles**: one ~100GB file = the entire planet's base map,
  self-hosted, serves atoc (MapLibre reads PMTiles natively), WebTAK, and ATAK
  (via tile endpoint). Runs anywhere — cloud hub now, Proxmox later, edge boxes.
- New workload: **om-maps** (tile server, trivial static serving) + GeoServer
  stays for tactical overlays (imagery, KML, mission layers).
- ATAK devices additionally get offline map packs sideloaded (standard ATAK flow).
This makes the ENTIRE stack's maps sovereign: no Google/Mapbox/OSM tile
dependency at runtime.

---

## 3. EXTERNAL REPO VERDICTS (the pasted research, adjudicated)

| Repo | Verdict | What we take | What we don't |
|---|---|---|---|
| **AMOS** (Apache-2.0) | **EXTRACT patterns** | Mission orchestration shapes for om-mission: typed mission contracts, plugin SDK boundaries, kill-web human-approval gate, atak_adapter CoT patterns | Don't adopt Flask runtime — we're NATS-native; their bus is HTTP. Enterprise features gated anyway. |
| **tak-ml** (Raytheon, GPL-3) | **EXTRACT patterns, NOT core** | Model packaging/metadata schema, KServe-style predict API shape, on-device inference lifecycle → informs om-overwatch-vision + future edge ML organs | NOT "our core add-on": it's Java/Android ATAK-plugin-centric, GPL-3 (viral — keep out of our Apache-style core), TAK-Server-coupled. Our data structure core is the Envelope + ORP graph, not tak-ml. |
| **tak-gpt** (Raytheon, GPL-3) | **EXTRACT pattern** | "LLM as TAK chat contact" — but implement OURS as om-tak-chat: a NATS organ bridging TAK chat ↔ agents.prompt.hermes.zeus.*. Hermes becomes a callsign in iTAK. No new LLM plumbing — reuse Hermes/Legacy. | Don't run their Java plugin in TAK Server; GPL + duplicate LLM path. |
| **TAK_OS** (GoTAK, product) | **PRODUCT-SHAPE benchmark** (already in memory) | Overwatch/COP, Tactical VMS, Feeds, Team Control module shape → already driving atoc IA | Nothing to run. |
| **constellation-overwatch** (community-copyleft) | **EXTRACT patterns only** | Entity ontology (CONSTELLATION_TAK_ONTOLOGY entity types), embedded-NATS microlith pattern, mavlink2constellation bridge patterns (device fingerprinting → entity_id) | License bars commercial use as dependency. Local dashboard clone already flagged; defense/ clone is delete-candidate per house rule once organs confirmed. |
| **ORP (shieldofsteel)** | **ALREADY ADOPTED — it's om-unity.** Deepen, don't re-evaluate. | It IS the world-graph engine: 39 protocol adapters (AIS/ADS-B/MAVLink/CoT/MQTT/Modbus/HL7/NATS), ORP-QL, entity graph. Next: turn on its NATS adapter → publish entities to bus (tel.*/ent.*, WORLD stream) instead of HTTP-only | Don't fork it into a distributed system; single-binary is its strength. Fix pagination issue (#17). |
| **NIZAM / OVERWATCH(mandarwagh9) / Palisade / ATAKRR / xTAK / VisionFrame / WarMatrix / SANKALP / DroneShield** | **SKIP all** | NIZAM: unproven solo scaffold — we already HAVE the Lattice-shape (om-overwatch). OVERWATCH: HUD visual reference only, bookmark it. Rest: README-ware/off-target. | Everything. No clones onto disk. |
| **TAKpilot** (not clonable) | **DESIGN north star** for omos-voice → CoT | Speak → marker/mission flow shape | Can't clone; don't chase the Brave1 marketplace. |

**Synadia agent SDK** (`hive/synadia-agents/agent-sdk`): this is the **contract
layer for the agent fabric** — request/reply, discovery, heartbeats, tool calls.
Hermes already follows the pattern. Doctrine: agents (Hermes, Legacy, sentinels,
specialists) speak Synadia-agent-protocol subjects on the AGENTS account;
NEX places them; K8s hosts the always-on gateway. No new framework needed.

---

## 4. CUT THE FAT (verdicts — nothing deleted without your GO)

### CUT NOW (dead/duplicate, zero extraction value left)
| Item | Why |
|---|---|
| `osiris-audit/` repo | Byte-duplicate clone of osiris/. One reference copy is enough. |
| murmur-test (:64739), murmur-stable (:64740) | Leftover debug containers, flagged since 7/2. |
| Python `om_task.py` remnants | Retired 7/1; Go om-task via NEX is canonical. |
| task-bridge, session-bridge, GatewaySyncNotifier | Docs: "dead / not rebuilt / retire." REST-bridge era is over. |
| Gateway :18789 | Superseded by NATS-native frontend (Path A) + cloud om-command-gateway. |
| nex/ scratch: om-overwatch-demo, om-bridge-constellation (RETIRED), logs mixed in | Segregate/remove from workload dir. |

### CUT AFTER REPLACEMENT PROVEN (your build→prove→delete rule)
| Item | Replaced by | Gate |
|---|---|---|
| Mac Docker stacks: logtrace-system (grafana/loki/jaeger) | Cloud hub observability (Prometheus/Loki/Tempo/Grafana) ✅ already live | Point NATS/NEX exporters at hub, verify dashboards, then decommission |
| go-events stack (postgres/redis/api/grafana/prometheus) | EVT stream + ClickHouse Time Machine + hub observability | Confirm nothing publishes to it that isn't on the bus |
| lets-bet / omos-missions (postgres/redis) | om-mission (cloud hub, planned) | Extract "mission economics" idea only if om-mission design wants it |
| peering stack (redis/minio/chitra/mudra) | om-media-registry + object stores (media OS bucket) + MediaMTX for video | Media path design first — video does NOT ride NATS payloads |
| defense/ extracted clones (SentinelX, constellation-overwatch, bachmach-pqc, om-overwatch-vision src) | Live overwatch organs | Your standing rule: delete when organs confirmed rock solid |
| GeoServer *as strategic map answer* | om-maps (PMTiles) + GeoServer demoted to tactical overlay server | Keep container; strategy shifts to self-hosted tiles |

### KEEP, EXPLICITLY (so nothing sacred gets swept)
- All 38 streams/KV landing zones (deliberate, awaiting publishers — never prune)
- wlow-core binaries (DAG workflow tooling — candidate engine inside om-router; evaluate when om-router builds)
- hive/ (Synadia agent SDK + hermes-nats-gateway — the agent fabric contract)
- Foreign clones marked reference (worldwideview, worldmonitor, Crucix, osiris — ONE copy each)
- omos-observ repo (ClickHouse-native observability — the Time Machine)

### STRUCTURAL FIXES (omos-core hygiene, from 7/2 gap analysis — still open)
1. `nats 2/` Finder-collision rename → restore proper `nats/` (HIGH)
2. `01 The Fabric (...)/` 87M docs dump at root → docs/, sane name
3. services/observability + services/omos-observ dead copies → confirm superseded by standalone repo, cut
4. fusion-platform/ (60M RU-language) + fluxrig/ → evaluate: nothing in the new picture needs them → likely CUT (pending your read)

---

## 5. WHERE EVERYTHING RUNS (final placement, one line each)

**CLOUD HUB (K8s, always-on, disposable):** NATS HA hub ✅ · observability ✅ ·
om-control-plane · om-router · om-command-gateway · om-mission · om-scheduler ·
om-feed-* (ais/adsb/weather/intel connectors) · Hermes/Legacy gateway + agent
workers · om-tak-chat · om-maps (planet tiles) · om-media-registry · Postgres/
ClickHouse (managed-first) · Argo CD/cert-manager/cloudflared ✅

**EDGE/SOVEREIGN (NEX + Docker local — Mac now, Proxmox later):** full
om-overwatch kill chain · om-unity (ORP) · TAK Server + takserver-nats-bridge ·
GeoServer · Murmur + omos-voice organs · mavlink/ros2 bridges · Home Assistant ·
local NATS leaf ✅(pending link) · om-maps replica (offline tiles)

**INTERFACES:** atoc (Flutter: Overwatch/Video/Feeds/Team Control/Work/Agents/
Systems) · iTAK/ATAK/WebTAK · voice PTT · later: AR overlay (ARKit vs entity
bus — the sovereign "EagleEye layer 1 substitute"), TAKVR same pattern

**FAMILY & FLEET MISSIONS:** om-mission (lifecycle: cmd.mission.* / evt.mission.*)
+ om-scheduler + policy-gated command path (proposal→policy→approval→execute→audit)
+ Team Control roles = DelaCruz family members as first-class entities with
simplified views; fleet (drones/robots/vehicles/tablets) as entities in the same
graph. One mission system serves both: "intercept drone" and "kids' pickup run"
are the same machinery, different policy + view.

---

## 6. WHAT CHANGES IN THE MIGRATION PLAN (the pivot's actual impact)

The pivot ADDS clarity but does NOT change the next physical step:

1. **Step 0 stands:** NATS leaf link Mac→hub. Nothing works split-brain without it.
2. Migration queue gains explicit later entries: om-mission, om-feed-*, om-tak-chat, om-maps.
3. om-unity gets a near-term task: enable its NATS adapter (publish entities → bus).
4. Fat-cutting runs AS deletions are approved, parallel to migration, one item at a time.
5. tak-ml/tak-gpt/AMOS get pattern-extraction issues, not code adoption.

**Proposed order after Step 0:** om-control-plane → om-router → om-command-gateway
→ Hermes gateway → om-feed skeleton (+INTEL wiring) → om-mission → om-maps →
om-tak-chat. Overwatch/TAK/voice never leave the edge.
