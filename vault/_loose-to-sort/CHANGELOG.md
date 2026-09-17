# Changelog

All notable changes to OneMind OS will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.0] - 2026-03-16

### Added
- **DEVELOPMENT.md** — Comprehensive development plan for v0.5 (repo restructure, SwarmClaw unification, iOS activation)
- **Repository restructure plan** — Fork SwarmClaw, Chowder-iOS, and OpenClaw Office into separate repos; delete vendored community projects
- **SwarmClaw unification plan** — 5 new pages (/system, /soul, /files, /nats, /terminal) + 7 enhancements
- **iOS app activation plan** — Production gateway wiring, macOS target addition
- **GitHub forks created** — `One-Mind-OS/swarmclaw` (fork of swarmclawai/swarmclaw), `One-Mind-OS/chowder-iOS` (fork of newmaterialco/chowder-iOS), `One-Mind-OS/openclaw-office` (fork of WW-AI-Lab/openclaw-office)
- **SwarmClaw customizations pushed** — OneMind-specific files (openclaw provider, logs component, provider-sets) pushed to fork

### Changed
- **Version bump v0.4 → v0.5** across all documentation (README, AGENTS.md, CLAUDE.md, ACCESS.md, codex/CLAUDE.md)
- **VPS cleanup** — Removed Star Office UI, Mission Control (5 Docker containers), Paperclip, OpenClaw Studio services
- **Caddy cleaned** — 4 remaining domains (gateway, swarm, office, nats)
- **Installed missing binaries** — `blu` (blucli v0.1.3) and `camsnap` (v0.2.0) on VPS
- **DEVELOPMENT.md updated** — Chowder-iOS replaces old openclaw-ios plan; Studio marked as removed

### Removed
- **command-center/openclaw-studio** — 350 tracked files removed (redundant to SwarmClaw)
- **command-center/swarmclaw** — 710 tracked files removed (moved to `One-Mind-OS/swarmclaw` fork)
- **command-center/ios-app** — 391 tracked files removed (replaced by Chowder-iOS fork)
- **command-center/design-reference** — Moved to SwarmClaw fork
- **command-center/Star-Office-UI** — Community clone deleted
- **command-center/openclaw-mission-control** — Community clone deleted
- **command-center/paperclip** — Community clone deleted
- Star Office UI, Mission Control, Paperclip, Studio services from VPS

## [0.4.0] - 2026-03-14

### Added
- **BLUEPRINT.md** — Full activation plan (16 sections, every service, every port, every command)
- **Custom domains** — gateway.onemindos.dev, swarm.onemindos.dev, office.onemindos.dev, nats.onemindos.dev via Caddy + Cloudflare
- **Skill Forge** (`skill-forge-1.0.0`) — Autonomous skill builder (4 tools: create, from_openapi, validate, list)
- **NATS Manager** (`nats-manager-1.0.0`) — Nervous system control skill (8 tools)
- **ACCESS.md** — Centralized credentials and access reference

## [0.3.0] - 2026-03-09

### Added
- **9-Agent Architecture** — Expanded from 4 agents (legacy, hp, le, ge) to 9 domain-scoped agents:
  - ⚡ Legacy (gpt-5.2) — Sovereign Orchestrator (default)
  - ⚔️ Spartan (gpt-4o) — Body & Physical Mastery (HP 27-28)
  - 🔮 Oracle (gpt-4o) — Mind, Money & Meaning (HP 25-26, 29-34) **NEW**
  - 🌳 Eden (gpt-4o) — Food Growing & Agriculture (LE farming)
  - 👑 Empire (gpt-5.2) — Business & Wealth (GE 75-84)
  - 🌿 Haven (gpt-4o) — Home, Family & Estate (LE 50-58) **NEW**
  - 🛡️ Guardian (gpt-4o) — Security & Legal (UI 11 + HP 31 + LE 57) **NEW**
  - ⚙️ Grid (gpt-4o) — Infrastructure Ops (UI 00-12) **NEW**
  - 🔺 Trinity (gpt-4o) — Community & Public (GE 75 + Community)
- **Oracle workspace** — 7 files (IDENTITY.md, SOUL.md, AGENTS.md, HEARTBEAT.md, BOOT.md, TOOLS.md, USER.md)
- **Discord server reorganization** — 5 categories aligned with Codex 4-pillar structure:
  - 🧭 Command Center: legacy-chat, agent-heartbeat, morning-brief, notifications, legacy-voice
  - 🧠 Sovereign Operations: grid-chat, guardian-chat
  - 💪 Holistic Performance: spartan-chat, oracle-chat
  - 🏠 Legacy Evolution: eden-chat, haven-chat
  - 🚀 Generational Empire: empire-chat, trinity-chat
- **Per-agent Discord bindings** — Each agent routes to its own Discord channel
- **Agent-to-agent communication** — All 9 agents can dispatch tasks to each other
- **Session dmScope** — `per-channel-peer` for proper session isolation
- **VPS config migration script** — `scripts/migrate-config.py`

### Changed
- **Spartan narrowed** — From HP 25-34 (all) to HP 27-28 (body/mastery only)
- **Eden redesigned** — From family steward to food growing/agriculture specialist
- **Haven redesigned** — From homestead/property to full home, family & estate steward
- **All agent workspace AGENTS.md files** — Updated with 9-agent system diagram and cross-domain routing
- **AGENT-ROSTER.md** — Version 3.0, all 9 agents documented with codex domains, skills, crons, escalation rules
- **openclaw.json.example** — Full 9-agent config template with bindings, guild channels, hooks
- **Root AGENTS.md + CLAUDE.md** — Updated repo structure, VPS paths, agent roster table

### Technical Details
- Gateway bind: `loopback` with Tailscale serve mode
- Memory search: hybrid (vector 0.7 + text 0.3) with MMR and temporal decay
- Heartbeat interval: 1h per agent to dedicated Discord channels
- Subagent pool: max 8 concurrent, gpt-4o-mini model
- All workspaces synced via onemind-pull.timer (GitHub → VPS every 1 min)

### Added
- **Tactical Solarpunk Theme** - New theme option blending tactical precision with organic solarpunk aesthetics
  - 50+ color palette with warm amber (#FFB86C) and moss green (#50FA7B) accents
  - Organic border radius system (10-16px) for softer, more natural UI elements
  - Breathing pulse animations for game nodes, connections, and status badges
  - Bezier curve connections in game visualization for organic flow
  - Warm amber particle system with organic movement patterns
  - `StatusBadge` component for reusable status indicators across the app
  - Theme persistence via SharedPreferences (survives app restarts)
  - Full WCAG AA accessibility compliance (4.5:1 contrast ratios)
  - Smooth theme transitions with 200ms animations

### Changed
- Extended `AppThemeMode` enum to support three themes (light/dark/solarpunk)
- Updated `TacticalColors` system with 3-way switch pattern for theme-aware colors
- Updated `TacticalTheme` with dedicated Solarpunk theme configuration
- Migrated 79 files from deprecated Flutter APIs (`.withOpacity` → `.withValues`)
- Updated sidebar navigation to use `TacticalColors` (removed hardcoded color values)
- Enhanced `ConnectionComponent` with support for Bezier curves
- Updated `ParticleSystem` with warm color palette for Solarpunk theme
- Modified `NodeComponent` with pulse animations and theme-aware rendering
- Updated all major UI components (cards, buttons, inputs, lists) for theme support

### Technical Details

#### Phase 1: Core Infrastructure
- Theme enumeration and mode switching
- Color system expansion (50+ named colors)
- Theme persistence layer with SharedPreferences
- Base theme configuration in `TacticalTheme`

#### Phase 2: Universal Components
- Card components (`TacticalCard`, info cards)
- Button components (primary, secondary, icon buttons)
- Input components (text fields, dropdowns)
- List components (`BaseListScreen`, list tiles)
- Navigation (sidebar, app shell)
- 79-file deprecated API migration

#### Phase 3: Specialized Components
- Game visualization (`NodeComponent`, `ConnectionComponent`)
- Particle system with organic movement
- Status badge system
- Real-time indicators
- Pulse animations and Bezier curves

### Files Changed
- **Frontend Configuration**: 3 files (theme provider, constants, theme config)
- **Universal Components**: 75+ files (screens, widgets, services)
- **Game Components**: 4 files (nodes, connections, particles, game screen)
- **Total**: 79+ files updated

---

## [0.2.0] - 2026-02-15

### Added
- Internal task system for task management
- Refactored frontend navigation hubs
- Implementation quick start guide

### Changed
- Repository cleanup and documentation organization
- Removed GitHub Actions workflows
- Updated CLAUDE.md with Agno + OM1 integration focus

---

## [0.1.0] - Initial Release

### Added
- FastAPI backend with Agno AgentOS integration
- Flutter web frontend with 50+ screens
- Multi-agent orchestration system
- PostgreSQL database with pgvector support
- NATS JetStream event bus
- WebSocket real-time updates
- Agent and team management
- Workflow builder
- Knowledge management
- Analytics and monitoring
- System topology visualization
