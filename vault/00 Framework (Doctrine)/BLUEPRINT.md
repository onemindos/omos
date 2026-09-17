# BLUEPRINT.md — One Mind OS Activation Plan

> **The complete build specification for taking Legacy from chatbot to ambient AI presence.**
> **Hardware**: MacBook Pro 16" M2 Max (32GB) · DigitalOcean VPS · Tailscale mesh
> **Date**: March 16, 2026

---

## Table of Contents

1. [Current State — What's Running Right Now](#1-current-state)
2. [Hardware Reality — MacBook Pro + Cloud Strategy](#2-hardware-reality)
3. [GPU Decision — Should Legacy Get GPU Power?](#3-gpu-decision)
4. [The Realtime Gap — Where Does Fast Processing Happen?](#4-the-realtime-gap)
5. [Office Project Comparison — Which One Do We Deploy?](#5-office-project-comparison)
6. [The Top 5 Picks — Detailed Activation Plans](#6-top-5-picks-activation)
7. [Beyond the Top 5 — Every Other Component](#7-beyond-the-top-5)
8. [Full Architecture — How Everything Connects](#8-full-architecture)
9. [Ambient Presence — How It Actually Works](#9-ambient-presence)
10. [Service Deployment Map — Every Process, Every Port](#10-service-deployment-map)
11. [NATS Connection — How to Get Live Events Flowing](#11-nats-connection)
12. [Activation Sequence — Phase by Phase (Detailed)](#12-activation-sequence)
13. [Cost Model — Monthly Burn](#13-cost-model)
14. [Skill Forge — Autonomous Skill Builder (MCP Builder Pattern)](#14-skill-forge)
15. [NATS Manager — Nervous System Control Skill](#15-nats-manager)
16. [Custom Domains — Caddy + Cloudflare Setup](#16-custom-domains)

---

## 1. Current State

### What's Live and Working

| Service | Where | Port | Status |
|---------|-------|------|--------|
| OpenClaw Gateway (9 agents) | VPS | 18789 | **LIVE** |
| NATS JetStream | VPS (Docker) | 4222 | **LIVE** |
| PostgreSQL/pgvector | VPS (Docker) | 5432 | **LIVE** |
| legacy-bridge (NATS↔OpenClaw) | VPS (Docker) | — | **LIVE** |
| nats-ingester | VPS (Docker) | — | **LIVE** |
| Discord (9 channels + voice) | VPS | — | **LIVE** |
| iMessage (BlueBubbles) | MacBook Pro | 1234 | **LIVE** |
| OpenClaw Studio | VPS | 3000 | **LIVE** |
| SwarmClaw v0.7.7 | VPS | 3456 | **LIVE** (outdated) |
| Codex sync (pull+watch) | VPS | — | **LIVE** |
| Memory search (hybrid) | VPS | — | **LIVE** |
| Datadog monitoring | VPS | 5000 | **LIVE** |

### What's NOT Working (and Should Be)

| Component | Problem | Fix |
|-----------|---------|-----|
| **Heartbeat** | Config is `{}` — empty. Zero proactive behavior. | Enable in openclaw.json |
| **Vision** | No legacy-eyes, no ambient sight | Deploy SmolVLM or API-based |
| **NATS events** | No legacy-pulse consuming events | Deploy as systemd service |
| **AMBIENT_STATE.md** | Template exists, never populated | Populated by eyes + pulse |
| **Lossless context** | Sliding-window compaction = memory loss | Install Lossless Claw plugin |
| **Office UI** | No visual agent status dashboard | Deploy office project |
| **Desktop presence** | No ambient overlay on desktop | Fork ClawX or deploy office |

---

## 2. Hardware Reality

### What We Actually Have

| Device | Specs | Tailscale IP | Role |
|--------|-------|-------------|------|
| **MacBook Pro 16"** | M2 Max · 32GB unified · macOS 26.3.1 | `100.114.142.55` | Dev machine, BlueBubbles, local vision |
| **VPS (onemind-main)** | DigitalOcean · Ubuntu 24.04 · CPU-only | `100.77.52.93` | OpenClaw runtime, NATS, databases |
| **VPS (onemind-sentinel)** | DigitalOcean · Ubuntu 24.04 | `100.67.223.5` | Secondary/standby |
| **iPhone** | iOS | `100.122.203.68` | ios-app client |

### Cloud-First Strategy (Correct for Now)

You're right — keep compute in the cloud. The MacBook Pro is a **development machine and local sensor node**, not a server. Here's why cloud-first makes sense:

1. **Uptime**: VPS runs 24/7. MacBook sleeps, travels, closes lid.
2. **Gateway stability**: OpenClaw, NATS, databases need persistent connectivity. Discord, iMessage bridges need always-on.
3. **NetworkEffects**: Tailscale mesh means all devices route through VPS as the hub.
4. **Migration path**: Start cloud → add local models incrementally → eventually hybrid.

### MacBook Pro's Role (Not a Server, But Powerful)

The M2 Max with 32GB is still critical:

| Task | Runs On MacBook | Why |
|------|-----------------|-----|
| SmolVLM-256M vision | Yes — 10-15fps, <1GB | Zero cost ambient vision from webcam |
| BlueBubbles iMessage | Yes — requires macOS | iMessage bridge to Legacy |
| Local Ollama models | Yes — if needed for offline | 32GB handles 30B parameter models |
| Desktop office/overlay | Yes — the screen you see | Ambient UI for command center |
| Development | Yes | VS Code, testing, Obsidian |

**Everything else stays on VPS or GPU Droplet.**

---

## 3. GPU Decision — Should Legacy Get GPU Power?

### The Short Answer

**Yes, but not the way you might think.** Legacy doesn't need a GPU for day-to-day operation. The OpenClaw Gateway, NATS, databases — all CPU. Where a GPU unlocks value is for:

1. **Local LLM inference** (run models yourself, stop paying API per-token)
2. **OpenViking context DB** (vector operations accelerated by GPU)
3. **Vision processing** (SmolVLM on GPU = faster, can run on VPS instead of MacBook)
4. **TTS** (Chatterbox self-hosted, GPU-accelerated = sub-100ms)

### DigitalOcean GPU Droplets — What's Available

DO has "Gradient AI GPU Droplets" in **NYC2, TOR1, ATL1, AMS3** datacenters:

| GPU | VRAM | Droplet RAM | vCPUs | Best For | Price Range |
|-----|------|-------------|-------|----------|-------------|
| **NVIDIA RTX 4000 Ada** | 20GB | 32GB | 8 | SmolVLM, Chatterbox, small models | ~$1.10-1.50/hr (~$800-1,100/mo) |
| **NVIDIA RTX 6000 Ada** | 48GB | 64GB | 8 | Ollama 70B models, OpenViking | ~$2.20-3.00/hr (~$1,600-2,200/mo) |
| **NVIDIA L40S** | 48GB | 64GB | 8 | Best inference value, all local models | ~$2.50-3.50/hr (~$1,800-2,500/mo) |
| **NVIDIA H100** | 80GB | 240GB | 20 | Overkill for us. Training tier. | ~$3.50-4.50/hr (~$2,500-3,300/mo) |
| **AMD MI300X** | 192GB | 240GB | 20 | Massive models, research | ~$4-6/hr |

### The Right GPU Strategy for One Mind OS

**Don't get a GPU Droplet running 24/7.** That's $800-3,000/mo for something you don't need full-time.

**Option A — Serverless Inference (Best Value Now)**
DO has "Gradient AI Platform" — serverless inference API. Pay per request, not per hour.
- Run Llama, Mistral, Qwen models via API
- No GPU Droplet needed
- Cost: pennies per request
- Like using OpenAI API but on DO infrastructure
- **Use your existing DO credit here**

**Option B — Spot/On-Demand GPU (When Needed)**
Spin up an RTX 4000 Ada GPU Droplet only when you need local inference:
- Run SmolVLM, Chatterbox TTS, OpenViking with GPU
- ~$1.10/hr — use for 4-8 hours, shut down
- Script to auto-start/stop based on schedule
- **$35-70/mo for 1-2 hours/day**

**Option C — Dedicated GPU (Month 2+, When Local Models Justify)**
When you're running Ollama 24/7 for all agent inference (replacing API calls):
- RTX 4000 Ada: ~$800/mo but saves $150-400/mo in API costs
- Break-even: only makes sense if API bill exceeds GPU cost
- **Not yet. API costs are ~$240-430/mo. GPU wouldn't save money yet.**

### Recommendation: Phase This

| Phase | GPU Strategy | Monthly Cost |
|-------|-------------|-------------|
| **Now (Phase 1-2)** | No GPU. Use API providers (Azure OpenAI, Anthropic). Use DO credit for Gradient Serverless for experiments. | $0 GPU, ~$300 API |
| **Phase 3 (Month 2)** | On-demand RTX 4000 Ada for SmolVLM + Chatterbox TTS. 2 hrs/day. | ~$65/mo GPU |
| **Phase 4 (Month 3+)** | Evaluate: if API bill > $600/mo, switch to dedicated RTX 6000 Ada running Ollama 24/7. Route all agent inference locally. | ~$2,200/mo GPU but $0 API |
| **Long-term** | Move to Mac Studio M4 Ultra at home (one-time $4k, runs 70B+ models 24/7 for free) | $0/mo |

### How GPU Would Work Inside DigitalOcean

```
Current Architecture:
┌─────────────────────────────┐
│  VPS (CPU Droplet)          │
│  OpenClaw → Azure/Anthropic │ ← all inference is remote API calls
│  NATS, Postgres, bridges    │
└─────────────────────────────┘

With GPU Droplet Added:
┌─────────────────────────────┐     ┌──────────────────────────────┐
│  VPS (CPU Droplet)          │     │  GPU Droplet (RTX 4000 Ada)  │
│  OpenClaw Gateway           │────→│  Ollama (Qwen3-30B, Llama)   │
│  NATS, Postgres, bridges    │     │  OpenViking server           │
│  Dashboard, Studio          │     │  SmolVLM vision              │
│                             │←────│  Chatterbox TTS              │
└─────────────────────────────┘     └──────────────────────────────┘
       Same VPC / Tailscale              Spun up on demand

OpenClaw model config change:
  "ollama/qwen3-30b": {
    "baseUrl": "http://gpu-droplet-tailscale-ip:11434",
    "alias": "local"
  }
```

Both droplets in the same DO region (NYC2) = internal network, ~1ms latency, no egress charges. Connect via Tailscale or DO VPC private networking.

---

## 4. The Realtime Gap — Where Does Fast Processing Happen?

### The Problem You Identified

A 5-minute heartbeat means Legacy checks the world every 5 minutes. That's not realtime. Between heartbeats, Legacy is deaf and blind. Here's the truth about timing:

### The 4 Speed Tiers of Legacy's Awareness

```
TIER 1: TRUE REALTIME (<1 second)
├── Discord voice — VAD detects speech in <100ms, responds in 1-3s
├── NATS event delivery — <10ms across Tailscale mesh
├── WebSocket messages — text chat arrives instantly
├── Gemini Live (future) — sub-200ms voice when glasses on
│
│   These work NOW. When Zeus talks, Legacy hears immediately.
│   The Gateway processes messages in real-time, no heartbeat needed.

TIER 2: NEAR-REALTIME (1-10 seconds)
├── SmolVLM local vision — scene digest every 2-5 seconds
├── NATS event buffer — legacy-pulse writes within 1 second of event
├── Tool execution — skill calls complete in <5 seconds typical
│
│   Between heartbeats, sensors are still writing to AMBIENT_STATE.md.

TIER 3: SCHEDULED (5 minutes)
├── Heartbeat tick — Legacy READS the accumulated state
├── Reads AMBIENT_STATE.md (vision + events since last check)
├── Decides: surface something? or HEARTBEAT_OK (silent)?
│
│   This is the GAP you're asking about. 5 minutes of accumulated
│   awareness read in one shot.

TIER 4: BATCH (hours/days)
├── Domain agent heartbeats (2-3 hours)
├── Codex sync (1 minute git pull, but contents change slowly)
├── Memory search reindexing
```

### Where Realtime ACTUALLY Happens

**Realtime isn't the heartbeat. Realtime is the event-driven layer.**

When Zeus sends a message (text or voice), the OpenClaw Gateway processes it **immediately** — no heartbeat involved. The heartbeat is for **proactive** behavior when nobody's talking to Legacy.

```
Zeus says something → Gateway routes immediately (realtime)
Zeus is silent → Heartbeat scans every 5 min (scheduled)
P0 alert fires → legacy-pulse writes to AMBIENT_STATE.md (near-realtime)
                 → Legacy sees it at next heartbeat (5 min worst case)
```

### How to Close the Gap

**Option 1 — Reduce Heartbeat to 1 Minute (Simple, Expensive)**
```json
{ "heartbeat": { "every": "1m" } }
```
- 1,440 LLM calls/day instead of 288
- Cost: ~$15-30/day with GPT-4o-mini
- P0 alert delivery: ≤60 seconds worst case
- **Good for Phase 1 testing, then tune**

**Option 2 — Webhook-Triggered Reads (Smart, Custom)**
Modify legacy-pulse to POST to OpenClaw's webhook endpoint when a P0 alert arrives:
```python
# In legacy-pulse.py, on P0 event:
if subject.startswith("system.alerts"):
    # Don't wait for heartbeat — poke Legacy NOW
    requests.post("http://localhost:18789/api/v1/hooks/alert",
                  json={"text": f"P0 ALERT: {summary}"},
                  headers={"Authorization": f"Bearer {HOOKS_TOKEN}"})
```
- P0 delivery: <5 seconds
- Regular events: still batched at heartbeat
- Cost: same as 5-min heartbeat (P0 alerts are rare)
- **Best of both worlds**

**Option 3 — OpenClaw Event Hooks (Built-in, Cleanest)**
OpenClaw has a hooks system. Configure a hook that fires when AMBIENT_STATE.md changes:
```json
{
  "hooks": {
    "fileChange": {
      "watch": ["AMBIENT_STATE.md"],
      "action": "agent.turn",
      "filter": "P0"
    }
  }
}
```
- Depends on OpenClaw version supporting this pattern
- Check docs for current hook capabilities

**Option 4 — Proactive Agent Skill (Already In Your Stack)**
You have `proactive-agent-3.1.0` skill installed. This skill enables event-driven agent turns — not just scheduled heartbeats. Configure it to fire on P0 NATS events.

### Recommended Realtime Strategy

| Event Type | Delivery Method | Latency |
|------------|----------------|---------|
| Zeus speaks (text/voice) | Direct Gateway message | <1 second |
| P0 system alert | Webhook from legacy-pulse | <5 seconds |
| Vision change (motion) | AMBIENT_STATE.md → next heartbeat | ≤5 minutes |
| NATS info events | Batched in AMBIENT_STATE.md | ≤5 minutes |
| Domain agent reports | Agent-to-agent or heartbeat | ≤3 hours |

**Start with 5-min heartbeat + P0 webhook. Tune from there.**

---

## 5. Office Project Comparison

### The 3 Projects Evaluated

| | **OpenClaw Office** | **Star Office UI** | **OpenClawfice** |
|---|---|---|---|
| **Repo** | WW-AI-Lab/openclaw-office | ringhyacinth/Star-Office-UI | openclawfice/openclawfice |
| **Stars** | 302 | 4,800 | 51 |
| **Forks** | 47 | 540 | 7 |
| **Contributors** | 2 (author + Cursor Agent) | 7 | 4 |
| **License** | MIT | MIT (code) / non-commercial (art) | AGPL-3.0 |
| **Language** | TypeScript 96% (React/Vite) | HTML 70% / Python 16% / JS 11% | TypeScript 78% (Next.js) |
| **Install** | `npx @ww-ai-lab/openclaw-office` | `git clone` + `python3 app.py` | `curl install.sh \| bash` |
| **Port** | 5180 | 19000 | 3333 |
| **Gateway Connection** | WebSocket direct + proxy | HTTP API (set_state.py) | Reads ~/.openclaw/openclaw.json |
| **Visual Style** | 2D SVG isometric + 3D R3F scene | Pixel art (Phaser-style) | Pixel art NPCs (Sims-style) |
| **Agent States** | idle/working/speaking/tool_calling/error | idle/writing/researching/executing/syncing/error | working/idle + room assignments |
| **Console/Admin** | Full: Dashboard, Agents, Channels, Skills, Cron, Settings | Minimal: state API + Gemini decor | Minimal: quest log, water cooler |
| **Chat** | Bottom-docked, streaming, Markdown | No embedded chat | Water cooler + DMs to NPCs |
| **Multi-Agent** | Yes — visualizes all Gateway agents | Yes — join keys for guests | Yes — auto-discovers agents |
| **i18n** | Chinese/English | Chinese/English/Japanese | English only |
| **Desktop Mode** | No (web only, could wrap in Electron) | Yes — Electron desktop pet (transparent) | No (Next.js web only) |
| **Mobile** | Responsive | Mobile-optimized | No |
| **Gamification** | None | None | XP, levels, trading cards, quest log, chiptune |
| **AI Features** | Token charts, cost pie charts, activity heatmaps, SubAgent graphs | AI-generated room decoration (Gemini) | Agent conversations (water cooler) |
| **Security** | Device auth bypass required | Session cookies, password protection | Token-based auth |

### The @bearclawbot Twitter Account

@bearclawbot is a "BearClaw CMO" account — a digital marketing agency bot. It references OpenClaw and agentic swarms but has **zero posts**. Not a project. Not relevant. Skip.

### Product Hunt — ClawOffice

The Product Hunt page returned a 403 error. Based on the openclawfice.com website, "ClawOffice" is the marketing name for the **OpenClawfice** project (openclawfice/openclawfice). Same thing.

### Winner: OpenClaw Office (WW-AI-Lab)

**For our use case — 9 production agents, need real monitoring, full console management — OpenClaw Office is the clear winner.** Here's why:

| Criteria | OpenClaw Office | Star Office UI | OpenClawfice |
|----------|----------------|----------------|--------------|
| **Connects to our Gateway** | Direct WebSocket to :18789 | HTTP state push (manual) | Reads config file |
| **Shows all 9 agents** | Auto-discovers via Gateway | Manual state per agent | Auto-discovers |
| **Console management** | Full: agents, channels, skills, cron, settings on every agent | None | None |
| **Realtime agent status** | WebSocket events (idle/working/speaking/tool_calling/error) | Polling-based state push | File-based polling |
| **Cost/token monitoring** | Token line charts, cost pie charts, activity heatmaps | None | None |
| **Tech stack alignment** | React 19 + Vite + TypeScript = matches our architecture | Flask (Python) = different stack | Next.js = close but AGPL |
| **Production readiness** | 6 releases, last updated last week | 2 releases, actively developed | No releases, v0.1.1 |
| **License** | MIT (full freedom) | MIT code / non-commercial art | AGPL-3.0 (viral copyleft) |

**Star Office UI** is charming and great for showing off (pixel art, desktop pet) but it's a **status board**, not a management tool. It doesn't connect to our Gateway via WebSocket — you manually push states via a Python script. No console, no skill management, no chat integration.

**OpenClawfice** has great gamification (XP, trading cards, quest log) but is AGPL-3.0 (restrictive), only 51 stars, no releases, and lacks the management console we need.

**OpenClaw Office** gives us: real-time WebSocket connection to all 9 agents, full management console (agents, channels, skills, cron, settings), token/cost monitoring, 2D AND 3D views, embedded chat, streaming tool calls — and it's MIT licensed.

### Deployment Plan for OpenClaw Office

```bash
# On VPS — one command
ssh root@104.248.4.212
su - onemind
npm install -g @ww-ai-lab/openclaw-office
openclaw-office --gateway ws://localhost:18789 --port 5180 --token "$(openclaw config get gateway.auth.token)"
```

Or as a systemd service:
```ini
# /etc/systemd/system/openclaw-office.service
[Unit]
Description=OpenClaw Office — Agent Visualization
After=network.target openclaw.service

[Service]
Type=simple
User=onemind
ExecStart=/usr/local/bin/openclaw-office --gateway ws://localhost:18789 --port 5180
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Tailscale HTTPS proxy for remote access:
```
https://onemind-main.tail717ba0.ts.net:5180 → OpenClaw Office
```

### Star Office UI — Deploy As Fun Bonus

Star Office UI is too cool NOT to deploy. Run it alongside OpenClaw Office as the "vibes" dashboard:

```bash
cd /home/onemind/onemind-os/command-center/
git clone https://github.com/ringhyacinth/Star-Office-UI.git star-office
cd star-office
python3 -m pip install -r backend/requirements.txt
cp state.sample.json state.json
cd backend && python3 app.py &
```

Add to Legacy's SOUL.md:
```markdown
## Star Office Status Sync
- When starting a task: `python3 set_state.py writing "Task description"`
- When researching: `python3 set_state.py researching "What I'm looking into"`
- When idle: `python3 set_state.py idle "Ready"`
```

Port 19000 — the pixel art view of your agents. Pure vibes.

---

## 6. Top 5 Picks — Detailed Activation Plans

### Pick 1: Lossless Claw (P0 — Install Day 1)

**What**: OpenClaw plugin. Replaces lossy sliding-window compaction with DAG-based summaries. Every message persisted in SQLite. Nothing ever lost.

**Why it's #1**: Without this, Legacy forgets every conversation. With it, Legacy can `lcm_grep("investor deck")` and find messages from 3 weeks ago.

**Activation**:
```bash
ssh root@104.248.4.212
su - onemind

# Install plugin
openclaw plugins install @martian-engineering/lossless-claw

# Edit openclaw.json — add to plugins section:
# {
#   "slots": { "contextEngine": "lossless-claw" },
#   "entries": {
#     "lossless-claw": {
#       "enabled": true,
#       "config": {
#         "freshTailCount": 32,
#         "contextThreshold": 0.75,
#         "incrementalMaxDepth": -1
#       }
#     }
#   }
# }
#
# Also add session reset:
# "session": { "reset": { "mode": "idle", "idleMinutes": 10080 } }

systemctl restart openclaw
journalctl -u openclaw -f  # verify LCM loads
```

**Gives Legacy**: `lcm_grep`, `lcm_expand`, `lcm_describe` tools. Persistent conversation memory. 7-day session window.

---

### Pick 2: OpenViking Context DB (P0 — Deploy Week 1)

**What**: Context Database for AI agents. Filesystem paradigm (`viking://`). L0/L1/L2 tiered loading. +46% accuracy, -91% tokens vs native OpenClaw memory.

**Why it's #2**: Cuts Legacy's token costs dramatically while improving retrieval accuracy. The codex (vault) becomes a structured knowledge hierarchy instead of flat vector search.

**Activation**:
```bash
ssh root@104.248.4.212
su - onemind

# Install
pip install openviking --upgrade --force-reinstall

# Configure
mkdir -p ~/.openviking
cat > ~/.openviking/ov.conf << 'EOF'
[server]
host = 127.0.0.1
port = 8900

[embedding]
provider = openai
model = text-embedding-3-small

[indexing]
paths = /home/onemind/onemind-os/codex
watch = true
tier_loading = true
EOF

# Create systemd service
cat > /etc/systemd/system/openviking.service << 'EOF'
[Unit]
Description=OpenViking Context Database
After=network.target

[Service]
Type=simple
User=onemind
ExecStart=/usr/local/bin/openviking-server --config /home/onemind/.openviking/ov.conf
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl enable --now openviking

# Add OpenViking plugin to openclaw.json
# {
#   "openviking-memory": {
#     "enabled": true,
#     "config": {
#       "serverUrl": "http://127.0.0.1:8900",
#       "tierLoading": true
#     }
#   }
# }

systemctl restart openclaw
```

**Gives Legacy**: 52% accuracy on knowledge retrieval (up from 35.65%), ~91% fewer tokens consumed on memory ops.

---

### Pick 3: Obsidian Skills by kepano (P1 — Install Week 1)

**What**: 5 official skills by Obsidian's creator for proper vault file handling. Teaches Legacy correct Obsidian markdown syntax, wikilinks, callouts, properties, Bases, Canvas, CLI, and Defuddle (clean web scraping).

**Why**: Legacy already writes to the codex but doesn't use proper Obsidian syntax. These skills ensure correct formatting.

**Activation**:
```bash
ssh root@104.248.4.212
su - onemind

# Clone and copy skills
cd /tmp
git clone https://github.com/kepano/obsidian-skills.git
cp -r obsidian-skills/skills/* /home/onemind/onemind-os/openclaw/workspaces/legacy/skills/
rm -rf /tmp/obsidian-skills

systemctl restart openclaw
```

**Gives Legacy**: `obsidian-markdown`, `obsidian-bases`, `json-canvas`, `obsidian-cli`, `defuddle` skills. Better vault writing. Token-efficient web scraping.

---

### Pick 4: ClawX Desktop Fork (P1 — Week 2-3)

**What**: Electron + React desktop app for OpenClaw. Already handles gateway connection, chat, skills, cron. We fork it and add MediaPipe + SmolVLM + ambient overlay.

**Why**: Instead of building a desktop command center from scratch, fork the 4.5k-star foundation.

**Activation**:
```bash
# On MacBook Pro
cd ~/Documents/onemindos-v0.5/command-center/
git clone https://github.com/ValueCell-ai/ClawX.git onemind-desktop
cd onemind-desktop
npm install

# Configure gateway connection
# Edit config to point to: onemind-main.tail717ba0.ts.net:18789
# Set gateway token

npm run dev  # test connection

# Phase 2: Add custom panels
# - MediaPipe Hands gesture recognition
# - SmolVLM camera panel
# - NATS event subscriber
# - Ambient overlay mode (transparent always-on-top)
```

**Gives You**: Desktop app with chat to all 9 agents, ambient overlay, gesture control (future).

---

### Pick 5: OpenClaw Office (P1 — Deploy Week 1)

**What**: Full visual monitoring + management frontend. 2D/3D office views, real-time agent status, embedded chat, console with agents/channels/skills/cron/settings management.

**Why**: Replaces the need for SwarmClaw as a management interface for OpenClaw agents. Better integrated. Designed specifically for OpenClaw Gateway.

**Activation**: (see Section 5 for full commands)

```bash
ssh root@104.248.4.212
su - onemind
npm install -g @ww-ai-lab/openclaw-office

# Create systemd service (see above)
systemctl enable --now openclaw-office
```

**Gives You**: Visual dashboard of all 9 agents, real-time status, token costs, management console, embedded chat — all via web.

---

## 7. Beyond the Top 5 — Every Other Component

### Heartbeat Activation (Day 1 — same as Lossless Claw)

```json
// Add to openclaw.json → agents.defaults.heartbeat
{
  "every": "5m",
  "target": "last",
  "directPolicy": "allow",
  "lightContext": false,
  "model": "azure-openai/gpt-4o-mini",
  "ackMaxChars": 500,
  "activeHours": {
    "start": "06:00",
    "end": "23:59",
    "timezone": "America/New_York"
  }
}

// Legacy agent override:
{
  "id": "legacy",
  "heartbeat": {
    "every": "5m",
    "model": "azure-openai/gpt-4o",
    "prompt": "Read HEARTBEAT.md and AMBIENT_STATE.md. Surface anything actionable. Otherwise reply HEARTBEAT_OK."
  }
}
```

### legacy-pulse (NATS → AMBIENT_STATE.md) — Week 1

Already written at `openclaw/workspaces/legacy/skills/legacy-pulse/legacy-pulse.py`.

```bash
ssh root@104.248.4.212

# Install dependencies
pip install nats-py

# Deploy service
cat > /etc/systemd/system/legacy-pulse.service << 'EOF'
[Unit]
Description=Legacy Pulse — NATS Event Bridge
After=network.target docker.service

[Service]
Type=simple
User=onemind
WorkingDirectory=/home/onemind/onemind-os/openclaw/workspaces/legacy/skills/legacy-pulse
ExecStart=/usr/bin/python3 legacy-pulse.py
Restart=always
RestartSec=5
Environment=NATS_URL=nats://localhost:4222
Environment=OPENCLAW_WORKSPACE=/home/onemind/.openclaw/workspace

[Install]
WantedBy=multi-user.target
EOF

systemctl enable --now legacy-pulse
```

### legacy-eyes (Vision → AMBIENT_STATE.md) — Week 2

SmolVLM on MacBook Pro (not Mac Mini — correcting all previous references):

```bash
# On MacBook Pro
pip install transformers torch accelerate pillow opencv-python

# SmolVLM-256M auto-downloads (~500MB)
cd ~/Documents/onemindos-v0.5/openclaw/workspaces/legacy/skills/legacy-eyes/

# Modify legacy-eyes.py to:
# 1. Use SmolVLM instead of Claude API
# 2. Set DIGEST_INTERVAL to 2 (every 2 seconds)
# 3. Output via NATS publish (legacy-pulse writes to AMBIENT_STATE.md)
#    OR direct SSH write to VPS workspace

python3 legacy-eyes.py &
```

### SwarmClaw Upgrade (v0.7.7 → v1.1.5) — Week 1

```bash
ssh root@104.248.4.212
cd /home/onemind/onemind-os/command-center/swarmclaw
npm update  # or npm install @swarmclawai/swarmclaw@latest
systemctl restart swarmclaw
```

New capabilities: Missions, Orchestrator agents, Structured Sessions, durable follow-ups.

### Skill Archival (12 Unused Skills) — Week 1

```bash
ssh root@104.248.4.212
su - onemind
cd /home/onemind/onemind-os/openclaw/workspaces/legacy/skills
mkdir -p archive
mv clickup-1.0.0 homestead-control desktop-control-1.0.0 browser-1.0.0 \
   credit-repair-1.0.0 lob-1.0.0 uber-1.0.0 transit-1.0.0 \
   zoho-mail-1.0.0 todoist-api find-skills-0.1.0 skill-creator-0.1.0 \
   archive/
systemctl restart openclaw
```

### ClawVault (Structured Memory) — Evaluate Week 3

After Lossless Claw + OpenViking are running, evaluate whether ClawVault's entity graph adds value:
```bash
openclaw hooks install clawvault
```

Unique value: multi-hop queries ("who did Zeus meet who works at X?"), cross-encoder reranking, Obsidian canvas dashboards, background workers (Curator, Janitor).

### Star Office UI (Fun Dashboard) — Week 2+

```bash
ssh root@104.248.4.212
su - onemind
cd /home/onemind/onemind-os/command-center/
git clone https://github.com/ringhyacinth/Star-Office-UI.git star-office
cd star-office
python3 -m pip install -r backend/requirements.txt
cp state.sample.json state.json

# Create systemd service for it
cat > /etc/systemd/system/star-office.service << 'EOF'
[Unit]
Description=Star Office UI — Pixel Art Agent Dashboard
After=network.target

[Service]
Type=simple
User=onemind
WorkingDirectory=/home/onemind/onemind-os/command-center/star-office/backend
ExecStart=/usr/bin/python3 app.py
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl enable --now star-office
# Accessible at http://localhost:19000
```

### NixClaw Fork (Gemini Live Bridge) — Month 2

```bash
# On dev machine
git clone https://github.com/arniesaha/NixClaw.git ~/legacy-claw
cd ~/legacy-claw

# Modify:
# 1. Agent name → "Legacy"
# 2. DEFAULT_OPENCLAW_HOSTNAME → onemind-main.tail717ba0.ts.net
# 3. DEFAULT_OPENCLAW_TOKEN → gateway token
# 4. GEMINI_API_KEY → your key
# 5. Add NATS publish on session start

# Build for Mason Android device (or iPhone)
```

### Chatterbox TTS (Self-Hosted Voice) — Month 2

```bash
# On GPU Droplet (or MacBook Pro)
pip install chatterbox-tts
chatterbox-server --host 0.0.0.0 --port 8200

# Then update OpenClaw voice config to use custom TTS endpoint
```

Replaces OpenAI TTS. Sub-200ms on CPU. Zero cloud dependency.

---

## 8. Full Architecture — How Everything Connects

```
╔════════════════════════════════════════════════════════════════════════════════╗
║                         ONE MIND OS — FULL ARCHITECTURE                       ║
║                         March 2026 · Blueprint v1.0                           ║
╚════════════════════════════════════════════════════════════════════════════════╝

                    ┌─────────────────────────────────┐
                    │        ZEUS (Human)              │
                    │                                  │
                    │  Voice → Discord / NixClaw       │
                    │  Text  → Discord / iMessage      │
                    │  Vision → Ray-Ban / Webcam       │
                    │  Gestures → MediaPipe (desktop)  │
                    │  Eyes  → OpenClaw Office UI      │
                    └──────────────┬───────────────────┘
                                   │
            ┌──────────────────────┼──────────────────────────┐
            │                      │                          │
            ▼                      ▼                          ▼
┌───────────────────┐  ┌──────────────────────┐  ┌────────────────────┐
│  MacBook Pro M2   │  │   VPS (CPU Droplet)  │  │  GPU Droplet       │
│  Max 32GB         │  │   104.248.4.212      │  │  (on-demand)       │
│                   │  │                      │  │                    │
│  ┌─────────────┐  │  │  ┌────────────────┐  │  │  ┌──────────────┐ │
│  │ SmolVLM-256M│  │  │  │ OpenClaw GW    │  │  │  │ Ollama       │ │
│  │ (vision)    │  │  │  │ (9 agents)     │  │  │  │ (local LLM)  │ │
│  └──────┬──────┘  │  │  │ :18789         │  │  │  │ :11434       │ │
│         │         │  │  └───────┬────────┘  │  │  └──────────────┘ │
│  ┌──────▼──────┐  │  │          │           │  │  ┌──────────────┐ │
│  │ BlueBubbles │  │  │  ┌───────▼────────┐  │  │  │ OpenViking   │ │
│  │ iMessage    │  │  │  │ Plugins:       │  │  │  │ (context DB) │ │
│  │ :1234       │  │  │  │ lossless-claw  │  │  │  │ :8900        │ │
│  └─────────────┘  │  │  │ openviking-mem │  │  │  └──────────────┘ │
│                   │  │  │ discord        │  │  │  ┌──────────────┐ │
│  ┌─────────────┐  │  │  │ bluebubbles    │  │  │  │ Chatterbox   │ │
│  │ onemind-    │  │  │  │ lobster        │  │  │  │ TTS :8200    │ │
│  │ desktop     │  │  │  └────────────────┘  │  │  └──────────────┘ │
│  │ (ClawX fork)│  │  │                      │  │  ┌──────────────┐ │
│  └─────────────┘  │  │  ┌────────────────┐  │  │  │ SmolVLM      │ │
│                   │  │  │ 29 Skills      │  │  │  │ (alt vision) │ │
│  ┌─────────────┐  │  │  │ (tools)        │  │  │  └──────────────┘ │
│  │ Obsidian    │  │  │  └────────────────┘  │  └────────────────────┘
│  │ (codex)     │  │  │                      │
│  └─────────────┘  │  │  ┌────────────────┐  │
│                   │  │  │ legacy-pulse    │  │
└───────────────────┘  │  │ (NATS→STATE)   │  │
                       │  └────────────────┘  │
                       │                      │
                       │  ┌────────────────┐  │
                       │  │ NATS JetStream │  │
                       │  │ :4222          │  │
                       │  └────────────────┘  │
                       │                      │
                       │  ┌────────────────┐  │
                       │  │ PostgreSQL     │  │
                       │  │ pgvector :5432 │  │
                       │  └────────────────┘  │
                       │                      │
                       │  ── DASHBOARDS ──    │
                       │  ┌────────────────┐  │
                       │  │ OpenClaw Office│  │
                       │  │ :5180          │  │
                       │  ├────────────────┤  │
                       │  │ Star Office UI │  │
                       │  │ :19000         │  │
                       │  ├────────────────┤  │
                       │  │ SwarmClaw      │  │
                       │  │ :3456          │  │
                       │  ├────────────────┤  │
                       │  │ OC Studio      │  │
                       │  │ :3000          │  │
                       │  └────────────────┘  │
                       │                      │
                       │  ── BRIDGES ──       │
                       │  legacy-bridge       │
                       │  nats-ingester       │
                       │  codex-pull/watch    │
                       │  openclaw-reload     │
                       │  datadog             │
                       └──────────────────────┘
```

---

## 9. Ambient Presence — How It Actually Works

### The AMBIENT_STATE.md Loop

This is the core mechanism. Every ambient service writes to one file. Legacy reads it on heartbeat.

```
┌──────────────────────────────────────────────────────────────────────┐
│                                                                      │
│                      AMBIENT_STATE.md                                 │
│                    (the single source of NOW)                         │
│                                                                      │
│  Updated: 2026-03-16 14:22:01                                       │
│                                                                      │
│  ## Vision Digest                                                    │
│  Zeus at desk, focused on dual monitors. Left: VS Code.             │
│  Right: Discord. Coffee cup present. Scene unchanged 3 ticks.        │
│                                                                      │
│  ## Active Session                                                   │
│  - Session ID: a3f7b2c1 | Location: home | OMCS: IT.31.P01         │
│                                                                      │
│  ## Recent NATS Events (last 20)                                    │
│  - system.alerts.cpu — VPS 34% @ 14:20:10                          │
│  - omcs.ge.todoist — Task synced: investor deck @ 14:19:40         │
│  - legacy.session.start — Session @ home @ 14:15:00                │
│                                                                      │
└──────────────────────┬──────────────────────────┬────────────────────┘
                       │                          │
              WRITTEN BY:                    READ BY:
              │                              │
              ├── legacy-eyes.py             ├── OpenClaw heartbeat (every 5 min)
              │   SmolVLM on MacBook Pro     │   Legacy reads → decides → surfaces/silent
              │   OR SmolVLM on GPU Droplet  │
              │   Writes vision digest       ├── OpenClaw Office (real-time display)
              │                              │
              ├── legacy-pulse.py            ├── onemind-desktop (ambient overlay)
              │   NATS subscriber on VPS     │
              │   Writes event buffer        └── P0 webhook (instant alerts)
              │
              └── legacy-session-start.py
                  Fires on Discord open / app launch
                  Writes session metadata
```

### A Day in the Life of Ambient Legacy

**6:00 AM** — Legacy activates (activeHours start). First heartbeat.
```
HEARTBEAT.md → check AMBIENT_STATE.md
Vision: "Office empty, lights off"
Events: nothing overnight
→ HEARTBEAT_OK (silent)
```

**8:15 AM** — Zeus opens Discord on MacBook Pro.
```
legacy-session-start fires → NATS: legacy.session.start
legacy-pulse writes: "Session started @ home @ 08:15"
SmolVLM starts seeing: "Zeus at desk, opening laptop"
```

**8:20 AM** — Next heartbeat.
```
Legacy reads AMBIENT_STATE.md:
  Vision: "Zeus at desk, opening laptop, coffee present"
  Events: "Session started @ home"
→ Legacy to Discord: "Morning, Zeus. No overnight alerts.
   3 Todoist tasks due today. VPS healthy at 34% CPU."
```

**10:45 AM** — Zeus is deep in work. SmolVLM sees "focused, coding."
```
Heartbeat after heartbeat: HEARTBEAT_OK (silent — nothing to surface)
SmolVLM keeps writing "scene unchanged" counts
```

**10:47 AM** — NATS: `system.alerts.disk` — VPS disk at 92%.
```
legacy-pulse writes P0 event to AMBIENT_STATE.md
legacy-pulse ALSO fires webhook to OpenClaw: POST /api/v1/hooks/alert
→ Legacy IMMEDIATELY (not waiting for heartbeat):
  "⚠️ VPS disk at 92%. Docker images likely. Want me to prune?"
```

**2:00 PM** — Zeus puts on Ray-Ban glasses, leaves desk.
```
SmolVLM: "Zeus left desk" → stops updating
NixClaw on Mason activates → Gemini Live sees + hears
Zeus speaks: "Legacy, what's my schedule?"
Gemini: execute("check afternoon schedule")
→ POST to OpenClaw Gateway
→ Legacy checks Zoho Calendar → returns result
→ Gemini speaks to glasses: "3pm contractor call, 5pm team review"
```

**11:59 PM** — activeHours end. Legacy sleeps until 6 AM.

---

## 10. Service Deployment Map

### Every Process After Full Activation

| Service | Host | Port | Protocol | Deployed | Phase |
|---------|------|------|----------|----------|-------|
| **openclaw.service** | VPS | 18789 | WebSocket | **Running** | 0 |
| **onemind-nats** | VPS | 4222, 1883, 7422, 8222 | Docker | **Running** | 0 |
| **onemind-omnidb** | VPS | 5432 | Docker | **Running** | 0 |
| **onemind-legacy-bridge** | VPS | — | Docker | **Running** | 0 |
| **onemind-nats-ingester** | VPS | — | Docker | **Running** | 0 |
| **openclaw-studio.service** | VPS | 3000 | HTTP | **Running** | 0 |
| **swarmclaw.service** | VPS | 3456 | HTTP | **Running** (upgrade) | 1 |
| **codex-pull.timer** | VPS | — | Git | **Running** | 0 |
| **codex-watch.service** | VPS | — | Git | **Running** | 0 |
| **openclaw-reload.service** | VPS | — | inotify | **Running** | 0 |
| **datadog-agent.service** | VPS | 5000 | HTTP | **Running** | 0 |
| lossless-claw (plugin) | VPS | — | In-process | **Deploy** | 1 |
| heartbeat (config) | VPS | — | Scheduler | **Enable** | 1 |
| **openclaw-office.service** | VPS | 5180 | HTTP+WS | **Deploy** | 1 |
| **openviking.service** | VPS | 8900 | HTTP | **Deploy** | 2 |
| **legacy-pulse.service** | VPS | — | NATS sub | **Deploy** | 2 |
| **star-office.service** | VPS | 19000 | HTTP | **Deploy** | 2 |
| legacy-eyes.py | MacBook Pro | — | OpenCV+ML | **Deploy** | 2 |
| BlueBubbles | MacBook Pro | 1234 | HTTP | **Running** | 0 |
| onemind-desktop | MacBook Pro | — | Electron | **Build** | 3 |
| NixClaw fork | Mason Android | — | WS+Gemini | **Build** | 4 |
| Chatterbox TTS | GPU/VPS | 8200 | HTTP | **Deploy** | 4 |
| Ollama | GPU Droplet | 11434 | HTTP | **Deploy** | 4+ |

### Tailscale HTTPS Endpoints (External Access)

| URL | Port | Service |
|-----|------|---------|
| `https://onemind-main.tail717ba0.ts.net` | 443→18789 | OpenClaw Gateway |
| `https://onemind-main.tail717ba0.ts.net:3443` | 3443→3000 | OpenClaw Studio |
| `https://onemind-main.tail717ba0.ts.net:8443` | 8443→3456 | SwarmClaw |
| `https://onemind-main.tail717ba0.ts.net:5180` | 5180 | OpenClaw Office (NEW) |
| `https://onemind-main.tail717ba0.ts.net:19000` | 19000 | Star Office UI (NEW) |

---

## 11. NATS Connection — How to Get Live Events Flowing

### The Reality: No Native "Bring Your Own NATS" in OpenClaw

OpenClaw Gateway has **no** `natsUrl`, `eventBus`, or BYON config key. The gateway speaks HTTP/WebSocket, not NATS natively. NATS integration works through two components you already have:

| Component | Direction | How It Works | Status |
|-----------|-----------|--------------|--------|
| **bridge.py** | NATS → OpenClaw | Python daemon subscribes to 11 NATS patterns, POSTs to `/hooks/wake` and `/hooks/agent` | Running in Docker as `onemind-legacy-bridge` |
| **nexus-bridge skill** | OpenClaw → NATS | Node.js skill inside gateway. Legacy calls `nexus_dispatch`, `nexus_query`, `nexus_publish` tools | Installed in skills/ |

**The problem**: NATS is live but **silent**. No event producers are publishing. Both the bridge and nexus-bridge have nothing to process.

### Step-by-Step: Make NATS Flow Daily

#### A. Install NATS CLI on VPS

```bash
ssh root@104.248.4.212
curl -sf https://binaries.nats.dev/nats-io/natscli/nats@latest | sh
mv nats /usr/local/bin/
nats --version  # verify
```

#### B. Create JetStream Streams (Events Persist Instead of Evaporating)

Without streams, a NATS event that fires when nobody's listening is gone forever.

```bash
# Create 6 persistent streams
nats stream add EVENTS   --subjects "events.>"   --retention limits --max-msgs 100000 --max-age 30d --storage file --replicas 1 --discard old
nats stream add SYSTEM   --subjects "system.>"   --retention limits --max-msgs 50000  --max-age 7d  --storage file --replicas 1 --discard old
nats stream add NEXUS    --subjects "nexus.>"    --retention limits --max-msgs 50000  --max-age 14d --storage file --replicas 1 --discard old
nats stream add DIGITAL  --subjects "digital.>"  --retention limits --max-msgs 50000  --max-age 14d --storage file --replicas 1 --discard old
nats stream add PHYSICAL --subjects "physical.>" --retention limits --max-msgs 100000 --max-age 30d --storage file --replicas 1 --discard old
nats stream add COMMS    --subjects "comms.>"    --retention limits --max-msgs 50000  --max-age 7d  --storage file --replicas 1 --discard old

nats stream ls  # verify all 6 exist
```

#### C. Create VPS Health Publisher (First Event Producer)

```bash
cat > /home/onemind/nats-health-pub.sh << 'SCRIPT'
#!/bin/bash
CPU=$(top -bn1 | grep "Cpu(s)" | awk '{print $2}')
MEM=$(free -m | awk 'NR==2{printf "%.1f", $3*100/$2}')
DISK=$(df -h / | awk 'NR==2{print $5}' | tr -d '%')
DOCKER=$(docker ps --format '{{.Names}}' | tr '\n' ',' | sed 's/,$//')

nats pub system.health.vps "{
  \"cpu_pct\": $CPU,
  \"mem_pct\": $MEM,
  \"disk_pct\": $DISK,
  \"docker_running\": \"$DOCKER\",
  \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
}" 2>/dev/null
SCRIPT
chmod +x /home/onemind/nats-health-pub.sh

# Run every 5 minutes via cron
(crontab -u onemind -l 2>/dev/null; echo "*/5 * * * * /home/onemind/nats-health-pub.sh") | crontab -u onemind -
```

#### D. Add Codex Sync Events (codex-pull publishes to NATS)

```bash
# Add NATS publish to codex-pull.service post-pull hook
cat > /home/onemind/codex-nats-notify.sh << 'SCRIPT'
#!/bin/bash
cd /home/onemind/onemind-os
CHANGES=$(git diff --stat HEAD~1 -- codex/ 2>/dev/null | tail -1)
if [ -n "$CHANGES" ]; then
  nats pub events.codex.updated "{
    \"changes\": \"$CHANGES\",
    \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
  }" 2>/dev/null
fi
SCRIPT
chmod +x /home/onemind/codex-nats-notify.sh
```

#### E. Make Legacy Publish Heartbeat to NATS

In Legacy's HEARTBEAT.md, add:
```markdown
## After completing heartbeat checks:
- Use nexus_publish: subject "nexus.heartbeat.legacy", payload { status: "ok", agents_checked: 9 }
```

Now every 30-minute heartbeat also publishes to NATS. Domain agents do the same.

#### F. Verify Events Are Flowing

```bash
# Watch ALL events in real time
nats sub ">" --count 10  # see next 10 events

# Check stream contents
nats stream info SYSTEM
nats stream info NEXUS

# Check specific subjects
nats sub "system.health.vps" --count 1  # wait for next health event
nats sub "nexus.heartbeat.>" --count 1  # wait for next heartbeat
```

### Daily Event Volume After Setup

| Producer | Subject | Frequency | Events/Day |
|----------|---------|-----------|------------|
| VPS health cron | `system.health.vps` | Every 5 min | 288 |
| Legacy heartbeat | `nexus.heartbeat.legacy` | Every 30 min | ~32 |
| 8 domain agent heartbeats | `nexus.heartbeat.{id}` | Every 1-2h | ~64 |
| Codex sync | `events.codex.updated` | On git pull (variable) | ~10-50 |
| Discord messages | `comms.text.human.zeus` | On message | ~20-100 |
| bridge.py forwarding | Various → `/hooks/` | On NATS event | same as above |
| **Total daily minimum** | | | **~400-500+** |

### Complete NATS Data Flow

```
OUTBOUND (Legacy → NATS):
  Legacy calls nexus_dispatch / nexus_query / nexus_publish
    │  (Node.js skill, runs inside gateway process)
    ▼
  NATS JetStream (:4222) → persisted in 6 streams
    ├──→ Agno agents (subscribed to digital.agent.{id}.task)
    ├──→ OM1 robots (subscribed to physical.robot.{id}.command)
    ├──→ nats-ingester (writes to PostgreSQL)
    ├──→ legacy-pulse (writes to AMBIENT_STATE.md)
    └──→ Any NATS subscriber on Tailscale mesh

INBOUND (NATS → Legacy):
  Something publishes to NATS (cron, sensor, agent, webhook)
    │
    ▼
  bridge.py (Docker: onemind-legacy-bridge)
    │  Matches 11 subject patterns
    │  POSTs to OpenClaw /hooks/wake or /hooks/agent
    ▼
  OpenClaw Gateway processes event
    │  /hooks/wake → adds to Legacy's context
    │  /hooks/agent → routes to specific agent, optionally delivers to Zeus
    ▼
  Legacy sees event and decides what to do

AMBIENT PATH (NATS → File → Heartbeat):
  NATS events → legacy-pulse.py → writes AMBIENT_STATE.md
    │
    ▼ (every 30 min)
  Legacy heartbeat reads AMBIENT_STATE.md
    │  Decides: surface something? Or HEARTBEAT_OK?
```

---

## 12. Activation Sequence — Phase by Phase (Detailed)

### Phase 0: VERIFY (Today — 30 minutes)

**Goal**: Confirm what's actually running vs broken.

```bash
ssh root@104.248.4.212

# 1. NATS alive?
docker exec onemind-nats nats-server --version
docker logs onemind-nats --tail 20

# 2. legacy-bridge running?
docker ps | grep legacy-bridge
docker logs onemind-legacy-bridge --tail 30

# 3. All 9 agents loaded?
curl -s -H "Authorization: Bearer $OPENCLAW_GATEWAY_TOKEN" \
  http://localhost:18789/api/v1/agents | python3 -c "import sys,json; [print(a['id']) for a in json.load(sys.stdin)]"

# 4. Test NATS pub/sub
nats pub test.hello "Hello from NATS"

# 5. nexus-bridge skill loaded?
journalctl -u openclaw --since "1 hour ago" | grep -i "nexus-bridge"

# 6. Current heartbeat config?
grep -A10 '"heartbeat"' /home/onemind/.openclaw/openclaw.json | head -15

# 7. NATS_URL set?
grep NATS_URL /opt/onemind/.env.sa
```

### Phase 1: FOUNDATION (This Week — Day 1-3)

```bash
# Step 1.1: Install NATS CLI + Create JetStream Streams
curl -sf https://binaries.nats.dev/nats-io/natscli/nats@latest | sh && mv nats /usr/local/bin/
# Create 6 streams (see Section 11B above)

# Step 1.2: Install Lossless Claw
su - onemind -c 'openclaw plugins install @martian-engineering/lossless-claw'
# Add plugin config to openclaw.json (see Section 6, Pick 1)

# Step 1.3: Enable Heartbeat
# Verify/update heartbeat config in /home/onemind/.openclaw/openclaw.json
# Target: "every": "30m", "target": "discord", active 07:00-23:00

# Step 1.4: VPS Health → NATS publisher (see Section 11C above)
# Creates system.health.vps events every 5 min

# Step 1.5: Install Obsidian Skills
cd /tmp && git clone https://github.com/kepano/obsidian-skills.git
cp -r obsidian-skills/skills/* /home/onemind/onemind-os/openclaw/workspaces/legacy/skills/
rm -rf /tmp/obsidian-skills

# Step 1.6: Deploy OpenClaw Office
npm install -g @ww-ai-lab/openclaw-office
# Create systemd service (see Section 5)
systemctl enable --now openclaw-office

# Step 1.7: Upgrade SwarmClaw
cd /home/onemind/onemind-os/command-center/swarmclaw && npm update

# Step 1.8: Restart + Verify
systemctl restart openclaw
journalctl -u openclaw -f  # watch for heartbeat, lossless-claw, nexus-bridge
nats sub ">" --count 5     # verify events flowing
```

### Phase 2: AMBIENT AWARENESS (Week 2)

```bash
# Step 2.1: Deploy legacy-pulse (NATS → AMBIENT_STATE.md writer)
pip install nats-py
# Create systemd service (see Section 7)
systemctl enable --now legacy-pulse

# Step 2.2: Deploy OpenViking Context DB
pip install openviking --upgrade
# Create config + systemd service (see Section 6, Pick 2)
systemctl enable --now openviking

# Step 2.3: Update HEARTBEAT.md to read AMBIENT_STATE.md
# Add ambient awareness instructions

# Step 2.4: Deploy Star Office UI
cd /home/onemind/onemind-os/command-center/
git clone https://github.com/ringhyacinth/Star-Office-UI.git star-office
cd star-office && pip install -r backend/requirements.txt
# Create systemd service, enable

# Step 2.5: Deploy legacy-eyes on MacBook Pro
# SmolVLM-256M → vision digest → NATS
pip install transformers torch accelerate pillow opencv-python
# Run legacy-eyes.py

# Step 2.6: Verify ambient loop
cat /home/onemind/.openclaw/workspace/AMBIENT_STATE.md  # should have content
nats stream info SYSTEM  # should show messages
```

### Phase 3: DESKTOP COMMAND CENTER (Week 3-4)

```bash
# On MacBook Pro
# Step 3.1: Fork ClawX → onemind-desktop
cd ~/Documents/onemindos-v0.5/command-center/
git clone https://github.com/ValueCell-ai/ClawX.git onemind-desktop
cd onemind-desktop && npm install
# Configure gateway: onemind-main.tail717ba0.ts.net:18789
npm run dev

# Step 3.2: Add NATS WebSocket panel
# Connect to ws://onemind-main.tail717ba0.ts.net:8080 for live events

# Step 3.3: Add SmolVLM camera panel
# Step 3.4: Add MediaPipe gesture recognition
# Step 3.5: Add ambient overlay mode (transparent always-on-top)
```

### Phase 4: VOICE + MOBILE + GPU (Month 2)

```bash
# Step 4.1: Fork NixClaw for Gemini Live voice
# Step 4.2: Deploy Chatterbox TTS (self-hosted)
# Step 4.3: Spin up RTX 4000 Ada GPU Droplet (on-demand, 2hr/day)
# Step 4.4: Install Ollama + SmolVLM on GPU Droplet
# Step 4.5: Update ios-app connections
```

### Phase 5: PHYSICAL + SCALE (Month 3+)

```bash
# Step 5.1: Deploy RosClaw (Docker) for robot skills
# Step 5.2: NATS leaf node at house (farm sensors, ESP32)
# Step 5.3: Home Assistant → MQTT → NATS bridge
# Step 5.4: G1 robot integration
```

### Phase 6: SOVEREIGNTY (Month 4+)

```bash
# Step 6.1: Mac Studio M4 Ultra at home (replaces GPU Droplet)
# Step 6.2: All inference via Ollama locally (zero API dependency)
# Step 6.3: Self-hosted TTS/STT/VLM stack
# Step 6.4: VPS becomes thin relay only (NATS + gateway routing)
```

---

## 13. Cost Model

### Current Monthly Spend

| Item | Cost |
|------|------|
| VPS (CPU Droplet) | ~$50-100 |
| Azure OpenAI (GPT-5.2, GPT-4o) | ~$150-200 |
| OpenAI (TTS, embeddings) | ~$20-30 |
| Datadog | ~$0 (free tier) |
| **Total current** | **~$220-330/mo** |

### After Phase 1-2 Activation

| Item | Cost | Delta |
|------|------|-------|
| VPS (same droplet) | ~$50-100 | — |
| Azure OpenAI (reduced by OpenViking) | ~$80-120 | -$70-80 saved |
| OpenAI (TTS, embeddings) | ~$20-30 | — |
| Heartbeat tokens (GPT-4o-mini, 5 min) | ~$90-180 | New cost |
| Lossless Claw summaries | ~$10-15 | New cost |
| OpenViking overhead | ~$5 | New cost |
| SmolVLM vision | $0 | Local MacBook Pro |
| **Total Phase 1-2** | **~$255-450/mo** | +$35-120 |

### After Phase 4 (GPU)

| Item | Cost | Delta |
|------|------|-------|
| GPU Droplet (RTX 4000, 2hr/day) | ~$65 | New cost |
| API costs (partially replaced) | ~$50-80 | Major savings |
| Chatterbox TTS | $0 | Replaces OpenAI TTS |
| **Total Phase 4** | **~$220-370/mo** | Breaks even or saves |

### Long-Term Target (Mac Studio at Home)

| Item | Cost |
|------|------|
| Mac Studio M4 Ultra (one-time ~$4k) | $0/mo after purchase |
| All inference local (Ollama) | $0 |
| VPS for NATS + Gateway only | ~$20-40 |
| **Total long-term** | **~$20-40/mo** |

---

## Key Decisions Locked In

| Decision | Answer | Reasoning |
|----------|--------|-----------|
| GPU now? | **No. Phase 2+ on-demand.** | API costs don't justify 24/7 GPU yet. Use DO credit for serverless inference experiments. |
| MacBook as server? | **No. Development + sensors only.** | VPS for uptime. MacBook for vision, iMessage, desktop UI, dev. |
| Which office UI? | **OpenClaw Office (primary) + Star Office UI (vibes)** | OpenClaw Office = management. Star Office = fun visual status. |
| Heartbeat frequency? | **5 minutes + P0 webhook** | Balances cost vs awareness. P0 alerts bypass the gap immediately. |
| Realtime processing? | **Event-driven (direct messages) + scheduled (heartbeat) + webhook (P0)** | Three tiers, not one monolithic solution. |
| Model for heartbeat? | **GPT-4o-mini (defaults) + GPT-4o (Legacy override)** | Mini is cheap enough for 288/day. Legacy gets the smart model. |
| ClawX or build custom? | **Fork ClawX** | 4.5k stars, 116 releases > building from scratch. |
| Replace SwarmClaw? | **No. Upgrade to v1.1.5.** | Already deployed, configured, backed up. Add OpenClaw Office alongside. |
| Cloud or local? | **Cloud-first, migrate local incrementally** | MacBook Pro is powerful but not a server. |

---

*This is the build. Every service, every port, every command, every dollar. Start Phase 1 on your word.*

---

## 14. Skill Forge — Autonomous Skill Builder (MCP Builder Pattern)

### The Concept

The n8n-claw project ships a pattern called **MCP Builder** — the agent can build its own MCP tools on demand by searching for API docs, generating code, and self-deploying. We've replicated this pattern natively for OpenClaw as an installable skill.

### What It Is

**Skill Forge** (`skill-forge-1.0.0`) is an OpenClaw skill that gives Legacy the ability to autonomously create new skills at runtime. Instead of a human developer writing `index.js` + `package.json` + `SKILL.md` for every new API integration, Legacy can:

1. **Research** an API (via `web_search` / `web_fetch`)
2. **Generate** a complete OpenClaw skill with all three files
3. **Deploy** it into the workspace skills directory
4. **Validate** it passes structural checks
5. **Reload** the gateway to make it live

### 4 Tools

| Tool | Purpose |
| --- | --- |
| `skill_forge_create` | Generate a skill from a manual API definition (name, base_url, auth, tool specs) |
| `skill_forge_from_openapi` | Generate a skill from an OpenAPI/Swagger JSON spec URL — auto-extracts endpoints, params, auth |
| `skill_forge_validate` | Check any skill directory for structural issues (missing files, export mismatches, missing deps) |
| `skill_forge_list` | List all installed skills with health status, version, tool count, forge-created flag |

### Workflow Example

**Legacy receives**: "I need to integrate the Linear API for issue tracking"

**Legacy executes**:
1. `web_search("Linear API documentation OpenAPI spec")`
2. Finds the spec URL or reads the API docs
3. Calls `skill_forge_from_openapi` or `skill_forge_create` with the API definition
4. Skill is written to `skills/linear-issues-1.0.0/`
5. Calls `skill_forge_validate` to confirm clean structure
6. `exec("openclaw reload")` to hot-load the new skill

**Result**: Legacy now has `list_issues`, `create_issue`, `update_issue` tools — built without human intervention.

### Auth Types Supported

| Type | Header Generated |
| --- | --- |
| `bearer` | `Authorization: Bearer <token>` |
| `api-key-header` | Custom header (configurable) |
| `api-key-query` | API key as URL query parameter |
| `basic` | `Authorization: Basic <base64>` |
| `none` | No auth headers |

### Generated Skill Structure

```
skills/<name>-<version>/
├── index.js          ← Node.js tool implementations (zero external deps by default)
├── package.json      ← Tool definitions with JSON Schema parameters
├── SKILL.md          ← Natural language instructions with YAML frontmatter
└── .forge-created    ← Marker file (distinguishes auto-generated from manual)
```

### What This Replaces

| Before | After |
| --- | --- |
| Developer writes index.js manually | Legacy generates it from API spec |
| Developer writes package.json tool defs | Auto-generated with full JSON Schema |
| Developer writes SKILL.md | Auto-generated with tool docs + auth docs |
| Developer runs `node -c` to syntax check | `skill_forge_validate` does it |
| No visibility into skill health | `skill_forge_list` shows all skills + issues |
| OpenAPI specs require manual parsing | `skill_forge_from_openapi` does it automatically |

### Security Notes

- Generated skills use Node.js built-in `https`/`http` — no npm dependencies by default
- API keys always come from environment variables, never hardcoded
- `.forge-created` marker lets you audit which skills were auto-generated
- Always review generated `index.js` before deploying to production VPS
- Add new API keys to 1Password + `env.tpl` before using on VPS

### Location

```
openclaw/workspaces/legacy/skills/skill-forge-1.0.0/
├── index.js       (839 lines)
├── package.json   (169 lines)
└── SKILL.md       (126 lines)
```

---

## 15. NATS Manager — Nervous System Control Skill

### The Concept

The best NATS dashboard is Legacy itself. Instead of only looking at a web UI, Legacy can query, inspect, and manage NATS through conversation. This skill wraps both the NATS HTTP monitoring API (read-only server stats) and the NATS JetStream client (stream/consumer management) into 8 tools.

### 8 Tools

**Read-Only Monitoring (HTTP API — port 8222)**

| Tool | What It Shows |
| --- | --- |
| `nats_health` | Server vitals: uptime, CPU, memory, message rates, connection count, JetStream status |
| `nats_connections` | All connected clients with name, IP, subscriptions, message rates, idle time |
| `nats_subscriptions` | Active subscription patterns — what subjects are being listened to |

**JetStream Management (NATS client over port 4222)**

| Tool | What It Does |
| --- | --- |
| `nats_streams` | List all streams with message counts, storage size, consumer counts |
| `nats_stream_detail` | Deep dive on one stream: config, state, subjects, consumer list |
| `nats_consumers` | List consumers on a stream with delivery status, ack pending |
| `nats_stream_create` | Create a new JetStream stream with subjects, retention, storage |
| `nats_stream_purge` | Delete all messages from a stream (destructive) |

### Relationship to nexus-bridge

| Skill | Purpose |
| --- | --- |
| **nexus-bridge** | Messaging — dispatch tasks, publish events, query agents via NATS |
| **nats-manager** | Infrastructure — monitor health, manage streams, inspect connections |

nexus-bridge talks through the nervous system. nats-manager manages the nervous system itself.

### Standard Stream Layout

| Stream | Subjects | Purpose |
| --- | --- | --- |
| `SYSTEM` | `system.>`, `nexus.>` | System alerts, health, heartbeats |
| `DIGITAL` | `digital.>` | Agent tasks and results |
| `PHYSICAL` | `physical.>` | Robot commands, sensor data, HA events |
| `COMMS` | `comms.>` | Human messages across channels |
| `EVENTS` | `events.>` | Community, task, integration events |

### Environment Variables

| Var | Default | Purpose |
| --- | --- | --- |
| `NATS_URL` | `nats://localhost:4222` | NATS server for JetStream operations |
| `NATS_MONITOR_URL` | `http://localhost:8222` | NATS HTTP monitoring endpoint |

### NATS Tooling Stack (Recommended)

| Tool | Type | Job |
| --- | --- | --- |
| **nats-manager skill** | OpenClaw skill | Legacy monitors + manages NATS through conversation |
| **nats-dashboard** | Docker container | Web UI for visual monitoring (one `docker run` command) |
| **Nui** | Desktop app | MacBook Pro app for stream management, KV inspection |
| **NATS CLI** | CLI | Already installed. `nats stream ls`, `nats sub ">"` |

Skip Synadia Cloud for now — add it when you need cross-team account management.

### Location

```text
openclaw/workspaces/legacy/skills/nats-manager-1.0.0/
├── index.js       (445 lines)
├── package.json   (135 lines)
├── SKILL.md       (64 lines)
└── node_modules/  (nats@2.28)
```

---

## 16. Custom Domains — Caddy + Cloudflare Setup

### Overview

All frontends and the gateway are accessible via HTTPS custom subdomains under `onemindos.dev`, served by Caddy v2.11.2 with Let's Encrypt certificates obtained via Cloudflare DNS challenge.

### Domain Map

| Subdomain | Backend | Port | Service |
|-----------|---------|------|---------|
| `gateway.onemindos.dev` | OpenClaw Gateway Control UI | 18789 | `openclaw.service` |
| `office.onemindos.dev` | OpenClaw Office | 5180 | `openclaw-office.service` |
| `swarm.onemindos.dev` | SwarmClaw Dashboard | 3456 | `swarmclaw.service` |
| `studio.onemindos.dev` | OpenClaw Studio | 3100 | `openclaw-studio.service` |
| `nats.onemindos.dev` | NATS Native Monitor | 8222 | `onemind-nats` Docker |

### DNS (Cloudflare)

- **Domain**: `onemindos.dev`
- **Cloudflare Account ID**: `ce13703e0b76938c179ebc51e2945c7e`
- **DNS Records**: 5 A records, each pointing to `100.77.52.93` (Tailscale VPS IP), proxy OFF (DNS only — Caddy handles TLS)
- **API Token**: Stored in `/opt/onemind/.env.caddy` (mode 600). Account API Token with Zone:DNS:Edit + Zone:Zone:Read permissions.

### Caddy Configuration

**Binary**: `/usr/local/bin/caddy` v2.11.2 with `dns.providers.cloudflare` module (custom build via `xcaddy`).

**Caddyfile** (`/etc/caddy/Caddyfile`):
```caddyfile
{
  admin off
}

(cf_tls) {
  tls {
    dns cloudflare {env.CLOUDFLARE_API_TOKEN}
    key_type rsa2048
  }
}

gateway.onemindos.dev {
  import cf_tls
  reverse_proxy localhost:18789
  bind 100.77.52.93
}

office.onemindos.dev {
  import cf_tls
  reverse_proxy localhost:5180
  bind 100.77.52.93
}

swarm.onemindos.dev {
  import cf_tls
  reverse_proxy localhost:3456
  bind 100.77.52.93
}

studio.onemindos.dev {
  import cf_tls
  reverse_proxy localhost:3100
  bind 100.77.52.93
}

nats.onemindos.dev {
  import cf_tls
  reverse_proxy localhost:8222
  bind 100.77.52.93
}
```

**Key decisions**:
- `key_type rsa2048` — Mac's LibreSSL 3.3.6 can't handle ECDSA certs from Let's Encrypt E7. RSA certs use R12 issuer which works everywhere.
- `bind 100.77.52.93` — Only listen on Tailscale IP. Not exposed to public internet.
- `admin off` — No Caddy admin API (security).
- DNS challenge — No HTTP challenge needed since Caddy doesn't listen on public IPs.

### Caddy systemd Service

**File**: `/etc/systemd/system/caddy.service`

```ini
[Unit]
Description=Caddy — HTTPS Reverse Proxy
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=caddy
Group=caddy
EnvironmentFile=/opt/onemind/.env.caddy
Environment=HOME=/var/lib/caddy
Environment=XDG_CONFIG_HOME=/var/lib/caddy/.config
Environment=XDG_DATA_HOME=/var/lib/caddy/.local/share
ExecStart=/usr/local/bin/caddy run --config /etc/caddy/Caddyfile --environ
ExecReload=/usr/local/bin/caddy reload --config /etc/caddy/Caddyfile
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

**Cert storage**: `/var/lib/caddy/.local/share/caddy/certificates/acme-v02.api.letsencrypt.org-directory/`

### Cloudflare API Token

Stored in `/opt/onemind/.env.caddy`:
```
CLOUDFLARE_API_TOKEN=<account-api-token>
```

Also referenced in `/opt/onemind/.env.tpl` for the `cloudflare-api-1.0.0` skill.

### Gateway Configuration for Custom Domains

In `/home/onemind/.openclaw/openclaw.json`, the `gateway.controlUi.allowedOrigins` must include all custom domains:

```json
{
  "gateway": {
    "controlUi": {
      "allowedOrigins": [
        "https://*.tail717ba0.ts.net",
        "http://*.tail717ba0.ts.net",
        "https://onemind-main.tail717ba0.ts.net",
        "http://localhost:3000",
        "http://localhost:3456",
        "http://127.0.0.1:3000",
        "http://127.0.0.1:3456",
        "https://gateway.onemindos.dev",
        "https://office.onemindos.dev",
        "https://swarm.onemindos.dev",
        "https://studio.onemindos.dev",
        "https://nats.onemindos.dev",
        "https://*.onemindos.dev"
      ]
    },
    "trustedProxies": ["100.64.0.0/10", "127.0.0.1", "::1", "::ffff:127.0.0.1"]
  }
}
```

### Studio Gateway Token

Studio reads the gateway token from `~/.openclaw/openclaw.json`, but the config has `"token": "${OPENCLAW_GATEWAY_TOKEN}"` (an env var reference). Studio doesn't resolve env vars — it reads the literal string.

**Fix**: Create `/home/onemind/.openclaw/openclaw-studio/settings.json` with the actual token:

```json
{
  "gateway": {
    "url": "ws://localhost:18789",
    "token": "<actual-64-char-gateway-token>"
  }
}
```

On rotation, update this file: `op read "op://aaeyris5z44pdqspk277u2oppi/VPS OpenClaw/gateway-token"` and write to the settings file.

### Tailscale — CRITICAL: No Serves on Port 443

Tailscale `serve` on `:443` conflicts with Caddy. If Tailscale intercepts `:443`, browsers get SSL errors when connecting to custom domains. **All Tailscale serves have been removed.** Caddy is the sole HTTPS handler.

To verify: `tailscale serve status` should return "No serve config".

If Tailscale serves re-appear (they can auto-restore), remove them:
```bash
tailscale serve reset  # Nuclear: removes ALL serves
# Or selectively:
tailscale serve --https 443 off
```

### Adding a New Subdomain

1. Create DNS A record in Cloudflare: `new.onemindos.dev → 100.77.52.93` (DNS only, proxy OFF)
2. Add server block to `/etc/caddy/Caddyfile`:
   ```caddyfile
   new.onemindos.dev {
     import cf_tls
     reverse_proxy localhost:<port>
     bind 100.77.52.93
   }
   ```
3. Add `https://new.onemindos.dev` to `gateway.controlUi.allowedOrigins` in `openclaw.json` (if it needs gateway access)
4. `systemctl reload caddy` — auto-obtains cert via DNS challenge
5. Verify: `curl -sS https://new.onemindos.dev/`
