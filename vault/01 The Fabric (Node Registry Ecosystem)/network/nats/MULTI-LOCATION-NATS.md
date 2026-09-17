# OneMind OS - Multi-Location NATS Architecture

> How NATS JetStream enables Legacy to be omnipresent across all locations.

---

## Table of Contents

1. [Overview](#overview)
2. [Why NATS?](#why-nats)
3. [Hub and Leaf Architecture](#hub-and-leaf-architecture)
4. [NATS Protocol Support](#nats-protocol-support)
5. [Subject Hierarchy](#subject-hierarchy)
6. [Cloud Hub Configuration](#cloud-hub-configuration)
7. [Edge Leaf Configuration](#edge-leaf-configuration)
8. [Security & Authentication](#security--authentication)
9. [Offline Resilience](#offline-resilience)
10. [Monitoring](#monitoring)
11. [Troubleshooting](#troubleshooting)

---

## Overview

NATS JetStream is the **nervous system** of OneMind OS. It enables:

| Capability | How NATS Enables It |
|------------|---------------------|
| Real-time events | Sub-millisecond pub/sub |
| Multi-location | Hub and leaf topology |
| IoT devices | Native MQTT protocol support |
| Offline resilience | JetStream persistence |
| Browser UI | WebSocket connections |
| Event replay | JetStream consumers |

---

## Why NATS?

### NATS vs Other Options

| Feature | NATS | Kafka | RabbitMQ | Redis Streams |
|---------|------|-------|----------|---------------|
| Latency | 200-400μs | 1-10ms | 1-5ms | 1-2ms |
| MQTT support | ✅ Native | ❌ No | ✅ Plugin | ❌ No |
| Leaf nodes | ✅ Built-in | ❌ No | ❌ No | ❌ No |
| Memory footprint | ~10MB | ~1GB | ~100MB | ~50MB |
| Persistence | JetStream | Yes | Yes | Yes |
| WebSocket | ✅ Native | ❌ No | ✅ Plugin | ❌ No |

### Why NATS Replaces Mosquitto + EMQX

**Before (Complex):**
```
IoT Device → Mosquitto (edge) → EMQX (cloud) → ??? → NATS
                3 services, broken bridge
```

**After (Simple):**
```
IoT Device → NATS Leaf (edge) → NATS Hub (cloud)
                1 service with MQTT support
```

NATS JetStream supports MQTT protocol natively on port 1883. This eliminates the need for separate MQTT brokers.

---

## Hub and Leaf Architecture

```
                              ╔═══════════════════════════════════╗
                              ║        CLOUD NATS HUB             ║
                              ║       (onemindos.dev)             ║
                              ╠═══════════════════════════════════╣
                              ║                                   ║
                              ║  ┌───────────────────────────┐   ║
                              ║  │     NATS JetStream        │   ║
                              ║  │                           │   ║
                              ║  │  :4222  NATS native       │   ║
                              ║  │  :1883  MQTT protocol     │   ║
                              ║  │  :7422  Leaf node accept  │   ║
                              ║  │  :9222  WebSocket (UI)    │   ║
                              ║  │                           │   ║
                              ║  │  JetStream Storage: 10GB  │   ║
                              ║  └───────────────────────────┘   ║
                              ║              │                    ║
                              ╚══════════════│════════════════════╝
                                             │
                    Tailscale VPN + NATS Leaf Protocol (:7422)
                                             │
           ┌─────────────────────────────────┼─────────────────────────────────┐
           │                                 │                                 │
           ▼                                 ▼                                 ▼
╔═══════════════════════╗       ╔═══════════════════════╗       ╔═══════════════════════╗
║    🏠 HOME LEAF       ║       ║    🌾 FARM LEAF       ║       ║   🚗 CAR LEAF         ║
╠═══════════════════════╣       ╠═══════════════════════╣       ╠═══════════════════════╣
║                       ║       ║                       ║       ║                       ║
║  ┌─────────────────┐  ║       ║  ┌─────────────────┐  ║       ║  ┌─────────────────┐  ║
║  │   NATS Leaf     │  ║       ║  │   NATS Leaf     │  ║       ║  │   NATS Leaf     │  ║
║  │                 │  ║       ║  │                 │  ║       ║  │                 │  ║
║  │ :4222 NATS      │  ║       ║  │ :4222 NATS      │  ║       ║  │ :4222 NATS      │  ║
║  │ :1883 MQTT      │  ║       ║  │ :1883 MQTT      │  ║       ║  │ :1883 MQTT      │  ║
║  │                 │  ║       ║  │                 │  ║       ║  │                 │  ║
║  │ Buffer: 256MB   │  ║       ║  │ Buffer: 256MB   │  ║       ║  │ Buffer: 64MB    │  ║
║  └─────────────────┘  ║       ║  └─────────────────┘  ║       ║  └─────────────────┘  ║
║          │            ║       ║          │            ║       ║          │            ║
║          ▼            ║       ║          ▼            ║       ║          ▼            ║
║  ┌─────────────────┐  ║       ║  ┌─────────────────┐  ║       ║  ┌─────────────────┐  ║
║  │ Home Assistant  │  ║       ║  │ Home Assistant  │  ║       ║  │ OBD-II Reader   │  ║
║  │ Frigate NVR     │  ║       ║  │ Frigate NVR     │  ║       ║  │ Dashcam         │  ║
║  │ RemoteAgent     │  ║       ║  │ RemoteAgent     │  ║       ║  │ GPS Tracker     │  ║
║  │ Ollama          │  ║       ║  │ Drone Control   │  ║       ║  │ RemoteAgent     │  ║
║  └─────────────────┘  ║       ║  └─────────────────┘  ║       ║  └─────────────────┘  ║
║                       ║       ║                       ║       ║                       ║
╚═══════════════════════╝       ╚═══════════════════════╝       ╚═══════════════════════╝
```

### How Leaf Nodes Work

1. **Edge NATS starts** and connects to cloud hub on port 7422
2. **Automatic subject bridging** - all messages flow both directions
3. **Local buffering** - if cloud disconnects, messages queue locally
4. **Reconnection** - when cloud is back, queued messages sync
5. **Local clients** - edge services connect to local NATS (fast, low latency)

---

## NATS Protocol Support

NATS JetStream supports multiple protocols on different ports:

### Port Overview

| Port | Protocol | Use Case |
|------|----------|----------|
| **4222** | NATS native | Services, RemoteAgent, high-performance |
| **1883** | MQTT | IoT devices, Home Assistant, legacy devices |
| **7422** | Leaf nodes | Edge location connections |
| **9222** | WebSocket | Browser UI, real-time dashboard |
| **8222** | HTTP | Monitoring, metrics, health checks |

### MQTT Compatibility

NATS MQTT support includes:

- QoS 0, 1, 2
- Retained messages (via JetStream)
- Last Will and Testament
- Session persistence
- Wildcard subscriptions (`+`, `#`)

```yaml
# Home Assistant MQTT configuration
mqtt:
  broker: 192.168.1.100  # NATS Leaf IP
  port: 1883
  discovery: true
```

---

## Subject Hierarchy

### Naming Convention

```
{category}.{location}.{source}.{type}.{detail}
```

### Event Subjects

```
events.>                                    # All events (Legacy subscribes)

events.home.>                               # All home events
events.home.ha.>                            # Home Assistant events
events.home.ha.sensor.motion.living_room    # Motion in living room
events.home.ha.light.bedroom.on             # Bedroom light turned on
events.home.ha.climate.thermostat.changed   # Thermostat changed

events.home.frigate.>                       # Frigate camera events
events.home.frigate.person.front_door       # Person at front door
events.home.frigate.car.driveway            # Car in driveway

events.farm.ha.sensor.soil.moisture         # Farm soil sensor
events.farm.weather.temperature             # Farm weather station

events.car.obd.speed                        # Car speed
events.car.obd.fuel                         # Fuel level
events.car.gps.location                     # GPS position
```

### Command Subjects

```
commands.>                                  # All commands

commands.home.ha.light.living_room.on       # Turn on living room light
commands.home.ha.climate.thermostat.set     # Set thermostat
commands.home.frigate.snapshot.front_door   # Take snapshot

commands.farm.irrigation.zone_a.start       # Start irrigation
commands.farm.drone.patrol                  # Start drone patrol

commands.car.lock                           # Lock car
commands.car.climate.start                  # Start car climate
```

### Query Subjects

```
queries.>                                   # All queries

queries.all.frigate.person.zeus             # Find Zeus (all locations)
queries.home.ha.state.all                   # Get all home states
queries.farm.sensors.current                # Get farm sensor readings
```

### System Subjects

```
system.>                                    # System messages

system.edge.home.heartbeat                  # Edge heartbeat
system.edge.home.status                     # Edge status
system.edge.farm.connected                  # Farm edge connected
system.edge.car.disconnected                # Car edge offline
```

---

## Cloud Hub Configuration

### nats.conf (Cloud)

```conf
# /infra/compose/cloud/nats.conf

# =============================================================================
# OneMind OS - NATS JetStream Cloud Hub Configuration
# =============================================================================

# Server identification
server_name: onemind-cloud-hub

# =============================================================================
# LISTENERS
# =============================================================================

# NATS native protocol
port: 4222

# HTTP monitoring
http_port: 8222

# MQTT protocol (for direct IoT connections)
mqtt {
    port: 1883
    no_auth_user: mqtt_anon
}

# WebSocket (for browser UI)
websocket {
    port: 9222
    no_tls: true  # TLS handled by Cloudflare/Caddy
}

# =============================================================================
# LEAF NODE LISTENER (Edge locations connect here)
# =============================================================================

leafnodes {
    port: 7422

    authorization {
        users = [
            { user: "edge-home", password: "$EDGE_HOME_PASSWORD" }
            { user: "edge-farm", password: "$EDGE_FARM_PASSWORD" }
            { user: "edge-greenhouse", password: "$EDGE_GREENHOUSE_PASSWORD" }
            { user: "edge-car", password: "$EDGE_CAR_PASSWORD" }
            { user: "edge-office", password: "$EDGE_OFFICE_PASSWORD" }
        ]
    }
}

# =============================================================================
# JETSTREAM (Persistence)
# =============================================================================

jetstream {
    store_dir: /data
    max_file_store: 10GB
    max_memory_store: 1GB
}

# =============================================================================
# AUTHORIZATION
# =============================================================================

authorization {
    # Anonymous MQTT users get limited permissions
    mqtt_anon = {
        publish = ["events.>", "sensors.>"]
        subscribe = ["commands.>"]
    }

    # Default user for services
    default_permissions = {
        publish = ">"
        subscribe = ">"
    }
}

# =============================================================================
# LOGGING & LIMITS
# =============================================================================

debug: false
trace: false
logtime: true

max_connections: 10000
max_payload: 1MB
write_deadline: 10s
```

---

## Edge Leaf Configuration

### nats-leaf.conf (Edge)

```conf
# /infra/compose/edge/nats-leaf.conf

# =============================================================================
# OneMind OS - NATS Leaf Node Configuration (Edge)
# =============================================================================

# Server identification (set via environment)
server_name: onemind-edge-$EDGE_LOCATION

# =============================================================================
# LOCAL LISTENERS
# =============================================================================

# NATS native (local services)
port: 4222

# MQTT (local IoT devices)
mqtt {
    port: 1883
}

# HTTP monitoring
http_port: 8222

# =============================================================================
# LEAF NODE CONNECTION TO CLOUD
# =============================================================================

leafnodes {
    remotes [
        {
            # Cloud NATS hub (via Tailscale)
            url: "nats://edge-$EDGE_LOCATION:$EDGE_PASSWORD@$CLOUD_NATS_HOST:7422"

            # Reconnection settings
            reconnect_time_wait: 2s
            reconnect_max_attempts: -1  # Infinite
        }
    ]
}

# =============================================================================
# LOCAL JETSTREAM (Offline Buffer)
# =============================================================================

jetstream {
    store_dir: /data
    max_file_store: 256MB
    max_memory_store: 64MB
}

# =============================================================================
# LOGGING & LIMITS
# =============================================================================

debug: false
trace: false
logtime: true

max_connections: 100
max_payload: 1MB
write_deadline: 10s
```

---

## Security & Authentication

### Network Security

```
┌─────────────────────────────────────────────────────────────┐
│                         INTERNET                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Cloudflare Tunnel (HTTPS only)
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      CLOUD VM                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    Tailscale                         │   │
│  │                  (100.x.x.x)                        │   │
│  └─────────────────────────────────────────────────────┘   │
│                         │                                   │
│                         ▼                                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │               NATS JetStream                         │   │
│  │  :7422 (leaf) ─── Only Tailscale IPs allowed        │   │
│  │  :4222 (nats) ─── Docker network only               │   │
│  │  :1883 (mqtt) ─── Docker network only               │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              ▲
                              │ Tailscale (encrypted)
                              │
┌─────────────────────────────────────────────────────────────┐
│                      EDGE DEVICE                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │               NATS Leaf Node                         │   │
│  │  Connects to 100.102.21.44:7422 (Tailscale)          │   │
│  │  Credentials: edge-home / password                   │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Authentication Methods

| Method | Use Case |
|--------|----------|
| **Username/Password** | Leaf nodes, simple setup |
| **NKey** | Service accounts, no password storage |
| **JWT + NKey** | Multi-tenant, fine-grained permissions |
| **TLS Certificates** | High-security environments |

### Per-Location Credentials

```bash
# Generate password for each edge location
openssl rand -hex 32 > edge-home.password
openssl rand -hex 32 > edge-farm.password
openssl rand -hex 32 > edge-greenhouse.password
```

---

## Offline Resilience

### How Buffering Works

```
         ONLINE                              OFFLINE

┌─────────────┐                      ┌─────────────┐
│  IoT Device │                      │  IoT Device │
└──────┬──────┘                      └──────┬──────┘
       │                                    │
       ▼                                    ▼
┌─────────────┐                      ┌─────────────┐
│  NATS Leaf  │                      │  NATS Leaf  │
│             │──────────────────────│  (buffering)│
│  JetStream  │    Cloud connected   │  JetStream  │
│  (pass-thru)│                      │  256MB queue│
└──────┬──────┘                      └─────────────┘
       │                                    │
       ▼                                    X (disconnected)
┌─────────────┐
│  Cloud Hub  │
└─────────────┘

When cloud reconnects:
┌─────────────┐        sync queued        ┌─────────────┐
│  NATS Leaf  │ ────────────────────────▶ │  Cloud Hub  │
│  (256MB)    │        messages           │             │
└─────────────┘                           └─────────────┘
```

### Buffer Configuration

| Location Type | Buffer Size | Rationale |
|---------------|-------------|-----------|
| Home | 256MB | ~1M events, handles 24h+ offline |
| Farm | 256MB | Sensor data, can be high volume |
| Car | 64MB | Less critical, limited storage |
| Greenhouse | 256MB | Sensor/automation data |

### Critical Events Priority

When buffer fills, older events are dropped. Configure priority:

```conf
# High-priority subjects (never dropped)
jetstream {
    priority_subjects = [
        "events.*.*.security.*",
        "events.*.*.alert.*",
        "commands.>"
    ]
}
```

---

## Monitoring

### NATS Monitoring Endpoints

| Endpoint | Purpose |
|----------|---------|
| `http://localhost:8222/healthz` | Health check |
| `http://localhost:8222/varz` | Server variables |
| `http://localhost:8222/connz` | Connections |
| `http://localhost:8222/routez` | Routes |
| `http://localhost:8222/leafz` | Leaf nodes |
| `http://localhost:8222/jsz` | JetStream stats |

### Key Metrics to Monitor

```bash
# Check connected leaf nodes
curl -s http://localhost:8222/leafz | jq '.leafs[] | {name, account, ip, rtt}'

# Check JetStream storage
curl -s http://localhost:8222/jsz | jq '.memory, .storage'

# Check connections by type
curl -s http://localhost:8222/connz | jq '.connections | group_by(.type) | map({type: .[0].type, count: length})'
```

### Grafana Dashboard Queries

```promql
# Connected leaf nodes
nats_leafnode_connections

# Messages per second
rate(nats_messages_total[1m])

# JetStream storage used
nats_jetstream_storage_bytes

# Consumer lag
nats_consumer_delivered_sequences - nats_consumer_acked_sequences
```

---

## Troubleshooting

### Common Issues

#### Leaf Node Won't Connect

```bash
# Check from edge
nats-server --config /etc/nats/nats.conf -DV

# Look for:
# - "leafnode connection established" = success
# - "authorization violation" = wrong password
# - "connection refused" = firewall/port issue
```

#### Messages Not Flowing

```bash
# Subscribe to all events (debug)
nats sub "events.>"

# Check if messages are being published
nats pub events.test.hello "test message"
```

#### JetStream Full

```bash
# Check storage
nats stream info EVENTS

# Purge old messages
nats stream purge EVENTS --keep 1000000
```

### Debug Commands

```bash
# NATS CLI subscribe (install: nats-io/natscli)
nats sub "events.home.>"

# NATS CLI publish
nats pub events.home.test '{"msg": "hello"}'

# Check leaf status
nats server info --leafz

# Check JetStream streams
nats stream ls
nats stream info EVENTS
```

---

## Quick Reference

### Adding a New Edge Location

1. **Generate credentials:**
   ```bash
   NEW_LOCATION=office
   NEW_PASSWORD=$(openssl rand -hex 32)
   echo "edge-$NEW_LOCATION: $NEW_PASSWORD"
   ```

2. **Add to cloud nats.conf:**
   ```conf
   leafnodes {
       authorization {
           users = [
               # ... existing users ...
               { user: "edge-office", password: "$NEW_PASSWORD" }
           ]
       }
   }
   ```

3. **Deploy edge:**
   ```bash
   curl -sSL https://onemindos.dev/edge-setup.sh | bash -s -- --location office
   ```

4. **Verify connection:**
   ```bash
   # On cloud
   curl -s http://localhost:8222/leafz | jq '.leafs[] | select(.name | contains("office"))'
   ```

### Subject Cheat Sheet

| What you want | Subject |
|---------------|---------|
| All events | `events.>` |
| Home events | `events.home.>` |
| All motion | `events.*.*.sensor.motion.>` |
| All cameras | `events.*.frigate.>` |
| Send command | `commands.{location}.{service}.{action}` |
| Find person | `queries.all.frigate.person.>` |
| Edge status | `system.edge.{location}.>` |

---

## Related Documentation

- [OMNIPRESENT-ARCHITECTURE.md](OMNIPRESENT-ARCHITECTURE.md) - Full system architecture
- [infra/compose/cloud/nats.conf](./infra/compose/cloud/nats.conf) - Cloud NATS config
- [infra/compose/edge/nats-leaf.conf](./infra/compose/edge/nats-leaf.conf) - Edge NATS config
