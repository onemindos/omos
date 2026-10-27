# cert-manager

Issues and renews in-cluster certificates. There is no `ClusterIssuer`: the unused `letsencrypt-dns` issuer and its Cloudflare token were removed. The only issuers are namespaced self-signed CAs owned by the components that use them: NATS server TLS, the CloudNativePG Barman plugin, and the Hermes Paperclip transport added in October 2026. A read-only check on October 27, 2026 found no certificate for any ingress hostname, so ingress HTTPS, public or tailnet, does not depend on cert-manager.

Sources: `docs/STACK.md` (cert-manager row), `platform/values/cert-manager.yaml`, `platform/argocd/edge-cert-manager.yaml`, `docs/gitops/HELM-ADOPTION.md`, `docs/gitops/APP-FINALIZERS.md`.

## Ownership

- `edge-cert-manager` adopted the Helm release in PR #176: chart `cert-manager` v1.20.3, values in `platform/values/cert-manager.yaml`. Never run `helm upgrade` on it; change the values file and merge.
- The six CRDs render from chart templates (`crds.enabled: true`, each with `helm.sh/resource-policy: keep`), are applied server-side and are never pruned or replaced. The Application has no finalizer, because deleting a CRD deletes every resource of that kind in the cluster.
- cainjector fills the webhook `caBundle` fields, which Argo CD ignores.
- Each Issuer and Certificate belongs to its component's manifests, for example `platform/nats/certs.yaml`. Change them there.

## Investigate (read-only)

```sh
kubectl --context do-nyc3-omos-cloud -n argocd get application edge-cert-manager
kubectl --context do-nyc3-omos-cloud -n cert-manager get deployments,pods
kubectl --context do-nyc3-omos-cloud get clusterissuers,issuers,certificates -A
kubectl --context do-nyc3-omos-cloud -n cert-manager logs deployment/cert-manager --tail=100
```

A Certificate's Ready condition, expiry and renewal time show whether renewal is on track; `kubectl describe certificate` shows the last issuance events.

## Recovery

cert-manager stores no application data. Certificate and Secret objects are included in the encrypted `backup-cluster-recovery` export. Fix a broken issuer or certificate through its component's Git change; do not disable TLS verification in clients as a workaround.

## Verify completion

The controller, cainjector and webhook are Ready, `edge-cert-manager` is Synced and Healthy, the affected Certificate is Ready with a future renewal time, and the consuming service accepts verified TLS connections.

Generated workload list: [cert-manager app runbook](../apps/cert-manager.md).
