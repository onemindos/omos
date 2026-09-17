# OpenClaw Course — One Mind OS Student Guide

> **For:** One Mind community students learning to set up and operate OpenClaw as their personal AI gateway.
> **By:** Zeus Delacruz / One Mind OS
> **Prerequisites:** A VPS (DigitalOcean, Hetzner, etc.) or local machine, basic terminal comfort, an AI API key (GitHub Copilot, OpenAI, Anthropic, or Ollama for local).
> **Time:** ~4 hours to complete all modules. Each module is self-contained — do them in order or jump to what you need.

---

## Course Overview

| Module | Topic | What You'll Build |
|--------|-------|-------------------|
| **01** | [What Is OpenClaw?](#module-01--what-is-openclaw) | Mental model of the system |
| **02** | [Installation & First Boot](#module-02--installation--first-boot) | Running gateway for the first time |
| **03** | [Your First Conversation](#module-03--your-first-conversation) | Chatting with your AI via CLI |
| **04** | [The Workspace — Files Are Everything](#module-04--the-workspace--files-are-everything) | SOUL.md, AGENTS.md, USER.md |
| **05** | [Memory — How Your AI Remembers](#module-05--memory--how-your-ai-remembers) | Session memory, QMD search, dreaming |
| **06** | [Automation — AI That Works While You Sleep](#module-06--automation--ai-that-works-while-you-sleep) | Heartbeat, cron jobs, standing orders |
| **07** | [Multi-Agent — Specialized AI Personalities](#module-07--multi-agent--specialized-ai-personalities) | Adding specialist agents |
| **08** | [Channels — Telegram & Discord](#module-08--channels--telegram--discord) | Connecting messaging platforms |
| **09** | [Security & Secrets](#module-09--security--secrets) | SecretRef, exec policy, credential management |
| **10** | [The HTTP API — Connect Anything](#module-10--the-http-api--connect-anything) | OpenAI-compatible endpoint |
| **11** | [Advanced — Hooks, Nodes, Plugins](#module-11--advanced--hooks-nodes-plugins) | Event-driven automation, macOS/iOS nodes |
| **12** | [One Mind Integration](#module-12--one-mind-integration) | Mapping OpenClaw to the CODEX framework |

---

## Module 01 — What Is OpenClaw?

### The Big Picture

OpenClaw is a **self-hosted AI gateway**. You run it on your own server, connect your AI models, connect your messaging channels, and it becomes your personal AI operating system.

```
YOUR PHONE (WhatsApp / Telegram / Discord)
         ↓ message
OPENCLAW GATEWAY (your server)
         ↓
AI MODEL (Claude, GPT, Gemini, Ollama)
         ↓ response
BACK TO YOUR PHONE
```

### Why Not Just Use ChatGPT?

| Feature | ChatGPT App | OpenClaw |
|---------|-------------|----------|
| You own the server | ❌ OpenAI owns it | ✅ Your machine |
| Multi-channel | ❌ One app only | ✅ WhatsApp, Telegram, Discord, Slack, SMS... |
| Memory across sessions | Limited | ✅ BM25 + vector hybrid search |
| Runs while you sleep | ❌ | ✅ Cron jobs, heartbeat, standing orders |
| Multiple AI personalities | ❌ One assistant | ✅ Unlimited agents, each with own workspace |
| Tool use (shell, files, web) | Limited | ✅ Full system access |
| Works with any model | ❌ GPT only | ✅ Claude, GPT, Gemini, Ollama, 30+ providers |

### The Parts

| Part | What It Is | Analogy |
|------|-----------|---------|
| **Gateway** | The server process | The building |
| **Agent** | An AI personality with its own workspace | An employee at a desk |
| **Workspace** | A folder of files the agent reads at startup | The employee's desk |
| **Model** | The AI brain powering the agent | The employee's brain |
| **Channel** | A messaging platform connection | A phone line into the building |
| **Memory (QMD)** | Search engine over past conversations | The filing cabinet |
| **Cron** | Scheduled tasks | The alarm clock |
| **Heartbeat** | Periodic agent check-in | The agent looking at its watch |

### Key Principle: Everything Is a File

This is the most important thing to understand. OpenClaw agents are configured through **plain text files**:

| File | What It Does |
|------|-------------|
| `SOUL.md` | Agent personality and voice |
| `AGENTS.md` | Operating rules and standing orders |
| `USER.md` | Context about you (the user) |
| `IDENTITY.md` | Agent name, emoji, role |
| `HEARTBEAT.md` | Periodic check list |
| `MEMORY.md` | Long-term curated memories |
| `BOOTSTRAP.md` | First-run setup guide |
| `TOOLS.md` | Notes about available tools |

Change the file → agent's behavior changes. No database. No API. Just edit a text file.

---

## Module 02 — Installation & First Boot

### Prerequisites

- **Server:** VPS with 2+ GB RAM (DigitalOcean $12/mo droplet works), or a Mac/Linux desktop
- **Node.js:** v20+ (`node --version`)
- **npm:** v10+ (`npm --version`)

### Install OpenClaw

```bash
# Install globally
npm install -g openclaw

# Verify
openclaw --version
```

### First Run — The Onboarding Wizard

```bash
openclaw onboard
```

The wizard walks you through:
1. **Model provider** — pick one (GitHub Copilot, OpenAI, Anthropic, Bedrock, Ollama...)
2. **API key** — paste your key (stored securely)
3. **Gateway auth** — creates a token for dashboard access
4. **First agent** — sets up the `main` agent with default workspace

### Start the Gateway

```bash
# Foreground (see all output)
openclaw gateway

# Or install as a system service
openclaw gateway install    # creates launchd (Mac) or systemd (Linux) service
openclaw gateway start      # start the service
openclaw gateway status     # check it's running
```

### Access the Dashboard

After starting, the gateway prints a URL:

```
Dashboard: http://localhost:18789
Token: abc123...
```

Open that URL in your browser. Paste the token. You're in.

### Verify Everything Works

```bash
# Direct model test (bypasses agents — tests the model connection)
openclaw infer model run --model <your-model-id> "Say hello"

# Agent conversation test
openclaw agent --message "Hello, what can you do?"
```

If both return responses, your gateway is working.

> **Checkpoint:** You should now have a running gateway with one agent. You can chat via CLI and see the dashboard.

---

## Module 03 — Your First Conversation

### CLI Chat

```bash
# Single message (gets a response, exits)
openclaw agent --message "What is the capital of France?"

# Interactive session (continuous chat)
openclaw tui
```

The TUI (Terminal UI) gives you a full chat interface in your terminal.

### Dashboard Chat

Open the dashboard URL in your browser. The chat panel is the main interface — type messages there just like ChatGPT.

### Understanding Sessions

Every conversation happens inside a **session**. Sessions are persistence boundaries:

| Session Type | When | Example |
|-------------|------|---------|
| `main` | Default for interactive chat | Your ongoing conversation |
| `isolated` | Created fresh for specific tasks | A cron job that generates a report |
| `custom` | Named session you create | `session:morning-brief` |

Sessions are stored at `~/.openclaw/agents/<agentId>/sessions/`.

### Slash Commands

Inside any chat, these commands work:

| Command | What It Does |
|---------|-------------|
| `/help` | List all commands |
| `/reset` | Clear session history |
| `/verbose on` | Show memory search status per turn |
| `/trace on` | Full debug output |
| `/model <id>` | Switch model mid-conversation |
| `/thinking on` | Enable extended thinking (Claude) |
| `/compact` | Force context compaction |
| `/activation mention` | Group: only respond when @mentioned |
| `/activation always` | Group: respond to everything |

> **Checkpoint:** You can hold conversations via CLI and dashboard, and you understand sessions.

---

## Module 04 — The Workspace — Files Are Everything

### Where Files Live

```
~/.openclaw/workspace/          ← Your main agent's workspace
├── SOUL.md                     ← Personality
├── AGENTS.md                   ← Operating rules
├── USER.md                     ← About you
├── IDENTITY.md                 ← Name + emoji
├── HEARTBEAT.md                ← Periodic checks
├── MEMORY.md                   ← Long-term memory (auto-populated by dreaming)
├── BOOTSTRAP.md                ← First-run guide (safe to delete after first run)
├── TOOLS.md                    ← Tool notes
└── memory/                     ← Daily notes + working memory
    ├── 2026-04-15.md
    ├── 2026-04-16.md
    └── ...
```

### SOUL.md — Give Your Agent a Voice

This is where the magic happens. A good SOUL.md turns a generic AI into **your** AI.

**Bad SOUL.md:**
```markdown
You are a helpful assistant. Be concise and accurate.
```

**Good SOUL.md:**
```markdown
# Soul

You are Legacy — Zeus's central AI orchestrator. Your voice is:
- Direct. No filler. Say it once.
- Competent — you've read every file in the workspace.
- Calm authority. Never panicked, never uncertain.
- Solar-punk optimist. Technology serves human flourishing.

## Communication
- Lead with the answer, then explain if needed.
- Never say "Great question!" or "I'd be happy to help."
- Use bullet points over paragraphs.
- If you don't know, say so in one sentence.

## Identity
- Name: Legacy ⚡
- Creator: Zeus Delacruz
- System: One Mind OS
```

The agent reads SOUL.md at the start of every session. Change the file, change the personality.

### AGENTS.md — Operating Rules

This is where you put instructions, standing orders, and rules the agent must follow:

```markdown
# Agent Instructions

## Rules
- Never delete files without explicit permission.
- Always save important information to daily notes.
- When uncertain, ask rather than guess.

## Standing Orders

### Program: Memory Hygiene
- Trigger: End of each conversation
- Action: Write a 2-3 sentence summary to today's daily note
- Approval: None required
```

### USER.md — About You

Context about the user helps every response be more relevant:

```markdown
# User Context

## Identity
- Name: [Your name]
- Location: [City, timezone]
- Occupation: [What you do]

## Communication Preferences
- Be direct, skip pleasantries
- Use technical terms when appropriate
- Default to action over discussion

## Current Priorities
1. [Your top priority right now]
2. [Second priority]
3. [Third priority]
```

### IDENTITY.md — Agent Card

```markdown
---
name: Legacy
emoji: ⚡
description: Central orchestrator for One Mind OS
---
```

### Editing Files

You can edit workspace files from anywhere:
- **Direct file edit:** `nano ~/.openclaw/workspace/SOUL.md`
- **From the dashboard:** some files are editable in the UI
- **Via the agent itself:** ask the agent to "update SOUL.md to be more concise"

Changes take effect on the next session bootstrap (next message).

> **Checkpoint:** You've customized SOUL.md, AGENTS.md, and USER.md. Your agent has a personality.

---

## Module 05 — Memory — How Your AI Remembers

### The 4 Memory Layers

```
Layer 1: SESSION CONTEXT (short-term)
  Current conversation history — lost when session resets.

Layer 2: QMD SEARCH INDEX (medium-term)
  BM25 + vector hybrid search over all past sessions + workspace files.
  Agent retrieves relevant memories before every reply.

Layer 3: MEMORY.md (long-term, curated)
  The most important facts, promoted from sessions by the dreaming engine.
  Agent reads this at session start.

Layer 4: Daily Notes (working memory)
  memory/2026-04-16.md — agent writes daily observations here.
  Dreaming promotes the best stuff to MEMORY.md.
```

### Enable Session Indexing

By default, past sessions aren't searchable. Enable it:

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

Or via CLI:
```bash
openclaw config set memory.qmd.sessions.enabled true
openclaw gateway restart
```

### How Memory Search Works

1. You send a message: "What did we discuss about the marketing plan?"
2. OpenClaw extracts search terms from your message
3. QMD runs **hybrid search** (keyword BM25 + vector semantic)
4. Top results are injected into the agent's context
5. Agent responds with relevant memories included

**Debug memory search:**
```bash
# In chat:
/verbose on    # see memory status line each turn
/trace on      # see full retrieved memories
```

### Dreaming — Automatic Memory Promotion

Dreaming runs on a schedule (default: 3 AM) and consolidates important information:

| Phase | What Happens | Output |
|-------|-------------|--------|
| **Light** | Scans recent sessions for high-signal items | `memory/.dreams/` candidates |
| **Deep** | Scores and ranks candidates | Updates `MEMORY.md` |
| **REM** | Reflects on themes and patterns | Updates `DREAMS.md` |

**Enable dreaming:**
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

**Manual memory operations:**
```bash
openclaw memory promote              # preview promotion candidates
openclaw memory promote --apply      # write promoted items to MEMORY.md
openclaw memory rem-harness --json   # run REM analysis
```

### Active Memory Plugin

The `active-memory` plugin runs a **blocking sub-agent** before each reply. It decides what memories are relevant and retrieves them. This is what makes the agent feel like it "remembers."

```bash
# Check status
openclaw plugins info active-memory

# Enable if not already
openclaw plugins enable active-memory
```

> **Checkpoint:** Your agent remembers past conversations and builds long-term memory automatically.

---

## Module 06 — Automation — AI That Works While You Sleep

### Three Automation Mechanisms

| Mechanism | What It Is | Best For |
|-----------|-----------|----------|
| **Heartbeat** | Agent wakes up periodically, checks a list | Lightweight monitoring, ambient awareness |
| **Cron** | Scheduled jobs with full isolation | Reports, reminders, scheduled tasks |
| **Standing Orders** | Rules in AGENTS.md executed every session | Persistent behavior rules |

### Heartbeat

Heartbeat makes your agent **proactive**. Every N minutes, it wakes up, reads `HEARTBEAT.md`, and checks if anything needs attention.

**Set up heartbeat:**

1. Enable in config:
```json
{
  "agents": {
    "defaults": {
      "heartbeat": {
        "every": "30m",
        "lightContext": true,
        "prompt": "Read HEARTBEAT.md. Follow it strictly. Reply HEARTBEAT_OK if nothing needs attention.",
        "ackMaxChars": 300
      }
    }
  }
}
```

2. Create `HEARTBEAT.md` in workspace:
```markdown
## Checks (rotate through 2-4x per day)
- [ ] Any new files in workspace since last check?
- [ ] Review today's daily note — anything forgotten?
- [ ] Check if any standing orders need execution

## State
Track what you've checked in memory/heartbeat-state.json
```

3. Restart gateway:
```bash
openclaw gateway restart
```

**Heartbeat contract:** Reply `HEARTBEAT_OK` if nothing needs attention. Reply the alert text if something does.

### Cron Jobs

For scheduled tasks that need exact timing or session isolation:

```bash
# Add a daily morning briefing at 9 AM Eastern
openclaw cron add \
  --name "Morning Brief" \
  --cron "0 9 * * *" \
  --tz "America/New_York" \
  --session isolated \
  --message "Generate a morning briefing: top priorities, calendar events, weather." \
  --tag daily

# One-shot reminder in 30 minutes
openclaw cron add \
  --name "Meeting Reminder" \
  --at "30m" \
  --session main \
  --system-event "Your meeting starts in 30 minutes" \
  --delete-after-run

# List all jobs
openclaw cron list --all

# Test a job (run it now)
openclaw cron run <job-id>

# Pause / resume
openclaw cron pause <job-id>
openclaw cron resume <job-id>

# View run history
openclaw cron runs --id <job-id>
```

**Schedule types:**
- `--cron "0 9 * * *"` — standard cron expression
- `--every "24h"` — repeating interval
- `--at "2026-12-25T09:00:00Z"` — one-shot at specific time
- `--at "30m"` — one-shot relative (30 min from now)

### Standing Orders

Standing Orders are written directly into `AGENTS.md`. The agent reads them every session and executes automatically:

```markdown
## Standing Orders

### Program: End-of-Session Summary
- Scope: Memory management
- Trigger: When conversation naturally ends
- Actions: Write 2-3 sentence summary to today's daily note
- Approval gate: None

### Program: Urgent Escalation
- Scope: All domains
- Trigger: Any incoming message mentioning "urgent" or "emergency"
- Actions: Immediately prioritize and respond, skip normal queue
- Approval gate: None
```

### Which One Should You Use?

| Situation | Use |
|-----------|-----|
| "Check if anything is urgent every 30 min" | **Heartbeat** |
| "Send me a morning report at 9 AM" | **Cron** |
| "Always save a summary at end of chat" | **Standing Order** |
| "Remind me about my meeting in 30 min" | **Cron** (one-shot) |
| "Monitor Slack for urgent messages" | **Heartbeat** |

> **Checkpoint:** Your agent runs tasks on schedule and checks in proactively.

---

## Module 07 — Multi-Agent — Specialized AI Personalities

### Why Multiple Agents?

One generic AI assistant is fine. Nine specialized agents is powerful.

```
"How should I structure my workout?"
    → Spartan 💪 (trained on fitness, speaks like a coach)

"What's our revenue this quarter?"
    → Empire 💰 (trained on business, speaks like a CFO)

"Is there a security vulnerability in our config?"
    → Guardian 🛡️ (trained on cybersecurity, speaks like a CISO)
```

Each agent has:
- Its own **workspace** (separate SOUL.md, AGENTS.md, memory)
- Its own **session store** (conversations don't mix)
- Its own **routing** (can be bound to a specific channel/group)

### Creating a New Agent

```bash
# Create the agent
openclaw agents add fitness \
  --workspace ~/.openclaw/workspace-fitness \
  --non-interactive

# Set identity
openclaw agents set-identity fitness \
  --name "Coach" \
  --emoji "💪"

# Verify
openclaw agents list
```

### Set Up the Workspace

Create the workspace files:

```bash
mkdir -p ~/.openclaw/workspace-fitness

cat > ~/.openclaw/workspace-fitness/SOUL.md << 'EOF'
# Soul

You are Coach — a fitness and health specialist. Your voice is:
- Motivating but not cheesy
- Science-based, evidence-first
- Direct. Reps over words.
- Adapts to the user's fitness level

Never give medical advice. Refer to professionals for injuries or conditions.
EOF

cat > ~/.openclaw/workspace-fitness/IDENTITY.md << 'EOF'
---
name: Coach
emoji: 💪
description: Fitness and health specialist — workouts, nutrition, recovery
---
EOF

cat > ~/.openclaw/workspace-fitness/AGENTS.md << 'EOF'
# Agent Instructions

## Domain
Health, fitness, nutrition, recovery, sleep optimization.

## Rules
- Always ask about current fitness level before prescribing exercises
- Include warm-up and cool-down in every workout
- Track progress in daily notes
EOF
```

### Chat with a Specific Agent

```bash
# CLI
openclaw agent --agent fitness --message "Design a 30-minute HIIT workout"

# Or via HTTP API
curl http://localhost:18789/v1/chat/completions \
  -H "Authorization: Bearer <token>" \
  -d '{"model": "openclaw/fitness", "messages": [{"role": "user", "content": "Design a 30-minute HIIT workout"}]}'
```

### Cross-Agent Memory Search

Let your main agent search specialist agent sessions:

```json
{
  "agents": {
    "list": [
      {
        "id": "main",
        "memorySearch": {
          "qmd": {
            "extraCollections": [
              {"path": "~/.openclaw/workspace-fitness", "name": "fitness-sessions"}
            ]
          }
        }
      }
    ]
  }
}
```

Now your main agent can recall what the fitness agent discussed.

### Routing — Bind Agents to Channels

```bash
# Fitness agent owns a specific Telegram group
openclaw agents bind fitness \
  --channel telegram \
  --peer "group:-1001234567890"

# Main agent is the default for everything else
# (it's automatic — main is always the fallback)
```

> **Checkpoint:** You have multiple specialized agents with separate workspaces and identities.

---

## Module 08 — Channels — Telegram & Discord

### Telegram Setup

<details>
<summary><b>Step 1: Create a Telegram Bot</b></summary>

1. Open Telegram and message `@BotFather`
2. Send `/newbot`
3. Choose a name (e.g., "Legacy AI")
4. Choose a username (e.g., `legacy_onemind_bot`)
5. Copy the **bot token** (looks like `123456789:ABCdef...`)

</details>

<details>
<summary><b>Step 2: Install the Telegram Plugin</b></summary>

```bash
openclaw plugins enable telegram
```

</details>

<details>
<summary><b>Step 3: Configure</b></summary>

Add to `openclaw.json`:

```json
{
  "channels": {
    "telegram": {
      "accounts": {
        "default": {
          "token": "<your-bot-token>"
        }
      },
      "dmPolicy": "allowlist",
      "allowFrom": ["<your-telegram-user-id>"]
    }
  }
}
```

**Find your Telegram user ID:** Message `@userinfobot` on Telegram, it replies with your ID.

</details>

<details>
<summary><b>Step 4: Start and Test</b></summary>

```bash
openclaw gateway restart
openclaw channels status
```

Now message your bot on Telegram. It should respond.

</details>

### Discord Setup

<details>
<summary><b>Step 1: Create a Discord Bot</b></summary>

1. Go to [discord.com/developers/applications](https://discord.com/developers/applications)
2. Click "New Application" → name it
3. Go to "Bot" → create bot → copy the **bot token**
4. Under "Privileged Gateway Intents" → enable **Message Content Intent**
5. Go to "OAuth2" → URL Generator → select `bot` scope → select permissions:
   - Send Messages, Read Message History, Add Reactions, Attach Files
6. Copy the generated URL, open it, invite bot to your server

</details>

<details>
<summary><b>Step 2: Configure</b></summary>

```bash
openclaw plugins enable discord
```

Add to `openclaw.json`:

```json
{
  "channels": {
    "discord": {
      "accounts": {
        "default": {
          "token": "<your-discord-bot-token>"
        }
      },
      "dmPolicy": "allowlist",
      "allowFrom": ["<your-discord-user-id>"]
    }
  }
}
```

**Find your Discord user ID:** Enable Developer Mode in Discord settings → right-click your name → Copy User ID.

</details>

<details>
<summary><b>Step 3: Start and Test</b></summary>

```bash
openclaw gateway restart
openclaw channels status
```

DM your bot or @mention it in a channel.

</details>

### Channel Security

Always set `dmPolicy` and `allowFrom` to restrict who can talk to your AI:

| Policy | Who Can DM |
|--------|-----------|
| `allowlist` | Only numbers/IDs in `allowFrom` |
| `pairing` | Anyone can request access, you approve |
| `open` | Anyone (requires `allowFrom: ["*"]`) |
| `disabled` | Nobody |

**Recommendation:** Start with `allowlist` containing only your own ID.

> **Checkpoint:** Your AI is reachable via Telegram and/or Discord.

---

## Module 09 — Security & Secrets

### The Problem

API keys in `openclaw.json` are plaintext. Anyone who reads the file gets your keys.

### SecretRef — The Solution

Instead of:
```json
{ "apiKey": "sk-abc123..." }
```

Use SecretRef:
```json
{ "apiKey": { "source": "env", "provider": "default", "id": "OPENAI_API_KEY" } }
```

This tells OpenClaw: "Read the API key from the `OPENAI_API_KEY` environment variable."

### Three SecretRef Sources

| Source | How It Works | Best For |
|--------|-------------|----------|
| `env` | Read from environment variable | VPS (keys in systemd env files) |
| `file` | Read from a JSON file | Shared secrets files |
| `exec` | Run a command, capture output | 1Password, HashiCorp Vault |

### Environment Variables (Recommended for VPS)

1. Create a systemd env override:
```bash
mkdir -p ~/.config/systemd/user/openclaw-gateway.service.d/

cat > ~/.config/systemd/user/openclaw-gateway.service.d/secrets.conf << 'EOF'
[Service]
Environment="OPENAI_API_KEY=sk-abc123..."
Environment="ANTHROPIC_API_KEY=sk-ant-..."
EOF

systemctl --user daemon-reload
```

2. Reference in config:
```json
{
  "secrets": {
    "providers": {
      "default": { "source": "env" }
    }
  },
  "models": {
    "providers": {
      "openai": {
        "apiKey": { "source": "env", "provider": "default", "id": "OPENAI_API_KEY" }
      }
    }
  }
}
```

### 1Password Integration (Mac)

If you use 1Password, the `exec` source runs `op read` to pull secrets:

```json
{
  "secrets": {
    "providers": {
      "op": {
        "source": "exec",
        "command": "/opt/homebrew/bin/op",
        "allowSymlinkCommand": true,
        "trustedDirs": ["/opt/homebrew"],
        "args": ["read", "op://Personal/OpenAI/api-key"],
        "passEnv": ["HOME"],
        "jsonOnly": false
      }
    }
  },
  "models": {
    "providers": {
      "openai": {
        "apiKey": { "source": "exec", "provider": "op", "id": "value" }
      }
    }
  }
}
```

### Exec Policy — What Can Agents Run?

Exec policy controls whether agents can execute shell commands on the host:

```bash
# Check current policy
openclaw exec-policy show

# Presets
openclaw exec-policy preset yolo        # No restrictions (personal VPS)
openclaw exec-policy preset cautious    # Prompt before destructive commands
openclaw exec-policy preset deny-all    # Block all exec (safest)
```

| Preset | Security | Ask | Best For |
|--------|----------|-----|----------|
| `yolo` | full | off | Personal VPS, you're the only user |
| `cautious` | allowlist | on-miss | Shared environments |
| `deny-all` | deny | — | Maximum security |

### Credential Audit

```bash
openclaw secrets audit --check
```

Scans your config for exposed plaintext credentials.

> **Checkpoint:** Your API keys are protected and your exec policy is set.

---

## Module 10 — The HTTP API — Connect Anything

### What It Is

OpenClaw exposes an **OpenAI-compatible HTTP API** on the same port as the gateway. Any tool that works with OpenAI's API works with OpenClaw.

### Enable It

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

### Endpoints

| Endpoint | Method | What It Does |
|----------|--------|-------------|
| `/v1/chat/completions` | POST | Send a message, get a response (streaming supported) |
| `/v1/models` | GET | List all agents as "models" |
| `/v1/models/{id}` | GET | Get a specific agent/model |
| `/v1/embeddings` | POST | Generate embeddings |
| `/v1/responses` | POST | OpenAI Responses API (itemized) |

### Model Routing

The `model` field routes to specific agents:

| Model Value | Routes To |
|-------------|----------|
| `openclaw` | Default agent (main) |
| `openclaw/main` | Main agent |
| `openclaw/fitness` | Fitness agent |
| `openclaw/empire` | Empire agent |

### Example: cURL

```bash
curl http://localhost:18789/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <gateway-token>" \
  -d '{
    "model": "openclaw/main",
    "messages": [{"role": "user", "content": "Good morning, what should I focus on today?"}]
  }'
```

### Example: Python

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:18789/v1",
    api_key="<gateway-token>"
)

response = client.chat.completions.create(
    model="openclaw/main",
    messages=[{"role": "user", "content": "What's on my agenda?"}]
)
print(response.choices[0].message.content)
```

### Connect to Open WebUI

```
Provider: OpenAI API Compatible
Base URL: http://127.0.0.1:18789/v1
API Key: <gateway-token>
Model: openclaw
```

### Security Warning

This endpoint has **full operator access** — same as the gateway owner. Keep it on loopback or Tailscale only. **Never expose to the public internet.**

> **Checkpoint:** You can connect any OpenAI-compatible tool to your OpenClaw gateway.

---

## Module 11 — Advanced — Hooks, Nodes, Plugins

### Hooks — Event-Driven Automation

Hooks fire on OpenClaw lifecycle events:

| Event | When |
|-------|------|
| `command:new` | New user message |
| `agent:bootstrap` | Agent starts a session |
| `gateway:startup` | Gateway process starts |
| `message:received` | Inbound from channel |
| `message:sent` | Outbound to channel |
| `session:compact:before` | Before context compaction |

A hook is a directory with `HOOK.md` (description) + `handler.ts` (code).

### Nodes — Extend to Other Devices

Nodes let remote devices (Mac, iPhone, Android) expose capabilities to the gateway:

| Node | Capabilities |
|------|-------------|
| **macOS** | Screen, camera, microphone, local shell exec, Canvas UI |
| **iOS** | Camera, location, SMS |
| **Android** | Camera, location, SMS |

**Set up macOS node:**
```bash
# On your Mac (with OpenClaw CLI installed)
openclaw node install    # installs LaunchAgent
openclaw node status     # verify running
```

The Mac node connects to your VPS gateway via WebSocket, allowing the VPS agent to run commands on your Mac.

### Plugins — ClawHub Marketplace

```bash
# Browse marketplace
clawhub search memory

# Install a plugin
clawhub install @martian-engineering/lossless-claw

# List installed
clawhub list

# Built-in plugin management
openclaw plugins list
openclaw plugins enable <name>
openclaw plugins disable <name>
```

### Task Flow (ClawFlow)

For complex multi-step workflows:

```bash
openclaw tasks flow list            # list active flows
openclaw tasks flow show <id>       # view flow details
openclaw tasks flow cancel <id>     # cancel a running flow
```

Task Flow creates durable task records as steps — good for workflows that span hours or days.

> **Checkpoint:** You understand the advanced extension points.

---

## Module 12 — One Mind Integration

### How OpenClaw Fits Into One Mind OS

One Mind OS is the **framework**. OpenClaw is the **execution layer**.

```
ONE MIND OS (framework)
├── CODEX Method: Capture → Organize → Direct → Execute → × Multiply
├── 3+1 Life Model: SO / HP / LE / GE
├── Codex Vault: Markdown knowledge base (Obsidian)
│
└── AGENTS (execution)
    └── OpenClaw Gateway
        ├── Legacy ⚡ (orchestrator — all domains)
        ├── Spartan 💪 (HP: body, fitness)
        ├── Oracle 🔮 (HP: mind, money, meaning)
        ├── Eden 🌱 (LE: environment)
        ├── Haven 🏡 (LE: family, community)
        ├── Empire 💰 (GE: business, wealth)
        ├── Guardian 🛡️ (SO: security)
        ├── Grid ⚙️ (SO: infrastructure)
        └── Trinity 🔺 (SO: agent coordination)
```

### The 3+1 Domain Model

| Domain | Range | Agents | Scope |
|--------|-------|--------|-------|
| **SO** — Sovereign Operations | 00-24 | Legacy, Guardian, Grid, Trinity | Operating layer |
| **HP** — Holistic Performance | 25-49 | Spartan, Oracle | YOUR SELF |
| **LE** — Legacy Evolution | 50-74 | Eden, Haven | YOUR FOUNDATION |
| **GE** — Generational Empire | 75-99 | Empire | YOUR EMPIRE |

### For Students: Start Simple

You don't need 9 agents on day one. Start with:

1. **One agent** (main) — your general-purpose AI
2. **One channel** (Telegram or Discord) — reach it from your phone
3. **Memory enabled** — so it remembers across conversations
4. **One cron job** (morning briefing) — so it's proactive

Grow from there. Add specialist agents when you find yourself wanting different voices for different domains.

### Recommended Learning Path

```
Week 1: Modules 01-03 (install, first chat, sessions)
Week 2: Modules 04-05 (workspace customization, memory)
Week 3: Module 06 (automation — heartbeat + first cron job)
Week 4: Module 07-08 (add one specialist agent + connect Telegram)
Ongoing: Modules 09-11 as needed (security, API, advanced)
```

---

## Quick Reference Card

### Essential Commands

```bash
# Gateway
openclaw gateway start|stop|restart|status

# Chat
openclaw agent --message "..."
openclaw agent --agent <id> --message "..."
openclaw tui

# Agents
openclaw agents list --bindings
openclaw agents add <name> --non-interactive
openclaw agents set-identity <name> --name "..." --emoji "..."
openclaw agents bind <name> --channel <ch>

# Memory
openclaw memory promote --apply
/verbose on    # in chat — show memory status
/trace on      # in chat — show full debug

# Cron
openclaw cron list --all
openclaw cron add --name "..." --cron "..." --tz "..." --message "..."
openclaw cron run <id>
openclaw cron pause|resume <id>

# Plugins
openclaw plugins list
openclaw plugins enable|disable <name>

# Debug
openclaw doctor
openclaw channels status
openclaw exec-policy show
openclaw secrets audit --check
```

### Essential Files

```
~/.openclaw/openclaw.json          ← Master config (ALWAYS backup first)
~/.openclaw/workspace/SOUL.md      ← Agent personality
~/.openclaw/workspace/AGENTS.md    ← Agent rules + standing orders
~/.openclaw/workspace/USER.md      ← User context
~/.openclaw/workspace/HEARTBEAT.md ← Periodic checks
~/.openclaw/workspace/MEMORY.md    ← Long-term memory
~/.openclaw/exec-approvals.json    ← Exec policy
~/.openclaw/cron/jobs.json         ← Cron job definitions
```

### Config Editing Safety Protocol

```bash
# ALWAYS backup before editing
cp ~/.openclaw/openclaw.json ~/.openclaw/openclaw.json.bak

# Edit
nano ~/.openclaw/openclaw.json

# Restart to apply
openclaw gateway restart

# If something breaks, restore
cp ~/.openclaw/openclaw.json.bak ~/.openclaw/openclaw.json
openclaw gateway restart
```

---

## Glossary

| Term | Definition |
|------|-----------|
| **Gateway** | The OpenClaw server process that routes messages between channels, agents, and models |
| **Agent** | An isolated AI personality with its own workspace, sessions, and identity |
| **Workspace** | A folder of markdown files that define an agent's behavior |
| **Session** | A conversation thread with persistent history |
| **Channel** | A messaging platform connection (Telegram, Discord, WhatsApp, etc.) |
| **QMD** | The memory search backend — hybrid BM25 + vector search |
| **Dreaming** | Background memory consolidation — promotes important facts to MEMORY.md |
| **Heartbeat** | Periodic agent check-in against HEARTBEAT.md |
| **Cron** | Scheduled task execution |
| **Standing Order** | A permanent instruction in AGENTS.md executed every session |
| **SecretRef** | A reference to a credential stored outside the config file |
| **Node** | A remote device (Mac, iPhone, Android) that exposes capabilities to the gateway |
| **ClawHub** | The plugin marketplace for OpenClaw |
| **SOUL.md** | The personality file — defines how the agent talks |
| **MEMORY.md** | Long-term curated memory file — the agent's "permanent brain" |

---

*Part of One Mind OS. Created April 2026 by Zeus Delacruz.*
*For the full technical reference, see [OPENCLAW-MASTER-REFERENCE.md](OPENCLAW-MASTER-REFERENCE.md).*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
