# llama-on-groq

Meta Llama models on Groq LPU — fastest inference for these models anywhere.

| Model ID | Context | Input | Output | Notes |
|----------|---------|-------|--------|-------|
| `llama-3.3-70b-versatile` | 128K | $0.59/1M | $0.79/1M | **Primary** — best quality on Groq |
| `llama-3.1-8b-instant` | 128K | $0.05/1M | $0.08/1M | Fastest + cheapest — ~1500 tok/s |
| `llama-3.2-90b-vision-preview` | 128K | $0.90/1M | $0.90/1M | Vision — image + text |

**Speed benchmark:** llama-3.1-8b on Groq: ~1200–1500 tok/s vs ~100 tok/s on GPU cloud.

**Stack use:** Legacy agent real-time streaming, fast sub-agent loops where raw speed matters more than frontier quality.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
