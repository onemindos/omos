---
id: arch-00.2-20260304
title: "OneMind System Primitives & Taxonomy"
type: architecture
status: active
created: 2026-03-04
modified: 2026-03-04
owner: zeus
tags: [type/architecture, status/active, layer/all]
ai_summary: "Complete vocabulary and layer map of all OneMind system primitives — the 55 concepts that make the OS run."
---

# OneMind System Primitives

> **"Every word below is a thing in the system. Not a metaphor. A node."**

This document is the canonical vocabulary of OneMind. Every concept listed here has a real home, a real function, and a real connection to the rest of the system. If it's in this list, it exists. If you can't find it in the Codex, build it.

Assembled from raw ideation (March 2026) — the paper becomes doctrine.

---

## The 7 Layers

The system runs in 7 layers, stacked from **why** to **feedback**.

```
┌─────────────────────────────────────────────────────────────────┐
│  LAYER 7: OUTPUT         results · metrics · analytics          │
├─────────────────────────────────────────────────────────────────┤
│  LAYER 6: COMMAND        command · navigation · launch pad      │
├─────────────────────────────────────────────────────────────────┤
│  LAYER 5: COORDINATION   agents · hivemind · a2a · mcp         │
├─────────────────────────────────────────────────────────────────┤
│  LAYER 4: EXECUTION      capture · workflow · tasks · runs      │
├─────────────────────────────────────────────────────────────────┤
│  LAYER 3: STRUCTURE      areas · projects · atlas · archive     │
├─────────────────────────────────────────────────────────────────┤
│  LAYER 2: INTELLIGENCE   knowledge · memory · rag · ai · ml    │
├─────────────────────────────────────────────────────────────────┤
│  LAYER 1: FOUNDATION     vision · principles · domains · paths  │
└─────────────────────────────────────────────────────────────────┘
```

Data flows **UP** through the layers (raw → meaning → action → result).
Authority flows **DOWN** (vision shapes what gets measured).

---

## Layer 1 — FOUNDATION
*The Why. The north star. Everything else serves this.*

| Primitive | Definition | Codex Home |
|-----------|-----------|------------|
| **vision** | The 50-year picture — what Zeus is building | `00.1 Core-Doctrine/onemind-vision.md` |
| **principles** | The 10 rules the system never breaks | `00.1 Core-Doctrine/Solar-Punk-Principles.md` |
| **missions** | Medium-term directional objectives per domain | `01 Command/` |
| **goals** | Specific, measurable outcomes with deadlines | `01 Command/` (metrics layer) |
| **governance** | Who decides what, at what autonomy level | `00 Framework/COMPLIANCE.md` |
| **guardrails** | Hard limits AI agents cannot cross (L0-L3) | `06 Protocols/` |
| **codex** | This vault — the canonical source of truth | Root `/` |
| **domains** | The 4 pillars: UI · HP · LE · GE | `00 Framework/` |
| **paths** | UI Path · Self Path · Legacy Path · Wealth Path | `00 Framework/` |

---

## Layer 2 — INTELLIGENCE
*The knowing layer. Raw inputs become understanding.*

| Primitive | Definition | Codex Home |
|-----------|-----------|------------|
| **knowledge** | Structured insight extracted from experience | `04 Resources/` |
| **intelligence** | Synthesized knowledge that can drive a decision | `02 Agents/` |
| **memory** | Persistent AI recall — stateful across sessions | `02 Agents/` (Letta blocks) |
| **rag** | Retrieval-Augmented Generation — Codex as AI context | `02 Agents/` (Letta archival) |
| **ai** | The AI agents layer (LEGACY + sub-agents) | `02 Agents/` |
| **ml** | Machine learning models for pattern recognition | `12 Research/` |
| **predictive analysis** | Forward-looking inference from signals + history | `12 Research/` |
| **insights** | Surfaced patterns delivered to Zeus from agents | `01 Command/` |
| **notes** | Raw, unprocessed captures or thinking-in-progress | `06 Inbox/` → `04 Resources/` |

---

## Layer 3 — STRUCTURE
*The skeleton. Where everything lives and how it's labeled.*

| Primitive | Definition | Codex Home |
|-----------|-----------|------------|
| **areas** | Ongoing responsibilities that don't expire | Cross-domain (HP/LE/GE/UI) |
| **projects** | Bounded work with a defined end state | Per-domain project folders |
| **tags** | The cross-cutting label system (path/category) | Obsidian tags + Todoist labels |
| **atlas** | The map — spatial/contextual navigation layer | `05 Environments/` |
| **domains** | (also L1) The 3 domains + 1 operating layer that organize all life | `00-24` · `25-49` · `50-74` · `75-99` |
| **archive** | Completed work, historical records, closed loops | `24 Archive/` |
| **docs** | Formal documents: plans, specs, reports | `04 Resources/` |
| **media** | Images, video, audio — non-text assets | `04 Resources/assets/` |
| **resources** | Reference material, templates, guides | `04 Resources/` |
| **systems** | Recurring operational setups (tech + non-tech) | `08 Infrastructure/` |

---

## Layer 4 — EXECUTION
*The bloodstream. How work actually moves through the system.*

| Primitive | Definition | Codex Home |
|-----------|-----------|------------|
| **capture** | The act of getting anything into the system | `06 Inbox/` |
| **inbox** | The universal entry point — all raw input lands here | `06 Inbox/` |
| **workflow** | A defined sequence that produces a predictable result | `09 Automation/` |
| **tasks** | Discrete, actionable units of work | Todoist (synced) |
| **action** | A single step executed by human or agent | Task subtask / agent run |
| **execution engine** | The runtime that converts tasks → results (N8N + Letta) | `09 Automation/` |
| **runs** | A single execution cycle of a workflow or agent | `09 Automation/` logs |
| **routines** | Scheduled recurring actions (daily/weekly/monthly) | `06 Protocols/` |
| **habits** | Behaviors locked into identity, tracked over time | `25-49 HP/` |
| **protocols** | Formal step-by-step SOPs for specific situations | `06 Protocols/` |

---

## Layer 5 — COORDINATION
*The nervous system. How agents, tools, and systems talk to each other.*

| Primitive | Definition | Codex Home |
|-----------|-----------|------------|
| **agents** | AI entities with memory, tools, and autonomy scope | `02 Agents/` |
| **hivemind** | The 4-agent collective (UI · HP · LE · GE) sharing memory | `02 Agents/` |
| **a2a** | Agent-to-Agent communication protocol | `02 Agents/` |
| **mcp** | Model Context Protocol — tool interface standard | `02 Agents/` |
| **tools** | Functions agents can call (Todoist, Discord, HA, etc.) | `02 Agents/letta-tools/` |
| **skills** | Packaged tool sets for specific agent capabilities | `02 Agents/letta-tools/` |
| **interfaces** | How Zeus talks to the system (Flutter, voice, WhatsApp) | `08 Infrastructure/` |
| **signals** | Real-time data events from the environment | `08 Infrastructure/` |
| **realtime** | The live data layer — sensor feeds, presence, alerts | `08 Infrastructure/` |
| **infrastructure** | The underlying stack: servers, APIs, networking | `08 Infrastructure/` |

---

## Layer 6 — COMMAND
*The cockpit. How Zeus sees and steers the system.*

| Primitive | Definition | Codex Home |
|-----------|-----------|------------|
| **command** | Zeus's primary interface to direct the system | `01 Command/` |
| **navigation** | How to move between domains, views, and contexts | `01 Command/` |
| **launch pad** | The daily home base — the first screen Zeus opens | `01 Command/` |
| **calendar** | Time-bound scheduling layer across all domains | `01 Command/` (Zoho Calendar) |
| **agenda** | Today's active committed items — the day's contract | `01 Command/` |
| **activity** | Live log of what happened — human + agent actions | `01 Command/` |
| **notifications** | Alerts routed to Zeus from agents or signals | `08 Infrastructure/` (Discord/Telegram) |
| **inbox** | (also L4) The view of what just arrived | `06 Inbox/` |

---

## Layer 7 — OUTPUT
*The feedback loop. Did it work? What do we know now?*

| Primitive | Definition | Codex Home |
|-----------|-----------|------------|
| **results** | The actual outcome of a run, task, or project | Per-domain project folders |
| **metrics** | Quantified measurements of system health + progress | `01 Command/` |
| **analytics** | Processed metrics → patterns → intelligence | `01 Command/` / `12 Research/` |
| **insights** | (also L2) Delivered intelligence — "here's what matters" | `01 Command/` |
| **goals** | (also L1) Goals live in L1 but are measured in L7 | `01 Command/` |

---

## The Full Map — All 55 Primitives

```
LAYER   PRIMITIVE              TYPE            HOME
──────  ─────────────────────  ──────────────  ─────────────────────────
L1      vision                 north star      00.1 Core-Doctrine
L1      principles             rules           00.1 Core-Doctrine
L1      missions               direction       01 Command
L1      goals                  target          01 Command
L1      governance             authority       00 Framework
L1      guardrails             limits          06 Protocols
L1      codex                  vault           Root /
L1      domains                domain groups   00 → 99
L1      paths                  trajectories    00 Framework

L2      knowledge              structured info 04 Resources
L2      intelligence           synthesized     02 Agents
L2      memory                 stateful recall 02 Agents (Letta)
L2      rag                    retrieval       02 Agents (Letta)
L2      ai                     agent layer     02 Agents
L2      ml                     models          12 Research
L2      predictive analysis    inference       12 Research
L2      insights               surfaced intel  01 Command
L2      notes                  raw capture     06 Inbox → 04 Resources

L3      areas                  ongoing scope   HP/LE/GE/UI folders
L3      projects               bounded work    Per-domain
L3      tags                   labels          Obsidian + Todoist
L3      atlas                  map             05 Environments
L3      archive                history         24 Archive
L3      docs                   formal text     04 Resources
L3      media                  assets          04 Resources/assets
L3      resources              reference       04 Resources
L3      systems                setups          08 Infrastructure

L4      capture                input           06 Inbox
L4      inbox                  entry point     06 Inbox
L4      workflow               sequence        09 Automation
L4      tasks                  work units      Todoist
L4      action                 single step     Tasks / Agent runs
L4      execution engine       runtime         09 Automation (N8N+Letta)
L4      runs                   cycle           09 Automation logs
L4      routines               scheduled       06 Protocols
L4      habits                 identity        25-49 HP
L4      protocols              SOPs            06 Protocols

L5      agents                 AI entities     02 Agents
L5      hivemind               collective      02 Agents
L5      a2a                    agent comms     02 Agents
L5      mcp                    tool protocol   02 Agents
L5      tools                  callable funcs  02 Agents/letta-tools
L5      skills                 packaged sets   02 Agents/letta-tools
L5      interfaces             entry channels  08 Infrastructure
L5      signals                live events     08 Infrastructure
L5      realtime               live layer      08 Infrastructure
L5      infrastructure         the stack       08 Infrastructure

L6      command                cockpit         01 Command
L6      navigation             movement        01 Command
L6      launch pad             home base       01 Command
L6      calendar               time layer      01 Command
L6      agenda                 today contract  01 Command
L6      activity               live log        01 Command
L6      notifications          alerts          08 Infrastructure

L7      results                outcomes        Per-domain
L7      metrics                measurements    01 Command
L7      analytics              patterns        01 Command / 12 Research
```

---

## The Flow

```
VISION (L1)
  └─ shapes what KNOWLEDGE (L2) is organized into
       └─ stored in STRUCTURE (L3) — areas, projects, atlas
            └─ moved by EXECUTION (L4) — capture → workflow → tasks → runs
                 └─ coordinated by AGENTS (L5) via A2A · MCP · tools
                      └─ commanded by Zeus via COMMAND layer (L6)
                           └─ measured in OUTPUT (L7) → feedback → back to L1
```

---

## Where Primitives Live in the Vault

| Codex Folder | Primitives that live here |
|--------------|--------------------------|
| `00 Framework` | codex, vision, principles, governance, domains, paths, guardrails |
| `01 Command` | command, launch pad, navigation, agenda, calendar, activity, metrics, analytics, insights, goals, missions |
| `02 Agents` | agents, hivemind, a2a, mcp, tools, skills, memory, rag, ai, intelligence |
| `06 Protocols` | protocols, routines, guardrails |
| `04 Resources` | knowledge, docs, notes, media, resources |
| `05 Environments` | atlas, signals |
| `06 Inbox` | capture, inbox, notes |
| `08 Infrastructure` | infrastructure, interfaces, realtime, signals, notifications, systems |
| `09 Automation` | workflow, execution engine, runs |
| `12 Research` | ml, predictive analysis, analytics |
| `24 Archive` | archive |
| `HP / LE / GE` | habits, areas, projects, results |

---

*"The paper becomes doctrine. The doctrine becomes the system. The system runs the life."*


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
