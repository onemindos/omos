---
id: so-02-omos-ha
title: "omos-ha — Home Assistant Edge Node"
type: reference
status: active
created: 2026-05-08
modified: 2026-05-08
owner: zeus
tags: [domain/so, infra/edge, infra/home-assistant, infra/tailscale]
---

# omos-ha — Home Assistant Edge Node

> Local LAN · Home Assistant OS · Smart home automation hub · Tailscale edge node

## Identity

| Attribute | Value |
|-----------|-------|
| Hostname | `omos-ha` |
| LAN IP | `192.168.1.180` |
| Tailscale IP | `100.64.25.12` |
| MAC | `e0:d4:64:ff:05:f0` |
| Platform | Home Assistant OS |
| HA OS Version | Current (Tailscale add-on `0.28.1`) |
| Location | Local LAN — no public IP |

## Access

```bash
# Via Tailscale (preferred — when on tailnet)
ssh hassio@100.64.25.12

# Via LAN (when on same network)
ssh hassio@192.168.1.180

# HA Web UI
http://192.168.1.180:8123
```

SSH key required: `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOtEsVg+0RB9Vf5kGrsnwW3iw2MFG0Jbg1PNC+dFrqT6 zeus@mac`

## Tailscale Setup

Tailscale runs as a **HA add-on** (`addon_a0d7b954_tailscale` — hassio-addons/tailscale v0.28.1), not as a native system service. The binary is at `/opt/tailscale` inside the container.

**Critical**: Do NOT run `tailscale set --ssh` manually on this node. The add-on schema (v0.28.1) does not include `ssh` as an option. If `RunSSH: true` ends up in `tailscaled.state` from a manual command, the add-on will crash on every restart with no way to fix it via the HA UI. See the fix runbook if this happens.

Options file location (host path):
```
/mnt/data/supervisor/addons/data/a0d7b954_tailscale/options.json
```

## Known Issue & Fix — Tailscale Add-on Crash on Restart

**Symptom**: `omos-ha` disappears from `tailscale status`, Tailscale ping times out.

**Root cause**: Add-on `options.json` missing `"ssh": true`, but `tailscaled.state` has `--ssh` from a prior manual `tailscale set --ssh` command. On startup, Tailscale refuses to run because the options conflict.

**Fix** (via LAN SSH):
```bash
# 1. SSH to omos-ha on LAN
ssh hassio@192.168.1.180

# 2. Write corrected options.json via temp Alpine container
sudo docker run --rm \
  -v /mnt/data/supervisor/addons/data/a0d7b954_tailscale:/tsdata \
  alpine sh -c 'cat > /tsdata/options.json << EOF
{
  "accept_dns": true,
  "accept_routes": true,
  "advertise_exit_node": false,
  "advertise_connector": false,
  "advertise_routes": [],
  "always_use_derp": false,
  "log_level": "info",
  "login_server": "https://controlplane.tailscale.com",
  "share_homeassistant": "disabled",
  "share_on_port": "443",
  "snat_subnet_routes": true,
  "stateful_filtering": false,
  "tags": [],
  "taildrive": {"addons":false,"addon_configs":false,"backup":false,"config":false,"media":false,"share":false,"ssl":false},
  "taildrop": true,
  "userspace_networking": false,
  "ssh": true
}
EOF'

# 3. Start the container
sudo docker start addon_a0d7b954_tailscale

# 4. Verify (from Mac)
/Applications/Tailscale.app/Contents/MacOS/Tailscale status | grep omos-ha
```

**Prevention**: Never use `docker exec ... tailscale set --ssh` manually. Always set `ssh: true` in the HA add-on UI (Settings → Add-ons → Tailscale → Configuration) so the options.json stays in sync.

## Docker Commands (from LAN SSH)

```bash
# Check container status
sudo docker ps | grep tailscale

# View logs
sudo docker logs --tail 30 addon_a0d7b954_tailscale

# Restart add-on
sudo docker restart addon_a0d7b954_tailscale
```

## Security Posture

- No public IP — LAN-only device
- Reachable from tailnet via `100.64.25.12`
- Reachable from LAN via `192.168.1.180`
- HA web UI on port 8123 (LAN only, not port-forwarded)

---

## Related

- [[TAILSCALE-NETWORK]] — Full tailnet map
- [[TAILSCALE-SSH]] — SSH runbook
- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
