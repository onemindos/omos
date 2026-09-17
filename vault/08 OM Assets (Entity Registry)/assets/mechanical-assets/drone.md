---
title: "Drone"
type: asset
subtype: drone
status: active
created: 2026-04-12
updated: 2026-04-12
owner: zeus
access_level: L0
tags: [type/asset, subtype/drone, domain/SO, status/active]
ai_summary: "Zeus's drone. Fill in make, model, FAA registration. Primary use: aerial photography / inspection."
---

# Drone

> **Action required:** Fill in make, model, serial, FAA registration number.

## Identity

| Field | Value |
|-------|-------|
| Name / callsign | {{e.g. "Eye in the Sky"}} |
| Make | {{DJI / Autel / custom}} |
| Model | {{Mini 4 Pro / Air 3 / Mavic 3 / etc.}} |
| Serial number | {{serial}} |
| FAA registration | {{FA-XXXXXXX}} |
| Max range | {{km}} |
| Max flight time | {{minutes}} |
| Camera | {{resolution}} |

## Software & Connectivity

| Field | Value |
|-------|-------|
| Controller app | DJI Fly / Autel Sky / custom |
| Live stream | RTSP / none |
| Stream endpoint | {{rtsp://... if configured}} |
| Home Assistant entity | {{entity_id if integrated}} |

## Regulatory

| Field | Value |
|-------|-------|
| FAA Part 107 pilot | {{Yes / No / cert number}} |
| FAA registration expires | {{date}} |
| LAANC authorization | {{Yes / No / zones}} |

## Equipment

- Batteries: {{count, mAh each}}
- Charging hub: {{model}}
- ND filters: {{sets}}
- Carry case: {{model / location}}
- Storage: {{where kept}}

## Notes

{{Flight log location, planned projects, preferred flying spots}}


---

## Related

- [[00-24 SO (Sovereign Operations)/08 OM Assets (Entity Registry) (The Registry) (Hive Mind)/08 OM Assets|08 OM Assets]]
