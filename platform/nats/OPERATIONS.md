# NATS workload ownership and memory

`platform-nats` now adopts the existing `nats` StatefulSet. The Helm release still
owns configuration, Services and related resources; do not run a Helm upgrade
that overwrites the Git-owned StatefulSet. Full Helm-to-Argo configuration
migration remains unfinished, especially authentication and certificate handling.
The manifest preserves selectors, service identity, TLS mounts, replica count,
PVC templates and retention. Existing volumes are not recreated.

The October 8 production snapshot attempt on the roughly 4.9 GB COT stream caused
OOM kills with the previous 1 GiB container limit. The workload sets a 2 GiB
container limit and `GOMEMLIMIT=1536MiB`, with a 512 MiB / 100m scheduling request and required hostname separation. The Go
budget is soft; monitor RSS and OOM events during actual snapshots.

The verified rollout is now at partition 0. For future high-risk rollouts, start at 2. Verify nats-2 readiness, route connections,
JetStream quorum and stream replication, then lower partition to 1 and finally 0
through reviewed Git changes. A Ready pod alone does not establish stream quorum.
Do not force-delete multiple replicas or delete PVCs. The existing PDB permits
only one voluntary disruption.

The encrypted production backup restored all 39 streams into an isolated server
without production restarts; the schedule is enabled. Repeat restore checks after
changes to versions, sizing or credentials.
The manifest's DigitalOcean storage class describes current volumes; a new home
cluster needs an explicit storage-class overlay and data restoration.

## Network access

`network-policy.yaml` is an ingress allow-list on the three server pods: routes
(6222) only between the servers; leaf nodes (7422), clients (4222), WebSocket
(8080) and monitoring (8222) only from the namespaces and pods it lists. A new NATS
client elsewhere needs a rule in the same change, otherwise its connection times
out without any error in the NATS logs; look for DROPPED flows into namespace
`nats` in Hubble. If a missing rule takes a client down, revert the change in Git.
Until the revert syncs, an operator can apply an extra allow-all ingress policy for
the server pods (Argo does not track it, so it is not pruned) and delete it after.
A policy removed from Git is deleted only while `platform-nats` prunes; otherwise
delete it by name.
