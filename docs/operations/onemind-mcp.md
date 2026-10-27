# OneMind application MCP for Legacy

User authorization: 2026-10-27, implement an extensible internal MCP and deploy it
for the cluster Hermes agent. Keep the existing services and their APIs. This is
a bounded addition, not a new cluster administration gateway or a replacement for
the app catalog, NATS, LiteLLM or Backstage's existing MCP.

## Source and activation status

- Source: `services/onemind-mcp/`, official Python MCP SDK 1.30.0, Streamable HTTP.
- Manifests: `platform/hermes/mcp/`, target namespace `hermes`. Use the dedicated
  `onemind-mcp-registry` image-pull Secret; do not rotate the shared registry credential.
- Initial tools: `apps_list`, `nats_server_stats`, `nats_jetstream_stats`,
  `prometheus_query`. App catalog resource: `onemind://apps`.
- 16 source/config tests pass; restricted container/official SDK smoke tests pass against
  **mock upstreams**, not live apps. Kubernetes manifests render and validate.
- **Deployed and live-verified**, 2026-10-27 04:00 EDT, from inside Legacy.
  [Acceptance run](https://github.com/onemindos/omos/actions/runs/40470531771)
  passed all four tools, resource reads, private authentication and rollout checks.
  MCP and Legacy each have one available replica; `platform-legacy` is Synced/Healthy.
- The deployed MCP uses the pinned image and a dedicated `onemind-mcp-registry`
  image-pull SealedSecret. The namespace's unrelated shared `ghcr` credential was
  left unchanged. The Legacy config contains only a runtime credential reference,
  with a private first-change backup.
- Only Legacy (`app.kubernetes.io/name=hermes,onemind.io/agent=legacy`) is the
  intended client. The separately isolated fresh Hermes is unchanged.

## Deploy through the existing Argo CD owner

1. Run the PR preflight on the existing `onemind-hermes-readonly` runner. Its
   existing scoped `pods/exec` identity reads only upstream availability, Hermes
   package/source feature indicators, and the SealedSecrets **public** certificate.
   It never reads config.yaml, environment credentials, user files or private keys.
   No new runner privileges or Kubernetes role are introduced. If the runner is
   absent/queued, report that gate instead of assuming cluster access exists.
2. Confirm the real endpoints: `nats-0.nats-headless.nats.svc.cluster.local:8222`
   and `monitoring-kube-prometheus-prometheus.monitoring.svc.cluster.local:9090`.
   Git records these names; the preflight/live calls must confirm reachability.
   NATS tools report one node and never claim a cluster aggregate.
3. The `OneMind internal MCP` workflow builds/pushes
   `ghcr.io/onemindos/onemind-mcp:sha-<source-sha>` on the feature branch and main.
   Obtain its `onemind-mcp-image` artifact and pin the returned `sha256` digest in
   the deployment. PR checks test without publishing. All actions are SHA-pinned.
4. Generate a random client token locally without printing it. Use the observed
   public controller certificate and `kubeseal` to create a namespace-scoped
   `SealedSecret/onemind-mcp-auth` in `hermes`, key `token`, alongside these
   manifests. Add it to the Kustomization. Commit only ciphertext. The server and
   Legacy consume this **same** namespaced Secret, avoiding cross-namespace copies.
   This task authorizes the new integration credential; unrelated credentials must
   not be read or rotated.
5. Add a third source to `platform/argocd/platform-legacy.yaml`, path
   `platform/hermes/mcp`, targetRevision `main`. Keep its existing sources and sync
   policy. Set the new SealedSecret's sync wave to -2 and MCP Deployment's to -1
   so the authenticated MCP is ready before Legacy's default wave. Argo owns all
   new resources; do not deploy them using an untracked `kubectl apply`.
6. Merge only the new MCP entry into Legacy's existing volume-managed config with
   the same atomic/preserve-other-keys pattern as `litellm_provider_config.py`.
   Add a dedicated config merge script/ConfigMap, init-container mount and client
   token environment reference through `platform/hermes/agent/`. Take a restrictive
   backup before editing; never export the full live config to CI or commit it.
   Keep unrelated providers, channels, skills and MCP servers unchanged.
7. Verify the **deployed** Hermes revision supports `${ENV_VAR}` in HTTP MCP headers
   and Tool Search. Current upstream documentation is not evidence for the pinned
   Legacy image. Enable `tools.tool_search.enabled: 'on'` only when supported. A
   small eagerly loaded MCP toolset is acceptable as a documented compatibility
   fallback; do not upgrade the whole agent as part of this integration silently.
8. Check the effective Cilium/Kubernetes network policy before sync. The supplied
   policy allows only the Legacy pod selectors to this MCP's port 8000 and allows
   MCP egress to DNS, NATS port 8222 and Prometheus port 9090. Add verified upstream
   pod selectors when observed. NetworkPolicies are additive; other policies can
   widen the effective allowance. Check upstream ingress and Legacy egress too.
9. Run server dry-run and review, then merge/sync via the existing GitOps owner.
   Legacy uses a Recreate strategy and one RWO volume; its configuration/env change
   restarts that agent. Do not delete the PVC, namespace or agent state.

## Client configuration

Merge with the live mapping; do not replace the configuration file:

```yaml
mcp_servers:
  onemind:
    url: http://onemind-mcp.hermes.svc.cluster.local:8000/mcp
    headers:
      Authorization: "Bearer ${ONEMIND_MCP_TOKEN}"
```

Legacy's pre-provisioned bearer credential (`ONEMIND_MCP_TOKEN` from
`onemind-mcp-auth`) is audited as client `hermes` and authorizes all configured read
tools. It is not browser OAuth or per-user authorization. The MCP has no Kubernetes
token mount or API privileges. Source supports a token per client
(`ONEMIND_MCP_TOKENS_FILE`, below); each client is audited by name but still gets
every read tool, so add separate scopes before enabling writes. Do not
expose this plaintext internal endpoint beyond the intended cluster trust boundary;
TLS/mTLS and an identity layer are prerequisites for wider exposure.

### Adding a client (not yet deployed)

The pinned image predates per-client tokens, and the deployment sets only
`ONEMIND_MCP_TOKEN` until an owner seals a client Secret. To add a client:

1. Generate one random token per client locally without printing it. Seal a
   namespace-scoped `SealedSecret/onemind-mcp-clients` in `hermes`, key `tokens`,
   holding a JSON object or `name:token` lines (format in
   `services/onemind-mcp/README.md`). Leave `hermes` out while `ONEMIND_MCP_TOKEN`
   is set. Commit only ciphertext; list it in `platform/hermes/mcp/kustomization.yaml`
   with sync wave -2.
2. In one commit, pin an image digest built from source with per-client tokens,
   mount the Secret read-only (for example at `/etc/onemind-clients`, keeping the
   default file mode because the container runs as uid 10001 without an fsGroup) and
   set `ONEMIND_MCP_TOKENS_FILE=/etc/onemind-clients/tokens`. Older images ignore it.
3. Allow the client's pods in `platform/hermes/mcp/network-policy.yaml` (today only
   Legacy may connect) and add the client's egress. A client in another namespace
   needs its own token sealed in that namespace.
4. Sync. Only the MCP pod restarts; Legacy keeps its token. Verify anonymous `/mcp`
   returns 401 and each client's calls appear in the audit log under its name.

To revoke or rotate one client, change its entry, reseal and roll the MCP; the
other clients keep working.

## Acceptance and rollback

Readiness is MCP process readiness, not upstream health. From the actual Legacy
pod, verify anonymous `/mcp` returns 401, then use the official SDK smoke script
with the runtime credential. Complete initialize/list-tools/resources, list apps,
read NATS node and JetStream stats, and run PromQL `up`. Check useful results and
the redacted tool/client/outcome/duration logs. Ask the real Hermes chat agent to
use the tools and confirm on-demand discovery if supported. Report each result;
neither a successful image build nor a healthy pod proves integration.

For rollback, restore only the narrowly added MCP config keys from the private
backup, revert the corresponding GitOps source/env/init changes and sync. Disable
the `onemind` MCP entry/reload to stop client use immediately. Keep the separate
fresh Hermes, app deployments, databases and PVCs unchanged. Rotate the new shared
token by resealing and rolling both consumers; no tokens appear in Git or logs.

## Adding applications

Use an existing maintained MCP directly from Hermes when it fits. For API-only
apps add a reviewed `register(add_tool, upstream)` module to
`services/onemind-mcp/onemind_mcp/adapters/` and its registry entry. Add a catalog
instance (unique `id`, description, `adapter`, `url_env`, optional `token_env`),
secret env refs and explicit network allowances. Existing adapters support another
app instance with catalog/env changes alone. Add success/failure/permission/output
tests, build a digest-pinned image, roll out and refresh Hermes discovery.

No arbitrary URL execution or automatic exposure of every Kubernetes Service/API
is provided. The current registrar is read-only: write tools need a separate
authorization/registration design plus idempotency and mutation-specific checks.

## Evidence log

Record image digest, source revision, Argo revision, pod rollout, verified Hermes
version/features and actual tool calls here when observed. Preflight run [40455726651](https://github.com/onemindos/omos/actions/runs/40455726651)
verified both upstreams from Legacy, Hermes 0.21.5 / MCP SDK 2.0.0, the Tool Search
module and actual nested HTTP header environment interpolation. The running
Hermes image digest matches the config init-container's pinned image.

Image candidate built from source `3c68986` in
[40455284585](https://github.com/onemindos/omos/actions/runs/40455284585):
`ghcr.io/onemindos/onemind-mcp@sha256:80fa67bc138006bd1efc0dc849504910b4068667036d880cba9d0c71afb762cc`.
Later changes add config integration/tests/docs; the server runtime source was
unchanged from that image until per-client tokens were added (not yet deployed).
The active SealedSecrets public certificate fingerprint
was `B8:21:29:16:FE:9B:B0:B7:4A:7B:29:BB:32:B6:F9:93:16:41:16:46:BC:F9:29:D4:69:15:E7:49:8C:9D:26:22`.
Live acceptance completed in
[40470531771](https://github.com/onemindos/omos/actions/runs/40470531771).
[Recorded result](onemind-mcp-acceptance-2026-10-27.json): all four tools and the app
resource pass; anonymous access returns 401; MCP generation 2 and Legacy generation
100 each have one available replica; Argo is Synced/Healthy. Progressive disclosure
is enabled in Legacy's config. Calls were made by the official SDK from the actual
Legacy container; this does not claim a separate end-to-end model/chat evaluation.

The private pull failure was resolved by #491 using repository Actions secret
`ONEMIND_MCP_GHCR_PULL_TOKEN`. Workflow `onemind-mcp-registry.yml` verifies access to
the pinned image and outputs only namespace-scoped SealedSecret ciphertext. To
rotate: update that repository secret, run the workflow, commit its sealed artifact
and sync GitOps. Track the token's chosen expiration in the password manager. The
registry token is used by Kubernetes image pulling, not loaded into the AI agent.
SDK 2 acceptance compatibility and fail-fast checks were completed by #492/#493.

## Broad app coverage expansion

The initial four-tool deployment covers NATS and Prometheus only. The expansion
adds searchable inventory and descriptions, STAC imagery collections/search,
Ollama installed/running models, active Alertmanager alerts, bounded Loki log
queries, OpenCost namespace allocations with MediaMTX stream status blocked by API authentication. This is 12
callable tools across seven configured apps, including the original integrations.
Configured is a deployment setting, not proof of upstream health. The acceptance
workflow checks each exposed tool without exporting application content.

`platform/hermes/mcp/apps.yaml` accounts for 32 application/access domains. Entries
without credentials, reviewed APIs or authorization stay **blocked**, are visible
to `apps_list`/`apps_describe`, and register no callable tools. This does not claim
that all applications are integrated. Nextcloud, Matrix, TAK HTTP, portal federation,
WebODM and other adapters remain explicit implementation work. Existing Legacy
bus/TAK tools remain separate. No user-file scraping or database backdoor is added.

Karakeep bookmark read/search/save, n8n workflow/execution summaries, Paperclip
company/issue read and backlog creation, Grafana dashboard search and Tempo trace
lookup have adapters in source. They are blocked pending the requirements recorded
on each catalog entry. Dedicated user/service identities are required; the registry
pull PAT, app admin credentials and model provider keys are not app API identities.
For n8n the transport uses `X-N8N-API-KEY`, not Bearer authentication.

Write tools are omitted unless the administrator explicitly sets `allow_writes:
true` on that app, supplies its dedicated credentials, and verifies network ingress
and end-to-end permissions. Write annotations are non-read-only/non-idempotent;
creation retries can duplicate records. The prepared deployment enables no writes.
Vault content, identity administration and physical actuation require separately
scoped authorization. "All apps" does not silently grant those privileges.

`coverage.yaml` maps all 56 Git-declared Argo Applications to their owning app or
supporting infrastructure disposition. The quality gate fails if a new Application
has no disposition, if an old one disappears without a coverage update, or if a
mapping points to a missing catalog entry. This is source inventory, not a count of
live Applications or ApplicationSet-generated children. Supporting Postgres,
Redis, operators, ingress and backup agents stay behind purpose-specific app APIs.

New app checklist: add the app access entry and coverage mapping, reuse a reviewed
adapter or implement typed operations, provision a dedicated least-privilege app
identity with SealedSecrets, allow both network directions only for the intended
pod/port, add boundary tests, pin the built image and run live acceptance from
Legacy. Optional/blocked entries must not make otherwise healthy tools unavailable.

Expansion preflight: STAC, Alertmanager, Loki and OpenCost returned HTTP 200 from
Legacy. MediaMTX returned 401 and stays blocked. Direct Legacy-to-Ollama timed out
under the gateway-only ingress policy; the expansion adds an MCP-pod-only ingress
exception for inventory reads, preserving the existing gateway rule. Legacy direct
access stays denied. Runtime acceptance from the MCP must confirm that exception.

The Legacy toolset-revision environment marker rolls its gateway after MCP readiness
so startup MCP discovery includes the new tools. Bump it when adding/removing tools.

Follow-through is tracked in [the complete app integration checklist](https://github.com/onemindos/omos/issues/543).
The expansion rollout serves the new catalog and tools; acceptance found that Loki's
Service port 80 forwards to gateway pod port 8080. The egress policy must allow
**8080**, scoped to Loki gateway pods. The pinned Loki 7.0.0 chart confirms this
mapping. Do not use Service ports as NetworkPolicy destination ports without
checking targetPort. Initial acceptance must wait for Legacy's toolset revision,
rather than testing its old startup cache while Argo is still rolling it.
