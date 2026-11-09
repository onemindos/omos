# Private Forgejo secondary — launch-ready 2026-11-04

User-selected route: **GitHub remains primary → Forgejo pull mirrors → independent
versioned encrypted backups**. This directory is deployable source for a small
private pilot, not a claim of deployment, owner enrollment, populated mirrors or
proven disaster recovery. `platform-forgejo` is manual-sync, with no finalizer.
Sync only after the capacity and backup-window gates below are confirmed.

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

**Launch state (2026-11-04):** both SealedSecrets are committed in `sealed/`,
sealed namespace-wide against the live controller certificate with
`scripts/prepare_forgejo_secrets.py` (generated from random values, never printed):

| Secret | Keys |
| --- | --- |
| `forgejo-keys` | `FORGEJO__security__SECRET_KEY`, `FORGEJO__security__INTERNAL_TOKEN`, `FORGEJO__oauth2__JWT_SECRET`, `FORGEJO__server__LFS_JWT_SECRET` |
| `forgejo-data-restic` | `RESTIC_REPOSITORY` (dedicated `.../backups/restic/forgejo/data` prefix), `RESTIC_PASSWORD` (fresh, unique), `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_DEFAULT_REGION` |

The Spaces access pair and region are the established platform backup storage
credentials (same pattern as `platform/vaultwarden/BACKUP.md`); only the prefix and
Restic password are new. The Restic password and application keys exist only as
SealedSecret ciphertext and, after sync, as the live Secrets. Recovery depends on the
sealed-secrets controller key backup (`backup-cluster-recovery`). **After the first sync the
owner should copy `RESTIC_PASSWORD` from the live Secret into the approved private password
store** (see below) so an offline copy exists. Never regenerate these: the application
`SECRET_KEY` decrypts stored mirror credentials, and the Restic password opens the backups.

Remaining gates before the first sync: confirm capacity (hub nodes had roughly 4 GiB
unrequested memory on 2026-11-04; the pod prefers the `onemind-hub-core-v2` pool),
the hostname, and a backup maintenance window. To regenerate for a different
environment, supply a private mode-0600 JSON input and the cluster public cert:

```bash
python3 scripts/prepare_forgejo_secrets.py --cert "$SEAL_CERT" --backup-json "$PRIVATE_BACKUP_JSON"
```

The helper refuses to replace existing sealed keys, writes no plaintext and performs
no cluster operation.

### Activation (coordinator, after merge)

```bash
argocd app sync platform-forgejo            # manual sync; no automated policy
kubectl --context do-nyc3-omos-cloud -n forgejo get sealedsecret,secret,pvc,pods,ingress
kubectl --context do-nyc3-omos-cloud -n forgejo rollout status deploy/forgejo
kubectl --context do-nyc3-omos-cloud -n forgejo get replicationsource forgejo-data
```

Verify PVC `Bound` on `do-block-storage-retain`, both SealedSecrets `Synced`, the
Deployment ready, `https://forgejo.tail717ba0.ts.net/` answering over the tailnet and
registration rejected. Health checks prove HTTP/DB startup, not mirror content. Do not
recreate claims or reset SQLite to fix readiness.

### Owner first sign-in

No account exists after sync (`INSTALL_LOCK=true` disables the web installer). The
owner creates the first administrator in a private operator terminal, so the generated
password appears only there, never in chat, PRs or logs:

```bash
kubectl --context do-nyc3-omos-cloud -n forgejo exec -it deploy/forgejo -- \
  forgejo admin user create --admin --username <owner> --email me@zeusdelacruz.com \
  --random-password --must-change-password
```

Check `forgejo admin user create --help` first, as flags can differ by version. If
`--random-password` is unavailable, generate a password in the terminal and pass it
through a shell variable, never inline in shared output. Store the printed password
directly in the approved password manager, open `https://forgejo.tail717ba0.ts.net/`
on a tailnet device, sign in and change it at the forced prompt, then enable 2FA.
No agent creates or holds this account.

### GitHub read-only token for mirrors

Create a GitHub **fine-grained personal access token** (owner: the organization or
account, repository access limited to the approved repositories, permissions
**Contents: Read-only** and **Metadata: Read-only** only; add **Issues/Pull requests: read**
only if metadata migration is wanted). Authorize it for organization SSO where
required. Enter it only in the Forgejo "New Migration → GitHub → Mirror" form
(private repository, interval 8h) for **one** harmless repository first. Forgejo stores it
encrypted with `SECRET_KEY`. Do not paste it into Git, chat or Kubernetes manifests.

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

## Authentik sign-in (added 2026-11-08)

Forgejo has an extra "Sign in with Authentik" button. The local admin login stays as the break-glass login, and registration stays closed.

- `auth-provider.yaml` (Argo Sync hook, namespace `auth`) creates the Authentik provider/application `onemind-forgejo`, bound to the owner account only. Callback: `https://forgejo.tail717ba0.ts.net/user/oauth2/authentik/callback`.
- The client secret is sealed twice: `auth/forgejo-oidc-bootstrap` (for the hook) and `forgejo/forgejo-oidc` (read by the pod as `FORGEJO_OIDC_CLIENT_SECRET`, only so the one-time command below can use it).
- Forgejo keeps login sources in its database, not in `app.ini`, so the source is created once with the CLI (the same command with `update-oauth --id N` changes it):

```sh
kubectl --context do-nyc3-omos-cloud -n forgejo exec deploy/forgejo -- sh -c \
  'forgejo admin auth add-oauth --name authentik --provider openidConnect \
   --key onemind-forgejo --secret "$FORGEJO_OIDC_CLIENT_SECRET" \
   --auto-discover-url https://auth.onemindos.dev/application/o/onemind-forgejo/.well-known/openid-configuration \
   --scopes "openid email profile"'
```

- First sign-in: click the Authentik button, sign in to Authentik, then choose "Link to existing account" and enter the local admin username and password once. Auto-registration is off and linking asks for the local password (`ACCOUNT_LINKING=login`).
- Rollback: delete the source in Site Administration -> Authentication Sources, or `forgejo admin auth delete --id N`. Local login is unaffected.
