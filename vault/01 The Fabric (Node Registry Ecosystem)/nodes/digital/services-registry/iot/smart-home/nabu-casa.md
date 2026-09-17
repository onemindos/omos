# Nabu Casa

> **Cloud access and subscription service for Home Assistant**

## Overview

Nabu Casa is the company behind Home Assistant and provides the optional cloud subscription (`home-assistant.cloud`). Enables remote access to HA without port forwarding, voice assistant integration (Alexa, Google), and funds open-source HA development.

## Status

| Attribute   | Value                          |
| ----------- | ------------------------------ |
| Status      | ✅ Active                       |
| Category    | IoT / Smart Home               |
| Type        | Cloud Service (HA companion)   |
| Integration | Home Assistant, Voice Assistants|

## Use Cases in OneMind

1. **Remote HA Access** — Access Home Assistant dashboard from anywhere (no VPN needed)
2. **Agent Access** — OpenClaw agents call HA REST API from VPS via Nabu Casa URL
3. **Voice Integration** — Connect HA to Amazon Alexa and Google Home
4. **Text-to-Speech** — Cloud TTS for HA notifications
5. **HA Funding** — Subscription directly funds open-source development

## Key Features

- Remote access URL (`*.ui.nabu.casa`)
- No port forwarding required
- End-to-end encrypted
- Alexa and Google Home integration
- Cloud TTS (text-to-speech for HA)
- $7.50/month subscription
- 1-month free trial

## Remote Access Example

```bash
# External access URL
curl -H "Authorization: Bearer HA_TOKEN" \
  https://YOUR_INSTANCE.ui.nabu.casa/api/states
```

## Links

- Website: [nabucasa.com](https://nabucasa.com)
- Setup: [home-assistant.io/integrations/cloud](https://home-assistant.io/integrations/cloud)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
