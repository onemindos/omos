---
status: selected-target-implementation-pending
user-confirmation: reiterated in the supplied application-flow list, 2026-10-21
---

# Application delivery and infrastructure ownership

## Decision and authority

Record the user's selected target: **Backstage + Git/Argo CD + KubeVela +
Crossplane**, with Terraform supplying the underlying clusters/network and NATS
transporting approved operational events. This corrects the older structure
plan's generic “comparison candidates” wording for KubeVela/Crossplane.
Selection is **not installation, version approval, coexistence acceptance,
provider access, a production cutover or authorization to remove OpenChoreo**.
The user requested the foundation for work, not those changes.

The historical running model remains native Argo/operator resources and the
existing OpenChoreo pilot/integrations. The October 20 baseline found no matching
KubeVela/Crossplane CRD groups. That is historical observation, not a fresh cluster
inventory. Keep current workloads and native Backstage integrations until a
separately accepted transition proves their replacement. The four-layer trust
rules and human approval for foundation/shared-service changes stay in force.

## Target flow

```text
Backstage request / authenticated agent (same reviewed application contract)
                    |
           validated Git proposal + required human approval
                    |
                 Argo CD
              /             \
     KubeVela intent      Crossplane intent
     app lifecycle        approved app infrastructure
              \             /
            Kubernetes + owning domain operators
                    |
          actual app/data/deployment evidence
                    |
        Backstage Operations + scoped agent tools

Terraform: clusters/network and explicitly retained infrastructure ownership
NATS: versioned events/projections, not another desired-state writer
```

The application contract is a logical specification mapped to each tool, **not a
new installed CRD/controller or separate database of deployment truth**. Begin
with the lifecycle proposal in #351 and map stable catalog refs to Git intent;
do not create another scaffolder/catalog or require all applications to use one
abstraction before they can appear in Backstage.

## One owner per resource and field

| Concern | Target boundary | Prohibited overlap / acceptance needed |
|---|---|---|
| Clusters, networks and bootstrap | Terraform retains explicitly recorded existing resources | Crossplane must not import/reconcile the same cloud object concurrently. State/backend recovery, import/hand-off and deletion protection need review before any transfer. |
| Git intent | Argo owns reviewed KubeVela Application and Crossplane composite/claim intent, plus approved native foundation manifests | Argo must not also apply their generated Deployment/managed resources. Audit tracking IDs, health, finalizers, prune and failure recovery. |
| App runtime | KubeVela definitions/traits/workflows own the approved generated app resources | Disable overlapping Flux/Fleet/Helm/other-controller ownership. Keep native/operator escape hatches and stateful/data ownership explicit. |
| App infrastructure | Crossplane provides narrowly scoped compositions/providers and reviewed lifecycle policy | It can compose Kubernetes resources too; our boundary deliberately excludes the KubeVela-owned runtime fields. No automatic cloud credentials or app-DB provisioning by selection alone. |
| Application data | The app/domain operator owns database consistency, migrations, credentials and restore | A Crossplane resource becoming Ready does not prove usable schema, a backup or app-consistent recovery. |
| Current OpenChoreo apps | Existing ownership remains until individually accepted transfer | Never have OpenChoreo and KubeVela reconcile the same workload. Inventory CRs/generated resources, quiesce ownership, test rollback/data consistency before cutover. Existing portal views stay. |
| Agent actions | Backend-verified principal/policy/approval/evidence; #339 remains gated | A portal button, CLI, plugin MCP tool or NATS subject is not execution authority. The offline evaluator grants no live permission. |
| Operational events | Existing NATS doctrine and app-specific authenticated adapters | Bus messages do not bypass Git review, grant device control or automatically reconcile infrastructure. Physical safety/local authority remain outside LLM decisions. |

## One application contract: mapping requirements

Version the contract before implementing adapters. Record:

- Stable app/catalog/API identities, accountable owner and intended domain use;
  reviewed primary/served [OneMind pillars](../../roadmap/onemind-catalog-pillars.md)
  are mission classification, not deployment ownership or authorization.
- Source, immutable artifact, architecture, environment, site/cluster/namespace.
- Deployment owner/controller per object/field and authoritative desired-state ref.
- Runtime and infrastructure dependencies, provisioning/configuration bindings,
  readiness meaning, deadlines and partial-failure/rollback behavior.
- Human access and authenticated AI capabilities; operation schema, actor, approval,
  run/correlation ID, limit/budget, expiry, revocation and audit/evidence location.
- Scoped metrics/logs/traces/alerts, observed time, source freshness and semantic
  application checks, not just pod/CR health.
- State/storage ownership, export, retention, independent keys/backup, verified
  restore and RPO/RTO; delete/orphan/retain rules and schema compatibility.
- Cost allocation window/currency/basis, provider-versus-estimated agent spending
  and unavailable data. Never derive a zero bill from missing samples.

Map these fields to catalog metadata, Git manifests, composition inputs, native
plugin selectors and evidence contracts. No credentials or private app content
in the public descriptor. A concrete versioned schema and migration policy are
future implementation work, not supplied by this prose.

## Unexecuted implementation gates (K08)

1. Select/pin versions, definitions/compositions and provider capabilities against
   the actual Kubernetes/Argo/Backstage versions; review license/maintenance,
   resources, network access and independent bootstrap/recovery. No “latest”.
2. Reproduce one disposable stateless app and one bounded infrastructure fixture;
   prove CRD/Helm/operator escape hatches without moving existing production.
3. Prove actual tracking/health, cross-controller identity and **no double owner**.
   Test dependency readiness, unavailable provider, denied credentials and broken
   bindings; do not publish success from merely submitted intent.
4. Test update/drift/restart and failure compensation, safe retention/finalizers,
   deletion/orphan behavior, independent restore and compatible rollback. A
   dependency graph or image rollback is not enough for stateful transitions.
5. Prove native Backstage mappings/permissions and user/AI reads at every layer;
   do not assume a maintained KubeVela plugin for our cohort exists. Use existing
   Kubernetes/Argo/API views where they meet the need, custom glue only for gaps.
6. Only then prepare an explicit app-specific transition PR. Keep the existing
   OpenChoreo pilot and its recovery path until its own cutover is accepted.

Argo reconciliation and provisioning can span multiple controllers and are not
an atomic transaction. Readiness sequencing, retries, idempotency, timeouts,
compensation and deletion safety must be explicit in the future implementation.
No existing cloud object, resource or workload is adopted by this decision file.

## Primary sources and limits

- [KubeVela Argo integration](https://kubevela.io/blog/2023/01/06/kubevela-argocd-integration/)
  documents an integration pattern, not compatibility of a newly selected release
  with this cluster. [Definitions](https://kubevela.io/docs/getting-started/definition/)
  describe extensible application building blocks; avoid enabling overlapping addons.
- [Crossplane/Argo integration](https://docs.crossplane.io/latest/guides/crossplane-with-argo-cd/)
  documents annotation tracking and health considerations. The consulted latest
  documentation identified v2.4; **we have not selected or installed that version**.
- [Current operating structure](../kubernetes-operating-structure.md) retains all
  twenty categories, foundation work and K00–K11 dependencies. This decision
  supersedes only KubeVela/Crossplane *selection wording*, not their missing gates.
