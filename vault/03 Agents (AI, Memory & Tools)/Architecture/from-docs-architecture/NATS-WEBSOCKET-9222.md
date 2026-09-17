# NATS WebSocket (Port 9222) - Educational Guide

> **Created:** 2026-01-25
> **Purpose:** Understanding real-time UI communication in OneMind OS

---

## What is Port 9222?

Port 9222 is the **NATS WebSocket** port - the real-time communication channel between the backend and frontend UI.

```
┌─────────────────────────────────────────────────────────────────┐
│                     NATS PORT ARCHITECTURE                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Port 4222  ─────►  NATS Native Protocol (backend services)     │
│  Port 1883  ─────►  MQTT Protocol (IoT devices, Home Assistant) │
│  Port 7422  ─────►  Leaf Node (edge location connections)       │
│  Port 8222  ─────►  HTTP Monitoring (NATS dashboard)            │
│  Port 9222  ─────►  WebSocket (UI real-time updates) ◄── THIS   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Why WebSocket?

Traditional HTTP is **request-response** - the client asks, server answers. But for real-time features, we need the server to **push updates** to the client without being asked.

| Protocol | Direction | Use Case |
|----------|-----------|----------|
| HTTP | Client → Server → Client | API calls, data fetching |
| WebSocket | Bidirectional (always open) | Real-time updates, notifications |

**WebSocket enables:**
- Instant inbox notifications (no page refresh needed)
- Live agent status updates
- Real-time awareness mode changes
- Live task progress tracking

---

## How It Works

### 1. Connection Flow

```
┌─────────────────┐          ┌─────────────────┐          ┌─────────────────┐
│   Flutter UI    │          │   NATS Server   │          │   Backend API   │
│                 │          │   (Port 9222)   │          │                 │
└────────┬────────┘          └────────┬────────┘          └────────┬────────┘
         │                            │                            │
         │ 1. ws://100.102.21.44:9222  │                            │
         │ ─────────────────────────► │                            │
         │                            │                            │
         │ 2. Connection established  │                            │
         │ ◄───────────────────────── │                            │
         │                            │                            │
         │ 3. Subscribe to topics     │                            │
         │ ─────────────────────────► │                            │
         │    (inbox.zeus, events.*)  │                            │
         │                            │                            │
         │                            │ 4. Backend publishes event │
         │                            │ ◄───────────────────────── │
         │                            │    (new inbox item)        │
         │                            │                            │
         │ 5. Push to subscriber      │                            │
         │ ◄───────────────────────── │                            │
         │    (instant notification)  │                            │
         │                            │                            │
         ▼                            ▼                            ▼
```

### 2. Message Topics (Subjects)

NATS uses **subjects** (topics) to route messages:

```
Subject Pattern              Description
─────────────────────────────────────────────────────────────────
inbox.{user_id}             Inbox updates for specific user
inbox.zeus                  Zeus's inbox notifications
events.awareness.*          Awareness mode changes
events.agent.*              Agent status updates
events.workflow.*           Workflow progress
events.task.*               Task updates
events.*.critical           Critical alerts (any source)
```

### 3. Backend Publishing (Python)

```python
# core/inbox/api.py
async def broadcast_item_update(item: InboxItem, event_type: str):
    """Broadcast inbox item update via WebSocket."""
    await ws_manager.broadcast(
        item.user_id,
        json.dumps({
            "type": event_type,          # "inbox.item.created"
            "data": item.model_dump()    # The actual item data
        })
    )

# Called after creating inbox item:
@router.post("", response_model=InboxItem)
async def create_inbox_item(data: InboxItemCreate):
    item = await service.create(data)
    await broadcast_item_update(item, "inbox.item.created")  # ← Push to UI
    return item
```

### 4. Frontend Subscribing (Flutter)

```dart
// ui/lib/core/services/nats_service.dart
class NatsService {
  WebSocket? _socket;
  final _streamController = StreamController<Map<String, dynamic>>.broadcast();

  Future<void> connect(String userId) async {
    _socket = await WebSocket.connect('ws://100.102.21.44:9222');

    // Subscribe to user's inbox
    _socket!.add(json.encode({
      'action': 'subscribe',
      'subject': 'inbox.$userId'
    }));

    // Listen for messages
    _socket!.listen((data) {
      final message = json.decode(data);
      _streamController.add(message);
    });
  }

  Stream<Map<String, dynamic>> get messages => _streamController.stream;
}

// In your widget:
natsService.messages.listen((message) {
  if (message['type'] == 'inbox.item.created') {
    // Add to inbox list, show notification, etc.
    setState(() => inboxItems.insert(0, message['data']));
  }
});
```

---

## NATS Configuration

### Docker Compose

```yaml
# infra/compose/docker-compose.core.yml
nats:
  image: nats:2.10-alpine
  container_name: onemind-nats
  ports:
    - "4222:4222"   # Native NATS protocol
    - "1883:1883"   # MQTT for IoT devices
    - "7422:7422"   # Leaf node connections (edge)
    - "8222:8222"   # HTTP monitoring dashboard
    - "9222:9222"   # WebSocket for UI ← THIS
  volumes:
    - ./../../configs/nats/hub.conf:/etc/nats/nats.conf:ro
  command: ["-c", "/etc/nats/nats.conf"]
```

### NATS Server Config

```conf
# infra/configs/nats/hub.conf
port: 4222

# WebSocket for browser/UI connections
websocket {
  port: 9222
  no_tls: true
}

# MQTT for IoT devices
mqtt {
  port: 1883
}

# Leaf node connections for edge locations
leafnodes {
  port: 7422
}

# HTTP monitoring
http_port: 8222

# JetStream for persistence
jetstream {
  store_dir: /data
  max_mem: 256MB
  max_file: 1GB
}
```

---

## What Uses Port 9222?

| Component | Purpose | Example |
|-----------|---------|---------|
| **Inbox Screen** | Real-time notifications | New email arrives → appears instantly |
| **Chat Screen** | Agent status updates | Agent starts typing → UI shows indicator |
| **Dashboard** | Live stats | Task completed → counter updates |
| **Awareness** | Mode changes | User changes mode → all screens update |
| **Workflows** | Progress tracking | Step completes → progress bar advances |

---

## Testing the Connection

### 1. Check NATS is Running

```bash
# On server
docker logs onemind-nats | grep "Listening for websocket"
# Should show: Listening for websocket clients on ws://0.0.0.0:9222

# Check ports
docker port onemind-nats
# 9222/tcp -> 0.0.0.0:9222
```

### 2. Test WebSocket Connection

```bash
# Using websocat (install: brew install websocat)
websocat ws://100.102.21.44:9222

# Or using curl (just checks if port is open)
curl -v http://100.102.21.44:9222
# Should fail with "Upgrade required" - that's correct for WebSocket!
```

### 3. Test from Browser Console

```javascript
// Open browser console on any page
const ws = new WebSocket('ws://100.102.21.44:9222');
ws.onopen = () => console.log('Connected to NATS WebSocket!');
ws.onmessage = (e) => console.log('Message:', e.data);
ws.onerror = (e) => console.error('Error:', e);
```

---

## Common Issues

### "Connection Refused"

```
Error: WebSocket connection to 'ws://100.102.21.44:9222' failed
```

**Causes:**
1. NATS container not running → `docker start onemind-nats`
2. Port not exposed → Check docker-compose ports section
3. Firewall blocking → Allow port 9222

### "Messages Not Arriving"

1. Check subscription subject matches publish subject
2. Verify `broadcast_item_update()` is being called
3. Check NATS logs: `docker logs onemind-nats -f`

---

## Why NATS Instead of Direct WebSocket?

| Feature | Direct WebSocket | NATS WebSocket |
|---------|------------------|----------------|
| Pub/Sub patterns | Manual implementation | Built-in |
| Message persistence | Custom code needed | JetStream |
| Multiple subscribers | Complex | Native |
| Message routing | Custom | Subject-based |
| Replay missed messages | Not possible | JetStream replay |
| Clustering | Very complex | Native |

**NATS gives us:**
- **Subject-based routing** - `inbox.zeus` vs `inbox.guest`
- **JetStream** - Messages persist if client is offline
- **Fan-out** - One publish → many subscribers
- **Edge connectivity** - Same protocol for cloud and edge devices

---

## Related Documentation

- [NATS Official Docs](https://docs.nats.io/)
- [NATS WebSocket](https://docs.nats.io/running-a-nats-service/configuration/websocket)
- `docs/architecture/MULTI-LOCATION-NATS.md` - Hub/leaf architecture
- `docs/UNIFIED-INBOX-SYSTEM.md` - How inbox uses WebSocket
- `infra/configs/nats/hub.conf` - Server configuration

---

*Document created for educational purposes - Zeus Delacruz*
