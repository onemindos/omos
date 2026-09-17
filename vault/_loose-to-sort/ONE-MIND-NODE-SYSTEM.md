---
id: doctrine-node-system-20260503
title: "ONE MIND NODE SYSTEM"
type: doctrine
status: active
created: 2026-05-03
modified: 2026-05-03
owner: zeus
pillar: SO
domain: "08 OM Assets (Entity Registry)"
ai_summary: "The One Mind Node System — every real-world entity agents need to know about, classified into 5 node categories (biological, digital, ecological, mechanical, physical) with access levels and a markdown-first registry architecture."
tags: [type/doctrine, zone/so, domain/so, status/active, framework/core]
related: ["[[ONE-MIND-ENTITY-SYSTEM]]", "[[ONE-MIND-FRAMEWORK]]", "[[ONEMIND-CODEX]]"]
depth: 1
---

# ONE MIND NODE SYSTEM

> **"Every noun in your world is a node. Nodes give agents context. Context makes AI useful."**
>
> By Zeus Delacruz | Created: May 2026

---

## What Is a Node?

A **node** is any real-world entity that agents need to know about and interact with.

If it exists in Zeus's world — a person, a machine, a building, a robot, a piece of land — it gets a node file. That file is the single source of truth for everything an agent needs to know about that entity.

Before the node system existed, agents had zero structured awareness of the world:
- "Remind my wife about dinner" → agent doesn't know who the wife is
- "What's running on the VPS?" → agent doesn't know the VPS IP or specs
- "Check if anyone's home" → agent can't query presence

**With the node system**, every entity has a markdown file with structured frontmatter. Any agent loads it and instantly knows who, what, where, and how.

---

## The 5 Node Categories

One Mind classifies every real-world entity into one of 5 categories:

```
┌─────────────────────────────────────────────────────────┐
│                  ONE MIND NODE SYSTEM                    │
│           Every noun in your world, typed and filed      │
├───────────┬───────────┬────────────┬──────────┬─────────┤
│BIOLOGICAL │  DIGITAL  │ ECOLOGICAL │MECHANICAL│PHYSICAL │
│           │           │            │          │         │
│ Humans    │ Machines  │ Land       │ Devices  │Locations│
│ Pets      │ Agents    │ Gardens    │ Robots   │Buildings│
│ Animals   │ Software  │ Ecosystems │ Vehicles │ Places  │
└───────────┴───────────┴────────────┴──────────┴─────────┘
```

| Category | What It Is | Examples |
|----------|-----------|---------|
| **Biological** | Living beings with presence and relationships | Zeus, Cielo, Liam, Noah, Soluna, family members, Cerberus (dog) |
| **Digital** | Compute infrastructure, AI agents, software systems | VPS, Mac, OpenClaw agents, APIs, software platforms |
| **Ecological** | Natural environment, land, ecosystems | Property land, gardens, outdoor spaces |
| **Mechanical** | Devices, robots, vehicles, hardware | iPhone, drones, DIMOS robot, OM1 robot, vehicles |
| **Physical** | Real-world locations and structures | Main house, barn, office, any named place |

---

## Node Registry Location

All nodes live in SO 08:

```
00-24 SO (Sovereign Operations)/
└── 08 OM Assets (Entity Registry) (The Registry) (Hive Mind)/
    └── nodes (Node Registry)/
        ├── biological (Biological Nodes)/     ← people, pets, living beings
        ├── digital (Digital Nodes)/           ← machines, agents, software
        ├── ecological (Ecological Nodes)/     ← land, nature, environment
        ├── mechanical (Mechanical Nodes)/     ← devices, robots, vehicles
        └── physical (Physical Nodes)/         ← locations, buildings, places
```

**Full path:** `00-24 SO/08 OM Assets/nodes/`

---

## Access Levels

Every node has an access level that controls what agents can share with whom and how much context they load.

| Level | Code | Who | What Agents Can Do |
|-------|------|-----|--------------------|
| **Owner** | L4 | Zeus | Full access — all agents, all memory, all commands, all domains |
| **Partner** | L3 | Cielo (wife) | Home agents, family memory, LE/HP domains |
| **Child** | L2 | Liam, Noah, Soluna | Age-gated — homework, schedules, safe content only |
| **Observer** | L1 | Guests, external contacts | Query-only, no writes, no personal context |
| **Entity** | L0 | Cerberus (dog), vehicles, assets | No chat — registry data and telemetry only |

Access level is set in the node's frontmatter:
```yaml
access: L3
```

---

## Node Frontmatter Standard

Every node file uses this base frontmatter:

```yaml
---
id: [category]-[name]-[YYYYMMDD]
title: "[Human readable name]"
type: [human|machine|device|robot|location|pet|digital-asset|vehicle]
node_category: [biological|digital|ecological|mechanical|physical]
status: [active|inactive|archived]
created: YYYY-MM-DD
updated: YYYY-MM-DD
owner: zeus
access: [L0|L1|L2|L3|L4]
pillar: SO
domain: "08 OM Assets (Entity Registry)"
ai_summary: "[One sentence — who/what this is and why it matters to agents.]"
tags: [type/[node-type], zone/so, domain/so, status/active]
related: []
depth: 3
---
```

### Category-specific extensions

**Biological (human):**
```yaml
relationship: [owner|partner|child|contact|guest]
telegram_id: "@handle"
language: en
timezone: America/Los_Angeles
birthday: YYYY-MM-DD
```

**Digital (machine):**
```yaml
subtype: [vps|mac|server|nas]
ip_public: "x.x.x.x"
ip_tailscale: "100.x.x.x"
os: ubuntu-24.04
services: [openclaw, postgres, ollama]
```

**Mechanical (device/robot):**
```yaml
subtype: [phone|tablet|robot|vehicle|drone]
model: "Device model name"
serial: "serial number"
os: iOS-18
```

**Physical (location):**
```yaml
address: "street address"
wifi_ssid: "network name"
home_assistant_area: "area_name"
type: [home|office|outbuilding|external]
```

---

## Current Node Registry

### Biological Nodes — People & Pets

| Node | Access | Status | Notes |
|------|--------|--------|-------|
| Zeus | L4 | Active | Owner — full profile |
| Cielo | L3 | Active | Partner — full profile |
| Liam | L2 | Active | Child — profile |
| Noah | L2 | Active | Child — profile |
| Soluna | L2 | Active | Child — profile |
| Wife | L3 | Template | Registry card — fill in details |
| Kids | L2 | Template | Registry card — fill in details |
| Cerberus (dog) | L0 | Template | Registry only — no chat |

### Digital Nodes — Machines & Infrastructure

| Node | Type | Status | Notes |
|------|------|--------|-------|
| VPS (onemindos-cloud) | Machine | Active | DigitalOcean, 104.248.4.212, Tailscale 100.77.52.93 |
| Mac (Zeus's MacBook) | Machine | Active | Source of truth, Claude Code primary |

### Mechanical Nodes — Devices & Robots

| Node | Type | Status | Notes |
|------|------|--------|-------|
| DIMOS | Robot | Registered | Robotic platform |
| OM1 | Robot | Registered | Robotic platform |
| Devices Registry | Device | Empty | iPhone not yet added |

### Physical Nodes — Locations

| Node | Type | Status | Notes |
|------|------|--------|-------|
| Main House | Home | Active | Primary residence, HA-integrated |
| Barn | Outbuilding | Active | Vehicle storage |
| Office | Work | Active | Work environment |

---

## How Agents Use Nodes

### 1. Context Loading

When a user message arrives, the agent identifies the sender and loads their node file:

```
Incoming message → identify sender → load biological/users/{name}/profile.md
                 → inject context → respond with full awareness
```

### 2. Multi-User Routing

Legacy routes responses based on who's talking. Each person's node file contains their Telegram ID, language, timezone, and access level — the agent loads the right context automatically.

```
Telegram message from Cielo's number
  → match telegram_id to cielo/profile.md
  → load L3 context (LE/HP domains)
  → respond as family assistant, not owner
```

### 3. Infrastructure Queries

```
"What's running on the VPS?"
  → load digital/machines-registry/vps.md
  → frontmatter has: IP, Tailscale address, services list
  → answer accurately from structured data
```

### 4. Location Awareness

```
"Who's home?"
  → query Home Assistant webhook → POST /api/entities/presence
  → cross-reference physical/locations-registry/main-house.md
  → return presence data
```

---

## Two-Layer Architecture

The node registry uses the same 2-layer architecture as the rest of One Mind:

### Layer 0 — Markdown (always active)

Plain `.md` files with YAML frontmatter. Every agent using OpenClaw's memory search already indexes and can query these files today with zero extra setup. Editing a file IS editing the registry — no sync, no migration, no database.

### Layer 1 — REST API (scaffold, not yet deployed)

A thin Express server in `08 OM Assets/api/` reads the same markdown files and serves them as structured JSON. Enables:
- Cross-entity queries ("all machines on Tailscale")
- Live data write-back from Bouncie GPS, Home Assistant
- Future: visual node graph in onemind-web

```
GET /api/entities/machines/vps    → full JSON with all frontmatter
GET /api/entities/presence        → who/what is currently home
POST /api/entities/assets/car-1/location  ← GPS update from Bouncie
```

**The markdown is always the source of truth. The API is a read layer on top.**

---

## Connection to the Entity System

The node system and the entity system work together:

| System | What It Classifies | Purpose |
|--------|--------------------|---------|
| **Entity System** | Content and work (Goal, Project, Task, Metric, Skill, Tool, Human, Agent, Calendar, Knowledge, Asset) | Routes everything you capture into the right pillar/domain |
| **Node System** | Real-world objects (biological, digital, ecological, mechanical, physical) | Gives agents structured awareness of the world Zeus operates in |

**The Entity System answers: "What type of content is this?"**
**The Node System answers: "What real-world thing is this?"**

A person like Cielo is both:
- An **Entity** of type `human` in the Entity System (routes to LE/51 Marriage domain)
- A **Node** of category `biological` in the Node System (profile.md for agent context loading)

The distinction: entity files track what you're doing WITH people (goals, projects, relationships). Node files track WHO people are (so agents can respond to them correctly).

---

## What Still Needs to Be Done

| Item | Status | Priority |
|------|--------|---------|
| Fill in wife.md (name, Telegram ID, language, timezone) | Empty | High |
| Fill in kids.md (names, ages, Telegram IDs) | Empty | High |
| Fill in dog.md (name, breed, age, vet) | Empty | Medium |
| Fill in machines/mac.md (hostname, Tailscale IP, specs, macOS) | Partial | High |
| Create devices/iphone.md | Missing | Medium |
| Populate digital-assets/ (domains, licenses, subscriptions) | Empty | Medium |
| Deploy `api/server.js` to VPS on port 18810 | Scaffold ready | Low |

---

## Full Registry Path Reference

```
00-24 SO/08 OM Assets/nodes/
├── biological/
│   └── users/
│       ├── zeus (Zeus)/profile.md            ← L4 owner
│       ├── cielo (Cielo)/profile.md          ← L3 partner
│       ├── liam (Liam)/profile.md            ← L2 child
│       ├── noah (Noah)/profile.md            ← L2 child
│       ├── soluna (Soluna)/profile.md        ← L2 child
│       └── family (Family Members)/
│           ├── wife.md                       ← L3 (fill in)
│           ├── kids.md                       ← L2 (fill in)
│           └── dog.md                        ← L0 (fill in)
├── digital/
│   └── machines-registry/
│       ├── vps.md                            ← DigitalOcean VPS
│       └── mac.md                            ← Zeus's MacBook
├── ecological/
│   └── (populate as needed)
├── mechanical/
│   ├── devices-registry/                     ← iphone.md not yet created
│   └── robots-registry/
│       ├── dimos (Dimos Robot)/
│       └── om1 (OM1 Robot)/
└── physical/
    └── locations-registry/
        ├── main-house.md
        ├── barn.md
        └── office.md
```

---

## Related Docs

| Document | What |
|----------|------|
| `ONE-MIND-ENTITY-SYSTEM.md` | The 11 entity types — what content IS |
| `ONE-MIND-FRAMEWORK.md` | The 4 pillars and CODE× method |
| `00-24 SO/08 OM Assets/_README.md` | Full technical registry docs, API scaffold, templates |
| `00-24 SO/01 Intelligence/openclaw/MULTI-USER.md` | How Legacy routes responses by user identity |

---

*One Mind Node System — created by Zeus Delacruz. May 2026.*
*Part of the One Mind ecosystem: [github.com/OneMind-OS/One-Mind](https://github.com/OneMind-OS/One-Mind)*
