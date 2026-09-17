---
title: "Interface — Telegram"
type: reference
status: active
created: 2026-04-12
modified: 2026-05-13
owner: zeus
tags: [type/reference, domain/SO, domain/openclaw, interface/telegram]
ai_summary: "Telegram as an OpenClaw interface — live config, forum topic map, super-max features, alert routing. Full audit & repair 2026-05-13."
---

# Interface: Telegram

**Status:** ✅ Active — primary mobile interface
**Bot:** `@onemindos_bot`
**Group:** `@OneMindOS` (`-1003905828404`) — Forum mode (supergroup with topics)
**Node:** omos-core (`100.123.126.34`) — ubuntu user service, port 18789
**Bot token:** 1Password → `op://aaeyris5z44pdqspk277u2oppi/c2hsruwzaum6z4d7qjjyduadpq/credential`
**Authorized user:** Zeus Telegram ID `577573558`

## Live Status (2026-05-13 — Post Full Audit)

| Item | Status | Notes |
|------|--------|-------|
| Bot connected | ✅ | `@onemindos_bot` long-polling active |
| Agent topics | ✅ 7 active | Topics 1, 256, 292, 318, 319, 320, 321 — each routed to correct agent |
| Delivery-only topics | ✅ 4 | Topics 325–328 `ingest: false` — cron output only, no replies |
| groupAllowFrom | ✅ | Locked to Zeus ID `577573558` only |
| allowFrom | ✅ | Locked to Zeus ID `577573558` only |
| gateway.auth.token | ✅ Fixed | Was in config (rotation loop bug) — now removed. Token comes from `gateway.systemd.env` only |
| Cron Discord refs | ✅ Fixed | 4 cron prompts had Discord channel names — all changed to Telegram |
| Orphan topics | ✅ Cleaned | 7 orphan topic entries (322–324, 329–332) removed from config |
| Skills Discord refs | ✅ Fixed | signals, onemind-maintainer, om-dashboard skills updated |
| Command menu | ✅ | 117 skills configured, top 100 registered (Telegram limit) |
| Super-max config | ✅ Applied | All advanced features enabled — see section below |

## Known Issues (Minor)

**topic:328 skill-workshop log error** — During heartbeat cron runs (~hourly), skill-workshop tries to create a review directory named `-1003905828404:topic:328` (colons are invalid in paths). The reviewer is gracefully skipped — this does NOT block agent replies or any other functionality. It's a cosmetic log entry only.

---

## Live Config (as of 2026-05-13)

> **Summary:** Super-max config applied. All advanced Telegram features enabled. Cleaned to 7 active topics. Old skill-agent topics removed. Delivery-only topics set to `ingest: false`.

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "dmPolicy": "allowlist",
      "allowFrom": ["577573558"],
      "groupPolicy": "allowlist",
      "groupAllowFrom": ["577573558"],
      "replyToMode": "off",
      "historyLimit": 100,
      "contextVisibility": "thread",
      "linkPreview": true,
      "errorPolicy": "once",
      "reactionLevel": "extensive",
      "reactionNotifications": "all",
      "threadBindings": {
        "enabled": true,
        "spawnSessions": true,
        "defaultSpawnContext": "isolated",
        "idleHours": 24,
        "maxAgeHours": 168
      },
      "actions": {
        "editMessage": true,
        "deleteMessage": true,
        "createForumTopic": true,
        "sendPoll": true,
        "sendSticker": true
      },
      "autoTopicLabel": { "enabled": true },
      "healthMonitor": { "enabled": true },
      "heartbeat": { "showAlerts": true, "showOk": true, "useIndicator": true },
      "execApprovals": { "enabled": true, "target": "both", "approvers": ["577573558"] },
      "dm": { "threadReplies": "always" },
      "markdown": { "tables": "code" },
      "commands": { "native": "auto", "nativeSkills": "auto" },
      "groups": {
        "-1003905828404": {
          "requireMention": false,
          "topics": {
            "1":   { "agentId": "main" },
            "256": { "agentId": "sentinel" },
            "292": { "agentId": "main" },
            "318": { "agentId": "zenith" },
            "319": { "agentId": "unity" },
            "320": { "agentId": "eternal" },
            "321": { "agentId": "sovereignty" },
            "325": { "ingest": false },
            "326": { "ingest": false },
            "327": { "ingest": false },
            "328": { "ingest": false }
          }
        }
      }
    }
  }
}
```

---

## Agent → Topic Map (Current)

| Agent | Topic Name | Topic ID | Notes |
|-------|-----------|----------|-------|
| Legacy (main) | General | `1` | Group general chat |
| Sentinel | 🛡 Sentinel | `256` | Monitoring |
| Legacy (main) | ⚡ Legacy | `292` | Primary orchestrator |
| Zenith | 🏔 Zenith | `318` | HP domain |
| Unity | 🌿 Unity | `319` | LE domain |
| Eternal | 👑 Eternal | `320` | GE domain |
| Sovereignty | 🌐 Sovereignty | `321` | SO domain |

### Delivery-Only Topics (cron output, no replies)

| Topic ID | Content |
|---------|---------|
| `325` | Daily domain reports |
| `326` | Weekly reports |
| `327` | OneMind status |
| `328` | Agent heartbeats + infra-check |

---

## Super-Max Features (Applied 2026-05-13)

Plain-language explanation of every feature enabled:

### threadBindings
Each Telegram topic gets its own memory slot. Messages in Legacy's topic (292) are remembered separately from Zenith's topic (318). Sessions expire after 24 hours idle (7 days max age). **Before:** Every message was a blank slate.

### reactionLevel: extensive
The bot reacts to your messages with emojis to confirm receipt. `extensive` = reacts to most message types, not just commands.

### execApprovals target: both
When an agent wants to run something risky (shell command, code, file edit), the approval buttons (✅ / ❌) appear in **both** your DM **and** the topic where the agent is working. **Before:** Only appeared in DM (easy to miss).

### actions (all enabled)
Agents can now: edit their own messages, delete their own messages, create new forum topics, send polls, send stickers.

### autoTopicLabel
Agents automatically rename forum topics based on the conversation context.

### healthMonitor
Gateway watches Telegram polling health every 5 minutes. Alerts if connection drops, confirms when it recovers.

### heartbeat
Agents report status to topic 328 on schedule. Alerts when an agent goes offline or recovers. Uses ✅/⚠️ indicators.

### replyToMode: off
Changed from `first`. Removes forced quote-reply structure — agents post replies directly. Fixes tool-progress preview in topics.

### nativeSkills: auto
Every installed skill is automatically a Telegram slash command. Type `/` in any topic to see all available skills.

### ingest: false on topics 325–328
These are delivery-only. Cron jobs post reports here, but your messages in these topics are ignored. No accidental agent invocations.

### errorPolicy: once
If an agent fails, it reports the error once and stops — no flood of repeated error messages.

### markdown.tables: code
Tables are wrapped in code blocks (Telegram can't render markdown tables natively). Keeps them readable.

### linkPreview: true
When an agent sends a URL, Telegram shows a card preview.

### dm.threadReplies: always
In DMs, agents always reply in a thread/reply chain. Keeps conversations organized.

---

## Capabilities

- **Thread isolation:** Each topic = isolated session with memory (via threadBindings)
- **Inline buttons:** Active for all chats — approve/reject flows work natively
- **Ack reaction:** ⚡ / 👀 appears on your message while agent processes
- **Voice notes (inbound):** Received as OGG/Opus → auto-transcribed via Deepgram → sent to agent as text
- **File uploads:** Sent to agent for processing
- **Native commands:** 100 slash commands auto-registered (Telegram limit — 117 configured)
- **Exec approvals:** Sent to DM + topic with ✅ / ❌ buttons (30-min expiry)
- **History:** 100 messages of group context loaded per turn

---

## Voice

Telegram Bot API does not support live voice calls. Voice is implemented as async voice note exchange:

```
Zeus sends voice note (OGG/Opus)
    → OpenClaw receives file
    → Deepgram transcribes (nova-3, real-time, already active)
    → Agent processes as text
    → Agent response → Azure Neural TTS → OGG/Opus file
    → Bot sends voice note reply
```

**Inbound:** Fully active — voice notes already transcribed by Deepgram on omos-core.

**Outbound TTS:** Pending activation. Config block:

```json
"channels": {
  "telegram": {
    "voice": {
      "tts": {
        "enabled": true,
        "provider": "azure",
        "voice": "en-US-AndrewMultilingualNeural",
        "format": "ogg-opus",
        "azure": {
          "endpoint": "onemind-openai.openai.azure.com",
          "region": "eastus2"
        }
      },
      "replyWithVoice": "always"
    }
  }
}
```

**Voice name options (Azure Neural):** `en-US-AndrewMultilingualNeural` (default), `en-US-BrianMultilingualNeural`, `en-US-EmmaMultilingualNeural`. All support multilingual output.

**Cost:** Azure Neural TTS ~$0.011/min of speech — negligible for personal use.

**Hard limit:** Live bidirectional voice (like a phone call) is not possible via Bot API. Use Discord voice channel or Talk mode (browser) for live duplex.

---

## Alert Routing

Telegram is the primary alert surface — always fires first.

| Urgency | Channel |
|---------|---------|
| Low / medium | Telegram only |
| High | Telegram + Discord |
| Critical / emergency | Telegram + Discord + TTS speaker + Twilio call |

Agents subscribe to `zeus.signal.*` subjects. Legacy dispatches based on urgency threshold.

---

## Owner Commands

Available in DM or any topic (gated by `commands.ownerAllowFrom`):

```
/diagnostics       — system health
/config set        — update config live
/session           — inspect sessions
/export-trajectory — export session history
/reasoning stream  — enable live reasoning stream
```

---

## Sentinel Fallback

`@omos_sentinel_bot` runs on omos-sentinel (`100.66.235.38`). It only activates in:
- Direct messages to the bot
- When omos-core is unreachable

Do not add sentinel to the group topics — it's a rescue channel only.

---

## Related

- [[CHANNELS]] — Full channel map for both platforms
- [[discord]] — Discord interface config
- [[MESSAGING]] — Bot registry
- [[omos-sentinel]] — Rescue/failover node
