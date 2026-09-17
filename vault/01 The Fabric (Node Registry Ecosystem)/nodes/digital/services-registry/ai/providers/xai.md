---
title: "xAI"
type: reference
status: active
tags: [platform, ai, provider, xai, grok]
---
# xAI

> Elon Musk's AI company. Provider of the Grok model family — real-time web access, strong reasoning, multimodal.

## API
- **Base URL:** https://api.x.ai/v1
- **Docs:** https://docs.x.ai
- **Models:** Grok-3, Grok-3 Mini, Grok-2 Vision, Grok-2

## Key Characteristics
| Trait | Notes |
|-------|-------|
| Real-time data | Trained on live X/web data |
| Context window | 131k tokens (Grok-3) |
| Vision | Grok-2 Vision supports image input |
| Speed | Grok-3 Mini is fast + cost-efficient |
| Function calling | Supported |

## Integration with OneMind OS
- Hermes Agent: configured as `xai` provider via OpenAI-compatible API
- TTS: Aurora voice (xAI native)
- Tools: Grok (assistant), Grok Build (coding agent)

## Cost
Pay-per-token. Credentials in 1Password (SO vault → xAI).
