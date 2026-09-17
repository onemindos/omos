---
title: "Composio — Universal SaaS Integration Layer"
type: reference
status: active
created: 2026-05-24
updated: 2026-05-24
owner: zeus
tags: [type/reference, domain/so, status/active, so/02-infrastructure, integration, composio]
ai_summary: "Composio SDK provides OAuth-managed access to 2000+ SaaS apps. Connected via NATS worker for fleet-wide use. 9 apps active. Dashboard at platform.composio.dev."
---

# Composio — Universal SaaS Integration Layer

## What It Is

Composio is a **managed OAuth + action execution platform** that gives OneMind OS access to 2000+ SaaS applications through a single Python SDK. Instead of building individual API integrations, we connect OAuth once and get pre-mapped actions for every app.

## Why Composio

| Alternative | Why Rejected |
|-------------|-------------|
| **Nango** | 19GB RAM, 9 CPU cores, 5 Docker containers. Insane for a single-user system |
| **Direct API integrations** | Thousands of hours maintaining OAuth flows, token refresh, schema changes |
| **MCP tools** | Framework-locked to one agent, no fleet-wide sharing, cold start |
| **Windmill** | NATS trigger is paid Enterprise feature |

Composio = **one SDK call** with managed OAuth, token refresh, and 10,000+ pre-mapped actions.

---

## Architecture in OneMind OS

```
Composio Cloud (OAuth tokens, action schemas)
        ↕ HTTPS
Composio SDK (Python 3.13, in-process)
        ↕ in-memory
SaaS Action Worker (~/.hermes/workers/saas_action_worker.py)
        ↕ NATS request-reply
Any Agent on the Fleet
```

The worker is the **bridge** — agents never touch Composio directly. They publish NATS requests, get responses. Pure decoupling.

---

## Credentials

| Item | Value / Location |
|------|-----------------|
| API Key | `~/.composio/user_data.json` (also in `~/.hermes/.env`) |
| User ID | `pg-test-dc73c699-2bd3-477d-9952-cf8bfe545bd9` |
| Dashboard | https://platform.composio.dev |
| API Base | https://backend.composio.dev |
| CLI Binary | `~/.composio/composio` (v0.2.28) |
| Python SDK | `composio` v0.13.1, `composio-client` v1.39.0 |

⚠️ The API key is also copied to all Hermes profiles: `~/.hermes/profiles/{unity,eternal,sentinel,sovereignty,zenith}/.env`

---

## Connected Apps (May 2026)

| App | Connection ID | Use Case |
|-----|---------------|----------|
| **Todoist** | `ca_wFtzOvD-wYnF` | Task management (NATS project) |
| **GitHub** | `ca_41kOMYQ0KTq0` | Repo/issue/PR management |
| **Cloudflare** | `ca_uArdAq7tKu3U` | DNS, tunnels, workers |
| **Webflow** | `ca_ECZVeKqSSMaX` | zeusdelacruz.com site management |
| **YouTube** | `ca_NrxyqiwaRrec` | Video/channel management |
| **Instagram** | `ca_DuTng_YD4YFH` | Social media posting |
| **Facebook** | `ca_rxE3XgEY1jmn` | Social media posting |
| **Google Maps** | `ca_209gwBwJCqXm` | Location/place queries |
| **Tavily** | `ca_2whmYaXjb_H6` | Web search/research |

### Priority Apps to Connect Next
- Gmail (email automation)
- Google Calendar (scheduling)
- Slack (team comms)
- Discord (community)
- Notion (knowledge base)
- Linear (project management)
- Stripe (payments)

---

## How to Connect a New App

1. Open https://platform.composio.dev
2. Click "Connected Accounts" → "Add Connection"
3. Select the app, complete OAuth
4. **Immediately available** via `saas.execute` on NATS — zero code changes

---

## SDK Surface (exposed via NATS)

| SDK Module | NATS Subjects | Operations |
|------------|---------------|------------|
| `tools` | `saas.execute`, `saas.tools`, `saas.tools.get`, `saas.tools.search` | Execute actions, browse/search action catalog |
| `connected_accounts` | `saas.connections`, `saas.connections.get`, `saas.connections.delete` | List/inspect/remove OAuth connections |
| `toolkits` | `saas.toolkits`, `saas.toolkits.get` | Browse 2000+ available apps |
| `triggers` | `saas.triggers`, `saas.triggers.active`, `saas.triggers.create/enable/disable/delete` | Event subscriptions |
| `auth_configs` | `saas.auth_configs`, `saas.auth_configs.get` | Auth configuration inspection |

---

## Common Patterns

### Execute an action (any agent)
```bash
nats req saas.execute '{"slug":"GITHUB_CREATE_ISSUE","args":{"repo":"OneMind-OS/One-Mind","title":"Bug fix","body":"Details..."}}'
```

### Discover what you can do
```bash
# What apps are available?
nats req saas.toolkits '{"limit":100}'

# What can I do with GitHub?
nats req saas.tools '{"toolkit":"github","limit":50}'

# Search for specific capability
nats req saas.tools.search '{"query":"create calendar event"}'
```

### Verify connection before executing
```bash
# Is GitHub connected?
nats req saas.connections '{"toolkit":"github"}'
```

---

## Important Caveats

1. **`saas.tools` ≠ connected**: Returns Composio's full action catalog regardless of whether you have OAuth. Always verify with `saas.connections` first.

2. **Action slugs are UPPERCASE**: `TODOIST_CREATE_TASK`, `GITHUB_CREATE_ISSUE` — never lowercase.

3. **Token refresh is automatic**: Composio handles OAuth token refresh internally. If you get auth errors, the OAuth grant itself may have been revoked — reconnect at the dashboard.

4. **v3 API quirk**: `/api/v3/connected_accounts` returns null `appName` — you must fetch each connection individually by ID to get the toolkit slug. The NATS worker handles this internally.

5. **SAAS_EVENTS stream**: Only captures `saas.event.>` (inbound webhook events). NOT `saas.>` which would intercept request-reply traffic.
