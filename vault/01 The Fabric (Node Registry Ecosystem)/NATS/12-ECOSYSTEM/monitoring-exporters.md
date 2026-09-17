---
title: "Monitoring & Exporters"
section: "12-ECOSYSTEM"
type: note
tags: [nats, monitoring, surveyor, prometheus, exporter, nats-top, grafana, observability]
source: https://github.com/nats-io/nats-surveyor
created: 2026-06-14
---

# Monitoring & Exporters

The OSS observability stack for NATS: the server's built-in HTTP endpoints, scraped by the **Prometheus exporter** or aggregated fleet-wide by **NATS Surveyor**, plus the **`nats-top`** live TUI. This is what feeds Grafana — and what OneMind already runs.

See also: [[NATS/07-DEPLOYMENT-OPS/monitoring|Deployment: Monitoring]] · [[NATS/12-ECOSYSTEM/repo-map|Repo Map]] · [[NATS/09-CLI-TOOLING/README|CLI & Tooling]]

---

## The foundation: built-in monitoring endpoints (:8222)

Every `nats-server` exposes an HTTP/JSON monitoring port (default **8222**) when enabled:

```text
http: 0.0.0.0:8222
```

| Endpoint | Reports |
|----------|---------|
| `/varz` | Server vitals: version, uptime, mem, CPU, connections, slow consumers |
| `/connz` | Per-connection detail (subs, pending, in/out bytes) |
| `/subsz` | Subscription/sublist stats |
| `/routez` | Cluster routes |
| `/gatewayz` | Supercluster gateways |
| `/leafz` | Leaf node connections |
| `/jsz` | JetStream: streams, consumers, bytes, replicas, RAFT |
| `/accountz` | Per-account detail |
| `/healthz` | Liveness/readiness (great for k8s probes & load balancers) |

Everything else here is built on these endpoints. See the full treatment in [[NATS/07-DEPLOYMENT-OPS/monitoring|Section 07: Monitoring]].

---

## `prometheus-nats-exporter`

Repo: `nats-io/prometheus-nats-exporter` (★487, Go). Scrapes the endpoints above and re-exposes them as **Prometheus metrics**.

```bash
prometheus-nats-exporter -varz -connz -jsz -leafz http://localhost:8222
# → exposes :7777/metrics for Prometheus to scrape
```

Pick which collectors to enable (`-varz`, `-jsz`, etc.). Point Prometheus at it, build Grafana dashboards on the resulting `gnatsd_*` / JetStream metrics. This is the **single-server / per-server** path.

---

## `nats-surveyor` — fleet-wide

Repo: `nats-io/nats-surveyor` (★317, Go). Where the exporter watches one server, **Surveyor** watches a whole **cluster/fleet** by subscribing to the **system account (`$SYS`)** and issuing service-discovery requests. It aggregates metrics across all servers and exports to Prometheus, with ready-made Grafana dashboards.

```bash
nats-surveyor -s nats://localhost:4222 \
  --creds /etc/nats/sys.creds \      # needs $SYS account access
  --count 3                          # expected number of servers
```

Surveyor is the right tool once you have **more than one server** — it gives a single pane across the cluster instead of N separate exporters.

> **OneMind runs this**: the `io.nats.surveyor` launchd agent (`~/Library/LaunchAgents/io.nats.surveyor.plist`, `go/bin/nats-surveyor`) is our fleet monitor. As the fleet grows past the single Mac hub, Surveyor + `$SYS` is how we keep one view.

---

## `nats-top` — the live TUI

Repo: `nats-io/nats-top` (★402, Go). A `top`-like terminal UI for a server — live connections, subscriptions, msg rates, slow consumers. Zero setup, perfect for "what's happening *right now*."

```bash
nats-top -s localhost -m 8222        # -m = monitoring port
```

No Prometheus needed — it reads `/varz` + `/connz` directly. The operator's first reflex when something feels off.

---

## Distributed tracing — `not.go` / `not.java`

Repos: `nats-io/not.go` (★104), `nats-io/not.java` (★2). **N**ATS **O**pen **T**racing references — patterns for propagating trace context (e.g. via message headers) so a request crossing services over NATS shows up as one trace in Jaeger/Tempo. Reference implementations, not turnkey products.

---

## Putting it together

```
                ┌──────────── Grafana ────────────┐
                │   dashboards (per-server + fleet) │
                └───────────────▲──────────────────┘
                                │ PromQL
                        ┌───────┴────────┐
                        │   Prometheus    │
                        └──▲──────────▲───┘
            scrape /metrics│          │scrape /metrics
   ┌────────────────────────┐   ┌─────┴───────────────┐
   │ prometheus-nats-exporter│   │   nats-surveyor      │
   │   (per server :8222)    │   │  ($SYS, whole fleet) │
   └─────────▲───────────────┘   └─────────▲────────────┘
             │ /varz /jsz …                 │ $SYS requests
        ┌────┴─────┐                   ┌─────┴─────┐
        │ server A │ ───── cluster ────│ server B  │
        └──────────┘                   └───────────┘

   nats-top → ad-hoc live view of any single server
```

---

## OneMind Relevance
- **Already wired:** `nats-surveyor` (launchd) + the exporter feed our metrics; Langfuse/Grafana-style dashboards visualize them. `/healthz` is the probe for service supervision.
- **Scaling path:** single-Mac hub → exporter is enough; the moment a second hub or leaf cluster appears, Surveyor on `$SYS` becomes the fleet view.
- **Tracing:** if cross-agent request tracing becomes a need, `not.go` headers are the pattern to copy.

## Sources
- nats-surveyor — https://github.com/nats-io/nats-surveyor
- prometheus-nats-exporter — https://github.com/nats-io/prometheus-nats-exporter
- nats-top — https://github.com/nats-io/nats-top
- Monitoring docs — https://docs.nats.io/running-a-nats-service/nats_admin/monitoring
