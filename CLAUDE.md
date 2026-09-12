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
infra/              DO VPC + DOKS cluster + container registry (Terraform)
edge/               Non-K8s hardware-bound nodes
  mac-leaf/         Mac NATS leaf node config + launchd
services/           Standalone NATS microservices (presence, tasks, bridges)
tak-server/         TAK Server 5.8 (self-contained)
  chart/            Helm chart (upstream + OneMind values)
  docker/           Dockerfiles (base, api, messaging, config, plugins, db-setup)
  deploy/           K8s overlay manifests (LB, NATS leaf, secrets)
  plugin-src/       NatsBridgePlugin Java source
  plugin/           Built plugin JAR + config
  edge/             EUD management, device registry, simulators
  build/            Staging + CI helpers
platform/           K8s manifests — organized by operational layer
  argocd/           App-of-apps + Application CRDs
  core/             Wave 0: ingress, cloudflare, tailscale, sealed-secrets, cert-manager, github-runner
  nats/             Wave 1: NATS HA helm values + TLS + streams
  onetak/           Wave 2: OneTAK (API/Events/Media/Tiles) — renamed from cloudtak/
  hermes/           Wave 3: Hermes/Legacy AI agent
  omos/             Wave 3: tak-projector, feeds, mission-engine
  intel/            Wave 3: OSINT services
  geo/              Wave 3: geospatial services
  data/             Wave 3: clickhouse/, redpanda-connect/
  monitoring/       Wave 4: grafana, alertmanager, nats-surveyor
  nats-tools/       Wave 4: NATS debug UIs
  values/           Shared Helm values (cert-manager, ingress-nginx, loki, tempo)
docs/               Architecture + runbooks
```

## The Deployment Pattern (how every service lands)

All services follow this exact pattern — no exceptions:

```
1. Image: ghcr.io/onemindos/{service}:latest  (GHCR, not DOCR — migrated 2026-07-05)
2. Directory: platform/{service}/
   - namespace.yaml
   - deployment.yaml  (Deployment + Service; stateful adds PVC)
   - ingress.yaml     (if public-facing — nginx class, Cloudflare IP whitelist)
3. Secrets: out-of-band via 1Password → kubectl create secret (never committed)
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

**External services** (Mac bridges, edge boxes, anything outside K8s) use a `nats-leaf` sidecar or local leaf node to forward to the cluster. The sidecar pattern is shown in `platform/hermes/deployment.yaml` — Legacy uses it because its pod needs to talk to both the cluster bus AND behave as an isolated agent identity. Most new K8s services do NOT need the sidecar.

```
K8s pod  ──────────────────────► nats.nats.svc.cluster.local:4222  (direct, fast)
Mac/edge ── leaf node ──────────► DOKS NATS :7422  (TLS, forwarding only)
```

### Public ingress pattern

Public hostnames are handled by `ingress-nginx` (LB `143.244.220.155`). All public DNS records are Cloudflare-proxied A records pointing at the LB. The ingress object lives in the same namespace as the service and uses the Cloudflare IP whitelist so only Cloudflare-edge IPs can reach the LB.

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {service}-ingress
  namespace: {namespace}
  annotations:
    nginx.ingress.kubernetes.io/whitelist-source-range: "173.245.48.0/20,..."
spec:
  ingressClassName: nginx
  rules:
    - host: {service}.onemindos.dev
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: {service}
                port:
                  number: {port}
```

The Cloudflare tunnel (`cloudflared`) is **scaled to 0** and kept only as a cold standby. All routing is via ingress-nginx.

## NATS Subject Grammar (LOCKED — do not invent new patterns)

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
tak.cot.>                            CoT from TAK       (cot-tap output)
```

## Services — Live State

| Service | Namespace | Image | Status |
|---------|-----------|-------|--------|
| NATS 3-node HA | nats | nats (helm) | ✅ live |
| Legacy (Hermes AI) | hermes | ghcr.io/onemindos/hermes-legacy:latest | ✅ live |
| copilot-runtime | hermes | ghcr.io/onemindos/copilot-runtime:latest | ✅ live |
| OneTAK API | onetak | ghcr.io/onemindos/onetak-api:latest | ✅ live |
| OneTAK Events | onetak | ghcr.io/onemindos/onetak-events:latest | ✅ live |
| OneTAK Tiles | onetak | ghcr.io/onemindos/onetak-tiles:latest | ✅ live |
| OneTAK Media | onetak | ghcr.io/onemindos/onetak-media:latest | ✅ live |
| TAK Server 5.8 | takserver | ghcr.io/onemindos/takserver-base:* | ✅ live |
| ClickHouse | clickhouse | altinity/clickhouse-server (operator) | ✅ live |
| Grafana/Loki/Tempo | monitoring | (helm) | ✅ live |
| Cloudflared | cloudflare | cloudflare/cloudflared | 🔴 scaled to 0 (standby) |
| nats-dashboard | nats | ghcr.io/onemindos/nats-dashboard:latest | ✅ live |
| agent-web-ui | nats | ghcr.io/onemindos/agent-web-ui:latest | ✅ live |
| agenttap | nats | ghcr.io/onemindos/agenttap:latest | ✅ live |
| heynats | nats | ghcr.io/onemindos/heynats:latest | ✅ live |
| nats-ui | nats | ghcr.io/gastbob40/nats-ui:latest | ✅ live |
| omos-tak-projector | omos | python:3.12-alpine (inline script) | ✅ live |
| feed-intel | omos | ghcr.io/onemindos/feed-intel:latest | ✅ live |
| feed-satellites | omos | ghcr.io/onemindos/feed-satellites:latest | ✅ live |
| sink-clickhouse | omos | ghcr.io/onemindos/sink-clickhouse:latest | ✅ live |
| redpanda-connect | omos | ghcr.io/redpandadata/connect:latest | ✅ live |
| omos-mission-engine | omos | ghcr.io/onemindos/omos-mission-engine | 🔴 scaled to 0 |
| valhalla | onetak | valhalla/valhalla:run-latest | 🔴 scaled to 0 (awaiting tile build) |
| nominatim | onetak | mediagis/nominatim:4.4 | 🔴 scaled to 0 (geofabrik rate-limited) |
| fedhub | takserver | ghcr.io/onemindos/fedhub:5.8 | ⏳ manifests deployed; awaiting image build + certs |
| tak-ldap | takserver | ghcr.io/onemindos/tak-ldap:5.8 | ⏳ manifest deployed; awaiting image build + SealedSecret |
| vrtak-releases | vrtak | nginx:1.27-alpine | ⏳ file server live; awaiting Unreal build upload |

## Namespaces

| Namespace | Purpose |
|-----------|---------|
| nats | NATS 3-node HA cluster + tools |
| hermes | Legacy AI agent + copilot-runtime |
| onetak | OneTAK — API/Events/Media/Tiles + geo services |
| takserver | TAK Server 5.8 |
| clickhouse | ClickHouse (Altinity operator) |
| monitoring | Prometheus + Grafana + Loki + Tempo |
| cloudflare | cloudflared (scaled to 0 — cold standby only) |
| omos | omos-* services (tak-projector, feeds, sinks) |
| geo | Geospatial stack (Tile38, STAC, farmOS, etc.) |
| takml | TAK-ML vision pipeline |
| vrtak | VRTAK distribution server (nginx file server + Pixel Streaming future) |
| argocd | ArgoCD GitOps controller |
| cert-manager | TLS cert automation |
| tailscale | Tailscale proxy sidecars |
| github-runner | Self-hosted GH Actions runner |
| ingress-nginx | Ingress controller |
| kube-system | sealed-secrets controller (secret decryption) |

## Secrets Convention

**All secrets are managed by Sealed Secrets.** Source of truth is encrypted SealedSecret YAMLs committed to this repo. The sealed-secrets controller (kube-system) decrypts them using the cluster private key. ArgoCD applies SealedSecrets automatically on sync.

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

| Secret name | Namespace | What it holds |
|-------------|-----------|---------------|
| hermes-env | hermes | AWS Bedrock + NATS creds (23 fields) |
| hermes-leaf-conf | hermes | NATS leaf.conf + ca.crt |
| hermes-nats-ca | hermes | NATS cluster CA cert |
| hermes-webui-auth | hermes | Hermes web UI password |
| hermes-dashboard-auth | hermes | Hermes dashboard auth |
| cloudtak-agent-token | hermes | CloudTAK API token (out-of-band — generate in CloudTAK UI) |
| onetak-secrets | onetak | signing secret, MinIO keys, postgres DSN |
| omos-nats-creds | omos | NATS user+password for omos namespace |
| tak-nats-leaf | takserver | leaf.conf + ca.crt for TAK NATS bridge |
| cloudflare-api-credentials | cloudflare | Cloudflare API key + email |
| cloudflare-api-token | cert-manager | Cloudflare DNS API token for cert-manager |
| tunnel-creds-json | cloudflare | Tunnel credentials.json (derived from TUNNEL_TOKEN — one-time) |
| nats/values.secret.yaml | nats | NATS ONEMIND + SYS passwords (Helm — out-of-band) |
| ghcr / ghcr-pull-secret | hermes/omos | GHCR pull secret (out-of-band) |

## Public Hostnames

All public DNS records are **Cloudflare-proxied A records → `143.244.220.155`** (ingress-nginx LB), except `tak.onemindos.dev` which is an unproxied A record direct to the TAK LB.

| Hostname | Ingress object | Backend | DNS status |
|----------|---------------|---------|------------|
| onemindos.dev | onetak/onetak-ingress | onetak-api:5000 | ⚠️ update DNS |
| tiles.onemindos.dev | onetak/onetak-ingress | onetak-tiles:5002 | ⚠️ update DNS |
| media.onemindos.dev | onetak/onetak-ingress | onetak-media:9997 | ⚠️ update DNS |
| ws.onemindos.dev | nats/nats-ws-ingress | nats:8080 (WebSocket) | ⚠️ update DNS |
| tak.onemindos.dev | — | tak-external LB 152.42.155.95 | A record, unproxied |

**Tailscale-only (never public):**

| Hostname | Service |
|----------|---------|
| argocd.tail717ba0.ts.net | ArgoCD |
| grafana.tail717ba0.ts.net | Grafana |
| hermes-agent.tail717ba0.ts.net | Hermes web UI |
| agent-web-ui-1.tail717ba0.ts.net | Agent Web UI |
| nats-dashboard-1.tail717ba0.ts.net | NATS Dashboard |

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
- [ ] Manifests added to the appropriate platform layer directory (e.g. `platform/omos/`, `platform/data/`)
- [ ] Cloudflare route ConfigMap added if public-facing (namespace: cloudflare)
- [ ] SealedSecret created and committed (see Secrets Convention above)
- [ ] ArgoCD Application updated or created in `platform/argocd/` with correct path
- [ ] Service added to Services table above with status

## What NOT to Do

- Never commit secrets or credentials (even test values) to this repo
- Never deploy imperatively without a matching manifest here
- Never create a new NATS subject pattern without following the grammar in §NATS Subject Grammar
- Never use DOCR (DigitalOcean registry) — all images are on GHCR now (`ghcr.io/onemindos/`)
- Never hardcode NATS server addresses in containers — always `nats://localhost:4222` via the leaf sidecar
- Never add services that duplicate an existing layer — check the Services table first

## NATS Stream Provisioning

Streams are defined as YAML CRDs in `platform/nats/streams/` (39 files — Stream + KeyValue).
The NACK JetStream controller (installed via Helm in the `nats` namespace) reconciles them
continuously. To provision into a fresh cluster:

```bash
helm upgrade --install nack nats/nack -n nats
kubectl apply -f platform/nats/streams/
```

## TAK Server

TAK Server 5.8 (hardened) runs via Helm chart in `tak-server/chart/`.
OneMind values override: `tak-server/values-onemind.yaml`.
Dockerfiles: `tak-server/docker/` (base, api, messaging, config, plugins, db-setup).
NATS bridge plugin source: `tak-server/plugin-src/` → built JAR in `tak-server/plugin/`.
The NatsBridgePlugin is the live CoT→NATS wire (`tak.cot.inject` in, `tak.cot.>` out).
CI builds all images via `.github/workflows/build-takserver.yml` (self-hosted runner).

Full deploy runbook: `tak-server/DEPLOY-RUNBOOK.md`.

## Service Roadmap

| Wave | Service | Purpose |
|------|---------|---------|
| 2 | omos-mission-engine | NATS mission orchestrator — workflow graphs dispatching cmd.* |
| 3 | omos-feed-world / omos-feed-intel | Automated data feeds (OSINT, world events, AIS, ADS-B) |
| 4 | Domain bridges | New bridges: ROS (robotics), wellness, scheduler, etc. |
