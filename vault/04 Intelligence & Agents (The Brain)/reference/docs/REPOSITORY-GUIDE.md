---
id: ui-02-20260111053925
title: "🧠 Legacy AI Hybrid - Repository Guide"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "🧠 Legacy AI Hybrid - Repository Guide"
---

# 🧠 Legacy AI Hybrid - Repository Guide

**Complete Map of What's Where and Why**

---

## 🎯 Quick Navigation

### 📍 For **Now** (Phase 2A - Production)
**What**: Running intelligent agent system  
**Where**: `legacy-ai-stacks/cognitive-architecture/`  
**Read**: `legacy-ai-stacks/cognitive-architecture/README.md`  
**Deploy**: `legacy-ai-stacks/cognitive-architecture/scripts/deploy.sh`

### 📍 For **Soon** (Phase 2B - Offline/Edge)
**What**: Same system, runs offline on Jetson Nano  
**Where**: `legacy-ai-stacks/cognitive-architecture/edge-deployment/`  
**Plan**: `legacy-ai-stacks/cognitive-architecture/FUTURE-ARCHITECTURE.md`

### 📍 For **Future** (Phase 3 - Multi-Agent)
**What**: Multiple collaborating AI agents  
**Where**: `legacy-ai-stacks/cognitive-architecture/multi-agent/`  
**Plan**: `legacy-ai-stacks/cognitive-architecture/FUTURE-ARCHITECTURE.md`

### 📍 For **Vision** (Robotics)
**What**: Physical embodiment with robotic capabilities  
**Where**: `legacy-ai-stacks/cognitive-architecture/robotics/`  
**Plan**: `legacy-ai-stacks/cognitive-architecture/FUTURE-ARCHITECTURE.md`

---

## 📂 Directory Structure Overview

```
Legacy-ai-Hybrid/
│
├── README.md                           ← Start here for overview
├── ARCHITECTURE-COMPLETE-MAP.md        ← Visual map of all phases
├── REPOSITORY-GUIDE.md                 ← This file
│
├── docs/                               📚 Documentation
│   ├── architecture/
│   │   └── DEVELOPMENT-MASTER.md      (Master reference)
│   ├── guides/
│   ├── current-state/
│   └── archive/reorganization-2025-10-24/
│
├── legacy-ai-stacks/                   🧠 THE INTELLIGENT BEING
│   ├── README.md
│   │
│   └── cognitive-architecture/         ⭐ MAIN PROJECT DIR
│       │
│       ├── ✅ PHASE 2A (Now)
│       ├── 🚀 PHASE 2B (Offline/Edge)
│       ├── 🤖 PHASE 3 (Multi-Agent)
│       └── 🦾 ROBOTICS
│
├── scripts/                            ⚙️ Utilities
│   ├── deployment/
│   ├── maintenance/
│   └── monitoring/
│
└── mcp/                                🔗 Tool integrations
```

---

## 🔍 What's In Each Major Directory?

### `legacy-ai-stacks/cognitive-architecture/` (THE CORE)

#### ✅ Phase 2A (Production Now)
```
brain/              🧠 Letta reasoning engine
senses/             👁️👂�� LiveKit perception (voice/vision)
tools/              🔧 n8n, Home Assistant, Node-RED, EMQX
memory/             💾 Supabase, PostgreSQL, Redis
agent/              🤖 Orchestration layer
scripts/            ⚙️ Deploy, backup, monitor
```

**Status**: ✅ Running in production  
**Latency**: <100ms voice, <1000ms vision  
**Internet**: Required  
**Hardware**: Cloud-based (AWS VM)

---

#### 🚀 Phase 2B (Offline/Edge - Coming Soon)
```
edge-deployment/
├── jetson-nano/
│   ├── models/      TensorRT quantized models
│   │   ├── whisper.cpp (STT <80ms)
│   │   ├── piper.onnx (TTS <50ms)
│   │   ├── yolov8-nano (Vision 60fps)
│   │   └── llama-3.2-1b (LLM <60ms)
│   └── local-services/
│       └── docker-compose-edge.yml
│
├── docker-config/   Edge optimizations
└── sync-strategies/ PowerSync, CRDT, fallback
```

**Status**: 🚀 Architecture designed, ready for development  
**Timeline**: 2-3 months of development  
**Hardware**: Jetson Nano 4GB ($99)  
**Cost**: $0/month (no cloud APIs needed)

**Key Difference from Phase 2A**:
- Runs completely offline
- SQLite instead of Supabase
- Local models instead of cloud APIs
- Jetson Nano instead of VM
- Sync protocols for cloud fallback

---

#### 🤖 Phase 3 (Multi-Agent Orchestration)
```
multi-agent/
├── orchestration/
│   ├── coordinator/   Routes tasks between agents
│   └── supervisor/    Monitors health & authority
│
├── agents/
│   ├── specialized-agents/   Research, planning, execution, learning
│   └── domain-specific/      IoT, enterprise, medical, finance
│
├── communication/     Inter-agent messaging
└── shared-resources/  Shared brain, memory, tools, senses
```

**Status**: 🏗️ Architecture designed, framework selection pending  
**Timeline**: 3-4 months after Phase 2B complete  
**Agents**: 3-5 specialized agents collaborating  
**Benefit**: Emergent capabilities, better decision-making

**Example Use Case**:
```
Research Agent: "Here's what I found..."
     ↓
Planning Agent: "Here's the strategy..."
     ↓
Execution Agent: "Here's what I'll do..."
     ↓
Learning Agent: "Here's what we learned..."
     ↓
Coordinator: Validates, logs, reports
```

---

#### 🦾 Robotics (Physical Embodiment)
```
robotics/
├── motion-control/
│   ├── arm-control/     Kinematics, dynamics, joints
│   └── locomotion/      Wheels, legs, humanoid gait
│
├── perception/
│   ├── computer-vision/ Object detection, segmentation
│   ├── lidar-processing/ 3D point clouds, mapping
│   └── sensor-fusion/   IMU, cameras, depth sensors
│
├── manipulation/
│   ├── object-grasping/ Grasp planning, force control
│   └── trajectory-planning/ Path smoothing, optimization
│
├── navigation/
│   ├── path-planning/   RRT*, Dijkstra, dynamic planning
│   └── localization/    SLAM, GPS, dead reckoning
│
└── hardware-abstraction/ ROS 2 adapters, drivers
```

**Status**: 🏗️ Architecture designed, hardware selection pending  
**Timeline**: Parallel with multi-agent, 4-6 months  
**Hardware**: Robot selection (Boston Dynamics, Unitree, Custom)  
**Integration**: Via brain's tool-calling mechanism

**Capabilities**:
- Real-time perception (60fps)
- Smooth motion control (<100ms)
- Accurate manipulation (±1cm)
- Autonomous navigation with SLAM

---

#### 💡 Hybrid Offline (Local Model Inference)
```
hybrid-offline/
├── local-models/
│   ├── llm/         Llama 3.2 1B, Mistral 7B
│   ├── vision/      YOLOv8-nano, MobileNet
│   ├── audio/       Whisper.cpp, wav2vec
│   └── tts/         Piper, Glow-TTS
│
├── sync-protocols/  PowerSync, CRDT, replication
└── fallback-strategies/ Graceful degradation
```

**Status**: 🏗️ Directory structure ready for implementation  
**Purpose**: Enable seamless online→offline→online transitions  
**Key Feature**: Zero data loss during sync

---

## 🎯 The Intelligent Being Model

### What is it?

A **complete autonomous system** with:

```
PERCEPTION      COGNITION       ACTION          MEMORY
(Senses)    →   (Brain)     →   (Tools)     →   (Storage)
├─ Voice        Letta           n8n             Supabase
├─ Vision       Reasoning       n8n             PostgreSQL
└─ Audio        Tool Calling    Home Asst       Redis
                Decisions       Node-RED
                Learning        EMQX
                                Robotics

                        ↓
                   EVOLUTION
        Offline (2B) + Multi-Agent (3) + Robotics
```

### How does it learn?

1. Perceives through senses (vision, voice, sensors)
2. Thinks through brain (reasoning, decisions)
3. Acts through tools (execute in world)
4. Stores in memory (learns from experience)
5. Improves over time (feedback loop)

### How does it scale?

- **Phase 2A**: Single agent, cloud-based
- **Phase 2B**: Same agent, runs offline anywhere
- **Phase 3**: Multiple agents, collaborative
- **Robotics**: Physical embodiment

Each phase **extends** capabilities without breaking previous ones.

---

## 📊 Status Dashboard

### Phase 2A ✅ Production
```
✅ Brain (Letta)                  Working, 25 containers running
✅ Senses (LiveKit + APIs)        Working, <100ms latency
✅ Tools (n8n, etc)         Working, processing requests
✅ Memory (Supabase)              Working, storing conversations
✅ Orchestration (Agent)          Working, coordinating all components
✅ Deployment                     Working, stable on AWS VM
✅ Monitoring                     Working, logs & alerts configured
```

### Phase 2B 🚀 Ready for Development
```
🏗️ Offline models               Architecture complete, ready to add Whisper/Piper/Llama
🏗️ Edge deployment              Directories ready, need Jetson setup & docker configs
🏗️ Sync protocols               Directory ready, need PowerSync adapter implementation
🏗️ Fallback strategies          Directory ready, need graceful degradation logic
```

### Phase 3 🤖 Ready for Design
```
🏗️ Multi-agent framework        Directories ready, need agent library selection
🏗️ Orchestration               Coordinator/supervisor structure ready
🏗️ Communication               Directory ready, need message protocol implementation
```

### Robotics 🦾 Ready for Hardware Selection
```
🏗️ Motion control              Structure ready, need hardware selection
🏗️ Perception                  Structure ready, need sensor integration
🏗️ Navigation                  Structure ready, need ROS 2 adapter
🏗️ Hardware abstraction        Directory ready, need driver development
```

---

## 🚀 How to Get Started

### If you're working on **Phase 2A** (Now):
```bash
cd legacy-ai-stacks/cognitive-architecture
ls -la brain/ senses/ tools/ memory/ agent/ scripts/
cat README.md
./scripts/deploy.sh
```

### If you're planning **Phase 2B** (Offline):
```bash
cd legacy-ai-stacks/cognitive-architecture/edge-deployment
ls -la jetson-nano/ hybrid-offline/
cat ../FUTURE-ARCHITECTURE.md | grep "Phase 2B" -A 50
```

### If you're designing **Phase 3** (Multi-Agent):
```bash
cd legacy-ai-stacks/cognitive-architecture/multi-agent
ls -la orchestration/ agents/ communication/ shared-resources/
cat ../FUTURE-ARCHITECTURE.md | grep "Phase 3" -A 50
```

### If you're planning **Robotics**:
```bash
cd legacy-ai-stacks/cognitive-architecture/robotics
ls -la motion-control/ perception/ manipulation/ navigation/
cat ../FUTURE-ARCHITECTURE.md | grep "Robotics" -A 50
```

---

## 📚 Documentation Map

| File | Purpose | Audience |
|------|---------|----------|
| `README.md` (root) | Project overview | Everyone |
| `ARCHITECTURE-COMPLETE-MAP.md` | Visual architecture | Visual learners |
| `REPOSITORY-GUIDE.md` | This file - navigation | New contributors |
| `legacy-ai-stacks/README.md` | Stack organization | DevOps/Deployment |
| `legacy-ai-stacks/cognitive-architecture/README.md` | Current system | Active developers |
| `legacy-ai-stacks/cognitive-architecture/FUTURE-ARCHITECTURE.md` | Future phases | Planners/Architects |
| `docs/architecture/DEVELOPMENT-MASTER.md` | Complete reference | Everyone eventually |

---

## 💡 Key Decisions

### Why organize by cognitive function?
- **Clarity**: Structure matches what it actually is (an intelligent being)
- **Scalability**: Easy to add new senses, tools, or memory
- **Maintainability**: Clear purpose for each component
- **Extensibility**: Phases 2B, 3, and Robotics fit naturally

### Why not change existing Phase 2A?
- ✅ Production system unaffected
- ✅ No breaking changes
- ✅ Explicit intention for future work
- ✅ Reduced risk and disruption

### Why these phases?
- **2A** (Now): Prove concept, production hardening
- **2B** (Next): Extend reach (edge computing)
- **3** (Then): Deepen capability (collaboration)
- **Robotics** (Parallel): Add embodiment (physical form)

---

## 🎯 Success Metrics

### Phase 2A
✅ <100ms voice latency  
✅ 99.9% uptime  
✅ Production stable

### Phase 2B
✅ Runs on $99 hardware  
✅ Works completely offline  
✅ Zero data loss on sync  
✅ <200ms voice latency offline

### Phase 3
✅ 3+ agents collaborating  
✅ Emergent capabilities  
✅ Conflict-free decisions  
✅ Better outcomes than single agent

### Robotics
✅ Real-time perception (60fps)  
✅ Autonomous navigation  
✅ Accurate manipulation  
✅ Integrated with multi-agent system

---

## 🤝 Contributing

### Adding to Phase 2A?
→ Modify in `legacy-ai-stacks/cognitive-architecture/{component}/`

### Building Phase 2B?
→ Add to `legacy-ai-stacks/cognitive-architecture/edge-deployment/`

### Designing Phase 3?
→ Add to `legacy-ai-stacks/cognitive-architecture/multi-agent/`

### Planning Robotics?
→ Add to `legacy-ai-stacks/cognitive-architecture/robotics/`

**All files in their respective directories. Don't mix phases.**

---

## ❓ FAQ

**Q: Where do I start?**  
A: `legacy-ai-stacks/cognitive-architecture/README.md` for Phase 2A overview

**Q: What's the current system doing?**  
A: Running on AWS VM, handling voice/vision conversations, storing in cloud

**Q: When does Phase 2B start?**  
A: After Phase 2A stabilizes (~6-8 weeks from now)

**Q: Why create empty directories now?**  
A: To establish clear intended structure and prevent ad-hoc organization later

**Q: Can I work on multiple phases?**  
A: Yes! They're isolated directories, work in parallel if desired

**Q: How do I know which file to edit?**  
A: Look at the phase you're working on, then the component, then the specific area

---

## 📞 Quick Commands

```bash
# See current Phase 2A structure
tree -L 2 legacy-ai-stacks/cognitive-architecture/

# See Phase 2B structure (empty, ready for files)
tree legacy-ai-stacks/cognitive-architecture/edge-deployment/

# See Phase 3 structure (empty, ready for files)
tree legacy-ai-stacks/cognitive-architecture/multi-agent/

# See Robotics structure (empty, ready for files)
tree legacy-ai-stacks/cognitive-architecture/robotics/

# Deploy Phase 2A
cd legacy-ai-stacks/cognitive-architecture && ./scripts/deploy.sh

# Check git status
git status

# View recent changes
git log --oneline -10
```

---

## �� Learning Path

### Week 1: Understand Current System
1. Read: `legacy-ai-stacks/cognitive-architecture/README.md`
2. Explore: Current Phase 2A components
3. Deploy: Try `./scripts/deploy.sh`

### Week 2: Understand Future Plans
1. Read: `ARCHITECTURE-COMPLETE-MAP.md`
2. Read: `legacy-ai-stacks/cognitive-architecture/FUTURE-ARCHITECTURE.md`
3. Explore: Future directory structures

### Week 3+: Start Contributing
1. Choose phase (2A, 2B, 3, or Robotics)
2. Pick component to work on
3. Add to appropriate directory
4. Commit with clear message

---

**Repository**: https://github.com/Legacyx-ai/Legacy-AI-Hybrid  
**Status**: ✅ Phase 2A Production | 🚀 Phase 2B Ready | 🤖 Phase 3 Designed | 🦾 Robotics Architected  
**Updated**: October 24, 2025  

🧠 Let's build an intelligent being together!


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
