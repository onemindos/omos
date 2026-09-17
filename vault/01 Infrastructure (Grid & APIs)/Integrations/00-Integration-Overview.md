---
id: ui-12-20260111053926
title: "One Mind Integration Overview"
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
ai_summary: "One Mind Integration Overview"
---

# One Mind Integration Overview

## Integration Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           ONE MIND HIVE MIND                            │
│                                                                         │
│   ┌─────────┐   ┌────────┐   ┌────────┐   ┌────────┐   ┌────────┐     │
│   │ LEGACY  │   │   HP   │   │   LE   │   │   GE   │   │   IT   │     │
│   │(Mother) │   │ Agent  │   │ Agent  │   │ Agent  │   │ Agent  │     │
│   └────┬────┘   └───┬────┘   └───┬────┘   └───┬────┘   └───┬────┘     │
│        │            │            │            │            │           │
│        └────────────┴────────────┼────────────┴────────────┘           │
│                                  │                                      │
│                         ┌────────▼────────┐                            │
│                         │   LETTA SERVER  │                            │
│                         │  (Agent Memory) │                            │
│                         └────────┬────────┘                            │
│                                  │                                      │
└──────────────────────────────────┼──────────────────────────────────────┘
                                   │
                          ┌────────┴────────┐
                          │                 │
              ┌───────────▼───┐    ┌───────▼───────┐
              │      N8N      │    │  MCP Servers  │
              │  (Workflows)  │    │   (Tools)     │
              └───────┬───────┘    └───────┬───────┘
                      │                    │
        ┌─────────────┼────────────────────┼─────────────┐
        │             │                    │             │
   ┌────▼────┐  ┌─────▼─────┐  ┌──────────▼┐  ┌────────▼────┐
   │ Todoist │  │  GitHub   │  │  Obsidian │  │  Calendar   │
   │ (Tasks) │  │  (Code)   │  │  (Notes)  │  │  (Events)   │
   └─────────┘  └───────────┘  └───────────┘  └─────────────┘
```

## Core Integrations

| Integration | Purpose | Protocol | Config File |
|-------------|---------|----------|-------------|
| **Todoist** | Task execution engine | REST API + MCP | [todoist-config.md](todoist-config.md) |
| **GitHub** | Version control & backup | REST API + MCP | [github-config.md](github-config.md) |
| **N8N** | Workflow automation | Webhooks + REST | [n8n-config.md](n8n-config.md) |
| **Obsidian** | Knowledge management | File system + Git | Native |
| **Letta** | Agent memory & state | REST API | Infrastructure |

## Data Flow Summary

### Task Flow (Todoist)

```
Zeus Request → LEGACY routes → Pillar Agent creates task →
Todoist → Zeus executes → N8N webhook → Agent context updated
```

### Knowledge Flow (Obsidian)

```
Agent insight → LEGACY formats → Git commit →
Obsidian syncs → Zeus reviews → Feedback loop
```

### Automation Flow (N8N)

```
Trigger (cron/webhook) → N8N workflow →
API calls to services → Results to agents →
Actions or notifications
```

## Setup Checklist

### 1. Environment Configuration

Create `.env` file in `00-Infrastructure/`:

```bash
# Letta Configuration
LETTA_URL=http://localhost:8283
LETTA_API_KEY=your_letta_api_key

# Todoist Integration
TODOIST_API_TOKEN=your_todoist_token
TODOIST_INBOX_PROJECT_ID=123456789

# GitHub Integration
GITHUB_TOKEN=ghp_your_token
GITHUB_OWNER=your-username
GITHUB_REPO=one-mind
GITHUB_WEBHOOK_SECRET=your_secret

# N8N Configuration
N8N_URL=https://n8n.your-domain.com
N8N_WEBHOOK_URL=https://n8n.your-domain.com/webhook
N8N_API_KEY=your_n8n_key

# Notifications
PUSHOVER_USER_KEY=your_pushover_user
PUSHOVER_API_TOKEN=your_pushover_token

# Optional
OPENWEATHER_API_KEY=your_weather_key
WEATHER_LOCATION=New York,US
```

### 2. Service Setup Order

1. **Letta Server** - Start first (agent memory)
2. **N8N** - Start second (automation hub)
3. **Configure webhooks** - Connect services to N8N
4. **Run health check** - Verify all connections

### 3. Verification Commands

```bash
# Check all services
python scripts/health_check.py --verbose

# Test Letta connection
curl http://localhost:8283/health

# Test Todoist connection
curl -H "Authorization: Bearer $TODOIST_API_TOKEN" \
  https://api.todoist.com/rest/v2/projects

# Test GitHub connection
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
  https://api.github.com/rate_limit
```

## Workflow Blueprints

Located in `workflows/blueprints/`:

| Blueprint | Purpose | Trigger |
|-----------|---------|---------|
| [morning-brief.json](../workflows/blueprints/morning-brief.json) | Daily morning brief | 6 AM cron |
| [task-sync.json](../workflows/blueprints/task-sync.json) | Todoist ↔ Agent sync | Webhook |
| [health-monitor.json](../workflows/blueprints/health-monitor.json) | System health checks | 5-min cron |

## Security Guidelines

1. **Never commit tokens** - Use `.env` files excluded from git
2. **Use fine-grained tokens** - Minimal necessary permissions
3. **Rotate regularly** - Set calendar reminders for token refresh
4. **Audit access** - Review connected apps periodically
5. **Encrypt at rest** - Use secure credential storage

## Troubleshooting

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| "Connection refused" | Service not running | Start the service |
| "401 Unauthorized" | Invalid/expired token | Regenerate token |
| "Rate limited" | Too many requests | Check rate limit, add delays |
| "Webhook not firing" | Misconfigured URL | Verify webhook URL in service |

### Debug Commands

```bash
# Check Letta logs
docker logs letta-server

# Check N8N workflow history
# Visit N8N UI → Executions

# Test webhook manually
curl -X POST -H "Content-Type: application/json" \
  -d '{"test": true}' \
  https://n8n.your-domain.com/webhook/test
```

## Maintenance Schedule

| Task | Frequency | Owner |
|------|-----------|-------|
| Token rotation | Quarterly | Grid Agent |
| Backup verification | Monthly | Grid Agent |
| Rate limit review | Weekly | Grid Agent |
| Integration audit | Quarterly | LEGACY |
