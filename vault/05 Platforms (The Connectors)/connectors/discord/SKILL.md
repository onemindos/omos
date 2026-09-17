---
name: discord
description: Interact with Discord via webhooks and bot API. Use when you need to send notifications, request approvals, post alerts, or manage bot messages.
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Read
  - Write
---

# Discord Integration Skill

Send messages and manage interactions via Discord webhooks and Bot API.

## Environment Variables

```bash
DISCORD_WEBHOOK_URL   # Webhook URL for notifications
DISCORD_BOT_TOKEN     # Bot token (optional, for advanced)
DISCORD_CHANNEL_ID    # Default channel ID
```

---

## Webhook Operations (Simple)

### Send Message
```bash
curl -s -X POST "$DISCORD_WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello from Claude!",
    "username": "Legacy AI"
  }'
```

### Send Embed
```bash
curl -s -X POST "$DISCORD_WEBHOOK_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "Legacy AI",
    "embeds": [{
      "title": "Alert Title",
      "description": "Alert description here",
      "color": 16711680,
      "fields": [
        {"name": "Field 1", "value": "Value 1", "inline": true},
        {"name": "Field 2", "value": "Value 2", "inline": true}
      ],
      "footer": {"text": "One Mind System"}
    }]
  }'
```

### Approval Request (with reactions)
```bash
curl -s -X POST "$DISCORD_WEBHOOK_URL?wait=true" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "Approval Bot",
    "embeds": [{
      "title": "Approval Required",
      "description": "Task: Deploy to production\n\nReact with ✅ to approve or ❌ to reject",
      "color": 16776960,
      "fields": [
        {"name": "Requested By", "value": "Legacy AI", "inline": true},
        {"name": "Priority", "value": "High", "inline": true}
      ]
    }]
  }'
```

---

## Bot API Operations (Advanced)

### Base URL
```
https://discord.com/api/v10
```

### Authentication
```bash
Authorization: Bot $DISCORD_BOT_TOKEN
```

### Send Message to Channel
```bash
curl -s -X POST "https://discord.com/api/v10/channels/$DISCORD_CHANNEL_ID/messages" \
  -H "Authorization: Bot $DISCORD_BOT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "Message text"}'
```

### Add Reaction
```bash
curl -s -X PUT "https://discord.com/api/v10/channels/{channel_id}/messages/{message_id}/reactions/{emoji}/@me" \
  -H "Authorization: Bot $DISCORD_BOT_TOKEN"
```

### Get Message Reactions
```bash
curl -s "https://discord.com/api/v10/channels/{channel_id}/messages/{message_id}/reactions/{emoji}" \
  -H "Authorization: Bot $DISCORD_BOT_TOKEN"
```

### Edit Message
```bash
curl -s -X PATCH "https://discord.com/api/v10/channels/{channel_id}/messages/{message_id}" \
  -H "Authorization: Bot $DISCORD_BOT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "Updated message"}'
```

### Delete Message
```bash
curl -s -X DELETE "https://discord.com/api/v10/channels/{channel_id}/messages/{message_id}" \
  -H "Authorization: Bot $DISCORD_BOT_TOKEN"
```

---

## Embed Colors

| Color | Hex | Decimal |
|-------|-----|---------|
| Red | #FF0000 | 16711680 |
| Green | #00FF00 | 65280 |
| Blue | #0000FF | 255 |
| Yellow | #FFFF00 | 16776960 |
| Orange | #FFA500 | 16753920 |
| Purple | #800080 | 8388736 |

---

## One Mind Approval Flow

```bash
# 1. Send approval request
response=$(curl -s -X POST "$DISCORD_WEBHOOK_URL?wait=true" \
  -H "Content-Type: application/json" \
  -d '{
    "embeds": [{
      "title": "🔐 Approval Required",
      "description": "**Action:** Create new project\n**Agent:** Legacy AI",
      "color": 16776960
    }]
  }')

# 2. Extract message ID
message_id=$(echo $response | jq -r '.id')

# 3. Check for approval (poll or webhook)
```

## Tips

- Webhooks are simpler but can't read reactions
- Bot API needed for interactive approvals
- Rate limit: 5 messages/5 seconds per channel
- Use `?wait=true` on webhooks to get message ID


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
