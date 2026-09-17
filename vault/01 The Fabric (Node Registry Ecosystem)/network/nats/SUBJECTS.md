# NATS Message Bus — The Nervous System

> Sub-millisecond message delivery between Agno (brain) and OM1 (body). All inter-layer communication flows through NATS subjects.

## Purpose

NATS is the decoupled communication backbone of One Mind OS. Agno and OM1 never call each other directly. They publish messages to NATS subjects and subscribe to subjects they care about. This ensures crash isolation (if one layer dies, the other continues) and enables multi-robot scaling (add a robot = add a subscriber).

## Subject Schema

All subjects follow the pattern: `onemind.{category}.{identifier}`

### Agno → OM1 (Commands)

```
onemind.commands.{robot_id}     # High-level commands to a specific robot
onemind.commands.fleet          # Broadcast command to all robots
onemind.personality             # Personality/identity sync
```

**Command payload:**
```json
{
  "action": "patrol",
  "params": "living room, check all entry points",
  "priority": "normal",
  "source": "agno-team-onemind",
  "timestamp": 1708300000.0
}
```

**Supported actions:** `walk`, `run`, `sit`, `stand`, `wave`, `patrol`, `go_to`, `speak`, `charge`, `stop`, `dance`, `shake_paw`, `backflip`

### OM1 → Agno (Sensor Data)

```
onemind.vision.{robot_id}       # What the robot sees (VLM text descriptions)
onemind.audio.{robot_id}        # What the robot hears (ASR transcripts)
onemind.status.{robot_id}       # Robot state (battery, position, current action)
onemind.alerts.{robot_id}       # High-priority events (unknown person, danger)
```

**Vision payload:**
```json
{
  "robot_id": "robot1",
  "description": "Adult male standing near front door, holding a package",
  "confidence": 0.92,
  "timestamp": 1708300000.0
}
```

**Status payload:**
```json
{
  "robot_id": "robot1",
  "battery": 78,
  "state": "patrolling",
  "position": { "room": "living_room", "x": 2.3, "y": 4.1 },
  "uptime_seconds": 3600,
  "timestamp": 1708300000.0
}
```

**Alert payload:**
```json
{
  "robot_id": "robot1",
  "alert": "Unknown person detected at front door",
  "priority": "high",
  "vision_context": "Adult male, unknown face, carrying backpack",
  "timestamp": 1708300000.0
}
```

### Aggregated Subjects

```
onemind.fleet.worldstate        # Combined world state from all robots
onemind.fleet.status            # Fleet-wide status summary
```

**World state payload (built by context aggregator):**
```json
{
  "rooms": {
    "living_room": { "status": "person detected", "source": "robot1" },
    "hallway": { "status": "clear", "source": "robot2" },
    "front_door": { "status": "package waiting", "source": "robot3" }
  },
  "alerts": ["unknown person at front door"],
  "fleet": [
    { "id": "robot1", "battery": 78, "state": "patrolling" },
    { "id": "robot2", "battery": 30, "state": "charging" },
    { "id": "robot3", "battery": 95, "state": "idle" }
  ],
  "timestamp": 1708300000.0
}
```

## Context Aggregator Service

A lightweight Python service that subscribes to all OM1 sensor feeds and builds a unified world state. This state is consumed by Agno's heartbeat and reactive agents.

```python
import asyncio, nats, json, time

class ContextAggregator:
    def __init__(self, nats_url="nats://localhost:4222"):
        self.nats_url = nats_url
        self.world_state = {
            "rooms": {},
            "alerts": [],
            "fleet": {},
        }

    async def start(self):
        self._nc = await nats.connect(self.nats_url)
        await self._nc.subscribe("onemind.vision.*", cb=self._on_vision)
        await self._nc.subscribe("onemind.status.*", cb=self._on_status)
        await self._nc.subscribe("onemind.alerts.*", cb=self._on_alert)

        # Respond to world state queries
        await self._nc.subscribe(
            "onemind.fleet.worldstate",
            cb=self._on_query,
        )

        # Publish aggregated state every 10 seconds
        asyncio.create_task(self._publish_loop())

    async def _on_vision(self, msg):
        data = json.loads(msg.data.decode())
        robot_id = data["robot_id"]
        self.world_state["rooms"][data.get("room", "unknown")] = {
            "status": data["description"],
            "source": robot_id,
            "timestamp": data["timestamp"],
        }

    async def _on_status(self, msg):
        data = json.loads(msg.data.decode())
        self.world_state["fleet"][data["robot_id"]] = {
            "battery": data["battery"],
            "state": data["state"],
            "position": data.get("position", {}),
        }

    async def _on_alert(self, msg):
        data = json.loads(msg.data.decode())
        self.world_state["alerts"].append({
            "alert": data["alert"],
            "robot_id": data["robot_id"],
            "timestamp": data["timestamp"],
        })
        # Keep only last 50 alerts
        self.world_state["alerts"] = self.world_state["alerts"][-50:]

    async def _on_query(self, msg):
        await msg.respond(json.dumps(self.world_state).encode())

    async def _publish_loop(self):
        while True:
            await self._nc.publish(
                "onemind.fleet.worldstate",
                json.dumps(self.world_state).encode(),
            )
            await asyncio.sleep(10)
```

## NATS Server Setup

```bash
# Install NATS
curl -L https://github.com/nats-io/nats-server/releases/download/v2.10.x/nats-server-v2.10.x-linux-amd64.tar.gz | tar xz

# Run with JetStream (for message persistence)
nats-server --jetstream --store_dir /data/nats

# Or via Docker
docker run -d --name nats -p 4222:4222 -p 8222:8222 nats:latest --jetstream
```

**Config for production (`nats.conf`):**
```
port: 4222
http_port: 8222

jetstream {
  store_dir: /data/nats
  max_mem: 1G
  max_file: 10G
}

# Optional: authentication
authorization {
  token: "onemind-secret-token"
}
```

## Latency Budget

| Hop | Latency |
|-----|---------|
| NATS publish → subscribe | < 1ms |
| Agno agent instantiation | ~3μs |
| LLM #1 inference (Agno) | 1-5s |
| LLM #2 inference (OM1) | 200-500ms |
| OM1 fuser cycle | 500ms-2s |
| End-to-end (command → physical action) | 2-5s |
| OM1 safety reflex (bypass Agno) | < 300ms |


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
