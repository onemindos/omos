# otel-edge — NATS/MQTT → OpenTelemetry Agent

**Language:** Rust (embedded-Linux-first)  
**Role in OneMind:** Bridge IoT telemetry from NATS/MQTT into the OpenTelemetry pipeline

Small embedded-Linux agent that ingests telemetry over MQTT and NATS, maps IoT topic/payload semantics onto OpenTelemetry data model, buffers through power/network loss, and emits OTLP (logs, metrics, traces).

## Why it matters for OneMind

OneMind telemetry currently lands in ClickHouse via NATS. otel-edge adds a second path:

```
NATS FABRIC tel.* subjects
        ↓
    otel-edge
        ↓ OTLP
  Grafana / Prometheus / Jaeger
```

Power-loss buffering means no telemetry gaps during outages — critical for homestead/field deployments.

## Integration point

otel-edge subscribes to `tel.>` on the NATS FABRIC account, maps to OTLP metrics, and emits to the observability stack. Runs on any embedded Linux (Raspberry Pi, Jetson Nano, etc.).
