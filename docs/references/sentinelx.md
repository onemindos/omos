# SentinelX — Autonomous C-UAS Interceptor

**Language:** C++ (ROS2) + Python (launch)  
**Stack:** PX4 SIH v1.15.4 → uXRCE-DDS → ROS2 Humble → SentinelX  
**Role in OneMind:** Drone interceptor autonomy layer for C-UAS operations

ROS2 node stack implementing an autonomous counter-UAS interceptor. Receives target tracks from a C2 system, executes intercept missions via PX4 offboard control, and reports engagement results.

## Nodes

| Node | Role |
|------|------|
| `launch_node` / `BaseNode` | Master interceptor lifecycle — ASSIGN_TARGET, START_INTERCEPT, ABORT state machine |
| `c2_guidance_node` | Mid-course guidance — navigates to target intercept point |
| `terminal_guidance_node` | Terminal phase — seeker lock and final intercept |
| `px4_proxy` | PX4 ↔ ROS2 bridge — publishes offboard setpoints via MAVLink/uXRCE-DDS |

## ROS2 Topics

```
IN:  /cuas/c2/command          → C2Command (ASSIGN_TARGET, START_INTERCEPT, ABORT)
IN:  /cuas/c2/target_track     → TargetTrack (position + velocity of threat)

OUT: /cuas/interceptor/ack     → MissionAck
OUT: /cuas/interceptor/result  → EngagementResult
OUT: /cuas/interceptor/fault   → FaultReport
OUT: /cuas/interceptor/snapshot → InterceptorSnapshot (100ms telemetry)
```

## Integration Path to OneMind

OneMind bridges ROS2 topics ↔ NATS subjects:

```
cmd.defense.interceptor.{id}.assign    → /cuas/c2/command (ASSIGN_TARGET)
cmd.defense.interceptor.{id}.intercept → /cuas/c2/command (START_INTERCEPT)
cmd.defense.interceptor.{id}.abort     → /cuas/c2/command (ABORT)

evt.defense.interceptor.{id}.ack       ← /cuas/interceptor/ack
evt.defense.interceptor.{id}.result    ← /cuas/interceptor/result
evt.defense.interceptor.{id}.snapshot  ← /cuas/interceptor/snapshot
tel.{node}.interceptor.position        ← /cuas/interceptor/snapshot (100ms)
```

## Full Kill Chain

```
APEX-SENTINEL (detect + track)
        ↓
  evt.defense.drone.tracked (NATS)
        ↓
  SentinelX target_track bridge
        ↓
  /cuas/c2/target_track (ROS2)
        ↓
  SentinelX autonomy → PX4 SIH → physical intercept
        ↓
  evt.defense.engagement.result (NATS)
        ↓
  cot2constellation → Constellation Overwatch COP
```

## Setup (verified environment)

| Component | Version |
|-----------|---------|
| Ubuntu | 22.04 |
| PX4 | v1.15.4 |
| ROS2 | Humble |
| Simulator | SIH QuadX |
| DDS | Micro XRCE DDS |

```bash
# PX4 SIH
cd ~/PX4-Autopilot/build/px4_sitl_default
PX4_SIM_MODEL=sihsim_quadx ./bin/px4 -s etc/init.d-posix/rcS etc

# DDS bridge
MicroXRCEAgent udp4 -p 8888

# SentinelX
ros2 launch sentinelx sentinelx_launch.py
```

## See Also

- [[apex-sentinel]] — upstream threat detection feeding target tracks
- [[CONSTELLATION-OVERWATCH]] — COP that receives engagement results via cot2constellation
