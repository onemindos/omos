# Catalog structure: how it is maintained

Proposal and decisions: [system-map.md](system-map.md). Specification: [onemind-catalog-pillars.md](../roadmap/onemind-catalog-pillars.md).

Reviewed sources (in `services/onemind-portal/catalog/`, not Backstage locations):

| File | What it decides |
|---|---|
| `structure.yaml` | Domains, Systems (title, description, pillar), new curated apps, placement of existing apps, old-System fallbacks, pillar overrides, title overrides |
| `workload-map.yaml` | `<namespace>/<workload>` to app, and every Argo CD application to exactly one app |
| `redirects.yaml` | Systems removed in the restructure and where they went |
| `ref-baseline.json` | Every entity ref that existed before the restructure; none may vanish without a redirect |

`scripts/curate-structure.cjs` applies them to `apps.yaml` and `cluster.yaml` (and the mounted copies in `platform/backstage/catalog/`). It is idempotent, byte-stable, never renames an entity and runs last in `refresh-inventory.cjs`. `snapshot-cluster.cjs` applies the workload part on a fresh snapshot.

Adding an app or Argo application: add the Argo name under `argoApplications` and its workloads under `workloads` in `workload-map.yaml`, add the app to `structure.yaml` if it has no Component yet, run `node scripts/curate-structure.cjs`, then `node --test scripts/test-topology.cjs`. A new Argo application with no mapping fails the test; a new workload with no mapping appears in the "Unassigned workloads" System.
