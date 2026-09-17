---
id: so-00-capability-roadmap-20260430
title: OpenClaw Capability Roadmap — Browser, Devices, Voice
type: design
domain: SO
area: "00"
agent: Grid
status: active
tags:
  - domain/so
  - area/00
  - type/design
  - status/active
  - openclaw
  - browser
  - devices
  - voice
  - nats
updated: 2026-04-30
review: 2026-Q3
related:
  - "[[OPENCLAW-FRAMEWORK-STRUCTURES]]"
  - "[[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/blueprint/agent-and-app/onemind-os-architecture-and-agno-memory]]"
  - "[[99-Review/24 Inbox (Queue) (The Dock)/history 1/operating-system-genesis/HIVE-PROTOCOL]]"
---
# OpenClaw Capability Roadmap — Browser, Devices, Voice

> Three capability layers that extend OpenClaw agents beyond conversation into the physical and digital world.

**Owner**: Grid | **Status**: Planned — not yet implemented in OpenClaw
**Source**: Extracted from pre-OpenClaw enhancement roadmap (Feb 2025). Architecture concepts valid; all implementation code must be rewritten for OpenClaw tool format.

---

## Overview

| Capability | Priority | NATS Integration | Complexity |
|-----------|---------|-----------------|-----------|
| Browser Automation | High | Optional (results published as events) | Medium |
| Device Event System | High | Core (devices ARE NATS publishers) | Low–Medium |
| Voice Pipeline | Medium | Optional (voice events can publish to NATS) | High |

**Key principle**: Don't rebuild what OpenClaw provides. Build only what's missing.

---

## 1. Browser Automation

### What It Enables
Agents can interact with web pages programmatically:
- **Form filling** — visa applications, job applications, registrations
- **Data extraction** — prices, flight status, competitor dashboards
- **Page monitoring** — detect changes, alert on updates
- **Screenshots** — capture dashboards and reports as images
- **JavaScript execution** — interact with dynamic content

### Architecture
```
OpenClaw Agent (Decision)
    ↓ tool call
Browser Tool Layer (Playwright)
    ↓ headless Chromium
Target Website
    ↓ result
Agent (interprets + acts)
```

### OpenClaw Tool Structure (to implement)
Four tools to register in OpenClaw:
- `browser_navigate(url, selector)` — extract text from CSS selector
- `browser_fill_form(url, fields, submit_selector)` — fill + submit form
- `browser_screenshot(url, selector?)` — full page or element screenshot
- `browser_extract_structured(url, {field: selector})` — batch extraction
- `browser_monitor(url, selector, interval_s)` — watch for changes over time
- `browser_execute_js(url, code)` — run custom JavaScript

### Agent Configuration
- Agent: `web-automator` (Grid toolset)
- Tools: `browser_*` preset + search
- Run tools sequentially (not parallel) for browser actions
- Confidence check before form submissions

### Primary Use Cases
1. Check in for a flight automatically
2. Monitor Amazon/product price and alert when it drops
3. Pull competitor pricing on a schedule
4. Take dashboard screenshots for weekly reports
5. Fill out visa/government applications

### Dependencies
```
playwright (Python) — Playwright for Python
playwright install chromium — headless browser
```

---

## 2. Device Event System

### What It Enables
Hardware devices become NATS publishers. Agents subscribe and react. Zero polling — pure event-driven:
- Cameras detect people/objects → agents decide on security response
- Temperature sensors alert → agents adjust smart home or notify
- Door/motion sensors → agents react contextually
- USB device plugged in → agents auto-register

### Why This Maps Perfectly to OpenClaw
NATS JetStream is already the backbone of OpenClaw agent-to-agent messaging. Devices become first-class citizens on the same bus — no new infrastructure needed.

### Architecture
```
Physical Device (camera, sensor, smart home)
    ↓ Device Adapter (normalize + validate)
NATS Event Bus (already running on omos-core)
    ↓ topic subscription
OpenClaw Agent(s) (multiple agents can react to same event)
    ↓ tool calls
Actions (notify, automate, log, escalate)
```

### NATS Topic Schema
```
onemind.devices.*                          # All device events
onemind.devices.camera.person_detected     # Vision detection
onemind.devices.camera.motion_detected     # Motion
onemind.devices.sensor.temperature_high    # Thermal alert
onemind.devices.sensor.humidity_low        # Environmental
onemind.devices.sensor.motion_detected     # PIR sensor
onemind.devices.usb.hotplug                # USB connected
onemind.devices.home.door_opened           # Door sensor
onemind.devices.home.light_turned_on       # Switch event
```

### Standard Event Schema
```json
{
  "event_type": "person_detected",
  "device_id": "living_room_camera_01",
  "device_type": "camera",
  "timestamp": "2026-04-30T14:35:22.123Z",
  "confidence": 0.87,
  "data": {
    "class": "person",
    "bounding_box": [100, 150, 300, 450],
    "snapshot_url": "http://camera/snapshot.jpg"
  },
  "metadata": {
    "location": "living_room",
    "zone": "indoor",
    "floor": 1
  }
}
```

### Device Abstraction Layer (to implement)
Three classes:
- `DeviceEvent` — standard event envelope (to_dict, to_nats_topic)
- `DeviceSource` — abstract base for all devices (start, stop, get_status)
- `DeviceRegistry` — singleton managing all registered device sources

### Priority Device Integrations
| Device | Type | Agent | Use Case |
|--------|------|-------|---------|
| MaixCAM ($50 AI camera) | Camera | Guardian | Security, person detection |
| Sipeed sensors | Environmental | Grid | Temperature, humidity alerts |
| Home Assistant | Smart home | Eden/Grid | Door, lights, energy |
| OM1 Robot vision | Camera (ROS2) | Grid | Obstacle detection, navigation |
| USB devices | Hotplug | Grid | Auto-register peripherals |

### MaixCAM Integration Note
MaixCAM runs on-device inference at 30fps. Detection API endpoint at `http://{ip}:8080/detect` returns JSON detections. Subscribe to its stream, filter by confidence threshold, publish to NATS. Total latency: <100ms from detection to NATS publish.

---

## 3. Voice Pipeline

### What It Enables
Full hands-free voice interaction loop anywhere on the network:
- Always-listening wake word ("Hey OneMind" / "Hey OM1")
- Fast speech-to-text transcription
- Agent processes + responds
- Audio response played through speakers or streamed to remote device

### Architecture
```
Microphone (always on, local)
    ↓ audio stream
Wake Word Detector (Porcupine — offline, <100ms)
    ↓ wake word triggered
Audio Capture (VAD stops on silence)
    ↓ audio buffer
STT (Groq Whisper — <300ms)
    ↓ transcribed text
OpenClaw Agent (voice-optimized prompt, short responses)
    ↓ text response
TTS (LiveKit — already configured)
    ↓ audio stream (PCM16)
Speaker / WebSocket to remote device
```

### LiveKit Integration
LiveKit is already the voice interface for OpenClaw (replaced Vapi). The TTS/audio streaming layer is already solved. What's missing is the **input side**: wake word + STT.

### Components to Implement
| Component | Tool | Notes |
|-----------|------|-------|
| Wake word detection | Porcupine (pvporcupine) | Offline, free tier, ~100ms |
| Custom wake words | "hey onemind", "hey om one", "hey legacy" | Configurable |
| Speech-to-text | Groq Whisper (whisper-large-v3) | ~300ms, free tier |
| TTS + streaming | LiveKit (already done) | Already configured |
| Remote voice | WebSocket `/voice` endpoint | Browser/mobile → agent |

### Voice-Optimized Agent Config
- Responses must be short and conversational (not markdown)
- No bullet points or headers in voice responses
- Complete thoughts in <3 sentences where possible
- Tools still execute fully — only the final response format changes

### Latency Breakdown
| Stage | Target | Notes |
|-------|--------|-------|
| Wake word | <100ms | Offline (Porcupine) |
| Audio capture | 1–3s | VAD stops on silence |
| STT | ~300ms | Groq fast inference |
| Agent | 1–2s | Depends on tool calls |
| TTS | Already handled by LiveKit | — |
| **Total** | **3–5s** | End-to-end |

### WebSocket Voice Endpoint
`ws://omos-core/voice` — bidirectional audio stream for browser/mobile clients.
Receives PCM audio chunks → transcribes → agent → streams TTS audio back.

### Use Cases
1. "Hey OneMind, what's the weather?" → spoken response
2. "Hey OM1, move forward" → robot motor control via NATS
3. "Hey OneMind, turn off the lights" → Home Assistant tool call
4. "Hey Legacy, what's on my calendar today?" → Todoist + spoken brief
5. Continuous always-listening mode on omos-core for hands-free operation

### Environment Variables (to add)
```bash
VOICE_ENABLED=true
PORCUPINE_ACCESS_KEY=...     # Free tier: picovoice.ai
GROQ_API_KEY=...             # Free tier: groq.com (STT only)
```

---

## 4. NATS Pulse Visualization (Future)

An OpenClaw canvas screen that renders live NATS message flow as animated particle effects — each NATS message = a spark traveling between agent nodes. Shows system health and agent activity as a living visual in real time.

- Every NATS subject becomes a labeled node
- Message publishes animate as particles between nodes
- Message rate = particle density
- Agent health = node color (green/yellow/red)
- Useful as a command center display wall

---

## Implementation Order

1. **Device Event System** — lowest complexity, highest leverage. NATS is already there. Build `DeviceSource` abstract class + MaixCAM adapter first.
2. **Browser Automation** — medium complexity. Register Playwright tools in OpenClaw. Start with `browser_navigate` and `browser_extract_structured`.
3. **Voice Pipeline** — highest complexity. Porcupine + Groq input side. LiveKit output already done. Wire together in `VoiceGateway` service.
4. **NATS Pulse** — future UI feature, no blockers.

---

## 5. Infrastructure Principles (Future Architecture)

> *Captured from ideation — architecture goals, not yet scoped to a sprint.*

### Air-Gapped Mode
One Mind OS should be operable with zero internet connectivity. All critical agent logic, memory, and tools run locally. Cloud (Bedrock, LiveKit cloud) degrades gracefully to local fallbacks. Use cases: travel, outages, privacy-critical operations, remote homestead.

- Local LLM fallback (Ollama on omos-core) when Bedrock unreachable
- NATS stays local-only in air-gap mode
- Voice pipeline falls back to local Whisper + local TTS
- Codex vault is always local-first (git sync is async)

### Decentralized / Privacy-First / No Vendor Lock-In
One Mind OS must never be hostage to a single vendor, cloud, or platform. Design principles:

- **Portable**: All config, memory, and skills are plain markdown — readable by any LLM, any tool, any future system
- **Hybrid**: Local-first with optional cloud acceleration (not cloud-dependent)
- **Privacy**: Zero telemetry by default. All data stays on Zeus's infrastructure.
- **Multi-provider LLM**: Bedrock today, any provider tomorrow. The `--model` flag is the only coupling point.
- **Blockchain/offline-sync** (long-term): Explore content-addressed storage (IPFS-style) for codex sync between nodes without centralized server dependency

### Persona Thinking Agent
An agent mode where Legacy (or any agent) adopts a specific thinking persona to stress-test ideas — e.g., "think like a skeptic", "think like a first-principles engineer", "think like Zeus in 10 years". Used for decision reviews, architecture critiques, and planning sessions. Implemented as a system prompt overlay in OpenClaw's agent config.

### Bidirectional Cognitive Communications
The communication model between Zeus and agents should be truly bidirectional — not just request/response. Agents proactively surface relevant information, flag anomalies, and initiate conversations when context warrants it. This is partially implemented via HEARTBEAT but should extend to:

- Agent-initiated context push when a NATS event is relevant to Zeus's current focus
- Ambient awareness layer: Legacy monitors signals and interrupts only when signal exceeds a relevance threshold
- Structured back-channel: Zeus can leave async notes for agents to process later (not just real-time voice/chat)
