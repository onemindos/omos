# Live inventory: the catalog follows the cluster (F11)

Status: **built, disabled by default.** Nothing here changes the catalog until `onemind.inventory.live.enabled` is set to `true`.
Work order: [F11](../roadmap/flagship/F11-live-inventory.md). Package: `services/onemind-portal/plugins/plugin-onemind-inventory-backend`.

## What it does

Every 2 minutes (configurable) the portal reads the cluster's Deployments and StatefulSets with its own read-only ServiceAccount and does three things:

| Case | Result |
|---|---|
| A workload exists in the cluster and in `catalog/cluster.yaml` (176 of 176 today) | The reviewed entity is kept exactly as in Git. A catalog processor only stamps `onemind.io/inventory-observed-at` (last seen) and `onemind.io/inventory-last-checked-at` (last check) and sets `onemind.io/inventory-state: observed`. |
| A workload exists in the cluster but not in `cluster.yaml` (new) | The provider publishes a Component with the same name scheme as the snapshot script, attached to its app by `catalog/workload-map.yaml`, or to **Unassigned workloads** when no mapping exists yet. |
| A workload is gone from the cluster | A provider-published entity gets `lifecycle: removed` and `onemind.io/inventory-removed-at` for 24 hours, then is dropped. A **reviewed** entity from `cluster.yaml` gets `lifecycle: removed` and `inventory-state: missing-review-required` and stays (see "Decision for Zeus" below). |

Excluded everywhere: namespaces `kube-node-lease` and `kube-public`, and every `dp-*` namespace (OpenChoreo data planes, which stay excluded until that platform is torn down). `kube-system` stays tracked because F2 maps its workloads to apps.

Safety rails: a failed or empty cluster read changes nothing (it never mass-removes); freshness is only stamped from an observation younger than 3 cycles; the provider never publishes a name that `cluster.yaml` already defines.

## Identity (no ref churn)

`entityName(ns, name)` is `k8s-<ns>-<name>` cut to 54 characters, `-`, and the first 8 hex of `sha256("<ns>/<name>")`: byte for byte the scheme of `scripts/snapshot-cluster.cjs`. Tests prove it for all 176 live workloads against `cluster.yaml`, against the script's own output and for the portal itself (`k8s-backstage-backstage-21b99354`). F2's `redirects.yaml` is not needed because no Component name changes; F2's Systems, Domains and `ref-baseline.json` are untouched.

Placement (`system`, `subcomponentOf`, title, pillar label and review, role label) is computed by a port of `scripts/curate-structure.cjs`'s rules over `src/mapping.generated.ts`, which `scripts/generate-inventory-mapping.cjs` writes from `workload-map.yaml` and `apps.yaml`. The mounted catalog ConfigMap carries only `apps.yaml`, `cluster.yaml` and `org.yaml`, so the bundle avoids reading `workload-map.yaml` at runtime. Tests run the port and the curator side by side on every live workload and compare the result with the committed file; another test fails when the bundle is stale. After editing `workload-map.yaml` or `apps.yaml` run:

```bash
cd services/onemind-portal
node scripts/curate-structure.cjs            # existing F2 flow
node scripts/generate-inventory-mapping.cjs  # refresh the bundle
node --test scripts/test-topology.cjs
```

## Turning it on (Zeus)

1. Merge and release the image that contains this package (a new backend module needs an image).
2. In `platform/backstage/app-config.cluster.yaml` (ConfigMap, no image) add:

   ```yaml
   onemind:
     inventory:
       live:
         enabled: true
   ```

   Restart Backstage (both replicas read config at startup). Other keys are in `app-config.yaml` as comments.
3. Check the log line `Live inventory enabled: observing every 120s` and, after a minute, `Live inventory: 175 workloads tracked, 175 defined in catalog/cluster.yaml, 0 published`.
4. Record the entity counts and the `inventory-checked-at` age on both replicas (acceptance: at most 5 minutes).
5. To go back: set `enabled: false` (or remove the key) and restart. Provider-published entities disappear on the next processing cycle; reviewed entities return to their file values.

## RBAC

No change. The existing `backstage-observer` ClusterRole already grants `get`, `list`, `watch` on `apps` `deployments` and `statefulsets` cluster-wide; the inventory uses only `list` on those two. It reads no Secrets, no ConfigMaps and no other kind, and writes nothing. A test (`scripts/test-live-inventory.cjs`) fails if that ClusterRole ever gains a write verb, a wildcard, Secrets or ConfigMaps.

## Decision for Zeus: reviewed entities are flagged, never silently dropped

`catalog/cluster.yaml` entities carry hand-reviewed links, dependencies, runbooks and Argo mapping. A catalog processor cannot delete a file-defined entity, and a provider that republished them would lose that review. So only provider-published entities follow the "24 hours, then dropped" rule. A reviewed workload that vanishes is flagged `lifecycle: removed` and `missing-review-required` (10 today) and removal from the file is a reviewed Git change (F8's hygiene report lists them). If you prefer automatic removal of reviewed entities too, that needs a larger cutover (the provider owns `cluster.yaml`'s workload documents and the file location is removed); say so and it becomes a follow-up.

## Not in this release

PVCs, Ingresses and Argo CD applications are not ingested (the snapshot script never read them either). Crossplane and KRO kinds are out of scope. `snapshot-cluster.cjs` stays for one release as the fallback, then goes in a follow-up once you agree.

## Why not the TeraSky Kubernetes Ingestor

See the evaluation at the top of the F11 work order's Status. In short: it cannot reproduce the `k8s-<ns>-<name>-<hash>` names (its `nameModel` offers only `name`, `name-kind`, `name-cluster`, `name-namespace`, `uid`), it creates one System per namespace (F2 caps the portal at 25 Systems), it takes owner, System and `subcomponentOf` only from annotations on the Kubernetes objects themselves (a change to about 176 manifests in other work orders' paths), it deletes immediately, and its newest release needs `@backstage/backend-plugin-api` 1.11 while the portal is pinned to 1.10.1.
