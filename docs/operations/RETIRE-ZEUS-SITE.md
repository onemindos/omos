# Retire zeus-site (post-merge operations)

Status: Git configuration removed 2026-10-10. Cluster objects, data and DNS are
still present until the steps below are run **with owner approval**. Nothing in
the PR deletes live data.

## What the merge does and does not do

| Argo Application | prune | On merge |
|------------------|-------|----------|
| `app-of-apps` (`platform/argocd`) | true | Deletes the child Application `platform-zeus-site` (its file is removed). The child has **no** `resources-finalizer`, so deleting it does **not** cascade: the namespace, Deployments, Services, Ingress, Certificate and PVC are **orphaned and keep running**. The site stays up until step (b). |
| `platform-backups` (`platform/backups`) | **false** | Nothing is pruned. The CronJob `backup-postgres-zeus-site`, the SealedSecrets (`platform-backup-credentials`, `ghcr`), NetworkPolicy and VolumeSnapshot `readiness-20261008-5cc6336f5206` in ns `zeus-site` are no longer in Git but stay live (nightly PG dump keeps running 06:20 UTC). |
| `platform-volume-backups` (`platform/backups/volsync`) | **false** | Nothing is pruned. ReplicationSource `volume-5cc6336f5206` (daily 11:20 UTC), its SealedSecret, ServiceAccount and NetworkPolicy stay live. |
| `platform-external-dns` | n/a | `zeusdelacruz.com` removed from `domainFilters`. Policy is `upsert-only`, so no DNS record is ever deleted by it. |

Step (b) removes all the orphans in one go. Backups keep running until then, which is
desirable: a fresh dump and VolSync snapshot exist right up to deletion.

## Live facts (read-only, 2026-10-10)

- PVC `zeus-site-postgres-pvc` (5Gi) -> PV `pvc-9f9d62d9-5dd8-41bc-b98c-b11fd297335e`,
  DO volume `96fe4b59-bdc3-11f1-83db-c284a1a8e782`.
- PVC `volsync-src-volume-5cc6336f5206-cache` (1Gi) -> PV `pvc-02f34a61-9e91-4c2d-8377-6102bc4a9f40`,
  DO volume `25d3845a-c30d-11f1-8ac9-daef5cd7beb0`.
- **PV reclaim policy: `Delete`** on both PVs (StorageClass `do-block-storage` is `Delete`).
  Deleting the PVC/namespace therefore **destroys the DO block volume and its data.**
- Last successful PG dump (`backup-postgres-zeus-site`): `2026-10-09T06:20:10Z`.
- One retained DO volume snapshot exists: `efd08439-c306-11f1-955c-3aab7d106f30`
  (VolumeSnapshotContent `snapcontent-e6ce20b5-2921-48a5-bdcd-cf82e46a8026`,
  deletionPolicy `Retain`, taken 2026-10-08, 0.07 GiB, source volume `96fe4b59-...`).
  It survives namespace deletion by design.
- DO droplet `wp-zeus` (id 605391425) is powered off.

## Steps (reviewer runs, owner approves each)

### (a) Confirm a final PG dump exists

Read only. Do not print credentials.

```sh
kubectl --context do-nyc3-omos-cloud -n zeus-site get cronjob backup-postgres-zeus-site \
  -o jsonpath='{.status.lastSuccessfulTime}{"\n"}'
```

Then confirm the newest `.verified.json` marker under the Spaces prefix
`backups/postgres-zeus-site/` (the CronJob's `BACKUP_PREFIX`; the offsite mirror is the
same prefix under `mirror/`). Do this immediately before (b) and, if the owner wants a
post-shutdown dump, trigger one manually first:
`kubectl -n zeus-site create job --from=cronjob/backup-postgres-zeus-site final-dump`.

### (b) Delete the namespace

Record the volume IDs first, then delete:

```sh
kubectl --context do-nyc3-omos-cloud get pv pvc-9f9d62d9-5dd8-41bc-b98c-b11fd297335e \
  pvc-02f34a61-9e91-4c2d-8377-6102bc4a9f40 \
  -o custom-columns=N:.metadata.name,RECLAIM:.spec.persistentVolumeReclaimPolicy,HANDLE:.spec.csi.volumeHandle
kubectl --context do-nyc3-omos-cloud delete namespace zeus-site
```

Expected: the PVs are `Delete`, so both DO volumes are removed within minutes. Verify
with `doctl compute volume list | grep -E '96fe4b59|25d3845a'` (should be empty). The
VolumeSnapshot object disappears but its `Retain` VolumeSnapshotContent and DO snapshot
remain; clean up in (c). Also confirm `platform-zeus-site` no longer exists in
`kubectl -n argocd get applications` before deleting (otherwise Argo would recreate it).

### (c) Leftover DO volume snapshots

```sh
doctl compute snapshot list --resource volume --format ID,Name,ResourceId,CreatedAt \
  | grep -E '96fe4b59-bdc3-11f1-83db-c284a1a8e782|25d3845a-c30d-11f1-8ac9-daef5cd7beb0'
```

Currently: `efd08439-c306-11f1-955c-3aab7d106f30`. Also delete the orphaned
`kubectl get volumesnapshotcontent snapcontent-e6ce20b5-2921-48a5-bdcd-cf82e46a8026`
object (Retain keeps it after namespace deletion), then the DO snapshot. Match strictly
by volume ID, not by name, and delete only with owner yes (this is the last block-level copy).

### (d) Cloudflare DNS (manual, owner)

Remove records for `zeusdelacruz.com` and `www` (and any `extdns-*` TXT ownership records)
in the Cloudflare zone, or remove the zone. external-dns is `upsert-only` and will not do it.
The Cloudflare API token scope still lists the zone; narrow it when the zone is gone.
The ACME account email in `platform/core/issuer-dns.yaml` (`zeus@zeusdelacruz.com`) is
a mailbox, not a solver; keep or change it independently.

### (e) Spaces objects

Prefixes (verify exact names with a listing first): `backups/postgres-zeus-site/` and
`backups/restic/zeus-site/` (restic repo `zeus-site-postgres-pvc`), in the nyc3 bucket
**and** the sfo3 offsite mirror (`mirror/...`, copy-only, never auto-deleted per
`platform/backups/OFFSITE-COPY.md`). Keep 30 days after step (b) (until 2026-11-09 if run
today), then delete on owner choice. Deleting restic data is irreversible and
`BACKUP_ENCRYPTION_KEY`/restic password are needed to ever read it.

### (f) Droplet `wp-zeus`

Powered off. Only with an explicit owner yes: `doctl compute droplet-action snapshot 605391425`
(wait for completion), then `doctl compute droplet delete 605391425`. Snapshots are billed;
set an expiry decision with the owner.

## Out of scope / left in place

- GitHub repo `onemindos/zeus-site` and its Backstage repository entry (repo not archived by this PR).
- Historical records (`reports/`, `docs/cluster-audit-*`, `docs/operations/*.json`, `docs/inventory/`,
  `docs/capacity/`) still mention zeus-site by design.
