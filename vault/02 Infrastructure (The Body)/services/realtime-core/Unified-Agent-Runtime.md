---
title: Unified Agent Runtime — Research & Architecture Decision
domain: 02 Infrastructure (The Body)
status: design-research
created: 2026-05-25
tags: [onemind-os, agent-runtime, hermes, openclaw, livekit, realtime, voice, unified]
---

# Unified Agent Runtime — Research & Architecture Decision

The fragmentation problem: Hermes (memory, skills, tools) and OpenClaw (voice, multi-channel gateway, realtime) are both installed but separate. The goal is ONE unified runtime where Legacy speaks, sees, hears, reasons, and acts — through any interface — with full memory and context.

---

## Critical Constraint: No Telegram/Discord

All interaction is direct — voice, WebUI, CLI, or device nodes. No messaging platforms as agent interfaces.

---

## The Research

### OpenClaw (TypeScript) — What It Already Has

**Version installed:** 2026.5.20
**Location:** `~/.nvm/versions/node/v24.15.0/lib/node_modules/openclaw`

OpenClaw is a **multi-channel AI gateway** with:

| Capability | Status | Details |
|---|---|---|
| Voice — Talk Mode (macOS/iOS/Android) | ✅ Built-in | Always-on overlay, STT → agent → TTS, interrupt detection |
| Realtime Voice (WebRTC) | ✅ Built-in | OpenAI Realtime API, Google Realtime, provider-websocket |
| Realtime Transcription | ✅ Built-in | Deepgram, ElevenLabs, xAI, Mistral streaming |
| TTS Providers | ✅ Built-in | ElevenLabs, Azure Speech, MLX local, system TTS |
| Multi-channel inbox | ✅ Built-in | WhatsApp, Telegram, Slack, Discord, Signal, iMessage, Matrix, IRC, 20+ |
| Node system (devices) | ✅ Built-in | macOS/iOS/Android nodes pair to gateway, expose camera/canvas/device |
| Canvas (live HTML panel) | ✅ Built-in | Agent-controlled live UI |
| Memory (active-memory, memory-wiki) | ✅ Extensions | Session memory, wiki-style long-term |
| Plugin SDK | ✅ Built-in | TypeScript, provider plugins, tool plugins, hooks |
| Local model support | ✅ Extensions | Ollama, LM Studio, vLLM, sglang |
| MCP support | ✅ Built-in | Native MCP client |
| Multi-agent routing | ✅ Built-in | Route channels to isolated agent workspaces |
| Device pairing (glasses/phone) | ✅ Built-in | WS node protocol, camera/screen commands |

**Key: OpenClaw's Talk Mode is EXACTLY the "join a room and stream" pattern you described.**

Talk mode architecture:
```
macOS/iOS/Android native app
    │
    │  WebSocket to Gateway
    ▼
Gateway (always-on process)
    │
    ├── STT (local speech recognition OR realtime provider)
    ├── Agent brain (LLM call with tools, memory, context)
    ├── TTS (ElevenLabs, MLX local, system, or realtime provider)
    │
    │  Modes:
    ├── "native" = local STT → Gateway chat → TTS (turn-based with streaming)
    ├── "realtime" = WebRTC/provider-websocket (full-duplex, like OpenAI Realtime API)
    └── "gateway-relay" = Gateway owns the audio stream, bridges to any provider
```

**Talk mode `brain` options:**
- `agent-consult`: Realtime voice routes tool calls through the full Gateway agent (memory, tools, policy) — THIS IS WHAT WE WANT
- `direct-tools`: Tool calls go straight to provider (no agent mediation)
- `none`: Transcription only, no agent response

### Hermes (Python) — What It Already Has

| Capability | Status | Details |
|---|---|---|
| Persistent memory (cross-session) | ✅ Built-in | SQLite + injected into every turn |
| Skills (procedural knowledge) | ✅ Built-in | 100+ skills, loadable, updatable |
| Tools (comprehensive) | ✅ Built-in | NATS, HA, SaaS, filesystem, web, browser, code_exec, terminal |
| Session history + search | ✅ Built-in | FTS5 across all past conversations |
| Multi-provider models | ✅ Built-in | Bedrock, Azure, Anthropic, OpenAI, OpenRouter, Groq, DeepSeek, local |
| Plugin SDK | ✅ Built-in | Python hooks (pre_llm_call, on_session_start, etc.) |
| Cron jobs | ✅ Built-in | Scheduled background tasks with full agent |
| Delegation (subagents) | ✅ Built-in | Spawn worker agents for parallel tasks |
| Profiles | ✅ Built-in | Isolated instances with separate config/memory |
| NATS native toolset | ✅ Built-in | request, publish, KV get/put/keys, stream read |
| Voice (TTS output) | ✅ Built-in | edge-tts, OpenAI TTS |
| Voice (STT input) | ❌ Not built | No mic → transcribe path |
| Realtime streaming | ❌ Not built | Turn-based only |
| Device pairing | ❌ Not built | No node/device protocol |
| Vision (continuous) | ❌ Not built | Can analyze images, not continuous frames |
| WebRTC/audio streaming | ❌ Not built | No real-time audio pipeline |

### Mark-XXXIX (Python) — Interface Reference

**What it is:** A Python desktop app using Gemini's `LiveConnectConfig` for real-time voice + vision + tool calling.

Key architecture:
```python
# Uses Google Gemini's native audio streaming
LIVE_MODEL = "models/gemini-2.5-flash-native-audio-preview"
config = types.LiveConnectConfig(
    response_modalities=["AUDIO"],
    output_audio_transcription={},
    input_audio_transcription={},
    system_instruction=system_prompt,
    tools=[{"function_declarations": TOOL_DECLARATIONS}],
)
# Opens a persistent bidirectional audio session
# Audio in → Gemini processes → Audio out + tool calls
```

- PyQt6 UI (desktop overlay, visualizations)
- `sounddevice` for mic capture
- Gemini handles STT + reasoning + TTS in ONE model call (native audio)
- Memory via JSON file (simple)
- Tools: screen capture, browser control, file ops, app launch

**Takeaway:** Mark-XXXIX is a good UI reference but architecturally simple — single model, single device, no multi-agent, no NATS, no persistent memory system. It's a Jarvis demo, not an infrastructure.

### LiveKit Agents (Python) — The Room Pattern

**What it is:** A framework for building real-time voice AI agents that connect to "rooms" via WebRTC.

Key pattern:
```python
server = AgentServer()

@server.rtc_session()
async def entrypoint(ctx: JobContext):
    session = AgentSession(
        vad=silero.VAD.load(),
        stt=deepgram.STT(model="nova-3"),
        llm=openai.LLM(model="gpt-4.1-mini"),
        tts=cartesia.TTS(model="sonic-3"),
    )
    agent = Agent(
        instructions="You are a voice assistant.",
        tools=[lookup_weather],
    )
    await session.start(agent=agent, room=ctx.room)
```

- **Room concept**: Client connects to a room → agent joins same room → bidirectional audio/video
- **MCP support**: Native MCP tool integration
- **Multi-agent handoff**: Agents can transfer to other agents mid-conversation
- **Self-hosted**: LiveKit server is open-source, runs locally
- **Console mode**: `python myagent.py console` — local mic/speaker, no server needed

**Takeaway:** LiveKit is the best "room-based streaming" framework. But it's a VOICE framework, not an agent framework — it doesn't have memory, skills, NATS integration, or persistent context.

---

## The Comparison Matrix

| Feature | Hermes (Python) | OpenClaw (TypeScript) | LiveKit Agents (Python) | Mark-XXXIX (Python) |
|---|---|---|---|---|
| **Persistent memory** | ✅ Best-in-class | ✅ Good (extensions) | ❌ None | ❌ JSON file |
| **Skills/procedures** | ✅ 100+ skills | ❌ Limited | ❌ None | ❌ None |
| **Tool ecosystem** | ✅ NATS, HA, SaaS, 50+ tools | ✅ MCP, plugins | ✅ MCP, function_tool | ✅ Basic (screen, browser) |
| **NATS integration** | ✅ Native toolset | ❌ None | ❌ None | ❌ None |
| **Voice — Talk mode** | ❌ TTS only | ✅ Full (STT+TTS+interrupt) | ✅ Full (STT+TTS+VAD) | ✅ Full (Gemini native) |
| **Realtime streaming** | ❌ Turn-based | ✅ WebRTC + relay | ✅ WebRTC rooms | ✅ Gemini LiveConnect |
| **Multi-device/node** | ❌ None | ✅ iOS/Android/macOS nodes | ✅ WebRTC clients | ❌ Desktop only |
| **Local model** | ✅ Any provider | ✅ Ollama, LMStudio, vLLM | ✅ Plugin system | ❌ Gemini only |
| **Multi-agent routing** | ✅ Profiles + delegation | ✅ Channel → agent routing | ✅ Agent handoff | ❌ Single agent |
| **Session flexibility** | ✅ Turn-based + cron | ✅ Turn + streaming Talk | ✅ Streaming only | ✅ Streaming only |
| **Plugin SDK** | ✅ Python | ✅ TypeScript | ✅ Python plugins | ❌ None |
| **Self-hosted** | ✅ Fully | ✅ Fully | ✅ Fully | ✅ Fully |

---

## The Answer

### Neither alone. Both together. Here's how:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         LEGACY AI RUNTIME                                │
│                                                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │              OPENCLAW GATEWAY (always-on, TypeScript)              │  │
│  │                                                                    │  │
│  │  ├── Talk Mode (voice I/O, interrupt detection, streaming)        │  │
│  │  ├── Realtime Voice (WebRTC, provider-websocket, gateway-relay)   │  │
│  │  ├── Device Nodes (glasses, phone, car — paired via WS)           │  │
│  │  ├── Canvas (live HUD/dashboard)                                  │  │
│  │  ├── TTS (ElevenLabs, MLX local, system)                         │  │
│  │  ├── STT (Deepgram realtime, system, Whisper)                     │  │
│  │  └── Agent Brain: calls Hermes via internal API                   │  │
│  │                      │                                             │  │
│  └──────────────────────┼────────────────────────────────────────────┘  │
│                          │                                               │
│                    Internal API (HTTP/NATS)                              │
│                          │                                               │
│  ┌──────────────────────┼────────────────────────────────────────────┐  │
│  │              HERMES AGENT (reasoning, Python)                      │  │
│  │                                                                    │  │
│  │  ├── Memory (persistent, cross-session, injected every turn)      │  │
│  │  ├── Skills (100+, loadable, self-updating)                       │  │
│  │  ├── Tools (NATS, HA, SaaS, filesystem, web, code_exec)          │  │
│  │  ├── Session history + search (FTS5)                              │  │
│  │  ├── Model routing (local Qwen for fast, Claude for deep)         │  │
│  │  ├── Watch system (background NATS monitoring)                    │  │
│  │  ├── Cron jobs (scheduled tasks)                                  │  │
│  │  └── Delegation (subagents for parallel work)                     │  │
│  │                                                                    │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                         NATS (nervous system)                      │  │
│  │                                                                    │  │
│  │  ├── cmd.ha.*         → Home Assistant                            │  │
│  │  ├── cmd.drone.*      → Fleet control                            │  │
│  │  ├── events.ha.*      → Home state changes                       │  │
│  │  ├── perception.*     → Vision/audio feeds                       │  │
│  │  ├── context.*        → Location, presence                       │  │
│  │  └── watches.*        → Background monitoring                    │  │
│  │                                                                    │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Why This Split Works

| Layer | OpenClaw handles | Hermes handles |
|---|---|---|
| **Interface** | Voice streaming, WebRTC, device pairing, audio pipeline | — |
| **Presentation** | Canvas HUD, UI overlays, TTS voice selection | — |
| **Reasoning** | — | LLM calls, tool use, decision making |
| **Memory** | — | Persistent cross-session, skills, context |
| **Tools/Actions** | — | NATS commands, HA control, SaaS, code exec |
| **Background work** | — | Watches, cron jobs, delegation |

**OpenClaw = the mouth, ears, and eyes.**
**Hermes = the brain, memory, and hands.**
**NATS = the nervous system connecting everything.**

### The Bridge: How They Connect

OpenClaw's `agent-consult` brain mode already supports calling an external agent for reasoning. The bridge:

**Option A: OpenClaw calls Hermes via HTTP API**
```
OpenClaw Talk → user speaks → STT → text
    → OpenClaw sends to Hermes API (localhost)
    → Hermes reasons (memory, tools, skills)
    → Returns response text
    → OpenClaw TTS → user hears
```

Hermes already has an HTTP API mode (the WebUI server at 8787). OpenClaw can call it as its "brain."

**Option B: OpenClaw calls Hermes via NATS**
```
OpenClaw Talk → user speaks → STT → text
    → nats_publish("cmd.agent.legacy.prompt", {text, context})
    → Hermes subscribes, processes with full memory/tools
    → nats_publish("cmd.agent.legacy.response", {text, actions})
    → OpenClaw receives, TTS → user hears
```

This keeps everything on the NATS fabric. Hermes subscribes to `cmd.agent.legacy.prompt` and publishes responses.

**Option C: OpenClaw plugin that wraps Hermes Python process**
An OpenClaw extension that spawns/connects to Hermes as a subprocess or RPC target. OpenClaw handles all I/O, Hermes handles all reasoning.

### The "Room" Pattern (LiveKit-style, via OpenClaw)

OpenClaw's Talk mode already IS the room:
- You "join" by activating Talk mode (menu bar toggle, or device node connects)
- The session is persistent and streaming
- Multiple devices can connect (phone node, glasses node, macOS app)
- The agent is always listening (with interrupt detection)
- Tool calls route through the agent brain (Hermes)

**You don't need LiveKit separately.** OpenClaw's gateway IS the LiveKit equivalent for your use case — it handles WebRTC, audio streaming, device pairing, and multi-node. LiveKit adds value only if you need PSTN telephony or massive multi-party rooms (which you don't).

### Session Flexibility (Turn-based + Streaming)

This architecture gives you BOTH:

| Mode | When | How |
|---|---|---|
| **Streaming (Talk)** | Voice conversation, real-time, continuous | OpenClaw Talk → Hermes reasoning |
| **Turn-based (Text)** | WebUI, complex tasks, code review | Direct to Hermes (current WebUI) |
| **Background (Watches)** | Monitoring, proactive | Hermes watches → OpenClaw TTS notification |
| **Scheduled (Cron)** | Reports, recurring tasks | Hermes cron → result to OpenClaw for delivery |

Nothing is forced. Voice is streaming. Text is turn-based. Background work runs independently. All share the same memory and context.

---

## Existing Projects / Prior Art

| Project | What it does | Relevant? |
|---|---|---|
| **[LiveKit Agents](https://github.com/livekit/agents)** | WebRTC room-based voice AI framework | ✅ Pattern reference (room concept) but OpenClaw already does this |
| **[Mark-XXXIX](https://github.com/FatihMakes/Mark-XXXIX)** | Gemini LiveConnect desktop assistant | ✅ UI reference (PyQt6 overlay) but no infra/memory |
| **[OpenClaw Talk](https://docs.openclaw.ai)** | Always-on voice + multi-device + WebRTC | ✅ THIS IS THE SOLUTION for I/O layer |
| **[Hermes Agent](https://github.com/NousResearch/hermes-agent)** | Memory + skills + tools + NATS | ✅ THIS IS THE SOLUTION for brain layer |
| **[Pipecat](https://github.com/pipecat-ai/pipecat)** | Python voice AI framework (similar to LiveKit) | Partial — simpler than LiveKit, no memory |
| **[Vocode](https://github.com/vocodedev/vocode-core)** | Voice agent platform | Partial — telephony focused, not personal assistant |
| **[Bolna](https://github.com/bolna-ai/bolna)** | Voice AI agent framework | Partial — call-center focused |
| **[Letta (MemGPT)](https://github.com/letta-ai/letta)** | LLM with persistent memory | Partial — memory only, no voice/device |

**Conclusion: No single existing project does what we need. But OpenClaw + Hermes together cover 95% — we just need the bridge.**

---

## Implementation Plan

### Phase 1: Bridge OpenClaw → Hermes (This Week)

Build a NATS-based bridge so OpenClaw's Talk mode uses Hermes as its brain:

1. **Hermes NATS listener** — Python service subscribing to `cmd.agent.legacy.prompt`
   - Receives: `{text, session_id, context}`
   - Processes with full memory, skills, tools
   - Publishes response to `cmd.agent.legacy.response`

2. **OpenClaw custom brain plugin** — TypeScript extension that:
   - On user speech → publishes to NATS
   - Subscribes to response subject
   - Returns text to TTS pipeline

3. **Test loop:** Talk mode → speak → Hermes reasons with memory → response spoken back

### Phase 2: Watch Notifications via OpenClaw TTS (Week 2)

When Hermes watch system triggers:
- Publishes to `cmd.tts.speak` on NATS
- OpenClaw subscribes and speaks the notification through active Talk session
- Respects interrupt/queue rules

### Phase 3: Device Nodes + Context (Week 3-4)

- Pair iPhone as OpenClaw node (already supported)
- Camera frames from glasses → NATS → Hermes perception
- Location context from device → NATS KV → Hermes context injection

### Phase 4: Local Model for Fast Interface (Month 2)

- Qwen2.5-Omni via local OpenAI-compatible server
- OpenClaw routes simple/fast queries to local model
- Complex queries route to Hermes → Claude/GPT-5
- Same Talk session, transparent routing

---

## Key Decisions

1. **OpenClaw for I/O, Hermes for brain** — don't rebuild voice/device infra, don't rebuild memory/tools
2. **NATS as the bridge** — not HTTP, stays on the fabric
3. **No Telegram/Discord** — stripped from architecture (per Zeus directive)
4. **Turn-based AND streaming coexist** — not forced into one mode
5. **OpenClaw's Talk mode IS the "room"** — no need for LiveKit separately
6. **Plugin SDK on both sides** — OpenClaw extension for NATS brain, Hermes plugin for watch→TTS
7. **Mark-XXXIX as UI inspiration only** — the PyQt6 overlay aesthetic, not the architecture

---

## Related

- [[OneMind-Realtime-Core]] — the vision document (always-on model, watches, omnipresence)
- [[HA-NATS-Bridge]] — existing NATS pattern (proven, sub-5ms)
- [[NATS]] — central nervous system
- [[NATS-Intent-Mesh]] — subject namespace

---

## Addendum: Model Flexibility & LiveKit Self-Hosting (2026-05-25)

### Zeus's Concern: "OpenClaw is limited to models they choose"

**This is partially true for REALTIME VOICE specifically:**

OpenClaw's realtime voice providers are **hard-coded to specific APIs**:
- `openai` → OpenAI Realtime API (gpt-realtime-2) — WebRTC
- `google` → Google Gemini LiveConnect — provider-websocket
- That's it for full-duplex native audio models

For **text inference** (the brain), OpenClaw supports ANY model:
- Ollama (local), LM Studio, vLLM, sglang, LiteLLM
- OpenRouter (hundreds of models)
- Anthropic, OpenAI, Groq, DeepSeek, Mistral, xAI, etc.
- Custom providers via plugin SDK

For **TTS/STT separately** (non-realtime Talk mode):
- STT: Deepgram, system (Apple), Whisper (via Ollama)
- TTS: ElevenLabs, Azure, MLX local, system
- These can be mixed with any text LLM

**The limitation:** If you want Qwen2.5-Omni running locally doing voice-to-voice natively (audio in → model → audio out), OpenClaw has no adapter for that. Its realtime voice assumes OpenAI or Google's proprietary APIs.

**To use a custom local model with full-duplex audio**, you'd need either:
1. Write an OpenClaw realtime voice provider plugin that wraps your local model's audio API
2. Use OpenClaw in "native Talk" mode (separate STT → text LLM → TTS) with local Whisper + local LLM + local TTS
3. Bypass OpenClaw for the audio path entirely and use LiveKit or a custom solution

### LiveKit Self-Hosting Complexity

**What you'd need to run:**

| Component | What | Complexity |
|---|---|---|
| `livekit-server` | Go binary, single process | ✅ Easy — just run the binary |
| Port requirements | TCP 7880 (API), UDP 50000-60000 (WebRTC media) | ⚠️ Need open UDP range |
| SSL/TURN | Required for production (not local dev) | Medium — Caddy/nginx + coturn |
| Agent process | Python, connects to server | ✅ Easy — `python agent.py start` |
| Client SDK | JS/Swift/Kotlin/Flutter apps | You'd need to build the client |

**Local dev (what you'd start with):**
```bash
# Install
brew install livekit

# Run server (local, no SSL needed)
livekit-server --dev

# Run agent (connects locally)
python my_agent.py dev
```

**That's it for local.** Server is a single binary. Agent is a Python script. Client connects via WebSocket + WebRTC.

**Production adds:**
- SSL certificates (Let's Encrypt via Caddy)
- TURN server for NAT traversal (coturn or LiveKit's built-in TURN)
- Redis for multi-instance coordination (optional, only for scaling)
- Kubernetes helm chart available for full deployment

**Self-hosted vs Cloud comparison (from their docs):**
- Self-hosted: up to ~3,000 users per room, single-home SFU, you manage it
- Cloud: unlimited, global mesh, managed — but costs money and sends data to them

**Complexity rating for your use case (single user, local):**
- **LOW** — literally one binary + one Python agent process
- No Redis needed (single room, single user)
- No SSL needed (all local/Tailscale)
- No TURN needed (same LAN or Tailscale direct)
- Console mode exists: `python agent.py console` — no server at all, just local mic/speaker

### LiveKit Agents: What It Actually Gives You

```python
from livekit.agents import Agent, AgentSession, AgentServer, function_tool, cli
from livekit.plugins import silero, deepgram, openai, cartesia

server = AgentServer()

@server.rtc_session()
async def entrypoint(ctx: JobContext):
    session = AgentSession(
        vad=silero.VAD.load(),           # Voice activity detection
        stt=deepgram.STT(model="nova-3"), # Speech-to-text
        llm=openai.LLM(model="gpt-4.1"), # ANY OpenAI-compatible endpoint
        tts=cartesia.TTS(),               # Text-to-speech
    )
    
    agent = Agent(
        instructions="You are Legacy AI.",
        tools=[control_home, check_drones, watch_door],
    )
    
    await session.start(agent=agent, room=ctx.room)
```

**Key:** The `llm=` parameter accepts ANY OpenAI-compatible API endpoint. You can point it at:
- Local Qwen2.5 via llama.cpp server
- Local vLLM
- Groq (fast cloud)
- Any custom endpoint

**LiveKit's MCP support:**
```python
from livekit.agents.llm import MCPServerTool

# Connect to any MCP server
mcp_tools = MCPServerTool(server_url="http://localhost:8080")
agent = Agent(tools=[mcp_tools])
```

**Multi-agent handoff:**
```python
class MonitorAgent(Agent):
    @function_tool
    async def escalate_to_reasoning(self, context, issue):
        """Hand off to deep reasoning agent"""
        return self.session.handoff(ReasoningAgent())
```

### Mark-XXXIX Codebase Analysis

**11,060 lines of Python total.** Here's what it actually has:

| Component | Lines | What it does |
|---|---|---|
| `main.py` | 880 | Gemini LiveConnect session, audio streaming, tool dispatch |
| `ui.py` | 1,503 | PyQt6 desktop UI — orb animation, chat panel, transparency |
| `agent/executor.py` | ~200 | Generates Python code via Gemini, runs it |
| `agent/planner.py` | ~150 | Multi-step task decomposition (max 5 steps) |
| `agent/task_queue.py` | ~100 | Priority queue with threading |
| `agent/error_handler.py` | ~100 | Error analysis and retry logic |
| `memory/memory_manager.py` | 214 | JSON file — categories: identity, preferences, projects, relationships |
| `actions/browser_control.py` | 892 | Playwright browser automation |
| `actions/screen_processor.py` | 441 | Screenshot → Gemini vision analysis |
| `actions/computer_control.py` | 499 | System commands (OS-specific) |
| `actions/file_controller.py` | 542 | File CRUD operations |
| `actions/code_helper.py` | 568 | Code generation and execution |
| `actions/dev_agent.py` | 596 | Multi-step coding tasks |
| 11 more action files | ~3,000 | Weather, flights, games, YouTube, etc. |

**Architecture pattern:**
```
Gemini LiveConnectConfig (native audio model)
    ↕ bidirectional audio stream (sounddevice)
    ↕ tool_calls come back in the stream
    → dispatch to actions/*.py
    → results fed back into session
    
Memory: JSON file (2200 char max, auto-trimming)
UI: PyQt6 overlay (always-on-top, transparent, orb visualizations)
Agent: Simple planner (Gemini generates step list) + executor (runs steps sequentially)
```

**What it DOESN'T have:**
- No NATS or message bus
- No persistent session history (just current conversation)
- No skills system
- No background monitoring/watches
- No multi-device support
- No custom model routing (Gemini only)
- No tool ecosystem beyond basic actions
- No cron/scheduled tasks
- No delegation/subagents

**What's USEFUL from it:**
1. The UI aesthetic (PyQt6 transparent overlay with orb animation)
2. The Gemini LiveConnect pattern (native audio in/out + tool calling in one stream)
3. The simplicity of the tool dispatch pattern
4. Proof that ~1000 lines of Python can do real-time voice + vision + tools

### Revised Architecture Decision

Given:
1. OpenClaw's realtime voice is locked to OpenAI/Google realtime APIs
2. LiveKit is trivial to self-host (one binary, local dev = zero config)
3. LiveKit agents support ANY OpenAI-compatible LLM endpoint
4. LiveKit has native MCP support
5. Mark-XXXIX shows the Gemini LiveConnect pattern works with minimal code

**Revised recommendation:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         LEGACY AI RUNTIME                                │
│                                                                          │
│  OPTION A: OpenClaw (for non-realtime voice)                            │
│  ├── Native Talk mode: System STT → Hermes brain → ElevenLabs TTS      │
│  ├── Works TODAY with any model                                          │
│  ├── macOS/iOS/Android nodes                                             │
│  └── Good for ambient/slow interactions                                  │
│                                                                          │
│  OPTION B: LiveKit (for true realtime streaming)                        │
│  ├── Self-hosted server (one binary, local)                              │
│  ├── Python agent with Whisper STT + local Qwen + TTS                   │
│  ├── Room concept = "join and it streams"                               │
│  ├── Full control over model choice                                      │
│  ├── MCP tools → connects to Hermes tools via MCP or NATS              │
│  └── Good for active voice conversation + low-latency                   │
│                                                                          │
│  OPTION C: Gemini LiveConnect (for native audio, simplest)              │
│  ├── One bidirectional audio stream                                      │
│  ├── Model handles STT+reasoning+TTS in one call                        │
│  ├── Lowest latency possible (no pipeline)                              │
│  ├── Limited to Gemini or OpenAI realtime                               │
│  └── Good for pure voice conversation without complex tools             │
│                                                                          │
│  SHARED BRAIN (always Hermes):                                          │
│  ├── Memory, skills, tools, NATS, watches, sessions                     │
│  └── Called via NATS or MCP from any interface above                    │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**The "middle man" is not one thing — it's a model router:**

| Interaction type | Best path | Why |
|---|---|---|
| Quick voice command ("turn off lights") | Gemini/OpenAI Realtime | Lowest latency, native audio |
| Complex reasoning ("plan drone patrol") | Local Qwen → Hermes tools | Needs memory, tools, NATS |
| Background notification ("Sarah's home") | Hermes watch → TTS out | Proactive, no user input |
| Extended conversation | LiveKit room + Hermes brain | Full streaming + full tools |
| Text/code work | Hermes direct (WebUI) | Turn-based is better here |

### The Actual Build Path

**Week 1: LiveKit local + Hermes brain**
```bash
brew install livekit
livekit-server --dev

# Agent connects Hermes via NATS
python legacy_voice_agent.py dev
```

Agent code (~200 lines):
- LiveKit STT (Deepgram or Whisper)
- LLM = Hermes via NATS (publish prompt, get response)
- TTS (ElevenLabs or local)
- Tools = NATS requests to ha-bridge, drone controller, etc.

**Week 2: OpenClaw integration for device nodes**
- OpenClaw handles macOS Talk overlay + iOS/Android nodes
- Brain plugin routes to Hermes via NATS
- Same memory, same session, different I/O

**Week 3: Model router**
- Fast queries → local Qwen (sub-500ms)
- Complex queries → Claude/GPT-5 via Hermes
- Voice commands → Gemini Realtime (if simple) or LiveKit pipeline (if needs tools)

This gives you FULL model flexibility, not locked to any provider's realtime API.
