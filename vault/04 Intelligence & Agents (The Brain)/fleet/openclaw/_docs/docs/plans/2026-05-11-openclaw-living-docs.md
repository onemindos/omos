# OpenClaw Living Documentation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Mirror the live omos-core OpenClaw state into complete, accurate documentation in the vault's `02 openclaw/` folder — adding only what's missing, never removing what exists.

**Architecture:** Pull live config from omos-core via SSH, write structured markdown docs into the correct vault folders. Three layers per section: inventory (what exists), architecture (how it works), runbooks (how to operate). Workspace files from the VM backup repo are mirrored into `workspaces/` for rebuild reference.

**Tech Stack:** SSH to omos-core (`ubuntu@100.123.126.34`), `openclaw.json` as source of truth, vault path `/Users/zeuslegacy/Documents/onemind/00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)/`

---

## Key Facts (live state as of 2026-05-11)

**Active agents in `openclaw.json`: 6**
| ID | Name | Model | Domain Skills |
|---|---|---|---|
| `main` | Legacy | default (github-copilot/claude-sonnet-4.6) | coding-agent, github, diffs, clawhub, 1password, oracle, healthcheck, taskflow, tavily, summarize, openai-whisper, video-frames, nano-pdf, camsnap, obsidian, wiki-maintainer, weather, voice-call, telegram, onemind-maintainer, onemind-inbox, onemind-sync |
| `zenith` | Zenith | Sonnet 4.6 [gh] | agent-zenith, spartan, oracle, identity, librarian, pulse + base |
| `unity` | Unity | Sonnet 4.6 [gh] | agent-unity, haven, helix, eden, nomad + base |
| `eternal` | Eternal | Sonnet 4.6 [gh] | agent-eternal, empire, empire-ops, creed, trinity, apex, finance + base |
| `sovereignty` | Sovereignty | Sonnet 4.6 [gh] | agent-sovereignty, grid, guardian, guardian-ops, overwatch, eagleeye, coding-agent, github, diffs + base |
| `sentinel` | Sentinel | default | onemind-maintainer, taskflow, tavily, summarize, telegram, onemind-inbox, onemind-sync |

**Base skills (all agents):** onemind-maintainer, taskflow, taskflow-inbox-triage, tavily, summarize, telegram, onemind-inbox, onemind-sync

**VPS paths:**
- Config: `/home/ubuntu/.openclaw/openclaw.json`
- Workspaces: `/home/ubuntu/.openclaw/workspace/` (Legacy) + `/home/ubuntu/.openclaw/workspace-{id}/`
- Backup repo: `github.com/OneMind-OS/openclaw-core-workspaces`
- Watch script: `/home/ubuntu/.openclaw/watch.sh` (runs as background process)
- Plugin skills: `/home/ubuntu/.openclaw/plugin-skills/`

**What already exists (do NOT modify):**
- `cron (Scheduled Jobs)/CRON.md` — complete, accurate
- `hooks (Event Hooks)/HOOKS.md` — complete, accurate
- `gateway/` — multiple architecture docs, complete
- `help/` — multiple runbooks, complete
- `automation/` — complete
- `install/` — complete
- `channels/` — complete
- `cli/` — full CLI reference, complete
- `workspaces (Agent Workspaces)/` — partial (most agents only have IDENTITY.md)

---

## Task 1: Fix `plugin-skills/PLUGIN-SKILLS.md` — update to 6-agent reality

**Files:**
- Modify: `plugin-skills (Custom Skills)/PLUGIN-SKILLS.md`

Current file says "19 agents share the same base skill set" — wrong. There are 6 agents. Domain expertise is embedded via skills into the 5 super-agents (zenith owns spartan/oracle/identity/librarian/pulse, etc.).

**Step 1: Rewrite PLUGIN-SKILLS.md with live data**

Replace the entire content with:

```markdown
---
title: "OpenClaw — Custom Plugin Skills"
type: reference
status: active
created: 2026-05-10
modified: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO-02, openclaw/skills]
---

# OpenClaw — Custom Plugin Skills

Custom skills installed on omos-core at `~/.openclaw/plugin-skills/`. Loaded into agent sessions to give agents specialized domain knowledge. Skills are plain markdown `SKILL.md` files — no code, just context.

Live path: `/home/ubuntu/.openclaw/plugin-skills/`  
Source: `github.com/OneMind-OS/openclaw-core-workspaces`

---

## Agent → Skills Map

### Legacy (`main`) — Orchestrator
Routing, coding, vault ops, health monitoring.

| Skill | Purpose |
|-------|---------|
| `coding-agent` | Full coding agent capability |
| `github` | GitHub repo operations |
| `diffs` | Diff and patch files |
| `clawhub` | ClawHub skill distribution |
| `session-logs` | Read/write session logs |
| `gh-issues` | GitHub issue management |
| `1password` | 1Password credential lookup |
| `oracle` | Oracle agent delegation |
| `healthcheck` | Node health checks |
| `taskflow-inbox-triage` | Inbox triage via Todoist |
| `openai-whisper` | Voice transcription |
| `video-frames` | Video frame extraction |
| `nano-pdf` | PDF parsing |
| `camsnap` | Camera snapshot |
| `obsidian` | Obsidian vault control |
| `wiki-maintainer` | Wiki knowledge base maintenance |
| `weather` | Weather data |
| `voice-call` | Twilio voice calls |

### Zenith — HP 25-49 (Body, Mind, Money, Spirit, Intelligence)

| Skill | Domain |
|-------|--------|
| `agent-zenith` | Zenith identity + operating context |
| `spartan` | HP 28 — Body, fitness, training |
| `oracle` | HP 25-34 — Mind, habits, personal finance |
| `identity` | HP 25-26 — Purpose, values, spirit |
| `librarian` | HP 27 — Intelligence pipeline, research |
| `pulse` | HP 27 — Real-time feed monitoring |

### Unity — LE 50-74 (Family, Home, Land, Estate, Generations)

| Skill | Domain |
|-------|--------|
| `agent-unity` | Unity identity + operating context |
| `haven` | LE 50-74 — Family, relationships, home |
| `helix` | LE 52-56 — Parenting, elder care, generational |
| `eden` | LE — Sustainability, environment |
| `nomad` | LE 50 — Off-grid, homestead, sovereign living |

### Eternal — GE 75-99 (Empire, Brand, Capital, Ventures, Wealth)

| Skill | Domain |
|-------|--------|
| `agent-eternal` | Eternal identity + operating context |
| `empire` | GE 75-99 — Business strategy, ventures |
| `empire-ops` | GE — Day-to-day ops, brand building |
| `creed` | GE 75-81 — Brand narrative, storytelling |
| `trinity` | SO 06 — Deep R&D, technology analysis |
| `apex` | GE 79 — Capital, fundraising, finance |
| `finance` | GE 79 — Financial modeling, analysis |

### Sovereignty — SO 00-24 (Infra, Security, Automation, Dev, Research)

| Skill | Domain |
|-------|--------|
| `agent-sovereignty` | Sovereignty identity + operating context |
| `grid` | SO 02 — VPS, servers, networking |
| `guardian` | SO 04 — Cybersecurity, access control |
| `guardian-ops` | SO 04 — Security operations |
| `overwatch` | SO 07 — n8n, automation workflows |
| `eagleeye` | SO 02 — Drone ops, computer vision |
| `coding-agent` | Full coding agent capability |
| `github` | GitHub operations |
| `diffs` | Diff and patch files |

### Sentinel — Watchdog (omos-sentinel)

| Skill | Purpose |
|-------|---------|
| `onemind-maintainer` | Vault read/write |
| `taskflow` | Todoist task management |
| `tavily` | Web search |
| `summarize` | Summarize documents |
| `telegram` | Send Telegram messages |
| `onemind-inbox` | Process inbox items |
| `onemind-sync` | Git sync vault changes |

---

## Base Skills (all agents)

Every agent loads these regardless of domain:

| Skill | Purpose |
|-------|---------|
| `onemind-maintainer` | Read/write vault files, maintain domain docs |
| `onemind-inbox` | Process inbox items, route captures to domains |
| `onemind-sync` | Git sync — pull/push vault changes |
| `taskflow` | Todoist task management |
| `taskflow-inbox-triage` | Inbox triage via Todoist (super-agents) |
| `tavily` | Web search |
| `summarize` | Summarize documents/pages |
| `telegram` | Send Telegram messages |

---

## All Available Plugin Skills on VPS

```
~/.openclaw/plugin-skills/
├── agent-eternal/     ├── agent-legacy/      ├── agent-sentinel/
├── agent-sovereignty/ ├── agent-unity/       ├── agent-zenith/
├── apex/              ├── browser-automation/├── creed/
├── eagleeye/          ├── eden/              ├── empire/
├── empire-ops/        ├── finance/           ├── grid/
├── guardian/          ├── guardian-ops/      ├── haven/
├── helix/             ├── homestead/         ├── identity/
├── librarian/         ├── nomad/             ├── obsidian-vault-maintainer/
├── om-dashboard/      ├── onemind-inbox/     ├── onemind-maintainer/
├── onemind-report/    ├── onemind-status/    ├── onemind-sync/
├── oracle/            ├── overwatch/         ├── pulse/
├── spartan/           ├── trinity/           └── wiki-maintainer/
```

Total: 37 skills available, subset loaded per agent.

---

## Managing Skills

```bash
# List all available skills
ssh ubuntu@100.123.126.34 "openclaw skills list"

# Check which skills an agent has loaded
ssh ubuntu@100.123.126.34 "openclaw skills list --agent main"

# Install a skill from ClawHub
ssh ubuntu@100.123.126.34 "openclaw skills install <skill-id>"

# Reload gateway after openclaw.json skill changes
ssh ubuntu@100.123.126.34 "kill -HUP \$(pgrep -f 'openclaw.*gateway')"
```

Config location: `~/.openclaw/openclaw.json` → `agents.list[].skills[]`

---

## Related

- [[../workspaces (Agent Workspaces)/README|Workspaces]] — Agent workspace files
- [[../cron (Scheduled Jobs)/CRON|CRON]] — Scheduled jobs per agent
- [[../hooks (Event Hooks)/HOOKS|HOOKS]] — Event hooks
```

**Step 2: Verify the file looks correct in Obsidian**

Open `plugin-skills (Custom Skills)/PLUGIN-SKILLS.md` and confirm tables render correctly.

**Step 3: Commit**
```bash
git add "00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)/plugin-skills (Custom Skills)/PLUGIN-SKILLS.md"
git commit -m "docs(openclaw): update plugin-skills to 6-agent architecture"
```

---

## Task 2: Create `agents/README.md` — master agent reference

**Files:**
- Create: `agents/README.md`

There is no `agents/` folder yet. This becomes the single authoritative page for all 6 active agents — their IDs, models, skills, channel bindings, workspace paths, and cron jobs.

**Step 1: Create the directory and file**

Content:

```markdown
---
title: "OpenClaw — Agents"
type: reference
status: active
created: 2026-05-11
modified: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO-02, openclaw/agents]
---

# OpenClaw — Agents

6 active agents on omos-core. All run `github-copilot/claude-sonnet-4.6` with Amazon Bedrock Sonnet 4.6 as fallback. Config: `~/.openclaw/openclaw.json` → `agents.list[]`.

---

## Active Agents

| ID | Name | Node | Role | Discord | Telegram |
|---|---|---|---|---|---|
| `main` | Legacy | omos-core | Orchestrator — Zeus's primary interface | `#legacy-chat` `1478604243509448705` | Topic `292` |
| `zenith` | Zenith | omos-core | HP 25-49 — body, mind, money, spirit, intelligence | `#zenith-chat` `1503272334990118952` | Topic `318` |
| `unity` | Unity | omos-core | LE 50-74 — family, home, land, estate, generations | `#unity-chat` `1503272335921123369` | Topic `319` |
| `eternal` | Eternal | omos-core | GE 75-99 — empire, brand, capital, ventures, wealth | `#eternal-chat` `1503272336978087986` | Topic `320` |
| `sovereignty` | Sovereignty | omos-core | SO 00-24 — infra, security, automation, dev, research | `#sovereignty-chat` `1503272338274390076` | Topic `321` |
| `sentinel` | Sentinel | omos-sentinel | External watchdog + failover | `#sentinel-chat` `1482355265092190319` | DM only |

---

## Agent Detail

### Legacy (`main`)
- **Workspace:** `~/.openclaw/workspace/`
- **Skills:** coding-agent, github, diffs, clawhub, session-logs, gh-issues, 1password, oracle, healthcheck, taskflow, taskflow-inbox-triage, tavily, summarize, openai-whisper, video-frames, nano-pdf, camsnap, obsidian, wiki-maintainer, weather, voice-call, telegram, onemind-maintainer, onemind-inbox, onemind-sync
- **Cron:** `legacy-daily-report` — `0 23 * * *` (11pm ET)
- **Purpose:** Routes all requests. Only agent Zeus talks to directly. Delegates to domain agents via skills.

### Zenith
- **Workspace:** `~/.openclaw/workspace-zenith/`
- **Skills:** agent-zenith + base + spartan, oracle, identity, librarian, pulse
- **Cron:** none (domain cron jobs run via spartan-weekly, oracle-weekly — these are old agent IDs; Zenith handles HP domain)
- **Purpose:** All HP 25-49 work — fitness, habits, mindset, personal finance, intelligence pipeline.

### Unity
- **Workspace:** `~/.openclaw/workspace-unity/`
- **Skills:** agent-unity + base + haven, helix, eden, nomad
- **Cron:** none (haven-weekly handled by Unity's haven skill)
- **Purpose:** All LE 50-74 work — family, relationships, home, sustainability, sovereign living.

### Eternal
- **Workspace:** `~/.openclaw/workspace-eternal/`
- **Skills:** agent-eternal + base + empire, empire-ops, creed, trinity, apex, finance
- **Cron:** none (empire-weekly, apex-daily handled via Eternal's domain skills)
- **Purpose:** All GE 75-99 work — business strategy, brand, capital, R&D.

### Sovereignty
- **Workspace:** `~/.openclaw/workspace-sovereignty/`
- **Skills:** agent-sovereignty + base + grid, guardian, guardian-ops, overwatch, eagleeye, coding-agent, github, diffs
- **Cron:** `grid-daily-report` — `0 22 * * *` (10pm ET)
- **Purpose:** All SO 00-24 work — VPS, security, automation, dev, research.

### Sentinel
- **Workspace:** `~/.openclaw/workspace-sentinel/` (on omos-sentinel)
- **Skills:** onemind-maintainer, taskflow, tavily, summarize, telegram, onemind-inbox, onemind-sync
- **Cron:** `sentinel-daily-report` — `55 22 * * *` (10:55pm ET)
- **Purpose:** External watchdog. Monitors omos-core from outside. Activates when core is unreachable. DMs Zeus via `@omos_sentinel_bot`.

---

## Cron Jobs by Agent

See [[../cron (Scheduled Jobs)/CRON|CRON]] for full schedule. Note: some cron job `agentId` values (spartan, oracle, empire, etc.) reference old domain agent IDs — these run as their own workspace contexts, separate from the 6 super-agents.

---

## Channel Bindings

Bindings live in `openclaw.json` → `bindings[]`. 19 Discord channel bindings exist — 6 for super-agents, 13 for old domain agent IDs (spartan, oracle, empire, etc.) that still have their own Discord channels routed directly.

To view all bindings:
```bash
ssh ubuntu@100.123.126.34 "openclaw agent bindings"
```

---

## Managing Agents

```bash
# List all agents
ssh ubuntu@100.123.126.34 "openclaw agents list"

# Check agent status
ssh ubuntu@100.123.126.34 "openclaw agent status main"

# Send a test message to an agent
ssh ubuntu@100.123.126.34 "openclaw message main 'hello'"

# View agent workspace files
ls ~/.openclaw/workspace-zenith/

# Reload after openclaw.json changes
ssh ubuntu@100.123.126.34 "kill -HUP \$(pgrep -f 'openclaw.*gateway')"
```

---

## Related

- [[../plugin-skills (Custom Skills)/PLUGIN-SKILLS|Plugin Skills]] — Skills loaded by each agent
- [[../workspaces (Agent Workspaces)/README|Workspaces]] — Workspace file structure
- [[../channels/CHANNELS|Channels]] — Full channel/topic map
- [[../cron (Scheduled Jobs)/CRON|Cron]] — Scheduled jobs
```

**Step 2: Commit**
```bash
git add "00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)/agents/README.md"
git commit -m "docs(openclaw): add agents/README.md — 6-agent master reference"
```

---

## Task 3: Create `workspaces/backup-repo.md` — document the workspace backup system

**Files:**
- Create: `workspaces (Agent Workspaces)/backup-repo.md`

This documents the `openclaw-core-workspaces` GitHub repo, `watch.sh`, and the divergence issue we fixed (for future reference).

**Step 1: Create backup-repo.md**

```markdown
---
title: "OpenClaw — Workspace Backup Repo"
type: reference
status: active
created: 2026-05-11
modified: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO-02, openclaw/workspaces]
---

# Workspace Backup Repo

All 15 agent workspaces on omos-core are version-controlled and auto-pushed to GitHub.

**Repo:** `github.com/OneMind-OS/openclaw-core-workspaces`  
**Branch:** `main`  
**Local path:** `/home/ubuntu/.openclaw/` (the repo root IS the `.openclaw` directory)

---

## What's Tracked

The repo tracks all `workspace*/` directories — every agent's 7 identity files:

| File | Purpose |
|------|---------|
| `IDENTITY.md` | Agent identity, name, role, personality |
| `SOUL.md` | Core values, mission, operating principles |
| `TOOLS.md` | Available tools and how to use them |
| `HEARTBEAT.md` | Heartbeat prompt — periodic self-check |
| `AGENTS.md` | Instructions from OpenClaw platform |
| `USER.md` | User context — Zeus's profile and preferences |
| `MEMORY.md` | Persistent agent memory (evolves over time) |

Also tracked: `.gitignore`, `wiki/`, `hooks/`, `plugin-skills/` (select files), `watch.sh`, `backup.sh`, `openclaw.json` (config snapshot), systemd templates, `bootstrap.sh`, `README.md`.

---

## Auto-Sync: watch.sh

A background process watches all `workspace*/` dirs and auto-commits + pushes on any change.

**Script:** `/home/ubuntu/.openclaw/watch.sh`  
**Process:** Runs as `ubuntu` (started at boot, verify with `ps aux | grep watch.sh`)  
**Debounce:** 3 second quiet period before committing  
**Commit format:** `watch: {EVENT} {FILE} — {ISO_TIMESTAMP}`

```bash
# Verify watch.sh is running
ps aux | grep watch.sh | grep -v grep

# View recent pushes
git -C ~/.openclaw log --oneline -10

# Manual push if watch is down
cd ~/.openclaw && git add workspace*/ && git commit -m "manual: workspace sync" && git push
```

---

## Health Check

**Signs watch is working:** `git -C ~/.openclaw log --oneline -3` shows recent `watch:` commits.

**Signs watch is broken:**
- No commits in 24h while agents are active
- `git -C ~/.openclaw status` shows `ahead of origin/main by N commits`
- `git -C ~/.openclaw status` shows diverged

**Fix diverged repo (pull + rebase):**
```bash
cd ~/.openclaw

# Stash any untracked files blocking the rebase
git stash push -u -m "pre-rebase save"

# Pull remote and rebase local commits on top
git pull --rebase origin main

# Pop stash (ignore "already exists" messages — files are on disk)
git stash pop

# Push
git push origin main
```

> Note: On 2026-05-11 the repo diverged (2 local commits, 16 remote from openclaw-config-sync.service). Fixed via stash + rebase. The `openclaw-config-sync.service` was writing to the same repo from a different service — now inactive.

---

## Workspaces on VPS

| Workspace Dir | Agent | Active in openclaw.json |
|---|---|---|
| `workspace/` | Legacy (main) | ✅ |
| `workspace-zenith/` | Zenith | ✅ |
| `workspace-unity/` | Unity | ✅ |
| `workspace-eternal/` | Eternal | ✅ |
| `workspace-sovereignty/` | Sovereignty | ✅ |
| `workspace-eden/` | Eden (old domain agent) | legacy |
| `workspace-empire/` | Empire (old domain agent) | legacy |
| `workspace-grid/` | Grid (old domain agent) | legacy |
| `workspace-guardian/` | Guardian (old domain agent) | legacy |
| `workspace-haven/` | Haven (old domain agent) | legacy |
| `workspace-librarian/` | Librarian (old domain agent) | legacy |
| `workspace-oracle/` | Oracle (old domain agent) | legacy |
| `workspace-overwatch/` | Overwatch (old domain agent) | legacy |
| `workspace-spartan/` | Spartan (old domain agent) | legacy |
| `workspace-trinity/` | Trinity (old domain agent) | legacy |

"Legacy" workspaces remain on disk and in the backup repo — they're not active in `openclaw.json` but workspace files still exist.

---

## Restore from Backup

To restore any agent's workspace files from the GitHub repo:

```bash
# Clone fresh (e.g., after VPS rebuild)
git clone https://github.com/OneMind-OS/openclaw-core-workspaces.git /home/ubuntu/.openclaw

# Or restore a single file
git -C ~/.openclaw checkout HEAD -- workspace-zenith/IDENTITY.md
```

---

## Related

- [[README|Workspaces README]] — Workspace file reference
- [[../agents/README|Agents]] — Agent roster
- [[../install (Install & Deploy)/README|Install]] — Full VPS bootstrap
```

**Step 2: Commit**
```bash
git add "00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)/workspaces (Agent Workspaces)/backup-repo.md"
git commit -m "docs(openclaw): add workspace backup-repo.md — repo health, watch.sh, restore"
```

---

## Task 4: Create `mcp/README.md` — document MCP servers

**Files:**
- Create: `mcp (MCP Servers)/README.md`

No MCP documentation exists. `circle-mcp` is the only active server.

**Step 1: Create mcp/README.md**

```markdown
---
title: "OpenClaw — MCP Servers"
type: reference
status: active
created: 2026-05-11
modified: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO-02, openclaw/mcp]
---

# OpenClaw — MCP Servers

MCP (Model Context Protocol) servers extend agent capabilities with external integrations. Config: `~/.openclaw/openclaw.json` → `mcp.servers{}`.

---

## Active MCP Servers

### circle-mcp

| Field | Value |
|-------|-------|
| **Status** | ✅ Active |
| **Purpose** | Circle community platform — manage spaces, posts, courses, members, events |
| **Config key** | `mcp.servers.circle-mcp` |

Provides 100+ tools for interacting with the Circle community platform. Agents with this MCP loaded can create posts, manage members, run events, and more.

To check Circle MCP config:
```bash
ssh ubuntu@100.123.126.34 "openclaw mcp list"
```

---

## Managing MCP Servers

```bash
# List all MCP servers
ssh ubuntu@100.123.126.34 "openclaw mcp list"

# Add a new MCP server (edit openclaw.json directly)
# Under mcp.servers, add:
# "my-server": { "command": "...", "args": [...], "env": {...} }

# Reload after changes
ssh ubuntu@100.123.126.34 "kill -HUP \$(pgrep -f 'openclaw.*gateway')"
```

Config location: `~/.openclaw/openclaw.json` → `mcp.servers{}`

---

## Adding a New MCP Server

1. SSH into omos-core
2. Edit `~/.openclaw/openclaw.json`
3. Add entry under `mcp.servers`:
```json
"server-name": {
  "command": "node",
  "args": ["/path/to/server.js"],
  "env": {
    "API_KEY": "..."
  }
}
```
4. Reload: `kill -HUP $(pgrep -f 'openclaw.*gateway')`
5. Verify: `openclaw mcp list`
6. Document here

---

## Related

- [[../agents/README|Agents]] — Which agents have MCP access
- [[../plugin-skills (Custom Skills)/PLUGIN-SKILLS|Plugin Skills]] — Skill-based integrations
```

**Step 2: Commit**
```bash
git add "00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)/mcp (MCP Servers)/README.md"
git commit -m "docs(openclaw): add mcp/README.md — circle-mcp and MCP server reference"
```

---

## Task 5: Mirror workspace files from VPS into vault

**Files:**
- Create/Update: `workspaces (Agent Workspaces)/{agent}/` — add all 7 files per active agent

The 5 active super-agent workspaces (Legacy, Zenith, Unity, Eternal, Sovereignty) need their full 7 files mirrored from the VPS into the vault as rebuild reference. Pull from omos-core via SSH.

**Step 1: Pull files for each active super-agent**

For each of: `workspace` (Legacy), `workspace-zenith`, `workspace-unity`, `workspace-eternal`, `workspace-sovereignty` — pull all 7 files (AGENTS.md, HEARTBEAT.md, IDENTITY.md, MEMORY.md, SOUL.md, TOOLS.md, USER.md) and write into the corresponding vault `workspaces/` folder.

Map:
- `workspace/` → `workspaces (Agent Workspaces)/legacy (Legacy Workspace)/`
- `workspace-zenith/` → `workspaces (Agent Workspaces)/zenith (Zenith Workspace)/`
- `workspace-unity/` → `workspaces (Agent Workspaces)/unity (Unity Workspace)/`
- `workspace-eternal/` → `workspaces (Agent Workspaces)/eternal (Eternal Workspace)/`
- `workspace-sovereignty/` → `workspaces (Agent Workspaces)/sovereignty (Sovereignty Workspace)/`

Also pull Sentinel from omos-sentinel: `workspace/` → `workspaces (Agent Workspaces)/sentinel (Sentinel Workspace)/`

**Step 2: Commit all workspace mirrors**
```bash
git add "00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)/workspaces (Agent Workspaces)/"
git commit -m "docs(openclaw): mirror active agent workspace files from omos-core"
```

---

## Task 6: Update `workspaces/README.md` — add backup repo reference and file guide

**Files:**
- Modify: `workspaces (Agent Workspaces)/README.md`

Current README is sparse. Add: what each of the 7 files does, backup repo link, agent list, and how to update files.

**Step 1: Read current README, then append missing sections**

Add below existing content:

```markdown
## The 7 Workspace Files

| File | Purpose | Who writes it |
|------|---------|---------------|
| `IDENTITY.md` | Agent name, role, personality, operating style | Zeus / Claude Code |
| `SOUL.md` | Core values, mission, principles that never change | Zeus |
| `TOOLS.md` | Available tools reference — what the agent can do | Zeus / OpenClaw |
| `HEARTBEAT.md` | Periodic self-check prompt — keeps agent calibrated | Zeus |
| `AGENTS.md` | Instructions from OpenClaw platform (routing, context) | OpenClaw auto-generated |
| `USER.md` | Zeus's profile, preferences, communication style | Zeus |
| `MEMORY.md` | Evolving memory — what the agent has learned about Zeus | Agent (auto-written) |

## Active Agent Workspaces

| Agent | VPS Path | Vault Mirror |
|-------|---------|--------------|
| Legacy | `~/.openclaw/workspace/` | `legacy (Legacy Workspace)/` |
| Zenith | `~/.openclaw/workspace-zenith/` | `zenith (Zenith Workspace)/` |
| Unity | `~/.openclaw/workspace-unity/` | `unity (Unity Workspace)/` |
| Eternal | `~/.openclaw/workspace-eternal/` | `eternal (Eternal Workspace)/` |
| Sovereignty | `~/.openclaw/workspace-sovereignty/` | `sovereignty (Sovereignty Workspace)/` |
| Sentinel | `~/.openclaw/workspace/` (omos-sentinel) | `sentinel (Sentinel Workspace)/` |

## Backup & Version Control

All workspace files are auto-committed and pushed to GitHub by `watch.sh`.  
Repo: `github.com/OneMind-OS/openclaw-core-workspaces`  
See [[backup-repo|Backup Repo]] for health checks and restore procedures.

## Editing Workspace Files

1. Edit the file directly in Obsidian (vault) or via SSH on the VPS
2. If editing in vault: sync will push to VPS via `onemind-watch.service`
3. If editing on VPS: `watch.sh` auto-commits and pushes to GitHub within ~3s
4. Reload agent after significant changes: `openclaw agent reload {id}`
```

**Step 2: Commit**
```bash
git add "00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/openclaw (OpenClaw Platform)/workspaces (Agent Workspaces)/README.md"
git commit -m "docs(openclaw): expand workspaces README — 7-file guide, backup repo, agent map"
```

---

## Task 7: Create Todoist tracking task

**Step 1: Create Todoist task to track backup repo health**

Via Legacy agent or directly in Todoist:
- Project: SO-02 Infrastructure
- Task: "Weekly: verify openclaw-core-workspaces backup repo is pushing"
- Description: Check `git -C ~/.openclaw log --oneline -3` for recent `watch:` commits. If stale, run: `cd ~/.openclaw && git pull --rebase origin main && git push origin main`
- Recurrence: Every Monday

---

## Execution Order

1. Task 1 (PLUGIN-SKILLS.md update) — standalone, no dependencies
2. Task 2 (agents/README.md) — standalone
3. Task 3 (backup-repo.md) — standalone
4. Task 4 (mcp/README.md) — standalone
5. Task 5 (workspace mirrors) — requires SSH, do last since most files
6. Task 6 (workspaces/README.md update) — after Task 5
7. Task 7 (Todoist task) — anytime

Tasks 1-4 can run in parallel. Tasks 5-6 are sequential.
