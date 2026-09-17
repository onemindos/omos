---
title: "Notification Router"
type: reference
status: active
tags: [platform, infrastructure, messaging, notifications, nats]
---
# Notification Router

> Universal notification delivery service — routes NATS messages to Telegram, Discord, Push, SMS, Email.

## Role in OneMind OS
- Single point for ALL outbound notifications
- Subscribes to `cmd.notify` (new) and `cmd.telegram.send` (legacy compat)
- Reads routing config + DND rules from NATS KV bucket `notifications`
- Supports priority-based routing (critical → all channels, low → push only)

## Access
- NATS subject: `cmd.notify`
- Config: NATS KV bucket `notifications`
- Source: `/Users/onemind/onemind/subscribers/notify_router.py`
- Logs: `/Users/onemind/onemind/subscribers/logs/notify.log`
- LaunchAgent: `com.onemind.notify-router`

## Supported Channels

| Channel | Method | KV Key | Status |
|---------|--------|--------|--------|
| Telegram | Bot API | `config/telegram` | 🔲 Token needed |
| Discord | Webhook | `config/discord` | 🔲 Not configured |
| Push | ntfy.sh | `config/push` | 🔲 Not configured |
| SMS | Twilio | `config/sms` | 🔲 Not configured |
| Email | SMTP | `config/email` | 🔲 Not configured |

## Usage

```bash
# Simple notification (goes to all enabled channels)
nats pub cmd.notify '{"message": "Hello from OneMind"}'

# With priority and channel routing
nats pub cmd.notify '{"message": "Security alert!", "priority": "critical", "channels": ["telegram", "sms"]}'

# Legacy format (backward compat with existing subscribers)
nats pub cmd.telegram.send '{"text": "This still works"}'
```

## Configure Channels

```bash
# Telegram
nats kv put notifications config/telegram '{"bot_token": "TOKEN", "chat_id": "CHAT_ID", "enabled": true}'

# Discord
nats kv put notifications config/discord '{"webhook_url": "https://discord.com/api/webhooks/...", "enabled": true}'

# DND / Quiet Hours
nats kv put notifications config/dnd '{"start": "23:00", "end": "06:00", "timezone": "America/New_York"}'
```

## Priority Routing

| Priority | Behavior |
|----------|----------|
| critical | All channels, ignores DND |
| high | All channels, respects DND |
| normal | Default channels only |
| low | Push/in-app only |

## Links
- [ntfy.sh](https://ntfy.sh) — self-hostable push notifications
- [Telegram Bot API](https://core.telegram.org/bots/api)
