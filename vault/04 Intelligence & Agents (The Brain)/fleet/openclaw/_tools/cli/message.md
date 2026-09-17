---
name: OpenClaw CLI — Message
updated: 2026-05-11
---

# `openclaw message`

Unified interface for sending messages and channel actions across all platforms.

**Official docs:** [docs.openclaw.ai/cli/message](https://docs.openclaw.ai/cli/message.md)

---

## Commands

```bash
# Send a message
openclaw message send --channel telegram --to "-1003905828404" --text "Hello"
openclaw message send --channel discord --to "channel:<id>" --text "Hello"

# With thread/topic
openclaw message send --channel telegram --to "-1003905828404" --thread-id 292 --text "Hello"

# Read/edit/delete
openclaw message read --channel telegram --to <chat-id> --id <message-id>
openclaw message edit --channel telegram --to <chat-id> --id <message-id> --text "Updated"
openclaw message delete --channel telegram --to <chat-id> --id <message-id>

# Reactions
openclaw message react --channel discord --to "channel:<id>" --id <message-id> --emoji "✅"

# Pin
openclaw message pin --channel telegram --to <chat-id> --id <message-id>

# Poll (Telegram)
openclaw message poll --channel telegram --to <chat-id> --question "..." --options "A" "B" "C"

# Broadcast to all configured channels
openclaw message broadcast --text "System announcement"
```

---

## Target Resolution by Platform

| Platform | Target Format |
|---|---|
| Telegram | Chat ID (`-1003905828404`), `@username`, or forum topic |
| Discord | `channel:<id>` or `user:<id>` |
| Slack | `user:U…` or `channel:C…` |
| WhatsApp | Phone number (`+15551234567`) |
| Matrix | `user:@user:server` |

---

## Presentation Payloads

A single `--presentation` JSON renders natively per platform:
- Discord → components
- Slack → blocks
- Telegram → inline buttons
- Teams → adaptive cards

---

## Platform-Specific Options

| Flag | Platform | Description |
|---|---|---|
| `--force-document` | Telegram | Suppress image compression |
| `--gif-playback` | WhatsApp | GIF playback mode |
| `--poll-duration-seconds` | Telegram | Poll duration |
| `--thread-id` | Telegram/Discord | Forum topic / thread |

---

## Our Channels

```bash
# Send to Telegram One Mind group (topic 292 = Legacy)
openclaw message send --channel telegram \
  --to "-1003905828404" --thread-id 292 \
  --text "Message from CLI"

# Send to Discord agent channel
openclaw message send --channel discord \
  --to "channel:<channel-id>" --text "Hello"
```

Full channel map: [`channels/CHANNELS.md`](../channels/CHANNELS.md)
