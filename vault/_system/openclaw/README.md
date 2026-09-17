# openclaw/ — OpenClaw Runtime

Legacy's sovereign AI runtime for OneMind OS. The **VPS is the source of truth** — this repo mirrors the native OpenClaw structure and syncs bidirectionally via GitHub. **9 agents, 30 skills, 9 lobster routines.**

---

## Architecture

```
┌──────────────────────────────────────────────────────┐
│  MAC (Obsidian + VS Code)                            │
│  Edit workspaces, skills, MDs locally                │
└────────────────┬─────────────────────────────────────┘
                 │ git push (Obsidian-git, 1-5 min)
          ┌──────▼──────┐
          │   GitHub    │
          └──────┬──────┘
                 │ git pull (codex-pull.timer, every 1 min)
┌────────────────▼─────────────────────────────────────┐
│  VPS (104.248.4.212)                                 │
│                                                       │
│  ~/.openclaw/workspace → repo/openclaw/workspaces/legacy  │
│  openclaw.service → gateway process                  │
│  codex-watch.service → auto-push VPS changes (3s)    │
│  openclaw-reload.service → restart on skill changes  │
│                                                       │
│  Docker: NATS JetStream | PostgreSQL+pgvector | Bridge│
└──────────────────────────────────────────────────────┘
```

---

## Directory Structure

```
openclaw/
├── openclaw.json.example         ← Config template (copy to ~/.openclaw/openclaw.json)
├── README.md                     ← This file
├── workspaces/                   ← Live agent workspaces (symlinked on VPS)
│   ├── legacy/                   ← Default agent — has ALL skills + lobster
│   │   ├── skills/               ← 30 OpenClaw skills (JS, CommonJS)
│   │   ├── lobster/              ← 9 scheduled workflow routines
│   │   ├── canvas/               ← HTML dashboards served by gateway
│   │   ├── scripts/              ← Shell utilities
│   │   ├── bridge.py             ← NATS ↔ OpenClaw relay
│   │   ├── SOUL.md               ← Personality
│   │   ├── AGENTS.md             ← Operating instructions
│   │   ├── USER.md               ← Zeus's profile
│   │   ├── HEARTBEAT.md          ← Proactive check checklist
│   │   └── ... (BOOT, TOOLS, IDENTITY, MEMORY, BOOTSTRAP)
│   ├── spartan/                  ← ⚔️  Body & Physical Mastery
│   ├── oracle/                   ← 🔮 Mind, Money & Meaning
│   ├── eden/                     ← 🌳 Food Growing & Agriculture
│   ├── empire/                   ← 👑 Business & Wealth
│   ├── haven/                    ← 🌿 Home, Family & Estate
│   ├── guardian/                 ← 🛡️  Security & Legal
│   ├── grid/                     ← ⚙️  Infrastructure Ops
│   └── trinity/                  ← 🔺 Community & Public
└── deploy/                       ← Deployment configs (systemd, scripts, env)
    ├── systemd/
    │   ├── openclaw.service
    │   ├── codex-pull.timer
    │   ├── codex-pull.service
    │   ├── codex-watch.service
    │   └── openclaw-reload.service
    ├── scripts/
    │   ├── setup.sh              ← Full VPS deploy from scratch
    │   ├── setup-symlinks.sh     ← Create workspace symlinks
    │   ├── codex-watch.sh        ← Real-time git sync script
    │   ├── openclaw-reload.sh    ← Skill change watcher
    │   └── workspace-backup.sh   ← Safety net cron (every 15 min)
    └── env.example               ← Environment variable template
```

> Deep-dive docs: `docs/openclaw/VPS-GATEWAY.md` and `docs/openclaw/VPS-SERVICES.md` at repo root.

---

## Agents

Configured in `openclaw.json` under `agents.defaults` (shared config) + `agents.list[]` (per-agent overrides).

| Agent | Model | Domain | Tool Profile | Workspace |
|-------|-------|--------|-------------|---------|
| `legacy` (default) | azure-openai/gpt-5.2 | UI 00-24 | `full` | `workspaces/legacy/` |
| `hp` | azure-openai/gpt-4o | HP 25-49 | `messaging` | `workspaces/hp/` |
| `le` | azure-openai/gpt-4o | LE 50-74 | `messaging` | `workspaces/le/` |
| `ge` | azure-openai/gpt-5.2 | GE 75-99 | `messaging` | `workspaces/ge/` |
| `community` | azure-openai/gpt-4o | Community | `messaging` | `workspaces/community/` |

### Tools vs Skills (OpenClaw native concepts)

OpenClaw separates these two concepts:

- **Tools** — Built-in functions the agent can call (`exec`, `browser`, `message`, `web_search`, `memory_search`, etc.). Controlled by `tools.profile` + `tools.allow`/`tools.deny` in `openclaw.json`.
- **Skills** — `SKILL.md` instruction bundles that teach the agent HOW to use external APIs. Loaded from workspace `skills/` directories. These inject guidance into the system prompt.

Tool profiles control what built-in capabilities each agent has:

| Profile | Includes |
|---------|----------|
| `full` | Everything — exec, browser, fs, memory, messaging, etc. |
| `messaging` | message + sessions tools only |
| `coding` | fs, runtime, sessions, memory, image |
| `minimal` | session_status only |

Per-agent tool access (from `openclaw.json`):

| Agent | Profile | Additional allowed tools |
|-------|---------|------------------------|
| `legacy` | `full` | `lobster` |
| `hp` | `messaging` | `web_fetch`, `web_search`, `group:memory`, `lobster` |
| `le` | `messaging` | `web_fetch`, `web_search`, `group:memory` |
| `ge` | `messaging` | `web_fetch`, `web_search`, `group:memory`, `lobster` |
| `community` | `messaging` | `web_fetch`, `web_search`, `group:memory` |

### Where do skills live?

OpenClaw loads skills from three locations (highest precedence first):

1. **`<workspace>/skills/`** — Per-agent only (highest priority)
2. **`~/.openclaw/skills/`** — Shared across all agents on the gateway
3. **Bundled skills** — Shipped with OpenClaw install (lowest priority)

Additional shared paths can be added via `skills.load.extraDirs` in `openclaw.json`.

**Current VPS state**: All 30 custom skills live in `workspaces/legacy/skills/` (per-agent). No `~/.openclaw/skills/` directory exists yet. The 8 domain agents operate with their `messaging` tool profile + allowed built-in tools.

To share skills across agents in the future, the OpenClaw-native approach is:
- Move shared skills to `~/.openclaw/skills/` (symlinked from repo)
- Or add `skills.load.extraDirs` pointing to a shared directory
- Keep agent-specific skills in `<workspace>/skills/`

### Subagent Architecture

```
GATEWAY (openclaw.json)
│
├── agents.defaults        ← shared config inherited by ALL agents
│   └── subagents:
│       maxConcurrent: 8        ← total subagents running at once, gateway-wide
│       maxSpawnDepth: 3        ← how many levels deep spawning can go
│       maxChildrenPerAgent: 5  ← how many children one parent can spawn
│       model: gpt-4o-mini
│
└── agents.list[]
    │
    ├── legacy  (explicit subagent override — same values, intentional)
    │   │
    │   ├── [subagent A]  ← spawned for subtask 1
    │   │   └── [subagent A1]  ← depth 2
    │   │       └── [subagent A1a]  ← depth 3 (maximum)
    │   ├── [subagent B]  ← subtask 2
    │   ├── [subagent C]  ← subtask 3
    │   ├── [subagent D]  ← subtask 4
    │   └── [subagent E]  ← subtask 5 (maxChildrenPerAgent hit)
    │
    ├── hp / le / ge / community  (inherit defaults)
    │   └── up to 5 children, 3 levels deep — but eat into the
    │       8 gateway-wide concurrent slots shared with legacy
```

**Key rules:**
- `maxChildrenPerAgent` — per agent. Legacy gets 5, HP gets 5, etc.
- `maxConcurrent` — gateway-wide total. All agents share the pool.
- `maxSpawnDepth` — depth 4+ produces cascading hallucinations and runaway token spend. 3 is the practical ceiling.

**Legacy theoretical ceiling by dimension** (for reference when upgrading compute):

| Dimension | Current | Useful ceiling | Beyond that |
|-----------|---------|---------------|-------------|
| `maxConcurrent` | 8 | 12–16 | Subagents fight for rate limits, no gain |
| `maxSpawnDepth` | 3 | 3 | Depth 4+ = cascading failures + cost explosion |
| `maxChildrenPerAgent` | 5 | 8–10 | Parent can't coherently synthesize >10 results |
| `model` | gpt-4o-mini | gpt-4o | Upgrade when subagents do complex reasoning, not just fetching |

### Per-agent config highlights

Each agent also gets:
- **`memorySearch.extraPaths`** — Scoped Codex access (HP sees HP domain, GE sees GE domain, Legacy sees all)
- **`heartbeat`** — Per-agent schedule, Discord channel, and prompt
- **`subagents`** — Legacy: explicit 8/3/5 override. HP/LE/GE/Community: inherit defaults

### Workspace contents by agent

| | legacy | hp / le / ge |
|---|---|---|
| Personality MDs | SOUL, AGENTS, USER, IDENTITY, HEARTBEAT, TOOLS, BOOT, BOOTSTRAP, MEMORY | SOUL, AGENTS, USER, IDENTITY, HEARTBEAT, TOOLS, BOOTSTRAP |
| `skills/` | 26 skills | — (none currently) |
| `lobster/` | 8 routines | — |
| `canvas/` | dashboard, holomat | — |
| `bridge.py` | NATS relay | — |
| `scripts/` | Shell utilities | — |
| `.gitignore` | state, memory, .openclaw, .env, secrets, node_modules | Same |

---

## Skills (26 total)

All skills are CommonJS (`index.js` + `package.json` + `SKILL.md`), pure HTTPS, zero npm dependencies unless noted.

### Core Platform
| Skill | Description |
|-------|-------------|
| `api-gateway-1.0.48` | Core API gateway for the OpenClaw platform |
| `proactive-agent-3.1.0` | Proactive behaviors (heartbeat, check-ins) |
| `nexus-bridge` | NATS JetStream dispatch to Agno agents |
| `find-skills-0.1.0` | Skill discovery and search |
| `skill-creator-0.1.0` | Create new skills from templates |
| `hitl-escalate-1.0.0` | Human-in-the-loop escalation workflows |
| `onepassword-1.0.0` | 1Password vault — read/write items, save config backups as secure notes |

### Productivity
| Skill | Description |
|-------|-------------|
| `todoist` | Todoist task management |
| `todoist-api` | Todoist REST API direct |
| `clickup-1.0.0` | ClickUp tasks, docs, search (21 actions) |
| `obsidian-1.0.0` | Obsidian vault read/write |
| `github-api-1.0.3` | GitHub repos, issues, PRs |

### Communication
| Skill | Description |
|-------|-------------|
| `twilio-1.0.0` | SMS/voice via Twilio |
| `resend-1.0.0` | Transactional email via Resend |

### Browser & Desktop
| Skill | Description |
|-------|-------------|
| `browser-1.0.0` | Web browsing automation |
| `agent-browser-0.2.0` | Agent-driven browser control |
| `desktop-control-1.0.0` | Desktop automation |

### Smart Home
| Skill | Description |
|-------|-------------|
| `homestead-control` | Home Assistant device control |
| `vision-1.0.0` | HA LLM Vision + Frigate NVR cameras |

### Personal
| Skill | Description |
|-------|-------------|
| `recipe-scraper-1.0.0` | Recipe extraction from URLs |
| `spotify-1.0.0` | Spotify playback + library |
| `weather-1.0.0` | Weather forecasts |

### Transportation
| Skill | Description |
|-------|-------------|
| `transit-1.0.0` | Public transit routes/schedules |
| `uber-1.0.0` | Ride estimates and booking |

### Business
| Skill | Description |
|-------|-------------|
| `zoho-cli` | Zoho CRM/Desk/Calendar |
| `lob-1.0.0` | Physical mail via Lob API |
| `credit-repair-1.0.0` | FCRA dispute letter generation |

---

## Lobster Routines (8)

Scheduled workflows that Legacy executes on a cron-like basis.

| Routine | Purpose |
|---------|---------|
| `morning-brief.lob` | Daily morning briefing |
| `health-check.lob` | System health monitoring |
| `git-sync.lob` | Repository sync verification |
| `habit-checker.lob` | HP habit tracking |
| `project-state.lob` | ClickUp/Todoist project status |
| `earnings-tracker.lob` | GE earnings monitoring |
| `second-brain.lob` | Codex knowledge sync |
| `tech-digest.lob` | Tech news digest |

---

## Hot vs Cold Reload

| What changed | Reload? | How |
|---|---|---|
| SOUL.md, AGENTS.md, USER.md, etc. | **Hot** — automatic | Read every turn |
| `skills/*` | **Cold** — restart needed | `openclaw-reload.service` auto-restarts |
| `lobster/*` | **Cold** — restart needed | `openclaw-reload.service` auto-restarts |
| `openclaw.json` | **Cold** — restart needed | `openclaw-reload.service` auto-restarts |
| `canvas/*` | **Live** — served directly | Static file server |

> You can edit SOUL.md, AGENTS.md, USER.md from Mac (Obsidian) and Legacy picks up changes on the next message — no restart needed.

---

## Systemd Services

| Service | Purpose | Type |
|---------|---------|------|
| `openclaw.service` | Gateway process (4GB mem limit) | always-on |
| `openclaw-reload.service` | Restart gateway on skill/config changes | always-on |
| `codex-pull.timer` | Git pull from GitHub every 1 min | timer |
| `codex-watch.service` | Push VPS changes to GitHub (3s debounce) | always-on |

Cron: `workspace-backup.sh` every 15 min (safety net for missed changes).

Unit files and scripts are in `deploy/systemd/` and `deploy/scripts/`.

---

## VPS Setup

| Property | Value |
|----------|-------|
| Provider | DigitalOcean |
| Region | NYC3 |
| Public IP | 134.122.8.135 |
| Tailscale IP | 100.104.61.29 |
| OS | Ubuntu 24.04 LTS |
| Specs | 8GB RAM, 4 vCPU, 240GB NVMe |
| SSH User | root |
| OpenClaw User | onemind |
| OpenClaw Version | 2026.3.2 |
| Gateway Port | 18789 (Tailscale only) |
| Channels | Discord only |

### Symlink structure on VPS

```
~/.openclaw/workspace         → /home/onemind/onemind-os/openclaw/workspaces/legacy
~/.openclaw/workspace-legacy  → /home/onemind/onemind-os/openclaw/workspaces/legacy
~/.openclaw/workspace-hp      → /home/onemind/onemind-os/openclaw/workspaces/hp
~/.openclaw/workspace-le      → /home/onemind/onemind-os/openclaw/workspaces/le
~/.openclaw/workspace-ge      → /home/onemind/onemind-os/openclaw/workspaces/ge
```

This means the gateway reads directly from the git repo. Git pull → instant update.

### What's on VPS but NOT in repo (correctly excluded)

| Path | Purpose |
|------|---------|
| `~/.openclaw/openclaw.json` | Gateway config (has API key refs) |
| `/opt/onemind/.env` | All secrets |
| `~/.openclaw/memory/legacy.sqlite` | Agent memory DB |
| `~/.openclaw/agents/*/` | Agent session state |
| `~/.openclaw/devices/paired.json` | Paired devices |
| `~/.openclaw/identity/device.json` | Ed25519 keypair |
| `~/.openclaw/exec-approvals.json` | Execution approvals |
| `~/.openclaw/logs/`, `cron/`, `media/` | Runtime data |

---

## Deploy From Scratch

If the VPS disappears, here's the recovery plan:

### Quick version

```bash
# 1. Spin up Ubuntu 24.04 droplet
# 2. Copy deploy scripts
scp -r openclaw/deploy/ root@NEW_IP:/tmp/deploy/

# 3. Run setup
ssh root@NEW_IP 'bash /tmp/deploy/scripts/setup.sh'

# 4. Configure secrets
ssh root@NEW_IP
vim /opt/onemind/.env                        # Fill in API keys
cp /tmp/deploy/../openclaw.json.example ~/.openclaw/openclaw.json
vim ~/.openclaw/openclaw.json                # Fill in config

# 5. Auth + tailscale
sudo -u onemind gh auth login
curl -fsSL https://tailscale.com/install.sh | sh && tailscale up

# 6. Start everything
cd /opt/onemind && docker compose --profile data up -d
systemctl start openclaw codex-watch openclaw-reload codex-pull.timer
```

### What you lose (and how to recover)

| Lost | Impact | Recovery |
|------|--------|----------|
| Agent memory (`legacy.sqlite`) | Conversation history | Starts fresh — acceptable |
| Paired devices | Need to re-pair Mac/phone | `openclaw pair` from each device |
| Identity keypair | New device identity | Auto-generated on first run |
| Session state | Active conversations reset | Resumes on next message |
| `.env` / `openclaw.json` | Must reconfigure | Templates in `deploy/env.example` + `openclaw.json.example` |

**What you DON'T lose**: All skills, lobster routines, personality MDs, canvas, bridge, scripts — everything in `workspaces/` is in the repo.

### Full setup script

`deploy/scripts/setup.sh` automates: user creation, package install, Node.js + OpenClaw install, repo clone, symlinks, systemd units, firewall. See the script for details.

---

## DigitalOcean Droplet Sizing

Current VPS (8GB/4vCPU) handles the gateway but will hit memory pressure when Legacy runs 8 concurrent subagents alongside the other agents + Docker services.

| Scenario | Droplet | RAM | vCPU | Cost/mo |
|----------|---------|-----|------|---------|
| Current | Basic Regular | 8GB | 4 | ~$48 |
| **Recommended next** | **General Purpose** | **16GB** | **4** | **~$126** |
| Full scale (multi-agent parallel) | CPU-Optimized | 32GB | 8 | ~$336 |
| Future (OM1 + full stack) | CPU-Optimized | 64GB | 16 | ~$672 |

**Why General Purpose 16GB/4vCPU:**
- 8 concurrent gpt-4o-mini subagents = ~1–2GB RAM each under load
- 5 agents + gateway process + Docker (NATS + PostgreSQL + bridge) = ~4GB baseline
- 16GB gives ~10GB headroom for subagent bursts
- Same 4 vCPU is fine — workload is I/O-bound (LLM API calls), not CPU-bound

Upgrade path: DigitalOcean → Resize → **General Purpose, 16GB, NYC3** (same region, zero data migration).

---

## Future: Local Primary + Cloud Fallback

Planned architecture for resilience:

```
┌──────────────────┐         ┌──────────────────┐
│  LOCAL GATEWAY   │ ──────► │  CLOUD GATEWAY   │
│  (Mac / NUC)     │  sync   │  (VPS fallback)  │
│  Primary runtime │ ◄────── │  Hot standby     │
└──────────────────┘         └──────────────────┘
```

- Local runs as primary — fastest response, no latency
- Cloud takes over if local goes offline
- Memory + state sync when back online
- Same repo, same workspaces, same skills

---

## Monitoring

```bash
# Service health
systemctl status openclaw openclaw-reload codex-watch codex-pull.timer

# Live gateway logs
journalctl -u openclaw -f

# Sync activity
journalctl -u codex-watch --since "1 hour ago" --no-pager

# Docker services
ssh root@134.122.8.135 'cd /opt/onemind && docker ps'

# Check symlinks
ssh root@134.122.8.135 'ls -la /home/onemind/.openclaw/workspace*'
```

---

## Related Documentation

| Doc | Location |
|-----|----------|
| VPS gateway details | `docs/openclaw/VPS-GATEWAY.md` |
| Systemd + sync deep dive | `docs/openclaw/VPS-SERVICES.md` |
| Skill format reference | Each skill has a `SKILL.md` |
| Main project README | `README.md` (repo root) |
