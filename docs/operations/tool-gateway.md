# Tool gateway runbook (agentgateway)

Design: [`docs/architecture/onemind-tool-gateway.md`](../architecture/onemind-tool-gateway.md). State: **Phase 1 (core)**. The gateway
serves `/ops/mcp` with one target (`onemind`, the existing `onemind-mcp` adapters). Hermes is not wired to it yet (Phase 5), and LiteLLM
still serves the MCP hub until Phase 6.

## What runs

| Thing | Where |
|---|---|
| Argo CD Application | `platform-tool-gateway` (`platform/argocd/platform-tool-gateway.yaml`), three sources: the chart, `platform/tool-gateway/manifests/` (recursed), the repo as `$values` |
| Chart | `agentgateway-standalone` v1.7.0 (OCI `cr.agentgateway.dev/charts`, digest `sha256:f8e74f989c6705e3e1aaaaabacf554d594e2a24867eb4097b216e019a88c14a6`) |
| Image | `cr.agentgateway.dev/agentgateway:v1.7.0@sha256:60216aafb03d37f94e952ff3d51d2f630a65b10ecaebdc0a864882f92610c45b` |
| Pod | Deployment `agentgateway` in `tool-gateway`, 1 replica, on `onemind-ecosystem-cpu` |
| Endpoint | `http://agentgateway.tool-gateway.svc.cluster.local:8080/<group>/mcp`; metrics `:15020`, readiness `:15021/healthz/ready` |
| Client key | Secret `tool-gateway/gateway-clients` key `hermes`; the same value is sealed as `hermes/onemind-gateway-client` key `ONEMIND_GATEWAY_KEY` |
| Hand-written values | `platform/tool-gateway/values.yaml` |
| Generated values | `platform/tool-gateway/values.generated.yaml` (never edit; see below) |

## How the config is built

Helm replaces lists instead of merging them, so everything list-shaped (routes, targets, `extraEnv`) is generated:

```
platform/tool-gateway/gateway.yaml         listener port, clients, groups      (shared; edited by the policy phase)
platform/tool-gateway/targets/<name>.yaml  one fragment per backend target     (owned by the backend's PR)
        |  python scripts/gateway/render_gateway_config.py
        v
platform/tool-gateway/values.generated.yaml
```

CI (`scripts/tests/test_gitops_tool_gateway.py`) fails when the committed output is stale. After merging a branch that changed
`targets/`, re-run the renderer and commit the result.

### Target fragment format

```yaml
name: grafana              # [a-z][a-z0-9]*  (no "-", "." or "_"); the file must be <name>.yaml
groups: [ops]              # group routes (keys of gateway.yaml `groups`) that expose this target
env:                       # optional; every ${VAR} used in `target` must be declared here
  TARGET_GRAFANA_TOKEN:    # must start with TARGET_<NAME>_
    secret: gateway-target-grafana   # Secret in tool-gateway (a SealedSecret under manifests/)
    key: token
target:                    # verbatim agentgateway MCP target, WITHOUT `name`
  mcp: {host: http://mcp-grafana.tool-gateway.svc.cluster.local:8000/mcp}   # or sse / openapi
  policies:
    backendAuth:
      key: "${TARGET_GRAFANA_TOKEN}"    # "Authorization: Bearer <value>"
    # other headers: requestHeaderModifier: {set: {x-api-key: "${TARGET_X_KEY}"}}
```

Rules enforced by the renderer: a literal dollar is written `$$`; any bare `$` or an undeclared `${VAR}` is an error (the gateway expands the
whole config text and a missing variable is a startup error). The Host header sent to a backend is the target host without its port
(`requestHeaderModifier` cannot change it), so a backend that checks `Host` must allow the port-less name (see `ONEMIND_ALLOWED_HOSTS`).

Every backend PR adds, and edits nothing shared except `values.generated.yaml` (regenerated):

1. `targets/<name>.yaml`
2. `manifests/egress-<name>.yaml` (a NetworkPolicy selecting the gateway pod, egress to the backend) and the ingress rule in the backend's own namespace
3. `manifests/gateway-target-<name>.sealedsecret.yaml` (one Secret per target, so PRs never share a file)
4. for a backend that runs in `tool-gateway`: its Deployment, Service and NetworkPolicy under `manifests/backends/<name>/`

## Operate

```bash
kubectl --context do-nyc3-omos-cloud -n tool-gateway get pods,svc
kubectl --context do-nyc3-omos-cloud -n tool-gateway logs deploy/agentgateway --tail=50     # JSON access log
```

Check from inside the cluster. The Hermes pod has no `kubectl` and is not wired to the key yet, so read the Secret on the workstation and
pipe it in as a curl config on stdin (the key is never printed or put in argv):

```bash
KEY=$(kubectl --context do-nyc3-omos-cloud -n hermes get secret onemind-gateway-client -o jsonpath='{.data.ONEMIND_GATEWAY_KEY}' | base64 -d)
printf 'header = "Authorization: Bearer %s"\n' "$KEY" | kubectl --context do-nyc3-omos-cloud -n hermes exec -i deploy/onemind-hermes -c hermes -- \
  curl -s -K - -H 'content-type: application/json' -H 'accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"check","version":"0"}}}' \
  -D - http://agentgateway.tool-gateway.svc.cluster.local:8080/ops/mcp
unset KEY
```

Without the `Authorization` header the same call returns `401 authentication failed`. `tools/list` needs the `mcp-session-id` response header
from `initialize` and a `notifications/initialized` first.

## Rotate the client key

Generate a new value, seal it into both namespaces (`tool-gateway/gateway-clients` key `hermes`, `hermes/onemind-gateway-client` key
`ONEMIND_GATEWAY_KEY`) with the `kubeseal` pattern in `CLAUDE.md`, commit, and restart `agentgateway` (env is read at start). Use
`--from-file=<key>=/dev/stdin` fed by `printf`, so the value never appears in argv or the shell history.

## Known behaviors (Phase 0)

- Sessions are encrypted with a key the gateway generates at start; with one replica a restart invalidates open sessions and Hermes
  reconnects. Add `session.secretName` before running a second replica.
- `apiKey` defaults to a mode that lets requests without a key through. Every route is rendered with `mode: strict`; do not remove it.
- With one dead target in a group, `failureMode: failOpen` keeps the group answering; a call routed to the dead target returns an error.
