---
title: "Tuning, Backup & Disaster Recovery"
section: "07-DEPLOYMENT-OPS"
type: note
tags: [nats, tuning, performance, backup, disaster-recovery, upgrades, signals]
source: https://docs.nats.io/running-a-nats-service/nats_admin
created: 2026-06-14
---

# Tuning, Backup & Disaster Recovery

Performance tuning, backups, upgrades, and the operational signals you control NATS with.

See also: [[NATS/07-DEPLOYMENT-OPS/configuration-reference|Configuration Reference]] · [[NATS/07-DEPLOYMENT-OPS/monitoring|Monitoring]] · [[NATS/02-JETSTREAM/jetstream-clustering|JetStream Clustering]]

---

## Performance Tuning

| Lever | Guidance |
|-------|----------|
| **`max_payload`** | Keep near default (~1MB) for low latency; big data → Object Store |
| **`max_pending`** | Raise per-connection outbound buffer for bursty subscribers (costs memory); or fix the slow consumer |
| **`write_deadline`** | Increase on slow networks so legit-but-slow clients aren't dropped |
| **File store disk** | **SSD/NVMe** for JetStream `store_dir`; disk I/O is the JetStream bottleneck |
| **OS file descriptors** | Raise `ulimit -n` / `LimitNOFILE` (e.g. 65536) for many connections |
| **GOMAXPROCS / CPU** | NATS is efficient; ensure enough cores for TLS + JetStream at scale |
| **Batching (clients)** | Pull-consumer `fetch(N)` batches and publish batching cut round-trips |
| **`MaxAckPending`** | Tune per-consumer in-flight window: higher = throughput, lower = back-pressure |
| **Replicas** | R3 costs ~3× write amplification vs R1 — only replicate what needs HA |

Benchmark with `nats bench` (see [[NATS/11-ADVANCED/README|Advanced]] / [[NATS/09-CLI-TOOLING/README|CLI]]).

---

## Signals (operational control)

`nats-server --signal <cmd>` (or `kill -<SIG>`):

| Signal | Effect |
|--------|--------|
| **reload** (SIGHUP) | Reload config in place (limits, logging, TLS certs, config users/perms) |
| **stop** (SIGINT/SIGTERM) | Graceful shutdown |
| **ldm** (SIGUSR2) | **Lame Duck Mode** — stop accepting new connections, gently drain existing, then exit. Use before maintenance so clients migrate to peers without a hard drop |
| **reopen** (SIGUSR1) | Reopen log file (log rotation) |

```shell
nats-server --signal ldm      # graceful drain before taking a node down
nats-server --signal reload   # apply config changes
```

> **Lame Duck Mode** is the right way to remove a cluster node: it signals clients to reconnect elsewhere *before* the server exits, avoiding a thundering reconnect storm.

> **Restart-required changes (verified pitfalls):** `system_account`, account topology, and JetStream `domain` do **not** apply on SIGHUP — they need a full restart.

---

## Backup & Restore (JetStream)

```shell
# Back up a single stream (config + data)
nats stream backup ORDERS ./orders-backup

# Restore it (to same or another server/cluster)
nats stream restore ./orders-backup

# Account-level / full JetStream snapshot
nats account backup --all ./js-backup
```

- Backups capture stream **config + messages + consumer state**.
- Restore to a **different cluster** for migration or DR rehearsal.
- For continuous DR, prefer **mirrors** to a second cluster/domain (live standby) over periodic snapshots — see [[NATS/02-JETSTREAM/mirrors-and-sources|Mirrors & Sources]].
- File-store data also lives under `store_dir` — snapshot the volume for a crash-consistent copy, but `nats stream backup` is the supported logical backup.

---

## Disaster Recovery Strategies

| Strategy | Mechanism | RPO/RTO |
|----------|-----------|---------|
| **Periodic backup** | `nats stream backup` to object storage | RPO = backup interval |
| **Live mirror (cross-cluster)** | Mirror stream in a second cluster/domain | Near-zero RPO, fast failover |
| **R3/R5 in-cluster** | RAFT replicas survive node loss | Zero RPO for node failure (not site failure) |
| **Supercluster + placement** | Regional quorums + cross-region mirrors | Site-failure tolerant |

Layer them: **R3** for node failure, **cross-cluster mirror** for site failure, **periodic backup** for "oops I deleted the stream."

---

## Upgrades

- **Rolling upgrade** a cluster one node at a time; R3 JetStream tolerates one node down. Use **Lame Duck Mode** on each node before replacing it.
- NATS maintains **client/server compatibility** across minor versions — upgrade servers and clients independently.
- Read the **release notes** for the target version (`nats-io/nats-server/releases`); current line **v2.14.2** (2026-06-02).
- Validate config (`nats-server -c conf -t`) before/after.

---

## OneMind Relevance

Single-server hub today → backups via `nats stream backup` of critical streams (AGENT_EVENTS, telemetry) to omos-cloud/object storage are the immediate DR win. As HA arrives (3-node cluster), R3 covers node failure; a **mirror to omos-cloud** would cover homestead-site failure. Use **LDM** whenever cycling the launchd service to avoid client reconnect storms. Restart (not reload) after any account/domain change.

## Sources
- Admin / signals — https://docs.nats.io/running-a-nats-service/nats_admin/signals
- Lame Duck Mode — https://docs.nats.io/running-a-nats-service/nats_admin/lame_duck_mode
- Backup/restore — https://docs.nats.io/running-a-nats-service/nats_admin/jetstream_admin/disaster_recovery
- Slow consumers / tuning — https://docs.nats.io/running-a-nats-service/configuration
