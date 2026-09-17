# OneMind OS — Production Plan

> **Goal**: Get Legacy into daily-driver production as a sovereign, omnipresent AI co-pilot.
> **Date**: March 15, 2026

---

## Part 1 — The Honest Audit (What's Real vs Dead Weight)

### What's Actually Connected and Working

| Component | Status | Notes |
|-----------|--------|-------|
| OpenClaw Gateway (9 agents) | **LIVE** | Running on VPS, responding, Discord connected |
| NATS JetStream | **LIVE** | Container healthy, 4222/1883/8222 open |
| PostgreSQL/pgvector | **LIVE** | Container healthy, 5432 open |
| legacy-bridge (NATS→OpenClaw) | **LIVE** | Container healthy, relaying events |
| nats-ingester | **LIVE** | Container healthy |
| Discord channels (9 agents) | **LIVE** | All 9 channels configured, threadBindings on |
| iMessage (BlueBubbles) | **LIVE** | Connected via Tailscale to Mac |
| Codex sync (pull+watch) | **LIVE** | Auto-pull every 1min, auto-commit changes |
| Memory search (hybrid BM25+vector) | **LIVE** | Indexes codex + memory + sessions |
| OpenClaw Studio | **LIVE** | Running, Tailscale proxied |
| SwarmClaw | **LIVE** | Running, Tailscale proxied |
| Datadog | **LIVE** | Monitoring active |

### What's Defined But NOT Connected (Dead Weight)

| Component | Status | Verdict |
|-----------|--------|---------|
| **Agno AgentOS** (Python, port 7777) | **NOT RUNNING** | Never started. No Python agent code deployed. |
| **LiveKit** (port 7880) | **NOT RUNNING** | Defined in docker-compose but never started |
| **n8n** (port 5678) | **NOT RUNNING** | Defined in docker-compose but never started |
| **OpenMind OM1** (robots) | **NOT RUNNING** | No hardware connected |
| **HuggingFace LeRobot** | **NOT RUNNING** | No hardware |
| **Heartbeat** | **NOT CONFIGURED** | Empty `{}` in openclaw.json — agents have no heartbeat |
| **Voice** | **PARTIAL** | Discord voice enabled (OpenAI TTS, alloy voice), but no ambient/always-on |
| **Vision** | **NOT RUNNING** | No legacy-eyes, no camsnap, no ambient state |

### Skills That Reference Dead Systems

| Skill | Problem |
|-------|---------|
| `nexus-bridge` | Dispatches to Agno agents at localhost:7777 — **Agno is not running** |
| `homestead-control` | References 6-building node grid — **no homestead yet** |
| `clickup-1.0.0` | References ClickUp workspace — **you use Todoist now** |
| `desktop-control-1.0.0` | Windows-focused (mspaint, calc) — **you're on Mac** |
| `browser-1.0.0` | Status: "pending-port" — **not implemented** |
| `vision-1.0.0` | Exists but no vision loop running |
| `credit-repair-1.0.0` | + `lob-1.0.0` — niche, not daily-driver |
| `uber-1.0.0` | Requires OAuth that's not configured |
| `spotify-1.0.0` | Requires OAuth refresh token setup |

---

## Part 2 — The Strategic Decision

### The Core Question: Where Does Orchestration Live?

**Answer: Keep it under OpenClaw. Here's why.**

OpenClaw is already running with 9 agents, agent-to-agent messaging, subagent spawning, Discord threading, iMessage, memory search, codex access, and 29 skills. It IS your operating system. The question was never "should we add another orchestration layer?" — it was "are we using what we already have?"

**What you have that you're NOT using:**

1. **Heartbeat** — Configured as `{}`. Your 9 agents have zero proactive behavior. They sit idle until spoken to.
2. **Agent-to-agent** — `agentToAgent.enabled: true` with all 9 agents allowed. Legacy can already dispatch to Spartan, Oracle, etc. natively. No NATS required.
3. **Subagents** — `maxConcurrent: 8, maxSpawnDepth: 3`. Legacy can spawn 8 parallel workers. Not using it.
4. **Thread bindings** — `spawnSubagentSessions: true`. Discord threads auto-spawn agent sessions. Working.

**What to stop chasing:**

| Don't | Why |
|-------|-----|
| Agno AgentOS | OpenClaw's 9 agents + subagents + agent-to-agent already does this natively |
| LangGraph/LangChain | OpenClaw has lobster routines (named workflows) — same concept, native |
| CrewAI/AutoGen | Multi-agent with roles = your 9 domain agents already |
| Paperclip | Governance layer for 50+ agents. You have 9. Premature. |
| SwarmClaw for orchestration | It bypasses OpenClaw features for OpenClaw agents. Wrong tool. |

### The Right Architecture (Lean)

```
Zeus (voice, text, vision)
    ↓
OpenClaw Gateway (Damien VPS)
    ├── Legacy ⚡ (sovereign orchestrator, gpt-5.2, full tools)
    │   ├── Can talk to all 8 domain agents via agent-to-agent
    │   ├── Can spawn 8 parallel subagents for deep work
    │   └── Has 29 skills (tools) for external actions
    │
    ├── 8 Domain Agents (gpt-4o, messaging+search tools)
    │   └── Each reads codex, searches memory, does domain-specific work
    │
    ├── NATS JetStream (event backbone)
    │   └── For IoT, sensors, future hardware — NOT for agent dispatch
    │
    └── PostgreSQL/pgvector (data store)
        └── Memory vectors, session data
```

**NATS stays but changes role:**

- **Now**: Supposed to connect Agno (dead) and OM1 (no hardware)
- **Production**: Event bus for vision feed, session events, sensor data, alerting
- **Future**: When you DO get hardware (homestead, robots), NATS is ready

**SwarmClaw gets demoted:**

- Keep it running for now (it's not hurting anything)
- Use it only for direct-API agents (Claude, Gemini, GPT without OpenClaw)
- Don't use it as the primary interface — use Discord + iOS app + voice

### Provider Strategy (Within OpenClaw)

OpenClaw already supports model aliases. Your config has:

```json
"models": {
  "azure-openai/gpt-5.2": { "alias": "gpt52" },
  "azure-openai/gpt-4o": { "alias": "gpt4o" },
  "azure-openai/gpt-4o-mini": { "alias": "gpt4omini" },
  "gradient/llama3.3-70b-instruct": { "alias": "llama70b" }
}
```

**Add these providers directly to OpenClaw** (not through SwarmClaw):

| Provider | Purpose | How |
|----------|---------|-----|
| Anthropic (Claude) | Give Legacy Claude for deep reasoning | Add `anthropic/claude-sonnet-4-6` to model aliases |
| Google Gemini | Long context + multimodal | Add `google/gemini-2.5-pro` to model aliases |
| Groq | Fast/cheap for subagent workers | Add `groq/llama-3.3-70b` to model aliases |
| Ollama (local) | Private/offline fallback | Install on VPS, add `ollama/qwen3-30b` |

Legacy can then use `@gpt52` for one task and `@claude` for another — all within OpenClaw.

---

## Part 3 — Ambient Presence (The Big Upgrade)

This is the real unlock. Legacy goes from "chatbot that waits" to "presence that sees and hears."

### Layer 1: Heartbeat (Enable What's Already There)

The heartbeat system exists in OpenClaw. It's just set to `{}`. Turn it on.

**Config change** (openclaw.json → agents.defaults.heartbeat):

```json
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
```

Legacy's heartbeat (override — more frequent, smarter model):

```json
{
  "every": "5m",
  "target": "last",
  "lightContext": false,
  "model": "azure-openai/gpt-4o"
}
```

**Why 5 minutes, not 1 minute**: 1-minute heartbeat = 1,440 LLM calls/day. At gpt-4o that's ~$15-30/day just for "checking in." 5-minute = 288 calls/day, ~$3-6/day. Start here, tune later.

Domain agents get slower heartbeats (or none — they respond when summoned via agent-to-agent).

### Layer 2: Vision (legacy-eyes)

The attached `legacy-eyes.py` is the right approach. Deploy it as a systemd service alongside OpenClaw.

**How it works:**

```
Webcam/RTSP (1fps) → capture frame → Claude Haiku vision API → digest
    → Write to AMBIENT_STATE.md in Legacy's workspace
    → Legacy reads on heartbeat → "Zeus is at his desk, focused on code"
```

**Cost**: Claude Haiku vision at 1 frame every 10 seconds = ~$0.50/day. Extremely cheap for environmental awareness.

**Phase 1 (now)**: Mac webcam via desktop app (stream to VPS)
**Phase 2 (later)**: Ray-Ban Meta glasses via RTSP
**Phase 3 (homestead)**: Multiple cameras per building

### Layer 3: NATS Events (legacy-pulse)

The attached `legacy-pulse.py` is correct. It subscribes to NATS subjects and writes events to `AMBIENT_STATE.md` alongside the vision digest. Legacy reads both on every heartbeat.

**What legacy-pulse captures:**
- Session lifecycle (Zeus opened/closed Discord)
- System alerts (droplet overload, service failures)
- Future: IoT sensor events, location changes, doorbell, etc.

### Layer 4: Voice (Without ElevenLabs Dependency)

**Current state**: Discord voice works with OpenAI TTS (alloy voice). This is already independent of ElevenLabs.

**Options for always-on voice:**

| Approach | Pros | Cons |
|----------|------|------|
| **Discord voice channel (current)** | Already works, VAD built-in, barge-in, dave encryption | Requires Discord open, not ambient |
| **LiveKit persistent room** | True ambient, multi-modal, supports video | Needs LiveKit deployed, client work |
| **OpenClaw iOS app voice** | Native, always with you, push-to-talk or VAD | Needs app development |
| **Web Speech API (browser)** | Zero infrastructure, works now | Requires browser tab open |

**Recommendation**:
- **Phase 1**: Discord voice (it works NOW, OpenAI TTS, no ElevenLabs)
- **Phase 2**: iOS app voice via OpenClaw Gateway WebSocket
- **Phase 3**: LiveKit for true ambient (when you need multi-room/multi-device)

**For TTS model**: OpenAI's `tts-1` is already fast and cheap ($0.015/1K chars). No ElevenLabs dependency. If you want a custom voice later, you can fine-tune via OpenAI's voice API or use a local model (Piper TTS on VPS — free, fast, runs locally).

---

## Part 4 — What to Cut (Run Lean)

### Services to Stop

```bash
# Kill orphan SwarmClaw on port 3457
ssh root@104.248.4.212 "kill \$(ss -tlnp | grep 3457 | grep -oP 'pid=\K[0-9]+')"
```

### Skills to Archive (move to `skills-archive/`)

| Skill | Reason |
|-------|--------|
| `clickup-1.0.0` | You use Todoist, not ClickUp |
| `homestead-control` | No homestead yet — bring back when you buy property |
| `desktop-control-1.0.0` | Windows-focused, not useful on Mac/VPS |
| `browser-1.0.0` | "pending-port" — use `agent-browser-0.2.0` instead |
| `credit-repair-1.0.0` | Niche tool, not daily-driver |
| `lob-1.0.0` | Physical mail — used only with credit-repair |
| `uber-1.0.0` | OAuth not configured, not daily-use |
| `transit-1.0.0` | Unlikely to be used regularly |
| `zoho-mail-1.0.0` | Duplicate — `zoho-cli` already covers mail |
| `todoist-api` | Duplicate — `todoist` (CLI-based) already works |
| `find-skills-0.1.0` | Meta-skill for discovering skills — not needed in production |
| `skill-creator-0.1.0` | Development tool, not runtime |

**Keep (13 core skills):**

| Skill | Purpose |
|-------|---------|
| `nexus-bridge` | NATS dispatch (keep for events, not Agno) |
| `home-assistant-1.0.0` | Smart home control |
| `todoist` | Task management |
| `github-api-1.0.3` | GitHub operations |
| `onepassword-1.0.0` | Secret management |
| `obsidian-1.0.0` | Vault operations |
| `rss-digest-1.0.0` | Newsletter digestion |
| `recipe-scraper-1.0.0` | Meal planning (HP domain) |
| `resend-1.0.0` | Email delivery |
| `twilio-1.0.0` | SMS escalation |
| `agent-browser-0.2.0` | Web automation |
| `api-gateway-1.0.48` | 100+ API proxy (Maton) |
| `proactive-agent-3.1.0` | Proactive behavior patterns |
| `spotify-1.0.0` | Music control (if OAuth is set up) |
| `weather-1.0.0` | Daily weather |
| `zoho-cli` | CRM + calendar |
| `vision-1.0.0` | Vision API (for legacy-eyes integration) |
| `hitl-escalate-1.0.0` | Human-in-the-loop escalation |

### Workspace Files to Update

Legacy's current AGENTS.md references Agno at localhost:7777 extensively. The HEARTBEAT.md references ClickUp. The TOOLS.md references Agno HTTP API. These need updating.

### Docker Services Decision

| Container | Keep? | Notes |
|-----------|-------|-------|
| `onemind-nats` | **YES** | Core event backbone |
| `onemind-omnidb` | **YES** | Data store |
| `onemind-legacy-bridge` | **YES** | NATS→OpenClaw relay |
| `onemind-nats-ingester` | **YES** | Event processing |
| n8n (not running) | **NO** | OpenClaw lobster routines replace this |
| LiveKit (not running) | **LATER** | Phase 3 of voice |

---

## Part 5 — The Production Deploy Sequence

### Phase 1: Clean House + Heartbeat (This Week)

1. **Archive unused skills** — move 12 skills to `skills-archive/` directory
2. **Update Legacy's AGENTS.md** — remove Agno references, add agent-to-agent dispatch as primary
3. **Update HEARTBEAT.md** — replace ClickUp with Todoist, remove Agno health check
4. **Update TOOLS.md** — remove Agno dispatch, document available skills accurately
5. **Update BOOT.md** — remove Agno health check
6. **Enable heartbeat in openclaw.json** — 5-minute for Legacy, off for domain agents initially
7. **Kill orphan SwarmClaw on port 3457**
8. **Test**: Restart OpenClaw, verify heartbeat fires every 5 min to Discord

### Phase 2: Ambient Vision + Events (Next Week)

1. **Deploy legacy-eyes.py** on Mac (webcam-based initially)
2. **Deploy legacy-pulse.py** on VPS alongside NATS
3. **Create AMBIENT_STATE.md template** in Legacy workspace
4. **Update HEARTBEAT.md** to read ambient state
5. **Update SOUL.md** with ambient presence behavior
6. **Create systemd service for legacy-pulse**
7. **Test**: Verify Legacy reads vision + events on heartbeat

### Phase 3: Voice + Multi-Provider (Week After)

1. **Set up Discord voice** — test always-listening in a dedicated channel
2. **Add Anthropic provider** to openclaw.json model aliases
3. **Add Groq provider** for fast subagent work
4. **Test model switching**: Legacy uses gpt-5.2 normally, Claude for deep reasoning
5. **Optional**: Install Ollama on VPS for local model fallback

### Phase 4: Polish + Scale (Ongoing)

1. **Tune heartbeat frequency** based on cost and usefulness
2. **Add vision to VPS** (RTSP camera or screen capture)
3. **iOS app voice** integration
4. **Evaluate Paperclip** if you grow beyond 9 agents
5. **Homestead setup** when property is acquired (bring back homestead-control, OM1)

---

## Part 6 — Updated Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────────┐
│              HOW ZEUS REACHES LEGACY                                  │
│  Discord  iMessage  Voice  iOS App  Browser  Webcam  NATS Events     │
└─────────────────────────┬────────────────────────────────────────────┘
                          │
              ┌───────────▼──────────────┐
              │     OPENCLAW GATEWAY      │  Port 18789
              │     (Damien VPS)          │
              │                           │
              │  Legacy ⚡ (gpt-5.2)      │  ← sovereign orchestrator
              │    ├── 29 skills          │     full tool access
              │    ├── agent-to-agent     │     dispatches to domain agents
              │    ├── 8 subagent slots   │     parallel deep work
              │    └── heartbeat (5m)     │     proactive checks
              │                           │
              │  8 Domain Agents:         │  ← gpt-4o, messaging tools
              │    Spartan ⚔️  Oracle 🔮   │     read/write codex
              │    Eden 🌳  Empire 👑      │     web search
              │    Haven 🌿  Guardian 🛡️   │     memory search
              │    Grid ⚙️  Trinity 🔺     │
              └───────┬──────┬────────────┘
                      │      │
         ┌────────────┘      └──────────────┐
         │                                  │
┌────────▼─────────┐            ┌───────────▼───────────┐
│  NATS JetStream   │            │  PostgreSQL/pgvector   │
│  (event backbone) │            │  (data + vectors)      │
│                   │            └───────────────────────┘
│  Subscriptions:   │
│  ├── legacy-pulse │ ← system events → AMBIENT_STATE.md
│  ├── legacy-eyes  │ ← vision feed → AMBIENT_STATE.md
│  └── (future IoT) │
└───────────────────┘

┌──────────────────────────────────────────────────────────────────────┐
│  AMBIENT PRESENCE (runs alongside gateway)                            │
│                                                                       │
│  legacy-eyes.py ──→ 1fps webcam → Claude Haiku ──→ AMBIENT_STATE.md │
│  legacy-pulse.py ──→ NATS events ─────────────────→ AMBIENT_STATE.md │
│  Heartbeat (5m) ──→ reads AMBIENT_STATE.md → surfaces what matters   │
└──────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────┐
│  EXTERNAL SERVICES (via skills)                                       │
│                                                                       │
│  Discord · iMessage · Todoist · GitHub · Home Assistant               │
│  1Password · Zoho · Resend · Twilio · Brave Search · Maton Gateway   │
│  Spotify · OpenAI TTS · RSS feeds                                     │
└──────────────────────────────────────────────────────────────────────┘
```

---

## Part 7 — Answers to Your Direct Questions

### "What can be removed?"

12 skills archived, Agno removed from all references, n8n dropped, ClickUp references replaced with Todoist. See Part 4.

### "Should we have multi-agent with OpenClaw, or keep it simple?"

**Keep it under OpenClaw.** You already have multi-agent (9 agents with agent-to-agent). Adding external frameworks (Agno, LangGraph, CrewAI) adds complexity with zero benefit — OpenClaw already does agent dispatch, subagent spawning, and domain routing natively.

### "Where should the orchestration layer be?"

**In OpenClaw itself.** Legacy IS the orchestrator. The `agent-to-agent` config is already enabled. Legacy can message Spartan, Oracle, etc. directly. For deeper work, Legacy spawns subagents (up to 8 concurrent, 3 levels deep). No external layer needed until you exceed 9 agents or need formal governance (budget caps, approval chains) — that's when Paperclip enters.

### "Different providers within the platform?"

**Yes, via OpenClaw model aliases.** Add Anthropic, Groq, Gemini, Ollama as model options. Legacy picks the right model per task. No external framework needed — OpenClaw supports any OpenAI-compatible API and has native Anthropic/Google support.

### "Giving Legacy Claude CLI is more important than Agno or LangChain"

**Correct.** Legacy can use Claude (or any model) directly through OpenClaw's model routing. Adding Claude as `anthropic/claude-sonnet-4-6` in the model config is 1 line. Integrating Agno or LangChain would be weeks of work for capabilities OpenClaw already has.

### "Voice without ElevenLabs dependency?"

**Already done.** Your config shows `tts.provider: "openai"` with voice `alloy`. ElevenLabs is removed. OpenAI TTS is fast ($0.015/1K chars), no external dependency beyond OpenAI API. For truly local: Piper TTS runs on VPS, free, offline.

### "Vision on desktop?"

**Deploy legacy-eyes.py.** Start with Mac webcam, capture 1fps, send to Claude Haiku, write digest to AMBIENT_STATE.md. Legacy reads it on heartbeat. Cost: ~$0.50/day. Later: add VisionClaw glasses, RTSP cameras, multi-camera feeds.

### "What about NATS — keep it?"

**Yes, but change its role.** NATS is NOT for agent-to-agent dispatch (OpenClaw does that natively). NATS IS for:
- Vision frame events
- Session lifecycle events
- System alerts
- IoT/sensor data (future)
- Cross-runtime communication (when you add hardware)

---

## Part 8 — Research Integration (14 GitHub Repos Evaluated)

> Full analysis in [GITHUB-RESEARCH.md](GITHUB-RESEARCH.md)

### Top Picks — Install / Fork / Deploy

| Priority | Repo | Stars | Action |
|----------|------|-------|--------|
| P0 | **lossless-claw** | 1.9k | Install NOW — replaces lossy sliding-window compaction with DAG summaries. `openclaw plugins install @martian-engineering/lossless-claw` |
| P0 | **OpenViking** | 12.7k | Deploy as context database — +46% accuracy, -91% token cost vs native OpenClaw memory |
| P1 | **obsidian-skills** (by kepano) | 14.1k | Install — proper Obsidian syntax for vault writes |
| P1 | **ClawX** | 4.5k | Fork → `onemind-desktop` — add MediaPipe + SmolVLM + ambient overlay |
| P2 | **ClawVault** | 584 | Evaluate after P0 — entity graph + Obsidian hooks |

### Key Decisions from Research

- **SwarmClaw**: Upgrade v0.7.7 → v1.1.5 (missions, orchestrators, structured sessions)
- **Desktop**: Fork ClawX, don't build from scratch
- **Memory**: Lossless Claw (conversation) + OpenViking (knowledge) — complementary, not competing
- **iOS**: Study Chowder patterns, don't replace our ios-app
- **Voice**: Hume AI emotional voice → Phase 4+ (luxury, not foundation)
- **ClawdAgent**: Inspiration only (9 stars, 2 contributors, different architecture)

### Revised Phase Plan

**Phase 1 — Memory Foundation (Week 1)**
1. Install Lossless Claw plugin
2. Configure session.reset.idleMinutes=10080
3. Install obsidian-skills (kepano)
4. Upgrade SwarmClaw to v1.1.5
5. Enable heartbeat

**Phase 2 — Context Intelligence (Week 2-3)**
1. Deploy OpenViking server on VPS
2. Configure OpenViking OpenClaw plugin
3. Evaluate ClawVault
4. Archive 12 unused skills

**Phase 3 — Desktop + Vision (Week 3-4)**
1. Fork ClawX → onemind-desktop
2. Add SmolVLM-256M camera panel
3. Add MediaPipe Hands gesture recognition
4. Add ambient overlay mode

**Phase 4 — Voice + Emotion (Month 2+)**
1. Evaluate Chatterbox TTS vs OpenAI TTS
2. Study Hume AI EVI for emotional voice
3. Implement Gemini Live bridge (NixClaw)

**Phase 5 — Physical (Month 3+)**
1. RosClaw → Unitree G1
2. DimensionalOS spatial memory
3. peaq DID machine identity

---

## Summary: The 80/20

The 80% of value comes from 20% of changes:

1. **Install Lossless Claw** (1 plugin → Legacy never forgets a conversation)
2. **Turn on heartbeat** (1 config change → Legacy becomes proactive)
3. **Deploy OpenViking** (1 server → +46% memory accuracy, -91% token cost)
4. **Upgrade SwarmClaw** (v0.7.7 → v1.1.5 → missions, orchestration, structured sessions)
5. **Fork ClawX** (desktop foundation → add MediaPipe + SmolVLM)

These 5 things get you from "chatbot that waits" to "AI co-pilot that remembers everything, sees, hears, and acts."
