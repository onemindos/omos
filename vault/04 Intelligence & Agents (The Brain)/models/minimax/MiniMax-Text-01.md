# MiniMax-Text-01

| Attribute | Value |
|-----------|-------|
| **Model ID** | `MiniMax-Text-01` |
| **Context** | **4,000,000 tokens** |
| **Input** | $0.20/1M |
| **Output** | $1.10/1M |
| **Tier** | STACK |

**Best for:** Ultra-long context — the largest context window in the stack. 4M tokens = ~3M words = ~6,000 pages = effectively the entire codex vault in one window.

## Context Comparison

| Model | Max Context |
|-------|------------|
| **MiniMax-Text-01** | **4,000,000** |
| Gemini 1.5 Pro | 2,000,000 |
| Qwen-Long | 1,000,000 |
| Claude 3.7 Sonnet | 200,000 |
| GPT-4o | 128,000 |

## Use Cases

- Full vault synthesis (entire codex domain in one call)
- Multiple books analyzed simultaneously
- Extended multi-session agent history
- Large codebase review with full context

> Note: Effective retrieval degrades in very long contexts — still performs well but precision decreases beyond ~500K tokens. Servers in China, 300–700ms from US — use async/batch.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
