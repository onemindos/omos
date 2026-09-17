# OneMind OS - Model Providers Reference

> Complete guide to all Agno native model providers and when to use each

**Last Updated:** 2026-01-27

---

## Overview

OneMind OS supports multiple model providers through the Agno framework. This document covers:

1. **Native Providers** - Direct API integration with full features
2. **Cloud Providers** - AWS, Azure, Google Cloud hosted models
3. **Local Providers** - Self-hosted models (Ollama, vLLM)
4. **Gateways** - Aggregators like Groq, Together, OpenRouter
5. **LiteLLM** - Universal proxy (convenience, but fewer native features)

### Key Insight: Native vs LiteLLM

| Approach             | Pros                                                            | Cons                                                                  |
| -------------------- | --------------------------------------------------------------- | --------------------------------------------------------------------- |
| **Native Providers** | Full features, provider-specific capabilities, better streaming | More imports to manage                                                |
| **LiteLLM Proxy**    | One API for all, easy switching                                 | Loses provider-specific features (grounding, extended thinking, etc.) |

**Recommendation:** Use native providers for production agents, LiteLLM for experimentation/fallback.

---

## 1. Native Model Providers (Direct API Integration)

### Anthropic Claude

| Import                         | Models                      | Notes                                |
| ------------------------------ | --------------------------- | ------------------------------------ |
| `agno.models.anthropic.Claude` | `claude-opus-4-1-20250805`  | Best model                           |
|                                | `claude-sonnet-4-20250514`  | Good for most use-cases, image input |
|                                | `claude-3-5-haiku-20241022` | Fastest model                        |

**Native Features:**

- Prompt caching
- Structured outputs
- Image input
- Extended thinking (reasoning tokens)

```python
from agno.models.anthropic import Claude

agent = Agent(
    model=Claude(id="claude-sonnet-4-20250514"),
    tools=[...],
)
```

---

### OpenAI

| Import                               | Models                     | Notes                   |
| ------------------------------------ | -------------------------- | ----------------------- |
| `agno.models.openai.OpenAIChat`      | `gpt-4o`, `gpt-4o-mini`    | Legacy chat completions |
| `agno.models.openai.OpenAIResponses` | `gpt-5.2`, `o1`, `o3-mini` | New Responses API       |

**Native Features:**

- Streaming
- Tool calling
- Audio input/output
- Vision (image understanding)
- Structured outputs

```python
from agno.models.openai import OpenAIResponses

agent = Agent(
    model=OpenAIResponses(id="gpt-4o"),
    tools=[...],
)
```

---

### Google Gemini

| Import                      | Models                  | Notes             |
| --------------------------- | ----------------------- | ----------------- |
| `agno.models.google.Gemini` | `gemini-2.0-flash`      | Fast, recommended |
|                             | `gemini-2.0-flash-lite` | Ultra-fast        |
|                             | `gemini-2.5-pro`        | Advanced          |
|                             | `gemini-3-pro-preview`  | Latest            |

**Native Features (ONLY available with native provider):**

- ⭐ **Grounding with Google Search** - Real-time web data
- Image generation
- Multimodal (image/video/audio/PDF input)
- Speech generation
- File search (native RAG)
- Thinking models (advanced reasoning)
- Context caching

```python
from agno.models.google import Gemini

agent = Agent(
    model=Gemini(
        id="gemini-2.0-flash",
        grounding=True,  # Native feature - Google Search grounding
    ),
    tools=[...],
)
```

**Env var:** `GOOGLE_API_KEY`

---

### Mistral

| Import                        | Models                 | Notes    |
| ----------------------------- | ---------------------- | -------- |
| `agno.models.mistral.Mistral` | `mistral-large-latest` | Flagship |
|                               | `mistral-small-latest` | Fast     |

```python
from agno.models.mistral import Mistral

agent = Agent(model=Mistral(id="mistral-large-latest"))
```

**Env var:** `MISTRAL_API_KEY`

---

### DeepSeek

| Import                          | Models              | Notes              |
| ------------------------------- | ------------------- | ------------------ |
| `agno.models.deepseek.DeepSeek` | `deepseek-chat`     | General purpose    |
|                                 | `deepseek-reasoner` | Advanced reasoning |
|                                 | `deepseek-coder`    | Code generation    |

**Env var:** `DEEPSEEK_API_KEY`

---

### Cohere

| Import                      | Models                   | Notes        |
| --------------------------- | ------------------------ | ------------ |
| `agno.models.cohere.Cohere` | `command-r-plus-08-2024` | Best for RAG |
|                             | `command-r-08-2024`      | Balanced RAG |

**Native Features:**

- Optimized for RAG/retrieval
- Citation generation

**Env var:** `COHERE_API_KEY`

---

### xAI (Grok)

| Import                | Models      | Notes     |
| --------------------- | ----------- | --------- |
| `agno.models.xai.xAI` | Grok models | Elon's AI |

**Env var:** `XAI_API_KEY`

---

### Perplexity

| Import                              | Models            | Notes             |
| ----------------------------------- | ----------------- | ----------------- |
| `agno.models.perplexity.Perplexity` | Perplexity models | Search-focused AI |

**Native Features:**

- Built-in web search
- Citation generation

**Env var:** `PERPLEXITY_API_KEY`

---

### DashScope (Alibaba Qwen)

| Import                            | Models      | Notes      |
| --------------------------------- | ----------- | ---------- |
| `agno.models.dashscope.DashScope` | Qwen models | Chinese AI |

**Env var:** `DASHSCOPE_API_KEY`

---

## 2. Cloud Model Providers

### AWS Bedrock

| Import                       | Models                            | Notes          |
| ---------------------------- | --------------------------------- | -------------- |
| `agno.models.aws.AwsBedrock` | `mistral.mistral-large-2402-v1:0` | Mistral on AWS |
|                              | `amazon.nova-pro-v1:0`            | Amazon Nova    |
|                              | All Bedrock foundation models     |                |

**Authentication:** AWS credentials (access key + secret key, SSO, or boto3 session)

```python
from agno.models.aws import AwsBedrock

agent = Agent(
    model=AwsBedrock(
        id="amazon.nova-pro-v1:0",
        aws_region="us-east-1",
    ),
)
```

---

### Claude via AWS Bedrock (Current OneMind Setup ✅)

| Import                       | Models                                       | Notes         |
| ---------------------------- | -------------------------------------------- | ------------- |
| `agno.models.aws.AwsBedrock` | `us.anthropic.claude-opus-4-5-20251101-v1:0` | Claude Opus   |
|                              | `anthropic.claude-3-5-sonnet-20241022-v2:0`  | Claude Sonnet |
|                              | `anthropic.claude-3-5-haiku-20241022-v1:0`   | Claude Haiku  |

```python
from agno.models.aws import AwsBedrock

# This is what Legacy uses
agent = Agent(
    model=AwsBedrock(
        id="us.anthropic.claude-opus-4-5-20251101-v1:0",
        aws_region="us-east-1",
    ),
)
```

**Env vars:** `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`

---

### Azure OpenAI

| Import                          | Models                  | Notes      |
| ------------------------------- | ----------------------- | ---------- |
| `agno.models.azure.AzureOpenAI` | `gpt-4o`, `gpt-4o-mini` | Enterprise |

**Native Features:**

- Prompt caching
- Content filtering (enterprise compliance)
- Azure AD authentication
- Private endpoints

```python
from agno.models.azure import AzureOpenAI

agent = Agent(
    model=AzureOpenAI(
        id="gpt-4o",
        api_base="https://your-resource.openai.azure.com",
        api_key="your-azure-key",
        api_version="2024-10-21",
    ),
)
```

**Env vars:** `AZURE_API_KEY`, `AZURE_API_BASE`, `AZURE_API_VERSION`

---

### Azure AI Foundry

| Import                             | Models          | Notes                   |
| ---------------------------------- | --------------- | ----------------------- |
| `agno.models.azure.AzureAIFoundry` | Azure AI models | Microsoft's AI platform |

---

### Vertex AI Claude

| Import                                | Models                  | Notes      |
| ------------------------------------- | ----------------------- | ---------- |
| `agno.models.vertexai.VertexAIClaude` | Claude via Google Cloud | Enterprise |

---

### IBM WatsonX

| Import                    | Models                | Notes      |
| ------------------------- | --------------------- | ---------- |
| `agno.models.ibm.WatsonX` | IBM foundation models | Enterprise |

---

## 3. Local Model Providers

### Ollama

| Import                      | Models                        | Notes           |
| --------------------------- | ----------------------------- | --------------- |
| `agno.models.ollama.Ollama` | `llama3.2`, `phi3`, `qwen2.5` | Local inference |
|                             | Any Ollama model              | Self-hosted     |

```python
from agno.models.ollama import Ollama

agent = Agent(
    model=Ollama(id="llama3.2"),
)
```

**Env var:** `OLLAMA_URL` (default: `http://localhost:11434`)

---

### vLLM

| Import                  | Models                | Notes                      |
| ----------------------- | --------------------- | -------------------------- |
| `agno.models.vllm.VLLM` | Any HuggingFace model | High-throughput production |

---

### LlamaCpp

| Import                          | Models      | Notes         |
| ------------------------------- | ----------- | ------------- |
| `agno.models.llamacpp.LlamaCpp` | GGUF models | CPU/GPU local |

---

### LM Studio

| Import                          | Models           | Notes       |
| ------------------------------- | ---------------- | ----------- |
| `agno.models.lmstudio.LMStudio` | LM Studio models | Desktop app |

---

## 4. Model Gateways & Aggregators

### Groq (Ultra-Fast Inference)

| Import                  | Models                         | Notes                    |
| ----------------------- | ------------------------------ | ------------------------ |
| `agno.models.groq.Groq` | `llama-3.3-70b-versatile`      | General use, recommended |
|                         | `llama-3.1-8b-instant`         | Fastest                  |
|                         | `llama-3.2-90b-vision-preview` | Vision support           |
|                         | `mixtral-8x7b-32768`           | Open source              |

**What is Groq?**

Groq is an AI inference company with custom hardware (LPU - Language Processing Unit) that runs LLMs **10-100x faster** than traditional GPUs.

| Traditional AI               | Groq LPU                |
| ---------------------------- | ----------------------- |
| ~50-100 tokens/sec           | **500-800+ tokens/sec** |
| ~8-12 seconds for 500 tokens | **~0.5-1 second**       |

**Native Features:**

- Multimodal (vision)
- Ultra-low latency
- Predictable performance

**Pricing (Very Cheap):**
| Model | Input | Output |
|-------|-------|--------|
| Llama 3.3 70B | $0.59/M tokens | $0.79/M tokens |
| Llama 3.1 8B | $0.05/M tokens | $0.08/M tokens |

Compare to Claude Opus at $15/M input, $75/M output - Groq is **20-100x cheaper**.

```python
from agno.models.groq import Groq

agent = Agent(
    model=Groq(id="llama-3.3-70b-versatile"),
)
```

**Env var:** `GROQ_API_KEY`

**Use Cases:**

- Real-time chat (instant responses)
- Voice/audio apps (low latency critical)
- High-volume tasks (cost-effective)
- Quick triage/routing agents

---

### Together AI

| Import                          | Models                                          | Notes              |
| ------------------------------- | ----------------------------------------------- | ------------------ |
| `agno.models.together.Together` | `meta-llama/Meta-Llama-3.1-405B-Instruct-Turbo` | Largest open model |
|                                 | `mistralai/Mixtral-8x22B-Instruct-v0.1`         | Open source        |

**Env var:** `TOGETHER_API_KEY`

---

### Fireworks

| Import                            | Models         | Notes            |
| --------------------------------- | -------------- | ---------------- |
| `agno.models.fireworks.Fireworks` | Fast inference | Production-ready |

**Env var:** `FIREWORKS_API_KEY`

---

### DeepInfra

| Import                            | Models             | Notes           |
| --------------------------------- | ------------------ | --------------- |
| `agno.models.deepinfra.DeepInfra` | Open source models | Cheap inference |

**Env var:** `DEEPINFRA_API_KEY`

---

### OpenRouter

| Import                              | Models      | Notes                  |
| ----------------------------------- | ----------- | ---------------------- |
| `agno.models.openrouter.OpenRouter` | 100+ models | Multi-provider gateway |

**Env var:** `OPENROUTER_API_KEY`

---

### Hugging Face

| Import                                | Models           | Notes         |
| ------------------------------------- | ---------------- | ------------- |
| `agno.models.huggingface.HuggingFace` | HF Inference API | All HF models |

**Env var:** `HUGGINGFACE_API_KEY`

---

### NVIDIA

| Import                      | Models     | Notes                   |
| --------------------------- | ---------- | ----------------------- |
| `agno.models.nvidia.NVIDIA` | NVIDIA NIM | GPU-optimized inference |

---

### Cerebras

| Import                          | Models         | Notes                  |
| ------------------------------- | -------------- | ---------------------- |
| `agno.models.cerebras.Cerebras` | Wafer-scale AI | Ultra-fast (like Groq) |

---

### SambaNova

| Import                            | Models            | Notes               |
| --------------------------------- | ----------------- | ------------------- |
| `agno.models.sambanova.SambaNova` | SambaNova systems | Enterprise hardware |

---

## 5. LiteLLM (Universal Proxy)

| Import                              | Models            | Notes           |
| ----------------------------------- | ----------------- | --------------- |
| `agno.models.litellm.LiteLLM`       | Any LiteLLM model | Direct LiteLLM  |
| `agno.models.litellm.LiteLLMOpenAI` | OpenAI-compatible | Generic wrapper |

**Current OneMind Usage:** Configured in `infra/configs/litellm/config.yaml`

**When to Use:**

- UI model dropdown convenience
- Cost comparison across providers
- Fallback routing (if Claude fails, try GPT)
- Experimentation with many models

**When NOT to Use:**

- Production agents (use native providers)
- Provider-specific features (grounding, audio, etc.)

---

## Recommended Configuration for OneMind OS

| Use Case                   | Provider              | Why                           |
| -------------------------- | --------------------- | ----------------------------- |
| **Primary Agent (Legacy)** | `AwsBedrock` (Claude) | ✅ Already configured         |
| **Google Integration**     | `Gemini` native       | Grounding, image gen, speech  |
| **Azure Enterprise**       | `AzureOpenAI` native  | Compliance, content filtering |
| **Fast Inference**         | `Groq` native         | 10x faster, cheap             |
| **Local/Edge**             | `Ollama` native       | Self-hosted, offline          |
| **UI Model Selector**      | `LiteLLM`             | Convenience                   |
| **Fallback**               | `LiteLLM`             | When primary fails            |

---

## Environment Variables

```bash
# AWS Bedrock (Current)
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
AWS_REGION=us-east-1

# Google Gemini
GOOGLE_API_KEY=...

# Azure OpenAI
AZURE_API_KEY=...
AZURE_API_BASE=https://your-resource.openai.azure.com
AZURE_API_VERSION=2024-10-21

# Groq (Fast Inference)
GROQ_API_KEY=gsk_...

# OpenAI (Direct)
OPENAI_API_KEY=sk-...

# Anthropic (Direct - not via Bedrock)
ANTHROPIC_API_KEY=sk-ant-...

# Others
MISTRAL_API_KEY=...
COHERE_API_KEY=...
DEEPSEEK_API_KEY=...
TOGETHER_API_KEY=...
OPENROUTER_API_KEY=...

# LiteLLM Proxy
USE_LITELLM=false
LITELLM_URL=http://litellm:4000
LITELLM_MASTER_KEY=sk-onemind-master-key
```

---

## Related Documentation

- [Agno Model Docs](https://docs.agno.com/models/overview)
- [Model Compatibility](https://docs.agno.com/models/compatibility)
- [Model Routing Architecture](./architecture/MODEL-ROUTING.md)
- [LiteLLM Config](../infra/configs/litellm/config.yaml)
