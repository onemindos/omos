---
title: "JetStream Domains"
section: "06-TOPOLOGY-SCALING"
type: note
tags: [nats, jetstream, domains, leaf-nodes, edge, isolation]
source: https://docs.nats.io/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes
created: 2026-06-14
---

# JetStream Domains

A **JetStream domain** is an isolated JetStream "universe" with its own name and `$JS.<domain>.API.>` namespace. Domains let an **edge/leaf** run its **own** JetStream independent of the hub's, while still allowing controlled data movement between them.

See also: [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]] · [[NATS/02-JETSTREAM/mirrors-and-sources|Mirrors & Sources]]

---

## The Problem Domains Solve

Without domains, a leaf node binding into a hub account would share **one** JetStream — the leaf's local streams and the hub's streams would collide in the same `$JS.API` namespace, and a leaf couldn't have truly local, offline-capable persistence. **Domains partition JetStream** so each side has its own:

```
Hub JetStream domain:  "hub"   → API on  $JS.hub.API.>
Leaf JetStream domain: "barn"  → API on  $JS.barn.API.>
```

A client talks to its **local** domain by default; to reach another domain it uses that domain's **API prefix**.

---

## Configuration

```text
# Hub
jetstream {
  domain: "hub"
  store_dir: "/var/lib/nats/hub"
}

# Leaf (barn-pi)
jetstream {
  domain: "barn"
  store_dir: "/var/lib/nats/barn"
}
leafnodes { remotes: [ { url: "nats-leaf://hub:7422", account: "FABRIC" } ] }
```

Each domain has its **own streams, consumers, KV, and Object Store**. `nats --js-domain barn stream ls` targets the leaf's domain; `--js-domain hub` targets the hub's.

---

## Moving Data Across Domains

Use **mirrors/sources** with an `External` reference that points at the other domain's **API prefix** and a **delivery prefix** (see [[NATS/02-JETSTREAM/mirrors-and-sources|Mirrors & Sources]]):

```
# Leaf sources a hub stream down for local offline read
stream TEL_LOCAL  ──source──▶  External{ api: "$JS.hub.API", deliver: "deliver.barn" }  ← hub TEL
```

Cross-domain access requires the in-between accounts to **export/import** the relevant `$JS.<domain>.API.*`, `$JS.ACK.*`, and `$JS.FC.*` (flow-control) subjects — exactly the export/import pattern in [[NATS/05-SECURITY/accounts-multitenancy|Accounts]].

---

## Domains vs Accounts vs Clusters

| Concept | Isolates | Scope |
|---------|----------|-------|
| **Cluster** | nothing (it's *one* logical bus) | servers in a location |
| **Account** | **subject space** (tenant wall) | within a server/cluster |
| **Domain** | **JetStream** (separate persistence universe) | typically hub vs each leaf/edge |

A leaf can be in account `FABRIC` *and* JetStream domain `barn` — accounts wall subjects, domains wall persistence. They compose.

---

## Why This Matters for the Edge

- **Local-first persistence** — the edge keeps a fully functional JetStream offline (store-and-forward).
- **No collision** — leaf streams (`barn`) never clash with hub streams (`hub`).
- **Selective sync** — only the streams you explicitly mirror/source cross the link.
- **Clean addressing** — `--js-domain` makes "which JetStream am I talking to" explicit.

---

## OneMind Relevance & Pitfall

Each edge node (barn, vehicles, robots) gets its **own domain** so it runs autonomous JetStream and **sources** up to the `hub` domain on reconnect. **Pitfall (verified):** changing the JetStream `domain` **orphans** existing streams/KV created under the old prefix — they become invisible under the new `$JS.<domain>.API`. Decide the domain name up front; if you must change it, plan to recreate or migrate assets. This is recorded in the OneMind KV/JetStream pitfalls.

## Sources
- JetStream on leaf nodes / domains — https://docs.nats.io/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes
- JetStream domains — https://docs.nats.io/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes#configuring-domains
- Cross-account JetStream — https://docs.nats.io/running-a-nats-service/nats_admin/jetstream_admin
