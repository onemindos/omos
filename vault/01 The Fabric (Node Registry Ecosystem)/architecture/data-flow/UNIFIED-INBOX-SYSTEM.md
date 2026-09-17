# Unified Inbox System

> **Status:** LIVE & WORKING
> **Last Updated:** 2026-01-23
> **Real-time:** WebSocket broadcasting enabled

---

## Overview

The Unified Inbox is the central notification and task management hub for OneMind OS. It provides:

- **Real-time notifications** via WebSocket (no page refresh needed)
- **Task assignment** to Legacy agent (human-in-the-loop)
- **Approval workflows** for sensitive actions
- **Multi-channel delivery** via Novu Cloud (push, email, SMS)
- **In-app inbox** in Flutter UI

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         UNIFIED INBOX ARCHITECTURE                           │
└─────────────────────────────────────────────────────────────────────────────┘

External Sources                Internal Sources
     │                               │
     ▼                               ▼
┌─────────────┐              ┌─────────────┐
│ Novu Cloud  │              │ Legacy Agent│
│ (Push/Email)│              │ (Notifs)    │
└──────┬──────┘              └──────┬──────┘
       │                            │
       └────────────┬───────────────┘
                    │
                    ▼
         ┌──────────────────┐
         │   INBOX API      │
         │  /api/inbox/*    │
         │                  │
         │ • CRUD operations│
         │ • Task assignment│
         │ • Approvals      │
         └────────┬─────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌───────────────┐   ┌───────────────┐
│  PostgreSQL   │   │   WebSocket   │
│  inbox_items  │   │  /api/inbox/ws│
│               │   │               │
│ • Persistence │   │ • Real-time   │
│ • History     │   │ • Broadcast   │
└───────────────┘   └───────┬───────┘
                            │
                            ▼
                    ┌───────────────┐
                    │  Flutter UI   │
                    │  Inbox Screen │
                    │               │
                    │ • Live updates│
                    │ • Actions     │
                    └───────────────┘
```

---

## Quick Start

### Send a Test Notification

```bash
# Create notification (appears instantly in Flutter inbox)
curl -X POST http://100.102.21.44:8081/api/inbox \
  -H "Content-Type: application/json" \
  -d '{
    "type": "notification",
    "title": "Test Notification",
    "body": "This appears in real-time!",
    "priority": "normal",
    "user_id": "zeus"
  }'
```

### Create a Task for Legacy Agent

```bash
curl -X POST http://100.102.21.44:8081/api/inbox \
  -H "Content-Type: application/json" \
  -d '{
    "type": "task",
    "title": "Research competitor pricing",
    "body": "Look at pricing pages for top 5 competitors",
    "priority": "high",
    "assignee_type": "agent",
    "user_id": "zeus"
  }'
```

### Request Approval (HITL)

```bash
curl -X POST http://100.102.21.44:8081/api/inbox \
  -H "Content-Type: application/json" \
  -d '{
    "type": "approval",
    "title": "Delete old log files",
    "body": "Agent wants to delete 50 log files older than 30 days",
    "priority": "high",
    "requires_approval": true,
    "user_id": "zeus"
  }'
```

---

## API Reference

### Base URL
- **Local:** `http://localhost:8081/api/inbox`
- **Tailscale:** `http://100.102.21.44:8081/api/inbox`
- **Production:** `https://api.onemindos.dev/api/inbox`

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/inbox` | List inbox items (with filters) |
| `GET` | `/api/inbox/stats` | Get inbox statistics |
| `GET` | `/api/inbox/approvals` | List pending approvals |
| `GET` | `/api/inbox/tasks` | List tasks |
| `GET` | `/api/inbox/agent-tasks` | List tasks assigned to agent |
| `GET` | `/api/inbox/novu-config` | Get Novu config for Flutter |
| `GET` | `/api/inbox/{id}` | Get single item |
| `GET` | `/api/inbox/{id}/history` | Get item history |
| `POST` | `/api/inbox` | Create new item |
| `PATCH` | `/api/inbox/{id}` | Update item |
| `POST` | `/api/inbox/{id}/read` | Mark as read |
| `POST` | `/api/inbox/read-all` | Mark all as read |
| `POST` | `/api/inbox/{id}/star` | Star item |
| `POST` | `/api/inbox/{id}/unstar` | Unstar item |
| `POST` | `/api/inbox/{id}/archive` | Archive item |
| `POST` | `/api/inbox/{id}/snooze` | Snooze until time |
| `POST` | `/api/inbox/{id}/assign` | Assign to human/agent |
| `POST` | `/api/inbox/{id}/approve` | Approve request |
| `POST` | `/api/inbox/{id}/deny` | Deny request |
| `POST` | `/api/inbox/{id}/complete` | Complete task |
| `POST` | `/api/inbox/{id}/fail` | Fail task |
| `DELETE` | `/api/inbox/{id}` | Delete (archive) item |
| `WS` | `/api/inbox/ws` | WebSocket for real-time |

### Query Parameters (List Endpoint)

| Parameter | Type | Description |
|-----------|------|-------------|
| `status` | string | Comma-separated: `new,reading,approved,denied,completed,failed,archived` |
| `type` | string | Filter: `notification`, `task`, `approval`, `message`, `alert` |
| `priority` | string | Filter: `critical`, `high`, `normal`, `low` |
| `assignee` | string | Filter: `human`, `agent`, `unassigned` |
| `is_read` | bool | Filter by read status |
| `is_starred` | bool | Filter by starred |
| `limit` | int | Max items (default: 50, max: 200) |
| `offset` | int | Pagination offset |
| `user_id` | string | User ID (default: "zeus") |

### Create Item Schema

```json
{
  "type": "notification | task | approval | message | alert",
  "title": "Item title",
  "body": "Item description/content",
  "priority": "critical | high | normal | low",
  "source": "legacy | todoist | github | etc",
  "user_id": "zeus",
  "assignee_type": "human | agent | null",
  "assignee_id": "optional specific assignee",
  "requires_approval": false,
  "metadata": {},
  "action_url": "/optional/deep/link",
  "task_due_at": "2026-01-24T17:00:00Z"
}
```

---

## WebSocket Real-Time Updates

### Connection

```javascript
// Connect to WebSocket
const ws = new WebSocket('ws://100.102.21.44:8081/api/inbox/ws?user_id=zeus');

ws.onopen = () => {
  console.log('Connected to inbox');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);

  switch (message.type) {
    case 'inbox.stats.updated':
      // Initial stats on connect
      updateStats(message.data);
      break;
    case 'inbox.item.created':
      // New item - add to list
      addItem(message.data);
      break;
    case 'inbox.item.updated':
      // Item changed - update in list
      updateItem(message.data);
      break;
  }
};

// Keep-alive ping
setInterval(() => {
  ws.send(JSON.stringify({ type: 'ping' }));
}, 30000);
```

### Event Types

| Event | When | Data |
|-------|------|------|
| `inbox.stats.updated` | On connect | Stats object |
| `inbox.item.created` | New item | Full InboxItem |
| `inbox.item.updated` | Any change | Full InboxItem |
| `pong` | Response to ping | `{}` |

### Flutter Provider

The Flutter app uses `UnifiedInboxNotifier` in:
- `ui/lib/features/inbox/providers/unified_inbox_provider.dart`

```dart
// Provider auto-connects WebSocket and handles events
final inboxProvider = ref.watch(unifiedInboxProvider);

// Access items
final items = inboxProvider.items;
final stats = inboxProvider.stats;
final isConnected = inboxProvider.isConnected;
```

---

## Novu Cloud Integration

### What Novu Handles

| Feature | Novu | Backend |
|---------|------|---------|
| Push notifications (iOS/Android) | ✅ | Trigger via API |
| In-app notifications | ✅ Inbox component | Event feed |
| Email notifications | ✅ | Configure providers |
| SMS notifications | ✅ | Configure providers |
| Subscriber preferences | ✅ UI | Sync to backend |
| Workflow orchestration | ✅ | Define in dashboard |

### Environment Variables

```bash
# Required
NOVU_API_KEY=your_api_key          # From https://web.novu.co/settings
NOVU_APP_ID=your_app_id            # From https://web.novu.co/settings

# Optional (defaults to NOVU_API_KEY)
NOVU_SECRET_KEY=your_secret_key    # For subscriber hash (HMAC)
```

### Setup Novu Workflows

Run the setup script to create required workflows:

```bash
# Set API key
export NOVU_API_KEY="your-api-key"

# Run setup
python scripts/novu/setup_workflows.py

# Or verify existing setup
python scripts/novu/setup_workflows.py --verify-only
```

This creates:
- `legacy-notification` - Standard notifications
- `legacy-approval` - HITL approval requests

### Novu Config Endpoint

Flutter fetches Novu config from:

```bash
curl http://100.102.21.44:8081/api/inbox/novu-config?user_id=zeus

# Returns:
{
  "app_id": "PFlDST_penFf",
  "subscriber_id": "zeus",
  "subscriber_hash": "hmac_hash_for_security"
}
```

---

## Database Schema

### inbox_items Table

```sql
CREATE TABLE inbox_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(100) NOT NULL DEFAULT 'zeus',

    -- Content
    type VARCHAR(50) NOT NULL,           -- notification, task, approval, message, alert
    title VARCHAR(500) NOT NULL,
    body TEXT,
    source VARCHAR(100) DEFAULT 'legacy',

    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'new',  -- new, reading, approved, denied, completed, failed, archived
    priority VARCHAR(20) DEFAULT 'normal',      -- critical, high, normal, low

    -- Assignment
    assignee_type VARCHAR(20),           -- human, agent, null
    assignee_id VARCHAR(100),

    -- Flags
    is_read BOOLEAN DEFAULT FALSE,
    is_starred BOOLEAN DEFAULT FALSE,
    requires_approval BOOLEAN DEFAULT FALSE,

    -- Metadata
    metadata JSONB DEFAULT '{}',
    action_url VARCHAR(1000),

    -- Task fields
    task_due_at TIMESTAMPTZ,
    task_completed_at TIMESTAMPTZ,
    task_result TEXT,

    -- Approval fields
    approved_at TIMESTAMPTZ,
    approved_by VARCHAR(100),
    denied_at TIMESTAMPTZ,
    denied_by VARCHAR(100),
    denial_reason TEXT,

    -- Novu integration
    novu_notification_id VARCHAR(100),
    novu_transaction_id VARCHAR(100),

    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    snoozed_until TIMESTAMPTZ,

    -- Indexes
    CONSTRAINT inbox_items_type_check CHECK (type IN ('notification', 'task', 'approval', 'message', 'alert')),
    CONSTRAINT inbox_items_status_check CHECK (status IN ('new', 'reading', 'approved', 'denied', 'completed', 'failed', 'archived'))
);

-- Indexes for common queries
CREATE INDEX idx_inbox_items_user_status ON inbox_items(user_id, status);
CREATE INDEX idx_inbox_items_user_type ON inbox_items(user_id, type);
CREATE INDEX idx_inbox_items_created_at ON inbox_items(created_at DESC);
CREATE INDEX idx_inbox_items_priority ON inbox_items(priority);
CREATE INDEX idx_inbox_items_assignee ON inbox_items(assignee_type, assignee_id);
```

### inbox_history Table

```sql
CREATE TABLE inbox_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES inbox_items(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    actor VARCHAR(100),
    old_value JSONB,
    new_value JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_inbox_history_item_id ON inbox_history(item_id);
```

---

## Agent Task Assignment Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    AGENT TASK ASSIGNMENT FLOW                                │
└─────────────────────────────────────────────────────────────────────────────┘

1. User creates task in Flutter Inbox
        │
        ▼
   POST /api/inbox
   {
     "type": "task",
     "title": "Research competitor pricing",
     "assignee_type": "agent"
   }
        │
        ▼
2. Task stored in PostgreSQL + WebSocket broadcast
        │
        ▼
3. Legacy Agent pre-hook: inject_agent_tasks()
   - Checks inbox for tasks assigned to agent
   - Injects into agent context
        │
        ▼
4. Agent sees in context:
   "[PENDING AGENT TASKS]
    You have 1 task(s) assigned to you:
    - [High] Research competitor pricing (ID: abc-123)"
        │
        ▼
5. Agent works on task, uses tools
        │
        ▼
6. Agent completes task:
   complete_inbox_task(task_id="abc-123", result="Found 5 competitors...")
        │
        ▼
7. Task marked completed + WebSocket broadcast
        │
        ▼
8. User sees completion in Flutter inbox (real-time)
```

### Agent Tools for Tasks

The agent has these inbox tools available:

```python
# List tasks assigned to agent
list_my_tasks(limit: int = 5) -> List[InboxItem]

# Complete a task
complete_inbox_task(task_id: str, result: str) -> InboxItem

# Fail a task
fail_inbox_task(task_id: str, reason: str) -> InboxItem

# Create inbox item (for notifications to user)
create_inbox_item(
    type: str,
    title: str,
    body: str,
    priority: str = "normal"
) -> InboxItem
```

---

## File Structure

```
core/inbox/
├── __init__.py           # Package init
├── api.py                # FastAPI router + WebSocket ⭐
├── models.py             # Pydantic models (InboxItem, etc.)
├── service.py            # Business logic (InboxService)
└── novu.py               # Novu client wrapper

core/tools/
└── inbox.py              # Agent tools (complete_task, fail_task, etc.)

core/agent/
└── hooks.py              # inject_agent_tasks() pre-hook

ui/lib/features/inbox/
├── screens/
│   └── inbox_screen.dart # Flutter inbox UI
├── providers/
│   └── unified_inbox_provider.dart  # Riverpod + WebSocket ⭐
└── widgets/
    └── inbox_item_card.dart

scripts/novu/
└── setup_workflows.py    # Novu workflow setup script
```

---

## Maintenance

### Check System Health

```bash
# Check inbox API
curl http://100.102.21.44:8081/api/inbox/stats

# Check WebSocket
websocat ws://100.102.21.44:8081/api/inbox/ws?user_id=zeus

# Check Novu config
curl http://100.102.21.44:8081/api/inbox/novu-config?user_id=zeus

# Check container logs
ssh ubuntu@100.102.21.44
docker logs onemind-agno -f --tail 100 | grep -i inbox
```

### Common Issues

| Issue | Cause | Fix |
|-------|-------|-----|
| WebSocket not connecting | Container restart | Refresh Flutter app |
| Novu config empty | Missing NOVU_SECRET_KEY | Add to env, restart container |
| Items not broadcasting | Old code deployed | Pull latest, rebuild image |
| Tasks not showing for agent | Wrong awareness mode | Set to ATTENTIVE or higher |

### Rebuild & Deploy

```bash
# SSH to server
ssh zeus@20.121.38.186

# Pull latest code
cd ~/onemind-os && git pull

# Rebuild image
docker build -f infra/containers/agno/Dockerfile -t onemind-os-agno:latest .

# Restart container (preserving env vars)
docker stop onemind-agno && docker rm onemind-agno

docker run -d \
  --name onemind-agno \
  --network onemind-os_onemind \
  -p 8081:8080 \
  -e DATABASE_URL=postgresql://onemind:PASSWORD@onemind-postgres:5432/onemind \
  -e REDIS_URL=redis://onemind-redis:6379/0 \
  -e NATS_URL=nats://onemind-nats:4222 \
  -e AWS_ACCESS_KEY_ID=YOUR_KEY \
  -e AWS_SECRET_ACCESS_KEY=YOUR_SECRET \
  -e AWS_DEFAULT_REGION=us-east-1 \
  -e NOVU_API_KEY=YOUR_NOVU_KEY \
  -e NOVU_APP_ID=YOUR_NOVU_APP_ID \
  -e NOVU_SECRET_KEY=YOUR_NOVU_KEY \
  --restart unless-stopped \
  onemind-os-agno:latest
```

### Verify Deployment

```bash
# Test notification (should appear instantly in Flutter)
curl -X POST http://100.102.21.44:8081/api/inbox \
  -H "Content-Type: application/json" \
  -d '{"type": "notification", "title": "Deploy Test", "body": "System working!", "user_id": "zeus"}'
```

---

## Awareness Mode Integration

The inbox respects the **awareness mode** to determine:
1. **What gets through** - Filtering based on priority thresholds
2. **How to notify** - Channels enabled per mode
3. **Who handles** - Auto-assignment logic

### Awareness Mode → Inbox Behavior

| Mode | Priority Threshold | Auto-Assign to Agent | Notifications |
|------|-------------------|---------------------|---------------|
| `dormant` | Critical only | No | Critical alerts only |
| `aware` | High + Critical | No | Critical + High + Approvals |
| `present` | All except Low | Yes (suggestions, low events) | All except Low |
| `omnipresent` | Everything | Yes (everything) | All |

---

## Verification Checklist

### Backend ✅
- [x] Database schema created
- [x] CRUD API working
- [x] WebSocket broadcasting
- [x] Novu notifications integration
- [x] Awareness mode integration
- [x] Agent task assignment
- [x] History/audit trail

### Flutter ✅
- [x] Freezed models generated
- [x] Real-time WebSocket updates
- [x] Filters working
- [x] Mark read/archive working
- [x] HITL approve/deny dialogs
- [x] Task complete/fail dialogs
- [x] Assign to human/agent
- [x] Snooze options

### Integration (To Test)
- [ ] GitHub webhook → Inbox
- [ ] Todoist webhook → Inbox
- [ ] Agent approval → Inbox → User decision → Agent continues
- [ ] Task assignment → Agent picks up → Completes → Notification

---

## Future Enhancements

| Feature | Status | Notes |
|---------|--------|-------|
| Real-time WebSocket | ✅ Done | Broadcasting on all actions |
| Novu push notifications | 🔜 Ready | Needs FCM/APNs provider setup |
| Email notifications | 🔜 Ready | Needs email provider in Novu |
| SMS notifications | 🔜 Ready | Needs Twilio in Novu |
| Approval workflows | ✅ Done | HITL in place |
| Agent task assignment | ✅ Done | Pre-hook injects tasks |
| Batch operations | 🔜 Planned | Mark multiple as read |
| Snooze reminders | 🔜 Planned | Background job for un-snooze |

---

## Related Documentation

- [Awareness System Plan](/.claude/plans/greedy-brewing-hollerith.md) - Full architecture
- [Notification Router](core/orchestrator/notifications.py) - Novu integration
- [Inbox Service](core/inbox/service.py) - Business logic
- [Agent Hooks](core/agent/hooks.py) - Task injection

---

## Changelog

### 2026-01-23
- ✅ Fixed WebSocket broadcasting (was never calling broadcast)
- ✅ Added `broadcast_item_update()` helper
- ✅ All action endpoints now broadcast changes
- ✅ Added NOVU_SECRET_KEY for subscriber hash
- ✅ Real-time notifications confirmed working
