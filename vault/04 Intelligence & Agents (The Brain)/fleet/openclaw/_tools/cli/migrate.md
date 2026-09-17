---
name: OpenClaw CLI — Migrate
updated: 2026-05-11
---

# `openclaw migrate`

Imports state from other agent systems (Claude, Codex, Hermes).

**Official docs:** [docs.openclaw.ai/cli/migrate](https://docs.openclaw.ai/cli/migrate.md)

---

## Commands

```bash
# Preview migration (dry run — no changes)
openclaw migrate claude --dry-run
openclaw migrate codex --dry-run
openclaw migrate hermes --dry-run

# Apply migration
openclaw migrate apply claude --yes
openclaw migrate apply codex --yes
openclaw migrate apply hermes --yes

# Selective skill migration
openclaw migrate apply codex --yes --skill skill-name
```

---

## Available Providers

| Provider | Source | Imports |
|---|---|---|
| `claude` | `~/.claude` | Projects, skills, MCP server definitions |
| `codex` | `~/.codex` or `$CODEX_HOME` | Skills, curated plugins |
| `hermes` | `~/.hermes` | Config, models, skills |

---

## Safety Features

- **Preview-first** — always shows itemized plan before applying
- Automatic backups before state changes
- Conflict detection — refuses to proceed until resolved
- Secrets excluded by default; use `--include-secrets` to import credentials

Items not imported automatically (included in migration report for manual review): hooks, permissions, cache state.
