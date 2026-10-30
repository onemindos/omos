# OneMind tool gateway: agentgateway as Legacy's single tool front door

Status: **design, in Git; nothing here is deployed.** Written 2026-10-30 from the owner decision of the same day. Live facts were read
from cluster `omos-cloud` (read-only) and from the Hermes `v0.22.5` image in `onemind-hermes` on 2026-10-30. Upstream facts come from
the agentgateway `v1.7.0` source tree and docs, and each upstream project's README (sources at the end).

Owner decision (2026-10-30): replace the LiteLLM-as-MCP-hub design in [`platform/hermes/mcp/HUB.md`](../../platform/hermes/mcp/HUB.md)
with **agentgateway** as the one tool front door for Legacy (Deployment `onemind-hermes`, namespace `hermes`). LiteLLM goes back to model
routing only. In the owner's words, Legacy can reach every app we host, every endpoint, read **and** write, fast. Deletes and physical
actions (doors, devices, scaling to zero, deleting data) still need the owner to type `execute` or `override` for that specific action.
API-key cloud services go behind the gateway. Browser-OAuth services and Composio connect to Hermes directly.

Contents: [1 Capabilities](#1-agentgateway-capabilities-verified) · [2 Inventory](#2-inventory-every-app-we-host) ·
[3 Credentials and tool classes](#3-credentials-and-tool-classes) · [4 Tool groups](#4-tool-count-strategy-groups) ·
[5 The execute gate](#5-the-execute-gate) · [6 Deployment design](#6-deployment-design) · [7 Rollout](#7-rollout-phases-and-verification) ·
[8 Risks](#8-risks-and-open-questions) · [Sources](#sources)

## 0. Summary

| Question | Answer |
|---|---|
| Gateway | agentgateway **v1.7.0** (released 2026-10-27, Linux Foundation project), **standalone Helm chart** `agentgateway-standalone` v1.7.0, no controller, no CRDs, no Gateway API |
| Why standalone, not the Kubernetes controller mode | Controller mode wants Gateway API CRDs 1.6.x; the cluster has Cilium's v1.2.1 bundle. Standalone is one Deployment plus a ConfigMap, Git is the config, `${ENV}` expansion fills credentials from SealedSecrets |
| Backends | 12 apps with an MCP server (11 first-party or project-official, 1 community), 6 apps via OpenAPI to MCP, 9 via the existing/extended `onemind-mcp` custom adapter, 13 skipped on purpose (40 rows) |
| Groups | `ops`, `dev`, `cloud`, `data`, `tak`, `home`, `work` (+ `identity`, deferred) and one `execute` route for destructive and physical tools |
| Hermes side | one `mcp_servers` entry per group; Hermes `tool_search` already defers MCP tools when they exceed 5% of context; per-server `tools.include`, `lazy`, `trust` exist in v0.22.5 |
| Execute gate | Layer 1 route split and CEL deny at the gateway; layer 2 an AuthZEN policy service whose approval arrives by the owner typing `execute <id>` in a Matrix room Hermes cannot write as the owner; layer 3 Hermes `trust: untrusted` on the `execute` server |
| Cost | gateway about 100m CPU / 128 Mi (chart default; measured footprint is a Phase 1 task), backends 25 to 100m / 64 to 256 Mi each, all on the ecosystem node pool |

## 1. agentgateway capabilities (verified)

**Version and packaging**

- Latest stable: **v1.7.0**, published 2026-10-27 (GitHub release `v1.7.0`; a `v1.7.0-rc.1` precedes it). The release notes say standalone MCP
  sessions are now always encrypted.
- Image: `cr.agentgateway.dev/agentgateway:v1.7.0`, multi-arch index digest
  `sha256:60216aafb03d37f94e952ff3d51d2f630a65b10ecaebdc0a864882f92610c45b` (resolved with `docker buildx imagetools inspect`, 2026-10-30).
  Runtime base is `cgr.dev/chainguard/glibc-dynamic`: **no shell, no node, no uvx**.
- Helm charts in OCI at `cr.agentgateway.dev/charts/` (all three pulled and inspected on 2026-10-30):

  | Chart | Version | Digest | Use |
  |---|---|---|---|
  | `agentgateway-standalone` | v1.7.0 | `sha256:f8e74f989c6705e3e1aaaaabacf554d594e2a24867eb4097b216e019a88c14a6` | **chosen**: Deployment + ConfigMap + Service + optional PodMonitor |
  | `agentgateway` | v1.7.0 | `sha256:40340d639e35206317c934e9af3e29575cf7c5377b2940679732baac6e108c15` | controller for Gateway API mode |
  | `agentgateway-crds` | v1.7.0 | (not pulled) | CRDs for controller mode |

- **CRDs and Gateway API.** The standalone chart ships no CRDs and needs no controller (`helm template` renders ServiceAccount, ConfigMap, Service,
  Deployment, PodMonitor). The controller mode needs `agentgateway-crds` plus the upstream Gateway API CRDs; the docs apply Gateway API 1.6.3.
  The cluster's `gateway.networking.k8s.io` CRDs are the **v1.2.1 standard** bundle that Cilium installed (`gatewayclass/cilium` is the only
  GatewayClass). Upgrading them would touch Cilium's CRDs, so controller mode is rejected. kgateway is not needed in either mode.
- The standalone chart's `config:` value is the same format as a local config file. Chart mode `readonly` serves it from a read-only
  ConfigMap; `database` mode adds an Admin UI backed by PostgreSQL (not used). Chart changes outside nested `config` apply without restart; other
  `config` settings restart the pod.

**Capability matrix**

| Need | Supported in v1.7.0 | Evidence | Gap or workaround |
|---|---|---|---|
| Multiplex several MCP backends on one endpoint | Yes: `backends[].mcp.targets[]`; tools prefixed with the target name (`prefixMode: always`, `conditional` or `never`) | `examples/mcp-multiplex`, `schema/config.md` | none. `failureMode` (`failClosed` default, `failOpen`) decides whether one dead target fails `initialize` for the group: use `failOpen` |
| Separate endpoints per group | Yes: one route (`matches.path.exact: /ops/mcp`) per backend group | `examples/mcp-authentication` (several MCP routes on one listener) | none |
| OpenAPI to MCP | Yes: target `openapi: {schema: {file\|url}, host, port, path}`; one tool per operation | `examples/mcp-openapi`, `schema/config.md` | **no per-operation include or exclude**, and no `readOnlyHint` annotations on generated tools (the source has no annotation handling). Workaround: commit trimmed specs generated by a script, and hide the rest with `mcpAuthorization`. OpenAPI **3.x only** (Swagger 2.0 specs need conversion) |
| stdio servers | Yes in the binary (`stdio: {cmd, args, env, clear_env}`) | `examples/mcp-multiplex` | **Not usable in the stock image** (no node or uvx) and it would run backend code and credentials inside the gateway process. Decision: **no stdio in the gateway**. Run each stdio-only server as its own Deployment behind `mcp-proxy` (`sparfenyuk/mcp-proxy` v0.13.0, stdio to streamable HTTP) |
| Streamable HTTP and SSE upstreams | Yes: `mcp: {host}`; transport policies per target (`backendAuth`, `backendTLS`, `requestHeaderModifier`) | `schema/config.md`, docs "MCP target policies" | MCP-level policies (`mcpAuthorization`, `mcpGuardrails`) are route or backend level, not per target; split routes to differ |
| Per-tool authorization | Yes: `mcpAuthorization.rules[]` with `allow`, `require`, `deny` CEL. Variables: `mcp.tool.name`, `mcp.tool.target`, `mcp.tool.arguments.*`, `mcp.methodName`, `jwt.*`, `request.headers`, `apiKey` metadata. A deny rule **filters the tool out of `tools/list`** (source test `authorization_deny_specific_tool_filters_only_that_tool`) | `schema/cel.md`, `crates/agentgateway/src/mcp/rbac.rs`, `mcp_tests.rs` | `deny` is documented as fragile (an expression failure fails to deny). Use `allow`/`require` lists. Spike A1: confirm that `allow`-list rules filter `tools/list` as deny rules do |
| External decision point | Yes: `extAuthz` with an HTTP **AuthZEN** protocol, `conditional: has(mcp.tool)` to call only on `tools/call`, request body built by CEL from `mcp.tool.*`, result read back into `extauthz.decision` and enforced with `require` | `examples/mcp-authzen` | needs a small PDP service (section 5). Spike A5: failure mode when the PDP is down must be fail-closed |
| Content guardrails | Yes: `mcpGuardrails.processors` (in-process CEL expressions or remote services) can reject or rewrite requests and responses | `examples/mcp-guardrails` | optional later: secret masking on tool results |
| Backend credentials from Secrets | Yes, by **environment expansion of the whole config file** (`${NAME}`, `${NAME:-default}`, shell-style; an unset variable is a startup error) and `backendAuth.key` (string or `{file: path}`) | `crates/agentgateway/src/config.rs` tests, `schema/config.md` | The chart's `extraEnv` accepts `valueFrom.secretKeyRef`, so SealedSecrets feed it. Non-Bearer headers (`X-N8N-API-KEY`, `x-api-key`, Basic) go through `requestHeaderModifier.set`. Spike A2: confirm that `backendAuth.key` sends `Authorization: Bearer` and that `set` values are literals. A literal `$` in CEL gets expanded: write `$$` |
| Client authentication | Yes: `apiKey` (keys with metadata), `jwtAuth` (JWKS), `basicAuth`, `mcpAuthentication` (OAuth resource metadata for browser flows; providers include authentik) | `schema/config.md`, `examples/mcp-authentication` | v1: **API key per client** (`hermes`, later `claude-code`), keys from env. JWT via Authentik is possible later |
| Observability | Yes: Prometheus metrics on `:15020/metrics` (chart `monitoring.enabled` creates a PodMonitor, scrape native or classic histograms); OTLP traces via `frontendPolicies.tracing.host`; structured access logs; chart ships `agentgateway-dashboard.json` | chart `values.yaml`, `examples/mcp-telemetry`, `schema/metrics.md` | Prometheus here selects `release: monitoring` (PodMonitor needs `monitoring.extraLabels`). No OTel Collector Service exists in `monitoring`; send traces to `tempo.monitoring.svc:4317` |
| Footprint | Rust single binary. The chart requests `cpu 100m` / `memory 128Mi` and sets no limits; the docs publish no measured numbers | chart `values.yaml` | measure in Phase 1 and record it here; set limits `1` CPU / `512Mi` until then |
| Replicas and sessions | Multiple replicas are supported by both storage modes; sessions are encrypted with a key from Secret `session.secretName` (entry `key`) | chart `values.yaml`, v1.7.0 notes | start with 1 replica (the ecosystem pool has one node); add the session key before a second replica |
| Admin UI | Yes (standalone) | docs | not exposed; no Service or Ingress for it |

**Gaps that change the design (summary).** (1) No stdio in the stock image: backends become their own Deployments. (2) OpenAPI targets cannot be
trimmed or annotated: trim specs in Git and classify tools ourselves. (3) CEL cannot see MCP tool annotations (`readOnlyHint`, `destructiveHint`),
so classification is an explicit, reviewed file, not discovery. (4) Per-target `mcpAuthorization` does not exist: classification is enforced per
route.

**Spikes that must pass before Phase 1 is built** (Phase 0, local binary with mock MCP servers, no cluster): A1 allow-rule list filtering;
A2 `backendAuth.key` header form and `requestHeaderModifier.set` literals; A3 `prefixMode: always` tool-name format as Hermes sees it
(`mcp_<server>_<target>_<tool>`, hyphens become underscores); A4 `$` escaping in CEL inside the env-expanded config; A5 `extAuthz` fail-closed
when the PDP is unreachable; A6 a failing target with `failOpen` does not break `tools/list` for the rest.

### 1.1 Phase 0 spike results (2026-10-30)

Run locally against the pinned image `cr.agentgateway.dev/agentgateway:v1.7.0@sha256:60216aafb03d37f94e952ff3d51d2f630a65b10ecaebdc0a864882f92610c45b`
(`linux/arm64` layer of the index) in Docker, with three mock MCP/OpenAPI backends (a ~60-line streamable-HTTP server) on a user network. The
config was the same shape as production: route per group, `${ENV}` expansion, read-only root filesystem, UID 65532, config mounted read-only.
The Phase 1 chart output (`helm template`, ConfigMap `config.yaml`) was also run through `--validate-only` and served end to end.

| # | Question | Result |
|---|---|---|
| A1 | Do `allow` rules filter `tools/list`, and are unlisted tools uncallable? | **Yes.** With two `allow` rules (per target, name `in [...]`) only the listed tools appear; an unlisted tool called directly returns JSON-RPC `-32602 unknown resource`. A route with **no** `mcpAuthorization` lists every tool. A single `require: "false"` rule gives an empty list (the deny-all `/execute` route works as designed). |
| A2 | Header forms | **`backendAuth: {key: "<string>"}` sends `Authorization: Bearer <string>`** (accepted by the schema as a plain string). `requestHeaderModifier.set` sends literals as written (`x-api-key: <value>`). The client's own `Authorization` header is **not** forwarded to the backend (the key is consumed by the `apiKey` policy). **The `Host` header sent to a backend is the target host without its port** and cannot be overridden with `requestHeaderModifier` (`host` is ignored); see deviation 3. |
| A3 | Tool-name format | **`<target>_<tool>`** with `prefixMode: always` (`m1_echo`). Hermes adds `mcp_<server>_`, so `onemind` in group `ops` becomes `mcp_ops_onemind_<tool>`. |
| A4 | `$` escaping | **`$$` becomes a literal `$`** in CEL rules and in `requestHeaderModifier` values (`lit$$eral` is sent as `lit$eral`; `matches("^(a\|b)$$")` works). A lone `$` is a config error, which the renderer rejects ahead of time. An unset `${VAR}` is a startup error (`environment variable not found`), so the env of the pod must be complete. |
| A5 | `extAuthz` (AuthZEN) failure mode | **Fail-closed.** A reachable PDP answering `{"decision": false}` blocks the call; an unreachable PDP blocks it too (`403 authorization failed`). `conditional: has(mcp.tool)` skips the PDP for `initialize` and `tools/list`. **Caveat for Phase 4:** a route that uses `require: '!has(mcp.tool) \|\| extauthz.decision == true'` evaluates the rule per tool during `tools/list` as well, where no decision exists, so the **list comes back empty** while an allowed call still works. Keep tool visibility in separate `allow` rules and use the PDP rule only to gate calls. |
| A6 | Dead target with `failOpen` | **Works.** `initialize` and `tools/list` succeed with the live target's tools; with `failClosed` (the default) `initialize` returns `500 failed to send message`. |
| A7 (new) | Does a route with `apiKey` reject a request with no key? | **No, not by default.** With `keys` only, a request with **no** `Authorization` header is served; only a wrong key is rejected with 401. **`apiKey.mode: strict` is required** to get 401 on a missing key. See deviation 1. |
| A8 (new) | OpenAPI target | **Works.** `openapi: {schema: {file}, host: "oa:9000"}` produced one tool per operation (`list_items`, `create_item`, `delete_item_by_id`) and executed a GET with the injected `x-api-key`. `host` takes `host:port`; setting `port` or `path` requires both (`Error: if port or path is set, both must be set`). The spec file path must exist in the pod (mount it from a ConfigMap). |
| A9 (new) | Session key and a read-only config | With a **writable** config file the gateway writes its generated session key into that file (a `config.session.key` entry). With the ConfigMap's read-only mount it starts normally and keeps the key in memory, so sessions do not survive a restart. Acceptable for one replica; set `session.secretName` before a second one. |
| A10 (new) | Footprint | Idle, three targets and about a dozen requests: **14 to 18 MiB RSS, 0.2% CPU**. `/healthz/ready` on `:15021` answers 200; 53 `agentgateway_*` series on `:15020/metrics`, including `agentgateway_requests_total` and `agentgateway_mcp_requests_total` by route and method. |
| H1 | Home Assistant `/api/mcp` with a long-lived token | **Not run:** needs the owner's Home Assistant. Stays open for Phase 3g. |

Deviations from this design that the spike forced (all applied in the Phase 1 PR):

1. **`apiKey.mode: strict` on every route** (the section 6 excerpt had `apiKey: {keys: [...]}` only, which would have accepted anonymous
   callers). The renderer emits it and a unit test fails if it is missing.
2. **Fragments, not one values file.** Helm replaces lists, so routes, targets and `extraEnv` come from `values.generated.yaml`, rendered by
   `scripts/gateway/render_gateway_config.py` from `platform/tool-gateway/gateway.yaml` and `targets/<name>.yaml` (format in
   [`docs/operations/tool-gateway.md`](../operations/tool-gateway.md)). Phase 2 extends the same renderer with classification instead of
   introducing a second mechanism. The `extraEnv` list in the section 6 excerpt therefore lives only in the generated file.
3. **One Secret per target** (`gateway-target-<name>`, SealedSecret beside the target) instead of the shared `gateway-cloud-credentials`,
   so the parallel Phase 3 PRs never edit the same sealed file. `gateway-clients` is unchanged. The Hermes copy of the client key is
   the separate SealedSecret `hermes/onemind-gateway-client` (key `ONEMIND_GATEWAY_KEY`), because `onemind-hermes-runtime` is a single
   sealed object and was not to be touched; Phase 5 can point Hermes' env at either.
4. **The backend must accept a port-less `Host`.** `onemind-mcp` checks `Host` against `ONEMIND_ALLOWED_HOSTS`; the gateway sends
   `onemind-mcp.hermes.svc.cluster.local` without `:8000`, so the port-less names were added to that list (Phase 1 PR, one env value on the
   existing Deployment). Every Phase 3 backend that validates `Host` needs the same.
5. **Resources:** requests `50m` / `64Mi`, limits `1` CPU / `256Mi` (measured idle 14 to 18 MiB), instead of `100m` / `128Mi` and
   `512Mi`. Re-measure once the Phase 3 targets exist.
6. **No Backstage Component in Phase 1.** Section 7.1 does not assign `platform/backstage/catalog/apps.yaml` to Phase 1 and the file is a
   paired, generated-runbook catalog; the Component and runbook page land with the catalog PR for the phase that follows. The `docs/STACK.md`
   row and the `coverage.yaml` disposition are in Phase 1.
7. **Egress policies per backend.** Each backend adds a `manifests/egress-<target>.yaml` NetworkPolicy (policies are additive) instead of
   editing one gateway policy, for the same parallelism reason. The Argo Application has a third, recursed `path` source for these plain
   manifests.
8. **Hermes is not touched in Phase 1.** The Phase 1 row of section 7 adds a Hermes `ops` server next to the LiteLLM entry; by owner
   instruction that moves to Phase 5. Phase 1 only adds the Hermes egress rule to the gateway and the sealed client key, and is verified with a
   curl from the Hermes pod instead of `hermes mcp test ops`.

## 2. Inventory: every app we host

Sources: `docs/STACK.md` (a 2026-10-09 snapshot, partly stale, see its addenda), `kubectl get deploy,sts,svc,ingress -A` on 2026-10-30 (read-only),
`platform/hermes/mcp/apps.yaml` and `coverage.yaml`. **Route** is the integration route chosen for the gateway: `MCP` (an MCP server exists),
`OpenAPI` (OpenAPI to MCP), `custom` (adapter in `onemind-mcp` or a new one), `skip` (reason given). **Group** is the section 4 group.
Versions are the latest upstream release seen on 2026-10-30 (GitHub API) and must be re-checked and digest-pinned at implementation.

A probe through the API-server service proxy showed most app namespaces answer nothing from outside their own NetworkPolicy (timeouts),
which is why OpenAPI availability below is marked **live** (a spec was fetched), **docs** (upstream documents it) or **verify** (to be fetched
during the owning phase, from the gateway pod once the NetworkPolicy exists).

### 2.1 Apps with an MCP server (12)

| # | App | Namespace / Service : port | MCP server (name, version) | How it runs behind the gateway | Group |
|---|---|---|---|---|---|
| 1 | Kubernetes API | `kube-apiserver` (in-cluster SA) | `containers/kubernetes-mcp-server` (release `v0.1.1`, 2026-10-29); HTTP mode with TOML (`read_only`, `toolsets`, `--config`); Helm chart `oci://ghcr.io/containers/charts/kubernetes-mcp-server` | Own Deployment in `tool-gateway`, ServiceAccount `k8s-operator`, toolsets `core,config` (+ `observability/*` optional). Not the Mac setup's `kube/ai-operator` token | ops |
| 2 | Argo CD | `argocd/argocd-server` :80/:443 | `argoproj-labs/mcp-for-argocd` (npm `argocd-mcp`; image `argoprojlabs/mcp-for-argocd`, verify); HTTP `POST /mcp` on 3000, token per request header | Own Deployment. Token sent as the `x-argocd-api-token` header by the gateway (`requestHeaderModifier`). Argo RBAC denies `delete` for the account | ops |
| 3 | Grafana (+ Prometheus, Loki, Tempo, Alertmanager datasources, alert rules) | `monitoring/monitoring-grafana` :80 (Grafana 13.1.0) | `grafana/mcp-grafana` v2.1.1; `-t streamable-http --address :8000`, `--disable-write`, `--enabled-tools`, `--allowed-hosts`; image `grafana/mcp-grafana` | Own Deployment, two instances if write is wanted: read (`--disable-write`) and write | ops |
| 4 | Backstage portal | `backstage/backstage` :7007 | Backstage MCP actions plugin (the portal already serves MCP; Secret `backstage-mcp` has `MCP_TOKEN`, `ONEMIND_MCP_HUB_TOKEN`) | Direct HTTP target, token resealed into `tool-gateway`. Read-only today (catalog, health, runbooks) | dev |
| 5 | n8n | `n8n/n8n` :5678 (n8n 2.42.4) | **Built-in instance-level MCP server**, `/mcp-server/http`, access-token auth, per-workflow opt-in | Direct HTTP target (`backendAuth.key`). The public REST API is a second OpenAPI target (see 2.2 note) | work |
| 6 | Karakeep | `karakeep/web` :3000 (0.34.1) | `@karakeep/mcp` (stdio; env `KARAKEEP_API_ADDR`, `KARAKEEP_API_KEY`; image `ghcr.io/karakeep-app/karakeep-mcp`, verify) | Own Deployment wrapping stdio with `mcp-proxy` | work |
| 7 | Home Assistant | `hermes/omos-ha-edge` ExternalName :8123 (edge host) | **Core integration "MCP Server"**, endpoint `/api/mcp`, exposure follows Assist exposure; docs describe **OAuth 2.0 + PKCE**, no mention of long-lived tokens | Gateway target with the edge address (own ExternalName in `tool-gateway`). **Spike H1**: does `/api/mcp` accept a long-lived token as Bearer? If not, run `homeassistant-ai/ha-mcp` (community, streamable HTTP, long-lived token) instead | home |
| 8 | ClickHouse "Time Machine" | `clickhouse/clickhouse-onemind` :8123 | `ClickHouse/mcp-clickhouse` v0.7.0 (PyPI `mcp-clickhouse`; HTTP/SSE transport with mandatory auth token; read-only by default, `CLICKHOUSE_ALLOW_WRITE_ACCESS` and a second flag for destructive DDL) | Own Deployment on a thin pinned Python image; read-only instance only | data |
| 9 | GitHub | `api.githubcopilot.com` (SaaS) | `github/github-mcp-server` v2.0.3; **remote hosted** `https://api.githubcopilot.com/mcp/` with a PAT; `--read-only`, `--toolsets` (remote: `X-MCP-*` headers) | Gateway target with `backendTLS`, no pod. Toolset headers fix scope | dev |
| 10 | DigitalOcean | `*.mcp.digitalocean.com` (SaaS) | `digitalocean-labs/mcp-digitalocean` v1.6.4; remote per-service endpoints (`apps`, `databases`, `accounts`, ...), PAT as Bearer | One gateway target per service endpoint; no pod | cloud |
| 11 | Cloudflare | `mcp.cloudflare.com` (SaaS) | `cloudflare/mcp` "API MCP Server" (code-execution style, a handful of tools), Bearer API token; domain servers in `cloudflare/mcp-server-cloudflare` | One target. **Token scope is the only classification** (see 3.2) | cloud |
| 12 | Nextcloud | `nextcloud/nextcloud` :8080 (35.0.2) | **Community** `cbcoutinho/nextcloud-mcp-server` v0.203.0 (image `ghcr.io/cbcoutinho/nextcloud-mcp-server`, stdio and HTTP, app password). No first-party standalone MCP server was found | Own Deployment, dedicated user and app password, never the admin credentials in `nextcloud-runtime` | work |

### 2.2 Apps through OpenAPI to MCP (6)

| # | App | Service : port | Spec | Credential model | Group |
|---|---|---|---|---|---|
| 13 | STAC API (`stac-api`, writer `stac-api-writer`) | `stac/stac-api` :8080 | stac-fastapi-pgstac publishes OpenAPI at `/api` (docs; **verify**). Trim to search, collections, items; the writer's transaction endpoints go to the `execute` side only | none for reads; writer token (verify) | data |
| 14 | TAK Portal | `onetak/tak-portal` :8080 | **live**: `/openapi.json` (OpenAPI 3.1.0, "OneMind TAK Portal" 0.1.0) | portal credentials (verify what the API accepts) | tak |
| 15 | CloudTAK | `cloudtak/cloudtak-api` :5000 (v13.113.0) | upstream API is schema-driven (verify `/api/schema` or equivalent) | existing agent token (`hermes/cloudtak-agent-token`, key `token`) | tak |
| 16 | onemind-core (+ ingestion) | `onemind-core/onemind-core` :8080 | FastAPI-style `/openapi.json` expected (**verify**; a probe was blocked by NetworkPolicy) | `core-runtime` has `READ_API_KEY`, `UI_GATEWAY_KEY`; reads use `READ_API_KEY`, writes need a new scoped key | data |
| 17 | MediaMTX | `onetak/onetak-media` | Control API v3 documents an OpenAPI file (docs, **verify** the served path and port) | existing basic credentials (`MEDIAMTX_API_CREDENTIALS`) | tak |
| 18 | Authentik | `auth/platform-auth-helm-authentik-server` :80 | `/api/v3/schema/` (documented). More than 1000 operations: **deferred**, read-only allow list only | service-account API token (owner) | identity (deferred) |

Note: n8n has both an MCP server (row 5) and a public REST API (`X-N8N-API-KEY`, Swagger at `/api/v1/docs`, **verify**) for workflow create,
update, activate and delete. Both go behind the gateway as two targets of one app; the REST target is trimmed to workflows, executions and tags.

### 2.3 Apps through a custom adapter (9)

`onemind-mcp` (`services/onemind-mcp`, Deployment `hermes/onemind-mcp` :8000, 32 catalog entries and 14 enabled read-only adapters on 2026-10-29)
**stays as a gateway target during the transition** (it needs no code change: a gateway target with its existing bearer). New write adapters
go in the same package, one module each, plus a new `scope` field (read, write, physical) so classification can be generated.

| # | App | Service : port | Today | Adapter work | Group |
|---|---|---|---|---|---|
| 19 | NATS | `nats/nats` :4222, monitor :8222 | `nats_server_stats`, `nats_jetstream_stats` (read) | **new** `nats_publish` restricted to the locked subject grammar (`evt.*`, `tak.cot.inject`; `cmd.*` is class `physical`), `nats_kv_get`, stream info. Authority bus stays in-cluster | ops (stats), tak (publish) |
| 20 | Ollama (`onemind-cpu-model`) | `onemind-ai/onemind-cpu-model` :11434 | model inventory (read) | none | data |
| 21 | OpenCost | `opencost/opencost` :9003 | namespace cost (read) | none | ops |
| 22 | LiteLLM (read only) | `onemind-ai/litellm` :4000 | model list (read) | none. Admin API (key management) is **not** exposed | ops |
| 23 | Tile38 | `onetak/tile38` :9851 | disabled | **live**: HTTP command API answers (`/server`). New adapter: `NEARBY`, `WITHIN`, `GET` read; `SET`, `FENCE` write; `DEL`, `DROP`, `FLUSHDB` destructive. No arbitrary command passthrough | data |
| 24 | WebODM | `onetak/webodm` :8000 | disabled | **live**: DRF root answers. No OpenAPI: projects and tasks adapter (list, create task = write, cancel/remove = destructive) | data |
| 25 | TAK Server (Marti API) | `takserver/takserver-api-service` :8443/:8446 | disabled | mTLS client certificate `legacy-gateway` (owner generates); missions, contacts, CoT query. Writes (mission edit, CoT) via NATS `tak.cot.inject` preferred | tak |
| 26 | Paperclip | `paperclip/paperclip` :3100 | disabled | company-scoped agent token; issues and backlog create. Check first whether Paperclip ships its own MCP server | work |
| 27 | TAK-ML | `takml/takml-server` :8234 | disabled | needs a reviewed detection/query contract first | data |

The existing read adapters for Prometheus, Alertmanager, Loki and Tempo (`onemind-mcp`) are superseded by Grafana's MCP datasource tools
(row 3) and are retired in the final phase; they are not counted again.

### 2.4 Skipped on purpose (13)

| # | App | Why |
|---|---|---|
| 28 | Vaultwarden | Excluded in `apps.yaml` and the mission: no agent access to the personal vault |
| 29 | FedHub (`fedhub-ui`, `fedhub-mongodb`) | Federation administration is a security setting; no broad admin inferred |
| 30 | OpenFGA | Authorization policy store: editing it is a privilege change |
| 31 | Matrix (Synapse, Element) | Hermes has a native Matrix platform (`MATRIX_ACCESS_TOKEN`); the Synapse admin API is not needed |
| 32 | Headlamp | UI over the Kubernetes API; row 1 covers the API |
| 33 | tileserver-gl, STAC Browser | Tile and browse UIs; tiles are binary |
| 34 | Agent Console, `agent-tap` | Bus UIs; Legacy uses NATS directly |
| 35 | NATS web UIs (`heynats`, `nats-dashboard`, `nats-ui`, `nats-surveyor`) | UIs and exporters |
| 36 | Redpanda Connect (`omos`) | Internal pipeline; a streams REST API exists but offers nothing Legacy needs |
| 37 | TAK internals (Ignite, messaging, plugins, config, LDAP, Mumble bridge, murmur) | Reached through TAK Server (row 25) and NATS; admin surface only |
| 38 | Controllers: cert-manager, ingress-nginx, Tailscale operator, KEDA, Dapr, Kyverno, Falco, CNPG, Barman, Velero, VolSync, metrics-server | No app API; their custom resources are readable and patchable through row 1 |
| 39 | Databases: Postgres (authentik, backstage, matrix, n8n, paperclip, takml, TAK, WebODM, CNPG `onemind-data`), Redis, MongoDB | No raw SQL for the agent (existing rule in `apps.yaml`); app APIs instead |
| 40 | GitHub runner, Hermes itself and its dashboard | CI and the agent are not tools for the agent |

**Totals: 40 rows = 12 MCP (11 first-party or project-official + 1 community), 6 OpenAPI, 9 custom, 13 skip.**

### 2.5 Direct to Hermes (not behind the gateway)

| Service | Why direct | Setup |
|---|---|---|
| Notion, Google Workspace (Drive, Gmail, Calendar), Todoist, Dropbox, Zoho (Mail, Calendar) | Browser OAuth with the owner's account; a gateway would need to hold a refresh token and replay a consent screen it cannot show. Listed from the connectors the owner already uses; confirm which Legacy needs | `hermes mcp add` and `hermes mcp login` / `hermes mcp reauth` (OAuth support is built in) |
| Composio | An aggregator with its own account linking (`composio link`) and its own tool surface | one `mcp_servers` entry with the Composio key; its tools are deferred by `tool_search` |

These need egress on 443 only, which `network.yaml` already allows. They are excluded from the execute gate, so the **owner must not link a
Composio or OAuth account with delete-capable scopes** unless the Hermes `trust: untrusted` prompt (section 5, layer 3) is turned on for it too.

## 3. Credentials and tool classes

### 3.1 Tool classes

| Class | Meaning | Where it is exposed | Gate |
|---|---|---|---|
| **read** | no state change | group route | none |
| **write** | reversible or low-blast-radius change (create, update, annotate, sync, trigger) | group route | none (the owner asked for fast read and write) |
| **destructive** | deletes data, scales to zero, drops, kills sessions, removes users, anything unrecoverable | `execute` route only | `execute <id>` typed by the owner |
| **physical** | moves something in the world (locks, doors, garage, alarm, valves, sirens, `cmd.*` bus subjects) | `execute` route only | `override <id>` typed by the owner |

**How destructive and physical tools are identified.** Never by discovery. A reviewed file per target lives in
`platform/tool-gateway/policy/classification/<target>.yaml` with `read`, `write`, `destructive`, `physical` lists (exact names, globs and
argument predicates). **A tool that appears in no list is treated as destructive** (fail closed). `scripts/gateway/classify_tools.py` (Phase 2)
diffs a live `tools/list` against the file and fails CI on drift. Sources by integration type:

- MCP servers: tool names, plus argument predicates where the name is not enough (`resources_scale` with `replicas == 0`; Home Assistant
  `domain` in `lock, cover, alarm_control_panel, siren, valve`).
- OpenAPI targets: generated from the HTTP method: `GET/HEAD` read, `POST/PUT/PATCH` write, `DELETE` destructive, plus explicit overrides
  (`POST .../kick`, `POST .../restart`, `DELETE`-like verbs).
- Custom adapters: a `scope` per tool in the adapter module.
- Code-execution tools (Cloudflare `execute`, ClickHouse with writes on): cannot be classified by name. They are connected **read-only by
  token scope** on the group route and, if a write token is ever created, exposed on `execute` only.

**Provider-side hard limits** (second line, independent of our classification): a DigitalOcean PAT created without delete scopes; GitHub
fine-grained PAT without admin and repo-deletion; Argo CD account policy without `delete`; Kubernetes RBAC with no `delete` verb outside pods
and no Secret access; ClickHouse user with `readonly=1`; Grafana service account role Editor, not Admin; Nextcloud app password for a normal user.

### 3.2 Per-app credentials

"Exists" means a Secret already in the cluster (names of keys only; values were not read). Secrets cannot cross namespaces, so a reused
credential is **resealed** into `tool-gateway` from the live Secret with the `kubeseal` pattern in `HUB.md` (no value printed).

| App | Credential | In-cluster today | Owner action | Read | Write | Destructive and how identified |
|---|---|---|---|---|---|---|
| Kubernetes | ServiceAccount `tool-gateway/k8s-operator` (view + patch deployments/statefulsets, pods delete, `scale`; no Secrets) | none | none: created by Git | `pods_*`, `resources_get/list`, events, logs, top | `resources_create_or_update`, rollout restart via patch | `*_delete`, `resources_delete`, `resources_scale` with `replicas == 0`, node drain, `helm_uninstall`; RBAC has no delete outside pods |
| Argo CD | account `legacy` API token | `argocd/argocd-secret` has `accounts.mcp-read.tokens`; `hermes/onemind-mcp-auth` key `ARGOCD_API_TOKEN` (read) | create account `legacy` (get, sync, action) and a token; keep `delete` off in `argocd-rbac-cm` | list/get/logs/events | sync, refresh, resource actions | `delete_application`, `terminate_operation`, resource delete |
| Grafana | service-account token | `monitoring-grafana` (`admin-user`, `admin-password`, `ldap-toml`), `onemind-mcp-auth` key `GRAFANA_API_TOKEN` (read) | create Editor SA token for writes | search, dashboards, datasource queries, alert rules | create/update dashboards, annotations, silences | `delete_dashboard`, `delete_*` (name pattern) |
| Backstage | `ONEMIND_MCP_HUB_TOKEN` | `backstage/backstage-mcp` keys `MCP_TOKEN`, `ONEMIND_MCP_HUB_TOKEN` | none | catalog, health, runbooks | none exposed | none |
| n8n | MCP access token; REST API key | `n8n-secrets` has app keys only | create both in n8n (Settings > Instance-level MCP, Settings > n8n API), enable the target workflows | list/get workflows, executions | execute workflow, create/update/activate | `DELETE /workflows/*`, `DELETE /executions/*`, `DELETE /credentials/*` (credentials endpoints are not exposed at all) |
| Karakeep | API key | `karakeep-secrets` has `MEILI_MASTER_KEY`, `NEXTAUTH_SECRET` only | create key `legacy-gateway` (Settings > API Keys) | search, list bookmarks, lists, tags | create bookmark, tag, list | delete bookmark/list/tag |
| Home Assistant | long-lived access token (or OAuth, see H1) | none in cluster | create dedicated non-owner user, token; decide admin requirement; expose only non-physical entities to Assist | `GetLiveContext`, state queries | lights, scenes, climate, media | **physical**: `lock`, `cover` (door and garage), `alarm_control_panel`, `siren`, `valve` by argument `domain` and area; names `HassTurnOn/Off`, `HassOpenCover`, `HassLockLock` |
| ClickHouse | user `legacy_ro` (profile `readonly=1`) | `omos-clickhouse-creds` (`CH_PASSWORD`), `clickhouse-backup-creds` (`CH_USER`, `CH_PASSWORD`) are other users' | create `legacy_ro` and seal its password | `run_select_query`, `list_databases/tables` | none (writes disabled) | DDL/DML are impossible by user profile |
| GitHub | fine-grained PAT, org `onemindos`, repo-limited | `github-runner-secret` `ACCESS_TOKEN` (runner, do not reuse) | create the PAT | read repos, issues, PRs, actions | issue/PR create and comment, branch push, label | `delete_*`, merge to default branch (class write, protected by branch rules), workflow dispatch is write |
| DigitalOcean | PAT with custom scopes | none | create PAT: read all, create/update chosen (apps, doks, networking), **no delete** | list/get | create, update, resize | `*_delete*`, `droplet_destroy`, `db_cluster_delete`, `doks_delete_cluster`; also impossible by token scope |
| Cloudflare | API token (read) and optionally a second write token | none | create read-only token (account, zone, DNS, workers read) | `search`, read-only `execute` calls | none now | `execute` is a free-form API call: write token, if ever, exposed on `execute` only |
| Nextcloud | dedicated user app password | `nextcloud-runtime` has **admin** and DB passwords: do not reuse | create user `legacy` + app password | files list/read, calendars, contacts | create/update files and events | delete file/folder/event/contact |
| STAC | none (read), writer token (verify) | none | decide on writer use | `GET /search`, `/collections*` | `POST/PUT` transactions on `stac-api-writer` | `DELETE` collection/item |
| TAK Portal / CloudTAK | portal credentials; CloudTAK agent token | `hermes/cloudtak-agent-token` (`token`) | confirm TAK Portal auth | list | enroll user/profile, create marker | delete user/connection, revoke certificates |
| onemind-core | `READ_API_KEY` (read); new write key | `onemind-core/core-runtime` has `READ_API_KEY`, `UI_GATEWAY_KEY`, `NATS_*` | create a scoped write key later | entity and ingestion reads | domain writes | `DELETE` endpoints |
| MediaMTX | basic credentials | `onemind-mcp-auth` key `MEDIAMTX_API_CREDENTIALS` | none | paths, readers, publishers | none | kick sessions, `config/paths/delete` |
| NATS | existing monitor access; NATS user for publish | `core-runtime` has `NATS_USER`, `NATS_PASSWORD` (core's, do not reuse) | create user `legacy-gateway` with publish only on `evt.>` and `tak.cot.inject` | stats, stream info | publish `evt.*`, `tak.cot.inject` | **physical**: any `cmd.>` publish; stream/consumer delete is not offered |
| Tile38, WebODM, TAK Server, Paperclip, TAK-ML | per adapter | none (TAK Server mTLS cert; WebODM user; Paperclip agent token) | create per row 2.3 | reads | writes per adapter | `DROP`, `FLUSHDB`, `DEL` bulk; WebODM task remove; TAK mission delete |
| Gateway client key (Hermes) | API key `hermes` | none | none (generated at implementation, sealed into `hermes` and `tool-gateway`) | n/a | n/a | n/a |
| Execute gate (Matrix) | bot account `@gate:onemindos.dev` + token, private room `approvals` | `matrix/matrix-user-credentials` has `legacy-password`, `zeus-password` (users', not the bot's) | create the bot and the room, invite only the owner | n/a | n/a | n/a |

### 3.3 Owner must create (checklist)

1. GitHub fine-grained PAT (org `onemindos`, no admin).
2. DigitalOcean PAT with custom scopes, no delete.
3. Cloudflare read-only API token.
4. Grafana service-account token (Editor) for writes; reads can reuse the existing `GRAFANA_API_TOKEN`.
5. Argo CD account `legacy` with a token and an RBAC policy without `delete`.
6. n8n MCP access token and n8n REST API key (and enable MCP on chosen workflows).
7. Karakeep API key `legacy-gateway`.
8. Home Assistant long-lived token for a dedicated user, and the Assist exposure review (no locks, doors or garage exposed until the gate is live).
9. ClickHouse user `legacy_ro` (`readonly=1`).
10. Nextcloud user `legacy` and an app password.
11. NATS user `legacy-gateway` (publish `evt.>` and `tak.cot.inject` only).
12. TAK Server client certificate `legacy-gateway`; TAK Portal credentials for the API; WebODM user; Paperclip company-scoped agent token (as each adapter lands).
13. Matrix bot `@gate:onemindos.dev` and the private `approvals` room with the owner only.
14. Authentik service-account token (only when the deferred `identity` group is taken up).
15. Browser OAuth consent for Notion, Google, Todoist, Dropbox, Zoho, and the Composio account link (direct to Hermes).

Reuse, no owner action beyond resealing: Backstage hub token, CloudTAK agent token, MediaMTX credentials, `READ_API_KEY`, existing Argo CD and Grafana read tokens.

## 4. Tool-count strategy (groups)

Hundreds of tools in one context hurt tool choice and cost tokens. Three controls stack:

1. **Group routes** on the gateway, each its own MCP endpoint with its own allow list: `http://agentgateway.tool-gateway.svc.cluster.local:8080/<group>/mcp`.
2. **Hermes per-server controls** (verified in `/opt/hermes/tools/mcp_tool_registration.py` and `mcp_tool_discovery.py`, v0.22.5):
   `tools.include` (whitelist, globs; `[]` registers nothing), `tools.exclude`, `enabled`, `lazy` (connect on first use, schema cache), `trust`,
   `timeout`, `connect_timeout`, `tools.resources`/`tools.prompts` switches.
3. **Hermes Tool Search** (`tools.tool_search.enabled: auto`, `threshold_pct: 5`, `listing_max_tokens: 4000`): MCP toolsets are deferrable, so
   when the MCP tools exceed the listing budget they are replaced by `tool_search`/`tool_describe`/`tool_call` bridges and surfaced on demand
   (tier 1 name and description manifest, tier 2 one line per server). No per-session toolset toggle was found for MCP servers; `enabled: false`
   is a config-level switch. So: all groups enabled, `lazy: true`, an `include` list per group, and Tool Search as the safety net.

Naming: `prefixMode: always` on the gateway gives `<target>_<tool>`; Hermes prefixes `mcp_<server>_` and replaces hyphens and dots with
underscores, e.g. `mcp_ops_grafana_query_prometheus`. Keep server and target names free of `-`.

| Group | Endpoint | Targets | Est. tools raw -> exposed (budget 40 eager) | Notes |
|---|---|---|---|---|
| `ops` | `/ops/mcp` | kubernetes, argocd, grafana (+ datasources), nats stats, opencost, litellm models | about 90 -> 40 | the daily driver: health, rollouts, alerts, logs, traces |
| `dev` | `/dev/mcp` | github, backstage | about 50 -> 40 | GitHub toolsets `repos,issues,pull_requests,actions,context` only |
| `cloud` | `/cloud/mcp` | digitalocean (apps, doks, networking, databases, accounts), cloudflare | about 60 -> 35 | DigitalOcean services chosen by endpoint, not by filtering |
| `data` | `/data/mcp` | clickhouse, stac, tile38, webodm, onemind-core, ollama, takml | about 50 -> 40 | |
| `tak` | `/tak/mcp` | tak server, cloudtak, tak-portal, mediamtx, nats publish | about 50 -> 40 | `cmd.*` and CoT physical actions only on `execute` |
| `home` | `/home/mcp` | home assistant | 15 to 25 | one server, exposure set in Home Assistant |
| `work` | `/work/mcp` | n8n (MCP + REST), karakeep, nextcloud, paperclip | about 70 -> 40 | |
| `identity` | `/identity/mcp` | authentik (read-only allow list) | deferred | off until asked for |
| `execute` | `/execute/mcp` | destructive and physical tools of every target | small, often under 30 | closed until the PDP is live (section 5); `trust: untrusted` in Hermes |

Tool counts are estimates from upstream READMEs and must be replaced with `hermes mcp test` numbers in Phase 3.

Hermes configuration shape (illustrative; values live in the runtime config of the PVC, not in `config.defaults.yaml`, see Phase 5):

```yaml
mcp_servers:
  ops:
    url: http://agentgateway.tool-gateway.svc.cluster.local:8080/ops/mcp
    headers:
      Authorization: Bearer ${ONEMIND_GATEWAY_KEY}
    lazy: true
    timeout: 60
    connect_timeout: 20
    tools:
      include: ["grafana_*", "kubernetes_pods_*", "argocd_*", "..."]   # generated from classification
  # dev, cloud, data, tak, home, work: same shape
  execute:
    url: http://agentgateway.tool-gateway.svc.cluster.local:8080/execute/mcp
    headers:
      Authorization: Bearer ${ONEMIND_GATEWAY_KEY}
    trust: untrusted      # every non-readOnly call asks the owner in Hermes first
    lazy: true
```

Caveat carried over from `HUB.md` ("Legacy 401"): with `gateway.multiplex_profiles: true` Hermes resolves `${VAR}` in headers only from the
profile's secret scope, not container env. `onemind-hermes` keeps multiplexing off.

## 5. The execute gate

Requirement: deletes and physical actions need the owner to type `execute` (destructive) or `override` (physical) **for that specific action**.
Three layers, from cheap to strong. Be honest about which is a security boundary.

### 5.1 What Hermes has natively (verified, v0.22.5)

- Per-server `trust: full | untrusted` (`tools/mcp_tool_handlers.py`, `_trust_gate_check`). On an `untrusted` server **every tool whose
  discovery-time `readOnlyHint` is not exactly `true` asks the human before the RPC fires**; malformed or missing hints count as write-capable; a
  typo in `trust` counts as untrusted. The prompt goes to the active surface (Telegram, Matrix, dashboard) through `request_elicitation_consent`;
  with no one to answer (cron, single-query, API-server without callback, timeout, error) it **declines** (fail closed).
- The answer is **Approve or Deny**, not a typed word, and it applies to one call.
- It is code in the Hermes process, not an LLM judgement, so the model cannot approve itself. But it is **not a security boundary against a
  compromised or prompt-injected agent**: Hermes holds the gateway key in its environment and its `terminal` tool runs in the same container
  with egress to the gateway. A shell `curl` with that key would skip the Hermes prompt. OpenAPI-generated tools carry no `readOnlyHint`, so on
  an untrusted server they would all prompt; that is acceptable only on the `execute` server.

### 5.2 Layer 1: route split and CEL at the gateway (cheap, enforceable)

- The group routes (`/ops/mcp`, ...) carry `mcpAuthorization.rules` generated from the classification file: an `allow` for every listed
  `read` and `write` tool (and argument predicates such as scale not to zero). Everything else, including unlisted tools, is **not in
  `tools/list` and not callable** on those routes.
- The `/execute/mcp` route lists exactly the `destructive` and `physical` tools. Until the PDP exists, this route has an explicit
  `deny`-all (`require: "false"`), so destructive tools simply do not exist for Legacy. Nothing regresses: Legacy has no destructive tool today.
- Client key `hermes` is accepted on all routes. A second key (`owner-console`) could be accepted on `/execute` only, for human tooling.
- Limit: this layer decides on tool name and arguments, not on a human. Someone holding the Hermes key can still try the `/execute/mcp` route
  directly. That is what layer 2 closes.

### 5.3 Layer 2: `execute-gate`, an AuthZEN policy service with out-of-band approval (the real gate)

A small service in `tool-gateway` (`services/execute-gate`, about 64 Mi, Phase 4) that the gateway calls on **every `tools/call` of the
`/execute/mcp` route** via `extAuthz` (AuthZEN `evaluation` request: subject `hermes`, resource `target/tool`, action `tools/call`, context =
canonical argument JSON). It holds the policy file rendered from the same classification YAML.

```
Legacy ──tools/call──► gateway /execute/mcp ──AuthZEN──► execute-gate ──► decision
                                                           │ no approval for (tool, args-hash): decision=false,
                                                           │   creates request A7F3 (TTL 10 min), posts it to the Matrix room "approvals"
owner (Matrix, @zeus) ─"execute A7F3" or "override A7F3"──►│ sender must be the owner, id must be open, word must match the class
Legacy retries ────────────────────────────────────────────► decision=true once, approval consumed
```

- The approval is bound to **tool + canonical arguments hash** and is **single use** with a 10-minute TTL. A changed argument needs a new
  approval. "execute" authorizes `destructive`, "override" authorizes `physical`; the wrong word does nothing.
- The request message names the target, tool, and a redacted argument summary, and shows the id. Hermes sees a 403 and is told (SOUL text in
  Phase 5) to report "waiting for owner approval A7F3" and not to retry in a loop. Whether the 403 body reaches the model depends on the gateway,
  so the Matrix message, not the error text, is what informs the owner.
- Storage: NATS KV bucket `exec-approvals` (a new `KeyValue` CR in `platform/nats/streams/`, R3, TTL) so a restart keeps pending approvals;
  audit as `evt.toolgw.approval.{requested,granted,denied,consumed}` on the bus, and decisions on stdout for Loki.
- Out-of-band channel: private Matrix room `approvals` on the in-cluster Synapse. Members: the owner and the bot `@gate`. Legacy's Matrix
  identity `@legacy` is not a member and cannot be invited by itself (Hermes joins invites from the owner only). The gate accepts approvals
  **only** from the owner's MXID, verified from the event sender as authenticated by Synapse. Telegram is not used: the bot token belongs to Hermes.
- Network: only the gateway pod may reach `execute-gate:8080`; only `execute-gate` may reach Synapse and NATS KV. Hermes has no route to it and
  no credential for the room. Fail closed: gateway `extAuthz` failure mode deny (spike A5), gate answers false on any error.
- Hole analysis, enforceable vs not:
  - **Enforced:** an injected or malicious agent, with the Hermes key and a shell, cannot run a classified destructive or physical tool
    without the owner typing the word for that exact call. It cannot forge the owner's Matrix message or read the room.
  - **Enforced by default deny:** an unclassified tool is destructive.
  - **Not enforced:** a destructive tool **misclassified as write** runs on a group route with no gate. Mitigations: classification is reviewed
    in Git and drift-checked against live `tools/list`; provider-side limits in 3.1; a periodic audit query in Loki for write calls with
    destructive-looking names.
  - **Not enforced:** side effects hidden inside a "write" (an n8n workflow that deletes things when run). The owner decides which workflows
    to enable for MCP; workflow execution is class write only for workflows listed in the allow list.
  - **Not enforced:** the owner's own Matrix account being compromised, or the owner typing `execute` without reading the summary.

### 5.4 Layer 3: Hermes `trust: untrusted` on the `execute` server (UX and defence in depth)

Kept permanently even though the gate exists. It produces the second prompt in the owner's chat before anything is sent, and it protects the
interim before Phase 4 and any misconfiguration where the gateway route is opened by mistake. Cost: two prompts for a destructive action
(Approve in chat, then `execute <id>` in `approvals`). Accepted; destructive actions are rare. Cron, API-server and single-query sessions never
get past it (they decline), which matches `approvals.cron_mode: deny` in `config.defaults.yaml`.

### 5.5 Summary

1. The gateway exposes destructive and physical tools only on `/execute/mcp`, and denies everything there until the gate exists.
2. The gate approves one exact call after the owner types `execute <id>` or `override <id>` in a room Legacy cannot write as the owner.
3. Hermes additionally prompts on that server (`trust: untrusted`); misclassification as "write" is the residual risk, bounded by
   provider-side token scopes and RBAC.

## 6. Deployment design

**Where.** New namespace `tool-gateway`, label `onemind.io/layer: decide`, Pod Security `restricted` (like `hermes` and `onemind-ai`). Everything
on the ecosystem pool (`nodeSelector: onemind.io/compute-role: ecosystem`, as LiteLLM does); `onemind-hub-core-v2` is memory-saturated
(57 to 67 percent used on 2026-10-30) and the ecosystem node was at 53 percent of 29 GiB. Nothing is added to hub-core.

**Argo CD.** Application `platform-tool-gateway` in `platform/argocd/` (sync wave 3, `prune: false`, `selfHeal: true`, `CreateNamespace`,
`ServerSideApply`), two sources like `ecosystem-litellm`:

```yaml
sources:
- repoURL: cr.agentgateway.dev/charts          # OCI, no oci:// prefix, as ghcr.io/berriai is used for LiteLLM
  chart: agentgateway-standalone
  targetRevision: v1.7.0                      # chart digest sha256:f8e74f98...14a6 recorded in docs and in the lock check
  helm:
    releaseName: agentgateway
    valueFiles:
    - $values/platform/tool-gateway/values.yaml              # hand-written: image, resources, service, monitoring
    - $values/platform/tool-gateway/values.generated.yaml    # generated: binds, routes, targets (Phase 2)
- repoURL: https://github.com/onemindos/omos.git
  targetRevision: main
  ref: values
```

and a second Application (or the same path as a kustomize source) for the plain manifests under `platform/tool-gateway/` (namespace,
NetworkPolicies, backends, SealedSecrets, `execute-gate`). The helm render check in CI (`scripts/validate_helm_renders.py`) covers the chart.
Every new Application must get a disposition in `platform/hermes/mcp/coverage.yaml` (the coverage CI) and a row in `docs/STACK.md`.

**Values (excerpt).**

```yaml
fullnameOverride: agentgateway
namespaceOverride: ""
replicaCount: 1
image:
  tag: v1.7.0@sha256:60216aafb03d37f94e952ff3d51d2f630a65b10ecaebdc0a864882f92610c45b
resources:
  requests: {cpu: 100m, memory: 128Mi}
  limits: {cpu: "1", memory: 512Mi}     # revise after the Phase 1 measurement
nodeSelector: {onemind.io/compute-role: ecosystem}
podSecurityContext: {runAsNonRoot: true, seccompProfile: {type: RuntimeDefault}}
gateway:
  service:
    type: ClusterIP
    ports: [{name: mcp, port: 8080, targetPort: 8080, protocol: TCP}]
monitoring:
  enabled: true
  extraLabels: {release: monitoring}
extraEnv:                                  # one secretKeyRef per credential; all optional: true
- {name: GATEWAY_KEY_HERMES, valueFrom: {secretKeyRef: {name: gateway-clients, key: hermes, optional: true}}}
- {name: GITHUB_TOKEN, valueFrom: {secretKeyRef: {name: gateway-cloud-credentials, key: github, optional: true}}}
# ...
config:
  frontendPolicies:
    tracing: {host: "tempo.monitoring.svc.cluster.local:4317", randomSampling: true}
  binds:
  - port: 8080
    listeners:
    - routes:
      - name: ops
        matches: [{path: {exact: /ops/mcp}}]
        policies:
          apiKey: {keys: [{key: "${GATEWAY_KEY_HERMES}", metadata: {client: hermes}}]}
          mcpAuthorization:
            rules:
            - allow: 'mcp.tool.target == "grafana" && mcp.tool.name in ["search_dashboards", "query_prometheus"]'   # generated
        backends:
        - mcp:
            prefixMode: always
            failureMode: failOpen
            targets:
            - name: grafana
              mcp: {host: "http://mcp-grafana.tool-gateway.svc.cluster.local:8000/mcp"}
```

(The `config` block is illustrative; field names come from `schema/config.md` v1.7.0 and are checked by the Phase 0 spike and `helm template`.)

**Image pins.** Gateway: the digest above. Every backend image is pinned by digest in the PR that adds it; the version column in section 2 is the
tag to resolve. Only the gateway image and charts were resolved in this design.

**How official MCP servers run.** Own Deployment per backend in `tool-gateway`, ClusterIP Service, 1 replica, restricted pod security,
read-only root filesystem, its own SealedSecret mounted as env, requests 25 to 100m CPU and 64 to 256 Mi memory, limits about 4x requests.
Not sidecars (a sidecar shares the gateway's blast radius and Secrets) and not stdio in the gateway (section 1). SaaS MCP servers (GitHub,
DigitalOcean, Cloudflare) and n8n, Backstage and Home Assistant are plain gateway targets with `backendTLS` and `backendAuth`/header injection.
Stdio-only servers (Karakeep, Nextcloud's stdio mode) get an `mcp-proxy` container as the entrypoint of the same pod.

**NetworkPolicies (default deny in `tool-gateway`).**

- `tool-gateway` namespace: default-deny ingress and egress, DNS egress.
- Gateway ingress: TCP 8080 from `hermes` pods labeled `app.kubernetes.io/name: onemind-hermes` only. (`hermes-legacy` is not a client.)
- Gateway egress: one rule per backend (the backend pod and port) and one rule per app Service. Public 443 only for the SaaS targets (GitHub,
  DigitalOcean, Cloudflare), with the private ranges excluded as in `onemind-hermes/network.yaml`.
- Each **app namespace** gets one ingress rule from the gateway or backend pod (the existing pattern: `platform/n8n/network.yaml`,
  `platform/karakeep/network.yaml`, `platform/onetak/media-mtx.yaml`, `platform/nats/network-policy.yaml`, `platform/values/argocd.yaml`; Grafana
  and Backstage need new ones). One file per app, so they parallelize.
- `onemind-hermes` egress (`platform/hermes/onemind-hermes/network.yaml`): **replace** the `onemind-ai/litellm:4000` MCP rule with
  `tool-gateway/agentgateway:8080`; keep DNS, Synapse and public 443. The LiteLLM rule goes when LiteLLM stops being the MCP hub (Phase 6).
- `execute-gate`: ingress only from the gateway pod; egress to Synapse (8008) and NATS (4222).

**Secrets (SealedSecrets, namespace-wide scope, files beside the workload).**

| Secret | Namespace | Keys (names) | Read by |
|---|---|---|---|
| `gateway-clients` | `tool-gateway` | `hermes`, optional `owner-console` | gateway |
| `gateway-cloud-credentials` | `tool-gateway` | `github`, `digitalocean`, `cloudflare`, `n8n-mcp`, `n8n-api`, `backstage`, `argocd`, `cloudtak`, `tak-portal`, `core-read`, `mediamtx` | gateway (injected into requests it forwards) |
| `mcp-<backend>-credentials` | `tool-gateway` | per backend (`grafana`, `clickhouse`, `karakeep`, `nextcloud`, `homeassistant`) | that backend only |
| `execute-gate` | `tool-gateway` | `matrix-token`, `owner-mxid` (not secret), `nats-creds` | gate only |
| `onemind-hermes-runtime` (existing) | `hermes` | add `ONEMIND_GATEWAY_KEY`; remove `ONEMIND_LITELLM_MCP_KEY` in Phase 6 | Hermes |

The Hermes `${ONEMIND_GATEWAY_KEY}` equals the `hermes` entry of `gateway-clients`, sealed twice from one generated value (never printed).

**Transition for `onemind-mcp`.** Gateway target `onemind` in the `ops` and `data` groups, bearer from a resealed `onemind-mcp-auth` token,
`platform/hermes/mcp/network-policy.yaml` gets an ingress rule from the gateway (replacing the LiteLLM one in Phase 6). Both hubs run in
parallel from Phase 1 to Phase 5, so Legacy never loses its current read adapters. Adapters are retired one by one after their replacement
passes its Phase 3 check (Prometheus, Loki, Tempo, Alertmanager by Grafana; Argo CD, Backstage by their MCP servers). Kept long term: NATS, STAC (if
OpenAPI is not worth it), Ollama, OpenCost, LiteLLM models, and the new write adapters.

**Observability wiring.** Metrics: PodMonitor (`release: monitoring`) on `:15020`; chart dashboard JSON loaded into Grafana by the existing
sidecar. Traces: OTLP to `tempo.monitoring.svc:4317`. Logs: JSON access log on stdout to Loki (per-call tool, target, client, outcome,
duration; **argument values are not logged**). Alerts: gateway down, backend 5xx ratio, `execute-gate` pending approval older than 10 minutes.
Audit of destructive attempts: bus `evt.toolgw.*` and Loki.

**Backstage and docs.** Components `agentgateway` and `execute-gate` in `platform/backstage/catalog/apps.yaml` (system `ai-agents`, owner
`group:default/onemind-admins`, annotations `argocd/app-name: platform-tool-gateway`, `backstage.io/kubernetes-namespace: tool-gateway`, Loki and
Grafana annotations as for `litellm`). A `docs/STACK.md` row under the platform tier (tier 2, namespace `tool-gateway`, state none, Git backup,
internal exposure). `platform/hermes/mcp/coverage.yaml` dispositions. `docs/operations/` runbook `tool-gateway.md` in Phase 1.

## 7. Rollout (phases and verification)

Conventions: "tool count" = number printed by `hermes mcp test <server>` inside `deploy/onemind-hermes -c hermes` (`hermes mcp list` shows
enabled servers). "Read call" = one real read through Legacy. "Safe write" = a reversible write on a dedicated test object; "none" means the
app has no safe write, so only the read is tested. Cleanup deletes are the **first destructive tests** (Phase 4).

| Phase | Deliverable | Verification |
|---|---|---|
| **0 Spike** (local binary, mock MCP servers, no cluster change) | Assumptions A1 to A6 and H1 (Home Assistant token) answered; findings written back to this doc | A1 an `allow` rule hides unlisted tools in `tools/list`; A2 backend header forms; A3 tool-name format seen by `hermes mcp test` against the local binary; A4 `$` escaping; A5 PDP outage returns an error, not an allow; A6 a dead target with `failOpen` keeps the others. H1 `curl` HA `/api/mcp` with a long-lived token (owner) |
| **1 Gateway core** | `tool-gateway` namespace, Application, standalone chart, client key, route `/ops/mcp` with the single target `onemind` (current adapters), PodMonitor, NetworkPolicies, Hermes `ops` server **added next to** the LiteLLM entry | pod Ready, `/healthz/ready`; `hermes mcp test ops` tool count equals what `hermes mcp test` reports today for the `onemind` server via LiteLLM (record both numbers); the same read call as via LiteLLM (`apps_list`) works; no key means 401; gateway metrics visible in Prometheus; record real CPU and memory |
| **2 Policy plane** | `classification/*.yaml` schema, `render_gateway_config.py`, `classify_tools.py`, CI drift check, generated values, route skeleton for all groups, `/execute/mcp` deny-all | `helm template` is clean in CI; renderer output is idempotent; a classified-tool fixture proves unlisted tool = not listed; calling any `/execute/mcp` tool returns denied |
| **3a..3g Backends, one app per PR** | per backend: Deployment/Service/NetworkPolicy/SealedSecret (or gateway-only target), the app-side ingress rule, `targets/<name>.yaml`, `classification/<name>.yaml`, route fragment | per app: `hermes mcp test <group>` tool count equals the exposed list; one read call through Legacy; one safe write where one exists; **one destructive tool shown as blocked** (not in `tools/list`, direct call returns denied). Safe writes: Kubernetes label on ConfigMap `tool-gateway/gateway-selftest`; Argo CD hard refresh of one app; Grafana annotation; GitHub issue comment on a test issue; DigitalOcean add a tag to a test resource; n8n create an inactive workflow `zz-gateway-test`; Karakeep bookmark with tag `gateway-test`; Nextcloud file in `legacy/selftest`; Home Assistant toggle a test `input_boolean`; NATS publish `evt.toolgw.selftest.ping`. None exist for Cloudflare (read-only token), ClickHouse (`readonly=1`), STAC, Backstage, MediaMTX, OpenCost, Ollama, LiteLLM |
| **4 Execute gate** | `services/execute-gate`, its build workflow, manifests, KV CR `exec-approvals`, Matrix bot and room, `extAuthz` wired on `/execute/mcp`, destructive tool lists populated | (a) call `resources_delete` on the self-test ConfigMap with no approval: denied and a request appears in `approvals`; (b) reply `execute <id>` as the owner: retry succeeds once; (c) retry again: denied; (d) change one argument: denied; (e) reply from another account or `override` on a destructive id: ignored; (f) stop `execute-gate`: calls denied; (g) a shell `curl` with the Hermes key fails the same way; (h) Hermes shows its own prompt first (`trust: untrusted`) and cron declines |
| **5 Hermes cutover** | `ONEMIND_GATEWAY_KEY` in `onemind-hermes-runtime`; all group servers and `execute` in the runtime config (via `hermes mcp add`, since edits to `config.defaults.yaml` do not touch a running instance) and mirrored into `config.defaults.yaml` for rebuilds; `network.yaml` egress; SOUL/skill text for the 403 behavior; Composio and OAuth servers added | `hermes mcp list` shows every group enabled; total eager tool count below the Tool Search listing budget (`tool_search` tier 1 or 2 as designed); an end-to-end task per group (health report from `ops`, issue from `dev`, query from `data`, entity lookup from `tak`, light toggle from `home`, bookmark from `work`) |
| **6 Retire LiteLLM-as-hub** | remove `mcp_servers` from `platform/ecosystem/values/litellm.yaml`, the `onemind-hermes-mcp` key, Secret `litellm-mcp-upstreams`, `ONEMIND_LITELLM_MCP_KEY`, the LiteLLM egress rule in `network.yaml`, the `litellm-clients` entry for Hermes, the LiteLLM ingress rule in the `onemind-mcp` policy; reduce `HUB.md` to a pointer; retire superseded `onemind-mcp` adapters | LiteLLM still serves models (a model call from `core`); Hermes has no MCP traffic to `onemind-ai`; `kubectl get netpol` matches Git; coverage CI green |
| **7 Optional** | `identity` group (read-only Authentik), Claude Code on the Mac through a tailnet ingress with its own client key, second gateway replica with `session.secretName` | per group verification as above |

### 7.1 What can run in parallel (separate agents, no shared files)

File ownership is the rule. The renderer assembles routes from per-target fragments, so backend PRs never edit a shared file.

| Phase | Owns (files) | Can start after | In parallel with |
|---|---|---|---|
| 0 Spike | `scripts/gateway/spike/` (scratch), doc edits to section 1 | now | 1 (design only), 4 (service code) |
| 1 Gateway core | `platform/argocd/platform-tool-gateway.yaml`, `platform/tool-gateway/{values.yaml,namespace,network,clients sealedsecret}`, `platform/hermes/mcp/network-policy.yaml` (+1 rule), `platform/hermes/onemind-hermes/` (not config), `STACK.md` row, coverage dispositions | 0 | 4 (code), 2 (design) |
| 2 Policy plane | `platform/tool-gateway/policy/`, `values.generated.yaml`, `scripts/gateway/*.py`, `.github/workflows/quality.yml` step | 1 merged (needs the chart in place) or developed against `helm template` | 3 (fixtures) once the schema is agreed |
| 3a Kubernetes | `platform/tool-gateway/backends/kubernetes/`, `targets/kubernetes.yaml`, `classification/kubernetes.yaml` | 1, 2 | every other 3x |
| 3b Argo CD + Grafana | `backends/argocd/`, `backends/grafana/`, `platform/values/argocd.yaml` rule, `platform/monitoring/` rule | 1, 2 | |
| 3c GitHub + DigitalOcean + Cloudflare (gateway-only targets) | `targets/{github,digitalocean,cloudflare}.yaml`, classification, `gateway-cloud-credentials` | 1, 2 | |
| 3d Data (ClickHouse, STAC, onemind-core, Tile38, WebODM, Ollama, TAK-ML) | `backends/clickhouse/`, `targets/`, adapters in `services/onemind-mcp/` | 1, 2 | |
| 3e TAK (TAK Portal, CloudTAK, TAK Server, MediaMTX, NATS publish) | `platform/onetak/` + `platform/cloudtak/` + `platform/takserver/` network rules, `targets/`, adapters | 1, 2 | |
| 3f Work (n8n, Karakeep, Nextcloud, Paperclip) | `platform/n8n/network.yaml`, `platform/karakeep/network.yaml`, `platform/nextcloud/`, `backends/{karakeep,nextcloud}/`, `targets/` | 1, 2 | |
| 3g Home (Home Assistant) | `targets/homeassistant.yaml`, edge ExternalName, `classification/homeassistant.yaml` | 0 (H1), 1, 2 | |
| 4 Execute gate | `services/execute-gate/`, `.github/workflows/build-execute-gate.yml`, `platform/tool-gateway/execute-gate/`, `platform/nats/streams/kv-exec-approvals.yaml` | 0 (A5), 1 | all of 3 (code); wiring to `/execute` after 2 |
| 5 Hermes cutover | `platform/hermes/onemind-hermes/{config.defaults.yaml,network.yaml,runtime.sealedsecret.yaml}`, SOUL/skill text | 3 (at least `ops`), 4 for `execute` | none that touch Hermes |
| 6 Retire | `platform/ecosystem/values/litellm.yaml`, `platform/ecosystem/resources/*`, `platform/hermes/mcp/*`, `HUB.md` | 5 verified | none |

Single-writer files (serialize PRs touching them): `docs/STACK.md`, `platform/hermes/mcp/coverage.yaml`, `PLATFORM-ROADMAP.md`, and
`platform/backstage/catalog/apps.yaml` (a 7,800-line file; one catalog PR per phase, not per backend).

## 8. Risks and open questions

- **Home Assistant auth (H1).** The documented `MCP Server` integration is OAuth with PKCE. If a long-lived token does not work as a Bearer,
  options are the community `ha-mcp` (long-lived token) or an HA add-on proxy. Decision after Phase 0.
- **One gateway, one failure domain.** A gateway restart drops every tool at once, as LiteLLM's restart did. Mitigation: `failOpen` per
  group (one dead backend does not break the rest), 2 replicas after the session key exists, PDB, and `onemind-mcp` staying reachable directly
  until Phase 6.
- **Classification is manual.** The drift check catches new or renamed tools, not a wrong decision. Review in Git is the control; the
  provider-side scopes in 3.1 bound the damage.
- **OpenAPI tools have no hints and no filtering.** Large specs (Authentik, n8n, CloudTAK) must be trimmed in Git; generated tool names are
  `operationId`-based and can collide or be too long for some model providers (check at Phase 3).
- **Env-expanded config.** The whole config text is expanded; a stray `$` in a CEL regex breaks startup (fails loudly, which is the safe
  direction). Keep secrets out of the rendered ConfigMap: only `${NAME}` placeholders are in Git and the ConfigMap, the values come from Secrets.
- **Approval fatigue and the typed word.** Two prompts for a destructive call. If it hurts, the Hermes layer can drop to `trust: full` on
  `execute` once Phase 4 is trusted; that is a one-line change and a conscious decision.
- **v1.7.0 is days old** (2026-10-27) and the MCP code is moving. Pin the digest, test the upgrade in Phase 0 style before bumping, and re-run
  the A1 to A6 spikes on every minor bump.
- **Roadmap decision log.** The 2026-10-18 entry left the tool gateway choice to Phase 8; this document records the choice (agentgateway), and
  the earlier "no separate Argo CD MCP" entry applied to Claude's read access, not to Legacy's write needs.
- **Hermes direct servers bypass the gate.** Composio and OAuth servers reach the internet directly. Scope them to read where the provider
  allows it, or accept `trust: untrusted` on them.

## Sources

- agentgateway: [site](https://agentgateway.dev/docs/standalone/latest/), [Kubernetes Helm install](https://agentgateway.dev/docs/kubernetes/latest/install/helm/),
  [MCP overview](https://agentgateway.dev/docs/standalone/latest/mcp/),
  [MCP configuration modes](https://agentgateway.dev/docs/standalone/latest/documentation/mcp/configuration-modes/),
  [MCP authorization](https://agentgateway.dev/docs/standalone/latest/documentation/mcp/mcp-authz/),
  [MCP target policies](https://agentgateway.dev/docs/standalone/latest/documentation/mcp/mcp-target-policies/),
  [releases](https://github.com/agentgateway/agentgateway/releases) (v1.7.0, 2026-10-27). Config schema, examples (`mcp-multiplex`, `mcp-openapi`,
  `mcp-authorization`, `mcp-authzen`, `mcp-authentication`, `mcp-telemetry`, `mcp-guardrails`), the `agentgateway-standalone` chart and the
  authorization tests were read from the `v1.7.0` source tarball.
- Hermes v0.22.5: `/opt/hermes/tools/mcp_tool*.py`, `tools/tool_search.py`, `hermes_cli/config_defaults.py`,
  `skills/autonomous-ai-agents/hermes-agent/references/native-mcp.md` (read-only in the running pod).
- MCP servers: [kubernetes-mcp-server](https://github.com/containers/kubernetes-mcp-server), [mcp-grafana](https://github.com/grafana/mcp-grafana),
  [mcp-for-argocd](https://github.com/argoproj-labs/mcp-for-argocd), [github-mcp-server](https://github.com/github/github-mcp-server),
  [mcp-digitalocean](https://github.com/digitalocean-labs/mcp-digitalocean), [Cloudflare MCP servers](https://github.com/cloudflare/mcp-server-cloudflare),
  [mcp-clickhouse](https://github.com/ClickHouse/mcp-clickhouse), [n8n instance-level MCP](https://docs.n8n.io/advanced-ai/mcp/accessing-n8n-mcp-server/),
  [Home Assistant MCP Server](https://www.home-assistant.io/integrations/mcp_server/), [Karakeep `@karakeep/mcp`](https://www.npmjs.com/package/@karakeep/mcp),
  [nextcloud-mcp-server](https://github.com/cbcoutinho/nextcloud-mcp-server), [mcp-proxy](https://github.com/sparfenyuk/mcp-proxy),
  [ha-mcp](https://github.com/homeassistant-ai/ha-mcp).
- Repo: `docs/STACK.md`, `PLATFORM-ROADMAP.md`, `platform/hermes/mcp/` (`HUB.md`, `apps.yaml`, `coverage.yaml`, `network-policy.yaml`),
  `platform/hermes/onemind-hermes/` (`config.defaults.yaml`, `network.yaml`, `README.md`), `services/onemind-mcp/README.md`,
  `platform/ecosystem/applications/ecosystem-litellm.yaml`.
