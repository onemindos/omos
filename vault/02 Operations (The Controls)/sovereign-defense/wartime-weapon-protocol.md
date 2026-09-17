---
domain: 04
name: Wartime Weapon Protocol — Tactical Drone & Robot Operating Procedures
status: reference
owner: zeus
tags: [robotics, drones, tactical, combat, protocol, security, meshtastic, eagleeye, px4, critical]
created: 2026-05-13
updated: 2026-05-13
classification: SOVEREIGN EYES ONLY
---

# Wartime Weapon Protocol

**Full tactical operating procedures for drone and robot deployment across all threat levels — from peacetime surveillance to active combat.**

This document defines exactly how Zeus's existing tools (OpenClaw, DimOS, EagleEye HUD, Meshtastic mesh, Tailscale, PX4/ArduPilot) combine into a unified tactical system with mode-switching based on threat level.

---

## 1. Threat Level Classification — OMOS DEFCON

Five operating modes. The system can escalate or de-escalate at any time. Escalation is automatic on trigger conditions. De-escalation requires Zeus's explicit command.

| Level | Name | Color | Trigger | Core Principle |
| --- | --- | --- | --- | --- |
| **LEVEL 5** | SOVEREIGN | 🟢 Green | Default state | Full intelligence, full connectivity |
| **LEVEL 4** | VIGILANT | 🔵 Blue | Anomaly detected | Enhanced monitoring, mesh standby |
| **LEVEL 3** | GUARDED | 🟡 Yellow | Confirmed threat / comms degraded | Mesh active, drones pre-staged |
| **LEVEL 2** | TACTICAL | 🟠 Orange | Active threat / comms denied | Physical controls, edge-only, no internet |
| **LEVEL 1** | WARFIGHTER | 🔴 Red | Under fire / life threat | One-button deploy, full autonomy, voice-only |

### Auto-Escalation Triggers

```
LEVEL 5 → LEVEL 4:
  - Perimeter sensor triggered (motion, thermal, acoustic)
  - Unknown person detected by facial recognition
  - Home Assistant anomaly alert
  - Unusual network activity / port scan detected

LEVEL 4 → LEVEL 3:
  - Multiple perimeter triggers in <60 seconds
  - Internet connectivity drops below threshold
  - GPS signal degradation detected
  - Unknown vehicle on property for >5 minutes
  - EagleEye HUD detects armed individual

LEVEL 3 → LEVEL 2:
  - Internet fully offline
  - Cell service unavailable
  - GPS denied / spoofed (position jumping)
  - Tailscale tunnel drops
  - Physical security breach (fence, door, window)

LEVEL 2 → LEVEL 1:
  - Gunshot detected (acoustic sensor)
  - Zeus triggers panic button (physical)
  - Dead-man switch timeout (Zeus unresponsive >30 seconds)
  - EagleEye HUD detects weapon pointed at Zeus
  - Manual escalation: voice command "WARFIGHTER" or physical switch
```

### De-Escalation

De-escalation is **never automatic**. Zeus must explicitly command:
- Voice: "STAND DOWN" + level name
- Physical: de-escalation switch on controller
- Telegram: Only works at LEVEL 4-5 (internet required)

This prevents the system from standing down while Zeus is incapacitated.

---

## 2. Operating Modes — Full Breakdown

### LEVEL 5: SOVEREIGN (Peacetime) 🟢

**Scenario:** Normal daily operations. Property patrols, farm monitoring, aerial photography.

#### Command Chain

```
Zeus (Telegram / Voice / Mac)
    │ message (~200ms)
    ▼
OpenClaw (omos-core VPS, Tailscale 100.98.25.10)
    │ EagleEye agent processes intent
    │ MCP tool call over Tailscale (~20-50ms)
    ▼
DimOS (Field Machine — Mac/Jetson, on Tailscale)
    │ localhost MCP server (0ms)
    │ MAVLink UDP 14550 (~5ms)
    ▼
Drone Hardware (DJI Mavic via RosettaDrone / Custom PX4)
```

#### Active Systems

| System | Status | Role |
| --- | --- | --- |
| OpenClaw (VPS) | ACTIVE | Mission planning, fleet coordination, Telegram interface |
| Legacy agent | ACTIVE | Orchestrator — routes Zeus's commands |
| EagleEye agent | ACTIVE | Drone specialist — translates commands to MCP calls |
| DimOS | ACTIVE | Real-time flight control, obstacle avoidance |
| Tailscale | ACTIVE | VPN mesh (VPS ↔ Mac ↔ Field machine) |
| Meshtastic | STANDBY | Mesh radios powered on, monitoring, not primary comms |
| EagleEye HUD | OPTIONAL | Worn during patrols for enhanced awareness |
| Home Assistant | ACTIVE | Property sensors feeding anomaly detection |
| ATAK | OFF | Not needed in peacetime |
| PX4 autonomous | STANDBY | Pre-loaded missions ready but not executing |

#### What Zeus Does

- Sends Telegram messages: "Patrol north field", "Follow that truck", "What's drone 3 seeing?"
- Uses DimOS CLI for direct commands: `dimos agent-send "orbit the barn"`
- Wears EagleEye HUD during property walks for enhanced awareness
- Monitors fleet via OpenClaw dashboard

#### Latency Budget

| Hop | Time | Acceptable? |
| --- | --- | --- |
| Zeus → Telegram → VPS | 200-2000ms | YES — no urgency |
| VPS → Tailscale → DimOS | 20-50ms | YES |
| DimOS → MAVLink → Drone | 5ms | YES |
| **Total** | **~300ms-2s** | **YES for peacetime** |

---

### LEVEL 4: VIGILANT (Anomaly Detected) 🔵

**Scenario:** Perimeter sensor tripped, unknown person detected, Home Assistant anomaly alert.

#### What Changes from LEVEL 5

```diff
+ Meshtastic mesh ACTIVE (primary backup comms confirmed working)
+ All drones powered on and pre-flighted
+ EagleEye HUD ON (if not already)
+ Home Assistant cameras streaming to HUD
+ Legacy auto-dispatches nearest drone to anomaly location
+ All drone batteries confirmed >50%
+ Acoustic sensors armed (gunshot detection active)
```

#### Automatic Actions (No Zeus Input Needed)

1. Legacy receives anomaly alert from Home Assistant
2. Legacy dispatches nearest drone to investigate:
   ```python
   await mcp_call("drone1:9990", "fly_to", {
       lat: anomaly.lat, lon: anomaly.lon, alt: 30
   })
   ```
3. Drone streams video to EagleEye HUD (picture-in-picture)
4. Legacy runs facial recognition on any detected persons
5. Legacy reports to Zeus via Telegram AND mesh radio simultaneously:
   ```
   "⚠️ LEVEL 4 VIGILANT — Unknown male at north gate.
   Drone 1 dispatched, ETA 45s. No match in face database.
   Mesh backup: ACTIVE. All drones: READY."
   ```
6. Meshtastic confirms connectivity to all mesh nodes

#### Zeus's Options

- **Watch and wait** — drone provides ISR, Legacy provides analysis
- **Challenge** — Legacy plays audio challenge through drone speaker
- **Escalate** — voice command "GUARDED" to move to LEVEL 3
- **De-escalate** — "STAND DOWN VIGILANT" if it's a known visitor

---

### LEVEL 3: GUARDED (Confirmed Threat / Comms Degraded) 🟡

**Scenario:** Multiple perimeter breaches, internet flickering, confirmed intruder, GPS degradation.

#### Command Chain Switches to Dual-Path

```
                    PRIMARY                          BACKUP
Zeus (Voice / Physical Controller)     Zeus (Voice / Physical Controller)
    │                                      │
    ▼                                      ▼
OpenClaw (VPS, if internet up)         Meshtastic Mesh (LoRa, always up)
    │ Tailscale MCP call                   │ Encrypted LoRa message
    ▼                                      ▼
DimOS (Field Machine)                  PX4 Companion (Jetson on drone)
    │ MAVLink                              │ Direct MAVLink serial
    ▼                                      ▼
Drone Hardware                         Drone Hardware
```

**Both paths are active simultaneously.** If internet drops, mesh takes over transparently.

#### What Changes from LEVEL 4

```diff
+ ATAK tablet ACTIVE — tactical map showing all assets
+ Drones airborne or on hot standby (rotors ready, 10-second launch)
+ Meshtastic promoted to PRIMARY backup (dual-path active)
+ PX4 autonomous behaviors armed:
    - Auto-RTH on signal loss
    - Auto-RTH on battery <20%
    - Geofence violation = auto-land
    - Obstacle avoidance = maximum sensitivity
+ Physical RC controller powered on and within reach
+ EagleEye HUD showing full tactical overlay:
    - All drone positions
    - All sensor contacts (thermal, UWB, camera)
    - Blue force tracking (Zeus's position on map)
    - Threat indicators pinned on map
+ Voice command engine ACTIVE (local ASR, no internet needed):
    - "ORBIT HERE" → nearest drone orbits Zeus's GPS
    - "EYES ON [direction]" → drone points camera [direction]
    - "OVERWATCH" → drone climbs to max altitude for wide view
    - "FOLLOW ME" → drone follows Zeus's GPS at 30m altitude
```

#### Meshtastic Mesh Network (Active Configuration)

```
┌──────────────────────────────────────────────────┐
│  MESH NODE LAYOUT                                 │
│                                                   │
│  Zeus (body-worn)                                 │
│    Heltec V3 ($25) + GPS module                   │
│    Channels:                                      │
│      ch0: LONGFAST (default) — telemetry/status  │
│      ch1: DRONE-CMD (encrypted) — drone commands │
│      ch2: ALERT (encrypted) — emergency triggers │
│                                                   │
│  Drone 1 (on-board)                              │
│    Heltec V3 ($25) + serial → PX4 companion      │
│    Receives ch1 commands → MAVLink translation   │
│                                                   │
│  Drone 2 (on-board)                              │
│    Same as Drone 1                               │
│                                                   │
│  Base Station (property)                         │
│    Heltec V3 ($25) + solar panel                 │
│    MQTT bridge to Home Assistant (when online)   │
│    Relays mesh → internet when internet available │
│                                                   │
│  Vehicle Node (truck/UTV)                        │
│    Heltec V3 ($25) + 12V power                   │
│    Mobile relay extends mesh range               │
│                                                   │
│  Range: 5-10km open terrain per hop              │
│  Encryption: AES-256 on ch1 + ch2               │
│  Latency: 200-500ms per hop (acceptable)         │
│  Cost: $25-40 per node × 5 = ~$150 total         │
└──────────────────────────────────────────────────┘
```

#### Mesh Command Protocol

Commands sent via Meshtastic channel `DRONE-CMD` (ch1, AES-256 encrypted):

| Mesh Message | Drone Action | Latency |
| --- | --- | --- |
| `D1:FLY:37.77,-122.41,50` | Drone 1 flies to GPS coordinate at 50m alt | <1s |
| `D1:ORBIT:37.77,-122.41,30,100` | Drone 1 orbits point at 30m alt, 100m radius | <1s |
| `D1:FOLLOW:ZEUS` | Drone 1 follows Zeus's GPS from mesh telemetry | <1s |
| `D1:RTH` | Drone 1 returns to home point | <1s |
| `D1:LAND` | Drone 1 lands immediately | <500ms |
| `D*:RTH` | ALL drones return to home | <1s |
| `D*:LAND` | ALL drones emergency land | <500ms |
| `D1:OBSERVE` | Drone 1 captures and relays camera snapshot | 2-5s |
| `PANIC` | All drones execute WARFIGHTER protocol | <500ms |

The PX4 companion computer (Jetson on each drone) translates these short mesh messages into MAVLink commands locally. No LLM involved. No internet needed.

---

### LEVEL 2: TACTICAL (Active Threat / Comms Denied) 🟠

**Scenario:** Internet fully down, cell service jammed, active intruder confirmed, GPS degraded.

#### Command Chain — Mesh Only, No Internet

```
Zeus (Physical Controller + Voice + EagleEye HUD)
    │
    ├─ Voice: Local ASR on Jetson (Whisper.cpp, no internet)
    │    "OVERWATCH" → Jetson translates → MAVLink command
    │
    ├─ Physical: Herelink RC / Custom controller with buttons
    │    Button 1: DEPLOY (launch + orbit Zeus's GPS)
    │    Button 2: RTH (all drones return home)
    │    Button 3: FOLLOW ME
    │    Button 4: EMERGENCY LAND
    │    Joystick: Manual override
    │
    ├─ Mesh: Meshtastic LoRa (encrypted ch1)
    │    Body-worn node → all drone nodes
    │    Zeus's GPS relayed every 5s for FOLLOW ME
    │
    └─ HUD: EagleEye AR overlay
         All drone positions on tactical map
         Thermal + night vision active
         UWB radar active (through-wall detection)
         Drone camera feeds in PiP
```

#### What Changes from LEVEL 3

```diff
- OpenClaw DISCONNECTED (no internet, VPS unreachable)
- Telegram DISCONNECTED
- Tailscale DISCONNECTED
- DimOS cloud LLM DISCONNECTED
+ ALL decision-making is LOCAL (edge compute only)
+ PX4 autonomous missions ARMED and HOT:
    - Pre-loaded patrol routes execute on one-button press
    - Obstacle avoidance runs on drone-local compute
    - Visual odometry active (GPS backup / GPS-denied nav)
+ Jetson Orin NX running full edge stack:
    - Whisper.cpp (voice → text, local)
    - YOLO v8 (object detection, local)
    - Qwen3-4B or Phi-4 (edge LLM for complex decisions, local)
    - InsightFace (facial recognition, local)
    - All inference <500ms per frame
+ ATAK is PRIMARY tactical display
    - All drone positions plotted
    - Zeus's position tracked (GPS + dead reckoning)
    - Threat positions marked
    - Shared with any allied ATAK users on mesh
+ Physical controller is PRIMARY input (not voice, not screen)
+ Dead-man switch ARMED:
    - If Zeus's mesh node goes silent >30 seconds → auto-escalate to LEVEL 1
    - Drone executes autonomous protective behavior
```

#### Edge AI Stack (What Runs Without Internet)

Everything on the Jetson Orin Nano (8GB, ~$250) or Orin NX (16GB, ~$400):

| Model | Task | Latency on Jetson | Size |
| --- | --- | --- | --- |
| **Whisper.cpp (tiny)** | Voice → text commands | ~200ms | 75MB |
| **YOLOv8n** | Object detection (people, vehicles, weapons) | ~15ms | 6MB |
| **Qwen3-4B (GGUF Q4)** | Complex decisions, scene understanding | ~500ms | 2.5GB |
| **InsightFace (buffalo_l)** | Facial recognition | ~50ms | 300MB |
| **Depth Anything v2 (small)** | Depth estimation for obstacle avoidance | ~30ms | 100MB |
| **Total VRAM** | | | **~3GB** |

All fit comfortably on an 8GB Jetson with room to spare.

#### Voice Command Dictionary (Local ASR, No Internet)

Whisper.cpp runs locally. A keyword-matching layer maps recognized phrases to MAVLink actions — no LLM in this path, just pattern matching for speed:

| Voice Command | Action | Latency |
| --- | --- | --- |
| "DEPLOY" | Launch nearest grounded drone, orbit Zeus at 30m | <500ms |
| "DEPLOY ALL" | Launch all available drones in formation | <500ms |
| "RTH" or "COME HOME" | All drones return to home coordinates | <500ms |
| "OVERWATCH" | Nearest drone climbs to max alt, wide-angle camera | <500ms |
| "FOLLOW ME" | Nearest drone follows Zeus's GPS at 30m offset | <500ms |
| "ORBIT HERE" | Nearest drone orbits Zeus's current position | <500ms |
| "EYES NORTH/SOUTH/EAST/WEST" | Drone points camera in cardinal direction | <300ms |
| "LIGHTS ON" | Drone activates spotlights/strobes | <200ms |
| "LIGHTS OFF" | Drone deactivates lights | <200ms |
| "LAND" | Nearest drone controlled landing | <500ms |
| "LAND ALL" | All drones emergency land | <500ms |
| "WARFIGHTER" | Escalate to LEVEL 1 | IMMEDIATE |
| "STAND DOWN [LEVEL]" | De-escalate | <500ms |
| "STATUS" | Audio report: battery, position, threats detected | ~2s |

Voice commands are processed locally on Zeus's body-worn Jetson or on the EagleEye HUD compute unit. Whisper.cpp + keyword matching = sub-500ms from speech to drone action.

#### EagleEye HUD at LEVEL 2

```
┌────────────────────────────────────────────────────────┐
│  TACTICAL HUD OVERLAY — LEVEL 2 ACTIVE                  │
│                                                          │
│  ┌─────────────────────────────────────┐                │
│  │  MAIN VIEW (what Zeus sees)          │  ┌──────────┐ │
│  │                                       │  │ DRONE 1  │ │
│  │  [Thermal overlay: 2 contacts]       │  │ camera   │ │
│  │  [UWB: 1 behind east wall, 3.2m]    │  │ feed     │ │
│  │  [Face: UNKNOWN — red tag]           │  │ PiP      │ │
│  │                                       │  └──────────┘ │
│  └─────────────────────────────────────┘                │
│                                                          │
│  ┌─────────────────────────────────────────────────────┐│
│  │  STATUS BAR                                          ││
│  │  D1: ■ AIRBORNE 47m | 73% | ORBIT                  ││
│  │  D2: ■ READY    gnd | 91% | STANDBY                ││
│  │  D3: ■ AIRBORNE 30m | 55% | FOLLOW-ME              ││
│  │  MESH: ■ 5/5 NODES | COMMS: ■ MESH-ONLY            ││
│  │  THREAT: 🟠 2 CONTACTS | GPS: ■ DEGRADED            ││
│  └─────────────────────────────────────────────────────┘│
│                                                          │
│  ┌───────────────┐                                      │
│  │  MINI-MAP      │   N                                 │
│  │  ● Zeus        │   ↑                                 │
│  │  ▲ D1 (orbit)  │                                     │
│  │  ▲ D3 (follow) │                                     │
│  │  ✖ Contact 1   │                                     │
│  │  ✖ Contact 2   │                                     │
│  │  ◎ Home base   │                                     │
│  └───────────────┘                                      │
└────────────────────────────────────────────────────────┘
```

---

### LEVEL 1: WARFIGHTER (Under Fire / Life Threat) 🔴

**Scenario:** Zeus is taking fire, active shooter, explosion, life-threatening emergency. Cannot use hands, cannot look at screens, cannot type. May be injured. May be unconscious.

#### This is Not Remote Control. This is Autonomous.

At LEVEL 1, the drones and robots act **independently** to protect Zeus. No commands needed after activation. The system defaults to pre-programmed autonomous behaviors and only Zeus's explicit "STAND DOWN" stops them.

#### Activation Methods

| Method | How | When to Use |
| --- | --- | --- |
| Voice | Shout "WARFIGHTER" | Zeus can speak but not use hands |
| Panic button | Press physical button (chest/belt) | Zeus can press but not speak |
| Dead-man switch | Zeus's mesh node silent >30s | Zeus is incapacitated / unconscious |
| Remote trigger | Allied ATAK user sends PANIC on mesh | Overwatch ally triggers for Zeus |

#### What Happens Automatically on LEVEL 1 Activation

**T+0 seconds — IMMEDIATE:**
1. All grounded drones launch immediately (PX4 auto-takeoff sequence)
2. All airborne drones switch to PROTECTIVE FORMATION around Zeus's last known GPS
3. Drone spotlights/strobes activate (disorientation + identification)
4. Audio: loud siren/alarm from drone speakers (deterrence)
5. Acoustic sensors lock on to gunshot direction vectors

**T+2 seconds — FORMATION:**
```
                    D1 (overwatch, 50m altitude)
                         │
                    ● Zeus
                   / │ \
                  /  │  \
               D2   D3   D4 (triangle at 20m alt, 30m radius)
```
- D1: Overwatch at maximum altitude, wide-angle camera, recording everything
- D2/D3/D4: Triangle formation around Zeus at 20m altitude, 30m radius
- All cameras streaming to EagleEye HUD (if Zeus is conscious and wearing it)
- All video recording locally to SD cards (evidence/replay)

**T+5 seconds — THREAT ASSESSMENT:**
- Edge YOLO detects and tags all humans in frame
- Known faces (family, allies) tagged BLUE
- Unknown faces tagged RED
- Weapons detected tagged RED + ALERT
- Threat bearing + distance calculated from multiple camera angles
- Audio triangulation: gunshot source direction plotted on ATAK map

**T+10 seconds — AUTONOMOUS BEHAVIORS:**

| Behavior | What the Drone Does | Why |
| --- | --- | --- |
| **Guardian Orbit** | Orbits Zeus at 20-30m radius, cameras inward | 360° overwatch, deters approach |
| **Spotlight Track** | Spotlights track detected threats | Dazzle/disorient + illuminate for cameras |
| **Siren Blast** | Continuous loud siren from all drones | Psychological deterrence, draws attention away from Zeus |
| **Video Record** | All cameras recording to local SD | Evidence preservation, no internet needed |
| **Threat Track** | Drones autonomously track moving threats | Maintains visual on hostiles |
| **RTH on Battery** | Auto-land at home when battery <15% | Prevents crash into Zeus's area |
| **Relay Position** | Zeus's position broadcast on mesh ch2 | Allies/emergency services can locate Zeus |

**T+30 seconds — SUSTAINED OPS:**
- Drones rotate: when one drone hits 20% battery, it RTH and a fresh one launches (if available)
- Continuous 360° recording of the scene
- Mesh network broadcasts Zeus's GPS and threat contact info on ch2 (ALERT channel)
- Any allied ATAK user on the mesh sees Zeus's position and threat contacts in real time

#### What This Does NOT Do

- **Does not fire weapons** — drones are ISR (intelligence, surveillance, reconnaissance) only
- **Does not ram targets** — no kamikaze behavior
- **Does not make lethal decisions** — deterrence and documentation only
- **Does not leave Zeus** — formation stays around Zeus's position
- **Does not contact police** — no internet, mesh only. Zeus or allies must call if possible

#### Ground Robot Behavior at LEVEL 1 (Future — Unitree Go2 / B2)

When ground robots are online:

| Robot | Behavior | Purpose |
| --- | --- | --- |
| **Scout** | Advances toward threat bearing, cameras streaming | Assess threat without Zeus moving |
| **Blocker** | Positions between Zeus and threat | Physical barrier, distraction |
| **Medic carrier** | Moves to Zeus with med kit | Automated supply delivery |
| **Relay** | Positions for optimal mesh coverage | Ensures comms don't drop |

---

## 3. Hardware Loadout by Threat Level

### Always Carried (EDC — Every Day Carry)

| Item | Purpose | Cost |
| --- | --- | --- |
| Meshtastic node (Heltec V3) | Mesh comms, GPS telemetry | $25 |
| Panic button (BLE → Meshtastic) | One-press WARFIGHTER trigger | $15 |
| Phone (iPhone) | Telegram + ATAK (when internet available) | — |

### LEVEL 4+ (Vigilant and Above)

| Item | Purpose | Cost |
| --- | --- | --- |
| EagleEye HUD (Xreal Air 2 Pro + Jetson) | Full tactical overlay | ~$850 |
| Herelink RC controller | Physical drone control + buttons | ~$500 |
| Body-worn Jetson Orin Nano | Edge AI + voice commands + mesh bridge | ~$250 |
| Extra drone batteries (×4) | Sustained ops, hot-swap | ~$200 |

### LEVEL 2+ (Tactical and Above)

| Item | Purpose | Cost |
| --- | --- | --- |
| ATAK tablet (ruggedized Android) | Tactical map, blue force tracking | ~$300 |
| Meshtastic repeater (solar, property-mounted) | Extended mesh range | ~$40 |
| Throat mic + bone conduction earpiece | Voice commands under fire | ~$50 |
| Chest-mount panic button (waterproof) | WARFIGHTER trigger | ~$15 |

### Full Tactical Kit (All Levels)

| Category | Items | Total Cost |
| --- | --- | --- |
| **Drones** | 2-4× Custom PX4 quads (non-DJI) | ~$2,000-4,000 |
| **Compute** | Jetson Orin NX (body) + Jetson Orin Nano per drone | ~$1,000 |
| **Comms** | 5× Meshtastic nodes + 1 solar repeater | ~$175 |
| **Display** | EagleEye HUD + ATAK tablet | ~$1,150 |
| **Control** | Herelink RC + throat mic + panic button | ~$565 |
| **Power** | 8× drone batteries + body battery pack | ~$500 |
| **Total** | | **~$5,400-7,400** |

---

## 4. Software Architecture by Mode

### LEVEL 5-4: Full Stack (Internet Available)

```
┌──────────────────────────────────────────────────────────┐
│  CLOUD LAYER                                              │
│  OpenClaw (omos-core VPS)                                │
│    ├── Legacy (orchestrator)                              │
│    ├── EagleEye (drone specialist)                       │
│    ├── Sovereignty (SO infra agent)                      │
│    └── Sentinel (omos-sentinel watchdog)                 │
│  Telegram bot (@onemindos_bot)                           │
│  Home Assistant (property sensors)                       │
├──────────────────────────────────────────────────────────┤
│  NETWORK LAYER                                            │
│  Tailscale VPN (primary)                                 │
│  Meshtastic mesh (standby/backup)                        │
│  WiFi / LTE (transport)                                  │
├──────────────────────────────────────────────────────────┤
│  EDGE LAYER                                               │
│  DimOS (robotics runtime, MCP server per drone)          │
│  Jetson Orin (edge AI — face rec, object detection)      │
│  PX4 autopilot (pre-loaded missions, failsafe armed)    │
│  EagleEye HUD (sensor fusion overlay)                    │
│  Whisper.cpp (voice command standby)                     │
├──────────────────────────────────────────────────────────┤
│  HARDWARE LAYER                                           │
│  Custom PX4 drones + DJI Mavic (peacetime)               │
│  Meshtastic radio nodes (body + drone + base)            │
│  RC controller (standby)                                 │
│  Cameras, thermal, UWB sensors (HUD)                     │
└──────────────────────────────────────────────────────────┘
```

### LEVEL 3: Dual-Path (Internet Degraded)

```
┌──────────────────────────────────────────────────────────┐
│  CLOUD LAYER (INTERMITTENT — may drop)                    │
│  OpenClaw sends missions when internet available          │
│  Telegram alerts queued for when connectivity returns     │
├──────────────────────────────────────────────────────────┤
│  MESH LAYER (PRIMARY BACKUP)                              │
│  Meshtastic LoRa mesh — always on, AES-256               │
│  MQTT bridge to Home Assistant (when internet up)         │
│  Drone command channel (ch1) active                       │
│  Alert channel (ch2) active                               │
├──────────────────────────────────────────────────────────┤
│  EDGE LAYER (AUTONOMOUS)                                  │
│  DimOS local mode (no cloud LLM, edge model only)        │
│  PX4 autonomous missions armed                           │
│  ATAK tactical map (mesh-networked to allies)            │
│  Voice commands active (Whisper.cpp)                     │
│  EagleEye HUD full tactical overlay                      │
├──────────────────────────────────────────────────────────┤
│  HARDWARE LAYER (NO CHANGE)                               │
│  Same as LEVEL 5 but RC controller in-hand               │
└──────────────────────────────────────────────────────────┘
```

### LEVEL 2-1: Edge Only (No Internet)

```
┌──────────────────────────────────────────────────────────┐
│  CLOUD LAYER: OFFLINE                                     │
│  OpenClaw: unreachable                                    │
│  Telegram: unreachable                                    │
│  Tailscale: disconnected                                 │
├──────────────────────────────────────────────────────────┤
│  MESH LAYER (PRIMARY COMMS)                               │
│  Meshtastic LoRa — only comms available                  │
│  Zeus GPS → drone nodes (every 5s)                       │
│  Drone telemetry → Zeus node (every 5s)                  │
│  ALERT channel broadcasting if LEVEL 1                   │
├──────────────────────────────────────────────────────────┤
│  EDGE LAYER (FULLY AUTONOMOUS)                            │
│  Jetson Orin: ALL AI runs locally                        │
│    ├── Whisper.cpp: voice → command translation          │
│    ├── YOLOv8: object/weapon/person detection            │
│    ├── InsightFace: friend-or-foe identification         │
│    ├── Qwen3-4B: complex decision making (if needed)     │
│    └── Depth Anything: obstacle avoidance                │
│  PX4 autopilot: pre-programmed behaviors executing       │
│  EagleEye HUD: full overlay, all local                   │
│  ATAK: tactical map via mesh network                     │
├──────────────────────────────────────────────────────────┤
│  HARDWARE LAYER (INDEPENDENT)                             │
│  Custom PX4 drones (no DJI — no geofencing/kill switch)  │
│  Meshtastic radios (AES-256, no internet needed)         │
│  RC controller (physical buttons primary input)          │
│  All sensors local (thermal, UWB, NV, cameras)           │
└──────────────────────────────────────────────────────────┘
```

---

## 5. Communication Failover Cascade

When systems fail, comms degrade gracefully:

```
Internet available?
    YES → Tailscane VPN active
        → OpenClaw + Telegram (full intelligence, LEVEL 5-4)
        → DimOS MCP over Tailscane (cloud LLM available)
    NO ↓

Cell service available?
    YES → LTE data (degraded)
        → Telegram may work (slow)
        → DimOS local mode (edge LLM only)
    NO ↓

WiFi available?
    YES → Local network only
        → DimOS on LAN (no cloud, edge LLM)
        → Home Assistant accessible on LAN
    NO ↓

Meshtastic mesh available?
    YES → LoRa mesh (LEVEL 2-1 mode)
        → Short command strings (D1:FLY:lat,lon,alt)
        → GPS telemetry (Zeus position → drones)
        → Drone telemetry (battery, status → Zeus)
        → ATAK over mesh (if allies on mesh)
    NO ↓

Radio link to drone available?
    YES → Direct RC control (LEVEL 1, manual override)
        → Herelink controller → drone
        → Physical sticks, no software in loop
    NO ↓

Nothing available:
    → Dead-man switch activates (Zeus silent >30s)
    → PX4 executes pre-loaded WARFIGHTER autonomous mission:
        1. Orbit last known Zeus GPS
        2. Strobes + siren active
        3. Record everything to SD card
        4. RTH when battery <15%
    → This runs ENTIRELY on the flight controller with ZERO connectivity
```

---

## 6. PX4 Pre-Loaded Mission Library

These missions live on the PX4 flight controller's SD card. They execute without any computer, phone, LLM, or internet. The flight controller runs them directly.

| Mission ID | Name | Trigger | What It Does |
| --- | --- | --- | --- |
| `M01` | GUARDIAN ORBIT | Panic button / dead-man | Takeoff → orbit Zeus GPS at 30m alt, 30m radius, indefinitely |
| `M02` | PROPERTY PATROL | Scheduled / button | Fly pre-loaded waypoint fence around property |
| `M03` | PERIMETER SWEEP | Button or mesh command | Fly to each perimeter sensor location, hover 10s, photo, next |
| `M04` | FOLLOW BEACON | Button | Follow Meshtastic GPS beacon (Zeus) at 30m offset |
| `M05` | EMERGENCY OVERWATCH | WARFIGHTER auto | Max altitude, wide camera, 360° slow orbit, record |
| `M06` | RTH SAFE | Any RTH trigger | Return to home, descend, land, disarm |
| `M07` | AREA DENIAL SIREN | WARFIGHTER auto | Orbit threat bearing at 20m, siren + strobe, aggressive maneuvers |
| `M08` | SEARCH GRID | Mesh command | Fly expanding square search from given GPS center |

All missions have built-in failsafes:
- Battery <15% → abort mission, execute M06 (RTH)
- Signal loss >60s → execute M01 (guardian orbit) then M06 after 10 min
- Geofence violation → auto-land immediately
- Crash detected (IMU) → disarm motors

---

## 7. Integration Map — How Everything Connects

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        ZEUS (OPERATOR)                                   │
│   Voice (throat mic) ─┐                                                  │
│   Panic button ───────┤                                                  │
│   RC controller ──────┤                                                  │
│   EagleEye HUD ───────┤                                                  │
│   Phone (Telegram) ───┤                                                  │
│   ATAK tablet ────────┘                                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                     COMPUTE (BODY-WORN JETSON)                           │
│                                                                          │
│   ┌─────────────┐   ┌──────────────┐   ┌──────────────┐                │
│   │ Whisper.cpp  │   │ YOLOv8       │   │ InsightFace  │                │
│   │ Voice → Cmd  │   │ Detection    │   │ Face Rec     │                │
│   └──────┬──────┘   └──────┬───────┘   └──────┬───────┘                │
│          │                  │                   │                         │
│          └──────────┬───────┴───────────────────┘                        │
│                     │                                                     │
│          ┌──────────▼──────────┐                                         │
│          │ Command Router       │                                         │
│          │ Voice match → action │                                         │
│          │ Mesh rx → action     │                                         │
│          │ Button → action      │                                         │
│          └──────┬──────┬───────┘                                         │
│                 │      │                                                  │
│     ┌───────────┘      └────────────┐                                    │
│     ▼                               ▼                                    │
│  ┌─────────┐                  ┌───────────┐                              │
│  │ Mesh TX  │                  │ MAVLink   │                              │
│  │ (LoRa)   │                  │ (serial)  │                              │
│  └────┬────┘                  └─────┬─────┘                              │
├───────┼────────────────────────────┼────────────────────────────────────┤
│       │    MESH NETWORK            │    DIRECT LINK                      │
│       │                            │                                     │
│  Meshtastic LoRa                   │                                     │
│  915MHz, AES-256                   │                                     │
│  5-10km per hop                    │                                     │
│       │                            │                                     │
│  ┌────▼─────────────┐        ┌─────▼──────────────┐                     │
│  │ Drone 1 Mesh Node │        │ Drone 1 PX4 FC     │                     │
│  │ Heltec V3         │──────►│ Companion Jetson    │                     │
│  │ Receives commands │ serial │ Edge AI + autopilot │                     │
│  └──────────────────┘        └────────┬────────────┘                     │
│                                       │ MAVLink                          │
│                                       ▼                                  │
│                              ┌────────────────────┐                      │
│                              │ Drone Motors/ESCs   │                      │
│                              │ GPS + IMU + Baro    │                      │
│                              │ Camera + Thermal    │                      │
│                              │ Spotlight + Speaker │                      │
│                              └────────────────────┘                      │
├─────────────────────────────────────────────────────────────────────────┤
│                     CLOUD (WHEN AVAILABLE)                                │
│                                                                          │
│  OpenClaw (VPS) ←── Tailscane VPN ──→ Jetson                            │
│    Legacy agent ← orchestrator                                           │
│    EagleEye agent ← drone specialist                                     │
│    Telegram ← Zeus's peacetime interface                                 │
│                                                                          │
│  Internet drops → cloud layer goes dark                                  │
│  Everything below this line keeps running                                │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Training & Drill Protocol

The system is useless if Zeus hasn't trained with it under stress. Quarterly drills per level:

### Drill Schedule

| Drill | Frequency | What Zeus Practices |
| --- | --- | --- |
| **SOVEREIGN patrol** | Weekly | Normal drone patrol via Telegram, review footage |
| **VIGILANT response** | Monthly | Respond to simulated anomaly, verify auto-dispatch works |
| **GUARDED failover** | Monthly | Intentionally disable internet, verify mesh takeover |
| **TACTICAL comms-denied** | Quarterly | Full blackout — mesh only, physical controls only, voice only |
| **WARFIGHTER full auto** | Quarterly | Trigger panic button, verify autonomous formation deploys |

### TACTICAL Drill Procedure (Quarterly)

1. Turn off all internet (WiFi router off, phone in airplane mode)
2. Verify Meshtastic mesh is communicating between all nodes
3. Launch one drone using physical controller button
4. Issue voice commands: DEPLOY, ORBIT HERE, FOLLOW ME, OVERWATCH
5. Verify EagleEye HUD shows drone position and camera feed
6. Verify ATAK shows blue force tracking over mesh
7. Issue RTH command via voice and verify drone returns
8. Restore internet and verify OpenClaw reconnects, DimOS re-syncs

### WARFIGHTER Drill Procedure (Quarterly)

1. All internet off, all screens down
2. Press panic button
3. Verify all drones auto-launch within 5 seconds
4. Verify formation establishes around Zeus within 15 seconds
5. Verify strobes and siren activate
6. Verify all cameras recording to SD card
7. Verify mesh is broadcasting Zeus's position on ch2
8. Issue "STAND DOWN WARFIGHTER" → verify controlled de-escalation
9. Review all SD card footage for coverage gaps
10. After-action report: what worked, what didn't, what needs adjustment

### Muscle Memory Requirements

Zeus must be able to do these **without thinking**, under stress, in the dark:
- Find and press panic button by feel (chest or belt mount)
- Issue voice commands while running/moving
- Switch RC controller from standby to active by feel
- Identify which drone is which on EagleEye HUD mini-map
- Issue "STAND DOWN" clearly under stress

---

## 9. Rules of Engagement (Drone Autonomous Behavior Limits)

**These rules are hard-coded into PX4 flight controller and cannot be overridden by any software, LLM, or remote command:**

### ALWAYS

- Record video to local SD card when airborne
- Broadcast GPS position on mesh when airborne
- Return to home when battery <15%
- Land immediately if flight controller detects crash/damage
- Respect geofence boundaries (property line + 100m buffer)

### NEVER

- Fly above 400ft AGL (FAA Part 107 limit)
- Fly beyond visual line of sight without active waiver
- Ram, collide with, or physically contact any person
- Drop objects on any person
- Carry weapons or weapon-like payloads
- Override geofence for any reason
- Fly into controlled airspace without authorization
- Operate without Zeus's knowledge (no fully rogue behavior)

### CONDITIONAL (WARFIGHTER only)

- MAY descend below 20m to illuminate threat with spotlight
- MAY activate siren at maximum volume for deterrence
- MAY follow threat at close range (10m minimum distance) for tracking
- MAY fly aggressive patterns (rapid orbits, quick altitude changes) for deterrence
- MAY NOT make physical contact with any person under any circumstances

---

## 10. Post-Incident Protocol

After any LEVEL 2+ event:

1. **SECURE:** Ensure all drones accounted for and landed safely
2. **PRESERVE:** Pull SD cards from all drones — do not overwrite footage
3. **DOCUMENT:** Legacy generates incident report (when internet restored):
   - Timeline of events
   - All sensor data logged
   - Drone flight paths
   - Detected threats and their movements
   - Audio recordings (if throat mic was active)
4. **REPORT:** If law enforcement involvement needed, provide SD card footage
5. **DEBRIEF:** After-action review — what worked, what failed, hardware damage assessment
6. **REARM:** Charge all batteries, replace any damaged components, reload PX4 missions, verify all mesh nodes operational
7. **LOG:** Save incident to vault: `00-24 SO/04 Research/Robotics/incidents/YYYY-MM-DD-incident.md`

---

## 11. Build Roadmap

### Phase 1: Foundation (Now — 3 months)

| Task | Status | Notes |
| --- | --- | --- |
| DimOS installed and tested with DJI Mavic | 🔴 Not started | Peacetime capability first |
| Meshtastic nodes purchased (5× Heltec V3) | 🔴 Not started | ~$125 total |
| Meshtastic mesh tested on property | 🔴 Not started | Verify range, relay, encryption |
| OpenClaw → DimOS MCP integration | 🔴 Not started | EagleEye agent sends commands |
| EagleEye HUD Phase 1 (Xreal + Jetson basic) | 🔴 Not started | See OMOS-EagleEye.md |
| ATAK installed on tablet | 🔴 Not started | Free download for Android |

### Phase 2: Tactical Capability (3-6 months)

| Task | Status | Notes |
| --- | --- | --- |
| Custom PX4 drone built (first unit) | 🔴 Not started | Non-DJI, full autonomous capable |
| PX4 mission library loaded (M01-M08) | 🔴 Not started | Pre-programmed behaviors |
| Whisper.cpp voice commands on Jetson | 🔴 Not started | Local ASR, keyword matching |
| Meshtastic → PX4 command bridge | 🔴 Not started | Mesh commands → MAVLink translation |
| Herelink RC controller configured | 🔴 Not started | Physical buttons mapped to actions |
| First TACTICAL drill (internet blackout) | 🔴 Not started | Quarterly schedule begins |

### Phase 3: Warfighter Ready (6-12 months)

| Task | Status | Notes |
| --- | --- | --- |
| Multi-drone PX4 fleet (2-4 units) | 🔴 Not started | Formation flying tested |
| Dead-man switch implemented | 🔴 Not started | Mesh heartbeat → auto-WARFIGHTER |
| Autonomous formation behaviors tested | 🔴 Not started | Guardian orbit, overwatch, follow |
| Siren + strobe modules integrated | 🔴 Not started | Deterrence hardware on drones |
| Full WARFIGHTER drill completed | 🔴 Not started | Panic button → auto formation |
| EagleEye HUD Phase 2-3 (thermal + radar) | 🔴 Not started | See OMOS-EagleEye.md |
| Ground robot integration (Unitree Go2) | 🔴 Not started | Scout/blocker behaviors |

### Phase 4: Full Integration (12+ months)

| Task | Status | Notes |
| --- | --- | --- |
| OpenClaw ↔ PX4 mode switching | 🔴 Not started | Peacetime (OpenClaw) ↔ Tactical (PX4 autonomous) |
| ATAK mesh integration with all assets | 🔴 Not started | All drones + robots + Zeus on ATAK |
| Blockchain governance (flight rules on-chain) | 🔴 Not started | Immutable geofence + altitude rules |
| Multi-property coverage | 🔴 Not started | Mesh repeaters at all 8 properties |
| Community mesh network | 🔴 Not started | Allied neighbors on shared mesh |

---

## Related Documents

- [[drone-stack-architecture-openclaw-dimos]] — Peacetime architecture (OpenClaw + DimOS)
- [[drone-combat-readiness-assessment]] — Why OpenClaw alone fails for combat
- [[openmind-om1-deep-analysis]] — OM1 patterns (dual LLM, modes, governance)
- [[om1-vs-dimos-vs-openclaw-comparison]] — Platform comparison
- [[OMOS-EagleEye]] — Wearable sensor fusion HUD
- [[meshtastic-mesh-integration]] — Mesh network integration
- [[meshtastic-starter-guide]] — Meshtastic hardware setup
- [[multi-drone-fleet-architecture]] — Fleet architecture design
- [[OMOS-Drone-Swarm]] — Swarm concepts
