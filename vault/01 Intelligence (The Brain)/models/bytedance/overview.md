---
id: so-01-ai-models-bytedance
title: "AI Models — ByteDance (Doubao + Seedance)"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/bytedance]
ai_summary: "ByteDance AI model reference — Doubao LLM API (text/vision) and Seedance video generation. Chinese AI lab behind TikTok. Strong vision-language models and video gen."
---

# ByteDance (Doubao + Seedance)

> STACK provider. ByteDance runs two AI product lines:
> **Doubao** — LLM/VLM API (text, vision, reasoning) via VolcEngine ARK platform.
> **Seedance** — text-to-video generation model (best-in-class Chinese video AI, 2025–2026).

---

## Overview

ByteDance (parent of TikTok) has a serious AI research division. Their models are split between the **Doubao** family (conversational LLMs + vision-language) and the **Seed** research line (advanced VLMs). **Seedance** is their video generation flagship. All accessed via VolcEngine's ARK inference platform.

- **LLM API Base URL:** `https://ark.cn-beijing.volces.com/api/v3`
- **SDK:** OpenAI-compatible — use OpenAI SDK with custom `baseURL`
- **Key Location:** 1Password → SO vault → `ByteDance ARK API Key`
- **Console:** https://console.volcengine.com/ark
- **Pricing:** https://www.volcengine.com/product/ark/pricing

---

## Models

### Doubao LLM (Text + Vision)

| Model ID | Context | Input ($/1M) | Output ($/1M) | Best Use |
|----------|---------|-------------|--------------|---------|
| `doubao-pro-32k` | 32K | $0.11 | $0.33 | Flagship chat, reasoning |
| `doubao-pro-128k` | 128K | $0.57 | $1.71 | Long context tasks |
| `doubao-pro-256k` | 256K | $1.14 | $3.43 | Ultra-long context |
| `doubao-lite-32k` | 32K | $0.043 | $0.043 | Fast, cheap, high throughput |
| `doubao-lite-128k` | 128K | $0.11 | $0.11 | Budget long context |

### Seed Vision-Language (VLM)

| Model ID | Context | Notes |
|----------|---------|-------|
| `Seed-1.5-VL` | 128K | Flagship VLM — images + text, top benchmark scores |
| `Seed-1.5-Thinking` | 128K | Extended reasoning variant |

### Seedance (Video Generation)

| Model | Resolution | Notes |
|-------|-----------|-------|
| `seedance-1-0-lite` | Up to 1080p | Fast video gen, lower compute |
| `seedance-1-0-pro` | Up to 1080p | Flagship video quality |

> Seedance 1.0 released 2025 — competes with Sora, Kling, Wan2.1 for video generation quality.

---

## API Setup

Doubao uses the OpenAI-compatible format:

```bash
export ARK_API_KEY="your-ark-api-key"
# Note: ByteDance ARK also requires an endpoint ID (model deployment ID)
# Get endpoint IDs from console.volcengine.com/ark after creating a deployment
```

```typescript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.ARK_API_KEY,
  baseURL: 'https://ark.cn-beijing.volces.com/api/v3',
});

// LLM
const response = await client.chat.completions.create({
  model: 'doubao-pro-32k',  // or use endpoint ID from console
  messages: [{ role: 'user', content: 'Hello' }],
});

// Vision
const vision = await client.chat.completions.create({
  model: 'Seed-1.5-VL',
  messages: [{
    role: 'user',
    content: [
      { type: 'image_url', image_url: { url: 'https://...' } },
      { type: 'text', text: 'What is in this image?' }
    ]
  }],
});
```

**Video generation (Seedance — separate REST API):**
```python
import requests

response = requests.post(
    'https://ark.cn-beijing.volces.com/api/v3/contents/generations/tasks',
    headers={'Authorization': f'Bearer {ARK_API_KEY}'},
    json={
        'model': 'seedance-1-0-lite',
        'content': [{'type': 'text', 'text': 'A golden retriever running on a beach at sunset'}]
    }
)
# Returns a task_id — poll for completion
```

---

## Usage in One Mind Stack

| Use Case | Model | Notes |
|----------|-------|-------|
| **Vision analysis** | Seed-1.5-VL | Strong image understanding, competes with GPT-4o vision |
| **Long context (budget)** | doubao-lite-128k | Cheapest 128K context option |
| **Video generation** | seedance-1-0-pro | Brand content, GE/78 marketing, course demos |
| **High-volume text** | doubao-lite-32k | Bulk classification, tagging |
| **Reasoning tasks** | Seed-1.5-Thinking | Extended thinking alternative |

---

## Notes

- **Endpoint IDs:** Unlike most providers, ByteDance ARK may require creating "endpoint deployments" in the console and referencing their deployment IDs rather than model names — check current console docs
- **Latency:** Servers primarily in Beijing (cn-beijing) — 300–700ms from US/EU; use for async/batch
- **International access:** VolcEngine has global regions but model availability varies; `cn-beijing` has most models
- **Seedance vs Kling:** Seedance (ByteDance) and Kling (Kuaishou) are the two top Chinese video models competing with Sora/Runway; Seedance edges Kling on speed, Kling edges on some quality metrics
- **Rate limits:** Check ARK console dashboard — free credits available on signup


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
