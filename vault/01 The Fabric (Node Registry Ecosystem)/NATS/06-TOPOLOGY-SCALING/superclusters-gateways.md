---
title: "Superclusters & Gateways"
section: "06-TOPOLOGY-SCALING"
type: note
tags: [nats, supercluster, gateways, geo, multi-region]
source: https://docs.nats.io/running-a-nats-service/configuration/gateways
created: 2026-06-14
---

# Superclusters & Gateways

A **supercluster** connects multiple **clusters** via **gateways**, forming a geo-distributed mesh that still behaves as one bus — but optimized so cross-region traffic only flows when there's real interest.

See also: [[NATS/06-TOPOLOGY-SCALING/clustering|Clustering]] · [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]]

---

## The Model

```
   Region US (cluster)            Region EU (cluster)
   ┌───────────────┐  gateway     ┌───────────────┐
   │ S1 — S2 — S3  │◀───────────▶ │ S4 — S5 — S6  │
   └───────────────┘              └───────────────┘
            ▲  gateway                    ▲
            └──────────── Region ASIA ────┘
                     ┌───────────────┐
                     │ S7 — S8 — S9  │
                     └───────────────┘
```

- Each **cluster** is a full mesh internally (low-latency, regional).
- **Gateways** connect clusters to each other — a **gateway mesh** between clusters.
- The result is one logical subject space spanning regions, but with **interest-aware** cross-region forwarding.

---

## Why Gateways (not just a giant cluster)?

A single stretched cluster would full-mesh every server across high-latency WAN links — expensive and fragile. A supercluster:
- Keeps **intra-region** traffic on fast local routes.
- Sends a message **across a gateway only when the remote cluster has a subscriber** (interest propagation). No interest → no cross-region bytes.
- Uses **optimistic sends with interest pruning** — gateways learn where interest exists and stop forwarding where it doesn't.
- Scales to **many** clusters without N² server-to-server links across the WAN (it's cluster-to-cluster, not server-to-server, across regions).

---

## Configuration

```text
# A server in the US cluster
gateway {
  name: "US"                       # this cluster's gateway name
  listen: "0.0.0.0:7222"
  gateways: [
    { name: "EU",   urls: ["nats://eu-host:7222"] }
    { name: "ASIA", urls: ["nats://asia-host:7222"] }
  ]
}
```

- **`gateway.name`** = the cluster's identity in the supercluster (all servers in a cluster share it).
- **`gateways`** lists the other clusters' gateway URLs.
- Default gateway port **7222**. Secure with `gateway { tls { } }`.

---

## Interest Propagation & Optimistic Mode

Gateways don't broadcast all subscriptions globally (that wouldn't scale). Instead they start **optimistic** — assume no interest — and build interest as subscriptions appear, sending **interest updates** across gateways. A publish crosses to a remote cluster only once that cluster has signaled interest in the subject. This keeps a global bus cheap: local stays local, global crosses only when needed.

---

## JetStream in a Supercluster

- Streams are placed in a **specific cluster** (RAFT quorum stays regional — you don't stretch a quorum across the WAN).
- Use **placement tags** (`--tag region:us-east`) to pin a stream's replicas to the right cluster.
- Use **mirrors/sources** with `External` to move data across regions for read-locality or aggregation (see [[NATS/02-JETSTREAM/mirrors-and-sources|Mirrors & Sources]]).

This is the correct geo pattern: **regional quorums + cross-region mirrors**, never one quorum spanning continents.

---

## When You Need a Supercluster

- Multiple **data centers / regions** that must share one bus.
- **Geo read-locality** — serve subscribers from their nearest cluster.
- **Disaster resilience** — a whole region can fail, others continue.
- **Global scale** beyond a single cluster's comfortable size.

For **edge devices** (vs whole clusters), use **leaf nodes** instead — lighter weight, store-and-forward, single-link.

---

## OneMind Relevance

A supercluster is the **far-future** shape (multiple homesteads / regions / a cloud presence each running a cluster, stitched by gateways). Near-term geo reach for OneMind is served by **leaf nodes** (barn, office, vehicles, robots phoning home to the hub) — lighter than standing up regional clusters. Keep gateways in mind as the scale-out beyond a single homestead cluster.

## Sources
- Gateways — https://docs.nats.io/running-a-nats-service/configuration/gateways
- Gateway config — https://docs.nats.io/running-a-nats-service/configuration/gateways/gateway
- Supercluster JetStream — https://docs.nats.io/running-a-nats-service/configuration/clustering/jetstream_clustering
