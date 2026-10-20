# Native portal visibility: first reviewed services

Status: scoped production acceptance passed, 2026-10-20 18:50–18:54 UTC.
Implementation: [PR #337](https://github.com/onemindos/omos/pull/337), merged
`e98f14366c53faa781933db87196d860fa952f77`.
Scope: existing Backstage host, its PostgreSQL dependency and NATS broker only.
This is Phase 6 stabilization, not a migration or phase completion.

## Baseline observed on 2026-10-20

- Signed-in Backstage Deployment Summary resolves `platform-backstage`, namespace
  `backstage`, instance `onemind-cloud`, Synced/Healthy with a real revision link.
- Native NATS Kubernetes view lists its StatefulSet and Services. The corresponding
  catalog entity has no Argo annotation, so its Argo view is not yet enabled.
- Native Backstage Kubernetes view has a cluster but no Deployment grouping.
  Read-only Kubernetes inspection confirms the Deployment lacks metadata labels;
  its pod selector/template already use `app=backstage`. The DB StatefulSet has the
  same controller-label omission for `app=backstage-postgres`.
- Argo resource tracking confirms Backstage belongs to `platform-backstage` and
  the NATS StatefulSet to `platform-nats`; both Applications are Synced/Healthy.
  Their manifest paths match the mappings. Namespace membership alone is not proof.
- Production portal generation 72, 2 ready replicas, image unchanged from #336:
  `sha256:6c729bd20c0a0a40330be8f0c7bd5213c59c7540bd1ea54c38e05f8c63fc3e7c`.

## Change and blast radius

- Add matching labels only to Deployment/StatefulSet **controller metadata**.
  No selector, pod labels, image or PostgreSQL pod template changes.
- Enrich exactly three existing catalog entities with explicit Argo/source
  annotations and native Kubernetes/Argo/runbook links. Preserve identities,
  owners and dependencies. Backstage still depends on its existing DB entity.
  The broker mapping does not silently claim the separate NATS UI Deployments.
- Reapply mappings during the existing inventory enrichment using a pure,
  idempotent curator that rejects identity/namespace/selector drift and duplicates.
  Preserve unrelated annotations, links and relations. No inventory resnapshot.
- No runtime plugin changes, new packages, credentials, permissions, controller
  installation, message-subject changes, DB migration or imperative deployment.
- Catalog ConfigMap hash changes from `backstage-catalog-2t276t2bd5` to
  `backstage-catalog-bbh6fk2g2t`; Kustomize rewrites the portal volume reference,
  causing its normal two-replica RollingUpdate **with the same image**. Render
  comparison confirms all other resources unchanged except those two controller
  metadata labels and the catalog ConfigMap. PostgreSQL does not roll for this.
- Revert the reviewed Git commit through a PR to undo; let Argo reconcile. Do not
  restore storage or force sync. No publish/promote step is appropriate: this is
  mounted-catalog metadata, not a new runtime image. The image pipeline's source
  equivalence verifier is not a metadata-release acceptance check.

## Pre-merge gates

- Catalog audit/curation gate: 11 tests (original 6 plus 5 curation contracts).
- GitOps wiring: 2 offline tests, included by the existing pipeline test discovery.
  Checks paired catalogs, controller/pod selectors, Argo instance/application/path.
  Kept outside the container-only catalog tests so future Linux image builds do
  not depend on files outside their `services/onemind-portal` build context.
- Semantic before/after Kustomize comparison: only expected changes above.
- `make validate` invoked: blocked by existing unrelated whole-tree YAML lint debt
  starting at `platform/onetak/media-mtx.yaml:86`, not claimed green.
- Required changed-file/pipeline/security gates and review recorded with the PR.
  GitHub checks may remain billing-blocked; do not bypass protections or report
  them passed. User authorized verified local/Linux checks instead.

## Production acceptance after normal PR merge

- [x] `platform-backstage` reconciles the actual merged revision Synced/Healthy.
- [x] Two ready/updated/available portal replicas on the unchanged digest; actual
      observed generation catches up. Database revision/template unchanged.
- [x] Signed-in native Backstage Kubernetes view shows Deployment and ready pods.
- [x] Signed-in native NATS Kubernetes view shows StatefulSet and three ready pods.
- [x] Both native Argo summaries resolve their correct Application, sync/health
      and revision links (not namespace association or empty success).
- [x] Overview dependency/source/native/runbook links are usable.
- [x] Existing service account remains read-only: no Secrets, exec or mutations.

Known remaining work: Grafana dashboard search 401; broader catalog/role/token/
recovery/launcher gates remain open. No blanket cluster health or full COP claim.


## Actual production proof

Argo naturally reconciled the merged revision at 18:48:01 UTC; no forced refresh,
manual apply, restart or sync. Portal generation/observed generation 73, two
updated/available/ready replicas and exactly two remaining pods after old-pod
termination completed:

- `backstage-695b9d4cd7-jzwzs`
- `backstage-695b9d4cd7-mhfww`

Both actual container `imageID` values match the unchanged digest above, both
Running/ready, zero restarts. The initial observation helper incorrectly counted
terminating old pods as active; corrected the helper to distinguish them, then
confirmed the final two-pod list. This was not a failed rollout. PostgreSQL's full
spec and `backstage-postgres-6dcc65dc8` controller revision are unchanged; 1/1 ready.
NATS remains 3/3 ready. No storage or workload-controller ownership changed.

Signed-in browser verification passed:

- Backstage **Deployment → Kubernetes** now shows its Deployment and both ready
  pods, restart counts and actual CPU/memory readings; no empty controller panel.
- PostgreSQL native Kubernetes view shows `backstage-postgres-0`, Running, 1/1,
  zero restarts, with live CPU/memory readings. No database rollout required.
- NATS native Kubernetes view shows `nats-0/1/2`, each Running, 2/2 containers,
  zero restarts, CPU/memory readings and the broker/headless Services.
- Both native Argo summaries resolve the correct Applications and report
  Synced/Healthy. Backstage shows the actual merged deployment revision `e98f143`.
  **Native summary revision is deployment history**, not necessarily the latest
  compared Git revision: NATS shows its last actual sync `ffb7915` (Oct 8), while
  its live `status.sync.revision` has advanced to `e98f143`. Confirmed both in the
  Argo object; no unnecessary NATS resync just to make the history look newer.
- Overview Links/source destinations match the reviewed mappings. The Backstage
  database dependency is clickable. Packaged NATS fabric and Availability/home
  migration recovery runbooks render actual documentation, not empty placeholders.
- Read-only Backstage MCP queries return all three enriched entities and the
  preserved Backstage-to-DB relation, with no entity-processing errors reported.
  This verifies AI access to documented mappings, not new deployment authority.
- Anonymous catalog GET still returns 401. Impersonated portal service-account
  access permits Deployment reads and denies Secrets reads, pod exec, Deployment
  patches and pod deletion. Existing permission/security regressions passed.

Required pipeline gates passed on exact source head
`91c37264941371ca7be41c4b2314a3acf150a4da`: 22 pipeline/GitOps tests, 7 image-release
contracts, immutable dependency install, TypeScript, lint, 82 security/STAC tests,
44 layout tests, 11 catalog/curation tests, changed YAML/schema/Helm, shell,
Terraform and baseline-aware commit/tracked-tree secret scans. Whole-tree
`make validate` remains failed on the historical YAML debt described above.
GitHub checks did not start: their annotation explicitly reports payment/spending
restriction. Normal exact-head PR merge used the authorized local-check path;
no admin bypass, billing, protection or check changes.

## Next visibility work, not silently completed

Grafana's existing bearer credential was retested through its fixed internal host:
`GET /api/search?type=dash-db&limit=100` returned **401**. No credential value was
printed, created or rotated. Its web/basic-health path is separate. Review the
proxy's API-path rewrite as well as the existing credential before claiming native
dashboard search works; do not fix this by enabling anonymous Grafana or forwarding
human tokens. Credential repair requires its own explicitly authorized secret
operation. Then add reviewed dashboard selectors for these services.

Topology, Prometheus and OpenCost remain next integrations, not installed by this
change. Existing repository-wide GitHub cards may show billing-blocked workflow
Errors; those are not NATS-specific build results or runtime health. Native metrics
and readiness do not certify NATS delivery/lag, app correctness, offsite recovery,
whole-cluster health, token-expiry/second-role coverage or Phase 6 completion.
