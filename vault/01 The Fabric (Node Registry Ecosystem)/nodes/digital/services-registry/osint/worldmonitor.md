---
title: "WorldMonitor"
type: service-card
status: active
created: 2026-06-14
updated: 2026-06-14
owner: zeus
category: osint
tags: [osint, situational-awareness, geospatial, monitoring, self-hosted, vps, docker]
ai_summary: "Open-source geospatial situational-awareness dashboard (aircraft, vessels, conflict, markets). Self-hosted on omos-cloud (Lightsail) Docker at :3002. Forked to onemindos/worldmonitor."
---

# WorldMonitor

> Open-source real-time situational-awareness dashboard — live aircraft, maritime vessels, conflict events, and market signals on a global map. Upstream reference app that WorldWideView (WWV) is being built to mimic.

## Role in OneMind OS
- Live global situational-awareness / OSINT feed
- **Reference architecture** for the in-house WWV build (`wwv-src` mimics this app)
- One of three OSINT dashboards on omos-cloud (with WWV + Osiris)

## Hosting
| Field | Value |
|-------|-------|
| Host | **omos-cloud** (AWS Lightsail) |
| Tailscale IP | 100.111.189.116 |
| Public IP | 3.224.147.141 |
| Path on host | `/home/ubuntu/worldmonitor` |
| Runtime | Docker Compose |
| URL (mesh) | `http://100.111.189.116:3002` |
| Port | 3002 → 8080 (container) |
| Health | `healthy` (Docker), app serves 200 |

## Containers
| Container | Role | State |
|-----------|------|-------|
| worldmonitor | main app (nginx :8080) | healthy |
| worldmonitor-redis | cache | up |
| worldmonitor-redis-rest | REST shim | up |
| worldmonitor-nats-bridge | NATS bridge | up |
| ~~worldmonitor-ais-relay~~ | maritime AIS feed | **disabled** — needs `AISSTREAM_API_KEY` (aisstream.io); was crash-looping |

## Source / Backup
| Field | Value |
|-------|-------|
| Upstream | `github.com/koala73/worldmonitor` (Elie Habib) |
| **Sovereign fork** | **`github.com/onemindos/worldmonitor`** (forked 2026-06-14, tracks upstream) |
| Local clone | `/Users/onemind/worldmonitor` (2.3G — reference/study copy, re-clonable) |

## Notes / Open Items
- Most market/data panels are empty pending API keys (FRED, Finnhub, etc.) — keyless run by design.
- AIS maritime relay disabled until an AISSTREAM key is provisioned.
- Healthcheck fixed 2026-06-14: was failing on IPv6 `localhost` mismatch → repointed to `127.0.0.1`.

---

## Related
- [[00-24 SO (Sovereign Operations)/01 The Fabric (Node Registry)/01 The Fabric|The Fabric]]
- [[worldwideview]] · [[osiris]]
