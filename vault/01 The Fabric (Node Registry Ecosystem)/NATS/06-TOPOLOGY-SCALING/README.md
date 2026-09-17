---
title: "Topology & Scaling — Section Index"
section: "06-TOPOLOGY-SCALING"
type: index
tags: [nats, clustering, supercluster, gateways, leaf-nodes, scaling]
created: 2026-06-14
---

# 06 — Topology & Scaling

How NATS grows from a single server to a planet-spanning mesh. Four building blocks, layered:

```
Single server
   └─ Cluster          — full-mesh routes, horizontal scale + HA within a location
        └─ Supercluster — gateways connect clusters across regions (geo)
             └─ Leaf nodes — extend the fabric to edge/remote with store-and-forward
```

## Notes
- [[NATS/06-TOPOLOGY-SCALING/clustering|Clustering]] — routes, full mesh, seed servers, horizontal scale
- [[NATS/06-TOPOLOGY-SCALING/superclusters-gateways|Superclusters & Gateways]] — connecting clusters, optimistic spread, geo
- [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]] — edge extension, store-and-forward, hub-and-spoke
- [[NATS/06-TOPOLOGY-SCALING/domains|JetStream Domains]] — isolating JetStream per leaf/edge
- [[NATS/06-TOPOLOGY-SCALING/scaling-patterns|Scaling Patterns]] — sizing, when to use which layer

## See also
- [[NATS/02-JETSTREAM/jetstream-clustering|JetStream Clustering]] (RAFT replicas ride on the cluster)
- [[NATS/02-JETSTREAM/mirrors-and-sources|Mirrors & Sources]] (data movement across the topology)
