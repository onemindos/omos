# LiveKit

**Type:** Real-time audio/video infrastructure + AI Agents framework
**SDK:** `livekit-agents` (Python)
**Key var:** `LIVEKIT_URL`, `LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET`
**Website:** livekit.io
**Hosting:** LiveKit Cloud (managed) or self-hosted

## What LiveKit Is

Open-source WebRTC infrastructure for real-time audio/video — but for OneMind OS, the key piece is **LiveKit Agents**: a Python framework for building voice AI agents that combines STT + LLM + TTS into a single real-time pipeline.

## Models

See model files in this folder.

## Architecture

```
User voice → LiveKit room → Agent pipeline:
  STT (Deepgram Nova-3)
  → LLM (GPT-4o / Legacy)
  → TTS (Cartesia Sonic-2)
  → Voice back to user
```

Latency target: ~800ms end-to-end (voice in → voice out).

## Stack Use Cases

- Legacy agent real-time voice mode
- Voice-enabled Oracle sessions (research Q&A by voice)
- Phone call handling via LiveKit SIP integration
- Multi-user voice rooms (future: Zeus + agents in the same room)


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
