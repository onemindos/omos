---
type: reference
domain: SO-06
status: active
created: 2026-05-14
updated: 2026-05-14
owner: zeus
tags: [type/reference, domain/SO, domain/06-openclaw, openclaw/memory, plugin/memory-core, plugin/memory-wiki, plugin/active-memory]
ai_summary: "The complete OpenClaw memory stack — memory-core (long-term dreaming), memory-wiki (topic index), active-memory (working RAM). How they work, how they chain, cron schedule, and what survives restarts."
---

# OpenClaw Memory Stack

Three memory plugins. Three layers. Together they form a complete RAM → disk → index chain — the same way human memory actually works.

---

## The Stack at a Glance

```
DURING SESSION          EVERY NIGHT (3AM ET)     ON RECALL
───────────────         ────────────────────      ──────────
active-memory           memory-core               memory-wiki
(RAM)           ──────► (writes to disk)  ──────► (indexes disk
                        "dreaming"                for fast lookup)
```

| Plugin | Analogy | Persists? | When It Runs | Purpose |
| ------ | ------- | --------- | ------------ | ------- |
| `active-memory` | RAM | ❌ Session only | Every turn | Working memory — what was said earlier this chat |
| `memory-core` | REM sleep / disk write | ✅ Survives restarts | Cron `0 3 * * *` (3AM ET) | Long-term durable memory — accumulates over time |
| `memory-wiki` | Encyclopedia index | ✅ Indexed from core | On recall queries | Structures core memory into navigable topic entries |

---

## Plugin 1: `active-memory` — Working Memory (RAM)

**What it does:** Tracks context within the current conversation. Injects relevant earlier-in-session context so the agent doesn't lose the thread mid-conversation.

**Analogy:** Short-term working memory — what's in your head right now, in this conversation.

**Persistence:** Ephemeral. Does not persist between sessions independently. When the session ends, active-memory's state is what gets handed off to memory-core for consolidation.

**When it runs:** Every turn — always active during a session.

**What it stores:** Recent messages, key facts mentioned in this session, decisions made, entities referenced, tasks in progress.

**Limits:** Session-scoped only. If you close and reopen a chat, active-memory starts fresh. The durable record lives in memory-core.

---

## Plugin 2: `memory-core` — Long-Term Memory (Disk / Dreaming)

**What it does:** Runs a "dreaming" consolidation pass over recent sessions at 3AM ET every night. Extracts important facts, decisions, entities, and patterns. Writes them to a persistent structured store that survives agent restarts, config changes, and version upgrades.

**Analogy:** REM sleep — the brain's nightly process of consolidating short-term experiences into long-term memory. The agent "dreams" about what happened today and decides what to keep.

**Cron schedule:** `0 3 * * *` — fires once per night at 3:00 AM ET.

**Persistence:** ✅ Fully durable. This is the primary long-term memory store. Survives restarts, crashes, config reloads.

**Resource cost:** Zero embeddings. Zero GPU/CPU cost during the day. All compute happens at 3AM in the background — the dreaming pass is the only expensive operation, and it's scheduled when the system is idle.

**What it extracts and stores:**
- Facts about Zeus (preferences, decisions, patterns)
- Project and domain context (active work, status, goals)
- Relationship context (people mentioned, roles, history)
- Agent learnings (what worked, what was corrected)
- Recurring topics and priorities

**Why this is the primary memory:** It accumulates over time. Every night adds to it. After weeks and months, the agent has a deep, rich memory of everything that's been discussed — without any manual work.

---

## Plugin 3: `memory-wiki` — Memory Index (Encyclopedia)

**What it does:** Organizes and surfaces memories as wiki-style topic entries. When an agent needs to recall something, it queries this layer like a mini-encyclopedia. It structures what memory-core stores into navigable, topic-indexed entries so recall is fast and relevant.

**Analogy:** The index of a book — or a personal Wikipedia built from your own history.

**Persistence:** ✅ Indexed from memory-core. Updates as memory-core consolidates new facts.

**When it runs:** On recall queries — when an agent needs to surface relevant memories for a conversation.

**Structure:** Topic-indexed. Instead of searching raw memory, the agent can look up: "Zeus / fitness / injuries" or "project / onemind-app / status" and get a structured entry back.

**What it surfaces:**
- Person profiles (Zeus, team members, contacts)
- Project histories and current status
- Domain-specific accumulated knowledge
- Past decisions and their outcomes
- Patterns and learnings over time

---

## How the Stack Chains Together

### During a Session (Active)
```
User message arrives
    ↓
active-memory checks: what was said earlier in this chat?
    ↓
Relevant context injected into agent's working context
    ↓
Agent responds with full conversation thread awareness
```

### Nightly Consolidation (3AM ET)
```
memory-core cron fires: 0 3 * * *
    ↓
Scans recent sessions (since last run)
    ↓
Dreaming pass: extracts facts, decisions, entities, patterns
    ↓
Writes structured entries to persistent store
    ↓
memory-wiki index updated with new topic entries
```

### On Recall (Query Time)
```
Agent needs to remember something
    ↓
memory-wiki queried: "what do I know about [topic]?"
    ↓
Relevant structured entries returned
    ↓
Injected into response context
```

---

## Memory Survival Matrix

| Event | active-memory | memory-core | memory-wiki |
| ----- | ------------- | ----------- | ----------- |
| Session ends | ❌ Lost | ✅ Consolidated at 3AM | ✅ Indexed |
| Agent restart | ❌ Lost | ✅ Survives | ✅ Survives |
| Config change | ❌ Lost | ✅ Survives | ✅ Survives |
| OpenClaw update | ❌ Lost | ✅ Survives | ✅ Survives |
| Server reboot | ❌ Lost | ✅ Survives | ✅ Survives |

**Rule:** If it matters, it needs to survive until 3AM for memory-core to consolidate it. Anything told to an agent that doesn't make it to memory-core is gone at session end.

---

## Operational Notes

**Check memory-core cron on VPS:**
```bash
openclaw cron list | grep memory
```

**Verify plugins are active:**
```bash
openclaw agent show main | grep -A5 plugins
```

**Force immediate consolidation (skip waiting for 3AM):**
```bash
openclaw cron run <memory-core-job-id>
```

**Memory store location on VPS:**
```
~/.openclaw/workspace/memory/     ← Legacy's memory
~/.openclaw/workspace-zenith/memory/   ← Zenith's memory
~/.openclaw/workspace-unity/memory/    ← Unity's memory
# etc. per agent workspace
```

---

## Design Principle

> "active-memory (RAM) → memory-core dreaming (writes to disk) → memory-wiki (indexes the disk for fast lookup)"

This mirrors human memory architecture intentionally:
- **active-memory** = working memory / prefrontal cortex — what you're actively thinking about
- **memory-core dreaming** = hippocampal consolidation during sleep — moving short-term to long-term
- **memory-wiki** = semantic memory / long-term knowledge store — facts and concepts you can look up

The agent doesn't need to be told the same things repeatedly. Over time, memory-core builds a complete picture of Zeus, One Mind, and everything discussed — automatically, every night at 3AM.

---

*See also: `AGENT-FILES.md` (workspace file load order), `MODEL-ROUTING.md` (which models run memory cron jobs)*
