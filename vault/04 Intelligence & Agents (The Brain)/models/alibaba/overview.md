---
id: so-01-ai-models-alibaba
title: "AI Models — Alibaba Cloud (Qwen)"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/alibaba]
ai_summary: "Alibaba Cloud Qwen model reference — context windows, pricing, API setup. In agent stack for long context (1M tokens) and multilingual tasks."
---

# Alibaba Cloud (Qwen)

> EXPERIMENTAL/STACK provider. Qwen is in the agent stack for ultra-long context tasks (up to 1M tokens) and multilingual work.
> Qwen-Long is notable: 1M context window at a fraction of the cost of Gemini 1.5 Pro.

---

## Overview

Alibaba's Qwen (通义千问) family offers competitive models with standout long-context capabilities. Qwen-Long provides a 1M token context window at very low cost — useful for processing entire books, large codebases, or extended multi-session conversation history.

- **API Base URL:** `https://dashscope.aliyuncs.com/compatible-mode/v1`
- **SDK:** OpenAI-compatible — use OpenAI SDK with custom `baseURL`
- **Key Location:** 1Password → SO vault → `Alibaba Cloud API Key`
- **Console:** https://dashscope.aliyun.com
- **Pricing:** https://help.aliyun.com/zh/dashscope/developer-reference/tongyi-thousand-tokens-price

---

## Models

| Model | Context | Input ($/1M) | Output ($/1M) | Best Use |
|-------|---------|-------------|--------------|---------|
| `qwen-max` | 32K | $2.40 | $9.60 | Best quality, flagship |
| `qwen-max-longcontext` | 32K | $2.40 | $9.60 | Stable long-context version of max |
| `qwen-plus` | 131K | $0.40 | $1.20 | Balanced quality/cost |
| `qwen-long` | **1M** | $0.05 | $0.20 | **Ultra-long context at lowest cost** |
| `qwen-turbo` | 131K | $0.05 | $0.20 | Fast, cheap, decent quality |
| `qwen2.5-72b-instruct` | 131K | $0.55 | $2.20 | Open-weight flagship, strong coding |
| `qwen2.5-coder-32b-instruct` | 128K | $0.28 | $0.84 | Code-specialized, strong on agent tasks |

> Prices are approximate USD equivalents (billed in CNY via Aliyun).

---

## API Setup

Qwen is OpenAI-compatible via DashScope:

```bash
export DASHSCOPE_API_KEY="sk-..."
```

```typescript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.DASHSCOPE_API_KEY,
  baseURL: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
});

const response = await client.chat.completions.create({
  model: 'qwen-long',
  messages: [{ role: 'user', content: 'Summarize this 500K-token document...' }],
});
```

**Alternative: DashScope native SDK**
```bash
pip install dashscope    # Python native SDK
npm install dashscope    # Not yet official — use OpenAI compat layer
```

---

## Qwen-Long: 1M Context Use Cases

| Scenario | Why Qwen-Long |
|----------|--------------|
| **Full vault synthesis** | Send entire codex domain as context ($0.05/1M vs $1.25 Gemini) |
| **Book summarization** | Process entire book in one call |
| **Multi-session history** | Long-running agent conversation memory |
| **Large codebase review** | Entire repo as context |
| **Research corpus** | Multiple papers + synthesis in one shot |

**Price comparison for 500K token context:**
- Qwen-Long: ~$0.025
- Gemini 1.5 Pro: ~$0.625
- Claude 3.7 Sonnet (200K limit, would need chunking)

---

## Usage in One Mind Stack

| Use Case | Model | Notes |
|----------|-------|-------|
| **Full vault context tasks** | qwen-long | Processing large sections of the codex |
| **Multilingual content** | qwen-max | Strong on Chinese, Japanese, Korean |
| **Code generation** | qwen2.5-coder-32b | Alternative to Claude for code-heavy tasks |
| **Cheap batch summaries** | qwen-turbo | High-volume, low-priority summaries |

---

## Notes

- **Latency:** DashScope servers are primarily in China — expect 200–500ms from US/EU vs 50–100ms domestically
- **Multimodal:** `qwen-vl-max` and `qwen-vl-plus` support image input
- **Search:** `qwen-max-internet` variant has Alibaba web search built in
- **Function calling:** Available on qwen-max, qwen-plus, qwen2.5 series
- **Rate limits:** Free tier is generous; pay-as-you-go after credit exhausted
- **Note on latency:** For real-time agent use, prefer Anthropic; use Qwen for async/batch jobs where latency is less critical


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
