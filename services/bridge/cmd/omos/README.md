# omos — One Mind Fabric CLI

`omos` is the developer CLI for inspecting and interacting with the **omos-bridge** fabric in real time.  
It connects to the NATS message bus used by omos-bridge and surfaces live bridge status, entity telemetry, and raw message traffic.

---

## Install

```bash
go install github.com/onemindos/omos-bridge/cmd/omos@latest
```

Or build from source:

```bash
git clone https://github.com/onemindos/omos-bridge
cd omos-bridge/cmd/omos
go build -o omos .
```

---

## Environment Variables

| Variable   | Default                  | Description                    |
|------------|--------------------------|--------------------------------|
| `NATS_URL` | `nats://localhost:4222`  | NATS server to connect to      |

```bash
export NATS_URL=nats://nats.internal:4222
```

---

## Commands

### `omos bridge list`

Print the static bridge catalog — no NATS connection required.

```
BRIDGE               PROTOCOL     DOMAIN     STATUS
mavlink              MAVLink v2   air        available
tak                  TAK/CoT      field      available
homeassistant        HA WS API    home       available
webhook              HTTP POST    net        available
mqtt                 MQTT 3.1/5   home       available
ros2                 ROS2/DDS     ground     planned
can                  CAN Bus      ground     planned
modbus               Modbus TCP   infra      planned
opcua                OPC-UA       infra      planned
rtsp                 RTSP         home       planned
```

---

### `omos bridge health`

Subscribe to `omos.bridge.*.health` for **35 seconds** (one full heartbeat cycle), then print bridge status and exit.

```
Collecting bridge health for 35 seconds…

BRIDGE                   STATUS     LAST SEEN
mavlink-bridge           ONLINE     2s ago
tak-bridge               ONLINE     28s ago
```

Bridge health messages are published on `omos.bridge.<bridge-name>.health` as JSON.  
Any message with a `"status"` field is honoured; absent field defaults to `ONLINE`.

---

### `omos fabric watch`

Subscribe to `omos.>` and stream every message as it arrives. Press **Ctrl+C** to exit.

```
Watching omos.> (Ctrl+C to exit)

[14:32:01.123] omos.air.drone-alpha-1.telemetry
{"id":"a1b2c3","ts":1718200321000,"source":"mavlink-bridge","entity":"air.drone-alpha-1","domain":"air","payload":{"msg_type":"position","data":{"lat":30.267,"lon":-97.743,"alt":120.5}}}

[14:32:01.456] omos.field.ANDROID-zeus.telemetry
{"id":"d4e5f6","ts":1718200321456,"source":"tak-bridge","entity":"field.ANDROID-zeus","domain":"field","payload":{"msg_type":"cot","data":{"callsign":"zeus","lat":30.268,"lon":-97.744}}}
```

#### Watch a specific domain

```bash
omos fabric watch --domain air
```

Subscribes to `omos.air.>` — only messages from the `air` domain are shown.

```
Watching omos.air.> (Ctrl+C to exit)

[14:32:01.123] omos.air.drone-alpha-1.telemetry
{"id":"a1b2c3",...}
```

#### Watch an explicit subject

```bash
omos fabric watch omos.home.>
```

Pass any NATS subject pattern directly as a positional argument.

---

### `omos entity list`

Subscribe to `omos.>.telemetry` for **5 seconds**, collect all unique entities seen, then print a summary and exit.

```
Collecting entities for 5 seconds…

ENTITY                         DOMAIN     LAST MSG TYPE     MSGS
air.drone-alpha-1              air        position          142
home.light-living-room         home       sensor            3
field.ANDROID-zeus             field      cot               7
```

Entity names are derived from the NATS subject:  
`omos.<domain>.<entity-id>.telemetry` → entity = `<domain>.<entity-id>`

---

## Message Format

omos-bridge publishes telemetry on subjects of the form:

```
omos.<domain>.<entity-id>.telemetry
```

Payload JSON structure:

```json
{
  "id":      "uuid",
  "ts":      1718200321000,
  "source":  "mavlink-bridge",
  "entity":  "air.drone-alpha-1",
  "domain":  "air",
  "payload": {
    "msg_type": "position",
    "data":     { ... }
  }
}
```

Bridge health messages are published on:

```
omos.bridge.<bridge-name>.health
```

```json
{
  "bridge": "mavlink-bridge",
  "status": "online",
  "ts":     1718200321000
}
```

---

## License

Apache-2.0 — see [LICENSE](../../LICENSE)
