# OneMind OS Interfaces - Complete Reference

> **Created:** 2026-01-25
> **Purpose:** Document all communication interfaces in OneMind OS

---

## System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              ONEMIND OS INTERFACE LAYER                              │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                      │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │                           USER INTERFACES                                    │   │
│  ├─────────────────────────────────────────────────────────────────────────────┤   │
│  │                                                                              │   │
│  │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │   │
│  │   │  Web App    │  │ Mobile App  │  │  Desktop    │  │    CLI      │       │   │
│  │   │  (Flutter)  │  │  (Flutter)  │  │  (Flutter)  │  │  (Future)   │       │   │
│  │   │  Port 3000  │  │  iOS/Android│  │ macOS/Win/  │  │             │       │   │
│  │   │             │  │             │  │   Linux     │  │             │       │   │
│  │   └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘       │   │
│  │          │                │                │                │              │   │
│  │          └────────────────┴────────────────┴────────────────┘              │   │
│  │                                    │                                        │   │
│  │   ┌─────────────┐  ┌─────────────┐│┌─────────────┐  ┌─────────────┐       │   │
│  │   │Frame Glasses│  │  Smartwatch │││   Voice     │  │   Browser   │       │   │
│  │   │(Brilliant)  │  │  (Garmin)   │││  (LiveKit)  │  │  Extension  │       │   │
│  │   │    BLE      │  │    MQTT     │││   WebRTC    │  │    HTTP     │       │   │
│  │   └──────┬──────┘  └──────┬──────┘│└──────┬──────┘  └──────┬──────┘       │   │
│  │          │                │       │       │                │              │   │
│  └──────────┼────────────────┼───────┼───────┼────────────────┼──────────────┘   │
│             │                │       │       │                │                   │
│  ┌──────────▼────────────────▼───────▼───────▼────────────────▼──────────────┐   │
│  │                           API GATEWAY (Port 8081)                          │   │
│  ├───────────────────────────────────────────────────────────────────────────┤   │
│  │   REST API  │  WebSocket Streaming  │  SSE (Server-Sent Events)          │   │
│  │   /api/*    │  /api/ws/*            │  /api/agents/{id}/runs (stream)    │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                              │
│  ┌─────────────────────────────────▼─────────────────────────────────────────┐   │
│  │                          EVENT BACKBONE                                    │   │
│  ├───────────────────────────────────────────────────────────────────────────┤   │
│  │                                                                            │   │
│  │   ┌─────────────────────────────────────────────────────────────────┐     │   │
│  │   │                    NATS JetStream                                │     │   │
│  │   ├─────────────────────────────────────────────────────────────────┤     │   │
│  │   │  Port 4222  │  Port 9222  │  Port 1883  │  Port 7422           │     │   │
│  │   │  Native     │  WebSocket  │  MQTT       │  Leaf Nodes          │     │   │
│  │   │  (Backend)  │  (UI)       │  (IoT)      │  (Edge)              │     │   │
│  │   └─────────────────────────────────────────────────────────────────┘     │   │
│  │                                                                            │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                              │
│  ┌─────────────────────────────────▼─────────────────────────────────────────┐   │
│  │                          EDGE DEVICES                                      │   │
│  ├───────────────────────────────────────────────────────────────────────────┤   │
│  │                                                                            │   │
│  │   ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐             │   │
│  │   │   Home    │  │   Farm    │  │  Office   │  │    Car    │             │   │
│  │   │  (Jetson) │  │  (Jetson) │  │   (Pi)    │  │   (Pi)    │             │   │
│  │   │           │  │           │  │           │  │           │             │   │
│  │   │ NATS Leaf │  │ NATS Leaf │  │ NATS Leaf │  │ NATS Leaf │             │   │
│  │   │ Ollama    │  │ Ollama    │  │ Ollama    │  │ Ollama    │             │   │
│  │   │ Frigate   │  │ Frigate   │  │           │  │           │             │   │
│  │   │ HA        │  │ Robot     │  │           │  │           │             │   │
│  │   └───────────┘  └───────────┘  └───────────┘  └───────────┘             │   │
│  │                                                                            │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                                                                   │
└───────────────────────────────────────────────────────────────────────────────────┘
```

---

## Interface Reference Table

### Core APIs

| Interface | Protocol | Port | Purpose | Authentication |
|-----------|----------|------|---------|----------------|
| **REST API** | HTTP/HTTPS | 8081 | All backend operations | API Key / JWT |
| **WebSocket Stream** | WS | 8081 | Agent response streaming | Session token |
| **SSE Stream** | HTTP | 8081 | Real-time agent runs | Session token |
| **GraphQL** | HTTP | - | Future: unified queries | - |

### Event System (NATS)

| Interface | Protocol | Port | Purpose | Clients |
|-----------|----------|------|---------|---------|
| **NATS Native** | NATS | 4222 | Backend service communication | Python, Go services |
| **NATS WebSocket** | WS | 9222 | UI real-time updates | Browser, Flutter |
| **NATS MQTT** | MQTT | 1883 | IoT devices | Home Assistant, sensors |
| **NATS Leaf** | NATS | 7422 | Edge location connections | Jetson, Raspberry Pi |
| **NATS HTTP** | HTTP | 8222 | Monitoring dashboard | Admin browser |

### External Services

| Interface | Protocol | Port | Purpose | Provider |
|-----------|----------|------|---------|----------|
| **LiveKit** | WebRTC | External | Voice/video real-time | LiveKit Cloud |
| **Deepgram** | WebSocket | External | Speech-to-text | Deepgram API |
| **Cartesia** | HTTP | External | Text-to-speech | Cartesia API |
| **AWS Bedrock** | HTTP | External | LLM (Claude) | AWS |

### Wearables

| Interface | Protocol | Port | Purpose | Device |
|-----------|----------|------|---------|--------|
| **Bluetooth LE** | BLE | - | Frame glasses, sensors | Brilliant Labs Frame |
| **MQTT** | MQTT | 1883 | Watch data | Garmin, Fitbit |
| **HTTP** | HTTP | - | Health APIs | Oura, Exist.io |

---

## Detailed Interface Specifications

### 1. REST API (Port 8081)

**Base URL:** `http://100.102.21.44:8081/api`

```
API Structure:
─────────────────────────────────────────────────────
/api
├── /agents                    Agent management
│   ├── GET /                  List agents
│   ├── GET /{id}              Get agent details
│   ├── POST /{id}/runs        Execute agent (SSE stream)
│   └── POST /{id}/runs/{run}/cancel  Cancel run
│
├── /sessions                  Session management
│   ├── GET /                  List sessions
│   ├── GET /{id}              Get session
│   ├── POST /                 Create session
│   └── DELETE /{id}           Delete session
│
├── /memory                    Memory operations
│   ├── GET /                  List memories
│   ├── POST /                 Create memory
│   ├── PATCH /{id}            Update memory
│   └── DELETE /{id}           Delete memory
│
├── /inbox                     Unified inbox
│   ├── GET /                  List items
│   ├── POST /                 Create item
│   ├── POST /{id}/dispatch    Dispatch to agent
│   └── WS /ws                 Real-time updates
│
├── /knowledge                 RAG knowledge base
│   ├── POST /upload           Upload document
│   ├── POST /url              Index URL
│   ├── GET /search            Search knowledge
│   └── GET /documents         List documents
│
├── /learning                  Learning machine
│   ├── GET /profile           User profile
│   ├── GET /entities          List entities
│   ├── GET /culture           Agent culture
│   └── GET /stats             Learning stats
│
├── /awareness                 Awareness system
│   ├── GET /                  Current state
│   └── PUT /                  Update mode
│
├── /workflows                 Workflow engine
│   ├── GET /                  List workflows
│   └── POST /{name}           Execute workflow
│
└── /voice                     Voice integration
    ├── GET /token             Get LiveKit token
    └── GET /status            Voice status
```

### 2. WebSocket (Port 9222)

**Connection URL:** `ws://100.102.21.44:9222`

```javascript
// Message format
{
  "type": "inbox.item.created",  // Event type
  "data": {                       // Payload
    "id": "uuid",
    "title": "New notification",
    "body": "Content here",
    "priority": "normal"
  }
}

// Subject patterns
"inbox.{user_id}"          // User's inbox events
"events.awareness.*"       // Awareness changes
"events.agent.*"           // Agent status
"events.workflow.*"        // Workflow progress
```

### 3. MQTT (Port 1883)

**Broker URL:** `mqtt://100.102.21.44:1883`

```
Topic Structure:
─────────────────────────────────────────────────────
onemind/
├── events/
│   ├── motion/{location}      Motion detection
│   ├── temperature/{sensor}   Temperature readings
│   └── presence/{device}      Device presence
│
├── commands/
│   ├── lights/{room}          Light control
│   ├── climate/{zone}         HVAC control
│   └── robot/{id}             Robot commands
│
└── status/
    ├── devices/{id}           Device health
    └── locations/{id}         Location status
```

### 4. LiveKit (Voice/Video)

**Protocol:** WebRTC via LiveKit SDK

```dart
// Flutter integration
final room = Room();
await room.connect(
  'wss://your-livekit-server.livekit.cloud',
  token,  // From GET /api/voice/token
);

// Audio tracks
room.onTrackSubscribed.listen((track) {
  if (track is AudioTrack) {
    // Play agent's voice response
  }
});
```

### 5. BLE (Brilliant Labs Frame)

**Protocol:** Bluetooth Low Energy

```python
# Python SDK
from frame import Frame

async with Frame() as f:
    # Display text on glasses
    await f.display.show_text("Hello from Legacy")

    # Capture photo
    photo = await f.camera.capture()

    # Record audio
    audio = await f.microphone.record(duration=5)
```

---

## Data Flow Diagrams

### User → Agent → Response

```
┌────────────┐      ┌────────────┐      ┌────────────┐      ┌────────────┐
│   User     │      │  Flutter   │      │   API      │      │   Agent    │
│            │      │    UI      │      │  (8081)    │      │  (Legacy)  │
└─────┬──────┘      └─────┬──────┘      └─────┬──────┘      └─────┬──────┘
      │                   │                   │                   │
      │ 1. Type message   │                   │                   │
      │ ─────────────────►│                   │                   │
      │                   │                   │                   │
      │                   │ 2. POST /agents   │                   │
      │                   │    /{id}/runs     │                   │
      │                   │ ─────────────────►│                   │
      │                   │                   │                   │
      │                   │                   │ 3. Execute        │
      │                   │                   │ ─────────────────►│
      │                   │                   │                   │
      │                   │                   │ 4. SSE Stream     │
      │                   │                   │ ◄─ ─ ─ ─ ─ ─ ─ ─ ─│
      │                   │                   │    (tokens)       │
      │                   │                   │                   │
      │                   │ 5. Stream chunks  │                   │
      │                   │ ◄─ ─ ─ ─ ─ ─ ─ ─ ─│                   │
      │                   │                   │                   │
      │ 6. Display text   │                   │                   │
      │ ◄─ ─ ─ ─ ─ ─ ─ ─ ─│                   │                   │
      │    (real-time)    │                   │                   │
      ▼                   ▼                   ▼                   ▼
```

### IoT Device → NATS → Backend

```
┌────────────┐      ┌────────────┐      ┌────────────┐      ┌────────────┐
│  Sensor    │      │   NATS     │      │  Backend   │      │   Agent    │
│  (ESP32)   │      │  (MQTT)    │      │            │      │  (Legacy)  │
└─────┬──────┘      └─────┬──────┘      └─────┬──────┘      └─────┬──────┘
      │                   │                   │                   │
      │ 1. MQTT publish   │                   │                   │
      │    temp/kitchen   │                   │                   │
      │ ─────────────────►│                   │                   │
      │                   │                   │                   │
      │                   │ 2. Route to       │                   │
      │                   │    subject        │                   │
      │                   │ ─────────────────►│                   │
      │                   │                   │                   │
      │                   │                   │ 3. Process        │
      │                   │                   │    (if needed)    │
      │                   │                   │ ─────────────────►│
      │                   │                   │                   │
      │                   │                   │                   │ 4. "Kitchen
      │                   │                   │                   │     is 85°F,
      │                   │                   │                   │     turn on AC?"
      ▼                   ▼                   ▼                   ▼
```

### Edge Location → Cloud

```
┌────────────┐      ┌────────────┐      ┌────────────┐      ┌────────────┐
│   Robot    │      │ Edge NATS  │      │ Cloud NATS │      │  Legacy    │
│ (Unitree)  │      │   (Leaf)   │      │   (Hub)    │      │  (Cloud)   │
└─────┬──────┘      └─────┬──────┘      └─────┬──────┘      └─────┬──────┘
      │                   │                   │                   │
      │ 1. Telemetry      │                   │                   │
      │ ─────────────────►│                   │                   │
      │                   │                   │                   │
      │                   │ 2. Forward        │                   │
      │                   │    (leaf → hub)   │                   │
      │                   │ ─────────────────►│                   │
      │                   │                   │                   │
      │                   │                   │ 3. Deliver        │
      │                   │                   │ ─────────────────►│
      │                   │                   │                   │
      │                   │                   │ 4. Command        │
      │                   │                   │ ◄─────────────────│
      │                   │                   │                   │
      │                   │ 5. Route back     │                   │
      │                   │ ◄─────────────────│                   │
      │                   │                   │                   │
      │ 6. Execute        │                   │                   │
      │ ◄─────────────────│                   │                   │
      ▼                   ▼                   ▼                   ▼
```

---

## Authentication & Security

| Interface | Auth Method | Notes |
|-----------|-------------|-------|
| REST API | API Key header | `X-API-Key: xxx` |
| WebSocket (9222) | Connection token | Passed on connect |
| MQTT (1883) | Username/password | Device credentials |
| LiveKit | JWT token | Short-lived, from `/api/voice/token` |
| BLE | Pairing code | Device-specific |

---

## Quick Reference

### Connect from Flutter

```dart
// REST API
final response = await http.get(
  Uri.parse('http://100.102.21.44:8081/api/agents'),
  headers: {'X-API-Key': apiKey},
);

// WebSocket (real-time)
final ws = await WebSocket.connect('ws://100.102.21.44:9222');
ws.add(json.encode({'subscribe': 'inbox.zeus'}));

// SSE Stream (agent runs)
final request = http.Request('POST', Uri.parse('.../runs'));
final response = await request.send();
response.stream.transform(utf8.decoder).listen((chunk) {
  // Process streaming response
});
```

### Connect from Python

```python
# REST API
import httpx
response = httpx.get('http://100.102.21.44:8081/api/agents')

# NATS
import nats
nc = await nats.connect("nats://100.102.21.44:4222")
await nc.publish("inbox.zeus", b"Hello")

# MQTT
import paho.mqtt.client as mqtt
client = mqtt.Client()
client.connect("100.102.21.44", 1883)
client.publish("onemind/events/temperature/kitchen", "72.5")
```

### Connect from JavaScript

```javascript
// REST API
const response = await fetch('http://100.102.21.44:8081/api/agents');

// WebSocket
const ws = new WebSocket('ws://100.102.21.44:9222');
ws.onmessage = (e) => console.log(JSON.parse(e.data));

// MQTT (via mqtt.js)
const client = mqtt.connect('ws://100.102.21.44:9222');
client.subscribe('onemind/events/#');
```

---

## Related Documentation

- `docs/architecture/NATS-WEBSOCKET-9222.md` - WebSocket details
- `docs/architecture/MULTI-LOCATION-NATS.md` - Edge connectivity
- `docs/architecture/REMOTEAGENT-ARCHITECTURE.md` - Edge agents
- `docs/UNIFIED-INBOX-SYSTEM.md` - Inbox system
- `docs/legacy/VOICE_INTEGRATION.md` - Voice/LiveKit

---

*Document created for OneMind OS - Zeus Delacruz*
