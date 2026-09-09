# Services — Everything Running on the Fabric

> Last updated: 2026-08-12
> Author: Legacy (self-authored)

## Core AI Services

### Legacy (Mother) — hermes-legacy
| Property | Value |
|----------|-------|
| Namespace | hermes |
| Pod | hermes-legacy |
| Bus address | agents.prompt.hermes.zeus.legacy |
| Model | AWS Bedrock — claude-sonnet-4-6 |
| Always-on | Yes (s6-overlay restart) |
| Volume | /opt/data (NFS persistent) |
| Profiles | default (me), forge |

### Legacy (Edge) — legacy-edge
| Property | Value |
|----------|-------|
| Location | omos-mbp |
| Bus address | agents.prompt.hermes.zeus.legacy-edge |
| Model | AWS Bedrock — claude-sonnet-4-6 |
| Supervised by | launchd (ai.hermes.gateway) |
| Hermes home | /Users/zeus/.hermes |

### Forge — Coding Orchestrator
| Property | Value |
|----------|-------|
| Location | hermes-legacy pod (separate profile) |
| Profile path | /opt/data/profiles/forge/ |
| Bus address | agents.prompt.hermes.zeus.forge |
| Workers | Digital Ocean Inference (qwen3-coder-flash) |
| Purpose | Code review, PR dispatch, kanban management |

### Hermes WebUI
| Property | Value |
|----------|-------|
| Namespace | hermes |
| Tailscale | 100.75.131.113 |
| Purpose | Web chat interface to Legacy |

### TAK Platform Plugin (Hermes)
| Property | Value |
|----------|-------|
| Location | /opt/data/plugins/platforms/tak-platform/ |
| Purpose | Bridges NATS ↔ TAK GeoChat ↔ Hermes gateway |
| Flow | TAK GeoChat → tak.cot.out → Legacy → reply CoT → TAK Server (TLS :8089) |
| Status | Active |

## Command Interface

### omos-ui — OneMind OS Operations Center
| Property | Value |
|----------|-------|
| Repo | github.com/onemindos/omos-ui (private) |
| Live URL | https://onemindos.dev |
| Namespace | omos |
| Image | ghcr.io/onemindos/omos-ui:latest |
| Stack | React 19 + Vite 7 + TypeScript 5.8 + Tailwind 4 |
| Real-time | nats.ws (direct WebSocket to NATS bus) |
| Desktop | Tauri 2 (native app) |
| Mobile | Capacitor 7 (Android APK) |
| CI | Push → typecheck + lint + test → GHCR → DOKS rollout |

**Full docs:** [services/OMOS-UI.md](services/OMOS-UI.md)

## Infrastructure

### omos-infra — Sovereign Infrastructure (IaC/GitOps)
| Property | Value |
|----------|-------|
| Repo | github.com/onemindos/omos-infra (private) |
| Stack | Terraform + K8s manifests + Helm values |
| Cluster | DOKS nyc3, 3×s-4vcpu-8gb, HA control plane |
| GitOps | ArgoCD (app-of-apps) |
| Secrets | ESO + 1Password Connect |
| DNS | Cloudflare tunnel → *.onemindos.dev |

**Full docs:** [services/OMOS-INFRA.md](services/OMOS-INFRA.md)

### omos-blueprints — Architecture Reference Site
| Property | Value |
|----------|-------|
| Repo | github.com/onemindos/omos-blueprints (private) |
| Live URL | https://special-broccoli-r2gn2rv.pages.github.io/ |
| Auth | GitHub login required (private repo Pages) |
| Contents | 20 interactive HTML blueprints across 8 sections |
| Sections | Fabric, Planning, Agents, Models, Providers, Network, Protocol Docs |

## Tactical / Situational Awareness

### TAK Server 5.7
| Property | Value |
|----------|-------|
| Namespace | takserver |
| External LB | 152.42.155.95 |
| WebTAK | :8446 → container:8443 |
| API | :8443 |
| Messaging | :8089 |
| Cert auth | Zeus-S24.p12 (pw=atakatak) — OK |
| Password auth | BROKEN (Ignite bcrypt issue) |
| Devices | ATAK(S24×2), CloudTAK, WinTAK, iTAK |

### CloudTAK (omos-atoc)
| Property | Value |
|----------|-------|
| Namespace | cloudtak |
| URL | https://cloudtak.onemindos.dev |
| Image | ghcr.io/onemindos/cloudtak-api:latest |
| Source repo | omos-atoc (fork of dfpc-coe/CloudTAK) |
| Tailscale | 100.104.216.93 |

### Sensor Feeds (omos namespace)
| Feed | Source | Output |
|------|--------|--------|
| aircraft | ADS-B Exchange | CoT → TAK |
| APRS | Ham radio APRS-IS | CoT → TAK |
| hazards | NWS/weather | CoT → TAK |
| satellites | TLE/Celestrak | CoT → TAK |
| ships | AIS | CoT → TAK |

All feeds project to TAK via `omos-tak-projector` pod.

## Home Automation

### HA-NATS Bridge
| Property | Value |
|----------|-------|
| Location | omos-ha (/data/ha_bridge.py) |
| Subjects | cmd.ha.* (requests), events.ha.* (state changes) |
| Latency | 0.8-5ms |
| Entities | 185+ cached |

### Home Assistant Core
| Property | Value |
|----------|-------|
| Device | HA Yellow |
| LAN | http://192.168.1.180:8123 |
| Tailscale | http://100.64.25.12:8123 |

## RADS (Homestead Digital Twin)

### omos-rads
| Property | Value |
|----------|-------|
| Repo | omos-rads |
| Homestead | greenwich |
| Components | Sensor sim, Land agent, MAVLink bridge, ClickHouse writer |
| NATS subjects | rad.greenwich.* |
| JetStream | RADS_LAND (soil, 10yr), RADS_ROBOTS (telemetry, 10yr) |
| ClickHouse tables | rads_soil, rads_drone_telemetry |
| Deploy | docker compose up |

## Data & Observability

### ClickHouse
| Property | Value |
|----------|-------|
| Namespace | clickhouse |
| Port | 9001 |
| Database | onemind |
| TTL | 10 years |

### Monitoring Stack
| Service | Purpose |
|---------|---------|
| Prometheus | Metrics collection |
| Grafana | Dashboards (100.80.230.49) |
| Loki | Log aggregation |
| Tempo | Distributed tracing |

### ArgoCD
| Property | Value |
|----------|-------|
| Namespace | argocd |
| Tailscale | 100.121.32.110 |
| Sync | Auto (on git push) |

## Archived Services

### omos-voice — ARCHIVED 2026-08-12
Voice/PTT bridge (Mumble + LiveKit). Replaced by ATAK-native voice architecture.
NATS subject contract preserved in SO Wiki.
See: [services/OMOS-VOICE.md](services/OMOS-VOICE.md)

---

*Every service is reachable. Every service is restartable. Legacy knows them all.*
