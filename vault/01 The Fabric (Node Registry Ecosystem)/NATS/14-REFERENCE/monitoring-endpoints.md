---
title: "Monitoring Endpoints"
section: "14-REFERENCE"
type: note
tags: [nats, monitoring, endpoints, varz, connz, jsz, healthz, http, reference]
source: https://docs.nats.io/running-a-nats-service/nats_admin/monitoring
created: 2026-06-14
---

# Monitoring Endpoints

The server's built-in HTTP/JSON monitoring surface (default port **8222**). Every endpoint, what it returns, and the query params that matter. This is the raw data the exporter and Surveyor consume — and what you `curl` when debugging live.

See also: [[NATS/12-ECOSYSTEM/monitoring-exporters|Monitoring & Exporters]] · [[NATS/07-DEPLOYMENT-OPS/monitoring|Deployment: Monitoring]] · [[NATS/14-REFERENCE/limits-and-defaults|Limits & Defaults]]

---

## Enabling

```text
http: "0.0.0.0:8222"          # plaintext
# or
https: "0.0.0.0:8222"         # TLS (uses the server tls{} or a monitor-specific one)
```

Then everything is `GET http://<host>:8222/<endpoint>`.

---

## The endpoints

| Endpoint | Returns | Key use |
|----------|---------|---------|
| `/healthz` | `{status:"ok"}` (200) or error | **Liveness/readiness probes** (k8s, LBs) |
| `/varz` | Server vitals: version, uptime, mem, cpu, conns, slow_consumers, in/out msgs+bytes | First stop — "how is this server?" |
| `/connz` | Per-connection: cid, ip, subs, pending, in/out, lang, uptime | "Who's connected? Who's slow?" |
| `/routez` | Cluster route connections | Cluster mesh health |
| `/gatewayz` | Gateway (supercluster) connections | Cross-cluster health |
| `/leafz` | Leaf node connections | Edge fleet health |
| `/subsz` | Subscription/sublist stats (counts, cache) | Routing/interest debugging |
| `/jsz` | JetStream: streams, consumers, messages, bytes, replicas, RAFT, account usage | **JetStream health** |
| `/accountz` | Per-account detail (when using accounts) | Multi-tenant inspection |
| `/accstatz` | Per-account stats (conns, msgs) | Per-tenant load |
| `/jszz` /`/healthz?js-enabled-only` | JS-focused health variants | JS liveness |

---

## Useful query params

**`/healthz`**
```
?js-enabled-only=true     # 200 only if JetStream is enabled
?js-server-only=true      # skip per-stream checks (faster)
```

**`/connz`**
```
?subs=true                # include subscription lists per connection
?subs=detail              # subs with stats
?sort=pending             # sort by: pending|msgs_to|bytes_to|subs|cid|uptime|last|idle
?limit=100&offset=0       # paginate
?state=open|closed|all    # include recently closed conns
?cid=<id>                 # one connection
?auth=true                # include auth (user) info
```

**`/jsz`**
```
?streams=true             # include stream details
?consumers=true           # include consumer details
?config=true              # include stream/consumer config
?accounts=true            # per-account JS usage
?account=<pubkey>         # one account
```

**`/subsz`**
```
?subs=true                # list subscriptions (can be large)
```

---

## Quick recipes

```bash
# Server alive?
curl -s localhost:8222/healthz

# JetStream enabled & healthy?
curl -s 'localhost:8222/healthz?js-enabled-only=true'

# Vitals (jq the bits you care about)
curl -s localhost:8222/varz | jq '{conns:.connections, slow:.slow_consumers, mem:.mem, msgs:.in_msgs}'

# Find slow consumers
curl -s 'localhost:8222/connz?sort=pending&limit=10' | jq '.connections[] | {cid,name,pending_bytes,lang}'

# JetStream usage per account
curl -s 'localhost:8222/jsz?accounts=true' | jq '.account_details[]? | {acc:.id, streams:.stream_detail|length, bytes:.store}'

# Leaf nodes connected to this hub
curl -s localhost:8222/leafz | jq '.leafs[] | {name,account,rtt,in_msgs,out_msgs}'
```

---

## Relationship to the tools
- **`prometheus-nats-exporter`** scrapes these endpoints → Prometheus metrics (per server).
- **`nats-surveyor`** uses `$SYS` requests (not just HTTP) to aggregate across the fleet.
- **`nats-top`** reads `/varz` + `/connz` for the live TUI.
- **`nats server report …`** (CLI) uses `$SYS` for cluster-wide views.

See [[NATS/12-ECOSYSTEM/monitoring-exporters|Monitoring & Exporters]] for the full stack.

---

## OneMind Relevance
- **`/healthz`** is the supervision probe for the launchd-managed hub and every `omos-<node>` — wire it into the watchdog. `?js-enabled-only=true` catches "server up but JetStream broken."
- **`/leafz`** is the fleet's edge dashboard primitive — one call shows every connected node, its RTT, and msg flow. Surface it in OneMind Command.
- **`/connz?sort=pending`** is the go-to when an agent is dropping messages (slow consumer hunt).
- Lock the monitoring port down (Tailscale-only, never public) — it leaks connection/topology detail.

## Sources
- Monitoring — https://docs.nats.io/running-a-nats-service/nats_admin/monitoring
- Health checks — https://docs.nats.io/running-a-nats-service/nats_admin/monitoring#health
