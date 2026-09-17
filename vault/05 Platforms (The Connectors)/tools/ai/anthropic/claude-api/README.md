---
type: platform
domain: SO-05
status: active
updated: 2026-05-14
owner: zeus
tags: [platform/anthropic, product/claude-api, domain/so]
ai_summary: "Anthropic API — direct model access, rate limits, pricing tiers, API keys, and available models."
---

# Anthropic API

Direct API access to Claude models. Used by OpenClaw agents (via Amazon Bedrock fallback), custom scripts, and any integration that calls Claude programmatically.

## Access

- Console: `console.anthropic.com`
- API key: 1Password → SO vault → "Anthropic API Key"
- Base URL: `https://api.anthropic.com/v1`

## Models Available (May 2026)

| Model | ID | Best For |
| ----- | -- | -------- |
| Claude Sonnet 4.6 | `claude-sonnet-4-6-20251101` | Primary — best balance of capability + cost |
| Claude Haiku 4.5 | `claude-haiku-4-5-20251001` | Fast, cheap — cron jobs, background tasks |
| Claude Opus 4 | `claude-opus-4-20251001` | Most capable — complex reasoning |

## Rate Limits (Tier 1 / Tier 2)

| Metric | Tier 1 | Tier 2 |
| ------ | ------ | ------ |
| RPM | 50 | 1,000 |
| TPM | 40K | 160K |
| TPD | 1M | 5M |

*Check current tier at `console.anthropic.com` → Limits*

## How One Mind Uses the API

| Usage | Via | Model |
| ----- | --- | ----- |
| OpenClaw primary | GitHub Copilot provider | Sonnet 4.6 |
| OpenClaw fallback | Amazon Bedrock | Sonnet 4.6 |
| OpenClaw cron | GitHub Copilot provider | Haiku 4.5 |
| Direct scripts | Anthropic API key | varies |

## Cost Tracking

- Console: `console.anthropic.com` → Usage
- Vault: `05 Platforms/tools/ai/anthropic/platform/`
