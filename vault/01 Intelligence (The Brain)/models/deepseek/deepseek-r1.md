# DeepSeek-R1

| Attribute | Value |
|-----------|-------|
| **Model ID** | `deepseek-reasoner` (R1) |
| **Context** | 128K |
| **Input** | $0.55/1M |
| **Output** | $2.19/1M |
| **Thinking** | Yes — CoT reasoning tokens visible |
| **Tier** | OPEN-WEIGHT |

**Best for:** Open-source chain-of-thought reasoning — matches o1 on math/coding/logic benchmarks at 10–20% of the cost. Reasoning tokens are visible and billed separately.

## Distilled Variants

| Model | Speed | Notes |
|-------|-------|-------|
| `deepseek-r1-distill-llama-70b` | Fast | On Groq — best price+speed combo |
| `deepseek-r1-distill-qwen-32b` | Medium | Good reasoning, smaller footprint |
| `deepseek-r1-distill-qwen-7b` | Very fast | Ollama self-host, minimal resources |

R1 distills are reasoning-capable smaller models — 70B distill runs on Groq at LPU speed with R1-class reasoning.

## Self-Hosting

Full R1: Requires significant GPU/RAM — use quantized version via Ollama or run on Groq/Cerebras.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
