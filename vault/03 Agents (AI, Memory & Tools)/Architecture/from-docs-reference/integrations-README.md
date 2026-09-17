# ClickUp Sync Service

> Local sync layer for ClickUp data ownership and agent access

## Overview

This service maintains a real-time local copy of all ClickUp data, ensuring:
1. **Data ownership** - You own your data locally, not just in the cloud
2. **Disaster recovery** - If ClickUp disappears, your data is safe
3. **Agent optimization** - Fast local queries for OneMind agents
4. **Offline access** - Query data without internet

## Architecture

```
clickup/
├── README.md              # This file
├── __init__.py            # Package init
├── config.py              # Configuration management
├── client/
│   ├── __init__.py
│   └── api.py             # ClickUp API wrapper
├── models/
│   ├── __init__.py
│   └── schemas.py         # Pydantic models for ClickUp data
├── sync/
│   ├── __init__.py
│   ├── webhook_receiver.py  # FastAPI webhook endpoint
│   ├── full_sync.py         # Complete workspace sync
│   └── export.py            # JSON/Markdown export
├── storage/
│   ├── __init__.py
│   ├── database.py        # PostgreSQL operations
│   └── filesystem.py      # File storage (attachments, docs)
└── tools/
    ├── __init__.py
    └── agno_tool.py       # Agno tool integration
```

## Quick Start

### 1. Environment Setup

Add to `.env`:

```bash
# Required
CLICKUP_API_TOKEN=pk_xxxxx_xxxxx
CLICKUP_TEAM_ID=12345678

# Optional
CLICKUP_WEBHOOK_SECRET=your_secret
CLICKUP_SYNC_INTERVAL_MINUTES=60
```

### 2. Database Setup

```sql
-- Run migrations
psql -d onemind -f clickup/migrations/001_initial.sql
```

### 3. Start Sync Service

```bash
# As part of OneMind OS
uvicorn backend.onemindos:app --reload

# Webhook endpoint available at:
# POST /webhooks/clickup
```

### 4. Register Webhook in ClickUp

```bash
# Using the CLI helper
python -m clickup.cli register-webhook https://your-domain.com/webhooks/clickup
```

## Features

### Real-time Sync (Webhooks)

ClickUp sends events to our webhook endpoint:
- Task created/updated/deleted
- Comments posted
- Status changes
- Assignment changes
- Due date changes

### Full Sync (Hourly)

Cron job verifies all data:
- Catches any missed webhook events
- Syncs new spaces/folders/lists
- Downloads new attachments
- Exports docs to markdown

### Export (Daily)

Full snapshot for disaster recovery:
- JSON export of entire workspace
- Markdown export of all docs
- Attachment backup

## API Usage

### Python Client

```python
from clickup.client import ClickUpClient

async with ClickUpClient() as client:
    # Get all tasks
    tasks = await client.get_tasks(list_id="123")

    # Create task
    task = await client.create_task(
        list_id="123",
        name="New task",
        description="Description here"
    )

    # Search across workspace
    results = await client.search("auth bug")
```

### Agno Tool

```python
from clickup.tools import clickup_tool

# Available in agents
agent = Agent(
    tools=[clickup_tool],
    instructions="You can manage tasks in ClickUp"
)
```

## Data Storage

### PostgreSQL Tables

| Table | Purpose |
|-------|---------|
| `clickup_workspaces` | Workspace metadata |
| `clickup_spaces` | Spaces within workspaces |
| `clickup_folders` | Folders within spaces |
| `clickup_lists` | Lists within folders/spaces |
| `clickup_tasks` | All tasks with full data |
| `clickup_comments` | Task comments |
| `clickup_docs` | Document metadata |
| `clickup_attachments` | Attachment metadata |
| `clickup_goals` | Goals and key results |
| `clickup_time_entries` | Time tracking data |
| `clickup_sync_log` | Sync audit trail |

### Filesystem

```
/data/clickup/
├── attachments/           # Downloaded files
│   └── {task_id}/
│       └── {filename}
├── docs/                  # Markdown exports
│   └── {doc_id}.md
└── exports/               # Daily snapshots
    └── {YYYY-MM-DD}/
        ├── workspace.json
        ├── tasks.json
        └── docs.json
```

## Webhook Events

Events we subscribe to:

```python
SUBSCRIBED_EVENTS = [
    "taskCreated", "taskUpdated", "taskDeleted",
    "taskPriorityUpdated", "taskStatusUpdated",
    "taskAssigneeUpdated", "taskDueDateUpdated",
    "taskCommentPosted", "taskCommentUpdated",
    "taskTimeTrackedUpdated",
    "listCreated", "listUpdated", "listDeleted",
    "folderCreated", "folderUpdated", "folderDeleted",
    "spaceCreated", "spaceUpdated", "spaceDeleted",
    "goalCreated", "goalUpdated", "goalDeleted",
]
```

## Rate Limits

| Plan | Limit | Our Usage |
|------|-------|-----------|
| Unlimited | 100 req/min | ~20-30 req/min typical |
| Business Plus | 1,000 req/min | Comfortable headroom |

Full sync of 1,000 tasks: ~15 minutes on Unlimited plan.

## Monitoring

### Sync Status

```bash
# Check sync health
python -m clickup.cli status

# Output:
# Last webhook: 2025-02-02 10:30:15 (2 min ago)
# Last full sync: 2025-02-02 10:00:00
# Local tasks: 1,234
# Local docs: 56
# Sync lag: 0 events pending
```

### Logs

```bash
# Tail sync logs
tail -f logs/clickup_sync.log
```

## Troubleshooting

### Webhook not receiving events

1. Check webhook URL is accessible
2. Verify CLICKUP_WEBHOOK_SECRET matches
3. Check ClickUp webhook status in settings

### Sync lag increasing

1. Check rate limits (429 errors in logs)
2. Verify database connectivity
3. Check disk space for attachments

### Missing data

1. Run manual full sync: `python -m clickup.cli full-sync`
2. Check sync_log for errors
3. Verify API token permissions
