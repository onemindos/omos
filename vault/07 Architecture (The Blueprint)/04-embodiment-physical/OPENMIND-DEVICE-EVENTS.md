# Device Event System

**Status**: ✅ Implemented
**Version**: 1.0
**Date**: 2025-02-17

---

## Overview

The Device Event System treats hardware as event sources that publish to the NATS message bus. This enables loose coupling between devices and agents, where agents subscribe to device events and react autonomously.

### Architecture

```
┌─────────────────────────────────────────┐
│       Device Sources (Hardware)         │
│  - Cameras (MaixCAM, USB, IP)          │
│  - Sensors (Temperature, Motion, etc.)  │
│  - Smart Home (Lights, Doors, etc.)     │
└─────────────────────────────────────────┘
                  ↓ DeviceEvent
┌─────────────────────────────────────────┐
│      NATS Message Bus (Topics)          │
│  onemind.devices.camera.*              │
│  onemind.devices.sensor.*              │
│  onemind.devices.home.*                │
└─────────────────────────────────────────┘
                  ↓ Subscribe
┌─────────────────────────────────────────┐
│   Agent Listeners (React to Events)     │
│  - Security Agent → Person detected    │
│  - HVAC Agent → Temperature high       │
│  - Monitor Agent → USB device added    │
└─────────────────────────────────────────┘
```

---

## Key Concepts

### 1. **DeviceEvent**
Standard event structure for all device types:

```python
from backend.devices import DeviceEvent

event = DeviceEvent(
    event_type="person_detected",
    device_id="living_room_camera",
    device_type="camera",
    confidence=0.87,
    data={
        "class": "person",
        "bounding_box": [100, 150, 300, 450],
        "snapshot_url": "http://camera/snapshot.jpg"
    },
    metadata={
        "location": "living_room",
        "zone": "indoor",
        "floor": 1
    }
)

# Publish to NATS
topic = event.to_nats_topic()  # "onemind.devices.camera.person_detected"
await publish(topic, event.to_dict())
```

### 2. **DeviceSource**
Abstract base class for device integrations:

```python
from backend.devices import DeviceSource
from backend.bus import publish

class MyCamera(DeviceSource):
    def __init__(self, device_id: str):
        super().__init__(device_id, "camera")

    async def start(self):
        """Start monitoring device"""
        self.running = True
        asyncio.create_task(self._monitoring_loop())

    async def stop(self):
        """Stop monitoring device"""
        self.running = False

    async def get_status(self):
        """Return device status"""
        return {"running": self.running}

    async def _monitoring_loop(self):
        """Monitor device and publish events"""
        while self.running:
            # Detect something
            event = DeviceEvent(...)
            await publish(event.to_nats_topic(), event.to_dict())
            await asyncio.sleep(1)
```

### 3. **DeviceRegistry**
Global registry for managing devices:

```python
from backend.devices import get_device_registry

registry = get_device_registry()

# Register device
registry.register(camera)

# Start all devices
await registry.start_all()

# Get device status
status = await registry.get_all_status()

# Stop all devices
await registry.stop_all()
```

### 4. **DeviceEventListener**
Agent that reacts to device events:

```python
from backend.devices.listeners import DeviceEventListener

listener = DeviceEventListener(
    agent_id="assistant",
    topic_pattern="onemind.devices.camera.person_detected"
)

await listener.start()
```

---

## Quickstart

### 1. Register a Device

```python
from backend.devices import get_device_registry
from backend.devices.sources.camera import MaixCAMSource

# Create device
camera = MaixCAMSource(
    device_id="front_door_camera",
    maixcam_ip="192.168.1.100"
)

# Set metadata
camera.set_metadata(
    location="front_door",
    zone="outdoor",
    floor=1
)

# Register
registry = get_device_registry()
registry.register(camera)
```

### 2. Start Devices

```python
# Start all registered devices
await registry.start_all()
```

### 3. Subscribe to Events

Agents automatically subscribe via listeners configured in [`listeners.py`](./listeners.py):

```python
# Pre-configured listeners
DEVICE_LISTENERS = [
    DeviceEventListener(
        agent_id="assistant",
        topic_pattern="onemind.devices.camera.person_detected"
    ),
    # ... more listeners
]
```

### 4. Manual Subscription

```python
from backend.bus import subscribe

async for event in subscribe("onemind.devices.camera.*"):
    print(f"Event: {event['event_type']}")
    print(f"Device: {event['device_id']}")
    print(f"Confidence: {event['confidence']}")
```

---

## Available Device Sources

### Camera Sources

#### **CameraSource** (Generic)
Base implementation for any camera:

```python
from backend.devices.sources.camera import CameraSource

camera = CameraSource(
    device_id="my_camera",
    camera_url="http://192.168.1.100:8080/stream",
    detection_classes=["person", "car"],
    confidence_threshold=0.7,
    fps=10
)
```

#### **MaixCAMSource** (Sipeed MaixCAM)
Specialized for MaixCAM AI camera ($50):

```python
from backend.devices.sources.camera import MaixCAMSource

maixcam = MaixCAMSource(
    device_id="front_door",
    maixcam_ip="192.168.1.100"
)
```

**MaixCAM Specs**:
- K230 Dual-Core RISC-V CPU (1.6GHz)
- KPU Neural Network Accelerator (1TOPS)
- 30 FPS inference
- Built-in person detection model

#### **MockCameraSource** (Testing)
For testing without hardware:

```python
from backend.devices.sources.camera import MockCameraSource

mock = MockCameraSource("test_camera")
mock.set_mock_detections([
    {
        "class": "person",
        "confidence": 0.85,
        "bbox": [100, 150, 300, 450]
    }
])
```

---

## NATS Topics

All device events publish to structured NATS topics:

```
onemind.devices.{device_type}.{event_type}
```

### Camera Events
```
onemind.devices.camera.person_detected
onemind.devices.camera.motion_detected
onemind.devices.camera.face_detected
onemind.devices.camera.vehicle_detected
```

### Sensor Events
```
onemind.devices.sensor.temperature_high
onemind.devices.sensor.temperature_low
onemind.devices.sensor.humidity_high
onemind.devices.sensor.motion_detected
onemind.devices.sensor.door_opened
onemind.devices.sensor.door_closed
```

### USB Events
```
onemind.devices.usb.hotplug
onemind.devices.usb.removed
```

### Smart Home Events
```
onemind.devices.home.light_on
onemind.devices.home.light_off
onemind.devices.home.door_locked
onemind.devices.home.door_unlocked
onemind.devices.home.window_opened
```

---

## Integration with main.py

The device system is automatically initialized in [`backend/main.py`](../main.py):

```python
# ─── 9. Start Device Event System ───
from backend.devices import get_device_registry
from backend.devices.listeners import start_all_listeners

# Initialize registry
registry = get_device_registry()

# Register devices (examples)
# maixcam = MaixCAMSource(...)
# registry.register(maixcam)

# Start devices
await registry.start_all()

# Start listeners
await start_all_listeners()
```

---

## Creating Custom Device Sources

### Step 1: Inherit from DeviceSource

```python
from backend.devices import DeviceSource, DeviceEvent
from backend.bus import publish

class TemperatureSensor(DeviceSource):
    def __init__(self, device_id: str, sensor_path: str):
        super().__init__(device_id, "sensor")
        self.sensor_path = sensor_path
        self._task = None

    async def start(self):
        """Start monitoring temperature"""
        self.running = True
        self._task = asyncio.create_task(self._monitoring_loop())

    async def stop(self):
        """Stop monitoring"""
        self.running = False
        if self._task:
            self._task.cancel()

    async def get_status(self):
        """Return sensor status"""
        return {
            "device_id": self.device_id,
            "device_type": "sensor",
            "running": self.running,
            "sensor_path": self.sensor_path
        }

    async def _monitoring_loop(self):
        """Monitor temperature and publish events"""
        while self.running:
            # Read temperature from sensor
            temp = await self._read_temperature()

            # Check threshold
            if temp > 85:  # High temperature
                event = DeviceEvent(
                    event_type="temperature_high",
                    device_id=self.device_id,
                    device_type="sensor",
                    confidence=1.0,
                    data={"temperature": temp, "unit": "F"},
                    metadata=self.metadata
                )

                await publish(event.to_nats_topic(), event.to_dict())

            await asyncio.sleep(60)  # Check every minute

    async def _read_temperature(self) -> float:
        """Read temperature from sensor"""
        # Implementation depends on sensor type
        return 75.0  # Placeholder
```

### Step 2: Register and Use

```python
from backend.devices import get_device_registry

# Create sensor
temp_sensor = TemperatureSensor(
    device_id="server_room_temp",
    sensor_path="/sys/class/hwmon/hwmon0/temp1_input"
)

temp_sensor.set_metadata(location="server_room", floor=2)

# Register
registry = get_device_registry()
registry.register(temp_sensor)

# Start
await registry.start_all()
```

---

## Testing

### Run Unit Tests

```bash
pytest tests/test_device_events.py -v
```

### Run Demo

```bash
python examples/device_events_demo.py
```

### Manual Testing

```python
import asyncio
from backend.devices.sources.camera import MockCameraSource
from backend.devices import get_device_registry

async def test():
    # Create mock camera
    camera = MockCameraSource("test")
    camera.set_mock_detections([
        {"class": "person", "confidence": 0.85}
    ])

    # Register and start
    registry = get_device_registry()
    registry.register(camera)
    await registry.start_all()

    # Wait
    await asyncio.sleep(5)

    # Stop
    await registry.stop_all()

asyncio.run(test())
```

---

## OM1 Robot Integration

The Device Event System is **the foundation** for OM1 robot integration:

```python
# backend/openmind/devices/om1_sensors.py

class OM1VisionSource(DeviceSource):
    """OM1 robot vision system"""

    async def _monitoring_loop(self):
        while self.running:
            # Get vision data from ROS2
            detections = await self.ros2_client.get_detections()

            # Publish to NATS
            for det in detections:
                event = DeviceEvent(
                    event_type="obstacle_detected",
                    device_id="om1_camera",
                    device_type="camera",
                    data=det
                )
                await publish(event.to_nats_topic(), event.to_dict())
```

**Benefits for OM1**:
- ✅ Loose coupling between sensors and agents
- ✅ Multiple agents can react to same sensor data
- ✅ Full event logging for debugging
- ✅ Easy to add new sensors without agent changes

---

## Troubleshooting

### No Events Received

1. **Check NATS is running**:
   ```bash
   docker-compose up -d
   ```

2. **Verify device is started**:
   ```python
   status = await registry.get_all_status()
   print(status)
   ```

3. **Check topic pattern**:
   ```python
   # Wildcard patterns
   "onemind.devices.camera.*"     # All camera events
   "onemind.devices.*.*"           # All device events
   ```

### Device Not Publishing

1. **Check device running**:
   ```python
   device = registry.get_device("my_device")
   print(device.running)
   ```

2. **Check for errors**:
   ```bash
   # View logs
   tail -f logs/app.log | grep "device"
   ```

3. **Verify NATS connection**:
   ```python
   from backend.bus import get_bus
   nc = await get_bus()
   print(f"Connected: {nc is not None}")
   ```

---

## Next Steps

1. **Add More Device Types**:
   - Temperature sensors
   - Motion sensors
   - Door/window sensors
   - Smart home integrations

2. **Enhance Camera Sources**:
   - RTSP stream support
   - USB camera support
   - YOLO inference integration

3. **Advanced Features**:
   - Device discovery (auto-detect cameras on network)
   - Event aggregation (combine multiple sensors)
   - Event replay (debugging)
   - Health monitoring dashboard

---

## References

- [ENHANCEMENT_ROADMAP.md](../../ENHANCEMENT_ROADMAP.md) - Full implementation guide
- [NATS-BUS.md](../../docs/architecture/NATS-BUS.md) - NATS architecture
- [Agno AgentOS](https://github.com/agno-oss/agno) - Agent framework
- [MaixCAM](https://wiki.sipeed.com/maixcam) - AI camera documentation

---

**Maintained By**: OneMindOS Core Team
**Last Updated**: 2025-02-17


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
