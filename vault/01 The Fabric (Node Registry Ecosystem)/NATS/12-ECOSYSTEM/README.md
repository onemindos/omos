---
title: "Ecosystem — Section Index"
section: "12-ECOSYSTEM"
type: index
tags: [nats, ecosystem, nats-io, repos, surveyor, exporter, kubernetes, integrations]
source: https://github.com/nats-io
created: 2026-06-14
---

# 12 — Ecosystem

The wider **`nats-io`** open-source world beyond the server and core clients: monitoring tools, Prometheus exporters, Kubernetes operators, bridges to other systems, security tooling, and the learning resources. As of this build the `nats-io` org has **100 repositories — 53 active, 47 archived** (the archived half is mostly the deprecated **NATS Streaming / STAN** lineage and old per-language example repos).

> This section maps the **community/OSS** ecosystem under `nats-io`. The **commercial** world (Synadia Cloud, Platform, Control Plane, NEX, Connect, Orbit) lives in [[NATS/13-SYNADIA/README|Section 13]].

See also: [[NATS/07-DEPLOYMENT-OPS/monitoring|Monitoring]] · [[NATS/09-CLI-TOOLING/README|CLI & Tooling]] · [[NATS/13-SYNADIA/README|Synadia]]

## Notes
- [[NATS/12-ECOSYSTEM/repo-map|The `nats-io` Repo Map]] — all 100 repos, categorized (server, clients, tooling, monitoring, k8s, bridges, security, streaming-legacy)
- [[NATS/12-ECOSYSTEM/monitoring-exporters|Monitoring & Exporters]] — `nats-surveyor`, `prometheus-nats-exporter`, `nats-top`
- [[NATS/12-ECOSYSTEM/bridges-connectors|Bridges & Connectors]] — Kafka bridge, NGINX module, Spring, Terraform, tracing
- [[NATS/12-ECOSYSTEM/kubernetes-iac|Kubernetes & IaC]] — Helm charts, NACK controllers, Terraform provider
- [[NATS/12-ECOSYSTEM/nats-by-example|Learning Resources]] — NATS by Example, docs repo, awesome-nats, examples

## The ecosystem at a glance

```
nats-io (100 repos)
├── Server            nats-server, nats-docker, nats-box
├── Core clients      nats.go, nats.py, nats.rs, nats.java, nats.net, nats.c,
│                     nats.rb, nats.ex, nats.deno, nats.cr  (+ examples repos)
├── CLI & tooling     natscli, nsc, nats-top
├── Security prims    nkeys, jwt, nkeys.js/.py/.rb
├── Identity/ID       nuid, nuid.js
├── Monitoring        nats-surveyor, prometheus-nats-exporter
├── Kubernetes        k8s (Helm), nack (controllers), [nats-operator † archived]
├── IaC               terraform-provider-jetstream
├── Bridges           nats-kafka, nginx-nats, spring-nats, [nats-mq/jms † archived]
├── JetStream libs    jsm.go, jetstream-gh-action
├── Tracing           not.go, not.java
├── Docs/learn        nats.docs, nats-site, nats-general
└── Streaming legacy † nats-streaming-server, stan.* (ALL archived — use JetStream)
```
† archived / deprecated

## Reading the map
- **Star counts** in this section are a **point-in-time snapshot** (build date 2026-06-14), used only for relative signal — not live.
- **Archived ≠ broken**, but archived ≠ maintained. For anything new, prefer active repos; treat the whole **STAN** family as historical (JetStream replaced it).

## OneMind Relevance
- We already run pieces of this ecosystem: `natscli`, `nats-top`, `nats-surveyor` (`io.nats.surveyor` launchd agent), and the Prometheus exporter feed our observability.
- The **k8s/NACK** and **Terraform** repos matter when the fleet outgrows the single Mac hub and JetStream assets need declarative management.
- Bridges (`nats-kafka`) are the escape hatch if we ever need to interop with an external Kafka — but the doctrine is NATS-native first.

## Sources
- nats-io org — https://github.com/nats-io
- Repo snapshot — `_assets/nats-io-repos.jsonl` (this folder)
- NATS by Example — https://natsbyexample.com
