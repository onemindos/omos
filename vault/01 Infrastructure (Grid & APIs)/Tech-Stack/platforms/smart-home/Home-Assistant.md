# Home Assistant

> **Smart Home Automation Platform**

## Overview

Home Assistant is an open-source home automation platform that puts local control and privacy first. It powers Legacy AI's smart home capabilities.

## Status

- **Active**: ✅ Yes
- **Deployment**: Local (Raspberry Pi / Home Server)
- **Last Updated**: January 25, 2026

## Use Cases

1. **Home Control** - Lights, switches, thermostats
2. **Presence Detection** - Who's home
3. **Automations** - Time/event-based actions
4. **Voice Control** - Via Legacy AI
5. **Monitoring** - Sensors, cameras

## Key Features

- 2000+ integrations
- Local-first (works offline)
- Privacy-focused (no cloud required)
- Powerful automation engine
- Dashboards (Lovelace)
- Add-on ecosystem

## Architecture

```
Devices (Zigbee, Z-Wave, WiFi)
    ↓
Home Assistant Core
    ↓
REST API / WebSocket
    ↓
Legacy AI Agent
    ↓
Voice/Chat Commands
```

## Integration Types

### Supported Protocols

| Protocol      | Hub Required          | Common Devices              |
| ------------- | --------------------- | --------------------------- |
| **Zigbee**    | Yes (ZHA/Zigbee2MQTT) | Lights, sensors, switches   |
| **Z-Wave**    | Yes                   | Locks, sensors, switches    |
| **WiFi**      | No                    | Smart plugs, bulbs, cameras |
| **Bluetooth** | No                    | Sensors, trackers           |
| **Matter**    | No                    | New standard devices        |
| **Thread**    | Border Router         | Low-power devices           |

## Configuration

### configuration.yaml

```yaml
homeassistant:
  name: OneMind Home
  unit_system: imperial
  time_zone: America/New_York
  latitude: !secret latitude
  longitude: !secret longitude

# Enable API
api:

# Enable frontend
frontend:

# Recorder for history
recorder:
  db_url: !secret db_url
  purge_keep_days: 30

# REST API for Legacy AI
rest_command:
  legacy_notify:
    url: http://legacy-ai:8080/webhooks/home
    method: POST
    payload: '{"event": "{{ event }}", "data": {{ data }}}'
```

## API Integration

### REST API Endpoints

```python
# Get all states
GET http://homeassistant.local:8123/api/states
Authorization: Bearer <long_lived_token>

# Get specific entity
GET http://homeassistant.local:8123/api/states/light.living_room

# Call service
POST http://homeassistant.local:8123/api/services/light/turn_on
{
  "entity_id": "light.living_room",
  "brightness": 255
}
```

### Python Integration

```python
from homeassistant_api import Client

client = Client(
    'http://homeassistant.local:8123/api',
    os.getenv('HOME_ASSISTANT_TOKEN')
)

# Get entity state
light = client.get_entity(entity_id='light.living_room')
print(light.state.state)  # 'on' or 'off'

# Turn on light
client.trigger_service(
    'light',
    'turn_on',
    entity_id='light.living_room',
    brightness=200
)
```

## Entity Categories

### Lighting

```yaml
light.living_room:
  - turn_on
  - turn_off
  - toggle
  - brightness (0-255)
  - color_temp
  - rgb_color

light.bedroom: ...
```

### Climate

```yaml
climate.thermostat:
  - set_temperature
  - set_hvac_mode (heat, cool, auto, off)
  - current_temperature
  - target_temperature
```

### Sensors

```yaml
sensor.outdoor_temperature:
sensor.indoor_humidity:
sensor.motion_living_room:
binary_sensor.front_door:
```

## Automations

### Example: Welcome Home

```yaml
automation:
  - alias: "Welcome Home"
    trigger:
      - platform: state
        entity_id: person.zeus
        to: "home"
    condition:
      - condition: sun
        after: sunset
    action:
      - service: light.turn_on
        target:
          entity_id: light.entryway
      - service: climate.set_temperature
        target:
          entity_id: climate.thermostat
        data:
          temperature: 72
      - service: notify.legacy_ai
        data:
          message: "Zeus is home"
```

### Example: Good Night

```yaml
automation:
  - alias: "Good Night"
    trigger:
      - platform: event
        event_type: legacy_command
        event_data:
          command: "good_night"
    action:
      - service: light.turn_off
        target:
          area_id: all
      - service: lock.lock
        target:
          entity_id: lock.front_door
      - service: climate.set_temperature
        data:
          temperature: 68
```

## Legacy AI Integration

### Agent Tools

```python
@Tool
def control_light(entity_id: str, action: str, brightness: int = None) -> str:
    """Control a light in the home."""
    if action == "on":
        client.trigger_service('light', 'turn_on', entity_id=entity_id, brightness=brightness)
    elif action == "off":
        client.trigger_service('light', 'turn_off', entity_id=entity_id)
    return f"Light {entity_id} turned {action}"

@Tool
def get_home_status() -> dict:
    """Get current status of the home."""
    return {
        "temperature": get_sensor("sensor.indoor_temperature"),
        "lights_on": count_lights_on(),
        "doors_locked": check_locks(),
        "people_home": get_people_home(),
    }

@Tool
def set_thermostat(temperature: int, mode: str = None) -> str:
    """Set the thermostat temperature and mode."""
    ...
```

### Voice Commands

| Command                     | Action                     |
| --------------------------- | -------------------------- |
| "Turn off the lights"       | All lights off             |
| "Set temperature to 72"     | Thermostat → 72°F          |
| "Is the front door locked?" | Check lock status          |
| "Who's home?"               | Check presence             |
| "Good night"                | Trigger bedtime automation |

## Dashboard (Lovelace)

### Custom Cards

- **Mini Graph Card**: Sensor history
- **Mushroom Cards**: Modern UI
- **Button Card**: Custom buttons
- **Layout Card**: Grid layouts

### Example Dashboard

```yaml
views:
  - title: Home
    cards:
      - type: entities
        title: Living Room
        entities:
          - light.living_room
          - switch.fan
          - sensor.temperature

      - type: thermostat
        entity: climate.main

      - type: weather-forecast
        entity: weather.home
```

## Add-ons

### Essential Add-ons

| Add-on             | Purpose                  |
| ------------------ | ------------------------ |
| **Zigbee2MQTT**    | Zigbee device management |
| **Node-RED**       | Visual automation        |
| **ESPHome**        | Custom sensors           |
| **Mosquitto**      | MQTT broker              |
| **File Editor**    | Edit configs in UI       |
| **Terminal & SSH** | Command line access      |

## Security

- ✅ Long-lived access tokens (stored in 1Password)
- ✅ HTTPS enabled
- ✅ Tailscale VPN for remote access
- ✅ Two-factor authentication
- ✅ Regular updates

## Best Practices

- ✅ Use areas to organize entities
- ✅ Name entities clearly (light.living_room_main)
- ✅ Group related automations
- ✅ Use scripts for reusable actions
- ✅ Back up configuration regularly
- ✅ Test automations before relying on them
- ✅ Monitor system resources

## Resources

- [Home Assistant Documentation](https://www.home-assistant.io/docs/)
- [Integrations](https://www.home-assistant.io/integrations/)
- [Community](https://community.home-assistant.io/)

---

**Related**: [NATS](./NATS.md) | [Agno](./Agno.md)
