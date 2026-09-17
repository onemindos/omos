# Meshtastic Mesh Network Integration with OneMind OS

> **Critical Resilience Layer for Drone/Robot Fleet**
>
> Offline-first communication backbone using LoRa mesh technology

---

## Table of Contents

1. [Overview](#overview)
2. [What is Meshtastic?](#what-is-meshtastic)
3. [Why OneMind OS Needs Meshtastic](#why-onemind-os-needs-meshtastic)
4. [Architecture](#architecture)
5. [MQTT Bridge to NATS](#mqtt-bridge-to-nats)
6. [Integration with Existing Fleet](#integration-with-existing-fleet)
7. [Implementation Plan](#implementation-plan)
8. [Cost Analysis](#cost-analysis)
9. [Code Implementation](#code-implementation)
10. [Emergency Scenarios](#emergency-scenarios)

---

## Overview

Meshtastic provides a **critical resilience layer** for the OneMind OS fleet by adding offline-capable, long-range communication that works when internet fails. This document covers complete integration with the existing drone/robot fleet architecture.

**Key Benefits:**
- 100% offline operation during internet outages
- 5-10km typical range (331km world record)
- Redundant communication paths (internet + mesh)
- Decentralized mesh (no single point of failure)
- Low power consumption (weeks to months on battery)
- Encrypted by default
- No recurring costs (no cellular plans)

---

## What is Meshtastic?

### Core Technology

**Meshtastic** is an open-source LoRa mesh network that enables long-range, off-grid communication:

| Specification | Details |
|--------------|---------|
| **Radio Technology** | LoRa (Long Range) |
| **Frequency** | 915MHz (US), 868MHz (EU), 433MHz (Asia) |
| **Typical Range** | 5-10km (open terrain) |
| **Max Range Record** | 331km (by MartinR7 & alleg) |
| **Network Type** | Decentralized mesh |
| **Power Consumption** | ~50mA active, <1mA sleep (nRF52) |
| **Battery Life** | 2-4 weeks typical, months with solar |
| **Encryption** | AES-256 (optional per channel) |
| **License Required** | No (ISM band) |
| **Cost** | $30-70 per device |

### How It Works

1. **Mesh Networking**: Devices automatically relay messages to extend range
   ```
   Device A → Device B → Device C → Device D
   ```
   If A and D are 20km apart, messages hop through B and C

2. **Multiple Channels**: Support for public/private channels with separate encryption

3. **GPS Tracking**: Optional position reporting for asset tracking

4. **Text Messaging**: Person-to-person or broadcast to all mesh users

5. **Telemetry**: Battery, signal strength, temperature sensors

6. **MQTT Bridge**: Internet-connected devices can bridge mesh ↔ MQTT server

---

## Why OneMind OS Needs Meshtastic

### Current Architecture Vulnerability

OneMind OS currently depends on internet connectivity:

```
Cloud NATS Hub (AWS)
      ↓
  (Internet)
      ↓
Edge Devices (Farm, Barn, House)
      ↓
Drones, Robots, Sensors
```

**Problem**: If internet fails, the entire fleet loses cloud connection.

### Meshtastic Solution

Add offline mesh layer that works independently:

```
Cloud NATS Hub (AWS)
      ↓
  (Internet - CAN FAIL)
      ↓
Edge Devices ←───(LoRa Mesh - ALWAYS WORKS)───→ All Assets
      ↓
  Drones, Robots, Sensors
```

**Result**: Fleet continues operation even during internet outages.

---

## Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Cloud (Azure VM)                    │
│  ┌──────────┐         ┌──────────┐         ┌──────────┐    │
│  │  Legacy  │←────────│ NATS Hub │────────→│    UI    │    │
│  │  Agent   │         │ :4222    │         │  :3000   │    │
│  └──────────┘         └──────────┘         └──────────┘    │
└─────────────────────────────────────────────────────────────┘
                              ↑
                    (Internet - Primary)
                              ↓
┌─────────────────────────────────────────────────────────────┐
│              House Edge Device (Jetson Orin Nano)            │
│  ┌──────────┐  ┌──────────────┐  ┌────────────────────┐   │
│  │   NATS   │  │   Frigate    │  │ Meshtastic Bridge  │   │
│  │   Leaf   │  │ (6 cameras)  │  │ (mesh ↔ NATS)      │   │
│  │  :4222   │  │   :5000      │  │    :1883           │   │
│  └──────────┘  └──────────────┘  └────────────────────┘   │
│       ↑                                    ↑                │
│       │                                    │                │
│       └────────────────────────────────────┘                │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │          Meshtastic Device (RAK Starter Kit)        │   │
│  │          - LoRa Radio (SX1262)                       │   │
│  │          - GPS Module                                │   │
│  │          - MQTT Server on :1883                      │   │
│  │          - WiFi → Bridge to internet                 │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                   (LoRa Mesh - Backup/Always)
                              │
        ┌─────────────────────┼──────────────────────────┐
        │                     │                          │
        ↓                     ↓                          ↓
┌───────────────┐   ┌─────────────────┐   ┌──────────────────┐
│ Farm Edge     │   │ Zeus Drone      │   │ Mobile Tracker   │
│ (Barn)        │   │ (in flight)     │   │ (Truck)          │
│               │   │                 │   │                  │
│ Jetson Nano   │   │ Raspberry Pi 4  │   │ Handheld Device  │
│ - NATS Leaf   │   │ - Flight Ctrl   │   │ - GPS            │
│ - Meshtastic  │   │ - Meshtastic    │   │ - Meshtastic     │
│ - Unitree x2  │   │ - Camera        │   │ - Display        │
│ - Sensors     │   │ - Telemetry     │   │ - Battery        │
└───────────────┘   └─────────────────┘   └──────────────────┘
        │                     │                          │
        ↓                     ↓                          ↓
┌───────────────┐   ┌─────────────────┐   ┌──────────────────┐
│ Meshtastic    │   │ Meshtastic      │   │ Meshtastic       │
│ RAK19007      │   │ LILYGO T3-S3    │   │ RAK19007         │
│ + SX1262      │   │ + SX1262        │   │ + SX1262         │
│ GPS, Battery  │   │ Compact module  │   │ GPS, Battery     │
└───────────────┘   └─────────────────┘   └──────────────────┘
```

### Communication Paths

**Normal Operation (Internet Available):**
```
Drone → Edge NATS Leaf → Cloud NATS Hub → Legacy Agent
```

**Backup Path (Internet Down):**
```
Drone → Meshtastic Mesh → House Meshtastic → NATS Bridge → Local Processing
```

**Hybrid (Redundant):**
```
Critical Alert:
├─ Send via Internet (fast)
└─ Send via Mesh (reliable)

Legacy receives via both paths (whichever arrives first wins)
```

---

## MQTT Bridge to NATS

### Bridge Architecture

Meshtastic devices with WiFi/Ethernet can run MQTT server on port 1883. The bridge service connects Meshtastic MQTT to OneMind NATS:

```
┌─────────────────────────────────────────────────────────┐
│              Meshtastic MQTT Topics                      │
├─────────────────────────────────────────────────────────┤
│ mesh/+/position     - GPS coordinates                    │
│ mesh/+/telemetry    - Battery, signal, sensors          │
│ mesh/+/text         - Text messages                      │
│ mesh/+/nodeinfo     - Device info                        │
└─────────────────────────────────────────────────────────┘
                          ↓
                 (Meshtastic Bridge)
                          ↓
┌─────────────────────────────────────────────────────────┐
│                   NATS Subjects                          │
├─────────────────────────────────────────────────────────┤
│ events.mesh.position.{device_id}                         │
│ events.mesh.telemetry.{device_id}                        │
│ events.mesh.text.{device_id}                             │
│ events.mesh.nodeinfo.{device_id}                         │
└─────────────────────────────────────────────────────────┘
```

### Message Flow

1. **Mesh → NATS (Inbound)**
   ```
   Drone sends GPS via LoRa
      → Farm Edge Meshtastic receives
      → Publishes to local MQTT
      → Bridge forwards to NATS
      → Legacy processes event
   ```

2. **NATS → Mesh (Outbound)**
   ```
   Legacy sends command
      → NATS publishes to commands.mesh.{device_id}
      → Bridge subscribes and receives
      → Bridge publishes to Meshtastic MQTT
      → Meshtastic broadcasts via LoRa
      → Drone receives and executes
   ```

---

## Integration with Existing Fleet

### 1. Zeus Drone (DIY Pixhawk)

**Hardware Addition:**
- LILYGO T3-S3 LoRa module: $35
- Connects to Raspberry Pi 4 via USB serial
- Powered from drone battery (5V regulated)

**Software Integration:**
```python
# agents/edge/drone_manager/bridges/mavlink_bridge.py

class MAVLinkBridge(DroneBridge):
    def __init__(self, connection_string: str, mesh_device: str = None):
        self.connection_string = connection_string
        self.drone = None

        # Add mesh backup
        if mesh_device:
            self.mesh = MeshtasticClient(mesh_device)

    async def send_telemetry(self, data: dict):
        """Send telemetry with automatic mesh fallback"""
        try:
            # Primary: Send via NATS
            await self.nats.publish("telemetry.zeus", json.dumps(data))
        except Exception as e:
            # Backup: Send via Meshtastic mesh
            self.mesh.send_text(f"ZEUS_TELEMETRY:{json.dumps(data)}")
            logger.warning(f"Using mesh backup: {e}")
```

### 2. DJI Agras T40 (Agricultural Drone)

**Hardware Addition:**
- RAK19007 Starter Kit: $65 (rugged for field use)
- Mounted on drone frame
- Connects to flight controller telemetry port

**Use Case:**
- Long-range field operations (5-10km from base)
- Continues reporting even if cellular fails
- Emergency abort commands via mesh

### 3. Skydio 2+ (Security Drone)

**Hardware Addition:**
- Compact LILYGO T3-S3: $35
- Lightweight (20g) for minimal impact on flight time
- Broadcasts security alerts via mesh

**Use Case:**
- Security patrol continues during internet outage
- Critical alerts sent via both internet AND mesh
- Tracks intruders even if they jam WiFi

### 4. Unitree Go1 Robots

**Hardware Addition:**
- RAK19007 kit mounted on robot back: $65
- Powered from robot battery (12V → 5V regulator)
- GPS for outdoor navigation

**Use Case:**
- Barn patrol continues offline
- Status updates via mesh if WiFi range exceeded
- Emergency stop commands via mesh

### 5. Fixed Sensor Nodes

**Hardware Addition:**
- ESP32-S3 + SX1262 LoRa: $40 per node
- Solar panel + battery: $30 per node
- Sensors (soil, weather, gates): $20-50 per node

**Deployment:**
- North Field: Soil moisture + weather
- South Pasture: Gate sensors + weather
- Barn: Temperature + door sensors
- Water Tank: Level sensor
- Equipment: GPS trackers on tractors, vehicles

**Total: 10 nodes × $90 = $900**

---

## Implementation Plan

### Phase 1: Test Mesh Network (Week 1)

**Goal**: Validate LoRa range and mesh reliability across property

**Hardware:**
- 2× RAK19007 Starter Kit ($130 total)
- Batteries and enclosures ($50 total)

**Tasks:**
1. Flash latest Meshtastic firmware (web flasher)
2. Configure channel with same PSK
3. Enable GPS on both devices
4. Test range: Walk from house to barn to fields
5. Measure signal strength and hop count
6. Test text messaging reliability
7. Map coverage zones

**Success Criteria:**
- 5km+ range achieved
- <5% packet loss
- <2s message latency
- GPS accuracy within 10m

### Phase 2: MQTT Bridge to NATS (Week 2)

**Goal**: Connect Meshtastic mesh to OneMind NATS backbone

**Software:**
- Meshtastic Bridge service (Python)
- Docker container for edge devices
- NATS JetStream persistence

**Tasks:**
1. Create MeshtasticBridge service (see code below)
2. Add to edge compose file
3. Configure Meshtastic device for MQTT mode
4. Test bidirectional messaging (NATS ↔ Mesh)
5. Implement message persistence (JetStream)
6. Add monitoring (Prometheus metrics)

**Success Criteria:**
- Messages flow NATS → Mesh → NATS
- <1s end-to-end latency
- 100% message delivery (with retries)
- Automatic reconnection on failure

### Phase 3: Drone/Robot Integration (Week 3-4)

**Goal**: Add mesh backup to all mobile assets

**Hardware:**
- 3× Drone LoRa modules ($150 total)
- 2× Robot LoRa modules ($130 total)

**Software:**
- TelemetryBackup class
- CommandFailover class
- Integration with existing bridges

**Tasks:**
1. Mount LoRa hardware on each asset
2. Add mesh telemetry backup to each bridge
3. Add mesh command failover
4. Test internet-down scenarios
5. Test emergency recall via mesh
6. Validate GPS tracking during flight

**Success Criteria:**
- Telemetry continues during internet outage
- Commands execute via mesh with <3s latency
- GPS tracking works for all assets
- Emergency recalls succeed 100%

### Phase 4: Full Mesh Deployment (Week 5-8)

**Goal**: Complete property coverage with sensor network

**Hardware:**
- 10× Sensor nodes ($900 total)
- 2× Mobile trackers ($120 total)
- 3× Additional edge modules ($150 total)

**Software:**
- Sensor data collectors
- Asset tracking dashboard
- Alert routing via mesh

**Tasks:**
1. Deploy sensor nodes at key locations
2. Configure solar power and enclosures
3. Add mobile trackers to vehicles/equipment
4. Implement sensor data processing
5. Create asset tracking UI
6. Test full-mesh resilience (disable internet)
7. Document coverage map and procedures

**Success Criteria:**
- 100% property coverage
- All sensors reporting every 5 minutes
- Asset tracking accurate to 10m
- Fleet operates 100% offline for 24+ hours
- Automatic sync when internet returns

---

## Cost Analysis

### Starter Phase (Test & Validate)

| Item | Qty | Unit Price | Total | Purpose |
|------|-----|------------|-------|---------|
| RAK19007 Starter Kit | 2 | $65 | $130 | House + Barn hubs |
| 18650 Batteries | 4 | $5 | $20 | Backup power |
| IP65 Enclosures | 2 | $15 | $30 | Weather protection |
| **Subtotal** | | | **$180** | Initial testing |

### Production Fleet

| Item | Qty | Unit Price | Total | Purpose |
|------|-----|------------|-------|---------|
| **Mobile Assets** | | | | |
| Zeus Drone module | 1 | $35 | $35 | LILYGO T3-S3 |
| DJI Agras module | 1 | $65 | $65 | RAK19007 (rugged) |
| Skydio module | 1 | $35 | $35 | LILYGO T3-S3 (light) |
| Unitree Robot modules | 2 | $65 | $130 | RAK19007 each |
| **Fixed Infrastructure** | | | | |
| Edge device modules | 3 | $50 | $150 | Barn, Field, Office |
| Sensor nodes (ESP32+LoRa) | 10 | $40 | $400 | Fields, gates, barn |
| Solar panels for sensors | 10 | $30 | $300 | Self-powered nodes |
| Mobile trackers | 2 | $60 | $120 | Trucks, equipment |
| **Subtotal** | | | **$1,235** | Full deployment |
| | | | | |
| **Grand Total** | | | **$1,415** | Complete mesh network |

### Comparison to Alternatives

| Solution | Setup Cost | Monthly Cost | Range | Offline |
|----------|------------|--------------|-------|---------|
| **Meshtastic** | **$1,415** | **$0** | **10km** | **Yes** |
| Cellular (1 modem per asset) | $600 | $200-300 | Unlimited | No |
| LoRaWAN (gateway) | $2,000 | $50-100 | 15km | Partial |
| Private LTE | $10,000+ | $100+ | 5km | No |
| Satellite (Starlink) | $5,000+ | $500+ | Unlimited | No |

**Winner**: Meshtastic - lowest cost, fully offline, adequate range

---

## Code Implementation

### Meshtastic Bridge Service

**File**: `agents/edge/meshtastic_bridge/main.py`

```python
#!/usr/bin/env python3
"""
Meshtastic MQTT Bridge to NATS
Connects Meshtastic mesh network to OneMind NATS backbone
"""

import asyncio
import json
import os
from typing import Dict, Optional
import paho.mqtt.client as mqtt
import nats
from nats.js.api import StreamConfig
from loguru import logger
from meshtastic.protobuf import mesh_pb2, portnums_pb2

class MeshtasticBridge:
    """Bridge between Meshtastic MQTT and NATS"""

    def __init__(
        self,
        nats_url: str = "nats://localhost:4222",
        mqtt_host: str = "localhost",
        mqtt_port: int = 1883,
        location: str = "house"
    ):
        self.nats_url = nats_url
        self.mqtt_host = mqtt_host
        self.mqtt_port = mqtt_port
        self.location = location

        # Clients
        self.nc: Optional[nats.NATS] = None
        self.js = None
        self.mqtt_client: Optional[mqtt.Client] = None

        # State
        self.device_map: Dict[str, str] = {}  # mesh_id → device_name
        self.running = False

    async def start(self):
        """Start the bridge"""
        logger.info(f"Starting Meshtastic bridge at {self.location}")

        # Connect to NATS
        self.nc = await nats.connect(self.nats_url)
        self.js = self.nc.jetstream()

        # Create stream for mesh events
        try:
            await self.js.add_stream(
                StreamConfig(
                    name="MESH_EVENTS",
                    subjects=["events.mesh.>"],
                    retention="limits",
                    max_age=86400 * 7,  # 7 days
                )
            )
        except Exception:
            pass  # Stream already exists

        # Connect to Meshtastic MQTT
        self.mqtt_client = mqtt.Client(client_id=f"bridge_{self.location}")
        self.mqtt_client.on_connect = self._on_mqtt_connect
        self.mqtt_client.on_message = self._on_mqtt_message

        logger.info(f"Connecting to Meshtastic MQTT at {self.mqtt_host}:{self.mqtt_port}")
        self.mqtt_client.connect(self.mqtt_host, self.mqtt_port, 60)

        # Subscribe to NATS commands (for outbound to mesh)
        await self.nc.subscribe("commands.mesh.>", cb=self._on_nats_command)

        # Start MQTT loop in thread
        self.mqtt_client.loop_start()

        self.running = True
        logger.info("Bridge started successfully")

    def _on_mqtt_connect(self, client, userdata, flags, rc):
        """MQTT connection callback"""
        if rc == 0:
            logger.info("Connected to Meshtastic MQTT")
            # Subscribe to all mesh topics
            client.subscribe("msh/+/+/#")
            client.subscribe("mesh/+/position")
            client.subscribe("mesh/+/telemetry")
            client.subscribe("mesh/+/text")
            client.subscribe("mesh/+/nodeinfo")
        else:
            logger.error(f"MQTT connection failed: {rc}")

    def _on_mqtt_message(self, client, userdata, msg):
        """MQTT message callback - forward to NATS"""
        try:
            # Parse topic
            parts = msg.topic.split("/")
            if len(parts) < 3:
                return

            # Get device ID and message type
            device_id = parts[1]
            msg_type = parts[2]

            # Parse payload (protobuf or JSON)
            try:
                payload = json.loads(msg.payload.decode())
            except:
                # Try protobuf decode
                payload = self._decode_protobuf(msg.payload, msg_type)

            if not payload:
                return

            # Get friendly device name
            device_name = self.device_map.get(device_id, device_id)

            # Construct NATS subject
            subject = f"events.mesh.{msg_type}.{device_name}"

            # Add metadata
            envelope = {
                "device_id": device_id,
                "device_name": device_name,
                "location": self.location,
                "type": msg_type,
                "data": payload,
                "timestamp": asyncio.get_event_loop().time()
            }

            # Publish to NATS
            asyncio.create_task(
                self.nc.publish(subject, json.dumps(envelope).encode())
            )

            logger.debug(f"Forwarded {msg_type} from {device_name} to NATS")

        except Exception as e:
            logger.error(f"Error processing MQTT message: {e}")

    async def _on_nats_command(self, msg):
        """NATS command callback - forward to mesh"""
        try:
            # Parse subject: commands.mesh.{action}.{device_name}
            parts = msg.subject.split(".")
            if len(parts) < 4:
                return

            action = parts[2]
            device_name = parts[3]

            # Parse payload
            try:
                data = json.loads(msg.data.decode())
            except:
                logger.error("Invalid command payload")
                return

            # Get mesh device ID
            device_id = self._get_device_id(device_name)
            if not device_id:
                logger.error(f"Unknown device: {device_name}")
                return

            # Forward to Meshtastic
            if action == "text":
                # Send text message
                self.mqtt_client.publish(
                    f"mesh/{device_id}/text",
                    json.dumps({"text": data.get("message", "")})
                )
            elif action == "position":
                # Request position update
                self.mqtt_client.publish(
                    f"mesh/{device_id}/cmd",
                    json.dumps({"cmd": "get_position"})
                )
            elif action == "emergency_recall":
                # Send emergency command
                self.mqtt_client.publish(
                    f"mesh/{device_id}/text",
                    "EMERGENCY_RECALL"
                )
                logger.warning(f"Emergency recall sent to {device_name}")

            logger.info(f"Forwarded command {action} to {device_name}")

        except Exception as e:
            logger.error(f"Error processing NATS command: {e}")

    def _decode_protobuf(self, data: bytes, msg_type: str) -> Optional[dict]:
        """Decode Meshtastic protobuf messages"""
        try:
            if msg_type == "position":
                pos = mesh_pb2.Position()
                pos.ParseFromString(data)
                return {
                    "latitude": pos.latitude_i / 1e7,
                    "longitude": pos.longitude_i / 1e7,
                    "altitude": pos.altitude,
                    "time": pos.time
                }
            elif msg_type == "telemetry":
                telem = mesh_pb2.Telemetry()
                telem.ParseFromString(data)
                return {
                    "battery_level": telem.device_metrics.battery_level,
                    "voltage": telem.device_metrics.voltage,
                    "channel_utilization": telem.device_metrics.channel_utilization,
                    "air_util_tx": telem.device_metrics.air_util_tx
                }
            # Add more types as needed
        except Exception as e:
            logger.error(f"Protobuf decode error: {e}")

        return None

    def _get_device_id(self, device_name: str) -> Optional[str]:
        """Get mesh ID from friendly device name"""
        for mesh_id, name in self.device_map.items():
            if name == device_name:
                return mesh_id
        return None

    async def stop(self):
        """Stop the bridge"""
        logger.info("Stopping bridge...")
        self.running = False

        if self.mqtt_client:
            self.mqtt_client.loop_stop()
            self.mqtt_client.disconnect()

        if self.nc:
            await self.nc.drain()
            await self.nc.close()

        logger.info("Bridge stopped")

async def main():
    """Main entry point"""
    # Get config from environment
    nats_url = os.getenv("NATS_URL", "nats://localhost:4222")
    mqtt_host = os.getenv("MESHTASTIC_HOST", "localhost")
    mqtt_port = int(os.getenv("MESHTASTIC_PORT", "1883"))
    location = os.getenv("LOCATION", "house")

    # Create and start bridge
    bridge = MeshtasticBridge(nats_url, mqtt_host, mqtt_port, location)
    await bridge.start()

    # Run until interrupted
    try:
        while bridge.running:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        logger.info("Interrupted by user")
    finally:
        await bridge.stop()

if __name__ == "__main__":
    asyncio.run(main())
```

### Telemetry Backup Class

**File**: `agents/edge/drone_manager/telemetry_backup.py`

```python
"""
Automatic telemetry backup via Meshtastic mesh
Falls back to LoRa when internet fails
"""

import asyncio
import json
from typing import Optional
from loguru import logger
import serial

class MeshtasticClient:
    """Simple client for sending data via Meshtastic serial"""

    def __init__(self, port: str = "/dev/ttyUSB0"):
        self.port = port
        self.serial: Optional[serial.Serial] = None

    def connect(self):
        """Connect to Meshtastic device"""
        try:
            self.serial = serial.Serial(self.port, 115200, timeout=1)
            logger.info(f"Connected to Meshtastic on {self.port}")
        except Exception as e:
            logger.error(f"Failed to connect to Meshtastic: {e}")

    def send_text(self, message: str):
        """Send text message via mesh"""
        if not self.serial:
            self.connect()

        try:
            # Send text command to Meshtastic
            cmd = f"--sendtext '{message}'\n"
            self.serial.write(cmd.encode())
            logger.debug(f"Sent via mesh: {message[:50]}...")
        except Exception as e:
            logger.error(f"Failed to send via mesh: {e}")

class TelemetryBackup:
    """Automatic failover for drone telemetry"""

    def __init__(
        self,
        nats_client,
        mesh_device: str = "/dev/ttyUSB0",
        device_name: str = "drone"
    ):
        self.nats = nats_client
        self.mesh = MeshtasticClient(mesh_device)
        self.device_name = device_name

        # State
        self.primary_failed = False
        self.fail_count = 0
        self.last_success = asyncio.get_event_loop().time()

    async def send_telemetry(self, data: dict):
        """Send telemetry with automatic mesh fallback"""

        # Always try primary first
        try:
            await self.nats.publish(
                f"telemetry.{self.device_name}",
                json.dumps(data).encode()
            )

            # Reset failure tracking
            if self.primary_failed:
                logger.info("Primary connection restored")
                self.primary_failed = False
            self.fail_count = 0
            self.last_success = asyncio.get_event_loop().time()

        except Exception as e:
            self.fail_count += 1

            # Switch to backup after 3 failures
            if self.fail_count >= 3 and not self.primary_failed:
                self.primary_failed = True
                logger.warning(
                    f"Primary connection failed ({self.fail_count} attempts), "
                    f"using mesh backup"
                )

            # Send via mesh backup
            try:
                # Compress data for LoRa (limited bandwidth)
                compact = {
                    "d": self.device_name,  # device
                    "lat": data.get("latitude"),
                    "lon": data.get("longitude"),
                    "alt": data.get("altitude"),
                    "bat": data.get("battery", 0),
                    "sig": data.get("signal_strength", 0),
                    "t": int(asyncio.get_event_loop().time())
                }

                # Format as compact JSON string
                mesh_msg = f"TELEM:{json.dumps(compact, separators=(',', ':'))}"

                # Send via mesh
                self.mesh.send_text(mesh_msg)

                logger.debug(f"Sent telemetry via mesh backup")

            except Exception as mesh_error:
                logger.error(f"Both primary and mesh backup failed: {mesh_error}")

    async def send_emergency(self, alert_type: str, data: dict):
        """Send emergency alert via BOTH paths simultaneously"""

        # Send via both paths for critical alerts
        tasks = []

        # Primary (NATS)
        try:
            tasks.append(
                self.nats.publish(
                    f"events.emergency.{self.device_name}",
                    json.dumps({
                        "type": alert_type,
                        "device": self.device_name,
                        "data": data,
                        "timestamp": asyncio.get_event_loop().time()
                    }).encode()
                )
            )
        except Exception as e:
            logger.error(f"Failed to send emergency via NATS: {e}")

        # Backup (Mesh) - always send
        try:
            mesh_msg = f"EMERGENCY:{alert_type}:{json.dumps(data, separators=(',', ':'))}"
            self.mesh.send_text(mesh_msg)
        except Exception as e:
            logger.error(f"Failed to send emergency via mesh: {e}")

        # Wait for both
        if tasks:
            await asyncio.gather(*tasks, return_exceptions=True)

        logger.critical(f"Emergency alert sent: {alert_type}")
```

### Docker Compose Integration

**File**: `agents/edge/docker-compose.edge.yml` (add to existing)

```yaml
services:
  # ... existing services ...

  meshtastic-bridge:
    build: ./meshtastic_bridge
    container_name: meshtastic-bridge
    restart: unless-stopped
    environment:
      - NATS_URL=nats://nats-leaf:4222
      - MESHTASTIC_HOST=192.168.1.100  # IP of Meshtastic device
      - MESHTASTIC_PORT=1883
      - LOCATION=house  # or farm, barn, etc.
    depends_on:
      - nats-leaf
    networks:
      - edge-network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

---

## Emergency Scenarios

### Scenario 1: Internet Outage During Drone Mission

**Situation**: Zeus drone is 3km away on agricultural survey mission. Internet connection fails at house edge device.

**Without Meshtastic:**
```
1. Drone continues mission (last commands)
2. Telemetry lost - no position updates
3. Cannot send new commands
4. Cannot recall drone remotely
5. Drone must complete mission or timeout
6. Risk of fly-away if mission logic has bugs
```

**With Meshtastic:**
```
1. Drone detects primary connection loss
2. Switches to mesh telemetry backup automatically
3. Position updates continue via LoRa → House Meshtastic → Local NATS
4. Legacy detects internet down, switches to local processing
5. Can send recall command via NATS → Mesh → Drone
6. Drone receives recall, returns home safely
7. Full telemetry maintained throughout
```

**Code Flow:**
```python
# Drone side (agents/edge/drone_manager/telemetry_backup.py)
async def monitor_connection():
    while True:
        try:
            # Try to send heartbeat via NATS
            await nats.publish("heartbeat.zeus", "alive")
            primary_ok = True
        except:
            # Switch to mesh backup
            mesh.send_text("ZEUS:HEARTBEAT")
            primary_ok = False

        await asyncio.sleep(5)

# Legacy side (agents/legacy/tools/drone_tool.py)
@tool
async def emergency_recall_drone(drone_name: str):
    """Recall drone immediately - uses all available paths"""

    # Try primary (NATS → Edge → Drone)
    try:
        await nc.publish(f"commands.drone.{drone_name}.recall", "")
    except:
        pass  # Primary failed

    # Send via mesh backup
    await nc.publish(f"commands.mesh.emergency_recall.{drone_name}", "")

    return f"Emergency recall sent to {drone_name} via all paths"
```

### Scenario 2: Fire Detected by Sensor

**Situation**: Remote field sensor detects fire. Internet may be spotty.

**Message Flow:**
```
Field Sensor (Meshtastic + Smoke Detector)
    → LoRa broadcast to all devices
        → House Meshtastic receives
            → NATS Bridge → events.mesh.sensor.fire
                → Legacy processes alert
                    → Sends alerts:
                        - Push notification to phone
                        - SMS via Twilio
                        - Voice call
                        - Recalls all drones immediately
                        - Sends robots to safe zone
                        - Activates fire suppression systems
```

**Why Mesh Matters:**
- Fire may damage internet infrastructure
- Cellular may be overloaded in emergency
- LoRa mesh continues working
- Multiple paths ensure alert delivery

### Scenario 3: Security Intruder with WiFi Jammer

**Situation**: Intruder jams WiFi/cellular to disable security systems.

**Without Meshtastic:**
```
1. Skydio security drone loses connection
2. Cannot send alerts to owner
3. Cannot track intruder
4. Security system blind
```

**With Meshtastic:**
```
1. Skydio detects intrusion via AI
2. Tries WiFi (jammed - fails)
3. Automatically switches to LoRa mesh
4. Alert sent via mesh to house
5. House Meshtastic has wired ethernet (no WiFi jamming)
6. Alert reaches Legacy via NATS
7. Owner notified immediately
8. Police called automatically
9. All events recorded with video evidence
```

**Jam-Proof Communication:**
- LoRa uses 915MHz (different from WiFi 2.4/5GHz)
- Encrypted by default (AES-256)
- Mesh routing finds alternate paths
- Cannot be jammed without military-grade equipment

### Scenario 4: Multi-Day Internet Outage

**Situation**: Rural area loses internet for 48 hours due to storm damage.

**Fleet Behavior:**
```
Hour 0: Internet fails
  - All devices detect failure
  - Switch to mesh communication automatically
  - Continue normal operations

Hour 2: Morning patrol
  - Security drone patrols via mesh commands
  - Robots patrol barn
  - All telemetry buffered locally

Hour 12: Agricultural task
  - Farm drone sprays field (pre-scheduled)
  - GPS coordinates via mesh
  - Mission status via mesh

Hour 24: Status check
  - Legacy queries all devices via mesh
  - Full status report received
  - All systems nominal

Hour 48: Internet returns
  - Primary connection restored
  - Buffered data synced to cloud
  - Resume cloud-first operation
  - Zero data loss
```

**Local Processing During Outage:**
```python
# agents/edge/drone_manager/main.py

class LocalDroneManager:
    """Operates drones without cloud connection"""

    async def handle_outage_mode(self):
        """Run local workflows during internet outage"""

        # Switch to local LLM (Ollama)
        self.llm = create_model("ollama/phi3:mini")

        # Continue scheduled tasks
        while not self.internet_restored:
            # Morning patrol (scheduled)
            if self.is_morning_patrol_time():
                await self.execute_patrol("security", local=True)

            # Agricultural tasks (scheduled)
            if self.has_scheduled_farm_tasks():
                await self.execute_farm_task(local=True)

            # Emergency monitoring
            alerts = await self.check_mesh_alerts()
            if alerts:
                await self.handle_local_emergency(alerts)

            await asyncio.sleep(60)

        # Sync buffered data when internet returns
        await self.sync_to_cloud()
```

---

## Security Considerations

### Encryption

**Default**: AES-256 encryption per channel
```python
# Configure channel with encryption
meshtastic --setchan psk random  # Random 256-bit key
meshtastic --set channel.name "onemind-private"
```

**Best Practice**: Separate channels for different security levels
- Public channel: General mesh coordination
- Private channel: Drone/robot commands (encrypted)
- Emergency channel: Critical alerts only

### Authentication

**Device Registration:**
```python
# Only accept commands from known devices
ALLOWED_DEVICES = {
    "!a1b2c3d4": "house_hub",
    "!e5f6g7h8": "farm_hub",
    "!i9j0k1l2": "zeus_drone",
    "!m3n4o5p6": "security_drone"
}

def verify_device(device_id: str) -> bool:
    """Only process messages from registered devices"""
    return device_id in ALLOWED_DEVICES
```

### Physical Security

**Tamper Detection:**
- Meshtastic devices are physical (can be stolen/tampered)
- Mount in locked enclosures
- GPS tracking on all devices
- Alert if device goes offline unexpectedly

**Recommendations:**
1. Lock enclosures with GPS trackers inside
2. Monitor device battery levels (low battery = possible theft)
3. Geofencing alerts (device moved outside expected area)
4. Regular heartbeat checks (missing heartbeat = alert)

---

## Monitoring and Maintenance

### Health Checks

**Metrics to Track:**
```python
# agents/edge/meshtastic_bridge/metrics.py

from prometheus_client import Counter, Gauge, Histogram

# Message counters
mesh_messages_received = Counter(
    'mesh_messages_received_total',
    'Total mesh messages received',
    ['device', 'type']
)

mesh_messages_sent = Counter(
    'mesh_messages_sent_total',
    'Total mesh messages sent',
    ['device', 'type']
)

# Connection status
mesh_connection_status = Gauge(
    'mesh_connection_status',
    'Mesh connection status (1=up, 0=down)',
    ['device']
)

# Signal strength
mesh_signal_strength = Gauge(
    'mesh_signal_strength_dbm',
    'Mesh signal strength in dBm',
    ['device']
)

# Message latency
mesh_message_latency = Histogram(
    'mesh_message_latency_seconds',
    'Message delivery latency',
    ['path']  # nats, mesh, hybrid
)
```

### Dashboard

**Grafana Panels:**
1. **Mesh Network Map**: Visual map of all devices and signal strength
2. **Message Throughput**: Messages/second on mesh vs NATS
3. **Device Battery Levels**: All mesh devices battery status
4. **Signal Strength Heatmap**: Coverage map of property
5. **Failover Events**: Count of primary → mesh failovers
6. **Alert Timeline**: Emergency alerts sent via mesh

### Maintenance Tasks

**Weekly:**
- Check battery levels on all devices
- Review failover logs
- Test emergency recall via mesh

**Monthly:**
- Walk property to verify coverage
- Update Meshtastic firmware
- Test full internet-down scenario

**Quarterly:**
- Replace batteries on solar nodes
- Physical inspection of outdoor enclosures
- Review and update device registry

---

## Troubleshooting

### Issue: Poor Signal Strength

**Symptoms:**
- High packet loss (>10%)
- Slow message delivery (>5s)
- Frequent disconnections

**Solutions:**
1. **Add relay node**: Place device between weak signal points
2. **Antenna upgrade**: Replace with higher-gain antenna (3dBi → 5dBi)
3. **Elevation**: Mount higher (roof vs ground = 3-5x range)
4. **Line of sight**: Clear trees/buildings between devices

### Issue: Messages Not Forwarding

**Symptoms:**
- Messages reach mesh but not NATS
- Commands sent but not executed

**Solutions:**
1. **Check bridge logs**: `docker logs meshtastic-bridge`
2. **Verify MQTT connection**: `mosquitto_sub -h localhost -t mesh/#`
3. **Check NATS connectivity**: `nats sub events.mesh.>`
4. **Restart bridge**: `docker restart meshtastic-bridge`

### Issue: High Battery Drain

**Symptoms:**
- Devices dying faster than expected
- Solar panels not keeping up

**Solutions:**
1. **Adjust transmit power**: Reduce if overkill (30dBm → 20dBm)
2. **Increase sleep intervals**: Less frequent updates
3. **Optimize routing**: Reduce unnecessary hops
4. **Upgrade batteries**: Larger capacity or more solar panels

---

## Future Enhancements

### Satellite Integration (Planned)

Combine Meshtastic with satellite for ultimate resilience:

```
Mesh Network (local, 10km)
    ↓ (fallback)
Satellite Uplink (global, unlimited)
    ↓ (fallback)
Internet (primary)
```

**Potential Integrations:**
- **Swarm** ($5/mo per device) - IoT satellite network
- **Starlink** (backup internet) - Mobile units on vehicles
- **Iridium** (emergency) - SOS beacon functionality

### Autonomous Mesh Repair

Drones automatically reposition to fix mesh holes:

```python
@workflow
async def mesh_health_monitor():
    """Monitor mesh coverage and auto-deploy relay drones"""

    # Detect coverage gap
    if mesh_signal_strength < -120:  # Very weak
        # Deploy Zeus drone as relay
        await command_drone("zeus", "goto", relay_position)
        await command_drone("zeus", "loiter", duration=3600)

        # Zeus now acts as mesh relay for 1 hour
        logger.info("Zeus deployed as mesh relay")
```

### Integration with Other Protocols

**Potential Additions:**
- **Zigbee**: Short-range sensors
- **LoRaWAN**: Long-range fixed infrastructure
- **BLE**: Ultra-low-power beacons
- **Satellite IoT**: Global coverage (Swarm, Skylo)

---

## Conclusion

Meshtastic provides a **critical resilience layer** for the OneMind OS fleet:

**Key Benefits:**
- ✅ 100% offline operation capability
- ✅ Zero monthly costs (no cellular plans)
- ✅ Long range (5-10km typical, 331km max)
- ✅ Low power (weeks to months on battery)
- ✅ Encrypted and secure
- ✅ Integrates seamlessly with NATS via MQTT bridge
- ✅ Automatic failover (internet → mesh)
- ✅ Property-wide sensor network
- ✅ Asset tracking (drones, robots, vehicles)
- ✅ Emergency communication (works when everything else fails)

**Investment:**
- Starter: $180 (test and validate)
- Full deployment: $1,415 (complete property coverage)
- Monthly cost: $0 (no subscriptions)

**Timeline:**
- Week 1: Test mesh network
- Week 2: MQTT bridge to NATS
- Week 3-4: Drone/robot integration
- Week 5-8: Full property deployment

**Next Steps:**
1. Order 2× RAK Starter Kits for testing
2. Map property to plan optimal device placement
3. Test LoRa range across property
4. Implement MQTT bridge service
5. Integrate with Zeus drone (when hardware arrives)

This completes the **communication resilience layer** of the OneMind OS architecture, ensuring the fleet can operate through internet outages, cellular failures, and even intentional jamming.

---

**Document**: `MESHTASTIC-MESH-INTEGRATION.md`
**Author**: OneMind OS Architecture Team
**Date**: 2026-01-22
**Version**: 1.0
**Related Docs**:
- `MULTI-DRONE-FLEET-ARCHITECTURE.md` - Primary drone integration
- `UNITREE-OPENMIND-INTEGRATION.md` - Robot integration
- `FULL-FABRIC-FLEET-INTEGRATION.md` - Unified fleet management


---

## Related

- [[00-24 SO (Sovereign Operations)/06 Research (R&D) (The Lab)/06 Research|06 Research]]
