---
name: farm-operations
description: Farm management including crops, livestock, equipment, and irrigation
license: Apache-2.0
metadata:
  version: "1.0.0"
  author: zeus
  tags: ["farming", "agriculture", "livestock", "homestead"]
  pillar: "LE"
---

# Farm Operations Skill

Use this skill for managing all aspects of farm and homestead operations.

## When to Use

- Planning crop rotations
- Scheduling farm tasks
- Managing livestock
- Equipment maintenance
- Irrigation planning
- Seasonal preparations
- Weather response

## Farm Calendar

### Spring (Mar-May)
- Soil preparation
- Seed starting
- Transplanting
- Equipment maintenance
- Irrigation system check

### Summer (Jun-Aug)
- Irrigation management
- Pest monitoring
- Harvest planning
- Heat stress management
- Equipment repairs

### Fall (Sep-Nov)
- Main harvest
- Soil amendments
- Cover crop planting
- Equipment winterization
- Storage preparation

### Winter (Dec-Feb)
- Planning next season
- Equipment repairs
- Seed ordering
- Infrastructure projects
- Record review

## Crop Management

### Planning
1. Review previous year's data
2. Plan rotations (3-4 year cycle)
3. Calculate seed quantities
4. Schedule planting dates
5. Plan succession planting

### Monitoring
- Daily: Visual inspection
- Weekly: Pest/disease check
- Monthly: Growth measurements
- Seasonal: Yield tracking

## Livestock Management

### Daily Tasks
- Morning: Feed, water, health check
- Evening: Feed, secure enclosures
- Record: Any abnormalities

### Weekly Tasks
- Enclosure cleaning
- Equipment check
- Supply inventory

### Health Indicators
- Appetite and water intake
- Movement and behavior
- Coat/feather condition
- Weight trends

## Equipment

### Critical Equipment
| Equipment | Check Frequency | Season |
|-----------|-----------------|--------|
| Tractor | Weekly | All |
| Irrigation | Daily (summer) | Growing |
| Mower | Bi-weekly | Growing |
| Chainsaw | Monthly | All |

### Maintenance Schedule
- Oil changes: Per hour meters
- Filter replacement: Seasonal
- Belt inspection: Monthly
- Full service: Annual (winter)

## Irrigation

### System Types
- Drip: Vegetables, orchards
- Sprinkler: Pasture, lawns
- Flood: Certain crops

### Scheduling Factors
- Soil moisture levels
- Weather forecast
- Crop stage
- Water availability

### Automation via Home Assistant
```yaml
# Example automation
automation:
  trigger:
    - platform: numeric_state
      entity_id: sensor.soil_moisture_garden
      below: 30
  action:
    - service: switch.turn_on
      entity_id: switch.garden_irrigation
```

## Decision Framework

### Priority Matrix
| Urgency | Importance | Action |
|---------|------------|--------|
| High | High | Immediate |
| High | Low | Delegate or quick fix |
| Low | High | Schedule |
| Low | Low | Batch or eliminate |

### Weather Response
- Frost warning: Cover crops, bring in livestock
- Heat wave: Increase irrigation, provide shade
- Heavy rain: Check drainage, delay fieldwork
- Wind: Secure structures, delay spraying

## Integration Points

### Home Assistant Entities
- `sensor.soil_moisture_*`
- `sensor.weather_*`
- `switch.irrigation_*`
- `binary_sensor.rain_*`

### Todoist Projects
- Farm Tasks (daily/weekly)
- Equipment Maintenance
- Seasonal Planning

### Linear Issues
- OMOS-LE-* for farm projects

## Best Practices

1. **Record everything** - Data drives decisions
2. **Plan ahead** - Seasonal prep is key
3. **Maintain equipment** - Prevention > repair
4. **Watch weather** - Adapt to conditions
5. **Rotate crops** - Soil health matters
6. **Automate where possible** - Sensors save time


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
