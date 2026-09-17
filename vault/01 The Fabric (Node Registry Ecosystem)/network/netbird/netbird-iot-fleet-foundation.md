---
title: "What Netbird Does for Your IoT Fleet"
created: 2026-05-23
tags: [infrastructure, networking, iot, netbird, sovereignty]
status: reference
pillar: SO
---

# What Netbird Does for Your IoT Fleet

## Overview

Netbird is a self-hosted WireGuard mesh VPN with built-in fleet management, access control, and public exposure. It replaces Tailscale as the sovereign networking foundation for OneMind OS.

**Repository:** https://github.com/netbirdio/netbird (BSD-3 License, 25k+ stars)
**Version:** v0.71.4 (May 2026)
**Self-hosted:** 100% — all features available, no cloud dependency

---

## Self-Hosted Architecture (MBP Clamshell)

```
MBP M2 Max (clamshell, always-on, Docker Desktop)
├── netbird-server        → Management + Signal server (gRPC)
├── netbird-dashboard     → Admin Web UI (port 80 behind Traefik)
├── netbird-proxy         → Reverse proxy for public HTTPS exposure
├── traefik               → TLS termination + routing
└── (optional) crowdsec   → DDoS/bot protection

Public access: Tailscale Funnel replacement
├── A record → MBP public IP (or DDNS)
├── Ports 80, 443 (TCP) + 3478 (UDP/STUN)
└── Traefik handles Let's Encrypt automatically
```

Docker containers on macOS:
- `netbird-server` — management + signal (all-in-one)
- `netbird-dashboard` — admin web UI
- `netbird-proxy` — reverse proxy service
- `traefik` — TLS/routing
- Total: ~4 containers, ~500MB RAM

---

## 1. IoT Hardware Support

Netbird builds static binaries for every IoT architecture:

| Architecture | Use Case |
|---|---|
| `linux/arm64` | Raspberry Pi 4/5, Jetson Nano, modern SBCs |
| `linux/arm` | Raspberry Pi Zero, older SBCs, drones |
| `linux/mips` (hardfloat + softfloat) | OpenWRT routers, cheap IoT gateways |
| `linux/mipsle` | MIPS little-endian routers |
| `linux/mips64` + `mips64le` | High-end routers, network appliances |
| `linux/386` | Legacy x86 embedded devices |
| `wasm` | Serverless / browser-based agents |
| OpenWRT package | Router-level mesh node |
| MikroTik RouterOS | Native integration |
| Synology, TrueNAS, Proxmox | NAS/hypervisor mesh nodes |
| pfSense, OPNsense | Firewall mesh nodes |

The MIPS builds cover commodity router silicon — any cheap OpenWRT gateway in a drone, robot, or field device runs Netbird natively. Single static binary, no dependencies.

---

## 2. Zero-Touch Provisioning (Setup Keys)

The killer feature for fleet scale:

```bash
# Flash firmware → device boots → joins mesh automatically
netbird up --setup-key YOUR_FLEET_KEY --management-url https://netbird.yourdomain.com:443
```

Setup key options:
- **Reusable** — one key for all devices of a type (all drones use the same key)
- **One-time** — single use, auto-expires after registration
- **Expiring** — time-limited enrollment window (e.g. 24h provisioning windows)
- **Auto-grouped** — devices auto-join a group on registration (e.g. `tag:drones`)
- **Ephemeral** — device removed from mesh when it goes offline (temp workers)

No manual approval, no GUI interaction. Device boots → joins mesh → inherits all ACL rules for its tag.

---

## 3. Access Control Per Device Group

```
Groups for OneMind OS:
├── tag:drones      → can reach telemetry + NATS only
├── tag:robots      → can reach NATS + control plane + SSH
├── tag:sensors     → read-only, publish to NATS, no outbound
├── tag:gateways    → full mesh routing, exit node capable
├── tag:agents      → MBP, agent nodes, full mesh
└── tag:admin       → Zeus's devices, unrestricted
```

Rules enforced at WireGuard level — cryptographic peer identity, not firewall rules. A compromised drone literally cannot reach services outside its group.

---

## 4. Posture Checks — Device Health Enforcement

Devices can be quarantined from the mesh if they fail checks:

| Check | What It Does | IoT Use Case |
|---|---|---|
| `NBVersionCheck` | Require minimum Netbird client version | Force firmware updates |
| `OSVersionCheck` | Require minimum OS version (per platform) | Block outdated kernels |
| `GeoLocationCheck` | Allow/deny by country code | Geofence devices |
| `PeerNetworkRangeCheck` | Enforce network range constraints | Detect stolen devices |
| `ProcessCheck` | Require specific process running | Ensure agent daemon is alive |

A robot running old firmware gets quarantined automatically until updated. No manual intervention.

---

## 5. Reverse Proxy — Public HTTPS Exposure

Replaces Tailscale Funnel with a more powerful alternative:

```
Public internet → Netbird Reverse Proxy (on MBP)
    ├── Auto TLS via Let's Encrypt (ACME)
    ├── Auth layer (configurable per route):
    │     ├── OAuth2/OIDC (SSO login)
    │     ├── Magic link (emailed one-time link)
    │     ├── PIN code (simple shared secret)
    │     ├── HTTP Basic Auth
    │     └── No auth (fully public APIs)
    └── WireGuard tunnel → target peer (any device on mesh)
```

Configured from the dashboard — no manual nginx/Caddy config files.

---

## 6. Private DNS Per Device Group

```
Drones    → resolve telemetry.onemind.local → 100.x.x.x
Robots    → resolve nats.onemind.local + control.onemind.local
Sensors   → no DNS access beyond their gateway
Agents    → resolve all *.onemind.local zones
```

Private DNS zones per group. Devices only resolve what they're permitted to reach.

---

## 7. Traffic Event Logging

Every peer connection is logged:
- Who connected to what
- When, for how long
- Bytes transferred
- Connection type (direct peer-to-peer or relayed)

Full audit trail across the fleet. Critical for sovereign operations.

---

## 8. Browser SSH & RDP

Access any device in the fleet from a browser:
- SSH terminal directly in the Netbird dashboard
- RDP to Windows devices
- No port forwarding, no SSH jump hosts

Centrally managed SSH policies — which users can SSH to which device groups.

---

## 9. Additional IoT Features

| Feature | Description |
|---|---|
| **Exit nodes** | Route device internet traffic through a specific peer |
| **Subnet routes** | Expose entire subnets (e.g. factory LAN) to the mesh |
| **IPv6 dual-stack** | Full IPv6 overlay support |
| **Quantum resistance** | Rosenpass post-quantum key exchange |
| **Lazy connections** | Devices don't hold persistent tunnels — connect on demand |
| **Periodic re-auth** | Force devices to re-authenticate on schedule |
| **Terraform provider** | Infrastructure-as-code for fleet configuration |
| **Ansible collection** | Automate client deployment at scale |
| **REST API** | Full programmatic control of all management functions |

---

## 10. Mobile Apps (iOS / Android)

- iOS app available on App Store — supports custom management URL for self-hosted
- Android app on Play Store — same capability
- Apple TV + Android TV clients available
- All support the same setup key provisioning

---

## Self-Hosted Feature Parity

Everything listed above is available self-hosted. Confirmed from source:
- Reverse proxy: open source (BSD-3), code in `proxy/` directory
- Posture checks: open source, code in `management/server/posture/`
- Dashboard: open source (separate repo: `netbirdio/dashboard`)
- iOS SDK: source in `client/ios/NetBirdSDK/`
- All client binaries: built from the same open-source repo

The only difference vs Netbird Cloud:
- You host the Docker containers yourself
- You manage TLS certificates (automated via Let's Encrypt)
- You provide an IdP (embedded Dex included in quickstart)

---

## OneMind OS Fleet Vision

```
┌─────────────────────────────────────────────────────────┐
│ MBP M2 Max (clamshell, always-on)                       │
│                                                         │
│  Docker Compose Stack:                                  │
│  ├── netbird-server (management + signal)               │
│  ├── netbird-dashboard (admin UI)                       │
│  ├── netbird-proxy (public HTTPS)                       │
│  ├── traefik (TLS)                                      │
│  ├── NATS server (agent backbone)                       │
│  └── Hermes (AI orchestrator)                           │
│                                                         │
│  Native services:                                       │
│  ├── Netbird client (MBP joins own mesh)                │
│  └── All existing services (Homepage, etc.)             │
└────────────────────┬────────────────────────────────────┘
                     │ WireGuard mesh
         ┌───────────┼───────────────────┐
         │           │                   │
    ┌────┴────┐ ┌────┴────┐      ┌──────┴──────┐
    │ iPhone  │ │ Drones  │      │ IoT Sensors │
    │ (iOS)   │ │ (arm64) │      │ (MIPS/ARM)  │
    └─────────┘ └─────────┘      └─────────────┘
```

---

## Migration Checklist

1. Install Docker Desktop on MBP (if not already)
2. Run Netbird self-hosted quickstart
3. Configure domain + DNS (DDNS or static)
4. Create setup keys per device group
5. Install Netbird client on MBP (joins own mesh)
6. Install Netbird on iPhone (App Store, custom URL)
7. Deactivate Tailscale on each device
8. Configure ACL rules per group
9. Set up reverse proxy routes for public services
10. Delete Tailscale account

---

## Key Consideration: Public IP Requirement

Self-hosted Netbird management needs to be reachable by all devices. Options for MBP in clamshell:

1. **Port forwarding** on home router (80, 443, 3478) → MBP local IP
2. **DDNS** service (e.g. DuckDNS, Cloudflare DDNS) for domain pointing
3. **Cloudflare Tunnel** as ingress (cloudflared on MBP, zero port-forward)
4. **Keep atlas temporarily** just as a tiny reverse proxy to MBP (cheapest Lightsail)

The relay/STUN server (coturn) on port 3478 UDP is needed for NAT traversal when devices can't connect peer-to-peer directly.
