# VPS Gateway — Setup & Operations

> Last updated: 2026-04-15 — onemind-main deleted, DO hardening applied

## Active Nodes

### onemindos-cloud (Primary — Active ✅)

| Property | Value |
|----------|-------|
| Provider | DigitalOcean Droplet |
| Region | NYC3 |
| Public IP | `64.227.8.23` |
| Tailscale IP | `100.125.47.103` |
| Tailscale Hostname | `onemindos-cloud` |
| Tailscale URL | `https://onemindos-cloud.tail717ba0.ts.net` |
| OS | Ubuntu 24.04 LTS |
| Specs | 31 GB RAM, CPU-only |
| SSH User | `root` |
| OpenClaw Version | `2026.4.14` (build 323493f) |
| Gateway Port | `18789` |
| Bind Mode | `loopback` (exposed via Tailscale serve) |
| Service Manager | systemd USER service (`openclaw-gateway.service`) |
| Config Path | `/root/.openclaw/` |
| Systemd Drop-ins | `/root/.config/systemd/user/openclaw-gateway.service.d/` |
| Gateway Token | 1Password: `00-24 SO / OpenClaw Gateway / gateway-token` |

### onemind-hermes (Secondary — Unconfigured ⚠️)

| Property | Value |
|----------|-------|
| Provider | DigitalOcean Droplet |
| Tailscale IP | `100.115.131.81` |
| Status | Joined Tailscale but never fully configured — no active workloads |

### onemind-main (DELETED — April 15, 2026)

| Property | Value |
|----------|-------|
| Provider | DigitalOcean Droplet |
| Droplet ID | `558147996` (was `565013040` in stale notes) |
| Public IP | `104.248.4.212` |
| Tailscale Hostname | `onemind-main.tail717ba0.ts.net` |
| OpenClaw Version | `2026.4.11` (system-level service, not user-level) |
| Status | **DELETED** — April 15, 2026 via `doctl compute droplet delete 558147996 --force` |
| NATS | Was running — see `nats/` folder for archived config + scripts |

This VPS had the NATS JetStream infrastructure and 23-skill legacy workspace.
All unique custom code extracted to codex before deletion.
See `openclaw/nats/NATS.md` for the full NATS architecture documentation.

---

## SSH Access

```bash
# Primary (onemindos-cloud)
ssh root@64.227.8.23

# Via Tailscale (preferred — no public exposure)
ssh root@100.125.47.103

# Secondary (onemind-hermes, unconfigured)
ssh root@100.115.131.81  # via Tailscale only
```

---

## Directory Structure (onemindos-cloud)

```
/root/.openclaw/
├── openclaw.json                     ← Main gateway config
├── agents/
│   └── legacy/
│       └── agent/
│           └── models.json           ← Per-agent model overrides
├── devices/paired.json               ← Paired devices (Mac: 5a4ff116, 100.114.142.55)
├── identity/
├── memory/
├── logs/
└── media/

/root/.config/systemd/user/openclaw-gateway.service.d/
├── ollama-env.conf    ← AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_DEFAULT_REGION
├── deepgram-env.conf  ← DEEPGRAM_API_KEY (STT)
├── azure-env.conf     ← ⚠️ STALE — Azure RG deleted, needs cleanup
└── copilot.conf       ← ⚠️ STALE — old PAT, OAuth now used instead
```

---

## Environment Variables (onemindos-cloud)

Injected via systemd drop-in files in `/root/.config/systemd/user/openclaw-gateway.service.d/`:

```bash
# ollama-env.conf — AWS Bedrock credentials
AWS_ACCESS_KEY_ID=AKIA6CU7...
AWS_SECRET_ACCESS_KEY=...
AWS_DEFAULT_REGION=us-east-1

# deepgram-env.conf — Speech-to-Text
DEEPGRAM_API_KEY=...  # 88 chars, 1Password: lrmv5ai5xauh7tghxnvmfnwqeu
```

Drop-ins to clean up:
```bash
# ⚠️ Remove after Azure fresh start:
rm /root/.config/systemd/user/openclaw-gateway.service.d/azure-env.conf
rm /root/.config/systemd/user/openclaw-gateway.service.d/copilot.conf
systemctl --user daemon-reload && systemctl --user restart openclaw-gateway
```

---

## Starting / Stopping the Gateway

> **Note:** On onemindos-cloud, OpenClaw runs as a **USER** systemd service under `root`.
> Use `systemctl --user` (NOT `sudo systemctl`). This is different from the old setup.

### systemd (preferred)

```bash
# Start
systemctl --user start openclaw-gateway

# Stop
systemctl --user stop openclaw-gateway

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

## Firewall

### DO Cloud Firewall — `onemindos-cloud-fw` (Active ✅)

Applied April 15, 2026 via `doctl compute firewall create`. Managed in DO control panel.

| Direction | Protocol | Port | Source/Dest |
|-----------|----------|------|-------------|
| Inbound | TCP | 22 | 0.0.0.0/0 (SSH) |
| Inbound | TCP | 443 | 0.0.0.0/0 (HTTPS / Tailscale serve) |
| Inbound | TCP | 18789 | 0.0.0.0/0 (OpenClaw gateway) |
| Inbound | UDP | 41641 | 0.0.0.0/0 (Tailscale WireGuard) |
| Inbound | ICMP | — | 0.0.0.0/0 (ping) |
| Outbound | All | All | 0.0.0.0/0 |

> DO cloud firewall is upstream of the VPS — packets never reach the OS if blocked here.

### UFW (OS-level)

Installed but **inactive** — DO cloud firewall handles perimeter protection.
Can be enabled as a second layer if needed: `ufw enable`

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

---

## DigitalOcean Hardening (Applied April 15, 2026)

Full hardening pass applied to `onemindos-cloud` (565014820) after decommissioning `onemind-main`.

### DO Platform Controls

| Control | Details |
|---------|---------|
| **Automated backups** | Weekly (Sundays) — enabled via `doctl compute droplet-action enable-backups 565014820` |
| **Snapshot** | `onemindos-cloud-hardened-2026-04-15` — taken as baseline before hardening |
| **Cloud Firewall** | `onemindos-cloud-fw` — see Firewall section above |
| **Monitoring alerts** | 4 policies → `zeus@zeusdelacruz.com` (see below) |
| **Uptime check** | `onemindos-cloud gateway` — HTTPS check from US-East, US-West, EU-West |
| **Old snapshots** | 3 orphaned `onemind-main` snapshots deleted (~150 GiB, ~$7.50/mo saved) |

### Monitoring Alert Policies

| UUID | Alert | Threshold | Window |
|------|-------|-----------|--------|
| `8d92e01f` | CPU > 80% | 80% | 5m |
| `171832fa` | Memory > 85% | 85% | 5m |
| `e513786b` | Disk > 80% | 80% | 5m |
| `e956e9ee` | Bandwidth spike | 500 Mbps outbound | 5m |

Uptime check ID: `606be2ce-d171-435c-a7d7-02ffbeaec1f1`
Uptime alert ID: `887773f0-79bb-4ff1-a041-d571f13879fd` (fires on `down_global`)

### VPS-Level Hardening

**SSH** (`/etc/ssh/sshd_config`):
- `PermitRootLogin prohibit-password` — key-only, no password auth
- `PasswordAuthentication no`
- `MaxAuthTries 3`
- `LoginGraceTime 30`
- `X11Forwarding no`
- Backup: `/etc/ssh/sshd_config.bak.20260415`

**Fail2ban** (`/etc/fail2ban/jail.local`):
- SSH jail: 3 attempts → 24-hour ban
- findtime: 600s, bantime: 86400s
- Service: `systemctl status fail2ban`

**Unattended-upgrades** (`/etc/apt/apt.conf.d/20auto-upgrades`):
- Security patches auto-apply daily
- Auto-cleans every 7 days
- No automatic reboots

**Sysctl** (`/etc/sysctl.d/99-onemind-harden.conf`):
- `rp_filter = 1` — IP spoofing protection
- `tcp_syncookies = 1` — SYN flood protection
- `tcp_max_syn_backlog = 2048`
- `log_martians = 1` — log suspicious packets
- `accept_source_route = 0` — disable source routing
- `send_redirects = 0`
- `tcp_rfc1337 = 1` — time-wait assassination protection

### Pending After Hardening

```bash
# Reboot to load new kernel (flagged during fail2ban install)
ssh root@64.227.8.23 'reboot'

# Verify after reboot
ssh root@64.227.8.23 'uname -r && systemctl --user is-active openclaw-gateway && fail2ban-client status sshd'
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
