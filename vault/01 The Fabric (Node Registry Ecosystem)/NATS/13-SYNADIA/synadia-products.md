---
title: "Synadia Products"
section: "13-SYNADIA"
type: note
tags: [synadia, cloud, ngs, platform, control-plane, connect, http-gateway, private-link, commercial]
source: https://docs.synadia.com
created: 2026-06-14
---

# Synadia Products

The **commercial** layer Synadia sells on top of open-source NATS. None of this is required to run NATS — it's the managed/enterprise convenience tier. Knowing the product map tells you what's paid vs. what you can self-host for free (almost everything core).

See also: [[NATS/13-SYNADIA/history-and-governance|History & Governance]] · [[NATS/13-SYNADIA/nex|NEX]] · [[NATS/13-SYNADIA/repo-map|synadia-io Repo Map]]

---

## Synadia Cloud (NGS)

**Synadia Cloud** (historically **NGS** — "NATS Global System") is a **managed, global, multi-tenant NATS** service — a planet-spanning NATS supercluster you connect to instead of running your own. You get an **account** on their global fabric; your apps connect with creds and publish/subscribe globally without operating any servers.

- **Pay-as-you-go** SaaS; global low-latency footprint.
- Full NATS + JetStream + KV/Object Store, managed.
- The fastest way to "just use NATS" with zero ops — the opposite of the self-hosted, sovereign posture.

---

## Synadia Platform

**Synadia Platform** is the **self-managed / on-prem** commercial distribution — NATS plus enterprise tooling you run in **your own** infrastructure (cloud, on-prem, air-gapped). It bundles the **Control Plane**, support, and enterprise features for organizations that want the management layer but not the multi-tenant cloud.

Think: "Synadia Cloud's management experience, but inside your perimeter."

---

## Control Plane

The **Synadia Control Plane** is the **management/observability UI + API** for operating NATS systems at scale — clusters, accounts, JetStream assets, users/JWTs, monitoring — through a single pane instead of raw config + CLI. Repos: `synadia-io/control-plane`, `control-plane-cli`, `control-plane-sdk-go`, and `synack` (a **Kubernetes operator** for Control Plane).

This is the heart of what Synadia sells: not the bus (that's free), but the **operational control surface** over a big NATS estate.

---

## Connect

**Synadia Connect** is the **managed connector framework** — move data between NATS and external systems (databases, SaaS, cloud services, Kafka, etc.) via configured connectors instead of hand-built bridges. It's the modern successor to the archived `nats-io` bridge zoo ([[NATS/12-ECOSYSTEM/bridges-connectors|see 12]]).

- Runtimes: **`connect-runtime-wombat`** (the **Wombat** runtime — a Benthos/Redpanda-Connect-style stream processor) and a vanilla reference runtime.
- Public assets: `connectors`, `connect-node`.
- **OneMind already uses this**: `~/bin/wombat` is built from `connect-runtime-wombat`.

---

## Supporting products
- **HTTP Gateway** (`synadia-io/http-gateway`) — expose NATS over HTTP/REST for systems that can't speak the native protocol or WebSocket.
- **Private Link** (`synadia-io/private-link`) — private connectivity into Synadia Cloud (cloud-provider private networking, no public exposure).
- **Schema Registry** (`synadia-io/schema-registry`) — manage message schemas across the fabric.
- **Insights** (`synadia-io/insights`) — analytics/observability product.

---

## What's paid vs. free

| Capability | Open source (self-host) | Synadia commercial |
|------------|------------------------|--------------------|
| NATS server, JetStream, KV/OS | ✅ free, Apache-2.0 | (managed in Cloud/Platform) |
| Clients, CLI, nsc/nk | ✅ free | — |
| Monitoring (surveyor, exporter) | ✅ free | Control Plane (nicer UI) |
| Global multi-region fabric | DIY (gateways/superclusters) | **Synadia Cloud** |
| Management UI / single pane | DIY | **Control Plane** |
| Managed connectors | DIY (or OSS Wombat) | **Connect** |
| Enterprise support / SLAs | community | **Platform** |
| NEX (workload execution) | ✅ free (OSS) | (supported in Platform) |

The pattern: **everything to *run* NATS is free; you pay for *managing it at scale without effort* (UI, global cloud, support, managed connectors).**

---

## OneMind Relevance
- **We self-host OSS NATS — no Synadia subscription.** Sovereign, local, portable. The commercial products are the thing we deliberately *don't* buy.
- **Exception (and it's OSS):** **Connect/Wombat** is already in the toolchain (`~/bin/wombat`) for stream processing/connectors, and **NEX** is on the radar for workload execution — both open source, both fine for the sovereign posture.
- Knowing this map prevents accidental lock-in: if a tutorial says "spin up an NGS account," we know the self-hosted equivalent and take it.

## Sources
- Synadia products — https://www.synadia.com/platform
- Synadia Cloud — https://www.synadia.com/cloud
- Synadia docs — https://docs.synadia.com
- Connect — https://docs.synadia.com/cloud/walkthrough/connect
- Control Plane — https://github.com/synadia-io/control-plane
