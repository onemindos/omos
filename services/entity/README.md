# omos entity service

The OneMind entity service (FEATURES A4). It reads observations from the bus, resolves which
real thing each one describes, fuses them, and publishes one record per entity on
`ent.{entity_type}.{entity_id}.state` into the `ENT` stream. Clients, Hermes and the TAK
projector read entities; only this service writes them.

- Contract: onemind-xr `docs/contracts/entity.md` (v 1). Field names are the contract's,
  exactly; the JSON Schema is `services/schemas/entity/entity.v1.schema.json`.
- Design: onemind-xr `docs/design/BUS-DESIGN.md` D2 and section 4.2, and
  `docs/research/fusion-and-tracking.md` "Recommendation".
- Stream and bucket: `platform/nats/streams/ent.yaml` and `kv-entity-correlation.yaml`, both
  owner **proposals** (`docs/proposals/2026-10-29-entity-service-ent-stream.md`).
- Deployment: `platform/omos/entity/`, Argo Application `platform-omos-entity` (manual sync).

## Subjects (all inside the locked grammar)

| Direction | Subject | Stream | Notes |
|---|---|---|---|
| in | `tak.cot.>` (`tak.cot.inject`, `tak.cot.out`, `tak.cot.mesh.rx`) | `COT` | CoT XML. The projector's own output (uid `OM.*` or `INTEL.*`, or `_flow-tags_ omos-tak-projector`) is ignored (cot-projection.md §8.2) |
| in | `tel.fleet.{class}.{id}.position\|status` | `TEL` | doctrine §9.2 geo-envelope |
| in | `evt.world.{air\|sea\|rail\|fire\|seismic\|conflict\|chokepoint}.{id}` | `WORLD` | military surge subjects are skipped (BUS-DESIGN D9). `rail` (#558 A1-P1, owner proposal BUS-P6) is a `vehicle` of kind `train`, keyed `gtfs-{feed}-{trip_id}` so it fuses with the legacy `MTA-*` CoT |
| in | `intel.>` | `INTEL` | the flat feed-intel payload or the omos-ingest bus envelope (`{v, id, ts, type: <subject>, source, data}`, #553); records without coordinates are skipped, so NWS zone-only alerts (no polygon) make no entity yet |
| in | `intel.space.celestrak.tle` | `INTEL` | one TLE per message (#558); propagated with SGP4 (`internal/orbit`) into `ent.vehicle.norad-{n}.state`, then re-propagated every `ORBIT_EVERY` until the next TLE |
| in | `det.{domain}.{type}.{id}` | `DET` | read only if the stream exists (it is proposed, not live) |
| out | `ent.{entity_type}.{entity_id}.state` | `ENT` | JetStream publish with ack; headers `Nats-Msg-Id: {entity_id}:{seq}`, `Nats-TTL`, `OneMind-Marking`, `OneMind-Releasable-To`, `OneMind-Schema: entity.v1` |
| out | `evt.sys.entity.rejected` | `SYS` | dead letters with a reason (stop-line SL9, unparseable input); BUS-P9 form |
| rpc | `svc.entity.get`, `svc.entity.query` | — | interim reads of entity.md section 8 (NATS micro service `entity`) |

`evt.entity.{id}.{event}` lifecycle events are **not** published: the `entity` event domain
waits for the owner (BUS-P5). Deletion still works: the final record has `live: false`.

## Pipeline

1. **Normalise** (`internal/adapt`): one adapter per subject family turns a message into an
   observation with hard keys, source time, position, `ce`/`le`, class, `stale` and its stream
   pointer. CoT `ce`/`le` 9999999 means unknown.
2. **Identity** (`internal/identity`): hard keys map to one `entity_id` through an alias map:
   `icao-{hex}`, `mmsi-{n}`, `norad-{n}`, `tak-{uid}`, `gtfs-{agency}-{trip}`, `fleet-{uid}`,
   `sim-*`. Soft keys (`mavlink_sysid`, `tak_callsign`) are kept in `aliases.ids` but never
   resolve or merge: they are not unique. A
   native id with characters a NATS token cannot carry is hashed
   (`{prefix}-{base32(sha256)[0:26]}`) and kept in `aliases.ids`. Declared links
   (`ENTITY_LINKS`) can join ids from different families; nothing else auto-merges. An AIS
   position that implies more than 50 m/s from the last fix is a duplicate MMSI and goes to
   `{id}-x2`, never averaged.
3. **Estimate** (`internal/fusion`): a constant-velocity Kalman filter in a local
   East-North-Up plane per entity, process noise by class. `Estimator` is the hook for the
   later IMM (CV + coordinated turn) and `Associator` the hook for GNN/JPDA on keyless
   detections; neither is implemented yet.
4. **Fuse** (`internal/engine`): per-source monotonic time (an older or equal `source_t` from
   the same source is dropped and counted), clock-skew replacement (more than 5 s in the
   future), per-attribute conflict resolution (human, then self-report, then sensor; newest
   within a tier; never downgrade to unknown; every resolution in `conflicts[]`), the overwatch
   rule (a sensor report does not move a track with a fresh self-report), decomposed confidence
   (`source_prior`, `track_quality`, `corroboration`), `egress_tak` and marking per
   cot-projection.md section 3.1, lifecycle `tentative/active/stale/lost/expired`.
5. **Publish** on change (moved more than `ce_m`, class, status or source set changed) or on a
   heartbeat of `min(stale/2, 120 s)`. `seq` increases by one per publish.

Safety rules in code: a `person` with `disposition: hostile` is rejected and dead-lettered
(SL9); a detected person is always a `track` of kind `person` with marking `AMBER+STRICT`,
`egress_tak: false` and a TTL of at most 60 s; hostile non-person assets are published as
`suspicious` with a conflict entry, because the contract has no `hostile`.

## Exactly-once effects, restart and rebuild

- Inputs are durable pull consumers (`entity-svc-{stream}`), explicit ack after the `ENT`
  PubAck. A new durable starts `START_LOOKBACK` back; an existing one keeps its position.
- On start the service restores every live entity from `ENT` (last message per subject):
  `seq`, aliases and the per-source last `source_t`. A redelivered input is then older than the
  stored source time and is dropped, so a restart never publishes a second record for the same
  observation. If the crash came between the `ENT` publish and the input ack, the republish
  would carry the same `seq`, so `Nats-Msg-Id` deduplication catches it too.
- A failed publish is kept and retried with the same `seq` and `Nats-Msg-Id`.
- **Rebuild from scratch**: stop the service, delete the `entity-svc-*` durables, purge `ENT`,
  **wait at least the `ENT` duplicate window (2 min)** (a purge keeps the window, so earlier
  `{entity_id}:{seq}` ids would be deduplicated away), set `START_LOOKBACK` to the replay
  start, and start the service.
- Single active replica: a lease in KV `entity-correlation`, key `lease.entity-service`,
  taken and renewed by compare-and-set; the value carries its own expiry. A replica that loses
  the lease exits. The service never creates streams or buckets; it fails if they are missing.

## Configuration (environment)

| Variable | Default | Meaning |
|---|---|---|
| `NATS_URL`, `NATS_USER`, `NATS_PASSWORD` | `nats://127.0.0.1:4222` | connection; credentials from the `omos-nats-creds` SealedSecret in the cluster |
| `ENT_MODE` | `jetstream` | `core` is the interim mode without `ENT`: core publish, no restore; `seq` starts at the epoch ms |
| `ENT_STREAM` | `ENT` | |
| `ENT_MSG_TTL` | `true` | send `Nats-TTL` (needs `allowMsgTtl` on the stream). The cluster manifest sets `false` until NATS 2.14.9; while `false`, each final `live:false` record is deleted from `ENT` 60 s after its ack |
| `INPUTS` | `COT=tak.cot.>;TEL=tel.fleet.>;WORLD=evt.world.>;INTEL=intel.>;DET=det.>` | missing streams are skipped |
| `START_LOOKBACK` | `10m` | where a new durable starts; `0` replays everything |
| `LEASE_BUCKET`, `LEASE_KEY`, `LEASE_TTL` | `entity-correlation`, `lease.entity-service`, `15s` | `none` disables the lease (local runs only) |
| `ENTITY_LINKS` | — | JSON `{"cot_uid:SOME-UID": "sim-x1"}` declared links |
| `SOURCE_PRIORS` | — | JSON per-integration `source_prior` overrides |
| `REJECT_SUBJECT` | `evt.sys.entity.rejected` | dead-letter subject; acked JetStream publish with `Nats-Msg-Id: {stream}:{seq}`, the input is redelivered if it fails |
| `SERVE_RPC` | `true` | `svc.entity.get` and `svc.entity.query` (`{types, bbox, limit, cursor}`; pages of at most 200, `cursor` is the last `entity_id`, `null` on the last page) |
| `ORBIT_EVERY` | `30s` | satellite re-propagation interval; negative disables. A TLE older than 7 days is no longer propagated |
| `HTTP_ADDR` | `:8080` | `/healthz`, `/readyz` (lease held and consuming), `/metrics` (counters) |

## Satellites (FEATURES A12)

`internal/orbit` propagates TLEs with SGP4 using
[`github.com/joshuaferrara/go-satellite`](https://github.com/joshuaferrara/go-satellite)
(BSD-2-Clause, a Go port of Vallado's reference SGP4; last commit 2022). The wrapper validates
every TLE column before calling the library (its parser calls `log.Fatal` on a bad number),
compensates for its whole-second epoch truncation, and converts TEME to WGS84 itself. Its test
compares 95 positions of 16 stations-group objects (including the ISS) out to 7 days
against python-sgp4 (`internal/orbit/testdata/gen_reference.py`): worst difference 3.7 m,
against the 5 km acceptance limit. The latest TLE per satellite is kept in memory and rebuilt
from the restored entity's `meta.line1/line2` on restart (JetStream mode only; in `core` mode
satellites reappear with the next TLE). Next-pass prediction and orbit rings are client work.

## Tests

```bash
cd services/entity
go vet ./...
NATS_SERVER_BIN=$(command -v nats-server) go test -race -count=1 ./...
```

The integration tests (`internal/bus`) start their own nats-server on a free port with a
temporary store, create the streams with the settings of the proposed CRs, replay
`testdata/fixtures/replay.jsonl` (generated by `gen_fixtures.py` in the wire formats of the live
omos feeds) and compare the `ENT` state with `testdata/golden/entities.json`
(`go test ./internal/bus -run TestGoldenReplay -update` rewrites it). They cover the entity.md
section 10 tests: two feeds give one aircraft with two sources, a ship updates in place, a
rebuild from scratch gives the same entities, a crash between publish and ack and a restart
publish no duplicate, a second replay of every input adds nothing to `ENT`, the hostile person
is dead-lettered, and the lease admits one holder. Without a nats-server binary these tests skip.

## Not in this slice

IMM estimator; GNN/JPDA association of keyless `det.*`; track-to-track correlation with
covariance intersection; `cmd.entity.*` edits and overrides; `svc.entity.search`,
`svc.entity.history`, `svc.fusion.explain`; `evt.entity.*` events (BUS-P5); per-service NATS
user. The per-type TTLs, stale values, process noise and source priors are starting values from
the design documents, not measured against feed cadence.
