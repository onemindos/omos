# DigitalOcean AI Models
*Last updated: 2026-06-15*

Pricing is per 1M tokens (input / output). Open/DO-hosted model prices are estimates based on market rates — verify exact rates at cloud.digitalocean.com → AI/ML → GenAI Platform → Pricing.

---

## ANTHROPIC (Claude)

| Model | Best For | Context | Input /1M | Output /1M |
|-------|----------|---------|-----------|------------|
| `anthropic-claude-opus-4.8` | Best all-around, 1M context, complex reasoning, long docs | 1M / 128K out | ~$15.00 | ~$75.00 |
| `anthropic-claude-opus-4.7` | High-quality reasoning, coding, analysis | 200K / 8K out | ~$15.00 | ~$75.00 |
| `anthropic-claude-opus-4.6` | Complex tasks, strategic analysis | 200K / 8K out | $15.00 | $75.00 |
| `anthropic-claude-opus-4.5` | Reasoning, long-form writing | 200K / 8K out | ~$15.00 | ~$75.00 |
| `anthropic-claude-opus-4` | General Opus-tier quality | 200K / 8K out | ~$15.00 | ~$75.00 |
| `anthropic-claude-4.1-opus` | Reasoning, long context | 200K / 8K out | ~$15.00 | ~$75.00 |
| `anthropic-claude-4.6-sonnet` | **Best value Anthropic** — coding, everyday tasks | 200K / 8K out | $3.00 | $15.00 |
| `anthropic-claude-4.5-sonnet` | Balanced tasks, coding | 200K / - | ~$3.00 | ~$15.00 |
| `anthropic-claude-sonnet-4` | General tasks, writing | 200K / 8K out | ~$3.00 | ~$15.00 |
| `anthropic-claude-haiku-4.5` | **Cheapest Anthropic** — fast routing, summaries, triage | 200K / 8K out | $0.80 | $4.00 |

---

## OPENAI

| Model | Best For | Context | Input /1M | Output /1M |
|-------|----------|---------|-----------|------------|
| `openai-gpt-5.5` | Best GPT, massive context | 1M / 128K out | ~$15.00 | ~$60.00 |
| `openai-gpt-5.4-pro` | Top quality + 1M context | 1M / 128K out | ~$10.00 | ~$40.00 |
| `openai-gpt-5.4` | Strong general purpose | 400K / 128K out | ~$5.00 | ~$20.00 |
| `openai-gpt-5.3-codex` | Coding focused | 400K / 128K out | ~$5.00 | ~$20.00 |
| `openai-gpt-5.2-pro` | Pro quality reasoning | 400K / 128K out | ~$5.00 | ~$20.00 |
| `openai-gpt-5.2` | Solid general use | 400K / 128K out | ~$3.00 | ~$12.00 |
| `openai-gpt-5.1-codex-max` | **Best coding on DO** | 400K / 128K out | ~$3.00 | ~$12.00 |
| `openai-gpt-5` | Flagship GPT-5, general | 400K / 128K out | $2.50 | $10.00 |
| `openai-gpt-5-mini` | Balanced cost/quality | 400K / 128K out | $0.40 | $1.60 |
| `openai-gpt-5-nano` | **Cheapest OpenAI** — fast tasks | 400K / 128K out | $0.15 | $0.60 |
| `openai-gpt-5.4-mini` | Fast mid-tier | 400K / 128K out | ~$0.40 | ~$1.60 |
| `openai-gpt-5.4-nano` | Ultra-fast cheap | 400K / 128K out | ~$0.15 | ~$0.60 |
| `openai-o3` | Deep reasoning, math, logic | 200K / 100K out | $10.00 | $40.00 |
| `openai-o1` | Reasoning, science, analysis | 200K / 100K out | ~$10.00 | ~$40.00 |
| `openai-o3-mini` | Fast reasoning | 200K / 100K out | $1.10 | $4.40 |
| `openai-gpt-4.1` | Reliable GPT-4, large context | 1M / 32K out | ~$2.00 | ~$8.00 |
| `openai-gpt-4o` | Multimodal, vision + text | 128K / 16K out | $2.50 | $10.00 |
| `openai-gpt-4o-mini` | Fast cheap multimodal | 128K / 16K out | $0.15 | $0.60 |
| `openai-gpt-oss-120b` | Open-source large, general | 128K / 4K out | ~$0.90 | ~$0.90 |
| `openai-gpt-oss-20b` | Open-source small, fast | 128K / 4K out | ~$0.20 | ~$0.20 |

---

## DEEPSEEK

| Model | Best For | Context | Input /1M | Output /1M |
|-------|----------|---------|-----------|------------|
| `deepseek-v4-pro` | **1M context**, coding, deep reasoning | 1M / 1M out | ~$0.27 | ~$1.10 |
| `deepseek-4-flash` | Fast DeepSeek, huge context | 262K / 262K out | ~$0.15 | ~$0.60 |
| `deepseek-3.2` | General reasoning, coding | 128K / 64K out | ~$0.27 | ~$1.10 |
| `deepseek-r1-distill-llama-70b` | Chain-of-thought reasoning | 32K / 8K out | ~$0.23 | ~$0.69 |

---

## META (Llama)

| Model | Best For | Context | Input /1M | Output /1M |
|-------|----------|---------|-----------|------------|
| `llama-4-maverick` | Fast, efficient, multimodal | 128K / 16K out | ~$0.20 | ~$0.65 |
| `llama3.3-70b-instruct` | Strong open-source general | 128K / 4K out | ~$0.23 | ~$0.40 |

---

## NVIDIA (Nemotron)

| Model | Best For | Context | Input /1M | Output /1M |
|-------|----------|---------|-----------|------------|
| `nvidia-nemotron-3-super-120b` | Strong reasoning, **1M context** | 1M / 32K out | ~$1.00 | ~$4.00 |
| `nemotron-3-ultra-550b` | Massive model, complex tasks | 131K / 131K out | ~$3.00 | ~$12.00 |
| `nemotron-nano-12b-v2-vl` | Vision + language, lightweight | 128K / 16K out | ~$0.10 | ~$0.30 |
| `nemotron-3-nano-omni` | Fast multimodal | 65K / 65K out | ~$0.10 | ~$0.30 |

---

## ALIBABA / QWEN

| Model | Best For | Context | Input /1M | Output /1M |
|-------|----------|---------|-----------|------------|
| `qwen3.5-397b-a17b` | Massive MoE, reasoning + coding | 131K / 81K out | ~$0.50 | ~$2.00 |
| `qwen3-coder-flash` | **Fast coding**, 262K context | 262K / 65K out | ~$0.15 | ~$0.60 |
| `alibaba-qwen3-32b` | General purpose, multilingual | 32K / 32K out | ~$0.15 | ~$0.60 |

---

## KIMI (Moonshot)

| Model | Best For | Context | Input /1M | Output /1M |
|-------|----------|---------|-----------|------------|
| `kimi-k2.6` | Long context, reasoning, agentic | 262K / 262K out | ~$0.50 | ~$2.00 |
| `kimi-k2.5` | Long context, general tasks | 256K / 32K out | ~$0.50 | ~$2.00 |

---

## OTHER

| Model | Best For | Context | Input /1M | Output /1M |
|-------|----------|---------|-----------|------------|
| `mistral-3-14B` | Lightweight, fast, cheap | 262K / - | ~$0.10 | ~$0.30 |
| `gemma-4-31B-it` | Google open model, general | 256K / 8K out | ~$0.15 | ~$0.30 |
| `glm-5` | Chinese/multilingual tasks | 128K / 128K out | ~$0.20 | ~$0.80 |
| `minimax-m2.5` | Long context, general | 200K / 128K out | ~$0.30 | ~$1.20 |
| `mimo-v2.5` | Math reasoning | 262K / 131K out | ~$0.20 | ~$0.80 |
| `arcee-trinity-large-thinking` | Thinking/reasoning tasks | 128K / 32K out | ~$0.30 | ~$1.20 |

---

## IMAGE GENERATION

| Model | Best For | Notes |
|-------|----------|-------|
| `openai-gpt-image-2` | Best image gen, photorealistic | Per image pricing |
| `openai-gpt-image-1.5` | High quality image gen | Per image pricing |
| `openai-gpt-image-1` | Standard image gen | Per image pricing |
| `stable-diffusion-3.5-large` | Open-source image gen | Per image pricing |
| `wan2-2-t2v-a14b` | Text-to-video | Per second pricing |

---

## AUDIO

| Model | Best For |
|-------|----------|
| `qwen3-tts-voicedesign` | Text-to-speech, voice design |

---

## EMBEDDINGS & RERANKING (no generation cost)

| Model | Best For |
|-------|----------|
| `bge-m3` | Multilingual embeddings, RAG |
| `bge-reranker-v2-m3` | Reranking search results |
| `e5-large-v2` | High quality English embeddings |
| `all-mini-lm-l6-v2` | Fast small embeddings |
| `gte-large-en-v1.5` | High quality embeddings |
| `qwen3-embedding-0.6b` | Qwen embeddings |
| `multi-qa-mpnet-base-dot-v1` | Q&A search embeddings |

---

## SMART ROUTERS (auto-select best model)

| Router | Routes to best model for... |
|--------|----------------------------|
| `router:software-engineering` | Coding & dev tasks |
| `router:writing` | Writing & editing |
| `router:general` | General purpose tasks |
| `router:knowledge-base-document` | Document Q&A |

---

## Quick Reference: Best Model by Task

| Task | Best Choice | Budget Choice |
|------|-------------|---------------|
| Coding | `openai-gpt-5.1-codex-max` or `anthropic-claude-4.6-sonnet` | `qwen3-coder-flash` |
| Deep reasoning | `openai-o3` or `deepseek-v4-pro` | `openai-o3-mini` |
| Long documents | `anthropic-claude-opus-4.8` (1M ctx) or `deepseek-v4-pro` | `nvidia-nemotron-3-super-120b` |
| Fast/cheap | `openai-gpt-5-nano` or `anthropic-claude-haiku-4.5` | `mistral-3-14B` |
| Multimodal/vision | `openai-gpt-4o` or `nemotron-nano-12b-v2-vl` | `llama-4-maverick` |
| Math | `openai-o3` or `mimo-v2.5` | `deepseek-r1-distill-llama-70b` |
| General best | `anthropic-claude-opus-4.8` or `openai-gpt-5.5` | `kimi-k2.6` |

---

*Prices marked ~$ are estimates. Verify at: cloud.digitalocean.com → AI/ML → GenAI Platform → Pricing*
