---
title: "Orbit — Client Extensions"
section: "13-SYNADIA"
type: note
tags: [synadia, orbit, client-extensions, helpers, jetstream, callout, patterns]
source: https://github.com/synadia-io/orbit.go
created: 2026-06-14
---

# Orbit — Client Extensions

**Orbit** is Synadia's family of **client-side extension libraries** — reusable helpers and patterns layered *on top of* the core NATS clients, for things the base client doesn't ship but many apps want (partitioned consumers, scheduled messages, distributed counters, JetStream publish helpers, etc.). Open source, one library per language.

See also: [[NATS/08-CLIENTS-SDKS/README|Clients & SDKs]] · [[NATS/13-SYNADIA/synadia-products|Synadia Products]] · [[NATS/13-SYNADIA/repo-map|synadia-io Repo Map]]

---

## What Orbit is (and isn't)

- **Is:** a set of *add-on* utilities that sit beside the official client in your app — convenience patterns extracted from real deployments, maintained by Synadia.
- **Isn't:** a replacement client. You still use `nats.go`/`nats.py`/etc. for the connection; Orbit adds helpers on top.

The `.net` description is the most explicit example of the flavor: *"JetStream publisher, partitioned consumers, scheduled messages, distributed counters."* These are exactly the "I keep rewriting this" patterns Orbit packages once.

---

## The Orbit family (one per language)

| Repo | ★ | Lang |
|------|---|------|
| [orbit.go](https://github.com/synadia-io/orbit.go) | 83 | Go |
| [orbit.rs](https://github.com/synadia-io/orbit.rs) | 19 | Rust |
| [orbit.js](https://github.com/synadia-io/orbit.js) | 16 | TypeScript |
| [orbit.java](https://github.com/synadia-io/orbit.java) | 8 | Java |
| [orbit.net](https://github.com/synadia-io/orbit.net) | 7 | C# |
| [orbit.py](https://github.com/synadia-io/orbit.py) | 0 | Python |
| [orbit.c](https://github.com/synadia-io/orbit.c) | 0 | C |

Coverage mirrors the core client languages, so whatever SDK you use ([[NATS/08-CLIENTS-SDKS/README|Section 08]]), there's a matching Orbit.

---

## Typical Orbit helpers
- **Partitioned consumers** — deterministically shard a stream's messages across N workers by a key (ordered processing per partition). A common JetStream need not built into the base client.
- **Scheduled / delayed messages** — publish-now, deliver-later patterns.
- **Distributed counters** — atomic counters over KV/JetStream.
- **JetStream publisher helpers** — batching, async-ack windows, retry ergonomics.
- **Active/Passive connection** (`synadia-io/java-active-passive`) — failover patterns for the JVM client.

Exact contents vary by language and evolve; check each repo's README for its current helper set.

---

## Related: auth & infra helper libs

Adjacent Synadia OSS client libraries (not branded Orbit but same "extend the client" spirit):
- **`callout.go` / `callout.net`** — **Auth Callout** libraries: implement external/delegated authentication (the server calls *your* service to authorize a connection). Ties into [[NATS/05-SECURITY/authentication|Authentication]].
- **`jwt-auth-builder.go`** — programmatic builder for NATS JWT auth configurations (accounts/users), complementing `nsc`.
- **`flink-connector-nats`** — Apache Flink connector for NATS (stream-processing integration).

---

## OneMind Relevance
- When a fleet service needs **partitioned/ordered processing** or **scheduled delivery**, reach for **`orbit.go`** before hand-rolling it — battle-tested helpers, Apache-licensed, no lock-in.
- **`callout.go`** is the path if we ever need **external auth decisions** (e.g., authorize an agent connection against our own identity service) rather than static JWT scopes.
- **`jwt-auth-builder.go`** could script our account/user JWT provisioning as the fleet's identity needs grow beyond manual `nsc`.
- All OSS — consistent with the sovereign, self-hosted posture.

## Sources
- orbit.go — https://github.com/synadia-io/orbit.go
- orbit.net (helper list) — https://github.com/synadia-io/orbit.net
- callout.go — https://github.com/synadia-io/callout.go
- jwt-auth-builder.go — https://github.com/synadia-io/jwt-auth-builder.go
