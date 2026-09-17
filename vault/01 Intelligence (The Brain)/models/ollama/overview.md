# Ollama

**Type:** Local model runner — desktop app + CLI + REST API
**API base:** `http://localhost:11434/v1` (OpenAI-compatible since v0.1.24)
**Key var:** None required (local)
**Cost:** Free — runs on your CPU/GPU

## What Ollama Is

CLI + daemon that downloads and runs quantized LLMs (GGUF format) locally. Models pulled from Ollama Hub or HuggingFace. Runs on Mac (Apple Silicon), Linux, Windows.

## Stack Integration

OpenClaw agents can use Ollama as a local model backend — zero API cost, zero data egress, offline capable. Essential for:
- Private vault processing (sensitive data never leaves machine)
- Dev/testing without burning API credits
- VPS inference (Ollama running on the DO droplet)

## Popular Models to Pull

```bash
ollama pull llama3.3          # Meta Llama 3.3 70B
ollama pull deepseek-r1       # DeepSeek R1 (multiple sizes)
ollama pull qwen2.5:72b       # Alibaba Qwen 2.5 72B
ollama pull mistral-nemo      # Mistral NeMo 12B
ollama pull nomic-embed-text  # Embeddings for pgvector
ollama pull phi4              # Microsoft Phi-4 14B
```

## Model Sizing Guide

| RAM/VRAM | Max Model Size |
|----------|---------------|
| 8GB | 7B models |
| 16GB | 13B models |
| 32GB | 34B models |
| 64GB | 70B models (quantized) |
| 128GB+ | 70B full / 123B+ |

## API Usage

```python
# OpenAI-compatible — works with any OpenAI SDK
from openai import OpenAI
client = OpenAI(base_url="http://localhost:11434/v1", api_key="ollama")
response = client.chat.completions.create(
    model="llama3.3",
    messages=[{"role": "user", "content": "Hello"}]
)
```

## Key Commands

```bash
ollama list           # see installed models
ollama ps             # see running models
ollama pull hf.co/org/model  # pull directly from HuggingFace
ollama serve          # start daemon if not already running
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
