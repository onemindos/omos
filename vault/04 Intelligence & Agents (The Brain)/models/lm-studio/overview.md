# LM Studio

**Type:** Desktop GUI application for local LLM inference
**API base:** `http://localhost:1234/v1` (OpenAI-compatible)
**Key var:** None required (local)
**Cost:** Free app — runs on Mac/Windows/Linux
**Download:** lmstudio.ai

## What LM Studio Is

GUI app that makes running LLMs locally accessible — browse and download GGUF models from HuggingFace, chat in a built-in UI, or expose an OpenAI-compatible local server.

## vs Ollama

| | Ollama | LM Studio |
|-|--------|----------|
| Interface | CLI / API | GUI + API |
| Automation | Excellent | Limited CLI |
| Model mgmt | CLI commands | Visual browser |
| Performance | Same (same GGUF) | Same |
| Best for | Server / agent use | Manual testing / chat |

## Stack Use Cases

- Manual model testing — try a new GGUF before integrating
- Demo environment — non-technical users can run models
- Quick local chat with sensitive documents (no API calls)
- Cross-reference: test a model in LM Studio, then deploy same GGUF in Ollama for agent use

## API Usage

Same as Ollama — OpenAI-compatible REST on port 1234. Set `base_url="http://localhost:1234/v1"` in any OpenAI SDK client.

## Model Sources

- HuggingFace GGUF models (browse inside app)
- TheBloke, bartowski repos — well-quantized model packs
- Same models as Ollama — pick Q4_K_M for best quality/size tradeoff


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
