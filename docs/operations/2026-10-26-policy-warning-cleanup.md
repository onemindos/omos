# Image and resource warning repair — 2026-10-26

This is the user-requested warning cleanup, not a policy enforcement change.
Telemetry was separately repaired: [receipt](2026-10-26-telemetry-pull-secret.md).

## Current workloads and budgets

The remaining real resource gaps were TAK internal NATS's config reloader, TAK
PostgreSQL and the existing retention CronJob. They were first adopted without
runtime changes after [recovery preflight](2026-10-26-tak-policy-preflight.md).
Whole-release Helm operations remain frozen to prevent a second writer.

Seven-day Prometheus maximum working sets were about 4.6Mi for the reloader and
1026.8Mi for PostgreSQL; PostgreSQL RSS peaked at96Mi, with the larger working set
including cache. The reloader now requests10m CPU/16Mi memory and limits memory to
64Mi. PostgreSQL keeps its250m CPU request, raises memory reservation from256Mi to
768Mi and receives a2Gi memory limit, leaving headroom above observed working set.
The psql retention client requests25m/64Mi and limits memory to256Mi; its SQL,
credentials, schedule and concurrency behavior are unchanged.

Argo sync waves roll the single internal NATS instance first, then PostgreSQL,
then update the retention template. A short TAK reconnect/database interruption
can occur. Images, data paths, existing claims and credential references stay the
same. Source Helm values mirror the adopted budgets but are not applied with Helm.
Standing memory reservations rise by528Mi, plus64Mi while a retention job runs;
no node-count change is included.

## Historical policy events

Current running app images already have explicit versions or digests. Most image
warnings were reissued against retained old zero-replica revisions. Keep those
revisions for rollback; do not delete their history or weaken admission checks.

Only the image/resource policies gain a background-history guard. Real admission
requests have a nonempty AdmissionReview UID and always run the original checks,
including updates of retired objects. Missing/empty UID is how Kyverno1.19.1's
background scan context differs; both contexts use operationCREATE, so operation
alone cannot safely distinguish them.

Background checks skip only observed zero-replica Deployment-owned ReplicaSets,
finished CronJob-owned Jobs without active/terminating pods, and terminal Pods
with recognized controller owners. Active, draining, newly created, unowned,
unknown/stale-status resources and scaled-to-zero top-level Deployments remain
assessed. The privileged/hostPath/namespace rules are unchanged. Audit mode,
failurePolicyIgnore, background scanning and automatic controller rules remain.

The exact deployed1.19.1 CLI tests46 explicit fixtures in four contexts (background,
empty UID, admissionCREATE and admissionUPDATE), including orphan/unknown cases.
The same checksum-pinned CLI is installed in CI so those tests cannot silently
skip. Independent before/after admission comparison found identical results for
all172 resource/rule evaluations per operation. An engine upgrade must revisit
these guards and fixtures. Old Event records remain until normal expiry; they are
not deleted to manufacture a clean result.

## Retired unbuilt scraper

`platform/monitoring/scraper-deployment.yaml` was permanently held at0 replicas
because its custom image had never been built. The registry lookup confirmed no
`ghcr.io/onemindos/nats-metrics-scrape:latest` artifact exists. Remove that inactive
placeholder rather than inventing a version or enabling it. It has no pods or
persistent claims. Existing NATS Surveyor, dashboard and Prometheus collectors
remain running. Source Python/Dockerfile remain in `platform/monitoring/` for a
future deliberate build with a real immutable image and secret references.

## Acceptance

Require Argo Synced/Healthy, all running containers Ready, unchanged source PVC/PV
identities, and the intended budgets on both new StatefulSet pods. Check PostgreSQL
readiness, NATS health/JetStream, public/private app response and new policy/pull-secret
warning timestamps. Scan all current top-level controllers and running app pods
against the original image/resource rules; provider-owned kube-system budgets are
outside this application repair. Evaluate the retained history with the exact engine
and verify admission remains unchanged. Record actual runtime results separately.


## Argo delivery correction

Runtime verification found `platform-security` repeatedly reporting successful
server-side sync while both changed ClusterPolicies stayed at generation1 with no
new context/preconditions. The existing `RespectIgnoreDifferences=true` plus
rule-list JQ ignore paths copied the old live rule list into the apply candidate.
Source render and dry-run success alone did not establish delivery.

Kyverno's observed default values are now explicit in all four source policies:
`admission=true`, `emitWarning=false`, `skipBackgroundRequests=true`, and
`allowExistingViolations=true`. Remove the default-field ignore/prepatch settings
while retaining server-side apply. This preserves existing effective policy
semantics and allows future intended rule changes to be reconciled. Verify actual
live contexts/preconditions and auto-generated guards, not only Argo's operation
message, before claiming completion.

## Explicit DigitalOcean budget ownership

The final whole-cluster review identified provider components in `kube-system`
whose resource budgets are DigitalOcean-managed defaults. Their controllers and
pods carry `doks.digitalocean.com/managed: "true"` (observed on Cilium, CSI,
telemetry, Hubble and Konnectivity). Do not take over the network/storage provider
DaemonSets or impose arbitrary caps through a second writer.

Only background resource-budget assessment has a narrow ownership exception:
both the system namespace and this exact provider-managed label must match.
There is no blanket namespace exception. All image checks, all nonempty-UID
admission checks, other security policies, unlabelled system workloads and labelled
workloads in application namespaces remain assessed. Provider defaults are an
explicit exception, not claimed to have been repaired with new application budgets.
The exact-engine fixtures test each boundary and retain original admission results.

The legacy metrics-server is different: it was manually applied in July and lacks
the provider label. Its existing Deployment is adopted into a dedicated
`platform-metrics-server` Argo Application with only a512Mi memory cap added.
The unchanged100m/200Mi requests already exceed its50.6Mi observed seven-day peak.
No Service/APIService/RBAC/image/argument change is included; no existing core
resource-ignore rule is broadened or removed. This adds one Argo ownership record,
not a new metrics workload.

The expanded exact-version suite now covers58 fixtures in four contexts. Provider
boundary comparison preserved all210 admission resource/rule results per CREATE
and UPDATE operation; exactly five background resource findings changed from fail
to skip for the explicit managed-system cases. Image checks remained assessed.
