---
title: "Scaling Patterns & Sizing"
section: "06-TOPOLOGY-SCALING"
type: note
tags: [nats, scaling, sizing, patterns, capacity-planning]
source: https://docs.nats.io/running-a-nats-service/configuration/clustering
created: 2026-06-14
---

# Scaling Patterns & Sizing

A decision guide for *which* topology layer to reach for, plus sizing rules of thumb.

See also: [[NATS/06-TOPOLOGY-SCALING/clustering|Clustering]] · [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]] · [[NATS/02-JETSTREAM/jetstream-clustering|JetStream Clustering]]

---

## The Decision Tree

```
Need more client/throughput capacity in ONE location?      → CLUSTER (add servers)
Need JetStream HA (survive a node failure)?                 → CLUSTER + R3/R5 streams
Need ONE bus across multiple regions/data centers?          → SUPERCLUSTER (gateways)
Need to reach EDGE devices / remote sites / offline ops?    → LEAF NODES (+ JS domains)
Need tenant isolation on the same infra?                    → ACCOUNTS (not a topology change)
Need separate persistence universes (edge vs hub)?          → JETSTREAM DOMAINS
```

---

## Sizing Rules of Thumb

| Concern | Guidance |
|---------|----------|
| **Core throughput** | A single server handles **millions** of msgs/sec for small messages; scale clients across cluster members before worrying |
| **JetStream cluster size** | **3** servers (tolerate 1 failure) or **5** (tolerate 2). Always **odd** for quorum |
| **Replicas** | **R3** for production streams; **R5** only for the most critical on 5+ nodes; **R1** for dev/edge-local |
| **`max_payload`** | Default ~**1 MB**; raise cautiously (memory per message). Big blobs → Object Store, not bigger payloads |
| **File store** | Size `max_file_store` to data × replicas × growth headroom; fast disk (SSD/NVMe) for JetStream |
| **Connections** | One server handles **tens of thousands+** of connections; spread across cluster for headroom |
| **Failure domains** | Distribute cluster nodes across racks/AZs so one failure domain ≠ quorum loss |

---

## Common Patterns

### 1. Single Server (homestead hub)
One always-on box. Simple, R1 JetStream. Fine until you need HA or scale. *(OneMind today: omos-mbp.)*

### 2. 3-Node Cluster (HA core)
Three servers, R3 streams. Survives one node failure with no data loss. The standard production baseline.

### 3. Hub + Leaf Edge (fleet)
A hub (single or cluster) + many **leaf nodes** at the edge, each with a local **JetStream domain** for offline operation, **sourcing** data up to the hub. The IoT/fleet pattern.

### 4. Multi-Region Supercluster
Several regional **clusters** stitched by **gateways**; **regional quorums**, **cross-region mirrors** for read-locality. Global scale + DR.

### 5. Stretch + Edge (full stack)
Supercluster of regional clusters **+** leaf nodes hanging off each region for the edge. The "planet-scale with edge" topology NATS is uniquely good at.

---

## Anti-Patterns

- **Stretching one cluster across the WAN** — routes assume low latency; use a **supercluster** instead.
- **Spanning a JetStream RAFT quorum across regions** — high-latency quorum = slow/unstable; keep quorums regional, replicate with **mirrors**.
- **Even-numbered JetStream clusters** — can't form a clean majority; use odd counts.
- **Giant messages instead of Object Store** — raising `max_payload` to hundreds of MB wrecks memory; chunk via Object Store.
- **One shared JetStream for hub + edge** — use **domains** so edge persistence is independent and offline-capable.

---

## OneMind Scaling Path (concrete)

1. **Now:** single server (omos-mbp), R1 — correct for one always-on hub.
2. **Next:** 3-node cluster (MBP + omos-cloud + 1) for Core HA + **R3** on critical streams (AGENT_EVENTS, telemetry).
3. **Fleet:** **leaf nodes** at House/Barn/Office and on vehicles/robots, each with a local **JS domain**, sourcing up to the hub — local-first, sync-on-reconnect, NAT-friendly over Tailscale.
4. **Far future:** multi-site **supercluster** if OneMind spans regions.

## Sources
- Clustering & sizing — https://docs.nats.io/running-a-nats-service/configuration/clustering
- JetStream clustering — https://docs.nats.io/running-a-nats-service/configuration/clustering/jetstream_clustering
- Leaf nodes — https://docs.nats.io/running-a-nats-service/configuration/leafnodes
- Slow consumers / limits — https://docs.nats.io/running-a-nats-service/configuration
