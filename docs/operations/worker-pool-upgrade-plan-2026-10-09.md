# NYC3 worker replacement plan

The authenticated provider offered Kubernetes 1.36.5-do.0, but the upgrade API
rejected the request with HTTP 422: insufficient NYC3 capacity for s-8vcpu-16gb.
No upgrade began. The provider size inventory excludes NYC3 for that SKU.

## Concrete alternative awaiting approval

Replace the three core workers with s-8vcpu-16gb-amd in the same NYC3 cluster/VPC.
The provider lists this alternative for NYC3; actual allocation is still subject
to capacity at execution time. It preserves eight vCPUs and 16 GiB per core node.
Current core-worker pricing is $96 each/month; the alternative is $112 each/month:
$288 to $336/month, a $48/month increase, excluding temporary overlap, volumes
and other cluster charges. The two existing ecosystem workers remain separate.

1. Export node-pool labels/taints and audit workload node selectors/affinity.
2. Verify all 49 original application PVCs remain Bound, all 49 retained baseline
   snapshots are Ready, fresh recovery exports and application backups pass, and
   representative restore checks are retained.
3. Create three replacement workers with equivalent workload labels/taints.
   Verify Ready nodes and sufficient scheduling reservations before draining.
4. Cordon and drain one old worker at a time, respecting disruption budgets.
   Verify PVC reattachment, NATS quorum/peer currency, database health and app
   readiness before proceeding. Single-instance apps can experience interruption.
   Never force-remove PVs or disable surge/PDB safeguards to bypass a failure.
5. Remove the empty old core pool only after workload recovery is verified.
6. Upgrade to the account-supported fixed patch release with surge enabled;
   reverify node versions, kubelet RSS, volumes, NATS, backups and Argo health.
7. Record resulting provider settings and reconcile infrastructure ownership.
   Existing Terraform defaults are not verified state for this live five-node
   cluster; do not run make apply as a migration shortcut.

This is a proposed operational plan, not authorization to create the paid pool
or evidence that the upgrade/replacement has occurred. Home deployment still needs
its own hardware/storage overlay and disconnected recovery exercise.
