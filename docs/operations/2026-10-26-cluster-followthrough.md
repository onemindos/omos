# Cluster follow-through — 2026-10-26

User authorized work on remaining capacity, independent recovery, namespace policy,
and GitOps ownership findings. Keep all four nodes during this work.

## Namespace labels and core drift

Twelve existing namespaces lack the policy-required `onemind.io/layer`. These are
functional classifications, not the numeric trust layers in the platform roadmap:

- `argocd`, `cert-manager`, `ingress-nginx`, `tailscale`: `guard`.
- `cnpg-system`, retained `onemind-search`: `data`.
- `dapr-system`, `keda`, retained `openchoreo-observability-plane`: `system`.
- `opencost`, `nextcloud`: `make-sense`.
- `cloudtak`: `act`.

Keep explicit existing Namespace resources with their current Applications. Dapr
keeps its existing `managedNamespaceMetadata` owner and VolSync annotation. The
other seven previously untracked shared namespaces receive metadata ownership from
`platform-security`; deletion/prune/replacement/force are disabled per namespace.
A server-side dry-run preserves their UIDs and existing labels/annotations. No live
NetworkPolicy uses `onemind.io/layer` as a selector, so this classification does not
change network access. No numeric trust level, RBAC or workload spec changes.

The three Deployments currently owned by `platform-core` have first-container
resource budgets identical to Git. Remove the blanket resource-budget ignore.
Preserve only `kubectl.kubernetes.io/restartedAt` on the two named Deployments that
currently have operator restart stamps: `kube-system/sealed-secrets-controller`
and `github-runner/github-runner`. Other annotations and all resource budgets
remain visible to reconciliation. No pod rollout is expected from this change.

## Acceptance

Source validation and server dry-runs precede merge. Live acceptance must check all
12 labels, original namespace/pod identities, absence of layer-policy failures,
actual core ignore rules, and Argo health after reconciliation. Historical warning
Events may remain until normal expiration.

`make validate` currently fails on existing full-tree YAML formatting debt. Changed
YAML/schema/render checks and the GitOps regression suite are separate gates.
