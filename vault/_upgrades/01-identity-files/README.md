---
id: upgrade-01-identity
title: "Upgrade 01: Structured Identity Files"
type: blueprint
status: ready
created: 2026-03-27
owner: zeus
tags: [type/blueprint, status/ready, domain/00-framework]
destination: "_codex/system/"
priority: "🔴 NOW"
---

# Upgrade 01: Structured Identity Files

> **Source**: PAI's TELOS system (17 structured personal context files)
> **Destination**: `_codex/system/` (expands from 2 files → 10 files)
> **Why**: Every agent interaction starts by loading `_codex/system/`. Right now that's just identity + goals. PAI loads mission, beliefs, wisdom, mental models, preferences, and more. This is the single highest-impact upgrade.

---

## What Changes

| Current (`_codex/system/`) | After Upgrade |
|---|---|
| `identity.md` — Basic identity card | `identity.md` — Enhanced with deeper context |
| `active-goals.md` — Q1 goals | `active-goals.md` — Unchanged (already good) |
| *(nothing else)* | `mission.md` — Life mission statement |
| | `beliefs.md` — Core beliefs / world model |
| | `wisdom.md` — Accumulated life wisdom |
| | `preferences.md` — Decision defaults + life preferences |
| | `models.md` — Mental models for decision-making |
| | `frames.md` — Thinking frameworks |
| | `narratives.md` — Personal stories / brand stories |
| | `predictions.md` — Future predictions + bets |

## Token Budget

Current: ~800 tokens (2 files)
After: ~2,000 tokens (10 files, each kept tight at ~200 tokens)

The `_codex/system/README.md` says max 2,000 tokens. We'll respect that by keeping each file to a focused summary. Deep content lives in domain folders and gets referenced.

---

## How to Deploy

1. Copy each `.md` file from this folder to `_codex/system/`
2. Review and personalize each one (templates have `{FILL IN}` markers)
3. Run a token count — total should stay under 2,000
4. Agents will auto-load these on every interaction

## Files in This Module

| File | Purpose | Token Est. |
|------|---------|------------|
| `mission.md` | Life mission statement | ~150 |
| `beliefs.md` | Core beliefs / world model | ~200 |
| `wisdom.md` | Top life lessons learned | ~200 |
| `preferences.md` | Decision defaults | ~200 |
| `models.md` | Mental models | ~200 |
| `frames.md` | Thinking frameworks | ~150 |
| `narratives.md` | Key personal/brand stories | ~200 |
| `predictions.md` | Future bets | ~150 |

---

## PAI Equivalent

PAI stores 17 files in `~/.claude/PAI/USER/TELOS/`:
- MISSION, BELIEFS, WISDOM, BOOKS, MOVIES, LEARNED, WRONG, FRAMES, MODELS, NARRATIVES, STRATEGIES, GOALS, PROJECTS, PROBLEMS, CHALLENGES, PREDICTIONS, TRAUMAS

We adapted this down to 8 new files (plus existing 2 = 10 total) because:
- BOOKS/MOVIES → belong in HP domain, not hot context
- LEARNED/WRONG → merged into WISDOM
- STRATEGIES → merged into MODELS
- PROBLEMS/CHALLENGES → covered by active-goals
- TRAUMAS → too personal for always-loaded context; belongs in HP domain
- PROJECTS → covered by active-goals
