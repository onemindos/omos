---
title: "{{Location Name}}"
type: location
status: active
created: {{DATE}}
updated: {{DATE}}
owner: zeus
tags: [type/location, domain/SO, status/active]
ai_summary: "{{Location description}} — {{home / office / travel hub / garage}}."
---

# {{Location Name}}

## Identity

| Field | Value |
|-------|-------|
| Name | {{e.g. "Home", "NYC Office", "Hangar"}} |
| Type | home / office / storage / hangar / travel-hub |
| Address | {{full address}} |
| Coordinates | {{lat, lon — optional}} |
| Timezone | America/New_York |

## Access & Entry

| Field | Value |
|-------|-------|
| Smart locks | Yes / No — Home Assistant entity: {{entity_id}} |
| Gate / garage | {{Home Assistant entity}} |
| Alarm system | {{provider / panel model}} |
| Alarm code | In 1Password: `{{path}}` |
| WiFi SSID | {{network name}} |
| WiFi password | In 1Password: `{{path}}` |

## Devices at This Location

{{List devices registered at this location}}
- [[zeus-iphone]]
- [[home-mac]]
- [[om1-robot]]

## People at This Location

{{Who lives/works here}}
- [[zeus]]
- [[wife-name]]

## Notes

{{Lease/mortgage info, property manager, emergency contacts for property}}
