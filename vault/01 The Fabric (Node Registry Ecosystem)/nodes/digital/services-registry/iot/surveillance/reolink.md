---
title: "Reolink"
type: reference
status: planned
tags: [platform, iot, surveillance, cameras, hardware]
---
# Reolink

> Affordable PoE cameras + NVR. Good image quality, local storage, RTSP/ONVIF support.

## Key Features
| Feature | Notes |
|---------|-------|
| PoE | Single-cable power + data |
| 4K | Up to 12MP resolution |
| RTSP/ONVIF | Works with Frigate/Scrypted/HA |
| NVR | Dedicated hardware NVR or microSD |
| AI detection | Person/vehicle (on-device) |
| Spotlight | Integrated deterrent lights |
| Local only | No cloud dependency (optional Reolink app) |

## RTSP URLs
```
Main: rtsp://<ip>:554/h264Preview_01_main
Sub:  rtsp://<ip>:554/h264Preview_01_sub
```

## Home Assistant
- Integration: Reolink native (REST API)
- Frigate: Add RTSP URL as camera source
