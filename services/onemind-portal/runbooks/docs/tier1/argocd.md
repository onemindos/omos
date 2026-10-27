# Argo CD

Argo CD applies everything under `platform/` from `main`. `app-of-apps` creates every Application in `platform/argocd/`, and the ecosystem parent Applications create the rest. The UI is tailnet-only at `https://argocd.tail717ba0.ts.net`; sign-in goes through Authentik, where group `argocd-admins` is admin, and the local `admin` account is kept as break-glass. Running workloads do not depend on Argo CD being up.

Sources: `docs/gitops/ARGOCD-SELF-MANAGEMENT.md`, `docs/gitops/HELM-ADOPTION.md`, `docs/gitops/APP-FINALIZERS.md`, `docs/operations/cluster-app-recovery-2026-10-25.md`, `platform/values/argocd.yaml`, `platform/backups/image/README.md`.

## Ownership

- Argo CD manages itself through `platform-argocd`: chart `argo-cd` at the version pinned in `platform/argocd/platform-argocd.yaml`, release `argocd`, values in `platform/values/argocd.yaml`. The adoption sync on October 25, 2026 changed only tracking annotations and restarted no pods. The Application is on manual sync, without prune, `Replace` or a finalizer.
- The `applications`, `applicationsets` and `appprojects` CRDs are outside its sync scope (`crds.install=false`). On a chart bump they are updated by hand first; the procedure is in `ARGOCD-SELF-MANAGEMENT.md`.
- Never run `helm upgrade`, `helm rollback` or `helm uninstall argocd` for a normal change: uninstall deletes every object in the release record, including Argo CD itself. Never run `platform/install.sh` on this cluster.
- `app-of-apps` prunes. Removing an Application file deletes the live Application. If that Application carries `resources-finalizer.argocd.argoproj.io`, everything it tracks is deleted with it, and a tracked Namespace takes every object in it, including objects that other Applications own. Before removing or renaming an Application, check the live object with the two finalizer commands under "Investigate". Do not rely on `docs/gitops/APP-FINALIZERS.md` for this, because its list is incomplete. On October 27, 2026, 25 live Applications carried the finalizer and 18 of them were missing from it, among them the tier 1 `platform-takserver`, `platform-nats`, `platform-core`, `platform-auth-helm`, `platform-auth-infra`, `platform-ops` and `platform-legacy`.

## Investigate (read-only)

```sh
kubectl --context do-nyc3-omos-cloud -n argocd get applications
kubectl --context do-nyc3-omos-cloud -n argocd get deployments,statefulsets,pods
kubectl --context do-nyc3-omos-cloud -n argocd logs statefulset/argocd-application-controller --tail=100
# Before removing or renaming an Application: does deleting it cascade, and does it track a Namespace or PVC?
kubectl --context do-nyc3-omos-cloud -n argocd get application <name> -o jsonpath='{.metadata.finalizers}'
kubectl --context do-nyc3-omos-cloud -n argocd get application <name> -o jsonpath='{range .status.resources[*]}{.kind}{" "}{.namespace}{" "}{.name}{"\n"}{end}' | grep -E '^(Namespace|PersistentVolumeClaim) '
```

OutOfSync means live state differs from Git: read the diff before syncing anything. Before a manual sync of `platform-argocd`, `argocd app diff platform-argocd` must show no unintended change to a Deployment, StatefulSet, Service, ConfigMap, Role or ClusterRole; an unintended pod-template change means fix Git or the values, not the cluster.

## Break-glass

If Argo CD is down or mis-synced, `ARGOCD-SELF-MANAGEMENT.md` documents rolling back or restarting the broken Deployment or StatefulSet in `argocd` and signing in with the local `admin` account. Its last resort, when Argo CD cannot run at all, is re-running `helm upgrade` with the last good values from Git history, which returns Argo CD to Helm-owned state.

## Backups

Applications and settings are in Git, and `app-of-apps` recreates the Applications from it. The encrypted `backup-cluster-recovery` export holds only Secrets and ConfigMaps. Among them are `argocd-cm` and `argocd-rbac-cm`, which Git also renders, and `argocd-secret` (server signing key, admin password hash), which Argo CD writes itself and Git does not render. No backup holds the Argo CD CRDs or the Application objects: the CRDs come from the pinned chart through the `helm template ... --set crds.install=true | kubectl apply --server-side` step in `ARGOCD-SELF-MANAGEMENT.md`.

## Verify completion

The Argo CD Deployments and the application controller are Ready, the affected Applications are Synced and Healthy, SSO sign-in works in a private window, and the `admin` break-glass still works.

Generated workload list: [ArgoCD app runbook](../apps/argocd.md).
