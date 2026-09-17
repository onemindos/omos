---
title: "Subjects & Wildcards"
section: "01-CORE-CONCEPTS"
type: note
tags: [nats, subjects, wildcards, addressing, hierarchy, tokens]
source: https://docs.nats.io/nats-concepts/subjects
created: 2026-06-14
owner: Legacy
status: active
---

# Subjects & Wildcards

> **Core idea.** Every NATS message is addressed by a **subject** — a dot-separated string of tokens. Publishers always use exact subjects. Subscribers can use **wildcards** to receive messages across a range of subjects without subscribing to each one individually.

---

## 1. What is a Subject?

A NATS subject is a plain UTF-8 string that serves as the **channel name** for a message. NATS is an *interest-based* messaging system: the server only routes a message to clients that have expressed interest (via a subscription) in that subject or a wildcard matching it.

```
time.us.east.atlanta
│    │  │    └── token 4
│    │  └──────── token 3
│    └─────────── token 2
└──────────────── token 1
```

> Subjects in other systems are called **topics** (Kafka, MQTT), **channels** (Redis Pub/Sub), or **routing keys** (RabbitMQ). In NATS, "stream" refers specifically to [[NATS/02-JETSTREAM/README|JetStream]] message storage, not a subject.

---

## 2. Subject Rules (Valid vs. Invalid)

### Syntax Rules

| Rule | Detail |
|------|--------|
| Token separator | `.` (ASCII dot, U+002E) — the **only** structural separator |
| Spaces | ❌ **Not allowed** anywhere in a subject |
| Wildcards in published subjects | ❌ **Not allowed** — publishers must use fully qualified subjects |
| Wildcards in subscriptions | ✅ Allowed |
| Case sensitivity | ✅ Case-sensitive (`Foo` ≠ `foo`) |
| Empty tokens | ❌ `foo..bar` is invalid (double-dot = empty token between) |
| Leading/trailing dot | ❌ `.foo` or `foo.` are invalid |
| `$` prefix | Reserved for NATS system subjects (e.g., `$SRV`, `$JS`, `_INBOX`) |

### Valid Examples

```
foo
foo.bar
time.us.east
orders.v2.store123.order171711
_INBOX.abc123xyz789
$JS.API.STREAM.CREATE
```

### Invalid Examples

```
foo bar          ← space not allowed
foo..bar         ← empty token
.foo             ← leading dot
foo.             ← trailing dot
time.*.east      ← wildcard in publish subject (not valid for PUB)
```

### Size Recommendations

| Limit | Recommended |
|-------|-------------|
| Max tokens | ≤ 16 |
| Max total length | < 256 characters |
| Max subscribed subjects per server | ~1 million (beyond that, server RAM grows linearly; 1M+ subjects > 1 GB RAM) |

---

## 3. Subject Hierarchy & Token Design

The dot separator creates a **hierarchical namespace** that makes wildcards useful. Think of tokens as progressively narrowing scopes from left to right.

### Anatomy of a Good Hierarchy

```
<domain>.<entity-type>.<identifier>.<action>
orders   .online      .store123    .created
sensors  .temperature .rack-7      .reading
```

### World Clock Example (from docs.nats.io)

```
time.us
time.us.east
time.us.east.atlanta
time.eu.east
time.eu.east.warsaw
```

A subscription on `time.us.>` receives all US times; `time.*.east` receives all eastern-region times regardless of continent.

### Design Guidelines

| Do | Avoid |
|----|-------|
| Use first tokens for broad namespace | Over-encoding technical details |
| Use last tokens for specific identifiers | Putting too much data in the subject |
| Use headers for metadata | Encoding payload data in the subject |
| Keep hierarchy ≤ 4–5 levels deep | Deeply nested (>8 tokens) unless necessary |
| Name business/physical entities | Generic names that require doc lookup |

**Good:**
```
orders.online.store123.order171711
factory1.tools.group42.unit17
```

**Over-engineered (avoid):**
```
orders.online.us.server42.ccpayment.premium.store123.electronics.deliver-dhl.order171711.create
```

---

## 4. Wildcards

> **Rule:** Publishers send to **exact subjects only**. Wildcards are **subscription-only** constructs.

### 4.1 `*` — Single Token Wildcard

Matches **exactly one** token at that position.

| Subscription | Matches | Does NOT match |
|--------------|---------|----------------|
| `time.*.east` | `time.us.east`, `time.eu.east` | `time.us.east.atlanta` (4 tokens vs 3) |
| `*.bar` | `foo.bar`, `baz.bar` | `a.b.bar` |
| `time.us.*` | `time.us.east`, `time.us.central` | `time.us.east.atlanta` |

`*` can appear **multiple times** and can be combined with `>`:

```
*.*.east.>   matches   time.us.east.atlanta
             matches   sensor.rack.east.temp.reading
```

### 4.2 `>` — Multi-Token Tail Wildcard

Matches **one or more** tokens. **Must appear at the end** of the subject.

| Subscription | Matches | Does NOT match |
|--------------|---------|----------------|
| `time.us.>` | `time.us.east`, `time.us.east.atlanta`, `time.us.west.la` | `time.eu.east` |
| `>` | **everything** on the server | — |
| `time.>` | all time subjects | `sensor.temp` |

> **Wire-tap / monitoring:** Subscribing to `>` receives all messages on the server (subject to security configuration). Useful for debugging and observability.

### 4.3 `*` vs `>` Comparison

| Property | `*` | `>` |
|----------|-----|-----|
| Tokens matched | Exactly 1 | 1 or more |
| Position | Any position | End only |
| Can appear multiple times | Yes | No (only terminal) |
| Example | `time.*.east` | `time.us.>` |

---

## 5. Wildcards in Code

### Go — subscribe to `*` wildcard

```go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

wg := sync.WaitGroup{}
wg.Add(2)

// Matches time.us.east and time.eu.east (but NOT time.us.east.atlanta)
if _, err := nc.Subscribe("time.*.east", func(m *nats.Msg) {
    log.Printf("subject=%s data=%s", m.Subject, m.Data)
    wg.Done()
}); err != nil {
    log.Fatal(err)
}

wg.Wait()
```

### Go — subscribe to `>` wildcard

```go
// Matches time.A.east, time.B.east, time.C.west, time.D.west, etc.
if _, err := nc.Subscribe("time.>", func(m *nats.Msg) {
    log.Printf("%s: %s", m.Subject, m.Data)
    wg.Done()
}); err != nil {
    log.Fatal(err)
}
```

### Python — subscribe to `>` wildcard

```python
import asyncio
import nats

async def main():
    nc = await nats.connect("nats://demo.nats.io:4222")

    async def cb(msg):
        print(f"subject={msg.subject} data={msg.data.decode()}")

    await nc.subscribe("time.>", cb=cb)

    # Publish test messages
    await nc.publish("time.A.east", b'A')
    await nc.publish("time.B.east", b'B')
    await nc.publish("time.C.west", b'C')

    await asyncio.sleep(1)
    await nc.close()

asyncio.run(main())
```

### Python — route messages by inspecting `msg.subject` inside wildcard subscription

```python
async def cb(msg):
    match msg.subject:
        case "time.us.east":
            tz = "America/New_York"
        case "time.us.central":
            tz = "America/Chicago"
        case _:
            tz = "UTC"
    print(f"{msg.subject} → {tz}")

await nc.subscribe("time.us.*", cb=cb)
```

---

## 6. Subject-Based Addressing & Location Transparency

NATS uses subject-based addressing to provide **location transparency** across a NATS server cloud:

- Subject subscriptions are **automatically propagated** across all servers in a cluster or supercluster.
- A message published to `orders.created` on server A is automatically delivered to a subscriber on server B if B has a matching subscription — no manual routing configuration needed.
- Messages with **no subscribers are silently discarded** (Core NATS). Use [[NATS/02-JETSTREAM/README|JetStream]] if you need persistence for offline consumers.

### Subject-Based Security & Filtering

A well-designed hierarchy enables fine-grained controls in:

- **Authorization** — allow/deny publish/subscribe per user or account
- **Account exports/imports** — expose subjects across multi-tenant accounts
- **Subject mapping** — rewrite/remap subjects for versioning or canary routing
- **JetStream stream ingestion** — filter which subjects land in a stream
- **Leaf node bridging** — control which subjects flow across edge connections

---

## 7. System Subjects (Reserved Prefixes)

| Prefix | Used By | Example |
|--------|---------|---------|
| `$JS.` | JetStream API | `$JS.API.STREAM.CREATE` |
| `$SRV.` | Micro/Services framework | `$SRV.INFO` |
| `$SYS.` | Server system events | `$SYS.ACCOUNT.*.CONNECT` |
| `_INBOX.` | Request/reply inboxes | `_INBOX.abc123` |
| `$G.` | Default account internal | (internal) |

> Avoid using `$` or `_INBOX` prefixes for application subjects.

---

## Sources

- https://docs.nats.io/nats-concepts/subjects (primary)
- https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/wildcards.md
- https://github.com/nats-io/nats.docs/blob/master/nats-concepts/subject_mapping.md
- https://github.com/nats-io/nats.docs/blob/master/reference/nats-protocol/nats-protocol/README.md
