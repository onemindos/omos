# Owner proposal: layer service, world-state and fleet-state (FEATURES A5)

Status: **proposed, needs owner approval**. Date: 2026-10-29. Filed under decision 0004 item 8
(onemind-xr `docs/decisions/0004-owner-defaults-for-the-run.md`). This PR changes **no live
subject, stream or bucket**: the two bucket CRs get comments only, and `platform-omos-layer`
is manual sync. It is stacked on PR #552 (entity service, `ENT`), whose proposals A4-P1..P5
must be ruled first.

Design sources (onemind-xr, branch `run/m0`): `docs/design/BUS-DESIGN.md` D3, D7, sections 3.2,
3.4, 5, 7 and 12; `docs/contracts/layer-snapshot.md`; `docs/contracts/layer-delta.md`.
Test evidence: onemind-xr `docs/evidence/M1/A5/`.

## Summary

| ID | Ask | Same as (BUS-DESIGN 10.1) | Used until ruled |
|---|---|---|---|
| A5-P1 | Ratify the purpose of the existing buckets: `world-state` holds layer manifests and density cells, `fleet-state` holds per-site status | BUS-DESIGN D3 | Nothing is written: the Application is not synced |
| A5-P2 | Add `maxBytes: 67108864` to `kv-world-state.yaml` and `maxBytes: 8388608` to `kv-fleet-state.yaml` | BUS-DESIGN section 5 | Not changed in this PR (live CRs, auto-synced Application) |
| A5-P3 | Keep the single-writer lease as key `_lease.layer-service` in `world-state` | — | As proposed (no new bucket) |
| A5-P4 | Ratify the endpoints `svc.layer.{list,meta,snapshot,watch,unwatch}` and the streamed reply of `snapshot` and `watch` | BUS-DESIGN 3.4 | As proposed |
| A5-P5 | Ratify `tel.fleet.node.{site}.status` as the site heartbeat the site monitor reads, and the lag rule | BUS-P1 | Subscribed; nothing publishes it yet |
| A5-P6 | Ratify the release rule and the principals ConfigMap | security-and-trust.md section 5 | As proposed; stage 1 has only same-account (trusted) callers |
| A5-P7 | For the `CLIENT` account (account split): export `svc.layer.>` with `response_type: stream` and a `response_threshold` of at least 2 minutes; import with `share: true` | BUS-DESIGN 7.1, 12 | Not configured (no `CLIENT` account yet) |
| A5-P8 | Additive contract details for onemind-xr `layer-snapshot.md` / `layer-delta.md` | — | Implemented as listed |

## A5-P1: purpose of `world-state` and `fleet-state`

Both buckets exist (`kubectl -n nats get keyvalues kv-world-state kv-fleet-state`, read-only,
2026-10-29: `world-state` TTL 5m history 1, `fleet-state` TTL 30s history 3) and have never
been written (`last_seq 0`, onemind-xr `docs/research/omos-live-state.md` section 5). No code in
`services/` or `platform/` uses them (`grep -rln 'world-state\|fleet-state' services platform`
finds only prose in `services/nex/overwatch/*` and the CRs).

| Bucket | Keys | Refresh | Why the TTL works |
|---|---|---|---|
| `world-state` | `{layer}.meta` (manifest, under 8 KiB), `{layer}.cells.z{z}` (cells, at most 512 KiB; a larger zoom is dropped from the manifest) | 30 s and on change; cells 30 s | 5 m TTL: a missing manifest means "layer down" |
| `fleet-state` | `{site}` = `{v, site, connected, last_contact, source_lag_s}` | 10 s | 30 s TTL: a missing key means "site disconnected" |

Single writer: one replica, `Recreate`, and the lease (A5-P3). Only the lease holder subscribes
to `svc.layer.*` and writes KV.

## A5-P2: caps

Proposed CR changes (not in this PR):

```yaml
# kv-world-state.yaml           # kv-fleet-state.yaml
  maxBytes: 67108864            #   maxBytes: 8388608
```

Sizing: 11 default layers, 2 of them with 3 cell zooms; the largest cell record measured in the
50,000-vessel test was 68,719 bytes (zoom 6). Changing `maxBytes` on an existing bucket is an
update NACK applies in place; it was **not** tried against a live server here.

## A5-P3: lease key

The service uses the entity service's lease design (compare-and-set on one key, expiry in the
value; `services/entity/internal/bus/lease.go`, copied). The key lives in the bucket it guards,
so no new bucket is needed. Clients reading `world-state` directly must ignore keys starting
with `_`. Test: `TestA5SingleWriter` (the standby stays unready while the lease is held and
takes over 20 ms after the holder stops).

## A5-P4: endpoints and streamed replies

All within `svc.{service}.{endpoint}`. `snapshot` answers with several messages on one reply
subject (request-many, terminated by `"last": true`); `watch` answers with an ack and then
batches on the reply subject until the lease ends. Both were run **across accounts** (an
`ONEMIND` export with `response_type: stream`, a `CLIENT` import with `share: true`):
`TestA5Releasability` received a 2-chunk snapshot and 25 s of watch batches (with renewals) in
the `CLIENT` account. This closes the BUS-DESIGN section 12 item "streamed service responses to
a `CLIENT` inbox" for nats-server 2.15.2 on this machine; the cluster runs 2.14.2 (unverified
there).

## A5-P5: site inputs

BUS-DESIGN section 3.1 lists `tel.fleet.node.{site}.status` as the site agent's status
(BUS-P1). The site monitor treats any message there as contact. For hub streams that source a
site stream through its JetStream domain (`$JS.{site}.API`, BUS-DESIGN sections 8-9),
`source_lag_s[{source}]` is 0 when nothing is pending, otherwise the seconds since the source
was last active. This is an interpretation of "source lag in seconds"; JetStream reports lag in
messages. Not exercised with a real site leaf (no domain-sourced stream on this machine).

## A5-P6: release rule

Caller identity follows the policy service (PR #554): the server-set `Nats-Request-Info`
header on a `share: true` import is the verified identity; same-account requests are trusted
services (they may read `ENT` directly anyway). Rule: CLEAR and GREEN to every known caller;
AMBER and AMBER+STRICT to callers with a group in `releasable_to`, or with the `org_group` when
the list is empty; RED only with an explicit group; unknown identities get `forbidden`. The
entity service defaults to `AMBER` with an empty list when no source sets a marking
(`services/entity/internal/engine/engine.go`), so most entities are organisation-only for
`CLIENT` users until markings are set per source. Test: `TestA5Releasability` (an
AMBER+STRICT entity for team-a is absent from team-b's snapshot and watch, and moves between
them when re-marked).

## A5-P7: account export settings

With a 0.5 s heartbeat and renewals every 5 s, watches ran 25 s without a gap both with
`response_threshold: 15s` and `3s` (`docs/evidence/M1/A5/experiment-response-threshold-3s.txt`),
so the mapping is not cut at the threshold in this setup. The mechanism was not traced in the
server source. Conservative setting: `response_threshold` of at least 2 minutes (more than the
10 s heartbeat and the default 60 s lease).

## A5-P8: additive contract details (onemind-xr)

| Item | Implemented as |
|---|---|
| Manifest `filter` for layers with several entity types | `"ent.*.*.state"` plus an added `filters` array with one filter per type |
| Error replies | add `"last": true` so a request-many read ends; new code `bad_request` for unreadable requests and bad AOIs |
| Watch with `from_seq: 0` | no replay; the current view is sent as `enter` ops |
| Renew | may carry a new reply subject; batches move there. An unknown `watch_id` gets `resnapshot` |
| Unknown values in render entities | `hae` and `expiry_t` are `null` (contract says number) |
| `aoi` `zone_id` | resolved from the `ent.zone.{id}.state` record's `geo.coords` (GeoJSON Polygon, or Point with `radius_m`) |
| `zoom` omitted | entities mode (the cap still applies) |
| Watch on a large layer in entities mode | not capped; ops are paced by `max_rate_hz` and coalesced |
