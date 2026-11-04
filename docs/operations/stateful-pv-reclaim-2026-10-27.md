# Stateful volumes on the Delete reclaim policy — 2026-10-27

Read-only finding and a proposed change. Nothing in the cluster was modified.
The inventory and reported backup-source flags below are observations dated
2026-10-27, not a current cluster census or a certification of recoverability.
Merging this document does not authorize any live patch, deletion or restore.

## Finding

The 2026-10-27 observation recorded the default StorageClass `do-block-storage`
with `reclaimPolicy: Delete`. Dynamically provisioned PersistentVolumes inherit the
StorageClass reclaim policy. With `Delete`, removing a claim can cause reclamation
of both the PV and its backing DigitalOcean volume: for example, after an Argo CD
prune, namespace deletion, Helm uninstall or mistaken `kubectl delete`.
PVC-in-use protection, finalizers and provider processing can delay deletion;
that delay is not a recovery guarantee. The same observation recorded retaining
classes `do-block-storage-retain` and `do-block-storage-xfs-retain`, but most
stateful workloads in this inventory did not use them.

On 2026-10-27, 62 of 85 PersistentVolumes were `Delete`. Most are VolSync Restic cache
volumes, which are disposable. These 21 hold application state:

| Namespace / claim | PersistentVolume | Size | VolSync source (reported 2026-10-27) |
|---|---|---:|---|
| `hermes/hermes-data` | `pvc-993e5ef7-9942-460b-801c-615eab02ae60` | 25Gi | yes |
| `hermes/hermes-fresh-data` | `pvc-4f05ab29-bf06-4c6c-8dec-c6521dd92a44` | 25Gi | yes |
| `karakeep/data-pvc` | `pvc-26e78589-9436-4c23-9c81-3c3c33738d88` | 20Gi | yes |
| `karakeep/meilisearch-pvc` | `pvc-7751b222-cac0-4a01-b923-c24f1decab73` | 5Gi | yes |
| `matrix/data-postgres-0` | `pvc-da722723-dc6a-45b2-88b8-41e92211e06d` | 20Gi | yes |
| `matrix/data-redis-0` | `pvc-66c94a56-6a97-41d5-8697-68d3099e93d0` | 10Gi | yes |
| `matrix/synapse-data` | `pvc-69d8e9ca-cd95-4bcc-91ac-13f7b85c55b5` | 50Gi | yes |
| `monitoring/prometheus-…-prometheus-0` | `pvc-4f24ac37-bc5c-4cc7-86dd-40141b07c350` | 20Gi | no |
| `monitoring/storage-loki-0` | `pvc-9b479c28-7a8c-4947-99b4-e66a146420a4` | 30Gi | yes |
| `monitoring/storage-tempo-0` | `pvc-4236680c-a00f-4d28-bdae-02a6391a033b` | 20Gi | yes |
| `nats/nats-js-nats-0` | `pvc-8fd0f545-e8e9-4146-9401-a68b5c9fdcdd` | 50Gi | yes |
| `nats/nats-js-nats-1` | `pvc-bdbf0a53-0b08-40a4-95ea-98c23f7911ed` | 50Gi | yes |
| `nats/nats-js-nats-2` | `pvc-79559c14-e5b1-4562-8ec5-2e579830cb54` | 50Gi | yes |
| `takml/takml-db-data` | `pvc-51e65cf5-7168-44af-81ff-492042f3d525` | 5Gi | yes |
| `takml/takml-models` | `pvc-362a3ed4-c5d7-4572-b98b-c22c6837af73` | 2Gi | yes |
| `takml/takml-server-jar` | `pvc-cc752a8e-1d96-45b0-b9e2-0817df61c510` | 1Gi | yes |
| `takserver/ldap-config-tak-ldap-0` | `pvc-d14388eb-7680-483d-8198-447151aa3d14` | 1Gi | yes |
| `takserver/ldap-data-tak-ldap-0` | `pvc-de870419-f658-4b80-8f46-5c13991cf6c0` | 1Gi | yes |
| `takserver/mongo-data-fedhub-mongodb-0` | `pvc-14f3be3b-f2e9-4228-96ef-edd973abb70f` | 5Gi | yes |
| `takserver/takserver-nats-js-takserver-nats-0` | `pvc-649550e9-edc5-4ee4-a8aa-ccae5407de1b` | 10Gi | yes |
| `vaultwarden/data-pvc` | `pvc-162b8f51-5099-4bc9-9878-f756ac068f3a` | 5Gi | yes |

## Retention and backup verification are separate

**Retention is not backup.** `Retain` preserves the same backing storage asset after
claim deletion for manual reclamation. It provides neither an independent recovery
copy nor protection against corruption, application-level deletion or storage
failure. Rebinding requires verified PV/PVC/CSI identity and application-consistent
recovery; this document promises neither a recovery time nor zero data loss.

A configured daily VolSync schedule does not guarantee a usable daily recovery
point. Recovery depends on the **last successful, accessible, restorable backup**,
available repository credentials and encryption keys, and a verified application
restore. A file backup may be only crash-consistent for a database or NATS; see the
[VolSync recovery caveats](../../platform/backups/volsync/README.md).

The table preserves the original reported yes/no flags, not newly verified backup
coverage. At the 2026-10-28 source review, four `yes` rows had no matching
ReplicationSource in [the committed source list](../../platform/backups/volsync/sources.yaml):
`hermes/hermes-fresh-data`, `karakeep/data-pvc`, `karakeep/meilisearch-pvc` and
`vaultwarden/data-pvc`. They may be managed elsewhere or exist only in live state;
absence from this file is not proof that no backup exists. Verify the current
source, last successful synchronization, repository/key access and restore evidence
for every claim before relying on its backups. This source review did not perform
a live backup-status check or a restore drill.

## Why not fix it in the manifests

`storageClassName` is immutable on a bound claim. Changing it in Git would make Argo
CD fail to sync those claims, and recreating them could delete the data this change
is meant to protect. The reclaim policy on an existing PersistentVolume is mutable.
Dynamic provisioning alone does not prove that no Argo Application or other
controller owns that field. Check the current manifests, ownership and
reconciliation behavior before any separately authorized patch; do not assume it
will be drift-free. Changing a StorageClass for future claims does not retroactively
change the reclaim policy of existing PVs.

## Proposed change (needs owner approval; not executed)

1. After separate owner approval, recheck each namespace/PVC/PV/CSI-volume binding,
   current policy and field ownership. Patch only the verified, explicitly approved
   PVs to `Retain`. This changes lifecycle policy; it does not create a backup or
   copy/migrate data. The table is historical, not an executable target list:

   ```sh
   for pv in <PersistentVolume names from the table>; do
     kubectl patch pv "$pv" -p '{"spec":{"persistentVolumeReclaimPolicy":"Retain"}}'
   done
   ```

2. Verify the approved PV identities, fully qualified claim bindings and policy:

   ```sh
   kubectl get pv -o custom-columns=PV:.metadata.name,NAMESPACE:.spec.claimRef.namespace,CLAIM:.spec.claimRef.name,VOLUME:.spec.csi.volumeHandle,RECLAIM:.spec.persistentVolumeReclaimPolicy
   ```

3. For new stateful claims, set `storageClassName: do-block-storage-retain` in the
   manifest. Leave VolSync cache claims on `do-block-storage`.

## Cost and cleanup

`Retain` keeps the backing DigitalOcean volume, and its billing, after the claim is
deleted. **Deleting a retained Kubernetes PV object does not delete that provider
volume.** Retirement therefore needs a separate, explicit owner-approved cleanup:

1. Before removing the PV object, record and independently verify the namespace,
   PVC, PV, CSI `volumeHandle` and exact DigitalOcean volume ID. Confirm it is unused
   and preserve the mapping and required recovery evidence.
2. Verify the required backups and restore evidence, and obtain explicit approval
   for the exact destructive targets. Retirement approval is not implied by this
   document or by approval to change reclaim policy.
3. Under that approved retirement procedure, remove the released Kubernetes PV
   object and separately delete the identified provider volume. Kubernetes does
   not perform the latter for `Retain`; do not delete anything if identity, usage,
   approval or backup checks are unresolved.
4. Read back both Kubernetes and DigitalOcean state to verify the intended objects
   are gone and no unintended volume was affected.

[RETIRE-ZEUS-SITE.md](RETIRE-ZEUS-SITE.md) describes a separate retirement of
`Delete`-policy volumes. It is not a cleanup procedure for `Retain` volumes.
No retirement action has been executed as part of this source-only correction.

## Primary policy references

- [Kubernetes PV reclamation and protection](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
- [Changing a PV reclaim policy](https://kubernetes.io/docs/tasks/administer-cluster/change-pv-reclaim-policy/)
