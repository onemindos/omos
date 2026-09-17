# Frigate

> **AI-powered NVR with real-time object detection**

## Overview

Frigate is an open-source NVR built around real-time AI object detection. Uses OpenCV and TensorFlow to detect people, cars, animals, and custom objects in camera feeds.

## Status

| Attribute   | Value                     |
| ----------- | ------------------------- |
| Status      | ✅ Active                 |
| Category    | Smart Home / Surveillance |
| Type        | Network Video Recorder    |
| Integration | MQTT, Home Assistant      |

## Use Cases in OneMind

1. **Security Monitoring** - Detect people, vehicles, packages
2. **AI Alerts** - Smart notifications with object context
3. **Event Recording** - Clip storage on detection
4. **Home Assistant Integration** - Automations on detection

## Key Features

- Real-time object detection (person, car, dog, cat, etc.)
- Google Coral TPU support for efficient inference
- RTSP camera support
- MQTT event publishing
- Home Assistant integration
- Clip & snapshot storage
- Zones and masks
- Re-streaming capabilities

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Cameras   │────▶│   Frigate   │────▶│    MQTT     │
│   (RTSP)    │     │  (Detect)   │     │   Events    │
└─────────────┘     └─────────────┘     └─────────────┘
                           │                   │
                           ▼                   ▼
                    ┌─────────────┐     ┌─────────────┐
                    │   Storage   │     │ Home Assist │
                    │  (Clips)    │     │ Automations │
                    └─────────────┘     └─────────────┘
```

## Installation

### Docker Compose

```yaml
version: "3.9"
services:
  frigate:
    container_name: frigate
    image: ghcr.io/blakeblackshear/frigate:stable
    privileged: true
    restart: unless-stopped
    shm_size: "256mb"
    devices:
      - /dev/bus/usb:/dev/bus/usb # Coral USB
    volumes:
      - /etc/localtime:/etc/localtime:ro
      - ./config:/config
      - ./media:/media/frigate
    ports:
      - "5000:5000" # Web UI
      - "8554:8554" # RTSP
      - "8555:8555" # WebRTC
    environment:
      FRIGATE_RTSP_PASSWORD: "your_password"
```

## Configuration

### frigate.yml

```yaml
mqtt:
  enabled: true
  host: mqtt.local
  port: 1883
  topic_prefix: frigate
  user: frigate
  password: "{FRIGATE_MQTT_PASSWORD}"

detectors:
  coral:
    type: edgetpu
    device: usb

cameras:
  front_door:
    ffmpeg:
      inputs:
        - path: rtsp://user:pass@camera-ip:554/stream
          roles:
            - detect
            - record
    detect:
      width: 1280
      height: 720
      fps: 5
    objects:
      track:
        - person
        - car
        - package
    zones:
      porch:
        coordinates: 0,300,400,300,400,720,0,720
    record:
      enabled: true
      retain:
        days: 7
      events:
        retain:
          default: 14

  backyard:
    ffmpeg:
      inputs:
        - path: rtsp://user:pass@camera2-ip:554/stream
          roles:
            - detect
    detect:
      width: 1920
      height: 1080
      fps: 5
    objects:
      track:
        - person
        - dog
        - cat
```

## MQTT Events

```json
{
  "type": "new",
  "before": {},
  "after": {
    "id": "1234567890.123456-abc123",
    "camera": "front_door",
    "label": "person",
    "score": 0.89,
    "start_time": 1234567890.123,
    "end_time": null,
    "zone": ["porch"],
    "thumbnail": "/api/events/1234567890/thumbnail.jpg"
  }
}
```

## API Integration

### REST API

```python
import httpx

FRIGATE_URL = "http://frigate.local:5000"

async def get_events(
    camera: str = None,
    label: str = None,
    limit: int = 20
):
    """Get recent events from Frigate."""
    params = {"limit": limit}
    if camera:
        params["camera"] = camera
    if label:
        params["label"] = label

    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{FRIGATE_URL}/api/events",
            params=params
        )
        return response.json()

async def get_snapshot(camera: str):
    """Get latest snapshot from camera."""
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{FRIGATE_URL}/api/{camera}/latest.jpg"
        )
        return response.content

async def get_event_clip(event_id: str):
    """Get video clip for an event."""
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{FRIGATE_URL}/api/events/{event_id}/clip.mp4"
        )
        return response.content
```

## Legacy AI Integration

```python
from agno import Tool

@Tool
async def check_camera(camera: str) -> dict:
    """Get current status from a security camera.

    Args:
        camera: Camera name (front_door, backyard, etc.)
    """
    events = await get_events(camera=camera, limit=5)
    snapshot = f"{FRIGATE_URL}/api/{camera}/latest.jpg"

    return {
        "camera": camera,
        "recent_events": events,
        "snapshot_url": snapshot
    }

@Tool
async def get_security_events(
    hours: int = 24,
    object_type: str = "person"
) -> list:
    """Get security events from all cameras.

    Args:
        hours: How many hours back to search
        object_type: Type of object (person, car, package)
    """
    events = await get_events(label=object_type, limit=50)
    return events

@Tool
async def who_is_at_door() -> dict:
    """Check who is currently at the front door."""
    events = await get_events(camera="front_door", label="person", limit=1)
    snapshot = await get_snapshot("front_door")

    return {
        "has_person": len(events) > 0,
        "event": events[0] if events else None,
        "snapshot": snapshot
    }
```

## Home Assistant Integration

```yaml
# configuration.yaml
frigate:
  host: frigate.local
  port: 5000

automation:
  - alias: "Person at front door"
    trigger:
      - platform: mqtt
        topic: frigate/events
    condition:
      - condition: template
        value_template: "{{ trigger.payload_json.after.label == 'person' }}"
      - condition: template
        value_template: "{{ trigger.payload_json.after.camera == 'front_door' }}"
    action:
      - service: notify.mobile_app
        data:
          title: "Someone at door"
          message: "Person detected at front door"
          data:
            image: "/api/frigate/notifications/{{ trigger.payload_json.after.id }}/thumbnail.jpg"
```

## Environment Variables

```bash
FRIGATE_URL=http://frigate.local:5000
FRIGATE_MQTT_HOST=mqtt.local
FRIGATE_MQTT_USER=frigate
FRIGATE_MQTT_PASSWORD=your_password
FRIGATE_RTSP_PASSWORD=camera_password
```

## Hardware Recommendations

| Component | Recommendation           |
| --------- | ------------------------ |
| Coral TPU | USB or M.2 for detection |
| CPU       | Intel with QuickSync     |
| RAM       | 4GB+ per 4 cameras       |
| Storage   | SSD for recordings       |

## Best Practices

1. **Use Coral TPU** - Dramatically improves performance
2. **Tune detection zones** - Reduce false positives
3. **Set object filters** - Min/max size, threshold scores
4. **Separate detect/record streams** - Use substreams for detection
5. **MQTT for events** - Real-time integrations

## Resources

- [Frigate Docs](https://docs.frigate.video/)
- [Frigate GitHub](https://github.com/blakeblackshear/frigate)
- [Home Assistant Integration](https://docs.frigate.video/integrations/home-assistant)

---

Last Updated: January 25, 2026
