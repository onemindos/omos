# Node-RED

> **Flow-based visual programming for IoT and automation**

## Overview

Node-RED is a flow-based development tool for wiring together hardware devices, APIs, and online services. Provides a visual browser-based editor for creating automation flows.

## Status

| Attribute   | Value                  |
| ----------- | ---------------------- |
| Status      | ✅ Active              |
| Category    | Automation             |
| Type        | Flow-Based Programming |
| Integration | MQTT, HTTP, WebSocket  |

## Role in OneMind OS

### Central vs Edge Architecture

**Central Node-RED (Mac)** — "The Orchestrator"
- **Role:** Fleet coordinator and cross-site automation
- **Subscribes:** `events.site.#` (all edge events)
- **Publishes:** `cmd.site.<location>.*` (commands to edges)
- **Does:** Aggregate telemetry, coordinate robots, cross-site workflows, n8n integration
- **Does NOT:** Read hardware directly (delegates to edge nodes)

**Edge Node-RED (Per Site)** — "The Translators"
- **Role:** Hardware protocol → NATS message translation
- **Reads:** ROS2, MAVLink, GPIO, Serial, Modbus
- **Publishes:** `events.site.<location>.*` (telemetry to central)
- **Subscribes:** `cmd.site.<location>.*` (commands from central)
- **Does:** Local autonomy, protocol translation, survive network outages
- **Does NOT:** Make fleet-wide decisions (central does this)

### NATS Subject Pattern
```
events.site.house.sensor.temperature     ← House sensor data
events.site.barn.robot.go2.telemetry     ← Unitree Go2 status
events.site.workshop.robot.g1.telemetry  ← Unitree G1 status
events.site.aerial.drone.telemetry       ← RECON drone position

cmd.site.house.ac.set_temp               ← Command house AC
cmd.site.barn.robot.go2.patrol           ← Command Go2 patrol
cmd.site.workshop.robot.g1.pick          ← Command G1 pick object
cmd.site.aerial.drone.goto_waypoint      ← Command drone waypoint
```

## Use Cases in OneMind

1. **Fleet Orchestration** - Coordinate multi-site robots and sensors
2. **Hardware Translation** - ROS2/MAVLink/GPIO → NATS messages
3. **Cross-Site Automation** - React to events across locations
4. **Robot Coordination** - Multi-robot task execution
5. **Edge Autonomy** - Local control survives network outages

## Key Features

- Visual flow editor
- 5000+ community nodes
- MQTT, HTTP, WebSocket support
- Function nodes (JavaScript)
- Dashboard UI
- Subflows for reuse
- Context storage

## Deployment Types

Node-RED supports multiple deployment patterns depending on use case:

### 1. Native npm (Recommended for Development)
**Best for:** Testing, development, central orchestrator

```bash
# Install globally
npm install -g --unsafe-perm node-red

# Run
node-red

# Access: http://localhost:1880
```

**Pros:**
- Direct access to system (GPIO, serial)
- Easy npm package installation
- Fast iteration
- Native performance

**Cons:**
- Manual process management
- No isolation
- Requires Node.js installed

### 2. Docker (Recommended for Production)
**Best for:** Production deployments, consistent environments

```yaml
version: "3.8"

services:
  node-red:
    image: nodered/node-red:latest
    container_name: node-red
    restart: unless-stopped
    ports:
      - "1880:1880"
    volumes:
      - node_red_data:/data
    environment:
      - TZ=America/Los_Angeles
      - NODE_RED_CREDENTIAL_SECRET=${NODE_RED_SECRET}

volumes:
  node_red_data:
```

**Pros:**
- Isolated environment
- Easy updates (pull new image)
- Portable across hosts
- Auto-restart on failure

**Cons:**
- No direct GPIO access (must use remote GPIO)
- Slightly higher resource usage
- Network overhead for hardware access

### 3. Raspberry Pi (Edge Deployment)
**Best for:** Edge nodes with sensors/GPIO

```bash
# Pre-installed on Raspberry Pi OS
# Or install manually:
bash <(curl -sL https://raw.githubusercontent.com/node-red/linux-installers/master/deb/update-nodejs-and-nodered)

# Enable as service
sudo systemctl enable nodered.service
sudo systemctl start nodered.service
```

**Pros:**
- Direct GPIO access
- Low power consumption
- Runs at the edge (low latency)
- Survives network outages

**Cons:**
- Limited resources (RAM/CPU)
- SD card reliability issues
- Manual updates per device

### 4. Home Assistant Add-on
**Best for:** HA-centric deployments

```yaml
services:
  node-red:
    image: nodered/node-red:latest
    ports:
      - "1880:1880"
    volumes:
      - node_red_data:/data
    environment:
      - TZ=America/Los_Angeles
    depends_on:
      - mosquitto
      - homeassistant
```

**Pros:**
- Tight HA integration
- Managed by HA supervisor
- Automatic backups
- Single UI for management

**Cons:**
- Tied to Home Assistant lifecycle
- Limited control over Node-RED version
- Supervisor overhead

### 5. Cloud Hosted (FlowFuse, AWS, Azure)
**Best for:** Team collaboration, multi-site management

**FlowFuse (Node-RED SaaS):**
```bash
# Managed Node-RED instances
# Visit: https://flowfuse.com
```

**AWS ECS:**
```yaml
# Run nodered/node-red:latest in ECS task
# Use EFS for persistent /data volume
```

**Pros:**
- No infrastructure management
- High availability
- Team collaboration features
- Centralized logging/monitoring

**Cons:**
- Subscription cost
- Internet dependency
- Higher latency to local hardware
- Data sovereignty concerns

### 6. Kubernetes (Enterprise)
**Best for:** Large fleets, multi-tenant, production scale

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: node-red
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: node-red
        image: nodered/node-red:latest
        ports:
        - containerPort: 1880
        volumeMounts:
        - name: data
          mountPath: /data
```

**Pros:**
- Horizontal scaling
- Load balancing
- Self-healing
- GitOps workflows

**Cons:**
- Complex setup
- Overkill for small deployments
- Kubernetes knowledge required

---

## Recommended Deployment for OneMind OS

### Current (n8n only)
```
Mac/VPS (Central)
  └── n8n (native npm, SQLite)
      └── NATS (native LaunchAgent)
```

### Future (n8n + Node-RED edge)
```
Mac/VPS (Central)
  ├── n8n (native npm) ← SaaS/webhooks/orchestration
  └── NATS (LaunchAgent) ← Message backbone

House (Raspberry Pi)
  └── Node-RED (systemd service) ← GPIO sensors → NATS

Barn (Edge device)
  └── Node-RED (Docker) ← Unitree Go2 ROS2 → NATS

Workshop (Edge device)
  └── Node-RED (Docker) ← Unitree G1 ROS2 → NATS

Aerial (Ground station)
  └── Node-RED (native) ← MAVLink parser → NATS
```

**Architecture:**
- **n8n** = Central, SaaS-facing
- **Node-RED** = Edge, hardware-facing
- **NATS** = Universal transport connecting both

---

## Installation

## Flow Examples

### MQTT to HTTP

```json
[
  {
    "id": "mqtt-in-1",
    "type": "mqtt in",
    "topic": "sensors/temperature",
    "broker": "mqtt-broker",
    "wires": [["http-req-1"]]
  },
  {
    "id": "http-req-1",
    "type": "http request",
    "method": "POST",
    "url": "https://api.example.com/data",
    "wires": [["debug-1"]]
  }
]
```

### Function Node

```javascript
// Transform incoming message
const data = msg.payload;

msg.payload = {
  timestamp: new Date().toISOString(),
  sensor: msg.topic,
  value: parseFloat(data),
  unit: "celsius",
};

return msg;
```

## API Integration

### Admin API

```python
import httpx

NODE_RED_URL = "http://localhost:1880"

async def get_flows():
    """Get all flows from Node-RED."""
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{NODE_RED_URL}/flows")
        return response.json()

async def deploy_flows(flows: list):
    """Deploy flows to Node-RED."""
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{NODE_RED_URL}/flows",
            json=flows,
            headers={"Node-RED-Deployment-Type": "full"}
        )
        return response.status_code == 204

async def inject_node(node_id: str):
    """Trigger an inject node."""
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{NODE_RED_URL}/inject/{node_id}"
        )
        return response.status_code == 200
```

### HTTP Endpoints

Create HTTP endpoints in Node-RED:

```json
[
  {
    "id": "http-in-1",
    "type": "http in",
    "url": "/api/trigger",
    "method": "post",
    "wires": [["function-1"]]
  },
  {
    "id": "function-1",
    "type": "function",
    "func": "msg.payload = {status: 'triggered', time: Date.now()}; return msg;",
    "wires": [["http-out-1"]]
  },
  {
    "id": "http-out-1",
    "type": "http response",
    "statusCode": "200",
    "wires": []
  }
]
```

```python
# Call from Python
async def trigger_nodered_flow(data: dict):
    """Trigger a Node-RED HTTP endpoint."""
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{NODE_RED_URL}/api/trigger",
            json=data
        )
        return response.json()
```

## MQTT Integration

```python
import asyncio
from nats.aio.client import Client as NATS

async def bridge_nats_to_nodered():
    """Bridge NATS messages to Node-RED via MQTT."""
    nc = await NATS().connect("nats://localhost:4222")

    async def message_handler(msg):
        # Forward to MQTT for Node-RED
        await mqtt_publish(
            topic=f"nodered/{msg.subject}",
            payload=msg.data.decode()
        )

    await nc.subscribe("events.*", cb=message_handler)
```

## Legacy AI Integration

```python
from agno import Tool

@Tool
async def trigger_nodered_automation(
    flow_name: str,
    data: dict = None
) -> dict:
    """Trigger a Node-RED automation flow.

    Args:
        flow_name: Name/endpoint of the flow to trigger
        data: Optional data to send to the flow
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{NODE_RED_URL}/api/{flow_name}",
            json=data or {}
        )
        return response.json()

@Tool
async def get_nodered_status() -> dict:
    """Get current Node-RED runtime status."""
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{NODE_RED_URL}/settings")
        settings = response.json()

        return {
            "version": settings.get("version"),
            "context_storage": settings.get("contextStorage"),
            "running": True
        }

@Tool
async def send_to_nodered(
    topic: str,
    payload: dict
) -> str:
    """Send data to Node-RED via MQTT.

    Args:
        topic: MQTT topic for Node-RED to receive
        payload: Data to send
    """
    await mqtt_publish(topic, payload)
    return f"Sent to Node-RED topic: {topic}"
```

## Installed Nodes (OneMind Stack)

### ✅ Currently Installed (13 packages)

**Core Infrastructure:**
| Node Package | Purpose | Status |
|-------------|---------|--------|
| **node-red-contrib-nats-suite** | NATS Core + JetStream + KV + Object Store | ✅ Installed |
| **node-red-contrib-home-assistant-websocket** | Home Assistant WebSocket API (live events) | ✅ Installed |
| **node-red-dashboard** | UI dashboards at :1880/ui | ✅ Installed |

**Communication:**
| Node Package | Purpose | Status |
|-------------|---------|--------|
| **node-red-contrib-telegrambot** | Telegram bot API | ✅ Installed |
| **node-red-contrib-discord-advanced** | Discord webhooks/bots | ✅ Installed |

**Hardware & IoT:**
| Node Package | Purpose | Status |
|-------------|---------|--------|
| **node-red-contrib-mf-mavlink** | MAVLink parser (drone telemetry) | ✅ Installed |
| **node-red-node-pi-gpio** | Raspberry Pi GPIO direct access | ✅ Installed |
| **node-red-node-serialport** | Serial/UART communication (Arduino) | ✅ Installed |
| **node-red-contrib-modbus** | Modbus industrial protocol | ✅ Installed |
| **node-red-contrib-onvif-nodes** | ONVIF camera control | ✅ Installed |

**AI & Data:**
| Node Package | Purpose | Status |
|-------------|---------|--------|
| **node-red-contrib-openai-ubos** | OpenAI API (GPT-4, etc.) | ✅ Installed |
| **node-red-contrib-influxdb** | InfluxDB time-series database | ✅ Installed |
| **node-red-contrib-spotify** | Spotify playback control | ✅ Installed |

### ⚠️ Not Installed (Requires Dependencies)

**ROS2 Robotics:**
| Node Package | Why Not Installed | Install When |
|-------------|-------------------|--------------|
| @chart-sg/node-red-ros2 | Requires ROS2 system installation (native build) | Install ROS2 first, then on edge nodes with robots |

**Note:** ROS2 nodes require ROS2 to be installed on the system (via apt/brew). When you deploy edge nodes for Unitree G1/Go2, install ROS2 first, then install this node package.

## Install Nodes

```bash
# Inside Node-RED container or npm
cd /data
npm install node-red-contrib-home-assistant-websocket
npm install node-red-dashboard
```

## Environment Variables

```bash
NODE_RED_URL=http://localhost:1880
NODE_RED_SECRET=your_credential_secret
NODE_RED_ADMIN_USER=admin
NODE_RED_ADMIN_PASSWORD=password
```

## Security

```javascript
// settings.js
module.exports = {
  adminAuth: {
    type: "credentials",
    users: [
      {
        username: "admin",
        password: "$2b$08$...", // bcrypt hash
        permissions: "*",
      },
    ],
  },
  httpNodeAuth: {
    user: "api",
    pass: "$2b$08$...",
  },
};
```

## Best Practices

1. **Use subflows** - Reusable flow components
2. **Error handling** - Add catch nodes
3. **Context storage** - Use file context for persistence
4. **Flow organization** - Separate flows by function
5. **Comments** - Document complex flows
6. **Version control** - Export flows to git

## Node-RED vs n8n — Protocol Breakdown

| Protocol | Used By | n8n Support | Node-RED Support | Winner |
|----------|---------|-------------|------------------|--------|
| **REST API** | FarmBot, HA, most SaaS | ✅ Native HTTP node | ✅ HTTP node | **n8n** (better auth/retry) |
| **WebSocket** | HA event stream | ✅ Via custom code | ✅ Native | **Tie** |
| **NATS** | Your backbone | ✅ `@zetanova/n8n-nodes-nats` | ✅ `node-red-contrib-nats-suite` | **Tie** |
| **ROS2** | Unitree G1/Go2 | ❌ No native | ✅ `node-red-contrib-ros2` | **Node-RED wins** |
| **MAVLink** | Drones (RECON) | ❌ No native | ✅ `node-red-contrib-mavlink` | **Node-RED wins** |
| **GPIO** | Raspberry Pi sensors | ❌ No | ✅ Native `rpi-gpio` | **Node-RED wins** |
| **Serial** | Arduino, sensors | ❌ No | ✅ Native `serialport` | **Node-RED wins** |
| **MQTT** | IoT devices | ✅ Native | ✅ Native | Tie |
| **Modbus** | Industrial sensors | ❌ No | ✅ `node-red-contrib-modbus` | **Node-RED wins** |

## When You Need Node-RED

### Scenario 1: Unitree G1/Go2 ROS2 Integration
When robots are deployed:
```
Unitree Robot (ROS2 topics)
    ↓
Node-RED (node-red-contrib-ros2)
    ↓
NATS (cmd.robot.g1.>)
    ↓
n8n / Hermes (orchestration)
```

**Why Node-RED:** Native ROS2 bridge exists, n8n doesn't have one.

### Scenario 2: RECON Drone (MAVLink)
When drone is operational:
```
PX4 Autopilot (MAVLink telemetry)
    ↓
Node-RED (node-red-contrib-mavlink)
    ↓
NATS (events.drone.telemetry)
    ↓
ClickHouse (time-series storage)
```

**Why Node-RED:** MAVLink is a low-level binary protocol, Node-RED has parsers.

### Scenario 3: Edge Sensor Network (Future)
Raspberry Pi + sensors (soil moisture, temp, etc.):
```
Raspberry Pi GPIO Sensors
    ↓
Node-RED (on the Pi, GPIO nodes)
    ↓
NATS (events.sensors.*)
    ↓
ClickHouse
```

**Why Node-RED:** Runs on the Pi, native GPIO, low resource usage.

### When to Use n8n Instead
- **FarmBot** — REST API → n8n handles HTTP better
- **Home Assistant** — REST API → n8n for orchestration
- **SaaS services** — Notion, Linear, GitHub → n8n has 500+ integrations
- **Webhooks** — Inbound from services → n8n native
- **Notifications** — Telegram, Discord → n8n has better retry/routing

## Architecture Decision

```
┌─────────────────────────────────────────────────────────────┐
│                      OneMind OS Fleet                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  HIGH-LEVEL (SaaS, APIs, webhooks)                          │
│  ════════════════════════════════════════════                │
│  n8n (primary)                                               │
│    • SaaS integrations (Notion, Linear, GitHub)             │
│    • Webhooks (inbound from services)                       │
│    • Notification routing (Telegram, Discord)               │
│    • FarmBot REST API                                        │
│    • Home Assistant REST API                                 │
│                                                               │
│  LOW-LEVEL (hardware, protocols, edge)                      │
│  ════════════════════════════════════════════                │
│  Node-RED (edge companion)                                   │
│    • ROS2 (Unitree G1/Go2 robots)                           │
│    • MAVLink (RECON drone)                                   │
│    • GPIO (Raspberry Pi sensors)                            │
│    • Serial (Arduino, custom hardware)                      │
│    • Modbus (industrial sensors)                            │
│                                                               │
│  UNIVERSAL BACKBONE                                          │
│  ════════════════════════════════════════════                │
│  NATS (central messaging)                                    │
│    • n8n publishes/subscribes                               │
│    • Node-RED publishes/subscribes                          │
│    • Hermes subscribes                                       │
│    • ClickHouse stores all events                           │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

## Deployment Pattern

**n8n** — Central orchestrator (runs on Mac/VPS)  
**Node-RED** — Edge instances (one per physical site)

```
House Node (Raspberry Pi)
  └── Node-RED (GPIO sensors + HA bridge) → NATS

Barn Node (Edge device)
  └── Node-RED (Unitree Go2 + climate sensors) → NATS

Workshop Node (Edge device)
  └── Node-RED (Unitree G1 + tools inventory) → NATS

Aerial Node (Drone)
  └── Node-RED (MAVLink parser on ground station) → NATS
```

### Current Status: n8n Only
**Today's needs:** SaaS integration (Notion, Linear, GitHub), webhooks, notifications  
**Today's tool:** n8n (already deployed)

**Deploy Node-RED when:**
1. Unitree G1/Go2 arrives (ROS2 bridge needed)
2. RECON drone is operational (MAVLink parsing)
3. You add edge sensors (Raspberry Pi GPIO)

## Node-RED + n8n: Complementary Tools

| Aspect | Node-RED | n8n |
| --- | --- | --- |
| Focus | IoT, hardware | Business, APIs |
| Execution | Always running | Trigger-based |
| Best for | Real-time flows | Scheduled tasks |
| Protocols | ROS2, MAVLink, GPIO, Serial, Modbus | REST, WebSocket, webhooks |
| Deployment | Edge devices (Raspberry Pi) | Central server (Mac/VPS) |
| Integration | Hardware-first | SaaS-first |

## Resources

- [Node-RED Docs](https://nodered.org/docs/)
- [Flow Library](https://flows.nodered.org/)
- [Node-RED Cookbook](https://cookbook.nodered.org/)
- [GitHub](https://github.com/node-red/node-red)

---

Last Updated: May 30, 2026
