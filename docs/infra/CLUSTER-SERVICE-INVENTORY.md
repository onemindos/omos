# Live Cluster Service Inventory

> Snapshot captured 2026-10-06 from Kubernetes context `do-nyc3-omos-cloud`.
> This describes live cluster state at capture time. Argo CD manifests in Git are the
> intended owners; service and pod state can change between snapshots.

## Changes after this snapshot

- Homepage was retired later on 2026-10-06 and replaced by `platform-backstage`
  (`platform/backstage`, namespace `backstage`). The historical Homepage entries below
  describe the earlier capture, not an active service or a recommendation to restore it.
- On 2026-10-07, all 19 current Argo CD Applications were checked: 16 were
  Synced/Healthy, `platform-takserver` was Synced/Degraded, and `app-of-apps` plus
  `platform-clickhouse-operator` were OutOfSync/Healthy before PRs #34 and #35.
- Backstage had two available application replicas and one ready PostgreSQL instance.
  Service and node-utilization counts below have not been recaptured and must not be
  treated as current measurements. PVC state was rechecked after the TAK-ML claim removal.
- On 2026-10-07, a read-only inspection found about 1.1 GiB of Python ML packages on
  `takml/takml-deps`. No current pod mounts it; the TAK-ML vision image now installs its
  dependencies in the image. PR #56 removed the claim; Argo CD synced revision
  `6078b366` as Healthy, and both the PVC and PV are now absent. A local recovery archive
  was validated before removal.
- Rechecked read-only on 2026-10-21: `takml/takml-deps` still returns NotFound, no pod in
  `takml` references the claim, no manifest under `platform/` mentions it, and
  `platform-takml` is Synced/Healthy. No further deletion is pending for this claim.

## Summary

- 19 Argo CD Applications are present: 16 `Synced/Healthy`, one `Synced/Degraded`,
  one `OutOfSync/Healthy`, and the `app-of-apps` root is `OutOfSync/Healthy`.
- 129 Kubernetes Services exist across 20 namespaces. 115 currently have endpoint
  addresses, 7 are `ExternalName` aliases, and 7 have no endpoint addresses.
- 5 Deployments have desired replicas set to zero.
- At the 2026-10-06 snapshot, 33 PVCs were `Bound`; three were unmounted (205 GiB
  requested). After PR #56 retired `takml-deps`, the 2026-10-07 check still found 33
  `Bound` PVCs, with two unmounted Onetak claims (200 GiB requested). Requested capacity
  is not measured filesystem usage.
- The three cluster nodes each have 8 vCPU and 16 GiB memory capacity. `kubectl top`
  showed 13%, 18%, and 21% CPU and 103%, 69%, and 85% of allocatable memory at capture.
  This is a point-in-time observation, not a sizing recommendation.

## Argo CD ownership

Application source paths and status below are read from the live Argo CD Application
objects. `HEAD` means the Git-tracked branch revision shown by Argo CD at capture.

| Application | Git path or chart | Destination | Sync / health |
|---|---|---|---|
| `app-of-apps` | `platform/argocd` | `argocd` | OutOfSync / Healthy |
| `platform-auth-helm` | Helm `authentik` `2026.8.2` | `auth` | Synced / Healthy |
| `platform-auth-infra` | `platform/auth` | `auth` | Synced / Healthy |
| `platform-clickhouse-operator` | Helm `altinity-clickhouse-operator` `0.27.3` | `kube-system` | OutOfSync / Healthy |
| `platform-core` | `platform/core` | Multiple namespaces | Synced / Healthy |
| `platform-fedhub` | `platform/takserver/fedhub` | `takserver` | Synced / Healthy |
| `platform-homepage` | `platform/homepage` | `homepage` | Synced / Healthy |
| `platform-kube` | `platform/kube` | `kube` | Synced / Healthy |
| `platform-hermes-base` | `platform/hermes/base` | `hermes` | Synced / Healthy |
| `platform-legacy` | `platform/hermes/agent`; `platform/hermes/bridge` | `hermes` | Synced / Healthy |
| `platform-onemind-mcp` | `platform/hermes/mcp` | `hermes` | Synced / Healthy |
| `platform-monitoring` | `platform/monitoring` | `monitoring` | Synced / Healthy |
| `platform-nats` | `platform/nats` | `nats` | Synced / Healthy |
| `platform-observability` | `platform/data/clickhouse` | `clickhouse` | Synced / Healthy |
| `platform-omos` | `platform/omos` | `omos` | Synced / Healthy |
| `platform-onetak` | `platform/onetak` | `onetak` | Synced / Healthy |
| `platform-ops` | `platform/ops` | `ops` | Synced / Healthy |
| `platform-redpanda-connect` | `platform/data/redpanda-connect` | `omos` | Synced / Healthy |
| `platform-takml` | `platform/takml` | `takml` | Synced / Healthy |
| `platform-takserver` | `platform/takserver/deploy` | `takserver` | Synced / Degraded |
| `platform-zeus-site` | retired 2026-10-10 (Application and `platform/zeus-site` removed; cleanup tracked in `docs/operations/RETIRE-ZEUS-SITE.md`) | `zeus-site` | Retired |

### Current Argo CD exceptions

- `platform-clickhouse-operator` is Healthy but OutOfSync. Its last sync condition says
  the CRD ConfigMap annotation exceeded Kubernetes' 262,144-byte annotation limit.
- `platform-takserver` is Degraded even though the resources listed in its Application
  status are Synced. The current Application status did not identify a degraded resource
  in its reported resource list; investigate the health source before changing TAK Server.
- `app-of-apps` is OutOfSync/Healthy. The live status has no condition explaining the
  difference in this snapshot.

## Kubernetes Services by namespace

Counts are live Service objects. “With endpoints” means the Service had one or more
addresses in its Endpoints object at capture. Endpoint presence does not establish that
a route is used by a human or application.

| Namespace | Services | With endpoints | No endpoint addresses | ExternalName |
|---|---:|---:|---:|---:|
| `argocd` | 5 | 5 | 0 | 0 |
| `auth` | 3 | 3 | 0 | 0 |
| `cert-manager` | 3 | 3 | 0 | 0 |
| `clickhouse` | 2 | 2 | 0 | 0 |
| `default` | 1 | 1 | 0 | 0 |
| `hermes` | 7 | 1 | 0 | 6 |
| `homepage` | 1 | 1 | 0 | 0 |
| `ingress-nginx` | 2 | 2 | 0 | 0 |
| `kube` | 2 | 2 | 0 | 0 |
| `kube-system` | 16 | 12 | 4 | 0 |
| `matrix` | 4 | 4 | 0 | 0 |
| `monitoring` | 13 | 13 | 0 | 0 |
| `nats` | 9 | 9 | 0 | 0 |
| `omos` | 2 | 1 | 1 | 0 |
| `onetak` | 19 | 17 | 2 | 0 |
| `openclaw` | 1 | 1 | 0 | 0 |
| `tailscale` | 20 | 20 | 0 | 0 |
| `takml` | 2 | 2 | 0 | 0 |
| `takserver` | 15 | 14 | 0 | 1 |
| `zeus-site` (retired 2026-10-10, pending namespace deletion) | 2 | 2 | 0 | 0 |
| **Total** | **129** | **115** | **7** | **7** |

The 7 Services without endpoint addresses are four kube-system control-plane metrics
Services, `omos/omos-mission-engine`, and `onetak/nominatim` plus `onetak/overpass`.
The six Hermes and one TAK Server `ExternalName` Services are aliases rather than
in-cluster workload endpoints.

## Deployments scaled to zero

These workloads consume no Deployment pod CPU or memory while at zero replicas. Their
Services and, where applicable, persistent volumes still exist.

| Namespace / Deployment | Git owner | Manifest note / observed dependency |
|---|---|---|
| `cloudflare/cloudflared` | `platform-core`, `platform/core/cloudflared.yaml` | Explicitly documented in `CLAUDE.md` as cold standby; routing failover must be understood before removal. |
| `nats/nats-metrics-scrape` | `platform-monitoring`, `platform/monitoring/scraper-deployment.yaml` | Manifest says its image is not built and to re-enable after adding it to CI. No PVC. |
| `omos/omos-mission-engine` | `platform-omos`, `platform/omos/mission-engine-deployment.yaml` | Manifest marks it planned; image is not built. No PVC. |
| `onetak/nominatim` | `platform-onetak`, `platform/onetak/nominatim.yaml` | 100 GiB retained claim; an Ingress route points to this Service. Confirm whether the geocoding endpoint is still required. |
| `onetak/overpass` | `platform-onetak`, `platform/onetak/overpass.yaml` | Manifest says it OOMed and is only to be started for mission planning; 100 GiB retained claim and Ingress route. |

## Persistent volume audit

All 33 claims are Bound. These three were not referenced by any currently running pod's
volumes. A claim absent from running pods can still be intentionally reserved or needed
by a future workload. The requested sizes below are not filesystem usage measurements.

| Namespace / PVC | Requested | Reclaim policy | Git owner / reference | Assessment |
|---|---:|---|---|---|
| `onetak/nominatim-data` | 100 GiB | Retain | `platform-onetak`; declared and mounted by the zero-replica Nominatim Deployment in `platform/onetak/nominatim.yaml` | Idle while Nominatim is stopped. Data and geocoding need are unverified. Retained PV means claim deletion alone would not release the volume. |
| `onetak/overpass-data` | 100 GiB | Retain | `platform-onetak`; declared and mounted by the zero-replica Overpass Deployment in `platform/onetak/overpass.yaml` | Idle while Overpass is stopped. Data and mission-planning need are unverified. Retained PV means claim deletion alone would not release the volume. |
| `takml/takml-deps` | 5 GiB | Delete | `platform-takml`; declared in `platform/takml/pvcs.yaml` at snapshot | Read-only inspection found 1.1 GiB of Python packages. The current vision image bakes dependencies into the image and no pod mounted this claim. Removed via PR #56 on 2026-10-07; Argo CD pruned the PVC and the PV was deleted. |

## Review queue

These are investigation items, not approved deletions. Start with one at a time and update
this inventory when an owner confirms the intended state.

1. **Decide whether the Nominatim service is needed.** Check route traffic and geocoding
   consumers, then decide whether to keep the 100 GiB data, rebuild it later, or retire it.
2. **Decide whether Overpass is still part of mission planning.** Check consumer usage and
   whether the data can be rebuilt before considering its 100 GiB retained volume.
3. **Remove only confirmed unbuilt features.** NATS metrics scrape and mission engine are
   manifest cleanup candidates if their features are no longer planned. Cloudflared is
   explicitly a standby and stays out of that queue until routing failover is reviewed.
4. **Investigate node memory and the TAK Server / Argo CD health exceptions.** Current node
   memory is high on two nodes; this snapshot alone is not enough to safely resize nodes.

The 2026-10-06 inventory capture changed no cluster resources. The `takml-deps` claim and
its DigitalOcean volume were retired later via PR #56 on 2026-10-07.
