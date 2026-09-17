---
title: "JetStream Retention & Limits"
section: "02-JETSTREAM"
type: note
tags: [nats, jetstream, retention, limits, workqueue, interest, discard]
source: https://docs.nats.io/nats-concepts/jetstream/streams
created: 2026-06-14
---

# JetStream Retention & Limits

How long a stream keeps messages, and what happens when it fills. Retention policy + discard policy + limits together define the stream's storage behavior.

See also: [[NATS/02-JETSTREAM/streams|Streams]] · [[NATS/02-JETSTREAM/consumers|Consumers]]

---

## The Three Retention Policies

| Policy | A message is removed when… | Classic use |
|--------|----------------------------|-------------|
| **Limits** (default) | A limit is hit (age, count, or bytes). Messages stay until then, regardless of consumption. | Event logs, audit trails, replayable history, time-series |
| **Interest** | **All** consumers that had interest have ack'd it (no interested consumer still needs it). | Fan-out to a known set of consumers; drop once everyone's done |
| **WorkQueue** | **Any** consumer ack's it — each message is delivered to exactly one consumer then deleted. | Task/job queues, work distribution |

### Limits (default)
Messages live until a configured limit forces eviction. This is the "it's a log" model — consumers can replay from the beginning, late subscribers still see history. Combine with `MaxAge`/`MaxMsgs`/`MaxBytes` to bound storage.

### Interest
The stream only retains a message while **at least one consumer still has interest and hasn't ack'd it**. Once every interested consumer acks, it's gone. Storage tracks consumer interest — there's no long-term history. Good for "deliver to the current set of services, then forget."

> **Subtlety:** with Interest retention, a message published when **no consumer exists** is dropped immediately — interest is evaluated at publish time. Consumers must exist *before* publishing to receive.

### WorkQueue
Each message is consumed **once** by one consumer, then deleted on ack. The stream acts as a distributed job queue. **Constraint:** consumer filter subjects must be **non-overlapping** — a given subject can be served by only one consumer, so messages aren't ambiguously claimed.

---

## Discard Policy — what happens at the limit

When a **Limits** stream hits `MaxMsgs` or `MaxBytes`:

| DiscardPolicy | Behavior |
|---------------|----------|
| **`old`** (default) | Evict the **oldest** messages to make room for new ones (ring-buffer feel) |
| **`new`** | **Refuse** new messages (publish errors) until space frees up |

`DiscardNewPerSubject` extends `new` to apply per-subject when `MaxMsgsPerSubject` is set.

---

## The Limit Knobs

| Limit | Meaning |
|-------|---------|
| **MaxAge** | Max age of any message (TTL for the whole stream). *In nats-py, set in **seconds**; in the raw stream config / protocol it's nanoseconds.* |
| **MaxMsgs** | Max total messages (`-1` = unlimited) |
| **MaxBytes** | Max total bytes stored (`-1` = unlimited) |
| **MaxMsgsPerSubject** | Max messages retained **per subject** (e.g. `1` ⇒ KV-like "latest only") |
| **MaxMsgSize** | Largest single message accepted (payload + headers) |
| **MaxConsumers** | Cap on consumers for the stream (`-1` = unlimited) |

> **OneMind pitfall (verified):** in `nats-py`, `StreamConfig.max_age` takes **seconds as int**, not nanoseconds. `max_age=30*24*3600` = 30 days. Multiplying by `1e9` expires the stream almost instantly.

---

## Per-Message TTL

Modern JetStream supports **per-message TTL** (the `Nats-TTL` header / stream `AllowMsgTTL`) so individual messages can expire independently of the stream's `MaxAge`. Useful for mixing short-lived and long-lived data in one stream.

---

## How They Combine

A stream applies **all** active limits simultaneously — the **first** limit reached triggers the discard policy. Example: a Limits stream with `MaxAge=7d`, `MaxBytes=5G`, `DiscardOld` keeps at most 7 days *and* at most 5 GB, evicting oldest when either is hit. This is exactly the OneMind `GNMI_TELEMETRY` pattern (5 GB max, 7-day retention).

---

## Choosing a Policy (decision guide)

```
Need replayable history / late subscribers see past?  ──▶ Limits
Each message handled by exactly one worker, then gone? ──▶ WorkQueue
Deliver to current consumers, drop once all ack'd?     ──▶ Interest
```

OneMind streams are mostly **Limits** (AGENT_EVENTS 30-day, EVT/TEL with byte+age caps) for replay/audit, with task queues as candidates for **WorkQueue**.

## Sources
- Streams (retention, limits, discard) — https://docs.nats.io/nats-concepts/jetstream/streams
- Retention policies — https://docs.nats.io/nats-concepts/jetstream/streams#retentionpolicy
- Stream limits — https://docs.nats.io/nats-concepts/jetstream/streams#stream-limits-and-message-retention
- Per-message TTL — https://docs.nats.io/nats-concepts/jetstream/streams#message-ttl
