# OMOS Global Infrastructure Architecture
**Domain:** SO-02 — Sovereign Infrastructure  
**Status:** Planning  
**Last Updated:** 2026-04-26 (rev 3 — DiMOS drone support added)  
**Author:** Legacy ⚡

---

## Vision

One central brain at the Homestead. Every location managed as an autonomous node reporting back to the primary gateway. All under One Mind OS.

---

## Core Principle

```
One Brain (Homestead) → Tailscale Mesh → All Locations
Each location = autonomous node + reports up to central gateway
```

---

## Three-Tier Node Model

### TIER 1 — Central Brain (Homestead)
- **Hardware:** GMKtec EVO-X2 (128GB RAM, Ryzen AI Max+ 395)
- OpenClaw Primary Gateway
- Ollama local LLMs (70B+ capable)
- **DiMOS** — AI robotics OS (github.com/dimensionalOS/dimos)
- Master Home Assistant instance (aggregates all locations)
- Tailscale hub / subnet router
- All agent sessions live here
- Meshtastic base node (USB LoRa radio)

### TIER 2 — Smart Nodes (Full Local Compute)
**Locations:** NYC Office | Greenwich House | Warehouse  
- **Hardware:** Mini PC — Beelink/Minisforum 16–32GB (~$300–500)
- OpenClaw Node → reports to Homestead gateway
- Local Home Assistant instance
- Frigate NVR (cameras, AI object detection — fully local)
- Tailscale connected

### TIER 3 — Edge Nodes (Lightweight, Remote)
**Locations:** Mexico Farm | DR Property | NY Farm | Barn  
- **Hardware:** Raspberry Pi 5 8GB (~$80)
- OpenClaw Node → reports to Homestead gateway
- Local Home Assistant instance
- Frigate (where cameras are needed)
- Starlink (Mexico + DR) → Tailscale
- Meshtastic LoRa for offline fallback (farms + barn)

---

## Location Matrix

| Location | Tier | Local Compute | Connectivity | Primary Use |
|---|---|---|---|---|
| **Homestead** | Central Brain | EVO-X2 128GB | Fiber/LTE | Gateway, AI, orchestration, greenhouse |
| **NYC Office** | 2 | Mini PC (16–32GB) | Fiber | Cameras, access control, work ops |
| **Greenwich House** | 2 | Mac Mini M4 / Mini PC | Fiber/Cable | Smart home, security |
| **Warehouse** | 2 | Mini PC (16GB) | Fiber/LTE | Inventory, cameras, access control |
| **Barn (Homestead area)** | 3 | Raspberry Pi 5 | LTE / Meshtastic | Sensors, animal monitoring |
| **NY Farm** | 3 | Raspberry Pi 5 | LTE / Meshtastic | Soil sensors, irrigation, weather |
| **Mexico Farm** | 3 | Raspberry Pi 5 | Starlink | Remote sensors, cameras, irrigation |
| **DR Property** | 3 | Raspberry Pi 5 | Starlink | Security, environment monitoring |

---

## Network Backbone

### Tailscale Zero-Trust Mesh
- Every node joins the existing Tailscale network
- All traffic encrypted end-to-end
- Works through NAT, Starlink, LTE — no port forwarding needed
- Homestead EVO-X2 acts as subnet router for all local devices
- Already operational on VPS and HA node — extend same pattern

### Starlink
- Required for: Mexico Farm, DR Property
- Provides reliable internet for Tailscale connectivity
- Hardware: ~$599 setup + ~$120/month per location

### Meshtastic LoRa Mesh
- Offline fallback for: Barn, NY Farm, Homestead Greenhouse
- No internet required — commands still reach local nodes
- Status pings back when online window occurs
- Hardware: Heltec V3 or LILYGO T-Beam (~$30–50 each)

---

## What Each Node Manages Locally

| System | Tool | Notes |
|---|---|---|
| **Cameras** | Frigate NVR | AI detection local, clips stored local, feeds to HA |
| **Sensors** | HA + ESPHome | Soil, temp, humidity, motion, door/window |
| **Access Control** | HA + smart locks / RFID | Gates, doors, entry points |
| **Energy** | HA energy monitoring | Per-location usage tracking |
| **Irrigation** | HA + relay boards | Automated by soil moisture thresholds |
| **Alerts** | Local node → Tailscale → Homestead → Legacy → Telegram | Unified alert routing |

---

## Robotics Layer — DiMOS

**Repo:** https://github.com/dimensionalOS/dimos  
**Role:** AI-native robotics OS that bridges Legacy → Ollama → physical robots  
**Key features:**
- Native Ollama integration (local LLM → robot agent pipeline)
- MCP support (Legacy can issue natural language commands directly)
- Supports Unitree Go2 quadruped, Unitree G1 humanoid, xArm7 manipulation arm
- **Drone support:** MAVLink (generic UAVs), DJI Mavic — natural language control
- Simulation mode via MuJoCo (test before deploying on real hardware)
- Camera + sensor stream pipelines
- WebRTC control for real robots over network

### Command Flow with DiMOS

```
Legacy (OpenClaw)
  └── natural language: "Plant row 3, greenhouse"
        └── DiMOS (EVO-X2)
              └── Ollama local LLM (interprets + plans)
                    └── Robot agent pipeline
                          └── xArm7 (planting) or Go2 (patrol)
                                └── Executes physically — no cloud
```

### Robot Hardware Options

| Robot / Platform | Use Case | Est. Cost | DiMOS Status |
|---|---|---|---|
| **Unitree Go2** | Farm/barn patrol, inspection, security | ~$1,600 | Beta |
| **Unitree G1** | Humanoid general labor (future) | ~$16,000 | Beta |
| **xArm7** | Greenhouse planting, manipulation | ~$5,000 | Beta |
| **xArm Lite 6** | Lighter manipulation tasks | ~$2,000 | Beta |
| **DJI Mavic** | Aerial farm/property surveillance | ~$1,000–2,500 | Alpha |
| **MAVLink drone** | Custom UAV (crop monitoring, delivery) | ~$500–1,500 | Alpha |

---

## Greenhouse Sub-System (Homestead)

Dedicated robotics + automation layer at the homestead greenhouse:

```
EVO-X2 (main house)
  ├── OpenClaw + Legacy
  ├── Ollama (local LLM)
  ├── DiMOS (robot agent OS)
  ├── Home Assistant
  └── Meshtastic base
        └── [LoRa Mesh]
              └── Raspberry Pi 5 (greenhouse)
                    └── HA + sensors (soil, temp, humidity, light)
                    └── Actuators (pumps, grow lights, vents)
                    └── xArm7 (planting automation)
                    └── Go2 (patrol + inspection)
                    └── DJI Mavic / MAVLink drone (aerial surveillance)
                    └── ESP32 nodes (field sensors, solar-powered)
```

---

## Command Flow (How Legacy Controls Everything)

```
Zeus: "Status on Mexico farm"
Legacy → queries Mexico Pi node via Tailscale
       → pulls HA sensor data
       → returns: temp, soil moisture, cameras online/offline

Zeus: "Lock the DR property"
Legacy → sends command to DR Pi node
       → HA executes lock automation
       → confirms back to Zeus

Zeus: "Start irrigation on NY Farm"
Legacy → Homestead EVO-X2 → Tailscale → NY Farm Pi
       → HA triggers relay → water pump ON
       → confirms
```

---

## Hardware Budget (Estimate)

### Infrastructure
| Component | Qty | Unit Cost | Total |
|---|---|---|---|
| GMKtec EVO-X2 128GB (Homestead) | 1 | $1,800–2,000 | $2,000 |
| Mini PC Tier 2 nodes (Office, Greenwich, Warehouse) | 3 | $350 | $1,050 |
| Raspberry Pi 5 8GB Tier 3 nodes (Barn, NY Farm, Mexico, DR) | 4 | $80 | $320 |
| Starlink setup (Mexico + DR) | 2 | $599 | $1,200 |
| Starlink monthly (Mexico + DR) | 2 | $120/mo | $240/mo |
| Meshtastic radios (Homestead, Barn, NY Farm, Greenhouse) | 6 | $40 | $240 |
| Sensors, relays, misc per location | — | — | ~$500 |
| **Infrastructure Total** | | | **~$5,300** |

### Robotics (Homestead Phase)
| Component | Qty | Unit Cost | Total |
|---|---|---|---|
| Unitree Go2 (farm/barn patrol) | 1 | $1,600 | $1,600 |
| xArm Lite 6 (greenhouse planting) | 1 | $2,000 | $2,000 |
| DJI Mavic 3 (aerial surveillance) | 1 | $2,200 | $2,200 |
| DiMOS software | — | Open source | $0 |
| **Robotics Total (starter)** | | | **~$5,800** |

| **Grand Total Hardware** | | | **~$11,100** |
| **Monthly recurring (Starlink ×2)** | | | **~$240/mo** |

---

## Phased Rollout

### Phase 1 — Homestead Core (Now)
- [ ] Order GMKtec EVO-X2 128GB
- [ ] Set up OpenClaw primary gateway
- [ ] Migrate from VPS to EVO-X2 or run hybrid
- [ ] Install DiMOS + Ollama on EVO-X2
- [ ] Greenhouse Pi 5 node + Meshtastic
- [ ] Test DiMOS in simulation (MuJoCo) before robot hardware

### Phase 2 — NYC + Greenwich (Month 1–2)
- [ ] Mini PC nodes at NYC Office + Greenwich House
- [ ] Frigate NVR at both locations
- [ ] Tailscale connected, OpenClaw nodes paired

### Phase 3 — Warehouse + Barn + NY Farm (Month 2–3)
- [ ] Pi 5 nodes deployed
- [ ] Meshtastic mesh for offline fallback
- [ ] HA automation per location

### Phase 4 — International (Month 3–6)
- [ ] Starlink ordered for Mexico + DR
- [ ] Pi 5 nodes shipped and deployed
- [ ] Remote management verified via Tailscale

---

## Open Questions

- VPS (onemindos-cloud-1): Keep as backup/cloud gateway or decommission once EVO-X2 is live?
- Gateway hybrid: Run EVO-X2 as primary + VPS as failover?
- Mac Mini vs EVO-X2: Revisit if Mac Studio M4 Max pricing becomes competitive (~$2,000)
- Solar power: For farms + barn nodes — Pi 5 draws ~5–10W, very viable on small solar panel

---

## EagleEye — Wearable Sensor Fusion HUD (Summary)

Civilian equivalent of Anduril EagleEye. AR glasses or custom helmet with thermal + UWB radar + night vision fused into one display. See through walls, see in the dark, command drones with hand gestures, Legacy whispers intel in real time.

**Stack:** Xreal Air 2 Pro (display) → Jetson Orin NX (on-body compute) → FLIR Boson+ (thermal) + Acconeer XM125 UWB (through-wall radar) + Sony STARVIS (night vision) → OpenXR + ROS2 + InsightFace → Legacy  
**Total hardware:** ~$2,059  
**Full doc:** `SO/02-infrastructure/OMOS-EagleEye.md`

---

## Personal Drone Swarm (Summary)

Custom-built drone swarm — facial recognition follow, hand gesture control, thermal + night vision, multi-drone formations. All commanded via Legacy through DiMOS. 3D printed chassis.

**Stack:** PX4 + MAVSDK + Jetson Orin Nano (per drone) + InsightFace + MediaPipe Hands + DiMOS  
**Lead drone:** RGB + FLIR Lepton 3.5 (thermal) + Sony STARVIS (night vision) — ~$860  
**Starter swarm (3 drones):** ~$2,020 | **Full swarm (5 drones):** ~$3,300  
**Full doc:** `SO/02-infrastructure/OMOS-Drone-Swarm.md`

---

## Related Codex Files
- `SO/02-infrastructure/VPS-Setup.md`
- `SO/02-infrastructure/OMOS-Drone-Swarm.md`
- `SO/02-infrastructure/OMOS-EagleEye.md`
- `SO/05-automation/HomeAssistant-Node.md`
- `LE/50-homestead/Greenhouse-Automation.md` *(to create)*
- `LE/50-homestead/Robotics-DiMOS.md` *(to create)*

## External References
- DiMOS repo: https://github.com/dimensionalOS/dimos
- Unitree Go2: https://www.unitree.com/go2
- xArm: https://www.ufactory.cc
- PX4 Autopilot: https://github.com/PX4/PX4-Autopilot
- InsightFace: https://github.com/deepinsight/insightface
