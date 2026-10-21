# F11 — Live inventory: the catalog follows the cluster automatically

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

## Build

1. **Entity provider** (in `plugin-onemind-machine-backend`, or a new
   `plugin-onemind-inventory-backend`). It uses the existing read-only cluster
   identity to list Deployments, StatefulSets, DaemonSets, CronJobs, PVCs,
   Ingresses/HTTPRoutes and Argo Applications on a 2-minute schedule, plus watch
   events where cheap.
2. **Same identities.** Emit the same entity names as `snapshot-cluster.cjs`
   (hash scheme), so refs, favorites and links do not change. Apply F2's
   `workload-map.yaml` for `subcomponentOf`. Unmapped objects go to "Unassigned
   workloads".
3. **Deletion.** When an object disappears, mark the entity `lifecycle:
   removed` for 24 hours with an annotation, then drop it (full-mutation
   provider semantics). Show removed items in F8's hygiene report.
4. **Exclusions.** Exclude system namespaces (`kube-system`, `kube-node-lease`,
   etc.) and OpenChoreo-generated `dp-*` namespaces, which the OpenChoreo
   provider already covers. Never emit both for one workload (test).
5. **Retire the script.** Keep it for one release as fallback, then remove it
   in a follow-up PR after Zeus agrees.
6. **RBAC.** List and watch only on the kinds above; no Secrets, no ConfigMap
   data. Write the ClusterRole into GitOps, reviewed in this PR as a separate
   commit, and do not apply it by hand.

## Files

`plugins/plugin-onemind-machine-backend/**` (or the new inventory backend),
`packages/backend/src/index.ts`, `platform/backstage/rbac.yaml` (list/watch
additions only), the catalog location config, tests with fixture API
responses, and `docs/catalog/live-inventory.md`.

## Acceptance

- [ ] On fixtures, the provider output equals the current `cluster.yaml` entity
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
