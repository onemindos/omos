---
title: "OpenClaw — Provider & Plugin Stack"
type: reference
status: active
created: 2026-04-14
updated: 2026-05-02
owner: zeus
tags: [type/reference, domain/SO, domain/openclaw, status/active]
ai_summary: "Current LLM, STT, TTS, and media provider configuration on both VPS nodes (omos-core, omos-sentinel) — plugins, credentials, models, drop-ins, and voice stack."
---

# OpenClaw — Provider & Plugin Stack

> Last verified: **2026-05-02 08:29 UTC** — full smoke test on omos-core
> OpenClaw version: **v2026.4.29** (a448042) on both boxes

---

## Nodes Overview

| Node | Tailscale IP | Role | Service Type |
|------|-------------|------|-------------|
| **omos-core** | `100.123.126.34` | Primary | System service (`sudo systemctl`) |
| **omos-sentinel** | `100.66.235.38` | Secondary | User service (`systemctl --user`) |

---

## Active Providers — omos-core

### Verified 2026-05-02

| # | Provider | Plugin | Status | Verified | Notes |
|---|----------|--------|--------|----------|-------|
| 1 | **LiteLLM / DigitalOcean GenAI** | `litellm` | ✅ PASS | API returns correct answers | 18 models via DO inference |
| 2 | **Google Gemini** | `google` | ✅ PASS | `gemini-2.5-flash` returns correct answers | 13 models in dropdown |
| 3 | **Groq** | `groq` | ✅ PASS | Raw curl verified; CLI has parser bug | 14 models in dropdown |
| 4 | **Deepgram** | `deepgram` | ✅ PASS | Auth valid, 1 project found | STT + real-time transcription |
| 5 | **Fal** | `fal` | ⚠️ AUTH VALID | Key accepted, balance exhausted | Needs top-up at fal.ai/dashboard/billing |
| 6 | **Amazon Bedrock** | `amazon-bedrock` | ✅ PASS | AWS identity confirmed | 129 models (us-east-1) |
| 7 | **Azure OpenAI** | `azure-openai` | 🔧 System drop-in only | Key + endpoint configured (root service) | `onemind-openai.openai.azure.com` |
| 8 | **Google Vertex** | `google-vertex` | 🔧 System drop-in only | Project + region configured (root service) | `legacy-aio-474306` / us-central1 |

### Disabled / Removed Providers (as of 2026-05-02)

| Provider | Reason | Date Removed |
|----------|--------|-------------|
| Cerebras | `payment_required` — free tier exhausted | 2026-05-01 |
| OpenRouter | No credits on account | 2026-05-01 |
| MiniMax | Insufficient balance | 2026-05-01 |
| Kimi / Moonshot | Account suspended | 2026-05-01 |

### Known Issues

- **Groq CLI parser bug**: OpenClaw v2026.4.29 `groq-native` response parser returns "No text output" even when raw API returns valid responses. Groq API works perfectly via raw curl. Fixed in a future OpenClaw update.
- **Fal billing**: Auth is valid but account balance is exhausted. Image/video generation requires top-up at `fal.ai/dashboard/billing`.

---

## Systemd Drop-in Configuration

### omos-core — System service (`/etc/systemd/system/openclaw-gateway.service.d/`)

| File | Provider | Key Variable(s) |
|------|----------|-----------------|
| `azure.conf` | Azure OpenAI | `AZURE_OPENAI_API_KEY`, `AZURE_OPENAI_ENDPOINT` |
| `bedrock.conf` | Amazon Bedrock | `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_DEFAULT_REGION` |
| `deepgram.conf` | Deepgram | `DEEPGRAM_API_KEY` |
| `digitalocean.conf` | LiteLLM / DO GenAI | `LITELLM_API_KEY`, `LITELLM_BASE_URL` |
| `fal.conf` | Fal | `FAL_KEY` |
| `google.conf` | Google Gemini | `GOOGLE_API_KEY` |
| `groq.conf` | Groq | `GROQ_API_KEY` |
| `vertex.conf` | Google Vertex | `GOOGLE_CLOUD_PROJECT`, `GOOGLE_CLOUD_LOCATION` |

### omos-core — User service (`~/.config/systemd/user/openclaw-gateway.service.d/`)

| File | Provider |
|------|----------|
| `bedrock.conf` | Amazon Bedrock |
| `deepgram.conf` | Deepgram |
| `digitalocean.conf` | LiteLLM / DO GenAI |
| `fal.conf` | Fal |
| `google.conf` | Google Gemini |
| `groq.conf` | Groq |

### omos-sentinel — User service (`~/.config/systemd/user/openclaw-gateway.service.d/`)

| File | Provider | Status |
|------|----------|--------|
| `bedrock.conf` | Amazon Bedrock | ✅ deployed |
| `deepgram.conf` | Deepgram | ✅ deployed |
| `digitalocean.conf` | LiteLLM / DO GenAI | ✅ deployed |
| `fal.conf` | Fal | ✅ deployed |

> **Not yet deployed to sentinel:** Google Gemini, Groq

---

## Auth Profiles

### omos-core — ubuntu/main (`/home/ubuntu/.openclaw/agents/main/agent/auth-profiles.json`)

| Profile Key | Provider |
|-------------|----------|
| `deepgram:api-key` | Deepgram |
| `fal:api-key` | Fal |
| `google:api-key` | Google Gemini |
| `groq:api-key` | Groq |
| `litellm:api-key` | LiteLLM / DO GenAI |

### omos-core — root/legacy (`/root/.openclaw/agents/legacy/agent/auth-profiles.json`)

| Profile Key | Provider |
|-------------|----------|
| `deepgram:api-key` | Deepgram |
| `fal:api-key` | Fal |
| `github-copilot:github` | GitHub Copilot (OAuth) |
| `google:api-key` | Google Gemini |
| `groq:api-key` | Groq |
| `litellm:api-key` | LiteLLM / DO GenAI |
| `microsoft-foundry:default` | Microsoft Foundry |

---

## Models in Dropdown (models.json)

Three copies of `models.json` exist on omos-core:

### 1. `/home/ubuntu/.openclaw/agents/main/agent/models.json`

| Provider | Model Count |
|----------|------------|
| `codex` | 3 |
| `bedrock` | 1 |
| `amazon-bedrock` | 129 |
| `litellm` | 0 (uses API discovery) |
| `google` | 13 |
| `groq` | 14 |

### 2. `/root/.openclaw/agents/legacy/agent/models.json`

Providers: `codex, ollama, google, bedrock, amazon-bedrock, microsoft-foundry, groq`

### 3. `/root/.openclaw/agents/main/agent/models.json`

Providers: `codex, ollama, google, bedrock, amazon-bedrock, microsoft-foundry, anthropic-vertex, litellm, cerebras, minimax, minimax-cn, moonshot, openrouter, minimax-portal, minimax-portal-cn, groq`

> Note: root/main still has disabled provider entries (cerebras, minimax, moonshot, openrouter) — harmless since no keys are configured, but could be cleaned up.

---

## Google Gemini — 13 Models

| Model ID | Notes |
|----------|-------|
| `gemini-2.5-flash` | Primary fast model |
| `gemini-2.5-flash-lite` | Lightweight variant |
| `gemini-2.5-pro` | Best quality |
| `gemini-2.0-flash` | Previous gen fast |
| `gemini-2.0-flash-lite` | Previous gen lightweight |
| `gemini-3-flash-preview` | Next gen preview |
| `gemini-3-pro-preview` | Next gen preview |
| `gemini-3.1-pro-preview` | Latest preview |
| `gemini-3.1-flash-lite-preview` | Latest lightweight preview |
| `gemma-4-31b-it` | Open-weight 31B |
| `gemma-4-26b-a4b-it` | Open-weight 26B |
| `gemini-1.5-pro` | Legacy |
| `gemini-1.5-flash` | Legacy |

---

## Groq — 14 Models

| Model ID | Notes |
|----------|-------|
| `llama-3.3-70b-versatile` | Best general Llama |
| `llama-3.1-8b-instant` | Fast small model |
| `deepseek-r1-distill-llama-70b` | Reasoning |
| `groq/compound` | Compound AI |
| `groq/compound-mini` | Compound AI (small) |
| `qwen/qwen3-32b` | Alibaba 32B |
| `qwen-qwq-32b` | Reasoning variant |
| `meta-llama/llama-4-maverick-17b-128e-instruct` | Llama 4 Maverick |
| `meta-llama/llama-4-scout-17b-16e-instruct` | Llama 4 Scout |
| `gemma2-9b-it` | Google 9B |
| `mistral-saba-24b` | Mistral 24B |
| `openai/gpt-oss-120b` | GPT OSS 120B |
| `openai/gpt-oss-20b` | GPT OSS 20B |
| `moonshotai/kimi-k2-instruct` | Kimi K2 |

---

## LiteLLM / DigitalOcean GenAI — 18 Models

Endpoint: `https://inference.do-ai.run/v1`

| Model ID | Size |
|----------|------|
| `llama3.3-70b-instruct` | 70B |
| `llama3.2-3b-instruct` | 3B |
| `llama3.2-1b-instruct` | 1B |
| `llama3.1-8b-instruct` | 8B |
| `llama3.1-70b-instruct` | 70B |
| `mistral-nemo-instruct-2407` | 12B |
| `Qwen2.5-VL-7B-Instruct` | 7B |
| `Qwen2.5-VL-72B-Instruct` | 72B |
| `Qwen2.5-Coder-32B-Instruct` | 32B |
| `s1-32B` | 32B |
| `deepseek-r1-distill-llama-70b` | 70B |
| `QwQ-32B` | 32B |
| `deepseek-r1-distill-qwen-32b` | 32B |
| `Phi-4-multimodal-instruct` | ~14B |
| `deepseek-r1` | 671B |
| `deepseek-v3-0324` | 685B |
| `Llama-4-Maverick-17B-128E-Instruct` | 17B |
| `Llama-4-Scout-17B-16E-Instruct` | 17B |

---

## Media / Multimodal Capabilities

### Deployed Media Providers

| Capability | Provider | Status |
|-----------|----------|--------|
| **Image Generation** | `fal` | ⚠️ Auth valid, needs billing top-up |
| **Video Generation** | `fal` | ⚠️ Auth valid, needs billing top-up |
| **Speech-to-Text** | `deepgram` | ✅ Working (nova-3) |
| **Real-time Transcription** | `deepgram` | ✅ Available |
| **Vision / Image Understanding** | `google` (26 models), `groq` (2 models), `amazon-bedrock` (69 models) | ✅ Working |

### Available but NOT Deployed (need API keys)

| Provider | Media Models | Capabilities |
|----------|-------------|-------------|
| `openai` | 39 | Vision |
| `anthropic` | 23 | Vision |
| `openrouter` | 134 | Vision (needs credits) |
| `vercel-ai-gateway` | 92 | Vision |
| `xai` | 10 | Vision |
| `mistral` | 11 | Vision |
| `alibaba/qwen` | — | Vision (API key pending) |

---

## Gateway Configuration

### omos-core

- **System service**: `sudo systemctl restart openclaw-gateway`
- **Gateway token**: `8759c8a97ba63a47fd868c8ea01e7c0307b380f8f656c908`
- **Gateway port**: 18789
- **Config path (ubuntu)**: `/home/ubuntu/.openclaw/`
- **Config path (root)**: `/root/.openclaw/`

### omos-sentinel

- **User service**: `systemctl --user restart openclaw-gateway`
- **Gateway token**: `YEDd55cuHTle0QFE8eBYgwGbzdfIYeQf`

---

## API Key References

All keys stored in 1Password vault: **00-24 SO (Sovereign Operations)**

| Provider | Key Prefix / Identifier |
|----------|------------------------|
| LiteLLM / DO GenAI | `sk-do-eAUC8q...` |
| Google Gemini | `AIzaSyAq...` |
| Groq | `gsk_c2zx...` |
| Deepgram | `caf1f180...` |
| Fal | `ddb90ddd-...` |
| AWS Bedrock | `AKIA6CU7...` (IAM key) |
| Azure OpenAI | `4rfbO5HN...` |

---

## Pending Work

- [ ] **Fal**: Top up billing at fal.ai/dashboard/billing for image/video generation
- [ ] **Sentinel**: Deploy Google Gemini and Groq (only litellm, deepgram, fal, bedrock deployed so far)
- [ ] **Alibaba/Qwen**: Deploy when API key is obtained (plugin: `qwen`, env: `QWEN_API_KEY` or `DASHSCOPE_API_KEY`)
- [ ] **Root models.json cleanup**: Remove disabled provider entries (cerebras, minimax, moonshot, openrouter) from `/root/.openclaw/agents/main/agent/models.json`
- [ ] **OpenClaw upgrade**: Next release should fix Groq CLI parser bug (`groq-native` endpoint class)

---

## Deployment Pattern Reference

To add a new provider to omos-core:

```bash
# 1. Create systemd drop-in
sudo tee /etc/systemd/system/openclaw-gateway.service.d/<provider>.conf << 'EOF'
[Service]
Environment="<ENV_VAR>=<value>"
EOF

# 2. Reload and restart
sudo systemctl daemon-reload
sudo systemctl restart openclaw-gateway

# 3. Configure auth profile
openclaw plugins enable <plugin-name>
# Follow interactive auth setup

# 4. Add models to dropdown (if applicable)
# Edit models.json at the appropriate agent path
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
