# Cerebras

**Type:** Cloud inference API (Wafer-Scale Engine — WSE hardware)
**API base:** `https://api.cerebras.ai/v1` (OpenAI-compatible)
**Key var:** `CEREBRAS_API_KEY`
**Free tier:** Yes — generous rate limits

## Why Cerebras

The CS-3 Wafer-Scale Engine is a single silicon wafer the size of a dinner plate. Entire 70B models fit on-chip — no inter-chip memory bottleneck. Result: **~2000+ tokens/sec** on Llama 70B, making it the world's fastest inference for that model class.

## Models

| Model ID | Context | Speed | Notes |
|----------|---------|-------|-------|
| `llama3.3-70b` | 128K | ~2000 tok/s | Fastest 70B inference available |
| `llama3.1-8b` | 128K | ~3000 tok/s | Fastest 8B inference |
| `llama-4-scout-17b-16e-instruct` | 128K | fast | Llama 4 MoE model |

## Stack Use Cases

- Sub-100ms first-token latency for agent steps
- Voice pipeline where latency directly impacts UX
- Streaming long outputs (fast even at 8K+ output tokens)
- Alternative to Groq when Groq rate limits are hit

## Comparison: Groq vs Cerebras

| | Groq | Cerebras |
|-|------|---------|
| Speed | ~500–1500 tok/s | ~2000+ tok/s |
| Models | More variety | Fewer (Llama-focused) |
| Free tier | Yes | Yes |
| Pricing | ~$0.05–0.59/1M | Competitive |


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
