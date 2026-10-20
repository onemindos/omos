# OpenChoreo: shared infrastructure, migration and the Kubernetes COP

Research date: 2026-10-20. Checked v1.3.x documentation, the v1.3.1 CLI and
installed CRD schemas. This corrects an overly broad recommendation to exclude
all shared platforms; it does not grant new AI write permissions or record a
completed migration.

## Confirmed capabilities

**Managed infrastructure is part of OpenChoreo's design.** ResourceTypes can
emit Kubernetes objects or operator CRs for infrastructure such as databases,
queues, caches and object stores. The documented Postgres/NATS/Valkey examples
are local-development samples, not production replacements for our existing
operators. Production templates can target native operators and expose connection
outputs while retaining state on deletion.
[ResourceType authoring](https://openchoreo.dev/docs/platform-engineer-guide/resource-types/).

**Stateful applications are not technically excluded.** The installed
ClusterComponentType CRD accepts `deployment`, `statefulset`, `cronjob`, `job`
and `proxy` workload types. Adopting stateful data still requires an explicit
storage, identity, deletion and recovery handoff.

**GitOps and OpenChoreo are complementary.** GitOps reconciles OpenChoreo
declarations; OpenChoreo renders/applies their resulting workloads. Keep one owner
per resulting object rather than pointing Argo at both the source declarations
and the generated Deployment.
[GitOps overview](https://openchoreo.dev/docs/platform-engineer-guide/gitops/overview/).

## What is not automatic

1. A ClusterResourceType is a reusable cluster-scoped **template**, not a single
   globally shared running database. A Resource belongs to one Project. Do not
   assume cross-project Resource references or transparent adoption of the existing
   NATS/Postgres clusters. Choose an explicit shared-service/consumer contract.
   [Resource API](https://openchoreo.dev/docs/reference/api/application/resource/).
2. Creating an OpenChoreo Resource does not prove every desired Backstage view
   exists. The documented catalog sync matrix explicitly includes Components,
   Projects and platform kinds; do not infer a complete managed-resource UI from
   catalog rules allowing Backstage's separate `Resource` kind.
   [Catalog sync](https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/catalog-sync/).
3. The control-plane MCP documents managed-resource tools. Their actual presence
   depends on version/toolsets and role permissions. Verify tools/list plus a real
   read call; don't expand the existing read-only identity to make a demo pass.
   [MCP reference](https://openchoreo.dev/docs/reference/mcp-servers/).
4. Release resource trees start at a release's rendered roots. Operator-created
   descendants may need extra discovery rules and narrow read RBAC. This is not
   automatic discovery of every unrelated workload in the cluster.
   [Resource tree rules](https://openchoreo.dev/docs/platform-engineer-guide/resource-tree-rules/).

The COP must present both OpenChoreo-managed releases and independently managed
foundation resources with explicit source/freshness/management ownership. The
existing Kubernetes MCP can observe either group; Backstage catalog presence alone
is not a live observation or evidence of functional health.

Live MCP preflight exposed 63 control-plane tools, including `list_resources` and
`get_resource`, under the existing read-only identity. This verifies the managed-
resource toolset exists here; it does not prove populated resources or write access.

## Corrected management policy

- Suitable applications: OpenChoreo Components/releases.
- Managed data/service dependencies: OpenChoreo Resources/ResourceTypes backed by
  appropriate provisioners. Evaluate ownership of the service instance separately
  from ownership of its operator/bootstrap.
- Shared platform instances are **eligible for reviewed integration/migration**,
  not permanently excluded merely because several apps use them.
- Argo, identity, essential networking, secrets/recovery, OpenChoreo itself and
  current protected TAK core retain an independent recovery path. This is our
  reliability policy, not a claim that OpenChoreo cannot render their objects.
- Existing machine identities remain read-only. Management visibility, human
  authorization and autonomous AI mutation are separate decisions.

## Existing-application migration

The official Workload model accepts a pre-built container image, so migration does
not require rebuilding or importing application source.
[Workload overview](https://openchoreo.dev/docs/developer-guide/workload/overview/).

The documented CLI can generate Workload, ComponentRelease and ReleaseBinding
files in filesystem/GitOps mode. Build-and-release workflows automate those files
and PR creation for future releases; they are not an in-place Helm/database import.
[Release generation](https://openchoreo.dev/docs/platform-engineer-guide/gitops/automations/build-and-release-workflows/).

I found no documented one-command Kubernetes/Helm adoption or application-migration
command in the v1.3.x CLI reference or the verified v1.3.1 binary's top-level help.
This is a bounded finding, not proof no other ecosystem tool exists. Upgrade guides
migrate OpenChoreo versions, not arbitrary existing workloads.
[CLI reference](https://openchoreo.dev/docs/reference/cli-reference/).

Procedure: inventory workload contract → pin the existing image → create Component
and Workload → generate immutable release/binding → deploy a parallel candidate
through Argo → verify functionality/UI/MCP → switch routes/consumers → retire the
old owner through Git. Stateful migrations add a backup/restore and schema plan.

## First quick win

Maputnik has no configured volumes, secrets or backend database. Use its existing
AMD64 image digest and existing `onetak` Project. Add startup/readiness checks and
fix its relative-asset URL canonicalization; preserve the old service during
verification. This limited quick-win exception precedes the original Phase 7
project sequence at the user's request.

Next candidates are chosen from manifests, not names: other small web services,
then feed workers with explicit credentials/contracts. STAC API is not a blind
quick win because its database/schema behavior needs separate verification.
Shared NATS/CNPG integration follows a reviewed ownership/dependency contract,
without swapping the existing servers for demo StatefulSets.
