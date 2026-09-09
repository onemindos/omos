# MAVLink → NATS Bridge (Phase 2)

## Overview

Bridges PX4 SITL drone telemetry into NATS JetStream so Legacy and other fabric agents
can see and command the drone in real time.

## Prerequisites

1. PX4 SITL running (via Docker):
```bash
docker run -d --name px4-drone \
  jonasvautherin/px4-gazebo-headless:1.14.0 \
  <YOUR_MAC_TAILSCALE_IP> <YOUR_MAC_TAILSCALE_IP>
```

2. Install mavsdk:
```bash
pip install mavsdk
```

3. NATS running with JetStream enabled.

## Run

```bash
MAVLINK_ADDRESS=udp://:14540 python run_mavlink.py
```

## Telemetry Flow

```
PX4 SITL → MAVLink UDP:14540 → mavsdk → NATS JetStream
  rad.greenwich.robots.drone1.telemetry
```

## Commanding the Drone via NATS

Publish an OMOS envelope to `rad.greenwich.robots.drone1.cmd`:
```json
{
  "id": "...",
  "ts": 1234567890.0,
  "source": "legacy",
  "event_type": "drone.cmd",
  "payload": { "action": "takeoff", "altitude_m": 10 }
}
```

## QGroundControl

QGC receives telemetry on UDP 14550 directly from PX4.
The NATS bridge runs alongside it — both can coexist.
