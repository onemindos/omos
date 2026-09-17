---
id: so-01-ai-models-minimax
title: "AI Models — MiniMax"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/minimax]
ai_summary: "MiniMax AI model reference — MiniMax-Text-01 (4M context), abab series. Chinese AI lab. In agent stack for ultra-long context tasks."
---

# MiniMax

> STACK provider. MiniMax-Text-01 has a **4 million token context window** — the largest available context of any model in the stack.
> Independent Chinese AI lab. Not affiliated with Alibaba, Moonshot, or other Chinese providers.

---

## Overview

MiniMax is a Chinese AI company whose `MiniMax-Text-01` model holds one of the longest context windows available (4M tokens). This makes it uniquely suited for full-vault synthesis, processing entire book collections, or any task requiring extreme context depth.

- **API Base URL:** `https://api.minimax.chat/v1`
- **SDK:** OpenAI-compatible — use OpenAI SDK with custom `baseURL`
- **Key Location:** 1Password → SO vault → `MiniMax API Key`
- **Console:** https://www.minimaxi.com/platform
- **Docs:** https://www.minimaxi.com/document/guides/chat-model/pro

---

## Models

| Model | Context | Input ($/1M) | Output ($/1M) | Best Use |
|-------|---------|-------------|--------------|---------|
| `MiniMax-Text-01` | **4M** | $0.20 | $1.10 | Ultra-long context — flagship |
| `abab6.5s-chat` | 245K | $0.10 | $0.10 | Fast, balanced quality |
| `abab6.5g-chat` | 8K | $0.005 | $0.005 | Cheapest, short tasks |
| `abab5.5-chat` | 16K | $0.015 | $0.015 | Legacy model |

> MiniMax-Text-01's 4M context = ~3M words = ~6,000 pages. Effectively the entire codex vault in one context window.

### Context Window Comparison (Stack)

| Provider | Model | Max Context |
|----------|-------|------------|
| **MiniMax** | MiniMax-Text-01 | **4,000,000** |
| Alibaba | qwen-long | 1,000,000 |
| Moonshot | kimi-k1.5 | 128,000 |
| Google | gemini-1.5-pro | 2,000,000 |
| Anthropic | claude-3.7-sonnet | 200,000 |
| OpenAI | gpt-4o | 128,000 |

---

## API Setup

OpenAI-compatible:

```bash
export MINIMAX_API_KEY="eyJ..."
export MINIMAX_GROUP_ID="..."  # Required for some endpoints
```

```typescript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.MINIMAX_API_KEY,
  baseURL: 'https://api.minimax.chat/v1',
  defaultHeaders: {
    'Authorization': `Bearer ${process.env.MINIMAX_API_KEY}`,
  },
});

const response = await client.chat.completions.create({
  model: 'MiniMax-Text-01',
  messages: [{ role: 'user', content: 'Analyze this 2M token document...' }],
});
```

> Note: Some MiniMax endpoints require a `GroupID` header or query param — check docs for the current requirement.

---

## Usage in One Mind Stack

| Use Case | Model | Notes |
|----------|-------|-------|
| **Full vault synthesis** | MiniMax-Text-01 | Entire codex in one context window |
| **Book collection analysis** | MiniMax-Text-01 | Multiple books at once |
| **Long agent history** | MiniMax-Text-01 | Extended multi-session conversation |
| **Cheap bulk tasks** | abab6.5g-chat | Cheapest model when quality isn't critical |

---

## Notes

- **4M context reality:** Effective retrieval degrades in very long contexts — MiniMax performs better than most but still has degradation beyond ~500K tokens for precise recall
- **Latency:** Servers in China — 300–700ms from US/EU; prefer async/batch use
- **Multimodal:** MiniMax has image/video generation (Hailuo) but text-only for chat API
- **Rate limits:** Check platform dashboard — free tier is limited
- **Authentication:** Some API paths need `GroupID` in addition to API key — store both in 1Password item
- **Best positioned for:** Async synthesis jobs where 4M context is the deciding factor over Gemini 1.5 Pro


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
