# Hermes

## Key Capabilities

- **Edge device access** — Native terminal tool with `watch_patterns` + SSH gives you direct, stateless access to any machine on your Tailscale mesh. No persistent connections to maintain, no reconnection logic needed.
- **HomeAssistant is a first-class platform** — Not just MCP, but `homeassistant` in `platform_toolsets`. HA events can trigger Hermes turns directly.
- **Plugin development** — Python plugins with no build step. You can write a NATS bridge, a custom device fleet plugin, or any custom process in hours, not days.

This is a Python script running as a Hermes cron/plugin. The ENTIRE logic chain — sensor reading, API call, memory lookup, decision, actuation, notification — is one Python process.

**Given:** drones + robots + wearables + agriculture + MediaPipe + light shows + LiveKit + NATS + Meshtastic + DimOS + OM1 + growing food + fleet ops

---

## HERMES (Brain)

```
HERMES (Brain)
├── Cron: "Every dawn, check all sensor readings + weather forecast → decide spray schedule"
├── Plugin: frame-bridge (Brilliant Labs glasses → direct Python SDK)
├── Plugin: nats-fleet (subscribe to all edge telemetry, publish commands)
├── Plugin: mediapipe-gesture (gesture → robot commands)
├── Plugin: dmx-shows (light automation)
├── MCP: DimOS (drone/robot runtime commands)
├── MCP: HomeAssistant (IoT devices)
├── Platform: Telegram (mobile alerts + commands)
├── Platform: HomeAssistant (bidirectional events)
├── SSH/Tailscale: direct server management
└── Webhook → OpenClaw: "display this on Zeus's glasses/phone/Mac"
```

## OPENCLAW (Interface/Display)

```
OPENCLAW (Interface/Display)
├── Talk mode: "Legacy, show me barn camera" → Hermes processes → OpenClaw renders
├── Mac node: desktop dashboard, quick commands
├── iOS node: mobile interface
├── Canvas: live fleet map, robot cameras, farm status
└── Voice: ambient presence (the voice you talk to)
```

---

## Full Stack Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│  BRAIN: Hermes Agent (Mac, always-on)                           │
│  ├── Decision engine: "soil dry + forecast hot → irrigate now"  │
│  ├── NATS subscriber: listens to all edge telemetry             │
│  ├── NATS publisher: sends commands to devices                  │
│  ├── Cron: scheduled checks, daily farm reports                 │
│  ├── Vision: analyze camera frames (YOLO/OpenCV via terminal)   │
│  ├── Drone missions: publishes MAVLink waypoints via NATS       │
│  └── Telegram/Discord: alerts you, takes voice commands         │
├─────────────────────────────────────────────────────────────────┤
│  MESSAGE BUS: NATS (central server on VPS or home NUC)          │
│  ├── Subjects: farm.barn.sensors.*, farm.barn.camera.*          │
│  ├── Subjects: office.docker.*, home.security.*                 │
│  ├── JetStream: persistent streams for history/replay           │
│  └── Leaf nodes at each location (auto-reconnect)              │
├─────────────────────────────────────────────────────────────────┤
│  DEVICE HUB: HomeAssistant (home NUC)                           │
│  ├── Zigbee/Z-Wave: switches, locks, HVAC                      │
│  ├── ESPHome: custom sensors → MQTT → NATS bridge              │
│  ├── Frigate: NVR + object detection → NATS publish            │
│  └── Automations: local failsafe (irrigate if agent offline)   │
├─────────────────────────────────────────────────────────────────┤
│  EDGE: Per-Location                                             │
│                                                                 │
│  BARN:                                                          │
│  ├── RPi 4: NATS leaf + Python sensor scripts                  │
│  ├── ESP32s: soil moisture, temperature, valve control          │
│  ├── IP cameras → Frigate (RPi or local NVR)                   │
│  └── Irrigation controller (GPIO or ESPHome relay)             │
│                                                                 │
│  OFFICE:                                                        │
│  ├── Linux server: NATS leaf + Docker services                 │
│  ├── publishes: docker events, system health, build status     │
│  └── GPU box: ML inference, video processing                   │
│                                                                 │
│  DRONE/ROBOTICS:                                                │
│  ├── PX4/ArduPilot: autonomous flight (local)                  │
│  ├── Companion computer (RPi): Python + NATS + camera          │
│  ├── publishes: telemetry, GPS, camera frames                  │
│  └── subscribes: mission waypoints, RTL commands               │
├─────────────────────────────────────────────────────────────────┤
│  NETWORK: Tailscale Mesh                                        │
│  └── All devices directly addressable by hostname              │
├─────────────────────────────────────────────────────────────────┤
│  VOICE/PERSONAL: OpenClaw (secondary)                           │
│  ├── Talk mode: "Hey, what's the barn temperature?"            │
│  ├── Mac/iOS node: canvas dashboards, A2UI                     │
│  └── Receives display commands from Hermes via webhook         │
└─────────────────────────────────────────────────────────────────┘
```
