# OneMind plugins (FEATURES K6)

One manifest, `onemind-plugin.yaml`, for every host: omos collectors, the Unreal client, the L0
web client, OneMindATAK and Hermes. A host ignores the blocks it cannot run. Schema:
`services/schemas/plugins/onemind-plugin.v1.schema.json`. Design: onemind-xr
`docs/research/video-and-plugins.md` "Plugin model" (worldwideview was read for design only; no
code copied).

| Kind | Block | What it is |
|---|---|---|
| `source` | `source` | a collector: an omos-ingest (A0) source entry (config only) or an own pinned GHCR image |
| `layer` | `layer` | declarative rendering of a subject or snapshot (entity type, fields, caps, LOD) |
| `action` | `action` | a consumer of an approved `cmd.*`; always policy + audit + `confirm: human` |
| `skill` | `skill` | a Hermes plugin exposing tools that call `svc.*` RPC |
| `video` | `video` | binds a MediaMTX path to its entity |

CI (`quality.yml`, step "Plugin manifests") runs `scripts/validate_plugin_manifests.py`:
schema; every `bus.*` subject against the locked grammar; no plugin publishes `ent.*`, `cmd.*`,
`evt.so.action.*`, `tak.cot.inject` or `tak.cot.mesh.tx` (their single writers are the entity,
policy and projector services); each block's subject is declared in `bus` (NATS permissions are
generated from `bus`); drone/robot actions only on `sim-*`; community plugins are declarative.
`scripts/tests/fixtures/plugins/` holds manifests that must fail.

Examples, one per kind: `examples/usgs-quakes` (source + layer: a new layer with no core
change), `examples/drone-orbit` (action), `examples/video-lookup` (skill), `examples/sim-x1-camera`
(video). The `kind: quake` filter in the USGS layer depends on the record type omos-ingest's
`geojson` parser writes; check it when the source is enabled.

Not in this PR (owner proposal K-P9): the plugin registry KV bucket and `svc.plugins.list`.
