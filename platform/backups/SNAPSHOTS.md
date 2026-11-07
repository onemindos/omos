# DigitalOcean volume recovery points

`snapshot-class.yaml` defines `onemind-backup-retain`, a Retain-policy
VolumeSnapshotClass for one-off safety snapshots before risky changes (shrinks,
upgrades, rebuilds). Deleting a VolumeSnapshot made with it keeps the
VolumeSnapshotContent and the provider snapshot. Remove the provider copy
separately once it is no longer needed. Retained provider snapshots consume
billable storage.

Recurring recovery points come from VolSync (`volsync/`), the database dumps and
Barman. They are not snapshots of this class.

## History

The October readiness baselines (`snapshots-20261008.yaml` and
`snapshots-20261009.yaml`, 38 recovery points), the 2026-10-17 pre-adoption and
pre-move snapshots, and the one-off Authentik snapshot restore test were retired
on 2026-11-07. Daily VolSync backups of the same volumes had replaced them. Their
provider snapshots were deleted by the owner with doctl.

When taking a new safety snapshot:

- Wait for `status.readyToUse=true`.
- Record the content name and provider handle.
- Test a restore to a separate PVC. Never change the production claim binding.

These are crash-consistent snapshots of single volumes. They are not atomic
across NATS replicas, database and media pairs, or multiple applications. Prefer
logical application archives for database recovery. They are also specific to
DigitalOcean. Moving off DigitalOcean still needs portable archives and
filesystem exports.

References:
- https://docs.digitalocean.com/products/kubernetes/how-to/create-snapshots/
- https://docs.digitalocean.com/products/kubernetes/how-to/restore-volumes/
- https://kubernetes.io/docs/concepts/storage/volume-snapshots/
