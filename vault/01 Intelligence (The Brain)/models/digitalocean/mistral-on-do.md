# mistral-on-do

Mistral models hosted on DigitalOcean Gradient.

| Model ID | Context | Price | Best Use |
|----------|---------|-------|---------|
| `mistral-7b-instruct` | 32K | ~$0.10/1M | Fast, cheap, code-capable |
| `mixtral-8x7b-instruct` | 32K | ~$0.27/1M | MoE — better quality than 7B |
| `mistral-nemo-12b-instruct` | 128K | ~$0.15/1M | Multilingual, function calling |

**Use case:** Fast reranking for pgvector RAG results, cheap fallback inference within DO network.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
