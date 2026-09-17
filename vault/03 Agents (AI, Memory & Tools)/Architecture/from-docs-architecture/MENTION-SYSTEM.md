# @ Mention System

## Overview

Universal @ mention system for referencing entities across OneMind OS. Supports tasks, pages, goals, projects, habits, people (family, contacts, collaborators), agents, teams, and workflows.

## Syntax

### Typed Mentions (Recommended)

```
@type:identifier
@type:"Identifier With Spaces"
```

| Prefix | Entity Type | Example |
|--------|-------------|---------|
| `@task:` | Task | `@task:Fix-login-bug` |
| `@page:` | Page/Document | `@page:Deployment-SOP` |
| `@goal:` | Goal | `@goal:Revenue-100k` |
| `@project:` | Project | `@project:OneMind-v2` |
| `@habit:` | Habit | `@habit:Morning-workout` |
| `@person:` | Any person (family, contacts, etc.) | `@person:Noah` |
| `@agent:` | AI Agent | `@agent:researcher` |
| `@team:` | Agent team | `@team:research-team` |
| `@workflow:` | Workflow | `@workflow:deep-research` |

### Aliases

| Alias | Resolves To |
|-------|-------------|
| `@contact:` | `@person:` (backwards compatibility) |

### Legacy Shorthand (Agents Only)

For convenience, agents can be mentioned without the `@agent:` prefix:

| Shorthand | Resolves To |
|-----------|-------------|
| `@researcher` | `@agent:research-agent` |
| `@coder` | `@agent:code-agent` |
| `@writer` | `@agent:writing-agent` |
| `@planner` | `@agent:planning-agent` |
| `@farm` | `@agent:farm-manager-agent` |
| `@health` | `@agent:health-agent` |
| `@finance` | `@agent:finance-agent` |

---

## Multiple Mentions

You can mention multiple entities in a single message:

```
@researcher @coder please review @task:Fix-login-bug for @project:OneMind-v2
```

This parses to:
- 2 agents: `research-agent`, `code-agent`
- 1 task: `Fix-login-bug`
- 1 project: `OneMind-v2`

The system returns all mentions grouped by type.

---

## Identifier Format

Identifiers use slug format:
- Words separated by dashes: `Fix-login-bug`
- Case-insensitive matching
- Spaces in names converted to dashes

**Converting Names to Identifiers:**

| Entity Name | Mention Identifier |
|-------------|-------------------|
| "Fix login bug" | `@task:Fix-login-bug` |
| "Q1 Revenue Report" | `@page:Q1-Revenue-Report` |
| "Morning Workout" | `@habit:Morning-workout` |

**Quoted Identifiers:**

For names with special characters, use quotes:

```
@task:"Fix bug #123"
@contact:"John Smith Jr."
```

---

## API Endpoints

### Parse Mentions

Parse text without database lookup.

```http
POST /api/mentions/parse
Content-Type: application/json

{
  "text": "@researcher find docs about @task:Fix-login-bug"
}
```

**Response:**

```json
{
  "original": "@researcher find docs about @task:Fix-login-bug",
  "clean_text": "find docs about",
  "mentions": [
    {
      "type": "agent",
      "identifier": "research-agent",
      "display_name": "Researcher",
      "raw": "@researcher",
      "start_pos": 0,
      "end_pos": 11
    },
    {
      "type": "task",
      "identifier": "Fix-login-bug",
      "display_name": "Fix Login Bug",
      "raw": "@task:Fix-login-bug",
      "start_pos": 28,
      "end_pos": 47
    }
  ],
  "has_mentions": true,
  "agents": [...],
  "tasks": [...]
}
```

### Resolve Mentions

Parse and look up entities in database.

```http
POST /api/mentions/resolve
Content-Type: application/json

{
  "text": "@task:Fix-login-bug for @project:OneMind",
  "user_id": "zeus"
}
```

**Response:**

```json
[
  {
    "mention": {
      "type": "task",
      "identifier": "Fix-login-bug",
      "display_name": "Fix Login Bug",
      "raw": "@task:Fix-login-bug",
      "start_pos": 0,
      "end_pos": 19
    },
    "resolved": true,
    "entity_id": "550e8400-e29b-41d4-a716-446655440000",
    "entity_data": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Fix login bug on mobile",
      "status": "in_progress",
      "priority": "high",
      "task_pillar": "ge"
    },
    "error": null
  },
  {
    "mention": {...},
    "resolved": true,
    "entity_id": "...",
    "entity_data": {...}
  }
]
```

### Autocomplete

Get suggestions while typing.

```http
POST /api/mentions/autocomplete
Content-Type: application/json

{
  "query": "@task:fix",
  "user_id": "zeus",
  "limit": 10
}
```

**Response:**

```json
{
  "query": "@task:fix",
  "suggestions": [
    {
      "mention": "@task:fix-login-bug",
      "display": "Fix login bug",
      "type": "task",
      "subtitle": "Task"
    },
    {
      "mention": "@task:fix-database-timeout",
      "display": "Fix database timeout",
      "type": "task",
      "subtitle": "Task"
    }
  ]
}
```

### Get Types

List all supported mention types.

```http
GET /api/mentions/types
```

---

## Usage Examples

### Chat Messages

```
User: @researcher @coder can you review the code in @project:OneMind-v2?
```

System parses:
- Routes to: `research-agent`, `code-agent`
- Context: `OneMind-v2` project

### Task Descriptions

```
This task is for @person:Noah's birthday party
Related to @goal:Family-events
```

### Agent Instructions

```
@farm check on the chickens and update @habit:Animal-care
```

### Multi-Agent Collaboration

```
@researcher gather information about AI agents
@writer create a summary
@planner create action items from the research
```

All three agents receive the message with their portion of the work.

---

## Code Usage

### Python

```python
from backend.agno.routing.mentions import (
    MentionParser,
    MentionResolver,
    parse_mentions,
    resolve_mentions,
    create_mention,
    slugify,
)

# Quick parse
parsed = parse_mentions("@researcher find @task:Fix-bug")
print(parsed.agents)  # [Mention(type='agent', identifier='research-agent', ...)]
print(parsed.tasks)   # [Mention(type='task', identifier='Fix-bug', ...)]
print(parsed.clean_text)  # "find"

# Full resolution
resolved = await resolve_mentions("@task:Fix-bug", user_id="zeus")
for r in resolved:
    if r.resolved:
        print(f"Found: {r.entity_data}")

# Create mentions programmatically
mention = create_mention(MentionType.TASK, "Fix login bug")
# Returns: "@task:fix-login-bug"
```

### Dart/Flutter

```dart
// Parse in frontend before sending to backend
final response = await api.post('/mentions/parse', {
  'text': messageController.text,
});

final mentions = response.data['mentions'];
final agents = mentions.where((m) => m['type'] == 'agent').toList();

if (agents.isNotEmpty) {
  // Route to specific agents
  await routeToAgents(agents, cleanText);
}
```

---

## Database Queries

### How Resolution Works

Each entity type has its own resolver that queries the appropriate table:

| Type | Table | Match Logic |
|------|-------|-------------|
| `task` | `inbox_items` | Title ILIKE or slug match |
| `project` | `projects` | Name ILIKE or slug match |
| `goal` | `goals` | Name ILIKE or slug match |
| `habit` | `habits` | Name ILIKE, active only |
| `person` | `learning_entities` | Name + type='person' OR bridge_entity_type IN ('child', 'contact', 'person') |
| `workflow` | `workflows` | Name ILIKE, active only |
| `agent` | (static) | Known agent IDs |
| `team` | (static) | Known team IDs |
| `page` | `pages` | Title ILIKE or slug match |

### Example Query (Tasks)

```sql
SELECT id, title, body, status, priority, task_pillar, project_id
FROM inbox_items
WHERE user_id = $1 AND type = 'task'
AND (
    title ILIKE '%fix login bug%'
    OR LOWER(REPLACE(title, ' ', '-')) = 'fix-login-bug'
)
ORDER BY
    CASE WHEN LOWER(title) = 'fix login bug' THEN 0 ELSE 1 END,
    created_at DESC
LIMIT 1
```

---

## Files

| File | Purpose |
|------|---------|
| [mentions.py](backend/agno/routing/mentions.py) | Parser & resolver classes |
| [mentions_api.py](backend/agno/routing/mentions_api.py) | REST API endpoints |

---

## Integration with Router

The mention system integrates with the existing message router:

```python
from backend.agno.routing.mentions import parse_mentions

# In message handler
parsed = parse_mentions(user_message)

if parsed.agents:
    # Route to specific agents
    for agent in parsed.agents:
        await route_to_agent(agent.identifier, parsed.clean_text)
else:
    # Route to default (Legacy)
    await route_to_legacy(parsed.clean_text)

# Inject resolved context for tasks/goals/etc.
if parsed.tasks or parsed.goals or parsed.projects:
    context = await build_entity_context(parsed)
    await inject_context(context)
```

---

## Future Enhancements

1. **Rich Previews** - Show entity cards on hover
2. **Bidirectional Links** - Track what mentions what
3. **Mention Notifications** - Alert when entity is mentioned
4. **User Mentions** - @user:name for multi-user support
