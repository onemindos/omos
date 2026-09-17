---
id: one-mind-entity-system-20260504
title: "ONE MIND ENTITY SYSTEM"
type: doctrine
status: active
created: 2026-05-04
modified: 2026-05-04
owner: zeus
tags: [type/doctrine, status/active, framework/core]
ai_summary: "The complete One Mind entity system — 11 entity types, where each lives in the 4 pillars, and how it maps against PARA, CODE, and Second Brain."
---

# ONE MIND ENTITY SYSTEM

> *The entities are the atoms of One Mind. Every piece of information, every action, every person, every tool — everything fits into one of these 11 types.*

---

## The Problem With PARA

PARA gives you 4 buckets:

| PARA | What It Is |
|------|-----------|
| Projects | Active work with an end date |
| Areas | Ongoing responsibilities |
| Resources | Reference material |
| Archives | Inactive stuff |

That is a **filing system**. It tells you WHERE to put things. It does not tell you WHAT things are, WHO acts on them, HOW they connect, or WHETHER they are working.

One Mind goes further. Everything has a type. Types have structure. Structure enables AI to act on them.

---

## The 11 Entity Types

| # | Entity | Symbol | What It Is | PARA Equivalent |
|---|--------|--------|-----------|-----------------|
| 1 | **Goal** | 🎯 | A desired future state with a deadline and metrics | Area (vague) |
| 2 | **Project** | 📁 | A time-bound effort with a clear outcome | Project |
| 3 | **Task** | ✅ | The smallest executable unit of work | (part of Project) |
| 4 | **Metric** | 📊 | A measurable signal — lead or lag | None |
| 5 | **Skill** | ⚡ | A capability — human or AI, portable, reusable | Resource |
| 6 | **Tool** | 🔧 | A system, app, or integration that does work | Resource |
| 7 | **Human** | 👤 | A person — you, a contact, a collaborator | None |
| 8 | **Agent** | 🤖 | An AI actor with a role, domain, and permissions | None |
| 9 | **Calendar** | 📅 | A time block, event, or recurring rhythm | (implicit in PARA) |
| 10 | **Knowledge** | 📖 | An atomic unit of information — a fact, concept, document | Resource / Archive |
| 11 | **Asset** | 💎 | A tangible or digital thing you own that holds value | Resource |

---

## Where Each Entity Lives in the 4 Pillars

```
                    SO          HP          LE          GE
                Operations   Personal     Home       Business
─────────────────────────────────────────────────────────────
Goal              ●           ●           ●           ●      all pillars
Project           ●           ●           ●           ●      all pillars
Task              ●           ●           ●           ●      all pillars
Metric            ●           ●           ●           ●      all pillars
Skill             ●●          ●           ○           ●●     SO + GE heavy
Tool              ●●          ○           ○           ●●     SO + GE heavy
Human             ○           ●●          ●●          ●      HP + LE heavy
Agent             ●●          ○           ○           ○      SO only
Calendar          ○           ●●          ●●          ●●     HP + LE + GE
Knowledge         ●           ●●          ●           ●●     everywhere
Asset             ○           ●           ●●          ●●     LE + GE heavy
```

---

## Entity Definitions — Full Spec

### 1. GOAL 🎯
*A desired future state with a deadline, owner, and success metric.*

```yaml
type: goal
pillar: HP          # HP / LE / GE / SO
domain: "30 Finance (Money)"
status: active      # active / achieved / abandoned / paused
horizon: Q3-2026    # quarter, year, or specific date
owner: zeus         # human or agent responsible
metric: "Credit score ≥ 750"
lead_metric: "Disputes sent per week"
lag_metric: "Score at 90-day check"
agent: Oracle
```

**Where it lives:** Every pillar has goals. HP goals are personal (health, finance, skills). LE goals are relational (marriage, family, home). GE goals are economic (revenue, brand, growth). SO goals are operational (infrastructure, automation, security).

**PARA comparison:** PARA has no goal entity. Everything goes into Areas (vague ongoing) or Projects (time-bound). One Mind separates Goals (outcomes) from Projects (efforts to reach them).

---

### 2. PROJECT 📁
*A time-bound effort with a defined outcome, owner, and linked goal.*

```yaml
type: project
pillar: GE
domain: "75 Zeus Delacruz (Brand & Empire)"
status: active      # active / complete / paused / cancelled
goal: "[[Launch One Mind Community]]"
owner: zeus
agent: Empire
start: 2026-05-01
due: 2026-06-30
todoist_project_id: "xxxxx"
```

**Where it lives:** Projects live inside domains. A project belongs to exactly one pillar. Cross-pillar work is coordinated through SO.

**PARA comparison:** Direct equivalent to PARA Projects. But in One Mind, every Project links to a Goal. In PARA, projects float independently.

---

### 3. TASK ✅
*The smallest executable unit of work. One action. One owner. One outcome.*

```yaml
type: task
project: "[[Credit Repair Sprint]]"
owner: zeus         # or agent name
status: open        # open / in-progress / done / cancelled
priority: p1        # p1 / p2 / p3 / p4
due: 2026-05-10
tid: "todoist_task_id"
```

**Where it lives:** Tasks live inside Projects or as standalone items in a domain's active list. They execute in Todoist; the source of truth is the Codex.

**PARA comparison:** PARA treats tasks as embedded in projects. One Mind promotes tasks to a first-class entity with explicit owner, priority, and sync ID.

---

### 4. METRIC 📊
*A measurable signal — lead (action you control) or lag (outcome you measure).*

```yaml
type: metric
pillar: HP
domain: "30 Finance (Money)"
metric_type: lag    # lead / lag
current: 680
target: 750
unit: "credit score points"
cadence: monthly
tracked_in: "[[Net Worth Tracker]]"
```

**Where it lives:** Metrics live inside domains, inside the `Performance-Metrics` or `_README.md` of any domain. HP metrics = personal KPIs. GE metrics = business KPIs. LE metrics = relationship and home health. SO metrics = system performance.

**PARA comparison:** PARA has no metric entity. This is a complete gap. Without metrics, you never know if PARA is working.

---

### 5. SKILL ⚡
*A portable capability — a set of instructions any agent (human or AI) can follow to accomplish a repeatable outcome.*

```yaml
type: skill
agent: any          # any / Legacy / Oracle / specific agent
requires: [read_file, list_directory, write_file]
tools: none_required
offline: true
pillar: SO          # where it lives
domain: "01 Intelligence (AI, Memory & Tools)"
```

**Where it lives:** Skills live in `SO 01 Intelligence / skills (Agent Skills)/`. Organized by pillar (SO, HP, LE, GE, onemind). Skills are portable — any agent that can read a file can follow them.

**PARA comparison:** PARA has no skill entity. Resources (in PARA) are static reference material. Skills are executable — they produce an outcome when followed.

---

### 6. TOOL 🔧
*A system, application, or API integration that agents and humans use to do work.*

```yaml
type: tool
category: AI        # AI / productivity / finance / communication / infrastructure
platform: claude
access: API
auth: "1Password SO vault"
agent: Legacy
docs: "[[Claude API Docs]]"
```

**Where it lives:** Tools live in `SO 01 Intelligence / tools (Tools & Utilities)/`. Organized by category (AI, productivity, financial, infrastructure). The tool entry is a pointer — it documents what the tool is and how to access it, not the tool itself.

**PARA comparison:** Tools map loosely to PARA Resources. But in One Mind, tools are structured with access info, agent assignments, and API docs — not just a bookmark.

---

### 7. HUMAN 👤
*A person — yourself, a contact, collaborator, mentor, family member, or vendor.*

```yaml
type: human
role: contact       # self / family / friend / mentor / vendor / collaborator
pillar: LE          # HP=self, LE=relationships, GE=professional network
relationship: spouse
tags: [zone/le, type/human]
birthday: YYYY-MM-DD
last_contact: YYYY-MM-DD
notes: "..."
```

**Where it lives:**
- You (Zeus) → `HP 25 Identity (Self & Identity)`
- Family / close relationships → `LE 51-57`
- Professional contacts → `GE 79 Partnerships`
- Mentors / coaches → `HP 31 Mastery / Mentorship-Coaching`

**PARA comparison:** PARA has no people entity. People get scattered across Projects, Areas, and Resources with no consistent structure. One Mind makes humans first-class — because relationships are a pillar (LE), not an afterthought.

---

### 8. AGENT 🤖
*An AI actor with a defined role, domain scope, autonomy level, and tool access.*

```yaml
type: agent
id: oracle
emoji: 🔮
domain_scope: [HP]
autonomy: L2        # L0=human only / L1=ask first / L2=act then notify / L3=full
model: claude-sonnet-4.6
workspace: "~/.openclaw/workspace-oracle"
tools: [read_file, write_file, todoist, memory_search]
```

**Where it lives:** Agent definitions live in `SO 01 Intelligence / Agents (Digital Arms)`. Each agent has an identity file, a workspace, and a skills directory.

**PARA comparison:** PARA has no agent entity. This is the deepest structural difference — One Mind has an entire operating layer (SO) designed around AI agents working alongside humans.

---

### 9. CALENDAR 📅
*A time block, event, ritual, or recurring rhythm.*

```yaml
type: calendar
event_type: ritual  # event / block / ritual / review / deadline
recurrence: weekly  # once / daily / weekly / monthly / quarterly
pillar: HP
day: Sunday
time: "19:00"
duration: 60min
purpose: "Weekly review — CODE× × step"
linked_skill: "[[codex-maintain]]"
```

**Where it lives:** Calendar structure lives in `HP 32 Rhythm / Time-Calendar`. Events that cross pillars (family dinner = LE, board meeting = GE) are referenced from both domains. The calendar is not just scheduling — it is the rhythm that runs the CODE× loop.

**PARA comparison:** PARA implies Projects have due dates but has no calendar entity. Time is invisible in PARA. In One Mind, time is a structured entity because rhythm is how the system compounds.

---

### 10. KNOWLEDGE 📖
*An atomic unit of information — a fact, concept, reference document, lesson, or synthesized insight.*

```yaml
type: knowledge
format: concept     # concept / fact / reference / lesson / synthesis / procedure
pillar: HP
domain: "27 Mind (Education & Learning)"
source: "[[Book: Atomic Habits]]"
ai_summary: "Habits form through cue-routine-reward loops."
related: ["[[Identity-Based Habits]]", "[[30-Day-Challenges]]"]
depth: 3
```

**Where it lives:** Knowledge lives everywhere — in the domain where it is most relevant. Finance knowledge lives in HP 30. Business knowledge lives in GE. Health knowledge lives in HP 28. Cross-domain knowledge lives in SO 06 Research or HP 27 Mind.

**PARA comparison:** PARA Resources and Archives are both knowledge. One Mind separates them — there is no Archive entity. Content is either active (status: active), completed (status: done), or deprecated (status: archived). The difference is status, not location.

---

### 11. ASSET 💎
*A tangible or digital thing you own that holds value — financial, creative, or physical.*

```yaml
type: asset
category: financial  # financial / digital / physical / intellectual / real_estate
name: "Brokerage Account - Fidelity"
value: 45000
currency: USD
pillar: GE           # or LE for home/property, HP for personal
domain: "30 Finance (Money)"
tracked_in: "[[Net Worth Tracker]]"
last_valued: 2026-05-01
```

**Where it lives:**
- Financial assets (investments, accounts, crypto) → `HP 30 Finance`
- Real estate, property → `LE 50 Home`
- Business assets (IP, brand, equity) → `GE 75-83`
- Digital assets (domains, repos, accounts) → `SO 01 The Fabric`

**PARA comparison:** PARA has no asset entity. Assets are treated as Resources or Area notes. One Mind treats assets as trackable, valued, owned things — because building wealth requires knowing what you have.

---

## One Mind vs PARA vs CODE vs Jobs-to-be-Done

| Concept | PARA | CODE (Forte) | Jobs-to-be-Done | **One Mind** |
|---------|------|-------------|-----------------|--------------|
| Where things go | 4 folders | Capture → Organize → Distill → Express | User stories | 4 pillars × 11 types |
| Goals | ❌ (implied in Projects) | ❌ | ❌ | ✅ first-class entity |
| Metrics | ❌ | ❌ | ❌ | ✅ lead + lag |
| People | ❌ | ❌ | ✅ (the "job" doer) | ✅ Human entity |
| AI Agents | ❌ | ❌ | ❌ | ✅ Agent entity |
| Skills | ❌ | ❌ | ❌ | ✅ portable skill format |
| Calendar | ❌ | ❌ | ❌ | ✅ rhythm as entity |
| Assets | ❌ | ❌ | ❌ | ✅ valued, tracked |
| Execution layer | ❌ | ❌ | ❌ | ✅ CODE× method |
| Compounding | ❌ | Weakly (Distill) | ❌ | ✅ × Multiply step |

---

## The Full Map: Entities × Pillars × CODE×

```
CODE×     →    C         O         D         E         ×
               CAPTURE   ORGANIZE  DIRECT    EXECUTE   MULTIPLY

Entity ↓
Goal           create    route     prioritize  track     review
Project        create    route     scope       execute   close
Task           capture   assign    prioritize  complete  archive
Metric         define    store     set target  measure   trend
Skill          document  file      assign      run       improve
Tool           document  register  configure   use       audit
Human          add       connect   assign      contact   nurture
Agent          define    deploy    set bounds  run       retrain
Calendar       block     schedule  protect     honor     audit
Knowledge      capture   route     distill     apply     compound
Asset          record    value     protect     leverage  compound
```

---

## Frontmatter Standard for All Entities

Every One Mind file should include these fields based on its type:

```yaml
---
# REQUIRED — ALL TYPES
id: [entity-type]-[slug]-[YYYYMMDD]
type: [goal|project|task|metric|skill|tool|human|agent|calendar|knowledge|asset]
status: [active|complete|paused|archived|planned]
created: YYYY-MM-DD
updated: YYYY-MM-DD
pillar: [SO|HP|LE|GE]
domain: "XX Domain Name"
ai_summary: "One sentence — what this is."
tags: [type/[entity], zone/[pillar], domain/[code], status/[status]]

# ENTITY-SPECIFIC (add as relevant)
owner: zeus           # or agent name
agent: Legacy         # assigned agent
related: ["[[link1]]", "[[link2]]"]   # explicit graph edges
depth: 2              # 1=pillar root, 2=domain, 3=subfolder, 4=deep
---
```

---

*One Mind Entity System — documented May 2026. Part of [ONE-MIND-FRAMEWORK.md](ONE-MIND-FRAMEWORK.md).*
