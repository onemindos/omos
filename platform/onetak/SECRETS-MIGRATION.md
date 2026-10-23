# Secrets Migration — cloudtak → onetak

SealedSecrets are namespace-scoped. The encrypted blobs from `cloudtak` won't
decrypt in `onetak`. Run the following to copy the live secrets directly from
the running cloudtak namespace into onetak, then seal them.

## Step 1 — Copy live secrets into onetak namespace

```bash
# Copy postgres credentials
kubectl get secret cloudtak-postgres-creds -n cloudtak -o json \
  | jq 'del(.metadata.resourceVersion,.metadata.uid,.metadata.creationTimestamp,.metadata.annotations,.metadata.labels) | .metadata.namespace = "onetak" | .metadata.name = "onetak-postgres-creds"' \
  | kubectl apply -f -

# Copy GHCR pull secret
kubectl get secret ghcr -n cloudtak -o json \
  | jq 'del(.metadata.resourceVersion,.metadata.uid,.metadata.creationTimestamp,.metadata.annotations,.metadata.labels) | .metadata.namespace = "onetak"' \
  | kubectl apply -f -
```

## Step 2 — Re-seal for future GitOps (optional but recommended)

Once the namespace is stable, re-seal the secrets properly:

```bash
# Export and re-seal onetak-postgres-creds
kubectl get secret onetak-postgres-creds -n onetak -o json \
  | kubeseal --namespace onetak --format yaml > onetak-postgres-creds.sealedsecret.yaml
```

Then replace the Step 1 manual copies with these SealedSecret files in GitOps.

## Step 3 — Data migration (PostGIS)

The `onetak-secrets` secret and the MinIO object store were removed from `onetak` (the
`cloudtak` bucket held no objects); that part of this migration no longer applies.

For PostGIS, the onetak-postgis starts fresh. To migrate existing data:

```bash
# Dump from cloudtak
kubectl exec -n cloudtak deploy/cloudtak-postgis -- \
  pg_dump -U docker gis > cloudtak-dump.sql

# Restore into onetak (once onetak-postgis is running)
kubectl exec -i -n onetak deploy/onetak-postgis -- \
  psql -U docker gis < cloudtak-dump.sql
```
