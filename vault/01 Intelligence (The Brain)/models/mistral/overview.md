---
id: so-01-ai-models-mistral
title: "AI Models — Mistral AI"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/mistral]
ai_summary: "Mistral AI model reference — French AI company (Paris), La Plateforme API. Mistral Large 3, Mistral Small 4, Codestral, Devstral, Magistral reasoning, Voxtral audio, OCR. OpenAI-compatible."
---

# Mistral AI

> STACK provider. **French AI company** (Paris, founded 2023). Home of Mistral Large, Codestral, and the Magistral reasoning line.
> API platform formerly called "La Plateforme" — now branded **AI Studio** (`console.mistral.ai`).
> OpenAI-compatible endpoint. Many models are open-weight and self-hostable.

---

## Overview

Mistral AI is Europe's leading LLM company. They produce some of the best efficiency-per-cost models available, are strong in multilingual (especially French/European languages), and ship full open-weight releases alongside their PREMIER closed API models.

- **API Base URL:** `https://api.mistral.ai/v1`
- **SDK:** OpenAI-compatible — use OpenAI SDK with custom `baseURL`
- **Key Location:** 1Password → SO vault → `Mistral API Key`
- **Console:** https://console.mistral.ai
- **Docs:** https://docs.mistral.ai

---

## Models

### Generalist (Text + Vision)

| Model ID (API) | Version | Context | Tier | Best Use |
|----------------|---------|---------|------|---------|
| `mistral-large-latest` | v25.12 | 128K | OPEN | Flagship reasoning, multimodal |
| `mistral-medium-3-1-latest` | v25.08 | 128K | PREMIER | Frontier multimodal, long tasks |
| `mistral-small-latest` | v26.03 | 128K | OPEN | Fast + cheap, instruct+reasoning hybrid |
| `ministral-3-14b-latest` | v25.12 | 128K | OPEN | Vision + text, efficient edge model |
| `ministral-3-8b-latest` | v25.12 | 128K | OPEN | Efficient, best-in-class small |
| `mistral-nemo-latest` | v24.07 | 128K | OPEN | Excellent multilingual (12B) |

### Reasoning (Magistral Line)

| Model ID | Version | Context | Notes |
|----------|---------|---------|-------|
| `magistral-medium-latest` | v25.09 | 128K | Frontier reasoning, chain-of-thought |
| `magistral-small-latest` | v25.09 | 128K | Open-weight reasoning, efficient |

### Specialist

| Model ID | Type | Notes |
|----------|------|-------|
| `codestral-latest` | Code completion | Best-in-class FIM (fill-in-middle), 80K ctx |
| `devstral-2-latest` | Code agent | Solves software engineering tasks end-to-end |
| `voxtral-tts-latest` | TTS | Zero-shot voice cloning, multilingual |
| `voxtral-mini-transcribe-latest` | Transcription | Audio → text, efficient |
| `mistral-ocr-latest` | OCR / Document AI | Best-in-class document parsing |
| `mistral-embed-latest` | Embeddings | 1024 dims, 8K ctx — text semantic search |
| `codestral-embed-latest` | Code embeddings | Code semantic search, retrieval |

> **OPEN** = open-weight (downloadable, self-hostable)
> **PREMIER** = API-only (proprietary)

---

## Pricing Estimates

| Model | Input ($/1M) | Output ($/1M) | Notes |
|-------|-------------|--------------|-------|
| mistral-large | ~$2.00 | ~$6.00 | Check console for latest |
| mistral-medium | ~$0.40 | ~$2.00 | |
| mistral-small | ~$0.10 | ~$0.30 | Great value |
| ministral-8b | ~$0.10 | ~$0.10 | Cheapest capable model |
| magistral-medium | ~$2.00 | ~$5.00 | Reasoning premium |
| codestral | ~$0.30 | ~$0.90 | Code focus |
| mistral-embed | ~$0.10 | — | Per 1M tokens |

> Verify at https://mistral.ai/technology/#pricing — Mistral reprices frequently.

---

## API Setup

```bash
export MISTRAL_API_KEY="your-key"
```

### OpenAI SDK (recommended — simplest)

```typescript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.MISTRAL_API_KEY,
  baseURL: 'https://api.mistral.ai/v1',
});

const response = await client.chat.completions.create({
  model: 'mistral-small-latest',
  messages: [{ role: 'user', content: 'Hello' }],
});
```

### Native Mistral SDK

```bash
npm install @mistralai/mistralai
```

```typescript
import Mistral from '@mistralai/mistralai';

const client = new Mistral({ apiKey: process.env.MISTRAL_API_KEY });

const response = await client.chat.complete({
  model: 'mistral-large-latest',
  messages: [{ role: 'user', content: 'Hello' }],
});
```

### Codestral (fill-in-middle)

```typescript
const completion = await client.fim.complete({
  model: 'codestral-latest',
  prompt: 'def fibonacci(',
  suffix: '\n    return result',
});
```

### Embeddings

```typescript
const embed = await client.embeddings.create({
  model: 'mistral-embed',
  input: ['Text to embed'],
});
// embed.data[0].embedding — 1024 dimensions
```

---

## Open-Weight Self-Hosting (via Ollama)

Most Mistral models marked OPEN can be run locally:

```bash
ollama pull mistral
ollama pull mistral-small   # smaller variant
ollama pull codestral       # if available

# Run API-compatible server
ollama serve
# Use at: http://localhost:11434/v1 with model 'mistral'
```

---

## Usage in One Mind Stack

| Use Case | Model | Notes |
|----------|-------|-------|
| **Code completion / FIM** | codestral-latest | Best for in-editor autocomplete |
| **Code agents** | devstral-2-latest | Agentic software engineering tasks |
| **Multilingual tasks** | mistral-nemo-latest | Best 12B multilingual model |
| **Budget reasoning** | magistral-small-latest | Open-weight chain-of-thought |
| **Document parsing** | mistral-ocr-latest | PDF → structured text |
| **Audio transcription** | voxtral-mini-transcribe-latest | Voice → text pipeline |
| **Fast agent turns** | mistral-small-latest | 128K, cheap, OpenAI-compatible |
| **Code search RAG** | codestral-embed-latest | Embed codebase for semantic search |

---

## Notes

- **European data residency:** API servers in EU — relevant for GDPR-sensitive tasks
- **Multilingual strength:** Exceptional French, German, Italian, Spanish — best non-US choice for European languages
- **La Plateforme:** Old name for the API. Now "AI Studio" at `console.mistral.ai`. Some older docs still say La Plateforme.
- **Open vs PREMIER:** Large 3 and Small 4 are open-weight — you can download and run them yourself (Ollama, HuggingFace, self-hosted)
- **Function calling:** Supported on mistral-large, mistral-small, mistral-medium (JSON mode native)
- **Ministral 3B/8B/14B:** Their new "edge" line — vision-capable, efficient, run on-device
- **Voxtral TTS:** Zero-shot voice cloning — can clone a voice from audio sample at inference time


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
