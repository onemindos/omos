# Portable filesystem backups

VolSync 0.16.0 uses Restic to export CSI snapshot clones into independent encrypted
S3 repositories. Source volumes stay attached to their applications. The operator
and Restic mover image are pinned by digest; Argo manages the operator in the
existing ops namespace. Controller RBAC manages temporary volume/snapshot/job
resources across the service namespaces. Backup credentials remain namespace-bound
SealedSecrets. This adds backup-controller access, not agent or portal access.

A separate Restic password is preserved outside Kubernetes at
`~/.config/onemind/recovery/restic-password` with restricted permissions. Keep an
independent offline copy. Preserve the sealing keys and existing logical-archive
key as well. Replacing the object-store access key does not replace encryption keys.

Each source has its own repository prefix. Snapshot/clone/cache classes describe
current DigitalOcean storage and must be replaced for a home cluster. Restic
archives themselves can restore onto another provider's filesystem. Raw database
and NATS volume copies remain crash-consistent; use verified logical archives for
coordinated database recovery. A completed file backup is not an application restore.

Begin with the Hermes volume and verify an isolated restore before enabling the
remaining sources. Do not claim coverage for paused sources. Stagger schedules and
monitor CPU, memory, temporary volumes, cache storage and provider snapshot cleanup.
Source PVCs must never be deleted during backup cleanup.

References:
- https://volsync.readthedocs.io/en/latest/usage/restic/
- https://github.com/backube/volsync/releases/tag/v0.16.0
