# OMOS Personal Drone Swarm
**Domain:** SO-02 — Sovereign Infrastructure  
**Status:** Planning  
**Last Updated:** 2026-04-26  
**Author:** Legacy ⚡

---

## Vision

Custom-built personal drone swarm — facial recognition follow, hand gesture control, thermal + night vision, multi-drone formation. All controlled through Legacy via DiMOS. 3D printed chassis. Fully sovereign, no cloud.

---

## Core Architecture

```
Legacy (OpenClaw)
  └── DiMOS (orchestration layer)
        └── Swarm Manager (EVO-X2)
              ├── Drone 1 — Lead (RGB + Thermal + Night Vision)
              ├── Drone 2 — Guard follower
              ├── Drone 3 — Scout follower
              └── Drone N — Specialist (farm, perimeter, etc.)
                    Each drone:
                    ├── PX4 flight controller (Pixhawk 6C Mini)
                    ├── Jetson Orin Nano (on-drone AI brain)
                    ├── Camera array (per role)
                    └── MAVLink ↔ DiMOS ↔ Legacy
```

---

## Software Stack

### Flight Controller OS — PX4
- Best swarm support (leader-follower proven)
- Python SDK: MAVSDK (DiMOS already speaks MAVLink)
- Simulation: PX4 SITL + Gazebo (test everything before flying)
- Repo: https://github.com/PX4/PX4-Autopilot

### On-Drone AI — NVIDIA Jetson Orin Nano
- Runs facial recognition + gesture detection locally on each drone
- No cloud latency — real-time tracking
- Confirmed compatible: FLIR Lepton 3.5 thermal camera
- ~$200, low power, mounts on frame

### Facial Recognition + Follow Me
- **InsightFace** — best accuracy face recognition, runs on Jetson
- **OpenCV** + PX4 offboard mode = drone locks Zeus's face, follows
- Zeus enrolls face once → all swarm drones recognize and track
- Repo: https://github.com/deepinsight/insightface

### Hand Gesture Control
- **MediaPipe Hands** (Google) — runs on Jetson camera feed
- Custom gesture → MAVLink command mapping

| Gesture | Command |
|---|---|
| Open palm | Hover / hold position |
| Point forward | Move forward |
| Fist | Land |
| Two fingers up | Ascend |
| Two fingers down | Descend |
| Wave left | Strafe left |
| Wave right | Strafe right |
| Circle | Orbit Zeus |
| Thumbs up | Return home |

### Swarm Coordination
- PX4 + MAVSDK leader-follower mode
- Lead drone tracks Zeus via facial recognition
- Follower drones maintain formation offsets from lead
- DiMOS manages swarm via MCP → Legacy issues natural language commands
- Formations: triangle, line, orbit, spread, escort

---

## Hardware Per Drone

| Component | Part | Est. Cost |
|---|---|---|
| Flight controller | Pixhawk 6C Mini | ~$150 |
| AI companion computer | NVIDIA Jetson Orin Nano 8GB | ~$200 |
| RGB camera | Sony IMX477 (12MP, 4K) | ~$60 |
| Thermal camera | FLIR Lepton 3.5 | ~$200 |
| Night vision camera | Sony STARVIS IMX327 (low-light) | ~$80 |
| Frame | 3D printed custom chassis (PETG/CF) | ~$20 material |
| Motors + ESCs | 4× brushless + 4× ESC | ~$80 |
| Battery | 4S LiPo 3000mAh | ~$40 |
| RC / Comms | ExpressLRS (long range) | ~$30 |
| **Per drone total** | | **~$860** |

---

## Camera Configuration by Drone Role

```
Lead Drone (Zeus Follow)
  ├── RGB — daytime facial recognition + gesture detection
  ├── FLIR Lepton 3.5 — thermal (night / heat signatures)
  └── Sony STARVIS — low light / true night vision

Guard / Follower Drones
  ├── RGB — formation hold + situational awareness
  └── Thermal — perimeter/security coverage

Scout Drone
  ├── RGB — wide aerial overwatch
  └── Night vision — nocturnal ops

Farm / Property Drone
  ├── RGB — crop health, property inspection
  └── Thermal — irrigation leak detection, livestock heat maps
```

---

## Swarm Roles

| Drone | Role | Sensor Suite |
|---|---|---|
| **Drone 1 — Zeus Lead** | Face-lock, follow, gesture control | RGB + Thermal + NV |
| **Drone 2 — Guard** | Flank coverage, perimeter | RGB + Thermal |
| **Drone 3 — Scout** | Aerial overwatch, wide view | RGB + NV |
| **Drone 4+ — Farm/Property** | Crop + property monitoring per location | RGB + Thermal |

---

## OpenClaw Command Flow

```
Zeus: "Swarm up"
Legacy → DiMOS → MAVSDK → all drones arm + hover in formation

Zeus: "Follow me, triangle formation"
Legacy → DiMOS → Drone 1 locks face → Drones 2+3 hold triangle offsets

Zeus: "Scout the north field"
Legacy → DiMOS → Scout drone → autonomous waypoint → returns footage

Zeus: "Night patrol, perimeter"
Legacy → DiMOS → thermal drones → sweep property boundary → alert on heat signature

Zeus: "Land all"
Legacy → DiMOS → MAVSDK → all drones return home + land
```

---

## Budget

### Starter Swarm (3 drones)
| Config | Drones | Cost |
|---|---|---|
| Lead drone (full sensor suite) | 1 | ~$860 |
| Follower drones (RGB + Thermal) | 2 | ~$580 each |
| **Starter swarm total** | 3 | **~$2,020** |

### Full Swarm (5 drones)
| Config | Drones | Cost |
|---|---|---|
| Lead (full suite) | 1 | ~$860 |
| Guard/Follower | 2 | ~$580 each |
| Scout | 1 | ~$700 |
| Farm/Property | 1 | ~$580 |
| **Full swarm total** | 5 | **~$3,300** |

---

## Build Phases

### Phase 1 — Lead Drone (Prove the Concept)
- [ ] Print chassis (PETG, test fit)
- [ ] Assemble Pixhawk 6C Mini + Jetson Orin Nano
- [ ] Test PX4 SITL simulation (facial recognition + follow)
- [ ] Install InsightFace → enroll Zeus's face
- [ ] Implement MediaPipe gesture controls
- [ ] First real-world hover + face follow test

### Phase 2 — Swarm (Add Followers)
- [ ] Build 2 follower drones
- [ ] PX4 leader-follower MAVSDK setup
- [ ] Test formation hold + transitions
- [ ] Integrate DiMOS swarm manager on EVO-X2

### Phase 3 — Full Integration
- [ ] Connect swarm to Legacy via DiMOS MCP
- [ ] Natural language command routing
- [ ] Deploy farm/property specialist drones per location
- [ ] Night / thermal ops testing

---

## 3D Chassis Notes
- Material: **PETG** (impact resistant, light) or **Carbon Fiber PLA** for rigidity
- Arm design: foldable preferred for transport between properties
- Mount points: Jetson Nano + camera array + battery tray
- Size: 5" to 7" frame recommended for payload capacity
- Suggest designing in **Fusion 360** or **OpenSCAD** for parametric tuning

---

## Related Codex Files
- `SO/02-infrastructure/OMOS-Global-Infrastructure.md`
- `LE/50-homestead/Greenhouse-Automation.md` *(to create)*

## External References
- DiMOS: https://github.com/dimensionalOS/dimos
- PX4 Autopilot: https://github.com/PX4/PX4-Autopilot
- MAVSDK Python: https://github.com/mavlink/MAVSDK-Python
- InsightFace: https://github.com/deepinsight/insightface
- MediaPipe: https://github.com/google-ai-edge/mediapipe
- FLIR Lepton: https://www.flir.com/products/lepton/
- Pixhawk 6C Mini: https://holybro.com/products/pixhawk-6c-mini
- Jetson Orin Nano: https://developer.nvidia.com/embedded/jetson-orin-nano-devkit
