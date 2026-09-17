# Tech Stack - OneMind OS

> **Complete registry of software, tools, APIs, and services powering OneMind OS**

## Quick Stats

- **Total Tools**: 80+
- **Monthly Cost**: ~$200
- **Active Integrations**: 25+
- **Programming Languages**: Python, TypeScript, Dart
- **Primary Frameworks**: Agno, Next.js, Flutter

## Documents in This Folder

| File                         | Purpose                                    |
| ---------------------------- | ------------------------------------------ |
| **SOFTWARE-CATALOG.md**      | Complete alphabetical list of all software |
| **INTEGRATIONS-REGISTRY.md** | APIs and external services                 |
| **TOOLS-BY-CATEGORY.md**     | Tools organized by function                |
| **COST-TRACKING.md**         | Monthly subscription costs                 |

## Core Stack Overview

### AI & Agents

- **Agno** - Agent framework (Python)
- **Claude Opus 4.5** - Primary AI model (via AWS Bedrock)
- **OpenAI GPT-4** - Alternative model
- **Letta** - Memory-augmented agents

### Backend

- **Python 3.12+** - Primary language
- **FastAPI** - API server
- **PostgreSQL + pgvector** - Database
- **NATS JetStream** - Message broker

### Frontend

- **Next.js 15** - Web UI (App Router)
- **Flutter** - Mobile/desktop UI
- **TailwindCSS** - Styling
- **shadcn/ui** - Component library

### Infrastructure

- **Docker** - Containerization
- **AWS** - Cloud hosting
- **Vercel** - Web deployment
- **GitHub Actions** - CI/CD

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
