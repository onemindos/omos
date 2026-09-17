# OneMind OS - Webhook System

> Last Updated: 2026-01-26
> Status: Production

---

## Overview

OneMind OS uses a custom webhook receiver for external service integrations. This replaces the previous Nango-based approach with a simpler, more direct solution.

**Supported Services:**
- Todoist (task management)
- GitHub (code repositories)
- Linear (issue tracking)

**Key Features:**
- Proper signature verification for each service
- Event publishing to NATS JetStream
- Automatic inbox item creation for important events
- Health check endpoint

---

## Architecture

```
External Service                    OneMind OS
     │                                  │
     │   Webhook Event                  │
     └──────────────────────────────────►
                                        │
                               ┌────────▼────────┐
                               │  Webhook API    │
                               │ /api/webhooks/* │
                               └────────┬────────┘
                                        │
                          ┌─────────────┼─────────────┐
                          │             │             │
                   ┌──────▼──────┐ ┌────▼────┐ ┌──────▼──────┐
                   │   Todoist   │ │  GitHub │ │   Linear    │
                   │   Handler   │ │ Handler │ │   Handler   │
                   └──────┬──────┘ └────┬────┘ └──────┬──────┘
                          │             │             │
                          └─────────────┼─────────────┘
                                        │
                               ┌────────▼────────┐
                               │  NATS JetStream │
                               │  events.*.*    │
                               └────────┬────────┘
                                        │
                               ┌────────▼────────┐
                               │  Inbox Service  │
                               │ (if important)  │
                               └─────────────────┘
```

---

## Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/api/webhooks/todoist` | Receive Todoist events |
| POST | `/api/webhooks/github` | Receive GitHub events |
| POST | `/api/webhooks/linear` | Receive Linear events |
| GET | `/api/webhooks/health` | Health check with status |
| GET | `/api/webhooks/test` | Verify router is mounted |

---

## Signature Verification

Each service uses different signature verification methods:

### Todoist
- **Header:** `X-Todoist-Hmac-SHA256`
- **Format:** Base64-encoded HMAC-SHA256
- **Secret:** `TODOIST_WEBHOOK_SECRET`

### GitHub
- **Header:** `X-Hub-Signature-256`
- **Format:** `sha256=` + hex-encoded HMAC-SHA256
- **Event Type Header:** `X-GitHub-Event`
- **Secret:** `GITHUB_WEBHOOK_SECRET`

### Linear
- **Header:** `Linear-Signature`
- **Format:** `t=timestamp,s=signature`
- **Event Type Header:** `Linear-Event`
- **Secret:** `LINEAR_WEBHOOK_SECRET`

---

## NATS Event Subjects

Events are published to NATS JetStream with these subject patterns:

```
events.todoist.item.added
events.todoist.item.updated
events.todoist.item.completed
events.todoist.item.deleted
events.github.push
events.github.pull_request.opened
events.github.pull_request.closed
events.github.issues.opened
events.github.issues.assigned
events.linear.Issue.created
events.linear.Issue.updated
events.linear.Comment.created
```

---

## Environment Variables

```bash
# Required for signature verification
TODOIST_WEBHOOK_SECRET=your_todoist_webhook_secret
GITHUB_WEBHOOK_SECRET=your_github_webhook_secret
LINEAR_WEBHOOK_SECRET=your_linear_webhook_secret

# NATS connection (optional, defaults to localhost)
NATS_URL=nats://localhost:4222
```

---

## Setting Up Webhooks

### Todoist

1. Go to [Todoist Developer Console](https://developer.todoist.com/appconsole.html)
2. Create or select your app
3. Go to **Webhooks** section
4. Add webhook URL: `https://api.onemindos.dev/api/webhooks/todoist`
5. Copy the **Webhook Secret** → set as `TODOIST_WEBHOOK_SECRET`
6. Select events to receive (item:added, item:completed, etc.)

### GitHub

1. Go to your repository → **Settings** → **Webhooks**
2. Click **Add webhook**
3. Payload URL: `https://api.onemindos.dev/api/webhooks/github`
4. Content type: `application/json`
5. Secret: Generate a secure secret → set as `GITHUB_WEBHOOK_SECRET`
6. Select events: Issues, Pull requests, Push

### Linear

1. Go to **Settings** → **API** → **Webhooks**
2. Click **New webhook**
3. URL: `https://api.onemindos.dev/api/webhooks/linear`
4. Select events: Issue, Comment, Project
5. Copy the **Signing Secret** → set as `LINEAR_WEBHOOK_SECRET`

---

## Testing

### Health Check
```bash
curl https://api.onemindos.dev/api/webhooks/health
```

Response:
```json
{
  "status": "healthy",
  "nats_connected": true,
  "handlers": ["todoist", "github", "linear"],
  "secrets_configured": {
    "todoist": true,
    "github": true,
    "linear": false
  }
}
```

### Test Endpoint
```bash
curl https://api.onemindos.dev/api/webhooks/test
```

### Monitor NATS Events
```bash
nats sub "events.>" --server=nats://localhost:4222
```

---

## Inbox Integration

The webhook handlers automatically create inbox items for important events:

| Service | Event | Creates Inbox Item? |
|---------|-------|---------------------|
| Todoist | item:added | Yes (all new tasks) |
| Todoist | item:completed | No |
| GitHub | pull_request.review_requested | Yes |
| GitHub | issues.assigned | Yes |
| GitHub | issues.opened | Yes |
| GitHub | push | No |
| Linear | Issue.created | Yes (if assigned or high priority) |
| Linear | Comment.created | No |

---

## File Structure

```
core/webhooks/
├── __init__.py              # Module exports
├── models.py                # WebhookEvent, WebhookResponse, WebhookSource
├── base.py                  # Signature verification utilities
├── api.py                   # FastAPI router
└── handlers/
    ├── __init__.py          # Handler exports
    ├── todoist.py           # Todoist webhook handler
    ├── github.py            # GitHub webhook handler
    └── linear.py            # Linear webhook handler
```

---

## Related Documentation

- [Composio Integration](../integrations/COMPOSIO_INTEGRATION.md) - OAuth for actions
- [Unified Inbox System](../UNIFIED-INBOX-SYSTEM.md) - Inbox item creation
- [NATS Architecture](CONSCIOUSNESS-ARCHITECTURE.md) - Event bus details
