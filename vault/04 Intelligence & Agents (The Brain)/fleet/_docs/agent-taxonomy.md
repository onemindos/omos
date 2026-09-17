---
title: "Agent Taxonomy"
type: reference
domain: SO
status: active
created: 2026-06-07
updated: 2026-06-07
tags: [agents, taxonomy, architecture, fleet, cli, sdk, api, infrastructure]
---

# Agent Taxonomy

OneMind OS organizes AI agents into four tiers. Each tier has a different access model, lifecycle, and role in the fleet.

```
┌─────────────────────────────────────────────────────────────────┐
│  TIER 4 — AUTONOMOUS INFRASTRUCTURE                             │
│  Hermes (Legacy) · OpenClaw · Pi Agent                          │
│  Always-on, persistent, NATS-wired, fleet-level                 │
├─────────────────────────────────────────────────────────────────┤
│  TIER 3 — VISUAL / PROTOTYPING AGENTS                           │
│  Grok Build                                                     │
│  Browser-based, stateless, full-stack generation                │
├─────────────────────────────────────────────────────────────────┤
│  TIER 2 — ON-DEMAND CODING AGENTS (Terminal)                    │
│  Model-locked: Claude Code · Codex CLI · Gemini CLI             │
│  Model-agnostic: OpenCode                                       │
├─────────────────────────────────────────────────────────────────┤
│  TIER 1 — MODELS / PROVIDERS → SDK / API                        │
│  Azure OpenAI · xAI · OpenRouter · Ollama · Qwen · MiniMax     │
│  Raw model access — the foundation everything above builds on   │
└─────────────────────────────────────────────────────────────────┘
```

---

## Tier 1 — Models · Providers → SDK / API

The raw intelligence layer. No interface, no UI — pure API or SDK access. Every tier above builds on one or more of these.

### Access Pattern

All active providers in this stack expose an OpenAI-compatible endpoint:

```python
from openai import OpenAI

client = OpenAI(
    api_key="...",          # from 1Password
    base_url="https://..."  # provider-specific
)

response = client.chat.completions.create(
    model="model-name",
    messages=[{"role": "user", "content": "Do this."}]
)
```

One SDK pattern, swappable providers.

### Active Provider Registry

| Provider | Models | API Base | Notes |
|----------|--------|----------|-------|
| **Azure OpenAI** | GPT-5.x, GPT-4.1, o4-mini, GPT-4o | `onemind-openai.openai.azure.com/openai/v1` | ⚠️ GPT-5.x + o4-mini require `max_completion_tokens`, not `max_tokens`. ~$30/mo |
| **xAI** | Grok-3, Grok-3 Mini, Grok-2 Vision | `api.x.ai/v1` | 131k context. Real-time web data. Aurora TTS. OpenAI-compatible. |
| **OpenRouter** | All of the above + 100+ | `openrouter.ai/api/v1` | Fallback routing. Single API key. ~$20/mo |
| **Ollama** | Llama, Mistral, Phi, Gemma, Qwen, CodeLlama | `localhost:11434/v1` | Local inference. No cloud, no cost. Metal GPU on Mac. |
| **Qwen** | Qwen series | — | See `ai/providers/Qwen.md` |
| **MiniMax** | MiniMax-Text, Speech | — | TTS + LLM. Hermes native provider. |
| **Nvidia NIM** | Llama, Mistral, Nemotron | — | GPU-optimized enterprise inference. See `ai/providers/nvidia.md` |

> **Note:** Anthropic is a primary provider (Legacy's main model is claude-sonnet-latest) but no clean provider card exists in the registry yet — the `ai/providers/anthropic.md` file currently contains Claude Code system prompts, not API reference. Credentials in 1Password.

### When to Use Tier 1 Directly

- Building a custom pipeline or script from scratch
- Embedding generation or batch inference
- Cost optimization — calling API directly skips agent overhead
- Local/air-gapped requirements → Ollama

---

## Tier 2 — On-Demand Coding Agents (Terminal)

Terminal-based agents you invoke for a specific task. They reason autonomously, use tools, write code, and return a result. Two sub-groups: **model-locked** (tied to one provider) and **model-agnostic** (you choose).

### Model-Locked

#### Claude Code (Anthropic)

> Autonomous coding agent from Anthropic. Full repo context, PR creation, code review.

```bash
claude                                       # Interactive REPL
claude "fix the auth bug"                   # One-shot task
claude --model opus                          # Force Opus model
claude --resume                              # Resume last session
claude --allowedTools bash,read,write        # Restrict tool surface
```

| | |
|---|---|
| **Context file** | `CLAUDE.md` at repo root |
| **Config** | `~/.claude/settings.json` |
| **Memory** | `~/.claude/projects/` per-project |
| **Model** | Claude Sonnet / Opus (Anthropic only) |
| **Cost** | ~$50/mo (Anthropic API usage) |

- **OneMind use:** Complex multi-file features, large refactors, PR workflows. Best when Anthropic model quality matters.
- **Hermes skill:** `autonomous-ai-agents/claude-code`

---

#### Codex CLI (OpenAI)

> OpenAI's terminal coding agent. Sandboxed shell execution.

```bash
codex                               # Interactive REPL
codex "implement the auth flow"     # One-shot task
codex --model o4-mini               # Specify model
codex --approval auto               # Skip per-action confirmation
codex --full-auto                   # Full autonomy mode
```

| | |
|---|---|
| **Context file** | `AGENTS.md` at repo root |
| **Config** | `~/.codex/config.yaml` |
| **Global instructions** | `~/.codex/instructions.md` |
| **Model** | o4-mini, GPT-4.1 (OpenAI/Azure only) |
| **Cost** | Included in Azure ~$30/mo |

- **Key trait:** Sandboxed — shell commands run in an isolated environment. Safe for untrusted ops.
- **OneMind use:** Tasks requiring shell execution during dev. Kanban worker mode via `kanban-codex-lane` skill.
- **Hermes skill:** `autonomous-ai-agents/codex`

---

#### Gemini CLI (Google)

> Google's coding agent. Massive context window for codebase-wide analysis.

```bash
gemini                              # Interactive REPL
gemini "analyze this codebase"      # One-shot task
```

| | |
|---|---|
| **Context file** | `GEMINI.md` at repo root |
| **Model** | Gemini (Google only) |
| **Context window** | 1M+ tokens |
| **Cost** | Free tier available, Pro ~$20/mo |

- **OneMind use:** Large context tasks — ingest an entire repo or docs set in one shot. Architecture review, documentation generation.
- **Hermes skill:** Terminal invocation only

---

### Model-Agnostic

#### OpenCode

> Open-source terminal coding agent. Any provider, any model.

```bash
opencode                                        # Interactive TUI
opencode --provider openrouter                  # Use OpenRouter
opencode --model deepseek/deepseek-coder        # Specific model
opencode --provider ollama --model llama3.1     # Local model
```

| | |
|---|---|
| **Context file** | `AGENTS.md` at repo root |
| **Config** | `~/.opencode/config.json` |
| **Model** | Any — your choice |
| **Cost** | Depends on provider — can be $0 with local models |

- **Key trait:** Provider flexibility — same interface regardless of model. Swap to cheapest/fastest/most private per task.
- **OneMind use:** Cost-optimized coding. Sensitive repos where cloud agents aren't acceptable. When a specific open model outperforms Claude for the task.
- **Hermes skill:** `autonomous-ai-agents/opencode`

---

### Tier 2 Comparison

| Agent | Provider | Context File | Sandboxed | Model Flexibility | Best For |
|-------|----------|-------------|-----------|------------------|---------|
| Claude Code | Anthropic only | `CLAUDE.md` | No | None | Complex features, PRs |
| Codex CLI | OpenAI/Azure only | `AGENTS.md` | Yes | None | Shell-heavy tasks, Kanban |
| Gemini CLI | Google only | `GEMINI.md` | No | None | Massive context, full-repo analysis |
| OpenCode | Any provider | `AGENTS.md` | No | Full | Cost control, private repos, local models |

---

## Tier 3 — Visual / Prototyping Agents (Browser)

Not terminal agents — browser-based tools for visual code generation. Stateless per session, no project context file.

### Grok Build (xAI)

> Browser-based autonomous full-stack code generator. Powered by Grok-3.

```
Access: https://build.x.ai
Model:  Grok-3 (xAI)
```

| Feature | Notes |
|---------|-------|
| Full-stack generation | Front-end + back-end from a single prompt |
| Live preview | In-browser rendering of generated app |
| Iterative editing | Prompt-driven edits on generated output |
| One-click deploy | Direct web deploy from the browser |

| | |
|---|---|
| **Context file** | None — stateless per session |
| **Model** | Grok-3 (xAI only) |
| **Interface** | Browser only — not a terminal agent |
| **Cost** | xAI subscription or pay-per-use API |

- **OneMind use:** Greenfield prototypes and rapid scaffolding before handing off to Hermes/Claude Code. Not for repo-based work — no persistent context.
- **Key distinction from Tier 2:** No terminal, no AGENTS.md, no project memory. Each session is fresh. Use it to generate a starting point, then bring the output into a real repo.

---

## Tier 4 — Autonomous Infrastructure

Always-on agents. Not tools you invoke — they **are** the infrastructure. They run persistently, own the fleet, and orchestrate everything below them.

### Hermes Agent (Legacy)

> THE brain of OneMind OS. Master orchestrator. Always running on `omos-mbp`.

```bash
hermes                          # Interactive WebUI / CLI
hermes chat "do this"           # One-shot from shell
hermes config set model ...     # Configure provider/model
hermes tools                    # List all available tools
hermes skills list              # List loaded skills
hermes cron list                # Active scheduled jobs
hermes session list             # Past session history
```

**What makes it infrastructure:**
- Persistent memory across all sessions (`agentmemory`, localhost:3111)
- Cron engine — scheduled jobs run autonomously 24/7 without Zeus present
- Multi-agent delegation — spawns subagents, routes to Tier 2/3 agents
- NATS pub/sub — publishes and consumes on the bus directly
- Platform delivery — Telegram, WebUI, Discord, CLI, email
- 53 MCP tools (Todoist, Fastmail, HA, NATS, filesystem, etc.)
- Skill library — procedural memory for recurring tasks
- Full cross-domain vault access (all 4 pillars)

**Key directories:**
```
~/.hermes/config.yaml       ← provider, model, tools config
~/.hermes/skills/           ← skill library (procedural memory)
~/.hermes/memories/         ← persistent memory store
~/.hermes/cron/             ← scheduled job definitions
~/.hermes/plugins/          ← plugin system
~/.hermes/profiles/         ← one dir per profile
```

**Profiles — specialized agents on the same runtime:**

| Profile | Role | Domain |
|---------|------|--------|
| `default` (Legacy) | Master orchestrator | All pillars |
| `sovereignty` | Strategy, doctrine | SO |
| `eternal` | Knowledge synthesis | SO + cross-pillar |
| `zenith` | Building, coding, execution | SO + GE |
| `unity` | Communications, outreach | LE + GE |
| `sentinel` | Monitoring, security | SO |

**Supported backends:** local, SSH, Docker, Modal

- **Dashboard:** localhost:9119
- **WebUI:** omos-mbp.tail717ba0.ts.net:8787
- **Skill:** `autonomous-ai-agents/hermes-agent`

---

### OpenClaw

> Self-hosted AI agent gateway and runtime. Zeus's interface and delivery layer.

```bash
# Runs on VPS (onemind-main) via systemd
systemctl status openclaw
systemctl restart openclaw
journalctl -u openclaw -f

# Agent workspaces on VPS
/home/onemind/onemind-os/openclaw/workspaces/{agent}/
```

**What it does:**
- Receives messages from Zeus via Telegram and direct channels
- Routes to the correct domain agent workspace
- Manages isolated agent workspaces — each agent has own `IDENTITY.md`, `SOUL.md`, shared `CODEX.md` + `USER.md`
- Memory indexing — embeds vault content into pgvector
- REST API + WebSocket streaming for real-time agent interaction

**Fleet agents running on OpenClaw:**
Legacy, Spartan, Oracle, Haven, Empire, Guardian, Grid, Trinity, Eden

**Key distinction from Hermes:**
```
Zeus → OpenClaw (interface, routing, delivery)
              ↓ NATS ↓
         Hermes/Legacy (reasoning, execution, memory)
```

- OpenClaw = what Zeus **talks to**
- Hermes = what actually **thinks and executes**
- NATS = the **nervous system** between them

- **Hosting:** VPS `onemind-main` (DigitalOcean), systemd
- **Skill:** `openclaw-imports/openclaw-fleet`

---

### Pi Agent

> Physical node agent. Always-on on Raspberry Pi hardware at each property node.

```bash
# Runs as a persistent service on each Pi node
# Communicates via NATS (nats-py)
```

**Hardware:**

| Node | Pi Model | Role |
|------|----------|------|
| House | Pi 5 | Central hub, Frigate NVR, voice satellite |
| Barn | Pi 4 | Greenhouse sensors, camera, actuators |
| Office | Pi 4 | Development, testing, edge inference |

**Stack:**
- OS: Raspberry Pi OS 64-bit
- Runtime: Python 3.11+
- Messaging: `nats-py`
- Voice: Wyoming protocol (HA-compatible)
- Vision: OpenCV + ONNX Runtime (local inference)
- Sensors: GPIO, I2C, SPI (BME280, soil moisture, etc.)

**NATS subjects (from vault):**
```
telemetry.{node}.temperature
telemetry.{node}.humidity
telemetry.{node}.camera.motion
cmd.{node}.relay.{id}
cmd.{node}.voice.speak
```

**Why Tier 4 (not Tier 3):**
Pi Agent is always-on physical infrastructure — it doesn't get invoked for a task and shut down. It runs continuously at each node, publishing telemetry and listening for commands, exactly like Hermes and OpenClaw do in their domains.

- **HA bridge:** MQTT on :1883
- **Skill:** `autonomous-ai-agents/pi-coding-agent`
- **Hardware docs:** `SO/03 Infrastructure/hardware/`

---

### Tier 4 Comparison

| System | Hosting | Always On | NATS | Memory | Orchestrates |
|--------|---------|-----------|------|--------|-------------|
| Hermes (Legacy) | `omos-mbp` (Mac), LaunchAgent | Yes | Yes — pub/sub | Yes — agentmemory | Tiers 1–3 |
| OpenClaw | `onemind-main` (VPS), systemd | Yes | Yes — NATS gateway | Yes — pgvector | Domain agents |
| Pi Agent | Pi hardware per node | Yes | Yes — nats-py | No | Sensors, GPIO, voice |

---

## How the Tiers Work Together

```
Zeus issues a command
        ↓
OpenClaw (Tier 4) — receives, routes to correct agent
        ↓
Hermes/Legacy (Tier 4) — reasons, plans, delegates
        ↓
    ┌──────────────────────────────────────────────┐
    │ Tier 2: Claude Code  → complex coding task   │
    │ Tier 2: Codex CLI    → sandboxed shell task  │
    │ Tier 2: Gemini CLI   → large context task    │
    │ Tier 2: OpenCode     → any-model coding      │
    │ Tier 3: Grok Build   → browser prototype     │
    │ Tier 1: Azure/xAI    → raw inference         │
    └──────────────────────────────────────────────┘
        ↓ results via NATS ↓
Pi Agents (Tier 4) — physical execution at nodes
        ↓
OpenClaw delivers back to Zeus (Telegram / WebUI)
```

---

## Decision Guide

| Situation | Reach For |
|-----------|-----------|
| Raw model output in a script or pipeline | Tier 1 — direct API |
| Local / air-gapped / zero-cost inference | Tier 1 — Ollama |
| Complex coding: multi-file, PR, refactor | Tier 2 — Claude Code |
| Coding with shell execution needed | Tier 2 — Codex CLI |
| Full-repo analysis, massive context | Tier 2 — Gemini CLI |
| Need specific model / cheapest path | Tier 2 — OpenCode |
| Greenfield prototype, visual generation | Tier 3 — Grok Build |
| Orchestration, cron, memory, delegation | Tier 4 — Hermes |
| Message routing, platform delivery | Tier 4 — OpenClaw |
| Physical node, sensors, GPIO, voice | Tier 4 — Pi Agent |

---

## Related Docs

- `fleet/_docs/` — This file and other fleet architecture docs
- `fleet/hermes/hermes-architecture-plan.md` — Hermes deep-dive
- `fleet/openclaw/06 OpenClaw.md` — OpenClaw full reference
- `nodes/digital/services-registry/ai/agents/` — Per-agent service cards
- `nodes/digital/services-registry/ai/providers/` — Per-provider API cards
- [[ONEMIND-NATS-OS-PLAN.md]] — NATS subject contract
