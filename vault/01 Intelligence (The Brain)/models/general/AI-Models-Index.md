---
id: so-01-ai-models-index
title: "AI Models — Master Index"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models]
ai_summary: "Quick-reference index of all AI model providers in the One Mind stack — Anthropic, OpenAI, Google, xAI, DigitalOcean, Alibaba Qwen, Moonshot Kimi, MiniMax, AWS Bedrock, ByteDance Doubao/Seedance, Gemma"
---

# AI Models — Master Index

> Single source of truth for every AI provider and model in the One Mind stack.
> Update when switching models, pricing changes, or new providers are added.

---

## Provider Quick Reference

| Provider | Primary Model | Tier | Key Use | File |
|----------|--------------|------|---------|------|
| **Anthropic** | claude-3.7-sonnet | PRIMARY | Agents, reasoning, code | [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The brain)/models/AI-Models/Anthropic]] |
| **OpenAI** | gpt-4o | SECONDARY | Broad tasks, embeddings | [[OpenAI]] |
| **Google** | gemini-2.0-flash | SECONDARY | Long context, multimodal | [[Google]] |
| **xAI** | grok-3 | EXPERIMENTAL | Real-time web, reasoning | [[xAI]] |
| **DigitalOcean** | llama-3.3-70b | INFRA | On-VPS inference, zero egress | [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The brain)/models/AI-Models/DigitalOcean]] |
| **Alibaba (Qwen)** | qwen-max | STACK | Long context (1M), multilingual | [[Alibaba]] |
| **Moonshot (Kimi)** | kimi-k1.5 | STACK | Extended reasoning, 128K CoT | [[Moonshot]] |
| **MiniMax** | MiniMax-Text-01 | STACK | Ultra-long context (4M tokens) | [[MiniMax]] |
| **AWS (Bedrock)** | claude / nova | STACK | AWS-integrated pipelines, IAM | [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The brain)/models/AI-Models/AWS]] |
| **ByteDance (Doubao)** | doubao-pro-128k | STACK | Video gen (Seedance), vision (Seed-VL), LLM | [[ByteDance]] |
| **Gemma (Google)** | gemma3:12b | OPEN-SOURCE | Local inference (Ollama), zero cost, privacy | [[Gemma]] |
| **Mistral AI** | mistral-large-latest | STACK | European/French AI, multilingual, code, open-weight | [[Mistral]] |
| **NVIDIA** | nemotron / cosmos / gr00t | SPECIALIST | Agentic AI, physical AI, robotics, AV, climate, biomedical | [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The brain)/models/AI-Models/NVIDIA]] |

---

## Agent → Model Mapping

| Agent | Provider | Model | Notes |
|-------|----------|-------|-------|
| **Legacy** (orchestrator) | Anthropic | claude-3.7-sonnet | Primary reasoning engine |
| **Spartan** (HP) | Anthropic | claude-3.5-haiku | Lightweight, fast responses |
| **Oracle** (HP) | Anthropic | claude-3.7-sonnet | Deep personal analysis |
| **Empire** (GE) | Anthropic | claude-3.7-sonnet | Business strategy |
| **Haven** (LE) | Anthropic | claude-3.5-haiku | Family/home tasks |
| **Guardian** (SO) | OpenAI | gpt-4o | Security analysis |
| **Grid** (SO) | DigitalOcean | llama-3.3-70b | Infra ops, low latency |
| **Trinity** (SO) | Anthropic | claude-3.7-sonnet | Meta-orchestration |
| **Eden** (LE) | Google | gemini-2.0-flash | Multimodal/environment |

> Update this table in openclaw.json and here when models change.

---

## Cost Overview (Monthly Estimate)

| Provider | Usage | Est. Cost/mo |
|----------|-------|-------------|
| Anthropic | Agents (Legacy, Oracle, Empire, Trinity) | $20–80 |
| OpenAI | Guardian + embeddings | $10–30 |
| Google | Eden + long-context tasks | $5–20 |
| DigitalOcean Gradient | Grid agent (bundled with VPS) | ~$0 extra |
| xAI | Experimental only | <$5 |
| Alibaba | Experimental only | <$5 |
| Moonshot | Experimental only | <$5 |
| **Total** | | **~$40–145/mo** |

---

## API Key Locations (1Password)

All keys stored in `1Password → 00-24 SO (Sovereign Operations)` vault.

| Provider | 1Password Item |
|----------|----------------|
| Anthropic | `Anthropic API Key` |
| OpenAI | `OpenAI API Key` |
| Google AI | `Google AI Studio API Key` |
| xAI | `xAI API Key` |
| DigitalOcean | `DigitalOcean API Token` (same token) |
| Alibaba Cloud | `Alibaba Cloud API Key` |
| Moonshot | `Moonshot API Key` |
| MiniMax | `MiniMax API Key` |
| AWS | `AWS Access Key ID` + `AWS Secret Access Key` |
| ByteDance ARK | `ByteDance ARK API Key` |
| Mistral AI | `Mistral API Key` |
| NVIDIA NIM | `NVIDIA NIM API Key` |
| Gemma (local) | No key — runs via Ollama on VPS |

---

## Decision Rules — Which Model to Use

```
REASONING / LONG TASK → claude-3.7-sonnet (extended thinking ON)
FAST / CHEAP RESPONSE → claude-3.5-haiku or gemini-2.0-flash
EMBEDDINGS          → text-embedding-3-small (OpenAI)
LONG CONTEXT (>200K) → gemini-1.5-pro OR moonshot-v1-128k OR kimi-k1.5
REAL-TIME WEB DATA  → grok-3
ON-VPS / ZERO COST  → llama-3.3-70b (Gradient) OR gemma3:12b (Ollama)
MULTILINGUAL        → qwen-max
IMAGE UNDERSTANDING → gemini-2.0-flash OR gpt-4o OR Seed-1.5-VL (ByteDance)
VIDEO GENERATION    → seedance-1-0-pro (ByteDance)
LOCAL / PRIVATE     → gemma3:12b via Ollama (no data leaves VPS)
EUROPEAN / GDPR     → mistral-small-latest (EU data residency)
CODE COMPLETION FIM → codestral-latest (Mistral)
MULTILINGUAL EURO   → mistral-nemo-latest (French / German / Italian)
AGENTIC REASONING   → nemotron-ultra-253b (NVIDIA) for reward modeling
PHYSICAL AI / ROBOT → Cosmos + GR00T (NVIDIA)
CLIMATE / PHYSICS   → Earth-2 FourCastNet (NVIDIA)
BIOMEDICAL / DRUG   → BioNeMo DiffDock (NVIDIA)
ULTRA-LONG CONTEXT  → MiniMax-Text-01 (4M tokens)
```

---

*See individual provider files for model specs, pricing tiers, and setup notes.*
