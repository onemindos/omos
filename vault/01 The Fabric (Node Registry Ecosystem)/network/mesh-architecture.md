---
title: "OneMind OS — Mesh Network Architecture"
type: reference
domain: SO
status: active
created: 2026-06-02
updated: 2026-06-02
tags: [networking, mesh, headscale, tailscale, wireguard, nats, meshtastic, iot, fleet]
---

# OneMind OS — Mesh Network Architecture

## Decision

**Headscale** (self-hosted coordination server) + **Tailscale clients** (on every node).

- Control plane: Headscale running on MBP (24/7 clamshell)
- Data plane: P2P WireGuard between all nodes
- Clients: Official Tailscale apps on iOS, Android, Linux, macOS
- Domain: onemindos.dev (permanent, owned, never changes)

---

## Mesh Layer (IP — Headscale/Tailscale)

All nodes that run Linux, macOS, iOS, or Android join this mesh via Tailscale client.

```
HEADSCALE MESH (onemindos.dev)
├── omos-mbp          MBP M2 Max — control center, Headscale host, subnet router
├── omos-dell         Dell workstation
├── omos-ha           Home Assistant — Zigbee/Z-Wave bridge
├── barn-node         Raspberry Pi — barn sensors
├── office-node       Raspberry Pi — office node
├── dimos             Robot (Linux onboard)
├── om1               Robot (Linux onboard)
├── g1                Robot (Linux onboard) — future
├── go2               Robot (Linux onboard) — future
├── recon-drone       Companion PC on drone (Linux)
└── greenhouse-node   Raspberry Pi — soil/climate sensors
```

MagicDNS hostname pattern: `<node>.onemindos.dev`
Subnet router on MBP exposes: `192.168.1.0/24` → entire LAN accessible from any mesh node

---

## Off-Mesh Layer (feeds INTO mesh nodes)

Devices that cannot run WireGuard — no OS, no Linux, raw firmware.
These report to a mesh node which bridges their data to NATS.

```
OFF-MESH                          BRIDGE PATH
───────────────────────────────────────────────────────────────
Meshtastic (LoRa radios)          BT/Serial → Pi bridge node → NATS
ESP32 / Arduino sensors           WiFi → NATS directly (no bridge needed)
Zigbee / Z-Wave devices           → Home Assistant → mesh → NATS
Pixhawk / ArduPilot (FC)          MAVLink → drone companion PC → mesh
Raw IoT sensors (soil, climate)   → Pi node → NATS
```

---

## Why These Devices Can't Run Headscale/Tailscale

| Device | Reason |
|---|---|
| Meshtastic | ESP32/nRF52 — bare metal LoRa firmware, no OS |
| ESP32/Arduino | No OS, no WireGuard stack |
| Zigbee/Z-Wave | Not IP-based, RF protocol |
| Pixhawk/ArduPilot | Bare metal flight controller firmware |
| LoRa sensors | Low-power embedded, no TCP/IP |

---

## Meshtastic Integration

Meshtastic and Headscale operate at different layers — they complement, not compete.

```
LAYER           TECHNOLOGY          PURPOSE
─────────────────────────────────────────────────────────────
Radio mesh      Meshtastic (LoRa)   Off-grid, no-internet comms
                                    915MHz, long range, low power
                                    Text messages, GPS position, telemetry

IP mesh         Headscale/Tailscale WireGuard P2P between Linux/mobile nodes
                                    Full TCP/IP, encrypted tunnel

Bridge          Raspberry Pi        Runs both:
                                    - Meshtastic Python client (BT/Serial)
                                    - Tailscale client (mesh node)
                                    - NATS publisher (routes LoRa data to fleet)
```

NATS subject pattern for Meshtastic data:
- `mesh.position.<node_id>` — GPS coordinates
- `mesh.telemetry.<node_id>` — battery, environment sensors
- `mesh.message.<channel>` — text messages

---

## Drone Integration

```
DRONE STACK
├── Pixhawk / ArduPilot (flight controller — bare metal)
│     ↕ MAVLink (serial/UDP)
├── Companion PC (Jetson Nano / RPi / Orange Pi — Linux)
│     ├── Tailscale client → Headscale mesh node (recon-drone.onemindos.dev)
│     ├── NATS publisher → telemetry, video streams, commands
│     └── MAVLink bridge → routes FC commands over mesh
│           ↕ WireGuard tunnel over mesh
└── MBP running QGroundControl
      └── Talks MAVLink to companion PC over encrypted mesh tunnel
```

No public IP needed. No port forwarding. Works from anywhere on the mesh.

---

## Robot Integration

All robots (DIMOS, OM1, G1, Go2) run Linux onboard computers.

```
ROBOT NODE
├── Tailscale client → mesh (dimos.onemindos.dev, etc.)
├── NATS client → publishes telemetry, subscribes to commands
└── ROS2 / custom stack → runs onboard

NATS subject pattern:
  cmd.<robot>.<action>         → commands to robot
  <robot>.<domain>.<sensor>    → telemetry from robot
  <robot>.status               → health / state

Examples:
  cmd.dimos.move_forward
  dimos.body.imu
  go2.leg.contact_sensors
  recon.camera.frame
```

---

## NATS as the Integration Bus

NATS bridges the IP mesh layer and the off-mesh layer.
Every sensor, robot, and service publishes and subscribes via NATS subjects.

```
NATS (127.0.0.1:4222, account: GLOBAL/onemind)
  ├── All mesh nodes have NATS client access
  ├── ESP32 sensors publish directly via WiFi → NATS
  ├── Meshtastic bridge Pi publishes LoRa data → NATS
  ├── Robots publish telemetry → NATS
  ├── Hermes / OpenClaw subscribe to relevant subjects
  └── Node-RED orchestrates IoT flows via NATS
```

---

## Subnet Routing

MBP runs as subnet router — exposes the entire home LAN to mesh nodes.

```
sudo tailscale up --advertise-routes=192.168.1.0/24
```

From iPhone, Android, or any mesh node:
- `omos-mbp.onemindos.dev:9119` → Hermes dashboard
- `omos-mbp.onemindos.dev:1880` → Node-RED
- `omos-mbp.onemindos.dev:18789` → OpenClaw gateway
- `omos-ha.onemindos.dev:8123` → Home Assistant
- Any `192.168.1.x` service → accessible via subnet route

---

## Why Headscale Over Tailscale Managed

| Criteria | Tailscale Managed | Headscale |
|---|---|---|
| Cost at 50+ nodes | Paid per device | Free forever |
| Domain URLs | tail717ba0.ts.net (theirs) | onemindos.dev (yours) |
| Vendor lock-in | High | Zero |
| If company shuts down | Mesh dies | Nothing changes |
| Control | Their rules | Your rules |
| iOS/Android | Same Tailscale app | Same Tailscale app |
| Dashboard | Polished | Headplane (community) |
| Maintenance | Zero | Minimal (1 container/binary) |
| Funnel (public URLs) | Yes | No |

Decision: Headscale wins at scale for a private fleet with permanent infrastructure.

---

## What Headscale Cannot Do (vs Tailscale Managed)

- Tailscale Funnel — expose local service to public internet without port forward
- Tailscale Serve — local reverse proxy feature
- Mullvad VPN integration
- Node sharing across organizations
- Official support / SLA

None of these are required for OneMind OS private fleet use.

---

## Current State (2026-06-02)

- Running: NetBird self-hosted (fragile — migrating to Headscale)
- NetBird peers: omos-dell (connected), omos-iphone (disconnected 15h), omos-ha (disconnected 15h)
- Root cause: STUN broken, relay address points to LAN IP not accessible externally
- Migration target: Headscale binary on MBP + Headplane dashboard + Tailscale clients

---

## Related

- [[SO/03 Infrastructure/robots/]] — individual robot node docs
- [[SO/01 The Fabric/infrastructure/networking/]] — netbird.md, tailscale.md
- [[SO/03 Infrastructure/services/nats.md]] — NATS server config
