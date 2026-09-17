---
type: index
domain: SO-01
status: active
updated: 2026-05-14
owner: zeus
tags: [runtime/claude-code, domain/so, type/index]
ai_summary: "Claude Code MCP plugins and tool extensions — servers, integrations, and tool configs wired into Claude Code sessions."
---

# Claude Code — Plugins (MCP)

MCP (Model Context Protocol) servers and tool extensions available to Claude Code. Configured in `~/.claude/settings.json` under `mcpServers`.

## What Goes Here

- MCP server configs and connection details
- Tool extension documentation
- Per-project MCP overrides
- Installation and setup notes for each plugin

## Active MCP Servers

| Server | Command | What It Provides |
| ------ | ------- | ---------------- |
| `context7` | `npx -y @upstash/context7-mcp@latest` | Up-to-date library docs |
| `github` | `npx -y @modelcontextprotocol/server-github` | GitHub repos, issues, PRs |
| `filesystem` | `npx -y @modelcontextprotocol/server-filesystem <path>` | Read/write local files |
| `desktop-commander` | `npx -y @wonderwhy-er/desktop-commander@latest` | Terminal, process management |
| *(add as configured)* | — | — |

## Config Location

```
~/.claude/settings.json   ← global MCP config
```

## Adding a New Server

1. Add entry to `mcpServers` in `~/.claude/settings.json`
2. Document it in the table above
3. Test with `/mcp` in a Claude Code session

---
*MCP server list also tracked in `06 OpenClaw/mcp/`. Cross-reference when adding servers used by both runtimes.*
