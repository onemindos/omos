---
id: so-01-claude-meta-layer-20260506
title: "Claude — One Mind Meta Layer"
type: reference
status: active
created: 2026-05-06
modified: 2026-05-06
owner: zeus
tags: [type/reference, domain/SO, domain/01-intelligence, agent/legacy, tool/claude-code]
ai_summary: "Design doc and config reference for Claude Code as the meta orchestration layer over One Mind and OpenClaw."
---

# Claude — One Mind Meta Layer

> Claude Code is the meta layer that sits above everything. It handles strategy, cross-domain reasoning, planning, and research — while OpenClaw agents handle persistent execution on the VPS.

**Equation:** `Zeus + Claude (meta) + OpenClaw (execution) = One Mind in motion`

---

## Layer Model

```
┌─────────────────────────────────────────────────┐
│           CLAUDE CODE  (Meta Layer)             │
│   Skills · MCP Tools · Hooks · Memory           │
├───────────────────┬─────────────────────────────┤
│  Vault            │  Execution Layer             │
│  (filesystem MCP) │  OpenClaw — 9 agents (VPS)  │
└───────────────────┴─────────────────────────────┘
       ↕ Todoist  ↕ Zoho  ↕ GitHub  ↕ Circle.so
```

---

## 1. MCP Stack

Config lives at: `~/.claude/mcp.json`

| Server | Status | What It Unlocks |
|--------|--------|-----------------|
| `filesystem` | active | Read/write vault directly |
| `github` | active | PRs, issues, commits |
| `circle` | active | Community posts, members, DMs |
| `tavily` | active | Deep web research |
| `context7` | active | Library/SDK docs lookup |
| `playwright` | active | Browser automation |
| `todoist` | **active** | Read tasks, create tasks, vault ↔ Todoist sync |
| `zoho` | deferred | Email, CRM, Forms, Desk |
| `openclaw-bridge` | **TODO** | Dispatch + observe all 9 OpenClaw agents |

### OpenClaw Bridge (custom MCP — to build)

A lightweight Node.js MCP server wrapping the OpenClaw gateway API.

- **Gateway:** `wss://onemindos-cloud.tail717ba0.ts.net`
- **Token:** stored in 1Password → `op://aaeyris5z44pdqspk277u2oppi/OpenClaw Gateway/token`
- **Tools to expose:**
  - `dispatch_agent(agent_id, message)` — send a task to any of the 9 agents
  - `agent_status(agent_id?)` — read active state / memory of one or all agents
- **Location when built:** `_system/openclaw/mcp-bridge/`

---

## 2. Skills — Vault → Claude Slash Commands

Vault skills live in `SO/02 Intelligence/Skills (Agent Skills)/`. Claude skills live at `~/.claude/skills/` (global) or `.claude/skills/` (project-scoped).

**Priority 1 — Daily Operations** ✅ Live in `.claude/skills/`

| Vault Skill | Claude Command | Pillar | Status |
|-------------|---------------|--------|--------|
| `hp/morning-briefing` | `/morning` | HP | ✅ active |
| `hp/daily-briefing` | `/brief` (universal, time-aware) | HP | ✅ active |
| `hp/evening-review` | `/evening` | HP | ✅ active |
| `so/research` | `/research` | SO | ✅ active |
| `onemind/onemind-architect` | `/architect` | SO | ✅ active |

> **Note:** `/brief` is the universal time-aware skill — auto-delivers morning intel (5am-12pm), midday status (12pm-5pm), or evening reflection (5pm-midnight) based on current time. Use it anytime.

**Priority 2 — Weekly Cadence** ✅ Complete (2026-05-10)

| Vault Skill | Claude Command | Pillar | Status |
|-------------|---------------|--------|--------|
| `hp/weekly-scorecard` | `/scorecard` | HP | ✅ active |
| `onemind/onemind-report` | `/report` | SO | ✅ active |
| `ge/task-management` | `/tasks` | GE | ✅ active |
| `so/signals` | `/signals` | SO | ⚡ VPS only — OpenClaw/Legacy (NATS, HA, Twilio) |

**Priority 3 — On-Demand** ✅ Complete (2026-05-10)

| Vault Skill | Claude Command | Pillar | Status |
|-------------|---------------|--------|--------|
| `hp/life-score-intelligence` | `/lifescore` | HP | ✅ active |
| `ge/project-planning` | `/plan` | GE | — generic, covered natively |
| `ge/writing` | `/write` | GE | — generic, covered natively |
| `so/problem-solving` | `/solve` | SO | — generic, covered natively |
| `so/code-review` | `/review` | SO | — covered by superpowers skill |

---

## 3. Hooks

Config lives at: `~/.claude/settings.json` → `hooks`

| Hook | Trigger | Purpose | Status |
|------|---------|---------|--------|
| `SessionStart` | Every session opens | Load Zeus persona + active goals from vault into context | active (pixel-agents hook) |
| `PreToolUse` | Before any vault write | Validate frontmatter, auto-route to correct domain | ✅ active |
| `PostToolUse` | After task created | Mirror new tasks to Todoist with `tid:` sync | ✅ active |
| `Stop` | Session ends | Log session summary to `_codex/memory/work/` | ✅ active |

---

## 4. Memory System

Claude's persistent memory lives at:
`~/.claude/projects/-Users-zeuslegacy-Documents-onemind/memory/`

Current memory files (auto-maintained by Claude):
- `user_*.md` — Zeus's profile, preferences, working style
- `feedback_*.md` — corrections and confirmed approaches
- `project_*.md` — active projects, decisions, timelines
- `reference_*.md` — pointers to external systems

This memory bridges sessions so Claude never starts cold in the One Mind context.

---

## 5. Agent Roster (OpenClaw — for dispatch reference)

| ID | Agent | Domain | Workspace |
|----|-------|--------|-----------|
| `main` | ⚡ Legacy | Orchestrator — all pillars | `~/.openclaw/workspace` |
| `spartan` | 💪 Spartan | HP 27-28 — body, fitness | `~/.openclaw/workspace-spartan` |
| `oracle` | 🔮 Oracle | HP 25-26, 29-34 — mind, money | `~/.openclaw/workspace-oracle` |
| `eden` | 🌱 Eden | LE — sustainability | `~/.openclaw/workspace-eden` |
| `haven` | 🏡 Haven | LE 50-74 — family, home | `~/.openclaw/workspace-haven` |
| `empire` | 💰 Empire | GE 75-99 — ventures, brand, wealth | `~/.openclaw/workspace-empire` |
| `guardian` | 🛡️ Guardian | SO 04 — security | `~/.openclaw/workspace-guardian` |
| `grid` | ⚙️ Grid | SO 01-02 — infrastructure | `~/.openclaw/workspace-grid` |
| `trinity` | 🔬 Trinity | SO 06 — research | `~/.openclaw/workspace-trinity` |

**Gateway:** `wss://onemindos-cloud.tail717ba0.ts.net`
**VPS:** `onemindos-cloud-1` / `64.227.8.23` / Tailscale `100.98.25.10`

---

## 6. Implementation Roadmap

### Phase 1 — Connect Tools ✅ Complete (2026-05-06)
- [x] Add Todoist MCP to `~/.claude/mcp.json` — token from `~/.secrets/onemind.env`
- [ ] Add Zoho MCP — deferred
- [x] GitHub MCP — wired via `gh` CLI OAuth token in `.claude/mcp.json`

### Phase 2 — Port Skills ✅ Complete (2026-05-06)
- [x] `/morning` — morning briefing (`morning-briefing` vault skill)
- [x] `/brief` — universal time-aware brief (morning/midday/evening auto-detected)
- [x] `/evening` — evening review (`evening-review` vault skill)
- [x] `/research` — deep research via Tavily + Context7
- [x] `/architect` — One Mind architecture review + audit

Skills location: `.claude/skills/` in vault root

### Phase 3 — Build OpenClaw Bridge
- [ ] Scaffold MCP server at `_system/openclaw/mcp-bridge/`
- [ ] Implement `dispatch_agent` tool
- [ ] Implement `agent_status` tool
- [ ] Add to `~/.claude/mcp.json`
- [ ] Test: Claude → Legacy → VPS → response back

### Phase 4 — Wire Hooks ✅ Complete (2026-05-06)
- [x] `PreToolUse` hook: frontmatter validator (warning on missing fields)
- [x] `PostToolUse` hook: Todoist mirror — auto-creates tasks + writes `tid:` back
- [x] `Stop` hook: session summary → `_codex/memory/work/YYYY-MM-DD.md`

Scripts: `_system/hooks/` | Registered: `.claude/settings.json` (project-scoped)

### Phase 5 — Meta Layer v2 ✅ Complete (2026-05-09)
- [x] `lib/runtime.js` — unified hook intelligence layer (getVaultContext, logEvent, writeMemory, writeSessionLog)
- [x] `lib/todoist.js` — extracted Todoist logic, fixed API path to `/rest/v2/tasks`
- [x] `UserPromptSubmit` — vault context injected into every prompt (date, inbox count, recent work)
- [x] `SubagentStart` + `SubagentStop` — agent activity observability with pillar mapping
- [x] `ContextCompaction` — state saved to memory before summarization fires
- [x] `SessionEnd` — rich session closeout log to `_codex/memory/sessions/`
- [x] `PostToolUseFailure` — tool error logging to event log + memory
- [x] `stop.js` refactored to use runtime (eliminated duplicate logic)
- [x] 9 OpenClaw-mirrored agents: Legacy, Spartan, Oracle, Eden, Haven, Empire, Guardian, Grid, Trinity
- [x] Permissions consolidated — broad wildcards, clean allowlist (~75 rules → clean)

Scripts: `_system/hooks/` | Agents: `.claude/agents/` | Registered: `.claude/settings.json`

### Phase 6 — Claude Code Skills ✅ Complete (2026-05-10)
- [x] `/tasks` — Todoist tasks grouped by pillar (SO/HP/LE/GE) + flat mode
- [x] `/sync` — vault ↔ GitHub ↔ VPS sync status with graceful SSH fallback
- [x] `/report` — dashboard (fast) + `/report deep` (full pillar breakdown)

### Phase 7 — OpenClaw Bridge ✅ Complete (2026-05-10)
- [x] `/dispatch` skill — sends messages to any agent via gateway OpenAI-compatible API
- [x] Gateway endpoint: `https://omos-core.tail717ba0.ts.net/v1/chat/completions`
- [x] Model routing: `openclaw/{agent_id}` (e.g. `openclaw/main`, `openclaw/spartan`)
- [x] Auth: gateway token (`Z6vJ49...`) in openclaw.json — `gateway.auth.token`
- [x] `chatCompletions` endpoint enabled: `gateway.http.endpoints.chatCompletions.enabled = true`
- [x] **LIVE:** Tested end-to-end Mac → Tailscale → VPS gateway → agent → response

No custom MCP server needed — gateway already speaks OpenAI protocol over Tailscale.

### Phase 8 — OpenClaw Agent Reporting ✅ Complete (2026-05-10)
- [x] `onemind-session-logger` hook pack — `_system/openclaw/hooks/onemind-session-logger/`
  - Fires on `message:sent` (every agent reply via Discord/Telegram channels)
  - Appends turn entry to `_codex/memory/agents/{name}/sessions/YYYY-MM-DD-{id}.md`
  - Installed via `openclaw plugins install` → listed as 6th hook, status ✓ ready
- [x] Cron deep reports — `_system/openclaw/hooks/setup-cron-reports.sh` — 8 jobs registered
  - Legacy: daily 11pm ET | Grid: daily 10pm ET | Spartan/Oracle/Empire: Mon 8am ET
  - Guardian: Wed 9am ET | Haven: Sun 9am ET | Eden: 1st of month 9am ET
- [x] Report files → `_codex/memory/agents/{name}/reports/YYYY-MM-DD.md`
- [x] onemind-watch auto-commits all session/report files

Scripts: `_system/openclaw/hooks/` | Design doc: `Agents/_runtimes/openclaw/AGENT-REPORTING.md`

### Phase 9 — Full Loop
- [ ] Morning: `/morning` pulls Todoist + vault state → Zeus gets brief
- [ ] During day: `/research`, `/plan`, `/dispatch` agents as needed
- [ ] Evening: `/evening` → logs to vault → syncs to VPS → agents pick it up
- [ ] Weekly: `/scorecard` + `/report` → full domain health view

---

## Related
- [[01 Intelligence]] — parent domain
- `_system/openclaw/` — OpenClaw agent configs
- `~/.claude/settings.json` — Claude harness config
- `~/.claude/mcp.json` — MCP server registry
- `AGENTS.md` — universal agent instructions
