# F2 — Catalog cleanup and SO / LE / HP / GE structure

Work order for one implementing agent. Branch: `flagship/f2-catalog-structure`.
Program index: [README](README.md). Read the rules there first.

## Outcome Zeus sees

- **Pillars.** The catalog opens on OneMind → **SO, LE, HP, GE**, each a real
  Domain page listing its Systems and apps.
- **Clean Systems.** About 15–20 meaningful Systems instead of about 70, with
  human names everywhere (no `k8s-backstage-backstage-21b99354` in the UI).
- **One app, not three.** Each app (for example Argo CD) is one entity. Its
  Kubernetes workloads, database and volumes hang under it, instead of appearing
  as separate unrelated rows.
- **Default view.** Noise (code repositories, mesh endpoints, reference/planned
  capabilities) is hidden by default and one filter away.

## Evidence (2026-10-21)

- 348 non-Component entities: about 70 Systems mixing `cluster-<namespace>` (generated
  by `scripts/snapshot-cluster.cjs`), `fabric-*` and capability names
  (`infrastructure`, `personal-finance`, `zoho-apps`, ...). One Domain, `default`.
- Resources: `documented-capability`, `reference`, `planned-capability`,
  `mesh-endpoint`, `code-repository`, all with hash suffixes.
- 325 Components. Curated apps (`catalog/apps.yaml`, 154 docs) and generated
  workloads (`catalog/cluster.yaml`, 514 docs) both describe the same software:
  `component:default/argocd` and `k8s-argocd-argocd-server-…` are unrelated rows.
  There are 0 `subcomponentOf` links.
- Only 6 Components carry `argocd/app-name`; there are 69 Argo applications.
- Specification to follow: `docs/roadmap/onemind-catalog-pillars.md` on PR #356
  (refs `domain:default/onemind`, `onemind-so|le|hp|ge`, label
  `onemind.io/pillar`, annotation `onemind.io/pillars-served`).

## Build

1. **Domains.** Add `domain:default/onemind` and four child Domains with
   `spec.subdomainOf`, titles "SO — Sovereign Operations", "LE — Legacy Evolution",
   "HP — Holistic Performance", "GE — Generational Empire". Owner
   `group:default/onemind-admins`.
2. **System model.** Write `docs/catalog/system-map.md`: every existing System,
   then the target System it maps to, then the pillar. Target set, for example:
   Platform core, Identity & access, GitOps & delivery, Observability,
   Data & storage, Messaging (NATS), TAK & field ops, Geospatial, AI & agents,
   Media, Productivity & comms (SaaS), Personal finance, Home & homestead, Ventures.
   Get Zeus's approval of the map in the PR before changing files.
3. **App-centric model.** The curated app Component is the primary entity. Make
   the snapshot generator emit workloads with `spec.subcomponentOf: <app ref>`
   using a reviewed mapping file (`catalog/workload-map.yaml`: namespace/name to
   app ref). Unmapped workloads stay visible under an "Unassigned workloads"
   System, so nothing disappears.
4. **Argo linking.** Add `argocd/app-name` to the app Component for all 69
   Argo applications (mapping in the same file). Add
   `backstage.io/kubernetes-id`/selector where missing.
5. **Names.** Every entity gets a human `metadata.title`; generated names stay
   as `metadata.name` (stable refs, favorites and links keep working).
6. **Noise.** Give Resource types a consistent set (`database`, `volume`,
   `bucket`, `cluster`, `network`, `code-repository`, `capability`,
   `endpoint`). The catalog index defaults to Kind = Component, owned-or-all,
   with saved filters "Apps", "Infrastructure", "Repos", "Everything".
7. **Pillar filters.** Add a pillar column and a filter (All / SO / LE / HP / GE /
   Shared / Unclassified) to the catalog table. Apps without a reviewed pillar
   show "Unclassified"; nothing is guessed.
8. **Guard tests.** Extend `scripts/test-topology.cjs` / `audit-catalog.cjs`:
   no duplicate titles inside a System, every Component has a System, every
   System has a Domain, every Argo app maps to exactly one Component, and no
   ref disappears without an entry in `catalog/redirects.yaml`.

Reuse #355's coverage report and #353's launcher report as before/after
evidence.

## Out of scope

Live ingestion from the Kubernetes API (F11), the app page layout (F5), new
plugins.

## Files

`services/onemind-portal/catalog/**`, `platform/backstage/catalog/**` (mounted
copies), `services/onemind-portal/scripts/{snapshot-cluster,curate-*,audit-catalog,test-*}.cjs`,
`docs/catalog/**`, and the catalog index page composition in
`packages/app/src/catalogCompositionModules.ts` for the pillar column and filter.

## Acceptance

- [ ] Five Domains render with correct parent/child links. Each pillar page
      lists its Systems.
- [ ] System count is at most 25. Each has a title, a description and a Domain.
- [ ] All 69 Argo apps resolve to one Component. A test fails if a new Argo
      app has no mapping.
- [ ] `argocd` shows its server, repo-server, dex and redis workloads as
      subcomponents on its page.
- [ ] No previously valid entity URL returns "Entity not found" (redirect file
      plus test), and starred entities still load.
- [ ] Default catalog view shows apps only; the "Everything" filter still
      reaches all entities.
- [ ] Required gates pass. Catalog-only changes roll out through Argo after
      merge; record both replicas' processed-entity counts.

## Status

Implementing agent: Claude Sonnet 5.5, 2026-10-21. Not merged, not marked ready, nothing applied to the cluster.
Approval gate: [`docs/catalog/system-map.md`](../../catalog/system-map.md) is committed first and lists the decisions for Zeus.
Everything below is built against that map and is `proposed` until Zeus approves it.

### Done

- **Domains.** `onemind` plus `onemind-so|le|hp|ge`, titles as specified, owner `group:default/onemind-admins`, children use `spec.subdomainOf`.
- **Systems.** 74 Systems became 22 (limit 25), each with title, description and a Domain. Systems with a reviewed pillar sit in that pillar's Domain; seven are Unclassified and sit under `onemind`. Surviving Systems keep their `metadata.name`; 61 removed Systems are in `catalog/redirects.yaml`.
- **App-centric model.** 186 generated workloads are `spec.subcomponentOf` their app (36 new curated apps plus the existing ones) via `catalog/workload-map.yaml`; they join the app's System. Unmapped workloads would land in `unassigned-workloads` (none today; proven by test and by a read-only comparison with the live cluster).
- **Argo linking.** All 69 live Argo applications map to exactly one app. The app carries `argocd/app-name` (primary) and `onemind.io/argocd-apps` (all). A test reads the Application definitions in Git (excluding `platform/parked`) and fails when one has no mapping.
- **Kubernetes selectors.** Apps with a single namespace and a common workload selector get `backstage.io/kubernetes-namespace` and `-label-selector` (only filled when missing).
- **Titles.** Every entity has a human title; workloads read `<App> — <Role>` (for example `ArgoCD — Repo Server`); no duplicate titles inside a System; names are untouched.
- **Pillars.** Label `onemind.io/pillar` (so, le, hp, ge) inherited from the reviewed System, one override (Mercury is Unclassified), annotation `onemind.io/pillar-review` (proposed | unclassified). No `pillars-served` values: no cross-pillar service is approved.
- **Default view and filters (app code).** `catalogViewFilters.tsx` adds saved views (Apps, Infrastructure, Repos, Everything; opens on Apps) and a Pillar filter (All, SO, LE, HP, GE, Shared, Unclassified) to the catalog index sidebar. Components carry label `onemind.io/catalog-role: app|workload` for this.
- **Generator.** `snapshot-cluster.cjs` now emits structured workloads and no namespace Systems; `curate-structure.cjs` runs last in `refresh-inventory.cjs`. Live check (read-only `kubectl get deployments,statefulsets -A`): 177 documents, every `metadata.name` identical to the committed file, 0 unassigned, 0 differences.
- **Guards.** `scripts/test-catalog-structure.cjs` and `test-snapshot-cluster.cjs` (run through `test-topology.cjs`, so the existing gate covers them): five Domains, <= 25 Systems with Domain, titles and no duplicates, every Component/Resource/API has a System, workload and Argo mapping, every ref in `ref-baseline.json` still exists or has a redirect, every relation target resolves, pillars never guessed, role labels, byte-stable rewrite, mounted copies equal.
- **Catalog files stay byte-faithful.** The curator patches changed fields only (the first version re-serialized everything and broke the launcher-review fixture; fixed and tested).

### Zeus's decisions applied (2026-10-21)

- Pillars: Intelligence & feeds and Communication are SO, Brand & social media is GE; Productivity & SaaS, Personal tools & devices and Reference & community are Shared (under `domain:default/onemind`, label `onemind.io/pillar: shared`, found by the Pillar filter's Shared view). Unassigned workloads stays Unclassified. Mercury moved to Business operations (GE). Argo convention approved. Domain names stay `onemind-so|le|hp|ge`; F1 fixes its links.
- Old System URLs: `packages/app/src/systemRedirects.tsx` registers one native `PageBlueprint` page per removed System at the static path `/catalog/default/system/<old>`, rendering react-router `<Navigate replace>` to the successor. A static path outranks the catalog's `/catalog/:namespace/:kind/:name`, and no entity exists under the old names, so live entities are never intercepted. Data comes from `systemRedirects.json`, generated from `catalog/redirects.yaml` by `curate-structure.cjs`. No third-party code. Tests: forwarding, no interception of live Systems, bundle equals YAML, every target exists, every removed System listed.

### Deferred / not done

- **Resource type rename** (`database`, `volume`, `bucket`, `cluster`, `network`, `code-repository`, `capability`, `endpoint`): needs edits to `plugin-onemind-common/AccessDeployments.tsx`, `report-service-coverage.cjs` and `scripts/tests/test_gitops_topology.py`, which other work orders touch.
- **Pillar column.** The new-frontend catalog table has no column extension point (plugin-catalog 2.0.9). Delivered the filter instead; a column needs replacing the index page.
- **The 10 not-observed workloads** (`missing-review-required`) are kept and mapped to their apps, not deleted.
- **Both-replica processed-entity counts** and signed-in acceptance: only possible after merge and rollout.
- **Before/after screenshots** (rule 10): the filter UI has not been rendered in a signed-in portal.
- **Pipeline wiring.** `pipeline.py` lists the app tests by name; the new filter tests live inside `catalogCompositionModules.test.ts` so the existing step runs them.

### Plugin-first record (rule 11)

- Native Backstage Domain / System / `subcomponentOf` model: used as is.
- `@terasky/backstage-plugin-kubernetes-ingestor` 4.3.0 (TeraSky, published 2026-10-21): generates Components from live Kubernetes objects. Not adopted here: F11 owns live ingestion, and F2's mapping must stay reviewed and offline.
- Built-in `@backstage/plugin-catalog-react` pickers (`EntityTagPicker`, kind and type pickers, `CatalogFilterBlueprint`): reused. They cannot filter on "label missing" or on saved presets, so a thin custom filter (about 150 lines, unit tested) covers only that gap.
- Argo CD plugin `@backstage-community/plugin-argocd` (installed): uses one `argocd/app-name` per Component, hence the primary-plus-list convention.
- Web directory (backstage.io/plugins) was not browsed in this pass; npm registry was queried for the ingestor only.

### Rollout notes

- Catalog files: `apps.yaml` 175 KB to 227 KB and `cluster.yaml` 454 KB to 503 KB; ConfigMap total 730,864 bytes of the 1 MiB limit.
- The `Domain` rule is in `platform/backstage/app-config.cluster.yaml` (ConfigMap, no image). Backstage must restart to load it. Until then the five Domains are rejected and the Systems' `spec.domain` references dangle.
- The filter UI is image-borne; the catalog restructure alone is Argo-only.

### Gate

- `yarn install --immutable`: ok.
- `pipeline.py check --profile portal --base origin/main`: `check_passed` (only the optional `whole-tree-diagnostic` fails, known debt).
- Tests: `test-topology.cjs` (5 topology + 15 structure + 1 snapshot guards) with `test-audit-catalog.cjs` (13): 34 pass; launcher-review and service-coverage: 17 pass; app jest `catalogCompositionModules.test.ts` 49 pass (10 new).
