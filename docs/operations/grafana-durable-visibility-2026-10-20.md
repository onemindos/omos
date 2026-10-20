# Grafana durable state and Backstage visibility

Status: Stage A preparation; cutover/browser acceptance pending. User explicitly
approved a 5 GiB persistent volume (DO storage charges), backed-up migration,
brief Grafana maintenance and a dedicated read-only Backstage credential.
No Prometheus/NATS storage changes or wider AI permissions are authorized here.

## Observed cause and limits

Grafana 13.1.0 is healthy but uses `emptyDir` for `/var/lib/grafana`. Admin API
inspection found zero service accounts; the existing Backstage `glsa_` credential
returns 401 for `/api/user` and `/api/search`. API health is a separate 200.
This explains why replacing a token without persistence is not a durable repair;
we do not claim the exact historical deletion/restart cause was observed.

Installed native Grafana plugin 1.1.0 uses proxy prefix `/grafana/api` and appends
`/api/search`. The proxy strips its configured prefix, leaving `/api/search` at
the internal Grafana root. The apparent doubled `api` in a browser request is
intentional here; do not blindly remove a segment. Authentication, authenticated
proxy routing, dashboard selection and actual panel data are separate gates.

## Stage A — prepare before disturbing Grafana

- `platform-monitoring` owns retained `grafana-data`, 5 GiB,
  `do-block-storage-retain`, `Prune=false`.
- Argo-owned `grafana-state-seed-20261020` waits for an explicit operator seed.
  No Kubernetes API token, no host mounts, non-root UID/GID 472, pinned existing
  platform-backup image. It does not fetch arbitrary URLs or mutate workloads.
- The operator briefly quiesces the existing Grafana process, captures exactly
  SQLite DB/journal/WAL files and resumes it with a guaranteed cleanup path.
  No active database copy without its journal/WAL; no kill/recreate of the old
  pod before the source is preserved. Liveness budget must accommodate capture.
- Stream that snapshot directly to the seed pod (no plaintext Git or logs).
  Strict archive allowlist rejects traversal, links, duplicates and oversize.
  Recover a hot journal/WAL in isolated scratch storage; SQLite backup produces
  a canonical DB. Check integrity and every table's row count.
- Encrypt using the existing OMBKP002 authenticated AES-GCM backup format/key;
  upload to the existing Spaces destination under `backups/grafana/`, download,
  authenticate/decrypt and open the isolated restored DB. Require hash/table
  equality before writing the new PVC or its seed-success receipt.
- Existing backup key/cloud/pull credentials are resealed to named monitoring
  Secrets; no new cloud IAM grants, bucket or encryption-key replacement.
- A new `onemind-backstage-observer` service account has Viewer role and a 90-day
  token. Only a dedicated namespace-bound SealedSecret is committed. Create it
  before capturing state so the migrated database retains the account/token.
  The old invalid integration credential and unrelated integration keys stay
  untouched. Backstage does not consume the new credential until Stage B.
- Daily `backup-grafana-sqlite` starts **suspended**. It uses SQLite online backup
  on the existing PVC, authenticated remote download/decrypt/integrity/count
  verification, bounded resources/deadline and no extra cache/clone volumes.

## Stage B — separately reviewed cutover

After Stage A seed/remote restore succeeds, point the existing Helm-owned Grafana
Deployment at the existing seeded claim, use Recreate for ReadWriteOnce storage,
disable root chown init and require the seed marker/database before startup.
Keep chart/image/datasources/identity/admin Secret stable. Unsuspend the backup;
its pod must colocate with Grafana for RWO attachment. A GitOps-owned one-time
verification job must exercise the same remote-restore backup after cutover.

Switch only Backstage's Grafana token reference, preserving other integration
credentials and read-only proxy methods/authentication. Add exact verified native
dashboard selectors and namespace/workload-scoped links for the first services;
preserve them through inventory refresh. Do not enable anonymous Grafana, forward
human/agent tokens or relax iframe protections merely to make a dashboard appear.
A native dashboard link and an authenticated embedded panel are different gates.

No new image is required for mounted catalog/config/credential changes. A frontend
composition change, if needed, requires the existing release pipeline and pinned
image promotion; never claim it is deployed from source alone.

## Recovery and operation

Never revert to emptyDir or delete/overwrite the source PVC as a rollback. If
cutover fails, retain the PVC and seed receipt, inspect readiness/permissions,
and use a separately reviewed restore onto a new target if needed. The verified
remote archive and existing Git/secret configuration remain recovery inputs.
Data-source secrets are encrypted inside Grafana's DB; preserve its configured
secret-key/admin/SSO configuration alongside the existing sealing/recovery keys.

The backup format is compatible with `platform/backups/tools/decrypt.py`; decrypt
only to protected isolated storage with the existing independently held backup key.
Restore `data/grafana.db` as UID/GID 472 into a new PVC, verify integrity/counts
before starting the same Grafana version with its original encryption config.
An archive is not independent offsite recovery: Spaces is still the same provider
and region. Initial migration seed plus isolated remote restore is not a full
hardware/WAN-loss drill or proof that every plugin/dashboard query works.

Rotate the Viewer token before its 90-day expiry: create a replacement for the
same account/role, seal it, switch through reviewed GitOps, prove both replicas
and native search, then explicitly revoke the old token. Do not expose values in
the browser, commits or receipts. Add expiry monitoring in a bounded follow-up;
this change does not silently install a credential-rotation service.

## Acceptance, all pending until actually observed

- [ ] Prepared PVC bound and retained; seed job uses no API token or broad RBAC.
- [ ] Original snapshot integrity/counts verified; encrypted remote restore passed.
- [ ] Helm-owned Grafana uses the seeded PVC, not emptyDir; healthy on same version.
- [ ] Migrated Viewer account survives cutover and authenticates successfully.
- [ ] Scheduled backup enabled; same backup code passes a post-cutover run.
- [ ] Two Backstage replicas use the dedicated credential; anonymous/machine
      denial and read-only methods remain intact.
- [ ] Signed-in native dashboard discovery resolves exact first-service dashboards.
- [ ] Scoped Grafana dashboard shows actual CPU/memory series for the intended
      namespace/workload, not just a 200 or a dashboard name.
- [ ] Grafana-managed alert rules versus Prometheus alerts are labelled honestly;
      an empty Grafana rule list is not proof of no cluster alerts.
- [ ] GitOps ownership, database recovery inputs and limits recorded with PRs.
