---
id: ui-02-20260111053925
title: "LiveKit Voice Agent - Future Reference"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "LiveKit Voice Agent - Future Reference"
---

# LiveKit Voice Agent - Future Reference

> Saved for future implementation when custom web interface is needed.
> Decision: Using Discord as primary interface first.

---

## Overview

LiveKit is the **official Letta voice integration** that handles:
- Real-time voice streaming
- Voice Activity Detection (VAD)
- Interruption handling
- Multi-platform support (web, mobile, custom apps)

---

## Your Available Services

| Service | Purpose | Status |
|---------|---------|--------|
| Deepgram | Speech-to-Text (STT) | Have account |
| Cartesia | Text-to-Speech (TTS) | Have account |
| Letta | Agent brain | Configured |
| LiveKit | Voice orchestration | Need account |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   LIVEKIT VOICE AGENT                        │
│                                                              │
│   ┌───────────┐    ┌───────────┐    ┌───────────┐          │
│   │ Deepgram  │    │   Letta   │    │ Cartesia  │          │
│   │   (STT)   │◀──▶│  Agent    │◀──▶│   (TTS)   │          │
│   └───────────┘    └───────────┘    └───────────┘          │
│         │                                   │                │
│         └───────────────┬───────────────────┘                │
│                         │                                    │
│              ┌──────────▼──────────┐                        │
│              │   LiveKit Session   │                        │
│              └──────────┬──────────┘                        │
└─────────────────────────┼────────────────────────────────────┘
                          │
           ┌──────────────┼──────────────┐
           │              │              │
    ┌──────▼──────┐ ┌─────▼─────┐ ┌─────▼─────┐
    │   Browser   │ │  Mobile   │ │  Custom   │
    │   Web App   │ │  Browser  │ │   App     │
    └─────────────┘ └───────────┘ └───────────┘
```

---

## Quick Start Code

```python
# voice_agent.py
from livekit.agents import AgentSession, WorkerOptions, cli
from livekit.plugins import openai, deepgram, cartesia

async def entrypoint(ctx):
    session = AgentSession(
        llm=openai.LLM.with_letta(
            agent_id="agent-da585f6f-189f-472a-974c-29b9578f9b55",  # Legacy-AI
            base_url="https://api.letta.com/v1"
        ),
        stt=deepgram.STT(),
        tts=cartesia.TTS()
    )
    await session.start(ctx.room)

if __name__ == "__main__":
    cli.run_app(WorkerOptions(entrypoint_fnc=entrypoint))
```

---

## Setup Steps

### 1. Create LiveKit Account
```
https://cloud.livekit.io
```

### 2. Install Dependencies
```bash
pip install livekit-agents livekit-plugins-openai livekit-plugins-deepgram livekit-plugins-cartesia
```

### 3. Environment Variables
```bash
export LETTA_API_KEY="sk-let-YjhiY2I3ZjItODZmOC00NTQ2LTlhYTYtNjg2NmNhNDA0ODMzOmM5NjlkODhhLTE0MmEtNGUzMy1iNzMwLTRkYjYzODI2MGZmYQ=="
export DEEPGRAM_API_KEY="your-deepgram-key"
export CARTESIA_API_KEY="your-cartesia-key"
export LIVEKIT_URL="wss://your-project.livekit.cloud"
export LIVEKIT_API_KEY="..."
export LIVEKIT_API_SECRET="..."
```

### 4. Run Agent
```bash
python voice_agent.py dev
```

---

## Web Interface Options

### Option A: LiveKit Playground (Fastest)
- Use LiveKit's built-in web interface
- No custom code needed
- Good for testing

### Option B: Custom Web App
```
onemind-voice/
├── backend/
│   └── voice_agent.py      # LiveKit agent
├── frontend/
│   ├── index.html          # Custom UI
│   ├── app.js              # LiveKit client SDK
│   └── styles.css          # Mobile-first design
└── docker-compose.yml
```

### Option C: PWA (Installable)
- Add manifest.json
- Add service worker
- Installable on phone home screen

---

## Why LiveKit (vs Custom Voice)

| Feature | LiveKit | Build Yourself |
|---------|---------|----------------|
| VAD | Built-in, tuned | Implement Silero |
| Interrupts | Automatic | Build stop/resume |
| Streaming | Handled | Tune buffers |
| Latency | Optimized | Trial and error |
| Maintenance | LiveKit team | You |

---

## Resources

- Letta Voice Docs: https://docs.letta.com/guides/voice/overview/
- Letta LiveKit Guide: https://docs.letta.com/guides/voice/livekit/
- LiveKit Letta Plugin: https://docs.livekit.io/agents/models/llm/plugins/letta/
- Example Repo: https://github.com/letta-ai/letta-voice

---

## When to Use This

Use LiveKit when you want:
- Custom branded web interface
- Mobile app without Discord
- Full control over UI/UX
- No dependency on Discord/Telegram

---

## Related

- [[Discord-Integration]] - Current primary interface
- [[Telegram-Integration]] - Text-only backup
- [[00-09 Command (Meta)/03 Legacy-AI (Co-Pilot)/01.02 Letta Agents/]] - Agent configurations
