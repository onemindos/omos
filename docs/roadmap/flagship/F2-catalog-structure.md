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

(Implementing agent: write progress here.)
