# Backstage application lifecycle workstream

Status: draft implementation; source/template improvement, not live acceptance.
Scope: K03 in the Kubernetes operating structure; Phase 6 preparation only.

## Implemented in this PR

Extend the existing `onemind-app-onboarding` proposal skeleton, not a new
scaffolder or deployment controller. Every new proposal carries a common review
contract for identity, ownership, artifact, access, runtime, signals,
relationships, AI access, data, recovery, lifecycle, cost, release, acceptance and
retirement. Existing parameters and credential handling are unchanged.

Correct historical instructions that treated the compatibility directory JSON as
live launcher discovery or required an image rebuild for every mounted catalog
change. Existing proposal outputs remain outside live catalog/GitOps paths.

## Still required before universal app lifecycle is claimed

- [ ] Select one deliberate reversible keeper and review its actual controller,
      data dependencies, state and architecture; no bulk namespace migration.
- [ ] Implement a separate allowlisted app profile in the existing portable
      pipeline, with tests for wrong target/source/digest and semantic failure.
      The current protected portal profile must not be repurposed.
- [ ] Connect reviewed catalog mappings to the coverage report and independently
      scoped acceptance receipts. Report metadata versus live acceptance distinctly.
- [ ] Test an actual rendered proposal through the installed scaffolder, including
      unchanged secret handling; merging this source alone is not browser proof.
- [ ] Prepare source/image/promotion PRs if the proposal skeleton is image-bundled,
      then verify both replicas and signed-in creation before claiming deployed.
- [ ] Prove one actual release/upgrade/recovery round trip, retain its evidence,
      then extend to another app. No imperative workload writes by this script.
- [ ] Keep removal app-specific, explicit and retention-aware. Do not implement
      generic deletion or enable AI mutation under the name of lifecycle work.

## Validation and boundaries

The offline contract test checks the reviewed skeleton and the existing template
still publishes only draft proposals, without new live deployment steps. It does
not prove human judgment, render/runtime security or an app's actual recovery.
No deployment, permission, credential, namespace, data migration or storage change.

Related: [release pipeline](gitops-release-pipeline.md),
[operating structure](../architecture/kubernetes-operating-structure.md),
[readiness tracker](../roadmap/migration-readiness-followups.md).
