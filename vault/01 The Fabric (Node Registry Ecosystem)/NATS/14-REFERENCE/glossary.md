---
title: "Glossary"
section: "14-REFERENCE"
type: note
tags: [nats, glossary, terms, definitions, reference]
created: 2026-06-14
---

# Glossary

Every NATS term, one line each. Linked to the section that explains it in depth.

See also: [[NATS/14-REFERENCE/README|Reference Index]] · [[NATS/01-CORE-CONCEPTS/README|Core Concepts]]

---

## Core messaging
- **Subject** — a dot-delimited name (`orders.us.new`) messages are published to and subscribed on. The address. See [[NATS/01-CORE-CONCEPTS/subjects-and-wildcards|Subjects]].
- **Token** — one dot-separated segment of a subject (`orders`, `us`, `new`).
- **Publish (PUB)** — send a message to a subject. Fire-and-forget in Core.
- **Subscribe (SUB)** — register interest in a subject (or wildcard) to receive matching messages.
- **Wildcard `*`** — matches exactly one token (`orders.*.new`).
- **Wildcard `>`** — matches one or more trailing tokens (`orders.>`); only valid at the end.
- **Queue group** — a named set of subscribers that **load-balance** messages (each msg to one member). Horizontal scaling. See [[NATS/01-CORE-CONCEPTS/queue-groups|Queue Groups]].
- **Request-Reply** — synchronous RPC pattern over pub/sub using a temporary reply **inbox**. See [[NATS/01-CORE-CONCEPTS/request-reply|Request-Reply]].
- **Inbox** — an auto-generated unique subject (`_INBOX.<nuid>`) for replies.
- **Core NATS** — the at-most-once, in-memory, no-persistence base layer. Fast, fire-and-forget.
- **Headers** — key/value metadata on a message (HPUB/HMSG), like HTTP headers. See [[NATS/10-PROTOCOL-INTERNALS/wire-protocol|Wire Protocol]].

## JetStream (persistence)
- **JetStream (JS)** — the persistence layer: streams, durable consumers, KV, Object Store, exactly-once. See [[NATS/02-JETSTREAM/jetstream-overview|JetStream]].
- **Stream** — a durable, replayable log of messages captured from subjects.
- **Consumer** — a stateful view/cursor over a stream that delivers messages to clients (push or pull).
- **Push consumer** — server pushes messages to a subject the client subscribes to.
- **Pull consumer** — client fetches batches on demand (preferred for scaling/flow control).
- **Durable** — a named consumer whose position survives disconnects/restarts.
- **Ack** — client acknowledgment that a message was processed (`ACK`, `NAK`, `TERM`, `WIP`).
- **Replay policy** — `instant` or `original` (preserve original timing).
- **Retention** — `limits`, `interest`, or `workqueue` — when messages are removed. See [[NATS/02-JETSTREAM/retention-and-limits|Retention]].
- **Discard policy** — `old` or `new` — what to drop when limits are hit.
- **R1/R3/R5** — JetStream replica count (1, 3, or 5) via RAFT. Odd numbers for quorum.
- **Mirror** — a stream that passively copies another stream. See [[NATS/02-JETSTREAM/mirrors-and-sources|Mirrors & Sources]].
- **Source** — a stream that aggregates messages *from* one or more other streams.
- **Nats-Msg-Id** — header used for **dedup** (exactly-once publish) within the dedup window.
- **KV (Key-Value)** — a bucket abstraction over a stream for key/value storage + watch. See [[NATS/03-KV-OBJECT-STORE/key-value-store|KV Store]].
- **Object Store** — chunked large-object storage over JetStream. See [[NATS/03-KV-OBJECT-STORE/object-store|Object Store]].
- **Domain** — a named, isolated JetStream instance (e.g. on a leaf) bridged to a hub. See [[NATS/06-TOPOLOGY-SCALING/domains|Domains]].

## Services / micro
- **Service API** — the micro framework: discovery, stats, versioning on top of request-reply. See [[NATS/04-SERVICES-MICRO/services-micro|Services API]].
- **Endpoint** — a named operation within a service, bound to a subject.

## Topology
- **Cluster** — a set of NATS servers fully meshed via routes, sharing subject interest. See [[NATS/06-TOPOLOGY-SCALING/clustering|Clustering]].
- **Route** — a server-to-server connection within a cluster (port 6222).
- **Gateway** — connection between clusters forming a **supercluster** (port 7222). See [[NATS/06-TOPOLOGY-SCALING/superclusters-gateways|Superclusters]].
- **Supercluster** — clusters of clusters, spanning regions, with optimized interest propagation.
- **Leaf node** — a satellite server that extends a remote cluster/account over one connection (port 7422). The edge primitive. See [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]].
- **Gossip** — the protocol by which cluster servers auto-discover each other.
- **Sublist** — the server's subscription-matching trie (subject → interested subs). See [[NATS/10-PROTOCOL-INTERNALS/server-architecture|Server Architecture]].

## Security & identity
- **NKEY** — an Ed25519 keypair used as a NATS identity (server/account/user). See [[NATS/10-PROTOCOL-INTERNALS/nuid-and-nkeys|NUID & NKEYS]].
- **Seed** — the private half of an NKEY (`S...`); guard it. Public half starts with the role letter (`U`,`A`,`O`,`N`,`C`).
- **JWT** — a NKEY-signed token encoding an account/user and its permissions. See [[NATS/05-SECURITY/jwt-nkeys-decentralized|JWT & NKEYS]].
- **Operator** — top of the JWT trust chain; signs accounts.
- **Account** — an isolation tenant; subjects are private to an account unless exported. See [[NATS/05-SECURITY/accounts-multitenancy|Accounts]].
- **Export / Import** — how accounts share streams or services across the account boundary.
- **`$SYS`** — the system account, exposing server events/monitoring. Used by Surveyor.
- **Auth Callout** — delegate auth decisions to an external service. See [[NATS/05-SECURITY/authentication|Authentication]].
- **NUID** — NATS Unique Identifier; fast collision-resistant IDs (inboxes, CIDs). See [[NATS/10-PROTOCOL-INTERNALS/nuid-and-nkeys|NUID & NKEYS]].

## Protocol & ops
- **PING/PONG** — keepalive/heartbeat between client and server.
- **RTT** — round-trip time; `nats rtt` measures it.
- **Slow consumer** — a subscriber that can't keep up; server drops messages to protect itself. See [[NATS/11-ADVANCED/benchmarking|Benchmarking]].
- **Drain** — graceful close: stop new messages, finish in-flight, then disconnect.
- **max_payload** — server cap on a single message size (default 1 MB).
- **varz / connz / jsz** — monitoring endpoints (server vitals / connections / JetStream). See [[NATS/14-REFERENCE/monitoring-endpoints|Monitoring Endpoints]].
- **Context** — a saved `nats` CLI connection profile (server + creds). See [[NATS/09-CLI-TOOLING/contexts|Contexts]].
- **creds file** — a file bundling a user JWT + NKEY seed for client auth.

## Legacy
- **STAN / NATS Streaming** — the **deprecated** predecessor to JetStream. Don't use for new work. See [[NATS/12-ECOSYSTEM/repo-map|Repo Map]].
- **gnatsd** — the old name/binary for `nats-server` (still seen in metric prefixes).

## OneMind Relevance
- Shared vocabulary for the fleet — when an SO doc or agent says "put it on a workqueue stream with an explicit-ack pull consumer," every term resolves here.
- The **identity** terms (NKEY/JWT/account/`$SYS`) are the ones to internalize first — fleet security is built on them.

## Sources
- NATS glossary/concepts — https://docs.nats.io/nats-concepts/overview
- Reference — https://docs.nats.io/reference
