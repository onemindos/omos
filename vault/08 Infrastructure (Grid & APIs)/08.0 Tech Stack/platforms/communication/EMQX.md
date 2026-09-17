# EMQX

> **Enterprise MQTT broker for IoT messaging**

## Overview

EMQX is a highly scalable, distributed MQTT message broker for IoT. Handles millions of concurrent connections with low latency, making it ideal for connecting devices, sensors, and agents.

## Status

| Attribute   | Value               |
| ----------- | ------------------- |
| Status      | ✅ Active           |
| Category    | Communication       |
| Type        | MQTT Broker         |
| Integration | IoT, Home Assistant |

## Use Cases in OneMind

1. **IoT Device Messaging** - Connect all sensors and devices
2. **Real-time Events** - Publish/subscribe architecture
3. **Home Automation** - MQTT backbone for smart home
4. **Agent Communication** - Inter-agent messaging
5. **FarmBot Integration** - Agricultural IoT

## Key Features

- MQTT 5.0 support
- Millions of concurrent connections
- Cluster and high availability
- Rule engine for data processing
- WebSocket support
- REST API
- Authentication/Authorization
- Message persistence

## Installation

### Docker Compose

```yaml
version: "3.8"

services:
  emqx:
    image: emqx/emqx:latest
    container_name: emqx
    restart: unless-stopped
    ports:
      - "1883:1883" # MQTT
      - "8083:8083" # WebSocket
      - "8084:8084" # WSS
      - "8883:8883" # MQTTS
      - "18083:18083" # Dashboard
    volumes:
      - emqx_data:/opt/emqx/data
      - emqx_log:/opt/emqx/log
    environment:
      - EMQX_NAME=onemind-emqx
      - EMQX_HOST=0.0.0.0
      - EMQX_DASHBOARD__DEFAULT_USERNAME=admin
      - EMQX_DASHBOARD__DEFAULT_PASSWORD=${EMQX_ADMIN_PASSWORD}

volumes:
  emqx_data:
  emqx_log:
```

## MQTT Topics Structure

```
onemind/
├── devices/
│   ├── {device_id}/status
│   ├── {device_id}/telemetry
│   └── {device_id}/command
├── agents/
│   ├── legacy/events
│   ├── legacy/commands
│   └── {agent}/status
├── home/
│   ├── {room}/sensors
│   ├── {room}/lights
│   └── {room}/climate
├── farm/
│   ├── sensors/soil
│   ├── sensors/weather
│   └── farmbot/status
└── events/
    ├── security
    ├── notifications
    └── system
```

## Python Integration

### paho-mqtt Client

```python
import paho.mqtt.client as mqtt
import json

class EMQXClient:
    def __init__(self, host: str, port: int = 1883):
        self.client = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
        self.host = host
        self.port = port

    def connect(self, username: str = None, password: str = None):
        if username and password:
            self.client.username_pw_set(username, password)

        self.client.on_connect = self._on_connect
        self.client.on_message = self._on_message
        self.client.connect(self.host, self.port, 60)
        self.client.loop_start()

    def _on_connect(self, client, userdata, flags, rc, properties):
        print(f"Connected with result code {rc}")

    def _on_message(self, client, userdata, msg):
        print(f"{msg.topic}: {msg.payload.decode()}")

    def publish(self, topic: str, payload: dict):
        self.client.publish(topic, json.dumps(payload))

    def subscribe(self, topic: str, callback=None):
        if callback:
            self.client.message_callback_add(topic, callback)
        self.client.subscribe(topic)
```

### Async with aiomqtt

```python
import asyncio
import aiomqtt
import json

async def mqtt_publish(topic: str, payload: dict):
    """Publish message to EMQX."""
    async with aiomqtt.Client(
        hostname="localhost",
        port=1883,
        username="onemind",
        password="password"
    ) as client:
        await client.publish(topic, json.dumps(payload))

async def mqtt_subscribe(topic: str, handler):
    """Subscribe to EMQX topic."""
    async with aiomqtt.Client(
        hostname="localhost",
        port=1883,
        username="onemind",
        password="password"
    ) as client:
        async with client.messages() as messages:
            await client.subscribe(topic)
            async for message in messages:
                payload = json.loads(message.payload.decode())
                await handler(message.topic, payload)
```

## REST API

```python
import httpx

EMQX_API_URL = "http://localhost:18083/api/v5"
EMQX_API_KEY = "your_api_key"
EMQX_API_SECRET = "your_api_secret"

async def get_clients():
    """Get all connected MQTT clients."""
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{EMQX_API_URL}/clients",
            auth=(EMQX_API_KEY, EMQX_API_SECRET)
        )
        return response.json()

async def publish_via_api(topic: str, payload: str):
    """Publish message via REST API."""
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{EMQX_API_URL}/publish",
            auth=(EMQX_API_KEY, EMQX_API_SECRET),
            json={
                "topic": topic,
                "payload": payload,
                "qos": 1
            }
        )
        return response.status_code == 200

async def get_subscriptions():
    """Get all active subscriptions."""
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{EMQX_API_URL}/subscriptions",
            auth=(EMQX_API_KEY, EMQX_API_SECRET)
        )
        return response.json()
```

## Legacy AI Integration

```python
from agno import Tool

@Tool
async def publish_mqtt_message(
    topic: str,
    payload: dict
) -> str:
    """Publish message to MQTT broker.

    Args:
        topic: MQTT topic path
        payload: Message payload as dictionary
    """
    await mqtt_publish(topic, payload)
    return f"Published to {topic}"

@Tool
async def get_iot_device_status(device_id: str) -> dict:
    """Get status of an IoT device.

    Args:
        device_id: Device identifier
    """
    clients = await get_clients()
    device = next(
        (c for c in clients["data"] if c["clientid"] == device_id),
        None
    )

    if device:
        return {
            "device_id": device_id,
            "connected": True,
            "ip": device.get("ip_address"),
            "uptime": device.get("connected_at")
        }
    return {"device_id": device_id, "connected": False}

@Tool
async def send_device_command(
    device_id: str,
    command: str,
    params: dict = None
) -> str:
    """Send command to IoT device.

    Args:
        device_id: Target device
        command: Command name
        params: Command parameters
    """
    topic = f"onemind/devices/{device_id}/command"
    payload = {
        "command": command,
        "params": params or {},
        "timestamp": datetime.utcnow().isoformat()
    }

    await mqtt_publish(topic, payload)
    return f"Command '{command}' sent to {device_id}"
```

## Rule Engine

Create data processing rules in EMQX Dashboard:

```sql
-- Forward temperature alerts
SELECT
  payload.temperature as temp,
  payload.device_id as device,
  timestamp
FROM
  "onemind/devices/+/telemetry"
WHERE
  payload.temperature > 30

-- Action: Forward to alert topic
```

## Authentication

### Username/Password

```yaml
# emqx.conf
authentication = [
{
mechanism = password_based
backend = built_in_database
}
]
```

### JWT

```yaml
authentication = [
{
mechanism = jwt
from = password
secret = "your-jwt-secret"
verify_claims = {
username = "${username}"
}
}
]
```

## Environment Variables

```bash
EMQX_HOST=localhost
EMQX_PORT=1883
EMQX_USERNAME=onemind
EMQX_PASSWORD=your_password
EMQX_API_KEY=your_api_key
EMQX_API_SECRET=your_api_secret
EMQX_DASHBOARD__DEFAULT_PASSWORD=admin_password
```

## Best Practices

1. **Topic hierarchy** - Use clear, consistent naming
2. **QoS levels** - 0 for telemetry, 1 for commands, 2 for critical
3. **Retained messages** - For status/state topics
4. **Client IDs** - Use unique, meaningful identifiers
5. **ACLs** - Restrict topic access per client

## EMQX vs Mosquitto

| Aspect      | EMQX                 | Mosquitto         |
| ----------- | -------------------- | ----------------- |
| Scale       | Millions connections | Thousands         |
| Clustering  | Built-in             | Manual            |
| Dashboard   | Rich web UI          | None              |
| Rule Engine | Yes                  | No                |
| Best for    | Production IoT       | Development/small |

## Resources

- [EMQX Docs](https://www.emqx.io/docs/en/latest/)
- [EMQX GitHub](https://github.com/emqx/emqx)
- [MQTT 5.0 Spec](https://docs.oasis-open.org/mqtt/mqtt/v5.0/mqtt-v5.0.html)
- [Rule Engine](https://www.emqx.io/docs/en/latest/data-integration/rules.html)

---

Last Updated: January 25, 2026
