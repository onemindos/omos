---
name: OpenClaw — VPS (Linux Server)
updated: 2026-05-11
---

# VPS — Linux Server

OpenClaw Gateway on Linux VPS. This covers our specific AWS Lightsail setup plus general Linux tuning.

**Official docs:** [docs.openclaw.ai/vps](https://docs.openclaw.ai/vps.md)

---

## Our Setup

| Node | Provider | Purpose |
|---|---|---|
| omos-core | AWS Lightsail | Primary gateway — 19 agents |
| omos-sentinel | AWS Lightsail | Watchdog + failover |
| omos-titan | AWS Lightsail | Development / secondary |

SSH: `ssh onemindos-core` / `ssh omos-sentinel` (see `~/.ssh/config`)  
Tailscale: omos-core at `100.123.126.34`

---

## How Cloud Setups Work

- The **Gateway runs on the VPS** and owns state + workspace
- Connect from Mac or phone via the **Control UI** or **Tailscale/SSH**
- Treat the VPS as the source of truth — **back up** state + workspace regularly
- Keep the Gateway on loopback; access via SSH tunnel or Tailscale Serve
- If binding to `lan` or `tailnet`, require `gateway.auth.token` or `gateway.auth.password`

---

## Security Baseline

Tailscale handles our VPS access — the Gateway stays loopback-bound, accessible only via Tailscale Serve or SSH tunnel. This is the correct setup for One Mind.

Do not expose the gateway directly to the public internet without auth token protection.

---

## Using Nodes with a VPS

The Gateway stays in the cloud; **nodes** on local devices (Mac, iOS) provide:
- Local screen/camera/canvas
- `system.run` capabilities (run commands on the local machine)

Docs: `cli/nodes.md`, `cli/node.md`

---

## Startup Tuning (Small VMs)

If CLI commands feel slow on Lightsail instances, enable Node's module compile cache:

```bash
grep -q 'NODE_COMPILE_CACHE=/var/tmp/openclaw-compile-cache' ~/.bashrc || cat >> ~/.bashrc <<'EOF'
export NODE_COMPILE_CACHE=/var/tmp/openclaw-compile-cache
mkdir -p /var/tmp/openclaw-compile-cache
export OPENCLAW_NO_RESPAWN=1
EOF
source ~/.bashrc
```

- `NODE_COMPILE_CACHE` — improves repeated command startup times
- `OPENCLAW_NO_RESPAWN=1` — avoids extra startup overhead from self-respawn path
- First run warms the cache; subsequent runs are faster

---

## systemd Tuning

Edit the user service unit:

```bash
systemctl --user edit openclaw-gateway.service
```

```ini
[Service]
Environment=OPENCLAW_NO_RESPAWN=1
Environment=NODE_COMPILE_CACHE=/var/tmp/openclaw-compile-cache
Restart=always
RestartSec=2
TimeoutStartSec=90
```

Reload:

```bash
systemctl --user daemon-reload
systemctl --user restart openclaw-gateway
```

---

## Our systemd Services

| Service | Purpose |
|---|---|
| `openclaw-gateway.service` | OpenClaw gateway process |
| `onemind-pull.timer` | Vault git pull every 10s |
| `onemind-watch.service` | Vault inotify → git push (3s debounce) |

Check status:
```bash
systemctl --user status openclaw-gateway
systemctl --user status onemind-pull.timer
systemctl --user status onemind-watch
```

---

## Backups

```bash
# Create a backup before any major change
openclaw backup create

# Verify backup
openclaw backup verify
```

Backups stored in: `install (Install & Deploy)/backups/`

---

## Related

- [gateway/README.md](../gateway/README.md)
- [install (Install & Deploy)/README.md](../install%20%28Install%20%26%20Deploy%29/README.md)
- [cli/gateway.md](../cli/gateway.md)
