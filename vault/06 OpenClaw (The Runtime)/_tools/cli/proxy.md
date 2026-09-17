---
name: OpenClaw CLI — Proxy
updated: 2026-05-11
---

# `openclaw proxy`

Validates operator-managed proxy routing and debugs transport-level traffic.

**Official docs:** [docs.openclaw.ai/cli/proxy](https://docs.openclaw.ai/cli/proxy.md)

---

## Commands

```bash
# Validation
openclaw proxy validate                         # Check proxy config + connectivity
openclaw proxy validate --apns-reachable        # Also verify APNs reachability
openclaw proxy validate --allowed-url <url>     # Test a URL that should succeed
openclaw proxy validate --denied-url <url>      # Test a URL that should be blocked

# Debugging
openclaw proxy start                            # Launch local proxy (default: 127.0.0.1)
openclaw proxy start --host 0.0.0.0 --port 8888
openclaw proxy run <command>                    # Run command with traffic capture
openclaw proxy sessions                         # List captures
openclaw proxy query <preset>                   # Analyze traffic patterns
openclaw proxy blob <id>                        # Read captured data
openclaw proxy purge                            # Clear all local captures
```

---

## Query Presets

`double-sends` · `retry-storms` · `cache-busting` · (others)

---

## Notes

- Captured data is local debugging only — purge after you're done
- Direct upstream forwarding for diagnostics is disabled by default in managed proxy mode
- Validates that public destinations succeed through the proxy, and that the proxy blocks loopback canaries
