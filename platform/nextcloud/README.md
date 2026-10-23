# Private Nextcloud vault sync

Status: prepared for review on 2026-10-23; not deployed or acceptance-tested.
This is a user-requested Layer 2 shared file service. Argo owns it; a human reviews
and merges the PR. It does not change the OpenChoreo migration phase gates.

## Deployment

- Application: `platform-nextcloud`, Nextcloud chart 9.4.1, server 35.0.2 Apache.
- Private planned URL: `https://onemind-files.tail717ba0.ts.net`.
- One application replica, Redis cache/locking, cron sidecar and a 20 GiB retained
  filesystem PVC. The initial size is for a selected vault, not the 112 GB archive.
- PostgreSQL: separate `nextcloud` role/database in existing CNPG `onemind-data`.
  Role credentials are sealed separately in the database and application namespaces.
- Separate bootstrap administrator and non-admin `zeus` sync user. First-install
  hooks create the sync user. Local bootstrap credentials are saved outside the
  vault under `~/.config/onemind/nextcloud/bootstrap.json` with mode 0600. Never
  copy that file into the vault, logs or PR. Use browser login to issue per-device
  app credentials; do not use the administrator account for synchronization.
- No public ingress, public file sharing route, Office, Talk or AI service is added.
  Authentik SSO is deferred; the first sync uses Nextcloud's own accounts.
- Apache starts as root to initialize its volume and switch to www-data. Baseline
  Pod Security and a limited capability set support that upstream startup behavior.
  Redis runs non-root. Workload service-account token mounting is disabled.

After merge, let the existing app-of-apps and ecosystem-resources Applications
reconcile. Do not run `helm install` or manually apply the workloads. Verify CNPG
role/database reconciliation, Secret unsealing, PVC binding, the application and
Redis rollout, and the Tailscale ingress certificate before connecting any vault.
Verify `/status.php` reports installed=true, maintenance=false and no DB upgrade
pending. Verify Login Flow v2 returns HTTPS URLs on the planned hostname.

## First two-device test

1. Keep Obsidian and Tailscale up to date on Android. Plugin 1.0.12 requires
   Obsidian 1.13 or newer. Mac Obsidian was observed at 1.14.4.
2. Open the prepared local `/Users/zeus/OneMind-Sync-Test` folder as a vault on Mac.
   Enable its Nextcloud Sync community plugin only after the endpoint is verified.
3. In plugin settings, use this full WebDAV endpoint:
   `https://onemind-files.tail717ba0.ts.net/remote.php/dav/files/zeus/`.
   Sign in using browser login as `zeus`, then run Sync now.
4. On Android create an **empty local vault** called exactly `OneMind-Sync-Test`.
   Install and enable Nextcloud Sync. Use the identical WebDAV endpoint and account.
   The Nextcloud Android files app is not required for this plugin workflow.
5. Run Sync now on Android and open `sync-check.md`. Add a line on Android, sync,
   then sync on Mac and check it arrived. Add a Mac line and check it on Android.
6. Test offline edits and reconnection, one rename, a small attachment and conflict
   recovery. Keep both apps open during the test and inspect errors/conflicts.
   Do not use Mirror from remote: it discards unsynced local-only files.
7. Turn on Sync on file change on Android. Mobile sync runs while Obsidian is
   open, plus startup/resume/manual sync; it is not an always-running background
   service. Confirm whether Wi-Fi-only sync is desired. Large mobile attachments
   are capped by default; choose explicit attachment rules before the real vault.

The main `/Users/zeus/OneMind` vault has not been enabled or uploaded. It contains
large binary archives, sensitive-file exclusions and substantial uncommitted Git
reorganization. Select and review the first sync set before connecting it. Git
ignore rules do not configure Nextcloud exclusions. Exclude `.git` at every depth,
credential/runtime state, archives and plugin secrets. Preserve existing attachment
links when moving files. Use one sync engine on each device; no second file-sync
engine or Android Git sync against the same live vault. Mac Git can record mobile
changes after they arrive, once its existing staged/deleted state is reconciled.

## Backup and recovery

The namespace has a dedicated VolSync repository for the complete application
PVC, including config, server encryption secrets and files. Daily snapshots are
scheduled at 07:10 UTC; CNPG already has base backups and WAL archiving. Backup
credentials are namespace-bound SealedSecrets. The repository is separate from
existing services; existing backup credentials and encryption password were reused
without rotation. Preserve existing recovery keys independently of the cluster.

Independent database and filesystem backups do not prove a consistent application
restore. Before treating this as the primary store, quiesce writes, capture the
coordinated recovery point, restore database plus config/files into an isolated
instance, and verify hashes, permissions, history and client behavior. Do not
claim a backup passed until a completed backup and restore have been observed.
An image rollback alone cannot reverse a Nextcloud database migration. Upgrade
one server major at a time, with a coordinated backup and restore plan.

## Acceptance record (pending)

- [ ] Human merge and Argo reconciliation.
- [ ] Non-admin browser login, HTTPS, private reachability and WebDAV write/read.
- [ ] Mac to server to Android and Android to server to Mac, verified by the user.
- [ ] Offline edit, rename, attachment and conflict checks.
- [ ] Scheduled cron, VolSync completion and coordinated restore drill.
- [ ] Reviewed main-vault exclusions and first-upload manifest.

References: [chart](https://github.com/nextcloud/helm),
[plugin](https://github.com/siosig/obsidian-nextcloudsync),
[Nextcloud backup](https://docs.nextcloud.com/server/stable/admin_manual/maintenance/backup.html).
