---
title: "WorldWideView (WWV)"
type: service-card
status: active
created: 2026-06-14
updated: 2026-06-14
owner: zeus
category: osint
tags: [osint, situational-awareness, geospatial, cesiumjs, nextjs, self-hosted, vps, docker]
ai_summary: "Open-source plugin-driven geospatial intelligence engine (Next.js 16 + CesiumJS 3D globe). Running on omos-cloud Docker at :3001 (prebuilt image). Forked to onemindos/worldwideview."
---

# WorldWideView (WWV)

> Open-source, plugin-driven geospatial intelligence engine — live data streams rendered as cinematic layers on a CesiumJS 3D globe.

## Role in OneMind OS
- Live global situational-awareness / OSINT feed on a 3D globe
- One of three OSINT dashboards on omos-cloud (with WorldMonitor + Osiris)

## Hosting
| Field | Value |
|-------|-------|
| Host | **omos-cloud** (AWS Lightsail) |
| Tailscale IP | 100.111.189.116 |
| Public IP | 3.224.147.141 |
| Path on host | `/opt/omos/worldwideview` |
| Runtime | Docker Compose (prebuilt image) |
| Image | `ghcr.io/silvertakana/worldwideview:latest` |
| URL (mesh) | `http://100.111.189.116:3001` |
| Port | 3001 → 3000 (container) |
| Health | `healthy` (Docker); app status `degraded` (data engine offline — separate service, non-blocking) |

## Containers
| Container | Role | State |
|-----------|------|-------|
| worldwideview-wwv-1 | Next.js app | healthy |
| worldwideview-db-1 | PostgreSQL 15 | healthy |
| worldwideview-redis-1 | Redis 7 (added 2026-06-14) | healthy |

## Source / Backup
| Field | Value |
|-------|-------|
| Upstream | `github.com/silvertakana/worldwideview` (Silvertakana) |
| **Sovereign fork** | **`github.com/onemindos/worldwideview`** (forked 2026-06-14, tracks upstream) |
| Local clone | `/Users/onemind/wwv-src` (1.5G — re-clonable, safe to remove) |

## Notes / Open Items
- Redis added to compose 2026-06-14 (app expected `localhost:6379`, none shipped) → `REDIS_URL=redis://redis:6379`.
- App health shows `engine:false` — WWV data-engine service (`:5001`) is a separate optional component, not deployed. Dashboard serves fine without it.
- Healthcheck fixed 2026-06-14: IPv6 `localhost` mismatch → repointed to `127.0.0.1`.

---

## Related
- [[00-24 SO (Sovereign Operations)/01 The Fabric (Node Registry)/01 The Fabric|The Fabric]]
- [[worldmonitor]] · [[osiris]]
