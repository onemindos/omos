# Voice Pipeline — Credit-Optimized Architecture

> Universal TTS/STT for OpenClaw Gateway using Azure, Deepgram, and DigitalOcean credits.
> Replaces direct OpenAI TTS with a separated pipeline where every component is credit-covered.

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                    VOICE PIPELINE (Always-On)                     │
│                                                                  │
│  User speaks                                                     │
│       │                                                          │
│       ▼                                                          │
│  ┌─────────────┐    ┌──────────────┐    ┌────────────────────┐  │
│  │ Silero VAD  │───▶│ Deepgram STT │───▶│ Azure OpenAI LLM   │  │
│  │ (local,free)│    │ (Nova-3)     │    │ (gpt-4o / gpt-5.2) │  │
│  └─────────────┘    └──────────────┘    └────────┬───────────┘  │
│                                                   │              │
│                                                   ▼              │
│                                          ┌────────────────────┐  │
│                        User hears  ◀─────│ Azure OpenAI TTS   │  │
│                                          │ (tts-1, alloy)     │  │
│                                          └────────────────────┘  │
│                                                                  │
│  Infrastructure: LiveKit on DigitalOcean (self-hosted)           │
└──────────────────────────────────────────────────────────────────┘
```

## Cost Breakdown (Always-On Ambient — 43,200 min/month)

| Component | Provider | Rate | Monthly | Credits |
|-----------|----------|------|---------|---------|
| **VAD** | Silero (local) | Free | $0 | N/A |
| **STT** | Deepgram Nova-3 | $0.012-0.024/min | $518-1,037 | **Deepgram credits** |
| **LLM** | Azure OpenAI | per-token | varies | **Azure credits** |
| **TTS** | Azure OpenAI tts-1 | ~$0.015/1K chars | $150-300 | **Azure credits** |
| **Infra** | LiveKit on DO | $24-48/mo | $24-48 | **DO credits** |
| **Total** | — | — | **$0/mo** | All covered by credits |

### Fallback Chain (when credits expire)

| Priority | TTS Provider | Cost | Notes |
|----------|-------------|------|-------|
| 1 | Azure OpenAI tts-1 | ~ $0.015/1K chars | Primary (Azure credits) |
| 2 | Deepgram Aura | ~ $0.0065/1K chars | Cheapest cloud (Deepgram credits) |
| 3 | Amazon Polly | ~ $4/M chars | Bedrock credits |
| 4 | Piper (self-hosted) | $0/request | DO droplet only |

---

## Phase 1: Azure OpenAI TTS (Config-Only)

**Status**: Implemented in `openclaw.json.example`

### What Changed

- `messages.tts.openai.apiKey` → `${AZURE_OPENAI_API_KEY}` (was `${APIKEY}`)
- `messages.tts.openai.baseUrl` → `${AZURE_OPENAI_ENDPOINT}/openai/deployments/tts/audio/speech?api-version=2024-12-17`
- `messages.tts.openai.model` → `tts-1` (was `gpt-4o-mini-tts`)
- Same change in `channels.discord.voice.tts.openai`

### Azure Deployment Required

You need a TTS model deployed on your Azure OpenAI resource:

```bash
# Option 1: Azure CLI
az cognitiveservices account deployment create \
  --name onemind-openai \
  --resource-group <your-rg> \
  --deployment-name tts \
  --model-name tts-1 \
  --model-version "1" \
  --model-format OpenAI

# Option 2: Azure Portal
# Azure Portal → Azure OpenAI Studio → Deployments → Deploy model → tts-1
# Deployment name: "tts"
```

### VPS Deployment

```bash
# 1. Update openclaw.json on VPS with the new TTS config
ssh root@104.248.4.212

# Edit the production config
vim /home/onemind/.openclaw/openclaw.json

# Update messages.tts block:
# {
#   "provider": "openai",
#   "openai": {
#     "apiKey": "${AZURE_OPENAI_API_KEY}",
#     "baseUrl": "${AZURE_OPENAI_ENDPOINT}/openai/deployments/tts/audio/speech?api-version=2024-12-17",
#     "model": "tts-1",
#     "voice": "alloy"
#   }
# }

# Same for channels.discord.voice.tts.openai

# 2. Restart
systemctl restart openclaw

# 3. Verify
journalctl -u openclaw -f  # Watch for TTS errors
# Test: send a message in Discord that triggers voice
```

---

## Phase 2: Full LiveKit Voice Pipeline

### Prerequisites

- Phase 1 working (Azure TTS confirmed)
- Deepgram API key (create at https://deepgram.com)
- DigitalOcean droplet for LiveKit (or use existing VPS)

### Service: `voice-agent`

A Python service using `livekit-agents` that runs alongside OpenClaw:

```
infra/services/livekit/voice-agent/
├── agent.py          ← LiveKit voice pipeline agent
├── requirements.txt  ← Python dependencies
├── Dockerfile        ← Container build
└── .env.example      ← Required env vars
```

### NATS Integration

The voice agent bridges LiveKit ↔ NATS ↔ OpenClaw:

```
LiveKit Room (user audio)
       │
       ▼
Voice Agent (Silero VAD + Deepgram STT)
       │
       ▼ transcribed text
NATS: comms.voice.human.zeus
       │
       ▼ (legacy-bridge forwards to OpenClaw)
OpenClaw Legacy processes + responds
       │
       ▼ response text
NATS: comms.voice.agent.legacy
       │
       ▼ (voice agent receives)
Voice Agent (Azure OpenAI TTS)
       │
       ▼ synthesized audio
LiveKit Room → User hears
```

### Docker Compose

The voice agent is added to `infra/docker-compose.yml` under the `full` profile:

```yaml
voice-agent:
  build: ./services/livekit/voice-agent
  profiles: ["voice", "full"]
  depends_on: [nats]
  environment:
    LIVEKIT_URL, LIVEKIT_API_KEY, LIVEKIT_API_SECRET,
    DEEPGRAM_API_KEY, AZURE_OPENAI_API_KEY, AZURE_OPENAI_ENDPOINT,
    NATS_URL
```

---

## Phase 3: Fallback & Self-Hosted

### Deepgram Aura TTS (Secondary)

Deepgram's Aura TTS at $0.0065/1K chars — cheapest cloud TTS:

```python
from livekit.plugins import deepgram
tts = deepgram.TTS(model="aura-asteria-en")
```

### Amazon Polly (Tertiary)

Via litellm proxy to present Polly as OpenAI-compatible:

```python
# litellm config
model_list:
  - model_name: polly-tts
    litellm_params:
      model: polly/neural/Joanna
      aws_region_name: us-east-1
```

### Piper (Self-Hosted, Zero Cost)

Open-source offline TTS on DO droplet:

```bash
docker run -p 10200:10200 rhasspy/piper --voice en_US-lessac-medium
```

---

## Environment Variables

| Variable | Source | Purpose |
|----------|--------|---------|
| `AZURE_OPENAI_API_KEY` | 1Password (existing) | TTS + LLM authentication |
| `AZURE_OPENAI_ENDPOINT` | 1Password (existing) | Azure OpenAI base URL |
| `DEEPGRAM_API_KEY` | 1Password (new) | Speech-to-text |
| `LIVEKIT_URL` | 1Password (existing) | LiveKit server URL |
| `LIVEKIT_API_KEY` | 1Password (existing) | LiveKit auth |
| `LIVEKIT_API_SECRET` | 1Password (existing) | LiveKit auth |
| `NATS_URL` | .env.sa (existing) | NATS connection |

### Adding Deepgram to 1Password

```bash
op item create --category=login \
  --vault="00-24 SO (Sovereign Operations)" \
  --title="Deepgram" \
  "api-key[password]=your-deepgram-api-key"
```

Already added to `env.tpl`:
```
DEEPGRAM_API_KEY=op://aaeyris5z44pdqspk277u2oppi/Deepgram/api-key
```

---

## Voice Comparison (March 2026)

| Provider | Cost/1K chars | Quality | Latency | Credits? |
|----------|--------------|---------|---------|----------|
| **Azure OpenAI tts-1** | ~$0.015 | High | ~300ms | ✅ Azure |
| **Azure OpenAI tts-1-hd** | ~$0.030 | Very High | ~500ms | ✅ Azure |
| **Deepgram Aura** | ~$0.0065 | Good | ~200ms | ✅ Deepgram |
| **Amazon Polly Neural** | ~$0.016 | High | ~200ms | ✅ Bedrock |
| **OpenAI tts-1** | ~$0.015 | High | ~300ms | ❌ Paid |
| **ElevenLabs** | ~$0.099 | Premium | ~500ms | ❌ Expensive |
| **Piper (self-hosted)** | $0 | Medium | ~100ms | ✅ Free |
