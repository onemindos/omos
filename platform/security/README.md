# Security layer (zero trust)

One place for the security tools, applied by the `platform-security` Application
(`platform/argocd/platform-security.yaml`). Each tool keeps its own namespace, labelled
`onemind.io/layer: guard`.

| Application | Role in zero trust | Namespace |
|---|---|---|
| `ecosystem-kyverno` | Policy guardrails (`kyverno-policies-audit.yaml`): image tags and namespace layer label enforced; requests/limits and privileged/hostPath audit only | `kyverno` |
| `ecosystem-falco` | Runtime threat detection; alerts not delivered yet | `falco` |
| `ecosystem-openfga` | Fine-grained authorization; used by Backstage Missions only | `onemind-core` (own namespace pending the OneMind Core decision) |

Chart values stay in `platform/ecosystem/values/` and the charts stay checksum-locked in
`platform/ecosystem/foundation-lock.json` and `chart-lock.json`.
