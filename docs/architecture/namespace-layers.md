# Namespace layer label

Every namespace in the cluster carries `onemind.io/layer` with one of seven values, set by a
Namespace manifest in Git. The exceptions are `kube-system`, `kube-public` and `kube-node-lease`,
which Kubernetes creates. The Kyverno policy `require-namespace-layer-label`
(`platform/security/kyverno-policies-audit.yaml`, PR #395) enforces it: admission denies creating a
namespace without a valid value and any update that removes or invalidates it. It excludes those
three `kube-*` namespaces. Argo CD `CreateNamespace=true` creates a bare namespace unless the
Application itself contains the labelled Namespace manifest, so a new Application must ship that
manifest or set `syncPolicy.managedNamespaceMetadata.labels`.

| Layer | Meaning | Namespaces (Git and cluster agree, checked 2026-10-27) |
|---|---|---|
| `data` | Event bus and the network the data moves on | `nats`, `cnpg-system`, `onemind-search` |
| `sense` | Collection and edge sensors | `omos`, `takml`, `takserver` |
| `make-sense` | Fusion, storage and analytics | `onetak`, `clickhouse`, `onemind-core`, `onemind-data`, `stac`, `stac-validation`, `karakeep`, `nextcloud`, `vaultwarden`, `opencost` |
| `decide` | AI models and agents | `onemind-ai`, `hermes`, `agent-fabric`, `n8n`, `paperclip` |
| `act` | Human-facing surfaces and delivery | `backstage`, `matrix`, `cloudtak` |
| `guard` | Security, identity, observability, backup and CI | `auth`, `monitoring`, `falco`, `kyverno`, `velero`, `ops`, `kube`, `github-runner`, `argocd`, `cert-manager`, `ingress-nginx`, `tailscale` |
| `system` | Cluster plumbing | `default`, `dapr-system`, `keda`, `openchoreo-observability-plane` |

`kyverno` and `falco` were `security`, and `agent-fabric` was `agents`; those values were
retired. `default` is labelled by `platform/ecosystem/resources/namespaces.yaml` (Application
`platform-ecosystem-resources`, which does not prune), with `Delete=false` so Argo CD never
deletes it. It moved there from the OpenChoreo defaults Application when that was removed.

## Where the label is set

Most Namespace manifests sit beside their workload (for example `platform/n8n/namespace.yaml`).
A namespace that another Application, a Helm chart or an operator creates is labelled from an
Application that does not prune it:

| File | Application | Namespaces |
|---|---|---|
| `platform/security/shared-namespaces.yaml` | `platform-security` | `argocd`, `cert-manager`, `ingress-nginx`, `tailscale`, `cnpg-system`, `keda`, `opencost`, `dapr-system` |
| `platform/security/namespaces.yaml` | `platform-security` | `kyverno`, `falco` |
| `platform/ecosystem/resources/namespaces.yaml` | `platform-ecosystem-resources` | `default`, `onemind-ai`, `onemind-core`, `onemind-data`, `velero` |
| `platform/ecosystem/resources/retained-recovery.yaml` | `platform-ecosystem-resources` | `onemind-search`, `openchoreo-observability-plane` |
| `platform/backups/volsync/namespaces.yaml` | `platform-volume-backups` | `clickhouse`, `matrix`, `monitoring`, `nats`, `takserver` |

To label another such namespace, add it to one of these files. Do not also set
`managedNamespaceMetadata` on an Application whose namespace is listed here: Argo CD 3.4 then
applies that namespace on every sync of the Application and overwrites its
`argocd.argoproj.io/sync-options` annotation with `ServerSideApply=true`. That removes
`Prune=false,Delete=false` from the namespaces in `shared-namespaces.yaml` and
`retained-recovery.yaml`, and `Delete=false` from `default`.

## Namespaces that Git does not define

Only `kube-system`, `kube-public` and `kube-node-lease`. Kubernetes creates them, they carry no
layer label, and the policy excludes them by name.

## Defined in Git but deliberately unlabelled

`geo` (`platform/geo/namespace.yaml`) and `nats-tools` (`platform/nats-tools/namespace.yaml`,
which Application `platform-nats-tools` excludes). Neither is applied and neither exists in the
cluster. Applying either file as it stands is denied, so add the label in the change that deploys
it. `onemind-search` and `openchoreo-observability-plane` were not removed with OpenSearch and
OpenChoreo: they remain, labelled and without workloads, to hold those services' recovery volumes
(`retained-recovery.yaml`).

## Creating a namespace by hand

`kubectl create namespace <name>` cannot set a label, so Kyverno denies it. Apply a manifest that
carries the label instead, with the layer of what the namespace will hold:

```bash
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Namespace
metadata:
  name: <name>
  labels:
    onemind.io/layer: <layer>
EOF
```

- Helm: create the namespace this way first, then run `helm install` (or `helm upgrade --install`
  for a release that is not installed yet) without `--create-namespace`. Helm creates the
  namespace with only a `name` label and sends that create even when the namespace already
  exists, so Kyverno denies it either way. Upgrading an installed release creates no namespace.
- A scratch namespace (restore drill, test) takes the layer of the service it holds, for example
  `decide` for a Hermes restore. Delete it when you are done.
- A namespace that outlives the task needs its manifest in Git, beside the workload or in one of
  the files above.
