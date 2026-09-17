---
title: "Queue Groups"
section: "01-CORE-CONCEPTS"
type: note
tags: [nats, queue-groups, load-balancing, competing-consumers, horizontal-scaling, jetstream]
source: https://docs.nats.io/nats-concepts/core-nats/queue.md
created: 2026-06-14
owner: Legacy
status: active
---

# Queue Groups

> **Core idea.** A queue group is a named group of subscribers sharing the same subject. When a message arrives, the NATS server delivers it to **exactly one** randomly chosen member of the group — providing built-in, zero-configuration load balancing. This is the NATS answer to competing consumers.

See [[NATS/01-CORE-CONCEPTS/publish-subscribe|Pub/Sub]] for the standard 1→N fan-out model.  
For durable, acknowledged work queues, see [[NATS/02-JETSTREAM/README|JetStream]] (WorkQueuePolicy).

---

## 1. How Queue Groups Work

```
Publisher                   NATS Server              Queue Group "workers"
    │                            │                         │
    │  PUB updates "task-1"      │   QueueSub updates ─────► worker-A ✓
    │───────────────────────────►│   QueueSub updates ─────  worker-B ✗ (not chosen)
    │                            │   QueueSub updates ─────  worker-C ✗ (not chosen)
    │
    │  PUB updates "task-2"      │   QueueSub updates ─────  worker-A ✗
    │───────────────────────────►│   QueueSub updates ─────► worker-B ✓
    │                            │   QueueSub updates ─────  worker-C ✗
```

**Key behaviors:**
- The server randomly selects **one** subscriber per queue group per message.
- Non-queue subscribers on the same subject **still receive every message** (fan-out still occurs for non-queue subs).
- Queue group membership is entirely **defined by the application** — no server-side configuration required.
- Scaling is immediate: add a process → it joins the group; stop it → it leaves.

### Mixed subscriptions on the same subject

```
PUB orders "order-99"
  → regular SUB orders     ← receives the message (fan-out)
  → regular SUB orders     ← receives the message (fan-out)
  → QueueSub orders "team" ← exactly ONE of these receives it
  → QueueSub orders "team"
  → QueueSub orders "team"
```

Regular (non-queue) subscribers and queue subscribers coexist on the same subject independently.

---

## 2. Queue Group Naming Rules

Queue group names follow the same rules as subjects:

| Rule | Detail |
|------|--------|
| Case-sensitive | `workers` ≠ `Workers` |
| No whitespace | Spaces are not allowed |
| Dot separator | Can use `.` for hierarchy: `payments.workers` |
| Wildcard matching | Server features like queue permissions can wildcard on the name |

---

## 3. Advantages

| Advantage | Description |
|-----------|-------------|
| **Fault tolerance** | If one instance crashes, others in the group keep processing |
| **Elastic scaling** | Scale up: start more instances. Scale down: stop them. No config changes. |
| **No duplicate delivery** | Only one instance receives each message |
| **Zero configuration** | No broker queues to declare; just subscribe with a queue group name |
| **Self-organizing** | Members join and leave the group dynamically |

---

## 4. Code Examples

### Go — `QueueSubscribe`

```go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

wg := sync.WaitGroup{}
wg.Add(10)

// Any number of processes can call this; each joins the "workers" group
if _, err := nc.QueueSubscribe("updates", "workers", func(m *nats.Msg) {
    log.Printf("Worker received: %s", m.Data)
    wg.Done()
}); err != nil {
    log.Fatal(err)
}

wg.Wait()
```

### Python — `subscribe(queue=...)`

```python
import asyncio, nats

async def main():
    nc = await nats.connect("nats://demo.nats.io:4222")

    future = asyncio.Future()

    async def worker(msg):
        print(f"Worker got: {msg.data.decode()}")
        future.set_result(msg)

    # Join the "workers" queue group on "updates"
    await nc.subscribe("updates", queue="workers", cb=worker)

    await nc.publish("updates", b"All is Well")

    msg = await asyncio.wait_for(future, 1)
    print("Received:", msg)
    await nc.close()

asyncio.run(main())
```

### C# — `queueGroup` parameter

```csharp
using NATS.Net;

await using var client = new NatsClient();
var count = 0;

await foreach (var msg in client.SubscribeAsync<string>(
    subject: "updates", queueGroup: "workers"))
{
    Console.WriteLine($"Received {++count}: {msg.Data}");
    if (count == 10) break;
}
```

---

## 5. Queue Groups in the CLI (Walkthrough)

Start three service instances in a queue group:

```bash
# Terminal 1
nats reply foo "service instance A Reply# {{Count}}"

# Terminal 2
nats reply foo "service instance B Reply# {{Count}}"

# Terminal 3
nats reply foo "service instance C Reply# {{Count}}"
```

Send 10 requests — each goes to one instance:

```bash
nats request foo --count 10 "Request {{Count}}"
```

The 10 replies will be distributed roughly evenly across the three instances.

---

## 6. Queue Groups + Request/Reply

Queue groups are the natural scaling mechanism for [[NATS/01-CORE-CONCEPTS/request-reply|Request/Reply]] services. Instead of a single responder, run N responders all subscribing with the same queue group:

```go
// Each instance of this service joins the "time-service" queue group
nc.QueueSubscribe("time", "time-service", func(m *nats.Msg) {
    m.Respond([]byte(time.Now().String()))
})
```

```python
# Python service — scale by running multiple processes
async def handle_time(msg):
    await msg.respond(datetime.now().isoformat().encode())

await nc.subscribe("time", queue="time-service", cb=handle_time)
```

The **no-responders** feature integrates seamlessly: if all instances of the queue group are down, the server returns `Status: 503` immediately to the requester. See [[NATS/01-CORE-CONCEPTS/request-reply#6-no-responders-instant-503|No Responders]].

---

## 7. Geo-Affinity

When connecting to a globally distributed NATS **supercluster**, queue groups exhibit automatic **service geo-affinity**:

- A request is first routed to a queue group member **in the same cluster** (same region/datacenter).
- Only if no local members are available does NATS route the request to another cluster (another region).
- This minimizes cross-region latency **automatically**, with zero configuration.

```
Cluster US-WEST                     Cluster EU-CENTRAL
┌──────────────────────────┐        ┌──────────────────────────┐
│  service QS "workers"    │        │  service QS "workers"    │
│  (instance A, B, C)      │        │  (instance D, E)         │
└──────────────────────────┘        └──────────────────────────┘
       ▲                                     ▲
       │ request from US client              │ request from EU client
       │ → served locally                    │ → served locally
```

---

## 8. Queue Permissions (Security)

The NATS server can restrict which queue group names a user may subscribe to:

```hcl
users = [
  {
    user: "worker-a", password: "secret", permissions: {
      sub: {
        allow: ["updates workers"]     # subject + queue group name
      }
    }
  },
  {
    user: "worker-b", password: "secret", permissions: {
      sub: {
        allow: ["updates", "updates v1", "updates v1.>", "updates *.dev"],
        deny: ["> *.prod"]             # deny prod queue groups
      }
    }
  }
]
```

Queue permission format: `"<subject> <queue-group>"`. Wildcards work on both parts.

---

## 9. Queue Groups + JetStream

### JetStream WorkQueuePolicy

For **persistent, acknowledged** work queues (at-least-once delivery), use [[NATS/02-JETSTREAM/README|JetStream]] with `RetentionPolicy: WorkQueuePolicy`:

```go
// Create a stream with work queue retention
js.AddStream(&nats.StreamConfig{
    Name:      "WORK",
    Subjects:  []string{"work.>"},
    Retention: nats.WorkQueuePolicy,
})

// Pull consumer — multiple workers compete for messages
sub, _ := js.PullSubscribe("work.>", "my-worker")
msgs, _ := sub.Fetch(10)
for _, msg := range msgs {
    // Process...
    msg.Ack()
}
```

### JetStream Push Consumer + Queue Group

Alternatively, use a **push consumer** with explicit ack and a queue group of subscribers for fan-out with acknowledgement:

```java
// Java JetStream push subscriber with durable consumer and queue group
PushSubscribeOptions pso = PushSubscribeOptions.builder()
    .durable("qdur-durable")
    .build();

JetStreamSubscription sub = js.subscribe(
    "qdur-subject",
    "qdur-queue",       // queue group name
    pso
);
```

| Approach | Delivery | Ack required | Replay on restart |
|----------|----------|--------------|-------------------|
| Core NATS queue group | At-most-once | No | No |
| JetStream pull + WorkQueue | At-least-once | Yes | Yes |
| JetStream push + queue group | At-least-once | Yes | Yes (durable) |

---

## 10. Drain Before Scale-Down

When scaling down (removing an instance from the queue group), use **drain** rather than an abrupt close. Drain:
1. Stops accepting new messages.
2. Processes all buffered/in-flight messages.
3. Then closes the connection.

```go
// Graceful scale-down
nc.Drain()   // blocks until all pending messages are processed
```

```python
await nc.drain()   # processes buffered messages then closes
```

This prevents message loss during rolling deployments and zero-downtime restarts.

---

## Sources

- https://docs.nats.io/nats-concepts/core-nats/queue.md (primary)
- https://docs.nats.io/nats-concepts/core-nats/queue/queues_walkthrough.md
- https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/queues.md
- https://github.com/nats-io/nats.docs/blob/master/running-a-nats-service/configuration/securing_nats/authorization.md
- https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/consumers.md
