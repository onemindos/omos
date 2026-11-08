# Owner proposals: feed grammar cut-over (FEATURES A1)

Date: 2026-10-29. PR: `feat/feed-grammar-fixes`. Runbook: `platform/data/redpanda-connect/CUTOVER.md`.

These are the grammar gaps and decisions the A1 change needs (decision 0004 #8: gaps go to the
owner, the PR changes no live subject). As merged, the PR publishes nothing new: the grammar
pipelines run in a Deployment at `replicas: 0`.

## A1-P1: `rail` world class (BUS-P6)

**Gap.** Doctrine §9.1 lists the moving world classes `air` and `sea` but nothing for trains;
§9.3 has no rail type. MTA trains need a subject and a type.

**Evidence.** omos `docs/NATS-DOCTRINE.md` §9.1 (`evt.world.air.{id}`, `evt.world.sea.{id}`), §9.3;
BUS-DESIGN §10 BUS-P6; ingestion migration plan step 8.

**Proposed ruling.** Add `evt.world.rail.{id}` to §9.1 (stream `WORLD`, already captured by
`evt.world.>`) and `ground.rail` to the §9.3 world types. Ids are `{operator}-{feed}-{trip}`,
lowercased, with any character outside `[a-z0-9_-]` replaced by `_` (MTA: `mta-mnr-1211`,
`mta-lirr-go202_26_5`). The record is the §9.2 envelope; `meta` carries `agency`, `feed`,
`trip_id`, `route_id`, `vehicle_id`, `vehicle_label`, `stop_id`, `current_status`, `start_date`.
Bus and light-rail feeds, if added later, use `transit` the same way.

**Used until ruled.** Nothing is published: Stage 1 of the runbook (replicas 1) waits for this ruling.

## A1-P2: who turns TLEs into positions

**Gap.** BUS-DESIGN §3.1 puts raw TLEs on `intel.space.celestrak.tle` and makes the entity
service propagate them (FEATURES A12). The entity service in PR #552 does not propagate yet, so
after the cut-over satellites would leave the picture.

**Proposed ruling.** Keep propagation in the entity service (A12). The legacy `feed-satellites`
CoT stays on until A12 lands, or until the owner accepts satellites off TAK (they reach no ATAK
today anyway, FEATURES A21). Not proposed: a collector-side `evt.world.space.{norad}` position
subject, which would put derived data on a raw-observation subject.

**Used until ruled.** Legacy `feed-satellites` keeps publishing; the TLE pipeline may run in
parallel (Stage 1) because `intel.space.celestrak.tle` is already in BUS-DESIGN §3.1.

## A1-P3: status and dead letters from the interim Redpanda pipelines

**Gap.** The grammar pipelines run on Redpanda Connect until MTA and CelesTrak migrate to
`omos-ingest` (PR #553). They follow its contracts with three differences.

**Proposed ruling (information).**
- `tel.feed.{mta|celestrak}.status` is published once per poll (30 s, 4 h), not every 10 s, with
  `collector: "redpanda-connect"`. `stale` is true only when the poll failed; the
  freshness alert itself comes with A3. Core publish, so it is stored once `TEL` adds `tel.feed.>`
  (A22).
- Dead letters go to `evt.sys.{source}.rejected` (BUS-P9 reading), as bare
  `{source, reason, raw, rawBytes}` records (`services/schemas/ingest/rejected.v1.schema.json` in
  #553) without the bus envelope #553 wraps them in.
- `Nats-Msg-Id`: `mta:{feed}:{trip}:{vehicle timestamp ms}` and `tle:{norad}:{epoch}` (the ids
  in the ingestion proof), so a re-poll inside the 2 min duplicate window is not stored twice.

## A1-P4: entity service follow-ups (not grammar; for PR #552 or its successor)

Before Stage 2, the entity service must accept what these pipelines publish:
1. `World` adapter: add class `rail` (entity type `vehicle`, kind `train`, key on `uid`). Today
   `evt.world.rail.*` returns `SkipUnsupported`.
2. `Intel` adapter: TLE records have no position by design and return `SkipNoPosition`; A12 adds
   SGP4 propagation from `intel.space.celestrak.tle`.
