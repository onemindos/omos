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
| omos-tak-projector | omos | python:3.12-alpine (inline script) | live |
| om-control-plane | omos | ghcr.io/onemindos/om-control-plane | PLANNED Wave 2 |
| wwv-bridge | omos | ghcr.io/onemindos/wwv-bridge | PLANNED Wave 2 |
| omos-feed-world | omos | ghcr.io/onemindos/omos-feed-world | PLANNED Wave 3 |
| omos-feed-intel | omos | ghcr.io/onemindos/omos-feed-intel | PLANNED Wave 3 |
| omos-mission | omos | ghcr.io/onemindos/omos-mission | PLANNED Wave 4 |

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
| argocd | ArgoCD GitOps |
| cert-manager | cert-manager |
| ingress-nginx | ingress-nginx |
| tailscale | Tailscale operator |

## Secrets Convention

All secrets are out-of-band (never committed). Source of truth: 1Password vault `00-24 SO`.

| Secret name | Namespace | What it holds | 1Password item |
|-------------|-----------|---------------|----------------|
| hermes-env | hermes | Full ~/.hermes/.env (AWS Bedrock + NATS creds) | hermes-env |
| hermes-leaf-conf | hermes | NATS leaf.conf + ca.crt for hermes pod | omos-hub-leaf |
| hermes-nats-ca | hermes | NATS cluster CA cert | nats-ca |
| nats/values.secret.yaml | nats | NATS ONEMIND + SYS passwords | omos-nats-secrets |
| tak-nats-leaf | takserver | leaf.conf + ca.crt for TAK NATS bridge | tak-leaf-creds |
| cloudtak-secrets | cloudtak | signing secret, MinIO, postgres DSN | cloudtak-secrets |
| wwv-env | wwv | WWV environment (auth, DB, etc.) | wwv-env |
| ghcr | hermes/orp/omos | GHCR pull secret | ghcr-onemindos |
| ghcr-pull-secret | wwv | GHCR pull secret | ghcr-onemindos |
| omos-nats-creds | omos | NATS user+password for omos namespace workloads | omos-nats-secrets |

## Key Hostnames (Cloudflare Tunnel)

| Hostname | Service |
|----------|---------|
| cloudtak.onemindos.dev | CloudTAK API |
| tiles.cloudtak.onemindos.dev | CloudTAK tile server |
| wwv.onemindos.dev | WorldWideView globe |
| tak.onemindos.dev | TAK Server (WebTAK / EUD connect) |

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
- [ ] Secrets documented in table above, values in 1Password
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
