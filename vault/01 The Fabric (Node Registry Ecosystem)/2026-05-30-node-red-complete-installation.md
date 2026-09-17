---
title: "Node-RED Complete Installation — All Nodes"
date: 2026-05-30
type: complete
status: ready
---

# Node-RED Complete Installation — All Nodes ✅

## What Was Installed

### ✅ 13 Node Packages — Ready to Use

**Core Infrastructure (3):**
- ✅ `node-red-contrib-nats-suite` — NATS Core + JetStream + KV + Object Store
- ✅ `node-red-contrib-home-assistant-websocket` — Home Assistant WebSocket API
- ✅ `node-red-dashboard` — UI dashboards at http://localhost:1880/ui

**Communication (2):**
- ✅ `node-red-contrib-telegrambot` — Telegram bot API
- ✅ `node-red-contrib-discord-advanced` — Discord webhooks/bots

**Hardware & IoT (5):**
- ✅ `node-red-contrib-mf-mavlink` — MAVLink parser (RECON drone telemetry)
- ✅ `node-red-node-pi-gpio` — Raspberry Pi GPIO direct access
- ✅ `node-red-node-serialport` — Serial/UART (Arduino, sensors)
- ✅ `node-red-contrib-modbus` — Modbus industrial protocol
- ✅ `node-red-contrib-onvif-nodes` — ONVIF camera control

**AI & Data (3):**
- ✅ `node-red-contrib-openai-ubos` — OpenAI API (GPT-4, etc.)
- ✅ `node-red-contrib-influxdb` — InfluxDB time-series database
- ✅ `node-red-contrib-spotify` — Spotify playback control

---

## Coverage: Your Entire Hardware Stack

From your `05 Stack (Arsenal Ecosystem)` documentation, Node-RED now covers:

### ✅ Active Systems (Now)
- **NATS** (infrastructure/messaging/) → `nats-suite`
- **Home Assistant** (iot/smart-home/) → `home-assistant-websocket`
- **Telegram** (communication/messaging/) → `telegrambot`
- **Discord** (communication/messaging/) → `discord-advanced`
- **Cameras** (iot/surveillance/ — Amcrest, Reolink, Wyze) → `onvif-nodes`
- **Spotify** (personal/productivity/) → `spotify`

### 🔜 Planned Hardware (Ready When Arrives)
- **RECON Drone** (hardware/drones/) → `mf-mavlink` (ready)
- **Raspberry Pi sensors** (iot/) → `pi-gpio`, `serialport` (ready)
- **Arduino/sensors** → `serialport` (ready)
- **Industrial sensors** → `modbus` (ready)

### ⚠️ Requires Pre-Installation
- **Unitree G1/Go2** (hardware/robotics/platforms/) → `@chart-sg/node-red-ros2` (install ROS2 first)

---

## What Each Node Does

### NATS Suite
**Use:** Central/edge messaging backbone
```
[NATS In: events.site.barn.robot.go2.telemetry]
    ↓
[Function: Process telemetry]
    ↓
[NATS Out: events.alert.robot_battery_low]
```

### Home Assistant WebSocket
**Use:** Real-time HA events
```
[HA Events: state] (sensor.temperature)
    ↓
[Function: Check threshold]
    ↓
[NATS Out: events.site.house.sensor.temperature]
```

### Dashboard
**Use:** Fleet monitoring UI (http://localhost:1880/ui)
```
[NATS In: events.site.#]
    ↓
[ui_gauge] (temperature)
[ui_chart] (battery levels)
[ui_text] (robot status)
```

### Telegram Bot
**Use:** Send/receive Telegram messages
```
[Telegram Receiver]
    ↓
[Function: Parse command]
    ↓
[NATS Out: cmd.site.barn.robot.go2.patrol]
    ↓
[Telegram Sender] "Patrol started"
```

### Discord Advanced
**Use:** Discord webhooks + bot commands
```
[NATS In: events.alert.#]
    ↓
[Discord Send Webhook] (#alerts channel)
```

### MAVLink Parser
**Use:** RECON drone telemetry
```
[Serial In] (MAVLink from drone)
    ↓
[MAVLink Parse]
    ↓
[NATS Out: events.site.aerial.drone.telemetry]
```

### Pi GPIO
**Use:** Raspberry Pi sensors (edge nodes)
```
[Pi GPIO In] (pin 17 — motion sensor)
    ↓
[Function: Detect motion]
    ↓
[NATS Out: events.site.house.sensor.motion]
```

### Serial Port
**Use:** Arduino/sensor communication
```
[Serial In] (/dev/ttyUSB0 — Arduino)
    ↓
[Function: Parse sensor data]
    ↓
[NATS Out: events.site.barn.sensor.soil_moisture]
```

### Modbus
**Use:** Industrial sensors
```
[Modbus Read] (Register 40001)
    ↓
[Function: Scale value]
    ↓
[NATS Out: events.site.workshop.sensor.pressure]
```

### ONVIF
**Use:** Camera control (Amcrest, Reolink, Wyze)
```
[NATS In: cmd.site.house.camera.ptz_preset]
    ↓
[ONVIF Command] (camera IP, preset 1)
```

### OpenAI
**Use:** AI workflows
```
[NATS In: events.site.barn.anomaly_detected]
    ↓
[OpenAI GPT-4] "Analyze this telemetry"
    ↓
[Telegram Sender] (send analysis to Zeus)
```

### InfluxDB
**Use:** Time-series storage (alternative to ClickHouse)
```
[NATS In: events.site.#]
    ↓
[InfluxDB Write] (measurement: telemetry)
```

### Spotify
**Use:** Music automation
```
[NATS In: cmd.home.music.play_morning_playlist]
    ↓
[Spotify Play] (playlist: Morning Focus)
```

---

## Not Installed: ROS2

**Why:** Requires ROS2 system installation (apt/brew)

**When to install:**
1. Deploy edge node for Unitree G1/Go2
2. Install ROS2 on that node: `brew install ros` (macOS) or `apt install ros-humble-desktop` (Ubuntu)
3. Then install: `cd ~/.node-red && npm install @chart-sg/node-red-ros2`

**Use case:**
```
[ROS2 Subscribe: /joint_states]
    ↓
[Function: Extract arm position]
    ↓
[NATS Out: events.site.workshop.robot.g1.arm_position]
```

---

## Installation Commands (Already Run)

```bash
cd ~/.node-red

# Core
npm install node-red-contrib-nats-suite
npm install node-red-contrib-home-assistant-websocket
npm install node-red-dashboard

# Communication
npm install node-red-contrib-telegrambot
npm install node-red-contrib-discord-advanced

# Hardware
npm install node-red-contrib-mf-mavlink
npm install node-red-node-pi-gpio
npm install node-red-node-serialport
npm install node-red-contrib-modbus
npm install node-red-contrib-onvif-nodes

# AI & Data
npm install node-red-contrib-openai-ubos
npm install node-red-contrib-influxdb
npm install node-red-contrib-spotify
```

---

## Current Status

| Component | Status | Access |
|-----------|--------|--------|
| Node-RED Central | ✅ Running | http://localhost:1880 |
| Dashboard | ✅ Ready | http://localhost:1880/ui |
| NATS Nodes | ✅ 13 installed | Ready to configure |
| HA Integration | ✅ Ready | Configure HA server |
| Telegram Bot | ✅ Ready | Need bot token from @BotFather |
| Discord Bot | ✅ Ready | Need webhook URL |
| MAVLink | ✅ Ready | When drone arrives |
| GPIO/Serial | ✅ Ready | When Pi/Arduino deployed |
| Modbus | ✅ Ready | When industrial sensors |
| ONVIF Cameras | ✅ Ready | Configure camera IPs |
| OpenAI | ✅ Ready | Need API key in env |
| Spotify | ✅ Ready | OAuth setup required |
| ROS2 | ⚠️ Not installed | Install ROS2 first |

---

## Next Steps

### 1. Configure NATS (Priority 1)
```
1. Open http://localhost:1880
2. Add NATS config node
3. Server: nats://localhost:4222
4. Name: OneMind NATS
5. Deploy
```

### 2. Test NATS Integration
```
1. Create flow: [NATS In: events.test] → [Debug]
2. Deploy
3. Terminal: nats pub events.test "hello from terminal"
4. Check debug panel
```

### 3. Create Fleet Dashboard
```
1. Add UI nodes (gauge, chart, text)
2. Connect to NATS events
3. Access: http://localhost:1880/ui
```

### 4. Connect Home Assistant (Optional)
```
1. Add HA Events node
2. Configure HA URL: http://homeassistant.local:8123
3. Subscribe to entity state changes
4. Publish to NATS
```

### 5. Setup Telegram Bot (Optional)
```
1. Create bot: @BotFather → /newbot
2. Copy token
3. Add Telegram Receiver node
4. Configure token
5. Test: send message to bot
```

---

## Architecture Now Complete

```
┌──────────────────────────────────────────────────────────────┐
│                    Central Node-RED (Mac)                     │
├──────────────────────────────────────────────────────────────┤
│                                                                │
│  ✅ NATS                    ← Events + commands               │
│  ✅ Home Assistant          ← Live events                     │
│  ✅ Dashboard               ← Fleet monitoring UI             │
│  ✅ Telegram/Discord        ← Notifications                   │
│  ✅ MAVLink                 ← Drone telemetry (ready)         │
│  ✅ GPIO/Serial/Modbus      ← Sensors (ready)                 │
│  ✅ ONVIF                   ← Camera control                  │
│  ✅ OpenAI/InfluxDB/Spotify ← AI + data + music              │
│                                                                │
└──────────────────────────────────────────────────────────────┘
           ↓ NATS (events.site.*, cmd.site.*)
┌──────────────────────────────────────────────────────────────┐
│                      Edge Nodes (Future)                      │
├──────────────────────────────────────────────────────────────┤
│                                                                │
│  House (Pi)      → GPIO + HA → NATS                          │
│  Barn (Docker)   → ROS2 (Go2) + MAVLink → NATS              │
│  Workshop (Docker) → ROS2 (G1) → NATS                        │
│  Aerial (Native) → MAVLink parser → NATS                     │
│                                                                │
└──────────────────────────────────────────────────────────────┘
```

---

## Files Updated

1. `/00-24 SO (Sovereign Operations)/05 Stack (Arsenal Ecosystem)/iot/automation/node-red.md`
   - Updated installed nodes section (13 packages)
   - Reorganized by category
   - Added ROS2 note (requires system install)

2. `/00-24 SO (Sovereign Operations)/24 Inbox/2026-05-30-node-red-complete-installation.md`
   - This summary document

---

## Summary

✅ **13 Node packages installed** covering your entire hardware stack  
✅ **Central orchestrator ready** for fleet coordination  
✅ **All active systems covered** (NATS, HA, Telegram, Discord, cameras, Spotify)  
✅ **Future hardware ready** (MAVLink, GPIO, Serial, Modbus for when hardware arrives)  
⚠️ **ROS2 deferred** (requires system installation — do on edge nodes when robots arrive)

**Node-RED is now fully equipped to orchestrate OneMind OS fleet operations.**

---

Last Updated: May 30, 2026
