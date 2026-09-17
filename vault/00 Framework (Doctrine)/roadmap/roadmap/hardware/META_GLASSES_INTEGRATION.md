# Meta Glasses Integration — Signal Path & Architecture

> **Status**: ROADMAP — Do not build yet. Meta Wearables DAT SDK is developer preview (v0.4.0), GA expected 2026.

## Overview

Meta Ray-Ban smart glasses become a first-person sensor array for OneMind entities. The glasses stream video and audio through the user's phone to a LiveKit room, where an entity processes them using Gemini Live's native video input.

## Signal Path

```
┌─────────────────────┐
│  Meta Ray-Ban       │
│  Smart Glasses      │
│  ┌───────────────┐  │
│  │ Camera (12MP) │  │
│  │ Mic (5-array) │  │
│  │ Speaker       │  │
│  └───────┬───────┘  │
└──────────┼──────────┘
           │ Bluetooth
           ▼
┌─────────────────────┐
│  Phone App          │
│  (Android/iOS)      │
│  ┌───────────────┐  │
│  │ Meta DAT SDK  │──┼── Camera frames (JPEG/YUV)
│  │ v0.4.0+       │──┼── Mic audio (PCM)
│  │               │──┼── Speaker output
│  └───────┬───────┘  │
│  ┌───────┴───────┐  │
│  │ LiveKit SDK   │──┼── WebRTC publish (video + audio tracks)
│  │ (Client)      │  │
│  └───────────────┘  │
└──────────┼──────────┘
           │ WebRTC (UDP/TURN)
           ▼
┌─────────────────────┐
│  LiveKit Server     │
│  (Cloud or Self)    │
│  Room: glasses-zeus │
└──────────┼──────────┘
           │ Server SDK
           ▼
┌──────────────────────────────┐
│  OneMind Entity              │
│  (LiveKit Agent participant) │
│  ┌────────────────────────┐  │
│  │ LiveKit AgentSession   │  │
│  │ video_input=True       │  │
│  │ STT: Deepgram          │  │
│  │ LLM: Gemini Live       │  │──── native video frame processing (1fps while speaking)
│  │ TTS: Cartesia/Gemini   │  │
│  └───────────┬────────────┘  │
│  ┌───────────┴────────────┐  │
│  │ Agno Agent (Brain)     │  │──── memory, tools, personality, knowledge
│  │ Nexus Awareness (NATS) │  │──── event bus integration
│  │ Heartbeat Loop (OM1)   │  │──── proactive monitoring (optional)
│  └────────────────────────┘  │
└──────────────────────────────┘
```

## Two Routes Evaluated

### Route A: RTMP Brute Force (REJECTED)
- Stream glasses video as RTMP via phone → LiveKit Ingress → Room
- **Why rejected**: Extra transcoding hop, higher latency, RTMP is dying, doesn't leverage DAT SDK properly

### Route B: Phone as LiveKit Participant (CHOSEN)
- Phone app captures frames from DAT SDK → publishes as WebRTC tracks via LiveKit SDK
- **Why chosen**: Native WebRTC, lowest latency, SDK-to-SDK clean path, phone handles BT→WebRTC bridge

## Meta Wearables DAT SDK

- **Version**: v0.4.0 (Developer Access Toolkit)
- **Platforms**: Android (primary), iOS (limited)
- **Access**: Developer preview — requires Meta developer account
- **GA Timeline**: Expected 2026
- **Capabilities**:
  - Camera: Live frames, photo capture, video recording
  - Microphone: 5-mic array, beam-forming, noise cancellation
  - Speaker: Direct audio output to glasses
  - Sensors: IMU, head tracking (future)
- **Limitations**:
  - BT bandwidth constrains video resolution (~720p practical max)
  - Battery life: ~4 hours with continuous streaming
  - SDK APIs may change before GA

## LiveKit Agent Vision Configuration

```python
from livekit.agents import AgentSession, Agent, room_io

class GlassesVisionAgent(Agent):
    def __init__(self, agno_agent):
        super().__init__(
            instructions="You can see through the user's smart glasses...",
        )
        self.agno_agent = agno_agent  # Shared Agno brain

session = AgentSession(
    stt="deepgram/nova-3:multi",
    llm="google/gemini-2.0-flash",  # Native video input support
    tts="cartesia/sonic-3:...",
    vad=silero.VAD.load(),
)

await session.start(
    room=ctx.room,
    agent=GlassesVisionAgent(entity.agent),
    room_options=room_io.RoomOptions(
        video_input=True,  # <-- Enable video frame processing
    ),
)
```

## Phone App Architecture (Future Build)

```
Phone App (Flutter or native)
├── Meta DAT SDK integration
│   ├── CameraStreamManager — BT → frame buffer
│   ├── AudioStreamManager — BT → PCM buffer
│   └── SpeakerManager — audio output → BT → glasses
├── LiveKit Client SDK
│   ├── Room connection (token-based auth)
│   ├── Video track publisher (frames → WebRTC)
│   ├── Audio track publisher (mic → WebRTC)
│   └── Audio track subscriber (TTS → speaker)
└── OneMind connection
    ├── Auth token from backend
    └── Room name: glasses-{user_id}
```

## NATS Integration

When glasses are active, events flow through the bus:

```
comms.vision.sensor.glasses-zeus     → Video analysis results
comms.voice.human.zeus-glasses       → Transcribed speech from glasses mic  
comms.audio.agent.entity-response    → TTS response audio
nexus.glasses.connected              → Glasses session started
nexus.glasses.disconnected           → Glasses session ended
physical.glasses.battery             → Battery level updates
physical.glasses.location            → GPS from phone
```

## Prerequisites Before Building

1. **LiveKit integration must be complete** — Entity needs LiveKit AgentSession working
2. **Meta DAT SDK GA** — Currently developer preview, APIs will change
3. **Phone app** — Need a companion app (Flutter recommended since frontend is Flutter)
4. **LiveKit Cloud or self-hosted** — Need a LiveKit server running
5. **Gemini API access** — For native video processing in LiveKit agent

## Build Order (When Ready)

1. Get LiveKit agent working with entity (text + voice first)
2. Test video input with webcam (simulate glasses)  
3. Build phone companion app with LiveKit SDK
4. Add Meta DAT SDK to phone app
5. Test full signal path: glasses → phone → LiveKit → entity
6. Add NATS event publishing for glasses telemetry
7. Add glasses-specific tools to entity (identify objects, read text, navigate)

## Estimated Effort

| Phase | Effort | Dependencies |
|-------|--------|-------------|
| LiveKit voice integration | 2-3 days | LiveKit account, entity system |
| LiveKit video input (webcam) | 1 day | Gemini API key |
| Phone companion app | 1-2 weeks | Flutter, LiveKit SDK |
| Meta DAT SDK integration | 1-2 weeks | Developer access, Android device |
| Full signal path testing | 3-5 days | All above |
| Glasses-specific tools | 1 week | Use cases defined |
