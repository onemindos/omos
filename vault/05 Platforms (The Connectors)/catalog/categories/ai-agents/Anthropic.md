# Anthropic

> **AI safety company behind Claude — primary LLM provider**

## Overview

Anthropic is the AI research company that builds Claude. As the maker of Claude 3.5 Sonnet, Claude 3 Opus, and Claude Code, Anthropic is the primary AI provider for One Mind. Direct API access enables agent integration beyond the consumer Claude.ai interface.

## Status

| Attribute   | Value                           |
| ----------- | ------------------------------- |
| Status      | ✅ Active                        |
| Category    | AI / Provider                   |
| Type        | AI Research Company / API       |
| Integration | REST API, Python SDK, Claude.ai |

## Use Cases in OneMind

1. **Primary LLM** — Claude powers Claude Code, the primary agent IDE tool
2. **Agent Backbone** — OpenClaw agents use Claude-family models via Anthropic API
3. **Claude Code** — Direct AI coding assistant in VS Code
4. **API Integration** — Programmatic access to claude-3-5-sonnet, claude-3-opus

## Key Features

- Claude 3.5 Sonnet — best coding + reasoning model
- Claude 3 Opus — deep analysis, long documents
- 200K context window (Opus/Sonnet)
- Tool use / function calling
- Vision (image understanding)
- Constitutional AI safety training

## Quick Start

```python
import anthropic

client = anthropic.Anthropic(api_key="ANTHROPIC_API_KEY")
message = client.messages.create(
    model="claude-3-5-sonnet-20241022",
    max_tokens=1024,
    messages=[{"role": "user", "content": "Hello"}]
)
```

## Links

- Website: [anthropic.com](https://anthropic.com)
- API Docs: [docs.anthropic.com](https://docs.anthropic.com)
- Console: [console.anthropic.com](https://console.anthropic.com)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
