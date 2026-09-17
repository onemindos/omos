# Multi-Agent Setup — OneMind Hive

## Overview

The OpenClaw gateway runs **4 isolated agents** (one per OneMind domain). Each agent has its own workspace, session store, auth profiles, and model assignment. They can communicate with each other via agent-to-agent messaging.

```
HUMAN (Zeus)
    │
    ▼
OpenClaw Gateway (:18789)
    │
    ├── legacy  (UI domain — orchestrator, Discord, BlueBubbles)
    ├── hp      (25-49 Holistic Performance)
    ├── le      (50-74 Legacy Evolution)
    └── ge      (75-99 Generational Endeavor)
```

---

## Agent Configuration

### legacy (Default Agent)

| Property | Value |
|----------|-------|
| ID | `legacy` |
| Default | Yes |
| Model | `openai/gpt-4o` |
| Workspace | `/home/onemind/.openclaw/workspace-legacy` |
| Agent Dir | `/home/onemind/.openclaw/agents/legacy/agent` |
| Domain | UI — 00-24 Sovereign Operations |
| Channel Routing | Discord, BlueBubbles |

**Skills loaded:**
Agent Browser, api-gateway, desktop-control, find-skills, github, obsidian,
proactive-agent, skill-creator, weather, 1password, apple-notes, apple-reminders,
bear-notes, bluebubbles, camsnap, discord, gemini, gh-issues, himalaya, imsg,
openai-image-gen, openai-whisper, slack, spotify-player, things-mac, tmux,
trello, voice-call, wacli, xurl (and more — 50+ total)

---

### hp (Holistic Performance Coach)

| Property | Value |
|----------|-------|
| ID | `hp` |
| Name | `HP \| Holistic Performance` |
| Model | `anthropic/claude-sonnet-4-5` |
| Workspace | `/home/onemind/.openclaw/workspace-hp` |
| Agent Dir | `/home/onemind/.openclaw/agents/hp/agent` |
| Domain | 25-49 HP — Body, Mind, Money, Law, Relationships, Daily Ops, Joy |

**SOUL.md personality:** Direct, evidence-based, zero fluff. Sports coach + life strategist. Tracks PI Score (Personal Intelligence).

---

### le (Legacy Evolution Steward)

| Property | Value |
|----------|-------|
| ID | `le` |
| Name | `LE \| Legacy Evolution` |
| Model | `anthropic/claude-sonnet-4-5` |
| Workspace | `/home/onemind/.openclaw/workspace-le` |
| Agent Dir | `/home/onemind/.openclaw/agents/le/agent` |
| Domain | 50-74 LE — Family, Home, Estate, Generational Wealth, Faith |

**SOUL.md personality:** Warm but strategic. 50-year horizon. Family office + estate planner + historian. Tracks LI Score (Legacy Intelligence).

---

### ge (Generational Endeavor Advisor)

| Property | Value |
|----------|-------|
| ID | `ge` |
| Name | `GE \| Generational Endeavor` |
| Model | `anthropic/claude-opus-4-6` |
| Workspace | `/home/onemind/.openclaw/workspace-ge` |
| Agent Dir | `/home/onemind/.openclaw/agents/ge/agent` |
| Domain | 75-99 GE — Business, Revenue, Clients, Ventures, Investments |

**SOUL.md personality:** Sharp, deal-oriented, ROI-first. McKinsey + startup founder. Tracks WI Score (Wealth Intelligence).

---

## Channel Routing (Bindings)

```json
"bindings": [
  { "agentId": "legacy", "match": { "channel": "discord" } },
  { "agentId": "legacy", "match": { "channel": "bluebubbles" } }
]
```

All Discord and BlueBubbles messages → `legacy` agent.
HP, LE, GE are accessible via the gateway UI session picker or by configuring additional bindings (e.g. separate WhatsApp numbers, Telegram bots).

---

## Agent-to-Agent Messaging

Enabled — all 4 agents can message each other:

```json
"tools": {
  "agentToAgent": {
    "enabled": true,
    "allow": ["legacy", "hp", "le", "ge"]
  }
}
```

---

## Subagent Configuration (Global Defaults)

```json
"subagents": {
  "maxConcurrent": 6,
  "maxSpawnDepth": 2,
  "maxChildrenPerAgent": 4,
  "model": "anthropic/claude-sonnet-4-5"
}
```

---

## Workspace Files

Each workspace contains (or will contain):
- `SOUL.md` — Agent personality, identity, tone
- `AGENTS.md` — Operating instructions, tool usage rules
- `USER.md` — Zeus's profile
- `MEMORY.md` — Long-term memory
- `memory/YYYY-MM-DD.md` — Daily session notes
- `skills/` — Per-agent custom skills

SOUL.md files are already written for hp, le, ge.
Copy AGENTS.md, USER.md from workspace-legacy to other workspaces to bootstrap them.

---

## CLI: Managing Agents

```bash
# List all agents and bindings
openclaw agents list --bindings

# Add a new agent interactively
openclaw agents add <id>

# Check agent status via gateway
openclaw gateway health

# Send a message to a specific agent (via gateway session)
openclaw gateway call sessions.send --agent hp --message "What is my PI Score?"
```

---

## Model Defaults

```json
"defaults": {
  "model": "anthropic/claude-opus-4-6",
  "models": {
    "anthropic/claude-opus-4-6": { "alias": "opus" },
    "anthropic/claude-sonnet-4-5": { "alias": "sonnet" },
    "openai/gpt-4o": { "alias": "gpt4o" }
  }
}
```

Switch model mid-session: type `/model sonnet` or `/model opus` in chat.

---

## Memory / Compaction

```json
"compaction": {
  "memoryFlush": {
    "enabled": true,
    "prompt": "Write any lasting notes to memory/YYYY-MM-DD.md; reply with NO_REPLY if nothing to store."
  }
}
```

Context is compacted automatically. Important context is flushed to daily memory files.

---

## Adding More Agent Bindings

### Route a specific WhatsApp DM to HP agent

```json
{
  "agentId": "hp",
  "match": {
    "channel": "whatsapp",
    "peer": { "kind": "direct", "id": "+1XXXXXXXXXX" }
  }
}
```

### Route a Telegram bot to GE agent

```json
{
  "agentId": "ge",
  "match": { "channel": "telegram", "accountId": "ge-bot" }
}
```

Add new Telegram bot:
```bash
channels.telegram.accounts.ge-bot.botToken = "YOUR_BOT_TOKEN"
```
