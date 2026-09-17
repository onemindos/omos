---
id: le-50-20260113225956
title: "Automated Farming System Research"
type: reference
path: LE
domain: 50-74
folder: "50"
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [path/le, path/le/50, type/reference, status/active]
ai_summary: "Automated Farming System Research"
---
# Automated Farming System Research

> *Research compiled by Legacy AI — 2026-01-10*
> *Status: Future Implementation*

---

## Primary System: FarmBot Genesis

**Website:** https://farm.bot
**Documentation:** https://software.farm.bot/docs
**Developer API:** https://developer.farm.bot
**GitHub:** https://github.com/FarmBot

### What is FarmBot?

FarmBot is humanity's first open-source CNC farming machine designed for at-home automated food production. It's a Cartesian robot (X/Y/Z axis movement) that automates precision gardening tasks.

### Current Version: v1.8 (Now Shipping)

- **90% pre-assembled** in box
- **~4 hour setup** time
- Simplified track design (5x faster assembly than v1.6)

### Models & Pricing

| Model | Bed Size | Track Length | Current Price | Regular Price |
|-------|----------|--------------|---------------|---------------|
| Genesis v1.7 | 1.5m x 3m (4.5m²) | 3m | ~$2,995 | $3,295 |
| Genesis XL v1.7 | 3m x 6m (18m²) | 6m | **$3,995** | $4,395 |
| Genesis v1.8 | 1.5m x 3m | 3m | Pre-order | TBD |

*Free global shipping included. US delivery: 2-6 business days.*

### Core Capabilities

| Function | Description |
|----------|-------------|
| **Planting** | Precision seed injection at exact coordinates |
| **Watering** | Per-plant watering with adjustable flow |
| **Weeding** | Autonomous weed detection and removal |
| **Soil Sensing** | Moisture and temperature monitoring |
| **Camera** | Plant health monitoring, weed detection, growth tracking |
| **Sequences** | Programmable action chains |

### Hardware Components

- **Tracks** — Aluminum rails (3m Genesis / 6m XL)
- **Gantry** — Moves along X-axis
- **Cross-slide** — Moves along Y-axis
- **Z-axis** — Vertical movement for tools
- **Universal Tool Mount** — Swappable tool heads
- **Farmduino** — Electronics control board
- **Raspberry Pi** — Runs FarmBot OS

### Tools Included

1. **Seed Injector** — Vacuum-based seed pickup and planting
2. **Watering Nozzle** — Precision water delivery
3. **Weeder** — Mechanical weed removal
4. **Soil Sensor** — Moisture/temperature probe
5. **Camera** — Computer vision for plant ID

---

## FarmBot REST API

**Base URL:** `https://my.farm.bot/api/`
**Auth:** Bearer token (JWT)
**Rate Limits:** Standard web API limits

### Key Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/plants` | GET/POST/PUT/DELETE | Manage plant records |
| `/sequences` | GET/POST/PUT/DELETE | Create/edit automation sequences |
| `/points` | GET/POST | Garden coordinates and points |
| `/tools` | GET/POST | Tool management |
| `/device` | GET/PUT | Device settings and status |
| `/farm_events` | GET/POST | Scheduled events |
| `/pin_bindings` | GET/POST | I/O pin to sequence bindings |
| `/peripherals` | GET/POST | Hardware output control |

### API Capabilities for Legacy AI

```
Legacy AI can:
├── GET /plants → Retrieve all plant data
├── POST /sequences → Create watering/care sequences
├── POST /farm_events → Schedule automated tasks
├── GET /device → Check FarmBot status
├── POST /celery_script → Execute commands directly
└── GET /sensor_readings → Pull soil moisture data
```

### Example: Create Watering Sequence via API

```json
POST /api/sequences
{
  "name": "Water Tomatoes",
  "body": [
    {"kind": "move_absolute", "args": {"location": {"x": 100, "y": 200, "z": 0}}},
    {"kind": "write_pin", "args": {"pin_number": 8, "pin_value": 1}},
    {"kind": "wait", "args": {"milliseconds": 3000}},
    {"kind": "write_pin", "args": {"pin_number": 8, "pin_value": 0}}
  ]
}
```

---

## Complementary System: Rachio (Smart Irrigation)

**Website:** https://rachio.com
**Home Assistant:** Native integration
**API Docs:** https://support.rachio.com/en_us/public-api-documentation

### Overview

Rachio is a WiFi-enabled smart sprinkler controller for lawn and garden irrigation beyond FarmBot's bed coverage.

### Models

| Model | Zones | Price |
|-------|-------|-------|
| Rachio 3 | 8 zones | ~$230 |
| Rachio 3 | 16 zones | ~$280 |
| Rachio Smart Hose Timer | 1 zone | ~$100 |

### Key Features

- **Weather Intelligence** — Auto-adjusts based on forecast
- **EPA WaterSense Certified** — 30-50% water savings
- **Zone Control** — Individual zone scheduling
- **Seasonal Shift** — Automatic seasonal adjustments
- **Flow Monitoring** — Detect leaks (with flow meter add-on)

### Home Assistant Integration

```yaml
# Auto-discovered or manual setup
# Settings → Devices & Services → Add Integration → Rachio
# Requires API key from Rachio app
```

**Entities created:**
- Binary sensors (system status)
- Switches (zone control, schedules)
- Calendar (irrigation schedule)

### API Details

- **Rate Limit:** 1,700 calls/day (~1/minute)
- **Webhooks:** Supported for real-time updates
- **Auth:** API key from app settings

---

## Complementary System: Ecowitt (Weather & Soil Sensors)

**Website:** https://www.ecowitt.com
**Home Assistant:** Native integration (HTTP callback)

### Overview

Ecowitt provides affordable weather stations and soil sensors that integrate with Home Assistant for environmental monitoring.

### Recommended Products

| Product | Function | Price |
|---------|----------|-------|
| **WH51 Soil Moisture Sensor** | Soil moisture monitoring | ~$16 |
| **GW1100 Gateway** | Connects sensors to WiFi | ~$23 |
| **GW2000 Gateway** | Advanced gateway + display | ~$70 |
| **WH65 Weather Station** | Full weather monitoring | ~$80 |
| **WS90 Weather Station** | Premium 7-in-1 sensor | ~$180 |

### WH51 Soil Moisture Sensor (Recommended)

- Measures soil moisture 0-100%
- 2-year battery life (AA)
- 433MHz wireless to gateway
- Up to 8 sensors per gateway
- Perfect for garden bed monitoring

### Home Assistant Integration

```yaml
# Ecowitt integration uses HTTP callback
# Settings → Devices & Services → Add Integration → Ecowitt
# Configure callback URL in Ecowitt app/console
# Note: Does NOT support HTTPS - must use HTTP
```

**Data available:**
- Temperature (indoor/outdoor)
- Humidity
- Soil moisture (per sensor)
- Rainfall
- Wind speed/direction
- UV index
- Barometric pressure

---

## Greenhouse Automation (Optional Expansion)

### DIY Smart Greenhouse Control

For greenhouse environments, a DIY approach using Home Assistant:

| Component | Purpose | Est. Cost |
|-----------|---------|-----------|
| **Raspberry Pi 4** | Controller | $55 |
| **Relay Board (8ch)** | Control fans/heaters/irrigation | $15 |
| **DHT22 Sensors** | Temp/humidity monitoring | $10 |
| **Soil Moisture Probes** | Multiple zones | $20 |
| **Smart Plugs** | Grow lights, heaters | $25 |
| **Exhaust Fans** | Ventilation | $50 |

### Climate Control Automation Logic

```yaml
# Example Home Assistant automation
automation:
  - alias: "Greenhouse Vent Fan"
    trigger:
      - platform: numeric_state
        entity_id: sensor.greenhouse_temperature
        above: 85
    action:
      - service: switch.turn_on
        entity_id: switch.greenhouse_vent_fan
```

### Commercial Options

- **Priva** — Enterprise greenhouse automation
- **Growlink** — Mid-tier greenhouse controllers ($500-2000)
- **Inkbird Controllers** — Budget temp/humidity control ($30-50)

---

## Full Stack Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      LEGACY AI                               │
│                  (Central Intelligence)                      │
└─────────────────────────┬───────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
        ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│   FarmBot     │ │    Rachio     │ │   Ecowitt     │
│   REST API    │ │     API       │ │  (via HA)     │
└───────┬───────┘ └───────┬───────┘ └───────┬───────┘
        │                 │                 │
        ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│  Garden Beds  │ │  Lawn/Orchard │ │    Sensors    │
│  - Planting   │ │  - Irrigation │ │  - Soil       │
│  - Watering   │ │  - Scheduling │ │  - Weather    │
│  - Weeding    │ │  - Zones      │ │  - Temp/Humid │
└───────────────┘ └───────────────┘ └───────────────┘
                          │
                          ▼
              ┌───────────────────────┐
              │    HOME ASSISTANT     │
              │   (Local Hub/Bridge)  │
              └───────────────────────┘
```

---

## Implementation Phases

### Phase 1: Foundation (Month 1-2)
- [ ] Purchase FarmBot Genesis XL v1.8 (~$4,000)
- [ ] Build raised bed to spec (3m x 6m)
- [ ] Assemble FarmBot (~4 hours)
- [ ] Configure FarmBot OS
- [ ] Connect to WiFi
- [ ] Grant Legacy AI API access
- [ ] Test basic sequences

### Phase 2: Sensor Network (Month 2-3)
- [ ] Purchase Ecowitt GW2000 gateway
- [ ] Deploy WH51 soil sensors (4-8 units)
- [ ] Install weather station
- [ ] Configure Home Assistant integration
- [ ] Create monitoring dashboard
- [ ] Set up alerts for soil moisture thresholds

### Phase 3: Irrigation Expansion (Month 3-4)
- [ ] Install Rachio 3 controller (8 or 16 zone)
- [ ] Map irrigation zones (lawn, orchard, beds)
- [ ] Configure weather intelligence
- [ ] Connect to Home Assistant
- [ ] Grant Legacy AI API access
- [ ] Create unified watering schedule

### Phase 4: Automation & Intelligence (Month 4-6)
- [ ] Create plant database in FarmBot
- [ ] Build watering sequences per crop type
- [ ] Set up camera-based monitoring
- [ ] Implement weather-adjusted schedules
- [ ] Create Legacy AI automation routines
- [ ] Track yield data

### Phase 5: Greenhouse (Optional, Month 6+)
- [ ] Build/install greenhouse structure
- [ ] Deploy climate control system
- [ ] Integrate with Home Assistant
- [ ] Extend FarmBot or add second unit

---

## Budget Estimate

| Item | Cost |
|------|------|
| FarmBot Genesis XL v1.8 | $4,000 |
| Raised bed materials | $300 |
| Rachio 3 (8-zone) | $230 |
| Ecowitt GW2000 + WH51 x4 | $150 |
| Misc (hoses, fittings, wire) | $200 |
| **Total Phase 1-3** | **~$4,880** |

---

## Legacy AI Capabilities Once Integrated

### Monitoring
- Real-time soil moisture across all beds
- Weather conditions and forecast
- Plant health via camera
- Irrigation zone status

### Control
- Execute planting sequences
- Trigger precision watering
- Run weeding routines
- Adjust irrigation schedules
- Control greenhouse climate

### Intelligence
- Recommend planting times based on weather
- Alert on low soil moisture
- Track growth progress
- Optimize water usage
- Predict harvest dates
- Generate yield reports

### Reporting
- Daily garden status brief
- Weekly growth summaries
- Water usage analytics
- Harvest tracking
- Season-over-season comparison

---

## Resources & Links

| Resource | URL |
|----------|-----|
| FarmBot Store | https://farm.bot |
| FarmBot Docs | https://software.farm.bot |
| FarmBot API | https://developer.farm.bot |
| FarmBot GitHub | https://github.com/FarmBot |
| FarmBot Forum | https://forum.farm.bot |
| Rachio | https://rachio.com |
| Rachio API Docs | https://rachio-api.readme.io |
| Ecowitt | https://www.ecowitt.com |
| Home Assistant Rachio | https://www.home-assistant.io/integrations/rachio |
| Home Assistant Ecowitt | https://www.home-assistant.io/integrations/ecowitt |

---

## Navigation

- Parent: [[_Delacruz - Homestead]]
- Pillar: LE (Legacy Evolution)
- Area: 50 Home (Sanctuary & Property)

---

**Last Updated:** 2026-01-10
**Research Sources:** FarmBot official docs, Rachio support, Ecowitt docs, Home Assistant integrations, web research
