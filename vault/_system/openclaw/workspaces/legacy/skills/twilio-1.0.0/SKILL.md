# Twilio Skill — SMS & Voice Call

Use this skill to send SMS messages or initiate voice calls to Zeus via Twilio.

## When to use

| Tool         | When                                                                       |
| ------------ | -------------------------------------------------------------------------- |
| `sms_send`   | Urgent alerts, HITL approval requests, status updates when Discord is down |
| `voice_call` | Critical escalation — Zeus hasn't responded to SMS/Discord within timeout  |

## Escalation ladder (follow this order)

1. Discord message (native channel — always first)
2. Telegram (if Discord fails)
3. `sms_send` (if no ack within 5–10 min on a critical item)
4. `voice_call` (last resort — only for L1/L0 decisions requiring immediate response)

## Examples

```
# HITL approval pending, 10min no response
sms_send("Legacy needs your approval: [approve/deny] the $500 server upgrade. Respond here or on Discord.")

# Financial decision, Zeus unreachable on all channels
voice_call("Legacy is calling. Critical approval needed for your server infrastructure. Check Discord.")
```

## Env vars

```
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_FROM_NUMBER=+12015551234
ZEUS_PHONE_NUMBER=+13475559876
```
