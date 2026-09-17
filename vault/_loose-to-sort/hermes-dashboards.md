---
title: "Hermes Dashboards & Agent UIs"
type: reference
status: active
tags: [platform, ai, hermes, dashboard, web-ui, agent-management]
updated: 2026-06-04
---

# Hermes Dashboards & Agent UIs

> All web-based and desktop interfaces for managing Hermes Agent and the One Mind agent fleet. Each serves a distinct purpose — use based on task.

---

## Desktop Apps (3)

| App | Bundle ID | Version | Type | Notes |
|-----|-----------|---------|------|-------|
| **HermesDesktop.app** | com.dodo-reach.hermesdesktop | latest | SwiftUI native | Primary daily driver — SSH-direct, no backend |
| **Hermes.app** | com.nousresearch.hermes | 0.15.1 | Electron | Official NousResearch desktop wrapper |
| **Hermes Agent.app** | com.nousresearch.hermes | 0.5.5 | Electron | Older NousResearch electron build |

All three are in `/Applications/`. All kept for evaluation.

---

## Web UIs (4)

### 1. Hermes WebUI — `localhost:8787`
- **Repo:** `github.com/nesquena/hermes-webui`
- **Install:** `~/dev/hermes-webui/` — fresh clone 2026-06-04
- **Start:** `cd ~/dev/hermes-webui && bash start.sh --no-browser`
- **State:** `~/.hermes/webui/`
- **Config:** `~/dev/hermes-webui/.env`
- **License:** MIT
- **What it is:** Lightweight vanilla Python + JS, three-panel layout (sessions / chat / file browser). No build step. Full CLI parity in the browser.
- **Key features:** Session management, profile switching, model selector, workspace file browser, token ring, Hermes Control Center launcher
- **Bot name set to:** Legacy
- **Best for:** Quick chat sessions, lightweight access, minimal overhead

```bash
# Start
cd ~/dev/hermes-webui && bash start.sh --no-browser

# Stop
kill $(lsof -ti:8787)
```

---

### 2. Hermes Workspace — `localhost:3040`
- **Repo:** `github.com/outsourc-e/hermes-workspace`
- **Install:** `~/dev/hermes-workspace/`
- **Start:** `cd ~/dev/hermes-workspace && PORT=3040 pnpm start`
- **Config:** `~/dev/hermes-workspace/.env`
- **License:** MIT
- **Stars:** 5,319
- **What it is:** Full-featured Node.js + Vite workspace UI. The most complete Hermes UI available.
- **Key features:**
  - **Conductor** — parallel sub-agent orchestration, swarm mission dispatch
  - **Swarm Mode** — tmux-backed workers with role dispatch
  - **Operations** — multi-agent panel with profile presets
  - **Agent View** — live agent panel
  - **2,000+ Skills marketplace**
  - **Monaco editor** file browser
  - **PWA** — installable on phone/tablet
  - **Tailscale support** — built-in native
  - **Cost ledger** — session cost tracking
  - **MCP page** — MCP server management
  - Local models: Ollama, LM Studio, vLLM, Atomic Chat
- **Gateway:** `HERMES_API_TOKEN=JdFfBxXEPLYy26flc6rdckkrvLp3DaEkKIQhebDiWTc`
- **Best for:** Primary daily workspace, multi-agent coordination, Conductor missions

```bash
# Start
cd ~/dev/hermes-workspace && PORT=3040 pnpm start

# Stop
kill $(lsof -ti:3040)
```

---

### 3. Hermes Studio — `localhost:8648`
- **Repo:** `github.com/EKKOLearnAI/hermes-web-ui`
- **Package:** `hermes-web-ui` (npm global)
- **Install:** `npm install -g hermes-web-ui`
- **Start:** `hermes-web-ui start` (default 8648)
- **State:** `~/.hermes-web-ui/`
- **License:** BSL-1.1 ⚠️ (not fully open source — converts after 4 years)
- **Stars:** 7,168
- **What it is:** Vue3 + Koa2 full-featured dashboard with the highest community adoption.
- **Key features:**
  - **Group Chat rooms** — multi-agent chat with @mention routing (unique feature — @Legacy, @Zenith, etc. in one room)
  - **8 Platform channels** — Telegram, Discord, Slack, WhatsApp, Matrix, Feishu, WeChat, WeCom config
  - **Multi-profile** — per-agent scoped config, sessions, usage, memory, skills
  - **Model management** — auto-discover from credential pool
  - **Usage analytics** — token breakdown, cost, cache hit rate, 30-day trends
  - **Web terminal** — node-pty + xterm
  - **File browser** — local, Docker, SSH, Singularity backends
  - **RBAC** — super admin + admin roles
- **Default login:** `admin` / `123456` — change immediately in settings
- **Best for:** Coordinating multiple agents simultaneously via @mention, platform channel config

```bash
# Start
hermes-web-ui start

# Stop
hermes-web-ui stop

# Status
hermes-web-ui status
```

---

### 4. Hermes Control Interface (HCI) — `localhost:10274`
- **Repo:** `github.com/xaspx/hermes-control-interface`
- **Install:** `~/dev/hermes-control-interface/`
- **Start:** `cd ~/dev/hermes-control-interface && PORT=10274 node server.js`
- **Config:** `~/dev/hermes-control-interface/hci.config.yaml`
- **License:** MIT
- **Stars:** 744
- **What it is:** Vanilla JS + Express single-server ops dashboard. 11 pages.
- **Key features:**
  - **Office v3 / ZOO Swarm Monitor** — 3-panel: health grid for all agents, 8-lane Kanban pipeline with dependency arrows, live event feed
  - **Per-agent sub-pages** — Legacy/Zenith/Unity/Eternal/Sovereignty/Sentinel each get their own dashboard, sessions, cron, config, memory, skills
  - **MCP control plane** — start/stop/restart MCP servers, live log tail, config editor
  - **Usage analytics** — token breakdown, cost projections, budget alerts, 30-day trends
  - **Skills** — browse, install, remove from browser
  - **RBAC** — 20 permissions across 3 roles (admin/viewer/custom)
  - **PWA** — offline caching
  - **Maintenance** — backup, restore, HCI update, system doctor
- **Login password:** `01bffe80c40c8578c36a61c1587ca682` (in hci.config.yaml)
- **Gateway key:** configured in hci.config.yaml
- **Best for:** Fleet ops monitoring, per-agent oversight, MCP management, cost alerts

```bash
# Start
cd ~/dev/hermes-control-interface && PORT=10274 node server.js

# Stop
kill $(lsof -ti:10274)
```

---

## Universal Dashboards (Single Pane of Glass)

### AionUi — Desktop App
- **Install:** `/Applications/AionUi.app`
- **Repo:** `github.com/iOfficeAI/AionUi`
- **License:** Apache-2.0
- **Stars:** 27,509
- **What it is:** Electron desktop app that auto-detects all installed CLI agents and wraps them in one unified interface.
- **Confirmed framework support:** OpenClaw ✓ Hermes ✓ Claude Code ✓ GitHub Copilot ✓ OpenCode ✓ Gemini CLI ✓ AWS Bedrock ✓ + 16 others
- **Note:** Sentry telemetry baked into official builds. Build from source (`~/dev/aionui-local`) to strip it.
- **Best for:** THE single pane of glass — interact with all frameworks from one place

### SwarmClaw — `localhost:3456`
- **Install:** `npm install -g @swarmclawai/swarmclaw` (already installed)
- **Start:** Already running as background process
- **License:** MIT
- **Stars:** 544 (early stage)
- **Framework support:** OpenClaw ✓ Hermes ✓ Claude Code ✓ Copilot ✓ OpenCode ✓ Gemini ✓
- **Best for:** Agent swarm orchestration, org chart visualization, OTLP metrics export

### Mission Control — `localhost:3030`
- **Repo:** `github.com/builderz-labs/mission-control`
- **Install:** `~/dev/mission-control/`
- **Start:** `cd ~/dev/mission-control && PORT=3030 pnpm start`
- **License:** MIT
- **Stars:** 5,161
- **What it is:** Next.js ops dashboard, best OpenClaw integration of any tool.
- **Framework support:** OpenClaw ✓ (deep native) — Hermes/Copilot/OpenCode need custom adapters
- **Best for:** OpenClaw fleet monitoring, heartbeat tracking, task dispatch, cost dashboards
- **First run:** Visit `localhost:3030/setup` to create admin account

---

## Gateway Configuration

All web UIs connect to the Hermes gateway on port `8642`.

```
Gateway URL:  http://127.0.0.1:8642
API Key:      API_SERVER_KEY in ~/.hermes/.env
              Value: JdFfBxXEPLYy26flc6rdckkrvLp3DaEkKIQhebDiWTc
```

| UI | Key Variable | Status |
|----|-------------|--------|
| Hermes Workspace | `HERMES_API_TOKEN` in `.env` | ✓ Configured |
| HCI | `gateway_api_key` in `hci.config.yaml` | ✓ Configured |
| Hermes WebUI | auto-discovers from `~/.hermes/config.yaml` | ✓ Auto |
| Hermes Studio | managed internally | ✓ Auto |

---

## All Ports Reference

| Port | Service | Process |
|------|---------|---------|
| 8642 | Hermes Gateway API | Hermes LaunchAgent (PID varies) |
| 8787 | Hermes WebUI | `~/dev/hermes-webui/start.sh` |
| 9119 | Hermes Dashboard (official) | `hermes dashboard` |
| 3040 | Hermes Workspace | `~/dev/hermes-workspace/` |
| 8648 | Hermes Studio | `hermes-web-ui start` |
| 10274 | Hermes Control Interface | `~/dev/hermes-control-interface/` |
| 3030 | Mission Control | `~/dev/mission-control/` |
| 3456 | SwarmClaw | npm global process |

---

## Start All Services

```bash
# Hermes WebUI (8787)
cd ~/dev/hermes-webui && bash start.sh --no-browser &

# Hermes Workspace (3040)
cd ~/dev/hermes-workspace && PORT=3040 pnpm start &

# Hermes Studio (8648)
hermes-web-ui start &

# Hermes Control Interface (10274)
cd ~/dev/hermes-control-interface && PORT=10274 node server.js &

# Mission Control (3030)
cd ~/dev/mission-control && PORT=3030 pnpm start &

# AionUi
open /Applications/AionUi.app
```

---

## Decision Framework — Which UI to Use

| Task | Use |
|------|-----|
| Quick chat with Legacy | Hermes WebUI (8787) or HermesDesktop.app |
| Multi-agent Conductor mission | Hermes Workspace (3040) |
| @mention multiple agents in one thread | Hermes Studio (8648) |
| Monitor fleet health / ZOO swarm view | HCI (10274) |
| MCP server management | HCI (10274) |
| All frameworks in one place | AionUi desktop |
| OpenClaw fleet ops | Mission Control (3030) |

---

## Notes
- **Phone access:** All web UIs are currently desktop-only (`127.0.0.1`). Tailscale reverse proxy (Caddy) planned — will expose all on Tailscale mesh for phone/remote access.
- **Dock shortcuts:** All 4 Hermes web UIs added to dock as `.webloc` files at `~/Applications/Hermes UIs/`
- **8787 instability history:** Previous install at `/Users/onemind/apps/hermes-webui` (unknown fork) was crashing with port conflicts. Replaced 2026-06-04 with fresh clone of `nesquena/hermes-webui` — stable.
