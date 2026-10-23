# DigitalOcean lock-in inventory

Verified read-only on 2026-10-10 against cluster `omos-cloud` (nyc3, Kubernetes
1.36.0-do.3) and `main`. "Effort" is for one operator: S = under a day,
M = 1-3 days, L = a week or more.

Plain-language idea: Kubernetes itself is portable. What is NOT portable is every
place where our files say "use DigitalOcean's disks / load balancers / buckets".
There are only a handful of such places, listed below.

## Summary

| # | Dependency | Count (live) | Portable replacement | Effort |
|---|---|---|---|---|
| 1 | DOKS managed control plane (HA), auto-upgrade | 1 cluster, 5 nodes, 49 namespaces | k3s or Talos (see HOME-CLUSTER-PLAN) | M |
| 2 | StorageClasses `do-block-storage*` (CSI `dobs.csi.digitalocean.com`) | 4 classes (2 used); **95 PVs, 1,182 GiB** | Longhorn / Rook-Ceph / local-path | L |
| 3 | DO LoadBalancers + `do-loadbalancer-*` annotations | 2 LBs (ingress-nginx, TAK) | Cloudflare Tunnel (web), Tailscale (admin/TAK), MetalLB or host ports | M |
| 4 | DO Spaces (S3) for every backup | 1 live bucket `onemind-backups` (nyc3); sfo3 mirror bucket prepared but copy job suspended; about 51 credential consumers | MinIO / Garage / Backblaze B2 / Cloudflare R2 | S-M |
| 5 | VolumeSnapshotClasses (CSI snapshots; used by VolSync) | 2 classes, 49 VolumeSnapshots, 46 VolSync sources using `copyMethod: Snapshot` | Longhorn/Ceph snapshot classes (or VolSync `Direct`/`Clone`) | M |
| 6 | DO Cloud Firewalls | 3 (2 DOKS-managed, 1 `onemindos-cloud-fw` attached to nothing) | Router/host firewall, Tailscale ACLs | S |
| 7 | Node pool labels in manifests | `onemind.io/compute-role: ecosystem` in about 20 manifest lines; `doks.digitalocean.com/node-pool` in 0 manifests | Label home nodes with the same custom label | S |
| 8 | DO-specific Service annotations | 2 live Services (1 defined in Git: `tak-lb-proxy.yaml`; 1 on the Helm-installed ingress-nginx Service) | Drop annotations | S |
| 9 | `doctl` in scripts | 0 in-cluster jobs; local use in `platform/install.sh`, `scripts/manage_ecosystem_workers.py`, `platform/matrix/deploy-matrix.sh`, legacy runbooks | Remove / replace per section 9 | S |
| 10 | Terraform `infra/` (DO provider) | 2 roots: `infra/` (VPC, cluster, project), `infra/ecosystem-workers/` (node pool) | Not needed for k3s/Talos | S |
| 11 | DO-hosted system images | 12 images from `ghcr.io/digitalocean-packages/*` (all DOKS-managed, kube-system) | Disappear with DOKS | none |
| 12 | DOCR (container registry) | 0 in use (images are GHCR + in-cluster Harbor) | none | none |
| 13 | Managed databases | 0 (CNPG + in-cluster Postgres) | none | none |

Everything else (Argo CD, Authentik, NATS, TAK, CNPG, VolSync, Velero, Tailscale,
Cloudflare, GHCR images) is already provider-neutral. Images: 152 of the running
container references come from `ghcr.io`, none from DOCR.

## 1. DOKS control plane

- Managed API server, etcd and upgrades. Cilium and CoreDNS are installed by DO
  (`kube-system`). Public API endpoint, no control-plane firewall (review H3).
- A new cluster must bring its own CNI, CoreDNS, metrics-server and CSI driver
  (k3s and Talos ship most of these).
- Alert rules KubeSchedulerDown / KubeControllerManagerDown / KubeProxyDown are
  false alarms on DOKS today (review A9) but become REAL signals on k3s/Talos,
  where the control plane is ours.
- etcd backup becomes our job at home (k3s: scheduled etcd snapshot to S3; Talos:
  `talosctl etcd snapshot`). Include it in the quarterly drill.

## 2. Storage

Live (`kubectl get sc`, `kubectl get pv`):

| StorageClass | Reclaim | PVs | GiB |
|---|---|---|---|
| `do-block-storage` (default) | Delete | 65 | 422 |
| `do-block-storage-retain` | Retain | 30 | 760 |
| `do-block-storage-xfs`, `-xfs-retain` | | 0 | 0 |
| **Total** | | **95** | **1,182** |

- `do-block-storage` appears in about 37 files (`grep -rIl do-block-storage .`):
  `platform/onetak/*`, `platform/auth/postgres.yaml`, `platform/n8n/*`,
  `platform/matrix/base/*`, `platform/backstage/postgres.yaml`, `platform/takml/pvcs.yaml`,
  `platform/hermes/agent/pvc.yaml`, `platform/data/clickhouse/clickhouse.yaml`,
  `platform/values/{loki,tempo}.yaml`, `platform/ecosystem/values/{ollama,dapr,opensearch,harbor}.yaml`,
  `platform/ecosystem/resources/postgres-cluster.yaml`, the VolSync sources and
  restore tests under `platform/backups/`, and `platform/takserver/DEPLOY-RUNBOOK.md`.
- DO volumes are ReadWriteOnce, single-region. They cannot be moved; data moves
  through backups (VolSync restic repos, logical dumps).
- 760 GiB is on the `-retain` class: deleting a PVC leaves a billed DO volume behind.
  That is a safety net while on DO, not a backup (it is not portable).
- Replacement, easiest first: on the new cluster create StorageClasses NAMED
  `do-block-storage` and `do-block-storage-retain` that point at the new provisioner.
  Then no manifest edits are needed for the first rebuild. Rename in Git later in one
  reviewed sweep.
  - 1 home node: `local-path` (built into k3s) or OpenEBS LVM. Fast, no replication.
  - 3 home nodes: Longhorn (simplest replicated option, snapshots, S3 backup) or
    Rook-Ceph (stronger but a lot more to learn; PR #142 stages it).
- Effort L: data must be restored per PVC (REBUILD-FROM-ZERO) and some sizes should be
  fixed during the move (`valhalla-data` 99% full, Loki 89%).
- Dormant volumes you may not want to carry: `onetak/nominatim-data` and
  `overpass-data` (2 x 100 GiB, unattached, apps scaled to 0). Do not delete them
  without explicit confirmation.

## 3. Load balancers and annotations

| Service | Kind | Notes |
|---|---|---|
| `ingress-nginx/ingress-nginx-controller` | **DO LB** `143.244.220.155` | annotation `service.beta.kubernetes.io/do-loadbalancer-type: REGIONAL_NETWORK`, Helm-CLI release |
| `nats/omos-nats`, `takserver/tak-tailnet`, `openchoreo-control-plane/cilium-gateway-gateway-default` | Tailscale (`loadBalancerClass: tailscale`) | already portable |

The TAK DO LB (`tak-external`, `152.42.155.95`) was removed on 2026-10-23; TAK now rides the
tailnet only (`tak-tailnet`). The remaining public IP belongs to DigitalOcean and cannot move;
Cloudflare DNS must be repointed when it does. EUDs use `tak.tail717ba0.ts.net`, not an IP.

Replacements:
- Web hosts (auth, ws, onetak, geo): Cloudflare Tunnel (the staged cloudflared
  manifests were removed on 2026-10-23; recover them from Git history, see NETWORKING.md). No inbound ports at home. Effort S.
- TAK raw TCP/UDP (8089 mTLS, 8443, 8446, 9101/9102, 64738 Mumble): a Cloudflare
  Tunnel does NOT carry these. Options: Tailscale only (`tak-tailnet`, staged);
  or a tiny VPS with a static IP forwarding to home over Tailscale (a cheap "front
  door"); or router port-forwarding with a static IP or dynamic DNS. Effort M. This
  is the hardest networking decision; settle it in Phase 2.
- LoadBalancer IPs at home: k3s `servicelb` (default) or MetalLB.

## 4. DO Spaces (all backups)

Endpoint `nyc3.digitaloceanspaces.com`, bucket `onemind-backups`:

| Prefix | Writer | Configured in |
|---|---|---|
| `velero/cloud/` | Velero | `platform/ecosystem/values/velero.yaml` |
| `cnpg/onemind-data/` | CNPG Barman plugin | `platform/ecosystem/resources/postgres-object-store.yaml` |
| `backups/<name>/` | platform-backup CronJobs (TAK Postgres, NATS, Authentik, n8n, matrix, takml, backstage, cluster-recovery) | `platform/backups/*.yaml` (`S3_REGION: nyc3`) + sealed `platform-backup-credentials` per namespace |
| `backups/restic/<ns>/<pvc>/` | VolSync, 46 sources | `platform/backups/volsync/sources.yaml` + sealed secrets |
| `backups/clickhouse/` | legacy ClickHouse native job (suspended) | `platform/backups/backups.yaml` |
| `mirror/` in `onemind-backups-sfo3` | `backup-offsite-copy` | `platform/backups/offsite-copy.yaml`. **Live: `suspend: true`, never run** |

Every archive is client-side encrypted (AES-256-GCM or restic), so the object store is
only dumb storage and can be swapped. Replacement is a change of endpoint, region and
credentials (OFFSITE-COPY.md "Moving to Backblaze B2 or a home server" has the recipe).
Recommended: enable the existing offsite copy to Backblaze B2 or Cloudflare R2 NOW
(before moving anything), and add Garage or MinIO at home later. Effort S-M; most of
it is resealing about 51 credential consumers.

## 5. Snapshots (CSI)

- VolumeSnapshotClasses: `do-block-storage` (Delete) and `onemind-backup-retain`
  (Retain), driver `dobs.csi.digitalocean.com`.
- 49 VolumeSnapshots exist (baseline sets of 2026-10-08/09). They only work on DO.
  They are a safety net while we stay, nothing more.
- All 46 VolSync sources use `copyMethod: Snapshot` with `do-block-storage` as storage,
  snapshot and cache class. VolSync's OUTPUT (restic repos in S3) is portable; only
  the source-side class names are DO-specific and need a home overlay.
- Velero runs with volume snapshots off: it backs up Kubernetes objects only.

## 6. Firewalls

`doctl compute firewall list`: two DOKS-managed firewalls and `onemindos-cloud-fw`
attached to nothing. No control-plane firewall. At home use the router firewall
(default deny inbound) plus Tailscale ACLs (`platform/core/policy.hujson`). Effort S.

## 7. Node pool labels

- Live pools: `onemind-hub-core-v2` (3 nodes) and `onemind-ecosystem-cpu` (2 nodes,
  label `onemind.io/compute-role=ecosystem`).
- `onemind.io/compute-role: ecosystem` is used as a nodeSelector in about 20 places
  (`platform/ecosystem/values/*.yaml`, `charts/onemind-core`,
  `resources/postgres-cluster.yaml`, `backups/volsync/restore-tests/cnpg.yaml`).
  It is a custom label set by `infra/ecosystem-workers/main.tf`. At home, label the
  node(s) `onemind.io/compute-role=ecosystem` and everything schedules.
- `doks.digitalocean.com/node-pool` is not referenced by any manifest.

## 8. DO-specific annotations

`service.beta.kubernetes.io/do-loadbalancer-name` (Git) and `-type` plus
`kubernetes.digitalocean.com/load-balancer-id` (written live; the last by DO's own
controller). Kubernetes ignores unknown annotations, so they are harmless elsewhere;
remove them in the home overlay.

## 9. doctl usage

| Where | Use | Action |
|---|---|---|
| `platform/install.sh` | legacy Helm-CLI bootstrap; header lists doctl | replace by REBUILD-FROM-ZERO |
| `scripts/manage_ecosystem_workers.py` | resizes the DO node pool via Terraform | DO-only; retire after migration |
| `platform/matrix/deploy-matrix.sh` | DNS check | move to Cloudflare |
| `platform/takserver/DEPLOY-RUNBOOK.md` | DOCR login (historical) | ignore |
| `platform/backups/OFFSITE-COPY.md` | `doctl spaces keys` | use the new provider's tool |
| In-cluster Jobs | none | none |

## 10. Terraform

`infra/` creates a VPC, a DOKS cluster and a DO project; `infra/ecosystem-workers/`
manages one node pool. STACK.md records that the defaults do not match the live
cluster. **Do not `terraform apply` without a clean `terraform plan`.** Not needed
for home clusters. If a cloud outpost stays, import the live cluster into state first.

## Re-verify these numbers (read-only)

```bash
CTX=do-nyc3-omos-cloud
kubectl --context $CTX get sc
kubectl --context $CTX get pv -o json | jq -r '.items[]|[.spec.storageClassName,.spec.capacity.storage]|@tsv' | sort | uniq -c
kubectl --context $CTX get volumesnapshotclass
kubectl --context $CTX get volumesnapshot -A --no-headers | wc -l
kubectl --context $CTX get svc -A | grep LoadBalancer
kubectl --context $CTX get nodes -L doks.digitalocean.com/node-pool,onemind.io/compute-role
doctl compute firewall list
grep -rIn 'service.beta.kubernetes.io/do-' . ; grep -rIl 'do-block-storage' .
```
