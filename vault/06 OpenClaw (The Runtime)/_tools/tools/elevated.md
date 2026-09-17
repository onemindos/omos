---
name: OpenClaw — Elevated Mode
updated: 2026-05-11
---

# Elevated Mode

When an agent runs inside a sandbox, its `exec` commands are confined to the sandbox. **Elevated mode** lets the agent break out and run commands on the host system, with configurable approval gates.

**Official docs:** [docs.openclaw.ai/tools/elevated](https://docs.openclaw.ai/tools/elevated.md)

> Elevated mode only changes behavior when the agent is **sandboxed**. For unsandboxed agents, exec already runs on the host.

---

## Directives

```bash
/elevated on      # Run outside sandbox, keep approvals
/elevated ask     # Same as on (alias)
/elevated full    # Run outside sandbox, SKIP approvals (break-glass)
/elevated off     # Return to sandbox-confined execution
/elevated         # Show current level (no argument)
/elev on|off|ask|full  # Short alias
```

---

## How It Works

**Step 1 — Enable in config + sender allowlist:**

```json
{
  "tools": {
    "elevated": {
      "enabled": true,
      "allowFrom": {
        "discord": ["user-id-123"],
        "telegram": ["+15555550123"]
      }
    }
  }
}
```

**Step 2 — Set the level:**

Session default (sticky for the conversation):
```
/elevated full
```

Inline (applies to that message only):
```
/elevated on run the deployment script
```

**Step 3 — Commands run outside the sandbox** with `exec`.

In `full` mode, exec approvals are skipped. In `on`/`ask` mode, configured approval rules still apply.

---

## Resolution Order

1. Inline directive on the message (applies to that message only)
2. Session override (set by sending a directive-only message)
3. Global default (`agents.defaults.elevatedDefault` in config)

---

## Four-Gate Authorization

All four gates must pass — failing any one means elevated is unavailable:

| Gate | Config key |
|---|---|
| Global enable | `tools.elevated.enabled` |
| Sender allowlist | `tools.elevated.allowFrom` |
| Per-agent restriction | `agents.list[].tools.elevated.enabled` |
| Per-agent allowlist | `agents.list[].tools.elevated.allowFrom` |

Allowlist entry formats: bare ID, `name:`, `username:`, `tag:`, `id:`, `from:`, `e164:`

---

## What Elevated Does NOT Control

- **Tool policy** — if `exec` is denied by tool policy, elevated cannot override it
- **Host selection** — elevated uses the configured/session exec target rules, not a free override
- **`/exec` directive** — adjusts per-session exec defaults separately, does not require elevated mode
- **`!` bash chat commands** — require both `tools.elevated.enabled` AND `tools.bash.enabled`

---

## One Mind Use Case

The fleet runs on Lightsail (omos-core). Elevated mode is how Sovereignty and Sentinel perform legitimate host-level operations — deploy scripts, manage systemd services, write to paths outside the sandbox — without dismantling the security boundary entirely.

The tiered allowlist ensures only Zeus (authorized sender) on authorized channels can trigger host-level execution.

```json
{
  "agents": {
    "list": [
      {
        "id": "sovereignty",
        "tools": {
          "elevated": {
            "enabled": true,
            "allowFrom": {
              "telegram": ["+1XXXXXXXXXX"],
              "discord": ["ZEUS_DISCORD_ID"]
            }
          }
        }
      }
    ]
  }
}
```

---

## Related

- [tools/exec.md](exec.md) — shell execution
- [cli/sandbox.md](../cli/sandbox.md) — sandbox configuration
- [gateway/sandboxing](../gateway/README.md)
