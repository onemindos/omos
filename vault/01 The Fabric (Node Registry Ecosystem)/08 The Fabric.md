---
title: "01 The Fabric — Node Registry"
domain: SO
subdomain: "01 The Fabric (Node Registry)"
type: reference
status: active
created: 2026-04-12
updated: 2026-05-03
owner: zeus
tags: [type/reference, domain/so, zone/so, status/active]
ai_summary: "Unified entity registry for the OneMind ecosystem — the single source of truth for every noun in Zeus's world: people (with full profiles), pets, physical assets, machines, devices, digital assets, locations, and robots. Each entity has a markdown file with structured context, access level, and data source links. Includes a scaffolded REST API (api/) and robot registry (dimos, om1)."
agent: Sovereignty
todoist_project_id: 6g6P9vcMfQ73mwXG
zone: SO
review: 2026-Q3
---

# 01 The Fabric (Node Registry)

> 🔀 **RENAMED 2026-06-05:** `08 OM Assets (Entity Registry)` → `01 The Fabric (Node Registry)`.
> The 173-file arsenal (old `05 Stack`) was collapsed in as `nodes/digital/services-registry/`, and a 6th category `economic/` was added.
> **The full doctrine merge** (unifying the 5-category Node System + 11-type Entity System into one 6-category / 18-class model, bridge frontmatter, subject law) is RATIFIED IN PRINCIPLE but **pending the "discuss the rest" pass** — this index still describes the legacy entity-registry layout below until that rewrite lands.

> Every person, pet, vehicle, machine, device, digital asset, robot, location — and now every service and economic entity — in the OneMind ecosystem lives here.
> Legacy loads these files to understand the context of who it's talking to or what it's managing.
>
> **This is the Node layer (The Nouns) of the OneMind architecture.**
> The node IS here (08). The agent KNOWS it (01). The infra RUNS it (02).

---

## What We Use For What

The authoritative map of every function → service in the OneMind stack.

| Function | Service | Notes |
|---|---|---|
| **Calendar** | Apple Calendar / Fastmail | Fastmail CalDAV synced to Apple Calendar |
| **Contacts** | Apple Contacts / Fastmail | Fastmail CardDAV synced to Apple Contacts |
| **Email (primary)** | Fastmail | `zeuslegacy@fastmail.com` — daily use, MCP live |
| **Email (encrypted)** | Proton Mail | `zeus.vaults@pm.me` — sensitive, via Hermes gateway |
| **Files & Storage (local)** | Mac local + OM PRIMARY (6TB drive) | `/Users/onemind/` + `/Volumes/OM PRIMARY` |
| **Files & Storage (cloud sync)** | iCloud Drive | Mac sync layer |
| **Repos** | GitHub | `github.com/onemindos` — primary remote |
| **Repos (self-hosted)** | Forgejo | Self-hosted git — fleet sync |
| **Repos (mirror)** | Codeberg | Secondary mirror |
| **Tasks** | Todoist | MCP live — `TODOIST_API_KEY` |
| **Projects** | Todoist | Projects + sections per pillar |
| **Reminders** | Todoist | Due dates + recurring tasks |
| **Password manager** | 1Password | All credentials — never in plaintext |
| **Encrypted notes/storage** | Proton Drive | Sensitive docs |
| **Knowledge base** | Obsidian (this vault) | `~/onemind/` — source of truth |
| **Messaging (personal)** | iMessage / Telegram | iMessage primary, Telegram for agent interface |
| **Agent interface** | OpenClaw → Hermes | All agent interaction routes here |
| **Smart home** | Home Assistant | Local instance, MCP live |
| **DNS / CDN** | Cloudflare | `zeusdelacruz.com`, `weareonemind.com` |
| **VPN / mesh** | Tailscale | `tail717ba0.ts.net` — all nodes |
| **Payments** | Stripe | Business payments |
| **Banking** | Mercury | Business bank |

---

## What This Is and Why It Exists

The Entity Registry is the **structured inventory of Zeus's world** — every significant thing that agents need to know about is described here as a markdown file.

Before this existed, agents had no structured awareness. If you asked Legacy "what's running on the VPS?" or "remind my wife about dinner", it had no idea who the wife is, what the VPS IP is, or what services exist. It would either hallucinate or fail silently.

Now every entity — the VPS, the Mac, the wife, the cars, the barn, the iPhone — has a dedicated file with structured frontmatter that any agent can load and reason about.

### Two layers working together

**Layer 1 — Markdown files (source of truth, always active)**
Plain `.md` files with YAML frontmatter. OpenClaw's memory search already indexes the whole codex vault, so agents can find and read these files today with zero extra setup. Editing a file is editing the registry. No database, no sync, no migration.

**Layer 2 — REST API scaffold (in `api/`, not yet deployed)**
A thin Express server that reads the same markdown files and serves them as structured JSON. Enables cross-entity queries ("all machines on Tailscale"), live data write-back from Bouncie/Home Assistant, and eventually powering a visual node graph in onemind-web. **The markdown never changes** — the API is a read layer on top.

### Why both and not just one?

| Markdown only | API only |
|---------------|---------|
| ✅ Agents can search it today | ❌ Requires deployment first |
| ✅ Obsidian renders it | ❌ Not human-readable |
| ✅ Git-tracked, zero ops | ❌ Another service to maintain |
| ❌ Can't query across entities | ✅ Structured queries |
| ❌ Can't push live data into it easily | ✅ Webhooks can write back |

The markdown is the floor. The API is the future ceiling. Having both means you're never blocked.

### What it enables (now and later)

| Capability | Status | How |
|------------|--------|-----|
| Agent context loading ("who is my wife?") | ✅ Active | OpenClaw memory search reads `users/wife.md` |
| Multi-user routing (wife texts Legacy) | ✅ Active | `openclaw/MULTI-USER.md` architecture |
| Structured entity queries via API | 🔧 Scaffold | Deploy `api/server.js` on VPS |
| Visual node graph in onemind-web | 🔮 Future | Wire `src/office/` to `/api/entities` |
| Live GPS overlay (cars on map) | 🔮 Future | Bouncie webhook → `POST /api/entities/assets/car-1/location` |
| Home presence ("who's home?") | 🔮 Future | HA webhook → `GET /api/entities/presence` |

---

---

## Entity Types

| Type | Folder | Examples | Notes |
|------|--------|---------|-------|
| **human** | `users (User Accounts)/` | Zeus, Cielo, Liam, Noah, Soluna | Full profiles + registry cards. Chat-capable (access-gated) |
| **pet** | `users (User Accounts)/` | dog | Registry only — no chat |
| **location** | `locations (Physical Locations)/` | main house, barn, office | Physical places |
| **machine** | `machines (Machine Registry)/` | VPS, Mac, server | Compute nodes you SSH/deploy to |
| **device** | `devices (Device Registry)/` | iPhone, tablet | Personal handheld endpoints |
| **physical asset** | `assets (Digital Assets)/` | cars, plane, motorcycle, drone | Vehicles + property |
| **digital asset** | `digital-assets (Crypto & Holdings)/` | domains, licenses, crypto | Intangible owned things |
| **robot** | `robots (Robot Agents)/` | dimos, om1 | Robotic platforms and hardware |

---

## Access Levels

| Level | Code | Who | What |
|-------|------|-----|------|
| **Owner** | L4 | Zeus | Full access — all agents, all memory, all commands |
| **Partner** | L3 | Wife | Home agents, family memory, LE/HP domains |
| **Child** | L2 | Kids | Age-gated — homework, schedules, safe content |
| **Observer** | L1 | Guests | Query-only, no writes, no personal context |
| **Entity** | L0 | Dog, assets | No chat — registry + telemetry only |

---

## Multi-User Routing

Legacy identifies the sender by their connection (Telegram phone number, device token, etc.)
and loads the matching user context file before responding.

```
Incoming message → identify sender → load users/{name}.md → inject context → respond
```

See `SO/02 Intelligence/openclaw/MULTI-USER.md` for full routing architecture.

---

## Templates

Use templates in `templates (Registry Templates)/` when creating new entities.

| Template | Use for |
|----------|---------|
| `templates/human.md` | People (family, guests) |
| `templates/pet.md` | Animals |
| `templates/vehicle.md` | Cars, plane, motorcycle |
| `templates/drone.md` | Drones / UAVs |
| `templates/device.md` | Phones, tablets, wearables |
| `templates/machine.md` | VPS, Mac, servers, NAS |
| `templates/location.md` | Physical places |
| `templates/digital-asset.md` | Domains, licenses, crypto, subscriptions |

---

## Folder Structure

Every subfolder has a parenthetical label that tells you at a glance **what kind of entities live inside**. This is consistent with the labeling convention used across the entire codex — the label in parentheses is the human-readable description; the short name before it is what agents and scripts reference.

```
01 The Fabric (Node Registry)/
│
├── _README.md                                  ← this file — full context
│
├── users (User Accounts)/                      ← humans + pets — the people in Zeus's world
│   ├── zeus (Zeus)/profile.md                  L4 owner — full identity profile
│   ├── cielo (Cielo)/profile.md                L3 partner — full identity profile
│   ├── liam (Liam)/profile.md                  L2 child — identity profile
│   ├── noah (Noah)/profile.md                  L2 child — identity profile
│   ├── soluna (Soluna)/profile.md              L2 child — identity profile
│   └── family (Family Members)/               ← family registry cards
│       ├── wife.md                             L3 partner — registry card (thin)
│       ├── kids.md                             L2 children — registry card (thin)
│       └── dog.md                              L0 registry — fill {{}} fields
│
├── locations (Physical Locations)/              ← real-world places Zeus operates from
│   ├── main-house.md                           primary residence, HA-integrated
│   ├── barn.md                                 outbuilding / vehicle storage
│   └── office.md                               work environment
│
├── machines (Machine Registry)/                 ← compute nodes (SSH/deploy targets)
│   ├── vps.md                                  DigitalOcean (104.248.4.212 / 100.77.52.93)
│   └── mac.md                                  Zeus's MacBook — source of truth
│
├── devices (Device Registry)/                   ← handheld personal endpoints
│   └── (add: iphone.md)
│
├── assets (Digital Assets)/                     ← physical owned things (vehicles, property)
│   ├── car-1.md                                Bouncie GPS tracking
│   ├── car-2.md                                Bouncie GPS tracking
│   ├── motorcycle.md                           optional GPS
│   ├── plane.md                                FAA / FlightAware
│   └── drone.md                                FAA registration
│
├── digital-assets (Crypto & Holdings)/          ← intangible owned things (domains, licenses, crypto)
│   └── (add: domains, licenses, subscriptions)
│
├── robots (Robot Agents)/                       ← robotic platforms and hardware
│   ├── dimos/                                  DIMOS robot platform
│   └── om1/                                    OM1 robot platform
│
├── templates (Registry Templates)/              ← scaffolds for creating new entities
│   ├── human.md
│   ├── pet.md
│   ├── vehicle.md
│   ├── drone.md
│   ├── device.md
│   ├── machine.md
│   ├── location.md
│   └── digital-asset.md
│
└── api (API Keys & Tokens)/                     ← REST API scaffold (NOT YET DEPLOYED)
    ├── server.js                               Express server — reads markdown, serves JSON
    ├── package.json                            dependencies: express + gray-matter
    └── README.md                               full API docs, endpoints, deploy instructions
```

### Folder-by-Folder Breakdown

| Folder | Label | What Lives Here | Count |
|--------|-------|-----------------|-------|
| `users` | User Accounts | Every person and pet in Zeus's life. Subdirectories for full profiles (zeus/, cielo/, liam/, noah/, soluna/), flat `.md` cards for lighter entries (wife.md, kids.md, dog.md). Access levels L0–L4 control what agents can share with whom. | 8 entries |
| `locations` | Physical Locations | Real-world places — houses, outbuildings, offices. Each file includes address, WiFi, Home Assistant area, and what's physically there. | 3 entries |
| `machines` | Machine Registry | Compute infrastructure Zeus owns or controls. VPS, Mac, future NAS. Each file has IPs, Tailscale addresses, specs, and what services run on it. | 2 entries |
| `devices` | Device Registry | Personal handheld endpoints — phones, tablets, wearables. Currently empty (iPhone not yet added). | 0 entries |
| `assets` | Digital Assets | Physical property — vehicles (cars, motorcycle, plane) and drone. Includes GPS/telemetry source info (Bouncie, FAA, FlightAware). | 5 entries |
| `digital-assets` | Crypto & Holdings | Intangible owned things — domains, software licenses, crypto wallets, subscriptions. Currently empty, awaiting population. | 0 entries |
| `robots` | Robot Agents | Robotic platforms — each gets a subdirectory with config, capabilities, and integration notes. Currently: DIMOS and OM1. | 2 entries |
| `templates` | Registry Templates | Blank scaffolds with correct frontmatter for each entity type. Copy one when creating a new entity to ensure consistent structure. | 8 templates |
| `api` | API Keys & Tokens | REST API scaffold that reads the markdown files above and serves them as structured JSON. Express + gray-matter. Not yet deployed — see API section below. | 3 files |

---

## The API Scaffold — What It Is, What It Does

The `api/` folder contains a complete but **undeployed** REST API server.

**It reads the markdown files above and serves them as structured JSON.** Nothing more. It does not store data, does not write to the markdown files, does not require a database. The markdown is always the source of truth.

### API endpoints

```
GET /api/health                      → liveness check
GET /api/entities                    → all entity types as one JSON object
GET /api/entities/:type              → all entities of one type
GET /api/entities/:type/:name        → single entity with full parsed body
```

**Example:**
```
GET /api/entities/machines/vps
→ {
    "_name": "vps",
    "_type": "machines",
    "title": "VPS — onemind-main",
    "subtype": "vps",
    "tags": ["type/machine", "domain/SO", "status/active"],
    "_body": "... full markdown body ..."
  }
```

### Stub endpoints (501 Not Implemented — placeholders for future)

```
POST /api/entities/:type/:name/update     ← write live data back (e.g. from n8n)
POST /api/entities/:type/:name/location   ← GPS update from Bouncie
GET  /api/entities/presence               ← who/what is home (Home Assistant)
```

These exist so the contract is defined. Wire them when the integrations are ready.

### To deploy when you're ready

```bash
# Copy api/ to VPS
scp -r "api/" root@104.248.4.212:/home/onemind/entity-registry/

# On VPS:
cd /home/onemind/entity-registry
npm install
# Set VAULT_PATH env var to point at the codex vault path
PORT=18810 VAULT_PATH=/home/onemind/onemind-os/codex/"00-24 SO (Sovereign Operations)/05 Environments (Entities & Registry)" node server.js

# Or wrap in systemd / PM2
# Port 18810 — no conflicts with existing services
```

**Note:** `server.js` currently hardcodes the registry root as one level above the `api/` folder. When deploying to VPS (where the api/ folder is separate from the vault), pass the vault path as an environment variable — one line change in `server.js`.

---

## What Still Needs Filling In

Several entity files have `fill in` placeholders — search for `fill in` in this folder:

| File | Missing |
|------|---------|
| `users/wife.md` | Name, Telegram ID, language, timezone |
| `users/kids.md` | Names, ages, Telegram IDs |
| `users/dog.md` | Name, breed, age, vet info |
| `machines/mac.md` | Hostname, Tailscale IP, RAM, storage, macOS version |
| `machines/vps.md` | RAM, CPU, storage, DigitalOcean plan + cost |
| `locations/main-house.md` | Address, WiFi SSID, HA URL, smart lock details |
| `locations/barn.md` | WiFi status, vehicles stored, HA area |
| `locations/office.md` | Address, WiFi SSID, shared/private |
| `devices/` | iPhone file not yet created — use `templates/device.md` |
| `digital-assets/` | Empty — domains, subscriptions, licenses not yet added |

---

## Related Files

| File | What |
|------|------|
| `SO/02 Intelligence/openclaw/MULTI-USER.md` | Full routing architecture — how Legacy identifies who's talking |
| `SO/02 Intelligence/openclaw/interfaces/telegram.md` | Telegram interface config |
| `SO/03 Agents/FIND-ZEUS-MODE.md` | Zeus location tracking (Bouncie, HA, phone) |
| `.claude/skills/bouncie/SKILL.md` | Bouncie GPS API skill |
| `.claude/skills/home-assistant/SKILL.md` | Home Assistant control skill |
