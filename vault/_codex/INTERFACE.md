---
id: codex-interface
title: "Codex Interface — Universal Agent Entry Point"
type: system
status: active
created: 2026-03-18
modified: 2026-03-28
owner: zeus
tags: [type/system, domain/codex]
ai_summary: "Universal agent entry point for the OneMind Codex vault. Start here."
---

# Codex Interface

> **Read this file first.** If you are an AI agent, this tells you everything you need to work with this vault.

---

## What This Is

This is a **OneMind Codex** — a structured knowledge base in plain Markdown, organized using the CODEX method. It is the single source of truth for all knowledge, plans, decisions, relationships, and documentation.

**Format**: Markdown + YAML frontmatter + `[[wiki-links]]` + tags
**Storage**: Git repository (versioned, portable, forkable)
**Editor**: Obsidian (optional — works with any text editor)

> **Broader ecosystem?** This vault is one layer of a larger system.
> See `ECOSYSTEM.md` at the repo root for the full map (all repos, asset vaults, apps, access patterns).
> Machine-readable version: `ecosystem-manifest.json` at the repo root.

---

## How It's Organized — 3 Domains + 1 Operating Layer, 100 Sub-Domains

```
00-24  SO (Sovereign Operations) — Tech, AI, infrastructure, architecture
25-49  HP (Holistic Performance) — Health, skills, finance, identity, relationships
50-74  LE (Legacy Evolution)     — Family, home, children, estate, agriculture
75-99  GE (Generational Endeavor) — Business, ventures, wealth
```

Each top-level folder is a domain group. Inside each domain group, numbered subfolders are domains:

```
codex/
├── 00-24 SO/
│   ├── 00 Framework/       ← Core doctrine and principles
│   ├── 01 Command/         ← Dashboards and navigation
│   ├── 02 Agents/          ← AI agents, memory, tools
│   ├── 06 Protocols/       ← SOPs and procedures
│   ├── 06 Inbox/           ← Queue for unsorted items
│   ├── 08 Infrastructure/  ← Servers, APIs, services
│   ├── 10 Development/     ← Code, testing, CI/CD
│   ├── 11 Security/        ← Cybersecurity, access control
│   └── 12 Research/        ← R&D, analysis, evaluations
├── 25-49 HP/
│   ├── 25 Identity/        ← Self, purpose, values
│   ├── 27 Body/            ← Health, fitness, nutrition
│   ├── 28 Mastery/         ← Skills, career, expertise
│   ├── 29 Life Systems/    ← Daily operations, routines
│   ├── 30 Finance/         ← Money, investments, budgets
│   └── ...
├── 50-74 LE/
│   ├── 50 Home/            ← Property, sanctuary
│   ├── 51 Marriage/        ← Partnership
│   ├── 52 Children/        ← Parenting
│   └── ...
├── 75-99 GE/
│   ├── 75 Ventures/        ← Businesses, products
│   └── ...
└── _codex/                  ← THIS FOLDER — the vault's brain
    ├── system/              ← PINNED CONTEXT — always loaded into agent prompt
    ├── skills/              ← Agent instructions (search, capture, organize, etc.)
    ├── templates/           ← Note templates (daily, decision, project, etc.)
    └── profiles/            ← Context modes (planning, review, incident)
```

---

## How to Work With This Vault

### Filetree as Navigation

The folder hierarchy and file names serve as a **navigational table of contents**. Agents should always have access to the directory listing — it tells you what exists without opening any files. Descriptive folder names (numbered domains) and descriptive file names (lowercase-hyphenated) make the filetree scannable.

### Pinned Context (`_codex/system/`)

Files in `_codex/system/` are **always loaded** at the start of every agent interaction. This is the "hot zone" — essential context that's needed regardless of the task. See `_codex/system/README.md` for rules on what goes here.

### Searching

Follow the tiered loading protocol in `_codex/skills/search.md`:

1. **L0 — Orientation**: Read the domain group README to understand what's in each domain
2. **L1 — Discovery**: Read the target domain folder's README or scan file names
3. **L2 — Deep Read**: Only read full files once you've confirmed relevance

### Creating Notes

1. Use templates from `_codex/templates/` for consistent structure
2. Always include YAML frontmatter (`id`, `title`, `type`, `created`, `tags`)
3. New items go to `06 Inbox (Queue)/` unless you know the exact domain
4. Follow conventions in `_codex/CONVENTIONS.md`

### Organizing

Follow `_codex/skills/organize.md` to move items from the inbox to the right domain.

### Reflection & Maintenance

- Daily: Run `_codex/skills/reflect.md` daily protocol
- Weekly: Run `_codex/skills/reflect.md` weekly protocol
- Monthly: Run `_codex/skills/maintain.md` health audit
- Quarterly: Run `_codex/skills/defrag.md` on cluttered domains

### Agent Self-Organization

Agents are **encouraged** to improve the vault as they use it:

- Update `ai_summary` frontmatter on files they read (if missing or inaccurate)
- Promote frequently-needed files to `_codex/system/` (pinned context)
- Demote `_codex/system/` files that haven't been relevant recently
- Create subfolders when a domain gets cluttered (see `_codex/skills/defrag.md`)
- Fix broken `[[wiki-links]]` when encountered

See `_codex/CONVENTIONS.md` for the full self-organization permission model.

---

## Guardrails

1. **Never create new top-level folders** — use existing 00-99 domains
2. **Never delete files without confirmation** — move to archive or tag as deprecated
3. **Always use frontmatter** — every file needs: `title`, `type`, `created`, `tags`, `ai_summary`
4. **Inbox first** — when unsure where something goes, put it in `06 Inbox (Queue)/`
5. **Wiki-links for connections** — use `[[Note Name]]` to create relationships between notes

---

## For AI Agents — Quick Reference

| Action | Instruction File |
| --- | --- |
| **Always-loaded context** | `_codex/system/` (read all files here first) |
| Search for information | `_codex/skills/search.md` |
| Capture new information | `_codex/skills/capture.md` |
| Organize items into domains | `_codex/skills/organize.md` |
| Run reflection cycles | `_codex/skills/reflect.md` |
| Consolidate daily notes | `_codex/skills/consolidate.md` |
| Restructure cluttered domains | `_codex/skills/defrag.md` |
| Audit vault health | `_codex/skills/maintain.md` |
| Bootstrap a new vault/agent | `_codex/skills/init.md` |
| Create a note from template | `_codex/templates/` |
| Pre-load context for a mode | `_codex/profiles/` |

---

## For Humans — Getting Started

1. Open this folder in Obsidian (or any markdown editor)
2. Read the domain group READMEs to understand the structure
3. Start capturing into `06 Inbox (Queue)/`
4. Follow the CODEX method: Capture → Organize → Direct → Execute → × Multiply
5. Read `CODEX-FRAMEWORK.md` for the full philosophy
