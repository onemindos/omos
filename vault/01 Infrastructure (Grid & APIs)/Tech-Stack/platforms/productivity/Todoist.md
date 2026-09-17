# Todoist

> **Task Management**

## Overview

Todoist is the primary task management system for OneMind OS, handling personal tasks, projects, and daily planning.

## Status

- **Active**: ✅ Yes
- **Plan**: Pro
- **Last Updated**: January 25, 2026

## Use Cases

1. **Task Management** - Daily to-dos
2. **Project Tracking** - Multi-step projects
3. **Recurring Tasks** - Habits and routines
4. **Quick Capture** - Fast task entry
5. **Planning** - Weekly/daily planning

## Key Features

- Natural language input
- Projects and sections
- Labels and filters
- Priority levels
- Recurring due dates
- Comments and attachments
- Integrations

## Structure

### Projects (Pillar-Based)

| Project                                   | Purpose             |
| ----------------------------------------- | ------------------- |
| **📥 Inbox**                              | Quick capture       |
| **🧠 UI (Sovereign Operations)**          | AI, tech, systems   |
| **💪 HP (Holistic Performance)**          | Health, learning    |
| **⚡ LE (Legacy Evolution)**              | Business, career    |
| **🌱 GE (Generational Endeavor)** | Family, legacy      |
| **📋 Routines**                           | Daily/weekly habits |

### Labels

| Label       | Purpose              |
| ----------- | -------------------- |
| `@home`     | Location-based       |
| `@computer` | Needs computer       |
| `@phone`    | Quick phone tasks    |
| `@errand`   | Outside tasks        |
| `@waiting`  | Waiting for response |
| `@legacy`   | Assign to Legacy AI  |

## API Integration

### Configuration

```python
from todoist_api_python import TodoistAPI
import os

api = TodoistAPI(os.getenv("TODOIST_API_KEY"))
```

### Get Tasks

```python
def get_tasks(filter: str = None, project_id: str = None):
    try:
        if filter:
            return api.get_tasks(filter=filter)
        elif project_id:
            return api.get_tasks(project_id=project_id)
        else:
            return api.get_tasks()
    except Exception as e:
        print(f"Error: {e}")
        return []

# Get today's tasks
today_tasks = get_tasks(filter="today")

# Get overdue
overdue = get_tasks(filter="overdue")
```

### Create Task

```python
def create_task(
    content: str,
    project_id: str = None,
    due_string: str = None,
    priority: int = 1,
    labels: list = None
):
    try:
        return api.add_task(
            content=content,
            project_id=project_id,
            due_string=due_string,
            priority=priority,
            labels=labels or []
        )
    except Exception as e:
        print(f"Error: {e}")
        return None

# Examples
create_task("Review PR #234", due_string="today", priority=2)
create_task("Call dentist", due_string="tomorrow at 10am", labels=["@phone"])
create_task("Weekly review", due_string="every friday at 5pm")
```

### Complete Task

```python
def complete_task(task_id: str):
    try:
        return api.close_task(task_id)
    except Exception as e:
        print(f"Error: {e}")
        return False
```

### Update Task

```python
def update_task(task_id: str, **updates):
    try:
        return api.update_task(task_id, **updates)
    except Exception as e:
        print(f"Error: {e}")
        return None

# Examples
update_task("123", content="Updated task name")
update_task("123", due_string="tomorrow")
update_task("123", priority=4)  # P1 (highest)
```

### Get Projects

```python
def get_projects():
    try:
        return api.get_projects()
    except Exception as e:
        print(f"Error: {e}")
        return []

def get_project_by_name(name: str):
    projects = get_projects()
    return next((p for p in projects if p.name == name), None)
```

## Natural Language

### Due Dates

| Input           | Result          |
| --------------- | --------------- |
| `today`         | Today           |
| `tomorrow`      | Tomorrow        |
| `next friday`   | Next Friday     |
| `in 3 days`     | 3 days from now |
| `jan 25`        | January 25      |
| `every day`     | Daily recurring |
| `every weekday` | Mon-Fri         |
| `every 2 weeks` | Bi-weekly       |

### Priority

| Priority | Level        | Flag |
| -------- | ------------ | ---- |
| 4        | P1 (highest) | 🔴   |
| 3        | P2           | 🟠   |
| 2        | P3           | 🔵   |
| 1        | P4 (lowest)  | ⚪   |

## Filters

### Common Filters

| Filter            | Purpose           |
| ----------------- | ----------------- |
| `today`           | Due today         |
| `tomorrow`        | Due tomorrow      |
| `overdue`         | Past due          |
| `7 days`          | Next 7 days       |
| `no date`         | No due date       |
| `@waiting`        | Waiting for label |
| `#UI`             | Project "UI"      |
| `p1`              | Priority 1        |
| `assigned to: me` | My tasks          |

### Custom Filters

```
# Focus for today
(today | overdue) & !@waiting

# Deep work
(today | overdue) & @computer & (p1 | p2)

# Quick tasks
@phone | @errand

# Weekly review
7 days & !subtask
```

## Legacy AI Integration

### Agent Tools

```python
@Tool
def get_todays_tasks() -> list:
    """Get all tasks due today."""
    tasks = get_tasks(filter="today | overdue")
    return [
        {
            "id": t.id,
            "content": t.content,
            "priority": t.priority,
            "due": t.due.string if t.due else None,
            "project": get_project_name(t.project_id)
        }
        for t in tasks
    ]

@Tool
def add_task(content: str, due: str = None, priority: int = 1, project: str = None) -> str:
    """Add a new task to Todoist."""
    project_id = None
    if project:
        proj = get_project_by_name(project)
        project_id = proj.id if proj else None

    task = create_task(
        content=content,
        due_string=due,
        priority=priority,
        project_id=project_id
    )

    return f"Created task: {task.content}" if task else "Failed to create task"

@Tool
def complete_todo(task_id: str) -> str:
    """Mark a task as complete."""
    if complete_task(task_id):
        return "Task completed!"
    return "Failed to complete task"
```

### Voice Commands

| Command                               | Action              |
| ------------------------------------- | ------------------- |
| "Add task: Buy groceries"             | Creates task        |
| "What's on my list today?"            | Lists today's tasks |
| "Complete the first task"             | Marks complete      |
| "Reschedule meeting prep to tomorrow" | Updates due date    |

## Webhooks

### Sync Events

Todoist doesn't have native webhooks, but you can:

1. **Poll API** - Check for changes periodically
2. **Use IFTTT** - Trigger on new tasks
3. **Zapier** - Connect to other services
4. **n8n** - Self-hosted automation

### Polling Example

```python
import time

last_sync = None

async def sync_tasks():
    global last_sync
    while True:
        tasks = get_tasks()

        # Check for changes since last sync
        if last_sync:
            new_tasks = [t for t in tasks if t.created_at > last_sync]
            for task in new_tasks:
                await handle_new_task(task)

        last_sync = datetime.utcnow()
        await asyncio.sleep(60)  # Poll every minute
```

## Integrations

### Calendar Sync

- Google Calendar (2-way sync)
- Tasks appear on calendar
- Events create tasks

### Other Integrations

| Service    | Integration                |
| ---------- | -------------------------- |
| **Slack**  | Create tasks from messages |
| **Email**  | Forward emails to Todoist  |
| **Alexa**  | Voice task creation        |
| **Zapier** | Automation                 |
| **IFTTT**  | Simple automations         |

## Best Practices

- ✅ Review Inbox daily (GTD capture)
- ✅ Use natural language for dates
- ✅ Set realistic due dates
- ✅ Use labels for context
- ✅ Weekly review of all projects
- ✅ Break large tasks into subtasks
- ✅ Use recurring tasks for habits
- ✅ Keep task names actionable (verbs)

## Keyboard Shortcuts

| Shortcut     | Action        |
| ------------ | ------------- |
| `q`          | Quick add     |
| `a`          | Add task      |
| `e`          | Edit task     |
| `c`          | Complete task |
| `/`          | Search        |
| `g` then `t` | Go to Today   |
| `g` then `i` | Go to Inbox   |

## Resources

- [Todoist Help](https://todoist.com/help)
- [API Documentation](https://developer.todoist.com/)
- [Python SDK](https://github.com/Doist/todoist-api-python)

---

**Related**: [Obsidian](./Obsidian.md) | [N8N](./N8N.md) | [Agno](./Agno.md)
