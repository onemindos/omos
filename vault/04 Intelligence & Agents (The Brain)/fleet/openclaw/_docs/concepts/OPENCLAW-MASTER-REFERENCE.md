# OpenClaw Master Reference — One Mind OS
*Compiled from: systematic read of all ~250 bundled docs at `/usr/lib/node_modules/openclaw/docs/` on VPS `onemindos-cloud-1` (64.227.8.23), OpenClaw v2026.4.14 (323493f). April 2026.*

> **For Zeus and One Mind students:** This file is the complete institutional memory of everything OpenClaw can do, how it works, and how this specific VPS is configured. Every section references the source doc so you can go deeper. Nothing is made up — everything here was extracted directly from the bundled source docs.

---

## Table of Contents

1. [What OpenClaw Is](#1-what-openclaw-is)
2. [VPS Setup — This Installation](#2-vps-setup--this-installation)
3. [Agent Architecture](#3-agent-architecture)
4. [Multi-Agent Routing](#4-multi-agent-routing)
5. [Memory Systems](#5-memory-systems)
6. [Automation Engine](#6-automation-engine)
7. [Channels](#7-channels)
8. [Node System (Clients)](#8-node-system-clients)
9. [HTTP APIs (OpenAI-Compatible)](#9-http-apis-openai-compatible)
10. [Gateway Configuration](#10-gateway-configuration)
11. [Security & Secrets](#11-security--secrets)
12. [Context Management](#12-context-management)
13. [Plugin System](#13-plugin-system)
14. [CLI Reference](#14-cli-reference)
15. [One Mind Agent Setup](#15-one-mind-agent-setup)
16. [Cron Jobs — One Mind Schedule](#16-cron-jobs--one-mind-schedule)
17. [Config Reference — This VPS](#17-config-reference--this-vps)
18. [Teaching Notes for Students](#18-teaching-notes-for-students)

---

## 1. What OpenClaw Is

**Source:** `README`, runtime observation, docs index

OpenClaw is a self-hosted AI gateway. You run it on your own server (or desktop), connect your AI models, connect your messaging channels (WhatsApp, Telegram, Slack, Discord, etc.), and it becomes your personal AI operating system.

Key mental model:
```
YOUR PHONE / LAPTOP
      ↓ (WhatsApp/Telegram/Slack message)
OPENCLAW GATEWAY (your VPS)
      ↓
AI MODEL (Claude, GPT, Gemini, Ollama local model)
      ↓ response
BACK TO YOUR PHONE
```

What makes it different from just using the ChatGPT app:
- **You own the server.** Your data stays on your machine.
- **Multi-channel.** One AI, reachable from WhatsApp, Telegram, Discord, Slack, SMS, email, etc.
- **Memory.** It remembers across sessions using a real memory system (QMD/BM25+vector).
- **Automation.** Cron jobs, heartbeats, standing orders — your AI runs while you sleep.
- **Multi-agent.** Multiple AI personalities, each with their own workspace, routing, and purpose.
- **Tool use.** Shell exec, file read/write, web search, browser control.

**Version:** 2026.4.14 (323493f)
**Install path:** `/usr/lib/node_modules/openclaw/`
**Docs path:** `/usr/lib/node_modules/openclaw/docs/`

---

## 2. VPS Setup — This Installation

**Source:** Direct VPS inspection + `OPENCLAW-OTEL-DATADOG.md` + session history

### Server Details
| Field | Value |
|-------|-------|
| VPS Provider | DigitalOcean |
| Public IP | `64.227.8.23` |
| Hostname | `onemindos-cloud-1` |
| Tailscale hostname | `omos-core.tail717ba0.ts.net` |
| Tailscale IP | `100.98.25.10` |
| OS | Ubuntu (hardened April 15, 2026) |
| RAM | 31 GB |
| Disk | 387 GB total, ~108 GB used |

### OpenClaw Service
| Field | Value |
|-------|-------|
| Service type | USER systemd service |
| Service name | `openclaw-gateway.service` |
| Port | 18789 |
| Gateway bind | loopback (`127.0.0.1`) |
| External access | Via Tailscale HTTPS serve |
| Dashboard URL | `https://omos-core.tail717ba0.ts.net` |
| Auth mode | Token |
| Gateway token | `8759c8a97ba63a47fd868c8ea01e7c0307b380f8f656c908` |

### Installed Software
| Tool | Version | Path | Purpose |
|------|---------|------|---------|
| OpenClaw | 2026.4.14 | `/usr/lib/node_modules/openclaw/` | AI gateway |
| QMD | 2.1.0 | `/usr/bin/qmd` | BM25+vector memory sidecar |
| ClawHub | v0.9.0 | `/usr/bin/clawhub` | Skills/plugins marketplace |
| Ollama | current | `http://127.0.0.1:11434` | Local model inference |

### Active Plugins (as of April 2026)
```
active-memory      — Blocking memory sub-agent before each reply
memory-core        — Core memory search (BM25+vector)
memory-wiki        — Memory with wiki-style lookup
lobster            — TBD
llm-task           — Background LLM task execution
openshell          — Shell exec capability
webhooks           — Inbound webhook receiver
diffs              — File diff tool
github-copilot     — GitHub Copilot model provider
amazon-bedrock     — AWS Bedrock model provider
microsoft          — Microsoft model provider
ollama             — Local Ollama inference
```

### OTEL / Observability
Telemetry flows: OpenClaw → OTLP (port 4317, loopback) → Datadog Agent → `us5.datadoghq.com`

Drop-in env files for the gateway service:
- `/root/.config/systemd/user/openclaw-gateway.service.d/copilot.conf`
- `/root/.config/systemd/user/openclaw-gateway.service.d/deepgram-env.conf`
- `/root/.config/systemd/user/openclaw-gateway.service.d/ollama-env.conf`

### File Layout
```
/root/.openclaw/
├── openclaw.json          ← Main config (DO NOT hand-edit without backup)
├── openclaw.json.bak*     ← Automatic backups
├── OPENCLAW-MASTER-REFERENCE.md  ← THIS FILE
├── workspace/             ← Default agent workspace (Legacy/main)
│   ├── AGENTS.md          ← Agent operating instructions
│   ├── SOUL.md            ← Agent personality
│   ├── IDENTITY.md        ← Agent name/emoji/description
│   ├── USER.md            ← User context (Zeus)
│   ├── TOOLS.md           ← Local tool notes
│   ├── BOOTSTRAP.md       ← First-run instructions
│   ├── HEARTBEAT.md       ← Periodic check list
│   └── memory/            ← Daily notes + MEMORY.md
├── workspace-spartan/     ← Spartan agent workspace
├── workspace-oracle/      ← Oracle agent workspace
├── workspace-eden/        ← Eden agent workspace
├── workspace-haven/       ← Haven agent workspace
├── workspace-empire/      ← Empire agent workspace
├── workspace-guardian/    ← Guardian agent workspace
├── workspace-grid/        ← Grid agent workspace
├── workspace-trinity/     ← Trinity agent workspace
├── agents/
│   ├── main/              ← Default agent state
│   ├── spartan/           ← Spartan agent state
│   ├── oracle/            ← Oracle agent state
│   └── ...
├── cron/
│   └── jobs.json          ← Persistent cron job definitions
├── hooks/                 ← Custom event hooks
├── qmd/                   ← QMD memory storage
└── logs/                  ← Gateway logs
```

---

## 3. Agent Architecture

**Source:** `concepts/multi-agent.md`, `concepts/soul.md`, `cli/approvals.md`

### What an Agent Is

An **agent** in OpenClaw is a fully isolated "brain":
- **Workspace** — a folder with SOUL.md, AGENTS.md, USER.md, memory files. This is the agent's home directory.
- **State directory** (`agentDir`) — auth profiles, model registry, per-agent config. Located at `~/.openclaw/agents/<agentId>/agent/`.
- **Session store** — chat history under `~/.openclaw/agents/<agentId>/sessions/`.

**Auth is per-agent.** Each agent reads from its own `~/.openclaw/agents/<agentId>/agent/auth-profiles.json`. Main agent credentials are NOT shared automatically.

### Bootstrap Files

Every agent workspace automatically gets injected these files at session start:
```
AGENTS.md    ← Operating rules for this agent
SOUL.md      ← Personality configuration
TOOLS.md     ← Local tool notes
IDENTITY.md  ← Name, emoji, role description
USER.md      ← User context
HEARTBEAT.md ← Heartbeat checklist (when heartbeat enabled)
BOOTSTRAP.md ← First-run setup guide (delete after first run)
MEMORY.md    ← Long-term curated memory (main session only)
```

### SOUL.md — The Personality File

**Source:** `concepts/soul.md`

SOUL.md is where the agent's voice lives. The "Molty prompt" approach: rewrite SOUL.md to give the agent an actual personality, not just instructions.

Key rules:
- Short/sharp beats long/vague
- Personality ≠ permission to be sloppy
- This is how the agent talks, not just what it knows

### Exec Policy (Approvals)

**Source:** `cli/approvals.md`

```bash
openclaw exec-policy preset yolo      # No confirmations (YOLO mode)
openclaw exec-policy preset standard  # Confirm destructive actions
openclaw exec-policy allowlist add /path/to/script.sh  # Allow specific scripts
```

---

## 4. Multi-Agent Routing

**Source:** `concepts/multi-agent.md`

### How it Works

Multiple agents can run in one Gateway process. Each has full isolation. Inbound messages are routed by **bindings** (most-specific wins).

#### Routing Tier Order (most-specific first)
1. `peer` match — exact DM/group/channel ID
2. `parentPeer` match — thread inheritance
3. `guildId + roles` — Discord role routing
4. `guildId` — Discord server
5. `teamId` — Slack workspace
6. `accountId` match — specific channel account
7. `accountId: "*"` — channel-wide fallback
8. Default agent (`main`) — fallback

#### Adding Agents

```bash
# Interactive wizard
openclaw agents add <name>

# Non-interactive (scripted)
openclaw agents add <name> \
  --workspace /root/.openclaw/workspace-<name> \
  --non-interactive

# View all agents and their bindings
openclaw agents list --bindings

# Update display identity
openclaw agents set-identity <name> --name "Display Name" --emoji "⚡"

# Add routing binding
openclaw agents bind <agentId> --channel whatsapp

# List all bindings
openclaw agents bindings

# Delete agent
openclaw agents delete <name>
```

#### Config Structure (in openclaw.json)
```json5
{
  agents: {
    defaults: {           // ← shared defaults for all agents
      workspace: "/root/.openclaw/workspace",
      model: { primary: "...", fallbacks: ["..."] }
    },
    list: [
      { id: "main" },     // ← default agent (Legacy)
      { id: "spartan", workspace: "/root/.openclaw/workspace-spartan" },
      { id: "oracle",   workspace: "/root/.openclaw/workspace-oracle" },
      // ... etc
    ]
  },
  bindings: [
    // Routes go here after channels are connected
    { agentId: "main", match: { channel: "telegram", accountId: "default" } }
  ]
}
```

### Cross-Agent Memory Search

```json5
{
  agents: {
    list: [
      {
        id: "main",
        memorySearch: {
          qmd: {
            extraCollections: [
              { path: "/root/.openclaw/agents/empire/sessions", name: "empire-sessions" }
            ]
          }
        }
      }
    ]
  }
}
```

---

## 5. Memory Systems

**Source:** `concepts/memory-search.md`, `concepts/active-memory.md`, `concepts/dreaming.md`, `concepts/memory-honcho.md`, `concepts/session-pruning.md`

### Overview

OpenClaw has 4 distinct memory layers:

| Layer | What | Scope |
|-------|------|-------|
| **Active Memory** | Blocking sub-agent retrieves memories before each reply | Per-session |
| **QMD / Memory Search** | Hybrid BM25+vector search over session transcripts + workspace files | Cross-session |
| **Dreaming** | Background consolidation → MEMORY.md | Long-term |
| **Honcho** | Cross-session AI-native user modeling (external, self-hostable) | Long-term |

### Memory Search (QMD)

**Source:** `concepts/memory-search.md`

QMD is the memory backend. It runs as a sidecar process (`/usr/bin/qmd`) alongside OpenClaw.

**How it works:**
1. Session transcripts are indexed by QMD automatically
2. When a user sends a message, OpenClaw queries QMD with a hybrid search (BM25 keyword + vector semantic)
3. Relevant memories are injected into context

**Embedding providers (auto-detected):**
- OpenAI / Gemini / Voyage / Mistral → auto-detect from auth
- Bedrock → auto-detect via AWS env
- Ollama → **must set explicitly** in config
- Local GGUF → ~0.6GB model, no API required

**Temporal decay:** Memories older than the 30-day half-life get down-ranked. "Evergreen" files (pinned workspace docs) are exempt.

**MMR diversity:** Ensures search results aren't all about the same topic.

**Multimodal:** Gemini Embedding 2 supports image embeddings.

**Session memory opt-in:** Session transcripts are NOT indexed by default — opt-in per session.

#### Config
```json
{
  "memory": {
    "backend": "qmd",
    "qmd": {
      "sessions": {
        "enabled": true
      }
    }
  }
}
```

### Active Memory Plugin

**Source:** `concepts/active-memory.md`

Active Memory runs a **blocking sub-agent** before each reply. This sub-agent retrieves specific memories relevant to the current message.

- Only runs in **persistent interactive sessions** — NOT in heartbeat, automation, or headless
- Query modes: `message` (default) / `recent` / `full`
- Prompt styles: `balanced` / `strict` / `contextual` / `recall-heavy` / `precision-heavy` / `preference-only`

**Debug commands:**
```
/verbose on    ← shows memory status line
/trace on      ← shows full debug summary
```

**Config:**
```json
{
  "plugins": {
    "entries": {
      "active-memory": {
        "enabled": true,
        "config": {
          "queryMode": "message",
          "promptStyle": "balanced"
        }
      }
    }
  }
}
```

### Dreaming

**Source:** `concepts/dreaming.md`

Background memory consolidation. Three phases:

| Phase | What it does | Writes to |
|-------|-------------|-----------|
| **Light** | Ingests recent signals, stages candidates | `memory/.dreams/` |
| **Deep** | Scores candidates, promotes to long-term | `MEMORY.md` |
| **REM** | Theme reflection, boosts deep ranking | `DREAMS.md` |

**Scoring signals (Deep phase):**
- Frequency: 0.24 weight
- Relevance: 0.30 weight
- Query diversity: 0.15 weight
- Recency: 0.15 weight
- Consolidation: 0.10 weight
- Conceptual richness: 0.06 weight

**Slash commands:**
```
/dreaming status
/dreaming on
/dreaming off
```

**CLI:**
```bash
openclaw memory promote             # preview what would be promoted
openclaw memory promote --apply     # apply promotion
openclaw memory promote-explain "router vlan"  # explain why a memory was/wasn't promoted
openclaw memory rem-harness --json  # run REM phase analysis
```

**Config:**
```json
{
  "plugins": {
    "entries": {
      "memory-core": {
        "config": {
          "dreaming": {
            "enabled": true,
            "timezone": "America/New_York",
            "frequency": "0 3 * * *"
          }
        }
      }
    }
  }
}
```

### Honcho Memory

**Source:** `concepts/memory-honcho.md`

Cross-session AI-native memory. Honcho builds a user model over time — it studies patterns across all conversations and develops understanding of who you are, your preferences, communication style.

- External service (honcho.dev) or self-hostable
- Works alongside QMD, not as a replacement

### Session Pruning

**Source:** `concepts/session-pruning.md`

Manages context window limits automatically.

- **In-memory only** — never modifies on-disk session transcript
- Trims old tool results first (soft-trim), then content (hard-clear)
- Auto-enabled for Anthropic models (which have strict context limits)
- Preserves head + tail of conversation

**Config:**
```json5
{ contextPruning: { mode: "cache-ttl" } }
```

---

## 6. Automation Engine

**Source:** `automation/index.md`, `automation/cron-jobs.md`, `automation/standing-orders.md`, `automation/taskflow.md`, `automation/hooks.md`, `gateway/heartbeat.md`

### Decision Guide: Cron vs Heartbeat

| Use heartbeat when... | Use cron when... |
|----------------------|----------------|
| Multiple checks can batch together | Exact timing matters ("9:00 AM on weekdays") |
| Need conversational context | Task needs session isolation |
| Timing can drift slightly | Different model or thinking level |
| Reduce API calls by combining checks | One-shot reminders |
| | Output directly to a channel |

**Tip:** Batch multiple periodic checks into `HEARTBEAT.md` instead of creating many separate cron jobs.

### Heartbeat

**Source:** `gateway/heartbeat.md`

Heartbeat runs **periodic main-session turns** — the agent wakes up and checks HEARTBEAT.md for its todo list.

- Does NOT create background task records
- Default interval: 30 minutes
- Response contract: reply `HEARTBEAT_OK` if nothing needs attention. Reply the alert text if something needs attention (do NOT include `HEARTBEAT_OK`).

**HEARTBEAT.md format:**
```markdown
## Checks (rotate through 2-4x per day)
- [ ] Email: urgent unread?
- [ ] Calendar: events in next 24-48h?
- [ ] Weather: relevant if going out?

## State
Track last-checked in memory/heartbeat-state.json
```

**Config:**
```json5
{
  agents: {
    defaults: {
      heartbeat: {
        every: "30m",
        lightContext: true,    // only inject HEARTBEAT.md, not full workspace
        isolatedSession: false,
        prompt: "Read HEARTBEAT.md if it exists. Follow it strictly. If nothing needs attention, reply HEARTBEAT_OK.",
        ackMaxChars: 300
      }
    }
  }
}
```

### Cron Jobs

**Source:** `automation/cron-jobs.md`

Full job scheduler. Jobs persist at `~/.openclaw/cron/jobs.json`.

#### Schedule Types
| Type | CLI flag | Example |
|------|----------|---------|
| One-shot | `--at` | `--at "2026-12-25T09:00:00Z"` |
| Interval | `--every` | `--every "24h"` |
| Cron expr | `--cron` | `--cron "0 9 * * 1-5" --tz America/New_York` |

**Important:** Day-of-month + day-of-week uses **OR** logic (standard Vixie cron). Use `+1` modifier for AND: `0 9 15 * +1`.

#### Execution Styles
| Style | `--session` value | Best for |
|-------|------------------|---------|
| Main session | `main` | Reminders, system events |
| Isolated | `isolated` | Reports, background chores |
| Current session | `current` | Context-aware recurring work |
| Custom session | `session:custom-id` | Workflows that build on history |

#### Delivery Modes
| Mode | What happens |
|------|-------------|
| `announce` | Post summary to a channel |
| `webhook` | POST to a URL |
| `none` | Internal only |

#### Key CLI Commands
```bash
# Add a recurring job
openclaw cron add \
  --name "Morning Brief" \
  --cron "0 9 * * *" \
  --tz "America/New_York" \
  --session isolated \
  --message "Generate the daily morning briefing for Zeus..." \
  --announce --channel telegram

# One-shot reminder
openclaw cron add \
  --name "Meeting Reminder" \
  --at "30m" \
  --session main \
  --system-event "Meeting in 30 minutes" \
  --wake now \
  --delete-after-run

# List all jobs
openclaw cron list

# View run history for a job
openclaw cron runs --id <job-id>

# Pause a job (disable without deleting)
openclaw cron pause <job-id>

# Resume a paused job
openclaw cron resume <job-id>

# Delete a job
openclaw cron delete <job-id>
```

### Standing Orders

**Source:** `automation/standing-orders.md`

Standing Orders are **permanent operating authority** written into AGENTS.md. The agent reads them every bootstrap and executes automatically — no human prompting needed.

**Structure in AGENTS.md:**
```markdown
## Standing Orders

### Program: Morning Intelligence Brief
- **Scope**: Information gathering only
- **Trigger**: Heartbeat between 7-9 AM
- **Actions**: Check email inbox, check calendar (72h), check weather
- **Approval gate**: Required before sending to Zeus
- **Escalation**: Direct message if urgent (health/safety/financial)
- **Pattern**: Execute → Verify → Report (max 3 attempts)

### Program: Memory Maintenance
- **Scope**: MEMORY.md and daily notes only
- **Trigger**: Heartbeat, low-activity periods
- **Actions**: Write/organize daily notes, update MEMORY.md
- **Approval gate**: None required
```

**Execute-Verify-Report pattern:** Try → verify success → report. Max 3 attempts before escalating.

### Task Flow (ClawFlow)

**Source:** `automation/taskflow.md`

For complex multi-step workflows with durable state.

Two modes:
- **Managed mode**: OpenClaw owns the lifecycle, creates task records as steps
- **Mirrored mode**: Observes external task system (e.g., Todoist, Linear)

```bash
openclaw tasks flow list
openclaw tasks flow show <flowId>
openclaw tasks flow cancel <flowId>
```

### Hooks

**Source:** `automation/hooks.md`

Event-driven hooks — code that runs on OpenClaw lifecycle events.

**Event types:**
| Event | When it fires |
|-------|--------------|
| `command:new` | New user message arrives |
| `command:reset` | Session reset |
| `command:stop` | Stop requested |
| `session:compact:before` | Before context compaction |
| `session:compact:after` | After context compaction |
| `session:compact:patch` | Patch for compaction |
| `agent:bootstrap` | Agent bootstraps (can mutate `bootstrapFiles` array) |
| `gateway:startup` | Gateway starts up |
| `message:received` | Message received from channel |
| `message:transcribed` | Audio transcribed |
| `message:preprocessed` | Message preprocessed |
| `message:sent` | Message sent to channel |

**Hook = directory with `HOOK.md` + `handler.ts`**

Discovery order: bundled → plugin → managed (`~/.openclaw/hooks/`) → workspace (workspace hooks disabled by default).

---

## 7. Channels

**Source:** `channels/` docs, `channels/slack.md`, `channels/broadcast-groups.md`, `channels/group-messages.md`

### Supported Channels
OpenClaw supports all of these out-of-the-box:
- WhatsApp, Telegram, Discord, Slack, Signal, iMessage
- IRC, LINE, Google Chat, Mattermost, Matrix, Nextcloud Talk
- BlueBubbles, Zalo, Zalo User, Nostr, Feishu
- SMS/MMS (via carriers), Voice (Deepgram, ElevenLabs TTS)

### Slack (Full Feature Set)

**Source:** `channels/slack.md`

Socket mode (recommended) or HTTP mode.

**Interactive Block Kit replies** — send buttons/menus to Slack:
```
[[slack_buttons: Option 1 | Option 2 | Cancel]]
```
Clicking a button sends the response back to the agent as a user message.

**Preview streaming modes:**
- `off` — standard block delivery
- `partial` — live typing indicator + partial text updates
- `block` — deliver complete blocks as they finish
- `progress` — progress messages with completion (good for long tasks)

**Native exec approvals:** Shell exec confirmations appear as Slack button prompts — no out-of-band confirmation needed.

**All Slack events handled:** messages, app_mentions, DMs, reactions, file shares, thread replies.

### Broadcast Groups

**Source:** `channels/broadcast-groups.md`

One message → multiple agents process it simultaneously (WhatsApp only).

```json5
{
  bindings: [
    {
      broadcast: ["agent1", "agent2"],  // ← takes priority over "agentId"
      match: { channel: "whatsapp", peer: { kind: "group", id: "GROUP_JID@g.us" } }
    }
  ]
}
```

- Each agent gets an **isolated session** for the same message
- Strategy: `parallel` (default) or `sequential`
- Use case: send a briefing request → multiple domain agents respond independently

### WhatsApp Groups

**Source:** `channels/group-messages.md`

Groups require explicit activation.

```bash
openclaw channels login --channel whatsapp --account default
```

In the group chat, type: `/activation`

Activation modes:
- `mention` — only responds when @mentioned
- `always` — responds to every message

**Context injection:** Last 50 pending messages from the group are included in the agent's context.

**Per-group session isolation:** Each group has its own session — the agent remembers group-specific context.

---

## 8. Node System (Clients)

**Source:** `nodes/index.md`, `nodes/audio.md`, `nodes/voicewake.md`

### Node Types

| Node | Platform | What it does |
|------|----------|-------------|
| iOS | iPhone | Camera, location, SMS read/send, screen share |
| Android | Android phone | Camera, location, SMS, screen share |
| macOS | Mac | Screen, camera, microphone, local exec |
| Headless | VPS/server | No GUI — pure automation |

### What Nodes Give You

From `nodes/index.md`:
- **Canvas** — Shared visual space (render images, maps, diagrams)
- **Camera** — Live camera feed to agent
- **Screen sharing** — Agent sees your screen
- **Location** — GPS coordinates piped to agent
- **SMS** — Read/send SMS messages
- **Local exec** — Run commands on the node's machine (not the server)

### Audio

**Source:** `nodes/audio.md`

Audio transcription auto-detection order:
1. Active session model (if it supports audio)
2. Local CLI tools: `sherpa-onnx`, `whisper-cpp`, `whisper` (Python)
3. Gemini CLI
4. Provider auth (priority): OpenAI → Groq → Deepgram → Google → Mistral

Audio scope can be gated by chat type (DM only vs groups).

### Voice Wake Word

**Source:** `nodes/voicewake.md`

Single global wake-word list owned by the Gateway at `~/.openclaw/settings/voicewake.json`.

- macOS node and iOS node have local enable toggles
- Android uses manual mic button
- Wake words broadcast to all nodes when changed

---

## 9. HTTP APIs (OpenAI-Compatible)

**Source:** `gateway/openai-http-api.md`, `gateway/openresponses-http-api.md`

### Enabling the Endpoint

Disabled by default. Enable in config:

```json
{
  "gateway": {
    "http": {
      "endpoints": {
        "chatCompletions": { "enabled": true }
      }
    }
  }
}
```

### What's Exposed

Once enabled, on the same port as the gateway:
- `POST /v1/chat/completions` — OpenAI chat completions
- `GET /v1/models` — List agents as OpenAI "models"
- `GET /v1/models/{id}` — Get specific model
- `POST /v1/embeddings` — Embeddings
- `POST /v1/responses` — OpenAI Responses API (itemized)

### Model Routing

OpenClaw uses the `model` field as an **agent target**, not a raw model ID:
- `model: "openclaw"` → default agent
- `model: "openclaw/main"` → main (Legacy) agent
- `model: "openclaw/spartan"` → Spartan agent
- `model: "openclaw/empire"` → Empire agent

**Override headers:**
- `x-openclaw-model: <provider/model>` — override backend model for the turn
- `x-openclaw-session-key: <key>` — pin to specific session
- `x-openclaw-message-channel: <channel>` — set synthetic channel context

### Open WebUI Integration

```
Base URL: http://127.0.0.1:18789/v1
Authorization: Bearer <gateway-token>
Model: openclaw
```

### Security

**CRITICAL:** This endpoint has full operator-access (same as owner/operator credential). Keep it on loopback or Tailscale only. Never expose to public internet.

Auth: `Authorization: Bearer <gateway-token>` (same token as gateway auth).

### OpenResponses API

`POST /v1/responses` — item-based inputs with `previous_response_id` for session continuity. SSE streaming. Useful for building agentic web apps on top of your gateway.

---

## 10. Gateway Configuration

**Source:** `gateway/secrets.md`, `gateway/sandboxing.md`, `gateway/heartbeat.md`, `gateway/multiple-gateways.md`, `gateway/local-models.md`

### Secrets (SecretRef)

**Source:** `gateway/secrets.md`

Instead of hardcoding API keys in config, use SecretRef:

```json5
{
  channels: {
    telegram: {
      accounts: {
        default: {
          token: { "$secret": "telegram/bot-token" }  // ← SecretRef
        }
      }
    }
  }
}
```

**Providers:**
- `env` — from environment variable: `{ "$secret": "env/MY_VAR" }`
- `file` — from file: `{ "$secret": "file:///path/to/secret" }`
- `exec` — run a command: `{ "$secret": "exec/op read op://vault/item/field" }` ← **1Password integration**
- Vault (HashiCorp), sops, Doppler, Infisical

**1Password Integration:**
```json5
{
  gateway: {
    auth: {
      token: { "$secret": "exec/op read op://Personal/openclaw-gateway/token" }
    }
  }
}
```

**Audit secrets:**
```bash
openclaw secrets audit --check
```

### Sandboxing

**Source:** `gateway/sandboxing.md`

Modes:
- `off` — no sandboxing (default)
- `non-main` — sandbox everything except main session
- `all` — sandbox all sessions

Backends: Docker, SSH (remote exec), OpenShell

Scope: per-agent, per-session, or shared.

### Multiple Gateways

**Source:** `gateway/multiple-gateways.md`

Run multiple isolated gateways on the same machine:

```bash
OPENCLAW_PROFILE=work openclaw gateway start
OPENCLAW_PROFILE=personal openclaw gateway start
```

**Requirements for isolation:**
- Unique `OPENCLAW_CONFIG_PATH` per profile
- Unique `OPENCLAW_STATE_DIR` per profile
- Unique workspace per profile
- Unique port per profile (spacing: ≥20, recommended: 18789, 18809, 18829...)

Derived port map (base port `P`):
- `P+2` = browser control
- `P+9` to `P+108` = CDP (Chrome DevTools Protocol) range

### Local Models (Ollama)

**Source:** `gateway/local-models.md`

Recommended stack: LM Studio or Ollama.

For Ollama, use `api: "openai-responses"` (NOT the `/v1` compat endpoint for tool calling — breaks tool JSON).

**Native Ollama API URL:** `http://127.0.0.1:11434` — do NOT append `/v1`.

Compatibility flags for models that don't fully support tools:
```json5
{
  models: {
    "ollama/qwen2.5:7b": {
      supportsTools: false,           // disable tool calling
      requiresStringContent: true,   // force string content type
    }
  }
}
```

Context window warnings:
- Warn if context window < 32k tokens
- Block if context window < 16k tokens

**Models on this VPS:**
| Model | Size | Context | Status |
|-------|------|---------|--------|
| qwen2.5:7b | 6.2 GB | 32k | Primary (slow with tools on CPU) |
| qwen3:4b-fast | 2.5 GB | 256k | Available |
| qwen3:4b | 2.5 GB | 256k | Available |
| qwen3:8b | 5.2 GB | 40k | Available |
| qwen2.5:3b | ~2 GB | 32k | Available |
| qwen2.5:1.5b | ~1 GB | 32k | Fastest option |
| gemma4:latest | 8.5 GB | 128k | Available |
| nomic-embed-text | small | 2k | Embeddings only |

---

## 11. Security & Secrets

**Source:** `gateway/secrets.md`, `gateway/sandboxing.md`

### Secrets Rule
**NEVER put raw API keys in `openclaw.json`.** Use SecretRef exec provider with 1Password:

```json5
{ "$secret": "exec/op read op://00-24 SO (Sovereign Operations)/item-name/field" }
```

This requires `op` CLI installed and authenticated on the VPS.

### Credential Audit
```bash
openclaw secrets audit --check    # scan for exposed credentials in config
```

### Firewall
VPS runs DigitalOcean Cloud Firewall `onemindos-cloud-fw`. Gateway is on loopback + Tailscale only (not public internet).

### Auth Modes
- `token` — shared secret token (current setup)
- `password` — password-based
- `none` — no auth (private ingress only)
- `trusted-proxy` — identity-bearing proxy

---

## 12. Context Management

**Source:** `concepts/compaction.md`, `concepts/context-engine.md`

### Auto-Compaction

OpenClaw compacts context automatically when approaching model limits.

- **Off by default:** actually auto-compaction is ON by default
- Before compacting, agent is reminded to save important things to memory
- Custom compaction model: `agents.defaults.compaction.model`
- Config: `identifierPolicy: "strict"` preserves opaque tool-call IDs across compaction

**Pluggable via custom plugin** — implement `registerCompactionProvider()` for custom logic.

### Context Engine

**Source:** `concepts/context-engine.md`

The context engine controls what gets assembled as the prompt for each turn.

Pluggable via `plugins.slots.contextEngine`. Built-in: `legacy` engine.

A plugin can implement:
- `ingest` — process incoming message
- `assemble` — build the prompt (returns `messages + estimatedTokens + systemPromptAddition`)
- `compact` — custom compaction logic
- `afterTurn` — post-turn processing

If `ownsCompaction: true`, disables built-in auto-compaction.

Example: `@martian-engineering/lossless-claw` (context-preserving alternative engine)

---

## 13. Plugin System

**Source:** Runtime `openclaw plugins list`, `cli/` docs

### Plugin Categories

**Stock plugins (bundled with OpenClaw):**
```
active-memory     — Blocking memory sub-agent (already enabled)
memory-core       — Core memory search BM25+vector
memory-wiki       — Wiki-style memory lookup
memory-lancedb    — LanceDB vector backend (disabled)
active-memory     — Active memory plugin
llm-task          — Background LLM task runner
lobster           — (TBD)
openshell         — Shell execution
webhooks          — Webhook receiver
diffs             — File diff viewer
```

**Provider plugins:**
```
github-copilot    — GitHub Copilot models (enabled)
amazon-bedrock    — AWS Bedrock (enabled)
microsoft         — Azure/Microsoft (enabled)
ollama            — Local Ollama (enabled)
```

**Disabled (China/regional — incompatible or not needed):**
```
alibaba, byteplus, chutes, kimi, minimax, moonshot,
qianfan, qwen, stepfun, volcengine, vydra, xiaomi, zai
```

**Disabled (cloud providers — API keys not configured or not needed):**
```
anthropic, anthropic-vertex, arcee, cloudflare-ai-gateway,
deepseek, fireworks, groq, kilocode, litellm, mistral,
nvidia, openai, openrouter, synthetic, together, venice,
vercel-ai-gateway, xai, comfy
```

### Plugin CLI
```bash
openclaw plugins list                     # list all plugins
openclaw plugins enable <name>            # enable a plugin
openclaw plugins disable <name>           # disable a plugin
openclaw plugins info <name>              # details about a plugin
```

### ClawHub (Plugin Marketplace)
```bash
clawhub search <keyword>                  # search marketplace
clawhub install <package-name>            # install a plugin
clawhub list                              # list installed
```

---

## 14. CLI Reference

**Source:** Runtime `openclaw --help`, all CLI docs

### Agents
```bash
openclaw agents list --bindings           # list agents + their bindings
openclaw agents add <name>                # add agent (interactive or --non-interactive)
openclaw agents bind <id> --channel <ch>  # add channel binding
openclaw agents bindings                  # list all bindings
openclaw agents set-identity <id>         # update display name/emoji
openclaw agents delete <id>              # delete agent
```

### Memory
```bash
openclaw memory promote                   # preview memory promotion candidates
openclaw memory promote --apply           # write promoted items to MEMORY.md
openclaw memory rem-harness --json        # run REM phase
openclaw memory promote-explain <query>   # explain promotion for a query
```

### Cron
```bash
openclaw cron list                        # list all jobs
openclaw cron add --name ... --cron ...   # add a new job
openclaw cron pause <id>                  # pause (keep job, stop running)
openclaw cron resume <id>                 # resume paused job
openclaw cron delete <id>                 # delete job
openclaw cron runs --id <id>              # view run history
```

### Gateway
```bash
openclaw gateway start                    # start gateway
openclaw gateway stop                     # stop gateway
openclaw gateway restart                  # restart gateway
openclaw gateway status                   # show status
openclaw gateway logs                     # tail logs
```

### Plugins
```bash
openclaw plugins list                     # list all
openclaw plugins enable <name>
openclaw plugins disable <name>
```

### Secrets
```bash
openclaw secrets audit --check            # scan for exposed credentials
```

### TUI (Terminal UI)
```bash
openclaw tui                              # launch terminal UI
```

### Infer (Direct model test)
```bash
openclaw infer model run --model <id> "message"   # test a model directly
openclaw infer model run --model ollama/qwen2.5:7b "say PONG"
```

---

## 15. One Mind Agent Setup

**Source:** `AGENTS.md` (onemind-codex), `CLAUDE.md`, multi-agent docs

### The 9-Agent Architecture

One Mind OS uses 9 agents organized by domain:

| Agent ID | Emoji | Domain | Tagline | Role |
|----------|-------|--------|---------|------|
| `main` | ⚡ | All | Legacy Orchestrator | Central hub, all-domain awareness |
| `spartan` | 💪 | HP 27-28 | Body & Fitness | Physical performance, workouts, health |
| `oracle` | 🔮 | HP 25-26, 29-34 | Mind, Money, Meaning | Mental performance, finance, purpose |
| `eden` | 🌱 | LE | Environment | Sustainability, ecology, home environment |
| `haven` | 🏡 | LE 50-74 | Family & Foundation | Relationships, family, community, home |
| `empire` | 💰 | GE 75-99 | Wealth & Endeavor | Business, ventures, brand, income |
| `guardian` | 🛡️ | SO 11 | Security | Cybersecurity, credentials, threat monitoring |
| `grid` | ⚙️ | SO 01, 09 | Infrastructure | VPS, automation, system health |
| `trinity` | 🔺 | SO 03 | Coordination | Agent meta-orchestration |

### Routing Architecture

```
INBOUND MESSAGE (Telegram/WhatsApp/etc.)
         ↓
    [Channel binding lookup]
         ↓
    Legacy (main) ← default for all channels
         ↓ (delegates via subagent)
    ┌────┬────┬────┬────┬────┬────┬────┬────┐
    Spartan Oracle Eden Haven Empire Guardian Grid Trinity
```

Legacy is the **front door**. Specialist agents are invoked via delegation (subagent/ACP calls from Legacy) or via separate channel accounts once those are connected.

### Workspace Files Per Agent

Each agent has its own workspace with domain-specific:
- **SOUL.md** — Personality and voice for that domain
- **IDENTITY.md** — Name, emoji, role description
- **AGENTS.md** — Domain-specific standing orders
- **USER.md** — Zeus context (same across all agents)

### Creating Agents (Native CLI)
```bash
# Add agents non-interactively
openclaw agents add spartan --workspace /root/.openclaw/workspace-spartan --non-interactive
openclaw agents add oracle  --workspace /root/.openclaw/workspace-oracle  --non-interactive
openclaw agents add eden    --workspace /root/.openclaw/workspace-eden    --non-interactive
openclaw agents add haven   --workspace /root/.openclaw/workspace-haven   --non-interactive
openclaw agents add empire  --workspace /root/.openclaw/workspace-empire  --non-interactive
openclaw agents add guardian --workspace /root/.openclaw/workspace-guardian --non-interactive
openclaw agents add grid    --workspace /root/.openclaw/workspace-grid    --non-interactive
openclaw agents add trinity --workspace /root/.openclaw/workspace-trinity --non-interactive

# Set display identities
openclaw agents set-identity main     --name "Legacy"   --emoji "⚡"
openclaw agents set-identity spartan  --name "Spartan"  --emoji "💪"
openclaw agents set-identity oracle   --name "Oracle"   --emoji "🔮"
openclaw agents set-identity eden     --name "Eden"     --emoji "🌱"
openclaw agents set-identity haven    --name "Haven"    --emoji "🏡"
openclaw agents set-identity empire   --name "Empire"   --emoji "💰"
openclaw agents set-identity guardian --name "Guardian" --emoji "🛡️"
openclaw agents set-identity grid     --name "Grid"     --emoji "⚙️"
openclaw agents set-identity trinity  --name "Trinity"  --emoji "🔺"

# Verify
openclaw agents list --bindings
```

---

## 16. Cron Jobs — One Mind Schedule

**Source:** `automation/cron-jobs.md`

These jobs are **configured but paused**. Activate when ready.

### Job List

| Job | Schedule | Agent | Purpose |
|-----|----------|-------|---------|
| Morning Brief | 9:00 AM daily | main | Daily context brief → Telegram |
| Evening Review | 9:00 PM daily | main | EOD debrief, memory write |
| Weekly Review | Sunday 8:00 AM | main | Weekly scorecard |
| Memory Promote | 3:00 AM daily | main | Dreaming / memory consolidation |
| Codex Health | Monday 10:00 AM | main | Vault health check |
| Workout Plan | Monday 7:00 AM | spartan | Weekly fitness plan |
| Finance Check | Friday 5:00 PM | oracle | Weekly money check |
| Security Scan | Daily 2:00 AM | guardian | Run security checks |
| Infra Health | Every 6h | grid | VPS + service health |

### Adding Jobs (CLI — execute disabled via pause immediately after)
```bash
# Morning Brief
openclaw cron add \
  --name "Morning Brief" \
  --cron "0 9 * * *" \
  --tz "America/New_York" \
  --session isolated \
  --message "You are Legacy, the One Mind central orchestrator. Generate Zeus's morning brief: check Todoist tasks due today, scan calendar for today's events, check any urgent emails, provide a 3-sentence situation report. Keep it concise." \
  --tag one-mind

# Evening Review
openclaw cron add \
  --name "Evening Review" \
  --cron "0 21 * * *" \
  --tz "America/New_York" \
  --session isolated \
  --message "You are Legacy. Generate Zeus's evening review: (1) wins today, (2) what was left undone, (3) write today's memory entry to daily note at memory/$(date +%Y-%m-%d).md, (4) set intention for tomorrow." \
  --tag one-mind

# Weekly Review — Sundays
openclaw cron add \
  --name "Weekly Review" \
  --cron "0 8 * * 0" \
  --tz "America/New_York" \
  --session isolated \
  --message "Generate the weekly One Mind scorecard. Review completed tasks, domain health (SO/HP/LE/GE), wins, lessons, and priorities for next week." \
  --tag one-mind

# Memory Promotion — 3 AM daily
openclaw cron add \
  --name "Memory Promote" \
  --cron "0 3 * * *" \
  --tz "America/New_York" \
  --session isolated \
  --message "Run memory promotion: review recent daily memory files, identify high-signal items, update MEMORY.md with any important long-term context. Run: openclaw memory promote --apply" \
  --tag one-mind

# Infra Health — every 6h
openclaw cron add \
  --name "Infra Health" \
  --cron "0 */6 * * *" \
  --tz "America/New_York" \
  --session isolated \
  --agent grid \
  --message "Check VPS health: disk usage, service status (openclaw-gateway, n8n, docker), memory, CPU load. Report any anomalies. Reply HEALTH_OK if all good." \
  --tag one-mind

# After adding — PAUSE ALL (do not activate yet)
# openclaw cron pause <id>  for each job above
```

> **Note:** Jobs are PAUSED upon creation. Run `openclaw cron resume <id>` to activate each one when ready.

### Activation Sequence (When Ready)
1. Confirm Telegram channel is connected (`openclaw channels status`)
2. Test each job manually first: `openclaw cron run <id>`
3. Resume jobs one at a time: `openclaw cron resume <id>`
4. Monitor first 3 runs: `openclaw cron runs --id <id>`

---

## 17. Config Reference — This VPS

**Source:** `/root/.openclaw/openclaw.json` read directly

### Current Config State (April 2026)

```json
{
  "agents": {
    "defaults": {
      "workspace": "/root/.openclaw/workspace",
      "model": {
        "primary": "github-copilot/claude-sonnet-4-6",
        "fallbacks": ["github-copilot/gpt-4o"]
      },
      "heartbeat": {
        "every": "30m",
        "lightContext": true,
        "prompt": "Read HEARTBEAT.md if it exists (workspace context). Follow it strictly. Do not infer or repeat old tasks from prior chats. If nothing needs attention, reply HEARTBEAT_OK.",
        "ackMaxChars": 300
      }
    },
    "list": [
      { "id": "main" },
      { "id": "spartan", "workspace": "/root/.openclaw/workspace-spartan" },
      { "id": "oracle",  "workspace": "/root/.openclaw/workspace-oracle" },
      { "id": "eden",    "workspace": "/root/.openclaw/workspace-eden" },
      { "id": "haven",   "workspace": "/root/.openclaw/workspace-haven" },
      { "id": "empire",  "workspace": "/root/.openclaw/workspace-empire" },
      { "id": "guardian","workspace": "/root/.openclaw/workspace-guardian" },
      { "id": "grid",    "workspace": "/root/.openclaw/workspace-grid" },
      { "id": "trinity", "workspace": "/root/.openclaw/workspace-trinity" }
    ]
  },
  "gateway": {
    "mode": "local",
    "auth": { "mode": "token" },
    "port": 18789,
    "bind": "loopback",
    "tailscale": { "mode": "serve" },
    "http": {
      "endpoints": {
        "chatCompletions": { "enabled": true }
      }
    }
  },
  "memory": {
    "backend": "qmd",
    "qmd": { "sessions": { "enabled": true } }
  },
  "plugins": {
    "entries": {
      "active-memory": { "enabled": true },
      "memory-core": {
        "config": {
          "dreaming": {
            "enabled": true,
            "timezone": "America/New_York",
            "frequency": "0 3 * * *"
          }
        }
      }
    }
  }
}
```

---

## 18. Teaching Notes for Students

> This section is for One Mind students learning OpenClaw for the first time.

### The Mental Model

Think of OpenClaw like a personal switchboard:

```
Your Phone (WhatsApp/Telegram)
         ↓
   OpenClaw Gateway
   (your brain's extension)
         ↓
One of your 9 AI agents
(Legacy for general, Spartan for fitness, Empire for business...)
         ↓
Claude / GPT / Local Model
(the actual intelligence engine)
```

The **gateway** is the switchboard. **Agents** are specialized personalities. **Models** are the brains inside. The **channels** (WhatsApp, Telegram) are the phones.

### The 5 Things to Understand First

1. **Everything is a file.** Agent personality = `SOUL.md`. Instructions = `AGENTS.md`. Memory = `MEMORY.md`. Plain text. No database to understand.

2. **Memory is tiered.** Short-term (session) → Medium (QMD search index) → Long-term (MEMORY.md file). Dreaming promotes the best stuff from medium to long-term automatically.

3. **Heartbeat = agent stays alive.** Every 30 minutes, the agent wakes up, reads `HEARTBEAT.md`, and checks if anything needs attention. This is how your AI notices things without you prompting it.

4. **Standing Orders = permanent authority.** Write rules in `AGENTS.md` like "every morning, check my calendar." The agent executes these automatically, every session, forever.

5. **Multiple agents = multiple personalities.** One server, 9 agents, each with its own workspace and personality. Your fitness agent talks like a trainer. Your wealth agent talks like a CFO. The same system, different voices.

### Common Mistakes

| Mistake | Fix |
|---------|-----|
| Putting API keys directly in config | Use SecretRef: `{ "$secret": "exec/op read op://vault/item" }` |
| Using Ollama `/v1` URL for agent API | Use `http://127.0.0.1:11434` (no `/v1`) |
| Editing `openclaw.json` by hand without backup | Always run `cp openclaw.json openclaw.json.bak` first |
| Activating all cron jobs at once | Add them all, test one at a time with `openclaw cron run <id>` |
| Expecting heartbeat to do heavy work | Heartbeat is for lightweight checks. Use cron with `--session isolated` for heavy tasks |

### Where Everything Lives

```
On the VPS:
  /root/.openclaw/           ← everything OpenClaw
  /root/.openclaw/workspace/ ← Legacy (main) agent home
  /usr/lib/node_modules/openclaw/docs/  ← full documentation

On your Mac (this repo):
  onemind-codex/_system/openclaw/  ← this reference file lives here
  onemind-codex/AGENTS.md          ← overall framework context
```

---

*Generated April 2026. Source: all docs in `/usr/lib/node_modules/openclaw/docs/` + direct VPS inspection. Updated by Zeus Delacruz / One Mind OS.*
