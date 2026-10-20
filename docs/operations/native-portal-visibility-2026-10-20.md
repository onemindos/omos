# Native portal visibility: first reviewed services

Status: implementation and pre-merge checks; production acceptance pending.
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

- [ ] `platform-backstage` reconciles the actual merged revision Synced/Healthy.
- [ ] Two ready/updated/available portal replicas on the unchanged digest; actual
      observed generation catches up. Database revision/template unchanged.
- [ ] Signed-in native Backstage Kubernetes view shows Deployment and ready pods.
- [ ] Signed-in native NATS Kubernetes view shows StatefulSet and three ready pods.
- [ ] Both native Argo summaries resolve their correct Application, sync/health
      and revision links (not namespace association or empty success).
- [ ] Overview dependency/source/native/runbook links are usable.
- [ ] Existing service account remains read-only: no Secrets, exec or mutations.

Known remaining work: Grafana dashboard search 401; broader catalog/role/token/
recovery/launcher gates remain open. No blanket cluster health or full COP claim.
