# OneMind OS — Sentinel VPS (159.203.122.175)

> DigitalOcean droplet · Ubuntu 24.04 · Rescue/failover gateway

---

## Purpose

Sentinel is a **lightweight rescue node** — it runs a stripped-down OpenClaw gateway (`openclaw-sentinel`) as a failover if the main VPS goes down. Legacy can be reached via Sentinel when the main node is unreachable.

---

## Access

| Method | Address |
|--------|---------|
| SSH | `ssh root@159.203.122.175` |
| Tailscale IP | `100.67.223.5` |
| OpenClaw Gateway | `http://localhost:19789` (internal, Tailscale only) |

---

## Running Services

| Service | Description |
|---------|-------------|
| `openclaw-sentinel.service` | Sentinel OpenClaw gateway (port 19789) |
| `tailscaled.service` | Tailscale VPN |

**No Docker** — intentionally lean. No database, no NATS, no dashboards.

---

## OpenClaw Config

| Setting | Value |
|---------|-------|
| Config path | `/home/onemind/.openclaw-sentinel/openclaw.json` |
| Workspace | `/home/onemind/onemind-os/openclaw/workspaces/sentinel` |
| Gateway port | `19789` |
| Profile | `sentinel` |
| NATS | `nats://104.248.4.212:4222` (points to main VPS) |
| State dir | `/home/onemind/.openclaw-sentinel/` |

---

## Key Directories

| Path | Purpose |
|------|---------|
| `/home/onemind/onemind-os/` | Monorepo clone |
| `/home/onemind/.openclaw-sentinel/` | Sentinel state + config |
| `/opt/sentinel/.env.sa` | Static env vars (no secrets template — standalone) |

---

## Service Definition

**File**: `/etc/systemd/system/openclaw-sentinel.service`

```ini
[Unit]
Description=OneMind OS — Sentinel (Rescue Gateway)
After=network-online.target tailscaled.service
Wants=network-online.target

[Service]
User=onemind
WorkingDirectory=/home/onemind
EnvironmentFile=/opt/sentinel/.env.sa
ExecStart=/usr/bin/openclaw gateway
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=openclaw-sentinel
```

---

## Ports

| Port | Service |
|------|---------|
| 22 | SSH |
| 19789 | OpenClaw sentinel gateway (Tailscale-only, bound to `::1`) |

---

## Service Management

```bash
systemctl status openclaw-sentinel
systemctl restart openclaw-sentinel
journalctl -u openclaw-sentinel -f
```

---

## Monorepo Sync

Same pattern as main VPS — `codex-pull.timer` pulls GitHub every 1 min.
Sentinel workspace: `openclaw/workspaces/sentinel/`

---

## Backups

| Backup | Status |
|--------|--------|
| **DO Auto Backups** | Enabled (weekly, 4 snapshots, ~$2.40/mo) |
| **Git sync** | `codex-pull.timer` pulls from GitHub every 1 min |

No database or Docker volumes to back up — Sentinel is stateless by design.

---

## Monitoring (Datadog)

**Agent**: Datadog v7.76.3 · **Site**: `us5.datadoghq.com`
**Config**: `/etc/datadog-agent/datadog.yaml`

**Delayed start**: DD agent waits 45s after openclaw-sentinel boots to avoid memory pressure.
Override: `/etc/systemd/system/datadog-agent.service.d/after-openclaw.conf`

| Monitor ID | Monitor | Threshold |
|------------|---------|-----------|
| 18513583 | Host Unreachable | 5 min no agent |
| 18513585 | High CPU | >85% (5m avg) |
| 18513587 | High Memory | <20% available |
| 18513589 | Disk Usage | >75% |
| 18513591 | OpenClaw Process Down | Not running |

All alert to `zeus@zeusdelacruz.me`. See `infra/ops/OBSERVABILITY.md` for full details.
