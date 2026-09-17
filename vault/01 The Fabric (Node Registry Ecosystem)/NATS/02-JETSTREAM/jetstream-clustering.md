---
title: "JetStream Clustering & High Availability"
section: "02-JETSTREAM"
type: note
tags: [nats, jetstream, clustering, raft, replicas, ha, meta-group]
source: https://docs.nats.io/running-a-nats-service/configuration/clustering/jetstream_clustering
created: 2026-06-14
---

# JetStream Clustering & High Availability

How JetStream replicates data and survives node failure: **RAFT** consensus over replica groups.

See also: [[NATS/06-TOPOLOGY-SCALING/clustering|Clustering]] · [[NATS/02-JETSTREAM/streams|Streams]]

---

## Replicas — R1 / R3 / R5

Every stream (and consumer) has a **replica count** (`R`). JetStream maintains `R` copies across distinct cluster nodes using RAFT.

| Replicas | Tolerates | Use |
|----------|-----------|-----|
| **R1** | 0 node failures (single copy) | Dev, non-critical, single-node, or where loss is acceptable |
| **R3** | **1** node failure | **Production default** — quorum of 2 of 3 |
| **R5** | **2** node failures | Critical data, larger clusters |

Replica count must be **odd** (1, 3, 5) so RAFT can form a majority quorum and cannot deadlock on ties. Max is **5**.

```shell
nats stream add ORDERS --replicas 3
```

---

## RAFT — how replication works

Each replicated stream is a **RAFT group**: one **leader** + followers. Writes go to the leader, which replicates to followers; a write is acknowledged once a **quorum** (majority) has persisted it. If the leader fails, the remaining nodes **elect a new leader** automatically (RAFT election) and service continues — this is the self-healing HA.

```
        ┌─────────── RAFT group (R3) ───────────┐
write ─▶ │  LEADER  ──replicate──▶ FOLLOWER       │
        │     │                    FOLLOWER       │
        │     └──ack once quorum (2/3) persisted──┘
        └────────────────────────────────────────┘
Leader dies ▶ followers elect a new leader ▶ no data loss (quorum had it)
```

NATS embeds a RAFT implementation (the `graft`/internal RAFT layer) — no external consensus service (contrast Kafka's historical Zookeeper dependency).

---

## The Meta-Group (cluster controller)

Above individual stream groups sits the **meta-group**: a cluster-wide RAFT group formed by all JetStream-enabled servers. It has its own **meta-leader** and manages:
- Stream/consumer **creation, deletion, placement**
- **Asset assignment** — deciding which nodes host which replica groups
- Cluster membership for JetStream

`nats server report jetstream` and the `$SYS`/`$JS.API` surfaces show the meta-leader and per-asset leaders.

---

## Placement & Tags

You can steer where replicas land using **placement**:
- **`cluster`** — pin a stream's replicas to a named cluster (in a supercluster).
- **`tags`** — match server tags (e.g. `tag:ssd`, `tag:region:us-east`) so data lands on appropriate hardware/region.

```shell
nats stream add ORDERS --replicas 3 --tag region:us-east
```

This is how you achieve **geo-aware placement** in a supercluster — keep a stream's quorum within one region while mirrors serve others.

---

## Failure Behavior

| Event | Result |
|-------|--------|
| 1 follower down (R3) | No impact — quorum (2/3) intact, writes continue |
| Leader down (R3) | Brief election (~sub-second to seconds), new leader, **no data loss** (quorum had the writes) |
| 2 nodes down (R3) | **Loss of quorum** — stream becomes unavailable for writes (read may continue on survivor) until quorum restored. Use R5 to tolerate 2 failures |
| Node rejoins | Catches up from the leader automatically |

> **Quorum math:** an R-replica group needs `floor(R/2)+1` healthy nodes. R3 needs 2, R5 needs 3. Lose quorum → the group pauses to protect consistency (CP-leaning under partition).

---

## Sizing & Best Practices

- **3 or 5 servers** for a JetStream cluster (odd, for meta-group quorum).
- Use **R3** for production streams; **R5** only for the most critical, on a 5+ node cluster.
- Keep **`max_file_store`/`max_memory_store`** consistent across nodes — replicas need room.
- Distribute nodes across **failure domains** (racks/AZs) so one failure domain ≠ quorum loss.
- For multi-region, use a **supercluster** with **placement tags** to keep quorums regional and **mirrors** for cross-region reads (avoid stretching a single RAFT quorum across high-latency links).

---

## OneMind Relevance

The current OneMind bus is **single-server** (omos-mbp hub) — effectively R1, fine for a homestead hub where the MBP is always-on. The growth path: as the fleet scales to "hundreds of nodes," a small JetStream **cluster** (3 always-on boxes: MBP + omos-cloud + one more) would give R3 HA for critical streams, while **leaf-node edge domains** stay R1 locally and **source** up to the clustered hub. See `~/onemindos-core/docs/NATS-DOCTRINE.md`.

## Sources
- JetStream clustering — https://docs.nats.io/running-a-nats-service/configuration/clustering/jetstream_clustering
- JetStream administration — https://docs.nats.io/running-a-nats-service/nats_admin/jetstream_admin
- Clustering — https://docs.nats.io/running-a-nats-service/configuration/clustering
- RAFT (nats-io/graft) — https://github.com/nats-io/graft
