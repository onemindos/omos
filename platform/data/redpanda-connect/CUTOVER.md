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

Three `replicas` fields in Git, each changed by its own reviewed PR. Argo applies them
(`platform-redpanda-connect` and `platform-omos` both sync automatically).

| Field | File | Merged value |
|---|---|---|
| `redpanda-connect-grammar` replicas | `platform/data/redpanda-connect/deployment-grammar.yaml` | `0` (off) |
| `redpanda-connect` replicas (legacy `feed-mta` CoT) | `platform/data/redpanda-connect/deployment.yaml` | `1` (on) |
| `feed-satellites` replicas (legacy satellite CoT) | `platform/omos/feeds.yaml` | `1` (on) |

`scripts/tests/test_gitops_feed_grammar.py` fails CI if a legacy output is set to 0 while the
grammar Deployment is still 0, and if any pipeline is named by two Deployments.

## Stage 0: merge this PR (no live change)

What changes live: the `redpanda-connect-pipelines` ConfigMap gains two keys and loses
`feed-hazards.yaml`; a `redpanda-connect-grammar` Deployment exists with 0 pods. The legacy pod
keeps its args and is not restarted (a ConfigMap change does not restart it).

Check (read-only):

```bash
kubectl --context do-nyc3-omos-cloud -n omos get deploy redpanda-connect redpanda-connect-grammar feed-satellites
# expected: redpanda-connect 1/1, redpanda-connect-grammar 0/0, feed-satellites 1/1
kubectl --context do-nyc3-omos-cloud get --raw \
  "/api/v1/namespaces/omos/pods/$(kubectl --context do-nyc3-omos-cloud -n omos get pod -l app.kubernetes.io/name=redpanda-connect -o name | cut -d/ -f2):4195/proxy/streams"
# expected: only {"feed-mta":{"active":true,...}}
```

## Stage 1: run both side by side

Prerequisites:
1. The owner has ruled on A1-P1 (`rail` world class, BUS-P6). The grammar MTA pipeline publishes
   `evt.world.rail.*`, so it does not start before the ruling.
2. The `omos-nats-creds` user may publish `evt.world.>`, `intel.>`, `evt.sys.>`, `tel.feed.>` and
   use the JetStream API. `feed-intel` already publishes `intel.>` with JetStream acks using the
   same Secret; the rest is confirmed by the log check below.

PR: set `redpanda-connect-grammar` `replicas: 1`. Nothing else.

Checks after Argo syncs (the `nats` commands need a NATS user; run them from the tailnet against
`omos-nats` or from a pod with the CLI):

```bash
kubectl --context do-nyc3-omos-cloud -n omos get pod -l app.kubernetes.io/name=redpanda-connect-grammar -o wide
# expected: 1/1 Running on the onemind-ecosystem-cpu node
kubectl --context do-nyc3-omos-cloud -n omos logs deploy/redpanda-connect-grammar | grep -iE "error|permission" | head
# expected: no output (a "Permissions Violation" line means prerequisite 2 is not met: go back to replicas 0)
nats sub 'tel.feed.>' --count 4
# expected: tel.feed.mta.status every 30 s with stale:false, published > 0, lastError null;
#           tel.feed.celestrak.status once per 4 h fetch (once right after the pod starts)
nats stream subjects WORLD 'evt.world.rail.>' | head
nats stream subjects INTEL 'intel.space.>'
# expected: one subject per train running now; intel.space.celestrak.tle present
nats sub 'evt.sys.*.rejected' --count 1 --timeout 60s
# expected: usually nothing; any record carries a reason
```

Compare with the legacy feed over the same window: trains in `evt.world.rail.>` against the
`MTA-MNR-*`/`MTA-LIRR-*` uids the legacy pipeline still puts on `tak.cot.inject`. Keep Stage 1
running for 7 days with no errors and matching coverage (ingestion migration plan, decision 0004 #18).

Rollback: revert the Stage 1 PR (replicas back to 0).

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
kubectl --context do-nyc3-omos-cloud -n omos get deploy redpanda-connect feed-satellites
# expected: 0/0 for both
nats sub tak.cot.inject --count 20 --timeout 90s | grep -cE 'MTA-(MNR|LIRR)-|SAT-'
# expected: 0 (only projector output, uids OM.*, remains)
nats sub 'tel.feed.mta.status' --count 1
# expected: still stale:false; the grammar path is unaffected
```

Rollback: revert the Stage 2 PR (both replicas back to 1). The grammar pipelines keep running.

## Stage 3: delete the legacy pipeline (later PR)

After Stage 2 has held for 7 days: delete `feed-mta.yaml` from the ConfigMap and `pipelines/`,
the `redpanda-connect` Deployment, the `feed-satellites` Deployment and `feed_satellites.py`, and
update `scripts/tests/test_gitops_feed_grammar.py`. When MTA and CelesTrak move to `omos-ingest`
(ingestion migration plan steps 7 and 8), the grammar pipelines here are retired the same way.

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
