# Owner proposal: ENT stream and entity-service bucket (FEATURES A4)

Status: **proposed, needs owner approval**. Date: 2026-10-29. Filed under decision 0004 item 8
(onemind-xr `docs/decisions/0004-owner-defaults-for-the-run.md`: "use the existing grammar only;
if something truly has no pattern, use the closest existing pattern and write a proposal in
omos `docs/`"). This PR changes **no live subject and no live stream**: both new CRs are
excluded from the `platform-nats-streams` Application, and `platform-omos-entity` is manual
sync.

Design sources (onemind-xr, branch `run/m0`): `docs/design/BUS-DESIGN.md` D2, sections 4.1, 4.2,
5 and 10; `docs/design/ONTOLOGY.md` section 9 (ONT-P1, ONT-P8); `docs/contracts/entity.md`;
`docs/research/fusion-and-tracking.md` and `docs/research/nats-platform.md`.

## Summary

| ID | Ask | Same as (combined register, BUS-DESIGN 10.1) | Used until ruled |
|---|---|---|---|
| A4-P1 | Approve the `ENT` stream CR `platform/nats/streams/ent.yaml` | ONT-P1, BUS-DESIGN D2 | Not applied. The service runs locally against a disposable server; in the cluster it waits (or runs `ENT_MODE=core`, see below) |
| A4-P2 | Approve the bucket CR `platform/nats/streams/kv-entity-correlation.yaml` | BUS-DESIGN section 5 row `entity-correlation` | Not applied. The service refuses to start without its lease bucket |
| A4-P3 | Ratify `evt.sys.{source}.rejected` as the dead-letter subject, with `{source}` = the rejecting service (`evt.sys.entity.rejected`) | BUS-P9 | `evt.sys.entity.rejected` (bound by `SYS`, `evt.sys.>`) |
| A4-P4 | Add an `entity_type` token for satellites (`satellite`), or confirm `vehicle` + `kind: satellite` | ONT-P8 | `vehicle`, `kind: satellite`, `environment: space` |
| A4-P5 | Accept domain `entity` for `evt.entity.{id}.{event}` and add `evt.entity.>` to `EVT` | BUS-P5, ONT-P2 | **Not published.** Deletion uses the final `live: false` record and `Nats-TTL` |

## A4-P1: the `ENT` stream

**Gap.** The grammar has `ent.{entity_type}.{entity_id}.state` ("fused entity state", omos
`CLAUDE.md` grammar block), but no stream stores `ent.>`: the union of `subjects:` in
`platform/nats/streams/*.yaml` has no `ent.` prefix (`grep -h -A3 subjects platform/nats/streams/*.yaml`).
A core-only `ent` publish is lost for any consumer that is down, and late joiners have no
snapshot.

**Proposed CR** (`platform/nats/streams/ent.yaml`): `ENT`, subjects `ent.*.*.state`, R3 file,
`maxMsgsPerSubject: 1`, `allowDirect: true`, `allowMsgTtl: true`, `subjectDeleteMarkerTtl: 15m`,
`duplicateWindow: 2m`, `maxAge: 168h` backstop, `maxBytes: 1 GiB`, `discard: old`,
`preventDelete: true`. These are the BUS-DESIGN section 4.1 values. The 1 GiB cap is included
in that section's per-node budget (about 32.4 GiB of 45 GiB, arithmetic, not a measurement).

**Evidence checked for this PR (read-only, 2026-10-29):**
- The live Stream CRD `streams.jetstream.nats.io` version `v1beta2` has `allowDirect`,
  `allowMsgTtl`, `subjectDeleteMarkerTtl` and `maxMsgsPerSubject`
  (`kubectl --context do-nyc3-omos-cloud get crd streams.jetstream.nats.io -o json`); the
  controller is `natsio/jetstream-controller:0.23.0` (`kubectl -n nats get pods`), matching
  chart `nack` 0.34.0 (`appVersion: 0.23.0` in the chart index). The same fields are in
  `nats-io/nack` `v0.23.0` `deploy/crds.yml`. This settles the "NACK accepts the TTL fields"
  item that BUS-DESIGN 4.1 prerequisite 2 marked unverified, at the CRD level; that NACK passes
  them to the server is shown only by applying the CR.
- The live server is `nats:2.14.2-alpine` (`kubectl -n nats get statefulset nats`). Per-message
  TTL has a stale-TTL fix in 2.14.8 (nats-platform.md R9). The deployment therefore ships with
  `ENT_MSG_TTL: "false"`; turn it on after the server moves to 2.14.9. Without `Nats-TTL` an
  ended entity is still removed for clients: the service publishes a final `live: false`
  record and, 60 s after that record is acked, deletes it from `ENT` by its stream sequence
  (`TestFinalRecordDeletedWhenTTLOff`). Final records found on restart are deleted the same
  way. `maxAge` 7 d stays the storage backstop.
- Locally (nats-server 2.15.2, and 2.14.9 in CI), the integration tests in
  `services/entity/internal/bus` create `ENT` with exactly these settings and pass: one message
  per entity, `Nats-Msg-Id` deduplication, no growth on a full replay.

**Eviction risk under `discard: old`.** `ENT` keeps one record per live entity, but the
1 GiB cap is shared. If live churn (aircraft, vessels; A14 plans 50k vessels at about 1.7 KB
per record) ever fills it, `discard: old` evicts the oldest records first, and those are the
long-lived and pinned entities (pinned heartbeat 3.5 d). A restart would then not restore
them. Deleting final records keeps ended entities from adding to this. If the owner wants a
margin, lower `maxAge` to at most the pinned heartbeat (84 h) or raise `maxBytes`; this PR
keeps the BUS-DESIGN values.

**Operational note.** A purge of `ENT` keeps its duplicate window. A rebuild from scratch must
wait 2 minutes after the purge before replaying, or earlier `{entity_id}:{seq}` ids are
deduplicated away (found by the rebuild test; `services/entity/README.md`).

**Activation, in order:** approve; remove `ent.yaml` and `kv-entity-correlation.yaml` from the
`exclude` of `platform/argocd/platform-nats-streams.yaml`; confirm NACK reconciled both
(`nats stream info ENT`, `nats kv info entity-correlation`); let `build-entity-service` push
the image and pin its digest in `platform/omos/entity/deployment.yaml`; sync
`platform-omos-entity` once under observation.

**Interim without `ENT`** (entity.md section 1): `ENT_MODE=core` publishes the same records on
core NATS only; restart cannot restore state, so `seq` starts at the epoch millisecond to keep
increasing across restarts. This mode is for local and transitional use; it is not set in the
manifest.

## A4-P2: the `entity-correlation` bucket

**Gap.** The service must run as a single active writer (entity.md section 2, "single active
replica, KV lease"). No existing bucket fits: `ingest-locks` is proposed for collectors, and
reusing an unrelated bucket would mix owners.

**Proposed CR** (`platform/nats/streams/kv-entity-correlation.yaml`): bucket
`entity-correlation`, R3 file, `history: 5`, `maxBytes: 64 MiB`, `preventDelete: true`
(BUS-DESIGN section 5). Key `lease.entity-service` holds `{holder, expires_ms}`; it is created
with `Create` and renewed or taken over only with `Update(key, value, lastRevision)`, so two
replicas cannot both hold it (`TestLeaseSingleActive`). The same bucket is where manual
correlation and `decorrelated_from` will live when `cmd.entity.*` edits are built.

## A4-P2b: identity keys and the record (deviations from entity.md, for the owner)

- **Soft keys.** `mavlink_sysid` and `tak_callsign` are kept in `aliases.ids` but never resolve
  or merge an entity: sysid 1 is the ArduPilot and PX4 SITL default, and two TAK devices can
  share a callsign. Only hard keys (`sim`, `icao`, `mmsi`, `norad`, `gtfs`, `cot_uid`,
  `fleet_uid`, `det_track`, `intel_id`) and declared links merge (entity.md section 5).
  Tests: `TestSameSysidTwoSimDronesAreTwoEntities`, `TestSameCallsignTwoEUDsAreTwoEntities`.
- **Fleet prefix.** A non-sim fleet asset is `fleet-{uid}` (for example `fleet-x47`), following
  entity.md section 3's `{prefix}-{native}` rule. entity.md section 3 does not list a `fleet`
  prefix; the owner should add it there or name another.
- **`sources[].tier`.** Each source carries its attribute tier (`human`, `self_report`,
  `sensor`) so a restart restores the conflict tiers exactly. This is an additive optional
  field: `sources` items in `entity.v1.schema.json` allow extra properties.
- **`svc.entity.query` paging.** Pages of at most 200 entities (the default) with a `cursor`
  (last `entity_id`; `null` on the last page). A page is halved until it fits the server's
  `max_payload` (1 MiB live; `TestQueryPagesLargeResult` reads 1,200 entities).

## A4-P3: dead-letter subject

Rejected inputs (stop-line SL9: a `person` with `disposition: hostile`; unparseable messages)
are published on `evt.sys.entity.rejected` with `{reason, service, subject, stream, seq,
integration, source_id, entity_type, disposition, t}`. `SYS` binds `evt.sys.>`
(`platform/nats/streams/sys.yaml`), so they are stored without a CR change. They are
JetStream publishes with `Nats-Msg-Id: {stream}:{seq}` of the input; if the publish is not
acked, the input is not acked either and is redelivered (`TestDeadLetterIsAckedOrInputRedelivered`). BUS-P9 proposes
`evt.sys.{source}.rejected`; this PR reads `{source}` as the service that rejects.

## A4-P4: satellite entity type

entity.md section 3 has no `satellite` token although its id examples include `norad-25544`
and the projector's `COT_TYPE_MAP` has `satellite` (`platform/omos/tak-projector.yaml`). Until
ruled, satellites from `feed-satellites` are `ent.vehicle.norad-{n}.state` with
`kind: satellite` and `environment: space`. MTA trains are `vehicle` with `kind: train`; BUS-P6
covers their world class.

The same applies to satellites propagated from `intel.space.celestrak.tle` (FEATURES A12) and to
trains on `evt.world.rail.*` (#558 A1-P4): both resolve to the entity the legacy CoT makes
(`norad-{n}`, `gtfs-{feed}-{trip}`), so the parallel run in #558's cut-over shows one entity each.

Intel events follow the same rule for the omos-ingest parallel run (#553, BUS-P3). The adapter
reduces a URL record id to its last path segment, so the legacy NWS `feature.id`
(`https://api.weather.gov/alerts/urn:oid:...`) and the omos-ingest `properties.id`
(`urn:oid:...`) give one entity. It also maps legacy `source` values to the grammar token
(`nasa_firms` to `firms`, `nasa-eonet` to `eonet`). Known migration gaps, for the owner:

- **FIRMS id shape**: the legacy collector's id is a SHA-1 of `lat:lon:acq_date:acq_time`; the
  #553 CSV source's id is `firms:{lat}:{lon}:{date}:{time}`. The prefix now matches, but the
  native ids differ, so a hotspot seen by both shows twice until one collector is retired or
  #553 adopts the legacy key.
- **Event time**: #553's GeoJSON parser defaults `timeField` to `time` and reads epoch
  milliseconds, so USGS envelopes already carry the quake time in `ts`. A source whose time
  field does not parse falls back to `fetchedAt`; the adapter therefore prefers a numeric
  `data.time` in milliseconds when present, as the flat path did.

## A4-P5: lifecycle events (information)

BUS-DESIGN BUS-P5 says lifecycle events are "not published until ruled". The service follows
that. If the owner accepts domain `entity`, add `evt.entity.>` to `EVT` and set the service to
emit `created`, `expired`, `deleted` and `conflict` events.

## A4-P6: SGP4 library (information)

FEATURES A12 needs SGP4 in Go. The most widely used pure-Go implementation is
`github.com/joshuaferrara/go-satellite` (BSD-2-Clause; last commit June 2022; cgo-free). It has
three defects the wrapper in `services/entity/internal/orbit` works around and tests: `log.Fatal`
on a malformed TLE number (prevented by validating every column first), whole-second epoch
truncation (compensated exactly), and a wrong longitude wrap in `LatLongDeg` (not used). Measured
against python-sgp4 the result is within 3.7 m over 7 days. **Proposed:** accept the dependency.
Alternative if the owner prefers no unmaintained dependency: vendor a reviewed copy under
`services/entity/third_party/` (BSD-2 permits it) or port Vallado's public-domain reference.

## Risks the owner should know

- **Duplicate dots in TAK** (fusion-and-tracking.md Risk 1): `feed-mta` and `feed-satellites`
  already put CoT on `tak.cot.inject`. The service marks entities sourced only from TAK with
  `indicators.egress_tak: false` (cot-projection.md section 3.1), but the live projector does
  not read that flag yet (FEATURES A6). Until A6 lands, a synced entity service would make the
  projector draw a second `OM.*` marker for those trains and satellites.
- **Consumers on `COT`**: the durable `entity-svc-cot` starts 10 minutes back
  (`START_LOOKBACK`), not at the 30-day start of `COT`, so the first sync does not replay
  9.5 million messages (omos-live-state.md section 2).
- **Shared credentials**: the service reuses `omos-nats-creds` until per-service users exist
  (BUS-DESIGN D6). Publishing `ent.>` is therefore not yet restricted to it.
