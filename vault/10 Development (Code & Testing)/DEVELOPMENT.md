# DEVELOPMENT.md — One Mind OS v0.5 Development Plan

> **The unified development roadmap for evolving One Mind OS from v0.4 → v0.5.**
> **Date**: March 16, 2026
> **Author**: Zeus + Copilot
> **Supersedes**: Prior ad hoc plans scattered across BLUEPRINT.md and conversation history.

---

## Table of Contents

1. [v0.5 Vision — What Changes](#1-v05-vision)
2. [Repository Restructure — Fork Strategy](#2-repository-restructure)
3. [SwarmClaw Unification — Feature Extraction Plan](#3-swarmclaw-unification)
4. [iOS App Activation](#4-ios-app-activation)
5. [OpenClaw Studio — Status & Decision](#5-openclaw-studio)
6. [VPS Cleanup — Already Done](#6-vps-cleanup)
7. [Infrastructure Activation — BLUEPRINT Phase 1-2](#7-infrastructure-activation)
8. [Full Service Map — Every Process After v0.5](#8-full-service-map)
9. [Cost Model](#9-cost-model)
10. [Implementation Sequence](#10-implementation-sequence)

---

## 1. v0.5 Vision

### What v0.4 Was

A monorepo holding everything — OpenClaw workspace, codex vault, 6+ vendored dashboard projects, iOS app, infra configs, roadmap docs — all in one repo (4,930 tracked files, 513 MB working directory). It worked for bootstrapping but became bloated:

- 3 cloned community repos sitting in `command-center/` with their own `.git` directories
- 6 different dashboard experiments (Gateway UI, SwarmClaw, Office, Studio, Mission Control, Paperclip)
- No clear separation between "Legacy's brain" and "things we build on top of it"
- Mixed toolchains (Swift, Node.js, Python) all tracked in one repo

### What v0.5 Is

A **lean monorepo** focused on what it should be — Legacy's brain, documentation, and infrastructure — with active development projects forked into their own focused repos.

**Before (v0.4):** 4,930 files, 6 dashboards, 3 foreign `.git` repos embedded

**After (v0.5):** ~2,000 files in the monorepo + 2-3 focused development repos

```
One-Mind-OS/onemindos-v0.5           ← Lean monorepo (brain + docs + infra)
One-Mind-OS/swarmclaw                 ← Fork of swarmclawai/swarmclaw (unified command center)
One-Mind-OS/chowder-iOS               ← Fork of newmaterialco/chowder-iOS (primary iOS/macOS client)
One-Mind-OS/openclaw-office            ← Fork of WW-AI-Lab/openclaw-office ("The Homestead" — 3D visualization)
```

---

## 2. Repository Restructure — Fork Strategy

### The Decision Matrix

| Component | Files | Action | Reason |
|-----------|-------|--------|--------|
| **SwarmClaw** | 710 | **Fork → `One-Mind-OS/swarmclaw`** | Heavy dev planned (6 new pages, three.js, etc.), has upstream (`swarmclawai/swarmclaw`), deploys independently, own Dockerfile/CI |
| **iOS App** | 391 | **Replace with Chowder fork** | Chowder-iOS (541 stars, MIT) has superior Protocol v3 impl. Our modules (Voice, Watch, Camera, Share) will be ported into `One-Mind-OS/chowder-iOS` |
| **OpenClaw Studio** | 350 | **Removed** | Redundant to SwarmClaw — removed from monorepo tracking |
| **openclaw-mission-control** | 762 | **Delete from repo** | Cloned community project (`abhi1693/openclaw-mission-control`), has own `.git`, already removed from VPS |
| **paperclip** | 1,209 | **Delete from repo** | Cloned community project (`paperclipai/paperclip`), has own `.git`, already removed from VPS |
| **Star-Office-UI** | 140 | **Delete from repo** | Cloned community project (`ringhyacinth/Star-Office-UI`), has own `.git`, already removed from VPS |
| **design-reference/** | 194 | **Move into SwarmClaw fork** | Flutter archive, UI specs, theme tokens — all feed into SwarmClaw design work |

### What Each Fork Looks Like

#### `One-Mind-OS/swarmclaw` (fork of `swarmclawai/swarmclaw`)

```
One-Mind-OS/swarmclaw/
├── upstream: swarmclawai/swarmclaw     ← Pull community updates
├── 30+ OneMind-specific files          ← src/app/api/openclaw/*, providers
├── 5 new pages (planned)               ← /system, /soul, /files, /nats, /terminal
├── design-reference/                   ← Moved from monorepo
├── Dockerfile + CI                     ← Independent deployment
└── README.md                           ← OneMind customization docs
```

**Workflow:**
```bash
# Pull upstream community updates
git fetch upstream && git merge upstream/main

# Deploy to VPS
git push origin main  # triggers build
ssh root@104.248.4.212 "cd /home/onemind/swarmclaw && git pull && npm run build && systemctl restart swarmclaw"
```

#### `One-Mind-OS/chowder-iOS` (fork of `newmaterialco/chowder-iOS`)

Chowder is the **primary iOS/macOS client** — a clean, MIT-licensed Swift app with 541 stars and full OpenClaw Gateway Protocol v3 implementation.

```
One-Mind-OS/chowder-iOS/
├── Chowder/
│   ├── ChowderApp.swift        ← App entry point
│   ├── Models/                 ← AgentActivity, BotIdentity, ConnectionConfig, Message, UserProfile
│   ├── Services/               ← ChatService (Protocol v3), KeychainService, LocalStorage
│   ├── ViewModels/             ← ChatViewModel
│   └── Views/                  ← ChatView, SettingsView, ActivityStepRow, etc.
└── README.md
```

**Why Chowder over our ios-app:**
- Full Protocol v3 WebSocket implementation (nonce handshake, scopes, operator mode)
- Streaming responses with thinking step visualization
- Live Activity with tool call + thinking step display
- Identity sync (IDENTITY.md + USER.md parsing)
- Custom avatar support
- Auto-reconnection with health monitoring
- Debug log viewer

**Modules to port from our ios-app into Chowder:**
- Voice: `TalkModeManager`, `VoiceWakeManager`, `TalkOrbOverlay`
- WatchApp: watchOS companion
- Camera: `CameraService`, vision capabilities
- ShareExtension: iOS Share Extension
- Screen: `ScreenRecordService`, `RootCanvas`
- Calendar/Contacts/Location/Motion: Hardware node capabilities

**macOS strategy:** Add macOS target to Chowder (pure SwiftUI, no Catalyst). The official OpenClaw macOS app (`ai.openclaw.mac`, v2026.3.2 with Sparkle auto-updates) can be forked later if needed.

#### `One-Mind-OS/openclaw-studio` — REMOVED

OpenClaw Studio (`grp06/openclaw-studio`) was a community-built Next.js dashboard. **Decision: Removed** — redundant to SwarmClaw. All useful features (agent management, chat, approvals, jobs) are already in SwarmClaw or planned for the 6 new pages.

### Monorepo After Restructure

```
onemindos-v0.5/
├── openclaw/                       ← Legacy's runtime (primary active code)
│   ├── workspaces/                 ← All 9 agent workspaces
│   │   ├── legacy/                 ← Default agent + 35 skills
│   │   ├── spartan/                ← ⚔️  Body & Physical Mastery
│   │   ├── oracle/                 ← 🔮 Mind, Money & Meaning
│   │   ├── eden/                   ← 🌳 Food Growing & Agriculture
│   │   ├── empire/                 ← 👑 Business & Wealth
│   │   ├── haven/                  ← 🌿 Home, Family & Estate
│   │   ├── guardian/               ← 🛡️  Security & Legal
│   │   ├── grid/                   ← ⚙️  Infrastructure Ops
│   │   └── trinity/                ← 🔺 Community & Public
│   └── deploy/                     ← VPS deployment (systemd, scripts, env)
├── codex/                          ← Obsidian vault (source of truth, 2,591 files)
├── infra/                          ← Docker compose, NATS config, VPS scripts
├── roadmap/                        ← Feature planning docs
├── command-center/
│   └── README.md                   ← Pointers to forked repos
├── AGENTS.md                       ← AI agent instructions
├── CLAUDE.md                       ← Claude Code instructions
├── BLUEPRINT.md                    ← Architecture blueprint
├── DEVELOPMENT.md                  ← THIS FILE
├── CHANGELOG.md                    ← Version history
├── ACCESS.md                       ← Credentials & access reference
└── README.md                       ← Project overview
```

### Execution Steps

**Completed:**
```bash
# 1. Forks created on GitHub (via GitHub API)
# One-Mind-OS/swarmclaw     ← fork of swarmclawai/swarmclaw ✅
# One-Mind-OS/chowder-iOS   ← fork of newmaterialco/chowder-iOS ✅

# 2. Removed from monorepo
git rm -r command-center/openclaw-studio           # 350 tracked files ✅
rm -rf command-center/Star-Office-UI               # untracked clone ✅
rm -rf command-center/openclaw-mission-control      # untracked clone ✅
rm -rf command-center/paperclip                     # untracked clone ✅

# 3. command-center/README.md updated with fork pointers ✅
```

**Remaining:**
```bash
# 4. Push OneMind customizations to SwarmClaw fork
#    Clone One-Mind-OS/swarmclaw → copy src/app/api/openclaw/* and other
#    OneMind-specific files from command-center/swarmclaw/ → push to fork
#    Then: git rm -r command-center/swarmclaw

# 5. Remove ios-app from monorepo (reference only — Chowder replaces it)
#    git rm -r command-center/ios-app

# 6. Move design-reference/ into SwarmClaw fork
#    git rm -r command-center/design-reference

# 7. Final commit
git add -A && git commit -m "refactor: v0.5 repo restructure — remove vendored projects, point to forks"
```

---

## 3. SwarmClaw Unification — Feature Extraction Plan

### The Goal

Consolidate all dashboard experiments into **one unified SwarmClaw** — the single command center for One Mind OS. Extract the best features from every community project evaluated, fold OpenClaw Office's 3D visualizations into SwarmClaw, and build 6 new pages.

### Current SwarmClaw State

- **Version**: v0.7.7 (Next.js 16.1.6, React 19.2.3, Tailwind v4)
- **22 existing views**: home, agents, chatrooms, schedules, memory, tasks, approvals, secrets, providers, skills, connectors, webhooks, mcp_servers, knowledge, plugins, usage, wallets, runs, logs, settings, projects, activity
- **140+ API routes** covering all OpenClaw operations
- **62 dependencies** including recharts, zustand, better-sqlite3, ws, qrcode, discord.js, grammy, baileys
- **30+ OneMind-specific files** already customized (openclaw provider, API routes, etc.)

### 5 New Pages to Build

> **Note**: The `/homestead` 3D visualization is now handled by the standalone `One-Mind-OS/openclaw-office` fork ("The Homestead") — a Vite + React Three Fiber app running at `office.onemindos.dev`. Absorbing R3F/three.js into Next.js SSR was architecturally problematic. The fork keeps Office as a focused visualization SPA while SwarmClaw stays focused on management.

#### Page 1: `/system` — VPS System Monitor

**Source**: mudrii (CPU/RAM/disk panels, sub-agent hierarchy tree, cost SVG charts, gateway liveness probes) + tugcantopaloglu (Docker container management, UFW/fail2ban status, audit log)

**What it does**: Real-time VPS health dashboard — CPU, memory, disk, Docker containers, firewall status, service health, NATS connection status.

**API**: New `/api/system/metrics` route polling VPS via SSH or agent.

#### Page 2: `/soul` — Identity Editor

**Source**: VidClaw (SOUL.md, IDENTITY.md, AGENTS.md, USER.md editor with version history and persona templates)

**What it does**: Edit any agent's identity files (SOUL.md, IDENTITY.md, AGENTS.md, USER.md, HEARTBEAT.md, BOOT.md, TOOLS.md) with a rich markdown editor. Version history via git log. Persona templates for quickly creating new agent identities.

**API**: Already exists — `/api/openclaw/agent-files/route.ts`

#### Page 3: `/files` — Workspace File Browser

**Source**: TenacitOS (file browser, live preview, in-browser edit) + tugcantopaloglu (path traversal protection, auto-backup)

**What it does**: Browse and edit any file in the OpenClaw workspace. Live preview for markdown. Syntax highlighting for JS/JSON. Path traversal protection (never escape workspace boundary).

**API**: Already exists — `/api/openclaw/directory/`, `/api/openclaw/agent-files/`

#### Page 4: `/nats` — NATS JetStream Monitor

**Source**: New build (unique to this stack)

**What it does**: Stream list with message counts and storage. Subject viewer with live event tailing. Sub-agent task log. legacy-bridge connection status. Manual publish for testing.

**Dependencies needed**:
```bash
npm install nats  # NATS client for Node.js
```

#### Page 5: `/terminal` — Safe Web Terminal

**Source**: Clawtrol + tugcantopaloglu (allowlist pattern, SSE streaming, re-auth gate)

**What it does**: Read-only web terminal for running safe commands (logs, status checks, NATS queries). Command allowlist prevents destructive operations. SSE streaming for live output. Re-authentication gate before execution.

### 7 Enhancements to Existing Views

| Enhancement | Target View | Source Project | Details |
|-------------|-------------|----------------|---------|
| **Global floating STT mic** | `app-layout.tsx` | actionagentai | Speech-to-text input, Cmd+Shift+M hotkey, floating button in layout |
| **Sub-agent cost breakdown** | `/usage` | mudrii | 7d/30d/all-time tabs, per-agent token/cost breakdown using existing recharts |
| **Animated flow diagram** | `/activity` | ClawMetry | Visual flow of agent interactions, tool calls, message routing |
| **WhatsApp QR login** | `/connectors` | actionagentai | QR code for WhatsApp Web pairing (`qrcode` package already installed) |
| **Cron heatmap** | `/schedules` | TenacitOS + mudrii | Activity heatmap showing when cron jobs fire, success/failure colors |
| **Soul version history** | `/agents` agent sheet | VidClaw | Git log of SOUL.md changes per agent, diff viewer |
| **AI chat about system** | `/home` | mudrii | "Ask about your system" chat panel on home page |

### How to Add a New View (Pattern)

```typescript
// 1. src/types/index.ts — extend AppView union
export type AppView = 'home' | 'agents' | ... | 'system' | 'soul' | 'files' | 'nats' | 'terminal'

// 2. view-route-paths.ts — add route
system: '/system',
soul: '/soul',
files: '/files',
nats: '/nats',
terminal: '/terminal',

// 3. src/components/layout/app-layout.tsx — add VIEW_LABELS + import + render
// Add to VIEW_LABELS map
// Add to component render switch
```

### OpenClaw Office — Forked as "The Homestead"

Instead of absorbing Office's R3F/three.js into SwarmClaw (Next.js SSR conflicts), Office is forked as a standalone visualization app:

- **Repo**: `One-Mind-OS/openclaw-office` (fork of `WW-AI-Lab/openclaw-office`)
- **Stack**: Vite 6, React 19, React Three Fiber, Zustand 5, Tailwind 4
- **Live**: `office.onemindos.dev` (:5180)
- **Customization plan**: Solarpunk theme, 9 custom agent avatars, NATS event feed panel, rename to "The Homestead"

---

## 4. iOS App Activation

### iOS App Activation — Chowder as Primary Client

The original ios-app/ (314 Swift files, 21+ modules) is superseded by `One-Mind-OS/chowder-iOS` (fork of `newmaterialco/chowder-iOS`, 541 stars, MIT license). Chowder has a cleaner Protocol v3 implementation with streaming, thinking steps, and Live Activity. Our unique modules (Voice, Watch, Camera, Share) will be ported into the Chowder fork.

### Activation Priority (P0)

Wire Chowder to production gateway:

```
Gateway URL: wss://gateway.onemindos.dev
Token: 6d74d6a39d9b3eba5b2104f17919d48ed6067401d1e63168d8a0c7afd400cd73
```

This single connection gives us:
- Voice chat with Legacy (TalkMode)
- Text chat with all 9 agents
- Canvas rendering on device
- Watch app connectivity
- Dynamic Island live activities
- Share Extension for sending content to Legacy

### macOS Target Addition

The `project.yml` currently targets iOS only. Adding macOS:

```yaml
# In project.yml targets section
OpenClaw-Mac:
  type: application
  platform: macOS
  deploymentTarget: "15.0"
  sources:
    - path: Sources
      excludes:
        - "Motion/**"  # No CoreMotion on Mac
  # ... same signing, capabilities
```

This gives a native Mac app using the same codebase — no Catalyst, no Electron, pure SwiftUI.

### Enhanced Features to Build (in the fork)

| Feature | Module | Description |
|---------|--------|-------------|
| **Ambient overlay** | New `Overlay/` module | Transparent always-on-top window showing agent status, slides in from edge |
| **SmolVLM camera** | `Camera/` extension | Periodic webcam capture → vision analysis → NATS publish |
| **NATS events** | New `NATS/` module | Subscribe to live NATS events, show notifications |
| **MediaPipe gestures** | New `Gestures/` module | Hand gesture recognition for controlling Legacy |
| **Menu bar app** | New `MenuBar/` module | Persistent menu bar icon with quick actions |

---

## 5. OpenClaw Studio — REMOVED

**What it was**: A community-built (`grp06/openclaw-studio`) Next.js web dashboard. Clean agent management, chat, approvals, jobs.

**Decision**: Removed from monorepo (350 files). Redundant to SwarmClaw. All useful features are already in SwarmClaw or planned for the 6 new pages. VPS systemd service was also removed.

---

## 6. VPS Cleanup — Already Done (v0.4 → v0.5 Transition)

These services and files were already removed from the VPS during this session:

| Removed | Type | Files | Reason |
|---------|------|-------|--------|
| **Star Office UI** | systemd + files | `/home/onemind/Star-Office-UI/` | Community project, replaced by Office fork |
| **Mission Control** | Docker (5 containers) + files | `/home/onemind/openclaw-mission-control/` | Community project, duplicate functionality |
| **Paperclip** | systemd + files | `/home/onemind/paperclip/` | Community project, not used |
| **OpenClaw Studio** | systemd service | Service only | Replaced by SwarmClaw |

**Caddy cleaned**: Only 4 domain blocks remain:
- `gateway.onemindos.dev` → OpenClaw Gateway (:18789)
- `swarm.onemindos.dev` → SwarmClaw (:3456)
- `office.onemindos.dev` → OpenClaw Office "The Homestead" (:5180) — forked to `One-Mind-OS/openclaw-office`
- `nats.onemindos.dev` → NATS Dashboard

**Binaries installed**: `blu` (blucli v0.1.3) and `camsnap` (v0.2.0) — resolved missing binary errors.

---

## 7. Infrastructure Activation — BLUEPRINT Phase 1-2

These are the immediate infrastructure improvements from BLUEPRINT.md that apply regardless of the repo restructure.

### Phase 1: Foundation (This Week)

| Step | Component | Status | Action |
|------|-----------|--------|--------|
| 1.1 | Heartbeat activation | **Not done** | Enable in `openclaw.json` — `"every": "30m"`, active 07:00-23:00 |
| 1.2 | Lossless Claw plugin | **Not done** | `openclaw plugins install @martian-engineering/lossless-claw` |
| 1.3 | NATS JetStream streams | **Not done** | Create 6 persistent streams (EVENTS, SYSTEM, NEXUS, DIGITAL, PHYSICAL, COMMS) |
| 1.4 | VPS health → NATS | **Not done** | Cron job publishing `system.health.vps` every 5 min |
| 1.5 | Obsidian skills (kepano) | **Not done** | Install `obsidian-markdown`, `obsidian-bases`, `json-canvas`, `obsidian-cli`, `defuddle` |
| 1.6 | SwarmClaw upgrade | **Not done** | v0.7.7 → latest (Missions, Orchestrator agents) |
| 1.7 | Skill archival | **Not done** | Move 12 unused skills to `archive/` |

### Phase 2: Ambient Awareness (Week 2)

| Step | Component | Action |
|------|-----------|--------|
| 2.1 | legacy-pulse | Deploy as systemd service — NATS subscriber → writes AMBIENT_STATE.md |
| 2.2 | OpenViking Context DB | Deploy on VPS — tiered context loading, hybrid search |
| 2.3 | HEARTBEAT.md update | Add ambient awareness instructions (read AMBIENT_STATE.md) |
| 2.4 | legacy-eyes | Deploy SmolVLM-256M on MacBook Pro — vision digest → NATS |

### Phase 3: Desktop + Mobile (Week 3-4)

| Step | Component | Action |
|------|-----------|--------|
| 3.1 | iOS app activation | Wire to `wss://gateway.onemindos.dev` with token |
| 3.2 | macOS target | Add macOS platform to project.yml |
| 3.3 | Office fork customization | Theme to "The Homestead", custom 9-agent avatars, NATS event feed |
| 3.4 | SwarmClaw `/soul` | Identity editor with version history |

### Phase 4: Voice + GPU (Month 2)

| Step | Component | Action |
|------|-----------|--------|
| 4.1 | Chatterbox TTS | Self-hosted voice synthesis (replaces OpenAI TTS) |
| 4.2 | GPU Droplet (on-demand) | RTX 4000 Ada, 2hr/day for SmolVLM + local inference |
| 4.3 | NixClaw fork | Gemini Live voice bridge for mobile |

---

## 8. Full Service Map — Every Process After v0.5

### Currently Running

| Service | Host | Port | Status |
|---------|------|------|--------|
| `openclaw.service` (9 agents) | VPS | 18789 | **LIVE** |
| NATS JetStream (Docker) | VPS | 4222 | **LIVE** |
| PostgreSQL/pgvector (Docker) | VPS | 5432 | **LIVE** |
| legacy-bridge (Docker) | VPS | — | **LIVE** |
| nats-ingester (Docker) | VPS | — | **LIVE** |
| SwarmClaw | VPS | 3456 | **LIVE** (v0.7.7) |
| OpenClaw Office "The Homestead" | VPS | 5180 | **LIVE** (forked to `One-Mind-OS/openclaw-office`) |
| Caddy (reverse proxy) | VPS | 443 | **LIVE** |
| Codex sync (pull + watch) | VPS | — | **LIVE** |
| Memory search (hybrid) | VPS | — | **LIVE** |
| BlueBubbles (iMessage) | MacBook | 1234 | **LIVE** |

### To Deploy (Phase 1-2)

| Service | Host | Port | Phase |
|---------|------|------|-------|
| Lossless Claw (plugin) | VPS | — | 1 |
| Heartbeat (config) | VPS | — | 1 |
| legacy-pulse | VPS | — | 2 |
| OpenViking | VPS | 8900 | 2 |
| legacy-eyes | MacBook | — | 2 |

### Custom Domains

| Domain | Target | Status |
|--------|--------|--------|
| `gateway.onemindos.dev` | OpenClaw Gateway (:18789) | **LIVE** |
| `swarm.onemindos.dev` | SwarmClaw (:3456) | **LIVE** |
| `office.onemindos.dev` | OpenClaw Office "The Homestead" (:5180) | **LIVE** (forked to `One-Mind-OS/openclaw-office`) |
| `nats.onemindos.dev` | NATS Dashboard | **LIVE** |

---

## 9. Cost Model

### Current (v0.4)

| Item | Monthly Cost |
|------|-------------|
| VPS (CPU Droplet) | ~$50-100 |
| Azure OpenAI (GPT-5.2, GPT-4o) | ~$150-200 |
| OpenAI (TTS, embeddings) | ~$20-30 |
| **Total** | **~$220-330** |

### After v0.5 Phase 1-2

| Item | Monthly Cost | Delta |
|------|-------------|-------|
| VPS (same) | ~$50-100 | — |
| Azure OpenAI (reduced by OpenViking) | ~$80-120 | -$70-80 |
| OpenAI (TTS, embeddings) | ~$20-30 | — |
| Heartbeat tokens (gpt-4o-mini) | ~$90-180 | New |
| Lossless Claw summaries | ~$10-15 | New |
| SmolVLM vision (local) | $0 | Free |
| **Total** | **~$250-445** | +$30-115 |

### Long-Term Target (Mac Studio at home)

| Item | Monthly Cost |
|------|-------------|
| VPS for NATS + Gateway relay | ~$20-40 |
| All inference local (Ollama) | $0 |
| Self-hosted TTS/STT | $0 |
| **Total** | **~$20-40** |

---

## 10. Implementation Sequence

### Week 1 — Repo Restructure + Foundation

- [ ] Create `One-Mind-OS/swarmclaw` fork on GitHub, push current code, add upstream remote
- [ ] Create `One-Mind-OS/openclaw-ios` repo on GitHub, push iOS app code
- [ ] Create `One-Mind-OS/openclaw-studio` fork on GitHub, push current code
- [ ] Remove vendored copies from monorepo (swarmclaw, ios-app, openclaw-studio, mission-control, paperclip, Star-Office-UI)
- [ ] Update monorepo `command-center/README.md` with pointers to forked repos
- [ ] Rename repo references from v0.4 → v0.5 across AGENTS.md, CLAUDE.md, README.md, ACCESS.md, codex/CLAUDE.md
- [ ] Enable heartbeat in `openclaw.json`
- [ ] Install Lossless Claw plugin
- [ ] Create 6 NATS JetStream streams
- [ ] Install Obsidian skills (kepano)

### Week 2 — SwarmClaw Phase 1 + Ambient

- [ ] Add global STT mic to SwarmClaw `app-layout.tsx`
- [ ] Build `/soul` page (identity editor)
- [ ] Add sub-agent cost tabs to `/usage`
- [ ] Deploy legacy-pulse (NATS → AMBIENT_STATE.md)
- [ ] Deploy OpenViking Context DB
- [ ] Update HEARTBEAT.md for ambient awareness

### Week 3 — iOS + SwarmClaw Phase 2

- [ ] Wire iOS app to production gateway (`wss://gateway.onemindos.dev`)
- [ ] Add macOS target to project.yml
- [ ] Build `/system` page (VPS monitor)
- [ ] Build `/nats` page (JetStream monitor)
- [ ] Build `/files` page (workspace browser)

### Week 4 — Office Fork + Terminal + Vision

- [ ] Customize `One-Mind-OS/openclaw-office` theme ("The Homestead", solarpunk)
- [ ] Add 9 custom agent avatars to Office fork
- [ ] Add NATS event feed panel to Office fork
- [ ] Build SwarmClaw `/terminal` page (safe web terminal)
- [ ] Deploy legacy-eyes (SmolVLM on MacBook Pro)

### Month 2 — Voice + GPU

- [ ] Deploy Chatterbox TTS (self-hosted)
- [ ] On-demand GPU Droplet setup
- [ ] NixClaw fork for Gemini Live voice
- [ ] iOS app enhanced features (ambient overlay, menu bar, NATS events)

---

## Key Decisions

| Decision | Answer | Reasoning |
|----------|--------|-----------|
| Fork SwarmClaw? | **Yes → `One-Mind-OS/swarmclaw`** | Heavy dev planned, has upstream, independent deployment |
| Fork iOS app? | **Yes → `One-Mind-OS/openclaw-ios`** | Different toolchain (Swift), App Store deployment, needs macOS CI |
| Fork OpenClaw macOS app? | **Not yet — extend iOS app to macOS first** | Our iOS app already has the full gateway stack; adding macOS target is faster |
| Fork Studio? | **Yes → `One-Mind-OS/openclaw-studio`** | Preserve for reference, no active dev planned |
| Keep design-reference? | **Move to SwarmClaw fork** | Feeds into SwarmClaw design work directly |
| GPU now? | **No. Phase 2+ on-demand.** | API costs don't justify 24/7 GPU yet |
| Which command center? | **SwarmClaw (unified)** | One dashboard to rule them all. Office, Studio absorbed. |
| Heartbeat frequency? | **30 min + P0 webhook** | Balance cost vs awareness. P0 alerts bypass immediately. |
| Cloud or local? | **Cloud-first, migrate local incrementally** | VPS for uptime. MacBook for vision, iMessage, dev. |

---

*This is the v0.5 development plan. Everything in one place. Start on your word.*
