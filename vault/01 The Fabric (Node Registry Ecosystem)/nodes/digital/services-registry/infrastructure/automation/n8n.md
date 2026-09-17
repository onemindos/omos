---
title: "n8n"
type: reference
status: active
tags: [platform, infrastructure, automation, webhooks, nats, saas-gateway]
---
# n8n

> Universal SaaS gateway — bidirectional bridge between NATS and external services.

## Role in OneMind OS
- **INBOUND (SaaS → NATS):** Receives webhooks from external services, validates, transforms, publishes to NATS
- **OUTBOUND (NATS → SaaS):** Subscribes to `cmd.>` subjects, executes actions on external SaaS APIs
- **Polling bridge:** For services without webhooks, polls APIs on schedule and publishes to NATS

## Does NOT Handle
- Fast rules / automations → subscriber.py
- AI reasoning → Hermes
- Outbound notifications (Telegram/Discord/SMS) → notify_router
- Long-term storage → ClickHouse

## Access
- UI: `http://localhost:5678` or `http://n8n.omos:8888`
- Webhook URL pattern: `http://localhost:5678/webhook/<workflow-id>`
- Config/data: `~/.n8n/`
- Database: `~/.n8n/database.sqlite` (SQLite mode)
- Logs: `~/.n8n/logs/`
- LaunchAgent: `com.onemind.n8n`

## Installed Community Nodes
- **n8n-nodes-synadia** — NATS Core, JetStream, KV, Object Store (publish + trigger)

## Architecture

```
INBOUND (SaaS → OneMind)
═══════════════════════════════════════════════
External SaaS ── webhook POST ──→ n8n (:5678)
                                    │
                                    │ Synadia NATS Publish
                                    ▼
                                  NATS (webhooks.>)
                                    │
                     ┌──────────────┼──────────────┐
                     ▼              ▼              ▼
               ClickHouse     subscriber.py    Hermes
               (stores)       (fast rules)    (reasoning)


OUTBOUND (OneMind → SaaS)
═══════════════════════════════════════════════
Hermes / Agent / Subscriber
    │
    │ nats pub cmd.<service>.<action> '{...}'
    ▼
  NATS (cmd.>)
    │
    │ Synadia NATS Trigger
    ▼
  n8n Workflow
    │
    │ Built-in SaaS node (Notion, Linear, GitHub, Google, etc.)
    ▼
  External SaaS API
```

## NATS Subject Conventions

### Inbound (webhooks arriving)
- `webhooks.github.<event>` — push, pr, issue
- `webhooks.stripe.<event>` — payment_succeeded, subscription_created
- `webhooks.todoist.<event>` — item_completed, item_added
- `webhooks.linear.<event>` — issue_created, issue_updated
- `webhooks.cal.<event>` — booking_created, booking_cancelled

### Outbound (commands to execute)
- `cmd.linear.create_issue` — create issue in Linear
- `cmd.notion.create_page` — create page in Notion
- `cmd.github.create_issue` — open GitHub issue
- `cmd.cal.create_event` — book calendar event
- `cmd.sheets.append_row` — add row to Google Sheets
- `cmd.slack.post` — post to Slack channel

## Resource Usage
- RAM: ~320MB (SQLite mode, no Redis/Postgres)
- Disk: ~500MB (binary + node_modules)
- CPU: near-zero idle

## External Access (future)
When ready to receive webhooks from the internet:
- Cloudflare Tunnel route: `webhooks.yourdomain.com → localhost:5678`
- Or NetBird public endpoint

## Links
- [n8n Docs](https://docs.n8n.io)
- [Synadia NATS Nodes](https://github.com/synadia-labs/n8n-nodes-synadia)
