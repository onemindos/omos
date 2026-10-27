# NATS

The authority bus: a three-node JetStream cluster in namespace `nats` that every cloud workload publishes to. Cloud workloads connect to `nats://nats.nats.svc.cluster.local:4222` with credentials from a Secret; Hermes is the one exception, through its own leaf sidecar. Edge and Mac bridges run leaf nodes that only forward, through the `omos-nats` Tailscale load balancer. If an edge machine is offline, the cloud bus keeps running.

Sources: root `CLAUDE.md` ("Bus Authority", "NATS Streams"), `platform/nats/OPERATIONS.md`, `docs/gitops/HELM-ADOPTION.md`, `docs/NATS-DOCTRINE.md`, `platform/backups/README.md`.

## Ownership

- `platform-nats` applies `platform/nats` and owns the `nats` StatefulSet. The Helm release `nats` still owns configuration, Services and related resources, with its passwords in an out-of-band values file. Never run a Helm upgrade that overwrites the Git-owned StatefulSet; the full Helm-to-Argo migration, especially authentication and certificates, is unfinished.
- Stream and KeyValue definitions are kept in `platform/nats/streams/`, and the NACK controller (`edge-nack`) reconciles Stream and KeyValue resources into JetStream. `platform-nats` does not recurse into that directory, so confirm how a definition reaches the cluster before relying on it. The older Python provisioning scripts are not the source of truth.
- Subjects follow the locked grammar in `docs/NATS-DOCTRINE.md`. Do not add a subject prefix without updating the doctrine and the affected schemas and consumers.
- Server TLS is issued by cert-manager inside the namespace; see [cert-manager](cert-manager.md).

## Investigate (read-only)

```sh
kubectl --context do-nyc3-omos-cloud -n argocd get application platform-nats edge-nack
kubectl --context do-nyc3-omos-cloud -n nats get statefulsets,pods,pvc,poddisruptionbudgets
kubectl --context do-nyc3-omos-cloud -n nats get events --sort-by=.metadata.creationTimestamp
kubectl --context do-nyc3-omos-cloud -n nats get cronjob backup-nats-streams
```

A Ready pod does not establish stream quorum: check route connections, JetStream quorum and stream replication in the dashboards. A snapshot of the large COT stream once caused OOM kills at a 1 GiB limit; the server container now has a 2 GiB limit and `GOMEMLIMIT=1536MiB`, which is a soft budget, so watch memory and OOM events during snapshots. Further symptoms are in [NATS fabric](../nats.md).

## Changes and recovery

- Roll out a risky StatefulSet change with the update partition at 2. Verify `nats-2` readiness, route connections, JetStream quorum and stream replication, then lower the partition to 1 and to 0 through reviewed Git changes. The PodDisruptionBudget permits one voluntary disruption.
- Never force-delete several replicas or delete PVCs. Do not reset streams, purge messages or recreate consumers as a generic fix.
- `backup-nats-streams` (namespace `nats`, daily at 03:30 UTC) writes an encrypted stream archive, restores it into a loopback-only disposable server, and only then writes its verification marker. In the readiness exercise recorded in `docs/operations/cluster-readiness-2026-10-08.md`, the production archive restored every stream into an isolated server without production restarts. Repeat the restore check after changes to versions, sizing or credentials.
- VolSync also copies the JetStream volumes. Those copies are crash-consistent and not atomic across replicas, so prefer the verified stream archive.
- A new cluster needs an explicit storage-class overlay and data restoration; the manifest's storage class describes the current DigitalOcean volumes.

## Verify completion

Brokers and storage are Ready and quorum and replication are confirmed. An authorized publish and consume round trip on a dedicated test subject succeeds with non-sensitive data, the affected consumer resumes and its lag falls. Dashboard availability is checked separately from message delivery.

Generated workload list: [NATS JetStream app runbook](../apps/nats.md).
