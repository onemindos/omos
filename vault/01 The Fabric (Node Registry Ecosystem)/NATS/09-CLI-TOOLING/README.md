---
title: "CLI & Tooling — Section Index"
section: "09-CLI-TOOLING"
type: index
tags: [nats, cli, nsc, nk, nats-top, bench, tooling]
created: 2026-06-14
---

# 09 — CLI & Tooling

The command-line tools every NATS operator lives in.

## Notes
- [[NATS/09-CLI-TOOLING/nats-cli|The `nats` CLI]] — the swiss-army tool: pub/sub, req, streams, consumers, kv, object, server, bench, contexts
- [[NATS/09-CLI-TOOLING/contexts|Contexts]] — saved connection profiles (`nats context`)
- [[NATS/09-CLI-TOOLING/nsc-nk|`nsc` & `nk`]] — security keys & JWT management
- [[NATS/09-CLI-TOOLING/other-tools|nats-top, bench, exporters]] — live monitoring TUI, benchmarking, Prometheus

## The toolbelt
| Tool | Repo | Purpose |
|------|------|---------|
| **`nats`** | `nats-io/natscli` | Primary CLI — everything |
| **`nsc`** | `nats-io/nsc` | Operator/account/user JWT + NKEY management |
| **`nk`** | `nats-io/nkeys` | Generate/inspect NKEYS |
| **`nats-top`** | `nats-io/nats-top` | `top`-style live server view |
| **`nats bench`** | (in `nats`) | Throughput/latency benchmarking |
| **prometheus-nats-exporter** | `nats-io/prometheus-nats-exporter` | Metrics for Prometheus |
| **nats-surveyor** | `nats-io/nats-surveyor` | Fleet-wide monitoring via `$SYS` |

## See also
- OneMind keeps a full CLI cheatsheet in the `nats-server` ops skill; [[NATS/14-REFERENCE/cli-cheatsheet|Reference: CLI cheatsheet]] distills it here.
