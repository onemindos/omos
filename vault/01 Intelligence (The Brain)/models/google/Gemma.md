---
id: so-01-ai-models-gemma
title: "AI Models — Gemma (Google Open-Source)"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/google, open-source]
ai_summary: "Google Gemma open-source model reference — Gemma 2B, 7B, 9B, 27B. Self-hostable on VPS or via Ollama/HuggingFace. Free inference, no API key required."
---

# Gemma (Google Open-Source)

> OPEN-SOURCE / SELF-HOSTED. Gemma models are Google's open-weight releases — run locally, on the VPS via Ollama, or via HuggingFace Inference API.
> No API key required for self-hosted. Zero cost beyond compute.

---

## Overview

Gemma is Google's family of open-weight models derived from the same research behind Gemini. Unlike Gemini (closed API), Gemma weights are freely downloadable and can run on the DO VPS via Ollama — zero inference cost, full privacy, no data sent to Google.

- **Self-hosted (VPS):** Ollama on DigitalOcean VPS
- **HuggingFace API:** `https://api-inference.huggingface.co`
- **Weights:** https://huggingface.co/google
- **Key Location:** HuggingFace token in 1Password (free tier works); no key for local Ollama
- **License:** Gemma Terms of Use (commercial use permitted with restrictions)

---

## Models

### Gemma 2 Family (Recommended)

| Model | Params | VRAM Required | Best Use |
|-------|--------|--------------|---------|
| `gemma2:2b` | 2B | ~2 GB | Embedded, very fast, lightweight |
| `gemma2:9b` | 9B | ~6 GB | **Sweet spot** — quality/speed on VPS |
| `gemma2:27b` | 27B | ~18 GB | High quality, needs bigger GPU |

### Gemma 3 Family (Latest)

| Model | Params | Context | Notes |
|-------|--------|---------|-------|
| `gemma3:1b` | 1B | 32K | Tiny, mobile-grade |
| `gemma3:4b` | 4B | 128K | Good multimodal support |
| `gemma3:12b` | 12B | 128K | **Best local option** — vision + text |
| `gemma3:27b` | 27B | 128K | Near-Gemini quality, large VRAM |

### CodeGemma

| Model | Notes |
|-------|-------|
| `codegemma:7b` | Code-specialized, strong on Python/JS/TS |
| `codegemma:7b-code` | Base model (no chat format) |

### PaliGemma (Vision)

| Model | Notes |
|-------|-------|
| `paligemma2:3b` | Image captioning, VQA, OCR |
| `paligemma2:10b` | Higher accuracy vision tasks |

---

## Running on VPS (Ollama)

```bash
# Install Ollama on VPS (if not already installed)
curl -fsSL https://ollama.ai/install.sh | sh

# Pull a model
ollama pull gemma3:12b
ollama pull gemma2:9b
ollama pull codegemma:7b

# Run inference
ollama run gemma3:12b "Explain the CODEX method"

# Serve as API (OpenAI-compatible on port 11434)
ollama serve
```

**OpenAI-compatible API via Ollama:**
```typescript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: 'ollama',  // placeholder — not validated
  baseURL: 'http://localhost:11434/v1',  // or VPS IP:11434
});

const response = await client.chat.completions.create({
  model: 'gemma3:12b',
  messages: [{ role: 'user', content: 'Hello' }],
});
```

**From Mac → VPS:**
```bash
# Access Ollama on VPS via SSH tunnel
ssh -L 11434:localhost:11434 root@104.248.4.212
# Then hit localhost:11434 from Mac
```

---

## VPS Sizing Guide (DO VPS)

| Model | DO Droplet | VRAM | Monthly |
|-------|-----------|------|---------|
| gemma2:2b | CPU-only // 8GB RAM | n/a | ~$48 |
| gemma2:9b | GPU (RTX 4000) | 8 GB | ~$162 |
| gemma3:12b | GPU (RTX 4000) | 12 GB | ~$162 |
| gemma3:27b | GPU (L40S) | 24 GB | ~$487 |

> Current DO VPS is CPU-only — gemma2:2b or gemma2:9b (quantized Q4) runs fine on CPU for non-latency-sensitive tasks.

---

## HuggingFace Inference API (Alternative)

```typescript
const response = await fetch(
  'https://api-inference.huggingface.co/models/google/gemma-2-9b-it',
  {
    method: 'POST',
    headers: { Authorization: `Bearer ${process.env.HF_TOKEN}` },
    body: JSON.stringify({ inputs: 'Hello' }),
  }
);
```

> HuggingFace free tier: rate-limited but no cost. Paid tier: $0.06/1M tokens for inference endpoints.

---

## Usage in One Mind Stack

| Use Case | Model | Notes |
|----------|-------|-------|
| **Private local inference** | gemma2:9b (Ollama) | Zero data egress, no API costs |
| **Code review/gen (offline)** | codegemma:7b | Runs on VPS without GPU |
| **Vision tasks (local)** | gemma3:12b | Image analysis without cloud API |
| **Embedded/edge scenarios** | gemma2:2b | Fast responses, minimal resources |
| **RAG summarization** | gemma2:9b | Summarize chunks for pgvector retrieval |

---

## Notes

- **Not the same as Gemini** — Gemma is open-weight, Gemini is closed-API. Both from Google.
- **License:** Gemma Terms of Use allow commercial use but prohibit certain applications (distillation into competing models, misrepresentation). Check before commercial deployment.
- **Quantization:** Ollama auto-quantizes to Q4_K_M by default — ~50% VRAM reduction with ~5% quality loss
- **Gemma 3 multimodal:** gemma3:4b+ supports image input natively via Ollama (pass base64 images)
- **OpenClaw integration:** Point an agent to `http://localhost:11434/v1` with model `gemma3:12b` — fully OpenAI-compatible


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Models (AI Models)/google (Google)/google|google]]
