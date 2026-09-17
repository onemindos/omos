---
title: "AI Agent Dashboards & UIs — Research Report"
type: reference
status: active
created: 2026-04-26
owner: zeus
tags: [type/reference, domain/SO, domain/02, openclaw, dashboard, ui, research]
ai_summary: "Comprehensive research on open-source AI agent dashboards and UIs — OpenClaw-specific (Clawtrol, ClawFleet), general agent UIs (Open WebUI, LibreChat, LobeHub, Dify, AnythingLLM), and AG-UI protocol. Ranked by real-time capability, feature density, and OpenClaw compatibility."
---

# AI Agent Dashboards & UIs — Research Report

> **Context:** Research conducted April 24–26, 2026. Goal: find the most feature-rich, real-time, context-aware frontends for managing OpenClaw agents.

---

## 🔴 TIER 0 — PROTOCOL-LEVEL (Game-Changing Find)

---

### `ag-ui-protocol/ag-ui` ⭐ 13.2k | TypeScript + Python | MIT
**URL:** https://github.com/ag-ui-protocol/ag-ui
**Updated:** April 24, 2026 (active daily)

AG-UI is an open standard protocol for **connecting AI agents to user-facing applications in real-time** — the missing link between the OpenClaw gateway and any custom frontend including onemind-app.

**Protocol stack context:**
- **MCP** = gives agents tools
- **A2A** = agents talk to agents
- **AG-UI** = agents talk to *your UI*

**What it does:**
- ~16 standard event types: streaming, state sync, tool output, generative UI, HITL interrupts
- Works over SSE, WebSocket, webhooks — matches existing stack
- Bi-directional state sync between agent and frontend
- Human-in-the-loop collaboration built in
- Generative UI (agents render components directly in chat)

**Why it matters for OneMind:**
- **Has a Dart/Flutter SDK** — can be wired directly into `onemind-app`
- Integrates with LangGraph, CrewAI, Google ADK, AWS Bedrock, Mastra, PydanticAI, LlamaIndex, AG2
- `npx create-ag-ui-app my-agent-app` — working app in 60 seconds
- 101 contributors, 1.2k forks

**AG-UI Dojo (live demo):** https://dojo.ag-ui.com

---

## 🟠 TIER 1 — DIRECTLY OPENCLAW-SPECIFIC

---

### `nachoiacovino/clawtrol` ⭐ 32 | TypeScript/Next.js | MIT
**URL:** https://github.com/nachoiacovino/clawtrol
**Version:** v0.1.2 | **Stack:** Next.js 16, React 19, Tailwind CSS 4, Radix UI

The most purpose-built OpenClaw dashboard. Designed for headless servers over Tailscale.

**Install:**
```bash
npm install -g clawtrol && clawtrol init && clawtrol start
# → http://localhost:4781
```

**Config:** Reads `~/.openclaw/config.json` automatically.

| Module | What It Does |
|--------|-------------|
| **Sessions** | View + chat with OpenClaw sessions **in real-time** |
| **Terminal** | Full web terminal (ttyd) |
| **Screen** | Remote screen viewer with click interaction |
| **Files** | Browse, read, download files |
| **Memory** | Browse agent memory + markdown files |
| **Sub-agents** | Monitor and manage spawned sub-agents |
| **Logs** | Gateway log viewer with filtering |
| **Network** | Tailscale peers, running processes |
| **Tasks** | Lightweight kanban board |
| **Cron** | View/manage cron jobs |

**Themes:** Nova (cyberpunk), Midnight, Catppuccin, Solar
**Tailscale-native:** Built specifically for headless machine access over Tailscale.

---

### `clawfleet/ClawFleet` ⭐ 99 | Go + JavaScript | MIT
**URL:** https://github.com/clawfleet/ClawFleet
**Version:** v1.2.1 (released April 23, 2026)

Fleet manager for running multiple OpenClaw (and Hermes) instances from one browser dashboard.

**Install:**
```bash
curl -fsSL https://clawfleet.io/install.sh | sh
# → http://localhost:8080
```

**Standout features:**
- Multi-runtime: OpenClaw + Hermes Agent in one unified UI
- Each instance isolated in its own Docker container with **full XFCE desktop via noVNC in browser**
- **Character system** — reusable personas assigned to instances (bio, backstory, style)
- **Soul Archive** — save/clone instance state (memory, model config, conversation history)
- **Skill management** — 52 built-in + 13,000+ from ClawHub
- Real-time resource charts + live logs per instance
- Channel management (Telegram, Discord, Slack)
- `clawfleet shell` — drop into any instance's terminal
- Fleet collaboration — multiple AI agents working in group chats

**Resource:** ~1.5 GB RAM per instance. Tested on M4 MacBook Air.

**CLI:** `clawfleet create`, `configure`, `list`, `desktop`, `logs`, `snapshot save/list`, `dashboard serve`

---

### `builderz-labs/marketing-dashboard` ⭐ 232 | TypeScript/Next.js
**URL:** https://github.com/builderz-labs/marketing-dashboard

"Open-source marketing operations control center for AI agent teams powered by **OpenClaw** + SQLite."
CRM, outreach automation, content, analytics. OpenClaw-native, more specialized use case.

---

## 🟡 TIER 2 — FEATURE-RICHEST GENERAL AGENT UIs

Ranked by real-time capability + feature density.

---

### `danny-avila/LibreChat` ⭐ 36k | TypeScript | MIT
**URL:** https://github.com/danny-avila/LibreChat
**Best for:** Real-time resilience + full-featured agent workspace

| Feature | Detail |
|---------|--------|
| **Resumable Streams** | Reconnect/resume on dropped connections — works across tabs and devices |
| **Agent Marketplace** | No-code custom agents + collaborative sharing |
| **Code Interpreter** | Python, Node.js, Go, C/C++, Java, PHP, Rust, Fortran — sandboxed |
| **Generative UI** | React/HTML/Mermaid artifacts rendered live in chat |
| **MCP Support** | Native tool integration |
| All major providers | Anthropic, AWS Bedrock, OpenAI, Azure, Google Vertex AI |
| **Conversation branching** | Fork/branch conversations |
| Web search | With scrapers + rerankers |
| Multi-user | OAuth2, LDAP, Entra ID groups, granular RBAC |
| Redis-backed | Horizontal scaling |

---

### `open-webui/open-webui` ⭐ 134k | Python + Svelte | MIT
**URL:** https://github.com/open-webui/open-webui
**Best for:** Maximum feature count + enterprise readiness

**Highlights:** Voice/video calls (Whisper/ElevenLabs/Deepgram), Python Function Calling (BYOF), Pipelines plugin framework, multi-model simultaneous conversations, MCP support, 15+ web search providers, image generation, Redis WebSocket for multi-tab sync, LDAP/SCIM/SSO, OpenTelemetry, PWA.

**Install:** `pip install open-webui && open-webui serve` or Docker.

---

### `lobehub/lobehub` ⭐ 75.6k | TypeScript | MIT
**URL:** https://github.com/lobehub/lobehub
**Best for:** Agent-as-coworker paradigm + multi-agent collaboration
**Activity:** Commits hourly, 2,710 releases

**Standouts:**
- **Agent Groups** — teams of agents that collaborate like real coworkers
- **Personal Memory** — white-box, editable; agents learn from how you work
- 10,000+ tools + MCP plugins
- **Schedules** — run agents at specific times
- Floating chat panel
- Pages, Projects, Workspaces — full productivity suite

---

### `langgenius/dify` ⭐ 139k | TypeScript + Python | Apache 2.0
**URL:** https://github.com/langgenius/dify
**Best for:** Visual workflow building + LLMOps observability

Visual canvas for AI pipelines, 100+ LLM providers, RAG pipeline (PDF/PPT/docs), 50+ built-in agent tools, Grafana monitoring integration, Kubernetes/Helm support. More of a platform than a chat UI — better for building/orchestrating than daily interaction.

**Deploy:** `docker compose up -d` → `http://localhost/install`

---

### `Mintplex-Labs/anything-llm` ⭐ 58.9k | JavaScript | MIT
**URL:** https://github.com/Mintplex-Labs/anything-llm
**Best for:** No-code agent builder + MCP + desktop app

Full MCP compatibility, no-code agent builder, 30+ LLM providers, 9+ vector DBs, browser extension, desktop app (Mac/Win/Linux), custom embeddable chat widget, **80% token reduction via Intelligent Skill Selection**.

---

### `livekit/agents-playground` ⭐ 367 | TypeScript | Apache 2.0
**URL:** https://github.com/livekit/agents-playground
**Best for:** Voice/video real-time agent interaction

WebRTC-based. Renders video, audio, and chat from agents simultaneously. Send video/audio/text to your agent. LiveKit has a **Flutter SDK** — relevant for real-time voice in onemind-app.

**Live demo:** https://agents-playground.livekit.io

---

## Summary Rankings

| Criterion | Winner |
|-----------|--------|
| **Most OpenClaw-native** | Clawtrol > ClawFleet > marketing-dashboard |
| **Most real-time** | LibreChat (resumable streams) + LobeHub (live collab) |
| **Most feature-rich overall** | Open WebUI > LibreChat > LobeHub |
| **Best for Flutter/onemind-app** | AG-UI Protocol (Dart SDK) > LiveKit (Flutter SDK) |
| **Best workflow/orchestration** | Dify > LobeHub |
| **Best desktop app** | AnythingLLM |
| **Most actively developed** | LobeHub (hourly), AG-UI (daily) |

---

## The Big Opportunity

**AG-UI + Clawtrol** is the combination worth pursuing:

- **Clawtrol** → OpenClaw-native real-time sessions today, deployable immediately on VPS
- **AG-UI** → emerging standard (13.2k★, Dart SDK exists) — if OpenClaw ships AG-UI compatibility, `onemind-app` gets native real-time agent streaming via the existing WebSocket transport, zero proxy needed

The AG-UI Dart SDK path is particularly compelling: it would let the Flutter app receive live streaming events, state sync, and generative UI from any AG-UI-compatible backend — all over the existing `omos-core` Tailscale connection.

---

*Research completed: April 24–26, 2026*
*Sources: GitHub (direct repo fetch). All star counts and dates accurate as of research date.*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
