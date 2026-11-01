# Nextcloud retirement — 2026-11-01

Zeus authorized retiring the service and deleting its 200 GiB file volume and
1 GiB VolSync cache to migrate to Dropbox. These archived manifests are not
referenced by an Argo Application. Restoring the original sealed credentials
and database is sensitive recovery work; do not redeploy this directory blindly.

Before live deletion: enable maintenance mode, export the Nextcloud PostgreSQL
database/configuration and recovery credentials privately, complete a fresh
DigitalOcean snapshot of the exact file volume, and record its recovery ID.
Never delete a volume merely because a snapshot request was accepted: verify
the completed snapshot refers to the correct source volume. Preserve existing
restic backups. Local Mac and Samsung originals are outside this teardown.

After Git merges, verify app-of-apps stops owning platform-nextcloud; remove
orphaned Nextcloud resources explicitly (the Application has no cascading
finalizer). The 200 GiB PV uses Retain, so namespace deletion alone does not
stop its DigitalOcean charge. Verify exact volume deletion at the provider.
Preserve the shared CNPG cluster and unrelated database roles.

A retained volume snapshot has a separate charge of $0.06/GiB/month; deleting
the 200+1 GiB live volumes removes $20.10/month of live-volume charges, not
all recovery storage costs. Snapshot restore/full Dropbox migration acceptance
remain separate from service retirement. Recovery records are private under
~/.config/onemind/nextcloud/retirement-20261101/ and are never committed.
