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

## What the engine pod needs for these sources

- ConfigMap `feed-sources` (generated) mounted at `/config/sources`.
- ConfigMap `feed-control-pipelines` (the controller and the seeder) mounted at `/config/pipelines`.
- `mta-rail` also needs the GTFS proto file: ConfigMap `redpanda-connect-proto`, key
  `gtfs-realtime.proto`, mounted at `/proto/gtfs-realtime.proto`. It needs no `-r` resource file
  (the dedupe cache of the old `feed-mta` pipeline is not used).
- Env `NATS_USER` and `NATS_PASSWORD` from the Secret `omos-nats-creds` (keys `user`, `password`).

None of this runs until the engine pod is wired to these ConfigMaps (a later change).
