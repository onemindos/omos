---
id: so-02-20260527-session-persistence-nats
title: "Session Persistence — Synadia Agent SDK Native Protocol"
type: reference
status: active
created: 2026-05-27
modified: 2026-05-27
owner: zeus
pillar: SO
domain: "02 Infrastructure (The Body)"
tags: [type/reference, zone/so, domain/so, status/active, infra/nats, agents/sessions, protocol/acp]
depth: 0
---

# Session Persistence — Native Protocol via Synadia Agent SDK

## Subject Pattern

```
agents.prompt.<agent>.<owner>.<session_name>
```

The `session_name` segment IS the persistence key. Each harness maps it to internal session storage.

## Lifecycle

| Action | Method |
|--------|--------|
| **Start session** | Send to `agents.prompt.onemind-ceo.zeus.q2-strategy` — creates if new |
| **Continue session** | Send to SAME subject — harness loads existing context |
| **List sessions** | `nats_request("agents.sessions.onemind-ceo")` → returns active sessions |
| **Close session** | `nats_request("agents.sessions.onemind-ceo.close", {session: "q2-strategy"})` |
| **Query session** | `nats_request("agents.sessions.onemind-ceo.get", {session: "q2-strategy"})` |

## How It Works

1. Legacy (or any authorized client) sends a message to a session subject
2. Target harness receives via Synadia Agent SDK subscription
3. Harness checks if `session_name` exists in its DB
   - If yes: loads full conversation context, appends new message, generates response
   - If no: creates new session, stores in DB, generates first response
4. Response streams back chunk-by-chunk via SDK protocol
5. Session state persists in harness DB (SQLite for Hermes, OpenClaw's internal DB)
6. Session remains active until explicitly closed or TTL expires

## Cross-Harness Consistency

Both Hermes and OpenClaw implement this same subject pattern via Synadia Agent SDK:

```
Legacy (Hermes)
    │
    │  agents.prompt.<agent>.<owner>.<session>
    │
    ├──→ OneMind-CEO (OpenClaw + Synadia SDK)
    ├──→ Greenhouse (Hermes + Synadia SDK)
    ├──→ Claude Code (ACP over NATS)
    └──→ Drone Fleet (Hermes leaf + NATS)

ALL agents are equal citizens on the bus.
Legacy doesn't know or care which harness runs them.
```

## Existing Implementation

Already built in: `~/.hermes/plugins/nats-platform/adapter.py`

The adapter registers as an AgentService and subscribes to `agents.prompt.<agent>.<owner>.<session_name>`. The SDK handles heartbeat, terminator tokens, and status endpoints automatically.

## Session Metadata (stored per session)

```json
{
  "session_id": "q2-strategy",
  "agent": "onemind-ceo",
  "owner": "zeus",
  "created": "2026-05-27T14:30:00Z",
  "last_active": "2026-05-27T16:45:00Z",
  "status": "active",
  "message_count": 47,
  "total_tokens": 128000,
  "cost_usd": 1.23
}
```
