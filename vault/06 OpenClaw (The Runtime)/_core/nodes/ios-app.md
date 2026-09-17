---
title: "Interface — iOS App (OpenClaw Native)"
type: reference
status: active
created: 2026-04-12
owner: zeus
tags: [type/reference, domain/SO, domain/openclaw, interface/ios]
ai_summary: "OpenClaw native iOS app — setup, Tailscale access, capabilities, push notifications."
---

# Interface: iOS App (OpenClaw Native)

**Status:** ✅ Available
**Type:** Native mobile app
**Primary use:** Full agent UI on iPhone — chat, sessions, push notifications

---

## Setup

1. Download OpenClaw iOS app from releases / TestFlight
2. Configure server URL: `http://100.77.52.93:443` (Tailscale required)
3. Enter access key (from VPS `.env.local` or 1Password → SO vault → OpenClaw Access Key)
4. Connect to any agent workspace

---

## Access Requirements

- **Tailscale must be active** on iPhone — `100.77.52.93` is the Tailscale IP of the VPS
- Alternative: Configure a Cloudflare Tunnel for a public HTTPS endpoint (no Tailscale needed)
- Auth key is set once during configuration

---

## Capabilities

- Full chat interface with all 9 agents
- Push notifications via APNS when app is backgrounded
- Voice input → auto-transcribed → sent to agent
- File picker → attach images, PDFs, documents to messages
- Session history — bookmark and resume conversations
- Agent switcher in-app

---

## Notes

- Use for on-the-go access when Mac is unavailable
- For quick messages to Legacy, prefer Telegram — faster with no app launch
- iOS app is best for longer structured sessions, reviewing history, file sharing
- Push notifications from Legacy fire even when app is closed (APNS background push)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
