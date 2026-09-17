---
name: OpenClaw CLI — Sandbox
updated: 2026-05-11
---

# `openclaw sandbox`

Manages isolated runtime environments for agent execution.

**Official docs:** [docs.openclaw.ai/cli/sandbox](https://docs.openclaw.ai/cli/sandbox.md)

---

## Commands

```bash
openclaw sandbox explain                    # Inspect effective sandbox config (mode, scope, policies)
openclaw sandbox list                       # All runtimes with status, backend type, idle duration
openclaw sandbox recreate                   # Force regeneration of runtimes with updated settings
openclaw sandbox recreate --all             # Recreate all runtimes
openclaw sandbox recreate --agent <id>      # Recreate for specific agent
```

---

## Backends

| Backend | Description |
|---|---|
| Docker | Default — containerized execution |
| SSH | Remote workspace execution |
| OpenShell | OpenShell runtime |

---

## Key Behaviors

- **Existing runtimes keep old settings** until explicitly recreated — always run `recreate` after config changes
- Remote backend recreation deletes the canonical workspace; it gets reseeded on next use
- **Prefer `sandbox recreate` over manual cleanup** to avoid registry mismatches when scopes change

---

## Config Location

```
~/.openclaw/openclaw.json
  └─ agents.defaults.sandbox    # Default sandbox config
  └─ agents.list[].sandbox      # Per-agent overrides
```

---

## Common Workflow

```bash
# After updating Docker image or SSH credentials:
openclaw sandbox recreate --all

# Check what sandbox config is active:
openclaw sandbox explain
```
