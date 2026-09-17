# OM1 Embodiment Layer — The Body

> Handles all physical perception, motor control, real-time safety, and hardware abstraction. This is Legacy's physical presence in the world.

## Role in the Architecture

OM1 is the embodiment runtime for One Mind OS. It runs a continuous perception-action loop at 0.5-2 Hz, processing sensor data through its Natural Language Data Bus (NLDB), making real-time physical decisions via LLM #2, and executing motor commands through the Hardware Abstraction Layer (HAL).

OM1 does NOT handle digital tasks (email, calendar, databases). Those belong to Agno. OM1's job is to see, hear, move, and react fast.

## Core Loop

Every 0.5-2 seconds, OM1 executes this pipeline:

```
Sensors → AI Captioning → NLDB → Fuser → Cortex LLM → Action Orchestrator → HAL → Robot
```

1. **Sensors** — Camera, microphone, LiDAR, IMU capture raw data
2. **AI Captioning** — VLM converts images to text, ASR converts audio to text
3. **NLDB** — All sensor data is now natural language on the data bus
4. **Fuser** — Combines all text inputs (sensors + any Agno commands from NATS) into one prompt
5. **Cortex LLM** — LLM #2 (fast, local) decides what action to take
6. **Action Orchestrator** — Parses structured output into specific commands
7. **HAL** — Translates commands to motor signals for the specific robot platform

## LLM #2 — The Reflex Model

OM1's cortex uses a fast, lightweight model optimized for physical decision-making.

- **Primary**: Qwen3-30B running locally on-device
- **Fallback**: GPT-4.1-nano via cloud (low latency)
- **Response time**: 200-500ms (critical for physical safety)
- **Output format**: Structured `CortexOutputModel` (Pydantic schema)

OM1 supports dual-model inference: local and cloud run in parallel, first good response wins. This maximizes both speed and quality.

See [LLM-STRATEGY.md](./LLM-STRATEGY.md) for full LLM architecture details.

## NLDB — Natural Language Data Bus

The defining architectural feature of OM1. ALL sensor data is converted to natural language before reaching the LLM:

- Camera feed → "adult male standing near couch, holding phone"
- LiDAR scan → "obstacle 1.2 meters ahead, wall to the right"
- Microphone → "doorbell sound detected from front of house"
- IMU → "robot tilted 15 degrees, possible uneven surface"

This makes OM1 hardware-agnostic. The LLM doesn't care if data comes from a $50 webcam or a $5,000 depth sensor — it's all text.

## Integration with Agno via NATS

OM1 connects to Agno through two custom plugins:

### Custom Input Plugin — Receives Agno Commands

```python
# src/inputs/agno_nats_input.py
import asyncio, nats, json

class AgnoNATSInput:
    """OM1 input plugin: receives high-level commands from Agno."""

    def __init__(self, config: dict):
        self.nats_url = config.get("nats_url", "nats://localhost:4222")
        self.robot_id = config.get("robot_id", "robot1")
        self.latest_directive = ""
        self._nc = None

    async def start(self):
        self._nc = await nats.connect(self.nats_url)
        await self._nc.subscribe(
            f"onemind.commands.{self.robot_id}", cb=self._on_command
        )
        await self._nc.subscribe(
            "onemind.personality", cb=self._on_personality
        )

    async def _on_command(self, msg):
        data = json.loads(msg.data.decode())
        self.latest_directive = data.get("action", "")
        if data.get("params"):
            self.latest_directive += f" — {data['params']}"

    async def _on_personality(self, msg):
        """Update the robot's personality from Agno's source of truth."""
        data = json.loads(msg.data.decode())
        # Inject into OM1's system_prompt_base and governance
        self.personality_update = data

    def get_input(self) -> str:
        """Called by OM1 fuser each cognitive cycle."""
        if self.latest_directive:
            directive = self.latest_directive
            self.latest_directive = ""  # consume it
            return f"[Cognitive Layer Directive]: {directive}"
        return ""
```

### Custom Action Plugin — Publishes Sensor Data to NATS

```python
# src/actions/nats_publisher.py
import nats, json

class NATSPublisherAction:
    """OM1 action plugin: publishes sensor summaries and events to NATS."""

    def __init__(self, config: dict):
        self.nats_url = config.get("nats_url", "nats://localhost:4222")
        self.robot_id = config.get("robot_id", "robot1")
        self._nc = None

    async def start(self):
        self._nc = await nats.connect(self.nats_url)

    async def publish_vision(self, description: str):
        await self._nc.publish(
            f"onemind.vision.{self.robot_id}",
            json.dumps({
                "robot_id": self.robot_id,
                "description": description,
                "timestamp": time.time(),
            }).encode()
        )

    async def publish_status(self, status: dict):
        await self._nc.publish(
            f"onemind.status.{self.robot_id}",
            json.dumps({
                "robot_id": self.robot_id,
                "battery": status.get("battery", 100),
                "state": status.get("state", "idle"),
                "position": status.get("position", {}),
                "timestamp": time.time(),
            }).encode()
        )

    async def publish_alert(self, alert: str):
        """High-priority events — unknown person, loud noise, etc."""
        await self._nc.publish(
            f"onemind.alerts.{self.robot_id}",
            json.dumps({
                "robot_id": self.robot_id,
                "alert": alert,
                "priority": "high",
                "timestamp": time.time(),
            }).encode()
        )
```

### OM1 Configuration (JSON5)

```json5
{
  "robot_id": "robot1",
  "llm": {
    "model": "qwen3-30b",
    "provider": "ollama",
    "fallback": {
      "model": "gpt-4.1-nano",
      "provider": "openai"
    }
  },
  "agent_inputs": [
    { "type": "AgnoNATSInput", "config": { "nats_url": "nats://localhost:4222", "robot_id": "robot1" } },
    { "type": "VLM_COCO_Local", "config": { "camera_index": 0 } },
    { "type": "GoogleASRInput" },
    { "type": "LidarInput", "config": { "enabled": true } }
  ],
  "agent_actions": [
    { "type": "NATSPublisherAction", "config": { "nats_url": "nats://localhost:4222", "robot_id": "robot1" } },
    { "type": "UnitreeG1Move" },
    { "type": "TTSSpeak", "config": { "voice": "alloy" } }
  ],
  "system_prompt_base": "You are Legacy, a unified AI with a physical robot body...",
  "history_length": 10,
  "cortex_rate_hz": 1.0
}
```

## Supported Robot Platforms

| Platform | Type | Connection | Capabilities |
|----------|------|-----------|-------------|
| Unitree G1 | Humanoid | CycloneDDS / Zenoh | Walk, wave, grasp, sit, speak |
| Unitree Go2 | Quadruped | CycloneDDS | Walk, run, sit, shake paw, backflip |
| TurtleBot 4 | Wheeled | ROS2 | Navigate, map, patrol |
| UbTech Yanshee | Humanoid | Custom API | Walk, dance, gesture |

## What OM1 Does NOT Handle

These belong to the Agno layer:

- Email, Slack, calendar, digital communications
- Long-term persistent memory across sessions
- Knowledge base / RAG retrieval
- Multi-agent team coordination
- MCP tool connections
- Complex reasoning and planning
- Webhook processing

OM1 focuses exclusively on physical perception and action.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
