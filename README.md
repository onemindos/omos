# omos-infra — OneMind Sovereign Infrastructure

> GitOps source of truth for everything deployed to the OneMind cloud hub.
> If a service runs in K8s, its manifest lives here. If infrastructure exists, Terraform declares it here.

## Quick Start

```bash
# Clone + enter
git clone git@github.com:onemindos/omos-infra.git && cd omos-infra

# Install pre-commit hooks (one-time)
pip install pre-commit && pre-commit install

# Show all developer commands
make help
```

## Architecture

```
MAC (cockpit / hardware edge)
  MAVLink bridge    ─┐
  HA bridge         ─┤── NATS leaf ──► K8s NATS HA (3-node) ◄── all cloud services
  Vision/cameras    ─┘   (forwarding      THE BUS              ORP, Legacy, projectors,
                          pipe only)                            connectors, ClickHouse
```

**The K8s NATS cluster is the authority.** The Mac is a forwarding pipe for hardware-bound bridges. If the Mac is offline, the cloud bus keeps running.

## Repo Layout

```
terraform/              DO VPC + DOKS cluster + container registry (Terraform)
platform/               K8s manifests for every service
  install.sh            Full cluster reprovisioning script
  argocd/               ArgoCD Applications (app-of-apps)
  nats/                 NATS HA helm values + TLS + streams/KV definitions
  legacy/               Hermes/Legacy AI agent
  cloudtak/             CloudTAK (WebTAK replacement)
  wwv/                  WorldWideView (CesiumJS 3D globe)
  clickhouse/           ClickHouse (time-machine)
  orp/                  ORP fusion engine
  authentik/            Authentik IdP (Forward Auth + OIDC)
  omos-ui/              OneMind OS command interface
  omos-mission-engine/  Mission planning + execution
  external-secrets/     ESO + 1Password Connect
  cloudflare/           Cloudflared tunnel
  mac-leaf/             Mac NATS leaf node config + launchd
  values/               Helm values (cert-manager, ingress-nginx, loki, tempo)
  ingress/              TLS issuers (cert-manager ClusterIssuers)
  tailscale/            Tailscale operator RBAC
nats-migration/         38-stream + 20-KV provisioning scripts
tak-cluster/            TAK Server 5.7 Helm chart + NATS bridge plugin
docs/                   Architecture + runbooks
uptime-kuma/            Uptime monitoring
Makefile                Developer commands (make help)
```

## Developer Workflow

```bash
make lint           # yamllint + terraform fmt + shellcheck
make fmt            # Auto-format Terraform
make plan           # Terraform plan (needs DIGITALOCEAN_TOKEN)
make install        # Deploy full platform to DOKS
make streams        # Provision NATS streams + KV
make validate       # Full quality gate (same as CI)
```

### Adding a New Service

Every service follows the same pattern — no exceptions:

```bash
# 1. Create service directory
mkdir -p platform/{service}

# 2. Required files
platform/{service}/
├── namespace.yaml        # Namespace declaration
├── deployment.yaml       # Deployment + Service (stateful adds PVC)
├── external-secret.yaml  # ESO → 1Password secret sync
└── cloudflare-route.yaml # If public-facing (via tunnel)

# 3. Add kubectl apply lines to platform/install.sh

# 4. Image convention: ghcr.io/onemindos/{service}:latest
```

### Secrets Convention

All secrets managed by **External Secrets Operator (ESO) + 1Password Connect**:
- Source of truth: 1Password vault `00-24 SO (Sovereign Operations)`
- Never `kubectl create secret` manually (except bootstrap)
- `external-secret.yaml` in each service directory declares what to sync

### NATS Access Pattern

```yaml
# K8s pods — direct connection (no sidecar needed)
env:
  - name: NATS_URL
    value: "nats://nats.nats.svc.cluster.local:4222"
  - name: NATS_USER
    valueFrom: { secretKeyRef: { name: nats-creds, key: user } }

# External/edge — NATS leaf node forwarding to :7422
```

## CI/CD

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `quality.yml` | All pushes + PRs | yamllint + shellcheck + terraform fmt/validate + secret scan |
| `build-cloudtak.yml` | CloudTAK changes | Build + push container |
| `build-orp.yml` | ORP changes | Build + push container |
| `build-synadia-bridge.yml` | Bridge changes | Build + push container |
| `build-tak-plugins.yml` | TAK plugin changes | Build + push container |
| `tailscale.yml` | Tailscale changes | Operator reconcile |

## Infrastructure Spec

| Item | Value |
|------|-------|
| Provider | DigitalOcean |
| Region | nyc3 |
| K8s version | 1.36.0-do.2 |
| Control plane | HA (99.95% SLA) |
| Node pool | 3 × s-4vcpu-8gb, autoscale 3→5 |
| Registry | GHCR (GitHub Container Registry) |
| DNS | Cloudflare → DO Load Balancer |
| Domain | *.onemindos.dev |
| Secrets | 1Password + ESO |
| Monitoring | Grafana + Loki + Tempo + Prometheus |
| GitOps | ArgoCD (app-of-apps pattern) |

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
| cloudflare | Cloudflared tunnel |
| orp | ORP fusion engine |
| omos | omos-* services |
| authentik | Authentik IdP |
| argocd | ArgoCD GitOps |
| cert-manager | TLS cert automation |
| ingress-nginx | Ingress controller |
| tailscale | Tailscale operator |
| external-secrets | ESO + 1Password Connect |

## Pre-commit Hooks

Installed via `.pre-commit-config.yaml`:
- Trailing whitespace + end-of-file fix
- YAML syntax check
- yamllint (style)
- ShellCheck (bash scripts)
- Terraform fmt + validate
- Private key detection

## Related Repos

| Repo | Purpose |
|------|---------|
| [omos-ui](https://github.com/onemindos/omos-ui) | Command interface (React SPA) |
| [omos-legacy](https://github.com/onemindos/omos-legacy) | Legacy's operational genome |
| [omos-flutter](https://github.com/onemindos/omos-flutter) | CloudTAK fork |
