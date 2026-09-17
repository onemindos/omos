---
title: "{{Device Name}}"
type: device
status: active
created: {{DATE}}
updated: {{DATE}}
owner: zeus
access_level: L0
tags: [type/device, domain/SO, status/active]
ai_summary: "{{Device description}} — {{owner name}}'s {{phone/Mac/robot/tablet}}."
---

# {{Device Name}}

## Identity

| Field | Value |
|-------|-------|
| Name | {{e.g. "Zeus's iPhone", "OM1 Robot"}} |
| Type | phone / mac / robot / tablet / server / smart-home |
| Make / Model | {{Apple iPhone 16 / MacBook Pro M4 / OpenMind OM1}} |
| Serial / IMEI | {{serial}} |
| OS / Firmware | {{iOS 18 / macOS 15 / custom}} |
| Owned by | {{zeus / wife / kid-name}} |

## Connectivity

| Field | Value |
|-------|-------|
| Tailscale IP | {{100.x.x.x}} |
| Local IP | {{192.168.x.x}} |
| MAC address | {{if relevant}} |
| OpenClaw node | Yes / No |
| Node token | In 1Password: `{{path}}` |
| Home Assistant entity | {{entity_id if integrated}} |

## Purpose

{{What this device is used for, what agents interact with it}}

## Notes

{{Repair history, warranties, important configs}}


---

## Related

- [[00-24 SO (Sovereign Operations)/08 OM Assets (Entity Registry) (The Registry) (Hive Mind)/08 OM Assets|08 OM Assets]]
