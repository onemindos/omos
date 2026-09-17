# Nova-3

| Attribute | Value |
|-----------|-------|
| **Model ID** | `nova-3` |
| **Type** | STT (Speech-to-Text) |
| **Latency** | ~300ms streaming |
| **Languages** | 30+ |
| **Tier** | STACK |

**Best for:** Real-time voice transcription — fastest + most accurate Deepgram model. Primary STT layer for voice agents.

## Key Features

- **Streaming:** WebSocket streaming with partial results — words appear as spoken
- **Smart formatting:** Auto-punctuation, numerals, capitalization
- **Diarization:** Speaker separation (who said what)
- **Keyterm boosting:** Boost recognition accuracy for custom vocabulary
- **Language detection:** Auto-detect spoken language

## Models

| Model ID | Notes |
|----------|-------|
| `nova-3` | Latest, best accuracy, recommended |
| `nova-3-medical` | Optimized for medical terminology |
| `nova-2` | Previous gen, stable fallback |
| `nova-2-phonecall` | Optimized for telephone audio (8kHz) |
| `enhanced` | Legacy high-accuracy model |
| `base` | Legacy fast/cheap model |

## Pricing

~$0.0043/min (Nova-3) — a 1-hour conversation costs ~$0.26.

## LiveKit Integration

```python
from livekit.plugins import deepgram
stt = deepgram.STT(model="nova-3", language="en-US")
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
