---
id: ui-02-20260111053925
title: "Letta Agent ↔ Todoist Integration"
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
ai_summary: "Letta Agent ↔ Todoist Integration"
---

# Letta Agent ↔ Todoist Integration

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                      LETTA AGENTS                               │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│  │ HP-Conductor│ │ LE-Conductor│ │ GE-Conductor│               │
│  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘               │
│         │               │               │                       │
│         └───────────────┼───────────────┘                       │
│                         │                                       │
│              ┌──────────▼──────────┐                           │
│              │   TODOIST TOOLS     │                           │
│              │  (MCP or Custom)    │                           │
│              └──────────┬──────────┘                           │
└─────────────────────────┼───────────────────────────────────────┘
                          │
              ┌───────────▼───────────┐
              │     TODOIST API       │
              │  - get_tasks          │
              │  - create_task        │
              │  - complete_task      │
              │  - update_task        │
              └───────────┬───────────┘
                          │
              ┌───────────▼───────────┐
              │    OMCS PROJECTS      │
              │  00-09 Command        │
              │  10-19 HP             │
              │  20-29 LE             │
              │  30-39 GE             │
              │  40-49 IT             │
              └───────────────────────┘
```

---

## 1. Todoist Tool for Letta

### Option A: MCP Server (Recommended)

```json
{
  "mcpServers": {
    "todoist": {
      "command": "npx",
      "args": ["-y", "@abhiz123/todoist-mcp-server"],
      "env": {
        "TODOIST_API_TOKEN": "3e3819010eeb7749f24409f1dd4bb1a2951d4c9b"
      }
    }
  }
}
```

**Available MCP Tools:**
- `todoist_get_tasks` - Fetch tasks by filter
- `todoist_create_task` - Create new task
- `todoist_update_task` - Modify task
- `todoist_complete_task` - Mark complete
- `todoist_get_projects` - List projects

### Option B: Custom Letta Tool

```python
# tools/todoist_tool.py
from letta import tool
import requests

TODOIST_TOKEN = "3e3819010eeb7749f24409f1dd4bb1a2951d4c9b"
HEADERS = {"Authorization": f"Bearer {TODOIST_TOKEN}"}

@tool
def get_my_tasks(project_name: str = None, filter_query: str = "today") -> str:
    """
    Get tasks from Todoist.

    Args:
        project_name: Optional project like "10 Base" or "30-39 GE"
        filter_query: Todoist filter like "today", "p1", "7 days"

    Returns:
        List of tasks with IDs, content, due dates
    """
    url = "https://api.todoist.com/rest/v2/tasks"
    params = {"filter": filter_query}

    response = requests.get(url, headers=HEADERS, params=params)
    tasks = response.json()

    # Filter by project if specified
    if project_name:
        tasks = [t for t in tasks if project_name in t.get('project_name', '')]

    return format_tasks(tasks)

@tool
def create_task(
    content: str,
    project: str,
    section: str = None,
    priority: int = 1,
    due_string: str = None,
    description: str = None,
    labels: list = None
) -> str:
    """
    Create a task in Todoist.

    Args:
        content: Task title
        project: Project name like "10 Base", "32 Operations"
        section: Section name within project
        priority: 1-4 (4 = urgent)
        due_string: "today", "tomorrow", "next monday"
        description: Task details with vault links
        labels: List of labels like ["@hp-conductor"]

    Returns:
        Created task ID and confirmation
    """
    # Get project ID
    project_id = get_project_id(project)
    section_id = get_section_id(project_id, section) if section else None

    data = {
        "content": content,
        "project_id": project_id,
        "priority": priority
    }

    if section_id:
        data["section_id"] = section_id
    if due_string:
        data["due_string"] = due_string
    if description:
        data["description"] = description
    if labels:
        data["labels"] = labels

    response = requests.post(
        "https://api.todoist.com/rest/v2/tasks",
        headers=HEADERS,
        json=data
    )

    return f"Created task: {response.json()['id']}"

@tool
def complete_task(task_id: str) -> str:
    """Mark a task as complete."""
    response = requests.post(
        f"https://api.todoist.com/rest/v2/tasks/{task_id}/close",
        headers=HEADERS
    )
    return "Task completed" if response.status_code == 204 else "Failed"
```

---

## 2. Agent ↔ Documentation Connection

### Task Description Template

When agents create tasks, include vault links in the description:

```python
description = """
## Context
HP-Conductor detected elevated stress levels from Oura data.

## Related Docs
- [[15 Performance (Body & Biology)/Recovery Protocols]]
- [[15 Performance (Body & Biology)/Mental Wellness & Practice]]

## Source
Agent: @hp-conductor
Skill: stress-detection
Timestamp: 2024-12-23T23:30:00Z

## Approval
Status: @approval-needed
"""
```

### Project-to-Folder Mapping

| Todoist Project | Vault Folder | Agent |
|-----------------|--------------|-------|
| 10 Base | `10-19 HP/10 Base (Personal Hub)` | @hp-conductor |
| 15 Performance | `10-19 HP/15 Performance (Body & Biology)` | @hp-conductor |
| 21 Sanctuary | `20-29 LE/21 Sanctuary (Home & Property)` | @le-conductor |
| 32 Operations | `30-39 GE/32 Operations (Execution & Entities)` | @ge-conductor |
| 44 Development | `40-49 IT/44 Development (Code & AI)` | @it-conductor |

---

## 3. Agent Workflow

### HP-Conductor Example

```python
# System prompt snippet for HP-Conductor

"""
You are HP-Conductor, managing Holistic Performance tasks.

TODOIST PROJECTS YOU OWN:
- 10 Base (Personal Hub)
- 11 Self (Identity & Purpose)
- 12 Lifestyle (Daily Operations)
- 13 Mind (Education & Learning)
- 14 Mastery (Skills & Training)
- 15 Performance (Body & Biology)
- 16 Finance (Money & Assets)
- 17 Sovereignty (Privacy & Protection)
- 18 Relationships (Social Network)

WORKFLOW:
1. Check tasks: get_my_tasks(filter_query="##10-19 HP & today")
2. For each task, read related vault docs
3. Execute or escalate based on approval status
4. Update task with progress
5. Complete when done

TASK CREATION RULES:
- Always include vault links in description
- Add @hp-conductor label
- Add @approval-needed for actions requiring Zeus
- Set appropriate priority (p1-p4)
- Link to relevant SOPs in 05 Protocols
"""
```

### Daily Agent Loop

```python
async def hp_conductor_daily_loop():
    # 1. Get today's HP tasks
    tasks = await get_my_tasks(filter_query="##10-19 HP & today")

    for task in tasks:
        # 2. Check if auto-approved or needs human
        if "@approval-needed" in task.labels:
            continue  # Skip, wait for Zeus

        # 3. Read related documentation
        vault_path = map_project_to_vault(task.project)
        context = await read_vault_folder(vault_path)

        # 4. Execute task
        result = await execute_task(task, context)

        # 5. Update task with result
        await update_task(task.id, description=f"Result: {result}")

        # 6. Complete if successful
        if result.success:
            await complete_task(task.id)
        else:
            # Create follow-up task
            await create_task(
                content=f"Follow up: {task.content}",
                project=task.project,
                labels=["@approval-needed"],
                description=f"Previous attempt failed: {result.error}"
            )
```

---

## 4. Task Routing by Label

| Label | Meaning | Action |
|-------|---------|--------|
| @hp-conductor | HP agent owns this | HP agent can execute |
| @le-conductor | LE agent owns this | LE agent can execute |
| @ge-conductor | GE agent owns this | GE agent can execute |
| @it-conductor | Grid agent owns this | Grid agent can execute |
| @approval-needed | Requires Zeus review | Agent waits |
| @auto-approved | Pre-approved action | Agent executes immediately |
| @blocked | Waiting on external | Agent monitors |
| @urgent | Time-sensitive | Agent prioritizes |

---

## 5. Integration Points

### Vault → Todoist (Agent Creates Tasks)

```python
# Agent finds something that needs action
await create_task(
    content="Schedule dentist appointment",
    project="15 Performance",
    section="Biometrics",
    priority=2,
    due_string="next week",
    labels=["@hp-conductor", "@approval-needed"],
    description="""
## Why
Last dental visit was 8 months ago per [[15 Performance/Labs & Testing]].

## Action Required
Zeus to call dentist and schedule.

## Related
- [[15 Performance/Holistic Performance Wiki]]
"""
)
```

### Todoist → Vault (Agent Reads Tasks)

```python
# Agent checks for work
tasks = await get_my_tasks(filter_query="##15 Performance & @hp-conductor")

for task in tasks:
    # Extract vault links from description
    vault_links = extract_links(task.description)

    # Read context
    for link in vault_links:
        doc = await read_file(link)
        context.append(doc)

    # Execute with context
    await execute_with_context(task, context)
```

### Obsidian ↔ Todoist (Bidirectional Sync)

The Todoist plugin syncs task status:
- Complete in Obsidian → Completes in Todoist
- Complete in Todoist → Updates in Obsidian view
- New task in Todoist → Appears in query blocks

---

## 6. Security & Approval Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Agent     │────▶│   Todoist   │────▶│    Zeus     │
│  Creates    │     │   Task      │     │   Reviews   │
│   Task      │     │ @approval-  │     │  in App     │
│             │     │   needed    │     │             │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                    ┌─────────────┐            │
                    │   Agent     │◀───────────┘
                    │  Monitors   │  Zeus removes
                    │  for        │  @approval-needed
                    │  approval   │  or adds
                    │             │  @auto-approved
                    └─────────────┘
```

---

## 7. Example: Full Flow

**Scenario**: HP-Conductor detects Zeus hasn't logged meals in 3 days.

1. **Agent creates task**:
```python
create_task(
    content="Log meals for past 3 days",
    project="15 Performance",
    section="Nutrition",
    priority=2,
    due_string="today",
    labels=["@hp-conductor", "@approval-needed"],
    description="""
## Detection
No meal entries in [[15 Performance/Nutrition & Supplement]] since Dec 20.

## Suggested Action
Open Cronometer and log meals for Dec 21-23.

## Health Impact
Tracking gap affects nutrition analysis and supplement timing.
"""
)
```

2. **Zeus sees in Todoist app** → Reviews → Completes task manually

3. **Agent detects completion** → Updates internal state → Logs to memory

---

## Related Files

- [[00-Master-Directory]] - Task engine overview
- [[../01.02 Letta Agents/Tool Registry]] - Available tools
- [[../01.02 Letta Agents/HITL]] - Human-in-the-loop protocols
- [[../../00 System (Core Infrastructure)/Api-integrations/todoist-config]] - API config
