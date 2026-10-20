# Post-migration repository cleanup plan

**Status: draft plan only. Recorded 2026-10-20. No cleanup or migration authorized.**

Requested by Zeus so repository organization and removal of superseded deployment
definitions are not forgotten after migration. This plan is a Phase 7 follow-up,
not an instruction to begin Phase 8 or change the cluster now. Approving this
document does not approve an application's retirement or a storage deletion.

## Dependencies and start gate

Execute cleanup only after the relevant migration PR is merged and its live
acceptance, rollback and observation-window evidence is recorded. Perform small
project-scoped cleanup PRs after each accepted handoff; perform the final cross-
project audit after the agreed migration scope is complete. Do not wait until
every possible future edge deployment exists to clean up an accepted cloud app.

- [ ] Agree which apps are kept, independently managed, on-demand or retired.
- [ ] Approve the target layout and management owner for each retained workload.
- [ ] Relevant migration and route/consumer cutover are complete, not just drafted.
- [ ] Application behavior, Backstage/COP visibility and relevant read-only
      MCP/CLI observations pass with source timestamps.
- [ ] Existing data, secrets, identity, backups and rollback path remain covered.
- [ ] Agree an observation window and confirm it has passed without regressions.
- [ ] Record the affected Argo Applications, pruning behavior and proposed diff.
- [ ] Obtain explicit approval for each destructive operation or app retirement.

The existing Maputnik draft PR #320 is preparation, not proof a migration happened
or authorization to delete its old resources. No application in this plan is
preselected for retirement.

## Target structure

Keep deployment configuration centralized without combining all application source
into one folder. Use the existing repository responsibilities:

| Area | Responsibility |
|---|---|
| `services/` or a dedicated source repository | Code we own; do not copy every upstream application's source into `omos` |
| `platform/openchoreo/platform/` | Reusable approved ComponentTypes, ResourceTypes, traits and platform policies |
| `platform/openchoreo/projects/<project>/` | Project declarations and approved managed resource definitions |
| `platform/openchoreo/projects/<project>/<app>/` | Component/Workload, generated immutable releases and environment bindings |
| Existing foundation/operator directories under `platform/` | Independently managed services and bootstrap; do not move them just for symmetry |
| `platform/argocd/` and existing ecosystem Application directories | Authoritative reconciliation entry points and pinned chart configuration |
| `infra/` | Infrastructure provisioning configuration |
| `scripts/` | Maintained automation and verification helpers |
| `docs/` | Architecture, app decisions, operational evidence, recovery and this plan |

Environment and site placement belong in approved bindings/configuration; folder
names alone do not establish isolation. Database contents, uploads, model weights,
private keys and backup archives do not belong in Git. A proposed `shared-services`
project is not mandatory: settle project-bound Resource ownership and sharing
contracts before imposing that directory.

## App-by-app decision register

Create a reviewed inventory before execution. One stable logical app identity may
have several site/environment deployments; do not turn each replica into a new app.

For every application record:

- Purpose and actual use; decision: **keep/OpenChoreo**, **keep/independent**,
  **on-demand**, or **retire proposed**.
- Current and target source/deployment paths and authoritative controller.
- Shared-service consumers, human URLs, machine APIs and event contracts.
- Runtime namespaces/sites/environments and data/secret/backup dependencies.
- Migration/cutover PR, observed acceptance, rollback and observation-window links.
- Exact files/resources proposed for cleanup, expected prune effects and approval.

An unused-looking manifest, zero replicas or an old timestamp is not sufficient
evidence that the app or its data is disposable.

## Cleanup work packages

### 1. Remove superseded deployment ownership

- [ ] Compare original manifests with the accepted OpenChoreo deployment contract.
- [ ] Prove traffic, scheduled work and consumers no longer use the old workload.
- [ ] Confirm the old and new reconcilers do not own the same generated object.
- [ ] Remove only the agreed old workload definitions and obsolete Application
      references through reviewed GitOps changes.
- [ ] Review Argo finalizers, prune settings, sync order, shared namespaces and
      shared ConfigMaps/Secrets before removing a parent Application or folder.
- [ ] Preserve controllers, operators and shared service instances still in use.

### 2. Consolidate conventions and automation

- [ ] Use consistent project/app naming, owner, labels and environment bindings.
- [ ] Document how immutable release snapshots are generated; preserve snapshots
      needed by active bindings, the rollback window or retained environments.
- [ ] Review duplicate Helm values, obsolete scripts and build/release workflows;
      remove only after checking CI, Argo, runbooks and external consumers.
- [ ] Update every affected path reference, image-build context and sync source.
- [ ] Keep working automation and intentional backups; do not replace them with
      ad hoc commands just to make the directory tree look shorter.

### 3. Align the catalog and Kubernetes COP

- [ ] Preserve one canonical logical app identity and its favorites/bookmarks.
- [ ] Reconcile curated catalog entries with OpenChoreo-generated entities;
      prevent entity ownership/name collisions and duplicate launcher cards.
- [ ] Update both deployed/source catalog copies where this is still required;
      consolidation of that duplication is a separate tested change, not assumed.
- [ ] Remove stale relationships only after confirming replacement entity refs,
      namespaces, labels, URLs and actual provider ownership.
- [ ] Keep independently managed foundation services visible with live evidence
      and freshness; catalog presence alone does not prove health.
- [ ] Verify human access and relevant AI read paths after catalog changes without
      expanding machine permissions or exposing credentials.

### 4. Retire only approved apps and data paths

- [ ] Obtain the user's app-retirement decision before uninstalling anything.
- [ ] Export required app data and verify the independent restore/recovery path.
- [ ] Inspect dependency relationships, DNS/ingress, certificates, scheduled jobs,
      dashboards, alerts, external clients and event consumers.
- [ ] Treat deletion of PVCs, namespaces, databases, NATS streams, model state,
      secrets/keys, backups or local worktrees as separately scoped operations
      requiring fresh recovery evidence and explicit approval.
- [ ] Do not fold historical database-consolidation proposals into file cleanup
      without their own reviewed migration, schema and rollback checks.

### 5. Refresh documentation and final inventory

- [ ] Update `docs/REPO-LAYOUT.md`, `docs/STACK.md`, relevant READMEs/runbooks and
      root instructions where verified facts or paths changed.
- [ ] Reconcile earlier consolidation proposals in `docs/REPO-LAYOUT.md`; mark
      completed items with evidence rather than executing stale suggestions blindly.
- [ ] Repair local links and remove obsolete instructions, not useful history.
- [ ] Record actual deployment/backup owners, retained exceptions and deferred work.
- [ ] Update the Phase 7 tracker only for completed, evidenced cleanup work.

## Validation and completion gate

Every implementation PR must list exact affected paths/resources and expected live
impact, link its migration prerequisite, and run `make validate` before opening.
If existing repository-wide debt blocks it, report the precise unrelated failures
and run the applicable changed-file checks; never weaken gates or claim green CI.
Also run secret scanning, relevant schema/render tests and reference/link checks.
Pure documentation planning does not require a deployment or Terraform apply.

Before completing each cleanup work package:

- [ ] Relevant Argo/OpenChoreo desired and observed state agrees after pruning.
- [ ] Functional app checks, URLs, dependencies and UI/MCP/CLI evidence still pass.
- [ ] No unintended storage loss, permission change, exposure change or new
      scheduling/capacity regression occurred.
- [ ] There is one documented deployment owner per resource and no superseded
      active definition for the migrated workload.
- [ ] Backups and retained rollback artifacts still protect the relevant state.
- [ ] Final file/path inventory and any intentionally retained exceptions are recorded.

**Not authorized by this draft:** moving files, deleting folders, pruning live
resources, merging implementation PRs, migration, credential changes, billing
changes or changing the existing read-only AI permissions.
