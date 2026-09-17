# ESP32 Watch Firmware - Legacy Neural Interface

> Voice + Biometric wearable firmware for direct LiveKit connection to Legacy

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ESP32 WATCH FIRMWARE ARCHITECTURE                         │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         HARDWARE LAYER                                 │  │
│  │                                                                        │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                │  │
│  │  │   ESP32-S3   │  │   MAX30102   │  │   BMA423     │                │  │
│  │  │   (CPU)      │  │   (HR/SpO2)  │  │   (IMU)      │                │  │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘                │  │
│  │         │                 │                 │                         │  │
│  │         │   I2C Bus       │                 │                         │  │
│  │         └─────────────────┴─────────────────┘                         │  │
│  │                           │                                           │  │
│  │  ┌──────────────┐  ┌──────┴───────┐  ┌──────────────┐                │  │
│  │  │   INMP441    │  │   ST7789     │  │   Battery    │                │  │
│  │  │   (Mic)      │  │   (Display)  │  │   (LiPo)     │                │  │
│  │  └──────────────┘  └──────────────┘  └──────────────┘                │  │
│  │                                                                        │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                                    │                                         │
│                                    ▼                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        SOFTWARE LAYER                                  │  │
│  │                                                                        │  │
│  │  ┌──────────────────────────────────────────────────────────────┐    │  │
│  │  │                    MAIN APPLICATION                           │    │  │
│  │  │                                                               │    │  │
│  │  │  ┌────────────┐  ┌────────────┐  ┌────────────────────────┐ │    │  │
│  │  │  │  Sensor    │  │   UI       │  │   Communication        │ │    │  │
│  │  │  │  Manager   │  │   Manager  │  │   Manager              │ │    │  │
│  │  │  │            │  │            │  │                        │ │    │  │
│  │  │  │ • HR       │  │ • Status   │  │ • LiveKit (voice)     │ │    │  │
│  │  │  │ • Motion   │  │ • Alerts   │  │ • MQTT (biometrics)   │ │    │  │
│  │  │  │ • Battery  │  │ • Touch    │  │ • WiFi management     │ │    │  │
│  │  │  └────────────┘  └────────────┘  └────────────────────────┘ │    │  │
│  │  │                                                               │    │  │
│  │  └──────────────────────────────────────────────────────────────┘    │  │
│  │                                                                        │  │
│  │  ┌──────────────────────────────────────────────────────────────┐    │  │
│  │  │                    RTOS TASKS                                 │    │  │
│  │  │                                                               │    │  │
│  │  │  [Sensor Task]    [UI Task]    [Comm Task]    [Audio Task]  │    │  │
│  │  │      100ms           50ms          10ms           1ms        │    │  │
│  │  │                                                               │    │  │
│  │  └──────────────────────────────────────────────────────────────┘    │  │
│  │                                                                        │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                                    │                                         │
│                                    ▼                                         │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        CLOUD LAYER                                     │  │
│  │                                                                        │  │
│  │  ┌──────────────────┐              ┌──────────────────┐              │  │
│  │  │    LiveKit       │              │      EMQX        │              │  │
│  │  │    (Voice)       │              │     (MQTT)       │              │  │
│  │  │                  │              │                  │              │  │
│  │  │  WebRTC          │              │  Biometrics      │              │  │
│  │  │  <200ms          │              │  100ms           │              │  │
│  │  └────────┬─────────┘              └────────┬─────────┘              │  │
│  │           │                                 │                         │  │
│  │           └────────────────┬────────────────┘                         │  │
│  │                            ▼                                          │  │
│  │                     ┌──────────────┐                                 │  │
│  │                     │    Legacy    │                                 │  │
│  │                     │   (Agent)    │                                 │  │
│  │                     └──────────────┘                                 │  │
│  │                                                                        │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Supported Hardware

| Board | Status | Notes |
|-------|--------|-------|
| LILYGO T-Watch 2020 V3 | Primary | Built-in mic, touch display |
| LILYGO T-Watch 2020 V2 | Supported | Add external mic |
| Watchy | Supported | Add MAX30102 + mic |
| Custom ESP32-S3 | Recommended | Full sensor suite |

---

## Features

- **Voice Commands**: LiveKit SDK for real-time voice to Legacy
- **Heart Rate**: MAX30102 sensor with HRV calculation
- **SpO2 Monitoring**: Blood oxygen saturation
- **Motion Tracking**: BMA423/MPU6050 accelerometer
- **Step Counting**: On-device pedometer
- **Activity Detection**: ML classification (sitting/walking/running)
- **Fall Detection**: Emergency alerts
- **Touch UI**: Status display and quick actions

---

## Project Structure

```
firmware/esp32-watch/
├── README.md                 # This file
├── platformio.ini            # PlatformIO configuration
├── src/
│   ├── main.cpp             # Entry point
│   ├── config.h             # Configuration constants
│   ├── sensors/
│   │   ├── sensor_manager.h
│   │   ├── sensor_manager.cpp
│   │   ├── heart_rate.h
│   │   ├── heart_rate.cpp
│   │   ├── motion.h
│   │   └── motion.cpp
│   ├── communication/
│   │   ├── comm_manager.h
│   │   ├── comm_manager.cpp
│   │   ├── livekit_client.h
│   │   ├── livekit_client.cpp
│   │   ├── mqtt_client.h
│   │   └── mqtt_client.cpp
│   ├── ui/
│   │   ├── ui_manager.h
│   │   ├── ui_manager.cpp
│   │   └── screens/
│   └── utils/
│       ├── battery.h
│       └── battery.cpp
└── lib/
    └── README.md
```

---

## Quick Start

```bash
# Install PlatformIO
pip install platformio

# Clone and build
cd firmware/esp32-watch
pio run

# Upload to device
pio run -t upload

# Monitor serial output
pio device monitor
```

---

## Configuration

Edit `src/config.h`:

```cpp
// WiFi
#define WIFI_SSID "your-ssid"
#define WIFI_PASSWORD "your-password"

// LiveKit
#define LIVEKIT_URL "wss://your-livekit-server.com"
#define LIVEKIT_TOKEN "your-token"  // Or fetch from server

// MQTT (EMQX)
#define MQTT_BROKER "your-emqx-server.com"
#define MQTT_PORT 1883
#define MQTT_USER "biometric-watch"
#define MQTT_PASSWORD "your-password"
#define MQTT_DEVICE_ID "watch-001"

// Sampling rates
#define HR_SAMPLE_RATE_MS 1000     // Heart rate every 1s
#define MOTION_SAMPLE_RATE_MS 100  // Motion every 100ms
#define MQTT_PUBLISH_RATE_MS 5000  // Publish every 5s
```
