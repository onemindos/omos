---
id: so-01-ai-models-moonshot
title: "AI Models — Moonshot AI (Kimi)"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/moonshot]
ai_summary: "Moonshot AI Kimi model reference — context windows, pricing, API setup. In agent stack for extended context reasoning and kimi-k1.5 long chain-of-thought."
---

# Moonshot AI (Kimi)

> EXPERIMENTAL/STACK provider. Kimi is in the agent stack for extended context reasoning and kimi-k1.5's long chain-of-thought.
> kimi-k1.5: 128K context + 64K thinking tokens — strong for deep multi-step reasoning.

---

## Overview

Moonshot AI (月之暗面) is a Chinese AI lab whose Kimi models excel at long context with reasoning. The flagship `kimi-k1.5` combines a 128K context window with long chain-of-thought reasoning (extended thinking), making it a competitive alternative to Claude 3.7 Sonnet for reasoning-heavy tasks at lower cost.

- **API Base URL:** `https://api.moonshot.cn/v1`
- **SDK:** OpenAI-compatible — use OpenAI SDK with custom `baseURL`
- **Key Location:** 1Password → SO vault → `Moonshot API Key`
- **Console:** https://platform.moonshot.cn
- **Pricing:** https://platform.moonshot.cn/docs/pricing

---

## Models

| Model | Context | Input ($/1M) | Output ($/1M) | Best Use |
|-------|---------|-------------|--------------|---------|
| `kimi-k1.5` | 128K | $2.50 | $10 | **Flagship** — long CoT reasoning, complex tasks |
| `moonshot-v1-128k` | 128K | $2.60 | $10.40 | Stable long-context general model |
| `moonshot-v1-32k` | 32K | $0.85 | $2.50 | Balanced — most use cases |
| `moonshot-v1-8k` | 8K | $0.17 | $0.50 | Fast, cheapest, short tasks |

> Prices are USD equivalents — billed in CNY.

### kimi-k1.5 Extended Thinking

kimi-k1.5 supports long chain-of-thought reasoning:
- **Thinking tokens:** Up to 64K thinking tokens
- **Context:** 128K for input, additional for thinking
- **Activation:** Set `thinking: { budget: "auto" }` or specify token budget
- **Use case:** Comparable to Claude 3.7 Sonnet with extended thinking

---

## API Setup

OpenAI-compatible:

```bash
export MOONSHOT_API_KEY="sk-..."
```

```typescript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.MOONSHOT_API_KEY,
  baseURL: 'https://api.moonshot.cn/v1',
});

// Standard call
const response = await client.chat.completions.create({
  model: 'moonshot-v1-128k',
  messages: [{ role: 'user', content: 'Analyze this document...' }],
});

// With extended thinking (kimi-k1.5)
const deepResponse = await client.chat.completions.create({
  model: 'kimi-k1.5',
  messages: [{ role: 'user', content: 'Solve this complex problem...' }],
  // extended thinking via additional parameters (check current docs)
});
```

**File upload (Kimi supports direct file API):**
```typescript
// Kimi has a Files API compatible with OpenAI's
const file = await client.files.create({
  file: fs.createReadStream('document.pdf'),
  purpose: 'assistants',
});
// Reference file_id in message content
```

---

## Usage in One Mind Stack

| Use Case | Model | Notes |
|----------|-------|-------|
| **Deep reasoning tasks** | kimi-k1.5 | Alternative to Claude extended thinking |
| **Long document analysis** | moonshot-v1-128k | Process long PDFs or research papers |
| **Multilingual tasks** | kimi-k1.5 | Strong on Chinese + English mixed content |
| **Agent fallback** | moonshot-v1-32k | When Anthropic quota is hit |
| **Cheap short tasks** | moonshot-v1-8k | Quick lookups, simple Q&A |

---

## kimi-k1.5 vs Claude 3.7 Sonnet (Extended Thinking)

| Factor | kimi-k1.5 | claude-3.7-sonnet |
|--------|-----------|-------------------|
| Context | 128K | 200K |
| Thinking tokens | 64K | ~32K default |
| Reasoning quality | Strong | Excellent |
| Code | Strong | Excellent |
| Price (output) | ~$10/1M | $15/1M |
| Latency (from US) | ~300-600ms | ~150-300ms |
| Recommendation | Budget alt. | Primary choice |

---

## Notes

- **Latency:** Moonshot servers in China — 200–600ms from North America; use async/batch jobs
- **Real-time use:** For synchronous agent turns, prefer Anthropic; use Kimi for async operations
- **Vision:** `kimi-vl-a3b` and other vision models available — check current platform
- **Rate limits:** Free tier with 75 RPM limit on 8k; paid tier significantly higher
- **Multimodal:** Kimi can process images and documents via Files API
- **Best positioned for:** Overnight batch reasoning jobs, research synthesis, cost-sensitive long-context tasks


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
