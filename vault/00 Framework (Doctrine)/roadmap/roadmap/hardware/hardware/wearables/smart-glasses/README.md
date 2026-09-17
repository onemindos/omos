# Smart Glasses Integration - OneMind OS

> Status: PLANNING
> Last Updated: 2026-01-19

---

## Overview

Smart glasses provide Legacy with first-person vision and always-available voice interface. This enables true omnipresent mode where Legacy sees what you see.

### Integration Points

| Component | OneMind Integration |
|-----------|---------------------|
| Camera | `multimodal.py` → Vision AI (Claude) |
| Microphone | `voice.py` → LiveKit STT (Deepgram) |
| Display | Legacy responses rendered on lens |
| Sensors | Feed registry as `FeedCategory.WEARABLE` |

---

## Option 1: Brilliant Labs Frame

### Specifications

| Spec | Value |
|------|-------|
| **Price** | ~$350 |
| **Camera** | 720p wide-angle |
| **Display** | 640x400 MicroOLED (20° FOV) |
| **Audio** | Microphone + Speaker |
| **Connectivity** | Bluetooth 5.0 (via phone) |
| **Battery** | ~30 min video, 10hr standby |
| **Weight** | 40g |
| **Open Source** | Full (hardware + firmware + SDK) |

### Why Frame

- **Open Source** - Full customization possible
- **Python SDK** - Direct Agno integration
- **Noa AI** - Can be replaced with Legacy
- **Active Development** - Regular firmware updates
- **Community** - Growing developer ecosystem

### Resources

- Docs: https://docs.brilliant.xyz
- GitHub: https://github.com/brilliantlabsAR
- SDK: `pip install frame-sdk`

### Integration Architecture

```
Frame Glasses
     │
     │ Bluetooth
     ▼
Phone (Bridge App)
     │
     │ WebSocket / WebRTC
     ▼
OneMind OS (Legacy)
     │
     ├── Camera → multimodal.py → Claude Vision
     ├── Audio → voice.py → LiveKit Pipeline
     └── Display ← Legacy Response
```

### TODO

- [ ] Order Brilliant Labs Frame
- [ ] Test Frame SDK with Python
- [ ] Build bridge app (or use official app)
- [ ] Create `frame_bridge.py` for Legacy integration
- [ ] Register as wearable feed in `feeds.py`
- [ ] Test voice commands end-to-end
- [ ] Test vision analysis latency

---

## Option 2: ESP32 DIY Camera Glasses

### Why DIY

- **Direct WiFi** - No phone bridge needed
- **Full Control** - Custom firmware
- **Lower Cost** - ~$50-100 in parts
- **LiveKit Native** - ESP32 SDK available
- **Matches Watch** - Same platform as ESP32 watch

### Hardware Options

#### Option A: XIAO ESP32S3 Sense (Recommended)

| Spec | Value |
|------|-------|
| **MCU** | ESP32-S3 (dual-core 240MHz) |
| **Camera** | OV2640 (2MP) or OV5640 (5MP) |
| **Microphone** | PDM digital mic |
| **Size** | 21x17.8mm (tiny!) |
| **WiFi** | 2.4GHz 802.11 b/g/n |
| **Bluetooth** | BLE 5.0 |
| **Price** | ~$15 |

**Purchase:** https://www.seeedstudio.com/XIAO-ESP32S3-Sense-p-5639.html

#### Option B: ESP32-CAM (Budget)

| Spec | Value |
|------|-------|
| **MCU** | ESP32 (dual-core 240MHz) |
| **Camera** | OV2640 (2MP) |
| **Microphone** | None (add I2S mic) |
| **Size** | 40x27mm |
| **WiFi** | 2.4GHz |
| **Price** | ~$8 |

#### Option C: Freenove ESP32-S3-WROOM CAM

| Spec | Value |
|------|-------|
| **MCU** | ESP32-S3 |
| **Camera** | OV2640 or OV5640 |
| **Display** | Optional LCD support |
| **Price** | ~$20 |

### Additional Components Needed

| Component | Purpose | Price |
|-----------|---------|-------|
| I2S Microphone (INMP441) | Audio input | ~$5 |
| I2S Amplifier (MAX98357) | Audio output | ~$5 |
| Small Speaker | Audio feedback | ~$3 |
| LiPo Battery (500mAh) | Power | ~$8 |
| TP4056 Charger | Battery charging | ~$2 |
| 3D Printed Frame | Glasses mount | DIY |
| **Total** | | **~$50-80** |

### Software Stack

```
┌─────────────────────────────────────────┐
│           ESP32-S3 Firmware             │
├─────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────────┐  │
│  │ ESP-IDF     │  │ Arduino Core    │  │
│  │ Framework   │  │ (alternative)   │  │
│  └─────────────┘  └─────────────────┘  │
├─────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────────┐  │
│  │ LiveKit     │  │ Camera Driver   │  │
│  │ ESP32 SDK   │  │ (esp32-camera)  │  │
│  └─────────────┘  └─────────────────┘  │
├─────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────────┐  │
│  │ WiFi/MQTT   │  │ I2S Audio       │  │
│  │ Connection  │  │ (mic + speaker) │  │
│  └─────────────┘  └─────────────────┘  │
└─────────────────────────────────────────┘
```

### Integration Architecture (Direct WiFi)

```
ESP32 DIY Glasses
     │
     │ WiFi (WebRTC/MQTT)
     ▼
OneMind OS (Legacy)
     │
     ├── Camera Stream → LiveKit Room → multimodal.py
     ├── Audio Stream → LiveKit Room → voice.py
     ├── Sensor Data → MQTT → consciousness/bus.py
     └── Commands ← NATS → ESP32
```

### Firmware Features (Planned)

```cpp
// Core features
- [ ] Camera streaming via LiveKit WebRTC
- [ ] Audio input via I2S microphone
- [ ] Audio output via I2S speaker
- [ ] WiFi connection management
- [ ] MQTT for commands/status
- [ ] OTA firmware updates
- [ ] Power management (deep sleep)
- [ ] Button/touch for wake word

// Advanced features
- [ ] On-device wake word detection
- [ ] Local face detection (ESP-WHO)
- [ ] Gesture detection (if IMU added)
- [ ] Location via WiFi triangulation
```

### Project Structure (Planned)

```
hardware/smart-glasses/
├── README.md                 # This file
├── brilliant-frame/          # Frame integration
│   ├── frame_bridge.py       # Legacy bridge
│   └── requirements.txt
└── esp32-cam-glasses/        # DIY glasses
    ├── platformio.ini
    ├── src/
    │   ├── main.cpp
    │   ├── config.h
    │   ├── camera/
    │   │   └── camera_manager.cpp
    │   ├── audio/
    │   │   ├── i2s_mic.cpp
    │   │   └── i2s_speaker.cpp
    │   ├── communication/
    │   │   ├── wifi_manager.cpp
    │   │   ├── mqtt_client.cpp
    │   │   └── livekit_client.cpp
    │   └── power/
    │       └── power_manager.cpp
    └── 3d-models/
        └── glasses_frame.stl
```

### TODO

- [ ] Order XIAO ESP32S3 Sense
- [ ] Order INMP441 microphone
- [ ] Order MAX98357 amplifier + speaker
- [ ] Design 3D printed glasses frame
- [ ] Set up PlatformIO project
- [ ] Test camera streaming to LiveKit
- [ ] Test audio I/O
- [ ] Integrate with Legacy feed system
- [ ] Power optimization

---

## Comparison Summary

| Feature | Brilliant Frame | ESP32 DIY |
|---------|-----------------|-----------|
| **Cost** | $350 | $50-80 |
| **Setup Time** | Hours | Days/Weeks |
| **Phone Required** | Yes | No |
| **Display** | Yes (MicroOLED) | No (unless added) |
| **Camera Quality** | 720p | 2-5MP |
| **Audio Quality** | Good | DIY dependent |
| **Form Factor** | Polished | DIY dependent |
| **Customization** | Firmware only | Full hardware |
| **WiFi Direct** | No | Yes |
| **LiveKit Native** | Via phone | Direct |

---

## Recommendation

### Short Term: Brilliant Labs Frame
- Quick to deploy
- Polished experience
- Good for testing Legacy integration

### Long Term: ESP32 DIY
- Full control
- No phone dependency
- Matches existing ESP32 watch platform
- Can iterate on hardware

### Hybrid Approach
1. Start with Frame to validate UX and integration patterns
2. Build ESP32 version with lessons learned
3. Use Frame for indoor/polished situations
4. Use ESP32 for outdoor/rugged/experimental

---

## Feed Registration

Both options will register as wearable feeds:

```python
# In feeds.py or during init
from agents.legacy.feeds import register_feed, FeedType, FeedCategory, FeedPriority

# Brilliant Labs Frame
register_feed(
    id="zeus-glasses-frame",
    name="Zeus Glasses (Frame)",
    type=FeedType.WEBRTC,
    url="lk://wearables/glasses-frame",
    category=FeedCategory.WEARABLE,
    livekit_room="lk-wearables",
    livekit_track="glasses-frame",
    is_mobile=True,
    has_audio=True,
    priority=FeedPriority.HIGH,
    tags=["glasses", "wearable", "first-person"],
    capabilities=["720p", "microphone", "display"],
)

# ESP32 DIY Glasses
register_feed(
    id="zeus-glasses-diy",
    name="Zeus Glasses (DIY)",
    type=FeedType.WEBRTC,
    url="lk://wearables/glasses-diy",
    category=FeedCategory.WEARABLE,
    livekit_room="lk-wearables",
    livekit_track="glasses-diy",
    is_mobile=True,
    has_audio=True,
    priority=FeedPriority.HIGH,
    tags=["glasses", "wearable", "first-person", "esp32"],
    capabilities=["1080p", "microphone", "wifi-direct"],
)
```

---

## References

### Brilliant Labs
- Documentation: https://docs.brilliant.xyz
- GitHub: https://github.com/brilliantlabsAR
- Frame SDK: https://github.com/brilliantlabsAR/frame-sdk-python
- Community: https://discord.gg/brilliant

### ESP32 Camera
- ESP32-Camera Driver: https://github.com/espressif/esp32-camera
- LiveKit ESP32 SDK: https://github.com/livekit/client-sdk-esp32
- ESP-WHO (Face Detection): https://github.com/espressif/esp-who
- XIAO ESP32S3: https://wiki.seeedstudio.com/xiao_esp32s3_getting_started/

### 3D Printing
- Printables: Search "ESP32 glasses"
- Thingiverse: Search "smart glasses frame"

---

## Notes

- Context7 library for LiveKit ESP32: `/livekit/client-sdk-esp32`
- Context7 library for ESP32: `/espressif/arduino-esp32`
- Existing ESP32 watch code in `hardware/esp32-watch/` can be referenced
