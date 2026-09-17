---
title: "Claude Desktop"
type: reference
status: active
tags: [platform, ai, tools, assistant, anthropic]
---
# Claude Desktop

> Anthropic's native desktop app. Chat with Claude, attach files, use MCP servers, extended thinking.

## Key Features
| Feature | Notes |
|---------|-------|
| MCP support | Connect local tools (filesystem, GitHub, etc.) |
| File upload | PDFs, images, code files |
| Projects | Organized conversations with shared context |
| Extended thinking | Deep reasoning mode |
| Artifacts | Code, documents, diagrams in side panel |
| Search | Web search integration |
| Computer use | Control your desktop (beta) |

## MCP Config
```json
// ~/Library/Application Support/Claude/claude_desktop_config.json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path"]
    }
  }
}
```

## Cost
Included with Anthropic Pro ($20/mo) or Team ($30/mo)
