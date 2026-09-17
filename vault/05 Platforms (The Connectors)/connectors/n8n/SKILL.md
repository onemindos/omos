---
name: n8n
description: Manage N8N workflow automation. Use when you need to trigger workflows, check workflow status, list executions, or manage automation pipelines.
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Read
  - Write
---

# N8N Workflow Automation Skill

Interact with N8N workflow automation platform.

## Environment Variables

```bash
N8N_URL        # Base URL (e.g., https://n8n.yourdomain.com)
N8N_API_KEY    # API key from N8N settings
```

## Authentication

```bash
X-N8N-API-KEY: $N8N_API_KEY
```

---

## Operations

### Workflows

**List all workflows:**
```bash
curl -s "$N8N_URL/api/v1/workflows" \
  -H "X-N8N-API-KEY: $N8N_API_KEY" | jq '.data[] | {id, name, active}'
```

**Get workflow:**
```bash
curl -s "$N8N_URL/api/v1/workflows/{workflow_id}" \
  -H "X-N8N-API-KEY: $N8N_API_KEY"
```

**Activate workflow:**
```bash
curl -s -X PATCH "$N8N_URL/api/v1/workflows/{workflow_id}" \
  -H "X-N8N-API-KEY: $N8N_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"active": true}'
```

**Deactivate workflow:**
```bash
curl -s -X PATCH "$N8N_URL/api/v1/workflows/{workflow_id}" \
  -H "X-N8N-API-KEY: $N8N_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"active": false}'
```

### Executions

**List executions:**
```bash
curl -s "$N8N_URL/api/v1/executions" \
  -H "X-N8N-API-KEY: $N8N_API_KEY" | jq '.data[] | {id, workflowId, status, startedAt}'
```

**Get execution:**
```bash
curl -s "$N8N_URL/api/v1/executions/{execution_id}" \
  -H "X-N8N-API-KEY: $N8N_API_KEY"
```

**Get execution with data:**
```bash
curl -s "$N8N_URL/api/v1/executions/{execution_id}?includeData=true" \
  -H "X-N8N-API-KEY: $N8N_API_KEY"
```

### Trigger Workflow via Webhook

**Trigger with GET:**
```bash
curl -s "$N8N_URL/webhook/{webhook_path}"
```

**Trigger with POST data:**
```bash
curl -s -X POST "$N8N_URL/webhook/{webhook_path}" \
  -H "Content-Type: application/json" \
  -d '{
    "task": "process_data",
    "data": {"key": "value"}
  }'
```

**Trigger test webhook:**
```bash
curl -s -X POST "$N8N_URL/webhook-test/{webhook_path}" \
  -H "Content-Type: application/json" \
  -d '{"test": true}'
```

### Credentials

**List credentials:**
```bash
curl -s "$N8N_URL/api/v1/credentials" \
  -H "X-N8N-API-KEY: $N8N_API_KEY" | jq '.data[] | {id, name, type}'
```

---

## One Mind Workflows

| Workflow | Purpose | Trigger |
|----------|---------|---------|
| Morning Briefing | Daily summary to Telegram | Schedule 7am |
| Task Sync | Todoist ↔ Vault sync | Every 15min |
| Approval Handler | Process Discord approvals | Webhook |
| Memory Consolidation | Daily memory archival | Schedule 11pm |
| Health Check | Monitor all services | Every 5min |

## Webhook Patterns

**Production webhook:**
```
$N8N_URL/webhook/one-mind-{workflow-name}
```

**Test webhook:**
```
$N8N_URL/webhook-test/one-mind-{workflow-name}
```

---

## Execution Status

| Status | Meaning |
|--------|---------|
| `success` | Completed successfully |
| `error` | Failed with error |
| `waiting` | Waiting for external trigger |
| `running` | Currently executing |

## Tips

- Use webhook triggers for real-time automation
- Schedule triggers for recurring tasks
- Check executions to debug failed workflows
- Test webhooks before activating production
- N8N stores execution data for debugging


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
