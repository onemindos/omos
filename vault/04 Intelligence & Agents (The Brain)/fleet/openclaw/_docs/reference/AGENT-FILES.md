# Agent Files — How They Work

Every OpenClaw agent has a workspace folder at `~/.openclaw/workspace-{id}/` (Legacy uses `~/.openclaw/workspace/`). These files are the agent's brain at startup and across sessions. Each file has a specific load time and purpose.

---

## Workspace File Reference

| File | Loaded When | Purpose |
| --- | --- | --- |
| `IDENTITY.md` | Every session (bootstrap) | Who the agent is — name, role, pillar, personality |
| `SOUL.md` | Every session (bootstrap) | Core directives, operating principles, non-negotiables |
| `USER.md` | Every session (bootstrap) | Who Zeus is — communication style, preferences, timezone, priorities |
| `AGENTS.md` | Every session (bootstrap) | Full fleet roster — who owns what, how to route |
| `TOOLS.md` | Every session (bootstrap) | Infrastructure map — server paths, services, key credentials references |
| `HEARTBEAT.md` | Cron job only | Template for scheduled health check messages — not loaded in conversations |
| `MEMORY.md` | Session start (if present) | Persistent memory context — current state, priorities, lessons learned |
| `memory/` directory | On demand | Timestamped session memory files — searched when agent needs historical context |
| `skills/` directory | On demand (skill invocation) | Domain skills — loaded only when the skill is triggered by keyword or command |

---

## Load Sequence — Every Session

When a new session starts, OpenClaw bootstraps the agent with the workspace files in this order:

```
1. IDENTITY.md      ← Who am I?
2. SOUL.md          ← What do I stand for?
3. USER.md          ← Who am I talking to?
4. AGENTS.md        ← Who else is on the team?
5. TOOLS.md         ← What infrastructure do I have access to?
6. MEMORY.md        ← What do I remember from before?
```

All 6 files go into the system prompt / context window at session start. They are **always present** for every conversation.

---

## What Does NOT Load Every Session

| What | When It Loads | Why |
| --- | --- | --- |
| `HEARTBEAT.md` | Only when a heartbeat cron job runs | It's a cron template, not conversational context |
| `memory/*.md` (timestamped files) | On demand via memory search | Too large to load every time — searched when needed |
| Skills | When triggered by keyword, `/skill`, or cron `--skill` flag | Each skill adds to context only when relevant |

---

## How Skills Work Per Turn

Skills are not passive — they activate dynamically:

```
User turn (or cron message) arrives
    ↓
OpenClaw scans for skill triggers:
    - Keyword match in message (e.g. "morning brief" → morning-briefing skill)
    - Explicit /skill command
    - Cron job with --skill flag
    ↓
If triggered: skill's SKILL.md injected into context for that turn
    ↓
Agent responds with skill's domain knowledge active
    ↓
Next turn: skill context may or may not persist (depends on session type)
```

**Skills are on-demand context injections.** They don't sit in the base context — they activate when needed, keeping the base context clean and fast.

### Skill file structure (each skill has):

```
skill-name/
├── SKILL.md       ← The actual instructions injected into context
├── README.md      ← Human-readable description
└── config.json    ← Trigger keywords, activation rules
```

---

## Memory Architecture

Three tiers of memory, each with different persistence and load behavior:

```
Tier 1: Working Memory (session-scoped)
    → Lives in the active context window
    → Cleared when session ends or /new is called
    → Fastest — already in the context

Tier 2: Episodic Memory (workspace/memory/*.md files)
    → Persisted timestamped files in workspace/memory/
    → Loaded on demand — agent searches when it needs historical context
    → Survives gateway restarts
    → Written by the session-memory hook when /new or /reset runs

Tier 3: Semantic Memory (vault — _codex/memory/)
    → Lives in the onemind vault
    → Permanent — survives everything
    → Agent reads via MCP codex-vault tool
    → Contains Zeus preferences, routing history, relationship context
```

---

## Session Types and Context Behavior

| Session Type | Context Scope | Use Case |
| --- | --- | --- |
| `isolated` | Fresh context per run — workspace files + cron prompt only | All cron jobs (heartbeats, reports) |
| `main` | Persistent conversation thread — carries full history | Direct conversations with Zeus |
| `fork` | Inherits parent session context | Subagents spawned from a parent task |

> All cron jobs run as `isolated` sessions — this is why heartbeats don't bleed context from the main conversation thread.

---

## File Locations (omos-core)

| Agent | Workspace Path |
| --- | --- |
| ⚡ Legacy | `~/.openclaw/workspace/` |
| 🏔 Zenith | `~/.openclaw/workspace-zenith/` |
| 🌿 Unity | `~/.openclaw/workspace-unity/` |
| 👑 Eternal | `~/.openclaw/workspace-eternal/` |
| 🌐 Sovereignty | `~/.openclaw/workspace-sovereignty/` |
| 👁️ Sentinel | `~/.openclaw/workspace-sentinel/` (omos-sentinel) |

Each agent also has a deeper path at `~/.openclaw/agents/{id}/workspace/` — this is the internal agent state used by the gateway. The workspace folder at the root level is what you edit.

---

## Vault Copies

Every workspace file has a vault counterpart kept in sync:

```
onemind/
└── 00-24 SO/
    └── 01 Intelligence/
        └── Agents (Digital arms)/
            └── _roster/_COMMAND/{agent}/
                └── openclaw/
                    ├── IDENTITY.md
                    ├── SOUL.md
                    ├── TOOLS.md
                    ├── MEMORY.md
                    ├── AGENTS.md
                    ├── USER.md
                    └── HEARTBEAT.md
```

**Edit the vault copy → git push → omos-core pulls within 10s** (via `onemind-pull.timer`).
Or edit directly on the VPS for immediate effect — then push to sync vault.

---

## Quick Edit Reference

```bash
# Edit a workspace file directly on VPS (takes effect next session)
ssh onemindos-core 'nano ~/.openclaw/workspace/SOUL.md'

# Check what's loaded for an agent
ssh onemindos-core 'ls ~/.openclaw/workspace/'

# View memory files
ssh onemindos-core 'ls ~/.openclaw/workspace/memory/'

# Trigger a skill manually (in agent chat)
/skill morning-briefing

# List all available skills
openclaw skills list
```

---

*Last updated: 2026-05-14*
*Related: [MODEL-ROUTING.md](MODEL-ROUTING.md) · [QUICK-REFERENCE.md](QUICK-REFERENCE.md)*
