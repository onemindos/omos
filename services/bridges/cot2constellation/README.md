# cot2constellation

TAK CoT → Constellation Overwatch bridge for OneMind OS.

Subscribes to `tak.cot.inject` on **omos-core NATS** (`:4222`), parses CoT XML,
auto-registers new callsigns as entities, and publishes telemetry to
**Constellation Overwatch NATS** (`:4299`) via NKey authentication.

---

## Quick Start

```bash
# Credentials live outside the repo (never commit secrets)
cp run.sh.example run.sh        # or edit run.sh directly
# Set OW_API_KEY, OW_ORG_ID, OW_NATS_SEED_FILE in run.sh

bash run.sh
tail -f bridge.log
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OMOS_NATS` | `nats://localhost:4222` | omos-core NATS |
| `OW_NATS` | `nats://localhost:4299` | Overwatch embedded NATS |
| `OW_TRANSPORT` | `nats` | `nats` (NKey) or `rest` (HTTP API) |
| `OW_NATS_SEED_FILE` | — | Path to NKey seed file (e.g. `~/tools/cot2constellation/overwatch.nk`) |
| `OW_API_KEY` | — | Overwatch API key (`c4_live_...`) |
| `OW_ORG_ID` | — | Overwatch organization UUID |
| `OW_API_BASE` | `http://localhost:8090` | Overwatch HTTP base URL |
| `COT_SUBJECT` | `tak.cot.inject` | NATS subject to subscribe to |

## Transport Modes

- **`nats`** (default when `OW_NATS_SEED_FILE` is set): Publishes directly to
  `constellation.telemetry.{org_id}.{entity_id}` on Overwatch's embedded NATS.
  Requires an NKey seed generated from the Overwatch Admin → API Keys panel.
- **`rest`**: POSTs position updates via the Overwatch REST API. Slower but
  requires only `OW_API_KEY` + `OW_ORG_ID`.

## Subject Schema

```
IN:   tak.cot.inject          (omos-core NATS :4222)
OUT:  constellation.telemetry.<org_id>.<entity_id>   (Overwatch NATS :4299)
```

## Full Runbook

See [`docs/infra/CONSTELLATION-OVERWATCH.md`](../../docs/infra/CONSTELLATION-OVERWATCH.md).
