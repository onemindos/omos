# Composio Integration Guide

> OneMind OS + Composio: External Service Integration Layer

## Overview

Composio provides managed OAuth, 500+ tool integrations, and webhook triggers for AI agents. This integration connects Composio to Legacy's consciousness layer, enabling:

1. **INBOUND** - External events flow into Legacy's awareness via triggers
2. **OUTBOUND** - Legacy can execute actions on external services

## Architecture

```
                         EXTERNAL SERVICES
                    (Todoist, GitHub, Slack, etc.)
                                │
          ┌─────────────────────┼─────────────────────┐
          │                     │                     │
     TRIGGERS                   │                ACTIONS
    (events IN)                 │             (commands OUT)
          │                     │                     │
          ▼                     │                     ▼
┌─────────────────────┐         │      ┌─────────────────────────┐
│  COMPOSIO TRIGGERS  │         │      │   COMPOSIO ACTIONS      │
│                     │         │      │                         │
│  - Webhook mgmt     │         │      │  - 10,000+ actions      │
│  - OAuth handling   │         │      │  - Managed auth         │
│  - Event normalize  │         │      │  - Python SDK           │
└─────────┬───────────┘         │      └─────────────────────────┘
          │                     │                     ▲
          │ POST /webhooks/composio                   │
          │                     │                     │
          ▼                     │                     │
┌─────────────────────────────────────────────────────────────────┐
│                     CONSCIOUSNESS LAYER                          │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    GATEWAY (gateway.py)                    │  │
│  │  /webhooks/composio  ← Composio triggers                  │  │
│  │  /webhooks/{source}  ← Direct webhooks (Home Assistant)   │  │
│  └─────────────────────────────┬─────────────────────────────┘  │
│                                │                                 │
│                   NATS: events.{source}.{type}                   │
│                                │                                 │
│  ┌─────────────────────────────▼─────────────────────────────┐  │
│  │                     BUS (bus.py)                           │  │
│  │           NATS JetStream - The Nervous System              │  │
│  └─────────────────────────────┬─────────────────────────────┘  │
│                                │                                 │
│                   Subscribe: events.>                            │
│                                │                                 │
│  ┌─────────────────────────────▼─────────────────────────────┐  │
│  │                 PROCESSOR (processor.py)                   │  │
│  │   Noise Filter → Awareness Filter → Batch → Enrich        │  │
│  └─────────────────────────────┬─────────────────────────────┘  │
│                                │                                 │
│                   NATS: legacy.inbox                             │
└────────────────────────────────┼─────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Legacy (Agno)                            │
│                                                                  │
│  Subscribes to: legacy.inbox                                     │
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    COMPOSIO TOOLS                          │  │
│  │                                                            │  │
│  │  create_task() ──────────────────────────────────────────────►
│  │  send_slack_message() ───────────────────────────────────────►
│  │  create_github_issue() ──────────────────────────────────────►
│  │                                                            │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Installation

```bash
# Core Composio SDK
pip install composio

# If using with OpenAI-compatible (we'll wrap for Agno)
pip install composio-openai
```

## Configuration

### Environment Variables

Add to your `.env` or environment:

```bash
# Composio
COMPOSIO_API_KEY=your_composio_api_key

# Optional: Entity ID for multi-user (default: "default")
COMPOSIO_ENTITY_ID=default
```

### Composio Dashboard Setup

1. Create account at [composio.dev](https://composio.dev)
2. Get API key from dashboard
3. Connect apps (Todoist, GitHub, Slack, etc.) via OAuth
4. Note your Entity ID

## File Structure

```
agents/legacy/
├── consciousness/
│   ├── __init__.py          # Updated exports
│   ├── bus.py               # NATS event bus (unchanged)
│   ├── gateway.py           # Updated with Composio endpoint
│   ├── processor.py         # Event processor (unchanged)
│   └── composio_bridge.py   # NEW: Composio integration
│
└── tools/
    └── composio_tools.py    # NEW: Composio tools for Legacy
```

## Integration Points

### 1. INBOUND: Triggers → Consciousness

Composio triggers send webhooks when events occur in external services.

**Flow:**
```
External Service → Composio Trigger → POST /webhooks/composio → NATS → Processor → Legacy
```

**Composio Trigger Payload Format:**
```json
{
  "trigger_id": "trig_abc123",
  "trigger_name": "GITHUB_PULL_REQUEST_EVENT",
  "connection_id": "conn_xyz789",
  "payload": {
    "action": "opened",
    "pull_request": {
      "number": 42,
      "title": "Add new feature",
      "html_url": "https://github.com/..."
    }
  },
  "metadata": {
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

### 2. OUTBOUND: Legacy → Actions

When Legacy decides to take action, it uses Composio tools.

**Flow:**
```
Legacy decision → Agno Tool → ComposioToolSet.execute() → External Service
```

**Example Tool Call:**
```python
# Legacy decides to create a task
result = create_todoist_task(
    title="Review PR #42",
    due_date="tomorrow",
    priority=2
)
```

## Supported Integrations

### Priority 1 (OneMind Core)

| Service | Triggers | Actions | Notes |
|---------|----------|---------|-------|
| **Todoist** | task.created, task.completed | create_task, complete_task, update_task | Core task management |
| **GitHub** | pr.opened, issue.created, push | create_issue, create_pr, comment | Code awareness |
| **Google Calendar** | event.created, event.updated | create_event, update_event | Schedule awareness |
| **Slack** | message.received | send_message, react | Team comms |
| **Discord** | message.received | send_message | Personal comms |

### Priority 2 (Extended)

| Service | Triggers | Actions |
|---------|----------|---------|
| Gmail | email.received | send_email, reply |
| Notion | page.created | create_page, update_page |
| Linear | issue.created | create_issue, update_issue |
| Spotify | - | play, pause, get_playing |

### Keep Direct (Not via Composio)

| Service | Reason |
|---------|--------|
| **Home Assistant** | Local network, custom events, low latency |
| **LiveKit** | Voice/video, real-time streaming |

## Usage Examples

### Creating a Task from Legacy

```python
from agents.legacy.tools.composio_tools import create_todoist_task

# Legacy can call this tool when it decides to create a task
result = await create_todoist_task(
    title="Follow up on PR review",
    description="Check the feedback on PR #123",
    due_date="2024-01-16",
    priority=2,
    project="Work"
)
```

### Sending a Slack Message

```python
from agents.legacy.tools.composio_tools import send_slack_message

result = await send_slack_message(
    channel="#general",
    message="Reminder: Team standup in 15 minutes"
)
```

### Handling Inbound Trigger

The gateway automatically routes Composio webhooks:

```python
# POST /webhooks/composio receives:
{
    "trigger_name": "TODOIST_TASK_CREATED",
    "payload": {"title": "New task", "due": "tomorrow"}
}

# Gateway publishes to NATS:
# Subject: events.todoist.task.created
# Processor filters, enriches, forwards to legacy.inbox
# Legacy becomes aware of the new task
```

## Security Considerations

1. **Webhook Verification**: Composio signs webhooks - verify signatures
2. **API Key Storage**: Store COMPOSIO_API_KEY securely (env vars, secrets manager)
3. **Entity Isolation**: Use separate entity IDs for different users/contexts
4. **Scope Limiting**: Only request necessary OAuth scopes

## Monitoring

### Metrics to Track

- Trigger events received (by source)
- Actions executed (by type)
- Action success/failure rate
- Latency (trigger → Legacy awareness)

### Logging

All Composio interactions are logged:
- `consciousness.composio` - Bridge operations
- `consciousness.gateway` - Webhook receipt
- `legacy.tools.composio` - Action execution

## Troubleshooting

### Triggers Not Arriving

1. Check Composio dashboard for trigger status
2. Verify webhook URL is publicly accessible
3. Check gateway logs for incoming requests
4. Verify NATS connection

### Actions Failing

1. Check Composio API key validity
2. Verify OAuth connection for the service
3. Check action parameters match expected schema
4. Review Composio dashboard for error details

### Authentication Issues

1. Re-authenticate the app in Composio dashboard
2. Check OAuth token expiry
3. Verify entity ID matches

## Migration Path

### Phase 1: Setup (Current)
- [x] Create Composio account
- [x] Install SDK
- [x] Create bridge module
- [x] Create tools module

### Phase 2: Connect Apps
- [ ] Connect Todoist via OAuth
- [ ] Connect GitHub via OAuth
- [ ] Connect Slack via OAuth
- [ ] Configure triggers

### Phase 3: Deploy
- [ ] Add environment variables to production
- [ ] Expose webhook endpoint publicly
- [ ] Test trigger → Legacy flow
- [ ] Test action execution

### Phase 4: Expand
- [ ] Add more services
- [ ] Create service-specific tools
- [ ] Optimize trigger filtering


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
