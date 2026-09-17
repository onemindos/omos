---
title: "Leaf Nodes"
section: "06-TOPOLOGY-SCALING"
type: note
tags: [nats, leaf-nodes, edge, store-and-forward, hub-and-spoke, iot]
source: https://docs.nats.io/running-a-nats-service/configuration/leafnodes
created: 2026-06-14
---

# Leaf Nodes

A **leaf node** is a NATS server that connects **outward** to a remote cluster/server as a client of that remote, extending the fabric to the **edge** while keeping a local bus. It's the lightweight, single-link, store-and-forward answer for remote sites, devices, and disconnected operation.

See also: [[NATS/06-TOPOLOGY-SCALING/superclusters-gateways|Superclusters]] · [[NATS/06-TOPOLOGY-SCALING/domains|JetStream Domains]] · [[NATS/02-JETSTREAM/mirrors-and-sources|Mirrors & Sources]]

---

## The Model — hub & spoke

```
                  ┌──────────── HUB (cluster) ───────────┐
                  │  S1 — S2 — S3   (accounts, JetStream) │
                  └───────────────▲──────────────────────┘
            leaf connection (7422)│  (outbound from leaf)
        ┌─────────────────────────┴───────────┐
   ┌────▼────┐                            ┌────▼────┐
   │ barn-pi │  local clients + JS domain │ vehicle │  local clients
   │ (leaf)  │                            │ (leaf)  │
   └─────────┘                            └─────────┘
```

- A leaf runs a **full local NATS server** — local clients connect to it normally and get a working bus even if the hub is unreachable.
- The leaf **dials out** to the hub on the **leafnode port (7422)** and binds into a hub **account**.
- Subjects bridge between leaf and hub per the account's permissions — local traffic stays local; only what's needed crosses the link.

---

## Why Leaf Nodes (vs cluster/gateway)?

| Property | Leaf node |
|----------|-----------|
| **Link** | A **single outbound** connection (NAT/firewall friendly — leaf dials hub) |
| **Local autonomy** | Local clients keep working during hub disconnects |
| **Store-and-forward** | With a local JetStream **domain**, messages buffer locally and sync on reconnect |
| **Weight** | Light — one server, one link; no full-mesh routes or gateway mesh |
| **Security** | Binds into a specific hub **account**; credentials per leaf |
| **Topology** | Hub-and-spoke (and chains) — ideal for edge/IoT/remote sites |

Leaf nodes are how you reach **devices and remote locations** that can't (or shouldn't) be full cluster members.

---

## Configuration

```text
# Leaf server (barn-pi)
port: 4222                      # local clients here
leafnodes {
  remotes: [
    {
      url: "nats-leaf://hub-host:7422"
      credentials: "/etc/nats/leaf.creds"   # how this leaf authenticates to the hub
      account: "FABRIC"                       # hub account to bind into
    }
  ]
}
```

```text
# Hub server — accept leaf connections
leafnodes {
  listen: "0.0.0.0:7422"
}
```

Secure the link with `leafnodes { tls { } }`. The leaf authenticates with hub **credentials**, landing in a hub **account** that scopes what it can see.

---

## Store-and-Forward (the killer edge feature)

Give the leaf its **own JetStream domain** (see [[NATS/06-TOPOLOGY-SCALING/domains|Domains]]). Then:
1. Local producers publish to **local** JetStream — durable even with the hub down.
2. A **source** on the hub (or a **mirror** on the leaf) syncs streams across the link with `External` references.
3. On reconnect, buffered messages flow; nothing is lost during the outage.

This is exactly what an intermittently-connected barn, vehicle, or robot needs: keep working offline, sync when back.

---

## Chaining & Hierarchies

Leaf nodes can **chain** (a leaf can itself accept leaf connections), building hierarchies: device → site gateway → regional hub → central hub. Each level keeps local autonomy and forwards upward. Great for large, layered fleets.

---

## OneMind Relevance (central)

Leaf nodes are the **fleet backbone** for "hundreds of nodes":
- Each site (**House / Barn / Office**) and mobile node (**vehicles, robots, drone**) runs a **leaf** (often via embedded NATS, `omd --embed`), binding into the right hub **account** (FABRIC/INTERFACE/etc.).
- Local-first: a barn Pi keeps automating during a hub outage; telemetry buffers in a **local JS domain** and **sources** up to omos-mbp on reconnect.
- NAT-friendly single outbound link = no inbound ports on edge devices, fits the Tailscale-only doctrine.

See `~/onemindos-core/docs/NATS-DOCTRINE.md` for the OneMind leaf/account plan.

## Sources
- Leaf nodes — https://docs.nats.io/running-a-nats-service/configuration/leafnodes
- Leaf node config — https://docs.nats.io/running-a-nats-service/configuration/leafnodes/leafnode_conf
- JetStream on leaf nodes — https://docs.nats.io/running-a-nats-service/configuration/leafnodes/jetstream_leafnodes
