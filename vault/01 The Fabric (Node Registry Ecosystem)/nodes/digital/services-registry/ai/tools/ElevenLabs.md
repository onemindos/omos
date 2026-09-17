# ElevenLabs

> **Ultra-realistic AI voice synthesis and voice cloning**

## Overview

ElevenLabs provides state-of-the-art text-to-speech and voice cloning. Supports 29+ languages, custom voice creation, and real-time streaming audio. Used for giving OneMind agents a voice and producing audio content.

## Status

| Attribute   | Value                        |
| ----------- | ---------------------------- |
| Status      | 🔄 Evaluating                |
| Category    | AI / Voice                   |
| Type        | Voice AI Platform            |
| Integration | REST API, Python SDK, iOS    |

## Use Cases in OneMind

1. **Agent Voice** — Give Legacy or other agents a spoken voice output
2. **Content TTS** — Convert articles, ebooks, or notes to audio
3. **Voice Cloning** — Create consistent brand voice for One Mind content
4. **Vapi Integration** — Potential voice backend for Vapi calls

## Key Features

- 1000+ voices in 29 languages
- Custom voice cloning from samples
- Real-time streaming TTS
- Voice design (create from description)
- Dubbing and translation
- Conversational AI voices

## Quick Start

```python
from elevenlabs import ElevenLabs

client = ElevenLabs(api_key="ELEVENLABS_API_KEY")
audio = client.generate(
    text="Hello from One Mind",
    voice="Rachel",
    model="eleven_monolingual_v1"
)
```

## Links

- Website: [elevenlabs.io](https://elevenlabs.io)
- API Docs: [elevenlabs.io/docs](https://elevenlabs.io/docs)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
