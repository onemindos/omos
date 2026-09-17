# Live Hermes Plugin Examples

> **Created:** 2026-05-19
> **Purpose:** Working plugin examples showing direct Python SDK imports for Zeus's hardware/robotics stack
> **Plugin Path:** `~/.hermes/plugins/<name>/`

---

## Plugin Structure Reference

Each plugin is a directory with two files:

```
~/.hermes/plugins/my-plugin/
├── plugin.yaml      # manifest (metadata, hooks, env requirements)
└── __init__.py      # must export register(ctx) function
```

**Available Lifecycle Hooks:**

| Hook | When it fires |
|------|---------------|
| `pre_tool_call` | Before any tool executes (can block) |
| `post_tool_call` | After a tool completes |
| `transform_terminal_output` | Reshape terminal output before conversation |
| `transform_tool_result` | Reshape any tool result |
| `transform_llm_output` | Reshape LLM response before display |
| `pre_llm_call` | Before LLM API call (can inject context) |
| `post_llm_call` | After LLM responds |
| `on_session_start` | Session begins |
| `on_session_end` | Session ends |
| `on_session_finalize` | Session finalized |
| `on_session_reset` | Session reset |
| `pre_gateway_dispatch` | Gateway incoming message (can skip/rewrite/allow) |

**PluginContext API:**

| Method | Purpose |
|--------|---------|
| `ctx.register_hook(name, callback)` | Register a lifecycle hook |
| `ctx.register_tool(name, toolset, schema, handler, ...)` | Register a tool |
| `ctx.register_command(name, handler, description)` | Register a `/slash` command |
| `ctx.dispatch_tool(tool_name, args)` | Call another tool programmatically |
| `ctx.inject_message(content, role)` | Inject message into conversation |
| `ctx.llm` | Host-owned LLM facade for model calls |

---

## 1. Brilliant Labs Frame (Wearables)

**What it does:** Gives the agent direct access to AR glasses — camera, microphone, display, and motion sensors.

**Use case:** "Take a photo of what I'm looking at" → agent calls `frame_take_photo` → gets image → analyzes with vision model → tells you what it sees. All in-process, no subprocess, no bridge.

```
~/.hermes/plugins/frame-bridge/
├── plugin.yaml
└── __init__.py
```

### plugin.yaml

```yaml
name: frame-bridge
version: 1.0.0
description: "Brilliant Labs Frame glasses — camera, mic, display, motion"
kind: standalone
hooks:
  - on_session_start
  - on_session_end
provides_tools:
  - frame_take_photo
  - frame_record_audio
  - frame_show_text
  - frame_wait_tap
```

### __init__.py

```python
import asyncio
from frame_sdk import Frame

_frame: Frame | None = None

async def _connect_frame():
    global _frame
    _frame = Frame()
    await _frame.__aenter__()

async def _disconnect_frame():
    global _frame
    if _frame:
        await _frame.__aexit__(None, None, None)
        _frame = None

async def _take_photo(args, ctx):
    """Capture what you're looking at through the glasses."""
    if not _frame:
        return {"error": "Frame not connected"}
    quality = args.get("quality", "medium")
    image_bytes = await _frame.camera.take_photo(quality=quality)
    # Return as base64 for the agent to analyze with vision
    import base64
    return {
        "image_base64": base64.b64encode(image_bytes).decode(),
        "format": "jpeg",
        "note": "Photo captured from Brilliant Frame glasses"
    }

async def _record_audio(args, ctx):
    """Record audio from glasses microphone."""
    if not _frame:
        return {"error": "Frame not connected"}
    duration = args.get("duration", 5)
    audio = await _frame.microphone.record_audio(max_length=duration)
    return {"audio_bytes": len(audio), "duration": duration}

async def _show_text(args, ctx):
    """Display text on the glasses OLED."""
    if not _frame:
        return {"error": "Frame not connected"}
    text = args["text"]
    align = args.get("align", "center")
    await _frame.display.show_text(text, align=align)
    return {"displayed": text}

async def _wait_tap(args, ctx):
    """Wait for tap gesture on glasses."""
    if not _frame:
        return {"error": "Frame not connected"}
    tap = await _frame.motion.wait_for_tap()
    return {"tap_detected": True, "direction": str(tap)}

def register(ctx):
    ctx.register_hook("on_session_start", lambda *a: asyncio.create_task(_connect_frame()))
    ctx.register_hook("on_session_end", lambda *a: asyncio.create_task(_disconnect_frame()))

    ctx.register_tool(
        name="frame_take_photo", toolset="frame-bridge",
        schema={"type": "object", "properties": {"quality": {"type": "string", "enum": ["low", "medium", "high"]}}},
        handler=_take_photo, emoji="📷"
    )
    ctx.register_tool(
        name="frame_record_audio", toolset="frame-bridge",
        schema={"type": "object", "properties": {"duration": {"type": "integer", "default": 5}}},
        handler=_record_audio, emoji="🎙️"
    )
    ctx.register_tool(
        name="frame_show_text", toolset="frame-bridge",
        schema={"type": "object", "properties": {"text": {"type": "string"}}, "required": ["text"]},
        handler=_show_text, emoji="👓"
    )
    ctx.register_tool(
        name="frame_wait_tap", toolset="frame-bridge",
        schema={"type": "object", "properties": {}},
        handler=_wait_tap, emoji="👆"
    )
```

---

## 2. MediaPipe (Gesture/Pose Detection)

**What it does:** Hand gesture recognition and body pose detection from images.

**Use case:** Frame takes photo → agent calls `detect_gesture` → sees you're pointing → sends command to robot "go where Zeus is pointing."

```
~/.hermes/plugins/mediapipe-vision/
├── plugin.yaml
└── __init__.py
```

### plugin.yaml

```yaml
name: mediapipe-vision
version: 1.0.0
description: "MediaPipe hand gesture and body pose detection"
kind: standalone
provides_tools:
  - detect_gesture
  - detect_pose
```

### __init__.py

```python
import mediapipe as mp
import cv2
import numpy as np
import base64

_hands = None
_pose = None

def _init_models():
    global _hands, _pose
    _hands = mp.solutions.hands.Hands(static_image_mode=True, max_num_hands=2)
    _pose = mp.solutions.pose.Pose(static_image_mode=True)

async def _detect_gesture(args, ctx):
    """Detect hand gestures from an image (base64 or file path)."""
    if not _hands:
        _init_models()

    if "image_base64" in args:
        img_bytes = base64.b64decode(args["image_base64"])
        img_array = np.frombuffer(img_bytes, dtype=np.uint8)
        image = cv2.imdecode(img_array, cv2.IMREAD_COLOR)
    else:
        image = cv2.imread(args["file_path"])

    rgb = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
    results = _hands.process(rgb)

    if not results.multi_hand_landmarks:
        return {"gesture": None, "hands_detected": 0}

    # Simple gesture classification from landmarks
    gestures = []
    for hand_landmarks in results.multi_hand_landmarks:
        # Thumb tip vs index tip distance = pinch detection
        thumb_tip = hand_landmarks.landmark[4]
        index_tip = hand_landmarks.landmark[8]
        pinch_dist = ((thumb_tip.x - index_tip.x)**2 + (thumb_tip.y - index_tip.y)**2)**0.5

        if pinch_dist < 0.05:
            gestures.append("pinch")
        elif all(hand_landmarks.landmark[tip].y < hand_landmarks.landmark[tip - 2].y
                 for tip in [8, 12, 16, 20]):
            gestures.append("open_palm")
        elif hand_landmarks.landmark[8].y < hand_landmarks.landmark[6].y and \
             all(hand_landmarks.landmark[tip].y > hand_landmarks.landmark[tip - 2].y
                 for tip in [12, 16, 20]):
            gestures.append("pointing")
        else:
            gestures.append("unknown")

    return {"gestures": gestures, "hands_detected": len(gestures)}

async def _detect_pose(args, ctx):
    """Detect body pose landmarks from image."""
    if not _pose:
        _init_models()

    image = cv2.imread(args["file_path"])
    rgb = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
    results = _pose.process(rgb)

    if not results.pose_landmarks:
        return {"pose_detected": False}

    landmarks = {
        "nose": {"x": results.pose_landmarks.landmark[0].x, "y": results.pose_landmarks.landmark[0].y},
        "left_wrist": {"x": results.pose_landmarks.landmark[15].x, "y": results.pose_landmarks.landmark[15].y},
        "right_wrist": {"x": results.pose_landmarks.landmark[16].x, "y": results.pose_landmarks.landmark[16].y},
    }
    return {"pose_detected": True, "landmarks": landmarks}

def register(ctx):
    ctx.register_tool(
        name="detect_gesture", toolset="mediapipe-vision",
        schema={"type": "object", "properties": {
            "image_base64": {"type": "string"},
            "file_path": {"type": "string"}
        }},
        handler=_detect_gesture, emoji="🤚"
    )
    ctx.register_tool(
        name="detect_pose", toolset="mediapipe-vision",
        schema={"type": "object", "properties": {
            "file_path": {"type": "string", "description": "Path to image file"}
        }, "required": ["file_path"]},
        handler=_detect_pose, emoji="🏃"
    )
```

---

## 3. NATS Fleet Coordination

**What it does:** Real-time pub/sub communication with all edge devices (drones, robots, sensors) over the Tailscale mesh.

**Use case:** "What's the battery on drone-1?" → `fleet_telemetry(device_id="drone-1")` → returns `{"battery": 72, "lat": 33.45, "lon": -112.07, "altitude": 45}`. No bridge, no subprocess — the NATS client is living in the same async event loop as the agent.

```
~/.hermes/plugins/nats-fleet/
├── plugin.yaml
└── __init__.py
```

### plugin.yaml

```yaml
name: nats-fleet
version: 1.0.0
description: "NATS pub/sub for fleet device coordination over Tailscale"
kind: standalone
hooks:
  - on_session_start
provides_tools:
  - fleet_publish
  - fleet_telemetry
  - fleet_status
requires_env:
  - name: NATS_URL
    description: "NATS server URL (default: nats://100.98.25.10:4222)"
```

### __init__.py

```python
import asyncio
import json
import nats
from nats.aio.client import Client as NATS

_nc: NATS | None = None
_subscriptions = {}
_telemetry_buffer = {}  # Latest telemetry per device

async def _connect(args=None, ctx=None):
    global _nc
    import os
    url = os.environ.get("NATS_URL", "nats://100.98.25.10:4222")  # Tailscale VPS
    _nc = await nats.connect(url)

    # Auto-subscribe to fleet telemetry
    async def _on_telemetry(msg):
        data = json.loads(msg.data.decode())
        device_id = data.get("device_id", msg.subject.split(".")[-1])
        _telemetry_buffer[device_id] = {**data, "_subject": msg.subject, "_ts": asyncio.get_event_loop().time()}

    await _nc.subscribe("fleet.telemetry.>", cb=_on_telemetry)
    await _nc.subscribe("fleet.alert.>", cb=_on_telemetry)

async def _publish_command(args, ctx):
    """Publish a command to a specific device or fleet-wide."""
    if not _nc:
        await _connect()
    subject = args["subject"]  # e.g. "fleet.cmd.drone-1" or "fleet.cmd.all"
    payload = json.dumps(args["payload"]).encode()
    await _nc.publish(subject, payload)
    return {"published": subject, "payload": args["payload"]}

async def _get_telemetry(args, ctx):
    """Get latest telemetry from all devices or a specific one."""
    device_id = args.get("device_id")
    if device_id:
        return _telemetry_buffer.get(device_id, {"error": f"No telemetry from {device_id}"})
    return {k: v for k, v in _telemetry_buffer.items()}

async def _request_status(args, ctx):
    """Request-reply status from a specific device (with timeout)."""
    if not _nc:
        await _connect()
    subject = f"fleet.status.{args['device_id']}"
    try:
        response = await _nc.request(subject, b"status", timeout=5.0)
        return json.loads(response.data.decode())
    except asyncio.TimeoutError:
        return {"error": f"Device {args['device_id']} did not respond within 5s"}

def register(ctx):
    ctx.register_hook("on_session_start", lambda *a: asyncio.create_task(_connect()))

    ctx.register_tool(
        name="fleet_publish", toolset="nats-fleet",
        schema={"type": "object", "properties": {
            "subject": {"type": "string", "description": "NATS subject (e.g. fleet.cmd.drone-1)"},
            "payload": {"type": "object", "description": "JSON payload to send"}
        }, "required": ["subject", "payload"]},
        handler=_publish_command, emoji="📡"
    )
    ctx.register_tool(
        name="fleet_telemetry", toolset="nats-fleet",
        schema={"type": "object", "properties": {
            "device_id": {"type": "string", "description": "Optional specific device"}
        }},
        handler=_get_telemetry, emoji="📊"
    )
    ctx.register_tool(
        name="fleet_status", toolset="nats-fleet",
        schema={"type": "object", "properties": {
            "device_id": {"type": "string"}
        }, "required": ["device_id"]},
        handler=_request_status, emoji="🏓"
    )
```

---

## 4. Meshtastic (Off-Grid LoRa Mesh)

**What it does:** Send/receive messages and track node positions over LoRa mesh when WiFi/Tailscale is unavailable.

**Use case:** When Tailscale/WiFi is down on the ranch, "send mesh broadcast: irrigation pump 3 activated" → goes over LoRa to all mesh nodes within radio range.

```
~/.hermes/plugins/meshtastic-mesh/
├── plugin.yaml
└── __init__.py
```

### plugin.yaml

```yaml
name: meshtastic-mesh
version: 1.0.0
description: "Meshtastic LoRa mesh — off-grid messaging and node tracking"
kind: standalone
hooks:
  - on_session_start
provides_tools:
  - mesh_send
  - mesh_nodes
  - mesh_position
requires_env:
  - name: MESHTASTIC_HOST
    description: "Meshtastic device IP or serial port"
```

### __init__.py

```python
import meshtastic
import meshtastic.serial_interface
import meshtastic.tcp_interface

_interface = None

async def _connect_mesh(args=None, ctx=None):
    global _interface
    import os
    host = os.environ.get("MESHTASTIC_HOST", "192.168.1.50")  # or serial
    _interface = meshtastic.tcp_interface.TCPInterface(hostname=host)

async def _send_mesh_message(args, ctx):
    """Send a text message over LoRa mesh network."""
    if not _interface:
        await _connect_mesh()
    text = args["text"]
    dest = args.get("destination", "^all")  # ^all = broadcast
    _interface.sendText(text, destinationId=dest)
    return {"sent": text, "to": dest, "channel": "default"}

async def _get_mesh_nodes(args, ctx):
    """List all nodes visible on the mesh network."""
    if not _interface:
        await _connect_mesh()
    nodes = _interface.nodes
    result = []
    for node_id, node in nodes.items():
        result.append({
            "id": node_id,
            "name": node.get("user", {}).get("longName", "Unknown"),
            "snr": node.get("snr"),
            "last_heard": node.get("lastHeard"),
            "position": node.get("position"),
            "battery": node.get("deviceMetrics", {}).get("batteryLevel")
        })
    return {"nodes": result, "count": len(result)}

async def _get_mesh_position(args, ctx):
    """Get GPS position of a specific mesh node."""
    if not _interface:
        await _connect_mesh()
    node_id = args["node_id"]
    node = _interface.nodes.get(node_id, {})
    pos = node.get("position", {})
    return {"lat": pos.get("latitude"), "lon": pos.get("longitude"),
            "alt": pos.get("altitude"), "time": pos.get("time")}

def register(ctx):
    ctx.register_hook("on_session_start", lambda *a: asyncio.create_task(_connect_mesh()))

    ctx.register_tool(name="mesh_send", toolset="meshtastic-mesh",
        schema={"type": "object", "properties": {
            "text": {"type": "string"}, "destination": {"type": "string"}
        }, "required": ["text"]},
        handler=_send_mesh_message, emoji="📻")

    ctx.register_tool(name="mesh_nodes", toolset="meshtastic-mesh",
        schema={"type": "object", "properties": {}},
        handler=_get_mesh_nodes, emoji="🗺️")

    ctx.register_tool(name="mesh_position", toolset="meshtastic-mesh",
        schema={"type": "object", "properties": {"node_id": {"type": "string"}}, "required": ["node_id"]},
        handler=_get_mesh_position, emoji="📍")
```

---

## 5. DMX Light Shows / Projection

**What it does:** Control DMX lighting fixtures — individual channels or predefined scenes (sunrise, alert, party).

**Use case:** "Turn on sunrise lighting in the barn" → `light_scene(scene="sunrise", duration=30)`. Or tied to cron: every 5:30 AM, agent runs sunrise scene automatically. If a drone alert fires, agent runs `alert_red`.

```
~/.hermes/plugins/dmx-lights/
├── plugin.yaml
└── __init__.py
```

### plugin.yaml

```yaml
name: dmx-lights
version: 1.0.0
description: "DMX lighting control — scenes, channels, automation"
kind: standalone
provides_tools:
  - light_set
  - light_scene
requires_env:
  - name: DMX_PORT
    description: "USB DMX interface serial port (e.g. /dev/tty.usbserial-EN123456)"
```

### __init__.py

```python
import asyncio
from dmxpy.DmxPy import DmxPy  # USB DMX interface

_dmx = None

async def _init_dmx(args=None, ctx=None):
    global _dmx
    import os
    port = os.environ.get("DMX_PORT", "/dev/tty.usbserial-EN123456")
    _dmx = DmxPy(port)

async def _set_light(args, ctx):
    """Set a DMX channel to a value (0-255). Controls individual fixtures."""
    if not _dmx:
        await _init_dmx()
    channel = args["channel"]
    value = args["value"]
    _dmx.setChannel(channel, value)
    _dmx.render()
    return {"channel": channel, "value": value}

async def _run_scene(args, ctx):
    """Run a predefined light scene (fade, chase, strobe, etc.)."""
    if not _dmx:
        await _init_dmx()
    scene = args["scene"]
    duration = args.get("duration", 5)

    if scene == "sunrise":
        # Warm fade up over duration
        for i in range(0, 256, 5):
            _dmx.setChannel(1, i)        # warm white
            _dmx.setChannel(2, i // 3)   # amber
            _dmx.setChannel(3, 0)        # cool white off
            _dmx.render()
            await asyncio.sleep(duration / 51)
        return {"scene": "sunrise", "completed": True}

    elif scene == "alert_red":
        # Flash red 5 times
        for _ in range(5):
            _dmx.setChannel(4, 255)  # red
            _dmx.render()
            await asyncio.sleep(0.3)
            _dmx.setChannel(4, 0)
            _dmx.render()
            await asyncio.sleep(0.3)
        return {"scene": "alert_red", "completed": True}

    elif scene == "all_off":
        for ch in range(1, 513):
            _dmx.setChannel(ch, 0)
        _dmx.render()
        return {"scene": "all_off", "completed": True}

    return {"error": f"Unknown scene: {scene}"}

def register(ctx):
    ctx.register_tool(name="light_set", toolset="dmx-lights",
        schema={"type": "object", "properties": {
            "channel": {"type": "integer", "minimum": 1, "maximum": 512},
            "value": {"type": "integer", "minimum": 0, "maximum": 255}
        }, "required": ["channel", "value"]},
        handler=_set_light, emoji="💡")

    ctx.register_tool(name="light_scene", toolset="dmx-lights",
        schema={"type": "object", "properties": {
            "scene": {"type": "string", "enum": ["sunrise", "alert_red", "all_off", "party", "focus"]},
            "duration": {"type": "number", "default": 5}
        }, "required": ["scene"]},
        handler=_run_scene, emoji="🎭")
```

---

## 6. MAVSDK (Direct Drone Control)

**What it does:** Direct MAVLink drone control for simple operations where DimOS is overkill.

**Use case:** For simple ops — "take off and hover at 10m" → direct MAVSDK call. For complex ops (visual tracking, obstacle avoidance), use the DimOS MCP server instead.

```
~/.hermes/plugins/drone-direct/
├── plugin.yaml
└── __init__.py
```

### plugin.yaml

```yaml
name: drone-direct
version: 1.0.0
description: "Direct MAVLink drone control via MAVSDK — takeoff, fly_to, land, telemetry"
kind: standalone
provides_tools:
  - drone_takeoff
  - drone_fly_to
  - drone_land
  - drone_position
requires_env:
  - name: DRONE_ADDRESS
    description: "MAVLink connection string (default: udp://:14540)"
```

### __init__.py

```python
from mavsdk import System
from mavsdk.offboard import VelocityNedYaw
import asyncio

_drone: System | None = None

async def _connect_drone(args=None, ctx=None):
    global _drone
    import os
    addr = os.environ.get("DRONE_ADDRESS", "udp://:14540")
    _drone = System()
    await _drone.connect(system_address=addr)
    async for state in _drone.core.connection_state():
        if state.is_connected:
            break

async def _takeoff(args, ctx):
    """Arm and takeoff to specified altitude."""
    if not _drone:
        return {"error": "Drone not connected"}
    altitude = args.get("altitude", 10)
    await _drone.action.set_takeoff_altitude(altitude)
    await _drone.action.arm()
    await _drone.action.takeoff()
    return {"status": "taking off", "target_altitude": altitude}

async def _fly_to(args, ctx):
    """Fly to GPS coordinates."""
    if not _drone:
        return {"error": "Drone not connected"}
    lat, lon, alt = args["lat"], args["lon"], args.get("alt", 30)
    await _drone.action.goto_location(lat, lon, alt, 0)
    return {"status": "en route", "destination": {"lat": lat, "lon": lon, "alt": alt}}

async def _land(args, ctx):
    """Land the drone at current position."""
    if not _drone:
        return {"error": "Drone not connected"}
    await _drone.action.land()
    return {"status": "landing"}

async def _get_position(args, ctx):
    """Get current GPS position and battery."""
    if not _drone:
        return {"error": "Drone not connected"}
    async for pos in _drone.telemetry.position():
        async for batt in _drone.telemetry.battery():
            return {
                "lat": pos.latitude_deg, "lon": pos.longitude_deg,
                "alt": pos.relative_altitude_m,
                "battery_pct": round(batt.remaining_percent * 100)
            }

def register(ctx):
    ctx.register_tool(name="drone_takeoff", toolset="drone-direct",
        schema={"type": "object", "properties": {"altitude": {"type": "number", "default": 10}}},
        handler=_takeoff, emoji="🛫")
    ctx.register_tool(name="drone_fly_to", toolset="drone-direct",
        schema={"type": "object", "properties": {
            "lat": {"type": "number"}, "lon": {"type": "number"}, "alt": {"type": "number"}
        }, "required": ["lat", "lon"]},
        handler=_fly_to, emoji="✈️")
    ctx.register_tool(name="drone_land", toolset="drone-direct",
        schema={"type": "object", "properties": {}},
        handler=_land, emoji="🛬")
    ctx.register_tool(name="drone_position", toolset="drone-direct",
        schema={"type": "object", "properties": {}},
        handler=_get_position, emoji="📍")
```

---

## 7. OpenCV + YOLO (Vision Analysis)

**What it does:** Object detection and counting from images — farm animals, vehicles, people, equipment.

**Use case:** Drone captures field image → `count_objects(file_path="/tmp/field.jpg", target_class="cow")` → returns `{"class": "cow", "count": 47}`. Or Frame glasses snap → detect everything in view.

```
~/.hermes/plugins/vision-detect/
├── plugin.yaml
└── __init__.py
```

### plugin.yaml

```yaml
name: vision-detect
version: 1.0.0
description: "YOLO object detection and counting from images"
kind: standalone
provides_tools:
  - detect_objects
  - count_objects
```

### __init__.py

```python
import cv2
import numpy as np
from ultralytics import YOLO
import base64

_model = None

def _load_model():
    global _model
    _model = YOLO("yolov8n.pt")  # Nano model, fast inference

async def _detect_objects(args, ctx):
    """Run YOLO object detection on an image."""
    if not _model:
        _load_model()

    if "image_base64" in args:
        img_bytes = base64.b64decode(args["image_base64"])
        img_array = np.frombuffer(img_bytes, dtype=np.uint8)
        image = cv2.imdecode(img_array, cv2.IMREAD_COLOR)
    else:
        image = cv2.imread(args["file_path"])

    results = _model(image, verbose=False)[0]

    detections = []
    for box in results.boxes:
        detections.append({
            "class": results.names[int(box.cls[0])],
            "confidence": round(float(box.conf[0]), 3),
            "bbox": [int(x) for x in box.xyxy[0].tolist()]
        })

    return {"detections": detections, "count": len(detections)}

async def _count_objects(args, ctx):
    """Count specific objects in an image (e.g. 'how many cows in field?')."""
    result = await _detect_objects(args, ctx)
    target = args.get("target_class", "").lower()
    if target:
        matches = [d for d in result["detections"] if d["class"].lower() == target]
        return {"class": target, "count": len(matches), "confidence_avg":
                round(sum(d["confidence"] for d in matches) / max(len(matches), 1), 3)}
    return result

def register(ctx):
    ctx.register_tool(name="detect_objects", toolset="vision-detect",
        schema={"type": "object", "properties": {
            "image_base64": {"type": "string"},
            "file_path": {"type": "string"}
        }},
        handler=_detect_objects, emoji="👁️")
    ctx.register_tool(name="count_objects", toolset="vision-detect",
        schema={"type": "object", "properties": {
            "image_base64": {"type": "string"},
            "file_path": {"type": "string"},
            "target_class": {"type": "string", "description": "e.g. 'cow', 'person', 'car'"}
        }},
        handler=_count_objects, emoji="🔢")
```

---

## 8. Combo Pipeline: Frame → Gesture → Robot Command

**What it does:** Chains multiple plugins together — captures from glasses, detects gesture, sends robot command. Demonstrates `ctx.dispatch_tool` for inter-plugin communication.

**Use case:** "Legacy, watch my hands" → agent runs `gesture_command` in a loop → you point at something → robot walks toward it → you open palm → everything stops. Zero latency penalty from bridges.

```
~/.hermes/plugins/gesture-to-robot/
├── plugin.yaml
└── __init__.py
```

### plugin.yaml

```yaml
name: gesture-to-robot
version: 1.0.0
description: "Chain: Frame camera → MediaPipe gesture → NATS robot command"
kind: standalone
provides_tools:
  - gesture_command
```

### __init__.py

```python
"""
Chain: Frame camera → MediaPipe gesture → NATS robot command
No bridges, no subprocesses — all in-process Python.
"""

async def _gesture_command_loop(args, ctx):
    """Take a photo, detect gesture, send appropriate robot command."""
    # 1. Capture from Frame
    photo_result = await ctx.dispatch_tool("frame_take_photo", {"quality": "medium"})
    if "error" in photo_result:
        return photo_result

    # 2. Detect gesture
    gesture_result = await ctx.dispatch_tool("detect_gesture", {
        "image_base64": photo_result["image_base64"]
    })

    gestures = gesture_result.get("gestures", [])
    if not gestures:
        return {"action": "none", "reason": "No gesture detected"}

    gesture = gestures[0]

    # 3. Map gesture → robot command
    command_map = {
        "pointing": {"subject": "fleet.cmd.unitree-g1", "payload": {"action": "go_to_point", "direction": "forward"}},
        "open_palm": {"subject": "fleet.cmd.all", "payload": {"action": "stop"}},
        "pinch": {"subject": "fleet.cmd.drone-1", "payload": {"action": "take_photo"}},
    }

    if gesture in command_map:
        cmd = command_map[gesture]
        await ctx.dispatch_tool("fleet_publish", cmd)
        return {"gesture": gesture, "command_sent": cmd}

    return {"gesture": gesture, "action": "unmapped"}

def register(ctx):
    ctx.register_tool(name="gesture_command", toolset="gesture-to-robot",
        schema={"type": "object", "properties": {}},
        handler=_gesture_command_loop, emoji="🤖")
```

---

## Comparison: Hermes vs OpenClaw For Same Tasks

| In Hermes | In OpenClaw |
|---|---|
| `import frame_sdk` | MCP server + IPC |
| `import mediapipe` | MCP server + IPC |
| `import nats` | Native (nats.js works) ✅ |
| `import meshtastic` | `exec('python3 mesh.py')` |
| `from mavsdk import System` | MCP server + IPC |
| `from ultralytics import YOLO` | MCP server + IPC |
| `from dmxpy import DmxPy` | `exec('python3 dmx.py')` |
| `ctx.dispatch_tool("other_plugin_tool")` | Internal tool call (works) ✅ |

NATS and inter-tool dispatch work fine in both. The difference is everything touching hardware SDKs — those are Python-native and Hermes avoids the bridge tax.

---

## Installation Commands

```bash
# Install Python dependencies for all plugins
pip install frame-sdk mediapipe opencv-python numpy ultralytics \
    nats-py meshtastic mavsdk dmxpy

# Enable plugins in Hermes
hermes plugins enable frame-bridge
hermes plugins enable mediapipe-vision
hermes plugins enable nats-fleet
hermes plugins enable meshtastic-mesh
hermes plugins enable dmx-lights
hermes plugins enable drone-direct
hermes plugins enable vision-detect
hermes plugins enable gesture-to-robot
```
