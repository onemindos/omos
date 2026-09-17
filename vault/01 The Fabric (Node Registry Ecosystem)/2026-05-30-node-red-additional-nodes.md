---
title: "Node-RED Additional Nodes — Installed"
date: 2026-05-30
type: update
status: complete
---

# Node-RED Additional Nodes — Installed ✅

## What Was Installed

### ✅ Now Available in Node-RED UI

| Node Package | Purpose | Use Case |
|-------------|---------|----------|
| **node-red-contrib-nats-suite** | NATS Core + JetStream + KV + Object Store | Central/Edge messaging |
| **node-red-contrib-home-assistant-websocket** | Home Assistant WebSocket API | Real-time HA events |
| **node-red-dashboard** | UI dashboards | Fleet monitoring at :1880/ui |
| **node-red-contrib-telegrambot** | Telegram Bot API | Telegram automation |

### Access
- **UI:** http://localhost:1880 (editor)
- **Dashboard:** http://localhost:1880/ui (after you create dashboard flows)

---

## Based on Your Hardware Stack

From `/00-24 SO (Sovereign Operations)/05 Stack (Arsenal Ecosystem)/`:

### Currently Installed (Matching Active Systems)
✅ **NATS** — You have NATS running, now Node-RED can publish/subscribe  
✅ **Home Assistant** — You have HA active, now Node-RED can read events  
✅ **Telegram** — In your messaging stack, now Node-RED can send/receive  
✅ **Dashboard** — For fleet monitoring UI

### Future Nodes (When Hardware Arrives)

**Robotics (hardware/robotics/):**
- `node-red-contrib-ros2` — For Unitree G1/Go2 (ROS2 topics)
- Install on edge nodes when robots arrive

**Drones (hardware/drones/):**
- `node-red-contrib-mavlink` — For RECON drone (MAVLink telemetry)
- Install on ground station when drone operational

**IoT Sensors (iot/):**
- `node-red-node-pi-gpio` — Raspberry Pi GPIO
- `node-red-node-serialport` — Arduino/serial sensors
- Install on edge Raspberry Pi nodes

**Additional Integrations (based on stack):**
- `node-red-contrib-discord-advanced` — Discord (communication/messaging/)
- `node-red-contrib-openai-ubos` — OpenAI (ai/providers/)
- `node-red-contrib-anthropic-ubos` — Anthropic Claude (ai/providers/)
- `node-red-contrib-onvif-nodes` — ONVIF cameras (iot/surveillance/)
- `node-red-contrib-spotify` — Spotify (personal/productivity/)

---

## What Changed in node-red.md

✅ **Added Central vs Edge Architecture section**
- Central Node-RED (Mac) — "The Orchestrator"
- Edge Node-RED (Per Site) — "The Translators"
- NATS subject pattern (`events.site.*`, `cmd.site.*`)

✅ **Updated Use Cases** to match fleet architecture:
- Fleet orchestration
- Hardware translation
- Cross-site automation
- Robot coordination
- Edge autonomy

✅ **Replaced "Useful Nodes" section** with:
- ✅ Currently Installed (4 packages)
- 📦 Recommended for Future Hardware (organized by category)

---

## Installation Commands

### Already Done
```bash
cd ~/.node-red
npm install node-red-contrib-nats-suite
npm install node-red-contrib-home-assistant-websocket
npm install node-red-dashboard
npm install node-red-contrib-telegrambot
```

### When Robots Arrive (Edge Nodes)
```bash
# On edge device/Pi
cd ~/.node-red
npm install node-red-contrib-ros2
npm install node-red-contrib-mavlink
npm install node-red-node-pi-gpio
npm install node-red-node-serialport
```

### Optional (As Needed)
```bash
cd ~/.node-red
npm install node-red-contrib-discord-advanced
npm install node-red-contrib-openai-ubos
npm install node-red-contrib-anthropic-ubos
npm install node-red-contrib-onvif-nodes
npm install node-red-contrib-spotify
```

---

## How to Use the New Nodes

### 1. NATS Integration
1. Open http://localhost:1880
2. Add "nats in" node (subscribe)
3. Add "nats out" node (publish)
4. Configure NATS server: `nats://localhost:4222`
5. Test: `nats pub events.test "hello from terminal"`

### 2. Home Assistant Events
1. Add "events: state" node
2. Configure HA server URL
3. Subscribe to entity state changes
4. Connect to NATS out → `events.site.house.sensor.*`

### 3. Dashboard
1. Add "ui_text", "ui_gauge", "ui_chart" nodes
2. Create groups/tabs
3. Access: http://localhost:1880/ui
4. Display fleet telemetry in real-time

### 4. Telegram Bot
1. Create bot with @BotFather
2. Add "telegram receiver" node
3. Add "telegram sender" node
4. Connect to your workflows

---

## Example Flow: HA → NATS → Dashboard

```
[HA Events: state] (temperature sensor)
    ↓
[Function: Transform]
    ↓
[NATS Out: events.site.house.sensor.temperature]
    ↓
[Dashboard: Gauge] (display current temp)
```

---

## Status

| Component | Status |
|-----------|--------|
| Node-RED Central | ✅ Running (:1880) |
| NATS Nodes | ✅ Installed |
| HA Nodes | ✅ Installed |
| Dashboard | ✅ Installed (:1880/ui) |
| Telegram Nodes | ✅ Installed |
| Robot Nodes | 🔜 When hardware arrives |
| Drone Nodes | 🔜 When hardware arrives |
| Edge GPIO Nodes | 🔜 When Pi deployed |

---

## Next Steps

1. **Configure NATS in Node-RED UI** (add server config)
2. **Test NATS pub/sub** (simple flow)
3. **Connect Home Assistant** (if you want live events)
4. **Create fleet dashboard** (http://localhost:1880/ui)

---

Last Updated: May 30, 2026
