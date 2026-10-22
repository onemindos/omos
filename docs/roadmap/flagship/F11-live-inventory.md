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

Built by an implementing agent, completed and reviewed by the coordinator (2026-10-22). Source only, **disabled by default**: it changes
nothing until Zeus sets `onemind.inventory.live.enabled: true` (runbook in [live-inventory.md](../../catalog/live-inventory.md)).

**Plugin evaluation (rule 11, stage 1)** TeraSky `@terasky/backstage-plugin-kubernetes-ingestor` (latest 4.3.0, 2026-10-21) was
evaluated and **not adopted**, for two reasons that configuration cannot remove:
1. *Version cohort.* No release matches our pinned backend packages. 4.0.0 requires `@backstage/backend-defaults ^0.17.4` (we pin 0.18.0;
   a caret on 0.x excludes it) and 4.3.0 requires `^0.18.1`, `backend-plugin-api ^1.11.0`, `config ^1.3.10`, `plugin-kubernetes-backend ^0.21.12`,
   `plugin-catalog-node ^2.2.6` and more. Adopting it means bumping the portal's whole backend cohort (auth, permissions, catalog) or running
   duplicate copies of core backend packages, which the roadmap risk register says to avoid without a reverification pass.
2. *Identity.* The portal's catalog identities, favorites, links, F2's workload map and redirects all rest on the stable
   `k8s-<namespace>-<name>-<hash8>` names. The ingestor names entities by its own model and cannot reproduce the hash scheme, so it would
   change every workload ref or create duplicates beside `cluster.yaml`.
The Kubernetes Entity Provider (Antoine Dao) and Red Hat options were not considered further for the same two reasons (no stable naming,
generic Component model). The custom piece is small: one catalog module (about 900 lines plus tests) and a generated mapping bundle.

**Done**
- `@internal/plugin-onemind-inventory-backend`: reads Deployments and StatefulSets (identity and selector only; no env, images, annotations
  or Secrets) with the pod's read-only ServiceAccount every 2 minutes. Publishes only workloads that `catalog/cluster.yaml` does not define yet,
  using the same name scheme as `snapshot-cluster.cjs` (proved byte for byte against all 176 live workloads), placed under their app and System
  by a port of F2's rules, or under Unassigned workloads. Removed provider-owned workloads stay `lifecycle: removed` for 24 h, then drop.
  Reviewed (file-defined) entities are never replaced: a processor only stamps last-seen and flags a vanished one `missing-review-required`.
- Safety rails: an empty or failed cluster read changes nothing; excluded namespaces (`kube-node-lease`, `kube-public`, `dp-*` which OpenChoreo
  already publishes); a full mutation that can only remove entities this provider owns; flag off registers nothing and makes no cluster call.
- Guards in the existing topology gate (`scripts/test-live-inventory.cjs`): mapping bundle matches `workload-map.yaml`, flag default is off,
  and `backstage-observer` never gains a write verb, wildcard, Secrets or ConfigMaps. They skip cleanly inside the image build context.
- **No RBAC change**: `backstage-observer` already grants get/list/watch on deployments and statefulsets cluster-wide (verified in
  `platform/backstage/rbac.yaml`).

**Tests** inventory backend 28, catalog node tests 65 on the host and 58 pass + 7 skipped in the simulated image context, whole-app render
and tsc clean.

**Needs Zeus** Turning it on is a config line plus a restart (steps in the doc). Decision recorded there: reviewed entities are flagged, never
silently dropped, because removing a reviewed entity should be a reviewed Git change; automatic removal would be a larger cutover.

**Not verified** Against the live cluster. Acceptance after enabling: both replicas' entity counts and an inventory-checked-at age under
5 minutes, plus the expected log lines in the doc.

**Enabled 2026-10-22 (Zeus's request)** in `platform/backstage/app-config.cluster.yaml` (`onemind.inventory.live.enabled: true`),
merged right after the image that contains the module was verified live. Rollback: set it to `false` and let Argo sync.
