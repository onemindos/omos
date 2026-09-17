# DO Gradient vs Anthropic Claude — Model Comparison

*Generated: 2026-06-19 | Pricing sourced directly from `doctl gradient list-models` API*

---

## Anthropic Claude Pricing (via AWS Bedrock)

| Model | Model ID | Input/1M | Output/1M |
|---|---|---|---|
| Claude Haiku 4.5 | `claude-haiku-4-5` | $1.00 | $5.00 |
| Claude Sonnet 4.6 | `claude-sonnet-4-6` | $3.00 | $15.00 |
| Claude Opus 4.8 | `claude-opus-4-8` | $5.00 | $25.00 |
| Claude Fable 5 | `claude-fable-5` | $10.00 | $50.00 |

---

## DO Gradient Accessible Models (Your Tier)

All models below are confirmed working on your current DO account. Pricing is per 1M tokens.

### Haiku-Tier Equivalents — Fast, Cheap, Simple Tasks

*Reference: Claude Haiku 4.5 MMLU Pro ~73 | $1.00 in / $5.00 out*

| Model | Hermes Slug | Params | Context | MMLU Pro | Input/1M | Output/1M | Savings vs Haiku |
|---|---|---|---|---|---|---|---|
| Deepseek V4 Flash | `deepseek-4-flash` | 284B MoE | 262k | — | $0.112 | $0.224 | ~9–22x cheaper |
| OpenAI GPT-oss-20b | `openai-gpt-oss-20b` | ~20B | 128k | 74.8 | $0.050 | $0.450 | ~11–20x cheaper |
| Gemma 4 | `gemma-4-31b-it` | 31B | 256k | — | $0.180 | $0.500 | ~5–10x cheaper |
| Ministral 3 14B | `mistral-3-14b` | 14B | 262k | — | $0.200 | $0.200 | ~5–25x cheaper |
| Nemotron-nano 12b v2 (vision) | `nemotron-nano-12b-v2-vl` | 12B | 128k | — | $0.200 | $0.600 | ~5–8x cheaper |
| MiMo V2.5 | `mimo-v2.5` | ~7B | 32k | — | $0.140 | $0.280 | ~7–18x cheaper |

**Best pick:** `deepseek-4-flash` — 262k context, MoE architecture punches above its weight, cheapest blended rate (~$0.15/1M).

---

### Sonnet-Tier Equivalents — Balanced, Most Workloads

*Reference: Claude Sonnet 4.6 MMLU Pro 84.2 | $3.00 in / $15.00 out*

| Model | Hermes Slug | Params | Context | MMLU Pro | Input/1M | Output/1M | Savings vs Sonnet |
|---|---|---|---|---|---|---|---|
| **Deepseek 3.2** | `deepseek-3.2` | 680B MoE | 65k | **86.2** | $0.425 | $1.360 | ~7–11x cheaper |
| Llama 4 Maverick | `llama-4-maverick` | 400B MoE | 128k | — | $0.250 | $0.870 | ~12–17x cheaper |
| Kimi K2.5 | `kimi-k2.5` | 1T MoE | 256k | — | $0.375 | $2.025 | ~8–7x cheaper |
| MiniMax M2.5 | `minimax-m2.5` | 230B | 200k | — | $0.225 | $0.900 | ~13–17x cheaper |
| OpenAI GPT-oss-120b | `openai-gpt-oss-120b` | 117B | 128k | 80.8 | $0.070 | $0.490 | ~43–31x cheaper |
| Qwen3 Coder Flash *(coding)* | `qwen3-coder-flash` | 30B | 262k | — | $0.450 | $1.700 | ~7–9x cheaper |
| Qwen3 32B | `alibaba-qwen3-32b` | 32.8B | 32k | 79.8 | $0.250 | $0.550 | ~12–27x cheaper |
| DeepSeek R1 Distill 70B *(reasoning)* | `deepseek-r1-distill-llama-70b` | 70B | 32k | 79.5 | $0.990 | $0.990 | ~3–15x cheaper |

**Best pick:** `deepseek-3.2` — MMLU **86.2 beats Sonnet 4.6's 84.2** at 7–11x lower cost. Strongest general-purpose model in the accessible catalog.  
**Coding pick:** `qwen3-coder-flash` — purpose-built for code, 262k context.

---

### Opus-Tier Equivalents — Complex Reasoning, Agents

*Reference: Claude Opus 4.8 MMLU Pro 87.3 | $5.00 in / $25.00 out*

| Model | Hermes Slug | Params | Context | MMLU Pro | Input/1M | Output/1M | Savings vs Opus |
|---|---|---|---|---|---|---|---|
| **Deepseek V4 Pro** | `deepseek-v4-pro` | 1.6T MoE | **1M** | — | $1.392 | $2.784 | ~3.5–9x cheaper |
| Kimi K2.6 | `kimi-k2.6` | 1T MoE | 96k | — | $0.760 | $3.200 | ~6.5–8x cheaper |
| MiMo V2.5 Pro | `mimo-v2.5-pro` | 1T MoE | 262k | — | $0.800 | $3.000 | ~6–8x cheaper |
| Nemotron 3 Ultra | `nemotron-3-ultra-550b` | 550B | 131k | — | $0.900 | $1.700 | ~5.5–15x cheaper |
| GLM 5 | `glm-5` | 744B | 64k | — | $0.750 | $2.400 | ~6.5–10x cheaper |
| Llama 3.3 70B Instruct | `llama3.3-70b-instruct` | 70B | 128k | 71.3 | $0.650 | $0.650 | ~8–38x cheaper |

**Best pick:** `deepseek-v4-pro` — largest model in the catalog (1.6T params), **1M context window** matching Claude's, best OSS reasoning available. Still 3–9x cheaper than Opus.

---

## Routing Decision Guide

| Task Type | Route To | Slug | Reason |
|---|---|---|---|
| Summarize / classify / extract | DO — Deepseek V4 Flash | `deepseek-4-flash` | Haiku-tier quality, ~$0.15/1M blended |
| General chat, drafting, analysis | DO — Deepseek 3.2 | `deepseek-3.2` | MMLU higher than Sonnet 4.6, 7–11x cheaper |
| Code generation / review | DO — Qwen3 Coder Flash | `qwen3-coder-flash` | Specialized coder, 262k context |
| Long context docs (100k–256k) | DO — Kimi K2.5 | `kimi-k2.5` | 256k context, 1T MoE |
| Very long context (256k–1M) | DO — Deepseek V4 Pro | `deepseek-v4-pro` | 1M context window |
| Complex reasoning / agentic | DO — Deepseek V4 Pro | `deepseek-v4-pro` | Largest OSS model |
| Vision / multimodal | DO — Nemotron-nano 12b vl | `nemotron-nano-12b-v2-vl` | Only vision-capable accessible model |
| Quality-critical / production | Bedrock — Sonnet 4.6 | `us.anthropic.claude-sonnet-4-6` | Instruction following, safety, nuance |

---

## Using DO Models in Hermes

### One-off via CLI flag
```bash
hermes chat --provider digitalocean --model deepseek-3.2 "your prompt"
```

### Set as default in config.yaml
```yaml
model:
  default: deepseek-3.2
  provider: digitalocean
```

### Provider config (already in your config.yaml)
```yaml
providers:
  digitalocean:
    name: DigitalOcean Gradient
    api: https://inference.do-ai.run/v1
    key_env: DO_INFERENCE_API_KEY
    transport: openai_chat
```

### Auth
- `DO_INFERENCE_API_KEY` is already set in `~/.hermes/.env`
- Uses your standard DO personal access token — same one doctl uses

---

## The Honest Tradeoff

DO open-source models match or beat Claude on raw MMLU benchmarks in some cases (`deepseek-3.2` at 86.2 vs Sonnet 4.6 at 84.2) but benchmarks are one axis.

**Claude wins on:**
- Instruction following precision and consistency
- Complex multi-step agentic tasks
- Safety guardrails and refusal behavior
- Nuanced tone and formatting control

**DO models win on:**
- Cost (3–43x cheaper depending on tier)
- Context windows (`deepseek-v4-pro` = 1M tokens)
- Volume workloads where rough edges are acceptable
- Routing non-critical tasks away from Bedrock spend

**Bedrock stays primary.** DO is the cost-saving secondary for appropriate workloads.
