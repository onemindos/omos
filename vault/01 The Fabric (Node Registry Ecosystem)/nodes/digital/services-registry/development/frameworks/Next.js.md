# Next.js

> **React framework powering onemind-web — the One Mind command center**

## Overview

Next.js is the full-stack React framework at the core of onemind-web. Provides server-side rendering, API routes, file-based routing, and edge functions. Built by Vercel; version 16.1.7 in production. The primary web interface for the One Mind platform.

## Status

| Attribute   | Value                                    |
| ----------- | ---------------------------------------- |
| Status      | ✅ Active                                 |
| Category    | Development / Framework                  |
| Type        | Full-Stack React Framework               |
| Integration | React 19, TypeScript, Tailwind, Vercel   |

## Use Cases in OneMind

1. **onemind-web** — Primary AI command center UI (port 3456)
2. **API Routes** — `/app/api/` routes for chat, agents, connectors
3. **Agent Chat UI** — Real-time streaming chat with OpenClaw agents
4. **Office Visualization** — 2D/3D floor plan in `src/office/`
5. **CLI Entry** — `bin/onemind.js` wraps Next.js for CLI access

## Stack Details

| Component    | Version / Tech                   |
| ------------ | -------------------------------- |
| Next.js      | 16.1.7 (standalone output)       |
| React        | 19                               |
| TypeScript   | Strict mode                      |
| Tailwind CSS | v4                               |
| State        | Zustand 5                        |
| DB           | better-sqlite3 (55 tables)       |

## Key Commands

```bash
npm run dev              # Turbopack dev server → localhost:3456
npm run build:ci         # Production build (NEXT_DISABLE_ESLINT=1)
npm run lint:baseline    # Lint gate (must pass before release)
```

## Repo

- onemind-web: `/Users/zeuslegacy/Documents/onemind-web/`

## Links

- Website: [nextjs.org](https://nextjs.org)
- Docs: [nextjs.org/docs](https://nextjs.org/docs)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
