# OpenClaw Gateway Config — Agent Handoff

**Date**: 2026-05-01  
**Status**: 77/77 models configured. DigitalOcean GenAI ✅ deployed 2026-05-01.  
**Purpose**: Brief the next agent on everything already done and what remains.

---

## Fleet

| Machine | Tailscale IP | SSH User | Service | Gateway Token |
|---|---|---|---|---|
| omos-core | 100.123.126.34 | ubuntu | `sudo systemctl [start/stop/restart/status] openclaw-gateway` | `8759c8a97ba63a47fd868c8ea01e7c0307b380f8f656c908` |
| omos-sentinel | 100.66.235.38 | ubuntu | `systemctl --user [start/stop/restart/status] openclaw-gateway` | `YEDd55cuHTle0QFE8eBYgwGbzdfIYeQf` |

- Config on core: `/root/.openclaw/openclaw.json` (requires sudo)
- Config on sentinel: `/home/ubuntu/.openclaw/openclaw.json`
- Model allowlist lives at: `d["agents"]["defaults"]["models"]` — a dict of `"provider/model-id": {"alias": "Display Name"}`

**CRITICAL**: Never edit `openclaw.json` while the gateway is running — it race-writes from memory. Always stop → edit → start.

---

## What Has Been Done

### Providers Fully Configured ✅

#### 1. GitHub Copilot [gh] — 26 models
- Auth: `github-copilot:github` token in `auth-profiles.json` on both boxes
- Models include: Claude Sonnet/Opus/Haiku (4–4.7), GPT-4.1/4o/5/5.x/5.x Codex, Gemini 2.5 Pro / 3 Flash/Pro/3.1 Pro, Grok Code Fast
- Smoke tested ✅

#### 2. AWS Bedrock [aws] — 15 models
- Auth: IAM credentials in systemd drop-in `/etc/systemd/system/openclaw-gateway.service.d/bedrock.conf` (core) and `/home/ubuntu/.config/systemd/user/openclaw-gateway.service.d/bedrock.conf` (sentinel)
- Env vars: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_DEFAULT_REGION=us-east-1`
- Models: Claude Sonnet 4.5/4.6, Opus 4.6, Haiku 4.5 [us.*], Nova Lite/Pro/Micro/Premier/2 Lite, Llama 4 Maverick/Scout, Llama 3.3 70B, DeepSeek R1, DeepSeek V3.2, Mistral Large 3
- Smoke tested ✅

#### 3. Azure OpenAI (microsoft-foundry) [az] — 10 models
- Auth: `api-key` header, endpoint `https://onemind-openai.openai.azure.com/openai/v1`
- API key in `models.json` providers section on both boxes
- Deployments: gpt-4.1, gpt-4.1-mini, gpt-4.1-nano, gpt-4o, gpt-4o-mini, gpt-5, gpt-5-chat-latest, gpt-5-mini, gpt-5-nano, o4-mini
- 1Password item: `fsa2tl43hguv4fbnmbglxrp5qy` (vault: "00-24 SO (Sovereign Operations)")
- Smoke tested ✅

#### 4. Google Gemini [gg] — 8 models
- Transport: `google-generative-ai` (NOT Vertex — OpenClaw has a bug parsing Vertex SSE streams)
- API key: `AIzaSyABJCwnwZ76G6ZqOnlAucyYDq7fdXz1Bzc` (project: `gen-lang-client-0596829426`)
- Auth configured 3 ways (belt-and-suspenders):
  1. `GEMINI_API_KEY` env var in systemd drop-in `google.conf` on both boxes
  2. `apiKey` field in `agents/main/agent/models.json` providers.google on both boxes
  3. `google:default` entry in `agents/main/agent/auth-profiles.json` on sentinel only
- **DO NOT TOUCH**: `AIzaSyA0ua8lZoX42vufWrdfNDaujbLwbell6mM` — that's the voice API key, separate project
- Billing: `01A74D-0EC544-797F99 / Legacy AI`
- 1Password item: `hjw67fllsja5hapvmuu2yiarq4`
- Models: gemini-2.0-flash, gemini-2.5-flash, gemini-2.5-flash-lite, gemini-2.5-pro, gemini-3-flash-preview, gemini-3-pro-preview, gemini-3.1-pro-preview, gemini-3.1-flash-lite-preview
- Smoke tested on both boxes ✅

#### 5. DigitalOcean GenAI (via litellm bridge) [do] — 18 models ✅ deployed 2026-05-01
- **Bridge**: Uses the `litellm` plugin (the only OpenClaw plugin with `allowExplicitBaseUrl: true`) pointed at DO's OpenAI-compatible endpoint
- **Base URL**: `https://inference.do-ai.run/v1` set in `models.providers.litellm.baseUrl`
- **API Key**: 1Password `op://aaeyris5z44pdqspk277u2oppi/DigitalOcean GenAI Inference Key/credential`
- **Auth architecture** (two layers required):
  1. `LITELLM_API_KEY` env var injected via systemd drop-in `digitalocean.conf` on both boxes (so the gateway service can call DO)
  2. `auth-profiles.json` entry `"litellm:api-key"` in `~/.openclaw/agents/main/agent/auth-profiles.json` on both boxes (so the CLI can authenticate)
- **Plugin enablement**: `plugins.entries.litellm = {"enabled": true}` in `openclaw.json` (do NOT use `openclaw plugins enable litellm` — it hangs; edit JSON directly)
- **Smoke tested**: All 18 models connect. 8 PASS on terse prompts; 10 prefer natural-language prompts (all confirmed functional with "What is 2+2?" test)
- **omos-core topology note**: Core has ONE config — `/home/ubuntu/.openclaw/openclaw.json` (user systemd service, `systemctl --user`). No system-level service exists. No `/root/.openclaw/` path.

**Drop-in locations**:
- Core (user systemd): `/home/ubuntu/.config/systemd/user/openclaw-gateway.service.d/digitalocean.conf`
- Sentinel (user systemd): `/home/ubuntu/.config/systemd/user/openclaw-gateway.service.d/digitalocean.conf`
- ~~Core system: `/etc/systemd/system/openclaw-gateway.service.d/digitalocean.conf`~~ — does not exist (no system service on core)

---

## What Remains ❌

Nothing — all 5 providers deployed and smoke tested. 77 models total.

---

## How to Add More Models (Pattern)

```bash
# 1. Stop gateway
ssh ubuntu@100.123.126.34 'sudo systemctl stop openclaw-gateway'
ssh ubuntu@100.66.235.38 'systemctl --user stop openclaw-gateway'

# 2. Edit config — add to d["agents"]["defaults"]["models"]
ssh ubuntu@100.123.126.34 'sudo python3 - << EOF
import json
d = json.load(open("/root/.openclaw/openclaw.json"))
d["agents"]["defaults"]["models"]["PROVIDER/MODEL_ID"] = {"alias": "Display Name [tag]"}
json.dump(d, open("/root/.openclaw/openclaw.json","w"), indent=2)
print("done")
EOF'

# Repeat for sentinel (no sudo, different path)

# 3. Restart
ssh ubuntu@100.123.126.34 'sudo systemctl start openclaw-gateway'
ssh ubuntu@100.66.235.38 'systemctl --user start openclaw-gateway'
```

## Current Full Model List (77 models)

### GitHub Copilot [gh] — 26
| OpenClaw Key | Alias |
|---|---|
| github-copilot/claude-sonnet-4.6 | Sonnet 4.6 [gh] |
| github-copilot/claude-opus-4.5 | Opus 4.5 [gh] |
| github-copilot/claude-haiku-4.5 | Haiku 4.5 [gh] |
| github-copilot/claude-opus-4.6 | Opus 4.6 [gh] |
| github-copilot/claude-opus-4.7 | Opus 4.7 [gh] |
| github-copilot/claude-sonnet-4 | Sonnet 4 [gh] |
| github-copilot/claude-sonnet-4.5 | Sonnet 4.5 [gh] |
| github-copilot/gemini-2.5-pro | Gemini 2.5 Pro [gh] |
| github-copilot/gemini-3-flash-preview | Gemini 3 Flash [gh] |
| github-copilot/gemini-3-pro-preview | Gemini 3 Pro [gh] |
| github-copilot/gemini-3.1-pro-preview | Gemini 3.1 Pro [gh] |
| github-copilot/gpt-4.1 | GPT-4.1 [gh] |
| github-copilot/gpt-4o | GPT-4o [gh] |
| github-copilot/gpt-5 | GPT-5 [gh] |
| github-copilot/gpt-5-mini | GPT-5 Mini [gh] |
| github-copilot/gpt-5.1 | GPT-5.1 [gh] |
| github-copilot/gpt-5.1-codex | GPT-5.1 Codex [gh] |
| github-copilot/gpt-5.1-codex-max | GPT-5.1 Codex Max [gh] |
| github-copilot/gpt-5.1-codex-mini | GPT-5.1 Codex Mini [gh] |
| github-copilot/gpt-5.2 | GPT-5.2 [gh] |
| github-copilot/gpt-5.2-codex | GPT-5.2 Codex [gh] |
| github-copilot/gpt-5.3-codex | GPT-5.3 Codex [gh] |
| github-copilot/gpt-5.4 | GPT-5.4 [gh] |
| github-copilot/gpt-5.4-mini | GPT-5.4 Mini [gh] |
| github-copilot/gpt-5.5 | GPT-5.5 [gh] |
| github-copilot/grok-code-fast-1 | Grok Code Fast [gh] |

### AWS Bedrock [aws] — 15
| OpenClaw Key | Alias |
|---|---|
| amazon-bedrock/us.anthropic.claude-sonnet-4-6 | Sonnet 4.6 [aws] |
| amazon-bedrock/us.anthropic.claude-opus-4-6 | Opus 4.6 [aws] |
| amazon-bedrock/us.anthropic.claude-haiku-4-5-20251001-v1:0 | Haiku 4.5 [aws] |
| amazon-bedrock/us.anthropic.claude-sonnet-4-5-20250929-v1:0 | Sonnet 4.5 [aws] |
| amazon-bedrock/amazon.nova-pro-v1:0 | Nova Pro [aws] |
| amazon-bedrock/amazon.nova-lite-v1:0 | Nova Lite [aws] |
| amazon-bedrock/amazon.nova-micro-v1:0 | Nova Micro [aws] |
| amazon-bedrock/amazon.nova-premier-v1:0 | Nova Premier [aws] |
| amazon-bedrock/amazon.nova-2-lite-v1:0 | Nova 2 Lite [aws] |
| amazon-bedrock/meta.llama4-maverick-17b-instruct-v1:0 | Llama 4 Maverick [aws] |
| amazon-bedrock/meta.llama4-scout-17b-instruct-v1:0 | Llama 4 Scout [aws] |
| amazon-bedrock/meta.llama3-3-70b-instruct-v1:0 | Llama 3.3 70B [aws] |
| amazon-bedrock/deepseek.r1-v1:0 | DeepSeek R1 [aws] |
| amazon-bedrock/deepseek.v3.2 | DeepSeek V3.2 [aws] |
| amazon-bedrock/mistral.mistral-large-3-675b-instruct | Mistral Large 3 [aws] |

### Azure OpenAI / microsoft-foundry [az] — 10
| OpenClaw Key | Alias |
|---|---|
| microsoft-foundry/gpt-4.1 | GPT-4.1 [az] |
| microsoft-foundry/gpt-4.1-mini | GPT-4.1 Mini [az] |
| microsoft-foundry/gpt-4.1-nano | GPT-4.1 Nano [az] |
| microsoft-foundry/gpt-4o | GPT-4o [az] |
| microsoft-foundry/gpt-4o-mini | GPT-4o Mini [az] |
| microsoft-foundry/gpt-5 | GPT-5 [az] |
| microsoft-foundry/gpt-5-chat-latest | GPT-5 Chat [az] |
| microsoft-foundry/gpt-5-mini | GPT-5 Mini [az] |
| microsoft-foundry/gpt-5-nano | GPT-5 Nano [az] |
| microsoft-foundry/o4-mini | o4-mini [az] |

### Google Gemini [gg] — 8
| OpenClaw Key | Alias |
|---|---|
| google/gemini-2.0-flash | Gemini 2.0 Flash [gg] |
| google/gemini-2.5-flash | Gemini 2.5 Flash [gg] |
| google/gemini-2.5-flash-lite | Gemini 2.5 Flash Lite [gg] |
| google/gemini-2.5-pro | Gemini 2.5 Pro [gg] |
| google/gemini-3-flash-preview | Gemini 3 Flash [gg] |
| google/gemini-3-pro-preview | Gemini 3 Pro [gg] |
| google/gemini-3.1-pro-preview | Gemini 3.1 Pro [gg] |
| google/gemini-3.1-flash-lite-preview | Gemini 3.1 Flash Lite [gg] |

---

## Systemd Drop-in Locations

**Core system** (need `sudo`):
- `/etc/systemd/system/openclaw-gateway.service.d/bedrock.conf`
- `/etc/systemd/system/openclaw-gateway.service.d/azure.conf`
- `/etc/systemd/system/openclaw-gateway.service.d/google.conf`
- `/etc/systemd/system/openclaw-gateway.service.d/digitalocean.conf`
- `/etc/systemd/system/openclaw-gateway.service.d/vertex.conf` (harmless, unused)

**Core ubuntu user** (no sudo, ubuntu's user service — currently FAILED due to port conflict with system service, but env is still read by CLI):
- `/home/ubuntu/.config/systemd/user/openclaw-gateway.service.d/digitalocean.conf`

**Sentinel user** (no sudo):
- `/home/ubuntu/.config/systemd/user/openclaw-gateway.service.d/bedrock.conf`
- `/home/ubuntu/.config/systemd/user/openclaw-gateway.service.d/azure.conf`
- `/home/ubuntu/.config/systemd/user/openclaw-gateway.service.d/google.conf`
- `/home/ubuntu/.config/systemd/user/openclaw-gateway.service.d/digitalocean.conf`
- `/home/ubuntu/.config/systemd/user/openclaw-gateway.service.d/vertex.conf` (harmless, unused)

Drop-in format:
```ini
[Service]
Environment="KEY=value"
```
After editing drop-ins: `sudo systemctl daemon-reload && sudo systemctl restart openclaw-gateway` (core) or `systemctl --user daemon-reload && systemctl --user restart openclaw-gateway` (sentinel).

---

### DigitalOcean GenAI (via litellm) [do] — 18

| OpenClaw Key | Alias | Pricing (in/out per M tokens) |
|---|---|---|
| litellm/alibaba-qwen3-32b | Qwen3-32B [do] | $0.25 / $0.55 |
| litellm/qwen3-coder-flash | Qwen3 Coder Flash [do] | ~$0.25 / $0.55 |
| litellm/qwen3.5-397b-a17b | Qwen3.5 397B [do] | $0.55 / $3.50 |
| litellm/deepseek-r1-distill-llama-70b | DeepSeek R1 Distill 70B [do] | ~$0.30 / $0.60 |
| litellm/deepseek-3.2 | DeepSeek V3.2 [do] | $0.50 / $1.60 |
| litellm/gemma-4-31B-it | Gemma 4 31B [do] | $0.18 / $0.50 |
| litellm/llama3.3-70b-instruct | Llama 3.3 70B [do] | ~$0.20 / $0.40 |
| litellm/llama-4-maverick | Llama 4 Maverick [do] | $0.25 / $0.87 |
| litellm/mistral-3-14B | Ministral 3 14B [do] | $0.20 / $0.20 |
| litellm/kimi-k2.5 | Kimi K2.5 [do] | $0.50 / $2.70 |
| litellm/minimax-m2.5 | MiniMax M2.5 [do] | $0.30 / $1.20 |
| litellm/nvidia-nemotron-3-super-120b | Nemotron-3 Super 120B [do] | ~$0.40 / $1.20 |
| litellm/nemotron-3-nano-omni | Nemotron Nano 3 Omni [do] | ~$0.10 / $0.30 |
| litellm/nemotron-nano-12b-v2-vl | Nemotron Nano 12B VL [do] | ~$0.10 / $0.30 |
| litellm/openai-gpt-oss-120b | GPT-OSS 120B [do] | ~$0.30 / $1.00 |
| litellm/openai-gpt-oss-20b | GPT-OSS 20B [do] | $0.05 / $0.45 |
| litellm/glm-5 | GLM 5 [do] | $1.00 / $3.20 |
| litellm/arcee-trinity-large-thinking | Arcee Trinity Large [do] | ~$0.50 / $1.50 |


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
