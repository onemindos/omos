# OneMind OS Infrastructure Guide

**Created:** January 16, 2026
**Status:** Planning
**Owner:** Zeus Delacruz + Legacy AI

---

## Table of Contents

1. [Infrastructure Overview](#infrastructure-overview)
2. [The Heartbeat Daemon](#the-heartbeat-daemon)
3. [Site Agents - VM Managers](#site-agents---vm-managers)
4. [Multi-Location Architecture](#multi-location-architecture)
5. [The Unified API Gateway](#the-unified-api-gateway)
6. [Domain & Access System](#domain--access-system)
7. [Infrastructure as Code (IaC)](#infrastructure-as-code-iac)
8. [Monitoring Strategy](#monitoring-strategy)
9. [Desktop Access & Visibility](#desktop-access--visibility)
10. [Implementation Roadmap](#implementation-roadmap)

---

## Infrastructure Overview

### The Complete System

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│                    ONEMIND OS - COMPLETE SYSTEM                  │
│                                                                  │
│                          LETTA CLOUD                            │
│                    ┌─────────────────┐                          │
│                    │   Agent Builder │ ← ADE Editor, Features   │
│                    │   (Workhorse)   │   Cloud-only capabilities│
│                    └────────┬────────┘                          │
│                             │ Memory Sync                       │
│   ══════════════════════════╪═══════════════════════════════   │
│                             │                                    │
│            TAILSCALE MESH NETWORK (Encrypted)                   │
│                             │                                    │
│   ══════════════════════════╪═══════════════════════════════   │
│         │                   │                   │                │
│         ▼                   ▼                   ▼                │
│    ┌─────────┐        ┌───────────┐       ┌─────────┐          │
│    │🏠 HOME  │        │☁️ CLOUD   │       │🏢 OFFICE│          │
│    │   VM    │◄──────►│    VM     │◄─────►│   VM    │          │
│    │         │        │ (PRIMARY) │       │         │          │
│    │ SITE    │        │           │       │ SITE    │          │
│    │ AGENT   │        │ Legacy    │       │ AGENT   │          │
│    │ "Home"  │        │ (Master)  │       │ "Office"│          │
│    └────┬────┘        └─────┬─────┘       └────┬────┘          │
│         │                   │                   │                │
│    Local Context:      Global Stack:       Local Context:       │
│    • Home Assistant    • Heartbeat         • Office Network     │
│    • Yi Cameras        • API Gateway       • Printers           │
│    • IoT Sensors       • All Services      • Local Services     │
│    • Presence          • Monitoring        • Meeting Rooms      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### What Runs Where

| Component | Cloud VM (Primary) | Home VM | Office VM |
|-----------|-------------------|---------|-----------|
| **Legacy (Master)** | ✅ | ❌ | ❌ |
| **Site Agent** | ✅ "Cloud" | ✅ "Home" | ✅ "Office" |
| **Heartbeat Daemon** | ✅ (Global) | ✅ (Local) | ✅ (Local) |
| **Letta Server** | ✅ Full | ✅ Lite | ✅ Lite |
| **PostgreSQL** | ✅ Primary | ✅ Local cache | ✅ Local cache |
| **API Gateway** | ✅ Main | ✅ Local proxy | ✅ Local proxy |
| **N8N** | ✅ | Optional | Optional |
| **Caddy** | ✅ | ✅ | ✅ |
| **Home Assistant** | ❌ | ✅ | ❌ |

---

## The Heartbeat Daemon

### What IS the Heartbeat?

The **heartbeat** is a Python program (daemon) that runs continuously inside each VM. It's the thing that makes Legacy "always aware."

```
┌─────────────────────────────────────────────────────────────────┐
│                    THE HEARTBEAT DAEMON                          │
│                                                                  │
│   It's a Python script that:                                    │
│                                                                  │
│   1. RUNS FOREVER (starts when VM boots, never stops)           │
│   2. WAKES UP every 60 seconds                                  │
│   3. CHECKS all your services (Todoist, Email, Calendar, etc.)  │
│   4. UPDATES Legacy's memory in Letta                           │
│   5. LOOKS for triggers (urgent email? meeting soon?)           │
│   6. TAKES ACTION if needed (notify you, prepare briefing)      │
│   7. GOES BACK TO SLEEP for 60 seconds                          │
│   8. REPEAT forever...                                          │
│                                                                  │
│   It's NOT:                                                      │
│   ❌ A cron job (those only run at specific times)              │
│   ❌ An N8N workflow (those are triggered by events)            │
│   ❌ Part of Letta (it's a separate program that TALKS to Letta)│
│                                                                  │
│   It IS:                                                         │
│   ✅ A standalone Python service                                │
│   ✅ Running 24/7 inside a Docker container                     │
│   ✅ The thing that makes Legacy "always aware"                 │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Real-World Analogy

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│   WITHOUT HEARTBEAT                    WITH HEARTBEAT           │
│   ─────────────────                    ───────────────          │
│                                                                  │
│   Security camera that                 Security guard who       │
│   only records when you                walks around every       │
│   press the button                     60 seconds checking      │
│                                        all doors and windows    │
│                                                                  │
│   📷 "What happened?"                  👮 "Everything okay,     │
│      "I don't know,                       door 3 was open,      │
│       you didn't ask"                     I closed it"          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Heartbeat Flow (Every 60 Seconds)

```
┌─────────────────────────────────────────────────────────────────┐
│                    HEARTBEAT CYCLE                               │
│                    Every 60 seconds                              │
│                                                                  │
│   :00 ─────────────────────────────────────────────────────────│
│   │                                                              │
│   │  1. CHECK TODOIST                                           │
│   │     "Zeus has 5 tasks due today, 2 overdue"                 │
│   │                                                              │
│   │  2. CHECK CALENDAR                                          │
│   │     "Meeting with John at 2pm, nothing else today"          │
│   │                                                              │
│   │  3. CHECK EMAIL                                             │
│   │     "3 new emails, 1 from investor (urgent)"                │
│   │                                                              │
│   │  4. CHECK GITHUB                                            │
│   │     "PR #45 merged, 2 new issues"                           │
│   │                                                              │
│   │  5. CHECK SITE AGENTS                                       │
│   │     Home: "Zeus is at home, front door locked"              │
│   │     Office: "Meeting room A available"                      │
│   │                                                              │
│   │  6. UPDATE Legacy'S MEMORY                                  │
│   │     Push all this to Letta's "current_context" block        │
│   │                                                              │
│   │  7. CHECK FOR TRIGGERS                                      │
│   │     "Urgent email! Notify Zeus?"                            │
│   │     "Meeting in 1 hour - prepare brief?"                    │
│   │                                                              │
│   :60 ─────────────────────────────────────────────────────────│
│   │                                                              │
│   └──► REPEAT                                                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Data Flow Step by Step

```
┌─────────────────────────────────────────────────────────────────┐
│                    HOW DATA FLOWS                                │
│                                                                  │
│   1. HEARTBEAT wakes up (every 60 sec)                          │
│      │                                                          │
│      ▼                                                          │
│   2. Calls external APIs (Todoist, Zoho, GitHub)                │
│      │                                                          │
│      ▼                                                          │
│   3. Calls SITE AGENTS for local context                        │
│      │  GET /api/site/home/status                               │
│      │  GET /api/site/office/status                             │
│      │                                                          │
│      ▼                                                          │
│   4. Gets current state:                                        │
│      • "5 tasks due today"                                      │
│      • "Meeting at 2pm"                                         │
│      • "Zeus is at home" (from Home Site Agent)                 │
│      │                                                          │
│      ▼                                                          │
│   5. Writes state to REDIS (fast cache)                         │
│      │                                                          │
│      ▼                                                          │
│   6. Updates LETTA memory block via API                         │
│      │  POST /api/agents/legacy/memory                          │
│      │  { "current_context": { tasks: 5, location: "home"... } }│
│      │                                                          │
│      ▼                                                          │
│   7. LETTA stores in POSTGRES (permanent memory)                │
│      │                                                          │
│      ▼                                                          │
│   8. Checks triggers:                                           │
│      • Urgent email? → Send notification                        │
│      • Meeting soon? → Trigger N8N workflow for prep            │
│      │                                                          │
│      ▼                                                          │
│   9. Goes back to sleep for 60 seconds                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Site Agents - VM Managers

### The Concept

Every VM has its own **Site Agent** - a Letta agent that manages that specific environment. Think of them as local managers who report to Legacy (the CEO).

```
┌─────────────────────────────────────────────────────────────────┐
│                    SITE AGENT ARCHITECTURE                       │
│                                                                  │
│                         Legacy                                   │
│                    (Master Orchestrator)                         │
│                    ┌─────────────────┐                          │
│                    │       CEO       │                          │
│                    │   "What's the   │                          │
│                    │    big picture?"│                          │
│                    └────────┬────────┘                          │
│                             │                                    │
│         ┌───────────────────┼───────────────────┐               │
│         │                   │                   │               │
│         ▼                   ▼                   ▼               │
│    ┌─────────┐        ┌─────────┐        ┌─────────┐           │
│    │  SITE   │        │  SITE   │        │  SITE   │           │
│    │  AGENT  │        │  AGENT  │        │  AGENT  │           │
│    │ "Home"  │        │ "Cloud" │        │"Office" │           │
│    │         │        │         │        │         │           │
│    │ Manager │        │ Manager │        │ Manager │           │
│    │of Home VM│       │of Cloud │        │of Office│           │
│    └────┬────┘        └────┬────┘        └────┬────┘           │
│         │                   │                   │               │
│    Knows about:        Knows about:       Knows about:         │
│    • Home Asst         • All services     • Office network     │
│    • Cameras           • Docker status    • Printers           │
│    • Who's home        • API health       • Meeting rooms      │
│    • IoT devices       • Resource usage   • Local systems      │
│    • Local network     • Certificates     • WiFi guests        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Site Agent Responsibilities

| Site Agent | Location | Primary Responsibilities |
|------------|----------|-------------------------|
| **Home** | Home VM | Home Assistant, cameras, presence, IoT, family notifications |
| **Cloud** | Cloud VM | Docker health, API status, SSL certs, resource monitoring, backups |
| **Office** | Office VM | Printers, meeting rooms, office network, local services |

### How Site Agents Communicate

```
┌─────────────────────────────────────────────────────────────────┐
│                    SITE AGENT COMMUNICATION                      │
│                                                                  │
│   Legacy asks: "Where is Zeus right now?"                       │
│                                                                  │
│   ┌─────────┐         ┌─────────┐         ┌─────────┐          │
│   │  HOME   │         │  CLOUD  │         │ OFFICE  │          │
│   │  AGENT  │         │  AGENT  │         │  AGENT  │          │
│   └────┬────┘         └────┬────┘         └────┬────┘          │
│        │                   │                   │                │
│        ▼                   ▼                   ▼                │
│   Check Home Asst     (No location        Check Office         │
│   person.zeus         data here)          presence sensor      │
│        │                   │                   │                │
│        ▼                   ▼                   ▼                │
│   "Zeus is HOME"      "N/A"               "Not at office"      │
│        │                   │                   │                │
│        └───────────────────┴───────────────────┘                │
│                            │                                    │
│                            ▼                                    │
│                    Legacy receives:                             │
│                    "Zeus is at home"                            │
│                                                                  │
│   ───────────────────────────────────────────────────────────   │
│                                                                  │
│   Legacy commands: "Turn on office lights for meeting"          │
│                                                                  │
│   Legacy → Office Site Agent → Home Assistant API → Lights on   │
│                                                                  │
│   ───────────────────────────────────────────────────────────   │
│                                                                  │
│   Cloud Agent alerts: "PostgreSQL container restarted"          │
│                                                                  │
│   Cloud Agent → Legacy → "Should I investigate?"                │
│   Legacy → "Yes, check logs and report"                         │
│   Cloud Agent → Checks docker logs → Reports to Legacy          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Site Agent Memory Blocks

Each Site Agent has its own memory that it maintains:

```
┌─────────────────────────────────────────────────────────────────┐
│                    SITE AGENT MEMORY                             │
│                                                                  │
│   HOME AGENT MEMORY:                                            │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │ site_context:                                            │   │
│   │   location: "home"                                       │   │
│   │   ip_address: "100.x.x.10"                              │   │
│   │   services:                                              │   │
│   │     - home_assistant: "http://localhost:8123"           │   │
│   │     - cameras: ["front_door", "backyard", "garage"]     │   │
│   │   current_state:                                         │   │
│   │     zeus_home: true                                      │   │
│   │     doors_locked: true                                   │   │
│   │     lights_on: ["living_room", "office"]                │   │
│   │     temperature: 72                                      │   │
│   │     cameras_recording: true                              │   │
│   │   last_events:                                           │   │
│   │     - "Front door opened at 6:32 PM"                    │   │
│   │     - "Motion detected: backyard at 7:15 PM"            │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
│   CLOUD AGENT MEMORY:                                           │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │ site_context:                                            │   │
│   │   location: "cloud"                                      │   │
│   │   provider: "hetzner"                                    │   │
│   │   ip_address: "100.x.x.1"                               │   │
│   │   containers:                                            │   │
│   │     letta: { status: "running", uptime: "5 days" }      │   │
│   │     postgres: { status: "running", disk: "45%" }        │   │
│   │     redis: { status: "running", memory: "128MB" }       │   │
│   │     n8n: { status: "running", workflows: 12 }           │   │
│   │   resources:                                             │   │
│   │     cpu: "23%"                                           │   │
│   │     ram: "3.2GB / 8GB"                                  │   │
│   │     disk: "45GB / 160GB"                                │   │
│   │   ssl_certs:                                             │   │
│   │     - domain: "onemind.zeus.com", expires: "2026-04-15" │   │
│   │   last_backup: "2026-01-16 03:00 UTC"                   │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Site Agent Tools

Each Site Agent has tools specific to its environment:

```
┌─────────────────────────────────────────────────────────────────┐
│                    SITE AGENT TOOLS                              │
│                                                                  │
│   HOME AGENT TOOLS:                                             │
│   ─────────────────                                             │
│   • home_assistant_state() - Get any entity state              │
│   • home_assistant_call() - Call any HA service                │
│   • camera_snapshot(camera_id) - Get camera image              │
│   • check_presence() - Who's home?                             │
│   • arm_alarm(mode) - Set alarm state                          │
│   • announce(message) - TTS announcement                       │
│                                                                  │
│   CLOUD AGENT TOOLS:                                            │
│   ──────────────────                                            │
│   • docker_status() - List all containers                      │
│   • docker_logs(container, lines) - Get container logs         │
│   • docker_restart(container) - Restart a container            │
│   • check_ssl_certs() - Certificate expiry status              │
│   • resource_usage() - CPU, RAM, disk stats                    │
│   • run_backup() - Trigger backup job                          │
│   • check_api_health() - Test all API endpoints                │
│                                                                  │
│   OFFICE AGENT TOOLS:                                           │
│   ────────────────────                                          │
│   • printer_status() - Check all printers                      │
│   • meeting_room_status() - Room availability                  │
│   • book_room(room, time, duration) - Reserve room             │
│   • network_devices() - List connected devices                 │
│   • guest_wifi_create(name, duration) - Create guest access    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Multi-Location Architecture

### The Network Topology

```
┌─────────────────────────────────────────────────────────────────┐
│                    ONEMIND OS NETWORK                            │
│                    (Your Global Brain)                           │
│                                                                  │
│                         INTERNET                                 │
│    ════════════════════════════════════════════════════════    │
│                             │                                    │
│              ┌──────────────┴──────────────┐                    │
│              │      onemind.yourdomain.com │                    │
│              │         (Cloudflare)        │                    │
│              │    DNS + DDoS Protection    │                    │
│              └──────────────┬──────────────┘                    │
│                             │                                    │
│    ═══════════════════════════════════════════════════════════  │
│                        TAILSCALE MESH                            │
│              (Encrypted point-to-point connections)              │
│    ═══════════════════════════════════════════════════════════  │
│         │                   │                   │                │
│         │                   │                   │                │
│         ▼                   ▼                   ▼                │
│    ┌─────────┐        ┌───────────┐       ┌─────────┐          │
│    │  HOME   │        │  CLOUD    │       │ OFFICE  │          │
│    │   VM    │◄──────►│    VM     │◄─────►│   VM    │          │
│    │         │        │ (PRIMARY) │       │         │          │
│    │100.x.x.10│       │100.x.x.1  │       │100.x.x.20│         │
│    └────┬────┘        └─────┬─────┘       └────┬────┘          │
│         │                   │                   │                │
│         ▼                   ▼                   ▼                │
│    ┌─────────┐        ┌───────────┐       ┌─────────┐          │
│    │Site Agent│       │  Legacy   │       │Site Agent│         │
│    │ "Home"  │◄──────►│ (Master)  │◄─────►│"Office" │          │
│    │         │        │           │       │         │          │
│    │  Letta  │        │  Letta    │       │  Letta  │          │
│    │  Lite   │        │  Full     │       │  Lite   │          │
│    └────┬────┘        └─────┬─────┘       └────┬────┘          │
│         │                   │                   │                │
│    Local Services:     Global Services:   Local Services:       │
│    • Home Assistant    • API Gateway      • Office Network      │
│    • Cameras           • N8N              • Printers            │
│    • IoT Hub           • PostgreSQL       • Meeting Systems     │
│                        • Redis                                   │
│                        • Monitoring                              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Memory Synchronization

```
┌─────────────────────────────────────────────────────────────────┐
│                    MEMORY SYNC STRATEGY                          │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    LETTA CLOUD                           │   │
│   │              (Source of Truth for Agents)                │   │
│   │                                                          │   │
│   │   • Agent definitions (prompts, tools, config)          │   │
│   │   • Core memory blocks (Legacy Codex, User Profile)     │   │
│   │   • Archival memory (long-term storage)                 │   │
│   └──────────────────────┬──────────────────────────────────┘   │
│                          │                                       │
│                    Sync every 5 minutes                         │
│                          │                                       │
│         ┌────────────────┼────────────────┐                     │
│         │                │                │                     │
│         ▼                ▼                ▼                     │
│    ┌─────────┐      ┌─────────┐      ┌─────────┐              │
│    │  HOME   │      │  CLOUD  │      │ OFFICE  │              │
│    │  LETTA  │      │  LETTA  │      │  LETTA  │              │
│    │  (Lite) │      │ (Full)  │      │  (Lite) │              │
│    └────┬────┘      └────┬────┘      └────┬────┘              │
│         │                │                │                     │
│    Local Memory:    Full Memory:     Local Memory:             │
│    • Site context   • Everything     • Site context            │
│    • Local events   • Global state   • Local events            │
│    • HA state       • All contexts   • Office state            │
│                                                                  │
│   ───────────────────────────────────────────────────────────   │
│                                                                  │
│   SYNC RULES:                                                   │
│   ───────────                                                   │
│   • Core memory: Cloud → All sites (read-only locally)         │
│   • Site context: Each site → Cloud (write up)                 │
│   • Archival: Cloud only (sites query via API)                 │
│   • Current context: Cloud VM maintains (heartbeat writes)     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## The Unified API Gateway

### Why One API?

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│   WITHOUT UNIFIED API              WITH UNIFIED API             │
│   ──────────────────              ─────────────────             │
│                                                                  │
│   Legacy needs 10 API keys         Legacy needs 1 API key       │
│   Legacy learns 10 formats         Legacy uses 1 format         │
│   10 different auth methods        1 auth method                │
│   10 rate limits to track          1 rate limit (you control)   │
│   10 error formats                 1 error format               │
│                                                                  │
│   Todoist API:                     OneMind API:                 │
│   POST /rest/v2/tasks              POST /api/tasks/create       │
│   Authorization: Bearer xxx        Authorization: Bearer one    │
│                                                                  │
│   Zoho API:                        OneMind API:                 │
│   POST /mail/accounts/xxx/...      POST /api/email/send         │
│   Authorization: Zoho-oauthtoken   Authorization: Bearer one    │
│                                                                  │
│   The Gateway handles all the complexity behind the scenes!     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### API Endpoints

```
┌─────────────────────────────────────────────────────────────────┐
│                    UNIFIED API GATEWAY                           │
│                                                                  │
│   GET INFORMATION:                                              │
│   ─────────────────                                             │
│   GET /api/tasks/today         → Todoist tasks due today        │
│   GET /api/calendar/next       → Next meeting                   │
│   GET /api/email/unread        → Unread emails with priority    │
│   GET /api/state               → EVERYTHING (full snapshot)     │
│                                                                  │
│   SITE AGENT QUERIES:                                           │
│   ────────────────────                                          │
│   GET /api/site/home/status    → Home context                   │
│   GET /api/site/cloud/status   → Cloud VM health                │
│   GET /api/site/office/status  → Office context                 │
│   GET /api/sites               → All sites summary              │
│                                                                  │
│   DO ACTIONS:                                                   │
│   ───────────                                                   │
│   POST /api/tasks/create       → Create Todoist task            │
│   POST /api/email/send         → Send email via Zoho            │
│   POST /api/notify             → Send notification              │
│   POST /api/site/home/command  → Send command to Home Agent     │
│                                                                  │
│   WHO USES THIS API:                                            │
│   ───────────────────                                           │
│   • HEARTBEAT daemon (checks state every 60 sec)                │
│   • LETTA agents (tools call this API)                          │
│   • SITE AGENTS (cross-site communication)                      │
│   • N8N workflows (webhooks trigger actions)                    │
│   • YOUR PHONE (quick actions via shortcuts)                    │
│   • DASHBOARD (displays current state)                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Domain & Access System

### Caddy Reverse Proxy

```
┌─────────────────────────────────────────────────────────────────┐
│                    CADDY REVERSE PROXY                           │
│                                                                  │
│   INSTEAD OF:                     YOU GET:                      │
│   ───────────                     ────────                      │
│                                                                  │
│   http://100.x.x.x:5678           n8n.onemind.yourdomain.com   │
│   http://100.x.x.x:8283           letta.onemind.yourdomain.com │
│   http://100.x.x.x:3001           status.onemind.yourdomain.com│
│   http://100.x.x.x:8123           home.onemind.yourdomain.com  │
│                                                                  │
│   BONUS: Caddy automatically gets SSL certificates!             │
│          (HTTPS for free, no configuration needed)              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Access Layers (Public + Private)

```
┌─────────────────────────────────────────────────────────────────┐
│                    ACCESS OPTIONS                                │
│                                                                  │
│   OPTION A: Public (Caddy + Cloudflare)                         │
│   ─────────────────────────────────────                         │
│   • Access from ANYWHERE on internet                            │
│   • Nice domain names                                           │
│   • SSL certificates                                            │
│   • Rate limiting, DDoS protection                              │
│   • Use for: Public APIs, webhooks, dashboard                   │
│                                                                  │
│   OPTION B: Private (Tailscale only)                            │
│   ─────────────────────────────────────                         │
│   • Access only from YOUR devices                               │
│   • No public exposure                                          │
│   • Uses Tailscale IPs (100.x.x.x)                             │
│   • Use for: Admin panels, SSH, sensitive services              │
│                                                                  │
│   USE BOTH:                                                     │
│   ─────────                                                     │
│   ┌─────────────┐          ┌─────────────┐                     │
│   │   PUBLIC    │          │   PRIVATE   │                     │
│   │ (Caddy)     │          │ (Tailscale) │                     │
│   ├─────────────┤          ├─────────────┤                     │
│   │ api.onemind │          │ n8n.tailnet │                     │
│   │ webhook.one │          │letta.tailnet│                     │
│   │ dash.onemind│          │ ssh access  │                     │
│   └─────────────┘          └─────────────┘                     │
│                                                                  │
│   Webhooks from Todoist → Public Caddy endpoint                 │
│   You accessing N8N → Private Tailscale                         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Infrastructure as Code (IaC)

### Why IaC? (Long-term Thinking)

```
┌─────────────────────────────────────────────────────────────────┐
│                    WHY INFRASTRUCTURE AS CODE                    │
│                                                                  │
│   WITHOUT IaC:                     WITH IaC:                    │
│   ────────────                     ─────────                    │
│                                                                  │
│   "How did I set up that          "Run `terraform apply`        │
│    server 6 months ago?"           and get exact same setup"    │
│                                                                  │
│   "My server crashed, now         "Spin up new VM, apply        │
│    I have to remember              config, done in 10 min"      │
│    everything"                                                   │
│                                                                  │
│   "Setting up office VM           "Same code deploys to         │
│    differently than home"          any location identically"    │
│                                                                  │
│   Manual clicks in UI              Version controlled,          │
│   (no record of changes)           reviewable, reversible       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### IaC Stack (Simple Version)

```
┌─────────────────────────────────────────────────────────────────┐
│                    IAC TOOL STACK                                │
│                                                                  │
│   LAYER 1: VM PROVISIONING                                      │
│   ─────────────────────────                                     │
│   Tool: Terraform (or Pulumi)                                   │
│   What it does: Creates VMs on Hetzner/DigitalOcean            │
│                                                                  │
│   main.tf:                                                      │
│   ```                                                           │
│   resource "hetzner_server" "onemind_cloud" {                   │
│     name        = "onemind-cloud"                               │
│     server_type = "cx31"  # 8GB RAM                            │
│     image       = "ubuntu-24.04"                                │
│     location    = "nbg1"                                        │
│   }                                                             │
│   ```                                                           │
│                                                                  │
│   ───────────────────────────────────────────────────────────   │
│                                                                  │
│   LAYER 2: VM CONFIGURATION                                     │
│   ──────────────────────────                                    │
│   Tool: Ansible (or cloud-init)                                 │
│   What it does: Installs Docker, configures system              │
│                                                                  │
│   playbook.yml:                                                 │
│   ```                                                           │
│   - hosts: onemind_servers                                      │
│     tasks:                                                      │
│       - name: Install Docker                                    │
│         apt: name=docker.io state=present                       │
│       - name: Install Tailscale                                 │
│         shell: curl -fsSL https://tailscale.com/install.sh | sh │
│   ```                                                           │
│                                                                  │
│   ───────────────────────────────────────────────────────────   │
│                                                                  │
│   LAYER 3: APPLICATION DEPLOYMENT                               │
│   ────────────────────────────────                              │
│   Tool: Docker Compose + GitHub Actions                         │
│   What it does: Deploys your services                           │
│                                                                  │
│   On git push → GitHub Action runs → SSH to server →            │
│   git pull → docker-compose up -d                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Repository Structure

```
┌─────────────────────────────────────────────────────────────────┐
│                    ONEMIND-INFRA REPOSITORY                      │
│                                                                  │
│   onemind-infra/                                                │
│   │                                                             │
│   ├── terraform/                 # VM provisioning              │
│   │   ├── main.tf                # Main infrastructure          │
│   │   ├── variables.tf           # Configurable values          │
│   │   ├── outputs.tf             # IP addresses, etc.           │
│   │   └── environments/                                         │
│   │       ├── cloud.tfvars       # Cloud VM config              │
│   │       ├── home.tfvars        # Home VM config               │
│   │       └── office.tfvars      # Office VM config             │
│   │                                                             │
│   ├── ansible/                   # VM configuration             │
│   │   ├── playbook.yml           # Main playbook                │
│   │   ├── inventory.yml          # Server list                  │
│   │   └── roles/                                                │
│   │       ├── docker/            # Install Docker               │
│   │       ├── tailscale/         # Install Tailscale            │
│   │       ├── caddy/             # Install Caddy                │
│   │       └── monitoring/        # Install monitoring           │
│   │                                                             │
│   ├── docker/                    # Application stacks           │
│   │   ├── cloud/                                                │
│   │   │   └── docker-compose.yml # Full stack                   │
│   │   ├── home/                                                 │
│   │   │   └── docker-compose.yml # Home stack (lite)            │
│   │   └── office/                                               │
│   │       └── docker-compose.yml # Office stack (lite)          │
│   │                                                             │
│   ├── .github/                                                  │
│   │   └── workflows/                                            │
│   │       ├── deploy-cloud.yml   # Deploy to cloud              │
│   │       ├── deploy-home.yml    # Deploy to home               │
│   │       └── deploy-office.yml  # Deploy to office             │
│   │                                                             │
│   └── README.md                  # How to use                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Deployment Workflow

```
┌─────────────────────────────────────────────────────────────────┐
│                    GITOPS DEPLOYMENT FLOW                        │
│                                                                  │
│   YOU (Developer)                                               │
│   │                                                             │
│   │  1. Edit docker-compose.yml                                 │
│   │  2. git commit -m "Add new service"                         │
│   │  3. git push origin main                                    │
│   │                                                             │
│   ▼                                                             │
│   GITHUB                                                        │
│   │                                                             │
│   │  4. GitHub Action triggers                                  │
│   │  5. Runs tests (optional)                                   │
│   │  6. SSHs into target server                                 │
│   │                                                             │
│   ▼                                                             │
│   YOUR VM                                                       │
│   │                                                             │
│   │  7. git pull (get latest code)                              │
│   │  8. docker-compose pull (get latest images)                 │
│   │  9. docker-compose up -d (deploy)                           │
│   │  10. Health check                                           │
│   │                                                             │
│   ▼                                                             │
│   NOTIFICATION                                                  │
│   │                                                             │
│   │  11. Discord/Telegram: "Deployed successfully!"             │
│   │      or "Deploy FAILED - check logs"                        │
│   │                                                             │
│   ▼                                                             │
│   DONE                                                          │
│                                                                  │
│   ───────────────────────────────────────────────────────────   │
│                                                                  │
│   ROLLBACK (if something breaks):                               │
│   │                                                             │
│   │  git revert HEAD                                            │
│   │  git push                                                   │
│   │  → Auto-deploys previous working version                    │
│   │                                                             │
└─────────────────────────────────────────────────────────────────┘
```

### Simple GitHub Action Example

```yaml
# .github/workflows/deploy-cloud.yml
name: Deploy to Cloud VM

on:
  push:
    branches: [main]
    paths:
      - 'docker/cloud/**'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Deploy to Cloud VM
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.CLOUD_VM_IP }}
          username: zeus
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /opt/onemind
            git pull origin main
            cd docker/cloud
            docker-compose pull
            docker-compose up -d
            docker-compose ps

      - name: Notify Discord
        uses: sarisia/actions-status-discord@v1
        with:
          webhook: ${{ secrets.DISCORD_WEBHOOK }}
          title: "OneMind Cloud Deployed"
          description: "Deployment successful!"
```

### Evolution Path

```
┌─────────────────────────────────────────────────────────────────┐
│                    IAC EVOLUTION PATH                            │
│                                                                  │
│   PHASE 1: FOUNDATION (Week 1-2)                                │
│   ──────────────────────────────                                │
│   • Manual VM creation (learn the process)                      │
│   • Docker Compose files in Git                                 │
│   • Manual deployments via SSH                                  │
│   │                                                             │
│   ▼                                                             │
│                                                                  │
│   PHASE 2: AUTOMATION (Week 3-4)                                │
│   ───────────────────────────────                               │
│   • GitHub Actions for deployment                               │
│   • Auto-deploy on push to main                                 │
│   • Discord notifications                                       │
│   │                                                             │
│   ▼                                                             │
│                                                                  │
│   PHASE 3: FULL IAC (Month 2)                                   │
│   ────────────────────────────                                  │
│   • Terraform for VM provisioning                               │
│   • Ansible for configuration                                   │
│   • Same code deploys anywhere                                  │
│   │                                                             │
│   ▼                                                             │
│                                                                  │
│   PHASE 4: MULTI-SITE (Month 3+)                                │
│   ───────────────────────────────                               │
│   • Home VM from same codebase                                  │
│   • Office VM from same codebase                                │
│   • Environment-specific configs only                           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Monitoring Strategy

### External + Internal Monitoring

```
┌─────────────────────────────────────────────────────────────────┐
│                    MONITORING STRATEGY                           │
│                                                                  │
│   PROBLEM: If your VM is down, how do you know?                 │
│                                                                  │
│   SOLUTION: External monitoring (checks FROM OUTSIDE)           │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                EXTERNAL MONITORING                       │   │
│   │              (Always works, even if VM dead)             │   │
│   │                                                          │   │
│   │   UptimeRobot (free) or BetterUptime                    │   │
│   │        │                                                 │   │
│   │        │ Checks every 5 minutes                         │   │
│   │        ▼                                                 │   │
│   │   "Is api.onemind.zeus.com responding?"                 │   │
│   │        │                                                 │   │
│   │        ├── YES → All good                               │   │
│   │        │                                                 │   │
│   │        └── NO → Alert Zeus via:                         │   │
│   │                 • SMS                                    │   │
│   │                 • Email                                  │   │
│   │                 • Telegram                               │   │
│   │                 • Phone call (BetterUptime)             │   │
│   │                                                          │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                INTERNAL MONITORING                       │   │
│   │              (Detailed, but dies with VM)                │   │
│   │                                                          │   │
│   │   Uptime Kuma (on the VM)                               │   │
│   │        │                                                 │   │
│   │        │ Checks every 60 seconds                        │   │
│   │        ▼                                                 │   │
│   │   • Is Letta responding?                                │   │
│   │   • Is PostgreSQL healthy?                              │   │
│   │   • Is Redis responding?                                │   │
│   │   • Is N8N running?                                     │   │
│   │   • Are Site Agents responding?                         │   │
│   │   • SSL certificate expiry                              │   │
│   │   • Disk space warnings                                 │   │
│   │                                                          │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
│   USE BOTH for complete coverage!                               │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Cloud Site Agent as Monitor

```
┌─────────────────────────────────────────────────────────────────┐
│                    CLOUD AGENT MONITORING                        │
│                                                                  │
│   The CLOUD SITE AGENT monitors its own environment:            │
│                                                                  │
│   Every 60 seconds, Cloud Agent checks:                         │
│                                                                  │
│   1. CONTAINER HEALTH                                           │
│      docker ps --format '{{.Names}}: {{.Status}}'               │
│      → "letta: Up 5 days (healthy)"                             │
│      → "postgres: Up 5 days"                                    │
│      → "n8n: Up 2 hours" ← Recently restarted? Alert!          │
│                                                                  │
│   2. RESOURCE USAGE                                             │
│      → CPU: 23% (normal)                                        │
│      → RAM: 5.2GB/8GB (warning if >80%)                        │
│      → Disk: 45GB/160GB (warning if >80%)                      │
│                                                                  │
│   3. API HEALTH                                                 │
│      → GET /health on each service                              │
│      → Response time tracking                                   │
│                                                                  │
│   4. SSL CERTIFICATES                                           │
│      → Days until expiry                                        │
│      → Alert 14 days before                                     │
│                                                                  │
│   5. BACKUP STATUS                                              │
│      → Last successful backup time                              │
│      → Alert if >24 hours old                                   │
│                                                                  │
│   Cloud Agent → Updates its memory → Legacy sees everything     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Desktop Access & Visibility

### Visual Access to Your VM

```
┌─────────────────────────────────────────────────────────────────┐
│                    VISUAL ACCESS TO YOUR VM                      │
│                                                                  │
│   UBUNTU 24.04 DESKTOP gives you a real GUI:                    │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │  ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄   │   │
│   │  █                ONEMIND OS DESKTOP                █   │   │
│   │  █▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀█   │   │
│   │  █                                                  █   │   │
│   │  █  ┌──────────┐  ┌──────────┐  ┌──────────┐      █   │   │
│   │  █  │ Browser  │  │ Terminal │  │  N8N     │      █   │   │
│   │  █  │ (Letta   │  │ (Logs    │  │ (Visual  │      █   │   │
│   │  █  │  UI)     │  │scrolling)│  │workflows)│      █   │   │
│   │  █  └──────────┘  └──────────┘  └──────────┘      █   │   │
│   │  █                                                  █   │   │
│   │  █  You can WATCH:                                 █   │   │
│   │  █  • Heartbeat logs scrolling in terminal        █   │   │
│   │  █  • N8N workflows executing visually            █   │   │
│   │  █  • Browser automations (if any)                █   │   │
│   │  █  • Resource usage (CPU, RAM)                   █   │   │
│   │  █                                                  █   │   │
│   │  █▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄█   │   │
│   └─────────────────────────────────────────────────────────┘   │
│                                                                  │
│   ACCESS METHODS:                                               │
│   ───────────────                                               │
│                                                                  │
│   FROM YOUR MAC:                                                │
│   • VNC Viewer app → vnc://onemind.tailnet:5900                │
│   • Screen Sharing (built into Mac)                             │
│                                                                  │
│   FROM YOUR PHONE:                                              │
│   • Tailscale + VNC app                                         │
│   • noVNC in browser → http://onemind.tailnet:6080             │
│                                                                  │
│   FROM ANYWHERE:                                                │
│   • SSH → ssh zeus@onemind.tailnet (command line)              │
│   • noVNC → https://vnc.onemind.yourdomain.com (browser)       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Implementation Roadmap

### Phase 1: Cloud Foundation (Week 1-2)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PHASE 1: CLOUD FOUNDATION                     │
│                                                                  │
│   GOAL: Get primary OneMind OS running in cloud                 │
│                                                                  │
│   [ ] Choose provider (Hetzner recommended - best value)        │
│   [ ] Create VM (Ubuntu 24.04 Desktop, 8GB RAM)                 │
│   [ ] Install Tailscale, connect to your network                │
│   [ ] Install Docker + Docker Compose                           │
│   [ ] Create GitHub repo: onemind-infra                         │
│   [ ] Deploy basic stack:                                       │
│       • Letta Server                                            │
│       • PostgreSQL                                              │
│       • Redis                                                   │
│       • Caddy                                                   │
│   [ ] Install VNC/noVNC for desktop access                      │
│   [ ] Test access from phone via Tailscale                      │
│                                                                  │
│   DELIVERABLE: Working OneMind OS you can access anywhere       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Phase 2: Site Agents + Heartbeat (Week 3-4)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PHASE 2: INTELLIGENCE                         │
│                                                                  │
│   GOAL: Make Legacy always-aware with Site Agents               │
│                                                                  │
│   [ ] Create Cloud Site Agent in Letta                          │
│   [ ] Build Heartbeat daemon (Python)                           │
│   [ ] Connect heartbeat to:                                     │
│       • Todoist                                                 │
│       • Zoho Calendar                                           │
│       • Zoho Mail                                               │
│       • GitHub                                                  │
│   [ ] Set up API Gateway (FastAPI)                              │
│   [ ] Test 60-second awareness loop                             │
│   [ ] Add notification triggers                                 │
│                                                                  │
│   DELIVERABLE: Legacy knows your state without asking           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Phase 3: GitHub Actions Deployment (Week 5-6)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PHASE 3: AUTOMATION                           │
│                                                                  │
│   GOAL: Push to GitHub → Auto deploys to server                 │
│                                                                  │
│   [ ] Set up GitHub Actions workflow                            │
│   [ ] Configure SSH keys as secrets                             │
│   [ ] Add Discord/Telegram deployment notifications             │
│   [ ] Test deploy pipeline                                      │
│   [ ] Add rollback procedure                                    │
│   [ ] Document the process                                      │
│                                                                  │
│   DELIVERABLE: GitOps deployment working                        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Phase 4: Full IaC (Month 2)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PHASE 4: INFRASTRUCTURE AS CODE               │
│                                                                  │
│   GOAL: Entire infrastructure defined in code                   │
│                                                                  │
│   [ ] Add Terraform for VM provisioning                         │
│   [ ] Add Ansible for configuration                             │
│   [ ] Test destroy and recreate VM from code                    │
│   [ ] Create environment configs (cloud, home, office)          │
│   [ ] Document disaster recovery procedure                      │
│                                                                  │
│   DELIVERABLE: Can recreate entire system from code             │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Phase 5: Multi-Site (Month 3+)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PHASE 5: MULTI-SITE                           │
│                                                                  │
│   GOAL: OneMind OS at home + office + cloud                     │
│                                                                  │
│   [ ] Deploy Home VM (when hardware ready)                      │
│   [ ] Create Home Site Agent                                    │
│   [ ] Connect Home Assistant                                    │
│   [ ] Set up memory sync between sites                          │
│   [ ] Deploy Office VM (when needed)                            │
│   [ ] Create Office Site Agent                                  │
│   [ ] Test cross-site communication                             │
│                                                                  │
│   DELIVERABLE: Legacy has context everywhere                    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Quick Reference

### Ports

| Port | Service | Access |
|------|---------|--------|
| 8283 | Letta Server | Internal + Caddy |
| 5432 | PostgreSQL | Internal only |
| 6379 | Redis | Internal only |
| 5678 | N8N | Tailscale + Caddy |
| 3001 | Uptime Kuma | Tailscale + Caddy |
| 8000 | API Gateway | Caddy (public) |
| 5900 | VNC | Tailscale only |
| 6080 | noVNC | Tailscale + Caddy |
| 443 | Caddy (HTTPS) | Public |

### Key Commands

```bash
# Start everything
cd /opt/onemind/docker/cloud
docker-compose up -d

# View logs
docker-compose logs -f heartbeat

# Check status
docker-compose ps

# Restart a service
docker-compose restart letta

# Deploy from GitHub
git pull && docker-compose up -d

# Check Tailscale
tailscale status

# Access VM desktop
# From Mac: open vnc://onemind.tailnet:5900
# From browser: http://onemind.tailnet:6080
```

### Tailscale Names

| Device | Tailscale Name | IP |
|--------|---------------|-----|
| Cloud VM | onemind-cloud | 100.x.x.1 |
| Home VM | onemind-home | 100.x.x.10 |
| Office VM | onemind-office | 100.x.x.20 |
| Your Mac | zeus-mac | 100.x.x.100 |
| Your Phone | zeus-phone | 100.x.x.101 |

---

## Related Documents

- [ONEMIND-OS-ARCHITECTURE.md](./ONEMIND-OS-ARCHITECTURE.md) - System architecture with visual diagrams
- [Legacy-OPERATIONS-CENTER.md](./Legacy-OPERATIONS-CENTER.md) - Original operations center plan
- [ONEMIND-CODEX.md](./ONEMIND-CODEX.md) - Master OneMind blueprint

---

**Document Status:** Living document - will be updated as infrastructure evolves

*Last Updated: January 16, 2026*
