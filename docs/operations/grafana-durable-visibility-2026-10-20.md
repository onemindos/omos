# Grafana durable state and Backstage visibility

Status: scoped production acceptance passed, 2026-10-20 22:29–22:35 UTC.
Preparation [#340](https://github.com/onemindos/omos/pull/340), cutover
[#341](https://github.com/onemindos/omos/pull/341), Argo ownership correction
[#342](https://github.com/onemindos/omos/pull/342) are merged and reconciled.
Stage A: [#340](https://github.com/onemindos/omos/pull/340), merged `9f9358f6`. User explicitly
approved a 5 GiB persistent volume (DO storage charges), backed-up migration,
brief Grafana maintenance and a dedicated read-only Backstage credential.
No Prometheus/NATS storage changes or wider AI permissions are authorized here.

## Observed cause and limits

Before this repair, Grafana 13.1.0 was healthy but used `emptyDir` for `/var/lib/grafana`. Admin API
inspection found zero service accounts; the existing Backstage `glsa_` credential
returned 401 for `/api/user` and `/api/search`. API health is a separate 200.
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
the browser, commits or receipts. A value-only Secret update does not refresh
existing process environment variables: use a new Secret name/reference or a
Git-owned rollout annotation in the reviewed change, not an imperative restart.
Add expiry monitoring in a bounded follow-up;
this change does not silently install a credential-rotation service.

## Acceptance observed after cutover

- [x] Prepared PVC bound and retained; seed job uses no API token or broad RBAC.
- [x] Original snapshot integrity/counts verified; encrypted remote restore passed.
- [x] Helm-owned Grafana uses the seeded PVC, not emptyDir; healthy on same version.
- [x] Migrated Viewer account survives cutover and authenticates successfully.
- [x] Scheduled backup enabled; same backup code passes a post-cutover run.
- [x] Two Backstage replicas use the dedicated credential; anonymous/machine
      denial and read-only methods remain intact.
- [x] Signed-in native dashboard discovery resolves exact first-service dashboards.
- [x] Scoped Grafana dashboard shows actual CPU/memory series for the intended
      namespace/workload, not just a 200 or a dashboard name.
- [x] Grafana-managed alert rules versus Prometheus alerts are labelled honestly;
      an empty Grafana rule list is not proof of no cluster alerts.
- [x] GitOps ownership, database recovery inputs and limits recorded with PRs.


## Stage A observed result and Stage B review

At 22:09 UTC, the original process was quiesced for 0.489 seconds and resumed.
2,582,016 bytes of DB/journal archive were streamed to the Argo-owned seed job.
An initial capture attempt failed closed because a watchdog inherited stdout and
held the stream open; no PVC data was seeded by that attempt, Grafana was resumed,
and stdout redirection corrected the operator helper before the successful retry.
The source database itself was not mutated by the capture.

Verified encrypted object:
`backups/grafana/20261020T220905Z-12b47eae2c094f3091abdbc64e5c58dd/archive.tar.gz.aesgcm`.
Remote download/decrypt/hash equality and isolated SQLite integrity/count checks
passed for all 87 tables before the PVC seed marker was written. Seed job succeeded.
At this Stage A observation, original Grafana remained available pending Stage B. The new Viewer
account was created before this snapshot, so its token state is included.

Full pinned kube-prometheus-stack 87.16.1 before/after Helm renders show **only**
`Deployment/monitoring-grafana` changed: retained existing claim, Recreate, required
seed guard and a persistent-pod label. No new Helm-owned PVC, no Prometheus/NATS
workload, data-source, chart version or other rendered resource changes. Backup
pods require the new persistent-pod label so they cannot race the old Grafana pod
onto another node during RWO cutover.

Three Git-owned dashboards are intrinsically scoped to Backstage, its PostgreSQL
workload and NATS. Their Prometheus datasource and CPU/memory/owner record queries
were tested against actual data before review. UID selectors and links are restored
by inventory refresh. No generic namespace-variable defaults or unsigned iframe
are presented as a service-specific view. Native alerts require an explicit label
selector, using the installed predicate schema; a test ensures every existing host
extension is retained by the runtime config's replacement array.

The older ClickHouse NATS dashboard/datasource remains separate: the original
Grafana had no installed ClickHouse datasource plugin. These new workload panels
use the existing working Prometheus path; no claim the old ClickHouse panels or
application-level/JetStream metrics have been repaired.


## Argo field-ownership correction during Stage B

At 22:21 UTC, Argo's server-side comparison rejected the union transition because
live Helm-owned `emptyDir` and default `rollingUpdate` fields were retained beside
the desired claim/Recreate. Nothing in Grafana was replaced by this failed diff;
the original pod remained healthy and the seed was preserved. Other Applications
can reconcile independently; never treat the partial portal rollout as cutover.

A narrow follow-up selects legacy diff for `monitoring-kube-prometheus-stack`
and annotates only its fully rendered Grafana Deployment `Replace=true`. This is
an Argo-managed API replace of the Deployment specification, not Force deletion
or a manual apply. The Deployment UID and source PVC must survive. All other
resources retain server-side apply, prune=false and existing ignore rules.
The exact complete Grafana replacement is server-dry-run validated first.

References: [Argo diff strategies](https://argo-cd.readthedocs.io/en/latest/user-guide/diff-strategies/)
and [resource-scoped replace](https://argo-cd.readthedocs.io/en/latest/user-guide/sync-options/#replace-resource-instead-of-applying-changes).
Once the transition is accepted, re-test server-side comparison before removing
these compatibility settings; do not restore them blindly and recreate the loop.


## Production acceptance receipt — 22:29–22:35 UTC

All three owning Applications (`monitoring-kube-prometheus-stack`,
`platform-monitoring`, `platform-backstage`) are **Synced / Healthy** against
`f145f14d690857e1d3d467c17478d16f6ee2a93f` (Helm source also 87.16.1).
No force/admin merge, check-policy change, manual workload apply, volume deletion
or Argo forced/selective sync was used. Argo naturally completed the transition.

- Grafana pod `monitoring-grafana-857cb95b54-nhlz6`: 3/3 ready, zero restarts.
  Deployment UID preserved through scoped API replace; same actual Grafana imageID
  and 13.1.0 version. `storage` uses retained/bound `grafana-data`, not emptyDir;
  Recreate and the verified-seed init guard are active.
- `grafana-state-verify-20261020` completed at 22:28:55 UTC using the enabled
  daily backup's code. Remote object:
  `backups/grafana/20261020T222853Z-d31d7772b0f4426abd666deb2bde993a/archive.tar.gz.aesgcm`.
  Download/decrypt/hash/isolated SQLite integrity and 87-table count verification
  passed (`restoreVerified: true`). Snapshot hash
  `315828849f573860b4835ffc97d5697bc24bd0cd7902bf575fbf4991de2879f9`.
  The recovery key at `~/.config/onemind/recovery/backup-encryption-key` already
  exists, is mode 0600 and matches the configured key (values never recorded).
  Daily 04:25 UTC schedule is enabled; one successful verification does not prove
  future recurring runs or independent offsite/hardware-loss recovery.
- Backstage generation 74, 2/2 ready/updated/available:
  `backstage-7c75c78b4-h2hp6`, `backstage-7c75c78b4-ncc7p`. Both actual imageIDs remain
  `ghcr.io/onemindos/onemind-portal@sha256:6c729bd20c0a0a40330be8f0c7bd5213c59c7540bd1ea54c38e05f8c63fc3e7c`.
  PostgreSQL specification and controller revision are unchanged; NATS remains 3/3.
- Each Backstage replica's installed credential returns 200 for Grafana `/api/user`,
  `/api/access-control/user/permissions` and `/api/search`. Identity remains
  `onemind-backstage-observer`, not Grafana admin; dashboards read allowed,
  dashboards write and service-account write absent. All three expected UIDs found.
  This proves the account/token persisted from original DB to the new PVC pod.
- Signed-in Backstage native dashboard discovery succeeded for Backstage,
  PostgreSQL and NATS: exact dashboard-title links resolve to their intended UIDs.
  A normal human catalog refresh was requested during source processing; catalog
  convergence, not the click itself, is the acceptance gate. No empty native
  Grafana Alerts card is rendered without explicit alert-label opt-in.
- Normal authenticated Grafana browser checks show ready replicas 2, 1 and 3,
  actual CPU/memory series for the intended pod names, and restart panels, without
  query/no-data/plugin errors on these three boards. Backstage's one-hour graphs
  correctly include prior rollout pods as historical series, not extra live replicas.
  All 12 underlying queries also returned finite live Prometheus series.
- Anonymous catalog and Grafana proxy GETs return 401. Existing machine identity
  can read catalog (200) but cannot reach Grafana proxy (403). Proxy remains
  authenticated, GET/HEAD-only; credentials stay backend-side. Impersonated portal
  Kubernetes service account can read Deployments but cannot read Secrets, exec,
  patch Deployments or delete pods. No runtime AI action permissions were widened.

Local screenshots (not committed binaries) are in
`/Users/zeus/.codex/visualizations/2026/10/20/grafana-visibility-2026-10-20/`:
`nats-panels.png`, `backstage-panels.png`, `postgres-panels.png`.
Redacted receipts are in `/tmp/onemind-grafana-ops/`: `final-checks.json`,
`denial-checks.json`, `seed-receipt.json`, `query-proof.json` and accepted objects.
These are local operator evidence, not a permanent remote artifact guarantee.

## Validation and remaining boundaries

Exact-head portable required gates passed for #340 (`dde10f9c`), #341
(`3a210492`) and #342 (`aad0704a`). The cutover cohort included immutable Yarn,
TypeScript, lint, 12 catalog checks, 30 pipeline/GitOps tests, 7 image contracts,
82 security/STAC and 44 layout checks plus YAML/schema/Helm and secret scanning.
Receipts: `/tmp/onemind-grafana-stage-a-check`,
`/tmp/onemind-grafana-stage-b-release-check`, `/tmp/onemind-grafana-sync-fix-check`.
GitHub billing prevented jobs starting; GitHub is **not** claimed passed.
`make validate` still reports pre-existing whole-tree YAML debt beginning at
`platform/onetak/media-mtx.yaml:86`; required changed-file gates passed separately.

Native plugin 1.1.0 renders an empty malformed folder link for root-level dashboard
search results (`...ts.netundefined/`). The named dashboard links and actual panels
work; folder-link presentation needs a bounded follow-up, not an assertion that
all native plugin UI is perfect. No iframe panels or Grafana SSO were introduced.
The old ClickHouse dashboard/plugin issue remains, as do application-level NATS
metrics/JetStream lag, full Prometheus alert integration and token-expiry monitoring.
The Viewer token has a 90-day lifetime from creation on October 20: rotate before
expiry in January 2027; the durable DB does not remove token lifecycle obligations.
No blanket all-app/COP certification, new controller, Phase 6 completion or Phase 7
migration is implied. Next scoped work: topology, alert integration and OpenCost.
