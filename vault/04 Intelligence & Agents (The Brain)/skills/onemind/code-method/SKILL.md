---
name: code-method
description: "Run the full CODE× operating loop — Capture → Organize → Direct → Execute → × Multiply. Use when processing inbox captures, routing content to the correct domain, prioritizing tasks against active goals, executing organized work, or running daily/weekly/monthly reflect cycles."
id: skill-code-method
title: "Skill: CODE× Method"
type: skill
status: active
created: 2026-03-27
modified: 2026-05-04
owner: legacy
tags: [type/skill, domain/onemind, stage/all]
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Bash(python3:*)
  - Read
  - Write
  - WebFetch
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# CODE× Method — Master Skill

> Capture → Organize → Direct → Execute → × Multiply. The complete One Mind operating loop.

For the full domain routing map see `references/domain-routing.md`. For Organize workflows and folder conventions see `references/organize-workflows.md`.

---

## C — CAPTURE

> Ingest raw input from any source and save it as a structured entry for processing.

### When to Activate Capture

| Signal | Examples |
|--------|---------|
| URL shared | Any `https://` link in message |
| Save intent | "save this", "capture this", "file this", "remember this" |
| Content forwarded | Forwarded message from another chat |
| Media shared | Image, screenshot, PDF, voice note |
| Recipe link | Any food/cooking URL |
| Product link | Amazon, store URLs, "want to buy" |
| Content idea | "remake this", "repurpose", "content like this" |
| Research item | Article, paper, video, podcast link |

### Capture Architecture

```
INPUT (URL, text, image, voice)
  ├── 1. DETECT type
  ├── 2. EXTRACT content (scrape/OCR/transcribe)
  ├── 3. CLASSIFY domain (SO/HP/LE/GE + subdomain)
  ├── 4. EXTRACT structured value
  ├── 5. DRAFT codex note from template
  ├── 6. PROPOSE routing + optional action
  └── 7. EXECUTE on approval (write vault + create task)
```

### Model Strategy

Run as a **subagent on the cheapest model** (gpt-4o-mini, Gemini Flash, Claude Haiku). Legacy detects capture intent and spawns a subagent with this skill. **Cost per capture: ~$0.001-0.005**

### Step 1-3: Detect, Extract, Classify

Input types: `url | text | image | voice | file | forwarded`

For URLs, scrape using curl + python. Extract: title, meta description, og:type. Auto-high patterns:
```
food blog domains → HP 29, recipe
amazon.com → product (ask domain)
instagram.com/reel → content-idea, GE
youtube.com/watch → reference (classify topic)
github.com → SO, reference
```

Confidence levels: **high** → auto-route, **medium** → ask Zeus to pick from 2-3 options, **low** → ask Zeus directly.

### Step 4: Present for Approval

```
📁 Route: [Domain + Subfolder]
🔗 Source: [Platform — Author]
Summary: [1-2 sentence summary]
📋 Suggested task: [optional Todoist task]
✅ Approve  ✏️ Re-route  ❌ Discard
```

### Step 5: Execute on Approval

- Write note to `[domain]/[subfolder]/[slug].md` with proper frontmatter
- Create Todoist task if suggested and approved
- Log capture in `_codex/memory/learning/capture-patterns.md`

### Batch Mode

Accept array of URLs → process each through steps 1-5 → group by domain for efficient approval → execute approved items in batch → report: `"{n} captured, {m} tasks created, {k} discarded"`

---

## O — ORGANIZE

> Route captured items to the correct domain + subfolder, maintain folder structure, update READMEs.

See `references/domain-routing.md` for the full domain routing map and `references/organize-workflows.md` for folder conventions, loose file protocol, and duplicate folder protocol.

**Core rule:** Read `CLAUDE.md`, `_codex/CONVENTIONS.md`, and the domain's `_README.md` BEFORE touching anything. Never delete without Zeus's explicit approval.

---

## D — DIRECT

> Prioritize organized items against active goals, assign to the right agent, set execution order.

### Triggers

- New items organized but lack priority or assignment
- Active goals have changed and existing items need re-ranking
- Agent asks "what should I do next?"
- Weekly or daily planning sessions

### Workflow 1: Priority Assessment

```
INPUT: Unranked items in any domain
1. Read _codex/system/active-goals.md
2. Score each item against active goals (high/medium/low alignment)
3. Factor in urgency (deadline) and effort (quick win vs. deep work)
4. Assign priority: P1 (do now), P2 (do soon), P3 (do later), P4 (delegate/defer)
OUTPUT: Items tagged with priority level
```

### Workflow 2: Agent Assignment

```
SO (00-24) → Legacy (orchestrator) or Grid/Guardian (specialist)
HP (25-49) → Zenith (body + mind + money)
LE (50-74) → Unity (home/family)
GE (75-99) → Eternal (ventures/brand/wealth)
Cross-pillar: Legacy coordinates, specialist executes
```

### Workflow 3: Execution Queue

Cap each agent's active queue at 3-5 items. Order by priority. Identify dependencies. Set execution windows (morning deep work, afternoon admin).

### Decision Framework

| Factor | Weight |
|--------|--------|
| Goal alignment | High |
| Urgency | High |
| Effort | Medium |
| Dependencies | Medium |
| Pillar balance | Low |

### Direct Rules

- Never queue more than 5 P1 items per agent at once
- Always check `active-goals.md` before assigning priority
- Items with no goal alignment drop to P3 or P4 automatically
- Re-run priority assessment when active goals change

---

## E — EXECUTE

> Load context, do the work, document the output, save to the right domain, and log the action.

### Workflow 1: Context Loading

```
1. Read task description and acceptance criteria
2. Load relevant domain files
3. Read _codex/system/active-goals.md
4. Read _codex/system/preferences.md
5. Check for related recent work
```

### Workflow 2: Work Execution

Task types: `CREATE | UPDATE | RESEARCH | COMMUNICATE | AUTOMATE`

Follow domain conventions and frontmatter requirements. Validate output against acceptance criteria.

### Workflow 3: Output Filing

1. Determine correct domain folder
2. Apply proper naming convention
3. Add required frontmatter (id, title, type, status, tags, owner)
4. Save file to the correct location (code → `_system/` only)

### Workflow 4: Action Logging

Record: what was done, where it was saved, which goal it served. The × stage depends on this trail.

### Execute Rules

- Always load context before starting — never execute blind
- Always save output to correct domain — never leave files in Inbox after execution
- Always add frontmatter to new markdown files
- If execution reveals new tasks, capture them (loop back to C)
- If blocked, mark as blocked and move to next item

---

## × — REFLECT & MULTIPLY

> Review completed work, extract lessons, score performance, and feed insights back into the system.

### Reflection Cadences

**Daily (5 min):** Review today's activity → log to `_codex/templates/daily.md` → extract facts as permanent knowledge → flag unfinished items

**Weekly (15 min):** Read last 7 daily logs → consolidate patterns → check goals in `active-goals.md` → process inbox → plan next week's top 3 priorities

**Monthly (30 min):** Read last 4 weekly reviews → audit each pillar (SO infra, HP health, LE family, GE business) → run consolidate + maintain skills → update goals

### Sleep-Time Reflection (Agent-Specific)

When agent has downtime: scan recent memory for promotable facts → check for contradictions → extract preferences ("Zeus said he prefers X") → flag knowledge gaps with `status/gap`

### Three Types of Reflection Output

| Output | Where It Goes |
|--------|--------------|
| **Insight** (new knowledge) | Codex domain folder |
| **Update** (corrected knowledge) | Edit existing note |
| **Action** (something to do) | Task system or inbox |

### The Multiply Effect

Every insight fed back into the system makes the next cycle faster, smarter, and more accurate. The loop compounds:

```
Better Capture → Better Organization → Better Direction → Better Execution → Better Reflection → ...
```

This is what separates CODE× from linear methods. The system learns.

---

## Related

- `references/domain-routing.md` — full SO/HP/LE/GE domain routing map
- `references/organize-workflows.md` — folder conventions, loose file protocol, duplicate protocol
- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
