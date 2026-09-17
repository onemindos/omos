# OpenClaw — Networking & Connectivity

## Architecture

```
iPhone (Tailscale 100.122.203.68)
    │ ws:// on port 18789
    ▼
VPS (Tailscale 100.104.61.29, Public 134.122.8.135)
    │
    └── OpenClaw Gateway :18789 (bound to 0.0.0.0 in config, but firewalled)
```

The gateway is **NOT exposed on the public internet**. Port `18789` is only open on the Tailscale interface (`tailscale0`). All clients must be on the same Tailscale network.

---

## Tailscale Setup

### VPS

```bash
# Install
curl -fsSL https://tailscale.com/install.sh | sh

# Authenticate
sudo tailscale up

# Confirm IP
tailscale ip -4
# → 100.104.61.29

# Hostname visible as
tailscale status | grep "onemind-gateway"
```

### iPhone

1. Install **Tailscale** from App Store
2. Log in with same Tailscale account
3. iPhone is assigned Tailscale IP `100.122.203.68` (hostname: `iphone182`)
4. Verify VPS reachable: `ping 100.104.61.29` in a terminal app

### Mac

```bash
# Tailscale for Mac via Homebrew
brew install --cask tailscale

# Or download from tailscale.com/download
# Then connect via menu bar icon
```

---

## Firewall Rules (UFW)

Port `18789` is allowed **only on the Tailscale interface**:

```bash
# Check current rules
sudo ufw status numbered

# Rule that exists:
18789 on tailscale0 ALLOW from any

# Public internet blocked (default):
# 18789 NOT in ALLOW for eth0/public
```

To update rules if needed:

```bash
# Allow on tailscale only
sudo ufw allow in on tailscale0 to any port 18789 proto tcp

# Deny from everywhere else (if not already default deny)
sudo ufw deny 18789

sudo ufw reload
```

---

## SSH Access

```bash
# Direct SSH
ssh \
    -o StrictHostKeyChecking=no \
    root@134.122.8.135

# As openclaw user (for agent commands)
ssh \
    root@134.122.8.135 \
    "sudo -u onemind env OPENCLAW_GATEWAY_TOKEN=onemind-legacy-2026 \
     OPENCLAW_HOOKS_TOKEN=onemind-hooks-2026 \
     openclaw gateway health"
```

### SSH Tunnel (for local browser access)

Forward the gateway port locally so you can open `http://localhost:18789` in your Mac browser:

```bash
ssh -f -N \
    -L 18789:127.0.0.1:18789 \
    \
    -o StrictHostKeyChecking=no \
    -o ServerAliveInterval=30 \
    root@134.122.8.135
```

After connecting: open `http://localhost:18789` in browser.

---

## Port Reference

| Port | Protocol | Interface | Purpose |
|------|----------|-----------|---------|
| 22 | TCP | Public | SSH access |
| 18789 | TCP/WS | tailscale0 only | OpenClaw gateway (agents) |
| 443 | TCP | Public | (Not currently used) |

---

## Connections Reference

Live `tailscale status` output (as of 2026-02-25):

| Device             | Tailscale Hostname  | Tailscale IP     | Public IP       | Platform | Status                             |
| ------------------ | ------------------- | ---------------- | --------------- | -------- | ---------------------------------- |
| VPS                | `onemind-gateway`   | `100.104.61.29`  | `134.122.8.135` | Linux    | Always on                          |
| Zeus's MacBook Pro | `zeuss-macbook-pro` | `100.114.142.55` | Dynamic         | macOS    | Active (relay: nyc)                |
| iPhone 17 Pro Max  | `iphone182`         | `100.122.203.68` | Dynamic         | iOS      | Offline (connects when app opened) |

> The MacBook Pro (`100.114.142.55`) also hosts the **BlueBubbles server** on port `1234`.

---

## OpenClaw Gateway Bind Config

In `/home/onemind/.openclaw/openclaw.json`:

```json
"gateway": {
  "port": 18789,
  "host": "0.0.0.0",
  "remote": {
    "token": "onemind-legacy-2026"
  }
}
```

`host: 0.0.0.0` means the process listens on all interfaces — but UFW ensures only Tailscale traffic reaches it.

---

## Troubleshooting Connectivity

```bash
# Is gateway process running?
ssh root@134.122.8.135 "ps aux | grep openclaw"

# Is port listening?
ssh root@134.122.8.135 "ss -tlnp | grep 18789"

# From VPS: is it reachable via Tailscale?
ssh root@134.122.8.135 "curl -s http://100.104.61.29:18789/health || echo FAIL"

# From Mac: test tunnel
curl -s http://localhost:18789/health

# From iPhone: test via Tailscale (use a terminal app)
curl http://100.104.61.29:18789/health
```

---

## Gateway Auto-Restart

The gateway runs as a `systemd` user service on the VPS:

```bash
# Status
ssh root@134.122.8.135 "sudo -u openclaw systemctl --user status openclaw-gateway"

# Restart
ssh root@134.122.8.135 "sudo -u openclaw systemctl --user restart openclaw-gateway"

# View logs
ssh root@134.122.8.135 "sudo -u openclaw journalctl --user -u openclaw-gateway -n 50"
```

The service file is at `/etc/systemd/system/openclaw-gateway.service`.

> Note: If the service was set up with `OPENCLAW_GATEWAY_TOKEN` and `OPENCLAW_HOOKS_TOKEN` sourced from `/home/onemind/.openclaw/env`, tokens persist across restarts without manual export.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
