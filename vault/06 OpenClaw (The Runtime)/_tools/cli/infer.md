---
name: OpenClaw CLI — Infer
updated: 2026-05-11
---

# `openclaw infer`

Unified interface for provider-backed inference tasks — text, images, audio, TTS, video, web, and embeddings.

**Official docs:** [docs.openclaw.ai/cli/infer](https://docs.openclaw.ai/cli/infer.md)

---

## Command Families

```bash
# Text inference
openclaw infer model run --prompt "Summarize this" --json
openclaw infer model run --prompt "..." --gateway   # Use full gateway routing

# Image generation
openclaw infer image generate --prompt "friendly lobster illustration" --json

# Image editing
openclaw infer image edit --image path/to/img.png --prompt "Add sunglasses" --json

# Image description
openclaw infer image describe --image path/to/img.png --json

# Audio transcription
openclaw infer audio transcribe --file path/to/audio.mp3 --json

# Text-to-speech
openclaw infer tts speak --text "Hello world" --output speech.mp3

# Video generation
openclaw infer video generate --prompt "..." --json

# Video description
openclaw infer video describe --file path/to/video.mp4 --json

# Web search
openclaw infer web search --query "latest AI news" --json

# Web fetch
openclaw infer web fetch --url https://example.com --json

# Embeddings
openclaw infer embedding create --text "..." --json
```

---

## Flags

| Flag | Description |
|---|---|
| `--json` | Structured JSON output (for scripting) |
| `--gateway` | Route through gateway for full provider selection |
| `--prompt` | Input prompt |
| `--model` | Override model |

---

## Key Advantage

Uses the providers and models already configured in OpenClaw — no separate API key wiring. Stateless commands run locally without a running gateway; use `--gateway` for full routing.

---

## Example: One-shot image gen from CLI
```bash
openclaw infer image generate \
  --prompt "One Mind OS logo, dark theme, lightning bolt" \
  --json
```
