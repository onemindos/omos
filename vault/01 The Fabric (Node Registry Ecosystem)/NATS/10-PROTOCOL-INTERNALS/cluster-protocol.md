---
title: "Cluster Protocol (Server ↔ Server)"
section: "10-PROTOCOL-INTERNALS"
type: note
tags: [nats, protocol, cluster, routes, gossip, rs, rmsg, interest-propagation]
source: https://docs.nats.io/reference/nats-protocol/nats-server-protocol
created: 2026-06-14
---

# Cluster Protocol (Server ↔ Server)

Separate from the client protocol, NATS servers speak a **route protocol** to each other to form a cluster. It handles three jobs: **mesh formation** (auto-discovery via gossip), **interest propagation** (who's subscribed to what, so messages only cross routes when needed), and **message forwarding**.

See also: [[NATS/06-TOPOLOGY-SCALING/clustering|Clustering]] · [[NATS/06-TOPOLOGY-SCALING/superclusters-gateways|Superclusters & Gateways]] · [[NATS/10-PROTOCOL-INTERNALS/wire-protocol|Client Wire Protocol]]

---

## Full-mesh routes + gossip

Cluster servers connect on the **route port (6222 by default)** and form a **full mesh** — every server has a direct route to every other server (one hop max for intra-cluster delivery).

You only need to seed each server with *one or two* peers; the rest is **automatic** via gossip:

```
# server A config
cluster {
  name: omos
  listen: 0.0.0.0:6222
  routes: [ nats://B:6222 ]      # seed: just point at B
}
```

When A connects to B, B's `INFO` advertises **all the other routes B knows** (`connect_urls`). A dials those too. Add server C pointing only at A, and the whole mesh re-forms automatically — every node learns every other node. This is the gossip protocol: topology converges without you enumerating every pair.

---

## Route verbs

The route protocol mirrors the client one but with an `R` family for routed subscriptions/messages, so interest and traffic can be distinguished from local:

| Verb | Purpose |
|------|---------|
| `INFO` | Exchanged on route connect; carries `connect_urls` for gossip discovery |
| `RS+` | **Route Subscribe** — propagate that *this server* has interest in a subject (for an account) |
| `RS-` | **Route Unsubscribe** — interest gone; stop forwarding that subject |
| `RMSG` | **Routed Message** — forward a message to a peer server that has interest, scoped by account |
| `PING`/`PONG` | Route keepalive |

The critical idea is **interest propagation**: a server only sends `RMSG` to a peer if that peer previously sent `RS+` for the subject+account. No interest → no cross-route traffic. This is why an idle subject costs nothing cluster-wide, and why NATS Core fan-out scales — interest is summarized, not flooded.

```
# B has a subscriber on orders.* in account ACME → B tells A:
RS+ ACME orders.*\r\n
# Now when A receives a client PUB on orders.new in ACME, it forwards:
RMSG ACME orders.new 5\r\nhello\r\n
# When B's last subscriber leaves:
RS- ACME orders.*\r\n
```

Account is part of every route operation because subjects are **namespaced per account** — interest and delivery never leak across account boundaries. See [[NATS/05-SECURITY/accounts-multitenancy|Accounts & Multi-tenancy]].

---

## One hop, and the no-loop rule

In a full mesh, a message takes **at most one route hop**: client→serverA→serverB→subscriber. Servers do **not** forward a routed message to *other* routes (that would loop) — `RMSG` goes only to servers with direct interest. This bounds latency and eliminates routing loops without a spanning-tree algorithm.

For **multi-cluster** (geo) topologies, you don't extend the mesh — you use **gateways** (supercluster) which have their own optimized interest protocol and a strict no-loop "no two gateway hops" rule. **Leaf nodes** use yet another, lighter protocol (single outbound link). See [[NATS/06-TOPOLOGY-SCALING/superclusters-gateways|Superclusters & Gateways]] and [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]].

---

## Where JetStream fits

The route protocol moves *messages*. JetStream replication (R3/R5 streams) rides a **separate RAFT layer** over NATS subjects (the `$JS`/`$NRG` system subjects), not the route verbs directly. RAFT elects a stream leader and replicates the append log to followers. So: routes = real-time interest graph; RAFT = durable replicated state. See [[NATS/02-JETSTREAM/jetstream-clustering|JetStream Clustering]].

---

## OneMind Relevance
- Our fleet is **leaf-node-first** (edge sites dial a hub), not a giant flat cluster — but the hub itself (omos-mbp + future peers) is a route mesh. Seeding one peer and letting gossip converge is exactly how we'll add the second hub node.
- Interest propagation (`RS+`/`RS-`) is the performance guarantee behind our "hundreds of nodes": subjects nobody listens to cost nothing across the fabric.
- Account-scoped routing is what keeps FABRIC / INTERFACE / etc. traffic isolated even as it crosses servers.

## Sources
- NATS server (route) protocol — https://docs.nats.io/reference/nats-protocol/nats-server-protocol
- Clustering — https://docs.nats.io/running-a-nats-service/configuration/clustering
- Gossip / auto-discovery — https://docs.nats.io/running-a-nats-service/configuration/clustering/cluster_config
