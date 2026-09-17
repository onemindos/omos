---
name: OpenClaw CLI — QR
updated: 2026-05-11
---

# `openclaw qr`

Generates mobile pairing QR codes and setup codes from Gateway configuration.

**Official docs:** [docs.openclaw.ai/cli/qr](https://docs.openclaw.ai/cli/qr.md)

---

## Commands

```bash
openclaw qr                          # Generate QR code + setup code
openclaw qr --remote                 # Use remote gateway URL if available
openclaw qr --url wss://my-gateway   # Override gateway URL
openclaw qr --public-url <url>       # Override public-facing URL
openclaw qr --setup-code-only        # Output setup code text only
openclaw qr --json                   # Structured output: setup code + gateway URL + auth
openclaw qr --token <tok>            # Specify auth token
openclaw qr --password <pw>          # Specify auth password (mutually exclusive with --token)
```

---

## After Scanning

The QR contains a **temporary bootstrap token**, not your shared credentials. After scanning, you must approve the pairing:

```bash
openclaw devices list
openclaw devices approve <requestId>
```

---

## Requirements & Constraints

- Requires a Gateway supporting `secrets.resolve` method
- Use `wss://` or Tailscale Serve/Funnel — `ws://` over public Tailscale will fail
- Private LAN and `.local` hosts can use `ws://` safely
- `--token` and `--password` are mutually exclusive

---

## Legacy Alias

`openclaw clawbot qr` is identical — maintained for backward compatibility only.
