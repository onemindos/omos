---
id: so-01-ai-models-openai
title: "AI Models — OpenAI"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/openai]
ai_summary: "OpenAI GPT/o-series model reference — context windows, pricing, API setup, and usage in One Mind stack. SECONDARY provider; primary for embeddings and Guardian agent."
---

# OpenAI

> SECONDARY model provider. GPT/o-series used for Guardian agent, embeddings, and broad-compatibility fallback.
> Also the source of `text-embedding-3-small` — the embedding model for pgvector RAG.

---

## Overview

OpenAI is the secondary provider in the One Mind stack. GPT-4o handles multimodal tasks and the Guardian security agent. The o-series models are reserved for deep reasoning when Anthropic extended thinking isn't available or a second opinion is needed.

- **API Base URL:** `https://api.openai.com`
- **SDK:** `openai` (Node.js / Python)
- **Key Location:** 1Password → SO vault → `OpenAI API Key`
- **Platform:** https://platform.openai.com
- **Pricing page:** https://openai.com/pricing

---

## Models

### Chat Models

| Model | Context | Input ($/1M) | Output ($/1M) | Best Use |
|-------|---------|-------------|--------------|---------|
| `gpt-4o` | 128K | $2.50 | $10 | Multimodal, vision, broad tasks |
| `gpt-4o-mini` | 128K | $0.15 | $0.60 | Cheap fallback, structured outputs |
| `o3` | 200K | $10 | $40 | Deep reasoning, complex problems |
| `o3-mini` | 200K | $1.10 | $4.40 | Reasoning on budget |
| `o1` | 200K | $15 | $60 | Archival — use o3 instead |

### Embedding Models

| Model | Dimensions | Price ($/1M tokens) | Best Use |
|-------|-----------|--------------------|---------|
| `text-embedding-3-small` | 1536 | $0.02 | **PRIMARY** — pgvector RAG, default |
| `text-embedding-3-large` | 3072 | $0.13 | Higher accuracy when needed |
| `text-embedding-ada-002` | 1536 | $0.10 | **Legacy** — avoid for new work |

> Use `text-embedding-3-small` by default. It outperforms ada-002 at 1/5 the price.

### Image & Audio

| Model | Use |
|-------|-----|
| `dall-e-3` | Image generation (brand assets, diagrams) |
| `whisper-1` | Audio transcription ($0.006/min) |
| `tts-1` | Text to speech (1M chars = $15) |

---

## API Setup

```bash
npm install openai        # Node.js
pip install openai        # Python

export OPENAI_API_KEY="sk-..."
```

**Chat example (Node.js):**
```typescript
import OpenAI from 'openai';
const client = new OpenAI();
const response = await client.chat.completions.create({
  model: 'gpt-4o',
  messages: [{ role: 'user', content: 'Hello' }],
});
```

**Embedding example:**
```typescript
const embedding = await client.embeddings.create({
  model: 'text-embedding-3-small',
  input: 'Your text here',
});
// embedding.data[0].embedding → float32[1536]
```

---

## Usage in One Mind Stack

| Use Case | Model | Notes |
|----------|-------|-------|
| **Guardian agent** | gpt-4o | Security analysis, threat evaluation |
| **RAG embeddings** | text-embedding-3-small | pgvector on VPS — all vault indexing |
| **Reasoning fallback** | o3-mini | When Anthropic quota hit or 2nd opinion needed |
| **Vision/images** | gpt-4o | Analyzing screenshots, diagrams, photos |
| **Audio transcription** | whisper-1 | Voice memos → text for capture pipeline |

---

## VPS Embedding Pipeline

The pgvector knowledge index on the VPS uses OpenAI embeddings:

```
Vault markdown → chunk → text-embedding-3-small → vector(1536) → pgvector
```

Table: `knowledge_chunks` with columns `content TEXT`, `embedding vector(1536)`, `source TEXT`, `domain TEXT`.

---

## Notes

- **Structured outputs:** Use `response_format: { type: 'json_schema' }` for guaranteed JSON from gpt-4o
- **Vision:** Pass images as `{ type: 'image_url', image_url: { url: '...' } }` in user message
- **o3 reasoning:** Does not support `temperature` or streaming in same way — check docs
- **Rate limits:** Tier 2 = 10K RPM for gpt-4o-mini, 5K for gpt-4o
- **Batch API:** 50% discount for async jobs (embeddings especially)


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Models (AI Models)/openai (OpenAI)/openai|openai]]
