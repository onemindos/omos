# ECLSS Digital Twin

**Language:** TypeScript/Node.js  
**Stack:** NATS 2.10 as event backbone  
**Role in OneMind:** Reference architecture for isolated habitat monitoring digital twin

Fully simulated, microservice-based digital twin of an Environmental Control and Life Support System (ECLSS) for isolated habitats. Autonomous power management with dynamic load shedding, critical life support subsystems, real-time monitoring dashboard.

## Architecture

- NATS (or RabbitMQ) as high-throughput pub/sub backbone
- WebSocket bridge: NATS → browser (port 9090)
- Simulates: power, atmosphere, water, thermal control, fire suppression subsystems
- Dynamic load shedding when power is constrained

## Why it matters for OneMind

Reference implementation for homestead environmental monitoring. The event subjects and subsystem model map directly onto OneMind's telemetry grammar:

```
tel.homestead.env.temperature
tel.homestead.env.co2
tel.homestead.power.load
evt.homestead.power.shedding_activated
```

Relevant for greenhouses, off-grid cabins, underground shelters — any isolated habitat.
