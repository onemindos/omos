# OneMind OS Senses - Multimodal Input/Output Architecture

> **Created:** 2026-01-25
> **Status:** PLANNING
> **Linear Issue:** OMOS-XXX (to be created)

---

## Overview

**OneMind Senses** defines how Legacy perceives the world (inputs), understands context (awareness), and takes action (outputs). This is the sensory system that makes Legacy truly present and aware.

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              ONEMIND OS SENSES                                       │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                      │
│                              ┌─────────────────────┐                                │
│                              │       LEGACY        │                                │
│                              │     (The Mind)      │                                │
│                              │                     │                                │
│                              │  Claude Opus 4.5    │                                │
│                              │  + Agno Framework   │                                │
│                              └──────────┬──────────┘                                │
│                                         │                                           │
│            ┌────────────────────────────┼────────────────────────────┐              │
│            │                            │                            │              │
│            ▼                            ▼                            ▼              │
│  ┌─────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐     │
│  │    INPUT SENSES     │    │   AWARENESS SENSES  │    │   OUTPUT SENSES     │     │
│  │  (Receive info)     │    │   (Perceive context)│    │   (Express/Act)     │     │
│  └─────────────────────┘    └─────────────────────┘    └─────────────────────┘     │
│                                                                                      │
│  ════════════════════════════════════════════════════════════════════════════════   │
│                                                                                      │
│                            INPUT SENSES                                             │
│                            ────────────                                             │
│                                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│  │   👂 HEAR    │  │   👁️ VISION  │  │   📝 TEXT    │  │   🎥 VIDEO   │            │
│  │              │  │              │  │              │  │              │            │
│  │ Voice input  │  │ Image input  │  │ Text input   │  │ Video input  │            │
│  │ from user    │  │ from cameras │  │ from chat/   │  │ from cameras │            │
│  │ or environ-  │  │ or uploads   │  │ extension/   │  │ or streams   │            │
│  │ ment         │  │              │  │ API          │  │              │            │
│  │              │  │              │  │              │  │              │            │
│  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │            │
│  │ │ Deepgram │ │  │ │ Claude   │ │  │ │ Chat UI  │ │  │ │ Frigate  │ │            │
│  │ │ STT      │ │  │ │ Vision   │ │  │ │          │ │  │ │ NVR      │ │            │
│  │ ├──────────┤ │  │ ├──────────┤ │  │ ├──────────┤ │  │ ├──────────┤ │            │
│  │ │ LiveKit  │ │  │ │ Frame    │ │  │ │ Browser  │ │  │ │ LiveKit  │ │            │
│  │ │ Audio    │ │  │ │ Camera   │ │  │ │ Extension│ │  │ │ Streams  │ │            │
│  │ ├──────────┤ │  │ ├──────────┤ │  │ ├──────────┤ │  │ ├──────────┤ │            │
│  │ │ Whisper  │ │  │ │ Yi       │ │  │ │ REST     │ │  │ │ Robot    │ │            │
│  │ │ (Offline)│ │  │ │ Cameras  │ │  │ │ API      │ │  │ │ Cameras  │ │            │
│  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │            │
│  │              │  │              │  │              │  │              │            │
│  │ Status: ✅   │  │ Status: ⚠️   │  │ Status: ✅   │  │ Status: ⚠️   │            │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘            │
│                                                                                      │
│  ════════════════════════════════════════════════════════════════════════════════   │
│                                                                                      │
│                          AWARENESS SENSES                                           │
│                          ────────────────                                           │
│                                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│  │  📍 PRESENCE │  │   🕐 TIME    │  │  🌡️ CONTEXT  │  │   💓 HEALTH  │            │
│  │              │  │              │  │              │  │              │            │
│  │ Where you    │  │ When things  │  │ Environment  │  │ Body state   │            │
│  │ are and on   │  │ happen and   │  │ conditions   │  │ and energy   │            │
│  │ what device  │  │ scheduling   │  │ around you   │  │ levels       │            │
│  │              │  │              │  │              │  │              │            │
│  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │            │
│  │ │ Device   │ │  │ │ Calendar │ │  │ │ Home     │ │  │ │ Oura     │ │            │
│  │ │ Tracking │ │  │ │ Events   │ │  │ │ Assistant│ │  │ │ Ring     │ │            │
│  │ ├──────────┤ │  │ ├──────────┤ │  │ ├──────────┤ │  │ ├──────────┤ │            │
│  │ │ Location │ │  │ │ Sun      │ │  │ │ Weather  │ │  │ │ Garmin   │ │            │
│  │ │ Tracking │ │  │ │ State    │ │  │ │ API      │ │  │ │ Watch    │ │            │
│  │ ├──────────┤ │  │ ├──────────┤ │  │ ├──────────┤ │  │ ├──────────┤ │            │
│  │ │ Room     │ │  │ │ Schedule │ │  │ │ Energy   │ │  │ │ Exist.io │ │            │
│  │ │ Presence │ │  │ │ Awareness│ │  │ │ Prices   │ │  │ │          │ │            │
│  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │            │
│  │              │  │              │  │              │  │              │            │
│  │ Status: ✅   │  │ Status: ✅   │  │ Status: ✅   │  │ Status: ⚠️   │            │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘            │
│                                                                                      │
│  ════════════════════════════════════════════════════════════════════════════════   │
│                                                                                      │
│                           OUTPUT SENSES                                             │
│                           ─────────────                                             │
│                                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│  │   🗣️ VOICE   │  │  📺 DISPLAY  │  │   ⚡ ACTION  │  │   🤖 ROBOT   │            │
│  │              │  │              │  │              │  │              │            │
│  │ Speak to     │  │ Show info    │  │ Take actions │  │ Physical     │            │
│  │ user via     │  │ on screens   │  │ in digital   │  │ actions in   │            │
│  │ speakers     │  │ and devices  │  │ world        │  │ real world   │            │
│  │              │  │              │  │              │  │              │            │
│  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │            │
│  │ │ Cartesia │ │  │ │ Flutter  │ │  │ │ Tools    │ │  │ │ Unitree  │ │            │
│  │ │ TTS      │ │  │ │ UI       │ │  │ │ (Todoist,│ │  │ │ Robot    │ │            │
│  │ ├──────────┤ │  │ ├──────────┤ │  │ │ GitHub)  │ │  │ ├──────────┤ │            │
│  │ │ LiveKit  │ │  │ │ Frame    │ │  │ ├──────────┤ │  │ │ OpenMind │ │            │
│  │ │ Audio    │ │  │ │ Glasses  │ │  │ │ HITL     │ │  │ │ OM1      │ │            │
│  │ ├──────────┤ │  │ ├──────────┤ │  │ │ Approvals│ │  │ ├──────────┤ │            │
│  │ │ Speakers │ │  │ │ Push     │ │  │ ├──────────┤ │  │ │ Drones   │ │            │
│  │ │ (Home)   │ │  │ │ Notifs   │ │  │ │ Workflows│ │  │ │          │ │            │
│  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │            │
│  │              │  │              │  │              │  │              │            │
│  │ Status: ✅   │  │ Status: ✅   │  │ Status: ✅   │  │ Status: ⚠️   │            │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘            │
│                                                                                      │
│  ════════════════════════════════════════════════════════════════════════════════   │
│                                                                                      │
│                     CONSCIOUSNESS PROCESSOR                                         │
│                     ───────────────────────                                         │
│                                                                                      │
│  All input senses flow through the Consciousness Processor which:                   │
│  • Filters based on awareness mode (dormant/aware/present/omnipresent)             │
│  • Prioritizes by urgency and relevance                                            │
│  • Routes to appropriate handlers                                                   │
│  • Prevents information overload                                                    │
│                                                                                      │
│  See: docs/legacy/CONSCIOUSNESS-PROCESSOR.md                                        │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Senses Inventory

### Input Senses

| Sense | Purpose | Technology | Status | Missing |
|-------|---------|------------|--------|---------|
| **👂 HEAR** | Voice commands, ambient audio | Deepgram STT, LiveKit | ✅ Working | - |
| **👁️ VISION** | See images, analyze visuals | Claude Vision | ⚠️ Partial | Continuous camera feeds |
| **📝 TEXT** | Read messages, documents | Chat UI, Extension, API | ✅ Working | - |
| **🎥 VIDEO** | Watch video streams | Frigate NVR, LiveKit | ⚠️ Planned | Frigate → NATS integration |

### Awareness Senses

| Sense | Purpose | Technology | Status | Missing |
|-------|---------|------------|--------|---------|
| **📍 PRESENCE** | Know where you are | Device tracking, room presence | ✅ Working | - |
| **🕐 TIME** | Understand scheduling | Calendar, sun state, schedule | ✅ Working | - |
| **🌡️ CONTEXT** | Environment awareness | Home Assistant, weather | ✅ Working | - |
| **💓 HEALTH** | Body state awareness | Oura, Garmin | ⚠️ Planned | Nango OAuth setup |

### Output Senses

| Sense | Purpose | Technology | Status | Missing |
|-------|---------|------------|--------|---------|
| **🗣️ VOICE** | Speak to user | Cartesia TTS, LiveKit | ✅ Working | - |
| **📺 DISPLAY** | Show information | Flutter UI, notifications | ✅ Working | - |
| **⚡ ACTION** | Digital actions | Tools, workflows, HITL | ✅ Working | - |
| **🤖 ROBOT** | Physical actions | Unitree, OM1, drones | ⚠️ Planned | Full robot integration |

---

## Components Needed

| Component | Sense | What It Does | Already Have? | Build Effort |
|-----------|-------|--------------|---------------|--------------|
| **Deepgram** | Hear | Speech-to-text | ✅ YES | - |
| **LiveKit** | Hear/Voice | Audio streaming | ✅ YES | - |
| **Claude Vision** | Vision | Image analysis | ✅ YES | - |
| **Frigate NVR** | Vision/Video | Camera AI | ⚠️ Ready, not deployed | 2-3 days |
| **Frigate→NATS** | Video | Stream events to NATS | ❌ NO | 1-2 days |
| **Frame SDK** | Vision | AR glasses camera | ⚠️ UI done | 2-3 days |
| **Cartesia** | Voice | Text-to-speech | ✅ YES | - |
| **Device Presence** | Presence | Track devices | ✅ YES | - |
| **Room Presence** | Presence | Track room location | ✅ YES | - |
| **Home Assistant** | Context | Smart home data | ✅ YES | - |
| **Oura Integration** | Health | Sleep, HRV data | ⚠️ Nango ready | 1 day |
| **Garmin Integration** | Health | Activity data | ⚠️ Nango ready | 1 day |
| **Unitree/OM1** | Robot | Physical robot control | ⚠️ Planned | 4-6 weeks |

**Total Build Effort for Full Senses:** ~6-8 weeks

---

## Data Flow

### Input Flow

```
┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐
│   HEAR     │    │   VISION   │    │    TEXT    │    │   VIDEO    │
│            │    │            │    │            │    │            │
│ Microphone │    │  Camera    │    │  Keyboard  │    │  Camera    │
│     ↓      │    │     ↓      │    │     ↓      │    │  Stream    │
│ Deepgram   │    │  Claude    │    │  API       │    │     ↓      │
│  (STT)     │    │  Vision    │    │            │    │  Frigate   │
└─────┬──────┘    └─────┬──────┘    └─────┬──────┘    └─────┬──────┘
      │                 │                 │                 │
      └─────────────────┴─────────────────┴─────────────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │   Consciousness       │
                    │   Processor           │
                    │                       │
                    │   • Filter by mode    │
                    │   • Prioritize        │
                    │   • Route             │
                    └───────────┬───────────┘
                                │
                                ▼
                    ┌───────────────────────┐
                    │       LEGACY          │
                    │     (Processing)      │
                    └───────────────────────┘
```

### Output Flow

```
                    ┌───────────────────────┐
                    │       LEGACY          │
                    │     (Decision)        │
                    └───────────┬───────────┘
                                │
      ┌─────────────────────────┼─────────────────────────┐
      │                         │                         │
      ▼                         ▼                         ▼
┌───────────┐           ┌───────────┐           ┌───────────┐
│   VOICE   │           │  DISPLAY  │           │  ACTION   │
│           │           │           │           │           │
│ Cartesia  │           │ Flutter   │           │ Tools     │
│ (TTS)     │           │ UI        │           │ Workflows │
│    ↓      │           │    ↓      │           │    ↓      │
│ LiveKit   │           │ WebSocket │           │ NATS      │
│    ↓      │           │    ↓      │           │    ↓      │
│ Speakers  │           │ Screen    │           │ Services  │
└───────────┘           └───────────┘           └───────────┘
```

---

## Sense Configurations

### Hear Configuration

```yaml
# core/senses/hear.yaml
hear:
  primary:
    provider: deepgram
    model: nova-2
    language: en-US
    punctuate: true
    diarize: false  # Single user

  streaming:
    provider: livekit
    room_prefix: onemind-voice
    codec: opus

  offline:
    provider: whisper
    model: base
    device: edge

  triggers:
    wake_word: "Legacy"
    always_listen: false  # Only when in voice mode
```

### Vision Configuration

```yaml
# core/senses/vision.yaml
vision:
  image_analysis:
    provider: claude
    model: claude-3-5-sonnet-20241022
    max_tokens: 1000

  cameras:
    provider: frigate
    detection:
      - person
      - car
      - dog
      - cat
    zones:
      - front_door
      - backyard
      - barn

  wearable:
    provider: brilliant_frame
    features:
      - capture
      - display
      - translate
```

### Presence Configuration

```yaml
# core/senses/presence.yaml
presence:
  device:
    sources:
      - tailscale  # Which device connected
      - homeassistant  # WiFi presence
    devices:
      - zeus-macbook
      - zeus-iphone
      - zeus-ipad

  location:
    sources:
      - homeassistant
      - gps  # When outdoor
    locations:
      - home
      - office
      - farm
      - car

  room:
    sources:
      - motion_sensors
      - bluetooth_beacons
    rooms:
      - office
      - bedroom
      - kitchen
      - living_room
```

---

## Implementation Phases

### Phase 1: Core Senses (Current)
- [x] Hear (Deepgram + LiveKit)
- [x] Text (Chat UI + API)
- [x] Voice (Cartesia + LiveKit)
- [x] Display (Flutter UI)
- [x] Action (Tools + Workflows)
- [x] Presence (Device + Location)
- [x] Time (Schedule + Sun state)
- [x] Context (Home Assistant)

### Phase 2: Enhanced Vision (2-3 weeks)
- [ ] Frigate NVR deployment
- [ ] Frigate → NATS event bridge
- [ ] Continuous camera monitoring
- [ ] Motion/person detection alerts
- [ ] Video clip analysis

### Phase 3: Health Awareness (1 week)
- [ ] Oura Ring OAuth via Nango
- [ ] Garmin OAuth via Nango
- [ ] Daily health metrics sync
- [ ] Sleep/HRV/activity tracking
- [ ] Health-aware recommendations

### Phase 4: Wearable Vision (2-3 weeks)
- [ ] Frame glasses BLE service
- [ ] Camera capture from glasses
- [ ] Display text on glasses
- [ ] Voice through glasses

### Phase 5: Robot Output (4-6 weeks)
- [ ] Unitree SDK integration
- [ ] OpenMind OM1 bridge
- [ ] Physical task execution
- [ ] Robot video streaming

---

## API Endpoints

```python
# core/senses/api.py

# Input Senses
@router.post("/senses/hear")
async def process_audio(audio: UploadFile):
    """Process audio input (STT)."""

@router.post("/senses/vision")
async def process_image(image: UploadFile):
    """Analyze an image."""

@router.post("/senses/video")
async def process_video(video_url: str):
    """Analyze video stream."""

# Awareness Senses
@router.get("/senses/presence")
async def get_presence():
    """Get current presence state."""

@router.get("/senses/context")
async def get_context():
    """Get current context (time, location, environment)."""

@router.get("/senses/health")
async def get_health():
    """Get health metrics (Oura, Garmin)."""

# Output Senses
@router.post("/senses/speak")
async def speak(text: str, voice: str = "default"):
    """Convert text to speech."""

@router.post("/senses/display")
async def display(content: dict, target: str = "ui"):
    """Display content on target device."""

# Status
@router.get("/senses/status")
async def get_senses_status():
    """Get status of all senses."""
```

---

## Related Documentation

- `docs/legacy/CONSCIOUSNESS-PROCESSOR.md` - Event filtering
- `docs/legacy/VOICE_INTEGRATION.md` - Voice system
- `docs/architecture/PRESENCE-SYSTEM.md` - Presence tracking
- `docs/planning/UNITREE-OPENMIND-INTEGRATION.md` - Robot integration
- `ui/lib/features/wearables/screens/frame_screen.dart` - Frame glasses UI

---

*OneMind Senses - How Legacy perceives and interacts with the world*
