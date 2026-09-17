# Offline Services & Resilience Infrastructure - Complete Overview

> **Master reference document for all offline/resilience systems**
>
> Last Updated: 2026-01-23

---

## Table of Contents

1. [What is This?](#what-is-this)
2. [The Problem We're Solving](#the-problem-were-solving)
3. [The Complete Solution](#the-complete-solution)
4. [System 1: Meshtastic Mesh Network](#system-1-meshtastic-mesh-network)
5. [System 2: OneMind Emergency Kit](#system-2-onemind-emergency-kit)
6. [How They Work Together](#how-they-work-together)
7. [All Documentation](#all-documentation)
8. [All Linear Issues](#all-linear-issues)
9. [Cost Summary](#cost-summary)
10. [Implementation Timeline](#implementation-timeline)
11. [Quick Reference](#quick-reference)

---

## What is This?

**Offline Services** is a **4-layer resilience architecture** that ensures OneMind OS continues operating even when internet, power, or infrastructure fails.

**Two complementary systems:**

1. **Meshtastic Mesh Network** (OMOS-132-134)
   - Keeps THE FLEET operating (drones, robots, sensors)
   - LoRa wireless mesh communication
   - 5-10km range across property

2. **OneMind Emergency Kit** (OMOS-135-136)
   - Keeps YOU informed (human operator)
   - Portable AI + knowledge on USB drive
   - Works on any laptop

**Together:** Complete offline capability for both autonomous systems and human operators.

---

## The Problem We're Solving

### Current Architecture Dependency

```
OneMind OS Today:
┌─────────────────────────────────────┐
│     Cloud NATS (Azure VM)      │ ← Single point of failure
└─────────────────────────────────────┘
                ↓
           (Internet) ← Can fail
                ↓
┌─────────────────────────────────────┐
│        Edge Devices (Farm)          │
└─────────────────────────────────────┘
                ↓
        Drones, Robots, Sensors
```

**What happens when:**
- ❌ Internet goes down → Fleet loses cloud connection
- ❌ Edge devices fail → Local control lost
- ❌ Power outage → All systems offline
- ❌ AWS outage → Cloud services unavailable

**Result:** Complete system failure

---

## The Complete Solution

### 4-Layer Resilience Architecture

```
Layer 1: Cloud NATS (Primary)
    ↓ (Internet fails)
Layer 2: Meshtastic Mesh (Fleet Backup) ← NEW
    ↓ (Edge devices fail)
Layer 3: Edge Ollama (Local AI)
    ↓ (Everything fails)
Layer 4: Emergency USB Kit (Personal Backup) ← NEW
```

**Key Benefits:**
- ✅ **Zero single points of failure**
- ✅ **Automatic failover** at each layer
- ✅ **Fleet keeps operating** during outages
- ✅ **Human stays informed** when everything fails
- ✅ **Gradual degradation** instead of total failure

---

## System 1: Meshtastic Mesh Network

### FOR THE FLEET: Keeps drones/robots operating

**Linear Issues:** OMOS-132, OMOS-133, OMOS-134

### What It Does

Provides backup communication network using LoRa radio technology when internet fails:

```
Normal Operation:
Drone → NATS Leaf → Cloud NATS → Legacy Agent

Internet Fails:
Drone → Meshtastic → House Hub → Local Legacy ✅
```

### Key Components

**Hardware:**
- RAK19007 LoRa radio modules ($65 each)
- LILYGO T3-S3 WiFi bridges ($35 each)
- GPS modules for tracking
- Solar panels for remote sensors

**Software:**
- Meshtastic firmware (open source)
- MQTT bridge to NATS
- Automatic failover logic
- RAG document search

**Coverage:**
- House hub (primary)
- Barn relay (secondary)
- Field sensors (×10)
- Mobile assets (drones, robots)

### Use Cases

**1. Drone Telemetry Backup**
```
Zeus drone 3km away → Internet fails
Without Mesh: ❌ Lost telemetry, cannot command
With Mesh: ✅ Telemetry continues, remote recall works
```

**2. Robot Status Updates**
```
Unitree robot beyond WiFi range
Without Mesh: ❌ No status, offline
With Mesh: ✅ Status reports, commands work
```

**3. Critical Security Alerts**
```
Skydio detects intruder
Without Mesh: ❌ Alert may fail if internet down
With Mesh: ✅ Dual-path (internet + mesh) guaranteed delivery
```

**4. Property-Wide Sensors**
```
10 sensors (soil, weather, gates) across 500 acres
Without Mesh: ❌ Cellular = $100/month ongoing
With Mesh: ✅ LoRa = $0/month, self-powered
```

### Architecture Diagram

```
┌──────────────────────────────────────────────────────┐
│              Cloud (Azure VM)                         │
│   ┌──────────┐      ┌──────────┐                    │
│   │  Legacy  │←────→│ NATS Hub │                    │
│   └──────────┘      └──────────┘                    │
└──────────────────────────────────────────────────────┘
                         ↓
                  (Internet - can fail)
                         ↓
┌──────────────────────────────────────────────────────┐
│         House Edge Device (Jetson Nano)               │
│   ┌──────────┐    ┌──────────────────────┐          │
│   │   NATS   │←──→│ Meshtastic Bridge    │          │
│   │   Leaf   │    │ (MQTT ↔ NATS)        │          │
│   └──────────┘    └──────────────────────┘          │
│                            ↓                          │
│                   ┌─────────────────┐                │
│                   │ Meshtastic RAK  │                │
│                   │ (LoRa Radio)    │                │
│                   └─────────────────┘                │
└──────────────────────────────────────────────────────┘
                         ↓
              (LoRa Mesh - always works)
                         ↓
        ┌────────────────┼────────────────┐
        ↓                ↓                ↓
  ┌──────────┐    ┌──────────┐    ┌──────────┐
  │  Zeus    │    │ Unitree  │    │  Field   │
  │  Drone   │    │ Robot    │    │ Sensors  │
  │ (LoRa)   │    │ (LoRa)   │    │ (LoRa)   │
  └──────────┘    └──────────┘    └──────────┘
```

### Cost Breakdown

| Phase | Items | Cost |
|-------|-------|------|
| **Starter** | 2× RAK kits, batteries | $180 |
| **Production** | 3 drones + 2 robots + 10 sensors | $1,235 |
| **Monthly** | Zero (no subscriptions) | $0 |
| **Total** | Complete property coverage | $1,415 |

### Documentation

**Complete Guides (2,400+ lines):**
- [MESHTASTIC-MESH-INTEGRATION.md](MESHTASTIC-MESH-INTEGRATION.md) - Architecture & code
- [MESHTASTIC-STARTER-GUIDE.md](MESHTASTIC-STARTER-GUIDE.md) - Hardware & setup
- [MESHTASTIC-USE-CASES.md](MESHTASTIC-USE-CASES.md) - Real-world scenarios

**Scripts:**
- `scripts/linear/create_meshtastic_issues.py` - Create Linear issues
- `scripts/linear/update_meshtastic_issues.py` - Add labels/project

---

## System 2: OneMind Emergency Kit

### FOR YOU: Keeps human operator informed

**Linear Issues:** OMOS-135, OMOS-136

### What It Does

Self-contained OneMind OS on USB drive - plug into any laptop for offline AI and knowledge access:

```
All Systems Down:
├─ Internet: ❌ Down
├─ Edge Devices: ❌ Failed
├─ Mesh Network: ❌ Out of range
└─ Emergency USB Kit: ✅ Still works!

Plug USB into laptop → Instant access to:
├─ 3× AI models (Phi-3, Llama, Qwen)
├─ 4GB survival manuals (medical, engineering, agriculture)
├─ OneMind OS documentation
├─ Meshtastic firmware/tools
└─ NATS portable server
```

### Key Components

**Hardware:**
- 256GB-512GB USB drive or SSD
- That's it! ($30-75)

**Software (all included on USB):**
- Ollama portable + 3 AI models (6GB)
- Knowledge base with RAG (4GB)
- Essential tools (Meshtastic, NATS, Python)
- Self-contained web UI
- Cross-platform startup scripts

**AI Models:**
- Phi-3 Mini (3.8B) - Best for reasoning
- Llama 3.2 (3B) - Balanced performance
- Qwen 2.5 (3B) - Best for coding

### Use Cases

**1. Bug-Out Bag Essential**
```
Natural disaster, evacuate to shelter
└─ Plug USB into any laptop at shelter
   ├─ Ask AI: "How to treat wounds without hospital?"
   ├─ Search: "Water purification methods"
   ├─ Reference: "Edible plant identification"
   └─ Guide: "DIY shelter construction"
```

**2. Remote Field Operations**
```
Installing edge devices in remote location (no internet)
└─ Plug USB into laptop
   ├─ Reference OneMind OS docs offline
   ├─ Troubleshoot NATS configuration
   ├─ Flash Meshtastic firmware
   └─ Test configurations locally
```

**3. System Recovery**
```
All infrastructure down
└─ Plug USB into laptop
   ├─ Ask AI: "How to restore NATS cluster?"
   ├─ Search: "Meshtastic troubleshooting"
   ├─ Reference: OneMind architecture diagrams
   └─ Rebuild system from scratch
```

**4. Learning/Training**
```
Teaching OneMind OS to new team member
└─ Hand them USB
   ├─ Works on their laptop
   ├─ No cloud account needed
   ├─ Instant demo
   └─ Complete documentation
```

### Architecture Diagram

```
┌─────────────────────────────────────────────────────┐
│          USB Drive (256GB-512GB)                     │
│                                                       │
│  ┌────────────────────────────────────────────┐    │
│  │         Portable Ollama (6GB)               │    │
│  │  ├── Phi-3 Mini (3.8B) - 2.3GB            │    │
│  │  ├── Llama 3.2 (3B) - 1.9GB               │    │
│  │  └── Qwen 2.5 (3B) - 1.9GB                │    │
│  └────────────────────────────────────────────┘    │
│                                                       │
│  ┌────────────────────────────────────────────┐    │
│  │         Knowledge Base (4GB)                │    │
│  │  ├── OneMind OS Documentation (500MB)      │    │
│  │  ├── Survival Manuals (2GB)                │    │
│  │  │   ├── Field Medicine                    │    │
│  │  │   ├── Engineering                       │    │
│  │  │   ├── Agriculture                       │    │
│  │  │   └── Tactical                          │    │
│  │  ├── Technical References (1GB)            │    │
│  │  └── Vector Database (1GB)                 │    │
│  └────────────────────────────────────────────┘    │
│                                                       │
│  ┌────────────────────────────────────────────┐    │
│  │         Essential Tools (400MB)             │    │
│  │  ├── Meshtastic CLI + Firmware             │    │
│  │  ├── NATS Portable Server                  │    │
│  │  ├── Python Portable                        │    │
│  │  └── ChromaDB (RAG)                         │    │
│  └────────────────────────────────────────────┘    │
│                                                       │
│  ┌────────────────────────────────────────────┐    │
│  │         Web Interface (50MB)                │    │
│  │  ├── Chat UI                                │    │
│  │  ├── Document Search                        │    │
│  │  ├── Tool Access                            │    │
│  │  └── System Status                          │    │
│  └────────────────────────────────────────────┘    │
│                                                       │
│  Startup: Double-click start.bat or start.sh        │
│  Browser: http://localhost:8000                     │
└─────────────────────────────────────────────────────┘
```

### Cost Breakdown

| Item | Option | Cost |
|------|--------|------|
| **USB Drive** | SanDisk 256GB | $30 |
| | Samsung T7 SSD 500GB | $75 |
| **Software** | Ollama + models | FREE |
| | Knowledge base | FREE |
| | Tools | FREE |
| **Total** | Complete system | $30-75 |

**vs. Bunker AI:** $89 + unknown specs + proprietary

### Documentation

**Complete Guide (1,600+ lines):**
- [ONEMIND-EMERGENCY-KIT.md](ONEMIND-EMERGENCY-KIT.md) - Complete build guide

**Scripts:**
- `scripts/linear/create_emergency_kit_issues.py` - Create Linear issues

**Bunker AI Research (pending):**
- Will create `BUNKER-AI-ANALYSIS.md` after purchase

---

## How They Work Together

### Complementary Systems

**They solve different problems:**

```
┌────────────────────────────────────────────────────┐
│              Problem: Offline Operation             │
└────────────────────────────────────────────────────┘
                        ↓
        ┌───────────────┴───────────────┐
        ↓                               ↓
┌─────────────────┐           ┌─────────────────┐
│   Meshtastic    │           │  Emergency Kit   │
│  (Fleet Comms)  │           │ (Human Knowledge)│
└─────────────────┘           └─────────────────┘
        ↓                               ↓
  Keeps robots                    Keeps you
  operating                       informed
```

### Failure Cascade Example

**Scenario: Total Infrastructure Failure**

```
Hour 0: Internet Outage
├─ Layer 1 (Cloud NATS): ❌ Down
└─ Layer 2 (Meshtastic): ✅ Activates
   ├─ Drones switch to mesh telemetry
   ├─ Robots receive commands via LoRa
   └─ Sensors keep transmitting
   Result: FLEET STILL OPERATING

Hour 6: Power Outage at Edge Devices
├─ Layer 2 (Meshtastic): ❌ Radios dead
├─ Layer 3 (Edge Ollama): ❌ No power
└─ Layer 4 (Emergency Kit): ✅ Activates
   ├─ Plug USB into laptop (battery powered)
   ├─ Ask AI: "How to restore power?"
   ├─ Search: "Generator troubleshooting"
   └─ Reference: OneMind architecture
   Result: YOU CAN SOLVE PROBLEMS

Hour 12: Power Restored
├─ Edge devices boot up
├─ Meshtastic network recovers
├─ Internet still down, but mesh works
└─ Fleet resumes operation via Layer 2
   Result: GRADUAL RECOVERY
```

### When to Use What

| Situation | Use Meshtastic | Use Emergency Kit |
|-----------|----------------|-------------------|
| Internet down | ✅ Yes | Maybe (if need docs) |
| Edge devices offline | No (needs power) | ✅ Yes |
| In the field (drone mission) | ✅ Yes | No (not needed) |
| Bug-out scenario | Maybe (if you bring radios) | ✅ Yes (in bag) |
| System troubleshooting | No (won't help) | ✅ Yes (docs/AI) |
| Learning OneMind OS | No (not relevant) | ✅ Yes (all docs) |
| Fleet operations | ✅ Yes | No (for humans) |
| Survival situation | Maybe (comms) | ✅ Yes (knowledge) |

---

## All Documentation

### Master Documents (This File)
- **OFFLINE-SERVICES-OVERVIEW.md** ← YOU ARE HERE
  - Complete overview of both systems
  - Use this as starting point for review

### Meshtastic Documentation (2,400+ lines)
1. **MESHTASTIC-MESH-INTEGRATION.md** (1,000+ lines)
   - Complete architecture
   - MQTT bridge implementation
   - Code examples
   - Emergency scenarios

2. **MESHTASTIC-STARTER-GUIDE.md** (800+ lines)
   - Hardware selection (RAK, LILYGO, Heltec)
   - Step-by-step setup
   - Range testing procedures
   - Troubleshooting guide

3. **MESHTASTIC-USE-CASES.md** (600+ lines)
   - Drone telemetry backup
   - Robot status updates
   - Critical security alerts
   - Property-wide sensors
   - 10 real-world scenarios

### Emergency Kit Documentation (1,600+ lines)
1. **ONEMIND-EMERGENCY-KIT.md**
   - Complete build guide
   - All components explained
   - Implementation plan (6-8 weeks)
   - Startup scripts (Windows/Mac/Linux)
   - Cost breakdown

### Pending Documentation
1. **BUNKER-AI-ANALYSIS.md** (after purchase)
   - Technical teardown
   - Performance benchmarks
   - Content inventory
   - Comparison to OneMind Kit

---

## All Linear Issues

### Project: "Offline Services & Resilience Infrastructure"

**Meshtastic Subsection (Fleet Communication):**

| Issue | Title | Status | Cost |
|-------|-------|--------|------|
| **OMOS-132** | Meshtastic Mesh Network Integration | 🔴 Not Started | - |
| **OMOS-133** | Meshtastic Starter Kit & Property Testing | 🔴 Not Started | $180 |
| **OMOS-134** | Fleet Integration with Mesh Backup | 🔴 Not Started | $1,235 |

**Emergency Kit Subsection (Personal Knowledge):**

| Issue | Title | Status | Cost |
|-------|-------|--------|------|
| **OMOS-135** | Research & Purchase Bunker AI | 🔴 Not Started | $89 |
| **OMOS-136** | Build OneMind Emergency Kit | 🔴 Not Started | $30-75 |

**Total Issues:** 5
**Total Cost:** $1,534-1,579

### To Create Issues

```bash
# Navigate to project root
cd ~/Documents/onemind-os

# Create Meshtastic issues (OMOS-132-134)
python3 scripts/linear/create_meshtastic_issues.py

# Create Emergency Kit issues (OMOS-135-136)
python3 scripts/linear/create_emergency_kit_issues.py

# Organize with labels and project
python3 scripts/linear/update_meshtastic_issues.py
```

---

## Cost Summary

### Meshtastic Mesh Network

**Starter (Testing):**
- 2× RAK19007 kits: $130
- Batteries: $50
- **Total: $180**

**Full Deployment:**
- Starter: $180
- 3× Drone modules: $150
- 2× Robot modules: $130
- 3× Edge hubs: $150
- 10× Sensor nodes: $900
- **Total: $1,415**

**Monthly:** $0 (no subscriptions)

### OneMind Emergency Kit

**Hardware:**
- USB drive (256GB): $30-75

**Software:**
- All FREE (open source)

**Total:** $30-75 (one-time)

**vs. Bunker AI:** $89 (proprietary, unknown specs)

### Grand Total

**Minimal:** $210 (starter mesh + basic USB)
**Complete:** $1,490 (full mesh + premium USB)
**Monthly:** $0 forever

**ROI:**
- vs. Cellular for 10 sensors: $100-300/month saved
- vs. Bunker AI: $14-59 saved
- vs. Total failure: Priceless

---

## Implementation Timeline

### Overview (14 weeks total)

```
Weeks 1-2:   Meshtastic Testing (OMOS-133)
Weeks 3-4:   MQTT Bridge Setup (OMOS-132)
Weeks 5-8:   Fleet Integration (OMOS-134)
Weeks 9-10:  Emergency Kit Build (OMOS-136)
Weeks 11-12: Testing & Polish
Weeks 13-14: Documentation & Training
```

### Detailed Timeline

**Phase 1: Meshtastic Starter (Weeks 1-2)**
- Week 1:
  - [ ] Order 2× RAK kits ($130)
  - [ ] While waiting: Read all documentation
  - [ ] Plan deployment locations
- Week 2:
  - [ ] Receive hardware
  - [ ] Unbox, assemble, flash firmware
  - [ ] Configure identical settings
  - [ ] Test range across property
  - [ ] Create coverage map

**Phase 2: MQTT Bridge (Weeks 3-4)**
- Week 3:
  - [ ] Deploy MeshtasticBridge service
  - [ ] Connect to NATS on house edge
  - [ ] Test bidirectional messaging
- Week 4:
  - [ ] Implement message persistence
  - [ ] Add monitoring (Prometheus)
  - [ ] Test failover scenarios
  - [ ] Document procedures

**Phase 3: Fleet Integration (Weeks 5-8)**
- Week 5:
  - [ ] Order production hardware ($1,235)
- Week 6:
  - [ ] Add LoRa to Zeus drone (when arrives)
  - [ ] Implement TelemetryBackup class
  - [ ] Test automatic failover
- Week 7:
  - [ ] Add LoRa to Unitree robots
  - [ ] Test status reporting
  - [ ] Test command execution via mesh
- Week 8:
  - [ ] Deploy 10 sensor nodes
  - [ ] Configure solar power
  - [ ] Test full mesh resilience

**Phase 4: Emergency Kit (Weeks 9-10)**
- Week 9:
  - [ ] Purchase USB drive ($30-75)
  - [ ] Install Ollama portable
  - [ ] Download 3 AI models
  - [ ] Build knowledge base
- Week 10:
  - [ ] Add essential tools
  - [ ] Create web UI
  - [ ] Write startup scripts
  - [ ] Test cross-platform

**Phase 5: Testing (Weeks 11-12)**
- Week 11:
  - [ ] Simulate internet outage (24h test)
  - [ ] Verify fleet operates via mesh
  - [ ] Test emergency USB kit
- Week 12:
  - [ ] Performance benchmarks
  - [ ] Security audit
  - [ ] User acceptance testing

**Phase 6: Finalization (Weeks 13-14)**
- Week 13:
  - [ ] Complete all documentation
  - [ ] Create training materials
  - [ ] Record demo videos
- Week 14:
  - [ ] Team training
  - [ ] Handoff procedures
  - [ ] Celebrate! 🎉

### Optional: Bunker AI Research (Parallel)

- Week 1: Order Bunker AI ($89)
- Weeks 2-4: Shipping
- Week 5: Receive and analyze
- Week 6: Document findings
- Week 7: Compare to OneMind Kit
- Week 8: Incorporate insights

---

## Quick Reference

### When You Review This Later

**Start Here:**
1. Read this document (OFFLINE-SERVICES-OVERVIEW.md)
2. Understand the two systems:
   - Meshtastic = Fleet comms
   - Emergency Kit = Personal knowledge
3. Review specific docs as needed

### Key Questions Answered

**Q: What is Meshtastic?**
A: LoRa mesh network for fleet communication when internet fails. See [MESHTASTIC-MESH-INTEGRATION.md](MESHTASTIC-MESH-INTEGRATION.md)

**Q: What is Emergency Kit?**
A: Portable AI + knowledge on USB drive for human use. See [ONEMIND-EMERGENCY-KIT.md](ONEMIND-EMERGENCY-KIT.md)

**Q: Are they the same thing?**
A: No! Different systems, different purposes, work together.

**Q: Which one do I need first?**
A: Meshtastic (fleet is more critical). Emergency Kit can wait.

**Q: What's Bunker AI?**
A: Competitor product ($89). We're buying to analyze and learn from.

**Q: Where are the Linear issues?**
A: Run `python3 scripts/linear/create_meshtastic_issues.py` and `create_emergency_kit_issues.py`

**Q: How much will this cost?**
A: $210 minimum (starter mesh + basic USB), $1,490 complete

**Q: How long to implement?**
A: 14 weeks for everything (can do in phases)

### Visual Cheat Sheet

```
╔═══════════════════════════════════════════════════════════╗
║         OFFLINE SERVICES - QUICK REFERENCE                ║
╚═══════════════════════════════════════════════════════════╝

📡 MESHTASTIC (FOR THE FLEET)
├─ Purpose: Drones/robots communicate when internet down
├─ Form: LoRa radios across property
├─ Cost: $180-1,415
├─ Issues: OMOS-132, 133, 134
└─ Docs: MESHTASTIC-*.md (3 files)

💾 EMERGENCY KIT (FOR YOU)
├─ Purpose: Access knowledge/AI when everything down
├─ Form: Portable USB drive
├─ Cost: $30-75
├─ Issues: OMOS-135, 136
└─ Docs: ONEMIND-EMERGENCY-KIT.md

🎯 TOGETHER: Complete 4-layer resilience stack
   Layer 1: Cloud NATS (primary)
   Layer 2: Meshtastic (fleet backup)
   Layer 3: Edge Ollama (local)
   Layer 4: Emergency Kit (personal backup)

📊 TOTAL INVESTMENT
   Hardware: $210-1,490 one-time
   Software: $0 (all open source)
   Monthly: $0 forever
   Timeline: 14 weeks
```

---

## Next Steps

### Immediate Actions

1. **Review this document** - Understand both systems
2. **Read detailed docs** - Dive into specifics
3. **Review Linear issues** - See implementation tasks
4. **Decide on budget** - Starter vs full deployment
5. **Set timeline** - When to start implementation

### To Start Implementation

**Meshtastic First:**
```bash
# Order hardware
1. Buy 2× RAK19007 kits ($130)
2. While waiting, read MESHTASTIC-STARTER-GUIDE.md

# Create Linear issues
python3 scripts/linear/create_meshtastic_issues.py

# Start Phase 1 (when hardware arrives)
Follow MESHTASTIC-STARTER-GUIDE.md
```

**Emergency Kit Second:**
```bash
# Order hardware
1. Buy USB drive ($30-75)

# Create Linear issues
python3 scripts/linear/create_emergency_kit_issues.py

# Start build
Follow ONEMIND-EMERGENCY-KIT.md
```

**Optional Bunker AI:**
```bash
# Order for research
1. Buy from https://bunkerai.io/ ($89)
2. Wait 2-4 weeks for shipping
3. Analyze and document findings
```

---

## Support & Questions

**Documentation Issues:**
- All docs in `docs/planning/`
- This overview: `OFFLINE-SERVICES-OVERVIEW.md`

**Implementation Questions:**
- Check detailed docs first
- Linear issues have step-by-step tasks
- Code examples in MESHTASTIC-MESH-INTEGRATION.md

**Hardware Questions:**
- Meshtastic: See MESHTASTIC-STARTER-GUIDE.md
- Emergency Kit: See ONEMIND-EMERGENCY-KIT.md

**Cost Questions:**
- See cost sections in this doc
- Detailed breakdown in each system's docs

---

## Document History

| Date | Change |
|------|--------|
| 2026-01-23 | Created master overview document |
| 2026-01-22 | Created all sub-documents |

---

**This Document:** Master reference for offline services
**Purpose:** Quick understanding when reviewing later
**Status:** Complete - ready for review
**Action:** Read this first, then dive into specific docs as needed

🚀 **Complete offline resilience for OneMind OS!**
