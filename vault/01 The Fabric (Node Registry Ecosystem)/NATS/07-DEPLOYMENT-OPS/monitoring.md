---
title: "Monitoring"
section: "07-DEPLOYMENT-OPS"
type: note
tags: [nats, monitoring, varz, connz, jsz, healthz, prometheus, observability]
source: https://docs.nats.io/running-a-nats-service/nats_admin/monitoring
created: 2026-06-14
---

# Monitoring

NATS exposes a built-in **HTTP monitoring** endpoint (default **`:8222`**) returning JSON about every aspect of the server. No agent required — just `curl`.

See also: [[NATS/09-CLI-TOOLING/README|CLI & Tooling]] · [[NATS/07-DEPLOYMENT-OPS/configuration-reference|Configuration Reference]]

---

## Enabling

```text
http_port: 8222           # plain HTTP
# or
https_port: 8222          # HTTPS (needs tls{})
```
Then browse/curl `http://localhost:8222/`.

---

## The Endpoints

| Endpoint | Returns |
|----------|---------|
| **`/varz`** | General server state: version, uptime, host/port, connection counts, msg/byte totals, mem/cpu, config flags, slow consumers |
| **`/connz`** | Per-connection detail: cid, ip, subscriptions, pending, msgs in/out, lang/version. Supports `?subs=1`, `?limit=`, `?sort=` |
| **`/subsz`** | Subscription routing table (the sublist); `?subs=1` for detail |
| **`/routez`** | Cluster route connections + state |
| **`/gatewayz`** | Gateway (supercluster) connections |
| **`/leafz`** | Leaf node connections |
| **`/jsz`** | **JetStream** state: streams, consumers, messages, bytes, memory/disk, cluster/RAFT info, per-account. `?streams=1&consumers=1&config=1` |
| **`/accountz`** | Accounts and their state |
| **`/accstatz`** | Per-account stats |
| **`/healthz`** | **Health check** — `{"status":"ok"}` (200) or error. Use for liveness/readiness probes. `?js-enabled=true`, `?js-server-only=true` variants |
| **`/`** | Index of all endpoints |

```shell
curl -s localhost:8222/varz | jq '.connections, .in_msgs, .out_msgs, .slow_consumers'
curl -s localhost:8222/jsz?streams=1 | jq '.streams, .messages, .bytes'
curl -s localhost:8222/healthz                 # liveness probe
curl -s 'localhost:8222/connz?subs=1&limit=5'  # who's connected + their subs
```

`/varz` sample fields: `server_id`, `version`, `max_payload` (1048576 = 1MB default), `connections`, `in_msgs/out_msgs`, `in_bytes/out_bytes`, `slow_consumers`, `mem`, `cpu`, `jetstream`.

---

## `/healthz` for Probes

```yaml
# Kubernetes liveness/readiness
livenessProbe:  { httpGet: { path: /healthz, port: 8222 } }
readinessProbe: { httpGet: { path: "/healthz?js-enabled-only=true", port: 8222 } }
```
Returns HTTP 200 + `{"status":"ok"}` when healthy. JetStream-aware variants check stream/consumer/RAFT health, not just process liveness.

---

## Prometheus

The server doesn't emit Prometheus format natively — use the **prometheus-nats-exporter** (`nats-io/prometheus-nats-exporter`), which scrapes the monitoring endpoints and exposes Prometheus metrics:

```shell
prometheus-nats-exporter -varz -connz -jsz -healthz http://localhost:8222
# exposes :7777/metrics for Prometheus to scrape
```
Pair with Grafana dashboards. For higher-level fleet observability, **NATS Surveyor** (`nats-io/nats-surveyor`) polls the system account across a cluster and exports rich metrics. See [[NATS/12-ECOSYSTEM/README|Ecosystem]].

---

## CLI Monitoring

The `nats` CLI wraps much of this (see [[NATS/09-CLI-TOOLING/README|CLI]]):
```shell
nats server info          # server snapshot
nats server report jetstream   # JetStream report (needs SYS creds!)
nats server report connections
nats server ls            # list servers (system account)
nats server check         # health checks (great for alerting)
nats stream report        # all streams: msgs, bytes, consumers
nats-top                  # live top-style TUI (separate tool)
```

> **Pitfall (verified):** `nats server report jetstream` needs **system-account** credentials, not regular bus creds. With normal creds it returns nothing/errors.

---

## What to Watch (alerting)

| Signal | Why |
|--------|-----|
| **slow_consumers** (`/varz`) | Clients can't keep up → message loss |
| **`/healthz` != 200** | Server/JetStream unhealthy |
| **JetStream disk/mem vs limits** (`/jsz`) | Approaching `max_file_store` → publishes will fail (DiscardNew) or evict (DiscardOld) |
| **No RAFT leader / under-replicated** (`/jsz`) | JetStream HA degraded |
| **connections trend** (`/varz`) | Capacity / leak detection |
| **route/leaf/gateway down** | Topology partition |

OneMind exposes this at the bus level; the Flutter `NatsMonitorService` and the live bus viewer read these endpoints + `$SYS` events.

## Sources
- Monitoring — https://docs.nats.io/running-a-nats-service/nats_admin/monitoring
- `/healthz` — https://docs.nats.io/running-a-nats-service/nats_admin/monitoring#health
- Prometheus exporter — https://github.com/nats-io/prometheus-nats-exporter
- Surveyor — https://github.com/nats-io/nats-surveyor
