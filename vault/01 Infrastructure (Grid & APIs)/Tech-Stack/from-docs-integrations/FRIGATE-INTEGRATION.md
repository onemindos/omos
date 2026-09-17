# Frigate Integration for Omnipresent Mode

> **Decision:** Use Frigate NVR for camera AI instead of building detection in Legacy

---

## Why Frigate?

### Frigate vs Pulling Frames

| Approach | How It Works | Pros | Cons |
|----------|--------------|------|------|
| **Pulling Frames** | Legacy requests frames from cameras, runs AI detection | Full control | High CPU/GPU, complex, reinventing wheel |
| **Frigate Events** | Frigate watches cameras 24/7, sends events when things happen | Efficient, battle-tested, pre-built | Dependency on Frigate |

### The Difference Explained

**Pulling Frames (DIY approach):**
```
Every 1 second:
  1. Legacy requests frame from camera (network I/O)
  2. Decode RTSP stream (CPU)
  3. Run YOLO/detection model (GPU/CPU intensive)
  4. Run face recognition (more GPU/CPU)
  5. Process results
  6. Repeat for ALL 6 cameras = 6x compute

Result: High resource usage even when nothing is happening
```

**Frigate Events (Recommended):**
```
Frigate runs 24/7:
  1. Watches all cameras efficiently (hardware decoding)
  2. Only runs AI when motion detected (smart)
  3. Sends event: "person detected at front_door"
  4. Legacy receives event, already knows what happened

Result: Low resource usage, events only when relevant
```

---

## Is Frigate the Best Option?

### Comparison of NVR/Detection Options

| Solution | Open Source | Detection | Face Recognition | HA Integration | Cost |
|----------|-------------|-----------|------------------|----------------|------|
| **Frigate** | ✅ Yes | ✅ YOLO | ✅ Via Double-Take | ✅ Native | Free |
| Blue Iris | ❌ No | ✅ DeepStack | ✅ Yes | 🟡 Plugin | $70 |
| Shinobi | ✅ Yes | 🟡 Limited | ❌ No | ❌ No | Free |
| ZoneMinder | ✅ Yes | 🟡 Basic | ❌ No | 🟡 Plugin | Free |
| AWS Rekognition | ❌ No | ✅ Excellent | ✅ Excellent | ❌ API only | $$$/mo |
| Wyze AI | ❌ No | ✅ Good | ❌ No | ❌ No | $2/mo/cam |

### Verdict: Frigate is Best for OneMind Because

1. **Open Source** - MIT license, full control, no vendor lock-in
2. **Home Assistant Native** - First-class integration, events flow naturally
3. **Efficient** - Only processes when motion detected
4. **Hardware Acceleration** - Supports Coral TPU, Intel QSV, NVIDIA
5. **Face Recognition** - Via CompreFace/Double-Take add-ons
6. **Active Community** - 10k+ GitHub stars, regular updates
7. **Self-Hosted** - Your data stays local

---

## Frigate Overview

### What is Frigate?

Frigate is a **Network Video Recorder (NVR)** with real-time AI object detection.

- **GitHub:** https://github.com/blakeblackshear/frigate
- **License:** MIT (fully open source)
- **Stars:** 10,000+
- **Language:** Python + C++ (for performance)

### What Frigate Does

```
┌─────────────────────────────────────────────────────────────────┐
│                          FRIGATE NVR                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  INPUTS                         OUTPUTS                         │
│  ───────                        ────────                        │
│  • RTSP streams (Yi cameras)    • MQTT events                   │
│  • USB cameras                  • HTTP API                      │
│  • RTMP streams                 • WebSocket live view           │
│                                 • Recordings (24/7 or events)   │
│                                 • Snapshots                     │
│  DETECTION                      • HA integration                │
│  ─────────                                                      │
│  • Person                       EVENT EXAMPLE                   │
│  • Car                          ─────────────                   │
│  • Dog/Cat                      {                               │
│  • Package                        "type": "person",             │
│  • Bicycle                        "camera": "front_door",       │
│  • Motorcycle                     "score": 0.92,                │
│  • Custom objects                 "box": [100, 200, 300, 400],  │
│                                   "timestamp": 1705612800       │
│  FACE RECOGNITION (via addon)   }                               │
│  ─────────────────────────                                      │
│  • CompreFace                                                   │
│  • Double-Take                                                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Do You Need a GPU?

### Short Answer: No, but it helps

### Hardware Options

| Hardware | Cost | Performance | Power | Recommendation |
|----------|------|-------------|-------|----------------|
| **CPU Only** | $0 | 2-5 FPS | High | Works for 1-2 cams |
| **Coral USB TPU** | $25-40 | 100+ FPS | 2W | **Best value** |
| **Coral M.2 TPU** | $25-40 | 100+ FPS | 2W | If M.2 slot available |
| **Intel QSV** | $0 (built-in) | 20-50 FPS | Medium | Good if Intel CPU |
| **NVIDIA GPU** | $200+ | 50-100+ FPS | High | Overkill for Frigate |

### Recommendation for OneMind

**Coral USB TPU ($25-40)** - Best option because:
- Handles 6 cameras easily
- 100+ inferences per second
- Only 2W power consumption
- Plug and play with Frigate
- Frees CPU for other tasks

**Where to buy:**
- Amazon: "Google Coral USB Accelerator"
- Coral.ai official store

### Without Any Accelerator

Frigate will still work on CPU for 6 cameras, but:
- Use lower resolution (720p instead of 1080p for detection)
- Set detect FPS to 5 instead of default
- May need beefy CPU (i5+ or Ryzen 5+)

```yaml
# Frigate config for CPU-only
cameras:
  front_door:
    detect:
      width: 1280
      height: 720
      fps: 5  # Lower FPS for CPU
```

---

## Frigate → NATS Bridge

### What It Does

The bridge subscribes to Frigate MQTT events and republishes them to your NATS consciousness stream.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   FRIGATE   │────►│   BRIDGE    │────►│    NATS     │
│   (MQTT)    │     │             │     │             │
│             │     │ • Transform │     │ • events.*  │
│ frigate/    │     │ • Filter    │     │ • presence.*│
│ events      │     │ • Enrich    │     │ • home.*    │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │   Legacy    │
                                        │             │
                                        │ Already     │
                                        │ knows who   │
                                        │ is where    │
                                        └─────────────┘
```

### Event Flow Example

```
1. Person walks to front door

2. Frigate detects → MQTT event:
   Topic: frigate/events
   Payload: {
     "type": "new",
     "after": {
       "camera": "front_door",
       "label": "person",
       "score": 0.92,
       "snapshot": "/api/events/123/snapshot.jpg"
     }
   }

3. Bridge transforms → NATS event:
   Subject: home.camera.person
   Payload: {
     "camera": "front_door",
     "label": "person",
     "confidence": 0.92,
     "location": "main-house",
     "timestamp": "2024-01-18T10:30:00Z"
   }

4. Legacy receives event:
   - Updates presence: "Someone at front door"
   - If face recognized: "Zeus is at front door"
   - Can proactively act: Turn on porch light, unlock door
```

### Bridge Features

| Feature | Description |
|---------|-------------|
| **Event Translation** | Converts Frigate MQTT → NATS format |
| **Location Mapping** | Maps camera IDs to locations |
| **Filtering** | Only forwards relevant events (person, car, package) |
| **Enrichment** | Adds metadata (location name, zone, time of day) |
| **Face Events** | Routes Double-Take face recognition events |
| **Snapshot URLs** | Includes links to event snapshots |

---

## Architecture with Frigate

### Complete Omnipresent Stack

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           OMNIPRESENT SYSTEM                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  CAMERAS (Yi x6)                                                        │
│  ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌────┐                             │
│  │Live│ │Kitc│ │Offi│ │Fron│ │Back│ │Gara│                             │
│  │Room│ │ hen│ │ ce │ │Door│ │yard│ │ ge │                             │
│  └──┬─┘ └──┬─┘ └──┬─┘ └──┬─┘ └──┬─┘ └──┬─┘                             │
│     │      │      │      │      │      │                                │
│     └──────┴──────┴──────┼──────┴──────┘                                │
│                          │ RTSP                                         │
│                          ▼                                              │
│  ┌───────────────────────────────────────────────────────────────────┐ │
│  │                        FRIGATE NVR                                 │ │
│  │                                                                    │ │
│  │  • Watches all 6 cameras 24/7                                     │ │
│  │  • Runs YOLO detection (person, car, package, animal)             │ │
│  │  • Records clips on events                                         │ │
│  │  • Coral TPU acceleration (optional)                               │ │
│  │                                                                    │ │
│  └───────────────────────────┬───────────────────────────────────────┘ │
│                              │ MQTT Events                              │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐ │
│  │                    DOUBLE-TAKE (Face Recognition)                  │ │
│  │                                                                    │ │
│  │  • Receives person snapshots from Frigate                         │ │
│  │  • Matches against known faces (Zeus, family)                     │ │
│  │  • Publishes face_recognized events                               │ │
│  │                                                                    │ │
│  └───────────────────────────┬───────────────────────────────────────┘ │
│                              │ MQTT Events                              │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐ │
│  │                      HOME ASSISTANT                                │ │
│  │                                                                    │ │
│  │  • Receives all events                                            │ │
│  │  • Triggers automations (lights, locks, notifications)            │ │
│  │  • UI for camera viewing                                          │ │
│  │  • Binary sensors for each camera's detection state               │ │
│  │                                                                    │ │
│  └───────────────────────────┬───────────────────────────────────────┘ │
│                              │ Webhook / MQTT                           │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐ │
│  │                   FRIGATE-NATS BRIDGE                              │ │
│  │                                                                    │ │
│  │  • Subscribes to frigate/# MQTT topics                            │ │
│  │  • Transforms events to NATS format                               │ │
│  │  • Publishes to consciousness stream                              │ │
│  │                                                                    │ │
│  └───────────────────────────┬───────────────────────────────────────┘ │
│                              │                                          │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐ │
│  │                         NATS JETSTREAM                             │ │
│  │                    (Consciousness Stream)                          │ │
│  │                                                                    │ │
│  │  Subjects:                                                         │ │
│  │  • home.camera.person      - Person detected                      │ │
│  │  • home.camera.car         - Vehicle detected                     │ │
│  │  • home.camera.package     - Package detected                     │ │
│  │  • presence.person.zeus    - Zeus face recognized                 │ │
│  │  • presence.person.unknown - Unknown person                       │ │
│  │  • legacy.>                - Legacy responses                     │ │
│  │                                                                    │ │
│  └───────────────────────────┬───────────────────────────────────────┘ │
│                              │                                          │
│                              ▼                                          │
│  ┌───────────────────────────────────────────────────────────────────┐ │
│  │                          Legacy                                    │ │
│  │                    (Agno + Claude Opus 4.5)                        │ │
│  │                                                                    │ │
│  │  • Receives all events in real-time                               │ │
│  │  • Updates presence state automatically                           │ │
│  │  • Can query: "Who's at the front door?"                          │ │
│  │  • Can act: "Turn on porch light when Zeus arrives"               │ │
│  │  • Has context: "Zeus entered through garage 5 min ago"           │ │
│  │                                                                    │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Setup Checklist

### Phase 1: Frigate Core

- [ ] Install Frigate add-on in Home Assistant
- [ ] Configure Yi camera RTSP streams
- [ ] Set up detection zones per camera
- [ ] Test person detection
- [ ] (Optional) Add Coral TPU for acceleration

### Phase 2: Face Recognition

- [ ] Install Double-Take add-on
- [ ] Install CompreFace (face recognition engine)
- [ ] Train faces (Zeus, family members)
- [ ] Test face recognition events

### Phase 3: NATS Bridge

- [ ] Deploy Frigate-NATS bridge
- [ ] Configure camera → location mapping
- [ ] Test events flowing to NATS
- [ ] Verify Legacy receives events

### Phase 4: Integration

- [ ] Update presence.py to consume Frigate events
- [ ] Remove frame-pulling code from feeds.py
- [ ] Test full flow: person → Frigate → NATS → Legacy
- [ ] Test voice: "Who's at the front door?"

---

## Configuration Examples

### Frigate Config (frigate.yml)

```yaml
mqtt:
  host: homeassistant.local
  port: 1883
  user: frigate
  password: !secret mqtt_password

detectors:
  coral:
    type: edgetpu
    device: usb  # or pci for M.2

cameras:
  front_door:
    ffmpeg:
      inputs:
        - path: rtsp://user:pass@192.168.1.101:554/live
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
        coordinates: 0,0,300,0,300,720,0,720
        objects:
          - person

  living_room:
    ffmpeg:
      inputs:
        - path: rtsp://user:pass@192.168.1.102:554/live
          roles:
            - detect
    detect:
      width: 1280
      height: 720
      fps: 5
    objects:
      track:
        - person

  # ... repeat for other cameras
```

### Double-Take Config

```yaml
frigate:
  url: http://frigate:5000
  cameras:
    - front_door
    - living_room

compreface:
  url: http://compreface:8000
  recognition_api_key: your-api-key

mqtt:
  host: homeassistant.local
  topics:
    frigate: frigate/events
    matches: double-take/matches
```

---

## Summary

| Question | Answer |
|----------|--------|
| **Frigate vs Pulling Frames?** | Frigate is way more efficient - only processes on motion |
| **Is Frigate best option?** | Yes, for HA + self-hosted + open source |
| **Is it open source?** | Yes, MIT license, 10k+ GitHub stars |
| **What does NATS bridge do?** | Converts Frigate MQTT → NATS for Legacy |
| **Need GPU?** | No, but $25 Coral TPU recommended for 6 cameras |

---

## Resources

- **Frigate Docs:** https://docs.frigate.video
- **Frigate GitHub:** https://github.com/blakeblackshear/frigate
- **Double-Take:** https://github.com/jakowenko/double-take
- **CompreFace:** https://github.com/exadel-inc/CompreFace
- **Coral TPU:** https://coral.ai/products/

---

## Frigate-NATS Bridge Code

The bridge is implemented at: `agents/legacy/consciousness/frigate_bridge.py`

### Bridge Features

| Feature | Description |
|---------|-------------|
| **MQTT → NATS** | Subscribes to Frigate MQTT, publishes to NATS |
| **Event Filtering** | Only forwards person, car, package, animal (configurable) |
| **Confidence Threshold** | Filters low-confidence detections (< 60%) |
| **Face Recognition** | Routes Double-Take matches to `presence.person.{name}` |
| **Camera Mapping** | Maps camera names to locations |
| **Snapshot URLs** | Includes Frigate snapshot/clip URLs in events |
| **Auto-Reconnect** | Handles MQTT and NATS disconnections |

### NATS Subjects Published

| Subject | When | Priority |
|---------|------|----------|
| `events.home.camera.person` | Person detected | 3 |
| `events.home.camera.car` | Vehicle detected | 2 |
| `events.home.camera.package` | Package detected | 2 |
| `events.home.camera.dog` | Dog detected | 2 |
| `events.home.camera.cat` | Cat detected | 2 |
| `events.presence.person.zeus` | Zeus face recognized | 4 |
| `events.presence.person.unknown` | Unknown face | 3 |
| `events.home.camera.status` | Camera online/offline | 2 |

### Event Payload Example

```json
{
  "type": "home.camera.person",
  "source": "frigate",
  "data": {
    "event_id": "1705612800.123456-abc123",
    "event_type": "new",
    "camera": "front_door",
    "label": "person",
    "confidence": 0.92,
    "location": "main-house",
    "zones": ["porch", "driveway"],
    "entered_zones": ["porch"],
    "has_snapshot": true,
    "has_clip": true,
    "snapshot_url": "http://frigate:5000/api/events/123/snapshot.jpg",
    "clip_url": "http://frigate:5000/api/events/123/clip.mp4",
    "box": [100, 200, 300, 400]
  },
  "timestamp": "2024-01-18T10:30:00Z",
  "priority": 3
}
```

### Running the Bridge

**Standalone:**
```bash
# Set environment variables
export MQTT_HOST=homeassistant.local
export MQTT_PORT=1883
export NATS_URL=nats://localhost:4222
export FRIGATE_URL=http://frigate:5000
export CAMERA_LOCATIONS='{"front_door":"main-house","backyard":"outdoor"}'

# Run bridge
python -m agents.legacy.consciousness.frigate_bridge
```

**As Part of Legacy:**
```python
# In main.py startup
from agents.legacy.consciousness.frigate_bridge import FrigateBridge, BridgeConfig

config = BridgeConfig.from_env()
bridge = FrigateBridge(config)
await bridge.start()
```

**Docker Compose:**
```yaml
frigate-bridge:
  build:
    context: .
    dockerfile: containers/Dockerfile.frigate-bridge
  environment:
    - MQTT_HOST=mosquitto
    - MQTT_PORT=1883
    - NATS_URL=nats://nats:4222
    - FRIGATE_URL=http://frigate:5000
  depends_on:
    - nats
    - frigate
  restart: unless-stopped
```

### Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `MQTT_HOST` | localhost | MQTT broker host |
| `MQTT_PORT` | 1883 | MQTT broker port |
| `MQTT_USER` | (none) | MQTT username |
| `MQTT_PASSWORD` | (none) | MQTT password |
| `NATS_URL` | nats://localhost:4222 | NATS server URL |
| `FRIGATE_URL` | http://localhost:5000 | Frigate API URL |
| `CAMERA_LOCATIONS` | {} | JSON mapping camera → location |

### Camera Location Mapping

Configure which location each camera belongs to:

```json
{
  "front_door": "main-house",
  "living_room": "main-house",
  "kitchen": "main-house",
  "office": "main-house",
  "backyard": "outdoor",
  "garage": "garage",
  "driveway": "outdoor"
}
```

This allows Legacy to know "person detected at main-house" instead of just "person detected at front_door".

### Dependencies

```bash
pip install aiomqtt nats-py
```

---

*This architecture lets Legacy be truly omnipresent without burning compute on frame analysis. Frigate handles the "eyes" efficiently, Legacy handles the "brain".*
