---
title: "Ollama"
type: reference
status: active
tags: [platform, ai, tools, local-inference, llm, self-hosted]
---
# Ollama

> Run LLMs locally. One command to download and serve any open model. No cloud, no API keys.

## CLI
```bash
ollama run llama3.1          # Chat with model
ollama run codellama         # Code completion
ollama pull mistral          # Download model
ollama list                  # Installed models
ollama serve                 # Start API server
ollama rm model-name         # Delete model
```

## API (OpenAI-compatible)
```bash
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.1",
  "prompt": "Hello"
}'

# OpenAI-compatible endpoint
curl http://localhost:11434/v1/chat/completions -d '{
  "model": "llama3.1",
  "messages": [{"role": "user", "content": "Hello"}]
}'
```

## Key Features
| Feature | Notes |
|---------|-------|
| Local inference | No internet needed |
| Model library | Llama, Mistral, Phi, Gemma, Qwen, CodeLlama |
| GPU acceleration | Metal (Mac), CUDA (NVIDIA) |
| OpenAI compatible | Drop-in replacement for API calls |
| Modelfile | Custom model configs (system prompt, params) |
| Memory | Context window up to 128K (model dependent) |

## Integration
- Hermes: use as local provider (`ollama/model-name`)
- Cursor: local model backend
- Open WebUI: ChatGPT-like interface for Ollama
