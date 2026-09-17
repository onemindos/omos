---
title: "Skills vs. Subagents: Architecture Decision Record"
domain: SO
subdomain: "01 Sovereign Core (The Brain)"
type: architecture-decision
status: active
created: 2026-05-12
updated: 2026-05-12
owner: zeus
tags: [domain/so, type/architecture, status/active, teaching/example]
ai_summary: "Architecture decision record analyzing three approaches to agent specialization in OpenClaw: monolithic skills, subagent-spawnable skills, and dedicated persistent agents. Includes tradeoff analysis, current platform limitations, and the recommended hybrid pattern. Designed as a teaching reference for OneMind students."
---

# Skills vs. Subagents: Architecture Decision Record

**Maintained by:** ⚡ Legacy  
**Decision date:** 2026-05-12  
**Status:** Accepted — Hybrid approach  

---

## Context

OneMind OS originally ran **10+ specialized agents** (Spartan, Oracle, Guardian, Grid, Empire, Haven, Eden, Trinity, Overwatch, etc.), each with its own identity, memory, and cron schedule. This caused:

- Fleet sprawl — too many agents eating resources
- Coordination overhead — all routing through Legacy
- Cron pile-ups — concurrent LLM calls causing event-loop starvation
- Maintenance burden — each agent needing heartbeats, daily reports, workspace files

The fleet was collapsed to **4 domain agents** (Zenith, Unity, Eternal, Sovereignty) + Legacy (orchestrator) + Sentinel (watchdog). The old specialist agents became **skills** loaded on demand by domain agents.

This document records the analysis of how to structure specialization going forward.

---

## The Three Options

### Option A: Fat Skills on Domain Agents (Current Approach)

Domain agents load skill files on demand. Example: Zenith loads `spartan`, `oracle`, `librarian` as needed.

**How it works:**
```
Zeus asks about fitness →
  Legacy routes to Zenith →
    Zenith loads spartan skill →
      SKILL.md instructions inject into Zenith's context →
        Zenith responds with Spartan-level expertise
```

**Pros:**
- Single session — no context overhead, no spawn latency
- Skill instructions inject directly into the agent's context window
- Agent maintains full memory and personality continuity
- Cheapest — one LLM call per task
- Agent retains USER.md, SOUL.md, IDENTITY.md — knows Zeus's preferences

**Cons:**
- Context window pressure — loading 3-4 skills simultaneously gets heavy
- Skills are static instructions — they don't *reason*, they're reference docs
- A finance skill inside Zenith still thinks like Zenith, not like a financial specialist
- Single-threaded — agent can't do two skill-heavy tasks in parallel

**Best for:** Most daily work. Quick lookups. Tasks where domain context matters more than deep specialization.

---

### Option B: Subagent-Spawnable Specialist Skills

Create skills designed to be spawned as isolated subagents — temporary workers for specific tasks.

**How it works:**
```
Zeus asks about credit repair →
  Legacy routes to Zenith →
    Zenith spawns subagent with task: "Analyze credit report, use finance skill" →
      Subagent runs in isolated session →
        Completes work, announces result back to Zenith →
          Zenith synthesizes and responds to Zeus
```

**Pros:**
- Clean context isolation — subagent only has task-relevant context
- Parallel execution — spawn credit repair + budgeting analysis simultaneously
- Cheaper model for routine work (e.g., Haiku 4.5 for data crunching)
- Each subagent gets a focused, single-purpose context window
- No permanent resource cost — subagent archives after completion

**Cons:**
- **Critical limitation (as of May 2026):** Subagents only inject AGENTS.md + TOOLS.md into context. No SOUL.md, IDENTITY.md, USER.md, or HEARTBEAT.md. They don't know Zeus's preferences, risk tolerance, or communication style.
- Extra token cost per spawn (new context window each time)
- No memory persistence between subagent runs — they're ephemeral
- Latency — spawn + LLM cold start + announce back adds seconds
- Cannot talk to each other directly — all coordination goes through parent
- `context: "fork"` copies the full parent transcript (expensive)
- `context: "isolated"` (default) means subagent has zero history

**Best for:** Parallel research. Long-running tasks that would block the main conversation. Computation-heavy work where a cheap model suffices.

---

### Option C: Dedicated Persistent Agents

Bring back specialist agents as full entries in `agents.list[]` — Finance Agent, Health Agent, Education Agent, etc.

**How it works:**
```
Zeus asks about credit repair →
  Legacy routes to Finance Agent →
    Finance Agent has full SOUL.md, IDENTITY.md, USER.md →
      Persistent session with memory of past financial conversations →
        Responds with deep financial expertise
```

**Pros:**
- Full identity stack (SOUL.md, IDENTITY.md, USER.md, HEARTBEAT.md, MEMORY.md)
- Persistent sessions with cross-conversation memory
- Can be directly routed to from Telegram topics
- Deepest specialization — the agent *is* the specialist
- Can have its own cron jobs, scheduled reports, proactive monitoring

**Cons:**
- Fleet sprawl — what we already experienced and abandoned
- Each agent needs its own heartbeats, daily reports, maintenance
- Resource cost — more concurrent sessions, more memory, more cron jobs
- Coordination overhead — all inter-agent communication through Legacy
- Diminishing returns — most tasks don't need this depth of specialization

**Best for:** Domains that truly need persistent context across conversations, their own scheduled operations, and deep autonomous behavior. Use sparingly.

---

## Decision: Hybrid Architecture

### The Pattern

```
┌─────────────────────────────────────────────────────────┐
│                    Legacy (Orchestrator)                  │
│              Routes all requests from Zeus                │
└────────────┬──────────┬──────────┬──────────┬───────────┘
             │          │          │          │
        ┌────▼───┐ ┌───▼────┐ ┌──▼─────┐ ┌──▼──────────┐
        │ Zenith │ │ Unity  │ │Eternal │ │ Sovereignty  │
        │  (HP)  │ │  (LE)  │ │  (GE)  │ │    (SO)      │
        └────┬───┘ └───┬────┘ └──┬─────┘ └──┬──────────┘
             │         │         │           │
        Skills loaded on demand (Tier 1 — primary work)
        ┌────┴────────┐
        │  spartan    │  ← Deep skill with references/
        │  oracle     │     health.md, fitness.md, etc.
        │  finance    │
        │  librarian  │
        └─────────────┘
             │
        Subagents spawned for parallel work (Tier 2 — when needed)
        ┌────┴────────┐
        │  Research   │  ← Isolated, cheap model, ephemeral
        │  Analysis   │     Used for parallelism, not specialization
        │  Data pull  │
        └─────────────┘
```

### Rules

1. **Skills handle specialization.** Build deep skills with progressive disclosure (reference files). The skill provides the expertise; the domain agent provides the identity and context.

2. **Subagents handle parallelism.** Spawn subagents when you need to do multiple things at once, run a long task without blocking, or use a cheap model for grunt work.

3. **Don't use subagents for specialization** — they lack identity context (no SOUL.md, USER.md). A spawned "finance subagent" doesn't know Zeus's risk tolerance or financial goals unless you encode it all in the task prompt every time.

4. **Don't create new persistent agents** unless the domain genuinely needs persistent cross-conversation memory AND its own scheduled operations AND autonomous behavior. The bar is high.

---

## The Deep Skill Pattern (Recommended)

Instead of many narrow skills, build fewer deep skills with progressive disclosure:

```
finance/
├── SKILL.md              ← Core workflow, decision routing (~200 lines)
│                            "If credit question → read references/credit.md"
│                            "If budget question → read references/budgeting.md"
└── references/
    ├── credit.md          ← Credit repair procedures, score factors
    ├── budgeting.md       ← Budget frameworks, envelope method, etc.
    ├── investing.md       ← Investment analysis, allocation models
    ├── tax.md             ← Tax planning strategies
    └── accounts.md        ← Account structure, bank relationships
```

**Why this works:**
- SKILL.md metadata (~100 tokens) is always in context → triggers correctly
- SKILL.md body (~500 lines max) loads only when triggered → lean
- Reference files load only when the specific sub-domain is needed → surgical
- One skill covers an entire domain without context bloat
- The domain agent (Zenith) retains its full identity while accessing deep expertise

**Compare to the alternatives:**
- ❌ 5 separate skills (credit, budgeting, investing, tax, accounts) → metadata for all 5 always in context, triggers compete
- ❌ 1 massive skill with everything inline → context window explosion when loaded
- ❌ 5 subagents for each sub-domain → no identity, high latency, expensive
- ✅ 1 deep skill with reference files → clean, efficient, extensible

---

## When to Use Each Pattern

| Scenario | Use |
|----------|-----|
| Zeus asks about fitness | Zenith loads `spartan` skill → direct response |
| Zeus wants a budget AND credit analysis | Zenith loads `finance` skill → spawns 2 subagents for parallel research → synthesizes |
| Zeus wants daily health tracking | Zenith + `spartan` skill + cron job (no subagent needed) |
| Zeus wants a 3-source research brief | Zenith spawns 3 research subagents (cheap model) → synthesizes results |
| Zeus wants ongoing financial monitoring with alerts | Consider a dedicated persistent agent (high bar) |
| Zeus asks a quick question | Domain agent handles directly, no skill or subagent needed |

---

## Subagent Configuration Reference

```jsonc
// In openclaw.json — agents.defaults.subagents
{
  "model": "github-copilot/claude-haiku-4.5",  // Cheap default for workers
  "thinking": "medium",
  "runTimeoutSeconds": 120,
  "archiveAfterMinutes": 60,
  "maxSpawnDepth": 1,  // Set to 2 for orchestrator pattern
  "allowAgents": ["*"]  // Allow spawning under any agent
}
```

**Key settings:**
- `model` — Use a cheap model for subagent workers. Save expensive models for domain agents.
- `maxSpawnDepth` — Default 1 (agents can spawn subagents). Set to 2 for agents to spawn orchestrators that spawn workers.
- `archiveAfterMinutes` — Auto-cleanup. Subagents are ephemeral.
- `context: "isolated"` (default) — Clean child session. Use this 90% of the time.
- `context: "fork"` — Copies parent transcript. Use only when child needs conversation history.

---

## Context Inheritance Matrix (as of May 2026)

| File | Domain Agent | Subagent (isolated) | Subagent (fork) | Persistent Agent |
|------|-------------|-------------------|-----------------|-----------------|
| SOUL.md | ✅ | ❌ | ❌ | ✅ |
| IDENTITY.md | ✅ | ❌ | ❌ | ✅ |
| USER.md | ✅ | ❌ | ❌ | ✅ |
| AGENTS.md | ✅ | ✅ | ✅ | ✅ |
| TOOLS.md | ✅ | ✅ | ✅ | ✅ |
| HEARTBEAT.md | ✅ | ❌ | ❌ | ✅ |
| MEMORY.md | ✅ | ❌ | ❌ | ✅ |
| Conversation history | ✅ | ❌ | ✅ (copied) | ✅ |
| Skills (triggered) | ✅ | ✅ | ✅ | ✅ |
| Persistent memory | ✅ | ❌ | ❌ | ✅ |

**This is the key constraint.** Subagents are capable workers but lack personality, user awareness, and cross-conversation memory. This is why specialization belongs in skills, not subagents.

---

## Future Evolution

As OpenClaw evolves, subagents will likely gain:
- Selective context injection (inherit USER.md without full transcript)
- Persistent subagent sessions (memory across spawns)
- Direct inter-subagent communication (skip parent routing)
- Skill-aware spawning (auto-load relevant skills into subagent context)

When these capabilities arrive, the subagent-spawnable specialist pattern becomes significantly more powerful. The hybrid architecture is designed to evolve into that future without requiring a full restructure.

---

## Summary

| Approach | Specialization | Parallelism | Identity | Cost | Recommendation |
|----------|---------------|-------------|----------|------|----------------|
| Deep skills on domain agents | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | **Primary — use for most work** |
| Subagents for parallel work | ⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐ | ⭐⭐⭐ | **Secondary — use for parallelism** |
| Dedicated persistent agents | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ | **Rare — high bar to justify** |

**The rule: Skills for depth. Subagents for breadth. Persistent agents only when both are needed simultaneously with autonomous behavior.**

---

## Related Documents

- [[agent-skills-architecture]] — Skill tiers, loading, and assignment matrix
- [[AGENT-TEAM-STRUCTURE]] — Current fleet roster and routing
- [[OPERATOR-VS-AUTONOMOUS-AGENTS]] — Autonomy levels and decision authority
- [[SOUL]] — Agent identity and personality framework
