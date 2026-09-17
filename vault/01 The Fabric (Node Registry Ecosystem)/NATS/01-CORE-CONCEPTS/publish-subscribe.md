---
title: "Publish / Subscribe"
section: "01-CORE-CONCEPTS"
type: note
tags: [nats, pubsub, publish, subscribe, fan-out, at-most-once, async, sync, unsubscribe]
source: https://docs.nats.io/nats-concepts/core-nats/pubsub
created: 2026-06-14
owner: Legacy
status: active
---

# Publish / Subscribe

> **Core idea.** NATS pub/sub is the fundamental messaging pattern: a publisher sends a message to a subject and **every** subscriber with a matching subscription receives it simultaneously. This is a **1-to-N fan-out** with **at-most-once** delivery semantics — the server makes no attempt to store, retry, or redeliver.

See [[NATS/01-CORE-CONCEPTS/subjects-and-wildcards|Subjects & Wildcards]] for subject naming rules.  
For at-least-once / exactly-once delivery, see [[NATS/02-JETSTREAM/README|JetStream]].

---

## 1. The Core Model

```
Publisher                    NATS Server                     Subscribers
    │                            │                                │
    │  PUB updates "All is Well" │                                │
    │───────────────────────────►│   SUB updates  ────────────────► sub-1 ✓
    │                            │   SUB updates  ────────────────► sub-2 ✓
    │                            │   SUB updates.* ───────────────► sub-3 ✓  (wildcard)
    │                            │   SUB orders   ────────────────  sub-4 ✗  (different subject)
```

Key properties:
- The publisher does **not know** how many subscribers exist (or if any exist at all).
- If no subscribers are present, the message is **silently discarded** — Core NATS is not a queue.
- Each active matching subscriber receives its own copy of the message.
- Subject routing is performed by the **server** — clients only publish and subscribe.

---

## 2. Delivery Semantics: At-Most-Once

NATS Core pub/sub is **at-most-once** (also called "fire-and-forget"):

| Property | Behavior |
|----------|----------|
| Durability | None — message lives only in-flight |
| Retry on failure | None — publisher gets no acknowledgement |
| Offline subscribers | Miss all messages published while disconnected |
| Persistence | None — use [[NATS/02-JETSTREAM/README\|JetStream]] |
| Ordering | Per-subject FIFO within a single publisher connection |
| Throughput | Up to ~15M+ msgs/sec (in-memory, single server) |

> **When at-most-once is fine:** telemetry, real-time dashboards, live sensor readings, notifications where staleness is acceptable, and any case where a missed message is harmless.

---

## 3. Messages

A NATS message has four fields:

| Field | Description |
|-------|-------------|
| **Subject** | Destination address (required) |
| **Payload** | Arbitrary byte array (can be empty) |
| **Headers** | Key-value metadata, NATS/1.0 format (optional) |
| **Reply-to** | Subject the publisher listens on for a response (optional) |

Max payload size is **1 MB by default**, configurable up to 64 MB via `max_payload` in the server config (though ≤ 8 MB is recommended). See [[NATS/01-CORE-CONCEPTS/messages-and-headers|Messages & Headers]] for full anatomy.

---

## 4. Subscriptions

### 4.1 Async Subscriptions (preferred)

Async subscriptions invoke a **callback** when a message arrives. The callback runs in a dispatcher goroutine (Go) or event loop (Python/JS), decoupled from the publishing path.

**Go — async subscribe:**

```go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

// Callback is invoked asynchronously for each matching message
sub, err := nc.Subscribe("updates", func(m *nats.Msg) {
    log.Printf("Received: %s", m.Data)
})
if err != nil {
    log.Fatal(err)
}

// Block until done
select {}
```

**Python — async subscribe:**

```python
import asyncio
import nats

async def main():
    nc = await nats.connect("nats://demo.nats.io:4222")

    async def message_handler(msg):
        print(f"Received on {msg.subject}: {msg.data.decode()}")

    # Async callback-based subscription
    sub = await nc.subscribe("updates", cb=message_handler)

    await nc.publish("updates", b"All is Well")
    await asyncio.sleep(1)
    await nc.drain()

asyncio.run(main())
```

### 4.2 Sync Subscriptions

A synchronous subscription buffers incoming messages; the application calls `NextMsg(timeout)` to retrieve them. Useful for test code or simple request/reply patterns.

**Go — sync subscribe:**

```go
sub, err := nc.SubscribeSync("updates")
if err != nil {
    log.Fatal(err)
}

// Block until a message arrives (or timeout)
msg, err := sub.NextMsg(2 * time.Second)
if err != nil {
    log.Fatal(err)
}
log.Printf("Got: %s", msg.Data)
```

### 4.3 Comparison: Async vs Sync

| | Async | Sync |
|-|-------|------|
| Mechanism | Callback / handler | `NextMsg()` polling |
| Use case | Production subscriptions | Tests, one-shot receive |
| Backpressure | Pending limits (message count + bytes) | Caller controls pacing |
| Goroutine | Dispatcher goroutine per subscription | Application goroutine |

---

## 5. Fan-Out (1→N)

Any number of independent subscribers can receive the same message simultaneously:

```python
# Three subscribers on the same subject — all receive every publish
await nc.subscribe("sensor.temp", cb=lambda m: log_to_db(m))
await nc.subscribe("sensor.temp", cb=lambda m: alert_if_hot(m))
await nc.subscribe("sensor.temp", cb=lambda m: stream_to_dashboard(m))

# This single publish reaches all three
await nc.publish("sensor.temp", b"42.7")
```

Compare with [[NATS/01-CORE-CONCEPTS/queue-groups|Queue Groups]] where only **one** of N subscribers receives each message.

---

## 6. Wildcard Subscriptions in Pub/Sub

Subscribers can use wildcards to receive messages on multiple subjects:

```go
// Receive ALL time messages regardless of region/city
nc.Subscribe("time.>", func(m *nats.Msg) {
    fmt.Printf("%s → %s\n", m.Subject, m.Data)
})

// Receive only direct-east messages (one region level)
nc.Subscribe("time.*.east", func(m *nats.Msg) {
    fmt.Printf("East time from %s: %s\n", m.Subject, m.Data)
})
```

See [[NATS/01-CORE-CONCEPTS/subjects-and-wildcards|Subjects & Wildcards]] for the `*` and `>` rules.

---

## 7. Unsubscribe & Auto-Unsubscribe

### Manual Unsubscribe

```go
sub, _ := nc.Subscribe("updates", func(m *nats.Msg) {
    fmt.Println(m.Data)
})

// Later: stop receiving
sub.Unsubscribe()
```

```python
sub = await nc.subscribe("updates", cb=handler)

# Drain buffered messages, then unsubscribe
await sub.drain()
# or immediately:
await sub.unsubscribe()
```

### Auto-Unsubscribe (after N messages)

A subscription can be set to **automatically unsubscribe** after receiving a specific number of messages. This is baked into the NATS protocol: the server enforces the limit and stops forwarding after the count is reached.

**Go:**
```go
// Unsubscribe after receiving exactly 5 messages
sub, _ := nc.Subscribe("updates", func(m *nats.Msg) {
    fmt.Println(m.Data)
})
sub.AutoUnsubscribe(5)
```

**Python (nats.py):**
```python
# Unsubscribe after receiving 5 messages
sub = await nc.subscribe("updates", cb=handler)
await sub.unsubscribe(limit=5)
```

**C# — via subscription options:**
```csharp
// dotnet add package NATS.Net
var opts = new NatsSubOpts { MaxMsgs = 10 };

await foreach (var msg in client.SubscribeAsync<string>("updates", opts: opts))
{
    Console.WriteLine($"Received: {msg.Data}");
}
// Loop exits automatically after 10 messages
```

> Auto-unsubscribe is used internally by the [[NATS/01-CORE-CONCEPTS/request-reply|Request/Reply]] pattern to receive exactly one reply.

---

## 8. Publishing

Publishing is a simple one-way send. No ack is returned by the server in Core NATS.

**Go — publish bytes:**
```go
if err := nc.Publish("updates", []byte("All is Well")); err != nil {
    log.Fatal(err)
}
```

**Go — publish with reply-to subject:**
```go
// Manually set a reply subject (basis of request/reply)
if err := nc.PublishRequest("time", "_INBOX.abc123", nil); err != nil {
    log.Fatal(err)
}
```

**Python:**
```python
await nc.publish("updates", b"All is Well")
```

The server confirms receipt with `+OK` only when verbose mode is enabled (off by default in modern clients). In normal mode, publish is **asynchronous and best-effort**.

---

## 9. Wire Protocol Recap

At the protocol level, pub/sub maps to two commands:

```
# Publisher → Server
PUB updates 11\r\n
All is Well\r\n

# Subscriber → Server (register interest)
SUB updates <sid>\r\n
+OK\r\n

# Server → Subscriber (message delivery)
MSG updates <sid> 11\r\n
All is Well\r\n
```

With headers (NATS/1.0), `PUB` becomes `HPUB` and `MSG` becomes `HMSG`. See [[NATS/01-CORE-CONCEPTS/messages-and-headers|Messages & Headers]] and [[NATS/10-PROTOCOL-INTERNALS/README|Protocol Internals]].

---

## 10. Common Patterns

| Pattern | Description | Section |
|---------|-------------|---------|
| Pure fan-out | All subscribers receive all messages | This page |
| Load-balanced consumers | Only one subscriber per message (queue group) | [[NATS/01-CORE-CONCEPTS/queue-groups\|Queue Groups]] |
| RPC / synchronous call | Publisher awaits a response | [[NATS/01-CORE-CONCEPTS/request-reply\|Request/Reply]] |
| Persistent delivery | Messages stored and replayed for offline consumers | [[NATS/02-JETSTREAM/README\|JetStream]] |

---

## Sources

- https://docs.nats.io/nats-concepts/core-nats/pubsub (primary)
- https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/wildcards.md
- https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/README.md
- https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/unsub_after.md
- https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md
