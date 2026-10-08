# NATS workload ownership and memory

`platform-nats` now adopts the existing `nats` StatefulSet. The Helm release still
owns configuration, Services and related resources; do not run a Helm upgrade
that overwrites the Git-owned StatefulSet. Full Helm-to-Argo configuration
migration remains unfinished, especially authentication and certificate handling.
The manifest preserves selectors, service identity, TLS mounts, replica count,
PVC templates and retention. Existing volumes are not recreated.

The October 8 production snapshot attempt on the roughly 4.9 GB COT stream caused
OOM kills with the previous 1 GiB container limit. The workload sets a 2 GiB
container limit and `GOMEMLIMIT=1536MiB`, with a 1 GiB scheduling request. The Go
budget is soft; monitor RSS and OOM events during actual snapshots.

Rolling-update partition starts at 2. Verify nats-2 readiness, route connections,
JetStream quorum and stream replication, then lower partition to 1 and finally 0
through reviewed Git changes. A Ready pod alone does not establish stream quorum.
Do not force-delete multiple replicas or delete PVCs. The existing PDB permits
only one voluntary disruption.

Scheduled snapshots remain suspended during sizing verification. Resume only
once the real encrypted download/restore job succeeds without server OOMs.
The manifest's DigitalOcean storage class describes current volumes; a new home
cluster needs an explicit storage-class overlay and data restoration.
