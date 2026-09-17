---
id: upgrade-04-skill-packs
title: "Upgrade 04: Self-Installing Skill Packs"
type: blueprint
status: deployed
created: 2026-03-27
deployed: 2026-03-27
owner: zeus
tags: [type/blueprint, status/deployed, domain/00-framework]
destination: "_codex/skills/"
---

# Upgrade 04: Self-Installing Skill Packs

> ✅ **DEPLOYED** to `_codex/skills/` — commit `28356c78` (2026-03-27)

> **Source**: PAI's Pack system (13 self-installing, composable skill modules)
> **Destination**: `_codex/skills/`
> **Why**: Current skills are single `.md` files with instructions. PAI's pattern adds INSTALL/VERIFY steps that let any AI agent self-install the skill. Makes skills portable, shareable, and community-ready.

---

## Current vs Upgraded

| Current | Upgraded |
|---------|----------|
| `capture.md` (flat instructions) | `capture/README.md` + `SKILL.md` + `INSTALL.md` + `VERIFY.md` |
| Agent reads one file | Agent reads README, installs via INSTALL, verifies via VERIFY |
| Not shareable | Community can fork individual packs |
| No verification | VERIFY.md confirms skill loaded correctly |

## Pack Structure Template

```
_codex/skills/{skill-name}/
├── README.md       ← What this skill does, why, how it works
├── SKILL.md        ← The actual skill instructions (what agents execute)
├── INSTALL.md      ← How an AI agent loads this skill into its context
└── VERIFY.md       ← Checklist to confirm the skill is working
```

## The 10 OneMind Skill Packs

| # | Pack | CODEX Stage | Current File |
|---|------|-------------|--------------|
| 1 | Capture | C | `capture.md` |
| 2 | Organize | O | `organize.md` |
| 3 | Direct | D | *(new)* |
| 4 | Execute | E | *(new)* |
| 5 | Multiply | × | `reflect.md` (renamed) |
| 6 | Search | — | `search.md` |
| 7 | Consolidate | — | `consolidate.md` |
| 8 | Defrag | — | `defrag.md` |
| 9 | Init | — | `init.md` |
| 10 | Maintain | — | `maintain.md` |

## How to Deploy

1. For each skill, create a folder: `_codex/skills/{name}/`
2. Move existing `.md` → `SKILL.md` inside the folder
3. Add `README.md`, `INSTALL.md`, `VERIFY.md` from templates
4. Test: Ask an agent to "install the capture skill" — it should read INSTALL.md and self-configure

## PAI Equivalent

PAI has 13 packs (Agents, ContentAnalysis, Investigation, Media, Research, Scraping, Security, Telos, Thinking, USMetrics, Utilities, ContextSearch). Each has README/INSTALL/VERIFY/src. We're applying the same structural pattern to OneMind's existing 8 skills + 2 new ones.
