# Model Routing in OpenClaw

One of OpenClaw's most powerful (and underused) features is **per-context model routing** — the ability to assign different AI models to different types of work. Not everything needs your best model. Smart routing saves cost, improves speed, and matches capability to task.

---

## The Core Idea

Think of it like a staffing strategy:
- **Opus / Sonnet** = your senior strategist — expensive, deep, reserved for what matters
- **Haiku** = your ops coordinator — fast, cheap, handles routine execution
- **Gemini Flash / Whisper** = your sensory layer — audio, voice, real-time I/O

You configure once. OpenClaw routes automatically.

---

## Contexts That Accept a Model

| Context | What It Is | Recommended |
|---|---|---|
| **Agent default** | Every conversation with an agent | Sonnet 4.6 |
| **Subagents** | Spawned sub-tasks from within an agent | Haiku 4.5 |
| **Cron jobs** | Scheduled tasks (heartbeats, reports, pulses) | Haiku 4.5 or per-job override |
| **Heartbeats** | Health check messages from each agent | Haiku 4.5 |
| **Hooks** | Event-triggered automations (message, file change, etc.) | Haiku 4.5 |
| **Talk (voice I/O)** | Real-time voice conversation | Gemini Flash Native Audio |
| **TTS (text-to-speech)** | Agent messages spoken aloud | Google TTS |
| **STT (speech-to-text)** | Your voice transcribed to text | Whisper (local, base model) |

---

## Current One Mind Config

```json
agents.defaults.model.primary         = github-copilot/claude-sonnet-4.6
agents.defaults.model.fallbacks[0]    = amazon-bedrock/us.anthropic.claude-sonnet-4-6

agents.defaults.subagents.model       = github-copilot/claude-haiku-4.5
agents.defaults.subagents.thinking    = low
agents.defaults.subagents.maxConcurrent = 12

talk.provider                         = google
talk.model                            = gemini-2.5-flash-native-audio-preview-12-2025
talk.voice                            = Kore

messages.tts.provider                 = google
tools.media.audio.models[0]           = whisper base (local CLI)
```

**All 18 cron jobs explicitly set to Haiku 4.5 (applied 2026-05-14):**

| Job | Schedule | Agent | Model |
| --- | --- | --- | --- |
| legacy-pulse | every 30 min | main | Haiku 4.5 |
| legacy-heartbeat | daily 8am | main | Haiku 4.5 |
| legacy-onemind-status | daily 7am | main | Haiku 4.5 |
| legacy-daily-report | daily 7:05pm | main | Haiku 4.5 |
| legacy-weekly-review | weekly Sun 9am | main | Haiku 4.5 |
| zenith-heartbeat | daily 8:05am | zenith | Haiku 4.5 |
| zenith-daily-report | daily 7:15pm | zenith | Haiku 4.5 |
| zenith-weekly-scorecard | weekly Mon 8am | zenith | Haiku 4.5 |
| unity-heartbeat | daily 8:10am | unity | Haiku 4.5 |
| unity-weekly-report | weekly Sun 10:30am | unity | Haiku 4.5 |
| eternal-heartbeat | daily 8:15am | eternal | Haiku 4.5 |
| eternal-daily-report | daily 7:20pm | eternal | Haiku 4.5 |
| eternal-weekly-scorecard | weekly Mon 9am | eternal | Haiku 4.5 |
| sovereignty-heartbeat | daily 8:20am | sovereignty | Haiku 4.5 |
| sovereignty-infra-check | daily 3am | sovereignty | Haiku 4.5 |
| sovereignty-daily-report | daily 7:10pm | sovereignty | Haiku 4.5 |
| sentinel-heartbeat | every 6h | sentinel | Haiku 4.5 |
| sentinel-daily-report | daily 7pm | sentinel | Haiku 4.5 |

---

## How to Set Models

### Agent default model
```bash
openclaw config set agents.defaults.model.primary "github-copilot/claude-sonnet-4.6"
openclaw config set agents.defaults.model.fallbacks '["amazon-bedrock/us.anthropic.claude-sonnet-4-6"]'
```

### Subagent model
```bash
openclaw config set agents.defaults.subagents.model "github-copilot/claude-haiku-4.5"
```

### Per-cron-job model override
When creating a cron job, pass `--model`:
```bash
openclaw cron add \
  --agent main \
  --schedule "0 8 * * *" \
  --model "github-copilot/claude-haiku-4.5" \
  --prompt "Send morning heartbeat to Telegram"
```

### Talk (voice) model
```bash
openclaw config set talk.provider "google"
openclaw config set talk.model "gemini-2.5-flash-native-audio-preview-12-2025"
openclaw config set talk.voice "Kore"
# Then restart:
systemctl --user restart openclaw-gateway
```

---

## Available Models (One Mind)

| Model | Alias | Use For |
|---|---|---|
| `github-copilot/claude-sonnet-4.6` | Sonnet 4.6 [gh] | Agent conversations, complex reasoning |
| `github-copilot/claude-haiku-4.5` | Haiku 4.5 [gh] | Subagents, cron, heartbeats, hooks |
| `amazon-bedrock/us.anthropic.claude-sonnet-4-6` | Sonnet 4.6 [bedrock] | Fallback for Sonnet |
| `amazon-bedrock/us.anthropic.claude-haiku-4-5-20251001-v1:0` | Haiku 4.5 [bedrock] | Fallback for Haiku |
| `amazon-bedrock/us.anthropic.claude-opus-4-6` | Opus 4.6 [bedrock] | Deep research, high-stakes reports |
| `gemini-2.5-flash-native-audio-preview-12-2025` | Gemini Flash Audio | Voice/talk mode |

---

## Recommended Routing (Live on One Mind as of 2026-05-14)

```
Interactive conversations  →  Sonnet 4.6     (best quality, worth the cost)
Subagents                  →  Haiku 4.5      (speed + cost)
Cron / scheduled tasks     →  Haiku 4.5      (applied to all 18 jobs)
Heartbeats                 →  Haiku 4.5      (template outputs, zero reasoning needed)
Hooks / event triggers     →  Haiku 4.5      (reactive, fast, cheap)
Voice I/O                  →  Gemini Flash   (real-time audio requirement)
STT                        →  Whisper base   (local, free, good enough)
```

This setup uses your expensive model only where it counts — in the moments you're actually talking to your agents.

### How to apply to your own instance

```bash
# Set subagent model globally
openclaw config set agents.defaults.subagents.model "github-copilot/claude-haiku-4.5"

# Patch all existing cron jobs to Haiku in one command
openclaw cron list --json | python3 -c "
import json, sys, subprocess
jobs = json.load(sys.stdin)
jobs = jobs if isinstance(jobs, list) else jobs.get('jobs', jobs.get('results', []))
for job in jobs:
    subprocess.run(['openclaw', 'cron', 'edit', job['id'], '--model', 'github-copilot/claude-haiku-4.5'])
    print('patched:', job.get('name'))
"

# Or set model on a single job
openclaw cron edit <job-id> --model "github-copilot/claude-haiku-4.5"
```

---

## Key Principle

> Match the model to the cognitive load of the task. A heartbeat message doesn't need Sonnet. A business strategy session does.

The same principle applies to subagents: when an agent spawns 6 child tasks to research something in parallel, those children run on Haiku — fast and cheap — while the parent synthesizes on Sonnet.

---

*Last updated: 2026-05-14*
*Related: [QUICK-REFERENCE.md](QUICK-REFERENCE.md) · [TALK-MODE-CONFIG.md](../../nodes/TALK-MODE-CONFIG.md)*
