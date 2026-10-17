# Dimos — Agentive OS for Physical Space

**Repo:** github.com/dimensionalOS/dimos  
**Language:** Python  
**License:** Pre-release beta  
**Role in OneMind:** Generalist robotics SDK — humanoids, quadrupeds, drones under agent control

Modern operating system for generalist robotics. No ROS required. Build physical applications in Python that run on any humanoid, quadruped, or drone. Agent-native: "vibecode" robots in natural language, build multi-agent systems that work seamlessly with hardware.

## Capabilities

- Navigation & Mapping — SLAM, dynamic obstacle avoidance, route planning, autonomous exploration
- Perception — detectors, 3D projections, VLMs, audio processing
- Agentic control — agents subscribe to embedded streams: perception (lidar, camera), spatial memory, control loops, motor drivers
- Agent CLI + MCP support
- Blueprints — reusable robot behavior definitions
- CUDA supported, NixOS supported, Docker-ready

## Why it matters for OneMind

OneMind's physical layer — when the fleet includes ground robots, drones, or humanoids, Dimos is the SDK that puts agents on the hardware. The agent model (subscribing to perception streams, publishing control) maps onto OneMind's NATS subject grammar:

```
tel.{node}.perception.camera   → agent subscription
tel.{node}.perception.lidar    → agent subscription
cmd.{node}.motion.{action}     → agent publishes
cmd.{node}.nav.waypoint        → agent publishes
```

## Integration path

Dimos agents run as modules on the robot, subscribing to NATS perception subjects published by IOnode/WireClaw sensors, and publishing control commands back to the NATS FABRIC account.

## See Also

- [[SentinelX]] — PX4-based drone interceptor (C-UAS specific)
- [[wireclaw]] — AI agent on ESP32 microcontroller
- [[ionode]] — sensor-only ESP32 NATS publisher
