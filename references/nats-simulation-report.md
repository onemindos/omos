# NATS.io Simulation-Related Repositories and Projects

**Research Date:** June 18, 2026  
**Compiled by:** Hermes Agent

This report documents open-source GitHub repositories and projects that combine NATS messaging with simulation systems across various domains.

---

## Key Findings Summary

- **Total repositories identified:** 15+ direct simulation-related projects
- **Most active domain:** Digital twins, IoT simulators, and multi-agent systems
- **Common languages:** Go, Python, TypeScript, C++
- **Production readiness:** Mix of production-ready tools and research prototypes

---

## 1. Digital Twins and Virtual Simulations

### UNS-Design-Studio
- **Repository:** https://github.com/ilja0101/UNS-Design-Studio
- **Stars:** 18
- **Language:** Python
- **Description:** Design, test, validate and demo your Unified Namespace architecture: topic trees, payload schemas and stateful simulation, without touching a single physical system.
- **Status:** Production-ready tool for UNS/ISA-95 architecture simulation
- **Last Updated:** June 2026

### go-digitaltwin/v2-experiment
- **Repository:** https://github.com/go-digitaltwin/v2-experiment
- **Stars:** 1
- **Language:** Go
- **Description:** Analytics-enabled digital-twins based on messaging patterns with NATS
- **Status:** Research/Prototype
- **Last Updated:** June 2026

---

## 2. Robot Simulators and Robotics

### eranet/rhombus
- **Repository:** https://github.com/eranet/rhombus
- **Stars:** 6
- **Language:** C++
- **Description:** ROS competitor, server for robotics microservices using NATS
- **Status:** Research/Prototype
- **Last Updated:** March 2023

---

## 3. IoT and Sensor Simulators

### rb-sensor-simulator
- **Repository:** https://github.com/robertobouses/rb-sensor-simulator
- **Stars:** 2
- **Language:** Go
- **Description:** IoT sensor simulator built in Go. Includes configurable sensor behavior, NATS-based messaging, and data persistence.
- **Status:** Production-ready for IoT testing
- **Last Updated:** June 2025

### go-nats-simulator
- **Repository:** https://github.com/AntonioBR9998/go-nats-simulator
- **Stars:** 0
- **Language:** Go
- **Description:** Small microservices architecture in Golang to manage a network of IoT devices and simulate periodic sending of events with NATS and TimescaleDB
- **Status:** Development/Testing

### unitedrhino/things
- **Repository:** https://github.com/unitedrhino/things
- **Stars:** 628
- **Language:** Go
- **Description:** SaaS cloud-native microservices IoT platform built on golang and gozero framework, supporting MQTT, CoAP, HTTP, TCP, Aliyun, Tencent Cloud, Tuya Cloud integration
- **Status:** Enterprise-grade IoT platform, Production-ready

---

## 4. Vehicle and Transit Simulators

### gtfs-simulator-go
- **Repository:** https://github.com/ponytojas/gtfs-simulator-go
- **Stars:** 0
- **Language:** Go
- **Description:** Golang service that uses GTFS data to simulate active vehicles and send data into NATS
- **Status:** Development/Testing
- **Last Updated:** November 2025

---

## 5. Multi-Agent Systems and AI Agents

### synadia-ai/synadia-agents
- **Repository:** https://github.com/synadia-ai/synadia-agents
- **Stars:** 58
- **Language:** TypeScript
- **Description:** SDKs and reference implementations for the Synadia Agent Protocol for NATS
- **Status:** Production framework for multi-agent systems
- **Last Updated:** June 2026

### cordum-io/cordum
- **Repository:** https://github.com/cordum-io/cordum
- **Stars:** 483
- **Language:** Go
- **Description:** Open agent control plane. Govern autonomous AI agents with pre-execution policy enforcement, approval gates, and audit trails. Works with LangChain, CrewAI, MCP
- **Status:** Enterprise AI agent governance, Production-ready

---

## 6. Test Harnesses and Mock Systems

### nats-memory-server
- **Repository:** https://github.com/ll1r1k-1337/nats-memory-server
- **Stars:** 7
- **Language:** TypeScript
- **Description:** Node.js package that provides an in-memory NATS server for testing
- **Status:** Production-ready testing utility
- **Last Updated:** June 2026

### PopulateEventToNats
- **Repository:** https://github.com/dharriso/PopulateEventToNats
- **Stars:** 0
- **Language:** Java
- **Description:** Test Harness to Publish Event History Message to Partitioned NATS Streams
- **Status:** Internal testing tool

---

## 7. Game/MMO Simulators and Virtual Environments

### bobohume/gonet
- **Repository:** https://github.com/bobohume/gonet
- **Stars:** 1,609
- **Language:** Go
- **Description:** Go distributed server, memory-based MMO simulation
- **Status:** Production-ready game server framework

---

## 8. Event-Driven and Discrete Event Simulation Frameworks

### ThreeDotsLabs/watermill
- **Repository:** https://github.com/ThreeDotsLabs/watermill
- **Stars:** 9,762
- **Language:** Go
- **Description:** Building event-driven applications the easy way in Go
- **Status:** Mature production framework with NATS support

### looplab/eventhorizon
- **Repository:** https://github.com/looplab/eventhorizon
- **Stars:** 1,670
- **Language:** Go
- **Description:** Event Sourcing framework for Go with NATS support
- **Status:** Production-ready CQRS and event sourcing

---

## 9. Industrial and Edge Computing

### wasmCloud/wasmCloud
- **Repository:** https://github.com/wasmCloud/wasmCloud
- **Stars:** 2,347
- **Language:** Rust
- **Description:** CNCF project for building polyglot apps across any cloud, K8s, or edge. Uses NATS as the lattice backbone.
- **Status:** Enterprise-grade, Production-ready

---

## 10. Network and Stream Processing

### liftbridge-io/liftbridge
- **Repository:** https://github.com/liftbridge-io/liftbridge
- **Stars:** 2,789
- **Language:** Go
- **Description:** Kafka-style message streaming built on NATS
- **Status:** Production-ready

### arkflow-rs/arkflow
- **Repository:** https://github.com/arkflow-rs/arkflow
- **Stars:** 1,283
- **Language:** Rust
- **Description:** High performance Rust stream processing engine with AI capabilities
- **Status:** Development/Beta

---

## 11. Additional NATS Ecosystem Tools

### ag2ai/faststream
- **Repository:** https://github.com/ag2ai/faststream
- **Stars:** 5,238
- **Language:** Python
- **Description:** Asynchronous Python framework for building event-driven applications
- **Status:** Production-ready

### shijuvar/go-distributed-sys
- **Repository:** https://github.com/shijuvar/go-distributed-sys
- **Stars:** 665
- **Language:** Go
- **Description:** Guidance for building event-driven distributed systems with NATS JetStream
- **Status:** Reference implementation

### openmessaging/benchmark
- **Repository:** https://github.com/openmessaging/benchmark
- **Stars:** 449
- **Language:** Java
- **Description:** OpenMessaging Benchmark Framework - includes NATS benchmarking
- **Status:** Production tool

---

## Domain Coverage Analysis

| Domain | Found Projects | Maturity Level |
|--------|----------------|----------------|
| Digital Twins | 2 | Production + Research |
| Robot Simulators | 1 | Research |
| IoT Sensors | 3 | Production-ready |
| Vehicle Simulators | 1 | Development |
| Flight Simulators | 0 | None found |
| Physics Simulators | 0 | None found |
| Network Simulators | 0 | None found |
| Multi-Agent Systems | 2 | Production |
| Test Harnesses | 2 | Production |
| Virtual Environments | 1 | Production |
| Event-Driven Systems | 3+ | Production |
| Industrial/Edge | 2 | Production |

---

## Notable Gaps and Opportunities

### Domains with NO dedicated NATS simulation projects:
1. Flight simulators (FlightGear, X-Plane, CARLA integration)
2. Physics simulators (Bullet, PhysX, MuJoCo)
3. Network simulators (ns-3, CORE, OMNeT++)
4. 3D simulators (Gazebo, Isaac Sim with NATS bridge)
5. Metaverse/Virtual worlds (dedicated NATS-based)

---

## Production Use Cases Identified

1. Smart Manufacturing: UNS-Design-Studio for ISA-95 architectures
2. IoT Testing: Sensor simulators for development and QA
3. Multi-Agent AI: Synadia agents for distributed AI systems
4. Edge Computing: wasmCloud for edge simulation and testing
5. Game Servers: gonet for MMO simulation
6. Public Transit: GTFS simulator for transit planning
7. Agent Governance: Cordum for AI agent control

---

## Quick Reference: Fork-Ready URLs

### Production-Ready Simulation Projects:
- https://github.com/ilja0101/UNS-Design-Studio
- https://github.com/robertobouses/rb-sensor-simulator
- https://github.com/synadia-ai/synadia-agents
- https://github.com/ll1r1k-1337/nats-memory-server
- https://github.com/unitedrhino/things
- https://github.com/cordum-io/cordum

### Research/Prototype Projects:
- https://github.com/go-digitaltwin/v2-experiment
- https://github.com/eranet/rhombus
- https://github.com/ponytojas/gtfs-simulator-go
- https://github.com/AntonioBR9998/go-nats-simulator

### Supporting Frameworks:
- https://github.com/ThreeDotsLabs/watermill
- https://github.com/looplab/eventhorizon
- https://github.com/wasmCloud/wasmCloud
- https://github.com/liftbridge-io/liftbridge
- https://github.com/ag2ai/faststream
