---
title: "Car 1"
type: asset
subtype: vehicle
status: active
created: 2026-04-12
updated: 2026-04-12
owner: zeus
access_level: L0
tags: [type/asset, subtype/vehicle, domain/SO, status/active]
ai_summary: "Primary car — Delacruz family vehicle. Fill in make, model, plate. Bouncie GPS tracker installed."
---

# Car 1

> **Action required:** Fill in year, make, model, VIN, plate, and Bouncie device ID.

## Identity

| Field | Value |
|-------|-------|
| Name / callsign | Car 1 |
| Year | {{year}} |
| Make | {{Make}} |
| Model | {{Model}} |
| Color | {{color}} |
| License plate | {{plate}} |
| VIN | {{VIN}} |
| Type | car |

## Live Tracking

| Field | Value |
|-------|-------|
| GPS tracker | Bouncie |
| Bouncie device ID | {{device ID — from Bouncie dashboard or skill}} |
| Live location query | Ask Legacy: *"Where is Car 1?"* |
| Skill | `.claude/skills/bouncie/SKILL.md` |

## Insurance & Registration

| Field | Value |
|-------|-------|
| Insurance provider | {{provider}} |
| Policy number | In 1Password: `{{path}}` |
| Registration expires | {{date}} |

## Maintenance

| Field | Value |
|-------|-------|
| Last oil change | {{date / mileage}} |
| Next service due | {{date / mileage}} |
| Service center | {{name + address}} |

## Notes

{{Primary driver, storage location, mods, quirks}}


---

## Related

- [[00-24 SO (Sovereign Operations)/01 The Fabric (Node Registry)/01 The Fabric|The Fabric]]
