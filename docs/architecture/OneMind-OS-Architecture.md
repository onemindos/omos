# OneMind OS — Unified Architecture

> **Review status (2026-10-17):** active OneMind OS design document, moved from the repo root and due for review and update. Until it is updated, where it disagrees with live state, `docs/STACK.md` (generated from the cluster) and `CLAUDE.md` win.

> **Last updated:** 2026-09-20
> **Status:** Active design document — source of truth for system unification

---

## Core Thesis

OneMind is not three repos bolted together. It is one system — a **sovereign, AI-native, real-time situational awareness platform** — with one bus, one operational picture, and one plugin model, deployable anywhere from a cloud cluster to a single laptop or forward edge node.

Think: **Ghost Recon + Watch Dogs** at the systems level. Every sensor, camera, robot, autonomous vehicle, and AI agent appears as an object on the same map. Voice commands trigger AI actions. AR overlays show bounding boxes from live ML inference. The operator sees everything; the AI acts on everything.

```
                    ┌───────────────────────────────────┐
                    │         UNIFIED COP (the truth)    │
                    │                                    │
                    │  Every sensor, every camera, every │
                    │  robot, every human, every AI      │
                    │  agent = a CoT marker on the map   │
                    │                                    │
                    └──────────────┬────────────────────┘
                                   │
                    ┌──────────────▼────────────────────┐
                    │          NATS (the nervous system)  │
                    │  sub-millisecond, edge-native,     │
                    │  offline-resilient, JetStream       │
                    └──────────────┬────────────────────┘
                                   │
              ┌────────────────────┼────────────────────────┐
              │                    │                         │
     ┌────────▼─────┐    ┌────────▼──────┐    ┌────────────▼──────┐
     │   CLIENTS     │    │   SERVERS      │    │   EDGE             │
     │               │    │                │    │                    │
     │ CloudTAK (web)│    │ TAK Server 5.8 │    │ mac-leaf           │
     │ ATAK (mobile) │    │ TAK-ML Server  │    │ robot-gateway      │
     │ AR glasses    │    │ Legacy AI      │    │ camera-bridge      │
     │ VR-TAK (sim)  │    │ NATS + CH      │    │ drone companion    │
     └───────────────┘    └────────────────┘    └────────────────────┘
```

---

## The Three Repos

| Repo | Path | Owns | Role |
|------|------|------|------|
| **omos-infra** (`omos/`) | `/Users/zeus/omos` | K8s cluster, NATS bus, TAK Server, all platform services, edge services, Go microservices, NATS stream definitions, Terraform, ArgoCD, monitoring, backups | GitOps source of truth |
| **onetak (CloudTAK)** | `/Users/zeus/onetak` | Web TAK client, ETL pipelines, REST API (Node.js/TS v13.78.0), browser-based SA display, MapLibre/CesiumJS rendering, plugin loader | Primary human display surface |
| **onemindos/VR** | `/Users/zeus/onemindos-VR` | Unreal Engine 5.5 immersive TAK client, VH_* closed-source plugin suite (Win64 DLLs), Quest 3 / Win64 builds, COT pipeline, Cesium integration | Simulation and VR training only |

---

## Repo Layout (omos-infra)

```
infra/              DO VPC + DOKS cluster + container registry (Terraform)
edge/               Non-K8s hardware-bound nodes
  mac-leaf/         Mac NATS leaf node config + launchd
  robot-gateway/    ROS2 ↔ NATS bridge (companion computer)
  omniverse-nats-bridge/  Isaac Sim USD prim mirror
services/           Go + Python NATS microservices
  task-engine/      Go — NATS-native task orchestration
  omd/              Go — OneMind Daemon (per-node supervisor)
  bridge/           Go — Universal protocol bridge (MAVLink, MQTT, HA, TAK, Webhook)
  nex/              Go — NATS Execution Engine workloads
  agents/meta/      Go — Meta agent service
  om_task.py        Python — Task subscriber
  om_presence.py    Python — Presence/voice bridge
  notify_router.py  Python — Routes evt.notify.> to notification channels
  graphite_bridge.py Python — NATS → Graphite/Carbon
platform/takserver/         TAK Server 5.8 (Helm chart + OneMind values)
  chart/            Helm chart + values-onemind.yaml
  docker/           Dockerfiles (base, api, messaging, config, plugins, db-setup)
  deploy/           K8s overlay manifests (LB, NATS leaf, secrets)
  plugin-src/       NatsBridgePlugin Java source
  plugin/           Built JAR + plugin config
  edge/             EUD management, device registry, simulators
platform/           K8s manifests by operational layer (see Deployment Waves)
docs/               Architecture + runbooks
```

---

## Deployment Waves

Everything in the cluster deploys in wave order. ArgoCD enforces sync order via App-of-apps.

| Wave | Directory | What Lives Here |
|------|-----------|----------------|
| 0 | `platform/core/` | ingress-nginx, cert-manager, sealed-secrets controller, Tailscale operator, GitHub Actions runner, cloudflared (standby) |
| 1 | `platform/nats/` | NATS 3-node HA (Helm), NACK JetStream controller, TLS, all 39 stream + KV definitions |
| 2 | `platform/onetak/` | OneTAK API/Events/Media/Tiles, PostGIS (onetak-db), Martin, Maputnik, Nominatim, Valhalla, WebODM, NodeODM, STAC |
| 2 | `platform/takserver/` | TAK Server 5.8 (Helm), FedHub, tak-ldap, NATS leaf sidecar |
| 3 | `platform/hermes/` | Legacy AI (Hermes) + copilot-runtime, Hermes web UI |
| 3 | `platform/omos/` | tak-projector, feed-intel, feed-satellites, sink-clickhouse, redpanda-connect, mission-engine |
| 3 | `platform/intel/` | OSINT services |
| 3 | `platform/geo/` | Tile38, STAC, additional geo tools |
| 3 | `platform/takml/` | TAK-ML Server (JAR), TAK-ML DB (PG16), Vision Worker |
| 3 | `platform/data/` | ClickHouse (Altinity operator), Redpanda Connect |
| 3 | `platform/vrtak/` | nginx file server (VR-TAK APK/ZIP distribution) |
| 4 | `platform/monitoring/` | Grafana, Prometheus, Loki, Tempo, Alertmanager, NATS Surveyor |
| 4 | `platform/nats-tools/` | nats-dashboard, agent-web-ui, agenttap, heynats, nats-ui |

---

## Live Cluster State

**Image registry:** All images at `ghcr.io/onemindos/` — DOCR was retired 2026-07-05.
**LB IP:** `143.244.220.155` (ingress-nginx)
**TAK LB:** `152.42.155.95` (direct unproxied)

| Service | Namespace | Image | Status |
|---------|-----------|-------|--------|
| NATS 3-node HA | nats | nats (Helm) | ✅ live |
| Legacy (Hermes AI) | hermes | ghcr.io/onemindos/hermes-legacy:latest | ✅ live |
| copilot-runtime | hermes | ghcr.io/onemindos/copilot-runtime:latest | ✅ live |
| OneTAK API | onetak | ghcr.io/onemindos/onetak-api:latest | ✅ live |
| OneTAK Events | onetak | ghcr.io/onemindos/onetak-events:latest | ✅ live |
| OneTAK Tiles | onetak | ghcr.io/onemindos/onetak-tiles:latest | ✅ live |
| OneTAK Media | onetak | ghcr.io/onemindos/onetak-media:latest | ✅ live |
| TAK Server 5.8 | takserver | ghcr.io/onemindos/takserver-base:* | ✅ live |
| ClickHouse | clickhouse | altinity/clickhouse-server (operator) | ✅ live |
| Grafana/Loki/Tempo | monitoring | (Helm) | ✅ live |
| nats-dashboard | nats | ghcr.io/onemindos/nats-dashboard:latest | ✅ live |
| agent-web-ui | nats | ghcr.io/onemindos/agent-web-ui:latest | ✅ live |
| agenttap | nats | ghcr.io/onemindos/agenttap:latest | ✅ live |
| heynats | nats | ghcr.io/onemindos/heynats:latest | ✅ live |
| nats-ui | nats | ghcr.io/gastbob40/nats-ui:latest | ✅ live |
| omos-tak-projector | omos | python:3.12-alpine (inline) | ✅ live |
| feed-intel | omos | ghcr.io/onemindos/feed-intel:latest | ✅ live |
| feed-satellites | omos | ghcr.io/onemindos/feed-satellites:latest | ✅ live |
| sink-clickhouse | omos | ghcr.io/onemindos/sink-clickhouse:latest | ✅ live |
| redpanda-connect | omos | ghcr.io/redpandadata/connect:latest | ✅ live |
| omos-mission-engine | omos | ghcr.io/onemindos/omos-mission-engine | 🔴 scaled to 0 |
| valhalla | onetak | valhalla/valhalla:run-latest | 🔴 scaled to 0 (awaiting tile build) |
| nominatim | onetak | mediagis/nominatim:4.4 | 🔴 crashloop — scale to 0 |
| cloudflared | cloudflare | cloudflare/cloudflared | 🔴 scaled to 0 (standby) |
| fedhub | takserver | ghcr.io/onemindos/fedhub:5.8 | ⏳ awaiting image + certs |
| tak-ldap | takserver | ghcr.io/onemindos/tak-ldap:5.8 | ⏳ awaiting image + SealedSecret |
| vrtak-releases | vrtak | nginx:1.27-alpine | ⏳ file server live; awaiting UE build upload |
| TAK-ML Server | takml | eclipse-temurin:21-jre | ✅ live (blocked by vision worker) |
| TAK-ML DB | takml | postgres:16 | ✅ live |
| TAK-ML Vision Worker | takml | python:3.11-slim | 🔴 Init:0/1 — pip-at-runtime broken |

---

## Namespaces

| Namespace | Purpose |
|-----------|---------|
| nats | NATS 3-node HA cluster + all NATS tooling |
| hermes | Legacy AI agent + copilot-runtime + web UI |
| onetak | OneTAK API/Events/Media/Tiles + all geo services |
| takserver | TAK Server 5.8 + FedHub + LDAP |
| clickhouse | ClickHouse (Altinity operator) |
| monitoring | Prometheus + Grafana + Loki + Tempo + Alertmanager |
| cloudflare | cloudflared (scaled to 0 — cold standby only) |
| omos | Runtime services: tak-projector, feeds, sinks, redpanda |
| geo | Geospatial stack (Tile38, STAC, additional tools) — **EMPTY, consolidate into onetak** |
| takml | TAK-ML vision pipeline |
| vrtak | VRTAK nginx distribution server |
| argocd | ArgoCD GitOps controller |
| cert-manager | TLS cert automation (Let's Encrypt via Cloudflare DNS) |
| tailscale | Tailscale proxy operator (~20 proxy pods) |
| github-runner | Self-hosted GitHub Actions runner |
| ingress-nginx | ingress-nginx LB controller |
| kube-system | sealed-secrets controller |

---

## Public Hostnames

All DNS: Cloudflare-proxied A records → `143.244.220.155` (ingress-nginx), except `tak.onemindos.dev`.

| Hostname | Backend | Status |
|----------|---------|--------|
| onemindos.dev | onetak-api:5000 | ⚠️ DNS update pending |
| tiles.onemindos.dev | onetak-tiles:5002 | ⚠️ DNS update pending |
| media.onemindos.dev | onetak-media:9997 | ⚠️ DNS update pending |
| ws.onemindos.dev | nats:8080 (WebSocket) | ⚠️ DNS update pending |
| auth.onemindos.dev | Authentik IdP | ✅ live |
| tak.onemindos.dev | tak-external LB `152.42.155.95` | ✅ live (unproxied, direct) |

**Tailscale-only (never public) — `*.tail717ba0.ts.net`:**

| Hostname | Service |
|----------|---------|
| argocd.tail717ba0.ts.net | ArgoCD |
| grafana.tail717ba0.ts.net | Grafana |
| hermes-agent.tail717ba0.ts.net | Hermes web UI |
| agent-web-ui-1.tail717ba0.ts.net | Agent Web UI |
| nats-dashboard-1.tail717ba0.ts.net | NATS Dashboard |

**Cluster access:**
```bash
export KUBECONFIG=~/.kube/onemind-hub
kubectl --context tailscale-operator.tail717ba0.ts.net get nodes
# or via doctl
doctl kubernetes cluster kubeconfig save onemind-hub
```

---

## NATS — The Nervous System

NATS is the only bus. No service-to-service HTTP in the data path. Every entity, every sensor, every command, every AI interaction flows through NATS.

### Bus Authority

**The K8s NATS cluster is the authority.** Everything publishes to it.

- **K8s pods** talk directly to `nats.nats.svc.cluster.local:4222` — no sidecar needed
- **Mac bridges** use the Mac NATS leaf node to forward to the K8s bus
- **Edge boxes / companion computers** use a local leaf node or sidecar to forward up
- The Mac is the **cockpit** — it runs hardware-bound bridges (MAVLink, HA, cameras) that must be local. The leaf node is a **forwarding pipe**, not a peer bus. No state, no logic at the edge.
- If the Mac is offline, the cloud bus keeps running — only hardware-bound feeds go dark

```
MAC (cockpit / hardware edge)
  MAVLink bridge    ─┐
  HA bridge         ─┤── NATS leaf ──► K8s NATS HA (3-node) ◄── all cloud services
  cameras/vision    ─┘   (forwarding      THE BUS              ORP, Legacy, projectors,
                          pipe only)                            connectors, ClickHouse
```

### NATS Access Pattern

**K8s pod (direct):**
```yaml
env:
  - name: NATS_URL
    value: "nats://nats.nats.svc.cluster.local:4222"
  - name: NATS_USER
    valueFrom:
      secretKeyRef:
        name: nats-creds
        key: user
  - name: NATS_PASSWORD
    valueFrom:
      secretKeyRef:
        name: nats-creds
        key: password
```

**Edge/Mac (leaf forwarding — always localhost):**
```yaml
env:
  - name: NATS_URL
    value: "nats://localhost:4222"   # leaf node always on localhost
```

The sidecar pattern is shown in `platform/hermes/deployment.yaml` — Legacy uses it because it needs isolated agent identity. Most K8s services do NOT need the sidecar.

---

## NATS Subject Grammar (LOCKED)

> Never invent new subject patterns. These cover every valid communication in the system.

```
tel.{node}.{class}.{metric}          telemetry          tel.drone.x47.position
det.{domain}.{type}.{id}             detections         det.air.visual.track42
evt.{domain}.{entity}.{event}        semantic events    evt.mission.patrol1.completed
cmd.{domain}.{entity}.{action}       commands           cmd.drone.x47.intercept
intel.{domain}.{source}.{type}       intelligence       intel.conflict.acled.event
ent.{entity_type}.{entity_id}.state  fused entity state ent.ship.mmsi123.state
svc.{service}.{endpoint}             NATS RPC           svc.hermes.query
agents.{verb}.{agent}.{owner}.*      agent fabric       agents.prompt.hermes.zeus.legacy
tak.cot.inject                       CoT → TAK inject   (NatsBridgePlugin input)
tak.cot.>                            CoT from TAK       (NatsBridgePlugin output)
voice.transcript / voice.say         voice layer        (omos-voice pipeline)
rad.{site}.{subsystem}.{id}.*        homestead/RADS     rad.homestead.soil.bed1.moisture
```

---

## NATS JetStream Streams (39 total)

All streams are NACK Stream CRDs in `platform/nats/streams/`. The NACK controller reconciles continuously.

### Persistent Data Streams (replicas: 3, storage: file)

| Stream | Name | Subjects | Retention | Purpose |
|--------|------|----------|-----------|---------|
| **TEL** | TEL | `tel.fleet.>` `tel.drone.>` `tel.robot.>` `tel.edge.>` `tel.iot.>` `tel.industrial.>` `tel.sat.>` | 168h | All fleet asset telemetry — GPS, heading, battery, speed |
| **COT** | COT | `tak.cot.>` | 720h | TAK CoT stream — ATAK position reports, bridge feeds |
| **CMD** | CMD | `cmd.>` | 720h | Every command to every actor — drones, robots, nodes, agents |
| **AGENTS** | AGENTS | `agents.hb.>` `agents.status.>` `agents.discover` `agents.heartbeat` | 168h | Agent fabric — prompts, heartbeats, status, tool calls |
| **EVT** | EVT | `evt.so.>` `evt.hp.>` `evt.le.>` `evt.ge.>` `evt.task.>` `evt.cloud.>` `evt.bio.>` `evt.robot.>` `evt.drone.>` `evt.sat.>` | 2160h (90d) | Semantic events — homestead, health, cloud, bio alerts, robot/drone/sat |
| **INTEL** | INTEL | `intel.>` | 2160h (90d) | Curated world intelligence — distilled alerts, geofence triggers, watch-zone events |
| **VOICE** | VOICE | `voice.>` | 720h | Voice sessions — transcripts, audio refs, intent routing |
| **ENV** | ENV | `tel.nature.>` `tel.bio.>` `tel.weather.>` `tel.soil.>` `tel.water.>` | TBD | Environmental + biometric telemetry — soil, weather, water, HRV, presence |
| **GEO** | GEO | `evt.tak.cot.>` `evt.tak.status.>` | TBD | TACMAP geo layer — normalized CoT + fleet positions |
| **MISSION** | MISSION | `evt.mission.>` `evt.map.>` | TBD | Mission lifecycle — status, waypoints, abort, geofence, map annotations |
| **SYS** | SYS | `evt.sys.>` `sys.>` `tel.host.>` | TBD | System health — node state, service status |
| **WORLD** | WORLD | `evt.world.>` | 24h / 2GB cap | Raw world intel — AIS, ADS-B, ACLED, wildfires, seismic |
| **ZEUS** | ZEUS | `zeus.>` `op.zeus.>` | 8760h (1yr) | Operator context log — durable record of Zeus's operational decisions |
| **SAAS** | SAAS | `saas.>` | TBD | SaaS I/O — Wombat inbound webhooks + outbound API request/response |

### Missing Streams (to add)

| Stream | Subjects | Why |
|--------|----------|-----|
| **ENT** | `ent.>` | Fused entity state persistence — currently fire-and-forget. Needed for replay and late-joining clients. |
| **DET** | `det.>` | Detection results from TAK-ML vision worker, camera inference, anomaly detection. |
| **CAM** | `tel.cam.>` | Camera metadata — feed status, PTZ position, active/inactive. Makes cameras appear as CoTs on map. |

### KV Buckets (25 total)

| Bucket | Name | Purpose |
|--------|------|---------|
| kv-agent-caps | agent-caps | Per-agent capability declarations (what each agent can do) |
| kv-agent-config | agent-config | Agent runtime configuration |
| kv-agent-registry | agent-registry | Live agent registry — what agents exist, their endpoints |
| kv-device-registry | device-registry | All registered devices (TAK EUDs, robots, cameras, sensors) |
| kv-feed-config | feed-config | Data feed configurations (intel feeds, world feeds) |
| kv-fleet-state | fleet-state | Latest known state for every fleet asset |
| kv-geofences | geofences | Defined geofence polygons and rules |
| kv-mission-state | mission-state | Active mission state (current waypoints, task graph, status) |
| kv-node-registry | node-registry | Edge node registry — mac-leaf, companion computers, Jetsons |
| kv-om-config | om-config | OneMind platform configuration |
| kv-om-meta-dedup | om-meta-dedup | Deduplication metadata for omni-media ETL |
| kv-om-schema | om-schema | ORP ontology schema definitions |
| kv-om-ui-projections | om_ui_projections | Retained latest ui.* projections from om-control-plane |
| kv-overwatch-intercept | OVERWATCH_INTERCEPT | om-overwatch interceptor engagement status |
| kv-overwatch-state | OVERWATCH_STATE | om-overwatch fused entity world-state (keyed by entity_id) |
| kv-overwatch-video | OVERWATCH_VIDEO | om-overwatch per-entity video stream availability |
| kv-projects | projects | Project definitions and metadata |
| kv-sys-health | sys-health | System health snapshots |
| kv-task-index | task-index | Task index for lookup and scheduling |
| kv-task-notify | task-notify | Task notification routing config |
| kv-webhook-registry | webhook-registry | Registered webhook endpoints |
| kv-world-state | world-state | Latest world state snapshot (AIS, ADS-B, etc.) |
| kv-zeus-context | zeus-context | Zeus operator context (current focus, active mission) |
| kv-zeus-focus | zeus-focus | Zeus current attention/focus target |
| kv-zeus-presence | zeus-presence | Zeus presence/online state |

---

## Entity State Authority Chain

Who owns entity state and how it flows through the system:

```
Physical World
      │
      │  GPS / sensors / TAK EUD / robot / camera / ML detection
      ▼
TAK Server 5.8  ◄────── NatsBridgePlugin (live, bidirectional)
   (WRITE MASTER                │
    for entity state)           │  tak.cot.> (all CoT out)
                                │
                         ┌──────▼──────┐
                         │    NATS      │  ←── also receives:
                         │    bus       │      ent.*.*.state (fused entities)
                         └──────┬───────┘      det.> (ML detections)
                                │              tel.drone/robot.> (platform assets)
              ┌─────────────────┼──────────────────────────┐
              │                 │                           │
    ┌─────────▼──────┐  ┌───────▼──────┐  ┌───────────────▼───────┐
    │  ClickHouse     │  │  KV Buckets  │  │  All Clients           │
    │  (archive)      │  │  (hot state) │  │  (CloudTAK, ATAK, VR) │
    │  every event    │  │  latest pos  │  │  render the COP        │
    │  ever written   │  │  per entity  │  │                        │
    └────────────────┘  └──────────────┘  └────────────────────────┘
              │
              │  replay queries, analytics, time-machine
              ▼
        Legacy AI (has full history)
```

**Rule:** TAK Server is the **write master** for CoT entity state. NATS is the transport. ClickHouse is the immutable append-only archive. KV buckets hold the hot/latest view. Clients never write state directly — they send commands via `cmd.>` or inject CoT via `tak.cot.inject`.

---

## Secrets Convention

All secrets are **Sealed Secrets** — encrypted SealedSecret YAMLs committed to repo. The `sealed-secrets` controller in `kube-system` decrypts with the cluster private key. ArgoCD applies automatically.

**To add or rotate a secret:**
```bash
# Fetch cluster public key
kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > /tmp/pub-cert.pem

# Create and seal
kubectl create secret generic <name> -n <namespace> \
  --from-literal=KEY=value \
  --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/<namespace>/<name>.sealedsecret.yaml

# Commit and push — ArgoCD applies it automatically
```

**SealedSecret location:** `platform/{namespace}/{secret-name}.sealedsecret.yaml`

| Secret | Namespace | Contents |
|--------|-----------|---------|
| hermes-env | hermes | AWS Bedrock creds + NATS creds (23 fields) |
| hermes-leaf-conf | hermes | NATS leaf.conf + ca.crt |
| hermes-nats-ca | hermes | NATS cluster CA cert |
| hermes-webui-auth | hermes | Hermes web UI password |
| hermes-dashboard-auth | hermes | Hermes dashboard auth |
| cloudtak-agent-token | hermes | CloudTAK API token (generate in CloudTAK UI) |
| onetak-secrets | onetak | Signing secret, MinIO keys, postgres DSN |
| omos-nats-creds | omos | NATS user+password for omos namespace |
| tak-nats-leaf | takserver | leaf.conf + ca.crt for TAK NATS bridge |
| cloudflare-api-credentials | cloudflare | Cloudflare API key + email |
| cloudflare-api-token | cert-manager | Cloudflare DNS API token |
| tunnel-creds-json | cloudflare | Tunnel credentials.json |
| nats/values.secret.yaml | nats | NATS ONEMIND + SYS passwords (Helm — out-of-band) |
| ghcr-pull-secret | hermes/omos | GHCR pull secret (out-of-band) |

---

## What to DROP

| Item | Why |
|------|-----|
| **omos-ui** | Dead — CloudTAK is the web client. All 20 omos-ui modules become CloudTAK plugins. |
| **`geo` namespace** | Empty — only holds secrets. Geo services already live in `onetak` ns. Move secrets, delete namespace. |
| **worldwideview** | CesiumJS globe absorbed into CloudTAK `plugin-3d-tiles`. Not a separate app. |
| **Cloudflared manifests** | Scaled to 0 for months — ingress-nginx handles everything. Keep as cold standby only. |
| **osiris** | Intelligence dashboard becomes CloudTAK `plugin-intel`. Not a separate Next.js app. |
| **VR-TAK as primary display** | Keep for VR/simulation training only. CloudTAK + CesiumJS 3D Tiles is the primary 3D experience for desktop/AR. |

---

## What to KEEP and Why

| Item | Why It Stays | Role |
|------|-------------|------|
| **TAK Server 5.8** | CoT protocol authority, military interop, federation, data packages | Write master for all entity state |
| **NATS 3-node HA** | The bus. Sub-ms pub/sub, JetStream persistence, edge leaf nodes, KV store | Real-time nervous system |
| **ClickHouse** | Time-series archive. Replay, analytics, audit trail — every NATS event ever | "Time machine" |
| **CloudTAK/OneTAK** | Web client, ETL pipelines, REST API, MapLibre rendering. CesiumJS 3D Tiles adds 3D. | Primary human display |
| **ATAK/iTAK** | Android field client. AR glasses run Android — ATAK plugins work on glasses. | Mobile/field/AR |
| **TAK-ML Server** | IoBT sensor + model execution framework. ONNX/PyTorch MX plugins. Connects to TAK Server natively. | AI/ML inference pipeline |
| **Legacy AI** | Sovereign AI agent (claude-opus-4 on Bedrock). Accessible from every surface via NATS. | Commander's AI |
| **Authentik** | Identity. Single SSO for all surfaces (OIDC + LDAP). Live at `auth.onemindos.dev`. | Auth authority |
| **FedHub** | TAK federation. Multi-site connectivity. | Inter-org bridge |
| **Monitoring stack** | Grafana/Loki/Tempo/Prometheus. | Observability |
| **VR-TAK** | Simulation, training, Isaac Sim integration — NOT everyday SA | Simulation/training only |

---

## Database Justification (no more, no fewer)

Six databases — each justified. The only consolidation is killing the empty `geo` namespace.

| Database | Namespace | Why Separate |
|----------|-----------|-------------|
| TAK PostgreSQL | takserver | CoT entities, missions, data packages. TAK Server owns schema — cannot merge. |
| ClickHouse | clickhouse | Time-series columnar. Access pattern (TSDB) is fundamentally different from relational. |
| PostGIS (onetak-db) | onetak | Tile sources, geo features, ETL state. CloudTAK's own schema. |
| TAK-ML DB | takml | Model registry, sensor metadata, MBT (Model Building Toolkit). ML lifecycle is its own domain. |
| FedHub MongoDB | takserver | Federation hub state. FedHub requires MongoDB specifically — cannot swap. |
| Authentik PG | auth | Identity, sessions, OIDC state. Authentik owns its schema. |
| ~~geo-db~~ | ~~geo~~ | **DROPPED** — duplicate of onetak PostGIS |

---

## Geo Server Map

| Service | Status | Verdict |
|---------|--------|---------|
| Martin (vector tiles from PostGIS) | Running | KEEP always-on — core map rendering |
| Maputnik (style editor) | Running | KEEP — dev tool, low resource |
| Tile38 (geospatial index) | Manifest only | KEEP — activate when geofencing live |
| Nominatim (geocoding) | Crashlooping | FIX: scale to 0, activate on-demand |
| Valhalla (routing) | Scaled to 0 | KEEP at 0 — on-demand only |
| Overpass (OSM queries) | Scaled to 0 | KEEP at 0 — on-demand only |
| STAC (imagery catalog) | Scaled to 0 | ACTIVATE when drone imagery pipeline live |
| NodeODM (photogrammetry) | Running | KEEP — drone image → 3D Tiles pipeline |
| WebODM | Running | KEEP — web UI for NodeODM (env: `WO_DB_HOSTNAME`) |

---

## The Unified Plugin Architecture

**One plugin model, six surfaces.** Every plugin follows the same pattern:

```
1. DECLARE  — capabilities manifest (what subjects you need, what you expose)
2. SUBSCRIBE — bus subjects you consume
3. PUBLISH  — results back to the bus
4. EXPOSE   — UI (client-side) or API (server-side)
```

Because NATS is the universal bus, a server-side plugin's output is immediately available to ALL clients without per-client wiring. Build once on the server, every surface sees it.

---

### Surface 1: TAK Server Plugins (Java, server-side)

Every capability added here is instantly visible to ATAK, CloudTAK, and any TAK client automatically.

| Plugin | Status | What It Does |
|--------|--------|-------------|
| **NatsBridgePlugin** | LIVE | TAK ↔ NATS bidirectional CoT wire. `tak.cot.inject` in, `tak.cot.>` out. Source in `services/takserver/plugin-src/`. |
| **TAK-ML Server Plugin** | DEPLOYED (blocked) | ML inference pipeline — IoBT sensors → MX models → detections → CoT injection. Blocked by vision worker init container. |
| **OneMind Fusion Plugin** | DESIGN | Single Java plugin: subscribes `ent.>`, `det.>`, `tel.cam.>` → injects as CoT markers. Makes cameras, robots, ML detections, and NATS entities all appear as TAK objects. Legacy AI appears as a TAK "contact" for chat. One plugin beats six separate JARs operationally (JAR loading is heavy). |

---

### Surface 2: ATAK/iTAK Plugins (Android, client-side)

AR glasses run Android — ATAK plugins run on glasses too (Xreal, RayNeo, Vuzix). This is the Watch Dogs AR overlay layer.

| Plugin | Priority | What It Does |
|--------|----------|-------------|
| **ATAK-Legacy** | HIGH | Chat with Legacy AI via TAK messaging. Voice PTT → omos-voice → NATS → Legacy → TAK chat reply. |
| **ATAK-Camera** | HIGH | View camera feeds as video windows within ATAK. Tap a camera CoT → opens RTSP stream. TAK-ML detection bounding boxes overlaid on feed. |
| **UAS Tool** | EXISTS upstream | Drone control panel. Wire to `cmd.drone.*` NATS subjects. |
| **ATAK-Fleet** | MEDIUM | Robot/drone fleet panel. Send `cmd.robot.{id}.move` commands. View telemetry. |

---

### Surface 3: CloudTAK/OneTAK Plugins (TypeScript, web client)

Desktop ops center. Replaces omos-ui entirely. CloudTAK plugin API: `install(app, api)` / `enable()` / `disable()`. `MenuItemConfig` uses `label`. Routes added in `enable()`, removed in `disable()` via `api.router.removeRoute()`.

| Plugin | Package | Status | Replaces | What It Does |
|--------|---------|--------|----------|-------------|
| **plugin-3d-tiles** | @onemindos/cloudtak-plugin-3d-tiles | DESIGN | globe module | CesiumJS 3D Tiles layer — photorealistic terrain, drone orthomosaics, OSM buildings. Self-hosted tiles from NodeODM → STAC pipeline. |
| **plugin-hermes** | @onemindos/cloudtak-plugin-hermes | Scaffolded | agents module | Full Legacy AI assistant — chat, tool approve/deny, context-aware queries. Float pane. NATS: `agents.prompt/response/tool.hermes.zeus.legacy` |
| **plugin-clickhouse** | @onemindos/cloudtak-plugin-clickhouse | Scaffolded | replay module | ClickHouse time-range replay — scrub timeline, see every entity position ever. "Time machine." |
| **plugin-mumble** | @onemindos/cloudtak-plugin-mumble | Scaffolded (needs proxy) | voice module | PTT interface, voice transcripts, Legacy voice responses. WebSocket proxy at `wss://mumble.onemindos.dev`. |
| **plugin-ops** | @onemindos/cloudtak-plugin-ops | Scaffolded | overview module | Inbox, schedule, feeds, log, projects. NATS: `onemind.inbox.>`, `onemind.tasks.>`, `onemind.alerts.>` |
| **plugin-cameras** | — | DESIGN | — | Multi-camera grid viewer. Camera CoTs on map → click to open live RTSP/HLS. ML detection bounding box overlay. |
| **plugin-fleet** | — | DESIGN | fleet module | Robot + drone control panel. Command, telemetry, waypoint mission planning. |
| **plugin-intel** | — | DESIGN | intel module | OSINT feeds, intelligence overlay — replaces standalone osiris. |
| **plugin-nats** | — | LIVE | — | NATS bus inspector (existing live plugin). |
| **plugin-legacy-ai** | — | LIVE | — | Existing Legacy AI CloudTAK plugin (live). |

---

### Surface 4: VR-TAK Plugins (Unreal Engine 5.5, simulation only)

VR-TAK is for **simulation and training only** — not everyday situational awareness. The primary 3D display is CloudTAK + CesiumJS 3D Tiles. VR-TAK runs on Win64 and Quest 3.

| Plugin | Status | Notes |
|--------|--------|-------|
| **VH_* plugin suite** | EXISTING (closed-source) | Win64-only DLLs. `VH_Globe` is the globe reference actor. Custom UE5 source build with `EngineAssociation: ""`. |
| **CesiumForUnreal v2.28.0** | Added (PR #3) | UE 5.5 compatible. Pin to v2.28.0 — v2.29.1 requires UE 5.6+. Submodule at `Game/Plugins/CesiumForUnreal`. |
| **UGame_CesiumGlobeAdapter** | Added (PR #3) | C++ adapter wiring `ACesiumGeoreference` ↔ `VH_GlobeReferenceActor`. Blueprint-callable `SyncOriginToCesium()`. |
| **VR-NATS Plugin** | FUTURE | Direct NATS from UE5 for live sim ↔ real-world bridge. Lets Isaac Sim entities appear as CoTs on the live COP. |

**VR build notes:**
- Upstream: `github.com/TAK-Product-Center/VR` — sync branch `5.5-RELEASE-0`
- LFS: `Data/dted0.exe` (170MB) not on upstream LFS server — use `GIT_LFS_SKIP_SMUDGE=1` on merge
- CI: `.github/workflows/build-android.yml` builds Quest 3 APK + Win64 on `[self-hosted, ue55, windows]` runner
- Android target: ARM64 Shipping, ASTC textures, OpenXR, Quest 3

---

### Surface 5: Edge Plugins (Python/Go, hardware bridges)

| Plugin | Status | What It Does |
|--------|--------|-------------|
| **robot-gateway** | PR open (omos-infra PR #12) | ROS2 ↔ NATS bridge. Companion computer on UGV/UAV. Publishes `tel.{type}.{id}.position/status`, `ent.{type}.{id}.state`. Subscribes `cmd.{type}.{id}.move/stop`. `SIM=true` for Isaac Sim mode. |
| **camera-bridge** | DESIGN | RTSP/ONVIF camera discovery → registers each camera as a CoT marker, streams frames to TAK-ML vision worker for inference, publishes `tel.cam.{id}.*`. |
| **omniverse-nats-bridge** | PR open (omos-infra PR #12) | Isaac Sim Omniverse Kit extension. Subscribes `ent.*.*.state`, `tel.robot.>`, `tel.drone.>` → creates/updates USD Xform prims at `/World/Entities/{type}/{id}`. Mirrors the live COP into the sim. |
| **services/bridge/** | EXISTING (Go) | Universal protocol bridge — MAVLink (drone/UAS control), MQTT (IoT sensors), HomeAssistant (smart home/Homestead), TAK (CoT ↔ NATS), Webhook (external event ingest). All translate to/from NATS subjects. |
| **mac-leaf** | LIVE | Mac NATS leaf node — `edge/mac-leaf/`. Forwards all hardware-bound bridge traffic to K8s NATS. launchd service for persistence. |

---

### Surface 6: Legacy AI Tools (Agent capabilities)

Legacy AI uses NATS-native tools. Each tool is a `svc.{service}.{endpoint}` RPC call that Legacy can invoke. This is the action layer — Legacy doesn't just answer questions, it acts.

| Tool | What It Does |
|------|-------------|
| **tak-tools** | Inject CoT markers, query entities, create missions, set geofences, send TAK messages |
| **camera-tools** | Request a frame from any camera, run inference, describe the scene in natural language |
| **fleet-tools** | Command robots/drones (`cmd.robot.*`/`cmd.drone.*`), query fleet state from KV |
| **nats-tools** | Publish to or subscribe any NATS subject, query KV buckets, inspect JetStream state |
| **map-tools** | Add overlays, bounding boxes, annotations, polygons to the COP |
| **voice-tools** | TTS responses via Mumble (`voice.say`), voice-targeted notifications |
| **clickhouse-tools** | Query time-series data, generate reports, replay events, answer "what happened at 14:23?" |

---

## omos-voice Pipeline

Voice is the primary interface for AR glasses and field ops. Fully sovereign — no cloud STT/TTS.

```
Mumble PTT ──► RealtimeSTT (VAD + faster-whisper) ──► voice.transcript (NATS)
                  0.36s STT latency (proven)                    │
                                                                 ▼
                                                      Legacy AI (subscribed)
                                                                 │
                       ┌─────────────────────────────────────────┘
                       │
              voice.say ──► Piper TTS (local, sovereign) ──► Mumble
                                                          (spoken in earpiece)
```

**Tech stack:**
- STT: `RealtimeSTT` library with `faster-whisper` backend — VAD-gated (only transcribes when speech detected)
- TTS: Piper TTS — local, runs on CPU, no cloud dependency
- Transport: Mumble (Murmur server in cluster) — PTT channel management, multi-user voice
- NATS subjects: `voice.transcript` (STT output → Legacy), `voice.say` (Legacy → TTS)
- Proven latency: **0.36s** from end of speech to transcript on NATS

**Next step:** Legacy subscribes `voice.transcript`, processes intent, replies via `voice.say` + optional `tak.cot.inject`.

---

## Communication with Legacy — Three Channels

All three channels converge on the **same Legacy AI instance** via NATS. Legacy has full context regardless of input channel.

| Channel | How | Best For |
|---------|-----|----------|
| **Voice (PTT)** | Mumble → faster-whisper STT → `voice.transcript` → Legacy → `voice.say` → Piper TTS → Mumble | Hands-free, AR glasses, field ops, driving |
| **TAK Chat** | ATAK/CloudTAK → TAK Server → `tak.cot.>` → NATS → Legacy → `tak.cot.inject` (reply CoT) | Text, shared group channels, auditable log trail |
| **Direct NATS** | CloudTAK plugin-hermes → `agents.prompt.hermes.zeus.legacy` → Legacy → `agents.response.*` | Power user, complex multi-step queries, tool chains |

---

## TAK-ML IoBT Framework

> Source: Raytheon BBN — TAK Conference 2021. Problem statement: TAK devices + IoBT sensors = massive amounts of data. That data is almost entirely underused. TAK-ML fills the gap.

**Three plugin types — all can run client-side (ATAK) OR server-side (TAK-ML Server) interchangeably:**

| Plugin Type | Role | Examples |
|-------------|------|---------|
| **Sensor Plugin** | Ingests data from physical or virtual sensors | Camera (RTSP/USB), accelerometer, GPS track, LIDAR point cloud, radar, SDR |
| **MX Plugin (Model Execution)** | Runs ML models against sensor data. Frameworks: ONNX, PyTorch, TensorFlow, Weka, SparxML, DL4J, Spark | YOLOv5 object detection (ONNX), terrain classification (TF), gait-based auth (Weka), edible plant classifier (TF) |
| **Application Plugin** | Consumes model outputs, drives UX or downstream actions | Bounding box overlays on ATAK map, geofence trigger on detection, auto-mission abort on anomaly |

**Deployed MX plugins (omos-infra `platform/takml/`):**
- ONNX Runtime — loaded, production-ready
- PyTorch — loaded, production-ready
- **Vision Worker** — YOLOv5s VisDrone model, publishes `fabric.vision.{feed_id}.detect` — **blocked** by init container pip-at-runtime (see Priority #1)

**MX Plugin performance (from TAK Conference data):**
- ONNX: ~3ms prediction
- Weka: ~31ms prediction
- TensorFlow: ~4210ms prediction (CPU; GPU dramatically faster)

**NodeDrop (Raytheon BBN):** Sparse network compression technique for MX plugins deployed to ATAK devices. Achieves 97–99% model weight pruning with minimal accuracy loss — allows running inference on Class I/II UAS companion computers and phone-class hardware. Critical for field/edge deployments.

**Model Building Toolkit (MBT):** TAK-ML includes a toolkit for custom model training: SPARQL queries for training data selection from the entity/sensor history, model registration API, training pipeline integration. The TAK-ML DB (PG16 in `takml` namespace) is the model registry.

---

## The Camera / Sensor Pipeline (IoBT Architecture)

```
ANY CAMERA
(body cam, head cam,          ┌──────────────────────────┐
 robot-mounted, security,     │   camera-bridge            │
 drone, vehicle, greenhouse)  │   (edge plugin — Python)   │
        │                     │                            │
        │ RTSP/ONVIF           │ 1. Discovers ONVIF cameras │
        └────────────────────►│ 2. Registers each as CoT   │
                              │    (tak.cot.inject)         │
                              │ 3. Streams frames to        │
                              │    TAK-ML vision worker     │
                              │ 4. Publishes tel.cam.{id}.* │
                              │    (feed status, PTZ pos)   │
                              └────────────┬───────────────┘
                                           │ frames + metadata
                                  ┌────────▼─────────┐
                                  │  TAK-ML Vision    │
                                  │  Worker           │
                                  │                   │
                                  │  YOLOv5s / ONNX  │
                                  │  object detect    │
                                  │                   │
                                  │  → det.cam.{id}   │
                                  │  → tak.cot.inject │
                                  └────────┬──────────┘
                                           │
              ┌────────────────────────────┼──────────────────────────┐
              │                            │                           │
     ┌────────▼──────┐          ┌──────────▼────────┐   ┌────────────▼──────┐
     │  TAK Server    │          │  ClickHouse        │   │  Legacy AI         │
     │  CoT marker    │          │  archived for      │   │  understands what  │
     │  on map        │          │  replay + queries  │   │  cameras see       │
     └──────────┬─────┘          └────────────────────┘   └────────────────────┘
                │
     Visible on ALL clients:
     ATAK, CloudTAK, AR glasses
     (bounding boxes on HUD)
```

**Full data flow: sensor → detection → human sees it in under 1 second.**

---

## AR Glasses Scenario (Watch Dogs / Ghost Recon Vision)

Android-based AR glasses (Xreal, RayNeo, Vuzix Blade 2) run ATAK with OneMind plugins. This is the primary field operator experience.

```
AR Glasses (Android-based: Xreal / RayNeo / Vuzix)
│
├── ATAK installed with:
│     ATAK-Legacy  (AI voice/chat)
│     ATAK-Camera  (video + bounding box overlays)
│     UAS Tool     (drone control)
│     ATAK-Fleet   (robot control)
│
│  Field operator says:
│  "Hey Legacy, what's going on in the greenhouse?"
│
│  PTT → Mumble → omos-voice (faster-whisper STT, 0.36s)
│       → voice.transcript (NATS)
│              ↓
│         Legacy AI queries:
│           tel.cam.greenhouse.*     (camera feed status)
│           det.cam.greenhouse.*     (recent ML detections)
│           ClickHouse               (last 30 min history)
│           kv-fleet-state           (any assets nearby)
│
│  Legacy responds:
│    voice.say → Piper TTS → Mumble (spoken in ear)
│    tak.cot.inject → bounding boxes appear on AR HUD
│
│  What operator sees on AR display:
│    - Bounding boxes on detected objects (person, vehicle, tool)
│    - Entity labels on every person / vehicle in FOV
│    - "Greenhouse: 72°F, 65% humidity.
│       Two people detected in Zone A.
│       Last motion event: 4 minutes ago.
│       No anomalies in last 30 min."
│    - Alerts auto-appear as CoT icons with priority colors
│
│  Ghost Recon / Watch Dogs parallel:
│    - Every building, vehicle, person = tagged entity on HUD
│    - AI tactical advisor responds in real-time to questions
│    - Drone overhead feeding computer vision to the same HUD
│    - No heads-down — eyes on target, AI in ear
```

---

## Agent Roster

Authority chain: `Zeus → Legacy Mother → All Servant Agents → Tools`

All servant agent cron jobs **PAUSED** pending explicit resumption. Legacy Mother is the sovereign meta-agent and the only always-on AI identity.

| Agent | Role | Model | Location |
|-------|------|-------|----------|
| **Legacy Mother** | Sovereign meta-agent — commander's AI, orchestrates all other agents | claude-opus-4 (AWS Bedrock) | DOKS (hermes namespace) |
| **Legacy Daughter/Edge** | Edge body — runs on Mac, local hardware access | claude-opus-4 (AWS Bedrock) | Mac |
| **Forge** | Coding orchestrator — multi-agent code generation | Opus + qwen3-coder-flash workers | DOKS |
| **Oracle** | Intelligence + analytics | — | — |
| **Guardian** | Security sentinel — monitors threats, IAM, anomalies | — | — |
| **Grid** | Infra manager — K8s, NATS, services health | — | — |
| **Spartan** | Personal performance — fitness, training protocols | — | — |
| **Haven** | Home / HA manager — Homestead automation | — | — |
| **Heritage** | Records keeper — archives, documentation | — | — |
| **Empire** | Business admin — scheduling, comms, projects | — | — |
| **Eden** | Sustainability / RADS — soil, irrigation, growing systems | — | — |
| **Trinity** | Cross-system integration — bridges between agent domains | — | — |

**NATS agent fabric subjects:**
- Heartbeat: `agents.heartbeat` (broadcast), `agents.hb.{agent}` (per-agent)
- Status: `agents.status.{agent}`
- Discovery: `agents.discover`
- Prompt/Response: `agents.prompt.{verb}.{agent}.{owner}` / `agents.response.{verb}.{agent}.{owner}`

---

## Go Services (services/)

All services connect to `nats://nats.nats.svc.cluster.local:4222`, speak the OneMind subject grammar, deploy via `ghcr.io/onemindos/{service}:latest` with manifests in `platform/`.

| Service | Path | Purpose |
|---------|------|---------|
| **task-engine** | `services/task-engine/` | NATS-native task orchestration — create, dispatch, track tasks. Successor to `om_task.py`. |
| **omd** | `services/omd/` | OneMind Daemon — runs on every node, supervises local agent harnesses |
| **bridge** | `services/bridge/` | Universal protocol bridge — all external protocols → NATS |
| **nex** | `services/nex/` | NATS Execution Engine (NEX) workloads — overwatch bridges, om-task |
| **agents/meta** | `services/agents/meta/` | Meta agent service |
| **cot2constellation** | `services/bridge/bridges/tak/` | TAK CoT → Overwatch bridge |

**Bridge types in `services/bridge/bridges/`:**

| Bridge | Protocol | What It Bridges |
|--------|----------|----------------|
| **mavlink** | MAVLink v2 | Drone/UAS flight controllers → `tel.drone.*`, `cmd.drone.*` |
| **mqtt** | MQTT | IoT sensors, Homestead devices → `tel.iot.*`, `tel.nature.*` |
| **homeassistant** | HA WebSocket API | Smart home entities → `tel.nature.*`, `evt.hp.*` |
| **tak** | CoT XML | TAK Server CoT ↔ NATS bidirectional (complement to Java NatsBridgePlugin for Python side) |
| **webhook** | HTTP POST | External event webhooks → NATS subjects |

**Python subscribers (legacy operational scripts, run via `subscriber_runner.py`):**

| Script | Purpose |
|--------|---------|
| `om_task.py` | Task Engine subscriber (Python predecessor of task-engine Go service) |
| `om_presence.py` | Presence / voice bridge subscriber |
| `notify_router.py` | Routes `evt.notify.>` events to notification channels |
| `graphite_bridge.py` | NATS telemetry → Graphite/Carbon |
| `netdata_exporter.py` | NATS metrics → Netdata |
| `create_streams.py` | Provisioning: creates all JetStream streams |
| `create_buckets.py` | Provisioning: creates all JetStream KV buckets |
| `seed_rules.py` | Seeds initial rules/config into KV buckets |

---

## Local / Offline Deployment Model

The same stack runs on a single machine (laptop, Intel NUC, Jetson AGX) via Docker Compose or k3s. This is how forward edge deployments work — the cloud hub is optional.

```
┌─────────────────────────────────────────────────────┐
│          OneMind Local (Docker Compose / k3s)        │
├─────────────────────────────────────────────────────┤
│  NATS (single node, JetStream enabled)               │
│  TAK Server (single pod, no HA)                      │
│  CloudTAK (local web client → localhost)             │
│  ClickHouse (single node — time machine still works) │
│  Legacy AI (local LLM: Llama 3.1 70B / Qwen 2.5)   │
│  TAK-ML (local inference: ONNX CPU or CUDA)         │
│  camera-bridge (local cameras via RTSP/ONVIF)        │
│  omos-voice (local faster-whisper + Piper TTS)       │
│  PostGIS + Martin (local tile serving)               │
│  Authentik (local IdP)                               │
└─────────────────────────────────────────────────────┘
         │
         └── NATS leaf → cloud hub when online
             Offline: fully functional standalone
             Back online: JetStream replays missed messages
             (JetStream is the sync mechanism — no custom sync code)
```

**Offline resilience logic:**
1. While offline: NATS JetStream stores everything locally in sequence
2. When reconnected: leaf node replicates missed messages to cloud hub automatically
3. Late-joining clients: subscribe to JetStream consumer, replay from last sequence
4. This works for all streams (TEL, COT, CMD, EVT, INTEL, etc.)

---

## 3D Tiles 2.0 + Drone Imagery Pipeline

3D Tiles 2.0 is the latest Cesium standard for photorealistic 3D terrain. OneMind generates its own 3D Tiles from drone imagery.

```
Drone flight (PX4 / DJI / ArduPilot)
    │
    │ captures overlapping RGB images
    ▼
NodeODM (photogrammetry — running in onetak ns)
    │
    │ generates orthomosaic + point cloud + DSM
    ▼
STAC catalog (imagery index — in geo ns)
    │
    │ tiles available via STAC API
    ▼
CesiumJS 3D Tiles layer (CloudTAK plugin-3d-tiles)
    │
    │ rendered in browser as photorealistic 3D terrain
    ▼
Operator sees: accurate, up-to-date 3D map of surveyed area
Legacy AI sees: current terrain via camera-tools + map-tools
```

**Client integration:**
- CloudTAK: `plugin-3d-tiles` — CesiumJS layer on top of MapLibre base map
- VR-TAK: CesiumForUnreal v2.28.0 — Cesium globe in UE5 scene, origin synced via `UGame_CesiumGlobeAdapter`

---

## TAK Server — NatsBridgePlugin Architecture

The NatsBridgePlugin is the live CoT↔NATS wire. It is the most important integration point in the system.

```
ATAK / iTAK / CloudTAK / any TAK client
    │
    │  CoT XML over TCP/TLS (port 8089/8443)
    ▼
TAK Server 5.8
    │
    │  NatsBridgePlugin (Java, loaded at TAK Server startup)
    │
    ├── OUTBOUND: all CoT → tak.cot.{callsign}.{type} (NATS JetStream COT stream)
    │     → consumed by: ClickHouse sink, Legacy AI, robot-gateway mirror, CloudTAK ETL
    │
    └── INBOUND: tak.cot.inject (NATS) → injects as CoT into TAK Server
          → produces by: TAK-ML vision worker, OneMind Fusion plugin, Legacy AI tak-tools
          → result: detection bounding boxes, robot/drone positions, AI-generated markers
               appear on ATAK/CloudTAK/AR HUD automatically
```

---

## Authentik — Identity Authority

Single IdP for all surfaces. Live at `auth.onemindos.dev`.

| Surface | Auth Method | Status |
|---------|-------------|--------|
| CloudTAK | OIDC | Configured |
| ArgoCD | OIDC | Live |
| Grafana | OIDC | Live |
| ATAK | LDAP outpost | Pending — outpost not yet deployed |
| TAK Server | LDAP | Pending — requires tak-ldap service |
| Legacy web UI | Password | Live (local Authentik user) |

**Authentik namespace:** `auth` (live at `auth.onemindos.dev`)
**LDAP outpost:** needed to wire ATAK EUD enrollment to Authentik user directory.

---

## Adding a New Service (Checklist)

Every new service follows this exact pattern — no exceptions:

```
1. Image: build and push to ghcr.io/onemindos/{service}:latest
2. Directory: platform/{wave}/{service}/
   - namespace.yaml      (if new namespace needed)
   - deployment.yaml     (Deployment + Service; stateful adds PVC)
   - ingress.yaml        (if public-facing — nginx class, CF IP whitelist)
   - *.sealedsecret.yaml (sealed secrets, committed to repo)
3. ArgoCD: add Application to platform/argocd/ pointing to platform/{wave}/{service}/
4. Secrets: seal with kubeseal, commit SealedSecret YAML
5. NATS subjects: follow subject grammar — check §NATS Subject Grammar before creating any new pattern
6. Services table: add to CLAUDE.md Services table with status
7. Check: does this duplicate an existing layer? Check Services table first.
```

---

## The Deployment Pattern (how every service lands)

```
1. Image: ghcr.io/onemindos/{service}:latest  (GHCR only — DOCR retired 2026-07-05)
2. Directory: platform/{service}/
   namespace.yaml + deployment.yaml + ingress.yaml (if public)
3. Secrets: sealed via kubeseal → *.sealedsecret.yaml → committed → ArgoCD applies
4. NATS: always nats://nats.nats.svc.cluster.local:4222 for K8s pods
         always nats://localhost:4222 for edge/leaf nodes
5. ArgoCD syncs continuously — no manual kubectl apply in production
```

---

## Priority Execution Order

| # | Task | Impact | Effort | State |
|---|------|--------|--------|-------|
| 1 | **Fix TAK-ML vision worker** — bake torch/torchvision/opencv into Docker image instead of pip-at-runtime init container. Init:0/1 means the pipeline is dead. | HIGH — unlocks entire IoBT camera→detection→CoT pipeline | LOW | BLOCKED |
| 2 | **Add ENT/DET/CAM NATS streams** — JetStream persistence for entity state, detections, camera metadata | HIGH — required by everything downstream (fusion plugin, ClickHouse sink, late-joiners) | LOW | PENDING |
| 3 | **CloudTAK plugin-3d-tiles** — CesiumJS 3D Tiles layer in web client | HIGH — primary 3D display surface, replaces worldwideview | MEDIUM | DESIGN |
| 4 | **camera-bridge edge service** — RTSP/ONVIF → NATS + TAK-ML pipeline. Every camera becomes a CoT. | HIGH — unlocks IoBT camera pipeline | MEDIUM | DESIGN |
| 5 | **CloudTAK plugin-hermes** — wire Legacy AI chat into CloudTAK (scaffolded, needs wiring) | HIGH — AI from primary client | MEDIUM | SCAFFOLDED |
| 6 | **ATAK-Legacy plugin** — TAK chat ↔ Legacy AI from field/AR glasses | HIGH — AI from field/AR | MEDIUM | DESIGN |
| 7 | **Wire Authentik LDAP outpost** — ATAK EUD enrollment via Authentik | MEDIUM — unified identity | MEDIUM | PENDING |
| 8 | **Consolidate geo namespace** — move secrets to onetak, delete empty ns | LOW — cleanup | LOW | PENDING |
| 9 | **Scale nominatim to 0** — stop the crashloop | LOW — cleanup | LOW | PENDING |
| 10 | **Local deployment Docker Compose** — single-machine OneMind for edge deployment | MEDIUM — field deployment model | HIGH | DESIGN |
| 11 | **OneMind Fusion TAK Server Plugin** — single Java plugin making all NATS entities visible as CoTs in TAK | HIGH — unified COP on server side | HIGH | DESIGN |
| 12 | **mumble-ws-proxy** — WebSocket proxy for CloudTAK plugin-mumble PTT | MEDIUM — voice from web client | LOW | PENDING |

---

## Service Roadmap

| Wave | Service | Purpose |
|------|---------|---------|
| 2 | omos-mission-engine | NATS mission orchestrator — workflow graphs dispatching `cmd.*`, currently scaled to 0 |
| 3 | camera-bridge | RTSP/ONVIF camera bridge → NATS + TAK-ML (see Priority #4) |
| 3 | omos-feed-world | Automated world data feeds (AIS, ADS-B, ACLED, wildfires, seismic) |
| 3 | omos-feed-intel | OSINT intel feeds — events, conflict, geopolitical |
| 4 | Domain bridges | New protocol bridges: ROS 2, wellness data, scheduler, WebRTC for cameras |
| 4 | mumble-ws-proxy | Mumble WebSocket proxy for CloudTAK PTT plugin |
| 5 | ATAK-Legacy plugin | Android TAK plugin for Legacy AI chat |
| 5 | OneMind Fusion Plugin | TAK Server Java plugin — all NATS entities as CoTs |

---

## Key Design Rules

1. **Build on the server, let every client just be a view.** TAK Server plugins and NATS microservices are the force multipliers. Every client sees the result automatically without per-client wiring.
2. **Every sensor, every device, every robot = a CoT on the map.** TAK is the common entity protocol. TAK-ML + NatsBridgePlugin + OneMind Fusion Plugin make this happen for every entity type.
3. **NATS is the only bus.** No direct service-to-service HTTP in the data path. Everything speaks NATS. Clients that can't speak NATS natively get a bridge (see Surface 5).
4. **One plugin model across all six surfaces.** Declare → Subscribe → Publish → Expose. Same pattern whether you're writing a Java TAK plugin, a TypeScript CloudTAK plugin, a Python edge bridge, or a UE5 component.
5. **Offline first.** Any deployment tier (cloud, local, edge node) must function standalone and sync when connected. NATS JetStream handles the sync — no custom sync code.
6. **No database without justification.** Six databases exist, each with a fundamentally different access pattern or locked schema. No new databases without clearing that bar.
7. **TAK Server is the write master for entity state.** Everything that produces a position or entity state funnels through TAK Server (directly or via `tak.cot.inject`). NATS is transport. ClickHouse is the archive. KV is hot cache.
8. **Sovereign everywhere.** No cloud STT (faster-whisper). No cloud TTS (Piper). No cloud LLM for operational use (local Llama/Qwen for offline, Bedrock for cloud). The system must work in a degraded/denied environment.
9. **Images always from GHCR.** All images at `ghcr.io/onemindos/` — DOCR was retired 2026-07-05. Never use DOCR.
10. **Never hardcode NATS addresses in edge containers.** Always `nats://localhost:4222` — the leaf sidecar handles routing. This makes containers portable across deployment tiers.

---

## What NOT to Do

- **Never commit secrets or credentials** (even test values) to any repo — use SealedSecrets
- **Never deploy imperatively** without a matching manifest in omos-infra
- **Never create a new NATS subject pattern** without following the grammar in §NATS Subject Grammar
- **Never use DOCR** — all images are on GHCR (`ghcr.io/onemindos/`)
- **Never hardcode NATS server addresses** in containers — always `nats://localhost:4222` via leaf
- **Never add a service that duplicates an existing layer** — check the Services table first
- **Never treat the Mac NATS leaf as the source of truth** — K8s NATS is the authority
- **Never deploy to VR as a primary display** — CloudTAK + CesiumJS is the everyday SA surface; VR is simulation/training only
