# Shared Vision — Multi-Robot Awareness

> Each OM1 publishes what it sees as text. Agno subscribes to all robots and builds a unified world state. No raw video streaming needed.

## How It Works

Each OM1 instance has a VLM (Vision Language Model) captioning layer that converts camera frames into natural language descriptions. These descriptions are published to NATS. Agno's Context Aggregator subscribes to all robots' feeds and builds a unified world map.

```
Robot 1 camera → VLM → "person in living room"  → NATS → onemind.vision.robot1
Robot 2 camera → VLM → "empty hallway"           → NATS → onemind.vision.robot2
Robot 3 camera → VLM → "package at front door"   → NATS → onemind.vision.robot3
                                                          │
                                                          ▼
                                                   Context Aggregator
                                                          │
                                                          ▼
                                                   Unified World State
                                                   {
                                                     "living_room": "person detected",
                                                     "hallway": "clear",
                                                     "front_door": "package waiting"
                                                   }
```

## Why Text, Not Video

OM1's NLDB architecture already converts video to text. This means:

- **Bandwidth**: Text descriptions are ~100 bytes vs. video frames at ~100KB each
- **Processing**: Agno doesn't need GPU for video processing — it reads sentences
- **Storage**: World state history is a JSON log, not a video archive
- **LLM-friendly**: LLM #1 reasons over text natively

## Robot-to-Robot Awareness

Each OM1 instance can also subscribe to other robots' vision feeds via NATS. This enables coordinated behavior without centralized planning:

```python
# In each OM1's NATS input plugin
async def start(self):
    self._nc = await nats.connect(self.nats_url)
    # Subscribe to commands from Agno
    await self._nc.subscribe(f"onemind.commands.{self.robot_id}", cb=self._on_command)
    # Subscribe to OTHER robots' vision (fleet awareness)
    await self._nc.subscribe("onemind.vision.*", cb=self._on_fleet_vision)

async def _on_fleet_vision(self, msg):
    """Other robots' observations become context for this robot."""
    data = json.loads(msg.data.decode())
    if data["robot_id"] != self.robot_id:  # skip own data
        self.fleet_context[data["robot_id"]] = data["description"]

def get_input(self) -> str:
    """Inject fleet awareness into this robot's fuser."""
    lines = []
    for rid, desc in self.fleet_context.items():
        lines.append(f"[{rid} sees]: {desc}")
    if self.latest_directive:
        lines.append(f"[Directive]: {self.latest_directive}")
        self.latest_directive = ""
    return "\n".join(lines)
```

Result: Robot 1 knows "Robot 2 is already covering the hallway" and focuses on a different area.

## When You Need Actual Video

For security camera feeds, live monitoring, or AR glasses overlay, use a separate media stream alongside NATS:

- **WebRTC** — for real-time browser-viewable streams
- **RTSP** — for security camera integration (Flock Safety, Frigate)
- **Bluetooth** — for Brilliant Labs Frame AR glasses

These video streams go directly to the Flutter app or web UI for display. Agno doesn't process them — it only reads the text descriptions from NATS.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
