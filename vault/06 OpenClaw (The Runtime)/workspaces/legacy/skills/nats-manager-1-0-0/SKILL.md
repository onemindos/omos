---
name: nats-manager
description: "NATS server monitoring and JetStream management. Use when Legacy needs to check nervous system health, inspect streams/consumers, view connections, create streams, or manage the NATS infrastructure. Complements nexus-bridge (which handles messaging/dispatch) with observability and management."
metadata: { "clawdbot": { "emoji": "🧠" } }
---

# NATS Manager — Nervous System Control

Monitor and manage the NATS JetStream infrastructure that connects all OneMind OS runtimes.

## When to Use

- "How's the nervous system?" / "Is NATS healthy?" → `nats_health`
- "Who's connected?" / "Show me NATS connections" → `nats_connections`
- "What streams exist?" / "Show JetStream status" → `nats_streams`
- "Tell me about the SYSTEM stream" → `nats_stream_detail`
- "What consumers are on DIGITAL?" → `nats_consumers`
- "What subjects are being listened to?" → `nats_subscriptions`
- "Create a stream for robot events" → `nats_stream_create`
- "Purge stale messages from SYSTEM" → `nats_stream_purge`

## Tools (8)

### Read-Only Monitoring (HTTP API — no NATS connection needed)

- `nats_health()` — Server vitals: uptime, CPU, memory, message rates, connection count
- `nats_connections(sort?, limit?)` — All connected clients with stats
- `nats_subscriptions(test_subject?)` — Active subscription patterns

### JetStream Management (NATS client connection)

- `nats_streams()` — List all streams with message counts and storage
- `nats_stream_detail(stream)` — Deep dive on one stream + its consumers
- `nats_consumers(stream)` — List consumers on a stream
- `nats_stream_create(stream, subjects, max_age_hours?, max_bytes?, storage?)` — Create a new stream
- `nats_stream_purge(stream)` — Delete all messages from a stream (destructive!)

## Relationship to nexus-bridge

| Skill            | Purpose                                                              |
| ---------------- | -------------------------------------------------------------------- |
| **nexus-bridge** | Messaging — dispatch tasks, publish events, query agents via NATS    |
| **nats-manager** | Infrastructure — monitor health, manage streams, inspect connections |

Use `nexus-bridge` to talk through the nervous system. Use `nats-manager` to manage the nervous system itself.

## Environment

- `NATS_URL` — NATS server URL for client connection (default: `nats://localhost:4222`)
- `NATS_MONITOR_URL` — NATS HTTP monitoring endpoint (default: `http://localhost:8222`)

## Standard Stream Names

| Stream     | Subjects              | Purpose                                 |
| ---------- | --------------------- | --------------------------------------- |
| `SYSTEM`   | `system.>`, `nexus.>` | System alerts, health, heartbeats       |
| `DIGITAL`  | `digital.>`           | Agent tasks and results                 |
| `PHYSICAL` | `physical.>`          | Robot commands, sensor data, HA events  |
| `COMMS`    | `comms.>`             | Human messages across channels          |
| `EVENTS`   | `events.>`            | Community, task, and integration events |

---

_Complements nexus-bridge. Legacy's window into the nervous system._


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
