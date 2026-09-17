# OneMind Notification System

> AI-aware notification routing with digest, throttle, and multi-channel delivery.

## Overview

The OneMind notification system replaces Novu with a custom NATS-based solution that understands:

- **Awareness modes** (dormant/aware/present/omnipresent)
- **Priority-based routing** (critical → all channels, low → in-app only)
- **Smart digest batching** (group low-priority notifications)
- **Escalation** (if no response, try SMS → voice)

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    NOTIFICATION ROUTER                           │
├─────────────────────────────────────────────────────────────────┤
│  1. Check awareness mode                                         │
│  2. Apply throttle limits                                        │
│  3. Add to digest OR send immediately                            │
│  4. Route to channels based on priority                          │
│  5. Schedule escalation if require_ack=True                      │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                      CHANNEL ROUTING                             │
├─────────────────────────────────────────────────────────────────┤
│  CRITICAL  →  SMS + Push + Discord + In-App (immediately)        │
│  HIGH      →  Push + Discord + In-App                            │
│  NORMAL    →  Push + In-App                                      │
│  LOW       →  In-App only (may be digested)                      │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                       CHANNELS                                   │
├─────────┬─────────┬─────────┬─────────┬─────────┬───────────────┤
│  FCM    │ Twilio  │ Discord │Telegram │ Resend  │  WebSocket    │
│ (push)  │ (SMS)   │(webhook)│  (bot)  │ (email) │  (in-app)     │
│  FREE   │ $0.01/  │  FREE   │  FREE   │ $0.001/ │    FREE       │
└─────────┴─────────┴─────────┴─────────┴─────────┴───────────────┘
```

## Features

### 1. Awareness Mode Integration

| Mode | Behavior |
|------|----------|
| **Dormant** | Only critical alerts get through (SMS + Push) |
| **Aware** | Critical + High priority only |
| **Present** | All priorities, normal routing |
| **Omnipresent** | Everything immediately, no batching |

### 2. Digest Engine

Batches multiple notifications into single messages:

```python
# 5 GitHub notifications in 5 minutes become:
# "5 new notifications from GitHub"
#   - PR #123 merged
#   - Issue #456 commented
#   - ...and 3 more
```

**Strategies:**
- **Regular**: Always batch for specified window (e.g., 5 minutes)
- **Look-back**: If no recent message, send immediately; else batch
- **Scheduled**: Send at specific times (e.g., 9 AM daily briefing)

### 3. Throttle Engine

Prevents notification overload:

| Priority | Default Limit | Window |
|----------|---------------|--------|
| Critical | 100/hour | 1 hour |
| High | 20/hour | 1 hour |
| Normal | 10/hour | 1 hour |
| Low | 5/hour | 1 hour |

### 4. Escalation

For notifications requiring acknowledgment:

```
Initial: Push + In-App
   ↓ (no ack in 5 min)
Escalate: SMS
   ↓ (no ack in 5 min)
Final: Voice call
```

## Channel Configuration

### Firebase Cloud Messaging (Push)

```bash
# Environment variables
FIREBASE_PROJECT_ID=your_project_id
GOOGLE_APPLICATION_CREDENTIALS=/app/firebase-credentials.json
```

**Cost:** FREE (unlimited)

### Twilio (SMS/Voice)

```bash
TWILIO_ACCOUNT_SID=ACxxxxxxxx
TWILIO_AUTH_TOKEN=your_auth_token
TWILIO_PHONE_NUMBER=+1234567890
```

**Cost:** ~$0.01/SMS, ~$0.02/min voice

### Discord (Webhook)

```bash
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/xxx/yyy
```

**Cost:** FREE

### Telegram (Bot)

```bash
TELEGRAM_BOT_TOKEN=123456:ABC-DEF...
TELEGRAM_CHAT_ID=your_chat_id
```

**Cost:** FREE

### Resend (Email)

```bash
RESEND_API_KEY=re_xxxxx
RESEND_FROM_EMAIL=Legacy <notifications@onemindos.dev>
```

**Cost:** ~$0.001/email (100 free/day)

## API Endpoints

### Send Notification

```bash
POST /api/notifications
{
  "title": "Server Alert",
  "body": "CPU usage above 90%",
  "priority": "high",
  "source": "monitoring",
  "require_ack": true
}
```

### Acknowledge

```bash
POST /api/notifications/acknowledge
{
  "notification_id": "uuid-here"
}
```

### Get Digest Stats

```bash
GET /api/notifications/digest/stats?user_id=zeus
```

### Flush Digests

```bash
POST /api/notifications/digest/flush?user_id=zeus
```

### Get Throttle Stats

```bash
GET /api/notifications/throttle/stats?user_id=zeus&priority=normal
```

### Update Preferences

```bash
PATCH /api/notifications/preferences?user_id=zeus
{
  "quiet_hours_enabled": true,
  "quiet_hours_start": "22:00",
  "quiet_hours_end": "07:00"
}
```

## Legacy's Tools

Legacy has four tools for contacting Zeus:

### contact_zeus

Main tool for any notification:

```python
await contact_zeus(
    message="Database backup completed successfully",
    priority="normal",  # critical, high, normal, low
    context="Backup Complete",
    require_response=False
)
```

### urgent_alert

For critical situations:

```python
await urgent_alert(
    message="Server disk space critical: 5% remaining",
    context="Disk Alert"
)
# Sends via ALL channels, schedules voice call escalation
```

### request_approval

For HITL scenarios:

```python
await request_approval(
    action_description="Delete 500 old log files",
    reason="Free up 2GB disk space",
    options=["Approve", "Deny", "Review Files"]
)
```

### send_update

For non-urgent FYI:

```python
await send_update(
    message="Morning briefing: 3 tasks due today",
    context="Daily Summary"
)
# In-app only, may be batched
```

## File Structure

```
core/notifications/
├── __init__.py           # Exports
├── models.py             # Pydantic models
├── router.py             # Main notification router
├── digest.py             # Digest engine (batching)
├── throttle.py           # Throttle engine (rate limiting)
├── api.py                # FastAPI endpoints
└── channels/
    ├── __init__.py
    ├── base.py           # Base channel class
    ├── websocket.py      # In-app notifications
    ├── fcm.py            # Firebase push
    ├── twilio.py         # SMS and voice
    ├── discord.py        # Discord webhook
    ├── telegram.py       # Telegram bot
    └── resend.py         # Email

core/tools/
└── contact.py            # Legacy's notification tools
```

## Cost Comparison

| Usage Level | Novu Cloud | OneMind Custom |
|-------------|------------|----------------|
| 1,000/month | ~$25 | ~$5 (SMS only) |
| 10,000/month | ~$100 | ~$20 |
| 100,000/month | ~$500+ | ~$50 |

**Breakdown (100k/month):**
- Push: FREE (FCM)
- In-app: FREE (WebSocket)
- Discord/Telegram: FREE
- SMS (1% critical): ~$10
- Email (10%): ~$10
- Voice (rare): ~$5

## Initialization

```python
from core.notifications.router import create_router

# In your app startup
async def startup():
    router = await create_router(
        redis_client=redis,
        nats_client=nats,
        ws_manager=ws_manager,
        db_pool=db_pool,
    )
```

## Why Not Novu?

| Aspect | Novu | OneMind Custom |
|--------|------|----------------|
| **AI Awareness** | None | Full integration |
| **Awareness Modes** | N/A | dormant/aware/present/omnipresent |
| **Context-Based** | Templates | AI decides channel |
| **Offline/Edge** | Cloud only | Works with NATS leaf |
| **Cost at Scale** | $200-500/mo | $20-50/mo |
| **Vendor Lock-in** | High | None |
| **Setup Time** | 2 hours | 2-3 days |

The key differentiator: **Novu sends notifications based on templates. OneMind sends notifications based on AI understanding of context.**

## Troubleshooting

### Notifications not sending

1. Check channel configuration:
```bash
curl http://localhost:8081/api/notifications/status
```

2. Check throttle:
```bash
curl "http://localhost:8081/api/notifications/throttle/stats?user_id=zeus"
```

3. Check digest queue:
```bash
curl "http://localhost:8081/api/notifications/digest/stats?user_id=zeus"
```

### FCM not working

1. Verify Firebase credentials:
```bash
echo $GOOGLE_APPLICATION_CREDENTIALS
cat /app/firebase-credentials.json | head -5
```

2. Check FCM token is registered in user preferences

### Escalation not triggering

1. Verify `require_ack=True` was set
2. Check Redis is running (escalation tasks stored there)
3. Check acknowledgment wasn't already sent

## Migration from Novu

1. Keep Novu running during transition
2. Update inbox service to use new router
3. Register FCM tokens from Flutter app
4. Test each channel independently
5. Disable Novu once verified

The inbox system already has WebSocket broadcasts working. The new router adds:
- Digest/throttle via Redis
- Multi-channel routing
- Legacy's contact tools
