---
title: "Amcrest"
type: reference
status: planned
tags: [platform, iot, surveillance, cameras, hardware]
---
# Amcrest

> HA-compatible IP cameras. Native RTSP, ONVIF, PoE options. Good bang for buck.

## Key Features
| Feature | Notes |
|---------|-------|
| RTSP | Native streaming (no cloud required) |
| ONVIF | Standard protocol — works with Frigate/Scrypted/HA |
| PoE | Power over Ethernet (single cable) |
| Pan/Tilt | PTZ models available |
| Local storage | microSD recording |
| Amcrest Cloud | Optional paid cloud (not needed) |

## RTSP URLs
```
Main stream: rtsp://<ip>:554/cam/realmonitor?channel=1&subtype=0
Sub stream:  rtsp://<ip>:554/cam/realmonitor?channel=1&subtype=1
```

## Home Assistant
- Integration: ONVIF or Amcrest native
- Frigate: Add RTSP URL as camera source
