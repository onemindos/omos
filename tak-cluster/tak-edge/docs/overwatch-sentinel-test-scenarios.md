# Overwatch Sentinel — ATAK Test Scenarios

Reference for manually triggering simulated threat tracks in ATAK via the overwatch pipeline.

## What it is

`apex-sentinel` is a scenario loop that generates simulated aircraft detections and publishes them to the NATS bus (`det.air.contact.*`). The overwatch pipeline fuses these into `OVERWATCH_STATE` KV, and `om-overwatch-tak` converts each track to CoT XML and injects it into ATAK via `tak.cot.inject`.

Useful for:
- Testing ATAK notification behavior
- Validating the overwatch → TAK pipeline end-to-end
- Demoing hostile track display without live sensor data

## Source location

```
/Users/zeus/omos-core/ignore /defense/apex-sentinel
```

Note: lives in the `ignore/` folder — intentionally excluded from main builds.

## How to start

```bash
# Start sentinel (from its directory)
cd "/Users/zeus/omos-core/ignore /defense/apex-sentinel"
node . >> /private/tmp/overwatch-sentinel.log 2>> /private/tmp/overwatch-sentinel.err &

# Start om-overwatch-tak (injects CoT into ATAK)
/Users/zeus/omos-core/nex/om-overwatch-tak/om-overwatch-tak &
```

Logs:
- `/private/tmp/overwatch-sentinel.log`
- `/private/tmp/overwatch-sentinel.err`

## How to stop

```bash
pkill -f apex-sentinel
pkill -f om-overwatch-tak
```

## Scenarios in the loop

The sentinel cycles through these scenarios continuously:

| Scenario | Description |
|---|---|
| `CHALLENGE_01_PERIMETER` | Single Shahed-136 approach on perimeter |
| `CHALLENGE_01_SWARM` | Multi-track Shahed-136 swarm |
| `CHALLENGE_02_TRAJECTORY` | Fixed-wing trajectory test |
| `CHALLENGE_02_URBAN` | Urban approach track |
| `FULL_PIPELINE` | End-to-end pipeline validation |

## Track characteristics

- Names: `Shahed-136`, `UNKNOWN`
- Types: `aircraft_fixed_wing`, `aircraft_unknown`
- Affiliation: `hostile`
- Source: `sentinel` (simulated), sometimes `vision+sentinel` (fused)
- Coordinates: ~44.1–44.5°N, 26.1–26.6°E (Romania/Ukraine area — test fixture)
- Status progresses: `tentative` → `active` → `lost`

## What you see in ATAK

Each track appears as a shared marker notification:
- Name: `UNKNOWN` or `Shahed-136`
- Remarks: `overwatch:aircraft_unknown status=lost sources=sentinel`
- CoT type: `a-h-A-M-H-Q` (hostile air) or `a-h-A-F-W` (hostile fixed-wing)

## Warning

Running the sentinel alongside a live session will flood ATAK with hundreds of notifications. Only run during dedicated testing — stop with `pkill` commands above when done.
