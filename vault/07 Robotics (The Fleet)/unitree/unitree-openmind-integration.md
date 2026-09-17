# Unitree Robot + OpenMind Integration Plan

> **Created:** 2026-01-22
> **Status:** Research & Planning
> **Priority:** High
> **Estimated Cost:** $5,000-$7,000 (hardware + services)
> **Timeline:** 4-6 weeks

---

## Executive Summary

Integration of Unitree robot with OneMind OS using OpenMind's OM1 platform to enable:
- Voice control via Legacy AI
- Autonomous navigation and task execution
- Real-time video streaming to UI
- Location tracking and presence awareness
- Coordinated fleet management (future)

---

## Part 1: OpenMind Platform Analysis

### What OpenMind Provides

**OpenMind** (openmind.org):
- $20M funded (Pantera, Coinbase Ventures, DCG)
- Stanford spin-off
- "Android for Robots" - universal robot OS

**Core Products:**

1. **OM1 (Robot OS)**
   - AI-native operating system
   - Hardware-agnostic (works with Unitree, Boston Dynamics, Figure, etc.)
   - LLM-agnostic (Claude, GPT, Gemini, local models)
   - Open source: https://github.com/openmind-agi/om1

2. **FABRIC Network**
   - Decentralized robot coordination
   - Shared world model across fleet
   - Sub-second latency coordination
   - Identity and location services

3. **Edge Compute Device ($997 - Assumed)**
   Based on industry standards, this likely includes:
   - **Hardware:**
     - NVIDIA Jetson Orin Nano/NX (8GB-16GB)
     - WiFi 6E + Ethernet
     - USB-C for robot connection
     - Mounting bracket for Unitree
   - **Software:**
     - Pre-installed OM1
     - Pre-configured FABRIC node
     - Docker runtime
     - Remote management dashboard
   - **Services:**
     - 1 year FABRIC network access
     - Cloud dashboard subscription
     - OTA updates
     - Developer support

---

## Part 2: Current OneMind OS Status

### ✅ What's Working

| Component | Status | Notes |
|-----------|--------|-------|
| Legacy AI Brain | ✅ Live | Claude Opus 4.5 via Bedrock |
| Voice Integration | ✅ Live | LiveKit + Deepgram + Cartesia |
| NATS Messaging | ✅ Live | Cloud hub operational |
| UI (Web/Mobile) | ✅ Live | Next.js + Flutter (24 screens) |
| Camera Integration | 🟡 Partial | Frigate ready, not deployed |
| Physical Presence | ✅ Live | Room tracking system |

### ❌ What's Missing for Unitree

| Component | Status | Blocker |
|-----------|--------|---------|
| RemoteAgent Software | ❌ Not built | Need to create `agents/edge/` |
| Unitree Control Tools | ❌ Not built | Need SDK integration |
| NATS Leaf Node | ❌ Not deployed | Need edge NATS setup |
| Robot Video Streaming | ❌ Not integrated | Frigate not deployed |
| Robot Location Tracking | ❌ Not implemented | GPS/odometry needed |
| OM1 Integration | ❌ Not started | API bridge needed |

---

## Part 3: Integration Architecture

### Proposed System Design

```
┌─────────────────────────────────────────────────────────────────────┐
│                    UNITREE + ONEMIND INTEGRATION                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  USER (You)                                                          │
│    │                                                                 │
│    │ Voice: "Legacy, have the robot patrol the barn"               │
│    ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  CLOUD (AWS - 100.102.21.44)                                  │   │
│  │                                                               │   │
│  │  ┌────────────┐    ┌────────────┐    ┌────────────┐         │   │
│  │  │  Legacy    │───►│   NATS     │    │  Flutter   │         │   │
│  │  │  (Agno)    │    │   Hub      │    │    UI      │         │   │
│  │  └────────────┘    └──────┬─────┘    └────────────┘         │   │
│  │                           │                                   │   │
│  └───────────────────────────┼───────────────────────────────────┘   │
│                              │ NATS over Tailscale                   │
│                              ▼                                        │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  EDGE (Barn/Farm)                                            │   │
│  │                                                               │   │
│  │  ┌────────────┐    ┌────────────┐    ┌────────────┐         │   │
│  │  │   NATS     │◄──►│  OpenMind  │◄──►│  Unitree   │         │   │
│  │  │   Leaf     │    │    OM1     │    │   Robot    │         │   │
│  │  └────────────┘    └────────────┘    └────────────┘         │   │
│  │                           │                  │                │   │
│  │  ┌────────────┐    ┌─────▼──────┐    ┌─────▼──────┐         │   │
│  │  │ RemoteAgent│    │   Frigate  │    │   Robot    │         │   │
│  │  │  (Python)  │    │    NVR     │    │   Cameras  │         │   │
│  │  └────────────┘    └────────────┘    └────────────┘         │   │
│  │                                                               │   │
│  └───────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  Message Flow:                                                       │
│  1. User speaks → Legacy (voice transcribed)                        │
│  2. Legacy decides action → NATS Hub                                │
│  3. NATS routes to Edge → NATS Leaf                                 │
│  4. RemoteAgent receives → OM1 API                                  │
│  5. OM1 controls → Unitree SDK                                      │
│  6. Robot executes → Sends video/telemetry back                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Part 4: Implementation Requirements

### 4.1 Hardware Needed

| Item | Source | Est. Cost | Purpose |
|------|--------|-----------|---------|
| **Unitree Robot** | Owned | $0 | Target robot |
| **OpenMind Edge Compute** | openmind.org/store | $997 | Pre-configured OM1 device |
| **OR DIY Edge Compute:** |
| - NVIDIA Jetson Orin Nano | Amazon | $499 | Edge AI compute |
| - Coral USB TPU (optional) | Amazon | $60 | Vision acceleration |
| - WiFi 6E USB Adapter | Amazon | $40 | Fast connectivity |
| **Network Infrastructure** |
| - Tailscale Pro (optional) | Tailscale | $60/yr | Secure networking |
| **Camera Upgrades (optional)** |
| - Additional Yi Cameras | Amazon | $30 each | More coverage |

**Total Hardware Cost:**
- **Option A (OpenMind):** $997 + $60/yr = ~$1,057
- **Option B (DIY):** $499 + $60 + $40 + $60 = ~$659

**Recommendation:** Start with **Option A (OpenMind $997 device)** for plug-and-play experience, then optimize later if needed.

### 4.2 Software Licenses & Services

| Service | Provider | Cost | Purpose |
|---------|----------|------|---------|
| OpenMind FABRIC | OpenMind | Included (1yr) | Robot coordination |
| AWS Bedrock | AWS | ~$50/mo | Claude Opus LLM |
| Deepgram STT | Deepgram | ~$20/mo | Voice transcription |
| Cartesia TTS | Cartesia | ~$10/mo | Voice synthesis |
| Tailscale | Tailscale | Free (personal) | VPN networking |
| LiteLLM (optional) | Self-hosted | Free | Model routing |

**Total Monthly Services:** ~$80/month

### 4.3 Development Time (Claude Code)

Assuming Claude Code does ALL implementation:

| Phase | Tasks | Est. Time |
|-------|-------|-----------|
| **Phase 1: RemoteAgent** | Build edge software | 1 week |
| **Phase 2: OM1 Integration** | API bridge, tools | 1 week |
| **Phase 3: NATS Leaf** | Deploy edge networking | 3 days |
| **Phase 4: Video Streaming** | Frigate + UI integration | 1 week |
| **Phase 5: Testing** | E2E validation | 1 week |

**Total Timeline:** 4-6 weeks (with Claude Code doing the work)

---

## Part 5: What You Get After Integration

### 5.1 Immediate Capabilities

**Voice Control:**
```
You: "Legacy, where's the robot?"
Legacy: "The robot is in the barn, battery at 78%. Want me to send it somewhere?"

You: "Have it patrol the greenhouse and report any issues"
Legacy: "Sending the robot to greenhouse now. I'll stream video to your dashboard."
```

**Autonomous Tasks:**
- Patrol routes (perimeter, barn, greenhouse)
- Object delivery (tools, supplies)
- Visual inspection (plants, equipment, animals)
- Security monitoring (motion detection, alerts)
- Follow mode (follows you around)

**Real-Time Awareness:**
- Robot location on map (GPS + odometry)
- Battery level monitoring
- Video stream in UI
- Task progress updates
- Collision avoidance status

### 5.2 Advanced Use Cases

**Farm Automation:**
```
Morning routine (automated):
1. Robot leaves charging dock at 7am
2. Patrols greenhouse, checks plant health
3. Inspects chicken coop, counts chickens
4. Checks perimeter fences for damage
5. Returns to dock, reports findings to Legacy
6. Legacy sends you voice summary: "Morning report ready"
```

**Voice-Commanded Inspections:**
```
You: "Legacy, I heard a noise in the barn. Can you check it out?"
Legacy: "Dispatching robot to barn. Streaming video now."
[Robot arrives, Legacy analyzes video]
Legacy: "It's just the cats playing. Everything looks normal."
```

**Coordinated Multi-Robot (Future):**
```
# When you add more robots:
You: "Legacy, we need to move all the hay bales from barn to storage"
Legacy: "I'll coordinate 3 robots for optimal efficiency. ETA 45 minutes."
[Robots work together via FABRIC network]
```

### 5.3 Integration with Existing OneMind OS

**Unified Dashboard:**
- Robot status card on main dashboard
- Live video feed widget
- Battery/health indicators
- Task queue and history

**Memory Integration:**
- Robot observations saved to Legacy's memory
- "Remember when the robot saw X" queries
- Pattern recognition (e.g., "cats in barn every evening")

**Workflow Automation:**
- Morning patrol → automated workflow
- Emergency response → robot dispatched automatically
- Low battery → return to dock triggered

**Presence Awareness:**
- Robot location tracked same as your presence
- "Where is everything?" shows you + robot + devices
- Collision avoidance (robot knows where you are)

---

## Part 6: Pricing Breakdown (Complete Deployment)

### One-Time Costs

| Item | Cost | Required? |
|------|------|-----------|
| OpenMind Edge Compute | $997 | Yes |
| Unitree Robot | Owned | - |
| Additional Cameras | $0-150 | Optional |
| **Total One-Time** | **$997-$1,147** | |

### Recurring Costs (Monthly)

| Service | Monthly | Annual |
|---------|---------|--------|
| AWS Bedrock (Claude) | $50 | $600 |
| Deepgram STT | $20 | $240 |
| Cartesia TTS | $10 | $120 |
| Tailscale Pro (optional) | $5 | $60 |
| OpenMind FABRIC (Year 2+) | ~$50 | ~$600 |
| **Total Recurring** | **$135/mo** | **$1,620/yr** |

### Development Cost (If Hiring vs. Claude)

| Approach | Cost | Timeline |
|----------|------|----------|
| Hire Developer | $50K-$100K | 3-6 months |
| **Claude Code** | **$20/mo** | **4-6 weeks** |
| DIY (Your Time) | Free | 6-12 months |

**Total First Year Cost:** ~$2,617 (hardware + services + Claude)

---

## Part 7: Vision Examples - What This Enables

### Example 1: Autonomous Farm Management

```
MORNING (7:00 AM - Automated):
─────────────────────────────────────────────────────────
Robot: Leaves charging dock
Legacy: Executes morning patrol workflow:
  1. Greenhouse inspection
     - Takes photos of all plants
     - Measures soil moisture via IR
     - Checks temperature sensors
  2. Chicken coop check
     - Counts chickens (computer vision)
     - Checks water level
     - Opens coop door if closed
  3. Perimeter patrol
     - Checks fences for damage
     - Looks for animal tracks
     - Inspects gates

Legacy: Analyzes all data, sends you summary:
  "Morning report: All 12 chickens accounted for.
   Greenhouse temp at 72°F. Tomatoes need watering today.
   Found rabbit tracks near south fence. Everything else normal."
```

### Example 2: Voice-Controlled Assistance

```
YOU: "Legacy, I'm working in the greenhouse and need my pruning shears.
      They're on the workbench in the barn."

LEGACY: "I'll send the robot to get them."
        [Robot navigates barn → workbench → greenhouse]
        "Robot is heading to you now with the shears. ETA 2 minutes."

[Robot arrives, you take shears]

YOU: "Thanks. Can you have it follow me for the next hour?
     I might need more tools."

LEGACY: "Sure, robot is now in follow mode."
        [Robot follows 6 feet behind you using SLAM + vision]
```

### Example 3: Security & Emergency Response

```
SCENARIO: Motion detected at 2 AM

Frigate NVR: Detects motion near barn
    ↓
NATS Event: events.motion.barn
    ↓
Legacy: Evaluates (awareness mode: DORMANT but motion = important)
    ↓
Legacy: Dispatches robot to investigate
        Switches robot cameras to IR night vision
        Streams video to cloud
        Analyzes video with vision model
    ↓
Legacy: Identifies coyote near chicken coop
        Determines: NOT emergency (coop is secure)
        Sends robot to make noise, scare away
    ↓
YOU: Receive notification:
     "Coyote near coop at 2:15am. Robot scared it off.
      Chickens safe. Video saved."
```

### Example 4: Multi-Robot Coordination (Future)

```
SCENARIO: You have 3 Unitree robots (2 Go1, 1 Go2)

YOU: "Legacy, we need to clear all the fallen branches
     from the pasture before the horses go out."

LEGACY: "I'll coordinate all 3 robots for optimal coverage."
        ↓
        [Via FABRIC Network - shared world model]
        ↓
        Robot 1: Patrols west section
        Robot 2: Patrols center section
        Robot 3 (Go2): Patrols east section + carries collection bin
        ↓
        [Robots detect branches, mark locations, Go2 collects]
        ↓
LEGACY: "Pasture cleared. Found 23 branches, all collected.
         Safe for horses now. Took 18 minutes."
```

### Example 5: Integration with Full OneMind Ecosystem

```
MORNING SEQUENCE (Fully Integrated):

6:30 AM - Oura Ring detects you waking up
    ↓
Legacy: Starts morning workflow
    ↓
    • Adjusts home thermostat (Home Assistant)
    • Starts coffee maker (IoT)
    • Dispatches robot to start farm patrol
    • Checks weather forecast
    • Reviews your calendar
    ↓
7:00 AM - You enter kitchen
    ↓
Legacy (via speaker): "Good morning! Coffee's ready.
                       Robot just finished greenhouse check -
                       everything looks good.

                       You have 3 meetings today.
                       Want your morning briefing?"
    ↓
YOU: "Yes please"
    ↓
Legacy: Provides briefing while showing:
    • Calendar on tablet
    • Robot patrol video clips
    • Plant health dashboard
    • Weather forecast
    • Todoist priority tasks
```

---

## Part 8: Technical Specifications

### Unitree Robot Models (Compatibility)

| Model | Price | OM1 Support | Best For |
|-------|-------|-------------|----------|
| **Go1 EDU** | $3,500 | ✅ Yes | Indoor/outdoor patrol, learning |
| **Go2** | $1,600-$5,000 | ✅ Yes | Lightweight tasks, education |
| **B1** | $50,000+ | ✅ Yes | Heavy payload, industrial |
| **H1** | TBD | 🟡 Future | Humanoid tasks |

### Performance Specs (Expected)

| Metric | Target | Notes |
|--------|--------|-------|
| Command latency | <500ms | Voice → robot movement |
| Video stream delay | <200ms | Robot camera → UI |
| Battery life | 2-4 hours | Depends on tasks |
| Patrol speed | 1-3 m/s | Unitree max speed |
| NATS message latency | <5ms | Edge → Cloud |
| GPS accuracy | <2m | With RTK: <2cm |

### API Endpoints (To Be Built)

```python
# Legacy → OM1 Control
POST /api/robot/command
{
  "robot_id": "unitree_go1_001",
  "action": "navigate",
  "target": {"lat": 37.123, "lon": -122.456},
  "speed": "normal"
}

# OM1 → Legacy Status
WebSocket /api/robot/telemetry
{
  "robot_id": "unitree_go1_001",
  "battery": 78,
  "location": {"lat": 37.123, "lon": -122.456},
  "speed": 1.2,
  "status": "navigating"
}

# Video Stream
GET /api/robot/video/{robot_id}
→ Returns: HLS/WebRTC stream URL
```

---

## Part 9: Risk Assessment

### Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| OM1 API incompatible | Low | High | Use open source OM1, fork if needed |
| Unitree SDK limitations | Medium | Medium | Direct ROS integration fallback |
| Network latency | Low | Medium | Local NATS leaf reduces latency |
| Robot navigation fails | Medium | Low | Add manual control override |
| Video streaming drops | Medium | Low | Record locally, sync later |

### Operational Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Robot gets stuck | High | Low | Timeout + return to dock |
| Battery dies mid-task | Medium | Medium | Always monitor battery, return at 20% |
| Weather damage | Low | High | Don't operate in rain/snow (Go1 not waterproof) |
| Animal interference | Medium | Low | Robot horn/noise to scare |

### Cost Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| API costs exceed budget | Low | Medium | Set billing alerts at $100/mo |
| FABRIC subscription increases | Medium | Low | Self-host NATS alternative |
| Hardware failure | Low | High | Warranty coverage + backup plan |

---

## Part 10: Decision Matrix

### Should You Proceed?

**YES, if:**
- ✅ You want immediate voice-controlled robot deployment
- ✅ You're comfortable with $1K hardware + $135/mo services
- ✅ You have Unitree robot ready to integrate
- ✅ You have 4-6 weeks for Claude Code to implement
- ✅ You want unified control of all farm/home systems

**WAIT, if:**
- ❌ Budget constraints (<$2K first year total)
- ❌ Unitree robot not yet purchased
- ❌ Higher priority OneMind OS features needed first
- ❌ Prefer to DIY everything (6-12 months timeline)

**ALTERNATIVE APPROACH:**
Start with **Phase 1 only** (RemoteAgent + NATS Leaf) without OM1:
- Cost: $0-$500 (DIY edge device)
- Timeline: 1-2 weeks
- Capability: Basic robot control via direct SDK
- Upgrade Path: Add OM1 later for advanced features

---

## Part 11: Next Steps

### Immediate Actions (This Week)

1. **Research OpenMind Pricing**
   - [ ] Visit openmind.org/store and confirm $997 product details
   - [ ] Contact sales for developer/education discount
   - [ ] Ask about FABRIC network pricing after Year 1

2. **Validate Unitree Integration**
   - [ ] Check Unitree model compatibility with OM1
   - [ ] Verify SDK access (Unitree Developer account)
   - [ ] Test basic Unitree connectivity

3. **Budget Approval**
   - [ ] Confirm $997 hardware budget
   - [ ] Confirm $135/mo service budget
   - [ ] Set AWS billing alert at $100/mo

### Phase 1: RemoteAgent (Week 1-2)

- [ ] Create `agents/edge/` directory structure
- [ ] Build RemoteAgent base class
- [ ] Implement NATS leaf connection
- [ ] Create shell command tools
- [ ] Deploy to test Jetson/Pi
- [ ] Test NATS cloud ↔ edge messaging

### Phase 2: OM1 Integration (Week 3-4)

- [ ] Install OM1 on edge device
- [ ] Create OM1 API client
- [ ] Build robot control tools (@tool functions)
- [ ] Add to Legacy agent toolset
- [ ] Test voice → robot command flow

### Phase 3: Full Deployment (Week 5-6)

- [ ] Deploy Frigate NVR for cameras
- [ ] Integrate robot video streams
- [ ] Build robot dashboard UI
- [ ] Add location tracking
- [ ] Create automated workflows
- [ ] End-to-end testing

---

## Part 12: Success Criteria

### Must Have (MVP)

- [ ] Voice command: "Move robot to barn" works
- [ ] Robot status visible in UI
- [ ] Battery level monitoring
- [ ] Manual emergency stop
- [ ] Return to dock command

### Should Have (V1)

- [ ] Live video stream in UI
- [ ] Autonomous patrol routes
- [ ] Location tracking on map
- [ ] Automated morning routine
- [ ] Collision avoidance

### Nice to Have (V2)

- [ ] Multi-robot coordination
- [ ] Computer vision analysis
- [ ] Voice from robot speakers
- [ ] Follow mode
- [ ] Advanced pathfinding

---

## Appendix A: Relevant Links

- **OpenMind:** https://openmind.org/
- **OM1 GitHub:** https://github.com/openmind-agi/om1
- **Unitree Robotics:** https://www.unitree.com/
- **OneMind OS Docs:** `docs/architecture/ONEMIND-ARCHITECTURE-PLAN.md`
- **RemoteAgent Architecture:** `docs/architecture/REMOTEAGENT-ARCHITECTURE.md`
- **NATS Multi-Location:** `docs/architecture/MULTI-LOCATION-NATS.md`

## Appendix B: Contact Info

- **OpenMind Sales:** (from website - need to verify)
- **Unitree Support:** support@unitree.com
- **NATS Community:** https://natsio.slack.com/
- **OneMind OS:** This repo

---

*Document created for Linear issue OMOS-88*
*Last updated: 2026-01-22*


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]
