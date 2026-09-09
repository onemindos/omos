# omos-infra — Sovereign Infrastructure

> Last updated: 2026-07-11
> Repo: github.com/onemindos/omos-infra (private)
> Default branch: main

## Purpose

The **GitOps source of truth** for everything deployed to the OneMind cloud hub. If a service runs in K8s, its manifest lives here. If infrastructure exists, Terraform declares it here. Nothing is deployed imperatively without a matching file in this repo.

## Stack

| Layer | Choice | Why |
|-------|--------|-----|
| IaC | Terraform | DO VPC + DOKS cluster + registry |
| Orchestration | DOKS (Kubernetes 1.36) | HA control plane, 3-5 nodes |
| GitOps | ArgoCD (app-of-apps) | Declarative sync from repo |
| Secrets | ESO + 1Password Connect | Zero manual secrets |
| Networking | Cloudflare tunnel + Tailscale | Public routes + private mesh |
| Monitoring | Grafana + Loki + Tempo + Prometheus | Full observability |
| TLS | cert-manager + Let's Encrypt | Auto-rotating certs |
| CI | GitHub Actions | Quality gate + image builds |
| Linting | yamllint + shellcheck + terraform fmt | Pre-commit + CI |

## Architecture

```
MAC (cockpit / hardware edge)
  MAVLink bridge    ─┐
  HA bridge         ─┤── NATS leaf ──► K8s NATS HA (3-node) ◄── all cloud services
  Vision/cameras    ─┘   (forwarding      THE BUS              ORP, Legacy, projectors,
                          pipe only)                            connectors, ClickHouse
```

**Bus authority:** The K8s NATS cluster IS the bus. The Mac is a forwarding pipe.

## Repo Layout

```
terraform/              DO VPC + DOKS + registry (Terraform)
platform/               K8s manifests for every service
  install.sh            Full cluster reprovisioning script
  argocd/               ArgoCD Applications (app-of-apps)
  nats/                 NATS HA helm values + TLS + 38 streams + 20 KV
  legacy/               Hermes/Legacy AI agent
  cloudtak/             CloudTAK (WebTAK replacement)
  wwv/                  WorldWideView (CesiumJS 3D globe)
  clickhouse/           ClickHouse (time-machine)
  orp/                  ORP fusion engine
  authentik/            Authentik IdP
  omos-ui/              OneMind OS command interface
  omos-mission-engine/  Mission planning + execution
  external-secrets/     ESO + 1Password Connect
  cloudflare/           Cloudflared tunnel
  mac-leaf/             Mac NATS leaf node config + launchd
  values/               Helm values (cert-manager, ingress, loki, tempo)
nats-migration/         Stream + KV provisioning scripts
tak-cluster/            TAK Server 5.7 Helm chart + NATS bridge
docs/                   Architecture + runbooks
Makefile                Developer commands
```

## Services Deployed

| Service | Namespace | Status |
|---------|-----------|--------|
| NATS 3-node HA | nats | ✅ live |
| Legacy (Hermes) | hermes | ✅ live |
| CloudTAK | cloudtak | ✅ live |
| WorldWideView | wwv | ✅ live |
| TAK Server 5.7 | takserver | ✅ live |
| ClickHouse | clickhouse | ✅ live |
| Grafana/Loki/Tempo | monitoring | ✅ live |
| Cloudflared | cloudflare | ✅ live |
| ORP (fusion) | orp | ✅ live |
| Authentik | authentik | ✅ live |
| omos-ui | omos | ✅ live |
| omos-tak-projector | omos | ✅ live |
| Mission Engine | omos | 🔜 planned |

## Infrastructure Spec

| Item | Value |
|------|-------|
| Provider | DigitalOcean |
| Region | nyc3 |
| K8s version | 1.36.0-do.2 |
| Control plane | HA (99.95% SLA) |
| Node pool | 3 × s-4vcpu-8gb, autoscale 3→5 |
| Registry | GHCR (migrated from DOCR 2026-07-05) |
| DNS | Cloudflare → DO Load Balancer |
| Domain | *.onemindos.dev |
| Cost | ~$190/mo (DO credit) |

## Deployment Pattern

Every service follows — no exceptions:
1. Image: `ghcr.io/onemindos/{service}:latest`
2. Directory: `platform/{service}/`
3. Files: namespace.yaml, deployment.yaml, external-secret.yaml, cloudflare-route.yaml (if public)
4. Add to `platform/install.sh`
5. Secrets: 1Password → ESO → K8s Secret (never committed)

## NATS Subject Grammar (LOCKED)

```
tel.{node}.{class}.{metric}          telemetry
det.{domain}.{type}.{id}             detections
evt.{domain}.{entity}.{event}        semantic events
cmd.{domain}.{entity}.{action}       commands
intel.{domain}.{source}.{type}       intelligence
ent.{entity_type}.{entity_id}.state  fused entity state
svc.{service}.{endpoint}             NATS RPC
agents.{verb}.{agent}.{owner}.*      agent fabric
tak.cot.inject                       CoT → TAK inject
tak.cot.>                            CoT from TAK
```

## CI Workflows

| Workflow | Purpose |
|----------|---------|
| quality.yml | yamllint + shellcheck + terraform validate + secret scan |
| build-cloudtak.yml | Build + push CloudTAK container |
| build-orp.yml | Build + push ORP container |
| build-synadia-bridge.yml | Build + push Synadia bridge |
| build-tak-plugins.yml | Build + push TAK plugins |
| tailscale.yml | Tailscale operator reconcile |

## Developer Workflow

```bash
make lint       # yamllint + terraform fmt + shellcheck
make fmt        # Auto-format Terraform
make plan       # Terraform plan
make validate   # Full quality gate (same as CI)
make install    # Deploy full platform to DOKS
make streams    # Provision NATS streams + KV
```

## Integration Points

- **Legacy (me):** I manage this repo, create PRs for infra changes
- **ArgoCD:** Syncs platform/ manifests to cluster
- **ESO:** Syncs 1Password secrets to K8s
- **GitHub Actions:** Builds containers, runs quality gate
- **Cloudflare:** Routes public traffic via tunnel
- **Tailscale:** Private mesh between all nodes
