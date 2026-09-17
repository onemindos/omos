# OpenClaw Quick Reference

## Credentials

| Item | Value |
|------|-------|
| Gateway Token | `onemind-legacy-2026` |
| Hooks Token | `onemind-hooks-2026` |
| VPS Public IP | `134.122.8.135` |
| VPS Tailscale IP | `100.104.61.29` |
| Gateway Port | `18789` |
| SSH Key | `N/A (key-based SSH via root)` |
| SSH User | `ubuntu` |
| OpenClaw User | `openclaw` |
| Apple Team ID | `9YMN85XAKP` |

---

## Node IDs

| Node | ID | Status |
|------|----|---------|
| Zeus-Mac | `206e8a72b7c6ce558139ad380e7c2eb240fc531b67c5f7296392d92be8d673b5` | paired · connected |
| Zeus's MacBook Pro | `50efbb0b87d9bf37bb732a000ad29aeef1fb42ffbdce9bf86c21394fc2722510` | paired · disconnected |
| iPhone 17 Pro Max | _(pending first connection)_ | not yet paired |

Full node registry with capabilities and connect commands: [NODES.md](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/openclaw/NODES.md)

---

## Agents

| ID | Name | Model | Domain |
|----|------|-------|--------|
| `legacy` | Legacy (default) | `openai/gpt-4o` | UI 00-24 |
| `hp` | HP Holistic Performance | `anthropic/claude-sonnet-4-5` | HP 25-49 |
| `le` | LE Legacy Evolution | `anthropic/claude-sonnet-4-5` | LE 50-74 |
| `ge` | GE Generational Empire | `anthropic/claude-opus-4-6` | GE 75-99 |

---

## SSH Commands

```bash
# Quick SSH
ssh root@134.122.8.135

# SSH alias (add to ~/.ssh/config):
# Host onemind-gateway
#   HostName 134.122.8.135
#   User ubuntu
#   IdentityFile N/A (key-based SSH via root)

# SSH tunnel for browser access
ssh -f -N -L 18789:127.0.0.1:18789 root@134.122.8.135
```

---

## Gateway Health

```bash
# From VPS (as openclaw user)
ssh root@134.122.8.135 \
  "sudo -u onemind env OPENCLAW_GATEWAY_TOKEN=onemind-legacy-2026 \
   OPENCLAW_HOOKS_TOKEN=onemind-hooks-2026 \
   openclaw gateway health"

# From Mac (local CLI after configuring ~/.openclaw/openclaw.json with remote.token)
openclaw gateway health

# Via curl (after SSH tunnel)
curl http://localhost:18789/health

# Web UI (after SSH tunnel)
open http://localhost:18789
```

---

## Gateway Management

```bash
# Status (systemd)
ssh root@134.122.8.135 "sudo -u onemind systemctl --user status openclaw-gateway"

# Restart
ssh root@134.122.8.135 "sudo -u onemind systemctl --user restart openclaw-gateway"

# Stop
ssh root@134.122.8.135 "sudo -u onemind systemctl --user stop openclaw-gateway"

# Logs (last 50 lines)
ssh root@134.122.8.135 "sudo -u onemind journalctl --user -u openclaw-gateway -n 50"

# Follow logs live
ssh root@134.122.8.135 "sudo -u onemind journalctl --user -u openclaw-gateway -f"
```

---

## Agent Commands

```bash
# Environment prefix for VPS commands
VPS_ENV="sudo -u onemind env OPENCLAW_GATEWAY_TOKEN=onemind-legacy-2026 OPENCLAW_HOOKS_TOKEN=onemind-hooks-2026 openclaw"

# List agents
ssh root@134.122.8.135 "$VPS_ENV agents list"

# List with bindings
ssh root@134.122.8.135 "$VPS_ENV agents list --bindings"

# Nodes status
ssh root@134.122.8.135 "$VPS_ENV nodes status"

# Check all
ssh root@134.122.8.135 "$VPS_ENV agents list --bindings && $VPS_ENV nodes status"
```

---

## iOS App Setup

```bash
# Build (from source)
cd ~/Documents/openclaw/apps/ios/
xcodegen generate
open OpenClaw.xcodeproj
# → Select iPhone target → Press Run ▶

# Gateway config for iOS app
Host:  100.104.61.29
Port:  18789
TLS:   OFF
Token: onemind-legacy-2026

# QR/import base64
eyJob3N0IjoiMTAwLjc2LjEwNC45MiIsInBvcnQiOjE4Nzg5LCJ0bHMiOmZhbHNlLCJ0b2tlbiI6Im9uZW1pbmQtbGVnYWN5LTIwMjYifQ==
```

---

## Important VPS Paths

```
/home/onemind/.openclaw/
├── openclaw.json          ← Main config (4 agents, channels, bindings)
├── env                    ← Token exports (sourced by systemd service)
├── workspace-legacy/      ← Legacy agent workspace (SOUL.md, AGENTS.md, skills/)
├── workspace-hp/          ← HP agent workspace
├── workspace-le/          ← LE agent workspace
├── workspace-ge/          ← GE agent workspace
├── agents/
│   ├── legacy/agent/      ← Agent session DB
│   ├── hp/agent/
│   ├── le/agent/
│   └── ge/agent/
└── .local/bin/openclaw    ← Binary

/etc/systemd/system/
└── openclaw-gateway.service   ← Systemd unit
```

---

## Local Mac Paths

```
~/.openclaw/
├── openclaw.json          ← Local node config (with remote.token for CLI auth)
└── workspace/             ← Local workspace (optional for local gateway)

~/Documents/openclaw/apps/ios/   ← iOS app source (patched)
```

---

## Troubleshooting One-Liners

```bash
# Is gateway running?
ssh root@134.122.8.135 "ps aux | grep openclaw | grep -v grep"

# Is port listening?
ssh root@134.122.8.135 "ss -tlnp | grep 18789"

# Is Tailscale up on VPS?
ssh root@134.122.8.135 "tailscale status"

# View openclaw config on VPS
ssh root@134.122.8.135 "sudo cat /home/onemind/.openclaw/openclaw.json"

# View openclaw env on VPS
ssh root@134.122.8.135 "sudo cat /home/onemind/.openclaw/env"

# Restart everything (if gateway hung)
ssh root@134.122.8.135 "sudo -u onemind systemctl --user restart openclaw-gateway && sleep 3 && sudo -u onemind systemctl --user status openclaw-gateway"
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
