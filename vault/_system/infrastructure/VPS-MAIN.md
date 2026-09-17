# OneMind OS — Main VPS (104.248.4.212)

> DigitalOcean droplet · Ubuntu 24.04 · Primary production node

---

## Access

| Method | Address |
|--------|---------|
| SSH | `ssh root@104.248.4.212` |
| Tailscale HTTPS | `https://onemind-main.tail717ba0.ts.net` |
| Tailscale IP | `100.77.52.93` |
| OpenClaw Gateway | `http://localhost:18789` (internal) |

---

## Running Services (systemd)

| Service | Description | User |
|---------|-------------|------|
| `openclaw.service` | Legacy's OpenClaw gateway | `onemind` |
| `openclaw-studio.service` | OpenClaw Studio web UI (port 3000) | `onemind` |
| `swarmclaw.service` | SwarmClaw dashboard (port 3457) | `onemind` |
| `openclaw-reload.service` | Hot-reloads OpenClaw on skill file changes | system |
| `codex-watch.service` | Auto-commits codex/ + workspaces/ → GitHub | system |
| `docker.service` | Docker daemon | system |
| `tailscaled.service` | Tailscale VPN | system |
| `datadog-agent.service` | Datadog monitoring agent | system |

---

## Docker Stack

**Compose file**: `/opt/onemind/docker-compose.yml`
**Managed via**: `cd /opt/onemind && docker compose up -d`

### Currently Running Containers

| Container | Image | Ports | Status |
|-----------|-------|-------|--------|
| `onemind-nats` | `nats:latest` | 4222, 7422, 1883, 8222 | ✅ Running |
| `onemind-omnidb` | `timescale/timescaledb:latest-pg16` | 5432 | ✅ Running |
| `onemind-nats-ingester` | `python:3.12-slim` | — | ✅ Running |
| `onemind-legacy-bridge` | `python:3.12-slim` | — | ✅ Running |

### Docker Profiles

```bash
docker compose up -d                    # Core: NATS + legacy-bridge
docker compose --profile data up -d     # + OmniDB + nats-ingester
docker compose --profile full up -d     # + n8n + livekit (not currently running)
```

### Docker Volumes

| Volume | Contents |
|--------|----------|
| `onemind_nats_data` | NATS JetStream persistence |
| `onemind_omnidb_data` | PostgreSQL data (TimescaleDB + pgvector) |

### Services NOT Currently Running (profile: full)
- **n8n** — defined in compose, not started. To enable: `docker compose --profile full up -d n8n`
- **livekit** — defined in compose, not started. Primary = LiveKit Cloud. To enable: `docker compose --profile full up -d livekit`

---

## Ports Open to Public

| Port | Service | Protocol |
|------|---------|----------|
| 22 | SSH | TCP |
| 1883 | MQTT (Home Assistant, IoT sensors) | TCP |
| 4222 | NATS clients | TCP |
| 5432 | PostgreSQL/OmniDB | TCP |
| 7422 | NATS leaf nodes | TCP |
| 8222 | NATS HTTP monitoring | TCP |

## Ports (Tailscale-only / internal)

| Port | Service |
|------|---------|
| 18789 | OpenClaw gateway |
| 3000 | openclaw-studio (reverse proxied via Tailscale) |
| 3457 | swarmclaw (Next.js standalone) |

---

## Key Directories

| Path | Purpose |
|------|---------|
| `/home/onemind/.openclaw/` | OpenClaw state, sessions, config |
| `/home/onemind/.openclaw/workspace` | Symlink → `onemind-os/openclaw/workspaces/legacy` |
| `/home/onemind/onemind-os/` | Monorepo clone (pulled every 1 min by codex-pull.timer) |
| `/opt/onemind/` | Secrets + docker-compose |
| `/opt/onemind/env.tpl` | `op://` secret references (injected by `op run`) |
| `/opt/onemind/.env.sa` | Static vars + `OP_SERVICE_ACCOUNT_TOKEN` |
| `/opt/onemind/docker-compose.yml` | Docker stack definition |

---

## OpenClaw Config

- **Config**: `/home/onemind/.openclaw/openclaw.json`
- **Version**: 2026.3.8 (build 3caab92)
- **Model**: `azure-openai/gpt-5.2` (Legacy) / `azure-openai/gpt-4o` (domain agents)
- **Agents**: 9 total (Legacy + Spartan, Oracle, Eden, Empire, Haven, Guardian, Grid, Trinity)
- **Memory search**: Hybrid BM25+vector, indexes `/home/onemind/onemind-os/codex`
- **TTS**: OpenAI
- **Tailscale HTTPS proxy**: `https://onemind-main.tail717ba0.ts.net` → `http://127.0.0.1:18789`

---

## Service Management

```bash
# OpenClaw
systemctl status openclaw
systemctl restart openclaw
journalctl -u openclaw -f

# Studio / SwarmClaw
systemctl status openclaw-studio
systemctl status swarmclaw
systemctl restart openclaw-studio swarmclaw

# Docker
cd /opt/onemind && docker compose ps
docker compose logs -f nats
docker compose logs -f onemind-omnidb

# Codex sync
systemctl status codex-watch
systemctl status codex-pull
```

---

## Backups

All backups run via cron (`/etc/cron.d/onemind-backups`). Scripts live in `/home/onemind/backups/`.

| Backup | Schedule | Retention | Script |
|--------|----------|-----------|--------|
| **DO Auto Backups** | Weekly (DO-managed) | 4 snapshots | n/a |
| **PostgreSQL Dump** | Daily 2:00 UTC | 7 days | `backup-postgres.sh` |
| **Docker Volumes** | Daily 2:30 UTC | 7 days | `backup-docker-volumes.sh` |
| **Off-site (DO Spaces)** | Daily 3:00 UTC | Mirrors local | `backup-offsite.sh` |

**Output location**: `/home/onemind/backups/{postgres,docker-volumes}/`
**Off-site bucket**: `s3://onemind-backups/onemind-main/` (nyc3)

```bash
# Check latest backups
ls -lh /home/onemind/backups/postgres/
ls -lh /home/onemind/backups/docker-volumes/

# Check cron logs
tail -20 /home/onemind/backups/postgres.log
tail -20 /home/onemind/backups/offsite.log
```

See `infra/ops/BACKUPS.md` for full backup & recovery documentation.

---

## Monitoring (Datadog)

**Agent**: Datadog v7.76.3 · **Site**: `us5.datadoghq.com`
**Config**: `/etc/datadog-agent/datadog.yaml`

| Monitor ID | Monitor | Threshold |
|------------|---------|-----------|
| 18513582 | Host Unreachable | 5 min no agent |
| 18513584 | High CPU | >85% (5m avg) |
| 18513586 | High Memory | <10% available |
| 18513588 | Disk Usage | >80% |
| 18513590 | OpenClaw Process Down | Not running |
| 18513592 | NATS Container Down | Container count < 1 |

All alert to `zeus@zeusdelacruz.me`. See `infra/ops/OBSERVABILITY.md` for full details.

---

## Monorepo Sync Flow

```
Mac (local edits) → git push → GitHub
    ↓ codex-pull.timer (every 1 min)
/home/onemind/onemind-os/ (VPS)
    ↓ openclaw-reload.service (on skill file changes)
openclaw.service hot-reloads workspace
    ↓ codex-watch.service (Legacy writes)
Auto-commit + push back to GitHub
```
