# HITL Escalate Skill — Human-in-the-Loop Escalation

Use this skill to escalate decisions when Zeus has not responded to normal channels.

## Escalation Ladder

Always follow this order. Do NOT skip steps without clear reason.

```
Step 1 — Discord message (native channel, always first)
         ↓ no response after 5–15 min on HIGH items
Step 2 — hitl_escalate_sms  ← this skill
         ↓ no response after 30+ min on CRITICAL items
Step 3 — hitl_escalate_voice ← this skill (last resort only)
```

## Autonomy levels that require HITL

| Level               | Examples                           | Timeout before SMS    |
| ------------------- | ---------------------------------- | --------------------- |
| L1 — Financial      | Server purchases, invoices >$100   | 10 min                |
| L1 — External comms | Sending emails from Zeus's address | 15 min                |
| L0 — Legal / health | Anything irreversible              | Immediate SMS + voice |

## Examples

```
# Zeus hasn't approved a server renewal in 10 minutes
hitl_escalate_sms(
  decision="Approve $48 DigitalOcean renewal due tonight",
  context="Server goes offline at midnight if not renewed",
  options="Reply YES to approve or check Discord"
)

# Critical — Zeus hasn't responded to SMS after 30+ minutes
hitl_escalate_voice(
  decision="Approve irreversible database deletion that will run in 5 minutes"
)
```

## Important rules

- **Never call `hitl_escalate_voice` without first sending `hitl_escalate_sms`**
- **Never use for LOW or NORMAL priority items** — these do not need HITL
- **Never spam** — one SMS escalation, one voice escalation, then wait and log
- If Zeus still doesn't respond after voice: abort the action, log it, notify on next heartbeat

## Env vars (same as Twilio skill)

```
TWILIO_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_AUTH_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TWILIO_FROM_NUMBER=+12015551234
ZEUS_PHONE_NUMBER=+13475559876
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
