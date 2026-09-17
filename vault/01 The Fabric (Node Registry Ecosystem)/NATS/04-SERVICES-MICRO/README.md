---
title: "Services / Micro — Section Index"
section: "04-SERVICES-MICRO"
type: index
tags: [nats, micro, services]
created: 2026-06-14
---

# 04 — Services / Micro Framework

NATS's built-in framework for discoverable, observable, load-balanced request-reply microservices.

## Notes
- [[NATS/04-SERVICES-MICRO/services-micro|Services / Micro Framework]] — service/group/endpoint model, `$SRV` discovery (PING/INFO/STATS), client API, auto stats, versioning, load balancing

## One-line model
> A Micro **service** is request-reply + queue-group load balancing + a standardized discovery/stats contract on the reserved `$SRV.>` subjects. Write a handler, get a discoverable, self-reporting, horizontally-scalable service.

## See also
- [[NATS/01-CORE-CONCEPTS/request-reply|Request/Reply]] · [[NATS/01-CORE-CONCEPTS/queue-groups|Queue Groups]]
- OneMind: `fabric_discover`, the `saas-actions` worker, agenttap, and the Synadia Agent Web UI all build on `$SRV`.
