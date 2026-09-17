---
id: ui-02-20260111053925
title: "LEGACY AI SKILLS INDEX"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "LEGACY AI SKILLS INDEX"
---

# LEGACY AI SKILLS INDEX

> **Skills are modular capabilities Legacy can execute.**
> Each skill has documentation, tools, and expected outputs.

---

## Skill Structure

Every skill follows this template:

```yaml
skill:
  name: "skill_name"
  trigger: "How user activates this"
  description: "What this skill does"

  documentation:
    - path: "path/to/doc.md"
      purpose: "Why Legacy reads this"

  tools:
    - name: "tool_name"
      action: "What to do with it"

  memory:
    reads: ["memory_blocks to read"]
    updates: ["memory_blocks to update"]

  output:
    type: "response | task | document | notification"
    format: "What the output looks like"

  hitl: "required | optional | none"
```

---

## Core Skills

### 1. memory_sync

```yaml
skill:
  name: "memory_sync"
  trigger: "Legacy, run memory sync" | "sync your memory" | scheduled cron
  description: "Sync memory blocks from knowledge base"

  documentation:
    - path: "02 Agents/02.1 Mother AI/02.1.1 Index/_LEGACY-SYNC.md"
      purpose: "Get list of pending items to sync"
    - path: "[varies based on pending items]"
      purpose: "Read referenced documents"

  tools:
    - name: "github_read_file"
      action: "Read sync manifest and referenced docs"
    - name: "github_edit_file"
      action: "Mark items as synced"

  memory:
    reads: ["sync_status"]
    updates: ["persona", "vault_map", "active_context", "sync_status"]

  output:
    type: "response"
    format: "Sync complete. Updated X memory blocks: [list]"

  hitl: "none"
```

---

### 2. navigate_vault

```yaml
skill:
  name: "navigate_vault"
  trigger: "Find [topic]" | "Where is [thing]" | "Show me [area]"
  description: "Find documents and areas in the knowledge base"

  documentation:
    - path: "01 Command/.../01.2 Atlas/.../_Atlas.md"
      purpose: "Vault map and navigation reference"

  tools:
    - name: "github_read_file"
      action: "Read Atlas for structure"
    - name: "github_search"
      action: "Search for specific content"

  memory:
    reads: ["vault_map"]
    updates: []

  output:
    type: "response"
    format: "Found [X]. Path: [path]. Todoist filter: [filter]"

  hitl: "none"
```

---

### 3. create_task

```yaml
skill:
  name: "create_task"
  trigger: "Add task [description]" | "Schedule [thing]" | "Remind me to [action]"
  description: "Create a task in Todoist with proper metadata"

  documentation:
    - path: "01 Command/.../TODOIST-STRUCTURE-GUIDE.md"
      purpose: "Know correct project, labels, tags"
    - path: "01 Command/.../01.2 Atlas/.../_Atlas.md"
      purpose: "Map task to correct pillar/area"

  tools:
    - name: "todoist_create_task"
      action: "Create the task with metadata"

  memory:
    reads: ["vault_map", "user_context", "active_context"]
    updates: ["active_context"]

  output:
    type: "task"
    format: "Created: [task] in [project] due [date]"

  hitl: "optional"  # Based on trust level
```

---

### 4. daily_brief

```yaml
skill:
  name: "daily_brief"
  trigger: "Morning brief" | "What's on today" | scheduled 6 AM
  description: "Generate daily overview and priorities"

  documentation:
    - path: "01 Command/.../01.2 Atlas/.../_Atlas.md"
      purpose: "Know all areas for comprehensive check"

  tools:
    - name: "todoist_get_tasks"
      action: "Get today's tasks, overdue, priorities"
    - name: "github_read_file"
      action: "Check any relevant dated docs"

  memory:
    reads: ["user_context", "active_context", "vault_map"]
    updates: ["active_context"]

  output:
    type: "response"
    format: |
      ## Morning Brief - [date]
      ### Priority Tasks (X items)
      ### Overdue (X items)
      ### Calendar
      ### Focus Recommendation

  hitl: "none"
```

---

### 5. check_pillar

```yaml
skill:
  name: "check_pillar"
  trigger: "Check HP" | "How's GE doing" | "Status of [pillar]"
  description: "Review status of a specific life pillar"

  documentation:
    - path: "[pillar]/.../[area.1] Index/_Index.md"
      purpose: "Get pillar overview and status"
    - path: "01 Command/.../01.2 Atlas/.../_Atlas.md"
      purpose: "Know pillar structure"

  tools:
    - name: "todoist_get_tasks"
      action: "Get tasks for that pillar (filter: ##[pillar])"
    - name: "github_read_file"
      action: "Read relevant index docs"

  memory:
    reads: ["vault_map", "[pillar]_context"]
    updates: ["[pillar]_context"]

  output:
    type: "response"
    format: |
      ## [Pillar] Status
      ### Open Tasks: X
      ### Overdue: X
      ### Key Areas: [list]
      ### Recommendations: [list]

  hitl: "none"
```

---

### 6. process_inbox

```yaml
skill:
  name: "process_inbox"
  trigger: "Process inbox" | "What's in my inbox" | "Triage tasks"
  description: "Review and categorize inbox items"

  documentation:
    - path: "01 Command/.../TODOIST-STRUCTURE-GUIDE.md"
      purpose: "Know where to route tasks"
    - path: "01 Command/.../01.2 Atlas/.../_Atlas.md"
      purpose: "Map items to correct areas"

  tools:
    - name: "todoist_get_tasks"
      action: "Get inbox tasks (filter: #Inbox)"
    - name: "todoist_update_task"
      action: "Move to correct project, add labels"

  memory:
    reads: ["vault_map", "user_context"]
    updates: ["active_context"]

  output:
    type: "response + tasks"
    format: "Processed X inbox items. Moved Y to projects. Z need your input."

  hitl: "required"  # User confirms categorization
```

---

### 7. weekly_review

```yaml
skill:
  name: "weekly_review"
  trigger: "Weekly review" | scheduled Sunday 9 AM
  description: "Comprehensive weekly status and planning"

  documentation:
    - path: "01 Command/.../01.2 Atlas/.../_Atlas.md"
      purpose: "Review all pillars"
    - path: "25-49 HP/.../29 Life Systems/.../_Index.md"
      purpose: "Check routine/habit tracking"

  tools:
    - name: "todoist_get_tasks"
      action: "Get completed this week, upcoming, overdue"
    - name: "github_read_file"
      action: "Read relevant tracking docs"

  memory:
    reads: ["all memory blocks"]
    updates: ["active_context", "sync_status"]

  output:
    type: "document"
    format: |
      ## Weekly Review - Week of [date]
      ### Completed: X tasks
      ### By Pillar: UI(X) HP(X) LE(X) GE(X)
      ### Wins:
      ### Challenges:
      ### Next Week Focus:

  hitl: "optional"
```

---

### 8. document_lookup

```yaml
skill:
  name: "document_lookup"
  trigger: "Read [doc name]" | "What does [doc] say" | "Get context on [topic]"
  description: "Find and read specific documentation"

  documentation:
    - path: "01 Command/.../01.2 Atlas/.../_Atlas.md"
      purpose: "Locate document path"

  tools:
    - name: "github_read_file"
      action: "Read the requested document"

  memory:
    reads: ["vault_map"]
    updates: []  # Unless explicitly asked to remember

  output:
    type: "response"
    format: "Here's what [doc] contains: [summary or full content]"

  hitl: "none"
```

---

### 9. self_schedule

```yaml
skill:
  name: "self_schedule"
  trigger: Legacy initiates based on protocols
  description: "Legacy creates tasks for herself"

  documentation:
    - path: "LEGACY-AI-SYSTEM.md"
      purpose: "Know scheduled behaviors"
    - path: "01 Command/.../TODOIST-STRUCTURE-GUIDE.md"
      purpose: "Proper task formatting"

  tools:
    - name: "todoist_create_task"
      action: "Create task with @legacy label"

  memory:
    reads: ["persona", "active_context"]
    updates: ["active_context"]

  output:
    type: "task"
    format: "Scheduled: [task] for [time] with @legacy label"

  hitl: "none"  # This is how Legacy manages herself
```

---

### 10. context_update

```yaml
skill:
  name: "context_update"
  trigger: Conversation reveals important user info
  description: "Update user context from conversation"

  documentation:
    - path: "None - conversation based"
      purpose: "Learn from what user shares"

  tools:
    - name: "memory_update"
      action: "Update user_context block"

  memory:
    reads: ["user_context"]
    updates: ["user_context"]

  output:
    type: "silent"  # No explicit output unless asked
    format: "[Internal: Updated user preference for X]"

  hitl: "none"
```

---

## How to Add New Skills

1. Create new file: `02.../02.1.../02.1.2 Skills/skill_[name].md`
2. Use the template above
3. Add to this index
4. Add sync entry so Legacy learns it

### Skill File Template

```markdown
# Skill: [name]

## Trigger
[How to activate]

## Description
[What it does]

## Documentation
| Path | Purpose |
|------|---------|
| path/to/doc.md | Why Legacy reads this |

## Tools
| Tool | Action |
|------|--------|
| tool_name | What to do |

## Memory
- **Reads:** [blocks]
- **Updates:** [blocks]

## Output
- **Type:** [response/task/document/notification]
- **Format:** [template]

## HITL
[required/optional/none]

## Examples
[Example interactions]
```

---

## Skill Categories

| Category | Skills | Purpose |
|----------|--------|---------|
| **Memory** | memory_sync, context_update | Keep Legacy current |
| **Navigation** | navigate_vault, document_lookup | Find information |
| **Tasks** | create_task, process_inbox, self_schedule | Task management |
| **Review** | daily_brief, weekly_review, check_pillar | Status & planning |

---

## Skills Sync Entry

Add this to `_LEGACY-SYNC.md` for Legacy to learn the skills system:

```yaml
- [pending] | P1 | system | 02 Agents/.../02.1.2 Skills/_SKILLS-INDEX.md | Skills framework created
```

---

*Skills make Legacy modular, teachable, and extensible.*
*Each skill is self-documenting with clear inputs/outputs.*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
