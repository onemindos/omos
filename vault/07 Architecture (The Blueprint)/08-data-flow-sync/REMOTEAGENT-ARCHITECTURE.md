# RemoteAgent + NATS Architecture

> Deep dive into how Legacy controls machines across multiple locations
> Last Updated: 2026-01-21

---

## Table of Contents

1. [Overview](#overview)
2. [Why NATS Instead of Direct Communication](#why-nats-instead-of-direct-communication)
3. [How It Works](#how-it-works)
4. [Message Flow](#message-flow)
5. [NATS Subjects (Addressing)](#nats-subjects-addressing)
6. [Real-World Examples](#real-world-examples)
   - [Morning Farm Check (Voice from Office)](#example-1-morning-farm-check-voice-from-office)
   - [Emergency Alert (Automatic)](#example-2-emergency-alert-automatic)
   - [Home Management from Office](#example-3-home-management-from-office-voice)
   - [Voice Command While Driving](#example-4-voice-command-while-driving)
7. [Installing RemoteAgent](#installing-remoteagent)
8. [RemoteAgent Code Structure](#remoteagent-code-structure)
9. [Full Architecture Diagram](#full-architecture-diagram)

---

## Overview

```
NATS = Nervous System (messages flow like nerve signals)
RemoteAgent = Legacy's hands/eyes at each location
```

Legacy (the AI brain in the cloud) needs to control devices, read sensors, and run commands on machines at multiple physical locations (home, office, farm, barn, etc.). Instead of connecting directly to each machine, Legacy uses NATS as a message broker - like a nervous system that carries signals to and from every location.

---

## Why NATS Instead of Direct Communication

**Question:** Why doesn't Legacy speak to agents directly?

**Answer:**

```
DIRECT (Legacy -> Agent)              NATS (Legacy -> NATS -> Agent)
---------------------------------------------------------------------------

X Legacy needs to know every          Legacy just publishes to
  agent's IP address                  "commands.farm.sensors"
                                      NATS routes it

X If agent is offline = error         If agent is offline =
                                      message queued, delivered later

X Agent must expose a port            Agent only connects OUT to NATS
  (security risk)                     (no open ports)

X One message = one destination       One message can go to ALL farms
                                      "commands.*.sensors" = broadcast

X Legacy must poll agents             Agents push events TO Legacy
  "any updates?"                      via same NATS

X Different protocols per agent       Everything speaks NATS
  (HTTP, SSH, MQTT, etc.)
```

**The real magic - events flow BACK without asking:**

```
WITHOUT NATS:                        WITH NATS:
---------------------------------------------------------------------------

Legacy: "Farm, any alerts?"          Farm Agent detects frost
Farm: "No"                                    |
Legacy: "Farm, any alerts?"                   v
Farm: "No"                           Agent publishes:
Legacy: "Farm, any alerts?"          events.farm.alert.frost
Farm: "YES FROST!"                            |
Legacy: "Oh no, too late"                     v
                                     Legacy receives INSTANTLY
                                     (was subscribed to events.*)
```

**NATS is the nervous system** - signals flow both directions without polling.

---

## How It Works

```
+---------------------------------------------------------------------------------+
|                         REMOTEAGENT + NATS EXPLAINED                            |
+---------------------------------------------------------------------------------+
|                                                                                 |
|  YOU (voice): "Hey Legacy, what's the temperature in the greenhouse?"          |
|                                                                                 |
|                              |                                                  |
|                              v                                                  |
|                                                                                 |
|   +------------------------------------------------------------------------+   |
|   |  1. Legacy (Cloud) receives voice, understands intent                  |   |
|   |                                                                        |   |
|   |  2. Legacy publishes to NATS:                                          |   |
|   |     Subject: "commands.farm.sensors"                                   |   |
|   |     Payload: {"action": "read", "sensor": "greenhouse_temp"}           |   |
|   |                                                                        |   |
|   |  3. NATS routes message to Farm's NATS Leaf                            |   |
|   |                                                                        |   |
|   |  4. RemoteAgent (Farm) receives, reads sensor, responds:               |   |
|   |     {"temperature": 78.5, "humidity": 65, "status": "optimal"}         |   |
|   |                                                                        |   |
|   |  5. Legacy receives response, speaks:                                  |   |
|   |     "The greenhouse is at 78.5 degrees with 65% humidity.              |   |
|   |      Looking good!"                                                    |   |
|   |                                                                        |   |
|   +------------------------------------------------------------------------+   |
|                                                                                 |
|   Total time: ~200ms (NATS is FAST)                                            |
|                                                                                 |
+---------------------------------------------------------------------------------+
```

---

## Message Flow

```
+---------------------------------------------------------------------------------+
|                              NATS MESSAGE FLOW                                  |
+---------------------------------------------------------------------------------+
|                                                                                 |
|   CLOUD VM                                           FARM EDGE                  |
|   --------------------------------------------------------------------------   |
|                                                                                 |
|   +--------------+      +--------------+            +--------------+           |
|   |    Legacy    |      |  NATS Hub    |            |  NATS Leaf   |           |
|   |              |      |              |<==========>|              |           |
|   |  "Turn on    |----->|              |  Tailscale |              |           |
|   |   sprinklers"|  (1) |              |  Encrypted |              |           |
|   |              |      |              |            |              |           |
|   +--------------+      +--------------+            +------+-------+           |
|                                                            |                    |
|   (1) Legacy publishes:                                    | (2)               |
|       nats.publish(                                        v                    |
|         "commands.farm.irrigation",             +------------------+           |
|         {action: "on", zone: "garden"}          |   RemoteAgent    |           |
|       )                                         |                  |           |
|                                                 | Subscribed to:   |           |
|                                                 | commands.farm.>  |           |
|   (4) Legacy receives response:                 |                  |           |
|       "Sprinklers on, zone: garden,             | (3) Executes:    |           |
|        duration: 30 min"                        | GPIO.on(pin=17)  |           |
|                                                 | Returns result   |           |
|                                                 +------------------+           |
|                                                                                 |
+---------------------------------------------------------------------------------+
```

---

## NATS Subjects (Addressing)

NATS subjects are like addresses. Legacy publishes to a subject, and any agent subscribed to that subject receives the message.

### Command Subjects (Legacy -> Agents)

| Subject | Purpose |
|---------|---------|
| `commands.farm.irrigation` | Water control at farm |
| `commands.farm.sensors` | Read temp/humidity/soil at farm |
| `commands.farm.cameras` | View/snapshot cameras at farm |
| `commands.farm.shell` | Run any shell command at farm |
| `commands.farm.coop` | Chicken coop controls |
| `commands.farm.greenhouse` | Greenhouse controls |
| `commands.home.lights` | Control lights at home |
| `commands.home.hvac` | Thermostat control |
| `commands.home.audio` | Play through speakers |
| `commands.home.locks` | Door lock control |
| `commands.office.printer` | Print documents |
| `commands.office.shell` | Run commands on office PC |
| `commands.*.sensors` | Broadcast to ALL locations |

### Event Subjects (Agents -> Legacy)

| Subject | Purpose |
|---------|---------|
| `events.farm.motion` | Motion detected (camera) |
| `events.farm.alert.critical` | Critical alert (frost, flood) |
| `events.farm.coop.door` | Coop door opened/closed |
| `events.home.doorbell` | Someone at door |
| `events.home.motion` | Motion detected |
| `events.home.alarm` | Alarm triggered |
| `agents.online` | Agent came online |
| `agents.offline` | Agent went offline |

---

## Real-World Examples

### Example 1: Morning Farm Check (Voice from Office)

```
+---------------------------------------------------------------------------------+
|  SCENARIO: You're at office, want to check on farm                              |
+---------------------------------------------------------------------------------+
|                                                                                 |
|  YOU (voice): "Legacy, give me a farm status report"                           |
|                                                                                 |
|  Legacy thinks:                                                                 |
|  +--------------------------------------------------------------------------+  |
|  | 1. Need greenhouse temp      -> commands.farm.sensors.greenhouse         |  |
|  | 2. Need chicken coop status  -> commands.farm.sensors.coop               |  |
|  | 3. Need water tank level     -> commands.farm.sensors.water              |  |
|  | 4. Need camera snapshot      -> commands.farm.cameras.barn               |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  Legacy sends 4 parallel NATS requests (all at once, ~50ms)                    |
|                                                                                 |
|  RemoteAgent (Farm) responds:                                                  |
|  +--------------------------------------------------------------------------+  |
|  | greenhouse: {temp: 72F, humidity: 70%, door: "closed"}                   |  |
|  | coop: {temp: 65F, door: "open", eggs_collected: false}                   |  |
|  | water: {level: 85%, pump: "off", last_filled: "2 days ago"}              |  |
|  | camera: {image_url: "https://...", motion_last_24h: 3}                   |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  Legacy (voice): "Good morning! Here's your farm status:                       |
|                                                                                 |
|                   The greenhouse is at 72 degrees, nice and humid.             |
|                   Chicken coop door is open, but eggs haven't been             |
|                   collected yet.                                               |
|                   Water tank is at 85% - you're good for a few more days.      |
|                   I saw 3 motion events at the barn in the last 24 hours,      |
|                   probably just the cats.                                      |
|                                                                                 |
|                   Want me to show you the barn camera?"                        |
|                                                                                 |
|  YOU: "Yes, and remind me to collect eggs when I get home"                     |
|                                                                                 |
|  Legacy: *shows camera on your screen*                                         |
|          *creates reminder for when you arrive at farm location*               |
|                                                                                 |
+---------------------------------------------------------------------------------+
```

### Example 2: Emergency Alert (Automatic)

```
+---------------------------------------------------------------------------------+
|  SCENARIO: Greenhouse temperature drops dangerously (you're asleep)            |
+---------------------------------------------------------------------------------+
|                                                                                 |
|  3:00 AM - RemoteAgent (Farm) detects:                                         |
|  +--------------------------------------------------------------------------+  |
|  |  Greenhouse temp: 38F (was 65F, dropping fast)                           |  |
|  |  Outside temp: 28F (unexpected frost)                                    |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  RemoteAgent publishes EVENT (not waiting for command):                        |
|  +--------------------------------------------------------------------------+  |
|  |  Subject: events.farm.alert.critical                                     |  |
|  |  Payload: {                                                              |  |
|  |    type: "temperature_drop",                                             |  |
|  |    location: "greenhouse",                                               |  |
|  |    current: 38,                                                          |  |
|  |    threshold: 45,                                                        |  |
|  |    trend: "dropping",                                                    |  |
|  |    risk: "frost_damage_imminent"                                         |  |
|  |  }                                                                       |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  Legacy (Cloud) receives, evaluates:                                           |
|  - Awareness mode = DORMANT but this is CRITICAL                               |
|  - Override dormant mode for critical alerts                                   |
|                                                                                 |
|  Legacy decides:                                                               |
|  +--------------------------------------------------------------------------+  |
|  |  1. Turn on greenhouse heater    -> commands.farm.greenhouse.heater ON   |  |
|  |  2. Close greenhouse vents       -> commands.farm.greenhouse.vents CLOSE |  |
|  |  3. Alert owner (SMS - urgent)   -> Direct Twilio API                    |  |
|  |  4. Log event                    -> PostgreSQL                           |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  YOU receive SMS: "Legacy: Greenhouse dropped to 38F. I've turned on the       |
|                    heater and closed vents. Current status: heating.           |
|                    Reply CALL for voice."                                      |
|                                                                                 |
|  YOU reply: "CALL"                                                             |
|                                                                                 |
|  Legacy calls you via Twilio, speaks:                                          |
|  "Hey, sorry to wake you. We had an unexpected frost. The greenhouse           |
|   dropped to 38 degrees but I caught it in time. Heater's been on for          |
|   5 minutes and we're already back up to 45 degrees. Your tomatoes should      |
|   be fine. Want me to keep monitoring and text you updates, or should I        |
|   call again if it drops?"                                                     |
|                                                                                 |
|  YOU: "Text updates are fine, thanks Legacy"                                   |
|                                                                                 |
|  Legacy: "Got it. Go back to sleep, I've got this."                            |
|                                                                                 |
+---------------------------------------------------------------------------------+
```

### Example 3: Home Management from Office (Voice)

```
+---------------------------------------------------------------------------------+
|  SCENARIO: You're at office, guests arriving at home in 1 hour                 |
+---------------------------------------------------------------------------------+
|                                                                                 |
|  YOU (voice): "Legacy, I have guests coming to the house in an hour.           |
|                Can you get things ready?"                                       |
|                                                                                 |
|  Legacy: "Sure! Let me check on the house first."                              |
|                                                                                 |
|  Legacy sends to home RemoteAgent:                                             |
|  +--------------------------------------------------------------------------+  |
|  |  commands.home.status.full                                               |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  RemoteAgent (Home) responds:                                                  |
|  +--------------------------------------------------------------------------+  |
|  |  thermostat: 68F (eco mode)                                              |  |
|  |  lights: all off                                                         |  |
|  |  robot_vacuum: docked, last_run: 3 days ago                              |  |
|  |  doorbell_camera: online                                                 |  |
|  |  music: off                                                              |  |
|  |  garage: closed                                                          |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  Legacy: "The house is at 68 degrees in eco mode, lights are off, and the      |
|           vacuum hasn't run in 3 days. Here's what I can do:                   |
|                                                                                 |
|           1. Bump the thermostat to 72 for comfort                             |
|           2. Run the vacuum now so it's done before they arrive                |
|           3. Turn on the porch lights and living room                          |
|           4. Start some background music when they're 10 minutes out           |
|                                                                                 |
|           Want me to do all of that?"                                          |
|                                                                                 |
|  YOU: "Yes, but skip the vacuum - I'll do that this weekend. Oh and can you   |
|        unlock the front door when they get there? It's my sister."             |
|                                                                                 |
|  Legacy: "Got it. I'll:                                                        |
|           - Set thermostat to 72 now                                           |
|           - Turn on porch and living room lights                               |
|           - Watch for your sister on the doorbell camera                       |
|           - When I see her, unlock the front door and start music              |
|                                                                                 |
|           Anything else?"                                                      |
|                                                                                 |
|  YOU: "Perfect, thanks!"                                                       |
|                                                                                 |
|  Legacy executes via NATS:                                                     |
|  +--------------------------------------------------------------------------+  |
|  |  commands.home.hvac       -> {temp: 72, mode: "heat"}                    |  |
|  |  commands.home.lights     -> {porch: "on", living_room: "on",            |  |
|  |                               brightness: 80}                            |  |
|  |  commands.home.doorbell   -> {watch_for: "person",                       |  |
|  |                               action: "notify_legacy"}                   |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  [45 minutes later - Frigate detects person at front door]                     |
|                                                                                 |
|  Legacy receives: events.home.doorbell.person_detected                         |
|  Legacy executes:                                                              |
|  +--------------------------------------------------------------------------+  |
|  |  commands.home.lock.front_door -> {action: "unlock"}                     |  |
|  |  commands.home.audio           -> {play: "jazz_playlist", volume: 30}    |  |
|  |  commands.home.lights          -> {entryway: "on"}                       |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  Legacy texts you: "Your sister just arrived. Door's unlocked, music's on."   |
|                                                                                 |
+---------------------------------------------------------------------------------+
```

### Example 4: Voice Command While Driving

```
+---------------------------------------------------------------------------------+
|  SCENARIO: Driving home from office, hands-free                                |
+---------------------------------------------------------------------------------+
|                                                                                 |
|  YOU (CarPlay/voice): "Hey Legacy"                                             |
|                                                                                 |
|  Legacy: "Hey, what's up?"                                                     |
|                                                                                 |
|  YOU: "I'm about 20 minutes out. Can you start the hot tub and preheat         |
|        the oven to 400? Oh and what's the chicken coop situation?"             |
|                                                                                 |
|  Legacy:                                                                        |
|  +--------------------------------------------------------------------------+  |
|  |  Parallel NATS requests:                                                 |  |
|  |  commands.home.hot_tub    -> {action: "on", target_temp: 102}            |  |
|  |  commands.home.oven       -> {action: "preheat", temp: 400}              |  |
|  |  commands.farm.coop       -> {action: "status"}                          |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  Legacy: "Hot tub's heating up - it was at 95, should be ready when you        |
|           get home. Oven's preheating now.                                     |
|                                                                                 |
|           As for the chickens - coop door is still open and it's getting       |
|           dark. Want me to close it? Or are you stopping by the farm first?"  |
|                                                                                 |
|  YOU: "Close it please, I'm going straight home"                               |
|                                                                                 |
|  Legacy:                                                                        |
|  +--------------------------------------------------------------------------+  |
|  |  commands.farm.coop.door  -> {action: "close"}                           |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  Legacy: "Done. Coop door is closed. I counted 6 chickens inside on the        |
|           camera before closing - that's all of them. Drive safe!"             |
|                                                                                 |
+---------------------------------------------------------------------------------+
```

---

## Installing RemoteAgent

You can install RemoteAgent on any machine via SSH - including having Legacy do it for you!

### Option 1: One-liner Install Script

```bash
# From your Mac (or Legacy could do this!)
ssh farm-server "curl -fsSL https://onemindos.dev/install-agent.sh | bash"

# Or with specific location name
ssh farm-server "curl -fsSL https://onemindos.dev/install-agent.sh | \
                 LOCATION=farm NATS_URL=nats://100.102.21.44:4222 bash"
```

### Option 2: Docker One-liner

```bash
ssh farm-server "docker run -d \
  --name onemind-agent \
  --restart unless-stopped \
  --privileged \
  -e LOCATION=farm \
  -e NATS_URL=nats://100.102.21.44:4222 \
  -e NATS_TOKEN=your_token \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/yourusername/onemind-agent:latest"
```

### Option 3: Legacy Installs It For You!

```
+---------------------------------------------------------------------------------+
|  SCENARIO: Legacy installs itself on a new device                              |
+---------------------------------------------------------------------------------+
|                                                                                 |
|  YOU: "Legacy, can you install yourself on the new Raspberry Pi at the barn?"  |
|                                                                                 |
|  Legacy: "Sure! I can see a new device on Tailscale at 100.87.87.50.          |
|           Let me SSH in and set up a RemoteAgent there."                       |
|                                                                                 |
|  Legacy executes:                                                              |
|  +--------------------------------------------------------------------------+  |
|  |  ssh_exec("100.87.87.50", """                                            |  |
|  |    curl -fsSL https://onemindos.dev/install-agent.sh | \                 |  |
|  |    LOCATION=barn \                                                       |  |
|  |    NATS_URL=nats://100.102.21.44:4222 \                                   |  |
|  |    bash                                                                  |  |
|  |  """)                                                                    |  |
|  +--------------------------------------------------------------------------+  |
|                                                                                 |
|  Legacy: "Done! I've installed a RemoteAgent on the barn Pi. It's now         |
|           connected and I can see these capabilities:                          |
|           - GPIO pins (for sensors/relays)                                     |
|           - Camera (USB webcam detected)                                       |
|           - Shell access                                                       |
|                                                                                 |
|           What should I call this location - 'barn' okay?"                    |
|                                                                                 |
+---------------------------------------------------------------------------------+
```

### Install Script (install-agent.sh)

```bash
#!/bin/bash
# install-agent.sh - Install OneMind RemoteAgent

set -e

LOCATION=${LOCATION:-"unknown"}
NATS_URL=${NATS_URL:-"nats://100.102.21.44:4222"}
NATS_TOKEN=${NATS_TOKEN:-""}

echo "Installing OneMind RemoteAgent..."
echo "   Location: $LOCATION"
echo "   NATS: $NATS_URL"

# Detect platform
if command -v docker &> /dev/null; then
    echo "Docker detected, using container install..."

    docker pull ghcr.io/yourusername/onemind-agent:latest

    docker run -d \
        --name onemind-agent \
        --restart unless-stopped \
        --privileged \
        --network host \
        -e LOCATION="$LOCATION" \
        -e NATS_URL="$NATS_URL" \
        -e NATS_TOKEN="$NATS_TOKEN" \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v /sys:/sys \
        -v /proc:/proc \
        ghcr.io/yourusername/onemind-agent:latest

elif command -v python3 &> /dev/null; then
    echo "Python detected, using pip install..."

    pip3 install onemind-agent

    # Create systemd service
    cat > /etc/systemd/system/onemind-agent.service << EOF
[Unit]
Description=OneMind RemoteAgent
After=network.target

[Service]
Environment="LOCATION=$LOCATION"
Environment="NATS_URL=$NATS_URL"
Environment="NATS_TOKEN=$NATS_TOKEN"
ExecStart=/usr/local/bin/onemind-agent
Restart=always

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable onemind-agent
    systemctl start onemind-agent
fi

echo "RemoteAgent installed and running!"
echo "   Test: curl http://localhost:8090/health"
```

---

## RemoteAgent Code Structure

```
agents/edge/
├── remote_agent.py      # Main agent
├── tools/
│   ├── shell.py         # Run commands
│   ├── files.py         # Read/write files
│   ├── docker.py        # Manage containers
│   ├── gpio.py          # Raspberry Pi GPIO
│   ├── sensors.py       # Temperature, humidity, etc.
│   ├── camera.py        # Capture images/video
│   ├── audio.py         # Play sounds, TTS
│   └── system.py        # Reboot, update, etc.
├── config.py            # Location-specific config
└── requirements.txt
```

### Main Agent Code (remote_agent.py)

```python
# agents/edge/remote_agent.py
import asyncio
import nats
import json
import os
from datetime import datetime
from tools import shell, sensors, camera, gpio, audio

LOCATION = os.getenv("LOCATION", "unknown")
NATS_URL = os.getenv("NATS_URL", "nats://localhost:4222")

class RemoteAgent:
    def __init__(self):
        self.nc = None
        self.tools = {
            "shell": shell.execute,
            "sensors": sensors.read,
            "camera": camera.capture,
            "gpio": gpio.control,
            "audio": audio.play,
        }

    async def connect(self):
        self.nc = await nats.connect(NATS_URL)

        # Subscribe to commands for this location
        await self.nc.subscribe(
            f"commands.{LOCATION}.>",
            cb=self.handle_command
        )

        # Announce presence
        await self.nc.publish("agents.online", json.dumps({
            "location": LOCATION,
            "capabilities": list(self.tools.keys()),
            "timestamp": datetime.now().isoformat()
        }))

        print(f"RemoteAgent online at {LOCATION}")

    async def handle_command(self, msg):
        """Handle incoming commands from Legacy"""
        subject = msg.subject  # e.g., "commands.farm.sensors"
        tool_name = subject.split(".")[-1]  # "sensors"

        data = json.loads(msg.data)

        if tool_name in self.tools:
            try:
                result = await self.tools[tool_name](data)
                await msg.respond(json.dumps({
                    "success": True,
                    "data": result
                }))
            except Exception as e:
                await msg.respond(json.dumps({
                    "success": False,
                    "error": str(e)
                }))
        else:
            await msg.respond(json.dumps({
                "success": False,
                "error": f"Unknown tool: {tool_name}"
            }))

    async def publish_event(self, event_type: str, data: dict):
        """Publish an event back to Legacy"""
        await self.nc.publish(
            f"events.{LOCATION}.{event_type}",
            json.dumps({
                "location": LOCATION,
                "timestamp": datetime.now().isoformat(),
                **data
            })
        )

# Run
if __name__ == "__main__":
    agent = RemoteAgent()
    loop = asyncio.get_event_loop()
    loop.run_until_complete(agent.connect())
    loop.run_forever()
```

### Example Tool: Sensors (tools/sensors.py)

```python
# agents/edge/tools/sensors.py
import os

# Try to import GPIO libraries (Raspberry Pi)
try:
    import board
    import adafruit_dht
    HAS_GPIO = True
except ImportError:
    HAS_GPIO = False

async def read(data: dict) -> dict:
    """Read sensor data"""
    sensor_name = data.get("sensor", "all")

    if sensor_name == "all" or sensor_name == "greenhouse_temp":
        if HAS_GPIO:
            dht = adafruit_dht.DHT22(board.D4)
            return {
                "temperature": dht.temperature,
                "humidity": dht.humidity,
                "unit": "celsius"
            }
        else:
            # Fallback: read from system sensors or return mock
            return {
                "temperature": 72.5,
                "humidity": 65,
                "unit": "fahrenheit",
                "source": "mock"
            }

    elif sensor_name == "water_level":
        # Read from water level sensor
        return {
            "level": 85,
            "unit": "percent"
        }

    else:
        raise ValueError(f"Unknown sensor: {sensor_name}")
```

### Example Tool: Shell (tools/shell.py)

```python
# agents/edge/tools/shell.py
import asyncio
import shlex

# Allowlist of safe commands (security)
ALLOWED_COMMANDS = [
    "ls", "cat", "df", "du", "ps", "top", "free",
    "docker", "systemctl", "journalctl",
    "ping", "curl", "wget",
    "date", "uptime", "hostname"
]

async def execute(data: dict) -> dict:
    """Execute a shell command"""
    command = data.get("command", "")

    # Security: check if command starts with allowed prefix
    cmd_parts = shlex.split(command)
    if cmd_parts and cmd_parts[0] not in ALLOWED_COMMANDS:
        # Check if it's a full path to allowed command
        base_cmd = os.path.basename(cmd_parts[0])
        if base_cmd not in ALLOWED_COMMANDS:
            raise PermissionError(
                f"Command not allowed: {cmd_parts[0]}"
            )

    # Execute
    proc = await asyncio.create_subprocess_shell(
        command,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE
    )

    stdout, stderr = await proc.communicate()

    return {
        "stdout": stdout.decode(),
        "stderr": stderr.decode(),
        "return_code": proc.returncode
    }
```

---

## Full Architecture Diagram

```
+---------------------------------------------------------------------------------+
|                              YOUR OMNIPRESENT SETUP                             |
+---------------------------------------------------------------------------------+
|                                                                                 |
|   YOU --------------------------------------------------------------------     |
|    |                                                                            |
|    |  Voice (anywhere)    Browser (Guacamole)    Mobile App                    |
|    |         |                    |                   |                         |
|    +---------+--------------------+-------------------+                         |
|                           |                                                     |
|                           v                                                     |
|   +-----------------------------------------------------------------------------+
|   |                           CLOUD VM                                          |
|   |  +---------+  +---------+  +---------+  +---------+  +---------+  +-------+ |
|   |  | Legacy  |  |  NATS   |  |Guacamole|  |  Nango  |  | Grafana |  |  UI   | |
|   |  | (brain) |  | (nerve) |  | (human) |  | (oauth) |  |(metrics)|  |(face) | |
|   |  +----+----+  +----+----+  +---------+  +---------+  +---------+  +-------+ |
|   |       |            |                                                        |
|   +-------+------------+--------------------------------------------------------+
|           |            |                                                        |
|           |            |  NATS Leaf Connections (Tailscale encrypted)          |
|           |            |                                                        |
|    +------+------------+------------------------+------------------------+      |
|    |                   |                        |                        |      |
|    v                   v                        v                        v      |
|  +-------------+  +-------------+  +-------------+  +-------------+            |
|  |    HOME     |  |   OFFICE    |  |    FARM     |  |    BARN     |            |
|  |             |  |             |  |             |  |  (new Pi!)  |            |
|  | RemoteAgent |  | RemoteAgent |  | RemoteAgent |  | RemoteAgent |            |
|  | - lights    |  | - printer   |  | - sensors   |  | - camera    |            |
|  | - hvac      |  | - scanner   |  | - irrigation|  | - gpio      |            |
|  | - cameras   |  | - shell     |  | - coop door |  | - shell     |            |
|  | - audio     |  |             |  | - cameras   |  |             |            |
|  | - locks     |  |             |  | - greenhouse|  |             |            |
|  +-------------+  +-------------+  +-------------+  +-------------+            |
|                                                                                 |
|   Legacy can:                                                                   |
|   - Run any command on any machine                                              |
|   - Read sensors across all locations                                           |
|   - Control devices (lights, locks, irrigation, etc.)                           |
|   - Capture camera images/video                                                 |
|   - Speak through any speaker                                                   |
|   - React to events automatically                                               |
|   - Install itself on new devices via SSH                                       |
|                                                                                 |
+---------------------------------------------------------------------------------+
```

---

## Two Separate Needs: Human vs AI Access

```
+---------------------------------------------------------------------------------+
|                     HUMAN ACCESS vs AI ACCESS                                   |
+---------------------------------------------------------------------------------+
|                                                                                 |
|  TWO SEPARATE NEEDS:                                                            |
|                                                                                 |
|  1. YOU accessing machines (human)     ->  Guacamole (browser GUI)             |
|  2. Legacy accessing machines (AI)     ->  RemoteAgent + NATS                  |
|                                                                                 |
|  These are DIFFERENT solutions for DIFFERENT purposes                           |
|                                                                                 |
|  +----------------------------------+  +----------------------------------+    |
|  |           GUACAMOLE              |  |        REMOTEAGENT + NATS        |    |
|  +----------------------------------+  +----------------------------------+    |
|  |                                  |  |                                  |    |
|  |  FOR: You (human)                |  |  FOR: Legacy (AI)                |    |
|  |                                  |  |                                  |    |
|  |  - Visual desktop access         |  |  - Programmatic commands         |    |
|  |  - Interactive sessions          |  |  - Automated responses           |    |
|  |  - Debugging, troubleshooting    |  |  - Event-driven actions          |    |
|  |  - Showing others your screen    |  |  - Parallel multi-location       |    |
|  |                                  |  |                                  |    |
|  |  ACCESS: Browser anywhere        |  |  ACCESS: NATS pub/sub            |    |
|  |                                  |  |                                  |    |
|  +----------------------------------+  +----------------------------------+    |
|                                                                                 |
|  EXAMPLE:                                                                       |
|  - You want to debug something visually?        -> Guacamole                   |
|  - Legacy needs to restart a service?           -> RemoteAgent                 |
|  - You want to show someone your screen?        -> Guacamole                   |
|  - Legacy needs to check disk space?            -> RemoteAgent                 |
|  - You need to install software manually?       -> Guacamole                   |
|  - Legacy needs to read sensor data?            -> RemoteAgent                 |
|                                                                                 |
+---------------------------------------------------------------------------------+
```

---

## Service Summary

| Profile | Services | Purpose |
|---------|----------|---------|
| **Core** | 5 | postgres, redis, nats, agno, ui |
| **manage** | 3 | Guacamole (YOUR GUI access) |
| **data** | 3 | Nango (OAuth for APIs) |
| **observe** | 3 | Monitoring (Grafana, Victoria, Netdata) |
| **webhooks** | 1 | External webhooks (Redpanda) |
| **archive** | 1 | TimescaleDB (event history) |
| **rag** | 1 | LlamaIndex (document search) |
| **models** | 1 | LiteLLM (model proxy) |
| **Edge (per location)** | 2-3 | RemoteAgent + NATS Leaf + (Ollama) |

---

## Quick Reference: Legacy Commands via NATS

```python
# Check greenhouse temperature
await nats.request("commands.farm.sensors", {
    "sensor": "greenhouse_temp"
})

# Turn on sprinklers
await nats.request("commands.farm.irrigation", {
    "action": "on",
    "zone": "garden",
    "duration_minutes": 30
})

# Close chicken coop
await nats.request("commands.farm.coop.door", {
    "action": "close"
})

# Set home thermostat
await nats.request("commands.home.hvac", {
    "temp": 72,
    "mode": "heat"
})

# Unlock front door
await nats.request("commands.home.lock.front_door", {
    "action": "unlock"
})

# Run shell command on office PC
await nats.request("commands.office.shell", {
    "command": "df -h"
})

# Broadcast to ALL locations
await nats.publish("commands.*.status", {
    "action": "report"
})
```

---

## Next Steps

1. **Create RemoteAgent package** - `agents/edge/`
2. **Create install script** - `infra/scripts/install-agent.sh`
3. **Add edge compose** - `infra/compose/edge/docker-compose.yml`
4. **Deploy to first edge location** - Test with Raspberry Pi
5. **Add tools for your specific devices** - GPIO, cameras, sensors

---

> Document maintained by Legacy for future reference


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
