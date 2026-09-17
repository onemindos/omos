# OpenClaw Multi-Agent Capabilities
*Captured: 2026-05-12 | Source: Live session with Legacy*

---

## What You Can Do ✅

**Parallel spawning** — spawn all 5 agents simultaneously in one turn. Legacy fires them all, they run concurrently, results stream back as completion events to the single Legacy session.

**Delegation with context** — two modes:
| Mode | How | Token Cost |
|------|-----|------------|
| `isolated` (default) | Agent starts fresh, reads its own SOUL/IDENTITY/workspace files | Low |
| `fork` | Agent gets full current transcript context | High — multiplied per agent |

**Skills** — each agent loads its own skills on demand per run (Zenith loads Spartan/Oracle, Unity loads Haven/Eden, etc.)

**Single pane** — all results surface back to Legacy. Zeus sees everything from one session.

---

## Fleet Limits (current config)

| Limit | Value | Notes |
|-------|-------|-------|
| Max children per agent | 6 | Can spawn all 5 fleet agents + 1 spare |
| Max concurrent subagents | 12 | No bottleneck for the fleet |
| Max spawn depth | 2 | Fleet agents can spawn sub-subs (one level deep) |
| Run timeout | 600s (10 min) | Raised from 300s — handles complex tasks |
| Archive after | 60 min | Sessions clean up after 1hr idle |

---

## Q&A — Architectural Decisions

### Q: Can Zenith call Unity if it needs home context to complete a task?
**Yes — as of 2026-05-12 this is now enabled.**

Fleet agents (Zenith, Unity, Eternal, Sovereignty) each have `allowAgents` set to all other fleet agents. So Zenith CAN spawn Unity as a sub-subagent to get home context, wait for Unity's response, then return the combined result to Legacy.

**Flow example:**
```
Legacy → spawns Zenith (task: HP decision needing home data)
         Zenith → spawns Unity (task: get current home sensor/schedule data)
                  Unity → returns home context
         Zenith → uses context to complete decision
         Zenith → returns result to Legacy
Legacy → delivers to Zeus
```

**Constraints:**
- Max spawn depth is 2, so Unity cannot spawn a further subagent from inside that call
- Sentinel cannot spawn others (watchdog only)
- No circular spawning (Zenith cannot re-spawn Zenith)

---

### Q: Are subagents ephemeral — does every new turn mean new context?
**Yes — subagents are ephemeral by design.**

Each `sessions_spawn` is a fresh run:
- The agent starts with its workspace files (SOUL.md, IDENTITY.md, AGENTS.md, etc.)
- It does NOT carry memory from a previous spawn unless it's written to persistent storage (Todoist, Obsidian vault, memory files)
- After the run completes it archives after 60 min idle

**How to persist context across turns:**
| Method | Use for |
|--------|---------|
| Write to Obsidian vault | Decisions, summaries, analysis |
| Write to Todoist | Tasks, action items, status |
| Write to `_codex/memory/` | Relationship context, standing state |
| Write to workspace files | Agent-level config/preferences |

If Zenith analyzes something important, it should write to memory before returning — otherwise it's gone.

---

### Q: Can the 5-min timeout be raised?
**Yes — raised to 10 min (600s) as of 2026-05-12.**

For tasks longer than 10 min, break into sequential delegations rather than one long run. The timeout is a safety rail — runaway agents get killed if they exceed it.

---

## Agent Spawn Allowlist (current)

| Agent | Can Spawn |
|-------|-----------|
| Legacy (main) | unity, zenith, eternal, sovereignty, sentinel |
| Zenith | unity, eternal, sovereignty, sentinel |
| Unity | zenith, eternal, sovereignty, sentinel |
| Eternal | unity, zenith, sovereignty, sentinel |
| Sovereignty | unity, zenith, eternal, sentinel |
| Sentinel | *(none — watchdog only)* |

---

## Recommended Patterns

### Pattern 1: Parallel Morning Brief
```
Legacy spawns → [Zenith, Unity, Eternal, Sovereignty] simultaneously
Each checks domain Todoist + priority state
All return → Legacy synthesizes → delivers to Zeus
~60s total
```

### Pattern 2: Cross-Domain Decision (Zenith needs home data)
```
Legacy spawns → Zenith (with task context)
Zenith spawns → Unity (sub-query for home context)
Unity returns → Zenith completes decision
Zenith returns → Legacy delivers
~90s total
```

### Pattern 3: Sequential Pipeline
```
Legacy spawns → Sovereignty (infra check)
Sovereignty writes result to shared file
Legacy spawns → Eternal (uses infra result to plan deployment)
Eternal returns → Legacy delivers
```

---

## Known Constraints

1. **No direct agent-to-agent messaging** — coordination goes through spawn/return only
2. **Fork is expensive** — forking full transcript to 5 agents = ~5x token burn, use isolated unless context is truly needed
3. **Ephemeral by default** — agents must explicitly write to persist state
4. **Depth limit = 2** — can't chain A→B→C→D, only A→B→C
5. **Sentinel is read-only spawner** — Sentinel monitors but cannot dispatch fleet agents
