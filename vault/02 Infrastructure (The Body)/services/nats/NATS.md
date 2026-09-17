---
title: "NATS — OneMind Nervous System"
type: reference
status: active
created: 2026-04-15
updated: 2026-05-24
owner: zeus
tags: [type/reference, domain/so, status/active, so/02-infrastructure]
ai_summary: "NATS JetStream running natively on Zeus's MBP M2 Max. Central message bus for the entire OneMind OS fleet — agents, SaaS integrations, IoT, events. 11 streams, 18 SaaS subjects, agent routing."
---

# NATS — OneMind Nervous System

> **Status: LIVE** — Running natively on Zeus's MBP M2 Max (localhost:4222). The real-time backbone of OneMind OS.

---

## What NATS Is

NATS is the **central nervous system** of OneMind OS. Every agent, every SaaS integration, every device command flows through NATS. It replaced the old VPS-hosted deployment (archived April 2026) with a native local install for zero-latency, always-on operation.

```
┌──────────────────────────────────────────────────────────┐
│                  NATS (localhost:4222)                     │
│                                                           │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌────────────┐  │
│  │ Hermes  │  │OpenClaw │  │ Pi Agent│  │Claude Code │  │
│  │(Legacy) │  │         │  │         │  │            │  │
│  └────┬────┘  └────┬────┘  └────┬────┘  └─────┬──────┘  │
│       │             │            │              │         │
│       └─────────────┴────────────┴──────────────┘         │
│                         │                                 │
│  ┌──────────────────────┼──────────────────────────────┐  │
│  │           JetStream (11 Streams)                    │  │
│  │  AGENT_EVENTS · COMMS · TASKS · SAAS_EVENTS         │  │
│  │  AGENT_PROMPTS · COMMANDS · HEALTH_EVENTS           │  │
│  │  SECURITY_EVENTS · SYSTEM · ZEUS_EVENTS             │  │
│  │  SUBSCRIBER_EVENTS                                  │  │
│  └─────────────────────────────────────────────────────┘  │
│                         │                                 │
│  ┌──────────────────────┼──────────────────────────────┐  │
│  │     SaaS Worker (18 endpoints)                      │  │
│  │  saas.execute · saas.tools · saas.connections       │  │
│  │  saas.toolkits · saas.triggers · health.check       │  │
│  │  → Composio SDK → 2000+ SaaS APIs                  │  │
│  └─────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

---

## Installation

NATS runs as a native process on macOS, managed by launchd.

- **Binary**: `/opt/homebrew/bin/nats-server`
- **Config**: `~/.nats/nats-server.conf` (JetStream enabled)
- **Data dir**: `~/.nats/jetstream/`
- **Port**: 4222 (client), 8222 (monitoring)

---

## JetStream Streams (11)

| Stream | Subjects | Purpose |
|--------|----------|---------|
| `AGENT_EVENTS` | `agent.>` | Agent activity log (30-day retention) |
| `COMMS` | `comms.>` | Inter-agent communication |
| `TASKS` | `tasks.>` | Task queue and results |
| `SAAS_EVENTS` | `saas.event.>` | Inbound webhook events from Composio triggers |
| `AGENT_PROMPTS` | `agents.prompt.>` | Prompt routing to agents |
| `COMMANDS` | `cmd.>` | Command dispatch |
| `HEALTH_EVENTS` | `health.>` | Health/heartbeat events |
| `SECURITY_EVENTS` | `security.>` | Security alerts |
| `SYSTEM` | `system.>` | System-level events |
| `ZEUS_EVENTS` | `zeus.>` | Zeus-specific notifications |
| `SUBSCRIBER_EVENTS` | `subscriber.>` | Subscription management |

⚠️ **Critical**: `SAAS_EVENTS` must use `saas.event.>` — NOT `saas.>`. The broader pattern intercepts request-reply traffic on `saas.execute`, `saas.tools`, etc.

---

## Agent Communication

### Synadia Agent Harnesses

Located at `~/.nats/synadia-agents/agents/`:

| Agent | Token | Subscribe Pattern |
|-------|-------|-------------------|
| Claude Code | `cc` | `agents.prompt.cc.*` |
| Hermes | `hermes` | `agents.prompt.hermes.*` |
| OpenClaw | `openclaw` | `agents.prompt.openclaw.*` |
| Pi Agent | `pi` | `agents.prompt.pi.*` |

Heartbeats: `agents.hb.<token>` (every 30s)

### Request-Reply Patterns

| Pattern | Used For |
|---------|----------|
| `agents.prompt.<token>.<owner>.<name>` | Route work to specific agent |
| `saas.*` | SaaS worker (18 subjects) |
| `health.check` | Worker health |
| `cmd.ha.*` | Home Assistant commands |

---

## SaaS Worker (v3.0.0)

See [[NATS-SaaS-Worker]] for complete documentation.

**Summary**: 18 NATS endpoints exposing the entire Composio SDK:
- Execute any of 10,000+ actions (`saas.execute`)
- Browse 2000+ toolkits (`saas.toolkits`)
- List/search actions (`saas.tools`, `saas.tools.search`)
- Manage OAuth connections (`saas.connections`)
- Full trigger lifecycle (`saas.triggers.*`)
- Auth config inspection (`saas.auth_configs`)

**Connected Apps**: Instagram, Facebook, Google Maps, Tavily, YouTube, Webflow, Cloudflare, GitHub, Todoist

---

## Related Services

| Service | Plist | Status |
|---------|-------|--------|
| SaaS Worker | `com.onemind.saas-action-worker` | ✅ Running |
| Synadia Agent Web UI | manual (port 3300) | ✅ Running |
| Composio NATS Bridge | `com.onemind.composio-nats-bridge` | ⏸️ Not loaded |
| Redpanda Connect | `com.onemind.rpk-connect` | ⏸️ Not loaded |

---

## Key Decisions

| Decision | Date | Rationale |
|----------|------|-----------|
| NATS over MCP for SaaS | 2026-05-24 | Sub-ms, multi-agent, event streaming, no cold start |
| Native over Docker | 2026-04 | Zero overhead, launchd KeepAlive, sovereignty |
| Composio over Nango | 2026-05-24 | Nango = 19GB RAM. Composio = one SDK call |
| `saas.event.>` stream fix | 2026-05-24 | Broader `saas.>` intercepted request-reply |
| Python 3.13 for worker | 2026-05-24 | Composio SDK requires 3.10+, async perf |
| Synadia Platform trial | Pending | 60-day trial available, control plane candidate |

---

## Operations

```bash
# NATS server
nats-server --config ~/.nats/nats-server.conf

# Check streams
nats stream ls

# Check subjects
nats sub "agent.>" --last=5

# SaaS worker
launchctl list | grep saas
launchctl unload ~/Library/LaunchAgents/com.onemind.saas-action-worker.plist
launchctl load ~/Library/LaunchAgents/com.onemind.saas-action-worker.plist

# Test
nats req health.check '{}'
nats req saas.connections '{}'
nats req saas.execute '{"slug":"TODOIST_CREATE_TASK","args":{"content":"Test"}}'
```

---

## Migration History

| Date | Event |
|------|-------|
| 2026-04-15 | NATS on VPS `104.248.4.212` documented |
| 2026-04 | VPS decommissioned, NATS moved to local MBP |
| 2026-05-24 | SaaS Worker v1→v2→v3 built (Composio SDK, 18 endpoints) |
| 2026-05-24 | This doc updated from "archived" to "active" |
