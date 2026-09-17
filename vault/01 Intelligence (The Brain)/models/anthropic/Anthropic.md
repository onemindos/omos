---
id: so-01-ai-models-anthropic
title: "AI Models — Anthropic"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/anthropic]
ai_summary: "Anthropic Claude model reference — context windows, pricing, API setup, and usage in One Mind agent stack. PRIMARY provider for most agents."
---

# Anthropic

> PRIMARY model provider. Claude powers Legacy, Oracle, Empire, and Trinity agents.
> Claude = the backbone of One Mind's intelligence layer.

---

## Overview

Anthropic is the primary AI provider for the One Mind stack. Claude models have best-in-class reasoning, strong code generation, long context support, and extended thinking mode for complex tasks.

- **API Base URL:** `https://api.anthropic.com`
- **SDK:** `anthropic` (Node.js / Python)
- **Key Location:** 1Password → SO vault → `Anthropic API Key`
- **Console:** https://console.anthropic.com
- **Pricing page:** https://www.anthropic.com/pricing

---

## Models

| Model | Context | Input ($/1M) | Output ($/1M) | Best Use |
|-------|---------|-------------|--------------|---------|
| `claude-opus-4` | 200K | ~$15 | ~$75 | Most complex tasks, maximum capability |
| `claude-3-7-sonnet-20250219` | 200K | $3 | $15 | Agent reasoning, extended thinking |
| `claude-3-5-sonnet-20241022` | 200K | $3 | $15 | Daily agent tasks, balanced speed/quality |
| `claude-3-5-haiku-20241022` | 200K | $0.80 | $4 | Fast tool calls, lightweight agents |
| `claude-3-opus-20240229` | 200K | $15 | $75 | Archival — use claude-opus-4 instead |

> ⚠️ Always use versioned model IDs in code (not aliases like `claude-sonnet`) to avoid accidental upgrades.

### Extended Thinking (claude-3-7-sonnet)

Activate via:
```json
{
  "thinking": {
    "type": "enabled",
    "budget_tokens": 10000
  }
}
```

Best for: multi-step reasoning, system design, research synthesis, complex decisions.
Add cost: thinking tokens billed at output rate.

---

## API Setup

```bash
# Install SDK
npm install @anthropic-ai/sdk          # Node.js
pip install anthropic                  # Python

# Environment variable
export ANTHROPIC_API_KEY="sk-ant-..."  # Never commit to code
```

**Minimal call example (Node.js):**
```typescript
import Anthropic from '@anthropic-ai/sdk';
const client = new Anthropic();
const msg = await client.messages.create({
  model: 'claude-3-7-sonnet-20250219',
  max_tokens: 1024,
  messages: [{ role: 'user', content: 'Hello' }],
});
```

**Rate limits (default tier):**
- claude-3-7-sonnet: 1K requests/min, 80K input tokens/min
- claude-3-5-haiku: 4K requests/min, 400K input tokens/min

---

## Usage in One Mind Stack

| Agent | Model | Why |
|-------|-------|-----|
| **Legacy** (orchestrator) | claude-3-7-sonnet | Extended thinking for complex planning |
| **Oracle** (HP mind) | claude-3-7-sonnet | Deep personal analysis, reflection |
| **Empire** (GE) | claude-3-7-sonnet | Business strategy, revenue reasoning |
| **Trinity** (SO meta) | claude-3-7-sonnet | Agent coordination decisions |
| **Spartan** (HP body) | claude-3-5-haiku | Fast fitness/habit responses |
| **Haven** (LE) | claude-3-5-haiku | Home/family — simple, quick |
| **Claude Code** (dev) | claude-opus-4 | Full codebase context, max capability |

---

## Notes

- **Token counting:** Use `client.messages.countTokens()` before large payloads
- **Streaming:** Prefer streaming for agent UX — `client.messages.stream()`
- **Images:** Pass as base64 or URL in content blocks; max 5MB per image
- **System prompt:** Keep under 1K tokens for haiku, up to 10K for sonnet/opus
- **Prompt caching:** Available for system prompts >1024 tokens — use for repeated agent context (60%+ cost reduction)
- **Batch API:** For offline/bulk jobs use Message Batches API (50% cost reduction)


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Models (AI Models)/anthropic (Anthropic)/anthropic|anthropic]]
