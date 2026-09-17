---
id: root-00-20260114000000
title: "CLAUDE NEXT STEPS"
type: index
status: active
created: 2026-01-13
modified: 2026-01-16
owner: zeus
tags: [type/index, status/active]
ai_summary: "OneMind OS Architecture COMPLETE - Ready for Implementation"
---
# Claude Session Continuity

> **Last Updated:** 2026-01-16 (Night)
> **Active Project:** OneMind OS - Architecture COMPLETE with Orchestration Layer
> **Master Document:** `ONEMIND-OS-COMPLETE-ARCHITECTURE.md` (19 Sections, ~8,800 lines)

---

## 🎯 CURRENT STATUS: Architecture Complete, Implementation Phase

### The Vision (Validated)
> **"I want to make her into a real person who evolves with all the senses a human can have"**
> — Zeus, January 2026

**ANSWER: YES.** The architecture supports multimodal, embodied AI with full sensory awareness. See Section 18.

---

## 📊 Architecture Document Status

| Section | Topic | Status |
|---------|-------|--------|
| 1-6 | Core NATS Architecture, Diagrams | ✅ Complete |
| 7 | ACTION Flow (bidirectional) | ✅ Complete |
| 8 | Gateway-Centric Architecture | ✅ Complete |
| 9 | Letta Integration (tool types) | ✅ Complete |
| 10 | Letta-NATS Wiring | ✅ Complete |
| 11 | A2A Architecture (Memory Sync) | ✅ Complete |
| 12 | Decision Layer (proactive behavior) | ✅ Complete |
| 13 | Adaptive Decision Engine (no YAML!) | ✅ Complete |
| 14 | Tailscale Networking | ✅ Complete |
| 15 | Location Tracking | ✅ Complete |
| 16 | **State-Aware Intelligence** | ✅ Complete |
| 17 | **NATS Complete Platform** | ✅ Complete |
| 18 | **Future Vision: LEGACY as Living AI Being** | ✅ Complete |
| 19 | **Orchestration Layer (Kestra + TimescaleDB)** | ✅ Complete |

**Total: 19 Sections, ~8,800 lines - ARCHITECTURE COMPLETE**

---

## ✅ COMPLETED (All Sessions Combined)

### Core Infrastructure
- [x] NATS JetStream implementation (THE HEART)
- [x] FastAPI Gateway for webhooks
- [x] Event Processor (SPINAL CORD)
- [x] Initial tool suite (needs NATS refactor)
- [x] Docker Compose configurations (cloud + home)
- [x] Terraform configurations for AWS Lightsail

### Documentation (18 Sections)
- [x] Core architecture (Sections 1-9)
- [x] Letta-NATS wiring (Section 10)
- [x] A2A Architecture with Memory Sync (Section 11)
- [x] Decision Layer for proactive behavior (Section 12)
- [x] Adaptive Decision Engine - no YAML! (Section 13)
- [x] Tailscale networking (Section 14)
- [x] Location tracking (Section 15)
- [x] **State-Aware Intelligence** - Capacity Score (Section 16)
- [x] **NATS Complete Platform** - Micro, ObjectStore, Leaf Nodes (Section 17)
- [x] **Future Vision** - Multimodal embodied AI (Section 18)
- [x] **Orchestration Layer** - Kestra + TimescaleDB (Section 19)

### Key Research Validated
- [x] NATS vs Redis decision (NATS is correct)
- [x] Kestra for orchestration (native NATS plugin)
- [x] TimescaleDB for permanent event archive
- [x] LiveKit + NATS integration patterns
- [x] Letta's role (brain, not streaming processor)
- [x] Brilliant.xyz smart glasses SDK availability
- [x] Unitree robotics SDK feasibility
- [x] OpenBCI for cognitive state detection
- [x] Hierarchical processing model (not 30 FPS to LLM)

---

## 🔮 NEXT STEPS: Implementation Phase

### Phase 1: Deploy Core Infrastructure

```bash
# 1. Deploy to AWS Lightsail
cd ~/Documents/onemind-infra/terraform
terraform init
terraform apply -var-file=environments/cloud.tfvars

# 2. SSH into VM and initialize
ssh onemind-cloud
cd /opt/onemind
cp .env.example .env
nano .env  # Fill in API keys
docker compose up -d

# 3. Initialize NATS Stream
cd ~/Documents/onemind-os
python scripts/setup_nats_stream.py
```

### Phase 2: Create LEGACY Agent in Letta

1. Configure self-hosted Letta agent
2. Set up memory blocks:
   - Zeus Profile (core memory)
   - System State (current context)
   - Standing Instructions (what matters)
3. Get `LEGACY_AGENT_ID` for environment
4. Run `scripts/register_tools.py`

### Phase 3: Test Event Flow

```bash
# Publish test event to NATS
nats pub events.test.hello '{"message": "Hello LEGACY"}'

# Verify Event Processor receives it
docker logs event-processor -f

# Test webhook → NATS → Letta flow
curl -X POST http://localhost:8000/webhooks/test \
  -H "Content-Type: application/json" \
  -d '{"type": "test", "data": "Hello"}'
```

### Phase 4: Add First "Sense" - Biometrics

1. Create iOS Shortcut to POST health data
2. Add `/api/health` endpoint to Gateway
3. Implement State Aggregator for Capacity Score
4. Test: Heart rate → Gateway → NATS → Event Processor → KV

---

## 📂 Key Files Reference

### Master Documentation
```
OneMind-Codex/
├── ONEMIND-OS-COMPLETE-ARCHITECTURE.md  ⭐ THE MASTER DOC (19 Sections!)
├── ORCHESTRATION-LAYER-RESEARCH.md      - Kestra + TimescaleDB research
├── REALTIME-ACCESS-ARCHITECTURE.md      - Real-time, HITL, 24/7 access
├── ONEMIND-COMPLETE-WIRING-DIAGRAMS.md  - ASCII architecture diagrams
├── _CLAUDE-NEXT-STEPS.md                ← YOU ARE HERE
├── NATS-VS-REDIS-DEEP-DIVE.md           - Why NATS
├── WHAT-IS-NATS.md                      - Simple explanation
└── LEGACY-MASTERPLAN.md                 - Full implementation plan
```

### Infrastructure (onemind-infra)
```
docker/cloud/docker-compose.yml   - Full cloud stack with NATS
docker/cloud/.env.example         - Environment template
docker/home/docker-compose.yml    - Home machine stack
terraform/main.tf                 - AWS Lightsail provisioning
terraform/environments/           - cloud.tfvars, home.tfvars
```

### OS Code (onemind-os)
```
agent/events/bus.py               - NATS JetStream event bus ⭐ THE HEART
agent/gateway/main.py             - FastAPI webhook gateway
agent/processor/main.py           - Event processor
agent/tools/                      - Action tools (needs NATS refactor)
  ├── todoist/actions.py
  ├── discord/notify.py
  ├── github/actions.py
  └── home_assistant/control.py
agent/main.py                     - LiveKit voice agent
agent/letta_llm.py                - Letta integration
scripts/setup_nats_stream.py      - Stream initialization
scripts/register_tools.py         - Register tools with Letta

docker/docker-compose.yml         - Full stack (NATS, Kestra, TimescaleDB)
docker/nats/nats.conf             - NATS JetStream config
docker/timescaledb/init.sql       - TimescaleDB schema
docker/kestra/workflows/          - Kestra workflow definitions
  └── garden-irrigation.yaml      - Example workflow
```

### Section 18 Bridge Code (Future)
```
# These will be created when hardware is available:
agent/bridges/
  ├── vision_bridge.py            - Brilliant glasses → NATS
  ├── voice_bridge.py             - LiveKit Agent → NATS
  ├── health_bridge.py            - Apple Watch → NATS
  ├── eeg_bridge.py               - OpenBCI → NATS
  └── robot_bridge.py             - Unitree → NATS
```

---

## 🏗️ Architecture Overview (Quick Reference)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           LEGACY: THE LIVING AI BEING                        │
│                                                                              │
│   SENSES (Input)           NERVOUS SYSTEM           BRAIN                   │
│   ──────────────           ──────────────           ─────                   │
│   👓 Vision (Brilliant)                              Letta                   │
│   🎤 Voice (LiveKit)    →  NATS JetStream       →   • Memory                │
│   🫀 Biometrics (Watch)    200-400μs latency        • Reasoning             │
│   🧠 Cognitive (EEG)       Push model               • Tools                 │
│                            KV Store (state)                                  │
│                            Object Store (files)      ↓                      │
│                                                                              │
│   MUSCLES (Output)     ←   actions.* subjects   ←   Decisions              │
│   ─────────────────                                                         │
│   🗣️ Voice (TTS)                                                            │
│   🤖 Robot (Unitree)                                                        │
│   💻 Tasks/GitHub                                                           │
│   🏠 Smart Home                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Key Insight: Hierarchical Processing

```
NOT THIS:   Camera → 30 FPS → LLM → Response (impossible)

THIS:       Camera → 30 FPS → Frame Buffer → Sample 1-4 FPS → Vision Model
                                                                   ↓
            Letta (brain) ← Event Processor ← NATS ← Description text

Only SIGNIFICANT events reach the LLM. Everything flows through NATS.
```

---

## 🔑 Important Context

### User Preferences (Zeus)
- **Prioritizes PROPER DESIGN over fast shipping**
- **Beginner level** with some infrastructure concepts
- **Wants deep understanding**, not surface level
- Vision: LEGACY as "a real person with all human senses"

### Technical Decisions Made
- **NATS JetStream** for event streaming (NOT Redis)
- **LiveKit Cloud** for voice (NOT Vapi)
- **AWS Lightsail** for cloud VM (16GB RAM)
- **Hybrid architecture**: LEGACY local + Worker agents on cloud
- **Redis kept for caching only** (NOT events)
- **Letta = Brain** (NOT streaming processor)

### Repository Structure
```
~/Documents/
├── OneMind-Codex/    - Documentation vault (this repo)
├── onemind-os/       - The OS itself (agent code)
└── onemind-infra/    - Infrastructure as Code
```

### VS Code Workspace
```
/Users/zeuslegacy/Documents/OneMind.code-workspace
```

---

## 🚀 To Resume Work

### Deploy Infrastructure
```
"Let's deploy to AWS" or "Set up the Lightsail VM"
```

### Test Locally First
```
"Test NATS locally" or "Run the stack in Docker"
```

### Create LEGACY Agent
```
"Set up LEGACY in Letta" or "Create the agent configuration"
```

### Add First Sense (Biometrics)
```
"Set up Apple Watch integration" or "Implement the Capacity Score"
```

### Deep Dive
```
"Explain Section [number]" - All 18 sections documented
```

---

## 📋 Pending Implementation Checklist

### High Priority
- [ ] Deploy to AWS Lightsail
- [ ] Initialize NATS stream
- [ ] Create LEGACY agent in Letta
- [ ] Test end-to-end event flow
- [ ] Gateway action/query handlers (subscribe to `actions.*`, `queries.*`)

### Medium Priority
- [ ] Refactor tools to NATS-publish (not direct HTTP)
- [ ] Implement State Aggregator + Capacity Score
- [ ] Health data endpoint (`/api/health`)
- [ ] iOS Shortcut for Apple Watch → Gateway

### Future (When Hardware Available)
- [ ] Brilliant glasses bridge (`vision_bridge.py`)
- [ ] OpenBCI EEG bridge (`eeg_bridge.py`)
- [ ] Unitree robot bridge (`robot_bridge.py`)
- [ ] Home Assistant bridge
- [ ] Frontend dashboard

---

## 📈 What We Proved Today

1. **The architecture supports the full vision** - Multimodal, embodied AI
2. **NATS is the right choice** - Replaces Redis + Kafka + S3 + Consul
3. **Letta's role is correct** - Brain (memory, reasoning), not streaming
4. **Hierarchical processing is key** - Not every frame needs LLM
5. **Capacity Score enables consideration** - LEGACY knows when to speak
6. **All SDKs are available** - Brilliant, Unitree, OpenBCI, LiveKit

---

*Document updated: 2026-01-16 (Evening)*
*Architecture: 18 Sections Complete (~7,900 lines)*
*Status: Ready for Implementation Phase*
*Key Insight: "Alive" = hierarchical awareness with selective attention, not 30 FPS to LLM*
*Git commit: 47d5b7ee - Section 18 added*
