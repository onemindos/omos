# OneMind OS — Recon 05: NEX Workload Topology (THE BIG ONE)
*Generated 2026-07-06 by Legacy (direct). Read-only. What runs where + cloud migration plan.*

## Full workload table (LIVE on Mac zeus-mac NEX node)

| Workload | Dir | Lang | Subscribes | Publishes | Stateful? | Recommendation |
|---|---|---|---|---|:--:|---|
| **om-unity (ORP)** | `ai-agents/orp/target/release/om-unity` (binary) + `nex/om-unity/` (source) | Rust | external protocols (AIS/ADS-B/MAVLink) | (not wired to NATS yet) → should be `tel.*` | ⚠️ ingest-stateful | **MOVE TO DOKS** (always-on fusion) + **DEDUP source** |
| om-overwatch-state | `nex/om-overwatch-state/` | Go | `det.air.*`, `tel.drone.*` | `ent.air.{id}.state` | no | **MOVE TO DOKS** (stateless aggregator) |
| om-overwatch-tak | `nex/om-overwatch-tak/` | Go | `ent.air.{id}.state` | `tak.cot.out` (CoT→TAK) | no | **MOVE TO DOKS** |
| om-overwatch-interceptor | `nex/om-overwatch-interceptor/` | Go | detections | intercept commands | no | **MOVE TO DOKS** |
| om-overwatch-relay | `nex/om-overwatch-relay/` | Go | relay in | relay out | no | **DUPLICATE?** (also listed nex/) — verify |
| om-overwatch-api | `nex/om-overwatch-api/` | Go | — | REST :port | no | **MOVE TO DOKS** |
| om-overwatch-video | `nex/om-overwatch-video/` | Go+MediaMTX | RTSP feeds | HLS/WebRTC | ⚠️ media | **STAY/HYBRID** (MediaMTX near cameras) |
| om-bridge-mavlink | `nats/nex/om-bridge-mavlink/` | Go | MAVLink UDP/serial | `tel.drone.*` | ✅ hardware | **STAY ON MAC** (local MAVLink) |
| om-bridge-constellation | `nats/nex/om-bridge-constellation/` | Go | CoT | Constellation entities | no | **MOVE TO DOKS** |
| om-bridge-homeassistant | `nex/om-bridge-homeassistant/` | Go | HA events (LAN) | `iot.*` | ✅ LAN | **STAY ON MAC** (Home Assistant LAN) |
| om-streaming-bridge | `nats/nex/om-streaming-bridge/` | Go | stream in | video subjects | ⚠️ media | **HYBRID** |
| om-http-gateway | `nats/nex/om-http-gateway/` | Go | HTTP in | NATS req/reply | no | **MOVE TO DOKS** |
| om-callout | `nats/nex/om-callout/` | Go | `callout.*` | notifications out | no | **MOVE TO DOKS** |
| om-nex-ui | `nats/nex/om-nex-ui/` | Go | workload status | UI | no | **MOVE TO DOKS** |
| om-task-nex | `nex/om-task/` | Go | `task.*` | task results | no | **MOVE TO DOKS** |
| resgate | `nex/resgate/` AND `nats/nex/resgate/` | Go | RES protocol | WS to clients | no | **DUPLICATE** — pick one, move to DOKS |
| om_presence_v2 | `nats/services/connect/` | Python | presence pings | `presence.*` | no | **MOVE TO DOKS** |
| nats-surveyor | `nex/nats-surveyor/` | Go | NATS monitoring | metrics | no | **DOKS** (observability) |
| apex-sentinel | `defense/apex-sentinel/` | script | sensor feeds | alerts | ⚠️ | verify hardware dep |
| overwatch-vision | `defense/om-overwatch-vision/` | script | camera frames | detections | ✅ camera | **STAY ON MAC** (local vision) |

## 2. Duplication findings (CRITICAL)

1. **ORP/om-unity source is in TWO places:**
   - `nex/om-unity/` — full source tree (crates, frontend, sdk, Cargo.lock) — NOT built
   - `ai-agents/orp/` — has the RUNNING binary (`target/release/om-unity`)
   - **Action:** these are the same project (ORP = om-unity). Consolidate to ONE canonical location. The running binary is from `ai-agents/orp`. Decide canonical home, delete the other, keep one source of truth.
2. **resgate defined twice:** `nex/resgate/` and `nats/nex/resgate/`. Only one should exist.
3. **om-overwatch-relay** appears in listings under `nex/` — verify no second copy.
4. **Two NEX root dirs:** `nex/` and `nats/nex/` both hold workloads. This split is confusing — workloads scattered across two trees. **Consolidate to one NEX workload root.**

## 3. Cloud migration plan (Mac → DOKS)

**Wave 1 — stateless data/logic (safe, high value, always-on):**
- om-overwatch-{state,tak,interceptor,api}, om-http-gateway, om-callout, om-task-nex, om-nex-ui, resgate, om_presence_v2, nats-surveyor

**Wave 2 — the fusion engine:**
- om-unity/ORP → DOKS (after dedup). This is the big win: always-on fusion of 5,938+ entities in the cloud, not dependent on the Mac being awake.

**Wave 3 — hybrid (media near source, control in cloud):**
- om-overwatch-video / om-streaming-bridge — MediaMTX stays near cameras; control plane in cloud.

**MUST STAY ON MAC (hardware deps):**
- om-bridge-mavlink (local MAVLink serial/UDP)
- om-bridge-homeassistant (Home Assistant LAN)
- overwatch-vision (local camera frames)

## 4. NATS subject graph (current, air-only kill chain)

```
[MAVLink drone] → om-bridge-mavlink → tel.drone.*
                                          ↓
[radar/sensors]  → det.air.*  ──────→ om-overwatch-state → ent.air.{id}.state
                                          ↓
                                    om-overwatch-tak → tak.cot.out → [ATAK map]
                                          ↓
                                    om-overwatch-interceptor → intercept.*
[TAK human] → takserver-bridge → tak.cot.* ↔ om-bridge-constellation → entities
[Home Assistant] → om-bridge-homeassistant → iot.*
[HTTP clients] → om-http-gateway → NATS req/reply
```

**THE GAP:** om-overwatch-state only handles `det.air.*` + `tel.drone.*` → only `ent.air.*`. Ships/ground/space/cyber are invisible. **Generalizing om-overwatch-state + om-overwatch-tak to all entity types is the single highest-leverage change.**

## 5. Recommended target topology

- **DOKS (cloud, always-on):** ORP fusion + all stateless overwatch/gateway/task workloads + WorldWideView map + NATS core
- **Mac (leaf, hardware edge):** MAVLink bridge, Home Assistant bridge, local vision, MediaMTX
- **One NEX workload root** (kill the nex/ vs nats/nex/ split)
- **One ORP source of truth** (kill the om-unity vs ai-agents/orp duplication)
