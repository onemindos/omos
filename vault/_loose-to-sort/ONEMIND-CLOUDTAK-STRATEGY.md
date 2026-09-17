# OneMind × CloudTAK — Master Strategy & Handoff Document

**Author:** Legacy (Letta) · **Date:** 2026-07-07 · **Audience:** external coding agent (no prior context)
**Status:** Strategy locked. This is the source-of-truth briefing for building it.

---

## 0. TL;DR (read this first)

Zeus is building **OneMind (OMOS)** — a sovereign, self-hosted, NATS-native unified infrastructure for all communication and command/control: humans, drones, robots, sensors, home automation, feeds, missions. Think **open-source Palantir AIP + Anduril Lattice + Maven OS**, but sovereign and running on his own infra.

The problem: OneMind had **no UI**. The half-built option (`omos-atoc`, a Flutter app) was re-implementing what TAK already does.

**The decision (locked):**
- **CloudTAK becomes the primary/main UI** — the platform everything lives inside. It's a mature, open-source, browser-based TAK Common Operating Picture (COP) already deployed live on his cluster.
- **`omos-atoc` is PARKED** (not deleted, not developed further).
- **The fork is THIN and config-first** — extend via CloudTAK's plugin + ETL systems, not by rewriting its core.
- **The clean architectural boundary is TAK Server.** CloudTAK never touches the NATS bus directly. Data flows: `NATS → NatsBridgePlugin → TAK Server 5.7 → CloudTAK`.

The rest of this document explains the full architecture, how CloudTAK's AWS-Lambda-shaped Task/ETL system works, exactly where NATS and data sources connect, the AWS-shape debt and how to run it on Kubernetes instead, the extension ecosystem, and the build roadmap.

---

## 1. The OneMind Architecture — Four Planes

OneMind is **one horizontal fabric**, not a stack of vertical silos. Every source publishes once; every consumer reads the same truth.

```
                        ┌─────────────────────────────────────────┐
   DATA SOURCES ───────▶│              NATS  (THE FABRIC)           │◀────── CONSUMERS
   (AIS, ADS-B,         │   the nervous system — HA, TLS, JetStream │        (CloudTAK,
    drones, sensors,    └───────────────────┬─────────────────────┘         Legacy AI,
    home automation,                        │                                ClickHouse,
    OSINT feeds…)                           ▼                                ATAK field)
                        ┌─────────────────────────────────────────┐
                        │           ORP  (THE BRAIN)                │
                        │  fusion + entity graph + ORP-QL query lang│
                        │  subscribes raw, republishes fused ent.*  │
                        └───────────────────┬─────────────────────┘
                                            │
                        ┌───────────────────▼─────────────────────┐
                        │   TAK Server 5.7  (OPERATIONAL BACKBONE)  │
                        │  CoT streaming, missions, video, federation│
                        │  the clean boundary CloudTAK talks to      │
                        └───────────────────┬─────────────────────┘
                                            │
                        ┌───────────────────▼─────────────────────┐
                        │   CloudTAK  (THE FACE — main UI)          │
                        │  browser COP, map, missions, feeds, video │
                        └───────────────────────────────────────────┘
```

**The four planes and their one job each:**

| Plane | Component | Job | Status |
|---|---|---|---|
| **Fabric** | **NATS** | The bus. Nothing talks app-to-app; everything talks bus. HA nervous system. | LIVE (DOKS 3-node cluster + Mac leaf) |
| **Brain** | **ORP** (Rust `om-unity`) | Fusion: ingest all streams, resolve entities, republish fused `ent.*`. Query via ORP-QL. | LIVE (5,938 AIS ships now) |
| **Backbone** | **TAK Server 5.7** | Operational C2 protocol: CoT, missions, video, federation, cert auth. | LIVE (DOKS, 8 pods) |
| **Face** | **CloudTAK** | The UI everyone uses. Browser COP + ETL connectors + missions + video. | LIVE (DOKS, 6 pods) |

**Critical principle — no single point of failure:** Everything publishes to NATS. ORP subscribes/fuses/republishes. Consumers read from NATS. **If ORP dies, TAK still gets raw CoT and the map still shows raw data — you lose *fusion*, not the *picture*.** ORP is the mind; NATS is the nervous system. That's the difference between "one mind" and "one point of failure."

---

## 2. NATS — The Fabric: Why It's the Superpower

NATS is the reason OneMind can be "one mind." Understand its power before wiring anything.

### 2.1 What NATS is
NATS is a **high-performance messaging system** (a "bus"). Publishers send messages to **subjects** (like `sensor.ais.position`); subscribers listen to subjects (with wildcards like `sensor.ais.*` or `sensor.>`). Publishers and subscribers never know about each other — they only know the subject. This is **publish/subscribe decoupling**, and it's what lets you add a new data source or a new UI without touching anything else.

### 2.2 Why it's fast (the speed that matters here)
- **Millions of messages/sec, sub-millisecond latency** on a single server. It's written in Go, tiny binary, no heavyweight broker.
- **Fire-and-forget core** (NATS core) for real-time telemetry where the newest value matters most.
- **JetStream** (persistence layer) when you need durability, replay, at-least-once delivery, or key/value storage. OneMind uses JetStream streams for TEL/DET/COT/GEO/INTEL/EVT and KV buckets for `ui.*` projections.
- **Leaf nodes** extend the bus across machines/networks with one config line — the Mac runs a full NATS server that leaf-connects to the DOKS cluster over Tailscale. Edge sites, robots, and field devices attach the same way. **This is the sovereignty moat:** the bus spans everywhere, works offline, and syncs when connected.
- **WebSocket + MQTT built in:** browsers and IoT devices can speak NATS directly.

### 2.3 The subject taxonomy — the CONTRACT that combines everything
This is the most important thing to lock. What actually unifies OneMind is not shared code — it's a **shared subject language**. Every source and consumer codes to this:

```
sensor.*   →  raw feeds        (sensor.ais.position, sensor.adsb.track)
ent.*      →  fused entities   (ORP output — one entity per real-world thing)
intel.*    →  curated intel    (analyst/agent-enriched)
world.*    →  geo/static       (basemaps, cables, static overlays)
cot.*      →  TAK Cursor-on-Target (the bridge to TAK Server)
ui.*       →  UI projections   (cached last-good state for frontends)
mission.*  →  tasking/orders   (missions assigned to entities)
```

Lock this taxonomy and **every new integration becomes copy-paste one line:** `nc.publish(subject, data)`. The hard part (fetch, parse, auth) is already written in existing source repos; the only reason a source isn't "on the bus" is that its last step prints to console instead of publishing to NATS.

---

## 3. Why CloudTAK (and why atoc is parked)

### 3.1 The evidence
- **atoc** = 116 Dart files (~1.4M), an early shell being actively built. Rebuilding TAK's capabilities in Flutter = 6–12 months.
- **CloudTAK** = 60+ API route files across ~20 route groups, **live on DOKS now**, AGPL-3.0, actively developed by a real consortium (Colorado DFPC Center of Excellence + community), already patched to talk to TAK Server 5.7.

The choice was never "my thing vs their thing." It was "finish someone else's mature C2 UI" vs "rebuild it from scratch." CloudTAK wins.

### 3.2 CloudTAK's route groups = the exact C2 surface OneMind needs
| Group | What it does | OneMind use |
|---|---|---|
| **Connection / Layer** | ETL pipelines — pull any data source → CoT → map. The plugin engine. | Where feeds plug in |
| **Marti** | Direct TAK Server passthrough — contacts, missions, packages, video | Missions/comms |
| **Basemap** | Map tile management — ZXY, PMTiles, Esri | Sovereign tiles (GeoServer) |
| **Server / Injector** | Programmatic CoT injection | NatsBridgePlugin already does this |
| **Task** | Serverless Lambda-style compute (AWS-oriented, pluggable) | ⚠️ AWS-shape — see §5 |
| **Video / Lease** | MediaMTX video stream management | Drone/CCTV video |
| **Profile** | Per-user overlays, saved features, assets, chatrooms | Family/roles/permissions |
| **Search** | Geocode, reverse-geocode, elevation, weather, declination | Field tools |
| **Config** | Server-wide settings, **whitelabel**, display options | "Dress it however we like" |

"Give missions and permissions to any entity, ingest outside info, control home automation" = **exactly this shape, already built.**

### 3.3 atoc verdict: PARK
- **Park, don't delete, don't keep building.** It's Zeus's code, in git — nothing lost.
- Reasons: (1) a second UI breaks unification; (2) kids/family aren't on yet, so atoc's only differentiator (family face) isn't needed; (3) field/mobile is already covered — CloudTAK has a Capacitor iOS build + ATAK/iTAK is the battle-tested field client.
- **Revisit only if** family goes live AND CloudTAK's Profile + whitelabel can't dress down far enough for a child. Likely they can → atoc never returns. One platform dressed per-role > two platforms.

---

## 4. Field Scan — Is CloudTAK One-of-a-Kind? (Yes, effectively)

Before committing, the landscape was searched. **The field confirms CloudTAK rather than threatening it.** When you require all four of: (1) TAK-native, (2) full browser COP, (3) real connector/ETL framework, (4) actively maintained by more than a solo dev — the field collapses to exactly one project.

**TAK-native COP competitors (all narrower):**
- **tak-webview-cesium** (sgofferj) — lightweight Cesium 3D CoT *viewer*. No ETL/missions/profiles. A lens, not a platform.
- **xTAK** (adamstern2000) — hardware gateways (LoRa/SDR/radio→TAK). Feeders, not a UI.
- **GroundWave** (MIT) — offline SA platform, own protocol, CoT bridge "coming soon." Not TAK-native.

**General C2 platforms (different protocol universe — disqualified):**
- **ODIN / ODINv2** (syncpoint, AGPL) — most mature open C2IS, but Matrix-protocol + desktop-first, not TAK, not NATS. Adopting = throwing away the TAK 5.7 + NATS investment.
- **Tower / Cockpit / Maven-System / HYDRA-C2 / NIZAM / Sovereign_Watch / CAIRN** — robotics fleet control, MAVLink GCS, or solo/demo-stage Lattice-clone experiments. None production-grade AND TAK-native AND plugin-driven.

**Conclusion:** Commit to CloudTAK. There is no better option for a sovereign, TAK-native, plugin-driven browser COP.

---

## 5. CloudTAK's Task/ETL System — The AWS-Lambda Architecture (DETAILED)

This is the part Zeus specifically wants documented in full, because it's both the power and the debt.

### 5.1 The three concepts: Connection → Layer → Task

CloudTAK's data-ingestion model has three nested pieces:

1. **Connection** — the core container. At its heart it's a **single private/public key certificate pair** that authenticates to TAK Server. Every integration pushes/pulls data through a Connection. (Think: "this is how I'm allowed to talk to TAK.")
2. **Layer** — a data pipeline attached to a Connection. A Layer pulls from a source, transforms to CoT (Cursor-on-Target), and submits it. Layers have direction (`Incoming` = data → TAK, `Outgoing` = TAK → data) and an invocation type (`Schedule` = cron-like, or event-driven).
3. **Task** — the actual runnable code unit for a Layer. **This is the "serverless Lambda-style compute."** Each ETL Task is a Node.js/TypeScript program packaged as a **Docker container image**, designed to run as an **AWS Lambda**.

### 5.2 How an ETL Task actually runs (grounded in the code)

From reading the CloudTAK source (`api/lib/aws/`), here's the real mechanism:

- **Task code** extends the `@tak-ps/etl` base class (`etl-base`). It implements a `control()` method that:
  1. Calls `this.fetchLayer()` to get its config (env vars set by the user in the CloudTAK UI).
  2. Fetches from its data source (HTTP API, etc.).
  3. Converts results to **GeoJSON FeatureCollection** (with `callsign` etc. in properties).
  4. Calls `this.submit(featureCollection)` — CloudTAK converts GeoJSON → CoT and injects into TAK Server.
- **Deployment of a Layer** (`api/lib/aws/layer-deploy.ts`): when a user adds/configures a Layer in the UI, CloudTAK **creates a dedicated AWS CloudFormation stack per layer** (`CloudFormation.create(config, layerid, stack)`). That stack provisions a **Lambda function** (from a Docker image in **ECR**) plus its schedule/trigger. If the Lambda needs replacing, it deletes and recreates the stack.
- **Queuing**: `api/lib/aws/queue.ts` uses **AWS SQS** (`SendMessageBatchCommand`) to fan hook/webhook messages.
- **Storage**: `api/lib/aws/s3.ts` for assets (swappable — see below).
- **Logs/alarms**: `lambda-logs.ts` (CloudWatch Logs) and `alarm.ts` (CloudWatch Alarms).
- **Scheduling**: scheduled Layers fire on an interval (EventBridge/CloudWatch schedule in AWS), invoking the Lambda, which runs `control()` once and exits. This is the "Lambda-style" nature — **short-lived, invoked on a schedule, stateless per run.**

**The AWS services the dynamic ETL system depends on:** CloudFormation (per-layer stacks), Lambda (the task runtime), ECR (task images), SQS (queue), S3 (assets), CloudWatch (logs/alarms), EventBridge (schedules).

### 5.3 The long-running tasks (different from ETL Lambdas)

Separately, CloudTAK has **always-on worker tasks** that run as containers, NOT Lambdas:
- **`events`** — event/webhook processing worker (uses S3 via `@aws-sdk/client-s3`).
- **`pmtiles`** — the tile server (serves PMTiles basemaps).
- **`retention`** — data retention/cleanup, runs on a cron (`RETENTION_SCHEDULE`, default `0 3 * * *`).

In `docker-compose.yml` these are plain services. In AWS they'd be ECS services. **These already run fine on Kubernetes** — they're just containers.

### 5.4 What's abstracted vs hard-coupled

| Piece | AWS service | Portable? |
|---|---|---|
| **api** (main server) | ECS container | ✅ Plain container — runs on K8s |
| **events / pmtiles / retention** | ECS containers | ✅ Plain containers — run on K8s |
| **S3 (assets)** | S3 | ✅ Already swapped for **MinIO** (CloudTAK deployed on DOKS with a `minio` pod; `AWS_S3_Endpoint` points at it) |
| **PostGIS (database)** | RDS | ✅ Runs as a pod |
| **Dynamic ETL Layer deploy** | CloudFormation + Lambda + ECR + SQS + EventBridge | ❌ **This is the hard AWS coupling** |

**Bottom line:** CloudTAK's *core* (UI, map, missions, video, the always-on workers, S3-via-MinIO) is already running sovereign on DOKS with zero AWS. The **only** AWS-locked feature is the *in-UI dynamic ETL Layer auto-deploy* — the ability to click "add a connector" in the web UI and have it spin up a Lambda automatically.

---

## 6. THE KEY QUESTION — Where NATS + Data Sources Connect to CloudTAK

Zeus's core question: *how do NATS and data sources connect / where do they fit into the OneMind fabric?* There are **three integration paths**. Choosing the right one per source is the whole game.

### Path A — Through TAK Server (the PRIMARY, clean-boundary path) ⭐

```
data source → NATS (sensor.*) → ORP (fuses → ent.*) → NatsBridgePlugin → TAK Server (CoT) → CloudTAK
```

- **NatsBridgePlugin** is Zeus's own TAK Server plugin. It subscribes to NATS `cot.*` / `tak.cot.inject` and injects CoT into TAK Server. TAK Server then broadcasts to **every connected client at once** — CloudTAK, ATAK/iTAK field devices, WinTAK, WebTAK.
- **cot-tap** handles the reverse (TAK → bus): a lightweight mTLS client on TAK's port 8089 that reads the CoT firehose and republishes to `tak.cot.out` on NATS.
- **Why this is primary:** CloudTAK never touches the NATS bus. It only speaks to TAK Server (which Zeus fully controls). **One path, one truth.** No risk of two data paths disagreeing. This is the boundary that keeps the fork thin.
- **Use for:** anything that should appear on the tactical map AND on field devices — entities, tracks, positions, missions. The vast majority of live operational data.

### Path B — CloudTAK ETL Layer (the connector path, for non-TAK sources)

```
external API → CloudTAK ETL Task (Layer) → CoT → TAK Server → all clients
```

- Use CloudTAK's built-in Connection/Layer/Task system to pull a data source directly.
- **For OneMind, invert it:** instead of an ETL Task that does `fetch(externalAPI)`, write one that does `nats.subscribe(subject)`. It's the same `etl-base` `control()` pattern — swap the data acquisition line. The 10+ existing dfpc-coe ETL connectors (see §8) are working templates.
- **Use for:** sources where you want CloudTAK to own the pull/schedule, or to reuse an existing dfpc-coe connector (ADS-B, GeoJSON, ArcGIS, CCTV) with minimal change.

### Path C — CloudTAK Plugin (UI/behavior, e.g. a NATS dashboard panel)

```
NATS (WebSocket) → CloudTAK Plugin (Vue) → sidebar panel / overlay / status widget
```

- CloudTAK's **Plugin API** lets you add custom menus, routes, map overlays, and status-bar widgets to the web client. NATS has a WebSocket port, so a plugin *could* subscribe to NATS directly in-browser for a dashboard/telemetry panel.
- **Use sparingly** — only for UI surfaces that aren't map entities (e.g., a "bus health" panel, stream monitors, an AI chat panel). Do NOT use this to render operational entities; that belongs on Path A so field devices see them too.

### Decision rule
- **Is it a map entity / operational track?** → **Path A** (NATS → TAK → CloudTAK). Field devices need it too.
- **Is it a non-TAK data feed you want CloudTAK to schedule/own?** → **Path B** (ETL Layer, NATS-subscribe variant).
- **Is it a UI panel / AI chat / dashboard, not a map entity?** → **Path C** (Plugin).

**For the OneMind super-fabric, Path A is the backbone.** ORP fuses everything into `ent.*`, the bridge pushes to TAK, and TAK fans out to CloudTAK + every field device simultaneously. That's the "one mind, many lenses" realized.

---

## 7. The AWS-Shape Debt — Concrete, and How to Run It on Kubernetes

This is the one real cost of adopting CloudTAK. Documented plainly so the other agent can plan for it.

### 7.1 What breaks on DOKS
The **dynamic in-UI ETL Layer deploy** assumes AWS: it creates a CloudFormation stack → Lambda → ECR image → EventBridge schedule → SQS. **None of that exists on DigitalOcean Kubernetes (DOKS).** So the "click to add a connector in the CloudTAK admin UI" flow will not function out of the box.

### 7.2 What does NOT break
Everything else. The api, the map, missions, video, Profile, Marti, Basemap, Search, Config, and the always-on workers (events/pmtiles/retention) all run as plain containers. S3 is already MinIO. This is proven — CloudTAK is live on DOKS right now with 6 pods and zero AWS.

### 7.3 The fix — run ETL connectors as native Kubernetes workloads
Instead of CloudTAK auto-deploying Lambdas, deploy each connector yourself as a **K8s CronJob** (for scheduled pulls) or **Deployment** (for streaming subscribers):

```
AWS model:   UI click → CloudFormation → Lambda (scheduled) → submit CoT to TAK
K8s model:   git manifest → CronJob/Deployment → same task code → submit CoT to TAK
```

- The **connector code ports unchanged** — it's the same `etl-base` `control()` that fetches (or `nats.subscribe`s) and submits GeoJSON. Only the *deploy wrapper* changes: a K8s CronJob YAML instead of a CloudFormation/Lambda stack.
- The task still needs `ETL_API` (points at CloudTAK's api service) and `ETL_LAYER` (the layer id) env vars — set them in the manifest instead of having Lambda inject them.
- This fits OneMind's **GitOps doctrine perfectly:** every connector becomes a version-controlled manifest in `omos-infra`, deployed via ArgoCD. More durable than click-to-deploy anyway.
- **Optional future work:** build a shim that maps CloudTAK's "deploy layer" API call to a Kubernetes Job creation instead of CloudFormation, restoring the in-UI experience. Not needed for v1 — deploy connectors as manifests.

### 7.4 Even simpler for OneMind
Because **Path A** (NATS → TAK → CloudTAK) is the backbone, most data never needs a CloudTAK ETL Layer at all. ORP + NatsBridgePlugin already put fused entities on TAK. CloudTAK just *displays* them. **The AWS-shaped ETL system is largely bypassed** — you only reach for Path B when you want CloudTAK to own a specific non-TAK pull. This significantly shrinks the debt in practice.

---

## 8. The Extension Ecosystem — What's Available to Build On

CloudTAK has **two distinct extension surfaces**. Don't conflate them.

### Surface 1 — Plugin API (UI/behavior)
Official, documented: *"extend the web client with custom menus, routes, map overlays, and status-bar widgets."* Vue/TypeScript, dropped in `./api/web/plugins/`, then `npm run build`.
- **dfpc-coe/CloudTAK-Plugin-Sample** (MIT) — official starter template from the maintainer (`ingalls`).
- **snstac/CloudTAK-Plugin-HelloWorld** (MIT) — working example by **Greg Albrecht (ampledata / Sensors & Signals LLC)**, the biggest independent name in open TAK tooling. snstac maintains the `*cot` gateway family: **aiscot** (AIS→CoT), **adsbcot** (ADS-B→CoT), **aprscot** (APRS→CoT), and more. Second ecosystem anchor besides dfpc-coe.
- **TAKpilot** (Corvus Intelligence, AGPL-3.0) — the AI copilot (see §9). No competing AI plugin exists.

### Surface 2 — ETL / Connector library (data in)
A whole library, all MIT, all from dfpc-coe, all `source → CoT → map`:
- **etl-adsbx** — ADS-B Exchange → air picture (directly relevant to air feeds)
- **etl-geojson** — any HTTP GeoJSON API (the generic workhorse)
- **etl-arcgis** — ESRI Feature/MapServer
- **etl-caltopo**, **etl-cotrip-weather**, **etl-cotrip-cameras** (CCTV), **etl-geotab**, **etl-optimus**, **etl-skymira** (fleet/GPS)
- **etl-template** + **etl-base** — scaffolding to write your own in ~one file
- **node-cot** + **node-tak** — the SDKs under everything (already patched for TAK 5.7)
- **TAK-NZ** maintains parallel etl-* forks with cleaner CI patterns worth studying

**For OneMind:** clone `etl-template`, replace `fetch(externalAPI)` with `nats.subscribe(subject)` → convert to GeoJSON → `submit()`. That's a NATS→CloudTAK connector. Deploy it as a K8s CronJob/Deployment (§7.3).

---

## 9. The AI Layer — TAKpilot on Bedrock

Zeus wants AI woven into the map (his agent is "Legacy"). **TAKpilot** is the anchor and validates the whole architecture.

- **What it is:** an AI chat copilot that embeds *inside* CloudTAK. Operator types or dictates natural language ("place a hostile infantry contact at 37U DP 12345 67890, callsign CONTACT-7") → TAKpilot translates it into a structured TAK API call → executes → marker on the map in under 2 seconds. No menus.
- **License:** AGPL-3.0. Full source available (Node.js service + CloudTAK UI extension + tool library + deployment docs).
- **Model-agnostic:** OpenAI-compatible endpoint abstraction. Runs on **AWS Bedrock** (Zeus's confirmed model path), **Ollama/vLLM** (sovereign/air-gapped), or Anthropic — no code changes. This enables classified/air-gapped deployments with the same features.
- **Integration:** runs as a Node.js service on the same host/LAN as CloudTAK, talks to CloudTAK's existing REST API using operator session creds — **no new TAK ports, no federation changes, no DB schema changes.** UI extended with a chat panel over WebSocket.
- **Why it matters:** it's the exact "Legacy-into-the-map" AI layer, already written, on Bedrock-compatible rails, using the same thin-fork/TAK-as-boundary pattern. An independent defense shop converged on the same architecture — strong validation. **Fork it, point it at Bedrock, extend its tool library.**

---

## 10. Fork Strategy — Thin, Config-First (the discipline)

| Approach | Verdict | Why |
|---|---|---|
| **Thin fork, config-first** ✅ | **CHOSEN** | Feeds via ETL/Layer + Path A (config, zero core changes). Fork frontend only for branding + optional NATS panel. Defer custom API route groups. Small diff = easy upstream rebase. Respects the NATS structure. |
| Full fork, go all-in | ❌ Rejected | Owning the whole fork = merge conflict on every upstream update (CloudTAK is actively developed). Integrating NATS *inside* CloudTAK's API would bypass TAK Server and create two disagreeing data paths — breaks the layered architecture. |

**Rules:**
1. **Don't fork the API yet.** Push feeds through Path A (NATS→TAK→CloudTAK) and Path B (ETL Layers). Both are config, not core edits.
2. **Fork the frontend only for branding** (CloudTAK has a central rename/whitelabel via the Config group) + later one NATS/AI sidebar panel. Keep the diff tiny. `git remote add upstream`, rebase regularly.
3. **Defer custom API route groups** (`/api/nats/*`, `/api/fleet`) until proven necessary. Most "fleet management" = Profile + Marti contacts + CoT already.
4. **The boundary is TAK Server** — CloudTAK never touches the bus. Non-negotiable.

**License note:** CloudTAK is **AGPL-3.0**. Fine for sovereign self-host. Copyleft only triggers if you offer *modified CloudTAK as a hosted service to others* — running it for yourself does not. TAKpilot is also AGPL-3.0 (same deal). The plugin *sample* is MIT.

---

## 11. Current Live State (as of 2026-07-07)

Everything below is deployed and running on **DOKS (DigitalOcean Kubernetes, `onemind-hub`, nyc3)** unless noted:

- **NATS** — 3-node HA cluster, R3 JetStream, TLS 1.3, account auth. Mac runs a full NATS server leaf-connected via Tailscale (`omos-nats.tail717ba0.ts.net:7422`). Bidirectional bus bridge PROVEN.
- **ORP / om-unity** — Rust fusion engine, 44 connector adapters, ORP-QL, DuckDB+RocksDB, 5,938 live AIS ships. (Note: currently a ghost-process recovery is in progress on the NEX edge stack — see internal ops notes; does not block CloudTAK work.)
- **TAK Server 5.7** — cluster on DOKS (namespace `takserver`, 8 pods), external LB `152.42.155.95` (ports 8089/8443/8446). NatsBridgePlugin (inbound) + cot-tap (outbound) both proven. Zeus's iPhone connected via iTAK, streaming live CoT.
- **CloudTAK** — namespace `cloudtak`, 6 pods (api, tiles, events, media, postgis, minio), at `cloudtak.onemindos.dev` via Cloudflare tunnel. Version 13.34.0. Patched for TAK 5.7 (node-tak SDK fixes). CoT display working (requires per-user ProfileConnConfig connection). Admin: `admin` / `OneMind-Admin-2026!`.
- **GeoServer** — DOKS, the sovereign geo/tile plane (kills Cesium Ion dependency; serves self-hosted tiles + static geo overlays).
- **Supporting:** ClickHouse (metrics/history), Grafana, ArgoCD (GitOps), Prometheus, Homepage dashboard, hermes-webui. All images on **GHCR** (`ghcr.io/onemindos/`). DOCR deleted.

---

## 12. Build Roadmap — The Order of Operations

### Phase 0 — Lock the contract (keystone)
Freeze the NATS subject taxonomy (§2.3) in `omos-infra` as the canonical interface. Every port after this aims at a fixed target.

### Phase 1 — Prove ONE vertical end-to-end (the proof)
Drive live AIS (om-unity's 5,938 ships) through the full loop:
```
AIS → NATS (sensor.ais.*) → ORP (ent.sea) → NatsBridgePlugin → TAK Server (CoT) → CloudTAK map + iTAK field
```
When ships appear on CloudTAK AND on the iPhone simultaneously, the architecture is proven. **This is the single most important next step** — it's mostly config (the pieces are all live).

### Phase 2 — Port sources one at a time
For each source (ADS-B, drones/MAVLink, home automation, OSINT feeds): wire to bus → confirm ORP fuses → confirm it renders on CloudTAK → done. One at a time, nothing missed. Reuse dfpc-coe ETL connectors where they fit (Path B); use Path A for everything operational.

### Phase 3 — AI layer
Fork TAKpilot, point it at AWS Bedrock, extend its tool library with OneMind-specific verbs (ORP-QL queries, mission tasking, home automation control). Legacy can now drive the map by voice/chat.

### Phase 4 — Branding + Profiles
Whitelabel CloudTAK as "OneMind" via the Config group. Configure Profile-based roles/permissions (operators now; family/kids later). Evaluate whether Profile+whitelabel covers the family case (likely retiring atoc permanently).

---

## 13. Open Questions / Decisions for the Implementing Agent

1. **ETL deploy shim** — build the CloudFormation→Kubernetes-Job shim to restore in-UI connector deploy, or just deploy connectors as GitOps manifests? (Recommendation: manifests first; shim only if the UI flow is needed.)
2. **NATS→CoT connector form** — a single generic `nats-to-cot` ETL Task (subject configurable per Layer) vs. one connector per subject? (Recommendation: one generic, parameterized by subject + a GeoJSON mapping.)
3. **ORP output → CoT mapping** — confirm the field mapping from ORP `ent.*` schema to CoT/GeoJSON (callsign, type, MIL-STD-2525 symbol code, position). This is the glue for Phase 1.
4. **CoT loop prevention** — TAK doesn't echo an injected CoT back to the injecting plugin; outbound proof needs a real device (cot-tap + a field EUD). Keep this in mind when testing round-trips.
5. **Branding central point** — locate CloudTAK's whitelabel/rename config (Config route group + any `brand.ts`-equivalent) before starting UI work.

---

## 14. Glossary (for the implementing agent)

- **OneMind / OMOS** — Zeus's sovereign unified C4ISR platform.
- **CoT (Cursor-on-Target)** — the TAK protocol's XML/protobuf message format for positions/events on a map.
- **TAK / ATAK / iTAK / WinTAK / WebTAK** — the Team Awareness Kit ecosystem (server + Android/iOS/Windows/web clients).
- **ORP** — Zeus's Rust fusion engine (`om-unity`); ingests protocols, fuses into an entity graph, queryable via ORP-QL.
- **NATS** — the messaging bus / fabric. JetStream = its persistence layer. Leaf node = a NATS server bridged into the cluster.
- **NatsBridgePlugin** — Zeus's TAK Server plugin; injects NATS `cot.*` messages into TAK (inbound).
- **cot-tap** — a mTLS client that reads TAK's CoT firehose and republishes to NATS (outbound).
- **ETL Layer / Task** — CloudTAK's data connector unit (AWS-Lambda-shaped by default; run as K8s workload for sovereignty).
- **DOKS** — DigitalOcean Kubernetes Service (the cloud hub).
- **GitOps / ArgoCD** — everything as version-controlled manifests, auto-deployed. OneMind's default doctrine.

---

*End of document. This is a complete, self-contained briefing — hand it to any agent and they have the full picture of the CloudTAK pivot, the NATS fabric, the AWS-Lambda ETL system, where data sources connect, and the build order.*
