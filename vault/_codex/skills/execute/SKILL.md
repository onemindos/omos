---
id: skill-execute
title: "Skill: Execute"
type: skill
status: active
created: 2026-03-27
modified: 2026-03-27
owner: legacy
tags: [type/skill, domain/codex, stage/execute]
---

# Skill: Execute

> Load context, do the work, document the output, save to the right domain, and log the action.

## CODEX Stage

**E — Execute**: The action step. This is where prioritized, assigned work actually gets done.

## Purpose

Execute turns plans into results. After items have been captured (C), organized (O), and directed (D), this skill handles the actual work — loading relevant context, performing the task, producing output, filing it correctly, and creating an audit trail.

## Triggers

- Agent receives an assigned, prioritized task from the Direct stage
- User explicitly requests an action ("do this", "create that", "update X")
- Scheduled task reaches its execution window
- Dependency resolved — blocked item is now unblocked

## Workflows

### 1. Context Loading

```
INPUT: Assigned task with priority and owner
PROCESS:
  1. Read the task description and acceptance criteria
  2. Load relevant domain files (the domain folder the task lives in)
  3. Read `_codex/system/active-goals.md` for alignment context
  4. Read `_codex/system/preferences.md` for style/format preferences
  5. Check for related recent work (avoid duplication)
OUTPUT: Fully loaded context — agent understands what to do and why
```

### 2. Work Execution

```
INPUT: Loaded context + task
PROCESS:
  1. Determine task type:
     - CREATE: Generate new content (note, document, code, plan)
     - UPDATE: Modify existing content (edit, append, revise)
     - RESEARCH: Gather and synthesize information
     - COMMUNICATE: Draft message, summary, or report
     - AUTOMATE: Build or modify a workflow/script
  2. Execute using the appropriate method for the task type
  3. Follow domain conventions and frontmatter requirements
  4. Validate output against acceptance criteria
OUTPUT: Completed work product
```

### 3. Output Filing

```
INPUT: Completed work product
PROCESS:
  1. Determine correct domain folder (by quadrant and subject)
  2. Apply proper naming convention (XX.Y.Z pattern if nested)
  3. Add required frontmatter (id, title, type, status, tags, owner)
  4. Save file to the correct location
  5. If code: save to `_system/`, add `_ops/README.md` pointer if needed
OUTPUT: Work product saved in the right place with proper metadata
```

### 4. Action Logging

```
INPUT: Completed and filed work
PROCESS:
  1. Record what was done (action summary)
  2. Record where it was saved (file path)
  3. Record which goal it served (active-goals.md reference)
  4. Record time spent (if trackable)
  5. Note any follow-up items created (feeds back to Capture)
OUTPUT: Action log entry — audit trail for the eXamine stage
```

## Task Type Reference

| Type | Action | Example |
|------|--------|---------|
| CREATE | Generate new file or content | Write a project plan, create a new SOP |
| UPDATE | Modify existing content | Update a domain README, revise a document |
| RESEARCH | Gather and synthesize | Investigate a tool, summarize findings |
| COMMUNICATE | Draft outbound content | Write a message, prepare a briefing |
| AUTOMATE | Build or modify automation | Create a script, configure a workflow |

## Rules

- Always load context before starting work — never execute blind
- Always save output to the correct domain — never leave files in Inbox after execution
- Always add frontmatter to new markdown files
- Always log what was done — the eXamine stage depends on this trail
- Code files go in `_system/`, not in knowledge domains
- If execution reveals new tasks, capture them (loop back to C) — don't scope-creep the current task
- If execution is blocked, mark the task as blocked with a reason and move to the next item
