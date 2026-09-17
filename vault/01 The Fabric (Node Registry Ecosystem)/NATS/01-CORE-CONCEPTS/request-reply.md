---
title: "Request / Reply"
section: "01-CORE-CONCEPTS"
type: note
tags: [nats, request-reply, rpc, inbox, no-responders, scatter-gather, timeout]
source: https://docs.nats.io/nats-concepts/core-nats/reqreply
created: 2026-06-14
owner: Legacy
status: active
---

# Request / Reply

> **Core idea.** NATS implements RPC-style request/reply directly on top of pub/sub. The requester publishes to a service subject and provides a unique **reply-inbox** subject; the responder listens on the service subject and publishes its answer back to that inbox. No dedicated middleware, sidecars, or proxies required.

See [[NATS/01-CORE-CONCEPTS/publish-subscribe|Publish/Subscribe]] for the underlying model.  
See [[NATS/01-CORE-CONCEPTS/queue-groups|Queue Groups]] for scaling responders.

---

## 1. The Pattern

```
Requester                        NATS Server                        Responder
    │                                 │                                  │
    │  PUB time                       │                                  │
    │  reply-to: _INBOX.abc123        │  SUB time ────────────────────── │
    │────────────────────────────────►│                                  │
    │                                 │──── MSG time _INBOX.abc123 ─────►│
    │                                 │                                  │ (processes)
    │  SUB _INBOX.abc123              │  PUB _INBOX.abc123 "12:34:56"   │
    │◄─────────────────────────────── │◄────────────────────────────────│
    │  MSG _INBOX.abc123 "12:34:56"   │                                  │
    │◄────────────────────────────────│                                  │
```

**Steps:**
1. Requester generates a unique **inbox subject** (e.g., `_INBOX.<nuid>`).
2. Requester subscribes to that inbox with `AutoUnsubscribe(1)` (receive exactly one reply).
3. Requester publishes to the service subject with `reply-to` set to the inbox.
4. Responder receives the message, inspects `msg.Reply`, publishes the answer there.
5. Requester receives the reply on its inbox; subscription auto-closes.

---

## 2. Inbox Subjects (`_INBOX`)

- The inbox prefix is `_INBOX.` by default (a reserved NATS prefix).
- Each inbox is a globally unique, randomly generated subject string using **NUID** (NATS UID — a 22-character base62 string).
- Client libraries generate inboxes automatically via helper functions:
  - Go: `nats.NewInbox()` → `_INBOX.abc123xyz789...`
  - Python: `nc.new_inbox()` → `_INBOX.abc123xyz789...`
  - JS: `createInbox()`

> **Security note:** Because inboxes are random, they are unguessable. In multi-tenant environments, the `_INBOX` prefix can be scoped per account or overridden via client options (`InboxPrefix`).

---

## 3. The High-Level `Request()` API

Most client libraries provide a `Request()` convenience function that creates the inbox, subscribes, publishes, and waits — all in one call.

### Go — `nc.Request()`

```go
nc, err := nats.Connect("demo.nats.io")
if err != nil {
    log.Fatal(err)
}
defer nc.Close()

// Sends to "time", waits up to 1 second for a single reply
msg, err := nc.Request("time", nil, time.Second)
if err != nil {
    // err == nats.ErrNoResponders  → no service listening (503)
    // err == nats.ErrTimeout       → service exists but took too long
    log.Fatal(err)
}

log.Printf("Reply: %s", msg.Data)
```

### Python — `nc.request()`

```python
import asyncio, nats

async def main():
    nc = await nats.connect("nats://demo.nats.io:4222")

    # Responder
    async def handle_time(msg):
        from datetime import datetime
        await nc.publish(msg.reply, datetime.now().isoformat().encode())

    await nc.subscribe("time", cb=handle_time)

    # Requester — blocks for up to 1 second
    try:
        reply = await nc.request("time", b"", timeout=1.0)
        print("Reply:", reply.data.decode())
    except asyncio.TimeoutError:
        print("No reply within timeout")

    await nc.close()

asyncio.run(main())
```

---

## 4. The Low-Level Pattern (Manual Inbox)

For more control (e.g., scatter-gather, custom timeout logic), use the underlying publish + subscribe approach directly.

### Go — manual inbox

```go
// Create a unique subject for replies
uniqueReplyTo := nats.NewInbox()

// Listen for a single response
sub, err := nc.SubscribeSync(uniqueReplyTo)
if err != nil {
    log.Fatal(err)
}

// Send the request
if err := nc.PublishRequest("time", uniqueReplyTo, nil); err != nil {
    log.Fatal(err)
}
nc.Flush()

// Read the reply (1 second timeout)
msg, err := sub.NextMsg(time.Second)
if err != nil {
    log.Fatal(err)
}
log.Printf("Reply: %s", msg.Data)
sub.Unsubscribe()
```

### Python — manual inbox

```python
future = asyncio.Future()

async def cb(msg):
    future.set_result(msg)

await nc.subscribe("time", cb=cb)

unique_reply_to = nc.new_inbox()
await nc.publish_request("time", unique_reply_to, b"")
await nc.flush()

msg = await asyncio.wait_for(future, timeout=1.0)
# Respond
time_bytes = datetime.now().isoformat().encode()
await nc.publish(msg.reply, time_bytes)
```

---

## 5. Responder Side

A responder is just a normal subscriber that reads `msg.Reply` and publishes to it:

```go
// Responder — subscribe to service subject and respond
nc.Subscribe("time", func(m *nats.Msg) {
    timeBytes := []byte(time.Now().String())
    m.Respond(timeBytes)   // convenience: publishes to m.Reply
})
```

```python
async def responder(msg):
    await msg.respond(datetime.now().isoformat().encode())

await nc.subscribe("time", cb=responder)
```

> If `msg.Reply` is empty, the message was published without a reply-to — the responder should handle this gracefully (either ignore or log).

---

## 6. No-Responders (Instant 503)

**Problem:** Without no-responders, a request to a subject with no subscribers blocks until the client timeout expires — potentially seconds of wasted latency.

**Solution:** NATS clients that opt into `no_responders` receive an **immediate synthetic reply** containing a `Status: 503` header when the server determines there are no subscribers on the target subject.

```
Requester                      NATS Server
    │                               │
    │  PUB time                     │  (no subscribers on "time")
    │  reply-to: _INBOX.abc         │
    │──────────────────────────────►│
    │                               │
    │◄── Status: 503 No Responders ─│  ← instant, no timeout wait
```

**Go — handling no-responders:**

```go
msg, err := nc.Request("time", nil, time.Second)
if err != nil {
    if errors.Is(err, nats.ErrNoResponders) {
        // Service is down or not deployed — react immediately
        log.Println("No service available")
    } else {
        log.Fatal(err)
    }
}
```

| Scenario | Error returned |
|----------|----------------|
| No subscribers on subject | `nats.ErrNoResponders` (instant) |
| Subscriber exists but doesn't reply in time | `nats.ErrTimeout` (after timeout) |
| Subscriber replies with error status | reply message with `Status` header |

> Requires server ≥ v2.2.0 and a client that sends `no_responders: true` in the CONNECT handshake (all modern official clients do this by default).

---

## 7. Timeouts

The `Request()` call always takes a timeout parameter. Best practices:

| Scenario | Recommended timeout |
|----------|---------------------|
| Internal microservice | 100–500 ms |
| Remote API / DB call | 1–5 s |
| Heavy computation | 10–30 s (consider async pattern instead) |
| Interactive user-facing | < 200 ms for good UX |

Set timeouts **defensively short** — no-responders will return instantly anyway, and the timeout only matters for services that are slow rather than absent.

---

## 8. Scatter-Gather (Multiple Responses)

NATS request/reply isn't limited to one responder. If multiple services subscribe to the same subject, you can **collect multiple responses** — useful for:

- Fan-out queries (ask all nodes, aggregate results)
- Quorum reads
- First-N-reply patterns (take the fastest K of N)

```go
// Scatter-gather: collect responses for up to 500ms
replyTo := nats.NewInbox()
sub, err := nc.SubscribeSync(replyTo)
if err != nil {
    log.Fatal(err)
}
nc.Flush()

nc.PublishRequest("service.query", replyTo, []byte("input"))

var responses []string
maxWait := 500 * time.Millisecond
deadline := time.Now().Add(maxWait)

for time.Now().Before(deadline) {
    msg, err := sub.NextMsg(time.Until(deadline))
    if err != nil {
        break   // timeout reached
    }
    responses = append(responses, string(msg.Data))
}
sub.Unsubscribe()

log.Printf("Got %d responses: %v", len(responses), responses)
```

> **First-responder wins:** The higher-level `nc.Request()` only receives the first reply and discards the rest. Use the manual inbox pattern above to collect multiple.

### Queue Groups + Scatter-Gather

When responders use a [[NATS/01-CORE-CONCEPTS/queue-groups|queue group]], each queue group returns exactly **one** response (load-balanced internally). Multiple queue groups on the same subject each return one response — useful for multi-service fan-out.

---

## 9. Automatic Routing & Location Transparency

The reply-to inbox is just a subject. NATS routes it automatically:

- Requester can be on server A, responder on server B — works transparently.
- **No service registry needed** — the subject is the service address.
- Multiple responders can subscribe to the same subject; NATS distributes via queue groups if configured.
- When a service scales down gracefully (drains), in-flight requests complete before the connection closes.

---

## 10. Request/Reply vs Direct Publish

| | `Request()` | Plain `Publish()` |
|-|-------------|-------------------|
| Waits for response | ✅ Yes | ❌ No |
| Timeout | ✅ Required | N/A |
| Reply subject | Auto-generated inbox | Not set |
| No-responders signal | ✅ ErrNoResponders | ❌ Silent discard |
| Use case | RPC / query | Notifications / events |

---

## Sources

- https://docs.nats.io/nats-concepts/core-nats/reqreply (primary)
- https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/replyto.md
- https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/sending/request_reply.md
- https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/receiving/reply.md
- https://github.com/nats-io/nats.docs/blob/master/_examples/request_reply.html
