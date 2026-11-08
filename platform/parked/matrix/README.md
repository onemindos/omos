# Matrix (Synapse + Element Web) — parked 2026-11-07

Retired by the owner on 2026-11-07: 2 users, 0 rooms, no events. Everything that ran in namespace `matrix`
(Synapse, Element Web x3, PostgreSQL, Redis) was removed from the cluster. These files are the old
`platform/matrix/` tree, moved here unchanged. `platform-matrix.application.yaml` is the Argo Application that
used to live in `platform/argocd/`. Nothing in this folder is deployed.

## What was kept

All objects are in the `onemind-backups` Spaces bucket (nyc3, `https://nyc3.digitaloceanspaces.com`).

| What | Where | Notes |
|------|-------|-------|
| Final PostgreSQL dump of the `synapse` database | `backups/matrix/final-2026-11-07/20261108T045506Z-4361f023d03c4326a391ec46a5d04e74/archive.tar.gz.aesgcm` (39,422 bytes) plus `.verified.json` marker | Made with the `platform-backup` image (same format as the nightly `backups/postgres-matrix/` archives). Encrypted with `BACKUP_ENCRYPTION_KEY` from `platform-backup-credentials`; the restore into a scratch Postgres was checked (151 tables) |
| Synapse volume files (`signing.key`, `log.config`; `media/` was empty) | `backups/matrix/final-2026-11-07/synapse-data.tar.gz.aesgcm` (1,770 bytes) | Same AES-GCM format and key. Tar of `/data` as `synapse-data/` |
| VolSync restic repositories (daily until retirement) | `backups/restic/matrix/data-postgres-0`, `backups/restic/matrix/data-redis-0`, `backups/restic/matrix/synapse-data` | Repository passwords were in the removed `volume-*-restic` SealedSecrets; they are in git history and in the `backup-cluster-recovery` export. Nothing prunes these repos |
| Nightly postgres archives | `backups/postgres-matrix/` | Stopped when the CronJob was removed; existing objects stay |

The Synapse config is the SealedSecret `argo/synapse-homeserver.sealedsecret.yaml`. It is bound to namespace
`matrix` and the cluster's sealed-secrets key, so it still unseals on this cluster. `matrix-secrets`,
`platform-backup-credentials`, `ghcr` and the VolSync restic secrets were deleted with the namespace and must be
resealed on restore (see `SECRETS-TODO.md`).

## Restore

1. Move this folder back to `platform/matrix/` and the Application back to `platform/argocd/platform-matrix.yaml`.
2. Recreate namespace `matrix` (it was labelled `onemind.io/layer: act` with pod-security warn/audit `restricted`; the
   old definition is in git history of `platform/security/shared-namespaces.yaml`), and reseal `matrix-secrets` and
   the credentials (`SECRETS-TODO.md`). Put back the `backup-postgres-matrix` CronJob and VolSync sources from git
   history of `platform/backups/` if backups are wanted again.
3. Sync `platform-matrix` by hand (it is manual-sync on purpose, see the Application's header comment).
4. Restore the database before starting Synapse: decrypt the dump with the `platform-backup` image helpers
   (`backup.decrypt`), untar, and load the `pg_dump` into the `synapse` database. Or restore the VolSync
   repositories into new PVCs (see `docs/operations/parked-services.md`).
5. Put `signing.key` back into the Synapse data volume so the server keeps its identity.

## Public hostnames

`matrix.onemindos.com` and `element.onemindos.com` were already gone from the ingress (2026-10-29). The tailnet hosts
`matrix.tail717ba0.ts.net` and `element.tail717ba0.ts.net` were removed with the workloads. The Cloudflare DNS records
for the two `onemindos.com` hosts (and `onemindos.dev` equivalents, if any exist) remain for the owner to delete.

## Known leftovers

The tool gateway execute-gate (`platform/tool-gateway/manifests/execute-gate/`, replicas 0) was designed to approve
calls through a Matrix room and its NetworkPolicy still allows egress to namespace `matrix`. It cannot work without
Matrix; decide separately whether to retire or rework it.
