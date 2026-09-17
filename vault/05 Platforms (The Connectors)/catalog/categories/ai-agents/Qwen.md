# Qwen

> **Alibaba's open-source LLM family — local and cloud inference**

## Overview

Qwen (通义千问) is Alibaba's family of open-source large language models. Available in sizes from 0.5B to 72B+. Can run locally via Ollama or via Alibaba Cloud API. Strong multilingual capabilities, especially Chinese. Viable cost-free alternative for local inference workloads.

## Status

| Attribute   | Value                              |
| ----------- | ---------------------------------- |
| Status      | 🔄 Evaluating                      |
| Category    | AI / LLM                           |
| Type        | Open-Source LLM                    |
| Integration | Ollama, HuggingFace, Alibaba API   |

## Use Cases in OneMind

1. **Local Inference** — Run on VPS via Ollama for cost-free agent tasks
2. **Cost Optimization** — Replace API calls for routine/lower-stakes tasks
3. **Multilingual** — Chinese-language processing for GE/LE international content
4. **Offline AI** — Part of offline AI stack (Ollama + ChromaDB on VPS)

## Key Features

- Qwen2.5 family: 0.5B–72B parameters
- Apache 2.0 license (commercial friendly)
- Strong code (Qwen2.5-Coder) and math models
- Function calling support
- 128K context (larger models)

## Quick Start

```bash
# Run locally via Ollama
ollama pull qwen2.5:7b
ollama run qwen2.5:7b

# Or via API
ollama pull qwen2.5-coder:7b
```

## Links

- HuggingFace: [huggingface.co/Qwen](https://huggingface.co/Qwen)
- Ollama: [ollama.com/library/qwen2.5](https://ollama.com/library/qwen2.5)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
