---
name: OpenClaw CLI — MCP
updated: 2026-05-11
---

# `openclaw mcp`

Operates OpenClaw as an MCP server, or manages MCP server definitions for runtimes to consume.

**Official docs:** [docs.openclaw.ai/cli/mcp](https://docs.openclaw.ai/cli/mcp.md)

---

## Two Modes

| Mode | What it does |
|---|---|
| `mcp serve` | Exposes OpenClaw as an MCP server (Claude Code can connect) |
| `mcp list/show/set/unset` | Manages MCP server definitions in OpenClaw config |

---

## Commands

```bash
# Run OpenClaw as MCP server (stdio bridge → Gateway)
openclaw mcp serve

# Manage MCP server definitions
openclaw mcp list
openclaw mcp show <name>
openclaw mcp set <name> --type stdio --command "npx -y @some/mcp-server"
openclaw mcp unset <name>
```

---

## When Serving: Tools Exposed

| Tool | Description |
|---|---|
| `conversations_list` | Discover routed sessions with channel metadata |
| `messages_read` | Retrieve transcript history |
| `events_poll` | Access live queued events |
| `events_wait` | Wait for queued events |
| `messages_send` | Reply through existing routes |
| `permissions_list_open` | List open approval requests |
| `permissions_respond` | Respond to approval requests |

---

## MCP Server Config (in openclaw.json)

Three transport types under `mcp.servers`:

```json
{
  "mcp": {
    "servers": {
      "my-server": {
        "type": "stdio",
        "command": "npx",
        "args": ["-y", "@some/mcp-server"]
      }
    }
  }
}
```

Transport types: `stdio`, `sse` (HTTP Server-Sent Events), `streamable-http`

---

## Notes

- Live event queue starts when the bridge connects — older history via `messages_read`
- Disconnecting disposes the event queue entirely
- Security relies on existing Gateway authentication — the bridge does not invent routing
