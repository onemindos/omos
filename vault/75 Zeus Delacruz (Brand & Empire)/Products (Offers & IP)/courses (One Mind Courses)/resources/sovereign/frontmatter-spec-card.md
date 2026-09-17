# Frontmatter Spec Card

> *Quick reference for One Mind OS vault frontmatter. Every file should start with this.*

---

## Universal Frontmatter

Every file in the vault should include at minimum:

```yaml
---
id: [domain-slug-YYYYMMDD]
domain: [XX]
type: [see types below]
status: [see statuses below]
tags: []
owner: Zeus
updated: YYYY-MM-DD
---
```

---

## Domain README Frontmatter

```yaml
---
domain: XX
name: "Domain Name (Short Title)"
pillar: SO | HP | LE | GE
agent: AgentName
status: active
todoist_project_id: XXXXXXXXXXXXXXX
review: YYYY-QQ
updated: YYYY-MM-DD
---
```

---

## Goal Frontmatter

```yaml
---
id: goal-[slug]-[YYYYMMDD]
domain: XX
type: goal
pillar: SO | HP | LE | GE
status: active | completed | paused | dropped
priority: p1 | p2 | p3 | p4
target_date: YYYY-MM-DD
completion_date: null
todoist_id: XXXXXXXXXXXXXXX
tags: [#goal, #[pillar-tag]]
owner: Zeus
created: YYYY-MM-DD
updated: YYYY-MM-DD
---
```

---

## Project Frontmatter

```yaml
---
id: project-[slug]-[YYYYMMDD]
domain: XX
type: project
status: planning | active | on-hold | completed | dropped
priority: p1 | p2 | p3 | p4
start_date: YYYY-MM-DD
target_date: YYYY-MM-DD
completion_date: null
agent: AgentName
todoist_project_id: XXXXXXXXXXXXXXX
tags: [#project]
owner: Zeus
updated: YYYY-MM-DD
---
```

---

## Note / Reference Frontmatter

```yaml
---
id: note-[slug]-[YYYYMMDD]
domain: XX
type: note | reference | resource | capture
status: active | archived | draft
tags: []
source: "[URL or citation if applicable]"
owner: Zeus
created: YYYY-MM-DD
updated: YYYY-MM-DD
---
```

---

## Daily Note Frontmatter

```yaml
---
date: YYYY-MM-DD
day: Monday
type: daily
energy: 1-5
mit-1: ""
mit-2: ""
mit-3: ""
week-score: null
tags: [#daily]
---
```

---

## Weekly Review Frontmatter

```yaml
---
date: YYYY-MM-DD
week: YYYY-WW
type: review-weekly
energy-this-week: 1-5
wins: 0
completed-tasks: 0
week-score: 1-10
tags: [#review, #weekly]
---
```

---

## File Types Reference

| Type | Use For |
|------|---------|
| `goal` | Active goals — things to achieve |
| `project` | Multi-step efforts with a defined outcome |
| `note` | General knowledge, captured ideas, observations |
| `reference` | External content you're citing or keeping |
| `capture` | Raw inbox item — route before fully forming |
| `resource` | Reusable template, checklist, or reference |
| `runbook` | Step-by-step operational procedure |
| `sop` | Standard operating procedure |
| `daily` | Daily note |
| `review-weekly` | Weekly review |
| `review-monthly` | Monthly review |
| `review-quarterly` | Quarterly review |
| `index` | Folder README or index document |

---

## Status Values Reference

| Status | Meaning |
|--------|---------|
| `active` | In use, current |
| `draft` | Work in progress, not ready |
| `planning` | Not started yet but scoped |
| `on-hold` | Paused temporarily |
| `completed` | Done — move to archive section |
| `dropped` | Consciously abandoned |
| `archived` | Old but keeping for reference |

---

## Priority Tags

| Tag | Level | Meaning |
|-----|-------|---------|
| `#p1` | Critical | Do today / this week — nothing else until done |
| `#p2` | High | Do this week / this sprint |
| `#p3` | Medium | Do this month |
| `#p4` | Low | Do when time / nice to have |

---

## Todoist Integration

Tasks link to the vault via `todoist_id`. When a task is created in Todoist for a vault item:

1. Copy the Todoist task ID from the URL or API
2. Add `todoist_id: XXXXXXXXXXXXXXX` to frontmatter
3. In the task description in Todoist, add: `vault: [domain]/[filename]`

This creates a bidirectional reference — agents can find the task from the vault file and the vault file from Todoist.

---

*Part of the One Mind OS Framework — Course 1 & Course 6*
