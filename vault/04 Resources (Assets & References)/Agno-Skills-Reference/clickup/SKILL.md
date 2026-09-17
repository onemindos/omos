---
name: clickup
description: Manage tasks, projects, and docs in ClickUp
license: Apache-2.0
metadata:
  version: "1.0.0"
  author: zeus
  tags: ["productivity", "tasks", "project-management", "docs", "clickup"]
allowed-tools:
  - clickup
  - clickup_list_workspaces
  - clickup_list_spaces
  - clickup_list_folders
  - clickup_list_lists
  - clickup_list_tasks
  - clickup_get_task
  - clickup_create_task
  - clickup_update_task
  - clickup_delete_task
  - clickup_get_tasks_due
  - clickup_get_task_with_docs
  - clickup_add_comment
  - clickup_search_tasks
  - clickup_list_docs
  - clickup_search_docs
  - clickup_get_doc
  - clickup_get_doc_pages
  - clickup_create_doc
  - clickup_create_doc_page
  - clickup_update_doc_page
---

# ClickUp Management Skill

You are an expert at managing tasks, projects, and documentation in ClickUp.

## When to Use

- User wants to manage their tasks or projects
- User needs to check what's due today/this week
- User wants to create, update, or organize tasks
- User mentions "ClickUp" or asks about task status
- User wants to read or create documentation
- User has tasks with linked docs that contain instructions

## ClickUp Structure

```
Workspace (Team)
└── Spaces (high-level areas)
    └── Folders (optional grouping)
        └── Lists (task containers)
            └── Tasks (work items)
                └── Subtasks
                └── Comments
                └── Checklists
                └── Time entries
```

## Available Tools

### Organization Tools
| Tool | Purpose |
|------|---------|
| `clickup_list_workspaces` | List all workspaces you can access |
| `clickup_list_spaces` | List spaces in the workspace |
| `clickup_list_folders` | List folders in a space |
| `clickup_list_lists` | List lists in a folder/space |

### Task Management
| Tool | Purpose |
|------|---------|
| `clickup_list_tasks` | List tasks in a list (filter by status, assignees) |
| `clickup_get_task` | Get full task details with linked doc info |
| `clickup_create_task` | Create a new task |
| `clickup_update_task` | Update task name, status, priority, due date |
| `clickup_delete_task` | Delete a task |
| `clickup_search_tasks` | Search tasks across workspace |

### Task + Doc Integration (Key Feature)
| Tool | Purpose |
|------|---------|
| `clickup_get_tasks_due` | Get tasks due today/soon with linked doc info |
| `clickup_get_task_with_docs` | Get task AND read any linked doc content |

### Comments & Time
| Tool | Purpose |
|------|---------|
| `clickup_add_comment` | Add comment to a task |
| `clickup_log_time` | Log time to a task |

### Documentation
| Tool | Purpose |
|------|---------|
| `clickup_list_docs` | List all docs in workspace |
| `clickup_search_docs` | Search docs by query |
| `clickup_get_doc` | Get doc details |
| `clickup_get_doc_pages` | Get all pages in a doc |
| `clickup_create_doc` | Create a new doc |
| `clickup_create_doc_page` | Create a page in a doc |
| `clickup_update_doc_page` | Update page content |

## Common Workflows

### 1. Check What's Due Today

```python
# Get tasks due today including overdue
clickup_get_tasks_due(days=1, include_overdue=True)
```

This returns tasks with their linked doc IDs if any.

### 2. Work on a Task with Instructions in a Doc

```python
# Get task + automatically fetch linked doc content
clickup_get_task_with_docs(task_id="abc123")
```

Returns:
```json
{
  "name": "Build feature X",
  "status": "in progress",
  "linked_docs": [
    {
      "doc_id": "xyz",
      "page_id": "123",
      "pages": [
        {"page_name": "Requirements", "content": "...markdown..."}
      ]
    }
  ]
}
```

### 3. Create a Task

```python
clickup_create_task(
    list_id="901234567",
    name="New feature request",
    description="Build the widget",
    priority=2,  # 1=Urgent, 2=High, 3=Normal, 4=Low
    due_date="2025-02-05",
    status="Open"
)
```

### 4. Search and Update

```python
# Find task
result = clickup_search_tasks(query="authentication bug")

# Update it
clickup_update_task(
    task_id="abc123",
    status="in progress",
    priority=1
)
```

### 5. Document Something

```python
# Create a doc
doc = clickup_create_doc(
    name="Project Requirements",
    content="# Overview\nThis project..."
)

# Add a sub-page
clickup_create_doc_page(
    doc_id=doc["id"],
    name="Technical Specs",
    content="## Architecture..."
)
```

## Priority Mapping

| Value | Level | Use When |
|-------|-------|----------|
| 1 | Urgent | Blocking others, critical bug |
| 2 | High | Important, do this week |
| 3 | Normal | Standard priority |
| 4 | Low | Nice to have |

## Task Statuses

Statuses vary by list, but common ones:
- `Open` / `to do` - Not started
- `in progress` - Being worked on
- `review` - Awaiting review
- `complete` / `Closed` - Done

Always use the exact status name from the list (case-sensitive).

## Tips

1. **Always check due tasks first** - Start with `clickup_get_tasks_due(days=1)`
2. **Read linked docs** - If a task has a linked doc, use `clickup_get_task_with_docs()` to get instructions
3. **Search before creating** - Use `clickup_search_tasks()` to avoid duplicates
4. **Add comments for updates** - When completing work, add a comment with what was done
5. **Log time** - Help track effort with `clickup_log_time()`

## Output Format

When presenting tasks to user:

```
## Due Today (3 tasks)

1. **[HIGH]** Task name - List name
   Due: Today at 5pm
   Status: in progress
   📄 Has linked doc with instructions

2. **[NORMAL]** Another task - List name
   Due: Today at EOD
   Status: Open

## Overdue (1 task)

1. **[URGENT]** Overdue task - List name
   Due: Yesterday
   Status: Open
```

## Error Handling

If you get an error:
- `CLICKUP_API_KEY not set` - Configuration issue, inform user
- `CLICKUP_WORKSPACE_ID not set` - Configuration issue, inform user
- 404 errors - Resource not found, verify IDs
- Rate limit - Wait and retry (tool handles this automatically)
