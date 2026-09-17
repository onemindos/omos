---
title: "Node-RED Deployment — Complete"
date: 2026-05-30
type: summary
status: complete
---

# Node-RED Central Orchestrator — Deployed ✅

## What Got Done

### 1. Installation
✅ **Installed Node-RED v4.1.10** (native npm, global)  
✅ **Installed node-red-contrib-nats-suite** (NATS Core + JetStream)  
✅ **Created LaunchAgent** `com.onemind.node-red`  
✅ **Running on :1880**

### 2. Stack Integration
Node-RED now joins your existing OneMind OS stack:

```
Mac (Central)
├── n8n (:5678)          ← SaaS gateway
├── Node-RED (:1880)     ← Hardware orchestrator (NEW)
├── NATS (:4222)         ← Message backbone
├── ClickHouse (:8124)   ← Event store
├── Langfuse (:3004)     ← Observability
└── Hermes Dashboard     ← Agent interface
```

### 3. Management
```bash
# Start
launchctl load ~/Library/LaunchAgents/com.onemind.node-red.plist

# Stop
launchctl unload ~/Library/LaunchAgents/com.onemind.node-red.plist

# Restart
launchctl restart com.onemind.node-red

# Logs
tail -f ~/.node-red/logs/node-red-stdout.log

# Install nodes
cd ~/.node-red && npm install <package-name>
launchctl restart com.onemind.node-red
```

### 4. Access
- **UI:** http://localhost:1880
- **Data:** ~/.node-red/
- **Flows:** Will be stored in ~/.node-red/flows.json
- **Logs:** ~/.node-red/logs/

---

## Why This Setup?

### Native npm + LaunchAgent (Not Docker)
**Reasoning:**
1. ✅ Matches your existing pattern (n8n, NATS, ClickHouse all native)
2. ✅ Direct NATS access (no Docker networking)
3. ✅ Easy community node installation
4. ✅ Native macOS integration
5. ✅ Simple management with `launchctl`

### Central vs Edge Architecture
**Central Node-RED (Mac):**
- Role: Orchestrator
- Subscribes to all edge events (`events.site.#`)
- Publishes commands to edges (`cmd.site.<location>.*`)
- Aggregates fleet telemetry
- Coordinates cross-site automation
- Integrates with n8n for SaaS workflows

**Edge Node-RED (Future):**
- Role: Hardware translator
- Reads local protocols (ROS2, MAVLink, GPIO)
- Publishes to NATS
- Executes local commands from central
- Survives network outages (local autonomy)

---

## Multi-Site Orchestration

### NATS Subject Pattern
```
events.site.house.sensor.temperature     ← House sensor
events.site.barn.robot.go2.telemetry     ← Barn robot
events.site.workshop.robot.g1.telemetry  ← Workshop robot
events.site.aerial.drone.telemetry       ← Aerial drone

cmd.site.house.ac.set_temp               ← Command to house
cmd.site.barn.robot.go2.patrol           ← Command to barn
cmd.site.workshop.robot.g1.pick          ← Command to workshop
cmd.site.aerial.drone.goto_waypoint      ← Command to drone
```

### Orchestration Example
**Scenario:** Barn temperature > 80°F

Central Node-RED:
1. Receives `events.site.barn.climate.temperature` (85°F)
2. Checks threshold (> 80°F)
3. Publishes `cmd.site.house.ac.set_temp` (cool house)
4. Calls n8n webhook (notify Zeus)
5. Logs to ClickHouse (analytics)

Edge Node-RED does NOT make this decision — central does.

---

## Next Actions

### Immediate (Today)
1. **Configure NATS broker in Node-RED UI**
   - Go to http://localhost:1880
   - Add NATS config node
   - Server: `nats://localhost:4222`
   - Name: `OneMind NATS`

2. **Create first flow (test NATS)**
   - NATS In node: `events.test`
   - Debug node
   - Test with terminal: `nats pub events.test "hello"`

3. **Install dashboard**
   ```bash
   cd ~/.node-red && npm install node-red-dashboard
   launchctl restart com.onemind.node-red
   ```
   Access: http://localhost:1880/ui

### Phase 2 (When Robots Arrive)
1. Deploy Node-RED on barn edge device (Docker)
2. Install `node-red-contrib-ros2` on edge
3. Connect edge to central NATS
4. Test Unitree Go2 ROS2 → NATS → Central

### Phase 3 (Full Fleet)
1. Deploy Node-RED on all sites
2. Create fleet monitoring dashboard
3. Build cross-site automation flows
4. Integrate with n8n workflows

---

## Key Decisions Made

### 1. "Do we need Node-RED if we have NATS?"
**YES.** NATS is the highway. Node-RED is the on-ramp for hardware.

NATS can't read ROS2, MAVLink, or GPIO. Node-RED translates hardware protocols into NATS messages.

### 2. "What deployment for central Node-RED?"
**Native npm + LaunchAgent.** Matches your existing stack pattern.

### 3. "How to orchestrate multiple Node-REDs?"
**Central/Edge pattern:**
- Central = Orchestrator (Mac)
- Edge = Translators (per site)
- NATS = Universal transport

Central subscribes to all edge events, makes fleet-wide decisions, publishes commands to edges.

---

## Files Created

1. `/Users/onemind/Library/LaunchAgents/com.onemind.node-red.plist`
   - LaunchAgent config

2. `/Users/onemind/.node-red/start-node-red.sh`
   - Wrapper script (handles nvm PATH)

3. `/Users/onemind/Documents/onemind/00-24 SO (Sovereign Operations)/05 Stack (Arsenal Ecosystem)/iot/automation/node-red-deployment-guide.md`
   - Complete deployment guide
   - Multi-site orchestration patterns
   - Edge deployment instructions
   - NATS subject conventions

---

## Current Stack Status

| Service | Port | Status | LaunchAgent |
|---------|------|--------|-------------|
| NATS | 4222 | ✅ Running | (brew service) |
| n8n | 5678 | ✅ Running | com.onemind.n8n |
| **Node-RED** | **1880** | **✅ Running** | **com.onemind.node-red** |
| ClickHouse | 8124 | ✅ Running | com.onemind.clickhouse |
| Langfuse | 3004 | ✅ Running | (Docker) |
| Hermes WebUI | 9119 | ✅ Running | com.onemind.hermes-webui |
| Notify Router | - | ✅ Running | com.onemind.notify-router |

---

## Resources

- **Node-RED UI:** http://localhost:1880
- **Deployment Guide:** `/00-24 SO (Sovereign Operations)/05 Stack (Arsenal Ecosystem)/iot/automation/node-red-deployment-guide.md`
- **Node-RED Docs:** https://nodered.org/docs/
- **NATS Nodes:** https://flows.nodered.org/node/node-red-contrib-nats-suite

---

**Status:** Ready for configuration. No edge nodes yet (waiting for robot hardware).

Last Updated: May 30, 2026
