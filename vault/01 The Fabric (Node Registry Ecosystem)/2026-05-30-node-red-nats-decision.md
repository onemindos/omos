---
title: "Node-RED + NATS Decision"
date: 2026-05-30
type: decision
status: active
tags: [node-red, n8n, nats, architecture, iot, robotics]
---

# Node-RED + NATS Architecture Decision

## The Question
**"If we have NATS for real-time communication, do we still need Node-RED?"**

## The Answer
**YES.** NATS and Node-RED solve different problems.

```
NATS = Message transport (the highway)
Node-RED = Protocol translator (the on-ramps)

Without Node-RED, your hardware can't speak NATS.
```

---

## What Each Tool Does

### NATS (Message Backbone)
- **Purpose:** Universal message transport
- **Speaks:** NATS protocol only
- **Role:** Connects all services, agents, and data stores
- **Cannot do:** Read ROS2, MAVLink, GPIO, Serial, Modbus

### Node-RED (Hardware Bridge)
- **Purpose:** Translate hardware protocols → NATS messages
- **Speaks:** ROS2, MAVLink, GPIO, Serial, Modbus, MQTT, HTTP
- **Role:** Edge gateway between physical hardware and NATS
- **Cannot do:** SaaS integrations, complex business logic

### n8n (SaaS Orchestrator)
- **Purpose:** Workflow automation for APIs/webhooks
- **Speaks:** REST, WebSocket, 500+ SaaS integrations
- **Role:** Central orchestrator for cloud services
- **Cannot do:** Read ROS2, GPIO, or low-level hardware protocols

---

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│                    OneMind OS Stack                       │
├──────────────────────────────────────────────────────────┤
│                                                            │
│  SaaS Layer (cloud-facing)                                │
│  ══════════════════════════════                           │
│  n8n                                                       │
│    • Notion, Linear, GitHub                               │
│    • Webhooks, notifications                              │
│    • FarmBot REST API                                     │
│                                                            │
│  Message Backbone (universal transport)                   │
│  ══════════════════════════════════════                   │
│  NATS                                                      │
│    • cmd.* (commands)                                     │
│    • events.* (telemetry)                                 │
│    • webhooks.* (inbound)                                 │
│                                                            │
│  Hardware Layer (edge-facing)                             │
│  ══════════════════════════════                           │
│  Node-RED (per site)                                      │
│    • ROS2 (Unitree robots)                                │
│    • MAVLink (drones)                                     │
│    • GPIO (sensors)                                       │
│                                                            │
└──────────────────────────────────────────────────────────┘
```

---

## Example: Unitree G1 Robot

**Without Node-RED (broken):**
```
Unitree G1 (publishes ROS2 topics)
    ↓
NATS ← CANNOT READ ROS2
    ↓
Hermes ← NO DATA
```

**With Node-RED (working):**
```
Unitree G1 (publishes ROS2 topics: /joint_states, /camera, /lidar)
    ↓
Node-RED (node-red-contrib-ros2)
    ├── Reads ROS2 topics
    ├── Transforms to JSON
    └── Publishes to NATS
        ↓
NATS (events.robot.g1.joint_states)
    ↓
Hermes subscribes + ClickHouse stores
```

---

## Deployment Strategy

### Current (n8n only)
```
Mac/VPS
  ├── n8n (SaaS/webhooks)
  └── NATS (message backbone)
```

**Status:** Correct for today — zero hardware deployed yet.

### Future (n8n + Node-RED edge)
```
Mac/VPS (Central)
  ├── n8n (SaaS orchestrator)
  └── NATS (backbone)

House (Raspberry Pi)
  └── Node-RED → GPIO sensors → NATS

Barn (Edge device)
  └── Node-RED → Unitree Go2 ROS2 → NATS

Workshop (Edge device)
  └── Node-RED → Unitree G1 ROS2 → NATS

Aerial (Ground station)
  └── Node-RED → MAVLink parser → NATS
```

**Deploy Node-RED when:**
1. Unitree G1/Go2 arrives (ROS2 required)
2. RECON drone operational (MAVLink required)
3. Edge sensors deployed (GPIO required)

---

## Decision

| Tool | Layer | Deployment | Status |
|------|-------|------------|--------|
| **n8n** | SaaS/webhooks | Mac/VPS (central) | ✅ Active now |
| **NATS** | Message backbone | Mac/VPS (central) | ✅ Active now |
| **Node-RED** | Hardware bridge | Edge devices (per site) | 🔜 Deploy with robots |

---

## What Got Done Today

✅ Installed `@zetanova/n8n-nodes-nats@1.1.9-rc1` in n8n  
✅ Updated `node-red.md` with:
  - Protocol breakdown (ROS2, MAVLink, GPIO, etc.)
  - When to use Node-RED scenarios
  - Architecture decision diagram
  - 6 deployment types (npm, Docker, Raspberry Pi, HA, Cloud, K8s)
  - Recommended deployment pattern for OneMind OS

**Next step:** Restart n8n to load NATS nodes

---

## Key Insight

**NATS doesn't eliminate Node-RED.**  
NATS is the *transport* that Node-RED publishes to.

Think of it like:
- **NATS** = The highway system
- **Node-RED** = The on-ramps (translates hardware → highway-compatible messages)
- **n8n** = The off-ramps (translates highway messages → SaaS APIs)

All three are complementary, not redundant.
