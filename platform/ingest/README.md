# platform/ingest

`omos-ingest`, the collector framework (feature A0; code and tests in
`services/bridge/ingest`). Applied by the Argo Application `platform-ingest`
(`platform/argocd/platform-ingest.yaml`), manual sync, namespace `omos`.

| File | What |
|---|---|
| `configmap.yaml` | Feed registry (`sources.yaml`). First source: NWS alerts on `intel.weather.noaa.alert` |
| `deployment.yaml` | One replica, ecosystem pool preferred, 20m/48Mi requests, 250m/256Mi limits (sized from measured peaks, see `deployment.yaml`), non-root, read-only root FS |
| `service.yaml`, `servicemonitor.yaml` | Metrics on port 9464, scraped by kube-prometheus-stack (`release: monitoring`) |
| `prometheusrule.yaml` | `OneMindFeedStale`, `OneMindFeedNeverFetched`, `OneMindFeedNoLeader`, `OneMindIngestDown`, `OneMindIngestPublishRetrying`, `OneMindIngestDeadLetters` |

Secrets: none new. It reuses the existing SealedSecret `omos-nats-creds` (keys `user`,
`password`) that `feed-intel` uses. Keyed sources later use `omos-ingest-keys`
(`services/bridge/ingest/README.md`, "API keys"). No PVC, so no backup is needed: state is
in NATS KV (`ingest-state`), which the NATS stream backup covers.

## Before the first sync

1. Merge this PR. `platform-nats-streams` (auto-sync) creates KV `ingest-locks` (TTL 10 s)
   and `ingest-state` and adds `tel.feed.>` to `TEL`. Check, read-only:
   `nats kv info ingest-locks`, `nats kv info ingest-state`, `nats stream info TEL`.
2. The `build-omos-ingest` workflow pushes `ghcr.io/onemindos/omos-ingest:0.1.0`. Take the
   digest from its `omos-ingest-image` artifact and pin it in `deployment.yaml`
   (`image: ghcr.io/onemindos/omos-ingest:0.1.0@sha256:...`) in a follow-up PR.
3. Sync `platform-ingest` in Argo CD.

## Check after sync (read-only)

```bash
kubectl --context do-nyc3-omos-cloud -n omos logs deploy/omos-ingest | head
nats sub 'tel.feed.nws-alerts.status' --count 1     # fetchedAt, count, lastError null, stale false
nats stream subjects INTEL 'intel.weather.noaa.alert'
nats sub 'evt.sys.nws-alerts.rejected'              # dead letters with reasons, if any
```

In Prometheus: `ingest_feed_stale{source="nws-alerts"} == 0` and
`ingest_feed_leader{source="nws-alerts"} == 1`.

## Comparing with the legacy NOAA task (7 days, decision 0004 #18)

The legacy `feed-intel` task publishes `intel.weather.noaa.alert.{state}.{event}` with
`Nats-Msg-Id` = alert id; `omos-ingest` publishes `intel.weather.noaa.alert` with
`Nats-Msg-Id` = `nws:{alert id}`. Count distinct alert ids per day on each and compare. The
projector (`omos-tak-projector`, durable `sink-tak` on `intel.>`) reads flat payloads with
top-level `lat`/`lon`; the new envelope has its data under `data`, so the projector skips it
and TAK gets no duplicate markers. `sink-clickhouse` stores both (distinguishable by subject).

## Rollback

Scale to zero or delete the Application. Nothing else depends on it yet; the legacy task
is unchanged.
