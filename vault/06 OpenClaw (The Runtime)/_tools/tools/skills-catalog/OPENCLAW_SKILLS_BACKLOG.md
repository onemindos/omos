# OpenClaw Skills Backlog

> Extracted from capabilities analysis and effort estimates (Feb 2026).
> Stack-agnostic skill targets for Legacy + OpenClaw + NATS.

---

## Memory Architecture (OpenClaw `memory/` folder structure)

Four memory types — map these to markdown files in `openclaw/workspace/memory/`:

| Type | What It Stores | File |
|------|---------------|------|
| **User Profile** | Name, timezone, Four Pillars (HP/LE/GE/UI), preferences, communication style, focus times | `user-profile.md` |
| **Entity Memory** | People (contacts, relationships, facts), Projects, Companies | `entities/` folder per entity |
| **Agent Culture** | Principles ranked by priority, communication style rules, behavioral constraints | `culture.md` (in SOUL.md) |
| **Session Memory** | Last 5 session summaries, user memories, decisions, entity tracking | `memory/YYYY-MM-DD.md` |

---

## Skills to Build — Effort Estimates

| Skill | Est. Effort | Priority | Notes |
|-------|------------|----------|-------|
| **Telegram Bot** | 3h | 🔴 High | Code pattern in `docs/OPENCLAW_SKILL_PATTERNS.md` |
| **Discord Bot** | 3h | 🔴 High | Code pattern in `docs/OPENCLAW_SKILL_PATTERNS.md` |
| **Browser Automation** | 4h | 🟡 Medium | Playwright; blueprint in `openclaw/tools/browser.py` |
| **Self-Modifying Skills** | 6h | 🟢 Advanced | Legacy writes its own skills via AI generation |
| **Skill Auto-Discovery** | 4h | 🟡 Medium | Reads YAML frontmatter from skills/ dir |
| **WhatsApp Integration** | 5h | 🟡 Medium | Via OpenClaw channel config |
| **Transit / Uber** | 4h | 🔴 High | Blueprints in `openclaw/tools/transit.py` + `uber.py` |
| **Lob (Physical Mail)** | 3h | 🟡 Medium | Blueprint in `openclaw/tools/lob.py` |
| **LiveKit Node (Smart Glasses)** | 6h | 🔴 High | Patterns in `docs/LIVEKIT_NODE_PATTERNS.md` |
| **Home Assistant** | 4h | 🔴 High | Direct HASS REST/WebSocket |
| **Frigate Camera** | 4h | 🟡 Medium | Local AI camera events via NATS |
| **FarmBot** | 5h | 🟢 Low | HTTP API for garden automation |
| **Oura / Apple Watch** | 4h | 🟡 Medium | Health biometrics → HP Agent |
| **Kokoro TTS** | 3h | 🔴 High | Local voice output |
| **openWakeWord** | 3h | 🔴 High | Local wake word detection |

---

## Integration Catalog (40+ tools — skills to build)

### Productivity
- Calendar (Google, Apple)
- Task management (Todoist — skill exists)
- Email (Gmail, IMAP)
- GitHub (skill exists)
- Obsidian notes (skill exists)

### Smart Home & IoT
- Home Assistant (lights, climate, locks, energy)
- Frigate NVR (cameras, facial recognition, motion zones)
- MQTT broker (generic IoT device control)
- FarmBot (garden automation)
- Solar/battery monitoring

### Health & Wellness
- Oura Ring (sleep, HRV, readiness score)
- Apple Watch / HealthKit
- Nutrition logging
- Workout tracking

### Communication
- WhatsApp (OpenClaw native channel)
- Telegram bot
- Discord bot
- iMessage (OpenClaw native — Apple Node)
- Voice (LiveKit + ElevenLabs + Whisper)

### Research & Data
- Web search (DuckDuckGo, Exa, Perplexity)
- Academic (ArXiv, PubMed)
- Financial data (Yahoo Finance)
- Browser automation (Playwright)

### Physical World
- LiveKit (smart glasses, robot video streams)
- Robot NATS bridge (ROS 2 command dispatch)
- Unitree / mobile robot control
- Drone telemetry

### Finance & Business
- Stripe / payment processing
- CRM (Zoho — existing skill)
- Physical mail / checks (Lob)
- Ride booking (Uber)

### Transit & Travel
- Google Maps / transit routing
- Uber ride management
- MTA / public transit

---

## Skill Priority Order (Recommended)

1. Kokoro TTS + openWakeWord — complete the voice loop
2. Transit + Uber — daily utility, blueprints ready
3. LiveKit Node for smart glasses — next major hardware integration
4. Telegram Bot — expand communication channels
5. Home Assistant — smart home control
6. Self-Modifying Skills — Legacy writes its own skills (6h, highest leverage)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
