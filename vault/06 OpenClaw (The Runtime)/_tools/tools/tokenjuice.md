---
name: OpenClaw — Tokenjuice
updated: 2026-05-11
---

# Tokenjuice

Optional bundled plugin that compacts noisy `exec` and `bash` tool results **after** the command has already run. Reduces tool output volume before it re-enters the session context — keeps long sessions lean without losing accuracy.

**Official docs:** [docs.openclaw.ai/tools/tokenjuice](https://docs.openclaw.ai/tools/tokenjuice.md)

> Tokenjuice changes the returned `tool_result`, not the command itself. It does not rewrite shell input, rerun commands, or change exit codes.

---

## Enable

```bash
openclaw config set plugins.entries.tokenjuice.enabled true
# or
openclaw plugins enable tokenjuice
```

No separate install step — the plugin ships with OpenClaw.

Config equivalent:
```json
{
  "plugins": {
    "entries": {
      "tokenjuice": { "enabled": true }
    }
  }
}
```

---

## What It Changes

- Compacts noisy `exec` and `bash` results before they feed back into the session
- Keeps the original command execution untouched
- **Preserves** exact file-content reads and other commands that should stay raw
- Stays opt-in — disable if you want verbatim output everywhere

---

## Verify It's Working

1. Enable the plugin
2. Start a session that can call `exec`
3. Run a noisy command: `git status` or `ls -la`
4. Check that the returned tool result is shorter and more structured than raw shell output

---

## Disable

```bash
openclaw config set plugins.entries.tokenjuice.enabled false
# or
openclaw plugins disable tokenjuice
```

---

## One Mind Use Case

On long infrastructure sessions (Sovereignty managing services, Guardian running audits), `git`, `systemctl`, and log commands generate verbose output that accumulates fast. Tokenjuice keeps the context window clean across multi-hour agent sessions.
