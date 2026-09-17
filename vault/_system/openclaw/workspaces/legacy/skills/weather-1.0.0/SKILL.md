---
name: weather
description: Get current weather and forecasts for any location (no API key required).
homepage: https://wttr.in/:help
metadata: { "clawdbot": { "emoji": "🌤️" } }
---

# Weather

Current conditions and multi-day forecasts. No API keys needed.

## Tools

- `get_weather(location)` — current conditions (temp, humidity, wind, UV, cloud cover, precipitation)
- `get_forecast(location, days?)` — multi-day forecast with highs/lows, precipitation, sunrise/sunset

## Data Sources

- **wttr.in** — quick human-readable summary (primary)
- **Open-Meteo** — structured JSON data with hourly/daily breakdowns (fallback + detailed)
- **Open-Meteo Geocoding** — free location lookup (no Google API key needed)

## Examples

- "What's the weather in Greenwich, CT?"
- "Will it rain this weekend in New York?"
- "Give me a 7-day forecast for Miami"
- "What's the UV index in Phoenix?"
