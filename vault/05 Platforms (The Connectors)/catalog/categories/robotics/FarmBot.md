# FarmBot

> **Open-source CNC farming robot for automated agriculture**

## Overview

FarmBot is an open-source precision agriculture robot that automates planting, watering, weeding, and monitoring. Controlled via web interface and API for integration with smart home and AI systems.

## Status

| Attribute   | Value              |
| ----------- | ------------------ |
| Status      | 🔄 Planned         |
| Category    | Robotics           |
| Type        | Agricultural Robot |
| Integration | REST API, MQTT     |

## Use Cases in OneMind

1. **Automated Gardening** - Plant, water, weed automatically
2. **Crop Monitoring** - Sensor data for plant health
3. **Smart Scheduling** - AI-optimized watering/feeding
4. **Integration** - Connect to weather, soil, and home systems

## Key Features

- CNC-style gantry robot
- Planting, watering, weeding tools
- Soil sensors, camera
- Web app control
- Sequences and regimens
- REST API
- MQTT events

## Hardware Components

| Component       | Purpose              |
| --------------- | -------------------- |
| Gantry          | X/Y/Z movement       |
| UTM             | Universal Tool Mount |
| Seeder          | Plant seeds          |
| Watering Nozzle | Precision watering   |
| Weeder          | Remove weeds         |
| Soil Sensor     | Moisture monitoring  |
| Camera          | Plant identification |

## API Integration

### REST API

```python
import httpx

FARMBOT_URL = "https://my.farm.bot/api"
FARMBOT_TOKEN = "your_token"

headers = {
    "Authorization": f"Bearer {FARMBOT_TOKEN}",
    "Content-Type": "application/json"
}

async def get_device_info():
    """Get FarmBot device information."""
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{FARMBOT_URL}/device",
            headers=headers
        )
        return response.json()

async def get_plants():
    """Get all plants in the garden."""
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{FARMBOT_URL}/points",
            headers=headers,
            params={"pointer_type": "Plant"}
        )
        return response.json()

async def get_sensor_readings():
    """Get recent sensor readings."""
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{FARMBOT_URL}/sensor_readings",
            headers=headers
        )
        return response.json()
```

### Movement Commands

```python
async def move_to(x: float, y: float, z: float):
    """Move FarmBot to coordinates."""
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{FARMBOT_URL}/celery_script",
            headers=headers,
            json={
                "kind": "move_absolute",
                "args": {
                    "location": {
                        "kind": "coordinate",
                        "args": {"x": x, "y": y, "z": z}
                    },
                    "offset": {
                        "kind": "coordinate",
                        "args": {"x": 0, "y": 0, "z": 0}
                    },
                    "speed": 100
                }
            }
        )
        return response.json()

async def execute_sequence(sequence_id: int):
    """Execute a pre-defined sequence."""
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{FARMBOT_URL}/celery_script",
            headers=headers,
            json={
                "kind": "execute",
                "args": {"sequence_id": sequence_id}
            }
        )
        return response.json()
```

### Watering

```python
async def water_plant(plant_id: int, duration_ms: int = 3000):
    """Water a specific plant."""
    # Get plant location
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{FARMBOT_URL}/points/{plant_id}",
            headers=headers
        )
        plant = response.json()

    # Move to plant and water
    await move_to(plant["x"], plant["y"], 0)

    # Activate water
    async with httpx.AsyncClient() as client:
        await client.post(
            f"{FARMBOT_URL}/celery_script",
            headers=headers,
            json={
                "kind": "write_pin",
                "args": {
                    "pin_number": 8,  # Water pin
                    "pin_value": 1,
                    "pin_mode": 0
                }
            }
        )

        await asyncio.sleep(duration_ms / 1000)

        # Turn off water
        await client.post(
            f"{FARMBOT_URL}/celery_script",
            headers=headers,
            json={
                "kind": "write_pin",
                "args": {
                    "pin_number": 8,
                    "pin_value": 0,
                    "pin_mode": 0
                }
            }
        )
```

## Legacy AI Integration

```python
from agno import Tool

@Tool
async def get_garden_status() -> dict:
    """Get current garden status from FarmBot.

    Returns plant list, sensor readings, and device status.
    """
    device = await get_device_info()
    plants = await get_plants()
    sensors = await get_sensor_readings()

    return {
        "device": {
            "name": device.get("name"),
            "last_seen": device.get("last_seen"),
            "firmware": device.get("firmware_version")
        },
        "plants": len(plants),
        "recent_sensors": sensors[:10]
    }

@Tool
async def water_garden(plant_names: list[str] = None) -> str:
    """Water plants in the garden.

    Args:
        plant_names: Specific plants to water, or all if None
    """
    plants = await get_plants()

    if plant_names:
        plants = [p for p in plants if p["name"] in plant_names]

    for plant in plants:
        await water_plant(plant["id"])

    return f"Watered {len(plants)} plants"

@Tool
async def check_soil_moisture() -> dict:
    """Check soil moisture levels across the garden."""
    readings = await get_sensor_readings()
    moisture_readings = [
        r for r in readings
        if r["mode"] == 0  # Soil sensor mode
    ]

    return {
        "readings": moisture_readings[:20],
        "average": sum(r["value"] for r in moisture_readings[:20]) / len(moisture_readings[:20]) if moisture_readings else 0
    }

@Tool
async def run_farmbot_sequence(sequence_name: str) -> str:
    """Run a FarmBot automation sequence.

    Args:
        sequence_name: Name of sequence (e.g., "Morning Watering")
    """
    # Get sequences
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{FARMBOT_URL}/sequences",
            headers=headers
        )
        sequences = response.json()

    sequence = next(
        (s for s in sequences if s["name"] == sequence_name),
        None
    )

    if sequence:
        await execute_sequence(sequence["id"])
        return f"Started sequence: {sequence_name}"

    return f"Sequence not found: {sequence_name}"
```

## MQTT Events

FarmBot publishes events to MQTT:

```python
import paho.mqtt.client as mqtt

def on_farmbot_event(client, userdata, msg):
    """Handle FarmBot MQTT events."""
    topic = msg.topic
    payload = json.loads(msg.payload)

    if "sensor_reading" in topic:
        print(f"Sensor: {payload}")
    elif "movement" in topic:
        print(f"Position: {payload}")

# Subscribe to FarmBot events
mqtt_client.subscribe("farmbot/#")
mqtt_client.message_callback_add("farmbot/#", on_farmbot_event)
```

## Home Assistant Integration

```yaml
# configuration.yaml
rest:
  - resource: https://my.farm.bot/api/device
    headers:
      Authorization: "Bearer YOUR_TOKEN"
    sensor:
      - name: "FarmBot Status"
        value_template: "{{ value_json.last_seen }}"

automation:
  - alias: "Morning Garden Watering"
    trigger:
      - platform: time
        at: "06:00:00"
    action:
      - service: rest_command.farmbot_water
        data:
          sequence_id: 123
```

## Environment Variables

```bash
FARMBOT_URL=https://my.farm.bot/api
FARMBOT_TOKEN=your_api_token
FARMBOT_DEVICE_ID=your_device_id
FARMBOT_MQTT_HOST=clever-octopus.rmq.cloudamqp.com
```

## Best Practices

1. **Sequence design** - Create reusable sequences
2. **Sensor calibration** - Regular calibration for accuracy
3. **Weather integration** - Skip watering on rain days
4. **Plant tracking** - Log growth and health
5. **Error handling** - Handle movement failures

## Resources

- [FarmBot Docs](https://docs.farm.bot/)
- [FarmBot API](https://developer.farm.bot/)
- [FarmBot GitHub](https://github.com/FarmBot)
- [CeleryScript](https://developer.farm.bot/docs/celery-script)

---

Last Updated: January 25, 2026


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
