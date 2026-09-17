---
title: "Node-RED Central Orchestrator — Deployment Guide"
date: 2026-05-30
type: guide
status: active
tags: [node-red, nats, orchestration, deployment, multi-site]
---

# Node-RED Central Orchestrator — Deployment Guide

## ✅ Installation Complete

**Status:** Node-RED running on Mac as LaunchAgent

### Access
- **UI:** http://localhost:1880
- **Data:** ~/.node-red/
- **Logs:** ~/.node-red/logs/
- **LaunchAgent:** com.onemind.node-red
- **Version:** Node-RED v4.1.10 + Node.js v24.15.0

### Installed Nodes
- ✅ **node-red-contrib-nats-suite** — NATS Core + JetStream + KV + Object Store

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  OneMind OS Central (Mac)                    │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  n8n (:5678)           ← SaaS gateway (webhooks, APIs)       │
│  Node-RED (:1880)      ← Hardware orchestrator (central)     │
│  NATS (:4222)          ← Message backbone                    │
│  ClickHouse (:8124)    ← Event store                         │
│  Langfuse (:3004)      ← Observability                       │
│                                                               │
└─────────────────────────────────────────────────────────────┘
           ↓ NATS (cmd.*, events.*)
┌─────────────────────────────────────────────────────────────┐
│                      Edge Nodes (Future)                     │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  House (Raspberry Pi)    → Node-RED → GPIO sensors          │
│  Barn (Edge device)      → Node-RED → Unitree Go2 (ROS2)    │
│  Workshop (Edge device)  → Node-RED → Unitree G1 (ROS2)     │
│  Aerial (Ground station) → Node-RED → MAVLink parser        │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

---

## Role: Central vs Edge

### Central Node-RED (Mac) — "The Orchestrator"
**Purpose:** Coordinate across all edge nodes, execute multi-site workflows

**Responsibilities:**
- Subscribe to NATS events from all edge nodes
- Aggregate telemetry (robot status, sensor readings, drone position)
- Execute cross-site automation (e.g., "if barn temp > 80°F, notify + activate house AC")
- Publish commands to edge nodes via NATS
- Dashboard/monitoring of entire fleet
- Integration with n8n (SaaS → Hardware workflows)

**Does NOT:**
- Read GPIO directly (delegates to edge)
- Parse ROS2 directly (delegates to edge)
- Parse MAVLink directly (delegates to edge)

### Edge Node-RED (Per Site) — "The Translators"
**Purpose:** Hardware protocol → NATS message translation

**Responsibilities:**
- Read local hardware (GPIO, ROS2, MAVLink, Serial)
- Publish to NATS (`events.site.<location>.*`)
- Subscribe to NATS (`cmd.site.<location>.*`)
- Execute local actions (drive motors, read sensors)
- Survive network outages (local autonomy)

**Does NOT:**
- Make cross-site decisions (central orchestrator does this)
- Store long-term data (NATS → ClickHouse does this)
- SaaS integration (n8n does this)

---

## NATS Subject Conventions for Multi-Site

### Events (Edge → Central)
```
events.site.house.sensor.temperature     ← House sensor data
events.site.house.sensor.motion
events.site.barn.robot.go2.telemetry     ← Unitree Go2 status
events.site.barn.climate.temperature
events.site.workshop.robot.g1.telemetry  ← Unitree G1 status
events.site.aerial.drone.telemetry       ← RECON drone position
events.site.aerial.drone.camera          ← Drone camera feed
```

### Commands (Central → Edge)
```
cmd.site.house.light.turn_on
cmd.site.house.ac.set_temp
cmd.site.barn.robot.go2.patrol           ← Tell Go2 to patrol
cmd.site.barn.robot.go2.return_home
cmd.site.workshop.robot.g1.pick          ← Tell G1 to pick object
cmd.site.workshop.robot.g1.navigate
cmd.site.aerial.drone.takeoff
cmd.site.aerial.drone.goto_waypoint
```

### Status (Heartbeats)
```
status.site.house.node_red               ← Edge health check
status.site.barn.node_red
status.site.workshop.node_red
status.site.aerial.node_red
```

---

## Orchestration Patterns

### Pattern 1: Event Aggregation (Central)
Central Node-RED subscribes to all edge events:

```
[NATS In: events.site.#]  ← Subscribe to ALL sites
    ↓
[Function: Parse + Route]
    ↓
[Switch: By Site]
    ├── House → [Dashboard: House Panel]
    ├── Barn → [Dashboard: Barn Panel]
    ├── Workshop → [Dashboard: Workshop Panel]
    └── Aerial → [Dashboard: Aerial Panel]
```

### Pattern 2: Cross-Site Automation (Central)
Central orchestrates based on multiple sites:

```
[NATS In: events.site.barn.climate.temperature]
    ↓
[Function: Check > 80°F]
    ↓ (if true)
[NATS Out: cmd.site.house.ac.set_temp] → Cool house
    +
[n8n HTTP: POST /webhook/barn-hot] → Notify Zeus
```

### Pattern 3: Robot Coordination (Central)
Central coordinates multiple robots:

```
[NATS In: events.site.workshop.robot.g1.task_complete]
    ↓
[Function: Next task = "transport to barn"]
    ↓
[NATS Out: cmd.site.barn.robot.go2.move_to_workshop_door]
    +
[NATS Out: cmd.site.workshop.robot.g1.move_to_barn_door]
```

### Pattern 4: Edge Autonomy (Edge)
Edge Node-RED runs local loops:

```
[NATS In: cmd.site.barn.robot.go2.patrol]
    ↓
[ROS2 Out: /cmd_vel] → Drive patrol route
    ↓
[ROS2 In: /lidar] → Read obstacles
    ↓
[Function: Avoid obstacles]
    ↓
[NATS Out: events.site.barn.robot.go2.telemetry] → Report status
```

---

## Edge Node Deployment

### Option 1: Raspberry Pi (Recommended for GPIO)
```bash
# On Raspberry Pi
bash <(curl -sL https://raw.githubusercontent.com/node-red/linux-installers/master/deb/update-nodejs-and-nodered)

# Enable as service
sudo systemctl enable nodered.service
sudo systemctl start nodered.service

# Install NATS nodes
cd ~/.node-red
npm install node-red-contrib-nats-suite

# Configure NATS server
# In Node-RED UI: Add NATS broker config → nats://central-mac-ip:4222
```

### Option 2: Docker (Recommended for Robots)
```yaml
# docker-compose.yml (barn edge node)
version: "3.8"

services:
  node-red:
    image: nodered/node-red:latest
    container_name: node-red-barn
    restart: unless-stopped
    ports:
      - "1880:1880"
    volumes:
      - node_red_barn:/data
    environment:
      - TZ=America/Los_Angeles
      - NATS_SERVER=nats://central-mac-ip:4222
      - SITE_NAME=barn
    # ROS2 bridge via host network
    network_mode: host

volumes:
  node_red_barn:
```

### Option 3: Native (Ground Station)
```bash
# On ground station Mac/Linux
npm install -g --unsafe-perm node-red
cd ~/.node-red
npm install node-red-contrib-nats-suite
npm install node-red-contrib-mavlink

# Start manually (or create systemd/LaunchAgent)
node-red
```

---

## Central Node-RED Configuration

### 1. NATS Broker Config
In Node-RED UI (http://localhost:1880):
1. Add NATS Config node
2. Server: `nats://localhost:4222`
3. Name: `OneMind NATS`

### 2. Create Dashboard
```bash
cd ~/.node-red
npm install node-red-dashboard
# Restart Node-RED
launchctl restart com.onemind.node-red
```

Access: http://localhost:1880/ui

### 3. Sample Central Flow (Fleet Monitor)
```json
[
  {
    "id": "nats-fleet-events",
    "type": "nats in",
    "name": "All Fleet Events",
    "server": "onemind-nats",
    "subject": "events.site.#",
    "queue": "",
    "wires": [["parse-event"]]
  },
  {
    "id": "parse-event",
    "type": "function",
    "name": "Parse Event",
    "func": "// Extract site, device, metric\nconst parts = msg.subject.split('.');\nmsg.site = parts[2];  // house, barn, workshop, aerial\nmsg.device = parts[3]; // sensor, robot, drone\nmsg.metric = parts[4]; // temperature, telemetry, etc\n\nmsg.payload = {\n  site: msg.site,\n  device: msg.device,\n  metric: msg.metric,\n  value: msg.payload,\n  timestamp: Date.now()\n};\n\nreturn msg;",
    "wires": [["dashboard-update", "clickhouse-store"]]
  },
  {
    "id": "dashboard-update",
    "type": "ui_text",
    "name": "Fleet Status",
    "group": "fleet-monitor",
    "wires": []
  }
]
```

---

## Management Commands

### Central Node-RED (Mac)
```bash
# Start
launchctl load ~/Library/LaunchAgents/com.onemind.node-red.plist

# Stop
launchctl unload ~/Library/LaunchAgents/com.onemind.node-red.plist

# Restart
launchctl restart com.onemind.node-red

# Logs
tail -f ~/.node-red/logs/node-red-stdout.log
tail -f ~/.node-red/logs/node-red-stderr.log

# Install nodes
cd ~/.node-red && npm install <package-name>
launchctl restart com.onemind.node-red
```

### Edge Node-RED (Raspberry Pi)
```bash
# Start/stop
sudo systemctl start nodered
sudo systemctl stop nodered
sudo systemctl restart nodered

# Logs
journalctl -u nodered -f

# Install nodes
cd ~/.node-red && npm install <package-name>
sudo systemctl restart nodered
```

---

## Security

### Central Node-RED
Create `~/.node-red/settings.js`:
```javascript
module.exports = {
  adminAuth: {
    type: "credentials",
    users: [{
      username: "zeus",
      password: "$2b$08$...",  // bcrypt hash
      permissions: "*"
    }]
  },
  httpNodeAuth: {
    user: "api",
    pass: "$2b$08$..."
  }
};
```

Generate password hash:
```bash
node-red admin hash-pw
```

### Edge Nodes
- No admin UI auth (local network only)
- NATS authentication via token (coming soon)

---

## Next Steps

### Phase 1: Central Setup (Now)
- [x] Install Node-RED on Mac
- [x] Install NATS nodes
- [x] Create LaunchAgent
- [ ] Build fleet monitoring dashboard
- [ ] Create example flows (NATS pub/sub)
- [ ] Test integration with n8n

### Phase 2: Edge Deployment (When Hardware Arrives)
- [ ] Deploy Node-RED on House Raspberry Pi
- [ ] Deploy Node-RED on Barn edge device (Docker)
- [ ] Deploy Node-RED on Workshop edge device (Docker)
- [ ] Deploy Node-RED on Aerial ground station

### Phase 3: Robot Integration
- [ ] Install ROS2 nodes on edge devices
- [ ] Create Unitree G1 control flows
- [ ] Create Unitree Go2 patrol flows
- [ ] Test central → edge command routing

### Phase 4: Production
- [ ] Enable admin auth on central
- [ ] NATS authentication tokens
- [ ] Automated backups of flows
- [ ] Monitoring/alerting for edge node health

---

## Resources

- Node-RED Docs: https://nodered.org/docs/
- NATS Suite Nodes: https://flows.nodered.org/node/node-red-contrib-nats-suite
- Dashboard: https://flows.nodered.org/node/node-red-dashboard
- ROS2 Nodes: https://flows.nodered.org/node/node-red-contrib-ros2

---

Last Updated: May 30, 2026
