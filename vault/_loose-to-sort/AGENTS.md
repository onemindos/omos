# agents.md

Universal instructions for AI coding agents (Claude Code, GitHub Copilot, Cursor, Gemini CLI, OpenCode, Codex, etc.) working in this repository. Read this before touching any file.

---

## What This Repo Is

**One Mind** is an Obsidian vault — the world's first AI-native human enhancement framework. It combines an enhancement method (CODE×), a knowledge architecture (3 life pillars + 1 operating layer × 100 domains), an AI memory layer, and operational code in one repo.

This is NOT a traditional software repository. Most files are markdown. Code lives in `_system/`.

Core equation: `HUMAN (Zeus) + AI (LEGACY) = ONE MIND`

| Context | Path |
|---------|------|
|  Local (Mac) | `/Users/onemind/onemind/` |
| GitHub | `github.com/Zeus-Delacruz/onemind` |
| Obsidian | Open repo root as a vault |

---

## Agent Context Files

Each platform auto-discovers its own context file. Do NOT consolidate these — they serve different roles.

| File | Read By | Purpose |
|------|---------|---------| 
| `agents.md` | **All agents** (universal) | Shared rules, vault structure, conventions |
| `claude.md` | Claude Code | Claude-specific context and instructions |
| `gemini.md` | Gemini CLI | Gemini-specific role and instructions |
| `.hermes.md` | Hermes Agent | Orchestrator role, profile mapping, scheduling |
| `.github/copilot-instructions.md` | GitHub Copilot | Copilot-specific instructions |

**Load order:** Agent reads its own overlay file → then `agents.md` for universal context.

---

## AI Agent Architecture

One Mind runs on **Hermes Agent** on the Mac as the primary orchestration engine, with coding agents as specialist hands.

```
┌─────────────────────────────────────────────────────────────┐
│                     ONE MIND SYSTEM                          │
├─────────────────────────────────────────────────────────────┤
│  Hermes (Brain / Workhorse)                                 │
│  → Legacy AI — primary command agent                        │
│  → Deep work execution, research, synthesis                 │
│  → Multi-agent orchestration                                │
│  → Scheduled automation (cron, recurring tasks)             │
│  → Persistent memory and skills (~/.hermes/)                │
├─────────────────────────────────────────────────────────────┤
│  Coding Agents (Hands)                                      │
│  → Claude Code, Copilot, Gemini CLI, Codex, OpenCode        │
│  → Code tasks, analysis, file operations in this vault      │
└─────────────────────────────────────────────────────────────┘
```

**Communication:** Shared filesystem (this vault) + MCP servers. All agents read and write the same vault.

---

## Repo Structure

```
onemind/
├── 00-24 SO (Sovereign Operations)/   ← Operating layer (AI, tech, infrastructure)
├── 25-49 HP (Holistic Performance)/   ← YOUR SELF (health, mind, skills, money)
├── 50-74 LE (Legacy Evolution)/       ← YOUR FOUNDATION (relationships, community, home)
├── 75-99 GE (Generational Empire)/    ← YOUR EMPIRE (career, ventures, brand, impact)
├── _system/                           ← Operational code (scripts, configs)
├── .hermes.md                         ← Hermes agent context
├── agents.md                          ← This file (universal, all agents)
├── claude.md                          ← Claude Code context
├── gemini.md                          ← Gemini CLI context
├── onemind.md                         ← THE canonical master document
├── codex.md                           ← CODE× method deep spec
└── README.md                          ← GitHub landing page
```

---

## Framework Architecture: 4 Pillars

One Mind is **NOT four equal boxes**. It's three life pillars running on one operating layer.

```
┌───────────────────────────────────────────────────────┐
│           SO — Sovereign Operations (00–24)            │
│                    OPERATIONS                          │
│           The operating layer that runs everything     │
├─────────────────┬──────────────────┬──────────────────┤
│  HP (25–49)     │  LE (50–74)      │  GE (75–99)      │
│  Holistic       │  Legacy          │  Generational    │
│  Performance    │  Evolution       │  Endeavor        │
│   PERSONAL      │      HOME        │    BUSINESS      │
└─────────────────┴──────────────────┴──────────────────┘
```

| Pillar | Range | What Lives Here |
|--------|-------|-----------------|
| **SO** — Operations | 00-24 | AI agents, infrastructure, protocols, automation (the operating layer) |
| **HP** — Personal | 25-49 | Health, mind, skills, money, identity |
| **LE** — Home | 50-74 | Relationships, community, home, environment |
| **GE** — Business | 75-99 | Career, ventures, brand, wealth, impact |

**SO is the platform. HP, LE, GE are the three life pillars.** Route by subject:
- AI agents, tech, automation, security → **SO (00-24)**
- Personal growth, fitness, habits → **HP (25-49)**
- Relationships, community, home → **LE (50-74)**
- Career, ventures, brand, income → **GE (75-99)**

**Never** rename the pillar folders without Zeus's explicit approval.

---

## The CODE× Method

**C**apture → **O**rganize → **D**irect → **E**xecute → **×** Multiply

- × is the multiplication sign (not the letter X, not "eXamine")
- Each review cycle compounds the system's intelligence
- The 5 steps form a LOOP (× feeds back to C)

Full spec: `codex.md`

---

## Working with the Vault

1. **Obsidian is the primary interface** — Open the repo root as an Obsidian vault
2. **Don't rename domain group folders** — SO/HP/LE/GE names are locked
3. **New content goes to Inbox first** — `00-24 SO/24 Inbox (Queue)/`, then route to proper domain
4. **Read before writing** — Many notes have embedded dataview queries and backlinks
5. **Preserve frontmatter** — Files have YAML frontmatter with `id`, `tags`, `owner`, `status`
6. **No secrets in markdown** — Document where credentials live, never the values
7. **Code goes in `_system/`** — Knowledge domains contain markdown only

---

## Secrets

All secrets are in **1Password vault `00-24 SO (Sovereign Operations)`**, vault ID `aaeyris5z44pdqspk277u2oppi`.

Reference format: `op://aaeyris5z44pdqspk277u2oppi/Item Name/field-name`

Do not store real API keys or passwords in any markdown file.

---

## Naming Decisions (FINAL — May 2026)

| ❌ Do NOT use | ✅ Use instead |
|---------------|----------------|
| "OneMind Codex" | "One Mind" |
| "Codex Method" or "CODEX Method" | "CODE× Method" or just "CODE×" |
| "quadrant" or "zone" | "pillar" |
| "four quadrants" | "4 pillars" or "3 life pillars + 1 operating layer" |
| "UI" for SO | "SO", "Sovereign Operations", or "Operations" |
| "Your Family" for LE | "Home" (public) |
| "empire" for GE tagline | "Endeavor" (internal code), "Business" (public) |
| "eXamine" | "× Multiply" |
| "Your Codex" for personal vault | "Your Mind" or "One Mind" |

---

## What NOT to Do

- Do not rename domain group folders (SO, HP, LE, GE) without Zeus's approval
- Do not commit API keys, passwords, or tokens in any markdown file
- Do not modify `.obsidian/` config files without understanding vault impact
- Do not add code files into knowledge domains — code goes in `_system/`
- Do not edit `onemind.md` casually — it is the canonical master document
- Do not use "quadrant", "UI", "enterprise", or "eXamine" — see Naming Decisions above

---

## Product & Distribution Strategy (May 2026)

### What One Mind Sells

The teachable product is NOT the 100-domain architecture. It's the transformation the framework enables. Two core offerings:

1. **The CODE× Method** (5 steps) — Capture → Organize → Direct → Execute → × Multiply. Anyone learns it in 5 minutes. It's a loop. It compounds. It goes beyond CODE (Forte) because it doesn't stop at "Express" — it goes through Execute and multiplies.

2. **The 4 Pillars** — Operations, Personal, Home, Business. The pitch: *"One Mind gives every part of your life a home — and AI runs them all together."*

**One-liner:** *"One Mind is how you and AI run your life together."*
**Counter-positioning:** *"Second Brain stores your knowledge. One Mind runs your life."*

### The 5 Product Layers

```
LAYER 1: THE IDEA (free, everywhere)
"You + AI = One Mind. Here's the 5-step loop."
→ Blog posts, social, podcast appearances

LAYER 2: THE METHOD (manifesto, $0-10)
CODE× method + 3+1 model + first 30 days playbook
→ Short ebook (50-80 pages)

LAYER 3: THE STARTER (free, template)
Empty vault with structure, templates, example content
→ GitHub repo, one-click Obsidian setup

LAYER 4: THE COMMUNITY ($10-30/mo)
Weekly calls, shared implementations, Zeus answers questions
→ Circle or Discord

LAYER 5: THE ARCHITECT (premium, $500+)
Zeus sets up YOUR vault with YOUR domains, YOUR agents
→ 1-on-1 or small cohort
```

### Competitive Position

| Framework | Created | Scope | AI-Native |
|-----------|---------|-------|-----------|
| **GTD** (David Allen) | 2001 | Tasks only | No |
| **PARA** (Tiago Forte) | 2017 | Information only | No |
| **BASB/CODE** (Forte) | 2022 | Knowledge only | No |
| **One Mind** (Zeus) | 2025 | **Entire life** | **Yes** |

The gap: **Nobody has a comprehensive, tool-agnostic, AI-native life framework.**
