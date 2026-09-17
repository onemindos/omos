---
id: ecosystem-map-20260328
title: "OneMind Ecosystem Map"
type: doctrine
status: active
created: 2026-03-28
modified: 2026-03-28
owner: zeus
pinned: false
tags: [type/doctrine, domain/so, status/active]
ai_summary: "Universal map of the full OneMind ecosystem — all repos, layers, access patterns for humans, AI agents, and automation."
---

# OneMind Ecosystem Map

> **Read this to understand how everything connects.**
> If you are an AI agent, robot, or automation entering the ecosystem, start here.
> If you are Zeus, this is the master index of every system you own.

---

## The Equation

```
HUMAN (Zeus) + AI (Legacy) = ONE MIND
```

The ecosystem is a **5-layer stack**. Each layer has a clear role and does not duplicate what the layer below it provides.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│  LAYER 5 — CONSUMER APPS                                     │
│  onemind-web (Next.js)  ·  onemind-mobile                   │
│  Access: End users, family, clients                          │
└──────────────────────────────┬──────────────────────────────┘
                               │ reads knowledge/assets via API
┌──────────────────────────────▼──────────────────────────────┐
│  LAYER 4 — PUBLIC SURFACES                                   │
│  onemind-codex-live  (published Quartz site on Cloudflare)  │
│  onemind-codex-template  (community fork, stripped of data) │
│  Personal_AI_Infrastructure  (PAI upstream skill packs)      │
│  Access: Community, public web                               │
└──────────────────────────────┬──────────────────────────────┘
                               │ pulls skills from PAI, publishes to live
┌──────────────────────────────▼──────────────────────────────┐
│  LAYER 3 — AI SKILL INFRASTRUCTURE                           │
│  PAI/Packs/  (13 general-purpose skill packs — upstream)    │
│  _codex/skills/  (personal skill overrides — private)        │
│  .claude/skills/  (Claude Code runtime — local)              │
│  Access: AI agents, Claude Code, OpenClaw                    │
└──────────────────────────────┬──────────────────────────────┘
                               │ skills inform + assets reference
┌──────────────────────────────▼──────────────────────────────┐
│  LAYER 2 — ASSET VAULTS                                      │
│  onemind-media  (278 GB — domain-organized media)           │
│  onemind-library  (66 GB — documents + image library)       │
│  Access: AI agents (read-only), Zeus (read/write), robots    │
└──────────────────────────────┬──────────────────────────────┘
                               │ all knowledge lives here
┌──────────────────────────────▼──────────────────────────────┐
│  LAYER 1 — THE BRAIN (single source of truth)                │
│  onemind-codex  (this repo — private Obsidian vault)        │
│  Access: Zeus, Legacy, all AI agents (read+write)            │
│  ├── 00-24 SO  (Sovereign Operations — the platform)        │
│  ├── 25-49 HP  (Holistic Performance — Zeus the human)      │
│  ├── 50-74 LE  (Legacy Evolution — family + foundation)     │
│  └── 75-99 GE  (Generational Empire — work + wealth)      │
└─────────────────────────────────────────────────────────────┘
```

---

## Complete Repo & Folder Inventory

| Name | Type | Local Path | Privacy | Size | Role |
|------|------|-----------|---------|------|------|
| **onemind-codex** | Git + Obsidian vault | `~/Documents/onemind-codex` | 🔒 Private | ~50 MB | **THE BRAIN** — single source of truth |
| **onemind-media** | Local folder (no git) | `~/Documents/onemind-media` | 🔒 Private | 278 GB | Domain-organized media vault |
| **onemind-library** | Local folder | `~/Documents/onemind-library` | 🔒 Private | 66 GB | Document + image library |
| **onemind-web** | Git repo (Next.js) | `~/Documents/onemind-web` | 🔒→🌐 | — | Consumer web application |
| **onemind-mobile** | Git repo | `~/Documents/onemind-mobile` | 🔒→🌐 | — | Mobile application |
| **onemind-codex-live** | Git repo (Quartz/Cloudflare) | `~/Documents/onemind-codex-live` | 🌐 Public | — | Published knowledge site |
| **onemind-codex-template** | Git repo | `~/Documents/onemind-codex-template` | 🌐 Public | — | Community fork template |
| **Personal_AI_Infrastructure** | Git repo (PAI) | `~/Documents/Personal_AI_Infrastructure` | 🌐 Public | 678 MB | Upstream AI skill packs (Daniel Miessler) |
| `Organize/` | Local scratch | `~/Documents/Organize` | Temp | — | Delivery scratch space — ephemeral |

---

## Skill System Hierarchy

Three skill layers compose from general → personal → runtime:

```
PAI/Packs/  ←  upstream community skills  (13 packs — general purpose)
    ↓ installed/adapted into
_codex/skills/  ←  personal vault skills  (8 skills — Zeus-specific)
    ↓ loaded by
.claude/skills/  ←  Claude Code runtime  (session-level)
```

**Rule:** Never duplicate a PAI skill into _codex/skills unless you're personalizing it.
PAI is upstream. The codex is the personalization layer.

---

## Agent Roster

| Agent | Runtime | Domain | Read Path | Write Path |
|-------|---------|--------|-----------|------------|
| **Legacy** | OpenClaw (VPS) | All | VPS: `/home/onemind/onemind-os/codex/` | Push → GitHub → auto-pull |
| **Spartan** | OpenClaw | HP (fitness/body) | Same as Legacy | Same as Legacy |
| **Oracle** | OpenClaw | HP (mind/money) | Same as Legacy | Same as Legacy |
| **Eden** | OpenClaw | LE (family/home) | Same as Legacy | Same as Legacy |
| **Empire** | OpenClaw | GE (business) | Same as Legacy | Same as Legacy |
| **Claude Code / Copilot** | Local (Mac) | Any | Direct file access | Direct file edit + git push |
| **Automation (N8N)** | VPS / cloud | Any | Webhook + GitHub API | Webhook → create/update file |

All agents read the same codex. No agent has a private knowledge silo.

---

## Access Patterns by Consumer

### 🧠 Zeus (Human — Desktop)
- **Knowledge writing/reading:** Obsidian → `onemind-codex`
- **Media triage:** Finder / Claude Code → `onemind-media`, `onemind-library`
- **Code work:** VS Code / GitHub Copilot → any git repo
- **Task management:** Todoist (synced to codex via N8N)
- **Entry point:** Open `onemind-codex` as Obsidian vault

### 📱 Zeus (Human — Mobile)
- **Quick capture:** `onemind-mobile` app → pushes to `24 Inbox`
- **Reference reading:** `onemind-codex-live` (public site) or mobile app
- **Entry point:** `onemind-mobile` app

### 🌐 Zeus (Human — Browser)
- **App interface:** `onemind-web`
- **Public knowledge:** `onemind-codex-live`
- **Entry point:** domain URLs

### 🤖 AI Coding Agents (Claude Code, GitHub Copilot, Cursor, Gemini CLI)
- **Start:** Read `AGENTS.md` → then `_codex/INTERFACE.md`
- **Context:** Load `_codex/system/` files (always-loaded pinned context)
- **Skills:** `_codex/skills/` for vault operations; `PAI/Packs/` for general tasks
- **Entry point:** `AGENTS.md` at repo root

### 🔁 OpenClaw Agents (Legacy, Spartan, Oracle, Eden, Empire)
- **Runtime:** VPS at `/home/onemind/onemind-os/codex/` (auto-pulled from GitHub)
- **Read:** Direct filesystem read of markdown files
- **Write:** Push commits to GitHub; auto-synced to Mac via pull cron
- **Workspaces:** `_system/openclaw/`
- **Entry point:** `AGENTS.md` + `_codex/system/` pinned context

### ⚙️ Automation (N8N, cron jobs, webhooks)
- **Read codex:** GitHub API → read file by path
- **Write codex:** GitHub API → create/update file → triggers auto-pull on VPS
- **Read media:** Direct filesystem (if running on Mac) or NAS mount
- **Events:** NATS JetStream topics (see `_system/infrastructure/`)
- **Entry point:** `ecosystem-manifest.json` at codex root

### 🏠 Physical / Robots (Home Assistant, IoT, future)
- **Query:** POST to OpenClaw API → Legacy agent → reads codex → responds
- **Data in:** NATS topic `onemind.physical.{device}.{event}`
- **Data out:** NATS topic `onemind.physical.{device}.{command}`
- **Entry point:** OpenClaw API endpoint (see `_system/infrastructure/openclaw/`)

### 👥 Community / Public
- **Framework:** `onemind-codex-template` (fork and personalize)
- **Knowledge:** `onemind-codex-live` (published notes)
- **AI skills:** `Personal_AI_Infrastructure` (PAI packs)
- **Entry point:** GitHub repos listed above

---

## PAI Relationship

**Personal_AI_Infrastructure** is [Daniel Miessler's](https://github.com/danielmiessler) open-source framework.
Zeus uses it as the **upstream AI skill infrastructure layer**.

| PAI Role | OneMind Role |
|----------|-------------|
| General-purpose skill packs | Personal skill overrides |
| Community-maintained tools | Zeus-specific adaptations |
| Platform-agnostic | Tailored to Zeus's stack |
| Public, MIT licensed | Private, personal |

**Do not** merge PAI into onemind-codex. Keep it as a separate upstream dependency.
Update PAI from the community. Adapt skills into `_codex/skills/` when personalizing.

---

## Naming Conventions (Canonical — March 2026)

| ❌ Old / Wrong | ✅ Current / Correct |
|---------------|---------------------|
| "quadrant" | "domain" or "domain group" |
| "UI" (for the SO domain) | "SO" or "Sovereign Operations" |
| "onemind-os" | "onemind-codex" |
| "four quadrants" | "3 domains + 1 operating layer" |
| "Your Business" | "Your Empire" (for GE) |
| "Your Family" | "Your Foundation" (for LE) |

> Note: `onemind-media` folder names still use "UI" in folder names (legacy). Do not rename — too risky for 278 GB of content. The folders function correctly; only the label is outdated.

---

## Machine-Readable Manifest

For automated consumers, a full machine-readable manifest is available at:

```
onemind-codex/ecosystem-manifest.json
```

This JSON file contains all paths, roles, privacy levels, and access endpoints in a parseable format.

---

## Sync Flow

```
Zeus writes in Obsidian (Mac)
    ↓ git commit + push
GitHub: OneMind-OS/onemind-codex
    ↓ VPS auto-pulls every 1 minute
Legacy + agents read from VPS filesystem
    ↓ agents write back (create/update notes)
GitHub receives push from VPS
    ↓ Mac auto-pulls
Obsidian reflects agent writes
```

Assets (onemind-media, onemind-library) do NOT sync to VPS — too large.
Assets are accessed by agents only when running locally (Mac-side Claude Code/Copilot).

---

*Part of the OneMind OS. See `ONEMIND-CODEX.md` for the framework doctrine.*


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
