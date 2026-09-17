# Browser Tool - Web Browsing & Data Extraction

> Replaces the browser extension with AI-driven on-demand extraction

## Overview

The Browser Tool gives Legacy the ability to:
- **Extract** structured data from any webpage (receipts, orders, bookings, articles)
- **Screenshot** pages for visual analysis
- **Interact** with pages (click, type, scroll, fill forms)
- **Search** the web and extract results

This replaces the passive browser extension with an active, AI-driven approach where Legacy fetches data on-demand rather than capturing everything you browse.

## Installation

### 1. Install Playwright

```bash
# Install Python package
pip install playwright>=1.40.0

# Install browser binaries (required!)
playwright install chromium
```

### 2. Verify Installation

```bash
python -c "from playwright.sync_api import sync_playwright; print('Playwright ready!')"
```

## Usage

### Extract Data from a URL

```python
# Auto-detect page type and extract
browser(action="extract", url="https://amazon.com/order/123-456")

# Force specific extraction type
browser(action="extract", url="https://example.com/receipt", extract_type="receipt")

# Wait for specific element before extracting
browser(action="extract", url="https://spa-app.com/dashboard", wait_for=".data-loaded")

# Include screenshot
browser(action="extract", url="https://example.com", screenshot=True)
```

### Extraction Types

| Type | Auto-Detected When | Extracts |
|------|-------------------|----------|
| `receipt` | URL contains: order, invoice, receipt, confirmation | order_number, total, date, merchant |
| `booking` | URL contains: booking, reservation, flight, hotel | confirmation, date, location, provider |
| `article` | URL contains: article, blog, post, news | title, author, date, content_preview, word_count |
| `page` | Default fallback | text_preview, links_count, images_count |

### Take Screenshots

```python
# Quick screenshot
browser(action="screenshot", url="https://example.com")

# Screenshot is included in response as base64
# result["screenshot_base64"] contains PNG data
```

### Interact with Pages

```python
# Login to a website
browser(
    action="interact",
    url="https://example.com/login",
    actions=[
        {"type": "type", "selector": "#email", "value": "user@example.com"},
        {"type": "type", "selector": "#password", "value": "secret123"},
        {"type": "click", "selector": "button[type=submit]"},
        {"type": "wait", "selector": ".dashboard"}
    ]
)

# Fill out a form
browser(
    action="interact",
    url="https://example.com/form",
    actions=[
        {"type": "type", "selector": "#name", "value": "John Doe"},
        {"type": "select", "selector": "#country", "value": "US"},
        {"type": "click", "selector": ".submit-btn"},
        {"type": "screenshot"}  # Capture result
    ]
)
```

### Action Types

| Action | Parameters | Description |
|--------|-----------|-------------|
| `click` | `selector` | Click an element |
| `type` | `selector`, `value` | Type text into input |
| `select` | `selector`, `value` | Select dropdown option |
| `scroll` | `value` (pixels) | Scroll down the page |
| `wait` | `selector` | Wait for element to appear |
| `screenshot` | (none) | Take screenshot at this point |

### Search the Web

```python
# Google search
browser(action="search", query="python playwright tutorial", num_results=5)

# Use DuckDuckGo
browser(action="search", query="privacy tools", engine="duckduckgo")

# Use Bing
browser(action="search", query="microsoft docs", engine="bing")
```

## Response Format

All responses are JSON strings with the following structure:

### Extract Response
```json
{
  "action": "extract",
  "type": "receipt",
  "url": "https://amazon.com/order/123",
  "title": "Amazon.com - Order Details",
  "data": {
    "order_number": "123-456-789",
    "total": "$49.99",
    "date": "January 20, 2026",
    "merchant": "amazon"
  },
  "metadata": {
    "og_title": "Order Confirmation",
    "og_description": "Your order has been placed"
  },
  "structured_data": [
    {"@type": "Order", "orderNumber": "123-456-789"}
  ]
}
```

### Interact Response
```json
{
  "action": "interact",
  "url": "https://example.com/dashboard",
  "actions_completed": [
    {"action": "type", "selector": "#email", "success": true},
    {"action": "click", "selector": "button", "success": true}
  ],
  "final_page": {
    "type": "page",
    "title": "Dashboard",
    "data": {...}
  }
}
```

### Search Response
```json
{
  "action": "search",
  "query": "python playwright",
  "engine": "google",
  "results_count": 5,
  "results": [
    {
      "title": "Playwright Python Documentation",
      "url": "https://playwright.dev/python/",
      "snippet": "Playwright enables reliable end-to-end testing..."
    }
  ]
}
```

## Example Conversations with Legacy

### Get Order Details
```
You: "Legacy, check my Amazon order from last week"
Legacy: [Uses browser tool to extract order data]
Legacy: "Your order #123-456-789 for $49.99 was placed on Jan 20.
         It includes 2 items and is expected to arrive by Jan 25."
```

### Research a Topic
```
You: "Find me the best Python testing frameworks in 2026"
Legacy: [Uses browser(action="search")]
Legacy: "Based on my search, the top frameworks are:
         1. Pytest - Most popular, extensive plugin ecosystem
         2. Playwright - Best for browser testing
         3. Robot Framework - Great for acceptance testing"
```

### Extract Article for Later
```
You: "Save that Hacker News article about AI agents"
Legacy: [Uses browser(action="extract", extract_type="article")]
Legacy: "I've extracted 'Building Autonomous AI Agents' by John Smith.
         2,450 words. Key topics: agent loops, tool use, memory systems.
         Want me to summarize it or save it to your knowledge base?"
```

## Alternative: agent-browser (NPM)

For **complex multi-step browser automation** that requires:
- Multiple tabs
- Persistent sessions
- Complex navigation flows
- Visual AI understanding

Use **agent-browser** instead:

```bash
# Install
npm install -g agent-browser
npx playwright install

# Use with natural language
agent-browser "Go to GitHub, search for 'agno framework', and star the first repo"
agent-browser "Log into my bank, check my balance, and take a screenshot"
```

### When to Use Which

| Scenario | Use |
|----------|-----|
| Extract data from a URL | `browser` tool |
| Take a screenshot | `browser` tool |
| Fill out a simple form | `browser` tool |
| Search the web | `browser` tool |
| Multi-step workflow across pages | `agent-browser` |
| Visual AI tasks ("find the red button") | `agent-browser` |
| Persistent login sessions | `agent-browser` |
| Complex navigation | `agent-browser` |

## Architecture

```
┌─────────────────────────────────────────────────┐
│                  Legacy Agent                    │
├─────────────────────────────────────────────────┤
│                                                  │
│   browser(action="extract", url="...")          │
│            │                                     │
│            ▼                                     │
│   ┌─────────────────────┐                       │
│   │   Browser Tool      │                       │
│   │   (browser.py)      │                       │
│   └─────────┬───────────┘                       │
│             │                                    │
│             ▼                                    │
│   ┌─────────────────────┐                       │
│   │     Playwright      │                       │
│   │  (headless Chrome)  │                       │
│   └─────────┬───────────┘                       │
│             │                                    │
│             ▼                                    │
│   ┌─────────────────────┐                       │
│   │    Target Website   │                       │
│   └─────────────────────┘                       │
│                                                  │
└─────────────────────────────────────────────────┘
```

## Comparison: Extension vs Browser Tool

| Aspect | Browser Extension | Browser Tool |
|--------|------------------|--------------|
| **Trigger** | Passive (auto-captures) | Active (on-demand) |
| **Auth** | Uses your browser sessions | Headless (no auth) |
| **Control** | Captures everything | You decide what to capture |
| **AI Integration** | Sends to backend | Native Legacy tool |
| **Maintenance** | Extension updates | Python package |
| **Privacy** | Sees all browsing | Only requested URLs |

## Troubleshooting

### "Playwright not installed"
```bash
pip install playwright
playwright install chromium
```

### "Browser launch failed"
```bash
# Install system dependencies (Linux)
playwright install-deps chromium

# Or use Docker
docker run -it mcr.microsoft.com/playwright/python:v1.40.0
```

### "Timeout waiting for selector"
- Increase wait time in actions
- Check if the selector is correct
- The page might require JavaScript to load

### "Access denied" or "CAPTCHA"
- Some sites block headless browsers
- Consider using `agent-browser` with visual AI
- Or authenticate manually and use session cookies

## Files

| File | Purpose |
|------|---------|
| [agents/legacy/tools/browser.py](../../agents/legacy/tools/browser.py) | Browser tool implementation |
| [agents/legacy/tools/__init__.py](../../agents/legacy/tools/__init__.py) | Tool exports |
| [agents/legacy/agent.py](../../agents/legacy/agent.py) | Agent with browser tool |
| [agents/legacy/requirements.txt](../../agents/legacy/requirements.txt) | Playwright dependency |

## Related

- [Extension API](../legacy/EXTENSION-API.md) - Legacy extension endpoints (archived)
- [Full Access Architecture](../architecture/full-access-architecture.md) - Data access strategy
- [Composio Tools](../../agents/legacy/tools/composio_tools.py) - OAuth-managed integrations
