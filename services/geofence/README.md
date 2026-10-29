# omos geofence and alert service

The OneMind geofence and alert service (FEATURES D1 and D6). It is the only writer of KV
`geofences`, applying zone commands that the policy service approved and signed. It watches
fused entities for zone crossings and raises alerts for crossings, stale own assets and feed
outages.

- Design: onemind-xr `docs/design/ONTOLOGY.md` (zone and alert rows, action catalogue §6.2:
  `zone.create|update|remove` → `cmd.map.geofence.{id}.{verb}`, level A2),
  `docs/design/BUS-DESIGN.md` (KV `geofences`, `evt.map.geofence.{id}.{enter|exit}` in
  `MISSION`), `docs/contracts/asset-command.md` §2 and §4 (the policy flow; the policy
  service never executes anything itself).
- Owner proposals: `docs/proposals/2026-10-29-geofence-alert-service.md` (D1-P1 … D1-P7).
- Schemas: `services/schemas/map/geofence.v1.schema.json` (KV value),
  `services/schemas/alert/alert.v1.schema.json` (alert facts),
  `services/schemas/actions/zone.{create,update,remove}.v1.schema.json` (action parameters).
- Deployment: `platform/omos/geofence/`, Argo Application `platform-omos-geofence` (manual sync).

## Subjects (all inside the locked grammar)

| Direction | Subject or bucket | Stream | Notes |
|---|---|---|---|
| in | `cmd.map.geofence.{id}.{create\|update\|delete\|remove}` | `CMD` | durable consumer `geofence-svc-cmd`; only commands with a valid policy `sig`, `expires_at` in the future and `approval.status: approved` are applied |
| in | KV `geofences` | `KV_geofences` | watched, so hand edits and other writers are seen too; an unreadable value is dropped and logged |
| in | `ent.*.*.state` | `ENT` | ordered consumer, last record per subject on start, then every change and delete marker; `ENT_MODE=core` for the interim core subject |
| in | `tel.feed.*.status` | core (`TEL` once it binds `tel.feed.>`, PR #553 A0-2) | the omos-ingest, Redpanda Connect and legacy feed-intel status forms |
| out | KV `geofences` | | value per D1-P1: what the projector (PR #563 A6-P2) and the policy service (PR #554 POL-9) read |
| out | `evt.map.geofence.{id}.{enter\|exit}` | `MISSION` | doctrine §9.7; `Nats-Msg-Id` deterministic |
| out | `evt.map.geofence.{created\|updated\|deleted\|refused}` | `MISSION` | the command result (four tokens; zone id in `data.zone_id`), `causation` = command id, `data.action_id` (D1-P2) |
| out | `evt.so.alert.{raised\|updated\|resolved}` | `EVT` | the alert stream (D1-P3); `Nats-Msg-Id: {alert_id}:{seq}` |
| rpc | `svc.geofence.alerts` | | `{kind?, entity_id?, feed?}` → `{v, alerts: [...]}` (open alerts) |

## Zones

A zone record has `geo{shape, kind, coords (GeoJSON), radius_m, alt_min_m, alt_max_m}`.
Shapes: `polygon` (with holes), `multipolygon`, `circle` (Point + `radius_m`), `corridor`
(LineString + `radius_m`). Kinds: `aoi`, `keep_in`, `keep_out`, `hazard`, `control_area`.
`keep_in` and `keep_out` must be `polygon` or `circle`, because the policy service fails closed
on any other keep-in/keep-out shape and would deny every movement (POL-9).

Crossings are published for zones with `watch: true` (default for `keep_in`, `keep_out` and
`hazard`; AOIs and control areas only filter what TAK receives and can be large).
A zone watches the entities in `properties.applies_to` if set, else the types in `watch_types`,
else every moving type for `keep_out`, `hazard`, `aoi` and `control_area`, and own assets
(`ASSET_TYPES`) for `keep_in`. **People are never watched**: `person` entities and person-class
tracks are skipped, and `watch_types: ["person"]` is refused `forbidden_action` (onemind-xr
JADC2-MAP §5 stop line; the projector excludes people too).

Alerts per zone kind: leaving a `keep_in` (warning), entering a `keep_out` (warning), entering a
`hazard` (caution), entering a watched AOI or control area (advisory). An alert is raised when
the condition starts (also when a zone is drawn around an asset that already violates it) and
resolved when it ends, when the zone is removed or deactivated, or when the entity is gone.

## Why an in-process index and not Tile38

Tile38 is MIT-licensed and runs as one binary without Docker, and omos already runs it in
`platform/onetak` (`tile38.yaml`). It is still not used here:

1. **Same answer as the other readers.** The policy service (keep-in/keep-out before a command)
   and the projector (AOI filter) evaluate the same KV records in process with planar ray
   casting and haversine circles. This service uses the same tests (and the projector's
   corridor buffer), so a point is inside for all three or for none. Tile38 has its own
   geometry engine and would disagree at edges.
2. **Exactly-once facts and restart.** Crossings and alerts are derived from our own state with
   deterministic `Nats-Msg-Id`s and rebuilt from `MISSION` and `EVT` on start. With Tile38 the
   membership would live in a second stateful store, its hooks would need a bridge into
   JetStream, and a Tile38 restart would replay or drop crossings.
3. **Size.** The `zones` layer stays under 1,000 zones (layer-snapshot.md §1). A uniform grid
   (0.1° cells, big zones on a short list) answers each entity update with a handful of exact
   tests. The memory requests stay tens of MiB on a cluster whose hub-core pool is saturated.
4. **Sites and DDIL.** One static binary runs on the Mac or a field kit next to a leaf node,
   with nothing else to deploy.

If zone counts grow past tens of thousands or nearest-neighbour queries are needed, Tile38 can
sit behind the `geo` package without changing any subject.

## Exactly-once effects and restart

- Every output is published with a JetStream ack, in order, and retried with the same
  `Nats-Msg-Id` until acked; a full queue blocks the inputs.
- Membership starts as "outside", so every inside membership has an `enter` fact. On start the
  service replays `MISSION` `evt.map.geofence.*.{enter,exit}` and `EVT` `evt.so.alert.>` (last
  `RESTORE_LOOKBACK`, default 720 h), then reads all zones, then the last record of every
  entity, before it serves. An unchanged world publishes nothing; zones deleted and entities
  expired while it was down are closed (`zone_removed`, `entity_gone`; an asset's stale alert
  becomes `expired`).
- Commands: a redelivered command is recognised by `command_id` in the stored record (KV
  history for a delete) and answered with the same result id, so the zone is written once.
  `create` uses KV create; `update` is compare-and-set on the read revision (and on
  `revision` when the editor sends it, so an edit from a stale view fails `conflict`).
- One replica with `Recreate`. Two replicas would publish the same deterministic ids and be
  deduplicated within the 2-minute window only, so do not scale it.
- Limit: an entity that stays inside a zone for longer than `RESTORE_LOOKBACK` without a crossing
  has its `enter` outside the replay window; after a restart it gets a second `enter`
  (`first_fix`). Set `RESTORE_LOOKBACK=0` to replay all of `MISSION` (90 d).

## Configuration (environment)

| Variable | Default | Meaning |
|---|---|---|
| `NATS_URL`, `NATS_USER`, `NATS_PASSWORD` | `nats://127.0.0.1:4222` | credentials from a SealedSecret-backed Secret |
| `POLICY_KEYS_DIR` | — | directory of `{kid}.pem` Ed25519 public keys (ConfigMap); required while `REQUIRE_SIGNATURE=true` |
| `REQUIRE_SIGNATURE`, `REQUIRE_APPROVAL` | `true`, `true` | `false` only on a local stack without the policy service |
| `GEOFENCE_BUCKET`, `ENT_STREAM`, `CMD_STREAM`, `MISSION_STREAM`, `EVT_STREAM` | `geofences`, `ENT`, `CMD`, `MISSION`, `EVT` | the service never creates streams or buckets; it fails if one is missing |
| `ENT_MODE` | `jetstream` | `core`: core subscription on `ent.*.*.state` |
| `CMD_DURABLE` | `geofence-svc-cmd` | |
| `RESTORE_LOOKBACK` | `720h` | `0` replays everything |
| `FEED_SUBJECT` | `tel.feed.*.status` | |
| `FEED_SILENCE` | `90s` | no status for this long raises `feed_outage` (`silent`) |
| `ASSET_TYPES` | `drone,robot,vehicle` | stale alerts and the default keep-in scope |
| `STALE_DEFAULT_S` | `60` | for an asset record without `stale` |
| `TICK` | `5s` | silence checks |
| `SERVE_RPC` | `true` | `svc.geofence.alerts` |
| `HTTP_ADDR` | `:8080` | `/healthz`, `/readyz` (zones and entities loaded), `/metrics` |

## Tests

```bash
cd services/geofence
go vet ./...
NATS_SERVER_BIN=$(command -v nats-server) go test -race -count=1 ./...
```

Unit tests cover the geometry and grid index, record building and refusals, the signature rule
(copied JCS from the policy service), and the engine (crossings, alert levels and reasons,
people never watched, stale assets, the three feed status forms, restore idempotence and
reconcile after downtime), and validate records, parameters and alerts against the schemas.
The integration tests (`internal/bus`) start their own nats-server with the streams of the CRs
and a test key standing in for the policy service: sim-x1 leaving a keep-in raises an alert, a
stale and a silent feed raise alerts, a restart publishes nothing new and applies a redelivered
command once, and the command gate refuses unsigned, wrongly signed, expired and unapproved
commands. Without a nats-server binary they skip.

## Not in this slice

- The policy catalogue entries `zone.create|update|remove` (D1-P4, a patch to PR #554).
- `alert.acknowledge` / `alert.resolve` by humans (catalogued in #554 as `cmd.entity.{id}.*`,
  which needs an alert entity owner; D1-P5).
- The `ent.zone.{id}.state` mirror of each zone (ONTOLOGY zone row) and `ent.alert.*` records:
  `ent.*` has one writer, the entity service (D1-P5).
- Shared Go module for the geometry and JCS code now copied between policy, projector and
  this service.
- A Prometheus rule on `omos_geofence_alerts_open`.
