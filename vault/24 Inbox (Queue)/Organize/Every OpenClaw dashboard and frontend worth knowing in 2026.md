# Every OpenClaw dashboard and frontend worth knowing in 2026

OpenClaw (formerly Clawdbot → Moltbot → OpenClaw) has spawned **at least 25 distinct frontend dashboards and web UIs** across its ecosystem, ranging from single-file HTML monitors to full enterprise orchestration platforms. The project's headless, WebSocket-based architecture on port **18789** has made it a magnet for community-built interfaces. **No dedicated NATS integration or "One Mind OS/OMCS" dashboard exists** in the ecosystem — these terms returned zero results across all sources searched. However, the ecosystem is rich in robotics, 3D visualization, observability, n8n integration, and mobile-first options. Below is the exhaustive inventory.

---

## The official built-in surfaces

OpenClaw ships with two official web surfaces that serve as the baseline for comparison.

**Control UI (Dashboard v2)** is the built-in SPA served directly by the Gateway at `http://127.0.0.1:18789/`. As of v2026.3.12, it offers modular views for overview, chat, config, agents, and sessions. Authentication uses `OPENCLAW_GATEWAY_TOKEN` passed via URL fragment. It includes a command palette, slash commands, search, export, pinned messages, and mobile bottom tabs. Launch it with `openclaw dashboard`. This is the minimum viable frontend — functional but deliberately minimal, leaving room for community dashboards to differentiate.

**Canvas + A2UI** runs as a separate server on port **18793** and provides an agent-driven visual workspace. The A2UI (Agent-to-UI) v0.8 component library lets agents push interactive HTML/CSS/JS surfaces over WebSocket to connected clients. Available in the macOS app (WKWebView), iOS, Android, and the web Control UI. Canvas is architecturally isolated from the Gateway — if it crashes, the Gateway keeps running. The bundled Canvas skill can render dashboards, management interfaces, and interactive tools. Known limitation: A2UI JSON push is still marked WIP; raw HTML files are recommended.

---

## The top-tier community dashboards

These are the most feature-rich, actively maintained projects that function as genuine alternatives to the built-in UI.

**Mission Control** (github.com/crshdn/mission-control) — **440 stars**, TypeScript. The best overall orchestration dashboard. Features a clarifying-questions flow before task dispatch, automatic agent spawning per task, a **7-stage kanban board** (Planning → Inbox → Assigned → In Progress → Testing → Review → Done), live WebSocket event feed, and multi-machine support over Tailscale. Uses SQLite with no external DB. Security via bearer token auth and Zod validation. Lacks cost tracking and memory file viewing.

**OpenClaw Studio** (github.com/grp06/openclaw-studio) — **649 stars**, Node.js. The developer power tool. Architecture uses server-owned WebSocket to Gateway with Browser → Studio communication via HTTP + SSE. Features a **visual node-based workflow editor**, raw JSON log viewer, **"time-travel" debugging** (rewind agent memory to see context at any step), agent creation, cron jobs, exec approvals, and runtime.db for durable state projection. Install: `npx -y openclaw-studio@latest`. Requires `STUDIO_ACCESS_TOKEN` for non-loopback access. Feels like VS Code crossed with Postman.

**LobsterBoard** (github.com/Curbob/LobsterBoard) — **382 stars**, vanilla JS/CSS, zero frontend dependencies. A dashboard *builder* with **50 widget definitions**: system monitoring, AI cost tracking (12 providers including Claude, Codex, Copilot, Cursor, Gemini), weather, stocks, crypto, cameras, smart home, calendars, RSS. Drag-and-drop editor with snap grid and resize handles. Template gallery for sharing layouts. Custom pages system. OpenClaw-specific widgets: auth status, cron jobs, system log, active sessions, release version check. Remote multi-server monitoring via `lobsterboard-agent` with **E2E encryption (ECDH P-256 + AES-256-GCM)**. License: BSL-1.1 (free personal, commercial requires license). Website: lobsterboard.com.

**AI Maestro** (github.com/23blocks-OS/ai-maestro) — **324 stars**, TypeScript. Built for large-scale multi-agent operations across **80+ agents on multiple machines**. Uses a **peer mesh network** with no central server. Features the Agent Messaging Protocol (AMP) for direct agent-to-agent communication with priority levels and crypto signatures. Works with any coding agent (Claude Code, Aider, Cursor, Copilot) — not locked to OpenClaw. Three memory layers: conversation memory, code graph, and auto-generated documentation. Interactive codebase visualization. Gateway integrations for Slack, Discord, Email, WhatsApp.

**Clawe** (github.com/getclawe/clawe) — **261 stars**, TypeScript. A pre-configured **4-agent squad** (Clawe/Squad Lead, Inky/Content Editor, Pixel/Designer, Scout/SEO) with Trello-style kanban coordination. Heartbeat scheduling wakes agents every 15 minutes, staggered to avoid rate limits. @mention notifications tag agents for next-wake pickup. Requires Convex backend (free tier). Unique in its opinionated multi-agent role structure.

---

## Mid-tier and specialized dashboards

**Clawd Control** (github.com/Temaki-AI/clawd-control) — **101 stars**, vanilla HTML/JS. The lightest fleet monitor with only one dependency (`ws`). Auto-discovers local OpenClaw agents, provides an agent creation wizard, host metrics (CPU/RAM/disk), fleet overview with health indicators, SPA navigation, and keyboard shortcuts. No build step needed. Legacy Clawdbot-era project, still functional.

**VidClaw** (github.com/madrzak/vidclaw) — **67 stars**, JavaScript. Combines a kanban task board (auto-executes via 2-minute cron or 30-minute heartbeat) with a **soul editor** for SOUL.md, IDENTITY.md, and USER.md with version history and persona templates. Model switching hot-reloads from the dashboard. Skills manager lets you view, enable/disable, and create custom skills. Usage tracking parses session transcripts for token costs.

**tugcantopaloglu/openclaw-dashboard** — **60 stars**, pure Node.js, zero npm dependencies. The **most security-hardened** community dashboard: PBKDF2 password hashing (100K iterations, SHA-512), **TOTP MFA** (Google Authenticator compatible), rate limiting (5 soft / 20 hard lockout), HSTS and CSP headers. Features cost analysis by model/session/time period, memory viewer (MEMORY.md, HEARTBEAT.md, daily notes), cron management, **Docker container management** (view/start/stop/restart containers and images), UFW/fail2ban security dashboard, and full audit logging. Tailscale integration shows status, IP, and connected peers; HTTPS enforced except from localhost and Tailscale IPs.

**Claw Dashboard** (github.com/spleck/claw-dashboard) — **49 stars**, JavaScript. A terminal-style TUI inspired by btop/htop using the `blessed` library. ASCII art with gradient colors and donut charts. Per-core CPU, memory gauges, GPU stats (Apple Silicon), top processes, active sessions. Auto-refreshes every 2 seconds. Keyboard controls: q to quit, r to refresh, s for settings.

**Clawtrol** (github.com/nachoiacovino/clawtrol) — **18 stars**, TypeScript. The most complete **all-in-one remote panel**: remote screen viewer (view/click headless Mac from browser), web terminal via ttyd, file browser, session viewer with live chat, kanban board, memory browser, cron manager, sub-agent monitor. Four themes (Nova, Midnight, Catppuccin, Solar). Modular — enable only the panels you need. Uses pm2. No built-in auth (designed for Tailscale/trusted networks).

---

## Observability and monitoring dashboards

**ClawMetry** (github.com/vivekchand/clawmetry) — Python/Flask, single-file `dashboard.py`, **~2MB footprint**. Described as **"Grafana but purpose-built for AI."** Features: animated SVG architecture diagram showing live message flow (the signature "Flow" tab), health checks, activity heatmap, token/cost tracking (daily/weekly/monthly), active sessions with model info, cron job monitoring, color-coded real-time log streaming, memory browser (SOUL.md, MEMORY.md, AGENTS.md), chat-bubble transcript viewer, sub-agent monitoring (files read, commands run, tools called). Supports OpenTelemetry for real token counts. Install: `pip install clawmetry && clawmetry` — zero config, auto-detects workspace. Runs on Raspberry Pi. macOS menu bar app available. **198+ Product Hunt upvotes**. MIT license.

**Grafana Lens** (github.com/awsome-o/grafana-lens) — An OpenClaw plugin providing **17 composable agent tools** for full Grafana integration: query PromQL/LogQL/TraceQL, create dashboards, set alerts, render charts, investigate incidents, push custom metrics. Includes **12 pre-built dashboard templates** and SRE investigation features. Not a standalone dashboard but the bridge between OpenClaw and the Grafana/Prometheus stack.

The Gateway itself exposes a native **Prometheus `/metrics` endpoint** and supports OTLP export, with metrics including `openclaw_gateway_requests_total`, `openclaw_gateway_response_time_seconds`, and `openclaw_agent_sessions_active`. A companion **monitor-gen skill** generates production-ready Prometheus alerting rules and Grafana dashboard JSON from plain-language descriptions.

---

## 3D, spatial, and gamified visualizations

**TenacitOS** (github.com/carlosazaustre/tenacitOS) — A mission control dashboard with an **interactive 3D office built with React Three Fiber** — one voxel avatar desk per agent. Also includes cost tracking, cron manager, memory browser, and file browser. Supports custom branding via environment variables (`src/config/branding.ts`), with per-agent emoji and color definitions. The 3D office visualization is unique in the ecosystem.

**openclaw-desktop-pet** (github.com/44-99/openclaw-desktop-pet) — A **Three.js desktop pet** running in Electron that integrates with OpenClaw. Shows real-time AI status, has an emotion system, and uses a Python backend for system monitoring. The pet reacts to agent state changes.

**Moltcraft** (github.com/askmojo/moltcraft, moltcraft.xyz) — A **gamified pixel-art isometric world** where agents become clickable pixel characters walking around, mining tokens, and completing tasks. Click buildings for real data: cron jobs, token usage, skills, channels. Features **voice I/O** (talk to agents, hear responses), day/night cycle, weather, procedural audio. Install: `npx @ask-mojo/moltcraft`. Pure HTML/CSS/JS, zero npm dependencies, ~2MB. Originally built for Moltbot but compatible with the broader ecosystem. Runs on Raspberry Pi.

**PhysicClaw-VEA** — Discussed on the Three.js forum. An interactive 3D app with a virtual entity that reacts in real-time through shaders, emotions, and voice. Stack: React 19, Three.js WebGPU, Supabase, TypeScript.

---

## Mobile-first and PWA dashboards

**ClawBridge** (github.com/dreamwing/clawbridge, clawbridge.app) — The dedicated **mobile-first PWA dashboard**. Installable via "Add to Home Screen." Features live activity feed, token economy tracking, memory timeline, mission control (trigger cron jobs, restart services). Uses Cloudflare Quick Tunnel for zero-config remote access. The strongest mobile-first option.

**ClawApp** (github.com/qingchencloud/clawapp) — Mobile chat client with offline message caching via IndexedDB, WebSocket proxy, PWA + APK packaging via Capacitor. Supports Chinese and English.

**OpenClaw Dashboard for Android** — Available on Google Play (com.cbstudio.openclaw_dashboard). Native Android app for Gateway and agent management.

**openclaw-termux** (github.com/mithun50/openclaw-termux) — A Flutter app that runs the **full OpenClaw Gateway on Android**. Built-in terminal, web dashboard via WebView, SSH access, and node device capabilities. Turns any Android phone into a full OpenClaw node.

---

## Additional dashboards and cyberpunk variants

**0xChris-Defi/openclaw-dashboard** — React 19, Tailwind CSS 4, shadcn/ui, Express 4, tRPC 11, MySQL/TiDB with Drizzle ORM. **Cyberpunk terminal aesthetic.** Features OAuth 2.0 + WalletConnect auth, multi-model streaming chat, gateway management, multi-channel support (Telegram, Discord, Slack, WhatsApp), and real-time metrics. The only dashboard with Web3 wallet authentication.

**mudrii/openclaw-dashboard** — **5 stars** but notable for its zero-dependency philosophy: Go server + Python data refresh, single-file HTML frontend — no npm, no pip, no CDN, no build tools. Ten dashboard panels including alerts, system health, cost cards, kanban, cron jobs, token usage, sub-agent activity. AI chat via Gateway's OpenAI-compatible endpoint. **6 themes** (Midnight, Nord, Catppuccin Mocha, GitHub, Solarized, Catppuccin Latte) with glassmorphism UI. Nix flake and Docker support.

**OpenMind** (github.com/JozefJarosciak/OpenMind) — Turns the OpenClaw workspace into a **fully interactive, live-editable mind map**. Edit any memory node directly, hot-swap logic, search across everything. Docker-based. Agent-installable (ask your OpenClaw to install it). Supports Tailscale network restriction. AGPL-3.0 license. Early alpha.

**MoltBoard** (github.com/TheSethRose/MoltBoard) — Next.js App Router, SQLite. Task and project management dashboard from the Moltbot era. Task lifecycle (backlog → ready → in-progress → completed), GitHub integration, system monitoring.

**Chrome Extension** — Available in the Chrome Web Store (nbnmkgfibnkhpmblefimbfppfoncblda). Browser sidepanel connecting directly to local Gateway via WebSocket for low-latency interaction without opening a separate tab.

---

## Robotics, vision, and n8n integrations

**RosClaw** (github.com/PlaiPin/rosclaw) is the leading **ROS2 + OpenClaw integration**, enabling natural language control of robots through messaging apps. Includes a TypeScript rosbridge WebSocket client and an OpenClaw plugin, with a planned real-time dashboard (`@rosclaw/openclaw-canvas` — Phase 3). **AgenticROS** (agenticros.com) offers a similar natural language → ROS2 bridge supporting multiple AI agents. The **OpenClaw Robotics Community** (openclawrobotics.com) documents **ClawBody**, a bridge connecting OpenClaw to physical hardware with MuJoCo simulation support.

**VisionClaw** (github.com/sseanliu/VisionClaw) provides **real-time AI vision** for Meta Ray-Ban smart glasses with voice + vision + agentic actions via Gemini Live and OpenClaw. Supports phone camera mode and WebRTC streaming to a browser viewer with bidirectional audio/video.

For **n8n workflow integration**, **openclaw-n8n-stack** (github.com/caprihan/openclaw-n8n-stack) provides a pre-configured Docker Compose stack combining OpenClaw + n8n with shared network and webhook templates. **n8n-claw** (github.com/freddy-schuetz/n8n-claw) recreates OpenClaw's autonomous agent loop entirely within n8n, with adaptive RAG-powered memory and MCP skill templates. An official n8n-workflow skill exists at playbooks.com.

---

## Voice, TTS, and macOS native apps

**TTS is a first-class built-in feature** supporting OpenAI TTS, ElevenLabs, and Microsoft Edge TTS (free). Configurable under `messages.tts` in openclaw.json. Supports voice notes on Telegram (Opus), Talk Mode on macOS/iOS/Android with wake words.

**VoxClaw** (github.com/malpern/VoxClaw) is a macOS menu bar app that lets OpenClaw speak aloud on any Mac on the network, supporting Apple TTS, OpenAI TTS, and ElevenLabs with a teleprompter overlay. Includes an iPhone companion app. **OpenClaw Voice** (openclawvoice.com) provides browser-based voice chat with Whisper STT and ElevenLabs TTS via self-hosted WebSocket audio streaming.

The **official OpenClaw.app** is a native Swift macOS menu bar companion: voice wake/push-to-talk, Talk Mode overlay, WebChat, Canvas, debug tools, remote gateway control via SSH tunnels and Tailscale discovery. Built with Swift Package Manager, universal binary (Apple Silicon + Intel), requires macOS 15+.

---

## Enterprise and white-label options

**AgentOps Dashboard** (referenced in Jumei.ai comparison) is the only true enterprise-grade option: RBAC, fleet management (restart/pause/update multiple instances simultaneously), compliance and audit logging. Requires a dedicated PostgreSQL database and has the highest setup complexity.

**ClawWrapper** (clawwrapper.com) is a **$149 commercial starter kit** for building an OpenClaw wrapper SaaS with white-label branding via environment variables, Fly.io multi-tenant deployment, Telegram bot integration, Google OAuth via Supabase, and Stripe/Lemon Squeezy billing. **openclawdashboard.com** claims enterprise-grade E2E encryption and RBAC but could not be verified as open-source.

---

## Comprehensive feature matrix

| Dashboard | Stars | Stack | Gateway WS | Auth | Cost Track | Multi-Agent | Mobile | Voice | 3D | Tailscale | Docker |
|---|---|---|---|---|---|---|---|---|---|---|---|
| **Official Control UI** | — | Built-in | ✅ Direct | Token | ❌ | Routing | Bottom tabs | ❌ | ❌ | ✅ Built-in | ❌ |
| **Mission Control** (crshdn) | 440 | TypeScript | ✅ | Bearer | ❌ | ✅ Auto-spawn | ❌ | ❌ | ❌ | ✅ | ❌ |
| **OpenClaw Studio** | 649 | Node.js | ✅ Server-owned | Token | ✅ Per-step | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **LobsterBoard** | 382 | Vanilla JS | Optional | API key | ✅ 12 providers | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **AI Maestro** | 324 | TypeScript | ✅ | Crypto-signed | ❌ | ✅ 80+ agents | ❌ | ❌ | ❌ | ❌ | ❌ |
| **Clawe** | 261 | TypeScript | ✅ | Convex | ❌ | ✅ Squad | ❌ | ❌ | ❌ | ❌ | ❌ |
| **ClawMetry** | PH:198 | Python/Flask | OTLP | Local-only | ✅ Deep | ✅ Sub-agents | ❌ | ❌ | ❌ | ❌ | ❌ |
| **Clawd Control** | 101 | HTML/JS | ✅ | ❌ | ❌ | ✅ Fleet | ❌ | ❌ | ❌ | ❌ | ❌ |
| **VidClaw** | 67 | JavaScript | ✅ | SSH tunnel | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **tugcantopaloglu** | 60 | Node.js | ✅ | ✅ MFA/TOTP | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ | ✅ Mgmt |
| **Clawtrol** | 18 | TypeScript | ✅ | Tailscale-only | ❌ | ✅ Sub-agents | ❌ | ❌ | ❌ | ✅ | ❌ |
| **TenacitOS** | — | React/R3F | ✅ | ❌ | ✅ | ✅ | ❌ | ❌ | ✅ 3D office | ❌ | ❌ |
| **ClawBridge** | — | PWA | ✅ | ❌ | ✅ Tokens | ❌ | ✅ Mobile-first | ❌ | ❌ | ❌ | ❌ |
| **Moltcraft** | — | HTML/JS | Moltbot | Token | ✅ Tokens | ✅ Visual | ❌ | ✅ Voice I/O | ✅ Isometric | ❌ | ❌ |
| **0xChris-Defi** | — | React 19/tRPC | ✅ | OAuth+Web3 | ❌ | ✅ Multi-channel | ❌ | ❌ | ❌ | ❌ | ❌ |
| **mudrii** | 5 | Go/Python | ✅ | Rate-limit | ✅ Deep | ✅ Sub-agents | ✅ Responsive | ❌ | ❌ | ❌ | ✅ |
| **AgentOps** | — | PostgreSQL | ✅ | ✅ RBAC | ✅ Org-level | ✅ Fleet | ❌ | ❌ | ❌ | ❌ | ❌ |

---

## What does not exist (yet)

Several categories the user asked about returned **no results across all platforms searched**:

- **NATS/JetStream integration**: Only one blog post by @ggondim mentioned considering NATS for event routing before abandoning it in favor of native OpenClaw features. No dashboard, plugin, or skill integrates with NATS. The NATS ecosystem has its own Grafana dashboards but none are OpenClaw-aware.
- **"One Mind OS" / OMCS**: These terms returned zero results. No project, feature, or concept with these names exists in the OpenClaw ecosystem. Some community members metaphorically describe OpenClaw as a "personal OS," and there's an `ai-persona-os` ClawHub skill, but nothing called "One Mind OS."
- **Canvas/A2UI extension dashboards**: While Canvas and A2UI are real and well-documented, no third-party project specifically *extends* A2UI with new components. The Telegram Canvas skill (`openclaw-tg-canvas`) is the closest, rendering Canvas content in Telegram Mini Apps.

---

## Recommended 3-tool stack for a sovereign AI life OS

For the described infrastructure (Mac Mini + DigitalOcean bare metal + Mason Android nodes + Tailscale mesh + Discord as primary channel), the optimal combination is:

**1. Clawtrol** as the primary command center. It is the only dashboard offering remote screen viewing, web terminal, file browsing, and session management in a single panel — essential for managing headless Mac Mini and bare metal nodes over Tailscale. Its explicit design for Tailscale-trusted networks means zero auth friction on your mesh. The modular panel system lets you enable only what you need per device.

**2. ClawMetry** as the observability layer. Its animated Flow visualization, token cost tracking, session monitoring, and sub-agent activity views provide the operational awareness a life OS demands. At ~2MB with zero config and auto-workspace detection, it runs comfortably on every node including Android via Termux. The OpenTelemetry support bridges to your Grafana/Prometheus stack if needed later. Its Raspberry Pi compatibility confirms it will run on resource-constrained Mason nodes.

**3. Mission Control** (crshdn) as the orchestration and task management layer. Its automatic agent spawning, 7-stage kanban, and multi-machine Tailscale support make it the natural fit for coordinating work across Mac Mini, DigitalOcean, and Android nodes. The clarifying-questions flow before task dispatch prevents runaway token costs — critical for a system running autonomously.

**Supplement with**: the built-in Canvas/A2UI for quick interactive surfaces, `openclaw-n8n-stack` if you need visual workflow orchestration, and RosClaw if any Mason nodes run robotics workloads. For Discord as primary channel, all three tools work alongside OpenClaw's native Discord bridge — they monitor and manage the agents, while Discord remains the conversational interface. For NATS JetStream event monitoring, you would need to build a custom Canvas skill or use Grafana Lens with the NATS JetStream Grafana dashboard, as no native bridge exists.