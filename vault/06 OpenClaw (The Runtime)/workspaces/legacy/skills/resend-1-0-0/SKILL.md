# Resend Skill — Email Delivery

Use this skill to send HTML emails to Zeus via Resend.

## When to use

Email is lower-urgency than Discord/Telegram. Use it for:

- Daily/weekly digest summaries
- Reports with structured data
- Non-urgent notifications Zeus should review later
- Anything with links or formatted content better suited to email

**Do not use** email for anything requiring immediate response — use Discord or `sms_send` instead.

## Priority routing guide

| Priority   | Channel order              |
| ---------- | -------------------------- |
| `critical` | Discord → SMS → Voice call |
| `high`     | Discord → Telegram → SMS   |
| `normal`   | Discord → Telegram         |
| `low`      | Email (digest)             |

## Examples

```
# Weekly HP summary
email_send(
  subject="Weekly Health Report",
  body="Score: 84/100. Sleep avg: 7.2h. Workouts: 4/5. Focus time: 31h.",
  priority="normal"
)

# Non-urgent invoice notification
email_send(
  subject="Invoice #1042 ready",
  body="Your invoice from DigitalOcean for $48 is ready to review.",
  action_url="https://cloud.digitalocean.com/billing",
  priority="low"
)
```

## Env vars

```
RESEND_API_KEY=re_xxxxxxxxxxxxxxxxxxxxxxxx
ZEUS_EMAIL=zeus@onemindos.dev
RESEND_FROM_EMAIL=Legacy <notifications@onemindos.dev>
```


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
