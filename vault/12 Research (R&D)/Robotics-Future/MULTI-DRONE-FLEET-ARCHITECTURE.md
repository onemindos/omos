# Multi-Drone Fleet Architecture (OpenMind Structure)

> **Created:** 2026-01-22
> **Architecture:** OneMind OS + OpenMind OM1 + FABRIC Network
> **Fleet:** Custom DIY (Zeus/Personal) + DJI (Farming) + Skydio (Security)

---

## Executive Summary

**Goal:** Multi-drone fleet coordination using OpenMind.org's FABRIC network architecture, integrated with OneMind OS (Legacy AI).

**Architecture:** One edge device controls multiple drone types through unified API, similar to how OpenMind OM1 controls heterogeneous robot fleets.

**Fleet Composition:**
1. **Zeus Drone** (DIY Pixhawk) - Personal tasks, custom payloads
2. **Farm Drone** (DJI Agras T40) - Agricultural operations, crop spraying
3. **Security Drone** (Skydio 2+) - Autonomous patrol, surveillance

**Cost:** ~$18K hardware + $997 edge compute + $135/mo services

---

## OpenMind Architecture Applied to Drones

### OpenMind's Robot Architecture

```
┌─────────────────────────────────────────────────────┐
│          OpenMind OM1 (Robot Operating System)      │
│                                                      │
│  Hardware-agnostic API for ANY robot type           │
│  - Unitree quadrupeds                               │
│  - Boston Dynamics Spot                             │
│  - Humanoids (Figure 01, Tesla Optimus)             │
│  - Custom robots                                    │
│                                                      │
│  Uses FABRIC Network for coordination               │
└─────────────────────────────────────────────────────┘
```

### Our Drone Fleet Architecture (Same Pattern)

```
┌─────────────────────────────────────────────────────────────────┐
│         OneMind Drone Manager (Drone Operating System)          │
│                                                                  │
│  Hardware-agnostic API for ANY drone type                       │
│  - Custom DIY (MAVLink)                                         │
│  - DJI Agricultural (Mobile SDK)                                │
│  - Skydio Security (Cloud API)                                  │
│  - Future: any other drone                                      │
│                                                                  │
│  Uses NATS + FABRIC pattern for coordination                    │
└─────────────────────────────────────────────────────────────────┘
```

---

## System Architecture Diagram

```
┌───────────────────────────────────────────────────────────────────────┐
│                       CLOUD (100.102.21.44)                            │
│                                                                       │
│  ┌────────────────────────────────────────────────────────────┐     │
│  │  Legacy Agent (Claude Opus 4.5)                            │     │
│  │  - Voice: "Launch Zeus for perimeter check"                │     │
│  │  - Voice: "Start farm spraying operation"                  │     │
│  │  - Voice: "Skydio, patrol the property"                    │     │
│  └─────────────────────┬──────────────────────────────────────┘     │
│                        │                                             │
│  ┌─────────────────────▼──────────────────────────────────────┐     │
│  │  NATS Hub (JetStream)                                      │     │
│  │  - commands.farm.drone.>                                   │     │
│  │  - events.farm.drone.>                                     │     │
│  │  - telemetry.farm.drone.>                                  │     │
│  └─────────────────────┬──────────────────────────────────────┘     │
│                        │                                             │
│  ┌─────────────────────▼──────────────────────────────────────┐     │
│  │  LiveKit Server (Video Streaming)                          │     │
│  │  - Room: drone-zeus-feed                                   │     │
│  │  - Room: drone-farm-feed                                   │     │
│  │  - Room: drone-security-feed                               │     │
│  └────────────────────────────────────────────────────────────┘     │
│                                                                       │
└────────────────────────────┬──────────────────────────────────────┬─┘
                             │ (Tailscale VPN)                      │
                             │                                      │
┌────────────────────────────▼──────────────────────────────────────▼─┐
│                   EDGE DEVICE (Farm Location)                        │
│                   OpenMind Edge Compute or Jetson Orin              │
│                                                                      │
│  ┌───────────────────────────────────────────────────────────┐     │
│  │  OneMind Drone Manager (Drone OS)                         │     │
│  │  - Hardware abstraction layer                             │     │
│  │  - Unified drone API                                      │     │
│  │  - Mission planning & coordination                        │     │
│  │  - Fleet state management                                 │     │
│  └───────────────────────────────────────────────────────────┘     │
│                              │                                      │
│  ┌───────────────────────────┼──────────────────────────────┐      │
│  │        NATS Leaf Node (Local Message Bus)                │      │
│  │        - Syncs with cloud NATS Hub                        │      │
│  │        - Offline-capable with replay                      │      │
│  └───────────────────────────┼──────────────────────────────┘      │
│                              │                                      │
│       ┌──────────────────────┼──────────────────────┐              │
│       │                      │                      │              │
│       ▼                      ▼                      ▼              │
│  ┌─────────┐           ┌─────────┐           ┌─────────┐          │
│  │ Zeus    │           │  Farm   │           │Security │          │
│  │ Bridge  │           │ Bridge  │           │ Bridge  │          │
│  │(MAVLink)│           │  (DJI)  │           │(Skydio) │          │
│  └────┬────┘           └────┬────┘           └────┬────┘          │
│       │                     │                     │                │
│       │ Serial/UDP          │ Mobile SDK          │ Cloud API      │
│       │                     │                     │                │
│       ▼                     ▼                     ▼                │
│  ┌─────────┐           ┌─────────┐           ┌─────────┐          │
│  │  ZEUS   │           │   DJI   │           │ SKYDIO  │          │
│  │  DRONE  │           │  AGRAS  │           │   2+    │          │
│  │ (Custom)│           │   T40   │           │(Security)│          │
│  └─────────┘           └─────────┘           └─────────┘          │
│   Pixhawk              Agricultural          Autonomous           │
│   DIY Build            Spraying Drone        Patrol Drone         │
│                                                                    │
│  All drones share:                                                │
│  - Same edge compute device                                       │
│  - Same NATS network                                              │
│  - Coordinated via Drone Manager                                  │
│  - Video streams via LiveKit                                      │
│  - Controlled by Legacy AI                                        │
│                                                                    │
└────────────────────────────────────────────────────────────────────┘
```

---

## Fleet Composition

### Drone 1: Zeus (Personal Tasks)

**Hardware:** Custom DIY Build with Pixhawk 6C
**Protocol:** MAVLink (native)
**Primary Use:** Personal errands, custom payloads, experimentation

**Specifications:**
| Component | Choice | Price |
|-----------|--------|-------|
| Flight Controller | Pixhawk 6C | $200 |
| Frame | S500 Quadcopter | $80 |
| Motors | 920KV (4x) | $80 |
| ESCs | 30A (4x) | $60 |
| Battery | 4S 5000mAh | $80 |
| GPS | UBlox M8N | $40 |
| Camera | Raspberry Pi HQ | $50 |
| Companion | Raspberry Pi 4 | $75 |
| Telemetry | 433MHz Radio | $60 |
| FPV Camera | RunCam | $40 |
| **Total** | | **$765** |

**Custom Payload Capabilities:**
- Package delivery (up to 1kg)
- Custom sensors (air quality, radiation, etc.)
- Experimental AI models (runs on Pi 4)
- RF mapping equipment
- Audio recording/broadcasting
- Emergency medical delivery

**Voice Commands:**
```
"Zeus, deliver this package to the workshop"
"Zeus, check air quality near the creek"
"Zeus, follow me for video"
"Zeus, map the property perimeter"
```

### Drone 2: Farm Operations (DJI Agras T40)

**Hardware:** DJI Agras T40 Agricultural Drone
**Protocol:** DJI Mobile SDK → Bridge
**Primary Use:** Crop spraying, agricultural operations

**Specifications:**
| Feature | Spec |
|---------|------|
| Price | $16,000 |
| Spray Capacity | 40kg (10.5 gallons) |
| Spray Width | 9-11 meters |
| Flight Time | 10-15 min (with full load) |
| Coverage | 15-20 acres/hour |
| Terrain Following | Yes (radar) |
| Obstacle Avoidance | 360° spherical radar |
| Refill Time | 30 seconds |
| Battery Swap | 10 seconds |

**Agricultural Workflows:**
```python
@workflow
def morning_crop_spray():
    """Automated crop spraying"""
    # Check weather conditions
    weather = await sensors.get_weather()
    if weather.wind_speed > 10:  # mph
        return "Wind too high for spraying"

    # Pre-flight checks
    farm_drone.preflight_check()
    farm_drone.fill_tank(chemical="pesticide", volume=40)

    # Execute spraying mission
    farm_drone.load_mission("greenhouse_rows.kml")
    farm_drone.takeoff()
    farm_drone.execute_mission()

    # Monitor progress
    while farm_drone.mission_active():
        telemetry = farm_drone.get_telemetry()
        await legacy.publish_status(telemetry)
        await asyncio.sleep(5)

    farm_drone.rtl()
    legacy.speak("Crop spraying complete. Covered 15 acres.")
```

**Voice Commands:**
```
"Start spraying the east field"
"Check spray levels"
"Return farm drone for refill"
"Show me crop coverage map"
```

### Drone 3: Security Patrol (Skydio 2+)

**Hardware:** Skydio 2+ Autonomous Drone
**Protocol:** Skydio Cloud API → Bridge
**Primary Use:** Security patrol, autonomous surveillance

**Specifications:**
| Feature | Spec |
|---------|------|
| Price | $1,099 |
| Flight Time | 27 minutes |
| Obstacle Avoidance | 360° (6 cameras) |
| AI Tracking | Best-in-class |
| Camera | 4K60 HDR |
| Range | 6 km |
| Autonomy | Fully autonomous |
| Night Vision | No (infrared optional) |

**Security Workflows:**
```python
@workflow
def perimeter_patrol():
    """Automated security patrol"""
    security_drone.takeoff()

    # Define patrol waypoints
    waypoints = [
        {"name": "barn_north", "lat": 40.123, "lon": -74.456, "alt": 15},
        {"name": "greenhouse", "lat": 40.124, "lon": -74.457, "alt": 15},
        {"name": "equipment_yard", "lat": 40.125, "lon": -74.458, "alt": 15},
        {"name": "property_gate", "lat": 40.126, "lon": -74.459, "alt": 20},
    ]

    for waypoint in waypoints:
        security_drone.goto(waypoint)

        # Hover and scan
        video = security_drone.capture_video(duration=10)

        # AI analysis
        analysis = await legacy.analyze_video(video,
            objectives=["detect_people", "detect_vehicles", "detect_animals"]
        )

        if analysis.threats_detected:
            await legacy.send_alert(f"Security alert at {waypoint['name']}")

    security_drone.rtl()
```

**Voice Commands:**
```
"Skydio, patrol the perimeter"
"Check the north fence line"
"Follow that vehicle"
"Return to base"
"Show me the security feed"
```

---

## OneMind Drone Manager (Core Software)

**Location:** `agents/edge/drone_manager/`

**Purpose:** Hardware-agnostic drone control layer (like OpenMind OM1 for robots)

### Architecture

```python
# agents/edge/drone_manager/core.py

from typing import Dict, Any, List
from enum import Enum
import asyncio
import nats
from dataclasses import dataclass

class DroneType(Enum):
    CUSTOM_MAVLINK = "custom_mavlink"
    DJI_AGRICULTURAL = "dji_agricultural"
    SKYDIO_SECURITY = "skydio_security"

class DroneStatus(Enum):
    IDLE = "idle"
    PREFLIGHT = "preflight"
    TAKING_OFF = "taking_off"
    FLYING = "flying"
    MISSION = "mission"
    RETURNING = "returning"
    LANDING = "landing"
    ERROR = "error"
    OFFLINE = "offline"

@dataclass
class DroneState:
    """Unified drone state (hardware-agnostic)"""
    drone_id: str
    drone_type: DroneType
    status: DroneStatus

    # Position
    latitude: float
    longitude: float
    altitude: float
    heading: float

    # Telemetry
    battery_percent: int
    signal_strength: int
    flight_time: int  # seconds

    # Mission
    mission_active: bool
    mission_progress: float  # 0-100

    # Video
    video_streaming: bool
    livekit_room: str

class DroneManager:
    """
    Hardware-agnostic drone fleet manager

    Similar to OpenMind OM1, provides unified API for all drone types.
    Each drone has a protocol bridge that translates between native
    protocol (MAVLink, DJI SDK, Skydio API) and this unified interface.
    """

    def __init__(self, nats_url: str):
        self.nats_url = nats_url
        self.nc = None
        self.drones: Dict[str, DroneState] = {}
        self.bridges: Dict[str, 'DroneBridge'] = {}

    async def connect(self):
        """Connect to NATS"""
        self.nc = await nats.connect(self.nats_url)

        # Subscribe to command channel
        await self.nc.subscribe(
            "commands.farm.drone.>",
            cb=self.handle_command
        )

        # Subscribe to telemetry from bridges
        await self.nc.subscribe(
            "telemetry.farm.drone.>",
            cb=self.handle_telemetry
        )

    async def register_drone(
        self,
        drone_id: str,
        drone_type: DroneType,
        bridge: 'DroneBridge'
    ):
        """Register a drone with the fleet"""
        self.drones[drone_id] = DroneState(
            drone_id=drone_id,
            drone_type=drone_type,
            status=DroneStatus.OFFLINE,
            latitude=0.0,
            longitude=0.0,
            altitude=0.0,
            heading=0.0,
            battery_percent=0,
            signal_strength=0,
            flight_time=0,
            mission_active=False,
            mission_progress=0.0,
            video_streaming=False,
            livekit_room=f"drone-{drone_id}-feed"
        )

        self.bridges[drone_id] = bridge
        await bridge.connect(self.nc, drone_id)

        # Publish registration event
        await self.nc.publish(
            f"events.farm.drone.{drone_id}.registered",
            json.dumps({
                "drone_id": drone_id,
                "drone_type": drone_type.value,
                "timestamp": time.time()
            }).encode()
        )

    async def handle_command(self, msg):
        """
        Handle unified drone commands (hardware-agnostic)

        Examples:
          commands.farm.drone.zeus.takeoff
          commands.farm.drone.farm.goto
          commands.farm.drone.security.patrol
        """
        subject_parts = msg.subject.split(".")
        # commands.farm.drone.{drone_id}.{action}
        drone_id = subject_parts[3]
        action = subject_parts[4]

        data = json.loads(msg.data.decode()) if msg.data else {}

        if drone_id not in self.bridges:
            await self.nc.publish(
                msg.reply,
                json.dumps({"error": f"Unknown drone: {drone_id}"}).encode()
            )
            return

        bridge = self.bridges[drone_id]

        # Translate to unified command
        try:
            if action == "takeoff":
                result = await bridge.takeoff(data.get("altitude", 10))

            elif action == "land":
                result = await bridge.land()

            elif action == "goto":
                result = await bridge.goto(
                    lat=data["latitude"],
                    lon=data["longitude"],
                    alt=data.get("altitude", 15)
                )

            elif action == "rtl":  # Return to launch
                result = await bridge.rtl()

            elif action == "orbit":
                result = await bridge.orbit(
                    radius=data.get("radius", 10),
                    altitude=data.get("altitude", 15)
                )

            elif action == "follow":
                result = await bridge.follow(
                    target=data.get("target", "person")
                )

            elif action == "mission":
                result = await bridge.execute_mission(
                    waypoints=data["waypoints"]
                )

            elif action == "stream_video":
                result = await bridge.start_video_stream()

            else:
                result = {"error": f"Unknown action: {action}"}

            # Send response
            if msg.reply:
                await self.nc.publish(
                    msg.reply,
                    json.dumps(result).encode()
                )

        except Exception as e:
            await self.nc.publish(
                msg.reply,
                json.dumps({"error": str(e)}).encode()
            )

    async def handle_telemetry(self, msg):
        """Update drone state from telemetry"""
        subject_parts = msg.subject.split(".")
        drone_id = subject_parts[3]

        if drone_id not in self.drones:
            return

        telemetry = json.loads(msg.data.decode())
        state = self.drones[drone_id]

        # Update state
        state.latitude = telemetry.get("latitude", state.latitude)
        state.longitude = telemetry.get("longitude", state.longitude)
        state.altitude = telemetry.get("altitude", state.altitude)
        state.heading = telemetry.get("heading", state.heading)
        state.battery_percent = telemetry.get("battery", state.battery_percent)
        state.signal_strength = telemetry.get("signal", state.signal_strength)
        state.flight_time = telemetry.get("flight_time", state.flight_time)

        # Publish unified state
        await self.nc.publish(
            f"events.farm.drone.{drone_id}.state",
            json.dumps(self._serialize_state(state)).encode()
        )

    def _serialize_state(self, state: DroneState) -> dict:
        """Serialize drone state to dict"""
        return {
            "drone_id": state.drone_id,
            "drone_type": state.drone_type.value,
            "status": state.status.value,
            "position": {
                "latitude": state.latitude,
                "longitude": state.longitude,
                "altitude": state.altitude,
                "heading": state.heading
            },
            "telemetry": {
                "battery": state.battery_percent,
                "signal": state.signal_strength,
                "flight_time": state.flight_time
            },
            "mission": {
                "active": state.mission_active,
                "progress": state.mission_progress
            },
            "video": {
                "streaming": state.video_streaming,
                "room": state.livekit_room
            }
        }

    async def get_fleet_status(self) -> List[dict]:
        """Get status of all drones"""
        return [
            self._serialize_state(state)
            for state in self.drones.values()
        ]

    async def coordinate_mission(
        self,
        mission_type: str,
        drones: List[str],
        params: dict
    ):
        """
        Coordinate multi-drone mission

        Examples:
          - Coordinated patrol (multiple security drones)
          - Search pattern (all drones grid search)
          - Formation flight
        """
        # This is where FABRIC-like coordination would happen
        # For now, simple sequential execution

        for drone_id in drones:
            if drone_id not in self.bridges:
                continue

            bridge = self.bridges[drone_id]

            # Assign mission waypoints based on mission type
            if mission_type == "patrol":
                # Divide patrol route among drones
                pass

            elif mission_type == "search":
                # Assign grid cells to each drone
                pass

# Base class for drone-specific bridges
class DroneBridge:
    """
    Abstract base for drone protocol bridges

    Each drone type implements this interface:
    - MAVLinkBridge (Zeus)
    - DJIBridge (Farm)
    - SkydioBridge (Security)
    """

    async def connect(self, nc, drone_id: str):
        """Connect to drone and NATS"""
        raise NotImplementedError

    async def takeoff(self, altitude: float) -> dict:
        raise NotImplementedError

    async def land(self) -> dict:
        raise NotImplementedError

    async def goto(self, lat: float, lon: float, alt: float) -> dict:
        raise NotImplementedError

    async def rtl(self) -> dict:
        """Return to launch"""
        raise NotImplementedError

    async def orbit(self, radius: float, altitude: float) -> dict:
        raise NotImplementedError

    async def follow(self, target: str) -> dict:
        raise NotImplementedError

    async def execute_mission(self, waypoints: List[dict]) -> dict:
        raise NotImplementedError

    async def start_video_stream(self) -> dict:
        raise NotImplementedError

    async def publish_telemetry(self, telemetry: dict):
        """Publish telemetry to NATS"""
        await self.nc.publish(
            f"telemetry.farm.drone.{self.drone_id}",
            json.dumps(telemetry).encode()
        )
```

---

## Protocol Bridges (One Per Drone Type)

### Bridge 1: MAVLink (Zeus Custom Drone)

```python
# agents/edge/drone_manager/bridges/mavlink_bridge.py

from pymavlink import mavutil
from drone_manager.core import DroneBridge
import asyncio

class MAVLinkBridge(DroneBridge):
    """Bridge for MAVLink-based drones (ArduPilot/PX4)"""

    def __init__(self, connection_string: str):
        """
        Args:
            connection_string: Serial or UDP connection
              - Serial: "serial:/dev/ttyUSB0:57600"
              - UDP: "udp:127.0.0.1:14550"
        """
        self.connection_string = connection_string
        self.drone = None
        self.nc = None
        self.drone_id = None

    async def connect(self, nc, drone_id: str):
        """Connect to MAVLink drone and NATS"""
        self.nc = nc
        self.drone_id = drone_id

        # Connect to drone
        self.drone = mavutil.mavlink_connection(self.connection_string)

        # Wait for heartbeat
        print(f"Waiting for heartbeat from {drone_id}...")
        self.drone.wait_heartbeat()
        print(f"Heartbeat from {drone_id} received")

        # Start telemetry loop
        asyncio.create_task(self.telemetry_loop())

    async def takeoff(self, altitude: float) -> dict:
        """Takeoff to specified altitude"""
        self.drone.mav.command_long_send(
            self.drone.target_system,
            self.drone.target_component,
            mavutil.mavlink.MAV_CMD_NAV_TAKEOFF,
            0,  # confirmation
            0, 0, 0, 0,  # params 1-4
            0, 0,  # lat, lon (use current)
            altitude  # altitude
        )

        return {"status": "taking_off", "altitude": altitude}

    async def land(self) -> dict:
        """Land at current position"""
        self.drone.mav.command_long_send(
            self.drone.target_system,
            self.drone.target_component,
            mavutil.mavlink.MAV_CMD_NAV_LAND,
            0, 0, 0, 0, 0, 0, 0, 0
        )

        return {"status": "landing"}

    async def goto(self, lat: float, lon: float, alt: float) -> dict:
        """Fly to GPS coordinates"""
        self.drone.mav.mission_item_send(
            self.drone.target_system,
            self.drone.target_component,
            0,  # seq
            mavutil.mavlink.MAV_FRAME_GLOBAL_RELATIVE_ALT,
            mavutil.mavlink.MAV_CMD_NAV_WAYPOINT,
            2,  # current (2=guided mode)
            1,  # autocontinue
            0, 0, 0, 0,  # params
            lat, lon, alt
        )

        return {"status": "navigating", "destination": {"lat": lat, "lon": lon, "alt": alt}}

    async def rtl(self) -> dict:
        """Return to launch"""
        self.drone.mav.command_long_send(
            self.drone.target_system,
            self.drone.target_component,
            mavutil.mavlink.MAV_CMD_NAV_RETURN_TO_LAUNCH,
            0, 0, 0, 0, 0, 0, 0, 0
        )

        return {"status": "returning"}

    async def orbit(self, radius: float, altitude: float) -> dict:
        """Orbit current position"""
        self.drone.mav.command_long_send(
            self.drone.target_system,
            self.drone.target_component,
            mavutil.mavlink.MAV_CMD_DO_ORBIT,
            0,
            radius,  # radius in meters
            1.0,  # velocity (m/s)
            0,  # yaw behavior (0=hold)
            0,  # reserved
            0, 0,  # lat/lon (use current)
            altitude
        )

        return {"status": "orbiting", "radius": radius, "altitude": altitude}

    async def follow(self, target: str) -> dict:
        """Follow mode not natively supported by MAVLink, would need custom implementation"""
        return {"error": "Follow mode requires custom implementation"}

    async def execute_mission(self, waypoints: List[dict]) -> dict:
        """Execute waypoint mission"""
        # Upload mission
        self.drone.waypoint_clear_all_send()

        for i, wp in enumerate(waypoints):
            self.drone.mav.mission_item_send(
                self.drone.target_system,
                self.drone.target_component,
                i,
                mavutil.mavlink.MAV_FRAME_GLOBAL_RELATIVE_ALT,
                mavutil.mavlink.MAV_CMD_NAV_WAYPOINT,
                0,  # current
                1,  # autocontinue
                0, 0, 0, 0,
                wp["latitude"],
                wp["longitude"],
                wp.get("altitude", 15)
            )

        # Start mission
        self.drone.mav.command_long_send(
            self.drone.target_system,
            self.drone.target_component,
            mavutil.mavlink.MAV_CMD_MISSION_START,
            0, 0, 0, 0, 0, 0, 0, 0
        )

        return {"status": "mission_started", "waypoints": len(waypoints)}

    async def start_video_stream(self) -> dict:
        """Start video streaming to LiveKit"""
        # This would be handled by companion computer (Raspberry Pi)
        # sending video to LiveKit
        return {"status": "streaming", "room": f"drone-{self.drone_id}-feed"}

    async def telemetry_loop(self):
        """Continuously publish telemetry to NATS"""
        while True:
            try:
                # Get GPS position
                msg = self.drone.recv_match(type='GLOBAL_POSITION_INT', blocking=False)
                if msg:
                    telemetry = {
                        "latitude": msg.lat / 1e7,
                        "longitude": msg.lon / 1e7,
                        "altitude": msg.relative_alt / 1000.0,
                        "heading": msg.hdg / 100.0
                    }

                # Get battery
                msg = self.drone.recv_match(type='BATTERY_STATUS', blocking=False)
                if msg:
                    telemetry["battery"] = msg.battery_remaining

                # Get signal strength
                msg = self.drone.recv_match(type='RADIO_STATUS', blocking=False)
                if msg:
                    telemetry["signal"] = msg.rssi

                # Publish to NATS
                await self.publish_telemetry(telemetry)

                await asyncio.sleep(1)  # 1Hz telemetry

            except Exception as e:
                print(f"Telemetry error: {e}")
                await asyncio.sleep(5)
```

### Bridge 2: DJI SDK (Farm Drone)

```python
# agents/edge/drone_manager/bridges/dji_bridge.py

from drone_manager.core import DroneBridge
# Note: DJI requires Mobile SDK running on Android/iOS device
# This bridge would communicate with a companion app

class DJIBridge(DroneBridge):
    """
    Bridge for DJI drones (Agras T40, Mavic, etc.)

    Requires DJI Mobile SDK running on tablet/phone
    Communicates via local HTTP API or WebSocket
    """

    def __init__(self, mobile_sdk_url: str):
        """
        Args:
            mobile_sdk_url: URL of DJI Mobile SDK companion app
              Example: "http://192.168.1.100:8080"
        """
        self.mobile_sdk_url = mobile_sdk_url
        self.nc = None
        self.drone_id = None

    async def connect(self, nc, drone_id: str):
        """Connect to DJI Mobile SDK and NATS"""
        self.nc = nc
        self.drone_id = drone_id

        # Connect to DJI Mobile SDK companion app
        # This app runs on Android/iOS device connected to DJI controller

        # Start telemetry loop
        asyncio.create_task(self.telemetry_loop())

    async def _send_command(self, command: str, params: dict = None) -> dict:
        """Send command to DJI Mobile SDK companion app"""
        import aiohttp

        async with aiohttp.ClientSession() as session:
            async with session.post(
                f"{self.mobile_sdk_url}/api/command",
                json={"command": command, "params": params or {}}
            ) as resp:
                return await resp.json()

    async def takeoff(self, altitude: float) -> dict:
        return await self._send_command("takeoff", {"altitude": altitude})

    async def land(self) -> dict:
        return await self._send_command("land")

    async def goto(self, lat: float, lon: float, alt: float) -> dict:
        return await self._send_command("goto", {
            "latitude": lat,
            "longitude": lon,
            "altitude": alt
        })

    async def rtl(self) -> dict:
        return await self._send_command("rtl")

    async def orbit(self, radius: float, altitude: float) -> dict:
        # DJI has Point of Interest (POI) mode
        return await self._send_command("poi", {
            "radius": radius,
            "altitude": altitude
        })

    async def follow(self, target: str) -> dict:
        # DJI has ActiveTrack
        return await self._send_command("active_track", {"target": target})

    async def execute_mission(self, waypoints: List[dict]) -> dict:
        return await self._send_command("mission", {"waypoints": waypoints})

    async def start_video_stream(self) -> dict:
        # DJI video stream forwarded to LiveKit
        return await self._send_command("start_stream")

    async def telemetry_loop(self):
        """Poll telemetry from DJI Mobile SDK"""
        import aiohttp

        while True:
            try:
                async with aiohttp.ClientSession() as session:
                    async with session.get(
                        f"{self.mobile_sdk_url}/api/telemetry"
                    ) as resp:
                        telemetry = await resp.json()
                        await self.publish_telemetry(telemetry)

                await asyncio.sleep(1)

            except Exception as e:
                print(f"DJI telemetry error: {e}")
                await asyncio.sleep(5)
```

### Bridge 3: Skydio Cloud API (Security Drone)

```python
# agents/edge/drone_manager/bridges/skydio_bridge.py

from drone_manager.core import DroneBridge
import aiohttp

class SkydioBridge(DroneBridge):
    """Bridge for Skydio drones via Cloud API"""

    def __init__(self, api_key: str, vehicle_id: str):
        """
        Args:
            api_key: Skydio Cloud API key
            vehicle_id: Skydio vehicle ID
        """
        self.api_key = api_key
        self.vehicle_id = vehicle_id
        self.base_url = "https://api.skydio.com/api/v1"
        self.nc = None
        self.drone_id = None

    async def connect(self, nc, drone_id: str):
        """Connect to Skydio Cloud API and NATS"""
        self.nc = nc
        self.drone_id = drone_id

        # Verify API access
        async with aiohttp.ClientSession() as session:
            async with session.get(
                f"{self.base_url}/vehicles/{self.vehicle_id}",
                headers={"Authorization": f"Bearer {self.api_key}"}
            ) as resp:
                if resp.status != 200:
                    raise Exception("Failed to connect to Skydio API")

        # Start telemetry loop
        asyncio.create_task(self.telemetry_loop())

    async def _api_request(self, method: str, endpoint: str, data: dict = None) -> dict:
        """Make API request to Skydio Cloud"""
        async with aiohttp.ClientSession() as session:
            kwargs = {
                "headers": {"Authorization": f"Bearer {self.api_key}"}
            }

            if data:
                kwargs["json"] = data

            async with session.request(
                method,
                f"{self.base_url}{endpoint}",
                **kwargs
            ) as resp:
                return await resp.json()

    async def takeoff(self, altitude: float) -> dict:
        return await self._api_request(
            "POST",
            f"/vehicles/{self.vehicle_id}/mission/takeoff",
            {"altitude": altitude}
        )

    async def land(self) -> dict:
        return await self._api_request(
            "POST",
            f"/vehicles/{self.vehicle_id}/mission/land"
        )

    async def goto(self, lat: float, lon: float, alt: float) -> dict:
        return await self._api_request(
            "POST",
            f"/vehicles/{self.vehicle_id}/mission/goto",
            {"latitude": lat, "longitude": lon, "altitude": alt}
        )

    async def rtl(self) -> dict:
        return await self._api_request(
            "POST",
            f"/vehicles/{self.vehicle_id}/mission/rtl"
        )

    async def orbit(self, radius: float, altitude: float) -> dict:
        return await self._api_request(
            "POST",
            f"/vehicles/{self.vehicle_id}/skills/orbit",
            {"radius": radius, "altitude": altitude}
        )

    async def follow(self, target: str) -> dict:
        # Skydio's best feature - AI tracking
        return await self._api_request(
            "POST",
            f"/vehicles/{self.vehicle_id}/skills/track",
            {"subject": target}
        )

    async def execute_mission(self, waypoints: List[dict]) -> dict:
        return await self._api_request(
            "POST",
            f"/vehicles/{self.vehicle_id}/mission/waypoints",
            {"waypoints": waypoints}
        )

    async def start_video_stream(self) -> dict:
        # Skydio video stream URL
        stream_url = await self._api_request(
            "GET",
            f"/vehicles/{self.vehicle_id}/video/stream_url"
        )

        # Forward to LiveKit (would need video pipeline)
        return {"status": "streaming", "url": stream_url}

    async def telemetry_loop(self):
        """Poll telemetry from Skydio Cloud API"""
        while True:
            try:
                telemetry = await self._api_request(
                    "GET",
                    f"/vehicles/{self.vehicle_id}/telemetry"
                )

                await self.publish_telemetry(telemetry)

                await asyncio.sleep(1)

            except Exception as e:
                print(f"Skydio telemetry error: {e}")
                await asyncio.sleep(5)
```

---

## Legacy Agent Integration (Voice Control)

```python
# agents/legacy/tools/drone_control.py

from agno import AgnoTool
import nats

class DroneControlTool(AgnoTool):
    """
    Universal drone control tool for Legacy agent

    Works with any drone type via unified OneMind Drone Manager API
    """

    def __init__(self):
        self.nc = None

    async def connect(self):
        """Connect to NATS"""
        self.nc = await nats.connect("nats://100.102.21.44:4222")

    async def command_drone(
        self,
        drone_name: str,
        action: str,
        params: dict = None
    ) -> str:
        """
        Control any drone in the fleet

        Args:
            drone_name: "zeus", "farm", or "security"
            action: "takeoff", "land", "goto", "rtl", "orbit", "follow", "patrol"
            params: Action-specific parameters

        Examples:
            await command_drone("zeus", "takeoff", {"altitude": 15})
            await command_drone("farm", "goto", {"latitude": 40.123, "longitude": -74.456, "altitude": 20})
            await command_drone("security", "patrol")
        """
        # Map friendly names to drone IDs
        drone_map = {
            "zeus": "zeus",
            "farm": "farm",
            "security": "security"
        }

        drone_id = drone_map.get(drone_name.lower(), drone_name)

        # Send command via NATS
        subject = f"commands.farm.drone.{drone_id}.{action}"

        try:
            response = await self.nc.request(
                subject,
                json.dumps(params or {}).encode(),
                timeout=10.0
            )

            result = json.loads(response.data.decode())
            return f"Drone {drone_name}: {result}"

        except Exception as e:
            return f"Error commanding {drone_name}: {str(e)}"

    async def get_drone_status(self, drone_name: str = None) -> str:
        """
        Get status of one or all drones

        Args:
            drone_name: Specific drone, or None for all drones
        """
        if drone_name:
            drone_id = {"zeus": "zeus", "farm": "farm", "security": "security"}.get(
                drone_name.lower(), drone_name
            )

            response = await self.nc.request(
                f"queries.farm.drone.{drone_id}.status",
                b"",
                timeout=5.0
            )

            status = json.loads(response.data.decode())

            return f"""
            Drone: {drone_name}
            Status: {status['status']}
            Position: {status['position']['latitude']}, {status['position']['longitude']}
            Altitude: {status['position']['altitude']}m
            Battery: {status['telemetry']['battery']}%
            Signal: {status['telemetry']['signal']}%
            """

        else:
            # Get all drones
            response = await self.nc.request(
                "queries.farm.drone.fleet.status",
                b"",
                timeout=5.0
            )

            fleet = json.loads(response.data.decode())

            summary = "Fleet Status:\n"
            for drone in fleet:
                summary += f"- {drone['drone_id']}: {drone['status']} ({drone['telemetry']['battery']}% battery)\n"

            return summary

    async def coordinate_mission(
        self,
        mission_type: str,
        drones: List[str],
        params: dict = None
    ) -> str:
        """
        Coordinate multi-drone mission

        Args:
            mission_type: "patrol", "search", "formation"
            drones: List of drone names
            params: Mission parameters
        """
        drone_ids = [
            {"zeus": "zeus", "farm": "farm", "security": "security"}.get(d.lower(), d)
            for d in drones
        ]

        response = await self.nc.request(
            "commands.farm.drone.fleet.coordinate",
            json.dumps({
                "mission_type": mission_type,
                "drones": drone_ids,
                "params": params or {}
            }).encode(),
            timeout=10.0
        )

        result = json.loads(response.data.decode())
        return f"Coordinated mission started: {result}"

# Add to Legacy agent
from agents.legacy.tools.drone_control import DroneControlTool

legacy_agent = Agent(
    name="Legacy",
    tools=[
        DroneControlTool(),
        # ... other tools
    ]
)
```

### Voice Command Examples

```python
# User says: "Legacy, launch Zeus for a perimeter check"
await drone_tool.command_drone("zeus", "takeoff", {"altitude": 15})
await drone_tool.command_drone("zeus", "mission", {
    "waypoints": [
        {"name": "north", "latitude": 40.123, "longitude": -74.456, "altitude": 15},
        {"name": "east", "latitude": 40.124, "longitude": -74.457, "altitude": 15},
        {"name": "south", "latitude": 40.125, "longitude": -74.458, "altitude": 15},
        {"name": "west", "latitude": 40.126, "longitude": -74.459, "altitude": 15},
    ]
})

# User says: "Start spraying the east field"
await drone_tool.command_drone("farm", "takeoff", {"altitude": 5})
await drone_tool.command_drone("farm", "mission", {
    "waypoints": get_spray_pattern("east_field")
})

# User says: "Skydio, check the north fence"
await drone_tool.command_drone("security", "takeoff")
await drone_tool.command_drone("security", "goto", {
    "latitude": 40.130,
    "longitude": -74.460,
    "altitude": 20
})

# User says: "What's the drone status?"
status = await drone_tool.get_drone_status()
# Legacy speaks: "Zeus is flying at 15 meters, 87% battery. Farm drone is idle on the pad. Security drone is on patrol, 62% battery."
```

---

## Cost Breakdown

### Hardware Costs (One-Time)

| Component | Cost | Notes |
|-----------|------|-------|
| **Zeus Drone (DIY)** | $765 | Pixhawk 6C + S500 frame + all components |
| **Farm Drone (DJI Agras T40)** | $16,000 | Agricultural spraying drone |
| **Security Drone (Skydio 2+)** | $1,099 | Autonomous patrol drone |
| **Companion Devices** | $300 | Tablet for DJI, radio telemetry spares |
| **Edge Compute Device** | $997 | OpenMind Edge OR Jetson Orin Nano ($500) |
| **Batteries (spare sets)** | $500 | Extra batteries for all drones |
| **Landing Pads** | $150 | 3x landing pads |
| **Charging Station** | $300 | Multi-drone charging hub |
| **Total Hardware** | **$20,111** | **Or $19,614 with DIY edge device** |

### Monthly Services

| Service | Cost | Notes |
|---------|------|-------|
| AWS Bedrock (Claude) | $50 | One AI brain for all drones |
| Deepgram STT | $20 | Voice transcription |
| Cartesia TTS | $10 | Voice synthesis |
| NATS/FABRIC | $0 | Self-hosted |
| LiveKit | $0 | Self-hosted video streaming |
| Tailscale Pro | $5 | Network for all locations |
| Skydio Cloud API | $0 | Included with drone |
| **Total Monthly** | **$85/mo** | **Does NOT scale with drone count** |

### ROI Analysis

**Use Cases:**
1. **Zeus (Personal):** Package delivery, errands, experimentation - **$50/task value**
2. **Farm Drone:** 15-20 acres/hour spraying vs manual ($200/acre manual) - **$3,000-4,000/day value**
3. **Security:** 24/7 patrol vs hiring guard ($20/hr × 24hr = $480/day) - **$14,400/month value**

**Break-even:** 1-2 months if using farm drone regularly, or 3-4 months with security + personal use only.

---

## Implementation Timeline

### Week 1-2: Foundation
- [ ] Set up edge device (OpenMind Edge or Jetson)
- [ ] Deploy NATS Leaf node
- [ ] Deploy OneMind Drone Manager base
- [ ] Test cloud ↔ edge messaging

### Week 3-4: Zeus Integration (MAVLink)
- [ ] Build Zeus drone (DIY Pixhawk)
- [ ] Deploy MAVLink bridge
- [ ] Test basic commands (takeoff, land, goto)
- [ ] Integrate video streaming to LiveKit
- [ ] Add to Legacy agent tools

### Week 5-6: Farm Drone (DJI)
- [ ] Set up DJI Agras T40
- [ ] Create DJI Mobile SDK companion app
- [ ] Deploy DJI bridge
- [ ] Test spraying mission
- [ ] Create agricultural workflows

### Week 7-8: Security Drone (Skydio)
- [ ] Set up Skydio 2+
- [ ] Deploy Skydio Cloud API bridge
- [ ] Test patrol missions
- [ ] Create security workflows
- [ ] Test coordinated missions

### Week 9-10: Integration & Testing
- [ ] Voice control testing (all drones)
- [ ] Multi-drone coordination testing
- [ ] Video streaming from all drones
- [ ] Emergency procedures
- [ ] Documentation

**Total: 8-10 weeks to full fleet operation**

---

## Next Steps

1. **Hardware Purchase Decision**
   - Zeus DIY parts ($765)
   - DJI Agras T40 ($16,000) OR start with cheaper DJI Mini for testing
   - Skydio 2+ ($1,099)
   - Edge device ($997 or $500 DIY)

2. **Create Linear Issue**
   - Import this document into Linear
   - Break into 4 main phases (Foundation, Zeus, Farm, Security)
   - Assign estimated effort

3. **Start with Zeus (Cheapest)**
   - Build DIY drone first
   - Validate entire architecture
   - Then add commercial drones

4. **Parallel Development**
   - Build Zeus hardware while coding Drone Manager
   - Test with simulated drones first
   - Deploy to real hardware incrementally

---

## Success Criteria

### MVP (Week 4)
- [x] Zeus drone responds to voice commands
- [x] "Legacy, launch Zeus" → drone takes off
- [x] Video stream visible in UI
- [x] Basic telemetry (battery, position)

### V1 (Week 8)
- [x] All 3 drones operational
- [x] Voice control for each drone type
- [x] Automated patrol workflow (Security)
- [x] Automated spraying workflow (Farm)
- [x] Multi-drone status dashboard

### V2 (Week 12+)
- [x] Coordinated multi-drone missions
- [x] Computer vision analysis (detect threats, count plants, etc.)
- [x] Predictive maintenance (battery cycles, flight hours)
- [x] Integration with ground robots (full FABRIC network)

---

**Ready to build a multi-drone fleet with OneMind OS! 🚁**

*Document inspired by OpenMind.org's OM1 architecture*
*Hardware-agnostic design for maximum flexibility*
*One edge device, multiple drones, unified control*
