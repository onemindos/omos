---
name: browser
description: Playwright-based browser automation — navigate pages, fill forms, extract data, take screenshots, monitor pages for changes, execute JavaScript.
metadata: { "status": "pending-port", "reference": "openclaw/tools/browser.py" }
---

# Browser Skill

Gives Legacy the ability to control a real web browser — navigate, click, fill forms, extract structured data, screenshot, and monitor pages.

## Capabilities (from browser.py blueprint)

- `navigate(url)` — open any URL
- `fill_form(url, fields)` — fill and submit web forms
- `screenshot(url)` — capture a page as PNG
- `extract_data(url, selectors)` — scrape structured data from pages
- `execute_js(url, script)` — run JavaScript on a page
- `monitor_page(url, check_interval)` — watch a page for changes

## To Build

Port `openclaw/tools/browser.py` → `index.js` using Playwright for Node.js:

```bash
npm install playwright
```

Key functions to port: `navigate`, `fill_form`, `screenshot`, `extract_data`, `execute_js`, `monitor_page`

## Use Cases for Legacy

- "Check if that website is back up"
- "Screenshot Zeus's ClickUp board"
- "Fill out this form for me"
- "Monitor this page and alert me when the price drops"
