# Legacy Voice & Multimodal Integration Plan

> Real-time voice, video, and screen sharing with Legacy
> Created: 2026-01-20

---

## Executive Summary

**Goal:** Enable real-time two-way communication with Legacy including voice, video, and screen sharing.

**Discovery:** ~80% of the code already exists! We need to:
1. Configure LiveKit Cloud credentials
2. Deploy the voice agent backend
3. Enhance the UI for video/screen share
4. Wire everything together

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              USER INTERFACE                                  │
│                         (Next.js + LiveKit React)                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │  Microphone  │  │    Camera    │  │ Screen Share │  │   Speaker    │    │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └──────▲───────┘    │
│         │                 │                 │                 │             │
└─────────┼─────────────────┼─────────────────┼─────────────────┼─────────────┘
          │                 │                 │                 │
          ▼                 ▼                 ▼                 │
┌─────────────────────────────────────────────────────────────────────────────┐
│                         LIVEKIT CLOUD                                        │
│                    (WebRTC Media Server)                                     │
│                                                                              │
│   Audio Track ──────────────────────────────────────────────► Audio Track   │
│   Video Track ──────────────────────────────────────────────► (to agent)    │
│   Screen Track ─────────────────────────────────────────────► (to agent)    │
│                                                                              │
└─────────────────────────────────────────────────────────────┬───────────────┘
                                                              │
                                                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      LIVEKIT VOICE AGENT                                     │
│                    (Python - agents/legacy/)                                 │
│                                                                              │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐  │
│   │  Silero VAD │ ──▶│ Deepgram    │ ──▶│   Legacy    │ ──▶│  Cartesia   │  │
│   │  (Voice     │    │ STT         │    │   (Agno)    │    │  TTS        │  │
│   │  Activity)  │    │ (Speech→Text)│   │   Claude    │    │ (Text→Speech)│ │
│   └─────────────┘    └─────────────┘    └──────┬──────┘    └─────────────┘  │
│                                                │                             │
│                                         ┌──────▼──────┐                     │
│                                         │ Vision Agent│                     │
│                                         │ (video/screen│                    │
│                                         │  analysis)   │                    │
│                                         └─────────────┘                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## What Already Exists

### Backend (agents/legacy/)

| File | Status | Description |
|------|--------|-------------|
| `voice.py` | ✅ Complete | LiveKit VoicePipelineAgent + AgnoLLM wrapper |
| `multimodal.py` | ✅ Complete | Vision, video, audio analysis agents |
| `agent.py` | ✅ Complete | Legacy agent definition |
| `main.py` | ✅ Complete | AgentOS entry point |

**voice.py Key Components:**
- `AgnoLLM` class - Wraps Legacy for LiveKit's LLM interface
- `AgnoChatStream` class - Streams responses from Legacy
- `create_voice_pipeline()` - Assembles VAD + STT + LLM + TTS
- Uses: Silero VAD, Deepgram STT (nova-2), Cartesia TTS

**multimodal.py Key Components:**
- `create_vision_agent()` - Claude Sonnet for images
- `create_video_agent()` - Gemini for video (Claude doesn't support video)
- `create_audio_agent()` - Gemini for audio analysis
- `HomeAssistantCameraIntegration` - Pull camera feeds
- `SecurityMonitor` - Continuous monitoring mode

### Frontend (ui/src/components/)

| File | Status | Description |
|------|--------|-------------|
| `voice/VoicePanel.tsx` | ✅ Complete | LiveKit room connection, audio only |
| `onemind/VoiceVisualizer.tsx` | ✅ Complete | Audio visualization (bars/wave/circle) |

**VoicePanel.tsx Key Components:**
- Fetches token from `/api/voice/token`
- LiveKitRoom with audio enabled
- Microphone toggle, connection states
- Voice assistant integration

---

## What Needs to Be Done

### 1. Credentials Setup (P0)

| Service | Required | Status |
|---------|----------|--------|
| LiveKit Cloud | LIVEKIT_URL, LIVEKIT_API_KEY, LIVEKIT_API_SECRET | ⚠️ Need to verify |
| Deepgram | DEEPGRAM_API_KEY | ⚠️ Need to verify |
| Cartesia | CARTESIA_API_KEY | ⚠️ Need to verify |

**1Password Vault:** `00-24 SO (Sovereign Operations)`
- Check for existing `Livekit API - OneMind`

### 2. Backend Deployment (P0)

| Task | Description |
|------|-------------|
| Create voice endpoint | `/api/voice/token` - Generate LiveKit room tokens |
| Deploy voice agent | Run `voice.py` as LiveKit agent worker |
| Add to docker-compose | New service for voice agent |

### 3. Frontend Enhancements (P1)

| Feature | Current | Needed |
|---------|---------|--------|
| Audio | ✅ Working | - |
| Video | ❌ Not enabled | Add camera track + display |
| Screen Share | ❌ Not enabled | Add screen share button + logic |
| Video Display | ❌ None | Grid/focus view for video tracks |

### 4. Video Processing Pipeline (P1)

| Task | Description |
|------|-------------|
| Receive video track | In voice agent, subscribe to video tracks |
| Frame extraction | Capture frames for vision agent |
| Vision analysis | Send frames to `create_vision_agent()` |
| Context injection | Add vision context to Legacy's conversation |

---

## Implementation Plan

### Phase 1: Voice Only (Today)

**Goal:** Get voice working end-to-end

```
Step 1: Verify Credentials
├── Check 1Password for LiveKit credentials
├── Check for Deepgram API key
├── Check for Cartesia API key
└── Add missing credentials to server .env

Step 2: Create Token Endpoint
├── Add /api/voice/token to Agno service
├── Generate LiveKit access tokens
└── Return token + server URL

Step 3: Deploy Voice Agent
├── Add voice agent service to docker-compose
├── Configure worker to connect to LiveKit Cloud
└── Test voice round-trip

Step 4: Test Frontend
├── Navigate to existing VoicePanel
├── Click connect
└── Speak and hear Legacy respond
```

### Phase 2: Video (This Week)

**Goal:** Legacy can see you via camera

```
Step 1: Enable Video in Frontend
├── Add video: true to LiveKitRoom props
├── Add camera toggle button
├── Display local video preview
└── Request camera permissions

Step 2: Receive Video in Agent
├── Subscribe to video tracks in voice agent
├── Set up frame capture interval (1fps for analysis)
└── Send frames to vision agent

Step 3: Integrate Vision Context
├── Vision agent analyzes frames
├── Inject visual context into conversation
└── Legacy describes what it sees when asked
```

### Phase 3: Screen Share (This Week)

**Goal:** Share screen with Legacy for help

```
Step 1: Add Screen Share UI
├── Add screen share button to VoicePanel
├── Call room.localParticipant.setScreenShareEnabled(true)
├── Show screen share preview
└── Handle permission prompts

Step 2: Process Screen Share in Agent
├── Detect screen share track (source: TrackSource.SCREEN_SHARE)
├── Capture frames at higher quality
├── Send to vision agent for analysis
└── Legacy can see and comment on screen content
```

---

## File Changes Required

### New Files

```
agents/legacy/
├── voice_worker.py          # LiveKit agent worker entry point
└── api/
    └── voice_routes.py      # Token generation endpoint

ui/src/components/voice/
├── VideoPanel.tsx           # Video + screen share panel (enhanced)
└── ScreenShareButton.tsx    # Screen share toggle component

infra/compose/cloud/
└── docker-compose.voice.yml # Voice agent service definition
```

### Modified Files

```
agents/legacy/voice.py
├── Add video track handling
├── Add frame capture for vision
└── Add screen share detection

ui/src/components/voice/VoicePanel.tsx
├── Enable video prop
├── Add camera toggle
├── Add screen share button
├── Add video grid display

agents/legacy/main.py
├── Add voice token endpoint
└── Add voice agent routes

infra/compose/cloud/docker-compose.yml
├── Add voice-agent service
└── Add required environment variables
```

---

## Required Credentials

### LiveKit Cloud

```bash
# Get from https://cloud.livekit.io
LIVEKIT_URL=wss://your-project.livekit.cloud
LIVEKIT_API_KEY=APIxxxxxxxx
LIVEKIT_API_SECRET=xxxxxxxxxxxxxxxx
```

### Deepgram (STT)

```bash
# Get from https://console.deepgram.com
DEEPGRAM_API_KEY=xxxxxxxxxxxxxxxx
```

### Cartesia (TTS)

```bash
# Get from https://cartesia.ai
CARTESIA_API_KEY=xxxxxxxxxxxxxxxx
```

---

## Testing Checklist

### Phase 1: Voice
- [ ] Token endpoint returns valid JWT
- [ ] Voice agent connects to LiveKit Cloud
- [ ] Frontend connects to room
- [ ] Speech is transcribed (Deepgram)
- [ ] Legacy responds with text
- [ ] Text is spoken back (Cartesia)
- [ ] Full round-trip < 3 seconds

### Phase 2: Video
- [ ] Camera enabled in frontend
- [ ] Local video preview shows
- [ ] Agent receives video track
- [ ] Frames captured successfully
- [ ] Vision agent analyzes frames
- [ ] Legacy describes what it sees

### Phase 3: Screen Share
- [ ] Screen share button works
- [ ] Screen preview shows
- [ ] Agent receives screen track
- [ ] Screen content analyzed
- [ ] Legacy can help with screen content

---

## Risk Assessment

| Risk | Mitigation |
|------|------------|
| Missing credentials | Check 1Password, create accounts if needed |
| High latency | Use Deepgram nova-2 (fastest), Cartesia (low latency) |
| Cost overrun | Monitor usage, set rate limits |
| Video quality | Use 720p max, reduce frame rate |
| Browser permissions | Handle gracefully, show instructions |

---

## Success Metrics

| Metric | Target |
|--------|--------|
| Voice latency (end-to-end) | < 3 seconds |
| Voice recognition accuracy | > 95% |
| Video analysis frequency | 1 frame/second |
| Connection success rate | > 99% |
| User satisfaction | "It feels like talking to a real assistant" |

---

## Next Actions

1. **Immediate:** Check 1Password for existing credentials
2. **Today:** Implement Phase 1 (voice only)
3. **Tomorrow:** Test and stabilize voice
4. **This week:** Add video and screen share

---

## Related Issues (Linear)

Create these issues to track progress:

- [ ] OMOS-XX: Set up LiveKit Cloud credentials
- [ ] OMOS-XX: Create voice token endpoint
- [ ] OMOS-XX: Deploy voice agent worker
- [ ] OMOS-XX: Test voice round-trip
- [ ] OMOS-XX: Add video support to frontend
- [ ] OMOS-XX: Add screen share support
- [ ] OMOS-XX: Integrate vision with voice agent
