# claude.md

Context for Claude Code when working in this repository.

---

## What This Repo Is

**One Mind** is an Obsidian vault — the world's first AI-native human enhancement framework. It is **not a software project**. Most files are markdown. Operational code lives exclusively in `_system/`.

Core equation: `HUMAN (Zeus) + AI (LEGACY) = ONE MIND`

| Context | Path |
|---------|------|
| Local (Mac) | `/Users/onemind/onemind/` |
| GitHub | `github.com/OneMind-OS/One-Mind` |

---

## Framework Architecture

One Mind is **3 life pillars running on 1 operating layer**, not 4 equal boxes.

```
┌───────────────────────────────────────────────────────┐
│           SO — Sovereign Operations (00–24)            │
│           The operating layer that runs everything     │
├─────────────────┬──────────────────┬──────────────────┤
│  HP (25–49)     │  LE (50–74)      │  GE (75–99)      │
│  Holistic       │  Legacy          │  Generational    │
│  Performance    │  Evolution       │  Empire          │
│   PERSONAL      │      HOME        │    BUSINESS      │
└─────────────────┴──────────────────┴──────────────────┘
```

Route new content by subject:
- AI agents, tech, automation, security → **SO (00-24)**
- Personal growth, fitness, habits, money → **HP (25-49)**
- Relationships, community, home → **LE (50-74)**
- Career, ventures, brand, wealth → **GE (75-99)**

### The CODE× Method

**C**apture → **O**rganize → **D**irect → **E**xecute → **×** Multiply

× is the multiplication sign (not the letter X). The 5 steps form a loop — × feeds back into Capture.

---

## Key Files — Read These First

| File | Purpose |
|------|---------|
| `onemind.md` | THE canonical master document — do not casually edit |
| `codex.md` | CODE× method deep spec — publishable IP |
| `agents.md` | Universal agent instructions (all AI agents read this) |
| `00-24 SO (Sovereign Operations)/00 Framework (The Constitution)/reference/SCHEMA.md` | Domain file format — frontmatter template, priority tags, sync rules |

---

## Directory Structure

```
onemind/
├── 00-24 SO (Sovereign Operations)/   ← Operating layer
│   ├── 00 Framework (The Constitution)/    ← Core philosophy, architecture, doctrine
│   ├── 01 Intelligence (The Brain)/        ← Agents, models, skills, memory, tools
│   ├── 03 Infrastructure (The Body)/       ← Servers, network, storage, services
│   ├── 04 Security (The Shield)/           ← Cybersecurity, access control, hardening
│   ├── 05 Operations (The Controls)/       ← Runbooks, SOPs, comms, hooks, scripts
│   ├── 05 Platforms (The Connectors)/      ← External tools, APIs, integrations catalog
│   ├── 06 Agents (The Fleet)/              ← Agent docs, roster, workspaces
│   ├── 06 Architecture (The Blueprint)/    ← System design, architectural layers
│   ├── 08 OM Assets (Entity Registry)/     ← System registry, nodes, templates
│   ├── 07 Development (The Forge)/         ← DevOps, deploy, scripts, guides
│   └── 08 Research (The Lab)/              ← R&D, robotics, AI agents, emerging tech
├── 25-49 HP (Holistic Performance)/   ← Personal pillar (10 domains: 25-34)
├── 50-74 LE (Legacy Evolution)/       ← Home pillar (8 domains: 50-57)
├── 75-99 GE (Generational Empire)/    ← Business pillar (9 domains: 75-83)
├── _system/                           ← Operational code ONLY (scripts, configs)
├── agents.md                          ← Universal agent instructions
├── claude.md                          ← This file
├── gemini.md                          ← Gemini CLI context
├── onemind.md                         ← THE canonical master document
├── codex.md                           ← CODE× method deep spec
└── README.md                          ← GitHub landing page
```

---

## Domain File Schema

Every domain folder has a `README.md` as its index. Standard frontmatter:

```yaml
---
domain: XX
name: Domain Name (Short Title)
agent: AgentName
status: active
todoist_project_id: XXXXXXXXXXXXXXX
review: YYYY-QQ
updated: YYYY-MM-DD
---
```

Standard sections: `## 🎯 Purpose`, `## 🏔 Goals`, `## ⚡ Active Tasks`, `## 📦 Completed Archive`, `## 📎 Context & Notes`, `## 🔗 Related`

Priority tags: `#p1` (Critical) → `#p4` (Low). Active tasks use `tid:TODOIST_TASK_ID` as the sync bridge to Todoist.

**Todoist is the execution layer. The vault is the source of truth for goals, context, and history.**

---

## Working Conventions

1. **New content goes to Inbox first** — `00-24 SO/24 Inbox (Queue)/`, then route to proper domain
2. **Preserve frontmatter** — Every file has YAML with `id`, `tags`, `owner`, `status` — do not strip it
3. **Code goes in `_system/` only** — Knowledge domain folders contain markdown only
4. **No secrets in markdown** — Document where credentials live (1Password), never the values
5. **1Password vault** — All secrets in vault ID `aaeyris5z44pdqspk277u2oppi` (`00-24 SO`). Reference format: `op://aaeyris5z44pdqspk277u2oppi/Item Name/field-name`
6. **Never rename domain group folders** — SO/HP/LE/GE names are locked; subdomains require Zeus's approval

---

## AI Agent Architecture

One Mind runs on **Hermes Agent** on the Mac as the primary orchestration engine (Legacy AI), with coding agents as specialist hands.

```
┌─────────────────────────────────────────────────────────────┐
│  Hermes (Brain)                                             │
│  → Legacy AI — command agent, full system access            │
│  → Persistent memory + skills (~/.hermes/)                  │
│  → Cron, orchestration, MCP servers                         │
├─────────────────────────────────────────────────────────────┤
│  Coding Agents (Hands)                                      │
│  → Claude Code ← YOU                                        │
│  → Copilot, Gemini CLI, Codex, OpenCode                     │
└─────────────────────────────────────────────────────────────┘
```

---

## Terminology Rules

| ❌ Do NOT use | ✅ Use instead |
|---------------|----------------|
| "quadrant" or "zone" | "pillar" |
| "four quadrants" | "4 pillars" or "3 life pillars + 1 operating layer" |
| "UI" for SO | "SO" or "Sovereign Operations" |
| "Codex Method" or "CODEX Method" | "CODE× Method" or "CODE×" |
| "eXamine" | "× Multiply" |
| "OneMind Codex" | "One Mind" |
| "empire" for GE tagline | "Empire" (internal), "Business" (public) |
| "Your Family" for LE | "Home" (public) |
