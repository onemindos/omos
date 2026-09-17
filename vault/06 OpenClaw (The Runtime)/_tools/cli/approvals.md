---
name: OpenClaw CLI — Approvals
updated: 2026-05-11
---

# `openclaw approvals`

Manages execution approval policies across local, gateway, and node hosts.

**Official docs:** [docs.openclaw.ai/cli/approvals](https://docs.openclaw.ai/cli/approvals.md)

---

## Commands

```bash
# Get current policy
openclaw approvals get                        # Local policy
openclaw approvals get --gateway              # Gateway host policy
openclaw approvals get --node <identifier>    # Node host policy

# Apply policy from file
openclaw approvals set --file policy.json
openclaw approvals set --stdin               # Pipe JSON policy

# Allowlist management
openclaw approvals allowlist add "~/Projects/**/bin/rg"
openclaw approvals allowlist remove "~/Projects/**/bin/rg"

# Preset: unrestricted execution (YOLO mode)
openclaw exec-policy preset yolo             # Local-only shortcut
```

---

## Precedence Model

**Host approvals file is the enforceable source of truth** — overrides requested policies even when they try to broaden or narrow permissions. Node approvals combine with gateway config at runtime.

---

## YOLO Mode Config

```json
{
  "security": "full",
  "ask": "off"
}
```

Or use the shortcut: `openclaw exec-policy preset yolo`

---

## Related

- [`security.md`](security.md) — security audit
- [Official sandbox docs](https://docs.openclaw.ai/gateway/sandboxing.md)
