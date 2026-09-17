# Apify

> **Web scraping and browser automation platform — extract data from any website**

## Overview

Apify is a cloud platform for web scraping, browser automation, and data extraction. Provides ready-made scrapers ("Actors") for popular sites plus infrastructure to build custom crawlers. Used for competitive intelligence, lead generation, content extraction, and research automation in the GE and SO domains.

## Status

| Attribute   | Value                               |
| ----------- | ----------------------------------- |
| Status      | 🔄 Evaluating                       |
| Category    | Automation / Data                   |
| Type        | Web Scraping Platform               |
| Integration | REST API, JavaScript SDK, CLI       |

## Use Cases in OneMind

1. **Research Automation** — Scrape web content for the CODEX Capture pipeline
2. **Competitive Intelligence** — Monitor competitor sites, pricing, content for GE domain
3. **Lead Data** — Extract contact and company data for GE/outreach workflows
4. **Content Extraction** — Pull articles, listings, or structured data into the vault
5. **N8N Integration** — Trigger Apify Actors from N8N workflows via HTTP node

## Key Features

- 2000+ pre-built Actors (scrapers for Google, LinkedIn, Amazon, etc.)
- Custom Actor development (Node.js / Python)
- Proxy rotation and CAPTCHA solving built-in
- Scheduled runs and webhook triggers
- Dataset storage and export (JSON, CSV, Excel)
- Puppeteer and Playwright support

## Quick Start

```bash
# Install CLI
npm install -g apify-cli

# Create a new Actor
apify create my-scraper
cd my-scraper
apify run      # Run locally

# Push and deploy to cloud
apify push
```

```javascript
// Call an Actor via API
const response = await fetch(
  'https://api.apify.com/v2/acts/apify~web-scraper/runs',
  {
    method: 'POST',
    headers: { Authorization: 'Bearer APIFY_TOKEN' },
    body: JSON.stringify({ startUrls: [{ url: 'https://example.com' }] })
  }
)
```

## Links

- Website: [apify.com](https://apify.com)
- Actor Store: [apify.com/store](https://apify.com/store)
- API Docs: [docs.apify.com](https://docs.apify.com)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
