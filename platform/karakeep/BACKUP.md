# Karakeep persistence and recovery gate

No Karakeep backup is claimed live or restore-tested from source alone. The user authorized first start and the existing VolSync/Restic pipeline now has two app-owned sources. The full `/data` volume is authoritative: SQLite, its journals/WAL, assets/attachments, background-job state and related app files belong together. Keep the strict SealedSecrets and the cluster sealing-key recovery separately. User bookmark JSON/ZIP exports produced by the app are useful portability exports, not a complete users/config/assets disaster-recovery backup.

## Deployed source policy

`karakeep-data` runs at 12:40 UTC daily and `karakeep-meilisearch` at 13:00 UTC daily. Both use the installed VolSync 0.16.0 operator, Snapshot copy method, existing DigitalOcean CSI `do-block-storage` snapshot/storage classes, one 1Gi cache each and non-root UID/GID1000 movers. Each mover requests 100m CPU / 256Mi memory with a 1 CPU / 1Gi limit. Namespace egress permits DNS and HTTPS to the existing S3 path, without giving the application access to backup credentials or the Kubernetes API.

Strict namespace-bound SealedSecrets `karakeep-data-restic` and `karakeep-meilisearch-restic` contain the existing authorized Spaces backup credential, fresh independent Restic encryption passwords and distinct repository prefixes:

- `s3:https://nyc3.digitaloceanspaces.com/onemind-backups/backups/restic/karakeep/data-pvc`
- `s3:https://nyc3.digitaloceanspaces.com/onemind-backups/backups/restic/karakeep/meilisearch-pvc`

The existing provider credential is reused without changing its permissions; these are separate repositories, not newly asserted provider-side prefix restrictions. Fresh encryption passwords and the app-only OIDC secret are preserved privately outside Git in the coordinator's task recovery material (mode0600). Preserve an independent recovery copy. The existing app/session and Meilisearch keys are preserved, not rotated. The encrypted cluster-recovery export must include these new secrets before cluster-loss recovery is considered covered.

Retention matches the existing pipeline: 7 daily / 4 weekly / 3 monthly / 1 yearly, with prune every 14 days. **Scheduled snapshots are crash-consistent** full-volume copies, including SQLite WAL/journals. There is no automatic app quiesce hook in the installed CRD, no backup account with workload-scaling rights and no application/DB coordinated-restore claim. Storage is in the same provider/region as the cluster; the independent offsite copy remains a separate platform prerequisite.

## Before real use

1. Record `karakeep/data-pvc` (20Gi) and `karakeep/meilisearch-pvc` (5Gi) UIDs/PVs after manual Argo creation; both use do-block-storage and explicit Prune=false, Delete=false, Replace=false, Force=false protections. Do not delete the namespace or claims on uninstall.
2. Confirm both dedicated SealedSecrets are Synced and both ReplicationSources use the expected repository secret. These AWS-named fields are the existing Restic S3-protocol contract, not AWS service credentials. See the [existing VolSync runbook](../backups/volsync/README.md). Do not print secret values or replace encryption passwords during redeploys.
3. Trigger the initial sources using new `spec.trigger.manual` identifiers through reviewed Git/Argo; retain the daily schedules. Wait for matching `status.lastManualSync`, recent successful `lastSyncTime` and no mover errors. Verify the encrypted repositories can be read with preserved keys. Snapshot allocation alone is not an off-cluster backup.
4. Restore the full data snapshot into an isolated new claim and run SQLite `PRAGMA integrity_check` on the restored database with its WAL/journals. Verify restored attachments and app behavior before accepting the crash-consistent policy for normal use. For an explicitly app-consistent/pre-upgrade copy, stop web through reviewed Git/Argo and wait until its sole writer has stopped before capturing/uploading the complete volume; restart via Git only after success. Never copy only live `db.db` and omit WAL/journals/assets.
5. Meilisearch is a derived index but its exact rebuild procedure must pass before calling it disposable. Its own whole-volume backup is retained initially. Pin the engine version and verify the native Karakeep reindex operation during restore acceptance; do not assume an image downgrade can read a newer index.
6. Restore into an isolated new namespace/claim with no production ingress, email, webhooks, crawler access or inference. Verify SQLite integrity, existing owner login, bookmark/attachment counts and checksums, search/reindex, app key recovery and export. Record snapshot ID/time, restore result and recovery scope privately; preserve the original claims. No such acceptance has happened in this source task.

Acceptance remains incomplete until initial encrypted backup and isolated restore evidence is recorded. Repeat isolated restores monthly and take a fresh app-consistent copy before upgrades. These schedules and retention do not establish an RPO/RTO guarantee. The app has single-node SQLite/RWO availability and short Recreate outages; HA would require a different supported architecture and is not silently introduced.
