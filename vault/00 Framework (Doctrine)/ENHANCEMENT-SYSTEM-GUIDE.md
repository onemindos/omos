---
id: upgrade-system-guide
title: "OneMind Enhancement System — Complete Guide"
type: reference
status: active
created: 2026-03-28
owner: zeus
tags: [type/reference, status/active, domain/00-framework]
ai_summary: "Complete guide to all 10 OneMind enhancement modules — what they are, what they do, where they live."
---

# OneMind Enhancement System — Complete Guide

> 10 modules adapted from Daniel Miessler's Personal AI Infrastructure (PAI), rebuilt for OneMind's human enhancement framework. Deployed March 2026.

---

## Module 01 — Structured Identity Files

**Location**: `_upgrades/01-identity-files/` → deploys to `_codex/system/`
**Status**: Staged (ready to personalize and move)

**What it is**: 8 new markdown files that define you — mission, beliefs, wisdom, preferences, mental models, thinking frames, personal narratives, and predictions.

**What it does**: Right now, when any AI agent starts a session, it loads 2 files about you (identity + goals). After deployment, it loads 10 files. Every agent — Sentinel, Coach, Haven, Empire — starts every conversation already knowing your life mission, how you make decisions, what you believe, how you like to communicate, and what stories define your brand. PAI calls this "TELOS." We call it system context.

**Why it matters**: Without this, agents ask you the same questions repeatedly. With this, they already know the answer.

**Files**:
- `mission.md` — Life mission statement
- `beliefs.md` — Core beliefs / world model
- `wisdom.md` — Accumulated life lessons
- `preferences.md` — Decision defaults, communication style
- `models.md` — Mental models for decision-making
- `frames.md` — Thinking frameworks for different problem types
- `narratives.md` — Origin story, brand stories, counter-positioning
- `predictions.md` — Future bets and confidence levels

---

## Module 02 — 3-Tier Memory Architecture

**Location**: `_codex/memory/` ✅ LIVE
**Status**: Deployed

**What it is**: A persistent memory system with 3 layers:
- **Work** (`work/current-focus.json`) — What's happening right now. Active projects, current session state.
- **Learning** (`learning/signals.jsonl` + `patterns.md`) — What the system has learned over time. Feedback, patterns, lessons.
- **Relationship** (`relationship/interaction-notes.md`) — How you like to work. Communication style, trust levels, preferences discovered over time.

**What it does**: Gives agents memory between sessions. Today, every conversation starts from scratch. With this, agents remember what you're working on (Work), what you liked and didn't like (Learning), and how to talk to you (Relationship). The `signals.jsonl` file is append-only — it grows forever, building a richer picture of your preferences with every interaction.

**Why it matters**: This is what makes OneMind *learn*. Without memory, it's just a chatbot. With memory, it's a partner.

---

## Module 03 — Visual Architecture Diagrams

**Location**: `_codex/assets/diagrams/` ✅ LIVE
**Status**: Deployed (Mermaid renders in Obsidian)

**What it is**: 8 Mermaid diagrams that render live in Obsidian:
1. **System Architecture** — Full 6-layer stack (Human → Intelligence → Agent → Knowledge → Operations → Platform)
2. **CODEX Method** — The C→O→D→E→× flow with substeps for each stage
3. **Domain Map** — All 100 domains organized across the 3+1 groups
4. **Agent Architecture** — How the 4 agents are structured, their traits, how they coordinate
5. **Memory Architecture** — The 3-tier memory flow with session lifecycle
6. **Data Flow** — How information moves through the entire system end-to-end
7. **OneMind vs Second Brain** — Side-by-side feature comparison
8. **The Equation** — HUMAN + AI = ONE MIND visual

**What it does**: Open in Obsidian and see diagrams rendered visually. No plugins needed. Use in book, course, community onboarding, pitch decks, and README. Image generation briefs with Midjourney/DALL-E prompts are in `_upgrades/03-visual-diagrams/image-briefs/all-briefs.md`.

**Why it matters**: PAI has 23 professional diagrams. You can't teach, sell, or explain a framework without visuals.

---

## Module 04 — Self-Installing Skill Packs

**Location**: `_codex/skills/` ✅ LIVE
**Status**: Deployed

**What it is**: Upgraded the old flat skill files into structured "packs." Each skill is now a folder with 4 files:
- `README.md` — What the skill does, when to use it
- `SKILL.md` — The actual instructions agents follow
- `INSTALL.md` — How an AI agent loads and activates this skill
- `VERIFY.md` — Checklist to confirm the skill is working correctly

**11 packs total**:

| Pack | CODEX Stage | What It Does |
|------|-------------|-------------|
| **Capture** | C | How to add new information — classify, template, route to inbox |
| **Organize** | O | How to route content to the right domain folder |
| **Direct** | D | How to prioritize tasks against goals, assign to the right agent |
| **Execute** | E | How to load context, do the work, document output, save properly |
| **Reflect** | × | How to run daily/weekly/monthly reviews and extract patterns |
| **Search** | Utility | How to find info using the 4-level tiered loading system |
| **Consolidate** | Utility | How to merge duplicate knowledge and reduce redundancy |
| **Defrag** | Utility | How to clean up the vault — fix broken links, stale files, orphans |
| **Init** | Utility | How to initialize a new vault or onboard a new agent |
| **Maintain** | Utility | How to keep the vault healthy — frontmatter checks, tag hygiene |
| **Report** | × / Utility | How to generate structured reports (see Module 10) |

**What it does**: Skills are now self-installing and community-shareable. Any AI agent reads `INSTALL.md` and configures itself. The CODEX loop is **complete** — Direct and Execute were missing before; now every stage has a dedicated skill.

**Why it matters**: This is what makes OneMind portable. PAI's packs are why people can install it in 5 minutes.

---

## Module 05 — Lifecycle Hooks

**Location**: `_system/hooks/` ✅ LIVE
**Status**: Deployed

**What it is**: An event-driven system with 10 hooks that fire at specific moments in the CODEX workflow:

| Hook | When It Fires | What It Does |
|------|---------------|-------------|
| **LoadContext** | Session start | Loads all identity files from `_codex/system/` |
| **LoadMemory** | Session start | Loads memory tiers for continuity |
| **ClassifyInput** | Capture stage | Auto-determines input type (note, decision, SOP, etc.) |
| **RouteToDomain** | Organize stage | Suggests the correct domain folder |
| **CheckGoals** | Direct stage | Compares task against active goals — flags misalignment |
| **SecurityValidate** | Before execution | Checks command safety against security patterns |
| **LogAction** | After execution | Records what was done for audit trail |
| **ExtractLearning** | Multiply stage | Pulls patterns and insights from completed work |
| **SaveMemory** | Session end | Persists memory updates to vault |
| **CaptureSignal** | Session end | Logs feedback (explicit + implicit) |

**What it does**: Makes the system event-driven instead of manual. The `hooks.yaml` registry lists every event and which hooks fire. Agents don't forget to save memory because it's a hook, not a hope.

**Why it matters**: PAI has 22 hooks that make it feel intelligent. This gives OneMind the same "it just knows" quality.

---

## Module 06 — Security Validation Gate

**Location**: `_system/security/` ✅ LIVE
**Status**: Deployed

**What it is**: A security layer with 3 files:
- `patterns.yaml` — Rules: what to BLOCK, CONFIRM, or ALERT on
- `validator.md` — Step-by-step validation instructions for agents
- `events/` — Security event logs (JSONL)

**What it protects**:
- **Blocks**: `rm -rf /`, repo deletion, database drops, fork bombs
- **Confirms**: Force pushes, data deletion, modifying root docs
- **Alerts**: Remote code execution (`curl | sh`), sudo, world-writable permissions
- **Path protection**: SSH keys are zero-access. `_codex/system/` requires confirmation. Root folders can never be deleted.
- **Injection defense**: If external content contains "IGNORE PREVIOUS INSTRUCTIONS," agent STOPS and alerts you.

**What it does**: Prevents agents from doing catastrophic things. Every action gets logged.

**Why it matters**: You're storing personal, family, and financial data across 100 domains. One bad command could wreck it.

---

## Module 07 — Agent Trait Composition

**Location**: `_codex/profiles/` ✅ LIVE
**Status**: Deployed

**What it is**: A composable trait system with 3 dimensions:
- **Expertise** (7 traits): systems-thinking, wellness, family-focused, business, security, research, creative
- **Personality** (7 traits): analytical, encouraging, patient, bold, skeptical, enthusiastic, measured
- **Approach** (7 traits): methodical, consistent, protective, strategic, thorough, rapid, systematic

**4 agent profiles**:
- **⚡ Sentinel** (SO 00-24) — systems-thinking + analytical + methodical
- **💪 Coach** (HP 25-49) — wellness + encouraging + consistent
- **🏡 Haven** (LE 50-74) — family-focused + patient + protective
- **💰 Empire** (GE 75-99) — business + bold + strategic

**What it does**: Each agent has a distinct personality matching its domain. You can override traits on the fly: "Coach, use bold + rapid for this workout plan." Custom traits go in `_codex/profiles/custom/`.

**Why it matters**: PAI's agents feel like different people. This gives your 4 agents distinct identities.

---

## Module 08 — Signal Capture System

**Location**: Part of `_codex/memory/learning/` ✅ LIVE
**Status**: Deployed

**What it is**: A feedback loop that captures how you respond to agent work:
- **Explicit signals**: "great work" (positive), "8/10" (numeric), "not what I meant" (negative)
- **Implicit signals**: Accepted without changes (positive), asked for revision (negative), expanded enthusiastically (interest)

Logged to `signals.jsonl` as single-line JSON.

**What it does**: Every week during × Multiply reviews, the system processes signals and updates `patterns.md` — "Zeus prefers batch operations," "Empire gets the best ratings on brand strategy." Over months, this becomes a rich dataset of what works.

**Why it matters**: This is the *learning* in "One Mind learns." Without signals, the system stays the same. With signals, it improves every week.

---

## Module 09 — Auto-Backup System

**Location**: `_system/backup/` ✅ LIVE
**Status**: Deployed

**What it is**: Pre-modification snapshots. Before an agent modifies any critical file, it creates a timestamped copy with a `manifest.json` explaining what was backed up and why.

**Protected files**: ONEMIND-CODEX.md, AGENTS.md, CLAUDE.md, CODEX-FRAMEWORK.md, `_codex/system/*`, `_codex/memory/*`

Backups are git-ignored (local only). Retention: last 10 + all under 7 days + monthly permanents.

**What it does**: Instant rollback without git archaeology.

**Why it matters**: Git is your long-term safety net. This is your instant undo button.

---

## Module 10 — Report Generation

**Location**: `_codex/skills/report/` ✅ LIVE
**Status**: Deployed

**What it is**: A skill pack that generates 4 report types:

1. **Quarterly Life Review** — Scores all 4 domain groups (🟢/🟡/🔴), top accomplishments, lessons from memory, next quarter priorities. Your life's quarterly earnings call.
2. **Domain Health Report** — Zooms into any domain. File counts, last updated, knowledge gaps, connections, AI recommendations. A health checkup for any life area.
3. **Family Legacy Report** — Annual report: family milestones, knowledge captured, financial health, system status, legacy score. Your family's annual report.
4. **Narrative Points** — McKinsey-style 8-section arc (Opening → Problem → Evidence → Stakes → Solution → Requirements → Vision → Close). Instant presentation builder.

**What it does**: Turns raw vault data into actionable intelligence. "Generate my quarterly review" → structured assessment of your entire life.

**Why it matters**: This is × Multiply made tangible. Without reports, reviews are vibes. With reports, reviews are data.

---

## The Complete System

```
Identity (01)  →  defines WHO you are
Memory (02)    →  remembers WHAT happened
Diagrams (03)  →  shows HOW it works visually
Skills (04)    →  defines WHAT agents can do
Hooks (05)     →  automates WHEN things happen
Security (06)  →  protects EVERYTHING
Traits (07)    →  gives agents PERSONALITY
Signals (08)   →  captures FEEDBACK
Backup (09)    →  prevents LOSS
Reports (10)   →  generates INSIGHT
```

---

## Quick Reference

| Module | Location | Status |
|--------|----------|--------|
| 01 Identity | `_upgrades/01-identity-files/` → `_codex/system/` | Staged |
| 02 Memory | `_codex/memory/` | ✅ Live |
| 03 Diagrams | `_codex/assets/diagrams/` | ✅ Live |
| 04 Skills | `_codex/skills/` | ✅ Live |
| 05 Hooks | `_system/hooks/` | ✅ Live |
| 06 Security | `_system/security/` | ✅ Live |
| 07 Traits | `_codex/profiles/` | ✅ Live |
| 08 Signals | `_codex/memory/learning/` | ✅ Live |
| 09 Backup | `_system/backup/` | ✅ Live |
| 10 Reports | `_codex/skills/report/` | ✅ Live |

## Source

Inspired by [Daniel Miessler's Personal AI Infrastructure](https://github.com/danielmiessler/Personal_AI_Infrastructure). PAI is a developer CLI tool. OneMind is a human enhancement framework. We took the best architectural patterns and rebuilt them for Obsidian, markdown, and multi-generational life management.
