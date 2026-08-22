# homeassistant-bridge

A standard-compliant OMOS bridge that connects a **Home Assistant** instance to the
One Mind OMOS message bus over **NATS JetStream**.  It is the replacement for the
earlier `ha-bridge` prototype and conforms to the same message-envelope contract
used by all other omos-bridge adapters (MAVLink, TAK, …).

---

## What it does

| Direction | How |
|-----------|-----|
| HA → NATS | Subscribes to `state_changed` events via the HA WebSocket API.  For every change it publishes two NATS messages: a **telemetry** envelope and an **event** envelope. |
| NATS → HA | Subscribes to `omos.home.*.command` on NATS and forwards each message to HA as a `call_service` WebSocket request. |

### Features
- Auto-reconnect on HA WebSocket disconnect (exponential backoff, max 30 s)
- NATS connection with unlimited reconnects (`MaxReconnects(-1)`)
- Heartbeat every 30 s → `omos.bridge.homeassistant-bridge.health`
- Error reporting → `omos.bridge.homeassistant-bridge.errors`
- Graceful shutdown on `SIGTERM` / `SIGINT`
- All configuration from environment variables (12-factor)

---

## Setup

### Prerequisites
- Go 1.22+
- A running NATS server
- Home Assistant with a **Long-Lived Access Token**
  (Profile → Security → Long-Lived Access Tokens → Create Token)

### Build

```bash
cd bridges/homeassistant
go mod tidy
go build -o homeassistant-bridge .
```

### Configure

Copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
```

### Run

```bash
export $(cat .env | xargs)
./homeassistant-bridge
```

---

## Environment variables

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `HA_TOKEN` | — | **Yes** | Home Assistant long-lived access token |
| `HA_HOST` | `localhost` | No | Hostname / IP of the HA instance |
| `HA_PORT` | `8123` | No | HTTP/WS port of the HA instance |
| `NATS_URL` | `nats://localhost:4222` | No | NATS server URL |
| `BRIDGE_NAME` | `homeassistant-bridge` | No | Bridge identifier (used in subjects & envelopes) |

---

## NATS subjects

### Telemetry (HA → NATS)

**Subject:** `omos.home.<entity-id-with-dashes>.telemetry`

_Example:_ `omos.home.light-living-room.telemetry`

```json
{
  "id": "homeassistant-bridge-1717612345000000000",
  "ts": 1717612345000,
  "source": "homeassistant-bridge",
  "entity": "light-living-room",
  "domain": "home",
  "msg_type": "sensor",
  "payload": {
    "entity_id": "light.living_room",
    "friendly_name": "Living Room",
    "state": "on",
    "unit": "",
    "device_class": "",
    "attributes": { "brightness": 128, "color_temp": 370 }
  }
}
```

### Event (HA → NATS)

**Subject:** `omos.home.<entity-id-with-dashes>.event`

```json
{
  "id": "homeassistant-bridge-1717612345000000001",
  "ts": 1717612345000,
  "source": "homeassistant-bridge",
  "entity": "light-living-room",
  "domain": "home",
  "msg_type": "event",
  "payload": {
    "event_type": "state_changed",
    "entity_id": "light.living_room",
    "old_state": "off",
    "new_state": "on",
    "context_id": "01HXYZ...",
    "occurred_at": 1717612345000
  }
}
```

### Command (NATS → HA)

**Subject:** `omos.home.*.command`

Send a JSON body that matches the `commandPayload` shape:

```json
{
  "domain": "light",
  "service": "turn_on",
  "service_data": { "entity_id": "light.living_room", "brightness": 200 }
}
```

Or wrap it in the standard `BridgeMessage` envelope — both are accepted.

### Health & Errors

| Subject | Description |
|---------|-------------|
| `omos.bridge.homeassistant-bridge.health` | Heartbeat every 30 s |
| `omos.bridge.homeassistant-bridge.errors` | Error reports |

---

## Watch commands

```bash
# Stream all HA telemetry
nats sub 'omos.home.*.telemetry'

# Stream all HA events
nats sub 'omos.home.*.event'

# Watch a single entity
nats sub 'omos.home.light-living-room.*'

# Send a command
nats pub omos.home.light-living-room.command \
  '{"domain":"light","service":"turn_on","service_data":{"entity_id":"light.living_room"}}'

# Monitor bridge health
nats sub 'omos.bridge.homeassistant-bridge.health'

# Monitor bridge errors
nats sub 'omos.bridge.homeassistant-bridge.errors'
```

---

## Replacing the legacy `ha-bridge`

This bridge supersedes the earlier `ha-bridge` prototype.  Key differences:

| Aspect | Legacy ha-bridge | homeassistant-bridge |
|--------|-----------------|----------------------|
| Message envelope | Custom / ad-hoc | Standard `BridgeMessage` (same as mavlink-bridge, tak-bridge) |
| Subject scheme | Varies | `omos.home.<entity>.telemetry` / `.event` / `.command` |
| Reconnect logic | None | Exponential backoff (max 30 s) |
| Command forwarding | Not implemented | `omos.home.*.command` → `call_service` |
| Health reporting | None | Every 30 s to `.health` subject |
| Config | Partial | 100 % env-var driven |

To migrate: stop the old `ha-bridge` process, set the same `HA_TOKEN`/`HA_HOST`
env vars, and start `homeassistant-bridge`.  No HA configuration changes needed.

---

## Roadmap

- [ ] TLS/HTTPS support for HA connections (`HA_TLS=true`)
- [ ] Subscribe to additional HA event types (e.g. `call_service`, `automation_triggered`)
- [ ] Filter entities by domain or glob pattern (`HA_ENTITY_FILTER`)
- [ ] NATS JetStream persistence for telemetry streams
- [ ] Prometheus metrics endpoint
- [ ] Docker image + Helm chart
- [ ] Area / floor metadata enrichment from HA entity registry
- [ ] Two-way attribute sync (NATS → HA `set_state`)
