---
name: entity-system
description: Classify any content into one of 11 One Mind entity types (goal, project, task, metric, skill, tool, human, agent, calendar, knowledge, asset), route to the correct pillar and domain, and generate correct frontmatter. Works offline, agent-agnostic.
skill: entity-system
version: 1.0.0
agent: any
requires: [read_file, write_file, list_directory]
tools: none_required
offline: true
pillar: SO
domain: "01 Intelligence (AI, Memory & Tools)"
created: 2026-05-04
tags: [type/skill, zone/so, status/active]
ai_summary: "Classify any content as one of 11 entity types, route to correct pillar/domain, generate frontmatter."
---

# SKILL: Entity System

> Classify any piece of content into the One Mind entity system. Route it to the correct pillar and domain. Generate correct frontmatter. Works offline. Agent-agnostic.

---

## Step 1 — Identify the Entity Type

Read the content. Ask these questions in order until one matches:

```
Is it something you want to ACHIEVE with a deadline?
    → GOAL 🎯

Is it a collection of work toward one outcome, with a start and end?
    → PROJECT 📁

Is it one specific action someone needs to do?
    → TASK ✅

Is it a number you track over time to know if something is working?
    → METRIC 📊

Is it a set of instructions any agent or human can follow repeatedly?
    → SKILL ⚡

Is it an app, platform, API, or system used to do work?
    → TOOL 🔧

Is it a person — real, named, with a relationship to Zeus?
    → HUMAN 👤

Is it an AI actor with a defined role and permissions?
    → AGENT 🤖

Is it a time block, scheduled event, or recurring ritual?
    → CALENDAR 📅

Is it information — a fact, concept, reference, or lesson?
    → KNOWLEDGE 📖

Is it something owned that holds financial, creative, or physical value?
    → ASSET 💎
```

If content spans multiple types, classify by PRIMARY purpose. A project plan that contains tasks is a PROJECT. A how-to guide is KNOWLEDGE. A person's contact record is HUMAN.

---

## Step 2 — Route to the Correct Pillar

```
Content about Zeus himself — health, mind, money, skills, identity?
    → HP (Holistic Performance) — Personal pillar

Content about relationships, family, home, community, heritage?
    → LE (Legacy Evolution) — Home pillar

Content about career, business, brand, ventures, revenue, wealth?
    → GE (Generational Endeavor) — Business pillar

Content about AI agents, tech, tools, automation, infrastructure, protocols?
    → SO (Sovereign Operations) — Operations pillar
```

**Cross-pillar rule:** If content serves two pillars equally, route to the pillar where it will be ACTED ON most often. A business contact (GE) is different from a family member (LE) even though both are HUMAN type.

---

## Step 3 — Route to the Correct Domain

Once the pillar is known, identify the domain number using this map:

### SO — Operations Pillar (00–24)
| Domain | What Goes Here |
|--------|---------------|
| 00 Framework | Doctrine, architecture, philosophy |
| 01 Intelligence | Agents, skills, tools, AI models, memory |
| 02 Infrastructure | VPS, servers, APIs, networking |
| 03 Command | Dashboards, SOPs, runbooks |
| 04 Security | Cybersecurity, credentials, access |
| 05 Development | Code, scripts, testing, DevOps |
| 06 Research | R&D, emerging tech, competitive analysis |
| 08 OM Assets | Entity registry, user accounts |
| 99 Inbox | Unrouted captures — triage queue |

### HP — Personal Pillar (25–49)
| Domain | What Goes Here |
|--------|---------------|
| 25 Identity | Self concept, values, beliefs, purpose |
| 26 Spirituality | Faith, philosophy, meaning |
| 27 Mind | Learning, courses, books, education |
| 28 Body | Health, fitness, nutrition, medical |
| 29 Life Systems | Life admin, legal, personal records |
| 30 Finance | Money — banking, credit, investments, tax |
| 31 Mastery | Skills — physical, technical, creative, personal |
| 32 Rhythm | Routines, productivity, time, energy |
| 33 Relationships | Personal network — friends, mentors, dating |
| 34 Joy | Recreation, travel, entertainment, hobbies |

### LE — Home Pillar (50–74)
| Domain | What Goes Here |
|--------|---------------|
| 50 Home | Property, sanctuary, household |
| 51 Marriage | Spouse, partnership, marriage |
| 52 Offspring | Children — development, education, activities |
| 53 Parents | Parents, elders, caregiving |
| 54 Heritage | Family history, genealogy, legacy |
| 55 Health (Family) | Family health — medical, wellness |
| 56 Wealth (Family) | Family finances — estate, wills, trusts |
| 57 Community | Extended network, neighbors, local |

### GE — Business Pillar (75–99)
| Domain | What Goes Here |
|--------|---------------|
| 75 Brand | Zeus Delacruz personal brand + products |
| 76 Real Estate | Property investments, landlord |
| 77 Acquisitions | Bought businesses |
| 78 Holdings | Legal entities, corporate structure |
| 79 Investments | Portfolio, equity, wealth vehicles |
| 80 Innovation | New venture pipeline |
| 81 Strategy | Portfolio intelligence, long game |
| 82 Partnerships | Network, deals, capital partners |
| 83 Philanthropy | Giving, impact |

---

## Step 4 — Generate Frontmatter

Use this template. Fill every field. Do not skip `ai_summary` or `related`.

```yaml
---
id: [type]-[slug]-[YYYYMMDD]
title: "[Human readable title]"
type: [goal|project|task|metric|skill|tool|human|agent|calendar|knowledge|asset]
status: [active|complete|paused|archived|planned]
created: YYYY-MM-DD
updated: YYYY-MM-DD
pillar: [SO|HP|LE|GE]
domain: "[XX Domain Name]"
ai_summary: "[One sentence — what this is and why it matters.]"
agent: [assigned agent name or 'any']
owner: [zeus or agent name]
tags: [type/[entity], zone/[pillar-lowercase], domain/[pillar-lowercase], status/[status]]
related: ["[[RelatedFile1]]", "[[RelatedFile2]]"]
depth: [1|2|3|4]
---
```

### Depth values:
- `1` = Pillar root or domain root README
- `2` = Top-level domain file
- `3` = Subfolder file
- `4` = Deep nested or highly specific file

### Tag format:
```
type/goal          type/project        type/task
type/metric        type/skill          type/tool
type/human         type/agent          type/calendar
type/knowledge     type/asset

zone/so            zone/hp             zone/le             zone/ge
domain/so          domain/hp           domain/le           domain/ge

status/active      status/complete     status/paused
status/archived    status/planned
```

---

## Step 5 — Place the File

File path convention:
```
[pillar folder]/[domain folder]/[subfolder if needed]/[slug].md
```

Examples:
```
25-49 HP (Holistic Performance)/30 Finance (Money)/Goals (Financial Goals)/credit-score-750-goal.md
75-99 GE (Generational Empire)/75 Zeus Delacruz (Brand & Empire)/Projects (Active)/one-mind-community-launch.md
00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/skills (Agent Skills)/onemind (One Mind System Skills)/[skill-name]/SKILL.md
```

If the correct subfolder does not exist, create it with a `(Descriptor)` in the name and add a `_README.md`.

If you are unsure where to route: place in `99 Inbox (Queue) (The Dock)/` with `status: planned` and a note in `ai_summary` indicating it needs routing.

---

## Entity-Specific Frontmatter Extensions

### GOAL
```yaml
horizon: Q3-2026
metric: "target outcome description"
lead_metric: "what you control"
lag_metric: "what you measure"
linked_project: "[[Project Name]]"
```

### PROJECT
```yaml
goal: "[[Goal Name]]"
start: YYYY-MM-DD
due: YYYY-MM-DD
todoist_project_id: "xxxxx"
```

### TASK
```yaml
project: "[[Project Name]]"
priority: p1
due: YYYY-MM-DD
tid: "todoist_task_id"
```

### METRIC
```yaml
metric_type: lag
current: 0
target: 0
unit: "unit description"
cadence: monthly
tracked_in: "[[tracker file]]"
```

### SKILL
```yaml
requires: [read_file, write_file, list_directory]
tools: none_required
offline: true
```

### TOOL
```yaml
category: AI
platform: claude
access: API
auth: "1Password SO vault"
docs: "[[docs link]]"
```

### HUMAN
```yaml
role: contact
relationship: friend
last_contact: YYYY-MM-DD
birthday: YYYY-MM-DD
```

### AGENT
```yaml
emoji: 🔮
domain_scope: [HP]
autonomy: L2
model: claude-sonnet-4.6
workspace: "~/.openclaw/workspace-oracle"
```

### CALENDAR
```yaml
event_type: ritual
recurrence: weekly
day: Sunday
time: "19:00"
duration: 60min
linked_skill: "[[skill name]]"
```

### ASSET
```yaml
category: financial
value: 0
currency: USD
last_valued: YYYY-MM-DD
tracked_in: "[[Net Worth Tracker]]"
```

---

## Quick Classification Reference

| If you have... | Type | Pillar hint |
|---------------|------|------------|
| "I want to achieve X by Y date" | Goal | Any |
| "I'm working on X project" | Project | Any |
| "I need to do X" | Task | Any |
| "Track my X score/number" | Metric | Any |
| "How to do X (repeatable)" | Skill | SO |
| "This app / platform" | Tool | SO |
| "This person" | Human | HP/LE/GE |
| "This AI agent" | Agent | SO |
| "This event / time block" | Calendar | HP/LE/GE |
| "This fact / concept / reference" | Knowledge | Domain match |
| "I own this / it has value" | Asset | LE/GE/HP |

---

## Validation Checklist

Before filing, verify:

- [ ] Entity type is assigned
- [ ] Pillar is assigned (SO / HP / LE / GE)
- [ ] Domain number and name is correct
- [ ] `ai_summary` is one clear sentence
- [ ] `related` links at least one connected file (if any exist)
- [ ] `depth` is set correctly
- [ ] Tags follow `type/`, `zone/`, `domain/`, `status/` format
- [ ] File is placed in the correct folder path
- [ ] If folder is new — `_README.md` was created with same frontmatter pattern

---

*One Mind Entity System Skill v1.0 — created 2026-05-04. Agent-agnostic. Offline-capable.*
