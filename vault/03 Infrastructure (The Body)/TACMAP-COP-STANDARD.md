---
title: "TACMAP — The Common Operating Picture Standard"
type: doctrine
status: draft
created: 2026-06-12
updated: 2026-06-12
owner: zeus
category: infrastructure
domain: SO/03
tags: [tacmap, cop, nats, cot, tak, atak, geo, fleet, command-center, doctrine, milsymbol]
related:
  - "[[ONEMIND-COMMAND-APP-SPEC]]"
  - "[[ONEMIND-OS-STACK]]"
  - "[[NATS-INFRASTRUCTURE]]"
  - "[[TAILSCALE-NETWORK]]"
ai_summary: "DRAFT doctrine for TACMAP — the map-centric Common Operating Picture (COP) layer of OneMind Command. Core thesis: GIVE NATS A WORLD MAP. NATS stays the single source of truth; we ADOPT the battle-tested Cursor-on-Target (CoT) protocol as the geo-interop standard and run a tested TAK Server as an optional EDGE BRIDGE (so standard ATAK/iTAK/WinTAK field clients work for free), rather than rebuilding any of it. We build/maintain only two thin things: a ~150-line pytak NATS<->CoT bridge, and the in-app MapLibre renderer. Defines the geo-envelope (position as a first-class bus field), the tel./evt.world./cmd. subject taxonomy, the two-ring model (inner=own fleet live, outer=world context), the use-vs-build line, IA mapping into the locked 19-page Command app, and a 5-phase rollout. D-name locked: TACMAP."
---

# TACMAP — The Common Operating Picture Standard

> **What this is:** the doctrine for turning OneMind Command from a dashboard into an **operational command center** — a single live map where our own fleet (inner ring) and the world (outer ring) share one canvas, and where the operator *acts on the map* (click an asset, draw a mission, E-STOP a robot), not just watches it.
>
> **Prime thesis: GIVE NATS A WORLD MAP.** NATS is already the nervous system. A Common Operating Picture is just the bus with two added ingredients — **position as a first-class field** and **a renderer**. TACMAP is those two ingredients, nothing more.
>
> **Use-don't-rebuild doctrine:** Cursor-on-Target (CoT) and the TAK client army (ATAK/iTAK/WinTAK) are battle-tested over 20+ years and free. We **adopt CoT as the interop standard** and run a **tested TAK Server as an optional edge bridge**. We do **not** rebuild CoT routing, symbology, certificate PKI, or video relay. We build and maintain only what's genuinely ours: a thin bridge and the in-app map.
>
> **Sovereignty rule:** everything Tailscale-bound (`tail717ba0.ts.net`). The TAK Server's public IP stays dark. No COP data touches the public internet.
>
> **Hard rule (inherited from [[ONEMIND-COMMAND-APP-SPEC]]):** No mock data. The inner ring is mostly empty today — assets that don't publish position render an honest "no signal" state. The map lights up as real nodes come online. The outer ring (world events) is real data now, so the map is never blank.

---

## 0. The one-paragraph architecture

A Common Operating Picture is **NATS + position + a map**. Every fleet asset already lives on the bus; we add a `lat/lon` to its telemetry and a subject convention (`tel.{class}.{id}.position`). The OneMind Command app subscribes over WebSocket and renders each asset as a live, military-symbol node on a MapLibre map. A **thin pytak bridge** translates those same NATS messages into CoT and feeds a **tested TAK Server**, so any standard ATAK phone in the field sees the identical picture — and any field client's position flows back onto the bus the same way. NATS is the source of truth; CoT is the lingua franca for the field; the map is the operator's interface. **That's the whole system.**

```
        OUTER RING — THE WORLD (observe)            INNER RING — OUR FLEET (participate)
        ────────────────────────────────           ──────────────────────────────────────
        WorldMonitor feed → evt.world.>            robots · drones · vehicles · agents
        • aircraft (ADS-B), vessels (AIS)          • nodes (House/Barn/Office) · edge devices
        • fire, conflict, radiation, markets       • LIVE position: tel.{class}.{id}.position
        • REAL DATA TODAY                          • COMMAND: cmd.{class}.{id}.{verb}
                        ╲                          ╱
                         ╲                        ╱
                    ┌──────────────────────────────────┐
                    │   ONE MAP — the COP (MapLibre GL) │
                    │   Command Center hero · HUD theme │
                    │   click → command · draw → mission │
                    └──────────────────────────────────┘
                                   │
                        NATS bus  (SOURCE OF TRUTH)
                          │                    │
            ┌─────────────┘                    └──────────────┐
            ▼                                                  ▼
   ClickHouse Time Machine                         pytak NATS↔CoT bridge (~150 LoC)
   (scrub/replay, 10-yr TTL)                                  │
                                                              ▼
                                                   TAK Server  :8089 TLS  (EDGE BRIDGE)
                                                              │
                                                   ATAK · iTAK · WinTAK field clients
```

---

## 1. What ATAK/TAK actually is (so we know exactly what we're adopting)

**TAK = Team Awareness Kit.** A US-government family of map clients (originally AFRL, now the TAK Product Center) built around one protocol. It is the de-facto open standard for tactical situational awareness — ~300k US federal users, deployed in disaster response, SAR, and military ops in 20+ countries.

| Component | What it is | License | Our use |
|---|---|---|---|
| **CoT (Cursor-on-Target)** | The protocol. A tiny, stable XML/protobuf event = *what / when / where*. Unchanged base class since 2002. | MITRE open standard | **ADOPT as our geo standard** |
| **ATAK-CIV** | Android client. Maps, BFT, geofences, routes, GeoChat, CasEvac, range/bearing, 3D, video, MIL-STD-2525 symbology — all core. Android-APK plugin SDK. | US public domain | **Field client (use as-is)** |
| **iTAK** | iOS client. Subset of ATAK, no plugins (iOS sandbox). | Free, closed binary | **Field client (iPhone/iPad)** |
| **WinTAK** | Windows client. Near-ATAK parity, command-post oriented. DLL plugins. | GOTS, free | **Optional command-post client** |
| **TAK Server (TPC)** | Java/Spring + PostGIS. CoT streaming :8089 TLS, Marti REST API :8443, WebTAK, federation, client-cert PKI. | **Apache-2.0** | **EDGE BRIDGE (run in Docker)** |
| **CloudTAK** | Open-source *browser* COP (Vue/TS) over the TAK Marti API. | AGPL-3.0 | **Reference / web fallback** |
| **MediaMTX** | RTSP/WebRTC video relay; the de-facto TAK video server. | MIT | **Adopt for drone/cam video** |
| **pytak / takproto** | Python libs to read/write CoT (the snstac gateway family: adsbcot, aiscot). | Apache-2.0 / MIT | **The bridge is built on these** |

**The key realization:** CoT *is* MIL-STD-2525 symbology expressed as a hyphen-delimited string (`a-f-G-E-V` = friendly-ground-equipment-vehicle). By adopting CoT type codes as our asset vocabulary from day one, we get military symbology for free and instant interop with the entire TAK client army — **zero client-side work to put our fleet on a field operator's phone.**

### The plugin ecosystem we inherit (don't build)
- **Drone C2 + FMV video:** UAS Tool (official), FreeTAKUAS (DJI), MAVLink bridges → CoT
- **Mesh / comms-denied:** Meshtastic ATAK plugin (LoRa), goTenna — CoT over off-grid radio
- **Airspace/maritime:** adsbcot (aircraft), aiscot (vessels) — already the exact pattern our bridge uses
- **Camera-to-map:** OpenTAK ICU streams a phone camera to the COP as a clickable video icon

These exist, are tested, and speak CoT. The moment we run a TAK Server, the fleet can use any of them.

---

## 2. The standard we DO define: the geo-envelope

This is the only protocol work that's ours — and it's an *extension* of the existing OneMind envelope, not a new protocol. It maps 1:1 onto CoT so the bridge is a pure field-rename.

```
# Geo-envelope — extends the base OneMind message envelope
lat            decimal degrees WGS-84        → CoT point/@lat
lon            decimal degrees WGS-84        → CoT point/@lon
hae            height above ellipsoid (m)    → CoT point/@hae
course         heading, deg true             → CoT detail/track/@course   ← "the lines"
speed          ground speed, m/s             → CoT detail/track/@speed     ← "the lines"
ce, le         circular/linear error (m)     → CoT point/@ce, @le
cot_type       MIL-STD-2525 type code        → CoT event/@type
uid            stable global id              → CoT event/@uid  = onemind-{class}-{id}
callsign       display label                 → CoT detail/contact/@callsign
ts             ISO-8601 UTC                  → CoT event/@time, @start
stale_s        freshness window (s)          → CoT event/@stale = ts + stale_s
```

**`course` + `speed` are "the lines."** A position with a vector is a heading arrow and a projected track — the visual primitive Zeus asked for to "use the lines to make decisions." Range-rings, ETAs, and intercepts all derive from these two fields.

### Subject taxonomy (binds into [[NATS-INFRASTRUCTURE]])

```
# Inner ring — our fleet (participate)
tel.{class}.{id}.position      class = robot | drone | vehicle | node | edge | agent
tel.{class}.{id}.status        battery, mode, link health (non-positional)
cmd.{class}.{id}.{verb}        estop | goto | rtl | arm | ...      ← act on the map
cmd.mission.{id}.*             mission lifecycle (waypoints, status, abort)
evt.geofence.{id}.{enter|exit} spatial Watch Job fires

# Outer ring — the world (observe)
evt.world.{domain}             domain = air | sea | fire | conflict | radiation | market

# Field interop (bridge-managed)
tak.cot.field.{uid}            raw CoT arriving FROM ATAK/iTAK/WinTAK clients → bus
```

**The standard, stated once:** *any actor or event on the bus that carries a geo-envelope is renderable on TACMAP, commandable from TACMAP, and visible to TAK field clients — with zero map-code changes when a new asset class appears.* Adding the Go2, the drone, or a new barn sensor to the COP = publishing to `tel.{class}.{id}.position`. Nothing else.

---

## 3. The two rings

### Inner ring — our fleet (we ARE the map)
The participatory layer. Every asset Zeus owns, live, as a military-symbol node:
- **Robots** (`a-f-G-E-V`): DIMOS, OM1, Go2, future UGVs
- **Drones** (`a-f-A-M-F-Q`): aerial, with FMV video click-through
- **Vehicles** (`a-f-G-E-V`): the 2 cars, motorcycle, tractor, the manned plane (as `a-f-A-C-F`)
- **Nodes** (`a-f-G-I-U-T-C`): House, Barn, Office — fixed installations
- **Edge devices** (`a-f-G-U-C`): Pis, sensors, cameras
- **AI agents** (`a-f-G-U-C-I` + lowercase ext): Legacy and the profile fleet, positioned at the node they run on
- **Family devices** (private group, restricted): opt-in BFT for Cielo + kids + Cerberus's collar

Most of this is roadmap → renders as honest "no signal yet" until the asset publishes. This is the participatory COP: missions run here, commands fire here.

### Outer ring — the world (we OBSERVE)
The context layer, fed by **WorldMonitor** (the repo we own, deepest analytics) republished to `evt.world.>`:
- Aircraft (ADS-B), vessels (AIS), wildfire, armed conflict, radiation, markets-as-geo
- **Real data, available now** — so TACMAP demos immediately while the inner ring fills in

### Why fuse them — cross-ring correlation
The payoff of one canvas: **"world event within Nkm of my asset" → automatic alert.** A conflict event near the plane's flight path, a fire near the barn, an unknown vessel near a coastal node. The outer ring *informs the inner ring's decisions*. That correlation — scoped to *your* assets — is the whole reason to fuse observer + participant instead of running them as separate apps.

---

## 4. Use vs. build — the maintenance line

The explicit point of this doctrine: **minimize what we maintain.** Everything below the line is adopted, tested, and someone else's upkeep.

| Capability | Decision | Why |
|---|---|---|
| CoT protocol | **USE** | 20-yr MITRE standard, never rewrite |
| MIL-STD-2525 symbology | **USE** (`milsymbol` JS / CoT type codes) | symbology is solved |
| TAK Server (routing, PKI, PostGIS, federation, WebTAK) | **USE** (Docker, Apache-2.0) | government-grade, free |
| ATAK / iTAK / WinTAK field clients | **USE** as-is | 300k-user tested, zero client work |
| Drone video relay | **USE** MediaMTX (MIT) | don't build an RTSP server |
| Drone/mesh/ADS-B/AIS → CoT | **USE** existing plugins/gateways | adsbcot, aiscot, UAS Tool, Meshtastic |
| Web COP reference | **REFERENCE** CloudTAK (AGPL-3.0) | study, maybe embed for web target |
| ── maintenance line ───────────── | | |
| **NATS↔CoT bridge** | **BUILD** (~150 LoC pytak) | the only glue; trivial, ours |
| **In-app MapLibre renderer** | **BUILD** | the OneMind face; HUD-themed, our product |
| **Geo-envelope + subject taxonomy** | **DEFINE** (this doc) | the standard that makes it all click |
| **Cross-ring correlation + Watch Jobs** | **BUILD** (NATS-native, L4) | our differentiator |

**Net new code we own: a thin bridge + a map widget + a doctrine.** Everything hard — the protocol, the clients, the server, the symbology, the video — is adopted.

> ⚠️ **Anti-pattern, explicitly rejected:** running the TAK Server as the *primary backend* and treating the bus as secondary. That creates a second nervous system competing with NATS — the exact parallel-backend mistake [[ONEMIND-COMMAND-APP-SPEC]] killed (Mission Control, SwarmClaw, WWV-as-app). **NATS is the source of truth. TAK is a bridge that NATS feeds.** All CoT servers speak the same protocol, so swapping OpenTAKServer → official TAK later is a server swap + cert regen, never a fleet rebuild.

---

## 5. How TACMAP lands in OneMind Command (no oversaturation)

Maps cleanly onto the **locked 19-page / 6-group IA** — TACMAP is a *function of the COMMAND group*, plus tabs, not a sprawl of new pages. Consistent with the bus-relationship model:

| Element | Bus relationship | Where in the app |
|---|---|---|
| **The COP map itself** | (renderer over FEED) | **COMMAND → Command Center** becomes map-first; the map IS the hero/home view |
| Fleet asset detail (click a node) | FEED | existing **FLEET** pages (robotics/drones/edge) = the detail views |
| World feed (outer ring) | FEED | a layer toggle on the map; sourced from WorldMonitor → `evt.world.>` |
| **NATS↔CoT bridge** | **BRIDGE** | a tab under **BUS → Connect** (same shelf as saas-worker, n8n, LiveKit, ha-bridge) |
| Mission run / waypoints | FEED + cmd | map draw tool → **OPS** (Task Engine) |
| Time Machine replay | PERSIST | scrub control on the map; ClickHouse `evt.>` |
| Voice-dispatched mission | BRIDGE (P5) | LiveKit intent → map action |

The map is the **face**, the Fleet pages are its **organs**, the bridge is a **Connect tab**. The three "observer" apps (WorldMonitor / WorldWideView / Crucix) collapse into **FEED sources**, not pages — WorldMonitor becomes the world publisher, Crucix's delta-engine becomes NATS-native Watch Jobs, WWV retires. One feed, not three apps.

---

## 6. Three capabilities almost no one else has

1. **Time Machine replay.** ClickHouse already records every `evt.>` with timestamps and a 10-yr TTL (see [[ONEMIND-COMMAND-APP-SPEC]]). TACMAP can **scrub the map backward in time** — replay where every asset and world event was at any past moment. Real command centers pay millions for this; we already have the substrate.
2. **Voice-dispatched missions.** LiveKit (P5) + map + Task Engine: *"Legacy, send drone 1 to the north field"* → route renders → confirm → `cmd.mission.*` fires.
3. **Cross-ring correlation.** Automatic proximity alerts between the world and your assets — the fusion payoff.

---

## 7. Rollout (phased by data-readiness, mirrors the app's P0→P5 discipline)

| Phase | Deliverable | Depends on |
|---|---|---|
| **T0 — Standard** | This doctrine + geo-envelope + subject taxonomy ratified into [[NATS-INFRASTRUCTURE]] | nothing |
| **T1 — Map shell** | MapLibre COP in Command Center, HUD theme, outer ring live (WorldMonitor → `evt.world.>`). Real data, never blank. | WorldMonitor NATS sink |
| **T2 — Inner ring** | Fleet assets publish `tel.{class}.{id}.position`; nodes/edge first (they have Tailscale positions today), robots/drones as they arrive | per-asset telemetry |
| **T3 — Bridge** | pytak NATS↔CoT bridge + TAK Server in Docker (Tailscale-bound); ATAK phone sees the fleet; field CoT flows back to the bus | T2 + TAK Server deploy |
| **T4 — Act on the map** | Click-to-command (`cmd.*`), draw-a-mission → Task Engine, geofence Watch Jobs, E-STOP | Task Engine (L4) |
| **T5 — Fusion** | Cross-ring correlation alerts, Time Machine scrub, voice-dispatched missions | T1+T2+ClickHouse+LiveKit |

**T1 is prototypable now and demos immediately.** T3 is the "give a field phone the fleet" milestone. T4 is where it becomes a true command center.

---

## 8. Open decisions (for Zeus)

1. **TAK Server choice for T3:** official TAK Server (TPC, Apache-2.0, full Marti API + federation, ~4Gi RAM) vs **OpenTAKServer** (MIT, lighter, PostGIS, Marti subset, edge-friendly). Recommendation: **OpenTAKServer first** (lighter, fully open, no tak.gov account), with a documented swap path to official TAK — all speak CoT, so it's reversible.
2. **MIL-STD-2525 in Flutter-native:** `milsymbol` is JS. Options: (a) webview the renderer, (b) port icon generation to Dart, (c) pre-render a sprite atlas. The interactive artifact uses `milsymbol` directly to de-risk the *look*; the Flutter path is a T1 spike.
3. **Host:** omos-cloud (always-on, 12Gi free, x86) is the natural T3 home for the bridge + TAK Server.
4. **Family BFT group:** opt-in, isolated TAK group/account, restricted to the inner-ring private layer. Confirm before any family device publishes position.

---

## 9. Source index (verified 2026-06-12)

| Reference | URL |
|---|---|
| TAK Server (official, Apache-2.0) | https://github.com/TAK-Product-Center/Server |
| ATAK-CIV source | https://github.com/TAK-Product-Center/atak-civ |
| pytak (NATS↔CoT bridge core) | https://github.com/snstac/pytak |
| takproto (CoT protobuf) | https://github.com/snstac/takproto |
| adsbcot / aiscot (the gateway pattern we copy) | https://github.com/snstac/adsbcot · https://github.com/snstac/aiscot |
| CloudTAK (open-source web COP) | https://github.com/dfpc-coe/CloudTAK |
| node-cot / node-tak (TS CoT libs) | https://github.com/dfpc-coe/node-CoT · https://github.com/dfpc-coe/node-tak |
| OpenTAKServer (light self-host) | https://docs.opentakserver.io |
| MediaMTX (video relay) | https://github.com/bluenviron/mediamtx |
| Meshtastic ATAK plugin (off-grid CoT) | https://github.com/meshtastic/ATAK-Plugin |
| MITRE CoT Developer's Guide | https://nps.edu/documents/104517539/109705106/COT+Developer+Guide.pdf |
| CoT type codes / MIL-STD-2525 | https://freetakteam.github.io/FreeTAKServer-User-Docs/About/architecture/mil_std_2525 |
| milsymbol (2525 in JS) | https://github.com/spatialillusions/milsymbol |
| MapLibre GL JS | https://maplibre.org |

---

*Doctrine status: DRAFT — pending Zeus review. The interactive TACMAP artifact (MapLibre + milsymbol + simulated NATS geo-feed, HUD theme) accompanies this doc as the visual proof-of-concept for §1, §3, and §6.*
