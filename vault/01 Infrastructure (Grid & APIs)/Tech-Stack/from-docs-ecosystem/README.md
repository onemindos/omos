# OneMind OS — OpenClaw Ecosystem Integration

> Everything the ecosystem offers, mapped to our architecture, organized by deployment priority.

## Folder Structure

```
ecosystem/
├── README.md                        ← You are here
├── ECOSYSTEM-42.md                  ← All 42 additions with full specs
├── DEPLOYMENT-SEQUENCE.md           ← Step-by-step deployment order
├── INVENTORY.md                     ← Original ecosystem research document
├── awesome-openclaw/                ← git subtree: rohitg00/awesome-openclaw (195K★)
└── awesome-openclaw-skills/         ← git subtree: VoltAgent/awesome-openclaw-skills (2,868 curated)
```

**Related folders:**
- `openclaw/docs/usecases/` — git subtree: hesamsheikh/awesome-openclaw-usecases (34 use cases)

## Priority Tiers

| Tier | Items | Effort | Description |
|------|-------|--------|-------------|
| **0** | 1-6 | Minutes | Install/enable immediately — security, cost savings, built-in features |
| **1** | 7-16 | Hours | Same-day additions — memory, workflows, ansible, skill catalogs |
| **2** | 17-22 | Days | Docker Compose additions — n8n, LiveKit, studio, MCP bridges |
| **3** | 23-32 | Weeks | Hardware + subtrees — smart glasses, ESP32, voice, Home Assistant |
| **4** | 33-42 | Months | Strategic — agent commerce, multi-framework, benchmarks, OM1 bridge |

## Current Status

- [x] awesome-openclaw-usecases subtree (34 use cases)
- [x] awesome-openclaw subtree (195K★ definitive collection)
- [x] awesome-openclaw-skills subtree (2,868 curated skills)
- [ ] Tier 0 deployed
- [ ] Tier 1 deployed
- [ ] Tier 2 deployed
- [ ] Tier 3 deployed
- [ ] Tier 4 deployed

## Quick Reference: What Connects Where

```
OpenClaw (bare metal)
├── Skills (install via openclaw skills install)
│   ├── #1  clawsec            — security suite
│   ├── #3  agent-protocol     — pub/sub agent chains
│   ├── #4  mcp-plugin         — 13,000+ MCP servers
│   └── #16 elevenlabs-voices  — 18 personas, 32 languages
│
├── Plugins (extend runtime)
│   ├── #2  ClawRouter         — 78% LLM cost savings
│   ├── #7  supermemory        — auto-recall on every turn
│   ├── #8  openclaw-foundry   — self-upgrading skills
│   └── #9  Composio           — 500+ app integrations
│
├── Workflows (Lobster .lob files)
│   └── #10 morning-brief, git-sync, health-check
│
├── Built-in (already in OpenClaw core)
│   └── #6  sessions_send / sessions_spawn
│
└── Infrastructure
    ├── #5  workspace backup   — cron → private GitHub
    ├── #11 openclaw-ansible   — automated hardened install
    └── #12 hardened-ansible   — rootless Podman + LiteLLM + egress

Docker Compose (containers)
├── #17 n8n                    — deterministic workflow engine
├── #18 LiveKit                — real-time AV / WebRTC
├── #19 openclaw-studio        — admin dashboard
├── #20 Mem0                   — persistent memory layer
├── #21 acpx                   — bridge Claude/Gemini/Copilot into OpenClaw
└── #22 openclaw-mcp           — Claude.ai → Legacy bridge

Hardware (edge nodes)
├── #23 VisionClaw             — Meta Ray-Ban smart glasses
├── #24 WireClaw               — $5 ESP32 NATS agents
├── #25 HeyClawy               — ESP32 voice companion per room
├── #26 MimiClaw               — ESP32-S3 pure C agent
├── #27 HA Add-on              — OpenClaw inside Home Assistant
├── #28 HA Conversation        — OpenClaw as HA Assist backend
├── #29 HomeClaw               — Apple HomeKit via MCP
├── #30 Jupiter Voice           — fully local STT/TTS fallback
├── #31 openclaw-voice          — browser-based voice chat
└── #32 voice-avatar            — LiveKit lip-synced Legacy face

Strategic (GE pillar + future)
├── #33 Virtual Protocol ACP   — agent-to-agent commerce
├── #34 Network-AI             — 12-framework orchestrator
├── #35 Langclaw               — LangGraph workflows
├── #36 ExoClaw                — WASM-sandboxed agent runtime
├── #37 BankrBot skills        — crypto/DeFi/Polymarket
├── #38 ClawWork benchmark     — economic capability testing
├── #39 CrewClaw               — visual multi-agent dashboard
├── #40 OM1 sub-repos          — robotics bridge (5 repos)
├── #41 Clawdstrike            — fleet security monitoring
└── #42 Community platforms    — Discord/Reddit/X/meetups
```

## Subtree Sync Commands

```bash
# Pull upstream updates for any subtree:
git subtree pull --prefix=openclaw/docs/usecases \
  https://github.com/hesamsheikh/awesome-openclaw-usecases.git main --squash

git subtree pull --prefix=openclaw/docs/ecosystem/awesome-openclaw \
  https://github.com/rohitg00/awesome-openclaw.git main --squash

git subtree pull --prefix=openclaw/docs/ecosystem/awesome-openclaw-skills \
  https://github.com/VoltAgent/awesome-openclaw-skills.git main --squash
```
