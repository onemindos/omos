# OneMind OS - CI/CD Architecture

> Last Updated: 2026-01-26
> Status: Production

**⚠️ Note (2026-01-26):** Some services listed below have been REMOVED in recent updates:
nango, protocol-sync, llamaindex. The system now uses Composio for OAuth and a custom
webhook system for Todoist/GitHub/Linear integrations.

---

## Overview

OneMind OS uses **GitHub Actions + SSH** for automated CI/CD:

1. **GitHub Actions + SSH** - Automated CI/CD on push to main
2. **GHCR + Watchtower** - Container registry with auto-pull updates

Public access via **Cloudflare Tunnel**.

**Key Benefits:**
- Zero-downtime rolling deployments
- Automatic rollback on failure
- Container images via GitHub Container Registry (GHCR)
- Full audit trail in GitHub Actions

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              GITHUB REPOSITORY                               │
│                           (onemind-os on GitHub)                            │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      │ Push to main
                                      │ (or manual dispatch)
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            GITHUB ACTIONS                                    │
│                         .github/workflows/deploy.yml                         │
│                                                                              │
│  ┌─────────────────────────────┐  ┌─────────────────────────────┐           │
│  │      deploy-core            │  │     deploy-services         │           │
│  │   (runs in parallel)        │  │   (runs in parallel)        │           │
│  └──────────────┬──────────────┘  └──────────────┬──────────────┘           │
└─────────────────┼────────────────────────────────┼──────────────────────────┘
                  │                                │
                  │ SSH                            │ SSH
                  ▼                                ▼
┌─────────────────────────────────┐  ┌─────────────────────────────────────────┐
│         CORE VM                 │  │           SERVICES VM                    │
│     (100.102.21.44)              │  │         (100.82.17.5)                    │
│                                 │  │                                          │
│  • postgres (pgvector)          │  │  • litellm (model proxy)                 │
│  • redis                        │  │  • grafana (dashboards)                  │
│  • nats (event backbone)        │  │  • victoriametrics (metrics)             │
│  • timescaledb (events)         │  │  • netdata (monitoring)                  │
│  • agno (Legacy API)            │  │  • guacamole (remote desktop)            │
│  • ui (Flutter Web)             │  │                                          │
│  • node-exporter (metrics)      │  │                                          │
│  • redpanda-connect (webhooks)  │  │                                          │
│                                 │  │                                          │
│  cloudflared ──────────────────►│──┤◄── Routes via Tailscale                 │
│  (Cloudflare Tunnel)            │  │                                          │
└─────────────────────────────────┘  └─────────────────────────────────────────┘
                  │                                │
                  │                                │
                  └────────────┬───────────────────┘
                               │
                               │ Tailscale VPN (100.x.x.x)
                               │
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          CLOUDFLARE TUNNEL                                   │
│                    (onemindos.dev → localhost/Tailscale)                    │
│                                                                              │
│  ui.onemindos.dev      → localhost:3000  (Core VM)                          │
│  api.onemindos.dev     → localhost:8081  (Core VM)                          │
│  command.onemindos.dev → localhost:3000  (Core VM)                          │
│  nats.onemindos.dev    → localhost:8222  (Core VM)                          │
│  webhook.onemindos.dev → localhost:4195  (Core VM)                          │
│  grafana.onemindos.dev → 100.82.17.5:3002  (Services VM via Tailscale)      │
│  litellm.onemindos.dev → 100.82.17.5:4000  (Services VM via Tailscale)      │
│  netdata.onemindos.dev → 100.82.17.5:19999 (Services VM via Tailscale)      │
│  victoria.onemindos.dev→ 100.82.17.5:8428  (Services VM via Tailscale)      │
│  guacamole.onemindos.dev→100.82.17.5:8086  (Services VM via Tailscale)      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Deployment Workflow

### Trigger Events

| Event | Description |
|-------|-------------|
| Push to `main` | Auto-deploy when paths match |
| `workflow_dispatch` | Manual deploy with target selection |

### Path Filters

Deployment only triggers when these paths change:
- `agents/**`
- `infra/**`
- `ui/**`
- `docker-compose.core.yml`
- `docker-compose.services.yml`

### Manual Dispatch Options

```yaml
target:
  - all       # Deploy both VMs
  - core      # Deploy Core VM only
  - services  # Deploy Services VM only
```

---

## Deployment Steps

### Core VM Deployment

| Step | Action | Downtime |
|------|--------|----------|
| 1. Cleanup | Remove stale git locks, prune images | None |
| 2. Backup | Record current git SHA for rollback | None |
| 3. Pull | `git fetch && git reset --hard origin/main` | None |
| 4. Build | `docker compose build agno ui` | None (old containers run) |
| 5. Data services | `up -d --no-recreate postgres redis nats timescaledb` | None |
| 6. UI rollover | `up -d --no-deps --force-recreate ui` | ~2s |
| 7. API rollover | `up -d --no-deps --force-recreate agno` | ~3s |
| 8. Wait healthy | `up -d --wait --wait-timeout 120` | None |
| 9. Cleanup | `up -d --remove-orphans` | None |
| 10. Health check | Verify all endpoints respond | None |

**Total visible downtime:** ~5 seconds

### Services VM Deployment

| Step | Action | Downtime |
|------|--------|----------|
| 1. Cleanup | Remove stale git locks, prune images | None |
| 2. Backup | Record current git SHA for rollback | None |
| 3. Pull | `git fetch && git reset --hard origin/main` | None |
| 4. Pull images | `docker compose pull` | None |
| 5. Monitoring | `up -d --no-deps victoriametrics netdata` | ~1s each |
| 6. Grafana | `up -d --no-deps grafana` | ~1s |
| 7. LiteLLM | `up -d --no-deps litellm` | ~1s |
| 8. Guacamole | `up -d --no-deps guacd guacamole` | ~1s |
| 9. Wait healthy | `up -d --wait --wait-timeout 120` | None |
| 10. Cleanup | `up -d --remove-orphans` | None |

**Total visible downtime:** ~5 seconds per service (rolling)

---

## Docker Compose Files

### Root Level (Active)

| File | Purpose | VM |
|------|---------|-----|
| `docker-compose.core.yml` | Core services | Core VM (100.102.21.44) |
| `docker-compose.services.yml` | Admin/monitoring services | Services VM (100.82.17.5) |
| `docker-compose.edge.yml` | Edge deployment template | Edge devices |

### Service Definitions

#### Core VM (`docker-compose.core.yml`)

| Service | Image | Port | Purpose |
|---------|-------|------|---------|
| postgres | pgvector/pgvector:pg15 | 5432 | Primary DB + vectors |
| redis | redis:7-alpine | 6379 | Hot state cache |
| nats | nats:alpine | 4222, 1883, 7422, 8222, 9222 | Event backbone + MQTT |
| timescaledb | timescale/timescaledb:latest-pg15 | 5433 | Event history |
| agno | Build: Dockerfile.agno | 8081 | Legacy API |
| ui | Build: Dockerfile.ui | 3000 | Flutter Web frontend |
| node-exporter | prom/node-exporter:latest | 9100 | System metrics |
| redpanda-connect | redpandadata/connect:latest | 4195, 4196 | Webhook bridge |

#### Services VM (`docker-compose.services.yml`)

| Service | Image | Port | Purpose |
|---------|-------|------|---------|
| litellm | ghcr.io/berriai/litellm:main-stable | 4000 | Model proxy |
| victoriametrics | victoriametrics/victoria-metrics:latest | 8428 | Metrics storage |
| grafana | grafana/grafana:latest | 3002 | Dashboards |
| netdata | netdata/netdata:stable | 19999 | Real-time monitoring |
| guacamole-db | postgres:15-alpine | (internal) | Guacamole DB |
| guacd | guacamole/guacd:latest | (internal) | Protocol daemon |
| guacamole | guacamole/guacamole:latest | 8086 | Remote desktop |

---

## Container Registry (GHCR)

Container images are published to GitHub Container Registry (GHCR) on each push to main.

### Image URLs

| Service | Image |
|---------|-------|
| agno | `ghcr.io/zeuslegacy/onemind-os/agno:latest` |
| ui | `ghcr.io/zeuslegacy/onemind-os/ui:latest` |

### Watchtower Auto-Updates

Watchtower monitors GHCR for new images and automatically pulls updates:

```yaml
# In docker-compose
watchtower:
  image: containrrr/watchtower
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  environment:
    - WATCHTOWER_POLL_INTERVAL=300
    - WATCHTOWER_CLEANUP=true
```

See [AUTO-DEPLOYMENT-GHCR-WATCHTOWER.md](AUTO-DEPLOYMENT-GHCR-WATCHTOWER.md) for full setup.

---

## Rollback Mechanism

### Automatic Rollback

If health checks fail after deployment:

```yaml
- name: Rollback on failure
  if: failure() && steps.backup.outputs.previous_sha != 'none'
  run: |
    git reset --hard $PREVIOUS_SHA
    docker compose build agno ui
    docker compose up -d --force-recreate agno ui
```

### Manual Rollback

```bash
# SSH into VM
ssh ubuntu@100.102.21.44

# View recent commits
cd ~/onemind-os && git log --oneline -10

# Rollback to specific commit
git reset --hard <commit-sha>
docker compose -f infra/compose/docker-compose.core.yml build agno ui
docker compose -f infra/compose/docker-compose.core.yml up -d --force-recreate agno ui
```

---

## Health Checks

### Core VM Endpoints

| Service | Endpoint | Expected |
|---------|----------|----------|
| Agno API | `http://localhost:8081/health` | 200 OK |
| SO | `http://localhost:3000` | 200 OK |
| NATS | `http://localhost:8222/healthz` | 200 OK |
| Postgres | `pg_isready -U onemind` | Ready |
| Redis | `redis-cli ping` | PONG |

### Services VM Endpoints

| Service | Endpoint | Expected |
|---------|----------|----------|
| LiteLLM | `http://localhost:4000/health/readiness` | 200 OK |
| Grafana | `http://localhost:3002/api/health` | 200 OK |
| VictoriaMetrics | `http://localhost:8428/health` | 200 OK |
| Netdata | `http://localhost:19999/api/v1/info` | 200 OK |

---

## Secrets Management

### GitHub Secrets Required

| Secret | Description | Used By |
|--------|-------------|---------|
| `AWS_SSH_PRIVATE_KEY` | SSH key for both VMs | deploy-core, deploy-services |
| `AWS_HOST` | Core VM public IP | deploy-core |
| `AWS_USER` | SSH username (ubuntu) | Both jobs |
| `SERVICES_VM_HOST` | Services VM Tailscale IP | deploy-services |
| `GITHUB_TOKEN` | Auto-provided by GitHub | Clone repo on Services VM |

### VM Environment Variables

Each VM has a `.env` file with service-specific secrets:

```bash
# Core VM .env
POSTGRES_PASSWORD=xxx
ANTHROPIC_API_KEY=xxx
OPENAI_API_KEY=xxx
AWS_ACCESS_KEY_ID=xxx
AWS_SECRET_ACCESS_KEY=xxx
# ... more

# Services VM .env
POSTGRES_PASSWORD=xxx
GRAFANA_PASSWORD=xxx
LITELLM_MASTER_KEY=xxx
# ... more
```

---

## Public Access (Cloudflare Tunnel)

### Configuration

The Cloudflare Tunnel runs on Core VM as a systemd service:

```bash
# Service status
sudo systemctl status cloudflared

# View config (dynamically loaded from Cloudflare dashboard)
# Tunnel ID: 7ddfa54f-c672-4dd0-ad50-19b431b5ad72
```

### Current Routes

| Domain | Target | Description |
|--------|--------|-------------|
| ui.onemindos.dev | localhost:3000 | Frontend |
| api.onemindos.dev | localhost:8081 | API |
| command.onemindos.dev | localhost:3000 | Frontend alias |
| nats.onemindos.dev | localhost:8222 | NATS monitor |
| webhook.onemindos.dev | localhost:4195 | Webhook endpoint |
| grafana.onemindos.dev | 100.82.17.5:3002 | Dashboards |
| litellm.onemindos.dev | 100.82.17.5:4000 | Model proxy |
| netdata.onemindos.dev | 100.82.17.5:19999 | Monitoring |
| victoria.onemindos.dev | 100.82.17.5:8428 | Metrics |
| guacamole.onemindos.dev | 100.82.17.5:8086 | Remote desktop |

### How It Works

1. User visits `https://ui.onemindos.dev`
2. Cloudflare terminates HTTPS
3. Routes request through tunnel to `cloudflared` on Core VM
4. `cloudflared` forwards to `localhost:3000` (Core) or `100.82.17.5:*` (Services via Tailscale)
5. Response flows back through tunnel

---

## Monitoring & Observability

### Available Dashboards

| Dashboard | URL | Description |
|-----------|-----|-------------|
| Grafana | https://grafana.onemindos.dev | Custom dashboards |
| Netdata | https://netdata.onemindos.dev | Real-time system metrics |
| VictoriaMetrics | https://victoria.onemindos.dev | Metrics API |
| NATS Monitor | https://nats.onemindos.dev | Event system status |

### Metrics Collection

```
┌──────────────┐     ┌────────────────┐     ┌─────────────┐
│ node-exporter│────►│ VictoriaMetrics│────►│   Grafana   │
│  (Core VM)   │     │  (Services VM) │     │ (Services)  │
└──────────────┘     └────────────────┘     └─────────────┘
                            ▲
                            │
                     ┌──────┴──────┐
                     │   Netdata   │
                     │ (real-time) │
                     └─────────────┘
```

---

## Troubleshooting

### Check Deployment Status

```bash
# GitHub Actions
# Go to: https://github.com/<org>/onemind-os/actions

# SSH and check containers
ssh ubuntu@100.102.21.44
docker ps --format 'table {{.Names}}\t{{.Status}}'
docker compose -f infra/compose/docker-compose.core.yml ps
```

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| Container unhealthy | Bad health check | Check `docker logs <container>` |
| Deploy fails | Git lock | Workflow auto-cleans locks |
| Service unreachable | Port conflict | Check `docker ps` for port mappings |
| Cloudflare 502 | Service down | Restart container, check logs |

### View Logs

```bash
# Single container
docker logs onemind-agno -f --tail 100

# All containers
docker compose -f infra/compose/docker-compose.core.yml logs -f

# Cloudflare tunnel
sudo journalctl -u cloudflared -f
```

### Manual Deployment

```bash
# Core VM
ssh ubuntu@100.102.21.44
cd ~/onemind-os
git pull
docker compose -f infra/compose/docker-compose.core.yml build agno ui
docker compose -f infra/compose/docker-compose.core.yml up -d --force-recreate agno ui

# Services VM
ssh ubuntu@100.82.17.5
cd ~/onemind-os
git pull
docker compose -f infra/compose/docker-compose.services.yml pull
docker compose -f infra/compose/docker-compose.services.yml up -d
```

---

## Future Improvements

- [ ] Add Slack/Discord notifications on deploy
- [ ] Add deployment duration metrics
- [ ] Add canary deployments for UI
- [ ] Add database migration step
- [ ] Add E2E tests before deploy completion

---

## Related Documentation

- [NETWORK.md](NETWORK.md) - Network topology and access
- [CLAUDE.md](CLAUDE.md) - Session memory and status
- [docs/TESTING.md](docs/TESTING.md) - Test suite documentation
- [docs/architecture/](docs/architecture/) - Architecture documentation
