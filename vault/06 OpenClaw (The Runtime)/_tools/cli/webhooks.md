---
name: OpenClaw CLI — Webhooks
updated: 2026-05-11
---

# `openclaw webhooks`

Webhook integration tools — currently focused on Gmail Pub/Sub.

**Official docs:** [docs.openclaw.ai/cli/webhooks](https://docs.openclaw.ai/cli/webhooks.md)

---

## Commands

### Gmail Webhooks
```bash
# Configure Gmail watch + Pub/Sub + webhook delivery
openclaw webhooks gmail setup --account you@example.com

# Start the Gmail watch service with auto-renewal
openclaw webhooks gmail run

# With JSON output
openclaw webhooks gmail setup --account you@example.com --json

# Custom host/port
openclaw webhooks gmail run --host 0.0.0.0 --port 3000
```

---

## Gmail Setup Options

| Option | Description |
|---|---|
| Tailscale funnel | Public internet exposure |
| Tailscale serve | Local network only |
| Off | No Tailscale, manual webhook URL |
| Custom topic/subscription | Override Pub/Sub resource names |
| Label filtering | Watch specific Gmail labels |

---

## Related

- [automation/README.md](../automation/README.md) — Cron + hooks + webhooks overview
- [Official webhooks plugin](https://docs.openclaw.ai/plugins/webhooks.md)
