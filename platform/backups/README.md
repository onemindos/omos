# Platform backups

Argo CD manages logical PostgreSQL/NATS jobs and encrypted Kubernetes recovery
exports through this directory's Kustomization. VolSync manages portable file
backups under the separate volsync application. Jobs run in their existing service
namespaces. Automatic pruning is disabled while adopting previously manual
resources; do not delete credentials or archives as part of adoption.

TAK and NATS use the tested, digest-pinned image described in
[image/README.md](image/README.md). Every run encrypts an archive, uploads it,
downloads it again, checks integrity, and restores into isolated temporary
storage before writing a verification marker. TAK restores into a loopback-only
PostGIS sidecar; NATS restores into a loopback-only disposable server. Neither
verification writes to the source service.

The namespace-bound SealedSecrets reuse the existing backup credentials and
encryption key. Preserve that encryption key and the Sealed Secrets controller
private keys independently of the cluster. Sealing is not credential rotation:
the previously committed Spaces credentials still require provider-side rotation.

## Verification and limits

After deployment, create a one-off Job from each CronJob, inspect its completion
and sanitized logs, and verify its `.verified.json` archive marker. Image CI
tests synthetic PostgreSQL and NATS restores; this does not establish successful
production backups. Record production results in the cluster review.

TAK has approximately 10 GB of source data and NATS approximately 5 GB at the
October 8 audit. Jobs request local ephemeral space and bound archive/restore
directories. Monitor disk and execution time as these grow. Backups are scheduled
in UTC. PostgreSQL 17 services require a matching client and restore image.

ClickHouse's old unverified native job is suspended. Its encrypted Restic export
passed an isolated server boot; all 35 original volumes completed initial checks.
Kubernetes recovery exports include identity/configuration and active sealing keys,
but their verification is decryption/equality, not replay onto a new cluster.
See [the evidence report](../../docs/operations/cluster-readiness-2026-10-08.md).
Do not treat this as complete disaster recovery or delete source volumes.

## Spaces credential rotation gate

The existing access key previously appeared in Git history. Removing plaintext
manifests and sealing live credentials does not revoke it. Provider authentication
still returns 401; complete these steps after local account access is restored:

1. Create a replacement key scoped to the backup bucket and required object
   operations. Preserve archive encryption keys and the Restic password unchanged.
2. Update all namespace-bound logical/recovery and Restic credential SealedSecrets
   through a protected local process. Never print values, use command-line literals,
   publish unredacted diffs or commit cleartext Secret manifests.
3. Merge reviewed ciphertext changes through GitOps. Verify a new logical upload,
   download and restore, plus a Restic backup and access to existing repositories.
4. Revoke the old key only after every consumer uses the replacement and independent
   recovery access is tested. Record revocation evidence, then review any remaining
   legacy consumers and historical exposure separately. Do not rewrite Git history
   as a substitute for revocation.

Provider Droplet backups, Kubernetes control-plane protection and PVC backups are
different mechanisms. Confirm account-specific coverage with authenticated API
access; a node backup cannot replace the application/volume recovery tested here.
Retained baseline snapshots have no automatic expiry. Portable Restic retention
is seven daily, four weekly, three monthly and one yearly snapshots, with pruning
every 14 days; observe its first real cycle and review costs before expansion.
