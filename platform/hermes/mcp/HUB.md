# OneMind MCP hub (LiteLLM)

Status: **intent, in Git; not live until merged and synced.** Verified against the running
LiteLLM `v1.104.2` image (`ghcr.io/berriai/litellm`, digest `sha256:3d88bd75…`) on
2026-10-29, read-only.

Goal: one MCP endpoint and one credential for the `onemind-hermes` agent, covering every
OneMind cluster app.

## Architecture

```
onemind-hermes (ns hermes)
   │  x-litellm-api-key: Bearer <onemind-hermes LiteLLM key>   ← the one credential
   ▼
LiteLLM proxy  http://litellm.onemind-ai.svc.cluster.local:4000/mcp/   (ns onemind-ai, Argo app ecosystem-litellm)
   │  per-key scoping: object_permission.mcp_servers / mcp_tool_permissions
   ├──► onemind      → http://onemind-mcp.hermes.svc.cluster.local:8000/mcp   (Bearer ONEMIND_MCP_TOKEN)
   │                   32 app adapters in apps.yaml (8 enabled, read-only)
   ├──► <app>        → future: an app's own MCP server (streamable HTTP)
   └──► <app>        → future: an app's OpenAPI spec, exposed as MCP tools by LiteLLM (spec_path)
```

- Where it is configured: `platform/ecosystem/values/litellm.yaml` → `proxy_config.mcp_servers`
  (rendered by the chart into the LiteLLM ConfigMap). Upstream credentials enter the pod as env
  via `extraEnvVars` (optional `secretKeyRef`) and are referenced as `os.environ/<NAME>`.
- Tools appear to the client as `<server>-<tool>`, e.g. `onemind-apps_list`. Server names
  therefore must not contain `-` (LiteLLM's `MCP_TOOL_PREFIX_SEPARATOR`).
- Network: `onemind-hermes` → `litellm:4000` is admitted by `litellm-clients`
  (`platform/ecosystem/resources/application-network.yaml`); `litellm` → `onemind-mcp:8000` is
  admitted by `platform/hermes/mcp/network-policy.yaml`. LiteLLM's own egress is unrestricted.
  `onemind-hermes` pod egress to `litellm:4000` lives in `platform/hermes/onemind-hermes/network.yaml`.

### Why LiteLLM and not agentgateway

- Already deployed and Argo-managed (`ecosystem-litellm`), already the model gateway Hermes uses,
  with its Postgres (`onemind-data`) and virtual-key system in place. No new layer, chart,
  database or CRDs.
- `onemind-hub-core-v2` is memory-saturated (CLAUDE.md, cluster review H1); LiteLLM runs on the
  ecosystem pool. A second gateway would add another always-on pod and another auth system.
- The capabilities needed are in this image: streamable-HTTP upstreams with static bearer auth,
  OpenAPI→MCP, per-key server and tool scoping, config-file (GitOps) server definitions.

### Verified capabilities (v1.104.2)

| Capability | Where verified |
|---|---|
| `mcp_servers` in `config.yaml` (GitOps), loaded at startup into memory; no DB row needed | `proxy_server.py` `_init_non_llm_configs` → `load_servers_from_config` |
| Streamable HTTP upstream, `auth_type: bearer_token` + `auth_value` (alias `authentication_token`) → `Authorization: Bearer …` | `mcp_server_manager.py` loader; `experimental_mcp_client/client.py` |
| `os.environ/NAME` resolved for every config value, MCP included | `ProxyConfig._check_for_os_environ_vars` in `get_config` (before MCP load) |
| `static_headers`, `extra_headers` (forward caller headers), `allowed_tools` / `disallowed_tools` per server | loader fields |
| OpenAPI → MCP via `spec_path` (file path or http(s) URL) | `openapi_to_mcp_generator.py` |
| Pinned `server_id` (else hashed from name/url/transport/auth, so grants break on edit) | `_pinned_config_server_id` |
| Per-key `object_permission.mcp_servers` and `mcp_tool_permissions` | `_types.py` (`LiteLLM_ObjectPermissionBase`), docs `/docs/mcp_control` |
| Endpoints: `/mcp/` (all permitted servers), `/<server>/mcp`, legacy `/mcp/sse` | `_experimental/mcp_server/server.py` mounts |
| Client auth header `x-litellm-api-key` (primary) or `Authorization` | `user_api_key_auth_mcp.py`; docs `/docs/mcp` |
| `available_on_public_internet: false` limits a server to internal client IPs | `mcp_server_manager.py` (default internal ranges, `mcp_internal_ip_ranges`) |
| DB: required for virtual keys and their `object_permission` (already present: `litellm` DB on CNPG); `store_model_in_db` only matters for UI/API-defined servers | chart values, docs `/docs/mcp` |

### Default access (read this)

- **Docs** (`/docs/mcp_control`) say a key with no MCP list at any level reaches every MCP server.
- **The v1.104.2 code** resolves a non-admin key with no key/team/access-group grant to *no* servers,
  plus any server marked `allow_all_keys: true` (`get_allowed_mcp_servers` /
  `get_mcp_server_access`). Admin keys (master key, `proxy_admin` roles) with no own
  `object_permission.mcp_servers` see **every** server.
- Design, regardless of which is right: `allow_all_keys: false` on every server,
  `general_settings.require_key_mcp_access_defined: true` (keys never inherit a team's MCP list),
  and only the `onemind-hermes` key carries `object_permission.mcp_servers`. The verification step
  below checks that an existing key (core) gets zero MCP tools; if it gets any, stop and add
  `"object_permission": {"mcp_servers": ["no-mcp-servers"]}` to the other keys.

## Owner steps after merge

### 1. Upstream credential for onemind-mcp (Secret in onemind-ai)

LiteLLM reads `ONEMIND_MCP_TOKEN` from Secret `litellm-mcp-upstreams` (optional; without it the
gateway lists no `onemind` tools and onemind-mcp answers 401). Until onemind-mcp gets a dedicated
`litellm` client (below), this reuses the existing `onemind-mcp-auth` token, sealed straight from
the live Secret so the value is never printed:

```bash
kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > /tmp/pub-cert.pem
kubectl --context do-nyc3-omos-cloud -n hermes get secret onemind-mcp-auth -o json \
  | jq '{apiVersion:"v1",kind:"Secret",type:"Opaque",
         metadata:{name:"litellm-mcp-upstreams",namespace:"onemind-ai"},
         data:{ONEMIND_MCP_TOKEN:.data.token}}' \
  | kubeseal --cert /tmp/pub-cert.pem --format json \
  | jq . > platform/ecosystem/resources/secrets/onemind-ai--litellm-mcp-upstreams.json
rm /tmp/pub-cert.pem
```

Commit it (Argo app `platform-ecosystem-resources` applies it), then restart LiteLLM once so the
optional env is picked up (`envFrom`/`env` are read at container start; Argo will not roll the
pod for a new Secret): `kubectl -n onemind-ai rollout restart deploy/litellm`.

Follow-up (needs a deployment change in `platform/hermes/mcp/`): give LiteLLM its own token via
`ONEMIND_MCP_TOKENS_FILE` (`name:token` lines, e.g. client `litellm`) so onemind-mcp's audit log
names the gateway instead of `hermes`, and so Legacy's token can rotate independently.

### 2. The onemind-hermes LiteLLM key (scoped to MCP)

Done 2026-10-29: alias `onemind-hermes-mcp` (the alias `onemind-hermes` already belongs to the Legacy
model key), `models: ["no-default-models"]`, `object_permission.mcp_servers: ["onemind"]`, sealed into
`onemind-hermes-runtime` as `ONEMIND_LITELLM_MCP_KEY`. The commands below are the original procedure.

Choose the key value locally, register it with the master key, then seal it for onemind-hermes.
No value is echoed.

```bash
KEY="sk-onemind-hermes-$(openssl rand -hex 24)"
MASTER="$(kubectl --context do-nyc3-omos-cloud -n onemind-ai get secret litellm-auth \
  -o jsonpath='{.data.master-key}' | base64 -d)"
kubectl --context do-nyc3-omos-cloud -n onemind-ai port-forward svc/litellm 4000:4000 &
PF=$!; sleep 2
jq -n --arg k "$KEY" '{
  key: $k,
  key_alias: "onemind-hermes",
  models: ["hermes-agent", "local-tools", "local-embed"],
  metadata: {service: "onemind-hermes"},
  object_permission: {
    mcp_servers: ["onemind"],
    mcp_tool_permissions: {}
  }
}' | curl -sS -o /dev/null -w '%{http_code}\n' -X POST http://127.0.0.1:4000/key/generate \
  -H "Authorization: Bearer $MASTER" -H 'Content-Type: application/json' --data-binary @-
kill $PF
```

- `mcp_tool_permissions` may list per-server tool names to narrow further, e.g.
  `{"onemind": ["apps_list", "apps_describe"]}`; empty means all of the server's tools.
- Adjust `models` to what onemind-hermes should call (or `[]` for MCP only).
- Put `$KEY` into the onemind-hermes runtime Secret (`onemind-hermes-runtime`, owned by
  `platform/hermes/onemind-hermes/`) with the same `kubectl create secret … --dry-run=client -o json
  | kubeseal` pattern, then `unset KEY MASTER`.
- A GitOps alternative is to add an `onemind-hermes` client to the
  `onemind-model-clients-v1` Job (`platform/ecosystem/resources/model-client-bootstrap.yaml`)
  with `object_permission` in its body; not done here.

### 3. Hermes config for onemind-hermes

```yaml
mcp_servers:
  onemind:
    url: http://litellm.onemind-ai.svc.cluster.local:4000/mcp/
    headers:
      x-litellm-api-key: "Bearer ${ONEMIND_LITELLM_MCP_KEY}"
```

Caveat learned from Legacy (see "Legacy 401" below): with `gateway.multiplex_profiles: true`,
Hermes resolves `${VAR}` only from the profile's `.env` / secret sources, never from container
env. Either keep multiplexing off for onemind-hermes, or make sure the variable is in the
profile's secret scope.

### 4. Verify

```bash
# From the onemind-hermes pod: initialize + tools/list through the gateway; expect onemind-* tools.
# From any other key (e.g. core): expect zero MCP tools. Master-key calls list everything (admin).
kubectl -n onemind-ai logs deploy/litellm | grep -i mcp | tail
kubectl -n hermes logs deploy/onemind-mcp   # audit lines per tool call (actor "hermes" until step 1 follow-up)
```

## Onboarding a new app (3 steps)

1. **Register.** Add an entry under `proxy_config.mcp_servers` in
   `platform/ecosystem/values/litellm.yaml`: either an MCP server (`url`, `transport: http`,
   `auth_type`, pinned `server_id`) or an OpenAPI app (`spec_path: <in-cluster URL or file>`,
   `url: <API base>`). Name without `-`. Keep `allow_all_keys: false`. Alternatively, if the app
   only needs read adapters, add it to `platform/hermes/mcp/apps.yaml` (onemind-mcp) instead.
2. **Credential.** Seal the app's token into `platform/ecosystem/resources/secrets/onemind-ai--litellm-mcp-upstreams.json`
   (or a new Secret) and add an `extraEnvVars` entry with `optional: true`; reference it as
   `os.environ/<NAME>`. Add the NetworkPolicy ingress on the app side for
   `onemind-ai`/`app.kubernetes.io/name: litellm`.
3. **Scope and verify.** Add the `server_id` to the onemind-hermes key:
   `POST /key/update {"key": …, "object_permission": {"mcp_servers": ["onemind", "<new>"]}}`
   (master key), restart LiteLLM if a new env var was added, and run the verify step.

## Risks

- **Admin keys reach every MCP server.** The master key and any `proxy_admin` key bypass per-key
  scoping unless they carry their own `object_permission.mcp_servers`. Treat the master key as a
  tool credential too.
- **Docs vs code on default access** (above). Re-check after every LiteLLM upgrade; the MCP code is
  under `_experimental/` and its permission model has changed between releases.
- **Single point of failure.** One LiteLLM replica (`replicaCount: 1`) now carries model traffic
  and all MCP tools for onemind-hermes. A restart drops both. Consider 2 replicas plus a PDB once
  memory on the ecosystem pool allows.
- **Shared upstream token.** Until the dedicated `litellm` client exists, LiteLLM and Legacy share
  onemind-mcp's `hermes` token; onemind-mcp audit cannot tell them apart, and anyone with an MCP-
  scoped LiteLLM key acts as "hermes" upstream.
- **Config changes need a rollout.** Server definitions live in the ConfigMap; the chart rolls the
  pod on config change, but a new Secret alone does not.
- **Version pin.** Verified only against `v1.104.2`; `server_id` pinning, `auth_value` alias and
  `require_key_mcp_access_defined` are recent fields.

## Legacy 401 (diagnosed 2026-10-29, not fixed)

`hermes-legacy` logs `MCP server 'onemind': … HTTP 401` every 5 minutes. The token is not stale:
`ONEMIND_MCP_TOKEN` has the same length (64) and SHA-256 prefix in the `onemind-mcp-auth` Secret,
the onemind-mcp container, the Legacy container and the gateway process, and a direct
`initialize` from the Legacy pod with it returns 200. Legacy's config holds the reference
`Authorization: Bearer ${ONEMIND_MCP_TOKEN}` (written by `onemind_mcp_config.py`).

Root cause: Legacy's `config.yaml` sets `gateway.multiplex_profiles: true`. Under multiplexing,
Hermes `v2026.9.24` resolves `${VAR}` in MCP headers through `agent.secret_scope.get_secret`,
which reads only the profile's secret scope (`/opt/data/.env` plus external secret sources) and
never falls back to `os.environ`. The token exists only as container env, so the placeholder
stays literal and onemind-mcp receives `Bearer ${ONEMIND_MCP_TOKEN}`. The dashboard process
(not multiplexed) resolves it from env and connects with 200, which is why both 200s and 401s
appear in `agent.log`. Fix options (owner/Legacy owner): add the token to the profile secret
scope (e.g. `/opt/data/.env` from the Secret in the init container), or disable multiplexing.
