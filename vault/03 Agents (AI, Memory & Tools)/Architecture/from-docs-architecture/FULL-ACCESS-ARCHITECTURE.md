# OneMind OS - Full Access Architecture

> Complete integration stack for personal data access across all platforms

---

## Overview

OneMind OS uses a **6-layer integration architecture** to provide full access to your personal data across all platforms. Each layer serves a specific purpose and together they create a unified data foundation for Legacy and custom dashboards.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    FULL ACCESS ARCHITECTURE                                  │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                     Layer 6: CUSTOM DASHBOARDS                         │  │
│  │  Your UI consuming unified data from all layers below                  │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                     ▲                                        │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                     Layer 5: UNIFIED DATA STORE                        │  │
│  │  PostgreSQL + Redis (hot/cold) | NATS JetStream (events)               │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                     ▲                                        │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │   Layer 4   │ │   Layer 3   │ │   Layer 2   │ │   Layer 1   │           │
│  │  Browser    │ │   Nango     │ │  Composio   │ │  Standard   │           │
│  │  Extension  │ │   (Data)    │ │  (Actions)  │ │  Protocols  │           │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘           │
│         │               │               │               │                   │
│         ▼               ▼               ▼               ▼                   │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│  │ Web Apps    │ │ OAuth APIs  │ │ OAuth APIs  │ │ CalDAV      │           │
│  │ (Any site)  │ │ (500+)      │ │ (250+)      │ │ IMAP        │           │
│  │             │ │             │ │             │ │ RSS/Atom    │           │
│  │             │ │             │ │             │ │ WebDAV      │           │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Layer 1: Standard Protocols (Universal Access)

**Purpose:** Access data without OAuth or API keys using open standards.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     STANDARD PROTOCOLS                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Protocol      │ Use Case           │ Services                              │
│  ─────────────────────────────────────────────────────────────              │
│  CalDAV        │ Calendars          │ Google, iCloud, Fastmail, any         │
│  CardDAV       │ Contacts           │ Google, iCloud, Fastmail, any         │
│  IMAP/SMTP     │ Email              │ Gmail, Outlook, any email             │
│  RSS/Atom      │ Feeds              │ Any blog, news, podcast               │
│  WebDAV        │ Files              │ Nextcloud, iCloud, any                │
│  MQTT          │ IoT                │ Home Assistant, sensors               │
│  SIP           │ Phone              │ VoIP providers                        │
│                                                                              │
│  WHY USE THIS LAYER:                                                         │
│  ✓ No API registration required                                             │
│  ✓ No rate limits (your server, your rules)                                 │
│  ✓ Full data access (not limited by API)                                    │
│  ✓ Works offline with local sync                                            │
│  ✓ Vendor agnostic - switch providers anytime                               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### CalDAV/CardDAV Configuration

```yaml
# Standard protocol endpoints
caldav:
  google:
    url: https://apidata.googleusercontent.com/caldav/v2/{email}/events
    auth: oauth2  # Use Nango for OAuth
  icloud:
    url: https://caldav.icloud.com/{user}/calendars
    auth: app_password
  fastmail:
    url: https://caldav.fastmail.com/dav/calendars/user/{user}/
    auth: app_password

imap:
  gmail:
    host: imap.gmail.com
    port: 993
    auth: oauth2  # Use Nango for OAuth
  fastmail:
    host: imap.fastmail.com
    port: 993
    auth: app_password
```

---

## Layer 2: Composio (AI Actions)

**Purpose:** Enable Legacy to take actions on your behalf.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     COMPOSIO (AI ACTIONS)                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  License: MIT (fully open source)                                           │
│  Self-hosted: Yes                                                            │
│  APIs: 250+ integrations                                                     │
│                                                                              │
│  USE FOR:                                                                    │
│  ✓ AI agent actions (complete task, send message, create issue)             │
│  ✓ OAuth flow for actions                                                   │
│  ✓ MCP tool server                                                          │
│                                                                              │
│  DO NOT USE FOR:                                                             │
│  ✗ Data sync (limited triggers)                                             │
│  ✗ Continuous data pull (use Nango instead)                                 │
│  ✗ Real-time events (use webhooks instead)                                  │
│                                                                              │
│  FLOW:                                                                       │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  User: "Complete my Todoist task about groceries"                      │ │
│  │                           │                                            │ │
│  │                           ▼                                            │ │
│  │  Legacy → Composio Tool → OAuth Token → Todoist API → Task Done        │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Layer 3: Nango (Data Sync)

**Purpose:** Continuously sync data from OAuth-protected APIs for dashboards.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     NANGO (DATA SYNC)                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  License: Elastic License (self-hosted free with limits)                    │
│  Self-hosted: Yes                                                            │
│  APIs: 500+ integrations                                                     │
│                                                                              │
│  USE FOR:                                                                    │
│  ✓ Continuous data sync to your database                                    │
│  ✓ OAuth credential management                                              │
│  ✓ Request proxying with auto-refresh                                       │
│  ✓ Dashboard data sources                                                   │
│                                                                              │
│  KEY FEATURES:                                                               │
│  • Syncs: Pull data on schedule (every 5 min, hourly, etc.)                 │
│  • Actions: One-time API calls with managed auth                            │
│  • Proxying: Make API calls with auto-injected credentials                  │
│  • Webhooks: Receive and normalize incoming webhooks                        │
│                                                                              │
│  DATA FLOW:                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                         │ │
│  │   Todoist API ──┐                                                       │ │
│  │   GitHub API  ──┼──▶ Nango ──▶ PostgreSQL ──▶ Dashboard                │ │
│  │   Notion API  ──┘      │                                                │ │
│  │                        │                                                │ │
│  │                        └──▶ NATS JetStream (sync events)                │ │
│  │                                                                         │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  SYNC SCRIPTS (Custom TypeScript):                                           │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  // nango-integrations/todoist/syncs/tasks.ts                          │ │
│  │  export default async function fetchData(nango: NangoSync) {           │ │
│  │    const tasks = await nango.get({ endpoint: '/tasks' });              │ │
│  │    await nango.batchSave(tasks, 'TodoistTask');                        │ │
│  │  }                                                                      │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Nango vs Composio Decision Matrix

| Scenario | Use |
|----------|-----|
| Legacy needs to complete a task | **Composio** |
| Dashboard shows all tasks | **Nango** |
| AI creates a GitHub issue | **Composio** |
| Dashboard shows repo activity | **Nango** |
| Real-time webhook events | **Direct → NATS** |
| Scheduled data refresh | **Nango** |

---

## Layer 4: Browser Extension (Everything Else)

**Purpose:** Capture data from any web app that doesn't have an API.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     BROWSER EXTENSION                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  WHY: Some data only exists in web apps with no API access                  │
│                                                                              │
│  CAPABILITIES:                                                               │
│  ✓ Capture page content from any site                                       │
│  ✓ Extract structured data (orders, receipts, bookings)                     │
│  ✓ Monitor DOM changes (new messages, notifications)                        │
│  ✓ Intercept network requests (internal APIs)                               │
│  ✓ Auto-fill forms from Legacy                                              │
│                                                                              │
│  ARCHITECTURE:                                                               │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                         │ │
│  │   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐             │ │
│  │   │   Content   │     │ Background  │     │  OneMind    │             │ │
│  │   │   Script    │────▶│   Worker    │────▶│   API       │             │ │
│  │   │ (per page)  │     │ (extension) │     │             │             │ │
│  │   └─────────────┘     └─────────────┘     └─────────────┘             │ │
│  │         │                    │                   │                     │ │
│  │         │ DOM access         │ Storage           │ Events              │ │
│  │         │ XHR intercept      │ Auth token        │ to NATS             │ │
│  │         ▼                    ▼                   ▼                     │ │
│  │   Any webpage         Credential store     Unified events              │ │
│  │                                                                         │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  EXAMPLE USE CASES:                                                          │
│  • Amazon: Capture order history, tracking, prices                          │
│  • Banking: Capture transactions (no open banking API)                      │
│  • Social: Capture DMs, notifications, activity                             │
│  • Travel: Capture bookings, confirmations, itineraries                     │
│  • Receipts: Capture any online receipt                                     │
│                                                                              │
│  SECURITY:                                                                   │
│  • Extension runs locally only                                              │
│  • Data goes to YOUR server (self-hosted)                                   │
│  • No third-party access                                                    │
│  • Optional: Encrypted local storage                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Layer 5: Unified Data Store

**Purpose:** Single source of truth for all integrated data.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     UNIFIED DATA STORE                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        PostgreSQL                                    │    │
│  │                                                                      │    │
│  │  unified_items (polymorphic table for all integrated data)          │    │
│  │  ┌────────────────────────────────────────────────────────────────┐ │    │
│  │  │ id          │ source      │ type       │ data (JSONB)          │ │    │
│  │  ├─────────────┼─────────────┼────────────┼───────────────────────┤ │    │
│  │  │ item_001    │ todoist     │ task       │ {title, due, ...}     │ │    │
│  │  │ item_002    │ github      │ pr         │ {repo, number, ...}   │ │    │
│  │  │ item_003    │ gmail       │ email      │ {from, subject, ...}  │ │    │
│  │  │ item_004    │ caldav      │ event      │ {title, start, ...}   │ │    │
│  │  │ item_005    │ extension   │ receipt    │ {vendor, amount, ...} │ │    │
│  │  └────────────────────────────────────────────────────────────────┘ │    │
│  │                                                                      │    │
│  │  Key columns:                                                        │    │
│  │  • source: Integration ID (todoist, github, extension, etc.)        │    │
│  │  • type: Item type within source (task, pr, email, etc.)            │    │
│  │  • external_id: ID in the source system                             │    │
│  │  • data: Full payload as JSONB (queryable)                          │    │
│  │  • embedding: pgvector for semantic search                          │    │
│  │  • synced_at: Last sync timestamp                                   │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        Redis                                         │    │
│  │                                                                      │    │
│  │  Hot data for real-time access:                                     │    │
│  │  • Current awareness state                                          │    │
│  │  • Active integrations status                                       │    │
│  │  • Recent items cache                                               │    │
│  │  • Dashboard widget data                                            │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     NATS JetStream                                   │    │
│  │                                                                      │    │
│  │  Real-time events from all layers:                                  │    │
│  │  • events.todoist.task.* (from webhooks)                            │    │
│  │  • events.github.pr.* (from webhooks)                               │    │
│  │  • events.sync.nango.* (from Nango sync completion)                 │    │
│  │  • events.extension.* (from browser extension)                      │    │
│  │  • events.caldav.* (from CalDAV sync)                               │    │
│  │                                                                      │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Layer 6: Custom Dashboards

**Purpose:** Your UI consuming unified data.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     CUSTOM DASHBOARDS                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  DATA SOURCES FOR DASHBOARDS:                                                │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                         │ │
│  │  Source              │ Access Method          │ Use Case               │ │
│  │  ───────────────────────────────────────────────────────────────────── │ │
│  │  PostgreSQL          │ Direct query           │ Historical data        │ │
│  │  Redis               │ Direct query           │ Real-time state        │ │
│  │  NATS                │ WebSocket subscription │ Live updates           │ │
│  │  Nango               │ Proxy API              │ Fresh data on demand   │ │
│  │                                                                         │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  DASHBOARD API ENDPOINTS:                                                    │
│                                                                              │
│  GET  /dashboard/unified          - All items across sources                │
│  GET  /dashboard/tasks            - All tasks (Todoist, GitHub, etc.)       │
│  GET  /dashboard/events           - All calendar events                     │
│  GET  /dashboard/emails           - Recent emails                           │
│  GET  /dashboard/activity         - Cross-source activity feed              │
│  GET  /dashboard/stats            - Aggregated statistics                   │
│  WS   /dashboard/stream           - Real-time updates                       │
│                                                                              │
│  EXAMPLE DASHBOARD WIDGETS:                                                  │
│                                                                              │
│  ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐            │
│  │   Today's Tasks  │ │  Upcoming Events │ │  Unread Emails   │            │
│  │   ────────────── │ │  ────────────── │ │  ────────────── │            │
│  │   □ Review PR    │ │  10:00 Standup   │ │  3 from John     │            │
│  │   □ Deploy v2    │ │  14:00 1:1       │ │  1 from GitHub   │            │
│  │   ✓ Fix bug      │ │  16:00 Review    │ │  5 newsletters   │            │
│  └──────────────────┘ └──────────────────┘ └──────────────────┘            │
│                                                                              │
│  ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐            │
│  │   GitHub         │ │  Home Status     │ │  Recent Activity │            │
│  │   ────────────── │ │  ────────────── │ │  ────────────── │            │
│  │   2 PRs open     │ │  Lights: ON      │ │  PR merged       │            │
│  │   1 needs review │ │  Temp: 72°F      │ │  Task completed  │            │
│  │   5 issues       │ │  Motion: None    │ │  Email received  │            │
│  └──────────────────┘ └──────────────────┘ └──────────────────┘            │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Data Flow: Complete Picture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                     COMPLETE DATA FLOW                                       │
└─────────────────────────────────────────────────────────────────────────────┘

                              EXTERNAL WORLD
    ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
    │ Todoist │ │ GitHub  │ │ Gmail   │ │ Calendar│ │Web Apps │
    └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘
         │           │           │           │           │
    ┌────┴────┐ ┌────┴────┐ ┌────┴────┐ ┌────┴────┐ ┌────┴────┐
    │Webhook  │ │Webhook  │ │ IMAP    │ │ CalDAV  │ │Extension│
    └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘
         │           │           │           │           │
         └─────────────────┬─────────────────┘           │
                           │                             │
                           ▼                             ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        ONEMIND OS                                            │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                 Redpanda Connect (Port 4196)                         │    │
│  │                 /webhooks                                            │    │
│  │                                                                      │    │
│  │  Transforms webhooks → Standard format → NATS                       │    │
│  └──────────────────────────────┬──────────────────────────────────────┘    │
│                                 │                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                 Protocol Sync Services                               │    │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐                          │    │
│  │  │CalDAV Sync│ │ IMAP Sync │ │ RSS Sync  │                          │    │
│  │  └─────┬─────┘ └─────┬─────┘ └─────┬─────┘                          │    │
│  └────────┼─────────────┼─────────────┼────────────────────────────────┘    │
│           │             │             │                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                 Nango (Port 3003)                                    │    │
│  │                                                                      │    │
│  │  Scheduled syncs → PostgreSQL                                       │    │
│  │  Sync completion events → NATS                                      │    │
│  └──────────────────────────────┬──────────────────────────────────────┘    │
│                                 │                                            │
│           ┌─────────────────────┼─────────────────────┐                     │
│           │                     │                     │                     │
│           ▼                     ▼                     ▼                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │   PostgreSQL    │  │     Redis       │  │  NATS JetStream │             │
│  │   (cold data)   │  │   (hot data)    │  │    (events)     │             │
│  └────────┬────────┘  └────────┬────────┘  └────────┬────────┘             │
│           │                    │                    │                       │
│           └────────────────────┼────────────────────┘                       │
│                                │                                            │
│                                ▼                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │              Consciousness Processor                                 │    │
│  │                                                                      │    │
│  │  Filters events by awareness level, time, location                  │    │
│  └──────────────────────────────┬──────────────────────────────────────┘    │
│                                 │                                            │
│                    ┌────────────┼────────────┐                              │
│                    │            │            │                              │
│                    ▼            ▼            ▼                              │
│  ┌─────────────────────┐ ┌───────────┐ ┌───────────────────────┐           │
│  │    Legacy Inbox     │ │    UI     │ │  Custom Dashboards    │           │
│  │  (AI consciousness) │ │  (Chat)   │ │  (Your widgets)       │           │
│  └─────────────────────┘ └───────────┘ └───────────────────────┘           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘

                              ACTIONS (Outbound)

┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  User: "Complete the grocery task"                                          │
│                     │                                                        │
│                     ▼                                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        Legacy Agent                                  │    │
│  │                                                                      │    │
│  │  Uses Composio tool: todoist_complete_task(task_id="12345")         │    │
│  └──────────────────────────────┬──────────────────────────────────────┘    │
│                                 │                                            │
│                                 ▼                                            │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        Composio                                      │    │
│  │                                                                      │    │
│  │  1. Retrieves OAuth token                                           │    │
│  │  2. Makes API call to Todoist                                       │    │
│  │  3. Returns success                                                 │    │
│  └──────────────────────────────┬──────────────────────────────────────┘    │
│                                 │                                            │
│                                 ▼                                            │
│                          Todoist API                                         │
│                          Task completed!                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Service Configuration

### Docker Compose Services

| Service | Port | Purpose | Layer |
|---------|------|---------|-------|
| `nango` | 3003 | OAuth + Data sync | Layer 3 |
| `nango-postgres` | 5433 | Nango's database | Layer 3 |
| `caldav-sync` | - | CalDAV/CardDAV sync | Layer 1 |
| `imap-sync` | - | IMAP email sync | Layer 1 |
| `redpanda-connect` | 4196 | Webhook ingestion | Events |
| `postgres` | 5432 | Unified data store | Layer 5 |
| `redis` | 6379 | Hot data cache | Layer 5 |
| `nats` | 4222 | Event backbone | Layer 5 |
| `agno` | 8080 | Legacy agent + Composio | Layer 2 |

### Environment Variables

```bash
# Nango
NANGO_SECRET_KEY=your-secret-key
NANGO_SERVER_URL=http://nango:3003
NANGO_CALLBACK_URL=http://your-domain.com/nango/callback

# Composio (for AI actions)
COMPOSIO_API_KEY=your-composio-key

# CalDAV (Google)
GOOGLE_CALDAV_CLIENT_ID=your-client-id
GOOGLE_CALDAV_CLIENT_SECRET=your-client-secret

# IMAP (Gmail)
GMAIL_APP_PASSWORD=your-app-password
GMAIL_EMAIL=your-email@gmail.com

# Browser Extension
EXTENSION_API_KEY=your-extension-api-key
```

---

## Custom Dashboard Data Access

### API Examples

```typescript
// Fetch all tasks across sources
const tasks = await fetch('/api/dashboard/tasks')
// Returns: { todoist: [...], github: [...], notion: [...] }

// Fetch unified timeline
const timeline = await fetch('/api/dashboard/activity?limit=50')
// Returns: [{ source, type, timestamp, data }, ...]

// Subscribe to real-time updates
const ws = new WebSocket('/api/dashboard/stream')
ws.onmessage = (event) => {
  const { source, type, data } = JSON.parse(event.data)
  // Update dashboard widget
}

// Query unified items
const items = await fetch('/api/dashboard/unified?source=todoist&type=task')
// Returns: [{ id, source, type, data, synced_at }, ...]
```

### Direct PostgreSQL Queries

```sql
-- All tasks due today
SELECT data FROM unified_items
WHERE source IN ('todoist', 'github', 'notion')
  AND type = 'task'
  AND (data->>'due_date')::date = CURRENT_DATE;

-- Activity in last 24 hours
SELECT source, type, data, created_at
FROM unified_items
WHERE created_at > NOW() - INTERVAL '24 hours'
ORDER BY created_at DESC;

-- Semantic search across all data
SELECT data, 1 - (embedding <=> $1) as similarity
FROM unified_items
ORDER BY similarity DESC
LIMIT 10;
```

---

## Summary: When to Use Each Layer

| Need | Layer | Tool |
|------|-------|------|
| AI takes an action | Layer 2 | Composio |
| Dashboard shows data | Layer 3 | Nango |
| Calendar sync | Layer 1 | CalDAV |
| Email sync | Layer 1 | IMAP |
| RSS feeds | Layer 1 | RSS |
| Real-time events | - | Direct webhooks → NATS |
| Data from web-only apps | Layer 4 | Browser extension |
| Query historical data | Layer 5 | PostgreSQL |
| Real-time state | Layer 5 | Redis |
| Live updates | Layer 5 | NATS WebSocket |

---

## Webhook Endpoints (Integration → NATS)

All integration events flow through NATS. Services can send webhooks to:

### Primary Path: Redpanda Connect (Port 4196)
```bash
# Universal webhook endpoint - auto-detects service
curl -X POST http://localhost:4196/webhooks \
  -H "X-Webhook-Source: todoist" \
  -d '{"event_name": "item:completed", "event_data": {...}}'
```

### Backup Path: Direct Agno Endpoints (Port 8080)
```bash
# Nango sync events
curl -X POST http://localhost:8080/webhooks/nango \
  -d '{"type": "sync.completed", "providerConfigKey": "google", "connectionId": "123"}'
# → Publishes to: events.nango.google.sync_completed

# Composio action events
curl -X POST http://localhost:8080/webhooks/composio \
  -d '{"app_name": "todoist", "trigger": "task.completed", "data": {...}}'
# → Publishes to: events.composio.todoist.task_completed

# Generic webhook (use X-Webhook-Source header)
curl -X POST http://localhost:8080/webhooks/generic \
  -H "X-Webhook-Source: myservice" \
  -d '{"event": "something.happened"}'
# → Publishes to: events.myservice.webhook
```

### NATS Subject Patterns
| Source | Subject Pattern | Example |
|--------|-----------------|---------|
| Protocol Sync | `events.sync.{protocol}.*` | `events.sync.caldav.new_event` |
| Extension | `events.extension.*` | `events.extension.capture` |
| Nango | `events.nango.{provider}.*` | `events.nango.google.sync_completed` |
| Composio | `events.composio.{app}.*` | `events.composio.todoist.task_completed` |
| Webhooks | `events.{service}.*` | `events.github.pr.opened` |

### Configure Webhook URLs in External Services
| Service | Webhook URL |
|---------|-------------|
| Todoist | `http://your-server:4196/webhooks` (with `X-Webhook-Source: todoist`) |
| GitHub | `http://your-server:4196/webhooks` (auto-detected) |
| Nango | `http://your-server:8080/webhooks/nango` |
| Composio | `http://your-server:8080/webhooks/composio` |
| Any Other | `http://your-server:8080/webhooks/generic` (with header) |

---

## Next Steps

1. Deploy Nango: `docker compose --profile data up -d`
2. Configure CalDAV/IMAP in `/settings/protocols`
3. Install browser extension from `/settings/extension`
4. Configure webhooks in external services (Todoist, GitHub, etc.)
5. Build dashboard widgets using the unified API
