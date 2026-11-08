# Feed platform: 100+ sources, each switchable live

Status: **design, waiting for owner approval**. Date: 2026-11-08. Nothing in this document is built
or deployed yet. Facts marked "verified" were checked read-only on `do-nyc3-omos-cloud` or by
running the pinned image locally on 2026-11-08.

Background reading: the research report (OSIRIS, God's Eye View, Argos Atlas, WorldWideView) in the
cleanup tracker, `docs/NATS-DOCTRINE.md` sections 2, 3 and 9, and the open PRs #553 (omos-ingest)
and #558 (feed grammar fixes), which this design builds on.

## 1. The idea in one paragraph

Each data source is **one small YAML file in Git** that says what it is, where it publishes, its
licence and how often it runs. A **NATS key-value bucket, `feed-control`**, holds the live switch
for each source (on or off, and its interval). A **controller** watches that bucket and starts or
stops the source inside a shared engine pod, so turning a feed on or off is one `nats kv put` and
needs no Git change, no Argo sync and no restart. Every source reports its health on the bus, and
a single archive pipeline copies feed data into ClickHouse for the time machine.

```
 Git (one file per source)            NATS KV feed-control (live switch)
 platform/data/sources/<id>.yaml      key = source id, value = {enabled, interval}
          │  Argo                               │ watch
          ▼                                     ▼
 ┌──────────────── feeds engine pod (Redpanda Connect, streams mode) ────────────────┐
 │  feed-controller stream ── starts / stops ──►  source streams (mta, ships, ...)   │
 │  feed-health stream (reads the engine's own API, publishes status)                │
 └───────────────────────────────────────────────────────────────────────────────────┘
          │ raw observations, geo-envelope (doctrine 9.2)        │ tel.feed.{id}.status
          ▼                                                      ▼
 WORLD / INTEL streams ──► entity service (#552) ──► ent.{type}.{id}.state ──► TAK projector ──► CoT ──► TAK
          │
          └──► archive pipeline (JetStream durable ──► ClickHouse onemind.events)

 Awkward sources (auth refresh, paging, scraping, life-safety): omos-ingest (Go, #553),
 same manifest format, same feed-control bucket, same status subject.
```

## 2. What exists today (verified 2026-11-08)

| Piece | State |
|---|---|
| Redpanda Connect pod `redpanda-connect` (namespace `feeds`, renamed from `omos` on 2026-11-08 by PR #760; hub-core pool) | Image `redpandadata/connect@sha256:b462…` = **4.103.1**. Limit 128 MiB, using about 90 MiB. Runs in streams mode, but its args start **only `feed-mta`**. |
| `feed-aircraft`, `feed-hazards`, `feed-ships` | In the ConfigMap, **not running**. `feed-hazards` duplicates the USGS and NWS collectors in `feed-intel`. |
| `feed-mta` | Running. Builds CoT XML inside the pipeline and publishes to `tak.cot.inject`. |
| Python pods | `feed-intel` (USGS, NOAA, EMSC, FIRMS, GDACS, EONET, ACLED), `feed-satellites`, `sink-clickhouse`, `omos-tak-projector`. |
| Health | The `feed-intel` collectors already publish `tel.feed.{id}.status` every poll (seen live: `emsc`, `firms`, `acled`, `gdacs`). |
| Live switch | `feed-intel` already watches the `feed-config` bucket (`feed.{id}.enabled`, `feed.{id}.poll_interval`). The bucket is **empty**, so everything runs on defaults. It is R1, app-created, and has no CR. |
| Entity service and `ENT` | `platform-omos-entity` is manual sync and not deployed. The `ENT` stream CR is excluded until the owner approves `docs/proposals/2026-10-29-entity-service-ent-stream.md`. |
| `WORLD` stream (`evt.world.>`) | Exists, R3, 24 h, **0 messages**. |
| TAK projector | Already reads `ent.*.*.state` and `intel.>` and builds CoT. So CoT-in-one-place is half done. |
| ClickHouse sink | Python. Reads `intel.>` (durable) and `ent.*.*.state` (core). **Discards a batch** when ClickHouse is down. |
| Streams API | On by default in 4.103.1 (no `--bind-http` gate in this version). Before the rename a hand-applied Service `redpanda-connect` (not in Git) exposed port 4195 to the whole cluster. In `feeds` (verified 2026-11-08) there is no Service and **no NetworkPolicy**, so the port is still reachable on the pod IP from any pod. |
| Open PRs | **#553** omos-ingest (Go collector framework, NWS first, merges cleanly). **#558** grammar MTA pipeline on `evt.world.rail.*` plus satellites on `intel.space.celestrak.tle`, in a new `redpanda-connect-grammar` Deployment at 0 replicas on the ecosystem pool (needs a small rebase). |

## 3. The pieces

### 3.1 Source manifest (one file per source)

Location: `platform/data/sources/<id>.yaml`. JSON Schema:
`services/schemas/feeds/source-manifest.v1.schema.json`, checked in CI by a new
`scripts/validate_feed_sources.py` (schema, unique ids, subject fits the grammar, file name equals
id, `redpanda-connect lint` on the `stream` block, credentials key exists in the named
SealedSecret's key list).

```yaml
id: mta-rail                 # lowercase, digits, dashes; also the KV key and stream id
name: MTA Metro-North and LIRR train positions
class: transit               # aircraft | vessel | transit | hazard | weather | space | intel | other
runner: redpanda-connect     # redpanda-connect | omos-ingest
mode: poll                   # poll | stream
schedule: 30s                # poll interval (default; the KV can override it live)
rate_limit: 2/min            # what upstream allows; the runner must stay under it
subject: evt.world.rail.{entity}   # must match the locked grammar; {entity} is filled per record
status_subject: tel.feed.mta-rail.status
credentials_ref: null        # or {secret: omos-feed-secrets, key: AISSTREAM_API_KEY}; never a value
license: "MTA developer terms: free, attribution required"
terms_url: https://new.mta.info/developers
key_required: none           # none | optional | required  (GEV's green/yellow/red column)
fallback: null               # another source id to switch on if this one fails
default_enabled: true        # used only when the KV has no entry yet
owner: zeus
stream:                      # runner: redpanda-connect only. Input and pipeline; no output.
  input: { generate: { interval: "${interval}", mapping: 'root = ""' } }
  pipeline: { processors: [ ... ] }
ingest_ref: null             # runner: omos-ingest only: its entry name in the ingest sources.yaml
```

Rules:
- **The controller adds the output**, built from `subject`: a JetStream publish with `Nats-Msg-Id`
  for dedup and the in-cluster NATS address and credentials. Source files never contain NATS URLs,
  TLS settings or passwords, and cannot publish to a subject other than the one in the manifest.
- The registry record in `services/schemas/registry/source.schema.json` (runtime state in
  `KV_SOURCE_REGISTRY`) stays as it is. The manifest is the Git half; the ids are the same.
- Bloblang mappings stay in the manifest, so a new simple source is **one file and one PR**.

### 3.2 `feed-control` bucket (the live switch)

A NACK `KeyValue` CR in `platform/nats/streams/kv-feed-control.yaml`: bucket `feed-control`, file
storage, R3, history 10, `preventDelete: true`.

- Key: the source id. Value: `{"enabled": true, "interval": "30s", "updated_by": "zeus", "reason": "..."}`.
- **Git holds the default, the KV holds the live value.** On start, a seeder writes the manifest
  default for any source that has no key yet (a KV *create*, so it never overwrites a live choice).
- Turning a feed off: `nats kv put feed-control mta-rail '{"enabled":false,"updated_by":"zeus"}'`.
  Later, `svc.feed.control` (section 3.5) and a Hermes tool do the same thing.
- History 10 gives an audit trail of who switched what.
- `feed-config` (empty, used by `feed-intel` code) is left alone for now. When `feed-intel` moves to
  manifests, it switches to `feed-control` and `feed-config` is retired. One bucket in the end.

### 3.3 Controller (no new image, no new pod)

The controller is **itself a Redpanda Connect stream** in the engine pod. It watches `feed-control`
with the `nats_kv` input and calls the engine's own streams API on `127.0.0.1:4195`:
on = `POST /streams/{id}` (then `PUT` if it already exists), off = `DELETE /streams/{id}`.
It reads the source's file from a mounted ConfigMap, fills in the interval and adds the output.

Verified with 4.103.1 and a throwaway local NATS (2026-11-08):
- The watch replays the current value of every key on start, so a restart restores the live state.
- On, off, interval change, key delete (= off) and six rapid flips all end in the right state.
- The pipeline must run with `threads: 1`, or changes are applied out of order.
- The API has no upsert: `PUT` on a missing stream is 404, `POST` on an existing one is 400.
- Stopping a stream takes 1 to 2 seconds (graceful drain).
- A bad id (rejected by a regex before any file is read) or a missing file is logged and skipped.

A Git change to a source file changes the ConfigMap hash (kustomize generator), which restarts the
pod; on start the controller replays the KV, so live on/off choices survive a deploy.

**Security (must ship with the controller):** the API fills in `${SECRET}` values on the server, and
`GET /streams/{id}` returns the config **with the secrets in plain text** (verified). So: put the
Service into Git without the API port (or delete the hand-applied one), and add a NetworkPolicy so
only the kubelet probe and, later, the metrics scraper reach 4195. Today any pod in the cluster can
read the AIS key and NATS password this way.

### 3.4 One shape for every position, CoT in one place

- Feeds publish **raw observations** with the doctrine 9.2 geo-envelope
  (`uid, type, lat, lon, hae, course, speed, callsign, stale, t, meta`) on the source subject:
  `evt.world.air.{id}`, `evt.world.sea.{id}`, `evt.world.rail.{id}` (#558), hazards on `intel.*`.
- **Only the entity service writes `ent.{type}.{id}.state`.** It fuses several sources about the same
  thing (one ship seen by two AIS feeds is one entity), which a feed cannot do on its own. This
  follows the entity proposal already filed; feeds publishing `ent.*` directly would give two writers.
- **Only the TAK projector builds CoT.** It already reads `ent.*`. Pipelines lose their inline CoT XML.
- Until the entity service and `ENT` are live, the legacy CoT feeds keep running beside the new
  ones (the #558 cut-over plan: 7 days in parallel, then switch off the old one).

### 3.5 Health and control on the bus

- Status: **`tel.feed.{id}.status`** every poll or every 30 s for streams:
  `{feed_id, enabled, interval, state, last_ok, last_error, items_last, items_total, lag_s}`.
  This is the subject `feed-intel` and #553 already use, so the brief's `evt.feed.{id}.health` is
  not needed as a second pattern. #553 adds `tel.feed.>` to `TEL`.
  Made by a `feed-health` stream in the engine pod that reads the engine's `/streams` and `/metrics`.
- State changes (switched on, switched off, failing, recovered) as events:
  `evt.feed.{id}.{enabled|disabled|failed|recovered}`. Needs `evt.feed.>` added to `EVT` (one line).
- Control RPC: `svc.feed.control` with `{op: list|enable|disable|set_interval, id, interval, reason}`.
  It only writes `feed-control`; the controller does the rest.
- Fallback: when a source reports `failed` N times in a row and its manifest names a `fallback`,
  the health stream switches the fallback on and says so in an event (GEV's fallback chains).

### 3.6 Time machine: archive pipeline to ClickHouse

A Redpanda Connect stream `archive-clickhouse`: `nats_jetstream` input with a durable consumer on
`WORLD` and `INTEL` (later `ENT`), batching, `sql_insert` with the ClickHouse driver into
`onemind.events`, and **ack only after the insert succeeds** (the Python sink drops the batch when
ClickHouse is down). It runs beside the Python sink on a different durable until row counts match
for 7 days, then `sink-clickhouse` goes to 0 replicas.

### 3.7 Two runners, one control plane

| | Redpanda Connect (feeds engine) | omos-ingest (Go, #553) |
|---|---|---|
| Good for | REST polls with a mapping, WebSocket, GTFS-RT, MQTT, the archive | OAuth or token refresh, paging, scraping, ETag and quotas, life-safety feeds |
| First sources | MTA (#558), ships, aircraft | NWS alerts |
| Switch | `feed-control` via the controller stream | `feed-control` (small change to #553: watch the bucket) |
| Status | `tel.feed.{id}.status` | `tel.feed.{id}.status` (already) |

No pod per source: 100 sources is about 100 streams in one or two engine pods, plus one ingest pod.

### 3.8 Memory and placement

The legacy pod runs near its 128 MiB limit on the memory-saturated hub-core pool. The engine runs in
the #558 Deployment `redpanda-connect-grammar` on the **ecosystem pool**, with requests 64 MiB and a
limit of 256 MiB to start. OpenSky's all-states response is several MB, so aircraft gets a bounding
box before it is switched on. If one pod grows past about 512 MiB, split by class (poll pod,
WebSocket pod), each still driven by the same bucket.

## 4. Build order

| Phase | Work | Live change |
|---|---|---|
| 1 (this stretch) | Schema and CI check; `feed-control` CR; controller and seeder streams; NetworkPolicy; MTA as the first manifest, in the engine pod; verify on/off live | New bucket. Engine pod from 0 to 1 replica, publishing `evt.world.rail.*` beside the legacy MTA (no TAK change) |
| 2 | `feed-health` stream, `evt.feed.*` events, `svc.feed.control`; ships and aircraft as manifests (default off) | Ships and aircraft only when switched on |
| 3 | Archive pipeline beside the Python sink; cut over after 7 days | New durable consumer |
| 4 | Entity service live, projector reads `ent.*` only, legacy CoT feeds off | TAK markers come from the projector |
| 5 | `feed-intel` collectors to manifests (simple ones to the engine, hard ones to omos-ingest); retire `feed-config` | Per collector, in parallel first |

Phase 1 is done when: the schema validates in CI, `feed-control` exists as a CR, MTA runs from its
manifest, and `nats kv put` turns it off and on again on the cluster within seconds with no restart.

## 5. Decisions for the owner

1. **Two runners** (Redpanda Connect for simple sources, omos-ingest #553 for hard ones) with one
   manifest format, one bucket and one status subject. *Recommended.* The alternative, Redpanda
   Connect only, means closing #553 and writing complex sources in Bloblang.
2. **Build on #558**: rebase and merge it, and make its `redpanda-connect-grammar` pod the
   controller-driven engine. *Recommended*: it already has the grammar MTA pipeline, tests and a
   cut-over plan. The alternative is to put the controller into the legacy pod, which is out of memory.
3. **Subjects**: status on the existing `tel.feed.{id}.status` instead of a new `evt.feed.{id}.health`;
   feeds publish raw observations and only the entity service writes `ent.*`. *Recommended.*
4. **Lock down the streams API** (Service without the API port, NetworkPolicy) in the same PR as the
   controller. *Recommended*; it closes an exposure that exists today.
5. **Restore the TAK-to-NATS wire.** Found by the cleanup session on 2026-11-07: the
   NatsBridgePlugin JAR is not in the `takserver-plugins` pod (no plugin has loaded since its restart
   on 2026-11-07 05:25 UTC), and `tak-nats-leaf` in `takserver` has 0 client connections, although
   its leaf link to the main bus is up. So TAK's own traffic (EUD positions, chat) does not reach
   NATS. The COT stream still fills, but only from the feeds publishing `tak.cot.inject`; with no
   subscriber on the TAK side, those markers probably do not reach TAK either (to check during
   phase 1: interest on `tak.cot.inject`). Two ways back:
   - **a. Adopt `tak-client`** (`platform/takserver/tak-client`, Argo app
     `platform-takserver-tak-client`, manual sync, never synced). It is a TAK streaming client on
     8089: TAK traffic out on `tak.cot.out` (inside the existing `tak.cot.>`), `tak.cot.inject` in.
     It is a direct client of the main bus, so `tak-nats-leaf` can be retired. It runs on the
     ecosystem pool (25m/32Mi requests). The owner mints its TAK client certificate and seals two
     Secrets (steps in its README). *Recommended*: Argo-managed, no change to the Helm-CLI release
     `takserver`, and it removes the last leaf inside the cluster.
   - **b. Restore the plugin**: put the JAR (`services/takserver/plugin/`) back into the plugins
     image or volume. This touches the Helm-CLI release `takserver` (no `helm upgrade` without the
     owner) and keeps `tak-nats-leaf`.
   Not both at once: each injects `tak.cot.inject` into TAK, so running both would show every
   marker twice. CLAUDE.md still calls the plugin the live wire and needs updating with whichever
   is chosen.

Namespace: `feeds`. The cleanup session renamed `omos` to `feeds` on 2026-11-08 (PR #760); paths
and Argo Application names are unchanged, and every manifest in this design targets `feeds`.

## 6. Risks and notes

- **Never run the Connect binary inside the engine pod** (for example `kubectl exec ... --version`).
  A second process in a 128 MiB container is OOM-killed and restarts the pod. Inspect the pinned
  image locally with Docker instead. (This happened once on 2026-11-08; the pod restarted cleanly.)
- Streams mode rejects `cache_resources` inside a stream config; shared caches go in a resources file
  loaded with `-r` (as `mta-resources.yaml` does today).
- `protobuf to_json` field names differ between Connect versions (#558 handles both cases).
- The NATS TLS posture (`skip_cert_verify: true`) is kept for now; the CA in `feeds` is #553's A0-6.
