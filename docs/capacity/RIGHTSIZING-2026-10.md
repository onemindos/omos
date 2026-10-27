# Right-sizing, placement and edge redundancy (October 2026)

Date: 2026-10-09. Cluster: `do-nyc3-omos-cloud`. Scope: Argo-managed workloads only.
Companion to `docs/cluster-audit-2026-10-07.md`.

## Why

- Pool `onemind-hub-core-v2` (3x s-8vcpu-16gb, fixed) was 80-88% memory *requested* (33.4 GiB of 39.9 GiB allocatable) while the
  pool `onemind-ecosystem-cpu` (2x s-8vcpu-32gb-amd, autoscale 2-4) sat at 45-59% requested and 13-25% actual.
- Cluster-wide requests were about 62 GiB against a 7-day peak working set of about 42 GiB. Most of the gap is stateless
  controllers and databases requesting 256 Mi to 1 Gi while using 10 to 190 Mi.
- Several tier-2 services with large actual memory (ClickHouse 1.7 GiB, WebODM webapp 1.6 GiB, the GitHub runner about 1.7 GiB)
  were running on the tier-1 pool.

Nothing is removed or scaled down. Every app keeps running; only requests, placement preference and Authentik replicas change.

## Method

1. Peak working set per container over 7 days:
   `max by (namespace,container)(max_over_time(container_memory_working_set_bytes{container!="",container!="POD"}[7d]))`.
2. Current requests and limits from `kube_pod_container_resource_{requests,limits}`, always filtered to
   `job="kube-state-metrics"` (OpenCost re-exports these series, so unfiltered queries double count).
3. CPU p95 over 7 days: `quantile_over_time(0.95, rate(container_cpu_usage_seconds_total[5m])[7d:5m])`.
4. Rule: new memory request = max(peak x 1.3, 64 Mi) rounded **up** to a clean step (64, 96, 128, 192, 256, 320, 384, 512 Mi).
   New CPU request = max(p95 x 1.5, 10m) rounded up to 10m/25m/50m/100m.
   Memory and CPU **limits are not changed anywhere** (none is lowered below peak x 1.5, and the few already below it are left as they are).
5. Each candidate was mapped to its Argo Application via the live `argocd.argoproj.io/tracking-id` annotation. Objects without
   the annotation were skipped (list below).
6. Only the repo source of each Application was edited (plain manifests, Helm values in `platform/ecosystem/values/`, the
   `onemind-core` chart, the Authentik `helm.values` block). Every edited file passes
   `yamllint -c .yamllint.yml`; parsed-YAML diffs against `main` show only the intended leaf changes.

### Deliberate deviations from the rule

| Case | Deviation | Reason |
|---|---|---|
| Everything in ecosystem namespaces (`onemind-*`, kyverno, spire, dapr, keda, kagent, openchoreo, ...) | 128 Mi floor for controllers, 64 Mi for tiny helpers | The ecosystem pool only exists since about 2026-10-08: **these containers have 15-27 h of history, not 7 days**. Re-check after a full week. |
| `onemind-data` (CNPG, 3 instances) | 384 Mi instead of 242 Mi | 27 h of data and `shared_buffers=256MB`; working set grows as buffers fill. |
| `nodeodm` | 512 Mi / 250m instead of 128 Mi / 10m | Burst photogrammetry job with an 8 Gi limit; idle peak is meaningless. |
| `harbor-trivy` | 256 Mi instead of 128 Mi | Loads the vulnerability DB only when scanning. |
| webodm `webapp` / `worker` | memory unchanged | Peak 1599 Mi is already above the request; the file comments record OOM kills during processing. CPU only. |
| Databases | 128 Mi floor | Cold buffers after restart. |

## Right-sizing table

Peak and p95 are the live 7-day values; "History" is how much data actually exists for that container (`27h` means the 7-day
peak is really a 27-hour peak). Limits are unchanged in this PR.

| Workload | Container | Peak (Mi) | History | Mem req (Mi) old -> new | CPU req old -> new | Mem limit (Mi) old -> new | CPU p95 | Notes |
|---|---|---|---|---|---|---|---|---|
| onetak-geo-db (TimescaleDB) | timescaledb | 138 | 7d | 1024 -> 256 | 500m -> 50m | 4096 -> 4096 (unchanged) | 17m | DB floor |
| nodeodm | nodeodm | 63 | 7d | 1024 -> 512 | 500m -> 250m | 8192 -> 8192 (unchanged) | 6m | burst exception: photogrammetry job, 8Gi limit kept |
| valhalla | valhalla | 39 | 7d | 1024 -> 128 | 300m -> 10m | 4096 -> 4096 (unchanged) | 1m | init container also 512Mi->128Mi (it set the pod's effective request) |
| tile38 | tile38 | 21 | 7d | 256 -> 128 | 200m -> 25m | 1024 -> 1024 (unchanged) | 9m | in-memory store, floor 128Mi |
| onetak-postgis | postgis | 49 | 7d | 256 -> 128 | 200m -> 25m | 512 -> 512 (unchanged) | 1m | DB floor |
| stac-api | stac-api | 56 | 7d | 256 -> 96 | 100m -> 10m | 512 -> 512 (unchanged) | 3m |  |
| tileserver-gl | tileserver-gl | 132 | 7d | 256 -> 192 | 200m -> 10m | 1024 -> 1024 (unchanged) | 0m |  |
| webodm-db | webodm-db | 38 | 7d | 256 -> 128 | 200m -> 10m | 512 -> 512 (unchanged) | 5m | DB floor |
| onetak-minio | minio | 174 | 7d | 256 -> 256 | 100m -> 10m | 512 -> 512 (unchanged) | 2m | memory unchanged (peak x1.3 = 226Mi), CPU only |
| webodm (webapp) | webapp | 1599 | 7d | 768 -> 768 | 300m -> 50m | 2048 -> 2048 (unchanged) | 7m | memory NOT changed: peak 1599Mi is above the 768Mi request; CPU only |
| webodm (worker) | worker | 242 | 7d | 537 -> 537 | 300m -> 50m | 1024 -> 1024 (unchanged) | 6m | memory NOT changed: burst OOM history (see file comments); CPU only |
| takml-vision-worker | vision-worker | 369 | 7d | 768 -> 512 | 500m -> 50m | 1024 -> 1024 (unchanged) | 2m |  |
| takml-db | postgres | 50 | 7d | 256 -> 128 | 100m -> 10m | 512 -> 512 (unchanged) | 1m | DB floor |
| zeus-site | zeus-site | 132 | 7d | 512 -> 192 | 200m -> 25m | 1024 -> 1024 (unchanged) | 13m |  |
| zeus-site-postgres | postgres | 48 | 7d | 256 -> 128 | 100m -> 10m | 512 -> 512 (unchanged) | 2m | DB floor |
| n8n postgres | postgres | 54 | 51h | 256 -> 128 | 100m -> 50m | 768 -> 768 (unchanged) | 17m | DB floor |
| onemind-data (CNPG x3) | postgres | 186 | 27h | 1024 -> 384 | 250m -> 50m | 4096 -> 4096 (unchanged) | 32m | 27h history; kept above rule because shared_buffers=256MB |
| onemind-core api (x2) | api | 79 | 27h | 512 -> 128 | 250m -> 50m | 2048 -> 2048 (unchanged) | 17m | 27h history, floor 128Mi |
| onemind-core daprd sidecar (x2) | daprd | 35 | 27h | 256 -> 128 | 100m -> 25m | 1024 -> 1024 (unchanged) | 3m | dapr.io/sidecar-* annotations; 27h history, floor 128Mi |
| onemind-ingestion | ingestion | 86 | 27h | 256 -> 128 | 100m -> 25m | 1024 -> 1024 (unchanged) | 10m | 27h history, floor 128Mi |
| openfga | openfga | 25 | 26h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 1m | 27h history, floor 128Mi |
| kyverno background/cleanup/reports (x3) | controller | 58 | 27h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 4m | 27h history, floor 128Mi |
| kyverno admission | kyverno | 60 | 27h | 256 -> 128 | 100m -> 25m | 1024 -> 1024 (unchanged) | 7m | 27h history; webhook path keeps 25m |
| spire-server | spire-server | 48 | 27h | 256 -> 128 | 100m -> 25m | 1024 -> 1024 (unchanged) | 13m | 27h history |
| spire-agent (DaemonSet x5) | spire-agent | 29 | 27h | 128 -> 64 | 50m -> 10m | 512 -> 512 (unchanged) | 2m | runs on every node incl. hub-core |
| dapr-scheduler (x3) | dapr-scheduler-server | 94 | 27h | 256 -> 128 | 100m -> 50m | 1024 -> 1024 (unchanged) | 32m | 27h history |
| dapr-placement | dapr-placement-server | 14 | 27h | 128 -> 64 | 50m -> 50m | 1024 -> 1024 (unchanged) | 21m | memory only; CPU p95 21m keeps 50m |
| dapr-sentry | dapr-sentry | 42 | 27h | 128 -> 64 | 50m -> 10m | 1024 -> 1024 (unchanged) | 2m |  |
| dapr-sidecar-injector | dapr-sidecar-injector | 16 | 27h | 128 -> 64 | 50m -> 10m | 1024 -> 1024 (unchanged) | 1m |  |
| openchoreo controller-manager | manager | 30 | 26h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 5m | 27h history, floor 128Mi |
| openchoreo api | api-server | 30 | 26h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 0m | 27h history, floor 128Mi |
| openchoreo cluster-gateway | cluster-gateway | 15 | 26h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 0m | 27h history, floor 128Mi |
| harbor-core | core | 45 | 27h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 4m | 27h history |
| harbor-jobservice | jobservice | 16 | 27h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 3m | 27h history |
| harbor-portal | portal | 7 | 27h | 128 -> 64 | 50m -> 10m | 512 -> 512 (unchanged) | 0m |  |
| harbor-trivy | trivy | 10 | 27h | 512 -> 256 | 100m -> 10m | 2048 -> 2048 (unchanged) | 1m | scan-time exception: loads the vuln DB when scanning, idle peak is 10Mi |
| kagent-controller | controller | 28 | 27h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 2m | 27h history, floor 128Mi |
| kagent-ui | ui | 70 | 27h | 256 -> 128 | 100m -> 10m | 768 -> 768 (unchanged) | 0m | 27h history |
| kagent-tools | tools | 20 | 27h | 128 -> 64 | 100m -> 10m | 512 -> 512 (unchanged) | 0m |  |
| keda-operator | keda-operator | 51 | 27h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 3m |  |
| keda-metrics-apiserver | keda-operator-metrics-apiserver | 56 | 27h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 5m |  |
| keda-admission-webhooks | keda-admission-webhooks | 15 | 27h | 256 -> 64 | 100m -> 10m | 1024 -> 1024 (unchanged) | 1m |  |
| agent-sandbox-controller | agent-sandbox-controller | 34 | 27h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 3m |  |
| argo-rollouts | argo-rollouts | 34 | 27h | 256 -> 128 | 100m -> 10m | 1024 -> 1024 (unchanged) | 2m |  |
| cnpg operator | manager | 76 | 27h | 256 -> 128 | 100m -> 25m | 1024 -> 1024 (unchanged) | 11m |  |
| opencost | opencost | 120 | 27h | 256 -> 192 | 100m -> 25m | 1024 -> 1024 (unchanged) | 8m | peak 120Mi x1.3 = 156Mi rounded up |

Sum of memory requests over the rows above (per replica x replicas): about 20505 Mi -> 9849 Mi
(**about 10.4 GiB freed** cluster-wide, including the unchanged-memory rows).

## Placement (preferred, never required)

All new affinity is `preferredDuringSchedulingIgnoredDuringExecution`, weight 100, on the DOKS pool label
`doks.digitalocean.com/node-pool`, so scheduling still works if a pool shrinks.

- Prefer **onemind-ecosystem-cpu**: `onetak` geo-db, valhalla, tile38, tileserver-gl, postgis, minio, stac-api, webodm,
  webodm-db, webodm-broker, nodeodm; `zeus-site`; `backstage` (merged into its existing pod anti-affinity); `n8n`;
  the GitHub runner (`platform/core/github-runner.yaml`); ClickHouse (`platform/data/clickhouse/clickhouse.yaml`, pod template).
- Prefer **onemind-hub-core-v2** (tier-1): `hermes-legacy` and `takserver-plugins`. (hermes-legacy is currently running on an
  ecosystem node; it will return to hub-core at its next restart.)
- Already hard-pinned by `nodeSelector: onemind.io/compute-role: ecosystem` in the existing manifests and left as is:
  `onemind-ai` (ollama, litellm), `onemind-core`, `onemind-data`, `onemind-registry`, `onemind-search` (opensearch),
  `kagent`, kyverno, spire-server, keda, velero, openchoreo, rollouts, agent-sandbox, opencost.
- Not given a preference: `takml-*` (`takml-server` is pinned by hostname to a hub-core node, probably to share an RWO model
  volume with the vision worker; needs its own change), small stateless services (martin, maputnik, tak-portal, mediamtx, ...).
- RWO-volume Deployments that previously used the default RollingUpdate now declare `strategy: Recreate`
  (geo-db, minio, tile38, tileserver-gl, valhalla, nodeodm, webodm-db). Without this, a preferred cross-node move can deadlock:
  the new pod waits on a Multi-Attach error while the old pod is kept alive by `maxUnavailable: 0`. Cost: a short outage per
  pod on the next rollout.

## Redundancy

- **Authentik server: 2 replicas** (`platform/argocd/platform-auth-helm.yaml`, `server.replicas`), chart default preferred pod
  anti-affinity already spreads them across nodes, plus a PDB `minAvailable: 1`
  (`platform/auth/authentik-server-pdb.yaml`, applied by `platform-auth-infra`). Adds 512 Mi of requests; the net change on
  hub-core is still strongly negative.
- **Authentik worker stays at 1**: it is not on the login path (background tasks, outpost sync); a pause during a drain delays
  jobs but does not take SSO down. Revisit when there is spare tier-1 headroom.
- Still single points of failure for auth: `authentik-postgres` (StatefulSet, 1 replica) and `authentik-redis`.
- The Application has `ignoreDifferences` on `/spec/template/spec/containers/0/resources`; the replica change alone will
  trigger the sync, so no manual step is needed.

## Expected requested memory per node after change

Estimated from live pod placement, new requests, and assuming every pod with a new ecosystem preference lands on an ecosystem
node when it is next recreated (the preference does not evict running pods; `Recreate` rollouts from this PR move the onetak
pieces immediately).

| Node | Before | After resizing only | After resizing and moves |
|---|---|---|---|
| hub-core ...37phka | 10.6 GiB (80%) | 9.9 GiB (74%) | 8.6 GiB (64%) |
| hub-core ...37phke | 11.1 GiB (83%) | 8.8 GiB (66%) | 6.9 GiB (52%) |
| hub-core ...37phkg | 11.7 GiB (88%) | 10.7 GiB (81%) | 7.5 GiB (56%) |
| ecosystem ...3oleu2 | 16.4 GiB (59%) | 14.2 GiB (51%) | 14.5 GiB (52%) |
| ecosystem ...3oleul | 12.5 GiB (45%) | 8.3 GiB (30%) | 14.5 GiB (52%) |
| **Total** | **62.4 GiB** | | **52.0 GiB** (hub-core 33.4 -> 22.9 GiB; +0.5 GiB if the second Authentik server lands on hub-core) |

Memory requested on hub-core never increases. CPU requests fall from 18.1 to 14.6 cores on hub-core before moves, and from 9.1 to 5.7 on
the ecosystem pool. Requests that are **still too low** (and therefore hide real usage) are deliberately left alone because the rule
forbids raising hub-core totals: ClickHouse (peak 1712 vs 512 Mi), argocd application-controller, nats (1675 vs 512), prometheus
(1908 vs 512), Authentik server (740 vs 512, limit 1 Gi), takserver postgresql (1016 vs 256), webodm webapp, takml-server, fedhub,
hermes, github-runner. Moving ClickHouse, webodm and the runner to the ecosystem pool is what actually relieves real memory
pressure on hub-core (actual use there was 74-77% on two nodes).

## Skipped (not Argo-managed, or out of scope)

No `argocd.argoproj.io/tracking-id` on the live object (applied by hand or by Helm CLI):

- `matrix` namespace: synapse (request 1 Gi, peak 184 Mi), postgres, redis, element-web. Largest single skipped saving (about 0.8 GiB).
- `takserver`: takserver-api, -messaging, -config (1.5 + 1.5 + 1 GiB requested; peaks 0.8, 0.8, 0.45 GiB), ignite, postgresql,
  nats, nats-box, `cot-retention` CronJob.
- `openclaw` gateway (request 1 Gi, actual 2.1 GiB), `argocd` itself, `cert-manager`, `ingress-nginx`, `tailscale` operator/proxies,
  most `nats` helper deployments (agent-web-ui, agenttap, heynats, nack, nats-dashboard, nats-ui).

Out of scope by instruction: `platform/monitoring/`, `platform/backups/`, `platform/values/*`, `platform/takserver/chart`,
`platform/nats` Helm values, velero (backups tier-1), cilium and other `kube-system` objects.

Not changed because they are already at or below their request, or too small to matter (under 128 Mi): martin, maputnik,
pg-featureserv, stac-browser, tak-portal, mediamtx (this file also has pre-existing yamllint errors), webodm-broker, harbor
registry/nginx/database/redis, kagent agent pods (sized by the Agent CRs, no values key), cronjob pods.

## Follow-ups (not in this PR)

1. **ingress-nginx and cert-manager** are Helm-CLI managed single replicas on the critical path. Bring them under Argo with
   2 replicas (ingress-nginx with a PDB and anti-affinity) before relying on any node-count reduction.
2. **Ecosystem autoscaler minimum 2 -> 1** once this is merged and stable for 48 h (about $168/month):
   `doctl kubernetes cluster node-pool update <cluster> onemind-ecosystem-cpu --min-nodes 1`.
   Check first: after this PR the hard-pinned ecosystem workloads need about 15.9 GiB of requests, which fits on one 28.1 GiB node;
   the 6.6 GiB of *preferred* movers would spill back to hub-core (which has the room: 22.9 of 39.9 GiB) rather than go Pending.
3. **Hub-core swap to s-4vcpu-16gb-amd** ($84 vs $96 per node, 16 GiB RAM is identical). Memory is no longer the constraint, but CPU
   requests on hub-core are still 14.6 cores before moves against about 11.8 allocatable on 4-vCPU nodes. Only do this if CPU
   requests after the moves settle under about 10 cores; create the new pool, cordon and drain the old one (one node at a time;
   NATS has 3 replicas with PDBs).
4. Raise the requests that are too low once the above moves have landed (see the list above); start with ClickHouse and Authentik
   server limits.
5. Re-run the table after 7 days of ecosystem history and tighten the 128 Mi floors if peaks stay under 50 Mi.
6. Put the skipped `matrix`, `takserver-*`, `openclaw` and `argocd` under Argo so they can be right-sized the same way.
7. `dapr-operator`, `dapr-sentry`, `dapr-sidecar-injector`, `dapr-placement` have `nodeSelector` in `platform/ecosystem/values/dapr.yaml`
   but the live pods carry none (chart value path mismatch); they run on the ecosystem pool only by luck.
8. `platform/data/clickhouse/clickhouse.yaml` contains a hard-coded ClickHouse user password in plain text (pre-existing); move it to a
   SealedSecret.

## Cost context

List price is about $828/month; the owner has about $2,500 of DigitalOcean credit, so there is roughly 3 months of runway at list. The
savings above are about capacity headroom first (the most-requested hub-core node goes from 88% to at most 64%) and about $168/month
(ecosystem min 1) later, and up to $36/month more with the hub-core swap (3 x $12).

## Rollout notes and risks

- Every changed pod template restarts its pod when Argo syncs. Suggested order: merge in a quiet window, let Argo sync, watch
  `kubectl get pods -A | grep -v Running`. The CNPG `onemind-data` cluster rolls replicas first and then switches over.
- ClickHouse, geo-db, minio, tile38, tileserver-gl, valhalla, nodeodm and webodm-db have a short outage while their `Recreate`
  rollout detaches and re-attaches the block volume.
- Low requests with unchanged high limits mean more overcommit on the ecosystem nodes. Total limits there were already far above
  capacity; node-level pressure shows up as eviction of the lowest-QoS pods first, so watch `node_memory_MemAvailable_bytes`
  on the ecosystem pool during the first week.
- 27 h of history on ecosystem containers: a rarely-used controller could burst above the new request. Their limits are unchanged, so
  the worst case is an unreserved spike, not an OOM.
- Rollback: revert the commit; Argo restores the old requests. No data migration is involved.

## Addendum: hub-core pass (2026-10-27)

The hub-core nodes were replaced on 2026-10-18 (`onemind-hub-core-v2-3rtf47/-3rtf4m/-3rtf4q`, 13,646 Mi allocatable each); the
ecosystem pool is now one node, `onemind-ecosystem-cpu-3rhwrp`. Requested memory was 79% / 96% / 71% while actual use was 57-66%. Data taken read-only at
2026-10-27 09:15 UTC: `kubectl describe nodes`, `kubectl top pods -A --containers`, and Prometheus 7-day and 10-day windows of
`container_memory_working_set_bytes` (a single kubelet/cAdvisor job, so nothing is double counted; retention is 10 days).

Rule: new memory request = 7-day p95 x 1.3, rounded up to the step list above (64 Mi floor), and never below current use x 1.25
or the 10-day peak. This is not the rule in the method above (7-day peak x 1.3, 64 Mi floor), but for the four containers changed
here both rules give the same values. Limits are not changed. CPU is changed only where it is grossly over (Synapse).

| Workload (node) | Container | 7-day p95 / 7-day peak / 10-day peak / now (Mi) | Mem req old -> new | CPU req old -> new | Owning app |
|---|---|---|---|---|---|
| matrix synapse (3rtf4m) | synapse | 181 / 182 / 199 / 169 | 1Gi -> 256Mi | 500m -> 100m (p95 10m, peak 29m) | platform-matrix (manual sync) |
| omos feed-intel (3rtf4m) | feeds | 42 / 47 / 51 / 41 | 192Mi -> 64Mi | 100m (unchanged) | platform-omos |
| onetak tak-portal (3rtf47) | portal | 44 / 44 / 54 / 43 | 128Mi -> 64Mi | 50m (unchanged) | platform-onetak |
| onetak webodm-broker (3rtf47) | broker | 4 / 5 / 9 / 3 | 128Mi -> 64Mi | 100m (unchanged) | platform-onetak |

Requested memory per hub-core node if every restarted pod comes back on its current node (see the placement note below):

| Node | Before | After |
|---|---|---|
| hub-core ...3rtf47 | 10,630 Mi (77.9%); 10,886 Mi (79.8%) while a VolSync mover runs | 10,502 Mi (77.0%); 10,758 Mi (78.8%) with the mover |
| hub-core ...3rtf4m | 13,136 Mi (96.3%) | 12,240 Mi (89.7%) |
| hub-core ...3rtf4q | 9,805 Mi (71.9%) | 9,805 Mi (71.9%) |
| **Pool** | **33,571 Mi (82.0%)** | **32,547 Mi (79.5%)** |

No node crosses 80% because of this pass. ...3rtf47 and ...3rtf4q were already under it, and ...3rtf4m, the only node over it,
stays above it.

Same placement is not the likely outcome. The scheduler puts each restarted pod on the node that scores best, and the default
scoring favours nodes with the most unrequested CPU and memory. Even after its old Synapse pod is gone, ...3rtf4m is the most
requested node (12,112 Mi, 88.8%), so it is the least likely target. Synapse uses `Recreate`, has no affinity, and its one block
volume only requires region `nyc3`. On the 2026-10-27 numbers it most likely lands on the untainted ecosystem node or on ...3rtf4q,
which score within one point of each other. If it leaves ...3rtf4m, that node ends at 11,984 Mi (87.8%; 11,920 Mi, 87.4%, if
feed-intel moves too) and the receiving node gains 256 Mi. If it lands on the ecosystem node, it takes that node's last free DigitalOcean volume slot (15 of 15). Until a volume
detaches there, any further block-volume pod that prefers the ecosystem pool is placed on hub-core instead, and the move of
hermes-legacy, ClickHouse or onemind-data-2 described below is blocked. feed-intel, tak-portal and webodm-broker have no volumes
and can move too; webodm-broker prefers the ecosystem pool (weight 100), so ...3rtf47 most likely loses its full 128 Mi.

...3rtf4m stays above 80% because what is left on it is tier-1 (TAK PostgreSQL, Ignite and NATS, the NATS cluster, Argo CD,
Authentik PostgreSQL, Hermes) or requests that are already **below** real use: hermes-legacy (tier-1) 2,128 Mi against a p95 of
3,404 Mi, ClickHouse 2 Gi against a 10-day peak of 2,644 Mi (raised on purpose in #227), Nextcloud 1,024 Mi and onemind-data-2
512 Mi against container p95s that sum to about 1,520 Mi and 615 Mi. Lowering requests cannot fix that node. Placement can:
hermes-legacy, ClickHouse and onemind-data-2 all prefer the ecosystem pool, but the single ecosystem node has 14 of its 15
DigitalOcean volume attachments in use, so pods with block volumes cannot follow that preference.

Not changed in this pass:

- `onemind-data` (CNPG): postgres is under-requested (p95 547 Mi against 384 Mi). `onemind-core` and kyverno run only on the
  ecosystem node, and SPIRE has no running pods, so none of them affects hub-core.
- matrix PostgreSQL (`postgres-0`, 256 -> 128 Mi possible, 7-day p95 67 Mi): it runs on ...3rtf47, which is already under 80%,
  and `docs/operations/matrix-drift-2026-10-23.md` (step 6) says not to restart Matrix Postgres or Redis without need. It can
  ride along with the next Postgres restart that is needed for another reason.
- `onetak-media` (mediamtx, 128 -> 64 Mi possible): `platform/onetak/media-mtx.yaml` has existing yamllint errors, and the
  changed-YAML gate would force an unrelated reformat of the file.
- vaultwarden (34 h of history, 64 Mi), n8n (the 10-day peak of its `n8n` container, 612 Mi, is above its 512 Mi request),
  matrix element-web and redis (already at the 64 Mi floor), and falcoctl (budget set on 2026-10-23 for download bursts).
- Everything tier-1 in `docs/STACK.md` (TAK Server, NATS, Authentik, ingress, cert-manager, sealed-secrets, backups and ops,
  Argo CD, Hermes). That includes hermes-legacy and hermes-fresh, and tak-ldap: `docs/STACK.md` lists tak-ldap under the tier-2
  `platform-fedhub`, but `platform/takserver/deploy/tak-ldap.yaml` is applied by `platform-takserver` (live tracking-id), the
  TAK Server overlay.
- takml (separate PRs) and DOKS-managed kube-system pods.

Rollout: no pod restarts on merge alone. `platform-matrix` is manual-sync. `platform-onetak` and `platform-omos` sync automatically
(with self-heal) but list `/spec/template/spec/containers/0/resources` under `ignoreDifferences` without
`RespectIgnoreDifferences`, so a change to resources alone leaves them Synced. Their next sync applies the new requests. Without a
deliberate sync, that happens whenever any other change under `platform/onetak` or `platform/omos` merges or self-heal reverts other
drift, at a time nobody chose, so sync both apps on purpose right after merge. Live resources matched Git for every container in
the three apps before this change, so a sync touches only the four Deployments above, and each restarts once:

- Synapse uses `Recreate`: Matrix is down for about 1-2 minutes, longer if the new pod lands on another node and its volume has
  to detach and re-attach. This follows step 6 of `docs/operations/matrix-drift-2026-10-23.md`: sync Synapse alone with
  server-side apply, and do not restart Postgres or Redis.
- The webodm Redis broker keeps no data on disk, so a restart drops queued Celery messages. Sync it while no WebODM task is running.
- feed-intel has no readiness probe, so the rolling update removes the old pod as soon as the new container starts. That container
  first runs an unpinned `pip install nats-py aiohttp` on the mutable `python:3.12-alpine` image. If PyPI is down or a new release
  breaks the script, feed-intel stays down until the install succeeds or the script is fixed. Reverting does not help, because the
  reverted pod runs the same install.
- tak-portal has readiness and liveness probes, so its surge rollout should not interrupt it.
