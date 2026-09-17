# OneMind OS — Recon 03: NATS / Bridge / Telemetry Overlap
*Generated 2026-07-06 by Legacy (direct). Read-only. Reference repos vs LIVE workloads.*

## 1. Duplication table (reference repo vs live workload)

| Reference project | Duplicates LIVE workload | Verdict |
|---|---|---|
| `tools/pulsar/` (MAVLink+RTSP+YOLO edge) | om-bridge-mavlink + overwatch-vision | **SUPERSEDED** — pulsar's jobs are split across live workloads |
| GitHub `mavlink2constellation` | om-bridge-mavlink | **SUPERSEDED** |
| GitHub `aero-arc-relay2constellation` | om-bridge-mavlink | **SUPERSEDED** |
| GitHub `ffmpeg2constellation` | om-overwatch-video / om-streaming-bridge | **SUPERSEDED** |
| GitHub `gocv2constellation` / `vision2constellation` | overwatch-vision | **SUPERSEDED** (patterns worth keeping) |
| `tools/cot2constellation/bridge.py` | om-bridge-constellation | **COMPLEMENTARY** — Python bridge still running; Go version is om-bridge-constellation. Pick one |
| `tools/constellation-overwatch/` (Go binary) | om-bridge-constellation consumes its entity model | **REFERENCE** — entity schema source of truth |
| GitHub `pulsar` (Constellation org) | om-bridge-mavlink | **SUPERSEDED** |
| GitHub `toolbelt` (IaC) | — | **UNIQUE** (IaC templates, reference only) |

## 2. Bridge inventory (protocol in → out)

| Bridge | In | Out | Status |
|---|---|---|---|
| om-bridge-mavlink | MAVLink UDP/serial | `tel.drone.*` | LIVE, keep (Mac) |
| om-bridge-constellation | CoT XML | Constellation entities | LIVE, keep |
| om-bridge-homeassistant | HA LAN events | `iot.*` | LIVE, keep (Mac) |
| om-streaming-bridge | RTSP/stream | video subjects | LIVE |
| om-http-gateway | HTTP | NATS req/reply | LIVE |
| takserver-bridge (NatsBridgePlugin) | NATS | TAK CoT | LIVE (your plugin) |
| cot-tap (DOKS) | TAK CoT | NATS | LIVE |
| cot2constellation/bridge.py | CoT | Constellation | LIVE — redundant with om-bridge-constellation |

## 3. Salvageable patterns (extract before deletion)

- **vision2constellation** — Python YOLO/Moondream edge detection → entity publish pattern. Good reference for overwatch-vision improvements.
- **constellation-overwatch** — the `entity_id` primitive schema + KV state aggregation + vendor-agnostic fleet ontology. This is the canonical entity model; document it.
- **cot2constellation** — the CoT type → entity_type mapping table (a-f-* → operator_station, a-.-A → aircraft_multirotor, etc.). Extract this mapping into shared schema.
- **aero-arc-relay2constellation** — pure-Go MAVLink fan-out pattern (fault tolerance). Reference for hardening om-bridge-mavlink.
- **omos-core/hive/orbit.go, jwt-auth-builder.go, callout.go** — reusable Go NATS libraries. Keep in hive/.
- **omos-core/integrations/nats/autonats** — codegen for NATS bindings. Keep as tool.

## 4. Delete-safe list (100% superseded by live workloads)

- `tools/pulsar/` — jobs covered by om-bridge-mavlink + overwatch-vision
- GitHub clones already extracted: mavlink2constellation, aero-arc-relay2constellation, ffmpeg2constellation (patterns documented → then remove local clones)
- `tools/dashboards/nats-dashboard-1/` — already running via launchd from its install location

## 5. Keep-as-reference (don't delete, extract to docs)

- `tools/constellation-overwatch/` → entity schema docs → `omos-core/docs/patterns/entity-model.md`
- `tools/cot2constellation/` → CoT mapping table → shared schema, THEN retire the .py (Go bridge wins)
- `omos-core/hive/` → keep entirely (NATS libraries)
- `omos-core/integrations/nats/` → keep autonats + UNS-Design-Studio as tools

## Verdict

The Constellation ecosystem repos are **mostly superseded** by your live om-bridge-* workloads — you already extracted the value into Go NEX workloads. Action: document the entity model + CoT mapping (the two genuinely valuable patterns), retire the duplicate Python bridge, then delete the superseded clones. Keep `hive/` and `integrations/nats/` libraries.
