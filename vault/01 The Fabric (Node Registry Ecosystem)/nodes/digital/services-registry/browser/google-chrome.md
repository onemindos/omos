---
title: "Google Chrome"
type: reference
status: active
tags: [platform, browser, google]
---
# Google Chrome

> Primary browser. DevTools, extensions, profile sync, PWAs.

## Key Features
| Feature | Notes |
|---------|-------|
| DevTools | F12 — inspect, debug, network, performance |
| Extensions | Largest extension ecosystem |
| Profiles | Multiple profiles (work, personal) |
| PWA | Install web apps as native |
| Sync | Bookmarks, passwords, tabs across devices |
| V8 | Fastest JS engine |

## CLI
```bash
# Open URL
open -a "Google Chrome" https://example.com

# DevTools remote debugging
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --remote-debugging-port=9222

# Headless
chrome --headless --dump-dom https://example.com
```

## Integration
- DevTools Protocol for browser automation
- Hermes CUA driver (computer use)
- Extensions: 1Password, Raindrop, uBlock Origin
