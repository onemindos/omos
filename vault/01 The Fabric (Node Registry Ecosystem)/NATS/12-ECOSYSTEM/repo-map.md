---
title: "The nats-io Repo Map (100 repos)"
section: "12-ECOSYSTEM"
type: note
tags: [nats, nats-io, repos, ecosystem, map, github]
source: https://github.com/nats-io
created: 2026-06-14
---

# The `nats-io` Repo Map

A categorized map of the **entire `nats-io` GitHub organization** — 100 repositories (53 active, 47 archived) as of 2026-06-14. Star counts are a point-in-time snapshot for relative signal only.

See also: [[NATS/12-ECOSYSTEM/README|Ecosystem Index]] · [[NATS/08-CLIENTS-SDKS/README|Clients & SDKs]] · [[NATS/13-SYNADIA/repo-map|synadia-io Repo Map]]

> **Legend:** ★ = snapshot stars · † = archived/deprecated. Languages in brackets.

---

## Server & runtime

| Repo | ★ | Lang | What |
|------|---|------|------|
| [nats-server](https://github.com/nats-io/nats-server) | 20019 | Go | The server. The whole thing in one binary. |
| [nats-docker](https://github.com/nats-io/nats-docker) | 166 | Dockerfile | Official Docker image |
| [nats-box](https://github.com/nats-io/nats-box) | 150 | HCL | Container with all NATS CLI utilities (`nats`, `nsc`, `nk`…) |
| [graft](https://github.com/nats-io/graft) | 189 | Go | RAFT election implementation (consensus primitive) |
| [official-images](https://github.com/nats-io/official-images) | 9 | Shell | Docker Official Images plumbing |

---

## Core clients (active)

| Repo | ★ | Lang | What |
|------|---|------|------|
| [nats.go](https://github.com/nats-io/nats.go) | 6650 | Go | **Reference client** — most complete, can embed the server |
| [nats.rs](https://github.com/nats-io/nats.rs) | 1487 | Rust | Async Rust client |
| [nats.py](https://github.com/nats-io/nats.py) | 1230 | Python | Asyncio client |
| [nats.rb](https://github.com/nats-io/nats.rb) | 892 | Ruby | Ruby client |
| [nats.java](https://github.com/nats-io/nats.java) | 665 | Java | JVM client |
| [nats.net.v1](https://github.com/nats-io/nats.net.v1) | 651 | C# | .NET client (v1 line) |
| [nats.c](https://github.com/nats-io/nats.c) | 494 | C | C client (embedded/firmware) |
| [nats.deno](https://github.com/nats-io/nats.deno) | 179 | TypeScript | Deno client |
| [nats.ex](https://github.com/nats-io/nats.ex) | 274 | Elixir | Elixir client |
| [nats-pure.rb](https://github.com/nats-io/nats-pure.rb) | 162 | Ruby | Pure-Ruby client (no C ext) |
| [nats.cr](https://github.com/nats-io/nats.cr) | 47 | Crystal | Crystal client |

See [[NATS/08-CLIENTS-SDKS/README|Section 08]] for usage of the major ones. (The JS/TS browser+node clients now consolidate under the `nats.js`/`@nats-io` family; `nats.ws`†, `nats.ts`†, `nats.node`† are archived predecessors.)

---

## CLI, tooling & JetStream libs

| Repo | ★ | Lang | What |
|------|---|------|------|
| [natscli](https://github.com/nats-io/natscli) | 785 | Go | The `nats` CLI — pub/sub, streams, KV, bench, contexts |
| [nats-top](https://github.com/nats-io/nats-top) | 402 | Go | `top`-style live server monitor |
| [nsc](https://github.com/nats-io/nsc) | 134 | Go | Operator/account/user JWT + NKEY management |
| [jsm.go](https://github.com/nats-io/jsm.go) | 186 | Go | JetStream management library (Go) |
| [jetstream-gh-action](https://github.com/nats-io/jetstream-gh-action) | 9 | Go | JetStream GitHub Actions |
| [cliprompts](https://github.com/nats-io/cliprompts) | 2 | Go | CLI prompt utils (used by nsc) |

See [[NATS/09-CLI-TOOLING/README|Section 09]].

---

## Security & identity primitives

| Repo | ★ | Lang | What |
|------|---|------|------|
| [nkeys](https://github.com/nats-io/nkeys) | 163 | Go | Ed25519 NKEYS (identity) |
| [jwt](https://github.com/nats-io/jwt) | 88 | Go | NATS JWT (NKey-signed Ed25519 tokens) |
| [nuid](https://github.com/nats-io/nuid) | 374 | Go | Fast unique IDs (inboxes, CIDs) |
| [nkeys.js](https://github.com/nats-io/nkeys.js) | 24 | TypeScript | NKEYS for JS/Deno/browser |
| [nkeys.py](https://github.com/nats-io/nkeys.py) | 19 | Python | NKEYS for Python |
| [nkeys.rb](https://github.com/nats-io/nkeys.rb) | 4 | Ruby | NKEYS for Ruby |
| [nuid.js](https://github.com/nats-io/nuid.js) | 14 | TypeScript | NUID for JS |

See [[NATS/10-PROTOCOL-INTERNALS/nuid-and-nkeys|NUID & NKEYS]] and [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS (Security)]].

---

## Monitoring & observability

| Repo | ★ | Lang | What |
|------|---|------|------|
| [nats-surveyor](https://github.com/nats-io/nats-surveyor) | 317 | Go | Fleet monitoring via `$SYS` → Prometheus/Grafana |
| [prometheus-nats-exporter](https://github.com/nats-io/prometheus-nats-exporter) | 487 | Go | Scrape `/varz`,`/connz`,`/jsz` into Prometheus |
| [not.go](https://github.com/nats-io/not.go) | 104 | Go | Distributed tracing reference (Go) |
| [not.java](https://github.com/nats-io/not.java) | 2 | Java | Distributed tracing reference (Java) |

Detailed in [[NATS/12-ECOSYSTEM/monitoring-exporters|Monitoring & Exporters]].

---

## Kubernetes & IaC

| Repo | ★ | Lang | What |
|------|---|------|------|
| [k8s](https://github.com/nats-io/k8s) | 635 | Go | **Helm charts** — the supported way to run NATS on k8s |
| [nack](https://github.com/nats-io/nack) | 244 | Go | NATS Controllers for Kubernetes — declarative streams/consumers (CRDs) |
| [terraform-provider-jetstream](https://github.com/nats-io/terraform-provider-jetstream) | 82 | Go | Manage JetStream assets via Terraform |
| [nats-operator](https://github.com/nats-io/nats-operator) † | 568 | Go | **Deprecated** operator — use `k8s` Helm + `nack` instead |

Detailed in [[NATS/12-ECOSYSTEM/kubernetes-iac|Kubernetes & IaC]].

---

## Bridges & connectors

| Repo | ★ | Lang | What |
|------|---|------|------|
| [nats-kafka](https://github.com/nats-io/nats-kafka) | 161 | Go | NATS ⇄ Kafka bridge |
| [nginx-nats](https://github.com/nats-io/nginx-nats) | 113 | C | NGINX module client for NATS |
| [spring-nats](https://github.com/nats-io/spring-nats) | 82 | Java | Spring Cloud Stream binder |

Detailed in [[NATS/12-ECOSYSTEM/bridges-connectors|Bridges & Connectors]]. (Archived bridges: `nats-mq`†, `nats-jms-bridge`†, `nats-connector-framework`†, `nats-connector-redis`†, `nats-replicator`†.)

---

## Docs & learning

| Repo | ★ | Lang | What |
|------|---|------|------|
| [nats.docs](https://github.com/nats-io/nats.docs) | 151 | HTML | Source of `docs.nats.io` — canonical markdown |
| [nats-site](https://github.com/nats-io/nats-site) | 104 | Markdown | nats.io website content |
| [nats-general](https://github.com/nats-io/nats-general) | 70 | — | General info, governance |
| [advisories](https://github.com/nats-io/advisories) | 3 | HTML | Security advisories |
| Example repos | — | various | `java-nats-examples`, `kotlin-nats-examples`, `homebrew-nats-tools` (taps), `nats-choco` (Chocolatey) |

See [[NATS/12-ECOSYSTEM/nats-by-example|Learning Resources]].

---

## Streaming legacy (STAN) — all archived †

The **NATS Streaming (STAN)** project is **end-of-life**; **JetStream replaced it**. Listed for recognition only — do not build on these:

`nats-streaming-server`† (2527★), `stan.go`† (707★), `stan.js`† (290★), `stan.py`† (113★), `stan.java`† (94★), `stan.net`† (137★), `stan.rb`† (20★), `nats-streaming-operator`† (177★), `nats-streaming-docker`†, `go-nats-streaming`†.

> If you encounter STAN in old code or docs, the migration target is **JetStream** ([[NATS/02-JETSTREAM/jetstream-overview|Section 02]]). Channels → streams, durable subscriptions → consumers.

---

## Other archived † (historical)

Old per-language client predecessors and demos, kept for history: `go-nats`† (now `nats.go`), `nats.ts`†/`nats.node`†/`nats.ws`† (now `nats.js`), `nats.py2`†, `elixir-nats`†, plus example/demo repos (`go-nats-examples`†, `asyncio-nats-examples`†, `ts-nats-examples`†, `node-nats-examples`†, `ruby-nats-examples`†), bridges, `nats-account-server`† (replaced by built-in resolver), `nats-rest-config-proxy`†, `sublist`† (history of the routing trie), `nats-on-a-log`†, `latency-tests`†.

---

## How to use this map
- **Building a client app?** Pick from *Core clients* → [[NATS/08-CLIENTS-SDKS/README|Section 08]].
- **Operating NATS?** *CLI/tooling* + *Monitoring* + *Kubernetes*.
- **Securing it?** *Security primitives* + [[NATS/05-SECURITY/README|Section 05]].
- **Saw a `stan.*` repo?** It's dead — use JetStream.

## Sources
- nats-io org — https://github.com/nats-io
- Snapshot data — `_assets/nats-io-repos.jsonl`
