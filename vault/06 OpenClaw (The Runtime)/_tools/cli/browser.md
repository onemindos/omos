---
name: OpenClaw CLI — Browser
updated: 2026-05-11
---

# `openclaw browser`

Manages browser control — lifecycle, profiles, tabs, automation, and debugging.

**Official docs:** [docs.openclaw.ai/cli/browser](https://docs.openclaw.ai/cli/browser.md)

---

## Commands

### Lifecycle & Profiles
```bash
openclaw browser start                          # Start managed Chrome
openclaw browser start --profile user           # Named profile
openclaw browser stop
```

Profiles: `openclaw` (managed), `user` (existing Chrome session), custom CDP endpoints.

### Tabs & Navigation
```bash
openclaw browser tabs                           # List open tabs (returns suggestedTargetId first)
openclaw browser open --url https://example.com
openclaw browser navigate --tab <id> --url https://example.com
openclaw browser close --tab <id>
```

### Page Interaction
```bash
openclaw browser click --tab <id> --ref <element-ref>
openclaw browser type --tab <id> --ref <element-ref> --text "Hello"
openclaw browser hover --tab <id> --ref <element-ref>
openclaw browser drag --tab <id> --from <ref> --to <ref>
```

### Content Inspection
```bash
openclaw browser snapshot --tab <id>            # Page snapshot
openclaw browser snapshot --tab <id> --links    # Include link destinations
openclaw browser screenshot --tab <id>          # Full page screenshot
openclaw browser screenshot --tab <id> --ref <el>  # Element screenshot
```

### State Control
```bash
openclaw browser emulate --device "iPhone 15"
openclaw browser cookies list --tab <id>
openclaw browser cookies clear --tab <id>
openclaw browser storage clear --tab <id>
```

### Debugging
```bash
openclaw browser logs --tab <id>                # Console output
openclaw browser network --tab <id>             # Network requests
openclaw browser errors --tab <id>              # Page errors
```

---

## Limitations (existing Chrome sessions via MCP)

- No CSS selector support
- No multiple file uploads
- Some timeout controls unavailable
- No batch actions

Use managed profiles or raw CDP for full functionality.
