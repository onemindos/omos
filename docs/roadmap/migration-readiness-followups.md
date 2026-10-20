# Migration readiness and outstanding follow-ups

Requested by Zeus on 2026-10-20 so necessary repairs are tracked without making
every platform enhancement a prerequisite for every application migration.
This is a planning/tracking change, not a release, cutover, credential change,
schema migration, waiver of safety gates or claim that Phase 6 is complete.

## Rule: core portal gate, then one application's safety gate

The immediate release gate is the existing
[PR #330](https://github.com/onemindos/omos/pull/330), its separately reviewed image
pin and live verification. Closing that gate does not fix every service in the
catalog. Once the required core checks pass and outstanding risks are explicitly
recorded in the roadmap, evaluate the selected app on its own dependencies,
permissions, data/recovery requirements and rollback. Do not wait for every link,
optional dashboard or future build facility to be perfect.

Until that live acceptance is recorded, Phase 6 remains In progress and Phase 7
Not started. This document does not itself authorize a migration or relax the
first-unfinished-phase rule. Do not move unaccepted global security failures to
the background merely by calling them follow-ups.

## A. Immediate release/core acceptance — before a new cutover

- [ ] Review and merge #330 after its source/release evidence is accepted.
- [ ] Open a **separate image-pin implementation PR** against current main; verify
      the chosen image is built from the accepted source and its remote digest.
- [ ] Verify Argo Synced/Healthy, two ready replicas, no rollout regression and
      signed-in Project/Component Operations pages, API relationships, source,
      runbook, logs and metrics links. Report failing STAC queries honestly.
- [ ] Check anonymous denial, current machine read-only boundaries, authentication
      refresh/recovery and expiry fail-closed behavior. Distinguish automated
      expiry tests from live session/role coverage. Do not grant a new role simply
      to bypass an untested or broken authorization path.
- [ ] Record known exceptions and the selected migration's acceptance criteria;
      update Phase 6/7 tracker status only from actual accepted live evidence.

Recorded release candidate, **not a deployed-image claim**:

- Source commit: `af7770fbd06d630ddd1a9256d759c1eb0995e921` (#330 head when checked).
- Linux/AMD64 image publication/remote verification completed locally.
- Candidate digest: `sha256:fab6522342b40e6ce595ab2e932c84399d3aefec843af45829bade802f955de3`.
- Recheck head/digest/provenance before pinning, especially if #330 changes.

GitHub jobs have been unable to start because of account billing restrictions;
that is not passing CI. Preserve all required checks and permissions. Where a
reviewed local validation path is permitted, attach its actual results; if a
required merge/release rule cannot be satisfied, that release stays blocked.
Do not change billing, disable checks or silently bypass policy. Existing
whole-tree YAML-formatting debt must remain separate from changed-file results.

## B. Required for each chosen application — not a whole-namespace move

Put these in that app's implementation PR, linked from the Phase 7 parent
[PR #283](https://github.com/onemindos/omos/pull/283):

- [ ] Confirm this is an application Zeus intends to keep; select it deliberately,
      not merely because its namespace contains several services.
- [ ] Record current owner/controller, canonical catalog identities, source and
      image, endpoints, dependencies, state/PVCs and intended target owner.
- [ ] Resolve provider/name collisions before creating managed entities. The
      existing Maputnik candidate [#320](https://github.com/onemindos/omos/pull/320)
      is not automatically ready: review its catalog identity collision first.
- [ ] Prove its actual user/API path, required dependencies, permitted human/AI
      reads, network isolation and useful timestamped observability.
- [ ] For stateful/dependent services, preserve a backup, verify isolated recovery
      and define data consistency, storage ownership, cutover and recovery limits.
      A ready DB pod or volume synchronization alone is not restore proof.
- [ ] Define a reversible cutover and test rollback. Avoid concurrent controllers
      or two writers to the same state. Image rollback need not reverse DB changes.
- [ ] Pin artifacts, render/admit configuration, run validation, merge through the
      approved path, let Argo reconcile, then verify app behavior and catalog state.

Prefer one important, reversible/stateless keeper after the core gate; do not
blindly bulk-migrate a namespace. Shared/foundation services can remain Argo-owned
and visible in the portal. Representation is not the same as deployment ownership.

## C. Tracked repairs and parallel work

Each row below needs its own bounded implementation PR when worked. The checkboxes
are outstanding work, **not claims that separate implementation PRs already exist**.
This tracker is the single place to retain the work until those PRs are opened.

| Follow-up | When it blocks a cutover | Acceptance / safe scope |
|---|---|---|
| STAC schema/query repair | STAC or an app that requires its failing collection/search path; not every unrelated app | Inventory existing data and exact image/schema versions; encrypted backup and isolated restore; review an Argo-owned initialization/migration; verify real collections/search and browser behavior. No blind schema write or demo-data fabrication. |
| Grafana dashboard-adapter authentication (observed 401) | If the selected app requires that dashboard and has no accepted monitoring alternative | Inspect actual adapter configuration; handle any needed credential work in its explicit scope without printing secrets; prove authenticated read-only queries and denial boundaries. Existing OpenChoreo metrics are a separate path. |
| Sixteen flagged launcher findings | A selected app's required access path must work; unrelated links can be repaired in parallel | Verify DNS/TLS/routes and intended hosting. Repair or explicitly classify each, preserve stable IDs/favorites, and never infer a dead workload from an unreachable launcher. Do not bypass TAK certificate checks. |
| Independent offsite backup copy | Resolve before higher-risk/stateful moves unless the app's reviewed recovery policy explicitly accepts the documented residual risk; not a blanket prerequisite for a disposable stateless pilot | Select independent destination, protect/export recovery keys, configure approved credentials, prove copy and restore. Do not blindly unsuspend an unconfigured CronJob. |
| Historical-backup warning presentation | Usually parallel; actual loss of required backup coverage still blocks the affected app | Distinguish superseded Backstage and retired ClickHouse jobs from current mechanisms. Keep the real offsite gap and unknown recovery evidence visible. |
| Remaining live token-expiry/role coverage | Any operation relying on an unproven security boundary; not permission to widen machine actions | Complete scoped negative/refresh/expiry checks, preserve fail-closed policy, and record what was tested versus still unknown. Source regressions are not universal live-role certification. |
| GitHub execution restriction and pre-existing validation debt | If required checks/release policy cannot be satisfied | Resolve account restriction with the owner without changing billing unasked; fix formatting in separate scoped PRs. Never describe local passes as GitHub green. |

- [ ] STAC repair implementation and acceptance evidence linked.
- [ ] Grafana adapter repair implementation and acceptance evidence linked.
- [ ] Launcher findings each repaired/classified and evidence linked.
- [ ] Offsite copy decision, implementation and restore evidence linked.
- [ ] Historical-backup UI clarification implemented and verified.
- [ ] Security session/role evidence and remaining limits reconciled.
- [ ] CI/validation follow-ups resolved or explicitly dispositioned per release policy.

For the observations behind these rows, see the
[October 20 trust/coverage review](../operations/portal-acceptance-followthrough-2026-10-20.md)
and the implementation/evidence in [#330](https://github.com/onemindos/omos/pull/330).
Re-observe live conditions before acting; this is not a fresh health audit.

## D. Already tracked separately — do not duplicate or silently deploy

- [#328: workflow-plane pilot](https://github.com/onemindos/omos/pull/328): high
  priority, optional for prebuilt-image migration; no build-plane installation here.
- [#283: Phase 7 app migrations](https://github.com/onemindos/omos/pull/283): parent
  migration work, still subject to core and per-app gates.
- [#320: Maputnik candidate](https://github.com/onemindos/omos/pull/320): a draft
  candidate, not the automatically selected first keeper or an accepted cutover.
- [#321: post-migration cleanup](https://github.com/onemindos/omos/pull/321): only
  remove superseded resources after proven cutover, dependency review and rollback
  retention. Do not delete storage or stable catalog identities as cosmetic cleanup.

**Exit from planning:** accepted core portal evidence plus an explicitly selected,
safe per-app migration PR. Recording backlog does not mean the repairs are done.
