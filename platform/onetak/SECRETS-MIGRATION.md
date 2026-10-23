# Secrets Migration — cloudtak → onetak

SealedSecrets are namespace-scoped. The encrypted blobs from `cloudtak` won't
decrypt in `onetak`. Run the following to copy the live secrets directly from
the running cloudtak namespace into onetak, then seal them.

The `onetak-secrets` secret, the MinIO object store, the `onetak-postgres-creds` secret
and the `onetak-postgis` database were removed from `onetak`; only the GHCR pull secret
below still applies.

## Step 1 — Copy live secrets into onetak namespace

```bash
# Copy GHCR pull secret
kubectl get secret ghcr -n cloudtak -o json \
  | jq 'del(.metadata.resourceVersion,.metadata.uid,.metadata.creationTimestamp,.metadata.annotations,.metadata.labels) | .metadata.namespace = "onetak"' \
  | kubectl apply -f -
```

Then seal it (`kubeseal --namespace onetak --format yaml`) and keep the SealedSecret in GitOps.
