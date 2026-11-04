# Owner proposals from feature A0 (omos-ingest)

Date: 2026-10-29. Filed under onemind-xr decision 0004 #8: where the locked grammar or the
stream set has no exact place, this PR uses the closest existing pattern and records it here
for the owner. None of these changes a live subject that another service publishes. IDs refer
to the proposal register in onemind-xr `docs/design/BUS-DESIGN.md` §10.1, so the A24 PR can
file them once.

| # | Gap | What this PR does | Proposed ruling |
|---|---|---|---|
| A0-1 (BUS-P9) | Dead-letter subject. The ingestion lane used `evt.feed.{source}.rejected`, which no stream binds | Publishes dead letters on `evt.sys.{source}.rejected`, stored by the existing `SYS` stream (`evt.sys.>`) | Ratify `evt.sys.{source}.rejected` as the dead-letter subject for every collector |
| A0-2 (BUS-P1, part) | Feed freshness heartbeats `tel.feed.{source}.status` are published today by `feed-intel` but stored by no stream | Adds `tel.feed.>` to `TEL` (`platform/nats/streams/tel.yaml`); the heartbeat is a core publish | Ratify `tel.feed.{source}.status` and its `TEL` binding. Note: once merged, the legacy `feed-intel` heartbeats (7 feeds, every 30 s) are stored too, under `TEL`'s 168 h `maxAge` |
| A0-3 (BUS-P3) | NWS alerts: the legacy task publishes 6-token `intel.weather.noaa.alert.{state}.{event}` | The new collector publishes the 4-token `intel.weather.noaa.alert` (feed registry row F094), state and event in the payload. The legacy subject is unchanged | Keep 4 tokens for new collectors; retire the legacy 6-token subject when `feed-intel`'s NOAA task is removed after the 7-day comparison |
| A0-4 | Checkpoint bucket. BUS-DESIGN §5 lists KV `ingest-locks` but no bucket for collector checkpoints | Adds KV `ingest-state` (history 1, no TTL, 64 MiB) beside `ingest-locks` (history 1, TTL 10 s, 8 MiB) | Add `ingest-state` to BUS-DESIGN §5 and the per-node budget (64 MiB × R3) |
| A0-5 | Collector identity. BUS-DESIGN §7.2 gives the collector group its own user (`intel.>`, `evt.world.>`, `tel.feed.>`, `evt.sys.*.rejected`, KV `ingest-locks`/`ingest-state`) | Reuses the shared `omos-nats-creds` that `feed-intel` uses, because per-service users do not exist yet | Create the collector user in the identities PR and switch `platform/ingest/deployment.yaml` to it |
| A0-6 | NATS client TLS. The hub client port requires TLS (`platform/nats/values.template.yaml`, `config.nats.tls`), and the CA is not available in namespace `omos` | Sets `NATS_TLS_INSECURE=true`, the same posture as `feed-intel` and `omos-tak-projector` (`ssl.CERT_NONE`) | Publish the NATS CA certificate to `omos` (it is public data) and switch to `NATS_CA` |
| A0-7 | `INTEL` duplicate window is 2 min, shorter than how long upstreams re-list items | Unchanged. Collector-side checkpoints (seen set after acks) prevent replays without a longer window (tested beyond the window) | Keep BUS-DESIGN's 1 h `INTEL` window and `maxBytes` caps in the stream caps PR, not here |
| A0-8 | WebSocket queue when full. The A0 acceptance test says "loses no data", but a stream without replay that is read slower than it sends must either stop reading (the upstream may disconnect the reader, which loses data upstream) or drop frames | Per-source `queueFull`: `block` (default: no received frame is dropped) or `drop-oldest` (counted in `ingest_dropped_total`). No websocket source is configured yet | Use `block` for feeds that carry records; allow `drop-oldest` only for position streams where a newer position supersedes an older one (AIS, ADS-B), and list them in the feed registry |

Information only (no ruling needed): every record is wrapped in the bus envelope
(`services/schemas/envelope`), with `source` = `ingest.{source id}` and two extra fields that
the envelope allows (`additionalProperties: true`): `fetchedAt` and `msgId`.
