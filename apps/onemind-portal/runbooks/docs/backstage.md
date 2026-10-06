# Backstage and catalog

## Access and dependencies

Open `https://onemind-preview.tail717ba0.ts.net/`. Backstage runs in namespace `backstage`, deployment `backstage`, with PostgreSQL in `backstage-postgres-0`. Private routing comes from Tailscale; login comes from Authentik. Source and image digest are in `apps/onemind-portal/` and `platform/backstage/`.

## Investigate

```sh
kubectl config current-context
kubectl -n argocd get application platform-backstage
kubectl -n backstage get deployments,pods,events
kubectl -n backstage logs deployment/backstage --tail=100
```

If the pod is waiting to mount a catalog ConfigMap, compare the referenced ConfigMap name with available resources. Catalog generation uses Argo server-side apply to avoid the last-applied annotation size limit. Do not disable validation or remove catalog data to mask a failed apply.

If new inventory is missing, confirm Argo synced the intended revision, the pod mounts the new catalog, and catalog processing has had time to refresh stored locations. Inspect processing errors in the catalog and logs. The reviewed enrichment command updates both catalog copies; plain cluster snapshot generation omits curated enrichment.

## Updates and rollback

The application currently has one replica and uses Recreate. A release has a short availability gap. Require a successful image build, pin the digest in Git, and let Argo deploy the merged change. To roll back, create a reviewed PR restoring the previous verified image digest and compatible configuration. Check migration compatibility before reverting application versions. Do not assume a rollback reverses database migrations.

## Backups and recovery

`backstage-backup` runs daily at 07:15 UTC. Each successful run dumps application databases, encrypts the archive, uploads to private Spaces storage, and restores the catalog to a temporary database to verify recovery. Inspect job status and the restore verification output; never print the backup Secret.

```sh
kubectl -n backstage get cronjob backstage-backup
kubectl -n backstage get jobs --sort-by=.metadata.creationTimestamp
```

Before a recovery exercise, confirm availability of database credentials, the backup encryption key, Spaces access and the SealedSecrets sealing key. Use a separate PostgreSQL instance and original database names. Stop application writes before final cutover. Follow `platform/backstage/README.md`; do not restore over a running database. Keep the original volume and keys until recovered data and login are verified.

## Verify completion

Argo must be Synced and Healthy. Confirm signed-in home, catalog, documentation, and operations panels render. Confirm anonymous operations API access returns 401. Check that the newest expected catalog entries and backup evidence are present.
