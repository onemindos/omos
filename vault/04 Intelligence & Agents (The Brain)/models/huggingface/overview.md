# HuggingFace

**Type:** Model hub + inference API
**Serverless API:** `https://api-inference.huggingface.co/models/{model-id}`
**Key var:** `HF_API_TOKEN`
**Free tier:** Yes — rate-limited serverless inference

## What HuggingFace Is

HuggingFace is both a model registry and inference platform. 500,000+ models available. For OneMind OS, HF serves three roles:

1. **Model source** — download GGUF files for Ollama/LM Studio
2. **Serverless inference** — free/cheap API for open models without self-hosting
3. **Dedicated endpoints** — deploy any model on private GPU, billed per hour

## Inference API Options

| Type | Cost | Latency | Best Use |
|------|------|---------|---------|
| Serverless (free tier) | Free | Variable | Dev/testing |
| Serverless (PRO) | ~$0.10–1.00/1M | Variable | Low-volume production |
| Dedicated Endpoint | $0.60–6.00/hr | Low | High-volume self-managed |

## Key Model Categories

- **LLMs:** Llama, Qwen, Mistral, Phi, DeepSeek (download or API)
- **Embeddings:** sentence-transformers/all-MiniLM, nomic-embed-text
- **Image gen:** FLUX, Stable Diffusion, SDXL
- **ASR:** Whisper large-v3
- **Vision:** LLaVA, PaliGemma

## Stack Use Cases

- Free inference fallback when paid APIs are rate-limited
- Source for pulling models into Ollama: `ollama pull hf.co/{org}/{model}`
- Embedding models for pgvector RAG (sentence-transformers are best-in-class free)


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
