---
title: "Interface — Twilio (Voice + SMS)"
type: reference
status: active
created: 2026-04-12
owner: zeus
tags: [type/reference, domain/SO, domain/openclaw, interface/twilio]
ai_summary: "Twilio as an OpenClaw interface — outbound calls, SMS fallback, emergency escalation, TTS."
---

# Interface: Twilio (Voice + SMS)

**Status:** ✅ Active
**Type:** Voice + SMS
**Primary use:** Emergency alerts, scheduled calls, voice briefings while driving

---

## Configuration

```yaml
# OpenClaw config.yaml — Twilio connector
connectors:
  - type: twilio
    account_sid: $TWILIO_ACCOUNT_SID   # 1Password
    auth_token: $TWILIO_AUTH_TOKEN     # 1Password
    from_number: "+1XXXXXXXXXX"        # 1Password → SO vault → Twilio Numbers
    to_number: "+1XXXXXXXXXX"          # Zeus's phone — 1Password
```

Credentials in **1Password → 00-24 SO vault → Twilio**.

---

## Capabilities

- **Outbound calls** — Legacy calls Zeus's iPhone and reads the message via TTS
- **SMS** — fallback when internet/Telegram is unavailable (160 char hard limit)
- **Emergency escalation** — fires automatically if no Telegram response in 5 min on emergency signal
- **Driving mode** — Twilio call is the preferred channel when Zeus is in a vehicle (hands-free)

---

## Making a Call

```python
from twilio.rest import Client

client = Client(TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN)

call = client.calls.create(
    to=ZEUS_PHONE,
    from_=TWILIO_FROM,
    twiml=f'<Response><Say voice="Polly.Matthew">{message}</Say></Response>'
)
```

Use voice `Polly.Matthew` — clear American English, natural pacing. Keep messages under 100 words.

---

## SMS Fallback

```python
msg = client.messages.create(
    to=ZEUS_PHONE,
    from_=TWILIO_FROM,
    body="[ALERT] " + message_text[:155]  # leave room for [ALERT] prefix
)
```

---

## When Twilio Fires

| Condition | Action |
|-----------|--------|
| Emergency-level signal | Immediate call |
| No Telegram response in 5 min (emergency) | Escalate to call |
| Zeus is driving (Bouncie speed > 5 mph) | Call preferred over Telegram for critical+ |
| Telegram unavailable / SMS-only scenario | SMS fallback |

---

## Notes

- Never call during quiet hours (22:00–07:00) unless `emergency` level
- TTS message should be complete in itself — no links, no markdown, no emojis
- Twilio call SID logged to `ZEUS-SIGNAL-LOG.md` for delivery audit


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
