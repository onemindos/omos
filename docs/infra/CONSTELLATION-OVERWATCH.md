# Constellation Overwatch — Integration Runbook
> Integrated: 2026-06-30 | Machine: macOS (Apple Silicon) | Author: Zeus / Legacy

---

## What Is It

[Constellation Overwatch](https://github.com/Constellation-Overwatch/constellation-overwatch) is a
sovereign, self-hosted Common Operating Picture (COP) platform. Single Go binary, no cloud dependency.
It provides a live tactical map, fleet panel, video feeds, and a NATS-backed telemetry bus — the
open-source equivalent of a Palantir Foundry / Lattice OS lightweight COP node.

OneMind OS uses it as the **visual intelligence layer** on top of the NATS spine.

---

## Stack Layout

```
┌─────────────────────────────────────────────────────────────┐
│                    OneMind OS Fabric                        │
│                                                             │
│  ATAK Client ──→ TAK Server :8446                          │
│                       │                                     │
│                  tak.cot.inject                             │
│               (omos-core NATS :4222)                        │
│                       │                                     │
│            cot2constellation bridge                         │
│            (NKey auth, NATS transport)                      │
│                       │                                     │
│         ┌─────────────▼──────────────────┐                 │
│         │   Constellation Overwatch       │                 │
│         │   :8090 (HTTP/WebAuthn)         │                 │
│         │   :4299 (embedded NATS)         │                 │
│         │   SQLite (./data/db/)           │                 │
│         └────────────────────────────────┘                 │
│                       ▲                                     │
│              Pulsar edge agent                              │
│              (fleet sync + NATS streaming)                  │
└─────────────────────────────────────────────────────────────┘
```

---

## Live Service Endpoints

| Service | URL | Notes |
|---------|-----|-------|
| Overwatch UI | `http://localhost:8090/overwatch` | WebAuthn passkey login |
| Overwatch API | `http://localhost:8090/api/v1/` | OpenAPI at `/api/openapi.yaml` |
| API Docs | `http://localhost:8090/api/docs` | Huma UI |
| Embedded NATS | `nats://localhost:4299` | NKey auth required for external clients |
| NATS Dashboard | `http://localhost:5173` | omos-core monitoring (nats-dashboard-1) |
| nui Dashboard | `http://localhost:5175` | NATS web admin (omos-core nui) |

---

## Credentials & Identity

| Item | Value | Location |
|------|-------|----------|
| Admin account | `admin@localhost` | WebAuthn passkey |
| API Key | `c4_live_120c6b027e91d35f5fc3f6df6abcd6976fa9a70512c885cab8b3c8a9978437c8` | Overwatch Admin → API Keys |
| NKey seed file | `~/tools/cot2constellation/overwatch.nk` | For external NATS clients |
| NKey public key | `UBQSTHSKZ6N6I32PNZLEUH2ZSUNK6AYDTOF4ZXQ4MLMH45RNETFH3JG4` | — |
| Org name | `OneMind Homestead` | — |
| Org ID | `9f2d2778-b03c-458f-afe0-ced5bf16179d` | — |
| Pulsar ID | `pulsar-7ab62769` | `~/tools/pulsar/config/c4.json` |

---

## Binary Locations

```
~/tools/constellation-overwatch/overwatch   # v0.1.5 darwin_arm64
~/tools/constellation-overwatch/.env        # runtime config
~/tools/constellation-overwatch/data/       # SQLite DB + state
~/tools/constellation-overwatch/logs/       # overwatch.log

~/tools/cot2constellation/bridge.py         # CoT→Constellation bridge
~/tools/cot2constellation/run.sh            # bridge launcher
~/tools/cot2constellation/overwatch.nk      # NKey seed (chmod 600)
~/tools/cot2constellation/bridge.log        # live bridge log

~/tools/pulsar/pulsar                        # v0.0.9 darwin_arm64
~/tools/pulsar/.env                          # Pulsar config
~/tools/pulsar/config/fleet.yaml            # fleet definition (source of truth)
~/tools/pulsar/config/c4.json               # Pulsar registration state
~/tools/pulsar/pulsar.log                   # live Pulsar log
```

---

## Starting the Stack

### 1. Constellation Overwatch
```bash
cd ~/tools/constellation-overwatch
./overwatch start 2>&1 | tee ./logs/overwatch.log
```

### 2. CoT → Constellation Bridge
```bash
cd ~/tools/cot2constellation
bash run.sh          # logs to bridge.log
```
Subscribes to `tak.cot.inject` on omos-core NATS `:4222`.
Auto-registers new ATAK callsigns as entities in Overwatch.
Publishes telemetry to `constellation.telemetry.{org_id}.{entity_id}` on Overwatch NATS `:4299`.

### 3. Pulsar Edge Agent
```bash
cd ~/tools/pulsar
C4_NATS_KEY="$(cat ~/tools/cot2constellation/overwatch.nk)" \
C4_NATS_URL="nats://localhost:4299" \
./pulsar 2>&1 | tee pulsar.log
```
Pulsar owns the fleet registry — it reconciles `config/fleet.yaml` with Overwatch on every start,
removes stale entities, and connects to Overwatch JetStream for event streaming.

> **Note:** Pulsar is authoritative for entities in `fleet.yaml`. Entities created manually via
> REST API that are NOT in `fleet.yaml` will be removed by Pulsar on next start. To persist an
> entity permanently, add it to `fleet.yaml`.

### 4. NATS Dashboards
```bash
# nats-dashboard-1 (metrics)
cd ~/omos-core/nats/dashboards/nats-dashboard-1
VITE_NATS_MONITOR_URL=http://localhost:8222 npm run dev -- --port 5173 --host

# nui (NATS web admin)
cd ~/omos-core/nats/dashboards/nui/frontend
npm run dev -- --port 5174 --host   # may land on :5175 if 5174 is taken
```

---

## Fleet Management

Fleet state lives in `~/tools/pulsar/config/fleet.yaml`. Pulsar reconciles this on every start.

```yaml
organization:
  name: OneMind Homestead
  type: civilian
  description: omos-core tactical fleet

entities:
  - name: Zeus-TAK
    type: operator        # gcs | operator | uav | fixed_wing | vtol | isr_sensor | sensor ...
    priority: normal

  - name: ADS-B Feeder
    type: isr_sensor
    priority: normal
```

**To add a drone or sensor:** append an entry to `fleet.yaml` and restart Pulsar. It will register
the new entity and assign it an RTSP video slot automatically.

---

## Registered Entities (as of 2026-06-30)

| Name | Type | Entity ID |
|------|------|-----------|
| Zeus-TAK | operator_station | `35c6caf2-67d2-43ca-8809-a3a41fc5efe5` |
| ADS-B Feeder | sensor_platform | `8b641a0b-243f-4064-a894-27cd522b39e2` |

CoT-sourced entities (ATAK callsigns) are auto-registered by the bridge and will appear under
`OneMind Homestead` as they are seen on `tak.cot.inject`. They are **not** in `fleet.yaml` and
will be cleaned up by Pulsar on next start — this is by design (bridge handles transient TAK tracks,
Pulsar handles permanent fleet assets).

---

## Overwatch API Quick Reference

Base: `http://localhost:8090/api/v1`  
Auth header: `X-API-Key: <key>`

```bash
API_KEY="c4_live_120c6b027e91d35f5fc3f6df6abcd6976fa9a70512c885cab8b3c8a9978437c8"
ORG="9f2d2778-b03c-458f-afe0-ced5bf16179d"
BASE="http://localhost:8090/api/v1"

# List organizations
curl -s "$BASE/organizations" -H "X-API-Key: $API_KEY" | jq .

# List entities
curl -s "$BASE/organizations/$ORG/entities" -H "X-API-Key: $API_KEY" | jq .

# Create entity
curl -s -X POST "$BASE/organizations/$ORG/entities" \
  -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" \
  -d '{"name":"Drone-1","entity_type":"aircraft_multirotor","status":"active"}' | jq .

# Update entity position
curl -s -X PATCH "$BASE/organizations/$ORG/entities/<entity_id>" \
  -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" \
  -d '{"position":{"latitude":36.1234,"longitude":-115.9876,"altitude":100.0}}' | jq .

# Health check
curl -s "$BASE/health" | jq .
```

---

## NATS Telemetry Subject Schema

Overwatch NATS `:4299` uses NKey auth. External clients need the seed from `overwatch.nk`.

```
constellation.telemetry.{org_id}.{entity_id}   # position + status updates
constellation.events.{org_id}                  # org-level events (JetStream)
constellation.video.{entity_id}                # video frame metadata
```

The bridge publishes GPS_POSITION messages:
```json
{
  "message_type": "GPS_POSITION",
  "entity_id":    "<uid from CoT>",
  "data": {
    "lat":      363xxx,
    "lon":      -1159xxx,
    "alt":      620000,
    "heading":  270.0,
    "speed":    1.5,
    "callsign": "Zeus-ATAK",
    "cot_type": "a-f-G-U-C",
    "status":   "active"
  },
  "timestamp": "2026-06-30T02:43:00Z"
}
```

---

## Publisher Ecosystem

Additional data sources that can push telemetry to Overwatch:

| Tool | Source | Use Case |
|------|--------|----------|
| `pulsar` | github.com/Constellation-Overwatch/pulsar | All-in-one fleet sync + MAVLink + video |
| `mavlink2constellation` | Constellation-Overwatch org | MAVLink → Overwatch (drones via MAVROS) |
| `ffmpeg2constellation` | Constellation-Overwatch org | FFmpeg video streams → Overwatch |
| `vision2constellation` | Constellation-Overwatch org | OpenCV computer vision → Overwatch |
| `cot2constellation` | `~/tools/cot2constellation/` | TAK CoT → Overwatch (this repo) |

---

## Troubleshooting

### Bridge fails with `nats: encountered error` on :4299
NKey is either missing or the wrong seed. Check:
```bash
cat ~/tools/cot2constellation/overwatch.nk
# Should start with: SUABAPF5SNIN6K...
python3 -c "import nkeys; kp = nkeys.from_seed(open('~/tools/cot2constellation/overwatch.nk','rb').read().strip()); print('OK:', kp.public_key)"
```

### Pulsar removes entities I created via REST
Expected — Pulsar is authoritative. Add them to `config/fleet.yaml` to persist.

### WebAuthn fails with origin error
Ensure `.env` has:
```
OVERWATCH_BASE_URL=http://localhost:8090
OVERWATCH_RPID=localhost
```

### Overwatch data directory
```bash
ls ~/tools/constellation-overwatch/data/db/
# constellation.db — SQLite, safe to inspect with sqlite3
```

### Regenerate bootstrap invite (new admin)
```bash
cd ~/tools/constellation-overwatch
./overwatch invite    # prints a new /invite/<token> URL
```

---

## Integration with omos-core

- **NATS subject**: `tak.cot.inject` — TAK Server NATS Bridge publishes all CoT here
- **Bridge source**: `~/tools/cot2constellation/` (to be moved into `nats/services/` or `integrations/defense/`)
- **Dashboards**: `nats-dashboard-1` (:5173) and `nui` (:5175) monitor the omos-core NATS spine
- **ClickHouse**: Overwatch events can be forwarded to `onemind.metrics` via omos-bridge if needed
