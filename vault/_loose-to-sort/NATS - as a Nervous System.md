# NATS — as a Nervous System

Think of your biological nervous system. Your brain doesn't call your hand via HTTP and wait for a 200 response. It fires a signal on a nerve pathway, and anything connected to that pathway reacts. That's NATS.

---

## Pub/Sub — The Sensory Nerves

```
Publisher (sensor)  →  Subject (nerve pathway)  →  Subscriber(s) (reaction)
```

**What it is:** You publish a message to a "subject" (just a string like `agent.legacy.inbox`). Anyone subscribed to that subject gets it instantly. One publisher, unlimited subscribers. Fire and forget.

**In OneMind terms:**
```
Legacy publishes:     "system.health.check"
    → Zenith hears it (subscribed to system.>)
    → Sovereignty hears it (subscribed to system.>)
    → Your dashboard hears it (subscribed to system.health.>)
    → Your phone notification service hears it
```

One publish. Everyone who cares gets it. Nobody polls.

**Why it matters:** Today OpenClaw agents communicate via webhooks — one sender, one receiver, you manually wire each connection. With pub/sub, an agent publishes once and the system routes it to everyone who needs it. Add a new consumer? Just subscribe. Zero config changes anywhere else.

**In OneMind terms:**
```python
# You (or any agent) ask:
reply = await nats.request("agents.prompt.hermes.zeus.default", "summarize my day")
# Hermes receives it, processes, streams reply back
# You don't import Hermes, don't know its IP, don't know its port
```

**Why it matters:** This is how Synadia Agents works. Discovery (`$SRV.INFO.agents`) is a request/reply — you ask "who's online?" and every agent replies. Prompting is request/reply with streaming. Your agents stop being hardcoded endpoints and become discoverable services.

---

## Queue Groups — Load Balancing (Built-in)

```
Message arrives → Subject → Only ONE member of the queue group gets it
```

**What it is:** Multiple instances of the same service subscribe to the same subject but form a "queue group." Each message goes to exactly one member. Round-robin distribution, zero config.

**In OneMind terms:**
```
You run 2 instances of Legacy (for high availability):
  Legacy-A subscribes to "agents.prompt.oc.zeus.legacy" (queue: "agents")
  Legacy-B subscribes to "agents.prompt.oc.zeus.legacy" (queue: "agents")

Prompt arrives → Only ONE Legacy gets it (no duplicates)
Legacy-A crashes → Legacy-B gets ALL prompts (auto-failover)
```

**Why it matters:** Scale out any agent without changing anything. The Synadia Agent Protocol uses this by default — two instances claiming the same identity automatically load-balance. This is how you get agent high-availability for free.

---

## JetStream — Memory (Persistence)

Messages aren't just fire-and-forget anymore — they're STORED in ordered streams.

**What it is:** A persistent log. Messages published to a subject are captured into a "stream" and kept (with configurable retention: time, count, or size). Consumers can replay from any point. Guaranteed delivery (at-least-once or exactly-once).

**In OneMind terms:**
```
Stream: AGENT_EVENTS
├── agent.legacy.completed.task-123  (stored 10:00am)
├── agent.zenith.workout.logged      (stored 10:05am)
├── agent.sovereignty.deploy.done    (stored 10:12am)
└── system.health.alert.disk-full    (stored 10:30am)
```

You were asleep? Doesn't matter.
Wake up → Consumer replays everything you missed since last seen.

**Why it matters:** Without persistence, pub/sub is volatile — miss a message and it's gone. JetStream gives you an event log of everything that happened in your system. Your daily briefing agent can scan the last 24h of events. Your phone can catch up after being offline. Nothing is lost.

### Stream Map

| Stream | Captures |
|--------|----------|
| EVENTS | System lifecycle events |
| SYSTEM | Health, alerts, heartbeats |
| NEXUS | Cross-agent coordination |
| DIGITAL | App/service signals |
| PHYSICAL | IoT, sensors, location |
| COMMS | Messages, notifications |

---

## Key/Value Store — Shared State

Like Redis, but inside NATS. No extra service.

**What it is:** A distributed key-value store with:
- **Watchers** — get notified the instant a key changes
- **Versioning** — every write creates a new revision (rollback possible)
- **TTL** — keys auto-expire after a time

**In OneMind terms:**
```python
# Agent state visible to ALL agents instantly:
kv = js.key_value("agent-state")

await kv.put("legacy.status", "busy")          # Legacy updates its status
await kv.put("legacy.current_task", "writing daily briefing")
await kv.put("zeus.location", "home-office")
await kv.put("zeus.focus_mode", "deep-work")

# Any agent watches for changes:
watcher = await kv.watch("zeus.>")
async for entry in watcher:
    # Fires INSTANTLY when Zeus's state changes
    print(f"{entry.key} → {entry.value}")
```

**Why it matters:** Today your agents read AMBIENT_STATE.md on heartbeat intervals. With KV, state changes propagate instantly to all watchers. Legacy knows the moment Zeus enters focus mode. Zenith knows the moment a workout starts. No polling. No file reads. Real-time shared consciousness.

---

## Object Store — Large Assets

Like S3, but inside NATS. Versioned. No extra service.

**What it is:** Store large blobs (files, images, PDFs) — anything over the message size limit (~1MB). Versioned, chunked, streamable.

**Why it matters:** Agents can share generated artifacts without touching the filesystem. A report generated by Eternal can be instantly available to Legacy for delivery to Telegram. No shared folders, no git sync delays.

---

## The Complete Picture: One Binary Replaces 5+ Services

| What NATS Replaces | Old Way | NATS Way |
|--------------------|---------|----------|
| Inter-agent messaging | Webhooks, HTTP calls | Pub/Sub subjects |
| Agent discovery | Hardcoded URLs, config files | `$SRV.INFO.agents` (one round-trip) |
| Shared state | Files (AMBIENT_STATE.md), Redis | KV Store with watchers |
| Event history | SQLite ingester, manual logs | JetStream streams |
| File sharing between agents | Git sync, shared filesystem | Object Store |
| Load balancing | Reverse proxy, manual routing | Queue Groups (built-in) |
| Health monitoring | Custom heartbeat scripts | NATS micro service heartbeats |

All of this in ONE 20MB binary. No Redis. No Kafka. No RabbitMQ. No separate DB for events. No API gateway. No service mesh.
