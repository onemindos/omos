# Matrix secrets: what you must do before syncing `platform-matrix`

> Current review: [Matrix drift and safe sync](../../docs/operations/matrix-drift-2026-10-23.md).
> Keep manual sync. Create and verify the sealed config before selecting Synapse;
> never force, replace, prune or recreate the data PVC. The current homeserver is
> `matrix.tail717ba0.ts.net` (tailnet-only since 2026-10-29); the active server identity is `onemindos.dev`.

`platform-matrix` (path `platform/matrix/argo/`) is deliberately **manual sync**. Reason: the live Synapse config
(`ConfigMap/synapse-config`, key `homeserver.yaml`) contains three plaintext secrets, so it cannot be committed:

| Value | Where it is in `homeserver.yaml` |
|-------|----------------------------------|
| PostgreSQL password | `database.args.password` |
| Redis password | `redis.password` |
| `registration_shared_secret` | top level |

The Argo-managed Deployment therefore mounts `/config/homeserver.yaml` from **`Secret/synapse-homeserver`** (key
`homeserver.yaml`) and `/config/log.config` from `ConfigMap/synapse-log-config`. The file Synapse reads is byte-for-byte
the one it reads today, so sealing the current live file means **no functional change**. The Secret does not exist yet:
syncing before step 1 would leave the new Synapse pod stuck in `ContainerCreating` (the old pod keeps serving).

Synapse has no env substitution in `homeserver.yaml`, and it has no `_path` option for the database or (as far as
verified for v1.106) Redis passwords, so moving only `registration_shared_secret` to `registration_shared_secret_path`
would still leave two plaintext passwords. Sealing the whole file is the lowest-risk option.

Other Secrets (`matrix-secrets`, `matrix-tls`, `ghcr`) are not touched by Argo. The files `base/manifests/secret.yaml`,
`base/redis/manifests/secret.yaml` and `overlays/` contain only placeholders and are not in the Argo path.

## 1. Seal the CURRENT live config (no rotation)

Requires `kubectl` on context `do-nyc3-omos-cloud` and `kubeseal`. The plaintext only travels through the pipe.

```bash
kubectl --context do-nyc3-omos-cloud -n matrix get configmap synapse-config \
  -o go-template='{{index .data "homeserver.yaml"}}' \
| kubectl create secret generic synapse-homeserver --namespace matrix \
    --from-file=homeserver.yaml=/dev/stdin --dry-run=client -o yaml \
| kubeseal --context do-nyc3-omos-cloud \
    --controller-name sealed-secrets-controller --controller-namespace kube-system \
    --format yaml > platform/matrix/argo/synapse-homeserver.sealedsecret.yaml
```

The SealedSecret uses the default strict scope (bound to name `synapse-homeserver` and namespace `matrix`). Commit only
the `.sealedsecret.yaml` output, never the plaintext. Confirm with `grep -c encryptedData` and run gitleaks.

## 2. Sync

1. Merge the sealed file in a separately reviewed PR, then selectively sync only that SealedSecret in Argo. The controller creates `Secret/synapse-homeserver`; no full application sync yet.
2. Check: `kubectl -n matrix get secret synapse-homeserver`.
3. In Argo, review the diff of `platform-matrix`, then sync manually. Expect one Synapse rollout (pod template volume
   changed), nothing else. Verify `https://matrix.tail717ba0.ts.net/_matrix/client/versions` and a login.
4. Keep `platform-matrix` manual. Automated sync is a separately reviewed decision.
5. Only after Synapse is healthy: remove the old `ConfigMap/synapse-config` (it still holds the plaintext) with
   `kubectl -n matrix delete configmap synapse-config`. Argo does not manage or prune it.

## 3. Rotating instead (recommended, the live values are weak and have appeared in `kubectl` annotations)

`Secret/matrix-secrets` was created with `kubectl apply`, so its `last-applied-configuration` annotation holds the
values in clear text. Anyone with `get secret` rights can read them. Rotate all four keys
(`postgres-password`, `redis-password`, `synapse-registration-secret`, `mas-postgres-password`):

```bash
# Generate new values locally (not printed to the repo):
NEWPG=$(openssl rand -base64 32 | tr -d '/+=\n'); NEWREDIS=$(openssl rand -base64 32 | tr -d '/+=\n'); NEWREG=$(openssl rand -base64 48 | tr -d '/+=\n')

# 1. Change the PostgreSQL role password inside the database:
kubectl --context do-nyc3-omos-cloud -n matrix exec postgres-0 -- psql -U synapse -d synapse -c "ALTER USER synapse PASSWORD '$NEWPG'"

# 2. Seal matrix-secrets (keep mas-postgres-password's current value or rotate it the same way):
kubectl create secret generic matrix-secrets --namespace matrix --dry-run=client -o yaml \
  --from-literal=postgres-password="$NEWPG" \
  --from-literal=redis-password="$NEWREDIS" \
  --from-literal=synapse-registration-secret="$NEWREG" \
  --from-literal=mas-postgres-password="<current or new value>" \
| kubeseal --controller-name sealed-secrets-controller --controller-namespace kube-system --format yaml \
  > platform/matrix/argo/matrix-secrets.sealedsecret.yaml

# 3. Build the new homeserver.yaml from the live one with the three values replaced, seal it as in step 1.
```

Then restart in this order: Redis (picks up `--requirepass` from `matrix-secrets`), PostgreSQL env is only read at first
init (hence the `ALTER USER`), then Synapse. Adopting an existing live `Secret/matrix-secrets` with a SealedSecret
requires annotating the live Secret with `sealedsecrets.bitnami.com/managed="true"` first (otherwise the controller
refuses to overwrite it); do this in a maintenance window. Also note that `platform/matrix/homeserver.yaml` in Git history contained
the old plaintext values (now redacted in the working tree); treat those values as burned and rotate them.
