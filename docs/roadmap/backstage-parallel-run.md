# Backstage completion run — eight workstreams

Requested by Zeus on 2026-10-21 UTC: separate PRs, bounded agent delegation and
coordinator verification. **Draft/source preparation is not production completion.**
Existing OneMind is preserved. Phase 6 remains In progress; Phase 7 Not started.

## Run ownership and status

| Workstream | PR | Assigned work completed in this batch | Still required before completion |
|---|---|---|---|
| Cost visibility | [#352](https://github.com/onemindos/omos/pull/352) | Volta: offline namespace allocation parser, window/currency/partial/freshness contracts and 51 synthetic tests | Real approved response/pricing evidence; existing native page/Observer adapter integration, unknown-state UI, source/image/GitOps/browser acceptance. Not a second Cost Insights page. |
| Complete app coverage | [#355](https://github.com/onemindos/omos/pull/355) | Wegener: deterministic redacted catalog integration-coverage report and 9 tests; coordinator ran on current authorized catalog | One reviewed app batch at a time: real mappings and user/API/log/metric/alert/AI permission checks, not annotation-only badges. |
| Catalog/navigation | [#353](https://github.com/onemindos/omos/pull/353) | Noether: launcher finding/closure report and 8 tests; coordinator confirmed 16 review tags still present | Repair/classify each actual endpoint, preserve canonical refs and favorites; no certificate bypass, invented URLs or blanket removal. |
| Trust/recovery | [#354](https://github.com/onemindos/omos/pull/354) | Ptolemy: offline hashed receipt validator and 19 tests, explicit fixture/observation and incomplete states | Authorized live expiry/second-role checks, independent recovery evidence and retained artifacts; no fabricated identities or production fault injection. |
| Standard app lifecycle | [#351](https://github.com/onemindos/omos/pull/351) | Coordinator: extend existing onboarding skeleton with per-app acceptance contract, 2 regression tests; stale live-directory instructions corrected | Signed-in rendered proposal/release proof, one explicit allowlisted app pipeline profile and semantic/recovery acceptance. No generic deletion executor. |
| Approved actions | [#339](https://github.com/onemindos/omos/pull/339) (existing) | Maxwell: first A0 offline request/approval evaluator and 16 tests prepared on existing draft | Verify principals, authorization, audit/locking/revocation/stop controls, selected pilot and human-merged release/recovery. No executor or live write permission. |
| Agents/conversation | [#332](https://github.com/onemindos/omos/pull/332) (existing) | Darwin: explicit conversation/identity/persistence/restore/agent-spend acceptance matrix on existing draft | Recorded provider 402 is unresolved, not freshly re-probed here. Approved inference remedy, identity, storage, backups and end-to-end isolated conversations precede activation. |
| Home/multicluster/edge | [#284](https://github.com/onemindos/omos/pull/284) (existing) | Darwin: revise Mac/home requirements and measured-versus-approved capacity gates | Approved ARM64 VM sizing/storage; lid/reboot/FileVault/WAN-loss/local recovery, real Argo/Backstage registration, separate OpenChoreo/Karmada/KubeEdge acceptance. No VM or storage allocated. |

The coordinator reviews diffs, reruns tests, updates the roadmap in each PR and
owns GitHub/release operations. Workers have disjoint worktrees/files and no
merge, deployment, credential or permission authority. No Ollama delegation was
used; workers inherited the coordinator model. Existing PRs are reused, not
replaced by duplicate plans. Agent assignments are bounded runs, not scheduled
unattended workers.

## Execution order and dependencies

1. Review the index and trust evidence gaps. Baseline is accepted #349/#350,
   not an assumption that every app/security/recovery path is already certified.
2. Prepare independent source/tests in parallel. Rebase/reconcile shared tracker
   edits and rerun exact-head checks **one PR at a time**; a parallel preparation
   request does not authorize simultaneous production cutovers.
3. Cost visibility and scoped coverage improvements can follow their own safe
   read-only release gates. The offline parsers/reports are useful preparations,
   not deployed adapters or a substitute for real app behavior.
4. Prove one important keeper's onboarding/operations/recovery contract. A Phase 7
   migration still requires accepted core and per-app gates and user approval;
   optional cost panels and unrelated launcher repairs are not blanket blockers.
5. Keep #339 executor activation behind trust/approval/audit gates. Keep #332
   behind inference/identity/persistence gates; keep #284 infrastructure activation
   behind approved capacity and local recovery. Planning can proceed now.
6. Separate source merge, immutable image/build proof, GitOps pin, Argo rollout,
   user/app/security acceptance and final receipt. No auto-merge or force sync.

## Current catalog observation, not complete live coverage

Coordinator fetched the authorized complete machine-readable processed catalog on
2026-10-21 UTC using the existing read-only identity and fixed catalog GET. The
private payload is not committed. Both new metadata reports ran successfully:

- 721 entities, 325 Components; 16 launcher review tags remain.
- 186 Component records have Kubernetes namespace/selector metadata, 6 have the
  reviewed Argo mapping evidence, 4 Grafana metadata, 3 Prometheus metadata and
  1 API metadata. These are report predicates, **not counts of working apps**.
- 325 Components include external links/provider applications: absent Kubernetes
  metadata does not automatically indicate a defect. Applicability needs review.
- Every reported runtime/authenticated capability remains `not-tested` unless
  separately observed and recorded. Report totals are not a readiness percentage.

Private report files: `/tmp/onemind-backstage-run/{coverage-report,launcher-report}.json`.
A catalog snapshot is not a deployment audit, recovery archive or live event stream.

## Validation / release boundary

Each PR records focused tests and required portable exact-head results where run.
New tests join existing `test_gitops*.py` discovery rather than remaining orphaned
scripts. `make validate` continues to fail historical whole-tree YAML debt starting
at `platform/onetak/media-mtx.yaml:86`. GitHub hosted CI has the recorded billing
restriction, not a passing status; preserve checks/protections and don't change
billing or credentials. No app/controller/storage/network/identity cutover or
new privilege is executed by this batch. Drafts need review and final exact-head
validation against the then-current main before sequential merges.

## Related foundation and later work — not lost

- [#334](https://github.com/onemindos/omos/pull/334): public ingress lifecycle and route parity.
- K02: verify Argo/TAK bootstrap ownership, Matrix drift, OpenClaw source and
  app-consistent recovery before affected changes; no blind NATS re-adoption.
- [#328](https://github.com/onemindos/omos/pull/328): optional workflow-plane pilot.
- [#283](https://github.com/onemindos/omos/pull/283): later per-app migrations.
- [#321](https://github.com/onemindos/omos/pull/321): cleanup only after accepted cutover/retention.
- [Operating structure](../architecture/kubernetes-operating-structure.md): K00–K11,
  all twenty capability categories, candidate-platform comparisons, ML and physical
  operations. This batch does not install Crossplane/KubeVela/Karmada/KubeEdge.

## Coordinator review follow-through

Independent static review identified trailing-whitespace/type-coercion edge cases
in the new cost contract and a UTC-offset overflow in the receipt parser. Both
were fixed with regression tests (51 cost / 19 receipt tests). A synthetic API
entity reference initially triggered the secret scanner; it was replaced with a
short unambiguous test ref, with no allowlist/baseline/check change.

The older #339/#332/#284 branches were merged with current main for preparation;
shared tracker conflicts were reconciled without dropping acceptance gates. #332
retains both its LibreChat safety/render tests and main's GitOps contract suite.
No draft was production-merged. Final exact-head receipts and focused tests are
reported on each PR; before later merging, recheck against then-current main.

## Experience and architecture foundation — user follow-through

The [capability/UI foundation](backstage-experience-foundation.md) compares the
supplied list against installed source and official Backstage features. It adds
U00–U06 future experience packages (navigation/home, consistent entity pages,
search/docs, events/evidence, agents/actions/sites and release/accessibility),
not a claim that another batch of implementation PRs already exists. The eight
workstream PRs above remain the same; this is an expansion of master #356.

The [ownership decision](../architecture/decisions/0002-application-delivery-and-infrastructure-boundaries.md)
records KubeVela/Crossplane as the user's selected target, implementation pending;
it supersedes old generic candidate wording without installing controllers,
transferring Terraform/OpenChoreo ownership or relaxing existing trust gates.
