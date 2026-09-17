---
id: upgrades-master-plan
title: "OneMind Enhancement Upgrades"
type: blueprint
status: active
created: 2026-03-27
modified: 2026-03-30
owner: zeus
tags: [type/blueprint, status/active, domain/00-framework]
source: "Adapted from Daniel Miessler's Personal AI Infrastructure (PAI)"
---

# 🔧 OneMind Enhancement Upgrades

> Adapted from PAI's architecture patterns, rebuilt for OneMind's human enhancement framework.
> These are NOT copies — they're OneMind-native implementations inspired by PAI's best ideas.

---

## What This Is

This folder contains **10 upgrade modules** — blueprints that were (or will be) integrated into the OneMind Codex vault. Each module is self-contained with:

- `README.md` — What it is, why it matters, where it goes in the vault
- Implementation files — Ready to move to their final vault locations (or reference files for deployed modules)
- Mermaid diagrams — Render natively in Obsidian

**9 of 10 modules are deployed.** Module 01 (Identity Files) is staged and waiting for Zeus to personalize `{FILL IN}` markers before deployment.

See also: [`ENHANCEMENT-SYSTEM-GUIDE.md`](../ENHANCEMENT-SYSTEM-GUIDE.md) at vault root for the full operational guide.

## Origin

Daniel Miessler's [Personal AI Infrastructure](https://github.com/danielmiessler/Personal_AI_Infrastructure) (PAI) is a developer-focused CLI tool for personal AI. We analyzed its architecture and extracted 10 patterns worth adopting — then rebuilt each one for OneMind's 100-domain, multi-generational, Obsidian-native framework.

**PAI = tool for developers. OneMind = framework for humans.**

---

## The 10 Upgrades

| # | Module | What It Does | Destination | Status |
|---|--------|-------------|-------------|--------|
| 01 | [Identity Files](01-identity-files/) | Structured personal context always loaded by agents | `_codex/system/` | ⏳ STAGED — needs personalization |
| 02 | [Memory Architecture](02-memory-architecture/) | 3-tier persistent memory (work/learning/relationship) | `_codex/memory/` | ✅ DEPLOYED |
| 03 | [Visual Diagrams](03-visual-diagrams/) | Architecture diagrams + image generation briefs | `_codex/assets/diagrams/` | ✅ DEPLOYED |
| 04 | [Skill Packs](04-skill-packs/) | Self-installing modular skill system (11 packs) | `_codex/skills/` | ✅ DEPLOYED |
| 05 | [Lifecycle Hooks](05-lifecycle-hooks/) | Event-driven automation on CODEX workflow stages | `_system/hooks/` | ✅ DEPLOYED |
| 06 | [Security Gate](06-security-gate/) | Pre-execution validation + audit trail | `_system/security/` | ✅ DEPLOYED |
| 07 | [Agent Traits](07-agent-traits/) | Composable trait system for domain agents | `_codex/profiles/` | ✅ DEPLOYED |
| 08 | [Signal Capture](08-signal-capture/) | Explicit + implicit feedback learning loop | `_codex/memory/learning/` | ✅ DEPLOYED |
| 09 | [Auto-Backup](09-auto-backup/) | Timestamped snapshots before critical changes | `_system/backup/` | ✅ DEPLOYED |
| 10 | [Report Generation](10-report-generation/) | Quarterly reviews + domain dashboards | `_codex/skills/report/` | ✅ DEPLOYED |

---

## Deployment History

| Phase | Modules | Commit | Date |
|-------|---------|--------|------|
| Phase 1 | Blueprints 01-10 staged | `dcf0c276` | 2026-03-27 |
| Phase 1b | Diagrams (03) deployed to `_codex/assets/diagrams/` | `eb2cfb78` | 2026-03-27 |
| Phase 2-4 | Modules 02, 04-10 deployed to vault | `28356c78` | 2026-03-27 |
| Docs | Enhancement System Guide created | `81eeae14` | 2026-03-27 |
| Docs | AGENTS.md + CLAUDE.md updated | `e1da988a` | 2026-03-28 |

## What's Left

1. **Module 01 — Identity Files**: Zeus needs to fill in the `{FILL IN}` markers in `01-identity-files/*.md`, then deploy to `_codex/system/`
2. **Image briefs**: `03-visual-diagrams/image-briefs/` contains Midjourney/DALL-E prompts for polished graphics (optional)
3. **Delete `_upgrades/`**: Once Module 01 is deployed and image briefs are used, this staging folder can be removed

---

## PAI vs OneMind Comparison

| Dimension | PAI | OneMind |
|-----------|-----|---------|
| **Interface** | Terminal (CLI) | Obsidian (Visual) |
| **Scope** | Digital AI assistant | 100-domain life framework |
| **Users** | Single developer | Multi-generational family |
| **Memory** | JSON files | Markdown + Letta agents |
| **Skills** | TypeScript packs | Markdown skill definitions |
| **Identity** | 17 TELOS files | Structured `_codex/system/` |
| **Security** | Hook-based validator | Policy + audit system |
| **Agents** | Trait-composed | Domain-assigned (SO/HP/LE/GE) |
| **Output** | Terminal responses | Obsidian notes + dashboards |
| **Philosophy** | Context-first | Enhancement-first |

> OneMind doesn't compete with PAI. PAI could run *inside* OneMind's SO layer.
