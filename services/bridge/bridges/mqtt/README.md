# mqtt-bridge

> **omos-bridge** component — MQTT ↔ NATS bidirectional  
> Bridges MQTT sensor data into NATS **One Mind envelopes** and routes NATS
> commands back to MQTT devices.

---

## What it does

`mqtt-bridge` maintains two concurrent data flows:

```
╔══════════════════════════════════════════════════════════════════╗
║                    MQTT → NATS  (telemetry)                      ║
║                                                                  ║
║  MQTT Device                                                     ║
║    │ publish  home/living_room/temp  {"temp":21.5}               ║
║    ▼                                                             ║
║  mqtt-bridge                                                     ║
║    │ entity = home.home-living_room-temp                         ║
║    │ subject = omos.home.home-living_room-temp.telemetry         ║
║    ▼                                                             ║
║  NATS → consumers / JetStream / dashboards                       ║
╠══════════════════════════════════════════════════════════════════╣
║                    NATS → MQTT  (commands)                       ║
║                                                                  ║
║  NATS publisher                                                  ║
║    │ publish  omos.home.living-room-light.command  {"on":true}   ║
║    ▼                                                             ║
║  mqtt-bridge                                                     ║
║    │ MQTT topic = omos/home/living-room-light/command            ║
║    ▼                                                             ║
║  MQTT Device receives command                                    ║
╚══════════════════════════════════════════════════════════════════╝
```

---

## Quick start

```bash
# 1. enter bridge directory
cd omos-bridge/bridges/mqtt

# 2. copy env
cp .env.example .env
# edit .env — point to your MQTT broker and NATS server

# 3. install deps & build
go get github.com/nats-io/nats.go@latest
go get github.com/eclipse/paho.mqtt.golang@latest
go build -o mqtt-bridge .

# 4. run
source .env && ./mqtt-bridge
```

---

## Environment variables

| Variable            | Default                  | Description                                               |
|---------------------|--------------------------|-----------------------------------------------------------|
| `NATS_URL`          | `nats://localhost:4222`  | NATS server URL                                           |
| `BRIDGE_NAME`       | `mqtt-bridge`            | Logical bridge name (used in health/error subjects)       |
| `MQTT_HOST`         | `localhost`              | MQTT broker hostname or IP                                |
| `MQTT_PORT`         | `1883`                   | MQTT broker port                                          |
| `MQTT_USERNAME`     | *(empty)*                | MQTT username (optional)                                  |
| `MQTT_PASSWORD`     | *(empty)*                | MQTT password (optional)                                  |
| `MQTT_CLIENT_ID`    | `omos-mqtt-bridge`       | MQTT client identifier                                    |
| `MQTT_TOPIC_FILTER` | `#`                      | MQTT wildcard subscription (`#` = all topics)             |

---

## Topic mapping examples

### MQTT → NATS

| MQTT topic                   | Domain  | Sanitized entity              | NATS subject                                    |
|------------------------------|---------|-------------------------------|-------------------------------------------------|
| `home/living_room/temp`      | `home`  | `home-living_room-temp`       | `omos.home.home-living_room-temp.telemetry`     |
| `air/sensor/co2`             | `air`   | `air-sensor-co2`              | `omos.air.air-sensor-co2.telemetry`             |
| `ground/field1/moisture`     | `ground`| `ground-field1-moisture`      | `omos.ground.ground-field1-moisture.telemetry`  |
| `weather/outdoor/humidity`   | `home`  | `weather-outdoor-humidity`    | `omos.home.weather-outdoor-humidity.telemetry`  |
| `infra/rack1/power`          | `infra` | `infra-rack1-power`           | `omos.infra.infra-rack1-power.telemetry`        |

**Domain detection**: the first segment of the MQTT topic is compared
against the known domains (`home`, `air`, `ground`, `field`, `infra`).
If it matches, that segment becomes the domain. Otherwise `home` is used.

### NATS → MQTT (commands)

| NATS subject                            | MQTT topic                          |
|-----------------------------------------|-------------------------------------|
| `omos.home.living-room-light.command`   | `omos/home/living-room-light/command` |
| `omos.infra.rack1-pdu.command`          | `omos/infra/rack1-pdu/command`      |

---

## Envelope schema (MQTT → NATS)

```json
{
  "entity":    "home.home-living_room-temp",
  "domain":    "home",
  "msg_type":  "sensor",
  "timestamp": 1718000000000,
  "data": {
    "mqtt_topic":     "home/living_room/temp",
    "payload_raw":    "{\"temp\":21.5}",
    "payload_parsed": { "temp": 21.5 },
    "qos":            1,
    "retained":       false
  }
}
```

---

## NATS subjects

| Subject                                | Purpose                         |
|----------------------------------------|---------------------------------|
| `omos.{domain}.{sanitized}.telemetry`  | MQTT sensor data forwarded      |
| `omos.*.*.command`                     | NATS commands forwarded to MQTT |
| `omos.bridge.mqtt-bridge.health`       | Heartbeat every 30 s            |
| `omos.bridge.mqtt-bridge.errors`       | Error reports                   |

---

## Home Assistant MQTT integration

`mqtt-bridge` works seamlessly with Home Assistant's built-in MQTT
integration:

1. Enable MQTT in HA (`Settings → Devices & Services → Add Integration → MQTT`)
2. HA publishes sensor state changes to topics like
   `homeassistant/sensor/living_room_temp/state`
3. Set `MQTT_TOPIC_FILTER=homeassistant/#` to capture all HA topics
4. All HA entities appear in NATS as
   `omos.home.homeassistant-sensor-<name>-state.telemetry`
5. Send commands from NATS → `omos.home.<entity>.command` → MQTT →
   HA automation triggers

> **Tip**: use `MQTT_TOPIC_FILTER=homeassistant/switch/#` to subscribe only
> to switchable devices.

---

## Auto-reconnect behaviour

- **MQTT**: Eclipse Paho `SetAutoReconnect(true)` + `SetConnectRetry(true)`
  — reconnects automatically with 3 s back-off on broker disconnect.
- **NATS**: `MaxReconnects(-1)` — retries forever with 2 s back-off.

Both reconnect silently and resume message flow without process restart.

---

## Roadmap

- [ ] TLS/mTLS MQTT connection
- [ ] MQTT v5 support
- [ ] Per-topic domain override YAML config
- [ ] JetStream publish with ack + dead-letter queue
- [ ] Prometheus `/metrics` endpoint
- [ ] Docker image & Helm chart
- [ ] Home Assistant MQTT Discovery auto-registration
