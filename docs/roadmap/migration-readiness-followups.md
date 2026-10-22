# Migration readiness and outstanding follow-ups

Requested by Zeus on 2026-10-20 so necessary repairs are tracked without making
every platform enhancement a prerequisite for every application migration.
This is a planning/tracking change, not a release, cutover, credential change,
schema migration, waiver of safety gates or claim that Phase 6 is complete.

## Current reconciliation — 2026-10-21 UTC

The [eight-workstream parallel run](backstage-parallel-run.md) is the current
execution index. #340–#345 resolved the scoped Grafana credential/durability/folder
issue; #346–#348 accepted native scoped Prometheus alerts; #349/#350 accepted
native dependency/location topology. Those receipts supersede the corresponding
old pending entries below, not the remaining expiry/role/offsite/launcher/STAC gaps.
New source drafts are not deployed capabilities. Phase 6 remains In progress.

## Resume here: full Kubernetes / Backstage execution map

Zeus subsequently requested PRs for the full Kubernetes structure, which apps we
are using, Backstage soundness and the next implementation steps. The linked PRs
below cover that request without pretending the entire stack is already integrated.

| Order | PR | What it owns / next step |
|---|---|---|
| 1 — portal truth | This tracker **#331**, merged source #330, pipeline #335, promotion #336 | Scoped portal release/rollout and signed-in UI/catalog/pod-boundary checks passed. STAC schema failure and wider trust/recovery/security coverage remain explicit; Phase 6 is not universally complete. |
| 2 — structure | [#333](https://github.com/onemindos/omos/pull/333) | Full 20-category capability map, observed versus proposed apps, controller boundaries, per-app Backstage/onboarding contract and K00–K11 execution packages. Planning only. |
| 2 — foundation assessment (parallel) | [#334](https://github.com/onemindos/omos/pull/334) | Retired community ingress-nginx assessment: eight live nginx-class Ingress resources; installed Cilium/Gateway API parity test, safe canary and per-route rollback plan. No traffic/DNS/identity change authorized. |
| 3 — agent surface | [#332](https://github.com/onemindos/omos/pull/332) | Staged private LibreChat → actual Legacy. Metadata/discovery work; configured DigitalOcean inference returned 402. Identity, sealed credentials, database, backup/restore and session/approval gates must pass before sync. |
| 3 — home preparation | Existing [#284](https://github.com/onemindos/omos/pull/284) | Revise old dedicated-host assumptions for the Mac: storage, ARM64 Linux VM/runtime, 24/7 closed-lid/reboot/FileVault, private access, local recovery and WAN-loss acceptance. Do not move files or create the VM yet. |
| Later — app and site expansion | Existing #283/#320/#321 and scoped follow-up PRs from #333 | One approved app/capability at a time; Karmada/KubeEdge and ML after prerequisites. No blanket namespace migration or competing app controllers. |

**Fresh read-only observation, 2026-10-20 around 16:33 UTC:** four Ready nodes,
Backstage 2/2, 69 Argo Applications with 67 Synced/Healthy; existing Matrix
OutOfSync/Healthy and OpenClaw Synced/Progressing exceptions. No Pending/unready
nonterminal pods in that sweep. Internal HTTP200 checks do not certify application
semantics, browser access, backup restore or AI tool permissions. Full evidence
and actual installed/proposed services are in #333, not inferred from this count.

**Do first:** Section A below, plus the read-only route/ownership/recovery
assessment. Do not begin by installing Crossplane, KubeVela, Karmada, KubeEdge,
Rancher/Fleet or a second mesh. Their inclusion in the blueprint means a bounded
evaluation/acceptance task, not current deployment or a final platform decision.
Review existing Helm-CLI ownership exceptions before adoption. Independent backup
and public ingress lifecycle risks remain visible even while portal work proceeds.

Mac internal free space was only 13–14 GiB; existing external data and Time Machine
volumes are **not** approved for repurposing. The latest GitHub jobs were blocked
from starting by account payment/spending restrictions. No provider, billing,
credential, permission or storage changes are part of these PRs.

All three new PRs are drafts. They may touch the same roadmap/handoff sections:
rebase/reconcile the tracker and rerun checks against current main before merging
them one at a time. Do not resolve conflicts by dropping another PR's gates or
copying a historical acceptance claim over current known failures.

## Rule: core portal gate, then one application's safety gate

The scoped release chain is now merged/verified: source
[#330](https://github.com/onemindos/omos/pull/330), pipeline #335 and image pin #336.
Broader security/recovery gates below remain open. Closing that gate does not fix every service in the
catalog. Once the required core checks pass and outstanding risks are explicitly
recorded in the roadmap, evaluate the selected app on its own dependencies,
permissions, data/recovery requirements and rollback. Do not wait for every link,
optional dashboard or future build facility to be perfect.

Until that live acceptance is recorded, Phase 6 remains In progress and Phase 7
Not started. This document does not itself authorize a migration or relax the
first-unfinished-phase rule. Do not move unaccepted global security failures to
the background merely by calling them follow-ups.

## A. Immediate release/core acceptance — before a new cutover

- [x] Review and merge #330 after its source/release evidence is accepted.
- [x] Open a **separate image-pin implementation PR** against current main; verify
      the chosen image is built from the accepted source and its remote digest.
- [x] Verify Argo Synced/Healthy, two ready replicas, no rollout regression and
      signed-in Project/Component Operations pages, API relationships, source,
      runbook, logs and metrics links. Report failing STAC queries honestly.
- [ ] Check anonymous denial, current machine read-only boundaries, authentication
      refresh/recovery and expiry fail-closed behavior. Distinguish automated
      expiry tests from live session/role coverage. Do not grant a new role simply
      to bypass an untested or broken authorization path.
- [ ] Record known exceptions and the selected migration's acceptance criteria;
      update Phase 6/7 tracker status only from actual accepted live evidence.

Verified deployment update, 2026-10-20, 17:29–17:32 UTC:

- Source #330, pipeline [#335](https://github.com/onemindos/omos/pull/335) and
  promotion [#336](https://github.com/onemindos/omos/pull/336) merged normally.
- Image source `c9a05f8d0b1fc1809b0496945d99593e884f561d`, Linux/AMD64 digest
  `sha256:6c729bd20c0a0a40330be8f0c7bd5213c59c7540bd1ea54c38e05f8c63fc3e7c`.
- Exact Argo revision `1545b50f4fb5521dbdfc9eb22fa302cf2e4b45d6` Synced/Healthy,
  2/2 ready/updated/available replicas. Both pods passed anonymous/machine denial
  checks and expose exactly six read-only MCP tools.
- Signed-in Component/Project Operations, real pilot logs/metrics and packaged
  runbook passed. Catalog: 716 entities, no duplicates/errors/orphans, reviewed
  API/dependency edges, only the known inactive workflow-plane reference missing.
- **Old source `af7770fb` / digest `fab652...` is superseded; do not promote it.**
- STAC collection/search still return 500 and data is unknown. Existing-session
  reload worked; second-role and production-expiry live coverage is not asserted.
- Full evidence and boundaries: [STAC release acceptance](../operations/stac-operational-slice-2026-10-20.md).
  Pipeline first profile is portal; removal remains plan-only, not a universal
  app onboarding/removal executor or a newly scheduled daemon.

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
- [x] Scoped Grafana adapter/durability/folder repairs and acceptance linked: #340–#345;
      see [durability receipt](../operations/grafana-durable-visibility-2026-10-20.md)
      and [folder receipt](../operations/grafana-native-folder-2026-10-20.md).
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

- [#333: full Kubernetes operating structure](https://github.com/onemindos/omos/pull/333)
  and [#334: public ingress lifecycle](https://github.com/onemindos/omos/pull/334):
  preparation and explicit acceptance tasks, not installed extra platforms.
- [#332: private LibreChat/Legacy pilot](https://github.com/onemindos/omos/pull/332):
  staged source, not usable chat or proven persistence until its gates pass.
- [#284: home expansion](https://github.com/onemindos/omos/pull/284): update against
  the user's actual Mac production requirement before infrastructure work.
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
