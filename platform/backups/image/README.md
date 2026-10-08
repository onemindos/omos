# Application backup image

This image replaces the legacy shell pipelines for PostgreSQL 15 and NATS.
It refuses empty NATS backups, propagates command failures, and uses authenticated
NATS connections configured through environment variables, including `NATS_CA`.
Archives include all visible streams, including KV streams and consumers. It does
not back up other NATS accounts, server credentials, or operator configuration.

Archives are streamed through AES-256-GCM encryption to temporary files rather than
held in memory. S3 endpoint, region, bucket and credentials are configurable, so
the storage provider can change independently of the cluster. Existing Backstage
archives retain their separate format and associated-data value.

Every run downloads the uploaded object, compares SHA-256, authenticates/decrypts
it, and performs a real restore. PostgreSQL uses a unique temporary database on
a separate loopback-only sidecar and drops it afterward. Source data and source
database disk capacity are never consumed by the restore check. The sidecar must
listen on port 15432 (or `PG_RESTORE_PORT`) with the required extensions installed.
NATS starts a loopback-only disposable JetStream server without production
credentials and restores each stream with one replica. A `.verified.json` marker
is uploaded only after all checks pass. Jobs need enough temporary disk for the
archive, remote copy, decrypted copy and isolated restore.

PostgreSQL dump/restore clients must match the database major version. The default
image is PostgreSQL 15; do not use it for PostgreSQL 17 databases. Temporary restore
databases need the original extensions available and appropriate creation rights.
Concurrent writes mean separate stream snapshots are not one atomic account-wide
snapshot. A successful restore verifies structure and snapshot readability; it
does not exercise all application workflows or replace disaster-recovery drills.

`BACKUP_ENCRYPTION_KEY` is base64 containing 32 random bytes. Preserve it outside
the cluster, with the S3 credentials and recovery instructions. Do not rotate it
without retaining old keys needed to decrypt old archives. The image does not
delete archives; configure reviewed retention policies at the storage provider.

The workflow runs integrity/isolation unit tests and synthetic PostgreSQL/NATS
backups against a local S3 emulator before publishing. Kubernetes adoption and
actual production restore checks are separate deployment steps. Keep backups of
all other databases, files, identity and sealing keys in the recovery inventory;
this image alone does not provide complete cluster backup coverage.

`BACKUP_KIND=postgres-owned` discovers all non-template databases owned by the
connected application role and restores each into the isolated sidecar. Backstage
uses this mode because its plugins own separate databases; an individual plugin
may legitimately have no tables. The complete set must contain application tables
before verification succeeds. Owned-database archives record a manifest mapping
names to hashed dump filenames. Roles and grants still come from deployment secrets
and bootstrap configuration. Dumps use compression level 1 to reduce backup CPU.
