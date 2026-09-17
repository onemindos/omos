# ElevenLabs TTS Models

| Model ID | Latency | Quality | Best Use |
|----------|---------|---------|---------|
| `eleven_flash_v2_5` | ~75ms | Good | Real-time agent voice (lowest latency) |
| `eleven_turbo_v2_5` | ~250ms | Better | Low-latency with higher quality |
| `eleven_multilingual_v2` | ~400ms | Best | Premium narration, 29 languages |
| `eleven_monolingual_v1` | ~400ms | High | English-only, legacy |

## Recommendation by Use Case

**Voice agent (real-time):** `eleven_flash_v2_5` — ~75ms, acceptable quality for conversation

**Content/narration:** `eleven_multilingual_v2` — highest quality, use for final output not real-time

## Voice Cloning

- **Instant cloning:** 1+ minute of audio → cloned voice in seconds
- **Professional cloning:** Longer samples → higher fidelity
- Cloned `voice_id` works with any model above

## Sound Effects

Model ID: `eleven_sfx_v1` — generate any sound effect from text description.

```
"heavy rain on a metal roof" → audio file
"notification chime, gentle" → audio file
```

## API Usage

```python
from elevenlabs import ElevenLabs
client = ElevenLabs(api_key="ELEVENLABS_API_KEY")
audio = client.text_to_speech.convert(
    voice_id="JBFqnCBsd6RMkjVDRZzb",
    model_id="eleven_flash_v2_5",
    text="Hello, this is Legacy."
)
```

## Pricing

~$0.15–0.30/1K chars depending on plan. Flash is cheapest per character.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
