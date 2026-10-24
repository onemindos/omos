# Security layer (zero trust)

One place for the security tools, applied by the `platform-security` Application
(`platform/argocd/platform-security.yaml`). Each tool keeps its own namespace, labelled
`onemind.io/layer: guard`.

| Application | Role in zero trust | Namespace |
|---|---|---|
| `ecosystem-kyverno` | Policy guardrails; existing audit policies in `kyverno-policies-audit.yaml` | `kyverno` |
| `ecosystem-falco` | Runtime threat detection; alerts not delivered yet | `falco` |
| `ecosystem-openfga` | Fine-grained authorization; used by Backstage Missions only | `onemind-core` (own namespace pending the OneMind Core decision) |
| `ecosystem-trivy-operator` | Prepared image-only vulnerability pilot for `nextcloud,n8n`; one scan Job; no enforcement | `trivy-system` |

Chart values stay in `platform/ecosystem/values/` and the charts stay checksum-locked in
`platform/ecosystem/foundation-lock.json` and `chart-lock.json`.

The Trivy pilot uses explicit namespace-bound RBAC instead of the chart's broad
operator binding. Its namespace/quota/network policy and RBAC are owned here;
the child owns its pinned Helm resources, and the existing monitoring app owns
the Grafana dashboard. [Runbook, capacity, acceptance and rollback](../../docs/operations/trivy-image-scan-pilot.md).
Prepared configuration is not live deployment or verified scanning coverage.
