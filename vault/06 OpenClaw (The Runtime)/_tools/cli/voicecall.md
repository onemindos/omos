---
name: OpenClaw CLI — Voice Call
updated: 2026-05-11
---

# `openclaw voicecall`

Manages phone calls through OpenClaw's voice-call plugin.

**Official docs:** [docs.openclaw.ai/cli/voicecall](https://docs.openclaw.ai/cli/voicecall.md)

---

## Commands

```bash
# Setup check
openclaw voicecall setup                      # Show readiness status

# Test (dry run or live)
openclaw voicecall smoke                      # Dry-run validation
openclaw voicecall smoke --to +15551234567 --yes  # Execute test call

# Active call management
openclaw voicecall call --to +15551234567 --message "Hello"
openclaw voicecall continue --message "Follow-up"
openclaw voicecall dtmf --digits "1234"       # Send touch-tone digits
openclaw voicecall end
openclaw voicecall status                     # Monitor active calls

# Webhook exposure (Tailscale)
openclaw voicecall expose --mode serve        # Local endpoint sharing
openclaw voicecall expose --mode funnel       # Public internet access
openclaw voicecall expose --mode off          # Disable
```

---

## Provider Requirements

For external carriers (Twilio, Telnyx, Plivo), a **publicly accessible webhook URL** is required. Local/private addresses won't work — carriers can't reach them.

Options:
- Tailscale Funnel (public)
- Tailscale Serve (LAN)
- Existing public URL

**Security:** Restrict webhook access to trusted networks. Prefer Serve over Funnel.

---

## Our Setup

Twilio integration documented in [`nodes/twilio.md`](../nodes/twilio.md)
