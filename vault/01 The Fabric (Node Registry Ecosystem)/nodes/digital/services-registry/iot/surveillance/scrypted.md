---
title: "Scrypted"
type: reference
status: planned
tags: [platform, iot, surveillance, cameras, homekit, nvr]
---
# Scrypted

> Camera streaming hub — makes any camera work with HomeKit Secure Video, Google Home, Alexa, and Home Assistant. Hardware transcoding, local processing.

## Key Features
| Feature | Notes |
|---------|-------|
| HomeKit Secure Video | Full HKSV support for any camera |
| Hardware transcoding | GPU/iGPU accelerated |
| HA integration | Direct camera feed to Home Assistant |
| Plugins | ONVIF, RTSP, Amcrest, Reolink, Wyze, Unifi |
| Object detection | Built-in or Frigate/CodeProject.AI backend |
| NVR | Recording with timeline + scrubbing |

## Deploy
```bash
docker run -d --name scrypted \
  --network host \
  -v scrypted-data:/server/volume \
  koush/scrypted
```

## Integration
- Pairs with Frigate (Frigate = AI detection, Scrypted = streaming/HomeKit bridge)
- ONVIF discovery for auto-adding cameras
