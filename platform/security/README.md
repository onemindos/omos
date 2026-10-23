# Security layer (zero trust)

One place for the security tools, applied by the `platform-security` Application
(`platform/argocd/platform-security.yaml`). Each tool keeps its own namespace, labelled
`onemind.io/layer: security`.

| Application | Role in zero trust | Namespace |
|---|---|---|
| `ecosystem-kyverno` | Policy guardrails (admission control); no policies yet | `kyverno` |
| `ecosystem-falco` | Runtime threat detection; alerts not delivered yet | `falco` |
| `ecosystem-openfga` | Fine-grained authorization; used by Backstage Missions only | `onemind-core` (own namespace pending the OneMind Core decision) |

Chart values stay in `platform/ecosystem/values/` and the charts stay checksum-locked in
`platform/ecosystem/foundation-lock.json` and `chart-lock.json`.
