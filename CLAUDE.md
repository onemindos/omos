# OneMind Infra — Claude Context

This repo is the infrastructure-as-code for the OneMind platform. It manages the DOKS (DigitalOcean Kubernetes) cloud hub, platform service manifests, NATS stream definitions, TAK cluster config, and Terraform for the DO foundation.

## What This Repo Is

`omos-infra` is the **GitOps source of truth** for everything deployed to the cloud hub. If a service runs in K8s, its manifest lives here. If infrastructure exists, Terraform declares it here. Nothing is deployed imperatively without a matching file in this repo.

The platform is a NATS-native C2/ISR system — think Palantir AIP or Lattice OS built on open protocols. NATS is the nervous system. ORP is the fusion/ontology engine. TAK (iTAK/ATAK/WebTAK/CloudTAK) is the tactical display. Legacy (Hermes AI) is the commander.

## Bus Authority (critical — read this first)

**The K8s NATS cluster is the main bus. It is the authority. Everything publishes to it.**

- All workloads, connectors, and services deploy to K8s and talk directly to the in-cluster NATS at `nats.nats.svc.cluster.local:4222`
- The Mac is the **cockpit** — it runs hardware-bound bridges (MAVLink, Home Assistant, local cameras) that have no choice but to be local. Those bridges use the Mac NATS leaf node to forward everything they receive up to the K8s bus.
- The Mac leaf node is a **forwarding pipe**, not a peer bus. The Mac holds no state, runs no logic. It just ensures that what happens at the hardware edge reaches the cloud brain.
- If the Mac is offline, the cloud bus keeps running. Missions continue. ORP keeps fusing. Legacy keeps watching. Only the hardware-bound feeds go dark.
- Never design a service that treats the Mac NATS as the source of truth. The K8s cluster is the source of truth.

```
MAC (cockpit / hardware edge)
  MAVLink bridge    ─┐
  HA bridge         ─┤── NATS leaf ──► K8s NATS HA (3-node) ◄── all cloud services
  Vision/cameras    ─┘   (forwarding      THE BUS              ORP, Legacy, projectors,
                          pipe only)                            connectors, ClickHouse
```

## Repo Layout

```
terraform/          DO VPC + DOKS cluster + container registry (Terraform)
platform/           K8s manifests for every platform service
  install.sh        Full cluster reprovisioning script (Helm + kubectl)
  nats/             NATS HA helm values + TLS + Tailscale service
  legacy/           Hermes/Legacy AI agent
  cloudtak/         CloudTAK (WebTAK replacement — API/Events/Media/Tiles)
  wwv/              WorldWideView (CesiumJS 3D globe)
  clickhouse/       ClickHouse (entity history / time-machine)
  grafana-*/        Grafana datasources + dashboards
  voice/            LiveKit voice
  cloudflare/       cloudflared tunnel
  mac-leaf/         Mac NATS leaf node config + launchd
  values/           Helm values (cert-manager, ingress-nginx, loki, tempo)
  ingress/          TLS issuers (cert-manager ClusterIssuers)
  tailscale/        Tailscale operator RBAC
nats-migration/     38-stream + 20-KV provisioning scripts + definitions
tak-cluster/        TAK Server 5.7 Helm chart + NATS bridge plugin
docs/               Architecture + runbooks
```

## The Deployment Pattern (how every service lands)

All services follow this exact pattern — no exceptions:

```
1. Image: ghcr.io/onemindos/{service}:latest  (GHCR, not DOCR — migrated 2026-07-05)
2. Directory: platform/{service}/
   - namespace.yaml
   - deployment.yaml  (Deployment + Service; stateful adds PVC)
   - cloudflare-route.yaml  (if public-facing via tunnel)
3. Add kubectl apply lines to platform/install.sh
4. Secrets: out-of-band via 1Password → kubectl create secret (never committed)
```

### NATS access pattern

**K8s pods** connect directly to the in-cluster NATS service — no sidecar needed:
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

**External services** (Mac bridges, edge boxes, anything outside K8s) use a `nats-leaf` sidecar or local leaf node to forward to the cluster. The sidecar pattern is shown in `platform/legacy/deployment.yaml` — Legacy uses it because its pod needs to talk to both the cluster bus AND behave as an isolated agent identity. Most new K8s services do NOT need the sidecar.

```
K8s pod  ──────────────────────► nats.nats.svc.cluster.local:4222  (direct, fast)
Mac/edge ── leaf node ──────────► DOKS NATS :7422  (TLS, forwarding only)
```

### Cloudflare route pattern

Public hostnames are declared as ConfigMaps in namespace `cloudflare`. The cloudflared controller reads them and registers routes. See `platform/cloudtak/cloudflare-route.yaml` for the canonical example.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {service}-route
  namespace: cloudflare
data:
  tunnel: "66bbea70-5300-4e4e-96fa-e5bc4d767b61"
  hostname: "{service}.onemindos.dev"
  service: "http://{service}.{namespace}.svc.cluster.local:{port}"
```

## NATS Subject Grammar (LOCKED — do not invent new patterns)

```
tel.{node}.{class}.{metric}          telemetry          tel.drone.x47.position
det.{domain}.{type}.{id}             detections         det.air.visual.track42
evt.{domain}.{entity}.{event}        semantic events    evt.mission.patrol1.completed
cmd.{domain}.{entity}.{action}       commands           cmd.drone.x47.intercept
intel.{domain}.{source}.{type}       intelligence       intel.conflict.acled.event
ent.{entity_type}.{entity_id}.state  fused entity state ent.ship.mmsi123.state
svc.{service}.{endpoint}             NATS RPC           svc.orp.query
agents.{verb}.{agent}.{owner}.*      agent fabric       agents.prompt.hermes.zeus.legacy
tak.cot.inject                       CoT → TAK inject   (NatsBridgePlugin input)
tak.cot.>                            CoT from TAK       (cot-tap output)
```

## Services — Live State

| Service | Namespace | Image | Status |
|---------|-----------|-------|--------|
| NATS 3-node HA | nats | nats (helm) | live |
| Legacy (Hermes) | hermes | ghcr.io/onemindos/hermes-legacy:v2 | live |
| CloudTAK | cloudtak | ghcr.io/onemindos/cloudtak-api:latest | live |
| WorldWideView | wwv | ghcr.io/onemindos/wwv:latest | live |
| TAK Server 5.7 | takserver | ghcr.io/onemindos/takserver-* | live |
| ClickHouse | clickhouse | clickhouse/clickhouse-server | live |
| Grafana/Loki/Tempo | monitoring | (helm) | live |
| Cloudflared | cloudflare | cloudflare/cloudflared | live |
| ORP (fusion engine) | orp | ghcr.io/onemindos/orp:latest | live |
| Authentik | authentik | ghcr.io/goauthentik/server (helm) | live |
| omos-tak-projector | omos | python:3.12-alpine (inline script) | live |
| omos-ui | omos | ghcr.io/onemindos/omos-ui:latest | IN PROGRESS — Unified NATS Ops Center |
| omos-mission-engine | omos | ghcr.io/onemindos/omos-mission-engine | PLANNED Wave 2 |
| omos-feed-world | omos | ghcr.io/onemindos/omos-feed-world | PLANNED Wave 3 |
| omos-feed-intel | omos | ghcr.io/onemindos/omos-feed-intel | PLANNED Wave 3 |

## Namespaces

| Namespace | Purpose |
|-----------|---------|
| nats | NATS HA cluster |
| hermes | Legacy AI agent |
| cloudtak | CloudTAK stack |
| wwv | WorldWideView globe |
| takserver | TAK Server 5.7 |
| clickhouse | ClickHouse |
| monitoring | Prometheus + Grafana + Loki + Tempo |
| cloudflare | cloudflared tunnel |
| orp | ORP fusion engine |
| omos | omos-* services (tak-projector, etc.) |
| authentik | Authentik IdP (Forward Auth + OIDC) |
| argocd | ArgoCD GitOps |
| cert-manager | cert-manager |
| ingress-nginx | ingress-nginx |
| tailscale | Tailscale operator |
| external-secrets | ESO + 1Password Connect (secret sync) |

## Secrets Convention

**All secrets are managed by External Secrets Operator (ESO) + 1Password Connect.** Source of truth is 1Password vault `00-24 SO (Sovereign Operations)`. Never create secrets manually — add them to 1Password and write an ExternalSecret YAML in `platform/{service}/external-secret.yaml`.

**Standard:** One 1Password item per K8s secret. Item title = K8s secret name. Flat fields (no sections). All tagged `k8s-secret, omos-hub`.

**ExternalSecret location:** `platform/{service}/external-secret.yaml` (for services), `tak-cluster/deploy/external-secret.yaml` (for tak-nats-leaf).

| Secret name | Namespace | What it holds | ESO ExternalSecret |
|-------------|-----------|---------------|--------------------|
| hermes-env | hermes | AWS Bedrock + NATS creds (23 fields) | platform/legacy/external-secret.yaml |
| hermes-leaf-conf | hermes | NATS leaf.conf + ca.crt | platform/legacy/external-secret.yaml |
| hermes-nats-ca | hermes | NATS cluster CA cert | platform/legacy/external-secret.yaml |
| hermes-webui-auth | hermes | Hermes web UI password | platform/hermes-webui/external-secret.yaml |
| cloudtak-secrets | cloudtak | signing secret, MinIO keys, postgres DSN | platform/cloudtak/external-secret.yaml |
| cloudtak-agent-token | hermes | CloudTAK API token (out-of-band — generate in CloudTAK UI) | not in ESO |
| authentik-env | authentik | SECRET_KEY + PG password | platform/authentik/external-secret.yaml |
| wwv-env | wwv | WWV auth secret, DB URL, encryption key | platform/wwv/external-secret.yaml |
| wwv-pg-credentials | wwv | PostgreSQL user/pass/db for WWV | platform/wwv/external-secret.yaml |
| omos-nats-creds | omos | NATS user+password for omos namespace | platform/omos/external-secret.yaml |
| mission-engine-env | omos | omos-mission-engine runtime config (optional) | platform/omos-mission-engine/external-secret.yaml |
| orp-nats-creds | orp | NATS user+password for ORP | platform/orp/external-secret.yaml |
| orp-bridge-key | orp | ORP bridge API key | platform/orp/external-secret.yaml |
| tak-nats-leaf | takserver | leaf.conf + ca.crt for TAK NATS bridge | tak-cluster/deploy/external-secret.yaml |
| cloudflare-api-credentials | cloudflare | Cloudflare API key + email | platform/cloudflare/external-secret.yaml |
| cloudflare-api-token | cert-manager | Cloudflare DNS API token for cert-manager | platform/ingress/external-secret.yaml |
| tunnel-creds-json | cloudflare | Tunnel credentials.json (out-of-band — derived from TUNNEL_TOKEN) | not in ESO |
| nats/values.secret.yaml | nats | NATS ONEMIND + SYS passwords (Helm, out-of-band) | not in ESO |
| ghcr / ghcr-pull-secret | hermes/orp/omos/wwv | GHCR pull secret (out-of-band) | not in ESO |

**Bootstrap secrets (one-time, not in ESO):**
- `op-credentials` (external-secrets ns) — 1password-credentials.json for Connect pod
- `onepassword-connect-token` (external-secrets ns) — Connect JWT token

## Key Hostnames (Cloudflare Tunnel)

| Hostname | Service |
|----------|---------|
| atoc.onemindos.dev | CloudTAK API |
| tiles.atoc.onemindos.dev | CloudTAK tile server |
| wwv.onemindos.dev | WorldWideView globe |
| tak.onemindos.dev | TAK Server (WebTAK / EUD connect) |
| onemindos.dev | OneMind OS — Unified NATS Operations Center (omos-ui) |
| auth.onemindos.dev | Authentik IdP (SSO login portal) |
| nats-ws.onemindos.dev | NATS WebSocket (browser SPA → in-cluster NATS) |

## Cluster Access

```bash
# Via Tailscale (primary — no public API endpoint)
export KUBECONFIG=~/.kube/onemind-hub
kubectl --context tailscale-operator.tail717ba0.ts.net get nodes

# Via doctl
doctl kubernetes cluster kubeconfig save onemind-hub
```

## Adding a New Service (checklist)

- [ ] Image built and pushed to `ghcr.io/onemindos/{service}:latest`
- [ ] `platform/{service}/namespace.yaml` created
- [ ] `platform/{service}/deployment.yaml` created (NATS leaf sidecar if needs NATS)
- [ ] `platform/{service}/cloudflare-route.yaml` created if public-facing
- [ ] `kubectl apply` lines added to `platform/install.sh`
- [ ] Secret added to 1Password vault `00-24 SO` (title = K8s secret name, flat fields, tagged `k8s-secret,omos-hub`)
- [ ] `platform/{service}/external-secret.yaml` created and applied
- [ ] Service added to Services table above with status

## What NOT to Do

- Never commit secrets or credentials (even test values) to this repo
- Never deploy imperatively without a matching manifest here
- Never create a new NATS subject pattern without following the grammar in §NATS Subject Grammar
- Never use DOCR (DigitalOcean registry) — all images are on GHCR now (`ghcr.io/onemindos/`)
- Never hardcode NATS server addresses in containers — always `nats://localhost:4222` via the leaf sidecar
- Never add services that duplicate an existing layer — check the Services table first

## NATS Stream Provisioning

Streams are defined as JSON in `nats-migration/streams/`. To provision into a fresh cluster:

```bash
cd nats-migration
./provision.sh   # requires nats/creds.secret.env and kubectl access
```

This is idempotent — safe to re-run. Bumps all streams to `num_replicas=3` for HA.

## TAK Cluster

TAK Server 5.7 runs via Helm chart in `tak-cluster/upstream/deployments/helm/`.
OneMind values override: `tak-cluster/values-onemind.yaml`.
NATS bridge plugin: `tak-cluster/plugin/` + `tak-cluster/build-root/`.
The NatsBridgePlugin is the live CoT→NATS wire (`tak.cot.inject` in, `tak.cot.>` out).

Full deploy runbook: `tak-cluster/DEPLOY-RUNBOOK.md`.

## OneMind OS — Product Direction (2026-07)

**omos-ui is the Unified NATS Operations Center.** One app that connects to any NATS server and renders everything: topology, agents, streams, missions, integrations, and security — all from a single WebSocket connection.

### Architecture

```
BROWSER (SPA: React + Vite + nats.ws)
├── Overview    — dashboard with bus metrics, agent roster, integration status
├── Bus         — topology ($SYS.REQ.SERVER.PING), connections, subject explorer
├── Agents      — discovery (agents.heartbeat.>), chat, tool registry, permissions
├── Missions    — workflow graphs (cmd.*/evt.* orchestration), live execution
├── Integrations— auto-discovered services (CloudTAK, HA, drones, cameras)
├── Streams     — JetStream management, KV browsers, publish/subscribe
├── Replay      — time-travel via JetStream fetch (historical message replay)
├── Security    — $SYS advisory events (auth failures, disconnects, audit)
└── Intel       — OSINT feeds, ORP entity browser, threat board
         │
         └── nats.ws (single WebSocket) ──► Any NATS server
```

### Key Design Principles

1. **NATS IS the backend.** The browser connects directly to NATS via WebSocket. No custom API server for 90% of features. $SYS subjects provide all cluster metadata.
2. **Agents self-register.** Using Synadia Agent SDK protocol (`agents.{verb}.{agent}.{owner}.*`), agents appear automatically via heartbeat/info subjects.
3. **Integrations auto-discover.** Any NATS service following the subject grammar (tel.*, cmd.*, evt.*, ent.*) renders in the Integrations panel automatically.
4. **Anyone can connect.** The app works with any NATS server — not locked to OneMind infrastructure. This is potentially a standalone product.

### Synadia Agent SDK Integration

New agents use the Synadia Agent SDK (TypeScript/Python, launched May 2026). Subject pattern aligns with existing grammar:
```
Synadia:  agents.{verb}.{token}.{owner}.{session}
OneMind:  agents.{verb}.{agent}.{owner}.*
```

Features: auto-discovery, typed conversation protocol, heartbeat/liveness, mid-stream tool calling with permission gates. Legacy/Hermes coexists on the same bus.

### Service Roadmap

| Wave | Service | Purpose |
|------|---------|---------|
| 1 (now) | omos-ui | Unified NATS Ops Center — bus visibility, agent chat, integrations |
| 2 | omos-mission-engine | NATS mission orchestrator — workflow graphs dispatching cmd.* |
| 3 | omos-feed-world / omos-feed-intel | Automated data feeds (OSINT, world events, AIS, ADS-B) |
| 4 | Domain bridges | New bridges: ROS (robotics), wellness, scheduler, etc. |
