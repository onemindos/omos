---
name: OpenClaw CLI — Doctor
updated: 2026-05-11
---

# `openclaw doctor`

Diagnoses and auto-fixes gateway and channel configuration issues.

**Official docs:** [docs.openclaw.ai/cli/doctor](https://docs.openclaw.ai/cli/doctor.md) · [gateway/doctor](https://docs.openclaw.ai/gateway/doctor.md)

---

## Commands

```bash
openclaw doctor                       # Run health checks
openclaw doctor --fix                 # Apply recommended fixes automatically
openclaw doctor --repair              # Alias for --fix
openclaw doctor --deep                # Comprehensive system service scanning
openclaw doctor --non-interactive     # No prompts — safe migrations only
openclaw doctor --generate-gateway-token   # Create + configure gateway auth token
openclaw doctor --force               # Aggressive repairs including config overwrites
```

---

## What Doctor Checks

- Gateway config validity
- Channel configuration
- Service installation state
- Authentication tokens
- Registry drift and legacy config
- Plugin load issues (via `openclaw plugins doctor`)

---

## What `--fix` Repairs

- Normalizes legacy job formats (cron)
- Repairs registry drift
- Fixes legacy configuration patterns

What it does NOT touch: credentials, tool disabling, plugin rewrites.

---

## Channel-Specific Checks

```bash
# Check Discord channel permissions
openclaw channels capabilities --channel discord --target channel:<channel-id>

# Live probe all channels
openclaw channels status --probe
```

---

## macOS Note

If you've previously set gateway env vars via `launchctl setenv`, remove them:

```bash
launchctl unsetenv OPENCLAW_GATEWAY_TOKEN
launchctl unsetenv OPENCLAW_GATEWAY_PASSWORD
```

---

## Nix Mode

In `OPENCLAW_NIX_MODE=1`, repair functionality is disabled — config files are immutable.
