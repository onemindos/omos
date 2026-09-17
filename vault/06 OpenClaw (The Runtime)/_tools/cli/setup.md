---
name: OpenClaw CLI — Setup
updated: 2026-05-11
---

# `openclaw setup`

Initialize the baseline config and agent workspace without the full guided onboarding flow. Use this for scripted or minimal installs.

**Official docs:** [docs.openclaw.ai/cli/setup](https://docs.openclaw.ai/cli/setup.md)

---

## Commands

```bash
openclaw setup                                    # Baseline config + workspace only
openclaw setup --workspace ~/.openclaw/workspace  # Custom workspace directory
openclaw setup --wizard                           # Run full onboarding
openclaw setup --non-interactive --mode remote --remote-url wss://gateway-host:18789 --remote-token <token>
```

---

## Options

| Option | Description |
|---|---|
| `--workspace <dir>` | Agent workspace directory |
| `--wizard` | Run full onboarding |
| `--non-interactive` | Run without prompts |
| `--mode <local\|remote>` | Onboarding mode |
| `--import-from <provider>` | Migration provider |
| `--import-source <path>` | Source agent home for migration |
| `--remote-url <url>` | Remote Gateway WebSocket URL |
| `--remote-token <token>` | Remote Gateway token |

---

## Notes

- Plain `openclaw setup` initializes config + workspace without the full onboarding flow
- For full guided setup, use `openclaw onboard`
- After plain setup: `openclaw configure` for targeted changes, `openclaw channels add` for channels
- Nix mode (`OPENCLAW_NIX_MODE=1`): setup writes are refused — config is managed by Nix
