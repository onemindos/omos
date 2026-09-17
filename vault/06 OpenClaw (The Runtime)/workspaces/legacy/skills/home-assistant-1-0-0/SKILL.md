# Home Assistant Skill

General-purpose Home Assistant control via REST API. Works with any HA instance.

## Environment Variables

| Var                    | Description                                          |
| ---------------------- | ---------------------------------------------------- |
| `HOME_ASSISTANT_URL`   | HA base URL (e.g. `http://homeassistant.local:8123`) |
| `HOME_ASSISTANT_TOKEN` | Long-lived access token from HA                      |

## Tools

| Tool            | Description                                                                 |
| --------------- | --------------------------------------------------------------------------- |
| `ha_control`    | Turn on/off, toggle, lock/unlock devices. Set brightness/temperature.       |
| `ha_status`     | Home summary: lights on/off, locks, doors, climate. Optional domain filter. |
| `ha_state`      | Get state + attributes of a specific entity.                                |
| `ha_scene`      | Activate a scene (e.g. `good_night`, `movie_time`).                         |
| `ha_automation` | Trigger an automation.                                                      |
| `ha_list`       | List all entities for a domain (light, lock, sensor, etc.).                 |

## Examples

```
ha_control(action="turn_on", entity_id="light.office_desk", brightness=128)
ha_control(action="lock", entity_id="lock.front_door")
ha_status()
ha_status(domain="light")
ha_state(entity_id="sensor.outdoor_temperature")
ha_scene(scene="good_night")
ha_automation(automation_id="morning_routine")
ha_list(domain="light")
```

## vs homestead-control

`homestead-control` is NATS-based, specific to the 6-building homestead edge grid.
`home-assistant` is REST-based, works with any HA instance for general device control.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
