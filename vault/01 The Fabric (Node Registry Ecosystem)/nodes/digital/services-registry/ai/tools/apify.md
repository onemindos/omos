---
title: "Apify"
type: reference
status: active
tags: [platform, ai, tools, web-scraping, automation, cloud]
---
# Apify

> Cloud web scraping and automation platform. Pre-built Actors for any website, headless browsers, proxy rotation, scheduling.

## Key Features
| Feature | Notes |
|---------|-------|
| Actors | Pre-built scrapers (1500+ in store) |
| Crawlee | Open-source crawling framework (JS/Python) |
| Proxy | Residential + datacenter rotation |
| Storage | Key-value, datasets, request queues |
| Scheduling | Cron-like actor runs |
| API | REST API for all operations |
| Webhooks | Notify on run complete |

## Use Cases
| Actor | Purpose |
|-------|---------|
| Web Scraper | Generic CSS/jQuery extraction |
| Cheerio Scraper | Fast HTML parsing (no browser) |
| Puppeteer Scraper | Full browser rendering (JS-heavy sites) |
| Google Search Scraper | SERP results |
| Instagram Scraper | Posts, profiles, hashtags |
| YouTube Scraper | Videos, channels, comments |
| Amazon Scraper | Products, reviews, prices |

## CLI
```bash
npm install -g apify-cli
apify login
apify create my-actor
apify run
apify push
apify call ACTOR_ID -i '{"url": "https://..."}'
```

## Python SDK
```python
from apify_client import ApifyClient

client = ApifyClient("YOUR_TOKEN")
run = client.actor("apify/web-scraper").call(
    run_input={"startUrls": [{"url": "https://example.com"}]}
)
dataset = client.dataset(run["defaultDatasetId"]).list_items()
```

## vs Firecrawl
| | Apify | Firecrawl |
|---|---|---|
| Approach | Full platform + actor marketplace | AI-optimized single-page extraction |
| Self-host | Yes (Crawlee open-source) | Yes |
| Pre-built scrapers | 1500+ actors | None (you define) |
| Best for | Complex multi-page crawls, scheduled scraping | LLM-ready page content, one-shot extractions |

## Integration
- Webhooks → NATS for scrape completion events
- Datasets → DuckDB/Parquet for analysis
- Crawlee (self-hosted) for no-cloud dependency
