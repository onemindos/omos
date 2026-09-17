# Biometric Awareness Architecture

> Legacy's Physiological Sensing Layer - The "Neural Interface"

---

## The Vision: Real-World Neural Interface

In Halo, Cortana connects to Master Chief via a neural interface implanted in his skull, giving her direct access to his vital signs, stress levels, and physical state. We achieve the same with **distributed wearable sensors**.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      CORTANA'S NEURAL INTERFACE (Halo)                       │
│                                                                              │
│   ┌─────────────┐     Direct Neural Link      ┌─────────────────────┐       │
│   │   CORTANA   │ ◄────────────────────────── │   Master Chief's    │       │
│   │    (AI)     │    < 1ms latency            │   Brain/Body        │       │
│   └─────────────┘    Full biometric access    └─────────────────────┘       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                    ║
                                    ║ Real-World Implementation
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Legacy'S DISTRIBUTED NEURAL INTERFACE                     │
│                                                                              │
│   ┌─────────────┐    NATS + LiveKit    ┌────────────────────────────────┐   │
│   │   Legacy    │ ◄────────────────────│   Wearable Sensor Mesh         │   │
│   │    (AI)     │    < 500ms latency   │                                │   │
│   └─────────────┘    Full biometric    │  ┌──────┐ ┌──────┐ ┌──────┐   │   │
│                      access            │  │Watch │ │Ring  │ │Glasses│  │   │
│                                        │  └──────┘ └──────┘ └──────┘   │   │
│                                        └────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## ESP32 + LiveKit: The Game Changer

LiveKit just released an SDK for ESP32, enabling **sub-$10 chips to do real-time voice AI**.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ESP32 + LIVEKIT ARCHITECTURE                         │
│                                                                              │
│                              WHAT THIS ENABLES                               │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        YOUR CUSTOM WEARABLE                          │    │
│  │                                                                      │    │
│  │   ┌──────────────┐   ┌──────────────┐   ┌──────────────┐            │    │
│  │   │   ESP32-S3   │   │   Sensors    │   │   LiveKit    │            │    │
│  │   │   or P4      │ + │   (HR/SpO2/  │ + │   SDK        │            │    │
│  │   │   (~$5-10)   │   │   Accel/Temp)│   │   (Built-in) │            │    │
│  │   └──────────────┘   └──────────────┘   └──────────────┘            │    │
│  │                                                                      │    │
│  │   Result: Voice AI wearable that talks directly to Legacy           │    │
│  │           with real-time biometric streaming                        │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│                                    │                                         │
│                                    ▼                                         │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         DATA FLOW                                    │    │
│  │                                                                      │    │
│  │   Wearable ──► LiveKit Room ──► Legacy Agent ──► Consciousness      │    │
│  │                                                                      │    │
│  │   Voice:     WebRTC audio      Real-time         patterns.py        │    │
│  │   Biometrics: Data channel     processing        temporal.py        │    │
│  │   Latency:   < 200ms           < 100ms           < 50ms             │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### LiveKit ESP32 SDK Features

| Feature | Status | Use Case |
|---------|--------|----------|
| Voice streaming | ✅ Ready | "Hey Legacy, how's my heart rate?" |
| Audio publishing | ✅ Ready | Send voice commands |
| Audio subscribing | ✅ Ready | Hear Legacy's responses |
| Data channels | ✅ Ready | Stream sensor data in real-time |
| Video (ESP32-P4) | Coming | Camera feed from glasses |

---

## Open Source Smartwatch Comparison

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      OPEN SOURCE SMARTWATCH OPTIONS                          │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                          PINETIME ($27)                               │   │
│  │                                                                       │   │
│  │  Processor: nRF52832 (ARM Cortex-M4)                                 │   │
│  │  Sensors:   ✅ Heart Rate (HRS3300)                                  │   │
│  │             ✅ Accelerometer (BMA425)                                │   │
│  │  Connectivity: Bluetooth 5.0 LE                                      │   │
│  │  Battery: 170mAh (~1 week)                                          │   │
│  │  Display: 240x240 IPS touch                                          │   │
│  │  Firmware: InfiniTime (open source)                                  │   │
│  │                                                                       │   │
│  │  PROS: Built-in HR, great battery, mature firmware                   │   │
│  │  CONS: No WiFi, no direct LiveKit (needs BLE bridge)                │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                          WATCHY ($50)                                 │   │
│  │                                                                       │   │
│  │  Processor: ESP32-PICO-D4 ⭐ LiveKit Compatible                      │   │
│  │  Sensors:   ❌ No Heart Rate                                         │   │
│  │             ✅ Accelerometer (BMA423)                                │   │
│  │  Connectivity: WiFi + Bluetooth LE                                   │   │
│  │  Battery: 200mAh (E-ink = weeks)                                    │   │
│  │  Display: 200x200 E-Paper                                           │   │
│  │  Firmware: Arduino/MicroPython/ESP-IDF                              │   │
│  │                                                                       │   │
│  │  PROS: ESP32 = direct LiveKit, WiFi, hackable                       │   │
│  │  CONS: No HR sensor (must add external)                             │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                     LILYGO T-WATCH 2020 V3 ($45)                      │   │
│  │                                                                       │   │
│  │  Processor: ESP32-D0WDQ6 ⭐ LiveKit Compatible                       │   │
│  │  Sensors:   ❌ No Heart Rate (add MAX30102)                          │   │
│  │             ✅ Accelerometer (BMA423)                                │   │
│  │             ✅ Microphone (SPM1423) ⭐                               │   │
│  │             ✅ Infrared transmitter                                  │   │
│  │  Connectivity: WiFi + Bluetooth                                      │   │
│  │  Battery: 380mAh                                                    │   │
│  │  Display: 240x240 IPS touch                                         │   │
│  │  Firmware: Arduino (TTGO library)                                   │   │
│  │                                                                       │   │
│  │  PROS: Built-in mic for voice, IR for home control, ESP32          │   │
│  │  CONS: No HR (easily added via I2C)                                 │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                     CUSTOM BUILD (Recommended)                        │   │
│  │                                                                       │   │
│  │  Processor: ESP32-S3 or P4 ⭐ Latest, best LiveKit support          │   │
│  │  Sensors:   ✅ MAX30102 (HR + SpO2) - $3                            │   │
│  │             ✅ MPU6050 (Accel + Gyro) - $2                           │   │
│  │             ✅ MAX30205 (Body Temp) - $4                             │   │
│  │             ✅ INMP441 (Mic) - $2                                    │   │
│  │  Connectivity: WiFi + Bluetooth 5.0                                  │   │
│  │  Cost: ~$20-30 total                                                │   │
│  │                                                                       │   │
│  │  PROS: Full sensor suite, latest chip, LiveKit native               │   │
│  │  CONS: Requires custom PCB/enclosure design                         │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Recommendation Matrix

| Your Goal | Best Option | Why |
|-----------|-------------|-----|
| Quick start, HR needed | PineTime + BLE Bridge | Works now, $27 |
| Voice commands from wrist | LILYGO T-Watch V3 | Built-in mic, ESP32 |
| Maximum hackability | Watchy + MAX30102 | Full ESP32, add sensors |
| Ultimate biometric suite | Custom ESP32-S3 Build | All sensors, LiveKit native |

---

## Complete Biometric Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Legacy BIOMETRIC AWARENESS SYSTEM                         │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        SENSOR LAYER                                  │    │
│  │                                                                      │    │
│  │   WRIST                    FINGER                   HEAD            │    │
│  │   ┌─────────────┐         ┌─────────────┐         ┌─────────────┐   │    │
│  │   │ ESP32 Watch │         │ Oura Ring   │         │ Brilliant   │   │    │
│  │   │             │         │             │         │ Labs Halo   │   │    │
│  │   │ • Heart Rate│         │ • HRV       │         │             │   │    │
│  │   │ • SpO2      │         │ • Sleep     │         │ • Vision    │   │    │
│  │   │ • Motion    │         │ • Temp      │         │ • Audio     │   │    │
│  │   │ • Voice Cmd │         │ • Activity  │         │ • Context   │   │    │
│  │   └──────┬──────┘         └──────┬──────┘         └──────┬──────┘   │    │
│  │          │                       │                       │          │    │
│  │          │ LiveKit               │ BLE/API               │ LiveKit  │    │
│  │          │ (direct)              │ (bridge)              │ (direct) │    │
│  └──────────┼───────────────────────┼───────────────────────┼──────────┘    │
│             │                       │                       │               │
│             └───────────────────────┼───────────────────────┘               │
│                                     ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      TRANSPORT LAYER                                 │    │
│  │                                                                      │    │
│  │   ┌─────────────────┐              ┌─────────────────┐              │    │
│  │   │    LIVEKIT      │              │   NATS + EMQX   │              │    │
│  │   │                 │              │                 │              │    │
│  │   │ • Voice streams │              │ • Sensor data   │              │    │
│  │   │ • Video feeds   │              │ • Events        │              │    │
│  │   │ • Real-time     │              │ • Persistence   │              │    │
│  │   │   data channels │              │ • IoT bridge    │              │    │
│  │   └────────┬────────┘              └────────┬────────┘              │    │
│  │            │                                │                        │    │
│  │            └────────────────┬───────────────┘                        │    │
│  │                             ▼                                        │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                │                                             │
│  ┌─────────────────────────────▼───────────────────────────────────────┐    │
│  │                    CONSCIOUSNESS LAYER                               │    │
│  │                                                                      │    │
│  │   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │    │
│  │   │  bus.py      │  │ archiver.py  │  │ patterns.py  │              │    │
│  │   │              │  │              │  │              │              │    │
│  │   │ Event ingest │─►│ Store in     │─►│ Detect:      │              │    │
│  │   │ < 1ms        │  │ TimescaleDB  │  │ • Stress     │              │    │
│  │   │              │  │              │  │ • Fatigue    │              │    │
│  │   └──────────────┘  └──────────────┘  │ • Anomalies  │              │    │
│  │                                       └──────┬───────┘              │    │
│  │                                              │                       │    │
│  │   ┌──────────────┐  ┌──────────────┐        ▼                       │    │
│  │   │ temporal.py  │  │ notifier.py  │  ┌──────────────┐              │    │
│  │   │              │  │              │  │ presence.py  │              │    │
│  │   │ Time context │  │ Alert user   │◄─│              │              │    │
│  │   │ + deadlines  │  │ via best     │  │ Device state │              │    │
│  │   │              │  │ channel      │  │ + location   │              │    │
│  │   └──────────────┘  └──────────────┘  └──────────────┘              │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                │                                             │
│  ┌─────────────────────────────▼───────────────────────────────────────┐    │
│  │                       Legacy AGENT                                   │    │
│  │                                                                      │    │
│  │   ┌──────────────────────────────────────────────────────────────┐  │    │
│  │   │                    BIOMETRIC CONTEXT                          │  │    │
│  │   │                                                               │  │    │
│  │   │  "Zeus's heart rate is 95bpm (elevated from baseline 68).    │  │    │
│  │   │   HRV is low (32ms vs usual 45ms). He's been at desk for     │  │    │
│  │   │   3.5 hours without movement. SpO2 normal at 98%.            │  │    │
│  │   │   Pattern match: Pre-stress state, similar to Jan 15th       │  │    │
│  │   │   before the investor call."                                  │  │    │
│  │   │                                                               │  │    │
│  │   │  RECOMMENDED ACTIONS:                                         │  │    │
│  │   │  • Suggest 5-minute walk                                      │  │    │
│  │   │  • Defer non-urgent notifications                             │  │    │
│  │   │  • Prepare calming music if requested                         │  │    │
│  │   └──────────────────────────────────────────────────────────────┘  │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Biometric Data Flow (Timing)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         BIOMETRIC DATA LATENCY                               │
│                                                                              │
│  TIER 1: REFLEXIVE (< 50ms) - Immediate Safety                              │
│  ════════════════════════════════════════════                               │
│                                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │
│  │   Sensor    │    │   ESP32     │    │   LiveKit   │    │   Action    │  │
│  │   Spike     │───►│   Process   │───►│   Transmit  │───►│   Execute   │  │
│  │             │    │             │    │             │    │             │  │
│  │   5ms       │    │   10ms      │    │   20ms      │    │   15ms      │  │
│  └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘  │
│                                                                              │
│  Use Case: Heart rate spikes to 180+ → Immediate alert                      │
│  Use Case: Fall detected → Emergency notification                           │
│                                                                              │
│  ───────────────────────────────────────────────────────────────────────    │
│                                                                              │
│  TIER 2: REACTIVE (50ms - 500ms) - Context Updates                          │
│  ══════════════════════════════════════════════════                         │
│                                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │
│  │   Sensor    │    │   NATS      │    │  Archiver   │    │  Patterns   │  │
│  │   Reading   │───►│   Publish   │───►│   Store     │───►│   Analyze   │  │
│  │             │    │             │    │             │    │             │  │
│  │   100ms     │    │   1ms       │    │   50ms      │    │   200ms     │  │
│  │  (sampled)  │    │             │    │             │    │             │  │
│  └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘  │
│                                                                              │
│  Use Case: HR trending up over 5 min → Pattern detection                    │
│  Use Case: Activity level low for 2 hours → Movement suggestion             │
│                                                                              │
│  ───────────────────────────────────────────────────────────────────────    │
│                                                                              │
│  TIER 3: COGNITIVE (500ms - 5s) - Agent Reasoning                           │
│  ═════════════════════════════════════════════════                          │
│                                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │
│  │  Pattern    │    │   Agent     │    │    LLM      │    │  Response   │  │
│  │  Detected   │───►│   Context   │───►│   Reason    │───►│  Generate   │  │
│  │             │    │   Build     │    │             │    │             │  │
│  │   instant   │    │   200ms     │    │   2-4s      │    │   500ms     │  │
│  └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘  │
│                                                                              │
│  Use Case: "You seem stressed, want me to reschedule the 3pm call?"        │
│  Use Case: "Your sleep was poor last night, should I block focus time?"    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Sensor Integration Details

### Heart Rate + SpO2 (MAX30102)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     MAX30102 INTEGRATION (HR + SpO2)                         │
│                                                                              │
│  HARDWARE CONNECTION                                                         │
│  ═══════════════════                                                        │
│                                                                              │
│   MAX30102              ESP32-S3                                            │
│  ┌─────────┐           ┌─────────┐                                          │
│  │ VCC     │───────────│ 3.3V    │                                          │
│  │ GND     │───────────│ GND     │                                          │
│  │ SDA     │───────────│ GPIO21  │  (I2C Data)                              │
│  │ SCL     │───────────│ GPIO22  │  (I2C Clock)                             │
│  │ INT     │───────────│ GPIO19  │  (Interrupt - new data ready)            │
│  └─────────┘           └─────────┘                                          │
│                                                                              │
│  DATA OUTPUT (every 100ms sample)                                           │
│  ════════════════════════════════                                           │
│                                                                              │
│  {                                                                           │
│    "type": "biometric.heart",                                               │
│    "timestamp": "2024-01-19T10:30:00.000Z",                                 │
│    "data": {                                                                │
│      "heart_rate": 72,           // BPM                                     │
│      "heart_rate_confidence": 95, // %                                      │
│      "spo2": 98,                 // %                                       │
│      "spo2_confidence": 92,      // %                                       │
│      "ir_value": 50000,          // Raw IR (finger presence)                │
│      "red_value": 48000          // Raw Red                                 │
│    },                                                                       │
│    "device_id": "wrist-watch-001"                                           │
│  }                                                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Motion + Posture (MPU6050/BMA423)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MOTION SENSOR INTEGRATION                                 │
│                                                                              │
│  DATA OUTPUT (every 50ms sample)                                            │
│  ═══════════════════════════════                                            │
│                                                                              │
│  {                                                                           │
│    "type": "biometric.motion",                                              │
│    "timestamp": "2024-01-19T10:30:00.050Z",                                 │
│    "data": {                                                                │
│      "accel_x": 0.02,           // g-force                                  │
│      "accel_y": -0.01,          // g-force                                  │
│      "accel_z": 1.01,           // g-force (gravity)                        │
│      "gyro_x": 0.5,             // degrees/sec                              │
│      "gyro_y": -0.3,            // degrees/sec                              │
│      "gyro_z": 0.1,             // degrees/sec                              │
│      "steps": 4523,             // Daily count                              │
│      "activity": "sedentary"    // sedentary/walking/running/sleeping       │
│    },                                                                       │
│    "device_id": "wrist-watch-001"                                           │
│  }                                                                          │
│                                                                              │
│  DERIVED INSIGHTS                                                           │
│  ════════════════                                                           │
│                                                                              │
│  • Fall detection: Sudden acceleration spike + no movement                  │
│  • Activity classification: ML model on ESP32                               │
│  • Posture estimation: Wrist angle relative to gravity                      │
│  • Sleep detection: Low movement + time of day                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Implementation Path

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      BIOMETRIC IMPLEMENTATION PHASES                         │
│                                                                              │
│  PHASE 1: Quick Win (1-2 weeks)                                             │
│  ══════════════════════════════                                             │
│                                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                     │
│  │  Buy        │    │  Connect    │    │  Display    │                     │
│  │  PineTime   │───►│  via BLE    │───►│  in Legacy  │                     │
│  │  ($27)      │    │  Bridge     │    │  Context    │                     │
│  └─────────────┘    └─────────────┘    └─────────────┘                     │
│                                                                              │
│  Components:                                                                 │
│  • PineTime with InfiniTime firmware                                        │
│  • Raspberry Pi / ESP32 as BLE-to-NATS bridge                              │
│  • Simple Python script to forward HR to NATS                               │
│                                                                              │
│  ───────────────────────────────────────────────────────────────────────    │
│                                                                              │
│  PHASE 2: Voice-Enabled Wearable (2-4 weeks)                                │
│  ═══════════════════════════════════════════                                │
│                                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                     │
│  │  Get        │    │  Flash      │    │  Talk to    │                     │
│  │  T-Watch V3 │───►│  LiveKit    │───►│  Legacy     │                     │
│  │  ($45)      │    │  Firmware   │    │  from Wrist │                     │
│  └─────────────┘    └─────────────┘    └─────────────┘                     │
│                                                                              │
│  Components:                                                                 │
│  • LILYGO T-Watch 2020 V3 (has mic)                                        │
│  • Custom firmware with LiveKit ESP32 SDK                                   │
│  • Add MAX30102 for heart rate (I2C breakout)                              │
│                                                                              │
│  ───────────────────────────────────────────────────────────────────────    │
│                                                                              │
│  PHASE 3: Custom Neural Interface (4-8 weeks)                               │
│  ════════════════════════════════════════════                               │
│                                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                     │
│  │  Design     │    │  Integrate  │    │  Deploy     │                     │
│  │  Custom PCB │───►│  Full       │───►│  Personal   │                     │
│  │             │    │  Sensor     │    │  Neural     │                     │
│  │             │    │  Suite      │    │  Interface  │                     │
│  └─────────────┘    └─────────────┘    └─────────────┘                     │
│                                                                              │
│  Components:                                                                 │
│  • ESP32-S3 or P4 (latest, best LiveKit support)                           │
│  • MAX30102 (HR + SpO2)                                                     │
│  • MPU6050 or BMA423 (Motion)                                               │
│  • MAX30205 (Body temperature)                                              │
│  • INMP441 (Microphone)                                                     │
│  • Custom 3D-printed enclosure                                              │
│  • 500mAh LiPo battery                                                      │
│                                                                              │
│  Total BOM: ~$30-40                                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Event Schema for Biometrics

```python
# In agents/legacy/consciousness/bus.py - add these event types

class EventType(str, Enum):
    # ... existing types ...

    # Biometric Events
    BIOMETRIC_HEART = "biometric.heart"           # HR, HRV, SpO2
    BIOMETRIC_MOTION = "biometric.motion"         # Steps, activity, posture
    BIOMETRIC_SLEEP = "biometric.sleep"           # Sleep stages, quality
    BIOMETRIC_STRESS = "biometric.stress"         # Derived stress score
    BIOMETRIC_TEMPERATURE = "biometric.temperature"  # Body temp

    # Biometric Alerts (Tier 1 - immediate)
    BIOMETRIC_ALERT_HR_HIGH = "biometric.alert.hr_high"
    BIOMETRIC_ALERT_HR_LOW = "biometric.alert.hr_low"
    BIOMETRIC_ALERT_SPO2_LOW = "biometric.alert.spo2_low"
    BIOMETRIC_ALERT_FALL = "biometric.alert.fall"

    # Biometric Patterns (Tier 2 - derived)
    BIOMETRIC_PATTERN_STRESS = "biometric.pattern.stress"
    BIOMETRIC_PATTERN_FATIGUE = "biometric.pattern.fatigue"
    BIOMETRIC_PATTERN_RECOVERY = "biometric.pattern.recovery"
```

---

## Integration with Existing Modules

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    BIOMETRIC + CONSCIOUSNESS INTEGRATION                     │
│                                                                              │
│  patterns.py                                                                 │
│  ══════════                                                                 │
│  Add biometric pattern detection:                                           │
│  • Stress detection (HR + HRV + activity correlation)                       │
│  • Fatigue detection (HR variability + time of day)                         │
│  • Recovery quality (sleep + morning HRV)                                   │
│                                                                              │
│  temporal.py                                                                 │
│  ═══════════                                                                │
│  Add biometric-aware scheduling:                                            │
│  • suggest_focus_time() considers energy levels                             │
│  • Meeting scheduling avoids low-energy periods                             │
│  • Deadline warnings factor in stress state                                 │
│                                                                              │
│  notifier.py                                                                 │
│  ═══════════                                                                │
│  Add biometric-aware notifications:                                         │
│  • Suppress non-urgent alerts during stress                                 │
│  • Escalate urgent alerts if user unresponsive                              │
│  • Emergency contact integration for health alerts                          │
│                                                                              │
│  presence.py                                                                 │
│  ═══════════                                                                │
│  Add biometric device tracking:                                             │
│  • Wearable as primary presence indicator                                   │
│  • Sleep detection = user unavailable                                       │
│  • Activity detection = context for interactions                            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Hardware Shopping List

### Phase 1: Quick Start ($27)

| Item | Price | Source |
|------|-------|--------|
| PineTime Dev Kit | $27 | pine64.org |
| **Total** | **$27** | |

### Phase 2: Voice Wearable ($60)

| Item | Price | Source |
|------|-------|--------|
| LILYGO T-Watch 2020 V3 | $45 | AliExpress/Amazon |
| MAX30102 Breakout | $8 | Amazon/AliExpress |
| Jumper wires | $5 | Amazon |
| **Total** | **$58** | |

### Phase 3: Custom Build ($40)

| Item | Price | Source |
|------|-------|--------|
| ESP32-S3 DevKit | $8 | AliExpress |
| MAX30102 (HR/SpO2) | $3 | AliExpress |
| MPU6050 (IMU) | $2 | AliExpress |
| MAX30205 (Temp) | $4 | AliExpress |
| INMP441 (Mic) | $2 | AliExpress |
| 500mAh LiPo | $5 | AliExpress |
| TP4056 Charger | $1 | AliExpress |
| PCB (JLCPCB) | $10 | jlcpcb.com |
| 3D Print Case | $5 | Local/Online |
| **Total** | **~$40** | |

---

## Next: EMQX Integration

For IoT sensors that can't run LiveKit (battery sensors, environmental monitors), EMQX provides the bridge:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      EMQX FOR IOT BIOMETRIC SENSORS                          │
│                                                                              │
│   Low-Power Sensors          EMQX                    NATS                   │
│   ┌─────────────────┐       ┌──────┐               ┌──────┐                 │
│   │ Bed sensor      │──────►│      │               │      │                 │
│   │ (sleep tracking)│ MQTT  │      │  Redpanda    │      │                 │
│   └─────────────────┘       │      │  Connect     │      │                 │
│   ┌─────────────────┐       │      │─────────────►│      │──► Legacy       │
│   │ Environment     │──────►│      │               │      │                 │
│   │ (temp/humidity) │ MQTT  │      │               │      │                 │
│   └─────────────────┘       └──────┘               └──────┘                 │
│                                                                              │
│   Use EMQX when: Battery life critical, simple sensors, many devices        │
│   Use LiveKit when: Voice needed, real-time critical, ESP32 available       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Summary

You now have a complete path to giving Legacy **biometric awareness**:

| Capability | Hardware | Transport | Latency |
|------------|----------|-----------|---------|
| Heart Rate | MAX30102 / PineTime | LiveKit / BLE | < 200ms |
| SpO2 | MAX30102 | LiveKit | < 200ms |
| Motion/Steps | BMA423 / MPU6050 | LiveKit | < 100ms |
| Voice Commands | INMP441 / T-Watch mic | LiveKit | < 500ms |
| Body Temp | MAX30205 | LiveKit/MQTT | < 500ms |
| Sleep Tracking | Oura / Bed sensor | API/MQTT | Minutes |

**The result:** Legacy knows your physiological state in real-time, just like Cortana knew Master Chief's.
