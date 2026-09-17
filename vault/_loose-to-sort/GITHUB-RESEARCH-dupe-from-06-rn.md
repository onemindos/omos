# GITHUB-RESEARCH.md — Ambient Presence Stack Research

> Researched 14 GitHub repositories against the **LEGACY_AMBIENT_PRESENCE_MASTERPLAN.md**.
> Date: 2026-03-14

---

## Executive Summary

Of 14 repos evaluated, **5 are strong candidates for immediate integration**, 3 are valuable reference resources, and 6 are lower priority (inspiration, already deployed, or tangential).

### Top 5 — Install / Fork / Deploy

| Priority | Repo | Stars | Verdict | Layer |
|----------|------|-------|---------|-------|
| 🔴 P0 | **lossless-claw** | 1.9k | Install NOW — replaces OpenClaw's sliding-window compaction | Memory |
| 🔴 P0 | **OpenViking** | 12.7k | Deploy as context database — +46% accuracy, -91% token cost | Memory |
| 🟠 P1 | **obsidian-skills** | 14.1k | Install — made by kepano (Obsidian creator), 5 skills for vault ops | Skills |
| 🟠 P1 | **ClawX** | 4.5k | Fork — desktop foundation, add MediaPipe + ambient overlay | Desktop |
| 🟡 P2 | **ClawVault** | 584 | Evaluate — structured memory with Obsidian integration + OpenClaw hooks | Memory |

### Reference Resources (Bookmark)

| Repo | Stars | Use |
|------|-------|-----|
| **awesome-openclaw-skills** | 37.9k | Skill discovery — 5,400+ curated skills across 30+ categories |
| **openclaw-api-list** | 126 | API discovery — 10,498+ APIs with recommended starter list |
| **chowder-iOS** | 541 | iOS client reference — compare against our ios-app architecture |

### Already Deployed / Low Priority

| Repo | Stars | Status |
|------|-------|--------|
| **swarmclawai/swarmclaw** | 140 | Already deployed (v0.7.7 on VPS). Now at v1.1.5 — significant upgrade available |
| **HumeAI/hume-python-sdk** | 165 | Emotion AI SDK — future voice layer enhancement, not immediate |
| **ClawdAgent** | 9 | Massive scope (73k lines) but only 9 stars, 2 contributors. Inspiration only |
| **ada_local** | 128 | Local voice assistant — different architecture, inspiration for local pipeline |
| **Mission Control** | 2.6k | Dashboard with 32 panels — impressive but SwarmClaw already deployed; revisit if SwarmClaw stalls |

---

## Detailed Evaluations

### 🔴 P0 — Lossless Claw (martian-engineering/lossless-claw)

**What**: OpenClaw plugin that replaces sliding-window context compaction with a DAG-based summarization system. Nothing is ever lost.

**Stars**: 1.9k | **License**: MIT | **Language**: TypeScript 70.8%, Go 28.7%

**How it works**:
1. Persists every message in SQLite, organized by conversation
2. Summarizes chunks of older messages into summaries via configured LLM
3. Condenses summaries into higher-level nodes as they accumulate (DAG)
4. Assembles context each turn: summaries + recent raw messages
5. Provides tools (`lcm_grep`, `lcm_describe`, `lcm_expand`) for agents to search/recall

**Why this matters for Legacy**:
- Legacy's PRODUCTION-PLAN.md identified "context death" as the #1 problem
- Current OpenClaw compaction = lossy truncation. Sessions reset and context vanishes
- LCM means Legacy never forgets a conversation. Raw messages stay in DB forever
- Agent tools let Legacy drill into any summary to recover original detail
- Compatible with the existing OpenClaw plugin system (`openclaw plugins install`)

**Install**:
```bash
openclaw plugins install @martian-engineering/lossless-claw
```

**Config** (add to openclaw.json):
```json
{
  "plugins": {
    "slots": { "contextEngine": "lossless-claw" },
    "entries": {
      "lossless-claw": {
        "enabled": true,
        "config": {
          "freshTailCount": 32,
          "contextThreshold": 0.75,
          "incrementalMaxDepth": -1
        }
      }
    }
  },
  "session": {
    "reset": {
      "mode": "idle",
      "idleMinutes": 10080
    }
  }
}
```

**Risk**: Low. It's a plugin — can be disabled without affecting core. MIT licensed. Active development (v0.3.0, 16 contributors).

**Master Plan Alignment**: Directly addresses Phase 1 "Memory Foundation" requirement. Install before any ambient presence work.

---

### 🔴 P0 — OpenViking (volcengine/OpenViking)

**What**: Open-source Context Database designed specifically for AI Agents. Replaces flat vector storage with a filesystem-paradigm hierarchical context system.

**Stars**: 12.7k | **License**: Apache 2.0 | **Language**: Python 82.8%, C++ 8.4%, Rust 3.6%

**Core Concepts**:
1. **Filesystem paradigm**: `viking://` protocol — memories, resources, skills are virtual directories
2. **Tiered context loading** (L0/L1/L2): Abstract → Overview → Full detail, loaded on demand
3. **Directory recursive retrieval**: Lock high-score directory → refine → recurse subdirectories
4. **Visualized retrieval trajectory**: Observable context debugging
5. **Automatic session management**: Memory self-iteration at session end

**OpenClaw Memory Plugin Benchmarks**:

| Configuration | Accuracy | Input Tokens |
|---|---|---|
| OpenClaw native (memory-core) | 35.65% | 24,611,530 |
| OpenClaw + LanceDB | 44.55% | 51,574,530 |
| **OpenClaw + OpenViking (-memory-core)** | **52.08%** | **4,264,396** |
| **OpenClaw + OpenViking (+memory-core)** | **51.23%** | **2,099,622** |

- **+46% accuracy** over native OpenClaw
- **-91% token cost** with native memory enabled
- **-83% token cost** with native memory disabled

**Why this matters for Legacy**:
- Current memory: hybrid BM25+vector with text-embedding-3-small
- OpenViking's hierarchical filesystem approach maps perfectly to the codex structure
- The L0/L1/L2 tiered loading could dramatically reduce heartbeat token costs
- VikingBot (built-in agent framework) could serve as the "brain" for ambient state queries
- Has OpenClaw memory plugin with proven benchmarks

**Integration Path**:
```bash
pip install openviking --upgrade --force-reinstall
openviking-server --with-bot
```

**Consideration**: Runs as a separate server (Python). Needs VPS resources. Consider running alongside NATS on the droplet. The OpenClaw plugin integrates via the plugin system.

**Master Plan Alignment**: Phase 1 "Memory Foundation" — this IS the next-gen memory layer. Could replace or complement current BM25+vector setup.

---

### 🟠 P1 — Obsidian Skills (kepano/obsidian-skills)

**What**: Official agent skills for Obsidian, made by **kepano** (Steph Ango — CEO/creator of Obsidian). 5 skills that teach AI agents to properly work with Obsidian vault files.

**Stars**: 14.1k | **License**: MIT | **Contributors**: 12 (including kepano himself)

**Skills included**:
| Skill | Purpose |
|---|---|
| `obsidian-markdown` | Create/edit Obsidian Flavored Markdown (.md) with wikilinks, embeds, callouts, properties |
| `obsidian-bases` | Create/edit Obsidian Bases (.base) with views, filters, formulas, summaries |
| `json-canvas` | Create/edit JSON Canvas files (.canvas) with nodes, edges, groups |
| `obsidian-cli` | Interact with Obsidian vaults via CLI including plugin/theme development |
| `defuddle` | Extract clean markdown from web pages (token-saving scraper) |

**Why this matters for Legacy**:
- Legacy already has an `obsidian-1.0.0` skill, but kepano's skills are the canonical reference
- The `obsidian-markdown` skill would improve how Legacy writes to the codex (proper wikilinks, callouts, properties)
- The `defuddle` skill saves tokens when fetching web content (strips clutter)
- `obsidian-cli` could enable Legacy to manage Obsidian plugins/themes on the VPS

**Install**:
```bash
# Copy to OpenClaw skills directory
git clone https://github.com/kepano/obsidian-skills.git
cp -r obsidian-skills/skills/* ~/.openclaw/workspace/skills/
```

Or via ClawHub: `clawhub install kepano/obsidian-skills`

**Risk**: Zero. These are read-only skill definitions — they teach the agent proper syntax, not execute code.

**Master Plan Alignment**: Enhances vault operations quality. Phase 1 "Codex Integration" improvement.

---

### 🟠 P1 — ClawX (ValueCell-ai/ClawX)

**What**: The Desktop Interface for OpenClaw AI Agents. Electron 40+ / React 19 / Tailwind / Zustand / Vite. Embeds an OpenClaw runtime directly.

**Stars**: 4.5k | **License**: MIT | **Releases**: 116 (v0.2.3)

**Features**:
- Zero-config setup wizard
- Intelligent chat with `@agent` routing
- Multi-channel management
- Cron automation
- Extensible skill system (pre-bundles document processing, search skills)
- Secure provider integration (keychain)
- Proxy settings
- Multi-platform (macOS/Windows/Linux)

**Why this matters for Legacy**:
- The Master Plan proposes building a desktop command center from scratch (Electron + React + MediaPipe)
- **ClawX already IS an Electron + React OpenClaw desktop app**
- Instead of building from zero, fork ClawX and add:
  - MediaPipe hand tracking overlay
  - Ambient state panel (AMBIENT_STATE.md viewer)
  - SmolVLM-256M vision integration
  - NATS connection for live sensor data
- 116 releases = actively maintained, battle-tested
- Already handles OpenClaw gateway connection, chat, skills, cron

**Fork Strategy**:
1. Fork ClawX
2. Add `legacy-eyes` camera panel with MediaPipe Hands
3. Add ambient overlay (transparent always-on-top)
4. Add NATS subscriber panel for live events
5. Add gesture recognition → command mapping
6. Keep upstream sync for core updates

**Risk**: Medium. Forking means maintenance burden. But starting from 4.5k-star foundation > building from scratch.

**Master Plan Alignment**: Phase 3 "Desktop Command Center" — this IS the foundation. Fork don't build.

---

### 🟡 P2 — ClawVault (Versatly/clawvault)

**What**: Persistent Memory for AI Agents. Structured markdown-based memory system with 8 core primitives and Obsidian integration.

**Stars**: 584 | **License**: MIT | **Tests**: 466

**Core Features**:
- 8 primitives: Goals, Agents, State Space, Feedback, Capital, Institution, Synthesis, Recursion
- Markdown-native vault with wiki-links
- Hybrid search: BM25 + semantic embeddings + cross-encoder reranking
- Fact store with conflict resolution
- Entity graph for multi-hop queries
- Session checkpoint/recover
- Context profiles (default/planning/incident/handoff/auto)
- Obsidian integration (graph themes, canvas dashboards, Kanban round-trip)
- Tailscale + WebDAV sync
- OpenClaw hook integration (`openclaw hooks install clawvault`)
- Background workers: Curator, Janitor, Distiller, Surveyor

**Why this matters**:
- Maps to existing Obsidian vault (codex/) and Tailscale infrastructure
- OpenClaw hooks = native integration without custom code
- Cross-encoder reranking is more sophisticated than current BM25+vector
- Entity graph enables multi-hop reasoning ("who did Zeus meet last week who works at X?")
- Background workers (Curator, Janitor) could replace manual vault maintenance

**Consideration**: With both OpenViking AND Lossless Claw, ClawVault may be redundant for memory. But its Obsidian integration and entity graph are unique. Evaluate after P0 items are installed.

**Master Plan Alignment**: Phase 1 memory enhancement. Lower priority than OpenViking/LCM because those have stronger benchmarks, but the Obsidian native integration is compelling.

---

### 📚 Reference — Awesome OpenClaw Skills (VoltAgent/awesome-openclaw-skills)

**Stars**: 37.9k | **License**: MIT

The **#1 community resource** for OpenClaw skills. 5,400+ curated skills from 13,729+ on ClawHub, organized into 30+ categories:

| Category | Count | Relevant Skills for Legacy |
|---|---|---|
| Coding Agents & IDEs | 1,222 | Claude Code task bridge, agent orchestration |
| Web & Frontend | 938 | Browser automation, scraping |
| DevOps & Cloud | 409 | Server management, Docker, SSH |
| Search & Research | 352 | Academic research, web search |
| Browser & Automation | 335 | Form filling, CAPTCHA solving |
| Productivity & Tasks | 206 | Task management, planning |
| AI & LLMs | 197 | Model management, prompt engineering |
| Communication | 149 | Messaging, notifications |
| Notes & PKM | 71 | Obsidian, knowledge management |
| Smart Home & IoT | 43 | Home Assistant, device control |
| Speech & Transcription | 45 | Voice, TTS, STT |

**Action**: Bookmark. When adding new capabilities to Legacy, search this list FIRST before building custom skills.

**Specific skills to investigate**:
- `active-maintenance` — automated system health for OpenClaw
- `agent-team-orchestration` — multi-agent workflows
- `arc-agent-lifecycle` — agent lifecycle management
- `agent-daily-planner` — structured daily planning for agents

---

### 📚 Reference — OpenClaw API List (cporter202/openclaw-api-list)

**Stars**: 126 | 10,498+ APIs catalogued

Massive curated list of APIs usable with OpenClaw. Has a `OPENCLAW_RECOMMENDED.md` with ~100 starter APIs. Categories: Agents, AI, scrapers, MCP servers, etc.

**Action**: Bookmark. Reference when expanding Legacy's API integrations.

---

### 📚 Reference — Chowder iOS (newmaterialco/chowder-iOS)

**What**: Native iOS chat client for OpenClaw. Swift, 541 stars, MIT.

**Stars**: 541 | **Language**: Swift 100%

**Key Features**:
- WebSocket → OpenClaw Gateway Protocol v3
- Real-time streaming with thinking steps + tool activity
- Live Activity tracking (polls `chat.history` every 500ms)
- Agent identity sync (IDENTITY.md, USER.md)
- Custom agent avatar
- Auto-reconnect with 3s backoff
- Tailscale connectivity

**Why it matters**:
- We already have `command-center/ios-app/` with similar goals
- Chowder is further along in protocol implementation (especially thinking steps + Live Activity)
- **Compare architecture**: Chowder's ChatService.swift vs our ChatService patterns
- The Live Activity implementation is exactly what the Master Plan describes for "ambient iOS presence"

**Action**: Study Chowder's protocol handling and Live Activity code. Consider adopting their patterns for our ios-app, or evaluate replacing our ios-app with a Chowder fork.

---

### 📋 Already Deployed — SwarmClaw (swarmclawai/swarmclaw)

**Current**: v0.7.7 on VPS | **Latest**: v1.1.5 (76 releases!)

**Major updates since our deployment**:
- v1.0.0: OpenClaw setup wizard, npm install support
- v1.1.0: Mission controller, Autonomy safety desk, durable queued follow-ups
- v1.1.1: Structured Sessions from context (chatrooms, tasks, missions), ProtocolRun orchestration
- v1.1.2: Nebius + DeepInfra as built-in providers, Ollama local/cloud routing
- v1.1.4: Orchestrator agents, parallel task execution, restart-safe swarm state

**New Capabilities since v0.7.7**:
- Missions (durable multi-step objectives)
- Orchestrator agents (wake cycles, governance, autonomy controls)
- Structured Sessions (bounded runs with templates, facilitators, transcripts)
- Conversation-to-skill drafts
- Chatroom breakout flows
- Provider health persistence across restarts

**Action**: **Upgrade SwarmClaw on VPS from v0.7.7 → v1.1.5**. This is a significant capability jump that brings missions, orchestration, and structured sessions.

---

### 📋 Future — Hume AI Python SDK

**Stars**: 165 | **License**: MIT | Python SDK for emotion AI

**Products**:
- **Empathic Voice Interface (EVI)**: Voice with emotional understanding
- **Text-to-Speech (TTS)**: Emotionally expressive speech
- **Expression Measurement**: Facial expression + prosody analysis

**Why it matters for future**:
- The Master Plan mentions upgrading voice with emotional intelligence
- Hume's EVI could replace or complement OpenAI TTS for emotionally-aware responses
- Expression Measurement could feed into AMBIENT_STATE.md (Legacy reads Zeus's emotional state from webcam)
- Prosody analysis on voice input = Legacy understands HOW Zeus says things, not just WHAT

**Action**: Phase 4+ consideration. After ambient presence is stable, evaluate Hume EVI integration for emotional voice. Not immediate priority.

---

### 📋 Inspiration — ClawdAgent (liortesta/ClawdAgent)

**Stars**: 9 | **License**: Apache 2.0 | 73,000+ lines TypeScript

Extremely ambitious autonomous agent system with 20 agents, 35 tools, 67 skills, 14-layer security, self-evolution, crypto trading, social publishing, and A2A/ACP/MCP protocols. Only 2 contributors (author + Claude).

**Interesting patterns to study**:
- 14-layer Zero Trust security architecture
- SHA-256 memory integrity checksums
- Self-evolution engine concept
- Intent classification (45+ intents)
- Smart model routing by complexity/cost/budget
- PostgreSQL + Redis + BullMQ architecture

**Why NOT to use directly**:
- 9 stars despite massive feature claims suggests very new, unvalidated
- Building for a completely different architecture (PostgreSQL + Redis vs OpenClaw + NATS)
- "Kitchen sink" approach — tries to do everything
- 2 contributors = sole maintainer risk

**Action**: Study specific patterns (security layers, model routing logic) as inspiration. Do not install or fork.

---

### 📋 Inspiration — ada_local (nazirlouis/ada_local)

**Stars**: 128 | **License**: MIT | Python

Local voice assistant: PySide6 GUI, Ollama LLM, Piper TTS, Whisper STT, wake word "Jarvis". Interesting for the fully-local pipeline concept.

**Relevant patterns**:
- Piper TTS (local, fast, free) — alternative to Chatterbox mentioned in Master Plan
- FunctionGemma router for intent classification
- TP-Link Kasa integration (we use Home Assistant)

**Action**: Reference only. Our architecture is fundamentally different (OpenClaw gateway, not standalone Python app). But Piper TTS is worth evaluating alongside Chatterbox.

---

### 📋 Already Covered — Mission Control (builderz-labs/mission-control)

**Stars**: 2.6k | **License**: MIT | Next.js 16 + React 19 + SQLite

32-panel dashboard for AI agent orchestration. Very impressive:
- Ed25519 device identity
- Agent SOUL system
- Skills Hub (ClawHub + skills.sh integration)
- Memory knowledge graph
- Claude Code session tracking
- GitHub Issues sync
- Framework adapters (OpenClaw, CrewAI, LangGraph, AutoGen)
- Agent eval framework (4-layer)
- Flight Deck (Tauri v2 desktop companion — private beta)

**Why not replacing SwarmClaw**:
- SwarmClaw already deployed and configured with our 9 agents
- SwarmClaw v1.1.5 now has missions, orchestration, structured sessions
- Mission Control would require full migration
- Both are Next.js + SQLite — similar tech stack
- Mission Control's agent eval framework is interesting but not critical right now

**Action**: Bookmark. If SwarmClaw development stalls or we need multi-framework support (CrewAI/LangGraph adapters), revisit Mission Control as replacement.

---

## Integration Architecture — How Everything Fits

```
┌─────────────────────────────────────────────────────────────────────┐
│                        LEGACY'S BRAIN                                │
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────────────┐  │
│  │ Lossless Claw│  │  OpenViking  │  │   ClawVault (optional)    │  │
│  │ (context     │  │  (context DB │  │   (structured memory      │  │
│  │  engine      │  │   + retrieval│  │    + entity graph         │  │
│  │  plugin)     │  │   server)    │  │    + Obsidian hooks)      │  │
│  │              │  │              │  │                           │  │
│  │ DAG summaries│  │ viking://    │  │ wiki-links, reranking    │  │
│  │ lcm_grep     │  │ L0/L1/L2    │  │ fact conflict resolution  │  │
│  │ lcm_expand   │  │ +46% recall  │  │ bg workers               │  │
│  └──────┬───────┘  └──────┬───────┘  └───────────┬───────────────┘  │
│         │                 │                       │                  │
│         └─────────────────┼───────────────────────┘                  │
│                           │                                          │
│  ┌────────────────────────▼─────────────────────────────────────┐    │
│  │              OpenClaw Gateway (9 agents)                      │    │
│  │  + obsidian-skills (kepano) — proper vault syntax             │    │
│  │  + awesome-skills catalog — discover new capabilities         │    │
│  │  + lossless-claw plugin — context engine slot                 │    │
│  └──────────────────────┬───────────────────────────────────────┘    │
│                         │                                            │
└─────────────────────────┼────────────────────────────────────────────┘
                          │
            ┌─────────────┼─────────────────┐
            │             │                 │
     ┌──────▼──────┐  ┌──▼──────────┐  ┌──▼────────────┐
     │  SwarmClaw   │  │   ClawX     │  │  Chowder iOS  │
     │  v1.1.5      │  │   (fork)    │  │  (reference)  │
     │  Dashboard   │  │  Desktop    │  │  for ios-app   │
     │  Missions    │  │  + MediaPipe│  │  Live Activity │
     │  Orchestrate │  │  + SmolVLM  │  │  thinking      │
     │  Connectors  │  │  + Ambient  │  │  steps         │
     └─────────────┘  └────────────┘  └───────────────┘
```

---

## Revised Deployment Phases (Incorporating Research)

### Phase 1: Memory Foundation (Week 1)
1. **Install Lossless Claw** — `openclaw plugins install @martian-engineering/lossless-claw`
2. **Configure session.reset.idleMinutes** to 10080 (7 days)
3. **Install obsidian-skills** — improve Legacy's vault writing quality
4. **Upgrade SwarmClaw** v0.7.7 → v1.1.5
5. **Enable heartbeat** on VPS (currently empty `{}`)

### Phase 2: Context Intelligence (Week 2-3)
1. **Deploy OpenViking** as context database server on VPS
2. **Configure OpenViking OpenClaw plugin** — benchmark against current BM25+vector
3. **Evaluate ClawVault** — install if entity graph adds value beyond OpenViking
4. **Archive 12 unused skills** (per PRODUCTION-PLAN.md)

### Phase 3: Desktop + Vision (Week 3-4)
1. **Fork ClawX** → `onemind-desktop`
2. Add SmolVLM-256M camera panel (from Master Plan)
3. Add MediaPipe Hands gesture recognition
4. Add ambient overlay mode (transparent, always-on-top)
5. Add NATS subscriber for live events

### Phase 4: Voice + Emotion (Month 2+)
1. Evaluate Chatterbox TTS (from Master Plan) vs current OpenAI TTS
2. Study Hume AI EVI for emotional voice
3. Study Chowder iOS patterns for our ios-app Live Activity
4. Implement Gemini Live bridge (from Master Plan — NixClaw)

### Phase 5: Physical (Month 3+)
1. RosClaw → Unitree G1 (from Master Plan)
2. DimensionalOS spatial memory
3. peaq DID machine identity

---

## Decision Matrix

| Question | Answer | Reasoning |
|---|---|---|
| Replace SwarmClaw with Mission Control? | **No** | SwarmClaw already deployed + configured. Upgrade to v1.1.5 instead. |
| Build desktop from scratch or fork ClawX? | **Fork ClawX** | 4.5k stars, 116 releases, already handles OpenClaw protocol. Add MediaPipe on top. |
| Lossless Claw or OpenViking first? | **Both, different layers** | LCM = context engine (conversation memory). OpenViking = knowledge DB (documents/resources). Complementary. |
| Replace ios-app with Chowder? | **Study, don't replace** | Our ios-app has more modules (Watch, Widget, Share Extension). Adopt Chowder's protocol patterns. |
| Install ClawdAgent? | **No** | Interesting ambition but 9 stars, different architecture. Study patterns only. |
| ClawVault worth it alongside OpenViking? | **Evaluate after Phase 2** | If OpenViking covers retrieval needs, ClawVault's entity graph may be redundant. But Obsidian hooks are unique. |
| Hume AI now or later? | **Later (Phase 4)** | Emotional voice is a luxury, not a foundation. Get memory + vision working first. |

---

## Cost Impact (Monthly)

| Component | Current | After Phase 1-2 |
|---|---|---|
| OpenClaw LLM (GPT-5.2) | ~$150-200 | ~$80-120 (OpenViking = -91% tokens on retrieval) |
| Heartbeat tokens | ~$30-60 | ~$15-30 (L0/L1 tiered loading) |
| SmolVLM-256M | $0 | $0 (local on MacBook Pro M2 Max) |
| OpenViking server | $0 | ~$5 VPS resources |
| Lossless Claw | $0 | ~$10-15 summarization LLM calls |
| **Total delta** | — | **~$50-100/mo savings** from token reduction |

---

## Files to Create/Modify

| Action | File | What |
|---|---|---|
| Modify | `/home/onemind/.openclaw/openclaw.json` | Add lossless-claw plugin config, extend session.reset |
| Create | VPS: OpenViking config | `~/.openviking/ov.conf` with text-embedding-3-small |
| Modify | `openclaw/workspaces/legacy/skills/` | Add obsidian-skills from kepano |
| Create | `command-center/onemind-desktop/` | ClawX fork with MediaPipe extensions |
| Modify | `PRODUCTION-PLAN.md` | Update with research findings and revised phases |
| Modify | `STATUS.md` | Add new services after deployment |
