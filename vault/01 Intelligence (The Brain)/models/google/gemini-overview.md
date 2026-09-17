# Google Gemini

**Type:** Cloud inference API (Google AI Studio + Vertex AI)
**API base (AI Studio):** `https://generativelanguage.googleapis.com/v1beta/openai/` (OpenAI-compatible)
**Key var:** `GOOGLE_API_KEY` (AI Studio) or service account (Vertex AI)
**Free tier:** Yes — Gemini 1.5 Flash and 2.0 Flash have generous free tiers

> Note: Root `Models/google/` has model documentation. This folder covers the API integration layer.

## API Options

| Option | Base URL | Auth | Notes |
|--------|----------|------|-------|
| AI Studio | `generativelanguage.googleapis.com/v1beta/openai` | API key | Simple, free tier, personal |
| Vertex AI | `us-central1-aiplatform.googleapis.com` | GCP service account | Enterprise, private, VPC |

## Models in This Folder

See per-model files for: Gemini 2.5 Pro, Gemini 2.5 Flash, Gemini 2.0 Flash, Gemini 1.5 Flash.

## Stack Use Cases

- Long document analysis (1M+ context on 2.5 Pro)
- Multimodal input (vision, audio, video natively)
- Free-tier dev inference (Flash free tier is very generous)
- Enterprise via Vertex when data must stay in GCP


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Models (AI Models)/google (Google)/google|google]]
