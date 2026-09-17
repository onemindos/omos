# Apify MCP

**Type:** Web scraping & automation platform
**MCP Package:** `apify-mcp-server` (official)
**Website:** apify.com
**Auth:** Apify API token (`APIFY_TOKEN`)

## What It Does

Apify hosts thousands of pre-built "Actors" — scrapers and automation bots for any website. The MCP server exposes a `run_actor` tool that lets agents trigger any Actor by ID and get results back, without building custom scrapers.

## Key Tools Exposed

| Tool | Description |
|------|-------------|
| `run_actor` | Run any Apify Actor with input params |
| `get_actor_info` | Get Actor details and input schema |
| `search_actors` | Find Actors by keyword |
| `get_dataset_items` | Fetch results from a previous run |
| `get_run_status` | Check if an Actor run finished |

## Most Useful Actors for Stack

| Actor ID | What It Does |
|----------|-------------|
| `apify/instagram-scraper` | Scrape IG profiles, hashtags, posts |
| `apify/twitter-scraper` | Scrape X/Twitter profiles and searches |
| `apify/linkedin-scraper` | LinkedIn profile/company data |
| `apify/google-maps-scraper` | Local business data |
| `apify/website-content-crawler` | Crawl entire sites for AI training |
| `apify/yt-channel-scraper` | YouTube channel/video data |

## Config

```json
{
  "mcpServers": {
    "apify": {
      "command": "npx",
      "args": ["-y", "apify-mcp-server"],
      "env": {
        "APIFY_TOKEN": "apify_api_..."
      }
    }
  }
}
```

## Stack Use Cases

- **Oracle agent:** Deep research on people/companies before Zeus's meetings
- **Empire agent:** Competitor analysis, market research, social monitoring
- **Legacy agent:** "Scrape this site" or "get all posts from this IG account"
- **Forge agent:** Gather dataset for AI training pipelines

## Pricing Note

Actors run on Apify's compute (not local). Costs are per-Actor-run, billed to Apify account. Free tier: $5/month credits.
