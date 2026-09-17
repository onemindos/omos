---
title: "{{Drone Name}}"
type: asset
subtype: drone
status: active
created: {{DATE}}
updated: {{DATE}}
owner: zeus
access_level: L0
tags: [type/asset, subtype/drone, domain/SO, status/active]
ai_summary: "{{Make Model}} drone — {{primary use}}. FAA registration: {{number}}."
---

# {{Drone Name}}

## Identity

| Field | Value |
|-------|-------|
| Name / callsign | {{e.g. "Eye in the Sky"}} |
| Make | {{DJI / custom / other}} |
| Model | {{Model}} |
| Serial number | {{serial}} |
| FAA registration | {{FA-XXXXXXX}} |
| Max range | {{km/miles}} |
| Max flight time | {{minutes}} |
| Camera | {{resolution / model}} |

## Software & Connectivity

| Field | Value |
|-------|-------|
| Controller app | DJI Fly / Autel Sky / custom |
| Live stream | RTSP / WebRTC / none |
| Stream endpoint | {{rtsp://...}} |
| Home Assistant entity | {{entity_id if integrated}} |
| Telemetry source | {{MAVLink / DJI SDK / none}} |

## Regulatory

| Field | Value |
|-------|-------|
| FAA Part 107 | Yes / No / In progress |
| Pilot cert | {{certificate number}} |
| Geofencing | DJI standard restrictions |
| Allowed zones | {{authorization zones if any}} |

## Batteries & Equipment

- Batteries: {{count, capacity mAh}}
- Charger: {{model}}
- Storage: {{location}}
- Carry case: {{model/location}}

## Notes

{{Flight logs location, repair history, planned upgrades}}


---

## Related

- [[00-24 SO (Sovereign Operations)/01 The Fabric (Node Registry)/01 The Fabric|The Fabric]]
