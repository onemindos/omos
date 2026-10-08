# Platform backups

Argo CD manages the retained ClickHouse job and repaired TAK PostgreSQL and NATS
jobs through this directory's Kustomization. Jobs run in their existing service
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

ClickHouse retains its previous backup implementation; it does not gain the new
encryption and restore assurance. This directory does not yet protect every PVC,
NATS account/server configuration, identity secret, certificate, or sealing key.
Do not treat it as complete disaster recovery or delete source volumes.
