# Offsite backup copy (second DigitalOcean region)

Every backup today lands in one Spaces bucket, `onemind-backups` in `nyc3` (same
provider and region as the cluster). `backup-offsite-copy` (namespace `ops`,
[offsite-copy.yaml](offsite-copy.yaml)) keeps a second copy in a bucket in another
region so a nyc3 outage, bucket deletion or compromise of the primary key does not
lose every backup.

## What is covered

One bucket-to-bucket `rclone copy` of the whole source bucket:

| Prefix in `onemind-backups` | Producer |
|---|---|
| `velero/cloud/` | Velero `default` BSL, schedule `velero-onemind-platform` 06:00 UTC, TTL 30d |
| `cnpg/onemind-data/` | CloudNativePG Barman base backups and WAL (retention 30d) |
| `backups/<name>/` | Encrypted logical archives and `.verified.json` markers (takserver-postgres, nats-onemind, postgres-{auth,n8n,matrix,takml,backstage}; `zeus-site` retired 2026-10-10, existing objects kept until cleanup, cluster-recovery, legacy clickhouse) |
| `backups/restic/<ns>/<pvc>/` | 46 VolSync ReplicationSources (Restic; 42 non-empty repositories) |

The mirror is written to `onemind-backups-sfo3/mirror/` with the same key layout.
Archives are already client-side encrypted (AES-256-GCM / Restic); the copy adds no
key material. Keep `BACKUP_ENCRYPTION_KEY` and the Restic password outside the
cluster as before: without them neither bucket is restorable.

Not covered: provider Droplet/volume snapshots, the laptop recovery mirror
(`tools/mirror.py`).

## Behaviour

- `rclone copy`, never `sync`: source deletions, Velero TTL expiry and Restic prune
  never delete anything in the mirror.
- Changed source objects overwrite the mirror; the previous mirror version is moved
  to `onemind-backups-sfo3/_overwritten/<run>/` (`--backup-dir`) rather than lost.
  `--immutable` was not used: Velero rewrites a few metadata objects after
  completion, which would fail the job on legitimate updates.
- Objects younger than 30 minutes are skipped so in-flight uploads are not copied
  half-written. Restic `locks/` are excluded.
- After copying, `rclone check --one-way --size-only` confirms every source object
  exists at the destination with equal size. Any difference fails the Job.
- A heartbeat object `_offsite/last-success.json` is written last.
- Schedule: daily 17:00 UTC. The last VolSync source starts 15:00 UTC and runs about
  a minute; TAK Postgres starts 02:00 and takes 1.5-2.5 h; NATS 03:30; cluster export
  04:00; Postgres exports 05:00-06:40; Velero 06:00. 17:00 leaves ~2 h slack after
  the latest writer. `activeDeadlineSeconds` 4 h; `backoffLimit` 1.
- Pod: non-root (65532), read-only root filesystem, all capabilities dropped,
  seccomp RuntimeDefault, 100m/256Mi requests, 1 CPU/1Gi limits, egress limited to
  DNS and HTTPS by `offsite-copy-egress`. Only `offsite-copy-credentials` is injected;
  the job never sees `BACKUP_ENCRYPTION_KEY`.
- Image: `rclone/rclone:1.75` pinned by multi-arch index digest (checked with
  `docker buildx imagetools inspect`; includes linux/amd64). It was not executed
  locally (no Docker daemon), so the first run is the functional test.

## Monitoring

- The CronJob is named `backup-offsite-copy`, so the existing `ScheduledBackupOverdue`
  alert (no successful run within 30 h, while not suspended) covers it unchanged.
  While it is suspended, `CriticalBackupScheduleSuspended` fires after 7 days instead.
- `OffsiteBackupCopyFailed` (added to `platform/monitoring/backup-health.yaml`) fires
  5 minutes after any failed Job; finished Jobs are garbage-collected after 24 h.
- A successful run proves presence and size, not restorability. Periodically run
  `restic check` against a mirrored repository and decrypt a mirrored archive with
  `tools/decrypt.py`.

## User steps

1. **Create the destination Space** in `sfo3`, named `onemind-backups-sfo3`, private,
   no CDN. `doctl` cannot create buckets (`doctl spaces` only manages keys), so use
   the control panel (Spaces Object Storage > Create) or an S3 client.
2. **Create a bucket-scoped Spaces key** (supported by `doctl spaces keys create`):
   ```sh
   doctl spaces keys create onemind-offsite-copy \
     --grants 'bucket=onemind-backups;permission=read,bucket=onemind-backups-sfo3;permission=readwrite'
   ```
   The access key ID and secret are shown once. Do not paste them into chat or Git;
   store them in a password manager or mode-0600 files. The key reads the source and
   writes the destination only. DO `readwrite` includes delete, so the provider does
   not enforce destination immutability; rely on `copy` semantics and keep this key
   only in the sealed secret.
3. **Optional, recommended**: enable versioning on the destination and add a
   lifecycle rule expiring `_overwritten/` after 90 days. This needs a temporary
   `fullaccess` key scoped to the destination bucket only (create, apply, then
   `doctl spaces keys delete`). Do not add age-based expiry to
   `mirror/backups/restic/`; Restic packs are long-lived regardless of age.
4. **Seal the credentials** (Secrets Convention in CLAUDE.md). Write each value to its
   own file with `umask 077` and `printf '%s' value > file` (no trailing newline):
   ```sh
   kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
     --controller-namespace=kube-system > /tmp/pub-cert.pem
   kubectl create secret generic offsite-copy-credentials -n ops \
     --from-file=RCLONE_CONFIG_SRC_ACCESS_KEY_ID=./src-id \
     --from-file=RCLONE_CONFIG_SRC_SECRET_ACCESS_KEY=./src-secret \
     --from-file=RCLONE_CONFIG_DST_ACCESS_KEY_ID=./dst-id \
     --from-file=RCLONE_CONFIG_DST_SECRET_ACCESS_KEY=./dst-secret \
     --dry-run=client -o json \
     | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/backups/sealed/ops-offsite-copy-credentials.yaml
   ```
   Then delete the plaintext files. With DigitalOcean at both ends SRC and DST hold the
   same key pair; they are separate so the destination can change provider.
5. **Enable it** in one commit: add `- sealed/ops-offsite-copy-credentials.yaml` to
   `platform/backups/kustomization.yaml` (deliberately not listed yet: a missing file
   would break the whole `platform-backups` sync) and set `suspend: false` on
   `backup-offsite-copy`.
6. **Seed and verify**: after Argo syncs, run once
   (`kubectl -n ops create job --from=cronjob/backup-offsite-copy offsite-seed-1`,
   then `kubectl -n ops logs -f job/offsite-seed-1`). The first run copies everything
   and may take an hour or more. Confirm `_offsite/last-success.json` exists, then
   spot-check `restic check` and `tools/decrypt.py` against mirrored objects.

## Retention policy

| Data | Source retention | Mirror retention |
|---|---|---|
| Velero backups | 30 d TTL | Kept until manually pruned; review quarterly |
| CNPG base backups/WAL | 30 d (Barman) | Same |
| Logical archives (`backups/<name>/`) | None (image never deletes) | Same as source |
| Restic repositories | 7 daily/4 weekly/3 monthly/1 yearly, prune every 14 d | Superset: old snapshots and packs stay; mirror only grows |
| `_overwritten/` | n/a | 90 d via lifecycle rule (step 3) |

Because `copy` never deletes, the mirror grows beyond the source. Reconcile
deliberately, never automatically: verify the source with `restic check`, then either
delete old mirror objects by hand or run a one-off `rclone sync` on a single prefix
after `rclone check`, with `--backup-dir` set. If age-based expiry is ever added to
the mirror, set `COPY_MAX_AGE` (for example `14d`) in `offsite-copy-config`;
otherwise `copy` re-uploads expired objects every night.

## Cost (assumptions; DigitalOcean Spaces pricing page, October 2026)

Spaces is one $5/month subscription per account: 250 GiB storage and 1,024 GiB
outbound transfer shared by all buckets; overage $0.02/GiB-month storage and
$0.01/GiB transfer; inbound free. The source bucket size cannot be read without its
credentials. Estimate from the 54.43 GiB local mirror of Restic repositories plus
archives, and Velero/CNPG data: roughly 60-100 GiB now.

| Scenario | Stored (primary + copy) | Incremental cost |
|---|---|---|
| Today, ~60-100 GiB each | 120-200 GiB, inside the 250 GiB allowance | $0 |
| Primary grows to 150 GiB | 300 GiB | ~$1/month |
| 500 GiB each | 1 TiB | ~$15/month |

nyc3 to sfo3 is a different region group, so the read counts against the 1 TiB
outbound allowance (about 100 GiB once to seed, then the daily change). Uploads to
sfo3 are free.

## Why sfo3

`sfo3` is the farthest standard DigitalOcean Spaces region from `nyc3` (opposite
coast, separate power and network region). `ams3`/`fra1` also work and suit EU
residency, but add transatlantic latency with no extra isolation benefit here.

## Moving to Backblaze B2 or a home server

Only `offsite-copy-config` and the credentials change:

- Backblaze B2 (S3 API): `RCLONE_CONFIG_DST_PROVIDER: Other`,
  `RCLONE_CONFIG_DST_ENDPOINT: https://s3.<region>.backblazeb2.com`,
  `RCLONE_CONFIG_DST_REGION: <region>`, `DST_BUCKET: <bucket>`; reseal the two DST
  values with a B2 application key.
- Home MinIO: `RCLONE_CONFIG_DST_PROVIDER: Minio` and its HTTPS URL (egress policy
  allows only 443).
- Home server pulling instead: run the same `rclone copy` there with the sfo3 bucket
  as source; `_offsite/last-success.json` shows freshness.
- A second copy job (duplicate the CronJob with its own DST config) can run to B2
  while sfo3 stays, giving three copies.

## Restoring from the mirror

Point the original tooling at `onemind-backups-sfo3` with prefix `mirror/`. Restic:
`s3:https://sfo3.digitaloceanspaces.com/onemind-backups-sfo3/mirror/backups/restic/<ns>/<pvc>`
(locks are not copied, so no `restic unlock` is needed). Velero: create a BSL with
`bucket: onemind-backups-sfo3`, `prefix: mirror/velero/cloud`, `region: sfo3`,
`s3Url: https://sfo3.digitaloceanspaces.com`, `accessMode: ReadOnly`.

## Verification reliability — 2026-10-26

The job selects eligible object paths once and uses that same literal manifest
for both copy and verification. Re-evaluating a relative age cutoff during a long
copy could previously select newer objects during verification that were absent
from the copy selection. The job now also compares the complete selected path
list against rclone's verified-match report; disappeared source objects cannot
produce a false successful receipt. Empty selections and listing/copy/check
failures preserve the previous success receipt. Source deletions still never
propagate, and replaced destination objects still go to `_overwritten/<run>`.

Nine local-remotes regression cases run the real pinned rclone 1.75.2 engine in CI.
Verification remains size-only for each selected object; an independent restore
and authenticated archive-decryption test is required after initial seeding.
A successful size check is not a complete recovery drill.

Current deployment remains suspended: no independent destination credentials are
configured. The source's existing SFO3 example is regional redundancy within
DigitalOcean, not protection against provider/account loss. Do not enable the job
until the owner selects a destination, scoped credentials are sealed, initial
copy passes, and recovery from that destination is verified.
