---
name: webhook-notifier
description: >
  Sends structured webhook notifications to any HTTP endpoint. Use when
  you need to notify an external service, trigger a webhook, ping a URL,
  send a Discord webhook, post to a Slack channel, or call an n8n/Zapier
  automation trigger. Supports custom headers, payloads, and retry logic.
metadata:
  author: zeus
  version: "1.0.0"
---

# Webhook Notifier Skill

## When to Use

Activate when the user wants to:

- Send a webhook notification or trigger
- Notify an external service (Discord, Slack, n8n, Zapier, Make, etc.)
- Ping a URL with data
- Fire an automation trigger
- Post a structured event to an HTTP endpoint

## Available Tools

| Tool              | Description                                      |
| ----------------- | ------------------------------------------------ |
| `webhook_send`    | Send a POST to any webhook URL with JSON payload |
| `webhook_discord` | Send a pre-formatted Discord webhook message     |
| `webhook_slack`   | Send a pre-formatted Slack webhook message       |
| `webhook_ping`    | Simple GET/HEAD ping to check if a URL is alive  |

## Quick Reference

### Send a generic webhook

```
webhook_send({
  url: "https://hooks.example.com/webhook/abc123",
  payload: { event: "task_completed", task: "Morning briefing", user: "legacy" }
})
```

### Discord notification

```
webhook_discord({
  url: "https://discord.com/api/webhooks/...",
  title: "System Alert",
  description: "Disk usage is approaching 90% on VPS",
  color: "red"
})
```

### Slack notification

```
webhook_slack({
  url: "https://hooks.slack.com/services/...",
  text: "Deployment completed successfully ✅",
  username: "Legacy"
})
```

### Ping a URL

```
webhook_ping({ url: "https://api.example.com/health" })
```

### Send with custom headers (e.g., API key auth)

```
webhook_send({
  url: "https://api.example.com/trigger",
  payload: { action: "deploy", environment: "production" },
  headers: { "X-API-Key": "my-secret-key", "X-Source": "legacy" }
})
```

## Retry Behavior

`webhook_send` and `webhook_discord` retry once on 5xx errors or network timeout.
They do NOT retry on 4xx errors (client errors — fix the request instead).

## Colors for Discord

| Color Name | Hex     | When to Use           |
| ---------- | ------- | --------------------- |
| `green`    | #57F287 | Success, completed    |
| `red`      | #ED4245 | Error, failure, alert |
| `yellow`   | #FEE75C | Warning, in-progress  |
| `blue`     | #5865F2 | Info, general         |
| `purple`   | #9B59B6 | Special, important    |
| `gray`     | #95A5A6 | Neutral, archived     |

## Notes

- Webhook URLs are passed at call time — never need to be configured in env vars
- If a webhook returns 4xx, the error message is included in the response (don't retry)
- `webhook_ping` is useful before sending data to verify an endpoint is alive
- Payloads are always JSON (`Content-Type: application/json`)
