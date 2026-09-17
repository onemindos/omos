---
id: ui-12-20260111053926
title: "N8N Integration Configuration"
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
ai_summary: "N8N Integration Configuration"
---

# N8N Integration Configuration

## Overview

N8N serves as the **Automation Orchestration Layer** for One Mind. It connects all systems and enables agents to trigger workflows without direct API access.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                          N8N HUB                                │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   Webhooks   │  │   Triggers   │  │   Schedules  │         │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘         │
│         │                 │                  │                  │
│         └─────────────────┼──────────────────┘                  │
│                           │                                     │
│                  ┌────────▼────────┐                           │
│                  │  Workflow Engine │                           │
│                  └────────┬────────┘                           │
│                           │                                     │
│    ┌──────────────────────┼──────────────────────┐             │
│    │                      │                      │             │
│ ┌──▼───┐  ┌──────▼──────┐  ┌──▼───┐  ┌──▼───┐  │
│ │Todoist│  │    Letta    │  │GitHub│  │Obsidian│              │
│ └──────┘  └─────────────┘  └──────┘  └───────┘               │
└─────────────────────────────────────────────────────────────────┘
```

## Core Workflows

### 1. Morning Brief Generation

```yaml
name: morning-brief
trigger:
  type: cron
  schedule: "0 6 * * *"  # 6 AM daily
steps:
  1. Fetch today's calendar events
  2. Get pending Todoist tasks
  3. Check system health (Grid agent)
  4. Query weather API
  5. Send context to LEGACY agent
  6. Generate morning brief
  7. Push to Obsidian daily note
  8. Send notification to Zeus
```

### 2. Task Sync Loop

```yaml
name: task-sync
trigger:
  type: webhook
  source: todoist
steps:
  1. Receive task event (created/completed/updated)
  2. Parse task metadata (project, labels)
  3. Route to appropriate agent context
  4. Update agent memory in Letta
  5. If completed: Update metrics
  6. Log to Obsidian activity feed
```

### 3. Vault Backup Automation

```yaml
name: vault-backup
trigger:
  type: cron
  schedule: "0 */1 * * *"  # Hourly
steps:
  1. Check for vault changes
  2. If changes detected:
     - Stage changes
     - Create commit
     - Push to GitHub
  3. Log backup status
  4. On failure: Alert Grid agent
```

### 4. HITL Approval Router

```yaml
name: hitl-approval
trigger:
  type: webhook
  source: letta
steps:
  1. Receive approval request from agent
  2. Parse request type and priority
  3. If auto-approvable:
     - Execute action
     - Log approval
  4. Else:
     - Create Todoist task with @review label
     - Set appropriate priority
     - Notify Zeus
  5. Wait for completion webhook
  6. Route result back to agent
```

### 5. Health Check Monitor

```yaml
name: health-check
trigger:
  type: cron
  schedule: "*/5 * * * *"  # Every 5 minutes
steps:
  1. Check Letta API health
  2. Check Todoist API health
  3. Check GitHub API health
  4. Verify vault sync status
  5. If any failure:
     - Alert Grid agent
     - Create incident ticket
     - Notify Zeus if critical
```

## Workflow Templates

### Workflow: Agent-to-Agent Communication

```json
{
  "name": "agent-communication",
  "nodes": [
    {
      "type": "webhook",
      "parameters": {
        "path": "agent-message",
        "method": "POST"
      }
    },
    {
      "type": "function",
      "parameters": {
        "code": "// Parse source and target agent\nconst {source, target, message, priority} = $input.body;\nreturn {source, target, message, priority};"
      }
    },
    {
      "type": "http",
      "parameters": {
        "url": "={{$env.LETTA_URL}}/agents/{{$json.target}}/messages",
        "method": "POST",
        "body": {
          "content": "={{$json.message}}",
          "metadata": {
            "source_agent": "={{$json.source}}",
            "priority": "={{$json.priority}}"
          }
        }
      }
    }
  ]
}
```

### Workflow: Obsidian Note Creation

```json
{
  "name": "create-obsidian-note",
  "nodes": [
    {
      "type": "webhook",
      "parameters": {
        "path": "obsidian/create",
        "method": "POST"
      }
    },
    {
      "type": "function",
      "parameters": {
        "code": "// Generate file path and content\nconst {folder, title, content, frontmatter} = $input.body;\nconst path = `${folder}/${title}.md`;\nconst fm = Object.entries(frontmatter).map(([k,v]) => `${k}: ${v}`).join('\\n');\nconst fullContent = `---\\n${fm}\\n---\\n\\n${content}`;\nreturn {path, content: fullContent};"
      }
    },
    {
      "type": "github",
      "parameters": {
        "operation": "create",
        "owner": "={{$env.GITHUB_OWNER}}",
        "repo": "={{$env.GITHUB_REPO}}",
        "path": "={{$json.path}}",
        "content": "={{$json.content}}"
      }
    }
  ]
}
```

## Webhook Endpoints

| Endpoint | Source | Purpose |
|----------|--------|---------|
| `/webhook/todoist` | Todoist | Task events |
| `/webhook/github` | GitHub | Repository events |
| `/webhook/letta` | Letta | Agent requests |
| `/webhook/agent-message` | Agents | Inter-agent comms |
| `/webhook/obsidian/create` | Agents | Note creation |
| `/webhook/health` | Monitors | Health reports |

## Environment Variables

```bash
# N8N Configuration
N8N_URL=https://n8n.your-domain.com
N8N_WEBHOOK_URL=https://n8n.your-domain.com/webhook
N8N_API_KEY=your_n8n_api_key

# Connected Services
LETTA_URL=http://localhost:8283
LETTA_API_KEY=your_letta_key
TODOIST_API_TOKEN=your_todoist_token
GITHUB_TOKEN=your_github_token

# Notification Services
PUSHOVER_USER_KEY=your_pushover_key
PUSHOVER_API_TOKEN=your_pushover_token
```

## Error Handling

| Error Type | Response | Escalation |
|------------|----------|------------|
| Webhook Timeout | Retry 3x with backoff | Log to error queue |
| API Rate Limit | Pause workflow, retry | Alert Grid agent |
| Auth Failure | Stop workflow | Critical alert to Zeus |
| Data Validation | Skip record, log error | Review queue |

## Workflow Categories

| Category | Prefix | Purpose |
|----------|--------|---------|
| Agent | `agent-*` | Agent-related workflows |
| Sync | `sync-*` | Data synchronization |
| Health | `health-*` | Monitoring workflows |
| Notify | `notify-*` | Notification workflows |
| Backup | `backup-*` | Backup workflows |

## Security Configuration

```yaml
security:
  # Webhook authentication
  webhook_auth:
    type: header
    header_name: X-Webhook-Secret
    secret: ${N8N_WEBHOOK_SECRET}

  # IP allowlist (optional)
  allowed_ips:
    - 127.0.0.1
    - ${LETTA_SERVER_IP}

  # Rate limiting
  rate_limit:
    requests_per_minute: 60
    burst: 10
```

## Monitoring

| Metric | Alert Threshold | Action |
|--------|-----------------|--------|
| Workflow Failures | >3 per hour | Grid agent notification |
| Execution Time | >30s average | Performance review |
| Queue Depth | >100 items | Scale warning |
| Memory Usage | >80% | Restart consideration |
