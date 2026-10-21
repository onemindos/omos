# Native service dependency and deployment-location topology

Status: scoped production acceptance passed, 2026-10-21 UTC.
Scope: next bounded Phase 6 visibility slice after Prometheus #346–#348.

## What this graph means

Use the installed Backstage Catalog Graph/Relations plugins and ordinary
`spec.dependsOn`/`spec.system`, not another dashboard or a custom Kubernetes
controller. Graphs describe reviewed source wiring, **not observed traffic**, live
health, complete blast radius, recovery success or authorization to operate apps.
Feature dependencies can degrade one view without stopping the whole portal.
[Backstage relation semantics](https://backstage.io/docs/features/software-catalog/well-known-relations/)
define dependency and membership separately; [Resource descriptors](https://backstage.io/docs/features/software-catalog/descriptor-format/)
allow storage/infrastructure dependencies without creating Kubernetes resources.

## Reviewed relationships

| Source | Target | Meaning and source evidence |
|---|---|---|
| Backstage | Backstage PostgreSQL | Runtime database; `services/onemind-portal/app-config.production.yaml`, `platform/backstage/postgres.yaml` |
| Backstage | Authentik server | Sign-in routing; `platform/backstage/network.yaml` / `dns.yaml`, existing identity configuration |
| Backstage | NATS | Bus summary GETs, not arbitrary message access; `services/onemind-portal/plugins/plugin-onemind-common-backend/src/natsSummary.ts` fixed `nats-headless.nats.svc:8222` target |
| Backstage | Grafana | Native dashboard visibility; `platform/backstage/app-config.cluster.yaml` authenticated fixed Grafana proxy |
| Backstage | Argo CD server | Delivery visibility; same runtime configuration's named Argo locator |
| Backstage PostgreSQL | `backstage/data-backstage-postgres-0` | StatefulSet `backstage-postgres` / claim template `data`; `platform/backstage/postgres.yaml` |
| Grafana | `monitoring/grafana-data` | Existing-claim persistence; `platform/values/kube-prometheus-stack.yaml`, `platform/monitoring/grafana-state.yaml` |
| NATS | `nats/nats-js-nats-{0,1,2}` | Three existing broker claim identities; `platform/nats/nats-statefulset.yaml` template `nats-js`, replicas 3 |
| Workloads/PVC identities | OneMind Cloud — NYC3 | Existing deployment location, **not** prospective Mac/home/edge sites |

The five PVC Resource entries are catalog records of existing claims, not new
PVCs, storage requests or additional charges. Read-only operator observations
confirmed those claims Bound and the source templates/mounts match; no portal PVC
API permission was added, so the graph is not live storage/backup health. Do not
include VolSync cache PVCs as primary application data volumes.

Snapshot-generated Systems previously titled `Cluster / namespace` are relabelled
`Namespace / namespace`. Preserve stable System refs and workload memberships;
there are not dozens of independent clusters. No home/Mac/edge deployment is
invented. Only six reviewed existing Component identities gain topology links;
all unrelated Component relationships and inventory timestamps are preserved.

## Navigation and refresh

Each reviewed Component links to `Dependency topology (documented)`,
`Dependent services (documented)` and its deployment-location Resource. Backstage
and the cluster Resource also link to a three-root core-service view. Native query
parameters are checked against the installed plugin: forward `dependsOn`, max
 depth 2; reverse `dependencyOf`, max depth 1; explicit unidirectional traversal.
This avoids expanding every workload sharing the cluster into an unreadable graph.
Reverse dependency view is **documented dependents**, not complete failure impact.
Existing Kubernetes, Argo, Grafana, Prometheus, source and recovery links remain.

`curate-topology.cjs` is a pure, idempotent part of the existing reviewed curator.
It rejects missing/duplicate targets, changed namespace/selector/location and
unreviewed storage identity collisions; preserves unrelated metadata/custom links.
Inventory refresh restores these mappings. Long quoted graph URLs round-trip
exactly without new YAML line-length failures. Offline source contracts check
PVC identities against the tracked templates; future storage/replica changes need
an explicit topology review instead of silently publishing stale volume identities.

The native graph code is already in the accepted Prometheus image. Only mounted
catalog data changes at runtime: **no image rebuild is needed** for this feature.
Curator/snapshot scripts and future Docker test gates change in source, not running
backend/frontend code. The catalog ConfigMap checksum naturally rolls Backstage
on the same pinned image; do not manually restart or apply a branch manifest.

## Acceptance gates

- Required exact-head portable gates, new Node/Python topology contracts and the
  existing catalog/security/layout cohorts; normal reviewed source PR merge.
- Natural Argo Synced/Healthy, 2/2 same actual portal imageIDs, DB specification/
  controller revision unchanged; Grafana/PVC/NATS unchanged.
- Processed catalog has the five PVC resources, preserved existing identities,
  expected forward/inverse relationships and namespace titles, no new errors.
- Signed-in native core graph shows source-backed cross-namespace/storage/location
  nodes; DB reverse graph shows the portal, and links navigate to real entities.
- Existing machine catalog reader sees reviewed relationships; anonymous catalog
  and machine proxy/write/Secret/exec boundaries stay denied. Catalog access is
  metadata, not authenticated app/device control or a grant of AI write access.

Phase 6 stays In progress, Phase 7 Not started; action draft #339 remains unmerged.
OpenCost and broader topology/traffic/multicluster coverage are separate work.

## Production acceptance

Source [#349](https://github.com/onemindos/omos/pull/349) merged normally at
`82c05d3ce54160696e4a9bb49b5465892a8be7e8`. Required portable exact-head gates
passed for implementation `524e9c75c5d8ee7335961b12fdbb1597bd3cd400`;
18 Node catalog/topology tests and 33 Python pipeline/GitOps contracts passed.
Hosted GitHub CI remains billing-blocked, not green. Historical whole-tree
`make validate` YAML failures begin at `platform/onetak/media-mtx.yaml:86`;
changed-file and required portable checks passed independently.

Natural Argo reconciliation was Synced/Healthy at the exact merge revision;
Backstage generation 77 observed, 2/2 ready/updated/available, zero restarts:

- `backstage-6ffb598ccc-9fklf`
- `backstage-6ffb598ccc-f4qbj`

Both actual imageIDs remain
`ghcr.io/onemindos/onemind-portal@sha256:e103bf7c0103b615bef2f4b9dfcb603ef952af7d9eb95d498117d1f48a5dc0e6`.
The full pod specification matches the pre-change specification after normalizing
only the reviewed generated catalog ConfigMap content hash. PostgreSQL's complete
specification/controller revision are unchanged. Grafana's Deployment UID/spec
are unchanged, its retained 5 GiB claim is Bound and backup schedule enabled;
NATS remains 3/3. No forced sync, imperative workload change or new image was used.

The paired source inventory grew from 509 to 514 entities with no removals.
Both replica machine catalog readers returned the same reviewed entities, all
15 selected references returned 200 with no catalog error items; forward/inverse
relationships, five storage records and namespace titles matched expectations.
This checks the selected slice, not every catalog entry or live traffic.

Signed-in browser acceptance confirmed:

- Backstage's Overview retains native Grafana folder/dashboard links, scoped
  Prometheus alerts and existing operations links alongside new topology links.
- The actual core-service link opens the installed native graph: exactly six
  reviewed Components, five existing PVC Resources and the NYC3 location, with
  bounded forward traversal. No unrelated workload fanout in this core view.
- The database's documented-dependent link shows the database and Backstage.
- The deployment-location link opens the real NYC3 Resource and its explicit
  documented-metadata disclaimer. Namespace membership reads `Namespace / ...`,
  not a claim of additional clusters.

Native generic location-entity Relations still lists its many existing workload
dependents; use its bounded core-service link for the readable scoped view.
This release does not redesign the generic graph or certify all dependencies.

Read-only boundaries were rechecked on both replicas: machine catalog GET 200;
anonymous scoped Prometheus GET 401; machine Prometheus GET/POST and Grafana
proxy GET 403. Kubernetes service-account reads remain allowed for Deployments
and PVCs; Secret reads, pod exec, Deployment patch and pod deletion remain denied.
No permissions or credentials changed. Fixture/source guards complement these
observations; this is not authorization for live AI mutations.

Local receipts: `/tmp/onemind-topology-check/run.json` and
`/tmp/onemind-topology-ops/{application-after,deployment-after,pods-after,catalog-accepted,safety-accepted,rbac-accepted}.json`.
Browser proof: `/Users/zeus/.codex/visualizations/2026/10/20/native-service-topology/core-topology.png`.
Temporary receipts/screenshots are local evidence, not durable backup storage.
The receipt PR changes documentation only; no second image build is needed.

Next bounded visibility slice: OpenCost. Broader alert coverage, native state-text
accessibility, observed network traffic and real multicluster topology remain
separate work. Phase 6 is still In progress; Phase 7 Not started; #339 stays draft.
