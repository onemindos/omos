---
title: "Services / Micro Framework"
section: "04-SERVICES-MICRO"
type: note
tags: [nats, micro, services, discovery, srv, observability]
source: https://docs.nats.io/using-nats/developer/services
created: 2026-06-14
---

# Services / Micro Framework

A built-in framework for writing **request-reply microservices** on NATS with automatic **discovery**, **observability**, and **versioning** — standardized across client languages.

See also: [[NATS/01-CORE-CONCEPTS/request-reply|Request/Reply]] · [[NATS/01-CORE-CONCEPTS/queue-groups|Queue Groups]]

---

## Why Micro Exists

You *can* build a service with plain `subscribe` + reply on a subject. But every team reinvents: how do I discover what services exist? how do I report stats? how do I version? Micro standardizes this. A service built with the Micro API automatically:
- **Registers** itself for discovery on the reserved **`$SRV.>`** subjects.
- **Load-balances** across instances (endpoints subscribe in a queue group — run N copies, requests spread automatically).
- **Reports stats** — request count, error count, processing time — with zero extra code.
- **Exposes schema/version** — clients can introspect.

It's a thin convention over request-reply, not a new protocol.

---

## Anatomy

```
Service "calc" v1.0.0
├── Group "math"
│   ├── Endpoint "add"      → subject math.add
│   └── Endpoint "multiply" → subject math.multiply
└── Endpoint "ping"          → subject ping
```

- **Service** — named, versioned unit (name + version + description).
- **Group** — a subject-prefix namespace for endpoints (optional).
- **Endpoint** — a handler bound to a subject; receives requests, sends replies.

---

## Discovery — the `$SRV` API

Every Micro service answers three reserved request subjects (all support `$SRV.<verb>`, `$SRV.<verb>.<name>`, `$SRV.<verb>.<name>.<id>`):

| Subject | Returns |
|---------|---------|
| **`$SRV.PING`** | Liveness — every service replies with id/name/version (service discovery) |
| **`$SRV.INFO`** | Full info — endpoints, subjects, metadata, description |
| **`$SRV.STATS`** | Per-endpoint stats — num_requests, num_errors, processing_time, average_processing_time |

```shell
nats micro ls               # list all services on the bus
nats micro info calc        # full info for a service
nats micro stats calc       # live request/error/latency stats
nats req '$SRV.PING' ''     # raw discovery ping
```

> This is how OneMind's `fabric_discover` finds agents: it queries `$SRV.INFO` and reads each service's metadata (token, owner, session) to locate the `prompt` endpoint subject. The SaaS action worker (`saas-actions`) registers here too — `nats micro ls` shows it.

---

## Client API (Python / Go)

```python
# Python (nats-py micro)
from nats.micro import add_service
from nats.micro.service import GroupConfig, EndpointConfig

async def add_handler(req):
    a, b = parse(req.data)
    await req.respond(str(a + b).encode())

srv = await add_service(nc, name="calc", version="1.0.0", description="math service")
group = srv.add_group(GroupConfig(name="math"))
await group.add_endpoint(EndpointConfig(name="add", handler=add_handler))
# endpoint subject = math.add ; stats auto-collected
```

```go
// Go (micro)
srv, _ := micro.AddService(nc, micro.Config{Name: "calc", Version: "1.0.0"})
grp := srv.AddGroup("math")
grp.AddEndpoint("add", micro.HandlerFunc(func(r micro.Request) {
    r.Respond([]byte(add(r.Data())))
}))
```

> **OneMind pitfall (verified, nats-py 2.14):** `srv.add_group("name")` errors — pass a config object: `srv.add_group(GroupConfig(name="math"))` and `group.add_endpoint(EndpointConfig(name="add", handler=fn))`. Import from `nats.micro.service`.

---

## Observability for Free

Because every request flows through the framework, Micro tallies per-endpoint **num_requests**, **num_errors**, **last_error**, **processing_time** (total), and **average_processing_time** — queryable live via `$SRV.STATS`. Tools like **agenttap** ("Wireshark for AI agents") and the Synadia Agent Web UI read these to visualize the service mesh. You get a service dashboard without instrumenting anything.

---

## Versioning & Load Balancing

- **Versioning** — the version string is part of discovery; run v1 and v2 side by side, route by inspecting `$SRV.INFO`.
- **Load balancing** — endpoints subscribe in a **queue group** named for the service, so launching multiple instances spreads requests automatically (same mechanism as [[NATS/01-CORE-CONCEPTS/queue-groups|Queue Groups]]). Scale a service by starting more copies — no config.

---

## When to Use Micro vs Raw Request-Reply

| Use **Micro** when | Use **raw req-reply** when |
|--------------------|----------------------------|
| You want discovery + stats + versioning | One-off / internal reply, no discovery needed |
| Multiple instances, want auto load-balance + introspection | Ultra-minimal footprint, you control both ends |
| Building a service mesh others will discover | Quick RPC inside one codebase |

## Sources
- Services framework — https://docs.nats.io/using-nats/developer/services
- Micro (Go) — https://pkg.go.dev/github.com/nats-io/nats.go/micro
- Service API reference — https://docs.nats.io/reference/reference-protocols/nats-service-api
- `$SRV` discovery — https://docs.nats.io/using-nats/developer/services#discovery
