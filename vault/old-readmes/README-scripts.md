# Infrastructure

Everything that runs on the VPS and the configs to deploy it. Single source of truth for what's running, where, and how to manage it.

---

## VPS Overview

| Item | Value |
|------|-------|
| Host | `root@104.248.4.212` |
| Droplet | `onemind-main` (ID: 558147996) |
| Spec | 4 vCPU / 16 GB RAM / 200 GB SSD — s-4vcpu-16gb-amd — $84/mo |
| OS | Ubuntu 24.04 |
| Tailscale | `100.77.52.93` — `onemind-main.tail717ba0.ts.net` |
| Gateway HTTPS | `https://onemind-main.tail717ba0.ts.net` |
| Security | UFW + Fail2ban, zero public ports except SSH |

```
┌─────────────────────────────────────────────────────────────────────┐
│   DigitalOcean Droplet  (Ubuntu 24.04 / s-4vcpu-16gb-amd / $84/mo) │
│                                                                     │
│  ┌───────────────────────────┐  ┌─────────────────────────────────┐│
│  │  HOST (bare metal)        │  │  DOCKER COMPOSE                 ││
│  │                           │  │                                  ││
│  │  OpenClaw Gateway :18789  │  │  NATS JetStream  :4222/:7422   ││
│  │  Node.js 22 LTS          │  │  MQTT bridge     :1883          ││
│  │  1Password CLI (op)       │  │  OmniDB (PG16)   :5432  [data] ││
│  │  Tailscale VPN            │  │  NATS Ingester           [data] ││
│  │  systemd services         │  │  legacy-bridge                  ││
│  │  codex-watch / pull       │  │  n8n             :5678  [full]  ││
│  │                           │  │  LiveKit          :7880  [full]  ││
│  └───────────────────────────┘  └─────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────────┘
```

---

## Port Map

| Port | Service | Protocol | Profile |
|------|---------|----------|---------|
| 18789 | OpenClaw Gateway | HTTP/WS | host |
| 4222 | NATS clients | NATS | default |
| 7422 | NATS leaf nodes | NATS | default |
| 1883 | MQTT (IoT/HA) | MQTT | default |
| 8222 | NATS monitoring | HTTP | default |
| 5432 | OmniDB (PostgreSQL) | TCP | data |
| 5678 | n8n workflows | HTTP | full |
| 7880 | LiveKit HTTP API | HTTP | full |
| 7881 | LiveKit WebRTC | TCP | full |
| 7882 | LiveKit WebRTC | UDP | full |

---

## Systemd Services

| Service | Purpose | Managed by |
|---------|---------|------------|
| `openclaw.service` | Legacy's runtime via `op run` | systemd |
| `codex-watch.service` | Auto-commits codex + workspace changes → GitHub | systemd |
| `openclaw-reload.service` | Hot-reloads Legacy when skill files change | inotify |
| `codex-pull.timer` | Pulls GitHub → monorepo every 1 min | systemd timer |

### Agent Architecture (9 agents)

The OpenClaw gateway runs 9 agents, each bound to a dedicated Discord channel:

| Agent | ID | Model | Discord Channel |
|-------|----|-------|-----------------|
| ⚡ Legacy | `legacy` | gpt-5.2 | `#legacy-chat` |
| ⚔️ Spartan | `spartan` | gpt-4o | `#spartan-chat` |
| 🔮 Oracle | `oracle` | gpt-4o | `#oracle-chat` |
| 🌳 Eden | `eden` | gpt-4o | `#eden-chat` |
| 👑 Empire | `empire` | gpt-5.2 | `#empire-chat` |
| 🌿 Haven | `haven` | gpt-4o | `#haven-chat` |
| 🛡️ Guardian | `guardian` | gpt-4o | `#guardian-chat` |
| ⚙️ Grid | `grid` | gpt-4o | `#grid-chat` |
| 🔺 Trinity | `trinity` | gpt-4o | `#trinity-chat` |

Config file: `/home/onemind/.openclaw/openclaw.json`
Per-agent state: `/home/onemind/.openclaw/agents/<id>/`
Workspaces: `/home/onemind/onemind-os/openclaw/workspaces/<id>/`

```bash
# Check all services
ssh root@134.122.8.135 "systemctl status openclaw codex-watch openclaw-reload codex-pull"

# Restart Legacy
ssh root@134.122.8.135 "systemctl restart openclaw"

# Follow logs
ssh root@134.122.8.135 "journalctl -u openclaw -f"
```

---

## Docker Compose Profiles

```bash
cd infra/
docker compose up -d                    # default: NATS + legacy-bridge
docker compose --profile data up -d     # + OmniDB + NATS ingester
docker compose --profile full up -d     # + n8n + LiveKit
docker compose logs -f                  # all logs
docker compose down                     # stop everything
```

---

## VPS Paths

| Path | Purpose |
|------|---------|
| `/opt/onemind/openclaw/workspace/` | Legacy's **active** OpenClaw workspace |
| `/opt/onemind/.env.tpl` | `op://` secret references (injected by `op run`) |
| `/opt/onemind/.env.sa` | Static vars + `OP_SERVICE_ACCOUNT_TOKEN` |
| `/home/onemind/onemind-os/` | Monorepo clone |
| `/home/onemind/onemind-os/codex/` | Codex vault — Legacy reads this path |

---

## Secrets (1Password)

All secrets stored in **1Password vault `00-24 SO (Sovereign Operations)`** (vault ID: `aaeyris5z44pdqspk277u2oppi`).

- `openclaw/deploy/env.tpl` — `op://` references, injected at runtime by `op run`
- `openclaw/deploy/env.sa.example` — template for static vars + SA token
- Never commit real secret values

```bash
# Add a new secret
op item create --category=login --vault="00-24 SO (Sovereign Operations)" \
  --title="My Service" "api-key[password]=actual-value"

# Add to env.tpl
echo 'MY_KEY=op://aaeyris5z44pdqspk277u2oppi/My Service/api-key' >> /opt/onemind/.env.tpl

# Restart to pick up
systemctl restart openclaw
```

---

## Directory Layout

Infra is shared services. OpenClaw-specific deploy config lives in `openclaw/deploy/`.

```
infra/                              ← Shared infrastructure (all runtimes)
├── docker-compose.yml              ← Core stack (NATS, OmniDB, bridge, n8n, LiveKit)
├── DO-DEPLOYMENT.md                ← Full DigitalOcean deployment guide
├── README.md                       ← This file
├── nats/                           ← NATS JetStream config + protocol docs
│   ├── hub.conf                    ← Cloud hub config
│   ├── leaf.conf                   ← Edge location config
│   ├── bus.py                      ← NATS bus utilities
│   ├── nats_bridge.py              ← NATS→webhook bridge
│   ├── nats_ingester.py            ← NATS→OmniDB pipeline
│   ├── nats_tools.py               ← NATS tool helpers
│   ├── NATS-BUS.md                 ← Bus architecture docs
│   ├── SUBJECTS.md                 ← Subject naming conventions
│   ├── COMMS_PROTOCOL.md           ← Communication protocol spec
│   ├── EVENT_FLOW.md               ← Event flow diagrams
│   └── timescaledb/                ← DB init scripts
└── scripts/                        ← DO droplet setup scripts
    ├── setup-droplet.sh            ← DO-specific base setup
    ├── setup-firewall.sh           ← UFW + security hardening
    ├── install-openclaw.sh         ← OpenClaw gateway install
    └── health-check.sh             ← VPS health check

openclaw/deploy/                    ← OpenClaw-specific deployment
├── env.tpl                         ← 1Password op:// secret references
├── env.sa.example                  ← Static vars template
├── env.example                     ← Human-readable env var list
├── systemd/                        ← Systemd service files
│   ├── openclaw.service
│   ├── openclaw-reload.service
│   ├── codex-watch.service
│   ├── codex-pull.service + timer
└── scripts/                        ← OpenClaw deploy scripts
    ├── setup.sh                    ← Full VPS setup from scratch
    ├── openclaw-reload.sh          ← Hot reload workspace on VPS
    ├── workspace-backup.sh         ← Backup VPS workspace → repo
    ├── setup-symlinks.sh           ← Sync repo workspace → VPS path
    ├── create-1p-items.sh          ← Migrate secrets → 1Password
    └── codex-watch.sh              ← Auto-commit codex changes
```
