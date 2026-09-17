---
name: todoist
description: Manage Todoist tasks, projects, sections, and labels for One Mind infrastructure. Use when you need to create tasks, query tasks, manage projects, set up pillar structure, or organize with labels. Makes API calls to Todoist.
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Read
  - Write
---

# Todoist Integration Skill

Interact with Todoist REST API for One Mind task and project management.

## Environment Variables Required

```bash
TODOIST_API_TOKEN   # Bearer token from todoist.com/app/settings/integrations/developer
```

## Quick Reference

| Operation | Command |
|-----------|---------|
| List projects | `curl -s "https://api.todoist.com/rest/v2/projects" -H "Authorization: Bearer $TODOIST_API_TOKEN"` |
| Get tasks | `curl -s "https://api.todoist.com/rest/v2/tasks" -H "Authorization: Bearer $TODOIST_API_TOKEN"` |
| Filter tasks | `curl -s "https://api.todoist.com/rest/v2/tasks?filter=today" -H "Authorization: Bearer $TODOIST_API_TOKEN"` |

## One Mind Structure

### Pillar Projects
```
00-24 SO (Sovereign Operations)  # The Brain
25-49 HP (Holistic Performance)  # Self Path
50-74 LE (Legacy Evolution)      # Legacy Path
75-99 GE (Generational Empire)  # Wealth Path
```

### Standard Sections (per project)
```
Backlog
This Sprint
In Progress
Blocked
Done (This Week)
```

### Key Labels
```
# Agent labels
@mother-ai, @hp-conductor, @le-conductor, @ge-conductor

# Status labels
@approval-needed, @blocked, @urgent, @waiting

# Context labels
@computer, @phone, @home, @office

# Type labels
@project, @routine, @review, @research
```

## Operations

### Tasks

**List all tasks:**
```bash
curl -s "https://api.todoist.com/rest/v2/tasks" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" | jq '.[] | {id, content, priority, due: .due.date}'
```

**Filter tasks (Todoist filter syntax):**
```bash
# Today's tasks
curl -s "https://api.todoist.com/rest/v2/tasks?filter=today" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"

# HP pillar tasks
curl -s "https://api.todoist.com/rest/v2/tasks?filter=%2325-49%20HP" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"

# Urgent tasks
curl -s "https://api.todoist.com/rest/v2/tasks?filter=p1" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"

# Approval needed
curl -s "https://api.todoist.com/rest/v2/tasks?filter=%40approval-needed" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"
```

**Create task:**
```bash
curl -s -X POST "https://api.todoist.com/rest/v2/tasks" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Task title",
    "project_id": "PROJECT_ID",
    "section_id": "SECTION_ID",
    "priority": 2,
    "due_string": "tomorrow",
    "labels": ["hp-conductor", "urgent"],
    "description": "Task details and context"
  }'
```

**Update task:**
```bash
curl -s -X POST "https://api.todoist.com/rest/v2/tasks/{task_id}" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Updated title",
    "priority": 3,
    "due_string": "next monday"
  }'
```

**Complete task:**
```bash
curl -s -X POST "https://api.todoist.com/rest/v2/tasks/{task_id}/close" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"
```

**Delete task:**
```bash
curl -s -X DELETE "https://api.todoist.com/rest/v2/tasks/{task_id}" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"
```

### Projects

**List all projects:**
```bash
curl -s "https://api.todoist.com/rest/v2/projects" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" | jq '.[] | {id, name, parent_id}'
```

**Create project:**
```bash
curl -s -X POST "https://api.todoist.com/rest/v2/projects" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Project Name",
    "parent_id": "PARENT_PROJECT_ID",
    "color": "blue"
  }'
```

**Get project by ID:**
```bash
curl -s "https://api.todoist.com/rest/v2/projects/{project_id}" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN"
```

### Sections

**List sections for project:**
```bash
curl -s "https://api.todoist.com/rest/v2/sections?project_id={project_id}" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" | jq '.[] | {id, name, order}'
```

**Create section:**
```bash
curl -s -X POST "https://api.todoist.com/rest/v2/sections" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Section Name",
    "project_id": "PROJECT_ID",
    "order": 1
  }'
```

**Create standard sections for a project:**
```bash
# Run for each section: Backlog, This Sprint, In Progress, Blocked, Done (This Week)
for section in "Backlog" "This Sprint" "In Progress" "Blocked" "Done (This Week)"; do
  curl -s -X POST "https://api.todoist.com/rest/v2/sections" \
    -H "Authorization: Bearer $TODOIST_API_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"$section\", \"project_id\": \"PROJECT_ID\"}"
done
```

### Labels

**List all labels:**
```bash
curl -s "https://api.todoist.com/rest/v2/labels" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" | jq '.[] | {id, name, color}'
```

**Create label:**
```bash
curl -s -X POST "https://api.todoist.com/rest/v2/labels" \
  -H "Authorization: Bearer $TODOIST_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "label-name",
    "color": "yellow"
  }'
```

## Filter Syntax

| Filter | Meaning |
|--------|---------|
| `today` | Due today |
| `tomorrow` | Due tomorrow |
| `overdue` | Past due |
| `p1` | Priority 1 (urgent) |
| `p2` | Priority 2 (high) |
| `@label-name` | Has label |
| `#project-name` | In project |
| `##25-49 HP` | In pillar project |
| `no due date` | No due date set |
| `assigned to: me` | Assigned to you |

**Combining filters:**
- `today & @urgent` - Today AND urgent
- `overdue | p1` - Overdue OR priority 1
- `##25-49 HP & @approval-needed` - HP pillar needing approval

## Priority Levels

| Value | Meaning | Color |
|-------|---------|-------|
| 1 | Normal | None |
| 2 | Medium | Blue |
| 3 | High | Orange |
| 4 | Urgent | Red |

## Error Handling

- **200**: Success (GET)
- **204**: Success (POST/DELETE, no content)
- **400**: Bad request (check JSON)
- **401**: Invalid token
- **404**: Resource not found
- **429**: Rate limited (max 450 requests/15 min)

## Tips

- URL-encode filter queries (`@` = `%40`, `#` = `%23`, space = `%20`)
- Use `jq` to parse and filter JSON responses
- Priority 4 is highest (urgent), 1 is lowest (normal)
- Labels don't include `@` when creating, only in filters
- `due_string` accepts natural language: "today", "tomorrow 3pm", "every monday"


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
