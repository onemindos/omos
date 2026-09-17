---
name: home-assistant
description: Control Home Assistant smart home devices, automations, scripts, and scenes. Use when you need to turn on/off lights, check sensors, trigger automations, or get home status. Supports multiple locations.
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Read
  - Write
---

# Home Assistant Skill

Control smart home devices via Home Assistant REST API.

## Environment Variables

```bash
HOME_ASSISTANT_URL     # e.g., https://your-ha-instance.ui.nabu.casa or http://homeassistant.local:8123
HOME_ASSISTANT_TOKEN   # Long-lived access token from HA
```

## Quick Reference

| Operation | Command |
|-----------|---------|
| Get all states | `curl -s "$HOME_ASSISTANT_URL/api/states" -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN"` |
| Turn on light | `curl -s -X POST "$HOME_ASSISTANT_URL/api/services/light/turn_on" -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" -H "Content-Type: application/json" -d '{"entity_id": "light.living_room"}'` |

## One Mind Locations

Multi-location support with entity prefixes:

| Prefix | Location | Role |
|--------|----------|------|
| `grw_` | Greenwich Home | Primary |
| `nyc_` | NYC Office | Work |
| `mex_` | Mexico Farm | Vacation |
| `ups_` | Upstate Homestead | Compound |

Example: `light.grw_office_desk`, `switch.nyc_desk_monitor`

## Operations

### Get All States
```bash
curl -s "$HOME_ASSISTANT_URL/api/states" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" | jq '.[].entity_id'
```

### Get Entity State
```bash
curl -s "$HOME_ASSISTANT_URL/api/states/light.grw_office_desk" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" | jq '{state, attributes}'
```

### Turn On
```bash
# Basic turn on
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/light/turn_on" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "light.grw_office_desk"}'

# With brightness
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/light/turn_on" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "light.grw_office_desk", "brightness": 128}'
```

### Turn Off
```bash
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/light/turn_off" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "light.grw_office_desk"}'
```

### Toggle
```bash
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/light/toggle" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "light.grw_office_desk"}'
```

### Trigger Automation
```bash
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/automation/trigger" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "automation.grw_morning_routine"}'
```

### Run Script
```bash
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/script/turn_on" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "script.grw_movie_mode"}'
```

### Activate Scene
```bash
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/scene/turn_on" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "scene.grw_dinner_time"}'
```

### Text-to-Speech
```bash
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/tts/speak" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "entity_id": "tts.home_assistant_cloud",
    "media_player_entity_id": "media_player.grw_speaker_echo",
    "message": "Dinner is ready"
  }'
```

### Climate Control
```bash
# Set temperature
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/climate/set_temperature" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "climate.grw_thermostat", "temperature": 72}'

# Set HVAC mode
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/climate/set_hvac_mode" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "climate.grw_thermostat", "hvac_mode": "heat"}'
```

### Call Any Service
```bash
curl -s -X POST "$HOME_ASSISTANT_URL/api/services/{domain}/{service}" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "domain.entity_name", ...service_data}'
```

## Common Domains

| Domain | Examples |
|--------|----------|
| `light` | turn_on, turn_off, toggle |
| `switch` | turn_on, turn_off, toggle |
| `climate` | set_temperature, set_hvac_mode |
| `cover` | open_cover, close_cover, set_position |
| `media_player` | play, pause, volume_set |
| `automation` | trigger, turn_on, turn_off |
| `script` | turn_on (runs script) |
| `scene` | turn_on (activates scene) |
| `tts` | speak |
| `notify` | notify (send notification) |

## Filter States by Domain
```bash
# Get all lights
curl -s "$HOME_ASSISTANT_URL/api/states" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" | \
  jq '[.[] | select(.entity_id | startswith("light."))]'

# Get all sensors for a location
curl -s "$HOME_ASSISTANT_URL/api/states" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" | \
  jq '[.[] | select(.entity_id | contains("grw_")) | select(.entity_id | startswith("sensor."))]'
```

## Tips

- Entity IDs follow pattern: `domain.location_name` (e.g., `light.grw_office_desk`)
- Check entity state before acting: `states/{entity_id}`
- Use `jq` to parse and filter large state lists
- Long-lived tokens: Settings → Security → Long-Lived Access Tokens


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
