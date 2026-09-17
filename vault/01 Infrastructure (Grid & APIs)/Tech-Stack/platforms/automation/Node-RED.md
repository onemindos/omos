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

## Use Cases in OneMind

1. **IoT Orchestration** - Connect sensors and devices
2. **Home Automation** - Complex automation flows
3. **API Integration** - Bridge between services
4. **Data Processing** - Transform and route data
5. **Prototyping** - Quick automation experiments

## Key Features

- Visual flow editor
- 5000+ community nodes
- MQTT, HTTP, WebSocket support
- Function nodes (JavaScript)
- Dashboard UI
- Subflows for reuse
- Context storage

## Installation

### Docker Compose

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

### With Home Assistant

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

## Useful Nodes

| Node Package                    | Purpose                    |
| ------------------------------- | -------------------------- |
| node-red-contrib-home-assistant | Home Assistant integration |
| node-red-contrib-mqtt-broker    | Built-in MQTT broker       |
| node-red-dashboard              | UI dashboards              |
| node-red-contrib-influxdb       | InfluxDB time series       |
| node-red-contrib-telegrambot    | Telegram integration       |
| node-red-contrib-discord        | Discord integration        |

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

## Node-RED + N8N Comparison

| Aspect    | Node-RED        | N8N             |
| --------- | --------------- | --------------- |
| Focus     | IoT, hardware   | Business, APIs  |
| Execution | Always running  | Trigger-based   |
| Best for  | Real-time flows | Scheduled tasks |

## Resources

- [Node-RED Docs](https://nodered.org/docs/)
- [Flow Library](https://flows.nodered.org/)
- [Node-RED Cookbook](https://cookbook.nodered.org/)
- [GitHub](https://github.com/node-red/node-red)

---

Last Updated: January 25, 2026
