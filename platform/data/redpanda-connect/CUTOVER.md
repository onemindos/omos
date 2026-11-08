# Feed grammar cut-over (FEATURES A1)

`feed-mta` and `feed-satellites` publish CoT XML straight to `tak.cot.inject`. Their
replacements publish raw observations on grammar subjects and leave CoT to the entity service
and the projector:

| Feed | Legacy output (stays on until Stage 2) | Grammar output | Stream |
|---|---|---|---|
| MTA Metro-North + LIRR | `feed-mta.yaml` in the `redpanda-connect` pod → `tak.cot.inject` | `feed-mta-grammar.yaml` → `evt.world.rail.mta-{mnr\|lirr}-{trip}` (doctrine §9.2 envelope, `type: ground.rail`) | `WORLD` |
| CelesTrak `stations` | `feed-satellites` Deployment (`platform/omos/feeds.yaml`) → `tak.cot.inject` | `feed-satellites-grammar.yaml` → `intel.space.celestrak.tle` (one TLE per satellite, no position) | `INTEL` |
| both | — | `tel.feed.{mta\|celestrak}.status` per poll; `evt.sys.{mta\|celestrak}.rejected` | core / `SYS` |

`feed-hazards` is removed (it duplicated the USGS and NWS collectors in `feed-intel` and was not
running). `feed-aircraft` and `feed-ships` stay off; they return only through `omos-ingest` with
an AOI box and quota.

Owner proposals this depends on: `docs/proposals/2026-10-29-feed-grammar-cutover.md`.

## The switch

MTA no longer switches with a Deployment field. `redpanda-connect-grammar` is now the feed engine
(`docs/architecture/feed-platform.md`, phase 1): it runs `feed-controller` and `feed-seeder`, and the
controller starts or stops each source stream from the NATS KV bucket `feed-control`. The MTA stream
comes from the manifest `platform/data/sources/mta-rail.yaml`, not from `feed-mta-grammar.yaml`.

| Switch | Where | Merged value |
|---|---|---|
| MTA grammar feed (`evt.world.rail.*`) | key `mta-rail` in KV bucket `feed-control` | `{"enabled":false}` (seeded from `default_enabled: false`) |
| `redpanda-connect-grammar` replicas (the engine, not a feed) | `platform/data/redpanda-connect/deployment-grammar.yaml` | `1` |
| `redpanda-connect` replicas (legacy `feed-mta` CoT) | `platform/data/redpanda-connect/deployment.yaml` | `1` (on) |
| `feed-satellites` replicas (legacy satellite CoT) | `platform/omos/feeds.yaml` | `1` (on) |

Switching MTA needs no Git change, no Argo sync and no restart (`nats` from the tailnet):

```bash
nats kv put feed-control mta-rail '{"enabled":true,"interval":"30s","updated_by":"you","reason":"stage 1"}'
nats kv put feed-control mta-rail '{"enabled":false,"updated_by":"you","reason":"rollback"}'
nats kv history feed-control mta-rail     # who switched it, and when
```

The satellites grammar pipeline (`feed-satellites-grammar.yaml`) is not started by the engine: phase 1
is MTA only, and satellites stays off until it has a manifest in `platform/data/sources/`. The legacy
`feed-satellites` Deployment keeps running.

A change to a manifest in `platform/data/sources/` changes the pod annotation
`onemind.io/feed-config-sha256` (written by `scripts/gen_feed_sources.py`), which restarts the engine
pod. On start the controller replays the bucket, so live on/off choices survive the restart.

`scripts/tests/test_gitops_feed_grammar.py` fails CI if a legacy output is set to 0 while the
grammar Deployment is still 0, and if any pipeline is named by two Deployments.

## Stage 0: merge (new bucket, engine pod up, MTA still off)

What changes live: the NACK controller creates the KV bucket `feed-control` (R3, from
`platform/nats/streams/kv-feed-control.yaml`); ConfigMaps `feed-sources` and `feed-control-pipelines`
are created; a NetworkPolicy denies all ingress to both Redpanda Connect pods (streams API lock-down,
decision 4); `redpanda-connect-grammar` goes from 0 to 1 pod on the ecosystem pool, runs the controller
and the seeder, and the seeder writes `mta-rail` as disabled. Nothing new is published. The legacy
pod keeps its args and is not restarted.

Check (read-only):

```bash
kubectl --context do-nyc3-omos-cloud -n feeds get deploy redpanda-connect redpanda-connect-grammar feed-satellites
# expected: all 1/1; redpanda-connect-grammar on the onemind-ecosystem-cpu node
kubectl --context do-nyc3-omos-cloud -n feeds logs deploy/redpanda-connect-grammar | grep -E 'feed-seeder|feed-control' | head
# expected: feed-seeder id=mta-rail processed err=none
nats kv get feed-control mta-rail
# expected: {"enabled":false,...,"updated_by":"seeder"}
kubectl --context do-nyc3-omos-cloud -n feeds get networkpolicy redpanda-connect-no-ingress
```

Do not `kubectl exec` into either pod (a second process OOM-kills it). The streams API is no longer
reachable from other pods; read the engine state from its logs.

## Stage 1: run both side by side

Prerequisites:
1. The owner has ruled on A1-P1 (`rail` world class, BUS-P6). The grammar MTA stream publishes
   `evt.world.rail.*`, so it does not start before the ruling.
2. The `omos-nats-creds` user may publish `evt.world.>`, `intel.>`, `evt.sys.>`, `tel.feed.>` and
   use the JetStream API. `feed-intel` already publishes `intel.>` with JetStream acks using the
   same Secret; the rest is confirmed by the log check below.

Switch it on (no PR):

```bash
nats kv put feed-control mta-rail '{"enabled":true,"interval":"30s","updated_by":"you","reason":"stage 1"}'
```

Checks (the `nats` commands need a NATS user; run them from the tailnet against `omos-nats` or from
a pod with the CLI):

```bash
kubectl --context do-nyc3-omos-cloud -n feeds logs deploy/redpanda-connect-grammar | grep -iE "error|permission" | head
# expected: no output (a "Permissions Violation" line means prerequisite 2 is not met: switch it off)
nats sub 'tel.feed.>' --count 4
# expected: tel.feed.mta.status every 30 s with stale:false, published > 0, lastError null
nats stream subjects WORLD 'evt.world.rail.>' | head
# expected: one subject per train running now
nats sub 'evt.sys.*.rejected' --count 1 --timeout 60s
# expected: usually nothing; any record carries a reason
```

Compare with the legacy feed over the same window: trains in `evt.world.rail.>` against the
`MTA-MNR-*`/`MTA-LIRR-*` uids the legacy pipeline still puts on `tak.cot.inject`. Keep Stage 1
running for 7 days with no errors and matching coverage (ingestion migration plan, decision 0004 #18).

Rollback: `nats kv put feed-control mta-rail '{"enabled":false,"updated_by":"you","reason":"rollback"}'`.
The stream stops within a few seconds. Setting the Deployment to 0 replicas also stops it, but the
bucket is the switch.

## Stage 2: turn the legacy CoT outputs off

Prerequisites:
1. Stage 1 ran clean for 7 days.
2. The entity service (omos PR #552) consumes the new subjects: its `World` adapter accepts the
   `rail` class and its `Intel` adapter propagates TLEs (FEATURES A12), or the owner accepts that
   trains and satellites are off TAK until then. Today #552 skips `evt.world.rail.*`
   (`SkipUnsupported`) and position-less `intel.*` records (`SkipNoPosition`).
3. For context: nothing delivers `tak.cot.inject` into TAK Server today (FEATURES A21), so this
   stage removes no picture from ATAK. The only live reader of that subject is the Hermes leaf
   (`tak.cot.>`).

PR (one PR, both fields, so the two feeds cut over together):
- `platform/data/redpanda-connect/deployment.yaml`: `redpanda-connect` `replicas: 0`
- `platform/omos/feeds.yaml`: `feed-satellites` `replicas: 0`

Checks:

```bash
kubectl --context do-nyc3-omos-cloud -n feeds get deploy redpanda-connect feed-satellites
# expected: 0/0 for both
nats sub tak.cot.inject --count 20 --timeout 90s | grep -cE 'MTA-(MNR|LIRR)-|SAT-'
# expected: 0 (only projector output, uids OM.*, remains)
nats sub 'tel.feed.mta.status' --count 1
# expected: still stale:false; the grammar path is unaffected
```

Rollback: revert the Stage 2 PR (both replicas back to 1). The grammar feed keeps running.

## Stage 3: delete the legacy pipeline (later PR)

After Stage 2 has held for 7 days: delete `feed-mta.yaml` from the ConfigMap and `pipelines/`,
the `redpanda-connect` Deployment, the `feed-satellites` Deployment and `feed_satellites.py`, and
update `scripts/tests/test_gitops_feed_grammar.py`. When MTA and CelesTrak move to `omos-ingest`
(ingestion migration plan steps 7 and 8), the grammar sources here are retired the same way.

## Test on a build machine

```bash
platform/data/redpanda-connect/test/replay.sh [evidence-dir]
```

It runs the pipelines exactly as they are in the ConfigMap against a throwaway `nats-server` and
recorded fixtures (MTA GTFS-RT and CelesTrak responses recorded 2026-10-29, trimmed, plus one bad
position and one bad checksum), captures everything with `nats sub '>'`, and checks: only grammar
subjects, nothing on `tak.cot.*`, every train once in `WORLD` across repeated polls (dedup), every
valid TLE once in `INTEL`, rejects with reasons, status per poll. A negative control runs the
legacy `feed-mta` pipeline and confirms the capture sees its `tak.cot.inject` traffic. Needs
`nats-server`, `nats`, `redpanda-connect` and Python with PyYAML. It was run with Redpanda Connect
4.115.1; the cluster image is 4.103.1, so the pipelines read GTFS field names in both the camelCase
and snake_case spellings that protobuf `to_json` produces across versions.
