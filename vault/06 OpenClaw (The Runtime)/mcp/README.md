---
title: "OpenClaw — MCP Servers"
type: reference
status: active
created: 2026-05-11
modified: 2026-05-11
updated: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO-02, openclaw/mcp]
---

# OpenClaw — MCP Servers

MCP (Model Context Protocol) servers extend agent capabilities with external integrations. Config: `~/.openclaw/openclaw.json` → `mcp.servers{}`.

---

## Active MCP Servers

### circle-mcp

| Field | Value |
|-------|-------|
| **Status** | ✅ Active |
| **Purpose** | Circle community platform — manage spaces, posts, courses, members, events |
| **Used by** | Legacy (main), Creed, Eternal |
| **Config key** | `mcp.servers.circle-mcp` |

Provides 100+ tools for interacting with the Circle community platform. Agents with this MCP loaded can create posts, manage members, run events, and more.

To check Circle MCP config:
```bash
ssh ubuntu@100.123.126.34 "openclaw mcp list"
```

---

## Managing MCP Servers

```bash
# List all MCP servers
ssh ubuntu@100.123.126.34 "openclaw mcp list"

# Check MCP server status
ssh ubuntu@100.123.126.34 "openclaw mcp status circle-mcp"

# Add a new MCP server (edit openclaw.json directly)
# Under mcp.servers, add:
# "my-server": { "command": "...", "args": [...], "env": {...} }

# Reload gateway after config changes
ssh ubuntu@100.123.126.34 "openclaw gateway reload"
```

Config location: `~/.openclaw/openclaw.json` → `mcp.servers{}`

---

## Adding a New MCP Server

1. SSH into omos-core
2. Edit `~/.openclaw/openclaw.json`
3. Add entry under `mcp.servers`:
```json
"server-name": {
  "command": "node",
  "args": ["/path/to/server.js"],
  "env": {
    "API_KEY": "..."
  }
}
```

4. Reload: `openclaw gateway reload`
5. Verify: `openclaw mcp list`
6. Document here

---

## Official Docs

- [MCP overview](https://docs.openclaw.ai/concepts/mcp.md)
- [Adding MCP servers](https://docs.openclaw.ai/gateway/mcp.md)

---

## Related

- [[../agents/README|Agents]] — Which agents have MCP access
- [[../plugin-skills (Custom Skills)/PLUGIN-SKILLS|Plugin Skills]] — Skill-based integrations
- [[../gateway/OPENCLAW-CONFIG-REFERENCE|Config Reference]] — Full openclaw.json reference
