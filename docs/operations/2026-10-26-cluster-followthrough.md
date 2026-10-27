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
managed metadata did not apply its intended label after sync; preserve its VolSync
annotation and SSA option in an explicit protected Namespace instead. The eight
previously untracked shared namespaces receive metadata ownership from
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

## Measured capacity and Nextcloud corrections

Seven-day measurements do not justify cutting the roughly 11GiB reserved-memory
failover deficit. Some existing services already peak above their requests, and
the sole ecosystem worker is required by AI workload node selectors. Keep four
nodes while the owner reviews the prepared alternative: raise only the ecosystem
pool minimum from 1 to 2, making five total workers. The isolated Terraform plan
contains one in-place update, no additions/replacements/deletions at the resource
level. The provider would provision one additional worker at its current quoted
$168/month ($0.25/hour); available capacity remains subject to NYC3 provisioning.
This would provide approximately15.45GiB of N-1 request headroom after daemon and
Nextcloud corrections. The plan is not applied or committed as an approved change.

Nextcloud's cron sidecar was OOM-killed on2026-10-23T19:55:28Z. Its sampled peak
was498.2Mi at a512Mi cap, and its mounted PHP configuration permits512M per process.
Raise the request128Mi→512Mi and cap512Mi→1Gi. The specific Nextcloud Tailscale
proxy last OOM-killed at2026-10-23T20:12:57Z; sampled peak253.1Mi at its256Mi cap.
Give only that ingress a dedicated ProxyClass with256Mi request/512Mi cap, retaining
existing CPU/init budgets. Total standing request increase:576Mi. The global proxy
class and all other proxies stay as configured.

The exact pinned Nextcloud chart render changes only the two cron memory fields.
The new ProxyClass passes the installed CRD's server dry-run. Rollout may briefly
interrupt this single-replica app/proxy; verify both Ready, original PVC/PV UIDs,
unchanged image IDs and authenticated-free endpoint reachability afterward.

## Offsite copy preparation

No independent backup destination is currently configured: primary objects remain
in DigitalOcean NYC3, and the connected Cloudflare account lists no R2 bucket.
The owner was asked to select an existing B2/R2/home-storage destination without
sharing credentials in chat. The scheduled copy stays suspended until destination,
sealed credentials, seed copy and recovery test are complete. The fixed-selection
verification repair and real-engine regression cases are safe preparation now.

Verified after PRs479/480: all policy-covered namespaces now carry a valid layer
label, every original namespace UID remains unchanged, and all58 Applications are
Synced/Healthy. TAK API/config/Ignite specs match their preflight SHA256 digests;
workload and pod UIDs/restart counts are unchanged, all Ready, and Argo tracking is
present. Thus the ownership sync did not restart these three running workloads.

## Live scanner warning correction

Fresh23:33Z warning Events disproved the earlier CLI-only history acceptance.
The reports controller uses native Kubernetes int64 values in its request object;
the CLI JSON fixture loader uses float64. Kyverno1.19.1's pinned JMESPath engine
compares native int64 zero unequal to its float64 numeric literal and returns null
for native generation ordering. Thus real retired ReplicaSets were still assessed
while equivalent CLI fixtures were skipped. This was not a second Kyverno release.

Canonicalize only the history guard's numeric comparisons through
`to_number(to_string(...))`. Missing/unknown replica or generation status remains
assessed, and optional Job counters preserve their missing→zero default while
unknown/boolean/object values fail closed. The admission UID discriminator,
provider budget boundary, audit/failure behavior and all image/budget validations
remain unchanged. No controller restart, event deletion or broad policy exception.

The regression uses the exact go-jmespath dependency pinned by Kyverno1.19.1 and
compares actual native int64 contexts against JSON float64 contexts, before the
CLI boundary. It reproduces the broken comparisons and verifies all80 history
fixtures in both modes. The exact Kyverno CLI additionally runs four contexts
(background absent/empty UID and admission CREATE/UPDATE). CI installs pinned Go
so native-number tests are mandatory. Source reference:
https://github.com/kyverno/kyverno/blob/v1.19.1/pkg/engine/context/utils.go

Live acceptance must verify actual policy/autogen generations, reports-controller
skip counters for retired revisions, no active-workload validation regressions,
and no newly created history warnings during an observed rescan. Local CLI success
alone is insufficient. Historical warning Events can expire normally.
