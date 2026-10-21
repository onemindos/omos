# DigitalOcean volume recovery points

`snapshots-20261008.yaml` requests one recovery point for each of the 35 bound
PVCs found during the readiness review. The dedicated snapshot class retains the
provider snapshot and VolumeSnapshotContent if a Kubernetes reference is removed.
These are a single dated baseline, not a recurring schedule or retention policy.
Do not prune provider snapshots until a newer verified recovery point exists.
Retained provider snapshots consume billable storage.

Wait for every VolumeSnapshot's `status.readyToUse=true`, retain its content name
and provider snapshot handle in the recovery inventory, and test a restore to a
separate PVC. Never change the production claim binding for a restore test.

These are crash-consistent snapshots of individual volumes. They are not atomic
across NATS replicas, database/media pairs or multiple applications. Prefer the
logical application archives for tested database recovery. A snapshot becoming
Ready does not prove application restoration.

DigitalOcean CSI supports snapshot and restore without the user's local doctl
session. These snapshots remain DigitalOcean-specific: home migration still
requires portable logical archives and filesystem exports, a home storage-class
overlay, recovery keys and independent verification. Droplet backups are not a
replacement for persistent-volume or application backups.

References:
- https://docs.digitalocean.com/products/kubernetes/how-to/create-snapshots/
- https://docs.digitalocean.com/products/kubernetes/how-to/restore-volumes/
- https://kubernetes.io/docs/concepts/storage/volume-snapshots/
