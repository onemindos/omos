---
title: "Interface — macOS App (OpenClaw Native)"
type: reference
status: active
created: 2026-04-12
owner: zeus
tags: [type/reference, domain/SO, domain/openclaw, interface/macos]
ai_summary: "OpenClaw native macOS app — setup, local and Tailscale access, capabilities."
---

# Interface: macOS App (OpenClaw Native)

**Status:** ✅ Available
**Type:** Native desktop app
**Primary use:** Full agent UI on Mac — menu bar access, drag-drop files, keyboard shortcuts

---

## Setup

1. Download OpenClaw macOS app from releases
2. Configure server:
   - Remote VPS: `http://100.77.52.93:443` (Tailscale required)
   - Local: `http://localhost:443` (if running locally)
3. Enter access key (1Password → SO vault → OpenClaw Access Key)

---

## Capabilities

- Full desktop-optimized chat interface with all 9 agents
- **Menu bar icon** → one-click to open Legacy chat from anywhere
- Native macOS notifications (alerts, briefings)
- Keyboard shortcuts throughout the UI
- Drag-and-drop file attachment (images, PDFs, documents)
- Multiple agent sessions in tabs
- Copy markdown, export sessions

---

## Notes

- Best for deep work sessions, long complex requests, file processing
- Menu bar integration makes the macOS app the fastest way to reach Legacy without Telegram
- Use alongside VS Code for coding sessions (drag file → agent processes it)
- Tailscale must be active for VPS access, or run OpenClaw locally


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
