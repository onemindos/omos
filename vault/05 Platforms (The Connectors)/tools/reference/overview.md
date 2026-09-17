# MCP Servers

**MCP** = Model Context Protocol — Anthropic's open standard for connecting AI agents to external tools, data sources, and services via a unified interface.

## What MCP Is

Instead of every agent having custom API integrations, MCP provides a standard protocol where:
- **MCP Server** = a service that exposes tools (e.g., Firecrawl exposes `scrape`, `crawl`, `search`)
- **MCP Client** = the AI agent/host that calls those tools (Cursor, Claude Desktop, OpenClaw)
- **Transport** = how they communicate (`stdio` for local, `SSE` or HTTP for remote)

## MCP Servers in This Stack

| Server | Type | Use Case |
|--------|------|----------|
| [circle-so](circle-so.md) | Community | Query/post in Circle communities |
| [apify](00-24%20SO%20(Sovereign%20Operations)/01%20Intelligence%20(AI,%20Memory%20&%20Tools)%20(The%20brain)/tools/mcp/apify.md) | Web scraping | Run any Apify actor from an agent |
| [firecrawl](firecrawl.md) | Web scraping | Scrape, crawl, and search the web |
| [context7](context7.md) | Docs | Fetch current library docs for coding |

## MCP vs Direct API Integration

| | MCP | Direct API |
|-|-----|-----------|
| Setup | One config block | Custom code per service |
| Discoverability | Tools listed automatically | Manual |
| Standardization | Yes — any MCP client works | No |
| Flexibility | Limited to exposed tools | Full API access |

## Config Pattern (openclaw / claude-desktop)

```json
{
  "mcpServers": {
    "firecrawl": {
      "command": "npx",
      "args": ["-y", "firecrawl-mcp"],
      "env": { "FIRECRAWL_API_KEY": "fc-..." }
    }
  }
}
```

## Where Configs Live

- **Claude Desktop:** `~/Library/Application Support/Claude/claude_desktop_config.json`
- **OpenClaw:** agent `mcp.json` config block
- **Cursor:** `.cursor/mcp.json` in project root or `~/.cursor/mcp.json` globally


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
