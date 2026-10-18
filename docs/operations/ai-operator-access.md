# AI operator access (read-only)

How Claude Code, Codex or another local AI agent gets read-only access to the cluster
and the portal. Part of phase 1 in [`PLATFORM-ROADMAP.md`](../../PLATFORM-ROADMAP.md).

Everything here runs on the admin Mac. Nothing is added to cluster memory except one
ServiceAccount, its token Secret and two ClusterRoleBindings.

## What the agent gets

| MCP server | Source | Sees | Cannot |
|---|---|---|---|
| `onemind-kubernetes` | [kubernetes-mcp-server](https://github.com/containers/kubernetes-mcp-server) 0.0.67, `--read-only` | Every namespace: workloads, pods, logs, events, nodes, ConfigMaps, and the platform CRDs (Argo CD Applications, OpenChoreo, CNPG, NATS JetStream, Gateway API, Prometheus, KEDA, Velero, VolSync, Tailscale, Cilium, Dapr, SPIRE, ClickHouse, Kyverno reports) | Read Secrets, create, patch, delete, exec, port-forward |
| `onemind-backstage` | Portal MCP actions (`/api/mcp-actions/v1`) | Catalog entities and relations, OneMind health checks, runbooks | Any mutation (machine identity is restricted by the portal) |

Two layers of protection apply to the Kubernetes server: the MCP server only exposes
tools marked read-only, and the ServiceAccount has no write permissions anyway.

The Argo CD state (sync and health of every Application) comes through
`onemind-kubernetes` by reading `applications.argoproj.io`, so no separate Argo CD
account is needed for observation.

## Cluster side (in Git)

[`platform/kube/ai-operator.yaml`](../../platform/kube/ai-operator.yaml), applied by the
`platform-kube` Argo Application:

- ServiceAccount `kube/ai-operator` and a token Secret `kube/ai-operator-token` that the
  cluster fills in.
- ClusterRoleBinding to the built-in `view` role (no Secrets).
- ClusterRole `ai-operator-read`: get, list and watch on PersistentVolumes, storage,
  CRDs, metrics and the platform CRD groups listed above.

ConfigMaps are readable through `view`. Do not put credentials in ConfigMaps.

## Mac side

Prerequisites: admin kube context `do-nyc3-omos-cloud`, Node.js (`npx`), Tailscale
connected, and the portal machine token at `~/.config/onemind/backstage-mcp.token`
(mode 600).

```bash
scripts/ai-operator/setup.sh
```

The script:

1. Writes `~/.kube/ai-operator.kubeconfig` (mode 600) from the ServiceAccount token.
2. Checks that the identity can list pods and Argo Applications, and cannot read
   Secrets, delete pods, exec or patch deployments. It stops if any check fails.
3. Registers both MCP servers in Claude Code at user scope.
4. Adds both servers to `~/.codex/config.toml` for Codex, replacing earlier copies.

Start a new Claude Code or Codex session afterwards. `claude mcp list` should show both
servers as connected.

## Rotating or revoking

- Rotate the Kubernetes token: delete the Secret `kube/ai-operator-token` (Argo CD
  recreates it and the cluster issues a new token), then re-run the script.
- Revoke completely: remove `platform/kube/ai-operator.yaml` in a PR. Argo CD prunes the
  ServiceAccount and bindings.
- The portal token is managed as described in `platform/backstage/README.md`.

## Not included yet

- **OpenChoreo MCP** (`/mcp` on the OpenChoreo API, Authentik login): added in phase 2,
  when OpenChoreo manages something. It needs a read-only OpenChoreo role and a
  dedicated Authentik client for the agent.
- **OpenChoreo observer MCP** (logs, metrics, traces, incidents, costs): phase 5.
- **Write access of any kind.** Changes go through pull requests and Argo CD.
