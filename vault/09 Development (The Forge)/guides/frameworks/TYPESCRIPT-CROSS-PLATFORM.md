# TypeScript Cross-Platform: Best Options for OneMind OS

> Decision doc created 2026-03-07. Revisit as ecosystem evolves.

## The Top Contenders

| Framework | Targets | Language | Native? | How It Works |
|---|---|---|---|---|
| **React Native + Expo** | iOS, Android, Web | TypeScript | Yes (native views) | Compiles to native UI components. Expo adds OTA updates, push notifications, easy builds |
| **Tauri** | macOS, Windows, Linux, iOS, Android | TypeScript (frontend) + Rust (backend) | Hybrid | Web view for UI, Rust for system access. Tiny bundle size (5-10MB vs Electron's 150MB) |
| **Capacitor (Ionic)** | iOS, Android, Web | TypeScript | Hybrid | Web app wrapped in native WebView. Access to native APIs via plugins |
| **Electron** | macOS, Windows, Linux | TypeScript | No (WebView) | Chromium + Node.js. Heavy (~150MB) but battle-tested |
| **NativeScript** | iOS, Android | TypeScript | Yes (native views) | Direct native API access from TS |

## Recommendation: Expo (React Native)

**Expo is the best fit for OneMind OS.** Reasons:

1. **Mobile-first** — Real native iOS + Android apps from TypeScript
2. **Same ecosystem** — Shares the React component model with Next.js (SwarmClaw/openclaw-studio). Skills transfer directly
3. **Hardware access** — Camera, voice, sensors, location, motion, contacts, calendar — Expo has mature APIs for all of these (exactly what the command-center Swift app does today)
4. **Web too** — Expo supports web targets: one codebase → iOS + Android + Web
5. **OTA updates** — Push UI updates without App Store review
6. **OpenClaw gateway** — WebSocket client works the same in React Native as in Next.js

## The Strategy: Two Codebases, Same Language

```text
TypeScript Everywhere:
├── SwarmClaw (Next.js)            → Admin/orchestration dashboard (web)
├── Expo App (React Native)        → Mobile command center (iOS/Android)
│   ├── Chat with Legacy
│   ├── Voice (like current Swift Voice module)
│   ├── Camera/Vision
│   ├── Location/Motion sensors
│   └── Gateway node connection
├── LobsterBoard (Node.js)         → Monitoring dashboard (web)
└── Shared:
    ├── OpenClaw Gateway client     → Same WS protocol lib
    ├── NATS client (nats.ws)       → Works in both Node.js and React Native
    └── Theme tokens                → Same solarpunk colors
```

## Language Comparison

| Approach | Languages | Effort |
|---|---|---|
| **Previous** (Swift + Flutter + Next.js) | Swift + Dart + TypeScript | 3 languages, 3 build systems, 3 mental models |
| **Expo + Next.js** (recommended) | TypeScript only | 1 language, shared React patterns, shared OpenClaw client |
| **Flutter for everything** | Dart only | 1 language, but fighting the OpenClaw (JS) ecosystem |

## Integration with OneMind Stack

- **OpenClaw Gateway**: WebSocket + JSON at port 18789. Any TS client sends `connect` + `chat.send` and gets agent responses. `GatewayBrowserClient` already exists in openclaw-studio
- **NATS**: `nats.ws` npm package works in both Node.js and React Native for direct JetStream access
- **PostgreSQL**: Server-side only (via API routes in Next.js or Convex Actions)
- **Agno AgentOS**: HTTP calls from Next.js API routes or Expo via fetch

## Convex (Optional Future Backend)

Convex is a reactive TypeScript BaaS (not a UI framework). Could replace/supplement PostgreSQL:

- Real-time reactivity (auto-updating queries via WebSocket)
- Built-in AI Agent component with threads, persistent messages, tool calls
- TypeScript server functions (no SQL)
- Client libs for React, React Native, Swift, Python
- Tradeoff: adds cloud dependency, no native pgvector equivalent

**Current recommendation**: Skip Convex. PostgreSQL + NATS already working.

## Decision

- **Primary UI platform**: TypeScript (Next.js for web, Expo for mobile)
- **Keep Swift ios-app**: As specialized hardware bridge (deep system integration)
- **Flutter**: Archive as design reference; stop active development
- **OpenClaw Studio**: Learning reference for gateway protocol patterns
- **SwarmClaw**: Deploy as admin/orchestration layer
- **LobsterBoard**: Deploy as monitoring dashboard


---

## Related

- [[00-24 SO (Sovereign Operations)/05 Development (Code & Testing) (The Forge)/05 Development|05 Development]]
