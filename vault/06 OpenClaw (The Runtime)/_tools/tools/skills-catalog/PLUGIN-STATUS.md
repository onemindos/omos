---
title: "OpenClaw Plugin Status"
type: reference
status: active
created: 2026-04-15
updated: 2026-04-15
owner: zeus
tags: [type/reference, domain/SO, domain/01, status/active, system/openclaw]
ai_summary: "Live plugin status snapshot for OpenClaw on onemindos-cloud-1. 36/98 loaded. Updated after each plugin audit session."
---

# OpenClaw Plugin Status

> **VPS**: `onemindos-cloud-1` (`64.227.8.23`)
> **OpenClaw version**: 2026.4.14
> **Last snapshot**: 2026-04-15 (session 2)
> **Status**: 36 / 98 loaded

Update this file after any plugin enable/disable session. Use `openclaw plugins list` on VPS to get current state.

---

## Summary

| Category | Loaded | Disabled | Notes |
| --- | --- | --- | --- |
| Core / System | 15 | 0 | All essential runtime plugins active |
| AI Providers — Western | 6 | 16 | Trimmed to: amazon-bedrock, copilot-proxy, github-copilot, google, huggingface, microsoft-foundry |
| AI Providers — Local | 6 | 0 | ollama, lmstudio, vllm, sglang, opencode, opencode-go |
| AI Providers — China/Regional | 0 | 13 | All disabled — not in use |
| Media / Creative | 4 | 2 | comfy + elevenlabs disabled |
| Search / Web | 3 | 3 | tavily + duckduckgo + browser active |
| Channels | 0 | 13 | telegram + discord disabled (malformed manifest / deferred) |
| **TOTAL** | **34** | **64** | — |

---

## LOADED Plugins (36)

### Core / System + Utility (15)

| # | ID | Name | Notes |
| --- | --- | --- | --- |
| 1 | `acpx` | ACPX Runtime | ACP runtime backend, session/transport management |
| 2 | `active-memory` | Active Memory | Memory sub-agent injected before replies |
| 3 | `codex` | Codex | Codex app-server harness + GPT model catalog |
| 4 | `device-pair` | Device Pairing | Mac ↔ VPS pairing (Mac device `5a4ff116` paired) |
| 5 | `diffs` | Diffs | Read-only diff viewer for agents |
| 6 | `llm-task` | LLM Task | Generic JSON-only LLM tool for workflows |
| 7 | `lobster` | Lobster | Typed workflow tool with resumable approvals |
| 8 | `memory-core` | Memory Core | Core memory search |
| 9 | `memory-wiki` | Memory Wiki | Persistent wiki / Obsidian-friendly knowledge vault |
| 10 | `openshell` | OpenShell Sandbox | SSH-based command execution with mirrored workspaces |
| 11 | `phone-control` | Phone Control | Arm/disarm high-risk phone node commands |

| 12 | `webhooks` | Webhooks | Inbound webhooks → TaskFlows |
| 13 | `open-prose` | OpenProse | Writing mode + `/prose` slash command |
| 14 | `diagnostics-otel` | OTEL Diagnostics | ✅ Verified working — traces to Datadog APM via OTLP gRPC:4317 |
| 15 | `talk-voice` | Talk Voice | Voice selection management |

### AI Providers — Western (6) — trimmed 2026-04-15

| # | ID | Name | Notes |
| --- | --- | --- | --- |
| 16 | `amazon-bedrock` | Amazon Bedrock | AWS creds in env |
| 17 | `copilot-proxy` | Copilot Proxy | GitHub Copilot provider |
| 18 | `github-copilot` | GitHub Copilot | GitHub Copilot |
| 19 | `google` | Google | Gemini models |
| 20 | `huggingface` | Hugging Face | HF models |
| 21 | `microsoft-foundry` | Microsoft Foundry | Azure AI Foundry |

### AI Providers — Local (6)

| # | ID | Name | Notes |
| --- | --- | --- | --- |
| 22 | `ollama` | Ollama | Primary local inference (qwen2.5:7b, qwen3:4b, etc.) |
| 23 | `lmstudio` | LM Studio | LM Studio local server |
| 24 | `vllm` | vLLM | vLLM inference server |
| 25 | `sglang` | SGLang | SGLang inference |
| 26 | `opencode` | OpenCode | OpenCode Zen provider |
| 27 | `opencode-go` | OpenCode Go | OpenCode Go provider |

### Media / Creative (4)

| # | ID | Name | Notes |
| --- | --- | --- | --- |
| 28 | `deepgram` | Deepgram | Audio transcription / STT |
| 29 | `fal` | fal.ai | Image/video generation |
| 30 | `microsoft` | Microsoft Speech | TTS via Microsoft |
| 31 | `runway` | Runway | Video generation |

### Search / Web (3)

| # | ID | Name | Notes |
| --- | --- | --- | --- |
| 32 | `browser` | Browser | Web browsing tool |
| 33 | `duckduckgo` | DuckDuckGo | Free search, no key needed |
| 34 | `tavily` | Tavily | Real-time structured search, key in env |

### Channels (2)

| # | ID | Name | Notes |
| --- | --- | --- | --- |
| 35 | `discord` | Discord | **Disabled** — malformed manifest blocks CLI; re-enable when token ready |
| 36 | `telegram` | Telegram | **Disabled** — malformed manifest blocks CLI; re-enable when token ready |

---

## DISABLED Plugins (43)

### China / Regional Providers (13) — Disabled 2026-04-15

| ID | Name | Reason |
| --- | --- | --- |
| `alibaba` | Alibaba Model Studio | China provider, not in use |
| `byteplus` | BytePlus | China provider (ByteDance), not in use |
| `chutes` | Chutes.ai | Not in use |
| `kimi` | Kimi (Moonshot AI) | China provider, not in use |
| `minimax` | MiniMax | China provider, not in use |
| `moonshot` | Moonshot | China provider, not in use |
| `qianfan` | Qianfan (Baidu) | China provider, not in use |
| `qwen` | Qwen Cloud (Alibaba) | China provider — Ollama handles local Qwen models instead |
| `stepfun` | StepFun | China provider, not in use |
| `volcengine` | Volcengine (ByteDance) | China provider, not in use |
| `vydra` | Vydra | Not in use |
| `xiaomi` | Xiaomi | China provider, not in use |
| `zai` | Z.AI | Not in use |

### Routing / Utility — Disabled 2026-04-15

| ID | Name | Reason |
| --- | --- | --- |
| `amazon-bedrock-mantle` | Bedrock Mantle | Cross-region routing layer — unnecessary complexity |
| `thread-ownership` | Thread Ownership | Slack-specific, not needed without Slack |

### Cloud Providers — Disabled 2026-04-15 (session 2)

| ID | Name | Reason |
| --- | --- | --- |
| `anthropic` | Anthropic | No API key on VPS — use locally via Mac if needed |
| `anthropic-vertex` | Anthropic Vertex | No Vertex setup |
| `arcee` | Arcee | Not in use |
| `cloudflare-ai-gateway` | Cloudflare AI Gateway | Not configured |
| `comfy` | ComfyUI | No ComfyUI server running on VPS |
| `deepseek` | DeepSeek | No API key |
| `fireworks` | Fireworks | No API key |
| `groq` | Groq | No API key |
| `kilocode` | Kilo Gateway | Not in use |
| `litellm` | LiteLLM | No LiteLLM proxy running |
| `mistral` | Mistral | No API key |
| `nvidia` | NVIDIA NIM | No key / not in use |
| `openai` | OpenAI | No API key on VPS |
| `openrouter` | OpenRouter | No API key |
| `synthetic` | Synthetic | Mock provider, not needed |
| `together` | Together | No API key |
| `venice` | Venice | No API key |
| `vercel-ai-gateway` | Vercel AI Gateway | Not configured |
| `xai` | xAI | No API key |

### Channels — Pending Setup

| ID | Name | What's Needed |
| --- | --- | --- |
| `bluebubbles` | BlueBubbles | iMessage via BlueBubbles server |
| `feishu` | Feishu/Lark | Bot token |
| `googlechat` | Google Chat | Bot setup |
| `imessage` | iMessage | Mac-local, VPS can't use |
| `irc` | IRC | IRC server config |
| `line` | LINE | Bot token |
| `matrix` | Matrix | Homeserver + credentials |
| `mattermost` | Mattermost | Server URL + token |
| `msteams` | Microsoft Teams | Bot registration |
| `nextcloud-talk` | Nextcloud Talk | Nextcloud server |
| `nostr` | Nostr | Keypair |
| `qqbot` | QQ Bot | China messaging |
| `signal` | Signal | Signal CLI setup |
| `slack` | Slack | Bot token |
| `synology-chat` | Synology Chat | Synology NAS |
| `tlon` | Tlon/Urbit | Urbit ship |
| `twitch` | Twitch | Bot token |
| `whatsapp` | WhatsApp | Complex — requires WA Business API |
| `voice-call` | Voice Call | ⚠️ Keep disabled for now — reserved |
| `zalo` | Zalo | Vietnam messaging |
| `zalouser` | Zalo Personal | Vietnam messaging |

### Providers — Disabled (Need Key / No Use Case)

| ID | Name | Reason |
| --- | --- | --- |
| `brave` | Brave Search | Needs API key — duckduckgo/tavily cover this |
| `elevenlabs` | ElevenLabs | TTS — needs API key, no current use case |
| `exa` | Exa | Needs API key — tavily covers this |
| `firecrawl` | Firecrawl | Web scraping — needs API key |
| `memory-lancedb` | Memory LanceDB | Local vector DB — redundant with memory-core |
| `perplexity` | Perplexity | Needs API key |
| `searxng` | SearXNG | Needs self-hosted SearXNG instance |

---

## Pending Actions

| Action | Plugin | Priority | Notes |
| --- | --- | --- | --- |
| Configure bot token | `telegram` | HIGH | Get from BotFather → set in OpenClaw config |
| Configure bot token | `discord` | HIGH | Get from Discord Developer Portal → set in config |
| Evaluate | `elevenlabs` | MED | Good TTS — worth enabling if voice workflow needed |
| Evaluate | `firecrawl` | MED | Better scraping than browser for structured content |
| Evaluate | `exa` | LOW | AI-powered search — tavily may be sufficient |
| Evaluate | `brave` | LOW | Tavily + duckduckgo sufficient for now |

---

## CLI Reference

```bash
# List all plugins
openclaw plugins list

# Enable a plugin
openclaw plugins enable <id>

# Disable a plugin
openclaw plugins disable <id>

# Restart gateway after changes
systemctl --user restart openclaw-gateway

# Check gateway status
systemctl --user status openclaw-gateway --no-pager
```

---

## Change Log

| Date | Change | By |
| --- | --- | --- |
| 2026-04-15 | Disabled 13 China/regional providers | Zeus |
| 2026-04-15 | Disabled amazon-bedrock-mantle, thread-ownership | Zeus |
| 2026-04-15 | Enabled tavily (+ API key in env) | Zeus |
| 2026-04-15 | Enabled duckduckgo, discord, telegram | Zeus |
| 2026-04-15 | Enabled open-prose, diagnostics-otel | Zeus |
| 2026-04-15 | OTEL→Datadog APM pipeline wired (see OPENCLAW-OTEL-DATADOG.md) | Zeus |
| 2026-04-15 | Disabled 19 cloud providers (anthropic, deepseek, groq, openai, etc.) — trimmed to local+essential only | Zeus |
| 2026-04-15 | OTEL verified: port 4317 open, DD Agent OTLP Status=Enabled+Running, env vars confirmed | Zeus |
| 2026-04-15 | Disabled telegram + discord — malformed plugin manifests blocked all CLI commands | Zeus |


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
