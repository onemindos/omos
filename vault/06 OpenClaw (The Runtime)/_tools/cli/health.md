---
name: OpenClaw CLI — Health
updated: 2026-05-10
---

# `openclaw health`

Retrieves the current health status of a running Gateway.

**Official docs:** [docs.openclaw.ai/cli/health](https://docs.openclaw.ai/cli/health.md)

---

## Commands

```bash
openclaw health                    # Health snapshot (may use cache)
openclaw health --verbose          # Live probe + full details per agent/account
openclaw health --json             # Machine-readable output
openclaw health --timeout 2500     # Custom connection timeout (ms, default 10000)
openclaw health --debug            # Debug-level output
```

---

## Notes

- Default may return a cached snapshot — use `--verbose` to force a live probe
- `--verbose` shows connection details and expands per-agent session stores
- Useful as a quick check before running commands against the gateway

```bash
# Quick sanity check on omos-core
ssh ubuntu@100.123.126.34 'openclaw health --json'
```
