---
id: so-01-ai-models-google
title: "AI Models — Google (Gemini)"
type: reference
status: active
created: 2026-04-12
modified: 2026-04-12
owner: zeus
tags: [type/reference, status/active, domain/so, ai/models, provider/google]
ai_summary: "Google Gemini model reference — context windows, pricing, API setup. Used for Eden agent, multimodal tasks, and ultra-long context."
---

# Google (Gemini)

> SECONDARY provider. Gemini powers the Eden agent (LE/sustainability) and handles multimodal + ultra-long context tasks.
> Best provider for: images + audio + video in one call, 1M+ token context windows.

---

## Overview

Google's Gemini family is the go-to choice when multimodal input (text + image + audio + video) is needed in a single call, or when context windows exceed what Anthropic/OpenAI offer. The free tier via AI Studio is generous for experimentation.

- **API Base URL:** `https://generativelanguage.googleapis.com`
- **SDK:** `@google/generative-ai` (Node.js) | `google-generativeai` (Python)
- **Key Location:** 1Password → SO vault → `Google AI Studio API Key`
- **Console:** https://aistudio.google.com
- **Pricing:** https://ai.google.dev/pricing

---

## Models

| Model | Context | Input ($/1M) | Output ($/1M) | Best Use |
|-------|---------|-------------|--------------|---------|
| `gemini-2.0-flash` | 1M | $0.075 | $0.30 | **Default Gemini** — fast, cheap, multimodal |
| `gemini-2.0-flash-lite` | 1M | $0.0188 | $0.075 | Ultra-cheap batch/summary tasks |
| `gemini-2.0-flash-thinking` | 32K | $0.075 | $0.30 | Flash with thinking steps |
| `gemini-1.5-pro` | 2M | $1.25 | $5 | True 2M context, research synthesis |
| `gemini-1.5-flash` | 1M | $0.075 | $0.30 | Archival — gemini-2.0-flash preferred |

> Free tier (AI Studio): 15 requests/min, 1M tokens/min, 1500 req/day for 2.0 Flash.

### Multimodal Input Support

| Capability | gemini-2.0-flash | gemini-1.5-pro |
|-----------|-----------------|----------------|
| Text | ✅ | ✅ |
| Images | ✅ | ✅ |
| Audio | ✅ | ✅ |
| Video | ✅ | ✅ |
| PDF | ✅ | ✅ |
| Code execution | ✅ | ✅ |
| Grounding (web search) | ✅ | ✅ |

---

## API Setup

```bash
npm install @google/generative-ai     # Node.js
pip install google-generativeai       # Python

export GOOGLE_API_KEY="AIza..."
```

**Basic call (Node.js):**
```typescript
import { GoogleGenerativeAI } from '@google/generative-ai';
const genAI = new GoogleGenerativeAI(process.env.GOOGLE_API_KEY!);
const model = genAI.getGenerativeModel({ model: 'gemini-2.0-flash' });
const result = await model.generateContent('Hello');
```

**With image:**
```typescript
const result = await model.generateContent([
  { inlineData: { data: base64ImageData, mimeType: 'image/jpeg' } },
  'Describe this image',
]);
```

**Vertex AI (enterprise/GCP billing):**
```bash
# Alternative endpoint — uses GCP project billing instead of API key
# API: https://{region}-aiplatform.googleapis.com
# Auth: Application Default Credentials (ADC)
```

---

## Usage in One Mind Stack

| Use Case | Model | Notes |
|----------|-------|-------|
| **Eden agent** (LE/environment) | gemini-2.0-flash | Sustainability analysis, environmental data |
| **Long context tasks** | gemini-1.5-pro | Documents >200K tokens, full vault synthesis |
| **Image analysis** | gemini-2.0-flash | Processing photos, screenshots, diagrams |
| **Audio transcription** | gemini-2.0-flash | Alternative to Whisper for long recordings |
| **PDF reading** | gemini-1.5-pro | Upload PDFs directly, no pre-extraction needed |
| **Gemini CLI** (local) | gemini-2.0-flash | Dev-time queries, GEMINI.md in codex |

---

## Notes

- **Context caching:** Available for inputs >32K tokens — major cost savings for repeated large prompts
- **Grounding:** Enable `tools: [{ googleSearch: {} }]` for real-time web answers
- **Thinking:** `gemini-2.0-flash-thinking-exp` adds visible reasoning steps
- **File API:** Upload files up to 2GB — persisted for 48 hours, shareable by URI across requests
- **Function calling:** Fully supported — use for tool-calling in multi-step agents
- **Safety filters:** May need `HarmBlockThreshold.BLOCK_NONE` for certain content analysis tasks
- **GEMINI.md** in codex root contains Zeus-specific Gemini CLI setup


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Models (AI Models)/google (Google)/google|google]]
