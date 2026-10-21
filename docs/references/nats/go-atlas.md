# go-atlas — NATS Go Toolkit Reference

**Repo:** https://github.com/altessa-s/go-atlas
**License:** MIT
**Install:** `go get github.com/altessa-s/go-atlas`
**Requires:** Go 1.25+

## What it is

Modular Go toolkit for distributed services. Each package covers one concern and is usable standalone — no framework bootstrap, no global state, interface-driven with dependency injection.

## Why it's relevant to OneMind

go-atlas includes a **NATS transport module** alongside auth, observability, secrets, caching, and config. For building OneMind Go services (`om-control-plane`, `om-router`, `om-connector-manager`) it covers the common boilerplate so each service only implements its domain logic.

## Modules relevant to OneMind

| Module | Use |
|--------|-----|
| `transport/nats` | NATS pub/sub, request/reply — connects services to the bus |
| `observability` | Prometheus metrics + tracing (already using Grafana/Tempo on DOKS) |
| `auth` | JWT/token validation for protected command subjects |
| `secrets` | Secret management (pairs with 1Password/K8s secrets) |
| `config` | Structured config loading with environment override |

## Status

Evaluated as a potential dependency for OneMind Go microservices. Not yet imported.
Use `go get` when building the first Go service that needs NATS transport — do not vendor/copy.

## Notes

- Interfaces are stable enough for production use
- Import only what you need — no monolithic bootstrap
- Complements (does not replace) the raw NATS Go client (`github.com/nats-io/nats.go`)
