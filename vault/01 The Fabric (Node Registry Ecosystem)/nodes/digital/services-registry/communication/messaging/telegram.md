# Telegram

> **Messaging platform with powerful bot API**

## Overview

Telegram provides secure messaging with a robust Bot API for automation, notifications, and interactive agent communication.

## Status

| Attribute   | Value              |
| ----------- | ------------------ |
| Status      | ✅ Active          |
| Category    | Communication      |
| Type        | Messaging Platform |
| Integration | Bot API            |

## Use Cases in OneMind

1. **Mobile Notifications** - Quick alerts to phone
2. **Bot Interface** - Chat with Legacy AI on mobile
3. **File Sharing** - Send/receive documents
4. **Group Management** - Team communication

## Key Features

- End-to-end encryption (Secret Chats)
- Powerful Bot API
- Inline keyboards
- File/media support
- Groups and channels
- Webhooks or polling

## API Integration

### Bot Setup

```python
from telegram import Update, InlineKeyboardButton, InlineKeyboardMarkup
from telegram.ext import Application, CommandHandler, MessageHandler, filters

TELEGRAM_BOT_TOKEN = "your_bot_token"

app = Application.builder().token(TELEGRAM_BOT_TOKEN).build()
```

### Send Message

```python
import httpx

async def send_telegram_message(
    chat_id: str,
    text: str,
    parse_mode: str = "HTML"
):
    """Send message via Telegram Bot API."""
    url = f"https://api.telegram.org/bot{TELEGRAM_BOT_TOKEN}/sendMessage"

    async with httpx.AsyncClient() as client:
        response = await client.post(url, json={
            "chat_id": chat_id,
            "text": text,
            "parse_mode": parse_mode
        })
        return response.json()
```

### Bot Handlers

```python
async def start(update: Update, context):
    """Handle /start command."""
    await update.message.reply_text(
        "👋 Welcome to Legacy AI!\n\n"
        "Commands:\n"
        "/ask - Ask a question\n"
        "/status - System status\n"
        "/tasks - View tasks"
    )

async def ask_handler(update: Update, context):
    """Forward question to Legacy AI."""
    question = " ".join(context.args)
    response = await legacy_ai.query(question)
    await update.message.reply_text(response)

# Register handlers
app.add_handler(CommandHandler("start", start))
app.add_handler(CommandHandler("ask", ask_handler))
```

### Inline Keyboards

```python
async def show_options(update: Update, context):
    """Show interactive buttons."""
    keyboard = [
        [
            InlineKeyboardButton("✅ Approve", callback_data="approve"),
            InlineKeyboardButton("❌ Deny", callback_data="deny")
        ],
        [InlineKeyboardButton("📊 Details", callback_data="details")]
    ]

    reply_markup = InlineKeyboardMarkup(keyboard)
    await update.message.reply_text(
        "Action required:",
        reply_markup=reply_markup
    )
```

## Legacy AI Integration

```python
from agno import Tool

@Tool
async def send_telegram_notification(
    message: str,
    chat_id: str = None,
    buttons: list[str] = None
) -> str:
    """Send notification via Telegram.

    Args:
        message: Message to send
        chat_id: Target chat (defaults to Zeus)
        buttons: Optional inline keyboard buttons
    """
    chat_id = chat_id or ZEUS_CHAT_ID

    await send_telegram_message(chat_id, message)
    return "Telegram notification sent"

@Tool
async def request_telegram_approval(
    action: str,
    details: str
) -> bool:
    """Request approval via Telegram with buttons.

    Args:
        action: What action needs approval
        details: Context about the action
    """
    # Send with approve/deny buttons
    # Wait for callback response
    pass
```

## Environment Variables

```bash
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_ZEUS_CHAT_ID=your_chat_id
TELEGRAM_WEBHOOK_URL=https://your-domain.com/telegram/webhook
```

## Webhook Setup

```python
from fastapi import FastAPI, Request

app = FastAPI()

@app.post("/telegram/webhook")
async def telegram_webhook(request: Request):
    """Handle incoming Telegram updates."""
    update = await request.json()
    # Process update
    return {"ok": True}

# Set webhook
async def set_webhook():
    url = f"https://api.telegram.org/bot{TOKEN}/setWebhook"
    await httpx.post(url, json={
        "url": "https://your-domain.com/telegram/webhook"
    })
```

## Best Practices

1. **Use webhooks in production** - More efficient than polling
2. **HTML parse mode** - Supports basic formatting
3. **Inline keyboards** - Great for approvals
4. **Rate limits** - 30 messages/second to same chat

## Resources

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [python-telegram-bot](https://python-telegram-bot.org/)
- [Bot Father](https://t.me/botfather)

---

Last Updated: January 25, 2026


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
