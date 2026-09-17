---
title: "Clustering"
section: "06-TOPOLOGY-SCALING"
type: note
tags: [nats, clustering, routes, full-mesh, high-availability]
source: https://docs.nats.io/running-a-nats-service/configuration/clustering
created: 2026-06-14
---

# Clustering

A **cluster** is a set of NATS servers connected by **routes**, forming a **full mesh** that behaves as one logical bus. Clients connect to any server; messages reach subscribers on any other. Clustering gives **horizontal scale** and **high availability** within a location.

See also: [[NATS/06-TOPOLOGY-SCALING/superclusters-gateways|Superclusters]] · [[NATS/02-JETSTREAM/jetstream-clustering|JetStream Clustering]]

---

## How It Works

- Servers connect to each other on a **route** port (default **6222**), separate from the client port (4222).
- The cluster forms a **full mesh** — every server has a route to every other server (one hop max between any two servers).
- Servers **gossip** topology: tell one server about a couple of peers and it discovers the rest. Clients learn the full server list via `INFO` and fail over automatically.
- A message published on server S1 is forwarded only to servers that have **interest** (a subscriber) — routing is interest-based, not blind flooding.

```
        client          client          client
          │               │               │
        ┌─▼─┐  route     ┌─▼─┐  route     ┌─▼─┐
        │S1 │◀──────────▶│S2 │◀──────────▶│S3 │
        └─┬─┘            └─┬─┘            └─┬─┘
          └──────route─────┴──────route─────┘     (full mesh)
```

---

## Configuration

```text
# server S1
port: 4222
cluster {
  name: "OMOS"
  listen: "0.0.0.0:6222"
  routes: [
    "nats-route://s2-host:6222"
    "nats-route://s3-host:6222"
  ]
}
```

- **`cluster.name`** must match across all members.
- **`routes`** list seed peers — you don't need every server in every config; gossip fills the rest. Listing 1–2 seeds is enough.
- Add TLS under `cluster { tls { } }` to encrypt route traffic (see [[NATS/05-SECURITY/tls-encryption|TLS]]).

CLI quick cluster (dev):
```shell
nats-server -p 4222 -cluster nats://0.0.0.0:6222 -cluster_name OMOS
nats-server -p 4223 -cluster nats://0.0.0.0:6223 -routes nats://0.0.0.0:6222 -cluster_name OMOS
```

---

## What Clustering Gives You

| Benefit | Detail |
|---------|--------|
| **Horizontal scale** | Spread client connections + subject traffic across N servers |
| **High availability** | A server dies → clients reconnect to a peer, subscriptions restored |
| **No single point of failure** | Full mesh, no leader for *Core* routing |
| **JetStream HA** | RAFT replica groups (R3/R5) ride on top — see [[NATS/02-JETSTREAM/jetstream-clustering|JetStream Clustering]] |
| **Transparent to clients** | Clients address subjects, not servers; topology changes are invisible |

---

## Sizing & Best Practices

- **Core** clustering scales horizontally with no required odd count. **JetStream** wants an **odd** number (3 or 5) for meta-group quorum.
- Keep cluster members on a **low-latency network** (same region/LAN) — routes assume fast links. For cross-region, use a **supercluster** with gateways, not a stretched cluster.
- Put members in different **failure domains** (racks/AZs).
- Secure route + monitoring ports; never expose 6222 publicly.
- Watch `/routez` and `nats server report` for route health.

---

## Cluster vs Supercluster vs Leaf (quick contrast)

| Layer | Connects | Link assumption | Use |
|-------|----------|-----------------|-----|
| **Cluster** | servers → one logical bus | low-latency (LAN/region) | scale + HA in one place |
| **Supercluster** | clusters → clusters via gateways | WAN/geo | multi-region |
| **Leaf node** | a server/edge → a remote cluster | any (incl. flaky) | edge extension, store-and-forward |

---

## OneMind Relevance

Today: **single server** (omos-mbp). Growth path: a small **3-node cluster** (MBP + omos-cloud + one always-on box) for Core HA and **R3 JetStream** on critical streams — keeping members on the **Tailscale** mesh as the "network." Cross-region/edge stays the job of **leaf nodes** and **superclusters**, not a stretched cluster.

## Sources
- Clustering — https://docs.nats.io/running-a-nats-service/configuration/clustering
- Cluster config — https://docs.nats.io/running-a-nats-service/configuration/clustering/cluster_config
- Routes — https://docs.nats.io/running-a-nats-service/configuration/clustering/v2_routes
