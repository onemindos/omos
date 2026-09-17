---
title: "Unified Azure Voice Architecture — One Mind"
domain: SO-02
id: unified-azure-voice-2026
type: architecture
status: active
created: 2026-05-11
updated: 2026-05-11
owner: zeus
tags: [voice, azure, realtime, discord, telegram, infrastructure, unified]
---

# Unified Azure Voice Architecture

**Objective:** Single Azure voice provider (OpenAI Realtime) across all surfaces — Discord voice channels, Telegram voice messages, OpenClaw gateway (nodes/native app), and native client apps.

**Why:** Consolidate 3 fragmented voice config trees into 1 authoritative source. Eliminate provider fragmentation (Google + Azure + ElevenLabs). Reduce maintenance. Ensure consistent voice/latency across all interfaces.

**Target Date:** 2026-05-15

---

## Current State (2026-05-11)

| Surface | Provider | Config Path | Status | Issue |
|---------|----------|-------------|--------|-------|
| **OpenClaw Gateway (realtime)** | Azure OpenAI Realtime | `talk.providers.openai.*` | ✅ Active | Working correctly |
| **Discord Voice Channels** | OpenAI Realtime (gpt-realtime-2) | `channels.discord.voice.*` | ✅ Configured | Separate config tree |
| **Telegram TTS** | Google Gemini 2.5 Flash | `messages.tts.providers.google` | ✅ Active | Wrong provider |
| **Nodes (macOS/iOS)** | Via gateway `talk.*` config | Inherits from gateway | ✅ Partial | Relies on gateway config |

**Problem:** Three different config sources for the same purpose.

---

## Unified Vision

```
                    Azure OpenAI (Single Source of Truth)
                    ┌──────────────────────────────────┐
                    │ onemind-voice.openai.azure.com   │
                    │ ├─ gpt-realtime-1.5   (primary)  │
                    │ ├─ gpt-realtime-mini  (fallback) │
                    │ ├─ gpt-4o-mini-tts    (backup)   │
                    │ └─ gpt-4o-transcribe  (ASR)       │
                    └─────────────┬──────────────────────┘
                                  │ AZURE_OPENAI_VOICE_KEY
                ┌─────────────────┼─────────────────┬─────────────────┐
                │                 │                 │                 │
         ┌──────▼────┐    ┌──────▼──────┐  ┌──────▼──────┐  ┌────────▼─────┐
         │ OpenClaw  │    │   Discord   │  │  Telegram   │  │ Native Apps  │
         │ Gateway   │    │  Voice CH   │  │  Voice Msg  │  │(iOS/Android) │
         │           │    │             │  │             │  │              │
         │ Realtime: │    │ Realtime:   │  │ TTS:        │  │ Talk Mode:   │
         │ gpt-realtime   │ gpt-realtime│  │ gpt-4o-mini │  │ via gateway  │
         │ -1.5      │    │ -1.5        │  │ -tts        │  │              │
         │           │    │             │  │             │  │              │
         │ Nodes     │    │ Voice conn: │  │ (discrete   │  │ Same as      │
         │ inherit   │    │ auto-join   │  │ voice msgs) │  │ OpenClaw     │
         │ this      │    │             │  │             │  │              │
         └───────────┘    └─────────────┘  └─────────────┘  └──────────────┘
         omos-core        Discord Guild    @onemindos_bot   Flutter/iOS nodes
         omos-sentinel    1453411552...    Group -1003905...
```

**One config change = all surfaces updated automatically.**

---

## Configuration Layers

### Layer 1: Master Voice Config (Top Level)

Location: `~/.openclaw/openclaw.json` on both nodes (omos-core + omos-sentinel)

```json
{
  "voice": {
    "provider": "azure-openai",
    "region": "eastus2",
    "resourceName": "onemind-voice",
    "apiVersion": "2025-04-01-preview",
    "apiKeyEnv": "AZURE_OPENAI_VOICE_KEY",
    "models": {
      "realtime": "gpt-realtime-1.5",
      "realtime_fallback": "gpt-realtime-mini",
      "tts": "gpt-4o-mini-tts",
      "transcribe": "gpt-4o-transcribe"
    },
    "defaults": {
      "voice": "azure-neural-voice-id",
      "language": "en-US",
      "sampleRate": 24000,
      "outputFormat": "audio-24khz-48kbitrate-mono-mp3"
    }
  }
}
```

### Layer 2: Per-Surface Inheritance

#### 2a. OpenClaw Gateway (`talk.*`)
```json
{
  "talk": {
    "provider": "openai",
    "providers": {
      "openai": {
        "model": "${voice.models.realtime}",
        "apiKey": { "source": "env", "id": "AZURE_OPENAI_VOICE_KEY" },
        "azureEndpoint": "https://${voice.resourceName}.openai.azure.com",
        "azureDeployment": "${voice.models.realtime}",
        "azureApiVersion": "${voice.apiVersion}",
        "voice": "${voice.defaults.voice}"
      }
    },
    "interruptOnSpeech": true,
    "silenceTimeoutMs": 1200
  }
}
```

#### 2b. Discord Voice Channels (`channels.discord.voice.*`)
```json
{
  "channels": {
    "discord": {
      "voice": {
        "enabled": true,
        "mode": "agent-proxy",
        "realtime": {
          "provider": "openai",
          "model": "${voice.models.realtime}",
          "voice": "${voice.defaults.voice}",
          "apiKey": { "source": "env", "id": "AZURE_OPENAI_VOICE_KEY" },
          "azureEndpoint": "https://${voice.resourceName}.openai.azure.com",
          "azureDeployment": "${voice.models.realtime}",
          "azureApiVersion": "${voice.apiVersion}"
        },
        "bargeIn": true
      }
    }
  }
}
```

#### 2c. Telegram Voice Messages (`messages.tts.*`)
```json
{
  "messages": {
    "tts": {
      "auto": "always",
      "provider": "azure-openai",
      "providers": {
        "azure-openai": {
          "model": "${voice.models.tts}",
          "voice": "${voice.defaults.voice}",
          "apiKey": { "source": "env", "id": "AZURE_OPENAI_VOICE_KEY" },
          "region": "${voice.region}",
          "resourceName": "${voice.resourceName}",
          "apiVersion": "${voice.apiVersion}",
          "outputFormat": "ogg-24khz-16bit-mono-opus",
          "lang": "${voice.defaults.language}"
        }
      }
    }
  }
}
```

#### 2d. Native Apps / Nodes (Flutter, iOS, Android)
Connected to OpenClaw gateway via `talk.*` config above. App receives:
- Realtime session endpoint from `talk.realtime.session` RPC
- Voice metadata from `voicewake.get` RPC
- Audio streams via `talk.speak` RPC (binary WebSocket frames)

---

## Azure Resources

### Deployed Endpoints

| Resource | Endpoint | Location | Deployments | Purpose |
|----------|----------|----------|-------------|---------|
| **onemind-voice** | `onemind-voice.openai.azure.com` | East US 2 | `gpt-realtime-1.5` ✅, `gpt-realtime-mini` ✅, `gpt-4o-mini-tts` ✅, `gpt-4o-transcribe` ✅ | Voice (realtime + TTS) |
| **onemind-openai** | `onemind-openai.openai.azure.com` | East US | `gpt-4o`, `gpt-4o-mini`, `gpt-5`, etc. | Chat models |

### Environment Variables

```bash
# On both omos-core and omos-sentinel
# ~/ .openclaw/gateway.systemd.env

AZURE_OPENAI_VOICE_KEY=<from 1Password>
AZURE_OPENAI_VOICE_REGION=eastus2
AZURE_OPENAI_VOICE_RESOURCE=onemind-voice
AZURE_OPENAI_VOICE_API_VERSION=2025-04-01-preview

# Also needed for chat fallback
AZURE_OPENAI_KEY=<separate chat key>
AZURE_OPENAI_ENDPOINT=https://onemind-openai.openai.azure.com
```

Retrieve from 1Password vault `aaeyris5z44pdqspk277u2oppi`:
- `AZURE_OPENAI_VOICE_KEY` → `op://aaeyris5z44pdqspk277u2oppi/[Item Name]/credential`
- Keep these separate from chat keys for security + independent rotation

---

## Voice Models & Use Cases

### Model Selection Matrix

| Model | Latency | Cost/min | Quality | Use Case | Status |
|-------|---------|----------|---------|----------|--------|
| **gpt-realtime-1.5** | 250ms avg | ~$0.21 | Highest | Discord voice channels, OpenClaw realtime, Flutter Talk mode | ✅ Primary |
| **gpt-realtime-mini** | 350ms avg | ~$0.09 | Good | Fallback (cost-sensitive), low-latency channels | ✅ Deployed |
| **gpt-4o-mini-tts** | N/A (TTS) | ~$0.10 | Good | Telegram voice messages, static TTS | ✅ Deployed |
| **gpt-4o-transcribe** | N/A (ASR) | ~$0.01 | Good | Voice message transcription (future) | ✅ Available |

**Recommendation:**
- **Default:** `gpt-realtime-1.5` (primary realtime model)
- **Fallback:** `gpt-realtime-mini` (if cost spikes or quota limits)
- **TTS (Telegram only):** `gpt-4o-mini-tts` (cheaper for discrete messages)
- **Transcription:** `gpt-4o-transcribe` (future ASR enhancement)

---

## Implementation Roadmap

### Phase 1: Consolidate Gateway Config (1 day)

**Goal:** Move `talk.*` config to unified `voice.*` structure with template substitution.

**Tasks:**
1. Create new `voice.*` top-level object in `openclaw.json`
2. Update `talk.providers.openai.*` to reference `${voice.*}` values
3. Verify gateway starts correctly with new config
4. Test Talk mode (macOS app, web browser)

**File:** `00-24 SO/03 Infrastructure/openclaw/gateway/AZURE-VOICE-CONFIG.md` (create)

---

### Phase 2: Update Discord Voice Config (1 day)

**Goal:** Align `channels.discord.voice.*` with unified voice structure.

**Tasks:**
1. Update Discord channel config to reference `${voice.*}`
2. Verify voice channels connect properly
3. Test barge-in and realtime audio

**Verification:**
```bash
ssh onemindos-core 'openclaw config get channels.discord.voice'
# Should show unified azure-openai reference
```

---

### Phase 3: Migrate Telegram TTS (1 day)

**Goal:** Switch from Google Gemini to Azure for Telegram voice messages.

**Tasks:**
1. Change `messages.tts.provider` from `google` to `azure-openai`
2. Update provider config with Azure credentials
3. Send test voice message in Telegram
4. Verify audio quality

**Config:**
```bash
ssh onemindos-core 'openclaw config set messages.tts.provider azure-openai'
ssh onemindos-core 'openclaw config set messages.tts.providers.azure-openai.outputFormat ogg-24khz-16bit-mono-opus'
ssh onemindos-core 'systemctl --user restart openclaw-gateway'
```

**Test:**
```bash
# In Telegram @onemindos_bot, send any message
# Legacy should reply with voice message (Opus format)
```

---

### Phase 4: Sync to Native Apps (2 days)

**Goal:** Flutter app inherits all voice settings from unified gateway config.

**What's already working:**
- iOS/macOS nodes read `talk.*` from gateway → auto-synced ✅
- Flutter app connects to gateway via `talk.realtime.session` RPC → auto-synced ✅

**What needs confirmation:**
1. Test Flutter app Talk Mode connects to Azure realtime
2. Verify latency < 300ms
3. Test interrupt-on-speech
4. Test across different network conditions

---

### Phase 5: Sync Both Nodes (1 day)

**Goal:** Ensure omos-sentinel has identical voice config for failover.

**Tasks:**
1. Copy unified `voice.*` config to sentinel
2. Verify environment variables are set
3. Test Talk mode failover: main node → sentinel node
4. Document failover procedure

---

### Phase 6: Documentation & Operations (1 day)

**Goal:** Document voice system for operators.

**Create:**
- `00-24 SO/03 Infrastructure/openclaw/voice/OPERATIONS.md` — how to change voice, diagnose issues, cost tracking
- `00-24 SO/03 Infrastructure/voice/ADMIN-COMMANDS.md` — CLI commands for voice management
- Update [this master config](#configuration-layers) with actual deployed values

---

## Cost Model

### Current (Fragmented)

| Surface | Model | $/min | Est. Usage | Est. $/mo |
|---------|-------|-------|-----------|-----------|
| Talk Mode (gateway) | gpt-realtime-1.5 | $0.21 | 10 min/day | ~$63 |
| Discord voice | gpt-realtime-1.5 | $0.21 | 5 min/day | ~$31 |
| Telegram voice | Google Gemini | $0.01 | 100 msg/day | ~$30 |
| **Total** | — | — | — | ~**$124** |

### Unified (Azure Only)

| Surface | Model | $/min | Est. Usage | Est. $/mo |
|---------|-------|-------|-----------|-----------|
| Talk Mode | gpt-realtime-1.5 | $0.21 | 10 min/day | ~$63 |
| Discord voice | gpt-realtime-1.5 | $0.21 | 5 min/day | ~$31 |
| Telegram voice | gpt-4o-mini-tts | $0.10 | 100 msg/day | ~$30 |
| **Total** | — | — | — | ~**$124** |

**Cost neutral** — same total. Benefit: **unified provider, reduced complexity.**

### Optimization Path (if needed)

If costs exceed budget, switch Discord/Talk to `gpt-realtime-mini`:
- Cost: $0.09/min (60% cheaper)
- Latency: +100ms (acceptable for most use cases)
- **Est. saving: ~$55/mo**

---

## Quick Reference: Commands

### View Current Voice Config
```bash
ssh onemindos-core 'openclaw config get voice'
ssh onemindos-core 'openclaw config get talk.providers.openai'
ssh onemindos-core 'openclaw config get channels.discord.voice'
ssh onemindos-core 'openclaw config get messages.tts.provider'
```

### Change Default Voice (all surfaces)
```bash
# Update master config
ssh onemindos-core 'openclaw config set voice.defaults.voice "azure-voice-id"'

# Restart gateway to apply
ssh onemindos-core 'systemctl --user restart openclaw-gateway'

# Verify all surfaces inherited change
ssh onemindos-core 'openclaw config get voice.defaults.voice'
```

### Switch to Cheaper Model (fallback)
```bash
ssh onemindos-core 'openclaw config set voice.models.realtime gpt-realtime-mini'
ssh onemindos-core 'systemctl --user restart openclaw-gateway'
```

### Test Realtime Voice (gateway)
```bash
# Via browser: openclaw dashboard → Talk button → speak
# Or via CLI (if available):
openclaw talk.test --text "Hello from Azure voice"
```

### Test Discord Voice
```bash
# In Discord #legacy-chat, enable voice and connect to a voice channel
# Join from web browser or iOS node
# Speak and verify audio is received
```

### Test Telegram Voice
```bash
# Send message to @onemindos_bot in Telegram
# Expect reply with voice message (Opus format)
# Verify audio quality
```

### Monitor Voice Usage
```bash
ssh onemindos-core 'cat ~/.openclaw/logs/voice-usage.log | tail -100'
# (log location may vary; check after first config)
```

---

## Troubleshooting

### Voice Not Working After Config Change

1. **Check environment variable is set:**
   ```bash
   ssh onemindos-core 'echo $AZURE_OPENAI_VOICE_KEY | head -c 20'
   # Should print first 20 chars of API key, not empty
   ```

2. **Verify gateway is running:**
   ```bash
   ssh onemindos-core 'systemctl --user status openclaw-gateway'
   # Should be active (running)
   ```

3. **Check logs for errors:**
   ```bash
   ssh onemindos-core 'journalctl --user -u openclaw-gateway -n 50 -f'
   # Look for Azure auth errors or model errors
   ```

4. **Restart gateway:**
   ```bash
   ssh onemindos-core 'systemctl --user restart openclaw-gateway'
   ```

### Audio Quality Issues

- **Choppy/glitchy audio:** Check sample rate matches model output (24kHz default). May be transport issue — test on different network.
- **High latency (> 500ms):** Fallback to `gpt-realtime-mini` or check network jitter. `talk.silenceTimeoutMs` may be too high.
- **Echo/feedback:** Verify `interruptOnSpeech: true` is set. On Discord, confirm echo cancellation is enabled.

### Cost Spike

1. Check usage logs
2. Identify high-usage surface (Talk Mode vs Discord vs Telegram)
3. Switch that surface to `gpt-realtime-mini` if acceptable
4. Monitor for 24 hours, verify cost impact

---

## Related Files

- `TALK-MODE-CONFIG.md` — Current (Google) config (to be deprecated after migration)
- `OPENCLAW-VOICE-RESEARCH.md` — Research doc (reference only)
- `VOICE-FIRST-SWARM-ARCHITECTURE.md` — Future vision (not yet implemented)
- `00-24 SO/03 Infrastructure/openclaw/nodes/conversal-voice-architecture.md` — Conversational voice patterns

---

## Next Steps

1. **Review this document** — validate the unified vision
2. **Phase 1 implementation** — consolidate gateway config (1 day)
3. **Test across all surfaces** — Talk Mode, Discord voice, Telegram voice (1 day)
4. **Document operations** — create runbook for Zeus and future operators (0.5 day)
5. **Archive old configs** — move Google/ElevenLabs refs to archive (0.5 day)

**Target completion: 2026-05-15**

---

## Appendix: Current Deployment (2026-05-11)

### omos-core (100.123.126.34)

**Active Config:**
```json
{
  "talk": {
    "provider": "openai",
    "providers": {
      "openai": {
        "model": "gpt-realtime-1.5",
        "voice": "coral",
        "apiKey": { "source": "env", "provider": "default", "id": "AZURE_OPENAI_VOICE_KEY" },
        "azureEndpoint": "https://onemind-voice.openai.azure.com",
        "azureDeployment": "gpt-realtime-1.5",
        "azureApiVersion": "2025-04-01-preview"
      }
    }
  }
}
```

**Status:** ✅ Gateway realtime voice working. Test: Open browser, click Talk, speak.

### omos-sentinel (100.66.235.38)

**Status:** ⏳ Assumed same config needed; not yet verified. Failover testing pending.

### Discord Guild (1453411552181883027)

**Voice Channels:** #legacy-chat voice room available
**Config:** Separate, needs alignment with unified structure
**Status:** ✅ Voice channels available; config fragmentation pending fix

### Telegram Group (-1003905828404)

**Voice Messages:** Supported via `[[audio_as_voice]]` tag
**Current TTS:** Google Gemini (should migrate to Azure)
**Status:** ⚠️ Works but uses wrong provider; migration pending

---

**Document Version:** 2026-05-11 (Initial)
**Author:** Zeus (via Legacy)
**Last Updated:** 2026-05-11
