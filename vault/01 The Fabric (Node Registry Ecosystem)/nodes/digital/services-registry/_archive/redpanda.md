---
title: "Redpanda"
type: reference
status: active
tags: [platform, infrastructure, messaging, streaming, kafka-compatible]
---
# Redpanda

> Kafka-compatible streaming platform. Single binary, no JVM, no ZooKeeper. 10x faster than Kafka.

## Key Features
| Feature | Notes |
|---------|-------|
| Kafka API | Drop-in Kafka replacement |
| No JVM | C++ — no garbage collection pauses |
| No ZooKeeper | Built-in Raft consensus |
| Single binary | Simple deployment |
| Schema registry | Built-in (Kafka SR compatible) |
| rpk CLI | Powerful management tool |
| Console | Web UI for topics/messages |

## CLI (rpk)
```bash
rpk topic create events
rpk topic list
rpk topic produce events
rpk topic consume events
rpk cluster info
```

## vs NATS
| | Redpanda | NATS |
|---|---|---|
| Model | Log-based (Kafka) | Pub/sub + queue |
| Ordering | Guaranteed per partition | Per subject (JetStream) |
| Replay | Full log replay | JetStream replay |
| Use case | Event sourcing, CDC, heavy streaming | Fleet messaging, request-reply, lightweight |

## Integration
- Event sourcing for audit trails
- CDC (Change Data Capture) from databases
- High-throughput telemetry pipelines
- Complement to NATS (different use cases)
