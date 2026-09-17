---
name: OpenClaw CLI — Config
updated: 2026-05-11
---

# `openclaw config`

Non-interactive editing of `openclaw.json` via get/set/patch/unset/validate. Use this to script config changes without opening the file directly.

**Official docs:** [docs.openclaw.ai/cli/config](https://docs.openclaw.ai/cli/config.md)

---

## Commands

```bash
# Read values
openclaw config get                           # Dump entire config as JSON
openclaw config get gateway.port              # Read a single key (dot-path)
openclaw config get agents.legacy.model       # Nested key lookup

# Write values
openclaw config set gateway.port 3000
openclaw config set agents.legacy.model "github-copilot/claude-sonnet-4-6"

# Patch (merge object at path)
openclaw config patch agents.legacy '{"thinking":"medium"}'

# Remove a key
openclaw config unset agents.legacy.thinking

# Validate config file
openclaw config validate                      # Check for schema errors
openclaw config validate --strict             # Also warns on unknown keys

# Show config file path
openclaw config path
```

---

## Dot-Path Syntax

Keys use dot notation to traverse nested objects:

```
gateway.port          → { "gateway": { "port": ... } }
agents.legacy.model   → { "agents": { "legacy": { "model": ... } } }
channels.telegram[0]  → first item in telegram array
```

---

## Config File

Active config lives at `/home/ubuntu/.openclaw/openclaw.json` on omos-core.

Never edit this file directly on the VPS while the gateway is running — changes require a gateway restart to take effect unless the key supports hot-reload. Use `openclaw config set` + `openclaw gateway reload` for safe live updates.

---

## Notes

- `config set` writes immediately but does **not** restart or reload the gateway automatically
- `config validate` is safe to run at any time — read-only
- String values containing spaces must be quoted in the shell: `openclaw config set gateway.host "my host"`
