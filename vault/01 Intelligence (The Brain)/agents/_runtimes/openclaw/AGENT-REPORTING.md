---
id: so-01-agent-reporting-20260510
title: "Agent Reporting System"
type: reference
status: active
created: 2026-05-10
modified: 2026-05-10
owner: zeus
tags: [type/reference, domain/SO, domain/01-intelligence, agent/legacy, tool/openclaw]
ai_summary: "Design spec for the OpenClaw agent session logging system — per-turn auto-log via Stop hook, deep reports via cron, session files committed to vault by codex-watch."
---

# Agent Reporting System

> Every OpenClaw agent automatically documents its work. Zero agent awareness required — pure infrastructure.

**Status:** Ready to deploy — scripts built, needs wiring on VPS.

## Activation (One-Time Setup on VPS)

```bash
# 1. Pull latest vault on VPS
cd /root/onemind && git pull

# 2. Copy stop hook to OpenClaw hooks dir
cp _system/openclaw/hooks/stop.js /root/.openclaw/hooks/stop.js

# 3. Wire hook in openclaw.json — add under "gateway":
#    "hooks": { "command:stop": { "script": "~/.openclaw/hooks/stop.js", "timeout": 5000 } }

# 4. Restart gateway
systemctl --user restart openclaw-gateway.service

# 5. Set up cron deep reports
bash _system/openclaw/hooks/setup-cron-reports.sh

# 6. Verify
openclaw cron list
```

Full install guide: `_system/openclaw/hooks/INSTALL.md`

---

## Architecture

```
Agent turn executes
    ↓
OpenClaw Stop hook fires (VPS-side)
    ↓
Appends turn summary to session file
    ↓
codex-watch detects file change (3s debounce)
    ↓
git commit → GitHub → Mac pulls
    ↓
_codex/memory/agents/{name}/sessions/YYYY-MM-DD-{session-id}.md
```

Deep reports run on a separate cron schedule and write to:
```
_codex/memory/agents/{name}/reports/YYYY-MM-DD.md
```

---

## Session File Format

**Location:** `_codex/memory/agents/{agent-name}/sessions/YYYY-MM-DD-{session-id-8}.md`

**Frontmatter:**
```yaml
---
agent: spartan
session_id: abc1234ef
pillar: HP
started: 2026-05-10T21:30:00Z
turn_count: 3
status: active | complete
---
```

**Body — one section per turn, session summary appended at end:**
```markdown
## Turn 1 — 21:30
Task: pull Strava week 2 data
Result: 4 workouts found, 2 PRs flagged
Files: 25-49 HP/28 Body/strava/2026-05-W2.md

## Turn 2 — 21:35
Task: analyze recovery metrics
Result: HRV trending down — flagged for Zeus
Files: none

## Turn 3 — 21:41
Task: write weekly summary
Result: summary written to domain README
Files: 25-49 HP/28 Body/README.md

## Session Summary
Completed weekly HP body analysis. Key flag: HRV declining — may need deload week.
```

---

## Deep Report Format

**Location:** `_codex/memory/agents/{agent-name}/reports/YYYY-MM-DD.md`

**Triggered by:** Cron (agent-specific schedule) or `/report deep` command.

**Content:** Agent reads its own session logs + current domain READMEs + active tasks → synthesizes a full domain health report. Not a raw log — a curated intelligence document.

---

## Directory Structure

```
_codex/memory/agents/
├── legacy/
│   ├── sessions/          ← per-session turn logs
│   └── reports/           ← deep domain reports
├── spartan/
│   ├── sessions/
│   └── reports/
├── oracle/
├── eden/
├── haven/
├── empire/
├── guardian/
├── grid/
└── trinity/
```

---

## Implementation — What Gets Built (Phase 6)

### VPS Stop Hook (per-turn auto-log)
- Script: `/root/.openclaw/hooks/stop.js` (mirrors Claude Code's stop.js pattern)
- Fires after every agent turn via OpenClaw hook config
- Appends structured turn entry to session file
- Creates session file on first turn (with frontmatter)
- Updates `turn_count` and `status: complete` on final turn

### Cron Deep Reports (scheduled)
- Each agent has a systemd timer (agent-specific cadence)
- Default cadences:
  | Agent | Cadence | Reason |
  |-------|---------|--------|
  | Legacy | Daily | Orchestrator — most activity |
  | Spartan | Weekly (Mon) | Fitness week review |
  | Oracle | Weekly (Mon) | Finance/habit week review |
  | Empire | Weekly (Mon) | Business week review |
  | Grid | Daily | Infra changes daily |
  | Guardian | Weekly | Security posture |
  | Eden | Monthly | Sustainability slow-moving |
  | Haven | Weekly | Family/home |
  | Trinity | On-demand | Research when triggered |

### codex-watch (already active)
- No changes needed — already commits all vault file changes
- Session files land in `_codex/memory/agents/` which is inside the vault

---

## Permission Model

| Action | Who can do it | How |
|--------|---------------|-----|
| Write session log | Any agent | Via Stop hook — automatic |
| Write deep report | Any agent | Via cron — automatic |
| Delete own session logs | Any agent | Direct |
| Delete another agent's logs | Legacy only | With explicit permission |
| Delete + escalation | Other agents | Request Legacy → Legacy executes |

**Git is the safety net.** All writes are committed. Any deletion can be reverted. Agents can act boldly.

---

## Chain of Thought Tracing

Session files are the connective tissue. Because every turn in a session is in one file, you can:
- Read one file and see exactly how the agent's thinking evolved
- Diff sessions to see how behavior changes over time
- Cross-reference session_id across agents to see coordinated work
- Feed session files to `/report deep` for instant domain intelligence

---

## Related
- [[CLAUDE-META-LAYER]] — Phase 6 roadmap
- `_system/openclaw/mcp-bridge/` — Phase 3 bridge (prerequisite)
- `_codex/memory/work/` — Claude Code session logs (parallel system, same pattern)
- `AGENTS.md` — universal agent instructions
