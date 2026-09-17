---
title: "The synadia-io Repo Map (40 repos)"
section: "13-SYNADIA"
type: note
tags: [synadia, synadia-io, repos, map, github, nex, connect, control-plane, orbit]
source: https://github.com/synadia-io
created: 2026-06-14
---

# The `synadia-io` Repo Map

The full **`synadia-io` GitHub organization** — 40 repositories (37 active, 3 archived) as of 2026-06-14. This is the **company's** OSS footprint: the engines (NEX), the connector runtimes (Connect/Wombat), Control Plane assets, the Orbit client extensions, and security/auth libraries. Star counts are a point-in-time snapshot.

See also: [[NATS/13-SYNADIA/README|Synadia Index]] · [[NATS/12-ECOSYSTEM/repo-map|nats-io Repo Map]] · [[NATS/13-SYNADIA/synadia-products|Synadia Products]]

> **Legend:** ★ = snapshot stars · † = archived. Languages in brackets. Many product repos are "public assets" (docs/config), since the products themselves are commercial.

---

## Execution & workloads

| Repo | ★ | Lang | What |
|------|---|------|------|
| [nex](https://github.com/synadia-io/nex) | 310 | Go | **NATS Execution Engine** — run distributed workloads over NATS ([[NATS/13-SYNADIA/nex\|detail]]) |
| [nex-ce](https://github.com/synadia-io/nex-ce) | 0 | — | NEX Community Edition public assets |
| [ai-workflows](https://github.com/synadia-io/ai-workflows) | 4 | — | AI workflow assets |

---

## Connect (connectors / stream processing)

| Repo | ★ | Lang | What |
|------|---|------|------|
| [connect](https://github.com/synadia-io/connect) | 19 | Go | Connect Client SDK and CLI |
| [connect-runtime-wombat](https://github.com/synadia-io/connect-runtime-wombat) | 3 | Go | **Wombat** runtime (`~/bin/wombat`) — Benthos-style stream processor |
| [connect-runtime-vanilla](https://github.com/synadia-io/connect-runtime-vanilla) † | 1 | Go | Reference runtime template (archived) |
| [connectors](https://github.com/synadia-io/connectors) | 0 | — | Public connector assets |
| [connect-node](https://github.com/synadia-io/connect-node) | 0 | — | Connect Node public assets |
| [flink-connector-nats](https://github.com/synadia-io/flink-connector-nats) | 34 | Java | Apache Flink connector for NATS |
| [connector-intro-series](https://github.com/synadia-io/connector-intro-series) | 1 | Go | Tutorial series |

See [[NATS/13-SYNADIA/synadia-products|Products: Connect]] and [[NATS/12-ECOSYSTEM/bridges-connectors|Ecosystem: Bridges]].

---

## Control Plane & Platform

| Repo | ★ | Lang | What |
|------|---|------|------|
| [control-plane](https://github.com/synadia-io/control-plane) | 8 | — | Synadia Control Plane (management UI/API) |
| [control-plane-sdk-go](https://github.com/synadia-io/control-plane-sdk-go) | 4 | Go | Go SDK for Control Plane |
| [control-plane-cli](https://github.com/synadia-io/control-plane-cli) | 3 | — | Control Plane CLI |
| [synack](https://github.com/synadia-io/synack) | 1 | Go | **K8s operator** for Control Plane |
| [platform-component](https://github.com/synadia-io/platform-component) | 0 | Go | Standardized platform component interactions |
| [platform-trial](https://github.com/synadia-io/platform-trial) | 1 | Shell | Platform trial config |
| [helm-charts](https://github.com/synadia-io/helm-charts) | 0 | Go Tmpl | Helm charts for Synadia software |
| [http-gateway](https://github.com/synadia-io/http-gateway) | 0 | — | HTTP Gateway public assets |
| [private-link](https://github.com/synadia-io/private-link) | 1 | — | Private connectivity into Synadia Cloud |
| [schema-registry](https://github.com/synadia-io/schema-registry) | 0 | — | Message schema registry |
| [insights](https://github.com/synadia-io/insights) | 1 | — | Analytics/observability product |

---

## Orbit — client extensions (one per language)

`orbit.go` (83), `orbit.rs` (19), `orbit.js` (16), `orbit.java` (8), `orbit.net` (7), `orbit.py` (0), `orbit.c` (0). Helper libraries on top of the core clients — partitioned consumers, scheduled messages, distributed counters, JetStream publishers. Detail: [[NATS/13-SYNADIA/orbit|Orbit]].

---

## Security & auth libraries

| Repo | ★ | Lang | What |
|------|---|------|------|
| [callout.go](https://github.com/synadia-io/callout.go) | 23 | Go | Auth Callout library (external/delegated auth) |
| [callout.net](https://github.com/synadia-io/callout.net) | 6 | C# | Auth Callout for .NET |
| [jwt-auth-builder.go](https://github.com/synadia-io/jwt-auth-builder.go) | 17 | Go | Programmatic NATS JWT auth builder |
| [java-active-passive](https://github.com/synadia-io/java-active-passive) | 0 | Java | Active/Passive connection failover (JVM) |

See [[NATS/05-SECURITY/authentication|Authentication]].

---

## Learning, content & misc

| Repo | ★ | Lang | What |
|------|---|------|------|
| [rethink_connectivity](https://github.com/synadia-io/rethink_connectivity) | 136 | Go | "Rethink Connectivity" screencast assets |
| [awesome-nats](https://github.com/synadia-io/awesome-nats) | 80 | — | Curated NATS ecosystem list |
| [nats-data-series](https://github.com/synadia-io/nats-data-series) | 1 | Go | Event-driven future blog/resources |
| [workflows](https://github.com/synadia-io/workflows) | 0 | — | Reusable GitHub Action workflows |
| [testing.java](https://github.com/synadia-io/testing.java) | 0 | Java | Java client testing programs |
| [synadia-server](https://github.com/synadia-io/synadia-server) † | 0 | — | Synadia's NATS server distribution (archived) |
| [json.java](https://github.com/synadia-io/json.java) † | 0 | Java | JSON util (archived) |
| [.github](https://github.com/synadia-io/.github) | 0 | — | Org profile |
| [orbit.py](https://github.com/synadia-io/orbit.py) | 0 | — | (listed under Orbit) |

---

## Reading the org
- **Engines & runtimes** (`nex`, `connect*`) are the strategically interesting, genuinely-OSS pieces — usable sovereignly.
- **`control-plane*`, `platform*`, `http-gateway`, `private-link`, `insights`** are mostly **public-asset shells** for **commercial** products (the code is paid).
- **`orbit.*`, `callout.*`, `jwt-auth-builder`** are reusable **client/auth libraries** — grab-and-use.
- **3 archived**: `connect-runtime-vanilla`†, `synadia-server`†, `json.java`†.

## OneMind Relevance
- **In use:** `connect-runtime-wombat` → `~/bin/wombat`. **On radar:** `nex` (workload execution).
- **Candidate libraries:** `orbit.go` (partitioned/scheduled patterns), `callout.go` + `jwt-auth-builder.go` (fleet identity automation).
- Everything we'd actually adopt here is **Apache-licensed OSS** — the commercial control-plane/platform repos we can read but won't buy, consistent with self-hosting.

## Sources
- synadia-io org — https://github.com/synadia-io
- Snapshot data — `_assets/synadia-io-repos.jsonl`
- Synadia docs — https://docs.synadia.com
