---
id: so-06-20260514-obsidian-openclaw-plugin
title: "Concept — Obsidian-OpenClaw Plugin (Offgrid Chatbot)"
type: idea
status: concept
created: 2026-04-01
modified: 2026-05-14
owner: zeus
tags: [path/so, domain/06-openclaw, type/idea, topic/obsidian, topic/openclaw, topic/offgrid, topic/plugin, status/concept]
ai_summary: "Obsidian plugin that embeds an OpenClaw chat panel directly in Obsidian. Online: WebSocket to VPS. Offgrid: local queue or Ollama fallback."
source: "Apple Notes — Scan 1, image 6"
---

## Concept: `obsidian-openclaw` Plugin

Embed a full OpenClaw chat panel inside Obsidian (Mac or Mobile). Vault context is already there — plugin is essentially a UI skin over the existing gateway protocol.

```
Obsidian (Mac or Mobile)
  └── OpenClaw Chat Panel
       [You]: How do I...?
       [Legacy]: Based on your codex note in HP/27...
            │
       HTTP/WebSocket
            │
       OpenClaw Gateway (VPS or local)
```

---

## Modes

### Online Mode
Chat panel sends messages to OpenClaw gateway on VPS via WebSocket (already works over Tailscale at `100.77.52.93`).

### Offgrid Mode (two approaches)

1. **Queue messages locally** — Plugin stores messages in `.openclaw-queue.json`, syncs when back online (like iMessage offline behavior)
2. **Local Ollama fallback** — Plugin detects no VPS connectivity, routes to local Ollama with vault as context. Lower quality but works anywhere.

---

## Why It's Feasible

- OpenClaw already has the `obsidian-cli` skill
- The vault is already the knowledge base
- Plugin = UI layer over existing gateway protocol — not a new backend

**Existing precedent:** The gateway already exposes a WebSocket interface. This is a frontend problem, not an infrastructure problem.

---

## Build Considerations

- Obsidian plugins are built in TypeScript (Obsidian Plugin API)
- WebSocket client in TS is straightforward
- Offline queue: write to `.openclaw-queue.json` in vault root (or `_system/`)
- Ollama fallback: detect connectivity → switch endpoint

---

## Next Steps

- [ ] Check existing `obsidian-cli` skill — what does it expose?
- [ ] Explore Obsidian Plugin API for chat/panel UI components
- [ ] Prototype: TypeScript → WebSocket → OpenClaw gateway
- [ ] Design offline queue format

---

## Related

- Domain: SO 06 OpenClaw
- [[06 OpenClaw.md]] — Runtime overview
