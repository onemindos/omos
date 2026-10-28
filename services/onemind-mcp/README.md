# OneMind internal application MCP

Give Hermes understandable tools for applications inside Kubernetes. Kubernetes
provides private Service networking; this process adapts application APIs to MCP.
Hermes's native Tool Search provides progressive disclosure. The server does not
replace app APIs, operate Kubernetes, or grant access to every app automatically.

## Current implementation

- Official Python MCP SDK, Streamable HTTP at `/mcp`, stateless JSON responses.
- `apps_list` searches all access domains; `apps_describe` explains available operations and blockers.
- STAC imagery, Ollama model inventory, Alertmanager alerts, Loki logs and OpenCost estimates.
- Credential-gated adapters for bookmarks, workflows, Paperclip, Grafana, Tempo and MediaMTX.
- `onemind://apps` exposes the same catalog as a resource.
- NATS: `nats_server_stats`, `nats_jetstream_stats`. Each reports one node, not a
  cluster aggregate, and names it in `node`. Set `NATS_MONITOR_URLS` (comma-separated,
  at most 8) to list every node; it falls back to `NATS_MONITOR_URL`. An optional
  `node` argument (e.g. `nats-1`) reads one node; without it the first reachable node
  answers and skipped nodes appear in `unavailable_nodes`. No message reading,
  publishing or administrative operations.
- Prometheus: `prometheus_query`. Instant PromQL, 5-second server-side timeout,
  100-series limit, 1 MiB response limit, no administrative endpoints.
- A bearer token per client, audited by client name; no anonymous MCP access. See
  [Clients](#clients).
- Typed tool inputs, fixed upstream paths, no caller-selected URLs, no redirects.
- Credentials come from runtime environment, never from catalog values or prompts.
- Tool-name/client/outcome/duration audit logs, without arguments or response bodies.
- Process health/readiness routes. An upstream outage fails its tool without taking
  other application tools out of service. Readiness does not prove app availability.
- Eight concurrent tool calls, a 15-second overall call budget and bounded HTTP
  connections. Every client has access to all configured read tools.

Deployment and live acceptance evidence are recorded separately in the runbook. The owning manifests are
in `platform/hermes/mcp/`; the activation gates and rollout evidence are recorded in
[the runbook](../../docs/operations/onemind-mcp.md). Source is not live acceptance.
Per-client tokens are source only: the pinned image predates them, so the live
service accepts just the `hermes` token.

## Run locally

Python 3.12 is required. Create a virtual environment, install `requirements.lock`
and `pip install --no-deps -e .`. Supply `ONEMIND_MCP_TOKEN` (a randomly generated
secret of at least 32 characters), `NATS_MONITOR_URL`, `PROMETHEUS_URL`, and:

```sh
export ONEMIND_APPS_FILE=config/apps.yaml
export ONEMIND_ALLOWED_HOSTS=127.0.0.1:8000,localhost:8000
python -m onemind_mcp.server
```

The example app URLs must reach actual upstreams for their tools to succeed.
Run tests with `python -m pytest -q` after installing `requirements-dev.txt`.

## Clients

Each client presents its own bearer token, and the audit log records the matched
client's name. `ONEMIND_MCP_TOKEN` remains client `hermes` (Legacy) with its existing
rule of at least 32 characters. `ONEMIND_MCP_TOKENS_FILE` names a mounted secret file
that adds clients, as a JSON object (`{"name": "token"}`) or as `name:token` lines
(blank lines and lines starting with `#` are ignored). Names match `[a-z][a-z0-9_-]{0,31}`, file
tokens have at least 32 visible ASCII characters, and no two clients share a token.
`hermes` may be in the file only when `ONEMIND_MCP_TOKEN` is unset. Startup fails
closed on any invalid entry, and errors cite the entry position, never a name or
token. The file is read at startup, so changes need a rollout. Tokens identify
callers and can be revoked one at a time; they do not scope tools, so every client
can call every configured read tool.

## Add another app

1. Prefer the application's existing maintained MCP server when it already exposes
   the required operations. Connect it directly to Hermes; Tool Search can discover
   tools across servers. This server intentionally does not proxy arbitrary MCPs.
2. For an API-only app, add a module under `onemind_mcp/adapters/`. Implement
   `register(add_tool, upstream)` with typed asynchronous functions and helpful
   descriptions. The supplied `upstream.get` bounds response size/time and accepts
   only fixed reviewed paths. Select output fields deliberately; avoid exposing
   credentials and unnecessary private configuration.
3. Register that module in `adapters/ADAPTERS`, then add an `id`, `adapter`,
   `description`, `url_env`, and optional `token_env` to `platform/hermes/mcp/apps.yaml`.
   Multi-node adapters (listed in `adapters.MULTI_NODE`) receive one upstream per URL
   in `<url_env>S`; other adapters refuse a URL list.
   IDs prefix tool names, allowing multiple instances without collisions. An
   existing adapter can be reused by adding only a new catalog entry and env vars.
4. Set the new URL and least-privilege credential through the deployment and a
   SealedSecret. Add the exact upstream namespace/port to egress rules;
   check whether that app's ingress policy also needs an allowance.
5. Test successful calls, auth failures, invalid arguments, upstream outages and
   output limits. Build an immutable image, update GitOps, roll out the deployment,
   and reload Hermes's MCP connection to refresh discovery.

The catalog is loaded at startup. ConfigMap hash changes roll pods automatically
through Kustomize; secret rotation needs a rollout/reload. No automatic import of
every OpenAPI endpoint, URL execution or Kubernetes scanning is included.
Write registration uses `add_tool(fn, read_only=False)` and is omitted unless
that app opts into `allow_writes: true`; the deployed catalog grants no writes.
Use dedicated per-app identities; creation operations are non-idempotent.
Update `coverage.yaml` whenever an Argo Application is added/removed. The quality
gate checks every Git-declared Application against this inventory.
Bump Legacy's `ONEMIND_MCP_TOOLSET_REVISION` when changing tools so its startup
MCP cache refreshes after server readiness.

## Deploy and connect Hermes

Follow [the deployment runbook](../../docs/operations/onemind-mcp.md). It covers registry access,
secrets in both namespaces, existing network policies, Hermes config compatibility,
protocol smoke checks, rollout and rollback.

Progressive disclosure is a Hermes client feature, not a guarantee of MCP itself.
Current Hermes documents `tool_search`, `tool_describe`, and `tool_call`; check your
deployed revision and configuration before enabling the example. If unavailable,
ordinary MCP discovery still works for this small toolset; upgrade Hermes as a
separate change before claiming progressive disclosure is active.

## Sources

- [MCP server concepts](https://modelcontextprotocol.io/docs/2026-07-28/learn/server-concepts)
- [Official Python SDK](https://github.com/modelcontextprotocol/python-sdk)
- [Hermes MCP](https://hermes-agent.nousresearch.com/docs/user-guide/features/mcp/)
- [Hermes Tool Search](https://hermes-agent.nousresearch.com/docs/user-guide/features/tool-search)
- [NATS monitoring](https://docs.nats.io/learn/monitoring/monitoring-endpoints)
- [Prometheus API](https://prometheus.io/docs/prometheus/latest/querying/api/)
