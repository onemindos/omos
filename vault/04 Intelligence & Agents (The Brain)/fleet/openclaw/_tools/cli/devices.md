---
name: OpenClaw CLI — Devices
updated: 2026-05-11
---

# `openclaw devices`

Manages paired device connections — approve pairing requests, revoke access, rotate tokens.

**Official docs:** [docs.openclaw.ai/cli/devices](https://docs.openclaw.ai/cli/devices.md)

> Note: `openclaw devices` manages **paired clients** (phones, desktops). For remote OpenClaw nodes (VPS agents), see [`nodes.md`](nodes.md).

---

## Commands

```bash
# List devices
openclaw devices list                         # All paired devices
openclaw devices list --pending               # Awaiting approval only
openclaw devices list --json

# Approve / deny pairing requests
openclaw devices approve <requestId>          # Approve a pending pairing
openclaw devices deny <requestId>             # Reject a pairing request

# Revoke access
openclaw devices remove <deviceId>            # Remove a paired device
openclaw devices revoke --all                 # Revoke all device tokens

# Token rotation
openclaw devices rotate <deviceId>            # Issue new token, invalidate old
openclaw devices rotate --all                 # Rotate all device tokens

# Inspect a device
openclaw devices inspect <deviceId>           # Show capabilities, last seen, token age
```

---

## Pairing Flow

1. On mobile: scan QR code from `openclaw qr`
2. Bootstrap token creates a pending pairing request
3. On gateway host: `openclaw devices list --pending`
4. Approve: `openclaw devices approve <requestId>`
5. Device now has a full paired token and can interact with the gateway

---

## Device Object Fields

| Field | Description |
|---|---|
| `id` | Unique device ID |
| `name` | User-provided label (set on mobile) |
| `platform` | `ios`, `android`, `desktop`, `cli` |
| `lastSeen` | ISO timestamp of last activity |
| `tokenAge` | Days since last token rotation |
| `status` | `pending`, `approved`, `revoked` |

---

## Notes

- Token rotation (`rotate`) issues a new token and immediately invalidates the old one — the device must re-authenticate
- `revoke --all` is a security response action — all devices must re-pair after this
- `devices remove` soft-deletes the record; historical session data is preserved
