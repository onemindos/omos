# CODEX.md Template

> *CODEX.md is the vault architecture document every agent reads.*
> *It tells agents how the vault is structured, how to navigate it, and how to write to it.*
> *Deploy to all workspaces: `~/.openclaw/workspace[-agent-id]/CODEX.md`*

---

# CODEX.md — One Mind OS Vault

## What This Is

The One Mind OS vault is a git-based knowledge system — the central brain of Zeus's life and empire. Every domain, goal, project, and capture lives here in plain markdown.

**Vault path (VPS):** `/root/onemind-codex/`
**Vault access:** Via `codex-vault` MCP server — use `read_file`, `write_file`, `list_directory` tools
**GitHub:** `github.com/OneMind-OS/onemind-codex` — source of truth

---

## Architecture — 4 Pillars

```
┌──────────────────────────────────────────────┐
│   SO — Sovereign Operations (00–24)           │
│   The operating layer that runs everything    │
├──────────────┬─────────────┬─────────────────┤
│ HP (25–49)   │ LE (50–74)  │ GE (75–99)      │
│ Personal     │ Home        │ Business        │
└──────────────┴─────────────┴─────────────────┘
```

**Route by subject:**
- AI agents, tech, automation, security → SO (00–24)
- Personal growth, fitness, money, habits → HP (25–49)
- Relationships, family, home → LE (50–74)
- Brand, ventures, investments → GE (75–99)

---

## Domain Map

### SO — Sovereign Operations (00–24)

| Domain | Name |
|--------|------|
| 00 | Framework (Doctrine) |
| 01 | Intelligence (AI, Memory & Tools) |
| 02 | Infrastructure (VPS, Servers & Networking) |
| 03 | Command (Operations & Navigation) |
| 04 | Security (Cyber & Digital Shield) |
| 05 | Development (Code & Testing) |
| 06 | Research (R&D) |
| 07 | Automations (Workflows & Triggers) |
| 08 | OM Assets (Entity Registry) |
| 99 | Inbox (Queue) |

### HP — Holistic Performance (25–49)

| Domain | Name |
|--------|------|
| 25 | Identity (Self, Purpose & Spirit) |
| 26 | Spirituality (Metaphysics) |
| 27 | Mind (Education & Learning) |
| 28 | Body (Health & Wellness) |
| 29 | Relationships (Social) |
| 30 | Finance (Money) |
| 31 | Mastery (Skills & Training) |
| 32 | Rhythm (Routines & Productivity) |
| 33 | Sovereignty (Legal) |
| 34 | Joy (Recreation & Travel) |

### LE — Legacy Evolution (50–74)

| Domain | Name |
|--------|------|
| 50 | Home (Sanctuary & Property) |
| 51 | Covenant (Relationship & Union) |
| 52 | Offspring (Parenting & Education) |
| 53 | Parents (Elder Care & Support) |
| 54 | Heritage (History & Traditions) |
| 55 | Health (Family Wellness) |
| 56 | Estate (Family Wealth & Legal) |
| 57 | Community (Extended Network) |

### GE — Generational Empire (75–99)

| Domain | Name |
|--------|------|
| 75 | Zeus Delacruz (Brand & Empire) |
| 76 | Real Estate (Landlord & Property) |
| 77 | Acquisitions (Bought Businesses) |
| 78 | Holdings (Legal & Entity Structure) |
| 79 | Investments (Wealth & Assets) |
| 80 | Innovation (New Ventures Pipeline) |
| 81 | Strategy (Portfolio Intelligence) |
| 82 | Partnerships & Capital |
| 83 | Philanthropy (Giving) |

---

## How to Write to the Vault

**New capture:** Write to `99 Inbox (Queue)/YYYY-MM-DD-[title].md`
**Domain content:** Write to the correct domain folder, follow existing file naming
**Goal:** Write to `[domain]/goals/goal-[slug]-[YYYYMMDD].md`
**Daily note:** Write to `HP/32 Rhythm/daily/YYYY-MM-DD.md`

**All files must have frontmatter.** Minimum:

```yaml
---
id: [slug-YYYYMMDD]
domain: XX
type: [note | goal | capture | resource | project]
status: active
owner: Zeus
updated: YYYY-MM-DD
---
```

---

## Naming Conventions

**Files:** `kebab-case-descriptive-name.md`
**Domains:** Never rename numbered domain folders — names are locked
**Do NOT use:** Spaces in filenames, underscores unless necessary, vague names like `notes.md`

---

## Terminology Rules

| Wrong | Correct |
|-------|---------|
| "quadrant" | "pillar" |
| "CODEX Method" | "CODE× Method" |
| "eXamine" | "× Multiply" |
| "OneMind Codex" | "One Mind" |
| "enterprise" for GE | "Empire" |

---

## After Writing

Auto-sync is active. Any file written to `/root/onemind-codex/` is auto-committed and pushed to GitHub within 3 seconds via `onemind-watch.service`. You do not need to run git commands.

---

*One Mind OS — Agent Vault Reference*
