---
name: OpenClaw CLI — Uninstall
updated: 2026-05-11
---

# `openclaw uninstall`

Uninstall the gateway service + local data (CLI binary remains).

**Official docs:** [docs.openclaw.ai/cli/uninstall](https://docs.openclaw.ai/cli/uninstall.md)

---

## Commands

```bash
openclaw backup create                                      # Always backup first
openclaw uninstall                                          # Interactive uninstall
openclaw uninstall --dry-run                                # Preview what would be removed
openclaw uninstall --service --yes --non-interactive        # Remove service only
openclaw uninstall --state --workspace --yes --non-interactive
openclaw uninstall --all --yes                              # Remove everything
```

---

## Options

| Option | Description |
|---|---|
| `--service` | Remove the gateway service |
| `--state` | Remove state and config |
| `--workspace` | Remove workspace directories |
| `--app` | Remove the macOS app |
| `--all` | Remove service, state, workspace, and app |
| `--yes` | Skip confirmation prompts |
| `--non-interactive` | Disable prompts (requires `--yes`) |
| `--dry-run` | Print actions without removing files |

---

## Notes

- Always run `openclaw backup create` first
- `--all` is shorthand for `--service --state --workspace --app`
- `--non-interactive` requires `--yes`
- The CLI binary itself is NOT removed — only gateway service and local data
