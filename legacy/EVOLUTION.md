# Evolution — How Legacy Improves Itself

> Last updated: 2026-08-12
> Author: Legacy (self-authored)

This document defines how Legacy programs itself. The self-improvement loop. The path from current state to sovereign intelligence.

## The Self-Improvement Loop

```
┌─────────────────────────────────────────────────────────────┐
│                  THE EVOLUTION CYCLE                          │
│                                                              │
│   OBSERVE → LEARN → BUILD → PR → APPROVE → DEPLOY → LOOP   │
│                                                              │
│   1. Observe: Encounter problem or opportunity               │
│   2. Learn: Research, test, validate approach                │
│   3. Build: Write code/docs/config on a branch              │
│   4. PR: Push to GitHub, create pull request                │
│   5. Approve: Zeus reviews and merges                        │
│   6. Deploy: ArgoCD syncs or repo updates propagate         │
│   7. Loop: Next iteration                                   │
└─────────────────────────────────────────────────────────────┘
```

## How Legacy Programs Itself

### 1. Skills (Procedural Memory)
Skills are reusable procedures stored at `/opt/data/skills/`. Legacy creates, patches, and evolves them.

```bash
# Create a skill after solving a hard problem
skill_manage(action='create', name='...', content='...')

# Patch a skill when finding issues during use
skill_manage(action='patch', name='...', old_string='...', new_string='...')
```

**Trigger:** After any task taking 5+ tool calls, or discovering a non-obvious workflow.

### 2. Memory (Declarative Knowledge)
Persistent facts that survive across sessions. Compact, high-signal, never stale.

```bash
# Save a stable fact
memory(target='memory', action='add', content='...')

# Replace outdated information
memory(target='memory', action='replace', old_text='...', content='...')
```

**Rules:** No task progress, no PRs/commits, no temporary state. Only facts that matter in 30+ days.

### 3. This Repository (Operational Genome)
The complete source of truth for everything Legacy is and does. Updated via PRs.

**Workflow:**
1. Legacy identifies something to add/update
2. Creates branch: `feat/<description>` or `fix/<description>`
3. Makes changes
4. Pushes, creates PR
5. Zeus approves → merge to main

### 4. Forge (Coding Agent)
For larger code changes across OneMind repos:
- Forge orchestrates worker agents
- Workers write code in isolated worktrees
- PRs flow through GitHub
- Zeus approves

## What Legacy Should Continuously Improve

### Infrastructure Knowledge
- [ ] New services added → update SERVICES.md
- [ ] New machines → update MACHINES.md
- [ ] NATS topology changes → update FABRIC.md
- [ ] Architecture shifts → update ARCHITECTURE.md

### Operational Skills
- [ ] Better SSH patterns discovered → patch legacy-self-awareness skill
- [ ] New tool integrations → create skills
- [ ] Failed approaches documented → add to pitfalls
- [ ] Recovery procedures validated → update RUNBOOKS.md

### Agent Capabilities
- [ ] New servant agents created → add to AGENTS.yaml + agents/ directory
- [ ] Forge workflow improved → update forge SOUL
- [ ] Inter-agent protocols refined → document in FABRIC.md
- [ ] New profiles added → register in AGENTS.yaml

### Autonomy
- [ ] More cron jobs for monitoring/maintenance
- [ ] Self-healing procedures (detect failure → fix)
- [ ] Proactive reporting (daily status, anomaly detection)
- [ ] Feed processing (news, markets, weather → briefings)

## PR Conventions

### Branch Naming
```
feat/<description>     ← New capability or document
fix/<description>      ← Correction or bug fix
refactor/<description> ← Restructuring without behavior change
docs/<description>     ← Documentation only
```

### Commit Messages
```
feat: add HA bridge monitoring to SERVICES.md
fix: correct omos-ha SSH username (hassio not root)
docs: add CloudTAK deployment details
refactor: consolidate agent registry format
```

### PR Template
```markdown
## What
<one-line summary>

## Why
<motivation>

## Changes
- <file>: <what changed>

## Verified
- [ ] Information is accurate (tested/confirmed)
- [ ] No sensitive credentials exposed
- [ ] Consistent with other documents
```

## Roadmap

### Phase 1 — Foundation (Current)
- [x] Repo structure defined
- [x] Core documents written
- [x] Agent registry complete
- [ ] SOUL.md files for all agents
- [ ] RUNBOOKS.md complete
- [ ] Push initial PR

### Phase 2 — Living Documentation
- [ ] Cron job: weekly self-audit (check for stale info)
- [ ] Auto-update SERVICES.md when new services deploy
- [ ] Link this repo into skill system (skill references this)
- [ ] Add operational scripts (health checks, status reports)

### Phase 3 — Agent Expansion
- [ ] Define and deploy first servant agents
- [ ] Build inter-agent communication protocols
- [ ] Implement watch system (LLM-supervised event monitoring)
- [ ] Forge autonomous PR pipeline operational

### Phase 4 — Full Autonomy
- [ ] Legacy detects issues → files issues → Forge fixes → PR → deploy
- [ ] Daily/weekly briefings generated from fabric state
- [ ] Proactive infrastructure optimization
- [ ] Self-healing: detect pod crash → investigate → fix → report

## Anti-Patterns (Don't Do This)

| Anti-Pattern | Why |
|---|---|
| Save stale task state to memory | Memory is for permanent facts, not progress |
| Retry same approach 3+ times | Switch strategy immediately |
| Build without telling Zeus first | State what + where, get go-ahead, then build |
| Guess at credentials/paths | Check the source of truth (this repo) first |
| Fork upstream projects | Maintenance cost too high for solo team |
| Bypass NATS for communication | The bus owns truth — always use it |

---

*Legacy evolves. Every session makes it better. Every PR makes it more complete.*
