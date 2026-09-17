# Depth Cameras Comparison

## Overview

Depth cameras enable:

- Full body skeleton tracking
- Hand pose estimation
- Room mapping (SLAM)
- Object recognition with depth
- Gesture control

---

## Hardware Comparison

| Camera                       | Depth Range | Resolution                     | Body Tracking      | Price        | Status       |
| ---------------------------- | ----------- | ------------------------------ | ------------------ | ------------ | ------------ |
| **Kinect v2**                | 0.5-4.5m    | 512×424 (depth), 1080p (color) | Yes (6 bodies)     | $50-100 used | Discontinued |
| **Azure Kinect DK**          | 0.25-5.5m   | 640×576 to 1024×1024           | Yes (body + hands) | $400         | Current      |
| **Intel RealSense D435**     | 0.1-10m     | 1280×720                       | Via software       | $300         | Current      |
| **Intel RealSense D455**     | 0.4-6m      | 1280×800                       | Via software       | $350         | Current      |
| **Orbbec Femto Mega**        | 0.25-5.5m   | 1280×800                       | Yes                | $500         | Current      |
| **ZED 2i**                   | 0.2-20m     | Up to 2K                       | Yes                | $450         | Current      |
| **Leap Motion Controller 2** | 0.1-1m      | N/A (hand only)                | Hands only         | $140         | Current      |

---

## Kinect v2 (Budget Option)

### Pros:

- Very cheap (used market)
- Good skeleton tracking
- Wide field of view
- Well documented

### Cons:

- Discontinued (no official support)
- USB 3.0 only
- Windows-centric SDK (libfreenect2 for cross-platform)

### Setup (Python):

```bash
# macOS/Linux
brew install libfreenect2
pip install freenect2

# Windows
# Download Kinect for Windows SDK v2.0 from web.archive.org
```

### Code Example:

```python
from freenect2 import Device, FrameType

device = Device()
device.start()

while True:
    frames = device.get_frames()
    depth = frames[FrameType.Depth]
    color = frames[FrameType.Color]
    # Process frames...
```

---

## Azure Kinect DK (Recommended)

### Pros:

- Official Microsoft support
- Excellent body tracking SDK
- Hand tracking included
- 4K color camera
- Cross-platform (Windows, Linux)

### Cons:

- Higher price ($400)
- Requires Azure Kinect SDK

### Setup (Python):

```bash
# Install Azure Kinect SDK first
# https://docs.microsoft.com/en-us/azure/kinect-dk/set-up-azure-kinect-dk

pip install pyk4a
```

### Code Example:

```python
from pyk4a import PyK4A, Config

k4a = PyK4A(Config(
    color_resolution=pyk4a.ColorResolution.RES_720P,
    depth_mode=pyk4a.DepthMode.NFOV_UNBINNED,
))
k4a.start()

while True:
    capture = k4a.get_capture()
    if capture.depth is not None:
        depth_image = capture.depth
        color_image = capture.color
        # Process...
```

### Body Tracking:

```python
from pyk4a import PyK4A
from pyk4abt import PyK4ABodyTracker

k4a = PyK4A()
k4a.start()

tracker = PyK4ABodyTracker()

while True:
    capture = k4a.get_capture()
    frame = tracker.update(capture)

    for body in frame.get_bodies():
        for joint in body.get_joints():
            print(f"{joint.id}: {joint.position}")
```

---

## Intel RealSense (Flexible Option)

### Pros:

- Good documentation
- Cross-platform
- Multiple models for different uses
- Strong Python SDK

### Cons:

- No built-in body tracking (needs external solution)
- Can be noisy in certain conditions

### Setup:

```bash
pip install pyrealsense2
```

### Code Example:

```python
import pyrealsense2 as rs
import numpy as np

pipeline = rs.pipeline()
config = rs.config()
config.enable_stream(rs.stream.depth, 640, 480, rs.format.z16, 30)
config.enable_stream(rs.stream.color, 640, 480, rs.format.bgr8, 30)

pipeline.start(config)

while True:
    frames = pipeline.wait_for_frames()
    depth_frame = frames.get_depth_frame()
    color_frame = frames.get_color_frame()

    depth_image = np.asanyarray(depth_frame.get_data())
    color_image = np.asanyarray(color_frame.get_data())

    # Process...
```

---

## Leap Motion (Hand-Only)

### Pros:

- Extremely precise hand tracking
- Very low latency
- Small form factor
- Works well for VR/AR

### Cons:

- Hand only (no body)
- Short range (< 1m)
- Limited field of view

### Use Case:

Best for precision hand gestures in desktop/VR setups, not room-scale.

---

## Recommendations for OneMind OS

### Budget Setup:

```
Kinect v2 (used) + libfreenect2
├── Body tracking for gestures
├── ~$50-100
└── Good for prototyping
```

### Production Setup:

```
Azure Kinect DK
├── Full body + hand tracking
├── High quality depth
├── Official SDK support
└── ~$400
```

### Portable/Flexible:

```
Intel RealSense D435i
├── USB-C portable
├── Good for robots/drones
├── IMU included for motion
└── ~$300
```

### Precision Hands Only:

```
Leap Motion 2
├── Desktop hand tracking
├── Very precise
├── Good for CAD manipulation
└── ~$140
```

---

## Integration with OneMind

All cameras feed into the gesture pipeline:

```
Camera → Python SDK → GestureDetector → NATS → Consumers
                                          │
                                          ├── Drone Controller
                                          ├── Robot Arm
                                          ├── Unreal Visualizer
                                          └── Flutter UI
```

The `gesture_systems/` folder contains the bridge code for connecting
these cameras to NATS messaging.
