# kagent and agent-sandbox (parked 2026-10-17)

kagent 0.10.3 (Kubernetes AI agent runtime) and the agent-sandbox controller (v1.0.6) ran from 2026-10-08 to 2026-10-17. Measured unused: 0 UI requests, 0 agent runs and 0 tool calls in 7 days, and no ingress. Removed because Hermes is the agent harness and OneMind's agents are Hermes profiles inside the Legacy pod. The full profile (features, our configuration, reinstall notes) is in Zeus's vault: `OneMind-NEW/Software/Self-Hosted/AI and Agents/Agent Platforms/kagent - Kubernetes AI Agent Runtime.md`.

| Folder | What it holds | Moved from |
|---|---|---|
| `applications/` | `ecosystem-kagent`, `ecosystem-kagent-crds`, `ecosystem-agent-sandbox` | `platform/ecosystem/applications/` |
| `values/` | `kagent.yaml`, `kagent-crds.yaml`, `agent-sandbox.yaml` | `platform/ecosystem/values/` |
| `charts/agent-sandbox/` | Vendored agent-sandbox chart 0.1.2 | `platform/ecosystem/charts/` |
| `resources/helm-oci-kagent.yaml` | Argo CD credential for the kagent OCI chart repository | `platform/ecosystem/resources/` |
| `secrets/` | `kagent` namespace SealedSecrets (`ghcr`, `kagent-db`, `kagent-model`, `kagent-oidc`) | `platform/ecosystem/resources/secrets/` |
| `live-export/` | The 7 live objects as they ran (4 Agents, 1 ModelConfig, 2 RemoteMCPServers), cluster fields stripped | Exported from the cluster |

## Database (dropped 2026-10-17)

The `kagent` database on `onemind-data` and its `kagent` role were dropped on 2026-10-17 at Zeus's request, instead of waiting 30 days. Its Database object, managed role, `db-kagent` secret (now in `secrets/onemind-data--db-kagent.json`) and network rule were removed from `platform/ecosystem/resources/`. The only copy of the data is the `pg_dump` on Zeus's Mac: `~/Backups/omos/kagent/kagent-2026-10-17.dump` (36,606 bytes, SHA-256 `d1ed3ba3ac4e9a089afa349e84aab17c1dc32a5677f0cce660c4456aa4221964`, 15 tables).

## Bringing it back

1. Move `applications/`, `values/`, `charts/`, `resources/` and `secrets/` back to where the table says.
2. Re-add the `kagent` Namespace to `namespaces.yaml`, the `kagent-tools-clients` policy and the `kagent` entry in the `litellm-clients` allow-list in `application-network.yaml`, `"kagent"` in the model-client and OIDC bootstrap Jobs, the two `chart-lock.json` entries, and the agent-sandbox lines in `scripts/validate_*.py`.
3. The `kagent` keys were removed from the `litellm-client-keys` and `onemind-platform-oidc` SealedSecrets: generate and seal new ones.
4. Re-add the Database object, the `kagent` managed role and network rule, then restore the dump with `pg_restore`.
5. Restrict the kmcp controller's ClusterRole first: it could create and delete Deployments, Services, ServiceAccounts and ConfigMaps in every namespace.
