---
title: Legacy UI — Frontend Architecture & Implementation
domain: 02 Infrastructure (The Body)
status: scaffolded
created: 2026-05-25
tags: [onemind-os, legacy-ui, frontend, livekit, pwa, dashboard]
---

# Legacy UI — Frontend Architecture

## Location

```
~/workspace/legacy-ui/
```

## Status: ✅ Built & Compiling

```
npm install && npm run dev  →  http://localhost:3000
npm run build              →  ✅ passes (6 routes, 146KB first load)
```

---

## What It Is

A desktop-first PWA command center for Legacy AI — the visual interface to the omnipresent agent. Dark HUD aesthetic inspired by Mark-XXXIX (FatihMakes/Mark-XXXIX), rewritten from PyQt6 to Next.js for cross-platform (desktop, mobile, office display).

**NOT a chatbot UI. A command center.**

---

## Tech Stack

| Layer | Technology | Why |
|---|---|---|
| Framework | Next.js 15 (App Router) | SSR, PWA, fast |
| Language | TypeScript | Type safety |
| Styling | Tailwind CSS v4 | Mark-XXXIX color system |
| Voice | @livekit/components-react + livekit-client | WebRTC voice room |
| Animation | framer-motion + Canvas API | 60fps HUD orb |
| Realtime data | nats.ws | Direct NATS WebSocket for dashboard |
| Future | Tauri wrapper | Native desktop app, transparent window |

---

## Color Palette (Mark-XXXIX)

```
BG:       #00060a    — deep space black
PANEL:    #010d14    — panel background
BORDER:   #0d3347    — subtle borders
PRI:      #00d4ff    — primary cyan (Legacy's color)
PRI_DIM:  #007a99    — dimmed cyan
ACC:      #ff6b00    — orange accent (alerts, active)
ACC2:     #ffcc00    — yellow accent (warnings, thinking)
GREEN:    #00ff88    — success, active, online
RED:      #ff3355    — errors, critical
TEXT:     #8ffcff    — primary text
TEXT_DIM: #3a8a9a    — secondary text
```

---

## Layout (3-Column)

```
┌──────────────────────────────────────────────────────────────┐
│  L.E.G.A.C.Y                                    HH:MM:SS     │ ← Header
├────────┬────────────────────────────────┬────────────────────┤
│ SYS    │                                │  ACTIVITY LOG      │
│ MONITOR│         HUD ORB                │  (scrolling)       │
│        │    (animated canvas)           │                    │
│ CPU ▓▓ │         ◉                      │  WATCHES           │
│ MEM ▓▓ │     ╱   ╲                     │  ● bar door        │
│ GPU ▓▓ │    ╲   ╱                      │  ● drones          │
│ NET ▓▓ │     LISTENING                  │  ● sarah home      │
│ TMP ▓▓ │                                │                    │
│        │  [voice controls]              │  [text input]      │
│ UP 04h │                                │                    │
├────────┴────────────────────────────────┴────────────────────┤
│  [F4] Mute  ·  [F11] Fullscreen  ·  HERMES CONNECTED        │ ← Footer
└──────────────────────────────────────────────────────────────┘
```

---

## Components

| Component | What it does |
|---|---|
| `HudOrb.tsx` | Animated canvas: 3 spinning arc rings, pulse rings, particles, crosshair, corner brackets. Reacts to agent state (idle/listening/speaking/thinking) |
| `SystemPanel.tsx` | CPU/MEM/GPU/NET/TEMP bars + uptime. Polls `/api/metrics` every 2s |
| `MetricBar.tsx` | Single metric bar with label, value, colored fill |
| `ActivityLog.tsx` | Auto-scrolling log. Color-coded: SYS=cyan, You=white, Legacy=green, WATCH=orange, ERR=red |
| `WatchPanel.tsx` | Active watches with status indicators (● idle, ◉ triggered, ○ dismissed) |
| `VoiceControls.tsx` | LiveKit room connect, mute toggle, VAD indicator |
| `TextInput.tsx` | Command input (monospace, cyan border on focus) |
| `Header.tsx` | "L.E.G.A.C.Y" title + live clock + date |
| `Footer.tsx` | Status bar with connection indicators + hotkey hints |

---

## Data Sources

### 1. LiveKit Room (Voice)
```typescript
// Connect to Legacy's voice room
const room = new Room();
await room.connect(LIVEKIT_URL, token);
// Audio streams automatically via WebRTC
// Data channel carries structured UI commands from agent
```

### 2. NATS WebSocket (Dashboard)
```typescript
// Direct subscription to NATS for real-time data
import { connect } from 'nats.ws';
const nc = await connect({ servers: 'ws://localhost:9222' });

// Subscriptions:
nc.subscribe('events.ha.>');        // HA state changes
nc.subscribe('watches.>');           // Watch registry changes
nc.subscribe('agent.status');        // Agent state changes
nc.subscribe('cmd.drone.*.telemetry'); // Drone positions
```

### 3. API Routes (System Metrics + Agent Text)
```
GET  /api/metrics  → { cpu, mem, gpu, net, temp, uptime }
POST /api/agent    → { message: string } → { response: string }
GET  /api/agent    → { state: AgentState, watches: Watch[] }
```

---

## HUD Orb Animation Details

Canvas-based 60fps animation with these layers (bottom to top):

1. **Grid dots** — subtle background texture
2. **Glow halo** — radial gradient, intensity = f(state)
3. **Pulse rings** — expanding circles, fade as they grow
4. **Arc rings** — 3 concentric spinning arcs (different speeds/directions)
5. **Scan arcs** — two sweeping arcs (cyan + orange)
6. **Tick marks** — 360° compass-style marks
7. **Crosshair** — center cross with gap
8. **Corner brackets** — frame the orb
9. **Core orb** — gradient circle, scale = f(state)
10. **Particles** — spawn from edge when speaking, drift outward
11. **State text** — "● LISTENING" / "● SPEAKING" / "◈ THINKING"
12. **Waveform** — audio level bars below orb

States:
- **Idle**: gentle pulse, slow arcs, dim glow
- **Listening**: medium glow, moderate arc speed, green status
- **Speaking**: intense glow, fast arcs, particles, orange status, waveform active
- **Thinking**: rotating arcs accelerate, yellow status, blink pattern

---

## What Was Taken From Mark-XXXIX

| Element | Adapted How |
|---|---|
| Color palette (C class) | → Tailwind theme tokens |
| HudCanvas (PyQt6) | → Canvas API + framer-motion (HudOrb.tsx) |
| MetricBar (PyQt6) | → React + Tailwind (MetricBar.tsx) |
| _SysMetrics (psutil) | → /api/metrics route (Node os module) |
| LogWidget | → ActivityLog.tsx (auto-scroll, color-coded) |
| MainWindow layout | → 3-column grid (page.tsx) |
| Header (clock, title) | → Header.tsx |
| File drop zone | → Future: drag-and-drop for context injection |
| Setup overlay | → Not needed (Hermes handles auth/config) |

---

## What's Different From Mark-XXXIX

| Mark-XXXIX | Legacy UI |
|---|---|
| Gemini only, forever | Any model via Hermes provider system |
| PyQt6 (desktop only) | Next.js PWA (any device) |
| No memory persistence | Hermes memory (cross-session, cross-device) |
| No background watches | Full watch dashboard + NATS events |
| No voice transport choice | LiveKit (WebRTC, console, or server) |
| Single file JSON memory | Hermes memory + session search + skills |
| No NATS | NATS everywhere (HA, drones, services) |
| No multi-device | Same UI on desktop, phone, display, glasses |

---

## Next Steps

1. Wire LiveKit voice connection to actual backend agent
2. Connect NATS WebSocket for live watch/HA data
3. Add drone map panel (MapLibre + telemetry)
4. Add camera tile panel (video tracks from agent)
5. PWA service worker for offline shell
6. Tauri wrapper for native desktop (transparent window mode)
7. Mobile-responsive layout (voice-first on small screens)

---

## Related

- [[OneMind-Realtime-Core]] — backend architecture
- [[Unified-Agent-Runtime]] — Hermes + LiveKit integration
- [[HA-NATS-Bridge]] — home data source
