---
title: "Entity Registry API"
type: reference
status: scaffold
created: 2026-04-14
owner: zeus
tags: [type/reference, domain/SO, status/scaffold]
ai_summary: "Entity Registry REST API — reads markdown entity files and serves them as structured JSON. Port 18810. Not yet deployed."
---

# Entity Registry API

> Static REST API over the `05 Environments/` markdown files.
> **Does not modify any markdown.** Markdown remains the source of truth.

---

## Status

**Scaffold only — not yet deployed to VPS.**

---

## How It Works

```
05 Environments/
  machines/vps.md      ←── reads YAML frontmatter
  machines/mac.md      ←── reads YAML frontmatter
  users/wife.md        ←── reads YAML frontmatter
  ...
        ↓
  api/server.js        parses all files on request
        ↓
  GET /api/entities    returns structured JSON
```

Files are read on every request (no cache). Changes to markdown are reflected immediately.

---

## Endpoints

| Method | Path | Returns |
|--------|------|---------|
| `GET` | `/api/health` | Liveness + registry path |
| `GET` | `/api/entities` | All entities, all types |
| `GET` | `/api/entities/:type` | All entities of one type |
| `GET` | `/api/entities/:type/:name` | Single entity + full body |

**Types:** `users`, `locations`, `machines`, `devices`, `assets`, `digital-assets`

### Examples

```
GET /api/entities/machines
→ [{ _name: "vps", _type: "machines", title: "VPS — onemind-main", ... }]

GET /api/entities/machines/vps
→ { _name: "vps", title: "VPS — onemind-main", subtype: "vps", ... _body: "..." }

GET /api/entities/users
→ [wife, kids, dog as JSON objects]
```

---

## Run Locally

```bash
cd api/
npm install
node server.js
# → Entity Registry API running on port 18810
```

---

## Deploy to VPS (when ready)

```bash
# On VPS: /home/onemind/onemind-os/entity-registry-api/
git pull
npm install --production
# Add systemd service or run behind PM2
```

Port `18810` — does not conflict with existing services (3000, 3456, 4222, 5432, 5678, 18800).

---

## Future: Live Data Stubs (not wired)

Stub endpoints exist in `server.js`, returning `501 Not Implemented`:

| Endpoint | Source | Purpose |
|----------|--------|---------|
| `POST /api/entities/:type/:name/update` | any | Write live data back |
| `POST /api/entities/:type/:name/location` | Bouncie / HA | Update GPS position |
| `GET /api/entities/presence` | Home Assistant | Who/what is home |

Wire these when integrations are ready. The markdown stays unchanged either way.

---

## Folder Structure

```
05 Environments (Entities & Registry)/
├── api/                ← THIS — API scaffold
│   ├── server.js       ← Express server, parse + serve
│   ├── package.json    ← express + gray-matter
│   ├── .env.example    ← PORT + future integration keys
│   └── README.md       ← this file
├── users/              ← markdown (unchanged)
├── locations/          ← markdown (unchanged)
├── machines/           ← markdown (unchanged)
├── devices/            ← markdown (unchanged)
├── assets/             ← markdown (unchanged)
├── digital-assets/     ← markdown (unchanged)
└── templates/          ← markdown (unchanged)
```
