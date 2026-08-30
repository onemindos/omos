# tak-bridge — TAK/CoT ↔ NATS Bridge

Part of the **One Mind OS (omos)** bridge ecosystem.

## What It Does

`tak-bridge` translates between the [Cursor-on-Target (CoT)](https://www.mitre.org/capabilities/cybersecurity/overview/cybersecurity-blog/cursor-on-target) XML protocol used by ATAK, WinTAK, and OpenTAK Server (OTS) and the omos NATS message bus.

**Inbound (ATAK → NATS):**
- Listens on a UDP port for CoT XML events from ATAK/OTS clients.
- Parses each `<event>` into a typed `CotPayload` struct.
- Wraps the payload in the standard One Mind envelope and publishes to `omos.field.<uid>.telemetry`.

**Outbound (NATS → ATAK):**
- Subscribes to `omos.field.*.command` on NATS.
- Serialises command payloads as CoT `<event>` XML and forwards them to the configured TAK server via UDP.

**Operations:**
- Heartbeat published every 30 s to `omos.bridge.tak-bridge.health`.
- Errors reported to `omos.bridge.tak-bridge.errors`.
- Graceful shutdown on `SIGTERM`/`SIGINT` with NATS drain.

---

## Quick Start

```bash
# 1. Clone / enter the bridge directory
cd omos-bridge/bridges/tak

# 2. Copy and edit environment
cp .env.example .env

# 3. Build
go build -o tak-bridge .

# 4. Run (with a local NATS server already running)
source .env && ./tak-bridge
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `NATS_URL` | `nats://localhost:4222` | NATS server connection URL |
| `BRIDGE_NAME` | `tak-bridge` | Bridge identifier used in log lines and NATS subjects |
| `TAK_LISTEN_PORT` | `8087` | UDP port to listen on for incoming CoT XML |
| `TAK_SERVER_HOST` | `localhost` | Hostname / IP of the TAK server to forward commands to |
| `TAK_SERVER_PORT` | `8087` | UDP port of the TAK server |

---

## NATS Subjects

| Subject | Direction | Description |
|---|---|---|
| `omos.field.<uid>.telemetry` | Publish | One Mind envelope containing `CotPayload` for each received CoT event |
| `omos.field.*.command` | Subscribe | Commands received here are wrapped in CoT XML and forwarded to the TAK server |
| `omos.bridge.tak-bridge.health` | Publish | Heartbeat JSON every 30 s: `{"source","ts","status":"ok"}` |
| `omos.bridge.tak-bridge.errors` | Publish | Error JSON: `{"source","ts","error":"…"}` |

### One Mind Envelope Shape

```json
{
  "id": "1718000000000000000",
  "ts": 1718000000000,
  "source": "tak-bridge",
  "entity": "field.ANDROID-abc",
  "domain": "field",
  "payload": {
    "msg_type": "cot",
    "data": {
      "uid": "ANDROID-abc",
      "cot_type": "a-f-G-U-C",
      "how": "m-g",
      "callsign": "ZEUS",
      "lat": 30.267,
      "lon": -97.743,
      "hae": 150.0,
      "ce": 9.9,
      "le": 9.9,
      "stale_at": 1751500800000,
      "detail_raw": "<detail><contact callsign=\"ZEUS\"/></detail>"
    }
  }
}
```

---

## Watch Commands

```bash
# Stream all field telemetry from ATAK units
nats sub 'omos.field.*.telemetry'

# Stream telemetry from a specific unit
nats sub 'omos.field.ANDROID-abc.telemetry'

# Monitor bridge health
nats sub 'omos.bridge.tak-bridge.health'

# Monitor bridge errors
nats sub 'omos.bridge.tak-bridge.errors'

# Send a test command to all field units
nats pub 'omos.field.ANDROID-abc.command' '{"action":"ping"}'

# Pretty-print incoming telemetry (requires jq)
nats sub 'omos.field.*.telemetry' | while read -r line; do echo "$line" | jq .; done
```

---

## Architecture

```
ATAK/WinTAK/OTS
      │  CoT XML (UDP)
      ▼
 tak-bridge (this binary)
      │  JSON One Mind Envelope
      ├──▶  NATS: omos.field.<uid>.telemetry
      │
      │  NATS: omos.field.*.command
      ◀──┤
      │  CoT XML (UDP)
      ▼
 TAK Server / OTS
```

---

## Roadmap

| Version | Feature |
|---|---|
| **v0.1** | UDP inbound CoT → NATS, NATS command → UDP outbound *(current)* |
| **v0.2** | TCP multicast / streaming socket support for high-density TAK nets |
| **v0.3** | Full bidirectional NATS→TAK command relay with ACK envelope |
| **v0.4** | `orbit.go` batch publish — aggregate CoT events and publish as bulk telemetry bundles |

---

## License

Apache-2.0 — see [`SPDX-License-Identifier`](../../../LICENSE) header in source.
