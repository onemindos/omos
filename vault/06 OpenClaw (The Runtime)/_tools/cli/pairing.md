---
name: OpenClaw CLI — Pairing
updated: 2026-05-11
---

# `openclaw pairing`

Manages approval and inspection of direct message pairing requests across channels.

**Official docs:** [docs.openclaw.ai/cli/pairing](https://docs.openclaw.ai/cli/pairing.md)

---

## Commands

```bash
openclaw pairing list --channel telegram       # List pending pairing requests
openclaw pairing list --channel discord --json

openclaw pairing approve <code> --channel telegram         # Approve a pairing code
openclaw pairing approve <code>                            # Auto-resolves if only one pairing channel

# With account (multi-account channels)
openclaw pairing list --channel telegram --account ops
openclaw pairing approve <code> --channel telegram --notify
```

---

## Bootstrap Behavior

On the **first approval**, if `commands.ownerAllowFrom` is empty, the approved sender automatically becomes the command owner using a channel-scoped entry like `telegram:123456789`. This only applies to the first approval — subsequent approvals require an existing owner.

Once approved, the owner can execute sensitive operations like diagnostics and config exports.

---

## Multiple Channels

If multiple pairing-capable channels are configured, you **must** specify `--channel` — there's no auto-selection.

---

## Related

- [`qr.md`](qr.md) — generate QR codes for mobile pairing
- [`devices.md`](devices.md) — device pairing and token management
- [`nodes.md`](nodes.md) — node pairing approval
