# omos-ui — Command Interface

> Last updated: 2026-07-11
> Repo: github.com/onemindos/omos-ui (private)
> Live: https://onemindos.dev

## Purpose

The sovereign operations center for OneMind OS. A NATS-native single-page application that renders the entire fabric as an operating system interface — agents, telemetry, missions, intel, and fleet management in one dark ops-center UI.

## Stack

| Layer | Choice | Why |
|-------|--------|-----|
| Framework | React 19 + Vite 7 | Bleeding edge, fast HMR |
| Language | TypeScript 5.8 (strict) | Full type safety |
| Styling | Tailwind 4 + CSS custom properties | Dark theme, no component lib overhead |
| Routing | TanStack Router | Type-safe, lazy-loaded per module |
| Server state | TanStack Query 5 | Caching, background refetch |
| Client state | Zustand (per-module) + Jotai (NATS atoms) | Minimal, fast, composable |
| Real-time | nats.ws (WebSocket) | Direct bus connection — no REST middleman |
| Charts | uPlot + ECharts | Time-series perf + rich visualization |
| Tables | TanStack Table + Virtual | Virtualized for large datasets |
| Editor | Monaco | In-app code/config editing |
| Terminal | xterm.js | Embedded shell |
| Tests | Vitest + Testing Library | Fast, jsdom, coverage |
| Lint | ESLint 9 (flat) + Prettier | Consistent quality |
| Pre-commit | Husky + lint-staged | Gate before push |
| Desktop | Tauri 2 | Native app (macOS/Win/Linux) |
| Mobile | Capacitor 7 | Android APK |

## Architecture

```
src/
├── main.tsx              Entry (router + query + ErrorBoundary)
├── index.css             Tailwind + CSS vars (dark theme)
├── routeTree.gen.tsx     All routes (generated, lazy-loaded)
├── shell/                OS chrome
│   ├── Shell.tsx         Root layout (TopBar + NavRail + Outlet)
│   ├── TopBar.tsx        Connection status, user, clock
│   ├── NavRail.tsx       Collapsible sidebar (sections + icons)
│   └── nats/             NatsManager, NatsContext, connection atoms
├── stores/               Global stores (NATS connection, profiles)
├── shared/
│   ├── components/       Card, Badge, Button, StatusDot, ErrorBoundary
│   ├── hooks/            useNats(), useTheme()
│   └── types/            Cross-module types (Agent, Entity, Telemetry, NATS)
├── lib/                  Utilities (cn(), formatters)
├── pages/                Module-per-directory (20 modules)
└── test/                 Test setup + mocks
```

### Module Pattern

Every page module follows:
```
pages/{module}/
├── {Module}Page.tsx      Route component (lazy-loaded)
├── components/           Module-specific UI
├── hooks/                Custom hooks (data fetching, subscriptions)
├── stores/               Zustand stores for this module
└── types.ts              Module-specific types
```

## Modules (20)

### COMMAND Section
| Module | Subject | Purpose |
|--------|---------|---------|
| overview | `$SYS.REQ.SERVER.PING` | Dashboard: bus metrics, agent status, system health |
| services | `meta.services.*` | Service directory + health checks |
| bus | `$SYS.>` | NATS topology, connections, leafnodes |
| agents | `agents.>` | Agent fleet: roster, chat panel, wire traffic, session logs |
| missions | `missions.>` | Mission planner + task dependency graph |

### OPERATIONS Section
| Module | Subject | Purpose |
|--------|---------|---------|
| integrations | — | External service configs |
| nats (streams) | `$JS.API.>` | Streams, KV, subject explorer, publish tool |
| intel | `intel.>`, `orp.>` | ORP entities, world feeds, threat board, AI analyst |
| replay | `timemachine.>` | Historical event replay from ClickHouse |

### SYSTEM Section
| Module | Subject | Purpose |
|--------|---------|---------|
| security | — | Access control, keys, certificates |
| wire | `agents.>` | Raw NATS traffic inspector (agenttap-in-browser) |
| terminal | — | xterm.js embedded shell |
| voice | `agents.voice.>` | Push-to-talk, transcript, voice bridge |
| settings | — | Theme, connections, notifications, profile |

### STANDALONE
| Module | Subject | Purpose |
|--------|---------|---------|
| map | `tak.>` | TAK COP (common operating picture) |
| fleet | `telemetry.>` | Device/node management + telemetry cards |
| tasks | `tasks.>` | Task engine UI (om_task integration) |
| connect | — | Onboarding new devices/agents |
| timemachine | `timemachine.>` | ClickHouse query builder + results |

## Deployment

```
Push to feat/agents → CI (quality gate) → GHCR → DOKS rollout
                      ↓
              typecheck + lint + format + test
```

- **Namespace:** `omos` on DOKS
- **Image:** `ghcr.io/onemindos/omos-ui:latest`
- **Serving:** nginx (SPA fallback to index.html)
- **Domain:** `onemindos.dev` (Cloudflare DNS → DO Load Balancer)

## Environment

| Variable | Value |
|----------|-------|
| VITE_NATS_WS_URL | wss://nats-ws.onemindos.dev |
| VITE_NATS_USER | onemind |
| VITE_CLOUDTAK_URL | https://cloudtak.onemindos.dev |
| VITE_WWV_URL | https://wwv.onemindos.dev |
| VITE_TIMEMACHINE_URL | https://timemachine.onemindos.dev |
| VITE_APP_NAME | OneMind OS |

## Developer Workflow

```bash
npm run dev             # Vite dev server :5173
npm run validate        # typecheck + lint + test (full gate)
npm run test:watch      # TDD mode
npm run format          # Prettier format all
```

Branch from `feat/agents` (default) → PR → quality gate → merge → auto-deploy.

## Integration Points

- **Legacy (me):** I build features, review PRs, delegate to Forge for heavy coding
- **Forge:** Primary coding agent for omos-ui feature work
- **NATS bus:** UI connects directly via WebSocket — zero REST APIs
- **CloudTAK:** JWT auth shared between WebTAK and omos-ui
- **ClickHouse:** TimeMachine queries for historical data
- **Tauri:** Desktop app auto-update via GitHub Releases
- **Capacitor:** Android APK loads from onemindos.dev origin

## Key Design Decisions

1. **NATS-native** — no REST backend. The UI is a NATS client.
2. **Zero component library** — custom primitives (Card, Badge, Button) for full control
3. **Module isolation** — each page is self-contained with its own state
4. **Lazy everything** — code-split per route, virtualized tables
5. **Dark-only** — ops-center aesthetic, high contrast for situational awareness
6. **Desktop + Mobile** — same codebase ships web, Tauri, and Capacitor
