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

VolSync namespace opt-in grants movers only DAC_OVERRIDE, CHOWN and FOWNER
capabilities, with privileged=false, no host mounts and read-only source snapshot
clones. These allow preserving arbitrary application ownership and permissions.
n8n retains its restricted namespace policy and uses the original application UID
(1000) or PostgreSQL UID (70), without added capabilities. A restore must verify
ownership and data-directory modes before starting databases.

Authentik's filesystem export was restored to a new PVC and PostgreSQL booted with
230 application tables. It now has a daily schedule. The next five-volume batch
uses manual initial triggers; other sources stay paused. TAK's original five-PVC
quota blocked backup caches and clones. The backup application now manages its
quota at 18 PVCs / 100 GiB, allowing the original 37 GiB, five caches, concurrent
clone space and a restore drill. This is an admission ceiling, not provisioned
storage. Preserve this override during any legacy TAK Helm upgrade.
