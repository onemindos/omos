---
id: ui-12-20260111053926
title: "Todoist Integration Configuration"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/12]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Todoist Integration Configuration"
---

# Todoist Integration Configuration

## Overview

Todoist serves as the **Task Execution Engine** for the One Mind Hive Mind. All agents route actionable tasks through Todoist for Zeus's review and completion.

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                      HIVE MIND AGENTS                        │
│  ┌─────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐   │
│  │ LEGACY  │ │   HP   │ │   LE   │ │   GE   │ │   IT   │   │
│  └────┬────┘ └───┬────┘ └───┬────┘ └───┬────┘ └───┬────┘   │
│       │          │          │          │          │         │
│       └──────────┴──────────┼──────────┴──────────┘         │
│                             │                                │
│                    ┌────────▼────────┐                      │
│                    │  HITL Approval  │                      │
│                    │     Queue       │                      │
│                    └────────┬────────┘                      │
│                             │                                │
└─────────────────────────────┼────────────────────────────────┘
                              │
                     ┌────────▼────────┐
                     │    TODOIST      │
                     │  (Task Engine)  │
                     └────────┬────────┘
                              │
                     ┌────────▼────────┐
                     │      ZEUS       │
                     │   (Execution)   │
                     └─────────────────┘
```

## Todoist Project Structure

| Project | Agent Owner | Purpose |
|---------|-------------|---------|
| `🎯 Inbox` | LEGACY | Unsorted/triaged items |
| `📊 HP - Performance` | HP Agent | Health, habits, learning |
| `👨‍👩‍👧‍👦 LE - Family` | LE Agent | Family, home, legacy |
| `💼 GE - Business` | GE Agent | Business, revenue, ops |
| `🖥️ IT - Technology` | Grid Agent | Infrastructure, tools |
| `🔴 @Critical` | LEGACY | L0 emergency escalations |

## Label Schema

| Label | Priority | Agent Action |
|-------|----------|--------------|
| `@critical` | P1 | Immediate Zeus attention |
| `@today` | P2 | Same-day completion |
| `@this-week` | P3 | Week completion target |
| `@someday` | P4 | Backlog/future |
| `@waiting` | — | Blocked on external |
| `@review` | — | Needs Zeus decision |
| `@auto-approved` | — | Agent can execute |

## Integration Endpoints

### MCP Server Configuration

```json
{
  "name": "todoist",
  "type": "mcp",
  "config": {
    "command": "npx",
    "args": ["-y", "@anthropic/mcp-todoist"],
    "env": {
      "TODOIST_API_TOKEN": "${TODOIST_API_TOKEN}"
    }
  },
  "tools": [
    "todoist_get_tasks",
    "todoist_create_task",
    "todoist_update_task",
    "todoist_complete_task",
    "todoist_get_projects"
  ]
}
```

### N8N Webhook Integration

```yaml
# Todoist → N8N Trigger
webhook_url: https://n8n.your-domain.com/webhook/todoist
events:
  - item:added
  - item:completed
  - item:updated
```

## Task Creation Protocol

When an agent creates a task:

1. **Determine Project**: Map to appropriate pillar project
2. **Set Priority**: Based on urgency assessment
3. **Add Labels**: Include agent source and approval status
4. **Set Due Date**: If time-sensitive
5. **Add Description**: Include context and rationale

### Example Task Creation

```python
# HP Agent creating a task
task = {
    "content": "Schedule annual physical exam",
    "project": "HP - Performance",
    "priority": 3,
    "labels": ["@this-week", "@hp-agent", "@review"],
    "due_string": "next Monday",
    "description": """
    ## Context
    Last physical was 11 months ago. Annual checkup recommended.

    ## Source
    HP Agent - health-management skill

    ## Approval
    Awaiting Zeus review
    """
}
```

## Sync Schedule

| Sync Type | Frequency | Purpose |
|-----------|-----------|---------|
| Task Pull | Every 5 min | Get latest task state |
| Completion Sync | Real-time | Update agent context |
| Daily Digest | 6:00 AM | Morning brief generation |
| Weekly Review | Sunday 9 AM | Week-ahead planning |

## Environment Variables

```bash
# .env configuration
TODOIST_API_TOKEN=your_todoist_api_token
TODOIST_INBOX_PROJECT_ID=123456789
TODOIST_HP_PROJECT_ID=123456790
TODOIST_LE_PROJECT_ID=123456791
TODOIST_GE_PROJECT_ID=123456792
TODOIST_IT_PROJECT_ID=123456793
```

## Error Handling

| Error | Response | Escalation |
|-------|----------|------------|
| API Rate Limit | Backoff + retry | Log to Grid agent |
| Auth Failure | Notify Grid agent | Critical alert |
| Sync Conflict | Todoist wins | Log discrepancy |
| Task Creation Fail | Retry 3x | Fall back to note |

## Security

- API token stored in secure environment variables
- Never log task content in plain text
- Rate limiting enforced on all API calls
- Audit trail maintained for all task operations


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
