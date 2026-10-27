# Private Nextcloud library sync

Observed 2026-10-26: Nextcloud is deployed under Argo CD, reports installed=true,
maintenance=false and no pending database upgrade, and its application and Redis
pods are healthy. The Mac client was running but unable to find its renamed local
folder. Its existing mapping and journal have been reconnected to the original
directory; real uploads resumed and a new Mac test edit matches the server by
SHA-256. The full-library transfer and final Android acceptance remain incomplete.
This user-directed Layer 2 service does not change the portal migration gates.

## Deployment and cost

- Argo Application: `platform-nextcloud`; chart 9.4.1; server 35.0.2 Apache.
- Private URL: `https://onemind-files.tail717ba0.ts.net/`.
- One application replica, Redis locking/cache, cron sidecar, retained file PVC.
  The live claim was 200 GiB at observation; this change requests 300 GiB for the
  larger audited library and temporary chunk assembly. Verify actual expansion
  after the reviewed merge and Argo/CSI reconciliation.
- Separate `nextcloud` database and role on the existing CNPG `onemind-data` cluster.
- Separate bootstrap administrator and non-admin `zeus` sync account. Local
  bootstrap and per-device credentials are outside the vault under
  `~/.config/onemind/nextcloud/`, with private directory/file permissions. Never
  copy them into the library, repository, logs or PR. Each device uses its own
  revocable app credential.
- Only private Tailscale ingress is configured. Authentik SSO remains deferred.
- DigitalOcean block storage is $0.10/GiB/month. The observed 200 GiB file volume
  plus 1 GiB backup cache costs $20.10/month; the requested 300 GiB plus cache
  costs $30.10/month, an additional $10/month, before tax, backups and any bandwidth
  overage. Existing compute/database costs are separate. The extra capacity needs
  human cost approval before merge/deployment. Zeus approved the original 200 GiB
  always-on hub on 2026-10-23; that approval is not recorded as approval of 300 GiB.
- The 2026-10-26 inventory using the installed ignore patterns is approximately
  159 GiB / 181,291 files, including a 36.67 GiB archive. This supersedes the earlier
  112.5 GB / 82,205-file estimate and remains an approximate source-side inventory.
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
Mac /Users/zeus/OneMind - 10-23-26
  <-> Nextcloud desktop client (classic sync)
  <-> Nextcloud /OneMind in Kubernetes
  <-> FolderSync Android (two-way)
  <-> /storage/emulated/0/Documents/onemind/OneMind
  <-> Obsidian Android
```

The Mac keeps the existing fully local folder. File Provider mode is disabled;
launch on login is enabled. No vault move or download-on-demand placeholders are
required. The desktop account is `zeus`, and the remote folder is `/OneMind`.
The dated local name does not rename the remote folder or the Samsung vault.

The Samsung S24 Ultra has a FolderSync account named `Nextcloud` and folder pair
`OneMind Full Sync`. It uses the private HTTPS server, port 443, and a dedicated
app password; certificate validation remains enabled. The local folder is the
nested `Documents/onemind/OneMind` directory opened by Obsidian, not its parent.
The pair is two-way with a five-minute schedule and deletion propagation enabled
for previously synchronized files. Background restrictions and actual scheduled
behavior must be checked separately from a foreground manual sync.

The official Nextcloud Android app browses the server directly. Its internal
two-way sync also exists: an isolated newer-build pilot on 2026-10-23 downloaded
a test file and uploaded an external edit and a newly created nested note. This
does not certify the production library, offline recovery, renames, deletions or
conflicts. FolderSync remains the production owner of the local Obsidian folder.
Keep the native pilot's `Native-Sync-Pilot-20261023` server/app-storage folder
separate. Do not enable another sync engine against the FolderSync folder.
The earlier `OneMind-Sync-Test` Obsidian plugin pilot is also separate.

### Path repair and device checks

The old `/Users/zeus/OneMind` directory no longer existed on 2026-10-26. The client
retried with `Unable to find local sync folder` and SQLite `CANTOPEN`/errno 2;
a running process and successful account HTTP requests did not establish syncing.
The renamed directory contained the original `.sync_ad7169db38e4.db` journal and
all 376 reference files retained their inode identities. With the client stopped,
its configuration and a consistent SQLite journal backup were saved privately,
then only the configured local path was changed. No vault move, empty replacement
folder, journal reset, Git reset or account recreation was performed.

For future folder moves, stop the client, identify the existing library and
journal, back up its configuration, update the local mapping, restart, and verify
actual file transfers. Never recreate an empty missing root to silence the error.
Preserve the vault's staged/uncommitted work. The Mac Obsidian registration is
separate from the desktop sync configuration; verify its path separately before
using Obsidian on the Mac.

On the physical Samsung S24 Ultra, `/sdcard` resolves to `/storage/emulated/0`.
The production vault appears in My Files as Internal storage → Documents →
onemind → OneMind. Map remote `/OneMind` to that exact folder, not its parent or
an extra `OneMind` child. It is separate from emulator/app-private storage.
ADB had no connected phone on 2026-10-26, so no fresh Samsung UI/offline acceptance
is claimed by the Mac repair. The prior human-entered `Samsung test` note is
present on the server; its download to the repaired Mac is still pending.

Android filename repairs on 2026-10-23 and 2026-10-26 preserve file content. The
latter repaired 18 remaining names in an archived sub-vault after checking 67,259
text files for old-name references (none found); 13 affected files retained inode,
size and modification time. Reversible mappings and audits are stored privately
outside the vault. Confirm the server and phone converge after the queued scan.

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

The library includes files up to 36.67 GiB. Android WebDAV clients may upload a
whole file in one PUT, while the Nextcloud desktop client uploads chunks. Values
set `APACHE_BODY_LIMIT=0`, PHP upload/post limits to 4 GiB and input/execution
limits to 3600 seconds. Verify these values after Argo rolls out the change, then
verify a real large-file upload/download and hash. Small-note success alone does
not prove large-file support.

The deployed local/Home storage does not implement `IChunkedFileWrite`, so the
optimized chunking-v2 path does not apply. `FutureFile`/`AssemblyStream` keeps the
uploaded chunks while a separate destination file is assembled; chunks are deleted
after the copy completes. The largest archive can therefore use approximately
73.34 GiB during assembly. In a worst-order initial transfer, 159 GiB of final data
plus 36.67 GiB of temporary chunks and about 7.6 GiB of other observed volume data
requires roughly 203 GiB before growth/retries. The observed 200 GiB PVC exposed
only 196.68 GiB of filesystem capacity. A small-note upload cannot validate this
headroom requirement. The 300 GiB request provides operating margin; it does not
prove large-file or recovery acceptance.

The storage class permits expansion and retains the volume. After approval and
Argo reconciliation, verify the PVC's reported capacity, resize conditions and
in-pod filesystem capacity before marking expansion complete. DigitalOcean volumes
cannot shrink in place: reverting the requested size is not a storage rollback;
reduction requires a separately reviewed migration to another volume.

## Acceptance and recovery

- [x] Human-authorized merge and Argo reconciliation of #424 and #425.
- [x] Private HTTPS, installed status, healthy pods and 200 GiB bound file PVC.
- [x] Non-admin Mac login and Samsung FolderSync credential test.
- [x] Existing Mac folder mapped to remote `/OneMind`; initial upload observed.
- [x] Samsung local Obsidian vault and two-way folder pair configured.
- [x] Original Samsung-origin test note reached Mac through the native clients.
- [x] Mac renamed-root failure diagnosed and existing journal preserved/reconnected.
- [x] New Mac repair edit reaches Nextcloud and matches by SHA-256 (2026-10-26).
- [x] Human-entered Samsung edit reached Nextcloud; background schedule observed.
- [ ] Human-entered Samsung edit downloaded to the repaired Mac.
- [ ] Mac edit returns to Samsung through the native clients.
- [ ] Offline edit, rename, small attachment and conflict preservation checks.
- [ ] Actual large-file upload/download and matching hashes verified.
- [ ] 300 GiB cost approval, reviewed merge and actual PVC/filesystem expansion.
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
