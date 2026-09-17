# NATS

> **Message Broker & Event Streaming**

## Overview

NATS is a high-performance message broker for event-driven communication between OneMind OS services, enabling real-time updates and loose coupling.

## Status

- **Active**: ✅ Yes
- **Deployment**: Docker
- **JetStream**: Enabled
- **Last Updated**: January 25, 2026

## Use Cases

1. **Event Streaming** - Real-time event flow
2. **Service Communication** - Microservice messaging
3. **Agent Events** - Agent activity broadcasting
4. **Notifications** - Real-time UI updates
5. **Task Queues** - Async job processing

## Key Features

- Ultra-fast (millions of msgs/sec)
- At-most-once & at-least-once delivery
- JetStream for persistence
- Subject-based routing
- Wildcard subscriptions
- Request-reply patterns

## Installation

### Docker

```bash
docker run -d \
  --name nats \
  -p 4222:4222 \
  -p 8222:8222 \
  nats:latest \
  -js -m 8222
```

### Docker Compose

```yaml
services:
  nats:
    image: nats:latest
    command: "-js -m 8222"
    ports:
      - "4222:4222" # Client
      - "8222:8222" # Monitoring
    volumes:
      - nats-data:/data
    restart: unless-stopped

volumes:
  nats-data:
```

## Python Client

### Installation

```bash
pip install nats-py
```

### Basic Connection

```python
import nats
import asyncio

async def main():
    nc = await nats.connect("nats://localhost:4222")

    # Publish
    await nc.publish("events.user.login", b'{"user_id": "123"}')

    # Subscribe
    async def handler(msg):
        print(f"Received: {msg.data.decode()}")

    sub = await nc.subscribe("events.*", cb=handler)

    # Keep running
    await asyncio.sleep(3600)

    await nc.drain()

asyncio.run(main())
```

## Subject Patterns

### Hierarchy

```
events.agent.started
events.agent.completed
events.agent.error
events.user.login
events.user.logout
events.home.light_on
events.home.temperature_changed
```

### Wildcards

```python
# Single token wildcard (*)
await nc.subscribe("events.agent.*")  # Matches events.agent.started

# Multi-token wildcard (>)
await nc.subscribe("events.>")  # Matches all events

# Specific patterns
await nc.subscribe("events.*.error")  # All error events
```

## JetStream (Persistence)

### Create Stream

```python
js = nc.jetstream()

# Create stream
await js.add_stream(
    name="EVENTS",
    subjects=["events.>"],
    retention="limits",
    max_msgs=1_000_000,
    max_age=86400 * 7  # 7 days
)
```

### Publish with Ack

```python
ack = await js.publish("events.agent.started", b'{"agent": "legacy"}')
print(f"Published to stream: {ack.stream}, seq: {ack.seq}")
```

### Durable Consumer

```python
# Create durable consumer
consumer = await js.pull_subscribe(
    "events.agent.>",
    durable="agent-processor"
)

# Fetch messages
msgs = await consumer.fetch(10, timeout=5)
for msg in msgs:
    print(f"Processing: {msg.data}")
    await msg.ack()
```

### Push Consumer

```python
async def handler(msg):
    print(f"Received: {msg.data.decode()}")
    await msg.ack()

await js.subscribe(
    "events.>",
    cb=handler,
    durable="event-handler",
    deliver_policy=nats.js.api.DeliverPolicy.NEW
)
```

## Event Patterns

### Agent Events

```python
import json
from datetime import datetime

async def publish_agent_event(nc, event_type: str, data: dict):
    subject = f"events.agent.{event_type}"
    payload = {
        "type": event_type,
        "timestamp": datetime.utcnow().isoformat(),
        "data": data
    }
    await nc.publish(subject, json.dumps(payload).encode())

# Usage
await publish_agent_event(nc, "run_started", {
    "agent_id": "legacy",
    "session_id": "sess-123",
    "input": "What's on my calendar?"
})

await publish_agent_event(nc, "run_completed", {
    "agent_id": "legacy",
    "session_id": "sess-123",
    "output": "You have 3 meetings today.",
    "duration_ms": 1234
})
```

### Home Events

```python
async def publish_home_event(nc, event_type: str, entity: str, data: dict):
    subject = f"events.home.{event_type}"
    payload = {
        "entity": entity,
        "timestamp": datetime.utcnow().isoformat(),
        **data
    }
    await nc.publish(subject, json.dumps(payload).encode())

# Usage
await publish_home_event(nc, "light_on", "light.living_room", {
    "brightness": 255
})
```

## Request-Reply

```python
# Responder
async def handler(msg):
    request = json.loads(msg.data.decode())
    response = await process_request(request)
    await nc.publish(msg.reply, json.dumps(response).encode())

await nc.subscribe("requests.agent.run", cb=handler)

# Requester
response = await nc.request(
    "requests.agent.run",
    json.dumps({"message": "Hello"}).encode(),
    timeout=30
)
print(f"Response: {response.data.decode()}")
```

## FastAPI Integration

```python
from fastapi import FastAPI
from contextlib import asynccontextmanager

app = FastAPI()
nc = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    global nc
    nc = await nats.connect("nats://localhost:4222")
    yield
    await nc.drain()

app = FastAPI(lifespan=lifespan)

@app.post("/events")
async def publish_event(event: dict):
    subject = f"events.{event['type']}"
    await nc.publish(subject, json.dumps(event).encode())
    return {"status": "published"}
```

## Monitoring

### Web UI

Access at `http://localhost:8222` for monitoring dashboard.

### Key Metrics

```bash
# Server info
curl http://localhost:8222/varz

# Connections
curl http://localhost:8222/connz

# Subscriptions
curl http://localhost:8222/subsz

# JetStream info
curl http://localhost:8222/jsz
```

### Prometheus Metrics

```yaml
# prometheus.yml
scrape_configs:
  - job_name: "nats"
    static_configs:
      - targets: ["nats:8222"]
```

## Configuration

### nats.conf

```conf
port: 4222
monitor_port: 8222

jetstream {
    store_dir: "/data"
    max_memory_store: 1G
    max_file_store: 10G
}

authorization {
    user: onemind
    password: $NATS_PASSWORD
}

cluster {
    name: onemind-cluster
    listen: 0.0.0.0:6222
    routes: [
        nats://nats-1:6222
        nats://nats-2:6222
    ]
}
```

## Event Schema

```python
from pydantic import BaseModel
from datetime import datetime
from typing import Any

class Event(BaseModel):
    id: str
    type: str
    timestamp: datetime
    source: str
    data: dict[str, Any]

    def to_nats(self) -> bytes:
        return self.model_dump_json().encode()

    @classmethod
    def from_nats(cls, data: bytes) -> "Event":
        return cls.model_validate_json(data)
```

## Best Practices

- ✅ Use descriptive subject hierarchies
- ✅ Enable JetStream for important events
- ✅ Set appropriate retention policies
- ✅ Use durable consumers for reliability
- ✅ Acknowledge messages promptly
- ✅ Monitor consumer lag
- ✅ Use connection pooling in high-throughput apps

## Resources

- [NATS Documentation](https://docs.nats.io/)
- [NATS Python Client](https://github.com/nats-io/nats.py)
- [JetStream Guide](https://docs.nats.io/nats-concepts/jetstream)

---

**Related**: [Redis](./Redis.md) | [Agno](./Agno.md) | [Home-Assistant](./Home-Assistant.md)
