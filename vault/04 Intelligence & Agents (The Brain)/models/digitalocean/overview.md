---
id: so-01-ai-models-digitalocean
title: "AI Models — DigitalOcean (Gradient)"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/digitalocean]
ai_summary: "DigitalOcean Gradient inference reference — open-weight models (Llama, Mistral) hosted on same DO infrastructure as VPS. Zero egress cost, private inference."
---

# DigitalOcean (Gradient)

> INFRA provider. Open-weight model inference on the same DigitalOcean infrastructure as the OpenClaw VPS.
> Key advantage: **zero data egress cost** — model runs in the same datacenter as the agents.

---

## Overview

DigitalOcean's Gradient platform provides serverless GPU inference for open-weight models. Since the One Mind VPS already runs on DigitalOcean, API calls to Gradient stay within DO's network — no external egress, faster latency, and no data leaves controlled infrastructure.

- **API Base URL:** `https://api.digitalocean.com/v1/inference` (or `https://api.inference.do`)
- **SDK:** OpenAI-compatible — use OpenAI SDK with custom `baseURL`
- **Key Location:** 1Password → SO vault → `DigitalOcean API Token` (same token as infra)
- **Console:** https://cloud.digitalocean.com/gen-ai
- **Docs:** https://docs.digitalocean.com/products/gen-ai/

---

## Models Available on Gradient

### Meta Llama

| Model | Context | Price ($/1M tokens) | Best Use |
|-------|---------|--------------------|---------|
| `llama3.3-70b-instruct` | 128K | ~$0.65 | **Primary** — best open-weight for agents |
| `llama3.1-70b-instruct` | 128K | ~$0.65 | Fallback |
| `llama3.1-8b-instruct` | 128K | ~$0.10 | Ultra-fast lightweight tasks |
| `llama3.2-90b-vision-instruct` | 128K | ~$0.90 | Vision + text (multimodal Llama) |

### Mistral

| Model | Context | Price ($/1M tokens) | Best Use |
|-------|---------|--------------------|---------|
| `mistral-7b-instruct` | 32K | ~$0.10 | Fast, cheap, code-capable |
| `mixtral-8x7b-instruct` | 32K | ~$0.27 | MoE — better quality than 7B |
| `mistral-nemo-12b-instruct` | 128K | ~$0.15 | Multilingual, function calling |

> Pricing approximate — check DO console. Gradient prices are generally 30–50% below Anthropic for comparable quality.

---

## API Setup

Gradient uses the OpenAI-compatible format:

```bash
# Use DigitalOcean Bearer token
export DO_API_TOKEN="dop_v1_..."
```

```typescript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.DO_API_TOKEN,
  baseURL: 'https://api.inference.do/v1',
});

const response = await client.chat.completions.create({
  model: 'llama3.3-70b-instruct',
  messages: [{ role: 'user', content: 'Analyze today's infrastructure status' }],
});
```

**Alternative: Direct REST**
```bash
curl -X POST https://api.inference.do/v1/chat/completions \
  -H "Authorization: Bearer $DO_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"model": "llama3.3-70b-instruct", "messages": [{"role":"user","content":"Hello"}]}'
```

---

## Usage in One Mind Stack

| Agent | Model | Why |
|-------|-------|-----|
| **Grid** (SO infra) | llama3.3-70b-instruct | Infrastructure ops on DO — zero latency |
| **Bulk summarization** | llama3.1-8b-instruct | Batch vault summaries — cheapest option |
| **RAG reranking** | mistral-nemo-12b | Fast reranker for pgvector results |
| **Private inference** | Any llama model | When data shouldn't leave DO network |

---

## Why This Matters

```
External API call:    [VPS] → [Internet] → [Anthropic/OpenAI] → [Internet] → [VPS]
Gradient call:        [VPS] → [DO internal network] → [Gradient GPU] → [VPS]
```

For Grid agent running infra checks every few minutes, Gradient is:
- **Faster:** ~50–100ms RTT vs ~300–600ms for external APIs
- **Cheaper:** No data transfer fees within DO datacenter
- **Private:** Sensitive infra data stays within DO

---

## Notes

- **auth**: Same API token used for DO droplets, volumes, and DNS — unified key management
- **Model list:** Check `GET https://api.inference.do/v1/models` for currently available models
- **GPU availability:** Gradient uses serverless GPUs — cold start ~1-2s on first call, warm after
- **No fine-tuning** via Gradient (inference only) — use DO's dedicated GPU droplets for fine-tuning
- **VPS region:** If VPS is in NYC1, prefer pulling from NYC inference endpoint for minimum latency


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
