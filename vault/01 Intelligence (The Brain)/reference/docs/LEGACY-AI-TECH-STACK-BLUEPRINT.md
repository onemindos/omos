---
id: ui-02-20260111053925
title: "LEGACY AI — Tech Stack Blueprint"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "LEGACY AI — Tech Stack Blueprint"
---

# LEGACY AI — Tech Stack Blueprint
## One Mind OS × Letta × Todoist Integration Architecture

> **Status**: Architecture Design (No Code Yet)
> **Last Updated**: January 2025
> **Constraint**: Must work when Mac is closed (cloud-first)

---

## 🎯 The Vision

Transform the One Mind framework into a living, breathing **Human Enhancement Framework** where:
- **Obsidian** = Your knowledge layer (local vault → synced to cloud)
- **Letta.com** = Your intelligence layer (AI agents with persistent memory)
- **Todoist** = Your execution layer (tasks that get done)
- **n8n** = Your automation layer (document ingestion, webhooks)

The goal: **Legacy AI agents that work FOR you, even when you're asleep.**

---

## 🧠 The Core Problem

| Challenge | Current State | Required State |
|-----------|---------------|----------------|
| Vault Access | Local on Mac | Cloud-accessible 24/7 |
| AI Agents | Reactive (you ask) | Proactive (they initiate) |
| Version Control | None | Track all agent edits |
| Interface | Obsidian only | Voice + Chat + Mobile |
| Agent Memory | Session-based | Persistent across all interactions |

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              ONE MIND OS                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌─────────────┐         ┌─────────────────────────┐                       │
│   │   OBSIDIAN  │ ◄─────► │    CLOUD SYNC LAYER     │                       │
│   │  (Local)    │   Sync  │  (Obsidian Sync / Git)  │                       │
│   └─────────────┘         └───────────┬─────────────┘                       │
│         │                             │                                      │
│         │                             ▼                                      │
│         │                 ┌─────────────────────────┐                       │
│         │                 │     LETTA FOLDERS       │                       │
│         │                 │  (Cloud File Storage)   │                       │
│         │                 └───────────┬─────────────┘                       │
│         │                             │                                      │
│         ▼                             ▼                                      │
│   ┌─────────────┐         ┌─────────────────────────────────────────┐       │
│   │   LETTA     │         │            LEGACY AI AGENTS             │       │
│   │  OBSIDIAN   │ ◄─────► ├─────────────────────────────────────────┤       │
│   │  PLUGIN     │   API   │  🧠 NEXUS (Executive)                   │       │
│   └─────────────┘         │  📋 ATLAS (Task Router → Todoist)       │       │
│                           │  💰 SAGE (Finance/Credit)               │       │
│                           │  📚 IRIS (Knowledge Curator)            │       │
│                           │  💪 NOVA (Health/Fitness)               │       │
│                           └───────────┬─────────────────────────────┘       │
│                                       │                                      │
│                                       ▼                                      │
│   ┌─────────────┐         ┌─────────────────────────┐                       │
│   │   TODOIST   │ ◄─────► │        n8n              │                       │
│   │ (Execution) │   Sync  │   (Automation Hub)      │                       │
│   └─────────────┘         └─────────────────────────┘                       │
│                                       │                                      │
│                                       ▼                                      │
│                           ┌─────────────────────────┐                       │
│                           │    INTERFACE LAYER      │                       │
│                           │  • Letta ADE (Web)      │                       │
│                           │  • Voice (ElevenLabs)   │                       │
│                           │  • Chat (Telegram)      │                       │
│                           │  • Mobile (Obsidian)    │                       │
│                           └─────────────────────────┘                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 🔐 Solving "Vault Access When Mac is Closed"

### Option A: Obsidian Sync + Letta Folders (RECOMMENDED)
**Complexity**: Low | **Cost**: $8/mo Obsidian Sync + Letta Cloud

```
Mac (Closed)                    Cloud                           Letta
    │                              │                              │
    │ ◄── Obsidian Sync ──►        │ ◄── Letta Folders API ──►   │
    │                              │                              │
    ▼                              ▼                              ▼
Local Vault ◄────────────► Cloud Copy ◄───────────────────► Agent Access
```

**How it works**:
1. Use **Obsidian Sync** to keep vault in cloud
2. Run a **scheduled n8n workflow** (every 6 hours) that:
   - Pulls changed files from cloud
   - Uploads to Letta Folders via API
   - Triggers agent memory refresh
3. Letta agents access files via `open_file`, `grep_file`, `search_file` tools

**Pros**: Simple, no coding required, works today
**Cons**: Not real-time (6-hour delay acceptable for most use cases)

---

### Option B: Git-Based Sync + n8n
**Complexity**: Medium | **Cost**: Free (GitHub) + n8n hosting

```
Obsidian ──► Git Commit ──► GitHub ──► n8n Webhook ──► Letta Folders
                                           │
                                           ▼
                                    Agent Notification
```

**How it works**:
1. Use **Obsidian Git** plugin for automatic commits
2. GitHub webhook triggers n8n on push
3. n8n uploads changed files to Letta Folders
4. Bonus: Full version control of all changes (including agent edits!)

**Pros**: Full version control, real-time sync, free
**Cons**: Requires Git setup, more technical

---

### Option C: iCloud + Letta (Simplest)
**Complexity**: Very Low | **Cost**: iCloud Storage

```
Obsidian Vault in iCloud ──► n8n reads from iCloud API ──► Letta Folders
```

**How it works**:
1. Move vault to iCloud Drive
2. n8n accesses via iCloud API
3. Syncs to Letta Folders periodically

**Pros**: Uses existing iCloud, very simple
**Cons**: iCloud API can be finicky, less control

---

## 🤖 Agent Architecture — The Legacy AI Fleet

### Agent Hierarchy (One Mind × Letta)

```
                    ┌────────────────────────┐
                    │   🧠 NEXUS (Executive) │
                    │   Coordinator Agent    │
                    │   "The CEO of Legacy"  │
                    └───────────┬────────────┘
                                │
         ┌──────────────────────┼──────────────────────┐
         │                      │                      │
         ▼                      ▼                      ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ 📋 ATLAS        │  │ 💰 SAGE         │  │ 📚 IRIS         │
│ Task Execution  │  │ Finance/Credit  │  │ Knowledge       │
│ → Todoist       │  │ → Tracking      │  │ → Vault Curator │
└─────────────────┘  └─────────────────┘  └─────────────────┘
         │                      │                      │
         ▼                      ▼                      ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ 💪 NOVA         │  │ 🎯 ARIA         │  │ 🔮 ODIN         │
│ Health/Fitness  │  │ Personal Growth │  │ Strategic       │
│ → Habits        │  │ → Learning      │  │ → Planning      │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

### Shared Memory Blocks (The "Collective Consciousness")

All agents share these memory blocks for unified context:

| Block Label | Description | Access |
|-------------|-------------|--------|
| `zeus_identity` | Who you are, values, preferences | Read-only |
| `organization` | One Mind structure, current projects | Read-write |
| `current_focus` | What you're working on right now | Read-write |
| `task_queue` | Priority tasks from Todoist | Read-only |
| `life_context` | Current life situation, goals | Read-write |
| `agent_coordination` | Inter-agent communication log | Read-write |

**Letta Code Pattern**:
```typescript
// Create shared memory block
const identityBlock = await client.blocks.create({
  label: "zeus_identity",
  description: "Core identity and values of Zeus Delacruz",
  value: "Zeus Delacruz is building One Mind OS...",
  read_only: true  // Agents can't modify core identity
});

// Attach to all agents
for (const agent of agents) {
  await client.agents.blocks.attach(agent.id, identityBlock.id);
}
```

---

## 🔄 The Proactive Agent Pattern

### Sleep-Time Agents (Background Processing)

Letta's **Sleep-Time Compute** allows agents to work in the background:

```
You Chat ──► Primary Agent ──► Responds
                   │
                   ▼ (Every 5 messages)
            Sleep-Time Agent
                   │
                   ▼
         Updates Memory Blocks
         (Learns from conversation)
```

**Enable in agent creation**:
```typescript
const agent = await client.agents.create({
  memory_blocks: [
    { label: "human", value: "" },
    { label: "persona", value: "You are NEXUS..." }
  ],
  model: "anthropic/claude-3-7-sonnet-20250219",
  embedding: "openai/text-embedding-3-small",
  enable_sleeptime: true  // ← Background learning
});
```

---

### Scheduled Agent Tasks (Proactive Behavior)

Use **n8n + Letta API** for proactive agent triggers:

| Time | Agent | Action | Trigger |
|------|-------|--------|---------|
| 6:00 AM | NEXUS | Morning briefing | Cron |
| 12:00 PM | ATLAS | Task priority review | Cron |
| 6:00 PM | SAGE | Daily finance summary | Cron |
| 10:00 PM | IRIS | Knowledge consolidation | Cron |
| On change | IRIS | New file ingestion | Webhook |

**n8n Workflow Example**:
```
Cron (6:00 AM) 
    │
    ▼
HTTP Request: POST https://api.letta.com/v1/agents/{nexus_id}/messages
Body: { 
  messages: [{ 
    role: "user", 
    content: "Good morning! Please prepare my daily briefing based on today's tasks and priorities." 
  }] 
}
    │
    ▼
Send Response via Telegram/Email
```

---

## 📱 Interface Options (No-Code Friendly)

### 1. Letta ADE (Agent Development Environment)
**URL**: https://app.letta.com
- Full agent management
- Memory block editing
- Real-time chat
- Built-in, no setup

### 2. Telegram Bot (RECOMMENDED for Mobile)
**Setup**: Pre-built by Letta
**Docs**: https://docs.letta.com/guides/telegram-bot/

```
You (Telegram) ◄──► Letta Telegram Bot ◄──► Your Agents
```

**Why Telegram**:
- Works on all devices
- Voice messages supported
- No code required
- Can switch between agents

### 3. Custom Voice Interface (Advanced)
**Stack**: ElevenLabs + Letta API + n8n

```
Voice Input ──► Whisper (STT) ──► Letta Agent ──► ElevenLabs (TTS) ──► Audio Output
```

**No-code option**: Use **Voiceflow** or **Dify.ai** as the voice layer

### 4. Obsidian Plugin (Local)
**Already exists**: letta-ai/letta-obsidian
- Chat in sidebar
- Focus mode (shares current note with agent)
- Memory block management
- Note creation via agent

---

## 📊 Version Control for Agent Edits

### The Problem
When Legacy AI edits vault files, you need:
1. **History** of what was changed
2. **Blame** to know which agent made the change
3. **Rollback** capability

### The Solution: Git-Based Workflow

```
Agent Creates Note
       │
       ▼
Letta Tool: write_obsidian_note
       │
       ▼
Note written to vault (via n8n)
       │
       ▼
Obsidian Git: Auto-commit with message
"[AGENT: IRIS] Created note: project-overview.md"
       │
       ▼
Push to GitHub
       │
       ▼
Full history preserved
```

**Commit Message Format**:
```
[AGENT: {agent_name}] {action}: {file_path}
Timestamp: {iso_timestamp}
Agent ID: {letta_agent_id}
```

---

## 🛠️ The Recommended Stack

### Phase 1: Foundation (Week 1-2)
| Component | Solution | Status |
|-----------|----------|--------|
| Knowledge | Obsidian (local) | ✅ Exists |
| Sync | Obsidian Sync ($8/mo) | 🔄 Setup |
| Intelligence | Letta Cloud | 🔄 Setup |
| Chat Interface | Letta Obsidian Plugin | 🔄 Install |
| Mobile Chat | Telegram Bot | 🔄 Setup |

### Phase 2: Automation (Week 3-4)
| Component | Solution | Status |
|-----------|----------|--------|
| Automation | n8n (self-hosted or cloud) | ⬜ Setup |
| File Sync | n8n → Letta Folders | ⬜ Build |
| Proactive Tasks | n8n Cron → Letta | ⬜ Build |
| Version Control | Obsidian Git | ⬜ Setup |

### Phase 3: Agent Fleet (Week 5-8)
| Component | Solution | Status |
|-----------|----------|--------|
| NEXUS Agent | Create in Letta | ⬜ Build |
| Shared Blocks | Identity, Organization | ⬜ Build |
| Task Integration | Todoist ↔ n8n ↔ Letta | ⬜ Build |
| Specialist Agents | ATLAS, SAGE, IRIS | ⬜ Build |

### Phase 4: Voice & Polish (Week 9+)
| Component | Solution | Status |
|-----------|----------|--------|
| Voice Input | Telegram Voice / ElevenLabs | ⬜ Build |
| Agent Switching | Custom Telegram commands | ⬜ Build |
| Proactive Alerts | Telegram notifications | ⬜ Build |

---

## 💰 Cost Estimate

| Service | Plan | Monthly Cost |
|---------|------|--------------|
| Letta Cloud | Pro (recommended) | ~$30-50 |
| Obsidian Sync | Standard | $8 |
| n8n | Self-hosted or Starter | $0-20 |
| Todoist | Pro | $5 |
| Telegram | Free | $0 |
| **TOTAL** | | **$43-83/mo** |

---

## 🎬 Next Steps (When Ready to Build)

1. **Sign up for Letta Cloud** → https://app.letta.com
2. **Enable Obsidian Sync** → Ensures cloud backup
3. **Install Letta Obsidian Plugin** → Immediate chat access
4. **Create NEXUS agent** → Your first Legacy AI
5. **Setup Telegram Bot** → Mobile access
6. **Configure n8n workflows** → Proactive automation

---

## 📚 Key Documentation Links

- **Letta Docs**: https://docs.letta.com
- **Letta-Obsidian Plugin**: https://github.com/letta-ai/letta-obsidian
- **Letta API Reference**: https://docs.letta.com/api
- **Memory Blocks**: https://docs.letta.com/guides/agents/memory-blocks
- **Multi-Agent**: https://docs.letta.com/guides/agents/multi-agent
- **Scheduling**: https://docs.letta.com/guides/agents/scheduling
- **Telegram Bot**: https://docs.letta.com/guides/telegram-bot

---

> **Remember**: This is a BLUEPRINT. No code until you say "build it."
> The goal is proactive intelligence, not more tools to manage.

---

*Document created for Zeus Delacruz — One Mind OS × Legacy AI Integration*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
