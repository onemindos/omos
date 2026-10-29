# omos-ingest

The OneMind collector framework (feature A0). Plain Go on `nats.go`, as the ingestion
lane recommends (onemind-xr `docs/research/ingestion-platform.md`, "Recommendation"; decision
0004 #18). One process runs every source in a feed registry (`sources.yaml`) and delivers
records into NATS JetStream under the locked omos subject grammar.

First source: **NWS alerts** (`nws-alerts`, registry row F094), published on
`intel.weather.noaa.alert`. The legacy `feed-intel` NOAA task keeps running beside it
until the new pipeline matches it for 7 days (decision 0004 #18).

Deployment: `platform/ingest/` and the Argo Application `platform-ingest` (manual sync).
Image: `.github/workflows/build-omos-ingest.yml`.

## Subjects

| What | Subject | Stream |
|---|---|---|
| NWS alert (envelope, data per `services/schemas/intel/weather-alert.v1.schema.json`) | `intel.weather.noaa.alert` | `INTEL` |
| Feed freshness heartbeat, every 10 s, lease holder only (`services/schemas/ingest/feed-status.v1.schema.json`) | `tel.feed.{source}.status` | `TEL` (adds `tel.feed.>`) |
| Dead letter with the reason (`services/schemas/ingest/rejected.v1.schema.json`) | `evt.sys.{source}.rejected` | `SYS` |
| Lease (value = owning instance) | KV `ingest-locks`, key `{source}` | NACK CR `kv-ingest-locks.yaml` |
| Checkpoint (ETag, Last-Modified, cursor, seen set) | KV `ingest-state`, key `{source}` | NACK CR `kv-ingest-state.yaml` |

The loader refuses any data subject outside `intel.{domain}.{source}.{type}`,
`evt.world.{class}.{id}`, `det.{domain}.{type}.{id}` or `tel.*`. Collectors never publish
`ent.*`, `cmd.*`, `tak.cot.*` or `svc.*`. Every record is wrapped in the bus envelope
(`services/schemas/envelope/envelope.schema.json`); its `id` is a UUIDv4-shaped hash of the
dedup key, so a replay produces the same envelope. Grammar gaps are in
`docs/proposals/2026-10-29-a0-omos-ingest.md`.

## Reliability list: mechanism and test

Tests run with `go test -race ./...` against an embedded JetStream server
(`internal/natstest`). The chaos and alert runs are scripts (below).

| Item (run prompt §10) | Mechanism | Test |
|---|---|---|
| Publish with server ack, retries with backoff, `Nats-Msg-Id` dedup | `Runtime.Publish`: `js.PublishMsg` with `Nats-Msg-Id`, waits for the PubAck, capped exponential backoff, never drops. Errors no retry can fix (over `max_payload` or the stream's `max_msg_size`, publish permission denied) dead-letter the record with the reason instead of blocking the feed | `TestPublishAckAndMsgIDDedup`, `TestPublishRetriesUntilStreamExists`, `TestPublishBlocksThroughNATSRestart`, `TestPermanentPublishErrorsDeadLettered`, `TestPermissionViolationDeadLettered` |
| Restart or replay never duplicates | Checkpoint (seen set, ETag, cursor) in KV only after every record of a fetch is acked; a failed checkpoint read skips the poll | `TestRestartBeyondDedupWindowDoesNotDuplicate`, `TestCheckpointOnlyAfterAck`, `TestCheckpointReadErrorSkipsPoll`, `TestCheckpointCompareAndSet`, chaos |
| Exactly one owner per source | KV lease: `Create`, renew with `Update(revision)` every TTL/3. If no renewal succeeds within 4/5 of the TTL (the leader is cut off, alone or with the whole bus), a local timer stops the source before the lease can expire, and the instance competes for the lease again with `Create` after reconnecting. The client keeps no reconnect buffer (`ReconnectBufSize(-1)`), so a stopped leader cannot replay publishes. Checkpoint writes are fenced: the lease must still be held, and the KV write is a compare-and-set on the revision read. Clean shutdown hands the lease over | `TestLeaseSingleOwnerAndTakeover`, `TestLeaderCutOffAloneNeverPublishesAfterTakeover` (one instance cut off by a TCP proxy), `TestLeaseOutageStopsSourceAndReacquires`, `TestLeaseLostStopsSource`, `TestCheckpointCompareAndSet`, chaos (two instances, leader killed) |
| Schema validation, dead letters with the reason | JSON Schema from `services/schemas` per source; bad records and unparseable bodies go to `evt.sys.{source}.rejected`; a snapshot that keeps listing a bad record dead-letters it once; stream frames are recorded per occurrence | `TestMalformedRecordsToDeadLetter`, `TestWholeBodyParseErrorDeadLettered`, `TestNWSMalformed`, chaos |
| Upstream politeness | Per-source rate limiter, fixed-window quota, `If-None-Match`/`If-Modified-Since`, `Retry-After` on 429/503, jittered backoff on 429/5xx, no retry on other 4xx, poll backoff after failures; terms of use required per source | `TestFetchConditionalHeadersAnd304`, `TestFetch429HonoursRetryAfter`, `TestFetch5xxBackoffThenGiveUp`, `TestFetch403NoRetry`, `TestFetchRateLimit`, `TestQuotaWindow`, `TestFetchQuotaStopsRequests`, `TestNextDelayBackoff`, `TestConfigDefaultsAndValidation` |
| Freshness: `fetchedAt`, record count, last error; stale alert | Status board → `tel.feed.{source}.status` and `ingest_feed_*` metrics; `OneMindFeedStale` in `platform/ingest/prometheusrule.yaml` | `TestStatusBoardFreshness`, `TestFreshnessMetrics`, `TestFreshnessStatusAndUpstreamOutage`, `test/alerts/run.sh` (promtool), `test/chaos/stale-alert.sh` (live Prometheus) |
| Backpressure, limits, bounded memory | Response size cap, bounded websocket queue (`queueFull: block`, the default, loses no received frame; `drop-oldest` for position streams, counted in `ingest_dropped_total`), bounded seen ring kept under 512 KiB, SSE line cap, `GOMEMLIMIT` and container limits; KV buckets have `maxBytes` | `TestFetchBodyLimit`, `TestDropQueueBoundedDropsOldest`, `TestDropQueueBlockModeLosesNothing`, `TestKeepCurrentLastKeepsSnapshotInRing` |
| Survives NATS down, upstream down, malformed data, clock skew, full disk | Blocking publish and reconnect forever; last good count kept and marked stale; dead letters; dedup keys from upstream ids only; a full stream (`discard: new`) blocks without checkpointing | `TestPublishBlocksThroughNATSRestart`, `TestFreshnessStatusAndUpstreamOutage`, `TestKeysIndependentOfClock`, `TestFullStreamBlocksWithoutLossOrCheckpoint`, chaos |
| Replayable: fixtures drive golden tests; a stream rebuilds state | Parsers are pure; recorded fixtures in `testdata/fixtures` (NWS: one real api.weather.gov response from 2026-10-29, trimmed to 12 alerts) and golden files in `testdata/golden` | `TestNWSGolden`, `TestRSSAndAtom`, `TestCSVAndGeoJSONAndJSON`, `TestReplayRebuildsRecords` |
| Observable | JSON logs (`log/slog`) with `source` and `instance`; Prometheus metrics; ServiceMonitor | `TestFreshnessStatusAndUpstreamOutage` (log lines), `TestFreshnessMetrics`. OpenTelemetry traces are **not built** (follow-up) |
| Deployable through Argo | Pinned GHCR image (digest after first build), Argo Application, requests and limits, ecosystem pool, SealedSecret for keys | `kubeconform -strict` on `platform/ingest` |

## Delivery modes

| `kind` | Adapter | Resume | Notes |
|---|---|---|---|
| `rest_poll` | poll loop | ETag / Last-Modified + seen set | NWS uses this |
| `rss` | poll loop with the `rss` parser (RSS 2.0, Atom, W3C geo, GeoRSS) | same | |
| `file` | poll loop; `http(s)://` or `file://` (mtime as Last-Modified); `-once` for a CronJob | same | CSV, GeoJSON |
| `sse` | one long-lived stream | `Last-Event-ID` from the checkpointed cursor (saved after acks, at most once per second) | |
| `websocket` | reader → bounded queue → publisher; `queueFull: block` (default) or `drop-oldest` | none (upstream has no replay) | optional `subscribe` message, env-expanded |

Webhook push and MQTT stay with the existing bridges (`bridges/webhook`, `bridges/mqtt`)
until they adopt this core. Generic parsers: `nws-alerts`, `rss`, `geojson`, `csv`, `json`
(options: `prefix`, `idField`/`idFields`, `latField`, `lonField`, `timeField`, `timeLayout`,
`subject`, `subjectTemplate` such as `evt.world.sea.{mmsi}`).

## Adding a source

1. Add a registry entry to `platform/ingest/configmap.yaml`: `id` (one subject token),
   `kind`, `parser`, `url`, `subject`, `schema`, `interval`, `staleAfter`, `rateLimit`,
   optional `quota`, and `terms` (licence and URL are required). Add `registry:` with the
   onemind-xr `sources.yaml` id.
2. If no generic parser fits, add a pure parser in `internal/parsers` with a recorded
   fixture and a golden test (`go test ./internal/parsers -update` writes the golden file;
   review the diff).
3. Add or reuse a JSON Schema in `services/schemas`.

## API keys

Keys never go in Git. Reference them in the registry as `${VAR}` (headers, `subscribe`), and
put them in the SealedSecret `omos-ingest-keys` in namespace `omos`, which the Deployment
loads with `envFrom` (optional). With cluster access, the owner seals it:

```bash
kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > /tmp/pub-cert.pem
kubectl create secret generic omos-ingest-keys -n omos --from-literal=EXAMPLE_API_KEY=<value> \
  --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/ingest/omos-ingest-keys.sealedsecret.yaml
```

NWS needs no key.

## Run locally

```bash
nats-server -js -p 4222 &
go run ./cmd/omos-ingest -config sources.local.yaml -schemas ../../schemas \
  -nats nats://127.0.0.1:4222 -create-buckets -listen 127.0.0.1:9464
curl -s 127.0.0.1:9464/metrics | grep ingest_feed
```

`-create-buckets` is for local runs only; in the cluster the buckets are NACK CRs.
`-check` validates the registry and schemas and exits; it fails when no enabled source is
found. CI runs it on `sources.local.yaml` and on `data['sources.yaml']` extracted from
`platform/ingest/configmap.yaml`.

## Chaos and alert tests

- `test/chaos/chaos.sh [out]`: two collectors, an upstream stub (50 NWS-shaped alerts/s
  and 50 SSE events/s); `kill -9` nats-server at 20 s (back at 35 s), upstream 503 from
  50 to 65 s, `kill -9` the leading collector at 80 s (back at 95 s), malformed records from
  110 to 120 s; `test/chaos/verify` then checks 0 duplicates, 0 missing and a dead letter
  for every malformed record.
- `test/chaos/stale-alert.sh [out]`: a real Prometheus with the unchanged alert rules; a
  3-minute upstream outage must take `OneMindFeedStale` to firing and back.
- `test/alerts/run.sh`: `promtool check rules` and `promtool test rules`.

## Known limits

- Quota counters are in memory; a restart starts a new window. Set quotas below the
  upstream's real limit.
- WebSocket sources have no replay: frames sent while no instance is connected are lost,
  and a full queue drops the oldest frame (counted).
- Zone-based NWS alerts (most of them) carry `affectedZones` and UGC/SAME codes, not a
  point. Resolving zone shapes is left to the entity service, so each poll is one upstream
  request. The legacy task resolved them with extra requests and fell back to state
  centroids or `0,0`.
- No OpenTelemetry traces yet.
