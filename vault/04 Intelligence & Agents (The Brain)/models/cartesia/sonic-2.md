# Sonic 2

| Attribute | Value |
|-----------|-------|
| **Model ID** | `sonic-2` |
| **Type** | TTS |
| **Latency** | ~90ms time-to-first-audio |
| **Languages** | English + Multilingual |
| **Tier** | STACK |

**Best for:** Real-time voice agent output — fastest TTS available. Sub-100ms means speech starts before the user finishes reading the response mentally.

## Variants

| Model ID | Notes |
|----------|-------|
| `sonic-2` | Latest, fastest, recommended |
| `sonic-2-2025-03-07` | Pinned version (use for production stability) |
| `sonic-english` | English-only, optimized |
| `sonic-multilingual` | Multi-language, ~140ms |

## Voice Cloning

Cartesia supports instant voice cloning from a short audio clip — no training required. Clone a voice in the Cartesia dashboard, get a `voice_id`, use it in API calls.

## API Usage

```python
import cartesia
client = cartesia.Cartesia(api_key="CARTESIA_API_KEY")
audio = client.tts.bytes(
    model_id="sonic-2",
    voice={"id": "voice_id_here"},
    transcript="Hello, I'm Legacy.",
    output_format={"container": "raw", "encoding": "pcm_f32le", "sample_rate": 44100}
)
```

## LiveKit Integration

Sonic 2 is the default TTS in LiveKit's `cartesia` plugin — zero config when using LiveKit Agents framework.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
