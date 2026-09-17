# Hermes vs OpenClaw: Architectural Analysis

> **Classification:** Core Education — AI Agent Architecture  
> **Date:** 2026-05-22  
> **Status:** Verified from source code review  
> **Purpose:** Teach students WHY we chose Hermes as our long-term agent runtime

---

## The #1 Architectural Difference

**Hermes has two loops. OpenClaw has one.**

| | OpenClaw | Hermes |
|---|---|---|
| Conversation loop | ✅ | ✅ |
| Maintenance loop | ❌ | ✅ |
| Self-improving | No | Yes |
| Architecture | Gateway (message router) | Agent OS (self-maintaining) |

OpenClaw runs a turn, delivers the output, and waits for the next input. That's it.

Hermes runs a turn, delivers the output, THEN spawns a background agent that reviews what just happened and updates the skills/memory that feed the NEXT turn's system prompt. On a 7-day cycle, a Curator prunes and archives stale knowledge. The system gets better at its job without human intervention.

---

## Source Code Proof (Verified May 2026)

### Hermes — Two-Loop Architecture

**Repository:** `~/.hermes/hermes-agent/` (NousResearch/hermes-agent)

#### Loop 1: Conversation (`run_agent.py` → `AIAgent.run_conversation()`)
```
User message → Build system prompt (skills injected as text)
    → pre_llm_call hook → LLM API call → post_llm_call hook
    → Tool dispatch (sequential or parallel, max 8 threads)
    → Loop until iteration budget exhausted
    → transform_llm_output hook → Deliver response
```

#### Loop 2: Maintenance (`agent/background_review.py` + `agent/curator.py`)
```
After EVERY turn:
    → Spawn daemon thread with forked AIAgent
    → Restricted to memory + skills tools ONLY
    → Reviews session, updates Markdown skill files
    → Results surfaced via background_review_callback

Every 7 days (when idle > 2h):
    → Curator scans agent-created skills
    → active → stale (30 days no use)
    → stale → archived (90 days)
    → Never deletes — only archives (recoverable)
    → Pinned skills bypass all transitions
```

### OpenClaw — Single-Loop Gateway

**Package:** `openclaw` (npm, compiled TypeScript → 4,160 minified JS files)

```
Inbound message → Channel adapter (plugin)
    → Gateway event bus
    → agent-runner.runtime.js (LLM call + tool dispatch)
    → Deliver response
    → STOP. Wait for next input.
```

**Verified absence:** `grep` across all 4,160 compiled files returns ZERO hits for:
- `selfImprove`, `self_improve`
- `curator`
- `backgroundReview`, `background_review`
- `skillReview`, `skill_review`

There is no maintenance loop. No self-improvement. No skill lifecycle management.

---

## What "Self-Improvement" Actually Means (Technically)

It does NOT mean:
- ❌ The agent modifies its own Python source code
- ❌ `exec()` or `eval()` running dynamic code
- ❌ Hot-reloading plugins at runtime
- ❌ The agent rewriting its own logic

It DOES mean:
- ✅ After each turn, a background agent reviews what happened
- ✅ It creates/updates **Markdown skill files** (`~/.hermes/skills/SKILL.md`)
- ✅ These skills are injected into the **system prompt** at next session start
- ✅ The Curator prunes skills that aren't being used (stale → archived)
- ✅ Skills have usage tracking (`last_activity_at`, `created_at`, state)
- ✅ The net effect: the system prompt evolves over time based on actual usage patterns

**Think of it like this:** The agent writes notes to its future self about what works. Those notes become instructions. Bad notes get cleaned up automatically. Good notes compound.

---

## Kanban Workers: Where It Compounds

Hermes spawns **real OS processes** for task execution:

```python
proc = subprocess.Popen(
    ['hermes', '-p', profile, '--skills', 'kanban-worker', '-q', prompt],
    cwd=workspace,
    start_new_session=True,
)
```

Each worker:
- Runs under a specialized **profile** (coder, researcher, writer, etc.)
- Gets its own skill set loaded
- Has access to the full tool ecosystem
- Reports results back to the Kanban board (SQLite, WAL mode, CAS locking)

The maintenance loop improves the skills → workers use those skills → better results → maintenance loop learns from better results → skills improve further.

**This is a compound interest loop for intelligence.**

OpenClaw has no equivalent. It dispatches a single LLM turn and delivers the response. No workers. No profiles. No skill evolution.

---

## Python vs TypeScript: Honest Assessment

| Factor | Python (Hermes) | TypeScript (OpenClaw) |
|---|---|---|
| AI/ML ecosystem | Native (pyrogram, boto3, etc.) | Requires bridges |
| MCP server ecosystem | Majority are Python (FastMCP) | Growing but smaller |
| LLM code generation quality | Best (most training data) | Good |
| Startup time | Slower | Faster |
| Concurrent I/O | Good (asyncio) | Excellent (native) |
| Type safety | Runtime errors | Compile-time checks |
| Gateway performance | Adequate | Superior |

**Verdict:** The language choice is NOT the fundamental differentiator. The architecture decisions are. The two-loop design COULD have been built in TypeScript — it just wasn't. Python has marginal ecosystem advantages for AI tooling, but the real win is the architectural decision to build a self-maintaining system.

---

## When to Use Each

| Use Case | Best Choice | Why |
|---|---|---|
| Long-term autonomous agent | **Hermes** | Self-improves, spawns workers, compounds |
| Multi-platform message delivery | **OpenClaw** | Purpose-built gateway, WebSocket, nodes |
| Task execution at scale | **Hermes** | Kanban + subprocess workers + profiles |
| Real-time voice/canvas/A2UI | **OpenClaw** | Node architecture, native client support |
| Learning from past sessions | **Hermes** | Background review + Curator |
| Routing to multiple services | **OpenClaw** | Event bus + channel adapters |

---

## Key Lesson for Students

> **An agent that only exists when spoken to is not an agent — it's a service.**

The question isn't "which has more features?" The question is: **does it get better when you're not looking?**

- OpenClaw: Same capability on day 1 and day 365 (unless YOU update it)
- Hermes: Measurably better on day 365 because 365 days of background reviews have refined its skills and pruned its failures

That compound effect is the entire thesis. Not Python. Not features. The feedback loop.

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│  HERMES (Agent OS)                                          │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Loop 1: CONVERSATION                               │    │
│  │  user msg → system prompt (skills) → LLM → tools    │    │
│  │  → response → deliver                               │    │
│  └──────────────────────────┬──────────────────────────┘    │
│                             │ triggers                       │
│  ┌──────────────────────────▼──────────────────────────┐    │
│  │  Loop 2: MAINTENANCE                                │    │
│  │  background_review.py: review turn → update skills  │    │
│  │  curator.py: 7-day prune cycle                      │    │
│  │  skill_usage.py: track what's working               │    │
│  └──────────────────────────┬──────────────────────────┘    │
│                             │ feeds back into               │
│  ┌──────────────────────────▼──────────────────────────┐    │
│  │  SKILLS (Markdown → System Prompt)                  │    │
│  │  Evolves over time. Compounds.                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  KANBAN WORKERS (subprocess.Popen)                  │    │
│  │  Profile-specialized. Independent processes.        │    │
│  │  Use evolved skills. Report back.                   │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  OPENCLAW (Gateway)                                         │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Loop 1: CONVERSATION (only loop)                   │    │
│  │  inbound → channel adapter → LLM → tools → deliver │    │
│  │  → STOP. Wait.                                      │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                             │
│  No maintenance loop. No skill evolution. No workers.       │
│  Static until manually updated.                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Summary (One Sentence)

**Hermes is an agent that maintains itself; OpenClaw is infrastructure that requires you to maintain it.**
