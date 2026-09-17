# AI Model API Costs

> Extracted from [awesome-openclaw](https://github.com/rohitg00/awesome-openclaw) — March 2026

The **real cost** of running OpenClaw is the AI model API, not infrastructure.

## Provider Pricing (per 1M tokens)

| Provider | Model | Input Cost | Output Cost | Notes |
|----------|-------|-----------|-------------|-------|
| **Anthropic** | Claude 3.5 Haiku | $0.80 | $4.00 | Best budget option |
| **Anthropic** | Claude 3.5 Sonnet | $3.00 | $15.00 | Best balance |
| **Anthropic** | Claude Opus 4.5 | $15.00 | $75.00 | Most capable, expensive |
| **OpenAI** | GPT-4o | $2.50 | $10.00 | Good alternative |
| **OpenAI** | GPT-4o-mini | $0.15 | $0.60 | Very cheap |
| **Google** | Gemini 2.0 Flash | $0.10 | $0.40 | Cheapest hosted |
| **Google** | Gemini Flash-Lite | **Free tier** | **Free tier** | Zero cost option |
| **DeepSeek** | DeepSeek V3 | $0.27 | $1.10 | Best value reasoning |
| **Moonshot** | Kimi K2.5 | $3.00 | $3.00 | Great agentic performance |
| **Grok** | Grok 4.1 mini | $0.20 | $0.50 | Budget alternative |
| **OpenRouter** | Various | Varies | Varies | Unified API for 200+ models |
| **Ollama** | Any | **$0** | **$0** | Requires hardware (16 GB+ RAM) |

> **75x price difference** between the most expensive (Claude Opus) and cheapest (Gemini Flash) options.

## Monthly API Cost Estimates

| Usage Level | Description | Monthly Cost |
|-------------|-------------|-------------|
| **Minimal** | Few messages/day, Gemini free tier | $0-5 |
| **Light** | ~50 messages/day, Claude Haiku | $5-15 |
| **Moderate** | ~200 messages/day, mixed models | $15-50 |
| **Heavy** | 500+ messages/day, Claude Sonnet | $50-150 |
| **Power User** | All-day usage, Claude Opus | $150-500+ |

## Cost Optimization Tips

1. **Smart model routing** — cheap models (Haiku/Flash) for simple tasks, expensive (Sonnet/Opus) for complex
2. **Prompt caching** — Anthropic auto-applies 5-min cache, reducing repeat costs
3. **Local LLMs** via Ollama for zero API cost (requires 16 GB+ RAM)
4. **Budget alerts** in your AI provider dashboard
5. [OpenClaw Cost Calculator](https://calculator.vlvt.sh/) to estimate spend
6. [ClawRouter](https://github.com/BlockRunAI/ClawRouter) — Smart LLM router (save 78% on costs)

## Total Real-World Cost Examples

| Setup | Infrastructure | API | Total/month | Notes |
|-------|---------------|-----|-------------|-------|
| **Free Tier Max** | Oracle Cloud ($0) | Gemini free ($0) | **$0** | Limited but functional |
| **Ultra Budget** | LumaDock ($2) | Haiku ($5) | **$7** | Basic personal assistant |
| **Sweet Spot** | Hetzner CX22 ($4) | Mixed models ($15) | **$19** | Best value-to-capability |
| **Standard** | DigitalOcean ($12) | Sonnet ($40) | **$52** | Production-ready |
| **Power User** | DigitalOcean ($24) | Opus ($200) | **$224** | Heavy professional use |

## Monthly Budget Targets

| Budget | Strategy |
|--------|----------|
| **$0** | Oracle Cloud + Gemini free tier + Ollama |
| **$6** | Agent37 ($1) + Claude Haiku ($5) |
| **$19** | Hetzner ($4) + mixed models ($15) |
| **$50** | DigitalOcean ($12) + Claude Sonnet ($38) |
| **$100** | Managed hosting ($24) + heavy API ($76) |

## Cost Resources

- [OpenClaw Cost Calculator](https://calculator.vlvt.sh/)
- [API Usage and Costs Docs](https://docs.openclaw.ai/reference/api-usage-costs)
- [Deploy Cost Guide: $0-8/month](https://yu-wenhao.com/en/blog/2026-02-01-openclaw-deploy-cost-guide/)
- [Cost Governor Skill](https://github.com/AtlasPA/openclaw-cost-governor) — Track LLM costs in real-time, enforce budget limits


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
