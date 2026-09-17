# Aura (Deepgram TTS)

Deepgram's text-to-speech voices — lower latency than ElevenLabs, higher latency than Cartesia.

| Voice ID | Gender | Accent | Notes |
|----------|--------|--------|-------|
| `aura-asteria-en` | Female | American | Warm, conversational — recommended |
| `aura-luna-en` | Female | American | Soft, calm |
| `aura-stella-en` | Female | American | Expressive |
| `aura-athena-en` | Female | British | Professional |
| `aura-hera-en` | Female | American | Confident |
| `aura-orion-en` | Male | American | Deep, authoritative |
| `aura-arcas-en` | Male | American | Neutral |
| `aura-perseus-en` | Male | American | Friendly |
| `aura-angus-en` | Male | Irish | Distinctive |
| `aura-orpheus-en` | Male | American | Resonant |
| `aura-helios-en` | Male | British | Formal |
| `aura-zeus-en` | Male | American | Bold |

## Pricing

~$0.015/1K characters — cheaper than ElevenLabs for high-volume.

## When to Use Aura vs Cartesia

- **Cartesia Sonic:** Real-time voice agents (lower latency ~90ms)
- **Aura:** When already using Deepgram for STT — single vendor, simpler billing

## API

```python
response = client.speak.rest.v("1").save("output.mp3",
    {"text": "Hello from Legacy."},
    {"model": "aura-asteria-en"}
)
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
