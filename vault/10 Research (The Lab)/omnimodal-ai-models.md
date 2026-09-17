---
title: "Omnimodal AI Models — Research"
type: reference
domain: SO
status: active
created: 2026-05-26
updated: 2026-05-26
tags: [ai, models, multimodal, omnimodal, audio, video, vision, research, intelligence]
---

# Omnimodal AI Models — Research

Omnimodal = models that accept AND/OR output across all major modalities:
**Audio · Video · Text · Images · Documents**

---

## Model Comparison Matrix

| Model | Text | Image | Audio In | Video | Docs/PDF | Audio Out | Notes |
|-------|------|-------|----------|-------|----------|-----------|-------|
| **GPT-4o** | ✅ | ✅ | ✅ | ✅ (frames) | ✅ | ✅ | Native audio I/O, real-time voice |
| **GPT-4o mini** | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ | Lighter, cheaper |
| **Gemini 2.0 Flash** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | 1M context, best overall omni |
| **Gemini 2.5 Pro** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Strongest reasoning + multimodal |
| **Gemini Omni** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Google's flagship omni stack |
| **Claude 3.7 Sonnet** | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | No native audio/video — text+vision |
| **Grok 2 Vision** | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ | Vision + text only |
| **Qwen2.5-Omni-7B** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Open source** — full omni, self-hostable |
| **Qwen2.5-Omni-72B** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | **Open source** — larger, stronger |
| **Phi-4 Multimodal** | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ | Microsoft, small/efficient |
| **Gemini Robotics-ER** | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | Spatial + physical reasoning for robots |

---

## Deep Dives

### GPT-4o (OpenAI)
- **Modalities:** Text, image, audio in/out, video (frame-by-frame)
- **Context:** 128K tokens
- **Audio:** Native real-time voice — not STT→LLM→TTS pipeline, actual audio tokens
- **Best for:** Real-time voice agents, live transcription + reasoning, vision tasks
- **API:** OpenAI API, Azure OpenAI (Zeus has both)
- **Limitation:** Video is frames not true streaming video understanding

### Gemini 2.0 Flash (Google)
- **Modalities:** Text, image, audio, video, PDF — all in. Text + audio out.
- **Context:** 1M tokens input, 64K output
- **Knowledge cutoff:** January 2025
- **Best for:** Long-context multimodal (hour-long videos, massive docs), agentic tasks
- **Tools:** Function calling, structured output, search, code execution
- **API:** Gemini API, Google AI Studio
- **Strength:** Best price/performance for omni tasks

### Gemini 2.5 Pro (Google)
- **Modalities:** Full omni — audio, video, text, image, docs
- **Context:** 1M+ tokens
- **Best for:** Complex reasoning over multimodal inputs, research-grade tasks
- **Limitation:** Higher cost than Flash

### Gemini Robotics-ER 1.6 (Google DeepMind)
- **Modalities:** Vision, text, spatial/physical input
- **Specialty:** Embodied reasoning — understands physical spaces, plans multi-step robot tasks
- **Capabilities:** Instrument reading, spatial reasoning, physical reasoning
- **Best for:** Robot fleet (Eden, drone agents) — physical world understanding
- **Note:** This is the model to watch for OneMind robot integration

### Qwen2.5-Omni (Alibaba — Open Source)
- **Modalities:** Text, image, audio in/out, video — end-to-end omni
- **Sizes:** 7B (edge/local) and 72B (server)
- **Audio output:** Natural speech synthesis built-in — not a separate TTS
- **Self-hostable:** Yes — runs on local hardware
- **Best for:** OneMind sovereign deployment — full omni without API dependency
- **7B:** Runs on MBP M2 Max 32GB — fits in memory
- **72B:** Would need VPS or quantized on MBP

### Claude 3.7 Sonnet (Anthropic)
- **Modalities:** Text, image, PDF/docs — NO audio, NO video
- **Best for:** Deep reasoning, coding, document analysis
- **Limitation:** Not omnimodal — text+vision only
- **Note:** Zeus's current primary model — strong but not the right pick for audio/video tasks

### Phi-4 Multimodal (Microsoft)
- **Modalities:** Text, image, audio in — no video, no audio out
- **Size:** Small/efficient — designed for edge deployment
- **Best for:** Edge nodes (Pi, drones) where compute is limited
- **Self-hostable:** Yes

---

## Use Case Routing for OneMind OS

| OneMind Use Case | Best Model | Why |
|-----------------|-----------|-----|
| Voice agent (live conversation) | GPT-4o or Qwen2.5-Omni | Native audio I/O, low latency |
| Eden (garden robot vision) | Gemini 2.0 Flash or Gemini Robotics-ER | Video + spatial reasoning |
| Drone surveillance feed | Gemini 2.0 Flash | Video + real-time analysis |
| Document/research ingestion | Gemini 2.5 Pro | 1M context, PDF native |
| Sensor fusion (audio + vision + data) | Qwen2.5-Omni-72B (self-hosted) | Full omni, sovereign |
| Edge node intelligence (Pi/drone) | Phi-4 or Qwen2.5-Omni-7B | Small, runs locally |
| Physical world robot planning | Gemini Robotics-ER | Built for embodied tasks |
| Coding + reasoning (current) | Claude 3.7 Sonnet | Best-in-class text/code |

---

## Self-Hostable Omni Models (Sovereign Priority)

| Model | Size | Modalities | Fits on MBP M2 Max 32GB |
|-------|------|------------|------------------------|
| Qwen2.5-Omni-7B | 7B | Full omni | ✅ Yes |
| Qwen2.5-Omni-7B Q4 | ~4GB | Full omni | ✅ Easily |
| Qwen2.5-Omni-72B Q4 | ~40GB | Full omni | ⚠️ Tight |
| Phi-4 Multimodal | 14B | Text+image+audio | ✅ Yes |
| LLaVA variants | 7-34B | Text+image | ✅ Yes |

---

## What Makes Omni Models Different

Traditional pipeline: Audio → STT → LLM → TTS → Audio
True omni model: Audio tokens → Model → Audio tokens (end-to-end)

**Why it matters:**
- Lower latency — no pipeline hops
- Richer understanding — model hears tone, emotion, pace — not just words
- Cross-modal reasoning — "what does that sound match in this image?"
- Unified context — all modalities share the same attention space

---

## OneMind Architecture Implications

- **Current:** Claude 3.7 Sonnet (text+vision) as brain — no audio/video
- **Near-term:** Add GPT-4o or Qwen2.5-Omni for voice agent layer (LiveKit)
- **Robot fleet:** Gemini Robotics-ER for physical world reasoning on Eden + drones
- **Sovereign stack:** Qwen2.5-Omni-7B self-hosted on MBP for offline/edge inference
- **Edge nodes:** Phi-4 Multimodal on Pi hardware for local sensor processing

---

## Sources
- https://deepmind.google/technologies/gemini/flash/
- https://ai.google.dev/gemini-api/docs/models
- https://huggingface.co/Qwen/Qwen2.5-Omni-7B
- https://platform.openai.com/docs/models
