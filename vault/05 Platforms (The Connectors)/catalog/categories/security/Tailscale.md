# Tailscale

> **Zero-config WireGuard VPN mesh — secure private network for all OneMind devices**

## Overview

Tailscale creates a private mesh network between all devices without port forwarding or complex VPN setup. OneMind runs a 10-node tailnet under the `Zeus-Delacruz@` account. All production Lightsail servers, workstations, and mobile devices are members.

## Status

| Attribute | Value |
|-----------|-------|
| Status | ✅ Active |
| Version | 1.96.5 |
| Tailnet | Zeus-Delacruz@ |
| Category | Security / Networking |
| Type | Zero-Config VPN (WireGuard) |
| Integration | CLI, macOS/iOS/Linux app |
| Last checked | 2026-05-08 |

## Use Cases in OneMind

1. **VPS Access** — SSH to Lightsail nodes without public key distribution
2. **Internal Services** — Access backend services without exposing ports
3. **Mobile Access** — iPhone/iPad connect to VPS services from anywhere
4. **Inter-agent Communication** — OpenClaw agents on VPS reachable from Mac

## Network

Full map → [[TAILSCALE-NETWORK]]

| Node | Tailscale IP | Role |
|------|-------------|------|
| zeus-macbook-pro | 100.114.142.55 | Primary workstation |
| omos-core | 100.123.126.34 | Lightsail — orchestration |
| omos-sentinel | 100.66.235.38 | Lightsail — sentinel |
| omos-titan | 100.103.197.23 | Lightsail — titan |
| omos-dell | 100.66.84.71 | Linux workstation |
| zeus-iphone-17 | 100.122.203.68 | Mobile |
| omos-ha | 100.64.25.12 | Home Assistant (offline ⚠️) |

## Quick Start

```bash
# macOS — full path (always works without setup)
/Applications/Tailscale.app/Contents/MacOS/Tailscale status

# Add wrapper script so `tailscale` works anywhere (run once, no sudo)
# NOTE: a symlink fails with bundle ID error — use wrapper script instead
printf '#!/bin/sh\nexec /Applications/Tailscale.app/Contents/MacOS/Tailscale "$@"\n' \
  > /opt/homebrew/bin/tailscale && chmod +x /opt/homebrew/bin/tailscale

# SSH to production nodes
ssh ubuntu@100.123.126.34   # omos-core
ssh ubuntu@100.66.235.38    # omos-sentinel
```

## Links

- Admin panel: https://login.tailscale.com/admin/machines
- ACLs: https://login.tailscale.com/admin/acls
- Docs: https://tailscale.com/kb

---

## Related

- [[TAILSCALE-NETWORK]] — Full network map and fleet table
- [[TAILSCALE-SSH]] — Keyless SSH runbook
- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
