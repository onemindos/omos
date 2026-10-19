# Observer MCP

`ops/observer-mcp` provides one tailnet-only Model Context Protocol endpoint for
read-only observability. It is owned by Argo CD Application `platform-ops`. The
Tailscale ingress publishes it as:

```text
https://observer-mcp.tail717ba0.ts.net/mcp
```

Connect Claude Code or Codex with an MCP streamable-HTTP client configuration at
that URL. The endpoint itself is anonymous to the monitoring APIs and has no
Secrets or Kubernetes write privileges; tailnet membership provides the network
boundary until an application-level MCP credential is added in a later gate.
Keep the hostname private and review authorized client identities before
expanding tool access.

The kagent `RemoteMCPServer` export remains parked because kagent was not
adopted; the new server is the active observer integration.

## Tools

| Tool | Read path |
|---|---|
| `query_metrics` | Prometheus instant query |
| `query_logs` | Loki query-range API |
| `query_traces` | Tempo search API |
| `get_alerts` | Alertmanager active alerts API |
| `search_incidents` | Loki keyword search (log-based MVP, not an incident database) |
| `query_costs` | OpenCost series through Prometheus |

## Examples

```text
query_metrics: promql=sum by (namespace) (rate(container_cpu_usage_seconds_total{namespace="onetak"}[5m]))
query_logs: namespace=onetak, text=error, since=1440, limit=100
query_traces: service=onetak, since=1440, limit=20
get_alerts: label=KubePodCrashLooping
search_incidents: text=error, namespace=onetak
query_costs: namespace=onetak
```

A cost query without `promql` uses the OpenCost hourly cost metric. If the
metric name changes with a future OpenCost upgrade, pass an explicit PromQL
expression and update the server defaults in the same PR.

## Verification

After Argo CD sync, inspect the Deployment and Ingress by Service name, then call
the endpoint with an authorized AI client and execute at least one metrics and
one log query. Do not run these manifests manually; let Argo CD own the change.

## Required follow-up

- Confirm the live OpenCost metric name and the `onetak` component's trace
  service name.
- Add an authenticated machine subject for MCP access after Claude/Codex
  connectivity is proven over the tailnet.
