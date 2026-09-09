# logtrace-system — NATS JetStream + Full OTel Observability Stack

**Language:** Go  
**Stack:** NATS JetStream · Loki · Grafana · Jaeger · OpenTelemetry  
**Role in OneMind:** Reference implementation for distributed tracing and log aggregation across NATS microservices

Comprehensive logging and tracing system for microservices using NATS JetStream as the event bus, with full OpenTelemetry instrumentation feeding Loki (logs), Grafana (dashboards), and Jaeger (distributed traces).

## Stack

| Component | Role |
|-----------|------|
| NATS JetStream | Event bus — services publish logs/traces as NATS messages |
| OpenTelemetry | Instrumentation SDK across all services |
| Loki | Log aggregation and storage |
| Grafana | Unified dashboards (logs + metrics + traces) |
| Jaeger | Distributed trace visualization |

## Why it matters for OneMind

OneMind's omos-services publish telemetry to NATS. This stack closes the observability loop:

```
omos-services (OTel instrumented)
    → tel.{service}.logs / tel.{service}.traces (NATS)
    → NATS JetStream consumer
    → Loki (logs) + Jaeger (traces)
    → Grafana unified view
```

Complements [[otel-edge]] which handles IoT→OTLP. This is the service-side observability path.

## See Also

- [[otel-edge]] — NATS/MQTT → OpenTelemetry for IoT edge nodes
- `NATS-SETUP.md` in omos-core docs/infra — live NATS server with stream/KV definitions
