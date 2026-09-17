---
name: simple-weather
description: >
  Gets current weather conditions and forecasts using the Open-Meteo free
  weather API (no API key required). Use when the user asks about weather,
  temperature, forecast, rain, conditions, or whether to bring an umbrella.
  Supports any location by name or coordinates.
metadata:
  author: zeus
  version: "1.0.0"
---

# Simple Weather Skill

## When to Use

Activate when the user:

- Asks "what's the weather" anywhere
- Wants a forecast ("will it rain tomorrow?")
- Asks about temperature, humidity, wind speed, or UV index
- Says "should I bring an umbrella?" or "is it hot outside?"

## Available Tools

| Tool               | Description                      |
| ------------------ | -------------------------------- |
| `weather_current`  | Current conditions at a location |
| `weather_forecast` | 7-day forecast for a location    |

## Quick Reference

### Current weather at a named city

```
weather_current({ location: "Lagos, Nigeria" })
```

### Current weather by coordinates

```
weather_current({ latitude: 6.5244, longitude: 3.3792 })
```

### 7-day forecast

```
weather_forecast({ location: "New York" })
```

### Tomorrow's forecast only

```
weather_forecast({ location: "London", days: 2 })
```

## Response Fields

`weather_current` returns:

- `temperature_c` — Current temp in Celsius
- `temperature_f` — Current temp in Fahrenheit
- `feels_like_c` / `feels_like_f` — Apparent temperature
- `humidity_pct` — Relative humidity (%)
- `wind_speed_kmh` — Wind speed
- `wind_direction` — Wind direction (e.g., "NE")
- `condition` — Human-readable condition (e.g., "Partly cloudy")
- `uv_index` — UV index (0-11+)
- `precipitation_mm` — Precipitation in last hour (mm)

## Notes

- This skill uses Open-Meteo (free, no API key needed)
- Location names are geocoded using the Open-Meteo geocoding API
- If the location is not found, suggest the user provide coordinates
- Temperature is always returned in both C and F for flexibility
- Forecasts show conditions at noon each day (good general representation)
