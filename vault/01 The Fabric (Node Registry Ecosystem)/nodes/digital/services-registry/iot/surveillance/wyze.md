---
title: "Wyze"
type: reference
status: planned
tags: [platform, iot, surveillance, cameras, budget]
---
# Wyze

> Budget cameras ($20-40). RTSP via firmware hack or docker-wyze-bridge for local streaming.

## Key Features
| Feature | Notes |
|---------|-------|
| Price | $20-40 per camera |
| RTSP | Via firmware flash or docker-wyze-bridge |
| Person detection | On-device (Cam v3+) |
| Color night vision | Starlight sensor |
| Indoor/Outdoor | Weatherproof models available |
| Pan cam | PTZ for $30 |

## Local Access (no cloud)
```bash
# docker-wyze-bridge — streams Wyze cams as RTSP locally
docker run -d --name wyze-bridge \
  -p 8554:8554 -p 8888:8888 -p 5000:5000 \
  -e WYZE_EMAIL=your@email.com \
  -e WYZE_PASSWORD=password \
  mrlt8/docker-wyze-bridge
```

## RTSP URLs (via bridge)
```
rtsp://localhost:8554/<camera-nickname>
```

## Integration
- docker-wyze-bridge → Frigate (RTSP source)
- Home Assistant: Wyze integration or via bridge
