# llama-on-do

Meta Llama models hosted on DigitalOcean Gradient — zero egress cost when called from DO VPS.

| Model ID | Context | Price | Best Use |
|----------|---------|-------|---------|
| `llama3.3-70b-instruct` | 128K | ~$0.65/1M | **Primary** — best open-weight for agents |
| `llama3.1-70b-instruct` | 128K | ~$0.65/1M | Fallback |
| `llama3.1-8b-instruct` | 128K | ~$0.10/1M | Ultra-fast lightweight tasks |
| `llama3.2-90b-vision-instruct` | 128K | ~$0.90/1M | Vision + text (multimodal) |

**Key advantage:** Called from the DO VPS, these stay within DO's internal network — no internet roundtrip, no egress fees, ~50–100ms RTT vs ~300–600ms for external APIs.

**Stack use:** Grid agent (SO infra ops), bulk vault summaries, private inference where data shouldn't leave DO.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
