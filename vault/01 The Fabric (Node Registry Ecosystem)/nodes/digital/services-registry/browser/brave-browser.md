---
title: "Brave Browser"
type: reference
domain: SO
status: active
created: 2026-05-28
updated: 2026-05-28
tags: [platform, browser, privacy, web3]
---

# Brave Browser

> Privacy-first Chromium browser with built-in ad blocking, Tor, and crypto wallet.

## Use Case in OneMind OS

- Primary browser for development (DevTools, extensions)
- Hermes `browser` toolset target (Chromium-compatible)
- Privacy browsing without tracking
- Leo AI assistant (built-in)

## Key Features

| Feature | Notes |
|---------|-------|
| Ad/tracker blocking | Shields — blocks by default |
| Brave Search | Privacy search engine |
| Leo AI | Built-in Claude-powered assistant |
| Sync | Cross-device via Brave Sync chain |
| Tor windows | Built-in onion routing |
| Web3 wallet | Ethereum/Solana native |

## CLI / Automation

```bash
# Launch with remote debugging (for Hermes browser tools)
/Applications/Brave\ Browser.app/Contents/MacOS/Brave\ Browser --remote-debugging-port=9222

# User data directory
~/Library/Application Support/BraveSoftware/Brave-Browser/
```

## Integration with Hermes

- `cua-driver` supports Brave (com.brave.Browser) for computer use
- `page` tool works via JavaScript from Apple Events
- Bundle ID: `com.brave.Browser`
