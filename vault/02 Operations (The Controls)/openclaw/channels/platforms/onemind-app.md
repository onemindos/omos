---
title: "Interface — OneMind App (Custom)"
type: reference
status: planned
created: 2026-04-12
owner: zeus
tags: [type/reference, domain/SO, domain/openclaw, interface/onemind-app]
ai_summary: "Custom Zeus-built OneMind app — vision, planned features, stack (onemind-mobile + onemind-web API)."
---

# Interface: OneMind App (Custom — Zeus-Built)

**Status:** 🔄 Planned
**Type:** Custom native app
**Primary use:** Purpose-built One Mind command center beyond standard OpenClaw UI

---

## Vision

A custom Zeus-built interface that unifies ALL One Mind systems in one place — not just agent chat:

- **Dashboard view** — presence (FIND ZEUS MODE), agent status, domain health scores
- **Quick-capture flow** — tap → capture → route to correct domain (the C in CODEX)
- **FIND ZEUS MODE indicator** — live Zeus location, last know zone, active mode state
- **Unified inbox** — tasks (Todoist) + alerts (signals) + messages (Legacy) in one feed
- **Life scores** — PI / WI / LI at a glance
- **Agent switcher** — switch between Legacy / Spartan / Oracle / Haven / Empire with context preview
- **Domain health** — SO / HP / LE / GE status cards with key metrics

---

## Stack

- **Framework:** `onemind-mobile` repo — React Native (Expo)
- **Backend:** `onemind-web` REST API + direct NATS WebSocket for live data
- **Auth:** Same OpenClaw access key as all other interfaces
- **State:** Zustand (mirroring onemind-web patterns)

---

## Relevant Repos

- `~/Documents/onemind-mobile` — primary mobile codebase
- `~/Documents/onemind-web` — backend API and data layer

---

## Difference From Standard iOS/macOS App

| Feature | OpenClaw Native App | OneMind Custom App |
|---------|--------------------|--------------------|
| Agent chat | ✅ | ✅ |
| FIND ZEUS MODE dashboard | ❌ | ✅ planned |
| Life scores (PI/WI/LI) | ❌ | ✅ planned |
| Unified inbox | ❌ | ✅ planned |
| Quick-capture CODEX flow | ❌ | ✅ planned |
| Domain health cards | ❌ | ✅ planned |
| Standard OpenClaw UI | ✅ | Modified / branded |

---

## Notes

- This is Zeus's primary build project for the One Mind UI layer
- Standard OpenClaw iOS/macOS apps remain available as fallbacks during development
- The custom app is not a replacement for Telegram — Telegram stays as the always-on background channel


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
