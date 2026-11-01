# Private Forgejo secondary — prepared 2026-11-01

User-selected route: **GitHub remains primary → Forgejo pull mirrors → independent
versioned encrypted backups**. This directory is deployable source for a small
private pilot, not a claim of deployment, owner enrollment, populated mirrors or
proven disaster recovery. `platform-forgejo` is manual-sync, with no finalizer.
Do not merge/sync until the credential and capacity gates below are complete.

## Layout and scope

- Rootless **Forgejo15.0.11LTS** (supported through15July2027), registry-verified
  image index `sha256:3cd6bc2b74c87cbde0a0fd527c9665cb416f5a0540dee54fe5afb3fa2df9f40c`.
- Single writer, Recreate; **SQLite** avoids another database for a mirror/standby
  pilot. Official documentation supports SQLite and whole-data-directory backup.
  Reassess PostgreSQL and storage sizing before large/team-primary workloads.
- One **20Gi** `forgejo-data` PVC on `do-block-storage-retain`. Git repositories,
  SQLite including journals, LFS, attachments, queues and configuration stay under
  `/var/lib/gitea`. The pinned rootless entrypoint uses
  `/var/lib/gitea/custom/conf/app.ini`, explicitly set in Git. Its `SECRET_KEY`
  decrypts stored mirror credentials and must survive restoration unchanged.
- Proposed private HTTPS **https://forgejo.tail717ba0.ts.net/** through Tailscale.
  The hostname is not registered/verified by this source. Native local accounts
  remain the application login; tailnet reachability is not account authorization.
- No public registration, public repositories, SSH listener, Actions runners,
  package registry, mail or push-create repositories. HTTPS Git/LFS is supported.
  No GitHub origin, Argo repo URL, Actions runner or GHCR reference changes.
- Requests: app100mCPU/256Mi, limits2CPU/1Gi; backup coordinator25m/48Mi, plus
  transient Restic mover100m/256Mi and1Gi cache. These are initial budgets, not
  measured production consumption. PVC expansion needs repository/LFS inventory.

The app has no Kubernetes API token. Default-deny networking permits only its
named Tailscale ingress proxy, DNS and public HTTPS downloads. Private/metadata/
tailnet target IPs are excluded. Migration URLs are restricted to GitHub domains;
LFS download redirects still need representative testing. No mirror access key is
configured by source. The backup coordinator is a different identity with no
Secret reads, data mount, exec, workload-create or namespace-delete permission.

## Credentials and activation

1. Review the exact-main diff, cluster capacity, Tailscale hostname availability,
   storage/billing and a backup maintenance window. Supply a cluster public sealing
   certificate and a **private mode0600 JSON** recovery input containing keys
   `RESTIC_REPOSITORY`, `RESTIC_PASSWORD`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.
   Repository must be a dedicated HTTPS S3 path ending
   `/backups/restic/forgejo/data`; use a new strong recovery password and preserve
   it outside Kubernetes, including an offline copy. Supply existing approved
   object-store credentials through a private file; do not paste them into Git,
   command arguments, logs or chat.
2. Run the helper with private file paths:

   ```bash
   python3 scripts/prepare_forgejo_secrets.py --cert "$SEAL_CERT" --backup-json "$PRIVATE_BACKUP_JSON"
   ```

   It creates only two strict-bound SealedSecrets and updates the
   sealed Kustomization; it refuses to replace existing keys, writes no plaintext
   and performs no cluster operation. Review the ciphertext, commit it and use Argo.
   Missing SealedSecrets are an explicit launch blocker, not dummy credentials.
3. Merge/sync the app through Argo. Verify PVC protection, SealedSecret status,
   deployment readiness, HTTPS and rejected registration. Health checks prove
   HTTP/DB startup, not useful mirror content. Run migrations only through the
   pinned normal startup; do not recreate claims or reset SQLite to fix readiness.
4. Create the first local administrator using the running image's native admin
   CLI in a private operator terminal. `INSTALL_LOCK=true` intentionally disables
   the web installer even before the first user. Use a generated password and
   `--must-change-password`; capture it only in the approved private password store.
   Inspect CLI help/version for exact arguments. Never print credentials into the
   PR or use a public registration/bootstrap window. OIDC can be a later change;
   do not mutate existing Authentik enrollment to accelerate this pilot.
5. Sign in, confirm access boundaries and prepare **one** harmless private GitHub
   repository as the first mirror. After backup and restore acceptance below,
   expand the approved list. No account or repository enrollment has been performed.

## Mirror enrollment and coverage

Create destination repositories through **New Migration → GitHub → Mirror**;
select private and a periodic pull interval (default8h; minimum1h). Pull-mirror
status is chosen during migration; do not create writable ordinary repositories
and expect them to become scheduled mirrors automatically. Use a separate scoped
GitHub fine-grained read token approved for the intended repositories, including
organization authorization/SSO when applicable. Preserve the application encryption
key so this stored token can be recovered. Do not enable Actions on mirror imports.

Inventory personal/org/private/archived repositories and submodule repositories;
record approved scope, discovery date, source URL, refs and destination. The pilot
has **no automatic new-repository discovery or failed-sync notification integration**.
Review mirror status and last-update timestamps manually until a separate automation
is implemented. A private fork/mirror is a second copy, not immutable recovery.

| Data | What must be verified independently |
| --- | --- |
| Git refs/history | Compare heads/tags and commits; fresh clone with GitHub unavailable; run Git integrity checks. Mirrors can reproduce upstream deletions/force pushes. |
| LFS | Enable LFS, fetch actual objects, prove fresh checkout for representative LFS commits. Pointer files alone are not the content. |
| Wikis/submodules | Capture and test each separate repository; a parent repo contains only submodule references. |
| PRs/issues/reviews/discussions | Native pull mirroring is Git synchronization, not continuous metadata backup. Supported migration/API exports and attachment capture are separate work. |
| Releases | Tags are Git data; release notes/assets need separate export/download and checksums. |
| Actions/secrets | Workflow YAML is Git data; logs, artifacts, rules and secret values are not. Keep original secret sources in encrypted recovery storage. |
| GHCR/packages | Copy required images/manifests by digest and package payloads separately. This PR does not turn Forgejo into a backup OCI registry. |
| Local changes | Uncommitted/ignored code requires a separate workstation file backup. |

## Consistent independent backups

The official Forgejo guide strongly recommends **shutting down the writer** to
coordinate SQL, Git and file state. This implementation follows that advice:

- CronJob runs daily **15:20UTC**, concurrency Forbid, no automatic job retry.
- It reads the exact `forgejo` Deployment scale and refuses counts outside0/1.
- If running, it scales only that writer to0 and waits for its pods to disappear.
- It triggers only `ReplicationSource/forgejo-data`, then waits for that exact
  unique manual tag in `status.lastManualSync`. No generic schedule triggers an
  unquiesced snapshot; `trigger: {}` is idle until the coordinator requests one.
- VolSync snapshots the whole PVC and Restic uploads the snapshot clone to the
  separate encrypted S3 repository, with7daily/4weekly/12monthly/1yearly retention.
  Clone/cache claims explicitly use delete-class `do-block-storage`; the original
  retained PVC is never deleted. Source UID1000 matches all app files.
- The coordinator restores the **prior** replica count in `finally`, including
  timeout/API failure/SIGTERM paths. It does not restart an intentionally stopped
  writer or overwrite a concurrently changed replica count. Kubernetes resource
  versions guard both scale changes. No Secret reads or image pulls occur in script. Do not launch overlapping manual
coordinator Jobs; CronJob Forbid only serializes its own scheduled children.

**This causes a nightly outage for the duration of backup export**, bounded by
3minutes for shutdown and30minutes for export, plus API/recovery overhead. The
first large upload can exceed that budget; do the initial backup in an attended
window and measure before adopting this schedule. Successful upload is not proof
of restoration. Network policy/Cilium kube-apiserver reachability, CSI snapshot
support, VolSync and the target credentials must be verified on the real cluster.

`activeDeadlineSeconds=2280` leaves margin beyond internal timeouts and a60second
termination grace period. SIGKILL, node loss, API loss or failed recovery PATCH can
still leave the app stopped. A failed coordinator Job requires inspection; never
claim a fresh backup from a failed/partial run. Recover the desired writer count
with the reviewed Argo desired replica1 **after** confirming no active backup owns
quiescence. Do not sync/restart or upgrade Forgejo during its backup window. Argo
manual sync is deliberate: automatic self-heal would fight the maintenance pause.

The Restic repository is outside the Forgejo PVC/cluster. Existing platform
second-provider backup-copy jobs must be verified to include this **new** prefix;
their mere presence is not a claim this repository is replicated. Keep a separate
offline copy of Restic/encryption/sealing keys and a portable restore procedure.
Provider Object Lock/WORM is **not** configured here, so dated retention is not an
immutability guarantee. No backup or restore has been executed by this PR.

## Restore and upgrade acceptance

Before normal mirror enrollment, trigger the attended coordinator through the
reviewed job workflow and verify that its manual tag completed. Restore a selected
Restic snapshot onto an isolated new retain-class20Gi claim, including config and
keys, owned UID/GID1000. Use the **same pinned Forgejo version** and no existing
production ingress. Inspect SQLite integrity, owner access, repository refs, a
fresh Git clone, LFS payload, wiki and any uploaded asset. Verify recovered mirror
credentials privately; do not contact/update upstream during the isolated drill.
Store the restore result/time and checksums, not secrets. Temporary drill claims
use `do-block-storage` if they will be destroyed; never delete/recreate the live claim.

Read release/migration notes and take an attended recoverable backup before every
upgrade. Reverting a container does not reverse database migrations. Preserve
app.ini/SECRET_KEY/INTERNAL_TOKEN/JWT keys; do not rotate them to fix login problems.
GitHub remains the working forge even if this standby is unavailable.

## Sources and verification

- [LTS release and support](https://forgejo.org/releases/)
- [Official rootless image and SQLite option](https://forgejo.org/docs/v15.0/admin/installation/docker/)
- [Backup scope and consistency](https://forgejo.org/docs/v15.0/admin/setup/backup/)
- [Pull mirrors](https://forgejo.org/docs/v15.0/user/repo-mirror/)
- [GitHub backup exclusions](https://docs.github.com/en/repositories/archiving-a-github-repository/backing-up-a-repository)

Local checks: `kubectl kustomize platform/forgejo`, changed YAML/schema validation,
MCP coverage/catalog checks, and `python3 -m unittest scripts.tests.test_gitops_forgejo`.
Registry digest verified2026-11-01. An ephemeral rootless container with exact
UID1000/read-only-root/dropALL/no-new-privileges and writable data/tmp passed
HTTP health200 and created SQLite/config under the retained path with no `/etc`
configuration fallback. No account, mirror or external inference was invoked.
The temporary container was removed; its data was synthetic, not production state. Real Argo rollout, account/mirror enrollment and recovery stay pending.
