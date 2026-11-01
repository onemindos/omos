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
- Shared objects it depends on are owned by `platform-hermes-base` (path `platform/hermes/base`): the
  `hermes` Namespace and the Secrets `hermes-dashboard-auth`, `omos-nats-creds` and `hermes-nats-ca`. They
  used to belong to `platform-legacy`; retiring that Application no longer touches them
  (`platform/hermes/base/README.md`). The tool gateway's `onemind` target, `onemind-mcp`, is
  `platform-onemind-mcp` (path `platform/hermes/mcp`).
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

## Tool gateway (MCP)

Legacy gets its tools from the agentgateway tool gateway (`platform/tool-gateway`, runbook
[`docs/operations/tool-gateway.md`](../../../docs/operations/tool-gateway.md)): servers `ops`, `data` and `tak`,
tools named `<target>_<tool>`. The Deployment passes `ONEMIND_GATEWAY_KEY` from Secret `onemind-gateway-client`;
`config.defaults.yaml` references it as `${ONEMIND_GATEWAY_KEY}` in the `Authorization` header. Keep
`gateway.multiplex_profiles` off, or `${VAR}` in headers is not filled.

`config.defaults.yaml` only seeds an empty volume. To change the live `/opt/data/config.yaml` run
`mcp/apply.sh` (add `--dry-run` first): it sets `ops`, `data` and `tak` exactly as in the defaults, removes the
old LiteLLM `onemind` entry, writes `config.yaml.bak-mcp-<timestamp>` before changing anything and touches no other
key. A second run prints only `unchanged`. Hermes reads `mcp_servers` at start, so restart the pod afterwards.
Check with `hermes mcp test ops|data|tak` inside the pod.

## Permanent Zoho MCP connections

The nine owner-provided Zoho endpoints are managed through strict-scoped
SealedSecrets and environment-only URL references. A Git-owned startup reconciler
uses the supported Hermes CLI to add only these nine entries on existing as well
as clean volumes, preserving other integrations. Initial exposure is an explicit
read-only allowlist, not all upstream actions. See [Zoho setup, verification and
recovery](mcp/zoho/README.md). Source preparation and live acceptance are separate.

## Permanent hosted Whop MCP

The existing hosted Whop MCP definition reconciles through the supported Hermes
config CLI on every start, scoped only to `mcp_servers.whop`. Private OAuth state
remains on the existing data PVC; its encrypted backups, not Git token snapshots,
are the authentication recovery path. No standalone Whop CLI install, startup
configuration or login is included. See [Whop MCP startup and recovery](mcp/whop/README.md).
Source tests do not establish live backup/restore success or authorize financial
actions.

## Native 1Password app credentials

**Source-only draft — NOT DEPLOYED.** Proposed native-source settings, pinned
`op` installer, reconciliation, runtime correction and tests remain under
[`onepassword/`](onepassword/README.md). Production Deployment, defaults and
Kustomization are byte-identical to the original base: no bootstrap injection,
installer, ConfigMap or native-source activation is added. The existing runtime
image is unchanged. Activation requires a built/tested corrected image pinned by
its real digest, a real rotated strictly scoped bootstrap SealedSecret, reviewed
GitOps wiring and authorized runtime/recovery acceptance. No new MCP/Connect
server is introduced.

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
| `plugins.enabled` | `[]` | No custom plugins in the seed. The apply scripts add `legacy-roles` (`roles/apply.sh`) and, at the NATS cutover, `nats-platform` (`nats/apply.sh`) |
| `mcp_servers` | `ops`, `data`, `tak` | The tool gateway routes (`http://agentgateway.tool-gateway.svc.cluster.local:8080/<group>/mcp`), key `Authorization: Bearer ${ONEMIND_GATEWAY_KEY}`. Eager connect (`lazy` unset), `trust` default. `/execute` is not added until the execute gate is deployed |

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
- egress on TCP 8080 to the tool gateway pod in `tool-gateway`
  (`app.kubernetes.io/name: agentgateway-standalone`), the MCP front door; its own
  NetworkPolicy admits this pod;
- egress on TCP 4222 to the NATS pods in `nats` (a direct client of the bus, see
  [NATS](#nats-legacys-bus-presence)); the `nats-clients` policy in
  `platform/nats/network-policy.yaml` admits this pod;
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

## Scheduled jobs

Four cron jobs run the daily brief, two cheap watchdogs and the nightly config snapshot.
Everything is in Git under `cron/` and applied by one idempotent script. Nothing is created
by hand in the pod.

| Job | Schedule (UTC) | Mechanism | Model |
|---|---|---|---|
| `legacy-sovereign-brief` | daily 07:00 | script gathers facts, model writes <180 words | `ollama-cloud` / `glm-5.3`, low effort |
| `infra-health` | every 15m | `--monitor-script`: silent unless the output changes | same, only on a change |
| `security-scan` | daily 06:30 | `--monitor-script`, same contract | same, only on a change |
| `config-export` | nightly 02:30 | `--no-agent`: the script is the job | none. **Paused**, see below |

All deliver to Telegram `8297858681`. The model jobs use the already-paid Ollama Cloud
subscription (the `ollama-pro` credential in the pool; verified with a one-token call), not
Codex, and never spend Codex quota. The old `so-infra-health` job (every 30 minutes, full model
run) hit `429 tokens_per_day`; here a healthy cluster costs zero tokens: the monitor script's
output is `OK`, its hash is unchanged, and Hermes skips the model. A problem changes the output,
wakes the model once, and stays silent until the output changes again (a new problem, or
recovery). The very first run after creation also counts as a change; the prompts answer
`[SILENT]` when everything is already `OK`.

Files:

| Path | Role |
|---|---|
| `cron/jobs.yaml` | Declarative job list: schedule, prompt, skills, model, delivery, `paused` |
| `cron/apply_jobs.py` | Applier, runs in the pod as `hermes`, through Hermes' own `cron.jobs` API |
| `cron/apply.sh` | Host wrapper: `cron/apply.sh [--dry-run]`. `--install-local` first copies this checkout's scripts into the pod, for testing before merge |
| `cron/scripts/legacy_common.py` | MCP client and Prometheus/Alertmanager/Argo helpers (stdlib only) |
| `cron/scripts/legacy_brief_facts.py` | Facts for the brief: nodes, bad workloads, volumes, alerts now and in 24h, Argo drift, cost, NATS, Codex usage |
| `cron/scripts/legacy_codex_usage.py` | Codex usage per credential (below) |
| `cron/scripts/legacy_infra_health.py`, `legacy_security_scan.py` | Monitor sources, print `OK` or one stable line per problem |
| `cron/scripts/legacy_config_export.py` | Snapshot, secret scan, commit and push |

How they reach the pod: `kustomization.yaml` builds the ConfigMap `onemind-hermes-cron` from
these files, and the `initial-config` init container copies them on **every** start (unlike
`config.yaml` they are Git-owned): scripts to `/opt/data/scripts` (the only directory the cron
runner accepts, symlinks out of it are refused) and `jobs.yaml` plus `apply_jobs.py` to
`/opt/data/cron-defs`. Changing a script therefore rolls the pod once. Jobs themselves live in
`/opt/data/cron/jobs.json` on the PVC (so VolSync backs them up) and are created only by
`cron/apply.sh`. A second apply prints `unchanged` for every job; fields that drifted are
updated, and `paused` in `jobs.yaml` is enforced both ways. Jobs not in the file are left alone.

Data sources and network: the pod has no Kubernetes token and egress only to DNS, the tool gateway
(`tool-gateway`, TCP 8080), Synapse and public 443. So the scripts call the read-only `onemind`
tools on the gateway's `/ops` route with `ONEMIND_GATEWAY_KEY` (`onemind_prometheus_query`,
`onemind_alertmanager_alerts`, `onemind_argocd_applications`, `onemind_nats_*`,
`onemind_opencost_allocation`; the gateway is stateful, so `legacy_common.py` runs
initialize and notifications/initialized before each call). Verified: Hermes' script runner passes that
variable to scripts (`build_subprocess_env`). The export's SSH goes to `ssh.github.com:443`,
which the public 443 rule already allows; the export's SSH goes to `ssh.github.com:443`,
which the public 443 rule already allows. Prometheus has no cert-manager or Falco metrics, so
the security scan reads Falco and certificate findings from Alertmanager only.

Approvals: `approvals.cron_mode: deny` governs commands an agent runs through its terminal tool.
The scripts run in the cron runner, not through that path, and the job prompts tell the model not
to run commands (the facts are already in the prompt), so nothing in these jobs can need a flagged
command. A flagged command would be denied, not asked.

### Codex usage per credential

`hermes usage --json` only reports the active credential. In Hermes v0.22.5
`agent/account_usage.fetch_account_usage(provider, base_url=, api_key=, read_only=)` queries
the ChatGPT usage endpoint for an explicit token, and `hermes_cli.auth.read_credential_pool()`
lists the pool entries (`codex-1`, `codex-2`). `legacy_codex_usage.py` calls the first for each
entry with `read_only=True`, the variant the model picker uses: a 401 is reported, never
repaired, and nothing is refreshed, rotated or written. Tokens are never printed. Output:

```
codex-1 (Plus): Session 0% (resets 2026-10-30 19:11 UTC), Weekly 31% (resets 2026-11-04 10:04 UTC)
codex-2 (Plus): Session 0% (resets 2026-10-30 19:45 UTC), Weekly 3% (resets 2026-11-04 15:00 UTC)
```

`--json` prints the full per-credential record. The brief script embeds it.

### config-export and the Git credential it needs

`legacy_config_export.py` writes an explicit allowlist, not `hermes profile export` (a 12 MB
tarball of the whole home, with sessions, hub caches and bundled skills): `config.yaml` with
secret-looking keys blanked and `${ENV}` references kept, `SOUL.md`, `memories/MEMORY.md` and
`USER.md`, `skills/index.json` (name, category, origin: bundled, user or external), the contents
of user-created skills, and the job definitions without run history. The whole snapshot is then
scanned: token shapes, private-key headers, credential assignments, and every secret value this
pod holds (environment and `auth.json`) verbatim. One hit aborts the run, names the file and rule
(never the value) and commits nothing. Tested here with a fake key and token (aborted) and
against a local bare repository (first run commits, second run is silent).

onemind-hermes has **no Git write credential** (no token variable, no SSH key, nothing under
`HOME`), so `config-export` is created **paused** with that reason, and the pod mounts an
optional Secret for it. Owner steps to enable it:

1. Create a private repository for the snapshots, for example `onemindos/legacy-snapshots`.
   (Keeping agent commits out of `onemindos/omos` is deliberate; the snapshot contains memory
   text, so it must be private.)
2. Create a **deploy key with write access on that one repository** (preferred over a token:
   it cannot touch any other repository). Generate it locally,
   `ssh-keygen -t ed25519 -N '' -C legacy-config-export -f ./legacy-export-key`, add the `.pub`
   under the repository's Settings, Deploy keys, "Allow write access". Restrict the branch with a
   ruleset if wanted (default branch `main`).
3. Seal it, namespace-wide like the other onemind-hermes secrets, into
   `platform/hermes/onemind-hermes/export.sealedsecret.yaml`:

   ```bash
   kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
     --controller-namespace=kube-system > /tmp/pub-cert.pem
   kubectl create secret generic onemind-hermes-export -n hermes \
     --from-file=ssh-privatekey=./legacy-export-key \
     --from-literal=repo=ssh://git@ssh.github.com:443/onemindos/legacy-snapshots.git \
     --dry-run=client -o json \
     | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/hermes/onemind-hermes/export.sealedsecret.yaml
   ```

   Add the file to `kustomization.yaml`, then delete the local key. An optional `branch` key
   (default `main`) can be added with `--from-literal=branch=...`. A fine-grained token over
   HTTPS would also work but needs the script changed; the deploy key needs no code change.
4. After Argo syncs (secret volumes are read at pod start, so the pod restarts once if it was
   already running), set `paused: false` for `config-export` in `cron/jobs.yaml`, merge, and
   run `cron/apply.sh`.

Until then the job stays paused; if it is ever run without the Secret it fails loudly rather
than silently skipping a backup.

## How Legacy works: doctrine and worker roles

One agent, one soul. GE, HP, LE and SO are modes of Legacy, not profiles. When work is big or parallel she
dispatches five worker roles on Codex and a `checker` verifies their output before she reports it. The doctrine
is the "How you work" section of `platform/hermes/seed/onemind-hermes/SOUL.md` (read every turn, about 400 more
tokens than before); the role definitions are `roles/roles.yaml`.

| Role | Model (`openai-codex`) | Toolsets | Job |
|---|---|---|---|
| `ops-cluster` | `gpt-6.1-sol` | terminal, file, web, skills, todo, session_search | cluster, Argo, K8s, NATS, Tailscale ops |
| `tak-mission` | `gpt-6.1-sol` | same | TAK, missions, CoT, video |
| `builder` | `gpt-6.1-sol` | + code_execution | code, manifests, branches, PRs (never merges) |
| `researcher` | `gpt-6.1-terra` | web, file, skills, todo, session_search | web and docs research |
| `checker` | `gpt-6.1-terra` | terminal, file, web, skills, todo, session_search | verifies other workers' output |

Models are those `provider_models_cache.json` lists for `openai-codex` in the pod. `terra` is the cheaper tier
(about half of `sol` per token in the models.dev cache); `luna` is cheaper still but too weak for review.

### What Hermes v0.22.5 actually supports (read from `/opt/hermes`)

- **`delegate_task` has no named roles.** The bundled `default` backend gives every child one route: the
  `delegation.provider` / `delegation.model` section or the parent's own (`plugins/delegation/default/__init__.py`
  `route_task`). `provider`, `model`, `role` and `toolsets` are reserved task names and not model-facing
  (`agent/delegation_provider.py` `RESERVED_TASK_PARAMETERS`; `role` is ignored, capability is depth-derived,
  `tools/delegate_tool.py:242`). A pluggable backend can add an enum parameter and answer with a route, including
  `profile` = run the child as that Hermes profile (`tools/delegate_tool_routes.py` `_profile_route`: the
  profile's own config, model, SOUL and toolsets). So a per-role model needs a backend plugin, and this PR ships
  one: `roles/plugin/legacy-roles/` adds `worker` (enum of the five) to every task; no `worker` = the default
  route. It is a user plugin (`$HERMES_HOME/plugins/`, `kind: backend`, gated by `plugins.enabled`,
  `hermes_cli/plugins_discovery.py` `gate_manifest`), selected by `delegation.backend: legacy-roles`.
- **`hermes kanban swarm` roles are profiles.** `--worker PROFILE:TITLE[:SKILL,..]`, `--verifier PROFILE`,
  `--synthesizer PROFILE` (`hermes_cli/kanban_parser.py:222`, `kanban_swarm.py`); the card's assignee must satisfy
  `profile_exists` or the dispatcher will not spawn it (`kanban_db_dispatch.py:2077`, `2329`). The worker runs
  `hermes -p <profile>` with that profile's `HERMES_HOME`, model and `platform_toolsets.cli`
  (`_default_spawn`, `_resolve_worker_cli_toolsets`). The dispatcher is embedded in the gateway (60 s tick; gateway
  log "kanban dispatcher: embedded in gateway"). `kanban swarm` is CLI-only (no tool), so Legacy runs it through
  her terminal tool.
- **Profiles share the Codex pool; no second sign-in.** `read_credential_pool` returns the default profile's
  `credential_pool` whenever the profile has zero entries for that provider (`hermes_cli/auth.py:948`, global-root
  `auth.json` read-only fallback `:507`); `_load_provider_state` does the same for OAuth state (`:836`), and a
  refresh that rotates the grant writes it back to the root store (`agent/credential_pool.py:724`). So role
  profiles carry no `auth.json`, and `credential_pool_strategies: {openai-codex: least_used}` in their own
  config spreads turns across both accounts. (`hermes profile create --clone` would copy the grant and fork it;
  `strip_cloned_single_use_oauth_grants` exists to undo that. The applier never clones.)
- **Where the soul comes from.** `/opt/data/SOUL.md`, a plain file. It is identical to
  `platform/hermes/seed/onemind-hermes/SOUL.md` (checked before this change); the init container does not seed it,
  `config.defaults.yaml` only seeds `config.yaml`, and `MIGRATION.md` step 2 copied it by hand. A role profile's
  `SOUL.md` is its whole brief (`roles.yaml`), so a role has no identity of its own.
- **Skills** are shared read-only: each role sets `skills.external_dirs` to `/opt/data/skills` and
  `/opt/data/onemind-skills`, so swarm cards that name a skill (`requesting-code-review`, `humanizer`) resolve
  without copying 56 MB per role.

### Applying it (idempotent)

`config.defaults.yaml` only seeds an empty volume, so a running instance gets this with one script:

```bash
platform/hermes/onemind-hermes/roles/apply.sh --dry-run   # shows create/update/unchanged
platform/hermes/onemind-hermes/roles/apply.sh
```

It pushes `roles/` and the repository `SOUL.md` to `/opt/data/roles-defs`, then runs `apply_roles.py` as `hermes`:
installs `SOUL.md` and the `legacy-roles` plugin; creates the five profiles (`--no-alias --no-skills`, no
credentials) with their model, `credential_pool_strategies`, `platform_toolsets.cli`, `skills.external_dirs`, and
the default profile's `approvals`, `security`, `agent.disabled_toolsets` and `mcp_servers` (copied raw, `${VAR}`
unexpanded, so it follows whatever the default profile has); sets `delegation.backend: legacy-roles` and adds
`legacy-roles` to `plugins.enabled`. Anything it overwrites is copied first to
`/opt/data/backups/legacy-roles-<UTC stamp>/`; a second run prints `unchanged` for everything. It never reads or
writes `auth.json` or any `.env`.

The `delegation` and `plugins` keys are deliberately **not** in `config.defaults.yaml`: on a fresh volume the
backend would be selected before its plugin exists, and an unavailable backend fails `delegate_task` closed.
After applying, the plugin loads at the next gateway start, and a new session (`/new`) is needed for the new soul
and tool schema. Rollback: `hermes config set delegation.backend default`, and restore `SOUL.md` from the backup
directory.

Workers run with `approvals.mode: manual`, which has no human in a worker session: a flagged command times out
and is denied, which is the safe direction. Workers get no authority Legacy lacks.

## NATS: Legacy's bus presence

Legacy is a **direct client of the cluster bus**, like every other cloud workload: no leaf sidecar (the
`hermes-legacy` sidecar is the one documented exception, `CLAUDE.md` "Bus Authority"; do not copy it).
Inventory of what Legacy does on NATS, what was ported and the switch order are in
[`MIGRATION.md`](MIGRATION.md#nats-inventory-and-cutover-step-7). Summary:

- Identity `agents.{prompt,status,hb}.hermes.zeus.legacy` through the Synadia agent protocol plugin
  `nats-platform` (vendored in `platform/hermes/agent/plugins/nats-platform/`; no start-up hot-patch, because
  Hermes v0.22.5 has no native NATS platform and the one line the patch fixed is now in the vendored file).
- Connection settings are in `deployment.yaml` (URL `tls://nats.nats.svc.cluster.local:4222`, user and password from
  Secret `omos-nats-creds`, CA from Secret `hermes-nats-ca`, identity, `HERMES_NATS_REFUSE_DUPLICATE=true`). They are
  inert until the plugin is listed in `plugins.enabled`.
- Python packages come from the `install-gateway-deps` init container (hash-pinned in `nats/gateway-deps.txt`,
  installed once to `/opt/data/.local/gateway-deps`).
- **Applying** (idempotent, backs up what it overwrites, restarts only the gateway process on request):

  ```sh
  platform/hermes/onemind-hermes/nats/apply.sh --status
  platform/hermes/onemind-hermes/nats/apply.sh --dry-run
  platform/hermes/onemind-hermes/nats/apply.sh --restart          # nats-platform only
  platform/hermes/onemind-hermes/nats/apply.sh --with-context     # also nats-context (off by default)
  platform/hermes/onemind-hermes/nats/apply.sh --disable --restart   # rollback switch
  ```

- While `hermes-legacy` still holds the identity, onemind-hermes connects, sees it, and waits instead of
  registering; it takes the identity over by itself once Legacy lets go. That makes enabling it first safe.
- `tak-platform` is not ported; PR #585 replaces it with a TAK-native client.

## TODO

- Ollama Cloud fallback: seal `OLLAMA_API_KEY` into `onemind-hermes-runtime` (or
  a new namespace-wide SealedSecret) before relying on the fallback.

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
