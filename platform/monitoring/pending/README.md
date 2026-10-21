# Pending Alertmanager receiver: OneMind portal

`alertmanagerconfig-onemind-backstage.yaml.disabled` sends critical and warning alerts to
the portal (`POST /api/events/http/alertmanager`). It is **not deployed**:

- The `.disabled` suffix means no Argo CD source or kustomization includes it. Argo's
  `platform-monitoring` Application reads `platform/monitoring` recursively and only treats
  `.yaml`, `.yml`, `.json` and `.jsonnet` files as manifests.
- It depends on a secret that does not exist yet. Do not enable it before that secret exists.
- `scripts/tests/test_gitops_alert_routing.py` fails if the file is renamed to a manifest
  extension, listed in a kustomization, or if an equivalent AlertmanagerConfig appears under
  `platform/monitoring` while the secret manifests are missing.

Alertmanager is managed by the kube-prometheus-stack release (Argo Application
`monitoring-kube-prometheus-stack`). Its `alertmanagerConfigSelector: {}` makes the operator
merge every `AlertmanagerConfig` in the `monitoring` namespace into the live config, with no
namespace matcher injected (`alertmanagerConfigMatcherStrategy: None`). The only other one
today is `openchoreo`. A new config with `continue: true` adds a route and leaves existing
ones untouched.

## Secret (Zeus creates it; the same random value goes in two namespaces)

| Namespace | Secret name | Key | Used by |
|---|---|---|---|
| `monitoring` | `onemind-alert-webhook` | `token` | this AlertmanagerConfig (`httpConfig.authorization`) |
| `backstage` | `onemind-alert-webhook` | `ONEMIND_ALERT_WEBHOOK_TOKEN` | portal env var, optional |

```bash
kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > /tmp/pub-cert.pem
TOKEN="$(openssl rand -hex 32)"   # 64 characters; keep it in the shell, never in a file or chat
kubectl create secret generic onemind-alert-webhook -n monitoring --from-literal=token="$TOKEN" \
  --dry-run=client -o json | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/monitoring/onemind-alert-webhook.sealedsecret.yaml
kubectl create secret generic onemind-alert-webhook -n backstage \
  --from-literal=ONEMIND_ALERT_WEBHOOK_TOKEN="$TOKEN" \
  --dry-run=client -o json | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/backstage/alert-webhook.sealedsecret.yaml
unset TOKEN
```

## Enable (one PR, after the secret files exist)

1. The portal image that contains F3 is released and pinned (the route does not exist before).
2. Commit the two SealedSecret files above and add `alert-webhook.sealedsecret.yaml` to
   `platform/backstage/kustomization.yaml`.
3. In `platform/backstage/deployment.yaml` add under the container's `envFrom` (keep it
   optional so the pod still starts if the secret is ever removed):

   ```yaml
   - secretRef:
       name: onemind-alert-webhook
       optional: true
   ```

4. `git mv platform/monitoring/pending/alertmanagerconfig-onemind-backstage.yaml.disabled platform/monitoring/alertmanagerconfig-onemind-backstage.yaml`
5. After Argo syncs: check the new route in Alertmanager (`/#/status`), then fire a harmless
   test alert and confirm one inbox item and its Resolved update. Rollback is removing step 4's
   file; the portal route then simply receives nothing.

If the secret in `monitoring` is missing when the AlertmanagerConfig is applied, the operator
skips that one config and logs the reason; it does not stop the other receivers.
