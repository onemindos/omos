# EMQX Integration - IoT Peripheral Nervous System

> Bridging constrained IoT devices to Legacy's consciousness

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    COMPLETE NERVOUS SYSTEM ARCHITECTURE                      │
│                                                                              │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                    HIGH-BANDWIDTH DEVICES                              │  │
│  │                    (LiveKit Direct)                                    │  │
│  │                                                                        │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐  │  │
│  │  │ ESP32    │  │ Brilliant│  │ Unitree  │  │ Desktop/Mobile       │  │  │
│  │  │ Watch    │  │ Labs     │  │ Robots   │  │ Apps                 │  │  │
│  │  │ (voice)  │  │ Halo     │  │ (G1/Go2) │  │                      │  │  │
│  │  └────┬─────┘  └────┬─────┘  └────┬─────┘  └──────────┬───────────┘  │  │
│  │       │             │             │                   │              │  │
│  │       └─────────────┴──────┬──────┴───────────────────┘              │  │
│  │                            │                                          │  │
│  │                            ▼                                          │  │
│  │                     ┌──────────────┐                                 │  │
│  │                     │   LIVEKIT    │  WebRTC, <200ms                 │  │
│  │                     │   CLOUD      │  Voice + Video + Data           │  │
│  │                     └──────┬───────┘                                 │  │
│  │                            │                                          │  │
│  └────────────────────────────┼──────────────────────────────────────────┘  │
│                               │                                              │
│                               ▼                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                         CORE NERVOUS SYSTEM                            │  │
│  │                                                                        │  │
│  │                     ┌──────────────────┐                              │  │
│  │                     │  NATS JETSTREAM  │  <1ms latency                │  │
│  │                     │  (Central Spine) │  Service mesh                │  │
│  │                     └────────┬─────────┘                              │  │
│  │                              │                                         │  │
│  │              ┌───────────────┼───────────────┐                        │  │
│  │              │               │               │                        │  │
│  │              ▼               ▼               ▼                        │  │
│  │  ┌─────────────────┐ ┌─────────────┐ ┌─────────────────┐             │  │
│  │  │ Consciousness   │ │   Agno      │ │  Event          │             │  │
│  │  │ Layer           │ │   (Legacy)  │ │  Archiver       │             │  │
│  │  │ (patterns.py)   │ │             │ │  (TimescaleDB)  │             │  │
│  │  └─────────────────┘ └─────────────┘ └─────────────────┘             │  │
│  │                                                                        │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                               ▲                                              │
│                               │                                              │
│                    ┌──────────┴──────────┐                                  │
│                    │   REDPANDA CONNECT  │  Bridge Layer                    │
│                    │   (MQTT → NATS)     │                                  │
│                    └──────────┬──────────┘                                  │
│                               │                                              │
│  ┌────────────────────────────┼──────────────────────────────────────────┐  │
│  │                            ▼                                          │  │
│  │                     ┌──────────────┐                                  │  │
│  │                     │    EMQX      │  MQTT 5.0                        │  │
│  │                     │  (IoT Hub)   │  10-100ms latency                │  │
│  │                     └──────┬───────┘                                  │  │
│  │                            │                                          │  │
│  │       ┌────────────────────┼────────────────────┐                    │  │
│  │       │                    │                    │                    │  │
│  │       ▼                    ▼                    ▼                    │  │
│  │  ┌──────────┐       ┌──────────┐        ┌──────────┐                │  │
│  │  │ ESPHome  │       │ Zigbee   │        │ LoRa     │                │  │
│  │  │ Sensors  │       │ Bridge   │        │ Nodes    │                │  │
│  │  └──────────┘       └──────────┘        └──────────┘                │  │
│  │                                                                        │  │
│  │                    LOW-POWER / CONSTRAINED DEVICES                     │  │
│  │                    (MQTT via EMQX)                                     │  │
│  │                                                                        │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## When to Use Which Protocol

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      PROTOCOL DECISION MATRIX                                │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        USE LIVEKIT WHEN:                             │    │
│  │                                                                      │    │
│  │  ✓ Voice interaction needed                                         │    │
│  │  ✓ Video/camera streaming                                           │    │
│  │  ✓ Real-time data channels (<200ms)                                 │    │
│  │  ✓ Device has WiFi and ESP32-S3/P4                                  │    │
│  │  ✓ Always-on power available                                        │    │
│  │                                                                      │    │
│  │  DEVICES: ESP32 watch, Brilliant Labs, Unitree robots, phones       │    │
│  │  LATENCY: 50-200ms end-to-end                                       │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                         USE NATS WHEN:                               │    │
│  │                                                                      │    │
│  │  ✓ Service-to-service communication                                 │    │
│  │  ✓ Event streaming between containers                               │    │
│  │  ✓ Sub-millisecond latency needed                                   │    │
│  │  ✓ Cloud/server infrastructure                                      │    │
│  │  ✓ Request-reply patterns                                           │    │
│  │                                                                      │    │
│  │  DEVICES: Agno, Event Processor, Action Gateway, TimescaleDB        │    │
│  │  LATENCY: 200-400µs (microseconds!)                                 │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                      USE EMQX (MQTT) WHEN:                           │    │
│  │                                                                      │    │
│  │  ✓ Battery-powered sensors                                          │    │
│  │  ✓ Sleep/wake duty cycles                                           │    │
│  │  ✓ Unreliable/cellular networks                                     │    │
│  │  ✓ Millions of simple devices                                       │    │
│  │  ✓ Legacy IoT equipment                                             │    │
│  │  ✓ QoS guarantees needed                                            │    │
│  │                                                                      │    │
│  │  DEVICES: ESPHome, Zigbee, bed sensors, environment monitors        │    │
│  │  LATENCY: 10-100ms                                                  │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Docker Compose Addition

Add this to your `docker-compose.yml`:

```yaml
  # ===========================================================================
  # 🌐 IOT HUB: EMQX (MQTT Broker)
  # ===========================================================================
  # MQTT 5.0 broker for IoT devices. Battery sensors, ESPHome, Zigbee bridges.
  # Low-power devices connect here, then bridge to NATS via Redpanda Connect.
  # Dashboard: http://localhost:18083 (admin:public)
  # ===========================================================================

  emqx:
    image: emqx/emqx:5.5
    container_name: onemind-emqx
    restart: unless-stopped
    environment:
      # Cluster name
      EMQX_NAME: onemind-mqtt
      # Dashboard credentials
      EMQX_DASHBOARD__DEFAULT_USERNAME: admin
      EMQX_DASHBOARD__DEFAULT_PASSWORD: ${EMQX_DASHBOARD_PASSWORD:-public}
      # Allow anonymous for local IoT (secure in production!)
      EMQX_ALLOW_ANONYMOUS: "true"
      # Listeners
      EMQX_LISTENER__TCP__EXTERNAL__ACCEPTORS: 8
      EMQX_LISTENER__TCP__EXTERNAL__MAX_CONNECTIONS: 1024000
      # WebSocket for browser clients
      EMQX_LISTENER__WS__EXTERNAL: 8083
    volumes:
      - emqx_data:/opt/emqx/data
      - emqx_log:/opt/emqx/log
      - ./emqx-acl.conf:/opt/emqx/etc/acl.conf:ro
    ports:
      - "1883:1883"     # MQTT TCP
      - "8083:8083"     # MQTT WebSocket
      - "8883:8883"     # MQTT SSL (when configured)
      - "18083:18083"   # Dashboard
    networks:
      - onemind
    healthcheck:
      test: ["CMD", "emqx", "ctl", "status"]
      interval: 30s
      timeout: 10s
      retries: 5
      start_period: 30s
```

Add to volumes section:
```yaml
volumes:
  # ... existing volumes ...
  emqx_data:
  emqx_log:
```

---

## Redpanda Connect Bridge Configuration

Update `connect.yaml` to bridge MQTT to NATS:

```yaml
# ===========================================================================
# Redpanda Connect Configuration
# Bridges: HTTP webhooks → NATS, MQTT (EMQX) → NATS
# ===========================================================================

input:
  broker:
    inputs:
      # Existing HTTP webhook input
      - http_server:
          address: "0.0.0.0:4196"
          path: /webhook
          sync_response:
            status: '${! json("status") }'
        processors:
          - bloblang: |
              root = this
              root.received_at = now()
              root.source = "webhook"

      # NEW: MQTT input from EMQX
      - mqtt:
          urls:
            - tcp://emqx:1883
          topics:
            # Biometric data
            - biometric/+/heart
            - biometric/+/motion
            - biometric/+/temperature
            - biometric/+/sleep
            # Environment sensors
            - environment/+/temperature
            - environment/+/humidity
            - environment/+/air_quality
            # Presence sensors
            - presence/+/motion
            - presence/+/door
            - presence/+/occupancy
            # Home automation
            - home/+/state
            - home/+/command
          client_id: onemind-connect
          clean_session: false
          qos: 1
        processors:
          - bloblang: |
              # Parse MQTT topic into event type
              let topic_parts = meta("mqtt_topic").split("/")
              let category = $topic_parts.index(0)
              let device_id = $topic_parts.index(1)
              let metric = $topic_parts.index(2)

              root.type = "%s.%s".format($category, $metric)
              root.source = "mqtt"
              root.device_id = $device_id
              root.data = this
              root.timestamp = now()
              root.priority = if $category == "biometric" { 2 } else { 1 }

pipeline:
  processors:
    # Route based on source
    - switch:
        - check: this.source == "mqtt"
          processors:
            - bloblang: |
                # Convert MQTT event to NATS subject
                root = this
                meta nats_subject = "events.%s".format(this.type)

output:
  nats_jetstream:
    urls:
      - nats://nats:4222
    subject: '${! meta("nats_subject").or("events.unknown") }'
    stream: ONEMIND
    max_in_flight: 1024
```

---

## MQTT Topic Schema

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           MQTT TOPIC STRUCTURE                               │
│                                                                              │
│  FORMAT: {category}/{device_id}/{metric}                                    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  BIOMETRIC TOPICS                                                    │    │
│  │  ══════════════════                                                  │    │
│  │                                                                      │    │
│  │  biometric/{device}/heart                                           │    │
│  │  └─ {"hr": 72, "hrv": 45, "spo2": 98, "confidence": 95}            │    │
│  │                                                                      │    │
│  │  biometric/{device}/motion                                          │    │
│  │  └─ {"steps": 4523, "activity": "walking", "posture": "standing"}  │    │
│  │                                                                      │    │
│  │  biometric/{device}/temperature                                     │    │
│  │  └─ {"body_temp": 36.8, "skin_temp": 32.5}                         │    │
│  │                                                                      │    │
│  │  biometric/{device}/sleep                                           │    │
│  │  └─ {"stage": "deep", "duration_minutes": 45, "quality": 0.85}     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  ENVIRONMENT TOPICS                                                  │    │
│  │  ══════════════════                                                  │    │
│  │                                                                      │    │
│  │  environment/{device}/temperature                                   │    │
│  │  └─ {"value": 22.5, "unit": "celsius", "location": "office"}       │    │
│  │                                                                      │    │
│  │  environment/{device}/humidity                                      │    │
│  │  └─ {"value": 45, "unit": "percent"}                               │    │
│  │                                                                      │    │
│  │  environment/{device}/air_quality                                   │    │
│  │  └─ {"pm25": 12, "co2": 450, "voc": 0.2}                          │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  PRESENCE TOPICS                                                     │    │
│  │  ═══════════════                                                     │    │
│  │                                                                      │    │
│  │  presence/{device}/motion                                           │    │
│  │  └─ {"detected": true, "location": "hallway"}                      │    │
│  │                                                                      │    │
│  │  presence/{device}/door                                             │    │
│  │  └─ {"state": "open", "location": "front_door"}                    │    │
│  │                                                                      │    │
│  │  presence/{device}/occupancy                                        │    │
│  │  └─ {"occupied": true, "count": 2, "room": "living_room"}          │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  HOME AUTOMATION TOPICS                                              │    │
│  │  ══════════════════════                                              │    │
│  │                                                                      │    │
│  │  home/{device}/state                                                │    │
│  │  └─ {"entity_id": "light.office", "state": "on", "brightness": 80} │    │
│  │                                                                      │    │
│  │  home/{device}/command   (Legacy → device)                          │    │
│  │  └─ {"action": "turn_on", "entity_id": "light.office"}             │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Data Flow Example

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    BIOMETRIC DATA FLOW: BED SENSOR                           │
│                                                                              │
│  ┌─────────────┐                                                            │
│  │  Withings   │  1. Sensor detects heart rate while sleeping              │
│  │  Sleep Mat  │     HR: 58 bpm, HRV: 52ms                                 │
│  └──────┬──────┘                                                            │
│         │                                                                    │
│         │ MQTT Publish                                                       │
│         │ Topic: biometric/withings-001/heart                               │
│         │ Payload: {"hr": 58, "hrv": 52, "spo2": 97}                       │
│         ▼                                                                    │
│  ┌─────────────┐                                                            │
│  │    EMQX     │  2. EMQX receives, applies QoS, stores if offline         │
│  │             │                                                            │
│  └──────┬──────┘                                                            │
│         │                                                                    │
│         │ MQTT Subscribe                                                     │
│         ▼                                                                    │
│  ┌─────────────┐                                                            │
│  │  Redpanda   │  3. Connect transforms MQTT → NATS event                   │
│  │  Connect    │     Subject: events.biometric.heart                        │
│  │             │     Payload: {type, source, device_id, data, timestamp}   │
│  └──────┬──────┘                                                            │
│         │                                                                    │
│         │ NATS JetStream Publish                                            │
│         ▼                                                                    │
│  ┌─────────────┐                                                            │
│  │    NATS     │  4. Event enters the consciousness stream                  │
│  │  JetStream  │                                                            │
│  └──────┬──────┘                                                            │
│         │                                                                    │
│         ├──────────────────┬──────────────────┐                             │
│         │                  │                  │                             │
│         ▼                  ▼                  ▼                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                     │
│  │ patterns.py │    │ archiver.py │    │ Legacy Agent│                     │
│  │             │    │             │    │             │                     │
│  │ Detects:    │    │ Stores in   │    │ Updates     │                     │
│  │ "Zeus is    │    │ TimescaleDB │    │ context:    │                     │
│  │  sleeping"  │    │ for history │    │ "User is    │                     │
│  │             │    │             │    │  asleep"    │                     │
│  └─────────────┘    └─────────────┘    └─────────────┘                     │
│                                                                              │
│  5. Legacy knows: "Zeus is sleeping. HR 58 (normal). Good recovery."       │
│     Actions: Suppress non-urgent notifications, note sleep quality          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## ESPHome Example Configuration

For ESP32-based sensors publishing to EMQX:

```yaml
# esphome/bedroom-sensor.yaml
esphome:
  name: bedroom-sensor
  platform: ESP32
  board: esp32dev

wifi:
  ssid: !secret wifi_ssid
  password: !secret wifi_password

mqtt:
  broker: emqx.local  # Or your EMQX IP
  port: 1883
  topic_prefix: environment/bedroom-001

sensor:
  - platform: dht
    pin: GPIO21
    temperature:
      name: "Bedroom Temperature"
      on_value:
        then:
          - mqtt.publish:
              topic: environment/bedroom-001/temperature
              payload: !lambda |-
                char buf[128];
                snprintf(buf, sizeof(buf),
                  "{\"value\": %.1f, \"unit\": \"celsius\", \"location\": \"bedroom\"}",
                  x);
                return std::string(buf);
    humidity:
      name: "Bedroom Humidity"
      on_value:
        then:
          - mqtt.publish:
              topic: environment/bedroom-001/humidity
              payload: !lambda |-
                char buf[128];
                snprintf(buf, sizeof(buf),
                  "{\"value\": %.0f, \"unit\": \"percent\"}",
                  x);
                return std::string(buf);

binary_sensor:
  - platform: pir
    pin: GPIO22
    name: "Bedroom Motion"
    on_press:
      then:
        - mqtt.publish:
            topic: presence/bedroom-001/motion
            payload: '{"detected": true, "location": "bedroom"}'
    on_release:
      then:
        - mqtt.publish:
            topic: presence/bedroom-001/motion
            payload: '{"detected": false, "location": "bedroom"}'
```

---

## Security Configuration

Create `emqx-acl.conf`:

```erlang
%% EMQX ACL Rules for OneMind OS
%% =============================================================================

%% Allow all local connections (for development)
{allow, {ipaddr, "127.0.0.1"}, all, all}.
{allow, {ipaddr, "172.0.0.0/8"}, all, all}.  %% Docker network

%% Biometric devices - can only publish to their own topics
{allow, {user, "biometric-*"}, publish, ["biometric/+/heart", "biometric/+/motion", "biometric/+/temperature", "biometric/+/sleep"]}.

%% Environment sensors - can only publish sensor data
{allow, {user, "env-*"}, publish, ["environment/+/temperature", "environment/+/humidity", "environment/+/air_quality"]}.

%% Presence sensors
{allow, {user, "presence-*"}, publish, ["presence/+/motion", "presence/+/door", "presence/+/occupancy"]}.

%% Home Assistant bridge - full access to home/ topics
{allow, {user, "homeassistant"}, all, ["home/#"]}.

%% Redpanda Connect - subscribe to everything, publish to nothing
{allow, {user, "onemind-connect"}, subscribe, ["#"]}.

%% Deny everything else
{deny, all, all, all}.
```

---

## EMQX Dashboard Setup

After deployment, access the EMQX dashboard:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        EMQX DASHBOARD SETUP                                  │
│                                                                              │
│  URL: http://localhost:18083                                                │
│  Default: admin / public                                                    │
│                                                                              │
│  1. Authentication Setup                                                     │
│     Dashboard → Access Control → Authentication                             │
│     Add: Built-in Database                                                  │
│     Create users: biometric-watch, env-bedroom, presence-hallway           │
│                                                                              │
│  2. Authorization Setup                                                      │
│     Dashboard → Access Control → Authorization                              │
│     Add: File (acl.conf)                                                    │
│                                                                              │
│  3. Bridge to NATS (Alternative to Redpanda Connect)                        │
│     Dashboard → Data Integration → Data Bridge                              │
│     Type: NATS                                                              │
│     URL: nats://nats:4222                                                   │
│     Topic: Forward all biometric/* → events.biometric.*                    │
│                                                                              │
│  4. Monitoring                                                               │
│     Dashboard → Monitoring → Clients                                        │
│     See all connected IoT devices                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Integration Verification

```bash
# Test MQTT publish (from any device with mosquitto-clients)
mosquitto_pub -h localhost -p 1883 \
  -t "biometric/test-device/heart" \
  -m '{"hr": 72, "hrv": 45, "spo2": 98}'

# Subscribe to see events arrive in NATS
nats sub "events.biometric.>" --server nats://localhost:4222

# Expected output in NATS:
# [events.biometric.heart] {
#   "type": "biometric.heart",
#   "source": "mqtt",
#   "device_id": "test-device",
#   "data": {"hr": 72, "hrv": 45, "spo2": 98},
#   "timestamp": "2024-01-19T10:30:00.000Z",
#   "priority": 2
# }
```

---

## Complete Stack Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      ONEMIND IOT COMPLETE STACK                              │
│                                                                              │
│                                                                              │
│     LAYER 4: CLOUD SERVICES                                                 │
│     ═══════════════════════                                                 │
│                                                                              │
│     ┌──────────────────────────────────────────────────────────────────┐    │
│     │                         AWS / CLOUD                               │    │
│     │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │    │
│     │  │ Bedrock  │  │ LiveKit  │  │ S3       │  │ Route53          │ │    │
│     │  │ (Claude) │  │ Cloud    │  │ Storage  │  │ DNS              │ │    │
│     │  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘ │    │
│     └──────────────────────────────────────────────────────────────────┘    │
│                                    │                                         │
│     LAYER 3: APPLICATION           │                                         │
│     ═════════════════════          │                                         │
│                                    │                                         │
│     ┌──────────────────────────────▼───────────────────────────────────┐    │
│     │                    LIGHTSAIL INSTANCE                             │    │
│     │                                                                   │    │
│     │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐            │    │
│     │  │  Agno   │  │ Action  │  │ Event   │  │   N8N   │            │    │
│     │  │ Legacy  │  │ Gateway │  │Processor│  │         │            │    │
│     │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘            │    │
│     │       └────────────┴────────────┴────────────┘                  │    │
│     │                           │                                      │    │
│     │                    ┌──────▼──────┐                              │    │
│     │                    │ NATS        │                              │    │
│     │                    │ JetStream   │                              │    │
│     │                    └──────┬──────┘                              │    │
│     │                           │                                      │    │
│     │                    ┌──────▼──────┐                              │    │
│     │                    │  Redpanda   │                              │    │
│     │                    │  Connect    │                              │    │
│     │                    └──────┬──────┘                              │    │
│     │                           │                                      │    │
│     │                    ┌──────▼──────┐                              │    │
│     │                    │    EMQX     │  ◄── NEW                     │    │
│     │                    │  MQTT Hub   │                              │    │
│     │                    └──────┬──────┘                              │    │
│     │                           │                                      │    │
│     └───────────────────────────┼──────────────────────────────────────┘    │
│                                 │                                            │
│     LAYER 2: EDGE GATEWAY       │                                            │
│     ═════════════════════       │                                            │
│                                 │                                            │
│     ┌───────────────────────────┼──────────────────────────────────────┐    │
│     │                           │          HOME NETWORK                 │    │
│     │       ┌───────────────────┼───────────────────┐                  │    │
│     │       │                   │                   │                  │    │
│     │       ▼                   ▼                   ▼                  │    │
│     │  ┌─────────┐       ┌─────────────┐     ┌─────────────┐          │    │
│     │  │ Home    │       │ Zigbee2MQTT │     │ ESPHome     │          │    │
│     │  │Assistant│       │ Bridge      │     │ Gateway     │          │    │
│     │  └────┬────┘       └──────┬──────┘     └──────┬──────┘          │    │
│     │       │                   │                   │                  │    │
│     └───────┼───────────────────┼───────────────────┼──────────────────┘    │
│             │                   │                   │                        │
│     LAYER 1: DEVICES            │                   │                        │
│     ════════════════            │                   │                        │
│                                 │                   │                        │
│     ┌───────────────────────────┼───────────────────┼────────────────────┐  │
│     │                           │                   │                    │  │
│     │  ┌─────────┐  ┌─────────┐ │ ┌─────────┐      │ ┌─────────┐       │  │
│     │  │  Hue    │  │ Sonoff  │ │ │ Aqara   │      │ │ ESPHome │       │  │
│     │  │ Lights  │  │ Switch  │ │ │ Sensors │      │ │ Sensors │       │  │
│     │  └─────────┘  └─────────┘ │ └─────────┘      │ └─────────┘       │  │
│     │                           │                   │                    │  │
│     │  ┌─────────┐  ┌─────────────────────────────────────────────┐    │  │
│     │  │ Withings│  │             WEARABLES                        │    │  │
│     │  │ Sleep   │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐      │    │  │
│     │  │ Mat     │  │  │ ESP32   │  │ PineTime│  │ Oura    │      │    │  │
│     │  └─────────┘  │  │ Watch   │  │         │  │ Ring    │      │    │  │
│     │               │  └─────────┘  └─────────┘  └─────────┘      │    │  │
│     │               └─────────────────────────────────────────────┘    │  │
│     │                                                                    │  │
│     │                        PHYSICAL WORLD                              │  │
│     └────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Summary

| Component | Purpose | Port | Protocol |
|-----------|---------|------|----------|
| EMQX | IoT MQTT broker | 1883 (TCP), 8083 (WS), 18083 (Dashboard) | MQTT 5.0 |
| Redpanda Connect | MQTT → NATS bridge | 4195, 4196 | Internal |
| NATS JetStream | Core event bus | 4222 | NATS |

**Data Path:**
```
IoT Device → MQTT → EMQX → Redpanda Connect → NATS → Consciousness Layer → Legacy
```
