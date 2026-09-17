---
title: "Stream Mirrors & Sources"
section: "02-JETSTREAM"
type: note
tags: [nats, jetstream, mirrors, sources, replication, aggregation]
source: https://docs.nats.io/nats-concepts/jetstream/source_and_mirror
created: 2026-06-14
---

# Stream Mirrors & Sources

JetStream can build streams **from other streams**. Two mechanisms — **mirrors** (1:1 copy) and **sources** (N:1 aggregation) — enable replication, aggregation, geo-distribution, and read scaling.

See also: [[NATS/02-JETSTREAM/streams|Streams]] · [[NATS/06-TOPOLOGY-SCALING/README|Topology & Scaling]] · [[NATS/06-TOPOLOGY-SCALING/domains|JetStream Domains]]

---

## Mirror — exact 1:1 copy

A **mirror** stream continuously copies messages from **one** origin stream, preserving sequence numbers. A mirror **cannot** have its own `Subjects` (it doesn't ingest from subjects directly — only from its origin).

```
Origin stream ORDERS  ──mirror──▶  ORDERS_MIRROR (read replica, same seqs)
```

Use cases:
- **Read scaling** — point read-heavy consumers at the mirror, sparing the origin.
- **Geo read-locality** — a mirror in another region/cluster serves local readers.
- **Backup / DR** — a mirror in a separate domain is a live standby copy.
- **Cross-domain bridging** — pull a hub stream down into a leaf-node domain.

```shell
nats stream add ORDERS_MIRROR --mirror ORDERS
```

---

## Source — N:1 aggregation

A **source** stream pulls from **one or more** origin streams (and can *also* have its own subjects), interleaving their messages into a single combined stream. Sequence numbers are **re-assigned** (the combined stream has its own sequence space).

```
ORDERS_US   ─┐
ORDERS_EU   ─┼─sources──▶  ORDERS_GLOBAL (aggregated)
ORDERS_ASIA ─┘
```

Use cases:
- **Aggregation** — combine per-region/per-edge streams into one global view.
- **Fan-in from the edge** — many leaf-node streams roll up to a hub stream.
- **Consolidated analytics** — one stream to feed a downstream sink (e.g. ClickHouse).

```shell
nats stream add ORDERS_GLOBAL \
  --source ORDERS_US --source ORDERS_EU --source ORDERS_ASIA
```

---

## Mirror vs Source — the decision

| | **Mirror** | **Source** |
|---|---|---|
| Origins | Exactly **one** | **One or many** |
| Own `Subjects`? | **No** | **Yes** (can mix sourced + direct) |
| Sequence numbers | **Preserved** from origin | **Re-assigned** in the new stream |
| Mental model | "Read replica / copy" | "Aggregator / combiner" |
| Direction typical | Hub → edge (pull down) or DR | Edge → hub (roll up) |

---

## Filtering & Transformation

Both mirrors and sources support:
- **Filter subjects** — copy only a subset of the origin's messages.
- **Subject transforms** — rewrite subjects as messages flow (e.g. `orders.*` → `archive.orders.*`).
- **Start position** — begin from a sequence or time (`OptStartSeq`, `OptStartTime`), so you can backfill or tail.
- **External origin (API prefix / domain)** — source/mirror from a stream in **another account or JetStream domain** via `External` config (domain API prefix + delivery prefix). This is the cross-domain replication primitive for leaf nodes.

---

## Cross-Domain / Cross-Account Replication

When the origin lives in a different **JetStream domain** (e.g. a hub) or a different **account**, the mirror/source declares an `External` reference: the **API prefix** of the origin's domain and an optional **delivery prefix**. This is how an edge leaf node mirrors a hub stream locally (read it offline) or sources its local data up to the hub. See [[NATS/06-TOPOLOGY-SCALING/domains|JetStream Domains]] and [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]].

```
Hub domain (omos-mbp)            Edge domain (barn-pi)
   stream TEL  ◀───── source (External: api_prefix=$JS.hub.API) ───  stream TEL_LOCAL
```

---

## OneMind Relevance

This is the backbone of the planned fleet topology: edge nodes (`omd --embed`) run isolated **JetStream domains**; **sources** roll their telemetry up to the hub, and **mirrors** push hub config/streams down to the edge for offline read. The hub (`omos-mbp`) stays the aggregation point. Combine with leaf-node store-and-forward for disconnect-tolerant geo-distribution.

## Sources
- Sources & Mirrors — https://docs.nats.io/nats-concepts/jetstream/source_and_mirror
- Replication / configuration — https://docs.nats.io/nats-concepts/jetstream/streams
- Cross-domain (External) — https://docs.nats.io/nats-concepts/jetstream/js_walkthrough
- Leaf nodes + JetStream — https://docs.nats.io/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes
