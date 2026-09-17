# OpenClaw Nodes Registry

Live status as of: 2026-02-25
Gateway: `100.104.61.29:18789`

To refresh: `ssh root@134.122.8.135 "sudo -u onemind env OPENCLAW_GATEWAY_TOKEN=onemind-legacy-2026 OPENCLAW_HOOKS_TOKEN=onemind-hooks-2026 openclaw nodes status"`

---

## Active Nodes

### Zeus-Mac

| Field | Value |
|-------|-------|
| Display Name | Zeus-Mac |
| Node ID | `206e8a72b7c6ce558139ad380e7c2eb240fc531b67c5f7296392d92be8d673b5` |
| Status | paired · connected |
| OpenClaw Version | v2026.2.23 (b817600) |
| Capabilities | browser, system |
| Machine | MacBook (Zeus's primary Mac) |
| How to connect | `openclaw node run --host 127.0.0.1 --port 18789 --display-name Zeus-Mac` |
| Notes | Requires SSH tunnel active on port 18789 |

---

### Zeus's MacBook Pro

| Field | Value |
|-------|-------|
| Display Name | Zeus's MacBook Pro |
| Node ID | `50efbb0b87d9bf37bb732a000ad29aeef1fb42ffbdce9bf86c21394fc2722510` |
| Status | paired · disconnected (OpenClaw node not running) |
| Tailscale IP | `100.114.142.55` (hostname: `zeuss-macbook-pro`, relay: nyc) |
| OpenClaw Version | — |
| Capabilities | — |
| Machine | MacBook Pro (secondary Mac) |
| BlueBubbles | Hosts BlueBubbles server on `http://100.114.142.55:1234` (iMessage relay) |
| How to connect | `openclaw node run --host 127.0.0.1 --port 18789` (via SSH tunnel) |
| Notes | Paired but not currently running as a node |

---

### iPhone 17 Pro Max _(pending)_

| Field | Value |
|-------|-------|
| Display Name | iPhone 17 Pro Max |
| Node ID | _(assigned on first OpenClaw connection)_ |
| Status | Tailscale registered, last seen ~11h ago (offline — appears when OpenClaw app launched) |
| Tailscale IP | `100.122.203.68` (hostname: `iphone182`) |
| Platform | iOS |
| App | OpenClaw iOS (sideloaded) — source at `openclaw/ios/` |
| Gateway Config | Host: `100.104.61.29`, Port: `18789`, TLS: OFF, Token: `onemind-legacy-2026` |
| Notes | Must be on Tailscale. Trust developer cert: Settings → General → VPN & Device Management |

---

## Node Commands

```bash
# List all nodes with status
openclaw nodes status

# Pair a new node (run on the node machine, with tunnel active)
openclaw node pair --host 127.0.0.1 --port 18789 --token onemind-legacy-2026

# Run node (connect a paired machine)
OPENCLAW_GATEWAY_TOKEN=onemind-legacy-2026 openclaw node run --host 127.0.0.1 --port 18789 --display-name Zeus-Mac

# Remove a node
openclaw nodes remove <NODE_ID>
```

---

## Adding a New Node

1. Ensure SSH tunnel is up: `ssh -f -N -L 18789:127.0.0.1:18789 root@134.122.8.135`
2. On the new machine, run: `openclaw node pair --host 127.0.0.1 --port 18789 --token onemind-legacy-2026`
3. Accept the pairing prompt on the gateway
4. Update this file with the new node's ID, capabilities, and status

---

## Node Capability Reference

| Capability | Description |
|------------|-------------|
| `browser` | Can open URLs, take screenshots, web scraping |
| `system` | Shell commands, file access, process control |
| `camera` | Camera capture (iOS) |
| `screen` | Screen recording / sharing |
| `motion` | Accelerometer / motion sensors (iOS/Watch) |
| `location` | GPS location (iOS) |
| `calendar` | Calendar read/write (macOS/iOS) |
| `reminders` | Reminders access |
| `contacts` | Contacts read |
