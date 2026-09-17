---
title: "Advanced — Section Index"
section: "11-ADVANCED"
type: index
tags: [nats, advanced, mqtt, websocket, embedded, benchmarking, performance]
created: 2026-06-14
---

# 11 — Advanced

The capabilities you reach for once the fundamentals are in place: alternative **transports** (MQTT for IoT, WebSocket for browsers), running NATS **embedded** in-process inside a Go binary, and **benchmarking / performance tuning** the bus.

See also: [[NATS/08-CLIENTS-SDKS/README|Clients & SDKs]] · [[NATS/06-TOPOLOGY-SCALING/leaf-nodes|Leaf Nodes]] · [[NATS/07-DEPLOYMENT-OPS/tuning-and-dr|Tuning & DR]]

## Notes
- [[NATS/11-ADVANCED/mqtt|MQTT Support]] — speak MQTT 3.1.1 directly to NATS; QoS 0/1 mapped onto JetStream; for IoT devices and existing MQTT fleets
- [[NATS/11-ADVANCED/websocket|WebSocket]] — browser-native NATS over WS/WSS; powers dashboards and web clients with no gateway shim
- [[NATS/11-ADVANCED/embedded-nats|Embedded NATS]] — run a full server in-process from Go (`server.NewServer`); zero-dependency edge nodes and tests
- [[NATS/11-ADVANCED/benchmarking|Benchmarking & Performance]] — `nats bench`, latency vs throughput, slow consumers, tuning levers

## When you need this section
| Need | Go to |
|------|-------|
| Connect an ESP32 / existing MQTT sensor fleet | [[NATS/11-ADVANCED/mqtt\|MQTT]] |
| A browser dashboard subscribing to live subjects | [[NATS/11-ADVANCED/websocket\|WebSocket]] |
| Ship one Go binary that *is* the broker (edge node, CLI, test) | [[NATS/11-ADVANCED/embedded-nats\|Embedded]] |
| Prove throughput/latency, find the bottleneck | [[NATS/11-ADVANCED/benchmarking\|Benchmarking]] |

## OneMind Relevance
This section is **disproportionately important** for the fleet: MQTT brings in off-the-shelf IoT hardware, WebSocket powers the OneMind Command dashboards, and **embedded NATS** is how edge nodes (`omd --embed`) run a broker with zero extra processes. These three are the edge story.

## Sources
- MQTT — https://docs.nats.io/running-a-nats-service/configuration/mqtt
- WebSocket — https://docs.nats.io/running-a-nats-service/configuration/websocket
- Embedding — https://docs.nats.io/using-nats/developer/embedding
- Benchmarking — https://docs.nats.io/using-nats/nats-tools/nats_cli/natsbench
