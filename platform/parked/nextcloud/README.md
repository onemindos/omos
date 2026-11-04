# Private Nextcloud library sync

Observed 2026-10-23: Nextcloud is deployed under Argo CD, reports installed=true,
maintenance=false and no pending database upgrade, and its application and Redis
pods are healthy. The Mac's initial full-library upload has started. Android is
configured with a native two-way folder pair; end-to-end acceptance is in progress.
This user-directed Layer 2 service does not change the portal migration gates.

## Deployment and cost

- Argo Application: `platform-nextcloud`; chart 9.4.1; server 35.0.2 Apache.
- Private URL: `https://onemind-files.tail717ba0.ts.net/`.
- One application replica, Redis locking/cache, cron sidecar, retained 200 GiB PVC.
- Separate `nextcloud` database and role on the existing CNPG `onemind-data` cluster.
- Separate bootstrap administrator and non-admin `zeus` sync account. Local
  bootstrap and per-device credentials are outside the vault under
  `~/.config/onemind/nextcloud/`, with private directory/file permissions. Never
  copy them into the library, repository, logs or PR. Each device uses its own
  revocable app credential.
- Only private Tailscale ingress is configured. Authentik SSO remains deferred.
- DigitalOcean block storage adds $20/month for the 200 GiB file volume and
  $0.10/month for the 1 GiB backup cache, before tax, backups and any bandwidth
  overage. Existing cluster compute/database costs are separate. Zeus chose to
  retain the always-on hub after discussing this cost on 2026-10-23.
- The full source library inventory is approximately 112.5 GB / 82,205 files.
  Allocated capacity is not the amount already transferred. Initial transfer must
  be allowed to finish before every file is available from the server or phone.

Argo owns all workload changes. Do not run `helm install` or manually apply them.
PRs [#424](https://github.com/onemindos/omos/pull/424) and
[#425](https://github.com/onemindos/omos/pull/425) deployed the service, corrected
Apache bootstrap capabilities and expanded the file PVC. Apache needs its limited
root bootstrap capabilities to initialize ownership, then serves as www-data.
Redis runs non-root; workload service-account token mounting is disabled.

## Device connections

The production path is:

```text
Mac /Users/zeus/OneMind
  <-> Nextcloud desktop client (classic sync)
  <-> Nextcloud /OneMind in Kubernetes
  <-> FolderSync Android (two-way)
  <-> /storage/emulated/0/Documents/onemind/OneMind
  <-> Obsidian Android
```

The Mac keeps the existing fully local folder. File Provider mode is disabled;
launch on login is enabled. No vault move or download-on-demand placeholders are
required. The desktop account is `zeus`, and the remote folder is `/OneMind`.

The Samsung S24 Ultra has a FolderSync account named `Nextcloud` and folder pair
`OneMind Full Sync`. It uses the private HTTPS server, port 443, and a dedicated
app password; certificate validation remains enabled. The local folder is the
nested `Documents/onemind/OneMind` directory opened by Obsidian, not its parent.
The pair is two-way with a five-minute schedule and deletion propagation enabled
for previously synchronized files. Background restrictions and actual scheduled
behavior must be checked separately from a foreground manual sync.

The official Nextcloud Android app can browse files, but does not yet implement
complete bidirectional arbitrary-folder synchronization. FolderSync provides that
local-folder behavior for Obsidian. Do not enable the Obsidian Nextcloud Sync
plugin against a folder already managed by these native clients. The earlier
`OneMind-Sync-Test` plugin pilot is not the production sync engine.

## Library scope and Git

PDFs, images, audio, video, installers and archives are included. Git remains on
the Mac; its existing staged/uncommitted reorganization must not be reset or
committed as part of sync setup. The Nextcloud journal pattern `.sync_*.db*` was
added to the vault's existing `.gitignore`; Git ignore rules do not configure
file-sync exclusions.

The Mac's native `sync-exclude.lst` excludes Git internals, dependency/cache
folders, `.env` files, credential/key containers, plugin `data.json`, workspace
state and the Git/Nextcloud sync plugins. Android has an exclude File path RegEx
for the corresponding runtime/credential files and the whole `.obsidian` folder,
so its mobile configuration stays device-local. Exact local policy and Android
regex are recorded privately under `~/.config/onemind/nextcloud/`.

Private notes and documents are part of the explicitly requested private library.
These exclusions are not a redaction of credentials embedded inside note content.
Keep native Nextcloud TLS and account access controls in place. Only one file-sync
engine should own the live vault on each device; Android Git automation is not
configured against this folder.

## Large files

The library includes files larger than 1 GiB. Android WebDAV clients may upload a
whole file in one PUT, while the Nextcloud desktop client uploads chunks. Values
set `APACHE_BODY_LIMIT=0`, PHP upload/post limits to 4 GiB and input/execution
limits to 3600 seconds. Verify these values after Argo rolls out the change, then
verify a real large-file upload/download and hash. Small-note success alone does
not prove large-file support.

## Acceptance and recovery

- [x] Human-authorized merge and Argo reconciliation of #424 and #425.
- [x] Private HTTPS, installed status, healthy pods and 200 GiB bound file PVC.
- [x] Non-admin Mac login and Samsung FolderSync credential test.
- [x] Existing Mac folder mapped to remote `/OneMind`; initial upload observed.
- [x] Samsung local Obsidian vault and two-way folder pair configured.
- [ ] Samsung-origin file reaches Mac through the native clients.
- [ ] Mac edit returns to Samsung through the native clients.
- [ ] Offline edit, rename, small attachment and conflict preservation checks.
- [ ] Actual background schedule and large-file transfer verified.
- [ ] Full initial transfer completed and source/destination inventories checked.
- [ ] VolSync completion and coordinated restore drill.

VolSync has a dedicated repository for the complete application PVC, including
config and server secrets, scheduled at 07:10 UTC. CNPG has existing base backups
and WAL archiving. Backup credentials are namespace-bound SealedSecrets; existing
backup credentials were reused without rotation. Independent filesystem/database
backups do not prove a consistent restore. Before treating this as the primary
store, quiesce writes, capture a coordinated point, restore into an isolated
instance, and verify file hashes, permissions, history and clients. Preserve
recovery keys independently of this cluster. Do not claim restore acceptance until
it has actually passed. Image rollback cannot reverse a database migration;
upgrade one server major at a time with a coordinated recovery plan.

References: [chart](https://github.com/nextcloud/helm),
[Android sync limitations](https://docs.nextcloud.com/server/stable/user_manual/en/files/desktop_mobile_sync.html),
[FolderSync](https://foldersync.io/docs/help/concept/),
[large files](https://docs.nextcloud.com/server/stable/admin_manual/configuration_files/big_file_upload_configuration.html),
[backup](https://docs.nextcloud.com/server/stable/admin_manual/maintenance/backup.html),
[DigitalOcean volume pricing](https://docs.digitalocean.com/products/volumes/details/pricing/).
