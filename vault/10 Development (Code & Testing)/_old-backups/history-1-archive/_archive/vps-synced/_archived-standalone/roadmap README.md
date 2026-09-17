# OneMind OS Upgrades

> Future integrations, extracted code, prototypes, and hardware projects.

This folder contains potential features to enhance OneMind OS capabilities - from external project extractions to hardware firmware.

---

## Status Overview

| Folder | Status | Description |
|--------|--------|-------------|
| `ada_v2_extraction/` | Ready to Integrate | CAD generation, 3D printer control, face auth |
| `concept_bytes/` | Reference Code | Spotify control, PyAutoGUI, UI concepts |
| `gesture_systems/` | Ready to Integrate | MediaPipe gestures, NATS bridge |
| `hardware/wearables/watches/` | **Working Code** | ESP32 watch firmware (PlatformIO) |
| `hardware/wearables/smart-glasses/` | Planning | Brilliant Labs Frame + ESP32 DIY |
| `hardware/drones/` | Placeholder | Future drone integration |
| `hardware/iot/` | Placeholder | Future IoT sensors |
| `hardware/robotics/` | Placeholder | Future robot control |
| `hardware/docs/` | Documentation | Depth cameras, robotics control guides |

---

## Directory Structure

```
upgrades/
├── README.md                    # This file
├── INTEGRATION_VISION.md        # How everything comes together
│
├── ada_v2_extraction/           # From nazirlouis/ada_v2 (MIT License)
│   ├── README.md                # Integration guide
│   ├── INTEGRATION_GUIDE.md     # Step-by-step integration
│   ├── backend/
│   │   ├── cad_agent.py         # 3D CAD generation (build123d)
│   │   ├── printer_agent.py     # 3D printer discovery & control
│   │   ├── authenticator.py     # Face authentication (MediaPipe)
│   │   ├── web_agent.py         # Browser automation
│   │   ├── kasa_agent.py        # TP-Link Kasa devices
│   │   ├── project_manager.py   # Project workspace management
│   │   └── hand_gesture_test.py # Gesture detection reference
│   ├── tools/                   # Agno-compatible tool wrappers
│   │   ├── cad_tools.py
│   │   ├── printer_tools.py
│   │   └── face_auth_tools.py
│   └── frontend/
│       └── gesture_notes.md
│
├── concept_bytes/               # From Concept-Bytes GitHub
│   ├── README.md                # Overview
│   ├── spotify_tools.py         # Spotify playback control
│   ├── pyautogui_fallback.py    # Desktop automation fallback
│   └── holomat_ui_concepts.md   # Radial UI inspiration
│
├── gesture_systems/             # Gesture control implementation
│   ├── README.md                # Overview + approach comparison
│   ├── kinect_pipeline.md       # Kinect + TouchDesigner + Unreal
│   ├── mediapipe_gestures.py    # Webcam-based gestures
│   └── nats_gesture_bridge.py   # Gesture → NATS integration
│
└── hardware/                    # Hardware projects & documentation
    ├── docs/                    # Hardware interface guides
    │   ├── README.md            # Hardware roadmap
    │   ├── depth_cameras.md     # Kinect, RealSense, Azure Kinect
    │   └── robotics_control.md  # Gesture → robot mappings
    │
    ├── drones/                  # Future: drone control
    ├── iot/                     # Future: IoT sensor integration
    ├── robotics/                # Future: robot arm control
    │
    └── wearables/
        ├── smart-glasses/       # AR glasses integration
        │   ├── README.md        # Planning doc (Frame + ESP32 DIY)
        │   ├── brilliant-frame/ # Brilliant Labs Frame SDK
        │   └── esp32-cam-glasses/ # DIY ESP32 glasses
        │
        └── watches/             # ESP32 Watch Firmware (WORKING)
            ├── README.md        # Architecture doc
            ├── platformio.ini   # PlatformIO project config
            ├── esp32-watches/   # Additional watch variants
            └── src/             # Firmware source code
                ├── main.cpp
                ├── config.h
                ├── config.h.template
                ├── sensors/     # Heart rate, motion
                ├── communication/ # LiveKit, MQTT
                └── ui/          # Display management
```

---

## Sources

| Source | Link | What We Got |
|--------|------|-------------|
| **ADA V2** | [github.com/nazirlouis/ada_v2](https://github.com/nazirlouis/ada_v2) | CAD generation, 3D printer control, face auth, Kasa devices |
| **JarvisV3** | [github.com/Concept-Bytes/JarvisV3](https://github.com/Concept-Bytes/JarvisV3) | OpenAI Realtime patterns, Bambu print automation |
| **Jarvis** | [github.com/Concept-Bytes/Jarvis](https://github.com/Concept-Bytes/Jarvis) | Spotify control, hotword detection |
| **HoloMat2** | [github.com/Concept-Bytes/HoloMat2](https://github.com/Concept-Bytes/HoloMat2) | Radial UI concepts for Flutter |
| **ai_plant** | [github.com/Concept-Bytes/ai_plant](https://github.com/Concept-Bytes/ai_plant) | IoT sensor + AI pattern |
| **Gesture Reel** | [instagram.com/reel/DTF-F-giHSp](https://www.instagram.com/reel/DTF-F-giHSp/?igsh=MWJibGptMHdpM2hwNw==) | Kinect + TouchDesigner + Unreal particle control |

---

## Priority Matrix

| Upgrade | Impact | Effort | Priority | Status |
|---------|--------|--------|----------|--------|
| **ESP32 Watch** | High | Done | - | Working code |
| **3D CAD Generation** | High | Medium | P1 | Ready to integrate |
| **3D Printer Control** | High | Medium | P1 | Ready to integrate |
| **Gesture Systems** | High | High | P2 | Code ready |
| **Smart Glasses** | High | High | P2 | Planning |
| **Kinect+TD+Unreal** | Very High | Very High | P2 | Research |
| **Face Auth** | Medium | Low | P3 | Code ready |
| **Spotify Control** | Low | Low | P3 | Code ready |
| **Drones** | High | High | P3 | Placeholder |
| **Robotics** | High | Very High | P4 | Placeholder |

---

## Quick Start

### ESP32 Watch (Working Now)

```bash
cd upgrades/hardware/wearables/watches
cp src/config.h.template src/config.h
# Edit config.h with your WiFi/LiveKit/MQTT settings
pio run -t upload
```

### MediaPipe Gestures

```bash
pip install mediapipe opencv-python
python upgrades/gesture_systems/mediapipe_gestures.py
```

### CAD Integration (When Ready)

```bash
pip install build123d cadquery-ocp zeroconf aiohttp
# Copy tools to backend/agno/tools/
cp upgrades/ada_v2_extraction/backend/cad_agent.py backend/agno/tools/cad/
cp upgrades/ada_v2_extraction/backend/printer_agent.py backend/agno/tools/hardware/
```

---

## Related Documentation

| Document | Location |
|----------|----------|
| Integration Vision | [INTEGRATION_VISION.md](./INTEGRATION_VISION.md) |
| ADA V2 Guide | [ada_v2_extraction/README.md](./ada_v2_extraction/README.md) |
| Gesture Systems | [gesture_systems/README.md](./gesture_systems/README.md) |
| Hardware Roadmap | [hardware/docs/README.md](./hardware/docs/README.md) |
| Watch Firmware | [hardware/wearables/watches/README.md](./hardware/wearables/watches/README.md) |
| Smart Glasses | [hardware/wearables/smart-glasses/README.md](./hardware/wearables/smart-glasses/README.md) |
