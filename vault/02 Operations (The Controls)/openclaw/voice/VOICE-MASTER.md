---
id: so-05-voice-master
title: "OneMind OS — Unified Voice Architecture"
type: doctrine
status: active
created: 2026-05-27
updated: 2026-06-05
owner: zeus
pillar: SO
domain: "05 Operations (The Controls)"
tags: [voice, livekit, realtime, local-ai, tool-calling, legacy-ai, stt, tts, llm, azure, nats, robotics]
sources:
  - local-realtime-voice.md
  - onemind-os-vision-realtime.md
  - research-local-realtime-voice-ai-2026.md
  - research-local-realtime-voice-ai-tool-calling.md
  - qwen2.5-omni-local-setup-guide.md
  - VOICE-PROJECT.md (SO/03 Infrastructure)
consolidated: 2026-06-05
---

# OneMind OS — Unified Voice Architecture

> All six source voice docs merged. Full content preserved — nothing summarized.

---



---

## PART 1: THE FULL VISION — Legacy as Nervous System + Realtime Voice

*Source: onemind-os-vision-realtime.md*

---

# OneMind OS — The Full Vision

## Legacy as the Nervous System

```
══════════════════════════════════════════════════════════════════════════
                    ZEUS (Human Operator)
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    INTERFACES (Flutter + Voice)                       │
│                                                                      │
│  onemind-app (Flutter)          LiveKit Voice Agent                  │
│  ├── Chat → Hermes WebSocket    ├── Wake word → STT → Hermes        │
│  ├── Fleet → NATS $SRV direct   ├── Hermes response → TTS → Speaker │
│  ├── Tasks → Hermes API         └── Works on phone, watch, home     │
│  ├── Infra → Hermes tools                                           │
│  └── Voice → LiveKit WebRTC     AirPods / HomePod / Barn speaker    │
│                                                                      │
│  ALL interfaces are just INPUT → LEGACY → OUTPUT                     │
└────────────────────────────────┬────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│                   LEGACY (Hermes — The Brain)                         │
│                                                                      │
│  "Zeus said 'send the drone to check the barn roof'"                 │
│                                                                      │
│  1. NLU: Understands intent (recon mission, drone, barn)            │
│  2. Plans: Creates mission (waypoints, camera angle, RTH)            │
│  3. Dispatches: nats_publish("cmd.drone.recon.mission", payload)     │
│  4. Monitors: Subscribes to "telemetry.drone.recon.>" for status    │
│  5. Responds: "Drone launching. ETA 2 min. I'll send you the feed." │
│  6. Follows up: When done, pushes video to Flutter + summarizes     │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │  PROFILES (same brain, different modes)                      │    │
│  │                                                              │    │
│  │  Legacy (orchestrator) — routes, decides, monitors           │    │
│  │  Zenith (builder) — writes code, deploys, fixes              │    │
│  │  Sentinel (watcher) — security, health, anomaly detection    │    │
│  │  Sovereignty (strategist) — long-term planning, doctrine     │    │
│  │  Unity (communicator) — external comms, social, outreach     │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │  PLUGINS (hot-loadable capabilities)                         │    │
│  │                                                              │    │
│  │  nats-platform     — ACP gateway (Synadia Agent SDK)         │    │
│  │  nats-watchers     — proactive event-driven sessions         │    │
│  │  kanban            — task board + worker dispatch             │    │
│  │  observability     — Langfuse traces                         │    │
│  │  fleet-manager     — $SRV discovery + health aggregation     │    │
│  │  robotics-bridge   — ROS2/MAVLink/Unitree SDK bindings       │    │
│  │  voice-dispatch    — LiveKit → intent → immediate execution  │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │  TOOLS (callable from any session)                           │    │
│  │                                                              │    │
│  │  nats_request / nats_publish / nats_kv_* / nats_stream_*    │    │
│  │  docker_status / service_health / system_metrics             │    │
│  │  drone_command / robot_command / device_control              │    │
│  │  langfuse_query / session_search / memory                    │    │
│  │  home_assistant (via cmd.ha.*)                               │    │
│  │  saas_execute (via saas.*)                                   │    │
│  │  terminal / file / web / browser / vision                    │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │  SKILLS (self-improving procedural memory)                   │    │
│  │                                                              │    │
│  │  "Last time I sent the drone at 4pm, sun glare ruined the   │    │
│  │   footage. Skill updated: prefer morning flights for roof    │    │
│  │   inspections. Auto-scheduling adjusted."                    │    │
│  │                                                              │    │
│  │  Skills accumulate. Legacy gets BETTER at operating the      │    │
│  │  fleet every time.                                           │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
└────────────────────────────────┬────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    NATS (The Nervous System)                          │
│                                                                      │
│  Subject hierarchy:                                                  │
│                                                                      │
│  cmd.drone.recon.>        — commands TO drone RECON                  │
│  cmd.drone.recon.mission  — launch mission                           │
│  cmd.drone.recon.rtl      — return to launch                         │
│  cmd.drone.recon.abort    — emergency stop                           │
│  cmd.robot.g1.>           — commands TO Unitree G1                   │
│  cmd.robot.go2.>          — commands TO Go2                          │
│  cmd.ha.>                 — commands TO Home Assistant                │
│                                                                      │
│  telemetry.drone.recon.>  — FROM drone (GPS, battery, altitude)      │
│  telemetry.robot.g1.>     — FROM G1 (joints, IMU, cameras)           │
│  telemetry.infra.>        — FROM infrastructure (CPU, RAM, disk)     │
│                                                                      │
│  events.ha.>              — Home events (motion, door, temp)          │
│  events.security.>        — Security events (intrusion, anomaly)     │
│  events.fleet.>           — Agent lifecycle (connected, lost, error) │
│                                                                      │
│  agents.prompt.<agent>    — ACP: talk to any Hermes profile          │
│  agents.hb.<agent>        — heartbeat from each agent                │
│                                                                      │
│  $SRV.INFO / $SRV.PING   — service discovery (who's alive?)         │
│  $SYS.ACCOUNT.*.DISCONNECT — service death notifications             │
│                                                                      │
│  KV: agents, tasks, config, presence, missions                       │
│  JetStream: ACTIVITY, TELEMETRY, COMMANDS                            │
│  Object Store: firmware, configs, ML models                          │
│                                                                      │
└───────────────┬────────────────┬───────────────────┬────────────────┘
                │                │                   │
    ┌───────────▼──┐    ┌───────▼───────┐   ┌──────▼──────────┐
    │ HA-NATS      │    │ SaaS Worker   │   │ Edge Nodes      │
    │ Bridge       │    │ (Composio)    │   │                 │
    │              │    │               │   │ ┌─────────────┐ │
    │ Translates   │    │ saas.execute  │   │ │ Drone RECON │ │
    │ cmd.ha.*     │    │ → 2000+ apps  │   │ │ Python      │ │
    │ → HA REST    │    │               │   │ │ PyMAVLink   │ │
    │              │    │ Todoist,      │   │ │ NATS leaf   │ │
    │ events.ha.*  │    │ Gmail,        │   │ │ Local JS    │ │
    │ ← HA events  │    │ Airtable...  │   │ │ Camera CV   │ │
    └──────────────┘    └───────────────┘   │ └─────────────┘ │
                                            │                 │
                                            │ ┌─────────────┐ │
                                            │ │ Unitree G1  │ │
                                            │ │ Python      │ │
                                            │ │ Unitree SDK │ │
                                            │ │ NATS leaf   │ │
                                            │ │ ROS2 bridge │ │
                                            │ │ Locomotion  │ │
                                            │ └─────────────┘ │
                                            │                 │
                                            │ ┌─────────────┐ │
                                            │ │ Go2 (dog)   │ │
                                            │ │ Python      │ │
                                            │ │ NATS leaf   │ │
                                            │ │ Patrol mode │ │
                                            │ └─────────────┘ │
                                            └─────────────────┘
```

---

## The Voice → Action Flow (Real-Time)

Zeus speaks, things happen:

```
Zeus: "Legacy, send the drone to check the barn roof and 
       have the G1 meet me at the front door"

                    ┌──────────────────────────────────┐
                    │         LiveKit Voice             │
                    │  STT: Deepgram/Whisper → text     │
                    └──────────────┬───────────────────┘
                                   │ (< 200ms)
                                   ▼
                    ┌──────────────────────────────────┐
                    │         LEGACY (Hermes)           │
                    │                                   │
                    │  Parse: 2 commands detected       │
                    │  1. drone.recon → barn roof scan  │
                    │  2. robot.g1 → navigate to front  │
                    │                                   │
                    │  Execute (parallel):              │
                    │  ├── nats_publish(                │
                    │  │     "cmd.drone.recon.mission", │
                    │  │     {waypoints: [barn_roof],   │
                    │  │      camera: "nadir",          │
                    │  │      on_complete: "rtl"})      │
                    │  │                                │
                    │  └── nats_publish(                │
                    │        "cmd.robot.g1.navigate",   │
                    │        {destination: "front_door",│
                    │         gait: "walk",             │
                    │         on_arrival: "wait"})      │
                    │                                   │
                    │  TTS Response: "On it. Drone is   │
                    │  launching for barn roof. G1 is   │
                    │  heading to the front door now.   │
                    │  I'll update you when both are    │
                    │  in position."                    │
                    └──────────────┬───────────────────┘
                                   │
                    ┌──────────────▼───────────────────┐
                    │         NATS (instant)            │
                    │                                   │
                    │  cmd.drone.recon.mission ──────►  Drone
                    │  cmd.robot.g1.navigate ────────►  G1
                    │                                   │
                    │  (simultaneously subscribes to):  │
                    │  telemetry.drone.recon.>          │
                    │  telemetry.robot.g1.>             │
                    └──────────────────────────────────┘

                    ... 90 seconds later ...

                    ┌──────────────────────────────────┐
                    │  NATS Watch Service triggers:     │
                    │                                   │
                    │  telemetry.drone.recon.mission    │
                    │  → {status: "complete",           │
                    │     footage: "obj://recon/barn"}  │
                    │                                   │
                    │  telemetry.robot.g1.navigation    │
                    │  → {status: "arrived",            │
                    │     location: "front_door"}       │
                    └──────────────┬───────────────────┘
                                   │
                                   ▼
                    ┌──────────────────────────────────┐
                    │  LEGACY (follow-up, autonomous):  │
                    │                                   │
                    │  "Zeus, both done. Drone footage  │
                    │  shows 3 missing shingles on the  │
                    │  east side of the barn roof —     │
                    │  I've marked the frames. G1 is    │
                    │  at the front door waiting.       │
                    │  Want me to schedule a roofer?"   │
                    │                                   │
                    │  [Pushes annotated footage to     │
                    │   Flutter app + voice notification]│
                    └──────────────────────────────────┘
```

**Latency breakdown:**
- Voice → STT: ~200ms (Deepgram streaming)
- STT → Legacy intent parsing: ~500ms (LLM)
- Legacy → NATS publish: ~1ms
- NATS → Drone/Robot received: ~5ms (local) / ~50ms (over NetBird)
- Total voice-to-action: **< 1 second**


---

## PART 2: CURRENT PRODUCTION SETUP — Azure Unified Voice Infrastructure

*Source: VOICE-PROJECT.md (SO/03 Infrastructure)*

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


---

## PART 3: LOCAL REALTIME VOICE — Architecture & Build Plan

*Source: local-realtime-voice.md*

---

# Local Realtime Voice — Architecture & Build Plan

> Goal: A fully local, sovereign realtime voice agent on Zeuss-MBP (M2 Max, 32GB) — comparable to Gemini Live and OpenAI Realtime API but running entirely on-device with full tool calling into NATS and the OneMind OS fabric.

---

## Why Local

| Concern | Cloud (OpenAI/Gemini Realtime) | Local (this plan) |
|---------|-------------------------------|-------------------|
| Cost | ~$0.06/min (realtime API) | $0 — all on-chip |
| Privacy | Audio sent to OpenAI/Google | Never leaves the Mac |
| Latency | 100-400ms round trip | 50-150ms on Apple Silicon |
| Offline | Requires internet | Fully offline capable |
| Tool calling | OpenAI/Google tools only | Full NATS fabric + any Hermes tool |
| Sovereignty | Platform-owned | Zeus-owned |

The M2 Max Neural Engine + GPU makes this viable. 32GB unified memory is more than enough.

---

## Current LiveKit Setup (as of May 2026)

### What Is Running

The Legacy OS already uses LiveKit as its voice transport layer. This is documented fully in `LEGACY-OS-STATUS.md`. Summary:

```
┌─────────────────────────────────────────────────────────────────────┐
│                        LEGACY AI — THE FABRIC                        │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐   ┌──────────────────┐   ┌──────────────────────┐ │
│  │   HERMES     │   │    LIVEKIT       │   │       NATS           │ │
│  │  (The Brain) │   │  (Transport)     │   │  (Nervous System)    │ │
│  │              │   │                  │   │                      │ │
│  │ Memory       │   │ WebRTC voice     │   │ HA commands          │ │
│  │ Skills       │   │ Video streams    │   │ Drone telemetry      │ │
│  │ Sessions     │   │ Data channels    │   │ Watch events         │ │
│  │ Tools/MCP    │   │ Console mode     │   │ IoT devices          │ │
│  │ Soul/identity│   │ (local solo)     │   │ Services fabric      │ │
│  └──────────────┘   └──────────────────┘   └──────────────────────┘ │
│                           ~700 lines of glue                          │
│              hermes_livekit_transport.py (300 lines)                  │
│              hermes_watch_daemon.py (200 lines)                       │
│              hermes_perception.py (200 lines)                         │
└─────────────────────────────────────────────────────────────────────┘
```

### LiveKit Config (Current)

| Component | Detail |
|-----------|--------|
| Mode | Console mode (local, no server) — single user on same machine |
| Server | Not deployed — not needed for solo use |
| Transport | Direct mic → STT → Hermes → TTS → speaker pipeline |
| LLM | Hermes Agent (cloud — Claude Sonnet 4.6 via AWS Bedrock) |
| STT | Configured but not locked to one provider |
| TTS | ElevenLabs / configured provider |
| Frontend | Legacy UI — Next.js 15 PWA at `~/workspace/legacy-ui/` |
| Frontend URL | http://localhost:3000 |
| NATS WS | ws://localhost:9222 (NATS server WebSocket listener) |
| Data channels | Backend → frontend via LiveKit data topic `"ui"` |

### Frontend Architecture (Legacy UI)

```
~/workspace/legacy-ui/
├── package.json          — Next.js 15, livekit-client, @livekit/components-react, nats.ws
├── app/
│   ├── page.tsx          — 3-column grid (metrics | HUD orb | logs+watches)
│   └── api/
│       ├── metrics/route.ts   — CPU/MEM polled every 2s
│       └── agent/route.ts     — Agent text/status stub
├── components/
│   ├── HudOrb.tsx        — Canvas 60fps: arcs, pulses, particles (states: idle/listening/speaking/thinking)
│   ├── SystemPanel.tsx   — 5 MetricBars (CPU, RAM, etc.)
│   ├── ActivityLog.tsx   — Color-coded: SYS=cyan, You=white, Legacy=green, WATCH=orange, ERR=red
│   ├── WatchPanel.tsx    — Watch events (● idle / ◉ triggered / ○ dismissed)
│   ├── VoiceControls.tsx — LiveKit connect/mute/VAD
│   └── TextInput.tsx     — Monospace text fallback
└── lib/
    ├── types.ts           — AgentState, Metric, LogEntry, Watch
    └── nats.ts            — NATS WebSocket singleton (ws://localhost:9222)
```

Color palette (Mark-XXXIX origin):
- bg: #00060a | panel: #010d14 | border: #0d3347
- primary: #00d4ff | accent: #ff6b00 | green: #00ff88 | red: #ff3355
- text: #8ffcff | text-dim: #3a8a9a

### The Gap

The current LiveKit setup uses Hermes as the brain, which routes to Claude Sonnet 4.6 (cloud). That means:
- Every voice turn makes a network call to AWS Bedrock
- Realtime "feel" is limited by cloud roundtrip
- Audio is sent to cloud STT (Deepgram or similar)

This plan closes that gap by moving STT, LLM, and TTS all on-device.

---

## The Local Realtime Voice Stack (Proposed)

### Hardware Context

| Spec | Value |
|------|-------|
| Machine | Zeuss-MBP — Apple M2 Max |
| RAM | 32GB unified (CPU + GPU share) |
| Neural Engine | 16-core, optimized for transformer inference |
| Storage | ~478GB free |

### Component Selection

#### STT — Speech to Text

**Winner: mlx-whisper (Apple Silicon optimized)**

- Runs entirely on M2 Max Neural Engine
- Whisper Turbo model: ~1.5GB, ~80ms latency on M2
- OpenAI-compatible API surface
- RAM: ~1.5GB loaded

```bash
pip install mlx-whisper
# model: mlx-community/whisper-turbo-mlx (recommended)
# or: mlx-community/whisper-large-v3-turbo-mlx for higher accuracy
```

Fallback: faster-whisper (CPU, portable, ~200ms on M2)

#### LLM — Tool Calling Brain

**Winner: Qwen2.5-7B-Instruct (MLX 4-bit) via Ollama**

Best-in-class tool calling at small size. Tested against Mistral, Llama3, Gemma. Qwen2.5 wins on structured function calling accuracy.

- Model size: ~4.5GB (4-bit quantized)
- RAM footprint: ~5GB loaded
- Tokens/sec on M2 Max: ~60-80 tok/s
- Tool calling: native JSON schema function calling (OpenAI format)
- Serves as OpenAI-compatible endpoint at http://localhost:11434/v1

```bash
brew install ollama
ollama pull qwen2.5:7b
# Starts server at localhost:11434 — OpenAI-compatible
```

Alternative if RAM is tight: Llama-3.2-3B (~2.5GB) — less accurate tool calling but faster.

#### TTS — Text to Speech

**Winner: Kokoro TTS**

- Most natural-sounding local TTS as of 2026
- Model: ~300MB (tiny footprint)
- Latency: ~50ms first token on M2
- ONNX runtime — no GPU required, runs on CPU efficiently
- Multiple voices available

```bash
pip install kokoro-onnx sounddevice
```

Fallback: Piper TTS (~40MB models, more robotic but extremely fast)

#### VAD — Voice Activity Detection

**Silero VAD** (ships with livekit-agents[silero])

- Detects when Zeus starts/stops speaking
- Handles turn detection for interruption
- ~30MB, CPU-only, near-zero latency

#### Orchestration: LiveKit Agents (Python)

Same framework already in the livekit-voice-agent skill. Console mode handles the full pipeline:

```
Mic → Silero VAD → mlx-whisper STT → Qwen2.5 LLM → Kokoro TTS → Speaker
                        ↑
                  Tool calls dispatched to NATS fabric
```

### RAM Budget

| Component | RAM |
|-----------|-----|
| Qwen2.5-7B (4-bit, Ollama) | ~5.0 GB |
| mlx-whisper turbo | ~1.5 GB |
| Kokoro TTS | ~0.5 GB |
| Silero VAD | ~0.1 GB |
| Python agent process | ~0.5 GB |
| **Total voice stack** | **~7.6 GB** |
| macOS + Hermes + other services | ~12 GB |
| **Available headroom** | **~12 GB** |

Comfortable on 32GB. No swapping expected.

### Latency Estimate (Console Mode, Local Only)

| Stage | Estimated Latency |
|-------|------------------|
| VAD: detect speech end | ~100ms |
| STT: transcribe utterance | ~80ms (whisper-turbo, M2 Max) |
| LLM: first token | ~150ms (Qwen2.5-7B, Ollama) |
| LLM: full response (20 tok) | ~250ms |
| TTS: first audio chunk | ~50ms (Kokoro) |
| **Total to first audio** | **~630ms** |

Comparable to Gemini Live (~400-600ms observed). Under 1 second is the target. Achievable.

---

## Tool Calling Architecture

Tool calls flow from Qwen2.5 → agent.py → NATS → OneMind OS fabric. Same NATS subjects already in use.

```python
# Tools defined as OpenAI-format JSON schemas
tools = [
    {
        "type": "function",
        "function": {
            "name": "control_lights",
            "description": "Turn lights on/off, set brightness or color",
            "parameters": {
                "type": "object",
                "properties": {
                    "entity_id": {"type": "string"},
                    "action": {"type": "string", "enum": ["on", "off", "dim"]},
                    "brightness": {"type": "integer", "minimum": 0, "maximum": 100}
                }
            }
        }
    }
]

# Tool dispatch via NATS
async def execute_tool(name: str, args: dict) -> str:
    if name == "control_lights":
        response = await nc.request("cmd.ha.service.call", json.dumps({
            "domain": "light",
            "service": args["action"],
            "entity_id": args["entity_id"]
        }).encode(), timeout=3)
        return response.data.decode()
```

### Available Tool Categories (via NATS)

| Category | NATS Subject | Examples |
|----------|-------------|---------|
| Home Assistant | `cmd.ha.*` | Lights, locks, climate, sensors (185 entities) |
| SaaS actions | `saas.execute.*` | Todoist, GitHub, Webflow, Cloudflare |
| Agent reasoning | `cmd.agent.legacy.*` | Delegate complex queries to Hermes |
| Drone telemetry | `telemetry.drone.*` | Position, camera, battery |
| Presence | `zeus.presence.*` | Location, active context |

---

## Two Build Paths

### Path A — Console Mode (Fastest, Recommended First)

No server. Local mic/speaker only. Full pipeline on-device.

```bash
# 1. Install dependencies
pip install livekit-agents[silero] mlx-whisper kokoro-onnx sounddevice
brew install ollama

# 2. Pull model
ollama pull qwen2.5:7b

# 3. Run agent
python local_voice_agent.py console
```

What you get:
- Full-duplex voice (interrupt mid-response)
- Tool calling via NATS
- No network required once running
- Single device (Mac)

Build time: 1-2 hours

### Path B — LiveKit Server (Remote Devices)

Adds livekit-server binary so iPhone, glasses, DIMOS can connect.

```bash
# Additional steps after Path A working:
brew install livekit
livekit-server --dev   # dev mode: no SSL, no TURN needed on LAN/Tailscale
```

What you get (on top of Path A):
- iPhone can connect to the voice room
- Glasses can connect (Future: Brilliant Frame)
- DIMOS can join as a room participant
- Legacy UI (localhost:3000) shows live room state

Build time: Path A + 2-4 hours

---

## Agent Code Skeleton (agent.py)

```python
import asyncio
import json
import nats
from livekit.agents import AutoSubscribe, JobContext, WorkerOptions, cli, llm
from livekit.agents.voice_assistant import VoiceAssistant
from livekit.plugins import openai as lk_openai, silero

# Local Ollama endpoint (Qwen2.5-7B)
LOCAL_LLM = lk_openai.LLM(
    model="qwen2.5:7b",
    base_url="http://localhost:11434/v1",
    api_key="ollama",  # dummy key, Ollama doesn't need it
)

# Local Whisper via mlx-whisper (exposed as OpenAI-compatible STT)
# See: mlx-whisper --server flag or use livekit-agents STT plugin
LOCAL_STT = lk_openai.STT(
    base_url="http://localhost:9090/v1",  # mlx-whisper server
    api_key="local",
    model="whisper-turbo",
)

# Kokoro TTS
LOCAL_TTS = lk_openai.TTS(
    base_url="http://localhost:9091/v1",  # Kokoro OpenAI-compat server
    api_key="local",
    voice="af_sky",  # Kokoro voice name
)

# Tool definitions
TOOLS = [
    {
        "type": "function",
        "function": {
            "name": "control_home",
            "description": "Control home automation — lights, locks, climate",
            "parameters": {
                "type": "object",
                "properties": {
                    "command": {"type": "string"},
                    "entity": {"type": "string"},
                    "value": {}
                },
                "required": ["command", "entity"]
            }
        }
    },
    {
        "type": "function",
        "function": {
            "name": "create_task",
            "description": "Create a task in Todoist",
            "parameters": {
                "type": "object",
                "properties": {
                    "content": {"type": "string"},
                    "priority": {"type": "integer", "minimum": 1, "maximum": 4},
                    "project": {"type": "string"}
                },
                "required": ["content"]
            }
        }
    }
]

nc = None  # NATS connection

async def execute_tool(name: str, args: dict) -> str:
    global nc
    if nc is None:
        nc = await nats.connect("nats://localhost:4222")
    
    if name == "control_home":
        msg = await nc.request("cmd.ha.service.call", json.dumps(args).encode(), timeout=3)
        return msg.data.decode()
    elif name == "create_task":
        msg = await nc.request("saas.execute.todoist.create_task", json.dumps(args).encode(), timeout=5)
        return msg.data.decode()
    return f"Unknown tool: {name}"


async def entrypoint(ctx: JobContext):
    await ctx.connect(auto_subscribe=AutoSubscribe.AUDIO_ONLY)

    system_prompt = """You are Legacy AI — Zeus's personal AI operating system.
You control the smart home, manage tasks, and act on his behalf.
Be concise. Respond as if speaking aloud (no markdown, no lists).
When Zeus gives a command, execute it immediately and confirm briefly."""

    assistant = VoiceAssistant(
        vad=silero.VAD.load(),
        stt=LOCAL_STT,
        llm=LOCAL_LLM,
        tts=LOCAL_TTS,
        chat_ctx=llm.ChatContext().append(
            role="system",
            text=system_prompt
        ),
        fnc_ctx=llm.FunctionContext(),  # attach tools here
    )

    assistant.start(ctx.room)
    await asyncio.sleep(1)
    await assistant.say("Legacy online. Ready.", allow_interruptions=True)

    await asyncio.Event().wait()  # hold until room closes


if __name__ == "__main__":
    cli.run_app(WorkerOptions(entrypoint_fnc=entrypoint))
```

Run:
```bash
python agent.py console
```

---

## Installation Steps (Full)

```bash
# 1. Homebrew dependencies
brew install ollama

# 2. Pull the LLM
ollama pull qwen2.5:7b
# Verify: curl http://localhost:11434/v1/models

# 3. Python environment
python3 -m venv ~/.venvs/livekit-voice
source ~/.venvs/livekit-voice/bin/activate

# 4. Core packages
pip install \
  livekit-agents[silero] \
  livekit-plugins-openai \
  mlx-whisper \
  kokoro-onnx \
  sounddevice \
  nats-py

# 5. mlx-whisper server (wraps as OpenAI-compatible STT)
# Start: mlx_whisper_server --model mlx-community/whisper-turbo-mlx --port 9090

# 6. Kokoro TTS server (OpenAI-compatible)
# pip install kokoro-fastapi
# Start: python -m kokoro_fastapi --port 9091

# 7. Test
python agent.py console
```

---

## Integration with Existing Legacy OS

### NATS Bridge (Already Live)

The NATS subjects are already deployed. The local voice agent just needs to connect:

```python
nc = await nats.connect("nats://localhost:4222")
# Then use cmd.ha.*, saas.execute.*, cmd.agent.legacy.* as needed
```

### Hermes as Fallback Brain

For complex reasoning (multi-step plans, vault writes, research) that Qwen2.5-7B can't handle:

```python
# Route to Hermes for heavy reasoning
msg = await nc.request("cmd.agent.legacy.reason", user_input.encode(), timeout=30)
return msg.data.decode()
```

This gives the best of both worlds: fast local for simple commands, Hermes for complex tasks.

### Legacy UI (Frontend)

The existing Next.js PWA at `~/workspace/legacy-ui/` already has:
- `VoiceControls.tsx` — LiveKit connect/mute/VAD
- `HudOrb.tsx` — Visual state (idle/listening/speaking/thinking)
- `ActivityLog.tsx` — Shows tool call results

When the local voice agent connects to a LiveKit server room, the frontend shows live state automatically.

---

## Milestones

| Milestone | What | Estimated Time |
|-----------|------|---------------|
| M1 | Ollama + Qwen2.5 running, tool calls via NATS working | 1 hour |
| M2 | mlx-whisper serving locally (STT) | 30 min |
| M3 | Kokoro TTS serving locally | 30 min |
| M4 | LiveKit agent.py console mode — full pipeline end to end | 1 hour |
| M5 | Test: voice → HA tool call → spoken confirmation | 30 min |
| M6 | livekit-server deployed, Legacy UI connected | 2 hours |
| M7 | iPhone can join the voice room | 1 hour |

**Total to M5 (fully functional local voice with tools): ~3.5 hours**

---

## Risks & Mitigations

| Risk | Likelihood | Mitigation |
|------|-----------|-----------|
| Qwen2.5-7B tool calling accuracy insufficient | Low | Fallback to Hermes for complex tools |
| mlx-whisper latency spikes under load | Low | whisper-turbo is small and fast; M2 Max handles it |
| Kokoro TTS voice sounds unnatural | Medium | Evaluate voice options; fallback to Piper |
| Ollama RAM contention with other services | Low | 32GB budget has ~12GB headroom |
| LiveKit console mode mic permissions on macOS | Low | Standard macOS privacy prompt, one-time grant |

---

## Related Docs

- `LEGACY-OS-STATUS.md` — Full current system architecture and running services
- `legacy.md` — Legacy AI identity, capabilities, phase roadmap
- `agent-infrastructure-architecture.md` — Full agent fleet architecture
- Skill: `devops/livekit-voice-agent` — LiveKit build patterns, console mode, server mode, Hermes integration
- `~/workspace/legacy-ui/` — Frontend PWA source

---

*Created: 2026-05-27 | Owner: Zeus | Agent: Legacy AI*


---

## PART 4: MODEL RESEARCH — Local Realtime Voice AI (May 2026)

*Source: research-local-realtime-voice-ai-2026.md*

---

# Local Realtime Voice AI Models - Research (May 2026)

## Target: MacBook Pro M2 Max 32GB | Sub-300ms latency | 24/7 local | No API costs

---

## 🏆 TOP RECOMMENDATION: Moshi by Kyutai

### 1. Moshi (Kyutai Labs)
- **GitHub**: https://github.com/kyutai-labs/moshi (10.3k ⭐, active May 2026)
- **Architecture**: TRUE audio-to-audio. Full-duplex speech-text foundation model. Models TWO audio streams simultaneously (user + model). Uses Mimi neural codec (12.5 Hz, 80ms frame size).
- **Parameters**: 7B (Temporal Transformer) + small Depth Transformer
- **Quantized sizes**: 
  - MLX int4: ~4GB (fits easily in 32GB)
  - MLX int8: ~8GB
  - bf16: ~14GB
- **Latency**: 160ms theoretical (80ms Mimi + 80ms acoustic delay), **200ms practical on L4 GPU**. On M2 Max with MLX q4, expect ~250-350ms.
- **Apple Silicon**: ✅ **First-class MLX support** (`moshi_mlx` package). Also Rust backend with Metal support.
- **Full-duplex**: ✅ TRUE simultaneous listening + speaking (like GPT Realtime)
- **Barge-in/Interruption**: ✅ Native - models user audio stream continuously
- **Function calling**: ❌ Not built-in. Would need to be added via text stream interception (model generates "inner monologue" text alongside audio)
- **Production readiness**: Strong. Three backends (PyTorch, MLX, Rust). Active development. CC-BY 4.0 license.
- **Key advantage**: Only open-source model with TRUE full-duplex like GPT Realtime. Has MLX native support specifically for Mac.

---

### 2. Sesame CSM (Conversational Speech Model)
- **GitHub**: https://github.com/SesameAILabs/csm (14.6k ⭐, active May 2026)
- **Architecture**: Text+audio → audio. Uses Llama backbone + Mimi codec. Generates RVQ audio codes from text AND audio context.
- **Parameters**: 1B (based on Llama-3.2-1B)
- **Size**: ~2GB quantized, easily fits 32GB
- **True audio-to-audio**: Partial. Takes audio context + text → generates speech. It's a TTS with audio context awareness, NOT a full conversational model. The fine-tuned version powers Sesame's voice demo but that version is NOT open.
- **Apple Silicon**: ✅ Supports MPS device explicitly in code
- **Barge-in**: ❌ Not a streaming/full-duplex model
- **Function calling**: ❌
- **Latency**: Good for TTS generation but not designed for realtime conversation loop
- **Production readiness**: Good as TTS component. Now in HuggingFace Transformers (v4.52.1+).
- **Verdict**: Excellent TTS component, NOT a full conversational model. Could be paired with a fast STT for a pipeline.

---

### 3. GLM-4-Voice (Zhipu AI / THUDM)
- **GitHub**: https://github.com/THUDM/GLM-4-Voice
- **Architecture**: Speech-to-speech via discrete tokens. Three parts: Tokenizer (Whisper encoder + VQ), 9B Chat Model, Decoder (CosyVoice flow matching).
- **Parameters**: 9B main model + tokenizer + decoder
- **Size**: bf16 ~18GB, int4 ~5GB for main model + decoder overhead. Tight but feasible on 32GB.
- **True audio-to-audio**: ✅ Yes - end-to-end speech tokens in → speech tokens out. But uses intermediate text ("streaming thinking" architecture).
- **Apple Silicon**: ❌ **CUDA required** explicitly. No Metal/MLX port exists.
- **Barge-in**: Limited - streaming output but not full-duplex
- **Function calling**: ❌
- **Latency**: Claims low latency (20 tokens to start synthesis) but requires CUDA GPU
- **Production readiness**: Research demo quality. Chinese-focused but supports English.
- **Verdict**: Won't run on M2 Max without significant porting effort. Skip for this use case.

---

### 4. Ultravox (Fixie AI)
- **GitHub**: https://github.com/fixie-ai/ultravox
- **Architecture**: Audio-in → TEXT-out only. Multimodal LLM that understands speech directly (no ASR step) but outputs text, not audio.
- **Parameters**: Default 70B (Llama 3.3), also 8B variant
- **True audio-to-audio**: ❌ **Audio → Text only**. "As we evolve the model, we'll train it to emit speech tokens" (future plan as of their README).
- **Apple Silicon**: Not optimized for it. 70B won't fit; 8B could work but still text-out only.
- **Function calling**: ✅ Yes (via their Realtime platform/API)
- **Verdict**: NOT suitable - it's audio-in but text-out. Would still need TTS on top.

---

### 5. Hertz-dev (Standard Intelligence)
- **GitHub**: https://github.com/Standard-Intelligence/hertz-dev (1.8k ⭐, last push Jan 2025 - **STALE**)
- **Architecture**: TRUE full-duplex conversational audio base model. First open-source full-duplex model.
- **Parameters**: ~8B estimated (details sparse)
- **True audio-to-audio**: ✅ Full-duplex audio generation
- **Apple Silicon**: ❌ **CUDA 12.1 required**. No Metal/MLX port.
- **Barge-in**: ✅ Full-duplex by design
- **Function calling**: ❌ Base model only - no instruction following
- **Latency**: Experimental, not benchmarked for production
- **Production readiness**: ❌ Research prototype. **Abandoned** (no updates since Jan 2025). Blog post: si.inc/hertz-dev/
- **Verdict**: Dead project. Was promising but abandoned. No Mac support.

---

### 6. Orpheus TTS (Canopy Labs)
- **GitHub**: https://github.com/canopyai/Orpheus-TTS (6.2k ⭐, last push Dec 2025)
- **Architecture**: Text → Audio TTS only. Built on Llama-3B backbone.
- **Parameters**: 3B
- **True audio-to-audio**: ❌ **TTS only** (text → speech)
- **Latency**: ~200ms streaming, reducible to ~100ms with input streaming
- **Apple Silicon**: Not explicitly supported (uses vLLM which is CUDA-focused). Could potentially run via llama.cpp.
- **Function calling**: N/A (TTS only)
- **Verdict**: Excellent TTS component (~200ms latency) for a pipeline approach. NOT a conversational model. Could pair with Whisper for STT → LLM → Orpheus TTS pipeline.

---

### 7. Qwen2-Audio (Alibaba)
- **Architecture**: Audio understanding model. Audio-in → Text-out.
- **Parameters**: 7B
- **True audio-to-audio**: ❌ Audio comprehension only, text output
- **Apple Silicon**: Can run via transformers/MLX
- **Verdict**: Not suitable for voice-to-voice. Audio understanding only.

---

### 8. Ichigo (Homebrew Ltd)
- **GitHub**: https://github.com/homebrewltd/ichigo
- **Architecture**: Speech package with ASR (22M tokenizer), TTS (coming soon), and experimental Speech LLM. NOT audio-to-audio.
- **True audio-to-audio**: ❌ It's a toolkit/pipeline (ASR + LLM + future TTS)
- **Verdict**: Utility toolkit, not a conversational voice model.

---

### 9. Gazelle (Tincans AI)
- **Website**: https://tincans.ai/slm
- **Architecture**: Speech-language model. Audio-in → Text-out (like Ultravox).
- **True audio-to-audio**: ❌ Audio understanding → text generation
- **GitHub**: Appears to be closed/unavailable
- **Verdict**: Referenced by Ultravox as prior work. Not a viable standalone option.

---

## Infrastructure & Runtimes

### llama.cpp Voice Support
- **Status**: llama.cpp has added multimodal support but primarily for vision. No native audio-to-audio model support as of mid-2026.
- Some community efforts to run audio token models (like Orpheus) via llama.cpp's token generation, but no integrated audio pipeline.

### MLX Voice Models
- **Moshi MLX**: ✅ The primary option. Official `moshi_mlx` package with q4/q8 quantization.
- **WhisperKit**: Apple's optimized Whisper for MLX (STT only)
- **No other native audio-to-audio models have MLX ports** as of mid-2026.

### Cactus Compute Engine
- **GitHub**: https://github.com/cactus-compute/cactus (5.2k ⭐, active May 2026)
- **What it is**: On-device AI runtime focused on **transcription (STT)** and **tool-calling agents**. Claims <120ms on-device latency.
- **Models**: Runs LLMs (e.g., LiquidAI/LFM2-2.6B) and speech models. Has "Needle" - a 26M parameter tool-calling model.
- **Full conversational**: ❌ Primarily STT + text agents. NOT audio-to-audio.
- **macOS**: ✅ Cross-platform including macOS (`brew install cactus-compute/cactus/cactus`)
- **Verdict**: Could be useful as the STT component with tool-calling, but NOT a full voice-to-voice solution. The "Needle" tool-calling model is interesting for NATS dispatch.

---

## GPT Realtime 2.0 - The Gap Analysis

**What GPT Realtime does that locals can't (yet):**
1. **Full-duplex with intelligence**: Simultaneously listens, thinks, and speaks with GPT-4 level reasoning
2. **Native function calling mid-speech**: Can pause, call a function, and resume naturally
3. **Emotion/tone awareness**: Understands and generates paralinguistic cues
4. **Sub-200ms latency** with high quality
5. **Seamless barge-in**: User interrupts → model stops immediately and responds to new input
6. **Voice continuity**: Maintains consistent voice identity across long conversations

**What Moshi achieves:**
- ✅ Full-duplex (unique among open models)
- ✅ Barge-in / interruption
- ✅ Sub-300ms on Mac (MLX q4)
- ✅ Voice continuity
- ❌ Intelligence is limited (~7B model, not GPT-4 level reasoning)
- ❌ No native function calling
- ❌ Limited emotion control

**The remaining gap**: Intelligence + function calling. Moshi can converse but can't reason deeply or dispatch tools natively.

---

## 🎯 RECOMMENDED ARCHITECTURE for Combat Voice Command System

### Option A: Moshi (Best for natural conversation + low latency)
```
[Mic] → Moshi MLX (q4, full-duplex) → [Audio out + Text stream]
                                              ↓
                                    Text stream monitor
                                              ↓
                                    Command parser → NATS dispatch
```
- **Pros**: True full-duplex, sub-300ms, native Mac, single model
- **Cons**: Limited reasoning, no native tool calling, command parsing needs custom work
- **Memory**: ~4-6GB (q4 + Mimi codec)

### Option B: Pipeline (Best for intelligence + tool calling)
```
[Mic] → WhisperKit/Cactus STT (~100ms) → Local LLM with tools (Llama 3 8B, ~5GB)
                                                    ↓
                                          Tool calls → NATS dispatch
                                                    ↓
                                          Orpheus TTS/CSM (~200ms) → [Speaker]
```
- **Pros**: Full function calling, better reasoning, modular
- **Cons**: Higher total latency (~400-600ms), no true full-duplex, no barge-in without extra work
- **Memory**: ~12-15GB total

### Option C: Hybrid (RECOMMENDED for combat system)
```
[Mic] → Moshi MLX (always listening, full-duplex) → conversation
              ↓ (when command detected via inner monologue text)
         Command intercept → Local LLM tool-call router → NATS
              ↓
         Moshi continues speaking confirmation
```
- Use Moshi for the conversational interface (always-on, full-duplex)
- Monitor Moshi's "inner monologue" text stream for command patterns
- Route detected commands to a small tool-calling model (Cactus Needle 26M or similar)
- Dispatch via NATS
- **Memory**: ~8GB total
- **Latency**: ~250ms voice, +50ms for command routing

---

## Summary Table

| Model | True A2A | Full-Duplex | Mac/MLX | Fits 32GB | Latency | Tool Calling | Status |
|-------|----------|-------------|---------|-----------|---------|--------------|--------|
| **Moshi** | ✅ | ✅ | ✅ | ✅ (4GB q4) | ~250ms | ❌ (hackable) | **Production** |
| CSM/Sesame | Partial | ❌ | ✅ (MPS) | ✅ (2GB) | Good | ❌ | TTS only |
| GLM-4-Voice | ✅ | ❌ | ❌ CUDA | Tight | ? | ❌ | Research |
| Ultravox | ❌ (text out) | ❌ | ❌ | 8B yes | ? | ✅ | API product |
| Hertz-dev | ✅ | ✅ | ❌ CUDA | ? | ? | ❌ | **Dead** |
| Orpheus TTS | ❌ (TTS) | ❌ | ❌ | ✅ (3B) | 200ms | N/A | Good TTS |
| Qwen2-Audio | ❌ (text out) | ❌ | ✅ | ✅ (7B) | ? | ❌ | Comprehension |
| Cactus | ❌ (STT+agents) | ❌ | ✅ | ✅ | <120ms | ✅ (Needle) | Active |

---

## Final Verdict

**Moshi is the clear winner** for local realtime voice on M2 Max. It's the ONLY open-source model that:
1. Does true audio-to-audio (not a pipeline)
2. Supports full-duplex conversation
3. Has native MLX/Apple Silicon support
4. Fits easily in 32GB RAM
5. Achieves sub-300ms latency
6. Is actively maintained (May 2026)
7. Has production-quality backends (Rust with Metal, MLX, PyTorch)

The main gap vs GPT Realtime is function calling and reasoning depth. This can be bridged with the hybrid architecture (Option C) by intercepting Moshi's text stream and routing commands to a lightweight tool-calling model.


---

## PART 5: DEEP RESEARCH — Tool Calling & Reasoning in Local Voice AI

*Source: research-local-realtime-voice-ai-tool-calling.md*

---

# Deep Research: Local Realtime Voice AI with Native Tool Calling & Reasoning

**Date:** May 28, 2026  
**Target Hardware:** MacBook Pro M2 Max 32GB  
**Use Case:** Combat-ready voice command system dispatching to drones/robots/infrastructure via NATS

---

## Executive Summary

**The honest answer: No single open-source model today does ALL of audio-in → reasoning → tool calling → audio-out in a single native generation loop with barge-in support, locally on M2 Max.**

However, **two models come very close** and represent viable paths:

1. **Qwen2.5-Omni-3B** — Best bet for M2 Max. Audio-in, audio-out, strong reasoning, streaming. Lacks *native* tool calling but its text output stream can be parsed for tool calls.
2. **MiniCPM-o 4.5** — Full-duplex omnimodal, runs on Mac via llama.cpp/Ollama. 9B params but quantized versions exist.

---

## 1. The State of Audio-Native Models with Function Calling

### Does ANY open model do all five requirements natively?

| Requirement | Qwen2.5-Omni | MiniCPM-o 4.5 | Moshi | LLaMA-Omni | Mini-Omni 2 |
|---|---|---|---|---|---|
| Audio input directly | ✅ | ✅ | ✅ | ✅ | ✅ |
| Reason about it | ✅ (MMLU via speech) | ✅ | ❌ weak | ❌ weak | ❌ weak |
| Native tool/function calling | ❌ not trained | ❌ not trained | ❌ | ❌ | ❌ |
| Audio output directly | ✅ | ✅ | ✅ | ✅ | ✅ |
| Barge-in/interruption | ✅ (chunked) | ✅ (full-duplex) | ✅ (full-duplex) | ❌ | ❌ |

**Key finding: ZERO open models have been explicitly trained with tool-calling tokens in their audio generation loop.** This is the gap.

---

## 2. Model-by-Model Analysis

### Qwen2.5-Omni (⭐ TOP CANDIDATE)

- **Architecture:** "Thinker-Talker" — Thinker generates text tokens (reasoning), Talker converts to speech tokens in streaming fashion
- **Sizes:** 3B and 7B
- **License:** Apache 2.0 (7B), Qwen-Research (3B)
- **Audio-in → Audio-out:** ✅ End-to-end, streaming
- **Reasoning:** ✅ Matches text-input performance on MMLU/GSM8K via speech input
- **Tool calling:** ❌ NOT explicitly supported, BUT the Thinker outputs text tokens first — **you can intercept the text stream for tool calls before the Talker generates audio**
- **Barge-in:** ✅ Chunked input, immediate output
- **Memory (M2 Max 32GB feasibility):**
  - 3B BF16: ~18GB for 15s audio context — **FITS on M2 Max**
  - 7B BF16: ~31GB minimum — tight, needs quantization
  - 7B GPTQ-Int4/AWQ: ~50% VRAM reduction → ~15-16GB — **FITS**
  - 3B via MNN: 3.6GB memory peak on mobile — **easily fits**
- **Edge deployment:** MNN support confirmed, runs on Snapdragon mobile SoCs
- **Real-time:** Yes, architecture designed for chunked streaming

**WHY THIS IS THE BEST PATH:** The Thinker-Talker architecture is a natural fit for tool calling. The Thinker generates text reasoning tokens. You can:
1. Fine-tune special `<tool_call>` tokens into the Thinker's vocabulary
2. OR parse the text output stream with a lightweight tool-call detector
3. Execute the tool while the Talker is still generating "thinking" filler audio
4. Inject tool results back into the Thinker's context

### MiniCPM-o 4.5 (⭐ STRONG ALTERNATIVE)

- **Architecture:** End-to-end omnimodal (SigLip2 + Whisper-medium + CosyVoice2 + Qwen3-8B)
- **Size:** 9B total parameters
- **License:** Open-source
- **Full-duplex:** ✅ Can see, listen, and speak simultaneously — output doesn't block input
- **Proactive interaction:** ✅ Can initiate without being asked
- **Reasoning:** ✅ "Approaches Gemini 2.5 Flash" on vision/speech
- **Tool calling:** ❌ Not explicitly trained for it
- **M2 Max feasibility:** Supported via llama.cpp and Ollama. Int4 quantized GGUF available in 16 sizes. Runs on MacBook per their docs.
- **Barge-in:** ✅ Full-duplex streaming

**Key advantage:** Based on Qwen3-8B which already has strong tool-calling capabilities in text mode. The pathway from "text tool calling" to "audio tool calling" is architecturally shorter.

### Moshi (Kyutai)

- Full-duplex audio, uses Mimi codec
- **Reasoning:** Weak — 7B Helium base model, not competitive on reasoning benchmarks
- **Tool calling:** None
- **Verdict:** Rejected by user. Correct assessment.

### LLaMA-Omni

- Built on Llama-3.1-8B-Instruct
- Speech-in → Speech-out, low latency
- **No tool calling, no barge-in**
- Last updated 2024, appears abandoned

### Mini-Omni / Mini-Omni 2

- Small speech-to-speech models
- No tool calling capability
- Weak reasoning
- Not competitive

### Meta Spirit LM / Voicebox

- Research models, not designed for tool calling
- Voicebox is TTS only, not bidirectional
- No open variants with tool use

### Moonshine (Useful Sensors)

- **STT only** — fast, small ASR model
- Not a full omnimodal model, no generation

### Mistral/Ministral Audio

- No public audio-native models from Mistral as of May 2026

### DeepSeek Audio

- No public audio-native omnimodal model released

### InternLM / InternVL

- InternVL focuses on vision, no audio-native models with tool calling

---

## 3. How GPT Realtime Handles Function Calling (Architecture Insight)

Based on OpenAI's documentation and developer resources:

1. **Single unified model** generates interleaved token types: audio tokens, text tokens, and function-call tokens
2. When the model decides to call a function, it emits a `function_call` event mid-stream
3. Audio generation **pauses** (or generates filler)
4. Client executes the function and returns results
5. Model resumes with function results in context, continues generating audio

**Key architectural insight:** The model's vocabulary includes both audio codec tokens AND function-call structured tokens. It was trained on sequences where tool calls appear interleaved with audio generation.

**What a local equivalent needs:**
- A model trained on interleaved sequences: `[audio_tokens] [text_reasoning] [<tool_call>name(args)</tool_call>] [<tool_result>...</tool_result>] [audio_tokens]`
- Streaming decoder that can detect tool-call tokens and pause/route appropriately
- Audio codec (like SNAC, Mimi, or the DiT-based codec Qwen uses) for input/output

---

## 4. The Viable Architecture for Your Use Case

### Option A: Qwen2.5-Omni-3B + Tool Call Interception (RECOMMENDED)

```
┌─────────────────────────────────────────────┐
│  Microphone → Audio Chunker (streaming)      │
│       ↓                                      │
│  Qwen2.5-Omni Thinker (3B, quantized)       │
│       ↓ (text token stream)                  │
│  ┌─── Tool Call Detector ───┐                │
│  │  if <tool_call> detected │                │
│  │  → Execute via NATS      │                │
│  │  → Inject result         │                │
│  └───────────────────────────┘               │
│       ↓                                      │
│  Qwen2.5-Omni Talker (speech synthesis)      │
│       ↓                                      │
│  Speaker (streaming audio out)               │
└─────────────────────────────────────────────┘
```

**Why this works:**
- Thinker generates text tokens (where tool calls live) BEFORE Talker generates audio
- You intercept at the text level, execute tools, then let Talker speak the result
- Latency: Thinker is fast (text generation speed), tool execution via NATS is ~ms
- The Talker can generate filler/acknowledgment while tool executes

**Implementation steps:**
1. Deploy Qwen2.5-Omni-3B via MNN or transformers on M2 Max
2. Use the streaming API with chunked audio input
3. Monitor the Thinker's text output for tool-call patterns
4. Fine-tune with LoRA to emit structured tool calls (the base model already reasons well from speech — it just needs the output format trained)

### Option B: Fine-tune Qwen2.5-Omni for Native Tool Calling

Since the Thinker is essentially a strong LLM that already handles speech instruction following at MMLU-level quality:

1. Create training data: speech commands → text reasoning → `<tool_call>dispatch_drone(target="building_alpha")</tool_call>` → spoken confirmation
2. LoRA fine-tune the Thinker on this format
3. The Talker doesn't need modification — it just speaks whatever text the Thinker produces

**This is architecturally feasible because the Thinker-Talker split means tool calling happens in the text domain where it's well-understood.**

### Option C: MiniCPM-o 4.5 via Ollama

- Already supported by Ollama (`ollama run openbmb/minicpm-o2.6` and newer)
- llama.cpp support confirmed
- Full-duplex — true barge-in
- Based on Qwen3-8B which has tool calling in text mode
- Quantized to run on Mac

---

## 5. Fine-Tuning Approaches Being Explored NOW

### Adding tool-calling to Moshi
- **Feasible but suboptimal** — Moshi's 7B Helium base is weak at reasoning
- You'd be fine-tuning a weak reasoner, better to start with a strong one

### Adding audio heads to a strong LLM
- This is exactly what Qwen2.5-Omni and MiniCPM-o already did
- Starting from scratch is months of work with massive compute
- Better to fine-tune an existing omni model

### Active projects:
- **Alibaba/MNN** — optimizing Qwen2.5-Omni for edge devices (active, May 2025)
- **OpenBMB** — MiniCPM-o 4.5 with Ollama/llama.cpp support (active, Feb 2026)
- **vLLM** — Added audio output support for Qwen2.5-Omni (April 2025)

---

## 6. Commercial Local Options

| Platform | Audio-Native Models | Tool Calling | Status |
|---|---|---|---|
| Ollama | MiniCPM-o 2.6/4.5 (GGUF) | Not built-in for audio mode | Available now |
| LM Studio | No omnimodal models yet | N/A | Not available |
| Jan.ai | No omnimodal models | N/A | Not available |
| llama.cpp | MiniCPM-o support | Manual via grammar | Available |

**No commercial local platform currently combines audio-native + tool calling out of the box.**

---

## 7. Concrete Recommendation for Your System

### Phase 1 (Deploy NOW — days):
1. Run **Qwen2.5-Omni-3B** via MNN or transformers on M2 Max
2. Implement a **text-stream interceptor** that parses Thinker output for command patterns
3. Map recognized commands to NATS tool dispatches
4. Use structured prompting: "When the user gives a command, output it as JSON: {action, target, params}"
5. The model already reasons at MMLU-level from speech — it will understand complex commands

### Phase 2 (Fine-tune — weeks):
1. LoRA fine-tune Qwen2.5-Omni-3B's Thinker on your specific command vocabulary
2. Training data: recordings of voice commands → structured tool-call outputs
3. This gives you near-native tool calling without modifying the architecture

### Phase 3 (Full native — months):
1. Wait for Qwen3-Omni or equivalent that includes native tool calling in the training mix
2. Or contribute to the open-source effort to add tool-call tokens to omni models

---

## 8. Key Technical Details for M2 Max Deployment

### Qwen2.5-Omni-3B on M2 Max:
- **Memory:** 3.6GB via MNN (mobile-optimized), ~18GB via transformers BF16
- **Latency:** Sub-second first token on edge SoCs per MNN benchmarks
- **Audio codec:** DiT-based speech decoder (not Mimi/Encodec — proprietary to Qwen)
- **Streaming:** Chunked input supported natively

### Qwen2.5-Omni-7B-GPTQ-Int4 on M2 Max:
- **Memory:** ~15-16GB after quantization + CPU offload
- **Better reasoning** than 3B
- Requires CUDA for GPTQ — may need MLX port or use the AWQ variant

### MiniCPM-o 4.5 on M2 Max:
- GGUF quantized versions available (int4 = ~5-6GB)
- llama.cpp native support
- Full-duplex capable

---

## 9. The Gap and Who Will Fill It

The specific capability gap is: **no model has been trained with audio-interleaved tool-call tokens in its generation vocabulary.**

GPT Realtime achieved this by training on interleaved audio/tool sequences. The open-source community hasn't replicated this yet because:
1. Training data doesn't exist publicly (audio conversations with tool calls)
2. The audio codec integration adds complexity
3. It requires training from scratch or extensive fine-tuning of the generation head

**Most likely to fill this gap first:**
- Alibaba (Qwen team) — already have the architecture, just need tool-call training data
- OpenBMB (MiniCPM team) — full-duplex + Qwen3 base = close to ready
- A community LoRA on Qwen2.5-Omni that adds tool-call patterns to the Thinker

---

## 10. Bottom Line

**For your combat-ready voice command system today:**

Use **Qwen2.5-Omni-3B** with text-stream tool-call interception. It:
- ✅ Accepts audio directly (streaming, chunked)
- ✅ Reasons at MMLU-level from speech
- ✅ Can be prompted/fine-tuned to emit tool calls in text stream
- ✅ Outputs audio directly (streaming speech synthesis)
- ✅ Supports interruption via chunked input architecture
- ✅ Fits on M2 Max (3.6GB via MNN, ~18GB via transformers)
- ✅ Open weight, Apache 2.0 license
- ✅ Active development (last update June 2025)

The tool-call interception happens at the Thinker's text output — it's fast, deterministic, and doesn't add latency because the Thinker generates text tokens before the Talker generates audio. Your NATS dispatch happens in parallel with audio generation.


---

## PART 6: QWEN2.5-OMNI LOCAL SETUP GUIDE — MacBook Pro M2 Max 32GB

*Source: qwen2.5-omni-local-setup-guide.md*

---

# Qwen2.5-Omni Local Setup Guide — MacBook Pro M2 Max 32GB

## Executive Summary

**Can you run it on M2 Max 32GB?**
- **3B model via MNN**: ✅ YES — 3.6GB memory, runs well on Apple Silicon
- **3B model via transformers (BF16)**: ⚠️ Barely — needs ~18GB for 15s video input
- **7B model via transformers**: ❌ NO — needs 31GB+ VRAM (GPU), won't fit in unified memory with overhead
- **7B GPTQ-Int4/AWQ**: ⚠️ Maybe text-only, but MPS support is NOT confirmed
- **GGUF via llama.cpp**: ⚠️ GGUF exists but audio generation likely won't work (text-only LLM path)

**Best path for M2 Max 32GB: MNN 3B model** (4-bit quantized, 3.6GB memory, native CPU+Metal)

---

## 1. Official GitHub Repo Requirements

**Repo**: https://github.com/QwenLM/Qwen2.5-Omni

### Dependencies
```bash
pip install transformers==4.52.3
pip install accelerate
pip install qwen-omni-utils[decord] -U  # or without [decord] on non-Linux
pip install soundfile librosa
```

**System requirement**: `ffmpeg` must be installed.

**Python**: 3.10+ (implied by transformers 4.52.3 requirement)

**Note**: No `requirements.txt` or `setup.py` in the repo — it's primarily documentation + cookbooks.

### Key packages needed:
- `transformers>=4.52.3`
- `accelerate`
- `qwen-omni-utils` (their utility package for processing audio/image/video)
- `soundfile` (for writing WAV output)
- `librosa` (for loading audio)
- `torch` (PyTorch — see section 8)
- `ffmpeg` (system-level)

---

## 2. Model Sizes & Memory Requirements

### Model Variants on HuggingFace:
| Model | Files | Precision | Min VRAM (15s video) |
|-------|-------|-----------|---------------------|
| Qwen2.5-Omni-3B | 3 safetensors shards | BF16 | 18.38 GB |
| Qwen2.5-Omni-7B | multiple shards | BF16 | 31.11 GB |
| Qwen2.5-Omni-7B-GPTQ-Int4 | quantized | 4-bit | ~15 GB (est.) |
| Qwen2.5-Omni-7B-AWQ | quantized | 4-bit | ~15 GB (est.) |

### BF16 Memory by Input Length:
| Model | 15s Video | 30s Video | 60s Video |
|-------|-----------|-----------|-----------|
| 3B BF16 | 18.38 GB | 22.43 GB | 28.22 GB |
| 7B BF16 | 31.11 GB | 41.85 GB | 60.19 GB |

**Note**: These are with `flash_attention_2` — which does NOT work on MPS/Apple Silicon.

---

## 3. MNN Deployment (BEST PATH for Apple Silicon)

### The 3B MNN Model
- **Repo**: `taobao-mnn/Qwen2.5-Omni-3B-MNN` on HuggingFace
- **Quantization**: 4-bit
- **Memory Peak**: 3.6 GB (benchmarked on Snapdragon, similar on Apple Silicon)
- **License**: Apache 2.0

### Exact Steps to Build & Run on macOS:
```bash
# 1. Install dependencies
brew install cmake

# 2. Download the model
pip install huggingface_hub
huggingface-cli download taobao-mnn/Qwen2.5-Omni-3B-MNN --local-dir ~/Models/Qwen2.5-Omni-3B-MNN

# 3. Clone and build MNN
git clone https://github.com/alibaba/MNN.git
cd MNN
mkdir build && cd build
cmake .. -DMNN_LOW_MEMORY=true -DMNN_CPU_WEIGHT_DEQUANT_GEMM=true -DMNN_BUILD_LLM=true -DMNN_SUPPORT_TRANSFORMER_FUSE=true -DMNN_METAL=true
make -j$(sysctl -n hw.ncpu)

# 4. Run inference
./llm_demo ~/Models/Qwen2.5-Omni-3B-MNN/config.json prompt.txt
```

### MNN Benchmark (mobile SoC — expect similar or better on M2 Max):
- 3B Thinker Prefill: ~55 tok/s
- 3B Thinker Decode: ~16-23 tok/s
- 3B Talker Decode: ~52-62 tok/s
- 3B Code2Wav: ~28 tok/s

**M2 Max should be FASTER than Snapdragon 8 Elite** given its superior memory bandwidth (400 GB/s vs ~50 GB/s on mobile).

---

## 4. MLX Port

**Status: NO MLX port exists as of this research.**

GitHub search for "qwen2.5-omni mlx" returns no relevant repositories. The Thinker-Talker architecture with audio codec output is non-trivial to port to MLX.

---

## 5. GGUF / Ollama / llama.cpp

### GGUF Models Available (unsloth/Qwen2.5-Omni-7B-GGUF):
- Q4_K_M, Q4_K_S, Q5_K_M, Q5_K_S, Q6_K, Q8_0, BF16
- IQ4_NL, IQ4_XS, Q2_K, Q3_K_M, Q3_K_S
- Ultra-low: UD-IQ1_M, UD-IQ2_M, etc.

### Ollama: NOT AVAILABLE
Ollama does NOT have `qwen2.5-omni` in its library (returns 404).

### Critical Limitation:
⚠️ **GGUF/llama.cpp only handles the TEXT generation (Thinker) part.** The audio output (Talker + DiT codec) is a separate model path that llama.cpp does not support. You'd get text-only responses — no speech synthesis.

---

## 6. Audio Hardware & Libraries

### Input (Microphone):
- Any macOS-compatible microphone works
- Python libraries: `sounddevice`, `pyaudio`, or `soundfile`
- Audio format: 16kHz sample rate for input

### Output (Speaker):
- Standard macOS audio output
- Output sample rate: 24kHz WAV
- Python libraries: `soundfile` (write), `sounddevice` (playback)

### Required Python audio packages:
```bash
pip install soundfile librosa sounddevice pyaudio
# pyaudio needs portaudio:
brew install portaudio
```

---

## 7. Inference Speed on M2 Max

### No official M2 Max benchmarks exist.

**Estimated based on MNN mobile benchmarks + M2 Max specs:**
- M2 Max has ~400 GB/s memory bandwidth vs ~50 GB/s on Snapdragon 8 Elite
- MNN 3B model (4-bit, 3.6GB): expect **60-100+ tok/s decode** on M2 Max
- The bottleneck will be the audio codec (Code2Wav), not the LLM

**For transformers path (if it works on MPS):**
- Expect significantly slower than GPU — likely 5-15 tok/s for 3B on MPS
- No flash_attention_2 on MPS = much higher memory usage

---

## 8. PyTorch & MPS Support

### PyTorch Required: YES
- Minimum: PyTorch 2.1+ (for transformers 4.52.3 compatibility)
- Recommended: `torch>=2.3` for better MPS support

### MPS (Metal Performance Shaders):
- **Theoretically works** with `device_map="mps"` or manual `.to("mps")`
- **NOT officially tested/supported** by the Qwen team (they only test CUDA)
- `flash_attention_2` does NOT work on MPS — must use default attention
- Likely issues with custom ops in the Talker/DiT audio generation

### Install PyTorch for macOS:
```bash
pip install torch torchvision torchaudio
# This automatically includes MPS support on Apple Silicon
```

---

## 9. Streaming/Chunked Input API — Actual Code

### Basic inference (audio in → text + audio out):
```python
import soundfile as sf
from transformers import Qwen2_5OmniForConditionalGeneration, Qwen2_5OmniProcessor
from qwen_omni_utils import process_mm_info

# Load model (3B fits in 32GB RAM on CPU/MPS)
model = Qwen2_5OmniForConditionalGeneration.from_pretrained(
    "Qwen/Qwen2.5-Omni-3B",
    torch_dtype="auto",
    device_map="auto"  # will use MPS on Mac
)
processor = Qwen2_5OmniProcessor.from_pretrained("Qwen/Qwen2.5-Omni-3B")

def voice_chat(audio_path):
    messages = [
        {"role": "system", "content": [
            {"type": "text", "text": "You are Qwen, a virtual human developed by the Qwen Team, Alibaba Group, capable of perceiving auditory and visual inputs, as well as generating text and speech."}
        ]},
        {"role": "user", "content": [
            {"type": "audio", "audio": audio_path},
        ]},
    ]
    text = processor.apply_chat_template(messages, tokenize=False, add_generation_prompt=True)
    audios, images, videos = process_mm_info(messages, use_audio_in_video=True)
    inputs = processor(text=text, audio=audios, images=images, videos=videos, 
                       return_tensors="pt", padding=True, use_audio_in_video=True)
    inputs = inputs.to(model.device).to(model.dtype)

    text_ids, audio = model.generate(**inputs, use_audio_in_video=True, return_audio=True)

    response_text = processor.batch_decode(text_ids, skip_special_tokens=True, 
                                           clean_up_tokenization_spaces=False)
    sf.write("output.wav", audio.reshape(-1).detach().cpu().numpy(), samplerate=24000)
    return response_text[0]

# Record from mic then call:
# result = voice_chat("my_recording.wav")
```

### Real-time mic recording + inference:
```python
import sounddevice as sd
import numpy as np
import soundfile as sf

# Record 5 seconds of audio
duration = 5  # seconds
sample_rate = 16000
print("Recording...")
audio_data = sd.rec(int(duration * sample_rate), samplerate=sample_rate, channels=1, dtype='float32')
sd.wait()
print("Done recording")

# Save to file
sf.write("/tmp/input.wav", audio_data, sample_rate)

# Run inference
response = voice_chat("/tmp/input.wav")
print(f"Response: {response}")

# Play output
output_audio, sr = sf.read("output.wav")
sd.play(output_audio, sr)
sd.wait()
```

### Note on "Streaming":
The real-time streaming interaction is only available via:
1. Qwen Chat web app (https://chat.qwen.ai/)
2. vLLM deployment (requires CUDA GPU)
3. MNN on-device (supports chunked input)

**There is no Python streaming API for transformers on local CPU/MPS.** The `model.generate()` call is batch — you feed complete audio and get complete response.

---

## 10. Known Issues on macOS / Apple Silicon

1. **flash_attention_2 NOT available** — must use default attention (higher memory)
2. **No official MPS testing** — the Qwen team only tests on NVIDIA GPUs
3. **`decord` not installable from PyPI on macOS** — must build from source or use torchvision fallback
4. **vLLM does NOT support macOS** — the deployment server requires CUDA
5. **Audio generation (Talker/DiT) may fail on MPS** — custom CUDA kernels may not have MPS equivalents
6. **Memory**: Without flash_attention_2, actual memory is ~1.2x the theoretical minimum
7. **GGUF path loses audio output** — only text generation works through llama.cpp

---

## Recommended Setup for M2 Max 32GB

### Option A: MNN (Best for real-time, lowest memory)
```bash
# Full setup
brew install cmake portaudio ffmpeg
pip install huggingface_hub sounddevice soundfile

# Download model
huggingface-cli download taobao-mnn/Qwen2.5-Omni-3B-MNN --local-dir ~/Models/Qwen2.5-Omni-3B-MNN

# Build MNN
git clone https://github.com/alibaba/MNN.git
cd MNN && mkdir build && cd build
cmake .. -DMNN_LOW_MEMORY=true -DMNN_CPU_WEIGHT_DEQUANT_GEMM=true -DMNN_BUILD_LLM=true -DMNN_SUPPORT_TRANSFORMER_FUSE=true -DMNN_METAL=true
make -j$(sysctl -n hw.ncpu)

# Run
./llm_demo ~/Models/Qwen2.5-Omni-3B-MNN/config.json prompt.txt
```

### Option B: Transformers + 3B model (Full Python API, higher memory)
```bash
# Full setup
brew install ffmpeg portaudio
pip install torch torchvision torchaudio
pip install transformers==4.52.3 accelerate
pip install qwen-omni-utils soundfile librosa sounddevice

# Download and run (will auto-download from HF)
python -c "
from transformers import Qwen2_5OmniForConditionalGeneration, Qwen2_5OmniProcessor
model = Qwen2_5OmniForConditionalGeneration.from_pretrained('Qwen/Qwen2.5-Omni-3B', torch_dtype='auto', device_map='auto')
processor = Qwen2_5OmniProcessor.from_pretrained('Qwen/Qwen2.5-Omni-3B')
print('Model loaded successfully!')
"
```

⚠️ **Option B may fail on MPS for audio generation.** If it does, fall back to CPU:
```python
model = Qwen2_5OmniForConditionalGeneration.from_pretrained(
    "Qwen/Qwen2.5-Omni-3B", torch_dtype=torch.float32, device_map="cpu"
)
```
This will be slow (~2-5 tok/s) but functional.

---

## Download Sizes
- **Qwen2.5-Omni-3B** (HF, BF16): ~6 GB
- **Qwen2.5-Omni-3B-MNN** (4-bit): ~2 GB  
- **Qwen2.5-Omni-7B** (HF, BF16): ~15 GB
- **Qwen2.5-Omni-7B-GPTQ-Int4**: ~4-5 GB
- **Qwen2.5-Omni-7B-Q4_K_M.gguf**: ~4.5 GB (text-only, no audio out)


---

*Consolidated: 2026-06-05 | Owner: Zeus | Agent: Legacy AI*
