# n8n Integration Bridge

n8n is a self-hosted workflow automation tool that can connect 400+ services to OneMind OS.

## Why n8n?

- **Self-hosted** - No data leaves your infrastructure
- **400+ integrations** - Todoist, Notion, Airtable, ClickUp, Spotify, Fitbit, etc.
- **Visual workflow builder** - Easy to configure
- **Free & open source** - No subscription fees

## Quick Setup

```bash
# Docker Compose (add to your docker-compose.yml)
n8n:
  image: n8nio/n8n
  ports:
    - "5678:5678"
  environment:
    - N8N_HOST=n8n.your-domain.com
    - N8N_PORT=5678
    - N8N_PROTOCOL=https
    - WEBHOOK_URL=https://n8n.your-domain.com/
  volumes:
    - n8n_data:/home/node/.n8n
```

## Integration Pattern

### Option A: Forward to Webhook Receiver
```
[Todoist] → [n8n Trigger] → [HTTP Request Node] → POST /webhooks/todoist
                                                   {
                                                     "event_type": "task_completed",
                                                     "task": {...}
                                                   }
```

### Option B: Trigger Agent Directly
```
[ClickUp Task] → [n8n Trigger] → [HTTP Request Node] → POST /webhooks/trigger/clickup-agent
                                                        {
                                                          "message": "New task: {{task.name}}",
                                                          "context": {{$json}}
                                                        }
```

## Example Workflows

### 1. Todoist → OneMind Activity Feed
```json
{
  "nodes": [
    {
      "name": "Todoist Trigger",
      "type": "n8n-nodes-base.todoistTrigger",
      "parameters": {
        "event": "taskCompleted"
      }
    },
    {
      "name": "Forward to OneMind",
      "type": "n8n-nodes-base.httpRequest",
      "parameters": {
        "method": "POST",
        "url": "https://api.onemind.local/webhooks/todoist",
        "body": {
          "event_type": "task_completed",
          "task_id": "={{$json.id}}",
          "content": "={{$json.content}}",
          "project": "={{$json.project_id}}"
        }
      }
    }
  ]
}
```

### 2. Notion Database Change → Agent
```json
{
  "trigger": "Notion Database Updated",
  "action": "HTTP POST /webhooks/trigger/researcher",
  "body": {
    "message": "Notion page updated: {{page.title}}",
    "context": {
      "page_id": "{{page.id}}",
      "properties": "{{page.properties}}"
    }
  }
}
```

### 3. Spotify Now Playing → Awareness Context
```json
{
  "trigger": "Every 5 minutes",
  "check": "Spotify Currently Playing",
  "action": "POST /webhooks/spotify",
  "body": {
    "event_type": "now_playing",
    "track": "{{track.name}}",
    "artist": "{{track.artist}}",
    "playing": true
  }
}
```

## Webhook URL Format

All n8n workflows should POST to one of these endpoints:

| Endpoint | Purpose |
|----------|---------|
| `POST /webhooks/{source}` | Generic events (goes to NATS → Activity Feed) |
| `POST /webhooks/trigger/{agent_id}` | Direct agent trigger (immediate response) |
| `POST /webhooks/trigger/{agent_id}/async` | Background agent trigger (returns immediately) |

## Event Format

Standard event format for `/webhooks/{source}`:
```json
{
  "event_type": "task_completed",  // Will be used as NATS subject suffix
  "timestamp": "2025-02-03T12:00:00Z",
  "source": "todoist",
  "data": {
    // Service-specific data
  }
}
```

## Alternative: Windmill

If you prefer code-first automation:
- **Windmill** (https://windmill.dev) - TypeScript/Python workflows
- Similar self-hosted, connects to same webhooks


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
