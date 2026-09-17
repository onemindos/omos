---
name: OpenClaw — Reactions
updated: 2026-05-11
---

# Reactions

The agent can add and remove emoji reactions on messages using the `message` tool with the `react` action. Lightweight async visibility into what agents are doing — no full reply needed.

**Official docs:** [docs.openclaw.ai/tools/reactions](https://docs.openclaw.ai/tools/reactions.md)

---

## How It Works

```json
{
  "action": "react",
  "messageId": "msg-123",
  "emoji": "thumbsup"
}
```

- `emoji` is required when adding a reaction
- `emoji: ""` (empty string) removes all of the bot's reactions on the message
- `remove: true` removes a specific emoji (requires non-empty `emoji`)
- `trackToolCalls: true` lets the runtime use the reacted message for subsequent tool-progress reactions during the same turn

---

## Channel Behavior

| Channel | Notes |
|---|---|
| Discord | Empty emoji removes bot's reactions. `remove: true` removes specific emoji |
| Telegram | Empty emoji removes bot reaction. `remove: true` maps to empty emoji internally |
| WhatsApp | Requires non-empty emoji. `remove: true` removes that specific emoji reaction |
| Signal | Reaction notifications controlled by `channels.signal.reactionNotifications`: `off`, `own` (default), or `all` |
| Feishu | Use the separate `feishu_reaction` tool with `add`/`remove`/`list` actions |

---

## Reaction Level Config

Per-channel `reactionLevel` controls how broadly the agent uses reactions:

| Value | Behavior |
|---|---|
| `off` | Agent never reacts |
| `ack` | React to acknowledge receipt |
| `minimal` | React only on key events |
| `extensive` | React broadly throughout task |

```json
{
  "channels": {
    "telegram": { "reactionLevel": "minimal" },
    "discord": { "reactionLevel": "ack" }
  }
}
```

---

## One Mind Use Case

Zeus sends a task to Legacy via Telegram. Legacy reacts with 🕐 when it starts working, then ✅ when complete — or ❌ if it fails. Zero message overhead for simple status tracking.

Useful for the 19 agents firing daily cron reports — a quick ✅ reaction on the cron trigger message confirms the agent ran without requiring Zeus to check the full output.
