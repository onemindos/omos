# Stateful volumes on the Delete reclaim policy — 2026-10-27

Read-only finding and a proposed change. Nothing in the cluster was modified.

## Finding

The default StorageClass `do-block-storage` has `reclaimPolicy: Delete`. Every claim
created from it gets a PersistentVolume that DigitalOcean destroys as soon as the claim
is deleted: by an Argo CD prune, a namespace deletion, a Helm uninstall or a mistaken
`kubectl delete`. The retaining classes `do-block-storage-retain` and
`do-block-storage-xfs-retain` exist but most stateful workloads do not use them.

On 2026-10-27, 62 of 85 PersistentVolumes were `Delete`. Most are VolSync Restic cache
volumes, which are disposable. These 21 hold application state:

| Namespace / claim | PersistentVolume | Size | VolSync source |
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

VolSync backups mean a deleted claim is recoverable, but only to the last Restic run
(daily) and only after a restore. NATS (tier 1), TAK LDAP and FedHub, Hermes and
Vaultwarden are on this list. A `Retain` volume survives the claim and can be rebound
in minutes with no data loss.

## Why not fix it in the manifests

`storageClassName` is immutable on a bound claim. Changing it in Git would make Argo
CD fail to sync those claims, and recreating them would delete the data this change
is meant to protect. The reclaim policy on an existing PersistentVolume is mutable
and is not owned by any Argo Application (the volumes are dynamically provisioned),
so patching it causes no drift.

## Proposed change (needs owner approval; not executed)

1. Patch each volume in the table to `Retain`. This is non-disruptive and reversible:

   ```sh
   for pv in <PersistentVolume names from the table>; do
     kubectl patch pv "$pv" -p '{"spec":{"persistentVolumeReclaimPolicy":"Retain"}}'
   done
   ```

2. Verify:

   ```sh
   kubectl get pv -o custom-columns=CLAIM:.spec.claimRef.name,RECLAIM:.spec.persistentVolumeReclaimPolicy
   ```

3. For new stateful claims, set `storageClassName: do-block-storage-retain` in the
   manifest. Leave VolSync cache claims on `do-block-storage`.

## Cost and cleanup

`Retain` keeps the DigitalOcean volume billing after its claim is deleted. When a
workload is retired on purpose, delete the released PersistentVolume and confirm the
DigitalOcean volume is gone, as `docs/operations/RETIRE-ZEUS-SITE.md` already does.
