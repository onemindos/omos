# VPS Gateway — Setup & Operations

## Infrastructure

| Property | Value |
|----------|-------|
| Provider | DigitalOcean Droplet |
| Region | NYC3 |
| Public IP | `134.122.8.135` |
| Tailscale IP | `100.104.61.29` |
| Tailscale Hostname | `onemind-gateway` |
| OS | Ubuntu 24.04 LTS |
| Specs | 8 GB RAM, 4 vCPU Intel, 240 GB NVMe |
| SSH User | `root` |
| OpenClaw User | `onemind` |
| OpenClaw Version | `2026.3.2` |
| Gateway Port | `18789` |
| Bind Mode | `tailnet` |
| Service Manager | systemd (`openclaw.service`) |
| Env File | `/opt/onemind/.env` |
| Channels | Discord only |

---

## SSH Access

```bash
# Direct
ssh root@134.122.8.135

# Via Tailscale
ssh root@100.104.61.29

# As onemind user (for openclaw operations)
ssh root@134.122.8.135
sudo -u onemind -i
```

---

## Directory Structure

```
/home/onemind/.openclaw/
├── openclaw.json                 ← Main gateway config (Discord-only, 4 agents)
├── openclaw.json.do-working      ← Backup of known-good DO config
├── openclaw.json.aws-multi-agent ← Reference: old AWS multi-agent config
├── openclaw.json.aws-pre-codex   ← Reference: old AWS pre-codex config
├── workspace/                    ← Main workspace (SOUL.md, AGENTS.md, etc.)
│   └── skills/                   ← 15 installed skills
├── workspace-legacy/             ← LEGACY agent workspace
├── workspace-hp/                 ← HP agent workspace
├── workspace-le/                 ← LE agent workspace
├── workspace-ge/                 ← GE agent workspace
├── agents/
│   ├── legacy/                   ← LEGACY sessions + auth
│   ├── main/agent/               ← Main agent state
│   ├── hp/agent/                 ← HP agent state
│   ├── le/agent/                 ← LE agent state
│   └── ge/agent/                 ← GE agent state
├── memory/legacy.sqlite          ← Agent memory DB
├── devices/paired.json           ← 5 paired devices
├── identity/device.json          ← Ed25519 keypair (ead40ef7...)
├── exec-approvals.json           ← Execution approval config
├── settings/voicewake.json       ← Voice wake config
├── logs/                         ← Gateway logs
├── canvas/                       ← HTML/CSS/JS surfaces
├── cron/                         ← Scheduled jobs
└── media/                        ← Media files

/opt/onemind/
├── .env                          ← Environment variables (loaded by systemd)
├── docker-compose.yml            ← NATS + OmniDB + Bridge
├── .do/                          ← Deployment scripts
│   ├── setup-droplet.sh
│   ├── install-openclaw.sh
│   ├── setup-firewall.sh
│   └── health-check.sh
└── openclaw/
    └── nats/hub.conf             ← NATS JetStream config
```

---

## Environment Variables (`/opt/onemind/.env`)

```bash
ANTHROPIC_API_KEY=sk-ant-api03-...       # Claude (Opus, Sonnet)
OPENAI_API_KEY=sk-proj-...               # GPT-4o
OPENCLAW_GATEWAY_TOKEN=...
OPENCLAW_HOOKS_TOKEN=...
DISCORD_BOT_TOKEN=...                    # Discord channel
ELEVENLABS_API_KEY=sk_...                # Voice
ELEVENLABS_VOICE_ID=...
OPENCLAW_GATEWAY_URL=http://localhost:18789
OPENCLAW_ALLOW_INSECURE_PRIVATE_WS=1
NODE_COMPILE_CACHE=/var/tmp/openclaw-compile-cache
OPENCLAW_NO_RESPAWN=1
```

> These are loaded via systemd `EnvironmentFile=/opt/onemind/.env`.

---

## Starting / Stopping the Gateway

### systemd (preferred)

```bash
# Start
sudo systemctl start openclaw

# Stop
sudo systemctl stop openclaw

# Restart
sudo systemctl restart openclaw

# Status
sudo systemctl status openclaw

# Enable on boot (already enabled)
sudo systemctl enable openclaw
```

### Logs

```bash
# Follow logs
journalctl -u openclaw -f

# Last 100 lines
journalctl -u openclaw -n 100 --no-pager

# Since last boot
journalctl -u openclaw -b
```

---

## Docker Services

```bash
# Start all services
cd /opt/onemind && docker compose --profile data up -d

# Check status
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Services:
#   onemind-nats           → 4222 (JetStream)
#   onemind-omnidb         → 5432 (PostgreSQL + pgvector)
#   onemind-legacy-bridge  → NATS ↔ OpenClaw relay
```

---

## Health Check

```bash
# Run the health check script
bash /opt/onemind/.do/health-check.sh

# Quick manual checks
systemctl is-active openclaw
docker ps
curl -s http://localhost:18789/health 2>/dev/null || echo "No HTTP health endpoint"
```

---

## Firewall (UFW)

```
22/tcp    ALLOW   ← SSH
4222/tcp  ALLOW on tailscale0  ← NATS (Tailscale only)
18789/tcp ALLOW on tailscale0  ← Gateway (Tailscale only)
```

> The gateway is NOT exposed on the public internet. Access is Tailscale-only.
> Public IP (134.122.8.135) is used only for SSH.

---

## Connected Channels

| Channel | Status | Notes |
|---------|--------|-------|
| Discord | Enabled | Set real `DISCORD_BOT_TOKEN` in `.env` |

Telegram, WhatsApp, iMessage are **disabled** by design. Discord is the sole channel.

---

## Gateway UI

Access via SSH tunnel:
```bash
# 1. Open tunnel
ssh -f -N -L 18789:127.0.0.1:18789 root@134.122.8.135

# 2. Open browser
open http://localhost:18789
```

Or via Tailscale directly:
```bash
open http://100.104.61.29:18789
```

---

## Nodes (Paired Devices)

| Node | ID (short) | Status |
|------|-----------|--------|
| Zeus-Mac | `206e8a72...` | paired |
| Zeus's MacBook Pro | `50efbb0b...` | paired |

5 total paired devices in `devices/paired.json`.

---

## Migration History

Migrated from AWS Lightsail (`44.222.153.60`) on 2026-03-03:
- All agent sessions, sub-workspaces, memory DB, paired devices, identity keypair
- AWS configs preserved as reference files (`.aws-multi-agent`, `.aws-pre-codex`)
- Full backup archived locally at `/tmp/aws-migration/openclaw-full-backup.tar.gz`
- AWS snapshot taken before deletion as safety net
