# Tech Stack — OneMind OS

> **Complete registry of software, tools, APIs, protocols, and services powering OneMind OS**

## Quick Stats

- **Total Tools**: 80+
- **Monthly Cost**: ~$200
- **Active Integrations**: 25+
- **Programming Languages**: Python, TypeScript, Dart
- **Primary Frameworks**: OpenClaw (agents), Next.js (web), Flutter (mobile)

---

## Folder Structure

```
Tech-Stack/
├── README.md                  ← You are here
├── SOFTWARE-CATALOG.md        ← Alphabetical list of all software
├── TOOLS-BY-CATEGORY.md       ← Tools organized by function + protocols
├── INTEGRATIONS-REGISTRY.md   ← APIs and external services
├── COST-TRACKING.md           ← Monthly subscription costs
├── STACK-DECISIONS.md         ← Architecture decisions (locked 2026-03)
├── DEPLOYMENT-SEQUENCE.md     ← Deployment priority order
├── ECOSYSTEM-42.md            ← All 42 ecosystem additions spec
├── INVENTORY.md               ← Ecosystem inventory
├── API-DOCUMENTATION.md       ← API reference docs
├── API-QUICK-START.md         ← API getting started
├── CODEX-INTEGRATION.md       ← Codex ↔ tech stack integration
├── DATABASE_SCHEMA.md         ← Database schema reference
├── HYBRID-SEARCH-SETUP.md     ← Vector + keyword search setup
├── OPENMIND-DB-REFERENCE.md   ← OpenMind database reference
└── platforms/                 ← Individual platform docs (61 files)
    ├── ai-agents/             ← Agno, Claude-Code, Letta, LiveKit, Whisper, etc.
    ├── automation/            ← N8N, Kestra, Node-RED
    ├── business/              ← Circle, Webflow, Zoho (16 sub-docs)
    ├── cloud/                 ← AWS, Azure, Cloudflare
    ├── communication/         ← Discord, Telegram, NATS, EMQX, Twilio, etc.
    ├── databases/             ← PostgreSQL, Redis, LanceDB, TimescaleDB
    ├── development/           ← Docker, FastAPI, GitHub, VS Code
    ├── payments/              ← Stripe, Whop, Fanbasis, Circle, Zoho
    ├── productivity/          ← Obsidian, Todoist, Raindrop, Habitica
    ├── robotics/              ← FarmBot, OM1, Unitree
    ├── security/              ← 1Password
    ├── smart-home/            ← Home Assistant, Frigate
    └── transport/             ← Uber
```

---

## Core Stack Overview

### AI & Agents

- **OpenClaw** — Agent gateway & orchestration (9 agents on VPS)
- **Claude Opus 4** — Primary AI model (via AWS Bedrock)
- **OpenAI GPT-4** — Alternative model
- **Letta** — Memory-augmented agents

### Backend

- **Python 3.12+** — Primary language
- **FastAPI** — API server
- **PostgreSQL + pgvector** — Database + vector embeddings
- **NATS JetStream** — Message broker

### Frontend

- **Next.js 16** — Web UI (onemind-web)
- **Flutter** — Mobile/desktop UI
- **TailwindCSS** — Styling
- **shadcn/ui** — Component library

### Infrastructure

- **Docker** — Containerization
- **DigitalOcean** — VPS hosting (primary)
- **AWS Bedrock** — LLM hosting
- **Cloudflare Pages** — Static site hosting (codex-live)
- **GitHub Actions** — CI/CD

### Communication

- **LiveKit** - Voice/video
- **MQTT** - IoT messaging
- **WebSockets** - Real-time updates

## Quick Reference

### Development

- **Editor**: VS Code / Cursor
- **Version Control**: Git + GitHub
- **Terminal**: iTerm2 + zsh
- **Package Managers**: pip, pnpm, brew

### Productivity

- **Tasks**: Todoist
- **Issues**: Linear
- **Notes**: Obsidian
- **Calendar**: Google Calendar

### Security

- **Secrets**: 1Password
- **VPN**: Tailscale
- **Auth**: OAuth 2.0

## Adding New Tools

When adding a new tool to the stack:

1. Add to **SOFTWARE-CATALOG.md** (alphabetical)
2. Add to **INTEGRATIONS-REGISTRY.md** (if API/service)
3. Add to **TOOLS-BY-CATEGORY.md** (appropriate category)
4. Update **COST-TRACKING.md** (if paid)
5. Update this README's Quick Stats

## Related Documentation

- **API Reference**: `08.4 API Reference/`
- **Integrations**: `08.5 Integrations/`
- **Cloud Infrastructure**: `08.1 Cloud Infrastructure/`
- **Platforms**: `platforms/`

---

Last Updated: January 25, 2026
