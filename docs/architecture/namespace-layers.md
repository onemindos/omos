# Namespace layer label

Every namespace that Git defines carries `onemind.io/layer` with one of seven values. The
Kyverno policy `require-namespace-layer-label` (`platform/security/kyverno-policies-audit.yaml`,
PR #395) enforces it: admission denies creating a namespace without a valid value and any update
that removes or invalidates it. It excludes `kube-system`, `kube-public` and `kube-node-lease`. Argo CD
`CreateNamespace=true` creates a bare namespace unless the Application itself contains the
labelled Namespace manifest, so a new Application must ship that manifest or set
`syncPolicy.managedNamespaceMetadata.labels`.

| Layer | Meaning | Namespaces labelled in Git |
|---|---|---|
| `data` | Event bus and the network the data moves on | `nats` |
| `sense` | Collection and edge sensors | `omos`, `takml`, `takserver` |
| `make-sense` | Fusion, storage and analytics | `onetak`, `clickhouse`, `onemind-core`, `onemind-data`, `stac` |
| `decide` | AI models and agents | `onemind-ai`, `hermes`, `agent-fabric`, `n8n` |
| `act` | Human-facing surfaces and delivery | `backstage`, `matrix` |
| `guard` | Security, identity, observability, backup and CI | `auth`, `monitoring`, `falco`, `kyverno`, `velero`, `ops`, `kube`, `github-runner` |
| `system` | Cluster plumbing | `default` |

`kyverno` and `falco` were `security`, and `agent-fabric` was `agents`; those values were
retired. `default` is labelled by `platform/ecosystem/resources/namespaces.yaml` (Application
`platform-ecosystem-resources`, which does not prune), with `Delete=false` so Argo CD never
deletes it. It moved there from the OpenChoreo defaults Application when that was removed.

## Namespaces that Git does not define

A namespace created by `CreateNamespace=true`, by a Helm chart, by DigitalOcean or by an
operator has no manifest, so it cannot be labelled through GitOps. Intended layer in brackets:
`argocd` (act), `cert-manager` (guard), `ingress-nginx` (guard), `keda` (guard), `opencost`
(guard), `tailscale` (data), `cnpg-system` (make-sense), `dapr-system` (decide),
`kube-system` (system), `kube-public`, `kube-node-lease`. To label them,
either add a Namespace manifest to an Application that does not prune (for example
`platform/ecosystem/resources/`) or set `syncPolicy.managedNamespaceMetadata` on the Application
that creates it.

## Defined in Git but deliberately unlabelled

`geo` and `nats-tools` (defined in Git, not
present in the cluster). (The `openchoreo-*` namespaces and `onemind-search` were removed with OpenChoreo and OpenSearch.)
