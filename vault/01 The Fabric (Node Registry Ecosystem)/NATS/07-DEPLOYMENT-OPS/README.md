---
title: "Deployment & Ops — Section Index"
section: "07-DEPLOYMENT-OPS"
type: index
tags: [nats, deployment, operations, config, monitoring, kubernetes]
created: 2026-06-14
---

# 07 — Deployment & Operations

Running NATS in the real world: install, configure, monitor, deploy on Kubernetes, tune, back up, and upgrade.

## Notes
- [[NATS/07-DEPLOYMENT-OPS/installation|Installation]] — binary, Docker, brew, package managers, service files
- [[NATS/07-DEPLOYMENT-OPS/configuration-reference|Configuration Reference]] — the full nats-server.conf surface
- [[NATS/07-DEPLOYMENT-OPS/monitoring|Monitoring]] — the `:8222` HTTP endpoints (/varz /connz /jsz /healthz …), Prometheus
- [[NATS/07-DEPLOYMENT-OPS/kubernetes|Kubernetes]] — Helm chart, NACK (controllers), operator patterns
- [[NATS/07-DEPLOYMENT-OPS/tuning-and-dr|Tuning, Backup & DR]] — performance tuning, limits, backups, upgrades, signals

## See also
- OneMind runs nats-server via **launchd** (binary install, no Docker) — see the `nats-server` ops skill + `~/onemindos-core/docs/infra/NATS-SETUP.md`.
