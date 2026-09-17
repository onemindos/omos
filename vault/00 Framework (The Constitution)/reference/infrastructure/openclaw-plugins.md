# OpenClaw Plugin Reference

> **Version**: OpenClaw v2026.4.29 (a448042)
> **Instances**: omos-core (100.123.126.34) · omos-sentinel (100.66.235.38)
> **Last Audit**: 2026-04-30
> **Initial Audit**: 69/117 → 21/117 — **Current**: 44/117 enabled

---

## Why This Matters

Every enabled plugin loads at gateway startup — consumes memory, adds latency, and registers tool/model hooks even if never called. Trimming from 69 → 21 cuts cold-start time and attack surface. Only enable what you actually route traffic through.

---

## Category 1 — Model Providers (LLM Inference Backends)

These register model catalogs so agents can route to them. Only providers with active API keys / auth need to be on.

| # | Plugin ID | Purpose | State |
|---|-----------|---------|-------|
| 1 | `amazon-bedrock` | AWS Bedrock — primary inference (Claude Sonnet 4.6, Haiku 4.5, Opus) | ✅ ON |
| 2 | `github-copilot` | GitHub Copilot Models — secondary inference (Claude, GPT, Gemini, Grok via GH) | ✅ ON |
| 3 | `google` | Google AI — Gemini models, powers Talk voice (bidiGenerateContent) | ✅ ON |
| 4 | `openai` | OpenAI direct — GPT-4.1, GPT-5, o-series reasoning models | ✅ ON |
| 5 | `microsoft-foundry` | Azure AI Foundry — Microsoft-hosted models (GPT, Phi, etc.) | ✅ ON |
| 6 | `copilot-proxy` | Internal proxy layer for GitHub Copilot token refresh and routing | ✅ ON |
| 7 | `anthropic` | Anthropic direct API — redundant since we route Claude via Bedrock | ❌ OFF |
| 8 | `anthropic-vertex` | Anthropic on Google Vertex — not using Vertex for Claude | ❌ OFF |
| 9 | `openrouter` | OpenRouter aggregator — pay-per-token multi-provider proxy | ❌ OFF |
| 10 | `groq` | Groq — ultra-fast Llama/Mixtral inference on LPUs | ❌ OFF |
| 11 | `together` | Together AI — open-model hosting (Llama, Mistral, etc.) | ❌ OFF |
| 12 | `fireworks` | Fireworks AI — fast open-model inference | ❌ OFF |
| 13 | `deepinfra` | DeepInfra — GPU cloud for open models | ❌ OFF |
| 14 | `deepseek` | DeepSeek — Chinese LLM provider (DeepSeek-V3, Coder) | ❌ OFF |
| 15 | `cerebras` | Cerebras — wafer-scale inference for open models | ❌ OFF |
| 16 | `mistral` | Mistral AI direct — Mistral/Mixtral models | ❌ OFF |
| 17 | `nvidia` | NVIDIA NIM — NVIDIA-hosted model inference | ❌ OFF |
| 18 | `ollama` | Ollama — local model runner (requires local GPU) | ❌ OFF |
| 19 | `lmstudio` | LM Studio — local model GUI/server | ❌ OFF |
| 20 | `litellm` | LiteLLM — unified proxy for 100+ providers | ❌ OFF |
| 21 | `huggingface` | Hugging Face Inference API | ❌ OFF |
| 22 | `alibaba` | Alibaba Cloud Qwen/Tongyi models | ✅ ON |
| 23 | `amazon-bedrock-mantle` | Bedrock Mantle — secondary Bedrock routing layer | ❌ OFF |
| 24 | `arcee` | Arcee AI — fine-tuned enterprise models | ❌ OFF |
| 25 | `byteplus` | BytePlus (TikTok/ByteDance) model API | ❌ OFF |
| 26 | `chutes` | Chutes — serverless GPU inference | ❌ OFF |
| 27 | `cloudflare-ai-gateway` | Cloudflare AI Gateway — edge-cached model proxy | ✅ ON |
| 28 | `comfy` | ComfyUI — Stable Diffusion workflow engine (https://comfy.org/) | ❌ OFF |
| 29 | `deepgram` | Deepgram — speech-to-text / transcription API | ✅ ON |
| 30 | `fal` | fal.ai — fast image/video generation API | ✅ ON |
| 31 | `kilocode` | Kilocode — code-focused model provider | ❌ OFF |
| 32 | `kimi` | Kimi (Moonshot) — Chinese long-context LLM | ❌ OFF |
| 33 | `minimax` | MiniMax — Chinese multimodal LLM | ❌ OFF |
| 34 | `moonshot` | Moonshot AI — Kimi model API | ❌ OFF |
| 35 | `opencode` | OpenCode — code generation provider | ❌ OFF |
| 36 | `opencode-go` | OpenCode Go — Go-native code provider | ❌ OFF |
| 37 | `qianfan` | Baidu Qianfan — ERNIE and Chinese LLMs | ❌ OFF |
| 38 | `qwen` | Alibaba Qwen direct API | ✅ ON |
| 39 | `runway` | Runway ML — video generation models (Gen-3) | ❌ OFF |
| 40 | `vidra` | VidraAI — video generation provider (https://vidraai.com/) | ❌ OFF |
| 41 | `senseaudio` | SenseAudio — audio processing models | ❌ OFF |
| 42 | `sglang` | SGLang — structured generation runtime for LLMs | ❌ OFF |
| 43 | `stepfun` | StepFun — Chinese multi-step reasoning LLM | ❌ OFF |
| 44 | `synthetic` | Synthetic provider — test/mock models | ❌ OFF |
| 45 | `tencent` | Tencent Cloud Hunyuan models | ❌ OFF |
| 46 | `venice` | Venice AI — privacy-focused uncensored models | ❌ OFF |
| 47 | `vercel-ai-gateway` | Vercel AI Gateway — edge model routing | ❌ OFF |
| 48 | `vllm` | vLLM — self-hosted high-throughput inference server | ❌ OFF |
| 49 | `volcengine` | Volcengine (ByteDance) — Doubao models | ❌ OFF |
| 50 | `voyage` | Voyage AI — embedding models | ❌ OFF |
| 51 | `vydra` | Vydra — inference provider | ❌ OFF |
| 52 | `xai` | xAI — Grok models | ❌ OFF |
| 53 | `xiaomi` | Xiaomi MiLM models | ❌ OFF |
| 54 | `zai` | ZAI — inference provider | ❌ OFF |

---

## Category 2 — Speech / Text-to-Speech

Voice synthesis for Talk mode and audio output.

| # | Plugin ID | Purpose | State |
|---|-----------|---------|-------|
| 55 | `azure-speech` | Azure AI Speech — high-quality neural TTS (MP3, Opus, PCM) | ✅ ON |
| 56 | `elevenlabs` | ElevenLabs — ultra-realistic voice cloning and TTS | ❌ OFF |
| 57 | `microsoft` | Microsoft Speech Services — alternative Azure TTS path | ✅ ON |
| 58 | `inworld` | Inworld — gaming/NPC streaming TTS | ❌ OFF |
| 59 | `tts-local-cli` | Local CLI TTS — system `say` / `espeak` commands | ✅ ON |
| 60 | `gradium` | Gradium — experimental speech provider | ❌ OFF |

---

## Category 3 — Core Infrastructure

Essential platform services that make OpenClaw work.

| # | Plugin ID | Purpose | State |
|---|-----------|---------|-------|
| 60 | `acpx` | ACP Runtime — embedded Agent Communication Protocol session/transport management | ✅ ON |
| 61 | `memory-core` | Memory Core — base memory system (read/write/search agent memories) | ✅ ON |
| 62 | `device-pair` | Device Pairing — generates setup codes, approves node pairing requests | ✅ ON |
| 63 | `file-transfer` | File Transfer — fetch/list/write files on paired nodes via base64 (up to 16MB) | ✅ ON |
| 64 | `talk-voice` | Talk Voice — manage voice selection for Talk mode (list/set voices) | ✅ ON |
| 65 | `document-extract` | Document Extraction — extract text + page images from PDF/DOCX/etc attachments | ✅ ON |
| 66 | `bonjour` | Bonjour Discovery — advertise gateway over mDNS for auto-discovery on LAN | ✅ ON |
| 67 | `browser` | Browser — headless Chrome browsing for agents (navigate, read, screenshot) | ✅ ON |
| 68 | `web-readability` | Web Readability — extract clean article text from HTML fetches (Readability.js) | ✅ ON |
| 69 | `phone-control` | Phone Control — arm/disarm high-risk phone node commands (camera/screen/writes) | ✅ ON |

---

## Category 4 — Memory Upgrades

Enhanced memory capabilities beyond the base `memory-core`.

| # | Plugin ID | Purpose | State |
|---|-----------|---------|-------|
| 70 | `active-memory` | Active Memory — runs a memory sub-agent before replies, injects relevant memories into prompt context automatically | ✅ ON |
| 71 | `memory-lancedb` | Memory LanceDB — vector-indexed memory search using LanceDB embeddings | ✅ ON |
| 72 | `memory-wiki` | Memory Wiki — persistent wiki compiler, Obsidian-friendly knowledge vault | ✅ ON |

---

## Category 5 — Messaging / Chat Integrations

Connect OpenClaw to external chat platforms. All OFF — enable only when integrating with a specific platform.

| # | Plugin ID | Purpose | State |
|---|-----------|---------|-------|
| 73 | `bluebubbles` | BlueBubbles — iMessage bridge for non-Apple devices | ❌ OFF |
| 74 | `discord` | Discord — bot integration for Discord servers | ✅ ON |
| 75 | `feishu` | Feishu/Lark — ByteDance enterprise messaging | ❌ OFF |
| 76 | `googlechat` | Google Chat — Google Workspace messaging | ❌ OFF |
| 77 | `google-meet` | Google Meet — join calls via Chrome or Twilio | ❌ OFF |
| 78 | `imessage` | iMessage — native Apple Messages integration | ❌ OFF |
| 79 | `irc` | IRC — Internet Relay Chat protocol | ❌ OFF |
| 80 | `line` | LINE — Japanese/Asian messaging platform | ❌ OFF |
| 81 | `matrix` | Matrix — decentralized encrypted messaging | ✅ ON |
| 82 | `mattermost` | Mattermost — self-hosted Slack alternative | ❌ OFF |
| 83 | `msteams` | Microsoft Teams — M365 team messaging | ❌ OFF |
| 84 | `nextcloud-talk` | Nextcloud Talk — self-hosted video/chat | ❌ OFF |
| 85 | `nostr` | Nostr — decentralized social protocol | ❌ OFF |
| 86 | `qqbot` | QQ Bot — Tencent QQ messaging | ❌ OFF |
| 87 | `signal` | Signal — encrypted messaging | ❌ OFF |
| 88 | `slack` | Slack — workspace messaging | ❌ OFF |
| 89 | `synology-chat` | Synology Chat — NAS-based team chat | ❌ OFF |
| 90 | `telegram` | Telegram — cloud messaging with bot API | ✅ ON |
| 91 | `tlon` | Tlon/Urbit — decentralized computing platform messaging | ❌ OFF |
| 92 | `twitch` | Twitch — live streaming chat integration | ❌ OFF |
| 93 | `voice-call` | Voice Call — outbound telephony via Twilio/SIP | ✅ ON |
| 94 | `whatsapp` | WhatsApp — Meta messaging via Business API | ❌ OFF |
| 95 | `zalo` | Zalo — Vietnamese messaging platform | ❌ OFF |
| 96 | `zalouser` | Zalo User — Zalo user-level integration | ❌ OFF |

---

## Category 6 — Search / Web Tools

External search APIs. All OFF — the `browser` + `web-readability` combo handles web access.

| # | Plugin ID | Purpose | State |
|---|-----------|---------|-------|
| 97 | `brave` | Brave Search API — privacy-focused web search | ❌ OFF |
| 98 | `duckduckgo` | DuckDuckGo — anonymous web search | ✅ ON |
| 99 | `exa` | Exa — semantic/neural web search | ❌ OFF |
| 100 | `firecrawl` | Firecrawl — web scraping and crawling API | ✅ ON |
| 101 | `perplexity` | Perplexity — AI-powered search answers | ❌ OFF |
| 102 | `searxng` | SearXNG — self-hosted metasearch engine | ❌ OFF |
| 103 | `tavily` | Tavily — AI-optimized search API | ✅ ON |

---

## Category 7 — Tools & Utilities

Specialized capabilities — enable as needed for specific workflows.

| # | Plugin ID | Purpose | State |
|---|-----------|---------|-------|
| 104 | `codex` | Codex — OpenAI Codex app-server harness and managed GPT catalog | ✅ ON |
| 105 | `diffs` | Diffs — read-only diff viewer and file renderer for agents | ✅ ON |
| 106 | `llm-task` | LLM Task — generic JSON-only LLM tool for structured workflow tasks | ✅ ON |
| 107 | `lobster` | Lobster — typed workflow tool with resumable human approvals | ✅ ON |
| 108 | `open-prose` | OpenProse — VM skill pack with `/prose` slash command for writing | ❌ OFF |
| 109 | `skill-workshop` | Skill Workshop — capture repeatable workflows as reusable skills | ✅ ON |
| 110 | `thread-ownership` | Thread Ownership — prevents multi-agent conflicts in Slack threads | ✅ ON |
| 111 | `tokenjuice` | TokenJuice — compacts bash/exec tool output to save tokens | ✅ ON |
| 112 | `webhooks` | Webhooks — inbound webhook endpoints that trigger TaskFlows | ✅ ON |

---

## Category 8 — Migration (One-Time Use)

Import config from other AI assistants. Run once then disable.

| # | Plugin ID | Purpose | State |
|---|-----------|---------|-------|
| 113 | `migrate-claude` | Claude Migration — import Claude Code/Desktop config, MCP servers, skills | ❌ OFF |
| 114 | `migrate-hermes` | Hermes Migration — import Hermes config, memories, skills, credentials | ❌ OFF |

---

## Category 9 — Observability

Telemetry and monitoring exports. Enable when debugging performance.

| # | Plugin ID | Purpose | State |
|---|-----------|---------|-------|
| 115 | `diagnostics-otel` | OpenTelemetry — export traces/spans to Jaeger/Grafana Tempo/etc. | ✅ ON |
| 116 | `diagnostics-prometheus` | Prometheus — export metrics for Prometheus/Grafana scraping | ❌ OFF |

---

## Category 10 — Sandbox

Isolated execution environments for agent code.

| # | Plugin ID | Purpose | State |
|---|-----------|---------|-------|
| 117 | `openshell` | OpenShell Sandbox — SSH-based sandboxed command execution with workspace mirroring | ✅ ON |

---

## Quick Reference — What's ON (21 plugins)

```
amazon-bedrock, github-copilot, google, openai, microsoft-foundry, copilot-proxy,
alibaba, cloudflare-ai-gateway, deepgram, fal, qwen,
azure-speech, microsoft, tts-local-cli,
acpx, memory-core, device-pair, file-transfer, talk-voice, document-extract,
bonjour, browser, web-readability, phone-control,
active-memory, memory-lancedb, memory-wiki,
discord, matrix, telegram, voice-call,
duckduckgo, tavily, firecrawl,
codex, diffs, llm-task, lobster, skill-workshop,
thread-ownership, tokenjuice, webhooks,
diagnostics-otel, openshell
```

## Re-Enable Cheat Sheet

```bash
# Enable a provider you need temporarily
openclaw plugins enable <plugin-id>

# Disable when done
openclaw plugins disable <plugin-id>

# Always restart gateway after changes
openclaw gateway stop && systemctl --user start openclaw-gateway
```


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
