---
name: telegram
description: Interact with Telegram Bot API. Use when you need to send messages, manage bots, handle updates, or create Telegram integrations.
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Read
  - Write
---

# Telegram Skill

Build and interact with Telegram bots via the Bot API.

## Environment Variables

```bash
TELEGRAM_BOT_TOKEN   # Bot token from @BotFather (format: 123456:ABC-DEF1234...)
TELEGRAM_CHAT_ID     # Default chat ID for sending messages
```

## Base URL

```
https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN
```

## Operations

### Bot Info

**Get bot info:**
```bash
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getMe" | jq '.result'
```

### Sending Messages

**Send text message:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "'"$TELEGRAM_CHAT_ID"'",
    "text": "Hello from Claude!",
    "parse_mode": "Markdown"
  }'
```

**Send with formatting:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "'"$TELEGRAM_CHAT_ID"'",
    "text": "*Bold* _italic_ `code` ```preformatted```",
    "parse_mode": "Markdown"
  }'
```

**Send with HTML:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "'"$TELEGRAM_CHAT_ID"'",
    "text": "<b>Bold</b> <i>italic</i> <code>code</code>",
    "parse_mode": "HTML"
  }'
```

**Send with inline keyboard:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "'"$TELEGRAM_CHAT_ID"'",
    "text": "Choose an option:",
    "reply_markup": {
      "inline_keyboard": [
        [{"text": "Option 1", "callback_data": "opt1"}, {"text": "Option 2", "callback_data": "opt2"}],
        [{"text": "Visit Website", "url": "https://example.com"}]
      ]
    }
  }'
```

**Send silently (no notification):**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendMessage" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "'"$TELEGRAM_CHAT_ID"'",
    "text": "Silent message",
    "disable_notification": true
  }'
```

### Media

**Send photo:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendPhoto" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "'"$TELEGRAM_CHAT_ID"'",
    "photo": "https://example.com/image.jpg",
    "caption": "Photo caption"
  }'
```

**Send document:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendDocument" \
  -F "chat_id=$TELEGRAM_CHAT_ID" \
  -F "document=@/path/to/file.pdf" \
  -F "caption=Document caption"
```

**Send location:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/sendLocation" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "'"$TELEGRAM_CHAT_ID"'",
    "latitude": 40.7128,
    "longitude": -74.0060
  }'
```

### Updates (Polling)

**Get updates:**
```bash
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getUpdates" | jq '.result'
```

**Get updates with offset (mark as read):**
```bash
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getUpdates?offset=12345"
```

**Long polling:**
```bash
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getUpdates?timeout=30"
```

### Webhooks

**Set webhook:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setWebhook" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://your-server.com/telegram/webhook",
    "secret_token": "your-secret-token"
  }'
```

**Get webhook info:**
```bash
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getWebhookInfo" | jq
```

**Delete webhook:**
```bash
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/deleteWebhook"
```

### Chat Management

**Get chat info:**
```bash
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getChat?chat_id=$TELEGRAM_CHAT_ID"
```

**Get chat member count:**
```bash
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getChatMemberCount?chat_id=$TELEGRAM_CHAT_ID"
```

**Leave chat:**
```bash
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/leaveChat?chat_id=$TELEGRAM_CHAT_ID"
```

### Message Management

**Edit message:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/editMessageText" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "'"$TELEGRAM_CHAT_ID"'",
    "message_id": 123,
    "text": "Updated message text"
  }'
```

**Delete message:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/deleteMessage" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "'"$TELEGRAM_CHAT_ID"'",
    "message_id": 123
  }'
```

**Pin message:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/pinChatMessage" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": "'"$TELEGRAM_CHAT_ID"'",
    "message_id": 123
  }'
```

### Commands

**Set bot commands:**
```bash
curl -s -X POST "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/setMyCommands" \
  -H "Content-Type: application/json" \
  -d '{
    "commands": [
      {"command": "start", "description": "Start the bot"},
      {"command": "help", "description": "Get help"},
      {"command": "status", "description": "Check status"}
    ]
  }'
```

**Get bot commands:**
```bash
curl -s "https://api.telegram.org/bot$TELEGRAM_BOT_TOKEN/getMyCommands" | jq '.result'
```

## Parse Modes

| Mode | Syntax |
|------|--------|
| `Markdown` | `*bold*` `_italic_` `` `code` `` |
| `MarkdownV2` | `*bold*` `_italic_` `` `code` `` (escape special chars) |
| `HTML` | `<b>bold</b>` `<i>italic</i>` `<code>code</code>` |

## Chat ID Types

| Type | Format |
|------|--------|
| User | Positive number (e.g., `123456789`) |
| Group | Negative number (e.g., `-123456789`) |
| Supergroup | Starts with `-100` (e.g., `-1001234567890`) |
| Channel | `@channelname` or `-100...` |

## Tips

- Get your chat ID: Send a message to your bot, then call `getUpdates`
- Bot tokens look like: `123456789:ABCdefGHIjklMNO-pqrSTUvwxYZ`
- Create bots via @BotFather on Telegram
- Use webhooks for production, polling for development
- Rate limits: ~30 messages/second to same chat


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
