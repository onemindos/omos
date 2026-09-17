---
name: OpenClaw CLI — Clawbot (Legacy Namespace)
updated: 2026-05-11
---

# `openclaw clawbot`

Legacy namespace. Previously the primary CLI entry point — now an alias layer maintained for backward compatibility only.

**Official docs:** [docs.openclaw.ai/cli/clawbot](https://docs.openclaw.ai/cli/clawbot.md)

---

## Still-Active Alias

```bash
openclaw clawbot qr        # Identical to: openclaw qr
```

All other `clawbot` sub-commands were migrated to top-level `openclaw` commands in v2025. The only surviving alias is `clawbot qr`.

---

## Migration Map

| Old command | New command |
|---|---|
| `openclaw clawbot qr` | `openclaw qr` |
| `openclaw clawbot start` | `openclaw gateway start` |
| `openclaw clawbot stop` | `openclaw gateway stop` |
| `openclaw clawbot status` | `openclaw status` |
| `openclaw clawbot agent` | `openclaw agent` |
| `openclaw clawbot logs` | `openclaw logs` |

---

## Notes

- Do not use `clawbot` in new scripts or automation — it will be removed in a future major version
- The `qr` alias is the only exception kept for mobile onboarding convenience
