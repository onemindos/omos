---
title: "Firecrawl"
type: reference
domain: SO
status: active
created: 2026-05-28
updated: 2026-05-28
tags: [platform, ai, web-scraping, mcp]
---

# Firecrawl

> AI-powered web scraping and crawling API. Turns any website into clean markdown or structured data.

## Use Case in OneMind OS

- Web content extraction for research agents
- Clean markdown output for Obsidian/vault ingestion
- Structured data extraction (tables, lists, metadata)
- Available as MCP server for Hermes

## API

```bash
# Base URL
https://api.firecrawl.dev/v1

# Scrape a page
curl -X POST https://api.firecrawl.dev/v1/scrape \
  -H "Authorization: Bearer $FIRECRAWL_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com", "formats": ["markdown"]}'

# Crawl a site
curl -X POST https://api.firecrawl.dev/v1/crawl \
  -H "Authorization: Bearer $FIRECRAWL_API_KEY" \
  -d '{"url": "https://example.com", "limit": 50}'
```

## MCP Integration

Configured in Hermes `config.yaml` as an MCP server. Provides:
- `firecrawl_scrape` — single page to markdown
- `firecrawl_crawl` — multi-page crawl
- `firecrawl_map` — sitemap discovery

## Pricing

Free tier: 500 credits/month. Pro: $19/mo for 3,000 credits.
