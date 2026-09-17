---
name: OpenClaw CLI — Reset
updated: 2026-05-11
---

# `openclaw reset`

Reset local config/state (keeps the CLI installed).

**Official docs:** [docs.openclaw.ai/cli/reset](https://docs.openclaw.ai/cli/reset.md)

---

## Commands

```bash
openclaw backup create                                      # Always backup first
openclaw reset                                              # Interactive reset
openclaw reset --dry-run                                    # Preview what would be removed
openclaw reset --scope config --yes --non-interactive
openclaw reset --scope config+creds+sessions --yes --non-interactive
openclaw reset --scope full --yes --non-interactive
```

---

## Options

| Option | Description |
|---|---|
| `--scope <scope>` | `config`, `config+creds+sessions`, or `full` |
| `--yes` | Skip confirmation prompts |
| `--non-interactive` | Disable prompts (requires `--scope` and `--yes`) |
| `--dry-run` | Print actions without removing files |

---

## Notes

- Always run `openclaw backup create` first
- Without `--scope`, reset uses an interactive prompt to choose what to remove
- `--non-interactive` is only valid when both `--scope` and `--yes` are set
