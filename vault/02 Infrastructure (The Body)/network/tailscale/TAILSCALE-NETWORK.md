# OneMind OS — Tailscale Network

> **Last Updated**: 2026-05-08
> **Status**: Production (Tailscale-only)
> **Version**: 1.96.5

## Overview

OneMind OS uses **Tailscale** as the sole networking layer across all servers, workstations, and mobile devices. All traffic is encrypted via WireGuard. No public ports, no DNS complexity.

## Network Map

```
┌─────────────────────────────────────────────────────────────────┐
│                    TAILSCALE MESH (zeus-delacruz)               │
│                                                                 │
│  WORKSTATIONS                  LIGHTSAIL (Production)           │
│  ┌──────────────────────┐      ┌──────────────────────────┐     │
│  │ zeus-macbook-pro     │─────►│ omos-core                │     │
│  │ 100.114.142.55       │      │ 100.123.126.34           │     │
│  └──────────────────────┘      │ 44.193.145.121           │     │
│                                └──────────────────────────┘     │
│  ┌──────────────────────┐      ┌──────────────────────────┐     │
│  │ omos-dell            │      │ omos-sentinel            │     │
│  │ 100.66.84.71         │      │ 100.66.235.38            │     │
│  └──────────────────────┘      │ 3.214.27.142             │     │
│                                └──────────────────────────┘     │
│  MOBILE                        ┌──────────────────────────┐     │
│  ┌──────────────────────┐      │ omos-titan               │     │
│  │ zeus-iphone-17       │      │ 100.103.197.23           │     │
│  │ 100.122.203.68       │      │ 44.214.141.64            │     │
│  └──────────────────────┘      └──────────────────────────┘     │
│                                                                 │
│  OFFLINE / INVESTIGATE                                          │
│  omos-ha (100.64.25.12) — unreachable as of 2026-05-08         │
└─────────────────────────────────────────────────────────────────┘
```

## Fleet

| Node | Tailscale IP | Public IP | Platform | Status |
|------|-------------|-----------|----------|--------|
| zeus-macbook-pro | 100.114.142.55 | — | macOS | ✅ active (this machine) |
| omos-core | 100.123.126.34 | 44.193.145.121 | Lightsail / Ubuntu | ✅ active, direct |
| omos-sentinel | 100.66.235.38 | 3.214.27.142 | Lightsail / Ubuntu | ✅ active, direct |
| omos-titan | 100.103.197.23 | 44.214.141.64 | Lightsail / Ubuntu | ✅ connected |
| omos-dell | 100.66.84.71 | — | Linux workstation | ✅ connected |
| zeus-iphone-17 | 100.122.203.68 | — | iOS | ✅ connected |
| omos-ha | 100.64.25.12 | — | Home Assistant OS | ⚠️ offline — investigate |
| omos-pad | 100.67.182.83 | — | Android | 💤 offline 10d |
| omos-tab | 100.81.60.35 | — | Android | 💤 offline 10d |
| omos-tv | 100.117.247.62 | — | Android | 💤 offline 13d |

## CLI

The Tailscale binary is not in `$PATH` on macOS by default. A symlink fails because the app bundle validates its own bundle ID when called from outside the bundle. Use a wrapper script instead:

```bash
# Full path (always works without setup)
/Applications/Tailscale.app/Contents/MacOS/Tailscale status

# Add wrapper script so `tailscale` works anywhere (run once, no sudo needed)
printf '#!/bin/sh\nexec /Applications/Tailscale.app/Contents/MacOS/Tailscale "$@"\n' \
  > /opt/homebrew/bin/tailscale && chmod +x /opt/homebrew/bin/tailscale
```

> **Note:** A symlink (`ln -sf`) looks correct but crashes with `Fatal error: The current bundleIdentifier is unknown to the registry`. Always use the wrapper script.

## Common Commands

```bash
# Network status
tailscale status

# Ping a node
tailscale ping omos-core

# This machine's Tailscale IP
tailscale ip

# SSH to a Lightsail node via Tailscale
ssh ubuntu@100.123.126.34   # omos-core
ssh ubuntu@100.66.235.38    # omos-sentinel
ssh ubuntu@100.103.197.23   # omos-titan
```

## Connectivity Check

```bash
# Quick reachability sweep
for node in 100.123.126.34 100.66.235.38 100.103.197.23; do
  /Applications/Tailscale.app/Contents/MacOS/Tailscale ping $node | tail -1
done
```

## Troubleshooting

| Issue | Fix |
|-------|-----|
| `command not found: tailscale` | Add wrapper: `printf '#!/bin/sh\nexec /Applications/Tailscale.app/Contents/MacOS/Tailscale "$@"\n' > /opt/homebrew/bin/tailscale && chmod +x /opt/homebrew/bin/tailscale` |
| Node shows offline | `tailscale ping <ip>` to confirm — if timeout, node is down |
| `omos-ha` unreachable | Local LAN device, no public IP — must be on same network or HA restarted |
| Relay instead of direct | Both peers need UDP 41641 open; check firewall |

## Security

- WireGuard encryption on all traffic
- Zero trust — only tailnet members can reach any node
- No port 22 or any service port exposed to the public internet
- Access control via ACLs: https://login.tailscale.com/admin/acls

---

**Maintained by**: Zeus
**Related**:
- [[TAILSCALE-SSH]] — Keyless SSH across the fleet
- [[02 Infrastructure]] — Full infrastructure overview
