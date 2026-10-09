# Cluster readiness follow-up — 2026-10-09

Verified against `do-nyc3-omos-cloud`; continues the October 8 audit. This is
observed evidence, not a guarantee of reliability or a completed home migration.

## Completed follow-through

- DigitalOcean authentication works. Five workers are Ready, still Kubernetes
  1.36.0. Two ecosystem workers were added outside this repair session.
- Overnight scheduled logical backups completed, including TAK and NATS. All
  46 filesystem backup sources have completed an initial/scheduled transfer.
  Coverage now includes the eleven new Dapr, model, Harbor, OpenSearch and SPIRE
  sources. Daily transfers are staggered through 15:00 UTC (PR 161).
- Fourteen new retained CSI snapshots are Ready, extending the original 35
  application-volume baseline to 49. CSI snapshots remain provider-specific
  and crash-consistent; they do not establish application restore correctness.
- CloudNativePG native backup completed using the replacement credentials.
  An isolated cluster restored the production backup and reached Ready. Every
  production database had matching user-table counts: 7 databases, 114 tables.
  The restore bootstrap also created an empty `app` database. This verifies
  startup and schema counts, not equality of every row or user workflow.
  Temporary test cluster and PVC were deleted. Evidence accompanies this report.
- Replaced the exposed Spaces credential across 51 existing consumers, including
  Velero's embedded credential file. New backups and provider access passed;
  the old bucket credential was revoked. New source credentials use the replacement.
  Encryption and Restic recovery keys were preserved. PR 155 contains sealed data.
- Harbor's RWO rollout deadlock was repaired using the already-merged Recreate
  strategy and removal of incompatible default rolling-update settings. Its
  active htpasswd Secret was adopted with unchanged bytes into SealedSecrets
  rather than pruned (PR 160). Harbor was verified Healthy and Synced.
- OpenChoreo assistant no longer crashes on an ambiguous model provider. Its
  private Gateway was created after chart/Gateway API compatibility repairs;
  assistant pods are Ready and application health is Healthy (PRs 157–158).
  Authenticated mesh HTTPS and full user workflows still need acceptance tests.
- Fresh cluster recovery export completed after credential adoption. Original
  checkout changes, including newer mission edits, were preserved.

## Remaining and blocked work

The authenticated upgrade API offers 1.36.5-do.0 but rejected the upgrade with
HTTP 422: insufficient NYC3 capacity for the current s-8vcpu-16gb worker size.
No permanent Kubernetes patch upgrade occurred. The kubelet restart mitigation
from October 8 is temporary. The attached worker replacement plan proposes
same-capacity AMD workers at an additional $48/month plus temporary overlap;
that spending decision is awaiting the user. Do not disable surge or force drain
through disruption budgets as a workaround.

Worker Droplet backups are disabled on all five nodes. These backups would not
replace PVC/application backups. Velero backups complete but have volume
snapshots and default filesystem backup disabled; their successful status covers
objects, not the application's persistent data. CNPG native and VolSync backups
are separate data protection mechanisms.

Repository-wide validation still reports 498 YAML errors and 7 warnings. Focused
changed manifests passed parsing/server dry runs and PR CI; the larger formatting
and workflow debt is unfinished. Existing Terraform defaults are not verified
state for this five-node cluster. Manual/Helm/Argo ownership reconciliation,
per-application user acceptance, and a disconnected home-hardware restore remain
open. No namespace was confirmed disposable, so none was removed.

The refreshed independent local recovery copy contains 54.43 GiB of encrypted
objects. All 42 nonempty filesystem repositories passed `restic check --read-data`;
four sources were confirmed empty. Eleven encrypted archives passed checksum
verification. The fresh Kubernetes export included 336 Secrets, 213 ConfigMaps
and two sealing keys; authenticated download/decrypt equality passed. These
exports have not been replayed into a new cluster. A separate physical/offsite
recovery destination remains needed for failure of this computer.

PR 163 repairs the identity bootstrap's missing required OAuth invalidation flow
and enables server-side diff for OpenChoreo CRD defaults. The corrected identity hook completed successfully and OpenChoreo was verified
Healthy and Synced. Live login/user workflow acceptance remains separate.
