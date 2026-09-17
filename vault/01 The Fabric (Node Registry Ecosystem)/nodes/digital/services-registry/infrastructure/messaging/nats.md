---
title: "NATS"
type: reference
status: active
tags: [platform, infrastructure, messaging, backbone]
---
# NATS

> Message backbone for OneMind OS. All fleet communication, event routing, telemetry, and AI experience storage.

## Ports
- **4222** — Client connections
- **8222** — HTTP monitoring
- **7422** — Cluster/leaf routes
- **1883** — MQTT bridge

## Key Concepts
- **JetStream** — Persistent streams (7d hot retention), KV buckets for state
- **Subjects** — Hierarchical routing (`cmd.ha.*`, `fleet.>`, `telemetry.>`)
- **Request-Reply** — Synchronous SaaS action execution via NATS

## CLI
```bash
nats pub subject "message"
nats sub "fleet.>"
nats stream ls
nats kv ls
nats server info
```

## OneMind Subjects
| Pattern | Purpose |
|---------|---------|
| `cmd.ha.*` | Home Assistant commands |
| `fleet.>` | Agent fleet communication |
| `telemetry.>` | Sensor/device telemetry |
| `saas.request` | SaaS action execution |
| `$SRV` | Fleet citizen service discovery |
