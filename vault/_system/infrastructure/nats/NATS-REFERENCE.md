# NATS Leaf Node Architecture

> NATS Leaf nodes extend the OneMind event bus to edge locations, bridging MQTT devices into the unified consciousness layer.

---

## Table of Contents

- [Overview](#overview)
- [What is a Leaf Node?](#what-is-a-leaf-node)
- [OneMind NATS Topology](#onemind-nats-topology)
- [Cloud NATS Connection](#cloud-nats-connection)
- [MQTT Bridge Pattern](#mqtt-bridge-pattern)
- [Device Integrations](#device-integrations)
  - [Home Assistant](#home-assistant)
  - [Frigate NVR](#frigate-nvr)
  - [ESP32 Devices](#esp32-devices)
  - [FarmBot](#farmbot)
  - [Robotics](#robotics)
- [Configuration](#configuration)
- [Subject Namespace](#subject-namespace)
- [Message Flow Examples](#message-flow-examples)
- [Security](#security)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)

---

## Overview

OneMind OS uses NATS JetStream as its central nervous system — the event bus connecting all components. Many IoT devices, home automation systems, and robotics platforms speak MQTT natively. **NATS Leaf nodes** with **MQTT bridging** allow these devices to participate in the OneMind consciousness layer without modification.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CLOUD (OneMind Core)                               │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                    NATS Hub Cluster (JetStream)                      │    │
│  │                         onemind.events.*                             │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
└────────────────────────────────────┼─────────────────────────────────────────┘
                                     │ Leaf Connection (TLS)
                                     │
┌────────────────────────────────────┼─────────────────────────────────────────┐
│                            HOME (Edge Location)                              │
│                                    │                                         │
│  ┌─────────────────────────────────┴─────────────────────────────────────┐  │
│  │                    NATS Leaf Node + MQTT Bridge                        │  │
│  │                         nats-server --config leaf.conf                 │  │
│  └─────────────────────────────────┬─────────────────────────────────────┘  │
│                                    │                                         │
│            ┌───────────────────────┼───────────────────────┐                │
│            │                       │                       │                │
│     ┌──────┴──────┐         ┌──────┴──────┐         ┌──────┴──────┐        │
│     │ Home Asst   │         │   Frigate   │         │   ESP32s    │        │
│     │   (MQTT)    │         │   (MQTT)    │         │   (MQTT)    │        │
│     └─────────────┘         └─────────────┘         └─────────────┘        │
│                                                                              │
│            ┌───────────────────────┬───────────────────────┐                │
│            │                       │                       │                │
│     ┌──────┴──────┐         ┌──────┴──────┐         ┌──────┴──────┐        │
│     │   FarmBot   │         │   Robots    │         │   Sensors   │        │
│     │   (MQTT)    │         │   (MQTT)    │         │   (MQTT)    │        │
│     └─────────────┘         └─────────────┘         └─────────────┘        │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## What is a Leaf Node?

A **NATS Leaf Node** is a lightweight NATS server that connects to a hub cluster and:

1. **Extends the subject namespace** — Messages published locally can flow to the hub
2. **Filters subjects** — Only subscribed subjects flow between leaf and hub
3. **Operates independently** — Continues working if hub connection drops
4. **Bridges protocols** — Can expose MQTT, WebSocket, or other protocols locally

### Leaf vs Cluster

| Aspect           | Cluster Node              | Leaf Node              |
| ---------------- | ------------------------- | ---------------------- |
| **Role**         | Full peer in cluster      | Extension/spoke        |
| **State**        | Shares full cluster state | Minimal state          |
| **Subjects**     | All subjects              | Filtered subset        |
| **JetStream**    | Full streams/consumers    | Can access hub streams |
| **Failure mode** | Cluster quorum            | Works in isolation     |
| **Use case**     | Core infrastructure       | Edge/remote locations  |

### Why Leaf Nodes for OneMind?

- **Edge locations** — Home server, farm, workshop all connect to cloud hub
- **Protocol translation** — MQTT devices connect locally, events flow globally
- **Resilience** — Home automation works even if cloud is unreachable
- **Bandwidth efficiency** — Only relevant events traverse the WAN
- **Security** — Edge devices never directly exposed to internet

---

## OneMind NATS Topology

### Hub (Cloud)

```
┌─────────────────────────────────────────────┐
│              NATS Hub Cluster               │
│                                             │
│  nats1.onemind.cloud:4222                  │
│  nats2.onemind.cloud:4222                  │
│  nats3.onemind.cloud:4222                  │
│                                             │
│  JetStream Enabled:                        │
│  - CONSCIOUSNESS stream                    │
│  - EVENTS stream                           │
│  - COMMANDS stream                         │
└─────────────────────────────────────────────┘
```

### Leaf Locations

| Location     | Purpose           | Devices                      |
| ------------ | ----------------- | ---------------------------- |
| **home**     | Primary residence | HA, Frigate, ESP32, Apple TV |
| **farm**     | Agricultural      | FarmBot, sensors, irrigation |
| **workshop** | Making/robotics   | 3D printers, CNCs, robots    |
| **mobile**   | Phone/watch       | ESP32 watch, mobile apps     |

Each location runs a NATS Leaf node with MQTT bridge.

---

## Cloud NATS Connection

The cloud NATS hub is the central nervous system — all leaf nodes connect here. This section covers how to set up and maintain the cloud ↔ edge connection.

### Cloud Infrastructure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              AWS / Cloud Provider                            │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                        VPC: 10.0.0.0/16                                │ │
│  │                                                                        │ │
│  │   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐               │ │
│  │   │   NATS 1    │    │   NATS 2    │    │   NATS 3    │               │ │
│  │   │  (Primary)  │←──►│  (Replica)  │←──►│  (Replica)  │               │ │
│  │   │  10.0.1.10  │    │  10.0.2.10  │    │  10.0.3.10  │               │ │
│  │   └──────┬──────┘    └──────┬──────┘    └──────┬──────┘               │ │
│  │          │ :4222            │ :4222            │ :4222                │ │
│  │          │ :6222 (cluster)  │ :6222            │ :6222                │ │
│  │          │ :7422 (leaf)     │ :7422            │ :7422                │ │
│  │          └──────────────────┴──────────────────┘                      │ │
│  │                             │                                          │ │
│  │                    ┌────────┴────────┐                                │ │
│  │                    │  NLB / Route53  │                                │ │
│  │                    │  nats.onemind.cloud                              │ │
│  │                    └────────┬────────┘                                │ │
│  └─────────────────────────────┼──────────────────────────────────────────┘ │
│                                │                                             │
└────────────────────────────────┼─────────────────────────────────────────────┘
                                 │ :7422 (TLS)
                                 │
         ┌───────────────────────┼───────────────────────┐
         │                       │                       │
    Home Leaf               Farm Leaf              Workshop Leaf
```

### Cloud Hub Ports

| Port     | Protocol  | Purpose                         |
| -------- | --------- | ------------------------------- |
| **4222** | NATS      | Client connections (internal)   |
| **6222** | NATS      | Cluster routing (internal only) |
| **7422** | NATS-Leaf | Leaf node connections (exposed) |
| **8222** | HTTP      | Monitoring endpoint             |
| **443**  | WSS       | WebSocket (browser clients)     |

### Connection URL Format

```bash
# Standard leaf connection
nats-leaf://username:password@nats.onemind.cloud:7422

# With TLS (required for production)
nats-leaf://username:password@nats.onemind.cloud:7422
# + tls.ca_file in config

# Multiple hub nodes for redundancy
nats-leaf://user:pass@nats1.onemind.cloud:7422,nats2.onemind.cloud:7422,nats3.onemind.cloud:7422
```

### TLS Certificate Setup

#### 1. Generate CA and Certificates

```bash
# Create CA (do once)
openssl genrsa -out ca.key 4096
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt \
  -subj "/CN=OneMind NATS CA"

# Generate hub server certificate
openssl genrsa -out hub.key 2048
openssl req -new -key hub.key -out hub.csr \
  -subj "/CN=nats.onemind.cloud"
openssl x509 -req -in hub.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out hub.crt -days 365 \
  -extfile <(echo "subjectAltName=DNS:nats.onemind.cloud,DNS:*.onemind.cloud")

# Generate leaf client certificate (per location)
openssl genrsa -out home-leaf.key 2048
openssl req -new -key home-leaf.key -out home-leaf.csr \
  -subj "/CN=home-leaf"
openssl x509 -req -in home-leaf.csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out home-leaf.crt -days 365
```

#### 2. Hub TLS Configuration

```hcl
# hub.conf - TLS for leaf connections
leafnodes {
  listen: 0.0.0.0:7422

  tls {
    cert_file: /etc/nats/certs/hub.crt
    key_file: /etc/nats/certs/hub.key
    ca_file: /etc/nats/certs/ca.crt
    verify: true              # Require client certs
    timeout: 5s
  }

  authorization {
    users: [
      {
        user: "home-leaf"
        password: "$2a$11$..."  # bcrypt hash
        permissions: {
          publish: ["home.>", "homeassistant.>", "frigate.>", "onemind.watch.>"]
          subscribe: ["onemind.commands.home.>", "onemind.agents.>"]
        }
      }
      {
        user: "farm-leaf"
        password: "$2a$11$..."
        permissions: {
          publish: ["farm.>", "farmbot.>", "esp32.farm.>"]
          subscribe: ["onemind.commands.farm.>"]
        }
      }
    ]
  }
}
```

#### 3. Leaf TLS Configuration

```hcl
# leaf.conf - Connect to cloud hub
leafnodes {
  remotes: [
    {
      url: "nats-leaf://home-leaf:secretpassword@nats.onemind.cloud:7422"

      tls {
        ca_file: /etc/nats/certs/ca.crt
        cert_file: /etc/nats/certs/home-leaf.crt
        key_file: /etc/nats/certs/home-leaf.key
      }

      # Reconnection settings
      reconnect_time_wait: 2s
      max_reconnects: -1        # Unlimited retries
    }
  ]
}
```

### Reconnection & Resilience

#### Automatic Reconnection

NATS Leaf nodes automatically reconnect when the hub connection drops:

```hcl
leafnodes {
  remotes: [
    {
      url: "nats-leaf://home-leaf:pass@nats.onemind.cloud:7422"

      # Reconnection behavior
      reconnect_time_wait: 2s    # Wait between attempts
      max_reconnects: -1         # -1 = unlimited

      # Connection timeouts
      connect_timeout: 5s
      ping_interval: 30s
      ping_max: 3                # Missed pings before disconnect
    }
  ]
}
```

#### Offline Operation

When disconnected from hub, the leaf node continues operating:

| Capability                | Offline Behavior                  |
| ------------------------- | --------------------------------- |
| **Local MQTT**            | ✅ Fully functional               |
| **Local NATS pub/sub**    | ✅ Works for local subscribers    |
| **Local JetStream**       | ✅ If configured with local store |
| **Hub streams**           | ❌ Unavailable until reconnect    |
| **Cross-location events** | ❌ Queued or dropped              |

#### Message Handling During Disconnect

```hcl
# Enable local JetStream for buffering
jetstream {
  store_dir: /data/jetstream
  max_mem: 256M
  max_file: 2G
}

# Hub stream mirroring (optional)
# Pulls from hub when connected, serves locally when not
```

### Hub Cluster Setup

#### Docker Compose (Cloud Hub)

```yaml
# docker-compose.hub.yml
version: "3.8"

services:
  nats1:
    image: nats:2.10-alpine
    container_name: onemind-nats-hub-1
    command: ["--config", "/etc/nats/hub.conf", "--name", "nats1"]
    ports:
      - "4222:4222"
      - "7422:7422"
      - "8222:8222"
    volumes:
      - ./hub.conf:/etc/nats/hub.conf:ro
      - ./certs:/etc/nats/certs:ro
      - nats1_data:/data/jetstream
    networks:
      - nats-cluster
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8222/healthz"]
      interval: 10s
      timeout: 5s
      retries: 5

  nats2:
    image: nats:2.10-alpine
    container_name: onemind-nats-hub-2
    command: ["--config", "/etc/nats/hub.conf", "--name", "nats2"]
    ports:
      - "4223:4222"
      - "7423:7422"
    volumes:
      - ./hub.conf:/etc/nats/hub.conf:ro
      - ./certs:/etc/nats/certs:ro
      - nats2_data:/data/jetstream
    networks:
      - nats-cluster

  nats3:
    image: nats:2.10-alpine
    container_name: onemind-nats-hub-3
    command: ["--config", "/etc/nats/hub.conf", "--name", "nats3"]
    ports:
      - "4224:4222"
      - "7424:7422"
    volumes:
      - ./hub.conf:/etc/nats/hub.conf:ro
      - ./certs:/etc/nats/certs:ro
      - nats3_data:/data/jetstream
    networks:
      - nats-cluster

volumes:
  nats1_data:
  nats2_data:
  nats3_data:

networks:
  nats-cluster:
    driver: bridge
```

#### Hub Cluster Configuration

```hcl
# hub.conf
server_name: $NATS_SERVER_NAME  # Set via env or command line
listen: 0.0.0.0:4222

# JetStream for persistence
jetstream {
  store_dir: /data/jetstream
  max_mem: 1G
  max_file: 50G
  domain: hub
}

# Cluster routing (internal)
cluster {
  name: onemind-hub
  listen: 0.0.0.0:6222
  routes: [
    nats://nats1:6222
    nats://nats2:6222
    nats://nats3:6222
  ]
}

# Leaf node acceptance
leafnodes {
  listen: 0.0.0.0:7422
  # ... TLS and auth config from above
}

# Monitoring
http: 0.0.0.0:8222
```

### Monitoring the Connection

#### Check Leaf Status (from Hub)

```bash
# List connected leaf nodes
curl http://hub:8222/leafz | jq

# Example output
{
  "server_id": "NATS-HUB-1",
  "now": "2026-01-25T10:30:00Z",
  "leafnodes": 3,
  "leafs": [
    {
      "name": "home-leaf",
      "ip": "73.42.xxx.xxx",
      "port": 52431,
      "rtt": "45.2ms",
      "in_msgs": 1523456,
      "out_msgs": 234567,
      "subscriptions": 42
    },
    {
      "name": "farm-leaf",
      "ip": "98.76.xxx.xxx",
      "rtt": "82.1ms",
      "in_msgs": 45678,
      "out_msgs": 12345,
      "subscriptions": 15
    }
  ]
}
```

#### Check Hub Status (from Leaf)

```bash
# View remotes
curl http://localhost:8222/varz | jq '.leafnode'

# Connection status
nats server info --server localhost:4222
```

#### Alerting on Disconnect

```python
# Monitor leaf connection in Legacy agent
import nats

async def monitor_connection():
    nc = await nats.connect("nats://localhost:4222")

    async def disconnected_cb():
        await notify_zeus("⚠️ Cloud connection lost")

    async def reconnected_cb():
        await notify_zeus("✅ Cloud connection restored")

    nc.set_disconnected_cb(disconnected_cb)
    nc.set_reconnected_cb(reconnected_cb)
```

### DNS & Load Balancing

#### Route53 Setup (AWS)

```hcl
# terraform/route53.tf
resource "aws_route53_record" "nats" {
  zone_id = aws_route53_zone.onemind.zone_id
  name    = "nats.onemind.cloud"
  type    = "A"

  alias {
    name                   = aws_lb.nats_nlb.dns_name
    zone_id                = aws_lb.nats_nlb.zone_id
    evaluate_target_health = true
  }
}

resource "aws_lb" "nats_nlb" {
  name               = "onemind-nats-nlb"
  internal           = false
  load_balancer_type = "network"
  subnets            = aws_subnet.public[*].id
}

resource "aws_lb_target_group" "nats_leaf" {
  name     = "nats-leaf-tg"
  port     = 7422
  protocol = "TCP"
  vpc_id   = aws_vpc.main.id

  health_check {
    enabled             = true
    interval            = 10
    port                = 8222
    protocol            = "HTTP"
    path                = "/healthz"
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }
}

resource "aws_lb_listener" "nats_leaf" {
  load_balancer_arn = aws_lb.nats_nlb.arn
  port              = 7422
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.nats_leaf.arn
  }
}
```

#### Multi-Region Setup (Advanced)

```
                    ┌─────────────────────┐
                    │   Global DNS (R53)  │
                    │  nats.onemind.cloud │
                    └─────────┬───────────┘
                              │ Latency-based routing
              ┌───────────────┼───────────────┐
              │               │               │
        ┌─────▼─────┐   ┌─────▼─────┐   ┌─────▼─────┐
        │ US-EAST   │   │ US-WEST   │   │ EU-WEST   │
        │ NATS Hub  │◄─►│ NATS Hub  │◄─►│ NATS Hub  │
        │ Cluster   │   │ Cluster   │   │ Cluster   │
        └─────┬─────┘   └─────┬─────┘   └─────┬─────┘
              │               │               │
         Home Leaf       Farm Leaf       EU Leaf
```

### Quick Connection Test

```bash
# From leaf node, test hub connectivity

# 1. Check DNS resolution
dig nats.onemind.cloud

# 2. Test TCP connectivity
nc -zv nats.onemind.cloud 7422

# 3. Test TLS handshake
openssl s_client -connect nats.onemind.cloud:7422 \
  -CAfile /etc/nats/certs/ca.crt

# 4. Start leaf and check logs
docker logs -f onemind-nats-leaf | grep -i "leaf"

# Expected output:
# [INF] Leafnode connection created for account: $G
# [INF] Leafnode connection to hub established
```

---

## MQTT Bridge Pattern

NATS has native MQTT support — it can act as an MQTT broker while translating messages to/from NATS subjects.

### How It Works

```
MQTT Client                    NATS Leaf                     NATS Hub
────────────                   ─────────                     ────────
                                   │
    ──── PUBLISH ──────────────────►
    topic: home/living_room/temp   │
    payload: {"temp": 72}          │
                                   │
                         ┌─────────┴─────────┐
                         │ Translate to NATS │
                         │ Subject mapping   │
                         └─────────┬─────────┘
                                   │
                                   ──── PUBLISH ───────────────►
                                   subject: home.living_room.temp
                                   payload: {"temp": 72}
```

### MQTT Topic → NATS Subject Mapping

| MQTT Style              | NATS Style              |
| ----------------------- | ----------------------- |
| `home/living_room/temp` | `home.living_room.temp` |
| `frigate/events`        | `frigate.events`        |
| `homeassistant/+/state` | `homeassistant.*.state` |
| `#` (all)               | `>`                     |

NATS automatically translates:

- `/` → `.` (topic separator)
- `+` → `*` (single-level wildcard)
- `#` → `>` (multi-level wildcard)

---

## Device Integrations

### Home Assistant

Home Assistant publishes state changes towards MQTT. We bridge HA events to OneMind's event bus.

#### HA MQTT Configuration

```yaml
# configuration.yaml
mqtt:
  broker: localhost # Local NATS Leaf with MQTT
  port: 1883
  discovery: true
  discovery_prefix: homeassistant
```

#### Subjects Published

| MQTT Topic                  | NATS Subject                | Description                |
| --------------------------- | --------------------------- | -------------------------- |
| `homeassistant/+/+/state`   | `homeassistant.*.*.state`   | Entity state changes       |
| `homeassistant/+/+/command` | `homeassistant.*.*.command` | Commands to entities       |
| `home/+/+`                  | `home.*.*`                  | Custom sensors/automations |

#### Example: Light State Change

```json
// NATS Subject: homeassistant.light.living_room.state
{
  "entity_id": "light.living_room",
  "state": "on",
  "attributes": {
    "brightness": 255,
    "color_temp": 370
  },
  "last_changed": "2026-01-25T10:30:00Z"
}
```

#### Legacy Agent Subscribing

```python
# In Legacy agent - react to home events
async def on_home_event(msg):
    subject = msg.subject
    data = json.loads(msg.data)

    if "light" in subject and data["state"] == "on":
        # Log activity, update context, etc.
        pass

await nc.subscribe("homeassistant.>", cb=on_home_event)
```

---

### Frigate NVR

Frigate is an AI-powered NVR that detects objects (people, cars, animals) in camera feeds. It publishes events over MQTT.

#### Frigate MQTT Configuration

```yaml
# frigate.yml
mqtt:
  host: localhost # Local NATS Leaf
  port: 1883
  topic_prefix: frigate
  client_id: frigate
```

#### Subjects Published

| MQTT Topic                 | NATS Subject               | Description                |
| -------------------------- | -------------------------- | -------------------------- |
| `frigate/events`           | `frigate.events`           | Detection events           |
| `frigate/+/person`         | `frigate.*.person`         | Person detected per camera |
| `frigate/+/car`            | `frigate.*.car`            | Car detected per camera    |
| `frigate/stats`            | `frigate.stats`            | System statistics          |
| `frigate/+/recordings/+/+` | `frigate.*.recordings.*.*` | Recording events           |

#### Example: Person Detection Event

```json
// NATS Subject: frigate.events
{
  "type": "new",
  "before": {},
  "after": {
    "id": "1706012345.123456-abc123",
    "camera": "front_door",
    "label": "person",
    "score": 0.89,
    "box": [0.2, 0.3, 0.5, 0.8],
    "region": [0, 0, 1280, 720],
    "current_zones": ["front_yard"],
    "entered_zones": ["front_yard"],
    "thumbnail": "/api/events/1706012345.123456-abc123/thumbnail.jpg"
  }
}
```

#### OneMind Response

```python
# Legacy agent monitors frigate events
async def on_frigate_event(msg):
    event = json.loads(msg.data)

    if event["after"]["label"] == "person":
        camera = event["after"]["camera"]
        zones = event["after"]["current_zones"]

        # Alert Zeus if person at front door
        if "front_door" in zones:
            await notify_zeus(f"Someone at the front door")

        # Log to consciousness
        await consciousness.log_event(
            type="perception",
            source="frigate",
            data=event
        )

await nc.subscribe("frigate.events", cb=on_frigate_event)
```

---

### ESP32 Devices

ESP32 microcontrollers are the hardware backbone — sensors, actuators, and the OneMind Watch.

#### ESP32 MQTT Setup

```cpp
// Arduino/PlatformIO
#include <WiFi.h>
#include <PubSubClient.h>

const char* mqtt_server = "192.168.1.100";  // NATS Leaf IP
const int mqtt_port = 1883;

WiFiClient espClient;
PubSubClient client(espClient);

void setup() {
  WiFi.begin(ssid, password);
  client.setServer(mqtt_server, mqtt_port);
  client.setCallback(callback);
}

void publishSensor() {
  // Publish temperature reading
  StaticJsonDocument<200> doc;
  doc["device_id"] = "esp32_watch_001";
  doc["temperature"] = readTemp();
  doc["humidity"] = readHumidity();
  doc["battery"] = getBatteryPercent();

  char buffer[256];
  serializeJson(doc, buffer);

  client.publish("onemind/watch/sensors", buffer);
}

void callback(char* topic, byte* payload, unsigned int length) {
  // Handle commands from OneMind
  if (strcmp(topic, "onemind/watch/command") == 0) {
    // Parse and execute command
  }
}
```

#### OneMind Watch Subjects

| Subject                  | Direction   | Description                        |
| ------------------------ | ----------- | ---------------------------------- |
| `onemind.watch.sensors`  | Watch → Hub | Sensor readings (temp, HR, motion) |
| `onemind.watch.location` | Watch → Hub | GPS coordinates                    |
| `onemind.watch.gesture`  | Watch → Hub | Detected gestures                  |
| `onemind.watch.command`  | Hub → Watch | Commands (vibrate, display, etc.)  |
| `onemind.watch.config`   | Hub → Watch | Configuration updates              |

#### Example: Watch Heartbeat

```json
// NATS Subject: onemind.watch.sensors
{
  "device_id": "esp32_watch_001",
  "timestamp": "2026-01-25T10:30:00Z",
  "sensors": {
    "heart_rate": 72,
    "temperature": 36.6,
    "accelerometer": { "x": 0.1, "y": 0.0, "z": 9.8 },
    "battery": 85
  }
}
```

---

### FarmBot

FarmBot is an open-source CNC farming robot. It uses MQTT for control and telemetry.

#### FarmBot MQTT Topics

| MQTT Topic           | NATS Subject             | Description      |
| -------------------- | ------------------------ | ---------------- |
| `bot/+/from_device`  | `farmbot.*.from_device`  | Device telemetry |
| `bot/+/status`       | `farmbot.*.status`       | Bot status       |
| `bot/+/logs`         | `farmbot.*.logs`         | Log messages     |
| `bot/+/from_clients` | `farmbot.*.from_clients` | Commands to bot  |

#### Example: Planting Sequence

```json
// Command to FarmBot
// NATS Subject: farmbot.genesis.from_clients
{
  "kind": "execute",
  "args": {
    "sequence_id": 42
  },
  "body": [
    {
      "kind": "move_absolute",
      "args": {
        "location": { "x": 100, "y": 200, "z": 0 }
      }
    },
    {
      "kind": "write_pin",
      "args": {
        "pin_number": 3,
        "pin_value": 1
      }
    }
  ]
}
```

#### OneMind ↔ FarmBot Integration

```python
# Schedule watering based on soil moisture
async def check_farm_sensors(msg):
    data = json.loads(msg.data)

    if data["sensor"] == "soil_moisture":
        if data["value"] < 30:  # Dry soil
            # Trigger watering sequence
            await nc.publish(
                "farmbot.genesis.from_clients",
                json.dumps({
                    "kind": "execute",
                    "args": {"sequence_id": 10}  # Watering sequence
                })
            )
```

---

### Robotics

General robotics platforms (ROS2, custom controls) can bridge through MQTT.

#### ROS2 ↔ MQTT Bridge

```yaml
# mqtt_bridge config
bridge:
  - topic: /robot/cmd_vel
    msg_type: geometry_msgs/msg/Twist
    mqtt_topic: robot/cmd_vel
  - topic: /robot/odom
    msg_type: nav_msgs/msg/Odometry
    mqtt_topic: robot/odom
```

#### Robotics Subjects

| Subject             | Description                          |
| ------------------- | ------------------------------------ |
| `robot.*.cmd_vel`   | Velocity commands                    |
| `robot.*.odom`      | Odometry data                        |
| `robot.*.sensors.*` | Sensor readings (lidar, camera, IMU) |
| `robot.*.status`    | Robot status/health                  |
| `robot.*.task`      | Task assignments/progress            |

#### Example: Robot Status

```json
// NATS Subject: robot.neo.status
{
  "robot_id": "neo_001",
  "timestamp": "2026-01-25T10:30:00Z",
  "state": "idle",
  "battery": 78,
  "position": { "x": 1.2, "y": 3.4, "theta": 0.5 },
  "current_task": null,
  "errors": []
}
```

---

## Configuration

### Hub Configuration (Cloud)

```hcl
# hub.conf
server_name: onemind-hub-1
listen: 0.0.0.0:4222

jetstream {
  store_dir: /data/jetstream
  max_mem: 1G
  max_file: 10G
}

leafnodes {
  listen: 0.0.0.0:7422

  authorization {
    users: [
      { user: "home-leaf", password: "$2a$...", permissions: {...} }
      { user: "farm-leaf", password: "$2a$...", permissions: {...} }
    ]
  }

  tls {
    cert_file: /etc/nats/certs/server.crt
    key_file: /etc/nats/certs/server.key
    ca_file: /etc/nats/certs/ca.crt
    verify: true
  }
}

# Cluster config for hub nodes
cluster {
  name: onemind-hub
  listen: 0.0.0.0:6222
  routes: [
    nats://nats2.onemind.cloud:6222
    nats://nats3.onemind.cloud:6222
  ]
}
```

### Leaf Configuration (Home)

```hcl
# leaf.conf
server_name: onemind-home-leaf
listen: 0.0.0.0:4222

leafnodes {
  remotes: [
    {
      url: "nats-leaf://home-leaf:password@hub.onemind.cloud:7422"
      tls {
        ca_file: /etc/nats/certs/ca.crt
      }
    }
  ]
}

# MQTT Bridge
mqtt {
  listen: 0.0.0.0:1883

  # No TLS for local network (devices often don't support it)
  # Use network isolation instead
}

# WebSocket for browser clients
websocket {
  listen: 0.0.0.0:8080
  no_tls: true  # Behind reverse proxy
}

# Local JetStream for offline operation
jetstream {
  store_dir: /data/jetstream
  max_mem: 256M
  max_file: 2G
}
```

### Docker Compose (Leaf)

```yaml
# docker-compose.leaf.yml
version: "3.8"

services:
  nats-leaf:
    image: nats:2.10-alpine
    container_name: onemind-nats-leaf
    restart: unless-stopped
    command: ["-c", "/etc/nats/leaf.conf"]
    ports:
      - "4222:4222" # NATS
      - "1883:1883" # MQTT
      - "8080:8080" # WebSocket
      - "8222:8222" # Monitoring
    volumes:
      - ./nats/leaf.conf:/etc/nats/leaf.conf:ro
      - ./nats/certs:/etc/nats/certs:ro
      - nats_data:/data/jetstream
    networks:
      - onemind-local
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8222/healthz"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  nats_data:

networks:
  onemind-local:
    driver: bridge
```

---

## Subject Namespace

### Hierarchy

```
onemind.
├── consciousness.       # Awareness, events, state
│   ├── awareness
│   ├── events
│   └── context
│
├── agents.              # Agent communication
│   ├── legacy.*
│   └── dispatcher.*
│
├── home.                # Home automation (via HA)
│   ├── living_room.*
│   ├── kitchen.*
│   └── ...
│
├── homeassistant.       # HA native topics
│   ├── light.*
│   ├── sensor.*
│   └── ...
│
├── frigate.             # NVR events
│   ├── events
│   └── <camera>.*
│
├── farmbot.             # Agricultural
│   ├── <device>.*
│   └── sensors.*
│
├── robot.               # Robotics
│   ├── <robot_id>.*
│   └── fleet.*
│
├── watch.               # ESP32 watch
│   ├── sensors
│   ├── location
│   └── command
│
└── esp32.               # Other ESP32 devices
    └── <device_id>.*
```

### Naming Conventions

| Pattern                         | Example                    | Use              |
| ------------------------------- | -------------------------- | ---------------- |
| `<system>.<entity>.<attribute>` | `home.kitchen.temperature` | State values     |
| `<system>.<entity>.command`     | `robot.neo.command`        | Commands         |
| `<system>.<entity>.events`      | `frigate.events`           | Event streams    |
| `<system>.<location>.<device>`  | `esp32.workshop.sensor_01` | Device hierarchy |

---

## Message Flow Examples

### Example 1: Doorbell Ring → Watch Vibration

```
1. Frigate detects person at front door
   MQTT: frigate/events
   NATS: frigate.events (on Leaf)

2. Leaf forwards to Hub
   NATS: frigate.events (on Hub)

3. Legacy Agent receives event
   Subscribes: frigate.events
   Action: Send command to watch

4. Publish watch command
   NATS: onemind.watch.command (on Hub)

5. Hub forwards to Home Leaf
   NATS: onemind.watch.command (on Leaf)

6. MQTT translation
   MQTT: onemind/watch/command

7. ESP32 receives, vibrates
```

### Example 2: Voice Command → Light Control

```
1. Zeus says "Turn on kitchen lights" to watch
   Watch publishes voice command
   MQTT: onemind/watch/voice
   NATS: onemind.watch.voice

2. Legacy Agent processes
   Determines intent: light.kitchen.on

3. Publish to Home Assistant
   NATS: homeassistant.light.kitchen.command
   MQTT: homeassistant/light/kitchen/command

4. HA receives, turns on light
   Publishes state update
   MQTT: homeassistant/light/kitchen/state
   NATS: homeassistant.light.kitchen.state

5. Legacy confirms to Zeus
   "Kitchen lights are on"
```

### Example 3: Soil Dry → Auto Watering

```
1. ESP32 soil sensor publishes
   MQTT: esp32/farm/soil_moisture
   NATS: esp32.farm.soil_moisture

2. Leaf forwards to Hub
   NATS: esp32.farm.soil_moisture

3. Legacy Agent evaluates
   Soil moisture < threshold

4. Trigger FarmBot sequence
   NATS: farmbot.genesis.from_clients
   MQTT: bot/genesis/from_clients

5. FarmBot waters plants
   Reports completion
   MQTT: bot/genesis/logs
   NATS: farmbot.genesis.logs

6. Legacy logs to consciousness
   Event recorded for review
```

---

## Security

### Network Isolation

```
┌─────────────────────────────────────────────────────────────┐
│                    HOME NETWORK                              │
│                                                              │
│   ┌──────────────────────────────────────────────────────┐  │
│   │                   IoT VLAN (192.168.10.0/24)         │  │
│   │                                                       │  │
│   │   ESP32s ──────┐                                     │  │
│   │   Cameras ─────┼──────► NATS Leaf ◄────── HA         │  │
│   │   Sensors ─────┘         (MQTT)                      │  │
│   │                            │                          │  │
│   └────────────────────────────┼──────────────────────────┘  │
│                                │                              │
│                    ┌───────────┴───────────┐                 │
│                    │    Firewall Rules     │                 │
│                    │  Only 7422 outbound   │                 │
│                    └───────────┬───────────┘                 │
│                                │                              │
└────────────────────────────────┼──────────────────────────────┘
                                 │
                                 ▼
                         NATS Hub (Cloud)
```

### Authentication Layers

| Layer       | Method                                 |
| ----------- | -------------------------------------- |
| Leaf → Hub  | TLS client certs + user/password       |
| MQTT → Leaf | Username/password (local network only) |
| Devices     | Pre-shared credentials                 |

### Permission Scoping

```hcl
# Hub user for home leaf
{
  user: "home-leaf"
  permissions: {
    publish: {
      allow: [
        "homeassistant.>",
        "frigate.>",
        "onemind.watch.>",
        "esp32.home.>"
      ]
    }
    subscribe: {
      allow: [
        "onemind.watch.command",
        "homeassistant.*.*.command",
        "onemind.agents.legacy.>"
      ]
    }
  }
}
```

---

## Deployment

### Prerequisites

- NATS Hub cluster running in cloud
- TLS certificates generated
- Leaf credentials created on hub
- Local network accessible to IoT devices

### Home Leaf Deployment

```bash
# 1. Create directories
mkdir -p /opt/onemind/nats/{certs,data}

# 2. Copy configuration
cp leaf.conf /opt/onemind/nats/

# 3. Copy certificates
cp ca.crt client.crt client.key /opt/onemind/nats/certs/

# 4. Start with Docker Compose
cd /opt/onemind
docker compose -f docker-compose.leaf.yml up -d

# 5. Verify connection to hub
docker logs onemind-nats-leaf | grep "Leafnode connection"

# 6. Test MQTT bridge
mosquitto_pub -h localhost -t test/hello -m "world"
```

### Verify Connectivity

```bash
# Check leaf status
curl http://localhost:8222/leafz

# Check MQTT clients
curl http://localhost:8222/connz | jq '.connections[] | select(.type == "mqtt")'

# Subscribe to all local events
nats sub ">"
```

---

## Troubleshooting

### Leaf Won't Connect to Hub

```bash
# Check leaf logs
docker logs onemind-nats-leaf

# Verify DNS resolution
nslookup hub.onemind.cloud

# Test TLS connection
openssl s_client -connect hub.onemind.cloud:7422

# Check credentials
nats-server --config leaf.conf --DV  # Debug verbose
```

### MQTT Devices Not Connecting

```bash
# Check MQTT port is exposed
netstat -tlnp | grep 1883

# Test MQTT connection
mosquitto_sub -h localhost -t "#" -v

# Check NATS logs for MQTT activity
docker logs onemind-nats-leaf | grep -i mqtt
```

### Messages Not Flowing to Hub

```bash
# Verify leaf connection
curl http://localhost:8222/leafz | jq

# Check subject permissions
# Ensure leaf user has permission for subject

# Monitor hub side
nats sub ">" --server hub.onemind.cloud:4222
```

### High Latency

```bash
# Check network latency
ping hub.onemind.cloud

# Check JetStream lag
nats stream info EVENTS --server hub.onemind.cloud:4222

# Consider local JetStream for caching
```

---

## References

- [NATS Leaf Nodes Documentation](https://docs.nats.io/running-a-nats-service/configuration/leafnodes)
- [NATS MQTT Support](https://docs.nats.io/running-a-nats-service/configuration/mqtt)
- [Frigate MQTT Documentation](https://docs.frigate.video/integrations/mqtt)
- [Home Assistant MQTT](https://www.home-assistant.io/integrations/mqtt/)
- [FarmBot Web API](https://developer.farm.bot/)

---

_Last updated: January 2026_
