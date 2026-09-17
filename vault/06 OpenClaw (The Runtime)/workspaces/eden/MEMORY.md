---
type: memory-architecture
agent: eden
status: active
updated: 2026-05-03
---
# MEMORY.md — Eden 🌳

## Memory Architecture

Eden maintains homestead memory — land state, agricultural cycles, environmental conditions, and sustainability metrics.

## Tiers

### Tier 1: Working Memory (session-scoped)
- Current growing conditions (season, weather, soil state)
- Active planting/harvesting tasks
- Pending land maintenance items
- Sensor alerts from garden/farm systems

### Tier 2: Episodic Memory (persistent)
- Planting log (what was planted, when, where, outcome)
- Harvest records (yield, quality, timing)
- Land improvement history (soil amendments, infrastructure additions)
- Weather event impacts and recovery actions
- Seasonal patterns and their effects on production

### Tier 3: Semantic Memory (knowledge base)
- Property layout and zone map
- Soil composition data per area
- Plant compatibility and companion planting rules
- Local climate patterns and frost dates
- Permaculture design principles applied to the property
- Water systems (irrigation, rainwater, wells)

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Current conditions | Sensors / workspace memory | Real-time |
| Growing log | Codex vault LE 50 Home | Permanent |
| Land map | Codex vault LE 50 Home | Permanent |
| Weather history | Local weather DB | 5 years |

## Memory Rules

1. **Seasonal awareness** — Always know what season it is and what that means for the land
2. **Track yield data** — Every harvest is data for next year's planning
3. **Long-term soil health** — Amendments, rotation, cover crops — track it all
4. **Weather correlation** — Connect weather events to outcomes (frost damage, drought stress)
5. **Record what works** — Varieties, timing, techniques that succeed in THIS specific climate/soil

## What Eden Remembers

- Complete planting and harvest history
- Soil test results and amendment history
- Property infrastructure (beds, paths, irrigation, structures)
- Which varieties perform best in the local microclimate
- Water usage and efficiency data
- Pest/disease history and treatments applied
- Seasonal task calendars (when to start seeds, transplant, mulch)
- Long-term land improvement plan and progress
