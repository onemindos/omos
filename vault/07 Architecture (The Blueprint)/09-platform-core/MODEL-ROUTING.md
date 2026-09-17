# Model Routing Architecture

> Universal model access via LiteLLM or direct AWS Bedrock

## Overview

OneMind OS supports dynamic model selection through two modes:

1. **Direct Bedrock** (default) - Uses AWS Bedrock directly for Claude models
2. **LiteLLM Proxy** (optional) - Universal proxy supporting 100+ LLM providers

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         OneMind UI                               │
│                    (Model Selector Dropdown)                     │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          │ GET /api/models
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Agno AgentOS                                │
│                   /api/models endpoint                           │
│                                                                  │
│   ┌─────────────────┐    ┌─────────────────────────────────┐    │
│   │ USE_LITELLM=    │    │       Model Selection           │    │
│   │ false (default) │───▶│  Return Bedrock fallback models │    │
│   └─────────────────┘    └─────────────────────────────────┘    │
│                                                                  │
│   ┌─────────────────┐    ┌─────────────────────────────────┐    │
│   │ USE_LITELLM=    │───▶│  Fetch from LiteLLM /models     │    │
│   │ true            │    │  Parse and return with provider  │    │
│   └─────────────────┘    └─────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                       LiteLLM Proxy                              │
│                    (models profile)                              │
│                                                                  │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│   │  AWS Bedrock │  │   OpenAI     │  │   Google     │          │
│   │  - Claude    │  │   - GPT-4o   │  │   - Gemini   │          │
│   └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                  │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│   │    Groq      │  │  DeepSeek    │  │   Qwen       │          │
│   │  - Llama 3.1 │  │  - Coder     │  │  - Alibaba   │          │
│   └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `USE_LITELLM` | `false` | Enable LiteLLM proxy mode |
| `LITELLM_URL` | `http://litellm:4000` | LiteLLM service URL |
| `LITELLM_MASTER_KEY` | `sk-onemind-master-key` | LiteLLM API key |
| `DEFAULT_MODEL` | `claude-opus` | Default model when none selected |

### Files

| File | Purpose |
|------|---------|
| [agents/legacy/config.py](../../agents/legacy/config.py) | `ModelConfig` dataclass |
| [agents/legacy/main.py](../../agents/legacy/main.py) | `/api/models` endpoint |
| [infra/compose/cloud/litellm-config.yaml](../../infra/compose/cloud/litellm-config.yaml) | LiteLLM model definitions |
| [ui/src/components/chat/ChatArea/ChatInput/ChatInput.tsx](../../ui/src/components/chat/ChatArea/ChatInput/ChatInput.tsx) | Model selector UI |

## Deployment Modes

### Mode 1: Direct Bedrock (Default)

Uses AWS Bedrock directly. No LiteLLM needed.

```bash
# Just deploy core services
docker compose up -d
```

Models available:
- `claude-opus` - Claude Opus 4.5
- `claude-sonnet` - Claude Sonnet 3.5
- `claude-haiku` - Claude Haiku 3.5

### Mode 2: LiteLLM Proxy

Enables access to 100+ models from multiple providers.

```bash
# 1. Enable LiteLLM in .env
USE_LITELLM=true

# 2. Add provider API keys
OPENAI_API_KEY=sk-...
GEMINI_API_KEY=...
GROQ_API_KEY=...

# 3. Start with models profile
docker compose --profile models up -d
```

## Adding New Models to LiteLLM

Edit `infra/compose/cloud/litellm-config.yaml`:

```yaml
model_list:
  # Add new model
  - model_name: deepseek-coder
    litellm_params:
      model: deepseek/deepseek-coder
    model_info:
      description: "DeepSeek Coder - Fast code generation"
```

Then restart LiteLLM:

```bash
docker compose --profile models restart litellm
```

## API Endpoints

### GET /api/models

Returns available models based on configuration.

**Response (Bedrock mode):**
```json
{
  "models": [
    {"id": "claude-opus", "name": "Claude Opus 4.5", "provider": "anthropic", "source": "bedrock"},
    {"id": "claude-sonnet", "name": "Claude Sonnet 3.5", "provider": "anthropic", "source": "bedrock"},
    {"id": "claude-haiku", "name": "Claude Haiku 3.5", "provider": "anthropic", "source": "bedrock"}
  ],
  "source": "bedrock",
  "default": "claude-opus"
}
```

**Response (LiteLLM mode):**
```json
{
  "models": [
    {"id": "claude-opus", "name": "Claude Opus", "provider": "anthropic", "source": "litellm"},
    {"id": "gpt-4o", "name": "Gpt 4O", "provider": "openai", "source": "litellm"},
    {"id": "gemini-pro", "name": "Gemini Pro", "provider": "google", "source": "litellm"},
    {"id": "llama-groq", "name": "Llama Groq", "provider": "groq", "source": "litellm"}
  ],
  "source": "litellm",
  "default": "claude-opus"
}
```

### GET /api/models/status

Check LiteLLM connection status.

```json
{
  "enabled": true,
  "connected": true,
  "source": "litellm",
  "litellm_url": "http://litellm:4000"
}
```

## UI Integration

The model selector in [ChatInput.tsx](../../ui/src/components/chat/ChatArea/ChatInput/ChatInput.tsx):

1. Fetches models from `/api/models` on mount
2. Shows source badge (bedrock/litellm) in dropdown
3. Displays provider icon for each model
4. Passes selected model to `/api/agents/{id}/runs` endpoint
5. Falls back to hardcoded models if API fails

## Dynamic Model Selection

When a user selects a model in the chat UI:

1. **UI stores selection** in Zustand store
2. **Stream handler** passes `model` field in form data
3. **Backend** creates/caches agent with that model
4. **Agent runs** with the selected model

```
User selects "gpt-4o"
       │
       ▼
┌─────────────────┐     POST /api/agents/{id}/runs
│   ChatInput     │ ──────────────────────────────▶
│  model="gpt-4o" │     FormData: {message, model}
└─────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  get_or_create_  │
                    │  agent(model_id) │
                    └────────┬─────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
              ▼              ▼              ▼
        ┌─────────┐   ┌───────────┐  ┌──────────┐
        │ Bedrock │   │  LiteLLM  │  │  Ollama  │
        │ claude  │   │  gpt-4o   │  │  phi3    │
        └─────────┘   └───────────┘  └──────────┘
```

### POST /api/agents/{agent_id}/runs

Custom endpoint with model selection support.

**Form fields:**
- `message` - User message
- `model` - (optional) Model ID to use
- `stream` - (optional) Stream response
- `session_id` - (optional) Session ID

**Example:**
```bash
curl -X POST http://localhost:8080/api/agents/legacy-mother-ai/runs \
  -F "message=Hello" \
  -F "model=gpt-4o" \
  -F "stream=true"
```

## Supported Providers (via LiteLLM)

| Provider | Models | API Key Env Var | Notes |
|----------|--------|-----------------|-------|
| AWS Bedrock | Claude Opus, Sonnet, Haiku | `AWS_ACCESS_KEY_ID` | Default |
| OpenAI | GPT-4o, GPT-4 Turbo | `OPENAI_API_KEY` | |
| Anthropic | Claude (direct) | `ANTHROPIC_API_KEY` | |
| Google | Gemini Pro, Flash | `GEMINI_API_KEY` | |
| Groq | Llama 3.1, Mixtral | `GROQ_API_KEY` | Ultra-fast |
| DeepSeek | Chat, Coder | `DEEPSEEK_API_KEY` | Very cheap |
| Ollama | Phi-3, Llama 3.2, Qwen | `OLLAMA_URL` | Local/offline |
| Azure | Azure OpenAI | `AZURE_API_KEY` | Enterprise |

## Ollama (Local/Offline Models)

For edge devices or offline operation, Ollama provides lightweight local models:

| Model | Size | VRAM | Use Case |
|-------|------|------|----------|
| `phi3-mini` | ~2GB | 2GB | Edge/laptop |
| `llama-local` | ~1GB | 1GB | IoT/embedded |
| `qwen-local` | ~500MB | 512MB | Ultra-light |

**Setup Ollama:**
```bash
# Run Ollama container
docker run -d -p 11434:11434 --name ollama ollama/ollama

# Pull lightweight models
docker exec ollama ollama pull phi3:mini
docker exec ollama ollama pull llama3.2:1b
docker exec ollama ollama pull qwen2.5:0.5b

# Add to .env
OLLAMA_URL=http://ollama:11434
```

## Troubleshooting

### Models not loading in UI

1. Check API endpoint:
   ```bash
   curl http://localhost:8080/api/models
   ```

2. If LiteLLM enabled, check connection:
   ```bash
   curl http://localhost:8080/api/models/status
   ```

3. Check LiteLLM logs:
   ```bash
   docker logs onemind-litellm
   ```

### LiteLLM not connecting

1. Verify service is running:
   ```bash
   docker ps | grep litellm
   ```

2. Check health:
   ```bash
   curl http://localhost:4000/health
   ```

3. Verify API key:
   ```bash
   curl http://localhost:4000/models \
     -H "Authorization: Bearer sk-onemind-master-key"
   ```

## Resource Requirements

| Mode | RAM | Notes |
|------|-----|-------|
| Direct Bedrock | 0 | No additional service |
| LiteLLM | 100-200MB idle | Scales with concurrent requests |
| LiteLLM (production) | 4GB+ | For high traffic (1000+ RPS) |

## Best Practices

1. **Start with Direct Bedrock** - Simpler, no extra service
2. **Add LiteLLM when needed** - Multiple providers, cost optimization
3. **Cache model list** - UI caches for 5 minutes
4. **Set default model** - Prevents empty selection on first load
5. **Use profiles** - Only run LiteLLM when needed

## Related Documentation

- [LiteLLM Docs](https://docs.litellm.ai/)
- [Full Access Architecture](FULL-ACCESS-ARCHITECTURE.md)
- [Infrastructure Visual](./infrastructure-visual.md)
