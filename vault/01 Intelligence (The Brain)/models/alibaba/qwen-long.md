# qwen-long

| Attribute | Value |
|-----------|-------|
| **Model ID** | `qwen-long` |
| **Context** | **1,000,000 tokens** |
| **Input** | $0.05/1M |
| **Output** | $0.20/1M |
| **Tier** | STACK |

**Best for:** Ultra-long context at lowest cost — full vault synthesis, entire books, large codebases.

## Why This Model Matters

1M context at $0.05/1M input — vs Gemini 1.5 Pro at $1.25/1M. Processing 500K tokens costs ~$0.025.

## Cost Comparison (500K token context)

| Provider | Model | Cost |
|----------|-------|------|
| **Qwen-Long** | qwen-long | **~$0.025** |
| Google | gemini-1.5-pro | ~$0.625 |
| Anthropic | claude-3.7-sonnet | needs chunking (200K limit) |

## Use Cases

- Full codex domain as single context window
- Book summarization in one call
- Multi-session agent conversation history
- Large codebase review
- Research corpus synthesis

> Note: Latency from US is 200–500ms (servers in China) — use async/batch jobs, not real-time turns.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
