# OneMind OS Wearables Integration Plan

> Last Updated: 2026-01-22
> Status: Research Complete - Ready for Implementation

---

## Overview

This document outlines the open source wearable hardware strategy for OneMind OS, enabling Legacy to be present on glasses (vision + voice) and watches (fitness + notifications) with Bluetooth connectivity to a phone hub for LiveKit voice rooms.

---

## Smart Glasses

### 1. Brilliant Labs Frame ⭐ RECOMMENDED

**The best open source AI glasses for OneMind integration.**

| Specification | Details |
|---------------|---------|
| **Camera** | 1280x720 RGB (OV09734 sensor from Omnivision) |
| **Display** | 640x400 micro OLED, 20° field of view, optically bonded prism |
| **Microphone** | ICS-41351 MEMS (TDK), -35.5dB to 129.5dB dynamic range, 4-20kHz sample rate, 4-16 bit recording |
| **Bluetooth** | 5.3, -20dBm to 8dBm power, -95dBm sensitivity |
| **Processor** | Nordic nRF52840 - 32-bit ARM Cortex-M4F @ 64MHz |
| **Memory** | 1MB flash, 256KB RAM |
| **FPGA** | Lattice CrossLink-NX for graphics/camera acceleration |
| **Battery** | 210mAh (2x 105mAh li-ion cells) |
| **Current Draw** | 45-100mA typical operation |
| **Sensors** | 3-axis accelerometer, 3-axis e-compass (MC6470 IMU), tap detection |
| **License** | **MIT (fully open source)** |
| **Price** | ~$349 USD |
| **Website** | https://brilliant.xyz |
| **Docs** | https://docs.brilliant.xyz |
| **GitHub** | https://github.com/brilliantlabsAR |

#### SDKs Available

| SDK | Language | Install | Use Case |
|-----|----------|---------|----------|
| frame-sdk | Python | `pip install frame-sdk` | Desktop/server integration |
| frame-sdk-flutter | Dart | pub.dev | Mobile app integration |
| Lua API | Lua | On-device | Embedded scripts |
| Raw Bluetooth | Any | BLE GATT | Low-level access |

#### Python SDK Capabilities

```python
from frame_sdk import Frame

async with Frame() as frame:
    # === CAMERA ===
    # Take photo with autofocus
    image = await frame.camera.take_photo(autofocus_seconds=2)
    image_bytes = await frame.camera.take_photo(quality="high")

    # Configure autofocus mode
    await frame.camera.set_autofocus_mode("center_weighted")

    # === MICROPHONE ===
    # Record audio with silence detection
    audio = await frame.microphone.record_audio(max_length=10)

    # Record to file
    await frame.microphone.record_audio(filepath="recording.wav")

    # === DISPLAY ===
    # Show text with alignment
    await frame.display.show_text("Hello from Legacy", align="center")

    # Draw graphics
    await frame.display.draw_rect(x=10, y=10, w=100, h=50, color=1)

    # Scroll screen
    await frame.display.scroll(direction="up", pixels=20)

    # === MOTION ===
    # Detect tap
    tap = await frame.motion.wait_for_tap()

    # Get motion direction
    direction = await frame.motion.get_direction()

    # === SYSTEM ===
    # Battery level
    battery = await frame.battery.level()

    # Run Lua script
    await frame.run_lua("print('Hello from Lua')")
```

#### OneMind Integration Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frame Glasses                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │
│  │  Camera  │  │   Mic    │  │  Display │  │  Accelerometer   │ │
│  │ 1280x720 │  │  MEMS    │  │ 640x400  │  │  Tap Detection   │ │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────────┬─────────┘ │
│       │             │             │                  │           │
│       └─────────────┴─────────────┴──────────────────┘           │
│                              │                                    │
│                         Bluetooth 5.3                            │
└──────────────────────────────┬──────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                    Phone Hub (Flutter App)                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐   │
│  │ Frame BLE   │  │  LiveKit    │  │     AgnoClient          │   │
│  │  Service    │  │   Client    │  │  (SSE Streaming)        │   │
│  └──────┬──────┘  └──────┬──────┘  └───────────┬─────────────┘   │
│         │                │                      │                 │
│         └────────────────┴──────────────────────┘                 │
│                              │                                    │
│                         5G / WiFi                                 │
└──────────────────────────────┬──────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                      OneMind Cloud                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐   │
│  │   Legacy    │  │   LiveKit   │  │         NATS            │   │
│  │   Agent     │  │   Server    │  │   (Events/Commands)     │   │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘   │
└──────────────────────────────────────────────────────────────────┘
```

#### Use Cases for OneMind

| Feature | Implementation |
|---------|----------------|
| **Vision AI** | Camera → phone → Claude Vision API → Legacy response → display |
| **Voice Commands** | Mic → phone → LiveKit room → Legacy → bone conduction response |
| **Notifications** | NATS event → phone → Frame display |
| **Quick Actions** | Tap gesture → predefined command to Legacy |
| **Context Awareness** | Accelerometer → presence/activity detection |
| **Real-time Translation** | Camera OCR → translate → display overlay |
| **Meeting Notes** | Continuous mic → transcription → summary |

---

### 2. Brilliant Labs Halo (Pre-Order)

**Newer glasses with on-device AI NPU, but less camera capability.**

| Specification | Details |
|---------------|---------|
| **Display** | Micro Color OLED |
| **Audio Output** | 2x ultra-compact bone conduction speakers |
| **Microphone** | 2x with audio activity detection |
| **Processor** | Alif B1 - Cortex-M55 CPU + NPU (on-device AI) |
| **Camera** | Ultra low-power optical sensor (inference only, not full camera) |
| **Sensors** | 6-axis IMU with tap detection |
| **Bluetooth** | 5.3 |
| **Battery** | 14-hour all-day |
| **Lenses** | Optical grade with anti-reflective coating, prescription available |
| **OS** | ZephyrOS with Lua interface |
| **AI** | Noa - cloud-based multimodal AI agent |
| **License** | **Open source** |
| **Price** | TBD (pre-order) |
| **Website** | https://brilliant.xyz |

#### Halo vs Frame Comparison

| Feature | Frame | Halo |
|---------|-------|------|
| Full camera | ✅ 1280x720 | ❌ Optical sensor only |
| On-device AI | ❌ | ✅ Cortex-M55 + NPU |
| Battery life | ~4-6 hours | 14 hours |
| Speakers | ❌ | ✅ Bone conduction |
| Best for | Vision AI, photos | Voice AI, all-day wear |

**Recommendation:** Get **Frame** first for vision capabilities. Consider Halo later for always-on voice assistant use case.

---

## Smartwatches

### 1. Bangle.js 2 ⭐ RECOMMENDED

**The most developer-friendly open source smartwatch.**

| Specification | Details |
|---------------|---------|
| **Processor** | Nordic nRF52840 ARM Cortex-M4 @ 64MHz |
| **Memory** | 256KB RAM, 1024KB on-chip flash, 8MB external flash |
| **Display** | 1.3" 176x176 always-on 3-bit color LCD, full touchscreen |
| **GPS** | ✅ GPS/GLONASS (AT6558/AT6558R) |
| **Heart Rate** | ✅ Vcare VC31/VC31B optical sensor |
| **Accelerometer** | ✅ 3-axis |
| **Magnetometer** | ✅ 3-axis compass |
| **Barometer** | ✅ Pressure + temperature sensor |
| **Bluetooth** | BLE (nRF52840) |
| **Battery** | 175mAh, ~4 weeks standby |
| **Water Resistance** | IP67 (everyday use, not swimming) |
| **Dimensions** | 36mm x 43mm x 12mm |
| **Strap** | Standard 20mm watch straps |
| **Price** | **£99.60 (~$125 USD)** |
| **Bulk Pricing** | £87.60 at 50+ units |
| **In Stock** | ✅ 352 units available |
| **Website** | https://banglejs.com |
| **Shop** | https://shop.espruino.com/banglejs2 |
| **Docs** | https://www.espruino.com/Bangle.js2 |

#### Developer Features

| Feature | Details |
|---------|---------|
| **Language** | JavaScript (Espruino) |
| **Visual Programming** | Blockly support |
| **App Store** | 600+ free apps |
| **App Install** | Wireless via web browser |
| **APIs** | Graphics, accelerometer gestures, Bluetooth |
| **AI Support** | Accelerometer gesture recognition with ML |
| **Companion Apps** | Gadgetbridge (Android), multiple iOS options |

#### JavaScript API Examples

```javascript
// === HEART RATE ===
Bangle.on('HRM', function(hrm) {
  console.log("BPM: " + hrm.bpm);
  console.log("Confidence: " + hrm.confidence);
});
Bangle.setHRMPower(1); // Turn on

// === GPS ===
Bangle.on('GPS', function(gps) {
  console.log("Lat: " + gps.lat);
  console.log("Lon: " + gps.lon);
  console.log("Speed: " + gps.speed);
});
Bangle.setGPSPower(1); // Turn on

// === ACCELEROMETER ===
Bangle.on('accel', function(acc) {
  console.log("X: " + acc.x);
  console.log("Y: " + acc.y);
  console.log("Z: " + acc.z);
});

// === STEP COUNTER ===
Bangle.on('step', function(steps) {
  console.log("Steps: " + steps);
});

// === GESTURES ===
Bangle.on('gesture', function(gesture) {
  // gesture is an array of accelerometer data
  // Can be classified with ML
});

// === COMPASS ===
Bangle.on('mag', function(mag) {
  console.log("Heading: " + mag.heading);
});

// === DISPLAY ===
g.clear();
g.setFont("6x8", 2);
g.drawString("Legacy says:", 10, 30);
g.drawString("Meeting in 5 min", 10, 50);
g.flip();

// === BLUETOOTH ===
// Send data to phone
Bluetooth.println(JSON.stringify({
  type: "health",
  bpm: 72,
  steps: 5432,
  location: {lat: 37.7749, lon: -122.4194}
}));

// === NOTIFICATIONS ===
E.showMessage("New message from Legacy", "OneMind");

// === VIBRATION ===
Bangle.buzz(200); // 200ms vibration
```

#### OneMind Integration Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Bangle.js 2                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │
│  │   GPS    │  │  Heart   │  │  Accel   │  │    Display       │ │
│  │ GLONASS  │  │   Rate   │  │ + Gyro   │  │   176x176        │ │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────────┬─────────┘ │
│       │             │             │                  │           │
│       └─────────────┴─────────────┴──────────────────┘           │
│                              │                                    │
│                    Bluetooth LE + JSON                           │
└──────────────────────────────┬──────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                    Phone Hub (Flutter App)                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐   │
│  │ Bangle BLE  │  │    MQTT     │  │     AgnoClient          │   │
│  │  Service    │  │   Client    │  │  (SSE Streaming)        │   │
│  └──────┬──────┘  └──────┬──────┘  └───────────┬─────────────┘   │
│         │                │                      │                 │
│         └────────────────┴──────────────────────┘                 │
│                              │                                    │
│                         5G / WiFi                                 │
└──────────────────────────────┬──────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                      OneMind Cloud                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐   │
│  │   Legacy    │  │    NATS     │  │      TimescaleDB        │   │
│  │   Agent     │  │ (MQTT:1883) │  │   (Health History)      │   │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘   │
└──────────────────────────────────────────────────────────────────┘
```

#### Data Flow: Watch → Legacy

```
Watch Sensors → BLE JSON → Phone App → MQTT → NATS → Legacy Agent
                                                         │
                                              ┌──────────┴──────────┐
                                              │                     │
                                              ▼                     ▼
                                        Store in DB          React to events
                                       (TimescaleDB)       (e.g., high heart rate)
```

#### MQTT Topics for Watch Data

| Topic | Data | Frequency |
|-------|------|-----------|
| `onemind/watch/health/heartrate` | `{"bpm": 72, "confidence": 95}` | Every 5s when active |
| `onemind/watch/health/steps` | `{"steps": 5432, "goal": 10000}` | Every minute |
| `onemind/watch/location/gps` | `{"lat": 37.77, "lon": -122.41, "speed": 0}` | Every 30s when moving |
| `onemind/watch/activity/gesture` | `{"type": "raise_wrist"}` | On event |
| `onemind/watch/battery` | `{"level": 85, "charging": false}` | Every 5 minutes |

#### Use Cases for OneMind

| Feature | Implementation |
|---------|----------------|
| **Fitness Tracking** | HR + steps + GPS → TimescaleDB → Dashboard |
| **Health Alerts** | High HR → NATS → Legacy → notification |
| **Location Presence** | GPS → physical presence system |
| **Activity Detection** | Accelerometer → sleeping/walking/running |
| **Quick Commands** | Gesture → predefined action |
| **Notifications** | Legacy response → watch display + vibration |
| **Meeting Reminders** | Calendar event → watch alert |

---

### 2. PineTime (Budget Option)

**Cheapest open source smartwatch, but limited capabilities.**

| Specification | Details |
|---------------|---------|
| **Processor** | Nordic nRF52832 |
| **Memory** | 64KB RAM, 512KB flash + 4MB SPI NOR |
| **Display** | 1.3" 240x240 IPS capacitive touch, RGB 65K colors (ST7789) |
| **Heart Rate** | ✅ HRS3300 optical sensor |
| **Accelerometer** | ✅ Bosch BMA421/BMA425 |
| **GPS** | ❌ None |
| **Bluetooth** | 5.0 BLE |
| **Battery** | 170-180mAh 3.8V LiPo, ~1 week |
| **Water Resistance** | IP67 (sealed edition) |
| **Dimensions** | 37.5 x 40 x 11mm, 38g |
| **Price** | **$26.99 USD** |
| **Availability** | ⚠️ Currently out of stock |
| **Website** | https://pine64.org/pinetime |
| **Wiki** | https://wiki.pine64.org/wiki/PineTime |

#### Firmware Options

| Firmware | Language | Features |
|----------|----------|----------|
| **InfiniTime** (default) | C++ | Full-featured, stable |
| **Wasp-OS** | MicroPython | Python scripting |
| **Custom** | Any | Full hardware access |

#### Companion Apps

| App | Platform | Features |
|-----|----------|----------|
| Gadgetbridge | Android | Notifications, music control, fitness |
| Amazfish | Linux/SailfishOS | Full sync |
| InfiniLink | iOS | Basic sync |
| ITD | Linux CLI | Terminal control |
| Siglo | Linux GTK | Desktop sync |

#### PineTime vs Bangle.js 2 Comparison

| Feature | PineTime | Bangle.js 2 |
|---------|----------|-------------|
| Price | $27 | $125 |
| RAM | 64KB | 256KB |
| Flash | 512KB | 1MB + 8MB |
| GPS | ❌ | ✅ |
| Display | 240x240 | 176x176 |
| Always-on | ❌ | ✅ |
| JavaScript | ❌ | ✅ |
| Developer-friendly | Medium | High |
| Availability | Out of stock | In stock |

**Recommendation:** Get **Bangle.js 2** unless budget is extremely tight. The GPS and better developer experience are worth the extra $100.

---

## 5G/LTE Connectivity Reality

### Current State

**There are NO open source 5G/LTE smartwatches available.**

All cellular-capable smartwatches are closed ecosystems:

| Watch | OS | Openness | Notes |
|-------|-----|----------|-------|
| Apple Watch Ultra | watchOS | Closed | No custom apps without Apple approval |
| Samsung Galaxy Watch | WearOS/Tizen | Partially open | Limited to Samsung ecosystem |
| Google Pixel Watch | WearOS | Partially open | Google services required |
| Huawei Watch | HarmonyOS | Closed | No western app support |

### Workaround: Phone as Hub

The recommended architecture uses the phone as the 5G/LTE hub:

```
┌─────────────┐     BLE      ┌──────────────────┐     5G/WiFi     ┌────────────────┐
│  Bangle.js  │◄────────────►│   Phone (Hub)    │◄───────────────►│  OneMind Cloud │
│  Watch      │              │                  │                 │                │
└─────────────┘              │  - 5G/LTE radio  │                 │  - Legacy      │
                             │  - LiveKit       │                 │  - NATS        │
┌─────────────┐     BLE      │  - MQTT bridge   │                 │  - API         │
│   Frame     │◄────────────►│  - BLE central   │                 │                │
│  Glasses    │              │                  │                 │                │
└─────────────┘              └──────────────────┘                 └────────────────┘
```

**Benefits of Phone Hub:**
- Always have 5G connectivity via phone
- No battery drain on wearables from cellular radio
- Single point of integration for all wearables
- LiveKit voice works through phone speakers/mic or Frame
- Can use existing phone (no extra hardware cost)

---

## Complete Hardware Stack

### Recommended Setup

| Device | Purpose | Price | Priority |
|--------|---------|-------|----------|
| **Brilliant Labs Frame** | Vision AI + Voice | $349 | ⭐ High |
| **Bangle.js 2** | Fitness + GPS + Notifications | $125 | ⭐ High |
| **Existing Phone** | 5G Hub + LiveKit + BLE Central | $0 | Required |
| **Brilliant Labs Halo** | Always-on voice (future) | TBD | Optional |
| **PineTime** | Backup/testing | $27 | Optional |

**Total Investment: ~$475 USD** (Frame + Bangle.js 2)

### Full Stack Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              WEARABLE LAYER                                      │
│                                                                                  │
│  ┌───────────────────────┐              ┌───────────────────────┐               │
│  │    Frame Glasses      │              │     Bangle.js 2       │               │
│  │  ┌─────┐ ┌─────┐     │              │  ┌─────┐ ┌─────┐     │               │
│  │  │Cam  │ │Mic  │     │              │  │GPS  │ │ HR  │     │               │
│  │  │720p │ │MEMS │     │              │  │GNSS │ │VC31 │     │               │
│  │  └─────┘ └─────┘     │              │  └─────┘ └─────┘     │               │
│  │  ┌─────┐ ┌─────┐     │              │  ┌─────┐ ┌─────┐     │               │
│  │  │OLED │ │IMU  │     │              │  │Accel│ │Baro │     │               │
│  │  │640x │ │Tap  │     │              │  │3-ax │ │Temp │     │               │
│  │  └─────┘ └─────┘     │              │  └─────┘ └─────┘     │               │
│  │     BLE 5.3          │              │     BLE (nRF52840)   │               │
│  └──────────┬───────────┘              └──────────┬───────────┘               │
│             │                                      │                           │
└─────────────┼──────────────────────────────────────┼───────────────────────────┘
              │              BLE                     │
              └──────────────┬───────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                               HUB LAYER (Phone)                                  │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐    │
│  │                     Flutter OneMind App                                  │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │    │
│  │  │ Frame BLE   │  │ Watch BLE   │  │  LiveKit    │  │ AgnoClient  │    │    │
│  │  │  Service    │  │  Service    │  │   Client    │  │    SSE      │    │    │
│  │  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘    │    │
│  │         │                │                │                │           │    │
│  │         └────────────────┴────────────────┴────────────────┘           │    │
│  │                                   │                                     │    │
│  │                          State Management                               │    │
│  │                           (Riverpod)                                    │    │
│  └─────────────────────────────────────────────────────────────────────────┘    │
│                                      │                                          │
│                               5G / WiFi / LTE                                   │
└──────────────────────────────────────┬──────────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              CLOUD LAYER                                         │
│                                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  Legacy Agent   │  │     NATS        │  │    LiveKit      │                  │
│  │  (Agno + Claude)│  │  (MQTT:1883)    │  │    Server       │                  │
│  │                 │  │  (WS:9222)      │  │                 │                  │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘                  │
│           │                    │                    │                           │
│           └────────────────────┴────────────────────┘                           │
│                                │                                                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐                  │
│  │   PostgreSQL    │  │  TimescaleDB    │  │     Redis       │                  │
│  │  (User Data)    │  │ (Health History)│  │  (Hot State)    │                  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘                  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Flutter App Integration

### Existing Scaffolded Screens

The Flutter app (`mobile/`) already has placeholder screens for wearables:

| Screen | Path | Purpose |
|--------|------|---------|
| Frame | `lib/screens/wearables/frame_screen.dart` | Frame glasses connection + control |
| Watch | `lib/screens/wearables/watch_screen.dart` | Watch connection + health data |
| Voice | `lib/screens/wearables/voice_screen.dart` | LiveKit voice rooms |

### Required Flutter Packages

```yaml
# pubspec.yaml additions
dependencies:
  # BLE for Frame and Watch
  flutter_blue_plus: ^1.31.0

  # Frame SDK (when available for Flutter)
  # frame_sdk_flutter: ^x.x.x

  # MQTT for watch data bridge
  mqtt_client: ^10.0.0

  # LiveKit for voice
  livekit_client: ^2.1.0

  # Health data
  health: ^10.0.0

  # Background services
  flutter_background_service: ^5.0.0
```

### BLE Service Architecture

```dart
// lib/services/ble_service.dart

abstract class WearableDevice {
  String get name;
  String get id;
  bool get isConnected;

  Future<void> connect();
  Future<void> disconnect();
  Stream<Map<String, dynamic>> get dataStream;
}

class FrameDevice implements WearableDevice {
  // Frame-specific BLE implementation
  // Camera capture
  // Audio recording
  // Display control
}

class BangleDevice implements WearableDevice {
  // Bangle.js BLE implementation
  // Health data streaming
  // GPS location
  // Notifications
}

class BleService {
  final List<WearableDevice> devices = [];

  Future<void> scanForDevices();
  Future<void> connectToFrame();
  Future<void> connectToBangle();

  // Unified data stream from all wearables
  Stream<WearableEvent> get eventStream;
}
```

### MQTT Bridge for Watch Data

```dart
// lib/services/mqtt_service.dart

class MqttService {
  final String brokerUrl = 'wss://mqtt.onemindos.dev:9222';

  Future<void> connect();

  // Publish watch data to NATS via MQTT
  void publishHealthData(HealthData data) {
    client.publish(
      'onemind/watch/health/heartrate',
      jsonEncode({'bpm': data.bpm, 'confidence': data.confidence}),
    );
  }

  void publishLocation(Location location) {
    client.publish(
      'onemind/watch/location/gps',
      jsonEncode({'lat': location.lat, 'lon': location.lon}),
    );
  }

  // Subscribe to Legacy responses for watch
  Stream<WatchNotification> get notifications;
}
```

---

## Implementation Roadmap

### Phase 1: Frame Integration (Priority)

| Task | Description | Effort |
|------|-------------|--------|
| 1.1 | Add flutter_blue_plus package | Small |
| 1.2 | Implement Frame BLE service | Medium |
| 1.3 | Camera capture → Vision AI | Medium |
| 1.4 | Audio recording → LiveKit | Medium |
| 1.5 | Display notifications from Legacy | Small |
| 1.6 | Tap gesture → quick commands | Small |

### Phase 2: Watch Integration

| Task | Description | Effort |
|------|-------------|--------|
| 2.1 | Implement Bangle.js BLE service | Medium |
| 2.2 | Heart rate streaming | Small |
| 2.3 | GPS location tracking | Small |
| 2.4 | Step counter integration | Small |
| 2.5 | MQTT bridge to NATS | Medium |
| 2.6 | Push notifications to watch | Medium |

### Phase 3: Unified Experience

| Task | Description | Effort |
|------|-------------|--------|
| 3.1 | Combined wearable dashboard | Medium |
| 3.2 | Health data visualization | Medium |
| 3.3 | Cross-device presence | Large |
| 3.4 | Voice across all devices | Large |

---

## Resources

### Brilliant Labs Frame
- Website: https://brilliant.xyz
- Documentation: https://docs.brilliant.xyz
- GitHub: https://github.com/brilliantlabsAR
- Discord: https://discord.gg/brilliantlabs
- Python SDK: `pip install frame-sdk`

### Bangle.js 2
- Website: https://banglejs.com
- Shop: https://shop.espruino.com/banglejs2
- Documentation: https://www.espruino.com/Bangle.js2
- App Store: https://banglejs.com/apps/
- GitHub: https://github.com/espruino/BangleApps

### PineTime
- Website: https://pine64.org/pinetime
- Wiki: https://wiki.pine64.org/wiki/PineTime
- InfiniTime: https://github.com/InfiniTimeOrg/InfiniTime
- Wasp-OS: https://github.com/wasp-os/wasp-os

### Flutter Packages
- flutter_blue_plus: https://pub.dev/packages/flutter_blue_plus
- mqtt_client: https://pub.dev/packages/mqtt_client
- livekit_client: https://pub.dev/packages/livekit_client

---

## Order Links

| Device | Link | Price |
|--------|------|-------|
| **Frame Glasses** | https://brilliant.xyz | ~$349 |
| **Bangle.js 2** | https://shop.espruino.com/banglejs2 | £99.60 (~$125) |
| **PineTime Dev Kit** | https://pine64.com/product/pinetime-dev-kit/ | $26.99 (out of stock) |

---

## Appendix: Alternative Glasses (Not Recommended)

These glasses were considered but not recommended for OneMind:

| Glasses | Why Not |
|---------|---------|
| **Meta Ray-Ban** | Closed ecosystem, Meta account required, no SDK |
| **Xreal Air** | Display-only (no camera), closed SDK |
| **Even Realities G1** | Limited SDK, not fully open source |
| **Vuzix Blade** | Enterprise-focused, expensive ($1000+), limited SDK |
| **Google Glass Enterprise** | Discontinued, enterprise-only |
| **Apple Vision Pro** | $3500, closed ecosystem, not glasses form factor |

**Frame is the only fully open source smart glasses with camera + display + mic + SDK access.**
