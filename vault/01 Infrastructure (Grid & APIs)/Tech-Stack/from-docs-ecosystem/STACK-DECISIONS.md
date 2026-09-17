# OneMind OS — Stack Decisions

> Architecture decisions for all 42 ecosystem items. Locked in 2026-03-03.
> Review quarterly or when hitting a specific wall.

---

## Core Principles

1. **One workspace per Droplet** — Legacy = one identity, one SOUL.md, one workspace
2. **Documentation before implementation** — every decision documented before code ships
3. **Memory stays basic** — OpenClaw built-in (SOUL.md, USER.md, memory/*.md, context window). No supermemory, Mem0, or pgvector recall until we hit a real "Legacy forgot" moment
4. **Cloud services first** — LiveKit Cloud, not self-hosted. Self-host only for privacy or cost at scale
5. **Add frameworks one at a time** — OpenClaw native → hit a wall → add ONE framework → hit another wall → add another. Never preemptive multi-framework.
6. **Subtrees for reference only** — read-only docs/catalogs. Active dev repos get forked to Zeus-Delacruz/
7. **n8n ships behind profile:full** — visual workflow builder for complex multi-API orchestration alongside Lobster

---

## Architecture: What Runs Where

```
┌──────────────────────────────────────────────────────────────────┐
│  DigitalOcean Droplet  (Ubuntu 24.04 / $28/mo)                  │
│                                                                  │
│  HOST (bare metal):                                              │
│    OpenClaw Gateway    :18789 (systemd, Tailscale-only)          │
│    Node.js 22 LTS                                                │
│                                                                  │
│  DOCKER COMPOSE:                                                 │
│    nats               :4222 / :7422 / :1883 / :8222              │
│    omnidb             :5432        [profile: data]               │
│    nats-ingester                   [profile: data]               │
│    legacy-bridge                   [always]                      │
│                                                                  │
│  DOCKER COMPOSE (profile: full):                                 │
│    n8n                :5678        [profile: full]               │
│    livekit (fallback) :7880        [profile: full]               │
│                                                                  │
│  NOT ON DROPLET (cloud):                                         │
│    LiveKit Cloud      → livekit.cloud (free tier: 50 hrs/mo)     │
│                                                                  │
│  SECURITY:                                                       │
│    Tailscale mesh VPN only. UFW deny all + Tailscale-only ports. │
│    Fail2ban. Workspace backup every 15m to private GitHub.       │
└──────────────────────────────────────────────────────────────────┘
```

---

## The 42: Decision Status

### Legend

| Status | Meaning |
| --- | --- |
| ✅ SHIP | Build/install now — in scope for v0.4 launch |
| 📋 DOCUMENTED | Decision made, documented, code not yet written |
| ⏸ DEFERRED | Clear plan exists, not needed yet |
| 🔬 EVALUATE | Needs hands-on testing before committing |
| ❌ SKIP | Not doing — doesn't fit current architecture |

---

### TIER 0 — Immediate

| # | Item | Decision | Method | Status |
| --- | --- | --- | --- | --- |
| 1 | clawsec | Install on first deploy | `openclaw skills install clawsec` | ✅ SHIP |
| 2 | ClawRouter | Install after clawsec, configure model tiers | Plugin in openclaw.json | ✅ SHIP |
| 3 | Agent Protocol | Install after testing sessions_spawn (#6) | `openclaw skills install agent-protocol` | 📋 DOCUMENTED |
| 4 | MCP Plugin | Install — unlocks 13K+ MCP servers | Plugin install | ✅ SHIP |
| 5 | Workspace Backup | Done — cron in install-openclaw.sh | Built into .do/install-openclaw.sh | ✅ SHIP |
| 6 | sessions_send/spawn | Built-in — learn and test, no install | Document patterns in workspace | 📋 DOCUMENTED |

**Tier 0 action**: Install #1, #2, #4, #5 on first deploy. Test #6 manually. Install #3 after #6.

---

### TIER 1 — Same Day

| # | Item | Decision | Method | Status |
| --- | --- | --- | --- | --- |
| 7 | supermemory | **SKIP for now** — memory stays basic per principles | — | ⏸ DEFERRED |
| 8 | openclaw-foundry | Install after 2 weeks of manual workflow patterns | Plugin | ⏸ DEFERRED |
| 9 | Composio | **SKIP for now** — use individual skills/MCP instead | — | ⏸ DEFERRED |
| 10 | Lobster Workflows | **SHIP** — write 3 core .lob files | Direct code in workspace/lobster/ | ✅ SHIP |
| 11 | openclaw-ansible | Reference only — our .do/ scripts cover this | Study patterns, don't run | 📋 DOCUMENTED |
| 12 | hardened-ansible | Evaluate after base deployment | Clone + evaluate | 🔬 EVALUATE |
| 13 | awesome-openclaw | Done — subtree in repo | Subtree (read-only reference) | ✅ SHIP |
| 14 | awesome-openclaw-skills | Done — subtree in repo | Subtree (read-only reference) | ✅ SHIP |
| 15 | Starter Kit Patterns | Extract: spend caps, failover, security checklist | Read + apply to openclaw.json | 📋 DOCUMENTED |
| 16 | ElevenLabs Voices | Install when voice channels are configured | `openclaw skills install elevenlabs-voices` | ⏸ DEFERRED |

**Tier 1 action**: Write Lobster workflows (#10). Study starter kits (#15) for spend caps. Everything else deferred.

---

### TIER 2 — Docker Compose Additions

| # | Item | Decision | Method | Status |
| --- | --- | --- | --- | --- |
| 17 | n8n | **SHIP** — deploy behind profile:full for complex multi-API workflows | Docker (profile: full) | ✅ SHIP |
| 18 | LiveKit | **USE CLOUD** — remove self-hosted container, use livekit.cloud free tier | Cloud service | ✅ SHIP (cloud) |
| 19 | openclaw-studio | **EVALUATE** — useful for debugging but not critical | Docker container | 🔬 EVALUATE |
| 20 | Mem0 | **SKIP** — memory stays basic | — | ⏸ DEFERRED |
| 21 | acpx | Install when multi-agent dispatch is needed | `npm install -g acpx` | ⏸ DEFERRED |
| 22 | openclaw-mcp | Deploy when Claude Desktop → Legacy bridge is useful | MCP server config | ⏸ DEFERRED |

**Tier 2 action**: Deploy n8n (#17) with `--profile full`. Configure LiveKit Cloud (#18) with cloud credentials. Skip everything else.

---

### TIER 3 — Hardware & Integrations

| # | Item | Decision | Method | Status |
| --- | --- | --- | --- | --- |
| 23 | VisionClaw | Fork to Zeus-Delacruz/, build when LiveKit Cloud works | Separate repo (fork) | ⏸ DEFERRED |
| 24 | WireClaw | Order ESP32 dev boards, fork repo | Separate repo (fork) | ⏸ DEFERRED |
| 25 | HeyClawy | After WireClaw — same platform | Separate repo (fork) | ⏸ DEFERRED |
| 26 | MimiClaw | Evaluate alongside WireClaw | Separate repo (fork) | 🔬 EVALUATE |
| 27 | OpenClawHomeAssistant | When Home Assistant is deployed at a location | HA add-on | ⏸ DEFERRED |
| 28 | openclaw_conversation | Lighter option — evaluate vs #27 | HACS component | 🔬 EVALUATE |
| 29 | HomeClaw | When Apple HomeKit devices are present | Mac Catalyst app | ⏸ DEFERRED |
| 30 | Jupiter Voice | When offline voice is a requirement | Local app | ⏸ DEFERRED |
| 31 | openclaw-voice | When browser-based voice is needed | Web app | ⏸ DEFERRED |
| 32 | openclaw-voice-avatar | After LiveKit Cloud + ElevenLabs working | Application | ⏸ DEFERRED |

**Tier 3 action**: None now. Fork repos when ready. Do NOT subtree these — they'll be actively developed.

---

### TIER 4 — Strategic / GE

| # | Item | Decision | Method | Status |
| --- | --- | --- | --- | --- |
| 33 | Virtual Protocol ACP | GE pillar — months out | Plugin + blockchain | ⏸ DEFERRED |
| 34 | Network-AI | **NOT NOW** — add only when running 3+ frameworks | Framework bridge | ⏸ DEFERRED |
| 35 | Langclaw | **NOT NOW** — add when Lobster can't handle a specific workflow | Python framework | ⏸ DEFERRED |
| 36 | ExoClaw | When running untrusted community skills | WASM runtime | ⏸ DEFERRED |
| 37 | BankrBot Skills | GE pillar — needs L1 HITL approval before ANY trade | Skill collection | ⏸ DEFERRED |
| 38 | ClawWork Benchmark | Run against Legacy to measure economic capability | Benchmark | 🔬 EVALUATE |
| 39 | CrewClaw | When managing 3+ agent teams | Web platform | ⏸ DEFERRED |
| 40 | OM1 (5 repos) | Strategic — when building OpenClaw↔OM1 bridge | Separate repos | ⏸ DEFERRED |
| 41 | Clawdstrike | When running multiple agent instances | Security tool | ⏸ DEFERRED |
| 42 | Community Platforms | Join Discord + Reddit now, rest later | Ongoing participation | 📋 DOCUMENTED |

**Tier 4 action**: Join community platforms (#42). Everything else is months out.

---

## Framework Progression Plan

This is the roadmap for when to add multi-framework capabilities:

```
NOW (v0.4)
│
│  OpenClaw native
│  ├── Skills (install from ClawHub)
│  ├── Lobster (.lob workflows)
│  ├── sessions_spawn (built-in multi-agent)
│  └── MCP Plugin (13K+ external tools)
│
├── WALL: "Lobster can't do complex branching with state"
│   → Add Langclaw (#35) for THAT workflow only
│
├── WALL: "Need LangChain RAG or DSPy extraction"
│   → Add that ONE framework directly
│
├── WALL: "Running 3+ frameworks, dispatch is messy"
│   → NOW add Network-AI (#34) as the router
│
└── WALL: "Untrusted skills need sandboxing"
    → Add ExoClaw (#36) for WASM isolation
```

**Rule**: Never add a framework preemptively. Each addition must solve a specific, documented problem.

---

## Multi-Framework Decision: Why Not Now

**Network-AI (#34)** connects 12 AI frameworks. Impressive on paper. Here's why it's wrong for v0.4:

| Concern | Reality |
| --- | --- |
| Dependencies | 12 frameworks = 12 Python/Node dependency trees |
| Debugging | Shared blackboard state across frameworks = hard to trace bugs |
| Maintenance | Each framework updates independently, breaking changes compound |
| Actual need | We haven't written our FIRST Lobster workflow yet |
| Cost | Every framework has its own model calls, memory usage, API keys |

**Langclaw (#35)** is LangGraph state machines. Better than Network-AI for a first addition because:
- Single framework, not 12
- Compiles to real LangGraph `CompiledStateGraph`
- Human-in-the-loop checkpoints (matches our L1 HITL pattern)
- BUT: requires Python runtime alongside Node.js
- BUT: we haven't outgrown Lobster yet

**Decision**: Use OpenClaw native for 2-3 months. Document every time Lobster/skills can't do something. When 3+ incidents point to the same gap → add the framework that fills it.

---

## LiveKit: Cloud vs Self-Hosted

| | Cloud (chosen) | Self-Hosted |
| --- | --- | --- |
| Setup | API key + SDK | Docker + TURN/STUN + SSL + public IP + UDP holes |
| Free tier | 50 participant-hours/mo | $0 but 4+ hours setup |
| For 1 user | **More than enough** | Over-engineered |
| WebRTC NAT | Built-in global TURN | You manage TURN servers |
| Latency | Edge servers worldwide | Your Droplet only |
| When to switch | >100 hrs/mo OR hard privacy requirement | — |

**Action to take in docker-compose.yml**:
- Keep livekit service behind `profile: full` but mark it as "self-hosted fallback"
- Primary integration uses `LIVEKIT_CLOUD_URL`, `LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET` from .env
- VisionClaw (#23) and voice-avatar (#32) point to cloud endpoint

---

## n8n: Ships with Profile Full

**Included in v0.4 deployment** behind `--profile full`.

n8n handles workflows that benefit from:
1. Touching 3+ external APIs in sequence with visual debugging
2. Retry logic with exponential backoff
3. Error branching (if step 2 fails → try alternate path)
4. Drag-and-drop node inspector for debugging complex flows

**Division of labor**: Lobster `.lob` files handle OpenClaw-native workflows (morning brief, health checks, git sync). n8n handles cross-platform orchestration that involves external APIs (Todoist → Zoho → Telegram → Discord chains).

Deploy: `docker compose --profile full up -d` starts n8n on `:5678`.

---

## Repo Structure (Final)

```
onemindos-v0.4/                          ← OneMind OS (the deployment)
├── .do/                                 ← Deployment scripts
│   ├── setup-droplet.sh                 ← Phase 1: base infra
│   ├── install-openclaw.sh              ← Phase 2: gateway + systemd
│   ├── setup-firewall.sh                ← Phase 3: Tailscale ports
│   ├── health-check.sh                  ← Phase 5: verify
│   └── DO-DEPLOYMENT.md                 ← Master deployment guide
├── docker-compose.yml                   ← All Docker services
├── .env.example                         ← Environment template
├── openclaw/                            ← OpenClaw-specific
│   ├── workspace/                       ← Legacy's identity + skills
│   │   ├── SOUL.md                      ← Personality
│   │   ├── AGENTS.md                    ← Instructions
│   │   ├── HEARTBEAT.md                 ← Proactive checks
│   │   ├── USER.md                      ← Zeus profile
│   │   ├── IDENTITY.md                  ← Name/vibe
│   │   ├── TOOLS.md                     ← Tool guidance
│   │   ├── BOOT.md                      ← Startup ritual
│   │   ├── bridge.py                    ← NATS→OpenClaw forwarder
│   │   ├── openclaw.json                ← Gateway config
│   │   ├── lobster/                     ← Workflow .lob files
│   │   │   ├── morning-brief.lob        ← Daily briefing (7AM)
│   │   │   ├── health-check.lob         ← System health (5 min)
│   │   │   └── git-sync.lob            ← Codex sync (15 min)
│   │   └── skills/                      ← Installed skills
│   ├── nats/                            ← NATS JetStream config
│   │   ├── hub.conf                     ← Hub config
│   │   ├── nats_ingester.py             ← NATS→OmniDB pipeline
│   │   └── timescaledb/init.sql         ← Database schema
│   ├── docs/                            ← Documentation
│   │   └── ecosystem/                   ← Ecosystem 42
│   │       ├── ECOSYSTEM-42.md          ← Full catalog
│   │       ├── STACK-DECISIONS.md       ← This file
│   │       ├── DEPLOYMENT-SEQUENCE.md   ← Install order
│   │       ├── INVENTORY.md             ← Raw research
│   │       ├── README.md                ← Overview
│   │       ├── awesome-openclaw/        ← Subtree (reference)
│   │       └── awesome-openclaw-skills/ ← Subtree (reference)
│   ├── hardware/                        ← Hardware integration docs
│   ├── nodes/                           ← Edge node configs (HA, etc.)
│   ├── reference/                       ← Reference implementations
│   └── tools/                           ← Custom tools
├── frontend/                            ← Flutter web app
└── legacy_docs/                         ← Historical documentation
```

**Why NOT flatten openclaw/ to root**: This repo is OneMind OS, not just OpenClaw. When Network-AI (#34), Flutter frontend, or other runtimes are added, they sit alongside openclaw/ at the same level. The repo is the deployment orchestrator; openclaw/ is one runtime within it.

---

## Subtree Policy

**Only subtree read-only reference material.** Currently:

| Subtree | Path | Purpose |
| --- | --- | --- |
| awesome-openclaw | `openclaw/docs/ecosystem/awesome-openclaw/` | Ecosystem discovery |
| awesome-openclaw-skills | `openclaw/docs/ecosystem/awesome-openclaw-skills/` | Skill catalog |

**Do NOT subtree**:
- Hardware repos (VisionClaw, WireClaw, etc.) → Fork to Zeus-Delacruz/, clone separately
- Active development repos → separate repos
- Frameworks (Langclaw, Network-AI) → add as dependencies when needed

**Update subtrees** quarterly:
```bash
git subtree pull --prefix=openclaw/docs/ecosystem/awesome-openclaw \
  https://github.com/rohitg00/awesome-openclaw.git main --squash
```

---

## What Ships in v0.4

| Category | Items |
| --- | --- |
| **Install on deploy** | clawsec (#1), ClawRouter (#2), MCP Plugin (#4) |
| **Already done** | Workspace backup (#5), awesome subtrees (#13, #14) |
| **Write code for** | Lobster workflows (#10): morning-brief, health-check, git-sync |
| **Configure** | LiveKit Cloud (#18): add cloud credentials to .env |
| **Document/test** | sessions_spawn (#6), starter kit patterns (#15) |
| **Join** | Community platforms (#42): Discord, Reddit |
| **Deploy with profile:full** | n8n (#17): multi-API workflow orchestration |
| **Everything else** | Deferred with clear trigger conditions |
