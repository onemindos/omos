# GitHub Copilot / GitHub Models

**Type:** AI coding assistant + model inference API (GitHub Models)
**API base (GitHub Models):** `https://models.inference.ai.azure.com`
**Key var:** GitHub Personal Access Token (PAT)
**Free tier:** Yes — GitHub Models is free with personal GitHub account (rate-limited)

## Two Distinct Products

### 1. GitHub Copilot (Coding Assistant)

IDE coding assistant — autocomplete, chat, inline edits, code review.

| Plan | Cost | Notes |
|------|------|-------|
| Individual | $10/mo | IDE integration only |
| Business | $19/mo per seat | Admin controls |
| Enterprise | $39/mo per seat | Fine-tuning, audit logs |

Copilot uses GPT-4o, Claude, and Gemini under the hood depending on context. Not a direct API — outputs are in-IDE only.

### 2. GitHub Models (API)

Marketplace to access top models via GitHub PAT — useful for dev/testing without separate API keys.

Available models include: GPT-4o, GPT-4o-mini, Claude 3.5 Sonnet, Llama 3.3 70B, Mistral Large, Phi-4, DeepSeek-R1, Cohere Embed.

```python
from openai import OpenAI
client = OpenAI(
    base_url="https://models.inference.ai.azure.com",
    api_key="github_PAT_token_here"
)
```

## Stack Use Cases

- **Copilot:** Forge agent / dev tooling — IDE assistant for code-heavy builds
- **GitHub Models:** Dev/testing new models without creating separate accounts
- **Free-tier fallback:** When OpenAI/Anthropic rate limits hit during dev sprints


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
