---
title: "Osiris"
type: service-card
status: active
created: 2026-06-14
updated: 2026-06-14
owner: zeus
category: osint
tags: [osint, situational-awareness, sanctions, intelligence, nextjs, self-hosted, vps, docker]
ai_summary: "Open-source OSINT situational-awareness dashboard with an OFAC sanctions intelligence layer (19,786 entities). Deployed to omos-cloud Docker at :3003 (built from source). Forked to onemindos/osiris."
---

# Osiris

> Open-source OSINT situational-awareness dashboard — aviation, camera feeds, and an intelligence layer indexing OFAC/OpenSanctions data. Palantir-style HUD ("Ghost/Core" modes).

## Role in OneMind OS
- OSINT dashboard with a dedicated **sanctions intelligence resolver** (`/resolve`)
- Third of three OSINT dashboards on omos-cloud (with WorldMonitor + WWV)
- Newest deploy — stood up net-new on omos-cloud 2026-06-14

## Hosting
| Field | Value |
|-------|-------|
| Host | **omos-cloud** (AWS Lightsail) |
| Tailscale IP | 100.111.189.116 |
| Public IP | 3.224.147.141 |
| Path on host | `/opt/omos/osiris` |
| Runtime | Docker Compose (**built from source** — ghcr image needs auth) |
| URL (mesh) | `http://100.111.189.116:3003` |
| Ports | 3003 → 3000 (app) · 8080 (cache) · 4000 (intel) |
| Health | app serves 200; intel layer `healthy` |

## Containers
| Container | Role | State |
|-----------|------|-------|
| osiris | Next.js app | up |
| osiris-cache | nginx:alpine cache | up |
| osiris-intel | intelligence layer (:4000) — 19,786 OFAC SDN entities indexed | healthy |

## Source / Backup
| Field | Value |
|-------|-------|
| Upstream | `github.com/simplifaisoul/osiris` |
| **Sovereign fork** | **`github.com/onemindos/osiris`** (forked 2026-06-14, tracks upstream) |
| Local clone | `/Users/onemind/osiris` (647M — re-clonable) |

## Notes / Open Items
- Runs **keyless** — most feeds work; Gemini/FIRMS/AIS panels dark without keys.
- Deploy required creating an external `umami_default` Docker network (compose references it) and `OSIRIS_PORT=3003` to avoid colliding with langfuse (:3000).
- intel-layer healthcheck fixed 2026-06-14 (IPv6 `localhost` → `127.0.0.1`).

---

## Related
- [[00-24 SO (Sovereign Operations)/01 The Fabric (Node Registry)/01 The Fabric|The Fabric]]
- [[worldmonitor]] · [[worldwideview]]
