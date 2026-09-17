# Todoist API Reference

## Base Configuration

```bash
# API Base URL
https://api.todoist.com/rest/v2

# Authentication header
Authorization: Bearer $TODOIST_API_TOKEN
```

## Endpoints

### Tasks

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks` | List tasks (optional `?filter=`) |
| GET | `/tasks/{id}` | Get task |
| POST | `/tasks` | Create task |
| POST | `/tasks/{id}` | Update task |
| POST | `/tasks/{id}/close` | Complete task |
| POST | `/tasks/{id}/reopen` | Reopen task |
| DELETE | `/tasks/{id}` | Delete task |

### Projects

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects` | List projects |
| GET | `/projects/{id}` | Get project |
| POST | `/projects` | Create project |
| POST | `/projects/{id}` | Update project |
| DELETE | `/projects/{id}` | Delete project |

### Sections

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/sections` | List all sections |
| GET | `/sections?project_id={id}` | List project sections |
| GET | `/sections/{id}` | Get section |
| POST | `/sections` | Create section |
| POST | `/sections/{id}` | Update section |
| DELETE | `/sections/{id}` | Delete section |

### Labels

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/labels` | List labels |
| GET | `/labels/{id}` | Get label |
| POST | `/labels` | Create label |
| POST | `/labels/{id}` | Update label |
| DELETE | `/labels/{id}` | Delete label |

### Comments

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/comments?task_id={id}` | List task comments |
| POST | `/comments` | Create comment |
| DELETE | `/comments/{id}` | Delete comment |

## Request/Response Formats

### Create Task Request
```json
{
  "content": "Task title",
  "description": "Task details",
  "project_id": "2203306141",
  "section_id": "7025",
  "parent_id": null,
  "order": 1,
  "labels": ["urgent", "hp-conductor"],
  "priority": 4,
  "due_string": "tomorrow at 3pm",
  "due_date": "2025-01-15",
  "due_datetime": "2025-01-15T15:00:00Z",
  "due_lang": "en",
  "assignee_id": null,
  "duration": 30,
  "duration_unit": "minute"
}
```

### Task Response
```json
{
  "id": "2995104339",
  "project_id": "2203306141",
  "section_id": "7025",
  "content": "Task title",
  "description": "Task details",
  "is_completed": false,
  "labels": ["urgent", "hp-conductor"],
  "priority": 4,
  "comment_count": 0,
  "created_at": "2025-01-10T12:00:00Z",
  "due": {
    "date": "2025-01-15",
    "string": "Jan 15",
    "datetime": "2025-01-15T15:00:00Z",
    "is_recurring": false
  },
  "url": "https://todoist.com/showTask?id=2995104339"
}
```

### Create Project Request
```json
{
  "name": "Project Name",
  "parent_id": "2203306141",
  "color": "blue",
  "is_favorite": false,
  "view_style": "list"
}
```

### Project Response
```json
{
  "id": "2203306141",
  "name": "25-49 HP (Holistic Performance)",
  "color": "blue",
  "parent_id": null,
  "order": 1,
  "comment_count": 0,
  "is_shared": false,
  "is_favorite": false,
  "is_inbox_project": false,
  "is_team_inbox": false,
  "view_style": "list",
  "url": "https://todoist.com/showProject?id=2203306141"
}
```

### Create Section Request
```json
{
  "name": "In Progress",
  "project_id": "2203306141",
  "order": 3
}
```

### Section Response
```json
{
  "id": "7025",
  "project_id": "2203306141",
  "order": 3,
  "name": "In Progress"
}
```

### Create Label Request
```json
{
  "name": "approval-needed",
  "color": "yellow",
  "order": 1,
  "is_favorite": false
}
```

### Label Response
```json
{
  "id": "2156154810",
  "name": "approval-needed",
  "color": "yellow",
  "order": 1,
  "is_favorite": false
}
```

## One Mind Project Structure

### Pillar Hierarchy
```
00-24 SO (Sovereign Operations)
├── UI.00 Framework
├── UI.01 Command
├── UI.02 Agents
└── ...

25-49 HP (Holistic Performance)
├── HP.25 Identity
├── HP.26 Mind
├── HP.27 Body
└── ...

50-74 LE (Legacy Evolution)
├── LE.50 Home
├── LE.51 Marriage
├── LE.52 Children
└── ...

75-99 GE (Generational Empire)
├── GE.75 Zeus Delacruz Enterprise
├── GE.76 Legacy Acadacom
├── GE.77 Legacy Platform Holdings
└── ...
```

### Standard Agent Labels

**Command (UI):**
- `mother-ai`
- `protocol-guardian`
- `approval-manager`
- `memory-curator`

**HP Path:**
- `hp-conductor`
- `hp-health`
- `hp-fitness`
- `hp-nutrition`
- `hp-sleep`
- `hp-mental`
- `hp-habits`

**LE Path:**
- `le-conductor`
- `le-identity`
- `le-relationships`
- `le-growth`
- `le-legacy`

**GE Path:**
- `ge-conductor`
- `ge-ventures`
- `ge-products`
- `ge-marketing`
- `ge-operations`
- `ge-finance`

### Status Labels

| Label | Color | Purpose |
|-------|-------|---------|
| `approval-needed` | yellow | Awaiting human approval |
| `blocked` | grey | Blocked by dependency |
| `urgent` | red | Requires immediate attention |
| `waiting` | grey | Waiting on external input |
| `scheduled` | sky_blue | Scheduled for specific time |

## Colors

Available colors for projects and labels:
```
berry_red, red, orange, yellow, olive_green, lime_green,
green, mint_green, teal, sky_blue, light_blue, blue,
grape, violet, lavender, magenta, salmon, charcoal, grey
```

## Rate Limits

- **450 requests per 15 minutes** per user
- Response includes `X-Limit-*` headers
- 429 response when exceeded

## URL Encoding for Filters

| Character | Encoded |
|-----------|---------|
| `@` | `%40` |
| `#` | `%23` |
| `&` | `%26` |
| space | `%20` |
| `(` | `%28` |
| `)` | `%29` |

Example: `@approval-needed & ##25-49 HP` becomes `%40approval-needed%20%26%20%23%2325-49%20HP`


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
