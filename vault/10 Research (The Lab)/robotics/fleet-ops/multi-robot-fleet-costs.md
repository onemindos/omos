# Multi-Robot Fleet: Costs & Compatible Robots

> **Created:** 2026-01-22
> **Key Question:** Is OpenMind Edge Compute $997 per robot or per location?

---

## Edge Compute Model: PER LOCATION, NOT PER ROBOT ✅

**One $997 OpenMind Edge Compute can control MULTIPLE robots at the same location.**

### How It Works

```
┌──────────────────────────────────────────────────────────────┐
│                  SINGLE EDGE DEVICE ($997)                    │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐  │
│  │        OpenMind OM1 + FABRIC Network                   │  │
│  │                                                         │  │
│  │    Acts as "Robot Hub" for all robots at location     │  │
│  └────────────────────────────────────────────────────────┘  │
│                              │                                │
│         ┌────────────────────┼────────────────────┐          │
│         │                    │                    │          │
│         ▼                    ▼                    ▼          │
│   ┌──────────┐         ┌──────────┐         ┌──────────┐    │
│   │ Unitree  │         │ Unitree  │         │   DJI    │    │
│   │  Go1 #1  │         │  Go1 #2  │         │  Drone   │    │
│   └──────────┘         └──────────┘         └──────────┘    │
│                                                               │
│  All robots share:                                           │
│  - Same edge compute                                         │
│  - Same FABRIC network node                                  │
│  - Coordinated via OM1                                       │
│  - Video streams via LiveKit                                 │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

**Think of it like:**
- 1 Router → Multiple WiFi devices
- 1 Hub → Multiple Bluetooth devices
- **1 OpenMind Edge → Multiple robots/drones**

---

## Cost Model: Multi-Robot Scenarios

### Scenario 1: Single Location (Barn/Farm)

| Item | Quantity | Unit Cost | Total |
|------|----------|-----------|-------|
| **OpenMind Edge Compute** | 1 | $997 | $997 |
| Unitree Go1 Robot #1 | 1 | Owned | $0 |
| Unitree Go1 Robot #2 | 1 | $3,500 | $3,500 |
| DJI Mini 4 Pro Drone | 1 | $759 | $759 |
| **TOTAL** | | | **$5,256** |

**Controls:** 2 quadrupeds + 1 drone from single edge device

### Scenario 2: Multi-Location Fleet

| Location | Edge Device | Robots | Cost |
|----------|-------------|--------|------|
| **Barn** | 1x $997 | 2x Unitree + 1x Drone | $5,256 |
| **Greenhouse** | 1x $997 | 1x Unitree + 1x Drone | $4,756 |
| **Home** | 1x $997 | 1x Humanoid | $1,997 |
| **TOTAL** | **3x $997** | **4 robots + 2 drones** | **$11,009** |

**Key Point:** You need 1 edge device PER LOCATION, not per robot.

### Monthly Service Costs (Shared Across Fleet)

| Service | Cost | Notes |
|---------|------|-------|
| AWS Bedrock (Claude) | $50 | One brain for all robots |
| Deepgram STT | $20 | Voice transcription |
| Cartesia TTS | $10 | Voice synthesis |
| OpenMind FABRIC | $50 | Entire fleet coordination |
| Tailscale Pro | $5 | Network for all locations |
| **TOTAL** | **$135/mo** | **Does NOT scale with robot count** |

---

## Compatible Robots & Drones

### 1. Quadruped Robots (OpenMind OM1 Native)

| Robot | Price | Speed | Payload | Battery | Best For |
|-------|-------|-------|---------|---------|----------|
| **Unitree Go1 EDU** | $3,500 | 3.5 m/s | 3 kg | 2-4 hrs | Patrol, inspection, learning |
| **Unitree Go2** | $1,600-$5,000 | 3 m/s | 5 kg | 1-2 hrs | Lightweight tasks, indoor |
| **Unitree B1** | $50,000+ | 5 m/s | 50 kg | 2-3 hrs | Heavy payload, industrial |
| **Boston Dynamics Spot** | $75,000 | 1.6 m/s | 14 kg | 90 min | Enterprise, construction |
| **Anymal (ANYbotics)** | $100,000+ | 1 m/s | 10 kg | 2-4 hrs | Industrial inspection |
| **Deeprobotics Lite3** | $15,000 | 2 m/s | 8 kg | 2 hrs | Mid-range, China market |

**Recommendation:** Start with **Unitree Go1 EDU** ($3,500) - best value, full SDK access, educational discount available.

### 2. Humanoid Robots (OpenMind OM1 Future)

| Robot | Price | Height | Weight | Status |
|-------|-------|--------|--------|--------|
| **Unitree H1** | TBD (est. $90K) | 180cm | 47kg | Pre-order 2024 |
| **Figure 01** | Not for sale | 168cm | 60kg | Enterprise only |
| **Tesla Optimus** | TBD | 173cm | 73kg | Future |
| **1X NEO** | TBD | 165cm | 30kg | Beta testing |
| **Agility Digit** | $250,000 | 175cm | 65kg | Enterprise rental |

**Recommendation:** Wait for **Unitree H1** release - same ecosystem as your Go1, likely best integration.

### 3. Drones (NATS/FABRIC Integration)

#### Consumer Drones (Ready Now)

| Drone | Price | Flight Time | Camera | Range | Best For |
|-------|-------|-------------|--------|-------|----------|
| **DJI Mini 4 Pro** | $759 | 34 min | 4K60 | 25 km | Aerial patrol, compact |
| **DJI Mavic 3** | $2,049 | 43 min | 4K120 | 30 km | Professional inspection |
| **DJI Air 3** | $1,099 | 46 min | Dual 4K | 20 km | Long flight, dual camera |
| **Autel EVO II Pro** | $1,795 | 40 min | 6K | 9 km | US-made alternative |
| **Skydio 2+** | $1,099 | 27 min | 4K60 | 6 km | Autonomous tracking |

**Recommendation:** **DJI Mini 4 Pro** - Under 250g (no FAA registration), excellent camera, great range.

#### Agricultural Drones (Future)

| Drone | Price | Purpose | Payload |
|-------|-------|---------|---------|
| **DJI Agras T40** | $16,000 | Crop spraying | 40 kg |
| **DJI Mavic 3 Multispectral** | $5,500 | Plant health analysis | Sensors |
| **Sentera 6X** | $25,000 | Precision agriculture | Multispectral |

### 4. Specialty Robots (NATS Direct Integration)

| Robot | Price | Purpose | Integration |
|-------|-------|---------|-------------|
| **FarmBot Genesis** | $4,500 | Automated gardening | HTTP API → NATS |
| **iRobot Roomba i7+** | $599 | Indoor cleaning | WiFi → MQTT → NATS |
| **Roborock S8 Pro** | $1,599 | Advanced cleaning | LAN API → NATS |
| **Husqvarna Automower** | $3,000 | Lawn mowing | Connect API → NATS |
| **Pool Patrol Dolphin** | $1,299 | Pool cleaning | Bluetooth → NATS |

### 5. Robotic Arms (Workshop Integration)

| Arm | Price | Reach | Payload | Best For |
|-----|-------|-------|---------|----------|
| **Dobot MG400** | $5,500 | 440mm | 750g | Desktop assembly |
| **Elephant Robotics myCobot** | $799 | 280mm | 250g | Education, testing |
| **Universal Robots UR5e** | $35,000 | 850mm | 5kg | Industrial, precise |
| **Franka Emika Panda** | $25,000 | 855mm | 3kg | Research, sensitive |

---

## LiveKit Integration: Video from ALL Robots

### How LiveKit Works with Multi-Robot Fleet

```
┌─────────────────────────────────────────────────────────────┐
│                    LIVEKIT VIDEO HUB                         │
│                                                              │
│  ┌────────────┐    ┌────────────┐    ┌────────────┐        │
│  │  Room:     │    │  Room:     │    │  Room:     │        │
│  │  barn-1    │    │  barn-2    │    │  drone-sky │        │
│  └─────┬──────┘    └─────┬──────┘    └─────┬──────┘        │
│        │                 │                  │               │
│        ▼                 ▼                  ▼               │
│  ┌────────────────────────────────────────────────────┐    │
│  │         OneMind OS UI (Web/Mobile)                 │    │
│  │                                                     │    │
│  │  📹 [Unitree #1]  📹 [Unitree #2]  📹 [Drone]     │    │
│  │                                                     │    │
│  │  All streams in one dashboard, switchable          │    │
│  └────────────────────────────────────────────────────┘    │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### Robot Camera → LiveKit → UI

| Robot | Cameras | LiveKit Room | Stream Quality |
|-------|---------|--------------|----------------|
| Unitree Go1 | 5x cameras (360°) | `robot-go1-001` | 1080p30 |
| Unitree Go2 | 5x cameras | `robot-go2-001` | 1080p30 |
| DJI Mini 4 Pro | 1x gimbal 4K | `drone-dji-001` | 4K60 (downscaled 1080p) |
| Boston Dynamics Spot | 5x cameras | `robot-spot-001` | 1080p30 |

**Key Features:**
- **Multi-stream dashboard** - See all robots at once
- **Voice switching** - "Legacy, show me drone view"
- **Low latency** - <200ms WebRTC
- **Bandwidth efficient** - Adaptive bitrate
- **Recording** - All streams saved for analysis

### LiveKit Cost (Multi-Robot)

**Good News:** LiveKit self-hosted = FREE unlimited streams

| Option | Cost | Streams | Storage |
|--------|------|---------|---------|
| **Self-hosted (Current)** | $0 | Unlimited | Your AWS |
| LiveKit Cloud (alternative) | $99/mo | 500 hours | 100GB |

**Your Setup:** Already self-hosting LiveKit on `100.102.21.44:7880` ✅

---

## Recommended Starter Fleets

### Budget Fleet ($5K-6K)

**Purpose:** Farm patrol + aerial inspection
```
1x OpenMind Edge Compute        $997
1x Unitree Go1 (you have)       $0
1x Unitree Go2                  $1,600
1x DJI Mini 4 Pro              $759
─────────────────────────────────────
TOTAL                          $3,356
+ Monthly services             $135/mo
```

**Capabilities:**
- 2 ground robots patrol barn/greenhouse
- 1 drone for aerial view
- All coordinated via single edge device
- All video streams to dashboard

### Farm Fleet ($10K-12K)

**Purpose:** Full farm automation
```
2x OpenMind Edge Compute        $1,994  (barn + greenhouse)
2x Unitree Go1                  $3,500
1x DJI Mini 4 Pro              $759
1x DJI Mavic 3                 $2,049
1x FarmBot Genesis             $4,500
─────────────────────────────────────
TOTAL                          $12,802
+ Monthly services             $135/mo
```

**Capabilities:**
- 2 ground robots (barn + greenhouse patrol)
- 2 drones (patrol + crop inspection)
- Automated garden (FarmBot)
- Multi-location coordination
- Voice control for entire fleet

### Enterprise Fleet ($50K-75K)

**Purpose:** Large property/commercial farm
```
3x OpenMind Edge Compute        $2,991  (barn + greenhouse + home)
2x Unitree Go1                  $7,000
1x Unitree H1 Humanoid         $90,000 (when available)
2x Boston Dynamics Spot        $150,000
2x DJI Mavic 3 Multispectral   $11,000
1x DJI Agras T40 (spraying)    $16,000
─────────────────────────────────────
TOTAL                          $276,991
+ Monthly services             $135/mo
```

---

## Drone-Specific Actions (What You Can Command)

### With DJI Drones + OneMind OS

**Voice Commands:**
```
"Legacy, launch drone for barn patrol"
"Drone, follow me at 10 feet"
"Scan the greenhouse for plant issues"
"Take aerial photo of entire property"
"Return to base, battery low"
"Show me drone camera feed"
```

**Automated Workflows:**
```python
# Morning aerial survey
@workflow
def morning_aerial_survey():
    """Automated daily property inspection"""
    drone.takeoff()
    drone.navigate_waypoints([
        "barn_roof",
        "greenhouse_perimeter",
        "pasture_fence_line",
        "equipment_yard"
    ])

    # Stream video to Legacy for analysis
    video = drone.stream_to_livekit()
    analysis = legacy.analyze_video(video,
        objectives=["roof_damage", "fence_breaks", "animal_count"]
    )

    drone.return_to_base()
    legacy.send_report(analysis)
```

**DJI SDK Features (Available):**
```
- Waypoint navigation
- Obstacle avoidance
- Follow mode
- Point of interest orbit
- RTH (Return to Home)
- Live video stream
- Camera gimbal control
- Flight telemetry
```

### Integration Methods

| Drone Type | SDK | Protocol | Integration |
|------------|-----|----------|-------------|
| **DJI Consumer** | Mobile SDK | UDP | NATS bridge → OM1 |
| **DJI Enterprise** | Onboard SDK | UART/USB | Direct to edge device |
| **Skydio** | Cloud API | REST | HTTP → NATS |
| **Autel** | Mobile SDK | WiFi | NATS bridge |
| **Custom/DIY** | MAVLink | Serial | Direct NATS |

---

## Final Cost Summary: Multi-Robot Scenarios

### Scenario 1: You + 1 More Robot + 1 Drone
```
Hardware (one-time):
  1x OpenMind Edge           $997
  1x Unitree Go2             $1,600
  1x DJI Mini 4 Pro          $759
  ───────────────────────────────
  TOTAL HARDWARE             $3,356

Monthly:
  Services                   $135/mo

First Year Total:           $4,976
```

### Scenario 2: Multi-Location, 3 Robots, 2 Drones
```
Hardware (one-time):
  2x OpenMind Edge           $1,994
  2x Unitree Go1             $3,500
  2x DJI Mini 4 Pro          $1,518
  ───────────────────────────────
  TOTAL HARDWARE             $7,012

Monthly:
  Services                   $135/mo

First Year Total:           $8,632
```

### Scenario 3: Full Farm Fleet
```
Hardware (one-time):
  3x OpenMind Edge           $2,991
  3x Unitree Go1             $10,500
  2x DJI Drones              $2,848
  1x FarmBot                 $4,500
  ───────────────────────────────
  TOTAL HARDWARE             $20,839

Monthly:
  Services                   $135/mo

First Year Total:           $22,459
```

---

## Key Insights

1. **Edge Compute scales efficiently:**
   - 1 device per location
   - NOT 1 per robot
   - Supports multiple robots/drones

2. **Monthly costs DON'T scale with fleet:**
   - $135/mo whether you have 1 robot or 10
   - One brain (Claude) controls all
   - One FABRIC network for coordination

3. **LiveKit self-hosted = FREE:**
   - Unlimited video streams
   - All robots/drones
   - No per-stream costs

4. **Best ROI progression:**
   ```
   Start:    1 robot (you have) + 1 edge device ($997)
   Phase 2:  +1 drone ($759)
   Phase 3:  +1 robot ($1,600-$3,500)
   Phase 4:  +1 location edge ($997) + robots
   ```

5. **Drones are CHEAP force multipliers:**
   - $759 drone = 10x coverage vs. ground robot
   - Aerial view invaluable
   - Easy to add to existing fleet

---

## Recommended Action Plan

### Week 1-2: Foundation
**Cost:** $997 (edge compute only)
- Deploy OpenMind Edge for your existing Unitree
- Test voice control: "Legacy, move robot to barn"
- Validate system before expanding fleet

### Week 3-4: Add Aerial Capability
**Cost:** +$759 (drone)
- Add DJI Mini 4 Pro
- Implement drone patrol tools
- Test coordinated ground + air patrols

### Month 2: Expand Ground Fleet
**Cost:** +$1,600-$3,500 (2nd robot)
- Add Unitree Go2 or Go1
- Test multi-robot coordination
- Build automated workflows

### Month 3+: Scale to Multi-Location
**Cost:** +$997 per location
- Add edge devices to other locations
- Full farm/home/property coverage
- Fleet coordination via FABRIC

---

*Document created to answer multi-robot fleet questions*
*Last updated: 2026-01-22*


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]
