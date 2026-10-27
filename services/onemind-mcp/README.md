# OneMind internal application MCP

Give Hermes understandable tools for applications inside Kubernetes. Kubernetes
provides private Service networking; this process adapts application APIs to MCP.
Hermes's native Tool Search provides progressive disclosure. The server does not
replace app APIs, operate Kubernetes, or grant access to every app automatically.

## Current implementation

- Official Python MCP SDK, Streamable HTTP at `/mcp`, stateless JSON responses.
- `apps_list` and `onemind://apps` describe configured applications and tools.
- NATS: `nats_server_stats`, `nats_jetstream_stats`. Each reports one node, not a
  cluster aggregate. No message reading, publishing or administrative operations.
- Prometheus: `prometheus_query`. Instant PromQL, 5-second server-side timeout,
  100-series limit, 1 MiB response limit, no administrative endpoints.
- Bearer authentication for the private Hermes client; no anonymous MCP access.
- Typed tool inputs, fixed upstream paths, no caller-selected URLs, no redirects.
- Credentials come from runtime environment, never from catalog values or prompts.
- Tool-name/client/outcome/duration audit logs, without arguments or response bodies.
- Process health/readiness routes. An upstream outage fails its tool without taking
  other application tools out of service. Readiness does not prove app availability.
- Eight concurrent tool calls, a 15-second overall call budget and bounded HTTP
  connections. The initial client has access to all configured read tools.

Deployment status: **source prepared and locally tested**. The owning manifests are
in `platform/hermes/mcp/`; the activation gates and rollout evidence are recorded in
[the runbook](../../docs/operations/onemind-mcp.md). Source is not live acceptance.

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
every OpenAPI endpoint, URL execution, Kubernetes scanning or write operations is
included. All tools are deliberately annotated read-only; before adding writes,
introduce separate authorization scopes and mutation-specific registration,
idempotency, approvals where needed, and rollback semantics. Do not add a write
tool through the existing read-only registrar.

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
