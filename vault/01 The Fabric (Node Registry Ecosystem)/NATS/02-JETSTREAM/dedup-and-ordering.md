---
title: "Deduplication & Message Ordering"
section: "02-JETSTREAM"
type: note
tags: [nats, jetstream, deduplication, ordering, exactly-once, msg-id]
source: https://docs.nats.io/nats-concepts/jetstream/streams
created: 2026-06-14
---

# Deduplication & Message Ordering

How JetStream achieves **exactly-once publish** and what ordering guarantees you actually get.

See also: [[NATS/02-JETSTREAM/streams|Streams]] · [[NATS/02-JETSTREAM/consumers|Consumers]]

---

## The Problem: Delivery Guarantees

- **NATS Core** = at-most-once (fire-and-forget).
- **JetStream** = at-least-once by default (a publisher retries on no-ack; a consumer redelivers on no-ack). At-least-once means **duplicates are possible** — a publish ack can be lost in transit, the publisher retries, and the same logical message lands twice.

JetStream closes this on **both ends**:
1. **Publish-side dedup** → exactly-once *write* into the stream.
2. **Consume-side ack tracking** → exactly-once *effective processing* (with idempotent handlers).

---

## Publish Deduplication (`Nats-Msg-Id`)

When a publisher sets the **`Nats-Msg-Id`** header on a message, JetStream tracks that ID within a sliding **duplicate window** (`Duplicates`, default **2 minutes**). If a second message arrives with the **same `Nats-Msg-Id`** inside the window, the server **discards the duplicate** and still returns a success ack (`duplicate: true`). The publisher can safely retry without creating duplicates.

```python
# Python — idempotent publish
ack = await js.publish("ORDERS.new", b"...", headers={"Nats-Msg-Id": "order-12345"})
# Retrying the exact same publish within the dedup window is a no-op write.
```

```shell
# CLI shows the window
nats stream info ORDERS | grep -i duplicate
# Duplicate tracking time window: 2m0s
```

| Knob | Meaning |
|------|---------|
| `Nats-Msg-Id` (header) | Caller-supplied unique ID for the message |
| `Duplicates` (stream config) | Length of the sliding dedup window (time) |

**Choosing the window:** it must be ≥ the longest time your publisher might retry. Too short → a late retry slips through as a duplicate. Longer windows cost more memory (the server tracks IDs for the window). 2 minutes is the default; tune to your retry policy.

> The dedup window is **per-stream** and lives in memory — it does not survive arbitrarily long. It guarantees exactly-once for *retries within the window*, not global dedup across all time.

---

## Exactly-Once Semantics (end to end)

JetStream's "exactly-once" is the combination of:

1. **Exactly-once publish** — `Nats-Msg-Id` + dedup window → the message is written once even if the publisher retries.
2. **Exactly-once consume** — the consumer uses **explicit ack**; on success it acks, advancing the ack floor. Redelivery only happens if the ack is missing. Make your handler **idempotent** (or use the message sequence as a dedup key downstream) so a redelivery before an ack lands is harmless.

> "Exactly-once" in any distributed system = at-least-once delivery + idempotent processing + dedup. JetStream gives you the primitives (`Nats-Msg-Id`, sequence numbers, explicit ack); the application supplies idempotency.

---

## Expected Last Sequence / Subject (optimistic concurrency)

Publishers can assert **preconditions** via headers, and the server rejects the publish if they don't hold — enabling optimistic concurrency control:

| Header | Asserts |
|--------|---------|
| `Nats-Expected-Stream` | Message is going to the stream you think it is |
| `Nats-Expected-Last-Sequence` | The stream's last sequence == N (no one else wrote since) |
| `Nats-Expected-Last-Subject-Sequence` | The last seq **for this subject** == N |
| `Nats-Expected-Last-Msg-Id` | The last message's `Nats-Msg-Id` == X |

These power the **compare-and-set** behavior under JetStream **KV** (which is built on a stream): a KV `update` is a publish with `Nats-Expected-Last-Subject-Sequence` so concurrent writers can't clobber each other.

---

## Ordering Guarantees

| Scope | Guarantee |
|-------|-----------|
| **Within a stream** | Messages get **monotonically increasing sequence numbers**; storage order is preserved. A single consumer reading with `AckPolicy=all` or sequential explicit acks sees them **in order**. |
| **Per subject** | Order within a subject is preserved (per-subject sequence tracked). |
| **Across multiple consumers** | If many workers share a pull consumer (load balancing), **per-message** processing parallelism breaks strict global order — each message is still delivered once, but worker A and worker B run concurrently. Use a single-threaded consumer or `MaxAckPending=1` if you need strict serial processing. |
| **Across streams** | No cross-stream ordering. Order is a per-stream property. |

> For strict in-order processing: one durable pull consumer, one worker, `AckPolicy=explicit`, `MaxAckPending=1`. For throughput: many workers on one consumer, accepting concurrent (still exactly-once) processing.

---

## OneMind Relevance

The KV buckets (`zeus-context`, `agent-state`, etc.) rely on this machinery: each KV bucket is a stream `KV_<bucket>`, writes are publishes, `update` uses expected-last-sequence for atomic CAS, and history is just the stream's retained revisions. The dedup/ordering rules above are *why* KV watch + history + atomic updates work.

## Sources
- Streams / deduplication — https://docs.nats.io/nats-concepts/jetstream/streams#message-deduplication
- Publish acks & headers — https://docs.nats.io/using-nats/developer/develop_jetstream/model_deep_dive
- KV (built on dedup/expected-seq) — https://docs.nats.io/nats-concepts/jetstream/key-value-store
- Exactly-once — https://docs.nats.io/using-nats/developer/develop_jetstream
