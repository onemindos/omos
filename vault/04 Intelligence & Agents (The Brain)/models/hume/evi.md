# EVI (Empathic Voice Interface)

| Attribute | Value |
|-----------|-------|
| **Model** | EVI 3 (latest) |
| **Type** | Voice AI with emotional intelligence |
| **Protocol** | WebSocket (real-time) |
| **Tier** | STACK |

**Best for:** Voice conversations that adapt to emotional state — EVI detects tone, stress, excitement, hesitation in the user's voice and adjusts its responses and vocal delivery accordingly.

## Capabilities

- **Expression Measurement:** Detects 48 emotional dimensions from voice in real time
- **Adaptive TTS:** EVI's own voice changes warmth, pacing, and tone based on detected emotions
- **Full pipeline:** STT → LLM → Expression analysis → emotionally-calibrated TTS (all in one API)
- **Custom system prompt:** Inject agent persona and behavior rules

## EVI vs Standard Voice Stack

| | EVI | Deepgram + LLM + Cartesia |
|-|-----|---------------------------|
| Emotional awareness | Built-in | None |
| Setup complexity | Single API | Three APIs |
| Voice quality | Good | Excellent (Cartesia) |
| Latency | ~500ms | ~300ms total |
| Best for | Empathic coaching | Speed-critical agents |

## Expression Measurement API (Standalone)

Separate API — analyze any audio file for emotions without the full voice pipeline.

```
POST /v0/batch/jobs (audio URL) → 48 emotion scores per time segment
```

Use case: analyze recorded calls/meetings to understand emotional trajectory.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
