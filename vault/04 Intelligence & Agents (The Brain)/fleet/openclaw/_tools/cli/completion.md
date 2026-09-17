---
name: OpenClaw CLI — Completion
updated: 2026-05-11
---

# `openclaw completion`

Generate and install shell completion scripts.

**Official docs:** [docs.openclaw.ai/cli/completion](https://docs.openclaw.ai/cli/completion.md)

```bash
openclaw completion                          # Print zsh completion to stdout (default)
openclaw completion --shell zsh              # Explicit shell
openclaw completion --shell bash
openclaw completion --shell fish
openclaw completion --shell powershell
openclaw completion --install               # Install by adding source line to shell profile
openclaw completion --shell fish --install
openclaw completion --write-state           # Write to $OPENCLAW_STATE_DIR/completions
openclaw completion -y                      # Skip install confirmation
```

---

## Notes

- `--install` writes an "OpenClaw Completion" block into your shell profile pointing at the cached script
- Without `--install` or `--write-state`, prints the script to stdout
- Completion generation eagerly loads command trees so all nested subcommands are included
