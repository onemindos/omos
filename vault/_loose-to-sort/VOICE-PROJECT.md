---
name: One Mind Unified Voice Infrastructure
description: Multi-surface voice integration (Discord, Telegram, native nodes) powered by Azure OpenAI
type: project
status: active
created: 2026-05-11
updated: 2026-05-11
owner: Zeus
domain: SO 02 Infrastructure
---

# One Mind Unified Voice Infrastructure

> Real-time bidirectional voice conversations across all One Mind surfaces using Azure OpenAI as the sole voice provider.

---

## Executive Summary

**Goal:** Enable seamless voice communication between Zeus and Legacy across all interaction surfaces — Discord voice channels, Telegram voice messages, and native OpenClaw nodes — through a unified, cost-efficient Azure-based voice infrastructure.

**Current Status:** ✅ Configuration unified and deployed to omos-core. Gateway running successfully with all three voice surfaces connected.

**Voice Identity:** Coral (female voice) across all surfaces for consistency.

---

## What We're Building

### Three Surfaces, One Voice

```
┌─────────────────────────────────────────────────────────────────┐
│                    ZEUS (User Interface)                         │
├─────────────────┬──────────────────┬───────────────────────────┤
│  Discord Voice  │  Telegram Voice  │  Native OpenClaw Nodes    │
│   (Channels)    │  (Messages)      │  (Talk Mode + CLI)        │
└────────┬────────┴────────┬─────────┴───────────────┬───────────┘
         │                 │                        │
         └─────────────────┼────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│              AZURE OPENAI VOICE RESOURCE                         │
│                (onemind-voice)                                   │
│                                                                  │
│  ├── gpt-realtime-1.5      (Real-time audio)                   │
│  ├── gpt-4o-mini-tts       (Text-to-speech)                    │
│  └── gpt-4o-transcribe     (Speech-to-text)                    │
└────────┬────────────────────────────────────────────────────────┘
         │
         ▼
    LEGACY (Agent)
  Real-time conversation
```

### Three Voice Modalities

| Surface | Protocol | Use Case | Status |
|---------|----------|----------|--------|
| **Discord Voice Channels** | WebRTC + gpt-realtime-1.5 | Live conversation in Discord | ✅ Active |
| **Telegram Voice Messages** | HTTP + TTS/transcription | Asynchronous voice via Telegram | ✅ Active (auto-TTS enabled) |
| **Native Nodes** | WebSocket + gpt-realtime-1.5 | Talk mode on omos-core, omos-sentinel | ✅ Configured |

---

## Architecture

### Voice Pipeline

```
INCOMING AUDIO                          OUTGOING AUDIO
User speaks ─────────────────┐         ┌──────► Agent speaks (TTS)
                              │         │
                    gpt-4o-transcribe   gpt-4o-mini-tts
                              │         │
                    Speech-to-Text  Text-to-Speech
                              │         │
                              └─► LEGACY ◄───┘
                            (Agent reasoning)
```

### Azure Voice Resource Configuration

**Endpoint:** `https://onemind-voice.openai.azure.com`  
**Region:** East US 2 (per existing Azure subscription)  
**Authentication:** Azure OpenAI API key (injected via systemd env)

**Deployments:**
- `gpt-realtime-1.5` — Real-time bidirectional audio (Discord voice, Talk mode)
- `gpt-4o-mini-tts` — Low-cost TTS for all auto-reply scenarios
- `gpt-4o-transcribe` — Speech transcription for message-based voice (Telegram, file uploads)
- `gpt-4o` — Chat completions (fallback, already deployed)

### OpenClaw Gateway Configuration

**File:** `/home/ubuntu/.openclaw/openclaw.json` (omos-core)

```json
{
  "talk": {
    "provider": "openai",
    "providers": {
      "openai": {
        "model": "gpt-realtime-1.5",
        "voice": "coral",
        "azureEndpoint": "https://onemind-voice.openai.azure.com",
        "azureDeployment": "gpt-realtime-1.5",
        "azureApiVersion": "2025-04-01-preview"
      }
    }
  },
  "messages": {
    "tts": {
      "auto": "always",
      "provider": "openai",
      "providers": {
        "openai": {
          "baseUrl": "https://onemind-voice.openai.azure.com/openai",
          "model": "gpt-4o-mini-tts",
          "voice": "coral"
        }
      }
    }
  },
  "channels": {
    "discord": {
      "voice": {
        "enabled": true,
        "mode": "agent-proxy",
        "model": "gpt-realtime-1.5",
        "realtime": {
          "provider": "openai"
        },
        "bargeIn": true,
        "minBargeInAudioEndMs": 250
      }
    }
  },
  "tools": {
    "media": {
      "audio": {
        "enabled": true,
        "maxBytes": 20971520,
        "models": [
          {
            "provider": "openai",
            "model": "gpt-4o-transcribe",
            "baseUrl": "https://onemind-voice.openai.azure.com/openai"
          }
        ]
      }
    }
  }
}
```

---

## Voice Identity & Behavior

### Voice: Coral (Female)

Consistent voice across all surfaces ensures:
- Recognizable agent identity (consistent with Legacy brand)
- Professional, natural-sounding interactions
- Cross-channel continuity

### Key Behaviors

| Behavior | Setting | Rationale |
|----------|---------|-----------|
| **Auto-TTS** | Always | All agent messages are spoken automatically (cost: ~$5-10/month) |
| **Barge-in** | Enabled (250ms threshold) | User can interrupt mid-sentence without waiting for response to finish |
| **Agent-proxy mode** (Discord) | Enabled | Voice routed directly to Legacy agent, not generic OpenAI model |
| **Audio length limit** | 20 MB | Reasonable cap for voice files (prevents abuse) |

---

## Technical Decisions

### Why Azure OpenAI (Not Alternatives)

| Criteria | Azure | Google Gemini | ElevenLabs |
|----------|-------|---------------|-----------+
| Real-time audio API | ✅ gpt-realtime-1.5 | ✅ (Gemini 2.5) | ❌ (latency-heavy) |
| TTS cost | ✅ $5-10/mo (auto) | ✅ Gemini free | ❌ $500+/mo |
| Existing credits | ✅ Available | ❌ Separate account | ❌ Separate account |
| Agent integration | ✅ OpenClaw native | ⚠️ Custom bridge | ⚠️ Custom bridge |
| **Decision** | **CHOSEN** | Fallback | Not viable |

### Why Coral Voice

- Natural, professional tone
- Female voice adds personality/brand consistency
- Performs well in real-time scenarios (low latency)
- Consistent across all models (realtime, TTS, text)

### Why Agent-Proxy Mode (Discord)

**Setting:** Discord voice conversations are routed through `openclaw_agent_consult` to Legacy, not directly to gpt-realtime-1.5.

**Why:**
- Discord voice feels like talking to Legacy (not generic OpenAI)
- Legacy can access full agent context (memory, skills, prior conversations)
- Agent-specific personality and reasoning applied
- Fallback to direct model if agent unavailable (graceful degradation)

### Why Auto-TTS Always

**Cost:** ~$5-10/month for continuous voice responses  
**Benefit:** Immersive experience (voice conversations feel more natural than text)  
**Precedent:** HumanAI systems (ChatGPT voice, Google Assistant) all use continuous TTS

---

## Implementation Status

### Phase 1: Configuration Unification ✅ Complete

- [x] Reviewed OpenClaw realtime audio architecture
- [x] Unified fragmented voice config across 3 separate providers (talk, messages.tts, channels.discord.voice)
- [x] Migrated all voice routing to single Azure resource (onemind-voice)
- [x] Resolved schema validation issues (audio vs TTS model config differences)
- [x] Deployed unified config to omos-core
- [x] Verified gateway restart and status

**Status:** Gateway running successfully as of 2026-05-11 09:27:56 UTC

### Phase 2: Testing & Validation 🔄 In Progress

- [ ] Test Discord voice channel (user speaks → Legacy responds)
- [ ] Test Telegram voice message (record voice → get TTS response)
- [ ] Test Talk mode on omos-core (`legacy talk`)
- [ ] Measure latency (target: < 2 seconds end-to-end)
- [ ] Verify Coral voice consistency across surfaces
- [ ] Test barge-in interruption
- [ ] Monitor Azure usage and cost

### Phase 3: Extended Surface Coverage (Planned)

- [ ] Telegram voice message transcription (auto-transcribe incoming voice)
- [ ] Native nodes (omos-sentinel Talk mode via agent-proxy)
- [ ] Desktop app integration (future: native voice via LibreChat or similar)

### Phase 4: Voice Activity Detection & Smart Reply (Planned)

- [ ] Voice activity detection (VAD) to reduce false TTS triggers
- [ ] Intelligent silence detection (don't speak if silence)
- [ ] Streaming response (speak as Legacy thinks, not wait-then-speak)

---

## Current Status

### ✅ Complete

- Unified configuration across all voice surfaces
- Azure voice resource (onemind-voice) with 4 deployments
- Gateway running on omos-core with config applied
- Coral voice identity established

### 🔄 In Progress

- Testing voice interactions across Discord, Telegram, native nodes
- Validating latency and audio quality
- Monitoring Azure costs and usage

### ⏳ Planned

- Extended surface coverage (omos-sentinel, desktop)
- Advanced features (VAD, streaming, smart replies)
- Voice recording playback improvements

---

## Cost Analysis

### Azure Voice Consumption (Monthly Estimate)

| Service | Usage | Cost |
|---------|-------|------|
| **gpt-realtime-1.5** | 4 hrs/month (Discord/Talk) | ~$1.20 |
| **gpt-4o-mini-tts** | 10,000 chars/month (auto TTS) | ~$3-5 |
| **gpt-4o-transcribe** | 2 hrs/month (Telegram voice) | ~$0.60 |
| **Total Estimated** | | **~$5-10/month** |

**Existing Azure Credits:** Sufficient for 6+ months at current consumption.

---

## Related Documentation

| Document | Purpose |
|----------|---------|
| `02 Infrastructure/voice/OPENCLAW-VOICE-RESEARCH.md` | Deep research on OpenClaw Talk Mode and ecosystem integrations |
| `02 Infrastructure/voice/VOICE_MULTIMODAL_PLAN.md` | Original multimodal voice + video/screen share integration plan |
| `00 Framework/AGENTS.md` | Agent architecture and skills system |
| `01 Intelligence/Agents/_COMMAND/{agent}/openclaw/` | Live agent workspace files on VPS |

---

## Deployment Commands

### Check gateway status (omos-core)

```bash
ssh onemindos-core 'systemctl --user status openclaw-gateway'
```

### View current config

```bash
ssh onemindos-core 'openclaw config get talk'
ssh onemindos-core 'openclaw config get messages.tts'
ssh onemindos-core 'openclaw config get channels.discord.voice'
```

### Restart gateway (after config changes)

```bash
ssh onemindos-core 'systemctl --user restart openclaw-gateway'
```

### Live logs

```bash
ssh onemindos-core 'journalctl --user -u openclaw-gateway -f'
```

---

## Success Metrics

| Metric | Target | Status |
|--------|--------|--------|
| End-to-end latency (speech→response) | < 2 seconds | 🔍 Testing |
| Voice recognition accuracy | > 95% | 🔍 Testing |
| Barge-in responsiveness | < 500ms | 🔍 Testing |
| Audio quality (Discord) | CD quality (44.1 kHz) | 🔍 Testing |
| Uptime (gateway) | > 99% | ✅ Running |
| Cost per month | < $15 | ✅ ~$5-10 |
| Voice consistency | Coral on all surfaces | ✅ Configured |

---

## Known Limitations & Future Work

### Current Limitations

1. **Discord voice:** Requires Discord bot to be in voice channel (can't dynamically join yet)
2. **Telegram:** Voice messages still async (working toward real-time integration)
3. **Native nodes:** Sentinel node not yet connected to unified voice (separate Talk mode instance)
4. **VAD:** No voice activity detection yet (continuous TTS on all agent responses)

### Future Enhancements

- **Multi-agent voice:** Different voice for each agent (Legacy = Coral, Zenith = different voice, etc.)
- **Streaming responses:** Speak as agent thinks (lower latency perception)
- **Voice cloning:** Optional custom voice (user's own voice option)
- **Transcription storage:** Archive all voice interactions for future training
- **Voice commands:** "Take a note," "Send a message," etc.

---

## Quick Start

To test the unified voice infrastructure:

### Discord Voice Channel

1. Join voice channel where OpenClaw bot is present
2. Speak: "Hey Legacy, what's 2+2?"
3. Hear response in realtime via Coral voice

### Telegram Voice Message

1. Send voice message in Telegram to @onemindos_bot
2. Bot transcribes and sends to Legacy
3. Legacy responds with auto-TTS (Coral voice)

### Native Talk Mode

```bash
ssh onemindos-core
legacy talk
# Start speaking to agent
```

---

**Last Updated:** 2026-05-11  
**Next Review:** 2026-05-18 (one week testing cycle)  
**Owner:** Zeus / Sovereignty Agent  
**Related Issues:** OMOS-XX (create tracking issues in Linear if needed)
