# STATUS.md — One Mind OS Live Status & Master Task List

> **Single source of truth** for what's running, what's broken, and what's planned.
> **Last updated**: March 17, 2026 (Phase 2 push — model providers, ambient awareness, VPS health)
> **Sources**: Live VPS audit + BLUEPRINT.md + DEVELOPMENT.md + PRODUCTION-PLAN.md + GITHUB-RESEARCH.md

---

## 1. Live System Inventory (What's Actually Running)

### VPS — `104.248.4.212` (Ubuntu 24.04, 16GB RAM, 193GB disk — 16% used)

| Service | Type | Port | Status | Version / Notes |
|---------|------|------|--------|---------|
| OpenClaw Gateway (9 agents) | systemd | 18789 | **LIVE** | v2026.3.13, 4 providers, 22 models |
| legacy-pulse (NATS→AMBIENT_STATE) | systemd | — | **LIVE** | 9 NATS subscriptions, writes AMBIENT_STATE.md |
| SwarmClaw | systemd | 3456 | **LIVE** | v0.7.7 (fork v1.1.6 build fails — needs dev work) |
| OpenClaw Office "The Homestead" | systemd | 5180 | **LIVE** | vanilla fork |
| Caddy (reverse proxy) | systemd | 443 | **LIVE** | 4 domains configured |
| onemind-watch (auto-commit) | systemd | — | **LIVE** | — |
| onemind-pull (git pull) | timer (1min) | — | **LIVE** | — |
| VPS health publisher | cron (*/5m) | — | **LIVE** | CPU/MEM/DISK → NATS `system.health.vps` |
| OpenViking Context DB | systemd | 1933 | **LIVE** | v0.2.6, local backend, OpenAI embeddings |
| NATS JetStream | Docker | 4222 | **LIVE** | 6 custom streams + 5 MQTT |
| PostgreSQL/pgvector | Docker | 5432 | **LIVE** | — |
| legacy-bridge (NATS→OpenClaw) | Docker | — | **LIVE** | — |
| nats-ingester | Docker | — | **LIVE** | — |
| NATS dashboard | Docker | 8222/18300 | **LIVE** | — |

### Custom Domains (Caddy → Tailscale IP 100.77.52.93)

| Domain | Target | Status |
|--------|--------|--------|
| `gateway.onemindos.dev` | `:18789` | **LIVE** |
| `swarm.onemindos.dev` | `:3456` | **LIVE** |
| `office.onemindos.dev` | `:5180` | **LIVE** |
| `nats.onemindos.dev` | `:8222` | **LIVE** |

### External Devices

| Device | Tailscale IP | Role | Status |
|--------|-------------|------|--------|
| MacBook Pro 16" M2 Max | `100.114.142.55` | Dev, BlueBubbles, local vision | **CONNECTED** |
| iPhone | `100.122.203.68` | iOS client | **CONNECTED** |
| VPS onemind-sentinel | `100.67.223.5` | Secondary/standby | Idle |

### Communication Channels

| Channel | Status | Notes |
|---------|--------|-------|
| Discord (9 channels + voice) | **LIVE** | All 9 agents bound to channels, voice enabled (Azure OpenAI TTS, alloy) |
| iMessage (BlueBubbles) | **LIVE** | Connected via Tailscale to MacBook at `100.114.142.55:1234` |
| WhatsApp | **NOT CONFIGURED** | Plugin available (`whatsapp`), not enabled |
| Telegram | **NOT CONFIGURED** | Plugin available (`telegram`), not enabled |

---

## 2. OpenClaw Configuration Audit

### What's Configured and Working

| Feature | Config | Status |
|---------|--------|--------|
| 9 domain agents | `agents.list` | **WORKING** — legacy, spartan, oracle, eden, empire, haven, guardian, grid, trinity |
| Default model | `azure-openai/gpt-4o` | **SET** |
| Legacy model | `azure-openai/gpt-5.2` | **SET** |
| Empire model | `azure-openai/gpt-5.2` | **SET** |
| Model providers | 4 active (22 models total) | **SET** — azure-openai (9), anthropic (3), openai-direct (7), gradient/DO (3) |
| Model aliases | 14 aliases | **SET** — gpt52, gpt5mini, gpt5nano, gpt4o, gpt4omini, o1, o3mini, sonnet4, sonnet, haiku, gpt41, gpt41mini, o4mini, llama70b |
| Session reset | `{mode: idle, idleMinutes: 10080}` | **SET** — 7-day window for Lossless Claw |
| Agent-to-agent | `agentToAgent.enabled: true` | **WORKING** — all 9 agents allowed |
| Subagents | `maxConcurrent: 8, maxSpawnDepth: 3` | **CONFIGURED** — gpt-4o-mini |
| Memory search | Hybrid BM25+vector | **WORKING** — indexes codex, memory, sessions |
| Compaction | `safeguard` mode | **SET** — 20k reserve floor, memory flush enabled |
| Discord | 9 channel bindings + voice | **WORKING** |
| BlueBubbles (iMessage) | Tailscale to MacBook | **WORKING** |
| Hooks | enabled with gmail preset + 3 mappings | **WORKING** |
| Cron | enabled, max 3 concurrent | **ENABLED** but no jobs found |
| TTS | Azure OpenAI tts-1, alloy (via Azure credits) | **SET** |
| Lobster | plugin enabled | **WORKING** — 9 routines in workspace |
| Skills auto-reload | `watch: true` | **WORKING** |
| Thread bindings | `spawnSubagentSessions: true` | **WORKING** |

### Model Provider Details

| Provider | Models | Status |
|----------|--------|--------|
| **azure-openai** | gpt-5.2, gpt-5.2-codex, gpt-5.1-codex-mini, gpt-5-mini, gpt-5-nano, gpt-4o, gpt-4o-mini, o1, o3-mini | **ACTIVE** — primary provider |
| **anthropic** | claude-sonnet-4, claude-3.5-sonnet, claude-3.5-haiku | **ACTIVE** |
| **openai** (direct) | gpt-4o, gpt-4o-mini, o3-mini, gpt-4.1, gpt-4.1-mini, gpt-4.1-nano, o4-mini | **ACTIVE** |
| **gradient** (DigitalOcean) | llama-3.3-70b-instruct, mistral-small-latest, qwen2.5-72b-instruct | **ACTIVE** — `GRADIENT_BASE_URL=https://inference.do-ai.run/v1` |
| **groq** | — | **READY** — env.tpl entry commented, needs 1Password item |
| **google** (Gemini) | — | **READY** — env.tpl entry commented, needs 1Password item |
| **bedrock** (AWS) | — | **READY** — env.tpl entry commented, needs 1Password item + AWS creds |

### What's Working Now (Fixed in Phase 1)

| Feature | Status |
|---------|--------|
| **Heartbeat** | ✅ Firing every 30m, target=last, activeHours 07:00-23:00 ET |
| **Cron jobs** | ✅ 3 enabled: morning-briefing (7am), system-health-check (6h), evening-wrap (9pm) |
| **NATS JetStream streams** | ✅ 6 custom: EVENTS (30d), SYSTEM (7d), NEXUS (30d), DIGITAL (30d), PHYSICAL (30d), COMMS (14d) |
| **Lossless Claw plugin** | ✅ v0.3.0, DAG-based lossless summarization, threshold 0.75 |

### What's NOT Configured (Should Be)

| Feature | What's Missing |
|---------|---------------|
| **legacy-eyes** | Not deployed — SmolVLM on MacBook Pro, needs setup |
| **Groq/Gemini/Bedrock providers** | env.tpl ready, needs 1Password items + keys |
| **WhatsApp/Telegram channels** | Plugins available but not enabled |
| **ClawVault** | Not evaluated — structured memory with entity graph |

---

## 3. Workspace File Problems (Stale References)

**ALL RESOLVED** — ClickUp → Todoist, Agno → agent-to-agent, WhatsApp → Discord, Agno dispatch removed, AMBIENT_STATE awareness added.

| File | Status |
|------|--------|
| **HEARTBEAT.md** | ✅ Clean — Todoist, NATS, AMBIENT_STATE, no ClickUp/Agno/WhatsApp |
| **BOOT.md** | ✅ Clean — NATS bridge, onemind-watch, AMBIENT_STATE, no Agno |
| **AGENTS.md** | ✅ Clean — agent-to-agent routing, no Agno dispatch |
| **TOOLS.md** | ✅ Clean — nexus-bridge updated for NATS events, agno-dispatch removed |
| **USER.md** | ✅ Clean — Todoist, Discord primary, correct GitHub org |

---

## 4. Skills Audit (27 active + 11 archived + 5 kepano)

### Active Skills (27)

| Skill | Purpose | Status |
|-------|---------|--------|
| `nexus-bridge` | NATS event publishing + robot commands | Working |
| `nats-manager-1.0.0` | NATS stream/subject management | Working |
| `home-assistant-1.0.0` | Smart home control (6 tools) | Working |
| `todoist` | Task management (CLI-based) | Working |
| `github-api-1.0.3` | GitHub operations | Working |
| `onepassword-1.0.0` | 1Password credential queries | Working |
| `obsidian-1.0.0` | Vault integration (legacy) | Working |
| `rss-digest-1.0.0` | Newsletter digestion | Working |
| `recipe-scraper-1.0.0` | Meal planning | Working |
| `resend-1.0.0` | Email delivery | Working |
| `twilio-1.0.0` | SMS escalation | Working |
| `agent-browser-0.2.0` | Web automation | Working |
| `api-gateway-1.0.48` | 100+ API proxy (Maton) | Working |
| `proactive-agent-3.1.0` | Proactive behavior patterns | Working |
| `weather-1.0.0` | Daily weather | Working |
| `zoho-cli` | CRM + calendar | Working |
| `vision-1.0.0` | Vision API | Working |
| `hitl-escalate-1.0.0` | Human-in-the-loop escalation | Working |
| `skill-forge-1.0.0` | Autonomous skill builder | Working |
| `cloudflare-api-1.0.0` | DNS/domain management | Working |
| `spotify-1.0.0` | Music control | Working |
| `legacy-eyes` | Vision loop (not deployed yet) | Placeholder |
| `legacy-pulse` | NATS ambient state (not deployed yet) | Placeholder |
| `defuddle` | Web content extraction (kepano) | Knowledge skill |
| `json-canvas` | JSON Canvas format (kepano) | Knowledge skill |
| `obsidian-bases` | Obsidian Bases format (kepano) | Knowledge skill |
| `obsidian-cli` | Obsidian CLI usage (kepano) | Knowledge skill |
| `obsidian-markdown` | Obsidian Markdown syntax (kepano) | Knowledge skill |

### Archived (11 — in `skills-archive/`)

| Skill | Reason |
|-------|--------|
| `browser-1.0.0` | Superseded by `agent-browser-0.2.0` |
| `credit-repair-1.0.0` | Niche, not daily-driver |
| `desktop-control-1.0.0` | Windows-focused, not useful on Mac/VPS |
| `find-skills-0.1.0` | Dev tool, not runtime |
| `homestead-control` | No homestead yet |
| `lob-1.0.0` | Physical mail — only with credit-repair |
| `skill-creator-0.1.0` | Dev tool, not runtime |
| `todoist-api` | Duplicate of `todoist` CLI |
| `transit-1.0.0` | Rarely used |
| `uber-1.0.0` | OAuth not configured |
| `zoho-mail-1.0.0` | Duplicate of `zoho-cli` |

### Config Cleanup ✅

| Item | Status |
|------|--------|
| `new-agent` in `agents.list` | ✅ Removed from openclaw.json |
| `studio.onemindos.dev` in `allowedOrigins` | ✅ Removed from openclaw.json |

---

## 5. GitHub Forks Status

| Repo | Upstream | Stars | Status | Next Action |
|------|----------|-------|--------|-------------|
| `One-Mind-OS/swarmclaw` | `swarmclawai/swarmclaw` | 140 | **Forked** + customizations pushed | Upgrade v0.7.7 → v1.1.5, build 5 new pages |
| `One-Mind-OS/chowder-iOS` | `newmaterialco/chowder-iOS` | 541 | **Forked** (vanilla) | Wire to gateway, sideload, port Voice/Watch/Camera |
| `One-Mind-OS/openclaw-office` | `WW-AI-Lab/openclaw-office` | 302 | **Forked** (vanilla) | Theme "The Homestead", 9 agent avatars, NATS events |
| `One-Mind-OS/VisionClaw` | `Intent-Lab/VisionClaw` | 1,600 | **Forked** (vanilla) | Cherry-pick NixClaw features, wire to gateway, Month 2 |

---

## 6. Master Task List — All Planned Work

Everything from BLUEPRINT, DEVELOPMENT, PRODUCTION-PLAN, and GITHUB-RESEARCH consolidated into one prioritized list.

### Phase 1: Fix What's Broken + Foundation (This Week)

| # | Task | Source | Priority | Status |
|---|------|--------|----------|--------|
| 1.1 | **Fix heartbeat** — added `target: "last"`, prompt, activeHours 07:00-23:00 EST, directPolicy allow, lightContext true | PRODUCTION-PLAN | P0 | **DONE** |
| 1.2 | **Update HEARTBEAT.md** — removed ClickUp, Agno, WhatsApp; added Todoist, NATS, AMBIENT_STATE | PRODUCTION-PLAN | P0 | **DONE** |
| 1.3 | **Update BOOT.md** — removed Agno health check; added NATS bridge, onemind-watch, AMBIENT_STATE | PRODUCTION-PLAN | P0 | **DONE** |
| 1.4 | **Update AGENTS.md** — removed Agno dispatch; documented agent-to-agent as primary | PRODUCTION-PLAN | P0 | **DONE** |
| 1.5 | **Install Lossless Claw** — `@martian-engineering/lossless-claw@0.3.0` installed + pinned, `[lcm] Plugin loaded` confirmed | GITHUB-RESEARCH P0 | P0 | **DONE** |
| 1.6 | **Create 6 NATS JetStream streams** — EVENTS (30d), SYSTEM (7d), NEXUS (30d), DIGITAL (30d), PHYSICAL (30d), COMMS (14d) — all verified | BLUEPRINT | P0 | **DONE** |
| 1.7 | **Archive 11 unused skills** — moved to `skills-archive/` (both local + VPS) | PRODUCTION-PLAN | P1 | **DONE** |
| 1.8 | **Install Obsidian skills** (kepano) — 5 skills installed: `obsidian-markdown`, `obsidian-bases`, `json-canvas`, `obsidian-cli`, `defuddle` (36 total skills) | GITHUB-RESEARCH P1 | P1 | **DONE** |
| 1.9 | **Remove `new-agent`** from openclaw.json agents list | Audit | P1 | **DONE** |
| 1.10 | **Clean up allowedOrigins** — removed `studio.onemindos.dev` | Audit | P2 | **DONE** |
| 1.11 | **Set up cron jobs** — morning brief (7am), health check (every 6h), evening digest (9pm) | PRODUCTION-PLAN | P1 | **DONE** |
| 1.12 | **Push monorepo restructure** — git push the restructure commit to GitHub | DEVELOPMENT | P0 | **DONE** |
| 1.13 | **Add multi-provider models** — 4 providers active: Azure OpenAI (9), Anthropic (3), OpenAI Direct (7), Gradient/DO (3) = 22 models, 14 aliases. Session reset for Lossless Claw. Groq/Gemini/Bedrock stubbed in env.tpl. | PRODUCTION-PLAN | P2 | **DONE** |

### Phase 2: Ambient Awareness (Week 2)

| # | Task | Source | Priority | Status |
|---|------|--------|----------|--------|
| 2.1 | **Deploy legacy-pulse** — systemd service, nats-py v2.14.0, 9 NATS subscriptions, writes to AMBIENT_STATE.md, running as `onemind` | BLUEPRINT | P0 | **DONE** |
| 2.2 | **Create AMBIENT_STATE.md template** — live and populating with VPS health events (CPU, MEM, DISK), vision digest placeholder, recent NATS events (last 20) | BLUEPRINT | P0 | **DONE** |
| 2.3 | **Deploy OpenViking Context DB** — v0.2.6 deployed as systemd service on port 1933, local backend with OpenAI text-embedding-3-small, codex indexing in progress | GITHUB-RESEARCH P0 | P1 | **DONE** |
| 2.4 | **VPS health → NATS** — `nats-health-pub.sh` cron every 5 min, publishes CPU/MEM/DISK/load/uptime/Docker to `system.health.vps` | BLUEPRINT | P1 | **DONE** |
| 2.5 | **Deploy legacy-eyes** — SmolVLM-256M on MacBook Pro, vision digest → NATS → AMBIENT_STATE.md | BLUEPRINT | P2 | NOT DONE |
| 2.6 | **Update HEARTBEAT.md** — Step 0: read AMBIENT_STATE.md FIRST, VPS thresholds (disk>85%, cpu>90%), Step 5: nexus_publish, Discord delivery | BLUEPRINT | P0 | **DONE** |
| 2.7 | **Evaluate ClawVault** — structured memory with Obsidian integration, entity graph | GITHUB-RESEARCH P2 | P2 | NOT DONE |

### Phase 3: SwarmClaw + Mobile (Week 3-4)

| # | Task | Source | Priority | Status |
|---|------|--------|----------|--------|
| 3.1 | **Upgrade SwarmClaw** — v0.7.7 → v1.1.6. Fork cloned to VPS, npm installed (1104 deps), but `next build` **FAILS** (`rpcOnConnectedGateway` export missing from `openclaw.ts`). Existing v0.7.7 build serves fine. Needs dedicated dev session to fix. | GITHUB-RESEARCH | P0 | **BLOCKED** |
| 3.2 | **Build `/system` page** — VPS monitor (CPU, RAM, disk, Docker, firewall) | DEVELOPMENT | P1 | NOT DONE |
| 3.3 | **Build `/soul` page** — Identity editor with version history | DEVELOPMENT | P1 | NOT DONE |
| 3.4 | **Build `/files` page** — Workspace file browser with live preview | DEVELOPMENT | P1 | NOT DONE |
| 3.5 | **Build `/nats` page** — JetStream monitor with live event tailing | DEVELOPMENT | P1 | NOT DONE |
| 3.6 | **Build `/terminal` page** — Safe web terminal with command allowlist | DEVELOPMENT | P2 | NOT DONE |
| 3.7 | **Wire Chowder-iOS to gateway** — `wss://gateway.onemindos.dev` + token | DEVELOPMENT | P1 | NOT DONE |
| 3.8 | **Sideload Chowder to iPhone** — build from Xcode, configure connection | DEVELOPMENT | P1 | NOT DONE |
| 3.9 | **Add global STT mic** to SwarmClaw `app-layout.tsx` | DEVELOPMENT | P2 | NOT DONE |
| 3.10 | **Add sub-agent cost breakdown** to `/usage` view | DEVELOPMENT | P2 | NOT DONE |

### Phase 4: Desktop + Office (Week 4)

| # | Task | Source | Priority | Status |
|---|------|--------|----------|--------|
| 4.1 | **Customize Office fork** — solarpunk theme, rename to "The Homestead" | DEVELOPMENT | P1 | NOT DONE |
| 4.2 | **Add 9 agent avatars** to Office fork | DEVELOPMENT | P2 | NOT DONE |
| 4.3 | **Add NATS event feed** to Office fork | DEVELOPMENT | P2 | NOT DONE |
| 4.4 | **Evaluate ClawX fork** — Electron desktop, add MediaPipe + ambient overlay | GITHUB-RESEARCH P1 | P2 | NOT DONE |

### Phase 5: Voice + Vision + GPU (Month 2)

| # | Task | Source | Priority | Status |
|---|------|--------|----------|--------|
| 5.1 | **Self-hosted Chatterbox TTS** — replace OpenAI TTS | BLUEPRINT | P2 | NOT DONE |
| 5.2 | **On-demand GPU Droplet** — RTX 4000 Ada, 2hr/day for local inference | BLUEPRINT | P2 | NOT DONE |
| 5.3 | **VisionClaw customization** — wire to gateway, cherry-pick NixClaw features (background audio, Live Activity, setup wizard) | BLUEPRINT | P2 | NOT DONE |
| 5.4 | **Hume AI emotion SDK** — emotional voice intelligence | GITHUB-RESEARCH | P3 | NOT DONE |
| 5.5 | **Port Voice module to Chowder** — TalkModeManager, VoiceWakeManager | DEVELOPMENT | P2 | NOT DONE |
| 5.6 | **Port Watch module to Chowder** — watchOS companion app | DEVELOPMENT | P3 | NOT DONE |

### Phase 6: Long-Term (Month 3+)

| # | Task | Source | Priority | Status |
|---|------|--------|----------|--------|
| 6.1 | **Mac Studio M4 Ultra** — run 70B+ models locally 24/7 | BLUEPRINT | P3 | NOT DONE |
| 6.2 | **WhatsApp channel** — enable whatsapp plugin in openclaw.json | PRODUCTION-PLAN | P3 | NOT DONE |
| 6.3 | **Telegram channel** — enable telegram plugin | PRODUCTION-PLAN | P3 | NOT DONE |
| 6.4 | **LiveKit persistent room** — true ambient multi-device voice | PRODUCTION-PLAN | P3 | NOT DONE |
| 6.5 | **Homestead IoT** — bring back homestead-control + OM1 when property acquired | PRODUCTION-PLAN | P3 | NOT DONE |

---

## 7. GitHub Research — Adoption Tracker

Items from GITHUB-RESEARCH.md and their status:

| Priority | Repo | Verdict | Status |
|----------|------|---------|--------|
| P0 | **Lossless Claw** (1.9k stars) | Install as plugin — lossless memory | **INSTALLED** ✅ v0.3.0, threshold 0.75 |
| P0 | **OpenViking** (12.7k stars) | Deploy as context DB — +46% accuracy | **DEPLOYED** ✅ v0.2.6, port 1933, local backend |
| P1 | **Obsidian skills** (14.1k stars) | Install kepano's 5 skills | **INSTALLED** ✅ 5 skills in workspace |
| P1 | **ClawX** (4.5k stars) | Fork for desktop command center | **NOT STARTED** → Task 4.4 |
| P2 | **ClawVault** (584 stars) | Evaluate after P0 memory items | **NOT STARTED** → Task 2.7 |
| Ref | **awesome-openclaw-skills** (37.9k) | Skill discovery catalog | Bookmarked |
| Ref | **openclaw-api-list** (126) | API discovery catalog | Bookmarked |
| Done | **chowder-iOS** (541 stars) | Forked as primary iOS client | **FORKED** → `One-Mind-OS/chowder-iOS` |
| Done | **SwarmClaw** (140 stars) | Forked as command center | **FORKED** → `One-Mind-OS/swarmclaw` |
| Done | **OpenClaw Office** (302 stars) | Forked as "The Homestead" | **FORKED** → `One-Mind-OS/openclaw-office` |
| Done | **VisionClaw** (1.6k stars) | Forked for Gemini Live voice+vision | **FORKED** → `One-Mind-OS/VisionClaw` |
| Future | **Hume AI SDK** (165 stars) | Emotion AI for voice | Phase 5 |
| Skip | **ClawdAgent** (9 stars) | Inspiration only, too niche | Skipped |
| Skip | **ada_local** (128 stars) | Different architecture | Skipped |
| Skip | **Mission Control** (2.6k) | Removed, SwarmClaw replaces | Removed from VPS |
| Skip | **NixClaw** (3 stars) | VisionClaw fork — cherry-pick features into our VisionClaw fork instead | Skipped |

---

## 8. Architecture Reference

```
Zeus (Discord · iMessage · Voice · iOS · Browser · Webcam · Ray-Ban)
    ↓
OpenClaw Gateway (VPS :18789, 9 agents, 4 providers, 22 models)
    ├── Legacy ⚡     gpt-5.2    full tools, 27 skills, 9 lobsters
    ├── Spartan ⚔️    gpt-4o     messaging + allow-list
    ├── Oracle 🔮     gpt-4o     messaging + allow-list
    ├── Eden 🌳       gpt-4o     messaging + allow-list
    ├── Empire 👑     gpt-5.2    messaging + allow-list
    ├── Haven 🌿      gpt-4o     messaging + allow-list
    ├── Guardian 🛡️    gpt-4o     messaging + allow-list
    ├── Grid ⚙️       gpt-4o     messaging + allow-list
    └── Trinity 🔺    gpt-4o     messaging + allow-list

NATS JetStream (:4222)
    ├── legacy-bridge (NATS → OpenClaw hooks)
    ├── nats-ingester (event processing)
    ├── legacy-pulse (NATS → AMBIENT_STATE.md) ✅ RUNNING
    ├── VPS health publisher (cron */5m → system.health.vps) ✅ RUNNING
    └── 6 custom streams: EVENTS, SYSTEM, NEXUS, DIGITAL, PHYSICAL, COMMS ✅

PostgreSQL/pgvector (:5432)

Dashboards:
    ├── gateway.onemindos.dev → OpenClaw Gateway UI
    ├── swarm.onemindos.dev   → SwarmClaw v0.7.7 (fork v1.1.6 build blocked)
    ├── office.onemindos.dev  → OpenClaw Office "The Homestead"
    └── nats.onemindos.dev    → NATS monitoring

Plugins:
    ├── Lossless Claw v0.3.0 (DAG-based lossless summarization)
    └── Session reset: idle 10080min (7-day window for Lossless Claw)

Model Providers:
    ├── azure-openai  → gpt-5.2, gpt-5.2-codex, gpt-5.1-codex-mini, gpt-5-mini, gpt-5-nano, gpt-4o, gpt-4o-mini, o1, o3-mini
    ├── anthropic      → claude-sonnet-4, claude-3.5-sonnet, claude-3.5-haiku
    ├── openai-direct  → gpt-4o, gpt-4o-mini, o3-mini, gpt-4.1, gpt-4.1-mini, gpt-4.1-nano, o4-mini
    └── gradient (DO)  → llama-3.3-70b, mistral-small, qwen2.5-72b

NOT RUNNING (planned):
    ├── legacy-eyes (SmolVLM → vision digest) — needs MacBook setup
    └── Groq / Google Gemini / AWS Bedrock (env.tpl ready, needs 1Password items)
```

---

## 9. Quick Reference

### Key Connections

| What | Value |
|------|-------|
| VPS IP | `104.248.4.212` |
| Tailscale IP | `100.77.52.93` |
| Gateway URL | `https://gateway.onemindos.dev` |
| Gateway Token | 1Password: `00-24 SO / OpenClaw Gateway / gateway-token` |
| OpenClaw version | v2026.3.13 |
| NATS port | 4222 |
| PostgreSQL port | 5432 |

### All GitHub Repos

| Repo | Purpose |
|------|---------|
| `One-Mind-OS/onemindos-v0.5` | Lean monorepo (brain + docs + infra) |
| `One-Mind-OS/swarmclaw` | Unified command center (Next.js) |
| `One-Mind-OS/chowder-iOS` | Primary iOS/macOS client (Swift) |
| `One-Mind-OS/openclaw-office` | "The Homestead" 3D visualization (Vite + R3F) |
| `One-Mind-OS/VisionClaw` | Gemini Live voice + vision for smart glasses (iOS + Android) |

### Cost (Current Monthly)

| Item | Cost |
|------|------|
| VPS | ~$50-100 |
| Azure OpenAI | ~$150-200 |
| OpenAI Direct (TTS, embeddings, gpt-4.1) | ~$20-50 |
| Anthropic (Claude Sonnet 4) | ~$10-30 |
| DigitalOcean Gradient (Llama, Qwen) | ~$5-20 |
| **Total** | **~$235-400/mo** |

### Cron Jobs (onemind user)

| Schedule | Script | Purpose |
|----------|--------|---------|
| `*/15 * * * *` | `workspace-backup.sh` | Workspace backup |
| `0 */6 * * *` | `backup-swarmclaw.sh` | SwarmClaw backup |
| `*/5 * * * *` | `nats-health-pub.sh` | VPS health → NATS |


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
