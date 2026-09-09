# orbit.go — Synadia NATS Go Extensions

**Repo:** https://github.com/synadia-io/orbit.go
**Publisher:** Synadia (makers of NATS)
**License:** Apache 2.0
**Install:** `go get github.com/synadia-io/orbit.go/<module>`

## What it is

Official Synadia toolkit of independent Go utility modules built on top of NATS. Pick only what you need — no monolithic import.

## Modules relevant to OneMind

| Module | Install | Use in OneMind |
|--------|---------|----------------|
| `natssysclient` | `go get github.com/synadia-io/orbit.go/natssysclient` | Query NATS server health, connections, JetStream stats programmatically from Go services |
| `jetstreamext` | `go get github.com/synadia-io/orbit.go/jetstreamext` | Fast batch publishing to JetStream — useful for high-volume feed ingestion (AIS, ADS-B, sensors) |
| `pcgroups` | `go get github.com/synadia-io/orbit.go/pcgroups` | Partitioned consumer groups — parallel processing of streams across multiple service instances |
| `kvcodec` | `go get github.com/synadia-io/orbit.go/kvcodec` | Transparent encode/decode for JetStream KV — cleaner access to registry buckets (`KV_ENTITY_REGISTRY`, `KV_SERVICE_REGISTRY`, etc.) |
| `natsext` | `go get github.com/synadia-io/orbit.go/natsext` | `RequestMany` pattern — fan-out requests to multiple NATS services and collect responses |
| `counters` | `go get github.com/synadia-io/orbit.go/counters` | Distributed counters on JetStream — useful for metrics, rate limiting, event counts |

## Status

Evaluated as utility dependencies for OneMind Go services. Not yet imported.
API not guaranteed stable until v1.0.0 per module — check before adopting in production.

## Notes

- All modules are independent — import only what you need
- Official Synadia project, actively maintained alongside NATS Server
- Complements the raw NATS Go client (`github.com/nats-io/nats.go`)
