# OneMind OS - Omnipresent Architecture

> Legacy's ability to be everywhere, see everything, and act anywhere.

---

## Table of Contents

1. [Overview](#overview)
2. [RemoteAgent vs API](#remoteagent-vs-api)
3. [Architecture Diagram](#architecture-diagram)
4. [Cloud Components](#cloud-components)
5. [Edge Components](#edge-components)
6. [Home Assistant Integration](#home-assistant-integration)
7. [USB Plug & Play Edge Deployment](#usb-plug--play-edge-deployment)
8. [Awareness Modes](#awareness-modes)
9. [Data Flow](#data-flow)

---

## Overview

Omnipresence means Legacy can:

| Capability | Implementation |
|------------|----------------|
| **See** all cameras | Frigate NVR → NATS → Legacy |
| **Hear** all spaces | Microphones → Voice Pipeline → Legacy |
| **Speak** anywhere | Legacy → RemoteAgent → Speakers |
| **Control** any device | Legacy → RemoteAgent → Home Assistant |
| **Find Zeus** | Query all Frigate instances for "person" + face recognition |
| **Work offline** | RemoteAgent → local Ollama when cloud unreachable |

---

## RemoteAgent vs API

There are two ways Legacy can interact with edge locations:

### Option 1: Direct API Calls (Cloud → Edge)

```
┌─────────────┐         HTTPS/API          ┌─────────────┐
│   Legacy    │ ─────────────────────────▶ │ Home Assist │
│   (Cloud)   │                            │   (Edge)    │
└─────────────┘                            └─────────────┘
```

**Pros:**
- Simple - just HTTP calls
- No extra software on edge
- Works with any API-enabled device

**Cons:**
- Requires edge devices to be internet-accessible (security risk)
- Higher latency (round-trip to cloud)
- No offline capability
- No local context awareness
- Each call is independent (no state)

### Option 2: RemoteAgent (Recommended)

```
┌─────────────┐    NATS (events/commands)   ┌─────────────┐
│   Legacy    │ ◀─────────────────────────▶ │ RemoteAgent │
│   (Cloud)   │                             │   (Edge)    │
└─────────────┘                             └─────────────┘
                                                   │
                                                   ▼
                                            ┌─────────────┐
                                            │ Home Assist │
                                            │  Frigate    │
                                            │  Ollama     │
                                            └─────────────┘
```

**Pros:**
- Secure - edge never exposed to internet
- Low latency - commands via NATS (sub-ms)
- Offline capable - falls back to local Ollama
- Local context - knows what's happening at that location
- Stateful - maintains presence and awareness
- Batched operations - can do complex multi-step tasks locally

**Cons:**
- Requires RemoteAgent software on each edge
- More infrastructure to maintain

### Comparison Table

| Factor | Direct API | RemoteAgent |
|--------|-----------|-------------|
| **Latency** | 100-500ms | <10ms (local) |
| **Security** | Edge exposed | Edge hidden |
| **Offline** | No | Yes (Ollama) |
| **Complexity** | Low | Medium |
| **Local awareness** | None | Full |
| **Multi-step tasks** | Multiple round-trips | Single command |
| **State management** | Cloud only | Local + Cloud |

### When to Use Each

| Use Case | Recommendation |
|----------|----------------|
| Quick status check | Either works |
| Control a light | Either works |
| Complex automation | RemoteAgent |
| Security-critical | RemoteAgent |
| Offline locations | RemoteAgent |
| Robot control | RemoteAgent (latency critical) |
| Camera analysis | RemoteAgent (bandwidth) |

### Hybrid Approach (What We Use)

Legacy uses **both**:

1. **RemoteAgent** for all edge locations (primary)
2. **Direct API** as fallback when RemoteAgent is down
3. **Direct API** for cloud services (no edge involved)

```python
# Legacy's decision logic
async def execute_at_location(location: str, action: Action):
    # Try RemoteAgent first (preferred)
    if remote_agent := get_remote_agent(location):
        if await remote_agent.is_healthy():
            return await remote_agent.execute(action)

    # Fallback to direct API
    if api_endpoint := get_direct_api(location):
        return await api_endpoint.call(action)

    raise LocationUnreachable(location)
```

---

## Architecture Diagram

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                              ☁️  CLOUD (Azure VM)                             ║
║                              onemindos.dev                                     ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║   ┌─────────────────────────────────────────────────────────────────────┐    ║
║   │                    🫀 NATS JetStream (THE HEART)                    │    ║
║   │                                                                      │    ║
║   │   • Port 4222: NATS native protocol                                 │    ║
║   │   • Port 1883: MQTT protocol (for IoT devices)                      │    ║
║   │   • Port 7422: Leaf node connections (edge locations)               │    ║
║   │   • Port 9222: WebSocket (browser UI)                               │    ║
║   │                                                                      │    ║
║   └─────────────────────────────────────────────────────────────────────┘    ║
║                    ▲                              ▲                           ║
║                    │ events                       │ commands                  ║
║                    ▼                              ▼                           ║
║   ┌──────────────────────────┐    ┌──────────────────────────┐               ║
║   │     🧠 Legacy (Agno)     │    │   📚 LlamaIndex (RAG)    │               ║
║   │                          │    │                          │               ║
║   │  • Claude Opus 4.5       │    │  • Document processing   │               ║
║   │  • Consciousness layer   │    │  • Vector search         │               ║
║   │  • Awareness modes       │    │  • Google Drive sync     │               ║
║   │  • RemoteAgent control   │    │  • Knowledge base        │               ║
║   └──────────────────────────┘    └──────────────────────────┘               ║
║              │                                                                ║
║              ▼                                                                ║
║   ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐            ║
║   │ PostgreSQL │  │   Redis    │  │   Novu     │  │  Kestra    │            ║
║   │ + pgvector │  │ (hot state)│  │ (notify)   │  │ (workflow) │            ║
║   └────────────┘  └────────────┘  └────────────┘  └────────────┘            ║
║                                                                               ║
║   ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐            ║
║   │ TimescaleDB│  │  Grafana   │  │ Guacamole  │  │ Redpanda   │            ║
║   │ (archive)  │  │ (metrics)  │  │ (remote)   │  │ (webhooks) │            ║
║   └────────────┘  └────────────┘  └────────────┘  └────────────┘            ║
║                                                                               ║
║   ┌─────────────────────────────────────────────────────────────────────┐    ║
║   │                      🖥️  OneMind UI (Flutter)                       │    ║
║   │   Dashboard • Chat • Voice • Cameras • Presence • Settings          │    ║
║   └─────────────────────────────────────────────────────────────────────┘    ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
                                      ▲
                                      │ NATS Leaf Nodes
                                      │ (Tailscale VPN)
                                      ▼
        ┌─────────────────────────────┼─────────────────────────────┐
        │                             │                             │
        ▼                             ▼                             ▼
╔═══════════════════╗       ╔═══════════════════╗       ╔═══════════════════╗
║   🏠 HOME EDGE    ║       ║   🌾 FARM EDGE    ║       ║  🌿 GREENHOUSE    ║
╠═══════════════════╣       ╠═══════════════════╣       ╠═══════════════════╣
║                   ║       ║                   ║       ║                   ║
║ ┌───────────────┐ ║       ║ ┌───────────────┐ ║       ║ ┌───────────────┐ ║
║ │  NATS Leaf    │ ║       ║ │  NATS Leaf    │ ║       ║ │  NATS Leaf    │ ║
║ │  • →Cloud     │ ║       ║ │  • →Cloud     │ ║       ║ │  • →Cloud     │ ║
║ │  • :1883 MQTT │ ║       ║ │  • :1883 MQTT │ ║       ║ │  • :1883 MQTT │ ║
║ └───────────────┘ ║       ║ └───────────────┘ ║       ║ └───────────────┘ ║
║        ▲          ║       ║        ▲          ║       ║        ▲          ║
║        │          ║       ║        │          ║       ║        │          ║
║ ┌──────┴────────┐ ║       ║ ┌──────┴────────┐ ║       ║ ┌──────┴────────┐ ║
║ │ RemoteAgent   │ ║       ║ │ RemoteAgent   │ ║       ║ │ RemoteAgent   │ ║
║ │ (Legacy local)│ ║       ║ │ (Legacy local)│ ║       ║ │ (Legacy local)│ ║
║ └───────────────┘ ║       ║ └───────────────┘ ║       ║ └───────────────┘ ║
║        │          ║       ║        │          ║       ║        │          ║
║        ▼          ║       ║        ▼          ║       ║        ▼          ║
║ ┌───────────────┐ ║       ║ ┌───────────────┐ ║       ║ ┌───────────────┐ ║
║ │Home Assistant │ ║       ║ │Home Assistant │ ║       ║ │Home Assistant │ ║
║ │  (HAOS VM)    │ ║       ║ │  (HAOS VM)    │ ║       ║ │  (HAOS VM)    │ ║
║ └───────────────┘ ║       ║ └───────────────┘ ║       ║ └───────────────┘ ║
║        │          ║       ║        │          ║       ║        │          ║
║ ┌──────┴────────┐ ║       ║ ┌──────┴────────┐ ║       ║ ┌──────┴────────┐ ║
║ │ Frigate NVR   │ ║       ║ │ Frigate NVR   │ ║       ║ │ Frigate NVR   │ ║
║ │ (camera AI)   │ ║       ║ │ (camera AI)   │ ║       ║ │ (camera AI)   │ ║
║ └───────────────┘ ║       ║ └───────────────┘ ║       ║ └───────────────┘ ║
║                   ║       ║                   ║       ║                   ║
║ ┌───────────────┐ ║       ║ ┌───────────────┐ ║       ║ ┌───────────────┐ ║
║ │ Ollama (LLM)  │ ║       ║ │ Drones/Robots │ ║       ║ │   FarmBot     │ ║
║ │ (offline AI)  │ ║       ║ │               │ ║       ║ │ (greenhouse)  │ ║
║ └───────────────┘ ║       ║ └───────────────┘ ║       ║ └───────────────┘ ║
║                   ║       ║                   ║       ║                   ║
║ Zigbee • Z-Wave   ║       ║ LoRa • Sensors    ║       ║ Sensors • Pumps   ║
║ Speakers • Mics   ║       ║ Weather Station   ║       ║ Climate Control   ║
╚═══════════════════╝       ╚═══════════════════╝       ╚═══════════════════╝
```

---

## Cloud Components

### Core (Always Running)

| Service | Port | Purpose |
|---------|------|---------|
| **NATS JetStream** | 4222, 1883, 7422, 9222 | Event backbone + MQTT |
| **Agno (Legacy)** | 8080 | AI brain |
| **PostgreSQL** | 5432 | Cold state + vectors |
| **Redis** | 6379 | Hot state cache |
| **UI** | 3000 | Flutter Web frontend |

### Optional (Enable via profiles)

| Service | Port | Purpose | Profile |
|---------|------|---------|---------|
| **Redpanda Connect** | 4195, 4196 | Webhook → NATS | `webhooks` |
| **Kestra** | 8088 | Workflow orchestration | `workflows` |
| **TimescaleDB** | 5433 | Event archive | `archive` |
| **Novu** | 3003, 4200 | Multi-channel notifications | `notify` |
| **LlamaIndex** | 8082 | Document RAG | `rag` |
| **Grafana** | 3002 | Metrics dashboards | `observe` |
| **VictoriaMetrics** | 8428 | Metrics storage | `observe` |
| **Guacamole** | 8085 | Remote desktop | `manage` |
| **Portainer** | 9000 | Docker UI | `manage` |

---

## Edge Components

### Core (Every Edge)

| Service | Purpose |
|---------|---------|
| **NATS Leaf** | Connects to cloud, provides local MQTT |
| **RemoteAgent** | Legacy's local presence |
| **Tailscale** | Secure VPN to cloud |

### Optional (Per Location)

| Service | Purpose | Use Case |
|---------|---------|----------|
| **Frigate** | Camera AI | Locations with cameras |
| **Ollama** | Offline LLM | Critical locations |
| **FarmBot** | Greenhouse automation | Greenhouse only |

### Home Assistant (Separate)

Home Assistant runs as **HAOS VM** (not Docker) to preserve:
- Add-on support
- Supervisor features
- USB device access (Zigbee, Z-Wave)
- One-click updates

---

## Home Assistant Integration

### Why HAOS VM Instead of Docker?

| Feature | HAOS (VM) | Docker |
|---------|-----------|--------|
| Add-ons | ✅ Full | ❌ None |
| Supervisor | ✅ Yes | ❌ No |
| USB (Zigbee/Z-Wave) | ✅ Easy | ⚠️ Complex |
| Backups | ✅ Built-in | ⚠️ Manual |
| Updates | ✅ One-click | ⚠️ Manual |

### Connection Flow

```
┌─────────────────────────────────────────────────────────┐
│                  Edge Hardware                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────────┐    ┌─────────────────────┐    │
│  │   Home Assistant    │    │      Docker         │    │
│  │   (HAOS VM)         │    │                     │    │
│  │                     │    │  ┌───────────────┐  │    │
│  │  • Zigbee2MQTT      │───▶│  │  NATS Leaf    │  │    │
│  │  • Z-Wave JS        │    │  │  (:1883 MQTT) │  │    │
│  │  • ESPHome          │    │  └───────────────┘  │    │
│  │  • Automations      │    │         │          │    │
│  │                     │    │         ▼          │    │
│  │  MQTT Integration   │    │  ┌───────────────┐  │    │
│  │  → nats-leaf:1883   │    │  │ RemoteAgent   │  │    │
│  └─────────────────────┘    │  └───────────────┘  │    │
│                             │         │          │    │
│                             │         ▼          │    │
│                             │  ┌───────────────┐  │    │
│                             │  │   Frigate     │  │    │
│                             │  └───────────────┘  │    │
│                             └─────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

### Home Assistant MQTT Configuration

```yaml
# configuration.yaml
mqtt:
  broker: 192.168.1.100  # NATS Leaf IP (Docker host)
  port: 1883
  discovery: true
  discovery_prefix: homeassistant
```

### Event Flow

1. **Device state changes** → Home Assistant detects
2. **HA publishes** → MQTT to NATS Leaf (:1883)
3. **NATS Leaf bridges** → Cloud NATS Hub
4. **Legacy receives** → Processes event
5. **Legacy decides** → Sends command (if needed)
6. **Command flows** → Cloud → NATS Leaf → RemoteAgent
7. **RemoteAgent calls** → Home Assistant API
8. **Device actuates** → Light turns on, etc.

---

## USB Plug & Play Edge Deployment

### Goal

Flash a USB drive, plug into any mini PC/Jetson, boot, done.

### What's on the USB

```
onemind-edge-usb/
├── docker-compose.yml      # Edge services
├── .env                    # Location config (set once)
├── nats-leaf.conf          # Auto-configured
├── setup.sh                # First-boot script
└── README.md               # Quick start guide
```

### First Boot Process

```bash
# 1. Plug USB into edge device
# 2. Boot from USB (or copy to SSD)
# 3. Run setup script

./setup.sh --location home

# This will:
# - Install Docker (if needed)
# - Configure NATS leaf for this location
# - Set up Tailscale (auto-join)
# - Start edge services
# - Register with cloud
```

### Setup Script Flow

```bash
#!/bin/bash
# setup.sh

LOCATION=${1:-home}

echo "🏠 Setting up OneMind Edge: $LOCATION"

# 1. Install Docker
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com | sh
fi

# 2. Configure environment
cat > .env << EOF
EDGE_LOCATION=$LOCATION
CLOUD_NATS_HOST=100.102.21.44
EDGE_PASSWORD=\$(openssl rand -hex 16)
TIMEZONE=America/New_York
EOF

# 3. Generate NATS config
envsubst < nats-leaf.conf.template > nats-leaf.conf

# 4. Join Tailscale network
if [ -n "$TS_AUTHKEY" ]; then
    tailscale up --authkey=$TS_AUTHKEY --hostname=onemind-$LOCATION
fi

# 5. Start services
docker compose up -d

# 6. Register with cloud
curl -X POST "https://onemindos.dev/api/edge/register" \
    -H "Content-Type: application/json" \
    -d "{\"location\": \"$LOCATION\", \"ip\": \"$(tailscale ip -4)\"}"

echo "✅ Edge setup complete!"
echo "📍 Location: $LOCATION"
echo "🔗 Tailscale IP: $(tailscale ip -4)"
```

### Remote Installation (Alternative)

```bash
# From any edge device with internet
curl -sSL https://onemindos.dev/edge-setup.sh | bash -s -- --location greenhouse
```

---

## Awareness Modes

Legacy operates in different awareness modes based on context:

| Mode | Heartbeat | Event Processing | Use Case |
|------|-----------|------------------|----------|
| **DORMANT** | 60s | Emergencies only | Night, vacation |
| **AWARE** | 20s | Passive monitoring | Normal background |
| **PRESENT** | 5s | Active engagement | Working with Zeus |
| **OMNIPRESENT** | 1s | Everything, everywhere | Full presence mode |

### Mode Transitions

```
                    ┌──────────────┐
                    │   DORMANT    │
                    │   (60s)      │
                    └──────┬───────┘
                           │ wake word / emergency
                           ▼
                    ┌──────────────┐
                    │    AWARE     │
                    │   (20s)      │
                    └──────┬───────┘
                           │ user interaction
                           ▼
                    ┌──────────────┐
                    │   PRESENT    │
                    │   (5s)       │
                    └──────┬───────┘
                           │ explicit request
                           ▼
                    ┌──────────────┐
                    │ OMNIPRESENT  │
                    │   (1s)       │
                    └──────────────┘
```

### What Changes Per Mode

| Mode | Cameras | Home Events | Notifications | Voice |
|------|---------|-------------|---------------|-------|
| DORMANT | Motion only | Security only | Critical only | Off |
| AWARE | All events | Filtered | Important | Wake word |
| PRESENT | All + analysis | All | All | Active |
| OMNIPRESENT | Real-time | Real-time | Real-time | Always on |

---

## Data Flow

### Event Ingestion (Edge → Cloud)

```
1. Sensor detects motion
   └─▶ Home Assistant
       └─▶ MQTT publish: homeassistant/sensor/motion/living_room
           └─▶ NATS Leaf (:1883)
               └─▶ Cloud NATS Hub
                   └─▶ Subject: events.home.ha.sensor.motion.living_room
                       └─▶ Legacy receives
                           └─▶ Consciousness processor evaluates
                               └─▶ Action decision
```

### Command Execution (Cloud → Edge)

```
1. Legacy decides: "Turn on living room lights"
   └─▶ Publish: commands.home.ha.light.living_room.on
       └─▶ Cloud NATS Hub
           └─▶ NATS Leaf (home)
               └─▶ RemoteAgent receives
                   └─▶ Calls Home Assistant API
                       └─▶ Light turns on
                           └─▶ HA publishes state change
                               └─▶ Confirmation flows back to Legacy
```

### Find Zeus (Multi-Location Query)

```
1. Legacy: "Where is Zeus?"
   └─▶ Publish: queries.all.frigate.person.zeus
       └─▶ All NATS Leaves receive
           └─▶ Each RemoteAgent queries local Frigate
               └─▶ Responses: events.{location}.frigate.person.detected
                   └─▶ Legacy aggregates
                       └─▶ "Zeus detected at home.living_room (85% confidence)"
```

---

## Related Documentation

- [MULTI-LOCATION-NATS.md](MULTI-LOCATION-NATS.md) - NATS hub/leaf architecture
- [docs/architecture/infrastructure-visual.md](./docs/architecture/infrastructure-visual.md) - Full infrastructure diagram
- [docs/architecture/remote-agent-use-cases.md](./docs/architecture/remote-agent-use-cases.md) - RemoteAgent examples
- [docs/legacy/CONSCIOUSNESS-PROCESSOR.md](./docs/legacy/CONSCIOUSNESS-PROCESSOR.md) - Event filtering
