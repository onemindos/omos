# Feed sources

One YAML file per data source. The file name is the source id (`mta-rail.yaml` has `id: mta-rail`).
The format is checked against `services/schemas/feeds/source-manifest.v1.schema.json`.
Design: `docs/architecture/feed-platform.md`.

## What a file says

Where the data comes from, which NATS subject it publishes to (only `evt.world.`, `intel.`, `tel.`
or `det.`), how often it runs, its licence, which secret keys it needs (names only, never values),
and for the `redpanda-connect` runner the whole stream config under `stream:`. Put the literal
`__INTERVAL__` where the poll interval goes. The controller fills it in.

## Add or change a source

1. Add or edit `platform/data/sources/<id>.yaml`.
2. Run `python scripts/gen_feed_sources.py`. It rewrites
   `platform/data/redpanda-connect/feed-sources.configmap.yaml` (do not edit that file by hand).
3. Open a PR. CI checks the schema, the subjects and that the generated file is up to date.

`default_enabled` is only the starting value. It is written to the bucket once, when the bucket has
no entry for the source. After that the bucket is the truth.

## Switch a feed on, off or change its speed

```
nats kv put feed-control mta-rail '{"enabled":true,"interval":"30s","updated_by":"you","reason":"why"}'
nats kv put feed-control mta-rail '{"enabled":false,"updated_by":"you","reason":"why"}'
nats kv get feed-control mta-rail
nats kv history feed-control mta-rail      # who changed it, and when
```

Deleting the key also switches the feed off. Keys that are not a known source id are ignored.

## First deploy of the bucket (once)

NACK needs a server URL on each KeyValue CR, and that URL is set out of band, never in Git
(`CLAUDE.md`, NATS Streams). After Argo first creates `kv-feed-control`, copy `spec.servers` from
another live KeyValue CR in `nats` (for example `kv-world-state`) without printing it, then check
`kubectl -n nats get keyvalue kv-feed-control` reports Ready. Until the bucket exists the engine pod
crash-loops on "bucket not found"; delete the pod once the bucket is Ready.

## What the engine pod needs for these sources

- ConfigMap `feed-sources` (generated) mounted at `/config/sources`.
- ConfigMap `feed-control-pipelines` (the controller and the seeder) mounted at `/config/control`.
- `mta-rail` also needs the GTFS proto file: ConfigMap `redpanda-connect-proto`, key
  `gtfs-realtime.proto`, mounted at `/proto/gtfs-realtime.proto`. It needs no `-r` resource file
  (the dedupe cache of the old `feed-mta` pipeline is not used).
- Env `NATS_USER` and `NATS_PASSWORD` from the Secret `omos-nats-creds` (keys `user`, `password`).

The engine pod `redpanda-connect-grammar` (`platform/data/redpanda-connect/deployment-grammar.yaml`)
has all of this wired: it starts only the controller and the seeder, and the pod annotation
`onemind.io/feed-config-sha256` (kept current by `scripts/gen_feed_sources.py`) restarts it when a
manifest changes.
