# omos-bridge Health Dashboard

A lightweight Go binary that subscribes to NATS subjects published by all
omos-bridge bridges and provides a live health view — both as an HTTP JSON
endpoint and as a periodic terminal table.

---

## What it does

| Feature | Detail |
|---|---|
| **NATS health tracking** | Subscribes to `omos.bridge.*.health`; records the last-seen timestamp per bridge name |
| **Error counting** | Subscribes to `omos.bridge.*.errors`; increments an error counter per bridge |
| **Telemetry counting** | Subscribes to `omos.>.telemetry`; counts messages per domain (e.g. `air`, `home`) |
| **Offline detection** | A bridge is **OFFLINE** if its last heartbeat was > 90 s ago (3 missed 30 s beats) |
| **HTTP /health** | JSON endpoint on `DASHBOARD_PORT` (default **8099**) |
| **Terminal table** | Printed every 10 s with a screen-clear; no TUI library required |
| **Graceful shutdown** | Handles `SIGTERM` / `SIGINT`; drains NATS and shuts down HTTP cleanly |
| **NATS reconnect** | `MaxReconnects(-1)` — retries indefinitely |

---

## How to run

### Prerequisites

- Go 1.22+
- A NATS server reachable at `NATS_URL` (default `nats://localhost:4222`)

### Build

```bash
cd /path/to/omos-bridge/internal/dashboard
go build -o dashboard .
```

### Run

```bash
# defaults: NATS nats://localhost:4222, HTTP :8099
./dashboard

# custom NATS + port
NATS_URL=nats://prod-nats:4222 DASHBOARD_PORT=9090 ./dashboard
```

### Docker one-liner

```bash
docker run --rm -e NATS_URL=nats://host.docker.internal:4222 \
  -p 8099:8099 omos-bridge-dashboard
```

---

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `NATS_URL` | `nats://localhost:4222` | NATS server URL |
| `DASHBOARD_PORT` | `8099` | Port the HTTP `/health` endpoint listens on |

---

## HTTP `/health` endpoint

`GET http://localhost:8099/health`

### Response shape

```json
{
  "ts": 1718630400000,
  "bridges": [
    {
      "name": "mavlink-bridge",
      "status": "online",
      "last_heartbeat": 1718630398000,
      "error_count": 0,
      "msg_count": 142
    },
    {
      "name": "tak-bridge",
      "status": "online",
      "last_heartbeat": 1718630395000,
      "error_count": 0,
      "msg_count": 89
    },
    {
      "name": "homeassistant-bridge",
      "status": "offline",
      "last_heartbeat": 1718630280000,
      "error_count": 3,
      "msg_count": 445
    }
  ],
  "domains": {
    "air":  {"msg_count": 142},
    "home": {"msg_count": 445},
    "field": {"msg_count": 89}
  },
  "total_msg_count": 1032
}
```

Field definitions:

| Field | Type | Description |
|---|---|---|
| `ts` | int64 (unix ms) | Snapshot timestamp |
| `bridges[].name` | string | Bridge identifier extracted from NATS subject |
| `bridges[].status` | `"online"` \| `"offline"` | Online if last heartbeat ≤ 90 s ago |
| `bridges[].last_heartbeat` | int64 (unix ms) | Timestamp of last `*.health` message |
| `bridges[].error_count` | int64 | Cumulative errors since start |
| `bridges[].msg_count` | int64 | Cumulative messages (health + errors) seen for this bridge |
| `domains.<name>.msg_count` | int64 | Cumulative telemetry messages for this domain |
| `total_msg_count` | int64 | Sum of all domain message counts |

---

## NATS subject conventions

| Subject pattern | Publisher | Dashboard action |
|---|---|---|
| `omos.bridge.<name>.health` | Each bridge (every 30 s) | Update last-heartbeat timestamp |
| `omos.bridge.<name>.errors` | Each bridge on error | Increment error counter |
| `omos.<domain>.*.telemetry` | Any bridge | Increment domain message counter |

---

## Example terminal output

```
=== omos-bridge Dashboard ===  2026-06-17 14:32:00
BRIDGE                   STATUS    LAST HB        ERRORS  MSGS
-----------------------------------------------------------------
mavlink-bridge           ONLINE    2s ago         0       142
tak-bridge               ONLINE    28s ago        0       89
homeassistant-bridge     OFFLINE   2m ago         3       445

DOMAIN    MSGS
--------------------
air       142
field     89
home      445

Total: 676 messages
```
