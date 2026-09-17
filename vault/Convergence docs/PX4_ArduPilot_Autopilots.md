# PX4 & ArduPilot: Open-Source Autopilot Systems for Autonomous Vehicles

**Links:**
- **PX4:** https://px4.io/ | https://github.com/px4/px4-autopilot/
- **ArduPilot:** https://ardupilot.org/ | https://github.com/ArduPilot/ardupilot

**Type:** Open-source flight control software / autopilot systems

**Status:** Production-ready, actively maintained, large community

**License:** BSD 3-Clause (PX4), GPLv3 (ArduPilot)

**Supported By:** 
- PX4: Dronecode Foundation (Linux Foundation non-profit)
- ArduPilot: Large open-source community

---

## What They Are

### PX4 Autopilot

**PX4 is an open-source flight control software** that provides autonomous control for drones and unmanned vehicles. It runs on flight control hardware (Pixhawk, etc.) and makes autonomous flight decisions based on sensor input, GPS, and operator commands.

**Core Purpose:**
- Autonomous flight stabilization
- Position control and waypoint navigation
- Mission planning and execution
- Multi-vehicle support
- Real-time sensor fusion
- Safety and failsafe systems

**Deployment Model:**
- Runs on Pixhawk and compatible flight computers
- Receives sensor data (IMU, GPS, barometer, magnetometer)
- Outputs motor commands via ESCs (Electronic Speed Controllers)
- Communicates via MAVLink protocol
- Integrates with ground station software (QGroundControl, etc.)

### ArduPilot Autopilot

**ArduPilot is a comprehensive open-source autopilot system** supporting diverse vehicle types: multirotors, helicopters, fixed-wing aircraft, boats, submarines, rovers, and more.

**Core Purpose:**
- Autonomous vehicle control across multiple platforms
- Trusted, battle-tested system (2009 to present)
- Flexible, modular architecture
- Community-driven development
- Production use worldwide

**Deployment Model:**
- Runs on wide variety of flight controller hardware
- Modular design supports custom platforms
- Extensive telemetry and logging
- Ground station integration (Mission Planner, etc.)
- Real-time operator communication

---

## What They Can Do: Core Capabilities

### PX4 Capabilities

**1. Autonomous Flight Control**
- Stabilized flight (roll, pitch, yaw)
- Altitude hold and position hold
- Attitude estimation (gyro + accelerometer fusion)
- Barometer-based altitude stabilization

**2. Navigation & Waypoints**
- GPS-based waypoint navigation
- Autonomous mission execution
- Return-to-home functionality
- Geofencing (virtual boundaries)

**3. Sensor Fusion**
- IMU (Inertial Measurement Unit) integration
- GPS positioning
- Magnetometer (compass) for heading
- Barometer for altitude
- Optional: LiDAR, optical flow, etc.

**4. Safety & Failsafe**
- Low battery detection
- GPS loss failsafe
- Communication loss recovery
- Propeller health monitoring
- Emergency landing procedures

**5. Modularity**
- Hardware abstraction layer (HAL)
- Pluggable modules (cameras, sensors, etc.)
- Custom vehicle support
- Simulation capability (SITL - Software In The Loop)

**Performance:**
- Stabilization loop rate: 400-1000 Hz (typically 1 kHz)
- GPS update rate: 10-50 Hz
- Latency: < 50ms typical
- Supports 100+ concurrent MAVLink connections

### ArduPilot Capabilities

**1. Multi-Vehicle Support**
- Multirotors (quadcopters, hexacopters, etc.)
- Fixed-wing aircraft
- Traditional helicopters
- Boats and submarines
- Ground rovers
- Balance-bots
- New vehicle types constantly emerging

**2. Advanced Flight Modes**
- Stabilize mode (manual with assistance)
- Altitude hold
- Loiter (hold position)
- Circle (orbit point)
- Follow (follow another vehicle)
- Autonomous missions
- Return-to-home
- Land (automated landing)

**3. Autonomous Mission Capability**
- Waypoint-based missions
- Conditional logic (IF/THEN branches)
- Parameter adjustments mid-mission
- Payload control (drop, camera, etc.)
- Multi-vehicle coordination

**4. Extensive Logging & Analysis**
- Comprehensive telemetry logging
- High-rate data recording (100+ Hz)
- Log playback and analysis
- Parameter logging for troubleshooting

**5. Hardware Support**
- Pixhawk family (most popular)
- Cube autopilots
- CUAV flight controllers
- Custom autopilot boards
- Zephyr OS HAL support (new)

**Performance:**
- Control loop: 400-1000 Hz
- Telemetry: Real-time MAVLink protocol
- Processing: Optimized for resource-constrained systems
- Scalability: Works on embedded Linux to ARM Cortex-M4

---

## Architecture: How They Work

### PX4 System Architecture

```
SENSORS
├─ IMU (accelerometer, gyroscope)
├─ Barometer (altitude)
├─ Magnetometer (heading)
├─ GPS (position)
└─ Optional (LiDAR, optical flow, etc.)
    ↓
FLIGHT CONTROL COMPUTER (Pixhawk)
├─ Sensor fusion (EKF - Extended Kalman Filter)
├─ Control loops (attitude, position, altitude)
├─ Navigation stack
├─ Safety monitoring
└─ Command arbitration
    ↓
ACTUATORS
├─ Motor commands (via ESCs)
└─ Servo outputs
    ↓
COMMUNICATION
├─ MAVLink protocol
├─ Ground station
├─ External systems (NATS, etc.)
└─ Real-time telemetry
```

### Flight Control Loops

```
Position Desired ──→ Position Controller ──→ Attitude Desired
                                                    ↓
                                            Attitude Controller
                                                    ↓
                                            Motor Commands
                                                    ↓
                                            Motors/ESCs
                                                    ↓
Sensors (IMU, GPS) ────────→ Sensor Fusion ─────── (feedback loop)
```

### ArduPilot System Architecture

```
VEHICLE INTERFACE
├─ Input from operator (RC, MAVLink)
├─ Mission/parameter input
└─ External commands
    ↓
FLIGHT LOGIC
├─ Flight mode selection
├─ Mission execution
├─ Autonomous navigation
├─ Vehicle-specific controllers
└─ Safety monitors
    ↓
SENSOR INTEGRATION
├─ IMU/gyro/accel
├─ Compass/magnetometer
├─ Barometer
├─ GPS
├─ Rangefinder (sonar, LiDAR)
└─ Optional payloads
    ↓
ACTUATOR CONTROL
├─ Motor/servo outputs
├─ Payload control
└─ ESC management
    ↓
TELEMETRY & COMMUNICATION
├─ MAVLink protocol
├─ Logging
├─ Ground station interface
└─ External system integration
```

### MAVLink Protocol

Both systems use **MAVLink** (Micro Air Vehicle Link) for communication:

```
Ground Station ←→ MAVLink Protocol ←→ Autopilot ←→ Sensors/Actuators
(QGroundControl,        (Telemetry link)     (ArduPilot/PX4)    (Motors, GPS,
Mission Planner)                                                   Compass, etc.)
```

**MAVLink provides:**
- Real-time vehicle state
- Mission upload/download
- Parameter adjustment
- Real-time telemetry (position, altitude, speed)
- Sensor data streaming
- Command injection (arm, disarm, mode change)

---

## Tactical Applications / Real-World Scenarios

### Scenario 1: FPV Attack Drone (Tenet Industries Model)

**Situation:**
```
Forward position needs real-time observation and strike capability.
Operator controls drone via FPV video feed.
Decision: Launch drone, observe area, attack target on command.
```

**How it applies:**
```
ArduPilot running on flight controller
  ↓ (Receives operator input via remote link)
  ↓ (Stabilizes aircraft while operator controls direction)
  ↓ (Maintains altitude and position)
  ↓ (Flies to target area)
  ↓ (Real-time telemetry back to operator)
  ↓ (Operator confirms target)
  ↓ (Executes attack)
Result: Precise FPV control with autopilot handling stabilization
```

**Specific use:**
- **Loiter mode:** Autopilot holds position while operator looks around
- **Stabilize mode:** Operator flies, autopilot prevents crashes
- **Waypoint mode:** Drone flies to location, operator takes FPV control

### Scenario 2: Autonomous ISR Mission

**Situation:**
```
Need surveillance of 10km² area.
No operator control needed.
Time-sensitive: must map before enemy moves.
```

**How it applies:**
```
PX4 running autonomous mission
  ↓
Operator uploads waypoint mission (grid pattern, camera triggers)
  ↓
Mission execution:
  - Takeoff to altitude
  - Fly waypoint 1 → trigger camera
  - Fly waypoint 2 → trigger camera
  - ... (repeat for all waypoints)
  - Return to home automatically
  ↓
Real-time telemetry to ground station
  ↓
Mission complete: full area photographed
Result: Complete coverage, no real-time operator input needed
```

**Specific use:**
- **Survey mode:** Automated grid-based scanning
- **Mission planner integration:** Pre-plan entire operation
- **Real-time monitoring:** See drone position and progress

### Scenario 3: Multi-Drone Coordination (GSoC 2026 Project)

**Situation:**
```
Multiple drones need coordinated swarm operation.
Limited communication bandwidth (mesh network).
Resilience needed: if one drone lost, others continue.
```

**How it applies:**
```
ArduPilot with AP_SwarmMesh (new project)
  ↓
Each drone running ArduPilot
  ↓
Ad-hoc mesh network between drones
  ↓
Resilient MAVLink networking
  ↓
Coordinated mission execution:
  - Drone 1: surveillance sweep
  - Drone 2: ISR observation
  - Drone 3: strike capability
  ↓
Loss of comms? Each drone continues autonomously
  ↓
Network recovers? All drone states sync
Result: Robust multi-vehicle mission in challenged comms environment
```

### Scenario 4: Boat/Submarine Operations

**Situation:**
```
Need autonomous underwater surveillance.
GPS unavailable underwater.
Compass/sonar-based navigation only.
```

**How it applies:**
```
ArduPilot Submarine mode
  ↓
Uses compass + sonar instead of GPS
  ↓
Depth hold via pressure sensor
  ↓
Autonomous patrol pattern
  ↓
Real-time telemetry via tether or mesh
  ↓
Pre-programmed mission or real-time pilot control
Result: Autonomous underwater ISR without GPS
```

---

## Integration with OneMind OS / Tactical C4ISR

### How They Fit in C4ISR

| C4ISR Component | PX4/ArduPilot Role |
|-----------------|-------------------|
| **Command** | Receives operator commands via MAVLink |
| **Control** | Executes coordinated actions (formation, swarm, etc.) |
| **Communications** | MAVLink telemetry to ground station/NATS |
| **Computers** | On-board autopilot (flight decision-making) |
| **Intelligence** | Sensor fusion, mission state awareness |
| **Surveillance** | Real-time position, payload data (camera, etc.) |
| **Reconnaissance** | Autonomous mission execution, area mapping |

### Integration Points with OneMind

**Scenario: Autonomous Drone Feed into OneMind TAK Server**

```
Drone (PX4/ArduPilot)
  ↓
  Position → MAVLink telemetry
  Sensor data → Real-time streaming
    ↓
  MAVLink-to-CoT Converter (PyTAK-based)
    ↓
  CoT Events (position, track, status)
    ↓
  NATS distribution
    ↓
  TAK Server ingestion
    ↓
  ATAK clients see:
  - Drone position (real-time)
  - Drone status (armed, flying, etc.)
  - Telemetry data (battery, signal, etc.)
  - Payload state (camera, armed, etc.)
    ↓
  Operator decision-making on unified COP
```

### Data Flow Model

```
PX4/ArduPilot Flight Computer
  ↓ (MAVLink protocol)
  ├─ Ground Station (QGroundControl, Mission Planner)
  ├─ Custom MAVLink bridge
  └─ MAVLink-to-CoT converter
      ↓ (CoT events)
      ├─ TAK Server
      ├─ NATS (for other systems)
      └─ ClickHouse (archival)
      ↓
      ATAK/TAK clients see drone on common operational picture
```

### Specific OneMind Integration

**With NATS:**
- MAVLink telemetry → NATS subject (`vehicles.drones.position`, etc.)
- Real-time position updates
- Telemetry available to all subscribers

**With TAK Server:**
- Drone as CoT entity on map
- Real-time position + track
- Payload status (camera, armed, etc.)
- Formation information

**With PostGIS:**
- Drone path history
- Operational area boundaries
- Waypoint storage (missions)
- Geofencing zones

**With DDF:**
- Drone capabilities catalog
- Payload information
- Hardware specifications
- Operational parameters

---

## Specifications / Technical Details

### PX4 Specifications

| Specification | Value |
|---------------|-------|
| **Flight Control Loop** | 1000 Hz (1ms cycle) |
| **Sensor Fusion** | Extended Kalman Filter (EKF3) |
| **Position Update Rate** | 10-50 Hz (GPS dependent) |
| **Typical Latency** | <50ms sensor to actuator |
| **Supported Vehicles** | Multirotor, fixed-wing, VTOL, ground rovers |
| **Flight Controller** | Pixhawk (ARM Cortex-M7), also ARM Cortex-M4 |
| **Memory** | Typical 256KB-512KB RAM, 2MB+ Flash |
| **Max Flight Time** | Vehicle dependent (30min+ typical) |
| **Communication Range** | 5-20 km typical (depends on radio link) |
| **Code Base** | C++ (150,000+ lines) |
| **Simulation** | SITL (Software In The Loop) with Gazebo |

### ArduPilot Specifications

| Specification | Value |
|---------------|-------|
| **Flight Control Loop** | 400-1000 Hz (vehicle dependent) |
| **Supported Vehicles** | 8+ major types (multirotor, fixed-wing, heli, boat, rover, submarine, etc.) |
| **Flight Controller** | Pixhawk, Cube, CUAV, and others |
| **Memory** | 256KB-2MB depending on platform |
| **Telemetry Rate** | 10-50 Hz typical |
| **Max Waypoints** | 3000+ (limited by memory) |
| **Code Base** | C++ (200,000+ lines) |
| **Logging Rate** | 10-100 Hz (configurable) |
| **Communication Range** | 5-40 km depending on radio |
| **Simulation** | SITL with Gazebo/Jsbsim |
| **Battery Life** | Vehicle dependent (15min-120min typical) |

### Real-World Performance (Ukraine Operations)

**FPV Drones on ArduPilot/PX4:**
- Flight time: 10-25 minutes
- Speed: 50-150 km/h
- Range (radio): 3-10 km (with extended antennas)
- Accuracy: ±2-5 meters (without RTK)
- Consumption rate: 1,000+ drones/week in high-intensity ops

---

## Comparison: PX4 vs. ArduPilot

| Aspect | PX4 | ArduPilot |
|--------|-----|-----------|
| **Vehicle Support** | Multirotor, fixed-wing, VTOL, rovers | All PX4 types + helicopters, boats, submarines, etc. |
| **Community Size** | Medium (academic + professional) | Large (30,000+ users worldwide) |
| **Commercial Use** | Yes (DJI heritage) | Yes (extensive production use) |
| **Learning Curve** | Moderate | Moderate to steep |
| **Documentation** | Good (developer-focused) | Excellent (extensive user guides) |
| **Customization** | Excellent (modular design) | Excellent (mature codebase) |
| **Hardware Support** | Pixhawk + some others | Very wide (100+ boards) |
| **Mission Planning** | QGroundControl | Mission Planner, QGC, etc. |
| **Active Development** | Rapid (bleeding edge) | Stable + incremental (tested features) |
| **Simulation** | SITL excellent | SITL excellent |
| **Real-Time OS** | NuttX (real-time) | Linux/Zephyr (real-time kernel options) |
| **Best For** | Research, cutting-edge, robotics | Production, reliability, diverse platforms |

### Why Choose PX4

- **Cutting-edge features** — Rapid development, new capabilities first
- **Research use** — Academic and experimental platforms
- **Performance tuning** — Fine-grained control over loops
- **Standardization** — Dronecode ecosystem support
- **Modularity** — Easy to add custom functionality

### Why Choose ArduPilot

- **Battle-tested** — Proven in production since 2009
- **Diverse platforms** — Helicopters, boats, submarines, rovers
- **Extensive documentation** — Beginner to advanced
- **Large community** — Many examples, tutorials, help
- **Reliability** — Conservative, tested releases
- **Commercial support** — Commercial options available

---

## Related Resources in Convergence

- **Tenet_Industries_Attack_Drones.md** — FPV drones using autopilots
- **Protocol_to_CoT_Converters.md** — MAVLink-to-CoT conversion (PyTAK)
- **NSTAC_Family.md** — PyTAK library for drone-to-TAK integration
- **Evolution_to_C4ISR.md** — Drone ISR role in C4ISR architecture
- **TAKpilot.md** — AI interface for drone command & control
- **ATAK_GIS_Data_Display.md** — Displaying drone tracks on TAK

---

**Captured:** 2026-07-11

**Status:** Production-ready, actively maintained

**Relevance:** Core autopilot systems for autonomous drone operations, essential for FPV attack drones, ISR missions, and multi-vehicle coordination

**Learning Curve:** PX4 (moderate), ArduPilot (moderate to steep)

**Key Insight:** Open-source autopilots (PX4/ArduPilot) enable rapid drone development at scale. ArduPilot's diversity supports any vehicle; PX4's modularity supports cutting-edge research. Together they power most autonomous drones worldwide, including modern FPV systems.

**Strategic Importance:** MAVLink protocol enables standardized drone integration with TAK Server and NATS. Autopilots handle vehicle control; OneMind handles coordination. This separation enables modular, resilient tactical systems.

**Integration Pattern:** Drone (autopilot) → MAVLink → CoT converter → NATS/TAK → Operators see unified COP with all vehicles (drones, aircraft, ground units, maritime) in real-time.
