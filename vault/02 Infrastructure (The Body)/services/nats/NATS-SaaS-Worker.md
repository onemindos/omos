---
title: NATS ↔ Composio Universal SaaS Worker
version: 3.0.0
status: active
created: 2026-05-24
updated: 2026-05-24
domain: "02"
tags: [nats, composio, saas, integration, fleet, worker]
---

# NATS ↔ Composio Universal SaaS Worker v3.0.0

## Overview

The SaaS Worker is the **universal execution bridge** between the OneMind NATS fabric and 2000+ SaaS applications via the Composio SDK. Any agent on the fleet (Hermes/Legacy, OpenClaw, Pi, Claude Code) can execute actions, browse toolkits, manage OAuth connections, and configure event triggers — **purely via NATS request-reply**. Zero HTTP, zero framework dependency.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        NATS Fabric                                │
│                                                                   │
│  Hermes ─┐                                                        │
│  OpenClaw─┼── request → saas.* subjects ──→ SaaS Worker ──→ Composio SDK ──→ SaaS APIs │
│  Pi Agent─┤                              ←── reply ──────────────────────────────────── │
│  Claude  ─┘                                                       │
└─────────────────────────────────────────────────────────────────┘
```

### Design Principles

1. **Framework-agnostic**: Any NATS client can use it — no Hermes/MCP dependency
2. **Fleet-wide**: One worker serves all agents simultaneously
3. **Zero code changes for new apps**: Connect OAuth at platform.composio.dev → immediately usable
4. **Native process**: Runs via launchd, not Docker — sovereign, local, fast
5. **Request-reply pattern**: Sub-millisecond routing, no polling

## File Locations

| File | Path | Purpose |
|------|------|---------|
| Worker script | `~/.hermes/workers/saas_action_worker.py` | Main Python worker (v3.0.0) |
| Launchd plist | `~/Library/LaunchAgents/com.onemind.saas-action-worker.plist` | Auto-start on boot |
| Logs | `/tmp/saas-action-worker.log` | Stdout/stderr |
| Config | `~/.hermes/.env` | `COMPOSIO_API_KEY` |
| SDK config | `~/.composio/user_data.json` | Composio API key + base URL |
| Skill doc | `~/.hermes/skills/devops/nats-saas-worker/SKILL.md` | Agent-readable reference |

## Requirements

- Python 3.13 (`/opt/homebrew/bin/python3.13`)
- `composio` SDK (v0.13.1+)
- `composio-client` (v1.39.0+)
- `nats-py` (async NATS client)
- NATS server running on `localhost:4222`

---

## Complete NATS API Reference (18 Endpoints)

### 1. Action Execution

#### `saas.execute` — Execute any SaaS action

The core endpoint. Execute any of 10,000+ pre-mapped actions across 2000+ apps.

```json
// Request
{
  "slug": "TODOIST_CREATE_TASK",
  "args": {
    "content": "Buy groceries",
    "project_id": "6ghf274mWPRfMV3h",
    "priority": 3
  },
  "user_id": "pg-test-dc73c699-2bd3-477d-9952-cf8bfe545bd9"  // optional, defaults to env
}

// Response
{
  "ok": true,
  "data": {
    "id": "6ghg3JVMVGPqjcv4",
    "content": "Buy groceries",
    "project_id": "6ghf274mWPRfMV3h"
  }
}
```

**Key points:**
- `slug` is UPPERCASE with underscores (e.g. `GITHUB_CREATE_ISSUE`, `TODOIST_CREATE_TASK`)
- `args` maps directly to the action's parameter schema (use `saas.tools.get` to discover)
- `user_id` is optional — defaults to the configured Composio user

#### `saas.tools` — List actions for a toolkit

```json
// Request
{"toolkit": "github", "limit": 50}
{"toolkit": ["todoist", "github"], "limit": 20}
{"search": "create issue", "limit": 10}

// Response
{
  "ok": true,
  "tools": [
    {"slug": "GITHUB_CREATE_ISSUE", "description": "Creates a new issue..."},
    {"slug": "GITHUB_CREATE_PR", "description": "Creates a pull request..."}
  ],
  "total": 50
}
```

#### `saas.tools.get` — Get full action schema

Returns the complete parameter schema for a specific action.

```json
// Request
{"slug": "GITHUB_CREATE_ISSUE"}

// Response
{
  "ok": true,
  "tool": {
    "slug": "GITHUB_CREATE_ISSUE",
    "description": "Creates a new issue in a GitHub repository...",
    "parameters": {
      "type": "object",
      "required": ["repo", "title"],
      "properties": {
        "repo": {"type": "string", "description": "Repository in owner/name format"},
        "title": {"type": "string", "description": "Issue title"},
        "body": {"type": "string", "description": "Issue body (markdown)"}
      }
    }
  }
}
```

#### `saas.tools.search` — Search actions globally

Search across ALL toolkits (not just connected ones).

```json
// Request
{"query": "send email", "limit": 10}

// Response
{
  "ok": true,
  "tools": [
    {"slug": "GMAIL_SEND_EMAIL", "description": "Send an email..."},
    {"slug": "MAILCHIMP_SEND_CAMPAIGN", "description": "..."}
  ],
  "total": 10
}
```

---

### 2. Connection Management

#### `saas.connections` — List connected accounts

```json
// Request (all optional filters)
{
  "toolkit": "github",        // filter by toolkit slug
  "status": "ACTIVE",         // ACTIVE | INITIALIZING | EXPIRED | FAILED
  "statuses": ["ACTIVE", "EXPIRED"],  // multiple statuses
  "limit": 50
}

// Minimal — list all active
{}

// Response
{
  "ok": true,
  "connections": [
    {"id": "ca_41kOMYQ0KTq0", "toolkit": "github", "status": "ACTIVE", "created_at": "2026-05-24T13:02:38.498Z", "user_id": "pg-test-..."},
    {"id": "ca_wFtzOvD-wYnF", "toolkit": "todoist", "status": "ACTIVE", "created_at": "2026-03-05T09:20:31.480Z", "user_id": "pg-test-..."}
  ],
  "total": 9
}
```

#### `saas.connections.get` — Get connection details

```json
{"id": "ca_41kOMYQ0KTq0"}
```

#### `saas.connections.delete` — Remove a connection

```json
{"id": "ca_xxx"}
```

---

### 3. Toolkit Discovery

#### `saas.toolkits` — List all available apps

Browse Composio's full catalog of 2000+ integrations.

```json
// Request
{"limit": 50, "category": "crm", "sort_by": "usage"}

// Response
{
  "ok": true,
  "toolkits": [
    {"slug": "gmail", "name": "Gmail", "description": "...", "logo": "...", "categories": [...]},
    {"slug": "github", "name": "GitHub", ...}
  ],
  "total": 50
}
```

#### `saas.toolkits.get` — Get toolkit details

```json
{"slug": "github"}
// Returns: auth requirements, available actions count, triggers, etc.
```

---

### 4. Trigger Management

Triggers let you subscribe to events from SaaS apps (new GitHub push, new email, etc.).

#### `saas.triggers` — List available trigger types

```json
// Request
{"toolkit": "github", "limit": 10}

// Response
{
  "ok": true,
  "triggers": [
    {"slug": "GITHUB_PUSH_TRIGGER", "name": "Push Event", "description": "Triggers on push...", "toolkit": "github", "type": "webhook"},
    {"slug": "GITHUB_ISSUE_CREATED_TRIGGER", "name": "New Issue", ...}
  ],
  "total": 10
}
```

#### `saas.triggers.active` — List active trigger instances

```json
{"show_disabled": true, "limit": 50}
```

#### `saas.triggers.create` — Create a trigger

```json
{
  "slug": "GITHUB_PUSH_TRIGGER",
  "trigger_config": {"repo": "OneMind-OS/One-Mind", "branch": "main"},
  "connected_account_id": "ca_41kOMYQ0KTq0"
}
```

#### `saas.triggers.enable` — Enable a trigger

```json
{"trigger_id": "trig_xxx"}
```

#### `saas.triggers.disable` — Disable a trigger

```json
{"trigger_id": "trig_xxx"}
```

#### `saas.triggers.delete` — Delete a trigger

```json
{"trigger_id": "trig_xxx"}
```

---

### 5. Auth Configuration

#### `saas.auth_configs` — List auth configurations

```json
{"toolkit": "github"}
```

#### `saas.auth_configs.get` — Get specific auth config

```json
{"id": "auth_config_id"}
```

---

### 6. System

#### `health.check` — Health check

```json
// Request
{}

// Response
{
  "ok": true,
  "service": "saas-actions",
  "version": "3.0.0",
  "composio_user_id": "pg-test-dc73c699-2bd3-477d-9952-cf8bfe545bd9",
  "status": "healthy",
  "timestamp": "2026-05-24T13:17:34.668944Z"
}
```

---

## Currently Connected Apps (May 2026)

| Toolkit | Connection ID | Status |
|---------|---------------|--------|
| Instagram | `ca_DuTng_YD4YFH` | ✅ ACTIVE |
| Facebook | `ca_rxE3XgEY1jmn` | ✅ ACTIVE |
| Tavily | `ca_2whmYaXjb_H6` | ✅ ACTIVE |
| Google Maps | `ca_209gwBwJCqXm` | ✅ ACTIVE |
| YouTube | `ca_NrxyqiwaRrec` | ✅ ACTIVE |
| Webflow | `ca_ECZVeKqSSMaX` | ✅ ACTIVE |
| Cloudflare | `ca_uArdAq7tKu3U` | ✅ ACTIVE |
| GitHub | `ca_41kOMYQ0KTq0` | ✅ ACTIVE |
| Todoist | `ca_wFtzOvD-wYnF` | ✅ ACTIVE |

### Not Connected (Available to Add)
- Gmail, Google Calendar, Slack, Discord, Notion, Linear, Airtable, Stripe, and 2000+ more

---

## Operations

### Start/Stop/Restart

```bash
# Check status
launchctl list | grep saas

# Restart worker
launchctl unload ~/Library/LaunchAgents/com.onemind.saas-action-worker.plist
launchctl load ~/Library/LaunchAgents/com.onemind.saas-action-worker.plist

# View logs
cat /tmp/saas-action-worker.log
tail -f /tmp/saas-action-worker.log
```

### Testing via nats CLI

```bash
# Health check
nats req health.check '{}'

# List connections
nats req saas.connections '{}'

# Execute action
nats req saas.execute '{"slug":"TODOIST_CREATE_TASK","args":{"content":"Test from CLI"}}'

# Search tools
nats req saas.tools.search '{"query":"create issue","limit":5}'
```

---

## Adding New SaaS Apps

1. Go to https://platform.composio.dev
2. Click "Connected Accounts" → "Add Connection"
3. Complete OAuth flow for the app
4. **Done** — immediately callable via `saas.execute` with zero code changes

The worker uses the Composio SDK which routes to the right OAuth token automatically based on the action's toolkit.

---

## How It Relates to Other NATS Services

| Service | Subject Pattern | Purpose |
|---------|----------------|---------|
| **SaaS Worker** (this) | `saas.*` | Execute SaaS actions via Composio |
| **Agent Harness** | `agents.prompt.*` | Route prompts to AI agents |
| **Agent Heartbeat** | `agents.hb.*` | Agent liveness monitoring |
| **SAAS_EVENTS stream** | `saas.event.>` | Inbound webhook events from triggers |
| **AGENT_EVENTS stream** | `agent.>` | Agent activity log |

---

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| **NATS over MCP** | Sub-ms latency, multi-agent, no cold start, event streaming |
| **Composio over Nango** | Nango = 19GB RAM, 9 CPU, 5 containers. Composio = one SDK call |
| **Python 3.13** | Composio SDK requires 3.10+; latest gives best async perf |
| **Launchd over Docker** | Native process, KeepAlive, zero overhead, sovereign |
| **Single worker, all subjects** | One process handles everything — simpler ops, one log file |
| **`saas.event.>` stream separate** | Was `saas.>` which intercepted request-reply — CRITICAL fix |

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-05-24 | Initial worker: `saas.execute`, `saas.search`, `saas.tools`, `health.check` |
| 2.0.0 | 2026-05-24 | Refactored to Composio SDK direct (was CLI subprocess) |
| 3.0.0 | 2026-05-24 | Full SDK surface: 18 endpoints — connections, toolkits, triggers, auth_configs |

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| `saas.tools` returns actions for apps not connected | Normal — it returns the catalog. Use `saas.connections` to verify live OAuth |
| Worker not responding | `launchctl list \| grep saas` — restart if exit code ≠ 0 |
| Action fails with auth error | Check `saas.connections` — OAuth may be expired. Reconnect at platform.composio.dev |
| Trigger not firing | Ensure trigger is created AND enabled. Check `saas.triggers.active` |
| SAAS_EVENTS capturing request-reply | Stream must be `saas.event.>` ONLY — not `saas.>` |

---

## Composio Credentials

| Item | Value |
|------|-------|
| User ID | `pg-test-dc73c699-2bd3-477d-9952-cf8bfe545bd9` |
| API Key | Stored in `~/.composio/user_data.json` and `~/.hermes/.env` |
| Dashboard | https://platform.composio.dev |
| Base URL | https://backend.composio.dev |
