---
name: signals
description: "Legacy's interrupt and alert layer — watches for conditions, evaluates urgency, and routes alerts to Zeus via Telegram, Discord, home speaker TTS, or Twilio call using FIND ZEUS MODE. Use when setting up a watch, routing an urgent event, or triggering a voice alert. Activates on: alert me when, watch for, notify me if, signal, urgent, speaker, TTS, FIND ZEUS."
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Signals Skill

> Alert Zeus of urgent events. Route to the right channel. Speak when needed.

## What This Is

The signals skill is Legacy's interrupt layer. It watches for conditions, evaluates urgency, and delivers alerts to Zeus through the appropriate channel — including voice (home speakers, Twilio call) when Zeus is reachable that way.

It integrates directly with FIND ZEUS MODE: it knows where Zeus is and routes accordingly.

---

## Trigger Patterns

### 1. Zeus Requests a Watch

```
"Watch for X, notify me when it happens"
"Alert me if disk hits 90%"
"Tell me when the deploy is done"
"Notify me if Legacy crashes"
"If I haven't worked out by 7pm, remind me"
```

### 2. Agent Publishes an Urgent Signal

Any OpenClaw agent can publish to `zeus.signal.*` subjects, which Legacy processes through the awareness filter.

```bash
nats pub zeus.signal.urgent '{
  "level": "critical",
  "message": "Server disk at 94% — action needed",
  "source": "guardian",
  "speak": true
}'
```

### 3. System Events (Auto-Watched)

These are always watched without Zeus needing to configure them:

| Event                            | Urgency  | Default Channel    |
| -------------------------------- | -------- | ------------------ |
| Any OpenClaw agent crash         | critical | Telegram + speaker |
| VPS disk > 90%                   | high     | Telegram           |
| VPS CPU > 95% for 5min           | high     | Telegram           |
| N8N workflow failure             | medium   | Telegram           |
| Service restart (>3x in 1hr)     | high     | Telegram           |
| New auth failure on VPS          | medium   | Telegram           |
| Calendar event starting in 5 min | low      | Telegram           |

---

## Urgency Levels

| Level         | Description                    | Channels Used                              |
| ------------- | ------------------------------ | ------------------------------------------ |
| **info**      | Log only — does not alert Zeus | None (silent)                              |
| **low**       | Non-urgent awareness           | Telegram (batched, 30min delay)            |
| **medium**    | Notice required                | Telegram (immediate)                       |
| **high**      | Action likely needed           | Telegram + Discord                         |
| **critical**  | Action required now            | Telegram + Discord + speaker TTS (if home) |
| **emergency** | Immediate attention            | All channels + Twilio call                 |

---

## FIND ZEUS MODE Integration

Signals always check Zeus's current location before routing:

```
Signal fires (any urgency)
    ↓
Read zeus.presence.zone (NATS last value)
    ↓
Apply zone routing rules:
    Home Greenwich:  can speak, Telegram primary
    NYC / Office:    Telegram primary, no speakers
    Driving:         Twilio call if critical+, otherwise Telegram
    Unknown:         Telegram only, raise threshold
    ↓
Apply quiet hours filter (22:00–07:00 at home → critical+ only)
    ↓
Deliver via routed channel
```

---

## Speaker TTS (Voice Alerts)

When Zeus is home and signal is critical+, Legacy speaks the alert through home speakers.

### Via Home Assistant

```bash
# From Legacy workspace or signaling agent:
POST $HA_URL/api/services/tts/google_translate_say
Authorization: Bearer $HA_TOKEN
{
  "entity_id": "media_player.home_speaker",
  "message": "Zeus, your server disk is at ninety-four percent. Immediate action required.",
  "language": "en-US"
}
```

### Via NATS (Preferred — Decoupled)

Publish to `zeus.alert.speak` and let the presence-pulse service handle delivery:

```bash
nats pub zeus.alert.speak '{
  "text": "Zeus, your server disk is at ninety-four percent.",
  "urgency": "critical",
  "timestamp": "2026-04-12T09:00:00Z"
}'
```

The VPS `speaker-bridge` service subscribes to `zeus.alert.speak` and calls HA TTS or Twilio depending on Zeus's current zone.

### TTS Rules

- Messages spoken at home: use HA TTS → `media_player.home_speaker`
- Messages while driving: use Twilio TTS call (only critical+)
- Messages during quiet hours: suppressed unless `emergency` level
- Max speech length: 100 words — Legacy truncates if longer
- Voice: calm English, no jargon, clear action if action needed

---

## Event Watch Registry

Zeus can register custom watches. These persist in `_codex/system/ZEUS-WATCHES.md`.

### Watch Format

```yaml
# ZEUS-WATCHES.md
watches:
  - id: deploy-complete
    description: "Notify when onemind-web deploy finishes"
    source: nats
    subject: "deploy.onemind.complete"
    condition: any
    urgency: medium
    message: "Deploy complete. onemind-web is live."
    expires: "2026-04-15" # one-time watch expires

  - id: workout-reminder
    description: "Remind Zeus to work out if not done by 7pm"
    source: calendar
    condition: "time=19:00 AND workout_logged=false"
    urgency: low
    message: "No workout logged yet today. Still time."
    recurring: daily

  - id: revenue-milestone
    description: "Alert when MRR crosses $1000"
    source: nats
    subject: "empire.revenue.mrr"
    condition: "value > 1000"
    urgency: high
    message: "MRR crossed $1000. That's the threshold."
    one_shot: true
```

### Registering a Watch (Via Chat)

Zeus can tell Legacy:

```
"Watch for X" → Legacy creates entry in ZEUS-WATCHES.md, confirms
"Remove the deploy watch" → Legacy removes that entry
"Show my active watches" → Legacy lists ZEUS-WATCHES.md entries
"Clear all watches" → Legacy deletes all non-system watches (confirms first)
```

---

## Cooldown Logic

To prevent alert spam, the signals skill maintains cooldown state:

```python
COOLDOWN_MINUTES = {
    "info":      0,    # no cooldown
    "low":       60,   # max 1 per hour
    "medium":    15,   # max 1 per 15 min per event type
    "high":      5,    # max 1 per 5 min per event type
    "critical":  2,    # rarely suppressed
    "emergency": 0,    # never suppressed
}

# Key = (event_type, source) — cooldowns are per-event-type
# If same event fires again within cooldown, it's dropped (not delivered)
# After 3x drops in 30 min for same event → escalate to next level
```

---

## NATS Subjects Summary

```
zeus.signal.urgent      → High-priority interrupt (any agent can publish here)
zeus.signal.event       → General event notification
zeus.signal.watch       → Watch condition fired
zeus.signal.error       → System error
zeus.alert.speak        → Trigger speaker TTS (text + urgency payload)
zeus.alert.call         → Trigger Twilio voice call
zeus.signal.status      → Legacy publishes current signals state (heartbeat)
```

---

## Execution Flow

```
1. Signal arrives (NATS publish OR Legacy detects event)
2. Legacy extracts: source, level, message, speak flag
3. Legacy reads zeus.presence.zone (last value cache)
4. Awareness filter: should_alert(signal, zeus_state) → {alert, channel, speak}
5. If alert=false: log to ZEUS-SIGNAL-LOG.md, done
6. If alert=true:
   a. Format message for channel (brief for voice/call, rich for Telegram)
   b. Deliver via channel
   c. If speak=true: publish zeus.alert.speak
   d. Log to ZEUS-SIGNAL-LOG.md with delivery receipt
7. Start cooldown timer for (event_type, source)
```

---

## Files and State

| File                      | Location          | Purpose                        |
| ------------------------- | ----------------- | ------------------------------ |
| `ZEUS-WATCHES.md`         | `_codex/system/`  | Registry of active watches     |
| `ZEUS-SIGNAL-LOG.md`      | `_codex/system/`  | Delivery history (last 7 days) |
| `FIND-ZEUS-MODE.md`       | `SO/03 Agents/`   | Presence architecture          |
| `home-assistant/SKILL.md` | `.claude/skills/` | HA API integration             |
| `telegram/SKILL.md`       | `.claude/skills/` | Telegram delivery              |

---

## Environment Variables Required

```bash
HOME_ASSISTANT_URL=https://homeassistant.local:8123
HOME_ASSISTANT_TOKEN=<long-lived-access-token>   # from 1Password SO vault
TWILIO_ACCOUNT_SID=<sid>                          # 1Password
TWILIO_AUTH_TOKEN=<token>                         # 1Password
TWILIO_FROM_NUMBER=+1XXXXXXXXXX                  # 1Password
ZEUS_PHONE_NUMBER=+1XXXXXXXXXX                   # 1Password
NATS_URL=nats://localhost:4222
```

---

## Usage Examples

```
Zeus: "Alert me when the VPS memory hits 80%"
Legacy: "Watch registered. I'll alert you when VPS memory exceeds 80%. Urgency: high."

Zeus: "Critical alert me if OpenClaw crashes"
Legacy: "Watch registered. Critical alert + speaker TTS if OpenClaw service stops."

Zeus: "Show my watches"
Legacy: lists ZEUS-WATCHES.md entries

Zeus: "Test signal — send me a medium alert"
Legacy: [sends Telegram test message] "Test delivered to Telegram. You're in Home Greenwich zone."
```

---

_Part of SO — Sovereign Operations. Lives at `.claude/skills/signals/SKILL.md`._


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
