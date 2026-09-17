---
name: transit
description: Multi-modal transit routing via Google Maps — subway, bus, train, walking, driving, bicycling.
metadata: { "clawdbot": { "emoji": "🚇" } }
---

# Transit

Plan trips, get directions, and find nearby transit stations via Google Maps.

## Tools

- `plan_trip(origin, destination, mode?, arrive_by?, depart_at?)` — full route planning with up to 3 alternatives
- `get_directions(origin, destination, mode?)` — simple turn-by-turn directions
- `nearby_transit(location, radius?)` — find subway/bus/train stations nearby

## Modes

`transit` · `walking` · `driving` · `bicycling`

## Env

`GOOGLE_MAPS_API_KEY` — needs Directions, Geocoding, and Places APIs enabled

## Examples

- "How do I get to Grand Central by subway?"
- "Plan my commute from Greenwich to Midtown arriving by 9am"
- "Find the nearest subway stop to Times Square"
- "Walking directions from Penn Station to Madison Square Garden"


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
