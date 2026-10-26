# TAK StatefulSet adoption recovery preflight — 2026-10-26

Recovery evidence was checked before the proposed adoption of the existing
`takserver-nats` and `takserver-postgresql` StatefulSets into
`platform-takserver`. This receipt does not claim the adoption or later resource
budget changes have happened. Adoption must preserve the live pod templates,
claims, images, configuration and credentials; resource changes belong in a
separate reviewed change after no-op adoption is verified.

The scope is **TAK's internal NATS in namespace `takserver`** and TAK PostgreSQL.
The authority NATS cluster in namespace `nats` was not accessed by this restore
exercise. The applicable adoption gate is
[`docs/gitops/HELM-ADOPTION.md`](../gitops/HELM-ADOPTION.md).

## Fresh encrypted volume archives

The existing VolSync sources use CSI `Snapshot` copies and their established
Restic repositories. Both acknowledged the Git-managed manual tag
`policy-preflight-20261026T073053Z` and reported `Successful`:

| Workload | Source PVC | ReplicationSource | New Restic snapshot | Completed UTC | Archived data |
|---|---|---|---|---|---|
| TAK internal NATS | `takserver-nats-js-takserver-nats-0` | `volume-2b88c0eca294` | `dc63b8ad` | 2026-10-26 07:33:58 | 5 files, 387 bytes |
| TAK PostgreSQL | `data-takserver-postgresql-0` | `volume-cccc1afaa021` | `50368b34` | 2026-10-26 07:34:02 | 2,128 files, 9.557 GiB |

Both repositories and their credentials were preserved. The one-time tags must
be removed through Git/Argo after acceptance so the existing daily schedules can
resume: PostgreSQL at 09:40 UTC and internal NATS at 11:00 UTC. A nonempty VolSync
manual tag takes precedence over the cron schedule.

Read-only inspection also confirmed the existing retained CSI checkpoints
`premove-20261017-data-takserver-postgresql-0` and
`premove-20261017-takserver-nats-js-takserver-nats-0` remain `readyToUse: true`, with
20Gi and 10Gi restore sizes respectively. These are October 17 checkpoints, not
new retained October 26 snapshots. The fresh October 26 archives above came from
VolSync's snapshot-copy workflow; its temporary CSI objects are not claimed as
retained checkpoints.

## PostgreSQL logical restore receipt

The latest daily `takserver/backup-takserver-postgres-29883000` Job started at
**2026-10-26 02:00:00 UTC** and succeeded at **04:20:19 UTC**. Its backup container's
JSON receipt reported `kind: postgres`, `restoreVerified: true` and **64 tables**
for this encrypted logical archive:

`backups/takserver-postgres/20261026T020038Z-a8a2ee2501b34fbfacd536bd473067c4/archive.tar.gz.aesgcm`

The deployed backup workflow creates a PostgreSQL custom-format dump of `cot`,
uploads the encrypted archive, downloads it again, checks its digest, authenticates
and decrypts it, and runs `pg_restore --exit-on-error` into a temporary database in
a separate PostGIS sidecar. That restore explicitly uses loopback
`127.0.0.1:15432` rather than production PG connection variables, counts restored
application tables, then drops its temporary database. The remote verification
marker is written only after this restore succeeds. The workflow is defined in
[`platform/backups/backups.yaml`](../../platform/backups/backups.yaml) and
[`platform/backups/image/backup.py`](../../platform/backups/image/backup.py).

This is factual evidence of the latest isolated **logical** PostgreSQL restore.
The fresh 9.557 GiB raw PostgreSQL volume archive was **not** downloaded or restored
again in this bounded preflight. No claim is made that this exercise tested a raw
PostgreSQL volume restore or a full running TAK application recovery.

## Actual internal NATS archive restore and startup

At **2026-10-26 07:36:46 UTC**, a read-only verifier retrieved the fresh encrypted
internal NATS archive:

`dc63b8ad55b98320001d2d27db5242aa2833f4d387f4f94c4f65ac70c410e23e`

The snapshot time was **07:33:53.558312252 UTC**. The verifier required the source's
expected PVC, completed manual tag and successful mover result, then used an
official Restic 0.18.0 binary with its release archive checksum revalidated. Only
`--no-lock snapshots`, `ls` and `dump` were used against the repository. Secret
values stayed in process memory and the Restic subprocess environment; they were
not printed, written to an environment file or passed as command arguments.

All **5 files / 387 bytes** were restored into task scratch with directories 0700 and
files 0600. The archive was below the authorized 200MiB retrieval limit. Links and
special nodes were rejected. The live source claim was neither mounted nor
written. No remote locks, initialization, pruning or repository writes occurred.

The isolated server used the exact amd64 image digest observed on the running
`takserver-nats-0` pod:

`docker.io/library/nats@sha256:511f5c4cfc6fdd61eb66afab99dfb38bed69aae630d8d5b36bc9bfc716723cd8`

It reported **NATS 2.9.21**. The restored store was mounted at `/data`, matching the
production JetStream storage path and 10Gi file-store budget. A minimal standalone
configuration was used for the recovery check; production configuration and
credentials were not copied into it. Docker used `network=none`, no published
ports, a read-only root filesystem and all Linux capabilities dropped. The only
writable application store was the disposable restored copy.

The server's loopback `/healthz?js-enabled-only=true` returned `status: ok`, and
`/jsz` confirmed JetStream enabled with **0 streams, 0 consumers, 0 messages and 0
stored bytes**. No local test messages or streams were created. Consequently this
proves actual archive decryption, recovery of the current metadata-only store and
startup with the exact server version. It does **not** prove recovery of populated
JetStream streams, consumers or messages; none were present in this archive.

The test container, decrypted restored files, isolated configuration and Restic
cache were removed after the check. A bounded receipt containing counts, times,
image/snapshot identifiers and outcomes was retained; no NATS payloads or secret
values were retained in the receipt.

## Remaining adoption verification

The recovery gate above does not replace the adoption checks. Before accepting
adoption, compare server-defaulted manifests to live state and require no pod
specification or volume-template change. After the first Argo sync, confirm the
same StatefulSet pod UIDs/restart counts and bound claims, healthy application
connections, and daily backup schedules after removal of the one-time manual
tags. Change resource budgets only in the subsequent reviewed rollout.


## No-op adoption acceptance

PR [#474](https://github.com/onemindos/omos/pull/474) brought the two existing
StatefulSets and retention CronJob under `platform-takserver`. The server-defaulted
specification of each resource matched its captured baseline and fresh live reads
exactly. The original PostgreSQL and internal NATS pod UIDs stayed unchanged;
both source PVC/PV identities and capacities were preserved. Argo reported
Synced/Healthy. No container was restarted by this ownership step.

Both temporary manual backup triggers were removed. Their daily source
`nextSyncTime` values resumed at 2026-10-26 09:40 UTC (PostgreSQL) and 11:00 UTC
(internal NATS). The following resource-budget rollout is separate and intentionally
restarts the single-replica workloads one at a time.
