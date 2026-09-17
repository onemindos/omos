# Telegram Bot API Reference

## Base URL
```
https://api.telegram.org/bot{token}/{method}
```

## Methods

| Method | Description | Key Params |
|--------|-------------|------------|
| `getMe` | Bot info | - |
| `sendMessage` | Send text | `chat_id`, `text`, `parse_mode` |
| `sendPhoto` | Send image | `chat_id`, `photo`, `caption` |
| `sendDocument` | Send file | `chat_id`, `document` |
| `sendLocation` | Send coords | `chat_id`, `latitude`, `longitude` |
| `getUpdates` | Poll messages | `offset`, `timeout` |
| `setWebhook` | Set webhook | `url`, `secret_token` |
| `deleteWebhook` | Remove webhook | - |
| `editMessageText` | Edit msg | `chat_id`, `message_id`, `text` |
| `deleteMessage` | Delete msg | `chat_id`, `message_id` |
| `setMyCommands` | Set bot cmds | `commands[]` |

## Parse Modes

| Mode | Bold | Italic | Code |
|------|------|--------|------|
| `Markdown` | `*text*` | `_text_` | `` `code` `` |
| `HTML` | `<b>text</b>` | `<i>text</i>` | `<code>code</code>` |

## Chat ID Types

- **User**: Positive integer
- **Group**: Negative integer
- **Supergroup/Channel**: Starts with `-100`

## Inline Keyboard

```json
{
  "reply_markup": {
    "inline_keyboard": [
      [{"text": "Button", "callback_data": "action"}],
      [{"text": "Link", "url": "https://..."}]
    ]
  }
}
```

## Rate Limits

- ~30 messages/second per chat
- ~20 messages/minute to same user


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
