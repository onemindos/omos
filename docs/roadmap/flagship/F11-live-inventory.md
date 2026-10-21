# F11 — Live inventory: the catalog follows the cluster automatically (Kubernetes Ingestor)

Work order for one implementing agent. Branch: `flagship/f11-live-inventory`.
Program index: [README](README.md). Depends on: F2 (mapping files).

## Outcome Zeus sees

- **Live catalog.** New workloads, Argo apps, PVCs and ingresses appear in the
  catalog within minutes, already attached to the right app. Removed ones are
  flagged, then disappear.
- **Freshness.** No more hand-run snapshot script, and no stale catalog. Each
  entity shows "seen in cluster 2 minutes ago".

## Evidence

- `services/onemind-portal/scripts/snapshot-cluster.cjs` is run **by hand**
  against the `do-nyc3-omos-cloud` context. It writes `catalog/cluster.yaml`
  (514 docs) and the mounted copy under `platform/backstage/catalog/`. The
  catalog drifts until someone reruns it.
- `docs/MISSION.md`: "A periodically refreshed inventory is not the live picture."
- The portal already has a read-only Kubernetes identity (Kubernetes plugin) and
  the Argo CD integration.

## Build (plugin first)

1. **Adopt TeraSky Kubernetes Ingestor** (`@terasky/backstage-plugin-kubernetes-ingestor`,
   listed in the Backstage plugin directory). It ingests workloads (and Crossplane
   claims/XRDs later) as catalog entities on a schedule, using the existing
   read-only cluster identity. Verify the version against Backstage 1.55.3 and
   the new backend system, and record the result. Only if it cannot meet the
   points below, fall back to a small custom provider and write down why.
2. **Configure it to keep identities stable.** Either configure the naming so it
   matches today's `k8s-<ns>-<name>-<hash>` entities, or map the old names in
   `catalog/redirects.yaml` (F2) so refs, favorites and links do not break. Use
   its annotation-driven mapping (or F2's `workload-map.yaml` through a small
   catalog processor) to set `spec.subcomponentOf` and `spec.system`. Unmapped
   objects go to "Unassigned workloads".
3. **Exclusions.** Exclude system namespaces and OpenChoreo `dp-*` namespaces,
   which the OpenChoreo provider already covers. Never emit two entities for one
   workload (test).
4. **Deletion.** Confirm the ingestor removes entities when objects disappear,
   and surface removals in F8's hygiene report.
5. **Retire the script.** Keep `snapshot-cluster.cjs` for one release as a
   fallback, then remove it in a follow-up after Zeus agrees.
6. **RBAC.** List and watch only what the ingestor needs; no Secrets. Put the
   ClusterRole in GitOps as a separate reviewed commit, never applied by hand.

## Files

`packages/backend/package.json` and `yarn.lock` (ingestor), a small processor in `plugins/plugin-onemind-machine-backend/**` only if needed,
`packages/backend/src/index.ts`, `platform/backstage/rbac.yaml` (list/watch
additions only), the catalog location config, tests with fixture API
responses, and `docs/catalog/live-inventory.md`.

## Acceptance

- [ ] The ingestor (or the fallback, with the reason recorded) produces output equal to the current `cluster.yaml` entity
      names (no ref churn).
- [ ] A new Deployment in a fixture appears; a deleted one is marked removed,
      then dropped.
- [ ] No duplicate with OpenChoreo entities (test).
- [ ] RBAC diff is list/watch-only on the named kinds (reviewed).
      Required gates pass.
- [ ] Live: after release, record entity counts on both replicas and the
      refresh age at most 5 minutes.

## Status

(Implementing agent: write progress here.)
