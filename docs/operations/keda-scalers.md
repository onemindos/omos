# KEDA scalers: what is safe today

Written 2026-10-23 from the repo and read-only `kubectl` against `do-nyc3-omos-cloud`. KEDA 2.21.0 runs in
namespace `keda` (Argo app `ecosystem-keda`). Before this change there were 0 ScaledObjects.

## What is deployed

`platform/omos/scaledobject-sink-clickhouse.yaml`: NATS JetStream trigger on `INTEL` / `sink-clickhouse`,
**min 1, max 1** (observe-only). It proves the monitoring path and exposes the lag as an HPA metric. It does
not change replica count.

The NATS monitoring endpoint (`nats-headless.nats.svc.cluster.local:8222`) needs no credentials, so no
TriggerAuthentication is included. `platform/nats/network-policy.yaml` admits the `keda-operator` pod
(`app: keda-operator` in namespace `keda`) to port 8222 on the NATS servers; update that rule if KEDA moves or is
relabelled. If monitoring ever gets auth, add a `TriggerAuthentication` and pass `natsServerMonitoringEndpoint`
via `authParams`.

## Who consumes JetStream

Live check of `/jsz?consumers=true`: 38 streams, **2 consumers**, both on `INTEL`, both **push** consumers
with 0 pending and 0 redelivered. The repo has no NACK `Consumer` resources (streams only), so consumers are
created by the workloads themselves.

| Workload (ns) | Consumer | Also holds | Verdict |
|---|---|---|---|
| `sink-clickhouse` (omos) | push durable `sink-clickhouse`, ack none | core sub `ent.*.*.state` | min 1 / max 1 |
| `omos-tak-projector` (omos) | push durable `sink-tak` | core subs `ent.*.*.state`, `cmd.projector.*` | keep 1 |
| `onemind-core`, `onemind-ingestion` | none (core subscribe, publishes to JS) | | not queue-driven |
| `hermes-legacy` | none (Synadia agent protocol on `agents.*`, core request/reply) | | not queue-driven |
| `takml-vision-worker` | none: it **publishes** `fabric.vision.*.detect` | | see below |

Push consumers with no deliver group allow one subscriber, so max replicas is 1. `services/subscriber.py`
uses a queue group (scalable) but nothing in the cluster runs it.

## Idle workloads

| Workload | Idle? | Why it cannot scale to zero yet |
|---|---|---|
| `takml-vision-worker` (takml) | Mostly: 3m CPU, 265Mi | `takml-server` has a **required** podAffinity on `app=takml-vision-worker` (shared RWO `takml-models` PVC). At 0 replicas, a `takml-server` restart cannot schedule. Fix the affinity first. The worker also reads a static `FEEDS` list and retries every 10s; it has no queue. A real trigger would be MediaMTX path readiness (`onetak-media:4000/v3/paths/list`, KEDA `metrics-api`). |
| `nodeodm` (onetak) | Yes, except during drone jobs | WebODM selects a processing node when a task is created. At 0, the node shows offline and job submission is likely to fail, and no signal exists before submission. Needs a test with the current WebODM image. Candidate signal: a `postgresql` scaler on `webodm-db` counting queued tasks, once WebODM is confirmed to hold tasks while no node is online. |
| `webodm` / `webodm-broker` / `webodm-db` | UI + DB | HTTP service and state; min 1. (The Celery queue in Redis is not a useful signal for NodeODM.) |

## Candidate scalers

| Candidate | Scaler | Needs before enabling |
|---|---|---|
| `nodeodm` 0..1 | `postgresql` (queued WebODM tasks) or `cron` as a fallback | Confirm WebODM queues tasks with no node online; keep PVC `nodeodm-data` (RWO, Recreate); give scale-up time (image pull + start). |
| `takml-vision-worker` 0..1 | `metrics-api` on MediaMTX paths | Remove the required podAffinity from `takml-server` (or move both to one pod); confirm FEEDS paths. |
| Any NATS consumer | `nats-jetstream` | A queue-group or pull consumer. Today's push consumers are single-subscriber and always busy. |

## Argo CD interaction

Argo apps here use `selfHeal: true`, and the Deployments hard-code `replicas`. Any ScaledObject that can
change replicas away from the manifest value needs `/spec/replicas` in the owning Application's
`ignoreDifferences`, or Argo will scale it back. The observe-only object avoids this because min = max =
manifest replicas.

## Verify

```
kubectl -n omos get scaledobject sink-clickhouse      # READY=True
kubectl -n omos get hpa                               # keda-hpa-sink-clickhouse, TARGETS shows lag
kubectl -n keda logs deploy/keda-operator | grep -i sink-clickhouse
```
