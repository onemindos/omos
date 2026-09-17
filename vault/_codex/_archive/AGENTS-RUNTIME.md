# AGENTS.md

Universal instructions for AI coding agents (GitHub Copilot, Cursor, Gemini CLI, Windsurf, Zed, etc.) working in this repository. Read this before touching any code.

---

## What This Repo Is

**OneMind OS** — a personal AI command center built as a monorepo. Not a single app. A network of 5 distinct runtimes bound together by NATS JetStream, with **Legacy** (an OpenClaw agent) as the central intelligence.

```
OpenClaw (Legacy's brain)
    ↕ webhooks/ws
NATS JetStream ── Agno AgentOS (worker agents)
                ── OpenMind OM1 (robots/sensors)
                └─ HuggingFace LeRobot (physical skills inside OM1)
```

**Full stack**: OpenClaw (Node.js) · FastAPI + Agno (Python) · PostgreSQL/pgvector · NATS · LiveKit

---

## Repo Structure

```
onemindos-v0.5/
├── openclaw/                   ← LEGACY'S RUNTIME (primary active code)
│   ├── workspaces/
│   │   ├── legacy/             ← Legacy's OpenClaw workspace (deployed to VPS)
│   │   │   ├── SOUL.md         ← Personality, tone, values
│   │   │   ├── AGENTS.md       ← Legacy's operating instructions
│   │   │   ├── USER.md         ← Zeus's profile and preferences
│   │   │   ├── HEARTBEAT.md    ← Proactive check list (runs every 30m)
│   │   │   ├── TOOLS.md        ← Tool usage guidance
│   │   │   ├── BOOT.md         ← Startup ritual
│   │   │   ├── IDENTITY.md     ← Name, vibe, emoji
│   │   │   ├── bridge.py       ← NATS → OpenClaw webhook relay
│   │   │   ├── memory/         ← Daily logs (YYYY-MM-DD.md)
│   │   │   ├── canvas/         ← HTML surfaces for device displays
│   │   │   ├── lobster/        ← Named routines (lobster run <name>)
│   │   │   └── skills/         ← 28 OpenClaw skills (Node.js packages)
│   │   │       ├── nexus-bridge/       ← NATS dispatch (nexus_dispatch, nexus_query)
│   │   │       ├── home-assistant-1.0.0/ ← Smart home control (6 tools)
│   │   │       ├── github-api-1.0.3/   ← GitHub API access
│   │   │       ├── todoist/            ← Todoist task management
│   │   │       └── */                  ← 24 more skills (see skills/ dir)
│   │   ├── spartan/            ← ⚔️  Body & Physical Mastery (HP 27-28)
│   │   ├── oracle/             ← 🔮 Mind, Money & Meaning (HP 25-26, 29-34)
│   │   ├── eden/               ← 🌳 Food Growing & Agriculture (LE farming)
│   │   ├── empire/             ← 👑 Business & Wealth (GE 75-84)
│   │   ├── haven/              ← 🌿 Home, Family & Estate (LE 50-58)
│   │   ├── guardian/           ← 🛡️  Security & Legal (UI 11 + HP 31 + LE 57)
│   │   ├── grid/               ← ⚙️  Infrastructure Ops (UI 00-12)
│   │   └── trinity/            ← 🔺 Community & Public (GE 75 + Community)
│   ├── deploy/                 ← VPS deployment files
│   │   ├── env.tpl             ← op:// secret references (NOT real values)
│   │   ├── env.example         ← Human-readable env var list
│   │   ├── env.sa.example      ← Static vars template
│   │   ├── systemd/            ← All systemd unit files
│   │   │   ├── openclaw.service        ← Legacy runtime (op run wrapper)
│   │   │   ├── codex-watch.service     ← Auto-commit codex/ + workspaces/ → GitHub
│   │   │   ├── codex-pull.service      ← Pull GitHub → VPS monorepo
│   │   │   ├── codex-pull.timer        ← Triggers codex-pull every 1 min
│   │   │   └── openclaw-reload.service ← Hot-reload on skill file changes
│   │   └── scripts/
│   │       ├── setup.sh              ← Full VPS setup from scratch
│   │       ├── create-1p-items.sh    ← Migrate secrets → 1Password
│   │       ├── openclaw-reload.sh    ← Hot reload workspace on VPS
│   │       ├── workspace-backup.sh   ← Backup VPS workspace → repo
│   │       ├── setup-symlinks.sh     ← Sync repo workspace → VPS path
│   │       └── codex-watch.sh        ← Auto-commit codex/ + workspaces/ changes to GitHub
│   └── openclaw.json.example   ← OpenClaw config template
├── command-center/             ← All UI/dashboard interfaces
│   ├── swarmclaw/              ← Agent orchestration dashboard (Next.js, v0.7.7)
│   ├── lobsterboard/           ← Monitoring dashboard (Node.js, 50 widgets, v0.3.2)
│   ├── openclaw-studio/        ← Next.js admin/dev studio (primary web UI)
│   ├── ios-app/                ← iOS/macOS native app (Swift/Xcode, project: OpenClaw)
│   │   ├── Sources/            ← 21+ modules (Chat, Gateway, Voice, Camera, etc.)
│   │   ├── OpenClaw.xcodeproj  ← Xcode project
│   │   ├── project.yml         ← XcodeGen spec
│   │   ├── WatchApp/           ← watchOS companion app
│   │   ├── ActivityWidget/     ← iOS Live Activity widget
│   │   ├── ShareExtension/     ← iOS Share Extension
│   │   └── fastlane/           ← App Store deployment
│   ├── OpenClawKit/            ← Swift Package for shared iOS/macOS logic
│   ├── flutter-archive/        ← Archived Flutter code (design reference only)
│   │   ├── SCREEN-INVENTORY.md ← 52-screen inventory for TypeScript rebuild reference
│   │   ├── flutter-lib/        ← Original Flutter web app (110 .dart files)
│   │   └── legacy/             ← Older Flutter code (53 .dart files)
│   ├── theme/                  ← Design tokens / shared theme
│   └── frontend design/        ← UI mockups and design specs
├── infra/                      ← Shared infrastructure (all runtimes)
│   ├── docker-compose.yml      ← NATS, OmniDB, n8n, LiveKit, legacy-bridge
│   ├── nats/                   ← NATS hub/leaf config, protocols, bridge scripts
│   ├── README.md               ← VPS reference (ports, services, paths)
│   ├── DO-DEPLOYMENT.md        ← DigitalOcean setup guide
│   ├── DROPLET-RESIZE.md       ← Droplet sizing analysis
│   └── scripts/                ← Droplet setup, firewall, health checks
├── codex/                      ← Obsidian vault — SOURCE OF TRUTH for all documentation
│   ├── .obsidian/              ← Full Obsidian config (plugins, themes, snippets — tracked for portability)
│   ├── 00-24 SO/               ← The Brain: AI, tech, infrastructure, architecture
│   ├── 25-49 HP/               ← Self Path: health, skills, personal finance
│   ├── 50-74 LE/               ← Legacy Path: family, home, estate
│   ├── 75-99 GE/               ← Wealth Path: business, ventures, entrepreneurship
│   ├── 99-organize/            ← Staging area (legacy docs, imports, unsorted)
│   ├── assets/                 ← Images and media
│   ├── ONEMIND-CODEX.md        ← Master document defining the vault
│   ├── CODEX-FRAMEWORK.md      ← Structural rules and conventions
│   ├── AGENTS.md               ← Vault-specific AI agent instructions
│   └── CLAUDE.md               ← Vault-specific Claude instructions
├── scripts/                    ← Top-level utility scripts
├── roadmap/                    ← Feature planning (consciousness, hardware, oneintel, openmind, voice)
├── *.af                        ← Letta agent export snapshots (ezra.af, legacy-ai-main.af)
└── TYPESCRIPT-CROSS-PLATFORM.md ← Cross-platform strategy doc
```

---

## VPS — Where Legacy Actually Runs

| Item | Value |
|------|-------|
| Host | `root@104.248.4.212` |
| OS | Ubuntu 24.04 |
| OpenClaw version | 2026.3.8 (build 3caab92) |
| Gateway config | `/home/onemind/.openclaw/openclaw.json` |
| Service | `openclaw.service` (systemd, runs as `onemind`) |
| op CLI | `/usr/bin/op` v2.32.1 |
| Tailscale HTTPS | `https://onemind-main.tail717ba0.ts.net` |
| Agent architecture | **9 agents** (Legacy + 8 domain agents) |

**Agent Roster (9 agents):**

| Agent | ID | Model | Heartbeat | Discord Channel |
|-------|----|-------|-----------|----------------|
| ⚡ Legacy | `legacy` | azure-openai/gpt-5.2 | 1h | `#legacy-chat` |
| ⚔️ Spartan | `spartan` | azure-openai/gpt-4o | 2h | `#spartan-chat` |
| 🔮 Oracle | `oracle` | azure-openai/gpt-4o | 3h | `#oracle-chat` |
| 🌳 Eden | `eden` | azure-openai/gpt-4o | 3h | `#eden-chat` |
| 👑 Empire | `empire` | azure-openai/gpt-5.2 | 2h | `#empire-chat` |
| 🌿 Haven | `haven` | azure-openai/gpt-4o | 3h | `#haven-chat` |
| 🛡️ Guardian | `guardian` | azure-openai/gpt-4o | 2h | `#guardian-chat` |
| ⚙️ Grid | `grid` | azure-openai/gpt-4o | 3h | `#grid-chat` |
| 🔺 Trinity | `trinity` | azure-openai/gpt-4o | 3h | `#trinity-chat` |

All subagent model: `azure-openai/gpt-4o-mini`

**Agent Tool Profiles:**

| Profile | Agents | Tools |
|---------|--------|-------|
| `full` | Legacy | All tools (coding, fs, browser, lobster, etc.) |
| `messaging` + allow-list | All 8 domain agents | `web_fetch`, `web_search`, `memory_search`, `memory_get`, `lobster`, `read`, `write`, `edit` |

All agents can search, read, and write the codex via `memory_search` (semantic/hybrid) and `read`/`write`/`edit` (direct file access).

**Memory Search Config:**
- Provider: `openai` (text-embedding-3-small)
- Search: Hybrid (BM25 + vector, 0.7/0.3 weight), MMR enabled
- All agents index: `/home/onemind/onemind-os/codex` via `agents.defaults.memorySearch.extraPaths`
- Sync: `watch: true` (auto-reindex on file changes)

**Voice Config:**
- TTS provider: OpenAI (ElevenLabs removed)
- Voice `autoJoin`: disabled (empty array — voice works on-demand only)
- Dave encryption: enabled

**Tailscale Access:**
- HTTPS: `https://onemind-main.tail717ba0.ts.net` → proxy to `http://127.0.0.1:18789`
- `trustedProxies`: `["100.64.0.0/10", "127.0.0.1"]`

**Key VPS directories:**

| Path | Purpose |
|------|---------|
| `/home/onemind/.openclaw/openclaw.json` | Gateway runtime config (all 9 agents) |
| `/home/onemind/.openclaw/workspace` | Symlink → `onemind-os/openclaw/workspaces/legacy` |
| `/home/onemind/.openclaw/agents/<id>/` | Per-agent state, sessions, auth |
| `/home/onemind/onemind-os/openclaw/workspaces/<id>/` | Agent workspaces (synced from GitHub) |
| `/opt/onemind/.env.tpl` | `op://` secret references, injected by `op run` |
| `/opt/onemind/.env.sa` | Static vars + `OP_SERVICE_ACCOUNT_TOKEN` |
| `/home/onemind/onemind-os/` | Monorepo clone (remote: `One-Mind-OS/onemindos-v0.5`) |
| `/home/onemind/onemind-os/codex/` | Codex vault — all agents read this path |

**Systemd services:**

| Service | Purpose |
|---------|---------|
| `openclaw.service` | Legacy's runtime, runs via `op run` |
| `codex-watch.service` | Watches monorepo paths (`onemind-os/codex/` + `workspaces/`) → auto-commits + pushes to GitHub |
| `openclaw-reload.service` | Hot-reloads Legacy when skill files change |
| `codex-pull.timer` | Pulls GitHub → `/home/onemind/onemind-os/` every 1 minute |

**Service lifecycle:**
```bash
ssh root@104.248.4.212 "systemctl status openclaw"
ssh root@104.248.4.212 "systemctl restart openclaw"
ssh root@104.248.4.212 "journalctl -u openclaw -f"
```

---

## Secrets — Critical Rules

All secrets are stored in **1Password vault `00-24 SO (Sovereign Operations)`** (vault ID: `aaeyris5z44pdqspk277u2oppi`).

The systemd service runs:
```
op run --env-file=/opt/onemind/.env.tpl -- /usr/bin/openclaw gateway
```

**Three files on VPS:**

| File | Purpose | Contains |
|------|---------|---------|
| `/opt/onemind/.env.tpl` | Secret template | `KEY=op://vault-id/Item/field` references |
| `/opt/onemind/.env.sa` | Static + SA token | `OP_SERVICE_ACCOUNT_TOKEN`, non-secret vars |
| ~~`/opt/onemind/.env`~~ | OLD — do not use | Plaintext secrets (deprecated) |

**Rules:**
- Never commit real secret values anywhere
- Never edit `.env.tpl` values directly — change the 1Password item instead
- `op://` format: `op://aaeyris5z44pdqspk277u2oppi/Item Name/field-name` (no `{{ }}`, use vault ID not vault name — parentheses in vault name break resolution)
- To add a new secret: create item in 1Password → add line to `env.tpl` → add to `create-1p-items.sh` for future reference → `systemctl restart openclaw`

**Adding a new secret:**
```bash
# 1. Create item in 1Password (in vault "00-24 SO (Sovereign Operations)")
op item create --category=login --vault="00-24 SO (Sovereign Operations)" \
  --title="My Service" "api-key[password]=actual-value"

# 2. Add to env.tpl
echo 'MY_SERVICE_API_KEY=op://aaeyris5z44pdqspk277u2oppi/My Service/api-key' \
  >> /opt/onemind/.env.tpl

# 3. Restart
systemctl restart openclaw
```

---

## Common Commands

### Deploy / Sync

```bash
# Full VPS setup from scratch
bash openclaw/deploy/scripts/setup.sh

# Push workspace changes to VPS and hot-reload
bash openclaw/deploy/scripts/openclaw-reload.sh

# Backup VPS workspace state back to repo
bash openclaw/deploy/scripts/workspace-backup.sh

# Link repo workspace dir to VPS path
bash openclaw/deploy/scripts/setup-symlinks.sh

# Verify deployed workspace codex path is correct
ssh root@104.248.4.212 "grep -n 'home/onemind/onemind-os/codex\|home/openclaw' /opt/onemind/openclaw/workspace/TOOLS.md"

# Check VPS service status
ssh root@104.248.4.212 "systemctl status openclaw codex-watch openclaw-reload codex-pull"
```

### Local Development

```bash
# Copy Legacy workspace locally and run
cp -r openclaw/workspaces/legacy ~/.openclaw/workspace
# Edit ~/.openclaw/openclaw.json with your credentials (see env.example)
openclaw gateway

# Install OpenClaw CLI
npm install -g openclaw

# Install a skill's dependencies
cd openclaw/workspaces/legacy/skills/nexus-bridge && npm install
```

### Skill Development

```bash
# Skills are Node.js packages in workspaces/legacy/skills/<name>/
# Each has: index.js, package.json, SKILL.md
# Functions exported from index.js become OpenClaw tools
npm install    # in skill dir
# Then restart gateway to pick up changes
```

### NATS Subject Conventions

```
digital.agent.{id}.task            → dispatch task to Agno agent
digital.agent.{id}.result.legacy   → agent result back to Legacy
physical.robot.{id}.command        → robot command (OM1)
comms.text.human.zeus              → Zeus's inbound messages
comms.text.agent.legacy            → Legacy's outbound messages
nexus.heartbeat.>                  → all runtime heartbeats
```

---

## OpenClaw Workspace Concepts

| Concept | What it is |
|---------|-----------|
| **Workspace** | Directory of markdown files + skills that define an agent |
| **Skill** | Node.js package exporting tool functions (like MCP servers but native) |
| **Lobster routine** | Named workflow in `lobster/` dir, run via `lobster run <name>` |
| **Canvas** | HTML/CSS/JS file served to connected display devices |
| **Node** | External device (Mac, iPhone, Android) connected to the gateway |
| **Channel** | Communication channel (WhatsApp, iMessage, Telegram, Discord, voice) |
| **Hook** | Triggered action on OpenClaw events |

**Key rule**: Legacy is a pure OpenClaw agent. Do NOT add Legacy back to Agno/Python presets. The `nexus-bridge` skill is the bridge between OpenClaw and NATS.

---

## Development Rules

1. **Legacy stays in OpenClaw** — do not add Legacy to any Python/Agno configs
2. **Secrets via 1Password only** — no `.env` files with real values in the repo
3. **Skill changes require gateway restart** — `openclaw reload` or `systemctl restart openclaw`
4. **SOUL.md / AGENTS.md / USER.md are Legacy's identity** — edit carefully, test immediately
5. **Workspace docs are live config** — changes to markdown files in `workspaces/legacy/` affect Legacy's behavior at runtime
6. **Read-only DB tools** — database query tools must not write/delete data
7. **NATS for cross-runtime** — don't make direct HTTP calls between runtimes; use NATS subjects
8. **codex/ lives at `onemindos-v0.5/codex/`** — this is the canonical local path for the Obsidian vault. When saving any file to the codex, always use `/Users/zeuslegacy/Documents/onemindos-v0.5/codex/`. NEVER save codex files to the separate `/Users/zeuslegacy/Documents/onemindcodex-v0.4/` directory.
9. **Legacy reads codex at `/home/onemind/onemind-os/codex/`** — this is the monorepo clone on VPS, kept fresh by `codex-pull.timer` + `codex-watch.service`. TOOLS.md must reference this exact path.
10. **Codex write guardrails** — agents must NEVER create new top-level folders in codex. New notes go to `06 Inbox (Queue)/` or `99-organize/`. Agents write to their domain subfolders only.

---

## Environment Variables Reference

See `openclaw/deploy/env.example` for full list. Key groups:

| Group | Vars |
|-------|------|
| LLM | `AZURE_OPENAI_*`, `ANTHROPIC_API_KEY`, `OPENAI_API_KEY` |
| Gateway auth | `OPENCLAW_GATEWAY_TOKEN`, `OPENCLAW_HOOKS_TOKEN` |
| Voice / TTS | `OPENAI_API_KEY` (OpenAI TTS; ElevenLabs removed) |
| Comms | `DISCORD_BOT_TOKEN` |
| Real-time AV | `LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET`, `LIVEKIT_URL` |
| Data | `OMNIDB_URL`, `OMNIDB_DB`, `OMNIDB_USER`, `OMNIDB_PASSWORD` |
| iMessage | `BLUEBUBBLES_PASSWORD` |
| GitHub | `GITHUB_TOKEN` |
| Automation | `N8N_ENCRYPTION_KEY` |
| Static | `NATS_URL`, `LOG_LEVEL`, channel IDs (in `.env.sa`) |

---

## Command Center Overview

UI interfaces have been forked to their own repos. `command-center/` now holds only a README pointer.

**Forked Repos:**
| Repo | Purpose |
|------|---------|
| [One-Mind-OS/swarmclaw](https://github.com/One-Mind-OS/swarmclaw) | Unified command center (Next.js 16, 15 AI providers, LangGraph, connectors) |
| [One-Mind-OS/chowder-iOS](https://github.com/One-Mind-OS/chowder-iOS) | iOS/macOS client (Protocol v3, Live Activity, streaming) |
| [One-Mind-OS/openclaw-office](https://github.com/One-Mind-OS/openclaw-office) | "The Homestead" — 3D agent visualization (Vite, R3F, Zustand) |

**SwarmClaw (local dev):**
```bash
git clone https://github.com/One-Mind-OS/swarmclaw.git && cd swarmclaw && npm install && npm run dev  # http://localhost:3456
```

**Chowder-iOS:**
```bash
git clone https://github.com/One-Mind-OS/chowder-iOS.git && cd chowder-iOS
# Open Chowder.xcodeproj in Xcode
```

- Prefix commits: `feat:`, `fix:`, `secrets:`, `skill:`, `deploy:`, `docs:`, `refactor:`, `chore:`
- Never force-push `main`
- PR title mirrors commit prefix
- Check `journalctl -u openclaw -f` after any deploy to confirm clean startup

---

## Codex — Obsidian Vault (Source of Truth)

The `codex/` directory is a **portable Obsidian vault** tracked in the monorepo. It is the single source of truth for all documentation, knowledge, and planning. Obsidian plugins, themes, and config are all tracked for portability — clone this repo on any machine, open `codex/` in Obsidian, and everything works.

### Codex Structure

```
codex/
├── .obsidian/              ← Full Obsidian config (tracked for portability)
│   ├── plugins/            ← All installed Obsidian plugins
│   ├── themes/             ← All installed themes
│   ├── snippets/           ← Custom CSS snippets
│   └── *.json              ← App config, hotkeys, etc.
├── 00-24 SO/               ← The Brain: AI, tech, infrastructure, architecture
├── 25-49 HP/               ← Self Path: health, skills, personal finance
├── 50-74 LE/               ← Legacy Path: family, home, estate
├── 75-99 GE/               ← Wealth Path: business, ventures, entrepreneurship
├── 99-organize/            ← Staging area (legacy docs, imports, unsorted content)
├── assets/                 ← Images and media
├── ONEMIND-CODEX.md        ← Master document defining the vault
└── CODEX-FRAMEWORK.md      ← Structural rules and conventions
```

### Codex Path Rules

| Context | Path |
|---------|------|
| Local (Mac) | `/Users/zeuslegacy/Documents/onemindos-v0.5/codex/` |
| VPS (Legacy reads) | `/home/onemind/onemind-os/codex/` |
| TOOLS.md reference | Must use `/home/onemind/onemind-os/codex/` |

### How Codex Syncs to VPS

```
Mac (local edits or Obsidian)
    ↓ git push (manual or Obsidian Git plugin)
One-Mind-OS/onemindos-v0.5 (GitHub) ← SOURCE OF TRUTH
    ↓ codex-pull.timer (every 1 min on VPS)
/home/onemind/onemind-os/codex/ ← LEGACY READS THIS PATH
    ↓ codex-watch.service (Legacy writes → auto-commit → GitHub)
```

Legacy writes (memory logs, vault edits) flow back via `codex-watch.service` → GitHub → local `git pull`.

---

## What NOT to Do

- Do not put real credentials in any file that gets committed
- Do not delete or overwrite `SOUL.md` / `USER.md` casually — they define Legacy's identity
- Do not use the old `/opt/onemind/.env` file on the VPS — it is deprecated
- Do not change the codex path in TOOLS.md from `/home/onemind/onemind-os/codex/` — this is where Legacy reads the vault
- Do not save codex files to the old `/Users/zeuslegacy/Documents/onemindcodex-v0.4/` directory — that standalone repo is archived

> **Note**: The local folder may still be named `onemindos-v0.4/` on disk. The version is v0.5 — folder rename happens when ready.
