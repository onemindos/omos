# Discord

> **Community chat, notifications, and bot integrations**

## Overview

Discord serves as a real-time communication hub for community engagement, team collaboration, and automated notifications via webhooks and bots.

## Status

| Attribute   | Value             |
| ----------- | ----------------- |
| Status      | ✅ Active         |
| Category    | Communication     |
| Type        | Chat Platform     |
| Integration | Webhooks, Bot API |

## Use Cases in OneMind

1. **Community Hub** - Engage with community members
2. **Notifications** - Send alerts from Legacy AI
3. **Bot Commands** - Interactive agent responses
4. **Approval Workflows** - Request approvals via reactions

## Key Features

- Real-time messaging and voice
- Webhooks for notifications
- Bot API for automation
- Role-based permissions
- Thread support
- Rich embeds

## API Integration

### Webhook Notifications

```python
import httpx

DISCORD_WEBHOOK_URL = "https://discord.com/api/webhooks/..."

async def send_discord_notification(
    content: str,
    username: str = "Legacy AI",
    embed: dict = None
):
    """Send notification to Discord channel."""
    payload = {
        "username": username,
        "content": content
    }

    if embed:
        payload["embeds"] = [embed]

    async with httpx.AsyncClient() as client:
        response = await client.post(
            DISCORD_WEBHOOK_URL,
            json=payload
        )
        return response.status_code == 204
```

### Rich Embed

```python
embed = {
    "title": "Task Completed",
    "description": "Legacy AI finished processing your request",
    "color": 0x9B59B6,  # Purple
    "fields": [
        {"name": "Task", "value": "Data sync", "inline": True},
        {"name": "Duration", "value": "2.3s", "inline": True}
    ],
    "footer": {"text": "OneMind OS"}
}

await send_discord_notification("", embed=embed)
```

### Bot Integration

```python
import discord
from discord.ext import commands

bot = commands.Bot(command_prefix="!", intents=discord.Intents.default())

@bot.command()
async def ask(ctx, *, question: str):
    """Ask Legacy AI a question."""
    # Forward to Legacy AI
    response = await legacy_ai.query(question)
    await ctx.reply(response)

@bot.command()
async def status(ctx):
    """Get system status."""
    await ctx.reply("🟢 OneMind OS is operational")
```

## Legacy AI Integration

```python
from agno import Tool

@Tool
async def send_discord_message(
    channel: str,
    message: str,
    mention_user: str = None
) -> str:
    """Send a message to Discord channel.

    Args:
        channel: Channel name or webhook identifier
        message: Message content
        mention_user: Optional user ID to mention
    """
    webhook_url = get_webhook_for_channel(channel)

    if mention_user:
        message = f"<@{mention_user}> {message}"

    await send_discord_notification(message)
    return f"Message sent to #{channel}"
```

## Environment Variables

```bash
DISCORD_BOT_TOKEN=your_bot_token
DISCORD_WEBHOOK_GENERAL=https://discord.com/api/webhooks/...
DISCORD_WEBHOOK_ALERTS=https://discord.com/api/webhooks/...
DISCORD_GUILD_ID=your_server_id
```

## Best Practices

1. **Use webhooks for notifications** - Simpler than bot for one-way messages
2. **Rate limiting** - Discord has strict rate limits (5 requests/second)
3. **Rich embeds** - Use for structured information
4. **Thread for conversations** - Keep channels clean

## Resources

- [Discord Developer Portal](https://discord.com/developers)
- [Discord.py Docs](https://discordpy.readthedocs.io/)
- [Webhook Guide](https://discord.com/developers/docs/resources/webhook)

---

Last Updated: January 25, 2026
