# OneMind catalog structure: SO / LE / HP / GE

User requirement, 2026-10-21: make the four OneMind pillars part of the **catalog**,
not merely launcher categories. This is a foundation/specification in master
#356, not applied catalog descriptors, new access rules or a deployed UI change.

## Names and source of truth

Use the existing [pillar doctrine](../../platform/hermes/agent/doctrine/PILLARS.md)
and [agent registry](../../platform/hermes/agent/AGENTS.yaml) for terminology:

| Display order requested | Name | Mission context |
|---|---|---|
| SO | Sovereign Operations | Technical infrastructure, security, development, automation, agents, TAK, drones and robotics |
| LE | Legacy Evolution | Family, home/land, homestead, off-grid systems, heritage and community |
| HP | Holistic Performance | Personal health/performance, habits, learning, development and personal finance |
| GE | Generational Empire | Ventures, products, brand, business operations, partnerships and long-term enterprise |

Do not rename them Human Performance, Life Engineering or Global Enterprise.
The existing doctrine explicitly says the historical numbered domain ranges and
routing rules still require Zeus's decision. Preserve that history: SO 00–24,
HP 25–49, LE 50–74, GE 75–99. The requested **display** order SO/LE/HP/GE does not
swap those ranges or authorize enforcing them. No numeric routing/ACL change,
agent activation, persona change, NATS subject change or vault reorganization.
An agent persona's scope is not evidence of actual credentials or live access.

## Native catalog model

Propose these collision-checked canonical refs (all remain unregistered here):

```text
OneMind                         Domain: domain:default/onemind
├── SO — Sovereign Operations   Domain: domain:default/onemind-so
├── LE — Legacy Evolution       Domain: domain:default/onemind-le
├── HP — Holistic Performance   Domain: domain:default/onemind-hp
└── GE — Generational Empire    Domain: domain:default/onemind-ge
      └── reviewed business Systems
            └── existing Components / APIs / Resources
```

The last two lines apply to **each** pillar, not just GE. Use native
[Domain/System/Component semantics](https://backstage.io/docs/features/software-catalog/system-model/):
child Domains use `spec.subdomainOf`, Systems use `spec.domain`, and existing
Components/APIs/Resources retain their reviewed `spec.system`. The standard
[descriptor fields](https://backstage.io/docs/features/software-catalog/descriptor-format/)
do not provide multiple primary Systems or multiple primary Domains per System.
Do not invent those fields or manually author processor-generated `relations`.

Illustrative child descriptor **inside this Markdown plan only**:

```yaml
apiVersion: backstage.io/v1alpha1
kind: Domain
metadata:
  name: onemind-so
  title: SO — Sovereign Operations
  labels:
    onemind.io/catalog-axis: mission
    onemind.io/pillar: so
spec:
  owner: group:default/onemind-admins
  subdomainOf: domain:default/onemind
```

The parent Domain must exist before processing children. Keep the existing human
accountable owner until an explicitly reviewed ownership decision; the pillar's
agent is **not** made the human owner/administrator. Check live provider refs and
source identities before reserving the names. No files under live catalog paths
are changed by this plan. The installed Domain schema accepted a synthetic child
Domain with `subdomainOf`; this is schema evidence, not registration/live graphs.

## Mission axis versus deployment axis

The current source has capability-oriented Systems such as `infrastructure`,
`observability`, `personal-finance`, `zoho-apps`, `communication` and `workspaces`,
plus `fabric-*`, `cluster-*` namespace Systems and OpenChoreo-managed entities.
Those are **not already the four pillars** and do not have source `spec.domain`
mappings to them. Inspect each one rather than bulk assigning all descendants.

Keep these independent:

- **Pillar:** why the app/system exists in OneMind.
- **System:** its reviewed logical capability context and current native membership.
- **Deployment:** site/cluster/environment/namespace, desired state and controller.
- **Owner/access:** accountable owner and actual authenticated authorization.
- **Launcher category:** existing practical app grouping, not automatically a pillar.

Do not overwrite a generated/OpenChoreo System's membership/provider ownership
just to make a four-column UI. Do not move workloads/namespaces into four buckets,
change Argo Applications or rename existing Components. Preserve stable refs,
favorites, links, selectors and current dependency/storage/location relationships.

Where an existing business System has a coherent reviewed pillar, map its
`spec.domain` without changing its name or its Components' `spec.system`. Where
membership represents a technical namespace or mixed mission use, preserve it
and use explicitly reviewed classification/projection until a separate logical
model migration is accepted. The UI must show native membership and classification
as different evidence, not manufacture Domain graph edges from a label.

## Cross-pillar services without duplicate entities

Propose a small, validated classification contract:

- `metadata.labels['onemind.io/pillar']`: one primary code `so|le|hp|ge`, only
  after review; leave unreviewed entries visibly **Unclassified**.
- `metadata.annotations['onemind.io/pillars-served']`: a JSON string array of
  unique reviewed codes, for additional supported mission contexts. Only approved
  consumers interpret it; it does not automatically create native relations.
- A reviewed manifest/projection records canonical ref, primary/served pillars,
  applicability, reviewer, non-secret evidence and review time. Schema validation,
  stable deduplication, inheritance/conflict rules and refresh tests are future work.

Example: NATS may have SO as its primary technical context and serve LE/HP/GE.
It must still be **one** NATS catalog service, not four copied Components or four
bus deployments. A platform-run Nextcloud may serve multiple missions; classify
real instances/capabilities without exposing personal file contents. Neither
example is an approved assignment or evidence of deployment.

Do not infer primary pillar solely from namespace, vendor, current host, launch
category or an agent name. In particular, shared `communication`, `workspaces`,
`external` and mixed provider suites need deliberate review. Not every externally
hosted app needs Kubernetes metadata; not every resource is a business application.

## Catalog and UI acceptance

The foundation's U00/U01/U02/U03 must include:

1. A OneMind catalog entry/landing and first-class SO/LE/HP/GE Domain pages in the
   requested order, with human names, purpose, accountable owner and clear scope.
2. **All / SO / LE / HP / GE / Shared / Unclassified** discovery controls that
   preserve per-user permissions. Primary/serves-pillar semantics and multi-pillar
   deduplication must be visible; Shared is a view, not an invented fifth pillar.
3. Filtered System/app/API/resource listings, native relationship drill-down and
   a clear route back to the global catalog. Counts must distinguish one app,
   multiple deployment instances and shared visibility; no false quadrupling.
4. App header/breadcrumb and onboarding proposal pillar selection, separate from
   hosting/site/environment/controller and existing launcher category.
5. Search facets, reviewed scorecard facts, attention/notification context and
   scoped AI discovery using the same mapping. No independent hard-coded pillar
   lists scattered across plugins or all-or-nothing agent access per pillar.
6. Denial tests proving classification/filtering does not grant data, infrastructure
   or device access. Preserve privacy boundaries within a pillar and across them;
   existing admin/read-only policies are not replaced by SO/LE/HP/GE labels.

## Bounded rollout work — not performed here

- [ ] Confirm exact refs/collisions, current memberships/providers and primary versus
      shared semantics against an authorized processed catalog; preserve an export.
- [ ] Review/validate one classification source with enums, immutable stable refs,
      owner/evidence, missing/ambiguous states and no credentials/content.
- [ ] Add paired Domain descriptors and location/config inclusion through reviewed
      GitOps, with parent/owner/reference/schema tests. Do not apply a branch file.
- [ ] Map a small deliberate set of business Systems and shared apps; keep other
      records unclassified, not falsely certified or moved wholesale.
- [ ] Extend curator/refresh and ingestion contracts so mappings survive inventory
      updates, deletions and OpenChoreo provider reconciliation without collisions.
- [ ] Implement native catalog filters/landing/entity presentation, then verify real
      graphs, search, privacy, no duplicate refs and persisted favorites.
- [ ] Build/release only required runtime changes; let Argo reconcile mounted catalog
      updates naturally and record both-replica catalog/UI acceptance separately.

Work with existing #351/#353/#355 and the experience packages, not a parallel
catalog/provider. This requirement does not start Phase 7 or move apps between
controllers. The initial deliverable is the clear plan; the catalog is not changed
or presented as live by documenting the desired hierarchy.
