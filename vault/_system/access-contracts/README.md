# Access Contracts — OneMind Ecosystem

Defines how each class of consumer reads from and writes to the OneMind system.
This is the authoritative reference for integration design.

For the machine-readable version, see `ecosystem-manifest.json` at the repo root.

---

## Consumer Classes

1. [Human — Desktop (Zeus)](#1-human--desktop-zeus)
2. [Human — Mobile (Zeus)](#2-human--mobile-zeus)
3. [Human — Browser (Zeus + public)](#3-human--browser-zeus--public)
4. [AI Coding Agents (Claude Code, Copilot, Cursor, etc.)](#4-ai-coding-agents)
5. [OpenClaw Agents (Legacy, Spartan, Oracle, Eden, Empire)](#5-openclaw-agents)
6. [Automation (N8N, cron, webhooks)](#6-automation-n8n-cron-webhooks)
7. [Physical / Robots / IoT](#7-physical--robots--iot)
8. [Community / Public](#8-community--public)

---

## 1. Human — Desktop (Zeus)

**Primary tool:** Obsidian (knowledge), VS Code (code), Finder (media)

| Action | How |
|--------|-----|
| Write / read knowledge | Obsidian → `onemind-codex` |
| Browse media assets | Finder / VS Code → `onemind-media`, `onemind-library` |
| Code work | VS Code + GitHub Copilot → any git repo |
| Tasks | Todoist → synced to codex via N8N |
| AI assistance (coding/research) | Claude Code / GitHub Copilot |
| AI assistance (life/planning) | Legacy agent (via any channel) |

**Entry point:** Open `~/Documents/onemind-codex` as Obsidian vault

**Read paths:**
- `~/Documents/onemind-codex/` (brain — all knowledge)
- `~/Documents/onemind-media/` (media assets, 278 GB)
- `~/Documents/onemind-library/` (document library, 66 GB)

**Write paths:**
- `~/Documents/onemind-codex/` → git commit → GitHub → auto-syncs to VPS

---

## 2. Human — Mobile (Zeus)

**Primary tool:** onemind-mobile app

| Action | How |
|--------|-----|
| Quick capture (thoughts, tasks) | onemind-mobile → pushes note to `00-24 SO/24 Inbox (Queue)/` |
| Read knowledge | onemind-mobile OR onemind-codex-live (public site) |
| View tasks | Todoist mobile app |
| Voice capture | onemind-mobile voice → transcription → inbox |

**Entry point:** onemind-mobile app

**Write path:** Mobile app → API → creates file in codex inbox OR Todoist task

---

## 3. Human — Browser (Zeus + public)

| Audience | URL | Access |
|----------|-----|--------|
| Zeus (app) | onemind-web domain | Full authenticated access |
| Public | onemind-codex-live domain | Published knowledge only |
| Community | github.com/OneMind-OS | Template + open source content |

---

## 4. AI Coding Agents

**Includes:** Claude Code, GitHub Copilot, Cursor, Gemini CLI, Windsurf, Codex

**Session start protocol:**
```
1. Read AGENTS.md (repo root)
2. Read _codex/INTERFACE.md
3. Load ALL files in _codex/system/ (pinned context — always load)
4. Read relevant domain README if working in a specific area
5. Check _codex/skills/ for applicable skill instructions
```

**Skill routing:**
- Vault operations → `_codex/skills/` (search, capture, organize, reflect, defrag, etc.)
- General AI tasks → `~/Documents/Personal_AI_Infrastructure/Packs/`
- Runtime skills → `.claude/skills/`

**Read access:**
- Full read access to entire `onemind-codex`
- Read `onemind-media/` and `onemind-library/` when working with assets (Mac-local only)

**Write access:**
- Direct file create/edit in `onemind-codex`
- Always follow frontmatter conventions (`_codex/CONVENTIONS.md`)
- New content → inbox first, then proper domain
- No code in knowledge domain folders → code goes in `_system/`

**Secrets:** Reference 1Password via `op://aaeyris5z44pdqspk277u2oppi/...` — never hardcode

**Commit convention:**
```
[domain] Brief description

Co-authored-by: <agent-name> <agent@onemind.os>
```

---

## 5. OpenClaw Agents

**Agents:** Legacy, Spartan, Oracle, Eden, Empire

**Runtime:** VPS at `/home/onemind/onemind-os/codex/`

**Sync flow:**
```
Mac (Zeus writes) → git push → GitHub
    ↓ cron: every 1 minute
VPS auto-pull → /home/onemind/onemind-os/codex/
    ↓ agents read from filesystem
Agents write files → git commit + push → GitHub
    ↓ Mac auto-pull hook
Obsidian reflects agent writes
```

**Domain routing by agent:**

| Agent | Reads | Writes to |
|-------|-------|-----------|
| Legacy | All | Any domain (central coordinator) |
| Spartan | HP 27-28 + SO | HP 27-28 (body/fitness notes, progress) |
| Oracle | HP 25-26, 29-34 + SO | HP 25-26, 29-34 (mind, money, meaning) |
| Eden | LE + SO | LE (family, home, sustainability) |
| Empire | GE + SO | GE (business, ventures, wealth) |

**Agent workspaces:** `_system/openclaw/` — each agent has a workspace config here

**Heartbeat behavior:**
- On each heartbeat: load `_codex/system/` + domain README + relevant recent notes
- Agent can promote/demote files to/from `_codex/system/` based on usage frequency
- New captures go to inbox; agents process inbox and route to proper domains

**What agents CANNOT do without human approval:**
- Delete any file (move to archive instead)
- Rename domain group folders
- Modify `ONEMIND-CODEX.md`
- Commit secrets or credentials

---

## 6. Automation (N8N, cron, webhooks)

**Runtime:** VPS or cloud-hosted N8N instance

**Reading the codex:**
```
GitHub API → GET /repos/OneMind-OS/onemind-codex/contents/{path}
Returns: base64-encoded file content + SHA
```

**Writing to the codex:**
```
GitHub API → PUT /repos/OneMind-OS/onemind-codex/contents/{path}
Body: { message, content (base64), sha (if updating existing) }
→ Triggers auto-pull on VPS → Agents see the update
```

**Reading media/library (Mac-local tasks only):**
```
Direct filesystem access on Mac
~/Documents/onemind-media/
~/Documents/onemind-library/
```

**Event bus (NATS JetStream):**
```
Subject pattern:  onemind.{domain}.{subdomain}.{event}
Examples:
  onemind.hp.body.workout-logged
  onemind.ge.ventures.new-lead
  onemind.so.agents.heartbeat
  onemind.physical.{device}.{event}
```

**Todoist sync:**
- Bidirectional: Todoist tasks ↔ codex notes in relevant domain
- Managed by N8N workflow in `_system/infrastructure/`

---

## 7. Physical / Robots / IoT

**Current status:** Design contract only — not yet fully implemented

**Query pattern (read):**
```
Device/Robot → POST https://api.onemind.os/query
Body: { agent: "legacy", domain: "HP", intent: "what should I do now?" }
→ Legacy agent reads codex → returns context-aware response
```

**Event pattern (write):**
```
Device publishes to NATS: onemind.physical.{device-id}.{event-type}
Example: onemind.physical.home-assistant.motion-detected
         onemind.physical.scale.weight-logged
         onemind.physical.calendar.event-starting
→ Legacy agent subscribes → takes action → may write to codex
```

**Device registry:** Define devices in `00-24 SO (Sovereign Operations)/01 Infrastructure (Grid & APIs)/`

**Supported integrations (planned):**
- Home Assistant → NATS bridge
- Wearables (Apple Watch, Oura) → health data → HP domain
- Smart home devices → LE domain events
- Physical robots → custom NATS bridge

**Security:** All physical devices must authenticate via JWT token issued by OpenClaw.
No device has direct write access to the codex — all writes go through Legacy agent.

---

## 8. Community / Public

**What's available:**

| Resource | URL | Purpose |
|----------|-----|---------|
| OneMind Codex Template | github.com/OneMind-OS/onemind-codex-template | Fork + build your own codex |
| OneMind Codex Live | onemind-codex-live domain | Zeus's published knowledge |
| PAI (Personal AI Infrastructure) | github.com/danielmiessler/Personal_AI_Infrastructure | Upstream AI skill packs |

**What's NOT available to community:**
- `onemind-codex` (private — contains personal data)
- `onemind-media` and `onemind-library` (local only)
- OpenClaw agent configs and workspaces (private)

---

## Access Contract Summary Matrix

| Consumer | Read Brain | Write Brain | Read Media | Write Media | Event Bus |
|----------|-----------|------------|------------|-------------|-----------|
| Zeus (Desktop) | ✅ Direct | ✅ Direct | ✅ Direct | ✅ Direct | — |
| Zeus (Mobile) | ✅ App | ✅ Inbox only | — | — | — |
| Zeus (Browser) | ✅ App | ✅ App | — | — | — |
| AI Coding Agent | ✅ Direct | ✅ Direct | ✅ Direct (Mac) | ❌ | — |
| OpenClaw Agent | ✅ VPS filesystem | ✅ git push | ❌ | ❌ | ✅ NATS |
| N8N Automation | ✅ GitHub API | ✅ GitHub API | ✅ Direct (Mac) | ❌ | ✅ NATS |
| Physical/Robot | ❌ | ❌ (via Legacy) | ❌ | ❌ | ✅ NATS |
| Community | ✅ Public site | ❌ | ❌ | ❌ | — |

---

*Updated: 2026-03-28 | See `ECOSYSTEM.md` at repo root for the full architecture map.*
