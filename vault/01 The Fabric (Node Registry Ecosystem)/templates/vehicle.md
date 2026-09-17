---
title: "{{Vehicle Name}}"
type: asset
subtype: vehicle
status: active
created: {{DATE}}
updated: {{DATE}}
owner: zeus
access_level: L0
tags: [type/asset, subtype/vehicle, domain/SO, status/active]
ai_summary: "{{Year Make Model}} — {{color}}. {{Primary use}}. Tracked via {{Bouncie / none}}."
---

# {{Vehicle Name}}

## Identity

| Field | Value |
|-------|-------|
| Name / callsign | {{e.g. "The Truck", "Black Bird"}} |
| Year | {{year}} |
| Make | {{Make}} |
| Model | {{Model}} |
| Color | {{color}} |
| License plate | {{plate}} |
| VIN | {{VIN}} |
| Type | car / truck / motorcycle / plane / boat |

## Tracking

| Field | Value |
|-------|-------|
| GPS tracker | Bouncie / AirTag / None |
| Bouncie device ID | {{device ID from Bouncie dashboard}} |
| Live location | Via `.claude/skills/bouncie/SKILL.md` |
| Bouncie trips | `GET /devices/{device_id}/trips` |

> For active location, ask Agent: *"Where is the [vehicle name]?"* — Legacy will call the Bouncie skill.

## Insurance & Registration

| Field | Value |
|-------|-------|
| Insurance provider | {{provider}} |
| Policy number | {{policy}} |
| Expiration | {{date}} |
| Registration expires | {{date}} |

## Maintenance

| Field | Value |
|-------|-------|
| Last oil change | {{date / mileage}} |
| Next service due | {{date / mileage}} |
| Service center | {{name + address}} |

## Notes

{{Any additional context — mods, quirks, storage location}}


---

## Related

- [[00-24 SO (Sovereign Operations)/01 The Fabric (Node Registry)/01 The Fabric|The Fabric]]
