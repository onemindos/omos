# Backup and recovery operations

The backup service is Git-owned under `platform/backups`; portable filesystem
exports are under `platform/backups/volsync`. A successful job is not sufficient:
logical archives write a verified marker only after downloading, authenticating
and restoring into an isolated process. Paused filesystem sources have no completed
backup coverage until an actual export succeeds.

Review `docs/operations/cluster-readiness-2026-10-08.md` for observed evidence and
remaining gaps, and `docs/operations/home-migration.md` for migration acceptance.
The local logical-archive key and Restic password live under
`~/.config/onemind/recovery` with restricted permissions. Keep independent offline
copies alongside the sealing-key recovery plan; do not upload cleartext keys to Git.

Read-only operational checks:

```sh
kubectl --context do-nyc3-omos-cloud get cronjobs,jobs -A
kubectl --context do-nyc3-omos-cloud get replicationsources -A
kubectl --context do-nyc3-omos-cloud get volumesnapshots -A
kubectl --context do-nyc3-omos-cloud -n argocd get applications
```

Inspect only the relevant backup container logs. Do not dump Secrets, full process
environments or old credential-bearing manifest diffs. Restore drills write only
to disposable processes or new explicitly named test PVCs. Preserve original PVCs,
retained snapshots, encryption keys and intentional mesh routes during cleanup.

The Backstage catalog provides operational links and inventory. Existing AI access
remains read-only; this service entry does not grant access to credentials or restore
operations. Provider-specific snapshots and portable encrypted archives have
different recovery uses. Consult the runbook before any production restore.
