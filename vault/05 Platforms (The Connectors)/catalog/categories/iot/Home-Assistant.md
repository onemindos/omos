# Home Assistant

> **Open-source home automation platform for smart home control**

## Overview

Home Assistant (HA) is the leading open-source home automation platform. Integrates with 3000+ devices and services. Runs locally for privacy, with optional cloud access via Nabu Casa. Agents can control home devices via REST API and webhooks, routing through the LE domain for environment management.

## Status

| Attribute   | Value                               |
| ----------- | ----------------------------------- |
| Status      | ✅ Active                            |
| Category    | IoT / Smart Home                    |
| Type        | Home Automation                     |
| Integration | REST API, Webhooks, MQTT, iOS App   |

## Use Cases in OneMind

1. **Smart Home Control** — Lights, locks, thermostat, security cameras
2. **Agent Commands** — Legacy agent turns on/off devices via REST API
3. **Environment Automation** — Scene triggers linked to codex events
4. **LE Domain** — Home environment management under LE domain
5. **N8N Integration** — Home events trigger N8N automation workflows

## Key Features

- Local-first, privacy by design
- 3000+ integrations (lights, cameras, sensors, etc.)
- Automation editor (visual + YAML)
- Dashboard (Lovelace UI)
- Voice assistant integration (Google, Alexa, Siri)
- REST API + WebSocket API
- HACS (community addon store)

## API Access

```bash
# Get all states
curl -H "Authorization: Bearer HA_TOKEN" \
  http://homeassistant.local:8123/api/states

# Control a device
curl -X POST \
  -H "Authorization: Bearer HA_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"entity_id": "light.living_room"}' \
  http://homeassistant.local:8123/api/services/light/turn_on
```

## Links

- Website: [home-assistant.io](https://home-assistant.io)
- Skill: `.claude/skills/home-assistant/SKILL.md`
- Cloud Access: See Nabu Casa


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
