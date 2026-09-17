---
title: "NATS Service Infrastructure — Evolution & Topology"
type: reference
domain: SO
status: active
created: 2026-05-26
updated: 2026-05-26
tags: [nats, infrastructure, cluster, leaf-node, super-cluster, gateway, topology]
source: https://docs.nats.io/nats-concepts/service_infrastructure
related: ["NATS-INFRASTRUCTURE.md", "synadia-research.md"]
---

# NATS Service Infrastructure — Evolution & Topology

Source: https://docs.nats.io/nats-concepts/service_infrastructure

---

## The Evolution Pattern

NATS infrastructure scales in stages — start simple, grow as needed.

### Stage 1 — Single Server (Local Dev)
```
[App] → [nats-server :4222]
```
- Single instance on local machine
- Used for development and local testing
- OneMind current state: MBP running single nats-server

### Stage 2 — Single Server or Cluster per Location
```
[Apps in VPC/Region] → [nats-server cluster]
```
- Deploy single server OR cluster of servers per VPC/region/on-prem location
- Apps connect to their local nats-server or cluster
- Multiple locations each get their own server/cluster

### Stage 3 — Connected Clusters (Global)
Two ways to connect location clusters together:

**Option A — Leaf Nodes**
```
[Local Cluster A] ──leaf──→ [Backbone Super-Cluster]
[Local Cluster B] ──leaf──→ [Backbone Super-Cluster]
[Local Cluster C] ──leaf──→ [Backbone Super-Cluster]
```

**Option B — Gateway Connections**
```
[Cluster A] ←──gateway──→ [Cluster B]
[Cluster B] ←──gateway──→ [Cluster C]
```

### Stage 4 — Global Super-Cluster (Enterprise/IoT Scale)
- Many clusters across multiple locations and cloud providers
- Connected into a global super-cluster
- Smart client routing to the "closest" NATS server cluster
- Target for: many IoT devices, applications deployed on end-user devices globally

---

## Topology Components

### Single Server
Simplest deployment. One `nats-server` instance. No redundancy.

### Cluster
Multiple `nats-server` instances in the same location working together.
- Full mesh of connections between servers in the cluster
- Built-in redundancy and load distribution
- All servers in cluster have access to all messages

### Super-Cluster
Multiple clusters connected via **gateway connections**.
- Each cluster is in a different region/cloud/on-prem
- Gateways connect clusters together into one global fabric
- Messages can flow across the entire super-cluster

### Leaf Node
A NATS server (or cluster) that connects to a backbone cluster/super-cluster.
- Extends the reach of the backbone into edge locations
- Edge location apps connect locally — traffic routes through leaf → backbone
- Lower latency for local communication, still connected globally
- Perfect for: IoT locations, edge nodes, remote facilities

---

## Running Your Own NATS Infrastructure

You can self-deploy any combination:
- Single servers
- Clusters of servers
- Super-clusters
- Leaf node servers

### Kubernetes
NATS recommends the official NATS Helm chart for K8s deployments.

### Containerization
Standard Docker/container deployment supported.

---

## OneMind OS Topology (Current + Future)

### Current
```
MBP M2 Max
  └── nats-server (single, localhost:4222)
        ├── Hermes (subscriber/publisher)
        ├── OpenClaw (subscriber/publisher)
        ├── HA Bridge (cmd.ha.*)
        ├── SAAS Worker (saas.execute)
        └── IoT/NATS.swift (iOS app)
```

### Future (as fleet grows)
```
MBP (Backbone / Hub Cluster)
  ├── VPS Leaf Node (onemindos-cloud)
  ├── Pi Leaf Nodes (IoT edge — Eden, sensors, drones)
  ├── iPhone/iOS (NATS.swift client)
  └── Drone nodes (when fleet active)
```

Leaf node pattern is the right evolution path for OneMind:
- Each robot/Pi/drone connects as a leaf to the MBP backbone
- Local comms stay local (low latency)
- All nodes visible on the global OneMind fabric
- No separate NATS cluster needed until fleet is very large

---

## Key Design Principles

1. **Start simple** — single server is fine for dev/small scale
2. **Clusters for redundancy** — within a location
3. **Leaf nodes for edge** — IoT, robots, remote locations
4. **Super-cluster for global** — multi-region, multi-cloud
5. **Client connects to nearest** — smart routing at scale

---

## Related Notes
- `NATS-INFRASTRUCTURE.md` — OneMind full NATS infra doc
- `synadia-research.md` — Synadia Platform (managed NATS)
- `synadia-agents-blog.md` — Synadia Agents SDK
