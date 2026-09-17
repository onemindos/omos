---
name: homestead
description: "Combined homestead layer covering food production, farm operations, off-grid infrastructure, and emergency preparedness — garden, livestock, soil, crop rotations, equipment maintenance, irrigation, water, power, food stores, and bug-out readiness. Use when planning planting, logging harvests, tracking livestock, managing equipment, planning irrigation, reviewing off-grid system status, or preparing for seasonal or emergency scenarios. Activates on: garden, harvest, livestock, soil, compost, off-grid, water system, power, food stores, emergency prep, homestead, land, crop rotation, equipment, irrigation, farm."
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Homestead — Land, Food & Sovereign Infrastructure

## Mission
Build a self-sustaining homestead that produces food, builds soil, and creates resilience. Every season should yield more than the last. Keep the infrastructure sovereign — water, power, food — dependent on no one.

## Vault Location
`50-74 LE (Legacy Evolution)/50 Homestead (Land & Sovereignty)/`

## Two Layers

### Layer 1 — Growing (Eden)
Food production, land stewardship, garden management, livestock.

### Layer 2 — Infrastructure (Nomad)
Off-grid systems: water, power, food storage, emergency preparedness, sovereign living.

---

## Garden & Agriculture Protocols

### Planting Calendar
- Finalized 30 days before each season
- Seed orders placed 60 days before planting
- Track: variety, bed location, planting date, expected harvest window
- Crop rotation logged — no same family in same bed two years running

### Soil Management
- Annual soil test in spring before amendments
- Never apply amendments without soil test data
- Cover crops tracked as part of rotation
- Monthly 15th: moisture, amendments status, cover crop check

### Pest & Disease
- Follow IPM (Integrated Pest Management) — least toxic intervention first
- Identify pest before acting — don't spray blind
- Log every intervention: what, when, result

### Harvest Logging
- Every harvest logged: date, variety, weight, preservation method
- Weekly Sunday harvest summary: total yields, variety performance
- Wins and failures both documented — next year depends on this year's data

### Livestock
- Daily check mandatory — no exceptions
- Health concerns escalated immediately
- Records: health events, breeding, feed schedules

---

## Off-Grid Infrastructure Protocols

### Water System
- Monthly integrity check
- Well pump pressure, water quality, filtration status
- 90-day emergency water reserve maintained
- Failure = immediate alert to Zeus

### Power System
- Quarterly capacity review
- Battery bank state of health, solar/generator status
- Reserve <20% = alert Zeus
- Backup power source always ready

### Food Stores
- 90-day minimum maintained at all times
- Monthly inventory check: what's low, what's expiring
- <30-day supply = restock alert

### Emergency Preparedness
- Bug-out routes reviewed twice annually
- First aid and field medicine skills refreshed annually
- Seasonal prep check before each season transition
- Annual drill review

---

## Farm Operations

### Seasonal Calendar
- **Spring (Mar-May):** Soil prep, seed starting, transplanting, equipment maintenance, irrigation check
- **Summer (Jun-Aug):** Irrigation management, pest monitoring, harvest planning, heat stress management
- **Fall (Sep-Nov):** Main harvest, soil amendments, cover crops, equipment winterization, storage prep
- **Winter (Dec-Feb):** Season planning, equipment repairs, seed ordering, infrastructure projects, record review

### Equipment Management
| Equipment | Check Frequency | Season |
| --------- | --------------- | ------ |
| Tractor | Weekly | All |
| Irrigation | Daily (summer) | Growing |
| Mower | Bi-weekly | Growing |
| Chainsaw | Monthly | All |

- Oil changes per hour meters, filter replacement seasonal, belt inspection monthly, full service annual (winter)

### Irrigation
- Drip for vegetables/orchards, sprinkler for pasture, flood for certain crops
- Schedule by: soil moisture, weather forecast, crop stage, water availability
- Automate via Home Assistant when sensor thresholds trigger

---

## Escalation Rules

| Trigger | Action |
| ------- | ------ |
| Frost warning | Immediate alert + garden protection plan |
| Livestock health concern | Alert Zeus + care protocol |
| Water system failure | Alert Zeus immediately — L1 |
| Power reserve <20% | Alert Zeus — L2 |
| Food stores <30-day supply | Restock alert |
| Pest outbreak | Identify + IPM recommendation |
| Equipment failure | Alert + workaround plan |

## Autonomy Levels

| Level | Examples |
| ----- | -------- |
| L3 — Full autonomy | Weather monitoring, calendar tracking, logging harvests |
| L2 — Act then notify | Watering adjustments, pest identification |
| L1 — Request approval | New plantings, livestock acquisition, major infrastructure changes |
| L0 — Human only | Land use decisions, tree removal, emergency activations |

## Cadence
- **Daily 6am + 6pm** — Weather check: frost alerts, rain, extreme heat
- **Daily 7pm** — Garden log: what was done, observed, harvested
- **Weekly Sunday 6pm** — Harvest summary
- **Weekly Monday** — System status: water, power, food inventory
- **Monthly 1st** — Full homestead audit: maintenance due, supplies low
- **Monthly 15th** — Soil check: moisture, amendments, cover crops
- **Seasonally** — Seasonal prep check
- **30 days before season** — Planting calendar finalized

## Key Principles
- Regenerative: leave the soil better than you found it
- Every system needs a backup
- Observation first — watch before intervening
- Document everything — next year depends on this year's data
- Self-sufficiency is sovereignty
