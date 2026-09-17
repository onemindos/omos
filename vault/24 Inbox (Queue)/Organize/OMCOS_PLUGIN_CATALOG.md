# One Mind OS — Plug-and-Play Project Catalog
**What already exists. What to clone. What to install. No wiring required.**
*Legacy Platform Holdings LLC · March 2026*

---

> **The point of this document**: Stop building what's already been built.
> The community has shipped working integrations for almost every layer of the stack.
> This is the map of what to pull in vs what to actually write yourself.

---

## How to Read This

Each entry has:
- **What it replaces** — what you would have had to custom-wire
- **Status** — Working / Fork / ClawHub install
- **Source** — exact repo or install command
- **Security note** — always review code before enabling

Three integration types in OpenClaw:
- **Skill** — SKILL.md in your workspace. Natural language instructions + optional scripts. Lowest trust boundary.
- **Plugin** — TypeScript/JS running inside the Gateway process. Deep access. Higher trust.
- **Fork** — Full repo you clone, modify, and own.

---

## Layer 1 — Realtime Voice + Vision Bridge

### VisionClaw
**What it replaces**: Building a Gemini Live → OpenClaw bridge from scratch

Ray-Ban Meta glasses stream at ~1fps to Gemini for visual context, while audio flows bidirectionally in real-time. Gemini Live supports function calling — a single `execute` tool routes everything through OpenClaw skills.

- **Type**: Fork (iOS + Android)
- **Repo**: `github.com/sseanliu/VisionClaw`
- **Stack**: Meta DAT SDK + Gemini Live WebSocket + OpenClaw `/v1/chat/completions`
- **What you modify**: `Secrets.swift` / `Secrets.kt` — gateway host, token, Gemini key
- **Gives you**: Glasses camera → Gemini sees → Legacy acts via all your skills
- **Security**: Review before build. App makes outbound calls to Gemini + your gateway.

### NixClaw ⭐ Recommended Fork
**What it replaces**: VisionClaw + background audio + named agent routing

NixClaw adds an `execute` tool that lets Gemini call your OpenClaw agent by name, routing actions through the gateway, plus background audio mode so the conversation keeps running when you lock your phone or switch apps.

- **Type**: Fork (iOS + Android, VisionClaw fork)
- **Repo**: `github.com/arniesaha/NixClaw`
- **What you modify**: Agent name → "Legacy", gateway host/token/Gemini key
- **Gives you**: Always-on ambient voice, background audio, agent-named routing, 3 modes (glasses / iPhone camera / audio-only)
- **Tailscale**: Add ATS exception in `Info.plist` for HTTP over Tailscale IPs — works over your mesh natively
- **Security**: Same as VisionClaw. Review Secrets config before building.

```swift
// Secrets.swift — what you change
static let geminiApiKey = "YOUR_GEMINI_KEY"
static let openClawHost = "http://your-tailscale-hostname"
static let openClawPort = 18789
static let openClawGatewayToken = "YOUR_GATEWAY_TOKEN"
```

---

## Layer 2 — Gateway Config (Enable This, Not Custom Build)

### OpenClaw Gateway — chatCompletions Endpoint
**What it replaces**: Writing a custom API bridge

This is a single config line. Enables the `/v1/chat/completions` HTTP endpoint that VisionClaw/NixClaw POST to. Legacy receives the call with full skill + workspace access.

```json
{
  "gateway": {
    "http": {
      "endpoints": {
        "chatCompletions": { "enabled": true }
      }
    }
  }
}
```

No code. One line. Done. The bridge maintains a rolling window of 20 messages to give OpenClaw context without blowing up token limits. Session continuity is handled via an `x-openclaw-session-key` header.

### OpenClaw Native Realtime Voice (OpenAI path)
**What it replaces**: Custom Twilio/WebRTC audio pipeline

A merged PR adds full conversation pipeline with `sttProvider: "openai-realtime-conversation"` that routes audio directly through OpenAI Realtime API for STT + LLM + TTS in one WebSocket hop, bypassing the Pi agent for sub-second response latency. Server-side VAD handles turn detection automatically.

- **Type**: Built-in OpenClaw plugin (voice-call)
- **Config**: Set `sttProvider: "openai-realtime-conversation"` in discord-voice or voice-call plugin
- **Tradeoff**: Realtime voice but OpenAI model handles voice turn — Claude still handles tool execution

---

## Layer 3 — Smart Home Integration

### OpenClaw Home Assistant Add-on (HAOS)
**What it replaces**: Custom HA API wiring

This repository contains a Home Assistant add-on that runs OpenClaw inside Home Assistant OS. Includes AI Gateway, web terminal embedded in HA, Assist Pipeline using OpenClaw as conversation agent via OpenAI-compatible API, browser automation with Chromium, and persistent storage across updates.

- **Type**: Home Assistant Add-on (HAOS native)
- **Repo**: `github.com/techartdev/OpenClawHomeAssistant`
- **Install**: Add repo to HA Add-on Store → install OpenClaw add-on
- **Gives you**: OpenClaw runs inside HA, accesses all 2,500+ HA integrations, Frigate NVR events, Zigbee/Matter devices

### OpenClaw HA Custom Integration (HACS)
**What it replaces**: Manually calling the HA REST API from skills

Connects Home Assistant to OpenClaw and adds: native conversation agent support for Assist, Lovelace chat card with message history, service calls for automations/scripts, event hooks for response-driven automations, optional voice mode via Browser or HA STT Assist, and expose-aware context support.

- **Type**: HACS custom integration
- **Source**: HA Community thread `community.home-assistant.io/t/989054`
- **Install**: HACS → Custom Repositories → add repo URL → install
- **Gives you**: Bidirectional HA ↔ OpenClaw. HA events trigger Legacy. Legacy controls any HA entity.

### Home Assistant Skill (ClawHub)
**What it replaces**: Writing curl scripts to hit the HA REST API

```bash
clawhub install homeassistant
# Set env vars:
export HA_URL=http://homeassistant.local:8123
export HA_TOKEN=your_long_lived_access_token
```

Gives full conversational access to your Home Assistant setup — sensor history, automation triggers, entity states, service calls. Supports context memory (say "Turn on AC" then "Lower it" — it knows what you mean).

- **Requires**: `curl`, `HA_URL`, `HA_TOKEN` env vars
- **Security**: L1 (low risk). Review before enabling.

---

## Layer 4 — Memory & Persistence

### agent-brain (ClawHub)
**What it replaces**: Custom SQLite memory layer

Local-first persistent memory with SQLite storage, orchestrated retrieve/extract loops, hybrid semantic + keyword search. The highest-starred skill on ClawHub (132 stars).

```bash
clawhub install agent-brain
```

- **Gives you**: Legacy remembers facts, preferences, and context across sessions with retrieval
- **Storage**: Local SQLite — sovereign, no cloud
- **Security**: Review before enabling. Has file write access.

### shared-memory (ClawHub)
**What it replaces**: Building cross-agent memory sync

Multiple OpenClaw agents sharing state — useful for multi-location setup (home Legacy ↔ farm Legacy ↔ office Legacy).

```bash
clawhub install shared-memory
```

- **Gives you**: Legacy instances across Tailscale mesh can share and sync memory state
- **Security**: Review carefully — cross-agent trust boundary.

### active-maintenance (ClawHub)
**What it replaces**: Manual memory cleanup and health monitoring

Automated system health and memory metabolism. Keeps Legacy's context clean, prunes stale memory, monitors agent health.

```bash
clawhub install active-maintenance
```

---

## Layer 5 — Automation & Workflow

### n8n-claw
**What it replaces**: Building n8n ↔ OpenClaw bridge from scratch

OpenClaw-inspired autonomous AI agent built entirely in n8n. Adaptive RAG-powered memory via Supabase, Skills via MCP templates, proactive task management, media understanding — self-hosted with one setup script. Install pre-built skills from skill catalog — just ask your agent: "What skills are available?" "Install weather-openmeteo."

- **Type**: Full repo (n8n + Supabase stack)
- **Repo**: `github.com/freddy-schuetz/n8n-claw`
- **Install**: Clone → run setup script → VPS ready
- **Gives you**: n8n workflows triggerable by Legacy, Supabase persistence, MCP skill builder
- **Note**: This is the n8n-native OpenClaw alternative. For your stack, use the **n8n ClawHub skill** instead (below) to connect your existing n8n to OpenClaw gateway.

### n8n Skill (ClawHub)
**What it replaces**: Custom webhook wiring between n8n and Legacy

```bash
clawhub install n8n
# Configure with your n8n instance URL
```

Connects Legacy directly to your existing n8n instance. Create, trigger, and manage n8n workflows via natural language. Set up a workflow that saves Gmail attachments to Dropbox and notifies you on Slack — entirely through a chat message. Because n8n runs locally, your data stays private.

### agent-autopilot (ClawHub)
**What it replaces**: Custom task scheduler + progress tracking

Self-driving agent workflow with heartbeat-driven task execution, day/night progress reports, and long-term memory.

```bash
clawhub install agent-autopilot
```

- **Gives you**: Legacy autonomously manages outstanding tasks between sessions, self-reporting

---

## Layer 6 — Robotics

### RosClaw
**What it replaces**: Building a custom ROS2 bridge from scratch

RosClaw connects OpenClaw to ROS2 through an intelligent plugin layer. User message → OpenClaw Gateway (AI Agent + Tools + Memory) → RosClaw Plugin → rosbridge WebSocket → ROS2 DDS → Nav2, MoveIt2, cameras, sensors.

- **Type**: OpenClaw plugin + ROS2 package
- **Repo**: `github.com/PlaiPin/rosclaw`
- **Install**: Clone → `docker compose up` (starts ROS2 + rosbridge + Gazebo) → add plugin to `openclaw.json`
- **Config**: Point plugin at `ws://localhost:9090`
- **Gives you**: Natural language → Legacy → ROS2 → Unitree G1 navigation and motion

```json
{
  "plugins": {
    "entries": {
      "rosclaw": {
        "enabled": true,
        "config": {
          "rosbridgeUrl": "ws://localhost:9090"
        }
      }
    }
  }
}
```

### peaq Robotics SDK (ClawHub)
**What it replaces**: Building DID machine identity and onchain robot registration

```bash
clawhub install peaq-robotics
```

Robots can now connect to OpenClaw agents and receive and deploy the peaq Robotics SDK directly via those agents. Fully ROS2 compatible, runs directly on robots as a lightweight node. Enables robot DID identity (ERC-8004), laying foundation for robots to become full economic agents.

### unitree-robot (ClawHub)
**What it replaces**: Manually coding Unitree SDK calls

```bash
clawhub install unitree-robot
```

Control Unitree G1, H1, GO1, GO2 via natural language through Legacy. "Forward 1m" → parsed → robot executes.

### robot (Official OpenClaw Skill)
**What it replaces**: Writing ROS2 / Arduino / motion planning code from scratch

```bash
clawhub install robot
```

Guides you from hardware wiring to ROS2 and safety in robotics projects, generating hardware plans, code, and debugging strategies. Covers Arduino/ESP32 wiring, ROS2 Humble configuration, MoveIt2 trajectory generation, safety constraints.

---

## Layer 7 — Security

### SecureClaw
**What it replaces**: Manual security auditing and permission management

Two-layer defense model: a code-level plugin that enforces hardening at the gateway and configuration level, combined with a behavioral skill that operates inside the agent's context window. Covers configuration-level checks and operational controls. Includes formal mappings to MITRE ATLAS agentic AI attack techniques.

- **Type**: Plugin + Skill (two components)
- **Repo**: Adversa AI — search `SecureClaw` on GitHub
- **Gives you**: Automated gateway hardening, permission auditing, behavioral guardrails
- **Install first** — before opening the gateway to external traffic

### agent-audit-trail (ClawHub)
**What it replaces**: Custom logging and audit system

```bash
clawhub install agent-audit-trail
```

Tamper-evident, hash-chained audit logging for all agent actions. Every tool call Legacy makes is logged with cryptographic integrity.

### agent-access-control (ClawHub)
**What it replaces**: Custom permission gating

```bash
clawhub install agent-access-control
```

Tiered access control. Define what Legacy can do without explicit approval vs what requires confirmation.

---

## Layer 8 — Productivity & Content

### Todoist (via Rube/Composio — already working)
Already confirmed working in your stack via `RUBE_MULTI_EXECUTE_TOOL`. Inbox project ID `6g6P3vvjx9643JrW`. No changes needed.

### gog — Google Workspace Suite (ClawHub)
**What it replaces**: Separate Gmail + Calendar + Drive integrations

```bash
clawhub install gog
```

Unified Google Workspace skill. Gmail triage, Calendar management, Sheets/Docs read/write — all through Legacy. Turns OpenClaw from assistant to time strategist — triaging emails, scheduling meetings, pulling spreadsheet data, drafting documents through natural conversation.

### adaptlypost (ClawHub)
**What it replaces**: Custom content distribution wiring

```bash
clawhub install adaptlypost
```

Schedule and manage posts across Instagram, X, Bluesky, TikTok, Threads, LinkedIn, Facebook — all via Legacy. Plugs directly into your content distribution machine.

### summarize (Official OpenClaw — bundled)
Already bundled. 10,000+ installs. Enable and use immediately.

---

## Layer 9 — Infrastructure & Monitoring

### OpenClaw on DigitalOcean 1-Click Deploy
**What it replaces**: Manual VPS setup and hardening

DigitalOcean offers 1-Click OpenClaw Deploy with a hardened security image. Security-hardened, production-ready agentic AI.

- **Gives you**: Pre-hardened Ubuntu image with OpenClaw pre-installed, firewall configured, ready in 2 minutes
- **URL**: digitalocean.com → Marketplace → OpenClaw
- **Cost**: Standard DO droplet pricing

### nix-openclaw
**What it replaces**: Manual dependency management and service configuration

```bash
# Nix flake approach — batteries-included
# Installs gateway + tools everywhere
# launchd service on macOS, systemd user service on Linux
```

- **Repo**: `github.com/openclaw/nix-openclaw`
- **Gives you**: Declarative config, all bundled plugins toggleable as Nix options, reproducible across all nodes
- **Ideal for**: Managing OpenClaw across your Tailscale mesh consistently

### Peekaboo (Official OpenClaw skill)
**What it replaces**: Custom macOS screenshot / vision capture

```bash
brew install steipete/tap/peekaboo
clawhub install peekaboo
```

Fast macOS screenshots with optional AI vision analysis. Uses macOS Accessibility APIs. Powers Bambu Labs 3D printer app automation and any Mac UI vision task.

---

## What You Actually Write Yourself

After pulling in everything above, this is all that's left to build:

| File | Why You Write It | Lines of Code |
|---|---|---|
| `SOUL.md` | Legacy's identity is yours alone | ~60 lines |
| `HEARTBEAT.md` | Your monitoring checklist | ~30 lines |
| `AGENTS.md` | Your tool access rules | ~50 lines |
| `legacy-eyes.py` (SmolVLM path) | Local vision → AMBIENT_STATE.md | ~100 lines |
| `legacy-pulse.py` | NATS → AMBIENT_STATE.md | ~120 lines |
| `legacy-session-start.py` | Session init event | ~40 lines |
| `launch.sh` | Stack startup orchestration | ~60 lines |

That's it. Everything else is already built and battle-tested by the community.

---

## Full Install Sequence

### Step 1 — Core Gateway

```bash
npm install -g openclaw@latest
openclaw onboard --install-daemon
```

### Step 2 — Essential Skills

```bash
# Memory
clawhub install agent-brain
clawhub install shared-memory
clawhub install active-maintenance

# Robotics
clawhub install peaq-robotics
clawhub install unitree-robot
clawhub install robot

# Vision
clawhub install @steipete/camsnap

# Automation
clawhub install n8n
clawhub install agent-autopilot

# Productivity
clawhub install gog
clawhub install adaptlypost

# Security
clawhub install agent-audit-trail
clawhub install agent-access-control
```

### Step 3 — Plugins (in openclaw.json)

```json
{
  "plugins": {
    "entries": {
      "discord-voice": { "enabled": true },
      "camsnap": { "enabled": true },
      "rosclaw": { "enabled": true }
    }
  }
}
```

### Step 4 — Home Assistant

```bash
# Option A: HACS custom integration (bidirectional)
# Go to: community.home-assistant.io/t/989054

# Option B: HA Add-on (OpenClaw runs inside HAOS)
# Repo: github.com/techartdev/OpenClawHomeAssistant
```

### Step 5 — Realtime Vision Bridge

```bash
# Fork NixClaw for Mason Android nodes
git clone https://github.com/arniesaha/NixClaw.git
cd NixClaw
# Edit Secrets.kt with your gateway host/token/Gemini key
# Build via Android Studio → deploy to Mason MASC09T
```

### Step 6 — Robotics Bridge

```bash
git clone https://github.com/PlaiPin/rosclaw.git
cd rosclaw
docker compose up  # ROS2 + rosbridge + Gazebo simulation
# Add rosclaw plugin to openclaw.json
```

### Step 7 — Security Layer

```bash
# Search GitHub for SecureClaw by Adversa AI
# Install plugin + skill following their README
# Run: openclaw security audit --deep
```

---

## Security Warning — Read Before Installing Any Community Skill

OpenClaw lets any developer publish a skill file on ClawHub, making it easy for developers to inject malicious instructions into markdown files and compromise systems. One known vulnerability (CVE-2026-25253) compromises the AI gateway.

**Rules before installing anything:**

1. Check the VirusTotal report on the skill's ClawHub page
2. Read the SKILL.md source yourself — it's a markdown file, takes 2 minutes
3. Check who published it and how many installs it has
4. For plugins: read the TypeScript source — plugins run inside your gateway process
5. Never install a skill that requests more permissions than it needs
6. Run `openclaw security audit --deep` after any new install

The skills marked ⭐ in this doc have been reviewed. Others require your own review.

---

## Repo Structure Recommendation

```
legacyx-ai/                          # Your GitHub org repo
├── workspace/
│   ├── SOUL.md                      # ← Write this yourself
│   ├── HEARTBEAT.md                 # ← Write this yourself
│   ├── AGENTS.md                    # ← Write this yourself
│   └── AMBIENT_STATE.md             # ← Auto-generated by daemons
├── skills/
│   ├── legacy-eyes/
│   │   ├── SKILL.md                 # ← Write this yourself
│   │   └── legacy-eyes.py           # ← Write this yourself (SmolVLM)
│   └── legacy-pulse/
│       ├── SKILL.md                 # ← Write this yourself
│       ├── legacy-pulse.py          # ← Write this yourself (NATS bridge)
│       └── legacy-session-start.py # ← Write this yourself
├── launch.sh                        # ← Write this yourself
├── openclaw.json                    # ← Configured, not written from scratch
│
├── submodules/                      # External projects as git submodules
│   ├── NixClaw/                     # github.com/arniesaha/NixClaw (fork)
│   ├── rosclaw/                     # github.com/PlaiPin/rosclaw
│   └── OpenClawHomeAssistant/       # github.com/techartdev/OpenClawHomeAssistant
│
└── docs/
    ├── LEGACY_AMBIENT_PRESENCE_MASTERPLAN.md
    └── PLUGIN_CATALOG.md            # This file
```

Community skills install into `~/.openclaw/workspace/skills/` via `clawhub install` — they don't live in your repo. Only your custom skills and submodules go in the repo.

---

## Quick Reference — Source of Truth

| What You Need | Go Here |
|---|---|
| Browse all skills | `clawhub.com` |
| Curated 5,400 skill list | `github.com/VoltAgent/awesome-openclaw-skills` |
| Top 339 curated skills | `github.com/LeoYeAI/openclaw-master-skills` |
| Official docs | `docs.openclaw.ai` |
| Full plugin architecture | `deepwiki.com/openclaw/openclaw` |
| VisionClaw | `github.com/sseanliu/VisionClaw` |
| NixClaw | `github.com/arniesaha/NixClaw` |
| RosClaw | `github.com/PlaiPin/rosclaw` |
| HA Add-on | `github.com/techartdev/OpenClawHomeAssistant` |
| nix-openclaw | `github.com/openclaw/nix-openclaw` |
| n8n-claw | `github.com/freddy-schuetz/n8n-claw` |
| SecureClaw | Search GitHub: `SecureClaw Adversa AI` |

---

*Last updated: March 2026. Verify all repos are still active before cloning.*
*Security review required before enabling any community skill or plugin.*
