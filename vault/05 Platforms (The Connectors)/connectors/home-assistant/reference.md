# Home Assistant API Reference

## Base URL
```
$HOME_ASSISTANT_URL/api/
```

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/` | API status |
| GET | `/api/states` | All entity states |
| GET | `/api/states/{entity_id}` | Single entity state |
| POST | `/api/services/{domain}/{service}` | Call service |
| POST | `/api/events/{event_type}` | Fire event |

## Common Services

| Domain | Service | Payload |
|--------|---------|---------|
| `light` | `turn_on` | `{"entity_id": "...", "brightness": 255}` |
| `light` | `turn_off` | `{"entity_id": "..."}` |
| `switch` | `toggle` | `{"entity_id": "..."}` |
| `climate` | `set_temperature` | `{"entity_id": "...", "temperature": 72}` |
| `automation` | `trigger` | `{"entity_id": "..."}` |
| `script` | `turn_on` | `{"entity_id": "..."}` |
| `scene` | `turn_on` | `{"entity_id": "..."}` |
| `tts` | `speak` | `{"entity_id": "tts.xxx", "media_player_entity_id": "...", "message": "..."}` |

## One Mind Locations

| Prefix | Location |
|--------|----------|
| `grw_` | Greenwich Home |
| `nyc_` | NYC Office |
| `mex_` | Mexico Farm |
| `ups_` | Upstate Homestead |

## Entity Domains

`light`, `switch`, `sensor`, `binary_sensor`, `climate`, `cover`, `media_player`, `automation`, `script`, `scene`, `camera`, `lock`


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
