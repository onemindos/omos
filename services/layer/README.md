# omos layer service (FEATURES A5)

Serves map layers to every client as **one snapshot read plus deltas**, so no client reads a
firehose and the `CLIENT` account never needs the JetStream API (BUS-DESIGN D7).

Contracts (onemind-xr, branch `run/m0`): `docs/contracts/layer-snapshot.md`,
`docs/contracts/layer-delta.md`; bus design `docs/design/BUS-DESIGN.md` D3, D7, sections 3.2,
3.4 and 5. Stacked on the entity service (omos PR #552), which writes `ENT`.

## What it reads and writes

| Direction | Subject / bucket | Notes |
|---|---|---|
| reads | `ENT` (`ent.*.*.state`) | one ordered ephemeral consumer: last message per subject on start, then every change, including `Nats-Marker-Reason` delete markers |
| reads | `tel.feed.*.status` (core) and the last one per feed in `TEL` at start | feed freshness for the manifests (omos-ingest status schema, PR #553) |
| reads | `tel.fleet.node.*.status` (core); hub stream source info | site contact and source lag |
| writes | KV `world-state`: `{layer}.meta`, `{layer}.cells.z{z}`, `_lease.layer-service` | single writer (lease); manifest every 30 s and on change, cells every 30 s |
| writes | KV `fleet-state`: `{site}` | every 10 s for each site in contact; a missing key means disconnected |
| serves | `svc.layer.snapshot` | request-many reply: chunks of at most 512 KiB and 2,000 items, the last has `"last": true` |
| serves | `svc.layer.watch`, `svc.layer.unwatch` | ack, then delta batches to the reply subject until the lease ends |
| serves | `svc.layer.list`, `svc.layer.meta` | layer summaries and the current manifest (Hermes) |

The service creates no stream and no bucket: `ENT` is a CR in PR #552, the two buckets are
the existing CRs `platform/nats/streams/kv-world-state.yaml` and `kv-fleet-state.yaml`.

## Behaviour that the contracts leave open (also in the owner proposal)

- **Replay window.** The watch path keeps the last 5 minutes (at most 200,000 changes) of `ENT`
  changes in memory. A `from_seq` older than that, or a `watch_id` the service does not know
  (for example after a restart), gets `{"error":"resnapshot"}`.
- **Expiry sweep.** Besides `ENT` delete markers, the service removes an entity when its
  `expiry_t` passes (checked every second). This covers `ENT` without per-message TTL
  (`ENT_MSG_TTL=false` in the entity Deployment).
- **Renew.** A renew may use a new reply subject; batches then go there. Clients that send the
  renew with the watch inbox as reply subject keep one inbox, and a cross-account response
  mapping is refreshed by every renew.
- **`from_seq: 0`** starts a watch without a snapshot: the current view comes as `enter` ops.
- **Releasability.** Same-account callers (no server-set `Nats-Request-Info`) are trusted
  services. Cross-account callers are mapped to groups by `PRINCIPALS_FILE`; CLEAR and GREEN
  are visible to every known caller, AMBER and AMBER+STRICT need a group in `releasable_to`
  (with an empty list, the `org_group`), RED needs an explicit group. Unknown identities get
  `forbidden`. A request with no header that crossed a service import (reply subject `_R_.`,
  an import without `share: true`) is refused with `forbidden`, not trusted.
- **Limits.** An entities-mode watch holds at most the layer `cap` (`too_many` at start or
  renew, and a `too_many` batch that ends the watch when the view grows past it). All watches
  together hold at most `WATCH_MEM_BUDGET_MB`; over it, the largest watch gets `resnapshot`.
  Queued ops point at the shared entity and are rendered when the batch is built. A cells
  watch takes a new diff only after the last one was sent.
- **Layer membership.** Default layers follow layer-snapshot.md section 1, mapped onto the
  ONTOLOGY types and kinds the entity service (omos #552) writes. The contract's `hazard`,
  `satellite`, `train` and `station` types do not exist on `ENT`, so a layer can narrow its
  entity types by `kinds`, by `source_subjects` (prefixes of a source's `subject`) and by
  `exclude_kinds`:

  | layer | ENT type | narrowed by |
  |---|---|---|
  | `hazards` | `event` | kind `weather_alert`, `earthquake`, `wildfire`, `flood`, `storm`, `volcano`, `tsunami`, or a source subject under `intel.weather.`, `intel.seismic.`, `intel.fire.`, `intel.disaster.`, `intel.natural.`, `evt.world.fire.`, `evt.world.seismic.` (conflict and chokepoint events stay out) |
  | `space` | `vehicle` | kind `satellite` (`norad-*`) |
  | `transit` | `vehicle` | kind `train` (`gtfs-*`) |
  | `fleet` | `drone`, `robot`, `vehicle`, `node`, `edge` | excludes kinds `satellite`, `train` |
  | `rf` | `sensor` | kind `aprs`, `meshtastic`, `rf`; no entity-service adapter emits these yet, so the layer is empty |

  The manifest carries `kinds`, `source_subjects` and `exclude_kinds` (additive) next to
  `filters`. A `zone` belongs to the `tak` layer only when one of its sources has
  `data_type: cot`; every zone is in `zones`.
- **Manifest `filters`.** A layer with several entity types has `filter: "ent.*.*.state"` and
  the exact list in `filters` (for a direct-path ordered consumer with `FilterSubjects`).

## Configuration (environment)

| Variable | Default | Meaning |
|---|---|---|
| `NATS_URL`, `NATS_USER`, `NATS_PASSWORD` | `nats://127.0.0.1:4222` | connection; credentials from the `omos-nats-creds` Secret |
| `SITE` | `hub` | writer site (`origin` default, own `fleet-state` key) |
| `LAYERS_FILE` | built-in table | JSON array of layer definitions (`internal/layer/config.go`) |
| `PRINCIPALS_FILE` | none | identity to groups map (ConfigMap); without it, cross-account callers are refused |
| `LEASE_BUCKET`, `LEASE_KEY`, `LEASE_TTL` | `world-state`, `_lease.layer-service`, `15s` | single-writer lease; `LEASE_BUCKET=none` for local runs |
| `ENT_STREAM`, `TEL_STREAM`, `WORLD_BUCKET`, `FLEET_BUCKET` | `ENT`, `TEL`, `world-state`, `fleet-state` | names |
| `REPLAY_WINDOW`, `REPLAY_MAX` | `5m`, `200000` | watch replay window |
| `MAX_WATCHES` | `200` | open watches (goroutines and timers) |
| `WATCH_MEM_BUDGET_MB` | `128` | estimated memory all watches may hold together (512Mi pod limit) |
| `WATCH_HEARTBEAT` | `10s` | empty batch when nothing changed |
| `MANIFEST_EVERY`, `CELLS_EVERY`, `SITE_EVERY` | `30s`, `30s`, `10s` | KV refresh |
| `SITE_SUBJECT`, `SITE_STREAMS` | `tel.fleet.node.*.status`, all streams | site inputs |
| `SITE_ALLOWLIST`, `MAX_SITES` | any valid site, `64` | sites accepted from `SITE_SUBJECT` (comma-separated list) |
| `HTTP_ADDR` | `:8080` | `/healthz`, `/readyz` (lease held and caught up), `/metrics` |

## Tests

```
NATS_SERVER_BIN=$(command -v nats-server) go test -race -count=1 ./...
```

The integration tests start a disposable `nats-server` (skipped if none is found), create
`ENT`, `TEL` and the two buckets with their CR settings, and run every acceptance test of both
contracts with `internal/layerclient`, a reference client that applies the contract's client
rules. Timers are shortened where the contract value would only lengthen the run;
`A5_LONG=1` runs the layer-down and dead-watch tests with the contract timers (5 minute bucket
TTL, 10 s heartbeat, 30 s dead watch). `A5_RESPONSE_THRESHOLD` changes the account export's
`response_threshold` in the cross-account test.

Local run against the local stack (no lease; the `world-state` and `fleet-state` buckets and
`ENT` must exist, for example created with the `nats` CLI from their CR settings):

```
LEASE_BUCKET=none NATS_URL=nats://127.0.0.1:4222 go run ./cmd/layer-service
```
