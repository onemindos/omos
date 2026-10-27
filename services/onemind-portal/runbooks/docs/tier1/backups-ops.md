# Backups and ops

The recovery machinery for the whole platform: logical database and stream archives, the encrypted cluster recovery export, portable volume copies, Velero namespace backups and a read-only health observer. A successful Job is not a restore: logical archives write their verification marker only after an isolated restore succeeds.

Sources: `platform/backups/README.md`, `platform/backups/volsync/README.md`, `platform/backups/OFFSITE-COPY.md`, `platform/backups/SNAPSHOTS.md`, `platform/backups/tools/README.md`, `platform/ops/README.md`, `docs/gitops/APP-FINALIZERS.md`, [Backup and recovery](../backups.md).

## What runs

| Mechanism | Argo CD application | Schedule (UTC) | What it proves |
| --- | --- | --- | --- |
| Logical archives: `backup-takserver-postgres`, `backup-nats-streams`, `backup-postgres-auth`, `-n8n`, `-matrix`, `-takml`, `-backstage` | `platform-backups` | 02:00 to 06:40 | Encrypted, uploaded, downloaded, checked and restored into an isolated sidecar or server before `.verified.json` is written |
| Cluster recovery export `backup-cluster-recovery` (namespace `ops`) | `platform-backups` | 04:00 | Every Secret and ConfigMap, including the sealing keys, and no other objects; verified by decryption and equality only |
| VolSync `ReplicationSource` objects (Restic) | `platform-volume-backups`, operator `platform-volsync` | Daily, staggered | Crash-consistent file copies; a completed copy is not an application restore |
| Velero schedule `velero-onemind-platform` | `ecosystem-velero` | 06:00, 30-day retention | Namespace backups of `onemind-core`, `onemind-data` and `onemind-ai` |
| `health-keeper` (namespace `ops`) | `platform-ops` | Every 5 minutes | Read-only observer that logs findings; it never repairs and is not an alert channel |
| Off-site copy `backup-offsite-copy` (namespace `ops`) | `platform-backups` | 17:00 | Suspended until the owner selects a destination and its sealed credentials, seed copy and recovery are verified |

Until the off-site copy runs, every backup lands in one Spaces bucket in `nyc3`, the cluster's own provider and region. The retained CSI snapshots in `platform/backups/snapshots-*.yaml` are dated, DigitalOcean-specific baselines, not a schedule.

## Investigate (read-only)

```sh
kubectl --context do-nyc3-omos-cloud get cronjobs,jobs -A
kubectl --context do-nyc3-omos-cloud get replicationsources -A
kubectl --context do-nyc3-omos-cloud get volumesnapshots -A
kubectl --context do-nyc3-omos-cloud -n velero get schedules.velero.io,backups.velero.io
kubectl --context do-nyc3-omos-cloud -n argocd get applications
```

The alerts in `platform/monitoring/backup-health.yaml` cover an overdue scheduled backup (`ScheduledBackupOverdue`, no success within 30 hours on a CronJob that is not suspended), VolSync sources out of sync or unscraped, and a failed off-site copy. Read only the relevant backup container's logs; do not dump Secrets, process environments or credential-bearing diffs.

## Restore rules

- Restore into new, explicitly named claims or disposable processes; never overwrite or rebind a source PVC. Keep the original volumes, retained snapshots and keys until the restored data and the user workflow are verified.
- Prefer the verified logical archives for databases and NATS; raw volume copies are crash-consistent. Restore clients must match the database major version.
- Recovery needs keys held outside the cluster: the logical-archive key, the Restic password and the Sealed Secrets keys ([Sealed Secrets](sealed-secrets.md)). Replacing the Spaces access key must preserve all of them.
- One-off restore exercises live in `platform/backups/restore-tests/` and `platform/backups/volsync/restore-tests/`. They are excluded from the recurring Kustomization and run only as a reviewed, explicit operation.
- Never delete PVCs, namespaces, archives or backup credentials as cleanup. `platform-backups`, `platform-volume-backups` and `platform-volsync` deliberately have no Argo CD finalizer: losing their sealed Restic credentials would make the repositories unreadable without the originals.

Restore evidence and remaining gaps are recorded in `docs/operations/`, for example `cluster-readiness-2026-10-08.md`, `hermes-restore-drill-2026-10-16.md` and `2026-10-26-tak-policy-preflight.md`. The outstanding Spaces credential rotation steps are in `platform/backups/README.md`.

## Verify completion

The scheduled Job completed and, for a logical archive, its `.verified.json` marker exists. The ReplicationSource reports a recent successful sync. After a real restore, the application's own workflow works against the restored data.

Generated workload lists: [Platform Backup and Recovery](../apps/platform-backup-recovery.md), [VolSync](../apps/volsync.md), [Velero](../apps/velero.md), [Cluster health observer](../apps/health-keeper.md).
