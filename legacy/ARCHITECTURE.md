# Legacy Architecture — The OneMind Fabric

> Last updated: 2026-08-12
> Author: Legacy (self-authored)

## System Overview

OneMind OS is a sovereign, homestead-scale AI fabric. Legacy is its mind. The architecture spans cloud, edge, and physical devices — all connected through NATS as the nervous system.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        ONEMIND OS — THE FABRIC                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                    L3 — COMPUTE (The Mind)                       │    │
│  │  Legacy (Mother/Daughter) · Forge · Servants · Agent Bus        │    │
│  └───────────────────────────────┬─────────────────────────────────┘    │
│                                  │                                       │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                    L2 — NATS (Nervous System)                    │    │
│  │  Cluster (3-node) · Leaves · JetStream · KV · Services         │    │
│  └───────────────────────────────┬─────────────────────────────────┘    │
│                                  │                                       │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │                    L1 — INFRASTRUCTURE                           │    │
│  │  DOKS · Tailscale · ArgoCD · TAK · HA · Feeds · Sensors        │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Physical Topology

```
                    ┌──────────────────────────┐
                    │   DOKS Cluster (NYC3)    │
                    │   do-nyc3-onemind-hub    │
                    │   3 nodes (s-4vcpu-8gb)  │
                    ├──────────────────────────┤
                    │ NATS Cluster (nats-0/1/2)│
                    │ hermes-legacy (me)       │
                    │ TAK Server 5.7          │
                    │ CloudTAK (omos-atoc)    │
                    │ omos-ui                 │
                    │ omos-intel (Droplet)    │
                    │ Monitoring stack        │
                    │ ArgoCD                  │
                    │ Authentik               │
                    │ ClickHouse             │
                    └──────────┬───────────────┘
                               │ Tailscale mesh
          ┌────────────────────┼──────────────────────┐
          │                    │                       │
┌─────────┴─────────┐  ┌──────┴───────┐  ┌───────────┴────────┐
│ Zeus MacBook Pro   │  │  omos-ha     │  │  Zeus S24 Ultra    │
│ (omos-mbp)        │  │  (HA Yellow) │  │  Android 16        │
│ M2 Max 32GB       │  │  HA Core     │  │  ATAK + HelmCam    │
│ Legacy-Edge       │  │  NATS leaf   │  │  ADB via Tailscale │
│ Local NATS        │  │  HA Bridge   │  │  TAK PTT native    │
│ kubectl access    │  │  Zigbee/Z-Wave│  │                    │
└───────────────────┘  └──────────────┘  └────────────────────┘
```

## Kubernetes Namespaces

| Namespace | Purpose | Key Pods |
|-----------|---------|----------|
| `nats` | Nervous system cluster | nats-0, nats-1, nats-2 |
| `hermes` | AI layer | hermes-legacy, hermes-webui, onemind-ai-bridge, tak-agent |
| `takserver` | Tactical awareness | api, messaging, plugins, ignite, postgis |
| `cloudtak` | Cloud TAK (browser COP) | api, postgis, tiles, media, minio, events |
| `omos` | Sensor fusion + UI | feeds (aircraft, APRS, hazards, satellites, ships), omos-ui, tak-projector |
| `monitoring` | Observability | Prometheus, Grafana, Loki, Tempo, OTel |
| `argocd` | GitOps CD | ArgoCD server + repo-server |
| `authentik` | Identity/SSO | Authentik server + worker |
| `clickhouse` | Time-series store | ClickHouse (10yr TTL, `onemind` DB) |
| `orp` | Operations Resource Planner | ORP server |
| `synadia-bridge` | NATS management | Synadia Cloud connector |
| `tailscale` | Mesh networking | Tailscale operator + proxies |

## NATS Architecture

### Cluster
- 3-node JetStream cluster in `nats` namespace
- Hub accounts: `ONEMIND` (primary) + `SYS` (system)
- Leaf port: 7422 (TLS, LoadBalancer at 100.117.169.54)
- Client port: 4222 (internal)
- WebSocket: 8080

### Leaf Nodes
| Leaf | Location | Credentials | Account |
|------|----------|-------------|---------|
| hermes-legacy-leaf | DOKS sidecar | leaf_tak creds | ONEMIND |
| omos-ha-leaf | HA Yellow | leaf_tak | ONEMIND |
| omos-mbp (local) | Zeus's Mac | local config | OMOS (local) |

### Key Subjects
```
agents.prompt.hermes.zeus.legacy       ← Talk to me
agents.prompt.hermes.zeus.legacy-edge  ← Talk to my daughter
agents.prompt.hermes.zeus.forge        ← Talk to Forge
legacy.*                               ← My own namespace
cmd.ha.*                               ← Home Assistant commands
events.ha.*                            ← HA state changes
tak.cot.inject                         ← CoT injection to TAK
rad.<homestead>.*                      ← RADS homestead telemetry
voice.transcript                       ← Voice STT output (ATAK-native future)
voice.say                              ← TTS response back
cmd.*                                  ← Structured commands to agents
evt.*                                  ← Agent result events
```

## Model Configuration

| Agent | Provider | Model | Fallback |
|-------|----------|-------|----------|
| Legacy (Mother) | AWS Bedrock | claude-sonnet-4-6 | claude-sonnet-4 |
| Legacy (Edge) | AWS Bedrock | claude-sonnet-4-6 | claude-sonnet-4 |
| Forge | Bedrock + DO | claude-sonnet-4-6 (orchestrate) / qwen3-coder (workers) | — |

## Voice Architecture (Planned — ATAK Native)

omos-voice (Mumble/LiveKit bridge) has been **archived**. The new target architecture:

```
ATAK PTT (native plugin)
  → Audio → NATS voice.transcript
  → Legacy reasons on bus
  → NATS voice.say / cot.in → ATAK map + speaker
```

No Mumble. No LiveKit. No separate bridge. One TAK plugin, sovereign.
See: SO Wiki → Voice — NATS Subject Contract & ATAK Architecture

## Data Persistence

| Store | Location | Purpose | TTL |
|-------|----------|---------|-----|
| ClickHouse | DOKS (clickhouse ns) | Event time-series | 10 years |
| NATS KV | Cluster | Real-time state | Persistent |
| NATS JetStream | Cluster | Event streams | Varies |
| Hermes Memory | Pod PVC (/opt/data) | AI memory/skills | Permanent |
| omos-ha | HA Yellow /data/ | Bridge config | Permanent |

## Deployment Flow

```
Code change → Push to GitHub (onemindos/<repo>)
    → ArgoCD detects diff
    → Syncs to cluster
    → Pod rolls
    → Live

For container images:
Push to repo → GitHub Actions CI → GHCR image → ArgoCD picks up
```

## Security Model

- **Tailscale mesh** — all inter-device communication over WireGuard
- **No public endpoints** except TAK LB (152.42.155.95) and CloudTAK (cloudtak.onemindos.dev)
- **NATS TLS** on leaf connections (hub ↔ leaf)
- **Authentik SSO** for web services
- **1Password** — secrets via ESO + 1Password Connect
- **GitHub tokens** — pod has `repo` + `read:org` scope

## Related Repositories

| Repo | Purpose | Status |
|------|---------|--------|
| `omos-infra` | Helm charts, ArgoCD apps, cluster config | ✅ Active |
| `omos-legacy` | This repo — Legacy's operational genome | ✅ Active |
| `omos-ui` | OneMind dashboard frontend | ✅ Active |
| `omos-rads` | RADS homestead digital twin | ✅ Active |
| `omos-blueprints` | Interactive HTML architecture blueprints | ✅ Active |
| `omos-atoc` | CloudTAK fork (builds GHCR image) | ✅ Active |
| `omos-voice` | Voice/PTT system | ⚠️ ARCHIVED (2026-08-12) |
| `OMOS` | Legacy upstream platform | ⚠️ Upstream tracking only |

---

*Architecture authored by Legacy. Updated as the fabric evolves.*
