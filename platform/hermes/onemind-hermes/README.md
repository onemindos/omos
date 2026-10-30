# onemind-hermes

Legacy (`onemind-hermes`): the new OneMind Hermes instance, running beside the current `hermes-legacy`
Deployment. It replaces the `hermes-fresh` launch of October 26, 2026 (same
image, same design, new names and new storage). It does not import any existing
agent data.

**Agent persona: Legacy, a single soul.** The pods carry
`onemind.io/agent: legacy`; the Deployment, Service and storage are named
`onemind-hermes`. The four profile souls from the old installation are not
separate agents here.

- Argo Application: `platform-onemind-hermes`, path `platform/hermes/onemind-hermes`.
- Deployment/Service/Ingress/NetworkPolicy: `onemind-hermes`, namespace `hermes`.
- Private dashboard: `https://onemind-hermes.tail717ba0.ts.net` once the
  Tailscale operator provisions it. Use the existing Hermes dashboard username
  and password.
- Image: official Hermes v0.22.5, pinned to the verified multi-architecture digest.
- New 25 GiB PVC `onemind-hermes-data` on `do-block-storage-retain`, so deleting
  the claim keeps the volume. No old volume is mounted or cloned; the unused
  `hermes-fresh-data` claim is not carried over.
- Initial config selects the existing preferred model/provider, but **no model
  credentials are imported**. Sign in to the model on the new instance.
- Native bundled skills/default features remain available. No custom plugins,
  MCP servers, channels, schedules, Legacy memory or conversation history are
  copied. `MIGRATION.md` lists what moves over and in which order.

## Labels and selectors

`app.kubernetes.io/name: onemind-hermes` is the selector for the Deployment,
Service and NetworkPolicy; it keeps Legacy's selectors
(`app.kubernetes.io/name: hermes` plus `onemind.io/agent: legacy`) separate.
Because the pods also carry `onemind.io/agent: legacy`, anything that selects on
that label alone matches both Deployments. Today that is only the
`hermes-legacy` PodDisruptionBudget (`maxUnavailable: 1`), which then spans both
pods. Keep new selectors on `app.kubernetes.io/name`.

## Initial configuration (`config.defaults.yaml`)

The init container copies it once into the empty PVC with exclusive creation.
Later pod starts keep configuration authored through Hermes, so changing this
file after first start has no effect on the running instance. Every key was
checked against the v0.22.5 code in the image and passed its
`validate_config_structure`:

| Key | Value | Why |
|---|---|---|
| `model.provider` / `model.default` | `openai-codex` / `gpt-6.1-sol` | Same route as Legacy |
| `credential_pool_strategies.openai-codex` | `least_used` | Spread turns across signed-in Codex credentials instead of draining the first (`fill_first` default) |
| `fallback_providers` | `ollama-cloud` / `glm-5.3` | Used only if the primary fails. Needs `OLLAMA_API_KEY`, not provisioned yet (see TODO) |
| `agent.disabled_toolsets` | `computer_use`, `browser`, `discord`, `discord_admin`, `yuanbao`, `feishu_doc`, `feishu_drive`, `x_search`, `video_gen` | Desktop, browser, chat-platform and media tools a cloud agent does not need |
| `approvals.mode` | `manual` | Dangerous commands wait for a human instead of the LLM judge (`smart` default) |
| `approvals.cron_mode`, `single_query_mode`, `unattended_mode` | `deny` | Cron, one-shot and API-server sessions never run flagged commands (matches the default; pinned) |
| `security.redact_secrets` / `allow_private_urls` | `true` / `false` | Defaults, pinned |
| `plugins.enabled` / `mcp_servers` | `[]` / `{}` | MCP is added later, through LiteLLM |

`agent.disabled_toolsets` is applied last and per tool. It lists `yuanbao` (the
Yuanbao platform tools) rather than `hermes-yuanbao`: the latter is the whole
Yuanbao platform composite and includes terminal, file, web, skills, delegation
and kanban tools, so disabling it would remove toolsets this instance keeps.

The terminal backend stays `local`. The container is the sandbox: no
service-account token, no privilege escalation, `NET_RAW` dropped, and the
NetworkPolicy below. Hermes' Docker, Singularity, SSH, Modal and Daytona backends
either do not run inside a pod or send commands to external hosts.

## Isolation and network

The pod has no Kubernetes service-account token, NATS sidecar, observer
credential or Paperclip credential. The NetworkPolicy allows:

- ingress on 9119 (dashboard) from the `tailscale` namespace;
- egress to cluster DNS;
- egress on TCP 4000 to the LiteLLM pods in `onemind-ai`
  (`app.kubernetes.io/name: litellm`), the planned MCP hub. LiteLLM's own
  `litellm-clients` ingress policy does not admit this pod yet; it needs a rule
  for `app.kubernetes.io/name: onemind-hermes` in namespace `hermes` before
  traffic flows;
- egress on 443 to public addresses, with private, cluster and tailnet ranges
  excluded.

The API is keyed independently and is not published by the Service or Ingress.
Dashboard session signing uses its own secret. Only the existing
`hermes-dashboard-auth` login Secret is shared. Plaintext credentials are not in
Git.

### Sealed secrets

Both SealedSecrets are namespace-wide scoped
(`sealedsecrets.bitnami.com/namespace-wide: "true"`), so they are bound to the
`hermes` namespace but not to a name. They were renamed by changing only
`metadata.name`, and the encrypted values were reused:

- `onemind-hermes-runtime` (was `hermes-fresh-runtime`): `API_SERVER_KEY`,
  `HERMES_DASHBOARD_BASIC_AUTH_SECRET`.
- `onemind-hermes-restic` (was `hermes-fresh-restic`): Restic storage
  credentials and password unchanged; `RESTIC_REPOSITORY` re-sealed with the new
  prefix (below) using `kubeseal --merge-into`, so no other value was decrypted.

Deployment uses one replica and `Recreate`. PVC annotations prevent Argo
prune/delete of its data.

## Backup

`ReplicationSource onemind-hermes-data` snapshots the PVC nightly at 01:25 UTC
using the existing VolSync mover/namespace policy, with explicit
`do-block-storage` snapshot and cache classes (temporary claims). Retention
follows Legacy's 7 daily / 4 weekly / 3 monthly / 1 yearly policy.

Its Restic repository is a new, independently initialized prefix,
`backups/restic/hermes/onemind-hermes-data` in the `onemind-backups` bucket. That
is the `backups/restic/<namespace>/<pvc>` layout every other VolSync source uses
and that `platform/backups/tools/mirror.py` expects. It is not shared with the
`hermes-fresh` repository or with Legacy's `hermes-data` repository. It shares
the existing storage credentials and recovery password; preserve the existing
off-cluster recovery password. First-backup completion and a restore test must
be observed; this source configuration alone does not establish recovery.

## Launch acceptance and gradual migration

After Argo reconciles, record the actual rollout, gateway health, dashboard auth
and tailnet access in the deployment PR. A ready pod is not proof of working
chat; model sign-in and a real response are a separate acceptance step.

Add selected memory, skill packages and integrations incrementally, following
`MIGRATION.md`. Inspect credential-bearing assets before importing them. Keep
Legacy running until the owner chooses the final cutover.

## TODO

- Ollama Cloud fallback: seal `OLLAMA_API_KEY` into `onemind-hermes-runtime` (or
  a new namespace-wide SealedSecret) before relying on the fallback.
- LiteLLM ingress: add the `onemind-hermes` rule to `litellm-clients` in
  `onemind-ai`.

## Cleanup after cutover

Removing `platform-hermes-fresh` from Git made the app-of-apps delete that
Application. It had no finalizer and `prune: false`, so its resources were
orphaned and then removed by hand.

Done (2026-10-29), with owner approval:

- Deployment, Service, Ingress and NetworkPolicy `hermes/hermes-fresh`, and the
  tailnet proxy for the `hermes-fresh` node.
- ConfigMap `hermes/hermes-fresh-defaults-*`.
- ReplicationSource `hermes/hermes-fresh-data` and its cache claim
  `volsync-src-hermes-fresh-data-cache`.
- SealedSecrets and Secrets `hermes/hermes-fresh-runtime` and
  `hermes/hermes-fresh-restic`.

Still pending:

- PVC `hermes/hermes-fresh-data` (`do-block-storage`, Delete reclaim policy; no
  sessions were ever recorded on it). Owner decision: keep for now. Do not
  delete it without a fresh snapshot and explicit owner approval.
- Remove the `hermes-fresh` device from the tailnet admin console.
- The `hermes-fresh` Restic prefix in the `onemind-backups` bucket, after the
  above.

## Agent name

The agent is called **Legacy** everywhere she is shown. `onemind-hermes` is only the
instance and infrastructure name. Hermes v0.22.5 has no agent-name config key, so the
name is set in three places: `SOUL.md` (identity), `API_SERVER_MODEL_NAME: Legacy` in
`config.yaml` (the id on `/v1/models`) and the Matrix display name of
`@legacy:onemindos.dev` (set once through the Synapse profile API). The Telegram bot
(`omos_hermes_bot`, first name `omos-hermes`) can only be renamed in BotFather. The old
`hermes-legacy` Deployment is the old Legacy instance, being retired.
