# Argo CD Application finalizers

`resources-finalizer.argocd.argoproj.io` on an Application makes Argo CD delete every
resource the Application tracks before the Application object itself goes away
(a "cascading delete"). Without it, deleting an Application only removes the Application
object and leaves its workloads running, unmanaged (orphaned).

This is the policy for which child Applications carry it. Last reviewed 2026-10-23
against live `do-nyc3-omos-cloud` (read-only).

## What the finalizer does here

- `app-of-apps` prunes (`prune: true`): removing an Application file from `platform/argocd/`
  deletes the live Application. With the finalizer that also deletes its workloads; without
  it the workloads are orphaned.
- The three ecosystem parents (`platform-ecosystem-applications`, `-foundation`,
  `-resources`) and `platform-security` have `prune: false`. Removing a child Application
  file there deletes nothing; the Application must be deleted by hand, and the finalizer
  decides whether that cascades.
- `app-of-apps` itself keeps its finalizer and is not touched.

## Has the finalizer (stateless, namespace-scoped)

| Application | Namespace | Tracks |
|-------------|-----------|--------|
| `ecosystem-falco` | falco | DaemonSet, ConfigMaps, ServiceAccount |
| `ecosystem-litellm` | onemind-ai | Deployment, Service, ConfigMap, ServiceMonitor |
| `ecosystem-opencost` | opencost | Deployment, Service, ClusterRole/Binding |
| `ecosystem-openfga` | onemind-core | Deployment, Service, Role |
| `ecosystem-onemind-core` | onemind-core | Deployments, Dapr Component, RBAC |
| `platform-nats-tools` | nats | UI Deployments, Services, Ingresses (not the JetStream cluster) |
| `platform-openchoreo-events`, `-finops`, `-logs`, `-traces` | openchoreo-observability-plane | adapters and collectors only |
| `ecosystem-opensearch` | onemind-search | StatefulSet; its PVC is a volumeClaimTemplate (not tracked, Retain) |
| `ecosystem-oc-control` | openchoreo-control-plane | OpenChoreo control plane **and 36 CRDs** (see below) |
| `ecosystem-harbor` | onemind-registry | being removed; tracks PVCs `harbor-registry` (Retain) and `harbor-jobservice` |

Cascade only deletes what the Application tracks. It never deletes the namespace unless the
Application renders the Namespace object, and it never deletes StatefulSet
volumeClaimTemplate PVCs (unless the StatefulSet sets `persistentVolumeClaimRetentionPolicy.whenDeleted: Delete`).

`ecosystem-oc-control` is the one deliberate exception to "no CRDs": it is part of the
OpenChoreo family (every other `platform-openchoreo-*` Application already cascades) and the
whole family is being retired. Deleting it removes the 36 OpenChoreo CRDs and therefore every
OpenChoreo custom resource. Delete the sibling Applications first so their custom resources
are gone before the CRDs, or the CRs hang in `Terminating` with no controller.

## Deliberately left without the finalizer

| Application(s) | Why |
|----------------|-----|
| `monitoring-kube-prometheus-stack`, `monitoring-loki`, `monitoring-tempo`, `monitoring-otel-collector`, `edge-cert-manager`, `edge-ingress-nginx`, `edge-nack`, `edge-tailscale-operator` | Adopted Helm releases. `docs/gitops/HELM-ADOPTION.md` requires "no finalizer" so a deletion can never cascade to workloads, PVCs or CRDs. Loki's StatefulSet has `persistentVolumeClaimRetentionPolicy: Delete`, so a cascade would delete its PVC. ingress-nginx owns the DigitalOcean load balancer (cascade releases the public IP). cert-manager and tailscale render CRDs from chart templates. |
| `platform-ecosystem-applications`, `platform-ecosystem-foundation`, `platform-ecosystem-resources`, `platform-security` | Parents of other Applications, Namespaces (12 in `-resources`) and 44+ SealedSecrets. A cascade would delete the children, and with them everything in this table. |
| `platform-backups`, `platform-volume-backups`, `platform-volsync` | Backup machinery: 54 SealedSecrets of restic credentials, five Namespaces, VolSync CRDs. Losing the credentials with the Application would make the restic repositories unreadable without the originals. |
| `ecosystem-cnpg`, `ecosystem-barman`, `ecosystem-external-secrets`, `ecosystem-keda`, `ecosystem-kyverno`, `ecosystem-velero`, `ecosystem-dapr` | Operators that own CRDs. Deleting a CRD deletes every custom resource of that kind cluster-wide (Postgres clusters, ExternalSecrets and the Secrets they own, ScaledObjects, ClusterPolicies, Backups, Rollouts). Kyverno already carries Argo's own `pre-delete-finalizer`. |
| `ecosystem-ollama` | Tracks PVC `onemind-cpu-model` (20Gi, Retain). The cascade would delete the PVC object and leave a Released PV needing manual recovery. |
| `platform-n8n`, `platform-matrix`, `platform-backstage`, `platform-agent-fabric` | Track their own Namespace and/or PVCs (`n8n-data`, `synapse-data`). Cascade would delete the whole namespace and its volumes. `platform-agent-fabric` documents the choice in its header. |
| `platform/parked/**` | Not synced. |

## Adding or removing it

Add it in Git, let the parent sync it, and check `kubectl get application <name> -n argocd -o jsonpath='{.metadata.finalizers}'`.
To make a single deletion non-cascading, use `argocd app delete --cascade=false` (or
`kubectl delete application --cascade=orphan`), which ignores the finalizer.
