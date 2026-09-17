---
title: "JetStream Consumers"
section: "02-JETSTREAM"
type: note
tags: [nats, jetstream, consumers, pull, push, ack, durable]
source: https://docs.nats.io/nats-concepts/jetstream/consumers
created: 2026-06-14
---

# JetStream Consumers

> A **consumer** is a stateful cursor into a stream. The stream stores; the consumer tracks *what a given application has read and acknowledged*. Many independent consumers can read the same stream at different positions without interfering.

See also: [[NATS/02-JETSTREAM/streams|Streams]] · [[NATS/02-JETSTREAM/retention-and-limits|Retention & Limits]] · [[NATS/02-JETSTREAM/dedup-and-ordering|Dedup & Ordering]]

---

## Push vs Pull

| | **Pull consumer** (preferred) | **Push consumer** |
|---|---|---|
| Who drives delivery | **Client fetches** batches when ready | **Server pushes** to a delivery subject |
| Flow control | Built-in — client asks for N messages | Requires explicit flow control / rate limits |
| Scaling workers | Many clients bind to **one** pull consumer → automatic load balancing | Needs a queue/deliver-group |
| Overwhelm risk | Low — consumer pulls only what it can handle | Higher — server can outrun a slow client |
| Modern recommendation | **Use pull** for almost everything | Legacy / specific streaming-to-subject needs |

```shell
# Pull consumer
nats consumer add ORDERS dispatch --pull --ack explicit --filter "ORDERS.new"
nats consumer next ORDERS dispatch --count 10     # fetch a batch
```

```python
# Python pull
js = nc.jetstream()
sub = await js.pull_subscribe("ORDERS.new", durable="dispatch")
msgs = await sub.fetch(10, timeout=5)
for m in msgs:
    await m.ack()
```

> **OneMind pitfall (verified):** in `nats-py` 2.x do **not** pass both `queue=` and `durable=` to `js.subscribe()` — the queue name *is* treated as the durable name; passing both errors.

---

## Durable vs Ephemeral

- **Durable** — has a name, **persists** its position across client restarts and disconnects. The server remembers the last ack'd sequence; a restarted worker resumes exactly where it left off. Use for workers/services.
- **Ephemeral** — unnamed, **auto-deleted** shortly after the client goes away (governed by an inactivity threshold). Use for one-off reads, debugging, temporary tailing.

---

## Acknowledgement Policies

The ack policy decides how delivery is confirmed:

| AckPolicy | Meaning |
|-----------|---------|
| **`explicit`** | Every message must be individually ack'd (most common, required for work queues) |
| **`all`** | Ack'ing message N implicitly acks all prior (≤N) — cheaper, ordered processing |
| **`none`** | No acks; fire-hose delivery, no redelivery (e.g. pure tailing/metrics) |

### Ack types (client → server)
| Ack | Effect |
|-----|--------|
| **`ACK`** (`+ACK`) | Success — advance, don't redeliver |
| **`NAK`** (`-NAK`) | Negative — redeliver (optionally after a delay) |
| **`TERM`** (`+TERM`) | Stop — never redeliver this message (poison message) |
| **`WPI`/`+WPI` (in progress)** | "Still working" — reset the AckWait timer for long jobs |

```python
await m.ack()          # success
await m.nak(delay=5)   # retry in 5s
await m.term()         # give up, don't redeliver
await m.in_progress()  # extend deadline
```

---

## Key Consumer Config

| Field | Purpose |
|-------|---------|
| **DeliverPolicy** | Where to start: `all` · `last` · `new` · `by_start_sequence` · `by_start_time` · `last_per_subject` |
| **AckPolicy** | `explicit` · `all` · `none` |
| **AckWait** | How long to wait for an ack before redelivering (default 30s) |
| **MaxDeliver** | Max redelivery attempts before giving up |
| **FilterSubject(s)** | Restrict the consumer to a subset of the stream's subjects |
| **ReplayPolicy** | `instant` (as fast as possible) · `original` (preserve original inter-message timing) |
| **MaxAckPending** | Max un-ack'd messages in flight (back-pressure throttle) |
| **Backoff** | Array of redelivery delays (escalating retry, e.g. `[1s,5s,30s]`) |
| **RateLimit** | Bits/sec cap (push consumers) |
| **SampleFrequency** | % of acks to sample for monitoring |
| **InactiveThreshold** | Auto-delete ephemeral after this idle time |
| **MaxWaiting** | Max outstanding pull requests (pull consumers) |
| **HeadersOnly** | Deliver only headers, not payloads (lightweight) |

---

## DeliverPolicy — where consumption starts

| Policy | Starts at |
|--------|-----------|
| `all` | The very first message in the stream (full replay) |
| `last` | The last message currently in the stream |
| `new` | Only messages arriving *after* the consumer is created |
| `by_start_sequence` | A specific sequence number |
| `by_start_time` | A specific timestamp |
| `last_per_subject` | The last message for **each** matching subject (great for KV-like "current state") |

---

## Redelivery & Poison Messages

If a message isn't ack'd within **AckWait**, JetStream **redelivers** it (up to **MaxDeliver**). Each redelivery increments a delivery count visible in metadata. Use **Backoff** for escalating retry delays. A message that keeps failing is a **poison message** — cap it with `MaxDeliver`, then `TERM` it or route to a dead-letter stream (via an advisory on `$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES.>`).

---

## Consumer Acknowledgement Advisory & Monitoring

JetStream publishes **advisories** (`$JS.EVENT.ADVISORY.>`) for max-deliveries, terminated messages, and more — subscribe to build dead-letter handling and alerting. `nats consumer info STREAM CONSUMER` shows live state: ack floor, pending, redelivered, num waiting.

## Sources
- Consumers — https://docs.nats.io/nats-concepts/jetstream/consumers
- Consumer config — https://docs.nats.io/nats-concepts/jetstream/consumers#configuration
- Acknowledgements — https://docs.nats.io/using-nats/developer/develop_jetstream/consumers
- Developing with JetStream — https://docs.nats.io/using-nats/developer/develop_jetstream
