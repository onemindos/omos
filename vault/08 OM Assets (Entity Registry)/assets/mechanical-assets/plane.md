---
title: "Plane"
type: asset
subtype: vehicle
status: active
created: 2026-04-12
updated: 2026-04-12
owner: zeus
access_level: L0
tags: [type/asset, subtype/vehicle, subtype/aircraft, domain/SO, status/active]
ai_summary: "Zeus's aircraft. Fill in N-number, make, model. Based at {{airport}}. FAA registration required."
---

# Plane

> **Action required:** Fill in N-number, make, model, home airport, and insurance.

## Identity

| Field | Value |
|-------|-------|
| Name / callsign | {{e.g. "The Bird", "N12345"}} |
| N-number | {{N-XXXXX}} |
| Year | {{year}} |
| Make | {{Cessna / Piper / Cirrus / etc.}} |
| Model | {{Model}} |
| Engine | {{type + hours}} |
| Color / markings | {{description}} |
| Type | fixed-wing / helicopter / light sport |

## Based At

| Field | Value |
|-------|-------|
| Home airport | {{ICAO code}} — {{airport name}} |
| Hangar | {{hangar number / tie-down}} |
| FBO | {{name + phone}} |
| Fuel type | 100LL / Jet-A / Auto |

## Tracking & Live Data

| Field | Value |
|-------|-------|
| ADS-B transponder | {{mode S / mode C / none}} |
| FlightAware | {{link or tail number}} |
| FlightRadar24 | {{link or tail number}} |
| Live tracking | Via FlightAware API (ask Legacy: *"Where is the plane?"*) |

## Regulatory

| Field | Value |
|-------|-------|
| FAA registration expires | {{date}} |
| Annual inspection due | {{date}} |
| Insurance | {{provider + policy}} |
| Pilot certificate | In 1Password: `{{path}}` |
| Medical certificate | {{class + expiry}} |

## Avionics & Equipment

- {{Garmin G1000 / GTN 750 / other avionics}}
- {{ADS-B in/out}}
- {{Emergency ELT}}

## Maintenance

| Field | Value |
|-------|-------|
| Last annual | {{date / tach time}} |
| Next annual due | {{date}} |
| A&P Mechanic | {{name + phone}} |
| Total time | {{airframe hours}} |
| Engine SMOH | {{hours since major overhaul}} |

## Notes

{{Fuel burn, range, typical missions, IFR capable yes/no}}


---

## Related

- [[00-24 SO (Sovereign Operations)/08 OM Assets (Entity Registry) (The Registry) (Hive Mind)/08 OM Assets|08 OM Assets]]
