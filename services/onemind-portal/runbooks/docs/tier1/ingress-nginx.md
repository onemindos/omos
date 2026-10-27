# Ingress (ingress-nginx)

The public HTTP and HTTPS entry point: ingress class `nginx` behind the DigitalOcean load balancer `143.244.220.155`. Public DNS names are Cloudflare-proxied A records pointing at that address, and public Ingresses carry `nginx.ingress.kubernetes.io/whitelist-source-range` with the Cloudflare ranges; the Matrix ingresses are the known exception without a source restriction. Admin UIs do not use it: they use the `tailscale` ingress class and are tailnet-only. TAK has no public listener.

Sources: root `CLAUDE.md` ("Public ingress"), `platform/values/ingress-nginx.yaml`, `platform/argocd/edge-ingress-nginx.yaml`, `docs/gitops/HELM-ADOPTION.md`, `docs/gitops/APP-FINALIZERS.md`, `docs/networking/NETWORKING.md`.

## Ownership

- `edge-ingress-nginx` adopted the hand-installed Helm release in PR #176: chart `ingress-nginx` 4.15.1, release name `ingress-nginx`, values in `platform/values/ingress-nginx.yaml`. Never run `helm upgrade` on it; change the values file and merge.
- Prune is off, and the Application deliberately has no finalizer: a cascading delete would release the load balancer and its public IP.
- Two controller replicas, spread across nodes where possible. The Service uses `externalTrafficPolicy: Local`, so the load balancer only sends traffic to nodes running a controller pod, and the client address reaches the allowlist check. With one replica that would be a single point of failure.
- The admission webhook certificate comes from the chart's hook Jobs, which Argo CD runs as short-lived Jobs during a sync; Argo CD ignores the injected `caBundle`.

## Investigate (read-only)

```sh
kubectl --context do-nyc3-omos-cloud -n argocd get application edge-ingress-nginx
kubectl --context do-nyc3-omos-cloud -n ingress-nginx get deployments,pods,services,poddisruptionbudgets
kubectl --context do-nyc3-omos-cloud get ingress -A
kubectl --context do-nyc3-omos-cloud -n ingress-nginx logs deployment/ingress-nginx-controller --tail=100
```

If every Cloudflare-only site answers 403, check that the Service still uses `externalTrafficPolicy: Local`: with `Cluster` the node replaces the client address with its own and the allowlist rejects every request.

## Recovery

Revert the offending values or Ingress change through Git and let Argo CD sync it. The controller holds no state; Git is its backup. Do not delete the LoadBalancer Service: that would release the address DNS points at.

## Verify completion

Both controller pods are Ready, `edge-ingress-nginx` is Synced and Healthy, and each affected public hostname answers through Cloudflare.

Generated workload list: [ingress-nginx app runbook](../apps/ingress-nginx.md).
