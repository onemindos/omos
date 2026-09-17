# LiveKit + NATS Architecture Boundaries

> Created: 2026-01-23
> Purpose: Clear separation of concerns between media (LiveKit) and data/control (NATS)

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    OneMind OS Architecture                      │
└─────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│ MEDIA LAYER (LiveKit - Real-Time Audio/Video)                   │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Voice & Video Only:                                            │
│  ✅ Human ↔ Legacy (Gemini Live voice)                          │
│  ✅ Multi-location speakers (home, office, farm)                │
│  ✅ Robot camera live view (on-demand, when human watching)     │
│  ✅ Multi-device voice (iPad, Mac, iPhone same session)         │
│  ✅ Video calls with humans                                     │
│                                                                  │
│  Protocol: WebRTC (UDP/TCP with TURN fallback)                 │
│  Latency: 50-150ms                                              │
│  Cost: FREE (self-hosted unlimited)                             │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
                           ↕ (Media streams)
┌──────────────────────────────────────────────────────────────────┐
│ APPLICATION LAYER (Agno Agents)                                 │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Legacy Agent:                                                  │
│  • Voice processing (Gemini Live integration)                   │
│  • Multimodal understanding (voice + context)                   │
│  • Command generation                                           │
│  • State management                                             │
│  • Tool execution                                               │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
                           ↕ (Commands, events, state)
┌──────────────────────────────────────────────────────────────────┐
│ EVENT/CONTROL LAYER (NATS - Everything Else)                   │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Commands & Control:                                            │
│  ✅ Robot movement commands (events.commands.robot.*)           │
│  ✅ Device control (lights, locks, etc.)                        │
│  ✅ Workflow triggers (events.workflow.*)                       │
│  ✅ Tool execution commands                                     │
│                                                                  │
│  Telemetry & Sensors:                                           │
│  ✅ Biometrics (heart rate, SpO2, steps)                        │
│  ✅ Robot sensors (battery, position, status)                   │
│  ✅ Camera AI detections (Frigate → NATS)                       │
│  ✅ Environmental sensors (temp, humidity, etc.)                │
│                                                                  │
│  State & Sync:                                                  │
│  ✅ Awareness mode changes                                      │
│  ✅ Device presence (which device user is on)                   │
│  ✅ Physical presence (which room user is in)                   │
│  ✅ Session state sync across devices                           │
│  ✅ Memory updates                                              │
│                                                                  │
│  Integrations:                                                  │
│  ✅ Todoist webhooks → NATS                                     │
│  ✅ GitHub webhooks → NATS                                      │
│  ✅ Calendar events → NATS                                      │
│  ✅ Email notifications → NATS                                  │
│                                                                  │
│  Protocol: NATS (pub/sub + JetStream persistence)              │
│  Latency: 5-20ms                                                │
│  Cost: FREE (already deployed)                                  │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
                           ↕ (Events, snapshots)
┌──────────────────────────────────────────────────────────────────┐
│ CAMERA AI LAYER (Frigate NVR)                                  │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  24/7 Recording & AI Detection:                                 │
│  ✅ 6x Yi cameras → local disk recording                        │
│  ✅ AI detection (person, vehicle, animal) → NATS events        │
│  ✅ Snapshots on motion → stored locally                        │
│  ✅ Live view on-demand → LiveKit (only when human requests)    │
│                                                                  │
│  Why NOT LiveKit for 24/7?                                      │
│  • Cloud: 6 cameras × 24hr = $129.60/day = $3,888/month 😱     │
│  • Self-hosted: Still unnecessary load on LiveKit               │
│  • Frigate: Purpose-built, efficient, FREE                      │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│ DECISION TREE: Which Layer?                                     │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Is it audio/video that a HUMAN needs to see/hear in real-time? │
│    → YES: LiveKit ✅                                            │
│    → NO: ↓                                                      │
│                                                                  │
│  Is it a command, event, sensor data, or state update?          │
│    → YES: NATS ✅                                               │
│    → NO: ↓                                                      │
│                                                                  │
│  Is it 24/7 camera recording?                                   │
│    → YES: Frigate ✅ (send events to NATS, live view via        │
│              LiveKit only when human requests)                  │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

## Real-World Examples

### ✅ Use LiveKit

**Voice conversation**:
```python
# User: "Hey Legacy, what's on my calendar today?"
# Flow: Human voice → LiveKit → Gemini Live → Legacy response voice → Human
```

**Watch robot camera while controlling**:
```python
# User: "Show me the barn camera"
# Flow: Robot camera feed → LiveKit → Human (while giving commands via NATS)
```

**Multi-location presence**:
```python
# Zeus speaks in office → Legacy responds through home speakers
# Flow: LiveKit connects all locations in same room
```

---

### ✅ Use NATS

**Robot command**:
```python
# User: "Legacy, patrol the barn"
# Flow: NATS publish("commands.farm.robot.move", {"action": "patrol"})
```

**Biometric alert**:
```python
# Watch detects heart rate >180 BPM
# Flow: NATS publish("events.biometric.alert", {"hr": 185, "type": "critical"})
```

**Camera AI detection**:
```python
# Frigate detects person in barn at 2 AM
# Flow: NATS publish("events.camera.detection", {"type": "person", "location": "barn"})
# → Legacy decides: Should I alert? Send snapshot? Turn on lights?
```

---

### ✅ Use Frigate

**24/7 camera recording** (NOT LiveKit!):
```python
# 6 cameras recording to disk 24/7
# Flow: Frigate records locally, detects with AI, sends events to NATS
# → LiveKit: ONLY when human requests "show me live camera"
```

---

## Cost Comparison: LiveKit Cloud vs Self-Hosted

### Full Omnipresent System

| Component | LiveKit Cloud | LiveKit Self-Hosted | Alternative |
|-----------|---------------|---------------------|-------------|
| **Voice (2hr/day)** | $54/mo | Server included | N/A |
| **6 cameras (24/7)** | $3,888/mo 😱 | Server included | Frigate FREE |
| **Robot control** | Not designed | Not designed | NATS FREE |
| **Telemetry** | Not designed | Not designed | NATS FREE |
| **Live robot view (1hr/day)** | $27/mo | Server included | N/A |
| **TOTAL** | $3,969/mo 😱 | ~$15-20/mo | $0* + voice |

*With Frigate + NATS (already in stack)

---

## Key Principles

### When to Use LiveKit
- Real-time audio/video for humans
- Low latency required (<200ms)
- WebRTC features needed (NAT traversal, adaptive bitrate)
- Interactive sessions

### When to Use NATS
- Commands, events, telemetry
- Pub/sub patterns
- High throughput (1000+ msg/sec)
- Persistence required (JetStream)
- MQTT device integration

### When to Use Frigate
- 24/7 camera recording
- AI detection on video
- Motion detection
- Local storage
- Cost-sensitive video needs

---

## Integration Points

### LiveKit → NATS
```python
# When LiveKit room events occur, publish to NATS
@livekit_room.on("participant_connected")
async def on_participant_connected(participant):
    await nats.publish(
        "events.voice.participant_joined",
        {"room": room_name, "participant": participant.identity}
    )
```

### NATS → LiveKit
```python
# When NATS event requires human notification, use voice
async def handle_critical_alert(msg):
    alert = json.loads(msg.data)
    if alert["type"] == "critical":
        # Send voice notification through LiveKit
        await legacy.speak_through_livekit(
            "Critical alert: " + alert["message"],
            location="all"  # All speakers
        )
```

### Frigate → NATS → LiveKit
```python
# Camera detects motion
Frigate AI Detection
    ↓
NATS event: events.camera.detection
    ↓
Legacy evaluates severity
    ↓
If human needs to see:
    ↓
Start LiveKit stream (on-demand)
```

---

## Summary

**LiveKit**: Media only (voice, video) - real-time human interaction
**NATS**: Everything else (commands, events, state, telemetry)
**Frigate**: Camera AI (24/7 recording, detection events)

**They complement each other** - don't try to use one for the other's job.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
