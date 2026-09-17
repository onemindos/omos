---
id: root-00-20260113225954
title: "CLAUDE"
type: agent
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/agent, status/active]
ai_summary: "CLAUDE"
---
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

**OneMind-Codex** is a solar-punk AI-native life operating system built as an Obsidian vault. It combines personal knowledge management, AI agent coordination, and life automation into a unified system. The core equation: `HUMAN (Zeus) + AI (LEGACY) = ONE MIND`.

This is NOT a traditional software project. It is:
- An **Obsidian vault** with structured markdown files
- A **Letta AI agent plugin** for Obsidian (TypeScript)
- **Python tools** for external service integrations
- **Infrastructure configs** for self-hosted services

## Architecture

### The Four Domains (25 slots each, 100 total)

```
00-24 UI (Unified Intelligence) — The Brain (AI coordination + tech infrastructure)
25-49 HP (Holistic Performance) — Self Path (PI Score)
50-74 LE (Legacy Evolution)     — Legacy Path (LI Score)
75-99 GE (Generational Entrepreneurship) — Wealth Path (WI Score)
```

### 00-24 UI Structure (Current)

```
00-24 UI (Unified Intelligence)/
├── 00 Framework (Doctrine)           # Core philosophy, architecture docs
├── 01 Command (Dashboard & Navigation) # Human command center, metrics
├── 02 Agents (AI, Memory & Tools)    # AI agents, Letta, memory blocks
├── 03 Protocols (SOPs)               # Standard operating procedures
├── 04 Resources (Assets & References) # Templates, guides, assets
├── 05 Environments (Contexts & Locations) # Context switching
├── 06 Inbox (Queue)                  # Capture point, unprocessed items
├── 08 Infrastructure (Grid & APIs)   # Servers, integrations, cloud
├── 09 Automation (Workflows)         # N8N, automations
├── 10 Development (Code & Testing)   # Code, repos, devops
├── 11 Security (Cyber)               # Cybersecurity, credentials
├── 12 Research (R&D)                 # Emerging tech, research notes
├── 24 Archive (Historical)           # Historical records
└── letta-tools/                      # Python tools for Letta agents
```

### AI Agent Architecture (Hive Mind)

Four Letta-powered agents share memory and coordinate:
- **UI (Mother AI)**: Central orchestrator + tech coordinator (00-24)
- **HP Agent**: Personal coach for self-optimization (25-49)
- **LE Agent**: Family steward for legacy preservation (50-74)
- **GE Agent**: Business advisor for wealth creation (75-99)

Agents route requests based on domain and share core memory blocks (Zeus Profile, Current Context, System State).

## Development Commands

### Letta AI Plugin (`.obsidian/plugins/letta-ai-agent/`)

```bash
cd .obsidian/plugins/letta-ai-agent
npm install                    # Install dependencies
npm run dev                    # Development mode with watch
npm run build                  # Production build with type check
```

The plugin uses `@letta-ai/letta-client` for API communication with Letta services.

### Letta Tools (`00-24 UI (Unified Intelligence)/letta-tools/`)

Python tools that give LEGACY AI ability to interact with external services:

```bash
cd "00-24 UI (Unified Intelligence)/letta-tools"
pip install -r requirements.txt
cp .env.example .env           # Configure API keys
```

**Available tool integrations:**
- `tools/discord/` — Notifications, human approval requests
- `tools/todoist/` — Task management
- `tools/zoho/` — CRM, calendar, support tickets
- `tools/circle/` — Community platform
- `tools/home_assistant/` — Smart home control
- `tools/1password/` — Credential queries, secure access, management (NEW - Phase 8)
- `tools/core/` — Agent utilities

**Setup scripts:**
- `scripts/letta_agent_setup.py` — Main agent configuration
- `scripts/attach_discord_tools.py` — Discord integration
- `scripts/todoist_setup.py` — Todoist structure setup

### 1Password Integration (NEW - Phase 7/8)

1Password has been restructured into 4 OneMind domain vaults with tag-based organization:

**Vaults:**
- `OneMind - Unified Intelligence` (UI) — 63 items
- `OneMind - Holistic Performance` (HP) — 103 items
- `OneMind - Legacy Evolution` (LE) — 14 items
- `OneMind - Generational Entrepreneurship` (GE) — 48 items

**Tag System:** Every credential has 2 tags:
- `path:ui|hp|le|ge` — Domain classification
- `category:finance|it|infrastructure|etc.` — Functional category

**Python Tools (Phase 8):**
- `tools/1password/query.py` — Query credentials by path/category
- `tools/1password/secure_access.py` — Approval workflows for sensitive items
- `tools/1password/manage.py` — Create/update/rotate credentials
- `tools/1password/sync.py` — Sync with vault and other systems

**Documentation:**
- Complete docs: `/1PASSWORD-CREDENTIALS/1PASSWORD-VAULT-RESTRUCTURING.md`
- Quick reference: `/1PASSWORD-CREDENTIALS/1password-quick-reference.md`
- Status: 94.3% tagged, ready for AI integration

## Key Files

| Purpose | Location |
|---------|----------|
| Master blueprint | `ONEMIND-CODEX.md` |
| Data architecture | `DATA-ARCHITECTURE.md` |
| Agent configurations | `00-24 UI/02 Agents (AI, Memory & Tools)/` |
| System prompts | `00-24 UI/02 Agents/02.3 Agents/02.3.18 System Prompts/` |
| Infrastructure configs | `00-24 UI/08 Infrastructure (Grid & APIs)/` |
| N8N workflows | `00-24 UI/09 Automation (Workflows)/` |
| Letta tools | `00-24 UI/letta-tools/` |

## Root Level Structure

```
OneMind-Codex/
├── 00-24 UI (Unified Intelligence)/  # The Brain
├── 25-49 HP (Holistic Performance)/  # Self Path
├── 50-74 LE (Legacy Evolution)/      # Legacy Path
├── 75-99 GE (Generational Entrepreneurship)/ # Wealth Path
├── Organize-important/               # Files pending relocation (important)
├── Organize-non-important/           # Archives and staging areas
├── ONEMIND-CODEX.md                  # THE canonical master document
├── README.md                         # GitHub landing page
├── CLAUDE.md                         # This file (AI instructions)
├── DATA-ARCHITECTURE.md              # Data systems architecture
└── _CLAUDE-NEXT-STEPS.md             # Session continuity (temporary)
```

## Naming Conventions

### Folders
Format: `XX Name (Subtitle)/` where XX is the domain number (00-99)

### Files within folders
Format: `XX.Y.Z Name/` for nested organization

### Content routing by domain
- **HP (25-49)**: Personal growth, health, skills, personal finance
- **LE (50-74)**: Family, home, estate planning, generational wealth
- **GE (75-99)**: Business operations, revenue generation, ventures
- **UI (00-24)**: AI agents, tech infrastructure, automation, security

## Integration Stack

```
Obsidian (Vault) ↔ Todoist (Tasks) ↔ Mobile
        ↓                ↓
    N8N (Sync Layer)
        ↓
GitHub (Backup) → Letta (AI Memory) ↔ 4 Agents
```

**Services (when self-hosted):**
- N8N: 5678
- Letta: 8283
- Langfuse: 3000
- Uptime Kuma: 3001
- Prometheus: 9090
- Grafana: 3002

## AI Agent Autonomy Levels

| Level | Trust | Actions |
|-------|-------|---------|
| L3 | Full | Read files, route requests, search |
| L2 | High | Create tasks, update metrics, log |
| L1 | Medium | HITL required (financial decisions, external comms) |
| L0 | None | Human only (legal, health emergencies, deletions) |

## Working with This Repository

1. **This is primarily documentation** — Most files are markdown, not code
2. **Obsidian is the primary interface** — Open the root folder as an Obsidian vault
3. **Code lives in specific places:**
   - Plugin: `.obsidian/plugins/letta-ai-agent/`
   - Python tools: `00-24 UI/letta-tools/`
   - Infrastructure scripts: `00-24 UI/08 Infrastructure/`
4. **Security**: Never commit API keys or secrets. Use `.env` files locally, GitHub Secrets for CI/CD.
5. **Backups**: The vault syncs to GitHub. Changes should be committed regularly.

## Common Tasks

**Modify AI agent behavior**: Edit system prompts in `02.3.18 System Prompts/` and memory blocks in `02 Agents/`

**Add new Letta tools**: Create Python functions in `00-24 UI/letta-tools/tools/[platform]/` following existing patterns

**Update plugin functionality**: Edit `.obsidian/plugins/letta-ai-agent/main.ts`, then run `npm run build`

**Infrastructure changes**: Configs in `08 Infrastructure/`, deployment scripts in `08 Infrastructure/`

## Session Continuity

When the user mentions continuing previous work, resuming, or asks "what were we working on":

1. **Check MCP Memory Graph**: Query `mcp__memory__read_graph` or `mcp__memory__search_nodes` for project context
2. **Read Next Steps File**: Check `_CLAUDE-NEXT-STEPS.md` at vault root (if it exists) for detailed task lists
3. **Check Git History**: Run `git log --oneline -5` to see recent work

Key memory entities to look for:
- `OneMind-Codex-Restructure` - Vault restructuring project status

This ensures continuity across Claude sessions without manual context re-entry.

## GitHub Copilot CLI & Agent Orchestration

Claude + Copilot CLI work together for 24/7 parallel development across all OneMind projects.

### Copilot CLI Installation

```bash
# Installed via Homebrew (already installed)
brew install --cask copilot-cli

# Location: /opt/homebrew/bin/copilot
# Version: 0.0.381+
```

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    DEVELOPMENT TEAM                             │
│                                                                 │
│   CLAUDE CODE          COPILOT CLI           COPILOT AGENTS    │
│   (Orchestrator)       (Interactive)         (Background)      │
│   • Plans work         • Edit files          • Process issues  │
│   • Creates issues     • Run commands        • Create PRs      │
│   • Reviews PRs        • /delegate tasks     • 24/7 execution  │
│   • Merges code        • MCP integrations    • Self-assign     │
└─────────────────────────────────────────────────────────────────┘
```

### Primary Repositories

| Repo | Purpose | Status |
|------|---------|--------|
| `Zeus-Delacruz/onemind-os` | Main app (Next.js portal) | ✅ Active |
| `Zeus-Delacruz/onemind-infra` | Infrastructure (Terraform/K8s) | ✅ Active |
| `Zeus-Delacruz/OneMind-Codex` | Obsidian vault + docs | 🔶 Docs only |

### Copilot CLI Commands (Interactive Terminal)

**Start Copilot CLI:**
```bash
# Launch interactive session
copilot

# Or in a specific directory
cd ~/path/to/repo && copilot
```

**Work on GitHub Issues:**
```
# In copilot interactive session:
> Work on issue #3 in Zeus-Delacruz/onemind-os
> Show me open issues assigned to Copilot
> What PRs are open from copilot-swe-agent?
```

**Delegate to Copilot Coding Agent (runs on GitHub):**
```
# Hand off current work to background agent
> /delegate implement the voice transcription panel from issue #3
> /delegate finish the API tests and fix edge cases
```

**Edit Files & Run Commands:**
```
# Copilot CLI can directly modify code
> Add error handling to the /api/agents route
> Run the tests and fix any failures
> Create a new component for the settings page
```

**MCP Server Integration:**
```
# Connect to custom MCP servers for extended context
> /mcp connect github
> /mcp list
```

### gh CLI Commands (for monitoring)

**Check Copilot agent work:**
```bash
# See all Copilot PRs
gh pr list --repo Zeus-Delacruz/onemind-os --author="app/copilot-swe-agent" --state=open

# See issues assigned to Copilot
gh search issues --assignee=Copilot --state=open --json number,title,repository,url

# Review a specific PR
gh pr diff <NUMBER> --repo Zeus-Delacruz/onemind-os
```

**Merge Copilot PR:**
```bash
gh pr merge <NUMBER> --repo Zeus-Delacruz/onemind-os --squash --delete-branch
```

### Issue Template for Copilot

```markdown
## Description
[Clear 1-2 sentence description of what to build]

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## Technical Details
- Framework: [Next.js/Python/etc]
- Files to modify: `path/to/file.ts`
- Dependencies: [any packages needed]

## Examples
[Code snippets or references to similar implementations]
```

### Workflow: Continuous Delegation

1. **Start Session** — Run `copilot` in terminal
2. **Review Issues** — Ask Copilot to show open issues
3. **Work or Delegate** — Either work interactively OR `/delegate` to background agent
4. **Monitor Progress** — Use `gh pr list` to check agent PRs
5. **Quality Gate** — Claude Code reviews all PRs before merge

### Quality Checklist for Copilot PRs

Before merging any Copilot PR:
- [ ] Code follows existing patterns
- [ ] No hardcoded secrets
- [ ] Error handling present
- [ ] Types correct (TypeScript)
- [ ] No security vulnerabilities
- [ ] Tests pass
- [ ] No unnecessary changes

### Scaling Strategy

**Phase 1 (Current):** Copilot CLI + Claude Code
- Claude creates detailed issue specs
- Use `copilot` CLI interactively or `/delegate` to agents
- Claude reviews and merges PRs

**Phase 2:** Parallel Agents
- Multiple `/delegate` tasks running simultaneously
- Automated PR notifications via Discord/Telegram
- Claude monitors and reviews

**Phase 3:** Full Automation (N8N + LEGACY)
- LEGACY AI creates issues autonomously
- Auto-delegate to Copilot agents
- Claude Code reviews all PRs
- Auto-merge on passing quality gate

## Recent Restructuring (2026-01)

The vault underwent significant restructuring:
- **00-24 UI** consolidated from 18 → 13 main folders
- **Root files** relocated to `Organize-important/` for cleanup
- **Archives** consolidated in `Organize-non-important/`
- **letta-tools/** moved inside `00-24 UI/`

See `_CLAUDE-NEXT-STEPS.md` for current restructuring status and next tasks.
