---
title: "Build Log: NATS SaaS Worker v1→v3"
type: log
status: complete
created: 2026-05-24
owner: zeus
tags: [type/log, domain/so, build-log, nats, composio]
ai_summary: "Full build session of NATS↔Composio SaaS Worker from v1 (4 endpoints) to v3 (18 endpoints). Covers architecture decisions, debugging, the SAAS_EVENTS stream fix, and final verification."
---

# Build Log: NATS ↔ Composio SaaS Worker (2026-05-24)

## Goal

Build a **universal SaaS execution layer** where any agent on the NATS fabric can execute 2000+ app actions without touching HTTP, MCP, or any framework-specific tooling.

## Timeline

### Phase 1: Foundation (v1.0.0)

**Decision**: Use Composio SDK as the universal executor. Rejected:
- Nango (19GB RAM, 9 CPU = absurd)
- Windmill (NATS trigger = paid Enterprise)
- Direct API integrations (maintenance nightmare)
- MCP (framework-locked, cold start, single-agent)

**Built**: `~/.hermes/workers/saas_action_worker.py`
- 4 endpoints: `saas.execute`, `saas.search`, `saas.tools`, `health.check`
- Initially used `composio` CLI subprocess (slow, brittle)

### Phase 2: SDK Direct (v2.0.0)

**Problem**: CLI subprocess added 2-3s latency per call, hard to parse output.

**Fix**: Rewrote to use `composio` Python SDK directly (in-process):
- `from composio import Composio`
- `c.tools.execute(slug=..., user_id=..., arguments=...)`
- Sub-100ms execution

**Critical Bug Found**: `SAAS_EVENTS` stream was configured as `saas.>` which captured ALL NATS traffic on `saas.*` subjects — including request-reply on `saas.execute`, `saas.tools`, etc. Messages were being consumed by JetStream before the worker could reply.

**Fix**: Changed SAAS_EVENTS stream subjects to `saas.event.>` (inbound webhooks only).

### Phase 3: Full Surface (v3.0.0)

**Goal**: Expose EVERY Composio SDK capability via NATS — not just execute.

**SDK Surface Discovered**:
- `c.tools` → execute, get, get_raw_composio_tool_by_slug, proxy
- `c.connected_accounts` → list, get, delete, initiate, refresh, enable, disable
- `c.toolkits` → list, get, authorize, list_categories
- `c.triggers` → list, list_active, create, enable, disable, delete, get_type
- `c.auth_configs` → list, get, create, delete, enable, disable

**Built 18 endpoints**:

| Subject | Handler | Purpose |
|---------|---------|---------|
| `saas.execute` | Execute any action | Core execution |
| `saas.tools` | List toolkit actions | Discovery |
| `saas.tools.get` | Get action schema | Schema lookup |
| `saas.tools.search` | Global action search | Cross-toolkit search |
| `saas.connections` | List OAuth accounts | Connection check |
| `saas.connections.get` | Get connection detail | Inspection |
| `saas.connections.delete` | Remove connection | Cleanup |
| `saas.toolkits` | List available apps | Catalog browse |
| `saas.toolkits.get` | Get toolkit detail | App info |
| `saas.triggers` | List trigger types | Trigger discovery |
| `saas.triggers.active` | List active instances | Active monitoring |
| `saas.triggers.create` | Create trigger | Event subscription |
| `saas.triggers.enable` | Enable trigger | Activation |
| `saas.triggers.disable` | Disable trigger | Deactivation |
| `saas.triggers.delete` | Delete trigger | Cleanup |
| `saas.auth_configs` | List auth configs | Auth inspection |
| `saas.auth_configs.get` | Get auth config | Detail |
| `health.check` | Health + version | Ops |

### Bug Fix: Tool Slug Parsing

**Problem**: `c.tools.get()` returns OpenAI function-calling format (nested `{"function":{"name":"...","description":"...","parameters":{...}}}`), not flat objects. The slug was being serialized as the entire dict.

**Fix**: Added dict detection — if result item is a dict, extract `function.name` as slug.

---

## Verification

All 18 endpoints tested and confirmed working:

```bash
# Health
nats req health.check '{}' → v3.0.0, healthy

# Connections (the key missing piece from v2)
nats req saas.connections '{}' → 9 active: instagram, facebook, tavily, google_maps, youtube, webflow, cloudflare, github, todoist

# Tools (clean slugs now)
nats req saas.tools '{"toolkit":"todoist","limit":10}' → TODOIST_ADD_WORKSPACE, TODOIST_ARCHIVE_PROJECT2, TODOIST_BULK_CREATE_TASKS...

# Search
nats req saas.tools.search '{"query":"create repository","limit":5}' → BITBUCKET_CREATE_BRANCH, GITHUB_CREATE_REPO...

# Triggers
nats req saas.triggers '{"toolkit":"github","limit":3}' → GITHUB_ARTIFACT_CREATED_TRIGGER, GITHUB_BRANCH_CHANGED_TRIGGER...

# Execute (live test)
nats req saas.execute '{"slug":"TODOIST_CREATE_TASK","args":{"content":"NATS worker pipeline fully operational ✅"}}' → task created
```

---

## What's Next

1. **Connect more apps**: Gmail, Calendar, Slack, Discord, Notion at platform.composio.dev
2. **Create triggers**: GitHub push → agent notification, new email → task creation
3. **Composio NATS Bridge**: `com.onemind.composio-nats-bridge` plist exists but not loaded — for inbound webhook events
4. **Synadia Platform**: 60-day trial available for control plane
5. **NetBird migration**: Replace Tailscale (not started)

---

## Files Modified/Created

| File | Action |
|------|--------|
| `~/.hermes/workers/saas_action_worker.py` | Rewritten v1→v2→v3 |
| `~/Library/LaunchAgents/com.onemind.saas-action-worker.plist` | Updated for Python 3.13 |
| `~/.hermes/skills/devops/nats-saas-worker/SKILL.md` | Created + updated to v3 |
| `~/.hermes/.env` | COMPOSIO_API_KEY added |
| `~/.hermes/profiles/*/. env` | API key propagated to all profiles |
