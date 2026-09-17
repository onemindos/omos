---
id: so-01-ai-models-xai
title: "AI Models — xAI (Grok)"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/xai]
ai_summary: "xAI Grok model reference — context windows, pricing, API setup. Experimental use for real-time web search augmented reasoning and X/Twitter data tasks."
---

# xAI (Grok)

> EXPERIMENTAL provider. Grok used for real-time web-augmented reasoning and tasks requiring live X/Twitter data.
> OpenAI-compatible API — easy drop-in via `apiEndpoint` override.

---

## Overview

xAI's Grok models are strong reasoning models with real-time internet access built in. The API is fully OpenAI-compatible (same endpoint shape), making it a zero-friction alternative when current information or X platform data is relevant.

- **API Base URL:** `https://api.x.ai/v1`
- **SDK:** Use OpenAI SDK with custom `baseURL` (OpenAI-compatible)
- **Key Location:** 1Password → SO vault → `xAI API Key`
- **Console:** https://console.x.ai
- **Pricing:** https://x.ai/api/pricing

---

## Models

| Model | Context | Input ($/1M) | Output ($/1M) | Best Use |
|-------|---------|-------------|--------------|---------|
| `grok-3` | 131K | $3 | $15 | Flagship — deep reasoning + live web |
| `grok-3-mini` | 131K | $0.30 | $0.50 | Fast reasoning on budget |
| `grok-3-fast` | 131K | $5 | $25 | Low-latency grok-3 |
| `grok-3-mini-fast` | 131K | $0.60 | $4 | Lowest latency, structured tasks |
| `grok-2-1212` | 131K | $2 | $10 | Archival — grok-3 preferred |

> Grok-3 includes live web search by default when enabled via the `search_enabled` parameter.

---

## API Setup

xAI uses the OpenAI SDK with a custom base URL:

```typescript
import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: process.env.XAI_API_KEY,
  baseURL: 'https://api.x.ai/v1',
});

const response = await client.chat.completions.create({
  model: 'grok-3',
  messages: [{ role: 'user', content: 'What happened in AI this week?' }],
  // Live search (when supported):
  // search_enabled: true,
});
```

**Environment variable:**
```bash
export XAI_API_KEY="xai-..."
```

---

## In OneMinID Web (OpenClaw)

Since xAI is OpenAI-compatible, it integrates via the standard provider pattern:

```typescript
// Patch apiEndpoint in session, then delegate to streamOpenAiChat
const patchedSession = {
  ...session,
  apiEndpoint: 'https://api.x.ai/v1',
  apiKey: process.env.XAI_API_KEY,
};
return streamOpenAiChat({ ...opts, session: patchedSession });
```

---

## Usage in One Mind Stack

| Use Case | Model | Notes |
|----------|-------|-------|
| **Real-time news/events** | grok-3 | "What happened in AI today?" |
| **X/Twitter analysis** | grok-3 | Has privileged access to X platform data |
| **Market/trend research** | grok-3-mini | Fast competitive research |
| **Current events briefing** | grok-3 | Morning briefing current context |

---

## Notes

- **Live search:** Grok-3 can search the web in real time — better than Perplexity for agent use
- **X data access:** Grok has unique access to real-time X/Twitter corpus
- **Rate limits:** Check console.x.ai — limits vary by tier
- **No vision yet (as of Q1 2026):** Grok doesn't process images in the standard API
- **Use case clarity:** Default to Anthropic for reasoning; switch to Grok when recency matters


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Models (AI Models)/xai (xAI)/xai|xai]]
