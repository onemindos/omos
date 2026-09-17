---
title: "Benchmarking & Performance"
section: "11-ADVANCED"
type: note
tags: [nats, benchmarking, performance, nats-bench, latency, throughput, tuning, slow-consumer]
source: https://docs.nats.io/using-nats/nats-tools/nats_cli/natsbench
created: 2026-06-14
---

# Benchmarking & Performance

NATS is fast by default — millions of Core msgs/sec, sub-100µs latency — but *your* numbers depend on payload size, fan-out, persistence, and network. The `nats bench` tool measures it directly; this note covers how to run it and which levers move the result.

See also: [[NATS/09-CLI-TOOLING/nats-cli|nats CLI]] · [[NATS/07-DEPLOYMENT-OPS/tuning-and-dr|Tuning & DR]] · [[NATS/10-PROTOCOL-INTERNALS/server-architecture|Server Architecture]]

---

## `nats bench` — the built-in benchmark

`nats bench` is part of the `nats` CLI. It spins up publishers/subscribers and reports throughput and latency.

```bash
# Core pub/sub: 5 publishers, 5 subscribers, 1M msgs, 16-byte payload
nats bench bench.test --pub 5 --sub 5 --msgs 1000000 --size 16

# Request-reply latency
nats bench svc.echo --pub 1 --sub 1 --request --msgs 100000

# JetStream throughput (persisted) — note the difference vs Core
nats bench bench.js --js --pub 5 --msgs 1000000 --size 256 \
  --stream BENCH --storage file

# Pull consumer throughput
nats bench bench.js --js --sub 5 --pull --msgs 1000000
```

Output reports **msgs/sec**, **MB/sec**, and latency distribution. Always benchmark against a server that resembles production (same storage, TLS on/off, same network).

> Older builds expose `nats bench <subject> --pub/--sub`; newer builds add subcommands like `nats bench pub|sub|service|jetstream`. Run `nats bench --help` to match your CLI version.

---

## The fundamental trade-off: throughput vs durability

| Mode | Semantics | Relative speed | Why |
|------|-----------|---------------|-----|
| **Core pub/sub** | at-most-once, in-memory | ⚡ fastest (millions/sec) | no disk, no ack, interest-only |
| **JetStream, memory storage** | at-least-once | fast | replication/ack overhead, no disk |
| **JetStream, file storage R1** | at-least-once, durable | moderate | fsync to disk |
| **JetStream, file R3/R5** | durable + replicated | slower | RAFT quorum across nodes |

The first rule of NATS performance: **don't pay for durability you don't need.** Telemetry that can drop a sample → Core. Orders that must never vanish → JetStream R3. Mixing them on the same bus is the whole point.

---

## Levers that move the numbers

**Payload size.** Small messages are latency/op-bound; large messages are bandwidth-bound. `max_payload` caps single messages (default 1 MB) — big blobs belong in the Object Store ([[NATS/03-KV-OBJECT-STORE/object-store|Object Store]]), chunked.

**Batching (JetStream).** Pull consumers fetching **batches** (`Fetch(500)`) massively out-throughput one-at-a-time. Async publish with a window of in-flight acks beats synchronous publish-and-wait.

**Fan-out.** Core fan-out is cheap (interest graph), but N subscribers means N deliveries — bandwidth scales with subscriber count. Queue groups load-balance instead of duplicating.

**Connections vs subscriptions.** Many subjects on one connection is fine (sublist is O(tokens)). Thousands of connections cost FDs/goroutines — size `max_connections` and the OS.

**TLS.** Adds handshake + per-message crypto cost. Real, but usually worth it; benchmark with it **on** if production uses it.

**Storage backend.** SSD vs spinning disk vs tmpfs dramatically changes JetStream file-storage throughput. RAFT R3 adds network round-trips per commit.

---

## Slow consumers — the #1 real-world gotcha

In **Core** NATS, if a subscriber can't keep up, the server's per-connection pending buffer fills and the server declares a **slow consumer**, dropping messages (and possibly the connection) to protect itself. This is by design — Core is at-most-once and won't let one laggard back-pressure the bus.

Fixes:
- Move the workload to **JetStream** with a durable **pull** consumer — flow-controlled, the consumer pulls at its own pace, nothing is dropped.
- Use **queue groups** to spread load across N workers.
- Do less in the message callback — hand off to a worker pool; never block the read loop.
- Raise client pending limits only as a band-aid, not a fix.

See the error semantics in [[NATS/10-PROTOCOL-INTERNALS/wire-protocol|Wire Protocol]] (`-ERR 'Slow Consumer'`).

---

## A sane benchmarking method
1. **Baseline Core** pub/sub at your real payload size — the ceiling.
2. **Add JetStream** (the persistence cost) — memory first, then file, then R3.
3. **Turn on TLS** if production uses it.
4. **Match the network** — localhost lies; benchmark across the real link (Tailscale, LAN, WAN).
5. **Watch the server** — `/varz`, `/jsz`, `nats-top` during the run to spot slow consumers, memory, CPU. See [[NATS/07-DEPLOYMENT-OPS/monitoring|Monitoring]].

---

## OneMind Relevance
- **Tiered by need:** fleet telemetry/heartbeats ride **Core** (drop-tolerant, maximal speed); commands, jobs, and ledgers ride **JetStream R3** (durable). Benchmark each tier at its real payload before committing subject law.
- Edge links (Tailscale/WAN) are the real bottleneck, not the broker — always bench across the actual link, and lean on **leaf + local JetStream** so edge throughput isn't hostage to the WAN.
- Watch for slow consumers on any agent doing heavy per-message work — move it to a pull consumer + worker pool.

## Sources
- nats bench — https://docs.nats.io/using-nats/nats-tools/nats_cli/natsbench
- Performance / tuning — https://docs.nats.io/running-a-nats-service/configuration/tuning
- Slow consumers — https://docs.nats.io/using-nats/developer/connecting/events/slow
- latency-tests — https://github.com/nats-io/latency-tests
