---
name: OpenClaw CLI — ACP
updated: 2026-05-11
---

# `openclaw acp`

Gateway-backed Agent Client Protocol bridge — enables IDEs and ACP clients to interact with OpenClaw via stdio.

**Official docs:** [docs.openclaw.ai/cli/acp](https://docs.openclaw.ai/cli/acp.md)

---

## Commands

```bash
# Start ACP bridge (stdio → Gateway WebSocket)
openclaw acp

# Remote gateway
openclaw acp --url wss://gateway-host:18789 --token-file ~/.openclaw/gateway.token
```

---

## What the Bridge Exposes

Fully implemented:
- `initialize` — session setup
- `newSession` — create session
- `prompt` — send a prompt
- `cancel` — cancel execution
- Session listing with workspace filtering
- Session resumption and closure

Partially supported:
- Session history replay (bridge-created sessions only)
- Prompt content with embedded resources and images
- Tool streaming (best-effort file locations)
- Execution approval relays

Not supported:
- Per-session MCP server configuration
- Client filesystem methods
- Client terminal operations
- Thought streaming / ACP plans

---

## Primary Use Cases

1. **IDE Integration** — Configure Zed or other IDEs to run `openclaw acp` over stdio
2. **Claude Code access** — Enables Claude Code to communicate with OpenClaw sessions
3. **Session targeting** — Route to specific agents: `agent:main:main`

---

## Security

Use `--token-file` or env vars rather than inline `--token`. If multiple ACP clients share the same Gateway session key, event and cancel routing are best-effort.

---

## Note

This is distinct from ACP Agents (`acpx`) — the bridge handles *external clients connecting to OpenClaw*, while ACP Agents have OpenClaw *launching external harnesses* like Claude Code.
