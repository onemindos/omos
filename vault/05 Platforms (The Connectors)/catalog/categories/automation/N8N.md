# N8N

> **Workflow Automation Platform**

## Overview

N8N is a self-hosted workflow automation platform for connecting services, automating tasks, and building integrations without code.

## Status

- **Active**: ✅ Yes
- **Deployment**: Self-hosted (Docker)
- **Last Updated**: January 25, 2026

## Use Cases

1. **Business Automation** - Automate repetitive tasks
2. **Integrations** - Connect services together
3. **Webhooks** - Receive and process webhooks
4. **Notifications** - Multi-channel alerts
5. **Data Sync** - Keep systems in sync

## Key Features

- Visual workflow builder
- 400+ integrations
- Self-hosted (data privacy)
- JavaScript code nodes
- Error handling & retry
- Credential management
- Version control via git

## Installation

### Docker

```bash
docker run -d \
  --name n8n \
  -p 5678:5678 \
  -v n8n_data:/home/node/.n8n \
  n8nio/n8n
```

### Docker Compose

```yaml
services:
  n8n:
    image: n8nio/n8n
    ports:
      - "5678:5678"
    volumes:
      - n8n_data:/home/node/.n8n
    environment:
      - N8N_BASIC_AUTH_ACTIVE=true
      - N8N_BASIC_AUTH_USER=admin
      - N8N_BASIC_AUTH_PASSWORD=${N8N_PASSWORD}
      - WEBHOOK_URL=https://n8n.onemind.io
      - N8N_ENCRYPTION_KEY=${ENCRYPTION_KEY}
    restart: unless-stopped

volumes:
  n8n_data:
```

## Workflow Structure

### Basic Workflow

```json
{
  "nodes": [
    {
      "name": "Webhook",
      "type": "n8n-nodes-base.webhook",
      "position": [250, 300],
      "parameters": {
        "httpMethod": "POST",
        "path": "process-event"
      }
    },
    {
      "name": "Process Data",
      "type": "n8n-nodes-base.code",
      "position": [450, 300],
      "parameters": {
        "jsCode": "return items.map(item => ({json: {processed: true, ...item.json}}));"
      }
    },
    {
      "name": "Send Notification",
      "type": "n8n-nodes-base.slack",
      "position": [650, 300],
      "parameters": {
        "channel": "#notifications",
        "text": "Event processed: {{ $json.event_type }}"
      }
    }
  ],
  "connections": {
    "Webhook": {
      "main": [[{ "node": "Process Data", "type": "main", "index": 0 }]]
    },
    "Process Data": {
      "main": [[{ "node": "Send Notification", "type": "main", "index": 0 }]]
    }
  }
}
```

## Common Triggers

### Webhook

```json
{
  "type": "n8n-nodes-base.webhook",
  "parameters": {
    "httpMethod": "POST",
    "path": "my-webhook",
    "responseMode": "responseNode"
  }
}
```

### Schedule

```json
{
  "type": "n8n-nodes-base.scheduleTrigger",
  "parameters": {
    "rule": {
      "interval": [{ "field": "hours", "value": 1 }]
    }
  }
}
```

### Cron

```json
{
  "type": "n8n-nodes-base.cron",
  "parameters": {
    "cronExpression": "0 9 * * 1-5"
  }
}
```

## Common Nodes

### HTTP Request

```json
{
  "type": "n8n-nodes-base.httpRequest",
  "parameters": {
    "method": "POST",
    "url": "https://api.example.com/endpoint",
    "headers": {
      "Authorization": "Bearer {{ $credentials.apiKey }}"
    },
    "body": {
      "data": "{{ $json.data }}"
    }
  }
}
```

### Code Node

```javascript
// Code node for custom logic
const items = $input.all();

const processed = items.map((item) => {
  return {
    json: {
      ...item.json,
      processed: true,
      timestamp: new Date().toISOString(),
    },
  };
});

return processed;
```

### IF Condition

```json
{
  "type": "n8n-nodes-base.if",
  "parameters": {
    "conditions": {
      "string": [
        {
          "value1": "{{ $json.status }}",
          "operation": "equals",
          "value2": "active"
        }
      ]
    }
  }
}
```

## Integrations

### Commonly Used

| Service            | Use Case         |
| ------------------ | ---------------- |
| **Slack**          | Notifications    |
| **Discord**        | Alerts           |
| **Google Sheets**  | Data storage     |
| **Postgres**       | Database         |
| **Todoist**        | Task creation    |
| **GitHub**         | Issue management |
| **Home Assistant** | Smart home       |
| **Zoho CRM**       | Lead management  |

## Error Handling

### Error Trigger

```json
{
  "type": "n8n-nodes-base.errorTrigger",
  "position": [250, 500]
}
```

### Retry Logic

- Configure retries on individual nodes
- Set max attempts and wait between retries
- Use error workflow for notifications

## Legacy AI Integration

### Trigger Agent via Webhook

```
N8N Workflow:
1. Webhook receives event
2. Format message for agent
3. HTTP Request to Legacy API
4. Process response
5. Take action (create task, send notification)
```

### Example: Email to Agent

```json
{
  "nodes": [
    { "name": "Email Trigger", "type": "n8n-nodes-base.emailTrigger" },
    { "name": "Extract Content", "type": "n8n-nodes-base.code" },
    { "name": "Call Legacy", "type": "n8n-nodes-base.httpRequest" },
    { "name": "Create Todoist Task", "type": "n8n-nodes-base.todoist" }
  ]
}
```

## API Access

### Trigger Workflow

```python
import requests

def trigger_n8n_workflow(workflow_id: str, data: dict):
    response = requests.post(
        f"https://n8n.onemind.io/webhook/{workflow_id}",
        json=data,
        headers={"Content-Type": "application/json"}
    )
    return response.json()
```

### Webhook URL Format

```
https://n8n.onemind.io/webhook/<workflow-id>
https://n8n.onemind.io/webhook-test/<workflow-id>  # Test mode
```

## Active Workflows

| Workflow           | Trigger        | Actions                              |
| ------------------ | -------------- | ------------------------------------ |
| **Daily Summary**  | Cron (9am)     | Gather data → Format → Send to Slack |
| **Lead Alert**     | Zoho webhook   | Process → Notify → Create task       |
| **GitHub Monitor** | GitHub webhook | Check PR → Notify if important       |
| **Error Handler**  | Error trigger  | Log → Alert → Create issue           |

## Environment Variables

```env
N8N_BASIC_AUTH_ACTIVE=true
N8N_BASIC_AUTH_USER=admin
N8N_BASIC_AUTH_PASSWORD=secure_password
WEBHOOK_URL=https://n8n.onemind.io
N8N_ENCRYPTION_KEY=your_encryption_key
N8N_LOG_LEVEL=info
```

## Backup & Restore

### Export Workflows

```bash
# Via CLI
n8n export:workflow --all --output=workflows.json

# Or use the UI: Settings → Download
```

### Import Workflows

```bash
n8n import:workflow --input=workflows.json
```

## Best Practices

- ✅ Use descriptive workflow names
- ✅ Add notes to explain complex logic
- ✅ Use sub-workflows for reusable logic
- ✅ Implement error handling
- ✅ Use credentials for secrets
- ✅ Test with webhook-test URLs first
- ✅ Version control workflow exports
- ✅ Set up monitoring for failed workflows

## Resources

- [N8N Documentation](https://docs.n8n.io/)
- [Node Library](https://n8n.io/integrations)
- [Community](https://community.n8n.io/)

---

**Related**: [Kestra](Kestra.md) | [Home-Assistant](./Home-Assistant.md) | [Todoist](./Todoist.md)
