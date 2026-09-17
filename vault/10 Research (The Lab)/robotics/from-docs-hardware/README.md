# Hardware Interfaces

## 🔗 Sources & References

| Hardware                 | Vendor/Docs                                                                                                    | Price   | Notes                          |
| ------------------------ | -------------------------------------------------------------------------------------------------------------- | ------- | ------------------------------ |
| **Azure Kinect DK**      | [Microsoft Store](https://www.microsoft.com/en-us/d/azure-kinect-dk/8pp5vxmd9nhq)                              | $400    | Recommended depth camera       |
| **Intel RealSense D435** | [Intel RealSense](https://www.intel.com/content/www/us/en/architecture-and-technology/realsense-overview.html) | $300    | Portable USB-C depth camera    |
| **Kinect v2**            | eBay/Used Market                                                                                               | $50-100 | Budget option (discontinued)   |
| **DJI Tello**            | [DJI Store](https://store.dji.com/product/tello)                                                               | $100    | Entry-level programmable drone |
| **Leap Motion 2**        | [Ultraleap](https://www.ultraleap.com/)                                                                        | $140    | Precision hand tracking        |
| **python-kasa**          | [Docs](https://python-kasa.readthedocs.io/)                                                                    | N/A     | TP-Link Kasa device control    |
| **STEMMA Sensors**       | [Adafruit](https://www.adafruit.com/category/1005)                                                             | $5-30   | I2C sensors for IoT            |

**Documentation Date:** 2026-01-30

---

## 📋 Overview

This folder contains documentation and code for integrating physical devices into OneMind OS:

- **Depth cameras** (Kinect, RealSense) - For gesture control
- **Robotics platforms** - Drones, robot arms, mobile robots
- **3D printers** - Creality, Bambu, Prusa
- **IoT sensors** - Environmental, proximity, etc.
- **Smart devices** - TP-Link Kasa (alternative to Home Assistant)

---

## 📁 Contents

| File                  | Purpose                                       | OneMind Use Case                           |
| --------------------- | --------------------------------------------- | ------------------------------------------ |
| `depth_cameras.md`    | Comparison of Kinect, RealSense, Azure Kinect | Choose right depth camera for your needs   |
| `robotics_control.md` | Gesture-to-robot mappings with safety rules   | Control drones and robot arms via gestures |

---

## 🗺️ Hardware Roadmap

### Phase 1: Input Devices ✅ (Partial)

- [x] Microphones (via LiveKit/Deepgram)
- [x] Webcam (via MediaPipe gestures)
- [ ] Kinect/Azure Kinect
- [ ] Intel RealSense
- [ ] Leap Motion

### Phase 2: Output Devices

- [ ] 3D Printers (Creality, Bambu)
- [ ] Smart Home (via Home Assistant)
- [ ] Display arrays
- [ ] Spatial audio

### Phase 3: Robotics

- [ ] Consumer drones (DJI Tello, Ryze)
- [ ] Robot arms (servo-based)
- [ ] Mobile robots
- [ ] Boston Dynamics Spot (stretch goal)

---

## 🎯 Integration Priority

| Device                 | Priority | Reason                                  |
| ---------------------- | -------- | --------------------------------------- |
| **Webcam + MediaPipe** | P1       | No additional hardware, immediate value |
| **3D Printer**         | P1       | CAD → Print workflow complete           |
| **Azure Kinect**       | P2       | Professional gesture control            |
| **DJI Tello Drone**    | P2       | Gesture-controlled drone demo           |
| **TP-Link Kasa**       | P3       | Alternative to Home Assistant           |
| **Spot Robot**         | P4       | Stretch goal, high cost                 |

---

## 💡 Key Insight

All hardware integrates through NATS messaging:

```
Hardware → Python Driver → NATS Topic → OneMind Consumers
   │                          │                │
 Kinect              gesture.classified    Drone Controller
 Drone                drone.telemetry      Unreal Visualizer
 Printer            print.status           Flutter UI
```

This decoupled architecture means adding new hardware is just:

1. Write a Python driver
2. Publish to appropriate NATS topic
3. Existing consumers automatically work
