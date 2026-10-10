# Secrets — External Secrets Operator + 1Password

> **Historical (checked 2026-10-09).** External Secrets Operator is not installed in the cluster (no `external-secrets` namespace). Secrets are committed as Bitnami SealedSecrets; the current procedure is in the repo-root `CLAUDE.md`, section "Secrets". Keep this page only as design history.

All K8s secrets in this cluster are managed declaratively via **External Secrets Operator (ESO) v2.7.0** + **1Password Connect Server**. No `kubectl create secret` for service secrets — everything syncs automatically from 1Password.

## How It Works

```
1Password vault          ESO ClusterSecretStore       K8s Secret
"00-24 SO (Sovereign     ←──── onepassword-connect ────→  (auto-created,
  Operations)"                  pod (in-cluster)            kept in sync)
```

1. 1Password Connect pod runs in `external-secrets` namespace, syncs vault data into the cluster.
2. `ClusterSecretStore` named `onepassword` points at the Connect pod (`platform/external-secrets/cluster-secret-store.yaml`).
3. Each service has an `external-secret.yaml` that declares which 1Password item maps to which K8s secret.
4. ESO reconciles every hour (and on demand via annotation).

## 1Password Item Standard

| Rule | Value |
|------|-------|
| Vault | `00-24 SO (Sovereign Operations)` |
| Item title | Exact K8s secret name (e.g. `cloudtak-secrets`) |
| Fields | Flat — no sections. Field label = K8s secret key |
| Tags | `k8s-secret`, `omos-hub` |

**One item per K8s secret.** Sections are not supported by the ESO Connect provider — it matches field labels only, and duplicate labels across sections cause `got N` errors.

## ExternalSecret Files

| File | K8s Secret(s) | Namespace |
|------|--------------|-----------|
| `platform/legacy/external-secret.yaml` | hermes-env, hermes-leaf-conf, hermes-nats-ca | hermes |
| `platform/hermes-webui/external-secret.yaml` | hermes-webui-auth | hermes |
| `platform/cloudtak/external-secret.yaml` | cloudtak-secrets | cloudtak |
| `platform/wwv/external-secret.yaml` | wwv-env, wwv-pg-credentials | wwv |
| `platform/omos/external-secret.yaml` | omos-nats-creds | omos |
| `platform/orp/external-secret.yaml` | orp-nats-creds, orp-bridge-key | orp |
| `platform/cloudflare/external-secret.yaml` | cloudflare-api-credentials | cloudflare |
| `platform/ingress/external-secret.yaml` | cloudflare-api-token | cert-manager |
| `tak-cluster/deploy/external-secret.yaml` | tak-nats-leaf | takserver |

## Secrets NOT Managed by ESO

Some secrets cannot be automated and must be created once out-of-band:

| Secret | Namespace | Why not in ESO |
|--------|-----------|----------------|
| `tunnel-creds-json` | cloudflare | Derived by decoding TUNNEL_TOKEN — not stored as-is in 1Password |
| `nats/values.secret.yaml` | nats | Helm values file injected at install time, not a K8s secret |
| `ghcr` / `ghcr-pull-secret` | hermes/orp/omos/wwv | Docker registry credentials |
| `cloudtak-agent-token` | hermes | Must be generated in CloudTAK UI → Settings → Tokens |
| `op-credentials` | external-secrets | Bootstrap: 1password-credentials.json for Connect pod itself |
| `onepassword-connect-token` | external-secrets | Bootstrap: Connect JWT token |

## Adding a New Secret

1. Add an item to 1Password vault `00-24 SO (Sovereign Operations)`
   - Title = the K8s secret name you want
   - Add one field per key, flat (no sections)
   - Tag it `k8s-secret` and `omos-hub`

2. Write `platform/{service}/external-secret.yaml`:

```yaml
apiVersion: external-secrets.io/v1
kind: ExternalSecret
metadata:
  name: my-secret
  namespace: my-namespace
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: onepassword
    kind: ClusterSecretStore
  target:
    name: my-secret
    creationPolicy: Owner
  dataFrom:
    - extract:
        key: my-secret   # 1Password item title
```

   Use `dataFrom.extract` when all field labels already match the K8s key names you want.
   Use `data[].remoteRef` when you need to rename keys:

```yaml
  data:
    - secretKey: user          # K8s key name
      remoteRef:
        key: my-secret         # 1Password item title
        property: username     # 1Password field label
```

3. Apply: `kubectl apply -f platform/{service}/external-secret.yaml`

4. Verify: `kubectl get externalsecret my-secret -n my-namespace`
   — status should show `SecretSynced: True` within ~10 seconds.

## Fresh Cluster Bootstrap

ESO requires two secrets before it can sync anything else (create these once, never committed):

```bash
# 1. Get the 1password-credentials.json from 1Password Connect server setup
kubectl create namespace external-secrets
kubectl create secret generic op-credentials \
  -n external-secrets \
  --from-file=1password-credentials.json=/path/to/1password-credentials.json

# 2. Get the Connect JWT token from 1Password
kubectl create secret generic onepassword-connect-token \
  -n external-secrets \
  --from-literal=token=<connect-jwt-token>
```

Then run `platform/install.sh` — step 0 installs ESO + Connect via Helm and applies all ExternalSecrets.

## Ops Runbook

**Force a secret to re-sync immediately:**
```bash
kubectl annotate externalsecret <name> -n <namespace> \
  force-sync="$(date +%s)" --overwrite
```

**Check sync status across all namespaces:**
```bash
kubectl get externalsecrets -A
```

**Connect pod updated a stale cache (after editing/recreating a 1Password item):**
```bash
kubectl rollout restart deployment/onepassword-connect -n external-secrets
# Wait for rollout, then force-sync the affected ExternalSecrets
```

**Check why a secret failed to sync:**
```bash
kubectl describe externalsecret <name> -n <namespace>
# Look at Events section for the error message
```
