# ${{ values.name }}

## Purpose
${{ values.reason }}

## Accountable reviewer
${{ values.reviewer }}

This is a proposal, not a release, approval, deployed app or accepted capability.
Never add credentials, private backup contents or temporary signed URLs here.
Unfilled items remain unknown; do not replace evidence with a green badge.

## Proposed app
- URL: ${{ values.url }}
- Hosting: ${{ values.hosting }}

## Application acceptance contract

For each item, record **evidence reference, observation time, reviewer and scope**.
Use passed / failed / not-tested / not-applicable; justify not-applicable. A
configuration value is not proof of live behavior. External-provider apps need
useful account/API/export checks, not invented Kubernetes objects.

- [ ] **Identity:** stable canonical catalog ref, accountable owner, intended use,
      duplicate/provider collision review and application we intend to keep.
- [ ] **Ownership:** existing/target controller, GitOps source, Argo Application,
      namespace, selector and actual cluster/site (or explicitly external).
      Exactly one owner per resource/field; no two controllers applying the same
      generated Deployment. Visibility does not require OpenChoreo migration.
- [ ] **Artifact:** immutable image digest/version, supported architecture and
      source/build provenance; external-provider version limitations recorded.
- [ ] **Access:** canonical route, verified TLS, intended private/public boundary,
      SSO/session behavior and roles; no certificate bypass or implicit admin.
- [ ] **Runtime:** resources, probes, network policy and dependency readiness;
      a ready pod is not proof of correct application data.
- [ ] **Signals:** working scoped Kubernetes/Argo, logs, metrics and alerts where
      applicable; timestamps and missing/unknown coverage shown honestly.
- [ ] **Relationships:** real dependencies, APIs, location and storage identities;
      documented relationships are not observed traffic or full blast radius.
- [ ] **AI access:** actual authenticated read/tool API, principal, exact permitted
      scope, denial tests and audit correlation. Catalog visibility alone does
      not authorize application access or writes. No new permissions by default.
- [ ] **Data:** authoritative store, volume ownership, data classification,
      encryption, retention and export path; sync is not backup.
- [ ] **Recovery:** backup owner/schedule, last app-consistent isolated restore,
      recovery keys independent of the failing system, RPO/RTO and offsite limits.
      Do not restore over production or treat a successful copy as a restore.
- [ ] **Lifecycle:** upgrade/schema compatibility, maintenance window, reversible
      cutover and tested rollback. Image rollback may not undo database changes.
- [ ] **Cost:** resource/storage/network and provider spend boundaries; estimates
      are not invoices, and absent cost data is unknown rather than zero.
- [ ] **Release:** reviewed app-specific pipeline profile and semantic checks;
      the existing protected portal profile is not a generic app deployer.
- [ ] **Acceptance:** exact merged revision/artifact, natural Argo reconciliation,
      real user/API behavior, retained state and negative security tests.
- [ ] **Retirement:** dependency and traffic review, explicit removal approval,
      data export/retention, Argo finalizer/prune review and rollback retention.
      No automatic namespace/PVC/Secret/backup deletion.

## Promotion — separate reviewed implementation

- [ ] Resolve required failed/unknown items or explicitly record the accepted
      residual risk with the accountable reviewer. Never waive a security boundary
      merely to finish this checklist.
- [ ] Add/update the stable entity in both
      `services/onemind-portal/catalog/apps.yaml` and
      `platform/backstage/catalog/apps.yaml` only after mapping review.
- [ ] For launcher discovery, use the catalog's reviewed launcher annotation or
      `Open application` link. The compatibility directory JSON is not the live
      card source; do not require an image rebuild solely for mounted catalog data.
- [ ] If runtime portal code changes, build/verify/promote its immutable image.
      Otherwise use the owning GitOps resource's normal reviewed release path.
- [ ] Record post-merge live acceptance separately. Merging this proposal does
      not register the descriptor, deploy the app or authorize a migration.

Use the existing deployment-change and recovery-preparation templates for linked
proposals. Retirement planning remains read-only; this file grants no authority.
