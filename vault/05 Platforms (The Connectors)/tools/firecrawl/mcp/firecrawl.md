# Firecrawl MCP

**Type:** Web scraping & content extraction
**MCP Package:** `firecrawl-mcp` (official)
**Website:** firecrawl.dev
**Auth:** `FIRECRAWL_API_KEY`

## What It Does

Firecrawl converts any webpage into clean markdown — handling JavaScript rendering, auth-walled pages, and complex layouts. The MCP server gives agents a `scrape` and `crawl` tool so they can pull live web content without leaving the chat context.

## Key Tools Exposed

| Tool | Description |
|------|-------------|
| `scrape` | Scrape a single URL → clean markdown |
| `crawl` | Recursively crawl a site and return all pages |
| `search` | Web search with full page content returned |
| `extract` | Extract structured data from a page using a schema |
| `map` | Get all URLs from a site (sitemap-style) |
| `batch_scrape` | Scrape multiple URLs in one call |

## Config

```json
{
  "mcpServers": {
    "firecrawl": {
      "command": "npx",
      "args": ["-y", "firecrawl-mcp"],
      "env": {
        "FIRECRAWL_API_KEY": "fc-..."
      }
    }
  }
}
```

## Firecrawl vs Apify

| | Firecrawl | Apify |
|-|-----------|-------|
| Best for | Clean content extraction | Structured data scraping |
| Output | Markdown | JSON |
| Custom logic | No | Yes (actors) |
| Speed | Fast (single pages) | Slower (actor startup) |
| Use when | Reading articles/docs | Scraping social/structured data |

## Stack Use Cases

- **Legacy agent:** "What does this page say?" — paste URL, get clean text back
- **Oracle agent:** Research deep-dive — crawl a domain and synthesize insights
- **Forge agent:** Feed web content into training pipelines
- **Any agent:** Real-time web context in conversation (live pricing, news, docs)

## Pricing

Free tier: 500 pages/month. Paid from $16/month (3,000 pages). API key in 1Password.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
