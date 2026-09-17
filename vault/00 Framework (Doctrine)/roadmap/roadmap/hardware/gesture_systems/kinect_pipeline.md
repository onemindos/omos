# Kinect + TouchDesigner + Unreal Engine Pipeline

## The Ultimate Gesture-to-Visualization System

This document describes how to build a professional-grade gesture control system
using depth cameras, visual programming, and real-time 3D rendering.

---

## Architecture Overview

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                    PHYSICAL → DIGITAL → PHYSICAL PIPELINE                    │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│   ┌─────────┐      ┌─────────────┐      ┌────────┐      ┌─────────────┐      │
│   │         │      │             │      │        │      │             │      │
│   │ KINECT  │─────►│TOUCHDESIGNER│─────►│  NATS  │─────►│   UNREAL    │      │
│   │         │      │             │      │        │      │   ENGINE    │      │
│   └─────────┘      └──────┬──────┘      └───┬────┘      └─────────────┘      │
│                           │                 │                                 │
│   Depth Camera            │                 │           3D Visualization     │
│   Body Tracking           │                 │           Particle Systems     │
│   Hand Detection          │                 │           Spatial Audio        │
│                           │                 │                                 │
│                           ▼                 ▼                                 │
│                    ┌─────────────────────────────────┐                       │
│                    │     OneMind OS Backend          │                       │
│                    │                                 │                       │
│                    │  Python Services:               │                       │
│                    │  • Gesture Classification       │                       │
│                    │  • Intent Recognition           │                       │
│                    │  • Robot Command Generation     │                       │
│                    │  • Memory/Context Integration   │                       │
│                    └─────────────────────────────────┘                       │
│                                  │                                            │
│                                  ▼                                            │
│                    ┌─────────────────────────────────┐                       │
│                    │      Physical Actuators         │                       │
│                    │                                 │                       │
│                    │  • Drones                       │                       │
│                    │  • Robot Arms                   │                       │
│                    │  • 3D Printers                  │                       │
│                    │  • Smart Home                   │                       │
│                    └─────────────────────────────────┘                       │
│                                                                               │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## Component Details

### 1. Kinect (Sensing Layer)

**Hardware Options:**

| Device                | Depth Resolution     | Body Tracking     | Price        | Status       |
| --------------------- | -------------------- | ----------------- | ------------ | ------------ |
| Kinect v2             | 512x424              | Yes (6 bodies)    | $50-100 used | Discontinued |
| Azure Kinect          | 640x576 to 1024x1024 | Yes (body + hand) | $400         | Current      |
| Intel RealSense D435i | 1280x720             | Via software      | $300         | Current      |
| Orbbec Femto Mega     | 1280x800             | Yes               | $500         | Current      |

**What It Provides:**

```python
# Skeleton data per frame
body = {
    "id": 0,
    "joints": {
        "head": {"x": 0.1, "y": 1.7, "z": 2.0, "confidence": 0.95},
        "neck": {"x": 0.1, "y": 1.5, "z": 2.0, "confidence": 0.92},
        "shoulder_left": ...,
        "shoulder_right": ...,
        "elbow_left": ...,
        "elbow_right": ...,
        "wrist_left": ...,
        "wrist_right": ...,
        "hand_left": {"x": 0.3, "y": 1.2, "z": 1.8, "state": "open"},
        "hand_right": {"x": -0.2, "y": 1.3, "z": 1.5, "state": "closed"},
        # ... more joints
    },
    "lean": {"x": 0.1, "y": 0.0},  # Body lean
    "engaged": True  # User attention
}
```

---

### 2. TouchDesigner (Processing Layer)

TouchDesigner is a node-based visual programming environment for real-time content.

**Why TouchDesigner?**

- Handles Kinect SDK integration out of the box
- Real-time data smoothing and filtering
- Easy OSC/Spout output to other applications
- Visual debugging of data flow

**Node Network Example:**

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Kinect     │────►│   Filter    │────►│  Classify   │────►│    OSC      │
│  CHOP       │     │   CHOP      │     │   (Script)  │     │   Out       │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
     │                                                             │
     │  Raw skeleton                                               │
     │  30-60 fps                                                  │
     │                                                             ▼
     │                                                      ┌─────────────┐
     │                                                      │   Python    │
     │                                                      │   Bridge    │
     │                                                      └─────────────┘
     │                                                             │
     ▼                                                             ▼
┌─────────────┐                                             ┌─────────────┐
│  Spout      │─────────────────────────────────────────────│   NATS      │
│  Out        │                                             │   Publish   │
└─────────────┘                                             └─────────────┘
     │
     ▼
  Unreal Engine
```

**TouchDesigner Python Script for Classification:**

```python
# In TouchDesigner Script CHOP

def onCook(scriptOp):
    # Get Kinect data
    kinect = op('kinectazure1')

    hand_left = {
        'x': kinect['hand_left_x'],
        'y': kinect['hand_left_y'],
        'z': kinect['hand_left_z'],
        'state': kinect['hand_left_state']
    }

    hand_right = {
        'x': kinect['hand_right_x'],
        'y': kinect['hand_right_y'],
        'z': kinect['hand_right_z'],
        'state': kinect['hand_right_state']
    }

    # Classify gesture
    gesture = classify_gesture(hand_left, hand_right)

    # Output to CHOP channels
    scriptOp.clear()
    scriptOp.appendChan('gesture_type').vals = [gesture['type_id']]
    scriptOp.appendChan('gesture_confidence').vals = [gesture['confidence']]
    scriptOp.appendChan('direction_x').vals = [gesture['direction'][0]]
    scriptOp.appendChan('direction_y').vals = [gesture['direction'][1]]
    scriptOp.appendChan('direction_z').vals = [gesture['direction'][2]]


def classify_gesture(left, right):
    # Simplified gesture classification

    # Point: Right hand closed except index
    if right['state'] == 'point':
        direction = calculate_point_direction(right)
        return {'type': 'point', 'type_id': 1, 'confidence': 0.9, 'direction': direction}

    # Grab: Both hands closed
    if left['state'] == 'closed' and right['state'] == 'closed':
        return {'type': 'grab', 'type_id': 2, 'confidence': 0.85, 'direction': [0, 0, 0]}

    # Open: Both hands open (stop gesture)
    if left['state'] == 'open' and right['state'] == 'open':
        return {'type': 'open_stop', 'type_id': 3, 'confidence': 0.95, 'direction': [0, 0, 0]}

    # Swipe: Fast horizontal hand movement
    if abs(right['velocity_x']) > 0.5:
        direction = [1 if right['velocity_x'] > 0 else -1, 0, 0]
        return {'type': 'swipe', 'type_id': 4, 'confidence': 0.8, 'direction': direction}

    return {'type': 'none', 'type_id': 0, 'confidence': 0, 'direction': [0, 0, 0]}
```

---

### 3. Python Bridge (NATS Integration)

**OSC → NATS Bridge:**

```python
# python_bridge.py
import asyncio
import json
from pythonosc import dispatcher, osc_server
import nats

class GestureNATSBridge:
    def __init__(self, osc_port=5005, nats_url="nats://localhost:4222"):
        self.osc_port = osc_port
        self.nats_url = nats_url
        self.nc = None

    async def connect(self):
        self.nc = await nats.connect(self.nats_url)
        print(f"[Bridge] Connected to NATS at {self.nats_url}")

    def handle_gesture(self, address, *args):
        """Handle incoming OSC gesture message."""
        # Parse OSC data
        gesture_data = {
            "type": args[0],
            "confidence": args[1],
            "hand": args[2],
            "position": {"x": args[3], "y": args[4], "z": args[5]},
            "direction": {"x": args[6], "y": args[7], "z": args[8]},
            "timestamp": args[9] if len(args) > 9 else None
        }

        # Publish to NATS
        asyncio.create_task(self._publish(gesture_data))

        # Log
        print(f"[Gesture] {gesture_data['type']} ({gesture_data['confidence']:.2f})")

    async def _publish(self, data):
        if self.nc:
            topic = f"gesture.classified"
            await self.nc.publish(topic, json.dumps(data).encode())

    async def run(self):
        await self.connect()

        # Set up OSC server
        disp = dispatcher.Dispatcher()
        disp.map("/gesture/*", self.handle_gesture)

        server = osc_server.AsyncIOOSCUDPServer(
            ("0.0.0.0", self.osc_port),
            disp,
            asyncio.get_event_loop()
        )
        transport, protocol = await server.create_serve_endpoint()

        print(f"[Bridge] OSC server listening on port {self.osc_port}")

        # Keep running
        while True:
            await asyncio.sleep(1)


if __name__ == "__main__":
    bridge = GestureNATSBridge()
    asyncio.run(bridge.run())
```

---

### 4. Unreal Engine (Visualization Layer)

**OSC Plugin Setup:**

1. Enable OSC plugin in Unreal
2. Create OSC Server actor
3. Bind to gesture messages
4. Drive Niagara particle systems

**Blueprint Overview:**

```
Event Graph:
───────────────────────────────────────────────────────────
[Event BeginPlay]
    │
    ▼
[Create OSC Server] ─────► [Bind to /gesture/point]
                               │
                               ▼
                          [On Message Received]
                               │
                               ▼
                          [Parse Position/Direction]
                               │
                               ▼
                          [Set Niagara Location/Direction]
                               │
                               ▼
                          [Spawn Particle Burst]
```

**Niagara System for Gesture Particles:**

```
Emitter Properties:
├── Spawn Rate: Linked to gesture confidence
├── Initial Velocity: Linked to gesture direction
├── Color: Mapped to gesture type
│   ├── Point: Blue trail
│   ├── Grab: Orange cluster
│   ├── Open: Green disperse
│   └── Swipe: Purple streak
└── Lifetime: 2-4 seconds with fade
```

---

## Latency Breakdown

```
┌─────────────────────────────────────────────────────────────────┐
│                    END-TO-END LATENCY                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Kinect Capture:        ~33ms (30fps) or ~16ms (60fps)          │
│  USB Transfer:          ~5ms                                     │
│  TouchDesigner Process: ~8ms (filter + classify)                │
│  OSC Network:           ~1ms (local) or ~5ms (network)          │
│  NATS Publish:          ~1ms                                     │
│  Unreal Receive:        ~8ms (next frame)                       │
│  Render:                ~8ms (120fps) or ~16ms (60fps)          │
│  ─────────────────────────────────────────────────────          │
│  TOTAL:                 ~70-100ms (perceptually instant)        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Example Use Cases

### 1. Drone Control via Pointing

```python
# drone_gesture_controller.py

async def handle_gesture(msg):
    data = json.loads(msg.data)

    if data['type'] == 'point' and data['confidence'] > 0.8:
        # Map point direction to world coordinates
        direction = data['direction']

        # Calculate target waypoint
        target = calculate_waypoint(
            current_drone_pos,
            direction,
            distance=5.0  # 5 meters in point direction
        )

        # Send command to drone
        await drone.fly_to(target)

# Subscribe to gestures
await nc.subscribe("gesture.classified", cb=handle_gesture)
```

### 2. 3D Model Manipulation

```python
# cad_gesture_controller.py

async def handle_gesture(msg):
    data = json.loads(msg.data)

    if data['type'] == 'grab':
        # Start rotation mode
        rotation_active = True

    elif data['type'] == 'open_stop':
        # Confirm current orientation
        rotation_active = False
        await save_model_state()

    elif data['type'] == 'pinch' and rotation_active:
        # Rotate model based on hand movement
        rotation = calculate_rotation(data['position'])
        await rotate_model(rotation)
```

### 3. Immersive Dashboard Navigation

```python
# dashboard_gesture_controller.py

async def handle_gesture(msg):
    data = json.loads(msg.data)

    if data['type'] == 'swipe':
        direction = data['direction']['x']
        if direction > 0.5:
            await switch_dashboard_panel('next')
        elif direction < -0.5:
            await switch_dashboard_panel('previous')

    elif data['type'] == 'point':
        # Select element at point target
        element = raycast_to_ui(data['direction'])
        if element:
            await select_element(element)
```

---

## Hardware Setup Guide

### Kinect v2 (Budget Option)

```bash
# macOS/Linux
brew install libfreenect2
pip install freenect2

# Windows
# Download Kinect for Windows SDK v2.0
```

### Azure Kinect (Recommended)

```bash
# Install Azure Kinect SDK
# https://docs.microsoft.com/en-us/azure/kinect-dk/set-up-azure-kinect-dk

# Python bindings
pip install pyk4a
```

### TouchDesigner

1. Download from [derivative.ca](https://derivative.ca/)
2. Free for non-commercial use (with watermark)
3. $600/year for commercial license

### Unreal Engine

1. Download from [unrealengine.com](https://www.unrealengine.com/)
2. Free (5% royalty after $1M revenue)
3. Enable OSC plugin in Edit → Plugins

---

## Alternative: Pure Python (No TouchDesigner)

If you want to skip TouchDesigner:

```python
# direct_kinect_to_nats.py
import pyk4a
from pyk4a import Config, PyK4A
import asyncio
import json
import nats

async def main():
    # Initialize Kinect
    k4a = PyK4A(Config(
        color_resolution=pyk4a.ColorResolution.RES_720P,
        depth_mode=pyk4a.DepthMode.NFOV_UNBINNED,
    ))
    k4a.start()

    # Initialize body tracker
    tracker = k4a.start_body_tracker()

    # Connect to NATS
    nc = await nats.connect("nats://localhost:4222")

    while True:
        # Get frame
        capture = k4a.get_capture()

        if capture.depth is not None:
            # Get body frame
            body_frame = tracker.update(capture)

            for body in body_frame.get_bodies():
                # Extract hand data
                gesture = classify_gesture(body)

                # Publish to NATS
                await nc.publish(
                    "gesture.classified",
                    json.dumps(gesture).encode()
                )

        await asyncio.sleep(0.016)  # ~60fps

asyncio.run(main())
```

---

## Resources

- [Azure Kinect DK Documentation](https://docs.microsoft.com/en-us/azure/kinect-dk/)
- [TouchDesigner Wiki](https://docs.derivative.ca/)
- [Unreal Engine OSC Plugin](https://docs.unrealengine.com/en-US/WorkingWithMedia/IntegratingMedia/OSC/)
- [PyK4A (Python Azure Kinect)](https://github.com/etiennedub/pyk4a)
- [python-osc](https://pypi.org/project/python-osc/)
