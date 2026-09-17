---
title: "Learning Resources (NATS by Example & more)"
section: "12-ECOSYSTEM"
type: note
tags: [nats, learning, nats-by-example, docs, awesome-nats, examples, resources]
source: https://natsbyexample.com
created: 2026-06-14
---

# Learning Resources

Where to actually learn and reference NATS beyond this KB: the runnable **NATS by Example** catalog, the canonical **docs**, the curated **awesome-nats** list, and the per-language example repos.

See also: [[NATS/12-ECOSYSTEM/repo-map|Repo Map]] · [[NATS/08-CLIENTS-SDKS/README|Clients & SDKs]] · [[NATS/README|KB Home]]

---

## NATS by Example — the runnable catalog

**https://natsbyexample.com** (org: `ConnectEverything`/`nats-io`). Short, **runnable, multi-language** examples organized by use case: pub/sub, request-reply, queue groups, JetStream (streams, consumers, KV, Object Store), and operations. Each shows the same pattern across Go, Python, JS, Java, Rust, C#, etc.

This is the **single best "how do I do X in language Y"** resource — when the conceptual docs aren't enough and you want working code, start here. It complements `docs.nats.io` (concepts) with copy-pasteable implementations.

---

## The canonical docs — `docs.nats.io`

Repo: `nats-io/nats.docs` (★151, HTML). The source of **https://docs.nats.io**, the authoritative reference this whole KB is built against. Structure:
- **NATS Concepts** — subjects, pub/sub, req/reply, JetStream, KV/OS
- **Using NATS** — developing with clients, CLI tools
- **Running a NATS Service** — config, clustering, security, k8s, MQTT/WS
- **Reference** — the wire protocol, ADRs

When this KB and the live docs disagree, **the live docs win** (note the build date here: 2026-06-14, server v2.14.2).

---

## `awesome-nats` — the curated list

Repo: `synadia-io/awesome-nats` (★80). A community-curated index of NATS libraries, tools, integrations, articles, and talks across the whole ecosystem — broader than just `nats-io` org repos. Good for discovering third-party clients and community projects. (Lives under the `synadia-io` org — see [[NATS/13-SYNADIA/repo-map|Section 13 repo map]].)

---

## Per-language example repos

The org keeps example repos for specific languages (handy as starting templates):

| Repo | Lang |
|------|------|
| [java-nats-examples](https://github.com/nats-io/java-nats-examples) | Java |
| [kotlin-nats-examples](https://github.com/nats-io/kotlin-nats-examples) | Kotlin |
| [java-nats-server-runner](https://github.com/nats-io/java-nats-server-runner) | Java (embed server in JVM tests) |

> Several older example repos (`go-nats-examples`†, `asyncio-nats-examples`†, `ts-nats-examples`†, `node-nats-examples`†, `ruby-nats-examples`†) are **archived** — their content largely migrated into **NATS by Example**. Prefer the website.

---

## Architecture & design — the ADRs

**https://github.com/nats-io/nats-architecture-and-design** holds the **ADRs** (Architecture Decision Records) — the deep, authoritative specs for JetStream, headers, the API, KV, Object Store, etc. When you need *exactly* how a feature is specified (not just how to use it), the ADRs are the ground truth beneath the docs.

---

## Community & support
- **Slack** — natsio.slack.com (the active community hub; invite via nats.io)
- **GitHub Discussions / Issues** — per-repo, on `nats-io`
- **YouTube** — "Rethink Connectivity" screencasts (`synadia-io/rethink_connectivity`), KubeCon talks
- **nats.io** — blog, adopters, comparisons

---

## A suggested learning path
1. **Concepts** → this KB Sections [[NATS/00-OVERVIEW/README|00]]–[[NATS/04-SERVICES-MICRO/README|04]] + `docs.nats.io` concepts.
2. **Hands-on** → **NATS by Example** in your language + `nats` CLI ([[NATS/09-CLI-TOOLING/README|09]]).
3. **Operate** → Sections [[NATS/05-SECURITY/README|05]]–[[NATS/07-DEPLOYMENT-OPS/README|07]] + monitoring tools.
4. **Go deep** → the **ADRs** + Section [[NATS/10-PROTOCOL-INTERNALS/README|10]] internals.
5. **Productionize** → Section [[NATS/13-SYNADIA/README|13]] (Synadia) if you want managed/commercial.

---

## OneMind Relevance
- **NATS by Example** is the go-to when implementing a new pattern in a fleet service — grab the canonical idiom rather than reinventing it.
- The **ADRs** are the citation of last resort when this KB or the docs are ambiguous about exact JetStream/KV semantics — important when we design durable subject law.
- `awesome-nats` is worth a periodic scan for community tools that could serve the fleet.

## Sources
- NATS by Example — https://natsbyexample.com
- Docs — https://docs.nats.io · repo https://github.com/nats-io/nats.docs
- ADRs — https://github.com/nats-io/nats-architecture-and-design
- awesome-nats — https://github.com/synadia-io/awesome-nats
