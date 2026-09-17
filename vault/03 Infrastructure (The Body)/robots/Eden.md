---
title: "Eden — Garden & Greenhouse Robot"
type: asset
domain: SO
status: planned
created: 2026-05-26
updated: 2026-05-26
tags: [robot, eden, garden, greenhouse, iot, fleet, onemind-os]
related: ["50 Home (Sanctuary & Property)"]
---

# Eden

**Type:** Autonomous Garden & Greenhouse Robot
**Role:** Head of Garden & Super Greenhouse Operations
**Status:** Planned

---

## Mission

Eden is the designated robot in charge of the garden and super greenhouse. She manages all growing operations — monitoring, maintenance, and automation of the full growing environment.

---

## Responsibilities

- Full oversight of the garden and super greenhouse
- Environmental monitoring (soil, humidity, temperature, light)
- Watering, feeding, and plant lifecycle management
- Pest and health detection
- Harvest tracking and scheduling
- Reporting growing conditions and yields back to OneMind OS

---

## OneMind OS Integration

Eden operates as a node on the NATS fabric:

- **Publisher (Sensor):** Reports environmental readings, plant health, alerts
- **Subject (Pathway):** `eden.garden.>` — all garden telemetry and events
- **Subscriber (Reaction):** Receives commands and schedules from OneMind OS

Example subjects:
- `eden.garden.soil.moisture` — soil moisture readings
- `eden.garden.alert.pest` — pest detection alert
- `eden.greenhouse.temp` — greenhouse temperature
- `cmd.eden.water` — trigger watering cycle

---

## Notes

- Part of the broader OneMind robot fleet
- Lives at the home property — see [[50 Home (Sanctuary & Property)]]
- Future: integrate with Home Assistant for unified home + garden automation

---
