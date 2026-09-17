# OMOS EagleEye — Wearable Sensor Fusion HUD
**Domain:** SO-02 — Sovereign Infrastructure  
**Status:** Planning  
**Last Updated:** 2026-04-26  
**Author:** Legacy ⚡  
**Inspired by:** Anduril EagleEye (https://www.anduril.com/eagleeye)

---

## Vision

A civilian sovereign equivalent of Anduril's EagleEye — wearable AR glasses or custom helmet with thermal, night vision, and UWB radar sensor fusion. See through walls, see in the dark, command the drone swarm with hand gestures, and have Legacy whispering intel in real time. All connected to OMOS.

---

## What EagleEye (Anduril) Does

Anduril (Palmer Luckey — Oculus founder) built EagleEye by fusing:
- Night vision + Thermal + AR overlay + RF/acoustic sensors + AI targeting
- Rendered into a single helmet HUD
- Connected to Lattice (their battlefield OS) — same concept as OpenClaw/Legacy
- Partners: Meta (AR display), Oakley (eyewear), Qualcomm (compute), Gentex (ballistics)

The "see through walls" capability = sensor fusion:
- UWB radar detects movement through walls
- Thermal sees heat signatures through thin walls
- All overlaid in AR in real time

---

## OMOS EagleEye Architecture

```
OMOS EagleEye (Wearable)
  ├── Display Layer     → AR glasses / custom helmet HUD
  ├── Sensor Layer      → Thermal + UWB Radar + RGB + Night Vision
  ├── Compute Layer     → Jetson Orin NX (on-body, fuses all sensors)
  ├── AI Layer          → Legacy (whispers intel, facial ID, object detection)
  └── Network Layer     → Tailscale → EVO-X2 → all drones + properties
```

---

## Display Options

| Option | Device | Cost | Notes |
|---|---|---|---|
| **Glasses (recommended dev start)** | Xreal Air 2 Pro | ~$449 | Most open AR SDK, developer-friendly |
| **Glasses (Meta ecosystem)** | Meta Ray-Ban Display | ~$799 | SDK just released, Meta AI native |
| **Custom helmet** | 3D printed (PETG/CF) | ~$200 | Max sensors, full integration, ballistic option |
| **Monocle** | OpenGlass (open source) | ~$100 | Lightweight, hackable, minimal |

**Build path:** Start Xreal Air 2 Pro for dev/testing → custom helmet for full sensor suite deployment.

---

## See-Through-Wall Sensor Stack

| Sensor | Capability | Range | Cost |
|---|---|---|---|
| **FLIR Boson+** | Heat signatures through drywall/thin walls | ~50m | ~$800 |
| **Acconeer XM125** (UWB Radar) | Movement detection through walls + concrete | ~10m | ~$50 |
| **Novelda X4** (UWB Radar) | Breathing / heartbeat detection through walls | ~5m | ~$80 |
| **Sony STARVIS IMX327** | Near-IR night vision, low light | — | ~$80 |
| **Sony IMX477** | RGB daylight camera | — | ~$60 |

**Fusion logic:**
- Thermal → heat body signatures
- UWB radar → confirms movement + distance behind wall
- Night vision → full scene in low light
- All fused into single AR overlay on display

---

## On-Body Compute

**NVIDIA Jetson Orin NX 16GB** (~$400)
- Worn on body (vest, belt mount, or backpack)
- Runs all sensor fusion in real time
- Facial recognition on all incoming camera feeds (InsightFace)
- Streams drone feeds into HUD (picture-in-picture)
- Tailscale connected → Legacy always accessible via EVO-X2
- Runs OpenXR / ROS2 sensor fusion pipeline

---

## What Zeus Sees in the HUD

```
[ Looking at a wall ]
  └── UWB radar overlay → glowing outline of person behind wall
  └── Thermal overlay   → heat signature, distance estimate
  └── AI tag (Legacy)   → "Unknown male, ~6ft, stationary, 3m"

[ Walking the farm at night ]
  └── Night vision       → full scene rendered
  └── Thermal overlay    → animals, humans, heat anomalies tagged
  └── Drone feed (PiP)   → Drone 2 aerial view bottom corner
  └── Legacy whisper     → "North perimeter clear. Drone 3 scanning east field."

[ Approaching the gate ]
  └── Facial recognition → tags known people (family, workers)
  └── Unknown person     → Legacy alerts, drone auto-dispatched
  └── Hand gesture       → command drone swarm without speaking

[ Property overview mode ]
  └── Map overlay        → all 8 properties shown
  └── Active alerts      → blinking pins on anomalies
  └── Sensor status      → all HA nodes live/offline indicators
```

---

## OpenClaw / Legacy Integration

### Legacy → EagleEye HUD
- Property alerts pushed to display in real time
- Drone feeds routed to HUD on demand
- Voice commands processed by Legacy → action executed
- Facial recognition results overlaid on detected faces
- Swarm status, battery levels, positions shown in HUD

### EagleEye → Legacy
- What Zeus sees is logged (thermal anomalies, intruders, events)
- Sensor data feeds back to Home Assistant (property awareness)
- Drone swarm commanded via hand gestures → DiMOS
- Legacy learns Zeus's patrol patterns and routes

### Example Command Flow
```
Zeus: "Legacy, who's at the gate?"
Legacy → pulls gate camera + facial rec → overlays result in HUD
       → "That's Marco, arrived 2 mins ago."

Zeus (hand signal: circle)
DiMOS → Drone 1 begins orbit pattern around Zeus

Zeus: "Legacy, thermal sweep, barn"
Legacy → routes barn drone to thermal mode → feed appears in HUD corner
       → flags any heat anomalies
```

---

## Full Hardware BOM

| Component | Part | Est. Cost |
|---|---|---|
| AR display | Xreal Air 2 Pro | ~$449 |
| On-body compute | Jetson Orin NX 16GB | ~$400 |
| Thermal camera | FLIR Boson+ | ~$800 |
| UWB radar ×2 | Acconeer XM125 | ~$100 |
| UWB heartbeat | Novelda X4 | ~$80 |
| Night vision camera | Sony STARVIS IMX327 | ~$80 |
| RGB camera | Sony IMX477 | ~$60 |
| Wearable battery | 20,000mAh slim pack | ~$60 |
| Mounting / housing | 3D printed (vest + helmet mounts) | ~$30 |
| **Total** | | **~$2,059** |

---

## Software Stack

| Layer | Tool | Notes |
|---|---|---|
| AR display SDK | OpenXR + Xreal SDK | Open standard, cross-platform |
| Sensor fusion | ROS2 | Industry standard robotics middleware |
| Facial recognition | InsightFace | Runs on Jetson, same as drone stack |
| Gesture control | MediaPipe Hands | Same as drone gesture system |
| Object detection | YOLOv8 (Jetson-optimized) | Real-time tagging in HUD |
| Thermal processing | FLIR Atlas SDK | FLIR's own SDK, Jetson compatible |
| Network | Tailscale | Always-on connection to EVO-X2 |
| Orchestration | Legacy via DiMOS MCP | Natural language commands |

---

## Build Phases

### Phase 1 — Basic HUD (Prove the Stack)
- [ ] Xreal Air 2 Pro + Jetson Orin NX
- [ ] RGB camera feed into AR overlay
- [ ] Legacy integration — whisper alerts, voice commands
- [ ] Drone feed routed to HUD (picture-in-picture)
- [ ] Basic facial recognition overlay

### Phase 2 — Thermal Vision
- [ ] Add FLIR Boson+ to Jetson
- [ ] Thermal feed fused with RGB in AR display
- [ ] Test through-wall heat detection (drywall first)
- [ ] Night vision via STARVIS camera

### Phase 3 — Through-Wall Radar
- [ ] Add Acconeer XM125 UWB radar
- [ ] Radar + thermal fusion → wall penetration overlay
- [ ] Tag detected targets in HUD with distance + movement

### Phase 4 — Custom Helmet
- [ ] Design 3D printable helmet housing (Fusion 360)
- [ ] Integrate all sensors into helmet shell
- [ ] Jetson NX mounted on vest/back with cable run
- [ ] Optional ballistic insert compatibility

### Phase 5 — Full OMOS Integration
- [ ] All 8 properties sensor feeds accessible from HUD
- [ ] Drone swarm commanded via hand gestures through HUD
- [ ] Legacy routes relevant intel to HUD automatically
- [ ] EagleEye event logs fed back to Home Assistant

---

## Legal / Regulatory Notes
- UWB radar: FCC Part 15 compliant devices available for civilian use
- Thermal cameras: No restriction on ownership in US (export restrictions apply for shipping abroad)
- AR glasses: No restrictions
- Drone integration: FAA Part 107 applies for any commercial drone ops
- "See through wall" radar: Legal for personal property use in US

---

## Related Codex Files
- `SO/02-infrastructure/OMOS-Global-Infrastructure.md`
- `SO/02-infrastructure/OMOS-Drone-Swarm.md`
- `LE/50-homestead/Greenhouse-Automation.md` *(to create)*

## External References
- Anduril EagleEye: https://www.anduril.com/eagleeye
- Xreal Air 2 Pro: https://www.xreal.com/air2pro
- Meta Ray-Ban Display: https://www.meta.com/glasses
- FLIR Boson+: https://www.flir.com/products/boson/
- Acconeer XM125: https://www.acconeer.com/products/xm125
- Jetson Orin NX: https://developer.nvidia.com/embedded/jetson-orin-nx
- OpenXR: https://www.khronos.org/openxr
- InsightFace: https://github.com/deepinsight/insightface
- DiMOS: https://github.com/dimensionalOS/dimos
